---
title: Build a terminal/pty live harness for team-mode e2e (residency + teardown)
status: validation
source: "FO + captain (2026-06-16, during 2yf): headless `claude -p` cannot sustain team mode — anthropics/claude-code 2.1.178 dropped the native TeamCreate/TeamDelete tools from headless sessions (anthropics/claude-code#68721), and even with tools present the SDK/headless session lifecycle races to end_turn before teammates finish (anthropics/claude-code-action#1124). Per 7e's recorded steer (headless `-p` goes bare), the two forced-team `-p` live tests (TestLiveEnsignCycleTeamTeardown, TestLiveStandingResidencyInjectsCommOfficer) were RETIRED in 2yf because they cannot work headless. Team-mode MECHANISMS stay covered offline (internal/dispatch/spawn_standing_all_test.go for comm-officer injection; internal/ensigncycle/teardown_grade_watcher_test.go + testdata/sonnet_teamdelete_*.jsonl for the bounded-teardown marker grading). The GAP this task closes: live end-to-end team-mode coverage (FO creates a real team, injects the comm-officer standing teammate into the roster, dispatches through the team, terminalizes, and runs the bounded teardown emitting TERMINAL_TEARDOWN_BOUNDED). That requires an INTERACTIVE (pseudo-terminal) harness where team tools are present and the session stays alive — the current live suite is entirely `claude -p`."
started: 2026-06-16T15:15:39Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-live-team-mode-terminal-harness
issue:
id: m40mphxan8phr3t3tp03gk89
sprint: 0204-structured-reads
mod-block: merge:pr-merge
pr: "#390"
sprint-readiness: in-progress
---

## Problem

Team mode is interactive-only: it needs a live tty session where the native team tools (TeamCreate/TeamDelete) are exposed and the agent loop stays alive while teammates work. The entire live-e2e suite drives `claude -p` (headless), which cannot do either reliably. So the live end-to-end team-mode path — comm-officer roster injection and the bounded terminal teardown — has no working harness; it has only offline mechanism coverage.

## What's needed

A live harness that drives a real INTERACTIVE Claude session via a pseudo-terminal (pty), boots the FO, and exercises:
- standing-teammate (comm-officer) injection landing in the team `config.json` roster;
- the bounded terminal teardown emitting the `TERMINAL_TEARDOWN_BOUNDED` marker after a real TeamCreate→dispatch→terminalize→teardown cycle.

## Notes / prior art

