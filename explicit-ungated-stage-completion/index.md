---
title: Make stage occupancy and ungated completion explicit
status: ideation
source: Captain discussion of email triage seed skipping intake, 2026-09-09
started: 2026-09-09T20:33:27Z
completed:
verdict:
score: 0.8
worktree:
issue:
pr:
id: 1ytadmcakh5s27r8wjn7qa6f
gates:
    version: 1
    records:
        - id: gate:1ytadmcakh5s27r8wjn7qa6f:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:1ytadmcakh5s27r8wjn7qa6f-backlog-1
              briefing:
                id: briefing:1ytadmcakh5s27r8wjn7qa6f:backlog:attempt-1:revision-1
                digest: sha256:ed30f7885e6d58e90499bdda191ab2f9e3f2cffb4e02c23f532845e0b98be368
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:1ytadmcakh5s27r8wjn7qa6f:backlog:1
                briefing: briefing:1ytadmcakh5s27r8wjn7qa6f:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-09-09T20:33:12.337022Z"
                decision: approve
                reason: 'Captain requested: file this and dispatch for ideation; then explicitly overrode the unrelated dirty-state startup blocker. Approval authorizes ideation only.'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:1ytadmcakh5s27r8wjn7qa6f:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:1ytadmcakh5s27r8wjn7qa6f-ideation-1
              briefing:
                id: briefing:1ytadmcakh5s27r8wjn7qa6f:ideation:attempt-1:revision-1
                digest: sha256:ca7b907305fc10019661d8878ce775883ff480119e1d5ea376f1a8f3f4de875a
                room-ref: '@review/ideation/briefing-1'
---

Design the smallest lifecycle change that makes stage occupancy distinct from completed stage work. The captain requested filing and ideation only; implementation requires a later design approval.

## Problem

A fresh email batch has status intake but no intake artifact. In intake -> triage -> done, the scheduler selects triage directly when intake is initial. Initial stages are exempt from completion proof. Non-initial ungated stages infer completion from committed reports. Gated work has a durable gate-attempt boundary. These meanings conflict.

The intended email flow is intake fetching a closed batch, triage classifying it and preparing human approval, then an execution hook reconciling approved changes before archive. Adding a queued stage hides the lifecycle inconsistency.

## Proposed direction for ideation

Prefer one ungated completion field, provisionally stage-complete: false|true. Status identifies the occupied stage. Initial selects the starting stage, not an exemption from work. An incomplete ungated stage dispatches itself. The FO validates the worker report and records completion through a guarded binary operation. Workers retain report ownership and never mutate lifecycle state.

Gate attempts remain authoritative for gated work. Do not add a competing gated completion flag or a new attempt registry. Determine the smallest representation and command surface; the field name and grammar are proposals, not approved implementation.

Actual entry or revision resets ungated completion atomically. Idempotent status stamps and dispatch retries must not reset completion. Specify how a new execution avoids accepting an old report. Decide absent-field semantics and migration explicitly so legacy backlog workflows do not suddenly execute their seed stages. Include initial gated seed behavior in that compatibility decision.

## Risk evidence

Current code: internal/status/entered_stage.go exempts initial, gated, and terminal stages from entered-stage completion checks. internal/status/format.go dispatchAnalysis projects the successor unless current-stage work is incomplete. internal/status/entered_stage_test.go explicitly tests initial-stage successor projection.

The email pilot observed missing intake artifacts with triage selected. Its no-op fixture covered approval-to-archive, not fresh-seed dispatch. Reproduce the initial-stage selection in a small isolated fixture before designing the change. Do not mutate live email workflows or call providers.

Relevant contracts: skills/first-officer/references/fo-dispatch-core.md, skills/commission/SKILL.md, docs/specs/gate-resolution-frontmatter-contract.md, and the state creation/mutation/dispatch implementations.

## Acceptance criteria

