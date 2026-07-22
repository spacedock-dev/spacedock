---
commissioned-by: spacedock@0.25.0
entity-type: task
entity-label: task
entity-label-plural: tasks
id-style: sd-b32
stages:
  # Stage names must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$ (kebab-case lowercase, no underscores or spaces); `status --validate` rejects others.
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

# Personal Site Development Workflow

Tasks for the personal site move from a captain-curated backlog through ideation, get built in a dedicated worktree, are independently validated against the acceptance criteria, and land via PR review. The repo-mutation layer is active on `implementation` and `validation`, and the `pr-merge` mod handles the PR lifecycle.

## File Naming

Each task lives as either:

- a flat markdown file `{slug}.md` (default — use this unless the entity produces many artifacts), or
- a folder `{slug}/` containing `index.md` as the canonical entity file, when the task produces per-stage artifacts (draft versions, transcripts, outputs) that belong alongside the tracker.

Slugs are lowercase, hyphens, no spaces. Example: `my-feature-idea.md` or `my-feature-idea/index.md`. The status scanner recognizes both forms; `--set` and `--archive` resolve the slug either way, and folder entities archive as a whole folder into the workflow's archive directory.

## Schema

Every task file has YAML frontmatter. Fields are documented below; see **Task Template** for a copy-paste starter.

### Field Reference

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier, format determined by id-style in README frontmatter |
| `title` | string | Human-readable task name |
| `status` | enum | One of: backlog, ideation, implementation, validation, done |
| `source` | string | Where this task came from |
| `started` | ISO 8601 | When active work began |
| `completed` | ISO 8601 | When the task reached terminal status |
| `verdict` | enum | PASSED or REJECTED — set at final stage |
| `score` | number | Priority score, 0.0–1.0 (optional). Workflows can upgrade to a multi-dimension rubric in their README. |
| `worktree` | string | Worktree path while a dispatched agent is active, empty otherwise |
| `issue` | string | GitHub issue reference (e.g., `#42` or `owner/repo#42`). Optional cross-reference, set manually. |
| `pr` | string | GitHub PR reference (e.g., `#57` or `owner/repo#57`). Set when a PR is created for this entity's worktree branch. |

### ID Style

The `id-style` frontmatter setting controls the operator-facing ID strategy:

- `sequential`: `id` is required and stores the next zero-padded numeric value from `status --next-id`, counting active and archived entities.
- `sd-b32`: `id` is required and stores the full stable 24-character lowercase SD-B32 stored ID from `status --next-id --id-seed <slug-or-title>`. SD-B32 is Spacedock Base32: SHA-256 digest material formatted with Spacedock's human-safe alphabet `0123456789abcdefghjkmnpqrstvwxyz`. Status tables show shorter display/address prefixes computed from active plus archived entities. `status --boot` reports `ID_STYLE: sd-b32`, `NEXT_ID: {candidate}`, and `MIN_PREFIX: 2`.
- `slug`: `id` is optional; the effective ID is the entity slug. `status --next-id is not applicable for id-style: slug` because the slug comes from the title.

SD-B32 display/address prefixes can lengthen after another branch adds a colliding prefix, while stored IDs remain stable. Use `status --validate` before trusting workflow state and `status --resolve <ref>` to resolve slugs, stored IDs, or sd-b32 address prefixes.

Copyable README frontmatter examples:

```yaml
id-style: sequential
```

```yaml
id-style: sd-b32
```

```yaml
id-style: slug
```

SD-B32 examples:

| Workflow size | Stored `id` examples | Display/address examples |
| --- | --- | --- |
| 10s of entities | `4k9q2m7x8c3v9r5t6w2p0n1h`, `8t5n0p2w6j9r4c8x1m7q3v5k` | `4k`, `8t` |
| 100s of entities | `9m2c7v4xq8j3h6t0p5w1r8n2`, `9m2cq8j3h6t0p5w1r8v7x4kn` | `9m2c7`, `9m2cq` |
| 1000s of entities | `v7k3q9x2m5c8h6t0p1w4r8n2`, `v7k3qrv5t9p3j6n2w8c4x1mk` | `v7k3q9`, `v7k3qr` |

