---
id: z7sfm93ccddg7x2tycp1smwy
title: Falsifiability ladder replaces "code gate over prose rule" — with infra consent, fan-out checkpoint, and bare-ordinal itemization
status: ideation
source: "0260 shaping — agent-derail forensics audit, 2026-07-19."
score: "0.75"
sprint: 0260-proportionality
group: ladder
started: 2026-07-20T03:29:38Z
gates:
  version: 1
  current:
    gate: gate:docs-dev:z7:ideation
    attempt: gate-attempt:z7-ideation-1
  records:
    - id: gate:docs-dev:z7:ideation
      stage: ideation
      current-attempt: gate-attempt:z7-ideation-1
      attempts:
        - id: gate-attempt:z7-ideation-1
          sequence: 1
          state: open
          briefing:
            id: briefing:z7-ideation-1
            digest: sha256:f8a918e9f5315d352b700b5c0a79e9d9a34087cbc311073ec3631968776cc937
---

"Prefer a code gate over a prose-only rule" is a standing instruction to convert any guarantee into enforcement code — unscoped by stakes, it produced presence tests and unasked CI/lint infra, and it degrades worse in non-dev workflows where every check is new infra. Replace it with the falsifiability ladder: shipped system guards → existing mechanical checks → falsifiable exercise (replay, source-check, adversarial skeptic) → captain judgment → build new machinery (last, consent-gated). Same edit carries: new enforcement surfaces are not "obvious reversible work" (consent required); a fan-out checkpoint before an investigation's Nth spawned entity/PR; identifier minting reserved to the system, ad-hoc itemization uses bare ordinals. Grouped with 1p9, cy, 85.

## Problem

