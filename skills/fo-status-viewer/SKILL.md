---
name: fo-status-viewer
description: "First-officer status query/display surface — the `status` command read flag docs, canonical captain-facing invocations, the Captain-Facing State Display rendering, and the GitHub-issue-filing approval gate. Invoke at the first ad-hoc status question, `--next-id`/`--resolve` lookup, or issue filing."
user-invocable: false
---

# First Officer Status Viewer

## Status Viewer

The `${SPACEDOCK_BIN:-spacedock} status` launcher owns path resolution and mutation guards; skill instructions stay declarative and never reference a plugin-private script path.

Invoke it as:
```
${SPACEDOCK_BIN:-spacedock} status --workflow-dir {workflow_dir} [--next-id|--next|--archived|--where ...|--boot|--validate|--resolve REF]
```

- `--boot` — startup roll-up (mods, ID style, next-ID candidate, orphans, PR state, dispatchables). Incompatible with `--next`, `--next-id`, `--archived`, `--where`.
- `--validate` — run before trusting manually edited workflow state.
- `--resolve REF` — deterministic lookup by slug, exact stored ID, or sd-b32 address prefix; `--root` rejects unqualified cross-workflow ambiguity rather than guessing.
- `--next-id` — preview the next-id candidate for `sequential` and `sd-b32` (n/a for `slug`). For `sd-b32`, pass `--id-seed "{slug-or-title}"` and optionally `--id-actor "{actor-or-agent}"` so creation context enters the candidate. To file a new entity, do NOT pair `--next-id` with a hand-written file — use `spacedock new` (see `Skill(skill="spacedock:fo-write-core")`), which mints the id and atomically writes the stamped entity in one call. `--next-id` is candidate-preview only.
- `--next` / `--where "pr !="` — targeted event-loop queries.

State-changing `status --set` and `status --apply-gate` invocations are write authority; use `Skill(skill="spacedock:fo-write-core")` before applying them.

### Captain-Facing State Display

The commissioned README directs the captain to dispatch the FO to inspect workflow state. Invoke `status` for captain-facing display on questions like:

- "what's the workflow state?" / "show me the workflow" / "what's going on?"
- "what's dispatchable?" / "what's ready?" / "what's next?"
- "what's archived?" / "show me the done entities"
- any ad-hoc question a `status` view answers (a single entity, entities in a stage, PR-pending).

**Canonical invocations** (all start with `${SPACEDOCK_BIN:-spacedock} status --workflow-dir {workflow_dir}`):
- Overview: no extra flags.
- Dispatchables: `--next`.
- Archive view: `--archived`.
- Single-entity: `--resolve {ref}` then `--where slug={resolved-slug}`.

**Output rendering guidance.** Forward `status` stdout verbatim inside a fenced code block, with a one-line preface naming the request ("Workflow overview:", "Dispatchable entities:", "Archived entities:"). On empty results, render a literal note ("No dispatchable entities right now.") instead of an empty fence. Do not paraphrase rows, omit columns, invent fields, summarize counts, or editorialize.

## Issue Filing

Do not file GitHub issues without explicit human approval.
