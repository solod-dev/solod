package clang

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"strings"
)

type symbolKind int

const (
	symbolFunc symbolKind = iota
	symbolMethod
	symbolType
	symbolVar
	symbolConst
)

type symbol struct {
	kind     symbolKind
	exported bool
	genDecl  *ast.GenDecl // parent GenDecl (for type symbols, enables comment lookup)
	typeSpec *ast.TypeSpec
	funcDecl *ast.FuncDecl
}

// collectSymbols gathers all top-level type, function, var, and const
// declarations into an ordered list. This list drives header emission
// (exported symbols), forward declarations, and hoisted vars/consts.
func (g *Generator) collectSymbols() {
	// Collect top-level types and functions and their export status.
	for _, file := range g.pkg.Syntax {
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				if found, _ := parseExternDirective(d.Doc); found {
					continue
				}
				switch d.Tok {
				case token.TYPE:
					for _, spec := range d.Specs {
						ts := spec.(*ast.TypeSpec)
						if g.hasExtern("", ts.Name.Name) {
							continue
						}
						g.symbols = append(g.symbols, symbol{
							kind:     symbolType,
							exported: ast.IsExported(ts.Name.Name),
							genDecl:  d,
							typeSpec: ts,
						})
					}
				case token.VAR:
					g.symbols = append(g.symbols, symbol{
						kind:     symbolVar,
						exported: hasExportedValueSpec(d),
						genDecl:  d,
					})
				case token.CONST:
					g.symbols = append(g.symbols, symbol{
						kind:     symbolConst,
						exported: hasExportedValueSpec(d),
						genDecl:  d,
					})
				}
			case *ast.FuncDecl:
				if d.Body == nil || d.Name.Name == "main" {
					continue
				}
				if g.hasExtern("", externFuncKey(d)) {
					continue
				}
				kind := symbolFunc
				exported := ast.IsExported(d.Name.Name)
				if d.Recv != nil {
					kind = symbolMethod
					if exported {
						exported = ast.IsExported(recvTypeName(d.Recv.List[0]))
					}
				}
				g.symbols = append(g.symbols, symbol{
					kind:     kind,
					exported: exported,
					funcDecl: d,
				})
			}
		}
	}

	// Validate that exported functions don't use unexported types.
	for _, sym := range g.symbols {
		if !sym.exported || (sym.kind != symbolFunc && sym.kind != symbolMethod) {
			continue
		}
		decl := sym.funcDecl
		if g.hasUnexportedTypes(decl) {
			g.fail(decl.Name, "exported function %s uses unexported types", decl.Name.Name)
		}
	}
}

// collectExterns scans all files for extern symbols and include directives.
// Body-less functions and declarations annotated with //so:extern are treated
// as external C symbols that should not be emitted.
func (g *Generator) collectExterns() {
	for _, file := range g.pkg.Syntax {
		// Collect include directives from the file.
		for _, cg := range file.Comments {
			for _, c := range cg.List {
				if include, ok := strings.CutPrefix(c.Text, "//so:include"); ok {
					g.includes = append(g.includes, strings.TrimSpace(include))
				}
			}
		}

		// Collect extern symbols from declarations.
		g.collectFileExterns("", file)
	}

	// Collect externs from imported packages so that callExtern
	// can identify cross-package extern calls (e.g. stdio.Printf).
	for _, imp := range g.pkg.Imports {
		for _, file := range imp.Syntax {
			g.collectFileExterns(imp.Name, file)
		}
	}
}

// emitPackageVars writes all package-level variable and constant
// declarations at the top of the .c file, before forward declarations.
// This ensures they are defined before any function that references them.
func (g *Generator) emitPackageVars(w io.Writer) {
	var decls []*ast.GenDecl
	for _, sym := range g.symbols {
		if sym.kind != symbolVar && sym.kind != symbolConst {
			continue
		}
		decls = append(decls, sym.genDecl)
	}
	if len(decls) == 0 {
		return
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "// -- Variables and constants --")
	for _, decl := range decls {
		g.emitComments(w, decl)
		switch decl.Tok {
		case token.CONST:
			for _, spec := range decl.Specs {
				g.emitConstSpec(spec.(*ast.ValueSpec))
			}
		case token.VAR:
			for _, spec := range decl.Specs {
				vs := spec.(*ast.ValueSpec)
				if len(vs.Names) > 0 && g.embeds.vars[vs.Names[0].Name] {
					continue
				}
				g.emitVarSpec(vs)
			}
		}
	}
}

// emitUnexportedTypes writes full type definitions for all unexported types.
// Emitted before package vars so that compound literals can reference them.
func (g *Generator) emitUnexportedTypes(w io.Writer) {
	var typeSyms []symbol
	for _, sym := range g.symbols {
		if sym.exported || sym.kind != symbolType {
			continue
		}
		typeSyms = append(typeSyms, sym)
	}
	if len(typeSyms) == 0 {
		return
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "// -- Types --")
	for _, sym := range typeSyms {
		hasDocs := g.emitComments(w, sym.genDecl, sym.typeSpec)
		if !hasDocs && isBlockTypeSpec(sym.typeSpec) {
			fmt.Fprintln(w)
		}
		g.emitTypeSpec(w, sym.typeSpec)
	}
}

// emitForwardFuncDecls writes forward declarations for unexported functions
// and methods so that they can be called before their definition.
func (g *Generator) emitForwardFuncDecls(w io.Writer) {
	var funcDecls []*ast.FuncDecl
	for _, sym := range g.symbols {
		if sym.exported || (sym.kind != symbolFunc && sym.kind != symbolMethod) {
			continue
		}
		funcDecls = append(funcDecls, sym.funcDecl)
	}
	if len(funcDecls) == 0 {
		return
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "// -- Forward declarations --")
	for _, decl := range funcDecls {
		g.emitFuncProto(w, decl)
		fmt.Fprintln(w, ";")
	}
}

// hasExportedValueSpec reports whether a GenDecl contains at least one
// exported name in its value specs.
func hasExportedValueSpec(d *ast.GenDecl) bool {
	for _, spec := range d.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		for _, name := range vs.Names {
			if ast.IsExported(name.Name) {
				return true
			}
		}
	}
	return false
}