Generated IDs make concurrent and offline creation safer because creators do not share a central counter. Migration from existing sequential workflows is manual migration in this release: validate the target style, update README/entity frontmatter deliberately, and defer rewrite automation to a separate tracked task.

## Stages

### `backlog`

A task enters backlog when it is first proposed: a seed description, no design work. Captain-curated holding stage — the gate decides which tasks advance to ideation.

- **Inputs:** The captain's seed note, a GitHub issue, or a retrospective item about the personal site.
- **Outputs:** A task file with `title`, `source`, and a one-paragraph problem statement in the body; no design work.
- **Good:** The problem is stated in terms of what a site visitor or the site owner experiences today.
- **Bad:** Backlog entries that already prescribe an implementation, or entries with no way to tell when they would be done.

### `ideation`

The captain greenlights a task for design: flesh out the problem, propose an approach, define acceptance criteria as entity-level end-state properties with `Verified by:` clauses, and write a test plan that matches the AC's level of abstraction.

- **Inputs:** The backlog problem statement, the current site codebase, and any linked issue.
- **Outputs:** `## Problem`, `## Proposed approach`, `## Acceptance criteria` (each AC an end-state property with a `Verified by:` clause), `## Test plan`, and `## Out of scope` filled in on the task body.
- **Good:** Each AC names a property of the finished task and cites a check a future reader can reproduce — a test name, a command, a file path.
- **Bad:** ACs that restate stage actions ("implement the toggle"), or whose only proof is a review of the task's own prose.

### `implementation`

The design is approved and the deliverable is built in a dedicated worktree on a feature branch — minimal changes that satisfy the AC, self-contained for validation.

- **Inputs:** The approved ideation AC and test plan; the site repo checked out in this task's worktree.
- **Outputs:** Commits on the task's feature branch satisfying every AC, with the tests named in the test plan present and passing; a PR opened for the branch and recorded on the `pr` field.
- **Good:** The failing test is written first and watched fail for the right reason; the change is the minimum that satisfies the AC.
- **Bad:** Scope beyond the AC, changes that leave the main checkout dirty, or a self-report of success with no test or command output behind it.

### `validation`

A `fresh` agent independently verifies the deliverable against the ideation AC, reproducing each `Verified by:` clause rather than trusting the implementation's self-report. The validator checks what was produced; it does not produce it. Either gate-approval to `done` or rejection back to `implementation` with concrete fixes.

- **Inputs:** The ideation AC list and the implementation's branch and PR — read without the implementation agent's context.
- **Outputs:** A per-AC verdict with the reproduced evidence (command run, exit code, output excerpt) recorded on the task body, and either an approval recommendation or a concrete fix list routed back to `implementation`.
- **Good:** Every `Verified by:` clause is re-run rather than re-read; a rejection names the specific AC that failed and what evidence would clear it.
- **Bad:** Writing or fixing code instead of checking it; accepting the implementation's summary as evidence; a rejection with no reproducible failing check.

### `done`

Terminal state: the task's PR is merged (tracked via the `pr` field and the `pr-merge` mod), `completed` set, `verdict: PASSED`, entity archived. Reached via real merge, not a manual flag flip.

- **Inputs:** The validation verdict and the merged PR.
- **Outputs:** `completed` timestamp set, `verdict: PASSED`, `worktree` cleared, task archived.
- **Good:** The terminal state is reached because the PR actually merged, so the site's main branch contains the change.
- **Bad:** Flipping `status: done` by hand while the PR is still open or closed-unmerged.

## Workflow State

Workflow state is read by the first officer at boot. To view current state, dispatch the first officer or run it directly:

```
spacedock claude
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
