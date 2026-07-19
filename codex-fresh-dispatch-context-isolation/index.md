---
id: rt8eywbyf3d7nyc9bsmrymnq
title: Codex fresh dispatch must isolate parent turns
status: validation
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
4. Add adapter and live behavior tests. Unit tests close every adapter escape;
   a Codex live scenario proves the instruction-driven First Officer sends the
   safe argument and the fresh child lacks a parent-only canary.

No spike is needed. The incident proves the unsafe explicit-`all` outcome, and
current upstream Codex source proves that `"none"` selects no fork. The live
scenario remains required because Spacedock's First Officer is instruction-
driven; adapter tests alone cannot prove the model-issued call.

## Why this is the smallest sufficient mechanism

- Setting the struct field's default to `"none"` while retaining a mutable
  `ForkTurns` field leaves an internal path back to `"all"` or a number. Removing
  the field makes unsafe output unrepresentable and serves AC-1 and AC-2.
- Adding a workflow or user-facing fork setting creates a choice Spacedock does
  not need. `spawn_agent` and `followup_task` already express isolation and
  continuity. The extra setting would weaken AC-1.
- Marking every implementation stage `fresh: true` changes workflow reuse
  policy, not spawn semantics. It would increase cost and still leave any other
  fresh spawn vulnerable. AC-4 preserves the existing stage policy.
- Contract wording alone cannot prove behavior. The adapter tests and live
  canary serve AC-1; the prose makes the instruction-driven call match the
  executable invariant.
- Rerunning only the affected validation repairs one report but leaves every
  later fresh dispatch exposed. AC-5 requires the rerun after AC-1 ships.

The `per-host-stage-model-override` task is compatible, not duplicative. Its
current design forces `"none"` only when model or effort overrides need it and
records plain-spawn isolation as later work. When rebased, that task should
forward overrides onto this invariant and must not restore a mutable fork mode.

## Acceptance criteria

- **AC-1 (VALUE): Fresh Codex dispatch inherits zero parent turns.** The actual
  `spawn_agent` call contains exactly
  `fork_turns: "none"`; the generated dispatch artifact omits a unique
  parent-only canary, and the child reports no canary in its initial context.
  Verified by unit tests for plain and override-bearing build envelopes plus a
  Codex live scenario that records the call arguments, dispatch-file bytes, and
  child result. A new handle alone fails this AC.
- **AC-2: Unsafe fork output is unrepresentable in the Codex adapter.**
  `CodexMultiAgentV2Spawn` exposes no fork-mode field, and every `ToolArgs` map
  contains exactly `"fork_turns": "none"`. Verified by table-driven tests that
  feed absent, `"all"`, and numeric `fork_turns` fields in helper JSON and assert
  the same safe output; an adversarial test mutation that restores conditional
  omission must fail.
- **AC-3: Fresh isolation is an invocation rule, not a default assumption.** The
  host-neutral core defines
  fresh as a new handle with no inherited parent turns, and the Codex binding
  supplies `"none"` on every spawn. Verified by AC-1's live behavior. Reading or
  grepping the contract does not satisfy this AC.
- **AC-4: Deliberate continuity keeps the existing worker and context.** An
  eligible non-fresh stage advancement and a feedback re-review use
  `followup_task` with the captured handle; a `fresh: true` transition uses a
  different handle with zero inherited parent turns. Verified by one live or
  durable structured scenario containing both paths. No implementation stage
  becomes globally fresh.
- **AC-5: `automate-beta-release` has independent replacement evidence.** A new
  validator is spawned with `"none"` from only its
  generated dispatch artifact. Its durable report records the replacement run
  and retains any valid original findings instead of discarding them solely
  because the first validator inherited context. Verified by the new worker's
  spawn record, report, state-checkout commit, and clean worktree status.
- **AC-6: Stable 0.25 and edge 0.26 both retain the fix without an edge rewind.**
  The fix lands on `main`, ships in annotated `v0.25.1`, and is integrated onto
  `next` while `next` keeps its 0.26 manifest, gate line, calendar lineage, and
  branch-exclusive commits. Verified by tests at both branch tips, exact patch
  equivalence, `edge-advance-decision v0.25.1` returning `skip` against
  `0.26.0-pre1`, and the remote `next` tip remaining unchanged by the patch-tag
  release job.

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

- Replace the omission assertion in `codex_v2_adapter_test.go` with exact map
  equality including `"fork_turns": "none"`.
- Add table rows for plain envelopes, future model/effort override envelopes,
  and helper JSON containing absent, `"all"`, or numeric fork fields. Every row
  yields the same isolation value. The test should remain compatible whether
  the model-override task lands before or after this fix.
- Remove the adapter's mutable `ForkTurns` member. Compile failures in callers
  expose any hidden override path.
- Update contractlint only for structural host-binding placement. Do not count a
  prose substring assertion as isolation proof.

### Codex live behavior — medium, required

- Add a focused fresh-dispatch scenario to the existing Codex live harness. Put
  a unique canary in the root conversation, keep it out of the generated
  dispatch file, advance to a `fresh: true` stage, and capture the actual
  `spawn_agent` arguments.
- Assert a new handle, exact `fork_turns: "none"`, canary absence from the
  dispatch-file bytes, and canary absence from the child's initial context
  report. Include an unsafe fixture or mutation using `"all"` to prove the
  canary detector turns red.
- In the same or a focused companion scenario, advance one eligible non-fresh
  stage with `followup_task` and prove the handle stays constant. Reuse evidence
  prevents a global-fresh implementation from false-passing.
- Rerun `TestLiveCodexSharedScenarios` after the focused scenario so existing
  rejection-flow reuse stays green.

### Regression gates — medium

- Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.
- Run focused adapter, contractlint, live fresh-isolation, and reviewer-reuse
  tests. Record exit status, exact source SHA, durable state, and clean status.
- Run the affected `automate-beta-release` validation again from its owning
  repository with a new isolated worker. Preserve both the original finding and
  the replacement evidence in its durable report.

## 0.25.1 release and 0.26 propagation plan

1. Land one self-contained fix commit on `main`. Run the deterministic and Codex
   live gates on that commit.
2. Propagate only the fix to a branch based on `origin/next` with
   `git cherry-pick -x <main-fix-commit>`. This avoids copying the 0.25.1 manifest
   stamp onto the 0.26 line. Run the same focused and full deterministic tests at
   the resulting `next` SHA, then push by fast-forward.
3. Run the deliberate `next-publish` path after the 0.26 fix lands so edge plugin
   users receive the corrected contract. Record the resulting `next` SHA and
   confirm its manifest remains `0.26.0-pre1` or later.
4. Cut `v0.25.1` from `main` by `docs/releasing.md`: stamp `0.25.1`, push the
   release commit to `main`, green Runtime Live E2E for that exact SHA, and create
   an annotated tag on that SHA.
5. Before pushing the tag, run the existing release guards, especially
   `TestEdgeAdvanceDecision`'s exact `v0.25.1` versus `0.26.0-pre1` equality case
   and `TestEdgeAdvancePatchDoesNotRegressNext`. The decision must be `skip`.
6. Record `origin/next` before the tag push and after the `edge-advance` job. The
   SHAs must match. The patch job must not merge the 0.25.1 tree, restamp the
   0.26 manifests, bump its calendar key, or create a new pre0 tag. The earlier
   explicit cherry-pick and `next-publish` are the propagation path.

Estimated complexity: small code change, medium live and release verification.
No release-workflow change is planned because the existing strict-greater
decision guard already contains the exact 0.25.1/0.26.0-pre1 regression case.

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
