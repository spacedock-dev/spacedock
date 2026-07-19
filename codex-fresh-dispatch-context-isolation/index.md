---
id: rt8eywbyf3d7nyc9bsmrymnq
title: Codex fresh dispatch must isolate parent turns
status: implementation
source: "/tmp/spacedock-codex-fresh-dispatch-incident.md; captain-requested legitimacy check and filing, 2026-07-19"
started: 2026-07-19T09:05:43Z
completed:
verdict:
score: "0.92"
worktree: .worktrees/spacedock-ensign-codex-fresh-dispatch-context-isolation
issue:
sprint:
---

## Problem

Spacedock promises that a fresh validation worker arrives without the maker's
reasoning. Codex fresh dispatch currently guarantees only a new handle. The
adapter omits `fork_turns`, its unit test requires that omission, and the First
Officer contract leaves the value open. A caller can therefore create a fresh
handle with inherited parent turns and silently defeat review independence.

The 2026-07-19 incident exercised that failure directly. The First Officer
spawned the `automate-beta-release` validator with `fork_turns: "all"`; the new
worker inherited the First Officer's conversation. The validator still found a
real private-checkout authentication defect. The incident weakens the
independence evidence, not the finding, so release validation needs one isolated
replacement run.

This is a review-integrity correctness fix. It does not add a confidentiality
boundary between agents that share a filesystem, process environment, or
repository. It prevents accidental parent-conversation inheritance at the
fresh-dispatch boundary.

## Evidence and the #414 gap

### Proven historical facts

