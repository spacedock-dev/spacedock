---
title: Refresh the v1 pilot manifest after archive moves
status: validation
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
        - id: gate:v72wj17717g4xpkhss7mhv06:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:v72wj17717g4xpkhss7mhv06-ideation-1
              briefing:
                id: briefing:v72wj17717g4xpkhss7mhv06:ideation:attempt-1:revision-1
                digest: sha256:20b082741d6ec6ed7a6befc6efe5eeb2a9cceb8c1fc1bd3033220801e546b3f0
                request-digest: sha256:cf1b0d6620fc5abd801c7200641c1e36d5ec2c78ca4d3d30899fdc4ac986f8a4
                room-ref: ./refresh-v1-pilot-manifest-after-archive-moves/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:v72wj17717g4xpkhss7mhv06:ideation:1
                briefing: briefing:v72wj17717g4xpkhss7mhv06:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-04T13:48:06.673607Z"
                decision: approve
                reason: 'Science officer concurs: the current-state audit is intentionally live, and the seven missing entries are ordinary archive moves. Rebinding them and changing the independent archive count to 22 preserves strict Read/Validate/application-node coverage and the 31-record invariant.'
              application:
                target-stage: implementation
                state: consumed
started: 2026-08-04T13:34:28Z
worktree: .worktrees/spacedock-ensign-refresh-v1-pilot-manifest-after-archive-moves
---

Restore truthful current-checkout validation after seven v1 pilot tasks moved from active state into the archive.

## Outcome

The checked-in v1 pilot manifest resolves every named task at its current path. The focused manifest test, the full Go suite, and the race-enabled suite pass against the current state checkout.

## Problem

`TestV1PilotManifestReadsAndValidates` fails on seven paths that were correct when code commit `9ff2aa50c` recorded the manifest. Normal workflow completion later moved those tasks under `_archive/`, but the manifest kept their active paths and its archive-count assertion stayed at 15.

The failure blocks 3d validation even though 3d does not own the manifest. A historical state snapshot is not a valid substitute because the test deliberately resolves the current `docs/dev/.spacedock-state` checkout.

## Root-cause evidence

The current focused test fails exactly these seven active bindings, and each destination is a normal state-history rename:

| Stale manifest binding | Current binding | Archive commit |
| --- | --- | --- |
| `bind-post-rework-briefing-at-rejection-regate.md` | `_archive/bind-post-rework-briefing-at-rejection-regate.md` | `a2052245c` |
| `collapse-gate-approval-ceremony/index.md` | `_archive/collapse-gate-approval-ceremony/index.md` | `80dc5cc4d` |
| `cut-workflow-specific-round-recorder-from-v1/index.md` | `_archive/cut-workflow-specific-round-recorder-from-v1/index.md` | `2987a085c` |
| `minimize-v1-gate-application-schema/index.md` | `_archive/minimize-v1-gate-application-schema/index.md` | `c15926037` |
| `shared-git-scaffold-helper.md` | `_archive/shared-git-scaffold-helper.md` | `fcfa65968` |
| `status-pagination-and-default-sorting.md` | `_archive/status-pagination-and-default-sorting.md` | `b67f466e6` |
| `status-where-robust-and-discoverable.md` | `_archive/status-where-robust-and-discoverable.md` | `596c23d1f` |

All seven destination files exist, and `git show --find-renames` reports each entity move as `R099`. The archive count therefore changed from 15 to 22 while the total manifest cardinality stayed 31. Baseline `go test ./...` and `go test ./... -race` runs pass every other package and fail `internal/gates` only on these same seven paths.

## Proposed approach

Replace only the seven stale manifest entries with the corresponding current bindings above, preserving manifest order. In `application_test.go`, change both the archive-count predicate and its failure diagnostic from 15 to 22; leave the independent total-cardinality assertion at 31.

Do not restore duplicate active records, change the strict gate reader, weaken path checks, add a compatibility layer, or pin tests to a historical state snapshot.

The `3d` candidate is an untouched downstream consumer: none of its worktree, entity, code, or evidence changes here. Its existing candidate is revalidated only after this prerequisite lands.

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

The baseline counts the archive-cardinality update as one logical substitution. Because the current literal appears in both the predicate and failure diagnostic, the correct Git diff may realize as +9/-9 overall; that stays inside the approved +10/-10 ceiling.

No command grammar, stored format, authority, or runtime semantics change. This task updates the checked-in current-state test oracle only.

## Acceptance criteria

**AC-1 (VALUE) - The v1 pilot manifest names all 31 records at their current paths, including exactly 22 archived records.**

Verified by: `go test -v ./internal/gates -run '^TestV1PilotManifestReadsAndValidates$' -count=1` passes against the current state checkout. Reverting any moved path or restoring the count to 15 makes the test fail.

**AC-2 (VALUE) - Mandatory current-checkout repository validation is green again.**

Verified by: `go test ./...` and `go test ./... -race` both exit 0 without `SPACEDOCK_STATE_ROOT` overrides or historical snapshots.

**AC-3 - The repair changes only the stale desired-state bindings and their independent archive-count oracle, while the `3d` candidate remains untouched.**

Verified by: the candidate diff contains only the two expected files and stays within the approved +10/-10 ceiling; the `3d` worktree and state entity have no task-attributable diff. This mechanism serves AC-1; broader parser, state, or downstream-candidate changes fail the surface check.

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

## Stage Report: ideation

