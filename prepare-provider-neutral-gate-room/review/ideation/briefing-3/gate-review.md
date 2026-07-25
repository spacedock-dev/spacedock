# Ideation gate: mechanical provider-neutral gate preparation

## Capability

An agent supplies human judgment and committed file choices; one Spacedock command
constructs and binds the canonical open gate room. The agent never authors
`briefing.json`, `request.json`, ids, digests, Git locators, or provider output paths.

## Agent-facing operation

```text
spacedock gate prepare ENTITY \
  --question TEXT \
  --artifact REVIEW.md \
  --summary TEXT \
  [--reference FILE ...] \
  [--workflow-dir DIR]
```

The First Officer first performs one fresh `spacedock gate --help` capability check
for `prepare`, the existing lifecycle operations, and every prepare flag it will use.
A stale launcher halts before source commits, preparation, or state mutation.

On success the command prints exactly:

```text
room=/clean/absolute/entity-room/path
briefing=briefing:<derived-id>
digest=sha256:<canonical-briefing-digest>
state=open
```

The emitted room is the sole public authority handed to a selected provider:

```text
/subspace:r gate <room>
```

The agent does not probe the provider, choose a probe-driven fallback, parse
`request.json`, or reconstruct entity, workflow, Briefing, actor, approver,
materialization, manifest, or output coordinates.

## What preparation writes

Both folder and flat entities converge on:

```text
<state-root>/<slug>/review/<stage>/briefing-N/
```

Immediately after preparation, the room contains exactly two authoritative regular
files:

```text
gate-briefing.json
request.json
```

`request.json` freezes gate, attempt, actor, approver, the arbitrary canonical
Briefing locator, Briefing id, and full canonical SHA-256. The Briefing contains the
exact caller-authored primary Artifact summary and Git-root identities for every
selected file:

```text
git-root://<main|state>/<full-commit>/<repository-path>
sha256:<raw-file-bytes>
```

Selected source bytes are not copied. A later provider may add its owned `provider/`
evidence subtree; it cannot replace or rewrite either authoritative file.
`association.json` is never written because the association is derived and verified
from the frozen request, Briefing, Result, and presented inventory.

## Durability and trust boundaries

- Flat `<slug>.md` and companion `<slug>/` are committed, archived, rolled back, and
  restored as one exact unit without sweeping siblings.
- Git-root resolution reopens only existing local objects. It accepts resolvable
  detached commits, requires the raw SHA to match, and never fetches, deepens,
  hydrates, substitutes worktree bytes, or imposes a branch/ref-retention policy.
- Request, Briefing, Result, and inventory reject duplicate JSON member names before
  typed decoding or mutation.
- Mandatory exact summary applies only to mechanically prepared, request-backed
  Briefings. Existing request-less and advisory Briefings remain summary-optional and
  byte-unchanged.
- File/LOC tables are planning evidence only. Implementation is judged by semantic
  ownership, authority, and acceptance behavior.

## Validation capability

The design reuses the existing CLI, recorder, eligibility, state-commit/archive,
contract, and recorded-gate lifecycle fixtures. It adds folder/flat preparation,
moved-root/detached-object, arbitrary-Briefing-name, duplicate-member, exact-summary,
flat rollback, stale-help, and forbidden-provider-probe controls. Full, race, docs,
and Roborev gates remain required. Provider materialization and the real moved-root
Subspace E2E remain the separately approved `rq` dependency.

## Review disposition

Independent review first required removal of numeric implementation authority and
alignment of the shipped skill owner map. Cycle 11 made those corrections. The
confirming review returned APPROVE with no material or polish findings.

## Recommendation

Approve Cycle 11 as the implementation contract and supersede the stale frozen-copy
attempt-2 approval.

## Decision

Approve to implement this exact preparation boundary; revise to change room,
authority, or owner semantics; or hold before implementation.
