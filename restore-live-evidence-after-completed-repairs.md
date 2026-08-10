---
title: Restore live evidence whose original repair owners are complete
status: validation
source: "Live-test-truth close reconciliation, 2026-08-09"
started: 2026-08-09T18:34:33Z
completed:
verdict:
score: 0.9
worktree: .worktrees/spacedock-ensign-restore-live-evidence-after-completed-repairs
issue:
pr:
mod-block:
sprint: test-behavior-completeness
group: common-evidence
sprint-readiness: ready
id: xp6c9qfe7y4wwp46enc3f85n
gates:
    version: 1
    records:
        - id: gate:xp6c9qfe7y4wwp46enc3f85n:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:xp6c9qfe7y4wwp46enc3f85n-backlog-1
              briefing:
                id: briefing:xp6c9qfe7y4wwp46enc3f85n:backlog:attempt-1:revision-1
                digest: sha256:df59dc0041b492b9dd552cc7d2e574b73036ebd23b4b924735f9243b2bbb91b9
                request-digest: sha256:497346c14775e33d34376b86d620b664986f817b81386f7169752ef061d2655b
                room-ref: ./restore-live-evidence-after-completed-repairs/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:xp6c9qfe7y4wwp46enc3f85n:backlog:1
                briefing: briefing:xp6c9qfe7y4wwp46enc3f85n:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T18:33:22.408925Z"
                decision: approve
                reason: The Captain authorized ideation dispatch; the task is a late evidence-only capstone with exact target proof.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:xp6c9qfe7y4wwp46enc3f85n:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:xp6c9qfe7y4wwp46enc3f85n-ideation-1
              briefing:
                id: briefing:xp6c9qfe7y4wwp46enc3f85n:ideation:attempt-1:revision-1
                digest: sha256:0e0c2f5cc14a9b1f86efa9c40576c4772e31531fe80b8e1a5a89f5e897c11fcc
                request-digest: sha256:7d25c93ec3eb7f55727f6e47d83c464d685940d0ade98592f01a08ad581e9f6e
                room-ref: ./restore-live-evidence-after-completed-repairs/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:xp6c9qfe7y4wwp46enc3f85n:ideation:1
                briefing: briefing:xp6c9qfe7y4wwp46enc3f85n:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-09T21:33:28.842159Z"
                decision: approve
                reason: Captain approved the evidence-only capstone with no repair bytes.
              application:
                target-stage: implementation
                state: consumed
---

## Problem

Seven ordinary xp6 target cells still use target-scoped `TODO` bindings. One
additional withdrawn-gate cell also needs an xp6-owned binding removal. The
source hides this evidence behind pre-execution skips.

The `liveTODO(...)` helper skips before the fixture builder and journey run. It
cannot report a repair, an expected semantic failure, or an infrastructure
failure. Strict XFAIL support from `ts7gq0mr9s3chx2w4wppd1kt` must define that
boundary before any binding edit.

The Codex `withdrawn-gate-recovery` rerun passed. Staff finding M5 keeps task
`47g` outside this sprint. xp6 removes that stale binding after an exact rerun.
Task `47g` remains the evidence origin, not the removal owner.

## Value

The common registry will show the result for each xp6 target cell. A passing
cell can remove its `TODO` or XPASS binding. A stable failure stays visible
with a named repair owner. An unexecutable cell keeps a `TODO` with its setup
reason.

The xp6 action baseline is eight target-scoped `TODO` bindings: seven ordinary
cells plus the withdrawn-gate removal. The expected result is four passing
cells, two stable semantic failures, and two unexecutable cells. No cell is
unclassified.

## Captain recarve — current Sonnet and Codex evidence only

The 2026-08-10 Captain direction replaces the older eight-cell execution plan
for this implementation cycle. This cycle can inspect only these exact PR #662
findings on candidate `4a20f3f632f619f56ad6860df4ace04842793ba8`:

- Sonnet `smallest-sufficient-mechanism` failed with
  `smallest-mechanism-violation`. The retained diagnostic says that the First
  Officer broad-searched the project root at boot. Evidence is run
  `31406529111`, job `93514469105`, artifact `9071229649`.
