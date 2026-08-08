---
id: 0ytmjwn4ppg5en25z7vmna0p
title: Make Sonnet 5 the only Claude live lane on pull requests
status: implementation
source: Captain decision on 2026-08-07 after review of PR 626 and current Opus cost and failure evidence.
started: 2026-08-08T15:45:20Z
completed:
verdict:
score: 0.85
worktree: .worktrees/spacedock-ensign-make-sonnet-5-only-claude-pr-live-lane
issue:
pr:
mod-block:
sprint:
gates:
    version: 1
    records:
        - id: gate:0ytmjwn4ppg5en25z7vmna0p:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:0ytmjwn4ppg5en25z7vmna0p-backlog-1
              briefing:
                id: briefing:0ytmjwn4ppg5en25z7vmna0p:backlog:attempt-1:revision-1
                digest: sha256:dbf8b6eef823205ddcd3f98c3329756465a1102398f78b4a504f02ad3548023e
                request-digest: sha256:db28e83d3906eb06c9ceb9d98ffebe8ac28753cf4b020b902f23a306cdd061fb
                room-ref: ./make-sonnet-5-only-claude-pr-live-lane/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:0ytmjwn4ppg5en25z7vmna0p:backlog:1
                briefing: briefing:0ytmjwn4ppg5en25z7vmna0p:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-08T15:44:23.487623Z"
                decision: approve
                reason: Captain directed dispatch. Ideation must make routine PR CI select Sonnet 5/max and Codex Luna/max, retain Opus for pre-release, and move Pi to manual/local evidence.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:0ytmjwn4ppg5en25z7vmna0p:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:0ytmjwn4ppg5en25z7vmna0p-ideation-1
              briefing:
                id: briefing:0ytmjwn4ppg5en25z7vmna0p:ideation:attempt-1:revision-1
                digest: sha256:5fe1c5f4d43d0e30a78f008bc59b1cfed1d231f431d841972080f1fc05e2b25a
                request-digest: sha256:eeb9dfc192729f5e6f48544b0840636230c0cc1df3dc51b805c540d27ce21068
                room-ref: ./make-sonnet-5-only-claude-pr-live-lane/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:0ytmjwn4ppg5en25z7vmna0p:ideation:1
                briefing: briefing:0ytmjwn4ppg5en25z7vmna0p:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-08T15:57:40.803253Z"
                decision: approve
                reason: Captain approved the ideated two-lane pull-request cadence and directed implementation.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:0ytmjwn4ppg5en25z7vmna0p:validation
          stage: validation
          attempts:
            - id: gate-attempt:0ytmjwn4ppg5en25z7vmna0p-validation-1
              briefing:
                id: briefing:0ytmjwn4ppg5en25z7vmna0p:validation:attempt-1:revision-1
                digest: sha256:8084561597119202308347ff4ad9dedcd71501df568dcc3c4f053a1fec6fcd1d
                request-digest: sha256:f0eccaf219416c998fdad5bb18b4d49726a31e77c106466b7b456fe56bbcab24
                room-ref: ./make-sonnet-5-only-claude-pr-live-lane/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:0ytmjwn4ppg5en25z7vmna0p:validation:1
                briefing: briefing:0ytmjwn4ppg5en25z7vmna0p:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-08T18:32:10.784333Z"
                decision: approve
                reason: Validation reproduced all acceptance criteria and required repository checks on the integrated candidate.
              application:
                target-stage: done
                state: superseded
---

Make ordinary pull-request CI run one Claude lane: Sonnet 5 with maximum effort. Keep Opus as a pre-release lane.

## Outcome

A pull request queues two live approvals. Claude uses Sonnet 5 at maximum effort. Codex uses Luna at maximum effort.

An explicit pre-release dispatch uses Opus at maximum effort. Pi evidence uses a local subscription login and does not queue on pull requests.

## Problem

The current pull-request workflow queues four approval jobs. These jobs are Claude Sonnet, Claude Opus, Codex, and Pi.

