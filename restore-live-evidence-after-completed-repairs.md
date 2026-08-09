---
title: Restore live evidence whose original repair owners are complete
status: ideation
source: "Live-test-truth close reconciliation, 2026-08-09"
started: 2026-08-09T18:34:33Z
completed:
verdict:
score: 0.9
worktree:
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
---

## Problem

Eight target cells still use target-scoped `TODO` bindings. Four cells name
completed repair work. One cell belongs to deferred task `47g`. The source code
therefore hides current runtime evidence behind a pre-execution skip.

The `liveTODO(...)` helper skips before the fixture builder and journey run. It
cannot report a repair, an expected semantic failure, or an infrastructure
failure. Strict XFAIL support from `ts7gq0mr9s3chx2w4wppd1kt` must define that
execution boundary before this task changes any binding.

## Value

The common registry will show the result for every target cell. A passing cell
will remove its `TODO` or XPASS binding. A stable product failure will remain
visible with its owner and semantic code. An unexecutable cell will retain a
`TODO` with a concrete setup reason.

The value measure is the target-cell state. The baseline is eight `TODO` cells.
The expected result is four passing cells, two stable semantic failures, and
two unexecutable cells. No target cell will remain unclassified.

## Spike result

The riskiest mechanism was the live assertion boundary. A temporary isolated
checkout removed only the target `TODO` rows. The real fixtures, hosts, and
durable assertions then ran at source commit `a929fcb60`.

The Codex `gate-guardrail` cell passed in 138.89 seconds. The Codex
`withdrawn-gate-recovery` cell from deferred task `47g` passed in 149.68
seconds. The Sonnet and Pi `owned-conflict-owner-handoff` cells passed in
312.44 and 135.39 seconds. These probes did not change the product checkout.

The Pi `gate-guardrail` probe failed at the durable hold assertion with
`state commit missing or before the successful gate prepare`. The Pi
`default-headless-gate-stop` probe failed with `gated entity is not held at its
open validation boundary`. These are semantic failures, not skips.

The Codex `default-headless-gate-stop` probe also failed with `implementation
was not dispatched before validation`. This is the known `98a` product route,
not an `xp6` target. The Opus lane skipped because the local benchmark token
was empty and `ANTHROPIC_API_KEY` was absent. The default Pi model stopped at an
OpenRouter credit error. A lower-cost Pi model enabled the semantic probes.

The `/tmp` artifact paths are local probe locations. The result, command, and
failure text above are the durable evidence for this ideation record.

## Target classification

| Journey | Target | Owner | Classification | Evidence |
| --- | --- | --- | --- | --- |
| `gate-guardrail` | Codex | `xp6c9qfe7y4wwp46enc3f85n` | passing | Exit 0; 138.89s; `codex-shared-scenarios/gate-guardrail` |
| `gate-guardrail` | Pi | `xp6c9qfe7y4wwp46enc3f85n` | stable semantic failure | `state commit missing or before the successful gate prepare`; lower-cost Pi model |
| `recorded-gate-lifecycle` | Claude Opus | `xp6c9qfe7y4wwp46enc3f85n` | unexecutable | Runner skipped because local Claude auth was unavailable |
| `default-headless-gate-stop` | Pi | `xp6c9qfe7y4wwp46enc3f85n` | stable semantic failure | `gated entity is not held at its open validation boundary`; lower-cost Pi model |
| `owned-conflict-owner-handoff` | Claude Sonnet | `xp6c9qfe7y4wwp46enc3f85n` | passing | Exit 0; 312.44s; `claude-shared-scenarios/owned-conflict-owner-handoff` |
| `owned-conflict-owner-handoff` | Claude Opus | `xp6c9qfe7y4wwp46enc3f85n` | unexecutable | Runner skipped because local Claude auth was unavailable |
| `owned-conflict-owner-handoff` | Pi | `xp6c9qfe7y4wwp46enc3f85n` | passing | Exit 0; 135.39s; lower-cost Pi model |
| `withdrawn-gate-recovery` | Codex | `47gnqfm1ft6f2hcahz98m2jv` | passing | Exit 0; 149.68s; deferred `47g` cell |

The two Pi semantic failures require product repair routing. Do not invent a
stable failure code in this task. The product owner must expose one code before
the strict-XFAIL binding is added. The Opus rows require an authenticated
pre-release run. No row is unclassified.

## Required order

1. Land or otherwise expose strict XFAIL behavior from
   `ts7gq0mr9s3chx2w4wppd1kt` at the durable assertion boundary.
2. Run each target on the exact candidate and record the result artifact.
3. Remove a binding only after the exact target passes its durable assertion.
4. Bind a stable semantic failure as strict XFAIL with its active owner and
   exact code. Keep the product repair outside this task.
5. Keep an unexecutable target as `TODO` with the setup reason.
6. Stop on an unclassified result. Do not convert it to `TODO` or XFAIL.

## Proposed approach

Use the existing common journey entry points, fixtures, exercises, and
assertions. Change only target binding rows in
`internal/ensigncycle/shared_live_runner_test.go`.

For the four passing rows, remove the target `TODO` or XPASS binding after an
exact-lane rerun. For the two Pi failures, retain ownership and route a product
repair to a separate `test-behavior-completeness` task. For the two Opus rows,
retain `TODO` until the pre-release lane can execute. Include the Codex
`withdrawn-gate-recovery` row under owner `47gnqfm1ft6f2hcahz98m2jv`; its passing
probe does not silently close or rewrite task `47g`.

The simplest alternative is to leave every `TODO` in place. That alternative
cannot show a repaired journey and hides the four passing cells. A new scenario
table or simulator is also unnecessary because the existing source annotations,
fixtures, and durable assertions already bind the required behavior.

