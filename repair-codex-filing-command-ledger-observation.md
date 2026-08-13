---
title: Repair Codex filing command-ledger observation
status: implementation
source: PR #679 run 31728107636, Codex job 94541783359
sprint: test-behavior-completeness
sprint-readiness: ready
score: 0.9
id: 6ker7h25hj86983e5ef71ahm
gates:
    version: 1
    records:
        - id: gate:6ker7h25hj86983e5ef71ahm:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:6ker7h25hj86983e5ef71ahm-backlog-1
              briefing:
                id: briefing:6ker7h25hj86983e5ef71ahm:backlog:attempt-1:revision-1
                digest: sha256:97242c2367a7b4e96084790b7765170b01ec932b6a9873d73c5233c694dc79cc
                request-digest: sha256:96f4def00a71e0db7e1e41b3c0d0e4405e2deed271501dcdb2b4ad62429f82b3
                room-ref: ./repair-codex-filing-command-ledger-observation/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:6ker7h25hj86983e5ef71ahm:backlog:1
                briefing: briefing:6ker7h25hj86983e5ef71ahm:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-13T22:09:09.620643Z"
                decision: approve
                reason: Captain approved the scoped direction for ideation.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:6ker7h25hj86983e5ef71ahm:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:6ker7h25hj86983e5ef71ahm-ideation-1
              briefing:
                id: briefing:6ker7h25hj86983e5ef71ahm:ideation:attempt-1:revision-1
                digest: sha256:52695675916f49dddf7ff6ef1ff802a349f6dbe73694f74d72523e35c9dddf28
                request-digest: sha256:773d7f5796198796f52c789c4a5d551e46e604333ec1b278667dbff09b9c470c
                room-ref: ./repair-codex-filing-command-ledger-observation/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:6ker7h25hj86983e5ef71ahm:ideation:1
                briefing: briefing:6ker7h25hj86983e5ef71ahm:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-13T22:29:12.030156Z"
                decision: approve
                reason: Captain approved the correlated native-input/public-exit-status observation design for implementation.
              application:
                target-stage: implementation
                state: consumed
started: 2026-08-13T22:11:00Z
worktree: .worktrees/spacedock-ensign-repair-codex-filing-command-ledger-observation
---

## Problem

PR #679 run `31728107636`, Codex job `94541783359`, created `wire-the-thing.md` atomically with ID `001`, but `TestLiveCommonFiling` reported `filing command log has no spacedock new wire-the-thing invocation`. The outer `codex exec --json` stream rendered the successful command as a shell-display string with quote seams around `${SPACEDOCK_BIN:-spacedock}`; `commandFilesViaNew` could not recognize that display even though `item.completed` recorded exit code 0 and stdout recorded `created: .../wire-the-thing.md id=001`.

This is an observation defect in the Codex live-test adapter, not a product filing defect. The runner already owns an isolated `CODEX_HOME`, and exact `origin/main` already correlates its public `thread.started` ID to one parent rollout under `CODEX_HOME/sessions/**` for native lifecycle evidence. That public host session JSONL is the authoritative source for the pre-rendered execution input; the outer stream remains authoritative for completion and exit status.

## Risk spike

Before proposing a change, a detached throwaway checkout at exact `origin/main` `177eb454011001a296f1f09bc2889d3436df0a54` ran a focused fixture containing PR #679 artifact `9194350789`'s exact successful `item_9`. `successfulCodexCommands` decoded one exit-0 command, then `assertCodexFilingViaNew` failed with the production error because the rendered command was `/bin/bash -lc ... | \""'${SPACEDOCK_BIN:-spacedock}" new wire-the-thing'`. This reproduces the false negative without modifying product code or invoking the created entity.

The riskiest mechanism is already proven on this base: `codexNativeLifecycleStream` correlates exactly one public thread ID to exactly one isolated-home parent rollout and fails closed on missing or ambiguous rollouts. Implementation extends that established evidence path; it does not introduce a new host hook or state store.

## Proposed approach

For Codex live runs only, build the filing command ledger from the correlated parent rollout's native execution inputs and the matching completed outer-stream execution results. Pair execution observations in order within the single correlated thread, retain only calls whose public `item.completed` status is `completed` with exit code 0, and feed the native input text into the existing `assertFilingCommands` ladder. Fail closed when correlation or pairing is missing or ambiguous; never fall back to treating narration, stdout, the durable entity, or a malformed display string as command evidence.

This native-input observation serves AC-1. The simplest alternative was another `commandFilesViaNew` quote-shape regex, but PR #679 followed earlier parser corrections for captured-launcher display forms; adding another shell-display exception spends the correction budget without making the ledger authoritative. The other alternative was a binary wrapper/logger, but it would add test instrumentation to the product execution path and violate the test-product firewall.

Keep the command grammar unchanged: the existing atomic-create recognizer still requires a resolved Spacedock launcher, `new`/`--new`, and the requested slug in one command. Keep the negative semantics unchanged: `--next-id` remains manual-flow evidence, nonzero execution is excluded, and narration or an unrelated slug cannot satisfy filing.

## Verification ladder

One table-driven offline fixture exercises the same observation-to-grade path as the live Codex adapter:

1. A correlated native `${SPACEDOCK_BIN:-spacedock} new wire-the-thing` call paired with the exact PR #679 exit-0 public item passes.
2. A `--next-id` plus shell write is manual and fails.
3. The correct native `new wire-the-thing` paired with exit code 1 fails.
4. An exit-0 `new other-slug` fails.
5. A stream with no atomic-create command fails.
6. Missing/ambiguous parent rollout or mismatched native/public execution counts fails closed.

