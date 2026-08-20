---
id: gf0jvhj4y8vjhd6ww62vsb2q
title: Gate rooms make a flat entity a hybrid, and converting it breaks every later gate
status: backlog
source: "Issue #739 (captain, 2026-08-20): gate prepare fails after a flat-to-folder conversion because a retained room-ref resolves relative to the new entity folder and duplicates the slug."
started:
completed:
verdict:
score:
worktree:
issue: "#739"
gates:
    version: 1
    records:
        - id: gate:gf0jvhj4y8vjhd6ww62vsb2q:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:gf0jvhj4y8vjhd6ww62vsb2q-backlog-1
              briefing:
                id: briefing:gf0jvhj4y8vjhd6ww62vsb2q:backlog:attempt-1:revision-1
                digest: sha256:2539b1d74e3671d51fc174a7b96ef4a4368e4ad0807828c8a567b02bb113ec08
                request-digest: sha256:6ba250fdd78031eabba8a169306935a272bb9c6d734c72192e859adfd856c308
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:gf0jvhj4y8vjhd6ww62vsb2q:backlog:1
                briefing: briefing:gf0jvhj4y8vjhd6ww62vsb2q:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-20T17:23:13.463297Z"
                decision: approve
                reason: 'Captain: ''dispatch gf0j''. Shape the fix, weighing enforcement at gate prepare against a tolerant read, and account for the 10 existing hybrids.'
              application:
                target-stage: ideation
                state: pending
---

Stop `gate prepare` from leaving an entity in a form whose retained references break on conversion, and unblock the entities already in that state.

## Problem

A retained `room-ref` is written relative to the entity file's directory. That directory changes meaning between the two entity forms, so the same room resolves to two different paths:

| form | directory | ref written |
|---|---|---|
| flat `<state>/<slug>.md` | `<state>` | `./<slug>/review/...` |
| folder `<state>/<slug>/index.md` | `<state>/<slug>` | `./review/...` |

Convert flat to folder and every retained ref resolves to `<state>/<slug>/<slug>/review/...`. Issue #739 reports the doubled slug in that path.

One old attempt blocks all later gates. `internal/gates/io.go:207` loops over every prior attempt that carries a `RequestDigest` and returns an error when any retained room is unreadable. A consumed historical attempt therefore makes every future gate unpreparable, which is what #739 observed with a green, captain-approved validation gate.

Five sites join a room-ref the same way: `internal/gates/io.go:207`, `internal/gates/application.go:240`, `internal/gates/operation.go:181`, `internal/gates/operation.go:560`, and the write side at `internal/gates/prepare.go:222` (`relativeRoomRef`, line 671).

The deeper fault is the form itself. `gate prepare` on a flat entity creates `<state>/<slug>/review/...`, so a `<slug>/` directory appears beside `<slug>.md`. Every flat gated entity is already this hybrid. In this repository all 10 flat entities that carry gates have a sibling directory of the same name, against 122 flat and 65 folder entities in total. The hybrid is one `git mv` away from folder form, and that move is exactly what breaks the refs.

`gate record --round` already refuses this: `internal/gates/operation.go:108` returns "gate record --round requires folder-form entity <slug>/index.md because review artifacts accumulate beside the entity". `gate prepare` creates review artifacts beside the entity and enforces nothing. Two commands in one package disagree about the same rule.

## Proposed approach

Ideation owns the choice. Three candidates, in the order the problem suggests:

1. Extend the round recorder's rule to `gate prepare`: refuse to create a room for a flat entity. This prevents the class, matches the existing precedent, and removes the hybrid. It blocks the 10 flat gated entities here on their next gate unless a migration ships with it.
2. Ship a conversion that moves `<slug>.md` to `<slug>/index.md` and rewrites retained refs from `./<slug>/review/...` to `./review/...`. The retained Briefing and request digests must survive unchanged, because a room's value is its immutability.
3. Resolve legacy refs tolerantly at read time: when a ref resolves to a missing path and its first segment equals the entity slug, retry without that segment. This unblocks existing state without rewriting any retained bytes, and it leaves the hybrid in place.

Option 1 with option 2 prevents recurrence and clears the backlog of hybrids. Option 3 alone is the smallest change that unblocks a stuck workflow. Weigh 1+2 against 3 on the cost of migrating entities in other checkouts, which cannot be enumerated from here.

## Out of scope

Do not change the room layout, the Briefing schema, digest computation, or the gate ceremony's ordering. Do not discard gate history. Do not create a replacement attempt to route around an unreadable room.

## Expected surface and tolerance

Ideation sets this. The five join sites plus `relativeRoomRef` bound the read-side work. A conversion path adds a command surface and its tests.

## Acceptance criteria

Ideation writes these. They must include a fixture that reproduces #739's doubled-slug path and fails before the change, and a case proving a consumed historical attempt no longer blocks a later gate.

## Test plan

Ideation writes this. The reproduction is offline and cheap: build a flat entity with a consumed gate, convert it to folder form, and prepare a later gate.
