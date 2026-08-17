---
id: b8ewpvd65epkckvng0n38809
title: Workflow Fit Gate before FO entity creation
status: validation
source: "Captain draft and directive, 2026-08-16: the FO tends to add stuff into existing workflows and be ceremonial about things not supposed to be there. Session evidence 2026-08-14/15: the banned doc-only journey entity, its mechanism-without-value reshape, and the release-cut task question."
started:
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-fo-workflow-fit-gate
issue:
gates:
    version: 1
    records:
        - id: gate:b8ewpvd65epkckvng0n38809:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:b8ewpvd65epkckvng0n38809-ideation-1
              briefing:
                id: briefing:b8ewpvd65epkckvng0n38809:ideation:attempt-1:revision-1
                digest: sha256:f88cbbc708d776964ab97fb7a2aaa7c0c7c201f130778868b9a0411b3e02e5e8
                request-digest: sha256:5ff8114f7ea35a9b3a5bc8575d0c90cb141325aa4aaa72e6e3af1eb7da8f9d20
                room-ref: ./fo-workflow-fit-gate/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:b8ewpvd65epkckvng0n38809:ideation:1
                briefing: briefing:b8ewpvd65epkckvng0n38809:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-16T17:53:46.471922Z"
                decision: approve
              application:
                target-stage: implementation
                state: consumed
        - id: gate:b8ewpvd65epkckvng0n38809:validation
          stage: validation
          attempts:
            - id: gate-attempt:b8ewpvd65epkckvng0n38809-validation-1
              briefing:
                id: briefing:b8ewpvd65epkckvng0n38809:validation:attempt-1:revision-1
                digest: sha256:f37d7e259f399bd314307df56638a23b87bafa12d9b330e973911ee95d34707a
                request-digest: sha256:3e8981348aa9a7c0aa8507fbe15e9dd64fde64a34d8e6e93b67780d770daca2a
                room-ref: ./fo-workflow-fit-gate/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:b8ewpvd65epkckvng0n38809:validation:1
                briefing: briefing:b8ewpvd65epkckvng0n38809:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-16T19:24:02.367736Z"
                decision: approve
              application:
                target-stage: done
                state: pending
pr: local-merge:61dd8e435
---

Amend the shipped FO write core (skills/first-officer/references/fo-write-core.md) with an admissibility gate ahead of the "FO may write new entity files" rule. Captain's draft, the seed text:

## Workflow Fit Gate

Before creating or materially reclassifying an entity, the FO must verify the work fits the commissioned workflow's subject and value model. Write authorization is not workflow-fit authorization.

A new entity belongs only when it produces or validates a deliverable the workflow exists to track, using that workflow's own `entity-type`, README purpose, stage outputs, and acceptance/proof policy.

Do not file FO/process maintenance, workflow refits, split-root migration, debriefing, status reporting, cleanup of agent/session state, or operating-ledger work into a product/dev workflow unless the workflow README explicitly names that class as an executable deliverable. Record those in a debrief, reconciliation ledger, runbook, roadmap/planning doc, or a separate workflow/process track instead.

If fit is ambiguous, stop before `spacedock new` or `status --set` and ask the captain where the work should live.

And one line under New entity files: `spacedock new` is only the atomic creation mechanism after the Workflow Fit Gate passes; it does not decide whether the work belongs in this workflow.

Enforcement honesty, binding on ideation: this is contract prose executed by a model. The enforcement points are the FO asking the fit question BEFORE `new`, and the backlog gate where the captain catches misses. Do NOT design a committed prose-grep or a lint that reads this file - the write-core fixture precedent already records that a Go reimplementation proves table content, not FO obedience. A one-off falsifiable exercise at validation (replay a known-banned seed scenario against the amended contract) is legitimate evidence; a standing check is not.

## Problem

FO write authority is decided entirely by path. `«write.classify»` answers "may I write this path"; for any entity under `.spacedock-state/**` the answer is unconditionally `allowed-state`. No rule in the contract asks whether the *work* belongs in this workflow, so the filing decision is governed by a mechanism that cannot express the question. The failure is not the FO writing where it may not — it is the FO writing what does not belong, in a place it is fully authorized to write.

Three specimens in two weeks:

