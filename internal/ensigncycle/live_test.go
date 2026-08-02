//go:build live

// ABOUTME: Live BEHAVIORAL test of the dispatch->ensign->stage cycle driven by a
// ABOUTME: REAL model through the spacedock claude front door (gated, -tags live only).
package ensigncycle

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/testgit"
)

// This is the live analog of TestEnsignCycleMechanicalOutputs (cycle_test.go).
// Where the skeleton stubs the LLM with a scripted Go ensign that works ONE stage
// in place, this test shells the v1 binary's REAL front door so an actual model
// drives the whole dispatch->ensign->stage protocol all the way to the terminal
// `done` stage. A real full-cycle FO makes MULTIPLE commits and ARCHIVES the
// terminal entity to `_archive/`, so the assertions target the REAL
// completed-and-archived end-state: the entity (located in place OR under
// `_archive/`) carries the anchored stage-report shape
// (liveStageReportHeading, doneMarker, `### Summary`, NOT checkboxBullet) AND the
// FO's terminal frontmatter (`status: done`, `verdict: passed`), and SOME commit
// in the history is path-scoped to the entity. The heading regex is stage-agnostic
// (the real cycle finishes at the TERMINAL stage, so the ensign writes
// `## Stage Report: done`); the remaining regexes are reused verbatim from the
// skeleton and the producer is the real runtime.
//
// The `//go:build live` tag keeps this out of the default `go test ./...` suite
// (the secret-free offline job). It compiles and runs ONLY under
// `go test -tags live`, the gated job's invocation that spends the live
// credential behind the CI-E2E approval gate.

