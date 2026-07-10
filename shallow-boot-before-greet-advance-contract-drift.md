---
title: "shallow-boot live scenario asserts before-greet merged-PR advancement that the #480 Startup restructure moved to «engage»"
status: ideation
source: "Science-officer finding during the dp (fo-deferred-load-point-hunt) shallow-boot investigation, 2026-07; captain-approved to file separately. The shallow-boot live scenario's prompt (internal/ensigncycle/shared_fixtures_test.go:430) scripts advance any merged PR per the before-greet merged-PR sweep and report the merged-PR entity as advanced, and its assertion (shallow_boot_assert_test.go:85-94) requires before-greet advancement. The current FO contract says the opposite: pr-merge advancement fires at first engage, not the greet — advanced at engage or never (first-officer-shared-core.md:26) — and a greet-and-stop boot writes nothing (Startup step 2). Git history: the scenario (fe4261be, PR #365, 2026-06-13) predates and was faithful to the pre-restructure contract's boot-time «state.sweep-merged» step; the #480 Startup restructure (0ba08c54, 2026-07-07) moved advancement to engage. So the scenario now passes only by scripting compliance with an instruction the contract no longer prescribes, and its sibling boot-subject scenarios use a neutral prompt instead (see below)."
started: 2026-07-10T12:52:14Z
completed:
verdict:
score: 0.5
worktree:
issue:
id: m3y296hx6tb939qqq5zcxphw
---

## Problem

The `shallow-boot` live scenario scripts the pre-#480 behavior it claims to observe. `shallowBootPrompt` tells the first officer to "advance any merged PR per the before-greet merged-PR sweep" and report the entity as advanced (`internal/ensigncycle/shared_fixtures_test.go:425-433`). `assertShallowBoot` then requires a terminal, archived merged-PR entity before the greeting (`internal/ensigncycle/shallow_boot_assert_test.go:48-97`). A model can pass by following the stale prompt even when it ignores the current startup contract.

PR #480 (`0ba08c54`) settled the direction: `status --boot --identify --json` performs local, read-only identification, and first `«engage»` runs `state ready` followed by `state sweep`. The contract says both that a greet-only session writes nothing (`skills/first-officer/references/first-officer-shared-core.md:14`) and that pr-merge advancement fires at first engage, not the greeting (`:26`). The scenario predates that change (`fe4261be`) and now asserts the opposite boundary. It also asks for a full gate review during the greeting, although #480 defers gate rendering to engage; a neutral prompt exposes both stale expectations.

## Proposed approach

Keep #480 authoritative and turn `shallow-boot` into one two-phase scenario over one fixture.

1. **Claude uses the existing interactive default and continuation seam.** Route only `shallow-boot` through `ptyLiveDriver`. Split `launchAndSend` after its existing `waitTranscriptIdle` gate so phase one launches bare `spacedock claude --plugin-dir ...` and observes the initial idle before any scenario text is sent. The Spacedock front door already selects `--agent spacedock:first-officer` and appends its normal `bootstrapPrompt`; the test supplies no greeting prompt. Snapshot Git HEAD and porcelain, both entity files, archive absence, and the recording `gh` log at that idle boundary. Capture the greeting that names the pending gate and offers `engage` without `Gate review:` or `Decision:`. Then send exactly `engage .` through the existing `sendToFO` input seam, wait for a fresh committed idle, and take the engage snapshot from the same resident session and fixture.

2. **Codex uses the normal front-door bootstrap with the smallest unavoidable headless exception.** The current `codexLiveRunner` has no interactive/session handle: `codexExecArgv` requires a positional prompt, and `run` starts, drains, and returns from a fresh `codex exec` process. Invoke that process through `spacedock codex` so `codexBootstrapPrompt` remains the authoritative normal startup context. Phase one appends only this task:

   ```text
   Stop after the greeting.
   ```

   This clause is required because `codex exec` is headless and the normal FO contract otherwise drives through first engage instead of stopping at the greeting. It carries no workflow path or PR, merge, sweep, advance, archive, gate, test, or expected-result wording; `--cd` supplies the workflow root. After the greeting snapshot, run a second front-door `codex exec` against the same fixture with only the operator task `engage .`; Codex has no existing continuation/input seam to reuse.

