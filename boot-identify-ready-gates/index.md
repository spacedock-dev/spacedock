---
title: Expose ready-gate entities in boot identify JSON
status: ideation
score: 0.9
source: PR #493 local-live m3 investigation on 2026-07-10. The default Claude greeting omitted the already-gated Gate Check entity because status --boot --identify --json returned only dispatchable entities; dispatchAnalysis intentionally suppresses current gate stages, so the authoritative boot record made the shipped greet requirement impossible to satisfy.
id: 8n55etrw9wj10jfejdq5f1s8
worktree: .worktrees/spacedock-ensign-boot-identify-ready-gates
gates:
    version: 1
    current:
        gate: gate:docs-dev:8n:validation
    records:
        - id: gate:docs-dev:8n:validation
          stage: validation
          attempts:
            - id: gate-attempt:8n-validation-1
              briefing:
                id: briefing:docs-dev:8n:validation:attempt-1:revision-1
                digest: sha256:74fb6c8bf5bac2091c553e58084933d1c50ada8329310f959e1e49480b88555e
                digest-domain: canonical-bytes
                room-ref: ./review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:8n:validation:1
                briefing: briefing:docs-dev:8n:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-23T08:32:11.092975Z"
                decision: revise
                reason: 'The live workflow counterexample shows current validation stage is not gate readiness: five tickets share the stage, only three have complete reports, and one complete ticket still points at its old ideation gate. The shipped projection must derive from durable current-stage gate lifecycle state.'
                adoption-note: 'helps only if gate readiness becomes first-class state—not something the FO infers from status: validation.'
              application:
                action: feedback
                target-stage: implementation
                state: pending
sprint: durable-decisions
---

## Problem

The interactive first-officer contract requires the greeting to name gates that are
ready for a Captain decision without rendering their reviews. `status --boot
--identify --json` exposes only `dispatchable`, which intentionally excludes entities
whose current stage is a gate. In the real PR #493 m3 shallow-boot record,
`gate-check.md` existed at `review`, but no boot field identified it.

Cycle 1 solved the missing-identity problem with commit `c5a96678`, but selected every
active entity whose current status named a declared gate stage. The Captain's live
counterexample disproves that rule: five entities all displayed `validation`; only
`mf`, `r4`, and `2n` had complete validation reports awaiting gate work or merge,
while `sp` and `qc` were still validating. The result was 5 projected rows for 3
ready entities, a 167% projection. `mf` also retained `gates.current` at its old
ideation gate, proving that neither status nor “some gate record exists” identifies
the current decision opportunity.

Readiness must therefore be a projection of canonical durable gate lifecycle state:
the entity's current workflow stage, its selected same-stage logical gate, that
gate's latest attempt and Briefing, and any Resolution/application. Boot must not
read Stage Report prose to decide readiness, and a second readiness field in
frontmatter would create two authorities that can drift.

## Landed vocabulary and preserved counterexample

PR #557 (`fa240a76`) landed the 3k canonical `gates` reader and the current gate
summary vocabulary: gate, attempt, open/closed attempt state, Briefing, Resolution,
decision, application/action/state, and target stage. PR #560 (`f06cce04`) landed
the h1 application conditions and one-use consumption. Status already exposes those
values as explicit human/JSON fields:

`gate`, `gate-attempt`, `gate-state`, `gate-briefing`, `gate-resolution`,
`gate-decision`, `gate-application`, `gate-application-state`, `gate-condition`,
`gate-eligible`, and `gate-target-stage`.

Commit `c5a96678` remains the executable, fixture-backed negative control: it proves
the identify-only additive schema and ordering work, while its stage-only selector
produces the Captain's 5/3 false-positive result. Implementation must preserve that
commit in branch history and correct it with a later commit; do not drop, rewrite, or
revert it.

No spike is needed. The design uses landed, exercised mechanisms: strict
`gates.Read`/`Validate`, `CurrentSummary`, status's existing gate projections,
`gate record --briefing`, application eligibility/consume, and the insertion-ordered
boot JSON builder. The riskiest state-selection defect is directly captured by the
stale-pointer control in the test plan rather than by an unverified parser or runtime
handoff.

