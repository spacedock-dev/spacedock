---
title: Repair Codex filing command-ledger observation
status: ideation
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
            - id: gate-attempt:6ker7h25hj86983e5ef71ahm-ideation-2
              briefing:
                id: briefing:6ker7h25hj86983e5ef71ahm:ideation:attempt-2:revision-1
                digest: sha256:893f3a6684790a5a2840dff4c3047d3ebbd84710463b70e6d89564a52f8ed076
                request-digest: sha256:35ae55712eebf94ed279a8f2a1d8dd1b2cb7fd03f2eb3ac27bc45ee86eceb842
                room-ref: ./repair-codex-filing-command-ledger-observation/review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:6ker7h25hj86983e5ef71ahm:ideation:2
                briefing: briefing:6ker7h25hj86983e5ef71ahm:ideation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-14T07:23:04.943409Z"
                decision: approve
                reason: The execution-grounded test-only shim directly falsifies the observed counterfeit, preserves product behavior, and has a bounded signed surface and validation ladder.
              application:
                target-stage: implementation
                state: consumed
            - id: gate-attempt:6ker7h25hj86983e5ef71ahm-ideation-3
              briefing:
                id: briefing:6ker7h25hj86983e5ef71ahm:ideation:attempt-3:revision-1
                digest: sha256:8dba3579ef91a9ced2995a667077516b1aa5eada6e9595612bf180aa0cde23a3
                request-digest: sha256:d418fe0de2549e9429e0e823de4878e867fbf1b04a5f7d6b771e1ba4acc89fa3
                room-ref: ./repair-codex-filing-command-ledger-observation/review/ideation/briefing-3
              withdrawal:
                by: agent:first-officer
                at: "2026-08-14T14:33:25.477984Z"
                reason: Captain rejected the platform-specific supervisor after confirming the existing correlated thread-log infrastructure; replace the briefing with supported thread-log plus durable-state evidence.
            - id: gate-attempt:6ker7h25hj86983e5ef71ahm-ideation-4
              briefing:
                id: briefing:6ker7h25hj86983e5ef71ahm:ideation:attempt-4:revision-1
                digest: sha256:7fbcbc97d3dc332f6bd273ccfd6274623bd4e9547a04ae7b37cc43b530de2f15
                request-digest: sha256:a549e3c5a21bff888200c1fb8ee8701fcdad79705032963fb0fb9b47ff011ea8
                room-ref: ./repair-codex-filing-command-ledger-observation/review/ideation/briefing-4
              resolution:
                type: Resolution
                id: resolution:spacedock:6ker7h25hj86983e5ef71ahm:ideation:4
                briefing: briefing:6ker7h25hj86983e5ef71ahm:ideation:attempt-4:revision-1
                by: agent:first-officer
                at: "2026-08-14T14:38:34.725209Z"
                decision: approve
                reason: Captain selected the existing correlated thread-log direction; the revised design preserves supported atomic-path evidence, rejects the known counterfeits, removes the forgeable ledger, and stays within five test files.
              application:
                target-stage: implementation
                state: consumed
            - id: gate-attempt:6ker7h25hj86983e5ef71ahm-ideation-5
              briefing:
                id: briefing:6ker7h25hj86983e5ef71ahm:ideation:attempt-5:revision-1
                digest: sha256:8d81b2b8a93e6c7ac79aa60525fb6d5f583465b7e7b2e8c1f895ae966adfe09f
                request-digest: sha256:ec47e1fdae1e9e12174a6b7dad778d5096373be4265293ed61bbbb122110bdc8
                room-ref: ./repair-codex-filing-command-ledger-observation/review/ideation/briefing-5
        - id: gate:6ker7h25hj86983e5ef71ahm:validation
          stage: validation
          attempts:
            - id: gate-attempt:6ker7h25hj86983e5ef71ahm-validation-1
              briefing:
                id: briefing:6ker7h25hj86983e5ef71ahm:validation:attempt-1:revision-1
                digest: sha256:b289e57bb28da68f5799407b79a6aae67e25f9847e336d3533f1af2f9ec49e06
                request-digest: sha256:1f5e63a4e07c38f1f00682b031e9fad4f0d16c9da639262f032dc801ef88af5d
                room-ref: ./repair-codex-filing-command-ledger-observation/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:6ker7h25hj86983e5ef71ahm:validation:1
                briefing: briefing:6ker7h25hj86983e5ef71ahm:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-14T04:36:48.680856Z"
                decision: approve
                reason: Captain approved the validated filing observer for delivery.
              application:
                target-stage: done
                state: superseded
started: 2026-08-13T22:11:00Z
worktree: .worktrees/spacedock-ensign-repair-codex-filing-command-ledger-observation
mod-block:
pr:
---

## Problem

PR #679 run `31728107636`, Codex job `94541783359`, created `wire-the-thing.md` atomically with ID `001`, but `TestLiveCommonFiling` reported `filing command log has no spacedock new wire-the-thing invocation`. The outer `codex exec --json` stream rendered the successful command as a shell-display string with quote seams around `${SPACEDOCK_BIN:-spacedock}`; `commandFilesViaNew` could not recognize that display even though `item.completed` recorded exit code 0 and stdout recorded `created: .../wire-the-thing.md id=001`.

This is an observation defect in the Codex live-test adapter, not a product filing defect. Candidate `470994d23` is frozen after three observation directions failed. Native JavaScript occurrence paired unrelated success with unreachable dot-form text. The launcher-entry ledger was model-writable and accepted forged, duplicate, and mixed evidence. Correlated rollout inspection then required a nested Codex workspace sandbox inside the outer agent-safehouse; the supported host rejected every exec with `sandbox_apply: Operation not permitted`, so no filing command or entity could exist.

The captain has selected the harness-owned public `codex exec --json` stream as the evidence boundary. Its top-level `item.completed` / `command_execution` envelope is serialized by the Codex process whose stdout the harness captures directly. A command may print arbitrary bytes, including JSON-looking text and a `created:` line, but those bytes remain JSON-escaped inside that event's `aggregated_output`; the grader decodes each public line once and never treats nested output, agent prose, native JavaScript, a rollout, or a file as a second event. Binding the displayed command, completion/exit, exact same-item receipt, and independently verified entity therefore supplies supported evidence of one atomic filing transaction without claiming OS child-`execve` identity.

## Risk spike

No model run was spent. The retained PR #679 public fixture is already a real harness artifact and contains the exact lossy display, including the seams `printf '"'%s\\n'` and `| \""'$sd_bin" new wire-the-thing`. On frozen `470994d23`, `go test ./internal/ensigncycle -run '^TestAssertCodexFilingViaNew$' -count=1` passed, proving that the bounded launcher recognizer already accepts those exact display bytes while rejecting its existing manual, wrong-slug, narration, malformed-quote, and cross-segment controls.