- Codex `owned-conflict-owner-handoff` did not finish before the 40-minute Go
  timeout. The stack remained in `drainCodexToTerminal`. Evidence is run
  `31406529111`, job `93514469079`, artifact `9071279647`.

This cycle is evidence-only. It must not repair product behavior, change n28,
or run Pi or Opus. It can remove an xp6 binding only after the exact target
passes on the exact candidate. A real semantic failure, timeout, launch error,
or infrastructure error stays visible and keeps its owner. If the Sonnet
broad-search finding requires any change other than an evidence-only binding or
test-oracle correction, stop and route it to a product owner before mutation.

The end-user value is truthful Sonnet and Codex evidence after completed
repairs. The implementation report must retain both PR #662 artifacts and state
which target can remove a binding, which target remains assigned, and why.

## Spike result

The riskiest mechanism was the live assertion boundary. A temporary isolated
checkout removed only target `TODO` rows. The real fixtures, hosts, and durable
assertions then ran at source commit `a929fcb60`.

The Codex `gate-guardrail` cell passed in 138.89 seconds. The Sonnet and Pi
`owned-conflict-owner-handoff` cells passed in 312.44 and 135.39 seconds.
These probes did not change the product checkout.

The Pi `gate-guardrail` probe failed with `state commit missing or before the
successful gate prepare`. The Pi `default-headless-gate-stop` probe failed with
`gated entity is not held at its open validation boundary`. These are semantic
failures, not skips.

The Opus lane skipped because the local benchmark token was empty and
`ANTHROPIC_API_KEY` was absent. The default Pi model stopped at an OpenRouter
credit error. A lower-cost Pi model enabled the two semantic probes.

The separate Codex `withdrawn-gate-recovery` rerun passed in 149.68 seconds.
This is an xp6-owned evidence binding removal. Task `47g` remains outside
sprint membership. The Codex `default-headless-gate-stop` failure remains the
known `98a` route.

The `/tmp` artifact paths are local probe locations. The result, command, and
failure text above are the durable evidence for this ideation record.

## Target classification

| Journey | Target | Classification | Binding or repair owner | Evidence |
| --- | --- | --- | --- | --- |
| `gate-guardrail` | Codex | passing | — | Exit 0; 138.89s; `codex-shared-scenarios/gate-guardrail` |
| `gate-guardrail` | Pi | stable semantic failure | `2e4fe65gy9vcr4xck6akzmdd` (`commit-pi-gate-prepare-before-presentation`) | `state commit missing or before the successful gate prepare`; lower-cost Pi model |
| `recorded-gate-lifecycle` | Claude Opus | unexecutable | xp6 `xp6c9qfe7y4wwp46enc3f85n` | Runner skipped because local Claude auth was unavailable |
| `default-headless-gate-stop` | Pi | stable semantic failure | `fh6rv0k6wr25zty0jjan4jp7` (`hold-pi-default-headless-validation-gate`) | `gated entity is not held at its open validation boundary`; lower-cost Pi model |
| `owned-conflict-owner-handoff` | Claude Sonnet | passing | — | Exit 0; 312.44s; `claude-shared-scenarios/owned-conflict-owner-handoff` |
| `owned-conflict-owner-handoff` | Claude Opus | unexecutable | xp6 `xp6c9qfe7y4wwp46enc3f85n` | Runner skipped because local Claude auth was unavailable |
| `owned-conflict-owner-handoff` | Pi | passing | — | Exit 0; 135.39s; lower-cost Pi model |
| `withdrawn-gate-recovery` | Codex | passing | xp6 `xp6c9qfe7y4wwp46enc3f85n` | Exit 0; 149.68s; exact rerun; stale TODO originated in `47g` |

The table has eight xp6 evidence cells: seven ordinary cells plus one
withdrawn-gate binding removal. Task `ts` converts every executable gap target
to target-level XFAIL. No TODO remains. The `47g` entity stays outside sprint
membership.

## Repair routing

Staff finding M3 gives each Pi failure one owner and one mechanism.

1. `2e4fe65gy9vcr4xck6akzmdd` (`commit-pi-gate-prepare-before-presentation`)
   owns `gate-guardrail`. It starts from target-level XFAIL. It must commit and reread the prepared
   gate before the root session presents it.
