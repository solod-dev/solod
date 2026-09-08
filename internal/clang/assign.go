package clang

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io"
	"maps"
	"slices"
)

// emitAssignStmt emits an assignment statement.
func (g *Generator) emitAssignStmt(w io.Writer, stmt *ast.AssignStmt) {
	stmt = unparenAssign(stmt)
	switch stmt.Tok {
	case token.DEFINE:
		g.emitDefine(w, stmt)

	case token.ASSIGN:
		g.emitAssign(w, stmt)

	case token.ADD_ASSIGN, token.SUB_ASSIGN, token.MUL_ASSIGN, token.QUO_ASSIGN,
		token.REM_ASSIGN, token.OR_ASSIGN, token.AND_ASSIGN, token.XOR_ASSIGN,
		token.SHL_ASSIGN, token.SHR_ASSIGN:
		g.checkMapIndex(stmt.Lhs[0])
		g.checkAssignCall(stmt)
		// String += uses so_string_add.
		if stmt.Tok == token.ADD_ASSIGN && g.hasStringType(stmt.Lhs[0]) {
			fmt.Fprint(w, g.indent())
			g.emitExpr(w, stmt.Lhs[0])
			fmt.Fprint(w, " = so_string_add(")
			g.emitExpr(w, stmt.Lhs[0])
			fmt.Fprint(w, ", ")
			g.emitExpr(w, stmt.Rhs[0])
			fmt.Fprint(w, ");\n")
			return
		}
		// Integer /= and %= guard against a zero divisor.
		if (stmt.Tok == token.QUO_ASSIGN || stmt.Tok == token.REM_ASSIGN) &&
			g.needsIntDivGuard(stmt.Lhs[0], stmt.Rhs[0]) {
			fmt.Fprint(w, g.indent())
			g.emitExpr(w, stmt.Lhs[0])
			if stmt.Tok == token.QUO_ASSIGN {
				fmt.Fprint(w, " = so_div(")
			} else {
				fmt.Fprint(w, " = so_mod(")
			}
			g.emitExpr(w, stmt.Lhs[0])
			fmt.Fprint(w, ", ")
			g.emitExpr(w, stmt.Rhs[0])
			fmt.Fprint(w, ");\n")
			return
		}
		fmt.Fprint(w, g.indent())
		g.emitExpr(w, stmt.Lhs[0])
		fmt.Fprintf(w, " %s ", stmt.Tok)
		g.emitExpr(w, stmt.Rhs[0])
		fmt.Fprint(w, ";\n")

	default:
		g.fail(stmt, "unsupported AssignStmt token: %s", stmt.Tok)
	}
}

// emitDefine emits a short variable declaration (:=).
func (g *Generator) emitDefine(w io.Writer, stmt *ast.AssignStmt) {
	if g.emitAssignSpecial(w, stmt, true) {
		return
	}
	a := g.assignment(w, stmt)
	a.check()
	a.hoist()
	for i, lhs := range stmt.Lhs {
		ident := lhs.(*ast.Ident)
		if ident.Name == "_" {
			// Blank identifier - the value is still evaluated.
			a.discardValue(i)
			continue
		}

		def := g.types.Defs[ident]
		if def == nil {
			// Redeclared variable - emit plain assignment.
			typ := g.types.Uses[ident].Type()
			fmt.Fprintf(w, "%s%s = ", g.indent(), ident.Name)
			g.emitExprAsType(w, stmt, a.value(i), typ)
			fmt.Fprint(w, ";\n")
			continue
		}

		typ := def.Type()
		ct := g.mapVarType(stmt, typ, true)
		if _, isArr := arrayType(typ); isArr {
			// C cannot assign an array, so it needs a declaration of its own.
			g.emitArrayVarDecl(w, ct, ident.Name, a.value(i))
			continue
		}

		fmt.Fprintf(w, "%s%s = ", g.indent(), ct.Decl(ident.Name))
		g.emitExpr(w, a.value(i))
		fmt.Fprint(w, ";\n")
	}
}