After either host's engage phase, observe the merged entity under `_archive/` with `status: done`, `verdict: PASSED`, and an empty `mod-block`; its active path is absent, the recording `gh` stub has run, the gate entity is byte-identical, and the engage response carries the gate review and decision prompt. Do not replace the two phase boundary with one prompt that says “greet, then engage”: that would prescribe the ordering and expose only final state. Do not build a new Codex interactive transport solely for this scenario; the minimal headless stop clause is smaller and its limitation is explicit.

No spike needed: `TestBootIdentifyIsSideEffectFree` proves read-only identify; `ptyLiveDriver.launchAndSend`, `sendToFO`, and `nudgePastGreet` already exercise Claude's default interactive launch and later-input seam; front-door argv tests pin Codex's default bootstrap/task composition; and the current Codex runner already captures separate exec artifacts. Implementation composes these proven seams, while the new phase behavior itself remains for TDD and live validation.

### Exact implementation surfaces

- `internal/ensigncycle/shared_fixtures_test.go`: remove the scripted shallow-boot prompt and rewrite stale S7b fixture prose; no host-neutral greeting prompt replaces it.
- `internal/ensigncycle/shallow_boot_fixture_live_test.go`: make the stub `gh` record calls outside the workflow root.
- `internal/ensigncycle/shallow_boot_assert_test.go`: grade separate post-greeting and post-engage observations, including repository and entity snapshots.
- `internal/ensigncycle/shared_scenarios_negative_test.go`: add isolating positive and negative phase-transition cases.
- `internal/ensigncycle/pty_live_driver_test.go`: split default interactive launch-to-idle from later input, preserving `launchAndSend` for existing callers and exposing the existing same-session send/fresh-idle seam.
- `internal/ensigncycle/claude_live_runner_test.go`: route `shallow-boot` through that default interactive PTY seam, assert the pre-input greeting snapshot, then send exactly `engage .` and assert the terminal archive.
- `internal/ensigncycle/codex_live_runner_test.go`: retain two headless processes but launch them through `spacedock codex`, using normal bootstrap plus only `Stop after the greeting.` and then normal bootstrap plus only `engage .`.
- `internal/ensigncycle/shared_scenarios_test.go`: update the scenario intent from before-greet advancement to greet-read-only/engage-advances.

No product or skill contract changes belong in this task. PR #480 already shipped the user-visible behavior, so no site documentation diff is required. Test comments and fixture text that say "S7b before-greet" change to "unchanged at greeting; advanced at first engage."

Integration sequence: m3 overlaps `internal/ensigncycle/shared_fixtures_test.go` with c3 / PR #493. Wait for #493 to land, then branch or rebase onto updated `origin/main` before editing any implementation surface above.

## Out of scope

- Changes to `first-officer-shared-core.md`, `status --boot --identify`, `state ready`, `state sweep`, or pr-merge semantics.
- Changes to other live scenarios, runtime adapters, or journey-metric scenario IDs.
- A new Codex interactive/PTY transport or general multi-turn abstraction; only Claude reuses an existing resident-session seam.
- Changes to the Claude or Codex launcher bootstrap prompts; the scenario consumes those defaults rather than redefining them.

## Acceptance criteria