## Approaches considered

1. **Recommended: one canonical readiness reducer over 3k/h1 state.** Project a
   selected same-stage attempt into a compact readiness value and reuse that one
   result for boot and human status. No new durable field or parser.
2. **Keep the stage-only selector and ask the FO to inspect reports.** This preserves
   `c5a96678`'s implementation but recreates the 5/3 counterexample and makes boot
   prose-dependent.
3. **Persist `ready: true` or a separate readiness ledger.** This makes readiness
   first-class by duplication rather than authority. It can disagree with the
   selected gate/attempt/application and requires another mutation protocol.
4. **Scan for any gate record whose stage matches status.** This would make `mf`
   appear ready while its durable selection still points to ideation, hiding the
   exact broken transition the projection must expose.

## Proposed durable readiness reducer

Add one pure reducer in `internal/gates` and call it from status entity discovery. It
accepts the already-validated canonical document plus the entity's current status;
it does not read the entity body or a retained room. The same reducer populates
`gate-readiness` for human/JSON status and the identify-only boot rows.

The reducer is fail-closed and uses this state machine:

| Durable current-stage state | `gate-readiness` | In `ready_gates` |
|---|---|---|
| Declared nonterminal gate stage, but no canonical selected same-stage attempt yet | `validating` | no |
| Canonical selected same-stage latest attempt is open and has its Briefing binding | `awaiting-captain` | yes |
| The same attempt is closed by `approve` and carries `advance/pending`, explicit empty blockers, no active execution hold, and a target stage | `approved-pending` | yes |
| Pending approval has blockers or an active execution hold | `blocked` or `held` | no |
| Closed `revise` has `feedback/pending` | `feedback-pending` | no; 6y routes it immediately |
| Application is `consumed`, `superseded`, or `not-applicable` | that canonical application state | no |
| The `gates` block is malformed or internally inconsistent | `invalid` | no |
| Current stage is not a declared nonterminal gate | empty | no |

“Approved awaiting merge” is not a second application state. It is the canonical
`approved-pending` row whose durable `target_stage` names a stage marked
`terminal:true` in the already-emitted taxonomy. `gate consume` remains the only
effect authorizer and rechecks eligibility/digest before atomically changing status
and application state.

## Exact boot and human surfaces

`ready_gates` remains identify-only, follows `stages`, emits `[]` at zero, and uses
the existing status order: declared stage, score descending with empty scores last,
then lexical discovery order. Every row has these string keys in this fixed order:

```json
{"id":"mf","slug":"mf","current":"validation","gate":"gate:docs-dev:mf:validation","attempt":"gate-attempt:mf-validation-1","state":"open","briefing":"briefing:docs-dev:mf:validation:attempt-1:revision-1","resolution":"","decision":"","application":"","application_state":"","target_stage":"","readiness":"awaiting-captain"}
```

The approval-pending terminal handoff is equally explicit:

```json
{"id":"2n","slug":"2n","current":"validation","gate":"gate:docs-dev:2n:validation","attempt":"gate-attempt:2n-validation-1","state":"closed","briefing":"briefing:docs-dev:2n:validation:attempt-1:revision-1","resolution":"resolution:spacedock:docs-dev:2n:validation:1","decision":"approve","application":"advance/pending","application_state":"pending","target_stage":"done","readiness":"approved-pending"}
```

The `mf` row is the proof that completion corrected selection: `current` is
`validation` and `gate` names the validation gate. A stale
`gates.current=...:ideation` must not yield that row even if a validation record or
complete report exists.

Human status keeps all existing `gate-*` fields and adds the computed
`gate-readiness`. It remains opt-in through `--fields`/JSON and participates in
`--all-fields`; the default five-column table stays byte-compatible. This makes the
counterexample directly visible:

- `sp`, `qc`: `status=validation`, `gate-readiness=validating`
- `mf`, `r4`: `gate-state=open`, `gate-readiness=awaiting-captain`
- `2n`: `gate-state=closed`, `gate-decision=approve`,
  `gate-application=advance/pending`, `gate-target-stage=done`,
  `gate-readiness=approved-pending`

