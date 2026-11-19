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
package (`<topic>/dispatch-sprint-execution.md`) the Commander runs from a cold boot,
not a context transfer.

## Sprint lifecycle checklist

The per-sprint operating sequence, owner-tagged. Copy into a sprint's `index.md` to track
progress. The two **independent reviews** (preflight + post-sprint) are the steps most
easily skipped — they are **not optional** and **never self-reviews**: the value is
refuting the FO's own assumptions, so a *fresh* reviewer runs them, never the FO and never
the ideation ensigns (they would grade their own work).

**Shape — Shaping FO**
- [ ] **Scope-lock** with the captain — which entities are in, which defer; bind the topic to a target release train HERE, as one movable line in `index.md` — never in the folder name, never in member labels *(captain decides)*
- [ ] **Carve** — stamp `sprint` / `group` / `sprint-readiness` on members; write `index.md` (goal, members-as-query, DoD, out-of-scope)
- [ ] **Ideate** each gated member — problem / approach / AC + test-plan, with the **riskiest mechanism exercised first** (a spike, or a recorded "no spike needed"); check existing ideation state first — never re-ideate a banked design
- [ ] **⚠️ Preflight staff review (sprint-wide)** — dispatch ONE *independent* reviewer (not the FO, not the ideation ensigns) to refute the **sprint as a whole, not individual tasks**: DoD coverage (every DoD bullet owned by an in-scope member), sequencing / dependency order, cross-member composition & wiring (shared-region collisions, seams), blast-radius across the set, scope (missing / over-scoped members), and Commander cold-boot readiness → `staff-review.md`. Per-task design soundness / AC quality / proof-gaps / riskiest-mechanism-first is the **ideation gate's** job — surface a per-task gap here only when it threatens the sprint's deliverable. Fold Material findings *before* the gates lock. **Implementability lens (added 2026-07-22 after the recorder's undefined-subcommand gap):** walk the first day of each member's implementation and enumerate the decisions the design does not answer — an interface named in the expected surface without a defined shape (a CLI verb with no usage/flags/exit codes, a schema field with no domain, a script with no contract) is a Material finding; refutation-style review audits only what is written, so this lens exists to catch what is absent
- [ ] **Present ideation gates** — checklist accounting + AC cross-check per member; never self-approve *(captain decides)*
- [ ] **Package** — write `dispatch-sprint-execution.md` (boot recipe, per-member build notes, in-drive gates, release-cut recipe)

**Drive — Commander (a separate, cold-booted session)**
- [ ] **Implementation → validation → done** per member; **detached adversarial audit at validation** for every high-stakes surface (front-door, status guards, shipped scaffolding, CI / release machinery)
- [ ] **Merge** each to `main` (PR-merge); keep state commits concurrency-safe
- [ ] **⚠️ Pre-cut antipattern audit** — with all members merged to `main` and the tag **not yet fired**, dispatch an *independent* reviewer (staff-eng persona; not the Commander, not the implementers) over the assembled sprint to catch cross-cutting antipatterns + integration holes **before they ship**. Ship-blockers are fixed before the cut; non-blockers are recorded for the next sprint. The whole value is being *before* the tag — run it after, and the antipattern has already shipped.
- [ ] **Cut the release** — `go test ./...` green from the root, then follow [`docs/releasing.md`](../releasing.md) (the authoritative cut procedure: manifest bumps, the `vN.N.N` tag, what the tag push fires, `next` publishing) *(captain authorizes)*

**Close — Shaping FO**
- [ ] **Seed the next sprint** — fold the pre-cut audit's recorded (deferred / non-blocking) findings into the next sprint's backlog (e.g. `kb` came straight out of the 019x sprint's antipattern findings) + a light post-cut release verification (some release-machinery issues only manifest when the tag actually fires).

## Sprint sequence

| # | Sprint | Deliverable (the value it unlocks) | State |
|---|--------|------------------------------------|-------|
| 1 | [`019x-pre-flip-cleanups`](019x-pre-flip-cleanups/) | spacedock **0.19.7** cut on `next` with the pre-flip cleanups landed | ✅ delivered — 0.19.7 cut (#322–326) |
| 2 | [`0198-pre-flip-hardening`](0198-pre-flip-hardening/) | spacedock **0.19.8** — binary/version/install UX (`qa`), Codex auto-install (`z9`), test-hygiene (`kb`), survey correctness (`vh`) | driving |
| 3 | [`0199-pre-flip-mechanics`](0199-pre-flip-mechanics/) | spacedock **0.19.9** — Linux binaries (`v3`) + dev-tooling quality (`th`, `jm`) | ✅ delivered — 0.19.9 cut + released (#332–336; darwin+linux) |
| 4 | [`0200-flip`](0200-flip/) (capstone) | **0.20.0** cut on `main` + marketplace flipped | driving — pre-flip set (`nzb`+`k6d`+`cmx`) [packaged](0200-flip/dispatch-sprint-execution.md), Commander driving (`nzb` in validation); `pj` the flip held for here |

Post-flip candidates (not yet sprints): a published docs site (`wv` + `e6`),
distribution (`5w` notarize + `44` bundle), and the
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

**The `sprint:` label is the topic slug — the sprint's identity. The release train is
not part of it.** A train rebind (a cut date moves, a sprint slips to the next minor)
changes the one target-train line in the sprint's `index.md` and nothing else — never
member frontmatter, never the folder, never the label. Reshuffling a member between
topics is one `--set`; moving a promised outcome between releases is a captain decision
recorded where the DoD line moves (captain ruling, 2026-07-21).

## Where things live

- **Strategy + sprint sequence:** this file.
- **Each sprint:** `<topic>/` — the topic slug IS the sprint's identity; its target
  release train is one movable line in `index.md`, bound at scope-lock and rebindable
  by captain decision without re-carve. Contents: `index.md` (goal / DoD / deliverable /
  members-as-query), `staff-review.md` (the readiness-review mechanism + findings), and
  `dispatch-sprint-execution.md` (the cold-boot Commander package). Older sprints used
  `NNN-<slug>` paths; their numbers are historical, not the convention.
- **Executable work:** Spacedock entities in `docs/dev/.spacedock-state/`, carrying
  `sprint` / `group` / `sprint-readiness` frontmatter. The build/factory workflow is
  `docs/dev`.
- **Original build roadmap:** `bootstrap-roadmap.md` (stage-specific required tests) —
  the strategy doc that predates this sprint layer.
