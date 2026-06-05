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

// TestNoUndeclaredHostneutralityTautology is the AC-3 sweep for this package,
// re-runnable offline. It parses every *_test.go and flags any test function that
// READS a markdown INSTRUCTION file's content — by any flow: a named reader helper,
// a tainted os.ReadFile/os.Open/io read, or a WalkDir-collected `.md` — UNLESS it
// self-classifies via markNonAC or markCodeBoundInvariant. The go/parser code-scan
// invariants (host_neutrality_test.go's scanFile over .go source via parser.ParseFile,
// NOT a content read sink) and the spanHostQualified unit test are NOT flagged: they
// read no instruction file. The undeclared-offender count is the AC-3 metric; it
// must be zero.
//
// The sweep keys on the READ, not on how the bytes are inspected: any inspection
// idiom (strings.Contains/Index/EqualFold, bytes.*, regexp.Regexp.Match, a bare ==)
// is covered because the trigger is the ingest itself. A legitimate content-read (a
// re-bound invariant, a text-hygiene lint) declares markCodeBoundInvariant/markNonAC.
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

	// First pass: discover this package's instruction-file reader helpers, then grow
	// the set to a fixpoint so a read cannot hide behind a helper chain. Seeded with
	// the named readers; a func is ALSO a reader if it ingests instruction content
	// directly (readsInstructionContent — a tainted ReadFile/Open/io read, or a
	// WalkDir-collected `.md`) OR (transitive) it calls a known reader. Methods are
	// NOT skipped: a reader can be a method on a fixture struct (the s.path /
	// method-receiver flow). The code-scan helper scanFile is NOT a reader: it uses
	// parser.ParseFile (not a content read sink) over a `../dispatch` path (no
	// instruction taint).
	taintedFields := instructionTaintedFields(files)
	readers := map[string]bool{}
	for r := range instructionFileReaders {
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
			readsInstruction := readsInstructionContent(fn, taintedFields)
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

// readsInstructionContent reports whether fn ingests an instruction file's content
// by any flow — the positive/taint replacement for the Cycle-1/2 allow-lists. It
// taints every string derived from an instruction-file path (a `.md` skill-tree
// literal/segment, an instructionPathIdent package var, a param, a package-wide
// struct field, a local built via +/strings.Join/filepath.Join/fmt.Sprintf) and
// reports a read when a tainted path flows into any read sink (ReadFile/Open/
// ReadAll/bufio), or when fn WalkDir/Walks a tree collecting instruction `.md`.
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
		if (sel.Sel.Name == "WalkDir" || sel.Sel.Name == "Walk") && fnFiltersInstructionMarkdown(fn) {
			found = true
		}
		return true
	})
	return found
}

// readsTaintedField reports whether expr reads a struct field whose name is in the
// package-wide instruction-tainted-field set — the s.path / method-receiver flow.
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

// instructionTaintedFields scans every struct composite literal and field
// assignment across the package, returning the set of FIELD NAMES ever assigned an
// instruction-file path. Keyed by field name (no type info at parse time) — an
// over-approximation that errs toward flagging, which the proof policy wants.
func instructionTaintedFields(files []*ast.File) map[string]bool {
	fields := map[string]bool{}
	for _, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.KeyValueExpr:
				if key, ok := node.Key.(*ast.Ident); ok && exprInstructionTainted(node.Value, nil) {
					fields[key.Name] = true
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

// instructionTaintedNames computes the set of names (params, locals, recv.field
// selectors) in fn holding a string derived from an instruction-file path. Every
// string parameter is tainted (a helper that reads a string param is a path-arg
// reader — the caller supplies the .md path). It propagates through := / = to a
// fixpoint, including a local assigned from a package-wide tainted field.
func instructionTaintedNames(fn *ast.FuncDecl, taintedFields map[string]bool) map[string]bool {
	tainted := map[string]bool{}
	if fn.Type.Params != nil {
		for _, field := range fn.Type.Params.List {
			if id, ok := field.Type.(*ast.Ident); ok && id.Name == "string" {
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

// lvalueName renders an assignable target as a taint key: a bare ident or a
// selector `recv.field`.
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

// exprInstructionTainted reports whether expr carries an instruction-file path
// taint: a tainted name, an instruction-path literal/segment, a known
// instructionPathIdent package var, or a build of any of these via the path-build
// idioms (+/strings.Join/filepath.Join/fmt.Sprintf/string(...)).
func exprInstructionTainted(expr ast.Expr, tainted map[string]bool) bool {
	hit := false
	ast.Inspect(expr, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.BasicLit:
			if x.Kind == token.STRING && isInstructionPathLiteral(strings.Trim(x.Value, "`\"")) {
				hit = true
			}
		case *ast.Ident:
			if tainted[x.Name] || instructionPathIdents[x.Name] {
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

// instructionPathSegments are skill-tree / contract path segments that mark a path
// literal as an instruction file. The positive instruction-surface predicate (a
// path with any segment is instruction however it is built); replaces the
// `.md`-suffix-only detection that a Join/split-built suffix evaded.
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

	// Shape 5 (Cycle-3 M1, match axis) — the ingested bytes are inspected with a
	// match idiom OUTSIDE the old matchFuncs allowlist. The positive rule keys on the
	// READ (readSkill), so HOW the bytes are inspected is irrelevant.
	matchIndex := `package fixture
import "strings"
func TestMatchIndexHN(t *T) {
	text := readSkill(t, foCorePath)
	if strings.Index(text, "HALT") < 0 { t.Error("x") }
}
`
	assertRedThenGreenHN(t, "match via strings.Index (no Contains)", "TestMatchIndexHN", matchIndex)

	matchBytesRegexp := `package fixture
import "regexp"
func TestMatchBytesRegexpHN(t *T) {
	text := readSkill(t, foCorePath)
	re := regexp.MustCompile("HALT")
	if !re.Match([]byte(text)) { t.Error("x") }
}
`
	assertRedThenGreenHN(t, "match via regexp.Regexp.Match([]byte)", "TestMatchBytesRegexpHN", matchBytesRegexp)

	// Shape 6 (Cycle-3 M2, reader axis) — the `.md` path is built with strings.Join,
	// not `+`. The base fragment carries a skill-tree segment so it taints before the
	// suffix is appended.
	joinPath := `package fixture
import (
	"os"
	"strings"
)
func TestJoinPathHN(t *T) {
	base := "../../skills/first-officer/references/first-officer-shared-core"
	p := strings.Join([]string{base, "md"}, ".")
	b, _ := os.ReadFile(p)
	if strings.Index(string(b), "x") < 0 { t.Error("y") }
}
`
	assertRedThenGreenHN(t, "strings.Join-built .md path", "TestJoinPathHN", joinPath)

	// Shape 7 (Cycle-3 M3, reader axis) — the `.md` path flows through a struct field
	// and the read happens via a METHOD on that struct; discovery must include methods
	// and taint the package-wide field.
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
func TestStructMethodHN(t *T) {
	f := &fixt{path: "../../skills/first-officer/references/first-officer-shared-core.md"}
	s := f.read(t)
	if strings.Contains(s, "x") { _ = s }
}
`
	assertRedThenGreenHN(t, "struct-field + method-receiver path flow", "TestStructMethodHN", structMethod)
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
