---
title: The lane model ids are written down five times, and the Pi lane runs an older Luna
status: implementation
score: 0.8
source: "Captain, 2026-10-04: pin gpt-6-luna for CI, and move the Claude lane to sonnet-5.5. Local smoke on openai-codex/gpt-6-luna passed (root and child, 166s, grade artifact written)."
id: q67rfhn33tncbvf9xqbq960r
started: 2026-10-04T05:48:52Z
worktree: .worktrees/spacedock-ensign-pin-lane-models-in-one-place
---

## Problem

Two problems with one cause.

Each lane model id is written down in several places. The Pi lane's
`openai-codex/gpt-5.6-luna:max` appears in `internal/ensigncycle/pi_liveenv.go`, in
three branches of the codex shim inside `.github/workflows/runtime-live-e2e.yml`, in
its step summary, in `docs/runtime-live-ci.md`, and in fixtures. The Claude lane's
`claude-sonnet-5` appears in the workflow matrix, in the harness's accepted-value
list, in the context-budget table, and in an artifact name.

And the Pi lane runs an older Luna. A local front-door smoke on
`openai-codex/gpt-6-luna:max` passed today: 5 reads, ensign skill at read 2, root and
child both `openai-codex/gpt-6-luna`, 166s, grade artifact written, no model error.

## Design

**One place for the lane models.** Put the ids in `internal/release` or
`internal/ensigncycle` as named constants, and print them for the workflow the way
`stamp-paths` prints the stamp targets: a small command that emits the model id per
lane and per auth path. The workflow matrix, the codex shim, the step summary, the
docs and the fixtures then consume that output instead of repeating literals. No
model id may remain in the workflow file.

**Pin the Pi lane to gpt-6-luna**, both auth paths: the OAuth spelling used with the
subscription, `openai-codex/gpt-6-luna:max`, and the API-key spelling beside it. Verify
the API-key spelling against the provider before pinning it; do not assume it from the
OAuth one.

**Claude moves only when the id is confirmed.** `claude-sonnet-5.5` is refused by the
installed CLI today: "There's an issue with the selected model (claude-sonnet-5.5). It
may not exist or you may not have access to it." Keep `claude-sonnet-5` until an
exact id is confirmed, and add the new value in the same single source when it is.

## Out of scope

- Recorded fixtures that name a past model on purpose. `auto_continue_replay_test.go`
  references real `claude-sonnet-5` streams and the CI artifact they came from; those
  are evidence, not configuration, and renaming them falsifies the record.
- Auth, quota, and thinking level.

## Expected surface and tolerance

`internal/ensigncycle/pi_liveenv.go`, the new print command, the workflow, the docs,
and the fixtures that assert the shim's argv. Estimate net +70, across 6 files.
Tolerance net +30 to +140, at most 8 files.

## Acceptance criteria

**AC-1 (VALUE) - One place holds each lane model.**
Verified by: the printed model list is what the workflow and the harness use, and a
test fails when they diverge. Falsifier: change the constant and require the workflow
site to follow.
**AC-2 - The Pi lane runs gpt-6-luna on both auth paths.**
Verified by: the harness constants and the shim's argv name it; a local front-door
smoke passes with root and child on it. Falsifier: leave the old literal and the shim
`--model` argument must differ from the constant.
**AC-3 - No model literal remains in the workflow.**
Verified by: searching the workflow for each id and requiring no match outside the
command that prints them. Falsifier: add one and the check must fail.
**AC-4 (no-regression) -** Real Pi discovery is unaffected, the Claude lane keeps its
current model until the id is confirmed, and the offline suite passes.

## Test plan

Behavioural tests only, with the model supplied independently. Every test must state
the reason it could fail. No test may assert a file's own text against a copy of it,
and no expected value may come from the thing it tests. The smoke is a local run with
the captain's credentials, never a CI lane run.

## Stage Report: implementation

- DONE: Put the lane model ids in one Go source and print them for the workflow, the way stamp-paths prints the stamp targets. Then make the workflow matrix, the codex shim, the step summary, the docs and the fixtures consume that output, so no model literal remains in .github/workflows/runtime-live-e2e.yml.
  `internal/release/livemodels.go` holds the five ids and `LiveModels()`; `spacedock live-models` prints `lane.auth=id` (`--get <key>` prints one) and gated `internal/release/livemodels_test.go` proves no literal remains and every resolved key exists. `grep -nE 'claude-sonnet-5|claude-opus-4-8|gpt-6-luna|gpt-5.6-luna' .github/workflows/runtime-live-e2e.yml` = no match; the offline job resolves the five ids into `needs.offline.outputs.*`, consumed by the claude matrix, the Codex shim (`SPACEDOCK_LIVE_CODEX_MODEL`), the Pi summary, and the journey-delta artifact name. Commit fbb6cc8df.
