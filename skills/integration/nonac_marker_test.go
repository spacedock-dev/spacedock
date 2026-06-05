// ABOUTME: The non-AC text-consistency marker + the AC-3 sweep meta-test — a
// ABOUTME: presence/absence check over an LLM-ingested instruction file is proof only if it declares its behavioral oracle.
package integration

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"testing"
)

// markNonAC is the explicit demotion seam required by the proof policy (f8b257cf):
// a string/substring/regex match over an instruction file the model reads NEVER
// satisfies a behavioral acceptance criterion. A test that matches such a file is
// legitimate ONLY as a text-consistency sanity check (the prose moved, a clause is
// present, a token is absent) — never as the proof of behavior. Calling
// markNonAC(t, oracle) declares that:
//   - this test is a non-AC text-consistency lint, NOT a behavioral proof, and
//   - `oracle` names where the behavior it touches is ACTUALLY proven (a live
//     drive, a code-side invariant, or "n/a — pure text property" when the claim
//     is itself about the text and has no behavior to drive).
//
// The AC-3 sweep meta-test (TestNoUndeclaredTautologicalProof) keys on this call:
// any test in this package that matches an ingested instruction file but does NOT
// call markNonAC is flagged as an undeclared tautology standing in for a
// behavioral claim. The call itself does nothing at runtime — its value is the
// declaration the sweep reads from the source.
func markNonAC(t *testing.T, oracle string) {
	t.Helper()
	if oracle == "" {
		t.Fatal("markNonAC requires a non-empty behavioral oracle reference")
	}
}

// markCodeBoundInvariant is the second of the two explicit classifications the
// proof policy's litmus demands ("does the expected value come from a source OTHER
// than the file under test?"). A test calls markCodeBoundInvariant(t, source) to
// declare that its expectation is NOT a literal hardcoded against the file under
// test but is read from `source` — an independent code-side value (a shared Go
// const, the seam target the Skill() invocation uses, a manifest the binary
// parses) that can DIVERGE from the file. That divergence is exactly what makes
// the check able to fail as an invariant, so it is a legitimate AC-2 invariant,
// not a tautology. The AC-3 sweep treats a markCodeBoundInvariant test as
// declared (not an offender), the same as a markNonAC text-consistency lint —
// every text-matching test must self-classify as one or the other. The call does
// nothing at runtime; its value is the source-level declaration the sweep reads.
func markCodeBoundInvariant(t *testing.T, source string) {
	t.Helper()
	if source == "" {
		t.Fatal("markCodeBoundInvariant requires a non-empty independent-source reference")
	}
}

// ingestedFileReaders are the named helper functions in this package that read an
// instruction file the model ingests (a contract, a workflow README, a skill
// body) — the recognized-reader allowlist. A test that calls one of these (the READ
// — how it then inspects the bytes is irrelevant under the match-axis positive rule)
// is a presence/absence check over an ingested file, exactly the shape the proof
// policy bans as standalone behavioral proof. The AC-3 sweep treats such a test as
// tautological unless it declares markNonAC, or unless it binds the expected value to
// a code-side source (a re-bound Bucket-B invariant — markCodeBoundInvariant). The
// sweep distinguishes the two by the declaration: a re-bound invariant whose
// expectation diverges from the file declares markCodeBoundInvariant; a pure
// text-consistency lint declares markNonAC.
var ingestedFileReaders = map[string]bool{
	"foCore":                            true,
	"foRuntime":                         true,
	"presentGateSkill":                  true,
	"feedbackRejectionFlowSkill":        true,
	"usingClaudeTeamSkill":              true,
	"vendoredSkillFiles":                true,
	"presentGateFrontmatterValue":       true,
	"feedbackRejectionFrontmatterValue": true,
}

