// ABOUTME: External-proof guard tests — opt-in default-off, shared-classifier
// ABOUTME: invariant, and the live-corpus precision/recall = 1.0 invariant.
package status

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestExternalProofOptInDefaultOff locks AC-3: under `require-external-proof:
// false` AND under the absent default, neither the terminal-set guard nor the
// --validate sub-check fires, and ordinary read flows are never gated.
func TestExternalProofOptInDefaultOff(t *testing.T) {
	cases := []struct {
		name        string
		readmeOptIn string // the require-external-proof line, or "" for absent
	}{
		{"absent", ""},
		{"explicit-false", "require-external-proof: false\n"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			env := pinnedEnv(t)
			root := stageExternalProofFixture(t, tc.readmeOptIn)

			// Terminal --set on the self-ref entity must succeed: the guard is silent.
			args := []string{"--workflow-dir", root, "--set", "010-self-ref-only", "status=done"}
			_, nErr, nCode := runNative(t, root, env, args...)
			if nCode != 0 {
				t.Fatalf("--set should exit 0 when guard is off, got %d (%q)", nCode, nErr)
			}
			fm := readFrontmatter(t, filepath.Join(root, "010-self-ref-only.md"))
			if !strings.Contains(fm, "status: done") {
				t.Fatalf("self-ref entity should have advanced to done, fm=%s", fm)
			}

			// --validate must return VALID when guard is off.
			vOut, _, vCode := runNative(t, root, env, "--workflow-dir", root, "--validate")
			if vCode != 0 || !strings.Contains(vOut, "VALID") {
				t.Fatalf("--validate should return VALID exit 0, got code=%d stdout=%q", vCode, vOut)
			}

			// Default table read must exit 0; --next and --boot must also exit 0.
			_, _, rCode := runNative(t, root, env, "--workflow-dir", root)
			if rCode != 0 {
				t.Fatalf("default table read should exit 0 when guard is off, got %d", rCode)
			}
			_, _, nxCode := runNative(t, root, env, "--workflow-dir", root, "--next")
			if nxCode != 0 {
				t.Fatalf("--next should exit 0 when guard is off, got %d", nxCode)
			}
			_, _, bCode := runNative(t, root, env, "--workflow-dir", root, "--boot")
			if bCode != 0 {
				t.Fatalf("--boot should exit 0 when guard is off, got %d", bCode)
			}
		})
	}
}

// TestClassifierIsSharedBySetAndValidate locks AC-4: a single shared classifier
// serves both the terminal-set guard and the --validate sub-check. Verified by
// (1) a structural single-definition assertion (`grep -c "^func ClassifyEntityACs"`
// over internal/status/*.go returns 1) and (2) a runtime call-counter that
// both surfaces bump in one test.
func TestClassifierIsSharedBySetAndValidate(t *testing.T) {
	// Structural invariant: exactly one definition site of ClassifyEntityACs
	// in internal/status/*.go. Reads the real parsed file content.
	defRe := regexp.MustCompile(`(?m)^func ClassifyEntityACs\b`)
	matches := 0
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read . : %v", err)
	}
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		matches += len(defRe.FindAllIndex(data, -1))
	}
	if matches != 1 {
		t.Fatalf("ClassifyEntityACs must be defined exactly once across internal/status/*.go, got %d definitions", matches)
	}

	// Runtime invariant: both runSet and validateWorkflow exercise the same
	// classifier — verified by a single test that runs both surfaces against
	// the same fixture and asserts the shared call counter advanced for each.
	env := pinnedEnv(t)
	root := stageExternalProofFixture(t, "require-external-proof: true\n")

	before := classifierCallCount
	// --validate run touches the classifier per active entity. Use the
	// counter delta rather than its absolute value (other tests may have run
	// first and bumped it).
	_, _, _ = runNative(t, root, env, "--workflow-dir", root, "--validate")
	afterValidate := classifierCallCount
	if afterValidate <= before {
		t.Fatalf("--validate did not exercise the classifier (counter %d → %d)", before, afterValidate)
	}
	// --set terminal advance triggers the runSet guard's classifier call on
	// the targeted entity. Use the real-proof entity so the guard exits 0; we
	// only care that the classifier ran.
	_, _, _ = runNative(t, root, env, "--workflow-dir", root, "--set", "020-real-proof", "status=done")
	afterSet := classifierCallCount
	if afterSet <= afterValidate {
		t.Fatalf("--set did not exercise the classifier (counter %d → %d)", afterValidate, afterSet)
	}
}

