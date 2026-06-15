// ABOUTME: AC-2 — binds the reconcile drift-class set in two independent sources
// ABOUTME: (helper driftClasses var via AST, contract step-0 JSON-shape token via regex) as equal sets.
package contractlint

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// reconcileSourcePath is the helper whose driftClasses var declares the emitted
// drift-class vocabulary.
func reconcileSourcePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(repoRoot(t), "internal", "dispatch", "reconcile.go")
}

// dispatchContractPath is the FO dispatch contract whose event-loop step-0 names
// the same class set in its JSON-shape token.
func dispatchContractPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(repoRoot(t), "skills", "first-officer", "references", "claude-fo-dispatch.md")
}

// helperDriftClasses extracts the emitted drift-class vocabulary from reconcile.go
// by AST: it reads the `driftClasses` var's composite-literal elements (which are
// const identifiers) and resolves each to its declared string value via the file's
// const declarations. This binds the var — the single helper-side source the
// rename introduced — to the actual emitted strings, not to a re-typed literal.
func helperDriftClasses(t *testing.T) []string {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, reconcileSourcePath(t), nil, 0)
	if err != nil {
		t.Fatalf("parse reconcile.go: %v", err)
	}

	// 1. const name -> string value, for every `name = "literal"` const spec.
	constVals := map[string]string{}
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.CONST {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok || len(vs.Names) != len(vs.Values) {
				continue
			}
			for i, name := range vs.Names {
				if lit, ok := vs.Values[i].(*ast.BasicLit); ok && lit.Kind == token.STRING {
					if v, err := strconv.Unquote(lit.Value); err == nil {
						constVals[name.Name] = v
					}
				}
			}
		}
	}

	// 2. driftClasses var's composite-literal elements -> resolved const values.
	var classes []string
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, name := range vs.Names {
				if name.Name != "driftClasses" || i >= len(vs.Values) {
					continue
				}
				lit, ok := vs.Values[i].(*ast.CompositeLit)
				if !ok {
					t.Fatalf("driftClasses is not a composite literal: %T", vs.Values[i])
				}
				for _, el := range lit.Elts {
					ident, ok := el.(*ast.Ident)
					if !ok {
						t.Fatalf("driftClasses element is not an identifier: %T", el)
					}
					v, ok := constVals[ident.Name]
					if !ok {
						t.Fatalf("driftClasses references %q with no resolvable const string value", ident.Name)
					}
					classes = append(classes, v)
				}
			}
		}
	}
	return classes
}

// The step-0 block bounds. The contract event-loop step-0 opens with the
// `0. **Reconcile sweep.**` heading and runs until the next numbered step
// (`1. **`). All three class-bearing surfaces (JSON-shape token, action bullets,
// one-line summary) live inside this slice; bounding to it keeps the extractors
// from picking up an identically-shaped token elsewhere in the contract.
var (
	step0HeadingRe = regexp.MustCompile(`(?m)^0\. \*\*Reconcile sweep\.\*\*`)
	nextStepRe     = regexp.MustCompile(`(?m)^1\. \*\*`)
)

// contractClassToken matches the step-0 JSON-shape `"class":"a|b|c"` token. The
// char class tolerates the hyphenated descriptive names (un-advanced-pr,
// stale-branch, local-main-drift) and the `|`-alternation.
var contractClassToken = regexp.MustCompile(`"class":"([A-Za-z|\-]+)"`)

// actionBulletRe matches one step-0 per-class action bullet head: a three-space
// indented `- ` marker, then the bolded class name(s) up to the `→` action arrow.
// The combined `- **lingering** / **superseded** →` bullet carries two names
// before the arrow; the captured prefix is re-scanned for every `**name**` span.
var actionBulletRe = regexp.MustCompile(`(?m)^   - (.*?)→`)

// boldNameRe extracts a `**name**` span's inner token (a descriptive class name).
var boldNameRe = regexp.MustCompile(`\*\*([a-z][a-z-]*)\*\*`)

// summaryClassRe matches a `name={N}` pair in the step-0 one-line drift summary
// (`reconcile: {N} entries: lingering={N} superseded={N} … — acting`). The `={N}`
// suffix anchors it to the summary template, not arbitrary prose.
var summaryClassRe = regexp.MustCompile(`([a-z][a-z-]*)=\{N\}`)

