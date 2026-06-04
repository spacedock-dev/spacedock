---
commissioned-by: spacedock@template
entity-type: task
entity-label: task
entity-label-plural: tasks
id-style: sd-b32
stages:
  defaults:
    worktree: false
    concurrency: 2
  states:
    - name: backlog
      initial: true
      gate: true
    - name: ideation
      gate: true
    - name: implementation
      worktree: true
    - name: validation
      worktree: true
      fresh: true
      feedback-to: implementation
      gate: true
    - name: done
      terminal: true
---

# Development Workflow Template

Refinement specialization for code that ships via PR/merge. Tasks move from a captain-curated backlog through ideation, get built in a dedicated worktree, are independently validated against the acceptance criteria, and land via PR review. The repo-mutation layer is active on `implementation` and `validation`; the `pr-merge` mod handles the PR lifecycle.

Use this template when the captain's mission is to ship code in a repo where work is reviewed and merged via PR. The bucket-noun stage names (`implementation`, `validation`, `done`) describe where the entity is sitting; the `pr-merge` mod removes any temptation to add `pr_open` or `awaiting_merge` stages because PR state is tracked on the `pr` field, not as a stage.

## File Naming

Each task lives as either:

- a flat markdown file `{slug}.md` (default), or
- a folder `{slug}/` containing `index.md` when the task produces sibling artifacts (transcripts, design notes, comparison tables) that belong with the tracker.

Slugs are lowercase, hyphens, no spaces. Example: `pluggable-id-style-collaboration-friendly-ids.md`.

## Schema

Every task file has YAML frontmatter. Fields are documented below; see **Task Template** for a copy-paste starter.

### Field Reference

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier; SD-B32 by default for collaborative cross-branch creation |
| `title` | string | Human-readable task name |
| `status` | enum | One of: backlog, ideation, implementation, validation, done |
| `source` | string | Where this task came from (issue, captain note, retrospective) |
| `started` | ISO 8601 | When active work began |
| `completed` | ISO 8601 | When the task reached terminal status |
| `verdict` | enum | PASSED or REJECTED — set at validation |
| `score` | number | Priority score, 0.0–1.0 (optional) |
| `worktree` | string | Worktree path while a dispatched agent is active, empty otherwise. Once set on first dispatch into a `worktree: true` stage, it stays set across all non-terminal advancements (stickiness) and clears at terminal merge. |
| `issue` | string | GitHub issue reference (e.g., `#42`) — optional cross-reference |
| `pr` | string | GitHub PR reference (e.g., `#57`) — set when a PR is opened for this task's branch |
| `mod-block` | string | Pending mod-declared blocking action, format `{lifecycle_point}:{mod_name}` |

## Stages

### `backlog`

A task enters backlog when it is first proposed. It carries a seed description but no design work. This is a captain-curated holding stage with a gate — the captain decides which tasks advance to ideation.

- **Inputs:** None — this is the initial state
- **Outputs:** A seed task file with title, source, brief description
- **Good:** Clear enough to recognize what the task is about
- **Bad:** Empty stub that even the captain cannot triage

### `ideation`

A task moves to ideation when the captain greenlights it for design work. The work here is to flesh out the problem, propose an approach, define acceptance criteria, and write a test plan.

- **Inputs:** The seed description and any relevant context (existing code, captain notes, related tasks)
- **Outputs:** A fleshed-out task body with problem statement, proposed approach, acceptance criteria (entity-level end-state properties with `Verified by:` clauses), and a test plan
- **Good:** Behavior-first, scoped, addresses a real need, AC items name end-state properties (not stage actions), test plan matches the AC's level of abstraction
- **Bad:** Vague hand-waving, scope creep, AC items written as imperative verbs, test plans that prove the wrong thing

### `implementation`

A task moves to implementation once its design is approved. The work happens in a dedicated worktree on a feature branch.

- **Inputs:** The fleshed-out task body from ideation with approach and acceptance criteria
- **Outputs:** The deliverable committed to the worktree branch with a stage report describing what was produced and where
- **Good:** Minimal changes that satisfy the AC, clean code, tests where appropriate, deliverable self-contained for validation
- **Bad:** Over-engineering, unrelated refactoring, skipping tests, leaving the deliverable incomplete for validation to finish

### `validation`

A task moves to validation after implementation is complete. A fresh agent independently verifies the deliverable meets the AC defined in ideation. The validator does not produce the deliverable — it checks what was produced.

- **Inputs:** The implementation summary, the AC from the task body, the worktree branch
- **Outputs:** A validation report covering each AC, with PASS/FAIL per criterion. Either gate-approval to `done` or rejection back to `implementation` with concrete fixes.
- **Good:** Reproduces each AC's `Verified by:` clause, reports actual evidence not assertions, exercises edge cases the AC implies
- **Bad:** Trusting implementation's self-report, skipping AC items, accepting "should work" without running the check

### `done`

Terminal state. The task's PR is merged and the entity is archived.