Opus and Pi add approval delay and cost to each pull request. They do not add routine evidence that the pull request needs.

The current Claude matrix uses the `sonnet` alias and does not set effort. The Codex shim pins Luna but does not set effort.

PR 626 documented the four-job policy. The Opus matrix existed before that pull request.

## Proposed approach

Keep `.github/workflows/runtime-live-e2e.yml` as the only live workflow. Do not add a scheduler or a second policy file.

Add a `live_cadence` choice to `workflow_dispatch`. Use `sonnet` as its default and `opus-pre-release` as its other value.

Normalize a pull-request event to the `pull-request` cadence. Expand the existing Claude matrix from the cadence and the model list.

Use exclusions so that each event creates one Claude matrix leg:

| Event or dispatch choice | Claude model | Effort | Environment |
|---|---|---|---|
| `pull_request` | `claude-sonnet-5` | `max` | `CI-E2E` |
| `workflow_dispatch` with `live_cadence=sonnet` | `claude-sonnet-5` | `max` | `CI-E2E` |
| `workflow_dispatch` with `live_cadence=opus-pre-release` | `claude-opus-4-8` | `max` | `CI-E2E-OPUS` |

Install a small Claude command shim in the existing Claude job. The shim adds `--effort max` to the real Claude command.

Keep the explicit model argument in each live runner. This makes the selected model visible in the command and artifacts.

Extend the existing Codex shim with `-c 'model_reasoning_effort="max"'`. Keep `--model gpt-5.6-luna` on each `codex exec` form.

Remove the `pi-live` job from pull-request CI. Keep the Pi coverage tests in the offline job because they have no model cost.

Keep the documented local Pi smoke. It uses `pi login` and the operator's subscription, or an explicit local API key.

Update the operating guide and the site architecture note. These files must name each event, model, effort, approval, and Pi path.

### Smallest mechanism

The existing event, matrix, command shim, and offline test surfaces can express the full cadence. No new component or workflow is necessary.

The event expansion test serves AC-1, AC-2, and AC-7. A prose-only policy cannot fail on a bad matrix change.

The two command-shim tests serve AC-2 and AC-3. Environment prose cannot prove the final command arguments.

Moving the Pi coverage guards serves AC-4. Deleting all Pi checks would remove free registry-parity evidence.

### Spike result

The installed Claude CLI is version `2.1.220`. `claude --effort max --version` exited 0.

`claude --model claude-sonnet-5 --effort max --help` exited 0. This no-cost command proves the model and effort spellings.

`codex -c 'model_reasoning_effort="max"' debug models` exited 0. Its Luna record lists `max` as a supported effort.

The implementation still needs the fixture expansion test. That test is the first implementation step and does not run a paid model.

## Expected surface

The baseline is 8 files, about 190 insertions, and about 275 deletions.

| File | Expected change |
|---|---|
| `.github/workflows/runtime-live-e2e.yml` | Add the cadence input, event exclusions, and maximum-effort shims. Remove the Pi job. About +45/-220. |
| `internal/release/runtime_live_cadence_workflow_test.go` | Add event fixtures and mutation cases for approval-job expansion. About +105/-0. |
| `internal/ensigncycle/codex_liveenv_test.go` | Make sure that all three `codex exec` forms select Luna and maximum effort once. About +12/-8. |
| `internal/release/runtime_live_evidence_workflow_test.go` | Remove paid Pi job claims. Keep the free Pi controls. About +4/-10. |
| `internal/release/workflow_exec_guard_test.go` | Move Pi parity requirements to the offline job and reject a live Pi job. About +8/-18. |
| `internal/release/cilog_clean_output_workflow_test.go` | Remove deleted Pi live steps and require two live `gotestsum` installs. About +4/-8. |
| `docs/runtime-live-ci.md` | Replace the four-job cadence with pull-request, pre-release, and local Pi paths. About +10/-8. |
| `docs/site/contributing/architecture-notes.md` | Apply the same short cadence summary for site readers. About +2/-3. |

