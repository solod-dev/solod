package clang

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io"
	"strings"
)

// emitFuncProto writes a full C function prototype (e.g. "static void main_foo(int x)")
// without a terminator. Returns the function's type signature for callers that need it.
func (g *Generator) emitFuncProto(w io.Writer, decl *ast.FuncDecl) *types.Signature {
	// Specifier: static inline for so:inline, static for unexported,
	// empty for exported and main.
	dirs := g.funcDirs[decl]
	spec := ""
	if dirs.inline {
		spec = "static inline "
	} else if decl.Name.Name != "main" {
		exported := ast.IsExported(decl.Name.Name)
		if exported && decl.Recv != nil {
			exported = ast.IsExported(recvTypeName(decl.Recv.List[0]))
		}
		if !exported {
			spec = "static "
		}
	}
	attr := dirs.attrString()
	if attr != "" {
		spec = spec + attr + " "
	}

	sig := g.funcSig(decl)

	// Return type.
	retType := "void"
	if isMainFunc(decl) {
		retType = "int"
	} else if decl.Type.Results != nil && len(decl.Type.Results.List) > 0 {
		retType = g.returnType(decl, sig)
	}

	// Name: methods use RecvType_Method, functions use symbolName.
	name := g.symbolName(g.types.Defs[decl.Name])
	if decl.Recv != nil {
		name = g.symbolName(g.recvTypeObj(decl.Recv.List[0])) + "_" + decl.Name.Name
	}

	// Parameters: methods prepend receiver
	// (void* self for pointer, T name for value).
	var parts []string
	if decl.Recv != nil {
		recv := decl.Recv.List[0]
		if _, ok := recv.Type.(*ast.Ident); ok {
			// Value receiver: pass struct by value.
			cStructType := g.symbolName(g.recvTypeObj(recv))
			recvName := "self"
			if len(recv.Names) > 0 {
				recvName = recv.Names[0].Name
			}
			parts = append(parts, cStructType+" "+recvName)
		} else {
			parts = append(parts, "void* self")
		}
	}
	if decl.Type.Params != nil {
		for _, field := range decl.Type.Params.List {
			typ := g.types.TypeOf(field.Type)
			ct := g.mapCType(decl, typ)
			for _, n := range field.Names {
				parts = append(parts, ct.Decl(n.Name))
			}
		}
	}
	params := "void"
	if isMainFunc(decl) && g.importsOS() {
		params = "int argc, char* argv[]"
	} else if len(parts) > 0 {
		params = strings.Join(parts, ", ")
	}

	fmt.Fprintf(w, "%s%s %s(%s)", spec, retType, name, params)
	return sig
}

// emitFuncTypeSpec emits a C function pointer typedef.
func (g *Generator) emitFuncTypeSpec(w io.Writer, spec *ast.TypeSpec) {
	named := g.types.Defs[spec.Name].Type().(*types.Named)
	sig := named.Underlying().(*types.Signature)

	retType := g.returnType(spec, sig)

	var params []string
	for parVar := range sig.Params().Variables() {
		params = append(params, g.mapType(spec, parVar.Type()))
	}

	name := g.declSymbolName(g.types.Defs[spec.Name])
	fmt.Fprintf(w, "%stypedef %s (*%s)(%s);\n", g.indent(), retType, name, strings.Join(params, ", "))
}

// emitFuncDecl emits a function declaration into the .c file.
// Inline functions are skipped here - they are emitted into the header
// by [Generator.emitInlineFuncDecl].
func (g *Generator) emitFuncDecl(decl *ast.FuncDecl) {
	if decl.Body == nil || g.hasExtern(g.types.Defs[decl.Name]) {
		return
	}
	if isInitFunc(decl) {
		return
	}
	if g.funcDirs[decl].inline {
		return
	}
	g.emitFuncBody(decl)
}

// emitInlineFuncDecl emits a so:inline function declaration into the header.
// Generic functions are emitted as #define macros; non-generic as static inline.
func (g *Generator) emitInlineFuncDecl(w io.Writer, decl *ast.FuncDecl) {
	if isGenericFunc(decl) {
		g.emitMacroFuncDecl(w, decl)
		return
	}
	saved := g.state.writer
	g.state.writer = w
	g.emitFuncBody(decl)
	g.state.writer = saved
}

