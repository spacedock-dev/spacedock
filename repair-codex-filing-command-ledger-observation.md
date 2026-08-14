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

This is an observation defect in the Codex live-test adapter, not a product filing defect. Candidate `77a71f8e03cc2f2280ad69d47be58c36c4c28a74` is frozen after validation proved that parsing native JavaScript is not execution evidence: an unrelated successful executed call plus an unreachable dot-form `tools.exec_command({cmd:"spacedock new wire-the-thing"})` passed the ledger.

The complete existing chain is non-injective. The correlated parent rollout proves which `functions.exec` payload Codex submitted; it does not prove which branch inside arbitrary JavaScript ran. The public `item.completed` proves one rendered shell command completed and carries its exit status/output, but its display is deliberately lossy. The final entity proves that some writer produced valid bytes, and the fixture Git history proves only the initial checkout plus later commits (the `new` operation itself does not commit). Even together, those observables cannot distinguish the detached counterfeit from a real atomic invocation without trusting either a JavaScript parser, a shell-display parser, or forgeable stdout. The cheapest execution-grounded discriminator is therefore a test-local launcher shim that records the real `spacedock` argv and exit code before returning control to the host.

## Risk spike

No new spike is needed and no live run was spent. Repository history already exercises the riskiest mechanism: commit `603af41123d1509d130d952eb169d76c7215e8f0` contains a test-local NUL-delimited argv ledger whose negative control echoes command-shaped text without producing a record, and commit `c293037ba38ab0b15e385116c1e57ec103f64b81` contains a Codex front-door probe showing the shim remains bound after the launcher repins `SPACEDOCK_BIN`. These are mechanism precedents only; implementation will not absorb their boot, Claude, Pi, or other-owner changes.

The detached counterfeit is the first implementation test: its unrelated successful public command and unreachable native dot-form filing text must produce zero successful launcher-ledger entries. It fails the frozen candidate and passes only when observation is grounded at executable entry.

## Proposed approach

For the Codex `filing` journey only, bind both `SPACEDOCK_BIN` and the front of `PATH` to a test-local `spacedock` shim. The shim writes one atomic, NUL-delimited record containing exact argv and the real binary's exit code, forwards stdin/stdout/stderr unchanged, and exits with the real code. The Codex host shim reapplies this binding after the Spacedock front door repins `SPACEDOCK_BIN`. After the run, grade only completed exit-0 launcher records: `new <slug>` or `status --new <slug>` passes; any executed `status --next-id`, wrong slug, missing invocation, or nonzero create fails. The existing final-file assertion remains the independent result check.

This execution ledger serves AC-1 and AC-2. The simplest alternative is the existing public completion plus exact `created:` receipt and on-disk entity; it is insufficient because a shell command can print that receipt and hand-write identical bytes. The native/public correlation alternative is insufficient because the detached counterfeit proves call occurrence is not call execution. A full shell parser is larger and still cannot prove which conditional branch ran. The shim is the smallest observable at the boundary where execution becomes real.

Product bytes remain unchanged, and no product hook, protocol, or stored format is introduced. Runtime execution bytes cannot remain literally untouched: without one test-only executable boundary, the available traces cannot prove actual invocation. The allowed semantic change is therefore narrowly expanded from passive Codex trace observation to Codex filing-test launcher instrumentation; real `spacedock` argv, stdin, stdout, stderr, exit code, CLI grammar, and filesystem behavior remain unchanged.

## Verification ladder

The focused falsifier matrix is exact and does not spend a model run:

1. **Bound positive:** `"$SPACEDOCK_BIN" new wire-the-thing` reaches the shim, the real stub exits 0, exact argv is recorded, and grading passes.
2. **PATH positive:** direct `spacedock status --new wire-the-thing` reaches the same shim with exit 0 and passes.
3. **Detached counterfeit:** unrelated executed `echo` plus unreachable dot-form filing text, a forged public exit-0 completion, and even a landed entity produce no launcher record and fail.
4. **Manual/non-atomic:** `status --next-id` plus a direct file write records no successful create and fails; `--next-id` also fails when a successful create is present.
5. **Failed create:** exact `new wire-the-thing` with real exit 1 fails even if the file or public receipt is fabricated afterward.
6. **Wrong/missing create:** exit-0 `new other-slug`, narration only, and no ledger file fail.
7. **Ledger integrity:** truncated, duplicate-terminal, unknown-tool, or malformed NUL records fail closed; concurrent records are independent atomic files.
8. **Front-door binding:** a live-tagged, no-model shell probe confirms the Codex launcher rebinds `SPACEDOCK_BIN` to the ledger shim after front-door pinning.