- **Inputs:** A merged PR (tracked via the `pr` field and the `pr-merge` mod's startup/idle hooks)
- **Outputs:** None — terminal. `completed` set, `verdict: PASSED`, entity archived.
- **Good:** Reached terminal via real merge, not by manual flag flip
- **Bad:** Marking done before the PR actually merged

## Recommended practices (opt-in)

These are proven dev-workflow disciplines a captain can adopt into this workflow's validation stage. They are recommended, not mandatory — the universal first-officer contract does not impose them, because a non-development workflow's acceptance proof may legitimately be a published artifact, a metric, or a human review rather than a test or command. Adopt them by copying the guidance into the `validation` stage's Outputs and Bad lists above when commissioning a code-shipping workflow.

### Test-first authoring

For a code or fixture deliverable, write the failing test first: write a test that captures the desired behavior, run it and watch it fail for the right reason, then write the minimum code to make it pass, and refactor green. The test is what the gate judges. This is a dev-workflow discipline — recommended, not mandatory; the universal contract does not impose it, because a non-development workflow's deliverable (a PRD, an analysis verdict, a published artifact) is judged against its own pre-fixed success criteria, not a test-first ritual.

### External-proof acceptance criteria

At the validation gate, require that each AC's cited evidence comes from a check OUTSIDE the task body — a test, a command's output or exit code, a file the change produces, or the resulting on-disk state. An AC whose only cited proof is review of the task's own prose ("verified by reviewing this task's decision section") proves only that the prose exists; it can never fail, so it does not satisfy the AC. Reject self-referential ACs. If the task's only deliverable is a decision with nothing shipped, do not recommend PASSED — the decision belongs in the roadmap, not a terminal dev task. (Behavioral enforcement of this rule, when a workflow wants it, is a workflow-opt-in `spacedock status --validate` guard — the workflow declares it; it is not universal binary behavior.)

### Detached adversarial audit

For high-stakes surfaces — a front-door launcher, the status/guard mutation paths, shipped contract/scaffolding, or CI/release machinery — treat a passing validation as necessary but not sufficient. Before merging, run (or dispatch) a read-only adversarial audit on a detached checkout of the merge result:

- **When:** the high-stakes surfaces above; routine low-blast-radius changes do not need it.
- **What:** the auditor works on a separate throwaway checkout (never the implementation worktree) and never mutates the deliverable. It tries to REFUTE the validation — construct an adversarial edit the deliverable's own tests should catch, and confirm they do. A test that stays green under a claim-breaking edit is a hole. Findings come in two tiers — `Material:` (a real correctness or test-strength hole) and `Polish:` (non-blocking); "refuted nothing material" is a valid recorded outcome.
- **How recorded:** material findings route back through the validation→implementation feedback flow (a `### Feedback Cycles` entry naming the audit and its adversarial edit); the gate is not clean until they close.

The audit catches the class of hole where the test passes but would also pass on a broken future edit — which a green suite cannot see itself.

## Workflow State

View the workflow overview:

```bash
spacedock status --workflow-dir {dir}
```

Output columns: ID, SLUG, STATUS, TITLE, SCORE, SOURCE.

Find dispatchable tasks ready for their next stage:

```bash
spacedock status --workflow-dir {dir} --next
```

## Task Template

```yaml
---
id:
title: Task title here
status: backlog
source:
started:
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
---

## Problem

{What is broken or missing, and why it matters now.}

## Proposed approach

{How the implementation will address the problem. Concrete enough that a worker can start.}

## Acceptance criteria

Each AC names a property of the finished task (not a stage action) and how it is verified.

**AC-1 — {End-state property.}**
Verified by: {grep / test name / file path / command a future reader can reproduce.}

## Test plan

{What tests verify the implementation, estimated cost, whether E2E is needed.}

## Out of scope

{What this task deliberately does not address.}
```

## Commit Discipline

- Commit status changes at dispatch and merge boundaries
- Commit task body updates when substantive
- Implementation commits land on the worktree branch; merge to main happens via the `pr-merge` mod after PR review

## Adoption

### Pre-fill stages

```yaml
- name: backlog
  initial: true
  gate: true
- name: ideation
  gate: true
- name: implementation
  worktree: true
- name: validation
  worktree: true
  fresh: true
  feedback-to: implementation
  gate: true
- name: done
  terminal: true
```

### Apply layers

- **repo-mutation**: fires on `implementation` and `validation` because both stages operate inside a worktree against the codebase. The `worktree: true` flag on those stages is the structural consequence; `validation` adds `fresh: true` to ensure independent perspective.

### Offer mods

- **pr-merge**: install by default. The shipping ritual for this template is open-a-PR-get-review-merge-to-main, which is exactly what `pr-merge` automates. Surface the offer in Phase 1 with this framing:

  > Because this workflow ships code via PR review, I'll install the **pr-merge** mod. This is the structural reason your stages can stay clean — you don't need a `pr_open` or `awaiting_merge` stage to model the PR step. The mod tracks PR state on the `pr` field, watches for merges in the background, and advances the entity to `done` when the PR lands. Stages describe where work is happening; the PR lifecycle is mod-managed.

  Skip the offer only if the captain explicitly says "no PR review" or "we commit straight to main."

### Inject entity-template snippet

Use the development snippet (problem / proposed approach / acceptance criteria / test plan / out-of-scope) shown in the Task Template section above.

### Surface variants

None. Development is a single coherent shape. Captains who want a different code-shipping cadence (e.g., direct-to-main, trunk-based, release-train) can adjust the stage list and skip the pr-merge mod, but those are edits to this template, not separate variants worth pre-baking.

### Confirmation prose

Surface this in Phase 1 once the template is selected:

> I'll set this up as a **development** workflow: tasks move through `backlog → ideation → implementation → validation → done`, with worktrees on `implementation` and `validation` and `validation` running with a fresh agent so it independently checks the AC.
>
> ID style is **sd-b32** by default, because development workflows usually involve multiple worktree branches creating tasks in parallel and you want the IDs to reconcile without coordination. If this workflow is single-writer (just you, no concurrent branches), I can switch to sequential — let me know.
>
> Default mod: **pr-merge** (PR lifecycle automation, framing above). I'll confirm the install at file generation time.