// emitAssign emits a regular assignment (=).
func (g *Generator) emitAssign(w io.Writer, stmt *ast.AssignStmt) {
	a := g.assignment(w, stmt)
	a.check()
	if g.emitAssignSpecial(w, stmt, false) {
		return
	}
	// Regular assignment.
	a.hoist()
	for i, lhs := range stmt.Lhs {
		// Blank identifier - emit a void expression.
		if ident, ok := lhs.(*ast.Ident); ok && ident.Name == "_" {
			a.discardValue(i)
			continue
		}

		// Map index assignment uses so_map_set.
		if idx, ok := lhs.(*ast.IndexExpr); ok {
			if isMapType(g.types.TypeOf(idx.X)) {
				g.emitMapIndexAssign(w, stmt, idx, a.value(i))
				continue
			}
		}

		// Array assignment uses memcpy.
		lhsType := g.types.TypeOf(lhs)
		if arr, ok := arrayType(lhsType); ok {
			// The array type is known at compile time, so sizeof takes the
			// type instead of the left side to avoid evaluating it twice.
			arrCType := g.mapTypeName(stmt, arr.Elem()) + arrayDims(arr)
			fmt.Fprintf(w, "%smemcpy(", g.indent())
			g.emitExpr(w, lhs)
			fmt.Fprint(w, ", ")
			g.emitArrayValue(w, stmt, a.value(i), arr)
			fmt.Fprintf(w, ", sizeof(%s));\n", arrCType)
			continue
		}

		// Non-array assignment.
		fmt.Fprint(w, g.indent())
		g.emitExpr(w, lhs)
		fmt.Fprint(w, " = ")
		g.emitExprAsType(w, stmt, a.value(i), lhsType)
		fmt.Fprint(w, ";\n")
	}
}

// emitAssignSpecial emits an assignment of a special form: a comma-ok type
// assertion, a comma-ok map read, or a multi-return destructuring. It reports
// whether the statement has one of these forms. define reports whether the
// statement declares the variables of its left side.
func (g *Generator) emitAssignSpecial(w io.Writer, stmt *ast.AssignStmt, define bool) bool {
	if len(stmt.Lhs) == 2 && len(stmt.Rhs) == 1 {
		// Comma-ok type assertion: v, ok := s.(Rect)
		if ta, ok := stmt.Rhs[0].(*ast.TypeAssertExpr); ok {
			g.emitTypeAssertion(w, stmt, ta)
			return true
		}
		// Comma-ok map read: v, ok := m[key]
		if idx, ok := stmt.Rhs[0].(*ast.IndexExpr); ok && isMapType(g.types.TypeOf(idx.X)) {
			if define {
				// The map read is inlined into the declaration of v, so a
				// self-shadowing key reads v instead of the outer variable.
				g.assignment(w, stmt).checkSelfShadow()
			}
			g.emitMapCommaOk(w, stmt, idx, define)
			return true
		}
	}
	// Multi-return destructuring: x, y := f()
	if len(stmt.Lhs) > 1 && len(stmt.Rhs) == 1 {
		if call, ok := stmt.Rhs[0].(*ast.CallExpr); ok {
			if define {
				g.emitMultiReturnDefine(w, stmt, call)
			} else {
				g.emitMultiReturnAssign(w, stmt, call)
			}
			return true
		}
	}
	return false
}

// checkAssignCall rejects a call in the left side of a compound assignment.
// A string += and a guarded /= or %= write the left side two times, so the
// call would run two times and could return a different result if it has
// side effects.
func (g *Generator) checkAssignCall(stmt *ast.AssignStmt) {
	ast.Inspect(stmt.Lhs[0], func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		g.fail(call, "call in the left side of a compound assignment is not supported")
		return false
	})
}

// callsFunc reports whether an expression calls a function with side effects.
func (g *Generator) callsFunc(expr ast.Expr) bool {
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if found {
			return false
		}
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		// A pure call still has operands, so the walk goes on.
		if g.isPureCall(call) {
			return true
		}
		found = true
		return false
	})
	return found
}

// isPureCall reports whether a call has no side effects. A type conversion and
// the len and cap builtins are pure. Every other call is not.
func (g *Generator) isPureCall(call *ast.CallExpr) bool {
	if tv, ok := g.types.Types[call.Fun]; ok && tv.IsType() {
		return true
	}
	ident, ok := call.Fun.(*ast.Ident)
	if !ok {
		return false
	}
	b, ok := g.types.Uses[ident].(*types.Builtin)
	return ok && (b.Name() == "len" || b.Name() == "cap")
}

// assignment emits the values of an assignment statement. Go evaluates the
// whole right side before it assigns any target, so a statement that reads
// what it assigns needs a temporary for every value.
type assignment struct {
	g    *Generator
	w    io.Writer
	stmt *ast.AssignStmt
	// values holds the value for each target. hoist replaces a value with
	// the temporary that holds it, or with nil for a discarded value.
	values []ast.Expr
	// assigned holds the names of the variables the statement sets.
	assigned map[string]bool
}

// assignment returns an emitter for a define or a plain assignment.
func (g *Generator) assignment(w io.Writer, stmt *ast.AssignStmt) *assignment {
	return &assignment{g: g, w: w, stmt: stmt, values: stmt.Rhs}
}

