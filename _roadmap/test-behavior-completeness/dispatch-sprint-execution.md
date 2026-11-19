# Commander dispatch: test-behavior-completeness

This package is the **activated cold-boot Commander handoff**. The staff review
has no open Material finding. The captain approved all nine ideation gates and
bound the target train to the `next` development line.

Under this activation, you are the Commander for this sprint. Assume
`spacedock:first-officer` for the complete session. Tell each worker and reviewer:
“We love you.” Require `$simple-english` in every report.

## Mission

Drive the approved sprint to current live evidence. Each executable cell passes
or runs as target-level XFAIL. Each remaining TODO names an execution path that cannot
run.

Do not implement, review, mutate, approve, or merge `g3` or durable-decisions
work.

## Cold boot

1. Read `docs/roadmap/test-behavior-completeness/index.md`.
2. Read `docs/roadmap/test-behavior-completeness/staff-review.md`.
3. Read `docs/runtime-live-ci.md` and the desired registry.
4. Run the membership query from the index.
5. Make sure that all 10 retained members are present.
6. If the staff review has an open Material finding, stop.
7. Make sure that the nine ideation gates remain closed-approved.
8. Read each approved task body before dispatch.
9. Reconcile `0a` before any new sprint dispatch.

## Authority boundaries

The Commander has the conn toward the sprint goal. The Commander can dispatch
approved implementation and validation work for sprint members.

The Commander can use judgment to approve sprint gates, pull requests, relevant
CI lanes, and merges. Escalate uncertainty to `/root`.

This authority has these limits:

- Do not change desired runtime targets, lane policy, release policy, acceptance
  criteria, or line tolerances.
- Do not add a product repair to `xp6`.
- Do not move `47g`, `g3`, or durable-decisions work into this sprint.
- Do not approve a protected environment. Environment approval needs a separate
  captain grant.
- Do not merge a task when a required lane is red, skipped, unapproved, stale,
  or built from another candidate.
- Do not waive the target-XFAIL baseline or exact-candidate proof.
- Escalate a third feedback cycle, a line-budget breach, an irrecoverable block,
  or a mechanism change.
- Do not create or push a stable release tag. Stable tagging needs a separate
  captain grant.

The Commander can merge a pull request only after its validation gate passes and
its merge mod is satisfied. Every diff-required lane must also be green on the
exact candidate. A human retains any repository or environment approval that the
hosting service reserves for a person.

## Required lanes

The offline suite runs for every landing:

```bash
gofmt -w ./cmd ./internal
go test ./...
go test ./... -race
go test ./internal/contractlint -run '^TestRuntimeLiveRegistryReconciliation$'
SPACEDOCK_LIVE_STATE_DIR=docs/dev/.spacedock-state \
  go test ./internal/contractlint -run '^TestRuntimeLiveTODOOwnersAreActive$'
```

Pull-request changes to common live behavior require Sonnet 5 at maximum effort
and Codex Luna at maximum effort. The exact target repair also needs its focused
live cell.

Use local subscription-backed live runs before paid CI whenever possible.

Opus runs as pre-release evidence. Pi runs through optional manual evidence.
These lanes do not replace required Sonnet or Codex evidence.

## Dispatch sequence

### Phase 1 — Restore CI and land target-level XFAIL

`0a` is in implementation and has no implementation Stage Report. Do not treat
its candidate as finished.

Converge these durable facts first:

- entity status and registered worktree
- branch `spacedock-ensign/restore-optional-manual-pi-common-live-ci`
- tip `e838fba69`
- active owner availability
- candidate status and missing implementation Stage Report

Use normal ownership and dispatch rules after convergence. Reuse an available
registered owner only when the runtime and worker identity permit reuse. Use a
fresh implementation dispatch when normal recovery requires it.

Require the implementation Stage Report before validation. Then validate the
event matrix and exact manual Pi cadence. Merge `0a` before `ts` because both
tasks edit the live guide.

Then implement `ts`. Its landing must convert every executable sprint-owned
target to XFAIL under its repair owner. Typed semantic failures are XFAIL.
XPASS and infrastructure failures stay red.

Do not dispatch product repair before `ts` merges.

### Phase 2 — Run four repair lanes

After `ts` merges, dispatch these lanes in parallel:

1. `98a`, then `fh6`.
2. `6x5`, then `9a`.
3. The recarved `zh`, then `dvd`.
4. `2e4`.

Each task follows this exact order:

1. Rebase onto the latest required predecessor.
2. Use the target XFAIL binding with its active repair owner.
3. Commit any owner transfer before product bytes.
4. Run the complete focused journey on that baseline commit.
5. If the result skips or fails on infrastructure, stop.
6. Apply the approved product repair without weakening the assertion.
7. Run the same target on the exact repair candidate.
8. Require XPASS failure while the binding remains.
9. Remove the binding.
10. Run the target again and require PASS.

Implementation work can proceed in parallel after each lane has its valid
baseline. Merge remains serial because tasks share source-binding and contract
files.

Use this merge order:

1. `98a`
2. `6x5`
3. `9a`
4. `zh`
5. `dvd`
6. `2e4`
7. `fh6`

Before each merge, rebase onto the prior landing. If the rebase changes product
bytes or a source binding, rerun the exact focused lane on the new candidate.

### Phase 3 — Evidence capstone

Dispatch `xp6` only after every product repair above merges.

The capstone can change target binding rows only. It removes a binding after an
exact passing run. It converts an executable TODO only when another active task
owns its product repair.

The passed Codex withdrawn-gate evidence permits removal of that TODO. Do not
bring `47g` into the sprint.

Exact run `31340713337` proves the Opus execution path. Keep those targets as
XFAIL until exact passing evidence removes each binding.

## Collision controls

These files are shared merge surfaces:

- `internal/ensigncycle/shared_live_runner_test.go`
- `skills/first-officer/references/fo-dispatch-core.md`
- `docs/runtime-live-ci.md`
- `skills/feedback-rejection-flow/SKILL.md`
- `skills/first-officer/references/pi-first-officer-runtime.md`

Assign one merge owner: the Commander. Workers do not resolve another task's
binding during a feature merge. After every conflict resolution, run registry
reconciliation and inspect all TODO and XFAIL rows.

Use separate artifact roots for concurrent live runs. Do not reuse metrics,
session logs, or temp workflow state across targets.

## Capstone close

The sprint is complete only when the completion definition in the index is true.
Then run one independent pre-cut antipattern audit on the assembled `main`.

Do not cut or tag a stable release. Stable releases come from `main` through the
annotated-tag procedure in `docs/releasing.md`, after a separate captain grant.
