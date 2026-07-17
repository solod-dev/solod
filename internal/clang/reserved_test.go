package clang

import (
	"go/ast"
	"go/token"
	"go/types"
	"testing"
)

func TestScopeDeclares(t *testing.T) {
	// Scope tree: root -> mid -> leaf, with root -> sibling.
	root := types.NewScope(nil, token.NoPos, token.NoPos, "root")
	mid := types.NewScope(root, token.NoPos, token.NoPos, "mid")
	leaf := types.NewScope(mid, token.NoPos, token.NoPos, "leaf")
	sibling := types.NewScope(root, token.NoPos, token.NoPos, "sibling")

	insert := func(s *types.Scope, name string) {
		s.Insert(types.NewVar(token.NoPos, nil, name, nil))
	}
	insert(root, "inAncestor")
	insert(mid, "inSelf")
	insert(leaf, "inDescendant")
	insert(sibling, "inSibling")

	tests := []struct {
		name string
		want bool
	}{
		{"inSelf", true},        // same scope: real collision
		{"inAncestor", false},   // enclosing scope: legal shadow
		{"inDescendant", false}, // nested scope: legal shadow
		{"inSibling", false},    // disjoint sibling scope
		{"missing", false},      // not declared anywhere
	}
	for _, tt := range tests {
		if got := scopeDeclares(mid, tt.name); got != tt.want {
			t.Errorf("scopeDeclares(mid, %q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestReservedLabel(t *testing.T) {
	// A reserved-word label is rejected rather than mangled: go/types records
	// the label as a *types.Label in both Defs and Uses, and isLocal returns
	// false for it. handleReservedNames then rejects in (fail at the label def).
	src := `package x
func f() {
	long:
	for {
		break long
	}
	goto long
}`
	info, file := checkSnippet(t, src)
	pkgScope := types.NewScope(types.Universe, 0, 0, "x")

	var defObj, useObj types.Object
	ast.Inspect(file, func(n ast.Node) bool {
		id, ok := n.(*ast.Ident)
		if !ok || id.Name != "long" {
			return true
		}
		if o := info.Defs[id]; o != nil {
			defObj = o
		}
		if o := info.Uses[id]; o != nil {
			useObj = o
		}
		return true
	})

	if defObj == nil {
		t.Fatal("label def not recorded in Defs; expected a *types.Label")
	}
	if useObj == nil {
		t.Fatal("label use not recorded in Uses; expected a *types.Label")
	}
	if _, ok := defObj.(*types.Label); !ok {
		t.Fatalf("label def object = %T, want *types.Label", defObj)
	}
	if isLocal(defObj, pkgScope) {
		t.Error("isLocal(label) = true; a label must not be treated as a manglable local")
	}
}
