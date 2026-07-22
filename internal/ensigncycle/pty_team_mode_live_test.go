//go:build live

// ABOUTME: Live team-mode e2e over the pty/tmux driver — resurrects the two retired forced-team
// ABOUTME: tests: comm-officer roster injection (AC-3) and team-mode dispatch-to-terminal (AC-4).
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
	// Legacy lane: this test drives native TeamCreate on an interactive session. On
	// a merged host (claude ≥2.1.178, native team tools gone) it cannot run, so SKIP
	// rather than RED when CI unpins past the merged floor (the merged lane owns the
	// in-process named-background-teammate coverage).
	skipUnlessTeamCreateCapable(t)
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

	// The boot prompt: name the run and the intent. An interactive FO greets and
	// PARKS for captain input after boot by contract (item 9 — correct behavior), so
	// this first message does NOT auto-dispatch; the captain nudge below pushes it
	// past the greet into the actual dispatch (which runs spawn-standing-all, the
	// residency injection under test).
	prompt := forceTeamModeCue +
		"Use spacedock:first-officer for this whole run. " +
		antiShutdownOverride
	session := driver.launchAndSend(t, "standing-residency", root, prompt)
	defer session.proc.kill()

	// Drive the FO past its contractual greet-stop like a captain: nudge with a "go +
	// conn" message and wait for it to actually create the team. nudgePastGreet is
	// bounded (≤3 nudges) and keys on the on-disk transcript carrying TeamCreate, then
	// the EXISTING expect(isTeamCreate) + roster assertions run unchanged.
	commOfficerNudge := "Yes — you have the conn for this run: dispatch make-it-work " +
		"into implementation now, create the team, and drive to completion, approving " +
		"gates yourself. Don't stop to ask."
	if err := driver.nudgePastGreet(session, commOfficerNudge, linesHaveTeamCreate, 3, ptyBootBudget); err != nil {
		t.Fatalf("live residency drive: the FO never created a team across captain nudges: %v\nFO pane:\n%s\nArtifacts: %s",
			err, driver.captureFOPane(session.tmuxName), session.artifactDir)
	}

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