// TestLiveEnsignCycle stages the skeleton's flat-entity backlog fixture, shells
// the real `spacedock claude` front door headless to drive the entity to done
// through a live model, then reads back the entity + git log and asserts the
// anchored mechanical contract. It is the smallest meaningful live mechanism
// proof: real binary front door + real plugin load + real model + real
// dispatch->ensign->stage cycle + a real path-scoped state commit.
//
// Auth + HOME isolation are resolved by isolatedClaudeEnv: an operator machine
// authenticates via the OAuth benchmark-token (~/.claude/benchmark-token), the
// CI runner via ANTHROPIC_API_KEY, and a machine with neither SKIPS (never
// fatals). The chosen credential runs against a fresh empty HOME so parallel
// `spacedock claude` invocations never collide in ~/.claude.
func TestLiveEnsignCycle(t *testing.T) {
	// The conn-cue drive prompt: the FO is given the conn to resolve gates from
	// each stage report's verdict (auto-approve) so the gateless realistic-lifecycle
	// fixture drives all the way to terminal `done`. The team-vs-bare mode is left
	// to the FO — the assertion below is TEAM-AGNOSTIC, so either choice passes.
	drivePrompt := "Drive the workflow to completion; you have the conn to resolve gates from each stage report's verdict (auto-approve). " + antiShutdownOverride
	watcher, root := startRealisticLifecycleDrive(t, drivePrompt)

	// The watch sequence asserts the dispatch→done INVARIANT, TEAM-AGNOSTICALLY —
	// it must green whether the headless `-p` FO drove TEAM or BARE (the captain's
	// determination sanctions a bare drive under `-p`; the team/bare choice is a
	// robustness detail of the dispatch mechanism, orthogonal to the cycle
	// completing). Two team-INDEPENDENT steps gate the smoke:
	//   1. the first ensign dispatch OPENS — the FO drove past boot into dispatch
	//      (an early fail-fast: a greet-stop or wrong-root never dispatches);
	//   2. the entity reaches its FULL on-disk TERMINAL end-state (`status: done`
	//      + a path-scoped commit) — the FO terminalized, archived, and committed.
	// Both hold in team AND bare mode. The barrier is the dispatch OPEN, NOT its
	// close: in current Claude Code the team-mode ensign completion arrives as a
	// `direct` message, not the `task_notification status=completed` anchor
	// expectDispatchClose keys on, so waiting for the close would FLAKE in team mode
	// (a healthy run leaves the dispatch "open" by that anchor's reckoning even after
	// the ensign finished). Step 2's terminalized end-state is the real completion
	// proof anyway — it STRICTLY implies the cycle ran past dispatch — so the open is
	// a sufficient early beat and the on-disk end-state is the load-bearing one.
	//
	// The team-only terminal-teardown MARKER (it fires only in team teardown) is NOT
	// gated here: headless `-p` goes bare, so a legitimate drive emits no marker and
	// gating on it would red an otherwise-correct cycle. The marker's offline coverage
	// is the fixture watcher suite (teardown_grade_watcher_test.go over the
	// sonnet_teamdelete_*.jsonl fixtures); the LIVE team end-to-end marker (and the
	// comm-officer roster injection) is owned by the pty/tmux harness
	// (pty_team_mode_live_test.go), which drives a real interactive session where
	// team mode is exposed. Each step is bounded by its own no-progress quiet budget;
	// a stalled step fails FAST and LOCALIZED.
	//
	// KNOWN GAP (conscious, by design — NOT a missing timeout): the quiet budget
	// resets on ANY drained line, so it catches a SILENT hang (no new stream
	// activity → tripped in ≤60s, localized) but NOT a "stuck-but-emitting"
	// (chatty) hang — an FO that keeps emitting unrelated stream lines without
	// ever reaching the next watched step. That case is deliberately left to
	// Go's BUILT-IN default test timeout rather than an explicit long `-timeout`:
	// the captain banned long individual timeouts, and a chatty hang is far rarer
	// than a silent one (which the quiet budget does catch). Documenting it here
	// keeps the trade-off explicit instead of silently absent.
	if _, err := watcher.expect(isEnsignDispatch, quietBudgetDispatchClose, "ensign dispatch open"); err != nil {
		// A dispatch that never opens is opaque on its own. The most common cause is
		// a wrong-root boot: a CI env leak lures the FO off `root` into the real repo,
		// it boots that workflow, finds nothing dispatchable, and greets-and-stops —
		// so it never dispatches and dispatch-open is exactly where it now fails.
		// Surface that explicitly — naming the expected fixture root vs the
		// wandered-to path — so the leak fails legibly instead of as a confusing
		// timeout.
		if wrongRoot := detectWrongRootBoot(watcher.fullTranscript(), root); wrongRoot != nil {
			t.Fatalf("live cycle failed waiting for the ensign dispatch to open due to a wrong-root boot: %v\nUnderlying watcher error: %v", wrongRoot, err)
		}
		t.Fatalf("live cycle failed waiting for the ensign dispatch to open: %v", err)
	}

	// Wait, TEAM-AGNOSTICALLY, for the entity to reach its on-disk terminal
	// end-state. This is the team-independent barrier that replaces the team-only
	// teardown-marker grade: the FO terminalizes (implementation→done) and
	// path-scoped-commits the entity in BOTH modes. The barrier requires the
	// MODE-INVARIANT durable facts only — entity locatable, `status: done`, AND a
	// path-scoped commit. It deliberately does NOT require a set `verdict:`: a live
	// `verdict:` is NOT mode-invariant — bare runs reliably write it, but team-mode
	// finalize non-deterministically OMITS it (the FO reaches the bounded-teardown
	// terminus without ever writing it; observed 1/3 team runs across two count=3
	// sonnet sets), so gating on it re-introduces a team-vs-bare verdict split — the
	// very thing this entity dissolves. The verdict gate is dropped here (captain's
	// Option A, 2026-06-15) → the team-mode verdict-omission is tracked as the
	// follow-up task `team-mode-verdict-omission` (reeppr990pyzzaejmbnyrvt7), not
	// silently relaxed. expectCondition drains the stream each poll (liveness; the
	// budget resets on activity) while checking the filesystem, bounded by the same
	// roomy silence budget — a live terminalize turn can be quiet between stream
	// lines. After the barrier lands, the clean-exit drain below keeps the fixture
	// alive until the subprocess finishes its terminal ceremony.
	terminalized := func() bool {
		body, _, found := locateEntity(root, "make-it-work")
		return found &&
			frontmatterField.MatchString(body) &&
			someCommitNamesOnly(t, root, "make-it-work")
	}
	if err := watcher.expectCondition(terminalized, quietBudgetDispatchClose, "entity terminalized"); err != nil {
		t.Fatalf("live cycle failed waiting for the entity to terminalize+commit (status: done + path-scoped commit): %v", err)
	}

	// The durable entity barrier can land before the headless FO finishes its
	// terminal ceremony. Drain the existing launcher to a clean exit before this
	// test returns: t.TempDir cleanup must not race a still-running Claude
	// descendant that can continue writing under root after the launcher is killed.
	if _, err := watcher.drainToExit(quietBudgetDispatchClose, "live cycle completion"); err != nil {
		t.Fatalf("live cycle reached its durable terminal state but the FO did not exit cleanly: %v", err)
	}
	if code, exited := watcher.proc.poll(); !exited || code != 0 {
		t.Fatalf("live cycle FO exit = (code=%d, exited=%t), want (0, true)", code, exited)
	}

	// Locate the entity at the REAL completed-cycle end-state. A full FO-to-done
	// cycle ARCHIVES the terminal entity: the flat `make-it-work.md` moves to
	// `_archive/make-it-work.md`. locateEntity searches the original path AND both
	// archive spellings; a missing entity everywhere is a hard FAIL (the cycle
	// neither completed nor left the entity in place).
	entity, where, found := locateEntity(root, "make-it-work")
	if !found {
		t.Fatalf("entity make-it-work not found in place or under _archive/ after the cycle")
	}
	t.Logf("located entity at %s", where)

	// The full-lifecycle END-STATE checks below (stage-report shape, terminal
	// frontmatter, path-scoped commit) are HARD assertions. They verify the REAL
	// completed-and-archived end-state the multi-agent FO+ensign cycle produces.
	// They sit AFTER the team-agnostic terminalized barrier, so a run only reaches
	// them once the entity reached `status: done` on disk — team or bare. Each
	// gates on a PRESENT-and-CORRECT end state: a present-but-wrong end state (e.g.
	// a malformed stage report, a wrong commit scope) is a real Spacedock
	// regression and fails immediately — it is NEVER retried or masked.

	// (a) the appended stage-report section has the protocol shape: heading, a
	// DONE accounting marker, a Summary, and NO checkbox-bullet form.
	if !liveStageReportHeading.MatchString(entity) {
		t.Errorf("entity missing anchored stage-report heading\n%s", entity)
	}
	if !doneMarker.MatchString(entity) {
		t.Errorf("entity missing anchored - DONE: marker\n%s", entity)
	}
	if !strings.Contains(entity, "### Summary") {
		t.Errorf("entity missing ### Summary\n%s", entity)
	}
	if checkboxBullet.MatchString(entity) {
		t.Errorf("entity contains forbidden checkbox-bullet stage-report markers\n%s", entity)
	}

	// (b) the FO finalized the cycle: the entity carries the terminal frontmatter
	// `status: done`. This is MODE-INVARIANT — both team and bare drives land it.
	// The `verdict:` field is NOT asserted here: team-mode finalize omits it
	// non-deterministically (captain's Option A, 2026-06-15 — verdict-presence
	// coverage moved to the follow-up task `team-mode-verdict-omission`,
	// reeppr990pyzzaejmbnyrvt7). Gating on it re-splits the smoke by mode, which is
	// exactly the team-vs-bare coin this entity dissolves.
	if !frontmatterField.MatchString(entity) {
		t.Errorf("entity missing terminal `status: done`\n%s", entity)
	}

	// (c) SOME commit in the history is path-scoped to the entity (names only the
	// entity), the concurrency-safe state-commit invariant at the cycle level.
	// HEAD itself is the FO's archive/finalize commit on a full cycle, so this
	// scans the whole log rather than pinning HEAD (the strict single-file HEAD
	// invariant is pinned deterministically by the skeleton's
	// TestEnsignCycleMechanicalOutputs).
	if !someCommitNamesOnly(t, root, "make-it-work") {
		t.Errorf("no path-scoped commit named only the entity in the cycle history")
	}
}

