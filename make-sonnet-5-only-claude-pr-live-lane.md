---
id: 0ytmjwn4ppg5en25z7vmna0p
title: Make Sonnet 5 the only Claude live lane on pull requests
status: ideation
source: Captain decision on 2026-08-07 after review of PR 626 and current Opus cost and failure evidence.
started: 2026-08-08T15:45:20Z
completed:
verdict:
score: 0.85
worktree:
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

- Do not change common journeys, fixtures, assertions, or target-specific TODO bindings.
- Do not change the Codex or Pi journey coverage.
- Do not repair the Opus product gap owned by `a7`.
- Do not remove Opus from the desired-state registry.
- Do not change `docs/runtime-live-ci-registry.md`.
- Do not add a new workflow.
- Do not add a cadence package or a reconciliation layer.
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