**AC-1 — Fresh batches execute intake before triage.**
Verified by: a fixture creates an incomplete intake seed and observes intake as the first dispatch target; triage becomes eligible only after validated intake completion. Restoring initial-stage successor projection must fail the check.

**AC-2 — Completion is durable and cannot leak across stage entry or revision.**
Verified by: fixture transitions exercise completion, restart, actual re-entry, and idempotent dispatch. A stale report or retained prior completion must not skip new work.

**AC-3 — Gated decision authority and legacy workflow behavior remain explicit.**
Verified by: fixtures cover gate approval/consume, ungated destinations, initial gated seeds, and the chosen absent-field migration rule. A duplicate completion authority or silent legacy seed dispatch must fail a check.

## Ideation deliverables

Recommend one minimal design, exact command and stored-field semantics, transition table, compatibility decision, affected contract wording, expected files/net LOC/tolerance, and a focused test plan. Name existing proof owners and the distinct failure each added test catches. Challenge whether the proposed field is necessary before expanding it. Stop at the ideation approval gate.

## Out of scope

Template-backed new command parameters, seed typo validation, email classification policy, generic external-proof registries, provider operations, and implementation during ideation.

## Ideation design

Recommend explicit completion for ungated work, with gradual entity migration. Keep `status` as occupancy and `stage-complete: false|true` as the sole ungated completion authority. Add one narrowly scoped freshness datum, `stage-report-baseline: N`, where N is the count of exact-current-stage report sections present at entry. This is a proposal for the ideation gate, not implementation authorization.

### Necessity and smallest alternatives

AC-1 could be fixed by removing the initial-stage exemption alone. That would silently execute legacy backlog seeds and would still equate a worker's report with the FO's acceptance. A boolean is necessary to distinguish occupied work, reported work, and accepted completion across an FO restart (AC-1/AC-2).

A boolean alone is insufficient for AC-2: resetting true to false leaves the previous report in the body, and the existing validator would accept it again. The baseline is freshness evidence, not another completion flag. Count exact-stage report headings using the same selector grammar as completion validation; the latest report must have an ordinal greater than the entry baseline. Editing an old report or adding unrelated body text cannot satisfy this check. Existing ensigns already append reports on correction rounds, so no new worker identifier or report format is needed. The FO still judges whether the newly appended evidence addresses the assignment; a copied report is not semantically sufficient merely because it is mechanically fresh.

Alternatives rejected: a synthetic queued stage hides the occupancy ambiguity; a timestamp relies on clocks and report timestamps that do not exist; scanning Git for dispatch commit messages makes message spelling, history depth, rebases, and unstamped paths part of lifecycle authority. A report digest plus entry identifier or an attempt registry adds more state than a baseline integer. The proposed baseline and restart operation serve AC-2 specifically, not an independent feature. There is no completion state for gated work.

### Stored fields and command grammar

