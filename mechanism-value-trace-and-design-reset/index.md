---
title: Generalize mechanism-value tracing and design-reset routing
status: backlog
source: Captain correction after the 7h, 1k, and e6j mechanism failures, 2026-07-14
started:
completed:
verdict:
score:
worktree:
issue:
milestone: 0.26.0
id: wywxp0zmxw33q5qfv1t3vyv6
---

Carry the development workflow's mechanism/value discipline into the shipped development template and First Officer rejection contract. This is separate from the direct `docs/dev/README.md` workflow change.

## Problem

Three implementation efforts exposed the same missing general rule. The clearest incident, 7h, expanded a disposable Zellij/WASM smoke into a custom PTY controller with terminal/session setup, foreground process-group checks, raw-byte canaries, lease publication, signal coordination, readiness races, and cleanup state. It consumed about 4.25 active hours across 13.5 elapsed hours, then the final packet failed 3/3 on harness readiness while focused profile runs passed; nothing merged. The successful replacement used an existing real-terminal boundary: isolated tmux -> real Zellij -> send literal keys -> capture screen/native state -> verify cleanup.

The reusable invariant is: a test harness observes the supported runtime; it does not become a second implementation of it. More generally, at dispatch, during implementation, and at rejection, every new mechanism names the user-visible/value AC it serves, the simplest available alternative, and why that alternative is insufficient. When rejection shows the mechanism architecture is wrong but the value remains reachable by a simpler route, stop the feedback loop and re-scope instead of repairing the mechanism.

## Required scope

- Update `skills/commission/references/templates/development.md` so newly commissioned development workflows carry the mechanism/value trace, boring-proof-first rule, implementation tripwire, 90-minute architecture-review timebox, and mechanism-failure design reset.
- Update the First Officer contract and `feedback-rejection-flow` so an architectural mechanism failure is classified before feedback-cycle accounting. It must not increment the cycle or redispatch implementation; it escalates a scope/design reset.
- Preserve ordinary feedback routing for defects in the intended product.
- Keep the harness-specific tripwire explicit: for multiplexers, `setsid`, process-group control, raw PTY writes, or a second lifecycle supervisor require architecture review and an existing real-terminal harness is tried first.
- Add no binary policy engine, new daemon, new lifecycle state, or parallel test controller.

## Acceptance criteria

- **AC-1 (VALUE):** In a live/fixture rejection journey where the product value remains valid but an enabling harness fails and a simpler supported-runtime route is available, the FO records rejected experimental evidence, adds no feedback cycle, dispatches no implementation repair, and surfaces a scope/design reset.
- **AC-2:** The paired product-defect journey still records the feedback cycle and routes concrete fix work to implementation.
- **AC-3:** A freshly commissioned development workflow includes the dispatch/implementation/rejection discipline and names the harness tripwire without requiring project-specific Zellij wording.
- **AC-4:** Proof reuses the existing smallest-mechanism and rejection-flow scenario infrastructure. It must grade tool/dispatch/state outcomes; contract-text matching may verify only generated structure and cannot satisfy AC-1 or AC-2.

## Test boundary

Extend existing scenarios; do not build a new harness. The negative controls are: a mechanism failure incorrectly bounced to implementation, a product defect incorrectly treated as a design reset, and generated workflow structure missing one of the three decision points. Run the focused scenario/contractlint suites, then the repository-required full and race suites.