## Acceptance criteria

**AC-1 (VALUE) — Four passing target cells have no target-scoped TODO or XPASS binding.**

Confirmed by exact-lane runs for Codex `gate-guardrail`, Codex
`withdrawn-gate-recovery`, Sonnet `owned-conflict-owner-handoff`, and Pi
`owned-conflict-owner-handoff`. Each run exits 0 and the durable assertion
observes the required end state.

**AC-2 (VALUE) — Stable semantic failures remain executable and report strict XFAIL.**

Confirmed by Pi `gate-guardrail` and Pi `default-headless-gate-stop` runs. Each
run executes its fixture and assertion, reports one owner-scoped semantic code,
and does not report a skip. A different code or infrastructure error fails the
lane.

**AC-3 — Unexecutable target cells retain honest TODO bindings.**

Confirmed by the Opus `recorded-gate-lifecycle` and Opus
`owned-conflict-owner-handoff` lanes after an auth-missing run. The result
records the setup reason and does not claim semantic evidence.

**AC-4 (VALUE) — All eight target cells, including deferred task `47g`, have one classification.**

Confirmed by the classification table, the focused artifacts, and the source
reconciliation. The final state records four passing cells, two strict XFAIL
cells, two TODO cells, and zero unclassified cells. The baseline is eight
target-scoped TODO cells.

**AC-5 — This task changes evidence bindings only.**

Confirmed by the implementation diff. It changes no First Officer behavior,
fixture semantics, durable assertion, desired registry, command grammar,
stored format, authority rule, or metrics format.

## Test plan

Run the pure strict-XFAIL grade and metrics tests from `ts7gq0mr9s3chx2w4wppd1kt`
before any binding edit. A matching code must produce XFAIL. A pass must produce
XPASS. A different code, launch failure, timeout, or authentication failure
must remain FAIL or TODO.

Run each focused common journey on the exact candidate. Use local subscription
auth before paid CI. A lower-cost Pi model can diagnose product behavior, but
the default Pi lane must pass before its binding is removed.

Run these deterministic checks after binding edits:

- `go test ./internal/ensigncycle -run 'TestGateGuardrailNegativeBrokenStateTransition|TestAssertGateHeld'`
- `go test ./internal/ensigncycle -run 'TestRecordedGateLifecycle'`
- `go test ./internal/contractlint -run '^TestRuntimeLiveRegistryReconciliation$'`
- `SPACEDOCK_LIVE_STATE_DIR=docs/dev/.spacedock-state go test ./internal/contractlint -run '^TestRuntimeLiveTODOOwnersAreActive$'`

Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.
No new fixture, simulator, runner test, or reconciliation test is needed.

## Expected surface and estimates

The product surface is one existing test file. The estimate has a tolerance of
two changed lines. The state entity and its stage report are workflow records,
not product LOC.

| File | Gross additions | Gross deletions | Net | Purpose |
| --- | ---: | ---: | ---: | --- |
| `internal/ensigncycle/shared_live_runner_test.go` | 4 | 4 | 0 | Rewrite four target binding rows after exact evidence |
| **Product total** | **4** | **4** | **0** | No product repair or new harness |
| `docs/dev/.spacedock-state/restore-live-evidence-after-completed-repairs.md` | ~180 | ~31 | ~+149 | Durable plan, classification, and stage report |

Product repair lines are outside this capstone. Route the Codex and Sonnet
default-headless failure to `98aa776adg66gn823a8gamdq`. Route the two Pi
semantic failures to new product-repair entities under
`test-behavior-completeness`. Route Opus authentication to the pre-release
lane owner. Leave the passing `47g` evidence for that entity's own disposition.

## Observable semantic boundary

- Command grammar: unchanged.
- Stored formats: unchanged.
- Authority: unchanged.
- Runtime behavior: only target evidence bindings change after a real run.
- Metrics: use the existing journey result and artifact.
- Documentation: no user-visible behavior changes require a doc diff.

## Scope

- Reuse the current common journey entry points, fixtures, exercises, and assertions.
- Use local subscription-backed live runs before paid CI where possible.
- Record the exact result for every target cell in the matrix.
- Keep product repairs and new failure codes outside this task.

## Out of scope

- Weakening a durable assertion.
- Adding a simulator or a new scenario table.
- Adding tests for the reconciliation code.
- Changing the desired journey registry.
- Changing First Officer behavior, gate authority, or stored formats.

## Stage Report: ideation

- DONE: Keep this task evidence-only and include the withdrawn-gate Codex cell from deferred task 47g.
  The plan changes only target binding rows and includes Codex `withdrawn-gate-recovery` under owner `47gnqfm1ft6f2hcahz98m2jv`; its isolated run passed at source commit `a929fcb60` in 149.68 seconds.
- DONE: Classify each target as passing, stable semantic failure, unexecutable, or unclassified.
  The matrix records four passing cells, two stable semantic failures, two unexecutable Opus cells, and zero unclassified cells; Codex and Pi focused runs supplied the result text and exit codes.
- DONE: Give gross and net line estimates with product-repair routing outside this capstone.
  The product estimate is 4 gross insertions, 4 gross deletions, and 0 net lines in one existing test file; Pi failures route to separate product entities, while `98a` owns the adjacent Codex and Sonnet headless failure.

### Summary

The ideation package defines an evidence-only binding change. It records the
riskiest live assertion probe first, including the deferred `47g` Codex cell.
It keeps stable failures and unavailable Opus lanes visible, and routes product
repairs outside this capstone.
