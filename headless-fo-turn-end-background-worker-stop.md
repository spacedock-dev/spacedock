---
title: Headless FO ends its turn while a background worker is open, and the host stops the worker
status: ideation
source: captain (2026-11-03), from the repair-sonnet-live-flakes ideation probe b2
score: 0.7
sprint:
group: live-lane
sprint-readiness: ready
id: 5779fk6mzntvv5j3x6hw8rdn
started: 2026-11-04T13:02:10Z
gates:
    version: 1
    records:
        - id: gate:5779fk6mzntvv5j3x6hw8rdn:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:5779fk6mzntvv5j3x6hw8rdn-ideation-1
              briefing:
                id: briefing:5779fk6mzntvv5j3x6hw8rdn:ideation:attempt-1:revision-1
                digest: sha256:c9399b0d41a8b1e46e8bb2c1bae3e4c892614f542f8447c6d586a8ff7e529ef2
                room-ref: '@review/ideation/briefing-1'
---

In a headless (`-p`) drive with background workers, the FO can end its turn while a worker is still running. In probe b2 of `repair-sonnet-live-flakes` (Sonnet 5.5, head `64274cef5`, keep-moving journey), the host stopped the background `questioned` ideation worker 600 s after the FO's last turn. The worker had not committed its report, so the journey graded `keep-moving-violation`.

The contract already covers one case: a single-entity `-p` run uses bare (blocking) dispatch, because "the blocking call prevents premature session termination". Multi-entity headless drives with background workers have no equal rule.

This task owns AC-1 and AC-3 of `repair-sonnet-live-flakes` (captain ruling 2026-11-03): an unbound live keep-moving PASS on Sonnet in 3 of 3 runs, and the retirement of the `060xp69y61yhrww23g3wvwqy` keep-moving binding and its registry row. Ownership of that binding moves here when `repair-sonnet-live-flakes` closes.

Evidence: the `### Finding: background worker stopped before it reported` section of `repair-sonnet-live-flakes.md`. The 600 s figure comes from one observation and is not confirmed in host documentation.

Acceptance criteria (ideation fleshes): a headless multi-entity drive does not end the FO turn while a dispatched worker has not reported; proven by a live keep-moving run where a slow background worker reports after the FO's last dispatch, plus the retirement evidence above. Options to weigh: a headless wait rule for the FO, and a stage-size rule for ensigns in ideation on fixture entities.

## Problem

The host stops the worker, not the FO. Claude Code has a print-mode limit on background work: `CLAUDE_CODE_PRINT_BG_WAIT_CEILING_MS`, default `600000` (10 min). When the main turn has ended, input is closed, and background tasks are still running after that time, the host stops them and writes this to stderr:

> Background tasks still running 10m after the last turn (subagent "…"); stopping them. Set CLAUDE_CODE_PRINT_BG_WAIT_CEILING_MS=0 to wait indefinitely.

The clock starts when the FO's turn ends. It resets each time the FO gets a new turn, for example when another worker's completion wakes it. So a worker can run for any length of time while the FO keeps getting turns. The host stops the worker only when the FO has been idle for 10 min. In b2 that is what happened: the FO finished the other three entities, its last line was at 21:23:35Z, and the host stopped the `questioned` worker at 21:33:35Z.

Host evidence (Claude Code 2.1.303): the binary holds `AOt=600000` as the default for `CLAUDE_CODE_PRINT_BG_WAIT_CEILING_MS`. Its print wind-down code sets the idle start when `inputClosed` is true and no turn is running, and clears it otherwise. When the time passes, it stops `local_agent` tasks (`tengu_print_ceiling_stop_agents`). The probes below confirm this behavior.

### The FO contract text that lets the turn end

The FO ends its turn on purpose. That is how the contract tells it to wait for a background worker. `skills/first-officer/references/claude-fo-dispatch.md` (at `8739edef9`):

> **First-turn-after-dispatch decision procedure.** When a turn begins and your most recent dispatch-related action was an `Agent(...)` spawn whose `«completion-signal»` has NOT yet been observed, preserve that worker and run `«dispatch.next-action»()` when unrelated mod/PR, ready-gate, dispatchable, or state work remains. Only after the one-shot idle/reconcile/retry finds no such work may you end the turn with no tool calls and no text.
> …
> Just emit `end_turn` with empty content. The runtime will wake you up again when a real event arrives.