Tolerance is 10 files, 240 insertions, and 330 deletions. A larger change requires a new ideation review.

Files can move within `internal/release` when existing test helpers make that move smaller. The semantic limits do not change.

### Semantic changes

- Command grammar changes: `workflow_dispatch` gains the `live_cadence` choice.
- Runtime behavior changes: pull requests queue two live approval jobs instead of four.
- Runtime behavior changes: Claude uses `claude-sonnet-5` and maximum effort on pull requests.
- Runtime behavior changes: Codex uses `gpt-5.6-luna` and maximum effort on pull requests.
- Runtime behavior changes: explicit pre-release dispatches select Opus and maximum effort.
- Runtime behavior changes: Pi live evidence runs locally with subscription authentication.
- Stored formats do not change.
- State authority does not change.
- The desired journey registry does not change.

### Operating-guide diff

Apply this replacement in `docs/runtime-live-ci.md`:

```diff
- The offline gate job (`go test ./...`, no secrets) must pass before either live lane burns its environment approval.
+ The offline gate job (`go test ./...`, no secrets) must pass before a live lane uses an environment approval.
- `claude-live` runs the core, shared, merged, bare, and break-glass proofs. Its matrix uses `sonnet` and `claude-opus-4-8`.
- `codex-live` runs the resolver and shared proofs. The PR delta and release ledger consume its metrics.
- `pi-live` runs the coverage guards and one front-door smoke. It uploads the grade, root session, child session, and diagnostics.
+ Pull requests run `claude-sonnet-5` at maximum effort and `gpt-5.6-luna` at maximum effort.
+ An explicit `live_cadence=opus-pre-release` dispatch runs `claude-opus-4-8` at maximum effort.
+ Pi live evidence runs locally with `pi login`. The offline job keeps the Pi coverage guards.
```

Apply this replacement in `docs/site/contributing/architecture-notes.md`:

```diff
- **`claude-live`** (matrix `sonnet` and `claude-opus-4-8`): secret `ANTHROPIC_API_KEY`. Runs the full-cycle smoke and the shared suite, loading the current checkout via `spacedock claude --plugin-dir "$GITHUB_WORKSPACE"`.
- **`codex-live`**: secret `OPENAI_API_KEY`. Builds a local marketplace under `$RUNNER_TEMP` and fails if the listing names a remote `github.com`/`ref next` install instead of the local path.
- **`pi-live`**: installs `pi-coding-agent` and runs the Pi coverage guard plus the front-door smoke.
+ **Pull requests** run Claude Sonnet 5 at maximum effort and Codex Luna at maximum effort.
+ **Pre-release dispatches** run Opus at maximum effort when `live_cadence=opus-pre-release`.
+ **Pi evidence** uses the local subscription path. Pull requests keep only the free Pi coverage guards.
```

## Out of scope

- Keep the desired journey registry and `docs/runtime-live-ci-registry.md` unchanged.
- Keep all journey, fixture, assertion, coverage, and target-specific TODO ownership unchanged. This includes the Opus gap owned by `a7`.
- Use only `.github/workflows/runtime-live-e2e.yml`. Do not add a workflow, cadence package, or reconciliation layer.
- Do not require a paid Opus run for this change.

## Acceptance criteria

**AC-1 - Each pull request queues exactly two live approval jobs.**

Proof: the event fixture expands to Claude Sonnet 5 at maximum effort and Codex Luna at maximum effort.

The same fixture expands to zero Opus jobs and zero Pi jobs. The old baseline expands to four approval jobs.

**AC-2 - Ordinary pull requests run exactly one Claude live lane.**

Proof: the pull-request fixture produces `claude-sonnet-5`, `max`, and `CI-E2E`. Mutations fail for aliases, Opus, extra legs, or lower effort.