The harness implementation opens `codex-exec.jsonl` itself and assigns the launched front door's stdout to an `io.MultiWriter` feeding that file and the line stream. The model controls command and command-output bytes, but not sibling bytes on that captured public pipe: Codex owns the public JSON serialization, and command output occupies `aggregated_output`. The remaining unverified mechanism is the one-item conjunction and exact create count, so implementation must first retain the exact PR #679 completed event and make all detached public-stream counterfeits pass before any further live run. The exact `TestLiveCommonFiling/codex` run remains reserved for validation.

## Proposed approach

For the Codex `filing` journey only, replace correlated rollout grading with one structural pass over the already-captured public JSONL. Extend the existing public `codexCommandItem` shape with `aggregated_output`; do not read a rollout or ledger, recursively decode output, parse native JavaScript, supervise processes, or add another event source. Remove the rejected filing-only nested sandbox override and return filing to the existing common Codex front door; this introduces no new sandbox posture.

Grade the public evidence as one bounded transaction:

1. Decode only top-level `type=item.completed`, `item.type=command_execution` lines. Count exact atomic `new wire-the-thing` / `--new wire-the-thing` invocations in each displayed command with the existing bounded launcher grammar, extended from boolean recognition to a count. The exact PR #679 display must count once.
2. Across all observed command-execution items, require exactly one create count. Reject a second recognized create within the same displayed command or in another item. Reject any observed `--next-id` in a successful or failed item, including one mixed into the filing item.
3. The sole create-owning item must be completed with exit code 0 and must own exactly one line `created: <canonical-entity-path> id=<landed-id>` in its `aggregated_output`. Another item's output, agent prose, a command-printed fake event, or a detached receipt cannot lend evidence.
4. Independently require the canonical landed entity with expected title, backlog status, one-line nonblank body, receipt-matching stamped ID, and the existing clean durable workflow result.

This serves AC-1 and AC-2 without trusting any individually model-writable byte. The top-level envelope proves which public command item Codex reports as executed; command, exit, receipt, and entity must agree inside the one bounded value transaction. Unreachable `tools.exec_command(...)` text remains agent/native content rather than a public command event, and an unrelated successful public command contains zero recognized creates. The simplest alternative is public command occurrence alone; it is insufficient because it can ignore exit, accept duplicates, or borrow a detached receipt/entity. Rollout correlation is also insufficient because securing rollout authority required the unsupported nested sandbox. Platform process supervision was rejected as unnecessary and platform-specific.

## Verification ladder

The seven-rung no-model public-fixture matrix exercises the selected transaction:

1. One completed exit-0 public command item with one direct or bound/PATH-resolved atomic alias, one exact same-item receipt, and the correct entity passes.
2. The exact retained PR #679 command display and `aggregated_output` pass byte-for-byte, including all quote seams; changing a seam, event type, status, exit, output, or entity turns it red.
3. Unreachable native `tools.exec_command` text plus an unrelated successful public execution fails with zero creates, even beside a command-printed fake event, forged receipt, and valid hand-written entity.
4. Duplicate create evidence within one displayed command or across two items fails the exactly-one count; mixed direct/bound aliases also fail.
5. Any observed `--next-id` fails regardless of status/exit and regardless of a valid create/receipt/entity beside it.
6. Failed, wrong-slug, manual, missing, malformed, started-only, and wrong-item commands; detached/missing/duplicate/wrong-path/wrong-ID receipts; and wrong/missing entities all fail.
7. A command that prints a complete fake `item.completed/command_execution` JSON object retains it only inside `aggregated_output` and fails recursive-injection controls; removing the rejected filing-only argv helper restores byte-identical common runner posture.

Implementation runs the exact retained fixture and every detached counterfeit before live-tag compile, `gofmt`, full, and race. No implementation-stage model run is allowed. Validation alone may spend one exact Codex filing run on frozen final bytes and requires one accepted public atomic-path transaction, exact owned receipt and entity, and zero `--next-id` evidence.

## Out of scope

- No product hook, command, protocol, state store, lifecycle guard, global hook, process supervisor, file ledger, rollout authority, or replacement for the compaction `SessionStart` hook.
- No change to `spacedock new`, command grammar, stored entity formats, write authority, or product runtime behavior.
- No work on Pi, Opus, rejection flow, supporting-evidence, mechanically-continue-Codex validation, another owner's target binding, or reconciliation row.
- PR #682's rejection-flow failure remains `live-evidence-followups` ownership.

## Acceptance criteria

- **AC-1 (VALUE):** Supported public Codex evidence establishes exactly one completed exit-0 atomic filing transaction, its exact same-item `created: <canonical-path> id=<landed-id>` receipt, and the correct durable entity; unreachable native filing text plus unrelated public success establishes zero. Verified by matrix rungs 1-3 and the sole exact Codex filing live run; accepting nested output as an event, detaching the receipt, or omitting the landed-entity check makes a named counterfeit pass.
- **AC-2:** Direct, bound, and PATH-resolved atomic aliases pass only as one public transaction; duplicate or mixed creates, any observed `--next-id` regardless of exit, manual/non-atomic, failed, wrong-slug, missing, malformed, command-printed-event, wrong-item, wrong-receipt, and entity-mismatch cases remain failures. Verified by matrix rungs 4-7 through the same top-level public-event decoder and grade.
- **AC-3:** The delivered diff changes only Codex filing-test observation and removes the rejected filing-only rollout/sandbox experiment; product bytes, CLI grammar, stored formats, authority, shared Codex front-door posture, other runtimes, and non-filing journeys are unchanged. Verified by exact changed paths and signed LOC, common-runner argv parity, live-tag compile, `go test ./...`, and `go test ./... -race`.

## Test plan

