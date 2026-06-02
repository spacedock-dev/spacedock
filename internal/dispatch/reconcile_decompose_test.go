// ABOUTME: table-driven decompose tests — covers prefix strip, cycle stripping,
// ABOUTME: longest-stage match, custom-stage support, and ambiguous-name rejection.
package dispatch

import (
	"path/filepath"
	"testing"
)

func TestDecomposeCanonicalStages(t *testing.T) {
	stages := readStageNames(t.TempDir()) // empty workflow dir → canonical only
	cases := []struct {
		name   string
		input  string
		slug   string
		stage  string
		cycle  string
		wantOk bool
	}{
		{"unsuffixed ideation", "spacedock-ensign-foo-ideation", "foo", "ideation", "", true},
		{"unsuffixed implementation",
			"spacedock-ensign-yaml-parser-migration-implementation",
			"yaml-parser-migration", "implementation", "", true},
		{"hyphenated slug + validation",
			"spacedock-ensign-ensign-lifecycle-reconcile-validation",
			"ensign-lifecycle-reconcile", "validation", "", true},
		{"cycleN suffix",
			"spacedock-ensign-foo-implementation-cycle2",
			"foo", "implementation", "cycle2", true},
		{"bare-digit cycle suffix",
			"spacedock-ensign-foo-validation-3",
			"foo", "validation", "3", true},
		{"backlog stage",
			"spacedock-ensign-bar-backlog",
			"bar", "backlog", "", true},
		{"done stage",
			"spacedock-ensign-baz-done",
			"baz", "done", "", true},
		{"longest-stage greedy: implementation wins over an inner match",
			// "implementation" is a real stage so the matcher peels the right suffix
			// rather than a false fragment.
			"spacedock-ensign-x-implementation",
			"x", "implementation", "", true},
		{"missing worker prefix", "foo-implementation", "", "", "", false},
		{"empty after prefix", "spacedock-ensign-", "", "", "", false},
		{"only stage suffix, empty slug", "spacedock-ensign-implementation", "", "", "", false},
		{"unknown stage", "spacedock-ensign-foo-deploy", "", "", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := decompose(c.input, stages)
			if got.ok != c.wantOk {
				t.Fatalf("ok=%v, want %v (input=%q)", got.ok, c.wantOk, c.input)
			}
			if !c.wantOk {
				return
			}
			if got.slug != c.slug || got.stage != c.stage || got.cycle != c.cycle {
				t.Fatalf("decompose(%q) = {%q %q %q}, want {%q %q %q}",
					c.input, got.slug, got.stage, got.cycle, c.slug, c.stage, c.cycle)
			}
		})
	}
}

// TestDecomposeWithCustomWorkflowStages confirms a workflow declaring custom
// stages adds them to the stage set so a non-canonical stage name still
// decomposes cleanly.
func TestDecomposeWithCustomWorkflowStages(t *testing.T) {
	wd := t.TempDir()
	writeFile(t, filepath.Join(wd, "README.md"), `---
entity-type: task
stages:
  states:
    - name: triage
      initial: true
    - name: design
    - name: build
      worktree: true
    - name: ship
      terminal: true
---
`)
	stages := readStageNames(wd)

	cases := []struct {
		name  string
		input string
		slug  string
		stage string
	}{
		{"triage custom", "spacedock-ensign-thing-triage", "thing", "triage"},
		{"design custom", "spacedock-ensign-pretty-design", "pretty", "design"},
		{"build custom", "spacedock-ensign-widget-build", "widget", "build"},
		{"ship custom", "spacedock-ensign-rocket-ship", "rocket", "ship"},
		{"canonical still works",
			"spacedock-ensign-rocket-validation",
			"rocket", "validation"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := decompose(c.input, stages)
			if !got.ok {
				t.Fatalf("decompose(%q): ok=false, want true", c.input)
			}
			if got.slug != c.slug || got.stage != c.stage {
				t.Fatalf("decompose(%q) = {%q %q}, want {%q %q}",
					c.input, got.slug, got.stage, c.slug, c.stage)
			}
		})
	}
}

// TestParseIncludeRoundTrip checks the --include flag parser accepts the empty
// default (= all classes) and a comma-separated subset, and rejects an unknown
// class.
func TestParseIncludeRoundTrip(t *testing.T) {
	got, err := parseInclude("")
	if err != nil {
		t.Fatalf("empty include: %v", err)
	}
	for _, c := range []string{"A", "B", "C", "D", "E"} {
		if !got[c] {
			t.Errorf("empty include should enable %s", c)
		}
	}
	subset, err := parseInclude("A,B")
	if err != nil {
		t.Fatalf("A,B include: %v", err)
	}
	if !subset["A"] || !subset["B"] || subset["C"] || subset["D"] || subset["E"] {
		t.Errorf("A,B include wrong: %#v", subset)
	}
	if _, err := parseInclude("Z"); err == nil {
		t.Errorf("Z include should error")
	}
}