- `stage-complete` accepts exactly `false` or `true`. The binary writes canonical scalar values; JSON field projection follows the existing string-valued status convention. Blank or malformed explicit values are invalid, not legacy.
- `stage-report-baseline` is a binary-owned nonnegative integer. Both keys occur together on an explicitly managed ungated nonterminal stage. Baseline without completion, completion without baseline in persisted state, or these keys on gated/terminal stages fail validation. `new` may accept seed input containing only `stage-complete: false` and generates the baseline atomically. Seed input `true` is refused.
- `spacedock status --workflow-dir DIR --set REF stage-complete=true` is the FO acceptance operation. It requires an ungated nonterminal current stage, a structurally complete latest exact-stage report, current entity bytes tracked and clean in local HEAD, and a report ordinal above the stored baseline. It changes only completion to true. Repeating it when already true is a byte-preserving success. It never advances status, commits, or pushes; the FO follows it with the existing `spacedock state commit REF --workflow-dir DIR` boundary before advancing.
- Setting `stage-complete=false` is allowed only to opt a legacy ungated entity into explicit incomplete work; it atomically captures the baseline. Repeating false on an already incomplete explicit entity is a byte-preserving success. It refuses true-to-false with an instruction to use restart, so a retry cannot accidentally re-open completed work.
- `spacedock status --workflow-dir DIR --restart-stage REF` explicitly starts another execution of the same ungated nonterminal stage. It requires a clean tracked entity and atomically writes false plus the current report count. It clears neither gate history nor worker reports. This is an intentional new execution, not a dispatch retry; do not invoke it automatically on worker/runtime retry. Repeating it with identical state and body is naturally a no-op. The FO owns its commit and subsequent dispatch.
- Direct writes to `stage-report-baseline`, clearing lifecycle keys, or mixing a completion/restart operation with `status=...` or other updates are refused byte-clean, including with `--force`. `--force` does not bypass completion/freshness or gate authority. Unsupported stage kinds and malformed state fail closed. Existing administrative overrides outside this lifecycle remain outside this task.
- Ordinary `status=DEST` and gate-consume paths share one destination-initialization helper. A real destination change writes status and destination completion fields in the same atomic entity replacement. An explicit incomplete source cannot advance through `status --set`, with or without a worktree, until accepted. Gated exits retain their current decision/application guards. Repeated same-stage stamps and `dispatch build --stamp` preserve completion and baseline.
- New operations use the existing entity-resolution/root-selection behavior. They retain the existing mutation text style (`stage-complete: false -> true`); existing default table columns and dispatch row JSON keys stay unchanged. `--fields stage-complete,stage-report-baseline` exposes stored values. Restart reports the affected field changes; errors use the existing nonzero status mutation convention.

### Transition and authority table

| Event | Completion and baseline | Dispatch/authority effect |
| --- | --- | --- |
| Explicit ungated initial seed created with false | false; count captured from seed body | Select current intake, even when successor is terminal |
| Explicit ungated occupied stage, fresh report committed | unchanged false | Current stage remains work target until FO acceptance |
| FO validates then records true and commits | true; baseline retained | Successor becomes eligible; restart sees the same result |
| Actual entry into an ungated nonterminal destination | false; destination-stage count captured | Dispatch destination; old destination reports cannot satisfy completion |
| Same-stage stamp, dispatch retry, unrelated field write | preserve both keys | No new execution and no lost completion |
| Explicit same-stage revision/restart | false; recapture count | Requires a newly appended report before acceptance |
| Gate approval recorded but not consumed | no ungated completion authority | Existing application remains pending; completion setter is refused |
| Nonterminal gate consume into ungated work | status plus false/baseline in one replacement | Approval is consumed once; destination must execute |
| Gate revise/hold recorded without a status change | unchanged gate-owned state | No implicit completion or restart; feedback routing owns subsequent entry |
| Actual entry into gated or terminal stage | remove both ungated keys atomically | Gate attempts or existing terminal/delivery rules alone apply |
| Initial gated seed, new or legacy | no ungated keys | Preserve committed-clean seed as gate artifact; no invented worker report |

Existing worktree, capacity, gate, terminal, and delivery suppression remains in force. Incomplete intake becomes its own eligible target; this does not bypass concurrency limits or create a second running worker. The FO keeps addressable-worker checks and existing dispatch ownership.

### Explicit legacy migration

Absence of both fields means legacy behavior on reads: initial ungated seeds keep successor projection; noninitial ungated stages keep committed-report inference; initial gated seeds keep committed-seed gate preparation. Do not mutate state during `--next`, boot, or validate. No schema-version field, bulk rewrite, or inferred migration from missing reports is introduced.

`new` preserves absent-field seeds as legacy so existing backlog templates and automation do not change behavior on upgrade. Commissioning must distinguish a seed-only backlog from executable intake: omit the fields for an explicitly retained legacy seed-only backlog; include `stage-complete: false` for a stage that owns initial work. This compatibility exception is documented as legacy, not presented as the new meaning of `initial`.

