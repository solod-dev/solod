package clang

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"
)

// emitSendStmt emits a channel send operation: ch <- val
func (g *Generator) emitSendStmt(stmt *ast.SendStmt) {
	w := g.state.writer

	// Get channel type to determine element type
	chType := g.types.TypeOf(stmt.Chan)
	elemType := chanElemType(chType)
	if elemType == nil {
		g.fail(stmt, "send on non-channel type")
	}

	// Emit: chs(ch, T, val);
	fmt.Fprintf(w, "%schs(", g.indent())
	g.emitExpr(stmt.Chan)
	fmt.Fprintf(w, ", %s, ", g.mapType(stmt.Value, elemType))
	g.emitExpr(stmt.Value)
	fmt.Fprintf(w, ");\n")
}

// emitChanRecv emits a channel receive expression: <-ch
func (g *Generator) emitChanRecv(expr *ast.UnaryExpr) {
	w := g.state.writer

	// Get channel type to determine element type
	chType := g.types.TypeOf(expr.X)
	elemType := chanElemType(chType)
	if elemType == nil {
		g.fail(expr, "receive from non-channel type")
	}

	// Emit: chr(ch, T)
	fmt.Fprintf(w, "chr(")
	g.emitExpr(expr.X)
	fmt.Fprintf(w, ", %s)", g.mapType(expr, elemType))
}

// emitMakeChan emits channel creation: make(chan T) or make(chan T, size)
func (g *Generator) emitMakeChan(call *ast.CallExpr, chanType *types.Chan) {
	w := g.state.writer
	elemType := chanType.Elem()

	// Determine buffer size
	size := "0" // unbuffered by default
	if len(call.Args) > 1 {
		// Buffered channel: make(chan T, size)
		// We need to emit the size expression
		var buf strings.Builder
		oldWriter := g.state.writer
		g.state.writer = &buf
		g.emitExpr(call.Args[1])
		g.state.writer = oldWriter
		size = buf.String()
	}

	// Emit: chmake(T, size)
	fmt.Fprintf(w, "chmake(%s, %s)", g.mapType(call, elemType), size)
}

// emitCloseChan emits channel close: close(ch)
func (g *Generator) emitCloseChan(call *ast.CallExpr) {
	w := g.state.writer

	if len(call.Args) != 1 {
		g.fail(call, "close() requires exactly one argument")
	}

	// Emit: chclose(ch)
	fmt.Fprintf(w, "chclose(")
	g.emitExpr(call.Args[0])
	fmt.Fprintf(w, ")")
}
