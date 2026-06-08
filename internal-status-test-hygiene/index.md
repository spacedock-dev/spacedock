---
id: qyc6g8bmvcdsj7bdz7sjwgbn
title: internal/status test-hygiene — go test from root fails on debrief frontmatter, plus gofmt-dirty files
status: validation
source: "captain (2026-06-04) — surfaced by the xa ideation ensign and verified this session. `go test ./...` from the project root fails because TestMigrationCheckFixturesParseConsistently scans the .spacedock-state debrief fixtures; two internal/status files are also gofmt-dirty on next."
score: "0.26"
started: 2026-06-08T05:11:44Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-internal-status-test-hygiene
issue:
sprint: 019x-pre-flip-cleanups
group: release-hygiene
sprint-readiness: ready
mod-block: merge:pr-merge
---

`go test ./internal/status/` (and thus `go test ./...`) FAILS when run from the project root, and two `internal/status` files are gofmt-dirty on `next`. Neither is caught by CI (the offline gate runs a fresh checkout with no `.spacedock-state`; worktree ensigns also lack it), so the failures only bite a developer or a non-worktree ensign running tests from the repo root — which is the common dev loop.

## The failures (verified 2026-06-04)

- **`TestMigrationCheckFixturesParseConsistently` fails from root.** It scans `docs/dev/.spacedock-state/_debriefs/*.md` and the debrief `session-date` field parses inconsistently between two paths — `migration_check_test.go:113`: `key "session-date" reader="2026-06-03" direct="2026-06-03T00:00:00Z"`. A bare YAML date (`session-date: 2026-06-03`) is read as a string by one parser and a full timestamp by the other. The test is checking parser consistency across fixtures + live state, and the debrief frontmatter (which is NOT entity frontmatter) trips it.
- **Two gofmt-dirty files on `next`:** `internal/status/external_proof.go` and `internal/status/no_yaml_silent_drop_test.go` (a curly-quote-in-comment quirk per prior validation notes). `gofmt -l internal/status/` lists both on a clean tree.

## Direction (for ideation to flesh out)

- The migration-check test should likely **exclude `_debriefs/` (and other non-entity dirs)** from its fixture scan — debriefs are session records with their own frontmatter shape (`session-date`, `sequence`, `first-commit`), not entity fixtures the migration check governs. Alternatively, make the date-field parse consistent for bare YAML dates. Pick the smaller, correct fix.
- `gofmt -w` the two dirty files (and consider why CI's offline gate doesn't gofmt-check — a `gofmt -l` guard in the offline lane would prevent recurrence).
- AC: `go test ./...` passes from the project root with `.spacedock-state` present, and `gofmt -l ./...` is empty — both checkable by command.

## Notes

Cleanup/hygiene task (the change is the fix plus the now-passing command). Surfaced by the `feedback-guarantee-binary-gate` (xa) ideation ensign's out-of-scope flag; verified by the FO. Provenance: session 2026-06-04 #8.

## Stage Report: implementation

- DONE: AC: go test ./... passes from the repo ROOT with docs/dev/.spacedock-state present — fix TestMigrationCheckFixturesParseConsistently (exclude _debriefs/ from the walk)
  Pruned the `_debriefs` subtree via `filepath.SkipDir` in the walk callback (worktree commit 6fe40471). Verified: copied the fix into the main repo, `go test ./...` from root went from `1 failed` to green for the migration-check test, then reverted the main-repo copy to pristine. With the machine-local `.claude/worktrees` scratch set aside, full `go test ./...` from root = `1141 passed`.
- DONE: gofmt: confirm gofmt -l ./... is empty
  `gofmt -l .` in the worktree is empty; external_proof.go and no_yaml_silent_drop_test.go are already clean (no `gofmt -w` needed, matching staff-review B1).
- DONE: Scope: migration-check fix only
  Single-file, 9-line change to migration_check_test.go. Did not touch the gofmt half or the contractlint sweep.

### Summary

Chose the exclude-`_debriefs/` direction (smaller, correct: debriefs are session records with a distinct frontmatter shape, not entity fixtures the migration check governs) over normalizing date scalars, which would have altered the consistency oracle. The proof — `go test ./...` from the repo root with `.spacedock-state` present — passes. NOTE for validation: a SEPARATE, pre-existing failure surfaces when running `go test ./...` from root on THIS machine — `TestNoInstructionReadsOutsideQuarantine` (internal/contractlint) walks into untracked machine-local agent-team checkouts under `.claude/worktrees/agent-*/` (its sweep skips `.worktrees` but not `.claude/worktrees`). That is out of this task's scope and not a repo defect (the `.claude/` tree is untracked scratch); flagging it for visibility, not fixing it. [SUPERSEDED — see cycle 2: team-lead extended the AC to cover this; now fixed.]