Implementation runs rungs 1-8, focused filing tests, live-tag compile, `gofmt`, full, and race. Final validation alone spends one exact `TestLiveCommonFiling/codex` run on frozen bytes; it must show one exit-0 argv record for `new wire-the-thing`, the exact entity on disk, and no `--next-id` record. Replacing real execution with the detached counterfeit makes AC-1 fail; accepting any rung 4-7 makes AC-2 fail.

## Out of scope

- No product hook, command, protocol, state store, lifecycle guard, global hook, or replacement for the compaction `SessionStart` hook; the ledger exists only under `t.TempDir()` for the Codex filing journey.
- No change to `spacedock new`, command grammar, stored entity formats, write authority, or product runtime behavior.
- No work on Pi, Opus, rejection flow, supporting-evidence, mechanically-continue-Codex validation, another owner's target binding, or reconciliation row.
- PR #682's rejection-flow failure remains `live-evidence-followups` ownership.

## Acceptance criteria

- **AC-1 (VALUE):** A real exit-0 `spacedock new wire-the-thing` execution produces exactly one successful argv-ledger observation and the expected on-disk entity, while the detached unrelated-command/unreachable-call counterfeit produces none. Verified by focused matrix rungs 1-3 and the sole exact Codex filing live run; replacing executable entry with native JavaScript occurrence makes the counterfeit pass and the test fail.
- **AC-2:** Bound and PATH-resolved aliases pass only on real exit-0 execution; manual/non-atomic, failed, wrong-slug, missing, malformed-ledger, duplicate-terminal, and `--next-id` cases remain failures. Verified by focused matrix rungs 2 and 4-8 through the same ledger reader and filing grade.
- **AC-3:** The delivered diff changes only Codex filing-test observation: product bytes, CLI grammar, stored formats, authority, other runtimes, and non-filing journeys are unchanged. Verified by a base-to-candidate changed-path check limited to the named `_test.go` surface, plus live-tag compile, `go test ./...`, and `go test ./... -race`.

## Test plan

- Add a default-tag Codex filing ledger helper and table-driven tests for the eight rungs above. Tests execute a local stub through the shim, so they prove argv/exit recording rather than inspect prose or arbitrary JavaScript.
- Add a live-tagged no-model front-door binding test. It fails if launcher pinning bypasses either `SPACEDOCK_BIN` or PATH instrumentation.
- Remove the PR #679 native/public parser fixtures and structural JavaScript decoder; retain the detached counterfeit as the negative control expressed against the executable ledger.
- Run the exact targeted local Codex filing live test once on the final candidate as the single revalidation. Do not spend another live correction on trace parsing.
- Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` before delivery.

## Expected surface and semantic budget

Expected implementation delta from frozen `77a71f8e0`: **-034 net LOC across 7 files** (about 190 insertions and 224 deletions):

- `internal/ensigncycle/codex_filing_invocation_ledger_test.go` — new Codex-only argv/exit ledger and fail-closed reader.
- `internal/ensigncycle/codex_live_runner_test.go` — bind the filing ledger across front-door pinning and grade its successful records.
- `internal/ensigncycle/shared_filing_negative_test.go` — replace native/parser ladders with the eight execution-grounded rungs.
- `internal/ensigncycle/shared_filing_test.go` — remove correlated-input and JavaScript decoder additions; retain shared command grammar.
- `internal/ensigncycle/claude_runtime_helpers_test.go` — restore the pre-task native-lifecycle helper boundary.
- `internal/ensigncycle/testdata/codex_filing_pr679/public.jsonl` — remove parser-only fixture.
- `internal/ensigncycle/testdata/codex_filing_pr679/parent-rollout.jsonl` — remove parser-only fixture.

Tolerance: ±40 net LOC within these same seven paths; zero non-test product files and no borrowed boot/Claude/Pi ledger surface. Relative to the original task base, the expected delivered surface is about **+190 net LOC across 3 test files** because the four parser-only candidate paths return to base. Allowed semantic change: the Codex filing live test routes launcher execution through a test-local transparent shim and grades exact successful argv. No allowed change to CLI output, command grammar, stored formats, write authority, product runtime behavior, Claude/Pi grading, other journeys, or documentation; no site doc diff is needed because user-visible behavior does not change.

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
