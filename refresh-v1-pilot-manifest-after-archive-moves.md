---
title: Refresh the v1 pilot manifest after archive moves
status: ideation
source: "Captain-directed live-test-truth prerequisite, 2026-08-04"
score: 1.0
sprint: live-test-truth
group: truthful-results
sprint-readiness: ready
id: v72wj17717g4xpkhss7mhv06
gates:
    version: 1
    records:
        - id: gate:v72wj17717g4xpkhss7mhv06:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:v72wj17717g4xpkhss7mhv06-backlog-1
              briefing:
                id: briefing:v72wj17717g4xpkhss7mhv06:backlog:attempt-1:revision-1
                digest: sha256:dc05dceae13f21d186bec0f876a5274a694ab6936a7087de5fd4ede203697f31
                request-digest: sha256:c3384e6b08937607a5425293f40d531a5a7a7fc026e316a7fc3413ea25157899
                room-ref: ./refresh-v1-pilot-manifest-after-archive-moves/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:v72wj17717g4xpkhss7mhv06:backlog:1
                briefing: briefing:v72wj17717g4xpkhss7mhv06:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-04T13:33:48.63634Z"
                decision: approve
                reason: Captain delegated the conn and directed this exact two-file manifest repair; reproduced failures and archive history show the plan restores current-checkout truth without widening semantics.
              application:
                target-stage: ideation
                state: consumed
---

Restore truthful current-checkout validation after seven v1 pilot tasks moved from active state into the archive.

## Outcome

The checked-in v1 pilot manifest resolves every named task at its current path. The focused manifest test, the full Go suite, and the race-enabled suite pass against the current state checkout.

## Problem

`TestV1PilotManifestReadsAndValidates` fails on seven paths that were correct when code commit `9ff2aa50c` recorded the manifest. Normal workflow completion later moved those tasks under `_archive/`, but the manifest kept their active paths and its archive-count assertion stayed at 15.

The failure blocks 3d validation even though 3d does not own the manifest. A historical state snapshot is not a valid substitute because the test deliberately resolves the current `docs/dev/.spacedock-state` checkout.

## Root-cause evidence

The current focused test fails exactly these seven active paths:

- `bind-post-rework-briefing-at-rejection-regate.md`
- `collapse-gate-approval-ceremony/index.md`
- `cut-workflow-specific-round-recorder-from-v1/index.md`
- `minimize-v1-gate-application-schema/index.md`
- `shared-git-scaffold-helper.md`
- `status-pagination-and-default-sorting.md`
- `status-where-robust-and-discoverable.md`

Each corresponding `_archive/` path exists. State history records a normal archive commit for each task. The archive count therefore changed from 15 to 22 while the total manifest cardinality stayed 31.

## Proposed approach

Replace only the seven stale manifest entries with their `_archive/` paths. Change the independent archive-count assertion from 15 to 22.

Do not restore duplicate active records, change the strict gate reader, weaken path checks, add a compatibility layer, or pin tests to a historical state snapshot.

No spike is needed: the existing focused test already fails on every stale path, all seven destination paths exist, and the same test validates path presence, strict decoding, gate validity, total cardinality, and archive cardinality.

## Out of scope

- Changes to 3d candidate code or evidence.
- State-history rewrites or active-record restoration.
- Gate schema, command behavior, stored formats, or runtime behavior.
- A new manifest generator, linter, or standing reconciliation process.

## Expected surface

| File | Expected change | Estimate |
| --- | --- | ---: |
| `internal/gates/testdata/v1_pilot_manifest.txt` | Move seven bindings to `_archive/` | +7 / -7 |
| `internal/gates/application_test.go` | Update archive cardinality from 15 to 22 | +1 / -1 |

Expected total: two files, 8 insertions, and 8 deletions. Tolerance: no additional files and at most two additional insertions or deletions.

No command grammar, stored format, authority, or runtime semantics change. This task updates the checked-in current-state test oracle only.

## Acceptance criteria

**AC-1 (VALUE) - The v1 pilot manifest names all 31 records at their current paths, including exactly 22 archived records.**

Verified by: `go test -v ./internal/gates -run '^TestV1PilotManifestReadsAndValidates$' -count=1` passes against the current state checkout. Reverting any moved path or restoring the count to 15 makes the test fail.

**AC-2 (VALUE) - Mandatory current-checkout repository validation is green again.**

Verified by: `go test ./...` and `go test ./... -race` both exit 0 without `SPACEDOCK_STATE_ROOT` overrides or historical snapshots.

**AC-3 - The repair changes only the stale desired-state bindings and their independent archive-count oracle.**

Verified by: the candidate diff contains only the two expected files and stays within the approved +10/-10 ceiling. This mechanism serves AC-1; broader parser or state changes fail the surface check.

## Test plan

The current focused test is the failing baseline and names all seven missing paths. After the two-file repair, run:

```bash
gofmt -w ./cmd ./internal
go test -v ./internal/gates -run '^TestV1PilotManifestReadsAndValidates$' -count=1
go test ./...
go test ./... -race
git diff --check
```

No live runtime lane is required because this task changes only an offline current-state test oracle. Relevant PR CI must run the repository offline gate before merge.