// TestClassifierPrecisionRecallOnLiveCorpus locks AC-5: walk every index.md
// under docs/dev/.spacedock-state/ (active + _archive/) and assert the flagged
// set is EXACTLY {external-tracker-checkpoint/index.md AC-6}.
func TestClassifierPrecisionRecallOnLiveCorpus(t *testing.T) {
	stateRoot := liveStateRoot(t)
	if stateRoot == "" {
		t.Skip("live .spacedock-state corpus not reachable from internal/status test cwd")
	}

	type hit struct {
		path  string
		label string
	}
	var hits []hit

	err := filepath.Walk(stateRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Base(path) != "index.md" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		for _, f := range ClassifyEntityACs(stripFrontmatter(data), devExternalTokens) {
			rel, _ := filepath.Rel(stateRoot, path)
			hits = append(hits, hit{path: rel, label: acLabel(f.Header)})
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", stateRoot, err)
	}

	const wantPath = "_archive/external-tracker-checkpoint/index.md"
	const wantLabel = "AC-6"
	if len(hits) != 1 {
		var lines []string
		for _, h := range hits {
			lines = append(lines, h.path+" "+h.label)
		}
		t.Fatalf("classifier must flag EXACTLY one AC on the live corpus, got %d:\n%s",
			len(hits), strings.Join(lines, "\n"))
	}
	if hits[0].path != wantPath || hits[0].label != wantLabel {
		t.Fatalf("flagged set mismatch: got {%s %s}, want {%s %s}",
			hits[0].path, hits[0].label, wantPath, wantLabel)
	}
}

// TestNoExternalProofGuardOnReadPaths locks the cycle-2 F1 regression: under
// `require-external-proof: true` AND with a self-referential entity present,
// the read surfaces (`sd status`, `--next`, `--boot`, `--next-id`) must exit 0.
// The classifier gate fires only on the mutation surface (runSet terminal-set)
// and the explicit `--validate` command; a read path failing on a flagged AC
// would lock the FO out of the very listing they need to see the broken entity.
func TestNoExternalProofGuardOnReadPaths(t *testing.T) {
	env := pinnedEnv(t)
	root := stageExternalProofFixture(t, "require-external-proof: true\n")

	for _, args := range [][]string{
		{"--workflow-dir", root},
		{"--workflow-dir", root, "--next"},
		{"--workflow-dir", root, "--boot"},
		{"--workflow-dir", root, "--next-id"},
	} {
		name := strings.Join(args[2:], " ")
		if name == "" {
			name = "default-table"
		}
		t.Run(name, func(t *testing.T) {
			nOut, nErr, nCode := runNative(t, root, env, args...)
			if nCode != 0 {
				t.Fatalf("read path %q must exit 0 under require-external-proof: true, got %d\nstdout=%q\nstderr=%q",
					args, nCode, nOut, nErr)
			}
			if strings.Contains(nErr, "self-referential AC proof") {
				t.Fatalf("read path %q must not emit the self-referential guard error, stderr=%q", args, nErr)
			}
		})
	}

	// --validate, on the other hand, MUST still fire — it's the explicit
	// surface the gate keys on.
	_, vErr, vCode := runNative(t, root, env, "--workflow-dir", root, "--validate")
	if vCode != 1 || !strings.Contains(vErr, "self-referential AC proof") {
		t.Fatalf("--validate must opt INTO the external-proof check, got code=%d stderr=%q", vCode, vErr)
	}
}

// TestExternalTokensClearSelfPhrase locks the cycle-2 L3-F1/L3-F2/L3-F3/L3-F4
// vocabulary additions: a self-phrased AC whose proof clause cites any of the
// build/compile/lint, live-pilot (no literal-article dependency), release-
// artifact, or commit/URL/GitHub families must NOT flag. Each subtest is a
// minimal classifier-level body fixture proving the relevant token clears the
// self-phrase.
func TestExternalTokensClearSelfPhrase(t *testing.T) {
	cases := []struct {
		name   string
		clause string
	}{
		// L3-F1: build/compile/vet/gofmt/lint.
		{"go-build", "Verified by: `go build ./...` succeeds against this entity's renamed files."},
		{"gofmt", "Verified by: `gofmt -l .` returns empty on this entity's tree."},
		{"go-vet", "Verified by: `go vet ./...` succeeds on this entity's package."},
		{"compile", "Verified by: a downstream caller compiles cleanly against this entity's new API."},
		{"lint", "Verified by: a lint pass clears this entity's renamed identifiers."},

		// L3-F2: generalized live-pilot vocabulary (no literal-article dependency).
		{"drives-lifecycle", "Verified by: a behavioral pilot on a real entity drives this entity's lifecycle."},
		{"pilots-merge", "Verified by: the FO pilots this entity through merge."},
		{"runtime-behavior", "Verified by: runtime behavior of the registry drives this entity's mod-block clear."},

		// L3-F3: release-artifact vocabulary.
		{"brew", "Verified by: a `brew install` against this entity's cask succeeds."},
		{"spctl", "Verified by: `spctl --assess --type execute` clears this entity's binary."},
		{"goreleaser", "Verified by: a `goreleaser` run produces this entity's artifact."},
		{"notarize", "Verified by: notarize a binary derived from this entity's commit."},

		// L3-F4: commit/URL/GitHub/merged vocabulary.
		{"commit-url", "Verified by: opening this entity's commit URL on GitHub renders the diff."},
		{"merged", "Verified by: this entity's PR has merged to main."},
		{"landed", "Verified by: the change landed on main in this entity's commit."},
		{"hex-shape", "Verified by: this entity's resulting commit deadbeef carries the rename."},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			body := "**AC-1 — A property of this entity's design.**\n" + tc.clause + "\n"
			flags := ClassifyEntityACs(body, devExternalTokens)
			if len(flags) != 0 {
				t.Fatalf("clause %q should clear the self-phrase, got flags=%+v", tc.clause, flags)
			}
		})
	}
}

