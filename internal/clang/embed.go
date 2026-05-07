package clang

import (
	"go/ast"
	"strings"
)

// embedFile represents a single embedded file.
type embedFile struct {
	name    string
	content string
}

// Embeds holds the embedded .h and .c files.
type Embeds struct {
	header []embedFile     // .h contents to inline in header
	impl   []embedFile     // .c contents to inline in impl
	vars   map[string]bool // var names to skip during emission
}

// embedDirective extracts the filename from a so:embed directive.
func embedDirective(doc *ast.CommentGroup) (string, bool) {
	if doc == nil {
		return "", false
	}
	for _, c := range doc.List {
		if filename, ok := strings.CutPrefix(strings.TrimSpace(c.Text), "//so:embed "); ok {
			return strings.TrimSpace(filename), true
		}
	}
	return "", false
}
