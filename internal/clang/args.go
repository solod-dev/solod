package clang

import (
	"fmt"
	"go/ast"
	"go/types"
	"io"
)

// callArgs writes the arguments of a call. It holds the separator between
// the arguments, so a caller emits an argument without knowing its position.
type callArgs struct {
	g     *Generator
	w     io.Writer
	node  ast.Node         // the position for a diagnostic
	macro bool             // the call emits as a macro
	count int              // the arguments written so far
	call  *ast.CallExpr    // the call, set by emit
	sig   *types.Signature // the callee signature, nil for a call through a value
}

// callArgs returns a writer for the arguments of a call.
// node gives the position of a diagnostic about an argument.
func (g *Generator) callArgs(w io.Writer, node ast.Node, macro bool) *callArgs {
	return &callArgs{g: g, w: w, node: node, macro: macro}
}

// emit writes one value argument. The write function emits the value itself.
func (a *callArgs) emitArg(write func()) {
	a.separate()
	if a.macro {
		fmt.Fprint(a.w, "(")
	}
	write()
	if a.macro {
		fmt.Fprint(a.w, ")")
	}
}

// emitType writes one type argument of a macro call.
func (a *callArgs) emitType(name string) {
	a.separate()
	fmt.Fprint(a.w, name)
}

// emitExpr writes one value argument coerced to the parameter type.
func (a *callArgs) emitExpr(arg ast.Expr, paramType types.Type) {
	a.emitArg(func() { a.g.emitCallArg(a.w, a.node, arg, paramType) })
}

// separate writes the comma before an argument that follows another one.
func (a *callArgs) separate() {
	if a.count > 0 {
		fmt.Fprint(a.w, ", ")
	}
	a.count++
}

// emit writes the arguments of a function or a method call.
// sig is nil for a call through a value with no signature.
func (a *callArgs) emit(call *ast.CallExpr, sig *types.Signature, extern externDecl, isExtern bool) {
	a.call, a.sig = call, sig
	if isExtern {
		if !extern.nodecay {
			// Extern C function: decay args to C-compatible types.
			a.emitDecayArgs()
			return
		}
		if a.sig != nil && a.sig.Variadic() {
			// Extern nodecay function: emit the variadic args flat.
			a.emitExternVarArgs()
			return
		}
	}
	if a.sig != nil && a.sig.Variadic() && !a.call.Ellipsis.IsValid() {
		// Variadic call with individual args: pack trailing args into a slice literal.
		a.emitVarArgs()
		return
	}
	// Non-variadic call or variadic call with ellipsis: emit all args directly.
	for i, arg := range a.call.Args {
		if a.sig != nil && i < a.sig.Params().Len() {
			a.emitExpr(arg, a.sig.Params().At(i).Type())
		} else {
			// No signature available (e.g. func literal), emit arg as-is.
			a.emitArg(func() { a.g.emitExpr(a.w, arg) })
		}
	}
}

// emitVarArgs packs the trailing arguments into an inline so_Slice literal.
func (a *callArgs) emitVarArgs() {
	// Emit fixed args first.
	fixedCount := a.sig.Params().Len() - 1
	for i := 0; i < fixedCount && i < len(a.call.Args); i++ {
		a.emitExpr(a.call.Args[i], a.sig.Params().At(i).Type())
	}

	// Emit variadic args as a so_Slice literal.
	variadicArgs := a.call.Args[fixedCount:]
	elemType := a.sig.Params().At(fixedCount).Type().(*types.Slice).Elem()
	cElemType := a.g.mapTypeName(a.node, elemType)
	count := len(variadicArgs)

	if count == 0 {
		// No variadic args: emit a nil slice.
		a.emitArg(func() { fmt.Fprint(a.w, "(so_Slice){}") })
		return
	}

	a.emitArg(func() {
		fmt.Fprintf(a.w, "(so_Slice){(%s[%d]){", cElemType, count)
		// The slice literal already protects its elements from the
		// preprocessor, so an element needs no parentheses of its own.
		elems := a.g.callArgs(a.w, a.node, false)
		for _, arg := range variadicArgs {
			elems.emitArg(func() { a.g.emitExprAsType(a.w, a.node, arg, elemType) })
		}
		fmt.Fprintf(a.w, "}, %d, %d}", count, count)
	})
}

// emitDecayArgs writes the arguments of a call to an extern C function,
// decayed to their C-compatible types.
func (a *callArgs) emitDecayArgs() {
	if a.call.Ellipsis.IsValid() {
		a.g.fail(a.call, "spreading variadic arguments to an extern function is not supported")
	}
	for i, arg := range a.call.Args {
		// Interface-typed parameters (e.g. Allocator) need emitExprAsType
		// to convert nil to a zero-initialized struct instead of NULL.
		if a.sig != nil && i < a.sig.Params().Len() && isNamedNonEmptyInterface(a.sig.Params().At(i).Type()) {
			paramType := a.sig.Params().At(i).Type()
			a.emitArg(func() { a.g.emitExprAsType(a.w, a.node, arg, paramType) })
		} else {
			a.emitArg(func() { a.g.emitCArg(a.w, arg) })
		}
	}
}

// emitExternVarArgs writes the arguments of a call to a variadic extern nodecay function.
func (a *callArgs) emitExternVarArgs() {
	if a.call.Ellipsis.IsValid() {
		a.g.fail(a.call, "spreading variadic arguments to an extern function is not supported")
	}
	for i := range a.call.Args {
		a.emitArg(func() { a.g.emitExternVarArg(a.w, a.node, a.call, a.sig, i) })
	}
}
