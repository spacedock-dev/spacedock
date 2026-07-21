# Gate review: Gate recorder (3k) — ideation, attempt 2

**What you are looking at:** the revision answering your two float annotations. Both asks are addressed; the full revised body is below the divider.

**Your ask 1 — "why are there still 14 ACs?"** The cut is now physical, not just declared: seven in-scope ACs keep full text (1, 4, the record-subset of 6, 10, 12, 13, 14); everything moved is a one-line pointer to its owner — ACs 2/3/5/11 and the eligibility half of 6 to the blockers task, 7/15 to the presentation task, 9 deferred. Scheduler rules and the test plan got the same treatment. Nothing moved reads as buildable here.

**Your ask 2 — PR #510 alignment.** The draft Ledger gate-binding boundary and this recorder sit at different layers describing adjacent facts: PR 510 is a thin provider pin inside a Helm Ledger gate slot plus Helm-owned application facts; the recorder is the whole gate tree inside entity frontmatter, which stays the workflow-owned authority (settled). Ten-element read: the digest discipline is already IDENTICAL (RFC 8785 canonical bytes, sha256 — no change needed), application states map cleanly (pending↔pending_apply, consumed↔applied, superseded shared), field names align modulo casing. **Four genuine forks are flagged for you, none blocking this cut:**

1. **Identity authority** — recorder ids are Spacedock-minted; Helm owns its own. Elective: namespace the recorder's id strings now so a later binding maps cleanly (a handful of lines, inside tolerance) — or leave purely internal. Annotate a preference or it stays internal.
2. **Who folds pending→applied** — recorder: the existing transition machinery; PR 510: Ledger receipt acceptance. The reconciling read is coexistence (recorder authoritative when no Ledger is present; a Ledger-bound deployment adds the receipt fold outside). Default: coexistence, nothing built now.
3. **Receipts/idempotency** — the recorder deliberately omits them (settled); PR 510 centers them. Default: the omission stands.
4. **Projection/quarantine lanes** — no recorder analog because frontmatter IS the authority, not a projection; maps to the banked event-design artifact, explicitly out of this cut. Note-only.

**Surface reconfirmed:** the alignment adds no build scope — ~600-900 Go LOC + equal tests, 2× tolerance, hard self-check untripped. The sprint's digest-verifiability goal is owned by AC-10 + AC-12; no new AC minted.

**Recommend approve.** Checklist: 3 done, 0 skipped, 0 failed.

**Decision:** approve = attempt 2 closes with a pending advance to implementation, and the three downstream ideations (presentation, blockers, triage reframe) dispatch against this contract. Annotate fork 1 if you want the id namespacing; the other forks default to the stated readings. Revise = annotate; hold = discuss.

---

# Full task body (frozen snapshot at attempt 2)

---
title: Gate recorder — durable gates records with binary-owned writes
status: ideation
score: "0.80"
source: "Captain design feedback, 2026-07-13."
id: 3kd1x1gfxr8mdwzbmnwtjbw8
started: 2026-07-18T08:58:53Z
gates:
    version: 1
    current:
        gate: gate:docs-dev:3k:ideation
        attempt: gate-attempt:3k-ideation-1
    records:
        - id: gate:docs-dev:3k:ideation
          stage: ideation
          current-attempt: gate-attempt:3k-ideation-1
          attempts:
            - id: gate-attempt:3k-ideation-1
              sequence: 1
              state: closed
              briefing:
                id: briefing:docs-dev:3k:ideation:attempt-1:revision-8
                digest: sha256:3a8fd6d6702d212d72b708a406549a3a4c1d3f81997887e36d3453755721825b
                room-ref: "./review/ideation/briefing-8"
                note: "Frozen at closure. The digest binds the briefing-8 gate-summary artifact (summary + full post-cut snapshot), byte-verifiable in the room. Provider result validated by digest equality and retained as provider-result-8.json; the provider envelope id (briefing:single-file:e63586cd350f4f7b6cdcaa074a1ff312) is normalized to this attempt briefing id per the recorded id-mapping practice."
              resolution:
                type: Resolution
                id: resolution:actor-1784592481316587000
                briefing: briefing:docs-dev:3k:ideation:attempt-1:revision-8
                by: person:reviewer
                at: 2026-07-21T00:08:01Z
                decision: revise
                reason: "1. why are there still 14 ACs? i thought we trimmed this. 2. take a look at PR#510 to see where things align"
              application:
                action: feedback
                target-stage: ideation
                state: consumed
              note: "Subspace advisory float on the rebuilt tip binary, probe-first ritual observed. Two asks: physically trim the body to the cut (the AC section still carries every pre-cut criterion in full; the scope-cut prose named the retained set but never restructured the sections), and produce an alignment read against open draft PR #510 (Ledger gate-binding boundary). Routed to a fresh ideation revision worker; attempt 2 opens at re-presentation."
