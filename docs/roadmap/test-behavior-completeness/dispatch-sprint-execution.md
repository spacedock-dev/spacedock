# Commander dispatch: test-behavior-completeness

This package is a **draft fail-closed handoff**. Do not start implementation
until `staff-review.md` has no open Material finding and all ideation gates are
closed-approved.

When the captain activates this package, you are the Commander for this sprint.
Assume `spacedock:first-officer` for the complete session. Tell each worker and
reviewer: “We love you.” Require `$simple-english` in every report.

## Mission

Drive the approved sprint to current live evidence. Each executable cell passes
or runs as strict XFAIL. Each remaining TODO names an execution path that cannot
run.

Do not implement, review, mutate, approve, or merge `g3` or durable-decisions
work.

## Cold boot

1. Read `docs/roadmap/test-behavior-completeness/index.md`.
2. Read `docs/roadmap/test-behavior-completeness/staff-review.md`.
3. Read `docs/runtime-live-ci.md` and the desired registry.
4. Run the membership query from the index.
5. Stop if any proposed task is absent, not ideation-approved, or deferred.
6. Stop if a Material staff finding lacks a durable task fold.
7. Read each approved task body before dispatch.

## Authority boundaries

Captain activation grants authority to dispatch approved implementation and
validation work for sprint members. It also grants ordinary execution-gate and
merge judgment within this package.

This authority has these limits:

- Do not change desired runtime targets, lane policy, release policy, acceptance
  criteria, or line tolerances.
- Do not add a product repair to `xp6`.
- Do not move `47g`, `g3`, or durable-decisions work into this sprint.
- Do not approve an environment deployment that needs a separate human approval.
- Do not merge a task when a required lane is red, skipped, unapproved, stale,
  or built from another candidate.
- Do not waive the strict-XFAIL baseline or exact-candidate proof.
- Escalate a third feedback cycle, a line-budget breach, an irrecoverable block,
  or a mechanism change.

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

Opus runs as pre-release evidence. Pi runs through optional manual evidence.
These lanes do not replace required Sonnet or Codex evidence.

## Dispatch sequence

### Phase 1 — Restore CI and land strict XFAIL

Finish `0a` first. Validate its event matrix and exact manual Pi cadence. Merge
it before `ts` because both tasks edit the live guide.

Then implement `ts`. Its first landing must contain both real
`default-headless-gate-stop` XFAIL results for Sonnet and Codex. The classifier
accepts only one sole matching semantic code. XPASS and every different or
additional code fail.

Do not dispatch product repair before `ts` merges.

### Phase 2 — Run four repair lanes

After `ts` merges, dispatch these lanes in parallel:

1. `98a`, then the proposed Pi headless-hold task.
2. `6x5`, then `9a`.
3. The recarved `zh`, then the proposed Codex rejection task.
4. The proposed Pi gate-commit task.

Each task follows this exact order:

1. Rebase onto the latest required predecessor.
2. Add the target XFAIL binding with its active owner and stable code.
3. Commit that baseline before product bytes.
4. Run the complete focused journey on that baseline commit.
5. Stop if the result skips, changes code, or adds another code.
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
4. recarved `zh`
5. proposed Codex rejection continuation
6. proposed Pi gate commit
7. proposed Pi headless hold

Before each merge, rebase onto the prior landing. If the rebase changes product
bytes or a source binding, rerun the exact focused lane on the new candidate.

### Phase 3 — Evidence capstone

Dispatch `xp6` only after every product repair above merges.

The capstone can change target binding rows only. It removes a binding after an
exact passing run. It converts a stable failure only when another active task
owns its product repair and supplies one stable code.

The passed Codex withdrawn-gate evidence permits removal of that TODO. Do not
bring `47g` into the sprint.

Keep an Opus TODO only when the authenticated execution path is unavailable.
Authentication failure does not become XFAIL.

## Collision controls

These files are shared merge surfaces:

- `internal/ensigncycle/shared_live_runner_test.go`
- `skills/first-officer/references/fo-dispatch-core.md`
- `docs/runtime-live-ci.md`
- feedback-flow process text

Assign one merge owner: the Commander. Workers do not resolve another task's
binding during a feature merge. After every conflict resolution, run registry
reconciliation and inspect all TODO and XFAIL rows.

Use separate artifact roots for concurrent live runs. Do not reuse metrics,
session logs, or temp workflow state across targets.

## Capstone close

The sprint is complete only when the completion definition in the index is true.
Then run one independent pre-cut antipattern audit on the assembled `main`.

Do not cut a release without separate captain authorization. Stable releases
come from `main` through the annotated-tag procedure in `docs/releasing.md`.
