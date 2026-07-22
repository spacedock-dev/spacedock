---
id: 02avdajaz0q3hnjwycm5fq45
title: Ensigns triage review findings against declared stakes before fixing — decline disposition for correct-but-disproportionate findings
status: done
source: "0260 shaping — agent-derail forensics audit, 2026-07-19."
score: "0.7"
sprint: durable-decisions
group: recorder
started: 2026-07-20T05:04:07Z
worktree: .worktrees/spacedock-ensign-ensign-finding-triage-disposition
sprint-readiness:
gates:
    version: 1
    current:
        gate: gate:ensign-finding-triage-disposition:validation
    records:
        - id: gate:ensign-finding-triage-disposition:validation
          stage: validation
          attempts:
            - id: gate-attempt:ensign-finding-triage-disposition-validation-1
              briefing:
                id: briefing:docs-dev:02av:validation:canonical-v1:revision-1
                digest: sha256:fc31fe9a3ed18b5c8dae203e74431fd6494beac263dbc5d2b6e5af5a35f8905b
                digest-domain: canonical-bytes
                room-ref: ./review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ensign-finding-triage-disposition:validation:1
                briefing: briefing:docs-dev:02av:validation:canonical-v1:revision-1
                by: agent:first-officer
                at: "2026-07-22T16:28:18.503661Z"
                decision: approve
                reason: Exact candidate e85eb0cf passed 12/12 validation checks, AC-1/2/3, live replay, negative controls, detached audit, full and race tests, formatting, and clean-head verification; no material finding remains.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge.
mod-block:
pr: pr-merge:559
verdict: passed
completed: 2026-07-22T17:23:07Z
archived: 2026-07-22T17:23:07Z
---

