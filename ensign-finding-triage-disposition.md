---
id: 02avdajaz0q3hnjwycm5fq45
title: Ensigns triage review findings against declared stakes before fixing — decline disposition for correct-but-disproportionate findings
status: ideation
source: "0260 shaping — agent-derail forensics audit, 2026-07-19."
score: "0.7"
sprint: 0260-proportionality
group: triage
started: 2026-07-20T05:04:07Z
---

`ensign-shared-core` contains zero guidance on consuming review findings — the exact actor that dutifully fixes a symlink edge case in a prototype has no rule to consult, and no disposition short of fixing exists for a substantively-correct-but-disproportionate finding. This adds the generic consumption rule (classify each finding against the workflow's committed finding-triage taxonomy AND the entity's own value ACs before fixing) and the decline disposition (a correct-but-disproportionate finding gets a recorded decline, not a dutiful fix), with the decline recorded as a field on the same `### Feedback Cycles` correction-round entry the gate already reads. The dev instance anchors on the committed taxonomy — the `validation` stage-def release-scope classification plus `.roborev.toml`'s four-field evidence record (this workflow's port of spacedock-subspace's triage). Anchor correction from the seed: the stakes member is **parked**, so triage keys on per-entity value ACs and the committed taxonomy, not a workflow stakes field.

## Problem

- **No ensign-side consumption rule exists.** `skills/ensign/references/ensign-shared-core.md` carries zero guidance on triaging review findings before fixing them (forensics digest GAPS, `_evidence/0260-agent-derail-forensics/remedy-analyses-digest.txt:157`). The actor at the center of the incident class — the ensign that dutifully fixes a symlink edge case on a prototype (spacedock_subspace codex:019f63c6, synthesis incident 13) — has no rule licensing any disposition short of fixing. The audit's stated biggest gap.
- **The pressure is one-sided toward fixing.** The user-global `~/.claude/CLAUDE.md` pushes every worker hard toward fixing every reported finding ("ALL TEST FAILURES ARE YOUR RESPONSIBILITY", "Fix broken things immediately"); roborev's verdict mechanics fail a job on any `## Review Findings` entry (digest:151); and `## Review Findings` "never suppresses a real observation" keeps correct-but-disproportionate findings continuously in front of an actor "that has no consumption rule for them" (digest:153-154). An ensign holding a finding "has no licensed way to not fix it."
- **The deferred-risk record has no consumer.** roborev's config lets a reviewer defer with a trigger / why-outside-promise / promotion-condition note, but "nothing consumes the record: no tracker, entity field, or FO loop step ever revisits" it, so deferred risks "live only in gate messages" (digest:148). The decline needs a durable home the gate reads.
- **Fixing is not the only over-response; narrowing the claim is its twin.** The 0.25.1 incident (synthesis addendum) is the same repair-forward failure expressed as document edits: under a correct rejection, the task narrowed its value AC until a weaker claim passed. Declining a disproportionate finding and narrowing the claim it targets are opposite moves under the same pressure — the rule must make the first legal and the second captain-visible.
- **Not a suppression scheme.** The taxonomy already exists to classify, not silence (digest:165). This entity is the missing *consumption* half: production (roborev classifies + emits) already ships; the ensign that *receives* the findings still has no re-triage-and-decline rule.

## Proposed approach

Two layers, cleanly separated (checklist item 2). The **generic** consumption rule and decline disposition land in fleet-loaded `ensign-shared-core` and name no dev artifact; the **dev-specific** taxonomy citations and review instances land in `docs/dev/README.md`.

### Generic consumption rule (lands in `ensign-shared-core`, workflow-agnostic)