// check rejects an unsupported assignment.
func (a *assignment) check() {
	if a.stmt.Tok == token.DEFINE {
		a.checkSelfShadow()
		return
	}
	a.checkTargets()
	a.checkAnonStruct()
}

// checkSelfShadow rejects a definition of a variable x that reads a variable
// with the same name from an outer scope, eg. `x := x + 1`.
func (a *assignment) checkSelfShadow() {
	rhsNames := collectIdents(a.stmt.Rhs...)
	for _, lhs := range a.stmt.Lhs {
		ident, ok := lhs.(*ast.Ident)
		if !ok || ident.Name == "_" {
			continue
		}
		if a.g.types.Defs[ident] == nil {
			continue
		}
		if rhsNames[ident.Name] {
			a.g.fail(a.stmt, "self-shadowing variable %q is not supported", ident.Name)
		}
	}
}

// checkTargets rejects a multiple assignment where a target on the left
// reads a variable the same statement assigns.
func (a *assignment) checkTargets() {
	if len(a.stmt.Lhs) < 2 {
		return
	}
	assigned := a.assignedVars()
	for _, lhs := range a.stmt.Lhs {
		// A plain identifier has no operand to evaluate.
		if _, ok := lhs.(*ast.Ident); ok {
			continue
		}
		for name := range collectIdents(lhs) {
			if assigned[name] {
				a.g.fail(a.stmt, "multiple assignment reads and assigns %q in the same statement", name)
			}
		}
	}
}

// checkAnonStruct rejects an assignment to an anonymous struct variable or field.
func (a *assignment) checkAnonStruct() {
	for _, lhs := range a.stmt.Lhs {
		// A blank target is not emitted, so its type is irrelevant.
		if ident, ok := lhs.(*ast.Ident); ok && ident.Name == "_" {
			continue
		}
		if isAnonStruct(a.g.types.TypeOf(lhs)) {
			a.g.fail(a.stmt, "cannot assign an anonymous struct; declare a named struct type instead")
		}
	}
}

// value returns the value for the i-th target.
func (a *assignment) value(i int) ast.Expr {
	return a.values[i]
}

// discardValue emits the value for a blank target. hoist discards the value
// in place, so this emits nothing after a hoist.
func (a *assignment) discardValue(i int) {
	if a.values[i] == nil {
		return
	}
	a.g.emitDiscard(a.w, a.values[i])
}

// hoist evaluates the right side into temporaries and replaces each value with
// a reference to its temporary. It does nothing when the statement keeps Go's
// evaluation order without temporaries.
//
// A blank target gets no temporary and no reference. Its value is discarded in
// place, which keeps the calls on the right in order.
func (a *assignment) hoist() {
	if !a.needsTemps() {
		return
	}
	g, stmt := a.g, a.stmt
	values := make([]ast.Expr, len(stmt.Rhs))
	for i, rhs := range stmt.Rhs {
		// A constant needs no temporary. No assignment can change it.
		if tv, ok := g.types.Types[rhs]; ok && tv.Value != nil {
			values[i] = rhs
			continue
		}
		typ := a.targetType(i)
		if typ == nil {
			g.emitDiscard(a.w, rhs)
			continue
		}
		ref := &ast.Ident{NamePos: rhs.Pos(), Name: g.newTemp(stmt, tempAssign)}
		// The reference has no declaration to look up, so record its type directly.
		g.types.Types[ref] = types.TypeAndValue{Type: typ}
		values[i] = ref

		ct := g.mapVarType(stmt, typ, true)
		if !ct.IsArray() {
			fmt.Fprintf(a.w, "%s%s = ", g.indent(), ct.Decl(ref.Name))
			g.emitExprAsType(a.w, stmt, rhs, typ)
			fmt.Fprint(a.w, ";\n")
			continue
		}
		// C cannot assign an array, so a composite literal initializes the
		// temporary and any other value is copied into it.
		if lit, isLit := ast.Unparen(rhs).(*ast.CompositeLit); isLit {
			fmt.Fprintf(a.w, "%s%s = ", g.indent(), ct.Decl(ref.Name))
			g.emitExpr(a.w, lit)
			fmt.Fprint(a.w, ";\n")
			continue
		}
		fmt.Fprintf(a.w, "%s%s;\n", g.indent(), ct.Decl(ref.Name))
		fmt.Fprintf(a.w, "%smemcpy(%s, ", g.indent(), ref.Name)
		g.emitExpr(a.w, rhs)
		fmt.Fprintf(a.w, ", sizeof(%s));\n", ref.Name)
	}
	a.values = values
}

