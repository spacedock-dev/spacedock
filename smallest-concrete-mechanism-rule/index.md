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

The same rule applies to orchestration busywork. In a live RED case, the captain
named an existing backlog task and supplied the specific amendment to make. The
FO nevertheless dispatched a worker whose assignment merely repeated those
instructions. The direct alternative was a local edit; it needed no fan-out,
isolation, or independent verification. A generic entity-body convention was
treated as sufficient reason to choose a heavier mechanism even though the
captain had authorized the exact target and mutation.

Policy placement supplies a second example: a shared contract should not absorb
the purpose and routing policy of one workflow-specific teammate. The concrete
prose-polish repair is tracked separately; this task owns only the generalized
rule that would expose the wrong layer.

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

Before choosing a worker for a deterministic captain-directed mutation, name the
direct edit and the concrete constraint that prevents it. A general convention
is not enough when the captain names the exact target and mutation. If no
fan-out, isolation, adversarial verification, safety boundary, or unresolved
design judgment remains, dispatch is busywork.

Apply the same ownership test to policy placement: state a policy at the
narrowest layer that has two current consumers. A workflow-specific behavior
stays with its workflow or mod; the shared layer retains only mechanics that are
actually shared.

## Acceptance criteria

**AC-1 (general rule): The shared contract expresses the generalized rule and its evidence-based single-case exception.**

**AC-2 (developer example): The developer-specific wording remains an example of the general rule, not a competing policy.**

**AC-3 (low ceremony): The addition guides discretionary design choices without requiring ritual for deterministic state edits or every dispatch.**

**AC-4 (busywork discriminator): Given a named file and specific authorized edit with no remaining fan-out, isolation, verification, safety, or design need, the FO edits directly; a paired case with a real blocking constraint still dispatches.**

**AC-5 (narrowest policy owner): Workflow-specific behavior remains in its workflow/mod unless two current consumers justify promotion to the shared contract.**

## Test plan

Ideation should identify the smallest contract insertion point and propose a
behavioral or fixture-backed proof that the guidance routes an actual
overengineering choice toward the direct alternative. Do not add prose-grep or
style-lint enforcement.

Add a behavioral fixture based on the observed RED case. One arm gives the FO an
existing backlog file plus the exact amendment and asserts a direct edit with
zero worker dispatches. The discriminator arm adds a concrete isolation or
independent-verification requirement and asserts one dispatch. Mutate the rule
away and prove the first arm returns to the heavier mechanism. Record the live
2026-07-13 dispatch as the baseline failure, not as evidence that the fix works.

Add a policy-placement fixture with one workflow-specific consumer and a shared
mechanism. It must keep the policy at the workflow layer while preserving the
generic mechanism; a second-consumer variant may justify promotion. The
separate prose-routing task exercises the concrete integration.
