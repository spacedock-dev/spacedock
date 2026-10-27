# Mods and standing teammates

A **mod** adds behavior to a workflow without adding a stage. Ask your agent to install an existing mod or create one.

## Lifecycle hooks

A mod can act at startup, while the workflow is idle or before work completes.

For example, the [pr-merge mod](https://github.com/spacedock-dev/spacedock/blob/main/docs/dev/_mods/pr-merge.md) manages pull requests and keeps work open until the PR merges.

## Standing teammates

A standing teammate is a specialist agent available throughout a session for recurring work.

For example, the [comm-officer mod](https://github.com/spacedock-dev/spacedock/blob/main/docs/dev/_mods/comm-officer.md) provides a prose editor for drafts and summaries.
