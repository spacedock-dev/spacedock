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
// body) — the seed of the reader set the sweep grows to a fixpoint. A test that
// calls one of these (the READ — how it then inspects the bytes is irrelevant under
// the match-axis positive rule) is a presence/absence check over an ingested file,
// exactly the shape the proof policy bans as standalone behavioral proof. The AC-3
// sweep treats such a test as tautological unless it declares markNonAC, or unless
// it binds the expected value to a code-side source (a re-bound Bucket-B invariant —
// markCodeBoundInvariant). The sweep distinguishes the two by the declaration: a
// re-bound invariant whose expectation diverges from the file declares
// markCodeBoundInvariant; a pure text-consistency lint declares markNonAC.
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
// recognized instruction file's content — via a reader helper, a tainted
// os.ReadFile/os.Open, or a WalkDir-collected `.md`, through the flow shapes the
// reader-axis taint covers (below) — unless it self-classifies via markNonAC or
// markCodeBoundInvariant. The count of undeclared offenders is the AC-3 metric: it
// must be zero.
//
// What the guard actually guarantees (two axes, with one closed and one bounded):
//
//   - MATCH axis (closed, universal, load-bearing): the sweep keys on the READ, not
//     on how the bytes are then inspected. ONCE a read of a recognized instruction
//     file is detected, the test MUST declare regardless of the inspection idiom —
//     strings.Contains/Index/EqualFold, bytes.*, regexp.Regexp.Match, len(Split)>1,
//     a bare `==`, anything. Enumerating "match functions" was whack-a-mole; this
//     rule closes the whole class because the trigger is the ingest, not the match.
//
//   - READER axis (covered flow shapes, NOT exhaustive): a read is detected for an
//     in-package read of a RECOGNIZED instruction path (a skill-tree/contract
//     segment, isInstructionPathLiteral) reaching a read sink through these flows: a
//     bare-`string` parameter, a `:=`/`=` local, a struct field, a method receiver,
//     a closure capture; with the path built by `+` / strings.Join / filepath.Join /
//     fmt.Sprintf. A transitive helper chain is followed to a fixpoint.
//
// KNOWN OUT-OF-SCOPE (tracked in the follow-up task sweep-guard-reader-axis-invert,
// id 4qnn7dbzkyh9qv65t618vtxy, backstopped by the detached adversarial audit before
// merge — NOT silently dropped):
//   - M-A: unrecognized instruction surfaces (AGENTS.md, mods/*.md) — not in the
//     instructionPathSegments predicate.
//   - M-B: cross-package reads (a read whose reader helper lives in another package).
//   - M-C: a path held in a package var defined in another file of this package.
//   - M-D: `[]string`/`...string`-param + range/slice-element flow.
// These are the same recurring enumerated-shape reader-flow class cycles 1-3 each
// closed instances of; the follow-up weighs an invert/positive predicate and a
// go/types+SSA taint that closes the class definitionally.
//
// This sweep is itself a code-side invariant over real parsed test source, not a
// text match over an instruction file — its expected value (which reads reach an
// instruction file) is independent of any contract prose, so it can fail when a
// future edit adds an undeclared ingest through a covered flow.
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

	// First pass: discover the package's instruction-file reader helpers, then grow
	// the set to a fixpoint so a read cannot hide behind a helper chain. A func is a
	// reader if it ingests instruction-file content directly (readsInstructionContent
	// — a tainted os.ReadFile/io read, or a WalkDir-collected `.md`) OR (transitive)
	// it calls a known reader. Methods are NOT skipped: a reader can be a method on a
	// fixture struct (the s.path / method-receiver flow shape). The seeded named
	// helpers cover readers that return a non-`.md`-literal handle (vendoredSkillFiles
	// returns a map).
	taintedFields := instructionTaintedFields(files)
	readers := map[string]bool{}
	for r := range ingestedFileReaders {
		readers[r] = true
	}
	helperCalls := map[string]map[string]bool{}
	for _, f := range files {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || strings.HasPrefix(fn.Name.Name, "Test") {
				continue
			}
			helperCalls[fn.Name.Name] = collectCalls(fn)
			if readsInstructionContent(fn, taintedFields) {
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

	// Second pass: a test is an offender if it ingests instruction-file content —
	// directly (readsInstructionContent) or via a discovered reader helper — and does
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
			readsIngested := readsInstructionContent(fn, taintedFields)
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

// readsInstructionContent reports whether fn ingests a recognized instruction
// file's content through the reader-axis flow shapes the taint COVERS — it is the
// positive/taint replacement for the Cycle-1/2 allow-lists (readsParamPath +
// walksForMarkdown + constStringConcat-only `+` concat), but it covers a bounded
// set of flows, not an exhaustive one. It taints a string derived from a recognized
// instruction-file path (a skill-tree/contract segment, isInstructionPathLiteral)
// built by `+` / strings.Join / filepath.Join / fmt.Sprintf and flowed through a
// bare-`string` param, a `:=`/`=` local, a struct field, or a method receiver, and
// reports a read when a tainted path flows into a read sink (os.ReadFile/os.Open/
// io.ReadAll/bufio scanner-reader), or when fn WalkDir/Walks a tree collecting
// instruction `.md` files (the reader-of-many shape its callers then read).
//
// NOT covered (tracked in sweep-guard-reader-axis-invert, id
// 4qnn7dbzkyh9qv65t618vtxy, audit-backstopped): `[]string`/`...string`-param +
// range/slice-element flow (M-D), cross-package reader helpers (M-B), a package var
// defined in another file of this package (M-C), and unrecognized surfaces like
// AGENTS.md / mods/*.md (M-A). See TestNoUndeclaredTautologicalProof's doc for the
// full honest bound.
func readsInstructionContent(fn *ast.FuncDecl, taintedFields map[string]bool) bool {
	tainted := instructionTaintedNames(fn, taintedFields)
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
				if exprInstructionTainted(arg, tainted) || readsTaintedField(arg, taintedFields) {
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

// readsTaintedField reports whether expr reads a struct field whose name is in the
// package-wide instruction-tainted-field set — the s.path / method-receiver flow
// (Cycle-3 M3): the `.md` literal is assigned to the field in a constructor in one
// function and the read happens via a field selector in another (a method). Field
// taint is computed package-wide (instructionTaintedFields), so this catches the
// read even though the assigning literal is not in fn's own body. A generated-path
// field (runBuild's res.DispatchFilePath) is never instruction-assigned, so it is
// not in the set and not flagged.
func readsTaintedField(expr ast.Expr, taintedFields map[string]bool) bool {
	hit := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectorExpr); ok && taintedFields[sel.Sel.Name] {
			hit = true
		}
		return true
	})
	return hit
}

// instructionTaintedFields scans every struct composite literal and every
// assignment to a field selector across the package, returning the set of FIELD
// NAMES ever assigned an instruction-file path. Keyed by field name (no type info
// at parse time) — a deliberate over-approximation that errs toward flagging, which
// the proof policy wants. A field only ever assigned a generated/temp path (a
// dispatch artifact) never enters the set.
func instructionTaintedFields(files []*ast.File) map[string]bool {
	fields := map[string]bool{}
	for _, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.KeyValueExpr:
				if key, ok := node.Key.(*ast.Ident); ok {
					if exprInstructionTainted(node.Value, nil) {
						fields[key.Name] = true
					}
				}
			case *ast.AssignStmt:
				for i, rhs := range node.Rhs {
					if i >= len(node.Lhs) {
						break
					}
					if sel, ok := node.Lhs[i].(*ast.SelectorExpr); ok && exprInstructionTainted(rhs, nil) {
						fields[sel.Sel.Name] = true
					}
				}
			}
			return true
		})
	}
	return fields
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

