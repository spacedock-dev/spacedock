//go:build live

// ABOUTME: Live merged-lane team-mode e2e — a headless `claude -p` FO on the current (merged) host
// ABOUTME: dispatches a named background Agent (no TeamCreate), asserting the merged on-disk shape.
package ensigncycle

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestLiveMergedTeamModeDispatch is the merged-lane live proof: on a current
// (merged, ≥2.1.178) Claude host, a HEADLESS `claude -p` FO drives the workflow as
// a named BACKGROUND teammate dispatch — NO TeamCreate, NO tmux — exactly the
// in-process team shape #396 (`using-claude-team` merged support) committed. It is
// the second live lane alongside the legacy interactive pty lane (which SKIPs on
// this host); together they realize 9243's OQ-6 host-regime split.
//
// It reuses the headless launch wiring of TestLiveEnsignCycle verbatim (the same
// `spacedock claude -- -p` front door, isolatedClaudeEnv auth/HOME isolation,
// streamWatcher liveness, detectWrongRootBoot, locateEntity terminal-state grade)
// and ADDS only the merged-shape assertions over the captured stream + the on-disk
// in-process member meta + a reconcile exit-0 check. No scenario/fixture/assertion
// is forked.
//
// GROUND-TRUTH SCOPE (verified by a live probe on local 2.1.181, corroborated by
// the team-lead's own earlier headless `-p` trials, recorded in the stage report):
// a headless `claude -p` merged host does NOT write a
// `~/.claude/teams/session-<id>/config.json` auto-team registry. That registry is
// the INTERACTIVE / TeamCreate-era artifact 9243 observed — it does not materialize
// for an in-process named-background `Agent(run_in_background=true)` dispatch under
// `-p`. The member record Claude Code DOES write headless is
// `projects/<cwd>/<session-id>/subagents/agent-*.meta.json`, carrying the same
// agentType="spacedock:ensign" + name + no-team_name shape. CONSEQUENCES, hence the
// scoping (deliberate, NOT a coverage gap):
//   - the on-disk member oracle reads that subagents meta (the real path) — #3;
//   - reconcile (no --team-name) is asserted to EXIT 0 only (the command contract
//     runs) — #4. Because no `teams/config.json` lands headless, the leadSessionId
//     auto-discovery has nothing on disk to match, so reconcile degrades to git-only
//     (team_name="", empty drift); a non-empty resolved roster is NOT asserted;
//   - the per-name shutdown_request teardown + members[] prune are interactive-only
//     beats and are NOT asserted here (there is no `teams/config.json` members[] to
//     prune) — #5; the kept terminus is the entity reaching its terminal state,
//     identical to TestLiveEnsignCycle.
//
// Skips (never fatals) without auth (isolatedClaudeEnv, the same AC-6 gate).
func TestLiveMergedTeamModeDispatch(t *testing.T) {
	// Merged lane: this asserts the in-process named-background shape. On a LEGACY
	// host (claude <2.1.178, native TeamCreate present) the FO drives the native team
	// registry instead, so SKIP there — the legacy interactive pty lane covers that
	// host. The mirror of skipUnlessTeamCreateCapable: exactly one team lane runs per
	// pinned/unpinned claude.
	skipUnlessMergedHost(t)
	binary := spacedockBinary(t)
	pluginDir := livePluginDir(t)
	model := envOr("SPACEDOCK_LIVE_MODEL", "sonnet")

	// Same isolated auth/HOME tree the headless suite uses (OAuth benchmark-token /
	// ANTHROPIC_API_KEY, else SKIP), plus the built binary first on PATH so the FO's
	// `spacedock --version` resolves the test binary.
	env := isolatedClaudeEnv(t, os.Getenv("HOME"))
	env = withBinaryOnPath(env, binary)
	configDir, _ := envValue(env, "CLAUDE_CONFIG_DIR")
	homeDir, _ := envValue(env, "HOME")
	effectiveConfigDir := configDir
	if effectiveConfigDir == "" {
		effectiveConfigDir = filepath.Join(homeDir, ".claude")
	}

	// The realistic ≥3-stage lifecycle fixture (backlog → implementation → done), a
	// flat entity at backlog — the SAME fixture the bare cycle and the pty lane use.
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeRealisticLifecycle())
	writeFile(t, filepath.Join(root, "make-it-work.md"), entityFixture())
	gitInit(t, root)

	// macOS t.TempDir() returns /var/folders/... while the FO boots the
	// EvalSymlinks-resolved /private/var/folders/... (same dir); resolve for the
	// wrong-root comparison AND the subagents-meta path (Claude Code encodes the
	// symlink-resolved cwd in the projects-dir name). No-op on CI Linux.
	resolvedRoot := root
	if r, err := filepath.EvalSymlinks(root); err == nil {
		resolvedRoot = r
	}

	// Force team mode: on a merged host "team mode" IS the named-background-teammate
	// dispatch (no TeamCreate). Without this cue the headless FO chooses the bare
	// sequential path (does the stage work inline) — verified live — so the cue makes
	// the FO take the merged dispatch shape this lane exists to cover. It is GENERIC
	// (names the dispatch mode only, no stage/task) and resilient on a merged host
	// where the named background teammate is exactly how the FO works concurrently.
	forceMergedTeamCue := "You MUST run in team mode for this run: dispatch every " +
		"worker as a named background teammate (an Agent with a name set and " +
		"run_in_background true) and do NOT do the stage work inline yourself. " +
		"Bare/sequential mode (doing the stage work yourself without dispatching a " +
		"named teammate) is NOT acceptable for this run. "
	drivePrompt := forceMergedTeamCue +
		"Drive the workflow to completion; you have the conn to resolve gates from " +
		"each stage report's verdict (auto-approve). " + antiShutdownOverride

	// The headless front door — identical to TestLiveEnsignCycle's launch — wrapped
	// in `env -u <nested-session markers>` so a child claude launched from inside a
	// Claude session (the common CI/teammate case) does not self-identify as nested
	// and suppress its transcript (CC 2.1.170+). The team flag is NOT set: the merged
	// channel is flag-free headless (verified) and the legacy native team tools it
	// would expose are gone on this host anyway. stdout+stderr fold into one pipe the
	// watcher drains line-by-line.
	launchArgs := unsetNestedSessionArgs(binary, "claude",
		"--plugin-dir", pluginDir,
		"--skip-compat-check",
		"--",
		"-p", drivePrompt,
		"--permission-mode", "bypassPermissions",
		"--output-format", "stream-json",
		"--verbose",
		"--model", model,
		"Use $spacedock:first-officer for this whole run.",
	)
	cmd := exec.Command("env", launchArgs...)
	cmd.Dir = root
	cmd.Env = env

	artifactDir := claudeLiveArtifactDir(t, filepath.Join("merged-team-mode", "dispatch"))

	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw
	started := time.Now()
	if err := cmd.Start(); err != nil {
		t.Fatalf("merged-lane spacedock claude failed to start: %v", err)
	}
	poller := newCmdPoller(cmd, pw)
	defer poller.kill()
	watcher := newStreamWatcher(newPipeLineSource(pr), poller, func(line string) { t.Log(line) })

	// Drive to the load-bearing END-STATE, then kill — NOT to FO subprocess exit.
	// A headless `-p` FO with the anti-early-shutdown override keeps emitting
	// post-completion turns (the terminal merge ceremony, cooperative worker
	// shutdown, a post-merge reconcile, idle-loop polling) AFTER the entity is
	// archived, so waiting for the process to EXIT can outrun the package timeout on
	// a slower run even though the deliverable was reached early (observed: one run
	// finished in 338s, a fuller-ceremony run crossed 600s still emitting). The
	// deliverable is "the merged FO dispatched a named background ensign and the
	// entity reached its terminal state" — so gate on THAT (the same expect-then-kill
	// discipline TestLiveEnsignCycle uses), then reap the still-chatty subprocess and
	// grade the accumulated transcript. Each step is bounded by the no-progress quiet
	// budget; the deferred poller.kill() reaps on every exit path.
	//
	// Step 1: the first merged ensign dispatch OPENS (the FO drove past boot into
	// dispatch — an early fail-fast that a greet-stop or wrong-root never reaches).
	if _, err := watcher.expect(isEnsignDispatch, quietBudgetDispatchClose, "merged ensign dispatch open"); err != nil {
		stream := watcher.fullTranscript()
		_ = os.WriteFile(filepath.Join(artifactDir, "merged-stream.jsonl"), []byte(stream), 0o644)
		if wrongRoot := detectWrongRootBoot(stream, resolvedRoot); wrongRoot != nil {
			t.Fatalf("merged-lane drive failed waiting for the ensign dispatch to open due to a wrong-root boot: %v\nUnderlying watcher error: %v\nArtifacts: %s", wrongRoot, err, artifactDir)
		}
		if _, extractErr := extractClaudeFinalMessage(stream); extractErr != nil {
			t.Fatalf("merged-lane launch failed (credential/launch error, not a behavior failure): %v\nArtifacts: %s\nStream tail:\n%s", extractErr, artifactDir, tail(stream, 4000))
		}
		t.Fatalf("merged-lane drive: the FO never dispatched a merged ensign (it fell to bare/sequential or stalled at boot): %v\nArtifacts: %s", err, artifactDir)
	}

	// Step 2: the entity reaches its on-disk terminal end-state (the mode-invariant
	// completion proof). expectCondition drains the stream each poll (liveness; the
	// budget resets on activity) while checking the filesystem.
	entityTerminal := func() bool {
		body, _, found := locateEntity(root, "make-it-work")
		return found && frontmatterField.MatchString(body) && someCommitNamesOnly(t, root, "make-it-work")
	}
	if err := watcher.expectCondition(entityTerminal, quietBudgetDispatchClose, "entity terminalized"); err != nil {
		stream := watcher.fullTranscript()
		_ = os.WriteFile(filepath.Join(artifactDir, "merged-stream.jsonl"), []byte(stream), 0o644)
		t.Fatalf("merged-lane drive failed waiting for the entity to terminalize+commit (status: done + path-scoped commit): %v\nArtifacts: %s", err, artifactDir)
	}

	// The end-state is reached; reap the still-running (post-completion-chatty)
	// subprocess and grade the accumulated transcript. The kill is what bounds
	// wallclock — the assertions below run on what the FO already emitted through the
	// terminal state, never waiting on its idle loop.
	poller.kill()
	stream := watcher.fullTranscript()
	duration := time.Since(started)
	t.Logf("merged-lane drive reached terminal state in %s, %d transcript lines", duration.Round(time.Second), len(strings.Split(stream, "\n")))
	_ = os.WriteFile(filepath.Join(artifactDir, "merged-stream.jsonl"), []byte(stream), 0o644)

	// A 401/is_error result is a LOUD launch failure, never fed into an assertion.
	// (Reached terminal state implies a real run, but keep the explicit guard so a
	// degenerate transcript fails loudly as a launch issue, not a behavior one.)
	if _, extractErr := extractClaudeFinalMessage(stream); extractErr != nil {
		t.Fatalf("merged-lane launch failed (credential/launch error, not a behavior failure): %v\nArtifacts: %s\nStream tail:\n%s",
			extractErr, artifactDir, tail(stream, 4000))
	}

	lines := strings.Split(stream, "\n")

	// Assertion #1 — NO TeamCreate. Two independent signals, both verified live:
	//   (a) the init-event tool surface carries SendMessage but NOT TeamCreate/
	//       TeamDelete (the static analog of the contract's ToolSearch(select:
	//       TeamCreate)-empty discriminator), and
	//   (b) the FO emitted no TeamCreate/TeamDelete tool_use anywhere in the stream.
	initTools := initEventToolNames(lines)
	if len(initTools) == 0 {
		t.Fatalf("no system/init event in the merged stream — cannot read the tool surface\nArtifacts: %s", artifactDir)
	}
	if stringInSlice("TeamCreate", initTools) || stringInSlice("TeamDelete", initTools) {
		t.Errorf("merged host init tool surface unexpectedly contains TeamCreate/TeamDelete — not a merged host\ninit tools: %v\nArtifacts: %s", initTools, artifactDir)
	}
	if !stringInSlice("SendMessage", initTools) {
		t.Errorf("merged host init tool surface missing SendMessage (the inter-agent communication tool)\ninit tools: %v\nArtifacts: %s", initTools, artifactDir)
	}
	if streamHasTeamCreateToolUse(lines) {
		t.Errorf("the merged FO emitted a TeamCreate/TeamDelete tool_use — it must use the named-background dispatch, not the native team registry\nArtifacts: %s", artifactDir)
	}

	// Assertion #2 — at least one merged-shape Agent dispatch: name set,
	// run_in_background true, NO team_name. This is the load-bearing proof the FO did
	// NOT fall to bare/sequential (a bare run emits no Agent dispatch at all).
	dispatches := mergedEnsignDispatches(lines)
	if len(dispatches) == 0 {
		t.Fatalf("the merged FO dispatched NO Agent(subagent_type=spacedock:ensign) — it fell to bare/sequential, not the merged named-background team shape\nArtifacts: %s", artifactDir)
	}
	mergedShaped := false
	for _, d := range dispatches {
		if d.name != "" && d.runInBackground && !d.hasTeamName {
			mergedShaped = true
			break
		}
	}
	if !mergedShaped {
		t.Errorf("no Agent dispatch had the merged shape (name set + run_in_background true + no team_name); dispatches=%+v\nArtifacts: %s", dispatches, artifactDir)
	}

	// Assertion #6 (completion target) — the FO forwarded a dispatch prompt whose
	// completion-signal block pins SendMessage(to="team-lead"), per #396's committed
	// build.go (NOT "main"). The dispatch prompt is a pointer the ensign Reads, so the
	// concrete completion target lives in the dispatch FILE under /tmp/spacedock-dispatch;
	// assert it there (the build helper wrote it during the FO's dispatch).
	if !mergedDispatchFileHasEnsignContract(t, dispatches, artifactDir) {
		t.Errorf("no merged dispatch file carried the ensign skill and pinned SendMessage(to=\"team-lead\")\nArtifacts: %s", artifactDir)
	}

	// Capture the FO's own session id (the init event's session_id == the
	// $CLAUDE_CODE_SESSION_ID Claude Code set, the leadSessionId reconcile matches).
	sessionID := initEventSessionID(lines)
	if sessionID == "" {
		t.Fatalf("no session_id in the merged init event — cannot locate the member meta or run reconcile by session\nArtifacts: %s", artifactDir)
	}

	// Assertion #3 — the on-disk in-process member record. Claude Code writes the
	// merged background teammate's meta at
	// {configDir}/projects/{encode(cwd)}/{sessionID}/subagents/agent-*.meta.json,
	// carrying agentType="spacedock:ensign" + the member name + NO team_name — the
	// merged-host analog of a legacy team config.json members[] entry (which is NOT
	// written headless). This is the independent on-disk oracle (read straight off
	// disk, not from the stream the test also grades).
	metas, err := readMergedMemberMetas(effectiveConfigDir, resolvedRoot, sessionID)
	if err != nil {
		t.Fatalf("reading merged member metas: %v\nArtifacts: %s", err, artifactDir)
	}
	metaPath := mergedMemberMetasPath(effectiveConfigDir, resolvedRoot, sessionID)
	if !hasMergedEnsignMember(metas) {
		t.Errorf("no in-process member meta with agentType=spacedock:ensign + a name + no team_name found under %s; metas=%+v\nArtifacts: %s",
			metaPath, metas, artifactDir)
	}

	// Assertion #4 — reconcile (no --team-name) auto-discovers by session id and
	// EXITS 0. Run the SAME binary with the SAME isolated HOME/CLAUDE_CONFIG_DIR and
	// the FO's session id in CLAUDE_CODE_SESSION_ID, so the leadSessionId
	// auto-discovery path runs under the FO's own identity. On a headless merged host
	// no teams/config.json lands, so reconcile degrades to git-only (team_name="",
	// empty drift) — a NON-empty resolved roster is NOT asserted; the command-contract
	// proof is the clean exit 0 (the sweep ran, auto-discovery was exercised).
	reconcileEnv := withClaudeConfigDir(env, effectiveConfigDir)
	reconcileEnv = upsertEnv(reconcileEnv, "CLAUDE_CODE_SESSION_ID", sessionID)
	reconcileEnv = upsertEnv(reconcileEnv, "HOME", homeDir)
	var rout, rerr bytes.Buffer
	reconcile := exec.Command(binary, "dispatch", "reconcile", "--workflow-dir", root)
	reconcile.Dir = root
	reconcile.Env = reconcileEnv
	reconcile.Stdout = &rout
	reconcile.Stderr = &rerr
	if err := reconcile.Run(); err != nil {
		t.Errorf("`spacedock dispatch reconcile` (no --team-name) did not exit 0: %v\nstdout:\n%s\nstderr:\n%s\nArtifacts: %s",
			err, rout.String(), rerr.String(), artifactDir)
	} else {
		t.Logf("reconcile (no --team-name) exit 0; stdout: %s", strings.TrimSpace(rout.String()))
	}

	// Kept terminus (the AC-4-style proof the dispatch DROVE THE ENTITY TO A TERMINAL
	// STATE) — the same mode-invariant on-disk facts TestLiveEnsignCycle asserts:
	// entity locatable, status: done, a path-scoped commit. verdict is NOT asserted
	// (team-finalize omits it non-deterministically — the spun-off verdict-omission
	// task owns that).
	entity, where, found := locateEntity(root, "make-it-work")
	if !found {
		t.Fatalf("entity make-it-work not found in place or under _archive/ after the merged drive\nArtifacts: %s", artifactDir)
	}
	t.Logf("located entity at %s", where)
	if !liveStageReportHeading.MatchString(entity) {
		t.Errorf("entity missing anchored stage-report heading\n%s", entity)
	}
	if !frontmatterField.MatchString(entity) {
		t.Errorf("entity missing terminal `status: done`\n%s", entity)
	}
	if !someCommitNamesOnly(t, root, "make-it-work") {
		t.Errorf("no path-scoped commit named only the entity in the merged drive history")
	}
}

