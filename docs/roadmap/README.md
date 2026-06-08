# Spacedock v1 — Roadmap

Durable strategy and the value-ordered sprint sequence for the Go launcher workflow.
This is the **strategy layer above** the per-entity workflow in `docs/dev` — it owns
outcome, scope, sequencing, and definition-of-done. It does **not** track task state;
that lives in `docs/dev/.spacedock-state`, queried by `spacedock status`.

> **Convention experiment (2026-06-08).** This roadmap formalizes the informal sprint
> tracking spacedock-v1 already improvises in
> `docs/dev/.spacedock-state/_sprint-notes.md`. It is the convention-only dry run of the
> sprint/roadmap construct (`docs/dev/_proposals/sprint-roadmap-construct.md`): prose +
> frontmatter + the native `status --where` query, **no new binary code**. See that
> proposal for the decision record and the graduation triggers.

## The two roles

- **Shaping FO** — owns strategy and shape: the roadmap, sprint *definition*
  (deliverable + DoD), the gating ideation of each sprint's entities, the cross-entity
  coherence check, the staff readiness review, and packaging each sprint as a Commander
  dispatch. Stays high-level; does **not** hand-drive stage execution.
- **Commander** — takes ONE packaged sprint and drives it to its deliverable: dispatches
  each stage, approves execution gates and merges with good judgment, runs the
  sprint-wide integration test, produces the report. Boots `spacedock:first-officer` and
  creates its own team. Escalates only on a 3rd feedback cycle, a budget blowout, an
  irrecoverable block, or a genuine scope fork.

The handoff between them is the **conn-to-drive dispatch** — a self-contained sprint
package (`NNN-<slug>/dispatch-sprint-execution.md`) the Commander runs from a cold boot,
not a context transfer.

## Sprint sequence

| # | Sprint | Deliverable (the value it unlocks) | State |
|---|--------|------------------------------------|-------|
| 1 | [`019x-pre-flip-cleanups`](019x-pre-flip-cleanups/) | spacedock **0.19.7** cut on `next` with the pre-flip cleanups landed | ✅ delivered — 0.19.7 cut (#322–326) |
| 2 | `0200-flip` (capstone) | **0.20.0** cut on `main` + marketplace flipped | backlog (shaped, not driven here) |

Post-flip candidates (not yet sprints): a published docs site (`wv` + `e6`),
distribution everywhere (`v3` linux + `5w` notarize + `44` bundle), and the
`spacedock:roadmap` skill graduation — only if this dry run proves the construct.

## How membership works

A sprint GROUPS entities by frontmatter query, never a hard-coded list:

```bash
# every member of a sprint
spacedock status --workflow-dir docs/dev --where sprint=019x-pre-flip-cleanups
# the drivable set (excludes deferred members)
spacedock status --workflow-dir docs/dev --where sprint=019x-pre-flip-cleanups --where 'sprint-readiness != defer'
```

Members carry `sprint:`, `group:`, and `sprint-readiness:` frontmatter. No contract
bump, no `sprint` recognizer, no `--sprint-validate` gate — the rollup is the native
`--where` query (verified on `spacedock 0.19.x` during this dry run).

## Where things live

- **Strategy + sprint sequence:** this file.
- **Each sprint:** `NNN-<slug>/` — `index.md` (goal / DoD / deliverable / members),
  `staff-review.md` (the readiness-review mechanism + findings), and
  `dispatch-sprint-execution.md` (the cold-boot Commander package).
- **Executable work:** Spacedock entities in `docs/dev/.spacedock-state/`, carrying
  `sprint` / `group` / `sprint-readiness` frontmatter. The build/factory workflow is
  `docs/dev`.
- **Original build roadmap:** `bootstrap-roadmap.md` (stage-specific required tests) —
  the strategy doc that predates this sprint layer.
