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
	var files []*ast.File
	for _, ent := range entries {
		name := ent.Name()
		if ent.IsDir() || !strings.HasSuffix(name, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, dir+"/"+name, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", name, err)
		}
		files = append(files, f)
	}

	// First pass: discover this package's instruction-file reader helpers, then grow
	// the set to a fixpoint so a tautology cannot hide one hop down. Seeded with the
	// named readers; a func is ALSO a reader if it reads an instruction file directly
	// (a `.md` skill-tree literal in its own body, OR a value built from its own
	// string parameter — the readSkill(t, path) shape, where the caller supplies the
	// `.md` path), OR it WalkDir/Walks a tree collecting `.md` files, OR (transitive)
	// it calls a known reader. The fixpoint closes the multi-hop-helper gap the
	// integration sweep already guards; without it a `wrap(t){return readSkill(t,
	// foCorePath)}` wrapper would leave a tautology behind it undetected.
	readers := map[string]bool{}
	for r := range instructionFileReaders {
		readers[r] = true
	}
	helperCalls := map[string]map[string]bool{}
	for _, f := range files {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || strings.HasPrefix(fn.Name.Name, "Test") {
				continue
			}
			calls, idents, strs := collectCallsAndIdents(fn)
			helperCalls[fn.Name.Name] = calls
			directLiteral := (calls["ReadFile"] || calls["Open"]) && (hasMarkdownLiteral(strs) || readsInstructionIdent(calls, idents))
			if directLiteral || readsParamPath(fn) || walksForMarkdown(calls, strs) {
				readers[fn.Name.Name] = true
			}
		}
	}
	for grew := true; grew; {
		grew = false
		for name, calls := range helperCalls {
			if readers[name] {
				continue
			}
			for r := range readers {
				if calls[r] {
					readers[name] = true
					grew = true
					break
				}
			}
		}
	}

	var offenders []string
	for _, f := range files {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || !strings.HasPrefix(fn.Name.Name, "Test") {
				continue
			}
			calls, idents, strs := collectCallsAndIdents(fn)
			readsInstruction := false
			for r := range readers {
				if calls[r] {
					readsInstruction = true
					break
				}
			}
			// os.ReadFile/os.Open of an instruction-path ident OR an inline `.md`
			// path literal (including a split `name + "." + "md"` reconstructed by
			// collectCallsAndIdents) is also a read of an instruction file.
			if (calls["ReadFile"] || calls["Open"]) && (readsInstructionIdent(calls, idents) || hasMarkdownLiteral(strs)) {
				readsInstruction = true
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

// readsInstructionIdent reports whether a function reads an instruction-path
// package variable (foCorePath, ensignCorePath, …) via os.ReadFile/os.Open.
func readsInstructionIdent(calls, idents map[string]bool) bool {
	if !calls["ReadFile"] && !calls["Open"] {
		return false
	}
	for id := range instructionPathIdents {
		if idents[id] {
			return true
		}
	}
	return false
}

// hasMarkdownLiteral reports whether any collected string (a literal or a
// reconstructed constant concatenation) ends in `.md` — a markdown instruction
// path literal.
func hasMarkdownLiteral(strs map[string]bool) bool {
	for s := range strs {
		if strings.HasSuffix(s, ".md") {
			return true
		}
	}
	return false
}

// readsParamPath reports whether fn os.ReadFile/os.Open's a value derived from one
// of its own string parameters — the path-arg reader shape (readSkill(t, path)
// reads `path`). The `.md` literal lives in the CALLER, so a body-literal scan
// misses it; this catches it by tracking parameter flow into the read argument.
func readsParamPath(fn *ast.FuncDecl) bool {
	params := map[string]bool{}
	if fn.Type.Params != nil {
		for _, field := range fn.Type.Params.List {
			if id, ok := field.Type.(*ast.Ident); ok && id.Name == "string" {
				for _, name := range field.Names {
					params[name.Name] = true
				}
			}
		}
	}
	if len(params) == 0 {
		return false
	}
	found := false
	ast.Inspect(fn, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || (sel.Sel.Name != "ReadFile" && sel.Sel.Name != "Open") {
			return true
		}
		for _, arg := range call.Args {
			ast.Inspect(arg, func(m ast.Node) bool {
				if id, ok := m.(*ast.Ident); ok && params[id.Name] {
					found = true
				}
				return true
			})
		}
		return true
	})
	return found
}

// walksForMarkdown reports whether fn WalkDir/Walks a tree collecting `.md` files —
// the shippedSkillText shape, a reader-of-many whose callers read+match each
// returned path. The signal: a filepath.WalkDir/Walk call AND a `.md` suffix the
// body filters on (carried in strs, which includes split-`.md` concatenations).
func walksForMarkdown(calls, strs map[string]bool) bool {
	if !calls["WalkDir"] && !calls["Walk"] {
		return false
	}
	return hasMarkdownLiteral(strs)
}

// collectCallsAndIdents walks a function body, returning the set of called
// function names (bare + selector trailing name) and the set of referenced
// identifier names. The ident set lets the sweep detect an os.ReadFile applied to a
// package-level instruction-path variable.
//
// The string set ALSO carries the constant value of any string-`+` concatenation
// (e.g. `name + "." + "md"`), so a constructed `.md` suffix cannot evade the
// `.md`-literal detection by splitting the suffix across operands.
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
		case *ast.BinaryExpr:
			if node.Op == token.ADD {
				if joined, ok := constStringConcat(node); ok {
					strs[joined] = true
				}
			}
		}
		return true
	})
	return calls, idents, strs
}

