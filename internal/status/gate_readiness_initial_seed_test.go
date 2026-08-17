// ABOUTME: A gated INITIAL stage has no prior stage to have written a report —
// ABOUTME: there the committed, clean-in-HEAD seed is itself the promotion proof.
package status

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/gates"
	"github.com/spacedock-dev/spacedock/internal/testgit"
)

// initialGateReadme declares a gated INITIAL stage and a gated non-initial one,
// so the promotion and the guard it must not weaken share a single workflow.
const initialGateReadme = `---
commissioned-by: spacedock@1
id-style: slug
state: .spacedock-state
stages:
  states:
    - name: backlog
      initial: true
      gate: true
    - name: implementation
    - name: validation
      gate: true
    - name: done
      terminal: true
---

# Initial Gated Workflow
`

// readinessByID projects gate-readiness per entity id.
func readinessByID(t *testing.T, def string) map[string]string {
	t.Helper()
	out, errOut, code := runNative(t, def, pinnedEnv(t), "--workflow-dir", def, "--fields", "id,gate-readiness", "--json")
	if code != 0 {
		t.Fatalf("fields exit=%d stderr=%q", code, errOut)
	}
	var projected struct {
		Entities []map[string]string `json:"entities"`
	}
	if err := json.Unmarshal([]byte(out), &projected); err != nil {
		t.Fatalf("parse fields: %v\n%s", err, out)
	}
	got := map[string]string{}
	for _, entity := range projected.Entities {
		got[entity["id"]] = entity["gate-readiness"]
	}
	return got
}

// readyGateReadinessBySlug maps each ready_gates row to its readiness.
func readyGateReadinessBySlug(t *testing.T, def string) map[string]string {
	t.Helper()
	out, errOut, code := runNative(t, def, pinnedEnv(t), "--workflow-dir", def, "--boot", "--identify", "--json")
	if code != 0 {
		t.Fatalf("boot exit=%d stderr=%q", code, errOut)
	}
	var boot struct {
		ReadyGates []map[string]string `json:"ready_gates"`
	}
	if err := json.Unmarshal([]byte(out), &boot); err != nil {
		t.Fatalf("parse boot: %v\n%s", err, out)
	}
	rows := map[string]string{}
	for _, row := range boot.ReadyGates {
		rows[row["slug"]] = row["readiness"]
	}
	return rows
}

// TestInitialGatedSeedReachesCaptainWithNoStageReport is AC-1, the deadlock
// fixture. The recorded red baseline was {"dispatchable":[],"ready_gates":[]}:
// dispatchable stays empty because a gated stage is never ordinarily
// dispatchable, while ready_gates gains the row whose absence was the deadlock.
// Reverting gatePreparable empties ready_gates again; any fix that instead
// demands a report makes the closing Stage Report count nonzero.
func TestInitialGatedSeedReachesCaptainWithNoStageReport(t *testing.T) {
	def, state := buildSplitRoot(t, initialGateReadme, map[string]string{
		"seed.md": "---\nid: seed\nstatus: backlog\ntitle: Seed\n---\n# Seed\n\n## Problem\n\nA committed seed carrying no report.\n",
	})
	entity := filepath.Join(state, "seed.md")
	testgit.InitRepo(t, state, "-q")
	gitC(t, state, "add", "--", "seed.md")
	gitC(t, state, "commit", "-q", "-m", "commit the seed", "--", "seed.md")
	// Preparation classifies the selected source against BOTH workflow roots, so
	// the definition checkout the state dir is nested in must be a repo too.
	testgit.InitRepo(t, def, "-q")
	gitC(t, def, "add", "--", "README.md")
	gitC(t, def, "commit", "-q", "-m", "commit the workflow definition", "--", "README.md")

	out, errOut, code := runNative(t, def, pinnedEnv(t), "--workflow-dir", def, "--next", "--json")
	if code != 0 {
		t.Fatalf("next exit=%d stderr=%q", code, errOut)
	}
	var next struct {
		Dispatchable []map[string]string `json:"dispatchable"`
		ReadyGates   []map[string]string `json:"ready_gates"`
	}
	if err := json.Unmarshal([]byte(out), &next); err != nil {
		t.Fatalf("parse next: %v\n%s", err, out)
	}
	if len(next.Dispatchable) != 0 {
		t.Fatalf("dispatchable=%v, want empty: a gated stage stays excluded from ordinary dispatch", next.Dispatchable)
	}
	if len(next.ReadyGates) != 1 || next.ReadyGates[0]["readiness"] != "needs-preparation" {
		t.Fatalf("ready_gates=%v, want one needs-preparation row\n%s", next.ReadyGates, out)
	}

	// Prepare directly from the seed, exactly as the amended gate lifecycle
	// instructs: --artifact names the entity, and no report is authored.
	if _, err := gates.Prepare(entity, gates.PrepareInput{
		WorkflowDir: def,
		Question:    "Approve this seed?",
		Artifact:    entity,
		Summary:     "Seed review",
	}); err != nil {
		t.Fatalf("prepare from the seed: %v", err)
	}
	if got := readinessByID(t, def)["seed"]; got != "awaiting-captain" {
		t.Fatalf("prepared seed gate-readiness = %q, want awaiting-captain", got)
	}

	body, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(string(body), "## Stage Report"); n != 0 {
		t.Fatalf("entity carries %d Stage Report sections, want 0: no report may be fabricated", n)
	}
}