// startRealisticLifecycleDrive stages the realistic ≥3-stage lifecycle fixture
// (backlog → implementation → done, a flat entity at backlog) in a fresh git root,
// launches the real `spacedock claude` front door headless with the given
// drivePrompt, and returns a streamWatcher over its stream-json plus the fixture
// root. The launch shape is identical across the team-agnostic default cycle and
// the team-FORCED teardown cycle — only the drivePrompt differs (the conn-cue vs.
// the team-mode cue) — so the env/fixture/launch/watcher wiring lives here once.
// The subprocess kill is registered via t.Cleanup so the caller never orphans a
// token-spending claude on any exit path.
func startRealisticLifecycleDrive(t *testing.T, drivePrompt string) (*streamWatcher, string) {
	t.Helper()
	binary := spacedockBinary(t)
	pluginDir := livePluginDir(t)
	model := envOr("SPACEDOCK_LIVE_MODEL", "sonnet")

	// Resolve the isolated child env (clean HOME + the authoritative credential)
	// or skip when no auth mechanism is available. The empty home argument means
	// "read the live $HOME"; the offline unit test drives isolatedClaudeEnv
	// directly with a fake home so it never touches the real ~/.claude.
	childEnv := isolatedClaudeEnv(t, os.Getenv("HOME"))

	// Put the built binary's directory first on the FO subprocess's PATH. The FO
	// contract's first step is `spacedock --version` and the FO knows `spacedock`
	// only by PATH name (not SPACEDOCK_BIN, which is the test's own resolution
	// hook); without this the runner PATH has no `spacedock` and the FO aborts at
	// the binary gate before ever reaching TeamCreate (CI run 26839572693).
	childEnv = withBinaryOnPath(childEnv, binary)

	// Stage a REALISTIC ≥3-stage lifecycle fixture: a git-init'd root with a
	// non-worktree workflow README declaring backlog (initial) → implementation
	// (work) → done (terminal), and a flat entity in the initial (backlog) stage.
	// The distinct WORK stage between initial and terminal is what every real
	// Spacedock workflow has: the FO dispatches backlog→implementation, the ensign
	// works, then the FO TERMINALIZES implementation→done with a verdict — a
	// distinct step from the dispatch (unlike a 2-stage backlog→done fixture, where
	// the dispatch lands status on the terminal stage and the FO never runs a
	// separate finalize, so no verdict is recorded). The real FO drives it to the
	// TERMINAL stage, so the ensign that finishes writes `## Stage Report: done`,
	// which the stage-agnostic liveStageReportHeading regex matches.
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeRealisticLifecycle())
	entityPath := filepath.Join(root, "make-it-work.md")
	writeFile(t, entityPath, entityFixture())
	gitInit(t, root)

	// The canonical neutral runtime-scenario prompt (README §Shared runtime
	// scenarios; the same string every shared fixture uses). The fixture supplies
	// the context — a single backlog entity in a backlog→implementation→done
	// workflow — and the FO acts from its contract; no coaching ("drive to done",
	// "dispatch for each stage", "do not stop") is injected. Coaching is exactly the
	// arbitrary instruction the scenario-prompt rule forbids, and it misframes the
	// FO's job (dispatch ensigns vs. drive the entity to a verdicted, archived
	// terminal state).
	task := "Use $spacedock:first-officer for this whole run."

	// The real front door: `spacedock claude --plugin-dir <repo> --skip-compat-check
	// -- -p <bootstrap> ... <task>`. --plugin-dir and --skip-compat-check are
	// spacedock-owned flags BEFORE `--`: --plugin-dir loads the local v1 plugin
	// checkout (and relaxes the contract gate, the keystone this entity fixes);
	// --skip-compat-check is belt-and-braces. Every host flag (-p, --permission-mode,
	// --output-format, --verbose, --model) rides AFTER `--` and forwards verbatim to
	// claude, ahead of the fenced task — only --plugin-dir is promoted before `--`
	// (the spike WINNER: every other host flag stays after `--`). stream-json +
	// bypassPermissions + the model pin mirror the headless launch the Python net
	// uses. CLAUDECODE is dropped by isolatedClaudeEnv so the binary takes the real
	// front-door path rather than a nested-session shortcut.
	//
	// The subprocess runs under a plain context.Background() with NO deadline ctx
	// (AC-1 bans the monolithic timeout); the streamWatcher's per-step quiet
	// budgets ARE the timeout discipline. Both stdout and stderr feed one io.Pipe
	// so the watcher reads claude's stream-json line-by-line as it arrives (a hang
	// leaves a partial transcript — AC-2) rather than CombinedOutput()'s zero-byte
	// block-until-exit.
	//
	// drivePrompt is the host `-p` input. The anti-early-shutdown clause (carried by
	// every caller via antiShutdownOverride) counters upstream claude-code bug
	// #55297 (a regression in 2.1.126; CI runs 2.1.161): in `claude -p` with an
	// active team the harness injects "you cannot return a response until your team
	// is shut down … shut down before your final response" EVERY turn, and the model
	// panic-shuts-down the team before finishing the work — the premature teardown
	// that left the entity un-terminalized. No FO-contract prose can out-argue a
	// per-turn harness reminder, so the override lives in the `-p` input instead. It
	// is GENERIC — it governs shutdown TIMING only, naming no stage or task — so it
	// does not coach workflow mechanics.
	cmd := exec.Command(binary, "claude",
		"--plugin-dir", pluginDir,
		"--skip-compat-check",
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

	// Run cmd.Wait() in the background; closing the pipe write-end when it
	// returns lets the watcher's line scanner reach EOF and drain the final
	// lines. cmdPoller.poll() reads the recorded exit non-blockingly so the
	// watcher's expect/expectDispatchClose loops can check for early exit, and
	// expectExit waits on it. The poller's kill() forcibly stops a hung claude.
	poller := newCmdPoller(cmd, pw)
	watcher := newStreamWatcher(newPipeLineSource(pr), poller, func(line string) { t.Log(line) })

	// Kill the subprocess on ANY exit path. Only expectExit kills on its own
	// timeout; an EARLY t.Fatalf (a dispatch-close-stall below) would otherwise
	// orphan a token-spending `claude`. kill() is a no-op once the process has
	// exited, so this is harmless on the clean path.
	t.Cleanup(poller.kill)

	return watcher, root
}

