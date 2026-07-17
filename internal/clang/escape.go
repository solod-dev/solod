package clang

import (
	"go/ast"
	"go/token"
	"go/types"
)

// Escape analysis
//
// So has no garbage collector and no heap by default. Several builtins allocate
// their result in the current function's stack frame: string concatenation, rune
// and []rune conversions, make for slices and maps, array literals, and taking
// the address of a local, etc. Such a value is a "frame value": it is valid only
// until the function returns. Returning one hands the caller a dangling pointer.
//
// This checker rejects any return of a frame value. It runs per function, in two
// passes over the body:
//
//  1. markFrameVars computes the set of locals that hold a frame value.
//  2. escapes scans every return and flags the results that are frame values.
//
// The analysis is intraprocedural. A call to a user or stdlib function is opaque:
// its result is assumed not to be a frame value (make, append and the builtin
// conversions are the known exceptions, handled explicitly). It does not descend
// into nested closures, which have their own frame.

// escapeMsg is the diagnostic reported for a return that escapes the frame.
const escapeMsg = "stack-allocated value escapes function frame"

// rejectEscapes fails the build on the first return that escapes the frame.
// Generic functions expand to macros inlined into the caller, so they never
// reach here (see emitMacroFuncDecl).
func (g *Generator) rejectEscapes(decl *ast.FuncDecl) {
	for _, node := range findReturnEscapes(g.types, decl) {
		g.fail(node, "%s", escapeMsg)
	}
}

// findReturnEscapes reports the return-value expressions in decl that escape the frame.
func findReturnEscapes(info *types.Info, decl *ast.FuncDecl) []ast.Node {
	if decl.Body == nil {
		return nil
	}
	c := &escapeChecker{
		info:   info,
		locals: map[types.Object]bool{},
		fvars:  map[types.Object]bool{},
	}
	c.collectLocals(decl)
	c.markFrameVars(decl.Body)
	return c.escapes(decl)
}

// escapeChecker holds the per-function analysis state.
type escapeChecker struct {
	info   *types.Info
	locals map[types.Object]bool // every object declared in the function
	fvars  map[types.Object]bool // locals that hold a frame value
}

// collectLocals records every object declared inside the function: receiver,
// parameters, results and body-local variables.
func (c *escapeChecker) collectLocals(decl *ast.FuncDecl) {
	ast.Inspect(decl, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok {
			if obj := c.info.Defs[id]; obj != nil {
				c.locals[obj] = true
			}
		}
		return true
	})
}

// markFrameVars fills fvars with the locals that hold a frame value,
// It iterates to a fixpoint so it follows assignment chains such as
//
//	t := a + b; s := t.
//
// Marking is monotonic, so the loop terminates.
func (c *escapeChecker) markFrameVars(body *ast.BlockStmt) {
	for {
		changed := false
		c.walk(body, func(n ast.Node) {
			switch s := n.(type) {
			case *ast.Ident:
				changed = c.markArrayDecl(s) || changed
			case *ast.AssignStmt:
				changed = c.markAssign(s) || changed
			case *ast.DeclStmt:
				changed = c.markDecl(s) || changed
			}
		})
		if !changed {
			return
		}
	}
}

// markArrayDecl marks a local array variable, which is itself a frame value: an
// array lowers to a bare C array, so returning it by value returns a pointer into
// the frame. It acts on the defining ident, which lives in the body, so it skips
// parameters (an array parameter is already a caller pointer, safe to return).
func (c *escapeChecker) markArrayDecl(id *ast.Ident) bool {
	obj := c.info.Defs[id]
	if obj == nil || !isUnderlyingArray(obj.Type()) {
		return false
	}
	return c.markVar(id)
}

// markAssign marks the LHS locals that an assignment turns into frame values.
func (c *escapeChecker) markAssign(s *ast.AssignStmt) bool {
	switch s.Tok {
	case token.ADD_ASSIGN:
		// s += x on strings emits so_string_add, a fresh frame value.
		if len(s.Lhs) == 1 && isStringExpr(c.info, s.Lhs[0]) {
			return c.markVar(s.Lhs[0])
		}
	case token.ASSIGN, token.DEFINE:
		return c.markPairs(s.Lhs, s.Rhs)
	}
	return false
}