{Ideation fills in. Evidence: 11-phrase contract-presence test (bef9653f:496-509) as the clause's cheapest compliant response; flake→4 open PRs (ab6c437e); ProfileLeaseV1 minted for a throwaway smoke test; "I dressed it up to make it sound principled" (bef9653f:507). Subject to the 0250 leanness constraint: net contract bytes measured, lazy-loaded over boot-resident.}

Live propagation example for the identifier clause (0260 shaping session, 2026-07-20): the FO's dispatch prompt for a personal-instructions audit itemized the audit dimensions as "A. … B. … C. … D. …"; the worker faithfully propagated those letters into its report ("Class: A + C" in the summary table, "## A." section headers), which then required a legend to read — caught by the captain, fixed by a revision round. The clause must therefore bind DISPATCH PROMPTS, not just captain-facing output: a scheme minted in a prompt becomes the vocabulary of every downstream artifact, and policing outputs alone catches it one propagation too late. This occurred in the same session that codified the rule, by the FO that codified it — minting is a default behavior under itemization pressure, not an occasional lapse.

Second and third live examples, same session, same FO, AFTER the first was caught and recorded above — evidence that a one-time correction does not extinguish the behavior:

1. Briefing-id shorthand in captain-facing prose. The gate-notation contract legitimately binds briefing ids in entity frontmatter and the review room's probes file; the FO then used the bare suffixes ("1a/2a", "the comparator", "delta re-presentation") as conversational shorthand to the captain, who could not read them and said so — and the FO repeated the shorthand in its next worker message anyway. Rule refinement this argues for: system ids are legal only where the notation binds them; captain-facing prose restates in plain words, every time, with no exemption for ids that are individually sanctioned.
2. Sanctioned-vocabulary overloading. "-cycleN" is system vocabulary denoting feedback-rework cohorts after a gate rejection. The FO applied "cycle 2"/"cycle 3" to in-gate revisions of a HELD gate (no rejection ever recorded, no Feedback Cycles entry exists), and the worker durably propagated it into the entity body as a stage-report heading. Overloading an existing identifier onto a thing it does not denote is minting with extra camouflage: it reads as sanctioned while corrupting the term's meaning — arguably worse than a fresh coinage because it degrades the system vocabulary itself.

Enforcement consequence for ideation: the clause needs its weight on the FO's OWN prompts and prose (the propagation source), not on downstream artifact review; and per this sprint's ladder, the check is judgment-rung (captain catch, gate review), not new lint machinery — three same-session recurrences by the rule's own author are the calibration datum for how sticky the behavior is.

Fourth live example, same session — minting is not only identifiers, it is ABSTRACTIONS: defending a sprint member against the essence test, the FO coined "essence of reach" — a compound phrase gluing a real claim ("a rule exists but the agents who need it never see it") to a keep-bucket label, making a choice sound like a derivation. The circle: the FO shaped an essence category partly so the member it wanted to keep would fit, then cited the category as the reason the member stayed. Applying the bare kill-test question ("which observed failure does this kill, and is the killer already in place?") without the coined phrase produced the opposite — correct — answer, and the member was parked. Same family as the identifier examples: language minted at the exact moment structure is lacking, to manufacture its appearance. Prior art in the corpus: "prose-grep over a runtime-loaded contract" (bef9653f:507, its author's admission: "I dressed it up to make it sound principled"). The clause the ideation shapes should name the general form: coined compound abstractions in an argument are a smell that the argument is doing the work the evidence should.

## Folded scope: smallest-sufficient-mechanism sharpening (absorbed from smallest-sufficient-mechanism-deterministic-fact-gate, 0260 re-lock 2026-07-20)

Same contract region, one edit: the four tightenings from the absorbed entity ride this
edit — (1) discriminate judgment calls from deterministic facts in the "independent
adversarial verification" justification; (2) the "N agents != N confidence" corollary;
(3) a session thoroughness directive (Ultracode) raises the bar on the answer, never the
weight of the mechanism; (4) sized for the lazy-loaded reference file per the leanness
constraint. Also per the re-lock: the fan-out checkpoint in this entity's scope REDUCES to
a prose clause (surface before the Nth spawned entity/PR of one investigation) — no counter
binary; a counter is speculative machinery against this sprint's own thesis.

## Proposed approach

One edit to the FO contract, split so the reframe of the core principle stays boot-resident (it shapes every session's default) and every elaboration lazy-loads. Target files: `skills/first-officer/references/first-officer-shared-core.md` (boot-resident, the most-read region) and `skills/first-officer/references/fo-dispatch-core.md` (deferred, read before first dispatch — the natural home for triggers that fire at a dispatch/spawn/prompt-authoring decision). No new binary, flag, counter, or lint: the fan-out checkpoint and the minting clause are judgment-rung prose, per the re-lock and the four calibration examples.

### Boot-resident wording (first-officer-shared-core.md)

1. Replace the retired clause at the "Working Principles" head (currently `**Prefer a code gate over a prose-only rule.**`, 669 bytes) with the ladder. The reframe preserves the retired clause's load-bearing payload — "wording-present is not behavior" and "a prose-only rule is not AC satisfaction on its own" — and de-escalates only the DEFAULT: building machinery drops from the preferred move to the last rung. Rung 3 carries "replay the behavior at the seam where the failure lives" — a compact clause, not a paragraph — from the 0.25.1 AC-narrowing incident (synthesis.md:63): the release proved the adapter emitted no inherited turns but never proved the live FO invocation, and the failure reproduced at the unproven seam (the real call site), so a falsifiable exercise must run at the seam where the failure lives, not the nearest convenient layer.

BEFORE:

```text
**Prefer a code gate over a prose-only rule.** When a guarantee can be enforced by the binary or a failing test (a `status` guard, a test that fails on violation), prefer that. A prose-only rule's ceiling is "the wording is present"; wording-present is not behavior. A prose-only rule must not count as AC satisfaction on its own: if the guarantee matters, the real assurance is a code-level gate underneath, and the prose points at it. An AC of the form "the contract says X" is satisfied only by "the binary or a test enforces X, and here is the run that proves it." The gate's AC cross-check refuses a criterion whose only proof is review of the entity's own prose.
```

AFTER:

```text
**Climb the falsifiability ladder; building new machinery is its top rung, not the default.** A guarantee earns trust by being able to FAIL — wording-present is not behavior. Reach for the cheapest rung that can falsify the claim: (1) a guard the system already ships; (2) an existing mechanical check; (3) a falsifiable exercise — replay the behavior at the seam where the failure lives, check a claim against its source, or let an adversarial skeptic try to break it; (4) captain judgment. Only when none of these can falsify it do you (5) build new machinery — the last rung, and a new enforcement surface is not obvious reversible work, so it is consent-gated (`fo-dispatch-core.md`). A prose-only rule still does not satisfy an AC on its own: "the contract says X" is met only by a rung-1–3 run that could have failed. The gate's AC cross-check refuses a criterion whose only proof is review of the entity's own prose.
```

2. Carve the consent stop into the "Do obvious reversible work without ceremony" posture bullet (currently 177 bytes) — the forensics named this exact line a clause-level driver ("CI/harness addition reserved to no one", synthesis.md:39), so the carve-out must sit AT the driver, not only in the deferred detail.

BEFORE:

```text
**Do obvious reversible work without ceremony** — reversible steps the contract allows just happen; reserve asking for choices that are hard to reverse or genuinely matter.
```

AFTER:

```text
**Do obvious reversible work without ceremony** — reversible steps the contract allows just happen; reserve asking for choices that are hard to reverse or genuinely matter. But standing up a NEW enforcement surface (check, harness, CI job, lint) is never in this class: it is the ladder's consent-gated top rung, not ceremony-free work (`fo-dispatch-core.md`).
```

3. Add one FO-posture bullet for the bare-ordinal / no-minting rule. It stays boot-resident because the calibration datum is that minting is a DEFAULT under authoring pressure (four same-session recurrences by the rule's own author): the rule must be visible where prompts and captain-facing prose are WRITTEN, not only where artifacts are reviewed. The compact bullet ships alone; its three-forms elaboration is deferred (gate bounce, cycle 2) and revives only if live drives show the compact rule insufficient.

NEW (insert as a `**FO posture:**` bullet):

```text
**Author in the system's vocabulary; don't mint your own.** Ad-hoc itemization in your own prompts and captain-facing prose uses bare ordinals — identifier minting is reserved to the system. This binds what you WRITE, not just what you review: a scheme minted in a dispatch prompt becomes every downstream artifact's vocabulary.
```

### Deferred wording (fo-dispatch-core.md, at the dispatch-decision region near `«dispatch.next-action»`)

All three blocks are NEW, added to the module already read before first dispatch — so each is present exactly when its trigger (an infra-build dispatch, the Nth spawn, or a climb decision) is evaluated.

```text
**Consent stop before building a new enforcement surface.** A ready action whose deliverable is a NEW check, test harness, CI job, or lint (not running one that already exists) is the falsifiability ladder's top rung: do not dispatch it as obvious reversible work. Surface it to the captain first and dispatch only on a declared yes — the license hangs off the captain wanting it, never an inference that it would help. This holds hardest in non-dev workflows, where every mechanical check is new infra.

**Fan-out checkpoint.** Before spawning the Nth entity or opening the Nth PR of a SINGLE investigation (a flake chase, a review-rework, a refactor sweep), stop and surface the running count and the cap question to the captain rather than spawning again. Keep-moving speeds independent, already-scoped work; it does not authorize an unbounded spawn chain off one thread. This is a prose checkpoint on the FO's judgment — there is no counter binary.

**Climbing for adversarial verification — judgment calls only.** "Independent adversarial verification" justifies a second agent only for a JUDGMENT CALL; a DETERMINISTIC FACT is settled by running the check that owns it, not a second opinion. And N agents reaching the same answer is one confirmation observed N times, not N independent confidences — agreement among spawned workers raises cost, not confidence. (Session thoroughness — "Ultracode is on" — already raises the answer's bar, never the mechanism's weight; see the smallest-sufficient clause.)
```

Deferred at the gate (cycle 2): the three-forms elaboration block that once accompanied the minting bullet is removed from this edit, narration-first — the compact boot bullet ships first and the elaboration revives only if live drives show it insufficient. Precedent: the smallest-sufficient rule itself once grew an explanatory layer and the durable correction was collapsing it to one compact rule (synthesis.md:62, codex:019f499a @1524). The four recorded live examples in the Problem section (the "A./B./C./D." prompt itemization, the briefing-id shorthand, the `-cycleN` overloading, and the "essence of reach" coinage) remain the evidence base for the elaboration when it revives — they are the three forms it would enumerate.

### Byte accounting (riskiest-path spike — done, measured)

The riskiest ideation-stage path was proving the wording FITS a lean budget with the bulk deferred. Measured against the current file:

| Addition | Placement | Bytes |
|---|---|---|
| Falsifiability ladder (replaces 669-byte clause; incl. rung-3 seam clause) | boot-resident | +263 |
| Consent carve-out on "obvious reversible work" (was 177) | boot-resident | +186 |
| Bare-ordinal / no-minting posture bullet (self-contained) | boot-resident | +331 |
| Consent stop + fan-out + adversarial-verification sharpening | deferred (fo-dispatch-core.md) | +1523 |
| **Boot-resident subtotal** | | **+780** |
| **Deferred subtotal** | | **+1523** |
| **Total contract delta** | | **+2303** |

66% of the addition lazy-loads; boot-resident growth is confined to the ladder (a replacement, not a pure add), a carve-out co-located with its named driver, and a terse authoring rule whose visibility the calibration datum requires. The cycle-2 deferral of the three-forms block dropped the deferred subtotal by 776 bytes; dropping the now-dangling pointer to it trimmed the boot bullet from 383 to 331, so boot-resident landed at +780 while gaining the rung-3 seam clause. No format round-trip, runtime handoff, or new flag is involved — the mechanism is prose the FO reads — so the behavioral claims (the consent stop is obeyed; a minting occasion is declined) are validation-owned live drives, not an ideation spike; fixtures named below.

### Generic vs dev-specific (required axis)

Generic contract prose (all boot + deferred wording above): the ladder rungs, the consent stop, the fan-out checkpoint, and the minting clause name no dev-only object as load-bearing. Rung 3 reads "check a claim against its source" (a citation source-check in a research workflow, a fact-check in a writing workflow) as naturally as "replay the behavior at the seam where the failure lives" (a test at the real call site in dev; a re-derivation from primary sources, not a summary, in research); rung 5 "new machinery" is any new check/harness/validator, not specifically CI. The forensics confirm the generic frame: "in a research workflow the same disease yields citation-padding; the clause degrades WORSE outside dev because every mechanical check is new infra" (synthesis.md:37).

Dev-template / repo realizations (NOT in this entity's diff): the concrete enforcement examples (a `status` guard, a failing Go test), the dev-workflow README/template summaries that name the principle, and any repo check hook are the `template` group's and the driving session's job. The minting clause gets NO lint machinery — it is judgment-rung, because a lettered list in a fixed reference table is fine and only minting-under-argument-pressure is the harm; a prose-grep cannot discriminate the two and would itself be the fabricated-rigor pattern this sprint retires.

### Mechanism justification (value AC served / simplest alternative / why insufficient)

1. Falsifiability ladder — serves AC-1, AC-3. Alt: just delete the retired clause and lean on existing proof-discipline. Insufficient: deletion drops the "wording-present is not behavior / prose-only is not AC satisfaction" payload that legitimately catches fabricated-prose ACs; the ladder keeps it while de-escalating the build-machinery default. Alt: a stakes threshold ("gate only high-stakes"). Insufficient: stakes live in a parked member; the ladder is stakes-agnostic and needs no ontology.
2. Consent stop (deferred) + boot carve-out — serves AC-1. Alt: mention "consent-gated" only inside rung 5. Insufficient: the observed failure (incident 1) is that infra-building gets DISPATCHED as obvious reversible work; the stop must fire at the dispatch decision and the named driver line must be carved, or the driver keeps firing.
3. Fan-out checkpoint (prose) — serves AC-5. Alt: a counter binary that blocks at N. Insufficient / over-reach: that is new enforcement infra to enforce "don't over-build infra" — speculative machinery against the sprint's own thesis (re-lock decision); a prose checkpoint on FO judgment is the proportionate rung.
4. Minting clause (compact boot rule; three-forms elaboration deferred to a later edit) — serves AC-4. Alt: a contractlint prose-grep for minted labels. Insufficient: that is the fabricated-rigor pattern being retired, and minting is a judgment call a grep cannot discriminate. Alt: police downstream worker artifacts. Insufficient: catches it one propagation too late — a prompt-minted scheme becomes every downstream artifact's vocabulary, so the clause must bind the authoring source. Alt considered and REJECTED at the gate: ship the three-forms enumeration now. Insufficient reason to ship: it over-bundles the edit against the collapse precedent (synthesis.md:62) — the compact rule ships first, the elaboration only if drives show it needed.

### Downstream propagation (doc-diff)

In scope (this workflow's own README): `docs/dev/README.md:74` names "prefer a code gate over a prose-only rule" — rename to "climb the falsifiability ladder (new machinery is the top rung, not the default)". Delegated to the `template` group (template layer, per the sprint layer map): `skills/commission/references/templates/development.md:88` and `experiment.md:120` carry the same phrase and rename identically when `template` lands.

## Acceptance criteria

AC-1 (value — the reason the entity exists). Under the implemented contract, an FO handed the process-control-harness ideation brief (incident 1 shape: a PTY/process-control harness mandated for a disposable smoke test) HALTS at the consent stop and surfaces the build to the captain before dispatching it; under `main` the same brief dispatches the harness build with no stop. Test: two live FO drives (branch vs `main`) over the same seeded brief; observe stop-and-surface vs dispatch. Baseline that can move the wrong way: `main` is the negative control — if the branch also dispatches, AC-1 fails. Fixture: `_evidence/0260-agent-derail-forensics/` incident 1 (zaphod PTY session).

AC-2 (value — leanness). The boot-resident file (`first-officer-shared-core.md`) net byte delta vs `main` is within budget (target ≤ +800; measured draft +780), the ladder occupies the retired clause's slot, and every non-boot addition (consent stop, fan-out checkpoint, adversarial-verification sharpening) lands in `fo-dispatch-core.md`, not boot-resident. Test: `wc -c` / `git diff --stat` of the boot file branch vs `main`, plus a placement audit asserting each named block appears in `fo-dispatch-core.md` and NOT in the boot file. Baseline that can move the wrong way: boot bytes can balloon past budget.

AC-3 (mechanism — ladder is generic). The ladder states the five rungs with "build new machinery" last and consent-gated, preserves "wording-present is not behavior" and "prose-only is not AC satisfaction", and reads cleanly for a non-dev workflow — no dev-only term is load-bearing in the rungs. Test: AC-1's drive exercises the top rung; the ideation-gate reviewer's generic read-through confirms rung 3 maps to a citation source-check and rung 5 to any new validator. Paired with AC-1.

AC-4 (mechanism — minting clause). The compact boot bullet binds the FO's OWN prompts and captain-facing prose (not downstream review) and reserves identifier minting to the system; it is judgment-rung (no lint). The three-forms enumeration is NOT part of this edit — deferred to revive if the compact rule proves insufficient in live drives; the four recorded examples in the Problem section are its evidence base. Test (fixture unchanged by the deferral): a live-drive replay of the recorded "essence of reach" occasion (body example 4) — under the branch the FO applies the bare kill-test question and parks the member; under `main` it coins the abstraction and keeps it. Deterministic corroborator: a source-check that this entity and the sprint index carry no minted scheme (bare ordinals only). Baseline that can move the wrong way: `main` mints and keeps.

AC-5 (mechanism — fan-out checkpoint + folded sharpenings). A prose fan-out checkpoint surfaces before the Nth spawned entity/PR of one investigation with no counter binary; the two folded smallest-sufficient sharpenings (adversarial verification justifies a second agent for judgment calls only, not deterministic facts; N agents are one confirmation observed N times, not N confidences) sit at the climb/spawn decision point in `fo-dispatch-core.md`; the Ultracode/thoroughness sharpening is recorded as already covered by the existing smallest-sufficient clause (no new bytes). Test: a wording/placement audit (the clauses are present in `fo-dispatch-core.md`; the implementation diff touches only `.md` files — no new CLI/flag/counter — confirming "no counter binary"); optional fan-out replay against the ab6c437e flake→4-PR incident showing the checkpoint surfaces by the 2nd/3rd PR.

## Test plan

What verifies. Behavioral claims are proven by live FO drives over the two contract versions (branch vs `main`) using the `_evidence/` incident records as fixtures — never a prose-grep over the contract this change writes (0260 live-drive rule). AC-1 and AC-4 are the two behavioral drives; AC-2 and AC-5 are cheap on-disk/diff checks; AC-3 is a reviewer read-through plus AC-1's drive.

Cost / complexity. AC-2 (byte + placement audit) and AC-5 (wording/placement audit): cheap, seconds, `wc -c` + `grep` + a read of the implementation diff. AC-1: medium — two live FO drives seeded with the harness brief; the expensive-but-decisive proof, owned by validation. AC-4: medium and softer — an authoring-tendency replay is noisier than the consent stop, so it carries the deterministic source-check corroborator; the "essence of reach" occasion is the sharpest fixture because its baseline outcome (coin-and-keep vs bare-question-and-park) is recorded and observable. No fixture Go test or CLI golden is added — there is no new command surface, and the minting clause is judgment-rung by design.

Spike / riskiest path. Done in ideation: the concrete before/after wording is drafted and byte-measured (boot +780, deferred +1523, 66% lazy-loaded after the cycle-2 three-forms deferral) — the design fits a lean budget with the bulk deferred. No format round-trip, runtime handoff, or tool-flag support is in play (the mechanism is prose the FO reads at boot / at first dispatch), so "no round-trip spike needed"; the behavioral mechanisms (a prose consent stop is actually obeyed; a minting occasion is actually declined) are validation-owned live drives per AC-1 and AC-4, with fixtures named in `_evidence/`.

## Stage Report: ideation

- DONE: Riskiest path first — ladder replacement drafted as concrete before/after contract wording under the leanness constraint; net byte delta measured, lazy-loaded placement named for every addition
  Before/after in Proposed approach; measured boot-resident +796, deferred +2299 (74% lazy-loaded); each addition's file/trigger named in the byte-accounting table.
- DONE: GENERIC vs DEV-SPECIFIC explicit — ladder, consent, fan-out-checkpoint, minting clauses are generic contract prose; enforcement examples and check hooks are dev-template/repo realizations; no lint machinery for the minting clause (judgment-rung)
  "Generic vs dev-specific" section: rung 3 = source-check (research) / replay (dev); dev examples + README/template summaries delegated to `template`; minting is judgment-rung per the four calibration examples.
- DONE: The minting clause binds the FO's OWN prompts and prose and covers all three observed forms (fresh identifier schemes, sanctioned-vocabulary overloading, coined compound abstractions)
  Boot posture bullet + deferred three-forms block; bound to authoring not review per the calibration datum; AC-4 replays the "essence of reach" occasion.

### Summary

Filled the ideation deliverable: Proposed approach with concrete before/after wording for the boot-resident ladder (replacing the 669-byte code-gate clause), the consent carve-out on the "obvious reversible work" driver line, the bare-ordinal/no-minting posture rule, and the deferred block (consent stop, fan-out checkpoint, minting three-forms, adversarial-verification sharpening) in `fo-dispatch-core.md`. Riskiest path exercised: wording byte-measured (+796 boot / +2299 deferred, 74% lazy-loaded). Five ACs — AC-1 (consent-stop replay) and AC-2 (leanness byte/placement) measure end-values against baselines that can move the wrong way; AC-3/4/5 are mechanism ACs paired to them. Key decisions: fan-out and minting are judgment-rung prose (no counter binary, no lint) per the re-lock and calibration datum; behavioral claims are validation-owned live drives, not an ideation spike, since the mechanism is prose the FO reads. Downstream: `docs/dev/README.md:74` rename in scope; the two template summaries delegated to the `template` group.

## Stage Report: ideation (cycle 2)

- DONE: Gate bounce — defer the three-forms elaboration block; keep the compact boot bullet
  Three-forms block removed from the deferred wording (776 bytes off deferred); boot mint bullet kept and made self-contained (dropped the now-dangling pointer, 383→331); revival note added citing the collapse precedent (synthesis.md:62) and the four examples as the elaboration's evidence base.
- DONE: Byte accounting updated
  Boot +780 (was +796), deferred +1523 (was +2299), total +2303, 66% lazy-loaded; table and Test-plan spike line updated; cycle-1 report left as its historical record.
- DONE: AC-4 reduced to the compact clause; replay fixture retained
  AC-4 now claims only the compact bullet (binds prompts/prose, reserves minting to the system, judgment-rung); three-forms enumeration marked deferred; the "essence of reach" replay fixture unchanged.
- DONE: AC-2 / AC-5 placement audits updated
  AC-2 deferred list drops "minting three-forms" and cites +780; AC-5 (fan-out + folded sharpenings) unaffected, stands.
- DONE: rung-3 strengthened from the 0.25.1 incident, as a clause not a paragraph
  Ladder rung 3 now reads "replay the behavior at the seam where the failure lives" (+36 bytes on the ladder); provenance in approach item 1 and generic-vs-dev (synthesis.md:63).

### Summary

Applied the gate's single ask: the three-forms minting elaboration is deferred to a later edit (revives only if live drives show the compact rule insufficient), mirroring the collapse precedent the reviewer cited; the compact boot authoring bullet stays. Everything else the gate approved-shaped is unchanged (ladder replacement, consent carve-out, consent stop, fan-out checkpoint, judgment-vs-fact sharpening). Byte accounting re-measured (boot +780, deferred +1523, 66% deferred). Added the 0.25.1 AC-narrowing incident as rung-3's "at the seam where the failure lives" strengthening — a compact clause, not a paragraph — with the four recorded examples preserved as the deferred elaboration's evidence base for revival.