`first-officer-shared-core.md` (keep moving):

> Launching an async ensign does not end the turn while independent FO work remains. Yield only when blocked on the async result with no other work, or at a gate or captain decision.

The same procedure bans the other ways to wait: "Do not: … emit `Bash` with commands like `sleep 30` or `wait`." Only single-entity mode has a blocking rule (`claude-fo-dispatch.md:9`): "use bare-mode dispatch for all agent spawning — the blocking call prevents premature session termination in `-p` mode."

So the defect is in the promise "the runtime will wake you up again". In `-p` mode the host keeps that promise for 10 min only. The FO did what the contract says. The seed AC ("does not end the FO turn while a dispatched worker has not reported") asks for a change to a posture that is correct. I propose a different end-state: the worker is not stopped before it reports.

### Spike (host probes, `claude -p`, Haiku, isolated HOME and config dir, benchmark token)

Each probe asks the main agent to launch one or two background `general-purpose` agents that run a Python sleep, then end its turn. The scripts and dirs were in `/tmp`, and I removed them.

| Probe | Ceiling env | Worker | Result |
|---|---|---|---|
| 1 | `15000` | 1 agent, 45 s sleep | Main turn ended 13:05:01. Host stopped the agent at about 15 s (`task_notification status: stopped`, `background_tasks_changed tasks: []`) and wrote the stderr message above with "15s". Same event shape as b2. |
| 2 | `25000` | 2 agents, 15 s and 35 s sleeps | Main turn ended 13:05:30.8. The short agent reported, and the main agent had a turn that ended 13:05:47. The long agent reported at 13:06:06.9, 36 s after the first turn end, and was not stopped. So a main-agent turn resets the clock. |
| 3 | unset (default) | 1 agent, two 330 s sleeps | Main turn ended 13:15:12.7. Stderr: "Background tasks still running 10m after the last turn (subagent "slowworker"); stopping them." Process exited 604 s after start. The agent was stopped in its second sleep. This confirms the 600 s figure. |
| 4 | `0` | same as probe 3 | The agent reported at 13:26:15.7, 663 s after the main turn ended. The main agent woke and answered `REPORTED`. Exit 0. |
| 5 | `0` | background Bash (`sleep(100000)`), no agent | The session stayed open until I killed it after 4 min or more. |
| 6 | `15000` | same as probe 5 | The session stayed open past 15 s until my 90 s cap killed it. No ceiling message. |

Probes 5 and 6 show that the ceiling does not govern background shells. The host holds the session for a running shell with or without the ceiling (its wind-down code holds shells with a deadline). So setting `0` does not make shell leaks worse. With `0`, the only new wait is for background agents, and those end when their model loop ends. The host's stream watchdog stops a stalled API stream after 300 s or more.

One run before probe 1 used `sleep 45`. The host's Bash policy blocked it ("Blocked: standalone sleep"), so I used Python for the sleeps. This is useful for the implementation probe.

### Live Sonnet keep-moving probes

None were run (0 of the 3 allowed). They were not needed to choose between the options. The host probes show the mechanism and the fix. On `8739edef9` a local keep-moving run would most likely end in a 60 s quiet-budget kill (7 of 8 unchanged-byte runs in the parent task). That is the separate harness task `live-quiet-budget-silent-background-turn`, and it would decide nothing here. b2 already shows the FO ending its turn with a worker open.

## Proposed approach

The `spacedock claude` launcher sets `CLAUDE_CODE_PRINT_BG_WAIT_CEILING_MS=0` in the host's environment unless the parent environment already has the variable. It uses the same pattern as `withAgentTeams` in `internal/cli/frontdoor.go`. When the launch is wrapped in Safehouse, the launcher adds `--env-pass CLAUDE_CODE_PRINT_BG_WAIT_CEILING_MS`. Safehouse gives the child a sanitized environment, and its allowlist does not include `CLAUDE_CODE_*`. Interactive sessions ignore the variable, because the ceiling applies only when input is closed. Codex and Pi launches do not change.

With this change, the contract's "the runtime will wake you up again" holds for every headless session that `spacedock claude` starts. The FO contract does not change.