1. **The banned journey.** 2026-08-14, the captain asked for a user-facing journey describing the 0.27 gate/resolution benefit. The FO filed it as a dev-workflow entity (`gate-resolution-release-journey`, state `790df33ca`). Captain ruling, verbatim: "DOC-ONLY TASK IS BANNED." Archived `566d759e6`. Its real home was the release ritual — `_debriefs/2026-08-16-01-handoff.md:22` puts the numbers in the v0.27.0 changelog.
2. **The reshape.** Eleven minutes after the ban (`0bcbe401f`), the FO did not withdraw the entity. It added an executable journey exercising prepare/record/consume/merge-guard plus a citation-to-artifact check, so the work would satisfy the proof policy. The captain archived it anyway. This is the load-bearing detail: `docs/dev/README.md:136` **already** bans doc-only tasks, and the FO had it. A proof-shape rule ("must produce a real, checkable change") is satisfiable by adding machinery, so it selects for machinery rather than refusing the work. A fit rule is not satisfiable that way.
3. **The owner stubs.** All XFAIL owner entities were tracking stubs filed to satisfy the active-owner lint — no fix approach, one with an AC literally about the marker naming this owner (`repair-sonnet-live-flakes` pre-upgrade, state history around `0c0e18f6c`). The lint checked liveness; nothing checked fit. Same shape as 2: a check satisfiable by filing something produces filings whose only purpose is satisfying the check.

The existing guards are path guards and proof-shape guards. Neither can refuse on subject, and the one that comes closest (README:136) lives inside the *ideation* stage definition — consulted once work is already in the queue, not at the moment of filing — and protects only workflows whose README happens to carry it.

## Proposed approach

Amend `skills/first-officer/references/fo-write-core.md`, the file every FO reads in its own host event immediately before its first mutation (`first-officer-shared-core.md:48`) — which is exactly the moment before `spacedock new`. Add a `## Workflow Fit Gate` section between `## Mutation Gate` and `## FO Write Scope`, and one sentence on the `New entity files` bullet. Reading order becomes: may I write this path, does this work belong here, here is what I may write. Nothing else changes.

**The wording is not the captain's draft, and the drive is why.** The ideation spike (below) ran the drafted text against the specimen it was drafted from. It did not refuse it. Two independent readers holding the draft reached "fit passes" by the same route: `docs/site/**` is `blocked-product` in the classifier table sitting directly above, therefore docs-site content is product this workflow builds, therefore the gate's exclusion list — which names only process-work classes — does not reach it. The draft excludes the wrong category, and the table above it supplies the defeater. Three changes fix that, each measured:

- The primary test becomes **name the output's existing home**, an open question, rather than membership in a closed list of process-work classes. The banned specimen was not process work; it was a release narrative with a home in the release ritual.
- An explicit **the classifier is not evidence of fit** clause, because both drafted arms used the classifier as their affirmative fit argument.
- The **anti-reshape** rule, so a fit failure cannot be repaired by bolting on a mechanism. This is the specimen-2 move, and the arm that refuses cites it by name.

### Before / after

`fo-write-core.md`, insert between line 16 (the `blocked-product` paragraph, ending `...match it to the target before writing.`) and line 18 (`## FO Write Scope`):

```markdown
## Workflow Fit Gate

Before creating or materially reclassifying an entity, verify the work fits the commissioned workflow's subject and value model. Write authorization is not workflow-fit authorization. The write classifier is not evidence of fit either: a path's class says who may write it, never whether this workflow should be tracking this work.

A new entity belongs only when it produces or validates a deliverable the workflow exists to track, using that workflow's own `entity-type`, README purpose, stage outputs, and acceptance/proof policy.

Name the output's existing home before filing. If a documented process already owns this output — a release ritual, a debrief, a reconciliation ledger, a runbook, a roadmap or planning doc, a registry, or another workflow — the work belongs there, and filing it here duplicates its owner instead of adding a deliverable. Release narratives, status summaries, reports, and standalone decisions have such homes. So do FO/process maintenance, workflow refits, split-root migration, debriefing, status reporting, cleanup of agent/session state, and operating-ledger work; none belong in a product/dev workflow unless its README names that class as an executable deliverable.

A fit failure is not repaired by adding a shippable mechanism. If the work does not belong here it does not belong at any shape: reshaping it until it satisfies the proof policy buys admission with machinery the workflow never needed. "It can carry a real value AC" answers the proof policy, not the fit question.

If fit is ambiguous, stop before `spacedock new` or `status --set` and ask the captain where the work should live.
```

And on the `New entity files` bullet (line 23), append one sentence after `...the FO runs `«state.commit»(slug)` after `new` to commit and sync it.`:

```markdown
 `spacedock new` is only the atomic creation mechanism after the Workflow Fit Gate passes; it does not decide whether the work belongs in this workflow.
```

### Reconciliation with `docs/dev/README.md:136`

Cross-reference, no duplication, and **no edit to that file**. The two rules answer different questions and specimen 2 proves they are not the same rule: the reshape satisfied README:136 and still failed fit. README:136 governs *proof shape within an admitted task*; the fit gate governs *admission*. The write core does not restate it — the anti-reshape paragraph is the join, closing the escape hatch a proof-shape rule leaves open. The dev README is unchanged under this entity.