// markDecl marks the locals that a var declaration
// with initializers turns into frame values.
func (c *escapeChecker) markDecl(s *ast.DeclStmt) bool {
	gd, ok := s.Decl.(*ast.GenDecl)
	if !ok || gd.Tok != token.VAR {
		return false
	}
	changed := false
	for _, spec := range gd.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok || len(vs.Values) != len(vs.Names) {
			continue
		}
		lhs := make([]ast.Expr, len(vs.Names))
		for i, name := range vs.Names {
			lhs[i] = name
		}
		changed = c.markPairs(lhs, vs.Values) || changed
	}
	return changed
}

// markPairs marks each LHS local whose paired RHS is a frame value. A single
// multi-value RHS (x, y := f()) is an opaque call, so it marks nothing.
func (c *escapeChecker) markPairs(lhs, rhs []ast.Expr) bool {
	if len(lhs) != len(rhs) {
		return false
	}
	changed := false
	for i := range lhs {
		if c.isFrameValue(rhs[i]) {
			changed = c.markVar(lhs[i]) || changed
		}
	}
	return changed
}

// markVar records the local named by expr as holding a frame value
// and reports whether that changed the set.
func (c *escapeChecker) markVar(expr ast.Expr) bool {
	id, ok := expr.(*ast.Ident)
	if !ok || id.Name == "_" {
		return false
	}
	obj := c.info.ObjectOf(id)
	if obj == nil || !c.locals[obj] || c.fvars[obj] {
		return false
	}
	c.fvars[obj] = true
	return true
}

// isFrameValue reports whether the value of expr lives in the current frame.
// It follows only the positions that carry the value onward (a subexpression,
// a slice of it), never the opaque result of a call.
func (c *escapeChecker) isFrameValue(e ast.Expr) bool {
	switch x := e.(type) {
	case *ast.ParenExpr:
		return c.isFrameValue(x.X)
	case *ast.Ident:
		obj := c.info.ObjectOf(x)
		return obj != nil && c.fvars[obj]
	case *ast.BinaryExpr:
		return c.isFrameConcat(x)
	case *ast.UnaryExpr:
		return x.Op == token.AND && c.isFrameAddress(x.X)
	case *ast.CompositeLit:
		return c.isFrameComposite(x)
	case *ast.SliceExpr:
		return c.isFrameValue(x.X)
	case *ast.CallExpr:
		return c.isFrameCall(x)
	}
	return false
}

// isFrameConcat reports whether a binary expression is a string + that emits
// so_string_add, whose result is a fresh frame value. A + of two string literals
// folds to a static constant and is safe.
func (c *escapeChecker) isFrameConcat(x *ast.BinaryExpr) bool {
	if x.Op != token.ADD || !isStringExpr(c.info, x.X) {
		return false
	}
	return !(isStringLit(x.X) && isStringLit(x.Y))
}

// isFrameAddress reports whether &operand points into the current frame: the
// address of a local or parameter, or of a composite-literal temporary.
func (c *escapeChecker) isFrameAddress(operand ast.Expr) bool {
	switch o := operand.(type) {
	case *ast.CompositeLit:
		return true
	case *ast.Ident:
		obj := c.info.ObjectOf(o)
		return obj != nil && c.locals[obj]
	}
	return false
}

// isFrameComposite reports whether a composite literal is a frame value.
func (c *escapeChecker) isFrameComposite(x *ast.CompositeLit) bool {
	t := c.info.TypeOf(x)
	// An array literal lowers to a bare C array in the frame.
	if isUnderlyingArray(t) {
		return true
	}
	// A non-empty slice literal lowers to a so_Slice over a compound-literal
	// backing array in the frame (see emitSliceLit); an empty one lowers to
	// a null slice and is safe.
	if _, ok := t.Underlying().(*types.Slice); ok && len(x.Elts) > 0 {
		return true
	}
	// A struct literal is a frame value when one of its elements is (e.g. Pair{s: a + b}).
	for _, el := range x.Elts {
		if kv, ok := el.(*ast.KeyValueExpr); ok {
			el = kv.Value
		}
		// An array element is copied by value into the struct or slice, so its
		// own frame storage travels with the copy and does not escape. (Frame
		// memory it references through a nested string or slice is not tracked.)
		if isUnderlyingArray(c.info.TypeOf(el)) {
			continue
		}
		if c.isFrameValue(el) {
			return true
		}
	}
	return false
}