- Replace the correlated-rollout matrix with table-driven public `item.completed/command_execution` transaction tests for all seven rungs. Preserve the exact retained PR #679 command and output fixture, not a normalized approximation.
- Extend the existing public command decoder with `aggregated_output`, exact create counting, same-item receipt ownership, and durable entity verification. Tests mutate event type/status/exit, alias/count, receipt ownership/bytes, and entity bytes independently; they never parse native JavaScript or infer child processes.
- Add detached agent-message, command-printed-event, unrelated-success, duplicate/mixed-create, failed-`--next-id`, and no-recursive-decode controls. Remove correlated-rollout-only and filing-only sandbox argv tests; existing front-door tests prove common posture.
- Run the exact targeted Codex filing live test once on frozen final bytes in validation, only after the retained fixture and every detached counterfeit pass. Do not spend another live correction on observation parsing.
- Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` before delivery.

## Expected surface and semantic budget

Expected implementation delta from frozen `470994d23`: **-055 net LOC across 6 test files** (about 115 insertions and 170 deletions), accepted range **-090 to -020**:

- `internal/ensigncycle/codex_filing_invocation_ledger_test.go` — delete the correlated-rollout transaction decoder/grade left under the obsolete filename.
- `internal/ensigncycle/codex_live_runner_test.go` — grade filing from the captured public stream and restore the common front-door argv path.
- `internal/ensigncycle/claude_runtime_helpers_test.go` — remove the filing-only rollout-return helper while preserving native lifecycle correlation used by other journeys.
- `internal/ensigncycle/shared_filing_test.go` — extend the public event shape and bounded filing transaction/count/receipt/entity grade.
- `internal/ensigncycle/shared_filing_negative_test.go` — replace correlated evidence cases with the exact PR #679 fixture and seven-rung public counterfeit matrix.
- `internal/ensigncycle/team_capability_test.go` — remove the rejected filing-only nested-sandbox argv helper/tests; retain the existing common front-door checks.

Allowed semantic change: Codex filing-test evidence moves from correlated/model-writable sources to one harness-owned public command-execution transaction with exact same-item receipt and durable entity; the rejected filing-only nested-sandbox override is removed so filing again uses the existing common front door. No allowed change to CLI output, command grammar, stored formats, write authority, product runtime behavior, Claude/Pi grading, other journeys, common Codex posture, or documentation. The candidate remains frozen at `470994d23`; the exact Codex filing model run remains reserved for validation.

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

## Stage Report: validation

- DONE: Independently reproduce AC-1 through AC-3 with the exact PR #679 correlated native/public evidence and all six fail-closed controls, confirming only test observation semantics changed.
  AC-1 byte-matched archived artifact 9194350789 `item_9` (SHA-256 `6b2fc61f…`) and passed focused/live grading; AC-2 rejected manual/non-atomic, exit-nonzero, wrong-slug, missing-command, missing/ambiguous rollout, and count mismatch; AC-3 found only six test/fixture paths in ee53f53d2.
- DONE: Validate the captain-approved +145 net LOC across 6 test-only files, run focused/full/race checks, detached adversarial audit, and the sole exact local Codex filing live target on frozen bytes.
  Commit ee53f53d2 is 153 insertions/8 deletions; focused ladder, `go test ./...`, `go test ./... -race`, detached audit `TestDetachedCodexLedgerAdversarialAudit`, and the sole exact live command passed (`TestLiveCommonFiling`, 46.81s) with clean HEAD before and after.
- DONE: Report PASSED or REJECTED with exact run/artifact identifiers, every failure classified from this run, and confirmation that no product file or other owner binding/reconciliation row changed.
  PASSED for PR #679 run 31728107636, job 94541783359, artifact 9194350789, and local candidate ee53f53d2; the exact local command intentionally retained no artifact directory, its XPASS alert was the expected repaired characterization, no test/audit failure or finding occurred, and no product, owner-binding, or reconciliation-row path changed.

### Summary

PASSED. Exact archived and fresh live evidence establish one native `spacedock new wire-the-thing` observation paired to public completed exit 0, while every specified adjacent failure remains closed. No material, deferred-risk, or polish finding was discovered; delivery remains confined to the captain-approved six-file test surface at net +145 LOC.

## Stage Report: implementation (cycle 2)

- SKIPPED: Reproduce the PR #686 Codex 0.147 failure in an isolated home, retain the actual parent-rollout tool input, and add the smallest focused failing fixture before changing the decoder.
  Artifact 9210380407 omitted `CODEX_HOME`, and local OAuth 0.147 emitted the original supported layout rather than the CI shape; the FO explicitly accepted this limitation, while `TestNativeCodexCommandStructuralDecoder` first reproduced the spelling/layout failure offline.
- DONE: Decode the observed native invocation structurally while preserving fail-closed behavior for malformed, missing, ambiguous, failed, manual, wrong-slug, and count-mismatch cases; do not change product/runtime bytes.
  Commit 77a71f8e0 tokenizes exactly one flat `tools.exec_command` call and one nonblank JSON-string `cmd`; the focused suite goes red if layout tolerance regresses or malformed, duplicate, multiple-call, manual, failed, wrong-slug, correlation, or pairing controls open.
- DONE: Run focused, full, race, detached controls, then one exact Codex filing target on frozen bytes; report the final signed LOC and stop on any new finding.
  Focused/live-tag/full/race and detached controls passed, then frozen commit 77a71f8e0 passed exact `TestLiveCommonFiling` in 68.78s; final surface is 232 insertions/8 deletions (net +224) across the same six test-only files, with no new finding.

### Summary

Cycle 2 replaces the exact marker with a bounded structural decoder and preserves every original filing and correlation control. The CI-native bytes remain unavailable by explicit FO acknowledgment; a pre-acknowledgment local OAuth diagnostic did not reproduce that layout, while the required frozen exact target passed on the committed correction.

## Stage Report: validation (cycle 2)

- DONE: Independently attack the bounded structural exec-command decoder across whitespace, pragma, variable-name, malformed, ambiguous, duplicate, multiple-call, non-string, blank, and unsupported shapes while preserving every filing control.
  Focused controls passed, but a detached end-to-end adversarial test failed: a computed real `echo` call plus one unreachable dot-form `spacedock new wire-the-thing` call was decoded as the filing command and passed the grade.
- DONE: Verify focused, live-tag compile, full, race, detached, diff, and current-origin/main merge-tree evidence; separate any current-main Pi compile regression from this six-file candidate.
  Focused, candidate live-tag compile, `go test ./...`, `go test ./... -race`, `gofmt -d`, `git diff --check`, and merge-tree `3c54efc6…` passed; the detached audit found the ledger defect described above.
- DONE: Classify the current-origin/main integration failure separately from filing ownership.
  Exact `origin/main` 024a507ab and its clean synthetic merge both fail only at `pi_live_runner_test.go:285:48: undefined: defaultPiLiveModel`; commit 77a71f8e0 changes no Pi path and does not own that compile regression.
- DONE: Validate the retained exact filing PASS, report +224 net LOC across 6 test-only files against the earlier +145 estimate, and recommend PASSED or REJECTED without another live run.
  Frozen HEAD remains 77a71f8e0; implementation state commit c8a3df62d retains the exact 68.78s `TestLiveCommonFiling` PASS. The base diff is 232 insertions/8 deletions, net +224 across six test-only files, +79 over the earlier +145 estimate; no live target was rerun.
- DONE: Reproduce AC-1 evidence.
  `TestCorrelatedCodexFilingPR679Ladder` passed the archived PR #679 native/public fixture and still fails when native input is replaced by the distorted public display.
- FAILED: Reproduce AC-2 evidence for every manual, non-atomic, failed, wrong-slug, missing-command, and ambiguous-correlation case.
  Existing controls pass, but the detached supported-shape boundary proves a missing atomic command can pass: the observer pairs public exit 0 from a computed unrelated call with an unexecuted dot-form filing call found by token scan.
- DONE: Reproduce AC-3 evidence that only Codex test observation changed.
  Base-to-candidate paths are exactly six `_test.go`/`testdata` files; candidate full and race suites pass, with no product/runtime file changed.
- FAILED: Recommend PASSED or REJECTED with findings classified independently by defect kind and release scope.
  REJECTED: Material evidence defect and chosen-mechanism failure. The promised unsupported/missing-command boundary and `docs/dev/README.md#validation` evidence-integrity contract fail; token occurrence cannot prove which JavaScript call executed, so recommend a scope/design reset rather than a third automatic parser correction.

