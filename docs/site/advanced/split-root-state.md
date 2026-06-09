# Multi-workflow & split-root state

A split-root workflow separates the workflow's definition from its runtime
state. The README and stage declarations stay on your main branch; the mutable
entities — frontmatter updates, stage reports, archive moves — live in a
separate state checkout. State transitions stop polluting your code branch's
history, and the same workflow definition can drive shared issues without each
status change landing as a commit on `main`.

You opt in with a single README field. Without it, Spacedock keeps the
default single-root behavior: entities sit beside the README on the same
branch.

## The two roots

A workflow resolves to two directory roles, derived from the README's `state:`
field (`internal/status/roots.go`):

- **`definition_dir`** — the directory containing `README.md`. It holds the
  workflow identity and the stage declarations. This is what you pass as
  `--workflow-dir`.
- **`state_dir`** — `definition_dir/<state>` when the README declares a
  non-empty `state:` value, otherwise `definition_dir` itself. It holds the
  active entities and the `_archive` directory.

`spacedock status` reads stage declarations from `definition_dir/README.md` and
entities from `state_dir`. It writes frontmatter updates and archive moves only
into `state_dir`. In single-root mode the two roots are the same path, matching
the original same-directory layout, so existing workflows are unaffected.

## Declare `state:` in the README

Add a top-level `state:` field to the README frontmatter. The value is a path
relative to the README directory:

```yaml
state: .spacedock-state
```

The path is resolved against the definition dir. The interpreter
(`internal/status/state.go`) rejects two classes of value rather than following
them silently:

- An **absolute path** fails: `state:` must be relative to the README directory.
- A path that **escapes the definition dir** via `..` fails: the v0 contract is
  a child checkout, not an arbitrary location.

An empty `state:`, an absent field, or the explicit `$inline` sentinel all
resolve to single-root (inline) mode. The shipped `docs/dev` workflow uses
`state: .spacedock-state`; see its README for a live example.

Active entities live directly under `state_dir` — there is no `entities/`
subdirectory. Archived entities move to `state_dir/_archive`. Read the state
with the launcher exactly as you would a single-root workflow; the split is
transparent to the command surface:

```bash
spacedock status --workflow-dir docs/dev
spacedock status --workflow-dir docs/dev --next
```

## The state branch

The state checkout lives on an orphan branch in the same repo — no second repo,
no second remote — and the checkout itself is a linked worktree of the main repo
at the gitignored `state:` path. State commits land on the orphan branch, so the
code branch never sees them. Spacedock derives the branch name from the workflow
dir's basename: `spacedock-state/<basename>` — so `docs/dev` maps to
`spacedock-state/dev`. An explicit `state-branch:` field in the README overrides
the derived name verbatim (`StateBranch` in `internal/status/state.go`).

Because the branch is shared through `origin`, multiple agents — and multiple
operators — can drive the same workflow concurrently. That makes the commit
discipline below a correctness requirement, not a style preference.

## Concurrency-safe state commits

The state checkout is a single, non-branched git index. A bare `git add -A`
followed by a bare `git commit` sweeps up a sibling writer's staged entity,
cross-attributing or clobbering it. **Every writer commits path-scoped**, naming
exactly the entity it touched:

```bash
git -C {state_checkout} add {entity_path}
git -C {state_checkout} commit -m "..." -- {entity_path}
```

Never a bare `git add -A` or a bare `git commit` against the state checkout.
On `index.lock` contention, retry after roughly two seconds. When the status
tool owns the `add`+`commit` under a lock, route through it instead — a
tool-managed atomic commit is preferred over the manual path-scoped fallback.

### Multi-writer sync

The path-scoped rule extends to three sync points against `origin` — not a pull
before every dispatch:

- **After a state commit, push.** `git -C {state_checkout} push origin {state_branch}`.
- **On a non-fast-forward rejection, rebase then re-push.**
  `git -C {state_checkout} pull --rebase origin {state_branch}` replays your
  single-file commit atop the peer's. Disjoint paths produce no conflict.
- **At first-officer boot, pull once.** Integrate peers' state at boot, not on
  every read.

If `pull --rebase` conflicts — two writers editing the same entity's frontmatter
at once — the first officer halts the dispatch, aborts the rebase, and surfaces
the conflicting entity and peer commit to the captain. It does not
`--force`-push and does not auto-resolve with `-X ours`/`-X theirs`, either of
which silently drops a peer's edit. A full lock model is out of scope; the halt
is the boundary behavior.

## Worktree stages under split-root

When a split-root workflow has a worktree stage, the worktree isolates the
deliverable work product only. The entity body and stage reports are still
written and committed to the state checkout at the entity's state-checkout path,
never a worktree copy — the dispatch helper hands the worker that path even under
a worktree stage. "Commits must be on this branch" applies to the deliverable
artifacts; entity state always lands in the state checkout.

## Bridging an external tracker

Split-root state is the integration point for external trackers — Linear, GitHub
Issues, kata, or another ticket ledger. The external system can own backlog
intake, discussion, and assignment while Spacedock remains the execution
workflow. The bridge uses flat top-level frontmatter fields (`issue:`,
`source:`) so the current line-oriented parser preserves them. See the
[external-tracker bridge](external-tracker.md) for the field contract and the
principles that keep Spacedock's stage semantics out of the tracker.