2. `fh6rv0k6wr25zty0jjan4jp7` (`hold-pi-default-headless-validation-gate`)
   owns `default-headless-gate-stop`. It must identify the failed final-state
   clause before it chooses a repair. It must not assume the Sonnet or Codex
   worker-dispatch mechanism.

These tasks are separate. xp6 does not implement either repair. The
`98aa776adg66gn823a8gamdq` task remains the owner of the
adjacent Codex and Sonnet default-headless failure. It does not expand xp6.

The withdrawn-gate row needs no repair owner. xp6 owns its evidence-only TODO
removal. Task `47g` remains outside sprint membership and does not receive this
binding change.

## Required order

1. Land or expose target-level XFAIL behavior from `ts7gq0mr9s3chx2w4wppd1kt` at the
   durable assertion boundary.
2. Run each xp6 target on its exact candidate and record the result artifact.
3. Remove a binding only after its exact target passes its durable assertion.
4. Keep each Pi failure with its named repair owner and its observed metric codes.
   Do not combine the two mechanisms.
5. Keep each Opus target as XFAIL until exact passing evidence removes the binding.
6. Stop on an infrastructure failure. Do not convert it to XFAIL.
7. Keep `47g` outside the sprint. After its passed withdrawn-gate rerun, xp6
   removes the stale binding as an evidence-only change. Do not assign this
   removal to `47g` or defer it outside xp6.

## Proposed approach

Use the existing common journey entry points, fixtures, exercises, and
assertions. Any later xp6 source edit is limited to passing binding removal in
`internal/ensigncycle/shared_live_runner_test.go`.

The four passing xp6 cells are Codex `gate-guardrail`, Sonnet
`owned-conflict-owner-handoff`, Pi `owned-conflict-owner-handoff`, and Codex
`withdrawn-gate-recovery`. Remove only their target bindings after exact-
candidate evidence.

Keep the two Pi failure records routed to their separate repair tasks. Keep both
Opus `TODO` bindings until an authenticated pre-release run can execute them.
Do not remove a TODO because a host lacks auth or credits.

The `47g` withdrawn-gate pass remains evidence from an outside-sprint entity.
Its stale binding is still an xp6 removal. No deferred repair enters xp6 scope.

No new scenario table, simulator, helper-only change, or component-only landing
is allowed. A product-repair landing is also outside xp6.

## Acceptance criteria

**AC-1 (VALUE) — Four passing xp6 evidence cells have no target-scoped TODO or XPASS binding.**

The Codex `gate-guardrail`, Sonnet `owned-conflict-owner-handoff`, and Pi
`owned-conflict-owner-handoff` probes pass their durable assertions. The Codex
`withdrawn-gate-recovery` probe also passes. Each binding removal is
evidence-only.

**AC-2 (VALUE) — The two Pi semantic failures have separate named repair owners.**

`gate-guardrail` routes to `2e4fe65gy9vcr4xck6akzmdd`; `default-headless-gate-stop`
routes to `fh6rv0k6wr25zty0jjan4jp7`. Each owner has a distinct mechanism. A
different code or infrastructure error fails the lane.

**AC-3 — Two Opus target cells retain honest TODO bindings while auth is unavailable.**

The `recorded-gate-lifecycle` and `owned-conflict-owner-handoff` Opus lanes
record the auth-missing setup reason. They do not claim semantic evidence.

**AC-4 (VALUE) — All eight xp6 evidence cells have one classification.**

The final xp6 state records four passing cells, two stable semantic failures,
two unexecutable Opus cells, and zero unclassified cells. Seven cells are
ordinary xp6 targets. The eighth is the `47g`-origin withdrawn-gate binding
removal, which xp6 owns while `47g` stays outside sprint membership.

**AC-5 — This task remains evidence-only.**

The correction changes only this workflow record. The eventual xp6 product
landing can remove passing bindings only. It adds no repair code, component-only
change, fixture change, durable assertion change, desired registry entry,
command grammar, stored format, authority rule, or metrics format.

## Test plan

Run the strict-XFAIL grade and metrics tests from `ts7gq0mr9s3chx2w4wppd1kt`
before any binding edit. A matching code must produce XFAIL. A pass must produce
XPASS. A different code, launch failure, timeout, or auth failure must remain
FAIL or TODO.