// TestNoUndeclaredTautologicalProof is the AC-3 sweep, re-runnable offline. It
// parses every *_test.go in this package and flags any test function that READS a
// recognized instruction file's content — via a named reader helper, a DIRECT
// os.ReadFile/os.Open of a recognized instruction literal/segment, or a
// WalkDir-collected `.md` — unless it self-classifies via markNonAC or
// markCodeBoundInvariant. The count of undeclared offenders is the AC-3 metric: it
// must be zero.
//
// What the guard guarantees, and what it deliberately does NOT (two axes):
//
//   - MATCH axis (closed, universal, load-bearing — STATICALLY GUARDED): the sweep
//     keys on the READ, not on how the bytes are then inspected. ONCE a read of a
//     recognized instruction file is detected, the test MUST declare regardless of
//     the inspection idiom — strings.Contains/Index/EqualFold, bytes.*,
//     regexp.Regexp.Match, len(Split)>1, a bare `==`, anything. Enumerating "match
//     functions" was whack-a-mole; this rule closes the whole class because the
//     trigger is the ingest, not the match.
//
//   - READER axis (does an undeclared instruction-file read hide via an undiscovered
//     read SHAPE? — DETACHED-AUDIT-BACKSTOPPED, NOT statically guarded): a read is
//     detected only when it is DIRECT — an in-package read sink whose path arg subtree
//     carries a recognized instruction literal/segment (isInstructionPathLiteral), a
//     named ingestedFileReaders helper, or a WalkDir `.md` collector. A read reached
//     through any other shape — param/local/field/method/closure taint flow, a
//     transitive helper chain, a `[]string`/range-element flow, a cross-package
//     reader, or a path in a package var defined in another file — is NOT statically
//     flagged. This is a deliberate, documented boundary: a per-package go/ast scan
//     structurally cannot see a cross-package read or a path built in another file, so
//     chasing reader-shape completeness was the enumeration trap the proof policy
//     itself warns against. The reader axis is covered by the detached adversarial
//     audit required at every high-stakes-surface gate (the validation-stage policy),
//     NOT by this sweep. The prior cycles' M-A (AGENTS.md, mods/*.md) / M-B / M-C / M-D
//     shapes are this audited boundary, not silent gaps.
//
// This sweep is itself a code-side invariant over real parsed test source, not a
// text match over an instruction file — its expected value (which reads directly
// reach an instruction file) is independent of any contract prose, so it can fail
// when a future edit adds an undeclared DIRECT ingest.
func TestNoUndeclaredTautologicalProof(t *testing.T) {
	offenders := sweepUndeclaredTautologies(t, ".")
	for _, o := range offenders {
		t.Errorf("%s reads an ingested instruction file's content without calling markNonAC or markCodeBoundInvariant — declare it a non-AC text-consistency lint (with its behavioral oracle) or re-bind its expectation to a code-side source; how the bytes are inspected does not matter", o)
	}
	if len(offenders) > 0 {
		t.Fatalf("AC-3 sweep: %d undeclared tautological-behavioral-proof test(s); the count must be zero", len(offenders))
	}
}

// sweepUndeclaredTautologies returns the sorted names of test functions in the
// package at dir that read an ingested instruction file, match it with a
// substring/regex call, and do NOT declare markNonAC. Exported as a helper so the
// sweep's own mutation control (TestSweepDetectsAnUndeclaredTautology) can call it
// against a synthetic fixture.
func sweepUndeclaredTautologies(t *testing.T, dir string) []string {
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

	// First pass: the recognized reader set is the named ingestedFileReaders allowlist
	// plus any non-test helper that DIRECTLY reads an instruction file
	// (readsInstructionContent — a read sink whose path arg carries a recognized
	// instruction literal/segment, or a WalkDir-collected `.md`). The seeded named
	// helpers cover readers that return a non-`.md`-literal handle (vendoredSkillFiles
	// returns a map).
	readers := map[string]bool{}
	for r := range ingestedFileReaders {
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
	// NOT declare its proof standing. The sweep keys on the READ, not on a match-func
	// allowlist: any inspection of ingested bytes (Contains/Index/EqualFold, bytes.*,
	// regexp.Match, a bare ==, …) is covered because the trigger is the ingest itself.
	var offenders []string
	for _, f := range files {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || !strings.HasPrefix(fn.Name.Name, "Test") {
				continue
			}
			calls := collectCalls(fn)
			readsIngested := readsInstructionContent(fn)
			for r := range readers {
				if calls[r] {
					readsIngested = true
					break
				}
			}
			declared := calls["markNonAC"] || calls["markCodeBoundInvariant"]
			if readsIngested && !declared {
				offenders = append(offenders, fn.Name.Name)
			}
		}
	}
	return sortedUnique(offenders)
}

