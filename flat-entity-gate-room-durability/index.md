---
id: b9pjkz3rv0svx9rt63kw8yg7
title: "A flat entity's gate room is neither committed by `state commit` nor migratable to folder form"
status: ideation
source: "Observed eight times driving gates on 79 and cn, 2026-07-27/28. Every gate room bound on a flat entity needed a manual path-scoped commit after `state commit` reported success; the obvious repair — converting to folder form — silently breaks prior closed gates."
started:
completed:
verdict:
score: 0.7
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:b9pjkz3rv0svx9rt63kw8yg7:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:b9pjkz3rv0svx9rt63kw8yg7-backlog-1
              briefing:
                id: briefing:b9pjkz3rv0svx9rt63kw8yg7:backlog:attempt-1:revision-1
                digest: sha256:5d63d8641afab33a2f57b0f4054d12583c8e1b6fd4e3d043f94ce0c6edfee97b
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:b9pjkz3rv0svx9rt63kw8yg7:backlog:1
                briefing: briefing:b9pjkz3rv0svx9rt63kw8yg7:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-27T17:14:04.981857Z"
                decision: approve
                reason: 'Captain approved the bound Subspace backlog review: support flat tickets directly, preserve historical room-ref meaning, and make sibling review rooms durable and round-recordable without migration.'
              application:
                target-stage: ideation
                state: consumed
---

Make a gate room on a flat entity durable through the supported command, and give a flat entity a migration path to folder form that does not corrupt its gate history.

## Two defects that compose into a cul-de-sac

A flat entity (`<slug>.md`) may legitimately carry a sibling room directory (`<slug>/review/...`). The shipped precedent is `source-build-compatibility-identity`, whose backlog briefing lives exactly that way. Both defects below apply to that shape.

### 1. `state commit` reports success having committed only the index

Observed on every gate room bound during the 2026-07-27/28 session — eight occurrences across `79` and `cn`, at backlog, ideation and validation gates. `spacedock state commit <slug>` commits `<slug>.md` and pushes, printing `Committed and pushed …`, while `<slug>/review/<stage>/briefing-N/` remains untracked.

The consequence is not cosmetic. The committed frontmatter carries a `room-ref` pointing at a room that does not exist on the state remote, so a peer or a later session cannot retrieve the Briefing that justified a recorded decision. That is the "digests no committed tree reproduces" failure class the durable-decisions sprint was created to eliminate, reproduced through the supported path.

Every occurrence was worked around by hand:

```
git add -- <slug>/ && git commit -m "…" -- <slug>/ && git push origin HEAD:spacedock-state/dev
```

`sync-merge-guard-archive-state` (`rd`) owns the folder-form version of this and explicitly scopes its work to folder-form entities. The flat-plus-sibling-room shape is uncovered by it.

### 2. Converting to folder form silently breaks every prior closed gate

The obvious repair is to convert the entity to folder form, which also unlocks advisory round recording (the recorder requires `<slug>/index.md` and refuses flat entities before locking or writing). It cannot be done safely today.

`room-ref` resolves relative to the **entity's own directory**. Proven against `prepare-provider-neutral-gate-room`, a folder-form entity whose stored `./review/backlog/briefing-1` maps to `<slug>/review/backlog/briefing-1`. For a flat entity the entity directory is the state root, so a correct stored ref reads `./<slug>/review/...`. After `git mv <slug>.md <slug>/index.md` that same ref resolves to `<slug>/<slug>/review/...` — one level too deep, and pointing at nothing.

Nothing catches it. `gate validate <entity>` exits 0 and reports gate, attempt, briefing, resolution and decision, because it reads frontmatter without dereferencing `room-ref`. Observed directly: after converting `79`, validate stayed green while the backlog gate's room pointer was broken.

There is no supported repair once broken. `gate record --briefing` will not rewrite a closed attempt, and hand-editing `gates:` is forbidden by the lifecycle skill. The conversion on `79` was reverted for exactly this reason.

### Why the two compose badly

Round recording requires folder form. A flat entity therefore cannot machine-record advisory rounds — and the migration that would fix that corrupts its closed gate records. An entity filed flat is permanently stuck with hand-authored rounds, and the only escape damages its decision history.

## What a fix needs to decide

- Whether `state commit` should discover and include an entity's sibling room, or whether flat-plus-sibling-room should stop being a legal shape. Those are different answers with different blast radius; `source-build-compatibility-identity` already relies on the current shape.
- Whether `room-ref` should be stored state-root-relative rather than entity-relative, which would make it survive a move — at the cost of changing an already-recorded field's meaning.
- Whether a supported flat-to-folder migration should exist at all, and if so whether it may rewrite `room-ref` on closed attempts without violating the frozen-record rule.
- Whether `gate validate` should dereference `room-ref`. It is the natural place to catch a broken pointer and currently does not look.

## Out of scope

- The folder-form `state commit` boundary, which is `rd`'s.
- Gate-room retention size, which is `9t`'s.
- The request-digest ordering trap, which is a separate defect in the same command family.

## Acceptance criteria

Ideation fills these in. At minimum: a room bound on a flat entity is retrievable from a fresh clone after the supported command alone, with a falsifier that reverting the fix leaves the room absent from the remote; and a broken `room-ref` is detected by a shipped check rather than by a reader failing to find the file.

## Test plan

Ideation fills this in. The existing two-clone real-Git harness is the substrate; do not build a second one.