// constStringConcat reconstructs the constant value of a `+` expression whose
// operands are string literals (or nested string-literal `+` expressions),
// treating a non-literal operand as an empty segment so the literal tail still
// reconstructs (`base + "." + "md"` -> `.md`-suffixed). Returns ok=false when no
// string literal participates.
func constStringConcat(expr ast.Expr) (string, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			return strings.Trim(e.Value, "`\""), true
		}
		return "", false
	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			return "", false
		}
		l, lok := constStringConcat(e.X)
		r, rok := constStringConcat(e.Y)
		if !lok && !rok {
			return "", false
		}
		return l + r, true
	default:
		return "", false
	}
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

	// A multi-hop-helper tautology — the read hidden one frame down behind a wrapper
	// that calls the named reader — must also be flagged: the transitive reader
	// fixpoint propagates reader-ness up the call chain. Before the fixpoint the HN
	// sweep left this GREEN (the validation Cycle-1 finding 1).
	multiHop := `package fixture
import "strings"
func wrapHop(t *T) string { return readSkill(t, foCorePath) }
func TestMultiHopUndeclaredHN(t *T) {
	text := wrapHop(t)
	if strings.Contains(text, "x") { _ = text }
}
`
	dir2 := t.TempDir()
	writeFixture(t, dir2+"/multihop_test.go", multiHop)
	off = sweepHostneutralityTautologies(t, dir2)
	if !containsStrHN(off, "TestMultiHopUndeclaredHN") {
		t.Fatalf("sweep failed to flag a multi-hop-helper tautology (transitive reader fixpoint not working); offenders=%v", off)
	}
}