### Summary

REJECTED. The correction handles the observed spelling/layout variations and retains the exact live PASS, but its token scan can attribute an unrelated successful execution to an unreachable filing call, allowing the filing grade to pass without atomic creation. The unrelated 024a507ab Pi live-tag compile defect remains separate, and candidate bytes were not changed or rerun live.

## Stage Report: validation (cycle 3)

- DONE: Preserve the detached finding that an unrelated executed command plus an unreachable dot-form filing call can pass the ledger grade without atomic creation.
  The cycle-2 finding remains unchanged: public exit 0 from the computed unrelated call was paired with the unexecuted dot-form `spacedock new wire-the-thing`, and the filing grade incorrectly passed.
- DONE: Record that the token-scan mechanism failed AC-2 and requires a normal design reset to ideation, with candidate 77a71f8e0 frozen.
  REJECTED remains the durable outcome: Material, task-owned evidence defect and chosen-mechanism failure; clean candidate HEAD is frozen at 77a71f8e03cc2f2280ad69d47be58c36c4c28a74 for the validation-to-ideation handoff.
- SKIPPED: Candidate edits and decoder corrections.
  The dispatch forbids candidate mutation and another parser correction; validation only preserves the rejected snapshot and routes a normal design reset.
- SKIPPED: Candidate and live test reruns.
  The dispatch forbids reruns; cycle-2 evidence and the retained exact filing PASS remain on record without claiming PASSED.
- DONE: Leave a complete validation-stage handoff that permits the First Officer to move back to ideation without claiming PASSED.
  The finding, AC-2 failure, Material evidence-defect classification, task ownership, mechanism-failure disposition, frozen SHA, and ideation reset are all explicit above.

### Summary

REJECTED is preserved without changing candidate bytes, evidence, or findings. Validation hands frozen commit 77a71f8e0 back to ideation for a normal design reset; no correction or test rerun occurred.

## Stage Report: ideation (cycle 2)

- DONE: Trace the complete filing evidence chain and identify the cheapest observable that distinguishes an actually executed atomic `spacedock new` from an unreachable or merely mentioned call.
  Native call occurrence, public completion, durable entity bytes, and Git history remain non-injective; exact successful argv at a test-local launcher shim is the first execution-grounded discriminator.
- DONE: Design the smallest mechanism that satisfies AC-1 and AC-2 without a permissive JavaScript token scan, preserving product/runtime bytes if the evidence can support that boundary; state clearly if it cannot.
  The design preserves product bytes and all real binary I/O/exit behavior but explicitly permits a Codex-filing-only shim because passive trace/runtime bytes cannot prove which JavaScript or shell branch executed.
- DONE: Specify the focused falsifier matrix, exact bound/unbound or live ladder, exact files, and signed net LOC estimate as +NNN or -NNN across M files; do not spend a live run.
  Eight offline/no-model rungs plus one validation-only live target are specified; the frozen-to-final estimate is -034 net LOC across 7 exact paths (base-to-final +190 across 3 test files), and no live run was spent.

### Summary

Ideation resets the observer from trace interpretation to an execution-grounded, Codex-filing-only argv/exit ledger. The detached counterfeit becomes the primary negative control, the frozen parser candidate remains unchanged, and the plan excludes Pi, BFM, rejection flow, timeout, Sonnet, and other-owner work.

## Stage Report: implementation (cycle 3)

- DONE: Build the executable-ledger falsifier first: the detached counterfeit records no successful create, while bound and PATH execution record exact argv and exit status.
  `TestCodexFilingInvocationLedgerExecutionMatrix` fails if mentioned/unreachable text creates a record, if bound/PATH argv or real exit status are lost, or if manual, failed, wrong, missing, `--next-id`, or concurrent cases are misgraded; `TestCodexFilingInvocationLedgerFailsClosed` rejects truncated, duplicate-terminal, unknown-tool, and malformed records.
- DONE: Replace the frozen parser candidate only within the approved seven-path surface; preserve real I/O, exit behavior, product bytes, and every negative control.
  Commit 4faa23656 adds the transparent Codex-filing shim, restores the pre-task parser/lifecycle boundary, and deletes parser-only fixtures; the frozen-to-final diff is exactly seven approved test paths and 214 insertions/208 deletions (net +006, within -034 +/-40), with no product file.
- DONE: Run all eight offline rungs plus focused, live-tag compile, full, race, gofmt, and diff checks; reserve the exact Codex filing live run for validation.
  The eight-rung matrix and one-iteration front-door binding probe pass; `go test -tags=live ./internal/ensigncycle -run '^$'`, `go test ./...`, `go test ./... -race`, `gofmt -d ./cmd ./internal`, and `git diff --check` pass, and replacing executable entry with mere command text makes the detached rung fail. The exact model-backed filing target was not run.

### Summary

Codex filing evidence now comes from an atomic NUL-delimited launcher ledger containing exact argv and the real exit status, with both `SPACEDOCK_BIN` and PATH rebound inside the Codex host shim after front-door pinning. The live journey retains real stdin/stdout/stderr, filesystem behavior, and the independent landed-entity assertion while removing the non-injective native JavaScript/public-display parser.

## Stage Report: validation (cycle 4)

- FAILED: Reproduce AC-1 and AC-2 with the eight-rung execution-ledger matrix, then spend the single exact Codex filing live run on frozen commit 4faa23656 and verify exact argv, exit 0, landed entity, and no --next-id record.
  The eight rungs and one-iteration front-door probe passed, but the detached audit falsified AC-1/AC-2 before the reserved model run; no live run was spent.
