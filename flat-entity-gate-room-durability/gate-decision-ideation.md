# Ideation decision: canonical review-room references

## Recommendation

Approve the narrowed design. Use `@review/...` as the canonical stored locator for
new gate and advisory-round rooms. Resolve it from the ticket's shared review home,
`<state-root>/<slug>/review/...`, for both flat and folder tickets.

## Selected approach

- Keep existing flat `./<slug>/review/...` and folder `./review/...` references
  readable and frozen. Do not migrate tickets or rewrite retained decisions.
- Remove the folder-only `gate record --round` refusal. Publish flat and folder
  rounds below the same physical review home and store `@review/...` for new rooms.
- Route new canonical refs through one small gate-package resolver. Non-`@review/`
  values continue through the existing legacy reader.
- Carry the derived ticket review home into round Artifact validation for both
  publication and replay. An Artifact outside that home, including a sibling ticket
  or the mutable ticket file, must fail without changing bytes.
- Reuse the existing flat state-commit unit. Current code already commits
  `<slug>.md` with its present or tracked `<slug>/` companion tree.

## Risk evidence

Current source already proves the common physical review home, literal flat companion
commit, tracked deletion, and byte-clean symlink ancestry behavior. The remaining
folder-only round guard is directly pinned by an existing test.

Independent review found and closed one in-scope authority gap: a flat round cannot
use the state root as its Artifact trust boundary. The corrected design uses the
derived ticket review home for both forms and both write and replay paths. No material
finding remains.

## Expected surface

- Estimated net change: **+180 LOC across 12 files**.
- Expected movement: approximately **270 insertions and 90 deletions**.
- Tolerance: **±70 net LOC and ±2 files**.
- Production changes are limited to five existing files. The remaining files are
  focused tests and two documentation updates.

Excluded: a new entity-identity package, migration, ref rewriting, dispatch-stamp
parity, general collision/symlink hardening, and archive/opaque-ref validation work.

## Semantic changes

- Command grammar and decision authority do not change.
- New room bindings use reserved `@review/<normalized-path>` syntax.
- Flat advisory rounds become supported and replayable.
- Existing relative refs and `git-root://` Artifact identities remain unchanged.
- State-commit output and ordinary status behavior remain unchanged.

## Proposed proof

1. A two-host Git test proves one supported state commit publishes a flat ticket and
   complete round room, with matching digest, clean replay, and no dirty sibling.
2. Flat and folder gate and round writers produce the same `@review/...` meaning and
   resolve to the same ticket-scoped physical layout.
3. Historical flat, folder, and legacy-round fixtures retain their stored refs,
   resolved bytes, and digests without migration.
4. Flat and folder rounds publish and replay identically; divergent replay is
   byte-clean.
5. Malformed canonical refs, out-of-home Artifacts, sibling tickets, and mutable
   ticket files fail byte-clean, while valid in-home Artifacts pass.

## Decision effect

Approval establishes this estimate and semantic boundary as the implementation
baseline. Rejection returns the design to ideation without creating a worktree.
