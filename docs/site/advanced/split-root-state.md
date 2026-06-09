# Multi-workflow & split-root state

Split-root state separates a workflow's *definition* from its *runtime state*. The README stays in the main repo; the mutable entities live in a per-workflow state checkout. This keeps noisy state transitions out of the code branch.

## Two roots

A split-root workflow composes two directories:

- `definition_dir` — the directory containing `README.md` (the workflow definition: stages, schema, gates). It stays in the main repo.
- `state_dir` — `definition_dir/<state>`, declared by the README's `state:` field (commonly `.spacedock-state`). It holds the active entities, archived entities under `_archive`, and stage reports under folder-form entities.

The launcher reads workflow identity and stage declarations from the definition README, and active entities from the state checkout. It writes frontmatter updates and archive moves to the state checkout — never to the code branch.

## Declaring split-root

The README opts in with one top-level frontmatter field:

```yaml
state: .spacedock-state
```

The path resolves relative to the README's directory. If `state` is absent, the workflow uses the same-directory layout (entities live beside the README).

## Concurrency-safe state commits

The state checkout is a single shared git index. When multiple agents write entities concurrently, a bare `git add -A` would sweep up a sibling writer's staged entity. State commits are therefore path-scoped — each writer adds and commits only its own entity path — and pushed with a rebase-on-rejection retry so disjoint single-file commits replay cleanly atop one another.

For the full storage-behavior contract, see the [external-tracker bridge](external-tracker.md) and the development workflow's state model.
