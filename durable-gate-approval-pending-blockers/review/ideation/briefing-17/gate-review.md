# Gate recorder ideation reset — v1 in place

Decision requested: approve this corrected v1 recorder design for implementation. It keeps
the binary-owned decision record, removes the agent-authored transaction envelope, and does
not add a schema version or migrate any existing gate history.

## Exact command surface

```text
spacedock gate record ENTITY --briefing FILE [--workflow-dir DIR]
spacedock gate record ENTITY --result FILE --association FILE --actor ID [--adoption-note TEXT] [--workflow-dir DIR]
spacedock gate record ENTITY --decision approve|revise|hold --actor ID [--reason TEXT] [--directive TEXT] [--workflow-dir DIR]
spacedock gate validate ENTITY [--workflow-dir DIR]
```

`record` accepts exactly one semantic source: `--briefing`, `--result`, or `--decision`.
`--association` is verification evidence, not another semantic source, and is valid only
with `--result`. `--operation`, expected pointers/digests, and caller-supplied recorder ids
do not exist. `validate` reads and checks; it never writes.

For `--briefing`, the binary derives the operation while holding the entity lock:

- no logical gate: open the first attempt;
- current attempt open + different complete Briefing: **rebind** that attempt by replacing
  its one binding; Git and the review room retain the old binding;
- current attempt closed: **supersede** by appending a new attempt with the new binding.

For either decision source, the binary closes only the current open attempt. A provider
Result supplies the exact portable Resolution; chat supplies decision semantics and
recording identity, while the binary mints id/time/shape. A delegated chat decision also
requires the quoted directive. `revise`/`hold` require a reason or valid included Annotation.

## What the commands write

The examples use fixed ids/times to make the projected bytes testable. In every case,
`status: validation`, `sprint: durable-decisions`, the Markdown body, process state, and
dispatch state remain byte-identical; only `gates` changes.

### 1. Bind a complete Briefing

Before:

```yaml
status: validation
sprint: durable-decisions
```

Command:

```text
spacedock gate record 3k --briefing review/validation/briefing-1/briefing.json --workflow-dir docs/dev
```

After (new v1 writes omit derivable attempt mechanics):

```yaml
status: validation
gates:
  version: 1
  current:
    gate: gate:docs-dev:3k:validation
  records:
    - id: gate:docs-dev:3k:validation
      stage: validation
      attempts:
        - id: gate-attempt:3k-validation-1
          briefing:
            id: briefing:docs-dev:3k:validation:attempt-1:revision-1
            digest: sha256:0a54f1baec0120c1c93523e6900a6ce28e025c570289e5dfa9835e28099042ac
            digest-domain: canonical-bytes
            room-ref: ./review/validation/briefing-1
sprint: durable-decisions
```

Re-running the same form with revision 2 while this attempt is open changes only its
`briefing` mapping (rebind); after it has a Resolution, the same form appends a new minimal
attempt (supersede). The agent never chooses either operation.

### 2. Consume an exact provider Result

Before: the attempt above is open on canonical Briefing revision 1. Command:

```text
spacedock gate record 3k --result ROOM/result.json --association ROOM/result-association.json --actor person:reviewer --adoption-note "Adopted by the captain" --workflow-dir docs/dev
```

After: the binary adds only the verified portable Resolution to that attempt:

```yaml
        - id: gate-attempt:3k-validation-1
          briefing:
            id: briefing:docs-dev:3k:validation:attempt-1:revision-1
            digest: sha256:0a54f1baec0120c1c93523e6900a6ce28e025c570289e5dfa9835e28099042ac
            digest-domain: canonical-bytes
            room-ref: ./review/validation/briefing-1
          resolution:
            type: Resolution
            id: resolution:actor-1784675152206198000
            briefing: briefing:docs-dev:3k:validation:attempt-1:revision-1
            by: person:reviewer
            at: "2026-07-21T23:05:52.206202Z"
            decision: revise
            reason: >-
              tell me exactly when is rebind and supersede used. and you should still link
              to the original entity and spec as part of the briefing package. btw do we
              have the semantic to support linking file in the repo (without duplicate it
              into the briefing package dir?
            adoption-note: Adopted by the captain
```

