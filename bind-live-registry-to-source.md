---
title: Bind the desired live-test registry to tests and fixture builders
status: ideation
source: "Captain semantic lock for docs/runtime-live-ci-registry.md, 2026-08-03"
started: 2026-08-03T10:41:35Z
completed:
verdict:
score: 0.95
worktree:
issue:
sprint: live-test-truth
group: registry
sprint-readiness: ready
id: 3w2rx3aw4vcympx84zt8mtv7
gates:
    version: 1
    records:
        - id: gate:3w2rx3aw4vcympx84zt8mtv7:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:3w2rx3aw4vcympx84zt8mtv7-backlog-1
              briefing:
                id: briefing:3w2rx3aw4vcympx84zt8mtv7:backlog:attempt-1:revision-1
                digest: sha256:31091557195718f78a0f4425d97cdb4fad7653b621456e030b063c5ff7e09fb1
                request-digest: sha256:6dbc7d42c596779e0d531fdd9a0dd73238546e9f3b6279f9ab3882dede2112dc
                room-ref: ./bind-live-registry-to-source/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:3w2rx3aw4vcympx84zt8mtv7:backlog:1
                briefing: briefing:3w2rx3aw4vcympx84zt8mtv7:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T10:33:36.145597Z"
                decision: approve
                reason: Captain approved the prepared Sol ideation cohort with make it so.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

## Problem

The desired live-test registry names journeys, proofs, entry points, and fixture IDs. Source declarations do not carry these stable IDs.

A reader cannot move from a registry promise to its test or fixture builder.

A reader also cannot distinguish an unused fixture from a deliberate helper.

The registry describes desired state. Therefore, reconciliation can find missing implementation. A missing binding must stay visible and must not become a false source annotation.

## Current evidence

The registry contains 16 common journeys, four runtime-specific proofs, and two non-gating experiments. The shared scenario table contains 10 journey IDs.

These six desired journeys do not have records in `sharedRuntimeScenarios()`:

- `full-ensign-cycle`
- `default-headless-gate-stop`
- `withdrawn-gate-recovery`
- `zero-discovery`
- `auto-continue-after-implementation`
- `ac-value-reanchor`

Related standalone tests exist for these journeys. They do not satisfy the canonical shared-scenario declaration rule in the registry.

The source has no function named `TestLiveSharedScenarios`. It has separate Claude and Codex suite functions. Pi has a coverage map with explicit gaps.

Reconciliation must report these facts. This task must not add fake bindings or new live behavior to hide them.

## Ideation spike

The spike used the registry entry `rejection-flow`. Its source sample used this adjacent declaration:

```go
//spacedock:live-journey id=rejection-flow
{
    name: "rejection-flow",
},
```

The one-off join compared the registry ID, the annotation ID, and the record name. It returned this result:

```text
JOIN PASS registry=rejection-flow annotation=rejection-flow declaration=rejection-flow
```

The path guard then compared two reconciliation bases with the approved path list. The old base failed and named all changed files:

```text
GUARD FAIL base=c25a4c1aaa6880f165ea2eb786d65810a29eeab2 changed=internal/ensigncycle/cycle_test.go,internal/ensigncycle/live_test.go,internal/ensigncycle/liveassert_unit_test.go,
GUARD PASS base=2997da9a243f36ad31f483f63ce29b1b838c3410
```

This spike proves the smallest required path. The semantic join works without builder names in the registry.

The guard detects a stale base and accepts a reconciled base.

## Proposed approach

### Annotation grammar

Keep the four annotation forms from the registry. Do not add status fields, quoted values, or builder locations.

```text
//spacedock:live-suite lanes=<lane>[,<lane>...]
//spacedock:live-journey id=<registry-id>
//spacedock:live-fixture id=<fixture-id>
//spacedock:live-proof id=<registry-id> lane=<lane>
```

An ID or lane is one nonempty ASCII token without whitespace. A lane list uses commas without spaces.

Each annotation is the last comment before its declaration. A blank line or another declaration breaks the attachment.

- `live-suite` attaches to a concrete shared-suite test function.
- `live-journey` attaches to one record in `sharedRuntimeScenarios()`.
- `live-fixture` attaches to the function that creates the semantic fixture variant.
- `live-proof` attaches to one runtime-specific live test function.

Use one annotation for one declaration. If one function creates multiple semantic fixture variants, extract a small wrapper.

Do not annotate non-gating experiment tests as proofs. Reconciliation joins their exact test names to the experiment section in the registry.

### Source binding rules

Add annotations only where the declared binding exists. Record a desired entry as `MISSING` when no valid declaration exists.

For a journey record, compare the annotation `id` with the adjacent `name` field. A mismatch is an invalid binding.

For a proof, compare the annotation `id` and `lane` with the registry entry. Then compare its test with the workflow selector.