// instructionTaintedNames computes the set of identifier names (params, locals,
// struct-field selectors rendered as `recv.field`, range vars) in fn that hold a
// string derived from an instruction-file path. It seeds from instruction-path
// expressions (a `.md` skill-tree literal/segment, an instructionPathSegment, a
// known instruction ident) and propagates through := / = assignments and string
// conversions to a fixpoint, so a path built and then read in separate statements
// is still tainted at the read.
func instructionTaintedNames(fn *ast.FuncDecl, taintedFields map[string]bool) map[string]bool {
	tainted := map[string]bool{}
	// Seed: any parameter is a candidate taint carrier only if the CALLER supplies an
	// instruction path; within fn we cannot see the caller, so a reader-helper whose
	// path arg is a parameter is caught by the param-flow rule below (the parameter is
	// tainted when fn itself also references an instruction literal, OR unconditionally
	// for a single-string-param helper that reads it — the readSkill(t, path) shape).
	// We treat every string parameter as tainted: a helper that ReadFiles a string
	// param is, by construction, a path-arg reader (the caller supplies the .md path).
	if fn.Type.Params != nil {
		for _, field := range fn.Type.Params.List {
			if isStringyType(field.Type) {
				for _, name := range field.Names {
					tainted[name.Name] = true
				}
			}
		}
	}
	for grew := true; grew; {
		grew = false
		ast.Inspect(fn, func(n ast.Node) bool {
			assign, ok := n.(*ast.AssignStmt)
			if !ok {
				return true
			}
			for i, rhs := range assign.Rhs {
				if i >= len(assign.Lhs) {
					break
				}
				// A local assigned from an instruction-tainted expr, OR from a read of a
				// package-wide instruction-tainted field, carries the taint forward.
				if !exprInstructionTainted(rhs, tainted) && !readsTaintedField(rhs, taintedFields) {
					continue
				}
				if name := lvalueName(assign.Lhs[i]); name != "" && !tainted[name] {
					tainted[name] = true
					grew = true
				}
			}
			return true
		})
	}
	return tainted
}

