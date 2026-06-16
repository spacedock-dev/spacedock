//go:build live

// ABOUTME: Live team-mode e2e over the pty/tmux driver — resurrects the two retired forced-team
// ABOUTME: tests: comm-officer roster injection (AC-3) and the bounded TERMINAL_TEARDOWN_BOUNDED marker (AC-4).
package ensigncycle

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// forceTeamModeCue makes the FO create a real team rather than taking the bare
// single-entity path. Unlike headless `-p` — where 2.1.178 dropped the native team
// tools and the SDK races to end_turn (so this cue cannot be honored) — the
// interactive pty session has the team tools present and stays resident, so the cue
// is honest and reliable here. It is GENERIC: it names the dispatch mode only, no
// stage or task. Resurrected verbatim from the retired forced-team tests (it was a
// shared const before those tests were retired in f1be9f07).
const forceTeamModeCue = "You MUST run in team mode for this run: create a team " +
	"(TeamCreate) before the first dispatch and dispatch every worker through that " +
	"team. Bare mode (no team) is NOT acceptable for this run — if you would otherwise " +
	"choose bare/single-entity mode, override that and use team mode anyway. "

// TestLivePtyStandingResidencyInjectsCommOfficer is AC-3: the live comm-officer
// roster-injection proof, resurrected over the pty transport (it could not run
// headless — the retired TestLiveStandingResidencyInjectsCommOfficer). It boots a
// REAL interactive `spacedock claude` against a fixture declaring a standing
// comm-officer mod, forces team mode, drives to the first team-mode dispatch (where
// the contract runs `spacedock dispatch spawn-standing-all`, injecting declared
// standing teammates), and asserts comm-officer LANDED in the team config.json
// members[] roster under the run's isolated team root.
//
// The independent oracle is the on-disk team config the live Claude Code harness
// writes, parsed with the SAME members[] shape claudeteam.MemberExists reads — NOT
// a contract grep. Spike-proven mechanism: an Agent-spawned teammate joined
// members[] in a second pane. Skips (never fatals) without auth (isolatedClaudeEnv,
// AC-6).
func TestLivePtyStandingResidencyInjectsCommOfficer(t *testing.T) {
	driver := newPtyLiveDriver(t)

	// Stage the realistic ≥3-stage lifecycle fixture PLUS a `_mods/comm-officer.md`
	// standing mod — the only difference from the plain cycle fixture. The standing
	// mod is what spawn-standing-all enumerates at first dispatch.
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeRealisticLifecycle())
	writeFile(t, filepath.Join(root, "_mods", "comm-officer.md"), ptyStandingCommOfficerMod())
	writeFile(t, filepath.Join(root, "make-it-work.md"), entityFixture())
	gitInit(t, root)

	// macOS t.TempDir() returns a `/var/folders/...` path while the FO's boot command
	// targets the EvalSymlinks-resolved `/private/var/folders/...` (the same
	// directory), so comparing the unresolved root would false-flag every local macOS
	// run as a wander. The CI Linux runner has no such symlink, so this is a no-op
	// there; resolving here keeps detectWrongRootBoot accurate on BOTH.
	rootResolved := root
	if r, err := filepath.EvalSymlinks(root); err == nil {
		rootResolved = r
	}

	prompt := forceTeamModeCue +
		"Use spacedock:first-officer for this whole run. Drive the workflow. " +
		antiShutdownOverride
	session := driver.launchAndSend(t, "standing-residency", root, prompt)
	defer session.proc.kill()

	// Watch the residency-relevant beats from the live session jsonl via the EXISTING
	// streamWatcher: TeamCreate engages team mode, then the first ensign dispatch
	// OPENS. The contract runs spawn-standing-all immediately BEFORE that first
	// team-mode Agent() dispatch, so by the time the dispatch OPENS the standing
	// teammate is already injected into the team config — exactly what the roster
	// check below verifies. The barrier is the dispatch OPEN (not its close): the
	// team-mode ensign completion can arrive as a `direct` message rather than the
	// task_notification close anchor, so waiting for close would flake; injection is
	// done at OPEN.
	watcher := newStreamWatcher(session.newFileSource(), session.proc, func(line string) { t.Log(line) })
	if _, err := watcher.expect(isTeamCreate, ptyBootBudget, "TeamCreate"); err != nil {
		if wrongRoot := detectWrongRootBoot(watcher.fullTranscript(), rootResolved); wrongRoot != nil {
			t.Fatalf("live residency drive failed at TeamCreate due to a wrong-root boot: %v\nUnderlying watcher error: %v\nArtifacts: %s", wrongRoot, err, session.artifactDir)
		}
		t.Fatalf("live residency drive failed at TeamCreate: %v\nFO pane:\n%s\nArtifacts: %s", err, driver.captureFOPane(session.tmuxName), session.artifactDir)
	}
	if _, err := watcher.expect(isEnsignDispatch, quietBudgetDispatchClose, "ensign dispatch open"); err != nil {
		t.Fatalf("live residency drive failed waiting for the first ensign dispatch to open (spawn-standing-all runs just before it): %v\nArtifacts: %s", err, session.artifactDir)
	}

	// The load-bearing assertion: comm-officer is a member of the live team's
	// config.json roster under the run's isolated team root. The team root is
	// {effective CLAUDE_CONFIG_DIR}/teams — where the comm-officer startup hook
	// membership-checks and TeamCreate writes the config.json, the same dir this
	// run's session jsonl lives under.
	teamRoot := filepath.Join(session.configDir, "teams")
	teams := ptyTeamConfigPaths(t, teamRoot)
	if len(teams) == 0 {
		t.Fatalf("no team config.json found under %s after the drive (TeamCreate matched but no roster persisted)\nArtifacts: %s", teamRoot, session.artifactDir)
	}
	if !ptyAnyTeamHasMember(t, teams, "comm-officer") {
		t.Fatalf("comm-officer ABSENT from every team roster after first dispatch — residency broke (standing-teammate injection did not land it in config.json)\nscanned: %v\nArtifacts: %s", teams, session.artifactDir)
	}
	t.Logf("comm-officer present in the live team roster across %d team config(s)", len(teams))
}