`gate-condition`/`gate-eligible` retain h1's effect-eligibility meaning and may perform
their existing deeper retained-Briefing check when explicitly requested.
`gate-readiness` is the local durable scheduling projection used at boot; it never
claims that an effect has already passed `gate consume`.

## Selection correction and lifecycle owners

Completion handling must prepare/bind and thereby select the current-stage gate
before presentation. 6y owns the FO procedure and its behavioral proof; this task
owns the narrow binary invariant needed by that procedure:

1. The stage worker writes its report and completes. State is `validating`; no boot
   row is inferred from the report.
2. 6y verifies the completion/checklist/ACs, prepares the canonical Briefing, and
   calls `gate record --briefing`. The recorder opens or rebinds the current-stage
   attempt and atomically sets `gates.current` to that gate even when the same
   Briefing bind is otherwise idempotent. This repairs `mf` without hand-editing and
   without creating a duplicate attempt.
3. State commit makes the folder package and selected attempt durable.
   `ready_gates` now reports `awaiting-captain`.
4. At engage, 6y consumes that row and invokes `present-gate`; xb alone owns the
   selected chat/provider transport and exact retained Result/association.
5. 6y records the exact Captain/delegated/provider result. The recorder closes the
   same attempt with a Resolution and derived application; state commit makes it
   durable. Approval to a terminal successor now reports `approved-pending` plus the
   terminal `target_stage`.
6. 6y validates, runs eligibility, and calls consume. h1 atomically marks the
   application consumed and advances status. The row disappears.
7. The existing terminal branch immediately runs merge hooks/default merge,
   terminalizes, archives, cleans up, and tears down. This task changes no merge,
   presentation, or FO judgment procedure.

These boundaries make restart recovery mechanical at steps 3 and 5. A boot after
binding re-presents the open attempt; a boot after approval resumes
eligibility/consume and then the terminal merge path. A boot before binding remains
honestly `validating`.

## Implementation scope

- Preserve `c5a96678` as an ancestor; integrate current `main` with a merge, then add
  correction commits. Do not rebase/drop/revert the preserved counterexample.
- `internal/gates/model.go` (and `io.go` only if needed for an absent-record error
  class): one pure current-stage readiness projection over the landed canonical model.
- `internal/gates/operation.go`: make same-Briefing/open-attempt binding select the
  current-stage logical gate atomically and idempotently.
- `internal/status/discover.go`, `format.go`, `boot.go`, and `json_commands.go`: retain
  the canonical projection on each active entity, materialize `gate-readiness`, replace
  the stage-only selector, and render the exact rows.
- `internal/gates/gates_test.go` or `internal/cli/gate_test.go`: stale-pointer selection
  correction, no duplicate attempt, unrelated bytes unchanged.
- `internal/status/boot_identify_test.go`, `gates_coexist_test.go`, and
  `json_boot_test.go`: the five-entity counterexample, human fields, ordering, ordinary
  boot omission, and dispatchability separation.
- `docs/site/reference/command-reference.md` and
  `docs/site/concepts/gates-and-decisions.md`: document the lifecycle row and the
  bind-before-ready boundary.

Expected production delta after integrating main: approximately 100-180 Go LOC;
tests approximately 250-450 LOC; docs 10-20 lines across 8-11 touched files.
Tolerance 2×. Reconfirm if implementation needs schema fields, report parsing, a new
command, presentation/FO skill changes, merge changes, or another durable readiness
store. No skill, provider transport, or live-presentation UI change belongs here.

## Out of scope

- Rendering gate reviews during boot.
- Parsing Stage Report prose or treating completeness wording as durable state.
- A second gate/readiness schema, compatibility parser, cache, or ledger.
- Changing dispatchability, gate judgment, application consumption, blockers, or
  merge semantics.
- 6y's First Officer procedure/live proof and xb's presentation/provider transport.
- Repairing or weakening the separate m3 scenario oracle in this change.

