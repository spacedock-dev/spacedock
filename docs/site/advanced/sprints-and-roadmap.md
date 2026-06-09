# Sprints & roadmap

A roadmap is the strategy layer above the per-entity workflow: it owns outcome, scope, sequencing, and definition-of-done; the workflow owns task state. The two never overlap. Spacedock's own build uses this split. `docs/roadmap/README.md` is the strategy layer, and `docs/dev/.spacedock-state/` holds the executable entities that `spacedock status` queries. The roadmap explicitly "does **not** track task state."

This page describes that construct as Spacedock practices it on itself. It is a convention (prose, frontmatter, and the native `status --where` query), not new binary behavior. There is no `sprint` recognizer, no `--sprint-validate` gate, and no contract bump. If you want this discipline for your own workflow, you adopt the convention; nothing in the launcher enforces it.

## The roadmap as strategy layer

The roadmap file owns four things the workflow does not: the **outcome** each sprint unlocks, the **scope** (which entities are in), the **sequencing** (value-ordered sprint list), and the **definition-of-done** per sprint. Task state, meaning which stage an entity is in and whether its gate passed, stays in the entities and is read with `spacedock status`. Keep the two separate: a roadmap that starts tracking stage transitions has duplicated the workflow and will drift from it.

A roadmap groups its work into sprints. A sprint is a set of entities driven to one deliverable, with its own `index.md` recording the goal, the members-as-query, the DoD, and what is out of scope.

## The shaping-FO / Commander split

A sprint is run by two distinct roles, and the boundary between them is the load-bearing rule of the construct.

- **Shaping FO** owns strategy and shape: the roadmap, each sprint's *definition* (deliverable + DoD), the gating ideation of the sprint's entities, the cross-entity coherence check, the staff readiness review, and packaging the sprint for execution. It stays high-level and does **not** hand-drive stage execution.
- **Commander** takes one packaged sprint and drives it to its deliverable: dispatches each stage, approves execution gates and merges, runs the sprint-wide integration test, and produces the report. The Commander boots `spacedock:first-officer` and creates its own team.

The handoff between them is a **conn-to-drive dispatch**: a self-contained sprint package at `NNN-<slug>/dispatch-sprint-execution.md` that the Commander runs from a cold boot. It is a package, not a context transfer. The Commander does not inherit the shaping FO's session, and escalates back to the shaping FO and captain only on a third feedback cycle, a budget blowout, an irrecoverable block, or a genuine scope fork.

## Sprint lifecycle: shape, drive, close

A sprint moves through three phases, owner-tagged. The per-sprint checklist lives in the sprint's `index.md`; the canonical template is the lifecycle checklist in `docs/roadmap/README.md`.

**Shape (shaping FO).**

1. **Scope-lock** with the captain: which entities are in, which defer. The captain decides.
2. **Carve.** Stamp `sprint`, `group`, and `sprint-readiness` frontmatter on the members; write `index.md` (goal, members-as-query, DoD, out-of-scope).
3. **Ideate** each gated member: problem, approach, acceptance criteria, and test plan, with the **riskiest mechanism exercised first** (a spike, or a recorded "no spike needed"). Check existing ideation state first; never re-ideate a banked design.
4. **Preflight staff review.** Dispatch an *independent* reviewer to refute the designs. This is not optional and never a self-review: the reviewer is neither the FO nor the ideation ensigns, because the value is refuting the FO's own assumptions. Findings land in `staff-review.md` and Material ones are folded before the gates lock.
5. **Present ideation gates**, with checklist accounting and acceptance-criteria cross-check per member. The captain decides; the FO never self-approves.
6. **Package.** Write `dispatch-sprint-execution.md`: the boot recipe, per-member build notes, in-drive gates, and the release-cut recipe.

**Drive (Commander, a separate cold-booted session).**

1. Move each member through implementation, validation, and done, with a **detached adversarial audit at validation** for every high-stakes surface (front door, status guards, shipped scaffolding, CI/release machinery).
2. Merge each member to `next` by PR; keep state commits concurrency-safe.
3. **Pre-cut antipattern audit.** With all members merged and the tag not yet fired, dispatch an *independent* reviewer over the assembled sprint to catch cross-cutting antipatterns and integration holes before they ship. Ship-blockers are fixed before the cut; non-blockers are recorded for the next sprint. Running it after the tag means the antipattern has already shipped.
4. **Cut the release.** Confirm `go test ./...` is green from the root, then follow the authoritative cut procedure (see [Releasing](../contributing/releasing.md)). The captain authorizes the cut.

**Close (shaping FO).**

1. **Seed the next sprint.** Fold the pre-cut audit's deferred and non-blocking findings into the next sprint's backlog, and run a light post-cut release verification, since some release-machinery issues only manifest once the tag actually fires.

## Membership is a query, not a list

A sprint groups its entities by frontmatter query, never a hard-coded roster. Members carry `sprint`, `group`, and `sprint-readiness` frontmatter, and the rollup is the native `--where` filter on `spacedock status`. Run it against the workflow that holds the entities (`docs/dev` in Spacedock's own build):

```bash
# every member of a sprint
spacedock status --workflow-dir docs/dev --where sprint=0200-flip

# the drivable set — excludes deferred members
spacedock status --workflow-dir docs/dev --where sprint=0200-flip --where 'sprint-readiness != defer'
```

Each `--where` clause is `field <op> value`, where `<op>` is `=` or `!=`. Stacking clauses ANDs them: an entity is listed only when it matches every clause. The operator forms also cover presence and absence: `field !=` matches a non-empty value, `field =` matches an empty one. This is the same filter any reader can run; the sprint is a convention layered on top of it, not a separate command.

`sprint-readiness: defer` is how a member stays in the sprint's definition but out of the Commander's drivable set. In the `0200-flip` capstone, `pj` is marked `defer` to mean "driven in the shaping session, not by the cold-boot Commander". That defers it from one driver; it is not dropped from the sprint.

## Adopting it for your own workflow

You do not need any launcher feature to use this. Add a roadmap file above your workflow, stamp `sprint` / `group` / `sprint-readiness` on the entities you group, and read membership with `status --where`. The construct buys you the strategy/state separation and the shaping-FO/Commander discipline; it costs you a convention you maintain by hand. Spacedock runs it on itself precisely to learn whether the convention earns a graduation to first-class support before any code is written for it.