An existing initial legacy backlog can leave normally; the first actual ungated destination entry initializes false/baseline, migrating forward without executing backlog. An existing working legacy entity can be adopted as pending via `stage-complete=false` (fresh report required), or explicitly accepted through `stage-complete=true` after the old committed report passes existing validation; that acceptance initializes baseline to one less than the current report count. This latter operation is deliberate adoption of legacy proof, not a new execution. A legacy initial stage without a report cannot use it to fabricate completion. To run a legacy intake seed now, explicitly opt it into false or use restart. Same-stage stamps never migrate an absent seed.

### Observed spike and reproducibility

On 2026-09-09, the assigned launcher reported `spacedock 0.28.0-pre2`, contract 3. Repository HEAD was `af70297ddae6ec64444849e8e3fcf57484bc16e1`. A temporary Git repository declared `intake (initial) -> triage (gate) -> done (terminal)`, worktree false and concurrency 2. The binary created a batch with `status: intake` and no report; the fixture committed it and called `status --next --json`. Creation, Git commit, and status all exited 0. Observed row: `{"id":"fresh-batch","slug":"fresh-batch","current":"intake","next":"triage","worktree":"no"}`. Expected after explicit opt-in: next is intake. No email workflow or provider was touched.

Reproduce with a source-built launcher (`go build -o /tmp/spacedock-intake-spike ./cmd/spacedock`) and a disposable directory. Create `README.md` with this frontmatter:

```yaml
commissioned-by: spacedock@1
entity-type: batch
id-style: slug
stages:
  defaults: {worktree: false, concurrency: 2}
  states:
    - {name: intake, initial: true}
    - {name: triage, gate: true}
    - {name: done, terminal: true}
```

Enclose that frontmatter in `---` delimiters. In the disposable directory run `git init`, pipe `---\ntitle: Fresh batch\nstatus: intake\n---\nNo intake artifact exists.\n` (actual newlines) to `/tmp/spacedock-intake-spike new fresh-batch --workflow-dir "$PWD"`, commit README and the created entity, then run `/tmp/spacedock-intake-spike status --workflow-dir "$PWD" --next --json`. Declare local fixture Git identity if needed. Implementation adds `stage-complete: false` to the seed and asserts intake first; retain the absent-field version as the compatibility control.

The existing deterministic owner also ran successfully: `go test ./internal/status -run 'TestEnteredStageLegacySuppressionControls/initial_stage_keeps_successor_projection' -count=1 -v`. It asserts legacy initial current=backlog,next=implementation through both next and boot, independently confirming the live fixture's interpretation. The proposed baseline/reset mechanism is not implemented or claimed as proven by this spike.

### Proof plan and semantic boundaries

Write the focused failing fixture for each stage before implementation. Reuse current proof owners rather than adding a parallel harness:

| AC | Primary existing owner; proposed behavior check | Distinct falsifying change | Cost/type |
| --- | --- | --- | --- |
| AC-1 | `internal/status/entered_stage_test.go`: create explicit intake seed; compare next and boot before report, after committed report but before acceptance, and after acceptance; add intake-to-terminal case | Restore unconditional initial successor projection, or infer explicit completion directly from report | Small deterministic Git fixtures |
| AC-1 | `internal/status/native_new_from_root_test.go`: real new path writes false plus baseline, rejects true seed, preserves absent seed | Make new silently omit or default completion incorrectly | Small deterministic creation fixture |
| AC-2 | `internal/status/entered_stage_test.go` mutation cases: accept/commit/reload; same-stage stamp preserves true; actual leave/re-enter and restart recapture baseline; old report, edited old report, unrelated edit, dirty report, wrong-stage report refuse; newly appended valid report succeeds | Retain true across entry, reset on a stamp, omit freshness comparison, or accept dirty report | Medium table of existing local-Git fixtures; no sleeps/network |
| AC-2 | `internal/dispatch/build_stamp_test.go`: repeated stamped dispatch preserves completion/baseline and keeps the existing entry commit ownership | Dispatch stamps reset the lifecycle or create fresh acceptance authority | Small deterministic integration case |
| AC-3 | `internal/gates/application_test.go`: approval consume initializes ungated destination in the same write; repeated consume preserves destination progress; failed/stale consume preserves bytes | Split reset from consume, reset on repeat, or let completion bypass application | Medium existing gate fixtures |
| AC-3 | `internal/gates/prepare_initial_seed_test.go` plus existing initial suppression case: committed initial gated seed still prepares without report; absent ungated backlog still skips itself; malformed partial opt-in refuses | Require a gated seed report, silently run absent backlog, or treat malformed explicit fields as legacy | Small deterministic controls |
| AC-1/2/3 | `internal/ensigncycle/shared_keep_moving_durable_test.go` existing durable journey owner: add explicit intake journey and require acceptance commit between report and successor; feedback return requires fresh report | FO prose omits acceptance, stamps reset completion, or feedback consumes old report | Medium shared behavior fixture; retain its existing runtime lane, no provider calls |