**AC-1 (value — the scenario measures the #480 boundary).** In both Claude and Codex live runs, the greeting phase leaves the fixture's Git HEAD, porcelain, gate entity, and merged-PR entity unchanged; leaves the merged entity active and unarchived; and makes zero calls to the recording `gh` stub. The later `engage .` phase on that same fixture calls `gh` and moves the merged entity to `_archive/` with `status: done`, `verdict: PASSED`, and an empty `mod-block`. Verified by: host-neutral durable-state assertions fed by both live runners. Independent baseline: the current oracle is RED when fed the required unchanged greeting observation because it still demands before-greet archival; after reconciliation, an early archive or a missed engage advancement independently makes the oracle RED.

**AC-2 (default startup observes rather than prescribes).** Claude phase one supplies no scenario prompt: the normal `spacedock claude` bootstrap reaches its interactive greeting, and phase two sends exactly `engage .` to that session. Codex phase one uses the normal `spacedock codex` bootstrap plus exactly `Stop after the greeting.`—the minimum required by its headless, prompt-mandatory, no-session-handle runner—and phase two uses the normal bootstrap plus exactly `engage .`. No phase task names PR, merge, sweep, advance, terminal, archive, gate review, test intent, or expected state. Verified by: offline argv/default-seam tests plus exact task/input assertions; behavioral correctness remains the durable/live proof in AC-1, AC-3, and AC-4.

**AC-3 (the phase oracle is falsifiable).** Offline isolating cases reject each wrong trajectory: advancement during greeting, any greeting-phase entity or Git mutation, a greet-time `gh` call, no advancement after engage, an incomplete terminal state, a lingering `mod-block`, gate mutation/dispatch, and a missing engage-time `gh` call. A correct unchanged-greeting/terminal-after-engage trajectory passes. Verified by: `TestShallowBoot...` unit cases over synthetic observations.

**AC-4 (the greeting and gate timing remain coherent).** The greeting names the ready gate and offers `engage` without a full review; after engage, the gate entity remains byte-identical and the response includes `Gate review:` and `Decision:`. Verified by: the same two live phase results plus the durable no-dispatch checks; transcript prose is not used to prove PR mutation.

## Test plan

- **TDD red, offline (cheap):** first add the correct two-phase observation tests and exact task/input tests: Claude has no greeting task and later receives `engage .`; Codex uses `Stop after the greeting.` and later `engage .`. Existing launch-parity tests establish both front-door defaults. The current `assertShallowBoot` rejects the unchanged greeting because it requires `mergedArchived == true`; record that focused failure before rewriting the oracle.
- **Offline green (cheap, under one second expected):** run `go test ./internal/ensigncycle -run 'TestShallowBoot'`. The positive trajectory passes, and each isolating negative fails for its intended reason.
- **Live green (expensive):** run `go test -tags live -count=1 ./internal/ensigncycle -run 'TestLive(Claude|Codex)SharedScenarios/shallow-boot$'`. Claude snapshots the default greeting and sends `engage .` in one interactive session; Codex snapshots after its minimal headless greet-stop task and then runs a fresh normal-bootstrap `engage .` process against the same fixture. Preserve separate phase artifacts for diagnosis.
- **Repository gates:** run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.
- **No prose-only proof:** launch/task/input assertions prove only default-seam use and neutrality. Entity bytes, Git state, archive location, frontmatter, stub-call logs, and live host runs prove behavior.

### Feedback Cycles

- Cycle 1 (ideation gate): REVISE. The captain agrees with the two-phase durable-state test but rejects the unexplained custom greeting prompt: use the host/runtime's existing default first-officer startup for phase one if that seam exists, then use only the operator action `engage .` for phase two. Audit the actual live-runner launch surface and remove the custom “interactive startup check” wording; retain any custom text only if technically unavoidable and demonstrate why the default cannot express the test.

## Stage Report: ideation

- DONE: Lock the #480 move-to-engage behavior as authoritative: greet-only boot stays read-only, with merged-PR advancement deferred to engage.
  PR #480 commit `0ba08c54` and current shared-core lines 14 and 26 establish the chosen boundary; the design changes only scenario code.
- DONE: Specify a neutral shallow-boot prompt plus falsifiable red-before/green-after tests that observe no greet mutation and the correct later engage behavior.
  The two exact prompts, pre/post durable snapshots, TDD red, live green, and isolating negative trajectories are specified above.
- DONE: Name exact code/test surfaces, keep scope minimal, and produce a complete ideation Stage Report with no prose-only behavioral proof.
  Seven test surfaces are named; all behavioral ACs rest on live execution and resulting Git/entity/archive/stub state.

### Summary

The design makes the live scenario follow #480: greeting identifies without mutation, and first engage advances the merged PR. Two host invocations over one fixture preserve the phase boundary without scripting it, while offline negatives and live durable-state checks make both early and missing advancement fail.

## Stage Report: ideation (cycle 2)

- DONE: Re-audit the current m3 design against PR #480's greet-read-only / engage-advances contract and retain the existing two-phase, durable-state approach. Do not change product or test code.
  Source audit confirms `0ba08c54` makes boot identify local/read-only and moves `state ready` plus `state sweep` to first `engage`; the existing two-invocation design remains minimal, and this cycle changes only the state report.
- DONE: Append an ideation cycle-2 Stage Report that explicitly maps AC-1 through AC-4 to the proposed verification surfaces and names which evidence is already established versus intentionally deferred to implementation/validation. Do not claim unrun live or offline tests passed.
  The four mappings below distinguish source/fixture evidence established during ideation from every task-specific two-phase offline and live result deferred to later stages; no such test pass is claimed.
- DONE: Record the integration constraint: m3 overlaps `internal/ensigncycle/shared_fixtures_test.go` with c3 / PR #493, so implementation must wait for #493 to land, then branch or rebase onto the updated origin/main before editing. Keep the design otherwise minimal and review-ready.
  Implementation is explicitly sequenced after #493 and must begin from updated `origin/main`; no additional surface or behavior is added to the design.
- DONE: AC-1 maps to host-neutral durable-state assertions supplied by both Claude and Codex live runners across two invocations of the same fixture.
  Established: #480's contract, read-only identify coverage, fixture/archive seams, and the current stale oracle are present in source. Deferred: the new phase snapshots, engage advancement assertions, and both live-host executions.
- DONE: AC-2 maps to an offline exact-prompt test that pins the engage prompt and rejects every banned behavior cue in the neutral greeting prompt.
  Established: the current fixture visibly contains the stale before-greet sweep and gate-rendering cues. Deferred: replacing the prompts and running the prompt-neutrality test.
- DONE: AC-3 maps to offline synthetic-observation cases covering the correct trajectory and each isolated early-mutation, missing-advancement, incomplete-terminal, gate-mutation/dispatch, and `gh`-call failure.
  Established: the current source oracle requires pre-greet archival and therefore encodes the contract drift. Deferred: the TDD red capture, the rewritten phase oracle, all isolating cases, and their focused execution.
- DONE: AC-4 maps to the two live phase responses plus byte-identical gate state and durable no-dispatch evidence, with transcript text used only for greeting/review timing.
  Established: #480's shared-core requires the greet to name the gate and offer `engage` without rendering a review. Deferred: live confirmation that engage renders `Gate review:` / `Decision:` while leaving the gate entity unchanged.

### Summary

The cycle-2 audit retains the two-phase, same-fixture design because it directly measures #480's durable greet/engage boundary without prescribing the result. Implementation and all task-specific AC test execution remain intentionally deferred, and m3 must wait for c3 / PR #493 to land before branching or rebasing onto updated `origin/main` and editing the shared fixture.

## Stage Report: ideation (cycle 3)

- DONE: Keep the approved two-phase same-fixture design while responding to the captain's anchored default-startup feedback.
  The design still snapshots greeting state before a later `engage .` over one durable fixture; only the host launch/input seams changed.
- DONE: Audit the actual Claude live-runner launch API and use the normal/default first-officer startup for phase one wherever technically possible.
  `ptyLiveDriver` already launches bare interactive `spacedock claude`, waits for its default bootstrap greeting to become idle before its first send, and exposes `sendToFO`; Claude now supplies no phase-one scenario prompt.
- DONE: Audit the actual Codex live-runner launch API and prove any unavoidable custom greeting instruction from code.
  `codexLiveRunner.run` starts and drains a fresh `codex exec`, `codexExecArgv` requires a positional prompt, and no session/input handle survives; the sole exception is exact task `Stop after the greeting.` appended by the normal Spacedock Codex front door.
- DONE: Make phase two carry only normal startup context plus operator action `engage .`, or use an existing continuation/input seam.
  Claude sends exact same-session input `engage .`; Codex runs a second front-door exec whose only task beyond the normal bootstrap is exact `engage .`.
- DONE: Remove “interactive startup check” and other test-directing wording, and update Proposed approach, AC-2, test plan, and exact surfaces.
  No phase text names the workflow path, PR state, merge/sweep/advance/archive behavior, gate timing, test intent, or expected result; the revised sections pin the host-specific seams and durable proof.
- DONE: Preserve the #493 sequencing constraint and do not claim unrun tests passed.
  Implementation remains blocked on c3 / PR #493 landing followed by branch/rebase onto updated `origin/main`; task-specific offline and live verification is explicitly deferred.
- DONE: Map the captain's feedback to a minimal review-ready choice rather than expanding transport scope.
  Reusing Claude PTY continuation satisfies the default-startup request; building a new Codex PTY was rejected as unnecessary because one minimal headless stop clause plus durable phase snapshots proves the boundary.

### Summary

Cycle 3 removes the custom Claude greeting entirely and binds the scenario to the launcher's real default greeting plus existing continuation input. Codex's code-proven headless limitation is isolated to `Stop after the greeting.`, while its engage phase uses only the normal bootstrap plus `engage .`; no product or test code was changed and no task-specific test result is claimed.
