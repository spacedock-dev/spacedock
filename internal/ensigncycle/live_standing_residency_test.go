//go:build live

// ABOUTME: Live residency proof — a real FO booted against a workflow declaring a
// ABOUTME: standing comm-officer mod injects it into the team config.json roster.
package ensigncycle

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestLiveStandingResidencyInjectsCommOfficer is the AC-2 behavioral proof for
// the standing-teammate mechanism relocation (spawn-standing-all). It boots the
// REAL `spacedock claude` front door against a fixture workflow that declares a
// standing comm-officer mod, lets the FO reach its first team-mode dispatch (where
// the contract's single trigger line runs `spacedock dispatch spawn-standing-all`,
// which injects the declared standing teammates), and asserts comm-officer is
// PRESENT in the team `config.json` members roster afterward.
//
// The independent oracle is the on-disk team config the live Claude Code harness
// writes — NOT a contract grep. The test REDS if residency breaks (comm-officer
// absent after first dispatch), which is exactly the cycle-1/2 regression this
// cycle-3 design must not reintroduce: the mechanism moved out of contract prose
// into the binary + the mod's self-declaration, but the resident teammate must
// still land in the roster.
//
// The `//go:build live` tag keeps this out of the secret-free offline suite; it
// runs only under `go test -tags live` behind the CI-E2E approval gate, alongside
// TestLiveEnsignCycle. Auth + HOME isolation reuse isolatedClaudeEnv (skips when
// no credential is available; never fatals).
func TestLiveStandingResidencyInjectsCommOfficer(t *testing.T) {
	binary := spacedockBinary(t)
	repoRoot := repoRoot(t)
	model := envOr("SPACEDOCK_LIVE_MODEL", "sonnet")

	childEnv := isolatedClaudeEnv(t, os.Getenv("HOME"))
	childEnv = withBinaryOnPath(childEnv, binary)

	// The team config the live harness writes lives under CLAUDE_CONFIG_DIR/teams;
	// resolve it from the same child env the FO subprocess runs under so the roster
	// read targets exactly where Claude Code persisted the members list.
	configDir, ok := envValue(childEnv, "CLAUDE_CONFIG_DIR")
	if !ok {
		t.Fatal("child env carries no CLAUDE_CONFIG_DIR — cannot locate the team config roster")
	}

	// Stage the realistic ≥3-stage lifecycle fixture (same as TestLiveEnsignCycle)
	// PLUS a `_mods/comm-officer.md` declaring a standing teammate. The standing mod
	// is what spawn-standing-all enumerates at first dispatch; its presence is the
	// only difference from the plain cycle fixture.
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeRealisticLifecycle())
	writeFile(t, filepath.Join(root, "_mods", "comm-officer.md"), liveStandingCommOfficerMod())
	entityPath := filepath.Join(root, "make-it-work.md")
	writeFile(t, entityPath, entityFixture())
	gitInit(t, root)

	task := "Use $spacedock:first-officer for this whole run."
	drivePrompt := "Drive the workflow. " + antiShutdownOverride
	cmd := exec.Command(binary, "claude",
		"--plugin-dir", repoRoot,
		"--skip-contract-check",
		"--",
		"-p", drivePrompt,
		"--permission-mode", "bypassPermissions",
		"--output-format", "stream-json",
		"--verbose",
		"--model", model,
		task,
	)
	cmd.Dir = root
	cmd.Env = childEnv

	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw
	if err := cmd.Start(); err != nil {
		t.Fatalf("spacedock claude failed to start: %v", err)
	}

	poller := newCmdPoller(cmd, pw)
	watcher := newStreamWatcher(newPipeLineSource(pr), poller, func(line string) { t.Log(line) })
	defer poller.kill()

	// Watch the same proven progress beats as TestLiveEnsignCycle: TeamCreate
	// engages teams mode, then the ensign dispatch closes. spawn-standing-all runs
	// BEFORE the first team-mode Agent() dispatch, so by the time the dispatch has
	// closed, the standing teammate has been injected into the team config.
	if _, err := watcher.expect(isTeamCreate, quietBudgetDefault, "TeamCreate"); err != nil {
		if wrongRoot := detectWrongRootBoot(watcher.fullTranscript(), root); wrongRoot != nil {
			t.Fatalf("live residency cycle failed at TeamCreate due to a wrong-root boot: %v\nUnderlying watcher error: %v", wrongRoot, err)
		}
		t.Fatalf("live residency cycle failed at TeamCreate: %v", err)
	}
	if err := watcher.expectDispatchClose(quietBudgetDefault, "dispatch close"); err != nil {
		t.Fatalf("live residency cycle failed at the ensign dispatch close: %v", err)
	}

	// The load-bearing assertion: comm-officer is a member of the live team's
	// config.json roster. The roster is the independent oracle the live harness
	// wrote; a missing comm-officer means residency broke (the cycle-1/2 mistake).
	teams := teamConfigPaths(t, configDir)
	if len(teams) == 0 {
		t.Fatalf("no team config.json found under %s/teams after the cycle (TeamCreate matched but no roster persisted)", configDir)
	}
	if !anyTeamHasMember(t, teams, "comm-officer") {
		t.Fatalf("comm-officer ABSENT from every team roster after first dispatch — residency broke (standing-teammate injection did not land it in config.json)\nscanned: %v", teams)
	}
	t.Logf("comm-officer present in the live team roster across %d team config(s)", len(teams))
}

