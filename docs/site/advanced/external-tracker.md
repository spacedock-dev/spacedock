# Bridge an external tracker

Your backlog may already live somewhere else: Linear, GitHub Issues, another ticket ledger. The bridge is asking, not configuring: tell the first officer to intake an external item ("intake GitHub PR #134") and it files an entity carrying the reference. The external system keeps owning intake, discussion, and assignment; Spacedock stays the execution workflow.

Two flat frontmatter fields carry the reference (the [frontmatter contract](../reference/frontmatter-contract.md) explains why they stay flat):

```yaml
issue: ENG-123
source: linear
```

`issue` is the human-facing external reference; `source` records where the entity came from. Keep the tracker out of Spacedock's stage semantics: sync through entity creation, state changes, and stage reports, not tracker-specific stage rules.
