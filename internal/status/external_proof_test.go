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
		for _, f := range ClassifyEntityACs(stripFrontmatter(data)) {
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