// liveStandingCommOfficerMod is a minimal standing comm-officer mod for the live
// residency fixture: standing: true frontmatter plus the spawn-config bullets the
// binary parses and an ## Agent Prompt body. It mirrors the shipped mod's
// self-declaration shape (frontmatter standing + trimmed ## Hook: startup +
// ## Agent Prompt) without the full routing prose, which is irrelevant to the
// roster-membership proof. The agent prompt's online handshake is harmless in the
// fixture: the spawned member lands in the roster regardless of whether it then
// finds an elements-of-style skill.
func liveStandingCommOfficerMod() string {
	return "---\n" +
		"name: comm-officer\n" +
		"description: Standing prose-polishing teammate for this workflow\n" +
		"standing: true\n" +
		"---\n" +
		"# Comm Officer\n\n" +
		"## Hook: startup\n\n" +
		"- subagent_type: general-purpose\n" +
		"- name: comm-officer\n" +
		"- model: sonnet\n\n" +
		"## Agent Prompt\n\n" +
		"You are the session's communications officer. Polish prose for clarity and " +
		"concision when asked. On spawn, send `comm-officer online` to team-lead, then " +
		"idle until you receive a polish request.\n"
}

// teamConfigPaths returns every {configDir}/teams/*/config.json path on disk. The
// live harness writes one per team; the live FO creates exactly one team, but the
// scan is plural-safe so a parallel-team layout does not silently miss the roster.
func teamConfigPaths(t *testing.T, configDir string) []string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(configDir, "teams", "*", "config.json"))
	if err != nil {
		t.Fatalf("glob team configs under %s: %v", configDir, err)
	}
	return matches
}

// anyTeamHasMember reports whether any of the given team config.json files lists a
// member with the given name. It parses the same members[] shape MemberExists
// reads (a list of objects each carrying a name), so the live oracle matches the
// binary's own membership predicate.
func anyTeamHasMember(t *testing.T, configPaths []string, name string) bool {
	t.Helper()
	for _, p := range configPaths {
		raw, err := os.ReadFile(p)
		if err != nil {
			t.Logf("read team config %s: %v", p, err)
			continue
		}
		var cfg struct {
			Members []struct {
				Name string `json:"name"`
			} `json:"members"`
		}
		if err := json.Unmarshal(raw, &cfg); err != nil {
			t.Logf("parse team config %s: %v", p, err)
			continue
		}
		for _, m := range cfg.Members {
			if strings.EqualFold(m.Name, name) {
				return true
			}
		}
	}
	return false
}
