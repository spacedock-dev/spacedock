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
        - id: gate:v72wj17717g4xpkhss7mhv06:validation
          stage: validation
          attempts:
            - id: gate-attempt:v72wj17717g4xpkhss7mhv06-validation-1
              briefing:
                id: briefing:v72wj17717g4xpkhss7mhv06:validation:attempt-1:revision-1
                digest: sha256:1f3746cbc8dc880458d56c04eff33dfd5d09045830aab3ac0dbcb79cd050e228
                request-digest: sha256:88e5098aa969046f49372d16a265a560b00466b0b079529fab73232b55ade9cd
                room-ref: ./refresh-v1-pilot-manifest-after-archive-moves/review/validation/briefing-1
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
  Canonical commit `0b8ca50dee7b8c6c6b2cf5e09a1580452f55b30f` makes the seven path substitutions and both 15-to-22 oracle-line updates; the focused test passed all 31 records, and reverting a binding or count makes its path/cardinality checks fail.
- DONE: Keep the strict Read, Validate, and approval-application-node assertions unchanged; do not restore active records or add a historical fixture snapshot.
  The focused pass exercised strict `Read`, `Validate`, and exact `[state target-stage]` application keys against state HEAD `ab33f713035bd1767df28b05663913ef606c5257`; the commit has no parser, state, snapshot, runtime, or `3d` file.
- DONE: Run the focused manifest test, `go test ./...`, `go test ./... -race`, `gofmt`, and `git diff --check`; commit only the two expected oracle files.
  `gofmt -w ./cmd ./internal`, `go test -v ./internal/gates -run '^TestV1PilotManifestReadsAndValidates$' -count=1`, both repository suites, and `git diff --check` exited 0 before the topology-only transplant; no long suite was rerun during this handback, as directed.
- DONE: Verify the approved surface after formatting and tests.
  Exact path/name/count assertions and the final diff confirmed 31 unique manifest lines, 22 `_archive/` lines, seven absent active bindings, seven present archive bindings, and no candidate edits outside the two oracle files.
- DONE: Record the topology-only candidate transplant onto `origin/main`.
  The FO rebased the sole repair commit off inherited local-main planning ancestry; the canonical candidate is 0 behind/1 ahead, remains exactly two files and +9/-9, and shares stable patch-id `164f505c3f1ae7f4df49d05e1590c6a2b777b91f` with superseded commit `3d5af7c92d7d17f6d0c5ac6742f14a2fa7729af2`.

### Summary

The canonical candidate `0b8ca50dee7b8c6c6b2cf5e09a1580452f55b30f` is topology-clean against `origin/main` and preserves the exact approved two-file +9/-9 repair patch. Existing focused, full, and race evidence remains attributed to the pre-transplant run; this handback verified topology and stable patch identity only, without rerunning suites or changing code.

## Stage Report: implementation (cycle 3)

- DONE: Finalize the two-file oracle repair after the authorized archived-state normalization, preserving 31 total entries, 22 archived entries, and strict application-node assertions.
  The implementation worktree is clean at canonical commit `0b8ca50de` (`test: refresh v1 pilot manifest archive paths`), a topology-only transplant of the same two-file repair formerly recorded as `3d5af7c92`. It changes only `internal/gates/testdata/v1_pilot_manifest.txt` and `internal/gates/application_test.go` with the approved +9/-9 surface. The manifest has seven current `_archive/` bindings; the test still requires 31 unique records, 22 archived paths, strict `Read`/`Validate`, and exact `[state target-stage]` approval application keys.
- DONE: Verify the authorized state normalization before claiming the focused oracle pass.
  State commit `ab33f7130` (`state: normalize archived gate application fields`) removed legacy `action: advance` and `blockers: []` from all three approval applications in `_archive/status-where-robust-and-discoverable.md`; the recorded pre-image hash is `3955b5c90873d3f09e2618aed7c34e126f7401767cae449b43cff7ddf49aee51` and the post-image hash is `8550aac8adb2c40b5d5b8575f98a168160d419e749fa1c5f55ae1ef9db06935f`.
- DONE: Run the required formatting, focused, full, race, and diff checks from the implementation worktree.
  `gofmt -w ./cmd ./internal` exited 0 with no candidate diff; `go test -v ./internal/gates -run '^TestV1PilotManifestReadsAndValidates$' -count=1` passed all 31 manifest subtests; `go test ./...` exited 0 (including `internal/ensigncycle` in 279.836s and `internal/release` in 23.423s); `go test ./... -race` exited 0 (including `internal/dispatch` in 42.583s, `internal/status` in 46.689s, and `internal/release` in 23.873s); and `git diff --check` exited 0.