// readsInstructionContent reports whether fn DIRECTLY ingests a recognized
// instruction file's content: a read sink (os.ReadFile/os.Open/io.ReadAll/bufio
// scanner-reader) whose path arg subtree carries a recognized instruction
// literal/segment (isInstructionPathLiteral), or a WalkDir/Walk over a tree
// collecting instruction `.md` files (the reader-of-many shape its callers read).
// This is the direct-read predicate — no taint-flow tracking. A read reached only
// through a param/local/field/method/closure flow, a transitive helper chain, or a
// range-element flow is NOT detected here; that reader-shape axis is the
// detached-audit-backstopped boundary documented on TestNoUndeclaredTautologicalProof.
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
		// A WalkDir/Walk that filters on an instruction `.md` path is a reader-of-many:
		// it collects the paths its callers read+inspect.
		if (sel.Sel.Name == "WalkDir" || sel.Sel.Name == "Walk") && fnFiltersInstructionMarkdown(fn) {
			found = true
		}
		return true
	})
	return found
}

// readSinks are the call selectors that ingest a file's content given a path: the
// os reads, io.ReadAll over an opened handle, and the bufio scanner/reader
// constructors. A tainted instruction path flowing into any of these is an ingest.
var readSinks = map[string]bool{
	"ReadFile":   true, // os.ReadFile
	"Open":       true, // os.Open
	"ReadAll":    true, // io.ReadAll
	"NewScanner": true, // bufio.NewScanner
	"NewReader":  true, // bufio.NewReader
}

// exprInstructionTainted reports whether an expression DIRECTLY carries an
// instruction-file path: it references an instruction-path string literal/segment
// (isInstructionPathLiteral) anywhere in its subtree — so the `+` / strings.Join /
// filepath.Join / fmt.Sprintf path-build idioms (whose instruction operand is a node
// in the subtree) are covered when their operands are literals.
func exprInstructionTainted(expr ast.Expr) bool {
	hit := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if x, ok := n.(*ast.BasicLit); ok {
			if x.Kind == token.STRING && isInstructionPathLiteral(strings.Trim(x.Value, "`\"")) {
				hit = true
			}
		}
		return true
	})
	return hit
}

// fnFiltersInstructionMarkdown reports whether fn's body filters paths by an
// instruction `.md` suffix — the WalkDir-collector signal. A `.md` HasSuffix check
// (or an instruction-`.md` literal) anywhere in a WalkDir helper marks it a
// reader-of-many over the instruction surface.
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

// collectCalls walks a function body and returns the set of called function names
// (bare `foo(...)` and selector `pkg.Foo(...)`/`recv.Method(...)` trailing name).
// Used to detect calls to discovered reader helpers and the markNonAC /
// markCodeBoundInvariant declarations.
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

