package routes

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// TestURLCalls reads the module's source (generated *_templ.go too) and checks that every
// routes.URL(routes.X, ...) passes as many args as X has params, which the compiler cannot see;
// a call whose pattern is a variable is not checked
func TestURLCalls(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "routes.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	consts := map[string]ast.Expr{}
	for _, d := range f.Decls {
		if g, ok := d.(*ast.GenDecl); ok && g.Tok == token.CONST {
			for _, s := range g.Specs {
				v := s.(*ast.ValueSpec)
				for i, n := range v.Names {
					consts[n.Name] = v.Values[i]
				}
			}
		}
	}
	var eval func(e ast.Expr) string
	eval = func(e ast.Expr) string {
		switch x := e.(type) {
		case *ast.BasicLit:
			s, err := strconv.Unquote(x.Value)
			if err != nil {
				t.Fatal(err)
			}
			return s
		case *ast.Ident:
			return eval(consts[x.Name])
		case *ast.BinaryExpr:
			return eval(x.X) + eval(x.Y)
		}
		t.Fatalf("routes.go: cannot evaluate %T", e)
		return ""
	}

	calls := 0
	err = filepath.WalkDir("../..", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}
		ast.Inspect(file, func(n ast.Node) bool {
			c, ok := n.(*ast.CallExpr)
			if !ok || !isRoutesSel(c.Fun, "URL") || len(c.Args) == 0 || c.Ellipsis.IsValid() {
				return true
			}
			pattern, ok := c.Args[0].(*ast.SelectorExpr)
			if !ok || !isRoutesSel(pattern, pattern.Sel.Name) {
				return true
			}
			calls++
			want := strings.Count(eval(consts[pattern.Sel.Name]), "/:")
			if got := len(c.Args) - 1; got != want {
				t.Errorf("%s: routes.URL(routes.%s, ...) has %d args, the pattern has %d params", fset.Position(c.Pos()), pattern.Sel.Name, got, want)
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if calls == 0 {
		t.Fatal("found no routes.URL calls, is the walk in the module root?")
	}
}

func isRoutesSel(e ast.Expr, name string) bool {
	s, ok := e.(*ast.SelectorExpr)
	if !ok || s.Sel.Name != name {
		return false
	}
	x, ok := s.X.(*ast.Ident)
	return ok && x.Name == "routes"
}
