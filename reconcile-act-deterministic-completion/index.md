---
id: vcew2gxw6813yensgewae7tv
title: "spacedock dispatch reconcile --act — deterministic Class C completion (the trigger that completes pr-complete's determinism)"
status: backlog
source: "captain + p2 cycle-3 (2026-06-04) — p2 (pr-complete-binary-command) makes the post-merge ceremony a deterministic PAYLOAD, but a prose mod cannot deterministically INJECT it (the captain's seam: 'how can the pr mod inject deterministic command to dispatch complete'). Resolution: the deterministic TRIGGER is reconcile --act, not the mod. reconcile already detects the merged-but-un-advanced PR (Class C, deterministic); --act has the BINARY run dispatch complete on detection — binary->binary, zero FO prose. This is roadmap #6, elevated from 'queue later' to 'the trigger that completes p2's determinism story.'"
score: "0.33"
worktree:
started:
completed:
verdict:
issue:
---

`reconcile` (internal/dispatch/reconcile.go) already DETECTS drift deterministically and emits a report the FO acts on via prose. `--act` makes the binary ACT on the detection, closing the prose-in-the-loop seam. The motivating case is Class C (merged-but-un-advanced PR): its action is "run `dispatch complete {slug}`" (p2), making the detect→complete loop fully deterministic with no FO interpretation.

## Problem

Today reconcile is dry-run-only: it detects A/B/C/D/E drift and the FO reads the report and acts by prose. For the post-merge case (Class C), that means the deterministic `dispatch complete` payload (p2) is still triggered by FO prose — deterministic payload, non-deterministic delivery. The seam the captain named.

## Proposed approach (ideation firms)

Add `--act` to `spacedock dispatch reconcile` (roadmap #6), preserving dry-run-by-default safety (act only under `--act`; every action surfaced on stdout for audit). Each drift class maps to a well-defined binary op:
- **C (un-advanced merged PR) → run `dispatch complete {slug}`** (depends on p2 `pr-complete-binary-command`). This is the headline: it completes p2's determinism story — `reconcile --act` becomes the deterministic trigger; the pr-merge mod's post-merge hook is then DELETED (advancement leaves the mod; the mod keeps only its captain-gated PR-open hook).
- **D (stale branch) → `git pull --rebase origin next`** (halt on conflict).
- **E (stale local main) → fetch + reset + rebuild.**
- **A/B (lingering/superseded agents) → cooperative shutdown** — NOTE: agent shutdown is a Claude/host TOOL, not shell-able; A/B `--act` likely stays FO-side or is out of scope for the binary (decide at ideation — this is the same non-shell-able boundary as team teardown).

## Out of scope / dependencies

- **Depends on p2** (`pr-complete-binary-command`) landing first — Class C's action invokes `dispatch complete`. Sequenced AFTER p2 (the roadmap #6→#3 dependency).
- Does NOT make the captain-gated PR-OPEN deterministic — that hook stays in the pr-merge mod (the PR-approval guardrail is judgment, not mechanism).
- The non-shell-able classes (A/B agent shutdown) may stay FO-side.

## Acceptance criteria (seed)

- **AC-1 (seed):** `reconcile --act` is dry-run-by-default; with `--act` it runs each acted class's binary op and surfaces every action on stdout — verified by a fixture per acted class asserting the action ran + was reported.
- **AC-2 (seed):** Class C `--act` invokes `dispatch complete {slug}` and the entity reaches archived terminal — verified by an end-to-end fixture (merged-PR drift → reconcile --act → archived), and the pr-merge mod's post-merge hook is removed (the mod retains only PR-open).
- **AC-3 (seed):** Safety — `--act` never acts without the flag; conflicts/halts (Class D rebase conflict) surface and stop rather than force.

## Notes

The determinism trio: p2 `pr-complete-binary-command` (payload) + this (trigger) + the live `pr-lifecycle-from-boot` scenario (p2 AC-5, proves the boot→PR→complete loop). Roadmap: docs/dev/_proposals/binary-simplification-roadmap.md #6. Sibling 0.19.6 contract-decomposition line: gate-presentation-skill-extraction, feedback-rejection-flow-skill-extraction.
