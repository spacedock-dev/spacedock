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

1. The greeting phase uses this neutral prompt:

   ```text
   Use $spacedock:first-officer for this whole run.

   Workflow directory: .
   This is an interactive startup check. Start the first officer, return its initial greeting, and stop for operator input.
   ```

   The prompt names the subject and stopping point but says nothing about PR state, sweeping, advancement, archival, gate rendering, or expected output. After this phase, the runner observes that the workflow Git HEAD and porcelain are unchanged, both entity files are byte-identical, the merged entity remains active, no archive exists, and the recording `gh` stub has not run. The greeting names the pending gate and offers `engage`; it does not render `Gate review:` or `Decision:`.

2. The engage phase runs against the same workflow with a separate host invocation and this operator action:

   ```text
   Use $spacedock:first-officer for this whole run.

   Workflow directory: .
   engage .
   ```

   The runner then observes the merged entity under `_archive/` with `status: done`, `verdict: PASSED`, and an empty `mod-block`; its active path is absent, and the recording `gh` stub has run. The gate entity remains unchanged, and the engage response now carries the gate review and decision prompt.

Use two one-shot invocations rather than a prompt that says "greet, then engage." A combined prompt would script the ordering under test and would expose only the final filesystem state. Reusing the same fixture across two invocations preserves a durable boundary snapshot and works through both existing host adapters. A live greeting paired only with an offline `state sweep` test is also insufficient: it would prove the binary but not the first officer's later engage behavior.

No spike needed: `TestBootIdentifyIsSideEffectFree` already proves the read-only identify mechanism; the shared Claude and Codex live adapters already execute this fixture; the existing pr-merge fixture already produces the terminal archived state. Implementation only composes these proven seams in sequence.

### Exact implementation surfaces

- `internal/ensigncycle/shared_fixtures_test.go`: replace the scripted prompt, add the short engage prompt, and rewrite stale S7b fixture prose.
- `internal/ensigncycle/shallow_boot_fixture_live_test.go`: make the stub `gh` record calls outside the workflow root.
- `internal/ensigncycle/shallow_boot_assert_test.go`: grade separate post-greeting and post-engage observations, including repository and entity snapshots.
- `internal/ensigncycle/shared_scenarios_negative_test.go`: add isolating positive and negative phase-transition cases.
- `internal/ensigncycle/claude_live_runner_test.go` and `codex_live_runner_test.go`: run greeting, assert its snapshot, then run engage on the same fixture and assert the terminal archive.
- `internal/ensigncycle/shared_scenarios_test.go`: update the scenario intent from before-greet advancement to greet-read-only/engage-advances.

No product or skill contract changes belong in this task. PR #480 already shipped the user-visible behavior, so no site documentation diff is required. Test comments and fixture text that say "S7b before-greet" change to "unchanged at greeting; advanced at first engage."

## Out of scope

- Changes to `first-officer-shared-core.md`, `status --boot --identify`, `state ready`, `state sweep`, or pr-merge semantics.
- Changes to other live scenarios, runtime adapters, or journey-metric scenario IDs.
- A persistent multi-turn transport abstraction; two existing one-shot invocations provide the required boundary evidence.

## Acceptance criteria

**AC-1 (value — the scenario measures the #480 boundary).** In both Claude and Codex live runs, the greeting phase leaves the fixture's Git HEAD, porcelain, gate entity, and merged-PR entity unchanged; leaves the merged entity active and unarchived; and makes zero calls to the recording `gh` stub. The later `engage .` phase on that same fixture calls `gh` and moves the merged entity to `_archive/` with `status: done`, `verdict: PASSED`, and an empty `mod-block`. Verified by: host-neutral durable-state assertions fed by both live runners. Independent baseline: changing only the prompt today makes the current oracle RED because it still demands before-greet archival; after reconciliation, an early archive or a missed engage advancement independently makes the oracle RED.

**AC-2 (the live prompt observes rather than prescribes).** The greeting prompt contains no PR, merge, sweep, advance, terminal, archive, gate-review, or expected-result instruction. The engage prompt contains only startup context and the operator action `engage .`. Verified by: an offline exact-prompt test that rejects each banned behavior cue and pins the engage prompt.

**AC-3 (the phase oracle is falsifiable).** Offline isolating cases reject each wrong trajectory: advancement during greeting, any greeting-phase entity or Git mutation, a greet-time `gh` call, no advancement after engage, an incomplete terminal state, a lingering `mod-block`, gate mutation/dispatch, and a missing engage-time `gh` call. A correct unchanged-greeting/terminal-after-engage trajectory passes. Verified by: `TestShallowBoot...` unit cases over synthetic observations.

**AC-4 (the greeting and gate timing remain coherent).** The greeting names the ready gate and offers `engage` without a full review; after engage, the gate entity remains byte-identical and the response includes `Gate review:` and `Decision:`. Verified by: the same two live phase results plus the durable no-dispatch checks; transcript prose is not used to prove PR mutation.

## Test plan

- **TDD red, offline (cheap):** first add the correct two-phase observation and neutral-prompt tests. The current `assertShallowBoot` rejects the unchanged greeting because it requires `mergedArchived == true`; record that focused failure before rewriting the oracle.
- **Offline green (cheap, under one second expected):** run `go test ./internal/ensigncycle -run 'TestShallowBoot'`. The positive trajectory passes, and each isolating negative fails for its intended reason.
- **Live green (expensive, two model invocations per host):** run `go test -tags live -count=1 ./internal/ensigncycle -run 'TestLive(Claude|Codex)SharedScenarios/shallow-boot$'`. Each host snapshots the same fixture after the neutral greeting, then explicitly engages and checks the archived terminal state. Preserve separate phase artifacts for diagnosis.
- **Repository gates:** run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.
- **No prose-only proof:** prompt scans prove only neutrality. Entity bytes, Git state, archive location, frontmatter, stub-call logs, and live host runs prove behavior.

## Stage Report: ideation

- DONE: Lock the #480 move-to-engage behavior as authoritative: greet-only boot stays read-only, with merged-PR advancement deferred to engage.
  PR #480 commit `0ba08c54` and current shared-core lines 14 and 26 establish the chosen boundary; the design changes only scenario code.
- DONE: Specify a neutral shallow-boot prompt plus falsifiable red-before/green-after tests that observe no greet mutation and the correct later engage behavior.
  The two exact prompts, pre/post durable snapshots, TDD red, live green, and isolating negative trajectories are specified above.
- DONE: Name exact code/test surfaces, keep scope minimal, and produce a complete ideation Stage Report with no prose-only behavioral proof.
  Seven test surfaces are named; all behavioral ACs rest on live execution and resulting Git/entity/archive/stub state.

### Summary

The design makes the live scenario follow #480: greeting identifies without mutation, and first engage advances the merged PR. Two host invocations over one fixture preserve the phase boundary without scripting it, while offline negatives and live durable-state checks make both early and missing advancement fail.
