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

The same mistake appears in policy placement. The shared First Officer dispatch
contract names a prose-polisher convention and when to use it, even though
`comm-officer` is declared only by this workflow's mod. During the `c6` ideation
drive, the generic contract was loaded but the mod-specific routing guidance did
not reach the writer; the task grew to roughly 9,200 words before a late polish
pass. A generic contract should provide standing-teammate mechanics, not absorb
the policy of one workflow-specific teammate.

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

Apply the rule to policy ownership as a second concrete example:

- The shared FO contract owns only generic standing-teammate discovery,
  addressing, lifecycle, and bounded routing mechanics.
- A workflow and its mod own the teammate's purpose, triggering conditions,
  required-versus-best-effort behavior, timeout, and fallback record. The shared
  contract must not name `comm-officer`, a prose-polisher convention, or prose
  categories.
- This dev workflow requires complex or long ideation prose to route through
  its `comm-officer` before commit. If unavailable after the existing bounded
  wait, the writer proceeds and records the fallback.
- Applicable workflow/mod guidance must reach the producing writer through the
  smallest existing dispatch/stage mechanism. Do not add a registry or policy
  engine without two present cases that need one.

## Acceptance criteria

**AC-1 (general rule): The shared contract expresses the generalized rule and its evidence-based single-case exception.**

**AC-2 (developer example): The developer-specific wording remains an example of the general rule, not a competing policy.**

**AC-3 (low ceremony): The addition guides discretionary design choices without requiring ritual for deterministic state edits or every dispatch.**

**AC-4 (workflow-owned teammate policy): The shared contract contains only generic standing-teammate mechanics; the workflow/mod owns prose-polish triggers, fallback, and writer guidance.**

**AC-5 (behavior reaches the writer): A complex or long ideation draft routes through this workflow's comm officer before commit, while an inapplicable short draft does not; teammate unavailability follows the bounded fallback and records it.**

## Test plan

Ideation should identify the smallest contract insertion point and propose a
behavioral or fixture-backed proof that the guidance routes an actual
overengineering choice toward the direct alternative. Do not add prose-grep or
style-lint enforcement.

Add a workflow fixture that drives three cases through the real dispatch/write
boundary: qualifying ideation routes before the state commit; a short,
non-applicable draft produces no polish route; and an unavailable teammate
proceeds after the bounded fallback with a durable note. A structural ownership
check may ensure workflow-specific names do not leak into the shared contract,
but it cannot substitute for the behavioral fixture.