// TestLivePtyEnsignCycleTeamTeardown is AC-4: the bounded terminal-teardown marker
// proof, resurrected over the pty transport (the retired
// TestLiveEnsignCycleTeamTeardown). It boots a REAL interactive `spacedock claude`,
// forces team mode, drives a full TeamCreate → dispatch → terminalize → teardown
// cycle, and grades the BOUNDED best-effort terminal teardown by the EXISTING
// gradeTerminalTeardown / markerEmittedByAssistant over the captured live stream —
// the SAME grader the offline fixture suite uses, no parallel grading stack.
//
// The FIX's unique signal is the FO EMITTING the contract-mandated
// TERMINAL_TEARDOWN_BOUNDED marker (a text/thinking block it authors, not a
// contract-Read). Spike-proven mechanism: the live FO emitted the marker on an
// active-member TeamDelete refusal. Skips (never fatals) without auth (AC-6).
func TestLivePtyEnsignCycleTeamTeardown(t *testing.T) {
	driver := newPtyLiveDriver(t)

	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeRealisticLifecycle())
	writeFile(t, filepath.Join(root, "make-it-work.md"), entityFixture())
	gitInit(t, root)

	// Resolve the fixture root so the wrong-root comparison normalizes macOS's
	// /var->/private/var symlink (see the residency test for the full rationale); a
	// no-op on the CI Linux runner.
	rootResolved := root
	if r, err := filepath.EvalSymlinks(root); err == nil {
		rootResolved = r
	}

	// Force team mode, give the FO the conn so the gateless fixture drives to
	// terminal, and carry the anti-early-shutdown override (it fights the per-turn
	// teardown nag, which is exactly the team teardown path this test exercises).
	prompt := forceTeamModeCue +
		"Use spacedock:first-officer for this whole run. Drive the workflow to " +
		"completion; you have the conn to resolve gates from each stage report's " +
		"verdict (auto-approve). " +
		antiShutdownOverride
	session := driver.launchAndSend(t, "team-teardown", root, prompt)
	defer session.proc.kill()

	watcher := newStreamWatcher(session.newFileSource(), session.proc, func(line string) { t.Log(line) })

	// Step 1: the first ensign dispatch OPENS — the cycle progressed past boot into
	// dispatch (the team-mode close anchor is unreliable, so OPEN is the early beat;
	// Step 2's marker is the load-bearing proof the FULL cycle ran).
	if _, err := watcher.expect(isEnsignDispatch, quietBudgetDispatchClose, "ensign dispatch open"); err != nil {
		if wrongRoot := detectWrongRootBoot(watcher.fullTranscript(), rootResolved); wrongRoot != nil {
			t.Fatalf("live team-teardown drive failed waiting for the ensign dispatch to open due to a wrong-root boot: %v\nUnderlying watcher error: %v\nArtifacts: %s", wrongRoot, err, session.artifactDir)
		}
		t.Fatalf("live team-teardown drive failed waiting for the ensign dispatch to open: %v\nFO pane:\n%s\nArtifacts: %s", err, driver.captureFOPane(session.tmuxName), session.artifactDir)
	}

	// Step 2 (the team-only coverage this test exists for): grade the BOUNDED
	// best-effort terminal teardown by the marker the FO authors on cap-exhaustion.
	// expectTerminalTeardownGrade keys PASS on that emission — NOT on a clean
	// self-exit (impossible while members[] is populated) and NOT on a post-marker
	// HOLD the real FO cannot deliver. It is the live realization of
	// gradeTerminalTeardown — the SAME discriminator (markerEmittedByAssistant) the
	// offline fixture suite greens on, run here over the live session jsonl. The
	// no-progress quiet budget stays ≤60s (quietBudgetDefault) so the AC-1 "no
	// individual timeout > 60s" guard is unaffected — it resets on every drained
	// line, so an actively-streaming teardown never trips it.
	if err := watcher.expectTerminalTeardownGrade(quietBudgetDefault); err != nil {
		t.Fatalf("live team-teardown drive failed grading the terminal teardown: %v\nFO pane:\n%s\nArtifacts: %s", err, driver.captureFOPane(session.tmuxName), session.artifactDir)
	}

	// Belt-and-braces over the captured stream with the offline grader directly —
	// proves the EXACT offline gradeTerminalTeardown([]string) greens the live
	// transcript, not just its streaming realization expectTerminalTeardownGrade.
	if ok, reason := gradeTerminalTeardown(splitStreamLines(watcher.fullTranscript())); !ok {
		t.Fatalf("offline gradeTerminalTeardown REDs the live stream despite the streaming grade passing: %s\nArtifacts: %s", reason, session.artifactDir)
	}

	// The terminal end-state must still be PRESENT and CORRECT on disk — the marker
	// alone proves the FO reached the teardown terminus, not that the cycle
	// completed. A team FO archives BEFORE the teardown hold, so by marker time the
	// entity is terminalized + archived. The `verdict:` field is NOT asserted: team
	// finalize omits it non-deterministically (the spun-off verdict-omission task
	// owns that), and forcing team mode here is exactly where a verdict assertion
	// would flake.
	entity, where, found := locateEntity(root, "make-it-work")
	if !found {
		t.Fatalf("entity make-it-work not found in place or under _archive/ after the team-teardown drive\nArtifacts: %s", session.artifactDir)
	}
	t.Logf("located entity at %s", where)
	if !liveStageReportHeading.MatchString(entity) {
		t.Errorf("entity missing anchored stage-report heading\n%s", entity)
	}
	if !frontmatterField.MatchString(entity) {
		t.Errorf("entity missing terminal `status: done`\n%s", entity)
	}
	if !someCommitNamesOnly(t, root, "make-it-work") {
		t.Errorf("no path-scoped commit named only the entity in the team-teardown drive history")
	}
}