The retained association must bind the exact Result digest, provider Briefing identity,
canonical Briefing identity and revision, and the complete provider presentation mapping.
The recorder verifies that evidence and actor authority before normalizing only the
Resolution's Briefing identity. Matching one primary artifact revision is insufficient.
The frozen late Result in this package deliberately has no such association, so it is a
negative fixture and must be rejected rather than laundered into the multi-artifact package.

### 3. Record a chat decision

Before: the same canonical attempt is open. Command:

```text
spacedock gate record 3k --decision approve --actor agent:first-officer --reason "All retained ACs reproduced" --directive "Captain: approve after the reset" --workflow-dir docs/dev
```

After: the binary constructs the portable Resolution; no Result-shaped file is fabricated:

```yaml
        - id: gate-attempt:3k-validation-1
          briefing:
            id: briefing:docs-dev:3k:validation:attempt-1:revision-1
            digest: sha256:0a54f1baec0120c1c93523e6900a6ce28e025c570289e5dfa9835e28099042ac
            digest-domain: canonical-bytes
            room-ref: ./review/validation/briefing-1
          resolution:
            type: Resolution
            id: resolution:spacedock:docs-dev:3k:validation:1
            briefing: briefing:docs-dev:3k:validation:attempt-1:revision-1
            by: agent:first-officer
            at: "2026-07-22T12:00:00Z"
            decision: approve
            reason: All retained ACs reproduced
            adoption-note: "Captain directive: approve after the reset"
```

### 4. Validate

```text
spacedock gate validate 3k --workflow-dir docs/dev
```

This reports the current gate/attempt/Briefing/Resolution and exits nonzero on an invalid
legacy or minimal v1 record. Before and after entity bytes are identical.

## Unreleased v1 compatibility

There is no v2 and no migration. Existing hand-authored and 0260 dogfood histories keep
their `current.attempt`, `current-attempt`, `sequence`, `previous-attempt`, `state`, notes,
applications, and unknown fields. Reads never rewrite them. A recorder mutation preserves
untargeted legacy nodes and extensions; where a present legacy pointer/state must reflect
the requested close or appended attempt, the binary updates that existing field but never
adds it to a minimal record. Contradictory legacy pointers fail closed. New attempts omit
derivable sequence/lineage/state fields; ordered attempts, Resolution absence/presence, and
the selected logical gate are normative.

The eight frozen production histories remain compatibility fixtures. Tests prove unchanged
status output, gates-only targeted mutation, application/extension preservation, and no
bulk rewrite. A/B/C Briefing changes retain one open attempt; closure freezes C; binding D
appends a new attempt. These behaviors continue to satisfy AC-1, AC-4, AC-6, AC-10, AC-12,
AC-13, and AC-14; only their command/fixture mechanics change.

## Scope and references

Expected incremental surface is edits to `internal/gates/{model,operation,io,gates_test}.go`,
`internal/cli/cli.go`, focused status projection tests if necessary, test-only exact
Result/association fixtures, and contract/reference docs: ~220-360 production LOC touched
with net production LOC no higher than `9d279b87`, ~300-500 test/fixture LOC, and ~80-150
documentation lines, tolerance 2x. Provider launch/retention stays with xb; application
lifecycle stays with h1.

Full references (not the presented body): the entity
[`../../../index.md`](../../../index.md) and the landed baseline contract
[`../briefing-16/contract-snapshot.md`](../briefing-16/contract-snapshot.md). The entity link
is navigational; the baseline contract is already a frozen package artifact. A direct
repository URI + SHA becomes a portable Briefing Reference only when the presentation
resolver can reproduce those exact bytes; otherwise the presentation side freezes a room
copy as it did for the contract.
