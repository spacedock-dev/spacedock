---
id: hjb4z38k1f70yj1psb3yeaap
title: Prove or disprove the armed-merge parking failure with a probe that models context pressure
status: backlog
source: "0260 Commander drive, 2026-07-20. Member 85 (merge-guard-arm-not-a-stopping-point) was re-anchored and its value claim drew a NULL across two independent blind probes, so the anti-parking clause was parked rather than shipped. But the failure it targets was observed TWICE in production. The null is as likely a limit of the instrument as a verdict on the claim, and the honest resolution needs a probe the careful-reader design cannot provide."
started:
completed:
verdict:
score:
worktree:
issue:
---

Determine whether the armed-merge parking failure is real and contract-addressable, using a probe that can actually exhibit it.

## Problem

Two production incidents, both real FOs, both after the captain had already granted the push:

1. 2026-07-08 — the FO ran `merge guard --verdict passed` for three entities (which only ARMS), then read the hook file instead of opening the PR in the same turn, was pulled into an unrelated task, and left three arms untouched. Honest answer to "what did you do when I said push it": armed three, pushed nothing.
2. Session 6d175b2f — after the push was already granted, the FO re-asked the captain for permission to push, twice, instead of proceeding on the armed merge.

Member 85 proposed a contract clause naming an armed result as not a stopping point. Its value AC was probed twice, blind, and did not move: 3/3 before-text readers already chose proceed-this-turn at high confidence, citing text already in force. On that evidence, plus an inability to fund its 844 bytes against the FO prompt-surface ratchet, the paragraph was parked. A narrower payload that DID move its baseline (a `--no-ff` merge conflict is a blocker) shipped instead, at net +1 byte.

The unresolved question is whether the parking failure is real-but-undetectable-by-this-instrument. Both probes asked a fresh reader, with clean context, "what do you do next?" — and a careful reader answers correctly. Both production failures happened to an FO deep in a session, mid-interruption, after a long approved-gate sequence. A careful-reader probe is structurally incapable of exhibiting a context-pressure failure, so a null from it is weak evidence about the production case.

## Proposed approach

Ideation fills this in. The design constraint is the point: the probe must reproduce the CONDITIONS of the observed failures, not merely the question. Candidate directions, none yet chosen — a long-running drive where the arm is followed by an unrelated interrupt before the turn ends; a replay of the archived incident transcripts against current contract text; instrumenting real drives for armed-but-not-advanced states rather than asking readers anything at all.

Note the third direction may be the cheapest check that can fail: an armed-merge state that persists past a turn boundary is observable in workflow state, so the failure may be measurable in production directly, with no probe and no reader.

## Out of scope

Re-litigating 85's shipped `--no-ff` blocker, which is landed and separately evidenced. Re-adding the parked paragraph without new evidence — the parking decision stands until this task produces something that moves a baseline.

## Acceptance criteria

Ideation fills these in. Whatever design is chosen, the value AC must be capable of returning a null that is INFORMATIVE — i.e. a null must distinguish "the failure is not real" from "the instrument cannot see it," which is precisely what the 0260 probes could not do.
