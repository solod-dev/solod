package clang

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"
)

// checkSnippet type-checks a self-contained package and returns its type info
// and the parsed file. The snippet must use only builtin types (no imports).
func checkSnippet(t *testing.T, src string) (*types.Info, *ast.File) {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "x.go", src, 0)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	info := &types.Info{
		Types: map[ast.Expr]types.TypeAndValue{},
		Defs:  map[*ast.Ident]types.Object{},
		Uses:  map[*ast.Ident]types.Object{},
	}
	conf := types.Config{Error: func(error) {}} // tolerate unused-var noise
	conf.Check("x", fset, []*ast.File{f}, info)
	return info, f
}

func funcByName(f *ast.File, name string) *ast.FuncDecl {
	for _, d := range f.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok && fd.Name.Name == name {
			return fd
		}
	}
	return nil
}

func TestFindReturnEscapes(t *testing.T) {
	// Each function is named after whether it should be flagged.
	const src = `package x

type Pair struct{ s string }
type Boxed struct{ a [3]int }
type Arr [3]int

// --- should be flagged (returns frame-bound memory) ---

func escConcat(a, b string) string        { return a + b }
func escConcatLit(a string) string        { return "x" + a }
func escConcatVar(a, b string) string     { s := a + b; return s }
func escConcatChain(a, b string) string   { t := a + b; s := t; return s }
func escAddAssign(s, a string) string     { s += a; return s }
func escRunes(s string) []rune            { return []rune(s) }
func escRunesToStr(rs []rune) string      { return string(rs) }
func escByteToStr(b byte) string          { return string(b) }
func escRuneToStr(r rune) string          { return string(r) }
func escArray() [2][3]int                 { return [2][3]int{} }
func escArrayPtr() *[3]int                { a := [3]int{}; return &a }
func escLocalArray() [3]int               { var a [3]int; return a }
func escLocalArrayDef() [3]int            { a := [3]int{}; return a }
func escNamedArray() Arr                  { var a Arr; return a }
func escMakeSlice(n int) []int            { return make([]int, n) }
func escMakeMap() map[string]int          { return make(map[string]int, 4) }
func escAddrLocal() *int                  { x := 5; return &x }
func escAddrParam(p int) *int             { return &p }
func escAddrComposite() *Pair             { return &Pair{} }
func escStruct(a, b string) Pair          { return Pair{s: a + b} }
func escSliceLit(a, b string) []string    { return []string{a + b} }
func escSliceLitConst() []int             { return []int{1, 2, 3} }
func escSliceLitVar() []int               { s := []int{1, 2, 3}; return s }
func escSubslice(a, b string) string      { s := a + b; return s[1:] }
func escAppend() []int                    { s := make([]int, 0, 4); return append(s, 1) }
func escAppendFrameDst() []byte           { var a [4]byte; s := a[:0]; return append(s, 1) }

// --- should NOT be flagged (safe) ---

func okConstConcat() string               { return "a" + "b" }
func okReturnParam(s string) string       { return s }
func okReturnPtrParam(p *int) *int        { return p }
func okParamSubslice(s string) string     { return s[1:] }
func okBytesFromParam(s string) []byte    { return []byte(s) }
func okStrFromByteParam(bs []byte) string { return string(bs) }
func okLenOfConcat(a, b string) int       { return len(a + b) }
func okReassignSafe(a, b, c string) string { s := a + b; s = c; _ = s; return c }
func okArrayParam(a [3]int) [3]int        { return a }
func okNamedArrayParam(a Arr) Arr         { return a }
func okMapParam(m map[string]int) map[string]int { return m }
func okSliceParam(s []int) []int          { return s }
func okEmptySliceLit() []int              { return []int{} }
func okStructParam(p Pair) Pair           { return p }
func okStructPtrParam(p *Pair) *Pair      { return p }
func okStructWithArray() Boxed            { return Boxed{a: [3]int{1, 2, 3}} }
func okAppendParam(dst []byte, a [3]byte) []byte { return append(dst, a[:]...) }

// --- methods (checked like functions) ---

func (p Pair) escMethodConcat(a string) string { return p.s + a }
func (p *Pair) escMethodAddr() *Pair           { return &Pair{} }
func (p Pair) okMethodField() string           { return p.s }
func (p *Pair) okMethodRecv() *Pair            { return p }
`

	info, f := checkSnippet(t, src)

	for _, d := range f.Decls {
		fd, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}
		name := fd.Name.Name
		wantEscape := len(name) >= 3 && name[:3] == "esc"
		got := findReturnEscapes(info, fd)
		if wantEscape && len(got) == 0 {
			t.Errorf("%s: expected an escape, found none", name)
		}
		if !wantEscape && len(got) != 0 {
			t.Errorf("%s: expected no escape, found %d", name, len(got))
		}
	}
}