For a fixture, compare the annotation `id` with every fixture reference in the registry. Keep its function name and path only in the one-off worksheet.

The union of all `live-suite` lane values describes implemented suite entry points. The workflow remains the authority for actual CI selection.

Do not convert standalone tests into shared journeys in this task. That work changes live behavior and needs separate approval.

### Registry reconciliation procedure

Add `## Registry reconciliation` to `docs/runtime-live-ci.md`. Name `docs/runtime-live-ci-registry.md` as the normative desired-state registry.

The procedure uses these steps:

1. Commit the candidate source and workflow changes.
2. Record that candidate commit as the reconciliation SHA.
3. List desired journey IDs, proof IDs and lanes, fixture IDs, and experiment test names.
4. List all four source annotation types under the two watched Go directories.
5. Compare each annotation with its adjacent declaration.
6. Join journey, proof, and fixture IDs in both directions.
7. Compare suite lanes and proof tests with the workflow selectors.
8. Inspect every test in a `//go:build live` file.
9. Classify each test as a journey, proof, experiment, or support test.
10. Inspect fixture builders that the live tests call.
11. Record bound, missing, duplicate, invalid, and orphan totals.
12. Update the reconciliation SHA in a second documentation-only commit.
13. Run the guard at the second commit.

The procedure records explicit ID lists for all nonzero diagnostic classes. A count without the IDs is not sufficient evidence.

### Diagnostics

Use these exact diagnostic classes:

- `MISSING journey <id>`: the registry ID has no valid annotated scenario record.
- `MISSING suite lane <lane>`: no suite annotation claims the required lane.
- `MISSING proof <id> lane=<lane>`: the proof has no valid annotated test.
- `MISSING fixture <id>`: the registry ID has no annotated builder.
- `DUPLICATE <kind> <id>`: more than one declaration claims a singular binding.
- `INVALID <file>:<line>`: the annotation grammar or adjacent declaration is wrong.
- `ORPHAN journey <id>`: an annotated journey is absent from the registry.
- `ORPHAN proof <id>`: an annotated proof is absent from the registry.
- `ORPHAN fixture <id>`: no registry journey, proof, or experiment references the annotated fixture.
- `UNACCOUNTED live test <name>`: a live-tagged test has no allowed classification.
- `UNACCOUNTED fixture builder <file>:<line>`: a live test calls an unbound semantic builder.
- `UNSELECTED proof <id> lane=<lane>`: the workflow does not invoke the annotated proof.

Missing and unselected diagnostics are desired-state gaps. Duplicate, invalid, orphan, and unaccounted diagnostics fail reconciliation.

### Standing SHA guard

Add the guard only to the secret-free `offline` job in `.github/workflows/runtime-live-e2e.yml`. Set `fetch-depth: 0` for that checkout.

The guard reads one 40-character SHA from `docs/runtime-live-ci.md`. It then runs this path-scoped comparison:

```sh
git diff --quiet "$reconciliation_sha"..HEAD -- \
  internal/ensigncycle/ \
  internal/livescenario/ \
  .github/workflows/runtime-live-e2e.yml
```

If the comparison fails, print the changed watched paths and stop the offline job.

If the SHA is absent or unavailable, stop with a specific error.

The watched paths are exactly:

- `internal/ensigncycle/`
- `internal/livescenario/`
- `.github/workflows/runtime-live-e2e.yml`

The guard does not parse Go, Markdown, or YAML semantics. It does not add a ratchet, generator, or second CI selector.

## Concrete documentation change

Before this task, `docs/runtime-live-ci.md` does not name the desired-state registry. It also has no reconciliation SHA.

After this task, the guide contains this normative statement:

> `docs/runtime-live-ci-registry.md` is the normative desired-state registry for runtime live CI. Source annotations bind its stable IDs to current declarations.

The new procedure also contains this record:

```text
Registry reconciliation SHA: `<40-character candidate SHA>`
```

The procedure states that source or selector changes need reconciliation in the same pull request. It states that the SHA guard detects stale review only.

## Implementation plan

1. Add annotations to the existing 10 shared journey records.
2. Add suite annotations to each concrete shared-suite function that the current workflow invokes.
3. Add proof annotations to the four registered runtime-specific tests.
4. Add fixture annotations to concrete semantic builders.
5. If an annotation needs one declaration, extract a small builder for inline or combined setup.
6. Run the one-off semantic inventory and record every current gap.
7. Add the reconciliation procedure and its first SHA placeholder.
8. Add the workflow guard and full-history checkout.
9. Commit the source, workflow, and procedure with the placeholder.
10. Replace the placeholder with that commit SHA in a documentation-only commit.
11. Run the guard with a stale SHA and make sure that it fails.
12. Run the guard with the recorded SHA and make sure that it passes.
13. Run the repository format and test commands.

## Expected surface

The baseline is 17 files and approximately 125 inserted lines. The tolerance is 20 files and 165 inserted lines.

