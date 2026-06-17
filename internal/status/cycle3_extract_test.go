// ABOUTME: Cycle-3 status --read additions — structured stages array (AC-3),
// ABOUTME: --fields projection (AC-6), --checklist (AC-1), --ac-scan (AC-2), loud failure (AC-5).
package status

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// devReadmePath returns the repo's real docs/dev/README.md — the external oracle
// for AC-3/AC-4 (its actual stages: block is the ground truth, so flipping a flag
// in the README would flip the emitted field, never a frozen byte golden).
func devReadmePath(t *testing.T) string {
	t.Helper()
	// internal/status -> repo root is ../..
	p, err := filepath.Abs(filepath.Join("..", "..", "docs", "dev", "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("real dev README not found at %s: %v", p, err)
	}
	return p
}

// readEnvelope is the --read --json shape carrying the cycle-3 additions: the
// existing path/total_lines/frontmatter/headings plus the new stages array.
type readEnvelope struct {
	Command     string              `json:"command"`
	Path        string              `json:"path"`
	TotalLines  string              `json:"total_lines"`
	Frontmatter map[string]string   `json:"frontmatter"`
	Headings    []map[string]string `json:"headings"`
	Stages      []map[string]string `json:"stages"`
}

// TestReadStagesArray (AC-3) asserts status --read <README> --json surfaces the
// nested stages: taxonomy as a structured array of ordered objects matching the
// REAL docs/dev/README.md field-by-field. The README is the oracle: its actual
// stages: block (backlog initial+gate, ideation gate, implementation worktree,
// validation worktree+fresh+feedback-to+gate, done terminal) is the expected set.
func TestReadStagesArray(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	readme := devReadmePath(t)

	out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--read", readme, "--json")
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	var doc readEnvelope
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}

	// The expected stages, in declaration order, with every leaf a string. Only
	// the flags the README actually sets are present; an unset bool flag is the
	// "false" string for the typed gate/terminal/initial/worktree fields (always
	// emitted), and the optional feedback-to/fresh keys appear only when set.
	type want struct {
		name       string
		worktree   string
		gate       string
		terminal   string
		initial    string
		feedbackTo string // "" means the key should be absent
		fresh      string // "" means absent
	}
	expected := []want{
		{name: "backlog", worktree: "false", gate: "true", terminal: "false", initial: "true"},
		{name: "ideation", worktree: "false", gate: "true", terminal: "false", initial: "false"},
		{name: "implementation", worktree: "true", gate: "false", terminal: "false", initial: "false"},
		{name: "validation", worktree: "true", gate: "true", terminal: "false", initial: "false", feedbackTo: "implementation", fresh: "true"},
		{name: "done", worktree: "false", gate: "false", terminal: "true", initial: "false"},
	}

	if len(doc.Stages) != len(expected) {
		t.Fatalf("stages count = %d, want %d\n%s", len(doc.Stages), len(expected), out)
	}
	for i, w := range expected {
		s := doc.Stages[i]
		if s["name"] != w.name {
			t.Fatalf("stages[%d].name = %q, want %q", i, s["name"], w.name)
		}
		check := map[string]string{
			"worktree": w.worktree, "gate": w.gate, "terminal": w.terminal, "initial": w.initial,
		}
		for k, v := range check {
			if s[k] != v {
				t.Errorf("stages[%d=%s].%s = %q, want %q", i, w.name, k, s[k], v)
			}
		}
		// Optional keys: present-with-value when the README sets them, absent otherwise.
		if w.feedbackTo == "" {
			if _, present := s["feedback-to"]; present {
				t.Errorf("stages[%d=%s] has feedback-to=%q, want absent", i, w.name, s["feedback-to"])
			}
		} else if s["feedback-to"] != w.feedbackTo {
			t.Errorf("stages[%d=%s].feedback-to = %q, want %q", i, w.name, s["feedback-to"], w.feedbackTo)
		}
		if w.fresh == "" {
			if _, present := s["fresh"]; present {
				t.Errorf("stages[%d=%s] has fresh=%q, want absent", i, w.name, s["fresh"])
			}
		} else if s["fresh"] != w.fresh {
			t.Errorf("stages[%d=%s].fresh = %q, want %q", i, w.name, s["fresh"], w.fresh)
		}
		// Every leaf is a string (the all-strings contract): no value is a JSON
		// bool/number. json.Unmarshal into map[string]string would have failed
		// already on a non-string, so reaching here proves it.
	}
}

// TestReadStagesArrayNoRegression (AC-3) asserts adding the stages array leaves
// the pre-existing default --read output untouched: the flat frontmatter still
// carries "stages":"" (the flattened scalar), and headings/total_lines for the
// README are intact.
func TestReadStagesArrayNoRegression(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	readme := devReadmePath(t)

	out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--read", readme, "--json")
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	var doc readEnvelope
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	// The flat frontmatter object STILL flattens the nested stages: block to "".
	if v, ok := doc.Frontmatter["stages"]; !ok || v != "" {
		t.Fatalf("flat frontmatter[stages] = %q (present=%v), want \"\" present — the flat map is unchanged", v, ok)
	}
	// The README has headings and a positive line count; their presence is the
	// default-read contract, which the stages addition must not disturb.
	if len(doc.Headings) == 0 {
		t.Fatal("headings empty — the default heading map regressed")
	}
	if doc.TotalLines == "" || doc.TotalLines == "0" {
		t.Fatalf("total_lines = %q, want the README's real line count", doc.TotalLines)
	}
}

// TestReadStagesArrayAbsentForPlainFile (AC-3) asserts a markdown file with NO
// stages: block emits no stages array at all (the array is keyed on the block
// existing, not always-present) — the section-reader fixture has no stages:.
func TestReadStagesArrayAbsentForPlainFile(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--read", fixturePath(t), "--json")
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	if strings.Contains(out, "\"stages\":[") {
		t.Fatalf("plain fixture (no stages: block) emitted a stages array: %s", out)
	}
}