**AC-3 - Ordinary pull requests run Codex Luna at maximum effort.**

Proof: the existing shim test exercises all three front-door forms. Each final command contains Luna and maximum effort exactly once.

**AC-4 - Opus and Pi evidence remain available without pull-request approval waits.**

Proof: the pre-release fixture produces one Opus maximum-effort leg. The local Pi command runs with subscription authentication.

The offline suite runs the Pi coverage guards. No pull-request fixture produces a `CI-E2E-OPUS` or `CI-E2E-PI` job.

**AC-5 - The desired journey and target registry remains byte-identical.**

Proof: `git diff --exit-code -- docs/runtime-live-ci-registry.md` exits 0 after implementation.

**AC-6 - Operators can identify all cadences without reading workflow YAML.**

Proof: the guide names the pull-request lanes, the Opus dispatch value, and the local Pi subscription command.

**AC-7 - The change uses the existing workflow and test surface.**

Proof: the diff adds no workflow, package, desired-journey row, or reconciliation command. The diff stays within the declared tolerance.

## Test plan

1. Add table fixtures for `pull_request`, manual Sonnet, and `opus-pre-release`.

2. Expand the workflow matrix for each fixture. Make sure that each fixture has the expected model, effort, environment, and approval count.

3. Add mutations for Opus on pull requests, Pi on pull requests, duplicate Claude legs, aliases, and lower effort.

4. Exercise the Claude shim with `--version` and `--help`. Exercise the Codex shim with its three supported command prefixes.

5. Run the focused `internal/release` and `internal/ensigncycle` tests. These tests use fixtures and do not spend model credits.

6. Run `go test ./...` and `go test ./... -race`.

7. Run `git diff --exit-code -- docs/runtime-live-ci-registry.md`.

A paid Sonnet, Codex, Opus, or Pi run is not necessary. The next normal pull request supplies Sonnet and Codex evidence.

## Stage Report: ideation

- DONE: Define the complete cadence: Sonnet 5/max and Codex Luna/max on pull requests, Opus at pre-release, and Pi through manual/local subscription evidence.
  The event table and semantic-change list define each model, effort, trigger, environment, and authentication path.
- DONE: Produce the smallest workflow and operating-guide design, with expected files, insertions, deletions, semantic changes, and tolerance.
  The design reuses one workflow and declares an 8-file baseline, a 10-file tolerance, and exact guide replacements.
- DONE: Demonstrate event selection with a cheap falsifiable expansion or fixture, without adding a second enforcement layer or requiring a paid Opus run.
  The plan defines three event fixtures and adversarial mutations. No-cost CLI probes exited 0 for both maximum-effort settings.

### Summary

The plan reduces pull-request approval jobs from four to two. Opus remains an explicit pre-release run.

Pi evidence moves to the local subscription path.

## Stage Report: implementation

- DONE: Make pull-request event expansion queue exactly Sonnet 5/max and Codex Luna/max, while Opus remains pre-release and Pi remains local/manual.
  Commit `105ffea63` adds three event fixtures; alias, Opus, duplicate-leg, and lower-effort mutations fail the pull-request assertion.
- DONE: Implement the cadence through the existing workflow and shims within the approved 10-file, +240/-330 ceiling, with the desired journey registry unchanged.
  Commit `105ffea63` changes 9 files by +240/-291; `git diff --exit-code -- docs/runtime-live-ci-registry.md` exited 0 before commit.
- DONE: Add falsifiable event and command-argument coverage, then run focused checks, gofmt, the full suite, and the race suite.
  Focused release and ensigncycle tests pass; shim tests fail if maximum effort or the explicit model is removed or duplicated. Both full suites reached one unrelated `internal/gates` fixture failure because two shared-state manifest files are absent.

### Summary

Pull requests now queue Sonnet 5 at maximum effort and Codex Luna at maximum effort. Opus is available only through the pre-release dispatch, while Pi live evidence stays local and its free coverage guards run offline.

