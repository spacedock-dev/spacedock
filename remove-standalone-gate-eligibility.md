---
title: Remove the standalone gate eligibility ceremony
status: validation
source: "Captain simplification review, 2026-08-03: eligibility is a necessary atomic safety predicate but an unnecessary operator preflight and readiness vocabulary."
score: 0.9
sprint: durable-decisions
group: gate-operator-ux
milestone: 0.27.0
sprint-readiness: ready
id: bv3hhbqr5spt1wn4557qyp8c
gates:
    version: 1
    records:
        - id: gate:bv3hhbqr5spt1wn4557qyp8c:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:bv3hhbqr5spt1wn4557qyp8c-backlog-1
              briefing:
                id: briefing:bv3hhbqr5spt1wn4557qyp8c:backlog:attempt-1:revision-1
                digest: sha256:4b936d38cb7ea50d8823bb73c50ef257b510d64736bf786950a32ff37c97ed58
                request-digest: sha256:098a9a069ace77b1c9afdc1e0166bb33788841d7d2ea8f37e8b14f43f41def03
                room-ref: ./remove-standalone-gate-eligibility/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:bv3hhbqr5spt1wn4557qyp8c:backlog:1
                briefing: briefing:bv3hhbqr5spt1wn4557qyp8c:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-03T23:30:19.040582Z"
                decision: approve
                reason: Captain directed dispatch of the next ideation; latest recorded SO concurrence places bv after g3 and preserves the private fail-closed predicate while removing only the public eligibility ceremony.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:bv3hhbqr5spt1wn4557qyp8c:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:bv3hhbqr5spt1wn4557qyp8c-ideation-1
              briefing:
                id: briefing:bv3hhbqr5spt1wn4557qyp8c:ideation:attempt-1:revision-1
                digest: sha256:51ad1e4dc1e6df9e6b757e77cc9644037d18afbf08ffe5ab56dc9f51432fb9ea
                request-digest: sha256:979bc8d6359e29a5c0c77a5ea03430b7379976d34f8d219e91a1f874e06a5bb6
                room-ref: ./remove-standalone-gate-eligibility/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:bv3hhbqr5spt1wn4557qyp8c:ideation:1
                briefing: briefing:bv3hhbqr5spt1wn4557qyp8c:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-03T23:53:15.718694Z"
                decision: approve
                reason: 'Science Officer concurs: approve the ideation design now; keep its implementation advance held until G3 completes. The Captain directed the next ideation dispatch.'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:bv3hhbqr5spt1wn4557qyp8c:validation
          stage: validation
          attempts:
            - id: gate-attempt:bv3hhbqr5spt1wn4557qyp8c-validation-1
              briefing:
                id: briefing:bv3hhbqr5spt1wn4557qyp8c:validation:attempt-1:revision-1
                digest: sha256:7fff755c501e0addc8182fe50c6756d1fbde859b21d233ca0a280b20c51839bd
                request-digest: sha256:19a6903b1b25bc0536c0881a37ae88a7ee87d77e5a5cfb45b2dbb4bdfa65eef0
                room-ref: ./remove-standalone-gate-eligibility/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:bv3hhbqr5spt1wn4557qyp8c:validation:1
                briefing: briefing:bv3hhbqr5spt1wn4557qyp8c:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-06T15:41:00.260459Z"
                decision: approve
                reason: Captain sprint conn authorizes gate approval. All four value criteria and both detached guard mutations pass. Exact candidate 013c8729e removes the public ceremony and retains locked authority checks. PR CI remains required because the local full and race suites are not green.
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:bv3hhbqr5spt1wn4557qyp8c-validation-2
              briefing:
                id: briefing:bv3hhbqr5spt1wn4557qyp8c:validation:attempt-2:revision-1
                digest: sha256:e99ac168d6cd8eccfb3e63dfa7efa03a95190334eacb9d9aa66beb2bcd815dcc
                request-digest: sha256:95c2238a79bedf0475363e3625e3ce0465a55650cdbc0d62adf6631892cfe6fa
                room-ref: ./remove-standalone-gate-eligibility/review/validation/briefing-2
