package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

// methodNameOverridesMap maps from generated method name to desired method name.
// This is used to fix up names that would be awkward or conflict with existing methods.
var methodNameOverridesMap = map[string]string{
	"UapSiaDb": "UapDb",
	"UapSiaVm": "UapVm",
}

type serviceInfo struct {
	Dir        string // filesystem path
	PkgName    string // package name
	ImportPath string // module import path
	Alias      string // import alias
	MethodName string // exported method name
	RetExpr    string // rendered return type of ServiceGenerator
}

func main() {
	moduleRoot, modulePath := mustFindModuleRootAndPath()
	servicesRoot := filepath.Join(moduleRoot, "pkg", "services")

	var svcs []serviceInfo
	seenAlias := map[string]bool{}

	_ = filepath.WalkDir(servicesRoot, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			// unreadable directory; skip
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		name := d.Name()
		// Skip common non-source dirs
		if name == "vendor" || name == "testdata" || strings.HasPrefix(name, ".") {
			return filepath.SkipDir
		}

		// Try to parse this directory as a Go package and find ServiceGenerator
		pkgName, retExpr := findServiceGeneratorSignature(p)
		if pkgName == "" || retExpr == "" {
			return nil // not a target package
		}

		rel, err := filepath.Rel(moduleRoot, p)
		if err != nil {
			return nil
		}
		if strings.Contains(rel, "models") {
			return nil // skip "models" dirs
		}

		// Create a readable method name from its path under pkg/services/
		methodName := exportPathAsMethod(strings.TrimPrefix(filepath.ToSlash(rel), "pkg/services/"))
		if override, ok := methodNameOverridesMap[methodName]; ok {
			methodName = override
		}

		// Build a unique alias (start from package name, then add numeric suffix)
		alias := pkgName
		base := alias
		for i := 2; seenAlias[alias]; i++ {
			alias = fmt.Sprintf("%s%d", base, i)
		}
		seenAlias[alias] = true

		svcs = append(svcs, serviceInfo{
			Dir:        p,
			PkgName:    pkgName,
			ImportPath: fmt.Sprintf("%s/%s", modulePath, filepath.ToSlash(rel)),
			Alias:      alias,
			MethodName: methodName,
			RetExpr:    retExpr,
		})
		return nil
	})

	// Qualify return types with the package alias
	for i := range svcs {
		svcs[i].RetExpr = qualifyTypeExpr(svcs[i].RetExpr, svcs[i].Alias)
	}

	// Sort methods for stable output
	sort.Slice(svcs, func(i, j int) bool { return svcs[i].MethodName < svcs[j].MethodName })

	importBlock := buildImportBlock(svcs)
	methods := buildMethods(svcs)

	// Load template
	tmplPath := filepath.Join(moduleRoot, "tools", "genapi", "api", "ark_api.tmpl.go")
	tpl, err := os.ReadFile(tmplPath)
	if err != nil {
		panic(err)
	}
	out := string(tpl)
	out = replaceMarker(out, "// +gen:imports", importBlock)
	out = replaceMarker(out, "// +gen:methods", methods)

	outPath := filepath.Join(moduleRoot, "pkg", "ark_api.go")
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(outPath, []byte(out), 0o644); err != nil {
		panic(err)
	}
}

func buildFuncMap(pkg *ast.Package) map[string]*ast.FuncType {
	funcs := make(map[string]*ast.FuncType)
	for _, f := range pkg.Files {
		for _, d := range f.Decls {
			if fd, ok := d.(*ast.FuncDecl); ok && fd.Recv == nil && fd.Type != nil && fd.Name != nil {
				funcs[fd.Name.Name] = fd.Type
			}
		}
	}
	return funcs
}

func findServiceGeneratorSignature(dir string) (pkgName, concreteRet string) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		n := fi.Name()
		return strings.HasSuffix(n, ".go") && !strings.HasSuffix(n, "_test.go")
	}, 0)
	if err != nil || len(pkgs) == 0 {
		return "", ""
	}

	var pkg *ast.Package
	for _, p := range pkgs {
		pkg = p
		break
	}
	pkgName = pkg.Name

	funcs := buildFuncMap(pkg)

	// Look for: var ServiceGenerator = MyFunc
	for _, f := range pkg.Files {
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.VAR {
				continue
			}
			// TODO: Also verify ServiceConfig exists so that we will know the service was exposed properly
			for _, spec := range gd.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok || !hasName(vs.Names, "ServiceGenerator") {
					continue
				}
				// We only support assignment to an ident: MyFunc
				for _, v := range vs.Values {
					if id, ok := v.(*ast.Ident); ok {
						if ft, ok := funcs[id.Name]; ok && hasErrorAsSecond(ft) {
							if r := firstConcreteReturn(ft); r != "" {
								return pkgName, r
							}
						}
					}
				}
			}
		}
	}
	return "", ""
}

func hasName(list []*ast.Ident, want string) bool {
	for _, n := range list {
		if n.Name == want {
			return true
		}
	}
	return false
}

func firstConcreteReturn(ft *ast.FuncType) string {
	if ft.Results == nil || len(ft.Results.List) < 1 {
		return ""
	}
	return renderType(ft.Results.List[0].Type)
}

func hasErrorAsSecond(ft *ast.FuncType) bool {
	if ft.Results == nil || len(ft.Results.List) < 2 {
		return false
	}
	if id, ok := ft.Results.List[1].Type.(*ast.Ident); ok && id.Name == "error" {
		return true
	}
	return false
}