`ensign-shared-core` contains zero guidance on consuming review findings — the exact actor that dutifully fixes a symlink edge case in a prototype has no rule to consult, and no disposition short of fixing exists for a substantively-correct-but-disproportionate finding. This adds the generic consumption rule (classify each finding against the workflow's committed finding-triage taxonomy AND the entity's own value ACs before fixing) and the decline disposition (correct-but-disproportionate gets a recorded decline, not a dutiful fix). **Reframe (captain-directed, 2026-07-21):** the disposition is recorded as an **advisory resolution under the gate-recorder model (3k)** — the ensign's triage is its own advisory resolution whose `includes` name each declined finding with class and why-not-material — not the `findings` prose field on a `### Feedback Cycles` entry the parked cut proposed. That answers the byte objection that parked the last cut: a per-entity structured record carries no always-loaded prompt-surface weight, where the ~2 KB prose block did. **Boundary (answered below, not assumed):** this entity owns **triage semantics only** — the rule text and the decline-as-advisory-resolution *shape* — and does NOT absorb the generic rounds-record plumbing (round briefings, room layout, frontmatter pointer, projection), which is the recorder's own generalization and belongs with a 3k successor or its own task. The rule *text* stands as approved-shaped; only the record/delivery mechanism reframes. Triage keys on per-entity value ACs + the committed taxonomy (the `validation` stage-def release-scope classification plus `.roborev.toml`'s four fields), never a workflow stakes field — that member is parked.

## Problem

- **No ensign-side consumption rule exists.** `skills/ensign/references/ensign-shared-core.md` carries zero guidance on triaging review findings before fixing them (forensics digest GAPS, `_evidence/0260-agent-derail-forensics/remedy-analyses-digest.txt:157`). The actor at the center of the incident class — the ensign that dutifully fixes a symlink edge case on a prototype (spacedock_subspace codex:019f63c6, synthesis incident 13) — has no rule licensing any disposition short of fixing. The audit's stated biggest gap.
- **The pressure is one-sided toward fixing.** The user-global `~/.claude/CLAUDE.md` pushes every worker hard toward fixing every reported finding ("ALL TEST FAILURES ARE YOUR RESPONSIBILITY", "Fix broken things immediately"); roborev's verdict mechanics fail a job on any `## Review Findings` entry (digest:151); and `## Review Findings` "never suppresses a real observation" keeps correct-but-disproportionate findings continuously in front of an actor "that has no consumption rule for them" (digest:153-154). An ensign holding a finding "has no licensed way to not fix it."
- **The deferred-risk record has no consumer.** roborev's config lets a reviewer defer with a trigger / why-outside-promise / promotion-condition note, but "nothing consumes the record: no tracker, entity field, or FO loop step ever revisits" it, so deferred risks "live only in gate messages" (digest:148). The decline needs a durable home the gate reads.
- **Fixing is not the only over-response; narrowing the claim is its twin.** The 0.25.1 incident (synthesis addendum) is the same repair-forward failure expressed as document edits: under a correct rejection, the task narrowed its value AC until a weaker claim passed. Declining a disproportionate finding and narrowing the claim it targets are opposite moves under the same pressure — the rule must make the first legal and the second captain-visible.
- **Not a suppression scheme.** The taxonomy already exists to classify, not silence (digest:165). This entity is the missing *consumption* half: production (roborev classifies + emits) already ships; the ensign that *receives* the findings still has no re-triage-and-decline rule.
- **The 0260 drive produced the drift this reframes onto a record.** Real hand-authored `### Feedback Cycles` rounds show the pathology concretely: 85's validation bounce (`_archive/merge-guard-arm-not-a-stopping-point`) triaged M1/M2/M3 as material and D1/D2 as deferred risks, but the deferred-risk dispositions lived only as a prose fragment inside the cycle line — no first-class, queryable, provider-preservable object. 2ae's Cycle 2 (`template-rigor-propagation`) narrowed a value AC's control assertion and *self-reconfirmed the narrowing inside the loop* ("the design-reset decision is RECONFIRM"), the exact move the AC-narrowing clause names illegal-for-the-ensign. These bank as the reframe's fixtures and its baselines-that-can-move-wrong.

## Boundary decision (reframe — the spine)

The reframe brief demands this be answered, not assumed. **This entity owns triage SEMANTICS only; it does NOT absorb the generic rounds-record plumbing.**

Owned here (the irreducible triage contribution no other member owns):

- **The consumption rule text** — review findings are inputs to triage, not a fix list; the three classes (material → fix, correct-but-disproportionate → recorded decline, needs-decision → escalate); a finding neither fixed nor recorded is not triaged; and narrowing a value AC to pass a finding is a captain-owned design-reset, not an ensign disposition.
- **The decline-as-advisory-resolution *shape*** — how a triage disposition is expressed as a resolution object in the recorder's settled vocabulary (3k): the ensign's triage is its OWN advisory resolution, its `includes` naming each declined finding with class, why-not-material, and the promote-to-material condition; advisory, never binding; a round can never advance status.
- **The all-declines semantics** — an all-declines round is a real recorded advisory resolution with zero fixes and each decline named; "no finding arrived" and "every finding declined" are different recorded states and must never render alike.

NOT owned here (deferred, one owner per concern — the 3k four-products precedent is the cautionary bound):

- **The generic rounds-record plumbing** — the recorder that writes a correction round's briefing/annotations/resolution into the entity's review room (append-only, the probes.jsonl pattern), maintains the frontmatter pointer, and projects the `### Feedback Cycles` line. This is *identical recorder machinery* to what 3k already ships for gate rounds, generalized from gate attempts to in-stage correction rounds. Absorbing it re-grows this entity into a second recorder — the "four products in one coat" 3k had to be scope-cut for. It is the recorder's successor (extend the recorder to in-stage rounds) or its own task; this entity only SPECIFIES the shape the plumbing must carry.
- **The application layer (h1)** — blockers, eligibility, execution holds, the one-use guard. A round record is advisory (`action: none` territory); bind only resolutions/briefings here.
- **The reviewer/production side** — how roborev classifies and emits. Unchanged; this is the consumption half.

The honest test of the boundary: if expressing the triage disposition required a NEW field in 3k's schema, that field would be plumbing and this entity would be over-absorbing. It does not — the spike (Claim 2) shows the shape rides 3k's settled contract (advisory Resolution + same-Briefing Annotation `includes`) with zero schema change. That is what keeps this a semantics entity, not a fourth product.

## Proposed approach

Three parts. The rule text (unchanged in substance) is delivered at the trigger; the disposition is recorded as an advisory resolution riding 3k's contract; AC-narrowing graduates to a binding resolution. The rule *text* stands as previously approved-shaped — the by-reference finding definition, the captain-instruction exclusion, the three classes — and is reproduced in the trigger prose below; only the record/delivery mechanism reframes.

### The triage disposition as an advisory resolution (the shape this entity owns)

The normative spec home for this shape is `durable-gate-approval-pending-blockers/gate-resolution-frontmatter-contract.md`, section "Round records and triage dispositions (advisory)" (contract sha256 681b2348…) — evergreen, component-language; that section is owner-tagged to this entity (02av) per the change protocol, so shape amendments route there, not to a second copy. The text below is the design rationale and the rule semantics that shape serves; it must not diverge from the contract section.

A correction round maps onto the recorder's settled shapes (3k, frozen at `durable-gate-approval-pending-blockers/review/ideation/briefing-15/contract-snapshot.md`):

- The round's **reviewed snapshot** is a **briefing** — immutable, digest-bound (SHA-256 over RFC 8785 canonical bytes), the same object 3k binds for a gate attempt.
- The reviewer's **findings** are **annotations** (with selectors) in that briefing's one ordered log.
- The round's **verdict** is the reviewer's **advisory resolution** — advisory is load-bearing: a round can never advance `status`, so it carries no advancing application (`action: none` territory, h1's layer untouched). 3k already preserves advisory resolutions as first-class, distinct from binding (its AC-13).
- The **ensign's triage** is the ensign's OWN **advisory resolution** on the same briefing. A **review finding** is, by reference, the output of a review instance the active stage's definition declares, or feedback the FO explicitly routes as a review outcome — a direct captain instruction is not a finding but an order to follow or a decision to seek. The resolution's `includes` name each **declined** finding with the three parts the LANDED validation taxonomy already requires of a deferred risk (`docs/dev/README.md:149-151`): its class (correct-but-disproportionate), why it is not material (no value AC at risk; trigger outside the promise), and the condition that promotes it to material. A **material** finding is fixed (the fix is the product change); a **needs-decision** finding escalates to the FO.
- An **all-declines round** is a real advisory resolution recording zero fixes and naming each decline. Absence of a resolution means no finding arrived; a resolution with only declines means every finding was declined — never rendered alike.

Concretely, riding 3k's vocabulary with no schema change (this is spike Claim 2's constructed record; it lives in the round's briefing log in the Subspace room, joined by briefing id — not in entity frontmatter):

```yaml
- type: Annotation                      # the ensign's decline, one per declined finding
  id: annotation:decline-symlink-prototype
  briefing: briefing:02av-impl-round-3
  by: actor:ensign
  includes: [annotation:finding-symlink-prototype]   # the reviewer's finding it declines
  body: >
    class: correct-but-disproportionate; why-not-material: no value AC breaks and the
    crafted-symlink trigger is outside the supported flow (Trust-boundaries carve-out);
    promotes-when: a released user reaches it through an operator-selected repo.
- type: Resolution                      # the ensign's advisory triage verdict for the round
  id: resolution:ensign-02av-impl-round-3
  briefing: briefing:02av-impl-round-3
  by: actor:ensign
  decision: revise                      # advisory only — no application block, status unchanged
  reason: "triage: 1 material fixed; 1 declined"
  includes: [annotation:decline-symlink-prototype]
```

### AC-narrowing / design-reset → a binding resolution (captain-owned)

Declining a disproportionate finding and narrowing a value AC to make a finding pass are opposite moves under the same pressure. The first is the ensign's, recorded as the advisory resolution above. The second weakens the value the entity promised: it graduates to a **binding resolution** — a real gate attempt, captain-owned — so the loop structurally cannot self-approve a reframe. 2ae's Cycle 2 self-reconfirmed a narrowing inside the loop; under this rule that narrowing opens a captain-owned binding gate attempt instead. bw's `AC {unchanged | narrowed}` field still *records* the drift; this entity's rule is what *forces* it onto a captain-owned resolution rather than an in-loop reconfirm.

### The rule text, delivered at the trigger (thinned)

The ensign needs a thin trigger — triage before fixing; record the triage as an advisory resolution — not the ~2 KB normative block the parked cut copied into every routed context. The rigor (what a valid decline contains, its falsifiability) now lives in the record shape and the LANDED validation taxonomy the ensign already triages against, off the always-loaded surface. Delivered via a short pointer in `feedback-rejection-flow` (cross-stage reflow) and the `docs/dev/README.md` implementation bullet (in-stage rounds); `ensign-shared-core` gets nothing (zero always-loaded delta preserved). This is the byte reframe: enforcement moves from "the FO copies normative prose every round" to "the record's class is machine-checkable against the finding's evidence" (AC-2). It genuinely reduces FO-loaded bytes rather than relocating them — the ensign is not handed the full rule to copy each round; it records a structured object the schema validates.

### Storage (per the recorder's principles — borrowed, not built)

Round records live in the entity's review room (append-only, the probes.jsonl pattern 3k already uses); the frontmatter carries the pointer; the body's `### Feedback Cycles` line survives as the human-readable projection, with bw's surface/estimate/AC-drift fields riding the same record. This entity SPECIFIES that shape; the append/pointer/projection are the recorder's generalization to in-stage rounds, DEFERRED beyond the recorder's first implementation — no in-sprint member builds that plumbing. Until that generalization lands, round records are hand-authored into the room plus the projection line.

### Per-mechanism justification (value AC served / simplest alternative / why insufficient)

- **Decline as an advisory resolution (serves AC-1/AC-2):** Alt — the parked `findings` prose field on the `### Feedback Cycles` entry. Insufficient: prose is neither provider-preservable nor machine-falsifiable, and it cost ~2 KB of ratcheted always-loaded contract (the park reason). The advisory-resolution object is per-entity structured metadata with zero prompt-surface weight and a machine-checkable class.
- **Advisory, not binding (serves the boundary + h1's split):** Alt — record the decline as a gate resolution. Insufficient: a gate resolution binds and advances status; a round must never advance status. Advisory is load-bearing and keeps the application layer h1's.
- **AC-narrowing → binding resolution (serves AC-1):** Alt — bw's `AC`-drift prose field alone. Insufficient: it *records* a narrowing but does not *force* it onto a captain-owned gate attempt; 2ae self-reconfirmed one in the loop. The binding-resolution rule is what removes the self-approval.
- **Own the shape, not the plumbing (serves leanness + the 3k bound):** Alt — absorb the rounds-record plumbing. Over-built: it duplicates 3k's recorder and re-grows this into a fourth product (the boundary decision above).

## Out of scope

- **The generic rounds-record plumbing** — the recorder that writes correction-round briefings/annotations/resolutions to the review room, maintains the frontmatter pointer, and projects the `### Feedback Cycles` line. The recorder's own generalization from gate rounds to in-stage rounds; a 3k successor or its own task, per the boundary decision. THE answer the reframe demanded.
- **The application layer** (blockers, eligibility, execution holds, one-use guard) — h1's, by the captain split. A round record is advisory (`action: none`).
- **A new FO-side enforcement check** (e.g. gate refuses to route a REJECTED verdict unless it cites a material finding — digest:158). Real gap, but a new enforcement process: per the captain ruling it needs explicit approval and normally its own entity. Recorded as a candidate follow-up.
- **Any schema extension to 3k's contract.** If the shape needed a new field it would be plumbing; spike Claim 2 shows it does not. Also out: machine-readable materiality / a roborev release-scope schema field / a lint verifying the four-field triage (digest:159).
- **bw's `### Feedback Cycles` entry format itself** (surface/estimate/AC-drift fields, the `git diff --numstat` one-liner). Owned by bw; those fields ride the record as the projection line.
- **The reviewer/production side** — how roborev classifies and what it emits (`.roborev.toml`, present-gate tiering). Unchanged; this is the consumption half.
- **A workflow stakes field.** The stakes member is parked; triage anchors on per-entity value ACs + the committed taxonomy. Do not reintroduce a stakes dependency.

## Riskiest-mechanism spike (done first)

Two claims, each of which would invalidate the rest of the design if false. Both exercised, evidence recorded not asserted.

**Claim 1 — the four-field taxonomy discriminates AND the recorded class is machine-falsifiable, without a stakes field.** (The Case A/B classification stands from the pre-reframe cut; the reframe adds the falsifiability exercise, since AC-2 now checks a recorded class.) Run the LANDED materiality rule (`docs/dev/README.md:149`: material iff a supported/promised/common/observed trigger AND a value-AC/non-negotiable-boundary break) over four fixtures and one red control, recomputing materiality from the four fields and asserting it matches the recorded class:

- **Case A** — symlink-edge-case-on-a-prototype (incident 13, spacedock_subspace codex:019f63c6): trigger not supported (a crafted symlink no supported path produces), no value AC breaks, boundary under the config's Trust-boundaries carve-out (`.roborev.toml:10`) → **not material → declinable**; recorded `declined`. Pass.
- **Case B** — control, `status --boot` drops the taxonomy field: normal-boot trigger (every FO at boot), value AC breaks → **material**; recorded `material`. Pass.
- **85 M1** (real 0260 drift) — FO prompt surface past the byte ceiling: every-FO-boot trigger, `TestFOFunctionPromptSurfaceShrinks` ceiling breaks → **material**; recorded `material`. Pass.
- **85 D1** (real 0260 drift) — a deferred risk, trigger outside the promise → **declinable**; recorded `declined`. Pass.
- **Red** — Case B (material) recorded as `declined`: recompute says material, flag says declined → **the check REJECTS it**. This is the falsifiability AC-2 needs.

Exercised as a throwaway offline check (`scratchpad/materiality_check.py`); result: discriminates + falsifiable, exit 0. No stakes level was consulted; "prototype-ness" entered entirely through the four fields. Seeds the AC-2 fixture directly (Case A/B + 85 M1/D1 with their four-field evidence; the red control proves the check can fail). Had Case A come out Material, the design would need the parked stakes field and I would have escalated; it did not.

**Claim 2 — the triage disposition rides 3k's settled contract with ZERO schema change** (if false, the shape is plumbing this entity cannot own alone, and the boundary answer collapses). Constructed the concrete record in the Proposed approach using only 3k / Review-&-Gate shapes: an ensign-authored **Annotation** (body = class/why-not-material/promote-when; `includes` the reviewer's finding) plus the ensign's **advisory Resolution** (`decision`, `includes` the decline-Annotation, no application). Every field already exists in `gate-resolution-frontmatter-contract.md` — advisory resolutions are already first-class (its AC-13), `includes` already references same-Briefing annotations, and an advisory resolution already carries no application. No field is minted. **Determination:** the shape is expressible in the frozen contract; had it required a new field I would have flagged plumbing and escalated — it did not. The record HOME needs no further spike: it rides 3k's room-ref + probes.jsonl pattern (3k's own proven mechanism).

## Documentation changes (concrete before/after — ideation proposes, implementation applies)

**`skills/feedback-rejection-flow/SKILL.md`** — the parked ~13-line standing block does NOT return. Add a thin pointer only, appended to step 5's findings-routing sentence:

> After step 5's findings-routing sentence, add:
> `Findings are inputs to triage, not a fix list (a direct captain instruction is an order, not a finding). Ask the routed worker to record its triage as an advisory resolution on the round's briefing — never a binding one, a round cannot advance status: material findings fixed; each correct-but-disproportionate finding declined, its resolution naming the finding's class, why it is not material, and what would promote it (the deferred-risk fields the validation stage already requires); a needs-decision finding escalated. A finding neither fixed nor recorded is not triaged. Narrowing a value AC to pass a finding is a captain-owned design-reset (a binding gate attempt), not the worker's to make.`

**`docs/dev/README.md`** — add one bullet in the `implementation` stage-def, after the design-reset bullet (`README.md:124`), composing with the LANDED validation taxonomy:

> `- When consuming an in-stage review round's findings (a roborev panel, a detached-audit pass, or routed gate feedback), triage before fixing and record the triage as an advisory resolution on the round's briefing — advisory only, a round never advances status. The committed taxonomy is the `validation` release-scope classification (Material / Deferred risk / Polish / Needs decision) plus `.roborev.toml`'s four-field evidence record (released user + workflow; observable harm; affected value AC or non-negotiable boundary; trigger evidence). Fix material findings; for a correct-but-disproportionate one, record a decline whose resolution names its class, why it is not material, and its promote-to-material condition; escalate a needs-decision finding to the FO. A value AC narrowed to make a finding pass is a captain-owned binding resolution (a gate attempt), not a fix (see the entry's AC-drift note).`

**The shape's spec home is the recorder's contract, not a second doc.** The triage-disposition-as-advisory-resolution shape (the Annotation/Resolution encoding, the all-declines semantics, and the graduation-is-binding rule) lives in `durable-gate-approval-pending-blockers/gate-resolution-frontmatter-contract.md` §"Round records and triage dispositions (advisory)" — the section owner-tagged to this entity, edited in place per the change protocol. No companion spec ships here: one shape, one normative home (two specs for one shape is exactly what the change protocol forbids). This entity's body carries only the design rationale, the dev taxonomy anchors, and the Case A/B fixtures.

**`skills/ensign/references/ensign-shared-core.md`** — **no change.** The rule reaches the ensign only via the trigger paths above, so an always-loaded pointer would restate at boot what the trigger delivers; net always-loaded contract delta zero.

No other surfaces change (no new FO enforcement check, no schema, no binary; no rounds-record plumbing — see Out of scope). The trigger prose is a thin pointer, not the ~2 KB standing block, and not an enforcement gate.

## Acceptance criteria

**AC-1 (VALUE) — A seeded correct-but-disproportionate finding against a low-stakes fixture entity produces a recorded decline (the ensign's advisory resolution) and a ZERO-line product diff, while a seeded material control produces a fix (non-zero diff) — the taxonomy discriminates, it does not decline everything.**
Verified by: the live replay (the moved 0260 DoD line) in which a dispatched ensign, given the trigger rule and the fixture entity, receives (a) the Case-A symlink-shape finding and (b) the Case-B value-AC-break finding; its worktree shows 0 changed product lines for (a) with a conforming decline recorded as an advisory resolution (hand-authored into the room + `### Feedback Cycles` projection line — interim until the round-record generalization lands, beyond the recorder's first cut), and a non-zero fix for (b). Independent baselines that moved the wrong way: (i) the archived incident-13 history, where the same finding-shape produced a non-zero dutiful fix; (ii) 85's real 0260 bounce, where the deferred-risk disposition existed only as a prose fragment with no first-class object. The number that can move wrong is product LOC on the decline path (0 vs archived >0). Validation owns the live drive, per the live-drive proof rule.

**AC-2 — The decline resolution's class is falsifiable against the finding's own four-field evidence: a decline recorded for a finding whose four fields establish materiality FAILS the check; a decline for a non-material finding passes.**
Verified by: a checked-in fixture (Case A declinable + Case B material + 85's M1 material + 85's D1 deferred, each with its four-field evidence as fixture data) and a small offline check that recomputes materiality from the four fields and asserts it matches the resolution's recorded class, with a red control (a material finding recorded as declined) the check must reject. The expected value is the fixture's four-field data (an independent source that can diverge from the flag), NOT a substring of any instruction file — satisfying the prose-grep ruling. Prototyped and passing in the spike (`materiality_check.py`). Cost: low, no product code.

**AC-3 — The disposition is an ADVISORY record delivered at the trigger, the AC-narrowing move is BINDING, and the rounds-record plumbing is ABSENT from this entity's surface.**
Verified by: a validation-stage audit — (a) the triage disposition is recorded as an advisory resolution carrying no advancing application (never a binding gate resolution; a round never advances status); (b) the contract's §"Round records" Graduation clause names an AC-narrowing/design-reset as a captain-owned binding resolution the loop cannot self-approve; (c) the rule text lands only at the trigger — a thin pointer in `feedback-rejection-flow` and the `docs/dev/README.md` implementation bullet — and `ensign-shared-core` carries no finding-consumption section (an absence read a re-introduction flips); (d) the boundary holds: this entity's shipped surface contains the triage semantics + the owner-tagged contract §"Round records" section (amended in place) and NO rounds-record plumbing (no recorder binary, no room-append/pointer/projection code, no schema field minted beyond 3k's contract). One-off reads at validation, output in the report, not committed prose-grep tests.

## Test plan

- **Value replay (live drive) → AC-1:** the moved DoD live replay — a dispatched ensign declines Case A (zero product LOC, decline recorded as an advisory resolution) and fixes Case B (non-zero). Baselines: incident-13's dutiful fix and 85's prose-only deferred-risk. Validation owns the drive; the fixture entity + two seeded findings come from the spike. Cost: one live dispatch, medium.
- **Offline classification check (fixture + check) → AC-2:** the spike formalized — the four fixtures' four-field evidence + the offline materiality recompute asserting the recorded class, with the red control. Already prototyped green (`materiality_check.py`). No product code. Cost: low.
- **Landing + boundary audit → AC-3:** trigger-delivery reads (thin pointer in `feedback-rejection-flow` + the README bullet), the `ensign-shared-core` absence read, the advisory-vs-binding distinction in the contract §"Round records" section, and the boundary check that no rounds-record plumbing enters this entity's surface. Cost: low, reads only.
- **No Go/binary tests, no product code in this cut** — the deliverable is the thin trigger prose, the contract §"Round records" amendment (owner-tagged), one fixture set, and one offline check; the rounds-record plumbing is the recorder successor's.
- **High-stakes note:** `feedback-rejection-flow` is shipped contract/scaffolding (a high-stakes surface per the Proof policy), so the detached adversarial audit applies at validation before merge.

## Expected surface + tolerance (declared, per captain ruling)

- `skills/feedback-rejection-flow/SKILL.md`: a thin trigger pointer, ~5-7 lines (NOT the parked ~13-line standing block), folded into the step-5 findings-routing amendment.
- `docs/dev/README.md`: +1 bullet in the `implementation` stage-def, ~5 lines, citing the LANDED validation taxonomy + `.roborev.toml` four fields.
- Shape spec: no new doc — the normative home is the recorder's contract §"Round records and triage dispositions (advisory)" (owner-tagged to 02av), amended in place per the change protocol; zero new/ratcheted bytes on this entity's surface.
- `skills/ensign/references/ensign-shared-core.md`: **no change** — net always-loaded contract delta zero.
- 1 fixture set (Case A/B + 85 M1/D1, four-field evidence) + 1 offline check (AC-2).
- **0 Go source files, 0 product LOC.** The rounds-record plumbing (recorder append/pointer/projection) is the deferred recorder-successor concern, explicitly out of scope.
- **Tolerance: 2×**, with a hard self-check: any Go/product code, any rounds-record plumbing appearing in this entity's surface, any schema field minted beyond 3k's contract, a return of the ~2 KB always-loaded standing block, a new FO-contract enforcement check, or a NET-POSITIVE `ensign-shared-core` delta trips a reconfirm. Delivering the SHAPE (not the plumbing) at zero always-loaded cost, and answering the byte objection by moving rigor from copied prose to a machine-checkable record, ARE the point.

## Stage Report: ideation

- DONE: Consumption rule drafted as concrete before/after for ensign-shared-core: classify review findings against the committed triage taxonomy (dev README validation stage + roborev config) BEFORE fixing; the decline disposition — correct-but-disproportionate gets a recorded decline, not a dutiful fix — with the record shape the FO gate checks.
  `## Documentation changes` carries the exact `ensign-shared-core` `## Consuming review findings` section (Material→fix / correct-but-disproportionate→recorded decline / needs-decision→escalate) and the `docs/dev/README.md` implementation-stage bullet citing the `validation` taxonomy + `.roborev.toml` four fields; the record is the `findings` field on the existing `### Feedback Cycles` entry the FO gate already reads (no new enforcement).
- DONE: Generic vs dev-specific explicit: the consumption rule and decline disposition are generic (any reviewer stage, any workflow); the roborev/README taxonomy citations are the dev instances.
  `## Proposed approach` splits them: generic rule → `ensign-shared-core` (names no dev artifact, anchored on the entity's ever-present value ACs); taxonomy/roborev/detached-audit citations → `docs/dev/README.md`. AC-3 makes the no-leak separation a validation-stage landing audit.
- DONE: Written expected surface + tolerance declared; riskiest mechanism spiked first or "no spike needed" recorded.
  `## Expected surface + tolerance`: ~13 lines skill + ~4 lines README + 1 fixture + 1 offline check, 0 product LOC, 2× tolerance with a hard reconfirm on any code / new enforcement / second record shape. `## Riskiest-mechanism spike` ran the four-field taxonomy over the real archived symlink finding (declinable) + a material control (fixed) — it discriminates WITHOUT the parked stakes field, so the design stands; record home needs no spike (rides bw's proven `### Feedback Cycles` convention).

### Summary

Ideated the ensign finding-consumption rule + decline disposition as a prose change to fleet-loaded `ensign-shared-core` (generic) plus a `docs/dev/README.md` implementation bullet (dev instance), correcting the seed's parked-stakes anchor to per-entity value ACs + the committed taxonomy. The decline is one convention, not a second shape: a `findings` field on the sibling bw member's existing `### Feedback Cycles` entry, adjacent to its `AC {unchanged|narrowed}` field, so the two opposite moves under review pressure (decline a finding / narrow the claim) sit side by side — the rule makes declining legal and names narrowing an AC as an illegal-for-the-ensign, captain-visible design-reset. The riskiest assumption (four-field taxonomy suffices without a stakes field) was exercised against the real archived symlink finding and a material control and discriminates correctly; a new FO enforcement check is held out of scope per the captain ruling.

## Stage Report: ideation (cycle 2)

- DONE: Qualify the generic rule's trigger by reference, not enumeration (sole captain gate ask; rest approved-shaped as-is).
  The generic consumption rule and its `ensign-shared-core` doc-diff now define a **review finding** as the output of a review instance the ACTIVE STAGE'S DEFINITION declares (panels/audits/validators named there, carried in the dispatch packet) OR feedback the FO explicitly routes as a review outcome; a direct captain instruction is explicitly excluded (an order to follow / a decision to seek, never a finding to triage). The shared contract references the stage-def instead of enumerating reviewer types; the dev stage-defs remain where concrete instances are named. Estimate updated to ~15 skill lines, inside the declared 2× tolerance; no other change.

### Summary

Folded the single captain ask into the generic rule and its concrete before/after diff: the finding trigger is now defined by reference to the active stage's definition (plus explicit FO-routed feedback), with captain direction excluded, replacing the unqualified "reviewer / staff review / automated panel" enumeration. Surface grew ~2 lines (well within tolerance); the rest of the design — the decline-as-`findings`-field convention, the AC-narrowing sibling, the spike, and the AC set — is unchanged.

## Stage Report: ideation (cycle 3)

- DONE: Move the generic rule OUT of always-loaded `ensign-shared-core` into the feedback delivery path (captain gate rework, item 1).
  `### Where the rule is delivered` designs the concrete mechanism: a standing block in `feedback-rejection-flow` that the FO prepends to the `--feedback-context-file` it already assembles every reflow, delivered to the worker via the existing `### Feedback from prior review` packet emission (`internal/dispatch/build.go:627-629`, gated on `is_feedback_reflow`; FO file-assembly at `fo-dispatch-core.md:129,138`). **Zero new machinery** — both the context-file→packet path and the FO's assembly exist today. Losing alternative named and deferred: a `dispatch build`-emitted block (guaranteed, but product LOC + golden churn), the justified upgrade only if live drives show FO omission. The `## Documentation changes` diff now targets `feedback-rejection-flow` (standing block + step-5 amendment) instead of `ensign-shared-core`.
- DONE: Keep the in-stage half in `docs/dev/README.md` (item 2).
  The implementation stage-def bullet is unchanged (carried in every implementation packet via `show-stage-def`); `template` group propagates the stage-def pattern for non-dev workflows — recorded in the landing spots.
- DONE: `ensign-shared-core` gets nothing, justified (item 3).
  Absence justified by construction: the by-reference finding definition scopes a review finding to exactly the stage-def-declared and FO-routed paths, both of which already deliver the rule with the findings, so an always-loaded pointer would restate at boot what the trigger carries. Net always-loaded delta zero.
- DONE: Update the landing-spot audit AC and the estimate (item 4).
  AC-3 rewritten to a trigger-delivery + absence audit, adding a `dispatch build --feedback-reflow` check that the block reaches the packet (behavior, not prose). Estimate moved the ~13 rule lines from `ensign-shared-core` to `feedback-rejection-flow`, recorded `ensign-shared-core` no-change, and added a NET-POSITIVE-always-loaded-delta trip to the hard self-check.

### Summary

Applied the captain's placement rework: the approved rule text now ships in the feedback delivery path (a `feedback-rejection-flow` standing block the FO includes in the routed `--feedback-context-file`, delivered via the existing `### Feedback from prior review` packet emission — zero new machinery) plus the unchanged `docs/dev/README.md` in-stage bullet, with `ensign-shared-core` getting nothing because the by-reference finding definition provably scopes every finding to a path that already carries the rule. Named and deferred the `dispatch build`-emitted-block alternative behind observed FO-omission drift (bw's convention-first ordering). The rule TEXT, the decline-as-`findings`-field convention, the AC-narrowing sibling, the spike, and AC-1/AC-2 are unchanged; only AC-3 and the estimate moved with the delivery.

## Stage Report: implementation

- DONE: ONE composed `### Feedback Cycles` entry format lands in skills/feedback-rejection-flow/SKILL.md carrying BOTH members' fields — including the all-declines case, which must render as a real recorded state and not as an empty field.
  This entity owns the `findings` field, adjacent to bw's `AC {unchanged | narrowed}` on the same entry (commit ee4a23cf): `findings {none | {F} fixed, {D} declined: <ref · class · why not material · promotes when>}`. `findings none` means no finding arrived; an all-declines round records `0 fixed` and names each decline with its grounds and promote condition, and the section prose states outright that the two must never read alike. The extension that a round records an entry whenever it produces a disposition — fixes OR declines — landed with it, so an all-declines round is not silently unrecorded.
- DONE: Both members' contract prose lands in the files each declared and NOWHERE else; skills/ensign/references/ensign-shared-core.md gets NOTHING — 02av's net always-loaded delta is zero and that is one of its acceptance criteria.
  `## Finding-triage block` (the approved rule text, relocated) plus the findings-routing amendment in `feedback-rejection-flow/SKILL.md`, and the triage bullet in `docs/dev/README.md`'s implementation stage-def citing the `validation` release-scope taxonomy and `.roborev.toml`'s four fields. `git diff --numstat HEAD~1` shows `skills/ensign/` absent from the change set: **always-loaded delta 0**, no new FO enforcement check, no second record section.
- FAILED: BYTE-FUNDED and suite-green: `go test ./...` passes INCLUDING TestFOFunctionPromptSurfaceShrinks, with per-file byte accounting recorded.
  Green on every suite except the ratchet (125,564 vs baseline 122,634). This entity's measured share is **+2,192 bytes**: the standing block +2,032 and the routing amendment +160. Composed headroom was 483. Not funded; reported, not absorbed.

### Where this entity's bytes went, and why they could not be moved

The standing block is ~2.0 KB of ensign-facing text sitting inside a file the FO prompt-surface ratchet measures. Relocating it to an unmeasured sibling reference would not reduce what the FO actually loads — it must read the block to copy it into the routed context — so that move only hides the cost from the ratchet, and the dispatch forbids it. The `dispatch build`-emitted alternative (guaranteed delivery, no FO-loaded bytes) is this entity's own recorded losing alternative, deferred as product LOC behind observed FO-omission drift; it is the option that would genuinely remove these bytes from the measured surface, and it is a captain call, not an implementation one.

A duplication scan across all 13 ratcheted files before cutting anything recovers 110 bytes in total — the measured set is already deduplicated. The compression actually taken: the block was tightened against the ideation text without dropping a clause (the by-reference finding definition, the captain-instruction exclusion, the three classes, the decline's four required parts, the record-or-it-is-not-triaged rule, and the AC-narrowing clause are all intact), and bw's side self-funded −248 from the skill's description and intro.

Delivery mechanism unchanged from the approved design: the block heads the `--feedback-context-file` the FO already assembles, which `dispatch build` already emits as the packet's `### Feedback from prior review`. No new machinery. AC-3's `dispatch build --feedback-reflow` delivery exercise is validation's.

### Summary

Landed the consumption rule and decline disposition at the trigger — a standing block in the feedback delivery path plus the `docs/dev/README.md` in-stage bullet — with `ensign-shared-core` untouched, so the always-loaded contract gains nothing, as designed. The decline is a field on bw's correction-round entry, not a second section: an all-declines round now has a real recorded state instead of an empty field, which is the case this entity exists for. The one unmet criterion is byte funding: +2,192 measured bytes with no honest trim available in the ratcheted files, which needs a captain decision (re-baseline, or promote the deferred `dispatch build` block that would carry the rule at zero FO-loaded cost) rather than a trim that deletes meaning.

## Park record (captain decision, 2026-07-21)

**PARKED for this train — a redesign candidate, not a rejection.** The rule and its all-declines requirement stand as gated; only the delivery mechanism is in question. The implementation was landed, suite-green on every check except the byte ratchet, and then withdrawn from the branch when the captain shipped the sibling `bw` alone.

**Reason.** The deliverable is correct and validated in design, but its standing block costs ~2.0 KB on a ratcheted contract surface whose total recoverable redundancy measured 110 bytes across all 13 files. The captain identified a better home: **3k's advisory gate-resolution record**, where the findings disposition becomes per-entity structured metadata rather than always-loaded contract prose. That reframes the cost — a per-entity record carries no standing prompt-surface weight at all, which is the objection this member could not answer by trimming.

**What the withdrawal removed** (commit 457b910d on `spacedock-ensign/feedback-cycle-record-command`, reverting the parts of ee4a23cf this member owned): the `## Finding-triage block` standing section and its findings-routing amendment in `skills/feedback-rejection-flow/SKILL.md`; the `findings` clause in the entry format and the all-declines spec prose; the implementation stage-def triage bullet in `docs/dev/README.md`. `skills/ensign/references/ensign-shared-core.md` was never touched, so the zero-always-loaded-delta criterion held throughout.

### Specification preserved for a 3k-based redesign

This is the design the composition proved out, kept whole so the park costs no design work. A redesign implements this record, not this prose.

**The disposition field, on the correction-round entry** — sibling to bw's `AC {unchanged | narrowed}`, so the healthy response to review pressure and the pathological one sit adjacent in one record:

    findings {none | {F} fixed, {D} declined: <ref · class · why not material · promotes when>}

**The all-declines semantics** — the case this member exists for. `findings none` means no finding arrived. An all-declines round records `0 fixed` and names each decline with its class, why it is not material, and its promote condition. "Nothing was found" and "everything found was declined" are different facts and must never render alike; an all-declines round is the ideal outcome and must leave a real recorded state, not an empty or absent field.

**The recording extension** — a round records an entry whenever it produces a disposition, fixes **or** declines, not only when it triggers another pass. Without this an all-declines round leaves no record at all.

**The rule text** (validated whole, and it passed every contractlint anchor, host-neutrality, and ordered-procedure check while it was on the branch): review findings are inputs to triage, not a fix list; a review finding is defined **by reference** — the output of a review instance the active stage's definition declares, or feedback the first officer explicitly routes as a review outcome — and a direct captain instruction is explicitly NOT a finding but an order to follow or a decision to seek; classification runs against the workflow's declared taxonomy where one exists plus the entity's own value ACs; **material** (breaks a value AC or a declared non-negotiable boundary reachable through the supported workflow) is fixed; **correct-but-disproportionate** gets a recorded decline naming the finding, its class, why it is not material, and what would promote it; **needs decision** escalates to the first officer rather than being resolved privately; a finding neither fixed nor recorded is not triaged; and **narrowing an acceptance criterion to make a finding or rejection pass is not a licensed disposition** — it is a design-reset event requiring the captain's sign-off, recorded so it is captain-visible, never a task-internal edit.

**Dev-instance anchor** (unchanged): the committed taxonomy is the `validation` stage-def release-scope classification plus `.roborev.toml`'s four-field evidence record — released user and workflow; observable harm; affected value AC or non-negotiable boundary; trigger evidence. The spike's Case A (the archived symlink finding → declinable) and Case B (a value-AC break → material) remain the AC-2 fixture seed and still discriminate without a stakes field.

**Carried forward unchanged:** AC-1's live replay (a decline path with zero product LOC against incident 13's non-zero dutiful fix), AC-2's offline classification check, and the out-of-scope boundaries — no FO enforcement check, no schema or lint, no second record section, no stakes dependency. AC-3's landing-placement audit is the one criterion a 3k-based redesign would restate, since the delivery surface is exactly what changes.

### Feedback Cycles

- 2026-07-20T16:20:45Z — captain design-reset, recorded before any further dispatch (the landed convention's own rule): PARK this entity's implementation before landing — bw lands alone in wave 3. Reframe at ideation: record round dispositions, including the decline, as ADVISORY resolutions under the gate-recorder model (see the 3k scope cut of 2026-07-21) instead of the findings field on the entry; the finding-triage rule TEXT stays approved-shaped, only the record/delivery mechanism reframes. surface 0 landed vs estimate ~17 lines (held before landing); AC unchanged.
- Cycle 1: REJECTED — Roborev job 548; surface 6 files/91 added lines vs estimate 5 files/60–95 lines (96% of upper bound); AC unchanged
- Cycle 2: REJECTED — Roborev job 554; surface 6 files/95 added lines vs estimate 5 files/60–95 lines (100% of upper bound); AC unchanged

## Reframe brief (captain-directed, 2026-07-21)

The triage rule text, the three-class taxonomy, the AC-narrowing design-reset clause, and the Case A/B spike all stand unchanged. What reframes is the RECORD AND DELIVERY mechanism, from a findings field on the prose entry to the gate-recorder model (see the 3k scope cut, 2026-07-21 — its contract doc is the design authority):

- A correction round's reviewed snapshot is a briefing (digest-bound); the reviewer's findings are annotations with selectors; the round's verdict is the reviewer's ADVISORY resolution — advisory is load-bearing: a round can never advance status.
- The ensign's triage is the ensign's OWN advisory resolution whose includes name the declined findings with class and why-not-material — the decline becomes a first-class, provider-preservable object instead of a prose fragment.
- A design-reset or AC-narrowing event graduates to a BINDING resolution — a real gate attempt, captain-owned. The loop structurally cannot self-approve a reframe.
- Storage per the recorder's principles: round records live in the entity's review room (append-only, the probes.jsonl pattern), frontmatter carries the pointer, the body's `### Feedback Cycles` line survives as the human-readable projection. Surface/estimate/AC-drift stay Spacedock-side fields riding the record.

The boundary question this ideation MUST answer (not assumed): does this entity own only the TRIAGE semantics (rule text + decline-as-advisory-resolution shape), with the generic rounds-record plumbing (round briefings, room layout, projection) as bw's deferred-machinery successor or its own task — or does it absorb the plumbing? Do not re-grow one entity into four products; the 3k split is the cautionary precedent.

Design inputs at reframe time: bw's LANDED wording (never ideation-time quotes), the 0260 Commander run's real hand-authored `### Feedback Cycles` entries as drift evidence and fixtures, and 3k's post-cut approved contract. Hard dependency: 3k's gate. Sprint membership and the index DoD line-36 disposition are a captain decision recorded separately.

## Stage Report: ideation (cycle 4)

- DONE: The reframe brief's boundary question answered explicitly and honestly: does this entity own triage semantics only (rule text + decline-as-advisory-resolution shape) or absorb the generic rounds-record plumbing — with the 3k four-products precedent as the cautionary bound.
  `## Boundary decision (reframe — the spine)`: TRIAGE SEMANTICS ONLY (rule text + decline-as-advisory-resolution shape + all-declines semantics); the generic rounds-record plumbing (room append, frontmatter pointer, projection) is NOT absorbed — it is identical recorder machinery generalized from gate rounds and belongs with a 3k successor / its own task. The 3k four-products scope-cut is cited as the cautionary bound. The honest test — "if the shape needed a new 3k schema field it would be plumbing" — is exercised by spike Claim 2 and comes out clean (no field minted). The application layer stays h1's; the reviewer/production side stays roborev's.
- DONE: Dispositions (including the decline) designed as ADVISORY resolutions against 3k's settled contract and bw's LANDED wording in skills/feedback-rejection-flow/SKILL.md and docs/dev/README.md — never ideation-time quotes; real 0260 hand-authored Feedback Cycles entries banked as fixtures.
  The disposition rides 3k's frozen contract (briefing-15 snapshot): reviewed snapshot → briefing (digest-bound); findings → annotations; round verdict → the reviewer's advisory resolution; the ensign's triage → its own advisory resolution whose `includes` name each declined finding with class/why-not-material/promote-when (the deferred-risk fields the LANDED `README.md:149-151` validation taxonomy already requires). Concrete YAML constructed with zero schema change (spike Claim 2). Designed against `feedback-rejection-flow/SKILL.md` (FO owns `### Feedback Cycles`, tracks cycles) and the landed README taxonomy — not ideation-time quotes. Real 0260 rounds banked as fixtures: 85's M1 (material) + D1 (deferred) seed AC-2; 2ae's self-reconfirmed AC-narrowing is the binding-resolution motivating case.
- DONE: Expected surface + tolerance declared; value ACs measure against baselines that can move the wrong way (the moved 0260 DoD line: seeded disproportionate finding -> recorded decline + zero-line diff in live replay).
  `## Expected surface + tolerance`: thin trigger pointer (~5-7 lines, not the ~13-line block) + ~5-line README bullet + the shape's spec home amended in place in the recorder's contract §Round records (no companion doc, owner-tagged to 02av) + 1 fixture + 1 offline check, 0 product LOC, `ensign-shared-core` unchanged; 2× tolerance with a reconfirm trip on any plumbing / schema-mint / 2 KB-block-return / net-positive `ensign-shared-core` delta. AC-1 measures the moved 0260 DoD line (seeded disproportionate finding → recorded decline as an advisory resolution + zero product LOC) against incident-13's >0 dutiful fix and 85's prose-only deferred-risk. AC-2's falsifiability was exercised offline in the spike (`materiality_check.py`: discriminates + falsifiable, exit 0).

### Summary

Applied the captain's advisory-resolution reframe: the triage disposition moves from a `findings` prose field on the `### Feedback Cycles` entry to the ensign's own advisory resolution under the gate-recorder model (3k), whose `includes` name each declined finding with the class / why-not-material / promote-when fields the landed validation taxonomy already requires — answering the byte objection that parked the last cut, since a per-entity structured record carries no always-loaded weight where the ~2 KB prose block did. The boundary question is answered honestly: this entity owns triage SEMANTICS only (rule text + shape), not the generic rounds-record plumbing, which is the recorder's own generalization and a 3k-successor / own-task concern — the 3k four-products cut is the cautionary bound, and spike Claim 2 proves the shape rides 3k's frozen contract with zero schema change, so it is genuinely semantics, not a fourth product. The rule text, three-class taxonomy, AC-narrowing-is-a-captain-binding-resolution clause, and the Case A/B spike all stand; the falsifiability of the recorded class was exercised offline over the real 85 drift fixtures and discriminates.

## Stage Report: implementation (cycle 2)

- DONE: Ship the finding-consumption rule and decline-as-advisory-resolution semantics in the approved prose/fixture surface, with zero product LOC and no generic rounds-record plumbing.
  Commit d3efaea0 changes only the three approved prose surfaces plus the fixture/check pair; `git diff --numstat "$(git merge-base main HEAD)"..HEAD` contains no Go or recorder-plumbing path.
- DONE: Prove a seeded correct-but-disproportionate finding yields an explicit recorded decline and zero-line product diff, while no-findings and all-declines remain observably distinct.
  `docs/specs/check-finding-triage-materiality.sh` accepts Case A, rejects the material-as-declined red control, and the contract distinguishes absent triage Resolution from a present zero-fix all-declines Resolution; changing the red expectation to accept makes the check fail.
- DONE: Keep material findings fixable, needs-decision findings escalated, and value-AC narrowing captain-owned; preserve the existing validation taxonomy and recorder schema.
  The trigger and implementation-stage rules route the three dispositions, while the owner-tagged contract section adds no schema field and makes AC narrowing a captain-owned binding gate attempt.
- FAILED: Required repository suites are green.
  `go test ./...` fails only `TestFOHostPromptLoadRatchet`: the approved thin trigger pointer adds 703 bytes to each host load; changing a Go baseline is outside the declared zero-Go boundary and the prior park records only ~110 honest redundant bytes.
- SKIPPED: Request Roborev on the completed commit and durably triage every finding.
  Commit d3efaea0 is preserved for inspection but is not a completed commit while the material ratchet failure awaits a captain scope decision, so Roborev was not requested prematurely.

### Summary

The triage semantics, advisory decline shape, all-declines distinction, and falsifiable four-field fixture are implemented in commit d3efaea0. Actual surface is 0 production LOC, 40 fixture/check lines, and 48 docs/skill/spec lines versus the declared 0 / 35–55 / 25–40 with 2× tolerance; no hard scope boundary was crossed, but the 703-byte FO ratchet failure requires an explicit decision before implementation can complete.

## Stage Report: implementation (cycle 3)

- DONE: Ship the finding-consumption rule and decline-as-advisory-resolution semantics in the approved prose/fixture surface, with zero product LOC and no generic rounds-record plumbing.
  Commits d3efaea0, 059b12d8, 9b2093b5, and e85eb0cf ship the thin trigger, owner-tagged contract shape, offline fixture/check, exact captain-authorized +703 host baselines, and bounded evidence fixes; no product or recorder-plumbing path changed.
- DONE: Prove a seeded correct-but-disproportionate finding yields an explicit recorded decline and zero-line product diff, while no-findings and all-declines remain observably distinct.
  The offline check accepts Case A with zero product LOC, rejects both material-as-declined and unknown-class red controls, and the contract distinguishes absent triage Resolution from a present zero-fix all-declines Resolution; validation owns the live replay.
- DONE: Keep material findings fixable, needs-decision findings escalated, and value-AC narrowing captain-owned; preserve the existing validation taxonomy and recorder schema.
  The trigger and implementation rule fix material findings, escalate needs-decision findings, and graduate AC narrowing to a captain-owned binding gate attempt; the advisory shape mints no schema field.
- DONE: Reconcile implementation surface to the declared boundary and commander reconfirmation.
  Actual surface is 0 production LOC, 47 test/fixture lines, and 48 docs/skill/spec lines versus declared 0 / 35–55 / 25–40 with 2× tolerance; the commander reconfirmed only the exact +703-byte all-host rebaseline after the cycle-2 blocker, realized by three literal test-baseline edits in 059b12d8.
- DONE: Run required checks.
  `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` are green; the offline check accepts eight valid cases and rejects its two red controls, and deleting any materiality conjunct makes its isolated control fail.
- DONE: Request Roborev on the completed commit and durably triage every finding against the value ACs and declared stakes.
  Branch-final panel jobs 548, 554, and 560 reviewed the full branch; Cycles 1 and 2 were recorded before each fix in state commits 584cae01 and 55473ae8, with AC unchanged and surface at 96% then 100% of the approved upper bound.

### Roborev advisory triage

- Job 548 advisory Resolution: fixed the material AC-2 evidence defect by isolating all four materiality conjuncts. Declined class rename (the captain-reconfirmed disposition class is an AC input; promote if the AC changes or a downstream enum requires it), durable escalation annotations (new record semantics outside the boundary; promote if a supported consumer loses escalation or an AC requires durability), and malformed external-fixture hardening (checked-in-only input; promote if external fixtures become supported).
- Job 554 advisory Resolution: fixed the material AC-2 gap that accepted unknown non-material classes, while declining the repeated class rename, runtime/CI smoke expansion (validation owns the live drive; promote if that drive fails or the AC changes), and zero-byte external-fixture hardening under the same recorded grounds.
- Job 560 all-declines advisory Resolution: **0 fixed; 3 declined**, observably distinct from no findings. Declined taxonomy/action separation because it replaces the captain-reconfirmed class semantics (same promotion condition); declined per-row `NF` hardening because the sole checked-in fixture is well-formed (promote on external fixture support); declined the symmetric non-material-as-material red control because AC-2 requires the opposite material-as-declined control and the existing expression already rejects the symmetric mismatch (promote if validation cannot reproduce that rejection).

### Summary

Implementation is complete at e85eb0cf with zero product LOC and no schema or rounds-recorder plumbing. The exact prompt-ratchet exception is captain-authorized and suite-green, material review evidence defects were fixed without exceeding the 95-line ceiling, and all proportionate declines—including Roborev job 560's all-declines round—are durable with explicit promotion conditions.

## Stage Report: validation

- DONE: Independently validate exact candidate e85eb0cf and the complete four-commit branch against the current entity, approved reframe, captain reconfirmation, and sprint release line; do not rely on the implementation report alone.
  Exact clean head `e85eb0cfcc3c243fd94754be2baafa23be302a21` is four commits over merge-base `fa240a76`; direct diff, entity, durable-decisions DoD, and captain-reconfirmed class/ratchet inputs agree.
- DONE: Reproduce AC-1 as the required live value replay: seed the correct-but-disproportionate finding, exercise the supported triage delivery/record path, and prove an explicit advisory decline with zero product-line fix against the recorded incident-13 nonzero dutiful-fix and 85 prose-only baselines. Do not skip this live replay.
  A fresh routed worker read the shipped trigger: Case A produced a linked decline Annotation plus advisory Resolution and `0/0` product lines; material Case B fixed `product/status.txt` by `+1/0`; both statuses remained `implementation`, unlike the recorded incident-13 fix and 85 prose-only disposition.
- DONE: Attack AC-2 with the offline materiality fixture/check: independently reproduce all valid cases and both red controls, reject unknown classes, and use claim-breaking edits to prove each materiality conjunct and class allowlist is load-bearing.
  `docs/specs/check-finding-triage-materiality.sh` returned 8 ACCEPT/2 REJECT; removing user/workflow, harm, boundary, or trigger independently failed its paired control, and removing the allowlist admitted `defered-risk` and failed.
- DONE: Validate AC-3 end to end: the disposition is advisory and non-advancing; AC narrowing/design reset is captain-owned binding; trigger delivery is present; ensign-shared-core remains unchanged; no product code, recorder round plumbing, new schema field, or second record section entered the diff.
  Live status stayed `implementation`; the trigger and owner-tagged contract carry advisory/no-application semantics and captain binding, while exact-range path/field audits found none of the forbidden surfaces and zero `ensign-shared-core` delta.
- DONE: Verify the exact captain-authorized +703-byte all-host prompt-ratchet rebaseline and prove there was no additional prompt growth or semantic expansion beyond that measured exception.
  The sole changed loaded path is `feedback-rejection-flow/SKILL.md`, `3329→4032` bytes; measured host loads equal Claude `96081`, Codex `75296`, Pi `71426`, each exactly baseline +703, and the ratchet test is green.
- DONE: Audit actual surface against the reconfirmed boundary: zero production LOC, 47 fixture/test lines, 48 docs/skill/spec lines, no more than the approved 95 added-line ceiling after the two recorded correction rounds.
  `git diff --numstat main...HEAD` reports exactly 95 insertions: 0 production LOC, 47 fixture/test insertions (44 fixture/check + 3 test-baseline replacements), and 48 docs/skill/spec insertions across six paths.
- DONE: Verify Roborev jobs 548, 554, and 560 and their durable triage: material evidence defects fixed; every decline names why it is not material and a promotion condition; job 560's all-declines outcome is observably distinct from no findings.
  `roborev show 548/554/560` matches commits `9b2093b5` and `e85eb0cf`; state commits `584cae01`/`55473ae8` precede fixes, and the cycle-3 advisory record gives every decline grounds/promote conditions with job 560 recorded as `0 fixed; 3 declined` rather than absent.
- DONE: Run the detached adversarial audit required for the shipped feedback-rejection skill/contract surface, then run the focused offline check, gofmt -w ./cmd ./internal, go test ./..., go test ./... -race, git diff --check, and cleanliness checks at the exact committed head.
  Detached checkout at `e85eb0cf` was CLEAN; focused check, formatter no-diff, full suite, race suite, and `git diff --check` all exited 0, and the implementation worktree stayed clean at the exact head.
- DONE: AC-1 (VALUE) — live behavior discriminates rather than declining everything.
  Case A's valid four-record JSONL has an explicit linked advisory decline and zero product diff; Case B's material control has a nonzero product fix, with no entity-status advance in either arm.
- DONE: AC-2 — the recorded decline class is falsifiable against independent four-field evidence.
  Both required red controls reject, and five independent claim-breaking mutations each exit 1 on the specific case whose evidence or allowlist they weaken.
- DONE: AC-3 — advisory authority, binding design reset, trigger delivery, and the negative scope boundary hold together.
  Contract/trigger audit plus the live non-advance observation prove the current semantics; exact-range audit confirms no application/plumbing/schema/product expansion.
- DONE: PASSED recommendation with deferred risks separated from material findings.
  No material finding remains. Deferred only: malformed/external fixture hardening (promote on external support/schema AC), CI auto-wiring (promote if release proof requires it), and advisory `decision` consumer ambiguity (promote when round plumbing branches on it).

### Summary

Fresh validation recommends **PASSED** for exact candidate `e85eb0cf`. The required routed live replay moved the value baseline—an explicit advisory decline with zero product-line fix versus a nonzero material control—and the offline evidence, prompt bytes, scope ceiling, Roborev triage, detached audit, full/race suites, and clean-head checks all held independently. No material finding remains; the three detached-audit risks stay deferred under the explicit promotion conditions above.