started: 2026-08-03T23:33:06Z
worktree: .worktrees/spacedock-ensign-remove-standalone-gate-eligibility
mod-block:
---

## Outcome

A First Officer uses `status` to see the next action and uses `gate consume` to apply an approval. The operator does not run a separate eligibility ceremony.

`gate consume` and terminal delivery must retain the internal fail-closed predicate. This predicate prevents stale approval, wrong-stage advancement, wrong-target advancement, and duplicate spending.

## Problem

The public `gate eligibility` command exposes a read-only preview of the predicate that `gate consume` recomputes under the entity lock. The normal First Officer path does not require this preview. The skill calls it an optional diagnostic.

The extra command creates another operator step and another set of condition labels. These labels can look like scheduler state even though they are derived diagnostics.

The safety mechanism and the command are different things. Removing the command must not remove or weaken the mutation guard.

## Proposed direction

Remove `gate eligibility` from the public CLI, help, command reference, completion, and First Officer procedure. Let `status` project the next action from durable facts.

Keep one internal predicate at each authority-spending boundary. `gate consume` must evaluate it while it holds the entity lock. Terminal delivery and rework must use the same authority checks.

On refusal, the acting command must report the concrete reason. Do not expose the predicate as a separate workflow phase or stored state class.

This is an unreleased-v1 cut. Do not add compatibility aliases, deprecation machinery, a replacement preflight command, or a new readiness class.

## Acceptance criteria

**AC-1 (VALUE) — The normal gate journey has no eligibility preflight.**
Verified by: a real CLI lifecycle runs `status`, close-and-consume, successor dispatch or terminal routing, and no `gate eligibility` command.

**AC-2 — Authority spending remains fail-closed and atomic.**
Verified by: stale input, wrong stage, wrong target, consumed authority, superseded authority, and malformed binding controls refuse without an unauthorized status transition.

**AC-3 — The operator receives an actionable refusal from the acting command.**
Verified by: `gate consume` and terminal delivery controls report the applicable refusal reason without requiring a prior diagnostic command.

**AC-4 — The public surface contains no standalone eligibility ceremony.**
Verified by: CLI help, command reference, completion, and First Officer instructions omit `gate eligibility`; source retains the private predicate only at mutation boundaries.

## Expected surface

Expected changes are limited to gate CLI routing and help, First Officer gate instructions, public documentation, completion, and focused lifecycle tests. The internal predicate and its table tests remain.

## Test plan

Start with the real operator lifecycle and negative controls. Remove the public command only after these tests prove that `consume` owns the same checks.

Run focused gate, CLI, status, terminal-delivery, and recorded-lifecycle tests. Then run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.

## Stage-specific test gates

- Ideation must trace every production caller of the public command and the private predicate. It must identify any consumer that cannot use `status` or the acting command.
- Implementation must remove no safety branch from the private predicate.
- Validation must use a detached mutation that weakens one consume guard and make sure that a behavioral test fails.

## Stage Report: ideation

- DONE: traced every production `gate eligibility` route and every
  `EvaluateEligibility`/`EligibilityFileAt` caller, including the status and merge
  consumers that are not test-only.
- DONE: recorded the concrete unreleased-v1 CLI/status/docs design, the no-alias/no-
  readiness-class boundary, and the operator replacements through status and acting
  commands.
- DONE: defined the real lifecycle, refusal matrix, detached mutation proof, and
  focused/full/race/formatting/live-lane gates, with baseline focused evidence and
  the unrelated shared-state manifest failure called out.

### Production caller trace and authority boundary

The tree-wide production trace is small and has two distinct classes of caller.

