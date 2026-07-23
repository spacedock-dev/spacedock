# IDEATION GATE — First Officer recorded gate lifecycle (`6y`)

Recommendation: **APPROVE the cycle-3 proof-boundary reset and resume implementation.**

## Capability preserved

Every supported host must still execute the six ordered gate commands, consume exactly once, issue exactly one successor `dispatch build`, and leave exactly one new durable successor effect. Static route coverage and six command-deletion controls remain mandatory.

## Proof correction

- Live journeys prove observable command order, consumed state, one dispatch build, and one durable stage report/commit.
- Host-native fixtures prove each runtime’s transport details and refusal controls.
- Codex no longer must expose a nonexistent public spawn/completion event.
- One representative live route per host replaces every-route-by-every-host multiplication; static tests continue to cover every gate-entry route.

## Surface reduction

The plan deletes 479 lines of cross-host forensic machinery and permits at most 80 replacement lines. The current 1,879-line branch must fall to at most 1,480 added lines: at least 399 lines, or 21.2%, removed. No production Go, recorder, runtime, schema, command, or public host interface is added.

## Findings

The reset removes an inferred proof obligation without narrowing the product. AC-2 through AC-7, the deferred lifecycle skill, documentation, byte limits, and full/race/three-host validation remain unchanged.

## Decision

Approve to implement only this measured proof cut. Revise for a concrete lost value guarantee. Hold only for a named prerequisite.
