# Sprints & roadmap

The roadmap is the strategy layer above the per-entity workflow. It owns outcome, scope, sequencing, and definition-of-done. It does not track task state — that lives in the workflow's state checkout and is queried with `spacedock status`.

## The two roles

- **The shaping first officer** owns strategy and shape: the roadmap, sprint definition (deliverable plus definition-of-done), the gating ideation of each sprint's entities, the cross-entity coherence check, the staff readiness review, and packaging each sprint as a dispatch. It stays high-level and does not hand-drive stage execution.
- **The Commander** takes one packaged sprint and drives it to its deliverable: dispatches each stage, approves execution gates and merges with judgment, runs the sprint-wide integration test, and produces the report. It boots the first-officer skill and creates its own team. It escalates only on a third feedback cycle, a budget blowout, an irrecoverable block, or a genuine scope fork.

The handoff between them is the **dispatch package** — a self-contained sprint package the Commander runs from a cold boot, not a context transfer.

## The sprint lifecycle

A sprint runs in three phases:

- **Shape** (shaping first officer) — scope-lock with the captain, carve the members, ideate each gated member with the riskiest mechanism exercised first, run an independent preflight staff review, present the ideation gates, and package the dispatch.
- **Drive** (Commander, a separate cold-booted session) — implementation → validation → done per member, with a detached adversarial audit at validation for every high-stakes surface; merge each member; run an independent pre-cut antipattern audit before the release tag fires; then cut the release.
- **Close** (shaping first officer) — fold the audit's deferred findings into the next sprint's backlog and run a light post-cut verification.

The two independent reviews — preflight and pre-cut — are never self-reviews. Their value is refuting the first officer's own assumptions, so a fresh reviewer runs them.

## Roadmap as a strategy layer

The roadmap holds the value-ordered sprint sequence: each sprint and the deliverable it unlocks. It is the durable strategy document, distinct from the executable work it sequences.

## A convention, not new code

Today this is a **convention-only dry run** of the sprint/roadmap construct: prose plus frontmatter plus the native `status --where` query, with no new binary code. A sprint *groups* entities by a frontmatter query rather than a hard-coded list:

```bash
# every member of a sprint
spacedock status --workflow-dir docs/dev --where sprint=019x-pre-flip-cleanups
# the drivable set (excludes deferred members)
spacedock status --workflow-dir docs/dev --where sprint=019x-pre-flip-cleanups --where 'sprint-readiness != defer'
```

Members carry `sprint`, `group`, and `sprint-readiness` frontmatter. There is no contract bump and no dedicated `--sprint` gate — the rollup is the native `--where` query. Whether this construct graduates into binary support depends on whether the dry run proves it.