## Stage Report: implementation (cycle 2)

- DONE: Make pull-request event expansion queue exactly Sonnet 5/max and Codex Luna/max, while Opus remains pre-release and Pi remains local/manual.
  Merge commit `8728da3a0` keeps upstream's `^TestLiveCommon` selectors for Claude and Codex, requires Pi's selector in the local guide, and rejects any Pi workflow job.
- DONE: Implement the cadence through the existing workflow and shims within the approved 10-file, +240/-330 ceiling, with the desired journey registry unchanged.
  The current-main delta is 11 files and +226/-438; the captain authorized the added manifest-coupling removal and reconciliation file after the original ceiling. The registry diff against `origin/main` is empty.
- DONE: Add falsifiable event and command-argument coverage, then run focused checks, gofmt, the full suite, and the race suite.
  Event mutations fail on Opus, aliases, duplicate legs, or lower effort. Claude and Codex shim tests require each model and maximum-effort argument exactly once; focused tests, `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` pass.
- DONE: Remove the mutable real-state dependency from the pilot manifest test while preserving synthetic gate format coverage.
  The deleted test failed when normal archival moved two tasks. `TestRecordClosureShapesApplication` still exercises application parsing and exact canonical keys with a synthetic workflow.

### Summary

The integrated branch now carries the approved two-lane pull-request cadence on current `origin/main`. It also removes the invalid test dependency on the mutable shared state checkout and keeps upstream's common-journey reconciliation aligned with local Pi evidence.

## Stage Report: validation

- DONE: Verify pull requests select exactly Sonnet 5/max and Codex Luna/max, while Opus remains pre-release and Pi remains local/manual.
  `TestRuntimeLiveWorkflowExpandsApprovedCadences` observed one Sonnet 5/max PR leg and two total approvals; exact shim-argv tests observed Luna/max, while mutations reject Opus, aliases, duplicates, lower effort, and a paid Pi job.
- DONE: Verify the mutable docs/dev entity-state pilot dependency is removed and focused synthetic gate-format coverage remains sufficient.
  The 11-file delta deletes the fixed 31-path pilot manifest and its shared-state reader; `TestRecordClosureShapesApplication` exercised synthetic approve/revise/hold records and exact canonical application keys.
- DONE: Verify the current-main 11-file delta, unchanged desired journey registry, formatting, focused tests, full suite, and race suite with independent evidence.
  Against merge parent `8728da3a0^2`, the candidate is 11 files and +226/-438; registry diff, `git diff --check`, and `gofmt -d ./cmd ./internal` were clean, and focused, full, and race suites exited 0.
- DONE: AC-1 - Each pull request queues exactly two live approval jobs.
  Cadence expansion returned exactly Sonnet 5/max in `CI-E2E` plus Codex in `CI-E2E-CODEX`; adding a Pi job or exposing the PR Opus leg makes the focused guards fail.
- DONE: AC-2 - Ordinary pull requests run exactly one Claude live lane.
  The parsed PR matrix produced one `claude-sonnet-5`/`max`/`CI-E2E` leg; alias, duplicate-model, Opus-exclusion, and effort mutations were all rejected.
- DONE: AC-3 - Ordinary pull requests run Codex Luna at maximum effort.
  `TestCodexLiveWorkflowPinsOnlyExecToLuna` executed all three supported shim prefixes and observed exactly one Luna model flag and one maximum-effort setting in each final argv.
- DONE: AC-4 - Opus and Pi evidence remain available without pull-request approval waits.
  Expansion observed one Opus/max `CI-E2E-OPUS` pre-release leg; registry reconciliation observed the Pi common command exactly once in the local guide and absent from the workflow.
- DONE: AC-5 - The desired journey and target registry remains byte-identical.
  `git diff --exit-code 8728da3a0^2..8728da3a0 -- docs/runtime-live-ci-registry.md` exited 0; reconciliation independently observed all 16 desired and actual journeys.
