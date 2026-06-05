// ABOUTME: The non-AC text-consistency marker + the AC-3 sweep meta-test for the
// ABOUTME: hostneutrality suite — a presence/absence check over an instruction file proves nothing unless it self-classifies.
package hostneutrality

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"testing"
)

// markNonAC declares a test a non-AC text-consistency lint (the prose carries a
// required clause / stays free of a banned token), NOT a behavioral proof, naming
// the behavioral oracle (a live drive, a code-side test, or "n/a — the claim is
// about the text"). The proof policy (f8b257cf) bans a string match over an
// instruction file the model reads as proof of any behavioral acceptance
// criterion. The AC-3 sweep (TestNoUndeclaredHostneutralityTautology) keys on this
// call; it does nothing at runtime.
func markNonAC(t *testing.T, oracle string) {
	t.Helper()
	if oracle == "" {
		t.Fatal("markNonAC requires a non-empty behavioral oracle reference")
	}
}

// markCodeBoundInvariant declares a test's expectation comes from an independent
// code-side source (a Go const, an env-var token the binary defines, a dispatch
// subcommand, a DIFFERENT file's n-gram) that can DIVERGE from the file under
// test — a legitimate invariant, not a tautology. The sweep treats it as declared.
func markCodeBoundInvariant(t *testing.T, source string) {
	t.Helper()
	if source == "" {
		t.Fatal("markCodeBoundInvariant requires a non-empty independent-source reference")
	}
}

// instructionFileReaders are the helpers that read a markdown instruction file the
// model ingests (a skill or contract). A test that calls one AND a substring/regex
// match is a presence/absence check over an ingested file. Tests that scan CODE
// (host_neutrality_test.go's scanFile over .go files) are NOT in this set, so the
// sweep does not flag the legitimate go/parser code invariants.
var instructionFileReaders = map[string]bool{
	"readSkill": true,
	"readText":  true,
	// The markdown-span parsers read an instruction file internally; a test that
	// drives one is reading an instruction file even though the os.ReadFile lives
	// one frame down.
	"parseSpans":                true,
	"parseProseSpansForOverlap": true,
}

// instructionPathIdents are the package-level path variables that resolve to a
// markdown instruction file. A test that reads one of these via os.ReadFile/os.Open
// and matches it is a presence/absence check over an ingested file. (Code-scanning
// tests reference ../dispatch, ../status package dirs, never these.)
var instructionPathIdents = map[string]bool{
	"foCorePath":               true,
	"ensignCorePath":           true,
	"commissionSkillPath":      true,
	"sharedCorePath":           true,
	"claudeRuntimePath":        true,
	"contractProseFiles":       true,
	"sharedCorePaths":          true,
	"runtimeAdapterPaths":      true,
	"devLeakageCorePaths":      true,
	"runtimeAdapterFieldPaths": true,
	"devHomePresence":          true,
}

var matchFuncs = map[string]bool{
	"Contains":              true,
	"Count":                 true,
	"HasPrefix":             true,
	"HasSuffix":             true,
	"FindString":            true,
	"FindStringIndex":       true,
	"FindStringSubmatch":    true,
	"MatchString":           true,
	"FindAllStringSubmatch": true,
	// Package match-helpers: a test that drives one of these over an instruction
	// file is still a presence/restatement check even though the strings.Contains
	// lives one frame down. assertAll matches a required-token list; the span
	// parsers set up a per-span token scan or a cross-file n-gram restatement check.
	"assertAll":                 true,
	"parseSpans":                true,
	"parseProseSpansForOverlap": true,
}

// TestNoUndeclaredHostneutralityTautology is the AC-3 sweep for this package,
// re-runnable offline. It parses every *_test.go and flags any test function that
// reads a markdown INSTRUCTION file (an instructionFileReader call, or os.ReadFile/
// os.Open of an instructionPathIdent) AND matches it with a substring/regex call,
// UNLESS the test self-classifies via markNonAC or markCodeBoundInvariant. The
// go/parser code-scan invariants (host_neutrality_test.go over .go source) and the
// spanHostQualified unit test are NOT flagged: they read no instruction file or
// match no ingested text. The undeclared-offender count is the AC-3 metric; it
// must be zero.
func TestNoUndeclaredHostneutralityTautology(t *testing.T) {
	offenders := sweepHostneutralityTautologies(t, ".")
	for _, o := range offenders {
		t.Errorf("%s reads a markdown instruction file and matches it as a standalone claim without self-classifying — call markNonAC (with its behavioral oracle) or markCodeBoundInvariant (with its independent source)", o)
	}
	if len(offenders) > 0 {
		t.Fatalf("AC-3 sweep: %d undeclared tautological-behavioral-proof test(s) in hostneutrality; the count must be zero", len(offenders))
	}
}