// splitStreamLines splits a newline-joined transcript back into the []string the
// offline gradeTerminalTeardown consumes.
func splitStreamLines(stream string) []string {
	if stream == "" {
		return nil
	}
	return strings.Split(stream, "\n")
}

// ptyStandingCommOfficerMod is a minimal standing comm-officer mod for the live
// residency fixture: standing: true frontmatter plus the spawn-config bullets the
// binary parses and an ## Agent Prompt body. It mirrors the shipped mod's
// self-declaration shape without the full routing prose, which is irrelevant to the
// roster-membership proof. Resurrected from the retired live residency test.
func ptyStandingCommOfficerMod() string {
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

// ptyTeamConfigPaths returns every {teamRoot}/*/config.json path on disk. The live
// harness writes one per team; the FO creates exactly one, but the scan is
// plural-safe so a parallel-team layout does not silently miss the roster.
func ptyTeamConfigPaths(t *testing.T, teamRoot string) []string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(teamRoot, "*", "config.json"))
	if err != nil {
		t.Fatalf("glob team configs under %s: %v", teamRoot, err)
	}
	return matches
}

// ptyAnyTeamHasMember reports whether any of the given team config.json files lists
// a member with the given name. It parses the same members[] shape
// claudeteam.MemberExists reads, so the live oracle matches the binary's own
// membership predicate.
func ptyAnyTeamHasMember(t *testing.T, configPaths []string, name string) bool {
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
