// ABOUTME: Code-derived spacedock CLI vocabulary — the independent source the
// ABOUTME: leakage checks bind to, AST-extracted from the dispatch router + the status stage-option keys.
package integration

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"testing"
)

// skillSeamRe captures the skill name from a `Skill(skill="spacedock:NAME")`
// invocation in an FO/ensign contract — the integration seam the FO actually
// invokes mid-run. The captured NAME is the independent source a skill's
// frontmatter `name:` binds to: if the contract's invocation and the skill's
// declared name drift apart, the seam breaks, and a check comparing the two REDs.
var skillSeamRe = regexp.MustCompile(`Skill\(skill="spacedock:([a-z0-9-]+)"\)`)

// invokedSeamName returns the skill name the given contract text invokes for the
// expected target. It scans every Skill(skill="spacedock:NAME") invocation and
// returns the matching NAME, or "" if the contract never invokes that seam. The
// expected value is read from the CONTRACT (the file that drives the FO), not from
// the skill file under test, so the two have independent sources that can diverge.
func invokedSeamName(contract, want string) string {
	for _, m := range skillSeamRe.FindAllStringSubmatch(contract, -1) {
		if m[1] == want {
			return m[1]
		}
	}
	return ""
}

// spacedockDispatchSubcommands AST-extracts the dispatch subcommand names the
// binary actually routes, from the `switch args[0] { case "..." }` in
// internal/dispatch/dispatch.go's Run. This is an independent code-side source:
// the binary parses these as commands, not from any instruction file the model
// reads, and a rename in the router shifts the set — which is exactly what lets a
// leakage check that binds to it diverge from a stale token frozen in a skill.
func spacedockDispatchSubcommands(t *testing.T) []string {
	t.Helper()
	src := filepath.Join(repoRoot(t), "internal", "dispatch", "dispatch.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, src, nil, 0)
	if err != nil {
		t.Fatalf("parse dispatch.go: %v", err)
	}
	var subs []string
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
			// Only the subcommand switch (`switch args[0]`) yields command names.
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
				for _, expr := range cc.List {
					if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
						subs = append(subs, trimQuotes(lit.Value))
					}
				}
			}
			return false
		})
		return false
	})
	if len(subs) == 0 {
		t.Fatal("extracted zero dispatch subcommands from dispatch.go — the AST source diverged from the router shape")
	}
	return subs
}

// spacedockBuildRequestFlags AST-extracts the dispatch-build request flag names
// the binary accepts, from the `case "--entity-path", ...` in dispatch.go's
// isBuildRequestFlag. These are the real flags the build path parses — an
// independent code-side source for the "docs teach file-backed dispatch input"
// invariant: a flag renamed in code shifts the set, so a docs check binding to it
// tracks the binary's actual flag surface rather than a frozen literal.
func spacedockBuildRequestFlags(t *testing.T) []string {
	t.Helper()
	src := filepath.Join(repoRoot(t), "internal", "dispatch", "dispatch.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, src, nil, 0)
	if err != nil {
		t.Fatalf("parse dispatch.go: %v", err)
	}
	var flags []string
	ast.Inspect(f, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "isBuildRequestFlag" {
			return true
		}
		ast.Inspect(fn, func(m ast.Node) bool {
			cc, ok := m.(*ast.CaseClause)
			if !ok {
				return true
			}
			for _, e := range cc.List {
				if lit, ok := e.(*ast.BasicLit); ok && lit.Kind == token.STRING {
					flags = append(flags, trimQuotes(lit.Value))
				}
			}
			return true
		})
		return false
	})
	if len(flags) == 0 {
		t.Fatal("extracted zero build-request flags from isBuildRequestFlag in dispatch.go")
	}
	return flags
}

// spacedockTopLevelCommands AST-extracts the binary's top-level command names
// from the `Use: "<name> ..."` fields of the cobra commands in internal/cli/cli.go.
// The first word of each Use string is the command verb (`status`, `dispatch`,
// `claude`, ...). This is the independent source for the `spacedock dispatch` /
// `spacedock status` leakage prefixes: the binary registers these verbs, and a
// rename in cli.go shifts the set, so a leakage check that binds to it tracks the
// real command surface rather than a frozen literal.
func spacedockTopLevelCommands(t *testing.T) []string {
	t.Helper()
	src := filepath.Join(repoRoot(t), "internal", "cli", "cli.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, src, nil, 0)
	if err != nil {
		t.Fatalf("parse cli.go: %v", err)
	}
	var cmds []string
	ast.Inspect(f, func(n ast.Node) bool {
		kv, ok := n.(*ast.KeyValueExpr)
		if !ok {
			return true
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok || key.Name != "Use" {
			return true
		}
		lit, ok := kv.Value.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}
		use := trimQuotes(lit.Value)
		if i := indexSpace(use); i >= 0 {
			use = use[:i]
		}
		if use != "" && use != "spacedock" {
			cmds = append(cmds, use)
		}
		return true
	})
	if len(cmds) == 0 {
		t.Fatal("extracted zero top-level commands from cli.go Use: fields")
	}
	return cmds
}

func indexSpace(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' {
			return i
		}
	}
	return -1
}

