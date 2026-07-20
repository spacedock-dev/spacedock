---
id: z7sfm93ccddg7x2tycp1smwy
title: Falsifiability ladder replaces "code gate over prose rule" — with infra consent, fan-out checkpoint, and bare-ordinal itemization
status: backlog
source: "0260 shaping — agent-derail forensics audit, 2026-07-19."
score: "0.75"
sprint: 0260-proportionality
group: ladder
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