- SKIPPED: Verify AC-3 with the exact changed-path and signed LOC checks, live-tag compile, full and race suites, gofmt, and unchanged product/runtime semantics.
  Static diff evidence is exactly 7 approved paths and 214 insertions/208 deletions (net +006), but the finding-stop rule prohibited further suite and formatting reruns.
- DONE: Run a detached adversarial audit that tries to make unreachable, forged, malformed, duplicate, failed, wrong-slug, manual, or concurrent evidence pass; report findings through the workflow disposition policy and recommend PASSED or REJECTED.
  Throwaway checkout `/tmp/spacedock-ledger-audit.f2Tey3/checkout` at 4faa23656 produced three failing counterexamples; recommendation is REJECTED.
- FAILED: Reproduce AC-1 evidence.
  The focused bound/PATH/counterfeit matrix passed, but a hand-written valid record passed without executable entry and two successful create records passed despite AC-1 requiring exactly one.
- FAILED: Reproduce AC-2 evidence.
  Existing malformed, wrong-slug, missing, manual, failed-create, and concurrent-integrity controls passed, but a nonzero `status --next-id` record alongside a successful create was ignored and graded successful.
- SKIPPED: Reproduce AC-3 evidence.
  The signed frozen-to-final diff is +006 across the exact seven test-only paths, with zero product paths; full/race/gofmt were not rerun after the finding.
- DONE: Preserve reviewer finding and recommend disposition without changing or rerunning the candidate.
  `go test ./internal/ensigncycle -run '^TestDetachedCodexFilingLedgerAdversarialAudit$' -count=1 -v` failed all three adversarial subtests; candidate HEAD remains 4faa23656.

### Review-finding disposition

- Released user and normal workflow: the supported Codex filing journey gives model-run shell commands the ledger directory through `SPACEDOCK_CODEX_FILING_LEDGER_DIR`, then trusts every matching `invocation.*` record.
- Observable harm: valid hand-written evidence can establish filing without launcher execution, and duplicate successful creates do not violate the promised exactly-one observation.
- Authority: value-ac[AC-1] requires one real exit-0 launcher execution, exactly one successful observation, and zero observations for detached counterfeit evidence.
- Trigger evidence: detached frozen-checkout subtests `well formed forged record` and `duplicate successful create` both reached a nil filing grade without the required unique launcher execution.
- Proposed classification: Material evidence defect, task-owned, chosen-mechanism failure; hold frozen bytes and route for First Officer disposition before any mutation or rerun.
- Separate AC-2 evidence: value-ac[AC-2] requires `--next-id` cases to remain failures; `failed next-id alongside successful create` instead received a nil filing grade because nonzero records are skipped before the guard.

### Summary

REJECTED. The intended eight-rung controls pass and the approved signed surface is exact, but the detached audit proves the ledger reader can accept forged and duplicate success and can ignore an executed failed `--next-id`. The reserved exact Codex filing model run remains unspent, and frozen candidate 4faa23656 was neither changed nor rerun after the finding.

## Stage Report: ideation (cycle 3)

- DONE: Trace why the launcher ledger remains forgeable and spike the cheapest observable that proves real executable entry against all three detached counterexamples.
  The ledger trusts a model-writable pathname; Codex 0.147 runtime `CommandExecution` was spiked without a model run and proved insufficient because unreachable, duplicate, and failed-then-successful children each collapse to one completed outer-shell event with exit 0.
- DONE: Design the smallest revised mechanism that satisfies AC-1 and AC-2 without trusting model-writable evidence; if no supported boundary can do this, state the exact limitation and required scope decision instead of adding another parser or shim patch.
  No supported Codex surface exposes child `execve` identity; implementation is held for a captain choice between new peer-verified test supervision, relaxed outer-command ACs, or a product receipt, with Option A recommended to preserve end value.
- DONE: Update the approach, acceptance evidence, semantic budget, exact paths, and signed net LOC estimate; keep the exact Codex filing live run unspent for validation.
  AC-1/AC-2 retain the three falsifiers; verification now tests authority, Options A/B carry exact provisional paths and +180/-050 signed estimates, candidate remains 4faa23656, and the exact live filing run was not executed.

### Summary

Cycle 3 establishes that both the file ledger and Codex's supported outer-command events stop short of trustworthy child executable entry. The task now requires an explicit scope decision rather than another parser or model-writable attestation patch; no candidate bytes or reserved live evidence changed.

## Stage Report: ideation (cycle 4)

- DONE: Replace the platform-supervisor options with one selected direction that reuses the existing correlated Codex thread-log infrastructure and the durable landed entity.
  The captain-selected design correlates the harness thread ID to one isolated rollout, grades runtime `CommandExecution` transactions plus the exact receipt/entity, and removes the rejected file ledger and process supervisor.
- DONE: Define supported evidence and negative controls: completed command execution, atomic new invocation, exact created receipt, correct entity, and rejection of any --next-id or unreachable/mixed-command counterfeit.
  Seven no-model rungs require one completed exit-0 atomic command, its owned exact receipt and correct entity; unreachable native text, unrelated success, duplicate/mixed create, any `--next-id`, failed/wrong/missing evidence, and ambiguous rollout fail.
- DONE: Update ACs, test plan, semantic budget, exact paths, and signed estimate to -050 net LOC across 5 files from frozen 4faa23656; keep the live run unspent.
  AC-1 now measures supported atomic-path evidence, AC-2 preserves fail-closed controls, AC-3 fences semantics; the exact five test paths budget about 150 insertions/200 deletions in the accepted -085..-015 range, and no model run occurred.

### Summary

Cycle 4 adopts the captain's supported-evidence boundary and removes the unresolved platform alternatives. The design reuses existing Codex thread correlation, adds receipt ownership and durable-state checks, keeps product and other journeys unchanged, and reserves the sole exact live filing run for validation.

## Stage Report: implementation (cycle 4)

- DONE: Build the seven-rung correlated CommandExecution transaction matrix first; every named unreachable, unrelated, duplicate, mixed, next-id, receipt, rollout, and entity counterfeit must fail.
  `TestCorrelatedCodexFilingTransactionMatrix`, structural-extraction controls, and correlated-rollout lookup fail if public/unreachable text becomes execution, if duplicate/mixed evidence is accepted, or if receipt/entity/correlation bytes are weakened.
- DONE: Replace the forgeable ledger with existing thread-to-rollout correlation and exact receipt/entity grading within the approved five test paths and -085 to -015 net LOC range.
  Commit e851be3fe deletes the launcher shim/ledger, reuses the lifecycle correlation, and delivers exactly 5 test paths with 174 insertions/190 deletions (net -016).
