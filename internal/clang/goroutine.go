package clang

import (
	"fmt"
	"go/ast"
	"go/types"
)

// emitGoStmt emits a go statement as a libmill go() macro call.
func (g *Generator) emitGoStmt(stmt *ast.GoStmt) {
	w := g.state.writer

	// Emit libmill go() macro
	fmt.Fprintf(w, "%sgo(", g.indent())

	// Mark the function as a coroutine
	g.markAsCoroutine(stmt.Call)

	// Emit the function call
	g.emitCallExpr(stmt.Call)

	fmt.Fprintf(w, ");\n")
}

// markAsCoroutine marks a function as needing the coroutine specifier.
func (g *Generator) markAsCoroutine(call *ast.CallExpr) {
	// Determine the function name
	var funcName string

	switch fun := call.Fun.(type) {
	case *ast.Ident:
		// Direct function call: go myFunc()
		if obj := g.types.Uses[fun]; obj != nil {
			funcName = g.symbolName(obj)
		}
	case *ast.SelectorExpr:
		// Package-qualified call: go pkg.MyFunc()
		if sel, ok := g.types.Uses[fun.Sel].(*types.Func); ok {
			if sel.Pkg() != nil {
				funcName = sel.Pkg().Name() + "_" + sel.Name()
			} else {
				funcName = sel.Name()
			}
		}
	}

	if funcName != "" {
		if g.coroutineFuncs == nil {
			g.coroutineFuncs = make(map[string]bool)
		}
		g.coroutineFuncs[funcName] = true
	}
}

// isCoroutine checks if a function should be marked as a coroutine.
func (g *Generator) isCoroutine(funcName string) bool {
	return g.coroutineFuncs != nil && g.coroutineFuncs[funcName]
}