- DONE: Confirm the focused current-checkout failure maps exactly to seven normal archive moves and changes archive cardinality from 15 to 22 while total cardinality stays 31.
  The focused test fails only the seven named active paths; all destinations exist as R099 archive moves, and an in-memory rebinding yields 22 archived entries out of 31.
- DONE: Produce a complete two-file design whose value proof is the existing focused oracle plus current-checkout full and race suites, without adding a generator, parser change, state rewrite, or standing check.
  Both baseline suites fail only on the seven stale bindings; AC-1 catches a reverted binding/count, and AC-2 requires both repaired suites against the default current state checkout.
- DONE: Keep the expected surface at `internal/gates/testdata/v1_pilot_manifest.txt` and `internal/gates/application_test.go`, +8/-8 with the declared tolerance, and preserve 3d as an untouched downstream consumer.
  The baseline remains two files and +8/-8 with a +10/-10 ceiling; both count-literal sites may realize +9/-9, and the disjoint `3d` worktree/entity remain out of scope.

### Summary

Ideation now records the exact seven archive rebindings and their rename evidence while preserving 31 total and 22 archived records. The implementation remains a two-file offline-oracle repair proved by the focused, full, and race suites; no generator, parser, state-history, runtime-lane, or `3d` change is introduced.

## Stage Report: implementation

- DONE: Apply the seven current `_archive/` bindings and the independent archive-count update in the candidate worktree, leaving the two expected oracle files as the only candidate edits (+9/-9).
  `internal/gates/testdata/v1_pilot_manifest.txt` now names all seven archived destinations, and `internal/gates/application_test.go` expects 22 archived paths while retaining the 31-entry invariant and strict Read/Validate/application-node assertions. The candidate changes remain uncommitted pending the finding below.
- DONE: Run a read-only focused falsifier against the intended two-file repair and preserve the exact new finding without widening the candidate.
  The focused test passes all seven path-presence/Read/Validate checks until `_archive/status-where-robust-and-discoverable.md`, then fails unchanged strict application-node coverage: `approval application keys = [action state target-stage], want [state target-stage]` at `application_test.go:145`. The extra `action: advance` is present on all three approval applications in the archived entity; the test therefore fails one subtest, not a stale-path or count assertion.
- DONE: Trace the finding to durable provenance and the merged legacy-warning reader/writer behavior.
  State commits `0da26c7ba` and `49329e0c9` explicitly authorized `action: advance` frontmatter repairs while the stale binary/schema required that key; `596c23d1f` archived the entity unchanged. Code commit `247940e4b` subsequently minimized the canonical application to `target-stage` and `state`; merged warning support in `e3a89ae8e` filters unknown application keys for strict decode, returns the original node for CAS/audit, and warns rather than mutating historical bytes. Canonical writers strip the extension only on a later entity mutation.

### Finding disposition

- Materiality: non-material to runtime authority and parser validity. `Read`/`Validate` accept this historical application through the warning path, and `action` is non-authoritative because an approval always derives `advance`; it is material only to this pilot's deliberate canonical-byte audit.
- Ownership: the archived state entity and its historical frontmatter belong to the state/FO maintainer. The two-file manifest candidate must not normalize that state or relax its independent audit assertion; the disjoint downstream `3d` candidate remains untouched.
- Proposed disposition: record the finding and defer a separately authorized, path-scoped state normalization (or an explicit audit-policy decision) to the state owner. Do not modify the archive, alter the strict assertion, rerun suites, or commit the candidate until the First Officer authorizes that follow-up. No race/full-suite result is claimed from this interrupted falsifier pass.

### Summary

The intended seven-path/22-archive oracle repair is staged but uncommitted. Its focused read-only falsifier exposes one independent historical `application.action` finding in the newly rebound archive entity. Provenance and current warning semantics show a compatibility-only legacy field; ownership and any normalization decision remain with the state/FO lane, so no candidate commit or further rerun is made.

## Stage Report: implementation (cycle 2)

- DONE: Rebind all seven stale manifest entries to their current `_archive` paths and update the archive count to 22 while preserving 31 total entries.
  Commit `3d5af7c92d7d17f6d0c5ac6742f14a2fa7729af2` makes the seven path substitutions and both 15-to-22 oracle-line updates; the focused test passed all 31 records, and reverting a binding or count makes its path/cardinality checks fail.
- DONE: Keep the strict Read, Validate, and approval-application-node assertions unchanged; do not restore active records or add a historical fixture snapshot.
  The focused pass exercised strict `Read`, `Validate`, and exact `[state target-stage]` application keys against state HEAD `ab33f713035bd1767df28b05663913ef606c5257`; the commit has no parser, state, snapshot, runtime, or `3d` file.
- DONE: Run the focused manifest test, `go test ./...`, `go test ./... -race`, `gofmt`, and `git diff --check`; commit only the two expected oracle files.
  `gofmt -w ./cmd ./internal`, `go test -v ./internal/gates -run '^TestV1PilotManifestReadsAndValidates$' -count=1`, both repository suites, and `git diff --check` exited 0; the commit is exactly two files and +9/-9.
- DONE: Verify the approved surface after formatting and tests.
  Exact path/name/count assertions and the final diff confirmed 31 unique manifest lines, 22 `_archive/` lines, seven absent active bindings, seven present archive bindings, and no candidate edits outside the two oracle files.

### Summary

The v1 pilot manifest now follows all seven durable archive moves and expects 22 archived records while preserving 31 total records and the strict current-checkout audit. The focused, full, and race suites pass against the repaired split-root state; the committed candidate is limited to the approved two-file +9/-9 oracle update.
