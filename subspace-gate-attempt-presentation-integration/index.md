---
title: Subspace gate-attempt presentation integration
status: backlog
score: "0.75"
source: "Split from 3k after captain clarification during Briefing 7 dogfood, 2026-07-19."
id: krdyxrhdcstqv543z70rekkm
---

# Subspace gate-attempt presentation integration

## Problem

The durable gate-state design in 3k should not also own Subspace transport integration. The dogfood path currently requires hand-built Zellij commands, can render an empty float when stderr is redirected, and the beta.5 `present` controller can fail with `present-child protocol ended early: EOF`. A raw launch may also return only an `open` handoff, leaving no blocking result boundary for the presenting ensign or waiting First Officer.

## Required capability

Ship the separate adapter that lets a gate-attempt ensign present one explicit immutable Briefing package through a single blocking command. Subspace owns floating/transport, continuation, ProbeResult/comparison joining, and atomic result retention. The ensign stays addressable until the review closes; the First Officer waits and consumes only the validated retained result. Pane creation, `open`, EOF, empty output, or malformed output are nonterminal/incomplete and must preserve diagnostics without retrying.

The presentation must surface any semantic Probe delta separately when the installed Subspace UI cannot render it. This task integrates with 3k's durable Briefing/Resolution binding seam but does not redesign that schema.

## Acceptance criteria

- **AC-1:** One ensign-facing command opens the complete explicit Briefing in the captain's active terminal transport and remains blocking until a terminal validated result or a retained incomplete continuation.
- **AC-2:** The command renders visibly with TUI stderr attached, handles beta.5 `open` continuation, and cannot report success on pane creation or controller EOF.
- **AC-3:** A completed result is atomically retained, exact-Briefing-bound, and delivered to the still-addressable ensign; the FO's 300-second waits continue until that completion signal.
- **AC-4:** ProbeResult/comparison presentation is owned by Subspace; until supported in-TUI, the adapter presents a separate frozen semantic-delta summary without changing the Briefing digest.
- **AC-5:** Failures retain exact diagnostics/result/continuation paths and do not silently retry or mutate Spacedock gate state.