// emitMacroFuncDecl emits a generic so:inline function as a #define macro.
func (g *Generator) emitMacroFuncDecl(w io.Writer, decl *ast.FuncDecl) {
	sig := g.funcSig(decl)
	g.rejectNamedReturns(decl, sig)

	// Build macro name.
	name := g.symbolName(g.types.Defs[decl.Name])
	if decl.Recv != nil {
		name = g.symbolName(g.recvTypeObj(decl.Recv.List[0])) + "_" + decl.Name.Name
	}

	// Build param list: type params, then receiver (for methods), then regular params.
	// Non-type params are suffixed with _ to avoid name collisions (b->val = val).
	// References are wrapped in parens to avoid syntax errors (&b->val).
	var params []string
	macroParams := make(map[string]bool)
	if decl.Type.TypeParams != nil {
		for _, field := range decl.Type.TypeParams.List {
			for _, n := range field.Names {
				params = append(params, n.Name)
			}
		}
	}
	if decl.Recv != nil {
		recv := decl.Recv.List[0]
		// Add receiver type params (no suffix - these are type names).
		params = append(params, recvTypeParams(recv)...)
		// Add receiver as parameter (suffixed).
		recvName := "self"
		if len(recv.Names) > 0 {
			recvName = recv.Names[0].Name
		}
		macroParams[recvName] = true
		params = append(params, recvName+"_")
	}
	if decl.Type.Params != nil {
		for _, field := range decl.Type.Params.List {
			for _, n := range field.Names {
				macroParams[n.Name] = true
				params = append(params, n.Name+"_")
			}
		}
	}

	// Capture body output.
	var buf strings.Builder
	savedState := g.state
	g.state.writer = &buf
	g.state.funcSig = sig
	g.state.tempCount = 0
	g.state.indent = 1
	g.state.inMacro = true
	g.state.macroParams = macroParams
	g.state.defers = nil
	g.walkStmts(decl.Body.List)
	g.state = savedState

	// Determine if returning or void.
	hasReturn := sig.Results() != nil && sig.Results().Len() > 0

	// Emit #define with line continuations.
	body := buf.String()
	// Trim trailing newline.
	body = strings.TrimRight(body, "\n")
	lines := strings.Split(body, "\n")

	if !g.emitComments(w, decl) {
		fmt.Fprintln(w)
	}
	if hasReturn {
		fmt.Fprintf(w, "#define %s(%s) ({", name, strings.Join(params, ", "))
	} else {
		fmt.Fprintf(w, "#define %s(%s) do {", name, strings.Join(params, ", "))
	}
	for _, line := range lines {
		fmt.Fprintf(w, " \\\n%s", line)
	}
	if hasReturn {
		fmt.Fprintln(w, " \\")
		fmt.Fprintln(w, "})")
	} else {
		fmt.Fprintln(w, " \\")
		fmt.Fprintln(w, "} while (0)")
	}
}

// emitFuncBody emits a function or method body. Shared by [Generator.emitFuncDecl]
// and [Generator.emitInlineFuncDecl].
func (g *Generator) emitFuncBody(decl *ast.FuncDecl) {
	if decl.Recv != nil {
		g.emitMethodDecl(decl)
		return
	}

	// Init emission state.
	w := g.state.writer
	sig := g.funcSig(decl)
	g.rejectNamedReturns(decl, sig)
	g.state.funcSig = sig
	g.state.tempCount = 0

	// Emit comments and function prototype.
	if !g.emitComments(w, decl) {
		fmt.Fprintln(w)
	}
	g.emitFuncProto(w, decl)
	fmt.Fprintln(w, " {")

	// Emit function body, handling deferred calls if needed.
	g.state.indent++
	if isMainFunc(decl) && g.importsOS() {
		fmt.Fprintf(w, "%sso_String _so_argv[argc];\n", g.indent())
		fmt.Fprintf(w, "%sso_args_init(argc, argv, _so_argv);\n", g.indent())
	}
	g.walkStmts(decl.Body.List)
	if !endsWithReturn(decl.Body.List) {
		g.emitDeferredCalls()
		if isMainFunc(decl) {
			fmt.Fprintf(w, "%sreturn 0;\n", g.indent())
		}
	}
	g.state.indent--
	fmt.Fprintf(w, "}\n")

	// Reset state.
	g.state.defers = nil
	g.state.funcSig = nil
}

