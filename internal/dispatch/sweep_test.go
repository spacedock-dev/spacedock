// ABOUTME: Sweep reuses reconcile's un-advanced-PR detection — the swept set is the
// ABOUTME: merged-but-not-terminalized entities, computed from an injected gh stub.
package dispatch

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSweepReportsMergedNotTerminalized pins AC-7's reported-set property: Sweep
// returns exactly the entities whose PR is MERGED and whose status is not done,
// reusing classC's detection. The gh stub pins a deterministic merged-state so the
// test is offline. pr-merged (merged, status!=done) IS swept; pr-open (OPEN) and
// pr-merged-done (status=done) are NOT — the two negative conjuncts.
func TestSweepReportsMergedNotTerminalized(t *testing.T) {
	repoRoot := t.TempDir()
	workflowDir := filepath.Join(repoRoot, "docs", "wf")
	writeFile(t, filepath.Join(workflowDir, "README.md"), `---
entity-type: task
state: .spacedock-state
stages:
  states:
    - name: ideation
      initial: true
    - name: done
      terminal: true
---
`)
	stateRoot := filepath.Join(workflowDir, ".spacedock-state")

	writeFile(t, filepath.Join(stateRoot, "pr-merged", "index.md"), reconcileEntityFM(map[string]string{
		"id": "id-m", "title": "merged", "slug": "pr-merged", "status": "implementation", "pr": "42",
	}))
	writeFile(t, filepath.Join(stateRoot, "pr-open", "index.md"), reconcileEntityFM(map[string]string{
		"id": "id-o", "title": "open", "slug": "pr-open", "status": "implementation", "pr": "43",
	}))
	writeFile(t, filepath.Join(stateRoot, "pr-merged-done", "index.md"), reconcileEntityFM(map[string]string{
		"id": "id-md", "title": "merged done", "slug": "pr-merged-done", "status": "done", "pr": "44",
	}))

	gh := func(pr string) (string, error) {
		switch pr {
		case "42", "44":
			return "MERGED", nil
		default:
			return "OPEN", nil
		}
	}

	var out, errBuf strings.Builder
	code := Sweep(workflowDir, gh, true, &out, &errBuf)
	if code != 0 {
		t.Fatalf("Sweep should exit 0; got %d stderr=%q", code, errBuf.String())
	}

	var res sweepResult
	if err := json.Unmarshal([]byte(out.String()), &res); err != nil {
		t.Fatalf("Sweep --json should emit valid JSON: %v\n%s", err, out.String())
	}
	gotSlugs := map[string]bool{}
	for _, s := range res.Swept {
		gotSlugs[s.Slug] = true
	}
	if !gotSlugs["pr-merged"] {
		t.Fatalf("sweep should report pr-merged (merged, status!=done); swept=%v", res.Swept)
	}
	if gotSlugs["pr-open"] {
		t.Fatalf("sweep must NOT report pr-open (PR not merged); swept=%v", res.Swept)
	}
	if gotSlugs["pr-merged-done"] {
		t.Fatalf("sweep must NOT report pr-merged-done (already terminalized); swept=%v", res.Swept)
	}
	if len(res.Swept) != 1 {
		t.Fatalf("sweep should report exactly the one merged-not-done entity; swept=%v", res.Swept)
	}
}

// TestSweepEmptyIsEmptyArray pins the empty case: no merged-not-done entities yields
// an empty swept array (never null), exit 0.
func TestSweepEmptyIsEmptyArray(t *testing.T) {
	repoRoot := t.TempDir()
	workflowDir := filepath.Join(repoRoot, "docs", "wf")
	writeFile(t, filepath.Join(workflowDir, "README.md"), `---
state: .spacedock-state
stages:
  states:
    - name: ideation
      initial: true
    - name: done
      terminal: true
---
`)
	stateRoot := filepath.Join(workflowDir, ".spacedock-state")
	if err := os.MkdirAll(stateRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(stateRoot, "no-pr", "index.md"), reconcileEntityFM(map[string]string{
		"id": "id-np", "title": "no pr", "slug": "no-pr", "status": "ideation",
	}))

	gh := func(string) (string, error) { return "OPEN", nil }
	var out, errBuf strings.Builder
	if code := Sweep(workflowDir, gh, true, &out, &errBuf); code != 0 {
		t.Fatalf("Sweep should exit 0; got %d stderr=%q", code, errBuf.String())
	}
	if !strings.Contains(out.String(), `"swept": []`) {
		t.Fatalf("empty sweep should emit an empty array, not null; json:\n%s", out.String())
	}
}