| Surface | Expected insertion | Purpose |
|---|---:|---|
| `docs/runtime-live-ci.md` | 45 lines | Normative link, procedure, SHA, and inventory result |
| `.github/workflows/runtime-live-e2e.yml` | 18 lines | Full history plus the path-scoped guard |
| `internal/ensigncycle/shared_scenarios_test.go` | 10 lines | Existing journey bindings |
| `internal/ensigncycle/claude_live_runner_test.go` | 1 line | Claude suite binding |
| `internal/ensigncycle/codex_live_runner_test.go` | 1 line | Codex suite binding |
| Three runtime-proof test files | 4 lines | Proof and lane bindings |
| Eleven fixture source files | 46 lines | Fixture bindings and small setup wrappers |

The fixture files include `shared_fixtures_test.go`, `live_test.go`, `recorded_gate_lifecycle_test.go`, and `auto_continue_fixtures_test.go`.

They also include `dispatch_recovery_live_test.go`, `pi_live_runner_test.go`, and `shallow_boot_fixture_live_test.go`.

The final four are `zero_discover_live_test.go`, `haiku_loop_spike_live_test.go`, `auto_continue_pi_live_test.go`, and `internal/livescenario/ac2_reanchor.go`.

The tolerance permits these builder extractions. It does not permit new live journeys, new runner behavior, or selector expansion.

## Semantic budget

- Command grammar: no change.
- Stored formats: no change.
- Workflow authority: no change.
- Runtime behavior: no live-test behavior change.
- CI behavior: if a watched path is newer than the recorded reconciliation, the offline job stops.
- Source structure: comments and small fixture-builder extractions only.

## Acceptance criteria

**AC-1 (VALUE)** Every registry ID has one inventory status, and no source binding is ambiguous.
Measured by: the inventory reports all desired IDs as `BOUND` or `MISSING`.

It reports zero duplicate, invalid, orphan, and unaccounted results. It lists all unselected proofs as desired-state gaps.

**AC-2** The operating guide names the desired-state registry and records one usable reconciliation SHA.
Measured by: the SHA guard fails against a known stale base and passes against the recorded candidate base.

**AC-3** Every test in a live-tagged file and every semantic builder that it calls has an allowed classification.
Measured by: the inventory assigns each item to a journey, proof, experiment, or support class. Any unclassified item fails reconciliation.

**AC-4** Source paths and builder names do not become registry data.
Measured by: the registry diff contains no new Go path or builder symbol.

Moving an annotated builder without changing its ID does not require a registry edit.

## Test plan

- Run the one-off join on all annotations. Changing an ID or moving an annotation away from its declaration must produce `INVALID` or `MISSING`.
- Add a duplicate annotation in a temporary worktree. The inventory must produce `DUPLICATE`.
- Add an unreferenced fixture annotation in a temporary worktree. The inventory must produce `ORPHAN fixture`.
- Use a stale reconciliation SHA. The workflow guard must exit nonzero and print the changed watched path.
- Use the candidate reconciliation SHA. The workflow guard must exit zero.
- Change `docs/runtime-live-ci.md` only. The path guard must stay green because the guide is not a watched source path.
- Run `gofmt -w ./cmd ./internal` after the source comments and builder extractions.
- Run `go test ./...`.
- Run `go test ./... -race`.
- Run the zero-cost live definition guards with `-tags live`. Do not spend model credentials for this binding task.

## Mechanism decisions

The annotations serve AC-1. A documentation link alone cannot identify the declaration that owns a stable ID.

The one-off inventory serves AC-1 and AC-3. A permanent AST linter is larger than the approved mechanism.

The path-scoped SHA guard serves AC-2. A registry-only diff guard cannot detect source or selector drift.

Builder names in the registry would make moves change desired state. Stable fixture IDs keep those moves local to source.

An unwired-count ratchet would turn current gaps into accepted debt. Explicit `MISSING` diagnostics preserve the desired state without a ceiling.

## Stage Report: ideation

- DONE: Prove one annotation-to-registry join end to end. Show the guard fail and then pass.
  The `rejection-flow` join passed. Base `c25a4c1a` failed on three watched files, while base `2997da9a` passed.
- DONE: Define the minimal grammar, attachment points, procedure home, watched paths, and semantic reconciliation steps.
  The plan fixes four annotation forms, adjacent declarations, `docs/runtime-live-ci.md`, three watched paths, and a bidirectional join.
- DONE: Produce a complete plan with expected files, line estimates, acceptance checks, and orphan or missing-binding diagnostics.
  The plan sets a 17-file baseline, 125 insertions, four value criteria, negative cases, and 12 diagnostic classes.

### Summary

The ideation defines a manual semantic reconciliation plus one path-scoped SHA guard.

It preserves missing desired bindings as visible results. It avoids a permanent parser, generator, ratchet, or second selector.
