---
title: Run the rejection journey in team mode
status: backlog
source: "Captain ruling, 2026-08-16: do not invoke a journey in single-entity bare mode unless it is specifically testing bare-mode behavior"
id: zqb683j8jth0tyr2eme231e2
---

## Problem

The live rejection-flow journey starts the FO in single-entity `-p` bare mode. In bare mode the contract's feedback flow is sequential fresh dispatch: the addressable-worker reuse path the rejection skill routes through cannot happen at all. So the journey runs in a mode that disables the mechanism it exists to prove. The lifecycle assertion (internal/ensigncycle/claude_runtime_helpers_test.go:209) compounds it: it demands one implementation dispatch and one implementation report section, while the fixture itself scripts a second report section on the rework round — correct two-cycle runs graded red (`implementation-worker-not-dispatched`) in 3 of 8 observed live runs.

## Proposed approach

Invoke the rejection-flow scenario in team mode on both runtimes, so the routing flow under test can actually run. Derive the team-mode determined shape from the fixture and skill (the fixture declares no context-budget probe, so reuse is the required path; the fixture mandates a second implementation report section) and align the lifecycle assertion and grading to that shape: strict, no multiple-path acceptance. A conforming two-cycle run must grade green; a fresh-dispatch-when-reuse-required run must grade red. If bare-mode rejection semantics deserve their own coverage, that is a separately named scenario decided at ideation, not this journey's default.

## Out of scope

- The FO never-prepares residual (separate behavior concern).
- Product code: this is scenario invocation and harness grading only.

## Evidence pointers (FO, post-filing)

Two captain-ordered investigations sharpen this entity's scope; ideation must read both:

- `_debriefs/2026-08-16-02-live-harness-audit.md` — full live-suite audit. Findings 1, 2, 7, 8, 10 are this journey's surface: the codex leg's two graded assertions are mutually contradictory on the fixture's own scripted end-state (exact duplicate heading mandated by the fixture, hard-errored by the section selector); the Claude leg's either/or reviewer acceptance has no ordering anchor and no implementation-producer check; two `observed` checks are tautologies; the fixture-determined Cycle 2 line is never asserted.
- `_debriefs/2026-08-16-03-rejection-shape-attribution.md` — blind two-sided shape analysis of the five preserved runs. Key corrections to this entity's filed text: on codex the followup surface was LIVE, so the determined fix route is REUSE (spawns=1), not bare-mode fresh dispatch — the fresh-dispatch derivation governs the Claude `-p` leg only. Both codex "green" runs used a single-worker chain that re-reviewed its own fix, violating the skill's reviewer-independence ban — undetected because the codex branch swapped out the reviewer-flow assert; the one fixture-literal ensign was graded red by the duplicate-heading error. Minimal corrections identified there (prefix-count sections, fixture-side distinct cycle-2 title, reviewer-topology assert for codex) are candidate design inputs, not pre-approved decisions.
