---
title: Generalize the smallest-concrete-mechanism rule
status: backlog
source: "Captain request 2026-07-13."
started:
completed:
verdict:
score:
worktree:
issue:
id: 5wmx24qwygh711wsmmqb1qwj
---

Capture a concise, reusable anti-overengineering rule in the shared Spacedock
contract, with a developer-specific adaptation as an example rather than a
second universal policy.

## Problem

The existing smallest-sufficient-mechanism guidance is sound but easy to apply
only at the workflow level. A recent Safehouse change turned “future tmux” into
a one-item registry and custom parser/merger, because no short rule required
the implementer or validator to justify the abstraction with a present case.

## Proposed direction

Add this generalized rule to the shared contract at the appropriate decision
point:

> Prefer the smallest concrete mechanism that achieves the current outcome. Do
> not add a layer, workflow, agent, stage, gate, state record, configuration
> surface, or extensibility scheme unless you name two current concrete cases
> that require it. Before acting, state the direct alternative and why it fails
> for this task.

Keep this narrow exception adjacent to it:

> A single case can justify structure only for a proven external contract,
> safety invariant, required isolation, or compatibility boundary; cite that
> evidence.

Retain this developer-specific adaptation as an illustrative example, not a
separate rule:

> Prefer the most direct implementation. Do not add a type, collection,
> parser, or extensibility layer unless you name two current concrete cases
> that require it. Before coding, show the simpler rejected alternative and why
> it fails.

## Acceptance criteria

**AC-1 (general rule): The shared contract expresses the generalized rule and its evidence-based single-case exception.**

**AC-2 (developer example): The developer-specific wording remains an example of the general rule, not a competing policy.**

**AC-3 (low ceremony): The addition guides discretionary design choices without requiring ritual for deterministic state edits or every dispatch.**

## Test plan

Ideation should identify the smallest contract insertion point and propose a
behavioral or fixture-backed proof that the guidance routes an actual
overengineering choice toward the direct alternative. Do not add prose-grep or
style-lint enforcement.