## Acceptance criteria

- **AC-1 (VALUE — exact ready population).** In the five-entity Captain fixture,
  identify reports exactly the three durable unresolved gates (`mf`, `r4`, `2n`) and
  excludes both still-validating entities (`sp`, `qc`): 3/3 true positives and 0/2
  false positives. Preserved `c5a96678` produces the 5/3 counterexample. *Verified
  by:* a native CLI fixture with exact raw rows and a mutant selector that keys only
  on `gate:true` stage and therefore fails with five rows.
- **AC-2 (selection is current-stage and completion-owned).** With `mf.status` at
  `validation`, an existing validation attempt, and `gates.current` still selecting
  ideation, boot excludes `mf`. Repeating `gate record --briefing` with the same
  current-stage Briefing changes only the selection pointer, creates no attempt, and
  makes the exact validation row appear. *Verified by:* before/after whole-file YAML
  comparison, attempt counts/ids, and the exact `mf` boot row; a no-pointer-update
  mutant remains at 2/3 and fails.
- **AC-3 (human distinctions share the reducer).** Status exposes
  `gate-readiness=validating` for `sp`/`qc`, `awaiting-captain` for open `mf`/`r4`,
  and `approved-pending` for closed approved `2n`, alongside the exact canonical
  gate/attempt/Briefing/Resolution/application fields. *Verified by:* text and JSON
  `--fields` tests fed the same fixture as boot; every entity's expected readiness
  is independently enumerated.
- **AC-4 (approved terminal work is restart-visible).** Before consumption, `2n`'s
  row preserves its approval Resolution, `advance/pending`, and terminal
  `target_stage=done`; after the real `gate consume`, status is `done`, the
  application is consumed, and `2n` is absent from `ready_gates`. *Verified by:* a
  command-level close→boot→eligibility→consume→boot sequence plus the combined 6y
  journey's observed entry into the existing terminal merge path.
- **AC-5 (fail-closed, prose-independent projection).** Malformed gates, a stale
  old-stage selection, a current gate without a Briefing attempt, terminal/unknown/
  archived/non-gate entities, blocked/held approvals, feedback pending, consumed,
  superseded, and not-applicable applications never appear. Adding/removing a
  complete Stage Report body without changing frontmatter cannot change any boot
  row. *Verified by:* table-driven frontmatter/body mutants and before/after exact
  JSON comparisons.
- **AC-6 (schema, ordering, and scheduling compatibility).** Rows use the fixed
  13-key order above and existing status ordering; identify keeps `ready_gates`
  after `stages` and renders `[]` at zero. Ordinary boot omits the field, default
  human status is unchanged, and raw `dispatchable` remains equal to `--next`.
  *Verified by:* raw key/order assertions, zero/one/many fixtures, default-table
  golden, ordinary-boot negative assertion, and exact dispatchable comparison.

## Test plan

1. Start with the five-entity native fixture. Give all five `status: validation`;
   `mf`/`r4` carry open current-stage attempts, `2n` carries a closed approval with
   empty blockers and terminal target, and `sp`/`qc` have no selected current-stage
   attempt. Initially leave `mf.gates.current` at ideation. Assert two ready rows,
   run the real same-Briefing record command, then assert the exact three rows.
2. Run human text/JSON projection over the same fixture and assert all five
   readiness values plus canonical ids. Body-only report-completeness edits must
   leave outputs byte-identical.
3. Add table controls for malformed/cross-stage/briefing-less records; blocked,
   held, feedback, consumed, superseded, and not-applicable applications;
   terminal-plus-gate, unknown, archived, and ordinary non-gates. Each is a negative
   exact-row oracle, not an absence-of-crash assertion.
4. Drive an open attempt through real decision record and consume. Pin the
   approval-pending terminal row before consumption, then the consumed status and
   missing row afterward. Removal of pointer correction, application-state checks,
   or current-stage agreement must make a named assertion fail.
5. Preserve the existing identify local-only/key-order/discovery tests, ordinary
   boot omission, default table golden, and exact `dispatchable == --next` control.