// spacedockStageOptionKeys AST-extracts the stage-option keys the binary parses,
// from the `[]string{"feedback-to", ...}` literal in internal/status/stages.go.
// These are real frontmatter keys the binary reads, not file prose — an
// independent source that shifts when a key is renamed in code.
func spacedockStageOptionKeys(t *testing.T) []string {
	t.Helper()
	src := filepath.Join(repoRoot(t), "internal", "status", "stages.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, src, nil, 0)
	if err != nil {
		t.Fatalf("parse stages.go: %v", err)
	}
	var keys []string
	ast.Inspect(f, func(n ast.Node) bool {
		cl, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}
		arr, ok := cl.Type.(*ast.ArrayType)
		if !ok {
			return true
		}
		if id, ok := arr.Elt.(*ast.Ident); !ok || id.Name != "string" {
			return true
		}
		var lits []string
		for _, e := range cl.Elts {
			lit, ok := e.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			lits = append(lits, trimQuotes(lit.Value))
		}
		// The stage-option list is the one containing "feedback-to" — pin to that
		// literal so an unrelated []string{...} does not contribute.
		for _, l := range lits {
			if l == "feedback-to" {
				keys = lits
				return false
			}
		}
		return true
	})
	if len(keys) == 0 {
		t.Fatal("did not find the stage-option-keys []string{...} (feedback-to) in stages.go")
	}
	return keys
}

// spacedockLeakageTokens is the canonical spacedock-specific token set the
// host-neutral / generic-skill bodies must NOT name. It is DERIVED from code (the
// dispatch router's subcommand surface + the status stage-option keys), not a
// literal frozen in a test, so the set and any skill body can diverge — that
// divergence is what makes a leakage check able to fail as an invariant rather
// than as a self-match.
//
// Only spacedock-SPECIFIC forms are included, so a generic-prose word never
// false-fires: the bare `spacedock dispatch` / `spacedock status` prefixes; the
// dispatch helper subcommands QUALIFIED with their `spacedock dispatch ` prefix
// (so a bare English `reconcile`/`build` in event-loop prose is fine, only the
// qualified invocation leaks); and the hyphenated stage-option keys (e.g.
// `feedback-to`), which are spacedock frontmatter vocabulary, not generic words —
// the single-word stage keys (agent/fresh/model) are excluded because they appear
// in legitimate generic prose.
func spacedockLeakageTokens(t *testing.T) []string {
	t.Helper()
	var tokens []string
	// The leak-prone top-level prefixes: dispatch + status, derived from the
	// binary's registered command verbs (not a frozen literal).
	for _, cmd := range spacedockTopLevelCommands(t) {
		if cmd == "dispatch" || cmd == "status" {
			tokens = append(tokens, "spacedock "+cmd)
		}
	}
	for _, sub := range spacedockDispatchSubcommands(t) {
		// build/show-stage-def are the host-neutral surface a skill may name; the
		// Claude-coupled helper subcommands are the leak-prone ones — banned only in
		// their qualified `spacedock dispatch <sub>` form so a generic English word
		// (a bare `reconcile`) does not false-fire.
		switch sub {
		case "build", "show-stage-def":
		default:
			tokens = append(tokens, "spacedock dispatch "+sub)
		}
	}
	for _, k := range spacedockStageOptionKeys(t) {
		// Only the hyphenated, spacedock-specific stage keys (feedback-to) are
		// leak-prone as a bare token; single-word keys appear in generic prose.
		if isHyphenated(k) {
			tokens = append(tokens, k)
		}
	}
	return tokens
}

func isHyphenated(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '-' {
			return true
		}
	}
	return false
}

// intersect returns the want values that appear in have — used to derive the
// load-bearing subset of a code-extracted set without hardcoding the full set.
func intersect(have []string, want ...string) []string {
	present := map[string]bool{}
	for _, h := range have {
		present[h] = true
	}
	var out []string
	for _, w := range want {
		if present[w] {
			out = append(out, w)
		}
	}
	return out
}

func trimQuotes(s string) string {
	if len(s) >= 2 && (s[0] == '"' || s[0] == '`') {
		return s[1 : len(s)-1]
	}
	return s
}
