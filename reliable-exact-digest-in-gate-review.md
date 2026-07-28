---
title: "Make exact Briefing digest presentation reliable in live gate reviews"
status: backlog
source: "se0 complete Sonnet live suite, 2026-07-28: the FO completed the recorded-gate lifecycle but presented only 'digest per briefing.json artifact set' instead of the bound sha256 value; the exact semantic oracle correctly rejected it."
score: 0.85
sprint: durable-decisions
group: gate
sprint-readiness: ready
issue:
id: w5bfnrvpcphw857nzz93340c
---

## Problem

A supported First Officer journey can approve and consume a recorded gate without presenting the exact bound Briefing digest to the captain. In the retained complete Sonnet run, the decision line was actionable, the exact Briefing ID was present, and the durable lifecycle and successor dispatch completed, but the review replaced the required digest with the phrase `(bound, digest per briefing.json artifact set, gate state open)`.

This is a genuine semantic conduct failure, not an oracle failure. The `present-gate` contract requires the reviewed snapshot's exact Briefing ID and `sha256:` digest so the captain knows which immutable package the decision spends. The failure occurred after the same focused Sonnet scenario had passed, so the remedy must make the behavior reliable rather than merely adding another example phrase.

The Sonnet execution of the shared recorded-gate lifecycle is temporarily TODO under this task. Keep Opus, Codex, Pi, rejection, keep-moving, and all deterministic gate lifecycle coverage active. Re-enable the Sonnet case when this task lands.

## Acceptance criteria

**AC-1 (VALUE) - Every captain-facing recorded-gate review identifies the exact immutable Briefing being decided.**
Verified by: repeated clean Sonnet recorded-gate journeys each display the canonical Briefing ID and exact `sha256:` digest after the retained Briefing commit and before decision mutation.

**AC-2 - The reliability fix preserves agent ergonomics and the provider-neutral gate lifecycle.**
Verified by: the FO derives the digest from the bound gate package without caller-authored reconstruction, prompt-specific hardcoding, or a transport dependency.

**AC-3 - The oracle remains semantic and adversarial.**
Verified by: exact ID+digest positives pass; ID-only, vague digest references, wrong digests, tool output, and post-decision summaries remain red.

**AC-4 - The quarantined Sonnet live journey is restored.**
Verified by: remove the linked TODO, then run the focused Sonnet recorded-gate case repeatedly and the complete affected Claude live suite at the exact candidate tip.

## Boundary

Do not weaken the exact-digest requirement and do not teach one fixture phrase. Determine why the loaded gate presentation path allowed the FO to omit a value already present in the canonical retained package, then fix the smallest authoritative contract, skill wiring, or presentation assembly seam. V1 is unreleased; add no compatibility behavior.