// step0Block returns the contract's event-loop step-0 slice, bounded by its
// heading and the next numbered step. Fails if either anchor is missing.
func step0Block(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(dispatchContractPath(t))
	if err != nil {
		t.Fatalf("read contract: %v", err)
	}
	loc := step0HeadingRe.FindIndex(data)
	if loc == nil {
		t.Fatalf("step-0 heading `0. **Reconcile sweep.**` not found in %s", dispatchContractPath(t))
	}
	rest := data[loc[0]:]
	if end := nextStepRe.FindIndex(rest); end != nil {
		rest = rest[:end[0]]
	}
	return string(rest)
}

// contractClassesFromToken extracts the class set from step-0's JSON-shape token.
func contractClassesFromToken(t *testing.T, block string) []string {
	t.Helper()
	m := contractClassToken.FindStringSubmatch(block)
	if m == nil {
		t.Fatalf("step-0 JSON-shape `\"class\":\"…\"` token not found")
	}
	return strings.Split(m[1], "|")
}

// contractClassesFromBullets extracts the class set from step-0's five per-class
// action bullets (the bolded name(s) before each `→`).
func contractClassesFromBullets(t *testing.T, block string) []string {
	t.Helper()
	var out []string
	for _, bullet := range actionBulletRe.FindAllStringSubmatch(block, -1) {
		for _, name := range boldNameRe.FindAllStringSubmatch(bullet[1], -1) {
			out = append(out, name[1])
		}
	}
	if len(out) == 0 {
		t.Fatalf("step-0 per-class action bullets yielded no class names")
	}
	return out
}

// contractClassesFromSummary extracts the class set from step-0's one-line drift
// summary (`name={N}` pairs).
func contractClassesFromSummary(t *testing.T, block string) []string {
	t.Helper()
	var out []string
	for _, m := range summaryClassRe.FindAllStringSubmatch(block, -1) {
		out = append(out, m[1])
	}
	if len(out) == 0 {
		t.Fatalf("step-0 one-line drift summary yielded no `name={N}` class tokens")
	}
	return out
}

// TestReconcileClassBinding (AC-2) asserts the helper's emitted drift-class set
// (driftClasses var, AST) and EACH of the FO dispatch contract's three step-0
// class-bearing surfaces — the JSON-shape token, the five per-class action
// bullets, and the one-line drift summary — are the SAME set. They are independent
// sources that red on drift in any direction: a class renamed, added, or dropped in
// the helper OR in any single step-0 surface reds the binding. It is a structural
// dual-extraction check (an AST literal scan and three delimited-token parses), NOT
// a prose-grep: it never asserts the doc contains a given word; it compares
// extracted enumerations. The behavior that the helper emits these strings is
// proven by the AC-1 behavioral test that runs the helper, not here.
func TestReconcileClassBinding(t *testing.T) {
	helper := helperDriftClasses(t)
	// Empty-set guard so the equality cannot pass vacuously (a broken extractor
	// yielding [] on both sides would otherwise "match").
	if len(helper) == 0 {
		t.Fatal("helper-side driftClasses extraction yielded zero classes — extractor bug; the binding would pass vacuously")
	}
	helperSet := toSet(helper)

	block := step0Block(t)
	surfaces := []struct {
		name    string
		classes []string
	}{
		{`JSON-shape token ("class":"…")`, contractClassesFromToken(t, block)},
		{"per-class action bullets", contractClassesFromBullets(t, block)},
		{"one-line drift summary", contractClassesFromSummary(t, block)},
	}
	for _, s := range surfaces {
		if len(s.classes) == 0 {
			t.Fatalf("contract step-0 %s extraction yielded zero classes — extractor bug; the binding would pass vacuously", s.name)
		}
		surfaceSet := toSet(s.classes)
		if !setEqual(helperSet, surfaceSet) {
			t.Errorf("drift-class set mismatch between the helper and the FO dispatch contract step-0 %s:\n  helper (reconcile.go driftClasses): %v\n  contract (claude-fo-dispatch.md step-0 %s): %v\nneither side may rename, add, or drop a class without the other",
				s.name, sortedSet(helperSet), s.name, sortedSet(surfaceSet))
		}
	}
}

// toSet builds a set from a slice, surfacing accidental duplicates implicitly
// (a duplicate collapses, so a set-equality mismatch catches a doubled member).
func toSet(xs []string) map[string]bool {
	out := map[string]bool{}
	for _, x := range xs {
		out[x] = true
	}
	return out
}

// setEqual reports whether two string sets hold the identical members.
func setEqual(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}

// sortedSet returns a set's members sorted, for deterministic error output.
func sortedSet(s map[string]bool) []string {
	out := make([]string, 0, len(s))
	for k := range s {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
