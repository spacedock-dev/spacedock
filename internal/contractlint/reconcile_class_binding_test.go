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

// contractClassToken matches the step-0 JSON-shape `"class":"a|b|c"` token in the
// FO dispatch contract. The char class tolerates the hyphenated descriptive names
// (un-advanced-pr, stale-branch, local-main-drift) and the `|`-alternation.
var contractClassToken = regexp.MustCompile(`"class":"([A-Za-z|\-]+)"`)

// contractDriftClasses extracts the drift-class vocabulary from the FO dispatch
// contract's step-0 JSON-shape token, splitting its `|`-alternation into members.
// This is the contract-side source — independent of the helper AST scan.
func contractDriftClasses(t *testing.T) []string {
	t.Helper()
	data, err := os.ReadFile(dispatchContractPath(t))
	if err != nil {
		t.Fatalf("read contract: %v", err)
	}
	m := contractClassToken.FindSubmatch(data)
	if m == nil {
		t.Fatalf("contract step-0 JSON-shape `\"class\":\"…\"` token not found in %s", dispatchContractPath(t))
	}
	return strings.Split(string(m[1]), "|")
}

// TestReconcileClassBinding (AC-2) asserts the helper's emitted drift-class set
// (driftClasses var, AST) and the FO dispatch contract step-0 class set (JSON-shape
// token, regex) are the SAME set — two independent sources that red on drift in
// either direction. It is a structural dual-extraction check (an AST literal scan
// and a delimited-token parse), NOT a prose-grep: it never asserts the doc contains
// a given word; it compares two extracted enumerations. The behavior that the
// helper emits these strings is proven by the AC-1 behavioral test that runs the
// helper, not here.
func TestReconcileClassBinding(t *testing.T) {
	helper := helperDriftClasses(t)
	contract := contractDriftClasses(t)

	// Empty-set guard on BOTH sides so the equality cannot pass vacuously (a broken
	// extractor yielding [] on both sides would otherwise "match").
	if len(helper) == 0 {
		t.Fatal("helper-side driftClasses extraction yielded zero classes — extractor bug; the binding would pass vacuously")
	}
	if len(contract) == 0 {
		t.Fatal("contract-side class-token extraction yielded zero classes — extractor bug; the binding would pass vacuously")
	}

	helperSet := toSet(helper)
	contractSet := toSet(contract)
	if !setEqual(helperSet, contractSet) {
		t.Errorf("drift-class set mismatch between the helper and the FO dispatch contract step-0:\n  helper (reconcile.go driftClasses): %v\n  contract (claude-fo-dispatch.md step-0): %v\nneither side may rename, add, or drop a class without the other",
			sortedSet(helperSet), sortedSet(contractSet))
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