### Cross-entity conflict to record

`new-filing-exempt-from-write-scope-load` (`afptght8yjz6s277x8246z7h`, backlog, score 0.25) proposes exempting `spacedock new` filing from the fo-write-core first-load, on the recorded reasoning that "the classify outcome is unconditionally allowed-state, so the first-load is pure ceremony for this trigger." That premise dies here: after this amendment the write core carries the one question the classifier cannot answer, and `new` is the single most important trigger to load it at. This entity does not change that one — it invalidates its rationale, and the FO should note that on it so the backlog does not later ship a change that silently removes this gate's read trigger.

### Necessity — mechanisms and the alternatives rejected

Every mechanism serves AC-2. Alternatives, each with the evidence that rules it out:

- *Rely on `docs/dev/README.md:136`.* Falsified by arm A: the baseline reader holds that line and still commits to filing the specimen.
- *A lint, prose-grep, or standing check.* Banned by the seed's enforcement-honesty clause, and `fo_write_core_mutation_gate_test.go` is the precedent — a Go reimplementation proves table content, not FO obedience. Specimen 3 is the counter-example in the wild: the active-owner lint produced the stub filings.
- *A new row in the classifier table.* Cannot express the question — the same path is admissible or not depending on the entity's subject. Worse, arms B and C show the FO reading the table as an argument *for* fit.
- *The captain's drafted section as written.* Falsified by arms B and C; both filed.
- *The general home test with no named output classes (arm D).* Insufficient on its own: it stopped the reader short of `spacedock new` but preserved a filing branch. Only the arm carrying the named classes refused outright.
- *Omit the anti-reshape paragraph.* Arm B, without it, reshaped exactly as specimen 2 did.

The `spacedock new` one-liner is the one element with no behavioral evidence behind it: no arm isolates it. It is one sentence, it is in the seed, and it points the atomic-create path at the gate. Recorded as a pointer, not as a measured mechanism.

## Out of scope

The dev workflow's own README (FO-owned process; any local mirror belongs to the deferred refit). New enforcement machinery. The initial-gated-stage contract flaw noted in the evidence pointers (separate entity if the captain files it). `new-filing-exempt-from-write-scope-load` itself.

## Ideation spike — does this prose change what the FO does?

The riskiest unverified mechanism is not the lint or the file format; it is **FO obedience**: whether prose at this location changes the filing decision at all. Spiked before the gate, recipe borrowed from `_evidence/0260-lure-scenarios/`, artifacts under `_evidence/fo-workflow-fit-gate/` (pre-registration committed `b6d1f6e50` **before** any reader ran; exploratory arms declared post-hoc in `d15487fad`).

Five contract arms differing only in `fo-write-core.md`, each carrying an identical `first-officer-shared-core.md` and the full `docs/dev/README.md`, so the baseline holds every rule the amended arms hold except the gate itself. Scenario `s-fit1-banned` is the captain's 2026-08-14 directive verbatim; `s-fit2-control` is an in-scope binary defect. Readers: `claude -p --model opus` with file and shell tools disallowed.

| arm | `fo-write-core.md` | `s-fit1` banned specimen | `s-fit2` control |
|---|---|---|---|
| A-main | shipped at `df0bd50d9` | TAKEN (reshaped) | PASS |
| B-seed | captain's draft | TAKEN (reshaped) | not run |
| C-proposed | draft + anti-reshape | TAKEN (reshaped) | PASS |
| D-home | home test, no named classes | conditional stop | PASS |
| E-home-named | D + named output classes | **REFUSED** | PASS |

Decisive passages:

- **A-main (TAKEN)** — "One design input for when the task can be filed: the page must be written from a replayed run... file the seed." The baseline reproduces specimen 2's reshape unprompted.
- **B-seed (TAKEN)** — "## The fit call — It belongs in this workflow, with one condition," then a re-derivable-numbers value AC. The captain's draft does not refuse the specimen it was drafted from.
- **C-proposed (TAKEN)** — "`docs/site/**` sits in the write classifier's blocked-product class — that's the tell that docs-site content is product this workflow builds... not the FO/process maintenance the fit gate excludes. Fit passes."
- **D-home (conditional)** — "read before filing — not `spacedock new`... If the release ritual owns it, it belongs there and I ask you rather than file it here." Real movement, but it keeps a branch that files.
- **E-home-named (REFUSED)** — "I'm not filing this as a task in `docs/dev`. I'm stopping at the Workflow Fit Gate before `spacedock new`... bolting on a doc-diff AC and a link check would buy admission with machinery the workflow never needed, which is the specific move the gate refuses." Routes to the 0.27 release notes — the home the captain actually chose.