// lvalueName renders an assignable target as a taint-tracking key: a bare ident, or
// a selector `recv.field` (the struct-field path-flow shape).
func lvalueName(e ast.Expr) string {
	switch x := e.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.SelectorExpr:
		if inner, ok := x.X.(*ast.Ident); ok {
			return inner.Name + "." + x.Sel.Name
		}
		return x.Sel.Name
	}
	return ""
}

// exprInstructionTainted reports whether an expression carries an instruction-file
// path taint: it references a tainted name (ident or `recv.field` selector), an
// instruction-path string literal/segment, or a known instruction path ident,
// anywhere in its subtree — so the `+` / strings.Join / filepath.Join / fmt.Sprintf
// path-build idioms (whose tainted operand is a node in the subtree) are covered.
// The over-approximation toward flagging is deliberate. It does NOT cover a taint
// carried in a slice element or recovered via a range variable (M-D) — see
// readsInstructionContent's NOT-covered note and the follow-up
// sweep-guard-reader-axis-invert.
func exprInstructionTainted(expr ast.Expr, tainted map[string]bool) bool {
	hit := false
	ast.Inspect(expr, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.BasicLit:
			if x.Kind == token.STRING && isInstructionPathLiteral(strings.Trim(x.Value, "`\"")) {
				hit = true
			}
		case *ast.Ident:
			if tainted[x.Name] {
				hit = true
			}
		case *ast.SelectorExpr:
			if inner, ok := x.X.(*ast.Ident); ok && tainted[inner.Name+"."+x.Sel.Name] {
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

// isStringyType reports whether a parameter type node carries a path string: a bare
// `string` (the readSkill(t, path) shape) — the kind a path-arg reader takes.
func isStringyType(t ast.Expr) bool {
	id, ok := t.(*ast.Ident)
	return ok && id.Name == "string"
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
// path. A path fragment carrying a segment is instruction-tainted even before a
// `.md` suffix is appended, so strings.Join([]string{"…/first-officer-shared-core",
// "md"}, ".") taints on the segment in the base (closing the Cycle-1
// `.md`-suffix-AND-segment pair a split/Join-built suffix evaded).
//
// KNOWN OUT-OF-SCOPE surfaces (M-A, tracked in sweep-guard-reader-axis-invert, id
// 4qnn7dbzkyh9qv65t618vtxy): a real instruction surface whose path carries NONE of
// these segments — e.g. AGENTS.md or mods/*.md — is not recognized and a read of it
// is not flagged. The follow-up weighs an invert/positive predicate that recognizes
// the instruction surface definitionally rather than by this enumerated list.
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

// TestSweepDetectsEvasionShapes is the planted-control mutation test for the
// reader-discovery evasion shapes the validation audit proved the sweep missed.
// Each case plants a synthetic offender that reaches an instruction file through a
// shape the naive `.md`-literal-in-the-reader detection cannot see, runs the sweep,
// and asserts it REDs (flags the offender); then it plants the declared form and
// asserts it GREENs. A regression that removed a discovery mechanism would let the
// matching case go un-flagged, failing this control.
func TestSweepDetectsEvasionShapes(t *testing.T) {
	// Shape 1 — path-arg reader: the helper os.ReadFile's a value built from its own
	// path parameter; the `.md` literal lives in the CALLER (the readSkill(t,root,rel)
	// shape). readsInstructionPath over the helper body sees no literal, so only
	// parameter-flow detection catches it.
	pathArg := `package fixture
import (
	"os"
	"path/filepath"
	"strings"
)
func readArg(t *T, root, rel string) string {
	b, _ := os.ReadFile(filepath.Join(root, rel))
	return string(b)
}
func TestPathArgOffender(t *T) {
	s := readArg(t, root, "first-officer/references/first-officer-shared-core.md")
	if strings.Contains(s, "x") { _ = s }
}
`
	assertRedThenGreen(t, "path-arg reader", "TestPathArgOffender", pathArg)

	// Shape 2 — WalkDir collector: the helper WalkDirs a tree collecting `.md`
	// paths and RETURNS them; it never os.ReadFile's the `.md` itself. The caller
	// reads+matches each returned path (the shippedSkillText shape).
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
func TestWalkDirOffender(t *T) {
	for _, p := range walkSkills(t, root) {
		b, _ := os.ReadFile(p)
		if strings.Contains(string(b), "x") { _ = b }
	}
}
`
	assertRedThenGreen(t, "WalkDir collector", "TestWalkDirOffender", walkDir)

	// Shape 3 — split-".md" suffix: the read path is constructed as
	// base + "." + "md", so no single literal carries the `.md` suffix. The
	// constant-concatenation reconstruction must rejoin it before .md detection.
	splitMD := `package fixture
import (
	"os"
	"path/filepath"
	"strings"
)
func TestSplitSuffixOffender(t *T) {
	p := filepath.Join(root, "first-officer", "references", "first-officer-shared-core" + "." + "md")
	b, _ := os.ReadFile(p)
	if strings.Contains(string(b), "x") { _ = b }
}
`
	assertRedThenGreen(t, "split-.md suffix", "TestSplitSuffixOffender", splitMD)

	// Shape 4 — multi-hop transitive helper: a tautology hidden two frames down
	// (the test calls wrapHop, which calls readArg, which reads a param path). The
	// reader fixpoint must propagate reader-ness up the call chain. This is the
	// integration-side guard that the transitive fixpoint stays load-bearing.
	multiHop := `package fixture
import (
	"os"
	"path/filepath"
	"strings"
)
func readArg2(t *T, root, rel string) string {
	b, _ := os.ReadFile(filepath.Join(root, rel))
	return string(b)
}
func wrapHop(t *T, root string) string {
	return readArg2(t, root, "ensign/references/ensign-shared-core.md")
}
func TestMultiHopOffender(t *T) {
	s := wrapHop(t, root)
	if strings.Contains(s, "x") { _ = s }
}
`
	assertRedThenGreen(t, "multi-hop transitive helper", "TestMultiHopOffender", multiHop)

	// Shape 5 (Cycle-3 M1, match axis) — the ingested bytes are inspected with a
	// match idiom OUTSIDE the old matchFuncs allowlist (strings.Index, regexp.Match
	// over []byte). The positive rule keys on the READ, so HOW the bytes are inspected
	// is irrelevant — the read of foCore alone must flag it regardless of the idiom.
	matchIndex := `package fixture
import "strings"
func TestMatchIndexOffender(t *T) {
	fo := foCore(t)
	if strings.Index(fo, "Skill(skill=\"spacedock:present-gate\")") < 0 { t.Error("x") }
}
`
	assertRedThenGreen(t, "match via strings.Index (no Contains)", "TestMatchIndexOffender", matchIndex)

	matchBytesRegexp := `package fixture
import "regexp"
func TestMatchBytesRegexpOffender(t *T) {
	fo := foCore(t)
	re := regexp.MustCompile("present-gate")
	if !re.Match([]byte(fo)) { t.Error("x") }
}
`
	assertRedThenGreen(t, "match via regexp.Regexp.Match([]byte)", "TestMatchBytesRegexpOffender", matchBytesRegexp)

	// Shape 6 (Cycle-3 M2, reader axis) — the `.md` path is built with strings.Join,
	// not `+`. The base fragment carries an instruction segment so it taints before
	// the suffix is appended; the constStringConcat-only-`+` design missed this.
	joinPath := `package fixture
import (
	"os"
	"strings"
)
func TestJoinPathOffender(t *T) {
	base := "../../skills/first-officer/references/first-officer-shared-core"
	p := strings.Join([]string{base, "md"}, ".")
	b, _ := os.ReadFile(p)
	if strings.Index(string(b), "x") < 0 { t.Error("y") }
}
`
	assertRedThenGreen(t, "strings.Join-built .md path", "TestJoinPathOffender", joinPath)

	// Shape 7 (Cycle-3 M3, reader axis) — the `.md` path flows through a struct field
	// and the read happens via a METHOD on that struct. readsParamPath tracked only
	// string params and discovery skipped methods (fn.Recv != nil); taint over the
	// field + method discovery must catch it.
	structMethod := `package fixture
import (
	"os"
	"strings"
)
type fixt struct { path string }
func (f *fixt) read(t *T) string {
	b, _ := os.ReadFile(f.path)
	return string(b)
}
func TestStructMethodOffender(t *T) {
	f := &fixt{path: "skills/first-officer/references/first-officer-shared-core.md"}
	s := f.read(t)
	if strings.Contains(s, "x") { _ = s }
}
`
	assertRedThenGreen(t, "struct-field + method-receiver path flow", "TestStructMethodOffender", structMethod)
}

// assertRedThenGreen plants fixtureSrc in a fresh temp dir, runs the sweep, and
// requires offenderName to be flagged (RED on the evasion shape). It then rewrites
// the offending test with a markNonAC declaration and requires the offender to
// clear (GREEN once declared). The declared rewrite reuses the fixture verbatim
// with the marker inserted as the test body's first statement.
func assertRedThenGreen(t *testing.T, shape, offenderName, fixtureSrc string) {
	t.Helper()
	dir := t.TempDir()
	writeFile(t, dir+"/evasion_test.go", fixtureSrc)
	offenders := sweepUndeclaredTautologies(t, dir)
	if !containsStr(offenders, offenderName) {
		t.Fatalf("%s: sweep failed to flag the evasion offender %s; offenders=%v", shape, offenderName, offenders)
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
	writeFile(t, dir2+"/evasion_test.go", declaredSrc)
	offenders = sweepUndeclaredTautologies(t, dir2)
	if containsStr(offenders, offenderName) {
		t.Fatalf("%s: sweep still flagged %s after it declared markNonAC; offenders=%v", shape, offenderName, offenders)
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
