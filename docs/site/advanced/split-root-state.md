# Multi-workflow & split-root state

A split-root workflow separates the workflow's definition from its runtime state. The README and stage declarations stay on your main branch; the mutable entities (frontmatter updates, stage reports, archive moves) live in a separate state checkout. State transitions stop polluting your code branch's history, and several agents or operators can drive the same workflow at once.

Reach for it when status changes would otherwise land as noisy commits on `main`, or when more than one writer needs to advance the same workflow concurrently. Without it, Spacedock keeps the default: entities sit beside the README on the same branch.

## Opt in with one field

Add a `state:` field to the README frontmatter, a path relative to the README directory:

```yaml
state: .spacedock-state
```

Spacedock then reads stage declarations from the README and entities from the state path, and writes entity changes only into the state path. The split is transparent to the command surface: read the workflow exactly as you would a single-root one.

```bash
spacedock status --workflow-dir docs/dev
```

The shipped `docs/dev` workflow runs split-root; see its README for a live example.

## Concurrent writers

The state checkout is a shared git index that multiple agents commit to, so the commit and sync discipline is a correctness requirement, not a style choice: writers commit path-scoped (never a bare `git add -A`), and conflicting edits to the same entity halt for the captain rather than auto-resolving. The exact protocol is owned by the launcher and the ensign skill; see [the split-root state contract](https://github.com/spacedock-dev/spacedock/blob/next/docs/dev/README.md) for the authoritative rules.

## Bridging an external tracker

Split-root state is the integration point for an external tracker (Linear, GitHub Issues, kata, or another ticket ledger). The external system can own backlog intake, discussion, and assignment while Spacedock stays the execution workflow. The bridge uses flat top-level frontmatter fields so the parser preserves them:

```yaml
issue: ENG-123
source: linear
```

`issue` is the human-facing external reference; `source` records where the entity came from. Keep the tracker out of Spacedock's stage semantics: a bridge should sync through entity creation, state changes, and stage reports, not by adding tracker-specific stage rules.
