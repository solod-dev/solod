package clang

import (
	"fmt"
	"go/ast"
	"io"
)

// emitSwitchStmt emits a switch statement, wrapping
// in a scope block if there's an init statement.
func (g *Generator) emitSwitchStmt(stmt *ast.SwitchStmt) {
	w := g.state.writer
	if stmt.Init != nil {
		fmt.Fprintf(w, "%s{\n", g.indent())
		g.state.indent++
		ast.Walk(g, stmt.Init)
		g.emitSwitchBody(w, stmt)
		g.state.indent--
		fmt.Fprintf(w, "%s}\n", g.indent())
	} else {
		g.emitSwitchBody(w, stmt)
	}
}

// emitSwitchBody emits the if/else-if/else chain for a switch statement.
func (g *Generator) emitSwitchBody(w io.Writer, stmt *ast.SwitchStmt) {
	var cases []*ast.CaseClause
	var def *ast.CaseClause
	for _, s := range stmt.Body.List {
		cc := s.(*ast.CaseClause)
		if cc.List == nil {
			def = cc
		} else {
			cases = append(cases, cc)
		}
	}

	// Empty switch.
	if len(cases) == 0 && def == nil {
		fmt.Fprintf(w, "%sif (false) {\n%s}\n", g.indent(), g.indent())
		return
	}

	// Default-only.
	if len(cases) == 0 {
		g.walkStmts(def.Body)
		return
	}

	// Emit if/else-if chain.
	isString := stmt.Tag != nil && g.hasStringType(stmt.Tag)
	for i, cc := range cases {
		if i == 0 {
			fmt.Fprintf(w, "%sif (", g.indent())
		} else {
			fmt.Fprintf(w, "%s} else if (", g.indent())
		}
		for j, expr := range cc.List {
			if j > 0 {
				fmt.Fprintf(w, " || ")
			}
			if stmt.Tag == nil {
				g.emitExpr(expr)
			} else if isString {
				fmt.Fprintf(w, "so_string_eq(")
				g.emitExpr(stmt.Tag)
				fmt.Fprintf(w, ", ")
				g.emitExpr(expr)
				fmt.Fprintf(w, ")")
			} else {
				g.emitExpr(stmt.Tag)
				fmt.Fprintf(w, " == (")
				g.emitExpr(expr)
				fmt.Fprintf(w, ")")
			}
		}
		fmt.Fprintf(w, ") {\n")
		g.state.indent++
		g.walkStmts(cc.Body)
		g.state.indent--
	}

	if def != nil {
		fmt.Fprintf(w, "%s} else {\n", g.indent())
		g.state.indent++
		g.walkStmts(def.Body)
		g.state.indent--
	}
	fmt.Fprintf(w, "%s}\n", g.indent())
}
