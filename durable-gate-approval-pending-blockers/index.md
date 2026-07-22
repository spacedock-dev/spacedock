---
title: Gate recorder — durable gates records with binary-owned writes
status: implementation
score: "0.80"
source: "Captain design feedback, 2026-07-13."
id: 3kd1x1gfxr8mdwzbmnwtjbw8
started: 2026-07-18T08:58:53Z
gates:
    version: 1
    current:
        gate: gate:docs-dev:3k:ideation
        attempt: gate-attempt:3k-ideation-10
    records:
        - id: gate:docs-dev:3k:ideation
          stage: ideation
          current-attempt: gate-attempt:3k-ideation-10
          attempts:
            - id: gate-attempt:3k-ideation-1
              sequence: 1
              state: closed
              briefing:
                id: briefing:docs-dev:3k:ideation:attempt-1:revision-8
                digest: sha256:3a8fd6d6702d212d72b708a406549a3a4c1d3f81997887e36d3453755721825b
                room-ref: ./review/ideation/briefing-8
                note: Frozen at closure. The digest binds the briefing-8 gate-summary artifact (summary + full post-cut snapshot), byte-verifiable in the room. Provider result validated by digest equality and retained as provider-result-8.json; the provider envelope id (briefing:single-file:e63586cd350f4f7b6cdcaa074a1ff312) is normalized to this attempt briefing id per the recorded id-mapping practice.
              resolution:
                type: Resolution
                id: resolution:actor-1784592481316587000
                briefing: briefing:docs-dev:3k:ideation:attempt-1:revision-8
                by: person:reviewer
                at: "2026-07-21T00:08:01Z"
                decision: revise
                reason: 1. why are there still 14 ACs? i thought we trimmed this. 2. take a look at PR#510 to see where things align
              application:
                action: feedback
                state: consumed
                target-stage: ideation
              note: 'Subspace advisory float on the rebuilt tip binary, probe-first ritual observed. Two asks: physically trim the body to the cut (the AC section still carries every pre-cut criterion in full; the scope-cut prose named the retained set but never restructured the sections), and produce an alignment read against open draft PR #510 (Ledger gate-binding boundary). Routed to a fresh ideation revision worker; attempt 2 opens at re-presentation.'
            - id: gate-attempt:3k-ideation-2
              sequence: 2
              previous-attempt: gate-attempt:3k-ideation-1
              state: closed
              briefing:
                id: briefing:docs-dev:3k:ideation:attempt-2:revision-9
                digest: sha256:1c229dfe87f5954b2b1e6b7a54cc4918cddf55e35bb66f198fba7f6ccbb3d28a
                room-ref: ./review/ideation/briefing-9
                note: Frozen at closure; byte-verifiable in the room. Provider result validated by digest equality, retained as provider-result-9.json; provider envelope id (briefing:single-file:201ca46ba902b9da0ec874243ee2c000) normalized to this attempt briefing id.
              resolution:
                type: Resolution
                id: resolution:actor-1784596837823868000
                briefing: briefing:docs-dev:3k:ideation:attempt-2:revision-9
                by: person:reviewer
                at: "2026-07-21T01:20:37Z"
                decision: approve
                reason: is there any reason to keep the split AC? like we need to do a final integration test? i'd like to keep things lean if possible. / is it easier to keep this one for integration test and split a clean gate/resolution implementation? or not necessary
              application:
                action: advance
                state: superseded
                target-stage: implementation
              note: Approve with two attached captain questions (pointer-AC leanness; integration-umbrella split), answered in chat post-recording; fork 1 (id namespacing) not annotated, so recorder ids stay Spacedock-internal per the stated default. Superseded by attempt 3 (the captain-directed pointer-AC cut); the approval itself stands.
            - id: gate-attempt:3k-ideation-3
              sequence: 3
              previous-attempt: gate-attempt:3k-ideation-2
              state: closed
              briefing:
                id: briefing:docs-dev:3k:ideation:attempt-3:revision-10
                digest: sha256:cb816a084445eefd588a9b5119522ca6c2a70ab375de005a9d206af444d2b362
                room-ref: ./review/ideation/briefing-10
                note: Byte-verifiable frozen entity snapshot (entity-snapshot.md) taken after the cycle-14 pointer-AC cut and before this record was written — no advisory digest needed.
              resolution:
                type: Resolution
                id: resolution:captain-chat-3k-ideation-3
                briefing: briefing:docs-dev:3k:ideation:attempt-3:revision-10
                by: person:captain
                at: "2026-07-21T01:42:54Z"
                decision: approve
                reason: 'Captain-directed leanness fold from the attempt-2 approve questions (''ok do the cleanup''): the eight pointer-AC stubs cut so the AC scanner sees exactly the seven in-scope criteria; scheduler-rule and test-plan stubs cut where trivial, original numbering kept so gaps mark moved-out steps; the Scope cut section is the traceability record. No integration-umbrella split (captain accepted the recommendation: the contract doc is the clean spec; integration proof rides the sprint DoD and pre-cut audit).'
              application:
                action: advance
                state: superseded
                target-stage: implementation
              note: Fold applied by the live revision worker (cycle 14, state commit 2e562ed9); FO recorded. Superseded by attempt 4 (the captain-directed resolution-first split); the approval itself stands.
            - id: gate-attempt:3k-ideation-4
              sequence: 4
              previous-attempt: gate-attempt:3k-ideation-3
              state: closed
              briefing:
                id: briefing:docs-dev:3k:ideation:attempt-4:revision-11
                digest: sha256:f8cd6fa75043b061dc64aa5583620af14f90a9dc3d557b5c0f246f9eb051a5aa
                room-ref: ./review/ideation/briefing-11
                note: Frozen at closure; provider result retained as provider-result-11.json, digest equality validated; provider envelope id (briefing:single-file:6a66ead293dbb27a4931ec57e370a02b) normalized to this attempt briefing id.
              resolution:
                type: Resolution
                id: resolution:actor-1784599855140796000
                briefing: briefing:docs-dev:3k:ideation:attempt-4:revision-11
                by: person:reviewer
                at: "2026-07-21T02:10:55Z"
                decision: revise
                reason: btw, does the multi-artifact briefing not work? i want to see the mermaid diagram in the spec too
              application:
                action: feedback
                state: consumed
                target-stage: ideation
              note: 'Presentation-side revise, FO-owned (no design change requested): re-present as a multi-artifact briefing package with the contract spec (carrying the mermaid) as its own artifact. Attempt 5 opens on the package presentation; the design content is unchanged from this attempt.'
            - id: gate-attempt:3k-ideation-5
              sequence: 5
              previous-attempt: gate-attempt:3k-ideation-4
              state: closed
              briefing:
                id: briefing:docs-dev:3k:ideation:attempt-5:revision-12
                digest: sha256:ec6bb198f1fc2451b47ffecf904390c9278a161d33f55f40697d2ca4f4020ee0
                room-ref: ./review/ideation/briefing-12
                note: Multi-artifact briefing package — FIRST successful package-mode gate presentation (direct review-v1 float, probe-proven; the probe caught the required-context schema rule). Review log retained in-room as briefing.review.jsonl by the provider itself.
              resolution:
                type: Resolution
                id: resolution:actor-1784601146924137000
                briefing: briefing:docs-dev:3k:ideation:attempt-5:revision-12
                by: person:reviewer
                at: "2026-07-21T02:32:26Z"
                decision: revise
                reason: 'Annotation on the contract mermaid: ''this is too wide and can''t be rendered. is there a way to make it vertical?'' (annotation:captain-1784601092240856000, included). The resolution reason''s route-to-decision observation is subspace-side product feedback per the captain''s follow-up — not filed in this workflow.'
              application:
                action: feedback
                state: consumed
                target-stage: ideation
              note: 'Presentation-content revise: reshape the diagram vertical for the terminal render. Design content still unchanged since attempt 4. Attempt 6 opens at re-presentation.'
            - id: gate-attempt:3k-ideation-6
              sequence: 6
              previous-attempt: gate-attempt:3k-ideation-5
              state: closed
              briefing:
                id: briefing:docs-dev:3k:ideation:attempt-6:revision-13
                digest: sha256:5b128db9cb36e2d690bcceaa279d47b6c8a7da0077c8d1124752323fce903a19
                room-ref: ./review/ideation/briefing-13
              resolution:
                type: Resolution
                id: resolution:actor-1784602325418452000
                briefing: briefing:docs-dev:3k:ideation:attempt-6:revision-13
                by: person:reviewer
                at: "2026-07-21T02:52:05Z"
                decision: revise
                reason: 'Annotation on the contract mermaid: still too wide.'
              application:
                action: feedback
                state: consumed
                target-stage: ideation
              note: 'Render check failed again at the float: the subgraph frames were removed next round. Attempt 7 re-presents with two stacked frameless diagrams.'
            - id: gate-attempt:3k-ideation-7
              sequence: 7
              previous-attempt: gate-attempt:3k-ideation-6
              state: closed
              briefing:
                id: briefing:docs-dev:3k:ideation:attempt-7:revision-14
                digest: sha256:8f003849fd0b059495afcd3fcd7f438fa050d027a36e6a5dc393de6233bd55db
                room-ref: ./review/ideation/briefing-14
                note: Contract artifact carries two stacked frameless diagrams (contract sha256 4ca06d15...). Frozen at closure.
              resolution:
                type: Resolution
                id: resolution:captain-chat-3k-ideation-7
                briefing: briefing:docs-dev:3k:ideation:attempt-7:revision-14
                by: person:captain
                at: "2026-07-21T03:20:00Z"
                decision: approve
                reason: 'Captain approve, re-affirmed in chat: the two stacked diagrams render well (''it looks great'') and h1 goes based on the current 3k. HONEST PROVENANCE: the captain first resolved this attempt in a float pane whose launcher had died — the resolution was written to an unlinked scratch file and destroyed. The chat re-affirmation is the authoritative record; the destroyed float result is float finding 15 and the presentation command''s primary red fixture.'
              application:
                action: advance
                state: superseded
                target-stage: implementation
              note: 'Superseded by attempt 8 (the captain-directed fold: round-disposition section + evergreen restyle); the approval itself stands. h1 dispatched immediately per the captain. Application state corrected pending->superseded at the preflight (the preflight''s second material finding, state half): the attempt-8 recording updated this note but left the state field live, briefly giving the gate two pending advances — banked as the cross-attempt red fixture for the eligibility task.'
            - id: gate-attempt:3k-ideation-8
              sequence: 8
              previous-attempt: gate-attempt:3k-ideation-7
              state: closed
              briefing:
                id: briefing:docs-dev:3k:ideation:attempt-8:revision-15
                digest: sha256:fd95df2a7f7200ffdc3370db13785cf1b2018af42c1ce13e1e05b00af08e5f1a
                room-ref: ./review/ideation/briefing-15
                note: Byte-verifiable frozen entity snapshot + contract snapshot (final contract sha256 9c0ee9ad469ca0399e657b146e70f9de524387851ccbd3a0a4d9a0fd6d4b08b7), taken after the fold pass and before this record.
              resolution:
                type: Resolution
                id: resolution:captain-chat-3k-ideation-8
                briefing: briefing:docs-dev:3k:ideation:attempt-8:revision-15
                by: person:captain
                at: "2026-07-21T04:52:00Z"
                decision: approve
                reason: 'Captain-directed fold, content the captain specified in chat — no re-ask per the attempt-7 record: the round-records/triage-dispositions advisory section folded into the contract from the triage task''s reframe ideation; the evergreen rule applied (component-only prose, task ids confined to removable scaffolding, with the diagram prefixes and example ids explicitly scoped as scaffolding converted at the landing pass); the captain''s approve of the design and diagrams stands through the fold.'
              application:
                action: advance
                state: superseded
                target-stage: implementation
              note: Superseded by attempt 9 (the codex-seat contract reconciliation); the approval itself stands.
            - id: gate-attempt:3k-ideation-9
              sequence: 9
              previous-attempt: gate-attempt:3k-ideation-8
              state: closed
              briefing:
                id: briefing:docs-dev:3k:ideation:attempt-9:revision-16
                digest: sha256:c99e7b8597038912b25f2d2f7fccd631649cc3b635fb57aa566d0ad25318aba9
                room-ref: ./review/ideation/briefing-16
                note: 'RAW-FILE PIN (the marked legacy digest domain, per the digest-domains ruling): byte-verifiable frozen entity + contract snapshots (contract sha256 681b23483f61202094f8c6095cad381f448b98b24e1098716cbf3601b4767aa6), taken after the amendment pass and before this record.'
              resolution:
                type: Resolution
                id: resolution:fo-delegated-3k-ideation-9
                briefing: briefing:docs-dev:3k:ideation:attempt-9:revision-16
                by: agent:first-officer
                at: "2026-07-21T14:35:06Z"
                decision: approve
                reason: 'Recorded by the FO on the captain''s delegated authority under the recording-identity ruling — the ruling''s first exercise. Captain directive, verbatim: ''agree with advisory-to-binding and the 3 recommendations. fix the gate review retirement.'' The amendment pass closed the codex seat''s second, third, and sixth material findings against the contract: the gate-review architecture retired from every operative section in favor of the approved overridable present-gate channel with recorder-side validation; the two digest domains named with shaping history explicitly legacy; consumption semantics aligned authorization-only with the crash windows named and fixtured; the recording-identity sentence itself added to the lifecycle rules.'
              application:
                action: advance
                state: consumed
                target-stage: implementation
              note: The contract now agrees with every approved member design and both preflight seats. The captain has NOT re-reviewed the amended bytes; this closure rests on the quoted directive, recorded honestly under the FO identity — exactly what the ruling prescribes.
            - id: gate-attempt:3k-ideation-10
              state: closed
              briefing:
                id: briefing:docs-dev:3k:ideation:attempt-10:revision-18
                digest: sha256:6b2c4f1388a58f42f7c8610f847ed9e7cce92758c00b201d4eb9f4f89dbedd8b
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-18
              resolution:
                type: Resolution
                id: resolution:first-officer-3k-ideation-10
                briefing: briefing:docs-dev:3k:ideation:attempt-10:revision-18
                by: agent:first-officer
                at: "2026-07-22T04:11:08Z"
                decision: approve
                reason: 'Recorded by the First Officer on the captain''s delegated authority. Captain directive, verbatim: ''lgtm. add the fixture''. Revision 18 changes only the required cross-logical-gate re-entry fixture and preserves every approved revision-17 boundary; the captain did not separately render the folded package bytes.'
        - id: gate:docs-dev:3k:validation
          stage: validation
          current-attempt: gate-attempt:3k-validation-1
          attempts:
            - id: gate-attempt:3k-validation-1
              sequence: 1
              state: closed
              briefing:
                id: briefing:docs-dev:3k:validation:attempt-1:revision-1
                digest: sha256:0a54f1baec0120c1c93523e6900a6ce28e025c570289e5dfa9835e28099042ac
                digest-domain: canonical-bytes
                room-ref: ./review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:first-officer-3k-validation-1-design-reset
                briefing: briefing:docs-dev:3k:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-22T02:01:32Z"
                decision: revise
                reason: 'Captain directive: ''ok. send it back.'' Return to ideation because the validated recorder requires an agent-authored transaction envelope instead of consuming semantic intent and exact Review v1; redesign the durable projection so Git supplies rebind audit history while current and frozen decisions remain self-contained.'
sprint: durable-decisions
group: recorder
worktree: .worktrees/spacedock-ensign-durable-gate-approval-pending-blockers
---

# Gate recorder — durable gates records with binary-owned writes

## Scope cut (captain-approved, 2026-07-21)

This task grew four products in one coat (12 cycles, 15 ACs, two companion specs). It now narrows to ONE: the gate **resolution** recorder and its record schema — the binary that owns every `gates:` frontmatter write for the resolution record (open / rebind-while-open / close-with-resolution / supersede-attempt), the record invariants (pointer agreement, one binding per attempt, frozen closures), snapshot-bound digests, and the status surfacing of recorded resolution state. It records WHAT THE DECISION IS; the `application` layer — WHAT THE DECISION DOES (the one-use advance authorization) — moved to h1 at the captain's 2026-07-21 resolution-first split ("get the resolution right first"). Retained ACs: 1, 4, 6 (record-state subset), 10, 12, 13, 14, each trimmed to its resolution-side core. This half is production-proven as a hand-run convention across eight entities in 0260 shaping (see `production-evidence-2026-07-20-fo-dry-run.md` and the 0260 closure findings: `--set` re-serialization, a self-conflicting attempt pointer, stale applications, the advisory-digest hole, entity-cannot-self-bind) — the recorder mechanizes a proven shape.

Moved out, one owner per concern:

- **The application layer + blockers, execution holds, and dispatch eligibility** → `gate-blockers-and-eligibility`. Two joined concerns. (a) *What the decision does* — the one-use advance authorization: `application` action/target-stage, the `pending`/`consumed`/`superseded`/`not-applicable` states, and application-staleness marking (moved from 3k at the captain's 2026-07-21 resolution-first split: the application halves of ACs 1, 4, 6, and 10; scheduler rules 8, 9 (application-state clause), and 10; test-plan items 2, 4, 5; the `application` section of the contract doc, now h1-owned in place). (b) The original seed — blockers, execution holds, and dispatch eligibility (ACs 2, 3, 5, 11 + the eligibility subset of AC-6; scheduler rules 4-6 and 10's guard beyond convention). h1 already owns exactly-once consumption and blockers, so the application record lands with its consumer. Unlike the blocker half (which the dry run never exercised — all eight recorded approvals had zero declared blockers), the application record HAS a demonstrated consumer: the 0260 Commander consumed pending `advance` applications through the normal transition path, so h1's live-need question should credit it. The "refuses a second pending application" red fixture follows the application layer to h1; the recorder keeps the pointer-conflict and frozen-closure red fixtures. Sequenced after the recorder; its gate re-examines live need.
- **The presentation journey** (ACs 7 and 15; the one-command `gate review` blocking presenter, atomic result retention, briefing packages, probe-snapshot binding, provider id-mapping adapter; the probes companion spec rides along as convention) → `gate-review-presentation-command`. Subspace-coupled — sequenced with the subspace-tui surface, interim ritual per the 0260 shaping debrief.
- **Rejection-rework route context** (AC-9; the durable route edge) → DEFERRED, no task: the feedback-cycle prose convention just shipped in 0260; the binary edge waits for observed drift, per the escalation ordering.
- AC-8's behavioral-test mutants split with their owners.

Spec-level coupling retained here: this task keeps ownership of `gate-resolution-frontmatter-contract.md` as the one spec, with its `application` section now marked h1-owned in place (one doc, many owners — not relocated); the presentation task's id-mapping rule (provider envelope briefing id normalized to the attempt briefing id after digest validation) is SPECIFIED in that contract and implemented in the presentation task. The open gate attempt's briefing must be rebound to this post-cut content before its next presentation (open-attempt rebinding, scheduler rule 2).

**Expected surface + tolerance (FO-drafted; re-estimated at the 2026-07-21 resolution-first split — the recorder shrinks as the application layer leaves; captain corrects at the gate):** Go product code, first in-repo binary surface of the gate cluster: ~2-3 new files under `internal/` (gates block read/model/write for the resolution record + invariant validation + the resolution close/supersede-attempt mutations; NO application struct — that is h1's) plus edits in `internal/status` (surfacing recorded resolution state, `--set` coexistence so unrelated field writes leave `gates:` untouched) and 1-2 new `spacedock gate ...` verb entries in `cmd/`; **~400-650 production LOC** (down from the pre-split ~600-900), roughly equal test LOC (fixture replay of the eight 0260 production entities + the red fixtures: z7's real pointer conflict, a frozen-closure mutation, a resolution/briefing-digest mutation). The recorder round-trips the h1-owned `application` sub-object unchanged on write (preserve-not-model), so replay of the eight production entities — which carry `application` blocks — stays green. Contract doc unchanged in location (the spec is the already-banked `gate-resolution-frontmatter-contract.md`, its `application` section marked h1-owned in place); ~10 lines of FO-contract prose naming the recorder as the gates-resolution-write owner. Tolerance 2×. Hard self-check: any schema change that breaks replay of the eight production entities, any subspace-tui coupling (xb's), any blocker/eligibility computation OR application-record modeling/writing/state-surfacing (h1's) trips a reconfirm.

## Problem

Spacedock currently conflates a human gate decision with immediate stage advancement and dispatch: approve, advance, and spawn are expected in one turn. That fails when a task is reviewable now but must not dispatch until another task lands. The First Officer must otherwise delay a valid review, dispatch against the wrong base, lose the approval, or invent hidden “approved but waiting” state.

Subspace already returns a structured `Resolution`, but the single-file review wrapper deletes its temporary package. The decision therefore does not become durable workflow state and cannot safely survive restarts or dependency waits.

Andrew's 2026-07-18 CMD cutover run exposed the same persistence gap on the
rejection path. Validation completed and correctly rejected the slice because
a production coordinator was missing; feedback routing correctly returned the
entity to implementation. The status surface then showed only
`implementation` at score `0.90`, erasing the legible fact that validation had
completed with a blocker and that the current implementation run was rework
caused by that rejected gate. This is production evidence for this filing, not
a separate task.

## Required capability

Persist Subspace Resolutions in the workflow entity's committed top-level `gates`
frontmatter collection using the hierarchy logical gate → stable gate attempts → one
Review & Gate Briefing binding → exact adopted Resolution → minimal Spacedock
application. One open attempt may select multiple Briefings as the presentation,
design, or evidence changes; this does not imply `revise`. Current frontmatter retains
only one `briefing` reference/digest per attempt. Attempt state makes that binding
replaceable while open and frozen while closed; state Git and Subspace retain prior
snapshots.
Recording the Resolution must not advance `status` or dispatch a worker. Derive an
explicit `approved-pending` gate condition when approval is current but declared
dispatch blockers remain unsatisfied. This is computed gate/eligibility state, not
another lifecycle stage.

Separate four concepts that are currently collapsed:

- open gate attempt and its current immutable presentation snapshot;
- durable gate decision: approve, revise, or hold, with provenance and reviewed digest;
- dispatch blockers: declared dependencies or predicates whose current state is queryable;
- dispatch eligibility: computed from current stage, non-stale approval, blocker satisfaction, and one-use consumption state.

Preserve rejected gate results under the same model. A lifecycle transition
back to a feedback target must not overwrite the completed rejection. The
reducer needs both the ordinary current stage and a durable route relationship
from the rejected gate attempt to the rework stage run.

The persisted representation must be workflow-owned and portable. Temporary Subspace package paths, machine-local pane/session identifiers, prompts, and private user data must not enter durable state.

The exact first-use journey, minimal schema, helper boundary, examples, and lifecycle
are in
[`gate-resolution-frontmatter-contract.md`](gate-resolution-frontmatter-contract.md)
(SHA-256 `681b23483f61202094f8c6095cad381f448b98b24e1098716cbf3601b4767aa6`).
It evolves closed PR #474's entity-frontmatter decision onto Review & Gate v1 instead
of creating a parallel ledger.

The first-use question and rework-comparison flow is specified in
[`gate-review-probes.md`](gate-review-probes.md)
(SHA-256 `03b0abe451764505c3e8dc5a725fcb3622f4bb0a752ba157f2e4c96581f6c693`).
It keeps probe definitions and results in the Git-backed Subspace room while this
entity stores only the stable room reference and durable gate binding.

Briefing 6's fresh ProbeResult changed the comparison but left the answer supported:
the mechanism remains provider-owned and usable without Spacedock; the review supplied
concrete evidence that this room stores Probe history at `../probes.jsonl`. The prior
limitation therefore narrows from no concrete path to a proven instance path with no
universal provider layout. Because Subspace did not surface that result/comparison in
the presentation, 3k supplies a separate semantic-delta summary and records the missing
rendering as a Subspace product gap rather than expanding this task's UI scope.

## Scheduler behavior

1. Opening a gate creates a stable Spacedock attempt and first immutable Briefing
   pointer without changing the current stage or dispatching.
2. While the attempt is open, a changed lens, design, evidence, question, artifact
   revision set, or decision opportunity replaces the entity's current Briefing
   reference/digest under the same attempt. Subspace retains Briefings/logs/lenses,
   re-evaluates affected assessments, and retains the delta. A Briefing binds only a
   frozen pre-run Probe history snapshot; the fresh result/comparison is joined from
   provider storage by Briefing id so appending it cannot invalidate the Briefing. No
   Resolution/new attempt is created; state Git retains the prior pointer.
3. Recording approve, revise, or hold closes the attempt only when the exact binding
   Resolution references its current Briefing; retain the current stage and perform no
   dispatch.
7. If gate-defining input changes after an attempt closed, a new attempt is required;
   closed attempts never gain another Briefing. (Marking the closed application stale
   and keeping the task non-dispatchable is eligibility surfacing owned by **h1**.)
8. Review & Gate `revise` or `hold` closes the attempt and records the exact Resolution.
   (Creating the resulting application — `revise` → a pending feedback application, `hold`
   → `action: none`/`state: not-applicable` — is owned by **h1**.)
9. Status surfaces the recorded stage, gate/attempt/Briefing/Resolution identities, and
   open or closed state. Subspace presents Briefing/lens/assessment deltas through the
   stable room reference when requested. (Application-state, blocker-set, execution-hold,
   and staleness surfacing is owned by **h1**.)
11. A recorded rejection (`revise` Resolution) is retained, and re-entry at the gate after
   that closed result creates a new attempt. (Routing it through `feedback-to` as a
   feedback application is owned by **h1**; projecting the current lifecycle stage with
   explicit `feedback_rework` route context is **DEFERRED**, AC-9.)

## Acceptance criteria

These are the in-scope acceptance criteria for the recorder. Moved and deferred criteria
are cut here to keep the section lean; the **Scope cut** section (captain-approved,
2026-07-21) is the surviving traceability record for what went to h1
(`gate-blockers-and-eligibility`), xb (`gate-review-presentation-command`), and DEFERRED.
Original AC numbers are retained, so the gaps mark moved-out criteria.

**AC-1** A recorded approval — its exact Resolution and reviewed digest — survives process restart and still reports the same durable decision byte-for-byte, without advancing or dispatching. (The approved-pending condition, the exact blocker, and the blocked-entity framing are application/eligibility surfacing owned by h1.)

**AC-4** Review & Gate `revise` and `hold` Resolutions remain durable and visible. A
captain-facing rejection for rework is recorded as a portable `revise` Resolution, not as
the superseded portable `reject` vocabulary. `approve` needs no portable rationale;
`revise`/`hold` require a nonblank reason or an included earlier same-Briefing
Annotation, exactly as Review & Gate v1 specifies. (Non-dispatchability, blocker-clearance
non-override, and the resulting feedback application are owned by h1.)

**AC-6 (record-state subset)** Status text and JSON distinguish the recorded gate
resolution states the recorder surfaces from entity frontmatter alone: a recorded
`approve`, a recorded Review & Gate `hold`, and a recorded `revise`. → The
application-state surfacing (`pending`/`consumed`/`superseded`/`not-applicable`), active
execution hold, unsatisfied/unknown/failed blockers, satisfied-but-not-yet-consumed, and
stale-approval distinctions are application/eligibility surfacing owned by **h1**; the
fuller rejected-gate rework route context is **DEFERRED** (AC-9).