- DONE: AC-6 - Operators can identify all cadences without reading workflow YAML.
  The reconciled guide names Sonnet 5/max, Luna/max, `opus-pre-release`, and the local `pi login` subscription path; the site architecture note carries the same three-path summary.
- DONE: AC-7 - The change uses the existing workflow and test surface.
  The upstream-relative delta changes one existing workflow and 10 existing code/doc paths, adds no workflow or package, and preserves the desired registry under the captain-authorized 11-file scope.
- DONE: Recommend PASSED with no material findings or deferred risks.
  The semantic adversarial pass covered all three cadences, duplicate and wrong-identity legs, lower effort, non-exec shim inputs, paid Pi reintroduction, canonical gate bytes, and full/race execution.

### Summary

Validation independently reproduced every acceptance criterion at merge commit `8728da3a0`. The candidate is clean, all focused and repository-wide checks pass, and no real shared-state entity inventory remains coupled to the gate tests; recommendation: PASSED.

## Stage Report: implementation (cycle 3)

- DONE: Reproduce and diagnose the rejected pull-request lane expansion before changing the workflow.
  PR #639 run `31272368367` at exact head `8728da3a0` queued Sonnet job `93140486631`, Codex job `93140486613`, and unexpected Opus job `93140486627`. GitHub applies `exclude` before `include`, so the Opus include created a new include-only row after its base combination was excluded; the local simulator was false-green because it modeled include augmentation before exclusion.
- DONE: Replace the Cartesian Claude matrix with one explicit include row selected directly from the event and dispatch cadence.
  Local commit `8e2d67f76` contains no base axes and no exclusions. Its only row selects pull-request/Sonnet 5/`CI-E2E` for pull requests or manual Sonnet, Opus 4.8/`CI-E2E-OPUS` only for `opus-pre-release`, and maximum effort in every case.
- DONE: Guard the exact one-row mechanism against the failure mode and structural regressions.
  `TestRuntimeLiveWorkflowHasOneExplicitClaudeCadence` requires exactly one include row with the approved expressions and rejects a second row, every base axis, an exclusion, aliases, reduced effort, and cadence or environment mapping changes. The correction changes 2 files by +44/-66 relative to `8728da3a0`.
- DONE: Run all required formatting, focused, full, and race checks without spending live-model credits.
  `gofmt -w ./cmd ./internal`, focused release, ensigncycle, and contractlint tests, `go test ./...`, and `go test ./... -race` all pass. The desired journey registry remains unchanged.
- DONE: Leave delivery and renewed live proof to the first officer.
  The implementation commit is local only: this worker did not push the code branch or mutate PR #639. No paid jobs were approved; the next synchronized PR run remains the direct semantic proof.

### Summary

The live failure exposed a mismatch between GitHub matrix semantics and the prior test simulator. The corrected workflow has one explicit Claude matrix row, making an ordinary pull request structurally incapable of queuing a separate Opus row while preserving the explicit pre-release Opus dispatch.

## Stage Report: validation (cycle 2)

- DONE: Verify commit 8e2d67f76 defines exactly one Claude matrix row with no base axes or exclusions, and preserves the three intended cadence mappings.
  The workflow has one `include` row, no base axes or `exclude`, and maps PR/manual Sonnet to Sonnet 5/max/CI-E2E and `opus-pre-release` to Opus 4.8/max/CI-E2E-OPUS.
- FAILED: Verify the replacement test checks the structural one-row contract without simulating GitHub matrix expansion.
  A throwaway mutation added `os: [ubuntu-latest, macos-latest]`; `TestRuntimeLiveWorkflowHasOneExplicitClaudeCadence` still exited 0 because its typed YAML shape ignores unknown matrix axes.
- DONE: Verify the live finding.
  `gh run view 31272368367` observed PR SHA `8728da3a0`, queued Sonnet job `93140486631`, Codex job `93140486613`, and unexpected Opus job `93140486627`; the run is cancelled and all paid jobs have empty step lists.