// emitFuncCall emits a regular function call.
func (g *Generator) emitFuncCall(call *ast.CallExpr) {
	w := g.state.writer
	if ident, ok := call.Fun.(*ast.Ident); ok {
		if bi, ok := g.types.Uses[ident].(*types.Builtin); ok {
			if g.emitBuiltin(call, ident, bi) {
				return
			}
		} else {
			g.emitExpr(call.Fun)
		}
	} else {
		g.emitExpr(call.Fun)
	}

	// Emit arguments, wrapping as interfaces if needed.
	var sig *types.Signature
	if funType := g.types.TypeOf(call.Fun); funType != nil {
		// Get the function signature to wrap value arguments as interfaces if needed.
		sig, _ = funType.Underlying().(*types.Signature)
	}
	fmt.Fprintf(w, "(")

	if ext, ok := g.callExtern(call); ok && !ext.nodecay {
		// Extern C call: decay all args to C-compatible types.
		// So wrapper types (so_String, so_Slice) must be unwrapped to their
		// underlying C representations for C function macros.
		if call.Ellipsis.IsValid() {
			g.fail(call, "spreading variadic arguments to an extern function is not supported")
		}
		g.emitCArgs(call)
	} else if sig != nil && sig.Variadic() && !call.Ellipsis.IsValid() {
		// Variadic call with individual args: pack trailing args into a slice literal.
		g.emitFixedArgs(call, sig)
		g.emitVariadicArgs(call, sig)
	} else {
		// Regular call: emit all args as-is.
		for i, arg := range call.Args {
			if i > 0 {
				fmt.Fprintf(w, ", ")
			}
			if sig != nil && i < sig.Params().Len() {
				// Emit arg, wrapping as interface if needed based on parameter type.
				g.emitExprAsType(call, arg, sig.Params().At(i).Type())
			} else {
				// No signature available (e.g. func literal), emit arg as-is.
				g.emitExpr(arg)
			}
		}
	}

	fmt.Fprintf(w, ")")
}

// emitFixedArgs emits the non-variadic arguments for a variadic call.
func (g *Generator) emitFixedArgs(call *ast.CallExpr, sig *types.Signature) {
	w := g.state.writer
	fixedCount := sig.Params().Len() - 1
	for i := 0; i < fixedCount && i < len(call.Args); i++ {
		if i > 0 {
			fmt.Fprintf(w, ", ")
		}
		g.emitExprAsType(call, call.Args[i], sig.Params().At(i).Type())
	}
}

// emitVariadicArgs packs trailing arguments into an inline so_Slice literal.
func (g *Generator) emitVariadicArgs(call *ast.CallExpr, sig *types.Signature) {
	w := g.state.writer
	fixedCount := sig.Params().Len() - 1
	variadicArgs := call.Args[fixedCount:]

	if fixedCount > 0 {
		fmt.Fprintf(w, ", ")
	}

	variadicParam := sig.Params().At(sig.Params().Len() - 1)
	elemType := g.mapType(call, variadicParam.Type().(*types.Slice).Elem())
	count := len(variadicArgs)

	fmt.Fprintf(w, "(so_Slice){(%s[%d]){", elemType, count)
	targetType := variadicParam.Type().(*types.Slice).Elem()
	for i, arg := range variadicArgs {
		if i > 0 {
			fmt.Fprintf(w, ", ")
		}
		g.emitExprAsType(call, arg, targetType)
	}
	fmt.Fprintf(w, "}, %d, %d}", count, count)
}

// emitCArgs emits arguments for an extern C function call.
func (g *Generator) emitCArgs(call *ast.CallExpr) {
	w := g.state.writer
	var sig *types.Signature
	if funType := g.types.TypeOf(call.Fun); funType != nil {
		sig, _ = funType.Underlying().(*types.Signature)
	}
	for i, arg := range call.Args {
		if i > 0 {
			fmt.Fprintf(w, ", ")
		}
		// Interface-typed parameters (e.g. Allocator) need emitExprAsType
		// to convert nil to a zero-initialized struct instead of NULL.
		if sig != nil && i < sig.Params().Len() && isNamedNonEmptyInterface(sig.Params().At(i).Type()) {
			g.emitExprAsType(call, arg, sig.Params().At(i).Type())
		} else {
			g.emitCArg(arg)
		}
	}
}

// emitCArg emits an expression decayed to its C-compatible type:
// string literals to raw C strings, strings to char*, slices to void*.
func (g *Generator) emitCArg(arg ast.Expr) {
	w := g.state.writer
	if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		fmt.Fprintf(w, "%s", rawStringValue(lit))
	} else if g.hasStringType(arg) {
		fmt.Fprintf(w, "so_cstr(")
		g.emitExpr(arg)
		fmt.Fprintf(w, ")")
	} else if _, ok := g.types.TypeOf(arg).Underlying().(*types.Slice); ok {
		fmt.Fprintf(w, "so_decay(")
		g.emitExpr(arg)
		fmt.Fprintf(w, ")")
	} else if isErrorType(g.types.TypeOf(arg)) {
		fmt.Fprintf(w, "errors_cstr(")
		g.emitExpr(arg)
		fmt.Fprintf(w, ")")
	} else {
		g.emitExpr(arg)
	}
}

