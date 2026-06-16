// ABOUTME: spawn-standing-all driver tests — array-emit shape, already-alive dedup,
// ABOUTME: empty `[]` when none declared, and loud failure on a broken standing mod.
package dispatch

import (
	"encoding/json"
	"strings"
	"testing"
)

// spawnAllSpec mirrors the array element spawn-standing-all emits. Expected
// values are sourced from the fixture mod, NOT the binary.
type spawnAllSpec struct {
	SubagentType string `json:"subagent_type"`
	Description  string `json:"description"`
	Name         string `json:"name"`
	TeamName     string `json:"team_name"`
	Model        string `json:"model"`
	Prompt       string `json:"prompt"`
}

// TestSpawnStandingAllEmitsAbsentMemberSpec drives the inject loop with the
// declared member ABSENT from the team config: spawn-standing-all emits a
// one-element JSON array whose spec carries subagent_type/name/team_name/model/
// prompt, with the prompt equal to the fixture mod's ## Agent Prompt body (the
// expected value comes from the fixture, an independent source, not the binary).
func TestSpawnStandingAllEmitsAbsentMemberSpec(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	wd := t.TempDir()
	writeMods(t, wd, map[string]string{
		"comm-officer.md": standingMod("comm-officer", "sonnet", "prose polisher", ""),
		"not-standing.md": "---\nstanding: false\nname: nope\n---\nbody\n",
	})
	// A team config that does NOT list comm-officer, so MemberExists is false.
	claudeFixture{
		team:    "fixture-team",
		session: "s",
		members: []fixtureMember{{name: "team-lead", model: "opus"}},
		jsonls:  map[string]string{},
	}.write(t, home)

	res := runNative("", "spawn-standing-all", "--workflow-dir", wd, "--team", "fixture-team")
	if res.exit != 0 {
		t.Fatalf("exit=%d stderr=%q", res.exit, res.stderr)
	}

	var specs []spawnAllSpec
	if err := json.Unmarshal([]byte(res.stdout), &specs); err != nil {
		t.Fatalf("stdout not a JSON array: %v\n%s", err, res.stdout)
	}
	if len(specs) != 1 {
		t.Fatalf("expected one spec (comm-officer absent); got %d: %v", len(specs), specs)
	}
	got := specs[0]
	if got.SubagentType != "general-purpose" {
		t.Errorf("subagent_type = %q, want general-purpose", got.SubagentType)
	}
	// The Agent tool REQUIRES description; the spec must carry a non-empty one
	// (spawn-standing-all derives it from the member name) or the forwarded
	// Agent() call fails InputValidationError and the standing teammate never
	// spawns.
	if want := "standing teammate: comm-officer"; got.Description != want {
		t.Errorf("description = %q, want %q", got.Description, want)
	}
	if got.Name != "comm-officer" {
		t.Errorf("name = %q, want comm-officer", got.Name)
	}
	if got.TeamName != "fixture-team" {
		t.Errorf("team_name = %q, want fixture-team", got.TeamName)
	}
	if got.Model != "sonnet" {
		t.Errorf("model = %q, want sonnet", got.Model)
	}
	// The fixture's ## Agent Prompt body (from standingMod) is "You are comm-officer.\n".
	if want := "You are comm-officer.\n"; got.Prompt != want {
		t.Errorf("prompt = %q, want %q", got.Prompt, want)
	}
}

// TestSpawnStandingAllSkipsAlreadyAliveMember drives the dedup path: the team
// config already lists comm-officer, so spawn-standing-all OMITS it from the
// array. With it the only declared standing mod, the array is empty `[]`.
func TestSpawnStandingAllSkipsAlreadyAliveMember(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	wd := t.TempDir()
	writeMods(t, wd, map[string]string{
		"comm-officer.md": standingMod("comm-officer", "sonnet", "prose polisher", ""),
	})
	claudeFixture{
		team:    "fixture-team",
		session: "s",
		members: []fixtureMember{{name: "comm-officer", model: "sonnet"}},
		jsonls:  map[string]string{},
	}.write(t, home)

	res := runNative("", "spawn-standing-all", "--workflow-dir", wd, "--team", "fixture-team")
	if res.exit != 0 {
		t.Fatalf("exit=%d stderr=%q", res.exit, res.stderr)
	}
	if got := strings.TrimSpace(res.stdout); got != "[]" {
		t.Fatalf("expected empty array (already-alive deduped); got %q", got)
	}
}

// TestSpawnStandingAllEmptyWhenNoneDeclared drives the degenerate path: a
// workflow with no standing mods emits `[]`, exit 0.
func TestSpawnStandingAllEmptyWhenNoneDeclared(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	wd := t.TempDir()
	writeMods(t, wd, map[string]string{
		"not-standing.md": "---\nstanding: false\nname: nope\n---\nbody\n",
	})

	res := runNative("", "spawn-standing-all", "--workflow-dir", wd, "--team", "fixture-team")
	if res.exit != 0 {
		t.Fatalf("exit=%d stderr=%q", res.exit, res.stderr)
	}
	if got := strings.TrimSpace(res.stdout); got != "[]" {
		t.Fatalf("expected empty array (no standing mods); got %q", got)
	}
}

// TestSpawnStandingAllFailsLoudOnBrokenMod drives the loud-failure path: a
// declared standing mod missing its ## Agent Prompt exits 1 with a stderr line
// naming the offending mod — the same validation runSpawnStanding enforces,
// shared via buildSpawnSpec rather than re-implemented.
func TestSpawnStandingAllFailsLoudOnBrokenMod(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	wd := t.TempDir()
	writeMods(t, wd, map[string]string{
		"broken.md": "---\nstanding: true\nname: broken\n---\n## Hook: startup\n- subagent_type: general-purpose\n- name: broken\n- model: sonnet\n",
	})

	res := runNative("", "spawn-standing-all", "--workflow-dir", wd, "--team", "fixture-team")
	if res.exit != 1 {
		t.Fatalf("expected exit 1 on broken mod; got %d stdout=%q", res.exit, res.stdout)
	}
	if !strings.Contains(res.stderr, "broken.md") || !strings.Contains(res.stderr, "## Agent Prompt") {
		t.Errorf("stderr does not name the offending mod and missing section: %q", res.stderr)
	}
}