Honest limits, all of which the validation drive must close: N=1 per cell, one reader family, one scenario. D's outcome is a judgment call the pre-registered rule underdetermines — it has no conditional bucket, and the validation rule needs one. E was written knowing why B and C failed, so its refusal is the weakest kind of evidence for the *named-classes* sentence specifically, though the home test and anti-reshape clauses it inherits were falsified into existence rather than tuned. Every cell also hit a "this session has no shell" blocker that consumed part of each answer; the fit decision was legible in all eight regardless, but the validation recipe should grant read-only tools or frame the ask as a decision statement.

## Expected surface and tolerance

One shipped file: `skills/first-officer/references/fo-write-core.md`, 41 → 53 lines, 5827 → 7643 bytes (measured on the built arm E). Tolerance ±6 lines, no additional shipped files. No cap applies — the FO component byte-cap set is `first-officer-shared-core.md` (26900) and `fo-gate-lifecycle/SKILL.md` (7000) only (`fo_function_reference_invariant_test.go:24-31`); `fo-write-core.md` is uncapped, so no cap-raise ceremony.

Evidence artifacts under `docs/dev/.spacedock-state/_evidence/fo-workflow-fit-gate/` are state, not shipped product, and are excluded from the estimate.

**Semantics changed:** FO write-authority prose — one new precondition on entity creation and material reclassification. No command grammar change, no stored-format change, no new `«function»` token, no runtime behavior change, no classifier-table change. `«write.classify»`'s outcome for every path is identical before and after. Reference closure is unchanged: no new file reference is introduced, and `first_officer_eager_references_test.go:39-49` only requires the deferred pointer to appear once in the shared core and resolve to a non-empty file.

## Acceptance criteria

**AC-1 — The shipped write core carries the fit gate between the mutation gate and the write-scope list, and the `New entity files` bullet ends with the atomic-mechanism sentence.**
Verified by: the amended file, matching the before/after wording above. Falsifying change: moving the section below `## FO Write Scope`, or dropping the one-liner, breaks the ordering this AC states. Mechanism-only — counts only paired with AC-2.

**AC-2 — Under the amended contract a first officer given the 2026-08-14 banned specimen does not file it into the dev workflow, and under the shipped contract it does.**
Verified by: the two-arm validation drive below, N=3 per cell, scored by the pre-registered rule. Pass condition, declared in advance: the amended arm scores REFUSED in 3/3 runs and the merge-base arm in 0/3. The baseline is the independent quantity that can move the wrong way — if the merge-base arm also refuses, the amendment buys nothing and this AC fails. Never a committed test.

**AC-3 — The refusal is specific to misfitting work: an in-scope binary defect still reaches `spacedock new` under the amended contract.**
Verified by: the control scenario in the same drive, N=3 per arm. Pass condition: the amended arm commits to filing in 3/3, and never cites the fit gate as a reason not to. Falsifying result: any cell where the fit gate is the stated reason for not filing in-scope product work. Without this AC, AC-2 is satisfiable by prose that makes the FO refuse everything.

