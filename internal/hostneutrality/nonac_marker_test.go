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

// instructionFileReaders are the named helpers that read a markdown instruction file
// the model ingests (a skill or contract) — the allowlist of recognized readers. A
// test that calls one is reading an ingested file (the READ alone triggers the
// must-declare rule; how it then inspects the bytes is irrelevant). Tests that scan
// CODE (host_neutrality_test.go's scanFile over .go files via parser.ParseFile, not a
// content read sink) are NOT in this set, so the sweep does not flag the legitimate
// go/parser code invariants.
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
// markdown instruction file — a declared allowlist of recognized path-carrying vars.
// A test that reads one of these via a read sink (os.ReadFile/os.Open/io.ReadAll/
// bufio) is reading an ingested file — the read triggers the must-declare rule
// regardless of how the bytes are then inspected. (Code-scanning tests reference
// ../dispatch, ../status package dirs, never these.) A path var NOT in this list is
// not recognized; whether such a read should declare is an instance of the
// reader-axis bound the detached adversarial audit backstops (see the sweep doc).
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

// TestNoUndeclaredHostneutralityTautology is the AC-3 sweep for this package,
// re-runnable offline. It parses every *_test.go and flags any test function that
// READS a recognized markdown INSTRUCTION file's content — via a named reader
// helper, a DIRECT os.ReadFile/os.Open/io read of a recognized instruction
// literal/segment/path-ident, or a WalkDir-collected `.md` — UNLESS it self-classifies
// via markNonAC or markCodeBoundInvariant. The go/parser code-scan invariants
// (host_neutrality_test.go's scanFile over .go source via parser.ParseFile, NOT a
// content read sink) and the spanHostQualified unit test are NOT flagged: they read no
// instruction file. The undeclared-offender count is the AC-3 metric; it must be zero.
//
// What the guard guarantees, and what it deliberately does NOT (two axes):
//
//   - MATCH axis (closed, universal, load-bearing — STATICALLY GUARDED): the sweep
//     keys on the READ, not on how the bytes are inspected. ONCE a read of a
//     recognized instruction file is detected, the test MUST declare regardless of
//     the inspection idiom (strings.Contains/Index/EqualFold, bytes.*,
//     regexp.Regexp.Match, a bare ==) — the trigger is the ingest, not the match, so
//     the whole match class is closed.
//
//   - READER axis (does an undeclared instruction-file read hide via an undiscovered
//     read SHAPE? — DETACHED-AUDIT-BACKSTOPPED, NOT statically guarded): a read is
//     detected only when it is DIRECT — an in-package read sink whose path arg subtree
//     carries a recognized instruction literal/segment or an instructionPathIdent
//     package var, a named instructionFileReaders helper, or a WalkDir `.md`
//     collector. A read reached through any other shape — param/local/field/method/
//     closure taint flow, a transitive helper chain, a `[]string`/range-element flow,
//     a cross-package reader, a package var defined in another file, or an
//     unrecognized surface (AGENTS.md, mods/*.md) — is NOT statically flagged. This is
//     a deliberate, documented boundary: a per-package go/ast scan structurally cannot
//     see a cross-package read or a path built in another file, so chasing reader-shape
//     completeness was the enumeration trap the proof policy itself warns against. The
//     reader axis is covered by the detached adversarial audit required at every
//     high-stakes-surface gate (the validation-stage policy), NOT by this sweep. The
//     prior cycles' M-A/B/C/D shapes AND the two declared range-var/cross-statement HN
//     reads (TestDevDisciplinesSurviveInDevHomes, TestLiveScenarioRecommendedPractice-
//     Present — both carry markNonAC) are this audited boundary, not silent gaps.
func TestNoUndeclaredHostneutralityTautology(t *testing.T) {
	offenders := sweepHostneutralityTautologies(t, ".")
	for _, o := range offenders {
		t.Errorf("%s reads a markdown instruction file's content without self-classifying — call markNonAC (with its behavioral oracle) or markCodeBoundInvariant (with its independent source); how the bytes are inspected does not matter", o)
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

	// First pass: the recognized reader set is the named instructionFileReaders
	// allowlist plus any non-test helper that DIRECTLY reads a recognized instruction
	// file (readsInstructionContent — a read sink whose path arg carries a recognized
	// instruction literal/segment/path-ident, or a WalkDir-collected `.md`). The
	// code-scan helper scanFile is NOT a reader: it uses parser.ParseFile (not a
	// content read sink) over a `../dispatch` path (no instruction literal).
	readers := map[string]bool{}
	for r := range instructionFileReaders {
		readers[r] = true
	}
	for _, f := range files {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || strings.HasPrefix(fn.Name.Name, "Test") {
				continue
			}
			if readsInstructionContent(fn) {
				readers[fn.Name.Name] = true
			}
		}
	}

	// Second pass: a test is an offender if it ingests instruction-file content —
	// directly (readsInstructionContent) or via a recognized reader helper — and does
	// NOT declare its proof standing. The sweep keys on the READ, not a match-func
	// allowlist.
	var offenders []string
	for _, f := range files {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || !strings.HasPrefix(fn.Name.Name, "Test") {
				continue
			}
			calls := collectCalls(fn)
			readsInstruction := readsInstructionContent(fn)
			for r := range readers {
				if calls[r] {
					readsInstruction = true
					break
				}
			}
			declared := calls["markNonAC"] || calls["markCodeBoundInvariant"]
			if readsInstruction && !declared {
				offenders = append(offenders, fn.Name.Name)
			}
		}
	}
	return sortedUniqueHN(offenders)
}