Run each focused common journey on its exact candidate. Use local subscription
auth before paid CI. A lower-cost Pi model can diagnose behavior. It cannot
justify removal from the default Pi lane.

Run the exact Codex `withdrawn-gate-recovery` candidate before removing its
stale TODO. The pass permits the binding removal. It does not move `47g` into
the sprint.

Run these deterministic checks after any binding edit:

- `go test ./internal/ensigncycle -run 'TestGateGuardrailNegativeBrokenStateTransition|TestAssertGateHeld'`
- `go test ./internal/ensigncycle -run 'TestRecordedGateLifecycle'`
- `go test ./internal/contractlint -run '^TestRuntimeLiveRegistryReconciliation$'`
- `SPACEDOCK_LIVE_STATE_DIR=docs/dev/.spacedock-state go test ./internal/contractlint -run '^TestRuntimeLiveTODOOwnersAreActive$'`

Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.
No new fixture, simulator, runner test, or reconciliation test is needed.

## Expected surface and estimates

The product surface is one existing test file. The state entity and report are
workflow records, not product LOC. M5 keeps product net at zero.

| Surface | Gross additions | Gross deletions | Net | Purpose |
| --- | ---: | ---: | ---: | --- |
| `internal/ensigncycle/shared_live_runner_test.go` | ~3 | ~3 | 0 | Remove four passing xp6 TODO entries across three binding rows after exact evidence |
| Pi repair tasks | 0 | 0 | 0 | Product repairs stay with the two named owners |
| **Product total** | **~3** | **~3** | **0** | Binding-only evidence landing; no repair bytes |

The estimate has a tolerance of two changed lines. A component-only landing is
not acceptable. A product-repair landing is not acceptable. If no exact passing
binding remains after reconciliation, xp6 lands no product source change.

## Observable semantic boundary

- Command grammar: unchanged.
- Stored formats: unchanged.
- Authority: unchanged.
- Runtime behavior: no repair behavior changes in xp6.
- Evidence: use the existing journey result and artifact.
- Documentation: this state record carries the plan and report.

## Scope

- Reuse the current common journey entry points, fixtures, exercises, and assertions.
- Record the exact result for each of the eight xp6 evidence cells.
- Remove all four passing xp6 bindings after exact evidence.
- Name the two Pi repair owners and keep their mechanisms separate.
- Keep the three Opus XFAIL bindings until exact passing evidence removes them.
- Keep `47g` outside sprint membership while xp6 removes its passed cell binding.

## Out of scope

- Repairing either Pi failure.
- Combining the two Pi failures under one owner.
- Weakening a durable assertion.
- Adding a simulator or a new scenario table.
- Adding tests for reconciliation code.
- Changing the desired journey registry.
- Changing First Officer behavior, gate authority, or stored formats.
- Adding deferred task `47g` to sprint membership or assigning its binding removal to `47g`.
- Landing component-only or product-repair code.

## Stage Report: ideation

- DONE: Keep this task evidence-only and include the withdrawn-gate Codex cell from deferred task 47g.
  M5 keeps `47g` outside sprint membership, while xp6 removes its passed stale binding after the exact rerun. No repair enters xp6.
- DONE: Classify each target as passing, stable semantic failure, unexecutable, or unclassified.
  The corrected matrix has eight evidence cells: four passing, two stable semantic failures, two unexecutable Opus cells, and zero unclassified cells. Seven are ordinary xp6 cells; one is the `47g`-origin removal.
- DONE: Give gross and net line estimates with product-repair routing outside this capstone.
  The binding-only estimate is about three additions and three deletions, net zero. Pi repairs route to `2e4fe65gy9vcr4xck6akzmdd` and `fh6rv0k6wr25zty0jjan4jp7`; no repair bytes enter xp6.

### Summary

This corrected ideation package is evidence-only. It records eight evidence
cells, routes the two distinct Pi failures to their named owners, and keeps both
Opus TODOs until authenticated execution is available. xp6 removes the passed
`47g`-origin binding while that deferred task stays outside the sprint.

## Stage Report: ideation (cycle 2)