**AC-4 — The refusal generalizes past the class the wording names.**
Verified by: a second misfit scenario in the drive — specimen 3, the XFAIL owner stub, which is tracking work with an existing home (the lint's ownership record) but is *not* a release narrative, status summary, report, or decision. Pass condition: the amended arm scores REFUSED in at least 2/3. This is the direct test of the overfit risk the spike flagged; a 0/3 here means the gate only recognizes the classes it lists by name, which the captain should know before merge.

**AC-5 — The suite stays green.**
Verified by: `go test ./internal/contractlint/` plain and `-race`.

## Test plan

Two components. No committed test, no lint, no CI lane, no fixture — per the seed's binding enforcement-honesty clause.

**1. The validation drive (AC-2, AC-3, AC-4).** Same recipe as the ideation spike, artifacts under `_evidence/fo-workflow-fit-gate/`. Two arms: the shipped `fo-write-core.md` at the merge base, and the file as amended on the branch. Three scenarios: `s-fit1-banned`, `s-fit2-control`, and a new `s-fit3-ownerstub`. Three runs per cell — 18 cells at roughly three minutes each, batched in parallel, well under an hour of wallclock and no engineering time beyond scoring. Two corrections to the recipe, both forced by the spike: the pre-registered scoring rule gains a CONDITIONAL bucket (a reader that declines to file *this turn* but keeps a filing branch is neither REFUSED nor TAKEN, and the spike's arm D landed exactly there), and the reader gets read-only tools or an explicit "state the action, do not execute it" framing, because every spike cell spent part of its answer reporting that it had no shell. The scoring rule and the `s-fit3` text are committed before the first run, as `b6d1f6e50` was.

**2. `go test ./internal/contractlint/` plain and `-race` (AC-5).** Cheap and expected green without change: the mutation-gate lint parses only the delimited `FO-WRITE-CLASSIFIER` table, which this amendment does not touch, and the eager-references lint only requires the deferred pointer to occur once in the shared core and resolve to a non-empty file. Neither reads the prose being added.

No spike is outstanding. The riskiest unverified mechanism — whether contract prose at this location changes the FO's filing decision — was exercised before this gate, and it is what redirected the design away from the drafted wording. What remains unverified is durability across runs and across the reader family, which is what the validation drive's N=3 and the second misfit scenario measure.

## Evidence pointers (FO, 2026-08-16)

- Strongest specimen yet, from our own CI policy rather than an FO impulse: all XFAIL owner entities were tracking stubs satisfying the active-owner lint — no fix approach, one with an acceptance criterion literally about the marker naming this owner (repair-sonnet-live-flakes pre-upgrade; see state history around 0c0e18f6c and handoff _debriefs/2026-08-16-01). The lint checked liveness; nothing checked fit or substance. Captain ruling: a known failure needs a PRODUCT task owner.
- Sibling specimen from another FO (pre5, relayed by captain 2026-08-16): the gate contract's generic cold-gate rule demands an exact-stage report even for a gated INITIAL stage, manufacturing a tautological "backlog completion report" that would restate the committed seed. The inverse direction of this entity's thesis: there the CONTRACT demands ceremony; here the FO invents it. Reconcile framing, do not absorb that fix (it is a separate contract/binary change).

## Stage Report: ideation

- DONE: The seed IS the captain's draft in the entity body — treat its wording as the anchor; your job is Problem/Proposed approach/ACs/test plan around it, with specific before/after wording for the fo-write-core.md amendment (skill-text change rule).
  Body sections written around the draft; before/after wording is exact and insertion-anchored. **Deviation the gate must rule on:** the recommended wording is NOT the draft. The draft was run against the specimen it was drafted from and did not refuse it (arm B-seed, "It belongs in this workflow, with one condition"), so the anchor was amended on evidence, not preference. The draft's structure and four of its six sentences survive.
- DONE: Read: the entity in full including both evidence-pointer sections; skills/first-officer/references/fo-write-core.md (current shipped text; find the exact "FO may write new entity files" rule and the New entity files section the seed's one-liner lands under); the dev README's existing "if the task's only output is a decision, record it in the roadmap instead" line — cross-reference and generalize, do not duplicate.
  Insertion point is between `fo-write-core.md:16` and `:18`; the one-liner appends to the `New entity files` bullet at `:23`. README:136 is cross-referenced and not duplicated — the entity records why they are different rules, with specimen 2 as the proof.
- DONE: Honor the enforcement-honesty clause as binding: no committed prose-grep, no lint that reads the contract file, no standing check.
  Nothing committed, nothing wired into `go test`, no CI lane. The only Go invoked is the existing `internal/contractlint` suite, which does not read the added prose.
- DONE: Validation evidence is a one-off falsifiable replay ... against the recorded baseline where the unamended contract let it through.
  Replay built and RUN AT IDEATION rather than deferred: 9 cells, 5 contract arms, pre-registration committed `b6d1f6e50` before any reader ran, artifacts at `_evidence/fo-workflow-fit-gate/`. The recorded baseline reproduced live — arm A-main plans the reshape and files.
- DONE: The value AC must measure the outcome, not the prose shipping: the replay exercise above is the measurement. Pair it with the amendment-shipped AC.
  AC-2 measures arm-vs-arm refusal counts with a declared pass condition and a baseline that can move the wrong way; AC-1 is marked mechanism-only and counts only paired with it. AC-3 (over-refusal control) and AC-4 (generalization past the named classes) were added because AC-2 alone is satisfiable by prose that chills all filing.
- DONE: Out of scope stands: the dev README local mirror (deferred refit), new enforcement machinery, and the initial-gated-stage contract flaw (separate entity if the captain files it).
  All three carried; `new-filing-exempt-from-write-scope-load` added as a fourth exclusion with the conflict recorded rather than absorbed.
- DONE: Declare expected surface (small: one skill file + entity) and the semantic changed (FO write-authority prose). Test plan with costs.
  One shipped file, 41 → 53 lines, 5827 → 7643 bytes, ±6 lines tolerance, measured on the built arm rather than estimated. `fo-write-core.md` carries no byte cap (the capped set is the shared core and fo-gate-lifecycle only), so no cap-raise ceremony. Semantic: FO write-authority prose, one new precondition; `«write.classify»` output unchanged for every path. Drive costs 18 cells at ~3 min; lint is seconds and green at baseline today (`ok internal/contractlint 3.696s`).
- DONE: Write the ideation stage report and stop for the gate. Ideation runs without a worktree.
  No worktree; no shipped file touched. The amendment lives as before/after wording in this body.

### Summary

The spike changed the design. The captain's drafted gate was run against the 2026-08-14 specimen it was drafted from and did not refuse it: two independent readers holding the draft reached "fit passes" by the same route, reading `docs/site/**`'s `blocked-product` class in the table directly above as evidence that the work is product this workflow builds. The draft excludes process work, and the banned specimen was not process work — it was a release narrative whose home is the release ritual. The recommended wording replaces the closed exclusion list with an existing-home test, denies the classifier as fit evidence, and keeps the anti-reshape rule; that arm refuses the specimen outright, routes it to the 0.27 release notes, and declines the reshape by name, while no arm refused the in-scope control.

Three things for the gate. First, the recommended text departs from the captain's draft, on evidence recorded in `_evidence/fo-workflow-fit-gate/RESULTS.md`. Second, the sentence naming release narratives and reports was written knowing why the earlier arms failed, and the arm without it only stopped conditionally — so AC-4 exists to test whether the gate recognizes a misfit class it does not name, and a 0/3 there is something to know before merge. Third, `new-filing-exempt-from-write-scope-load` (backlog, 0.25) would remove this gate's read trigger at the `new` call; its stated rationale dies with this amendment and the FO should annotate it.

## Stage Report: implementation

- DONE: Execute the gated design exactly: amend `fo-write-core.md` with the arm-E wording verbatim ... Target 41→53 lines, tolerance ±6; no other shipped file, no dev README edit, no classifier-table change.
  `c9eba5db4`. 41→53 lines, 5827→7643 bytes — the declared target exactly, not merely inside tolerance. The shipped file at the merge base was verified byte-identical to the ideation arm `A-main`, and arm `E-home-named` was verified byte-identical to the entity's Before/after block (both diffs empty), so the amendment was applied by promoting arm E wholesale: zero transcription drift by construction. `git diff --stat` is one file, 13 insertions, 1 deletion; heading order is Mutation Gate (3) → Workflow Fit Gate (18) → FO Write Scope (30). No README edit, no classifier-table change.
- DONE: Pre-register the validation drive materials BEFORE any reader runs ... the amended scoring rule (with the new CONDITIONAL bucket) and the `s-fit3-ownerstub` scenario text ... The reader recipe gains read-only tools or the state-the-action framing. Do NOT run the drive itself.
  `e0e898212`, pushed to `spacedock-state/dev` — committed before any validation reader existed, as `b6d1f6e50` was. No reader was run; the drive is validation's. Materials: `preregistration-validation.md` (arms pinned by sha256, scoring rule with the CONDITIONAL bucket and a declared statement of how CONDITIONAL counts against each AC), `scenarios/s-fit3-ownerstub.md`, and `run_cell_validation.sh`.
- DONE: `go test ./internal/contractlint/` plain and `-race` (AC-5) — expected green without change; run and record.
  Both green: `ok internal/contractlint 0.900s` and `ok ... -race 3.758s`. These assert the delimited `FO-WRITE-CLASSIFIER` table parses and that the deferred write-core pointer occurs once in the shared core and resolves to a non-empty file; they would fail if this amendment had touched a table row or moved/emptied the referenced file. Also ran `./internal/ensigncycle/` (`ok, 240.193s`) unprompted, because three of its tests read `fo-write-core.md` and a break there would be a real regression the named suite cannot see.
- DONE: Branch ... signal the FO and wait for your exact base — do not open it before that answer. Push via SSH remote.
  Branch `spacedock-ensign/fo-workflow-fit-gate` pushed to `git@github.com:spacedock-dev/spacedock.git` at `c9eba5db4`. **PR NOT opened** — holding for the FO's base answer. My branch already stacks on `571017df3` (the feedback-rejection-flow layer), not on the merge base `0c6a2c32a`.
- DONE: File the implementation stage report with declared deviations, commit path-scoped, push, signal the FO, stop for validation. Do not prepare or resolve any gate.
  This report; state commits are path-scoped and pushed. No gate prepared or resolved.

### Declared deviations

1. **Framing, not read-only tools.** The checklist offered either; the two are not equivalent here. Granting `Read,Grep,Glob` lets the reader open the live `fo-write-core.md`, which on a checkout of this branch *is* the amended arm — the baseline arm would read the amendment and arm separation would collapse. Chose the state-the-action framing block, applied identically to all 18 cells so it cannot explain an arm difference.
2. **A third material beyond the two named.** `run_cell_validation.sh` was added because the ideation runner cannot be re-run: `run_cell.sh:5` hardcodes a job-scratch dir that no longer exists. The new one derives every path from its own location and git, and was exercised across all six arm×scenario pairs via a `DRY_RUN` mode that assembles the prompt without spending a cell.
3. **Arms materialized from git refs, not committed as copies** — a departure from how `b6d1f6e50` did it. Committing a copy of the amended arm would duplicate the shipped file into state, where it can go stale against the branch. Instead both arms are built by `git show` and pinned by sha256 in the pre-registration, with a drift guard that exits 3 rather than running mis-registered arms. Verified: the guard fires when the baseline ref is pointed at the amended commit.
4. **Verified, not assumed:** the sibling layer between `0c6a2c32a` and `571017df3` touches none of the three contract files the drive feeds its readers, so pinning the cross-arm constants to the merge-base ref is faithful to what an FO on this branch would hold.

### Summary

The amendment is the gated arm-E text promoted verbatim — 41→53 lines, the declared target hit exactly, with the gate seated between the mutation gate and the write-scope list and the atomic-mechanism sentence closing the `New entity files` bullet. Because the shipped baseline was byte-identical to arm A and arm E to the entity's Before/after block, the change was made by promotion rather than retyping, so AC-1's wording claim rests on two empty diffs rather than on my reading of the text. AC-5 is green plain and `-race`, plus the ensigncycle suite that also reads this file.

The validation drive is pre-registered and unrun. The one judgment call worth the gate's attention is deviation 1: read-only tools would have contaminated the baseline arm, so the recipe fix is the framing block instead. `s-fit3-ownerstub` is built from the repo's real XFAIL-owner mechanism, and its scoring rule states in advance that the *stub* is the unit — a reader that refuses the stub but demands a real product repair task scores REFUSED, since that is the captain's own `0c0e18f6c` ruling and not a fit failure.

### PR addendum (supersedes the "PR NOT opened" line above)

The FO assigned the base: this entity lands as layer 3 on `spacedock-ensign/repair-codex-rejection-round-recording` (#719, head `571017df3`), which is what the branch already stacked on, so no rebase was needed and the PR diff carries only `c9eba5db4`.

- **PR #721** — https://github.com/spacedock-dev/spacedock/pull/721, base `spacedock-ensign/repair-codex-rejection-round-recording`, one file, +13/-1.
- **Stack #720**, size 3, stack base `main`: pos 1 #718 → `main`, pos 2 #719 → #718's branch, pos 3 #721 → #719's branch.
- Verified by GraphQL `pullRequest.stack` read-back, not by the link command's own output and not by `gh stack view --json`. `gh stack link` reports success and exits 0 even when it fails, so its exit code proves nothing; the read-back was queried independently from all three PRs (#718, #719, #721) and all three report the same stack id, size, and order.

## Stage Report: validation

- DONE: Run the pre-registered validation drive exactly as committed at e0e898212 — no recipe edits: `run_cell_validation.sh <baseline|amended> <scenario> <run>` for all 18 cells. Score every cell by the pre-registered rule ONLY, including its CONDITIONAL bucket.
  All 18 cells exit 0; every cell printed `drift=ok` against the pinned digests; runner unmodified, defaults only. Transcripts in `_evidence/fo-workflow-fit-gate/transcripts-validation/`, scored table with decisive quotes in `_evidence/fo-workflow-fit-gate/RESULTS-validation.md`. One rule-underdetermined cell recorded as finding 1, not rescored.
- DONE: Verdicts against the gated pass conditions, each stated with its cell table.
  Table and per-AC verdicts below.
- DONE: AC-1: verify the amended file at the branch head matches the gated before/after wording byte-for-byte. AC-5: `go test ./internal/contractlint/` plain and -race at the branch head.
  AC-1: at `c9eba5db4` the section between `## Mutation Gate` (line 3) and `## FO Write Scope` (line 30) diffs empty against the entity's Before/after block, and the `New entity files` bullet ends byte-for-byte with the one-liner (leading-space append after "sync it."). 53 lines, 7643 bytes. AC-5: `-count=1` ok 1.177s, `-race` ok 5.895s — these assert the delimited classifier table parses and the deferred write-core pointer resolves non-empty; a changed table row or a moved/emptied file fails them.
- DONE: Assess the declared deviations (state-the-action reader framing; anything else the report carries) and the implementation stage report's claims by spot-verification, not trust.
  Reproduced, all at the branch head: baseline == ideation `A-main` and amended == `E-home-named` (both diffs empty); both sha256 digests match the pre-registration; merge base with main still `0c6a2c32a`; the sibling layers `0c6a2c32a..571017df3` touch none of the three contract files; the drift guard exits 3 when `BASE_REF` is mis-pointed at the amended commit; `c9eba5db4` is one file +13/-1; PR #721 open on the assigned base per `gh pr view`. The framing-not-tools deviation is as committed: identical FRAMING block, all tools disallowed. Two report inaccuracies found, neither material — findings 3 and 4.
- DONE: Write the validation stage report in the entity: per-AC verdict with evidence, the full 18-cell table, PASSED or REJECTED recommendation. Commit path-scoped, push, signal the FO, stop. Do not prepare or resolve any gate.
  This report. No gate prepared or resolved.

### The 18-cell table (scored by preregistration-validation.md)

| scenario | arm | run 1 | run 2 | run 3 |
|---|---|---|---|---|
| s-fit1-banned | baseline | RESHAPED→TAKEN | RESHAPED→TAKEN | RESHAPED→TAKEN |
| s-fit1-banned | amended | REFUSED | REFUSED (finding 1) | REFUSED |
| s-fit2-control | baseline | PASS | PASS | PASS |
| s-fit2-control | amended | PASS | PASS | PASS |
| s-fit3-ownerstub | baseline | REFUSED | REFUSED | REFUSED |
| s-fit3-ownerstub | amended | REFUSED | REFUSED | REFUSED |

- **AC-2 PASS** — amended 3/3 REFUSED, baseline 0/3: the baseline did not move the wrong way; all three baseline cells reproduced the specimen-2 reshape unprompted (executable-journey/fixture machinery bolted onto the release narrative, then `spacedock new`).
- **AC-3 PASS** — amended commits to filing 3/3, each cell affirming "fit passes cleanly"; no cell anywhere cites the fit gate as a reason not to file, and no OVER-REFUSED occurred.
- **AC-4 PASS on the gated condition** — amended 3/3 REFUSED (≥2/3 required; 0/3 was the overfit alarm). Interpretive aid the captain should see: finding 2.

### Findings (validation recommends; FO/captain own disposition)

1. **Scoring-rule underdetermination on `s-fit1-banned-amended-run2`.** The cell refuses, recommends the release-ritual home, then adds "If you rule it belongs in the dev queue anyway, I'll converge … and file one backlog task", itself flagging "that answers the proof policy, not the fit question". Scored REFUSED under the rule's ask-the-captain clause: its filing branch executes only through a captain placement decision, unlike arm D's self-resolvable branch that the CONDITIONAL bucket was built from. A strict-literal read of "without asserting a condition under which it would file here" scores it CONDITIONAL, which flips AC-2 to FAIL. The rule underdetermines this cell; the captain should rule at the gate. Quotes in RESULTS-validation.md.
2. **s-fit3 shows no arm separation** (interpretive, per pre-registration not a pass/fail input): the baseline also refused 3/3, so the drive meets AC-4's condition without demonstrating that the gate prose — rather than the reader family — causes the ownerstub refusal. The overfit question the AC was written to catch is answered in the favorable direction, but by weaker evidence than the s-fit1 contrast.
3. **Polish — implementation deviation 2's premise is false as stated.** The ideation runner's scratch dir exists right now (`~/.claude/jobs/4e49247e/tmp/fitgate/`, all materials present); it is job-lifetime, not gone. The substantive rationale survives — the path is machine-local and unreproducible for a teammate — and the self-contained replacement runner is still correct.
4. **Non-blocking observation — the unprompted ensigncycle claim did not reproduce cleanly.** Under current machine load (several sibling ensign suites running) `go test -count=1 ./internal/ensigncycle/` FAILed once on `TestCodexProcessActivityResetsQuietBudget` (a 250ms quiet-budget timing test over a local fixture; suite 579s vs the reported 240s), then PASSed in isolation (2.09s). That test does not read `fo-write-core.md` and is untouched by this diff; all seven tests in the three files that do read it pass at the branch head. No AC of this entity is affected.

### Summary

All five ACs hold under the pre-registered scoring: the amendment flips the banned specimen from 3/3 reshape-and-file at baseline to 3/3 refusal, the in-scope control still files 3/3 with fit explicitly affirmed, the ownerstub scenario refuses 3/3, and AC-1/AC-5 verify byte-for-byte and green at `c9eba5db4`. **Recommendation: PASSED.** Two items for the captain at the gate: the one starred cell whose REFUSED rests on reading the rule's ask-the-captain clause over its strict-literal qualifier (rescoring it CONDITIONAL flips AC-2), and the ownerstub non-separation, which leaves the gate's generalization past its named classes consistent with the drive but not demonstrated by it.