func sweepHostneutralityTautologies(t *testing.T, dir string) []string {
	t.Helper()
	fset := token.NewFileSet()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read package dir %s: %v", dir, err)
	}
	var offenders []string
	for _, ent := range entries {
		name := ent.Name()
		if ent.IsDir() || !strings.HasSuffix(name, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, dir+"/"+name, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", name, err)
		}
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || !strings.HasPrefix(fn.Name.Name, "Test") {
				continue
			}
			calls, idents, strs := collectCallsAndIdents(fn)
			readsInstruction := false
			for r := range instructionFileReaders {
				if calls[r] {
					readsInstruction = true
				}
			}
			// os.ReadFile/os.Open of an instruction-path ident OR an inline `.md`
			// path literal is also a read of an instruction file.
			if calls["ReadFile"] || calls["Open"] {
				for id := range instructionPathIdents {
					if idents[id] {
						readsInstruction = true
					}
				}
				for s := range strs {
					if strings.HasSuffix(s, ".md") {
						readsInstruction = true
					}
				}
			}
			matches := false
			for m := range matchFuncs {
				if calls[m] {
					matches = true
				}
			}
			declared := calls["markNonAC"] || calls["markCodeBoundInvariant"]
			if readsInstruction && matches && !declared {
				offenders = append(offenders, fn.Name.Name)
			}
		}
	}
	return sortedUniqueHN(offenders)
}

// collectCallsAndIdents walks a function body, returning the set of called
// function names (bare + selector trailing name) and the set of referenced
// identifier names. The ident set lets the sweep detect an os.ReadFile applied to a
// package-level instruction-path variable.
func collectCallsAndIdents(fn *ast.FuncDecl) (calls map[string]bool, idents map[string]bool, strs map[string]bool) {
	calls = map[string]bool{}
	idents = map[string]bool{}
	strs = map[string]bool{}
	ast.Inspect(fn, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.CallExpr:
			switch f := node.Fun.(type) {
			case *ast.Ident:
				calls[f.Name] = true
			case *ast.SelectorExpr:
				calls[f.Sel.Name] = true
			}
		case *ast.Ident:
			idents[node.Name] = true
		case *ast.BasicLit:
			if node.Kind == token.STRING {
				strs[strings.Trim(node.Value, "`\"")] = true
			}
		}
		return true
	})
	return calls, idents, strs
}

func sortedUniqueHN(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j-1] > out[j]; j-- {
			out[j-1], out[j] = out[j], out[j-1]
		}
	}
	return out
}

func containsStrHN(in []string, want string) bool {
	for _, s := range in {
		if s == want {
			return true
		}
	}
	return false
}

// TestHostneutralitySweepDetectsAnUndeclaredTautology is the mutation control for
// the sweep: it must RED on the shape it polices and GREEN once that shape
// self-classifies. Writes synthetic fixtures to a temp dir and runs the sweep.
func TestHostneutralitySweepDetectsAnUndeclaredTautology(t *testing.T) {
	dir := t.TempDir()
	undeclared := `package fixture
import "strings"
func TestUndeclaredHN(t *T) {
	text := readSkill(t, foCorePath)
	if strings.Contains(text, "x") { _ = text }
}
`
	declared := `package fixture
func TestDeclaredHN(t *T) {
	markNonAC(t, "live split-root-halt scenario")
	text := readSkill(t, foCorePath)
	if strings.Contains(text, "x") { _ = text }
}
`
	codeScan := `package fixture
func TestCodeScanHN(t *T) {
	leaks := scanFile(t, "../dispatch/x.go")
	if strings.Contains(leaks[0].text, ".claude") { _ = leaks }
}
`
	writeFixture(t, dir+"/undeclared_test.go", undeclared)
	off := sweepHostneutralityTautologies(t, dir)
	if !containsStrHN(off, "TestUndeclaredHN") {
		t.Fatalf("sweep failed to flag an undeclared instruction-file presence check; offenders=%v", off)
	}

	writeFixture(t, dir+"/declared_test.go", declared)
	writeFixture(t, dir+"/codescan_test.go", codeScan)
	off = sweepHostneutralityTautologies(t, dir)
	if containsStrHN(off, "TestDeclaredHN") {
		t.Fatalf("sweep wrongly flagged a declared lint; offenders=%v", off)
	}
	if containsStrHN(off, "TestCodeScanHN") {
		t.Fatalf("sweep wrongly flagged a code-scanning invariant (reads no instruction file); offenders=%v", off)
	}
	if !containsStrHN(off, "TestUndeclaredHN") {
		t.Fatalf("adding declared/codescan fixtures must not stop the sweep flagging the undeclared one; offenders=%v", off)
	}
}

func writeFixture(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