A small shared helper may hold field parsing, exact-stage report counting, and entry/reset calculation for status and gates; it must not import CLI or own dispatch, Git history, gate decisions, or a general proof registry. Move the existing report selector grammar into that helper rather than creating competing report matching. Full tests and race tests remain implementation checks; no fresh live host/auth experiment is needed because this task changes lifecycle state rather than runtime integration.

Approved semantic boundary sought: two optional stored scalars; one guarded existing `--set` field and one new `--restart-stage REF` action; explicit ungated current-stage scheduling and acceptance; atomic entry/reset through status and gate consumption; explicit absence compatibility; FO acceptance/restart instruction updates; commissioning examples and lifecycle/command/schema docs. Gated decision vocabulary, attempt identities, terminal approval consumption, delivery hooks, worker body ownership, runtime adapters, and provider behavior do not change. A gate writer's preservation contract gains only the named destination lifecycle-field exception.

### Expected surface and estimate

Estimate net LOC change: +650, across 24 files (approximately 820 insertions, 170 deletions). Tolerance: net +450 to +850 and 20 to 28 files. This is the implementation baseline for gate approval, excluding this ideation-only entity body. Crossing the semantic boundary requires review even within the numeric tolerance.

Expected files: `internal/stagecompletion/stagecompletion.go` (new small shared helper); `internal/status/{entered_stage.go,entered_stage_test.go,gate_extract.go,handlers.go,status.go,new.go,native_new_from_root_test.go,validate.go}`; `internal/gates/{application.go,application_test.go,io.go,prepare_initial_seed_test.go}`; `internal/dispatch/build_stamp_test.go`; `internal/cli/{help.go,status_help_test.go}`; `internal/ensigncycle/shared_keep_moving_durable_test.go`; `skills/first-officer/references/fo-dispatch-core.md`; `skills/commission/SKILL.md`; `docs/schema/entity.mdschema.yml`; `docs/specs/gate-resolution-frontmatter-contract.md`; `docs/site/concepts/stage-lifecycle.md`; `docs/site/reference/command-reference.md`. The frontmatter reference page is an expected optional 25th file if its preservation summary needs the same explicit exception. No changes to runtime-specific ensign adapters are planned.

### Proposed contract and documentation wording