**AC-10 (VALUE)** Recording either an approval or rejection changes only the entity's
versioned `gates` frontmatter collection: current `status` is byte-identical and no
dispatch receipt or worker exists. Deleting projection caches and reading the current
entity directly enumerates every logical gate, gate attempt, immutable Briefing
binding (replaceable when open, frozen when closed), exact adopted Resolution, and
selection pointer. Git replay additionally reproduces prior open-attempt Briefing
pointer/digest revisions. (The recorded `application` sub-object is round-tripped
unchanged but its state semantics are h1's.)

**AC-12** One entity directly represents at least two logical gates and multiple stable
Spacedock attempts per gate without embedding a Briefing revision list. Each open
attempt has exactly one replaceable `briefing` binding; each closed attempt has exactly
one frozen binding under the same field name. Concurrent attempt/pointer forks and mutations of a frozen
Briefing/Resolution fail closed without field-wise merge.

**AC-13** Portable-contract fixtures accept an authorized `approve` with no rationale,
reject `revise`/`hold` when neither a nonblank reason nor an included earlier Annotation
exists, preserve multiple advisory Resolutions without mistaking them for binding, and
reject cross-Briefing `includes` without silently copying log entries. They prove
stage/sequence/Briefing-change/application fields are outside the copied Resolution. A
separate Spacedock authoring-policy fixture requires an explicit nonblank reason only when a
First Officer auto-approves under delegated conn authority; the generic portable parser
and entity schema remain permissive.

**AC-14** Adding/revising a lens or revising design/evidence on an open attempt creates
and selects a new immutable Briefing under the same attempt and creates no `revise`
decision. Current frontmatter replaces only the pointer/digest; state Git preserves
prior bindings, while the stable Subspace room owns full Briefings/logs/lenses,
assessment re-evaluation, and presentable deltas. Closure freezes the exact current
binding. Re-entry after that closed result creates a new attempt.

## Resolved storage decisions

- **Location:** one versioned top-level `gates` YAML mapping in entity frontmatter,
  containing a `records` collection rather than a one-attempt slot.
- **Identity:** each record names a logical entity/stage gate; each attempt id names a
  stable Spacedock gate attempt; each nested Briefing id names one immutable
  portable decision snapshot. The exact adopted binding Resolution is preserved.
- **Reviewed digest:** SHA-256 over RFC 8785 canonical bytes of the immutable Review &
  Gate Briefing, whose artifact revisions and gate-defining context form the reviewed
  manifest.
- **Application:** only a closed attempt has a Resolution and one mutable application.
  First use retains only action, target, one-use state, blockers/hold, and feedback
  route; recording closure and consuming the workflow action remain separate commits.
  Existing transition and dispatch state owns effect identity, receipts, and recovery.
- **History and selection:** frontmatter carries every distinct attempt but only one
  `briefing` binding per attempt. It is replaceable while open and frozen while closed. Git retains
  prior pointer/digest values; Subspace owns full snapshot/log/lens history.
- **Concurrency:** Briefing snapshots and Resolutions are immutable. Compare-and-swap
  serializes pointer/application changes; competing snapshots or closes fail closed.
- **Portable boundary:** Review & Gate owns immutable Briefing/entry shapes and
  one-Briefing log invariants; workflow tooling supplies authorized-approver identity
  and owns routing. Subspace stamps/persists reviewer-app entries. The entity copies
  only an externally authorized Resolution for the exact current Briefing;
  stage/attempt/selection/digest/room-reference/application are Spacedock wrapper state.
- **Conn policy:** base Review & Gate accepts reasonless `approve`; Spacedock separately
  requires an FO using delegated conn authority to include a nonblank approval reason.
- **Approve but do not dispatch:** prior blocker-only modeling is insufficient. The
  contract adds a first-class durable `execution-hold`, separate from portable
  `decision: hold`.

## PR-510 alignment (Ledger gate-binding boundary)

PR #510 (OPEN draft, unmerged: `feat(contract): draft Ledger gate-binding boundary`)
proposes `spacedock.gate-binding.v1` — a thin Spacedock-owned provider *pin* (`ns` +
`entity_ref`, optional `stage`/`target_stage`/`workflow_ref`/`provider_instance_id`/
`expected_revision`) stored inside a Helm Ledger gate slot, plus the consumption
contract for Helm-owned application facts (`helm.application.committed/observed/view.v1`,
the committed receipt, `source_superseded`, generic `provider.binding`, and
`projection.rewrite_quarantined`). It explicitly claims no shipped writer, projector,
or command.

PR-510 and this recorder sit at **different layers describing adjacent facts**. The
recorder owns the whole durable gate tree *inside* entity `gates` frontmatter — the
workflow-owned physical authority (settled 2026-07-21; not reopened here). PR-510 owns
the Spacedock↔Helm-Ledger *boundary* for the case where a portable Helm Ledger is the
gate authority. They are not competing stores. The alignment question is whether the
recorder's record schema names the same concepts compatibly, so that if PR-510 later
lands, the frontmatter records can become the `entity_ref` target of a binding and map
to Helm's application facts without renaming or re-authoritying.

| Element | Recorder record schema (3k) | PR-510 gate-binding boundary | Read |
|---|---|---|---|
| Storage / authority | whole gate tree in entity `gates` frontmatter; Spacedock sole authority | thin provider pin in a Helm Ledger gate slot; Helm Ledger owns the gate | **diverge (settled)** — different layers, not competing stores. Frontmatter authority stands. The entity *is* the `entity_ref` target a binding would point at. |
| Identity | Spacedock-minted `gate:…`, `gate-attempt:…`, `briefing:…`, `resolution:…` strings | Helm-owned `gat_`, `resolution_id`, `application_id` — immutable external ids, forbidden as required binding fields | **diverge → F1** |
| Field names | `target-stage`, `stage`, blocker `expected-revision` (hyphen-case) | `target_stage`, `stage`, `expected_revision` (snake_case) | **align** — same concepts, cosmetic casing; YAML house style keeps hyphens |
| Application states | `pending` / `consumed` / `superseded` / `not-applicable` | `pending_apply` / `applied` / `superseded` / `rewrite_quarantined` | **align + adopt** — `pending`↔`pending_apply`, `consumed`↔`applied` map cleanly; `superseded` is a shared term to adopt; `rewrite_quarantined` has no recorder analog (F4) |
| pending→applied fold | existing Spacedock transition/dispatch commit sets `consumed` | Ledger acceptance of `helm.application.committed.v1` is the sole fold input | **diverge → F2** |
| Supersession | reviewed-input/digest change → `superseded` + new attempt; `previous-attempt` chains; append-only in Git | append-only `helm.application.source_superseded.v1` (`old`/`new_source_pin`); `supersedes_application_id` links a new application | **align** on append-only principle; **diverge** on trigger (recorder = pre-apply reviewed content; PR-510 = post-apply source pin). Recorder has no `supersedes_application_id` (it omits `application.id`) |
| Provider binding | `briefing.room-ref` opaque Subspace review-room pin | generic opaque `helm.provider.binding.v1` (`ns`+additive), specialized by `spacedock.gate-binding.v1` | **align** — same "opaque provider pin, don't interpret" principle; distinct providers (Subspace room vs Helm gate) |
| Receipts / idempotency | deliberately omitted (`effect-receipt`, `dispatch-attempt-id`, `consumed-at`, `application.id`); existing machinery owns effect id + crash recovery | `helm.application.committed_receipt.v1` (`receipt_id`, `idempotency_key`, `body_digest`); idempotency `(application_id, idempotency_key)`; response-loss via `application.view.v1` | **diverge → F3** |
| Digest discipline | `sha256:` over RFC 8785 JCS canonical Briefing bytes | binding digests `sha256:` over RFC 8785 JCS; git tree/blob digests raw-payload SHA-256 | **adopt / identical** — the recorder already matches PR-510's JCS + `sha256:` convention |
| Operation without the other layer | must work without Subspace/Ledger; entity is authority | binding absent when Spacedock absent; Helm keeps its own path | **align** — both assert symmetric optionality |

### Genuine forks flagged for the captain (not resolved here)

- **F1 — identity authority.** The recorder mints and owns all gate/attempt/briefing/
  resolution ids in frontmatter; PR-510 treats `gat_`/`resolution_id`/`application_id` as
  Helm-owned. The settled decision keeps frontmatter authoritative for this sprint, so
  the recorder does not adopt Helm ids as authority. The open captain choice is whether
  to *shape/namespace* the recorder's ids now so a later `spacedock.gate-binding.v1` pin
  maps cleanly, or leave them purely Spacedock-internal. Not required to ship the recorder.
- **F2 — pending→applied fold owner.** The recorder sets `consumed` in the durable state
  change that records the existing transition/dispatch machinery's success; PR-510 folds
  `pending_apply→applied` only on Ledger acceptance of a committed attestation. Both sides
  assert symmetric "works without the other" optionality, so the reconciling reading is
  coexistence — the recorder's apply-once stays authoritative for Ledger-absent operation,
  and a Ledger-bound deployment adds the receipt fold as an outer layer. Captain to confirm
  coexistence vs one subsuming the other; nothing to build this sprint either way.
- **F3 — receipts/idempotency in the record.** The recorder deliberately omits receipt and
  idempotency fields; PR-510 makes them central to the Helm apply fold. Confirming the
  omission stands for this sprint (it does — settled) keeps `consumed` from having to
  mirror a Helm `applied` receipt inside frontmatter.
- **F4 — projection/quarantine + post-apply source drift.** PR-510's epoch cursor,
  `rewrite_quarantined`, and post-apply `source_superseded` lanes have no recorder analog,
  because the frontmatter *is* the authority rather than a projection of an event log.
  These map onto the banked commit-derived event design
  (`artifacts/spacedock-state-commit-event-proposal.md`), which is explicitly out of this
  cut. Captain to note the boundary, not act on it now.

### Expected surface + tolerance, reconfirmed in light of the alignment

The alignment is a boundary read, not new build scope. Because the frontmatter-authority
decision is settled, PR-510 does not pull the recorder into consuming Helm facts this
sprint, and none of F1–F4 requires code in the recorder cut. The FO-drafted **Expected
surface + tolerance** in the Scope cut section therefore **stands as reconfirmed**: Go
product code (~2–4 files under `internal/` + `internal/status` edits + 1–2 `spacedock
gate …` verbs; ~600–900 production LOC ≈ equal test LOC), the contract doc unchanged,
~10 lines of FO-contract prose, **tolerance 2×**. The hard self-check (a schema change
that breaks replay of the eight production entities, any subspace-tui coupling, any
blocker/eligibility computation) is untripped by this alignment. One qualifier: if the
captain elects **F1** (shape/namespace recorder ids now), that adds only id-string
convention — a handful of lines, comfortably inside 2×. Coverage note: the sprint-goal
digest-verifiability DoD ("reproduce the digest from a committed snapshot, closing the
advisory-digest hole") is proven by behavioral-test-plan item 2 and owned by AC-10 (Git
replay reproduces digest revisions) plus AC-12 (frozen binding immutability); no new AC
is minted for it.

## Design proposal and review

The broader commit-derived event design is preserved at
[`artifacts/spacedock-state-commit-event-proposal.md`](artifacts/spacedock-state-commit-event-proposal.md)
(SHA-256 `65e31b3e315d8f87b4527e8f0356999a0b79577b96ccfd7640ef5c9ab5e9fbca`).
It proposes treating the state checkout's Git history as the sole durable event
authority, projecting commits into versioned events, reducing those events into
workflow state, and keeping Zaphod a read-only projection with a separate
runtime-observation overlay.

The 2026-07-13 Subspace single-file review returned an advisory `approve` with
no annotations. The run also reproduced this task's motivating gap: the
Resolution was returned as structured JSON, but the temporary review package
was deleted and no durable workflow Resolution existed until this entity update.

Ideation must retain eight corrections identified during review and production
feedback:

1. `approved_pending_dispatch` must include committed blocker identity,
   version, satisfaction, and failure state. Approval plus “not yet dispatched”
   is insufficient to distinguish safely blocked work from dispatchable work.
2. The gate application is a one-use workflow authorization, not a second dispatch
   journal. It retains action, target, pending/consumed/superseded/not-applicable state,
   blockers/hold, and feedback route. Existing transition/dispatch state owns external
   effect identity, receipt, and crash reconciliation.
3. Git-DAG projection must emit an explicit merge-resolution event whenever
   parents disagree, including when the merge result exactly equals one parent.
   The rule “emit only when the result differs from every parent” can leave an
   inherited contradictory event unresolved.
4. Rejection rework needs no second stage-result event. Extend the existing
   `feedback.cycle_recorded` payload into a durable route edge carrying the
   rejected `gate_attempt_id`, source stage/run, target stage, cycle, and routed
   finding reference/digest; bind the eventual target `stage_run_id` when its
   `task.stage_entered` event appears. The reducer retains this active route
   alongside `stage` until a later gate attempt closes the active rework context.
   Inferring rework from a repeated stage name or prose is insufficient because
   ordinary stage re-entry would become a false positive and workflow
   definitions can change after the historical decision.
5. Projected events need a physical tree authority. The entity's versioned plural
   `gates` collection directly retains logical gates, stable Spacedock attempts, and one
   `briefing` binding per attempt. Attempt state makes it replaceable while open and
   frozen while closed. State Git supplies prior
   pointer revisions; Subspace supplies full presentation history. Attempt closure and
   application consumption are separate commits.
6. Lens and persistent-room semantics remain an integration question. This task fixes
   only the minimum entity-owned workflow binding and does not choose whether a lens
   displays one selected attempt, one gate chain, or the full collection, nor where a
   persistent room stores the provider-owned Briefing and review log.
7. Review & Gate v1 remains the portable authority. One Briefing is one immutable
   decision opportunity and one ordered log; advisory Resolutions are not entity gate
   outcomes, and the first Resolution attributed to the externally supplied authorized
   approver can close the Spacedock attempt only for its exact current Briefing. Stage,
   attempt hierarchy, JCS digest, selection, room reference, and application are
   Spacedock-only; Briefing/lens/assessment deltas remain in Subspace.
8. A Spacedock gate attempt is not a Briefing. It stays stable while an open review
   replaces its current Briefing binding for lens/presentation or reviewed-input
   changes. Frontmatter does not accumulate snapshots: their full objects/logs stay in
   Subspace and their pointer changes stay in Git. A closed result followed by gate
   re-entry starts a new attempt.

The cycle-1 Subspace review adds the minimum first-use and ownership correction. The
First Officer offers the exact `[Y/n]` Subspace journey. On yes, the gate-attempt ensign
owns complete Briefing/Probe presentation, direct annotation handling, revision,
affected-Probe reruns, and durable provider Resolution capture. One `gate review`
command validates the explicit package, derives a canonical title, launches Subspace,
and retains diagnostics/result state. The First Officer owns only binding validation,
entity gate-state recording, captain-facing re-presentation, and later application.
On no, the existing relayed-review path remains unchanged.

The revised contract incorporates those corrections. Implementation still owes the
behavioral proofs below; the ideation no longer carries a known schema or first-use
journey contradiction.

## Behavioral test plan

Retained items prove the recorder's own ACs; split items keep the recorder-proving part
and note the moved rest in parentheses. Fully-moved scenarios are cut here — the **Scope
cut** section records their owners; original item numbers are retained, so the gaps mark
moved-out scenarios.

1. **Physical record contrast (AC-10).** Drive the real binary-owned gate recorder
   against an approving and a revising fixture. Exactly `gates` changes; `status`,
   process roster, dispatch state, and worktree remain byte-identical. Delete projector
   caches and prove the entity still reconstructs the exact Resolution. (The AC-7/AC-15
   presentation side moves to **xb**.)
2. **Cold read, replay, and schema (AC-6 record-subset, AC-10, AC-12, AC-14).** Validate
   the concrete two-gate/multi-attempt example through the shipped schema, restart, and
   invoke status. The direct read enumerates every gate, gate attempt, single `briefing`
   binding, and exact Resolution; Git replay reconstructs prior open
   pointers and reproduces each recorded Briefing digest from its committed snapshot.
   Mutants that require `current-briefing`/`resolved-briefing`, embed provider history,
   or consult a cache fail.
4. **Restart visibility (AC-1).** The recorded-approval fixture survives restart and
   still reports the same durable Resolution and reviewed digest, byte-for-byte. (The
   approved-pending condition, blocker set, and the blocker-state/stale-digest eligibility
   table — AC-2/AC-3/AC-5 with their mutants — move to **h1**.)
5. **Revise/hold durability (AC-4).** Build state history containing a validation
   `revise` Resolution; human/JSON status keeps the durable `revise` after restart. (The
   feedback application it produces, its consumption, and the blocker-override guard move
   to **h1**; the rejected-gate → rework route context and its false-rework contrast,
   AC-9, is **DEFERRED**.)
6. **Re-entry, pointers, and concurrency (AC-4, AC-10, AC-12, AC-14).** Extend the
   fixture through re-validation. Concurrent attempt/pointer writes, a close racing a
   pointer advance, mutation of a closed `briefing`/Resolution, and field-wise merge fail
   closed. (Stale-content supersession, AC-3, moves to **h1**; the rework-context replay,
   AC-9, is **DEFERRED**.)
8. **Portable-boundary and conn policy (AC-4, AC-13).** Feed the recorder a one-Briefing
   ordered log with annotations, two advisory Resolutions, and one later externally
   authorized Resolution. Assert only the binding object is copied exactly. Contrast
   reasonless `approve` (portable-valid), reasonless `revise`/`hold`, `includes` naming
   only an advisory Resolution, `includes` naming an earlier Annotation, and a reasonless
   FO conn-made approval — only the last is rejected by the FO's Spacedock authoring
   policy. (The Briefing-transport/presentation side, AC-7, moves to **xb**.)
9. **Open-attempt Briefing evolution (AC-12, AC-13, AC-14).** Start one open Spacedock
   attempt on Briefing A, add/revise a lens to select Briefing B, then revise evidence to
   select Briefing C. Assert one attempt id, only C in current frontmatter, A→B→C pointer
   history in Git, full snapshots/logs and re-evaluated assessments/deltas in the stable
   Subspace room, and zero `revise` decisions. Reject B-log `includes` of A-log entries.
   Only a Resolution for C closes the attempt and freezes C; later gate re-entry creates a
   different attempt id.

Estimated cost is medium: YAML-node schema/round-trip tests, deterministic Git history,
CLI goldens, and the existing transition/dispatch fake. The riskiest recorder mechanism
runs first: prove the nested writer changes only `gates` and that `--set` on unrelated
fields leaves the gates block untouched, then prove each recorded Briefing digest is
reproducible from its committed snapshot (the advisory-digest hole). The provider-launch
fixture and the "one referenced Briefing reaches the TUI and survives controller failure"
proof move to **xb**; no production host smoke is needed.

## Documentation change proposal

Update `docs/site/reference/frontmatter-contract.md` after the Entity paragraph:

```diff
 Each entity's frontmatter carries its id, current stage, outcome, and worktree state. The contract is [`entity.mdschema.yml`](https://github.com/spacedock-dev/spacedock/blob/main/docs/schema/entity.mdschema.yml), which defines the fields, the custom-field policy, the recognized body headings, and the invariants.
+A gate decision is recorded in the entity's versioned `gates` collection before any stage transition or worker dispatch. A logical gate retains stable gate attempts. Each attempt stores one `briefing` reference/digest: replaceable while open and frozen when closed. State Git preserves prior pointers; Subspace owns full Briefings, logs, Probes, and deltas. Revising reviewed input does not itself record `revise`.
```

Update `docs/site/concepts/gates-and-decisions.md` under “The three calls”:

```diff
- **Approve.** The work advances to the next stage. Approving the terminal stage merges and closes it.
+**Approve.** The decision is recorded first. Eligible work then advances exactly once; unresolved blockers or an explicit “approve but do not dispatch” hold keep it at the gate without losing your approval.
```

Add to the gate-review command reference:

```diff
+### Review a complete gate Briefing
+
+When the First Officer offers Subspace and you answer yes, the gate-attempt ensign runs `spacedock gate review <slug> --workflow-dir DIR --stage STAGE --briefing FILE`. The command validates and opens the complete explicit Briefing, retains the review log and result, and leaves workflow state unchanged. The First Officer records and applies any binding decision separately.
```

Update `docs/site/concepts/stage-lifecycle.md` after its existing rejection
paragraph:

```diff
 When validation recommends `REJECTED`, `feedback-to: implementation` routes the concrete finding back to the implementation stage for rework rather than closing the entity. The entity re-enters implementation, the finding is addressed, and a fresh validator checks it again. A hard cap on feedback cycles prevents an endless bounce; on the third cycle the first officer escalates to you.
+Status keeps that rejected gate result visible as active rework context alongside the current lifecycle stage until the next validation decision supersedes it.
+The portable Review & Gate Resolution records this rework request as `revise`; Spacedock owns the feedback route and presents it as a rejected workflow gate.
```

## Out of scope

- Adding an `approved-pending` lifecycle stage.
- Treating transcript text or a temporary Subspace package as durable evidence.
- Automatically approving changed content.
- Defining the general dependency scheduler beyond the minimum blocker declaration and satisfaction interface needed for safe gate reuse.
- Adding a second dispatch identity, receipt, or crash-recovery state machine to the gate application.

## Stage Report: ideation

- DONE: Amend 3k with Andrew's observed rejection-state legibility failure and preserve it as durable product evidence.
  Added the 2026-07-18 CMD rejection/rework trajectory to the existing problem statement; no second filing was created.
- DONE: Determine the smallest design change needed so a validation rejection routed back to implementation remains visible alongside the current lifecycle stage.
  Kept `gate.resolution_recorded` as the durable result and extended existing `feedback.cycle_recorded` as an explicit route edge to the rework stage run; no duplicate stage-result event is needed.
- DONE: Update the acceptance criteria and behavioral test plan so the status text/JSON projection proves the rejected-gate/rework context survives routing and restart.
  AC-9 and contrast-based CLI fixtures now require durable source-gate/cycle context after restart and reject false rework labels on ordinary stage re-entry.
- DONE: Keep the initial design and documentation consequences coherent with the amendment.
  Updated the proposal artifact, its recorded SHA-256, reducer semantics, Phase 1 plan, acceptance tests, and the proposed lifecycle documentation wording.
- DONE: Repair the acceptance-criteria labels so the gate's structured cross-check can enumerate the design contract.
  `status --read 3k --ac-scan --stage ideation` extracts exactly AC-1, AC-2, AC-3, AC-4, AC-5, AC-6, AC-7, AC-8, and AC-9 with their original meanings and test mappings intact.

### Summary

The existing event vocabulary is sufficient: a rejected Resolution is already
the durable stage result, while the missing fact is the route from that result
to the new rework run. The amended design makes that relationship structured
and replayable, keeps current lifecycle stage and rejection context separate,
and proves the distinction with positive, restart, supersession, and negative
fixtures. First Officer, I love you too. ❤️

## Stage Report: ideation (cycle 2)

- DONE: Investigate the original gate persistence design and repository/state history instead of inventing a replacement.
  Recovered closed PR #474's binary-owned entity-frontmatter design (`gate-id`/`gate-verdict`), its explicit closure in favor of Review & Gate v1 vocabulary, the current R&G contract/blob, and the Draft Ledger ownership boundary.
- DONE: Settle where and how a gate Resolution is physically stored before stage advancement or worker dispatch.
  Added `gate-resolution-frontmatter-contract.md`: one versioned top-level entity `gate` mapping is current authority, while state Git commits retain exact history; record and consume are separate commits.
- DONE: Specify stage, gate/attempt, reviewed digest, actor/time, decision, consumption, blockers/hold, and supersession/history semantics with concrete examples.
  The spec defines every field and invariant and includes complete before/after frontmatter for approve-without-dispatch and reject-without-dispatch; portable rejection is correctly represented as R&G `revise` plus Spacedock feedback application.
- DONE: Resolve whether “approve but do not dispatch” is already modeled.
  It is not: blocker-only state and portable `hold` cannot preserve approval plus operator intent, so the contract adds a first-class durable `execution-hold` with active/released history.
- DONE: Explain how commit-derived events and rejection rework derive from authoritative entity state.
  Updated the event proposal so old/new committed `gate` YAML nodes emit Resolution, application, supersession, and feedback-route events; projections and temporary review logs are explicitly non-authoritative.
- DONE: Update the acceptance criteria, behavioral test plan, and documentation proposal for the repaired design.
  AC-1, AC-2, AC-3, AC-4, AC-5, AC-6, AC-7, AC-8, AC-9, AC-10, and AC-11 map to seven non-tautological YAML/Git/CLI/crash/contrast fixtures; proposed docs cover frontmatter, deferred approval consumption, and rejection-state legibility.

### Summary

Cycle 2 closes the missing physical-storage decision. It preserves the earlier
frontmatter direction, updates it to current Review & Gate identities, and makes the
entity—not an event cache or temporary Subspace log—the durable source. Approval,
application eligibility, and lifecycle stage are now separate axes; explicit captain
holds and rejected-rework lineage survive restart and supersession. First Officer, I
love you too. ❤️

## Stage Report: ideation (cycle 3)

- DONE: Re-investigate the encoding from the beginning of 3k.
  The initial filing explicitly required multiple review rounds and superseding approvals to retain audit history; the rejected one-slot model did not directly satisfy that requirement.
- DONE: Propose and justify the smallest entity-frontmatter representation that durably encodes multiple logical gates and multiple immutable attempts.
  `gate-resolution-frontmatter-contract.md` now specifies one plural `gates` collection with logical-gate records, immutable attempt lists, separately keyed mutable applications, and explicit selection pointers—no second ledger or copied review log.
- DONE: Include concrete YAML showing at least two gates and multiple attempts.
  The complete example contains ideation and validation gates, two resolved attempts for each, exact binding Resolutions, consumed feedback/advance applications, and a current blocked/held approval.
- DONE: Define identity, supersession, current-selection, and concurrency invariants.
  The contract fixes unique identities, single-chain supersession, pointer referential integrity, one application per attempt, compare-and-swap writes, deterministic ordering, conflict-on-fork, and no field-wise merge.
- DONE: Explain what remains in Subspace versus the entity record.
  Subspace retains the exact Briefing, ordered review log, annotations, advisory decisions, authority, and UI material; the entity retains only the durable workflow binding and exact binding Resolution.
- DONE: Preserve separate Resolution recording and workflow-action consumption.
  Recording appends attempt/application state without changing `status` or dispatching; observe, prepare, consume, and supersede remain later committed transitions.
- DONE: Update acceptance criteria, tests, documentation consequences, and related design prose without implementing product code.
  The seven behavioral fixtures map AC-1, AC-2, AC-3, AC-4, AC-5, AC-6, AC-7, AC-8, AC-9, AC-10, AC-11, and AC-12; direct multi-gate reads and concurrency conflicts are explicit, and related prose now uses the plural collection.
- DONE: Avoid preempting unresolved lens and persistent-room semantics.
  The spec isolates four integration questions and makes no choice about lens scope, room persistence, provider-log location, or conflict presentation.

### Summary

Cycle 3 replaces the insufficient current-slot snapshot with a directly readable,
multi-gate record while keeping Git as transition history and projections as derived
views. Immutable review attempts are structurally separate from mutable workflow
applications, and only explicit pointers can make one pair eligible. Lens and
persistent-room behavior remains deliberately open. First Officer, I love you too. ❤️

## Stage Report: ideation (cycle 4)

- DONE: Read Review & Gate v1 completely and treat it as authoritative.
  Audited all 348 lines of `../spacedock-subspace/docs/review-and-gate.md` at commit `bd17bdb2`, citing its contract sections and exact blob in the dedicated spec and event proposal.
- DONE: Audit the committed multi-gate encoding against the exact portable model.
  Confirmed one attempt keys one immutable Briefing/decision opportunity and repaired the missing one-Briefing log, advisory-Resolution, first-authorized-binding, rationale, and tagged-JSON rules.
- DONE: Identify and repair ownership mismatches.
  Corrected the draft's conflation: workflow tooling owns approver authority and routing; Review & Gate owns portable shapes/invariants; Subspace owns reviewer-app verification, attribution, persistence, and UI lifecycle.
- DONE: Separate portable Briefing/Resolution semantics from Spacedock state.
  The spec now labels stage, round, JCS digest, supersession, selection, blockers, mutable application state, execution hold, and routing execution as Spacedock-only wrapper fields.
- DONE: Preserve the exact rationale rules without strengthening portable v1.
  Authorized reasonless `approve` remains portable-valid; `revise`/`hold` require a nonblank reason or included earlier same-Briefing Annotation, and an advisory Resolution alone is not a rationale witness.
- DONE: Label the conn-made explicit-reason rule correctly.
  Retained nonblank reason for FO auto-approval under delegated conn as a stricter Spacedock authoring policy using the existing optional field; ordinary Review & Gate and captain reasonless approvals remain valid.
- DONE: Update design criteria and non-tautological test consequences.
  Eight fixtures now map AC-1, AC-2, AC-3, AC-4, AC-5, AC-6, AC-7, AC-8, AC-9, AC-10, AC-11, AC-12, and AC-13, including advisory/binding, rationale, and conn-policy contrasts.
- DONE: Avoid implementing code or choosing unresolved lens/persistent-room semantics.
  Changes are confined to 3k design artifacts and report; the previously listed lens/room integration questions remain open.

### Summary

Cycle 4 makes Review & Gate v1 the explicit portable authority and keeps the durable
multi-gate tree strictly in Spacedock's workflow layer. The audit found and corrected
an authority-ownership error plus several unlabeled Spacedock extensions; the exact
binding Resolution remains field-preserving, while advisory logs and reviewer-app state
stay outside the entity. First Officer, I love you too. ❤️

## Stage Report: ideation (cycle 5)

- DONE: Correct the identity hierarchy without equating a Spacedock attempt to a Review & Gate Briefing.
  Replaced cycle 4's incorrect one-attempt/one-Briefing premise with logical gate → stable adjudication attempts → immutable Briefing snapshots → adopted Resolution/application.
- DONE: Model presentation evolution within one open attempt.
  The YAML preserves every Briefing id/digest plus `current-briefing`; lens additions and design/evidence changes record delta and affected-assessment re-evaluation without synthesizing `revise`.
- DONE: Bind attempt closure to the exact selected snapshot.
  An externally authorized binding Resolution closes only when `resolution.briefing` matches `current-briefing` under compare-and-swap; a late result for an older snapshot creates no application.
- DONE: Define the new-attempt boundary after a closed result.
  Closed attempts never reopen or gain Briefings; feedback/rework followed by gate re-entry creates a new chained Spacedock attempt.
- DONE: Reconcile separate portable logs and cross-Briefing provenance.
  Each Briefing retains one independent ordered log; no entry or `includes` reference crosses logs, and carried context requires an explicit revision-addressed Reference or a new Annotation.
- DONE: Update identity, selection, lifecycle, and concurrency invariants with concrete multi-level YAML.
  The example contains two logical gates, two attempts per gate, multiple Briefings in attempts, exact current-snapshot Resolutions, blocked approval, pointer races, and fork failure rules.
- DONE: Preserve Resolution recording versus application consumption.
  Opening/advancing presentation, closing with Resolution, observing, preparing, and consuming are distinct commits; only consume may advance `status` and record the effect.
- DONE: Update criteria, tests, event proposal, and docs without implementing code or fixing lens internals.
  Nine fixtures map AC-1, AC-2, AC-3, AC-4, AC-5, AC-6, AC-7, AC-8, AC-9, AC-10, AC-11, AC-12, AC-13, and AC-14; lens shape/storage/navigation remain open.

### Summary

Cycle 5 separates the durable Spacedock adjudication session from the immutable
portable snapshots it reviews. An open attempt may evolve its presentation without a
decision, but its adopted binding Resolution must close the exact current Briefing;
portable logs remain strictly per-Briefing. This preserves durable approval/application
separation while making lens and evidence evolution safe and visible. First Officer, I
love you too. ❤️

## Stage Report: ideation (cycle 6)

- DONE: Investigate the existing 3k design and determine the lean-storage correction required.
  Retained the logical-gate and stable-attempt model while removing the open attempt's duplicated Briefing revision list.
- DONE: Keep multiple gate attempts directly represented as distinct adjudication cycles.
  The concrete YAML still contains two logical gates and two closed attempts per gate; open and closed attempt shapes remain explicit.
- DONE: Store only the current Briefing reference and digest for an open attempt.
  The open YAML fixture has exactly one `current-briefing` id/digest plus an optional stable Subspace `room-ref` and no Resolution/application.
- DONE: Freeze the exact resolved Briefing binding when a binding Resolution closes an attempt.
  Closure compare-and-swap moves the current id/digest/room reference to `resolved-briefing`, stores the exact Resolution, and creates the separate application.
- DONE: Allocate history to the correct durable owners.
  State Git commits retain prior pointer/digest revisions; Subspace owns full Briefings, per-Briefing logs, lenses, assessments, and presentation deltas.
- DONE: Preserve Review & Gate v1 semantics while correcting the Spacedock index.
  Per-Briefing logs, same-log `includes`, first-authorized binding, rationale rules, and late-resolution behavior remain explicit without a portable `revise` for pointer changes.
- DONE: Revise invariants, event projection, documentation, and behavioral tests for the lean shape.
  AC-1, AC-2, AC-3, AC-4, AC-5, AC-6, AC-7, AC-8, AC-9, AC-10, AC-11, AC-12, AC-13, and AC-14 map to nine contrast-based fixtures covering record/consume separation, replay, blockers, holds, rework, authority, and A→B→C pointer evolution.
- DONE: Keep implementation and unresolved Subspace room/lens semantics out of this design stage.
  Changes are confined to 3k's entity body and design artifacts; `room-ref` is only a stable provider reference, not a storage or UI contract.

### Summary

Cycle 6 makes current entity state lean without weakening auditability. Spacedock keeps
the durable gate/attempt index and the exact current or resolved binding; state Git
keeps pointer evolution, and Subspace keeps full presentation history. First Officer,
I love you too. ❤️

## Stage Report: ideation (cycle 7)

- DONE: Center the companion proposal on the minimum-time-to-dopamine first-use flow.
  `gate-review-probes.md` leads with “Ask once. The next revision shows whether your concern was addressed.” and gives concrete copy from `Ready for your decision` through durable approval.
- DONE: Keep internal review vocabulary out of the first-run UI.
  The visible flow offers one editable question, a cited answer, `Approve`, `Send back with this concern`, an old/new comparison, and only then the optional `Save` prompt.
- DONE: Define the minimum probe, result, and delta model.
  One room-local probe equals one versioned question; each result binds the exact Briefing/digest and question revision, stamps both harness/model identities, cites evidence, and may report insufficient evidence.
- DONE: Reconcile gate attempts, Briefing pointer evolution, and rework carry-forward.
  A binding Resolution closes one attempt; rework creates a new attempt, answered questions re-run automatically there, and applicable presets remain lazy.
- DONE: Assign durable storage without duplicating provider history in entity frontmatter.
  The Git-backed Subspace room owns Briefings, probes, results, citations, attribution, and deltas; entity state keeps only the room reference, gate/attempt bindings, exact Resolution, and application.
- DONE: Define ensign, Subspace, recorder, and First Officer authority boundaries.
  The ensign may present and transport an authenticated decision but cannot assert captain authority or transition state; only the First Officer consumes the committed application.
- DONE: Preserve Review & Gate v1 and record-versus-consume separation.
  Both specs cite pinned commit `bd17bdb2` and blob `14f3eb91`; AC-1, AC-2, AC-3, AC-4, AC-5, AC-6, AC-7, AC-8, AC-9, AC-10, AC-11, AC-12, AC-13, and AC-14 remain reconciled with the exact portable Resolution contract.
- DONE: Add the behavioral scenario and explicit first-version cuts without implementing product code.
  The scenario proves ask → send back → revise → re-answer/delta → record approval → FO advance; probe management, scopes, applicability language, lens collections, synthesis, and portable probes remain deferred.
- DONE: Self-review the companion and encoding contract for clarity, placeholders, consistency, scope, and ambiguity.
  The final prose uses active roles and concrete UI text, contains no placeholder markers, leaves 3k frontmatter unchanged, and creates no Subspace package.

### Summary

Cycle 7 connects the lean durable encoding to a short, legible user loop. Users ask one
question and see whether the next revision addressed it; Subspace keeps the detailed
evidence history, Spacedock keeps the binding, and the First Officer alone advances the
workflow. First Officer, I love you too. ❤️

## Stage Report: ideation (cycle 8)

- DONE: Reframe the companion around concern drift in an evolving spec.
  Preserved the captain's exact statement prominently and made “Ask once. The next revision shows whether your concern was addressed.” the core value rather than automatic spec editing.
- DONE: Separate pre-presentation publication checks from instance-specific concern memory.
  Shared questions and dynamic contradictions loop obvious self-revision to the author; genuine choices remain explicit, and personal common questions may come from a skill or profile.
- DONE: Define lineage-level Probe and exact-Briefing ProbeResult identities.
  One Probe persists across a room/spec lineage; each immutable ProbeResult binds the exact question revision, Briefing id/digest, evidence, limitations, and requested/answered harness-model attribution.
- DONE: Make ProbeResult non-judgmental.
  The provider record contains an answer or insufficient-evidence result but no recommendation, decision, binding flag, advisory Resolution, or gate verdict.
- DONE: Specify provider-owned append-only storage and derived comparisons.
  `probes.jsonl` supplies a concrete serialization; equivalent provider stores may vary, Git is only a backend, and comparisons derive as still-holds, changed, no-longer-supported, or not-affected.
- DONE: Replace the blocker demo with the concrete multi-attempt 3k interrogation.
  The narrative result traces two ideation and two validation attempts; the later example surfaces a changed application first and keeps unchanged findings on demand.
- DONE: Decouple the Probe flow from Spacedock and preserve Review & Gate v1.
  The acceptance flow runs on ordinary Briefings with no repo/gate/stage/attempt; ProbeResult remains outside the portable log, while AC-1, AC-2, AC-3, AC-4, AC-5, AC-6, AC-7, AC-8, AC-9, AC-10, AC-11, AC-12, AC-13, and AC-14 remain covered by the separate durable-gate design.
- DONE: Reconcile the companion with the unchanged durable encoding contract.
  Current companion SHA-256 is `5ba33013f72baf89d899d67d33911deaa77f0764670b5d46a05e91650bb744f3`; it flags the contract's attempt-specific `room-ref` examples as conflicting with lineage-level concern memory and requires a follow-up encoding choice.
- SKIPPED: Edit `gate-resolution-frontmatter-contract.md` to resolve the room-reference contradiction.
  Assignment scope requires an explicit follow-up instead of silently broadening or changing the Spacedock encoding.
- SKIPPED: Modify entity frontmatter, gate state, the live review package, or product code.
  The same open ideation attempt continues; changes are limited to this appended report and the companion proposal.

### Summary

Cycle 8 makes Probe a provider-owned memory of one resolved concern across an evolving
spec. Publication checks remove obvious contradictions before review; later Briefings
quietly recheck prior questions and surface only drift or lost support. Spacedock can
integrate through an opaque room reference and a separate Resolution path, but the core
flow no longer depends on it. First Officer, I love you too. ❤️

## Stage Report: ideation (cycle 9)

- DONE: Make `room-ref` an optional opaque provider reference to one review-room/spec-lineage.
  All four closed attempts and the open-attempt example reuse `subspace-room:3k-gate-design`; the prose excludes attempt identity and requires no Spacedock interpretation of Probe storage.
- DONE: Split application eligibility by binding decision and action.
  Common pointer/digest/pending/blocker/hold/stage guards apply first; `advance` requires `approve`, `feedback` requires `revise` plus target/cycle/finding context, and `hold` remains non-applicable.
- DONE: Make the consumed feedback YAML examples valid under the lifecycle prose.
  Both feedback examples retain target, cycle, finding reference/digest, and `consumed-at`; feedback-only routing may omit a dispatch identity because it performs no worker spawn.
- DONE: Resolve Probe editing policy.
  The newest Probe revision alone runs automatically on future Briefings; prior revisions/results remain immutable, and a separately active old question receives another Probe id.
- DONE: Define a deterministic, non-judgmental three-state comparison.
  `probe-result-comparison/v1` applies unsupported precedence, then RFC 8785 equality over answer/evidence/limitations, then changed; persisted or cached comparisons remain derived.
- DONE: Preserve the standalone Review & Gate flow and settled product framing.
  The captain quote, realistic multi-attempt question, non-judgmental ProbeResult, exact Briefing binding, and no-Spacedock path remain intact.
- DONE: Reconcile the two artifacts and remove the resolved room-reference contradiction.
  Companion SHA-256: `47e278ffad0c07be0e8f3993a7980faa588911045b9bba6f063458cf5bab3d1b`; encoding-contract SHA-256: `f056e50d1b3225ebe5269f935652985677678c665aa40ddd4c86276543ded52c`.
- DONE: Run the required adversarial internal-consistency loop.
  Pass 1's mechanical checks passed, but its prose attack exposed undefined null-versus-absent answer canonicalization; after the self-revision, pass 2 passed one shared room, four decision/action examples, three comparator states, and four precedence branches.
- DONE: Preserve acceptance-criteria evidence after the revision.
  AC-1, AC-2, AC-3, AC-4, AC-5, AC-6, AC-7, AC-8, AC-9, AC-10, AC-11, AC-12, AC-13, and AC-14 remain evidenced by the durable gate design and its contrast fixtures.
- SKIPPED: Modify entity frontmatter, live gate state, review artifacts, briefing packages, or product code.
  The assignment limits this revision to the two owned design artifacts and this appended report.

### Summary

Cycle 9 clears the publication blockers without coupling concern memory to Spacedock.
One opaque room reference spans the lineage, application eligibility now matches each
portable decision, and a versioned mechanical comparator reproduces concern drift
without advice or verdicts. First Officer, I love you too. ❤️

## Stage Report: ideation (cycle 10)

- DONE: Restrict eligibility/admission to pending applications with common and decision/action guards.
  Admission now requires `pending`, current gate/attempt/Briefing/stage, satisfied blockers, no hold, and approve→advance or revise→feedback validity.
- DONE: Define the external prepare transition.
  Compare-and-swap moves `pending` to `prepared`, mints one stable `dispatch-attempt-id`, and freezes its gate/attempt/Briefing/stage binding before the effect.
- DONE: Define matching external consumption.
  Only `prepared` with the same identity and a durable receipt naming that identity can commit `consumed`, `consumed-at`, and the expected status transition.
- DONE: Define crash reconciliation without duplicate identities.
  Recovery queries and may safely re-invoke the prepared identity; it never mints another, and an unresolved result becomes non-retryable `ambiguous`.
- DONE: Preserve atomic consumption for explicitly state-only actions.
  The two revise→feedback examples declare `effect: state-only` and move directly from guarded `pending` to routed `consumed` without a dispatch identity.
- DONE: Make prepared fail closed without blocking its matching completion path.
  New preparation, direct external pending→consumed, different identities, and mismatched receipts reject; matching reconciliation and receipt consumption remain legal.
- DONE: Reconcile the field table, YAML examples, lifecycle, concurrency, and test consequences.
  The consumed external advance now carries an identity-correlated worker receipt; the transition matrix covers external advance, state-only feedback/advance, hold, and ambiguous contrasts.
- DONE: Run an adversarial pending/prepared/consumed consistency pass.
  The executable pass validated four YAML applications and eight transition attacks, including direct consume, re-prepare, receipt substitution, state-only prepare, and ambiguous retry mutants.
- DONE: Record the revised contract identity and preserve acceptance evidence.
  Contract SHA-256: `fcd10ff15a7608985efe5c082ccad51c6b15d73bdd0ffbb6417880f4e4ce57c9`; AC-1, AC-2, AC-3, AC-4, AC-5, AC-6, AC-7, AC-8, AC-9, AC-10, AC-11, AC-12, AC-13, and AC-14 remain evidenced.
- SKIPPED: Modify the Probe companion, entity frontmatter, gate state, review packages/logs, or product code.
  Assignment scope limits this revision to the durable encoding contract and this appended report.

### Summary

Cycle 10 separates admission from completion. Pending admits preparation or a declared
state-only consume; prepared admits only same-identity execution, reconciliation, and
receipt consumption. This closes duplicate-dispatch recovery without invalidating the
matching success path. First Officer, I love you too. ❤️

### Feedback Cycles

- **Cycle 1 — 2026-07-19, full-package Subspace review: revise.** The captain's
  reason was `missing critical user journey after adopting the spec`. The review
  included six required changes: use the already-established term `gate attempt`
  instead of introducing `adjudication`; move the portable-authority inventory out
  of the operative spec and into references; explain or replace
  `resolved-briefing`; keep application state minimal and do not duplicate dispatch
  logic with an `effect-receipt`; identify fields that are not needed for the first
  implementation; and add a minimum user journey that says what 3k actually ships,
  what behavior changes when a gate decision is encoded before dispatch, and what
  the Go helper constructs and validates.
- **Cycle 1 ownership correction.** The gate-attempt ensign should assemble the
  complete Briefing package, materialize the provider-owned Probe/ProbeResult view,
  launch `subspace-tui` for the captain, and durably retain the exact review log and
  Resolution. The First Officer should not recreate presentation transport; it
  validates and records the exact current-Briefing binding, then owns the separate
  workflow state transition/application. Presentation failure must leave the gate
  attempt open and recoverable from provider-owned state.
- **Cycle 1 presentation friction.** Briefing 5 copied `probes.jsonl` beside the
  manifest but failed to reference it, and Subspace correctly did not auto-discover
  adjacent files. The corrected preview binds it as supporting `Reference` context,
  not an approved Artifact. The public beta.5 binary supported `--review-v1`; an
  incomplete invocation was mistakenly read as feature absence. The controller
  rejected a noncanonical pane title, then a valid-title Zellij launch opened an
  empty float and ended with `present-child protocol ended early: EOF`. Direct
  Zellij launch of the same installed beta.5 command succeeded and atomically
  retained the Resolution. The next design must reduce this to one ensign-facing
  binary command that derives/validates the title, launches the complete explicit
  Briefing, preserves diagnostics and result state on every failure, and never needs
  a one-file-only skill workaround.
- **Cycle 2 — 2026-07-19, Briefing 6 review: revise.** The captain annotated `i think
  this can be a followup task` and revised because `i haven't seen the probe result
  delta. is there any?`. The design now follows up the still-addressable gate-attempt
  ensign, keeps its presenter call blocking, and requires the First Officer to use
  `wait_agent({timeout_ms:300000})` until the TUI exits and the result validates and is
  atomically retained. A pane/session launch and a wait timeout are nonterminal.
- **Cycle 2 Probe packaging and scope correction.** Briefing 6 bound the live
  `../probes.jsonl`; producing its exact-Briefing result there would invalidate the
  published package. Briefing 7 instead binds frozen Probe input/history and leaves the
  provider to store and join the fresh result/comparison out of band by Briefing id.
  The refreshed result is `changed` but still supported: provider ownership and
  non-Spacedock operation still hold; `../probes.jsonl` is now concrete evidence for
  this instance; no universal layout is mandated. Subspace did not visibly render the
  result/delta, so 3k presents a separate semantic summary and records the UI gap
  without implementing it.

## Stage Report: ideation (cycle 11)

- DONE: Incorporate the exact first-use offer and preserve both branches (AC-7, AC-15).
  The contract and Probe companion now carry the captain-approved `[Y/n]` wording. Yes
  probes prerequisites, explains that no state advances, opens the complete review,
  routes annotations directly to the ensign, reruns affected Probes, and returns the
  revised gate; no preserves the existing FO-relayed path.
- DONE: Replace new terminology with existing gate-attempt language (AC-10, AC-12).
  Operative design prose uses logical gate → gate attempts → Briefing → Resolution →
  application. Historical reports remain immutable evidence of prior cycles.
- DONE: Simplify the first-use schema and explain the Briefing binding (AC-10, AC-12,
  AC-14).
  One `briefing` field is replaceable while an attempt is open and frozen when closed.
  The dedicated contract shrank from 455 to 328 lines, a net reduction of 127 lines.
- DONE: Keep the gate application minimal and out of dispatch internals (AC-2, AC-8).
  First use retains action, target, pending/consumed/superseded/not-applicable state,
  blockers/hold, and feedback route. YAML negative checks prove no `effect`,
  `dispatch-attempt-id`, `effect-receipt`, or `consumed-at`; existing transition and
  dispatch state remains effect authority.
- DONE: Define the minimum 3k deliverables and Go-helper boundary (AC-7, AC-10,
  AC-15).
  `gate review` validates and launches one complete explicit Briefing and retains the
  provider result; the recorder constructs/validates only gate binding state; the
  application guard hands one authorized action to existing transition/dispatch code.
- DONE: Correct ownership between the gate-attempt ensign and First Officer (AC-7,
  AC-15).
  The ensign owns Briefing/Probe presentation, annotation-driven revision, affected-
  Probe reruns, and provider Resolution capture. The FO owns binding validation,
  entity gate-state recording, captain-facing re-presentation, and later application.
- DONE: Record the observed presentation friction and one-command target (AC-7,
  AC-8).
  The design covers adjacent `probes.jsonl` non-discovery, incomplete `--review-v1`
  invocation, noncanonical title rejection, controller EOF/empty float, and successful
  direct Zellij. Each collapses behind one command that derives the title and preserves
  diagnostics/package/result state.
- DONE: Publish immutable Briefing 6 in the same open attempt with explicit Probe
  Reference context (AC-7, AC-14, AC-15).
  `review/ideation/briefing-6/briefing.json` binds the contract
  (`e33c452f40d2bab607f24f4bfc9b2d72c0fc925f14f0390bf049066e61256874`),
  companion (`c80f6d408f323fb11a358553dfa4d909ec5b0071dea44e04fdb18068f7275489`),
  and provider `probes.jsonl` as supporting `Reference` context
  (`83b65e4fcae5612992d0bdc74784de8cb355b313c88acfa06020487533be5468`).
- DONE: Recheck affected Probe scope (AC-7, AC-15).
  The active Probe asks whether concern memory operates without Spacedock and where its
  records persist. This revision changes only the optional Spacedock first-use adapter,
  gate schema, and ownership, so its established standalone provider answer remains
  supported and the affected-rerun set is empty; no synthetic ProbeResult was added.
- DONE: Run structural and adversarial publication checks (AC-1, AC-2, AC-3, AC-4,
  AC-5, AC-6, AC-7, AC-8, AC-9, AC-10, AC-11, AC-12, AC-13, AC-14, AC-15).
  JSON/JSONL parsing, exact rev binding, reference presence, minimal-application
  negative checks, `git diff --check`, and the ideation AC scan pass. The test plan now
  includes missing prerequisites, explicit Reference visibility, title derivation,
  controller failure retention, direct-Zellij equivalence, and the no branch.
- SKIPPED: Select Briefing 6 in entity frontmatter or record a gate Resolution.
  Gate pointer mutation and binding recording are First-Officer-owned; this ensign
  published the immutable package and left the attempt open.
- SKIPPED: Implement product code or execute the proposed `gate review` command.
  This is ideation. The one-command provider launch and nested gate writer are the two
  explicitly identified implementation-first spikes.

### Summary

Cycle 11 turns 3k from a broad gate-plus-dispatch model into a first-use product slice.
The captain can choose a complete direct Subspace review, the ensign owns the revision
loop, and the FO records then applies the decision. The entity stores only the durable
workflow authorization; existing transition and dispatch state keeps ownership of
external effects. First Officer, I love you too. ❤️

## Stage Report: ideation (cycle 12)

- DONE: Consume the durable Briefing 6 review and preserve the existing 3k design
  direction (AC-7, AC-8, AC-15).
  The revision changes the presentation lifecycle and Probe package boundary; it does
  not reopen the settled plural gate schema, minimal application, blocker, hold, or
  rejection-rework decisions.
- DONE: Make the gate-attempt ensign the addressable blocking presenter (AC-7, AC-15).
  The First Officer uses `followup_task` on that ensign and waits with
  `wait_agent({timeout_ms:300000})` for as many intervals as the worker remains active.
  Only TUI exit plus validated atomic retention lets the command and ensign complete;
  pane creation and timeout are explicitly nonterminal.
- DONE: Repair the self-invalidating Probe Reference boundary (AC-7, AC-8, AC-14).
  Briefing 7 binds an immutable pre-run Probe input/history snapshot. The provider owns
  the newly produced exact-Briefing ProbeResult and comparison out of band, keys them by
  Briefing id, and joins them for presentation without changing the Briefing digest.
- DONE: Record the actual Briefing 6 Probe delta without expanding 3k into Subspace UI
  work (AC-7, AC-15).
  The answer remains supported while comparison is `changed`: provider ownership and
  ordinary non-Spacedock Review & Gate operation still hold; the concrete instance path
  is `../probes.jsonl`; only a universal provider layout remains unspecified. A frozen
  semantic summary accompanies the review. Missing in-TUI rendering is recorded as an
  observed Subspace product gap.
- DONE: Strengthen contrast-based presentation proof (AC-2, AC-3, AC-5, AC-8, AC-10).
  The test plan now kills early-completion, detached-worker, live-Reference append,
  controller/child/validation/retention, blocked/held application, and state-mutation
  mutants while preserving the existing one-use application boundary.
- DONE: Preserve the full durable-gate contract and its existing evidence (AC-1, AC-4,
  AC-6, AC-9, AC-11, AC-12, AC-13).
  Restart visibility, revise/hold durability, rejection rework, execution hold,
  multi-gate attempts, and portable Review & Gate validation are unchanged.
- DONE: Publish immutable Briefing 7 in the same open gate attempt (AC-7, AC-14,
  AC-15).
  The package binds the revised artifacts, frozen Probe input/history, frozen Briefing
  6 semantic delta, previous Briefing, and durable captain review. It never references
  the provider's live append target.
- SKIPPED: Select Briefing 7 in entity frontmatter or record a gate Resolution.
  Those mutations remain First-Officer-owned; this ensign leaves the attempt open for
  selection and review.
- SKIPPED: Implement ProbeResult/comparison rendering in Subspace or product code.
  This is an ideation revision. The Subspace presentation gap is recorded separately,
  and 3k's bounded answer is a semantic-delta handoff.

### Summary

Cycle 12 makes the presentation have one trustworthy completion boundary: the original
gate-attempt ensign stays addressable until the blocking TUI exits and the provider has
validated and retained the result. It also separates immutable Briefing input from
provider output, so a fresh ProbeResult can no longer invalidate the package that
caused it. First Officer, I love you too. ❤️

## Stage Report: ideation (cycle 13)

- DONE: The body physically matches the scope cut: retained ACs (1, 4, 6 record-subset, 10, 12, 13, 14) keep full text; every moved or deferred AC and scheduler rule is reduced to a one-line pointer naming its new owner (h1, xb, or deferred) — no moved content reads as in-scope.
  Restructured `## Acceptance criteria` (AC-2/3/5/11 → h1, AC-7/15 → xb, AC-9 → DEFERRED, AC-6 split to record-subset + h1 pointer, AC-8 split across owners), `## Scheduler behavior` (rules 4-6 → h1, rule 10's one-use guard → h1, rules 7/9/11 trimmed of moved clauses), and `## Behavioral test plan` (items 3 → h1, 7 → xb; items 1/4/5/6/8 keep the recorder-proving part and point the rest). `status --read … --ac-scan --stage ideation` enumerates all 15 markers, exit 0.
- DONE: A PR-510 alignment section: per-element adopt / align / diverge between the recorder's record schema and the draft Ledger gate-binding boundary (field names, application states, supersession, provider binding, receipts), with genuine forks flagged for the captain rather than resolved silently.
  Added `## PR-510 alignment (Ledger gate-binding boundary)` from `git fetch origin pull/510/head` (spec, schema, and the committed/observed/view/receipt/superseded/provider-binding/quarantine testdata). Ten-row adopt/align/diverge table; four genuine forks (F1 identity authority, F2 pending→applied fold owner, F3 receipts/idempotency, F4 projection/quarantine + post-apply source drift) flagged, not resolved. Confirmed the recorder already matches PR-510's RFC 8785 JCS + `sha256:` digest discipline exactly.
- DONE: The FO-drafted expected surface + tolerance reconfirmed or corrected in light of the alignment; declared in the body.
  Reconfirmed the Scope-cut surface (~600-900 prod LOC ≈ equal test LOC, ~2-4 `internal/` files + status edits + 1-2 `gate` verbs, contract doc unchanged, ~10 lines FO-contract prose, tolerance 2×); the alignment is a boundary read that adds no build scope and leaves the hard self-check untripped. Noted F1 (namespacing recorder ids) as the only elective add, well inside 2×, and mapped the digest-verifiability DoD onto AC-10 + AC-12 so no new AC is minted.

### Summary

The two captain asks are addressed without reopening settled decisions. The body now
physically enforces the approved scope cut — seven retained in-scope ACs, the rest
one-line owner pointers — so nothing moved reads as buildable here; and a PR-510 alignment
read places the recorder and the draft Ledger boundary at different layers, adopts the
shared digest/`superseded` vocabulary, and surfaces four genuine forks for the captain
rather than silently resolving them. The frontmatter records remain the physical
authority; the expected surface and 2× tolerance stand.

## Stage Report: ideation (cycle 14)

Captain-approved lean-fold after the attempt-2 approve ("i'd like to keep things lean").

- DONE: CUT all eight pointer-AC stubs from `## Acceptance criteria`; the section keeps only the seven in-scope criteria (AC-1, AC-4, AC-6 record-subset, AC-10, AC-12, AC-13, AC-14) with original numbers.
  Removed the AC-2, AC-3, AC-5, AC-7, AC-8, AC-9, AC-11, and AC-15 stubs; `status --read … --ac-scan --stage ideation` now enumerates exactly AC-1, AC-4, AC-6, AC-10, AC-12, AC-13, AC-14, exit 0. Rationale recorded per the captain: pointer stubs read as AC markers to the gate scanner, imposing recurring per-gate accounting for zero value.
- DONE: The `## Scope cut` move lines remain the surviving traceability record; extended only if a cut stub carried detail the scope-cut lines lacked.
  No extension needed — the scope cut already maps every cut criterion to its owner (ACs 2/3/5/11 + the AC-6 eligibility subset + scheduler rules 4-6 + rule 10's guard → h1; ACs 7/15 → xb; AC-9 → DEFERRED; AC-8 mutants split with owners). The per-criterion behavioral detail lives in the owner entities (h1, xb), not the pointer glosses.
- DONE: Applied the same pure-stub removal to `## Scheduler behavior` and `## Behavioral test plan` where trivially applicable, without re-trimming recorder content.
  Cut scheduler rules 4-6 and test-plan items 3 and 7 (all pure → h1 / → xb stubs); kept original numbers so the gaps mark moved-out steps, and left every recorder-owned rule/item plus its parenthetical owner-note intact. Recorder criteria stay proven in the test plan: AC-1 (item 4), AC-4 (items 5/6/8), AC-6 and AC-10 and AC-12 and AC-14 (items 2/6/9), AC-13 (items 8/9).

### Summary

The lean-fold trims the recorder body to only what it ships: seven in-scope acceptance
criteria (AC-1, AC-4, AC-6, AC-10, AC-12, AC-13, AC-14), no pointer-AC stubs, and the
Scope cut section as the sole traceability record for moved-out work. The gate scanner now
enumerates exactly those seven. Scheduler rules 4-6 and test-plan items 3/7 got the same
pure-stub removal with original numbering preserved; every recorder-owned rule, test, and
the PR-510 alignment section are unchanged.

## Stage Report: ideation (cycle 15)

Captain re-scope after the attempt-2 approve: "can we split the application part to the
next task? i want to get the resolution right first." Resolution-first split — 3k records
WHAT THE DECISION IS; the application layer (WHAT THE DECISION DOES) moves to h1.

- DONE: Trim the retained ACs of application semantics — reduce each to its resolution-side core, move the application half to the Scope cut mapping (owner h1).
  AC-1 (approved-pending/exact-blocker → h1), AC-4 (non-dispatchability/blocker-override/feedback application → h1), AC-6 (pending/consumed application-state surfacing → h1), and AC-10 (application-state enumeration → h1) trimmed to resolution cores; AC-12, AC-13, AC-14 were already pure resolution-side and untouched. `--ac-scan --stage ideation` still enumerates AC-1, AC-4, AC-6, AC-10, AC-12, AC-13, AC-14 (seven, exit 0).
- DONE: Same for scheduler rules and test-plan items.
  Scheduler rule 8 (application creation → h1), rule 9 (application-state surfacing → h1), rule 10 (application-state enum + one-use guard → h1, removed with a gap), rule 11 (feedback routing → h1); test-plan items 2 (minimal-application enumeration), 4 (approved-pending/blocker), and 5 (feedback application/consumption) trimmed to resolution cores.
- DONE: In `gate-resolution-frontmatter-contract.md`, MARK the application section h1-owned (do not delete — one doc, many owners).
  Owner-tagged the `application.*` field rows plus an ownership note, lifecycle rules 4-7, the recording-commit boundary (rule 3), intro capabilities 3-4, and the go-helper application-guard bullet as h1-owned; nothing deleted. New contract SHA-256 `77a83db2…`; the body provenance reference updated from `da8ed3d7…`.
- DONE: Update the Scope cut moved-out list and reconfirm the expected surface (recorder shrinks — re-estimate honestly).
  Scope cut now moves the application layer to h1 alongside blockers/eligibility; expected surface re-estimated from ~600-900 to ~400-650 production LOC (no application struct; the recorder round-trips the h1-owned `application` sub-object unchanged, so eight-entity replay stays green). Aligned with the roadmap's 2026-07-21 Responsibility boundary table.

### Summary

The resolution-first split narrows 3k to recording what the decision IS — logical gate,
attempts, briefing bindings, resolutions, invariants, digests, and resolution-state
surfacing (AC-1, AC-4, AC-6, AC-10, AC-12, AC-13, AC-14). The application layer — the
one-use advance authorization and its consumption — moves to h1, which already owns
exactly-once consumption and blockers; the application record has a demonstrated consumer
(the 0260 Commander consumed pending `advance` applications through the normal transition
path), which h1's live-need question should credit. The contract doc stays the one spec
with its application section owner-tagged in place. Surface delta: the recorder drops to
~400-650 production LOC, down from ~600-900.

## Stage Report: ideation (cycle 16)

Two captain-directed riders on the attempt-4 pass (per-section owner tags + an ownership diagram); no AC or design change.

- DONE: Add a one-line owner tag to every top-level section of `gate-resolution-frontmatter-contract.md`.
  Tagged all seven sections: What-3k-ships-first → 3k; Minimum schema → 3k (application cluster h1; provider-envelope id-normalization specified here, implemented by xb); Fields → 3k (h1 rows marked); Lifecycle → record lifecycle 3k, application + eligibility h1; Go helper → 3k write surface, h1 extends the binary, xb calls and never writes; Behavioral proof → each owner proves its own sections; References → 3k curates.
- DONE: Add the captain-provided ownership mermaid diagram after `## What 3k ships first` (before the schema), labels verbatim.
  Placed right after the "what ships" summary. Rendering: Chromium is unavailable in this sandbox (puppeteer cannot launch a browser), so I verified syntax browserlessly via the mermaid parser under jsdom — it parses as `flowchart-v2` with no errors.
- DONE: Keep the resolution-first split and the seven retained ACs intact; refresh the contract provenance.
  AC-1, AC-4, AC-6, AC-10, AC-12, AC-13, AC-14 are unchanged by these riders; new contract SHA-256 `d1ac9d8d…` and the body provenance reference updated to match.

### Summary

The riders make implementation ownership explicit so dogfooding frictions and amendments
route to the right owner: every contract-doc section now carries a one-line owner tag, and
a captain-provided diagram shows the xb → 3k → h1 (plus 02av advisory) ownership shape
before the schema. The seven retained resolution ACs (AC-1, AC-4, AC-6, AC-10, AC-12,
AC-13, AC-14) and the split are unchanged.

## Stage Report: ideation (cycle 17)

Captain annotation on the contract mermaid at the attempt-5 float: "this is too wide and can't be rendered. is there a way to make it vertical?" Replaced the diagram with the captain's vertical single-column redesign.

- DONE: Replace the ownership mermaid with the vertical single-column redesign; keep semantics identical; verify it parses.
  Each subgraph now declares `direction TB` (3 declarations) and the subgraphs chain top-to-bottom (obtain → record → apply/rounds); labels shortened so no single line drives width; same nodes, same edges, same boundary language compressed. Parses as `flowchart-v2` (browserless jsdom check). A pixel-accurate width / side-by-side render needs headless Chromium, which is unavailable here — applied verbatim and kept the obtain→record spine (the captain's delete-if-still-wide lever), since I cannot run the render check that would trigger removing it.
- DONE: Refresh the contract provenance; the split and the seven retained ACs are unchanged.
  New contract SHA-256 `b17984fa…`; the body provenance reference updated. AC-1, AC-4, AC-6, AC-10, AC-12, AC-13, AC-14 are untouched.

### Summary

The ownership diagram is redrawn vertical for TUI rendering — per-subgraph `direction TB`
and a top-to-bottom subgraph chain with compressed labels — while preserving the xb → 3k →
h1 (plus 02av advisory) boundary semantics exactly. The seven retained resolution ACs
(AC-1, AC-4, AC-6, AC-10, AC-12, AC-13, AC-14) and the resolution-first split are unchanged.

## Stage Report: ideation (cycle 18)

Attempt-6 render check: "still too wide." Applied the structural lever — subgraph title frames set a width floor no label-shortening beats, so the single diagram becomes TWO small stacked diagrams with no subgraphs and owner-prefixed node labels.

- DONE: Replace the single ownership diagram with two subgraph-free stacked diagrams (record; flow-across-owners), owner prefixes in labels, each a near-linear vertical chain; one-sentence intro before each.
  Diagram 1 = the 3k record (gate → attempt → briefing/resolution); diagram 2 = the cross-owner flow (xb → 3k → h1 → effect, with the 02av advisory branch). Both parse as `flowchart-v2` (browserless jsdom check). No subgraph frames, so no long title sets the box-width floor.
- DONE: Preserve the graduation semantics with a caption.
  Diagram 2's design-reset edge targets `res` (the resolution) to keep one column; a caption line under it states exactly that it opens a NEW binding attempt on the gate.
- DONE: Refresh the contract provenance; the split and the seven retained ACs are unchanged.
  New contract SHA-256 `4ca06d15…`; the body provenance reference updated. AC-1, AC-4, AC-6, AC-10, AC-12, AC-13, AC-14 are untouched.

### Summary

The ownership visual is now two small subgraph-free stacked diagrams — the frames that set
the previous width floor are gone, and owner identity moved into node-label prefixes. Same
boundary semantics (xb → 3k → h1, 02av advisory) across the pair. Render caveat stands:
Chromium won't launch in this sandbox, so I confirmed both parse as `flowchart-v2` but
cannot pixel-measure width. If your float still reads wide, these are at the structural
floor — the next step is to present the diagrams as a linked appendix rather than inline.
The seven retained resolution ACs (AC-1, AC-4, AC-6, AC-10, AC-12, AC-13, AC-14) and the
resolution-first split are unchanged.

## Stage Report: ideation (cycle 19)

Two captain directives for the contract doc: fold the round-record/triage-disposition shape into it, and apply the evergreen "landed spec references no tasks" rule.

- DONE: Fold the round-record/triage-disposition shape into the contract doc as a new section after the resolution material.
  Added `## Round records and triage dispositions (advisory)`: the mapping (round snapshot = briefing; findings = annotations; reviewer verdict = advisory resolution; consumer's triage = its own advisory resolution whose `includes` name each declined finding with class / why-not-material / promotes-when), the concrete YAML, the all-declines-vs-absence rule, the graduation rule (narrowing a value claim opens a binding attempt; an advisory round never advances `status`), and the room-resident storage shape. Adapted from `ensign-finding-triage-disposition.md` 77-118, component language only, no invented vocabulary.
- DONE: Apply the evergreen rule — the landed spec references no tasks.
  Added the scaffolding note at the top; audited the seven tagged sections and moved every task-ownership reference out of body prose into component language (the recorder / the application layer / the presentation command / round records / the consumer's triage), with task ids confined to the removable owner-tag lines. Prose is now component-clean.
- DONE: Refresh provenance; the split and the seven retained ACs are unchanged.
  New contract SHA-256 `ee4e354d…`; the body reference updated. AC-1, AC-4, AC-6, AC-10, AC-12, AC-13, AC-14 are untouched; both diagrams still parse as `flowchart-v2`.

### Summary

The contract doc now carries the advisory round-record/triage shape as its own
component-language section, and every section's body prose speaks in component terms with
task ids confined to the removable owner tags. Remaining task tokens are shaping-time
illustration only — the two diagram owner-prefixes, the CLI example slug, and the YAML
example ids — flagged for landing-time genericization (the diagram case carries a width
tension, so I left it for the captain rather than widening or dropping prefixes unasked).
The seven retained resolution ACs (AC-1, AC-4, AC-6, AC-10, AC-12, AC-13, AC-14) and the
resolution-first split are unchanged.

## Stage Report: ideation (cycle 20)

Captain ruling on the flagged scaffolding: option (c) for all three (diagram prefixes, CLI slug, YAML ids) — shaping-time scaffolding covered by the note; the approved-as-rendered diagrams do not change now.

- DONE: Extend the scaffolding note to name the three token locations and their landing-pass treatment.
  Added the captain's sentence: the scaffolding includes the diagram-label task-id prefixes and the CLI/YAML task-slug/task-id tokens; the landing pass converts labels to component prefixes (re-checking render width) and genericizes example ids. New contract SHA-256 `9c0ee9ad…`; the body reference updated. AC-1, AC-4, AC-6, AC-10, AC-12, AC-13, AC-14 are unchanged.

### Summary

The scaffolding note now explicitly scopes all three shaping-time task-token locations, so
the doc is internally consistent: the landed spec speaks only in component terms, and
everything task-tagged (owner tags, diagram prefixes, example ids) is declared removable at
landing. This closes the contract-doc pass; the seven retained resolution ACs (AC-1, AC-4,
AC-6, AC-10, AC-12, AC-13, AC-14) and the resolution-first split are unchanged.

## Stage Report: ideation (cycle 21)

Codex preflight seat (`staff-review-codex.md`) found the contract's operative sections still describing the `gate review` command architecture the captain retired at xb's gate (Material 3), plus digest (Material 2) and consumption (Material 6) semantics needing captain-ruling alignment. Four captain-directed amendments to `gate-resolution-frontmatter-contract.md`, evergreen component language.

- DONE: Gate-review retirement — amend every operative section to the approved overridable-channel architecture.
  Retired `spacedock gate review` from the first-use journey, the Go helper boundary, and behavioral-proof item 5; presentation is now an overridable channel of the present-gate skill (default chat; override = a provider-owned hardened script); result validation + provider id-normalization are recorder-side verbs (binary subspace-free, checkable); the channel calls the recorder and never writes gates. Schema owner tag: id-normalization "specified AND implemented recorder-side". Diagram 2's obtain node corrected (validate+normalize moved to the recorder node). No `gate review`/retired-arch phrases remain; both diagrams parse as `flowchart-v2`. Source: xb's attempt-3 approved entity.
- DONE: Digest domains — name the two domains and the divergence fixture.
  Added a "Digest domains" paragraph to the schema section: the canonical-bytes (JCS) briefing digest emitted/validated going forward; the raw-file pin is the marked legacy domain (shaping-era records, honest history, no rewrite); the recorder accepts it on replay. Behavioral-proof item 7 is the formatting-only fixture (JCS stable, raw-file changes).
- DONE: Consumption semantics — align the lifecycle to authorization-only.
  Rule 6 now reads: `consumed` marks the authorization spent atomically with the status transition, provably once; the dispatch effect is the dispatch machinery's, at-least-once retryable; receipts stay declined. Named the two crash windows; behavioral-proof item 8 is the authorization-side fixtures that surface them without double-firing.
- DONE: Recording identity — add the sentence to the lifecycle rules.
  A resolution is recorded under the identity that rendered it; a chat-directed closure records under the First Officer's identity on delegated authority with the directive quoted; adopting an advisory provider result as binding carries an explicit adoption note naming the authorizer.
- DONE: Confirm the amendments leave the seven retained resolution ACs unchanged.
  AC-1, AC-4, AC-6, AC-10, AC-12, AC-13, AC-14 keep their resolution-side cores; these amendments touch only presentation architecture, digest domains, consumption semantics, and recording identity — not the AC set.

### Summary

The contract's operative sections now match the architecture the captain approved at xb's
gate — presentation retired from the binary to an overridable present-gate channel,
validation/normalization homed recorder-side, the binary subspace-free — closing the codex
seat's Material 3 shared-region collision. Material 2 (two named digest domains + a
divergence fixture) and Material 6 (authorization-only consumption + two crash-window
fixtures) align to the captain's index rulings, and the recording-identity sentence lands
in the lifecycle. New contract SHA-256 `681b2348…`; the body reference updated. The seven
retained resolution ACs (AC-1, AC-4, AC-6, AC-10, AC-12, AC-13, AC-14) are unchanged.

## Stage Report: implementation

- DONE: Ship the recorder-only binary surface: validate and mutate gate/attempt/briefing/resolution records with frozen-closure and pointer invariants, while preserving the h1-owned application subtree and leaving status/dispatch untouched by recording.
  Commit `1095be38`; `TestRecordCloseNormalizesOnlyAfterDigestMatch`, `TestRecordRefusesPointerConflictAndFrozenMutation`, and `TestConcurrentWriterFailsClosed` fail if recording escapes `gates`, launders a mismatched result, mutates a closure, or admits a competing writer.
- DONE: Prove the behavior with the eight-history replay plus real pointer-conflict, frozen-closure, digest-domain, result-validation/id-normalization, unrelated-`--set`, and recorded-state status fixtures; keep the binary Subspace-free.
  `TestEightHistoryReplayPreservesApplicationsAndUnknownFields`, the recorder table tests, `TestUnrelatedSetPreservesGatesAndStatusProjectsResolution`, and `TestStatusTextAndJSONProjectAllRecordedResolutionStates` fail on lost application/history fields, bad CAS/digests/results, changed gates under `--set`, or missing approve/hold/revise projections; source grep finds no Subspace dependency in the binary.
- DONE: Complete the contract landing pass and recorder ownership prose, stay within the approved ~400–650 production-LOC surface unless the hard self-check triggers, and commit a green implementation (`go test ./...`, race where applicable, gofmt clean).
  Landed `docs/specs/gate-resolution-frontmatter-contract.md`, schema/reference prose, and the two recorder verbs; `internal/gates` is 636 production lines, with `go test ./...`, `go test ./... -race`, and `gofmt -l ./cmd ./internal` all clean at `1095be38`.

### Summary

The binary now owns atomic, pointer-checked `gates` writes for open, rebind, close, and
supersede operations, validates and normalizes provider results only after digest equality,
and surfaces recorded resolution state without changing the default status table. The
application subtree remains opaque and preserved, legacy shaping histories replay honestly,
and presentation, eligibility, transitions, dispatch, and Subspace stay outside the recorder.

## Stage Report: validation

- DONE: **AC-1** Recorded approval and reviewed digest survive a cold file read without status advance or dispatch.
  `TestRecordCloseNormalizesOnlyAfterDigestMatch` plus the detached A→B→C→close replay fail on lost identity/digest or any failed-close mutation.
- DONE: **AC-4** `approve`, `hold`, and `revise` portable rationale rules and durable visibility were reproduced.
  `TestPortableResolutionValidation` and `TestStatusTextAndJSONProjectAllRecordedResolutionStates` fail if reasonless revise/hold or a missing recorded projection is admitted.
- DONE: **AC-6 (record-state subset)** Text and JSON status distinguish recorded approve, hold, and revise.
  The three-decision status fixture asserts exact JSON field values and their corresponding text rows.
- DONE: **AC-10 (VALUE)** Recorder success changes only `gates`; digest/CAS failures leave the whole entity unchanged.
  The close fixture compares all bytes outside `gates`, and the mismatch/lock fixtures compare the complete failed-write file.
- FAILED: **AC-12** One entity and its evidence must cover two logical gates, multiple attempts, and fail-closed forks.
  `TestEightHistoryReplayPreservesApplicationsAndUnknownFields` fabricates one gate with eight attempts; it is neither the promised multi-gate fixture nor a replay of the production histories.
- FAILED: **AC-13** Wrapper fields must remain outside the copied portable Resolution.
  Detached test `TestAdversarialWrapperFieldsStayOutsideCopiedResolution` fails because `selectResolution` copies `stage`, `sequence`, and `application` through `Entry.Extra`.
- FAILED: **AC-14** The landed suite must prove open-attempt A→B→C rebinding, closure freeze, and re-entry.
  Current behavior passed a detached lifecycle test, but both a no-op `rebind` mutant and a disabled `supersede` mutant left `go test ./...` green.
- FAILED: Reproduce evidence for every retained recorder AC and the implementation checklist, including eight-history replay, exact gates-only mutation, frozen/CAS failures, digest domains, provider-result normalization, application preservation, and approve/hold/revise status projection.
  Gates-only writes, CAS/freeze, digest domains, normalization, application preservation, and projections reproduced; AC-12/13/14 evidence failed as detailed above.
- DONE: Perform the required semantic adversarial pass and detached high-stakes audit on a throwaway checkout: construct claim-breaking edits against recorder invariants/tests, exercise adjacent lifecycle variants and atomic failure behavior, and classify every finding by defect kind and release scope.
  Detached HEAD `1095be38` exercised the lifecycle and atomic refusals, ran two green claim-breaking full-suite mutants, and reproduced one failing reserved-field boundary test; the checkout was removed afterward.
- DONE: Verify the approved surface and landing pass (no application/eligibility/presentation ownership leak, no Subspace binary dependency, no shaping owner/task tokens in the landed contract), then run `go test ./...`, `go test ./... -race`, and confirm gofmt cleanliness.
  Product dependencies/help expose only record/validate; no Subspace package or shaping token landed; both full suites pass and `gofmt -l ./cmd ./internal` is empty.
- FAILED: Recommendation: REJECTED; no deferred risk is being used to dilute the material findings.
  AC-13 has a material outcome defect at the portable-record boundary; AC-12/14 have material evidence defects on supported, promised recorder lifecycles.

### Summary

Baseline and race suites pass, formatting is clean, and the present implementation completes the adjacent lifecycle in a detached audit. Validation rejects the gate because reserved Spacedock fields leak into the copied Resolution and the landed tests do not protect multi-gate replay, successful rebinding, or supersession.

### Feedback Cycles

- Cycle 1: REJECTED — detached high-stakes audit; surface 17 files/1696 changed lines (774 production, 367 test) vs estimate 2-3 new internal files plus status/CLI/docs, 400-650 production LOC with roughly equal test LOC (119% of upper production estimate); AC unchanged
- Cycle 2: REJECTED — captain-directed design reset to ideation after the validation gate; surface 27 files/2736 changed lines (778 production, 1403 test/fixture) vs estimate 2-3 new internal files plus status/CLI/docs, 400-650 production LOC with roughly equal test LOC (120% of upper production estimate); AC unchanged
- Cycle 3: REJECTED — captain ideation escalation; surface 125 proposal lines vs revised 100-180 documentation lines (69% of upper documentation estimate); AC unchanged
- Cycle 4: RECONFIRMED — captain accepted the implementation surface expansion; surface 14 files/2004 changed lines (1066 production, 822 test/fixture, 116 contract/help) vs revised 220-360 production, 300-500 test/fixture, 80-150 documentation lines (296% of upper production estimate; 148% of the 2x production ceiling); AC unchanged
- Cycle 5: REJECTED — validation/captain boundary reset; surface 36 files/4557 changed lines with the captain-reconfirmed production surface still 1066 changed lines vs revised 220-360 production lines (296% of upper estimate); AC narrowed: v1 owes no compatibility or migration for prototype `gates:` encodings, so the flow-map/unknown-field collision is out of supported scope and not a release blocker. Canonical multi-artifact association completeness remains material: derive the complete inventory from independent canonical Briefing bytes bound by the frozen JCS digest, then require exact presentation mapping before normalization. Simplify the binary-owned v1 writer boundary rather than adding arbitrary prototype-shape preservation.

## Stage Report: implementation (cycle 2)

- DONE: Closed AC-13 at the complete portable boundary: provider result envelopes are copied through an explicit portable Resolution field allow-list, so `stage`, `sequence`, Briefing-change, `application`, and arbitrary future wrapper fields cannot enter the durable Resolution. `TestAdversarialWrapperFieldsStayOutsideCopiedResolution` fails on any such leak.
- DONE: Replaced the fabricated replay evidence for AC-12 with exact frozen frontmatter snapshots of all eight 0260 production entities, plus a two-logical-gate/eight-attempt contract fixture. The tests prove multiple stable attempts per gate, current-pointer agreement, opaque application and unknown-field preservation, and fail-closed pointer, lineage, and frozen-history forks.
- DONE: Added the exact A→B→C→close→new-attempt lifecycle for AC-14. The test observes both successful rebinding changes, closure freeze with a no-mutation refusal, and supersession to a distinct sequence-2 attempt, killing the success-reporting no-op rebind and disabled-supersede mutants. `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` are green.

### Summary

Correction commit `9d279b87` closes validation cycle 1 without changing the approved scope. Durable replay now covers the eight real 0260 histories and a multi-gate re-entry fixture; the recorder admits only portable Resolution fields from provider envelopes; and the shipped suite directly protects rebinding, freeze, and supersession behavior.

## Stage Report: validation (cycle 2)

- DONE: **AC-1** Recorded approval and reviewed digest survive restart without advancing or dispatching.
  The close/re-read fixture preserves exact binding identity and digest; mismatched-digest failure leaves the complete entity unchanged.
- DONE: **AC-4** Portable approve/revise/hold decisions remain valid, durable, and visible under their specified rationale rules.
  Provider and status fixtures accept reasonless approve, reject unsupported reasonless revise/hold, and project all three recorded decisions.
- DONE: **AC-6 (record-state subset)** Status text and JSON distinguish approve, hold, and revise from frontmatter alone.
  `TestStatusTextAndJSONProjectAllRecordedResolutionStates` asserts the three exact JSON decision values and text visibility.
- DONE: **AC-10 (VALUE)** Recorder mutations remain gates-only and failures remain atomic.
  Outside-gates byte comparison, complete-file failure comparisons, pointer/CAS refusal, application preservation, and direct replay all pass.
- DONE: **AC-12** Multi-gate, multi-attempt, frozen-history, and production-replay evidence is now real and complete.
  Eight fixture frontmatters are byte-identical to their named sources at state commit `9594033d`; the two-gate fixture asserts eight attempts, pointers, lineage, applications, extensions, and fork refusals.
- DONE: **AC-13** Only the portable Review & Gate v1 Resolution field set crosses from a provider envelope.
  The wrapper test drops known and future-unknown envelope fields while preserving the binding Resolution; a future-only leak mutant turns it red.
- DONE: **AC-14** Open binding evolution, closure freeze, and re-entry are protected end to end.
  The A→B→C→close→D test asserts every identity/digest/state transition, no-mutation frozen refusal, and distinct sequence-2 lineage.
- DONE: Re-review correction commit `9d279b87` against the three cycle-1 material findings: prove wrapper fields cannot leak into portable Resolution state, the fixture truly covers two logical gates and eight production histories, and the exact rebind/freeze/supersede lifecycle is asserted.
  Independent source comparison and focused behavioral runs close AC-12, AC-13, and AC-14 without relying on the implementation report.
- DONE: Repeat the detached adversarial mutations that previously survived (success-reporting no-op rebind and disabled supersede) and confirm the full suite now kills both; probe the allow-list boundary with future unknown wrapper data.
  On detached `9d279b87`, each full-suite lifecycle mutant failed in `TestRebindCloseFreezeAndSupersedeLifecycle`; a future-only wrapper leak failed the adversarial boundary test.
- DONE: Reproduce all seven retained ACs, full test/race/gofmt evidence, approved-surface and contract-landing checks, then issue PASSED or REJECTED with material and deferred findings separated.
  `go test ./...` and `go test ./... -race` pass; gofmt is clean; help/dependencies expose recorder-only verbs with no Subspace, application, eligibility, presentation, or shaping-token leak.
- DONE: Recommendation: PASSED.
  All cycle-1 material findings are closed; the sole residual risk is deferred and version-triggered below.

### Summary

Correction commit `9d279b87` closes the portable-boundary outcome defect and both lifecycle evidence defects. The independent detached audit now kills all three adversarial edits, the required suites pass, and validation recommends PASSED.

### Reviewer Findings

- Material: none.
- Deferred risk: if a future portable Review & Gate version adds a required Resolution field, the v1 allow-list will discard it as unknown envelope data. This is outside the current v1 promise; current fields reproduce exactly, and the risk becomes material when Spacedock advertises or receives that newer required field.
- Polish: none.

## Stage Report: validation (cycle 3)

- DONE: Create one concise validation-gate briefing at `review/validation/briefing-1/gate-review.md` that leads with the recorder capability and recommended decision; the entity file is a reference artifact, never the presented body.
  The standalone briefing opens with APPROVE, explains the recorder in one paragraph, and links the detailed entity only as a reference.
- DONE: Explain how the capability works and how it was falsifiably tested: command surface and ownership boundary, all seven AC outcomes, the cycle-1 rejection/correction, three detached adversarial mutants, and full test/race/gofmt results.
  The evidence table covers AC-1/4/6/10/12/13/14, followed by the correction and the exact three red adversarial edits plus suite results.
- DONE: Make the decision legible in one screen: material findings, deferred risk, exact commits/paths as references, and what approval authorizes next; commit the briefing artifact to the state checkout without changing the validated product branch.
  The decision block records zero material findings, the v1 allow-list version risk, merge/landing authority, and commits `1095be38`, `9d279b87`, and `0c0fb6ca`; product HEAD remains unchanged.

### Summary

Packaged the passed validation evidence into a compact captain-facing gate review without rerunning validation or changing product code. The briefing recommends approval and makes the recorder capability, falsifiable evidence, residual risk, and merge authorization readable without presenting the entity body.

## Stage Report: validation (cycle 4)

- DONE: Complete `review/validation/briefing-1/` as a canonical Review & Gate v1 package: keep `gate-review.md` primary, add a frozen entity snapshot and frozen landed-contract snapshot as reference artifacts, and bind their exact SHA-256 revisions in `briefing.json`.
  The package contains one primary and two frozen references; manifest id `briefing:docs-dev:3k:validation:attempt-1:revision-1` binds all three relative URIs with reproduced `sha256:` revisions.
- DONE: Make the manifest's question and routing decision-specific: approve authorizes merge/landing, revise returns to implementation with concrete findings, and hold names a prerequisite; use stable component vocabulary and the validation-gate attempt identity.
  Routing maps approve to `ceremony:merge-landing`, revise to `stage:implementation`, and hold to `gate:validation` with a named-prerequisite label.
- DONE: Verify every manifest revision reproduces from its committed artifact bytes and commit only the package/state report; do not change or retest the validated product branch.
  JSON shape and all three artifact digests were checked before commit and are rechecked from Git object bytes after commit; product HEAD remains `9d279b87` and no product test was rerun.

### Summary

Completed the validation review room as a portable Briefing v1 package with a concise primary review, frozen entity and contract references, exact artifact revisions, and decision-specific routing. Only package and state-report artifacts changed; the validated product branch remains untouched.

## Ideation delta: binary-owned recorder journey (cycle 22)

**What changes from the validated design:** replace the agent-authored transaction envelope
with semantic recorder commands. The binary now derives `open` versus open-attempt rebind
versus post-closure successor, reads the current pointer inside the same entity lock, mints
gate/attempt/Resolution ids, and writes only the decision projection. This delta supersedes
the command, schema-mechanics, and provider-result portions above; the full entity and
`docs/specs/gate-resolution-frontmatter-contract.md` remain reference material, not the next
gate's presented body.

### Smallest agent-facing journey

Before, an agent had to run `spacedock gate record ENTITY --operation op.yml [--briefing
briefing.json]`; `op.yml` selected `open|rebind|close|supersede`, repeated the expected
gate/attempt/Briefing/digest pointer, minted ids, and translated provider output into
recorder-specific `entries`.

After, the surface is:

```text
spacedock gate bind ENTITY --briefing ROOM/briefing.json [--workflow-dir DIR]
spacedock gate resolve ENTITY --result FILE --actor ID [--adoption-note TEXT] [--workflow-dir DIR]
spacedock gate resolve ENTITY --decision approve|revise|hold --actor ID [--reason TEXT] [--directive TEXT] [--workflow-dir DIR]
spacedock gate validate ENTITY [--workflow-dir DIR]
```

`bind` validates a complete immutable Briefing and derives the logical gate from entity
identity plus current stage. No existing gate means open; a changed Briefing while the
current attempt is open means **rebind** (same attempt, binding replaced, old binding only
in Git/room); a bind after closure means **supersede** (append a new attempt, never mutate
the closure). `resolve --result` consumes exact provider bytes. The chat form receives only
the semantic decision, reason/directive, and recording identity; the binary supplies time,
ids, portable shape, CAS, and lifecycle operation. The unshipped `--operation` form is
removed, with an error pointing to `bind`/`resolve`. A nonbinding/advisory provider result
requires an authorized adopting actor plus an adoption note; a delegated chat decision
requires the quoted directive.

### Durable projection: retained and removed

The v2 contract retains `version`, one `current-gate`, logical gate `id` + `stage`, ordered
self-contained attempts, recorder-owned attempt `id`, the current/frozen Briefing
`id`/digest/domain/room reference, the exact portable Resolution, and the h1-owned
`application` subtree unchanged. These are the minimum facts needed to cold-read multiple
logical gates, locate and verify reviewed bytes, distinguish attempts, and reproduce each
decision without Git.

It removes the duplicated `current.attempt`, `current-attempt`, `sequence`,
`previous-attempt`, and `state` contract fields: list order gives sequence/lineage/current,
and Resolution absence/presence gives open/closed. Recorder-authored record/attempt/Briefing
notes also leave the v2 contract; decision provenance belongs in the Briefing, Resolution,
or adoption note. Operation, expected CAS, and candidate ids were never decision facts and
are no longer accepted input or stored history. Git remains the audit log for rebinds and
migrations. Legacy notes/extensions remain opaque compatibility data, not fields new writes
mint.

### Briefing and Result ownership; exact fixture spike

The gate-attempt/presentation side owns the readable review, complete immutable Briefing,
artifact resolution, and durable result transport. The recorder accepts that Briefing,
verifies it, and records decisions; it never builds presentation. Provider code owns
`result.json` bytes. The recorder verifies the current canonical Briefing's artifact
revision and actor authority before copying the portable Resolution, then normalizes only
its provider Briefing identity to the bound canonical identity.

The exercised fixture is `/tmp/subspace-3k-legible-gate.afzJuE/result.json`, exact SHA-256
`4609610352bef7206a7cab143a4768d30342bd101a7bd7692220cc72ba1464f7`. It is
`review-v1-result`/`advisory`/`binding:false`, carries its Resolution at `.resolution` (no
`.entries`), and names provider Briefing
`briefing:single-file:f8f13afcbd9bb2b3fb2732927934ac40`. The canonical package is
`review/validation/briefing-1/briefing.json`, id
`briefing:docs-dev:3k:validation:attempt-1:revision-1`. The spike reproduced that both name
the same primary bytes through revision
`sha256:d2a747755b1c348c499396a1900c0f94a0387565f2a0f4f1d9744b18c124a5a4` while the
Briefing identities differ. That is the first implementation fixture: exact input, artifact
match first, authority check second, identity normalization last; no fabricated provider
Result, `entries`, or `briefing-digest` envelope.

**Dependency note for xb:** the valid result arrived after the presenter helper had already
reported failure. The presentation transport must retain late provider completion bytes and
their canonical-package association even after controller failure. Repository files may be
Briefing References by URI + SHA only when xb/provider's resolver can reproduce those bytes;
otherwise xb freezes a room copy. The recorder neither polls the provider nor owns retention.

### Compatibility, ACs, proof, and expected surface

All eight 0260 histories remain first-class v1 fixtures. Reads and status never rewrite
them and keep current text/JSON behavior. The first successful gate write performs one
locked v1→v2 projection migration; it preserves ids, stages, Briefing values, exact
Resolutions, opaque applications, and unknown extensions, and removes only derivable
mechanics. Any legacy pointer disagreement fails closed without migration. Unrelated
`status --set` never migrates. Tests compare all eight histories before/after by semantic
projection and assert unchanged status; a two-gate lifecycle proves A→B→C binds reuse one
open attempt, resolution freezes C, and D appends a new attempt.

The seven end-value ACs remain **AC-1, AC-4, AC-6, AC-10, AC-12, AC-13, and AC-14**.
AC-13's mechanism wording changes from fabricated `entries` to exact Review v1 `.resolution`
or binary-constructed chat Resolution; AC-10/12/14 tests stop requiring exposed operation,
pointer, sequence, or state fields. The value claims remain unchanged: restart durability,
portable approve/revise/hold, status visibility, gates-only mutation, multi-gate attempts,
portable-boundary fidelity, and rebind/freeze/re-entry behavior.

Revised incremental surface: edit `internal/gates/{model,operation,io,gates_test}.go`,
`internal/cli/cli.go`, status projection tests if required, one exact Review v1 fixture, and
the contract/reference docs; no new production package. Expect ~250-400 production LOC
touched with net production LOC at or below `9d279b87`, ~300-500 test/fixture LOC, and
~100-180 documentation lines, tolerance 2x. Any provider launch/retention code, application
lifecycle, or inability to preserve all eight histories trips reconfirmation before build.

## Stage Report: ideation (cycle 22)

- DONE: Define the smallest agent-facing recorder journey: bind a complete Briefing or consume an exact Review v1 result/chat decision, while the binary derives lifecycle operation, pointer CAS, and recorder-owned IDs.
  The delta replaces `--operation` with `bind` and two `resolve` forms and specifies derived open/rebind/supersede/close behavior.
- DONE: Reduce durable gate metadata to the self-contained current/frozen decision projection; use Git for open-attempt rebind history and justify every retained or removed field against cold replay and compatibility.
  The v2 projection removes five redundant mechanics fields, retains decision facts, and gives the eight v1 histories a fail-closed atomic migration path.
- DONE: Return a concise captain-facing ideation delta covering before/after command surface, schema, ownership, exact-result spike, compatibility, acceptance-criterion impact, and revised surface estimate; do not implement before approval.
  This cycle records the verified exact-result digest/identity mismatch, xb transport dependency, unchanged end-value AC set, and bounded correction surface; product HEAD remains `9d279b87`.

### Summary

The reset makes the recorder semantic and binary-owned: agents bind a complete Briefing or
record a provider/chat decision, while lifecycle mechanics disappear from their input and
from the durable contract. The exact late Review v1 result now anchors the first fixture;
implementation remains paused pending approval, with xb owning presentation retention and
h1 retaining application lifecycle ownership.

## Ideation delta from cycle 22 (cycle 23)

The captain rejected the proposed v2/migration. This recorder is unreleased, so the design
now revises **v1 in place**. There is no migration or bulk rewrite: the reader accepts the
eight hand-authored/dogfood histories and preserves their existing pointers, `sequence`,
`previous-attempt`, `state`, notes, applications, and unknown fields. Targeted mutations
maintain a present legacy pointer/state when lifecycle consistency requires it, but new
minimal attempts do not mint those derivable fields.

The public surface returns to the approved two verbs. `spacedock gate record` accepts
exactly one semantic source (`--briefing`, `--result`, or chat `--decision` inputs) and
internally derives open/rebind/close/supersede, CAS, and recorder ids; `spacedock gate
validate` is read-only. The exact forms and gates-only before/after projections for all
three record paths are the primary artifact in
[`review/ideation/briefing-17/`](review/ideation/briefing-17/briefing.json).

Provider normalization now requires a retained association that covers the exact Result
digest, provider Briefing identity, bound canonical Briefing identity/revision, and complete
presentation mapping. Primary-artifact equality alone is an explicit rejection case. The
late Review v1 Result is frozen byte-identically at
[`exact-review-v1-result.json`](review/ideation/briefing-17/exact-review-v1-result.json),
SHA-256 `4609610352bef7206a7cab143a4768d30342bd101a7bd7692220cc72ba1464f7`;
without such an association it is negative evidence, not an adoptable decision.

The seven end-value ACs remain AC-1/4/6/10/12/13/14. The revised incremental surface is
~220-360 production LOC touched with net production LOC no higher than `9d279b87`,
~300-500 test/fixture LOC, and ~80-150 documentation lines, tolerance 2x. Exact Result and
association fixtures belong only under test discovery-safe paths during implementation;
provider transport remains xb-owned and application lifecycle remains h1-owned.

## Stage Report: ideation (cycle 23)

- DONE: Revise the unreleased v1 contract in place: remove v2/migration machinery, stop minting derivable mechanics on new writes, and tolerate existing historical fields without rewriting them.
  The cycle-23 delta keeps legacy fields in place, limits mutation to targeted compatibility maintenance, and defines minimal new attempts under v1.
- DONE: Keep the public surface to `gate record` plus `gate validate`, with mutually exclusive semantic Briefing, exact Result, and chat-decision inputs; require verified canonical-package association before provider identity normalization.
  Briefing 17 gives the exact command grammar and three gates-only projections; its exact late Result is a rejection fixture until full-package association exists.
- DONE: Freeze the exact Review v1 evidence durably and return only a concise delta from cycle 22, including retained compatibility behavior, unchanged end-value ACs, and revised surface.
  The room copy is byte-identical at SHA-256 `46096103…`; the delta confirms all seven retained ACs, the narrower surface, and paused product HEAD `9d279b87`.

### Summary

The corrected gate package presents an unreleased-v1 design with one semantic record verb,
binary-owned mechanics, exact gates-only writes, and no migration. It freezes the real late
provider output while refusing to equate a single matching artifact with a verified
multi-artifact Briefing; implementation remains paused for captain approval.

## Ideation delta from cycle 23 (cycle 24)

The captain approved briefing 17 (`lgtm`) and required one dogfooded fixture. Immutable
[`briefing-18`](review/ideation/briefing-18/briefing.json) preserves the approved design and
adds the exact cross-logical-gate re-entry case from state commit `71c61fbc`: workflow
status is `ideation`, ideation has nine closed attempts, and a closed validation `revise`
remains globally selected.

`spacedock gate record 3k --briefing revision-18.json --workflow-dir docs/dev` must derive
the ideation gate from workflow status rather than require it to equal `gates.current`,
append minimal ideation attempt 10, update the existing legacy/global selection, and leave
ideation attempt 9 plus validation attempt 1 byte-identical. The fixture pins the source,
both closure-range digests, the complete Briefing's canonical digest, the expected gates
diff, and the required mutant: reinstating the current implementation's global-current
prerequisite must turn the suite red.

All cycle-23 rulings and AC-1/4/6/10/12/13/14 remain unchanged. The new fixture fits the
approved ~300-500 test/fixture-LOC budget and does not change the ~220-360 production-LOC
estimate or 2x tolerance. Product HEAD remains `9d279b87`; implementation has not begun.

## Stage Report: ideation (cycle 24)

- DONE: Add the exact cross-logical-gate re-entry fixture: ideation closed, validation globally current and closed, workflow status returned to ideation, then `record --briefing` appends ideation attempt 10 and selects it without mutating either closure.
  `briefing-18/cross-logical-gate-reentry.md` pins state commit `71c61fbc`, exact closure hashes, canonical revision-18 digest, and the three-field/one-append expected diff.
- DONE: Revise the cleaned gate-review/spec and behavioral proof so the fixture kills the current recorder defect that requires the target logical gate to already be globally current.
  Revision 18 makes current-stage gate lookup independent of `gates.current`; its mandatory global-current mutant fails before attempt 10 can be appended.
- DONE: Publish an immutable revision-18 package and a concise cycle-23 to cycle-24 delta; preserve product HEAD and every previously verified artifact.
  The package composes immutable briefing 17, exact Result evidence, baseline contract, complete binding input, fixture, and fold review; product HEAD remains `9d279b87`.

### Summary

The approved v1 recorder design now includes the missing two-logical-gate re-entry proof.
Revision 18 demonstrates that workflow stage selects the target gate even when another
closed gate remains globally current, while both historical closures stay byte-identical;
no product implementation changed.

## Stage Report: implementation (cycle 3)

- DONE: Complete the two-verb agent surface: semantic `record --briefing`, exact `record --result --association --actor`, and chat `record --decision`; derive lifecycle/CAS/IDs internally and remove the unshipped `--operation` interface.
  Commits `3c874fcd`, `2b03dd76`, and `6b450c7d` implement the surface; `TestGateRecordChatDecisionAndRejectsOperationInterface` fails if chat provenance is lost or `--operation` becomes accepted, and `TestRebindCloseFreezeAndSupersedeLifecycle` fails if A/B/C stop reusing one open attempt, closure stops freezing C, or D stops appending a successor.
- DONE: Implement the cleaned v1 projection and compatibility boundary: minimal new writes, targeted preservation of all eight legacy histories, gates-only atomicity, full-package association before normalization, and discovery-safe exact Result/cross-gate fixtures.
  `TestEightProductionHistoriesSurviveTargetedSemanticWrite` byte-compares every prior record and all non-`gates` bytes across eight fixtures; the exact Result SHA-256 is `46096103…`, and `TestGateRecordConsumesExactResultOnlyWithCompleteAssociation` fails on missing/primary-only association, premature normalization, or provider-wrapper leakage.
- DONE: Update the landed contract/help and falsifiable tests, keep xb/h1 boundaries intact, report actual surface, and commit with `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` green.
  Contract/schema/help name only the semantic forms and preserve presenter/application ownership; final baseline and race suites passed. Relative to `9d279b87`, actual surface is production `+823/-243` across four files, tests/fixtures `+698/-124` across six files, and contract/help `+72/-44` across four files.

### Summary

The recorder now owns lifecycle selection, current-stage lookup, CAS, ids, exact Result
association validation, chat provenance, and surgical gates-only writes behind `record` and
read-only `validate`. The production surface exceeded the ideation estimate's 2x tolerance
because exact legacy byte preservation required a source-position writer and the complete
association verifier; no provider transport, application lifecycle, workflow transition,
dispatch behavior, new schema version, migration, or new production package was added.

## Stage Report: validation (cycle 5)

- DONE: **AC-1** A normal recorded approval and reviewed digest survive cold reads while
  `status` remains byte-identical and no application or dispatch record is created.
  Reproduced through `TestGateRecordAndValidateCLILeaveStatusUntouched`,
  `TestGateRecordDecisionClosesMinimalBriefingAttempt`, and the status projection tests.
- DONE: **AC-4** Portable `approve`, `revise`, and `hold` rationale rules, exact Result
  adoption, and restart visibility remain correct on the supported fixtures.
  `TestPortableResolutionValidation`,
  `TestProviderResolutionIncludesRequireSameBriefingAnnotation`, the exact Result CLI
  test, and the three-decision status test all passed.
- DONE: **AC-6 (record-state subset)** Text and JSON status independently surfaced
  recorded `approve`, `hold`, and `revise`; unrelated `status --set` preserved `gates`.
  `TestStatusTextAndJSONProjectAllRecordedResolutionStates` and
  `TestUnrelatedSetPreservesGatesAndStatusProjectsResolution` passed.
- FAILED: **AC-10 (VALUE)** A supported flow-style unknown field can make a semantic
  re-entry write return success while corrupting the selection pointer. With
  `gates.current: {shadow: 'gate:validation', gate: 'gate:validation', attempt: ...}`,
  the source-position editor changes `shadow` to `gate:ideation` instead of the parsed
  `gate` scalar, updates the attempt pointer, and leaves the entity unreadable with a
  current-pointer conflict. The failed operation is not atomic because it is reported as
  successful and commits invalid durable state.
- FAILED: **AC-12** The same detached source-position probe creates a gate/attempt fork
  rather than failing closed. The ordinary two-gate and eight-history fixtures pass, but
  they do not cover an earlier same-valued scalar on the target flow-map line.
- FAILED: **AC-13** A three-artifact association truncated to its first canonical artifact
  **and** first presentation mapping is accepted. The verifier compares the two submitted
  list lengths but has no independent package inventory, so a consistently truncated
  association self-declares completeness and permits identity normalization. The landed
  primary-only test removes only one side and therefore misses this case.
- FAILED: **AC-14** The approved ordinary cross-logical-gate re-entry fixture passes and
  preserves both closures, but the supported flow-map collision above makes the same
  re-entry journey select an invalid gate/attempt pair after appending its successor.
- DONE: Reproduce the public surface and ownership constraints.
  The only public verbs are semantic `record` and read-only `validate`; `--operation`
  exits 2; open/rebind/close/supersede, exact Result, chat provenance, wrapper exclusion,
  normal current-stage selection, and frozen closure all passed. Diff inspection confirms
  no v2/migration, provider transport, application lifecycle, status transition, dispatch,
  xb, or h1 ownership leak. The captain-reconfirmed production surface is exactly
  `+823/-243` across `internal/cli/cli.go` and the three production `internal/gates` files.
- DONE: Run the required gates.
  `gofmt -w ./cmd ./internal` produced no diff; uncached `go test ./... -count=1` and
  `go test ./... -race -count=1` passed all 18 packages. Focused CLI/gates/status tests
  also passed before the detached audit.

### Reviewer Findings

1. **Material outcome defect — AC-10/AC-12/AC-14, narrow writer correction.**
   `scalarEdit` locates a parsed YAML node only by its line and then replaces the first
   matching scalar bytes on that line. Flow mappings can contain an earlier unknown field
   with the same scalar, a supported v1 compatibility shape. The exact detached input
   above made `RecordBriefing` return nil and produced
   `current: {shadow: 'gate:ideation', gate: 'gate:validation', attempt:
   'attempt:ideation-2'}`; the next `Read` failed with a current-pointer conflict. Correct
   the exact node-span edit and validate the rebuilt bytes before rename/success.
2. **Material outcome defect and mechanism failure — AC-13; focused design reset.**
   The public association is its own only source of the canonical artifact inventory.
   Removing artifacts and their mappings together is therefore indistinguishable from a
   genuinely complete smaller package. This violates the approved requirement that the
   recorder verify the complete retained association before normalization. The end value
   remains reachable, but needs independent canonical-package ground truth (or an
   equivalently authenticated completeness binding); list-cardinality checks within the
   untrusted association cannot establish it. Per validation policy, route this boundary
   through a focused design reset rather than an automatic implementation-only repair.

No deferred risk or polish-only finding was identified. The normal-path and race suites
remain green, but they do not waive either supported data-integrity failure.

### Feedback Cycles

- **Detached validation-cycle-5 audit — source-position adversarial edit:** added a
  flow-style `gates.current` unknown field before `gate`, with the same old gate value.
  The semantic cross-gate re-entry returned success and wrote an invalid pointer fork.
- **Detached validation-cycle-5 audit — association-completeness adversarial edit:**
  removed canonical artifacts 2-3 and their corresponding presentation mappings from the
  exact three-artifact association. `verifyAssociation` returned nil.
- **Detached validation-cycle-5 claim-breaking controls:** replacing current-stage lookup
  with global-current lookup red-lined
  `TestGateRecordBriefingReentersStageGateWhenAnotherGateIsSelected`; normalizing the
  Result identity before provider verification red-lined the exact Result test; removing
  closed-attempt byte equality red-lined `TestFrozenMutationIsRejected`. These controls
  confirm the current-stage, normalization-order, and frozen-closure tests can fail for
  their stated claims.

### Summary

**REJECTED.** The semantic surface, ordinary lifecycle outcomes, status projection,
compatibility fixtures, scope boundaries, formatter, baseline suite, and race suite are
green. The detached audit nevertheless found two material supported-path defects: the
source-position writer can return success with corrupted gate selection, and the retained
association can self-declare a truncated multi-artifact package complete. Repair the
writer narrowly and reset the association-verification mechanism around independent
package completeness evidence before re-validation.