- DONE: Pin the Pi lane to gpt-6-luna on both auth paths. The OAuth spelling openai-codex/gpt-6-luna:max passed a local front-door smoke today with root and child on it. Verify the API-key spelling against the provider before pinning it; do not assume it from the OAuth one.
  `release.PiOAuthModel`/`PiAPIKeyModel` = `openai-codex/gpt-6-luna:max` / `openai/gpt-6-luna:max`; `pi_liveenv.go` aliases them. API-key spelling verified against the installed pi-ai catalog: `providers/data/openai.json` declares `gpt-6-luna` with `thinkingLevelMap.max="max"` (and `openai-codex.json` the same), so `openai/gpt-6-luna:max` is real, not assumed.
- DONE: Keep the Claude lane on claude-sonnet-5. The id claude-sonnet-5.5 is refused by the installed CLI today, so record it as an input the captain must confirm and do not pin it.
  `release.ClaudeSonnetModel = "claude-sonnet-5"` (opus stays `claude-opus-4-8`); `docs/runtime-live-ci.md` records `claude-sonnet-5.5` as refused by the installed Claude CLI and pending captain confirmation. Not pinned.
- DONE: Do not rename the recorded fixtures that reference real claude-sonnet-5 streams; they are evidence of past runs.
  `auto_continue_replay_test.go`, `shared_filing_pr399_test.go`, `journey_workflow_test.go`, `dispatch/contextbudget_parity_test.go`, `journeymetrics/claude_test.go`, `shallow_boot_measure_unit_test.go` untouched (`git show --stat fbb6cc8df` lists none of them).
- DONE: Prove it: the focused tests, the print command's output, and one local front-door smoke on the pinned Pi model with a distinct and EMPTY SPACEDOCK_LIVE_ARTIFACT_DIR. State the reason each test you add could fail.
  Print output: `claude.sonnet=claude-sonnet-5`, `claude.opus=claude-opus-4-8`, `codex.exec=gpt-6-luna`, `pi.oauth=openai-codex/gpt-6-luna:max`, `pi.api-key=openai/gpt-6-luna:max`. Live smoke is PENDING (see Summary). Falsifiers, each run:
  - `TestLiveModelsCommandPrintsPinnedLaneModels` (cli): fails if a lane constant changes without the print — verified by reverting `PiOAuthModel` to `gpt-5.6-luna`, which failed the exact-stdout assert.
  - `TestLiveModelsCommandGetPrintsOneID` (cli): fails if a key is renamed or an unknown key exits 0/prints empty.
  - `TestRuntimeLiveWorkflowCarriesNoLaneModelLiteral` (release): fails if a literal is re-added — verified by injecting `claude-sonnet-5` into the workflow, which failed.
  - `TestRuntimeLiveWorkflowResolvedKeysExist` (release): fails if a workflow `--get` key has no `LiveModels()` entry — verified by renaming `--get claude.sonnet` to `claude.sonnetX`, which failed.
  - `TestRuntimeLiveWorkflowLaneModelWiring` (release): fails if an offline output is renamed on one side only — verified by suffixing the `pi_api_key` consumer, which failed.
  - `TestCodexLiveWorkflowPinsOnlyExecToLuna` (updated): supplies `SPACEDOCK_LIVE_CODEX_MODEL=luna-sentinel-model`; fails if the shim hardcodes a model — verified by replacing `$live_model` with `hardcoded-luna`, which failed.
  - `TestRuntimeLiveClaudeShimSetsMaximumEffort` (updated): uses `release.ClaudeSonnetModel` only as the arbitrary passthrough arg; fails if the Claude shim stops prepending `--effort max`.

### Summary

Pinned the lane model ids in one Go source (`internal/release/livemodels.go`), added `spacedock live-models` to print them, and rewired the workflow matrix, the Codex exec shim, the Pi step summary, the docs, and the harness/fixtures to consume that source — the live workflow now carries no model literal. The Pi lane runs `gpt-6-luna` on both auth paths; the API-key spelling is verified against the installed provider catalog, and the Claude lane stays on `claude-sonnet-5` with `claude-sonnet-5.5` recorded as pending captain confirmation.

**Live proof is PENDING, not passed.** A local front-door smoke (`TestLivePiFrontDoorSmoke`, a distinct empty `SPACEDOCK_LIVE_ARTIFACT_DIR="/tmp/spacedock-pin-lane-smoke-1791093534"`) finished `ok` in 193.660s and its grade artifact read `verdict: pass` with `root_session.model = openai-codex/gpt-6-luna` (179506ms) and `child_session.model = openai-codex/gpt-6-luna` (77719ms), 5 reads and ensign skill at read 2. The FO's mid-run steering says the machine was saturated (1-min load ~118 earlier, ~85 at steering) and that a live result taken now must be reported PENDING; the smoke had already completed when that arrived, so it was not killed and its numbers are not treated as a quiet-machine conclusion. Load at stop: `8.48 21.01 21.19` on 10 cores. The captain should re-run the smoke on a quiet machine.

Deviation: 13 files, net +289 (design estimate +70, tolerance +30..+140, at most 8 files). The extra surface is forced by the mandated behavioral proof (new print command + constants + two new test files) and by the fixtures the Problem names; no scope was added beyond the design's named sites.
