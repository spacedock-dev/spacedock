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

This is an observation defect in the Codex live-test adapter, not a product filing defect. Candidate `4faa23656` is frozen after two observation mechanisms failed. Native JavaScript occurrence could pair an unrelated successful execution with an unreachable dot-form filing call. The replacement launcher shim reached executable entry, but wrote its authority to a model-writable directory; a hand-written valid record passed without entry, duplicate successful creates passed, and a failed `--next-id` beside a successful create was ignored.

The captain has re-anchored AC-1 from OS child-`execve` identity to proof that Codex executed the atomic filing path. Public `codex exec --json` remains too lossy by itself, but the already-supported thread correlation resolves its `thread.started` ID to exactly one Codex session rollout. That rollout contains runtime-generated `CommandExecution` items with the executed outer command, completion status, exit code, stdout, and stderr. Combined with the exact `created:` receipt and independently verified landed entity, this supported evidence distinguishes an executed atomic path from unreachable native JavaScript text and unrelated success without requiring OS-specific supervision.

## Risk spike

No model run was spent. Codex CLI `0.147.0` reports hooks as stable and persists runtime-owned `event_msg.item_completed` / `CommandExecution` records. A no-model spike against this session observed exact outer argv and exit status, then exercised the three detached shapes:

- `if false; then /nonexistent/spacedock new wire-the-thing; fi` produced one completed `/bin/zsh -lc` event with exit 0 despite zero child entry.
- `/usr/bin/true; /usr/bin/true` produced one completed outer event, so the observable cannot count duplicate child executions.
- `/usr/bin/false; /usr/bin/true` produced one completed outer event with exit 0, so a later success masks the earlier child failure.

The spike establishes the boundary of the selected claim: correlated `CommandExecution` can prove which outer command Codex actually ran, but not child `execve` identity or arbitrary shell branch execution. The revised grade therefore requires one unambiguous atomic-`new` filing observation, rejects mixed `--next-id` or duplicate create evidence within or across observed commands, and requires the exact successful receipt plus durable entity. A filing-only workspace-write/never launch removes model write access to the correlated session log. The reserved exact `TestLiveCommonFiling/codex` run remains unspent.

## Proposed approach

For the Codex `filing` journey only, replace the file ledger with the existing `codexNativeLifecycleStream` correlation: take the sole `thread.started` ID from the harness-owned public stream, resolve exactly one rollout under the isolated `CODEX_HOME`, and inspect only runtime-generated completed `CommandExecution` items. Launch only this journey with Codex `workspace-write` and approval `never`, keeping the rollout outside the writable workflow root.

Grade the correlated execution evidence as one transaction:

1. Exactly one completed exit-0 observed command must contain exactly one atomic `new wire-the-thing` / `status --new wire-the-thing` invocation recognized by the existing `commandFilesViaNew` grammar.
2. No observed command, successful or failed, may contain `--next-id`; no command may contain a second atomic create for the requested slug. Counts apply within one shell command and across commands.
3. The same completed command must carry exactly one receipt line `created: <canonical-entity-path> id=<landed-id>` in runtime-recorded stdout. A receipt from an unrelated command, a forged public event, or command narration does not correlate.
4. The existing filing assertion independently requires the canonical entity path, expected title/status/body, stamped ID, and clean durable workflow result.

This serves AC-1 and AC-2 without a new JavaScript or shell parser: it structurally decodes the supported Codex event schema and reuses the existing bounded filing recognizer. Unreachable dot-form `tools.exec_command` text has no `CommandExecution` item and fails. Unrelated success cannot lend status or receipt to another item. Mixed `new`/`--next-id`, duplicate creates, failed commands, wrong slug, missing or duplicate receipt, and wrong entity fail closed.

The simplest alternative is the frozen launcher ledger; it is insufficient because the model can write its records. Public completion plus durable state is also insufficient because public command display is lossy and stdout/state can be forged independently. Platform process supervision was rejected by the captain as unnecessary for the revised atomic-path claim.

## Verification ladder

The focused no-model matrix exercises the selected transaction:

1. One correlated completed exit-0 command with one recognized atomic create, one exact receipt, and the correct entity passes.
2. PR #679's native command shape passes even though the public display alone remains unrecognized.
3. Unreachable native dot-form filing text plus an unrelated successful `CommandExecution` produces zero creates and fails, even with a forged public receipt and landed entity.
4. Duplicate create evidence within one command or across two completed commands fails the exactly-one count.
5. Any `--next-id` evidence fails regardless of command exit or a valid create/receipt/entity beside it.
6. Failed create, wrong slug, missing create, receipt from another command, missing/duplicate/wrong-path/wrong-ID receipt, malformed/missing/ambiguous rollout, and wrong entity fail.
7. Filing-only argv proves `workspace-write` plus `never`; other journey argv remains byte-identical.

Focused tests, live-tag compile, `gofmt`, full, and race run before validation. Validation alone spends the one exact Codex filing model run and requires one accepted atomic-path observation, the exact receipt and entity, and zero `--next-id` evidence.

## Out of scope

- No product hook, command, protocol, state store, lifecycle guard, global hook, process supervisor, file ledger, or replacement for the compaction `SessionStart` hook.
- No change to `spacedock new`, command grammar, stored entity formats, write authority, or product runtime behavior.
- No work on Pi, Opus, rejection flow, supporting-evidence, mechanically-continue-Codex validation, another owner's target binding, or reconciliation row.
- PR #682's rejection-flow failure remains `live-evidence-followups` ownership.

## Acceptance criteria

- **AC-1 (VALUE):** Supported Codex evidence establishes exactly one completed exit-0 atomic filing-path observation, its exact `created: <canonical-path> id=<landed-id>` receipt, and the correct durable entity; unreachable native filing text plus unrelated success establishes zero. Verified by matrix rungs 1-3 and the sole exact Codex filing live run; removing thread correlation, receipt correlation, or the landed-entity check makes a named negative control pass.
- **AC-2:** Bound and PATH-resolved atomic aliases pass only with one correlated transaction; duplicate create evidence, any observed `--next-id` regardless of exit, manual/non-atomic, failed, wrong-slug, missing, malformed, forged-public, mixed-command, wrong-receipt, and ambiguous-rollout cases remain failures. Verified by matrix rungs 4-7 through the same observer and grade.
- **AC-3:** The delivered diff changes only Codex filing-test observation and its filing-only sandbox posture; product bytes, CLI grammar, stored formats, authority, other runtimes, and non-filing journeys are unchanged. Verified by exact changed paths and signed LOC, filing/other-journey argv controls, live-tag compile, `go test ./...`, and `go test ./... -race`.

## Test plan

- Replace the ledger matrix with table-driven correlated `CommandExecution` transaction tests for all seven rungs. Fixtures structurally encode supported Codex events; they do not parse native JavaScript or infer child processes.
- Extend the existing correlated-session fixture with PR #679 command/stdout evidence and the unreachable/unrelated-success control. Mutating thread ID, event type/status/exit, create count, receipt ownership, receipt bytes, or entity bytes must turn the corresponding positive red.
- Add filing-only argv tests for workspace-write/never and a non-filing parity control.
- Run the exact targeted Codex filing live test once on frozen final bytes in validation. Do not spend another live correction on trace parsing.
- Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` before delivery.

## Expected surface and semantic budget

Expected implementation delta from frozen `4faa23656`: **-050 net LOC across 5 test files** (about 150 insertions and 200 deletions), accepted range **-085 to -015**:

- `internal/ensigncycle/codex_filing_invocation_ledger_test.go` — delete the model-writable file ledger and reader.
- `internal/ensigncycle/codex_live_runner_test.go` — use correlated native execution evidence and receipt/entity grading; remove launcher/host shim binding.
- `internal/ensigncycle/claude_runtime_helpers_test.go` — extend the existing thread-to-rollout correlation with structural `CommandExecution` transaction extraction and exact receipt ownership.
- `internal/ensigncycle/shared_filing_negative_test.go` — replace ledger integrity tests with the seven-rung correlated evidence matrix.
- `internal/ensigncycle/team_capability_test.go` — add filing-only workspace-write/never argv plus non-filing parity controls.

Allowed semantic change: Codex filing-test evidence moves from a model-writable launcher ledger to correlated runtime `CommandExecution` plus exact receipt and durable entity; only the filing journey changes from sandbox bypass to workspace-write/never. No allowed change to CLI output, command grammar, stored formats, write authority, product runtime behavior, Claude/Pi grading, other journeys, or documentation. The candidate remains frozen at `4faa23656`; the exact Codex filing model run remains reserved for validation.

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