// needsTemps reports whether a multiple assignment must evaluate its right
// side into temporaries to match Go's semantics (evaluate the whole right side
// before assigning any variable on the left).
func (a *assignment) needsTemps() bool {
	if len(a.stmt.Lhs) < 2 || len(a.stmt.Rhs) < 2 {
		return false
	}
	changed := a.changedNames()
	if len(changed) == 0 {
		return false
	}
	if slices.ContainsFunc(a.stmt.Rhs, a.g.callsFunc) {
		return true
	}
	for name := range collectIdents(a.stmt.Rhs...) {
		if changed[name] {
			return true
		}
	}
	return false
}

// changedNames returns the names of the variables the statement changes. It
// adds the root of every target that writes into memory: arr in arr[i], p in
// *p and p.f.
func (a *assignment) changedNames() map[string]bool {
	names := maps.Clone(a.assignedVars())
	for _, lhs := range a.stmt.Lhs {
		if _, ok := lhs.(*ast.Ident); ok {
			continue
		}
		if ident := rootIdent(lhs); ident != nil {
			names[ident.Name] = true
		}
	}
	return names
}

// assignedVars returns the names of the variables the statement sets.
func (a *assignment) assignedVars() map[string]bool {
	if a.assigned != nil {
		return a.assigned
	}
	a.assigned = map[string]bool{}
	for _, lhs := range a.stmt.Lhs {
		// Only a plain identifier target sets a variable. Any other target
		// writes into memory, and the variable itself keeps its value.
		ident, ok := lhs.(*ast.Ident)
		if !ok || ident.Name == "_" {
			continue
		}
		// Redeclaration does not introduce a new variable.
		if a.stmt.Tok == token.DEFINE && a.g.types.Defs[ident] != nil {
			continue
		}
		a.assigned[ident.Name] = true
	}
	return a.assigned
}

// targetType returns the type of the i-th target of an assignment.
// It returns nil for a blank identifier, which has no type.
func (a *assignment) targetType(i int) types.Type {
	lhs := a.stmt.Lhs[i]
	ident, isIdent := lhs.(*ast.Ident)
	if isIdent && ident.Name == "_" {
		return nil
	}
	if a.stmt.Tok != token.DEFINE {
		return a.g.types.TypeOf(lhs)
	}
	if def := a.g.types.Defs[ident]; def != nil {
		return def.Type()
	}
	return a.g.types.Uses[ident].Type()
}

// rootIdent returns the identifier at the root of an assignment target:
// arr in arr[i], p in *p and p.f. It returns nil for any other target.
func rootIdent(expr ast.Expr) *ast.Ident {
	for {
		switch e := expr.(type) {
		case *ast.Ident:
			return e
		case *ast.IndexExpr:
			expr = e.X
		case *ast.StarExpr:
			expr = e.X
		case *ast.SelectorExpr:
			expr = e.X
		case *ast.ParenExpr:
			expr = e.X
		default:
			return nil
		}
	}
}

// collectIdents returns the set of identifier names in the given expressions.
// The blank identifier is excluded.
func collectIdents(exprs ...ast.Expr) map[string]bool {
	names := map[string]bool{}
	var visit func(ast.Node) bool
	visit = func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.Ident:
			if n.Name != "_" {
				names[n.Name] = true
			}
		case *ast.KeyValueExpr:
			// Only recurse into Value, skip Key (struct field names
			// are not variable references; map key variables are
			// also skipped but self-shadowing there is unlikely).
			ast.Inspect(n.Value, visit)
			return false
		case *ast.SelectorExpr:
			// Only recurse into X, skip Sel (field/method names
			// are not variable references).
			ast.Inspect(n.X, visit)
			return false
		}
		return true
	}
	for _, expr := range exprs {
		ast.Inspect(expr, visit)
	}
	return names
}

// unparenAssign removes the parentheses around each assignment target.
func unparenAssign(stmt *ast.AssignStmt) *ast.AssignStmt {
	var lhs []ast.Expr
	for i, expr := range stmt.Lhs {
		bare := ast.Unparen(expr)
		if bare == expr {
			continue
		}
		if lhs == nil {
			lhs = slices.Clone(stmt.Lhs)
		}
		lhs[i] = bare
	}
	if lhs == nil {
		return stmt
	}
	clone := *stmt
	clone.Lhs = lhs
	return &clone
}

// isTypeParam reports whether typ is a type parameter of a generic function.
func isTypeParam(typ types.Type) bool {
	_, ok := types.Unalias(typ).(*types.TypeParam)
	return ok
}
