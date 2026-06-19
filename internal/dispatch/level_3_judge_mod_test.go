// ABOUTME: AC-3 — the real docs/dev level-3-judge mod is discovered + spawned by
// ABOUTME: the existing spawn-standing-all path with model: opus, no new machinery.
package dispatch

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

// TestLevel3JudgeModSpawnsWithOpus is AC-3: the level-3-judge standing mod that
// ships in the real docs/dev/_mods is discovered and emitted by the SAME
// spawn-standing-all path as comm-officer (no new spawn machinery), carrying
// model: opus. The expected model comes from the mod file's contract, asserted
// against the binary's emitted array. The workflow dir is the real repo tree
// (../../docs/dev from this package), so the test fails if the mod is ever removed
// or its model declaration drifts.
func TestLevel3JudgeModSpawnsWithOpus(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	wd, err := filepath.Abs(filepath.Join("..", "..", "docs", "dev"))
	if err != nil {
		t.Fatal(err)
	}

	// A real team name with no matching config → the member is not already-alive,
	// so its spec is emitted.
	res := runNative("", "spawn-standing-all", "--workflow-dir", wd, "--team", "ac3-fixture-team")
	if res.exit != 0 {
		t.Fatalf("spawn-standing-all exit=%d stderr=%q", res.exit, res.stderr)
	}

	var specs []spawnAllSpec
	if err := json.Unmarshal([]byte(res.stdout), &specs); err != nil {
		t.Fatalf("stdout not a JSON array: %v\n%s", err, res.stdout)
	}

	var judge *spawnAllSpec
	for i := range specs {
		if specs[i].Name == "level-3-judge" {
			judge = &specs[i]
			break
		}
	}
	if judge == nil {
		t.Fatalf("level-3-judge not emitted by spawn-standing-all over docs/dev/_mods; got %d specs", len(specs))
	}
	if judge.Model != "opus" {
		t.Errorf("level-3-judge model = %q, want opus (the stronger-model declaration is the safety property)", judge.Model)
	}
	if judge.SubagentType != "general-purpose" {
		t.Errorf("level-3-judge subagent_type = %q, want general-purpose", judge.SubagentType)
	}
	if judge.TeamName != "ac3-fixture-team" {
		t.Errorf("level-3-judge team_name = %q, want ac3-fixture-team", judge.TeamName)
	}
	if judge.Prompt == "" {
		t.Error("level-3-judge prompt is empty — the ## Agent Prompt body was not extracted")
	}
}
