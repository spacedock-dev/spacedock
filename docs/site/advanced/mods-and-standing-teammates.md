# Mods & standing teammates

Mods extend a workflow without touching the binary. A mod is a markdown file under `{workflow_dir}/_mods/`, and it adds behavior that is the same across stages and across workflows: a step that must happen at a fixed point in every run, or a specialist that should persist across the whole session. The first officer reads the mod and acts on it; nothing compiles.

A mod does one of two things, or both.

## Lifecycle hooks

A hook runs first-officer prose at a fixed point in the run: at `startup`, on an `idle` pass when nothing is ready to dispatch, or at the `merge` boundary when an entity terminalizes. The point is workflow-independent: any workflow can register the same hook to get the same behavior.

The canonical example is the [`pr-merge` mod](https://github.com/spacedock-dev/spacedock/blob/next/docs/dev/_mods/pr-merge.md): it opens the code-branch PR at merge, records the PR on the entity, and holds the terminal transition until the PR merges. A merge hook can block that transition, and the binary enforces the block so a half-merged entity cannot slip past the gate.

## Standing teammates

A standing teammate is a long-lived specialist agent declared by a mod. It lives in the team for the session and is addressed by name. Reach for one when the same specialist judgment recurs across entities and is worth a persistent agent rather than a fresh dispatch each time.

The canonical example is the **comm-officer**, a prose-polisher the first officer routes deliberate drafts through (PR bodies, gate summaries, debriefs) before they reach you. Routing is best-effort: if the teammate is absent or slow, the work proceeds without it.

## The exact format

The mod file format, the hook points, and the `spacedock dispatch` subcommands that read standing-teammate mods are defined in the skills and binary that own them. See [the mods reference](https://github.com/spacedock-dev/spacedock/blob/next/docs/dev/README.md) for the authoritative contract.
