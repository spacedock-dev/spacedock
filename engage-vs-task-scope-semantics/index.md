---
title: Define engage versus task-scoped First Officer semantics
status: backlog
score: "0.70"
source: "Captain filing request."
id: syhj6n4v9rb4ayasvjwf4zmx
---

# Define engage versus task-scoped First Officer semantics

## Problem

The First Officer contract currently makes workflow engagement and a task-specific request too easy to conflate. Explicitly engaging a workflow reasonably authorizes workflow-wide event-loop behavior. Filing, reviewing, dispatching, or resuming one named task does not necessarily authorize touching every ready or merge-blocked task in that workflow.

Today, unscoped instructions such as “keep moving,” workflow-wide mod-block scans, status next selection, and idle hooks can make a task-scoped interaction spill into unrelated work. The same ambiguity distorts async waiting: the First Officer may refuse to enter idle monitoring because unrelated workflow work exists, or may stop after a worker completes instead of advancing the next in-scope action.

## Proposed direction

Introduce an explicit captain-authorized active-scope concept across the shared First Officer and runtime contracts:

- engage workflow establishes workflow scope and authorizes the existing workflow-wide event loop.
- A request naming or creating specific tasks establishes task scope containing only those task references.
- Keep moving means advance ready work inside active scope; it never widens scope.
- A captain correction, review response, or external handoff preserves or narrows scope unless the captain explicitly expands it.
- Status selection, merge-block recovery, gates, idle hooks, and next-action routing filter through active scope.
- A task handed to an external runtime may end the local task loop even while unrelated workflow work remains.
- Async idle monitoring begins when no actionable in-scope work remains. Unrelated workflow work does not prevent a task-scoped wait; explicit workflow engagement still requires ready workflow work to advance first.
- Captain-facing state makes the current scope visible enough to explain why a task is or is not being advanced.

Prefer a mechanism or executable behavioral test over prose alone. The contract should identify the authoritative scope carrier and filtering seam rather than relying on the model to remember a conversational distinction.

## Acceptance criteria

1. Task scope excludes unrelated work. In a fixture containing a target task, an unrelated ready task, and an unrelated merge-blocked task, a task-specific filing, review, or dispatch trace advances only the target. It performs no status mutation, merge guard, GitHub action, or dispatch for unrelated tasks.
2. Explicit engagement retains workflow behavior. The equivalent fixture under explicit workflow engagement continues to discover and advance workflow-wide ready and merge-blocked work according to the existing event loop.
3. Waiting is scope-aware and completion advances. Task scope may enter async idle monitoring when its own work is exhausted even if unrelated workflow work exists. Workflow scope must advance ready workflow work before monitoring. A worker completion triggers durable verification and the next in-scope action rather than a completion-only stop.
4. Corrections and handoffs do not widen scope. A captain correction preserves or narrows the active task set. An external-runtime handoff terminates the local loop for that task unless further local work was explicitly retained.
5. Scope is observable and enforced. Behavioral coverage records the selected scope and proves task filtering. Contract-lint or wording checks may support this, but prose presence alone does not satisfy the criterion.
6. Existing focused workflows remain compatible. Single-task and explicitly workflow-engaged behavior remain unchanged except for the new scope observability and filtering guarantees.

## Test plan

Build a host-neutral scenario with one target task, one unrelated ready task, and one unrelated merge-blocked task. Exercise it once with a task-scoped request and once with explicit workflow engagement. Extract or simulate the runtime event trace and assert the exact task IDs passed to status mutation, dispatch, merge guard, gate, and wait transitions.

Add negative mutants that remove the scope filter from next-action selection, merge-block recovery, or post-completion routing; each mutant must fail. Add adapter-specific trace coverage for supported multi-agent runtimes where their event-loop bindings differ. Keep the async timeout value owned by the existing timeout task; this task owns scope and next-action semantics.

## Out of scope

- Changing the five-minute async wait timeout policy.
- Multi-workflow engagement in one active scope.
- Processing any currently ready or merge-blocked task.
- Changing the product behavior of tasks used only as scope fixtures.