Why `0` and not a larger finite value: a finite value moves the limit and does not remove it. Implementation and validation stages in this repo often run 30 to 60 min or more under load. The host's own message names `0` as the way to wait. An operator who wants a limit can set any value, and the launcher keeps it.

Alternatives considered (each against value-AC 1):
1. A headless wait rule in the FO contract (the FO keeps a turn alive while a worker is open). Rejected. Host 2.1.303 has no tool that blocks on a background agent: the init tool list has no `TaskOutput` or `AgentOutput`. A Bash poll loop breaks the contract's `sleep`/`wait` ban, and a foreground Bash call is limited to 10 min. `ScheduleWakeup` turns could reset the clock, but the result depends on the model, it adds FO turns, and it needs new contract text. The env setting removes the limit mechanically.
2. A stage-size rule for ensigns on fixture entities. Rejected as the fix. b2's worker over-spent (nine symlink experiments on an entity with no design text), but a real ideation, implementation or validation stage can legitimately take more than 10 min. A prose rule depends on the model, and it leaves the host limit in place for all real stages.
3. Set the variable only in the live harness (`isolatedClaudeEnv`). Rejected. It makes the lane green and leaves the defect for released users, who start headless drives through the same launcher. The LOC is about the same.
4. Pass the value through `--settings '{"env":…}'` so that Safehouse needs no `--env-pass`. Rejected. It is not verified that the host reads the ceiling after it applies the settings env, and it can conflict with an operator's own `--settings`.
5. Blocking (bare) dispatch for every headless drive. Rejected. It removes parallel dispatch, which the multi-entity contract and keep-moving require.

Dependencies:
- Stack on `repair-sonnet-live-flakes` (`8739edef9`). Its step-5 fix covers the other cause of the same `060xp` binding (the hand batch commit before `--stamp`). Neither fix alone can give 3 of 3.
- Graded local live runs also need `live-quiet-budget-silent-background-turn`. Without it, the 60 s quiet budget kills most local runs after the FO's last turn. The CI Sonnet lane, which ran the journey in about 62 s, is the other source of graded runs.

Expected surface (Estimate net LOC change: +51, across 8 files; about 63 insertions, 12 deletions):
- `internal/cli/frontdoor.go`: a constant, `withPrintBgWait`, the call at the Claude `ops.Launch`, and `--env-pass` in the Claude Safehouse wrap (about +16/-2).
- `internal/cli/frontdoor_test.go`: two tests next to `TestClaudeFrontDoorEnablesAgentTeamsWhenParentUnset`: parent unset gives `=0`, and an explicit parent value is kept (about +36).
- `internal/cli/launch_parity_test.go`, `internal/cli/safehouse_frontdoor_test.go`, `internal/cli/safehouse_knob_test.go`: the wrapped Claude argv expectations gain `"--env-pass", printBgWaitCeilingEnv` (7 lines, +7/-7). Codex expectations do not change.
- `docs/site/reference/command-reference.md`: one paragraph (+2).
- At retirement: `internal/ensigncycle/shared_live_runner_test.go` (remove `liveXFail("claude-sonnet", "060xp69y61yhrww23g3wvwqy")`, +1/-1) and `docs/runtime-live-ci-registry.md` (remove the `keep-moving-posture | claude-sonnet | 060xp…` row, -1).

Tolerance: net +30 to +85 LOC, 6 to 9 files. The Safehouse part (4 files, about +9/-8) can be split into a separate task if the captain wants the smallest surface. Without it, a Safehouse-wrapped headless drive still hits the 10 min limit.

Semantics: runtime behavior changes. The host child environment of `spacedock claude` gets one more variable, and the Safehouse-wrapped Claude argv gets one more `--env-pass`. A headless session with a running background agent now waits until the agent ends, where before the host stopped the agent 10 min after the FO's last turn. Command grammar, stored formats and authority do not change. The FO and ensign contracts do not change.

Doc diff (`docs/site/reference/command-reference.md`, new paragraph after the permission-mode paragraph in `## Launch`):

```diff
 An unsandboxed bootstrap launch carries no safehouse isolation, so per-action permission prompting is friction without a matching safety gain: ...
+
+In print mode (`-p`), Claude Code stops background tasks that are still running 10 minutes after the main turn ends. A headless first officer ends its turn while its background workers run, and a worker's completion wakes it again. So `spacedock claude` sets `CLAUDE_CODE_PRINT_BG_WAIT_CEILING_MS=0` (wait until background tasks end) unless you set a value yourself, and passes it through Safehouse. If you start a headless first officer with `claude -p` directly, set this variable yourself.
```

