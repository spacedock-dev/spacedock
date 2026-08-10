// ABOUTME: Cold gated-stage recovery promotes only mechanically complete,
// ABOUTME: committed reports to the First Officer's preparation queue.
package status

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/testgit"
)

const needsPreparationReport = `
## Stage Report: validation

- DONE: Verify the validation obligations.
  Commit abc123 contains the required evidence.

### Summary

Validation is complete and ready for review.
`

func TestGateReadinessPromotesMechanicallyCompleteColdReports(t *testing.T) {
	cases := []struct {
		name   string
		entity string
	}{
		{
			name:   "absent gate authority",
			entity: "---\nid: absent\nstatus: validation\nscore: 90\n---\n# Absent\n" + needsPreparationReport,
		},
		{
			name:   "prior-stage authority",
			entity: strings.Replace(openGateEntity("prior", "draft", "80"), "status: draft", "status: validation", 1) + needsPreparationReport,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			def, state := buildSplitRoot(t, identifyReadyGatesReadme, map[string]string{
				filepath.Base(tc.name) + ".md": tc.entity,
			})
			// The entity must be tracked and clean: the gqs proof is path-scoped
			// and deliberately does not treat an uncommitted report as durable.
			entity := filepath.Join(state, filepath.Base(tc.name)+".md")
			testgit.InitRepo(t, state, "-q")
			gitC(t, state, "add", "--", filepath.Base(entity))
			gitC(t, state, "commit", "-q", "-m", "seed complete gated report", "--", filepath.Base(entity))

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
			if len(boot.ReadyGates) != 1 || boot.ReadyGates[0]["readiness"] != "needs-preparation" {
				t.Fatalf("ready_gates=%v, want one needs-preparation row\n%s", boot.ReadyGates, out)
			}

			nextOut, nextErr, nextCode := runNative(t, def, pinnedEnv(t), "--workflow-dir", def, "--next", "--json")
			if nextCode != 0 {
				t.Fatalf("next exit=%d stderr=%q", nextCode, nextErr)
			}
			var next struct {
				Command      string              `json:"command"`
				Dispatchable []map[string]string `json:"dispatchable"`
				ReadyGates   []map[string]string `json:"ready_gates"`
			}
			if err := json.Unmarshal([]byte(nextOut), &next); err != nil {
				t.Fatalf("parse next: %v\n%s", err, nextOut)
			}
			if next.Command != "next" || len(next.Dispatchable) != 0 || len(next.ReadyGates) != 1 || next.ReadyGates[0]["readiness"] != "needs-preparation" {
				t.Fatalf("next envelope dispatchable=%v ready_gates=%v, want empty dispatchable + one candidate", next.Dispatchable, next.ReadyGates)
			}
			fieldOut, fieldErr, fieldCode := runNative(t, def, pinnedEnv(t), "--workflow-dir", def, "--fields", "id,gate-readiness", "--json")
			var projected struct {
				Entities []map[string]string `json:"entities"`
			}
			if fieldCode != 0 || json.Unmarshal([]byte(fieldOut), &projected) != nil || len(projected.Entities) != 1 || projected.Entities[0]["gate-readiness"] != "needs-preparation" {
				t.Fatalf("field projection exit=%d stderr=%q output=%q", fieldCode, fieldErr, fieldOut)
			}
		})
	}
}

func TestGateReadinessRejectsDirtyAndMalformedColdReports(t *testing.T) {
	def, state := buildSplitRoot(t, identifyReadyGatesReadme, map[string]string{
		"dirty.md":     "---\nid: dirty\nstatus: validation\n---\n# Dirty\n" + needsPreparationReport,
		"malformed.md": "---\nid: malformed\nstatus: validation\ngates:\n  version: 1\n  records: []\n---\n# Malformed\n" + needsPreparationReport,
	})
	testgit.InitRepo(t, state, "-q")
	gitC(t, state, "add", "-A")
	gitC(t, state, "commit", "-q", "-m", "seed report controls")
	// Dirty target bytes are never mechanical proof, even if the report itself
	// is still structurally complete. The malformed gates document is fail-closed.
	writeFile(t, filepath.Join(state, "dirty.md"), "---\nid: dirty\nstatus: validation\n---\n# Dirty\n"+needsPreparationReport+"\nUncommitted dirt.\n")
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
	for _, row := range boot.ReadyGates {
		if strings.HasPrefix(row["slug"], "dirty") || strings.HasPrefix(row["slug"], "malformed") {
			t.Fatalf("invalid control appeared in ready_gates: %v", row)
		}
	}
}