// mergedDispatchFileHasEnsignContract reports whether any merged dispatch the FO
// made wrote a dispatch file under /tmp/spacedock-dispatch that invokes the ensign
// skill and whose completion-signal block pins SendMessage(to="team-lead"). The
// dispatch file is keyed on the session id +
// the derived member name (build.go's mergedMode path), so it is located by the
// dispatched member's name. The completion target is the assertable code surface
// (#396's build.go emits to="team-lead"); a `claude-fo-dispatch.md` prose note still
// frames it as "main", but the code target is the one this asserts.
func mergedDispatchFileHasEnsignContract(t *testing.T, dispatches []mergedAgentDispatch, artifactDir string) bool {
	t.Helper()
	const dispatchDir = "/tmp/spacedock-dispatch"
	for _, d := range dispatches {
		if d.name == "" {
			continue
		}
		// The merged dispatch filename is {sessionToken}-{derivedName}.md or, when no
		// session token resolves, {derivedName}.md — both END with the derived name +
		// .md, so a suffix glob finds either.
		matches, _ := filepath.Glob(filepath.Join(dispatchDir, "*"+d.name+".md"))
		for _, p := range matches {
			raw, err := os.ReadFile(p)
			if err != nil {
				continue
			}
			if mergedDispatchArtifactHasEnsignContract(string(raw)) {
				t.Logf("merged dispatch file %s carries the ensign skill and pins completion target to team-lead", p)
				return true
			}
		}
	}
	return false
}

// upsertEnv returns env with key set to value — replacing an existing entry or
// appending a new one. The merged-lane reconcile uses it to ride the FO's session
// id (and the isolated HOME) into the reconcile subprocess so auto-discovery runs
// under the FO's own leadSessionId identity.
func upsertEnv(env []string, key, value string) []string {
	prefix := key + "="
	out := make([]string, 0, len(env)+1)
	replaced := false
	for _, e := range env {
		if strings.HasPrefix(e, prefix) {
			out = append(out, prefix+value)
			replaced = true
			continue
		}
		out = append(out, e)
	}
	if !replaced {
		out = append(out, prefix+value)
	}
	return out
}
