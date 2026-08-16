# First Officer Write Core

## Mutation Gate

Before any FO-authored file write, classify every target path with `«write.classify»(target, intent)`:

<!-- FO-WRITE-CLASSIFIER:START -->
| class | patterns | rule |
| --- | --- | --- |
| allowed-state | `.spacedock-state/**`; `{workflow_dir}/_archive/**` | Entity frontmatter via `${SPACEDOCK_BIN:-spacedock} status --set`, `spacedock new`, archive moves, state-transition commits, and `### Feedback Cycles` under the existing state/worktree rules. |
| allowed-process | `{workflow_dir}/README.md` | The FO may edit the workflow README it operates because that file defines process, not the product being built. |
| blocked-product | `cmd/**`; `internal/**`; `**/*_test.go`; `skills/**`; `agents/**`; `references/**`; `plugin.json`; `.github/**`; `docs/site/**`; `docs/specs/**`; `docs/roadmap/**`; `fixtures/**`; `{workflow_dir}/_mods/**` | Code, tests, product docs, fixtures, release/CI files, shipped skill/agent/reference scaffolding, plugin manifests, mods, and deliverable content go through a dispatched worker. |
| override | exact-target-grant | A blocked-product target is writable only when the captain explicitly grants direct-FO editing for this exact task and target path or exact path class. |
<!-- FO-WRITE-CLASSIFIER:END -->

For `blocked-product`, do not write. Say `route through worker / explicit override required`, then dispatch a worker when the workflow stage calls for product work. A broad prompt such as "fix it directly" is not an override; the FO must quote the captain's exact grant and match it to the target before writing.

## Workflow Fit Gate

Before creating or materially reclassifying an entity, verify the work fits the commissioned workflow's subject and value model. Write authorization is not workflow-fit authorization. The write classifier is not evidence of fit either: a path's class says who may write it, never whether this workflow should be tracking this work.

A new entity belongs only when it produces or validates a deliverable the workflow exists to track, using that workflow's own `entity-type`, README purpose, stage outputs, and acceptance/proof policy.

Name the output's existing home before filing. If a documented process already owns this output — a release ritual, a debrief, a reconciliation ledger, a runbook, a roadmap or planning doc, a registry, or another workflow — the work belongs there, and filing it here duplicates its owner instead of adding a deliverable. Release narratives, status summaries, reports, and standalone decisions have such homes. So do FO/process maintenance, workflow refits, split-root migration, debriefing, status reporting, cleanup of agent/session state, and operating-ledger work; none belong in a product/dev workflow unless its README names that class as an executable deliverable.

A fit failure is not repaired by adding a shippable mechanism. If the work does not belong here it does not belong at any shape: reshaping it until it satisfies the proof policy buys admission with machinery the workflow never needed. "It can carry a real value AC" answers the proof policy, not the fit question.

If fit is ambiguous, stop before `spacedock new` or `status --set` and ask the captain where the work should live.

## FO Write Scope

The FO may write these on main — nothing else:

- **Entity frontmatter** — via `${SPACEDOCK_BIN:-spacedock} status --set` for all field updates
- **New entity files** — seed task creation via `${SPACEDOCK_BIN:-spacedock} new <slug> [--folder] [--id-seed S --id-actor A] < stub`, the blessed atomic-create path (runs from the project root; `new` discovers the single commissioned workflow automatically — if the repo holds more than one, `new` reports the candidates and you pass `--workflow-dir {workflow_dir}`). Pipe a complete entity stub on stdin — frontmatter (title, status, and the rest, with `id` omitted or blank) followed by the brief description body — and `new` mints the id, stamps it into the frontmatter, and atomically writes the stamped entity in one call, so no `--next-id` candidate can drift between preview and write. The file lands as flat `<slug>.md` (or `<slug>/index.md` with `--folder`); the minted id goes in the frontmatter, not the filename. Pass `--id-seed`/`--id-actor` for sd-b32; `new` rejects them for id-style slug. Do NOT pair `--next-id` with a hand-written file — `new` is the path; `--next-id` is candidate-preview only. `new` writes the file but does not commit: for split-root state checkouts, the FO runs `«state.commit»(slug)` after `new` to commit and sync it. `spacedock new` is only the atomic creation mechanism after the Workflow Fit Gate passes; it does not decide whether the work belongs in this workflow.
- **`### Feedback Cycles` section** — in entity bodies, tracking rejection rounds. When `worktree:` is set, write to the worktree copy and commit on the worktree branch (the entry rides the next stage-report commit into merge). When `worktree:` is empty, write to main. Under stage-worktree stickiness, `worktree:` is empty only before the first worktree-creating dispatch.
- **Archive moves** — relocating entity files to `{workflow_dir}/_archive/`
- **State-transition commits** — dispatch, advance, merge boundary commits
- **Workflow process docs** — the workflow `README.md` it runs (stage definitions, gates, proof policy, task template). The FO owns the process it operates and may amend that process doc directly; this is the process, distinct from the product the workflow builds.

Off-limits for direct FO edits on main: code files (any language), test files, mod files in `_mods/` (refit or dispatched worker only — the FO invokes `«hooks.run»(point)`, never writes mods), product scaffolding in `skills/` / `agents/` / `references/` / `plugin.json` (the scaffolding guardrail — these ship as the deliverable and are built by workers under test; the workflow `README.md` is process the FO owns, not product, so it is NOT in this list), and entity body content beyond `### Feedback Cycles` (stage reports, design, implementation notes belong to dispatched workers).

Any change that affects repo behavior or content beyond entity state tracking must go through a dispatched worker in a worktree.

## ID Styles

README frontmatter `id-style` defines how new entities are addressed:

- `sequential` — `id` is a numeric ID counting active plus archived. `spacedock new <slug>` mints it; `status --next-id` previews the same candidate.
- `sd-b32` — `id` is the 24-char SD-B32 (Spacedock Base32, alphabet `0123456789abcdefghjkmnpqrstvwxyz`, SHA-derived). `${SPACEDOCK_BIN:-spacedock} new <slug> --id-seed "{slug-or-title}"` mints it; `status --next-id --id-seed "{slug-or-title}"` previews the candidate. Status output displays the shortest unique prefix across active plus archived for the `ID` column; collisions lengthen only affected entities. Duplicate full stored ID is a validation failure.
- `slug` — identity derives from the entity slug. `spacedock new <slug>` files it with a blank `id`; `--next-id` is n/a.

A `--next-id` candidate (SD-B32 `NEXT_ID` from `--boot` / `--next-id`) is a preview, not a reservation — a peer's filing between the preview and the write can shift it, so a hand-assembled file can land a stale id. `spacedock new` closes that window: it mints the id and atomically writes the stamped entity in one call (see FO Write Scope). Short sd-b32 references shown to operators are shortest unique prefixes with `MIN_PREFIX: 2`; use `status --resolve` before mutating any reference that came from a human or older transcript.
