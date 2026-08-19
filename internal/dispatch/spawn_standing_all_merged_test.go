// ABOUTME: merged-mode spawn-standing-all — no --team: each standing spec is the
// ABOUTME: merged dispatch shape (name present, team_name absent, run_in_background).
package dispatch

import (
	"encoding/json"
	"strings"
	"testing"
)

// mergedSpawnSpec decodes a spawn-standing-all array element to the fields the
// merged shape constrains. TeamName/RunInBackground are pointers so absent (key
// omitted) is distinguishable from present-with-zero-value.
type mergedSpawnSpec struct {
	SubagentType    string  `json:"subagent_type"`
	Description     string  `json:"description"`
	Name            string  `json:"name"`
	TeamName        *string `json:"team_name"`
	RunInBackground *bool   `json:"run_in_background"`
	Model           string  `json:"model"`
	Prompt          string  `json:"prompt"`
}

// TestSpawnStandingAllMergedEmitsBackgroundSpec is the merged fixture: with NO
// --team (the merged .178+ shape) the declared standing teammate is emitted in
// the SAME merged dispatch shape build.go produces — name present, team_name
// absent, run_in_background true — with the mod-derived subagent_type/name/model/
// prompt. Expected values come from the fixture mod, not the binary.
func TestSpawnStandingAllMergedEmitsBackgroundSpec(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	wd := t.TempDir()
	writeMods(t, wd, map[string]string{
		"comm-officer.md": standingMod("comm-officer", "sonnet", "prose polisher", ""),
		"not-standing.md": "---\nstanding: false\nname: nope\n---\nbody\n",
	})

	// No --team: the merged host has no TeamCreate name to pass.
	res := runNative("", "spawn-standing-all", "--workflow-dir", wd)
	if res.exit != 0 {
		t.Fatalf("merged spawn-standing-all exit=%d stderr=%q", res.exit, res.stderr)
	}

	var specs []mergedSpawnSpec
	if err := json.Unmarshal([]byte(res.stdout), &specs); err != nil {
		t.Fatalf("stdout not a JSON array: %v\n%s", err, res.stdout)
	}
	if len(specs) != 1 {
		t.Fatalf("expected one spec (comm-officer); got %d: %v", len(specs), specs)
	}
	got := specs[0]
	if got.SubagentType != "general-purpose" {
		t.Errorf("subagent_type = %q, want general-purpose", got.SubagentType)
	}
	if want := "standing teammate: comm-officer"; got.Description != want {
		t.Errorf("description = %q, want %q", got.Description, want)
	}
	if got.Name != "comm-officer" {
		t.Errorf("name = %q, want comm-officer", got.Name)
	}
	if got.TeamName != nil {
		t.Errorf("merged spec must omit team_name (.178+ has no TeamCreate name); got %q", *got.TeamName)
	}
	if got.RunInBackground == nil || !*got.RunInBackground {
		t.Errorf("merged spec must emit run_in_background=true (named background teammate); got %v", got.RunInBackground)
	}
	if got.Model != "sonnet" {
		t.Errorf("model = %q, want sonnet", got.Model)
	}
	if want := "You are comm-officer.\n"; got.Prompt != want {
		t.Errorf("prompt = %q, want %q", got.Prompt, want)
	}
}

// TestSpawnStandingAllMergedEmitsAllDeclared is the no-dedup guard: the merged
// path has no team config to read by name, so it emits a spec for EVERY declared
// standing teammate (idempotency is the FO's own-roster concern on merged, not a
// config probe). Two declared mods => two specs, both merged-shaped.
func TestSpawnStandingAllMergedEmitsAllDeclared(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	wd := t.TempDir()
	writeMods(t, wd, map[string]string{
		"comm-officer.md": standingMod("comm-officer", "sonnet", "prose polisher", ""),
		"scribe.md":       standingMod("scribe", "haiku", "note taker", ""),
	})
	// A team config naming comm-officer must NOT dedup it on the merged path —
	// merged has no team name, so the config is irrelevant here.
	claudeFixture{
		team:    "fixture-team",
		session: "s",
		members: []fixtureMember{{name: "comm-officer", model: "sonnet"}},
		jsonls:  map[string]string{},
	}.write(t, home)

	res := runNative("", "spawn-standing-all", "--workflow-dir", wd)
	if res.exit != 0 {
		t.Fatalf("merged spawn-standing-all exit=%d stderr=%q", res.exit, res.stderr)
	}
	var specs []mergedSpawnSpec
	if err := json.Unmarshal([]byte(res.stdout), &specs); err != nil {
		t.Fatalf("stdout not a JSON array: %v\n%s", err, res.stdout)
	}
	if len(specs) != 2 {
		t.Fatalf("merged path must emit all declared specs with no config dedup; got %d: %v", len(specs), specs)
	}
	for _, s := range specs {
		if s.TeamName != nil {
			t.Errorf("merged spec %q must omit team_name; got %q", s.Name, *s.TeamName)
		}
		if s.RunInBackground == nil || !*s.RunInBackground {
			t.Errorf("merged spec %q must emit run_in_background=true; got %v", s.Name, s.RunInBackground)
		}
	}
}

// TestSpawnStandingAllMergedEmptyWhenNoneDeclared guards the degenerate merged
// path: no standing mods => `[]`, exit 0, just like the legacy empty case.
func TestSpawnStandingAllMergedEmptyWhenNoneDeclared(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	wd := t.TempDir()
	writeMods(t, wd, map[string]string{
		"not-standing.md": "---\nstanding: false\nname: nope\n---\nbody\n",
	})

	res := runNative("", "spawn-standing-all", "--workflow-dir", wd)
	if res.exit != 0 {
		t.Fatalf("merged spawn-standing-all exit=%d stderr=%q", res.exit, res.stderr)
	}
	if got := strings.TrimSpace(res.stdout); got != "[]" {
		t.Fatalf("expected empty array (no standing mods); got %q", got)
	}
}

