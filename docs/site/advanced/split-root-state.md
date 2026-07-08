# Split-root state

A busy workflow generates constant small state commits. Split-root keeps them off your code branch. The README and stage declarations stay on your main branch; the mutable entity state (frontmatter updates, stage reports, archive moves) lives in a separate state checkout. Your code history stays clean, and several agents or operators can drive the same workflow at once.

The README [opts in with one `state:` field](../concepts/workflows-and-entities.md#keep-workflow-state-off-your-code-branch); the split is transparent, and you read the workflow exactly as you would any other. On a fresh clone the state checkout is absent; run `spacedock state init` to restore it. The checkout always lands at the workflow's declared path under the repository's main worktree, wherever the command runs from; a repo with no `origin` remote resumes from the local state branch. The shipped [`docs/dev` workflow](https://github.com/spacedock-dev/spacedock/tree/main/docs/dev) runs split-root, a live example.

## Concurrent writers

The agents follow the commit and sync discipline that keeps concurrent writers from clobbering each other; it is theirs to follow, not yours. The one case that reaches you: conflicting edits to the same entity halt for your call rather than auto-resolving.
