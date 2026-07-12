---
title: Move prose-polish routing policy into the workflow
status: backlog
score: "0.70"
source: "Captain correction after c6 ideation review."
id: csb4c89dteavbq1htdac7fwm
---

# Move prose-polish routing policy into the workflow

## Problem

The shared First Officer dispatch contract names a standing prose-polisher convention and describes when drafts may route through it. That policy is not universal: comm-officer exists because this workflow declares it in its comm-officer mod.

The split caused a live failure. The FO loaded the generic dispatch contract but did not consume the mod-specific guidance or pass it to the c6 ideation writer. The task grew to roughly 9,200 words before the comm officer was used late for cleanup. The workflow and mod contained the relevant scope, but routing was optional, late, and absent from the writer's assignment.

## Proposed direction

Keep the shared contract generic. It may define standing-teammate discovery, addressing, lifecycle, and bounded best-effort routing, but it must not name a prose-polisher convention, comm-officer, prose categories, or workflow-specific triggers.

Make this workflow and its mod the sole owners of prose-polish behavior:

- Define which artifacts qualify, including complex or long ideation bodies before commit.
- Define required versus best-effort routing, the existing bounded timeout, and the durable fallback note.
- Ensure applicable routing guidance reaches the producing writer through the smallest existing stage or dispatch mechanism.
- Avoid a generalized policy registry or new routing engine unless two current workflows require one.
- Preserve direct captain chat, short statuses, logs, and other excluded content as non-routed.

## Acceptance criteria

**AC-1 (ownership): The shared First Officer contract contains only generic standing-teammate mechanics; prose-polish purpose and triggers live in this workflow and mod.**

**AC-2 (writer routing): A qualifying ideation writer receives the workflow's routing policy and contacts the comm officer before committing the task body.**

**AC-3 (bounded fallback): If the comm officer is unavailable beyond the declared bound, the writer proceeds and records the fallback without blocking the workflow.**

**AC-4 (no unnecessary routing): Short or excluded content produces no comm-officer route.**

**AC-5 (no new abstraction): The change reuses the existing workflow, mod, dispatch package, and standing-teammate surfaces without adding a registry or generalized policy engine.**

## Test plan

Create a fixture workflow that declares the comm-officer mod and drives the real assignment and write boundary:

- A complex or threshold-crossing ideation body must produce a comm-officer route before the state commit.
- A short or explicitly excluded artifact must produce zero polish routes.
- An unavailable teammate must hit the bounded fallback, commit successfully, and leave the required durable note.

Add a host-neutral dispatch-package test proving applicable workflow or mod guidance reaches the writer without hard-coding comm-officer into the shared contract. A structural ownership check may guard against workflow-specific names in the shared core, but it cannot substitute for the behavioral fixture.

## Out of scope

Changing standing-teammate spawn timing, making prose polish load-bearing beyond this workflow, redesigning the comm officer's prose rules, or implementing a general routing-policy registry.
