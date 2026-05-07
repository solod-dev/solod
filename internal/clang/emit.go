package clang

import (
	"fmt"
	"go/ast"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// EmitOptions holds the options for code generation.
type EmitOptions struct {
	Pkg    *packages.Package
	OutDir string
}

// Emit generates C code for the given Go package and all its subpackages,
// and writes it to the specified output directory. Creates a single header
// file with typedefs (.h) and a single implementation file (.c) for each package.
func Emit(opts EmitOptions) error {
	var err error
	if err = os.MkdirAll(opts.OutDir, 0o755); err != nil {
		return fmt.Errorf("create output directory %s: %w", opts.OutDir, err)
	}

	// Initialize the generator with package information.
	g := newGenerator(opts.Pkg)
	g.collect()

	// Emit header file.
	hPath := filepath.Join(opts.OutDir, g.pkg.Name+".h")
	hFile, err := os.Create(hPath)
	if err != nil {
		return fmt.Errorf("create header file %s: %w", hPath, err)
	}
	defer hFile.Close()
	g.emitHeader(hFile)

	// Emit implementation file.
	cPath := filepath.Join(opts.OutDir, g.pkg.Name+".c")
	cFile, err := os.Create(cPath)
	if err != nil {
		return fmt.Errorf("create C file %s: %w", cPath, err)
	}
	defer cFile.Close()
	g.emitImpl(cFile)

	return nil
}

// State holds the code generation state for the current scope.
type State struct {
	writer io.Writer

	// Current indentation level (number of tabs).
	indent int
	// Current function's signature (for multi-return).
	funcSig *types.Signature
	// Deferred generic calls to emit before returns, panics, and function end.
	defers []string
	// Counter for unique temp variable names.
	tempCount int
	// Whether we are emitting inside a #define macro body.
	inMacro bool
	// Non-type macro parameter names. They are suffixed with _ and parenthesized
	// to avoid name collisions (b->val = val) and syntax errors (&b->val) in macro bodies.
	macroParams map[string]bool
}

// Includes holds the C headers to be included in the emitted .h and .c files.
type Includes struct {
	header []string // so:include -> emitted in .h
	impl   []string // so:include.c -> emitted in .c
}

// Generator is responsible for generating C code from Go ASTs.
type Generator struct {
	pkg      *packages.Package
	types    *types.Info
	state    State
	externs  map[types.Object]externInfo  // symbols provided by C headers
	includes Includes                     // included headers from so:include
	symbols  []symbol                     // pre-collected top-level declarations
	embeds   Embeds                       // embedded C files from so:embed
	comments ast.CommentMap               // all comments across all files
	funcDirs map[*ast.FuncDecl]directives // parsed directives per function decl
	initFunc *ast.FuncDecl                // package init() function, if any
	panicked bool                         // true after first panic caught in Visit
}

// newGenerator creates a new Generator instance.
func newGenerator(pkg *packages.Package) *Generator {
	return &Generator{
		pkg:      pkg,
		types:    pkg.TypesInfo,
		externs:  make(map[types.Object]externInfo),
		funcDirs: make(map[*ast.FuncDecl]directives),
	}
}

// emitHeader creates the .h file with typedefs, includes, and extern declarations.
func (g *Generator) emitHeader(w io.Writer) {
	fmt.Fprintf(w, "#pragma once\n")
	fmt.Fprintf(w, "#include \"so/builtin/builtin.h\"\n")
	for _, inc := range g.includes.header {
		fmt.Fprintf(w, "#include %s\n", inc)
	}
	g.emitImports(w)
	g.emitEmbeds(w, g.embeds.header)
	g.emitHeaderDecls(w)
}

// emitImpl creates the .c implementation file by walking the AST.
func (g *Generator) emitImpl(w io.Writer) {
	g.state.writer = w

	fmt.Fprintf(w, "#include \"%s.h\"\n", g.pkg.Name)
	for _, inc := range g.includes.impl {
		fmt.Fprintf(w, "#include %s\n", inc)
	}

	g.emitEmbeds(w, g.embeds.impl)
	g.emitUnexportedTypes(w)
	g.emitPackageVars(w)
	g.emitForwardFuncDecls(w)

	multiFile := len(g.pkg.Syntax) > 1
	if !multiFile {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "// -- Implementation --")
	}
	for _, file := range g.pkg.Syntax {
		if multiFile {
			pos := g.pkg.Fset.Position(file.Pos())
			fmt.Fprintf(w, "\n// -- %s --\n", filepath.Base(pos.Filename))
		}
		ast.Walk(g, file)
	}
	g.emitInitFunc(w)
}

// emitEmbeds writes the content of embedded files, separated by blank lines.
func (g *Generator) emitEmbeds(w io.Writer, files []embedFile) {
	if len(files) == 0 {
		return
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "// -- Embeds --")
	for _, ef := range files {
		fmt.Fprintf(w, "\n%s\n", strings.TrimRight(ef.content, "\n"))
	}
}

// emitInitFunc emits the package init() function as a GCC constructor
// that runs automatically before main().
func (g *Generator) emitInitFunc(w io.Writer) {
	if g.initFunc == nil {
		return
	}
	decl := g.initFunc
	g.state.funcSig = g.funcSig(decl)
	g.state.tempCount = 0

	fmt.Fprintf(w, "\nstatic void __attribute__((constructor)) %s_init() {\n", g.pkg.Name)
	g.state.indent++
	g.walkStmts(decl.Body.List)
	g.emitDeferredCalls()
	g.state.indent--
	fmt.Fprintf(w, "}\n")

	g.state.defers = nil
	g.state.funcSig = nil
}

// indent returns the current indentation string based on the indent level.
func (g *Generator) indent() string {
	return strings.Repeat("    ", g.state.indent)
}
