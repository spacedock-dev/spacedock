# Split-root state

A busy workflow generates constant small state commits. Split-root keeps them off your code branch: the README and stage declarations stay on your main branch, while the mutable entity state (frontmatter updates, stage reports, archive moves) lives in a separate state checkout. Your code history stays clean, and several agents or operators can drive the same workflow at once.

The README [opts in with one `state:` field](../concepts/workflows-and-entities.md#keep-workflow-state-off-your-code-branch); the split is transparent, and you read the workflow exactly as you would a single-root one. On a fresh clone the state checkout is absent; run `spacedock state init` to restore it. The shipped [`docs/dev` workflow](https://github.com/spacedock-dev/spacedock/tree/main/docs/dev) runs split-root, a live example.

## Concurrent writers

The agents follow the commit and sync discipline that keeps concurrent writers from clobbering each other; it is theirs to follow, not yours. The one case that reaches you: conflicting edits to the same entity halt for your call rather than auto-resolving.

## Bridging an external tracker

Split-root state is the integration point for an external tracker (Linear, GitHub Issues, or another ticket ledger). The external system can own backlog intake, discussion, and assignment while Spacedock stays the execution workflow. The bridge is two flat frontmatter fields (the [frontmatter contract](../reference/frontmatter-contract.md) explains why they stay flat):

```yaml
issue: ENG-123
source: linear
```

`issue` is the human-facing external reference; `source` records where the entity came from. Keep the tracker out of Spacedock's stage semantics: a bridge should sync through entity creation, state changes, and stage reports, not by adding tracker-specific stage rules.