// isGenericFunc reports whether a function declaration is generic
// (has type params on the function itself or on its receiver type).
func isGenericFunc(decl *ast.FuncDecl) bool {
	if decl.Type.TypeParams != nil && len(decl.Type.TypeParams.List) > 0 {
		return true
	}
	if decl.Recv != nil {
		recv := decl.Recv.List[0]
		typ := recv.Type
		if star, ok := typ.(*ast.StarExpr); ok {
			typ = star.X
		}
		switch typ.(type) {
		case *ast.IndexExpr, *ast.IndexListExpr:
			return true
		}
	}
	return false
}

// isMainFunc reports whether a function declaration is the main function.
func isMainFunc(decl *ast.FuncDecl) bool {
	return decl.Name.Name == "main" && decl.Recv == nil
}

// isInitFunc reports whether a function declaration is the init function.
func isInitFunc(decl *ast.FuncDecl) bool {
	return decl.Name.Name == "init" && decl.Recv == nil
}

// hasUnexportedTypes reports whether a function declaration
// references any unexported types from the current package.
func (g *Generator) hasUnexportedTypes(decl *ast.FuncDecl) bool {
	sig := g.funcSig(decl)
	for p := range sig.Params().Variables() {
		if g.isUnexportedType(p.Type()) {
			return true
		}
	}
	for r := range sig.Results().Variables() {
		if g.isUnexportedType(r.Type()) {
			return true
		}
	}
	return false
}

// importsOS reports whether the current package imports "os",
// which determines whether we need to initialize argc/argv in main().
func (g *Generator) importsOS() bool {
	// Only check the main package for simplicity. If "os" is imported in
	// a non-main package, the user will have to import "os" in main too
	// to signal that they want argc/argv support.
	_, ok := g.pkg.Imports["solod.dev/so/os"]
	return ok
}

// funcSig returns the types.Signature for a function or method declaration.
func (g *Generator) funcSig(decl *ast.FuncDecl) *types.Signature {
	if decl.Recv != nil {
		return g.types.ObjectOf(decl.Name).Type().(*types.Signature)
	}
	return g.types.Defs[decl.Name].Type().(*types.Signature)
}

// endsWithReturn reports whether a statement list ends with a return statement.
func endsWithReturn(stmts []ast.Stmt) bool {
	if len(stmts) == 0 {
		return false
	}
	_, ok := stmts[len(stmts)-1].(*ast.ReturnStmt)
	return ok
}

// recvTypeName returns the Go type name from a method receiver field.
// Handles both pointer receivers (*Rect) and value receivers (Rect).
func recvTypeName(recv *ast.Field) string {
	typ := recv.Type
	// Unwrap pointer receiver.
	if star, ok := typ.(*ast.StarExpr); ok {
		typ = star.X
	}
	// Unwrap generic type parameters.
	switch t := typ.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.IndexExpr:
		return t.X.(*ast.Ident).Name
	case *ast.IndexListExpr:
		return t.X.(*ast.Ident).Name
	}
	panic(fmt.Sprintf("unsupported receiver type: %T", recv.Type))
}

// recvTypeObj returns the types.Object for the receiver type of a method.
func (g *Generator) recvTypeObj(recv *ast.Field) types.Object {
	typ := recv.Type
	if star, ok := typ.(*ast.StarExpr); ok {
		typ = star.X
	}
	switch t := typ.(type) {
	case *ast.Ident:
		return g.types.Uses[t]
	case *ast.IndexExpr:
		return g.types.Uses[t.X.(*ast.Ident)]
	case *ast.IndexListExpr:
		return g.types.Uses[t.X.(*ast.Ident)]
	}
	g.fail(recv, "unsupported receiver type: %T", recv.Type)
	return nil // unreachable
}

// recvTypeParams extracts type parameter names from a generic receiver field.
func recvTypeParams(recv *ast.Field) []string {
	typ := recv.Type
	if star, ok := typ.(*ast.StarExpr); ok {
		typ = star.X
	}
	switch t := typ.(type) {
	case *ast.IndexExpr:
		if ident, ok := t.Index.(*ast.Ident); ok {
			return []string{ident.Name}
		}
	case *ast.IndexListExpr:
		var names []string
		for _, idx := range t.Indices {
			if ident, ok := idx.(*ast.Ident); ok {
				names = append(names, ident.Name)
			}
		}
		return names
	}
	return nil
}
