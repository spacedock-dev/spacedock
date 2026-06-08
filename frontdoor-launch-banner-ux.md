---
id: cmxchb8y1y2m455xhx7ce87g
title: Launch-banner UX — first-officer framing, status-command overload, multi-workflow limbo
status: backlog
source: "FO + captain launch-banner review (2026-06-08), following yq (frontdoor-launch-ux). The 3-line pre-launch banner (frontdoor.go:139-141) reads confusingly on the new-user / multi-workflow path."
started:
completed:
verdict:
score:
worktree:
issue:
sprint:
group:
sprint-readiness:
---

A followup on yq (frontdoor-launch-ux): the pre-launch banner (`launchBanner`, `internal/cli/frontdoor.go:139-141`) is clean on the common single-workflow path but wobbles on the new-user / multi-workflow path. Three issues worth fixing, worst first.

## Problem

The banner emits three lines:
1. `spacedock {Version} · first officer launching {host}`
2. `Workflow: {detectedWorkflow}`  (e.g. `5 workflows detected (run spacedock status to pick)`)
3. `{host} is starting as your first officer; run spacedock status inside the session for the queue.`

1. **The "first officer" metaphor flips within three lines.** Line 1 ("first officer **launching** claude") reads as *spacedock is the FO launching claude*; line 3 ("claude **is** your first officer") reads as *claude is the FO*. A new user can't form a stable model. Intended meaning is line 3; line 1 shouldn't call the launcher a "first officer." Candidate: `spacedock {v} · launching {host} as your first officer`.
2. **`spacedock status` is overloaded and undercuts the value prop.** The same command appears twice for two jobs — "run `spacedock status` to pick" (line 2) and "run `spacedock status` ... for the queue" (line 3). Worse: the whole pitch is "claude is your first officer" — the FO runs `status` *for* you. Telling the user to run it themselves (and "inside the session," ambiguous in a chat) contradicts that. Candidate for line 3: *"...ask your first officer for the queue."*
3. **Multi-workflow "pick" then launches anyway.** "N workflows detected (run `spacedock status` to pick)" announces an unresolved choice, then claude starts without one — leaving the user unsure which workflow the FO is on. Per the FO contract the disambiguation actually happens in-session (multiple → the FO presents the list on its first turn), so the banner instruction is redundant + misleading. Either say nothing actionable, or *"your first officer will help you pick."*

## Out of scope

The single-workflow happy path (clean today). The host's own banner ordering.

## Minor / to confirm at ideation

- The displayed `Version` is the compiled-in constant — confirm it reflects the real installed release (a dev build shows a stale `0.19.0`); a wrong version is a poor first impression and is exactly the version/channel correctness the 0.20.0 flip is about.
- `Workflow:` (singular label) + `N workflows detected` (plural content) reads slightly off.

## Acceptance criteria

(Ideation fills in — each verified by an `internal/cli` test over `launchBanner`'s rendered output for the single-workflow, multi-workflow, and none-detected cases, not by a source-grep.)