// TestSpawnStandingAllMergedFailsLoudOnBrokenMod guards that the merged path
// keeps the same loud validation the legacy path has: a standing mod missing its
// ## Agent Prompt exits 1 naming the offending mod, with or without --team.
func TestSpawnStandingAllMergedFailsLoudOnBrokenMod(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	wd := t.TempDir()
	writeMods(t, wd, map[string]string{
		"broken.md": "---\nstanding: true\nname: broken\n---\n## Hook: startup\n- subagent_type: general-purpose\n- name: broken\n- model: sonnet\n",
	})

	res := runNative("", "spawn-standing-all", "--workflow-dir", wd)
	if res.exit != 1 {
		t.Fatalf("expected exit 1 on broken mod (merged); got %d stdout=%q", res.exit, res.stdout)
	}
	if !strings.Contains(res.stderr, "broken.md") || !strings.Contains(res.stderr, "## Agent Prompt") {
		t.Errorf("stderr does not name the offending mod and missing section: %q", res.stderr)
	}
}

// TestSpawnStandingAllMergedFableModel is AC-7 ported to the merged path (the
// legacy-era TestSpawnStandingParityFableModel drove the retired `spawn-standing`
// singular subcommand): a standing mod declaring `model: fable` in its
// ## Hook: startup spawns successfully through spawn-standing-all's shared
// buildSpawnSpec validation.
func TestSpawnStandingAllMergedFableModel(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	wd := t.TempDir()
	writeMods(t, wd, map[string]string{
		"fable-officer.md": standingMod("fable-officer", "fable", "fable ensign", ""),
	})

	res := runNative("", "spawn-standing-all", "--workflow-dir", wd)
	if res.exit != 0 {
		t.Fatalf("spawn-standing-all exit=%d stderr=%q", res.exit, res.stderr)
	}
	if !strings.Contains(res.stdout, `"model": "fable"`) {
		t.Errorf("spec array missing model=fable:\n%s", res.stdout)
	}
}

// TestSpawnStandingAllMergedNonASCIIPrompt is the A-2 non-ASCII parity case
// ported to the merged path (the legacy-era TestSpawnStandingParitySpecNonASCIIPrompt
// drove the retired `spawn-standing` singular subcommand): the FO forwards
// spec.prompt VERBATIM, so a non-ASCII Agent Prompt (em-dash U+2014 here) must
// serialize through the shared ensure_ascii escaping spawn-standing-all's array
// emission uses.
func TestSpawnStandingAllMergedNonASCIIPrompt(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	wd := t.TempDir()
	mod := "---\nstanding: true\nname: comm-officer\n---\n" +
		"## Hook: startup\n" +
		"- subagent_type: general-purpose\n" +
		"- name: comm-officer\n" +
		"- model: sonnet\n" +
		"## Agent Prompt\nYou are comm-officer — the prose polisher.\n"
	writeMods(t, wd, map[string]string{"comm-officer.md": mod})

	res := runNative("", "spawn-standing-all", "--workflow-dir", wd)
	if res.exit != 0 {
		t.Fatalf("spawn-standing-all exit=%d stderr=%q", res.exit, res.stderr)
	}
	assertEmDashEscaped(t, res.stdout)
}

// TestSpawnStandingAllMergedLoudFailures ports the mode-neutral validation cases
// from the retired `spawn-standing` singular subcommand's parity fixture
// (missing model, bad model enum, and a trailing heading after ## Agent Prompt)
// onto spawn-standing-all: each still exits 1 with a stderr line naming the
// offending mod through the shared buildSpawnSpec validation.
func TestSpawnStandingAllMergedLoudFailures(t *testing.T) {
	cases := []struct {
		name       string
		modName    string
		mod        string
		wantStderr []string
	}{
		{
			name:       "missing-model",
			modName:    "nomodel.md",
			mod:        "---\nstanding: true\nname: nomodel\n---\n## Hook: startup\n- subagent_type: general-purpose\n- name: nomodel\n## Agent Prompt\nx\n",
			wantStderr: []string{"nomodel.md", "no 'model'"},
		},
		{
			name:       "bad-model-enum",
			modName:    "badmodel.md",
			mod:        "---\nstanding: true\nname: badmodel\n---\n## Hook: startup\n- subagent_type: general-purpose\n- name: badmodel\n- model: gpt-4\n## Agent Prompt\nx\n",
			wantStderr: []string{"badmodel.md", "invalid model 'gpt-4'"},
		},
		{
			name:       "trailing-heading",
			modName:    "trailer.md",
			mod:        "---\nstanding: true\nname: trailer\n---\n## Hook: startup\n- subagent_type: general-purpose\n- name: trailer\n- model: opus\n## Agent Prompt\nbody\n## Trailing Section\noops\n",
			wantStderr: []string{"trailer.md", "trailing top-level heading"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			wd := t.TempDir()
			writeMods(t, wd, map[string]string{tc.modName: tc.mod})

			res := runNative("", "spawn-standing-all", "--workflow-dir", wd)
			if res.exit != 1 {
				t.Fatalf("expected exit 1 for %q; got %d stdout=%q", tc.name, res.exit, res.stdout)
			}
			for _, want := range tc.wantStderr {
				if !strings.Contains(res.stderr, want) {
					t.Errorf("%s: stderr %q missing %q", tc.name, res.stderr, want)
				}
			}
		})
	}
}