// TestLivePtyEnsignCycleTeamTeardown is AC-4: the live proof that a real team-mode
// dispatch DRIVES THE ENTITY TO A TERMINAL STATE, resurrected over the pty transport
// (the retired TestLiveEnsignCycleTeamTeardown). It boots a REAL interactive
// `spacedock claude`, forces team mode, drives a full TeamCreate → dispatch →
// terminalize cycle, and asserts the clean dispatch-to-terminal terminus: the FO
// drove the entity to its terminal/archived state (status: done + a stage-report + a
// path-scoped commit). In the interactive/tmux world the launcher just kills the
// session once the work is finished, so a clean terminalize is the valid terminus.
//
// The load-bearing behavior the live test proves is that the DISPATCH MECHANISM
// reaches a gate/terminal state. (The legacy cap-exhaustion TERMINAL_TEARDOWN_BOUNDED
// marker terminus was retired with the legacy bounded TeamDelete teardown; it was a
// headless-`-p` edge-case fallback for a stuck-roster process that cannot self-exit,
// and the interactive launcher's clean-terminal happy path never needed it.) Skips
// (never fatals) without auth (AC-6).
func TestLivePtyEnsignCycleTeamTeardown(t *testing.T) {
	// Legacy lane: drives native TeamCreate on an interactive session. On a merged
	// host (claude ≥2.1.178, native team tools gone) it cannot run, so SKIP rather
	// than RED when CI unpins past the merged floor.
	skipUnlessTeamCreateCapable(t)
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

	// The boot prompt names the run and force-team mode; the interactive FO greets and
	// PARKS by contract (item 9), so the captain nudge below pushes it past the greet.
	// The anti-early-shutdown override fights the per-turn teardown nag — exactly the
	// team teardown path this test exercises.
	prompt := forceTeamModeCue +
		"Use spacedock:first-officer for this whole run. " +
		antiShutdownOverride
	session := driver.launchAndSend(t, "team-teardown", root, prompt)
	defer session.proc.kill()

	// Drive the FO past its contractual greet-stop like a captain: nudge with a "go +
	// conn, auto-approve gates" message until it actually dispatches the ensign (the
	// first beat of the TeamCreate→dispatch→terminalize→teardown cycle). Bounded
	// (≤3 nudges); the EXISTING expect(isEnsignDispatch) + teardown grade run unchanged
	// after. If the FO re-parks at a gate, the same nudge ("approving gates yourself")
	// covers it.
	teardownNudge := "Yes — you have the conn for this run: dispatch make-it-work into " +
		"implementation now, create the team, drive to completion, and approve gates " +
		"yourself from each stage report's verdict. Don't stop to ask."
	if err := driver.nudgePastGreet(session, teardownNudge, linesHaveEnsignDispatch, 3, quietBudgetDispatchClose); err != nil {
		t.Fatalf("live team-teardown drive: the FO never dispatched an ensign across captain nudges: %v\nFO pane:\n%s\nArtifacts: %s",
			err, driver.captureFOPane(session.tmuxName), session.artifactDir)
	}

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

	// Step 2 (the team-only coverage this test exists for): wait for the dispatch to
	// drive the entity to its terminal state, then assert that terminus is real and
	// correct.
	//
	// The terminus proving the dispatch drove the entity to a terminal state is the
	// clean dispatch-to-terminal happy path: the on-disk entity reached its
	// terminal/archived state (status: done + a stage-report). A team FO archives the
	// entity as part of finalize, so this is observable on disk. (The legacy
	// cap-exhaustion TERMINAL_TEARDOWN_BOUNDED marker terminus was retired with the
	// legacy bounded TeamDelete teardown; the interactive launcher just kills the
	// session on clean finalize, which is the interactive happy path anyway.)
	//
	// expectCondition drains the stream each poll (keeping the watcher transcript and
	// the no-progress budget fresh) and returns once the entity is terminal. The
	// dispatch-close budget bounds STREAM SILENCE, not wallclock — it resets on every
	// drained line, so a multi-turn finalize that keeps streaming never trips it.
	entityTerminal := func() bool {
		entity, _, found := locateEntity(root, "make-it-work")
		return found && liveStageReportHeading.MatchString(entity) && frontmatterField.MatchString(entity)
	}
	if err := watcher.expectCondition(entityTerminal, quietBudgetDispatchClose, "team-mode terminus"); err != nil {
		if wrongRoot := detectWrongRootBoot(watcher.fullTranscript(), rootResolved); wrongRoot != nil {
			t.Fatalf("live team-teardown drive reached no terminus due to a wrong-root boot: %v\nUnderlying watcher error: %v\nArtifacts: %s", wrongRoot, err, session.artifactDir)
		}
		t.Fatalf("live team-teardown drive reached no clean terminal entity state: %v\nFO pane:\n%s\nArtifacts: %s", err, driver.captureFOPane(session.tmuxName), session.artifactDir)
	}

	// Whichever terminus the FO took, the entity must be in its terminal/archived
	// end-state on disk — the load-bearing proof that the dispatch DROVE THE ENTITY TO
	// A TERMINAL STATE (the assertion this test exists for). A team FO archives the
	// entity during finalize, so by terminus time it is terminalized + archived under
	// either path. The `verdict:` field is NOT asserted: team finalize omits it
	// non-deterministically (the spun-off verdict-omission task owns that; verdict is
	// NOT a mode-invariant terminal fact), and forcing team mode here is exactly where
	// a verdict assertion would flake — so this gates on the mode-invariant terminal
	// facts (status: done + a stage-report + a path-scoped commit).
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
// line-oriented offline graders consume.
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

// linesHaveTeamCreate reports whether any raw transcript line is a TeamCreate
// tool_use — the signalSeen predicate the residency nudge keys on (the FO created
// the team, so spawn-standing-all has run). It reuses the EXISTING isTeamCreate
// streamEntry predicate over each parsed line.
func linesHaveTeamCreate(lines []string) bool {
	return linesMatch(lines, isTeamCreate)
}

// linesHaveEnsignDispatch reports whether any raw transcript line is an
// Agent(subagent_type="spacedock:ensign") tool_use — the signalSeen predicate the
// teardown nudge keys on (the cycle progressed past boot into dispatch). It reuses
// the EXISTING isEnsignDispatch streamEntry predicate over each parsed line.
func linesHaveEnsignDispatch(lines []string) bool {
	return linesMatch(lines, isEnsignDispatch)
}

// linesMatch parses each raw stream-json line into a streamEntry and reports
// whether any satisfies the given streamEntry predicate (the same predicates the
// streamWatcher's expect() uses), so the nudge loop and the watcher key on the
// identical signal.
func linesMatch(lines []string, predicate func(streamEntry) bool) bool {
	for _, line := range lines {
		var e streamEntry
		if json.Unmarshal([]byte(line), &e) != nil {
			continue
		}
		if predicate(e) {
			return true
		}
	}
	return false
}