// instructionPathSegments are the skill-tree / contract path segments that mark a
// path literal as targeting an instruction file the model ingests (a skill,
// contract, agent, or runtime adapter) rather than a binary-parsed artifact (a
// manifest .json) or a dev-only doc (docs/dev/*.md recipes are NOT an LLM
// instruction surface and are intentionally out of scope — Cycle-2 P1 divergence:
// the sweep scopes to the shipped skill/contract surface, not every `.md` in the
// repo).
//
// This is the RECOGNIZED-instruction-surface predicate (a deliberate bound, not a
// universal one): a path carrying one of these listed segments is an instruction
// path. A path fragment carrying a segment is recognized even before a `.md` suffix
// is appended, so a literal "…/first-officer-shared-core" matches on its segment.
//
// A real instruction surface whose path carries NONE of these segments — e.g.
// AGENTS.md or mods/*.md — is not recognized and a read of it is not statically
// flagged. That reader-shape gap is the detached-audit-backstopped boundary
// documented on TestNoUndeclaredTautologicalProof — covered by the adversarial audit
// at the high-stakes gate, not by enumerating more segments here.
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
// manifest path or a docs/dev recipe path carries none and is not instruction.
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

// TestSweepDetectsAnUndeclaredTautology is the mutation control for the AC-3
// sweep itself: the sweep is the AC-3 oracle, so it must be demonstrated to RED on
// the exact shape it polices and GREEN once that shape declares its demotion.
// Without this, the sweep could silently degrade to a no-op (e.g. an ingested-file
// reader renamed out of the map) and pass vacuously. It writes two synthetic test
// files to a temp dir and runs the sweep against it:
//   - undeclared: a test that reads an ingested file (foCore) and matches it
//     (strings.Contains) but never calls markNonAC -> MUST be flagged.
//   - declared: the same shape plus a markNonAC call -> MUST NOT be flagged.
func TestSweepDetectsAnUndeclaredTautology(t *testing.T) {
	dir := t.TempDir()
	undeclared := `package fixture
import "strings"
func TestUndeclaredFixture(t *T) {
	fo := foCore(t)
	if strings.Contains(fo, "x") { _ = fo }
}
`
	declared := `package fixture
func TestDeclaredFixture(t *T) {
	markNonAC(t, "behavioral oracle: live gate-guardrail scenario")
	fo := foCore(t)
	if strings.Contains(fo, "x") { _ = fo }
}
`
	writeFile(t, dir+"/undeclared_test.go", undeclared)
	offenders := sweepUndeclaredTautologies(t, dir)
	if !containsStr(offenders, "TestUndeclaredFixture") {
		t.Fatalf("sweep failed to flag an undeclared presence-check over an ingested file; offenders=%v", offenders)
	}

	writeFile(t, dir+"/declared_test.go", declared)
	offenders = sweepUndeclaredTautologies(t, dir)
	if containsStr(offenders, "TestDeclaredFixture") {
		t.Fatalf("sweep wrongly flagged a declared (markNonAC) text-consistency lint; offenders=%v", offenders)
	}
	if !containsStr(offenders, "TestUndeclaredFixture") {
		t.Fatalf("adding a declared fixture must not stop the sweep flagging the undeclared one; offenders=%v", offenders)
	}

	// A code-bound invariant (expectation from an independent source) is the other
	// valid self-classification and must clear the sweep too.
	codeBound := `package fixture
func TestCodeBoundFixture(t *T) {
	markCodeBoundInvariant(t, "shared const presentGateSeamName")
	fo := foCore(t)
	if strings.Contains(fo, presentGateSeamName) { _ = fo }
}
`
	writeFile(t, dir+"/codebound_test.go", codeBound)
	offenders = sweepUndeclaredTautologies(t, dir)
	if containsStr(offenders, "TestCodeBoundFixture") {
		t.Fatalf("sweep wrongly flagged a code-bound invariant; offenders=%v", offenders)
	}
}

func containsStr(in []string, want string) bool {
	for _, s := range in {
		if s == want {
			return true
		}
	}
	return false
}

func sortedUnique(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	// simple insertion sort — the lists are tiny
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j-1] > out[j]; j-- {
			out[j-1], out[j] = out[j], out[j-1]
		}
	}
	return out
}