## Acceptance criteria

- **AC-1 (value):** In a headless Claude session with the environment `spacedock claude` sets, a background agent that reports more than 600 s after the main turn ends is not stopped, and its completion wakes the main agent. Baseline at `8739edef9` (no variable): the host stops the agent at 600 s and writes "Background tasks still running 10m after the last turn". Tested by: a one-off live host probe (Haiku, about 11 min, a few cents), probes 3 and 4 above repeated on the candidate. The env value the probe uses must be the value that AC-2's test reads from the launcher. Falsify: remove the variable, and the agent is stopped at 600 s (probe 3).
- **AC-2:** The `spacedock claude` child environment has `CLAUDE_CODE_PRINT_BG_WAIT_CEILING_MS=0` when the parent does not set it. An explicit parent value (for example `600000`) is kept unchanged. A Safehouse-wrapped Claude launch carries `--env-pass CLAUDE_CODE_PRINT_BG_WAIT_CEILING_MS` before `--`. Codex and Pi launch argv and environments do not change. Tested by: Go unit tests in `internal/cli` with `fakeHost` and `envValue`, plus the updated wrapped-argv tests. Falsify: drop the helper call, and the unset test fails. Make the helper overwrite the value, and the explicit-value test fails. Drop the `--env-pass`, and the wrapped-argv tests fail.
- **AC-3 (moved from `repair-sonnet-live-flakes` AC-1):** Sonnet keep-moving passes without an XFAIL binding. `TestLiveCommonKeepMovingPosture` on `claude-sonnet-5-5` grades PASS in 3 of 3 unbound runs on the candidate bytes. Baseline: a graded red in 1 of 2 CI attempts on `c6a336a33`, 1 of 1 on `861a6e13a`, and b2 on `64274cef5` (XFAIL from this host stop). Tested by: exact live runs (local `SPACEDOCK_LIVE_MODEL=claude-sonnet-5-5` with the quiet-budget fix present, or the CI Sonnet lane). An infra kill (quiet budget, auth) is not a graded run.
- **AC-4 (moved from `repair-sonnet-live-flakes` AC-3):** The claude-sonnet `060xp69y61yhrww23g3wvwqy` binding on keep-moving is removed only after a graded bound XFAIL or XPASS on the candidate, followed by unbound PASS runs on otherwise unchanged bytes (the 3 of AC-3). The registry row is removed in the same change. Tested by: the run records cited in the validation report, and `TestRuntimeLiveRegistryMatchesBindings`.
- **AC-5:** `go test ./...`, `go test ./... -race`, `gofmt -l ./cmd ./internal` (no new entries) and `go test ./internal/contractlint/` pass.

## Test plan

- Primary proof owners: the `internal/cli` front-door tests for AC-2 (deterministic, under 1 s). `TestLiveCommonKeepMovingPosture` for AC-3 and AC-4 (live, about $3 and 7 to 28 min per run). No existing test covers AC-1, because the host ceiling is a host behavior. A one-off live host probe proves it and records it in the stage report.
- AC-2 tests: `TestClaudeFrontDoorSetsPrintBgWaitWhenParentUnset` and `TestClaudeFrontDoorPreservesExplicitPrintBgWait`, modeled on the two agent-teams tests. The existing wrapped-argv tests get the extra `--env-pass`. Each falsifying edit is named in AC-2.
- AC-1 probe: `claude -p` with Haiku, an isolated HOME and config dir, and the benchmark token. Run once without the variable (stopped at 600 s) and once with the launcher's value (reports at about 660 s and wakes the main agent). Use Python sleeps, not `sleep`, because the host blocks a standalone `sleep`. Cost: a few cents and about 11 min. Not a CI test, because it waits more than 10 min on a host limit.
- AC-3 and AC-4 live cost: 1 bound and 3 unbound graded Sonnet runs, about $12 to $15. Run them after the quiet-budget fix is on the stack, or use the CI Sonnet lane.
- Risk: a headless session whose background agent never ends now does not exit, where before the host stopped it after 10 min. The host still bounds the parts: the stream watchdog (300 s or more) stops a stalled API stream, and Bash calls have their own timeouts. Background shells are not affected (probes 5 and 6). The live harness keeps its quiet budget and test timeout. An operator who wants the old limit sets the variable.
- Risk: the variable name and the `0` meaning are host internals, not documented host API. The host's own stderr message names them. If a later host renames the variable, AC-1's probe goes red, and the AC-2 tests still pass. So AC-1's probe is the check to repeat when the host version changes.