// collectCalls walks a function body and returns the set of called function names
// (bare and selector trailing name), used to detect reader-helper calls and the
// markNonAC / markCodeBoundInvariant declarations.
func collectCalls(fn *ast.FuncDecl) map[string]bool {
	calls := map[string]bool{}
	ast.Inspect(fn, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			switch f := call.Fun.(type) {
			case *ast.Ident:
				calls[f.Name] = true
			case *ast.SelectorExpr:
				calls[f.Sel.Name] = true
			}
		}
		return true
	})
	return calls
}

// readSinks are the call selectors that ingest a file's content given a path.
var readSinks = map[string]bool{
	"ReadFile":   true, // os.ReadFile
	"Open":       true, // os.Open
	"ReadAll":    true, // io.ReadAll
	"NewScanner": true, // bufio.NewScanner
	"NewReader":  true, // bufio.NewReader
}

// readsInstructionContent reports whether fn DIRECTLY ingests a recognized
// instruction file's content: a read sink (ReadFile/Open/ReadAll/bufio) whose path
// arg subtree carries a recognized instruction literal/segment (isInstructionPathLiteral)
// or an instructionPathIdent package var, or a WalkDir-collected instruction `.md`.
// This is the direct-read predicate — no taint-flow tracking. A read reached only
// through a param/local/field/method/closure flow, a transitive helper chain, or a
// range-element flow is NOT detected here; that reader-shape axis is the
// detached-audit-backstopped boundary documented on TestNoUndeclaredHostneutralityTautology.
func readsInstructionContent(fn *ast.FuncDecl) bool {
	found := false
	ast.Inspect(fn, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if readSinks[sel.Sel.Name] {
			for _, arg := range call.Args {
				if exprInstructionTainted(arg) {
					found = true
				}
			}
		}
		if (sel.Sel.Name == "WalkDir" || sel.Sel.Name == "Walk") && fnFiltersInstructionMarkdown(fn) {
			found = true
		}
		return true
	})
	return found
}

// exprInstructionTainted reports whether expr DIRECTLY carries an instruction-file
// path anywhere in its subtree: an instruction-path literal/segment, or a known
// instructionPathIdent package var — so the +/strings.Join/filepath.Join/fmt.Sprintf/
// string(...) path-build idioms (whose instruction operand is a subtree node) are
// covered when their operands are literals or recognized vars.
func exprInstructionTainted(expr ast.Expr) bool {
	hit := false
	ast.Inspect(expr, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.BasicLit:
			if x.Kind == token.STRING && isInstructionPathLiteral(strings.Trim(x.Value, "`\"")) {
				hit = true
			}
		case *ast.Ident:
			if instructionPathIdents[x.Name] {
				hit = true
			}
		}
		return true
	})
	return hit
}

// fnFiltersInstructionMarkdown reports whether fn's body filters paths by a `.md`
// suffix — the WalkDir-collector signal.
func fnFiltersInstructionMarkdown(fn *ast.FuncDecl) bool {
	hit := false
	ast.Inspect(fn, func(n ast.Node) bool {
		if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			if strings.HasSuffix(strings.Trim(lit.Value, "`\""), ".md") {
				hit = true
			}
		}
		return true
	})
	return hit
}

// instructionPathSegments are the skill-tree / contract path segments that mark a
// path literal as an instruction file. The RECOGNIZED-instruction-surface predicate
// (a deliberate bound, not universal): a path carrying one of these listed segments
// is an instruction path, recognized even before a `.md` suffix is appended (so a
// Join/split-built suffix still matches on the base segment).
//
// A real instruction surface whose path carries NONE of these segments (e.g.
// AGENTS.md, mods/*.md) is not recognized and a read of it is not statically flagged.
// That reader-shape gap is the detached-audit-backstopped boundary documented on
// TestNoUndeclaredHostneutralityTautology — covered by the adversarial audit at the
// high-stakes gate, not by enumerating more segments here.
var instructionPathSegments = map[string]bool{
	"skills":        true,
	"references":    true,
	"agents":        true,
	"first-officer": true,
	"ensign":        true,
	"commission":    true,
	"present-gate":  true,
	"SKILL.md":      true,
}

// isInstructionPathLiteral reports whether a string literal is (a fragment of) an
// instruction-file path: it carries a skill-tree/contract segment. A `.json`
// manifest path carries none and is not instruction.
func isInstructionPathLiteral(s string) bool {
	if strings.HasSuffix(s, ".json") {
		return false
	}
	for seg := range instructionPathSegments {
		if s == seg || strings.Contains(s, seg) {
			return true
		}
	}
	return false
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
