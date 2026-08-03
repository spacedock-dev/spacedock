---
title: Remove the standalone gate eligibility ceremony
status: ideation
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
started: 2026-08-03T23:33:06Z
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

### Ideation completion checklist

- DONE — traced every production `gate eligibility` route and every
  `EvaluateEligibility`/`EligibilityFileAt` caller, including the status and merge
  consumers that are not test-only.
- DONE — recorded the concrete unreleased-v1 CLI/status/docs design, the no-alias/no-
  readiness-class boundary, and the operator replacements through status and acting
  commands.
- DONE — defined the real lifecycle, refusal matrix, detached mutation proof, and
  focused/full/race/formatting/live-lane gates, with baseline focused evidence and
  the unrelated shared-state manifest failure called out.
