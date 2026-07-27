---
id: t25b6hdevexavss2j3thjqy0
title: "Gate legibility — what is the First Officer's summary, and what must be excerpted from the stage result the captain judges"
status: backlog
source: "Captain question 2026-07-27 on the spacedock-subspace present-gate proposal, retained as proposal-subspace-fo.md beside this entity. Filed as the broader question the proposal is one symptom of; the proposal's own Journey-section shape is subspace-specific and is input, not the answer."
started:
completed:
verdict:
score: 0.75
worktree:
issue:
---

Make a gate presentation legible and meaningful: separate what the First Officer summarises from what must reach the captain as the worker's own words, and let each workflow declare which is which.

## The question

A gate message today is entirely First-Officer-mediated. Assembly rule 3 mandates it explicitly: the checklist renders "one bullet per DONE/SKIPPED/FAILED item as a verb-noun gist (≤10 words, FO paraphrase, no new facts)". So every fact the captain votes on arrives as the FO's summary of the worker's claim, never the worker's own text.

Two different things are being conflated, and naming them is the point of this task:

1. **The gate summary** — the FO's own framing: the recommendation, why, what is at risk, the decision ask. This is irreducibly judgment and belongs to the FO.
2. **The excerpted stage result** — the content the captain must actually judge, produced by the worker, surfaced as written rather than paraphrased.

The open question is where the line falls, how a workflow declares it, and what stays universal.

## Why this is desirable — evidence, not preference

**The template has no slot for content, and the length cap enforces that.** Rule 1 makes the spine title / chosen-direction / recommend / decision; rule 9 caps the message at 15-25 lines. For a mechanical stage that is correct — the checklist *is* the content. For a stage whose entire output is a design, the captain is asked to vote on a one-line `Chosen direction:` and a verdict. An FO obeying the skill has no room for substance and no slot to put it in. The template and the requirement conflict, and the template wins.

**A First Officer following the skill faithfully produced a rejected briefing, twice.** From the spacedock-subspace `one-question-in-a-review` ideation gate, the captain's verbatim rejection: *"I ASKED FOR A JOURNEY DESCRIPTION AND ESTIMATED SCOPE OF CHANGES / WHERE THIS WILL BE TOUCHING. WHY ARE YOU ARGUING WITH ME ABOUT THE APPROACH WHEN YOU HAVE NOT SURFACED THE JOURNEY"*. That is not a discipline failure better FO judgment fixes.

**Independently reproduced in this repo.** At every substantial gate in the 2026-07-27 session the FO appended prose *after* the template block, because the template alone could not carry the decision. The workaround was universal, which makes it a property of the template rather than of any one FO.

**The float channel outperforms chat precisely by bypassing paraphrase.** On this repo's `cn` ideation gate the captain annotated exact quoted text from the worker's design — `Runtime: claude (CLAUDECODE)` and an outside sandbox string. One annotation killed three proposed output strings and led the worker to identify that the proposed fix reproduced the original defect's own category error, answering a launch question in a line whose purpose was to report the session. That finding was anchored to the worker's exact bytes and could not have surfaced from a ≤10-word FO gist. Subspace measured the same effect at larger scale: three floats, six inline annotations, which changed the entry point from a new script to a flag on an existing entry, inverted the process topology, and cut roughly 160 lines of over-design before any of it was written.

**The excerpt is what the captain reaches for.** In both repos the decisive annotations landed on worker-authored text. Nothing decisive landed on the FO's summary.

## Why "add a Journey section" is not the general answer

The attached proposal asks for a required `Journey:` and `Where it touches:` block. That is the right content for a workflow building a TUI product, and the wrong universal: a research workflow's gate wants question, method and confidence; a content workflow's wants a draft excerpt, audience and what changed since the last round; an ops workflow's wants blast radius and rollback. Spacedock's product is commissioned workflows, so a gate shape hardcoded to one domain's vocabulary imposes that domain on every other.

## The structural defect this sits on

The dispatch side of the loop is already workflow-declarable and the gate side is not. `internal/status/stages.go:110` parses per-stage `context-sections` with `stages.defaults` and per-state overrides, strictly validated — a workflow already declares which sections a *worker* receives. There is no equivalent for what a *gate* must carry. That asymmetry is the structural defect; "no journey section" is one workflow's symptom of it. Any mechanism proposed here should reuse that precedent rather than invent a second declaration surface.

## Hypotheses for ideation to confirm or refute, not decisions

Probably universal, and belonging in the shared skill:

- The spine — title, chosen direction, recommend, decision. It concerns *deciding*, which every workflow does.
- An OBSERVED / DESIGNED marking distinguishing what was exercised from what was only designed. This is exercised-versus-asserted and applies to any claim in any domain; it also promotes the riskiest-mechanism question from a claim buried in a spike section to something a captain can scan.
- The length cap applying to FO prose only. Declared or excerpted content is content, not narration.

Probably workflow-declared:

- Which sections are excerpted, and from where.
- Domain precision rules such as the proposal's "programs, not roles", which is a software-domain rule.
- Whether failure modes must be rendered in the same terms as the success path — general in spirit, journey-shaped as written.

## Open questions

- Does the excerpt come from named entity-body sections, or from something the worker explicitly marks captain-facing? Section names reuse the `context-sections` precedent; worker marking puts the selection where the knowledge is.
- Verbatim or bounded? An unbounded excerpt makes the gate unreadable; a bound reintroduces the paraphrase this task exists to remove.
- May the FO editorialise an excerpt at all, or is the summary the only place judgment appears?
- If the float shows the captain the whole artifact anyway, is the chat gate a degraded mode, and should the two be specified together rather than separately?
- Does FO process accounting belong in a gate message at all? The proposal argues for banning it; a narrower rule may be right, since an FO error that changes what the captain is voting on is decision-relevant.

## Out of scope

- The gate-room ordering trap and its four sub-frictions (proposal section 4). That is a real defect in shared machinery, independently reproduced in both repos, but it is a recorder bug rather than a legibility question. It belongs with `sk` (gate-agent-ergonomics).
- Changing the float provider, the recorder's write surface, or the decision affordance.

## Acceptance criteria

Ideation fills these in. The end state is a gate presentation whose FO-authored and worker-excerpted parts are distinguishable by rule rather than by FO discretion, with the excerpt selection declared by the workflow, and at least one criterion measuring legibility against something that can move the wrong way rather than asserting the rule shipped.

## Test plan

Ideation fills this in.
