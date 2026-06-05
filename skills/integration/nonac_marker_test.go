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

// ingestedFileReaders are the helper functions in this package that read an
// instruction file the model ingests (a contract, a workflow README, a skill
// body). A test that calls one of these AND matches its result with
// strings.Contains / strings.Count / a regexp is a presence/absence check over an
// ingested file — exactly the shape the proof policy bans as standalone
// behavioral proof. The AC-3 sweep treats such a test as tautological-unless it
// declares markNonAC, or unless it instead binds the expected value to a code-side
// source (a re-bound Bucket-B invariant — those read the constant, not just the
// file). The sweep distinguishes the two by the markNonAC declaration: a re-bound
// invariant whose expectation diverges from the file does not need the marker,
// because its expected value is independent; a pure text-consistency lint does.
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

// matchFuncs are the substring/regex match calls that, applied to an ingested
// file's text, make a test a presence/absence check.
var matchFuncs = map[string]bool{
	"Contains":              true, // strings.Contains
	"Count":                 true, // strings.Count
	"HasPrefix":             true,
	"HasSuffix":             true,
	"FindString":            true, // regexp
	"FindStringSubmatch":    true,
	"MatchString":           true,
	"FindAllStringSubmatch": true,
}

// TestNoUndeclaredTautologicalProof is the AC-3 sweep, re-runnable offline. It
// parses every *_test.go in this package and, for each test function, decides
// whether it is a presence/absence check over an ingested instruction file (it
// both calls an ingestedFileReader AND a substring/regex matchFunc). Such a test
// is permitted only if it ALSO calls markNonAC — the explicit demotion declaring
// it a non-AC text-consistency lint and naming its behavioral oracle. The sweep
// fails (lists the offenders) if any test matches an ingested file as a standalone
// claim without that declaration. The count of undeclared offenders is the AC-3
// metric: it must be zero.
//
// This sweep is itself a code-side invariant over real parsed test source, not a
// text match over an instruction file — its expected value (which functions read
// ingested files, which functions match) is independent of any contract prose, so
// it can fail when a future edit adds an undeclared tautology.
func TestNoUndeclaredTautologicalProof(t *testing.T) {
	offenders := sweepUndeclaredTautologies(t, ".")
	for _, o := range offenders {
		t.Errorf("%s reads an ingested instruction file and matches it as a standalone claim without calling markNonAC — declare it a non-AC text-consistency lint (with its behavioral oracle) or re-bind its expectation to a code-side source", o)
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

	// First pass: discover the package's instruction-file reader helpers — any
	// non-test func that directly reads a `.md` skill-tree file. Seeded with the
	// named helpers (some, like vendoredSkillFiles, return a map rather than
	// reading a `.md` literal directly). This auto-detection keeps the sweep from
	// silently drifting when a new foSharedCore/claudeFORuntime-style reader is
	// added — the new helper is discovered, so a test calling it is still caught.
	readers := map[string]bool{}
	for r := range ingestedFileReaders {
		readers[r] = true
	}
	// Collect every non-test func's calls + strings once, then grow the reader set
	// to a fixpoint: a func is a reader if it directly reads a `.md` skill-tree file
	// OR it calls a known reader (transitive, e.g. a test → startupStep1 →
	// foSharedCore → os.ReadFile chain). The fixpoint closes the multi-hop helper
	// gap an adversary could otherwise hide a tautology behind.
	helperCalls := map[string]map[string]bool{}
	for _, f := range files {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || strings.HasPrefix(fn.Name.Name, "Test") {
				continue
			}
			calls, strs := collectCallsAndStrings(fn)
			helperCalls[fn.Name.Name] = calls
			if (calls["ReadFile"] || calls["Open"]) && readsInstructionPath(strs) {
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
			calls, strs := collectCallsAndStrings(fn)
			readsIngested := false
			for r := range readers {
				if calls[r] {
					readsIngested = true
					break
				}
			}
			// A direct os.ReadFile/os.Open in the test body whose path is an
			// instruction-file location is also an ingested-file read. The path is an
			// instruction file when the test carries a `.md`-suffixed string literal
			// AND a skill-tree path segment. A `.json` manifest read (parsed by the
			// binary, Bucket D) does not match.
			if (calls["ReadFile"] || calls["Open"]) && readsInstructionPath(strs) {
				readsIngested = true
			}
			matches := false
			for m := range matchFuncs {
				if calls[m] {
					matches = true
					break
				}
			}
			declared := calls["markNonAC"] || calls["markCodeBoundInvariant"]
			if readsIngested && matches && !declared {
				offenders = append(offenders, fn.Name.Name)
			}
		}
	}
	return sortedUnique(offenders)
}

// collectCallsAndStrings walks a function body and returns the set of called
// function names (bare `foo(...)` and selector `pkg.Foo(...)` trailing name) plus
// the set of string-literal values. The call set detects ingested-file reads,
// substring/regex matches, and the markNonAC/markCodeBoundInvariant declarations;
// the string set detects a direct os.ReadFile of an instruction-file path literal.
func collectCallsAndStrings(fn *ast.FuncDecl) (calls map[string]bool, strs map[string]bool) {
	calls = map[string]bool{}
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
		case *ast.BasicLit:
			if node.Kind == token.STRING {
				strs[strings.Trim(node.Value, "`\"")] = true
			}
		}
		return true
	})
	return calls, strs
}

// readInstructionPathSegments are skill-tree path segments that, alongside a `.md`
// literal in the same test, mark a direct os.ReadFile as targeting an instruction
// file the model ingests (a skill, contract, agent, or runtime adapter) rather
// than a binary-parsed artifact (a manifest .json).
var readInstructionPathSegments = map[string]bool{
	"references":    true,
	"agents":        true,
	"first-officer": true,
	"ensign":        true,
	"commission":    true,
	"present-gate":  true,
	"SKILL.md":      true,
}

// readsInstructionPath reports whether a test's string literals indicate a read of
// an instruction file: a `.md`-suffixed literal plus a skill-tree segment (or a
// full skill-tree path in one literal). A `.json` read never matches.
func readsInstructionPath(strs map[string]bool) bool {
	hasMD := false
	hasSegment := false
	for s := range strs {
		if strings.HasSuffix(s, ".md") {
			hasMD = true
		}
		for seg := range readInstructionPathSegments {
			if s == seg || strings.Contains(s, seg) {
				hasSegment = true
			}
		}
	}
	return hasMD && hasSegment
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