6. Run focused gate/status/CLI tests, `gofmt -w ./cmd ./internal`, `go test ./...`,
   `go test ./... -race`, and live-tag compilation. After 6y lands, its combined
   fixture/live journey owns the FO presentation→decision→consume→terminal-merge
   proof; this task supplies the durable rows it consumes and does not duplicate
   that model-costly run.

## Documentation change proposal

In `docs/site/reference/command-reference.md`:

```diff
-folds in the stage taxonomy, and reports the boot sections
+folds in the stage taxonomy and canonical ready-gate lifecycle rows, and reports
+the boot sections. A row names the selected current-stage gate, attempt, Briefing,
+Resolution/application state, target stage, and readiness; a gate stage alone is
+not ready
```

After the recorded-decision paragraph in `docs/site/concepts/gates-and-decisions.md`:

```diff
+After completion verification, the First Officer binds the retained Briefing before
+presenting the gate. That bind selects the current-stage gate attempt. Startup can
+then distinguish work still validating, an open attempt awaiting the Captain, and
+a recorded approval awaiting one-use consumption. Approval to a terminal target is
+consumed before the existing merge and terminalization path begins.
```

## Stage Report: ideation

- DONE: Audit the current boot JSON schema, dispatchAnalysis gate suppression, status stage/entity ordering, and existing boot identify tests. Confirm the observed m3 record truly cannot expose `gate-check` through any current field.
  The code audit found gate suppression before worktree/concurrency, status ordering by declared stage then score then slug, and no boot inventory beyond dispatchable; the real m3 record names only `merged-pr`, while `gate-check.md` independently records `status: review`.
- DONE: Choose the smallest stable schema that lets the greet name ready gates without changing dispatchability. Specify exact fields, deterministic ordering, zero-gate behavior, compatibility impact, and exact implementation/test files; avoid entity-body reads or prompt changes.
  The recommended identify-only `ready_gates` field carries ordered `id`/`slug`/`current` rows after `stages`, emits `[]` at zero, leaves ordinary boot and dispatchable bytes unchanged, and scopes implementation to five Go files plus the command reference.
- DONE: Refine AC-1 through AC-4 with independent offline evidence and the m3 live end-value handoff. Append a complete ideation Stage Report and recommendation; do not change product code or claim unrun tests passed.
  AC-1 measures 3/3 ready gates against the current 0/3 baseline; AC-2/AC-3 pin dispatch and schema compatibility offline; AC-4 hands the authoritative row to m3's real live greet without claiming its separate oracle is fixed here. Recommend approval for implementation.

### Summary

The design adds one identify-only, append-only `ready_gates` array and reuses existing active-entity parsing and status ordering. It supplies the missing m3 gate identity while preserving ordinary boot output, dispatchability, prompts, and entity-body read boundaries.

## Stage Report: implementation

- DONE: TDD the missing identify schema on coupled base 5ec370c0. Add a real native boot-identify mixed fixture that is RED because 3/3 current-gate entities are absent, then cover deterministic status order, exact row keys (`id`,`slug`,`current`), m3-shaped `gate-check`/`review`, and zero-gate output.
  RED: `TestBootIdentifyReadyGates` received an absent `ready_gates` value instead of the exact three rows, and `TestBootIdentifyFoldsDiscoveryTaxonomyLocalPR` reported the zero-gate key missing; ordinary boot already passed its omission pin. GREEN: the focused `TestBoot(Identify|JSON)` command passed all 11 selected tests.
- DONE: Implement the smallest append-only identify surface: `ready_gates` appears only in `--boot --identify --json`, after `stages`; ordinary `--boot --json` remains byte-compatible; dispatchable/--next semantics and bytes remain unchanged; terminal, terminal+gate, unknown, archived, and ordinary non-gates are excluded. Reuse parsed active entities/stages and existing ordering; no entity-body rereads, prompts, stage mutations, or m3 assertion changes.
  Commit `c5a96678` adds the status-ordered selector in `format.go`, identify-only boot data in `boot.go`, fixed-row rendering in `json_commands.go`, native controls in `boot_identify_test.go`/`json_boot_test.go`, and the command-reference wording; no other path changed.
