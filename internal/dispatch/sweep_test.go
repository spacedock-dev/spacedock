// ABOUTME: Sweep reuses reconcile's un-advanced-PR detection — the swept set is the
// ABOUTME: merged-but-not-terminalized entities, computed from an injected gh stub.
package dispatch

import (
	"encoding/json"
	"fmt"
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

// TestSweepGhUnavailableReportsUnknown pins AC-4 (D2): with >=1 PR-pending entity
// and every gh probe erroring, Sweep must declare merge state UNKNOWN rather than
// silently reporting "0 entity(ies) merged" — a real empty sweep and a gh-outage
// sweep are indistinguishable without this. classC itself stays best-effort (it
// swallows the error); Sweep counts probes around it.
func TestSweepGhUnavailableReportsUnknown(t *testing.T) {
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
	writeFile(t, filepath.Join(stateRoot, "pr-pending", "index.md"), reconcileEntityFM(map[string]string{
		"id": "id-p", "title": "pending", "slug": "pr-pending", "status": "implementation", "pr": "42",
	}))

	ghUnavailable := func(string) (string, error) { return "", fmt.Errorf("gh: command not found") }
	var out, errBuf strings.Builder
	code := Sweep(workflowDir, ghUnavailable, true, &out, &errBuf)
	if code != 0 {
		t.Fatalf("Sweep should exit 0 even when gh is unavailable; got %d stderr=%q", code, errBuf.String())
	}
	var res sweepResult
	if err := json.Unmarshal([]byte(out.String()), &res); err != nil {
		t.Fatalf("Sweep --json should emit valid JSON: %v\n%s", err, out.String())
	}
	if res.Gh != "unavailable" {
		t.Fatalf(`sweep JSON should carry "gh": "unavailable", got %q (full: %s)`, res.Gh, out.String())
	}
	if !strings.Contains(res.Reason, "UNKNOWN") {
		t.Fatalf("reason should declare merge state UNKNOWN, got %q", res.Reason)
	}
	if strings.Contains(res.Reason, "0 entity(ies)") {
		t.Fatalf("gh-unavailable reason must NOT read as a real empty sweep, got %q", res.Reason)
	}
	if len(res.Swept) != 0 {
		t.Fatalf("gh-unavailable sweep should report zero swept entities (all probes errored), got %v", res.Swept)
	}
}

// TestSweepGhPartiallyAvailableStillReportsNormally pins the complement: when at
// least one gh probe SUCCEEDS, the sweep is NOT gh-unavailable — a mixed
// success/failure batch still resolves to the normal count-based reason (the
// binary is reachable; a single flake should not mask real merged entities).
func TestSweepGhPartiallyAvailableStillReportsNormally(t *testing.T) {
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
	writeFile(t, filepath.Join(stateRoot, "pr-erroring", "index.md"), reconcileEntityFM(map[string]string{
		"id": "id-e", "title": "erroring", "slug": "pr-erroring", "status": "implementation", "pr": "43",
	}))

	gh := func(pr string) (string, error) {
		if pr == "42" {
			return "MERGED", nil
		}
		return "", fmt.Errorf("rate limited")
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
	if res.Gh == "unavailable" {
		t.Fatalf("a partial gh failure must NOT report unavailable, got %q", res.Gh)
	}
	if !strings.Contains(res.Reason, "1 entity(ies)") {
		t.Fatalf("reason should still report the one real merged entity, got %q", res.Reason)
	}
}

// TestSweepNonEmptyNamesRegisteredStartupModNextStep pins AC-4/D2's non-empty
// next-step: a workflow with a registered startup-hook mod names the mod file to
// advance per, not a hardcoded procedure — the mod file is the per-workflow
// authority (shipped pr-merge.md advances directly; this repo's local one
// delegates to sentinel+merge-guard, and the binary must not pick a side).
func TestSweepNonEmptyNamesRegisteredStartupModNextStep(t *testing.T) {
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
	writeFile(t, filepath.Join(workflowDir, "_mods", "pr-merge.md"), "## Hook: startup\n\nAdvance a merged PR.\n")
	stateRoot := filepath.Join(workflowDir, ".spacedock-state")
	writeFile(t, filepath.Join(stateRoot, "pr-merged", "index.md"), reconcileEntityFM(map[string]string{
		"id": "id-m", "title": "merged", "slug": "pr-merged", "status": "implementation", "pr": "42",
	}))

	gh := func(string) (string, error) { return "MERGED", nil }
	var out, errBuf strings.Builder
	if code := Sweep(workflowDir, gh, true, &out, &errBuf); code != 0 {
		t.Fatalf("Sweep should exit 0; got %d stderr=%q", code, errBuf.String())
	}
	var res sweepResult
	if err := json.Unmarshal([]byte(out.String()), &res); err != nil {
		t.Fatalf("Sweep --json should emit valid JSON: %v\n%s", err, out.String())
	}
	if !strings.Contains(res.Next, "_mods/pr-merge.md") {
		t.Fatalf("next-step should name the registered startup mod file, got %q", res.Next)
	}
	if strings.Contains(res.Next, "merge guard") || strings.Contains(res.Next, "sentinel") {
		t.Fatalf("next-step must point at the mod file, not prescribe a procedure, got %q", res.Next)
	}
}

// TestSweepNonEmptyNamesGenericModPointerWhenNoneRegistered pins the fallback: a
// non-empty sweep with NO startup-hook mod registered still names a next-step,
// pointing generically at _mods/ rather than a specific (nonexistent) file.
func TestSweepNonEmptyNamesGenericModPointerWhenNoneRegistered(t *testing.T) {
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

	gh := func(string) (string, error) { return "MERGED", nil }
	var out, errBuf strings.Builder
	if code := Sweep(workflowDir, gh, true, &out, &errBuf); code != 0 {
		t.Fatalf("Sweep should exit 0; got %d stderr=%q", code, errBuf.String())
	}
	var res sweepResult
	if err := json.Unmarshal([]byte(out.String()), &res); err != nil {
		t.Fatalf("Sweep --json should emit valid JSON: %v\n%s", err, out.String())
	}
	if !strings.Contains(res.Next, "_mods/") {
		t.Fatalf("next-step should still point generically at _mods/, got %q", res.Next)
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