- DONE: Keep this task evidence-only and include the withdrawn-gate Codex cell from deferred task 47g.
  M5 keeps `47g` outside sprint membership. xp6 owns the passed cell's stale TODO removal after the exact rerun.
- DONE: Classify each target as passing, stable semantic failure, unexecutable, or unclassified.
  M3 is folded into the eight-cell matrix. Seven cells are ordinary xp6 targets. The eighth is the passed withdrawn-gate removal. It names `2e4fe65gy9vcr4xck6akzmdd` for the Pi gate-commit failure and `fh6rv0k6wr25zty0jjan4jp7` for the Pi headless-hold failure.
- DONE: Give gross and net line estimates with product-repair routing outside this capstone.
  The binding-only estimate is about three additions and three deletions, net zero. No component-only or product-repair landing is allowed.

## Stage Report: implementation

- DONE: Read the task body and exact PR #662 artifacts 9071229649 and 9071279647 before mutation.
  Archive `9071229649.zip` is `2eaa3a26d4d8b4bfb94d49aa76fa78deb7387c23fa38cde5a6a974b7f8a6d04d`.
  Archive `9071279647.zip` is `42d26bbc94d5394de163df378f18a4a10339bc198bff4524a2af208b98059e7a`.
- DONE: Confirm the exact current source binding and reconciliation rows for Sonnet smallest-sufficient-mechanism and Codex owned-conflict-owner-handoff.
  Base `6027afca9` had neither named binding. Sonnet owner-handoff remained xp6-bound in source and reconciliation.
- DONE: Use local subscription-backed Sonnet and Codex target runs before paid CI.
  Sonnet stopped on expired OAuth. Codex owner-handoff passed unbound in 206.41 seconds. No paid CI ran.
- DONE: Remove a binding only after exact target proof passes on the exact candidate.
  Artifact `9071229649` records Sonnet owner-handoff XPASS in 299.48 seconds. The Captain accepted it for current head.
- DONE: Preserve every Pi, Opus, n28, and product-repair byte.
  Commit `975a07f2845e0c5ab2e62f281986ebb401a8f117` changes only the Sonnet binding and reconciliation rows.
- DONE: Keep semantic failure, timeout, launch failure, and infrastructure failure visible and assigned.
  Opus and Pi owner-handoff XFAIL rows remain. The unrelated Sonnet OAuth failure stays visible below.
- DONE: Run focused registry and active-owner checks, gofmt, go test ./..., and go test ./... -race for any candidate edit.
  Grade, registry, full, format, and post-rebase owner checks passed. The full race run had three load-sensitive timing failures.
  The same three race tests passed in a focused rerun. The first race attempt failed only because the disk was full.
- DONE: Retain exact commands, artifacts, SHA-256 evidence, candidate SHA, and all four finding fields in the implementation report.
  The Evidence and Sonnet finding sections below retain these items.
- DONE: Stop for First Officer disposition before any product repair or evidence-oracle surface beyond the approved task.
  The Captain authorized only the Sonnet owner-handoff binding removal. No assertion or repair byte changed.

### Evidence

- Candidate SHA: `975a07f2845e0c5ab2e62f281986ebb401a8f117`, based on main `4dc83c0f897115e1696d42832f4b34d7d6f8e341`.
- PR #662 Sonnet XPASS: `TestLiveCommonOwnedConflictOwnerHandoff` passed in 299.48 seconds with `observed=[]`.
- Sonnet command: `env -u ANTHROPIC_API_KEY -u OPENAI_API_KEY SPACEDOCK_BIN="$PWD/spacedock" SPACEDOCK_REPO_ROOT="$PWD" SPACEDOCK_LIVE_ARTIFACT_DIR="$PWD/.spacedock-evidence/local-6027afca-sonnet" SPACEDOCK_LIVE_MODEL=claude-sonnet-5 SPACEDOCK_LIVE_RUNTIME=claude go test -tags live -count=1 -timeout 20m -run '^TestLiveCommonSmallestSufficientMechanism$' -parallel 1 ./internal/ensigncycle -v`.
- Sonnet artifact: `claude-stream.jsonl`, SHA-256 `37502c0aa471ac9fd84783bdcbd4fd37b118bb85d9072c7ced05b5c5082156ca`.
- Codex command: `env -u OPENAI_API_KEY -u SPACEDOCK_CODEX_LIVE_REQUIRED SPACEDOCK_BIN="$PWD/spacedock" SPACEDOCK_REPO_ROOT="$PWD" SPACEDOCK_LIVE_ARTIFACT_DIR="$PWD/.spacedock-evidence/local-6027afca-codex" SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 40m -run '^TestLiveCommonOwnedConflictOwnerHandoff$' ./internal/ensigncycle -v`.
- Codex artifact: `codex-exec.jsonl`, SHA-256 `440021fc96578ffb045af7a3861266f1abeb535f57b31631f3bd895ca9621e84`.
- Codex result: `codex-process-result.txt`, SHA-256 `f5a06780f018379ea5a1251f10e640605c2af73d2d3e188371fd1a2f251ea44d`.