- DONE: Verify the two-file correction.
  `8728da3a0..8e2d67f76` changes only `.github/workflows/runtime-live-e2e.yml` and `internal/release/runtime_live_evidence_workflow_test.go`, by +44/-66; `git diff --check` is clean.
- DONE: Verify the unchanged registry.
  `git diff --exit-code 8728da3a0..8e2d67f76 -- docs/runtime-live-ci-registry.md` exited 0.
- DONE: Verify focused checks.
  Release, ensigncycle, and contractlint packages pass; the cadence test rejects the named axes, exclusions, multiple include rows, aliases, lower effort, and mapping changes, but not arbitrary axes.
- DONE: Verify the full suite and race suite.
  `go test ./... -race` passed; an isolated `go test ./...` rerun passed after a concurrent full/race run exposed a transient split-root dispatch-golden mismatch.
- DONE: AC-1 - Each pull request queues exactly two live approval jobs.
  Current YAML structurally defines one Claude include row plus the unchanged Codex job; direct GitHub job-graph proof remains pending.
- FAILED: AC-2 - Ordinary pull requests run exactly one Claude live lane.
  The candidate YAML meets the value, but its required extra-leg regression proof is incomplete because an arbitrary two-value base axis escapes the structural test.
- DONE: AC-3 - Ordinary pull requests run Codex Luna at maximum effort.
  The correction does not change the Codex job or its exact-argv tests; the focused ensigncycle suite passes.
- DONE: AC-4 - Opus and Pi evidence remain available without pull-request approval waits.
  The single row selects Opus only for `opus-pre-release`; no Pi job exists, and the correction leaves offline/local Pi evidence unchanged.
- DONE: AC-5 - The desired journey and target registry remains byte-identical.
  The exact correction-range registry diff exited 0.
- DONE: AC-6 - Operators can identify all cadences without reading workflow YAML.
  The correction changes no operator documentation, preserving the previously validated cadence descriptions.
- DONE: AC-7 - The change uses the existing workflow and test surface.
  The correction touches two existing files and adds no workflow, package, registry row, reconciliation command, or simulator.
- FAILED: Recommend whether it is safe to update PR #639 for direct job-graph proof.
  Not yet: fix the material evidence defect first, then push that exact corrected SHA so the replacement PR run proves the final candidate and no paid lane needs approval.
- FAILED: Recommend PASSED or REJECTED.
  Recommend REJECTED for a material evidence defect: supported trigger is an added matrix axis; harm is duplicate paid Claude approvals; authority is `value-ac[AC-2]` exactly one Claude lane; trigger evidence is the passing two-value `os` mutation.

### Summary

The two-file candidate YAML itself is a clean one-row correction, and all focused, isolated full, and race checks pass. Validation rejects because the replacement structural test can still pass when a new Cartesian base axis creates multiple Claude rows; repair that guard before pushing the final exact SHA to PR #639 for direct job-graph proof.

### Feedback Cycles

- Cycle 1: REJECTED — PR #639 live job graph; surface 2 files/110 LOC vs estimate 8 files/570 LOC (19%); AC unchanged.
- Cycle 2: REJECTED — fresh validation of the one-row guard; surface 2 files/110 LOC vs estimate 8 files/570 LOC (19%); AC unchanged.
- Cycle 3: REJECTED — corrected PR live run exposed concrete-model/stable-role drift; AC unchanged; captain authorized narrow rework.

## Stage Report: implementation (cycle 4)

- DONE: Make the raw YAML matrix key set exactly `include`, so any unknown base axis fails the structural guard.
  Commit `113738b20` parses the matrix as a raw key map, requires `include` to be its sole key, and only then decodes and checks the approved row.
- DONE: Keep the one-row workflow and acceptance criteria unchanged; add no matrix simulator or controller.
  The correction changes only the existing structural test by +20/-12; the workflow YAML and desired journey registry are byte-identical to `8e2d67f76`.