- DONE: Run focused, filing/non-filing argv, live-tag compile, gofmt, full, race, and exact diff checks; preserve the sole exact Codex filing model run for validation.
  Focused matrix/argv tests, live-tag compile, `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, `git diff --check`, and exact path/LOC checks passed; the model-backed filing target was not run.

### Summary

Codex filing evidence now comes from the single correlated rollout's runtime-owned completed `CommandExecution`, with exactly one atomic command, its exact same-item receipt, and an independently checked landed entity. Only filing uses `workspace-write`/`never`; non-filing argv is byte-identical, and the reserved validation run remains unspent.

## Stage Report: validation (cycle 5)

- FAILED: Reproduce AC-1 with the seven-rung correlated transaction matrix, then spend the single exact Codex filing run on frozen commit e851be3fe and verify one owned atomic command, exact receipt, correct durable entity, and zero unreachable/unrelated evidence.
  The matrix passed, but the sole exact run exited 2 before model execution because Codex CLI 0.147.0 does not accept the candidate's filing-only `--ask-for-approval never`; no transaction, receipt, or entity could be observed and no rerun was made.
- DONE: Reproduce AC-2 by attacking duplicate/mixed creates, any --next-id regardless of exit, failed/wrong/missing commands, forged public evidence, receipt ownership/bytes, malformed or ambiguous rollout, and entity mismatch.
  Focused tests plus detached `TestDetachedCodexFilingCorrelatedTransactionAudit`/`CorrelationAudit` rejected all named counterfeits and accepted bound, PATH, captured-launcher, and `status --new` aliases; weakening uniqueness, receipt ownership/bytes, rollout identity, or entity bytes makes a named case pass.
- FAILED: Verify AC-3 with exact five-path -016 LOC evidence, filing-only workspace-write/never and non-filing argv parity, live-tag compile, gofmt, full, race, and a detached adversarial audit; recommend PASSED or REJECTED without editing candidate bytes.
  Five test paths and 174 insertions/190 deletions give exact net -016; argv parity, live-tag compile, gofmt, full suite, and detached audit passed, but the filing argv's approval flag is incompatible with live Codex and the finding-stop rule left race unrun.
- DONE: Reproduce AC-1 evidence.
  Offline rungs 1-3 passed and reject public/unreachable text plus unrelated success, but the required sole live proof failed at host argument parsing before a `thread.started` or `CommandExecution` could exist.
- DONE: Reproduce AC-2 evidence.
  Rungs 4-7 and the detached matrix reject duplicate creates within/across commands, successful or failed `--next-id`, failed/wrong/manual/missing creates, forged public evidence, wrong-item or inexact receipts, correlation ambiguity, and durable mismatch.
- FAILED: Reproduce AC-3 evidence.
  Static surface, formatting, compile, full suite, and non-filing parity are exact, but a unit test asserted an argv spelling unsupported by the installed live-tag host and therefore did not prove filing-only workspace-write/never compatibility.
- DONE: Run a detached adversarial audit and recommend PASSED or REJECTED without editing candidate bytes.
  Throwaway checkout `/tmp/spacedock-codex-filing-audit.9pvbEM/checkout` at e851be3fe passed 31 alias/counterfeit/entity/correlation cases; the subsequent exact live run independently found a material host-argv defect, so recommendation is REJECTED.
- DONE: Preserve the new finding and stop before candidate mutation or rerun for First Officer disposition.
  Candidate HEAD remains e851be3fe; the one authorized model-run command was invoked exactly once, exited before model execution, and was not retried.

### Review-finding disposition

- Released user and normal workflow: the required supported Codex filing validation path constructs `codex exec --sandbox workspace-write --ask-for-approval never` through the normal Spacedock front door on installed Codex CLI 0.147.0.
- Observable harm: the exact filing journey exits 2 during host argument parsing, so it cannot execute an atomic filing command or produce the correlated receipt and durable entity required for release evidence.
- Authority: value-ac[AC-1] requires the sole exact Codex filing live run to establish one completed exit-0 atomic filing transaction, its exact receipt, and the correct durable entity.
- Trigger evidence: `SPACEDOCK_LIVE_RUNTIME=codex go test -tags=live -count=1 -timeout 40m -run '^TestLiveCommonFiling$' ./internal/ensigncycle -v` failed at `codex_live_runner_test.go:38` with `codex exec exited 2`; `codex exec --help` exposes `--sandbox` but no `--ask-for-approval` option.
- Proposed classification: Material evidence defect, task-owned, narrow candidate fix affecting AC-1's live observation boundary; hold e851be3fe and route for distinct First Officer authorization before mutation or rerun.

### Summary

REJECTED. The correlated transaction observer and detached falsifier matrix are clean, and the candidate has the exact approved five-file net -016 surface, but its filing-only approval argv cannot launch the installed supported Codex host. The candidate was not edited, the exact command was not rerun, and the race suite was intentionally stopped after the new finding.

## Stage Report: implementation (cycle 5)

- DONE: Move the filing-only sandbox and approval flags to the Codex top-level position before exec, as proven by the existing launcher-order test; change no evidence semantics.
  Commit 354f4e362 changes only two argv/test lines so the filing host suffix starts `--sandbox workspace-write --ask-for-approval never exec`; correlated transaction files and semantics are byte-identical to e851be3fe.
- DONE: Add or correct the argv test so the exact Spacedock host fence produces --sandbox workspace-write --ask-for-approval never exec, while non-filing argv remains byte-identical.
  `TestCodexLiveRunnerUsesRestrictedPostureOnlyForFiling` compares non-filing argv with `slices.Equal` and fails if the exact filing prefix order changes or bypass posture returns.
- DONE: Run the focused argv/matrix checks, live-tag compile, gofmt, full, race, and exact five-path/LOC checks; do not spend another model run.
  Focused argv/matrix/extraction and launcher-order tests, live-tag compile, `gofmt -w ./cmd ./internal`, full, race, diff, and surface checks passed; aggregate delta remains 174 insertions/190 deletions (net -016) across 5 paths, and no model run occurred.

### Summary

The authorized correction moves only the filing permission flags to Codex's supported top-level position before `exec`. The transaction observer, ACs, product bytes, and every non-filing journey remain unchanged; validation retains ownership of the model-backed rerun.

## Stage Report: validation (cycle 6)

- FAILED: Reproduce AC-1 with the seven-rung correlated transaction matrix, then spend the single exact Codex filing run on frozen commit e851be3fe and verify one owned atomic command, exact receipt, correct durable entity, and zero unreachable/unrelated evidence.
  On corrected frozen commit 354f4e362, rungs 1-3 passed and a no-model front-door probe proved the fixed argv launches, but the authorized replacement live run found no completed `CommandExecution` item and therefore established no atomic command, receipt, or entity evidence.
- DONE: Reproduce AC-2 by attacking duplicate/mixed creates, any --next-id regardless of exit, failed/wrong/missing commands, forged public evidence, receipt ownership/bytes, malformed or ambiguous rollout, and entity mismatch.
  The unchanged focused matrix, structural extraction, and correlated-rollout tests passed; cycle-5's detached 31-case audit already covers the unchanged transaction semantics and fails if any named counterfeit is accepted.
- FAILED: Verify AC-3 with exact five-path -016 LOC evidence, filing-only workspace-write/never and non-filing argv parity, live-tag compile, gofmt, full, race, and a detached adversarial audit; recommend PASSED or REJECTED without editing candidate bytes.
  The aggregate diff remains exactly 5 approved test paths and net -016, focused argv parity passed, and the candidate launcher accepted the exact top-level permission order; the new live finding stopped this pass before live-tag/full/race/gofmt and a fresh detached audit.
- FAILED: Reproduce AC-1 evidence.
  `TestLiveCommonFiling` launched the host on 354f4e362 but failed at `codex_live_runner_test.go:48` because the correlated rollout decoded to zero completed command executions.
- DONE: Reproduce AC-2 evidence.
  Rungs 4-7, structural extraction, missing/ambiguous rollout correlation, and the prior detached same-semantics audit remain green; the argv-only correction does not touch these bytes.
- FAILED: Reproduce AC-3 evidence.
  The exact signed surface and corrected top-level `--sandbox workspace-write --ask-for-approval never exec` order are proven, but the required live behavior failed and the finding-stop rule prohibited the remaining validation reruns.
- DONE: Preserve the new finding and stop before candidate mutation or another live run for First Officer disposition.
  Candidate HEAD remains 354f4e362; the authorized replacement model run was invoked exactly once and was not retried.

### Review-finding disposition

- Released user and normal workflow: the required exact Codex filing journey on the supported Codex CLI 0.147.0, using the corrected top-level workspace-write/never host argv.
- Observable harm: the host launched, but its correlated rollout contained no completed `CommandExecution`; the live target could not establish the promised atomic filing path, same-item receipt, or durable entity.
- Authority: value-ac[AC-1] requires the exact Codex filing live run to establish one completed exit-0 atomic filing transaction, its exact receipt, and the correct durable entity.
- Trigger evidence: `SPACEDOCK_LIVE_RUNTIME=codex go test -tags=live -count=1 -timeout 40m -run '^TestLiveCommonFiling$' ./internal/ensigncycle -v` failed after 42.22s with `correlated Codex rollout has no completed CommandExecution items`; the preceding candidate front-door `exec --help` probe exited 0 with the corrected flags.
- Proposed classification: Material outcome defect; current task ownership requires First Officer investigation because the observer correctly rejected an actual supported live run with no command execution; hold 354f4e362 and do not mutate or rerun pending distinct disposition.

### Summary

REJECTED. The narrow argv-order correction fixes host launch and preserves the exact approved net -016 surface, while all unchanged offline transaction controls remain green. The authorized live rerun nevertheless produced zero completed command executions, so AC-1 lacks its required live proof; candidate bytes were not edited and no second run was attempted.

## Stage Report: validation (cycle 7)

- DONE: Inspect the retained cycle-6 public stream, correlated rollout, final message, stderr, and durable workflow root to identify why zero completed CommandExecution items were observed.
  The original `t.TempDir` artifacts were gone, so the First Officer authorized one unchanged diagnostic rerun; durable bytes are under `/tmp/spacedock-6k-cycle7-live-artifacts.dpYyUp`, and the copied workflow root is `retained-go-temp/002`.
- DONE: Classify the failure as model behavior, sandbox/write posture, correlation/extraction mismatch, or another exact cause using retained bytes; do not edit candidate bytes or rerun the model.
  Sandbox/write posture: the turn completed normally, but the model reported that its first contract read hit the environment path sandbox and that no command or contract read was permitted; this was one authorized diagnostic rerun on unchanged 354f4e362.
- DONE: Report the smallest evidence-backed next action and whether the live invocation produced any entity, tool call, refusal, error, or alternate event schema.
  Next action is a narrowly authorized filing-sandbox access correction that makes the first-officer contract path readable while keeping the rollout outside the writable workflow; then retain the isolated rollout and rerun once. This invocation produced no entity, tool event, refusal, or alternate schema—only a sandbox-error final message.
- DONE: Preserve exact public stream, stderr, final message, process result, and filesystem outcome.
  `codex-shared-scenarios/filing/` retains `codex-exec.jsonl`, `codex-exec.stderr.txt`, `codex-final-message.txt`, and `codex-process-result.txt`; SHA-256 values are b3273b9, a7ef0f9, be2c8b7, and 8fc3006 respectively.
- DONE: Preserve relevant event types, status, and output.
  Public JSONL has exactly `thread.started`, `turn.started`, three `item.completed/agent_message`, and `turn.completed`; process result is terminal=true, timed_out=false, duration=51.058505291s, while the harness reports exit_code=-1 after the front door remained alive past terminal.
- DONE: Verify the durable workflow outcome.
  `retained-go-temp/002` is clean at commit 886b79b, tracks only `README.md`, and has no `wire-the-thing.md`; the final message says no files changed, no workers dispatched, and no advancement.
- SKIPPED: Preserve the correlated rollout file itself.
  The harness resolved the sole rollout for thread `01a000e2-5fb2-7230-a0e0-5f7968d14be0` and decoded zero completed `CommandExecution` items, but its isolated `_codex-home` was removed by registered test cleanup before the external copier captured it; the public bytes contain no alternate command/tool item.

### Summary

Root cause is the filing-only workspace sandbox denying the model's first contract read, after which the model terminated without attempting the atomic filing command. This is not a correlation or extraction mismatch: the retained public stream contains only agent messages and a normal terminal event, and the retained workflow proves no entity landed. Candidate 354f4e362 remains byte-identical and no unapproved model run or workflow mutation occurred.

## Stage Report: implementation (cycle 6)

- DONE: Add only the isolated CODEX_HOME/plugins subtree to the filing journey's sandbox roots so the first-officer contract is readable; keep CODEX_HOME/sessions outside every added/writable root.
  Commit 470994d23 adds one top-level `--add-dir <isolated-CODEX_HOME>/plugins`; it never adds CODEX_HOME, sessions, danger-full-access, or bypass approvals, so the correlated rollout stays outside writable roots.
- DONE: Strengthen argv tests to prove the exact plugin-root add-dir placement, reject adding CODEX_HOME or sessions, and preserve byte-identical non-filing argv.
  `TestCodexLiveRunnerUsesRestrictedPostureOnlyForFiling` fails unless the exact host prefix is workspace-write/never plus one plugins-only add-dir before `exec`, and still compares non-filing argv byte-for-byte.
- DONE: Run focused sandbox/argv and transaction matrices, live-tag compile, gofmt, full, race, and exact five-path/LOC checks; do not run the model-backed filing target.
  Focused sandbox/argv, transaction, extraction, and launcher-order tests, live-tag compile, `gofmt -w ./cmd ./internal`, full, race, diff, and surface checks passed; final delta is 175 insertions/190 deletions (net -015) across 5 paths, with no model run.

### Summary

The filing sandbox now exposes only the installed plugin subtree needed for the first contract read while preserving the runtime-owned sessions rollout outside model-writable roots. Transaction semantics, product bytes, ACs, and non-filing argv are unchanged; validation owns the durable-artifact live rerun.

## Stage Report: validation (cycle 8)

- FAILED: Reproduce AC-1 with the seven-rung correlated transaction matrix, then spend the single exact Codex filing run on frozen commit e851be3fe and verify one owned atomic command, exact receipt, correct durable entity, and zero unreachable/unrelated evidence.
  On updated frozen commit 470994d23, rungs 1-3 passed, but the exact live target again produced zero `CommandExecution` items; nine runtime `custom_tool_call/exec` attempts all failed before execution with `sandbox-exec: sandbox_apply: Operation not permitted`.
- DONE: Reproduce AC-2 by attacking duplicate/mixed creates, any --next-id regardless of exit, failed/wrong/missing commands, forged public evidence, receipt ownership/bytes, malformed or ambiguous rollout, and entity mismatch.
  The focused matrix/extraction/correlation controls and detached audit at `/tmp/spacedock-codex-filing-cycle8-audit.RR2Kpy/checkout` passed; weakening any transaction, receipt, entity, or rollout invariant makes a named case pass.
- FAILED: Verify AC-3 with exact five-path -016 LOC evidence, filing-only workspace-write/never and non-filing argv parity, live-tag compile, gofmt, full, race, and a detached adversarial audit; recommend PASSED or REJECTED without editing candidate bytes.
  Updated approved surface is exactly 5 test paths and net -015; filing root containment, non-filing parity, live-tag compile, gofmt, diff, and detached audit passed, but the new live finding stopped full/race before rerun.
- FAILED: Reproduce AC-1 evidence.
  The retained rollout for thread `01a000f2-e990-77e0-a2c1-b14bf71385b7` contains zero `CommandExecution`, nine completed custom exec calls, and repeated sandbox-apply failures; no atomic command, receipt, or entity exists.
- DONE: Reproduce AC-2 evidence.
  Rungs 4-7 reject duplicate/mixed creates, every `--next-id` exit, failed/wrong/manual/missing commands, forged public evidence, wrong-item/inexact receipts, ambiguous rollout, and durable mismatch.
- FAILED: Reproduce AC-3 evidence.
  Exactly one `--add-dir <CODEX_HOME>/plugins` is present; detached containment proves `<CODEX_HOME>/sessions` is a sibling outside that and all other added roots, while non-filing argv is byte-identical. The supported host still cannot apply its nested workspace sandbox.
- DONE: Preserve public stream, stderr, final message, process result, isolated rollout, and workflow outcome before cleanup.
  Durable evidence is under `/tmp/spacedock-6k-cycle8-live-artifacts.qxb445`; rollout SHA-256 is d2f70fe9, public stream is 041a865e, and the clean copied workflow `retained-go-temp/002` tracks only `README.md` with no `wire-the-thing.md`.
- DONE: Stop before candidate mutation, full/race reruns, or another live target after the new finding.
  Candidate HEAD remains 470994d23 and the exact model-backed target was invoked once on these frozen bytes.

### Review-finding disposition

- Released user and normal workflow: the required Codex filing journey under the normal outer agent-safehouse, with filing-only top-level workspace-write/never and plugins-only added root.
- Observable harm: Codex cannot apply its nested execution sandbox, so every shell attempt fails before command execution and the promised filing transaction/entity never occurs.
- Authority: value-ac[AC-1] requires one completed exit-0 atomic filing transaction, exact same-item receipt, and correct durable entity from the exact supported live run.
- Trigger evidence: retained rollout lines 35-50 show nine `custom_tool_call/exec` attempts and `sandbox_apply: Operation not permitted`, including `pwd` and `/bin/ls`; public/final bytes report every command rejected, and the retained workflow has no entity.
- Proposed classification: Material outcome defect and chosen-mechanism failure; the plugins-only root is not causal because nested `sandbox-exec` cannot apply inside the outer agent-safehouse. Recommend scope/design reset and `Needs decision`, since bypass would violate rollout authority and changing the outer host/sandbox is outside this candidate.

### Summary

REJECTED. The plugins-only correction is narrow, keeps sessions outside every writable/add-dir root, preserves non-filing argv, and passes all offline and detached falsifiers, but the live host cannot nest Codex workspace sandboxing inside the existing agent-safehouse. This is a mechanism-level incompatibility rather than another missing path; candidate bytes remain unchanged and no additional validation runs followed the finding.

## Stage Report: ideation (cycle 5)

- DONE: Re-anchor the observer on harness-owned public item.completed/command_execution events and prove why the model cannot forge or inject those event envelopes.
  `runCodexProcess` captures the front door's stdout directly; Codex serializes command-controlled output inside `aggregated_output`, and the proposed decoder reads each top-level public JSONL line once without recursive decoding.
- DONE: Design one bounded transaction grade that handles PR #679's exact lossy command display, same-item exit/receipt, durable entity, duplicate/mixed create, any --next-id, and unreachable native-text counterfeits without rollout authority or nested sandboxing.
  AC-2 evidence: the seven-rung fail-closed controls bind exact create count, completed exit 0, receipt, and entity in one public item; the retained quote-seam fixture passed `TestAssertCodexFilingViaNew`, while duplicate/mixed creates, any `--next-id`, detached receipts, command-printed events, and nested-output evidence are specified to fail.
- DONE: Update ACs, seven-rung offline fixtures, semantics, exact paths, and signed net LOC estimate from frozen 470994d23; forbid another live run until the exact retained PR #679 fixture and all detached counterfeits pass.
  AC-1 through AC-3 now name the public-envelope value proof and boundaries; expected surface is six test files, about 115 insertions/170 deletions, net -055 with accepted range -090..-020, and validation alone owns the next exact live run.

### Summary

The selected design grades one harness-owned public command-execution item instead of trusting a rollout, ledger, or nested sandbox. It preserves PR #679's exact lossy display as the positive fixture and requires every detached counterfeit to pass offline before the single validation live run.