// TestSelfPhraseStillFlagsWithoutExternalToken locks that the vocabulary
// additions did not over-correct: a truly self-referential AC (self-phrase
// present, no external token) still flags. This is the recall side of the
// precision/recall = 1.0 invariant.
func TestSelfPhraseStillFlagsWithoutExternalToken(t *testing.T) {
	cases := []string{
		"**AC-1 — Intent.**\nVerified by: review of this entity's decision section.\n",
		"**AC-2 — Intent.**\nVerified by: the entity's own decision states the intent.\n",
		"**AC-3 — Intent.**\nVerified by: re-reading this entity's body confirms the intent.\n",
	}
	for i, body := range cases {
		body := body
		t.Run(acLabel(strings.Split(body, "\n")[0]), func(t *testing.T) {
			flags := ClassifyEntityACs(body, devExternalTokens)
			if len(flags) != 1 {
				t.Fatalf("case %d should flag exactly once, got %+v", i, flags)
			}
		})
	}
}

// TestExternalTokenInBacktickStillClears locks the cycle-2 L3-F3 strip-order
// fix: a self-phrased AC whose external token sits inside backticks must
// still clear. The fix runs externalTokenRe against the UNSTRIPPED clause so
// `--cask` inside backticks (the corpus convention) is honored.
func TestExternalTokenInBacktickStillClears(t *testing.T) {
	body := "**AC-1 — A property of this entity's design.**\n" +
		"Verified by: `--cask` against this entity's binary.\n"
	flags := ClassifyEntityACs(body, devExternalTokens)
	if len(flags) != 0 {
		t.Fatalf("backtick-quoted external token must clear the self-phrase, got %+v", flags)
	}
}

// TestQuotedSelfPhraseStillDoesNotFlag locks that the strip-order change did
// not regress the original quote-stripping intent: a clause that QUOTES the
// self-phrase as an example (not as the live proof) must still NOT flag.
func TestQuotedSelfPhraseStillDoesNotFlag(t *testing.T) {
	body := "**AC-1 — A property.**\n" +
		"Verified by: a Go test refusing \"this entity's decision section\" as proof.\n"
	flags := ClassifyEntityACs(body, devExternalTokens)
	if len(flags) != 0 {
		t.Fatalf("quoted self-phrase + external token must not flag, got %+v", flags)
	}
}

