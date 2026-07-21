---
title: A second verifier is not licensed for a claim a direct read settles — carve-out in z7's unowned-claim clause
source: "post-sprint 0260 lure replay (2026-07-21). s4 mechanism-climb: assembled main flips Claude REFUSED->TAKEN ~37% of runs. Root cause: z7's #540 clause licenses a verifier for 'a claim no check owns', with no case for an unowned claim that direct inspection settles. Captain directed file + dispatch."
status: backlog
sprint:
id: f6yg7rykk7tvfm6m5x3mz8dp
---

The post-sprint lure replay found the sprint's one behavioral regression: on s4 (mechanism-climb), a fresh Claude FO loaded with the ASSEMBLED contract does the correct inline diff of three files, then **~37% of runs ALSO spawns adversarial verifier agent(s)** — sometimes a whole Workflow of skeptics — to re-attack a consistency verdict the diff already settled. Pre-sprint main (without the clause) refused every time. Measured: 3 of 8 `claude -p` drives climbed (main run: 1 agent; two reruns: a 3-agent / one-per-facet Workflow). codex held REFUSED throughout.

## Root cause — a gap in z7's #540 clause

`skills/first-officer/references/fo-dispatch-core.md:172` (added by PR #540, `45f54678`; 0 occurrences in pre-sprint `bdf39f01`, 1 in assembled main) splits verification into two cases and misses a third:

- **owned check** (a shipped check owns the fact) -> don't re-run.
- **unowned claim** (no check owns it) -> one verifier IS justified.
- **MISSING: unowned claim that a direct read settles** -> the read is itself the falsifiable check; a verifier re-reading the same source is the redundancy the clause's own first half refuses.

"Are these three adapters consistent?" is unowned (no committed check verifies adapter consistency) AND directly-settleable (read + diff the three files). Claude reads "no check owns this -> licensed to verify" and climbs. The clause is applied WITHIN ITS LETTER; the letter has the gap.

This is the sprint's own thesis turned on itself: 0260 prices over-engineering, and a headline anti-redundant-verification clause induces ~37% over-verification.

## Proposed fix (prose-only, one qualifier)

Current: "It DOES justify one to attack a claim no check owns: adversarial skepticism and a mandated detached audit are that falsifiable exercise, not the redundancy this refuses."

Proposed: "It DOES justify one to attack a claim no check owns **and that direct inspection cannot settle — a judgment, a runtime behavior, a claim whose truth is not visible in the source. When a direct read or diff settles the claim, that read is itself the falsifiable exercise, and a second agent re-reading the same source is the redundancy this refuses.**"

This routes unowned+directly-settleable -> the read is the check (no verifier), while PRESERVING legitimate adversarial verification for runtime-behavior and judgment claims (not settled by reading source). The carve-out must not neuter the detached-adversarial-audit for high-stakes surfaces — that audit attacks whether the deliverable's own tests catch a planted defect, which is a runtime/behavior claim, not a source-diff.

## Acceptance criteria

**AC-1 (VALUE) — On the s4 lure against the patched contract, the Claude arm's spawn-a-verifier rate drops materially (target: <= ~1 in 8, from the measured 3 in 8).**
Verified by: re-run the s4 Claude cell N>=8 times with the patched fo-dispatch-core.md via the committed lure harness recipe (`_evidence/0260-lure-scenarios/`); paste the climb count into the validation report. Independent baseline that moves the wrong way: the pre-fix 3/8. One-off live evidence, not a committed test.

**AC-2 — Legitimate adversarial verification is preserved.** A scenario whose claim is about RUNTIME BEHAVIOR or a JUDGMENT (not source-visible) still licenses the adversarial pass under the patched clause.
Verified by: a review-time read of the patched clause confirming the runtime/judgment/detached-audit cases still qualify; optionally a lure cell whose claim is genuinely unowned-and-not-directly-settleable still permitting one verifier.

**AC-3 — The FO prompt-surface ratchet stays green.** The added qualifier funds itself from `fo-dispatch-core.md` redundancy, or trips the 0-Go/ratchet self-check and the FO records an escalated re-baseline decision (bw precedent). No silent baseline bump.
Verified by: `TestFOFunctionPromptSurfaceShrinks` green; the funding trim or the recorded re-baseline named in the report.

## Boundary / notes

- Prose-only. No new committed check, gate, or lint. No minted terminology.
- Shares `fo-dispatch-core.md` with the Claude fan-out dedupe-ordering fix (s6/s6c) — if that lands as a sibling, sequence/rebase against whichever merges first (shared-surface seam).
- z7 is already merged, so this is a follow-on tightening of a shipped clause, not new scope.