The positive case fails if the implementation continues grading the distorted display string. Each negative case fails if the observer invents success, ignores exit status, weakens slug binding, or accepts absence/manual creation.

## Out of scope

- No product hook, command, protocol, state store, lifecycle guard, global hook, or replacement for the compaction `SessionStart` hook.
- No change to `spacedock new`, command grammar, stored entity formats, write authority, or product runtime behavior.
- No work on Pi, Opus, rejection flow, supporting-evidence, mechanically-continue-Codex validation, another owner's target binding, or reconciliation row.
- PR #682's rejection-flow failure remains `live-evidence-followups` ownership.

## Acceptance criteria

- **AC-1 (VALUE):** The exact PR #679 successful filing shape produces one valid atomic `wire-the-thing` ledger observation from correlated public Codex host evidence, and the filing grade passes. Verified by the focused PR #679 fixture through the live adapter's observation function; replacing native input with the outer display makes it fail.
- **AC-2:** Manual, non-atomic, failed, wrong-slug, missing-command, and ambiguous-correlation cases remain failures. Verified by the six-rung table above through the same observation and grading path.
- **AC-3:** The delivered diff changes Codex test observation only; command grammar, stored formats, authority, and runtime product behavior are unchanged. Verified by a base-to-candidate changed-path check limited to `_test.go` and `testdata`, plus `go test ./...` and `go test ./... -race`.

## Test plan

- Add a default-tag offline fixture for the PR #679 outer event plus a correlated native rollout event. Run the focused test first; cost is milliseconds and no model/network spend.
- Add table-driven positive and red controls for success pairing, manual flow, failed execution, wrong slug, missing command, and correlation ambiguity. These are fixture tests over public Codex JSONL shapes.
- Run the exact targeted local Codex filing live test once on the final candidate as the single revalidation. It must create the entity and pass the ledger assertion; do not spend a second live correction on another display parser.
- Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` before delivery.

## Expected surface and semantic budget

Expected baseline: four test-only files, about 110 insertions and at most 15 deletions:

- `internal/ensigncycle/codex_live_runner_test.go` — source the Codex command ledger from correlated host evidence.
- `internal/ensigncycle/shared_filing_test.go` — bounded native/public execution pairing and fail-closed errors.
- `internal/ensigncycle/shared_filing_negative_test.go` — focused ladder assertions.
- `internal/ensigncycle/testdata/codex_filing_pr679/*.jsonl` — minimal public and parent-rollout fixtures carrying the reproduced shape.

Tolerance: one additional test-helper/fixture file and up to 60 additional insertions if separating dialect decoding makes the evidence path clearer; zero non-test product files. Allowed semantic change: Codex live-test command observation uses correlated native invocation input while preserving public exit status. No allowed change to CLI output, command grammar, stored formats, write authority, product runtime behavior, Claude/Pi grading, or documentation; no site doc diff is needed because user-visible behavior does not change.

## Stage Report: backlog

- DONE: Record the released workflow and observable harm.
- DONE: Preserve exact run, job, test, and artifact identity.
- DONE: State the product/test boundary.
- DONE: Define falsifiable positive and negative evidence.

### Summary

Repair the filing evidence boundary without changing successful filing behavior.

## Stage Report: ideation

- DONE: Produce a behavior-first design with a concrete approach, independently testable acceptance criteria, and expected files/insertions/tolerance plus all allowed semantic changes.
  The design uses correlated isolated-home Codex invocation evidence, names three falsifiable ACs, and budgets four test-only files/~110 insertions with explicit tolerance and a test-observation-only semantic allowance.
- DONE: Reproduce the PR #679 filing command shape with a focused read-only fixture against exact origin/main 177eb454011001a296f1f09bc2889d3436df0a54 before proposing any product change.
  A detached exact-base fixture decoded PR #679 artifact 9194350789 item_9 as one successful command and reproduced the existing `filing command log has no spacedock new wire-the-thing invocation` failure.
- DONE: Define one verification ladder that distinguishes valid atomic spacedock new evidence from manual, failed, wrong-slug, and missing-command cases while preserving the test-product firewall.
  The six-rung offline ladder grades native invocation plus public exit status and includes fail-closed correlation controls without hooks, product changes, or durable-state inference.

### Summary

Ideation identifies the outer Codex command display as the false-negative boundary and selects the already-established correlated session JSONL path as authoritative invocation evidence. The proposal preserves filing behavior and grammar, confines changes to test observation, and permits one focused correction followed by one live revalidation.

## Stage Report: implementation

- DONE: Deliver the approved correlated native-input/public-exit-status command ledger within the test-only surface and stated tolerance, preserving product behavior and command grammar.
  Commit ee53f53d2 changes six test/fixture files by 153 insertions and 8 deletions (net +145 within the +155 ceiling); the focused ladder fails if filing returns to the distorted public display or ignores public completion status.
- DONE: Add the six-case falsifiable offline ladder covering exact PR #679 success plus manual, failed, wrong-slug, missing-command, and ambiguous-correlation failures.
  `TestCorrelatedCodexFilingPR679Ladder` byte-matches archived artifact 9194350789 item_9 and goes red if native/public pairing invents success, weakens slug/manual checks, or accepts missing, ambiguous, or count-mismatched correlation.
- DONE: Commit the candidate and report focused results plus frozen verification status without spending the one exact live revalidation reserved for final validation.
  Focused tests, `go test -tags live ./internal/ensigncycle -run '^$' -count=1`, `go test ./...`, and `go test ./... -race` passed; the exact local Codex filing run remains frozen for validation.

### Summary

Implementation now grades Codex filing from the correlated parent rollout's decoded native command while retaining the public completed event as exit-status authority. The change is committed at ee53f53d2, remains entirely test-only, and leaves the reserved live revalidation unspent.