| Caller | Current use | Decision for the cut |
| --- | --- | --- |
| `internal/cli/cli.go:newGateCommand` (`Use`, help, subcommand allow-list, and the `gates.EligibilityFileAt` branch) | Exposes the read-only `gate eligibility` ceremony and prints `condition`/`eligible` labels. | Delete the public route and all help/usage/allow-list bytes. Keep `gate consume` as the only acting command and retain its refusal line. An unknown `eligibility` operand must be a usage error; no alias or deprecation path. |
| `internal/status/discover.go:materializeGateEligibility` | Computes `gate-condition`/`gate-eligible` only when an explicit status field or filter requests them. This is a second read-only projection of the authority predicate, not a scheduler input. | Remove this diagnostic projection and its two derived field names. `status --boot --identify` and `gate-readiness` already project the durable next action; no operator needs the predicate labels. Keep the canonical gate fields (`gate-state`, `gate-application`, `gate-target-stage`) and the single readiness vocabulary. |
| `internal/status/merge.go:pendingTerminalApproval` | Classifies the pending terminal approval before clearing `mod-block`, including the legacy gate-less fallback. | Keep. This is part of the merge authority boundary, not an operator preflight. It must refuse unreadable, stale, consumed, superseded, wrong-target, and malformed authority before any delivery mutation; `gates.FinalizeTerminalApproval`/`SupersedeTerminalApproval` repeat the check under the entity lock. |
| `internal/gates/application.go:ConsumeAt` | Calls `EvaluateEligibility` after taking the entity lock, then atomically spends `pending` or supersedes stale input. | Keep every branch, including the second successor check, `Wrote` distinction, and one locked status/application write. |
| `internal/gates/delivery.go:lockPendingTerminalApproval` | Calls the same pure predicate while holding the lock for terminal finalize and rework. | Keep unchanged; both terminal writers must continue to share the fail-closed predicate. |

`gates.EligibilityFile` has no production caller (only its `EligibilityFileAt`
delegation and tests). Tests, examples, and prose references are not authority
consumers. No production path needs to run a preflight: the First Officer can use
`status --boot --identify --json`/`gate-readiness` to choose the next action, then
use the acting command (`gate consume` or `merge guard`) to obtain the authoritative
result. The merge classifier is the sole consumer that cannot be replaced by a
status read: status cannot prove retained-room currency while deciding whether it
is safe to clear `mod-block`, so it remains an in-boundary fail-closed check.

### Concrete unreleased-v1 design and expected surface

The implementation is a deletion, not a compatibility layer:

1. Remove `eligibility` from `newGateCommand`'s `Use` string, help text, accepted
   subcommands, unknown-subcommand diagnostic, and dispatch branch. Preserve the
   existing `consume` output (`condition`, target, `consumed`, and terminal route)
   and its exit-code contract. A refusal must remain byte-clean except for its
   actionable stdout/stderr diagnostic; a real advance or stale supersede alone
   emits `sync=... phase=consume`.
2. Remove `materializeGateEligibility` and `gate-condition`/`gate-eligible` from
   status's derived-field registry and focused status fixtures. `gate-readiness`
   (`awaiting-captain`, `withdrawn-awaiting-prepare`,
   `approved-awaiting-advance`, `approved-awaiting-merge`, or `validating`) is the
   only scheduling projection. Do not add a readiness value for eligibility and do
   not make any status projection an authorization input.
3. Keep `EvaluateEligibility`, `ConsumeAt`, `lockPendingTerminalApproval`, and the
   terminal delivery writers. Do not move the predicate outside the entity lock,
   weaken `validateRetainedAuthority`, remove the successor/terminal-target checks,
   or infer authority from `status` alone. `pending -> consumed` (or stale
   `pending -> superseded`) and the status/terminal delivery write remain one
   compare-before-replace candidate.
4. Update the contract and operator surfaces: remove the command from
   `docs/specs/gate-resolution-frontmatter-contract.md` and
   `docs/site/reference/command-reference.md`; change the concept/FO prose to
   say that status projects the next action and the acting verb reports refusal;
   remove eligibility from the FO preflight/resume checklist and from the shared
   FO phrase that says revise “repeats eligibility.” The static bash/zsh completion
   already has no gate subcommands; add a negative assertion rather than adding a
   completion alias.

No replacement diagnostic command, compatibility alias, deprecation warning,
stored state class, or new readiness vocabulary is part of this v1 cut. The private
predicate remains an implementation detail of authority-spending boundaries; the
operator sees durable readiness and the refusal emitted by the command that would
spend authority.