// spacedockBinary resolves the built v1 binary the test shells. SPACEDOCK_BIN
// (set by the CI job after `go build -o ./spacedock`) takes precedence; locally
// it falls back to a `spacedock` on PATH. The test fails loudly when neither
// resolves rather than silently shelling a stale or absent binary.
func spacedockBinary(t *testing.T) string {
	t.Helper()
	if p := os.Getenv("SPACEDOCK_BIN"); p != "" {
		abs, err := filepath.Abs(p)
		if err != nil {
			t.Fatalf("SPACEDOCK_BIN=%q is not resolvable: %v", p, err)
		}
		if _, err := os.Stat(abs); err != nil {
			t.Fatalf("SPACEDOCK_BIN=%q does not exist: %v", abs, err)
		}
		return abs
	}
	p, err := exec.LookPath("spacedock")
	if err != nil {
		t.Fatal("no spacedock binary: set SPACEDOCK_BIN to the built binary or put spacedock on PATH")
	}
	return p
}

// repoRoot resolves the plugin-checkout root passed to --plugin-dir. The
// ensigncycle package lives at internal/ensigncycle, so the repo root is two
// directories up from the test's working directory.
func repoRoot(t *testing.T) string {
	t.Helper()
	if p := os.Getenv("SPACEDOCK_REPO_ROOT"); p != "" {
		abs, err := filepath.Abs(p)
		if err != nil {
			t.Fatalf("SPACEDOCK_REPO_ROOT=%q is not resolvable: %v", p, err)
		}
		return abs
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root, err := filepath.Abs(filepath.Join(wd, "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

// livePluginDir stages an ISOLATED plugin checkout for `--plugin-dir` and returns
// its path. It exists to stop a wrong-root boot: the real repo root carries a
// discoverable `docs/dev` workflow (with live entities), so an FO that anchors its
// `git rev-parse --show-toplevel` + `status --discover` on the plugin path — instead
// of its isolated cwd fixture — finds and drives the REAL workflow. Staging copies
// ONLY the plugin scaffolding (`.claude-plugin/`, `skills/`, `agents/`) into a temp
// dir with NO `docs/dev` sibling, then `git init`s it so a `rev-parse` from the
// plugin path resolves to a workflow-free root. An FO that boots from here discovers
// zero workflows and falls back to the cwd fixture. The result is cached per repo
// root so parallel scenarios share one staged copy.
func livePluginDir(t *testing.T) string {
	t.Helper()
	return cachedLivePluginDir(t, repoRoot(t))
}

var (
	livePluginOnce sync.Once
	livePluginPath string
	livePluginErr  error
)

func cachedLivePluginDir(t *testing.T, repo string) string {
	t.Helper()
	livePluginOnce.Do(func() {
		// MkdirTemp (not t.TempDir) so the staged plugin outlives the first test's
		// cleanup and the cached path stays valid for every scenario in the run.
		staged, err := os.MkdirTemp("", "spacedock-live-plugin-")
		if err != nil {
			livePluginErr = err
			return
		}
		for _, sub := range []string{".claude-plugin", "skills", "agents"} {
			src := filepath.Join(repo, sub)
			if _, statErr := os.Stat(src); statErr != nil {
				continue // optional members (e.g. a layout without a top-level agents/)
			}
			if copyErr := copyTree(src, filepath.Join(staged, sub)); copyErr != nil {
				livePluginErr = copyErr
				return
			}
		}
		// git init so the FO's `git rev-parse --show-toplevel` resolves to this
		// workflow-free root, not an enclosing checkout that has a docs/dev.
		testgit.InitRepo(t, staged, "-q")
		livePluginPath = staged
	})
	if livePluginErr != nil {
		t.Fatalf("stage isolated live plugin dir: %v", livePluginErr)
	}
	return livePluginPath
}

// copyTree recursively copies src to dst, preserving file modes. Symlinks are
// resolved to real files so the staged plugin has no path back into the real repo.
func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(src, path)
		if relErr != nil {
			return relErr
		}
		target := filepath.Join(dst, rel)
		info, infoErr := os.Stat(path) // Stat (not Lstat) resolves symlinks to real content
		if infoErr != nil {
			return infoErr
		}
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if mkErr := os.MkdirAll(filepath.Dir(target), 0o755); mkErr != nil {
			return mkErr
		}
		return os.WriteFile(target, data, info.Mode().Perm())
	})
}

// readmeRealisticLifecycle is the live cycle's workflow README: a realistic
// ≥3-stage lifecycle — backlog (initial) → implementation (work) → done
// (terminal) — matching every real Spacedock workflow. The distinct WORK stage
// between initial and terminal makes the FO's TERMINALIZE step
// (implementation→done with a verdict) DISTINCT from its DISPATCH step
// (backlog→implementation): the FO records a verdict naturally at terminalization
// and the M1 verdict gate has a real trigger. A 2-stage backlog→done fixture
// collapses dispatch and terminalize onto the same stage, so the FO never runs a
// distinct finalize and no verdict is recorded (the failure this fixture fixes).
// It is live-only (kept beside the //go:build live test) so the offline minimal
// fixture readmeNonWorktree — which the single-stage mechanical test pins — is
// untouched. All stages are non-worktree to keep the flat-entity path.
func readmeRealisticLifecycle() string {
	return "---\n" +
		"commissioned-by: spacedock@1\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: backlog\n" +
		"      initial: true\n" +
		"    - name: implementation\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Fixture Workflow\n" +
		"\n" +
		"### backlog\n\nseed.\n\n- **Outputs:** a one-line note.\n\n" +
		"### implementation\n\nDo the trivial work and write the note.\n\n- **Outputs:** the note recorded.\n\n" +
		"### done\n\nterm.\n"
}

// isTeamCreate matches the FO's TeamCreate assistant tool_use — the first
// progress beat of the teams-mode live cycle (the team is engaged).
func isTeamCreate(e streamEntry) bool {
	b := e.toolUseBlock()
	return b != nil && b.Name == "TeamCreate"
}

// isEnsignDispatch matches the FO's first ensign dispatch — an
// Agent(subagent_type="spacedock:ensign") assistant tool_use. The contract runs
// `spacedock dispatch spawn-standing-all` immediately before this dispatch, so its
// OPEN is the reliable barrier for "standing teammates have been injected" (used by
// the residency test, which only needs injection to have run, not the dispatch to
// close). It scans ALL tool_use blocks so an Agent dispatch riding as a second
// block in a multi-tool turn is not missed.
func isEnsignDispatch(e streamEntry) bool {
	for _, b := range e.toolUseBlocks() {
		if b.Name == "Agent" && b.Input.SubagentType == "spacedock:ensign" {
			return true
		}
	}
	return false
}

// cmdPoller is the live procPoller: it Waits the exec.Cmd in the background,
// records the exit code, and reports it non-blockingly via poll() so the
// watcher's expect loops can detect an early crash and expectExit can wait for a
// clean exit. Closing the pipe write-end on exit unblocks the line scanner so
// the watcher drains the final stream-json lines. kill() forcibly stops a hung
// claude on a quiet-budget timeout.
type cmdPoller struct {
	cmd *exec.Cmd

	mu     sync.Mutex
	exited bool
	code   int
}

// newCmdPoller starts the background Wait. pw is the pipe write-end shared by
// cmd.Stdout/Stderr; it is closed when Wait returns so the reader reaches EOF.
func newCmdPoller(cmd *exec.Cmd, pw io.Closer) *cmdPoller {
	p := &cmdPoller{cmd: cmd}
	go func() {
		err := cmd.Wait()
		code := 0
		if exit, ok := err.(*exec.ExitError); ok {
			code = exit.ExitCode()
		} else if err != nil {
			code = -1
		}
		pw.Close()
		p.mu.Lock()
		p.exited = true
		p.code = code
		p.mu.Unlock()
	}()
	return p
}

func (p *cmdPoller) poll() (int, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.code, p.exited
}

func (p *cmdPoller) kill() {
	if p.cmd.Process != nil {
		_ = p.cmd.Process.Kill()
	}
}