- `skills/first-officer/references/fo-dispatch-core.md`, dispatch selection: replace “If current is initial and next is terminal, set dispatch_stage = current” with “For explicit ungated lifecycle state, use the scheduler's next target, including current-stage intake. Retain the initial-to-terminal fallback only for absent-field legacy entities.”
- Same file, completion: replace “If not gated: terminal -> merge; else decide reuse-or-fresh” with “For ungated nonterminal work, validate the report, run status --set REF stage-complete=true, and commit that entity through state commit before successor selection. A report or worker completion message does not itself mark explicit stage work complete. A refused acceptance stops advancement. Gates retain their existing gate lifecycle.”
- Same file, revision: add “A new execution in the same ungated stage uses status --restart-stage REF and a state commit before dispatch. Ordinary same-stage status stamps and worker dispatch retries preserve lifecycle state. Workers append their canonical report and never write completion fields.”
- `skills/commission/SKILL.md`, seed template guidance: add “Status names the stage the entity occupies. For an initial ungated stage that performs work, include stage-complete: false in seed input; new captures the report baseline. For an intentionally seed-only legacy backlog, omit completion fields and explain that compatibility choice. Initial gated seeds remain committed artifacts for gate review and carry no ungated completion fields.”
- `docs/site/concepts/stage-lifecycle.md`, after opening paragraph: add “Occupying an ungated stage does not mean its work is complete. Explicit stages dispatch themselves until the first officer accepts a fresh committed report and records stage-complete: true. Entering or restarting work resets completion; retrying dispatch does not. Initial chooses where work starts. Older entities with no completion fields retain their legacy seed behavior until explicitly adopted or advanced.”
- `docs/site/reference/command-reference.md`, status entry: add “--set REF stage-complete=true accepts a fresh committed ungated report; --set REF stage-complete=false adopts legacy work as pending; --restart-stage REF starts another execution of the current ungated stage. Each is followed by state commit. Completion cannot be forced or applied to gates.” Update ungated-to-terminal example to show acceptance and commit before the existing finalization command; preserve the warning that terminal-target gate approvals are consumed only by merge guard.
- `docs/schema/entity.mdschema.yml`: add optional canonical `stage-complete` (boolean) and `stage-report-baseline` (nonnegative integer), pair/stage-kind validation, binary ownership, absence migration, and the reset/acceptance rules above. Stored scalars remain projectable through the existing status string representation.
- `docs/specs/gate-resolution-frontmatter-contract.md`: replace the consumed-successor sentence “After the worker report is durable, ordinary atomic terminal fields can complete that successor without --force” with “After fresh ungated work is accepted and stage-complete is committed, ordinary atomic terminal fields can complete that successor without --force. Consuming entry initializes false and its destination report baseline atomically. Repeated consume preserves the destination's work state.” Qualify preservation wording: “Status-changing gate writes may initialize or remove only stage-complete and stage-report-baseline with the destination status; other top-level fields and the Markdown body retain their existing preservation guarantee.”

## Stage Report: ideation

- DONE: Recommend the smallest consistent lifecycle design with exact field/command semantics, gate authority, reset behavior, and explicit legacy migration; challenge the proposed field before adding machinery.
  AC-1/AC-2/AC-3: Ideation design specifies one completion boolean and a justified report-count freshness baseline, guarded acceptance/restart, an authority table, and gradual absent-field migration; no attempt registry or runtime mechanism.
- DONE: Exercise the riskiest stage-selection path in an isolated fixture and record observed evidence plus behavioral proof plans for AC-1, AC-2, and AC-3, reusing existing proof owners.
  AC-1: new/commit/next fixture exited 0 and wrongly selected triage from reportless intake; existing initial projection test passed, and the proof table names reset, gate, legacy, dispatch, and durable-journey falsifiers for AC-2/AC-3.
- DONE: Commit the ideation design and canonical stage report with AC citations, expected files and net LOC with tolerance, semantic boundaries, and proposed contract/doc wording; stop before implementation.
  This path-scoped state artifact contains the design and report, +650 net LOC across 24 expected implementation files with explicit tolerance, and concrete wording for the affected contracts and docs; no implementation files were changed.

### Summary

Recommend explicit ungated completion plus a current-stage report baseline, preserving gate authority and legacy seeds while requiring executable intake to opt in. The isolated observed failure and existing proof-owner plan are recorded above; implementation awaits ideation approval. Required full and race test commands encountered unrelated compile errors in untracked `tmp/dispatch-027/stash-recovery-copy/untracked-files/archived_annotation_decode_test.go` (`decodeProviderResult` and `verifyProviderResolution` undefined); the focused status test passed, and required gofmt ran with its unrelated whitespace-only changes restored.