// isFrameCall reports whether a call produces a frame value. Conversions and the
// make and append builtins can; every other call is opaque and assumed not to.
func (c *escapeChecker) isFrameCall(call *ast.CallExpr) bool {
	if tv, ok := c.info.Types[call.Fun]; ok && tv.IsType() {
		return c.isFrameConversion(tv.Type, call)
	}
	id, ok := call.Fun.(*ast.Ident)
	if !ok {
		return false
	}
	b, ok := c.info.Uses[id].(*types.Builtin)
	if !ok {
		return false
	}
	switch b.Name() {
	case "make":
		// make of a slice or map emits so_make_slice / so_make_map.
		switch c.info.TypeOf(call).Underlying().(type) {
		case *types.Slice, *types.Map:
			return true
		}
	case "append":
		// append writes into the destination slice's existing backing (it never
		// reallocates, only asserts capacity), so its result aliases that backing.
		// It is a frame value only when the destination is.
		return len(call.Args) > 0 && c.isFrameValue(call.Args[0])
	}
	return false
}

// isFrameConversion reports whether a type conversion produces a frame value.
// []rune(s) allocates a fresh buffer in the frame; []byte(s) is a zero-copy view,
// so it is a frame value only when s is. Conversions to string are handled by
// isFrameStringConv. Any other conversion aliases its argument.
func (c *escapeChecker) isFrameConversion(target types.Type, call *ast.CallExpr) bool {
	if len(call.Args) != 1 {
		return false
	}
	arg := call.Args[0]
	argT := c.info.TypeOf(arg)
	switch t := target.Underlying().(type) {
	case *types.Slice:
		if b, ok := t.Elem().Underlying().(*types.Basic); ok && isStringType(argT) {
			switch b.Kind() {
			case types.Int32:
				return true // []rune(s): so_string_runes
			case types.Byte:
				return c.isFrameValue(arg) // []byte(s): view
			}
		}
	case *types.Basic:
		if t.Kind() == types.String {
			return c.isFrameStringConv(argT, arg)
		}
	}
	return c.isFrameValue(arg)
}

// isFrameStringConv reports whether a conversion to string produces a frame
// value. string([]rune), string(byte) and string(rune) allocate a fresh buffer
// in the frame. string([]byte) is a zero-copy view, so it is a frame value only
// when its argument is.
func (c *escapeChecker) isFrameStringConv(argT types.Type, arg ast.Expr) bool {
	switch a := argT.Underlying().(type) {
	case *types.Slice:
		if b, ok := a.Elem().Underlying().(*types.Basic); ok {
			switch b.Kind() {
			case types.Int32:
				return true // string([]rune): so_runes_string
			case types.Byte:
				return c.isFrameValue(arg) // string([]byte): view
			}
		}
	case *types.Basic:
		switch a.Kind() {
		case types.Byte, types.Int32:
			return true // string(byte) / string(rune)
		}
	}
	return c.isFrameValue(arg)
}

// escapes returns the return-value expressions that are frame values. It skips
// nested closures. A bare return has no results, so it never escapes (So has no
// named results, so bare returns occur only in functions that return nothing).
func (c *escapeChecker) escapes(decl *ast.FuncDecl) []ast.Node {
	var found []ast.Node
	c.walk(decl.Body, func(n ast.Node) {
		ret, ok := n.(*ast.ReturnStmt)
		if !ok {
			return
		}
		for _, r := range ret.Results {
			if c.isFrameValue(r) {
				found = append(found, r)
			}
		}
	})
	return found
}

// walk visits every node under root but does not enter nested closures, which
// have their own frame.
func (c *escapeChecker) walk(root ast.Node, visit func(ast.Node)) {
	ast.Inspect(root, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		if _, ok := n.(*ast.FuncLit); ok {
			return false
		}
		visit(n)
		return true
	})
}

// isStringExpr reports whether expr has string type.
func isStringExpr(info *types.Info, expr ast.Expr) bool {
	return isStringType(info.TypeOf(expr))
}

// isStringType reports whether t is a string type.
func isStringType(t types.Type) bool {
	if t == nil {
		return false
	}
	b, ok := t.Underlying().(*types.Basic)
	return ok && (b.Kind() == types.String || b.Kind() == types.UntypedString)
}

// isUnderlyingArray reports whether t is an array, named or not. So returns an
// array by value as a pointer into the frame (see returnType), so a local array
// is a frame value. A struct that wraps an array is copied by value and stays
// safe.
func isUnderlyingArray(t types.Type) bool {
	if t == nil {
		return false
	}
	_, ok := t.Underlying().(*types.Array)
	return ok
}
