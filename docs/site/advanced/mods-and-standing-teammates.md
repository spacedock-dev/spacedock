# Mods & standing teammates

Two needs recur as a workflow matures: a step that must happen at a fixed point in every run (open a PR when work merges), and a specialist whose judgment should persist across the whole session (a prose polisher). A mod adds either: a markdown file in the workflow's `_mods/` directory that the first officer reads and acts on.

## Lifecycle hooks

A hook adds a step the first officer performs at a fixed point in the run: at `startup`, on an `idle` pass when nothing is ready to dispatch, or at the `merge` boundary when an entity reaches its final stage. The point is workflow-independent: any workflow can register the same hook to get the same behavior.

The canonical example is the [`pr-merge` mod](https://github.com/spacedock-dev/spacedock/blob/main/docs/dev/_mods/pr-merge.md): it opens the code-branch PR at merge, records the PR on the entity, and holds the terminal transition until the PR merges. The block is enforced; a half-merged entity cannot slip past the gate.

## Standing teammates

A standing teammate is a long-lived specialist agent declared by a mod. It lives in the team for the session and is addressed by name. Reach for one when the same specialist judgment recurs across entities and is worth a persistent agent rather than a fresh dispatch each time.

The canonical example is the **comm-officer**, a prose-polisher the first officer routes deliberate drafts through (PR bodies, gate summaries, debriefs) before they reach you. Routing is best-effort: if the teammate is absent or slow, the work proceeds without it.

Ask the agent to install a shipped mod or write a new one; the file format is its job.