### Falsifiable lifecycle and refusal proof

The first implementation fixture must be a real CLI lifecycle with a command log:

* `status --boot --identify --json` (or `status --fields gate-readiness,...`) shows
  `approved-awaiting-advance` or `approved-awaiting-merge`.
* The log then records close-and-consume (`gate record ... --consume`, or a durable
  `gate record` followed by `gate consume`) with no eligibility invocation.
* A nonterminal approval advances and is handed to `dispatch build --stamp`; a
  terminal approval remains pending and is handed to `merge guard`. The final
  entity snapshot must match the existing consumed/terminal predicates and the
  command log must contain zero `gate eligibility` argv entries.

The negative matrix must run through the acting command and assert both refusal
reason and unchanged unauthorized state: stale retained input (pending becomes
superseded with no status advance), wrong stage/removed successor, wrong target,
already consumed authority, already superseded authority, and malformed or
digest-stale retained binding. Terminal delivery must separately prove that a
stale/malformed binding refuses before clearing `mod-block`, and that `--rework`
refuses a missing, undefined, or terminal `feedback-to`. Repeated consume and
repeated terminal delivery remain byte-clean and emit no sync line. These controls
exercise `TestConsumeRefusesTargetRemovedFromCurrentWorkflow`,
`TestConsumeStaleSupersedesWithoutEffect`,
`TestConsumeRepeatAfterStaleSupersedeReportsNoWrite`,
`TestMergeGuardRefusesDigestStaleAuthorityByteClean`,
`TestTerminalDeliveryRefusalsByteClean`, and the recorded lifecycle refusal/resume
matrix while leaving `TestEligibilityFailClosedTable` as the private predicate
table.

The detached-mutation proof is required before validation is called green. In a
throwaway worktree, remove the `applicationTargetMatches` guard in `ConsumeAt` and
run:

```text
go test ./internal/gates -run TestConsumeRefusesTargetRemovedFromCurrentWorkflow -count=1
```

The test must turn red by observing an unauthorized advance; restore the worktree.
Then (in a fresh throwaway worktree) bypass the retained-authority check in
`lockPendingTerminalApproval` and run:

```text
go test ./internal/cli -run TestMergeGuardRefusesDigestStaleAuthorityByteClean -count=1
```

It must turn red by allowing terminal delivery or clearing delivery state. These
mutants are deliberately exercised before broad cleanup so a green result cannot
be mistaken for proof that the public-command deletion preserved authority.

### Evidence and implementation-gate recommendation

The riskiest existing mechanisms were exercised first with focused tests (no
production source was changed during ideation):

* `go test ./internal/gates -run 'Test(EligibilityFailClosedTable|Consume|TerminalSpend|ApplicationExtension|ReadDiagnostics|Finalize|Supersede)' -count=1` — PASS.
* Focused CLI gate/terminal consume and refusal/recovery tests — PASS (the selected
  package run completed in 107.140s).
* Recorded real-CLI lifecycle/refusal/terminal/provenance tests — PASS (67.047s).
* Focused status readiness/readiness-projection tests — PASS (0.648s).

The combined baseline `go test ./internal/gates ./internal/cli ./internal/status
./internal/ensigncycle` reached the existing gates manifest test but failed because
six entries named by `internal/gates/testdata/v1_pilot_manifest.txt` are absent from
the checked-out shared state (`TestV1PilotManifestReadsAndValidates`); this is a
pre-existing checkout/state fixture condition, not a predicate or CLI failure. The
implementation and validation stages must still run the required focused suites,
`gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`, recording
the manifest condition if the shared state is still incomplete. Host live lanes
should run the real lifecycle when their runtime notes require it; no host-specific
lane is needed for a pure command-surface deletion.

Recommendation: **APPROVE ideation** and dispatch implementation with the
command-log lifecycle first, then the detached guard mutants, followed by the
documentation/help sweep. The intended semantic change is only removal of the
operator ceremony and its redundant diagnostic projection; all authority and
atomicity branches remain.