func renderType(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + renderType(t.X)
	case *ast.ArrayType:
		if t.Len == nil {
			return "[]" + renderType(t.Elt)
		}
		return "[" + renderExpr(t.Len) + "]" + renderType(t.Elt)
	case *ast.MapType:
		return "map[" + renderType(t.Key) + "]" + renderType(t.Value)
	case *ast.SelectorExpr:
		return renderExpr(t.X) + "." + t.Sel.Name
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.ChanType:
		prefix := "chan "
		if t.Dir == ast.SEND {
			prefix = "chan<- "
		} else if t.Dir == ast.RECV {
			prefix = "<-chan "
		}
		return prefix + renderType(t.Value)
	case *ast.Ellipsis:
		return "..." + renderType(t.Elt)
	case *ast.FuncType:
		var b bytes.Buffer
		b.WriteString("func(")
		if t.Params != nil {
			for i, f := range t.Params.List {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(renderType(f.Type))
			}
		}
		b.WriteString(")")
		if t.Results != nil && len(t.Results.List) > 0 {
			b.WriteString(" (")
			for i, f := range t.Results.List {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(renderType(f.Type))
			}
			b.WriteString(")")
		}
		return b.String()
	default:
		return fmt.Sprintf("%T", e)
	}
}

func renderExpr(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.BasicLit:
		return t.Value
	case *ast.SelectorExpr:
		return renderExpr(t.X) + "." + t.Sel.Name
	default:
		return renderType(e)
	}
}

var predecl = map[string]bool{
	"bool": true, "byte": true, "complex64": true, "complex128": true,
	"error": true, "float32": true, "float64": true, "int": true,
	"int8": true, "int16": true, "int32": true, "int64": true,
	"rune": true, "string": true, "uint": true, "uint8": true,
	"uint16": true, "uint32": true, "uint64": true, "uintptr": true,
	"any": true,
}

func qualifyTypeExpr(s, alias string) string {
	type tok struct{ v string }
	var toks []tok

	seps := "[]*(),{} \t\r\n"
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || strings.ContainsRune(seps, rune(s[i])) {
			if start < i {
				toks = append(toks, tok{s[start:i]})
			}
			if i < len(s) {
				toks = append(toks, tok{s[i : i+1]})
			}
			start = i + 1
		}
	}

	var out strings.Builder
	for i := 0; i < len(toks); i++ {
		t := toks[i].v
		if isIdent(t) && !predecl[t] {
			if i+1 < len(toks) && toks[i+1].v == "." {
				out.WriteString(t)
				continue
			}
			out.WriteString(alias)
			out.WriteString(".")
			out.WriteString(t)
			continue
		}
		out.WriteString(t)
	}
	return out.String()
}

func isIdent(s string) bool {
	if s == "" {
		return false
	}
	runes := []rune(s)
	if !(s[0] == '_' || unicode.IsLetter(runes[0])) {
		return false
	}
	for _, r := range runes[1:] {
		if !(r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)) {
			return false
		}
	}
	return true
}

func buildImportBlock(svcs []serviceInfo) string {
	if len(svcs) == 0 {
		return ""
	}
	var b bytes.Buffer
	for _, s := range svcs {
		fmt.Fprintf(&b, "\n\t%s \"%s\"", s.Alias, s.ImportPath)
	}
	return b.String()
}

func buildMethods(svcs []serviceInfo) string {
	var b bytes.Buffer
	for _, s := range svcs {
		fmt.Fprintf(&b, "func (api *ArkAPI) %s() (%s, error) {\n", s.MethodName, s.RetExpr)
		fmt.Fprintf(&b, "\tif serviceIfs, ok := api.services[%s.ServiceConfig.ServiceName]; ok {\n", s.Alias)
		fmt.Fprintf(&b, "\t\treturn (*serviceIfs).(%s), nil\n", s.RetExpr)
		fmt.Fprintf(&b, "\t}\n")
		fmt.Fprintf(&b, "\tservice, err := %s.ServiceGenerator(api.loadServiceAuthenticators(%s.ServiceConfig)...)\n", s.Alias, s.Alias)
		fmt.Fprintf(&b, "\tif err != nil {\n")
		fmt.Fprintf(&b, "\t\treturn nil, err\n")
		fmt.Fprintf(&b, "\t}\n")
		fmt.Fprintf(&b, "\tvar baseService services.ArkService = service\n")
		fmt.Fprintf(&b, "\tapi.services[%s.ServiceConfig.ServiceName] = &baseService\n", s.Alias)
		fmt.Fprintf(&b, "\treturn service, nil\n")
		b.WriteString("}\n\n")
	}
	return b.String()
}

func exportPathAsMethod(rel string) string {
	// rel like: "x", "x/sub", "x/sub/deeper"
	parts := strings.Split(rel, "/")
	for i, p := range parts {
		parts[i] = exportName(p)
	}
	return strings.Join(parts, "")
}

func exportName(name string) string {
	if name == "" {
		return ""
	}
	r := []rune(name)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func replaceMarker(src, marker, insertion string) string {
	lines := strings.Split(src, "\n")
	var out []string
	found := false

	for _, line := range lines {
		if strings.TrimSpace(line) == strings.TrimSpace(marker) {
			if insertion != "" {
				out = append(out, insertion)
			}
			found = true
			continue // drop the marker line
		}
		out = append(out, line)
	}

	// If marker wasn't found, append insertion at the end.
	if !found && insertion != "" {
		if len(out) > 0 && out[len(out)-1] != "" {
			out = append(out, "")
		}
		out = append(out, insertion)
	}

	return strings.Join(out, "\n")
}

func mustFindModuleRootAndPath() (string, string) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for {
		data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
		if err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "module ") {
					return dir, strings.TrimSpace(strings.TrimPrefix(line, "module "))
				}
			}
			panic("module line not found in go.mod")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("go.mod not found; run from within a module")
		}
		dir = parent
	}
}
