package clang

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"maps"
	"os"
	"path/filepath"
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
	dirs     directives   // parsed so: directives
	genDecl  *ast.GenDecl // parent GenDecl (for type symbols, enables comment lookup)
	typeSpec *ast.TypeSpec
	funcDecl *ast.FuncDecl
}

// collect performs a single pass over all package files, collecting:
// - Comment map (for doc comment emission)
// - Include directives (so:include, so:include.c)
// - Embed directives (so:embed) with file reads
// - Extern symbols (so:extern, body-less functions)
// - Symbol list (types, vars, consts, functions) with parsed directives
func (g *Generator) collect() {
	g.comments = ast.CommentMap{}
	g.embeds = Embeds{vars: make(map[string]bool)}

	srcDir := ""
	if len(g.pkg.GoFiles) > 0 {
		srcDir = filepath.Dir(g.pkg.GoFiles[0])
	}

	for _, file := range g.pkg.Syntax {
		// Merge file comments into the global comment map.
		fileComments := ast.NewCommentMap(g.pkg.Fset, file, file.Comments)
		maps.Copy(g.comments, fileComments)

		// Collect include directives from file-level comments.
		for _, cg := range file.Comments {
			for _, c := range cg.List {
				if path, ok := strings.CutPrefix(c.Text, "//so:include.c"); ok {
					g.includes.impl = append(g.includes.impl, strings.TrimSpace(path))
				} else if path, ok := strings.CutPrefix(c.Text, "//so:include"); ok {
					g.includes.header = append(g.includes.header, strings.TrimSpace(path))
				}
			}
		}

		// Collect extern symbols and build the symbol list.
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				g.collectGenDecl(srcDir, d)
			case *ast.FuncDecl:
				g.collectFuncDecl(d)
			}
		}
	}

	// Collect externs from imported packages so that callExtern
	// can identify cross-package extern calls (e.g. stdio.Printf).
	for _, imp := range g.pkg.Imports {
		for _, file := range imp.Syntax {
			g.collectFileExterns(imp.TypesInfo, file)
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

// collectGenDecl processes a GenDecl for externs, embeds, and symbol collection.
func (g *Generator) collectGenDecl(srcDir string, d *ast.GenDecl) {
	// Handle so:extern declarations.
	foundExtern, externInf := parseExtern(d.Doc)
	if foundExtern {
		for _, spec := range d.Specs {
			switch s := spec.(type) {
			case *ast.TypeSpec:
				g.markExtern(g.types.Defs[s.Name], externInf)
				g.markExternFields(g.types, s, externInf)
			case *ast.ValueSpec:
				for _, name := range s.Names {
					g.markExtern(g.types.Defs[name], externInf)
				}
			}
		}
		return
	}

	// Handle so:embed on variable declarations.
	if d.Tok == token.VAR {
		if filename, ok := embedDirective(d.Doc); ok {
			content, err := os.ReadFile(filepath.Join(srcDir, filename))
			if err != nil {
				g.fail(d, "error reading file %s: %v", filename, err)
			}
			ef := embedFile{name: filename, content: string(content)}
			switch filepath.Ext(filename) {
			case ".h":
				g.embeds.header = append(g.embeds.header, ef)
			case ".c":
				g.embeds.impl = append(g.embeds.impl, ef)
			}
			for _, spec := range d.Specs {
				vs := spec.(*ast.ValueSpec)
				for _, name := range vs.Names {
					g.embeds.vars[name.Name] = true
				}
			}
			return
		}
	}

	// Parse directives for non-extern, non-embed GenDecls.
	dirs := parseDirectives(d.Doc)

	// Validate directive/declaration-kind compatibility.
	if dirs.inline {
		g.fail(d, "so:inline is only allowed on functions")
	}
	switch d.Tok {
	case token.TYPE:
		if dirs.volatile {
			g.fail(d, "so:volatile is not allowed on type declarations")
		}
		if dirs.threadLocal {
			g.fail(d, "so:thread_local is not allowed on type declarations")
		}
		for _, spec := range d.Specs {
			ts := spec.(*ast.TypeSpec)
			if g.hasExtern(g.types.Defs[ts.Name]) {
				continue
			}
			g.symbols = append(g.symbols, symbol{
				kind:     symbolType,
				exported: ast.IsExported(ts.Name.Name),
				dirs:     dirs,
				genDecl:  d,
				typeSpec: ts,
			})
		}
	case token.VAR:
		g.symbols = append(g.symbols, symbol{
			kind:     symbolVar,
			exported: hasExportedValueSpec(d),
			dirs:     dirs,
			genDecl:  d,
		})
	case token.CONST:
		if dirs.volatile {
			g.fail(d, "so:volatile is not allowed on const declarations")
		}
		if dirs.threadLocal {
			g.fail(d, "so:thread_local is not allowed on const declarations")
		}
		g.symbols = append(g.symbols, symbol{
			kind:     symbolConst,
			exported: hasExportedValueSpec(d),
			dirs:     dirs,
			genDecl:  d,
		})
	}
}

// collectFuncDecl processes a FuncDecl for externs and symbol collection.
func (g *Generator) collectFuncDecl(d *ast.FuncDecl) {
	// Handle extern functions (body-less or so:extern).
	foundExtern, externInf := parseExtern(d.Doc)
	if d.Body == nil || foundExtern {
		g.markExtern(g.types.Defs[d.Name], externInf)
		return
	}

	if isMainFunc(d) {
		return
	}
	if isInitFunc(d) {
		if g.initFunc != nil {
			g.fail(d.Name, "multiple init functions in package %s", g.pkg.Name)
		}
		g.initFunc = d
		return
	}
	if g.hasExtern(g.types.Defs[d.Name]) {
		return
	}

	dirs := parseDirectives(d.Doc)

	// Validate directive/declaration-kind compatibility.
	if dirs.volatile {
		g.fail(d, "so:volatile is not allowed on functions")
	}
	if dirs.threadLocal {
		g.fail(d, "so:thread_local is not allowed on functions")
	}

	g.funcDirs[d] = dirs

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
		dirs:     dirs,
		funcDecl: d,
	})
}