- [PR #414](https://github.com/spacedock-dev/spacedock/pull/414), merged as
  `97eff5ff`, introduced `CodexMultiAgentV2Spawn` and the current Codex tool
  binding. Its design said inherited context should be deliberate and that
  file-pointer dispatch should normally have none.
- The implementation made `ForkTurns` optional. `ToolArgs` adds `fork_turns`
  only when the field is non-empty.
- `TestCodexMultiAgentV2SpawnInputMapsBuildOutput` requires `ForkTurns == ""`.
  The test checks task-name sanitation, prompt identity, and omitted unsupported
  helper fields; it never checks context isolation.
- #414's validation checked the mapping tests and a live feedback-rejection
  scenario. That scenario proved deliberate reuse of the same validation
  worker after rejection. It did not inspect a fresh spawn's `fork_turns` value
  or the child-visible parent context. The report's caveat concerned the
  `send_input` event label, not spawn isolation.

### Proven current facts

- Current upstream Codex parses an absent or blank `fork_turns` as `"all"`,
  maps `"all"` to full history, and maps `"none"` to no fork. See
  [`SpawnAgentArgs::fork_mode`](https://github.com/openai/codex/blob/main/codex-rs/core/src/tools/handlers/multi_agents_v2/spawn.rs).
- The incident record proves the harmful behavior for an explicit `"all"`
  spawn: the replacement handle inherited the First Officer's turns.
- Spacedock already uses a different operation for deliberate continuity:
  `followup_task` addresses the existing worker handle. Fresh dispatch uses
  `spawn_agent`.

### Unproven historical claim

No surviving #414 artifact proves how the exact Codex build used on 2026-06-20
interpreted an omitted `fork_turns`. The design, tests, and live report never
observed that value or the child's starting context. We must not retroactively
describe omission as full-history inheritance in that historical run.

The fix does not need that claim. An explicit `"none"` preserves isolation
whether Codex's default is `"all"`, `"none"`, or changes again.

## Proposed approach

Make one invariant own the boundary:

> Every Spacedock Codex `spawn_agent` call is a fresh dispatch and carries
> `fork_turns: "none"`. Deliberate continuity never calls `spawn_agent`; it uses
> `followup_task` with the existing handle.

Implement the invariant at the narrowest executable seam:

1. Remove `ForkTurns` from `CodexMultiAgentV2Spawn`. Make `ToolArgs` always add
   `"fork_turns": "none"`. The adapter then has no state or input channel that
   can emit omission, `"all"`, or a numeric turn count.
2. Bind the Codex First Officer's `«worker.spawn»` call to the same literal
   value. Keep the generic tool-shape probe separate from the invocation policy:
   the probe discovers that `fork_turns` exists; the runtime binding always
   invokes it with `"none"`.
3. Define fresh dispatch in the host-neutral dispatch core as a new handle with
   no inherited parent turns. Preserve the existing reuse decision and
   `followup_task` route. A stage without `fresh: true` may still advance on the
   same handle when all reuse conditions pass.
4. Add one compact adapter test that closes every adapter escape. Retain a
   one-off Codex canary probe as host-behavior evidence; do not add a bespoke
   live harness for this three-key map invariant.

No spike is needed. The incident proves the unsafe explicit-`all` outcome, and
current upstream Codex source proves that `"none"` selects no fork. The retained
one-off probe confirms that behavior on the live host and separately confirms
that `followup_task` preserves an existing worker's context.

## Why this is the smallest sufficient mechanism

- Setting the struct field's default to `"none"` while retaining a mutable
  `ForkTurns` field leaves an internal path back to `"all"` or a number. Removing
  the field makes unsafe output unrepresentable and serves AC-1 and AC-2.
- Adding a workflow or user-facing fork setting creates a choice Spacedock does
  not need. `spawn_agent` and `followup_task` already express isolation and
  continuity. The extra setting would weaken AC-1.
- Marking every implementation stage `fresh: true` changes workflow reuse
  policy, not spawn semantics. It would increase cost and still leave any other
  fresh spawn vulnerable. AC-3 preserves the existing stage policy.
- Contract wording alone cannot prove behavior. The compact adapter test proves
  the shipped map invariant, while the retained one-off canary proves the live
  host meaning of `"none"`; no second lifecycle harness is needed.
- Replacing the incident's `automate-beta-release` validation is follow-up work
  in that workflow after this fix is available. It is not part of this adapter
  patch's acceptance boundary.

The `per-host-stage-model-override` task is compatible, not duplicative. Its
current design forces `"none"` only when model or effort overrides need it and
records plain-spawn isolation as later work. When rebased, that task should
forward overrides onto this invariant and must not restore a mutable fork mode.

## Acceptance criteria

- **AC-1 (VALUE): Every Spacedock Codex fresh-spawn payload is explicitly
  isolated.** `ToolArgs` contains exactly `fork_turns: "none"` for plain and
  override-bearing helper envelopes. Verified by exact map equality in the
  compact adapter test and by the retained one-off Codex probe, where the actual
  spawn used `"none"` and the child rollout lacked the exact parent-only canary.
- **AC-2: Unsafe fork output is unrepresentable in the Codex adapter.**
  `CodexMultiAgentV2Spawn` exposes no `ForkTurns` field, and absent, `"all"`, or
  numeric helper fields cannot change the three-key output map. Verified by the
  same table test plus a reflection guard; restoring conditional omission or a
  mutable field makes the test fail.
- **AC-3: Fresh isolation does not change deliberate continuity or stage
  selection.** The host-neutral contract defines a fresh handle as inheriting no
  parent turns, the Codex binding supplies `"none"`, and eligible reuse remains
  on `followup_task`; no stage becomes globally fresh. Verified by the focused
  diff and the retained one-off probe, whose follow-up addressed the same task
  and thread and recalled its prior-turn marker.
- **AC-4: The patch is pre-merge release-safe without claiming post-merge
  publication state.** The self-contained commit passes focused, full, race, and
  existing no-rewind release guards; `edge-advance-decision` returns `skip` for
  v0.25.1 against 0.26.0-pre1. Landing on `main`, cherry-picking to `next`,
  publishing edge, cutting the annotated tag, and observing remote branch tips
  are post-merge release ceremony, not implementation acceptance.

## Scope

Do not force every implementation stage to be fresh. A workflow that requires
role separation should declare `fresh: true`; ideation-to-implementation reuse
remains valid when it does not. Do not add a user-facing `fork_turns` field or
support `"all"` or numeric modes. This task fixes the meaning of every actual
Codex fresh spawn.

## Contract and documentation delta

`skills/first-officer/references/fo-dispatch-core.md`, after `## Reuse and Fresh
Dispatch`:

```diff
+**Freshness invariant.** A fresh dispatch creates a new worker handle with no
+inherited parent turns. Runtime adapters enforce that boundary with the host's
+spawn mechanism. This invariant does not change stage selection: a reusable
+worker may still advance through `«addressable-worker»` when every reuse
+condition passes and the next stage does not declare `fresh: true`.
```

`skills/first-officer/references/codex-first-officer-runtime.md`, in
`«worker.spawn»` runtime implementation:

```diff
-call `spawn_agent(task_name,message,fork_turns)` with the helper-emitted prompt
-unchanged as `message`.
+call `spawn_agent(task_name,message,fork_turns="none")` with the helper-emitted
+prompt unchanged as `message`. Every spawn is a fresh dispatch; deliberate
+continuity uses `followup_task` with the existing handle.
```

Keep the live tool-surface probe's generic signature because it describes the
discovered tool shape, not invocation policy. Update its structural contractlint
expectation only as needed to accommodate the exact runtime binding.

No public site wording changes. `docs/site/concepts/stage-lifecycle.md` already
promises that `fresh: true` arrives without the implementer's reasoning, and
`docs/site/concepts/operating-model.md` already promises fresh context with no
access to that reasoning. The implementation makes those existing statements
true. It must not expose `fork_turns` as a user setting.

## Test plan

### Adapter and contract tests — small

- Assert exact three-key map equality including `"fork_turns": "none"` for
  plain, future-override, absent-fork, `"all"`, and numeric helper envelopes.
- Remove the adapter's mutable `ForkTurns` member and retain the reflection guard
  that fails if an override channel returns. A conditional-omission mutation
  must fail the same compact table.
- Update contractlint only for structural host-binding placement. Do not count a
  prose substring assertion as isolation proof.

### One-off Codex behavior evidence — retained, not a permanent harness

- Retain `implementation-live-evidence.md`: the actual spawn arguments contain
  exact `"none"`, the child rollout omits the exact parent-only canary, and
  `followup_task` reaches the same task/thread for a second turn that recalls its
  prior marker.
- Keep the raw artifact hashes for auditability. Do not add or restore a bespoke
  fresh-isolation workflow, rollout parser, or live-runner scenario.

### Regression gates — small

- Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.
- Run focused adapter and contractlint tests plus
  `TestEdgeAdvanceDecision`/`TestEdgeAdvancePatchDoesNotRegressNext`; record the
  exact source SHA and clean status.
- Treat the `automate-beta-release` replacement validation as a separate
  follow-up in its owning workflow after this fix lands.

## 0.25.1 release and 0.26 propagation plan

### Pre-merge proof — this task's acceptance boundary

1. Keep one self-contained fix commit based on current `origin/main`.
2. Run focused, full, race, and existing release-guard tests. The exact
   `v0.25.1` versus `0.26.0-pre1` decision must be `skip`, and
   `TestEdgeAdvancePatchDoesNotRegressNext` must keep the simulated `next` tip
   byte-identical.
3. Retain the one-off live canary/reuse evidence without checking a new harness
   into the product patch.

### Post-merge release ceremony — explicitly not implementation acceptance

1. Land the fix on `main`, then propagate only that fix to a branch based on
   `origin/next` with `git cherry-pick -x <main-fix-commit>`. Run the required
   tests there and publish the updated edge plugin through `next-publish`.
2. Cut `v0.25.1` from `main` per `docs/releasing.md`: stamp the release, green
   Runtime Live E2E for the exact release SHA, and create the annotated tag.
3. Record `origin/next` before and after the tag's `edge-advance` job. The tips
   must match because the equality decision skips; the deliberate cherry-pick
   and `next-publish` are the propagation path.

No release-workflow change is planned. The existing strict-greater decision
guard already covers the exact v0.25.1/0.26.0-pre1 no-rewind case.

### Feedback Cycles

- Cycle 1: REJECTED (fresh Codex validation, 2026-07-19) — the adapter invariant and its compact adversarial test pass, and the retained one-off live evidence is legitimate. Rejection comes from the entity contract demanding a large integrated live journey, unrelated `automate-beta-release` replacement work, and already-consumed release state before this small patch can merge. Captain direction is “simple fix, equivalent simple test.” Route to implementation to narrow ACs and the test plan to the shipped adapter boundary plus one-off live evidence; split unrelated replacement work and treat actual 0.25.1/next publication as the post-merge release ceremony guarded by the existing no-rewind checks. Do not restore a bespoke live harness.

## Stage Report: ideation

- DONE: Trace the #414 design, implementation, and validation gap; distinguish proven facts from the unproven historical Codex default.
  `Evidence and the #414 gap` ties design, code, tests, and live validation to durable sources and marks the June omission semantics unproven.
- DONE: Choose the smallest invariant that makes every fresh Codex spawn isolated without changing deliberate followup_task reuse.
  The adapter makes `fork_turns: "none"` unconditional and unrepresentable otherwise; `followup_task` remains the sole continuity route.
- DONE: Produce biting ACs and a release/backport test plan for 0.25.1 that also preserves the fix on the 0.26 prerelease line.
  Six ACs cover live isolation, replacement evidence, explicit next propagation, and no rewind; focused release, full, and race tests pass.

### Summary

Ideation closes the #414 gap with one invariant: every Codex spawn is isolated,
and every deliberate continuation reuses a handle through `followup_task`. The
plan ships the fix on `main` as v0.25.1, cherry-picks it onto `next`, republishes
the edge plugin deliberately, and proves the patch-tag job leaves 0.26 untouched.

## Stage Report: implementation

- DONE: Every fresh Codex spawn emits exactly fork_turns="none", with unsafe helper input unable to change it; focused adapter tests include a mutation that would expose conditional omission.
  Commit `2148c80c` removes `ForkTurns`; table tests prove absent, `"all"`, numeric, and future-override helper fields all yield the same exact three-argument map.
- DONE: Live Codex evidence proves parent-canary absence on a fresh handle while deliberate followup_task reuse preserves the existing handle and context.
  `implementation-live-evidence.md` records exit 0, literal `"fork_turns":"none"`, exact-canary absence, one child thread reused across two triggered turns, and SHA-256 evidence hashes.
- DONE: The implementation is a self-contained patch suitable for main and clean cherry-pick to next, with release guards proving v0.25.1 cannot rewind the 0.26 prerelease line.
  Rebased commit `2148c80c` changes five files; full and race suites pass, release guard tests pass, and the CLI decision for v0.25.1 versus 0.26.0-pre1 prints `skip`.

### Summary

Codex fresh dispatch now has one executable invariant: `ToolArgs` always emits
`fork_turns: "none"`, and no mutable adapter field can override it. Minimal
contract alignment keeps deliberate continuity on `followup_task`; the shipped
patch contains no bespoke live harness.

## Stage Report: validation

- DONE: Adversarially verify the small adapter invariant: every fresh spawn emits exact fork_turns="none", unsafe helper fields cannot override it, and the compact test fails under a realistic conditional-omission mutation.
  Focused tests passed all absent/`"all"`/numeric/override rows; a throwaway mutation restoring mutable `ForkTurns` plus conditional omission failed all four rows and the structural field guard.
- DONE: Validate the retained one-off Codex evidence for both boundaries: parent canary absent on a fresh child, and followup_task preserves the existing handle/context without shipping a bespoke live harness.
  All three recorded SHA-256 hashes match raw artifacts; spawn used exact `"none"`, the child rollout lacked the exact parent canary, and one child thread handled two turns and recalled `CHILD_CONTINUITY_ONEOFF_6D42`.
- FAILED: Cross-check all six acceptance criteria, including automate-beta-release replacement evidence and the v0.25.1-to-next propagation/no-rewind release proof; report missing evidence as FAILED rather than inferring it.
  AC-2 passes, but AC-1/3/4 lack their specified integrated FO/stage evidence and AC-5/6 promised durable release state does not exist.
- FAILED: **AC-1 (VALUE): Fresh Codex dispatch inherits zero parent turns.**
  Adapter output and the direct one-off host probe pass, but no live run records a generated dispatch artifact's bytes, the FO-issued spawn arguments, and its child result as one reproducible observation.
- DONE: **AC-2: Unsafe fork output is unrepresentable in the Codex adapter.**
  Commit `2148c80c` removes the field, emits one exact three-key map, ignores unsafe helper JSON, and the conditional-omission adversarial mutation fails.
- FAILED: **AC-3: Fresh isolation is an invocation rule, not a default assumption.**
  The executable adapter and direct host probe agree, but the instruction-driven First Officer binding has no recorded live invocation; contract text alone is expressly insufficient.
- FAILED: **AC-4: Deliberate continuity keeps the existing worker and context.**
  The one-off proves raw `followup_task` continuity, but not an eligible stage advancement plus feedback re-review and a `fresh: true` transition; the live rejection-flow run abstained because no structured validation spawn/follow-up handle was correlatable.
- FAILED: **AC-5: `automate-beta-release` has independent replacement evidence.**
  Its durable entity ends with validation cycle 2 from the inherited validator; there is no replacement spawn record/report/state commit produced from an isolated generated artifact, though its worktree is clean.
- FAILED: **AC-6: Stable 0.25 and edge 0.26 both retain the fix without an edge rewind.**
  After fetch, neither `origin/main` (`235e7636`) nor `origin/next` (`53beac31`) contains `2148c80c`, remote `v0.25.1` is absent, and both remote adapters remain unsafe; only the existing CLI guard returns `skip` for next's `0.26.0-pre1` manifest.
- FAILED: Run applicable regression gates at exact source SHA `2148c80c`.
  Full and race suites pass, but `TestLiveCodexSharedScenarios` times out at 10m during `filing`; `gofmt -w ./cmd ./internal` also reformats unrelated upstream `internal/release/journeydelta.go` (reverted to preserve the clean task worktree).

### Summary

Validation verdict: **REJECTED**. AC-1/3/4 have material evidence defects at the exact integrated FO/stage observation boundaries, while AC-5/6 are material outcome defects because the replacement validation and promised release/propagation state are absent; no implementation repair was made. The adapter mechanism itself is strong and mutation-resistant, the retained lower-level Codex one-off is legitimate, and there are no deferred risks asserted from missing proof.

## Stage Report: implementation (cycle 2)

- DONE: Narrow the acceptance criteria and test plan to the simple shipped invariant plus retained one-off live evidence; do not add or restore a bespoke live harness.
  Four reduced ACs now cover the exact adapter map, removal of unsafe state, unchanged reuse policy, and pre-merge release guards; the test plan explicitly retains the one-off instead of shipping a harness.
- DONE: Remove or split the unrelated automate-beta-release replacement requirement, and distinguish pre-merge no-rewind proof from the post-merge 0.25.1/next publication ceremony.
  Beta replacement is named as owning-workflow follow-up, while the release plan separates this task's deterministic no-rewind proof from later main/next/tag operations.
- DONE: Keep commit 2148c80c behavior unchanged unless the contract reshaping exposes a real defect; rerun only checks required by the resulting changes and report the reduced evidence boundary clearly.
  Code HEAD remains clean at `2148c80c`; only the entity changed, and `status --read --json` plus entity `diff --check` verified the 312-line reshaped contract before this appended report.

### Summary

Cycle 2 resolves the validation rejection by aligning the entity with the
captain's intended small patch and equivalent small permanent test. No product
code or test harness changed; remote publication and incident-workflow repair
remain explicit follow-up ceremony after merge.
