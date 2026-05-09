package clang

import (
	"fmt"
	"go/ast"
	"go/token"
)

// emitSelectStmt emits a select statement as a libmill choose block.
func (g *Generator) emitSelectStmt(stmt *ast.SelectStmt) {
	w := g.state.writer

	// Emit: choose {
	fmt.Fprintf(w, "%schoose {\n", g.indent())
	g.state.indent++

	// Emit each case
	for _, clause := range stmt.Body.List {
		cc := clause.(*ast.CommClause)
		g.emitCommClause(cc)
	}

	g.state.indent--
	fmt.Fprintf(w, "%send\n", g.indent())
}

// emitCommClause emits a single case in a select statement.
func (g *Generator) emitCommClause(clause *ast.CommClause) {
	w := g.state.writer

	if clause.Comm == nil {
		// default case → otherwise:
		fmt.Fprintf(w, "%sotherwise:\n", g.indent())
		g.state.indent++
		g.walkStmts(clause.Body)
		g.state.indent--
		return
	}

	switch comm := clause.Comm.(type) {
	case *ast.SendStmt:
		// case ch <- val:
		g.emitSelectSendCase(comm, clause.Body)

	case *ast.ExprStmt:
		// case <-ch:
		if unary, ok := comm.X.(*ast.UnaryExpr); ok && unary.Op == token.ARROW {
			g.emitSelectRecvCase(unary, clause.Body, nil)
		} else {
			g.fail(comm, "invalid select case")
		}

	case *ast.AssignStmt:
		// case v := <-ch: or case v, ok := <-ch:
		if unary, ok := comm.Rhs[0].(*ast.UnaryExpr); ok && unary.Op == token.ARROW {
			g.emitSelectRecvCase(unary, clause.Body, comm.Lhs)
		} else {
			g.fail(comm, "invalid select case")
		}

	default:
		g.fail(clause, "unsupported select case type: %T", comm)
	}
}

// emitSelectSendCase emits a send case: case ch <- val:
func (g *Generator) emitSelectSendCase(send *ast.SendStmt, body []ast.Stmt) {
	w := g.state.writer

	// Get channel element type
	chType := g.types.TypeOf(send.Chan)
	elemType := chanElemType(chType)
	if elemType == nil {
		g.fail(send, "send on non-channel type")
	}

	// Emit: out(ch, T, val):
	fmt.Fprintf(w, "%sout(", g.indent())
	g.emitExpr(send.Chan)
	fmt.Fprintf(w, ", %s, ", g.mapType(send.Value, elemType))
	g.emitExpr(send.Value)
	fmt.Fprintf(w, "):\n")

	// Emit body
	g.state.indent++
	g.walkStmts(body)
	g.state.indent--
}

// emitSelectRecvCase emits a receive case: case v := <-ch: or case <-ch:
func (g *Generator) emitSelectRecvCase(recv *ast.UnaryExpr, body []ast.Stmt, lhs []ast.Expr) {
	w := g.state.writer

	// Get channel element type
	chType := g.types.TypeOf(recv.X)
	elemType := chanElemType(chType)
	if elemType == nil {
		g.fail(recv, "receive from non-channel type")
	}

	// Determine variable name for received value
	varName := "_"
	if len(lhs) > 0 {
		if ident, ok := lhs[0].(*ast.Ident); ok {
			varName = ident.Name
		}
	}

	// Emit: in(ch, T, varName):
	fmt.Fprintf(w, "%sin(", g.indent())
	g.emitExpr(recv.X)
	fmt.Fprintf(w, ", %s, %s):\n", g.mapType(recv, elemType), varName)

	// Emit body
	g.state.indent++

	// If there's a second variable (ok), set it to true
	// TODO: Implement proper closed-channel detection
	if len(lhs) > 1 {
		if ident, ok := lhs[1].(*ast.Ident); ok {
			fmt.Fprintf(w, "%s%s = true;\n", g.indent(), ident.Name)
		}
	}

	g.walkStmts(body)
	g.state.indent--
}
