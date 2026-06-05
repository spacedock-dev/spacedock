// ABOUTME: Code-derived sources the hostneutrality re-binds bind to — host env-var
// ABOUTME: names and dispatch subcommands AST-extracted from the binary, so a check's expectation has an independent source.
package hostneutrality

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"
)

// repoRoot is the project root (two levels up from this package's source dir).
func repoRoot() string {
	return filepath.Join("..", "..")
}

// hostEnvVar AST-extracts the host-derivation env-var name the binary reads from
// internal/dispatch/build.go (the `getenv("CODEX_THREAD_ID")` / "CLAUDECODE"
// selectors). It returns the name if the binary reads it, else "". This is the
// independent source for "the skill branches on the same env var the binary reads":
// if the binary stops reading the var, or the skill stops branching on it, the two
// diverge and a check binding to this reds.
func hostEnvVar(t *testing.T, name string) string {
	t.Helper()
	src := filepath.Join(repoRoot(), "internal", "dispatch", "build.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, src, nil, 0)
	if err != nil {
		t.Fatalf("parse build.go: %v", err)
	}
	found := ""
	ast.Inspect(f, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}
		if trimLit(lit.Value) == name {
			found = name
			return false
		}
		return true
	})
	return found
}

// dispatchSubcommands AST-extracts the dispatch subcommand names the binary routes
// from internal/dispatch/dispatch.go's Run switch — the independent source for the
// claude-helper / relocated-command checks (so a renamed subcommand shifts the set
// rather than the test self-matching a frozen literal).
func dispatchSubcommands(t *testing.T) map[string]bool {
	t.Helper()
	src := filepath.Join(repoRoot(), "internal", "dispatch", "dispatch.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, src, nil, 0)
	if err != nil {
		t.Fatalf("parse dispatch.go: %v", err)
	}
	subs := map[string]bool{}
	ast.Inspect(f, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "Run" {
			return true
		}
		ast.Inspect(fn, func(m ast.Node) bool {
			sw, ok := m.(*ast.SwitchStmt)
			if !ok {
				return true
			}
			tag, ok := sw.Tag.(*ast.IndexExpr)
			if !ok {
				return true
			}
			if id, ok := tag.X.(*ast.Ident); !ok || id.Name != "args" {
				return true
			}
			for _, stmt := range sw.Body.List {
				cc, ok := stmt.(*ast.CaseClause)
				if !ok {
					continue
				}
				for _, e := range cc.List {
					if lit, ok := e.(*ast.BasicLit); ok && lit.Kind == token.STRING {
						subs[trimLit(lit.Value)] = true
					}
				}
			}
			return false
		})
		return false
	})
	if len(subs) == 0 {
		t.Fatal("extracted zero dispatch subcommands from dispatch.go")
	}
	return subs
}

func trimLit(s string) string {
	if len(s) >= 2 && (s[0] == '"' || s[0] == '`') {
		return s[1 : len(s)-1]
	}
	return s
}