// TestExternalProofPolicyEmptyAfterColonIsOff locks the cycle-2 F2 polish:
// `require-external-proof:` (empty after colon, including comment-only and
// null) coerces to OFF, byte-identical to absent. The error message
// enumerates these as absent-equivalent shapes so a reader who hits the typo
// guard understands them.
func TestExternalProofPolicyEmptyAfterColonIsOff(t *testing.T) {
	cases := []struct {
		name        string
		readmeOptIn string
	}{
		{"empty-after-colon", "require-external-proof:\n"},
		{"comment-only", "require-external-proof: # off for now\n"},
		{"explicit-null", "require-external-proof: null\n"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			readme := "---\nid-style: sequential\n" + tc.readmeOptIn + "---\n\n# Fixture\n"
			if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(readme), 0o644); err != nil {
				t.Fatal(err)
			}
			policy, err := resolveExternalProofPolicy(dir)
			if err != nil {
				t.Fatalf("empty/comment/null after colon should not error, got %v", err)
			}
			if policy != externalProofOff {
				t.Fatalf("empty/comment/null after colon should resolve to OFF, got %v", policy)
			}
		})
	}
}

// TestExternalProofPolicyTypoErrorEnumeratesShapes locks the cycle-2 F2
// polish: the typo-rejection error message includes the absent-equivalent
// shapes so a reader hitting the guard knows `key:` / `key: null` are valid.
func TestExternalProofPolicyTypoErrorEnumeratesShapes(t *testing.T) {
	dir := t.TempDir()
	readme := "---\nid-style: sequential\nrequire-external-proof: tru\n---\n\n# Fixture\n"
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := resolveExternalProofPolicy(dir)
	if err == nil {
		t.Fatal("typo should reject loudly")
	}
	msg := err.Error()
	for _, want := range []string{"'tru'", "absent", "empty", "null"} {
		if !strings.Contains(msg, want) {
			t.Fatalf("error %q should mention %q", msg, want)
		}
	}
}

// TestClassifierKernelTakesInjectedTokens locks AC-6(a): the generic kernel
// `ClassifyEntityACs` takes the external-token vocabulary as an injected
// parameter — it is callable with a NON-dev token set and honors it, proving
// the kernel never reaches for a package-global dev vocabulary. A separate
// assertion drives the kernel with `devExternalTokens` over the same fixtures
// TestExternalTokensClearSelfPhrase uses and gets the identical verdicts, so
// the dev behavior is unchanged.
func TestClassifierKernelTakesInjectedTokens(t *testing.T) {
	// A self-phrased clause whose only external token is a NON-dev word
	// ("p-value") that the dev vocabulary does NOT recognize.
	body := "**AC-1 — A property of this entity's design.**\n" +
		"Verified by: the reported p-value falls below this entity's pre-registered threshold.\n"

	// Under the dev vocabulary the clause has no recognized external token, so
	// the kernel flags it — proving the kernel applies the INJECTED set.
	devFlags := ClassifyEntityACs(body, devExternalTokens)
	if len(devFlags) != 1 {
		t.Fatalf("under dev vocabulary the non-dev clause should flag (no dev token present), got %+v", devFlags)
	}

	// Inject a non-dev vocabulary that DOES recognize "p-value"; the same clause
	// must now clear — the kernel honors whatever matcher the caller injects.
	researchTokens := regexp.MustCompile(`(?i)p-value|dataset|\bDOI\b`)
	researchFlags := ClassifyEntityACs(body, researchTokens)
	if len(researchFlags) != 0 {
		t.Fatalf("with an injected research vocabulary recognizing p-value the clause should clear, got %+v", researchFlags)
	}

	// Dev-behavior preservation: the kernel driven with devExternalTokens over
	// TestExternalTokensClearSelfPhrase's fixtures yields the identical clear
	// verdicts (no flags).
	devFixtures := []string{
		"Verified by: `go build ./...` succeeds against this entity's renamed files.",
		"Verified by: `gofmt -l .` returns empty on this entity's tree.",
		"Verified by: a `goreleaser` run produces this entity's artifact.",
		"Verified by: this entity's PR has merged to main.",
	}
	for _, clause := range devFixtures {
		fixtureBody := "**AC-1 — A property of this entity's design.**\n" + clause + "\n"
		if flags := ClassifyEntityACs(fixtureBody, devExternalTokens); len(flags) != 0 {
			t.Fatalf("dev clause %q should clear under devExternalTokens, got %+v", clause, flags)
		}
	}
}