- DONE: Confirm no unrelated candidate or state changes remain and leave the implementation ready for fresh validation.
  The code worktree has no uncommitted changes, the candidate commit is path-scoped to the two oracle files, and the state checkout is clean after the report update. No parser, runtime, historical snapshot, or downstream `3d` candidate files were changed.

### Summary

The v1 pilot manifest now truthfully resolves all 31 current records, including 22 archived records, and the strict application-node oracle passes after the separately authorized normalization commit `ab33f7130`. The canonical two-file implementation commit `0b8ca50de` (topology-only transplant of `3d5af7c92`) is complete; focused, full, race, gofmt, and diff-check evidence is recorded, with no unrelated changes.

## Stage Report: validation

- DONE: Independently verify immutable candidate 0b8ca50dee7b8c6c6b2cf5e09a1580452f55b30f is exactly the approved two files and +9/-9, with no parser, snapshot, runtime, state, or 3d change.
  `git show --format=fuller --stat --oneline --decorate --no-renames HEAD`, `git diff --numstat HEAD^!`, and `git diff --name-status HEAD^!` exited 0 and showed only `internal/gates/application_test.go` (+2/-2) and `internal/gates/testdata/v1_pilot_manifest.txt` (+7/-7); final HEAD/status were the requested SHA and clean.
- DONE: Reproduce the 31-record focused current-checkout oracle against split-root state HEAD ab33f713035bd1767df28b05663913ef606c5257, including 22 archived paths and exact canonical approval-application keys.
  Detached state `git rev-parse HEAD` returned `ab33f713035bd1767df28b05663913ef606c5257`; `SPACEDOCK_STATE_ROOT=/tmp/spacedock-v72-validation.sMFpfp/state go test -v ./internal/gates -run '^TestV1PilotManifestReadsAndValidates$' -count=1` exited 0 for all 31 subtests, while an independent count returned `total=31 archived=22 duplicates=0` and the unchanged oracle enforced exact `[state target-stage]` keys.
- DONE: Run independent negative controls showing a reverted archive binding and a 22-to-15 count regression fail; keep the candidate commit immutable by using a temporary detached worktree or equivalent disposable copy.
  In detached candidate `0b8ca50de`, the same focused command exited 1 after one archive binding was reverted (`pilot manifest has 21 archived paths, want 22`) and exited 1 after the oracle was restored from 22 to 15 (`pilot manifest has 22 archived paths, want 15`); the disposable worktrees were removed and the candidate stayed clean.
- DONE: Run gofmt verification, go test ./..., go test ./... -race, and git diff --check; report exact commands, exit status, and any material finding.
  Detached `gofmt -w ./cmd ./internal` plus `git diff --exit-code` exited 0; candidate `go test ./...`, `go test ./... -race`, `git diff --check`, and `git diff --check HEAD^ HEAD` each exited 0 without a state-root override, failure, race report, formatting delta, or material finding.
- DONE: Verify the candidate is 0 behind/1 ahead of origin/main and its stable patch-id equals superseded pre-transplant commit 3d5af7c92d7d17f6d0c5ac6742f14a2fa7729af2 (`164f505c3f1ae7f4df49d05e1590c6a2b777b91f`), so no unrelated sprint-planning commits enter the PR.
  `git fetch origin main` exited 0; `git rev-list --left-right --count origin/main...HEAD` returned `0 1`, HEAD's parent/merge-base was `20cfc809a7342b2428509c76ddd6e91423db39b7`, and `git patch-id --stable` returned `164f505c3f1ae7f4df49d05e1590c6a2b777b91f` for both commits.
- DONE: Verify AC-1, including identity, order, and current-path semantics rather than cardinality alone.
  The default no-override focused command also exited 0; stripping `_archive/` produced the same ordered-manifest SHA-256 before and after, and all seven old active paths were absent while their archive destinations existed at `ab33f7130`.
- DONE: Verify AC-2 and AC-3 with independent repository and surface evidence.
  Both mandatory suites exited 0 against the live shared checkout, while the immutable two-file +9/-9 commit contains no production reader, parser, snapshot, runtime, state, gate-authority, or downstream-3d path.
- SKIPPED: Run a live runtime lane.
  The approved scope explicitly requires no live lane because the repair changes only an offline current-state oracle; no runtime hot path, allocation, I/O, or format behavior changed.
- DONE: Recommend PASSED, with findings classified separately.
  PASSED: all three ACs have falsifiable evidence; material findings: none; deferred risks: none; polish findings: none.

### Summary

Validation independently reproduced the exact-state and live-checkout oracle, both regression controls, full and race suites, formatting/diff hygiene, immutable surface, clean PR topology, and stable patch identity. Candidate `0b8ca50de` satisfies AC-1 through AC-3 with no findings, so the validation recommendation is PASSED.