## Stage Report: implementation

- DONE: Remove the public `gate eligibility` CLI, status projection, help,
  completion, documentation, and First Officer ceremony without adding an
  alias, replacement command, or readiness class.
  Commit `013c8729e` deletes the CLI route/help and both derived status fields;
  unknown-command, unknown-filter, help, and completion tests fail if any public
  surface returns.
- DONE: Preserve the private fail-closed predicate and every locked authority
  check in `gate consume`, terminal finalize, and rework; prove stale,
  wrong-stage, wrong-target, consumed, superseded, and malformed inputs refuse
  safely.
  No production file under `internal/gates/` or `internal/status/merge.go`
  changed; focused gate/CLI suites exercise the refusal matrix and byte-clean
  state, while detached guard mutants made the removed-target consume and
  digest-stale terminal tests fail.
- DONE: Run the real status-to-record-to-consume or merge journey without an
  eligibility preflight, plus focused, formatting, full, race, and required
  detached evidence.
  `TestGateJourneyUsesStatusAndActingCommandsWithoutEligibilityPreflight`
  exercises `status --boot --identify`, `gate record --consume`, and stamped
  dispatch, and fails if its command log contains the removed preflight.
  Focused gates, CLI, status, and ensign-cycle suites pass; `gofmt -w ./cmd
  ./internal` and `git diff --check` pass. Both `go test ./...` and `go test
  ./... -race` ran: all packages except `internal/gates` pass, where the
  pre-existing `TestV1PilotManifestReadsAndValidates` fixture fails because nine
  manifest paths are absent from the shared state checkout.

### Summary

The public eligibility ceremony and its redundant status projection are gone;
operators now read durable readiness from status and receive refusal reasons
from the acting command. Commit `013c8729e` preserves all private locked
authority-spending code and adds lifecycle, absence, and mutation-sensitive
evidence for the clean cut.

### Feedback Cycles

- Cycle 1: REJECTED — merge-tree / moving-target conflicts in `internal/status/handlers.go` and `skills/fo-gate-lifecycle/SKILL.md`; surface 19 files/+100/-102 vs estimate not declared; AC unchanged

## Stage Report: validation

- DONE: Prove that the public gate eligibility command, readiness projection, help, completion, documentation, and First Officer ceremony are absent, with no alias or replacement readiness class.
  Candidate `013c8729e`: behavioral unknown-verb/filter, help, four-shell completion, and status tests pass; scoped public-surface search returns no `gate eligibility`, `gate-eligible`, or `gate-condition` occurrence.
- DONE: Prove that private locked authority checks still refuse stale, wrong-stage, wrong-target, consumed, superseded, and malformed inputs without unsafe writes.
  Focused CLI/status/gates/terminal suites pass their exact-state and byte-clean refusal matrices; the private predicate remains called under acting-command locks.
- DONE: Run the real status-to-record-to-consume or merge journey, detached guard mutants, focused, formatting, full, and race evidence; classify current failures and recommend PASSED or REJECTED.
  The real journey and focused suites pass; two detached mutants are killed; formatting is clean; full/race fixture failures are classified below; recommendation: PASSED.
- DONE: **AC-1 (VALUE) — The normal gate journey has no eligibility preflight.**
  `TestGateJourneyUsesStatusAndActingCommandsWithoutEligibilityPreflight` observes `status --boot --identify`, record-and-consume, and successor dispatch, and fails if the command trace contains an eligibility preflight.
- DONE: **AC-2 — Authority spending remains fail-closed and atomic.**
  `TestEligibilityFailClosedTable`, consume once/stale/repeat/removed-target tests, and terminal refusal tests assert exact authority state, status, cardinality, and unchanged bytes across the required variants.
- DONE: **AC-3 — The operator receives an actionable refusal from the acting command.**
  `TestGatePreparedBriefingLocatorLifecycleAndRefusals` requires concrete missing/tampered locator and digest reasons from `gate consume`; terminal controls require their applicable pending/stale/route reason.