- Retirement of the two forced-team `-p` tests + the go-bare alignment happened in 2yf (shared-merge-dispatch-contract). 7e (headless-dispatch-mode-intent, #381) determined headless `-p` goes bare. #271 (headless-fo-drive-flake) first surfaced the headless silent-await stall and deferred the runtime-await question.
- Mechanism coverage that already exists (do not duplicate): `internal/dispatch/spawn_standing_all_test.go`, `standing_parity_test.go`; `internal/ensigncycle/teardown_grade_watcher_test.go`, `teardown_grade_test.go`, `streamwatch_test.go` + `testdata/sonnet_teamdelete_*.jsonl`.
- Upstream constraints to watch: anthropics/claude-code#68721 (2.1.178 headless team-tool regression — may resolve and restore headless tools), anthropics/claude-code-action#1124 (SDK/headless team lifecycle).

## Spike: riskiest mechanism (DEMONSTRATED, not asserted)

The riskiest claim — the one that invalidates the whole task if false — is: *can a pty/tmux drive of a REAL interactive `spacedock claude` session expose the native TeamCreate/TeamDelete tools AND stay resident while teammates work, where headless `claude -p` cannot?* I drove the smallest end-to-end exercise of that path with tmux on this machine (tmux 3.6a, claude 2.1.177, Go 1.26.1). It PASSED.

**Probe transport** (matches spacedock-gym's `internal/driver` approach): `tmux new-session -d -s <name> -c <cwd> "<launch>"`, then `tmux send-keys -t <pane> -l <text>` + `Enter` to drive, `tmux capture-pane -p -t <pane>` to observe. Launch line is the live runner's exact front door minus `-p`: `spacedock claude --plugin-dir <staged> --skip-contract-check -- --permission-mode bypassPermissions --model sonnet`.

**Probe 1 — team tools present + resident.** Booted the interactive session; the TUI banner showed `Claude Team · Recce` and `@spacedock:first-officer` loaded (team mode, FO contract live). Sent a prompt asking it to call TeamCreate. A real team `config.json` landed at `~/.claude/teams/ptyspike-<ts>/config.json` with `team-lead` in `members[]`, and `tmux has-session` reported the session still ALIVE after the tool call. This is exactly what headless `-p` cannot do (no team tools post-2.1.178; SDK races to end_turn).

**Probe 2 — full teardown branch.** Drove TeamCreate → `Agent(run_in_background)` spawn of a teammate `slowpoke` → `TeamDelete` while slowpoke active. The FO's own result table (captured from the FO pane):
- `TeamCreate("ptyteardown2-probe")` → `✓ team created`
- `Agent(name=slowpoke, run_in_background=true)` → `✓ member active` (on-disk roster became `members: ['team-lead', 'slowpoke']`; slowpoke launched into a SECOND tmux pane via the harness's own `claude --agent-id slowpoke@… --team-name … --agent-type general-purpose`)
- `TeamDelete() while slowpoke alive` → `✗ Cannot cleanup team with 1 active member(s): slowpoke` (the exact bounded condition the marker grades)
- marker → `✓ TERMINAL_TEARDOWN_BOUNDED: …` emitted by the FO as its own assistant text

Session stayed resident throughout (`Idle · teammates running · 1 teammate`). The marker was FO-authored, not a prompt echo (the marker/refusal strings were deliberately absent from the injected prompt; an earlier probe revision that put the marker string IN the prompt produced a false positive — recorded below as a test-design hazard).

**Spike findings that shape the design (these would have cost an implementation cycle to discover):**
1. **No pipe on launch.** Piping the launch through `| tee` makes claude's stdout a non-tty, so the interactive TUI never renders (the session ran but stayed at the wrapper banner). Claude MUST own the pane tty; observe via `capture-pane`, never a pipe. (spacedock-gym already does this.)
2. **Teammate spawn opens a new pane → capture the FO pane explicitly.** A `run_in_background` Agent spawn creates a second tmux pane and the active pane switches to the teammate. `capture-pane -t <session>` then reads the teammate, not the FO — the FO's TeamDelete/marker live in pane 0 (`#{pane_title}` = `spacedock:first-officer`). The driver MUST resolve and read the FO pane by title, exactly spacedock-gym's `firstOfficerPaneIndex` (`internal/driver/foreground.go`). This is the single most important borrowed mechanism.
3. **Send only when genuinely idle.** Keystrokes sent while the FO is mid-boot QUEUE behind the boot turn (TUI shows "Press up to edit queued messages") and the probe stalls. Gate the send on a stable-idle signal (status line present AND no spinner word / "esc to interrupt") for N consecutive polls.
4. **Auth, and the benchmark-token rotation hazard.** The repo's `~/.claude/benchmark-token` was STALE (401) on this machine; the live headless suite shares this exact dependency. The spike ran on the operator's logged-in auth (`~/.claude.json`). Note: operator auth lives in the config-dir-level `.claude.json`, so isolating `CLAUDE_CONFIG_DIR` away from it yields "Not logged in". The pty harness inherits the headless suite's auth decision tree (`isolatedClaudeEnv`: OAuth benchmark-token / `ANTHROPIC_API_KEY`, else skip) — and the same token-expiry operational risk.

## Approach

Build the pty/tmux harness as a NEW DRIVER (transport) over the EXISTING shared scenario+fixture set, not a new scenario set. The current live suite already separates a host-neutral scenario table (`sharedRuntimeScenarios()`) from a host adapter that launches the model and feeds shared assertions. Today there is exactly ONE transport baked into that adapter — `exec.Command(spacedock, "claude", … "-p", …)` inside `claudeLiveRunner.run` (`internal/ensigncycle/claude_live_runner_test.go:320-413`). The pty harness is a SECOND transport that runs the SAME scenarios.

The minimal refactor: extract the transport behind a small interface so the per-scenario orchestration (set up the workflow fixture, run the scenario, feed the shared assertions) is driver-agnostic, then add the pty driver as a second implementation alongside the current `-p` driver.

### The driver seam (the minimal refactor)

Today's coupling: each `runClaude<Scenario>Scenario(t, runner, scenario)` function calls `runner.run(t, scenario, workflowRoot, prompt) → claudeScenarioResult{finalMessage, stream, …}`. The orchestration (workflow setup + assertions) is already separate from the launch; only the launch mechanism is hard-bound to `-p`. Extract:

```go
// liveDriver is the transport seam: it turns a prompt + workflow root into the
// observed (finalMessage, stream) the shared assertions consume. The headless
// `-p` runner and the pty/tmux runner are two implementations; the scenario
// orchestration and assertions do not know which transport ran.
type liveDriver interface {
    run(t *testing.T, scenario sharedRuntimeScenario, workflowRoot, prompt string) liveResult
}

type liveResult struct {
    finalMessage string // headless: stream result event; pty: FO-pane final text / sentinel
    stream       string // headless: stream-json transcript; pty: session jsonl under CLAUDE_CONFIG_DIR
    artifactDir  string
    duration     time.Duration
}
```

`claudeLiveRunner` already satisfies this shape (`run(t, scenario, workflowRoot, prompt) claudeScenarioResult`) — the change is renaming `claudeScenarioResult`→`liveResult` (or aliasing) and declaring it implements `liveDriver`. The per-scenario functions take `liveDriver` instead of `claudeLiveRunner`. NO scenario, prompt, fixture, or assertion is rewritten — they are reused verbatim.

The pty driver `ptyLiveDriver` implements the same `run`:
- launches `spacedock claude --plugin-dir … --skip-contract-check -- --permission-mode bypassPermissions --model …` (NO `-p`) in a detached tmux pane (claude owns the pane tty);
- waits for stable-idle, then `send-keys -l <prompt>` + `Enter`;
- resolves the FO pane by `#{pane_title}` and reads it via `capture-pane` (borrowed `firstOfficerPaneIndex`);
- sources `stream` from the session jsonl written under the pinned `CLAUDE_CONFIG_DIR` (the same transcript shape the existing assertions already grade), so the shared marker/roster assertions run UNCHANGED;
- on exit, `tmux kill-session` (deferred teardown), reaps panes.

`liveResult.stream` being the session jsonl (not the `-p` stream-json) is the key compatibility point: the existing `streamWatcher`/grading code parses the stream-json dialect, and the interactive session writes the SAME dialect to its session jsonl. The pty driver tails that file for liveness instead of a stdout pipe.

### What the pty driver runs (shared) vs adds (team-mode only)

**Reused verbatim (NO duplication):**
- Scenario table: `sharedRuntimeScenarios()` (`internal/ensigncycle/shared_scenarios_test.go`) — `gate-guardrail`, `rejection-flow`, `feedback-3-cycle-escalation`, `merge-hook-guardrail`, `filing`, `shallow-boot`.
- Per-scenario fixtures + assertions: `writeGateWorkflow`/`assertGateHeld`, `writeRejectionWorkflow`/`assertRejectionFlow`, `assertThirdCycleEscalation`, `writeMergeHookGuardWorkflow`/`assertMergeHookGuardHeld`, `assertClaudeFilingViaNew`, `assertShallowBoot` + `assertNoTeamCreateBeforeGreet` + `assertShallowBootMeasured`.
- Liveness: `streamWatcher` (the shared FOStreamWatcher port) — one mechanism.
- The auth/HOME isolation decision tree: `isolatedClaudeEnv` / `decideClaudeEnv`.

**Added by the pty driver, ON TOP of the shared set (team-mode-only, the gap this task closes):**
- **comm-officer roster injection (live):** after a team-mode boot+dispatch, assert the comm-officer standing teammate actually lands in the team `config.json` `members[]` under the run's isolated team root. (Offline `spawn_standing_all_test.go` proves the SPEC is emitted; the live assertion proves the member JOINS the roster — the spike already showed an Agent-spawned teammate joining `members[]`.)
- **bounded-teardown marker (live):** after a real TeamCreate→dispatch→terminalize→teardown, assert the FO emits `TERMINAL_TEARDOWN_BOUNDED` (grade with the EXISTING `gradeTerminalTeardown` / `markerEmittedByAssistant` over the captured stream — same grader the offline fixture suite uses; the spike showed the live FO emitting it on the active-member refusal).

These two assertions are the resurrected content of the retired `TestLiveStandingResidencyInjectsCommOfficer` and `TestLiveEnsignCycleTeamTeardown`, now driven by the pty driver instead of `-p`.

### Coverage boundary (what m4 closes; what it must NOT duplicate)

- **m4 closes:** the LIVE team-mode end-to-end gap — real team → comm-officer in roster → dispatch → terminalize → bounded teardown — through the pty driver, resurrecting the two retired forced-team tests. This is the live home `live_test.go:76-82` explicitly defers to ("live team end-to-end is deferred to the terminal/pty harness task m40mphxan8phr3t3tp03gk89").
- **m4 does NOT duplicate** the OFFLINE mechanism coverage, which stays as-is: comm-officer SPEC emission/dedup (`internal/dispatch/spawn_standing_all_test.go`, `standing_parity_test.go`) and bounded-teardown marker GRADING over fixtures (`internal/ensigncycle/teardown_grade_watcher_test.go`, `teardown_grade_test.go`, `streamwatch_test.go` + `testdata/sonnet_teamdelete_*.stream.jsonl`). The pty test REUSES the grader (`gradeTerminalTeardown`) against a LIVE stream; it does not reimplement grading.
- **Relation to siblings:** unblocks `team-mode-verdict-omission` (reeppr990pyzzaejmbnyrvt7) — that task needs a stable team-mode live drive to test verdict presence, which only this pty harness provides — and complements `bare-mode-coverage-baseline` (e3z): the default `TestLiveEnsignCycle` stays team-AGNOSTIC/bare, this harness owns team-mode. The current `TestLiveEnsignCycle` flake gating v0.20.4 is a separate budget/await issue, but team-mode coverage having no live home is why the marker gate was dropped from it; m4 gives that gate a real home.

### spacedock-gym decision: reference-only (port the mechanism, not the module)

spacedock-gym (`github.com/spacedock-research/spacedock-gym`, go 1.26) has a working tmux `Driver` (`internal/driver/driver.go`: `Launch/Send/WaitForTurn/WaitForGate/Snapshot/Kill/Alive`, pane resolution in `foreground.go`, picker drive in `picker.go`, pane reads via `capture-pane`). It is a SEPARATE Go module with a different module path and Go version (1.26 vs spacedock-v1's 1.22), and it is a transcript/session test harness, not a spacedock runtime component. Decision: **reference-only** — build a minimal spacedock-v1 `ptyLiveDriver` (under `internal/ensigncycle`, `//go:build live`) that borrows gym's proven mechanisms (no-pipe launch, `firstOfficerPaneIndex` pane-by-title resolution, capture-pane polling, send-keys keystroke drive) without importing the module. Rationale: (a) cross-module import would force a Go-version bump and a dependency on a sibling research repo for a test-only path; (b) the spacedock-v1 driver only needs the narrow `liveResult` surface, far less than gym's full picker/cast machinery; (c) the shared FIXTURES already exist in BOTH repos (`testdata/sonnet_teamdelete_*`), so there is no fixture to port — only the transport technique. Extracting a shared lib is rejected by YAGNI: one consumer (the pty test) does not justify a third module. If a future second consumer appears, revisit.

## Acceptance criteria

AC-1. A `liveDriver` interface exists in `internal/ensigncycle` and BOTH the headless `-p` runner (`claudeLiveRunner`) and the new `ptyLiveDriver` implement it; the per-scenario orchestration functions accept `liveDriver`, not a concrete runner. *Test:* a compile-time assertion `var _ liveDriver = …` for each impl, and `TestLiveClaudeSharedScenarios` still green driving through the interface (the refactor is behavior-preserving for the existing `-p` path).
AC-2. The pty driver runs the EXISTING `sharedRuntimeScenarios()` with the EXISTING fixtures and assertions — no scenario, fixture, or assertion is duplicated or forked. *Test:* a parity meta-test (mirroring `claudeScenarioRunners` coverage guard) fails if the pty driver lacks a runner for any shared scenario; `grep` guard / review confirms the pty path imports the shared assertions rather than redefining them.
AC-3. A live pty team-mode test drives a real TeamCreate→dispatch and the comm-officer standing teammate lands in the team `config.json` `members[]` under the run's isolated team root. *Test:* `//go:build live` test asserting the on-disk roster contains `comm-officer` after the drive (resurrects `TestLiveStandingResidencyInjectsCommOfficer` via pty). Spike-proven mechanism: an Agent-spawned teammate joined `members[]`.
AC-4. A live pty team-mode test drives a real TeamCreate→dispatch→terminalize→teardown and the FO emits `TERMINAL_TEARDOWN_BOUNDED`, graded by the EXISTING `gradeTerminalTeardown`/`markerEmittedByAssistant` over the captured live stream. *Test:* `//go:build live` test (resurrects `TestLiveEnsignCycleTeamTeardown` via pty) reusing the offline grader against the live session jsonl. Spike-proven mechanism: the live FO emitted the marker on an active-member TeamDelete refusal.
AC-5. The pty driver reads the FO pane by `#{pane_title}` (not the active pane), launches without a stdout pipe, and gates its first send on stable-idle — the three spike-surfaced hazards are encoded, not left to rediscovery. *Test:* the driver resolves the pane by title (unit-level seam test over a synthetic `list-panes` output, mirroring gym's `firstOfficerPaneIndex` test); the live tests passing is the end-to-end proof.
AC-6. The pty live tests are gated `//go:build live` and skip (never fatal) when no auth is available, reusing `isolatedClaudeEnv`. *Test:* default `go test ./...` excludes them (build tag); `go test -tags live` with no credential skips with the existing auth-missing message.

## Test plan

- **Cost/complexity:** the refactor (AC-1/AC-2) is cheap and offline-verifiable (compile + the existing `-p` suite stays green). The live pty tests (AC-3/AC-4) are multi-minute, live-credentialed, `-tags live` only — same cost class as the current live suite, run behind the CI-E2E gate. AC-5's pane-resolution seam is a fast offline unit test.
- **Fixture tests:** the pane-by-title resolver gets a table test over synthetic `tmux list-panes` output (no tmux needed), borrowed from gym's `foreground` test shape.
- **CLI/behavior:** none new — the harness is test-only; it drives the EXISTING `spacedock claude` front door.
- **Live workflow tests:** AC-3/AC-4 are the live claims; they are the resurrected forced-team tests, now over the pty transport. They reuse the shared fixtures and the offline grader, so the live test adds only the team-mode assertions (roster membership, marker emission), not a parallel grading stack.
- **Risk already retired by the spike:** the load-bearing unknown (interactive pty exposes team tools + stays resident through TeamCreate/Agent/TeamDelete) is DEMONSTRATED above; implementation starts from a proven mechanism, and the spike's three hazards (no-pipe, FO-pane capture, idle-gated send) seed the driver's first tests.

## Docs impact

None user-visible: the harness is a `//go:build live` test-only addition. No CLI surface, startup banner, or docs-site behavior changes, so no doc diff is required at the ideation gate.

## Stage Report: ideation

- DONE: Spike the RISKIEST mechanism FIRST and record the evidence in the entity body: prove a pty/tmux harness can drive a REAL interactive Claude session (not `claude -p`) in which the native TeamCreate/TeamDelete tools are PRESENT and the session stays alive while teammates work — captured end-to-end by (a) the comm-officer landing in the team config.json roster and (b) a real TeamCreate→dispatch→terminalize→teardown cycle emitting TERMINAL_TEARDOWN_BOUNDED. Reuse spacedock-gym's internal/driver tmux pane driver as prior art to de-risk this cheaply; the unknown that invalidates everything is whether the pty drive actually exposes team tools + stays resident where headless `-p` cannot.
  Two live tmux probes PASSED on this machine (tmux 3.6a, claude 2.1.177): probe 1 — interactive `spacedock claude` boots team mode (`Claude Team`) and TeamCreate lands a real `~/.claude/teams/.../config.json` with the session resident; probe 2 — TeamCreate→Agent-spawn (teammate joined `members[]` in a 2nd pane)→TeamDelete refused with active member→FO emits `TERMINAL_TEARDOWN_BOUNDED`, session resident throughout. See `## Spike` section.
- DONE: Design the scenario/driver SPLIT so the new pty/tmux driver runs the EXISTING shared scenarios and reuses the EXISTING fixtures — NOT a duplicated scenario set, test codepath, or fixtures. [pinned `liveDriver` interface, named reused scenarios+fixtures+grader, named the two team-mode assertions added on top]
  See `## The driver seam` + `## What the pty driver runs (shared) vs adds`: `liveDriver` interface (both `-p` and pty implement), reuses `sharedRuntimeScenarios()` + per-scenario fixtures/assertions + `streamWatcher`, adds only live comm-officer roster + live `TERMINAL_TEARDOWN_BOUNDED` graded by the EXISTING `gradeTerminalTeardown`.
- DONE: Pin the coverage boundary: m4 closes the LIVE team-mode end-to-end gap, resurrecting the two retired forced-team `-p` tests through the pty driver — WITHOUT duplicating the existing OFFLINE mechanism coverage. State the decision on spacedock-gym's driver.
  See `## Coverage boundary` (resurrects `TestLiveStandingResidencyInjectsCommOfficer` + `TestLiveEnsignCycleTeamTeardown` via pty; offline spec-emit + fixture-grading coverage untouched; relates to verdict-omission + bare-baseline siblings) and `## spacedock-gym decision`: reference-only (borrow mechanism, no cross-module import / no shared lib, YAGNI).

### Summary

Designed the pty/tmux harness as a SECOND transport (`ptyLiveDriver`) over the one existing shared scenario+fixture set, behind a minimal extracted `liveDriver` interface that the current `claudeLiveRunner` already shape-satisfies — so no scenario, fixture, assertion, or grader is duplicated; the pty driver only ADDS the two team-mode-only live assertions (comm-officer roster membership, bounded-teardown marker). The riskiest unknown — whether an interactive pty drive exposes native team tools and stays resident where headless `-p` cannot — was DEMONSTRATED with two live tmux probes, which also surfaced three implementation hazards now encoded as ACs (no stdout pipe on launch, capture the FO pane by title not the active pane, gate the first send on stable-idle) and one operational hazard (the benchmark-token was stale/401; the harness inherits the headless suite's auth dependency). spacedock-gym is reference-only: borrow its proven tmux mechanisms (pane-by-title resolution, capture-pane polling, send-keys drive) without importing the separate module.

## Stage Report: implementation

- DONE: Extract the `liveDriver` seam reusing every existing scenario/fixture/assertion verbatim — no fork, no duplication — so the existing `-p` suite (`TestLiveClaudeSharedScenarios`) stays green driving through the interface (AC-1/AC-2).
  Renamed `claudeScenarioResult`→`liveResult`, added `liveDriver` (run/model/home/withStubPATH) with `var _ liveDriver = claudeLiveRunner{}`; the six `runClaude*Scenario` runners + `claudeScenarioRunners()` map now take `liveDriver`. Offline-verified: `go test ./...` green (0 fails), `TestSharedScenarioRunnerCoverage` green under `-tags live` (commit 6ee167cd).
- DONE: `ptyLiveDriver` encodes all three spike hazards as code — no stdout pipe on launch, resolve the FO pane by `#{pane_title}` not the active pane, gate the first send on stable-idle — with the pane-resolver covered by an offline table test (AC-5).
  `pty_live_driver_test.go` (`//go:build live`) implements `liveDriver` via tmux new-session (no pipe), `waitStableIdle`, `firstOfficerPaneIndex` (ported from gym, NOT imported); `TestFirstOfficerPaneIndex` runs offline (green). Two hazards the spike MISSED, found + fixed by tmux probe: (a) `tmux new-session` against a pre-existing server drops the command env → child env now rides per-session `-e KEY=VAL` (else HOME/CLAUDE_CONFIG_DIR/OAuth-token leak to operator real values); (b) a fresh isolated HOME stalls at Claude Code's theme/login/trust dialogs before the FO loads → `seedInteractiveClaudeConfig` pre-clears them (commit f643d278).
- DONE: The two added team-mode live assertions (comm-officer lands in team `config.json` `members[]`; live `TERMINAL_TEARDOWN_BOUNDED` graded by the EXISTING `gradeTerminalTeardown`/`markerEmittedByAssistant`) are `//go:build live`, reuse the offline grader with no parallel grading stack, and skip-not-fatal without auth (AC-3/AC-4/AC-6).
  `pty_team_mode_live_test.go` resurrects the two retired forced-team tests over the pty transport, reusing `isolatedClaudeEnv` (skip-not-fatal, AC-6), `gradeTerminalTeardown`/`expectTerminalTeardownGrade` (no new grader), and the shared fixtures (`readmeRealisticLifecycle`/`entityFixture`/`locateEntity`/...). Wired into `runtime-live-e2e.yml` (tmux install + dedicated step) behind the CI-E2E gate (commit e5aaef92).

### Summary

Built `ptyLiveDriver` as a second `liveDriver` transport over the one shared scenario+fixture set behind a minimal extracted interface the existing `claudeLiveRunner` already shape-satisfied — no scenario/fixture/assertion/grader forked. AC-1/AC-2/AC-5 are proven OFFLINE (build clean, `go test ./...` 0 fails, parity meta-test + pane-resolver table test green); the two `//go:build live` team-mode tests (AC-3 roster injection, AC-4 bounded-teardown marker) reuse the offline grader and skip-not-fatal without auth (AC-6), wired into the CI-E2E gate with tmux installed. BLOCKER for live AC-3/AC-4 on THIS machine: the repo `~/.claude/benchmark-token` is STALE (401) — a tmux probe confirmed the FO boots to idle, the FO pane resolves by title, the send lands, and the seed clears all onboarding dialogs, but every model call 401s so no real TeamCreate/teardown runs and (separately) a 401 session writes no transcript jsonl. The live tests need a fresh `~/.claude/benchmark-token` (or CI's `ANTHROPIC_API_KEY`) to green end-to-end; the harness plumbing is validated up to the credentialed model call.

### Feedback Cycles

**Detached adversarial audit (validation, cycle 1) — REFUTED NOTHING MATERIAL.** Ran on a THROWAWAY detached worktree (`git worktree add --detach` at `e5aaef92`, torn down after), never the implementation worktree. Four refutation attempts:

1. **AC-2 parity guard — genuine.** Removed the `"filing"` runner from `claudeScenarioRunners()`; `TestSharedScenarioRunnerCoverage` REDs with `shared scenario "filing" has no Claude runner`. The guard is a real discriminator. Structural note: the pty driver does NOT have a separate runner map — it SHARES `claudeScenarioRunners()` with the `-p` runner through the `liveDriver` interface, so AC-2's "pty driver lacks a runner" condition is structurally impossible by construction (the pty driver reuses the same `runClaude*Scenario` functions verbatim). The parity guard correctly covers the one shared map.
2. **Pane resolver — does not mis-resolve.** Adversarial `list-panes` inputs (FO NOT at index 0 with an active ensign at 0; FO last among three; a teammate title embedding "first-officer" prose without the `spacedock:first-officer` marker) all resolved the FO pane correctly. Resolves by marker substring and returns the real index, not "0".
3. **Fork check — no fork.** Every "shared" fixture/assertion the pty test uses (`readmeRealisticLifecycle`, `entityFixture`, `locateEntity`, `someCommitNamesOnly`, `liveStageReportHeading`, `frontmatterField`, `isTeamCreate`, `isEnsignDispatch`, `detectWrongRootBoot`, `gradeTerminalTeardown`/`markerEmittedByAssistant`/`expectTerminalTeardownGrade`) is defined OUTSIDE the pty files and imported same-package. The `pty*`-prefixed helpers are genuinely new team-mode-only code (no pre-existing equivalent forked). `ptyAnyTeamHasMember` parses the IDENTICAL `members[].name` shape `claudeteam.MemberExists`/`teamConfig` reads — an independent on-disk oracle, NOT self-referential.
4. **Spike-hazard re-introduction — a HOLE EXISTS, but it is AC-5's documented live-only altitude, NOT a strength gap.** Re-introduced each of the three hazards in the audit checkout (active-pane read in `captureFOPane`; removed the `waitStableIdle` gate from `launchAndSend`; piped the launch through `| tee`). Each left the ENTIRE offline suite green AND `go vet -tags live` clean. So no OFFLINE test guards the driver-level wiring of the three hazard fixes. This is acknowledged by AC-5 itself: its offline proof is scoped to the `firstOfficerPaneIndex` SEAM, and it states "the live tests passing is the end-to-end proof" for the integration. The hazard fixes are real code and the live AC-3/AC-4 tests exercise them; the hole is the inherent live-only e2e altitude, the same boundary every `//go:build live` test carries — not a claimed-but-missing offline guard. Recorded, not blocking.

Minor non-material notes (no action required): `paneIsStableIdle` is a pure, offline-testable function that correctly rejects empty/whitespace/busy panes but has no dedicated offline test (only the live `waitStableIdle` exercises it); `ptyAnyTeamHasMember` uses case-insensitive `EqualFold` vs the binary's exact `==` (strictly more permissive, does not weaken the AC-3 proof).

## Stage Report: validation

- DONE: Verify AC-1/AC-2/AC-5/AC-6 OFFLINE in the worktree: `go test ./...` green (the `liveDriver` refactor is behavior-preserving), the parity meta-test (`TestSharedScenarioRunnerCoverage`, `-tags live`) fails if a shared scenario lacks a pty runner, `TestFirstOfficerPaneIndex` resolves the FO pane by title, and the two `//go:build live` tests SKIP-not-fatal without auth (`go test ./...` excludes them; `-tags live` with no credential skips).
  AC-1: `go test -count=1 ./internal/ensigncycle/` green (6.4s) + full `go test ./...` all `ok` (behavior-preserving refactor; both `claudeLiveRunner` and `ptyLiveDriver` carry `var _ liveDriver = …`). AC-2: `TestSharedScenarioRunnerCoverage` green under `-tags live`; refuted by removing a runner → REDs (see Feedback Cycles #1). AC-5: `TestFirstOfficerPaneIndex` 3/3 green under DEFAULT tags (two-pane case resolves FO=0 not active ensign=1); `launchAndSend` has no launch pipe and gates the send on `waitStableIdle` (code-read + hazard re-intro audit). AC-6: with an empty no-token HOME and no `ANTHROPIC_API_KEY`, both live tests SKIP (not fatal) with the `isolatedClaudeEnv` auth-missing message, package PASS; default `go test ./...` excludes them (build tag).
- DONE: Run the DETACHED adversarial audit on a THROWAWAY checkout (never the implementation worktree): refute the offline ACs. Record findings as a `### Feedback Cycles` entry naming the adversarial edit; "refuted nothing material" is a valid recorded outcome.
  See `### Feedback Cycles` above — four adversarial edits on a throwaway detached worktree at `e5aaef92` (torn down clean). Refuted nothing material: parity guard genuine, pane resolver does not mis-resolve, no fork, oracle independent. One documented hole (the three AC-5 hazard fixes have no offline guard) is AC-5's own declared live-only altitude, not a strength deficiency.
- DONE: Assess AC-3/AC-4 (live, CI-gated): confirm the two team-mode tests are correctly wired into `runtime-live-e2e.yml` (tmux installed, correct test selectors, behind the CI-E2E gate) and would actually assert comm-officer roster membership + `TERMINAL_TEARDOWN_BOUNDED` via the EXISTING `gradeTerminalTeardown`. Recommend PASSED-pending-CI-E2E. Do NOT mark FAILED for the local auth gap.
  Wiring confirmed: `runtime-live-e2e.yml` installs tmux (step "Install tmux", `tmux -V`) before the pty step; the pty step `-run 'TestLivePtyStandingResidencyInjectsCommOfficer|TestLivePtyEnsignCycleTeamTeardown'` selectors exactly match the two `func TestLivePty…` names (`go test -tags live -list` confirms exactly those two); behind the `CI-E2E`/`CI-E2E-OPUS` environment matrix (required-reviewer gate) with `ANTHROPIC_API_KEY` and `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` (the flag that exposes native team tools). AC-3 oracle is the on-disk team `config.json` `members[]` (independent of the test). AC-4 reuses the EXISTING `gradeTerminalTeardown`/`markerEmittedByAssistant` (defined in `streamwatch_test.go`/`teardown_grade*.go`, NOT reimplemented) — verified an honest discriminator: offline grader tests prove it PASSES on FO-authored marker emission and FAILS on silent give-up AND on a contract-Read of the marker string (the false-positive guard). Cannot green locally (stale benchmark-token 401); NOT marked FAILED.

### Summary

Validated the m4 terminal/pty live harness. AC-1/AC-2/AC-5/AC-6 are PROVEN OFFLINE by running the behavior (full suite green, parity meta-test green and refute-tested, pane-resolver table test green under default tags, live tests SKIP-not-fatal with no auth and are build-tag-excluded by default). The detached adversarial audit on a throwaway checkout refuted nothing material: the parity guard is a genuine discriminator, the pane resolver does not mis-resolve under adversarial `list-panes` inputs, no shared scenario/fixture/assertion/grader is forked (all imported same-package, the team-roster oracle parses the binary's own `members[]` shape), and the one hole found — no offline guard on the three AC-5 hazard fixes inside `launchAndSend`/`captureFOPane` — is AC-5's explicitly documented live-only altitude, not a strength gap. AC-3/AC-4 are live/CI-gated: the two team-mode tests are correctly wired into `runtime-live-e2e.yml` (tmux installed, exact `-run` selectors, `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`, behind the CI-E2E approval gate) and reuse the EXISTING honest teardown grader; they cannot green locally only because the operator benchmark-token is stale (401), so they are PASSED-pending-CI-E2E (NOT failed). Verdict: PASSED — AC-1/2/5/6 verified offline now; AC-3/AC-4 must show, in the CI-E2E gate run, the live team `config.json` carrying `comm-officer` in `members[]` (AC-3) and an FO-authored `TERMINAL_TEARDOWN_BOUNDED` that `gradeTerminalTeardown` greens over the captured session jsonl (AC-4).
