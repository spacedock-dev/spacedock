// ABOUTME: Behavioral proof for the edge-advance job's `decision` step — runs
// ABOUTME: the REAL release.yml shell (not a text-shape assertion) against a
// ABOUTME: fixture git repo, closing the two falsifiers the validator found
// ABOUTME: dead: a stable-only-narrowed candidate pool, and a fail-open empty pool.
package release

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/testgit"
)

// edgeAdvanceDecisionStep returns the edge-advance job's `id: decision` step,
// or nil when none — the step whose run script this file executes for real.
func edgeAdvanceDecisionStep(job workflowJob) *workflowStep {
	for i := range job.steps {
		if job.steps[i].id == "decision" {
			return &job.steps[i]
		}
	}
	return nil
}

// tagFixtureRepo builds a throwaway git repo with one commit and the given
// tags (lightweight — the decision script only ever reads `git tag --list`,
// never a tag's own annotation), so its `.git` can be swapped in for the
// running shell via GIT_DIR/GIT_WORK_TREE without touching this checkout's own
// history.
func tagFixtureRepo(t *testing.T, tags []string) string {
	t.Helper()
	dir := t.TempDir()
	testgit.InitRepo(t, dir, "-q", "-b", "main")
	if err := os.WriteFile(filepath.Join(dir, "seed.txt"), []byte("seed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	git("add", "-A")
	git("commit", "-q", "-m", "seed")
	for _, tag := range tags {
		git("tag", tag)
	}
	return dir
}

// runDecisionStepScript executes the ACTUAL `decision` step's run script
// extracted from the on-disk release.yml — not a copy, not a re-derivation —
// against a fixture repo's tags. cmd.Dir stays this checkout's repo root (the
// script's `go run ./cmd/spacedock-release ...` calls resolve relative to it),
// while GIT_DIR/GIT_WORK_TREE redirect every `git` call the script makes to
// the fixture repo instead. It returns the $GITHUB_OUTPUT file's contents and
// the script's own stdout/stderr (for `::notice::` diagnostics on failure).
func runDecisionStepScript(t *testing.T, refName, fixtureRepo string) (githubOutput, scriptOutput string) {
	t.Helper()
	workflow := readWorkflow(t, "release.yml")
	job := edgeAdvanceJob(workflow)
	if job == nil {
		t.Fatal("release.yml has no edge-advance job")
	}
	step := edgeAdvanceDecisionStep(*job)
	if step == nil {
		t.Fatal("edge-advance has no step with id: decision")
	}

	repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	outFile := filepath.Join(t.TempDir(), "github_output")
	if err := os.WriteFile(outFile, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("bash", "-c", step.run)
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(),
		"GIT_DIR="+filepath.Join(fixtureRepo, ".git"),
		"GIT_WORK_TREE="+fixtureRepo,
		"GITHUB_REF_NAME="+refName,
		"GITHUB_OUTPUT="+outFile,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("real release.yml decision step script failed: %v\n%s", err, out)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	return string(data), string(out)
}

// TestDecisionStepScriptIncludesPrereleaseCandidates makes the first of the
// two falsifiers the validator found dead actually fire: "fails if the
// candidate pool silently narrows to stable-only." The fixture's only bare
// stable is v0.25.0; the newer 0.27 line exists ONLY as a prerelease
// (v0.27.0-pre5), no bare v0.27.0 tag at all. Deciding v0.26.0 (a stable cut
// slotting between them) must SKIP — 0.27.0-pre5 already outranks v0.26.0's
// own target (0.27.0-pre1) by construction (same core, pre1 < pre5) — which is
// only possible if the script's `git tag --list 'v*'` candidate scan actually
// carries prereleases through to highest-known-edge-version. A candidate-pool
// edit that narrows the scan to bare stables only (e.g. a stricter tag glob or
// an inline grep) would find only v0.25.0, wrongly rank v0.26.0 as the latest
// line, and print advance=true here instead — this test would catch that
// regression the moment it landed in release.yml, not just in a Go unit test
// of the library function called with a hand-picked candidate list.
func TestDecisionStepScriptIncludesPrereleaseCandidates(t *testing.T) {
	fixture := tagFixtureRepo(t, []string{"v0.25.0", "v0.27.0-pre5", "v0.26.0"})
	output, log := runDecisionStepScript(t, "v0.26.0", fixture)
	if !strings.Contains(output, "advance=false") {
		t.Fatalf("decision step wrote %q, want advance=false (a stable-only-narrowed candidate pool would print advance=true here)\nscript output:\n%s", output, log)
	}
	if strings.Contains(output, "advance=true") {
		t.Fatalf("decision step wrote both advance=true and advance=false: %q", output)
	}
}

// TestDecisionStepScriptFailsClosedOnEmptyCandidatePool makes the second dead
// falsifier fire: "fails if an empty pool stops failing closed." The fixture
// repo's only tag is the one being decided, so the candidate pool (every
// OTHER tag) is empty after the script's own `grep -v "^$GITHUB_REF_NAME\$"`
// exclusion — highest-known-edge-version prints nothing, and the script's `[
// -z "$HIGHEST_KNOWN" ]` branch must write advance=false, never fall open to
// "nothing to compare against, so anything advances" (a wrongly-cut lower pre0
// publishes a regression; a missed one is recoverable by hand). A regression
// back to the first cut's `HIGHEST_OTHER="0.0.0"`-style fail-open fallback
// would make v99.9.9 compare against 0.0.0 and print advance=true here.
func TestDecisionStepScriptFailsClosedOnEmptyCandidatePool(t *testing.T) {
	fixture := tagFixtureRepo(t, []string{"v0.25.0"})
	output, log := runDecisionStepScript(t, "v0.25.0", fixture)
	if !strings.Contains(output, "advance=false") {
		t.Fatalf("decision step wrote %q, want advance=false (empty candidate pool must fail CLOSED)\nscript output:\n%s", output, log)
	}
	if strings.Contains(output, "advance=true") {
		t.Fatalf("decision step wrote both advance=true and advance=false: %q", output)
	}
}