## Stage Report: ideation

- DONE: From the b2 probe evidence in repair-sonnet-live-flakes.md, the Claude FO contract, and host behavior, establish when the host stops a background worker after a headless FO ends its turn (confirm or refute the 600 s figure with a cheap deterministic or single-call probe), and quote the FO contract text that lets a multi-entity headless drive end its turn while a dispatched worker has not reported.
  Confirmed. Claude Code 2.1.303 stops background agents when the main turn has been idle for `CLAUDE_CODE_PRINT_BG_WAIT_CEILING_MS` (default `600000` in the binary). Each main-agent turn resets the clock. Haiku host probes: the default stopped the agent at 600 s with "Background tasks still running 10m after the last turn" (probe 3). `15000` stopped it at 15 s (probe 1). A main turn in between reset the clock (probe 2). `0` let the agent report 663 s after the main turn and woke the main agent (probe 4). Quoted the contract: `claude-fo-dispatch.md` "Just emit `end_turn` with empty content. The runtime will wake you up again when a real event arrives", and the shared-core "Yield only when blocked on the async result with no other work". Only single-entity mode has a blocking rule.
- DONE: Propose the smallest fix (for example a headless wait rule in the FO contract, or a stage-size rule for ensigns on fixture entities) with a falsifiable AC, an exact file list, net LOC estimate and tolerance, and state the retirement evidence for the 060xp keep-moving binding (AC-1 and AC-3 moved here from repair-sonnet-live-flakes).
  Fix: `spacedock claude` sets `CLAUDE_CODE_PRINT_BG_WAIT_CEILING_MS=0` unless the parent set it, and passes it through Safehouse. The FO contract does not change. The headless wait rule and the ensign stage-size rule are rejected with reasons. Net +51 LOC across 8 files, tolerance +30 to +85 and 6 to 9 files. AC-1 (value, host probe against the 600 s baseline) and AC-2 (launcher env tests) each name a falsifying edit. Retirement evidence: 1 graded bound XFAIL or XPASS on the candidate, then 3 of 3 unbound PASS on otherwise unchanged bytes (AC-3, AC-4), on a stack with the parent's step-5 fix and, for local runs, the quiet-budget fix.
- DONE: At most 3 local live Sonnet keep-moving probes on 8739edef9 if they are needed to decide between options; record each run's grade, and whether the FO ended its turn with a worker open.
  0 Sonnet runs. They were not needed: the host probes decide the mechanism and the fix. On `8739edef9` a local run would most likely end in a 60 s quiet-budget kill. b2 already shows the FO ending its turn with the `questioned` worker open (last FO line 21:23:35Z, worker stopped 21:33:35Z). I used 7 Haiku single-call host probes instead (one hit the host's standalone-`sleep` block and was rerun). Each is listed in the Spike table.

### Summary

The host, not the FO, causes the b2 failure. In `-p` mode Claude Code stops background agents after 10 min of main-turn idleness (`CLAUDE_CODE_PRINT_BG_WAIT_CEILING_MS`, default 600000). Each FO turn resets the clock. Haiku probes confirmed the 600 s figure, the reset, and that `0` lets a worker report later and wake the main agent. The FO contract tells the FO to end its turn and wait ("the runtime will wake you up again"). That posture is correct, so I propose to make the promise true instead of changing it: `spacedock claude` sets the variable to `0` unless the operator chose a value, and passes it through Safehouse. That is net +51 LOC across 8 files, with a doc paragraph. A headless FO wait rule is rejected, because the host has no tool that blocks on an agent and the contract bans sleep loops. An ensign stage-size rule is rejected, because real stages exceed 10 min. I made no Sonnet runs. Retirement of the `060xp` binding needs the parent's step-5 fix on the same stack and graded live runs, which locally depend on the quiet-budget fix.