// TestNonInitialGatedStageStillOwesItsReport is AC-2. The exception keys on
// stage.initial and never on "no report exists"; keying it on the latter would
// promote the report-less non-initial entity here.
func TestNonInitialGatedStageStillOwesItsReport(t *testing.T) {
	def, state := buildSplitRoot(t, initialGateReadme, map[string]string{
		"bare.md":     "---\nid: bare\nstatus: validation\n---\n# Bare\n",
		"reported.md": "---\nid: reported\nstatus: validation\n---\n# Reported\n" + needsPreparationReport,
	})
	testgit.InitRepo(t, state, "-q")
	gitC(t, state, "add", "-A")
	gitC(t, state, "commit", "-q", "-m", "commit both guard controls")

	readiness := readinessByID(t, def)
	if readiness["bare"] != "validating" {
		t.Fatalf("report-less non-initial gated stage = %q, want validating", readiness["bare"])
	}
	if readiness["reported"] != "needs-preparation" {
		t.Fatalf("reported non-initial gated stage = %q, want needs-preparation", readiness["reported"])
	}
	rows := readyGateReadinessBySlug(t, def)
	if _, present := rows["bare"]; present {
		t.Fatalf("report-less non-initial stage appeared in ready_gates: %v", rows)
	}
	if rows["reported"] != "needs-preparation" {
		t.Fatalf("ready_gates[reported]=%q, want needs-preparation", rows["reported"])
	}
}

// TestInitialGatedSeedMustBeCommittedAndClean is AC-3. Durability is the entire
// proof at an initial stage, so dropping the entityPathCleanInHEAD call would
// promote both the untracked and the dirty seed here.
func TestInitialGatedSeedMustBeCommittedAndClean(t *testing.T) {
	const seed = "---\nid: fresh\nstatus: backlog\n---\n# Fresh\n"
	def, state := buildSplitRoot(t, initialGateReadme, map[string]string{"fresh.md": seed})
	testgit.InitRepo(t, state, "-q")
	if got := readinessByID(t, def)["fresh"]; got != "validating" {
		t.Fatalf("untracked seed = %q, want validating", got)
	}
	gitC(t, state, "add", "--", "fresh.md")
	gitC(t, state, "commit", "-q", "-m", "commit the seed", "--", "fresh.md")
	if got := readinessByID(t, def)["fresh"]; got != "needs-preparation" {
		t.Fatalf("committed clean seed = %q, want needs-preparation", got)
	}
	writeFile(t, filepath.Join(state, "fresh.md"), seed+"\nUncommitted dirt.\n")
	if got := readinessByID(t, def)["fresh"]; got != "validating" {
		t.Fatalf("dirty seed = %q, want validating", got)
	}
}