// TestHostneutralitySweepDetectsEvasionShapes is the planted-control mutation test
// for the reader-discovery evasion shapes the validation audit proved the HN sweep
// missed (the integration sweep already guarded these; this ports the guard). Each
// case plants a synthetic offender reaching an instruction file through a shape the
// naive named-reader/`.md`-literal detection cannot see, runs the sweep, asserts it
// REDs, then plants the declared form and asserts it GREENs. A regression removing a
// discovery mechanism leaves the matching case un-flagged, failing this control.
func TestHostneutralitySweepDetectsEvasionShapes(t *testing.T) {
	// Shape 1 — multi-hop transitive helper: the tautology hides one hop down behind
	// a wrapper that calls the named reader readSkill. The fixpoint must propagate
	// reader-ness up the chain. This is the finding-1 control: before the fixpoint,
	// the HN sweep left this GREEN.
	multiHop := `package fixture
import "strings"
func wrapHop(t *T) string {
	return readSkill(t, foCorePath)
}
func TestMultiHopHN(t *T) {
	text := wrapHop(t)
	if strings.Contains(text, "x") { _ = text }
}
`
	assertRedThenGreenHN(t, "multi-hop transitive helper", "TestMultiHopHN", multiHop)

	// Shape 2 — path-arg reader: a NEW helper (not in the named map) os.ReadFile's a
	// value built from its own path parameter; the `.md` literal lives in the caller.
	pathArg := `package fixture
import (
	"os"
	"strings"
)
func readArg(t *T, path string) string {
	b, _ := os.ReadFile(path)
	return string(b)
}
func TestPathArgHN(t *T) {
	text := readArg(t, "../../skills/first-officer/references/first-officer-shared-core.md")
	if strings.Contains(text, "x") { _ = text }
}
`
	assertRedThenGreenHN(t, "path-arg reader", "TestPathArgHN", pathArg)

	// Shape 3 — WalkDir collector: a helper WalkDirs a tree returning `.md` paths the
	// caller reads+matches.
	walkDir := `package fixture
import (
	"os"
	"path/filepath"
	"strings"
)
func walkSkills(t *T, base string) []string {
	var out []string
	filepath.WalkDir(base, func(p string, d os.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(p, ".md") { out = append(out, p) }
		return nil
	})
	return out
}
func TestWalkDirHN(t *T) {
	for _, p := range walkSkills(t, "../../skills") {
		b, _ := os.ReadFile(p)
		if strings.Contains(string(b), "x") { _ = b }
	}
}
`
	assertRedThenGreenHN(t, "WalkDir collector", "TestWalkDirHN", walkDir)

	// Shape 4 — split-".md" suffix: the read path is built as base + "." + "md", so
	// no single literal carries the `.md` suffix; constStringConcat must rejoin it.
	splitMD := `package fixture
import (
	"os"
	"path/filepath"
	"strings"
)
func TestSplitSuffixHN(t *T) {
	p := filepath.Join("..", "..", "skills", "first-officer", "references", "first-officer-shared-core" + "." + "md")
	b, _ := os.ReadFile(p)
	if strings.Contains(string(b), "x") { _ = b }
}
`
	assertRedThenGreenHN(t, "split-.md suffix", "TestSplitSuffixHN", splitMD)
}

// assertRedThenGreenHN plants fixtureSrc, runs the HN sweep, requires offenderName
// flagged (RED on the evasion shape), then rewrites the test with a markNonAC
// declaration and requires it cleared (GREEN once declared).
func assertRedThenGreenHN(t *testing.T, shape, offenderName, fixtureSrc string) {
	t.Helper()
	dir := t.TempDir()
	writeFixture(t, dir+"/evasion_test.go", fixtureSrc)
	off := sweepHostneutralityTautologies(t, dir)
	if !containsStrHN(off, offenderName) {
		t.Fatalf("%s: HN sweep failed to flag the evasion offender %s; offenders=%v", shape, offenderName, off)
	}

	declaredSrc := strings.Replace(
		fixtureSrc,
		"func "+offenderName+"(t *T) {",
		"func "+offenderName+`(t *T) {
	markNonAC(t, "declared evasion-shape fixture")`,
		1,
	)
	if declaredSrc == fixtureSrc {
		t.Fatalf("%s: could not inject markNonAC into the fixture for %s (signature not found)", shape, offenderName)
	}
	dir2 := t.TempDir()
	writeFixture(t, dir2+"/evasion_test.go", declaredSrc)
	off = sweepHostneutralityTautologies(t, dir2)
	if containsStrHN(off, offenderName) {
		t.Fatalf("%s: HN sweep still flagged %s after it declared markNonAC; offenders=%v", shape, offenderName, off)
	}
}

func writeFixture(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