- DONE: **AC-4 — The public surface contains no standalone eligibility ceremony.**
  `TestRemovedGateVerbsAreAbsentAndSideEffectFree`, help/completion tests, and removed-field filter/status tests reject the old surface while source inspection retains only private mutation-boundary checks.
- DONE: Detached consume successor-target mutant.
  Disabling `applicationTargetMatches` made `TestConsumeRefusesTargetRemovedFromCurrentWorkflow` fail after an unauthorized status/application write; this kills weakening of the locked consume guard.
- DONE: Detached terminal retained-authority/digest mutant.
  Bypassing retained-request validation and reviewed-input currency made `TestMergeGuardRefusesDigestStaleAuthorityByteClean` fail after delivery state changed; this kills stale-authority acceptance.
- DONE: Focused and formatting evidence.
  Focused CLI/status/gates/ensigncycle behavior passes; `gofmt -w ./cmd ./internal`, `git diff --check`, and candidate worktree cleanliness pass.
- FAILED: Repository-wide suites complete without fixture failures.
  Both `go test ./...` and `go test ./... -race` fail only `TestV1PilotManifestReadsAndValidates`: nine unchanged manifest paths are absent from the external shared state checkout; all other packages pass.
- DONE: Classify the repository-wide failure and verify the unaffected package boundary.
  Evidence defect/deferred infrastructure risk, not an outcome defect: candidate changes neither manifest nor test; `internal/gates` passes normal and race runs with that environment-coupled test skipped; promote if a candidate-touched durable fixture fails.
- DONE: Semantic adversarial pass and recommendation.
  Empty/malformed, stale, wrong-stage, wrong-target, consumed, superseded, repeat, terminal, and removed-public-input variants preserve one invariant: only current binding authority at the acting boundary can write; recommend PASSED with no material finding.

### Summary

Validated exact candidate `013c8729ee8ae6e792fd421b40304008f79e3c93` as a smaller public gate surface with the private locked predicate intact. All value ACs have behavioral evidence and both required detached guard weakenings are caught; the sole full/race failure is pre-existing shared-state fixture drift, classified as a non-material evidence defect, so validation recommends PASSED.

## Stage Report: implementation (cycle 2)

- DONE: Reconcile only the two reported moving-target conflicts against current
  origin/main: internal/status/handlers.go and skills/fo-gate-lifecycle/SKILL.md;
  preserve both current-main behavior and BV's public-surface deletion.
  Merge commit `698867babe7d57eb309dca476ae91187e92a3a57` retains
  current-main's unconditional `--next` readiness materialization and
  `needs-preparation` First Officer flow while omitting the removed eligibility
  projection and preflight.
- DONE: Keep the shipped change limited to public gate eligibility removal;
  preserve private locked authority checks and add no alias, readiness class,
  or replacement ceremony.
  The reconciliation changes no private gate/terminal authority file; scoped
  source search finds no public eligibility command or derived status field,
  and focused unknown-verb/filter, status, lifecycle, completion, and skill
  integration tests pass.
- DONE: Return a clean fresh head with a clean merge-tree against current
  origin/main, focused evidence, formatting, full, and race results; report the
  new exact surface and head.
  Exact head is `698867babe7d57eb309dca476ae91187e92a3a57` with
  `origin/main` `8b5af99baa5c37fe7c969a904819041688420e22` as an ancestor;
  `git merge-tree --write-tree origin/main HEAD` exits 0 with tree
  `ed8b3948d858f4d7d5e374976b90c2c8b02fba53`. Formatting and focused
  suites pass. Full and race runs pass every package except the pre-existing
  manifest fixture: two shared-state paths remain absent
  (`codex-launch-multi-agent-v2.md`, `gate-agent-ergonomics.md`).

### Summary

The two moving-target conflicts are semantically reconciled on current main.
The fresh candidate keeps current scheduler/First Officer behavior, removes the
public eligibility ceremony, and leaves private atomic authority checks intact.

## Stage Report: validation (cycle 2)