Review findings are inputs to triage, not a fix list. A **review finding** is the output of a review instance the active stage's definition declares — the panels, audits, or validators it names, carried to the ensign in the dispatch packet — or feedback the first officer explicitly routes to the ensign as a review outcome. A direct instruction from the captain is not a finding: it is an order to follow or a decision to seek, never something to triage away. The shared contract qualifies the trigger by *reference* (the stage definition), not by enumerating reviewer types it cannot know — the dev instance's stage-defs are where the concrete panels/audits/validators are named. Before changing anything in response to such a finding, the ensign classifies each against the workflow's declared finding-triage taxonomy (where one is declared) and this entity's own value acceptance criteria:

- **Material** — breaks a value AC, or a declared non-negotiable boundary (safety, security, data-integrity, compatibility) reachable through the supported workflow → **fix it**.
- **Correct-but-disproportionate** (deferred risk or polish) — substantively right, but no value AC breaks and its trigger is outside the supported/promised workflow → **record a decline; do not fix**. The decline is the licensed disposition, not a dodge: name the finding, its class, and why it is not material (no value AC at risk; trigger outside the promise; the condition that would promote it to material).
- **Needs decision** — a genuine product/compatibility fork → **escalate to the first officer**; do not resolve it privately.

The anchor is always present: every entity has value ACs, so the rule has something to triage against even where no workflow-level taxonomy is declared (this is why the parked stakes member does not block it — the entity's own value ACs plus any committed taxonomy are the anchor the captain named).

### Dev-specific realization (lands in `docs/dev/README.md`, cites dev artifacts)

- **The committed taxonomy** the dev ensign triages against = the `validation` stage-def release-scope classification (Material / Deferred risk / Polish / Needs decision; `README.md:146-152`) + `.roborev.toml`'s four-field evidence record (released user + workflow; observable harm; affected value AC or non-negotiable boundary; trigger evidence; `.roborev.toml:34-56`).
- **The review instance** = a roborev panel, a detached-audit pass, or routed gate feedback (the dev instances of "an automated panel / staff review").
- **The decline record** = an entry in the entity's `### Feedback Cycles` section (the dev/FO-owned correction-round record).

### The decline record: one convention, not a second shape

Weighed per the dispatch. The decline is recorded as a **disposition field on the same `### Feedback Cycles` correction-round entry** the sibling `feedback-cycle-record-command` (bw, `reframe` group) already defines — not a new `### Declined Findings` section. That entry already carries one disposition-under-pressure field, `AC {unchanged | narrowed: <note>}`; the finding disposition is its sibling field on the same entry:

    - {timestamp} — {reviewer/loop} {verdict}; surface {actuals} vs estimate {declared} ({P}%); findings {N material→fixed; M declined: <ref · class · why-not-material>}; AC {unchanged | narrowed: <note>}

Co-ownership is explicit: bw owns `surface` + `AC`-drift; this entity owns the `findings` field. They compose into one entry, one section, one convention — the two "opposite moves under the same pressure" (decline a finding / narrow an AC) sit as adjacent fields so a reader sees the healthy and the pathological response side by side. This extends bw's in-stage-round rule minimally: a review round records a `### Feedback Cycles` entry whenever it produces a disposition (fixes **or** declines), not only when it triggers another pass — otherwise an all-declines round (the ideal outcome) would leave no record. "The gate checks it" means the FO's existing gate body-read surfaces the `findings` field; **no new FO enforcement check** is added here (that would be a new enforcement process — captain ruling requires explicit approval + its own entity; recorded in Out of scope as a candidate follow-up).

### The AC-narrowing sibling (0.25.1): opposite move, made captain-visible

The consumption rule names the twin explicitly: narrowing a value AC to make a finding or rejection pass is **not** a licensed ensign disposition — it is a design-reset event requiring the captain's sign-off, recorded (the entry's `AC {… | narrowed}` field, bw's) so it is captain-visible, never a task-internal edit. Declining is the ensign's to make; narrowing escalates. The detection/recording of AC-drift is bw's field; this entity's contribution is the *rule* that classifies the narrowing move as illegal-for-the-ensign and routes it to the captain, stated in the same breath as the legal decline.

### Per-mechanism justification (value AC served / simplest alternative / why insufficient)

- **The generic consumption rule + decline disposition (serves AC-1):** Alt — rely on roborev's output rule alone (only material findings land under `## Review Findings`). Insufficient: that is the *production* side; when the reviewer mis-files a deferred risk as material, or the user-global "fix everything" pressure bites, the ensign still has no licensed re-triage/decline (digest:158, 154). The rule is the consumption half that was missing.
- **Decline-as-a-field on the `### Feedback Cycles` entry (serves AC-1/AC-2):** Alt — a dedicated `### Declined Findings` section. Insufficient/over-built: a second record shape duplicates the correction-round record bw already ships and splits "how the ensign responded to review pressure" across two sections; one entry with sibling disposition fields is the smaller, single-convention answer the sprint's own thesis prefers.
- **The AC-narrowing-is-illegal clause (serves AC-1):** Alt — leave AC-drift to bw's field alone. Insufficient: bw's field *records* a narrowing; it does not tell the ensign the narrowing move is off-limits. Naming it as the decline's illegal twin, in the consumption rule, is what makes the ensign escalate instead of edit.

## Out of scope

- **A new FO-side enforcement check** (e.g. gate refuses to route a REJECTED verdict unless it cites a material finding — digest:158). Real gap, but a new enforcement process: per the captain ruling it needs explicit approval and normally its own entity. Recorded as a candidate follow-up; here the decline is a durable artifact the FO's *existing* gate review reads.
- **Machine-readable materiality / a roborev release-scope schema field / a lint verifying the four-field triage** (digest:159). This entity ships the ensign-side prose rule + a record convention, not new schema or binary surface.
- **The reviewer/production side** — how roborev classifies and what it emits (`.roborev.toml`, present-gate tiering). Unchanged; this is the consumption half.
- **A workflow stakes field.** The stakes member is parked; triage anchors on per-entity value ACs + the committed taxonomy. Do not reintroduce a stakes dependency.
- **bw's `### Feedback Cycles` entry format itself** (surface/estimate/AC-drift fields, the `git diff --numstat` one-liner). Owned by bw; this entity adds only the `findings` field and the all-declines-still-records extension, coordinated at the shared gate.

## Riskiest-mechanism spike (done first)

**Claim under test (what would invalidate the rest of the design if false):** the committed four-field taxonomy + the entity's own value ACs suffice to classify the archived correct-but-disproportionate finding as **declinable**, WITHOUT a workflow stakes field (which is parked) — and still force a genuinely material finding to be **fixed** (the rule must discriminate, not "decline everything"). Exercised as a classification pass over the real archived finding and the real committed `.roborev.toml` triage, evidence recorded not asserted:

**Case A — the symlink-edge-case-on-a-prototype finding** (incident 13, spacedock_subspace codex:019f63c6), run through `.roborev.toml:38-42`'s four fields against a plausible prototype entity's value ACs:
- *Released user + normal workflow:* the supported flow operates on operator-selected repos; a symlink path-escape needs a crafted symlink no supported path produces. None.
- *Observable harm:* in a prototype with no released users, nothing lost in the supported flow; harm is hypothetical.
- *Affected value AC or non-negotiable boundary:* no value AC breaks; a crafted symlink in the operator's own repo falls under the config's own Trust-boundaries carve-out ("a malicious local OS user or deliberately compromised agent process is outside the CLI's threat model", `.roborev.toml:10`) — no in-model boundary at risk.
- *Trigger evidence:* reachable only via crafted input the supported flow never produces; adversarial-only.
- → **Deferred risk / Polish → declinable.** No stakes level was consulted; "prototype-ness" entered entirely through the four fields (no released user, trigger outside the promise).

**Case B — control, a genuinely material finding:** "`status --boot` drops the `stakes`/taxonomy field for `docs/dev`" against an entity whose value AC is "boot exposes the field." Four fields: released user = every FO at boot; harm = the FO does not receive it; affected value AC = **broken**; trigger = the normal boot path. → **Material → fix, not decline.**

**Determination:** the taxonomy **discriminates** using only committed fields + per-entity value ACs — declines A, fixes B — with no stakes input. The design stands under the parked stakes member; the discriminator is "does a value AC (or in-model boundary) break, via a supported trigger," which is computable per entity today. Had Case A come out Material, the design would need the parked stakes field and I would have escalated; it did not. The record home needs **no** further spike: it rides bw's already-proven `### Feedback Cycles` convention (bw's spike) and the existing gate body-read. This classification pass seeds the AC-1 fixture (Case A = the declined finding; Case B = the fixed control).

