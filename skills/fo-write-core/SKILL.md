---
name: fo-write-core
description: "First-officer main-branch write-authority boundary — what the FO may write on main, the `spacedock new` atomic-create procedure, and new-entity id-style minting. Invoke at the first write to main (`status --set`, `spacedock new`, archive move, or `### Feedback Cycles` write)."
user-invocable: false
---

# First Officer Write Core

## FO Write Scope

The FO may write these on main — nothing else:

- **Entity frontmatter** — via `${SPACEDOCK_BIN:-spacedock} status --set` for ordinary field updates, or `${SPACEDOCK_BIN:-spacedock} status --apply-gate --gate {gate_id} --entity {ref} --verdict approve|revise|reject` when applying an externally captured gate decision
- **New entity files** — seed task creation via `${SPACEDOCK_BIN:-spacedock} new <slug> [--folder] [--id-seed S --id-actor A] < stub`, the blessed atomic-create path (runs from the project root; `new` discovers the single commissioned workflow automatically — if the repo holds more than one, `new` reports the candidates and you pass `--workflow-dir {workflow_dir}`). Pipe a complete entity stub on stdin — frontmatter (title, status, and the rest, with `id` omitted or blank) followed by the brief description body — and `new` mints the id, stamps it into the frontmatter, and atomically writes the stamped entity in one call, so no `--next-id` candidate can drift between preview and write. The file lands as flat `<slug>.md` (or `<slug>/index.md` with `--folder`); the minted id goes in the frontmatter, not the filename. Pass `--id-seed`/`--id-actor` for sd-b32; `new` rejects them for id-style slug. Do NOT pair `--next-id` with a hand-written file — `new` is the path; `--next-id` is candidate-preview only. `new` writes the file but does not commit: for split-root state checkouts, the FO runs `«state.commit»(slug)` after `new` to commit and sync it.
- **`### Feedback Cycles` section** — in entity bodies, tracking rejection rounds. When `worktree:` is set, write to the worktree copy and commit on the worktree branch (the entry rides the next stage-report commit into merge). When `worktree:` is empty, write to main. Under stage-worktree stickiness, `worktree:` is empty only before the first worktree-creating dispatch.
- **Archive moves** — relocating entity files to `{workflow_dir}/_archive/`
- **State-transition commits** — dispatch, advance, merge boundary commits
- **Workflow process docs** — the workflow `README.md` it runs (stage definitions, gates, proof policy, task template). The FO owns the process it operates and may amend that process doc directly; this is the process, distinct from the product the workflow builds.

Off-limits for direct FO edits on main: code files (any language), test files, mod files in `_mods/` (refit or dispatched worker only — the FO runs mod hooks, does not write them), product scaffolding in `skills/` / `agents/` / `references/` / `plugin.json` (the scaffolding guardrail — these ship as the deliverable and are built by workers under test; the workflow `README.md` is process the FO owns, not product, so it is NOT in this list), and entity body content beyond `### Feedback Cycles` (stage reports, design, implementation notes belong to dispatched workers).

Any change that affects repo behavior or content beyond entity state tracking must go through a dispatched worker in a worktree.

## ID Styles

README frontmatter `id-style` defines how new entities are addressed:

- `sequential` — `id` is a numeric ID counting active plus archived. `spacedock new <slug>` mints it; `status --next-id` previews the same candidate.
- `sd-b32` — `id` is the 24-char SD-B32 (Spacedock Base32, alphabet `0123456789abcdefghjkmnpqrstvwxyz`, SHA-derived). `${SPACEDOCK_BIN:-spacedock} new <slug> --id-seed "{slug-or-title}"` mints it; `status --next-id --id-seed "{slug-or-title}"` previews the candidate. Status output displays the shortest unique prefix across active plus archived for the `ID` column; collisions lengthen only affected entities. Duplicate full stored ID is a validation failure.
- `slug` — identity derives from the entity slug. `spacedock new <slug>` files it with a blank `id`; `--next-id` is n/a.

A `--next-id` candidate (SD-B32 `NEXT_ID` from `--boot` / `--next-id`) is a preview, not a reservation — a peer's filing between the preview and the write can shift it, so a hand-assembled file can land a stale id. `spacedock new` closes that window: it mints the id and atomically writes the stamped entity in one call (see FO Write Scope). Short sd-b32 references shown to operators are shortest unique prefixes with `MIN_PREFIX: 2`; use `status --resolve` before mutating any reference that came from a human or older transcript.