sprint: durable-decisions
group: recorder
---

# Gate recorder — durable gates records with binary-owned writes

## Scope cut (captain-approved, 2026-07-21)

This task grew four products in one coat (12 cycles, 15 ACs, two companion specs). It now narrows to ONE: the gate recorder and its record schema — the binary that owns every `gates:` frontmatter write (open / rebind-while-open / close-with-resolution / supersede / consume), the record invariants (one pending application per record, pointer agreement, frozen closures), snapshot-bound digests, and the status surfacing of recorded gate state. Retained ACs: 1, 4, 6 (record-state subset), 10, 12, 13, 14. This half is production-proven as a hand-run convention across eight entities in 0260 shaping (see `production-evidence-2026-07-20-fo-dry-run.md` and the 0260 closure findings: `--set` re-serialization, a self-conflicting attempt pointer, stale applications, the advisory-digest hole, entity-cannot-self-bind) — the recorder mechanizes a proven shape.

Moved out, one owner per concern:

- **Blockers, execution holds, and dispatch eligibility** (ACs 2, 3, 5, 11 + the eligibility subset of AC-6; scheduler rules 4-6 and 10's guard beyond convention) → `gate-blockers-and-eligibility`. The original seed concern — and the one part the production dry run never exercised (all eight recorded approvals had zero declared blockers). Sequenced after the recorder; its gate re-examines live need.
- **The presentation journey** (ACs 7 and 15; the one-command `gate review` blocking presenter, atomic result retention, briefing packages, probe-snapshot binding, provider id-mapping adapter; the probes companion spec rides along as convention) → `gate-review-presentation-command`. Subspace-coupled — sequenced with the subspace-tui surface, interim ritual per the 0260 shaping debrief.
- **Rejection-rework route context** (AC-9; the durable route edge) → DEFERRED, no task: the feedback-cycle prose convention just shipped in 0260; the binary edge waits for observed drift, per the escalation ordering.
- AC-8's behavioral-test mutants split with their owners.

Spec-level coupling retained here: this task keeps ownership of `gate-resolution-frontmatter-contract.md`; the presentation task's id-mapping rule (provider envelope briefing id normalized to the attempt briefing id after digest validation) is SPECIFIED in that contract and implemented in the presentation task. The open gate attempt's briefing must be rebound to this post-cut content before its next presentation (open-attempt rebinding, scheduler rule 2).

**Expected surface + tolerance (FO-drafted at re-presentation per the standing ruling — the entity predates it; captain corrects at the gate):** Go product code, first in-repo binary surface of the gate cluster: ~2-4 new files under `internal/` (gates block read/model/write + invariant validation + the record mutations) plus edits in `internal/status` (surfacing, `--set` coexistence so unrelated field writes leave `gates:` untouched) and 1-2 new `spacedock gate ...` verb entries in `cmd/`; ~600-900 production LOC, roughly equal test LOC (fixture replay of the eight 0260 production entities + the red fixtures: z7's real pointer conflict, a second pending application, a frozen-closure mutation). Contract doc unchanged (the spec is the already-banked `gate-resolution-frontmatter-contract.md`); ~10 lines of FO-contract prose naming the recorder as the gates-write owner. Tolerance 2×. Hard self-check: any schema change that breaks replay of the eight production entities, any subspace-tui coupling (that is xb's surface), or any blocker/eligibility computation (h1's) trips a reconfirm.

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
(SHA-256 `da8ed3d7cf6a580913179f60e27698845c1bf98fe226cf7b7db5e55e17b179cb`).
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
4. **→ h1.** Approval current with a blocker present (`approved-pending`, retain stage, no dispatch) is eligibility computation owned by `gate-blockers-and-eligibility`.
5. **→ h1.** Approval current with a captain execution hold active (`approved-held`, distinct from a Review & Gate `hold`) is eligibility computation owned by `gate-blockers-and-eligibility`.
6. **→ h1.** Applying the approval exactly once when the final blocker clears and any hold is released is eligibility/consumption computation owned by `gate-blockers-and-eligibility`.
7. If gate-defining input changes after an attempt closed, a new attempt is required;
   closed attempts never gain another Briefing. (Marking the closed application stale
   and keeping the task non-dispatchable is eligibility surfacing owned by **h1**.)
8. Review & Gate `revise` or `hold` closes the attempt. `revise` creates a pending
   feedback application; `hold` creates `action: none`, `state: not-applicable`.
9. Status surfaces the recorded stage, gate/attempt/Briefing/Resolution identities, open
   or closed state, and application state. Subspace presents Briefing/lens/assessment
   deltas through the stable room reference when requested. (The blocker-set,
   execution-hold, and staleness surfacing is owned by **h1**.)
10. The gate schema records only `pending`, `consumed`, `superseded`, or
   `not-applicable` application state; the existing transition/dispatch coordinator owns
   effect identity and crash reconciliation. (The one-use guard — a current approval
   authorizes exactly one successful application — is eligibility enforcement owned by
   **h1**.)
11. A rejection routed through `feedback-to` retains the rejected gate result, and
   re-entry at the gate after that closed result creates a new attempt. (Projecting the
   current lifecycle stage with explicit `feedback_rework` route context is **DEFERRED**,
   AC-9.)

## Acceptance criteria

Retained criteria keep full text. Every moved or deferred criterion is a one-line
pointer naming its new owner: **h1** (`gate-blockers-and-eligibility`), **xb**
(`gate-review-presentation-command`), or **DEFERRED**. The scope cut (captain-approved,
2026-07-21) is the authority for the split; nothing moved reads as in-scope here.

**AC-1** An approved blocked entity survives process restart and still reports the same durable approval, exact blocker, reviewed digest, and `approved-pending` condition without advancing or dispatching.

**AC-2 → h1.** Blocker-clearance eligibility (one advance + one dispatch on final-blocker clear, no redispatch on repeated passes) now lives in `gate-blockers-and-eligibility`.

**AC-3 → h1.** Digest-bound staleness (any reviewed-input change before clearance marks the approval stale with zero advance/dispatch effects) now lives in `gate-blockers-and-eligibility`.

**AC-4** Review & Gate `revise` and `hold` Resolutions remain durable, visible, and
non-dispatchable; blocker clearance cannot override them. A captain-facing rejection
for rework is stored as portable `revise` plus a Spacedock feedback application, not as
the superseded portable `reject` vocabulary. `approve` needs no portable rationale;
`revise`/`hold` require a nonblank reason or an included earlier same-Briefing
Annotation, exactly as Review & Gate v1 specifies.

**AC-5 → h1.** Fail-closed blocker evaluation (missing/ambiguous/unqueryable blocker state never reads as satisfied and never consumes approval) now lives in `gate-blockers-and-eligibility`.

**AC-6 (record-state subset)** Status text and JSON distinguish the recorded gate
states the recorder surfaces from entity frontmatter alone: pending approval, consumed
approval, a recorded Review & Gate `hold`, and a recorded `revise` with its feedback
application. → The active-execution-hold, unsatisfied/unknown/failed-blocker,
satisfied-but-not-yet-consumed, and stale-approval distinctions are computed eligibility
surfacing owned by **h1**; the fuller rejected-gate rework route context is **DEFERRED**
(AC-9).

**AC-7 → xb.** The one-command blocking gate-review presentation (explicit Briefing + frozen Probe Reference, canonical title, blocking child, atomic log/Resolution/diagnostics retention, ensign-unresolved-until-exit) now lives in `gate-review-presentation-command`.

**AC-8 (mutants split with owners).** The recorder retains its frontmatter
record/replay, open-attempt Briefing advancement, and concurrency/frozen-mutation
mutants (behavioral-test-plan items 2 and 9). The blocked/held/stale/duplicate-pass
mutants move to **h1**; the presentation mutants (early completion, detached worker,
live-Reference append, controller/child/validation/retention) move to **xb**.

**AC-9 → DEFERRED.** The durable rejected-gate → rework route edge (status projecting `feedback_rework` context with cycle and source-gate identity after restart) is deferred with no task; the feedback-cycle prose convention shipped in 0260 and the binary edge waits for observed drift.

**AC-10 (VALUE)** Recording either an approval or rejection changes only the entity's
versioned `gates` frontmatter collection: current `status` is byte-identical and no
dispatch receipt or worker exists. Deleting projection caches and reading the current
entity directly enumerates every logical gate, gate attempt, immutable Briefing
binding (replaceable when open, frozen when closed), exact adopted Resolution,
selection pointer, and latest application state. Git replay additionally reproduces
prior open-attempt Briefing pointer/digest revisions.

**AC-11 → h1.** Approve-but-do-not-dispatch (durable `approve` plus an active workflow-owned `execution-hold` that survives restart, distinct from a portable `hold` decision) now lives in `gate-blockers-and-eligibility`.

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

**AC-15 (VALUE) → xb.** The first-use value measurement (answering yes reaches the complete multi-source review through one blocking command; both branches leave `status`, dispatch roster, and worktree byte-identical until the FO records and applies a decision) now lives in `gate-review-presentation-command`.

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

Retained items prove the recorder's own ACs. Moved items are one-line pointers to their
owner (**h1**, **xb**, or **DEFERRED**); split items keep the recorder-proving part and
point the rest.

1. **Physical record contrast (AC-10).** Drive the real binary-owned gate recorder
   against an approving and a revising fixture. Exactly `gates` changes; `status`,
   process roster, dispatch state, and worktree remain byte-identical. Delete projector
   caches and prove the entity still reconstructs the exact Resolution. (The AC-7/AC-15
   presentation side moves to **xb**.)
2. **Cold read, replay, and schema (AC-6 record-subset, AC-10, AC-12, AC-14).** Validate
   the concrete two-gate/multi-attempt example through the shipped schema, restart, and
   invoke status. The direct read enumerates every gate, gate attempt, single `briefing`
   binding, exact Resolution, and minimal application; Git replay reconstructs prior open
   pointers and reproduces each recorded Briefing digest from its committed snapshot.
   Mutants that require `current-briefing`/`resolved-briefing`, embed provider history,
   or consult a cache fail.
3. **→ h1.** Approve-but-do-not-dispatch (record approve + active execution hold, restart, repeated passes yield zero effects; release the hold and observe exactly one application) is proven in `gate-blockers-and-eligibility` (AC-11).
4. **Restart visibility, blocker table split (AC-1).** The approved-blocked fixture
   survives restart and still reports the same durable approval, exact blocker, digest,
   and `approved-pending` condition. (The blocker-state/stale-digest eligibility table —
   `unsatisfied`/`satisfied`/`unknown`/`failed` and changed-digest, AC-2/AC-3/AC-5 with
   their mutants — moves to **h1**.)
5. **Revise/hold durability (AC-4).** Build state history containing a validation
   `revise` and its consumed feedback application; human/JSON status keeps the durable
   `revise` and feedback application after restart, and blocker clearance cannot override
   it. (The rejected-gate → rework route context and its false-rework contrast, AC-9, is
   **DEFERRED**.)
6. **Re-entry, pointers, and concurrency (AC-4, AC-10, AC-12, AC-14).** Extend the
   fixture through re-validation. Concurrent attempt/pointer writes, a close racing a
   pointer advance, mutation of a closed `briefing`/Resolution, and field-wise merge fail
   closed. (Stale-content supersession, AC-3, moves to **h1**; the rework-context replay,
   AC-9, is **DEFERRED**.)
7. **→ xb.** First-use presentation (`gate review` with a complete Briefing whose frozen Probe input/history is an explicit supporting `Reference`, title derivation, controller-failure retention, direct-Zellij equivalence, no branch) is proven in `gate-review-presentation-command` (AC-7, AC-15, presentation-side AC-8 mutants).
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