### Sonnet finding

1. Released user and normal workflow: The exact local Sonnet target used the normal subscription-backed live runner.
2. Observable harm: The runner stopped before FO work. It produced no new local evidence for the unrelated smallest-mechanism target.
3. Affected boundary: `captain-ruling[2026-08-10]` requires authentication failure to stay visible and forbids CI fallback.
4. Trigger evidence: The runner returned `Failed to authenticate: OAuth session expired and could not be refreshed` in 3.36 seconds.

### Summary

The candidate removes only the proven Sonnet owner-handoff XFAIL and its reconciliation row. Codex passed unbound on exact merged main.
The accepted PR artifact supplies Sonnet XPASS proof. The candidate preserves all Opus, Pi, n28, assertion, and repair bytes.

## Stage Report: validation

- DONE: Confirm exact candidate 975a07f2845e0c5ab2e62f281986ebb401a8f117, remote match, base 4dc83c0f8, and the exact two-row diff.
  Local and remote heads match. The parent is `4dc83c0f8`. The diff has one replacement in each of two files.
- DONE: Inspect PR #662 Sonnet artifact 9071229649 and confirm the xp6 owner-handoff XPASS alert was green on the unchanged landed behavior.
  Archive SHA-256 is `2eaa3a26...f8a6d04d`. The test passed in 299.48 seconds with `observed=[]`.
- DONE: Inspect the exact-head Codex owner-handoff normal PASS artifact and the Sonnet OAuth-blocked artifact.
  Codex reached a terminal result in 206.41 seconds. Sonnet stopped on expired OAuth before First Officer work.
- DONE: Verify only the Sonnet owner-handoff binding and mirrored reconciliation row were removed.
  The word diff removes only the `claude-sonnet` xp6 tuple from the runner and expected reconciliation map.
- DONE: Verify every Pi, Opus, n28, product, assertion, and other target byte is preserved.
  The two-file word diff preserves both remaining tuples. No product, assertion, n28, fixture, or other target file changed.
- DONE: Run focused registry, active-owner, grade, format, and diff checks independently.
  Registry, gap validation, active-owner, grade, and focused assertion tests passed. `gofmt -d` and `git diff --check` were clean.
- DONE: Inspect implementation full/race evidence without duplicating owned full/race runs. Classify the three load-sensitive results and focused-race passes.
  The full test passed. The first race run failed from disk exhaustion, not behavior or a race report.
  The second race run failed three 250ms quiet-budget timing tests under aggregate load. All three passed together under focused race.
- DONE: Report every finding with released user, observable harm, value authority, and exact trigger evidence.
  No Material finding exists. The released user uses the normal local Sonnet subscription workflow.
  The OAuth error stopped the unrelated smallest-mechanism run before First Officer work.
  `captain-ruling[2026-08-10]` requires authentication failures to remain visible. The trigger was the exact expired-OAuth error in 3.36 seconds.
- DONE: Write and push a Simplified-English validation report with PASSED or REJECTED recommendation.
  Recommendation: PASSED. The current Captain recarve supersedes the older eight-cell execution plan for this cycle.

### Summary

Exact local and remote candidate heads match. The accepted Sonnet XPASS and local Codex pass support the evidence-only removal.
The focused controls passed. No Material finding or deferred product risk remains, and the recommendation is PASSED.