- DONE: Run focused boot/status/CLI tests including `TestBootJSONDispatchableMirrorsNext`, exact key-order/zero-gate/ordinary-boot controls, git diff --check, scoped gofmt, go test ./internal/status ./internal/cli, go test ./..., go test ./... -race, and live-tag compile. Commit only the ready-gate schema/tests/reference changes above base 5ec370c0; exclude the known journeydelta formatting defect.
  `git diff --check` and scoped `gofmt -d` emitted no output; status/CLI passed 1,049 tests in 2 packages; full and race gates each passed 2,116 tests in 17 packages; live-tag compile exited 0. `git merge-base HEAD 5ec370c0` equals `5ec370c0`, and the six-path diff excludes `journeydelta.go`.
- DONE: Append a complete implementation Stage Report with RED/GREEN evidence, exact files and commit/ancestry, AC-1 through AC-3 evidence, AC-4 explicitly handed to the combined m3 live rerun, and clean status.
  AC-1: the mixed native fixture returns 3/3 exact ordered gate rows, including `gate-check` at `review`. AC-2: its independently specified raw dispatchable value equals `--next`, while terminal, terminal+gate, unknown, archived, and non-gate rows stay absent. AC-3: zero gates render `[]`, `ready_gates` follows `stages`, and ordinary boot omits the field. AC-4 remains assigned to the combined PR #493 m3 live rerun; this commit changes no live oracle.
  Independent review found 0 Critical, 0 Important, and 0 Minor issues and returned Ready: Yes; the code worktree was clean at report time.

### Summary

Boot identify JSON now exposes active, non-terminal ready gates as fixed `id`/`slug`/`current` rows in deterministic status order. The additive identify-only field preserves ordinary boot and dispatchable output, and the committed native tests cover the m3-shaped value, exclusions, ordering, and zero-gate behavior; the combined m3 live run remains the AC-4 handoff.

### Feedback Cycles

- Cycle 1: REVISE — Captain live-workflow counterexample; surface 5 current-validation rows vs 3 durably complete unresolved gates (167%); AC narrowed: readiness requires a prepared current-stage gate attempt after completion verification, not merely `gate: true` stage membership

## Stage Report: ideation (cycle 2)

- DONE: Define one durable readiness state machine that mechanically distinguishes validating, awaiting captain, and approved awaiting merge from canonical current-stage gate/attempt/application data.
  The revised body defines one pure 3k/h1-backed reducer: absent/mismatched selection is `validating`, a selected open Briefing attempt is `awaiting-captain`, and a closed approve plus unblocked `advance/pending` is `approved-pending`; terminal target taxonomy makes the merge handoff explicit without another stored state.
- DONE: Specify the exact boot JSON and human status fields, ordering, transition owners, and backward route from completion verification through Briefing preparation, engage presentation, decision recording, consumption, merge, and terminalization.
  The fixed 13-key rows and `gate-readiness` human field are specified byte-order-first; the lifecycle assigns verification/bind/present/record/consume to 6y with xb transport, h1 effect consumption, this task's projection/pointer invariant, and the existing terminal merge path.
- DONE: Replace the stage-only fixtures with controls covering three complete unresolved gates, two incomplete validations, a stale old-stage gates.current pointer, approval pending merge, and mutation-proof negative cases while preserving dispatchability separation.
  The test plan starts from the Captain's five-entity 5/3 counterexample, requires same-Briefing bind to repair `mf` without a duplicate attempt, pins `2n` before/after consume, rejects body-only and malformed-state mutants, and keeps raw `dispatchable == --next`.

### Summary

Re-ideation replaces `c5a96678`'s stage-only meaning with a canonical selected-attempt readiness projection while preserving that commit as the executable counterexample. The design makes completion binding correct `gates.current`, exposes restart-safe open and approval-pending rows, and leaves FO procedure, presentation transport, consumption safety, dispatchability, and merge semantics with their existing owners.
