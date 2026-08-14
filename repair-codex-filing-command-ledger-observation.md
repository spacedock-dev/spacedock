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

The supported Codex observables stop one process too early. Public `codex exec --json` is a harness-owned pipe but exposes a lossy shell display. The correlated Codex session rollout exposes a runtime-generated `CommandExecution` with exact outer argv, status, and exit code, but the executable is `/bin/zsh -lc`; it does not expose child `execve` events. The rollout file is also model-writable under the filing journey's current sandbox-bypass posture. Durable entity bytes, stdout, and Git history cannot supply the missing child-process identity. Therefore no existing supported Codex surface proves real `spacedock` executable entry against the required counterexamples.

## Risk spike

No model run was spent. Codex CLI `0.147.0` reports hooks as stable and persists runtime-owned `event_msg.item_completed` / `CommandExecution` records. A no-model spike against this session observed exact outer argv and exit status, then exercised the three detached shapes:

- `if false; then /nonexistent/spacedock new wire-the-thing; fi` produced one completed `/bin/zsh -lc` event with exit 0 despite zero child entry.
- `/usr/bin/true; /usr/bin/true` produced one completed outer event, so the observable cannot count duplicate child executions.
- `/usr/bin/false; /usr/bin/true` produced one completed outer event with exit 0, so a later success masks the earlier child failure.

Codex hooks observe tool use at the same shell boundary and do not add child-process identity. Reading the session rollout from a workspace-write sandbox would remove model write authority but would not repair these three non-injective results. The reserved exact `TestLiveCommonFiling/codex` run remains unspent.

## Proposed approach

Implementation is held for a captain scope decision; another parser, file ledger, or shim-only patch cannot satisfy AC-1 and AC-2.

**Option A — preserve the end value (recommended):** authorize a new test-only child-execution authority outside the model's write capability. The smallest credible design is a filing-only process supervisor that receives records over a private Unix socket and accepts them only after verifying the peer process is the exact harness-created wrapper executable; it then counts every create and every `--next-id` attempt before forwarding to the real binary. This is custom OS/runtime instrumentation, requires a platform spike and platform-specific peer-identity support, and is outside the currently approved “existing supported observable” boundary.

**Option B — preserve the current scope:** change AC-1/AC-2 from real `spacedock` executable entry to Codex outer-command execution. Run only the filing journey in `workspace-write` with `never` approval so the correlated session rollout is not model-writable, then grade its runtime-owned `CommandExecution` events. This is cheaper and uses a supported observable, but explicitly permits unreachable, duplicate, and masked-failure child commands; it abandons the stated end value and the three required falsifiers.

**Option C — product receipt:** authorize a product-owned, harness-verifiable execution receipt from `spacedock new`. This would be portable but changes product command/runtime surface and is currently forbidden.

No option is selected in this ideation. The captain must either authorize Option A/C's expanded authority or explicitly relax AC-1/AC-2 to Option B. Until then there is no sound implementation target.

## Verification ladder

After the scope decision, the first test remains the exact falsifier matrix:

1. One real exit-0 `new wire-the-thing` entry plus the landed entity passes.
2. Hand-written valid old-ledger bytes, a forged public receipt, narration, and unreachable JavaScript or shell text produce zero accepted entries.
3. Two successful creates fail the exactly-one rule, including concurrent creates.
4. Any executed `status --next-id` fails regardless of its exit code or a successful create beside it.
5. Failed create, wrong slug, missing create, malformed observer input, and manual entity write fail.
6. The model cannot write, replace, truncate, or invoke the observation authority directly; the negative test exercises that capability boundary rather than checking pathname placement.

Focused tests, live-tag compile, `gofmt`, full, and race run before validation. Validation alone spends the one exact Codex filing model run and requires one accepted real create, the exact entity on disk, and zero `--next-id` attempts.

## Out of scope

- No product hook, command, protocol, state store, lifecycle guard, global hook, or replacement for the compaction `SessionStart` hook; the ledger exists only under `t.TempDir()` for the Codex filing journey.
- No change to `spacedock new`, command grammar, stored entity formats, write authority, or product runtime behavior.
- No work on Pi, Opus, rejection flow, supporting-evidence, mechanically-continue-Codex validation, another owner's target binding, or reconciliation row.
- PR #682's rejection-flow failure remains `live-evidence-followups` ownership.

## Acceptance criteria

- **AC-1 (VALUE):** A real exit-0 `spacedock new wire-the-thing` child execution produces exactly one accepted observation and the expected on-disk entity, while hand-written evidence and unreachable JavaScript or shell commands produce none. Verified by matrix rungs 1-3 and the sole exact Codex filing live run; substituting outer-shell occurrence or model-writable bytes makes the negative controls fail.
- **AC-2:** Bound and PATH-resolved aliases pass only on real child execution; duplicate creates, any executed `--next-id` regardless of exit, manual/non-atomic, failed, wrong-slug, missing, malformed, forged, and concurrent cases remain failures. Verified by matrix rungs 3-6 through the same observer and grade.
- **AC-3:** The selected implementation changes only Codex filing-test observation within the captain-approved expanded authority; product grammar, stored formats, write authority, other runtimes, and non-filing journeys remain unchanged unless the captain explicitly selects Option C. Verified by an exact changed-path/semantic-budget check, live-tag compile, `go test ./...`, and `go test ./... -race`.

## Test plan

- Before implementation, spike the captain-selected authority against peer forgery/model write access and all six verification rungs. A failed capability-boundary spike returns to ideation; it does not trigger another parser or file-ledger correction.
- Add table-driven observer/grade tests that exercise real child processes and the three exact validation counterexamples.
- Retain the correlated PR #679 fixture only as the historical false-negative shape; it is not execution authority.
- Run the exact targeted Codex filing live test once on frozen final bytes in validation. Do not spend another live correction on trace parsing.
- Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` before delivery.

## Expected surface and semantic budget

No implementation estimate is authoritative until the captain selects the required scope; the alternatives have materially different authority and platform cost.

- **Option A provisional surface from frozen `4faa23656`: +180 net LOC across 6 test files** (about 260 insertions and 80 deletions): `internal/ensigncycle/codex_filing_execution_observer_test.go` (replacement supervisor and peer check), `codex_filing_peer_darwin_test.go`, `codex_filing_peer_linux_test.go`, `codex_live_runner_test.go`, `shared_filing_negative_test.go`, and `team_capability_test.go`. Tolerance: ±60 after the required platform spike. Allowed semantic change: filing-only test process supervision and peer identity; no product or other-journey behavior.
- **Option B provisional surface from frozen `4faa23656`: -050 net LOC across 5 test files** (about 120 insertions and 170 deletions): delete `codex_filing_invocation_ledger_test.go`; update `codex_live_runner_test.go`, `claude_runtime_helpers_test.go`, `shared_filing_negative_test.go`, and `team_capability_test.go`. Tolerance: ±35. Allowed semantic change: filing-only workspace-write sandbox and outer-command grading, plus explicit relaxation of AC-1/AC-2.
- **Option C:** intentionally unestimated until product receipt authority, command surface, and portability are approved; it breaches the current zero-product-file semantic budget.

For Options A/B, CLI output, command grammar, stored formats, write authority, Claude/Pi grading, other journeys, and site documentation remain unchanged. The current frozen candidate remains exactly `4faa23656` and the exact Codex filing model run remains reserved.

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