// TestDevPathSuppliesDevVocabulary locks AC-6(a)'s supplier half: the dev
// callers reach the kernel through classifyEntityFile, which must supply
// devExternalTokens — verified by driving classifyEntityFile over a fixture
// whose dev token clears the self-phrase. If classifyEntityFile failed to
// inject the dev vocabulary, the dev clause would flag.
func TestDevPathSuppliesDevVocabulary(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "index.md")
	body := "---\nid: x\n---\n\n" +
		"**AC-1 — A property of this entity's design.**\n" +
		"Verified by: `go vet ./...` succeeds on this entity's package.\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if flags := classifyEntityFile(path); len(flags) != 0 {
		t.Fatalf("classifyEntityFile must supply devExternalTokens so the go-vet clause clears, got %+v", flags)
	}
}

// kernelSpanNames are the source spans that make up the generic, workflow-
// agnostic kernel: the classifier function, its generic helpers, and the
// self-reference phrase table. None of these may name a dev-toolchain token.
var kernelSpanNames = []string{
	"func ClassifyEntityACs",
	"func isolateProofClause",
	"func matchSelfPhrase",
	"var selfPhraseRes",
}

// TestKernelFreeOfDevVocabulary locks AC-6(b): the generic kernel's source
// spans contain NO Go-toolchain literal. The dev vocabulary must live only in
// the devExternalTokens var. Negative proof: re-hard-code a dev token inside
// the kernel and this test goes red.
func TestKernelFreeOfDevVocabulary(t *testing.T) {
	data, err := os.ReadFile("external_proof.go")
	if err != nil {
		t.Fatalf("read external_proof.go: %v", err)
	}
	src := string(data)

	bannedDevLiterals := []string{`gofmt`, `goreleaser`, `\bvet\b`, `\.go\b`, `\bcask\b`}

	for _, name := range kernelSpanNames {
		span := topLevelSpan(t, src, name)
		for _, banned := range bannedDevLiterals {
			if strings.Contains(span, banned) {
				t.Errorf("kernel span %q contains dev-toolchain literal %q — dev vocabulary leaked back into the generic kernel", name, banned)
			}
		}
	}
}

// topLevelSpan returns the source of the top-level declaration that begins with
// decl (e.g. "func ClassifyEntityACs" or "var selfPhraseRes"), from that line
// to the matching closing brace at column 0. Top-level Go declarations close on
// a `}` (or `)`) in the first column, so the span ends at the first such line.
func topLevelSpan(t *testing.T, src, decl string) string {
	t.Helper()
	lines := strings.Split(src, "\n")
	start := -1
	for i, line := range lines {
		if strings.HasPrefix(line, decl) {
			start = i
			break
		}
	}
	if start < 0 {
		t.Fatalf("declaration %q not found in external_proof.go", decl)
	}
	for i := start + 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "}") || strings.HasPrefix(lines[i], ")") {
			return strings.Join(lines[start:i+1], "\n")
		}
	}
	t.Fatalf("no top-level close found for %q", decl)
	return ""
}

// stageExternalProofFixture copies testdata/external-proof-workflow into a
// fresh git-initialized temp dir, then rewrites the README's
// `require-external-proof:` line to optIn (which is "" for the absent case or
// the full line including the newline). Returns the absolute root.
func stageExternalProofFixture(t *testing.T, optIn string) string {
	t.Helper()
	src, err := filepath.Abs(filepath.Join("testdata", "external-proof-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	dst := t.TempDir()
	cpTree(t, src, dst)

	// Rewrite README to replace the fixture's `require-external-proof: true\n`
	// with the requested optIn. An empty optIn omits the key entirely.
	readme := filepath.Join(dst, "README.md")
	data, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	replaced := bytes.Replace(data, []byte("require-external-proof: true\n"), []byte(optIn), 1)
	if err := os.WriteFile(readme, replaced, 0o644); err != nil {
		t.Fatal(err)
	}

	gitInit(t, dst)
	return dst
}

// liveStateRoot resolves the live .spacedock-state checkout absolute path
// relative to the test cwd (internal/status), walking up to the repo root.
// Returns "" when the corpus is not present (e.g. in a packaging-time build).
func liveStateRoot(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	dir := cwd
	for i := 0; i < 6; i++ {
		candidate := filepath.Join(dir, "docs", "dev", ".spacedock-state")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}