- DONE: Re-review exact head 698867ba and confirm that the two conflict resolutions preserve current-main scheduler and First Officer behavior while removing the public eligibility ceremony.
  Exact head `698867babe7d57eb309dca476ae91187e92a3a57` keeps unconditional `--next` readiness materialization and the FO `needs-preparation` cold-report flow, with no eligibility route, projection, or preflight.
- DONE: Re-run the four acceptance criteria, both detached authority mutants, focused and formatting evidence, and confirm a clean merge-tree against current origin/main.
  All AC-focused suites and both mutants pass their falsification requirement; read-only `gofmt -l`, `git diff --check`, candidate cleanliness, ancestor check, and merge-tree exit 0 pass.
- DONE: Run full and race evidence, classify only current-head failures, and recommend PASSED or REJECTED without mutating the candidate.
  Both commands ran on the exact clean head; only two shared-state manifest paths fail as classified below; recommendation: PASSED.
- DONE: **AC-1 (VALUE) — The normal gate journey has no eligibility preflight.**
  `TestGateJourneyUsesStatusAndActingCommandsWithoutEligibilityPreflight` observes status, record-and-consume, and successor dispatch, and fails if the command trace contains the removed ceremony.
- DONE: **AC-2 — Authority spending remains fail-closed and atomic.**
  Consume and terminal matrices pass stale, wrong-stage, wrong-target, consumed, superseded, repeat, malformed, terminal, exact-state, cardinality, and byte-clean controls.
- DONE: **AC-3 — The operator receives an actionable refusal from the acting command.**
  Acting-command tests require concrete locator, digest, pending-authority, stale-authority, and route reasons rather than a collapsed readiness diagnostic.
- DONE: **AC-4 — The public surface contains no standalone eligibility ceremony.**
  Unknown-verb/filter, help, four-shell completion, status, lifecycle, and skill integration tests pass; scoped public-surface search returns no old command or derived fields.
- DONE: Preserve current-main scheduler semantics in the conflict resolution.
  `TestGateReadinessPromotesMechanicallyCompleteColdReports` proves both boot and unconditional `--next` emit one `needs-preparation` row; dirty/malformed reports remain excluded.
- DONE: Preserve current-main First Officer semantics in the conflict resolution.
  The reconciled lifecycle routes `needs-preparation` through one cold-report review/prepare flow while retaining the five acting verbs and omitting eligibility preflight.
- DONE: Detached consume successor-target mutant.
  Disabling locked `applicationTargetMatches` makes `TestConsumeRefusesTargetRemovedFromCurrentWorkflow` fail after unauthorized status/application movement.
- DONE: Detached terminal retained-authority/digest mutant.
  Bypassing retained-request validation and reviewed-input currency makes `TestMergeGuardRefusesDigestStaleAuthorityByteClean` fail after delivery state changes.
- DONE: Clean current-main ancestry and merge-tree.
  `origin/main` `8b5af99baa5c37fe7c969a904819041688420e22` is an ancestor; `git merge-tree --write-tree HEAD origin/main` exits 0 with tree `ed8b3948d858f4d7d5e374976b90c2c8b02fba53`.
- FAILED: Repository-wide suites complete without fixture failures.
  `go test ./...` and `go test ./... -race` fail only two manifest entries now archived in shared state: `codex-launch-multi-agent-v2.md` and `gate-agent-ergonomics.md`; every other package passes.
- DONE: Classify current-head failures and unaffected boundary.
  Evidence defect/deferred infrastructure risk, not an outcome defect: active-root manifest paths lag archive moves; `internal/gates` passes normal/race with that environment-coupled test skipped; promote if a candidate-owned gate behavior or durable fixture fails.
- DONE: Semantic adversarial recommendation.
  The reconciled scheduler and FO paths retain current-main semantics, and only current binding authority at an acting boundary can write; no material finding remains, so recommend PASSED.

### Summary

Re-validated exact merged head `698867babe7d57eb309dca476ae91187e92a3a57` without candidate mutation. The two conflict resolutions preserve current-main `--next`/`needs-preparation` behavior, the four gate-removal ACs and both detached authority mutants hold, and the clean merge-tree plus non-material shared-state fixture classification support PASSED.