// emitPackageVars writes all package-level variable and constant
// declarations at the top of the .c file, before forward declarations.
// This ensures they are defined before any function that references them.
func (g *Generator) emitPackageVars(w io.Writer) {
	var symbols []symbol
	for _, sym := range g.symbols {
		if sym.kind != symbolVar && sym.kind != symbolConst {
			continue
		}
		symbols = append(symbols, sym)
	}
	if len(symbols) == 0 {
		return
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "// -- Variables and constants --")
	for _, sym := range symbols {
		g.emitComments(w, sym.genDecl)
		switch sym.genDecl.Tok {
		case token.CONST:
			for _, spec := range sym.genDecl.Specs {
				g.emitConstSpec(spec.(*ast.ValueSpec))
			}
		case token.VAR:
			for _, spec := range sym.genDecl.Specs {
				vs := spec.(*ast.ValueSpec)
				if len(vs.Names) > 0 && g.embeds.vars[vs.Names[0].Name] {
					continue
				}
				g.emitVarSpec(vs, sym.dirs)
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
	g.emitForwardTypeDecls(w, typeSyms)
	for _, sym := range typeSyms {
		hasDocs := g.emitComments(w, sym.genDecl, sym.typeSpec)
		if !hasDocs && isBlockTypeSpec(sym.typeSpec) {
			fmt.Fprintln(w)
		}
		g.emitTypeSpec(w, sym.typeSpec, sym.dirs)
	}
}

// emitForwardTypeDecls writes forward declarations for struct types
// so that self-referencing and out-of-order references resolve.
func (g *Generator) emitForwardTypeDecls(w io.Writer, typeSyms []symbol) {
	hasDecls := false
	for _, sym := range typeSyms {
		if _, ok := sym.typeSpec.Type.(*ast.StructType); ok {
			cName := g.declSymbolName(g.types.Defs[sym.typeSpec.Name])
			fmt.Fprintf(w, "\ntypedef struct %s %s;", cName, cName)
			hasDecls = true
		}
	}
	if hasDecls {
		fmt.Fprintln(w)
	}
}

// emitForwardFuncDecls writes forward declarations for unexported functions
// and methods so that they can be called before their definition.
func (g *Generator) emitForwardFuncDecls(w io.Writer) {
	var funcDecls []*ast.FuncDecl
	for _, sym := range g.symbols {
		if sym.kind != symbolFunc && sym.kind != symbolMethod {
			continue
		}
		if sym.exported || sym.dirs.inline {
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
