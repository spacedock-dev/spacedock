---
title: Pi bootstrap prompt parity — match claude/codex warmth; "Use $spacedock:first-officer" is the cold outlier
status: backlog
source: "Captain (2026-06-20): the pi bootstrap prompt (internal/cli/pi.go:20) is 'Use $spacedock:first-officer for this whole Pi session.' — a bare mechanism trigger. claude (frontdoor.go:25) and codex (frontdoor.go:434) get 'You totally got this. Take your time. I love you. And tell all subagents and team members you love them too. Engage.' (codex appends 'Assume $spacedock:first-officer for the entire session.'). Pi is the cold outlier — pure mechanism, zero warmth. The skill is the contract (single source of truth), but the launch moment is the one chance to frame the commissioning, and pi's reads like a config line."
score:
started:
completed:
verdict:
worktree:
issue:
sprint:
sprint-readiness:
id: 7vtn8yda8vn0p7y8am3f43c8
---

# Pi bootstrap prompt parity

## End value

`spacedock pi` launches the FO with a bootstrap prompt that matches claude/codex's warmth and commissioning framing — not the current bare mechanism trigger. The launch moment frames the officer's posture (you got this; love your crew; engage) and triggers the skill, instead of reading like a config line. Pi is no longer the cold outlier among the three hosts.

## Problem — root cause already determined

The three bootstrap prompts (verified in source):

- **claude** (`internal/cli/frontdoor.go:25`): `"You totally got this. Take your time. I love you. And tell all subagents and team members you love them too. Engage."`
- **codex** (`internal/cli/frontdoor.go:434`): `"You totally got this. Take your time. I love you. And tell all subagents and team members you love them too. Engage. Assume $spacedock:first-officer for the entire session."`
- **pi** (`internal/cli/pi.go:20`): `"Use $spacedock:first-officer for this whole Pi session."`

Pi's prompt is pure mechanism (trigger the skill) with zero warmth or commissioning framing. Claude/codex share a warm core ("You totally got this. Take your time. I love you. And tell all subagents and team members you love them too. Engage.") — codex appends the skill-trigger clause; claude's skill-trigger comes from the `--agent spacedock:first-officer` flag instead. Pi has neither the warm core nor the flag-based trigger, so it falls back to the bare sentence.

## Approach

Bring pi's prompt to parity. Two options:

- **(a) Mirror codex exactly:** `"You totally got this. Take your time. I love you. And tell all subagents and team members you love them too. Engage. Assume $spacedock:first-officer for the entire session."` — byte-identical to codex (codex already carries the skill-trigger clause, which pi needs since pi has no `--agent` flag). Recommend — true parity, zero drift risk, the warm core + the trigger in one.
- **(b) Pi-specific warm prompt:** a pi-flavored variant of the warm core. More character, but introduces a third prompt to maintain and risks drift. Reject unless the captain wants a distinct pi persona.

Pick (a). The prompt is a single `const` swap in `internal/cli/pi.go:20`; no other change. The skill remains the contract (single source of truth); the prompt is the launch-moment framing + trigger, exactly as on codex.

## Scope

In scope:
- Replace `piBootstrapPrompt`'s value with the codex prompt (the warm core + `Assume $spacedock:first-officer for the entire session.`).

Out of scope:
- Changing the skill or the contract — the prompt is the trigger, not the contract.
- claude/codex prompts — they're already the warm reference.
- The `--agent` vs skill-trigger mechanism difference — pi has no `--agent` flag, so the prompt carries the trigger (as codex does).

## Acceptance criteria (provisional — ideation finalizes; proof = behavior)

**AC-1 — `piBootstrapPrompt` matches codex's prompt (warm core + skill-trigger clause).**
Verified by: a Go test asserting `piBootstrapPrompt == codexBootstrapPrompt` (byte-identical) — OR, if a distinct pi persona is chosen, asserting the warm core is present + the skill-trigger clause is present. (Binding two independent values — the two prompt constants — that can diverge: legitimate structural check, not prose-grep.)

**AC-2 — `spacedock pi` launches with the warm prompt (the smoke / a live launch observes it).**
Verified by: the pi-live smoke or a live launch capturing the appended prompt and confirming the warm core + skill-trigger appear in the launch argv. Behavior-bound, not a prose claim.

## Test plan

- Go test (AC-1): `piBootstrapPrompt == codexBootstrapPrompt` (or the warm-core + trigger-clause assertion).
- Live/harness (AC-2): the smoke's captured launch argv contains the warm core + skill-trigger.
- `pi-live` lane (touches `internal/cli/pi.go` — pi-only surface).

## Related

- `internal/cli/pi.go:20` (`piBootstrapPrompt`) — the source of truth.
- `internal/cli/frontdoor.go:25,434` (`bootstrapPrompt`, `codexBootstrapPrompt`) — the warm reference.
- The 0223 Shaping FO debrief — surfaced the "without love" observation during the fnm-fix smoke.