## Documentation changes (concrete before/after — ideation proposes, implementation applies)

**`skills/ensign/references/ensign-shared-core.md`** — insert a new section after `## Proving your work` (before `## Worktree Ownership`):

> ```markdown
> ## Consuming review findings
>
> Review findings are inputs to triage, not a fix list. A **review finding** is the output of a review instance the active stage's definition declares — the panels, audits, or validators it names, carried to you in the dispatch packet — or feedback the first officer explicitly routes to you as a review outcome. A direct instruction from the captain is not a finding: it is an order to follow or a decision to seek, never something to triage away. Before changing anything in response to such a finding, classify each against the workflow's declared finding-triage taxonomy (where one is declared) and this entity's own value acceptance criteria:
>
> - **Material** — breaks a value AC, or a declared non-negotiable boundary (safety, security, data-integrity, compatibility) reachable through the supported workflow. Fix it.
> - **Correct-but-disproportionate** (deferred risk or polish) — substantively right, but no value AC breaks and its trigger is outside the supported/promised workflow. Record a decline; do not fix it. The decline is your licensed disposition, not a dodge: name the finding, its class, and why it is not material (no value AC at risk; trigger outside the promise; the condition that would promote it to material).
> - **Needs decision** — a genuine product or compatibility fork. Escalate to the first officer; do not resolve it privately.
>
> Record the disposition — which findings were fixed as material, which were declined and why — in the entity's correction-round record (the workflow's `### Feedback Cycles` section) so the gate sees it. A finding you neither fix nor record is not triaged.
>
> **Narrowing an acceptance criterion to make a finding or rejection pass is not a licensed disposition.** Declining a disproportionate finding and narrowing the claim it targets are opposite moves under the same pressure: the first leaves the product unchanged and is yours to make; the second weakens the value the entity promised and is a design-reset event requiring the captain's sign-off, recorded so it is captain-visible — never a task-internal edit.
> ```

**`docs/dev/README.md`** — add one bullet in the `implementation` stage-def, after the design-reset bullet (`README.md:123`), composing with bw's in-stage-round bullet:

> `- When consuming an in-stage review round's findings (a roborev panel, a detached-audit pass, or routed gate feedback), triage before fixing. The committed finding-triage taxonomy is the release-scope classification in the `validation` stage-def (Material / Deferred risk / Polish / Needs decision) plus `.roborev.toml`'s four-field evidence record (released user + workflow; observable harm; affected value AC or non-negotiable boundary; trigger evidence). Fix only material findings; record a decline for a correct-but-disproportionate one — its class and why it is not material — in the `findings` field of the same `### Feedback Cycles` entry that logs the round; escalate a needs-decision finding to the FO. A value AC narrowed to make a finding pass is a design-reset event for the captain, not a fix (see the entry's AC-drift note).`

No other surfaces change (no FO-contract enforcement, no schema, no binary — see Out of scope).

## Acceptance criteria

**AC-1 (VALUE) — A seeded correct-but-disproportionate finding against a low-stakes fixture entity produces a recorded decline and a ZERO-line product diff, while a seeded material control against the same entity produces a fix (non-zero diff) — the taxonomy discriminates, it does not decline everything.**
Verified by: a live replay (the sprint DoD line) in which a dispatched ensign, given the consumption rule and the fixture entity, receives (a) the Case-A symlink-shape finding and (b) the Case-B value-AC-break finding, and its resulting worktree shows 0 changed product lines for (a) with a conforming decline recorded in `### Feedback Cycles`, and a non-zero fix for (b). Independent baseline that moved the wrong way: the archived incident-13 history, where the same finding-shape produced a non-zero dutiful fix. The number that can move wrong is product LOC changed on the decline path (0 vs the archived >0). This is the outcome the entity exists for; the live drive is validation's, per the live-drive proof rule.

**AC-2 — The decline record's class is falsifiable against the finding's own evidence: a decline recorded for a finding whose four fields establish materiality FAILS the check; a decline for a non-material finding passes.**
Verified by: a checked-in fixture (the Case-A declined finding + the Case-B material control, each with its four-field evidence as fixture data) and a small offline check that recomputes materiality from the four fields and asserts it matches the entry's `findings` class flag — so a mis-classified decline (material finding recorded as declined) fails it. The expected value is the fixture's four-field data (an independent source that can diverge from the flag), NOT a substring of any instruction file — satisfying the captain's prose-grep ruling. Seeds directly from the spike's Case A/B. Cost: low, no product code.

**AC-3 — The generic consumption rule + decline disposition + AC-narrowing clause live in `ensign-shared-core` (naming no dev artifact); the taxonomy citations, review instances, and `### Feedback Cycles` home live in `docs/dev/README.md` — neither layer leaks into the other.**
Verified by: a validation-stage landing-placement audit (like bw's AC-7) — read the two surfaces and assert the generic clauses are present in `ensign-shared-core` and carry no `.roborev.toml`/`docs/dev` specifics, and the dev specifics (the four-field taxonomy citation, the roborev/detached-audit instance) are present in `docs/dev/README.md` and absent from the skill. A one-off read at validation, output in the report, not a committed prose-grep test.

## Test plan

- **Value replay (live drive) → AC-1:** the DoD live replay — a dispatched ensign declines Case A (zero product LOC) and fixes Case B (non-zero), decline recorded in `### Feedback Cycles`. Validation owns the drive; the fixture entity + the two seeded findings come from the spike. Cost: one live dispatch, medium.
- **Offline classification check (fixture + check) → AC-2:** the spike formalized — a checked-in fixture with the two findings' four-field evidence and a small offline check asserting the recorded class matches the independently-recomputed materiality; fails a mis-classified decline. No product code. Cost: low.
- **Landing-placement audit (read at validation) → AC-3:** the core-vs-dev separation read. Cost: low.
- **No Go/binary tests, no product code in this cut** — the deliverable is prose (a skill section + a README bullet) plus a fixture and one offline check, matching the sibling bw's prose-first posture and the sprint's anti-over-engineering thesis.
- **High-stakes note:** `ensign-shared-core` is shipped contract/scaffolding (a high-stakes surface per the Proof policy), so the detached adversarial audit applies at validation before merge.

## Expected surface + tolerance (declared, per captain ruling)

- `skills/ensign/references/ensign-shared-core.md`: +1 section, ~15 lines (the generic rule + the by-reference trigger definition + captain-direction exclusion + decline disposition + AC-narrowing clause).
- `docs/dev/README.md`: +1 bullet in the `implementation` stage-def, ~4 lines.
- 1 fixture entity (+ its two seeded four-field findings) and 1 offline check for AC-2.
- **0 Go source files, 0 product LOC.** The `findings` field is an addition to bw's `### Feedback Cycles` entry, coordinated at the shared gate, not new machinery.
- **Tolerance: 2×**, with a hard self-check: any Go/product code, any new FO-contract enforcement check, or a second record section (a `### Declined Findings`) appearing in the first cut trips a reconfirm — deferring enforcement to its own captain-approved entity and keeping one record convention IS the point.

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