## Stage Report: implementation (cycle 2)

- DONE: AC: go test ./... passes from the repo ROOT (extended) — fix TestNoInstructionReadsOutsideQuarantine (the contractlint sweep tripping on `.claude/worktrees/agent-*`)
  Added `.claude` to the boundary sweep's `SkipDir` case in `boundary_guard_test.go:88`, alongside `.worktrees` (worktree commit 160780cc). `.claude/worktrees/agent-*` are untracked Claude-Code agent-team checkouts — the same class of agent scratch as `.worktrees`, not the repo's shipped instruction surface.
- DONE: Regression test that does not depend on the live `.claude` tree
  TDD: added `TestBoundaryGuardSweepSkipsAgentWorktrees` (`boundary_guard_control_test.go`) — plants an instruction-read fixture inside both a `.worktrees/agent-*/skills/integration/` and a `.claude/worktrees/agent-*/skills/integration/` temp-dir path, asserts neither is reported while a real out-of-quarantine read still is. Confirmed RED before the fix (flagged the `.claude` path), GREEN after. Full contractlint package: 21 passed.
- DONE: End-to-end proof, `go test ./...` green from root with live `.claude` scratch present
  Copied all three changed files into the main repo (which has `.spacedock-state` and 8 live `.claude/worktrees/agent-*` dirs), ran `go test ./...` from root = `1142 passed` (no files moved aside), then reverted the main-repo copies to pristine. The whole AC is now satisfied without machine-local set-aside.
- DONE: gofmt / vet clean
  `gofmt -l .` empty in the worktree; `go vet ./internal/contractlint/` clean.

### Summary (cycle 2)

Per team-lead's extension, fixed the sibling failure I had flagged — the contractlint boundary sweep now prunes `.claude` the same minimal, same-theme way it already prunes `.worktrees`. With both worktree commits (6fe40471 migration-check, 160780cc contractlint), `go test ./...` from the repo root is genuinely green (`1142 passed`) with `.spacedock-state` and live `.claude/worktrees` scratch present. No new AC was needed; this completes the same "go test green from root" criterion. Both changes are test-only and minimal.

## Stage Report: validation

- DONE: AC (inspect + confirm correctness): confirm migration_check_test.go excludes _debriefs subtree (filepath.SkipDir), and contractlint boundary sweep now skips .claude the same way as .worktrees. Both minimal and correct.
  migration_check_test.go:75-77 prunes `_debriefs` via `filepath.SkipDir` (9-line add, commit 6fe40471). boundary_guard_test.go:88 adds `.claude` to the same `SkipDir` case as `.worktrees` (1-line change, commit 160780cc). Inspected both; minimal and on-theme.
- DONE: Regression has teeth — TestBoundaryGuardSweepSkipsAgentWorktrees plants BOTH .worktrees and .claude scratch and asserts neither flagged; sanity-check it goes RED without the .claude skip.
  Test plants both scratch trees + a real offender (boundary_guard_control_test.go:138-177): asserts offender IS flagged (sweep live) and no `.worktrees/`/`.claude/` path is. Revert-and-rerun: removing the `.claude` skip → RED with msg `sweep flagged an agent-worktree scratch read ".claude/.../skill_surface_test.go"`; restored → GREEN. Bonus probe: removing the `.worktrees` skip also REDs → teeth on both halves, not trivially passing. Tree restored clean (empty `git diff`).
- DONE: Run go test ./... in the worktree — green; gofmt -l ./... empty.
  `go test ./...` exit 0, all 17 packages `ok` (incl. internal/status 8.186s, internal/contractlint). `gofmt -l .` empty (0 dirty files). Worktree `git status` clean.

### Summary

PASSED. Both changes are minimal, correct, and on-theme — `_debriefs` pruned from the migration-check walk; `.claude` pruned from the boundary sweep alongside `.worktrees`. The cycle-2 regression test has genuine teeth: revert-and-rerun confirmed it REDs when either skip is removed (caught the exact scratch-read path) and GREENs when restored. Worktree `go test ./...` is green (17/17 packages) and `gofmt -l` is empty. The high-stakes contractlint boundary-guard surface was adversarially probed inline (revert-and-rerun on both skip halves) — no material test-strength hole found. The live-state DoD (`go test ./...` from the main repo root with `.spacedock-state` + `.claude` present) is the FO's post-merge integration test; the impl ensign cited 1142 green via the copy-to-main-revert dance.