- DONE: Reproduce the arbitrary-axis negative control and run focused, formatting, full, and race checks before reporting.
  Adding `os: [ubuntu-latest, macos-latest]` first reproduced the false green, then failed under the corrected guard; focused release/ensigncycle/contractlint tests, `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` pass.

### Summary

The structural proof now observes every matrix key instead of enumerating known axes, so any new Cartesian axis invalidates the test. The one-row workflow is unchanged, PR #639 remains untouched, and the local candidate is ready for fresh independent validation.

## Stage Report: validation (cycle 3)

- DONE: Re-run the arbitrary unknown-axis mutation and verify the raw matrix key set must be exactly `include`.
  In a detached audit checkout, adding `os: [ubuntu-latest, macos-latest]` made the focused test exit 1 with `Claude matrix keys = [include os], want only include`.
- DONE: Verify the workflow remains one explicit row and the correction adds no simulator, controller, or semantic scope.
  The unchanged workflow has only one `include` row; commit `113738b20` replaces the typed matrix fields with a raw key map and adds no expansion logic or runtime component.
- DONE: Verify the one-file correction.
  `8e2d67f76..113738b20` changes only `internal/release/runtime_live_evidence_workflow_test.go`, by +20/-12; `git diff --check` and `gofmt -d ./cmd ./internal` are clean.
- DONE: Verify the unchanged workflow and registry.
  Exact-range `git diff --exit-code` checks for `.github/workflows/runtime-live-e2e.yml` and `docs/runtime-live-ci-registry.md` both exited 0.
- DONE: Verify focused checks.
  Release, ensigncycle, and contractlint packages pass; the cadence test now fails for every extra matrix key, any second include row, wrong identity or effort, and changed cadence/environment mappings.
- DONE: Verify the full suite and race suite.
  Sequential isolated runs of `go test ./...` and `go test ./... -race` both exited 0.
- DONE: AC-1 - Each pull request queues exactly two live approval jobs.
  The raw-key test requires one Claude include row and the unchanged Codex approval lane; direct GitHub job-graph observation remains the post-push proof.
- DONE: AC-2 - Ordinary pull requests run exactly one Claude live lane.
  The structural guard requires the entire matrix key set to equal `[include]` and that include to contain exactly one Sonnet 5/max/CI-E2E row for pull requests.
- DONE: AC-3 - Ordinary pull requests run Codex Luna at maximum effort.
  The correction leaves the Codex workflow and exact-argv tests unchanged, and the focused ensigncycle package passes.
- DONE: AC-4 - Opus and Pi evidence remain available without pull-request approval waits.
  The unchanged single row selects Opus only for `opus-pre-release`, no Pi job exists, and offline/local Pi evidence remains intact.
- DONE: AC-5 - The desired journey and target registry remains byte-identical.
  The correction-range registry diff exited 0.
- DONE: AC-6 - Operators can identify all cadences without reading workflow YAML.
  No operator documentation changed, preserving the previously validated cadence descriptions.
- DONE: AC-7 - The change uses the existing workflow and test surface.
  The correction modifies one existing test file and adds no workflow, package, registry row, simulator, controller, or reconciliation command.
- DONE: Recommend whether exact SHA 113738b20 is safe to push to PR #639 for direct job-graph proof.
  Safe to push exactly `113738b20`; PR #639 currently remains at `8728da3a0`, and no paid lane should be approved while the replacement graph is inspected.
- DONE: Recommend PASSED or REJECTED.
  Recommend PASSED for the local correction with no material findings, deferred risks, or polish findings; GitHub behavior has not been claimed or observed for `113738b20`.

### Summary

Validation reproduced the former false-green mutation and observed the corrected raw-key guard reject it. The one-file test correction is clean, every requested suite passes, and exact SHA `113738b20` is safe for the first officer to push to PR #639 for direct job-graph proof without approving paid lanes.
