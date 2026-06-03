---
id: f4h107he3s9swwz51ymr022m
title: Enable Codex live CI by reusing shared runtime scenarios on top of r0
status: validation
source: "captain (2026-06-02) — revive the useful old Python Codex live CI setup, ignore tautological static grep tests, and document setup/principles in docs/dev/README.md"
started: 2026-06-02T17:59:45Z
completed:
verdict:
score: "0.72"
worktree: .worktrees/spacedock-ensign-codex-live-ci
issue:
mod-block:
---

The old Python-era live CI had the right Codex shape: a `codex-live` job with `CI-E2E-CODEX`, `OPENAI_API_KEY`, Codex CLI install, `codex login --with-api-key`, preserved artifacts, and a `make test-live-codex` target that reused scenario files also run by Claude. The useful overlap was three shared scenarios:

- `test_gate_guardrail.py --runtime codex`
- `test_rejection_flow.py --runtime codex`
- `test_merge_hook_guardrail.py --runtime codex`

The static workflow-grep tests around that harness are not worth reviving.

This task ports the useful live behavior to the current Go repo as a shared runtime scenario library. Codex CI should exercise the same common user journeys Claude uses where the behavior is runtime-neutral, and f4 must be stacked on r0 so the live run proves the Codex runtime adapter and Codex-shaped dispatch contract.

## Problem

Current CI has a Claude live lane but no Codex live lane. r0 adds the Codex runtime contract and Codex-shaped dispatch text, but Codex needs proof beyond local CLI/plugin probing and unit tests. The proof must cover runtime regressions such as skill loading, `codex exec --json` event shape, multi-agent notification handling, install/login drift in GitHub Actions, and the r0 adapter actually being used by a real Codex first officer.

The first f4 implementation was too narrow: a Codex-only gate-hold smoke proved the live test library and local auth plumbing, but it did not reuse the old shared scenarios and it was not stacked on r0. That makes it a useful harness spike, not a sufficient CI deliverable.

The old Python live harness proves the intended scenario overlap, but it cannot be copied wholesale: `tests/` is gone, the repo is Go-first, the old harness used synthetic `$HOME/.agents/skills/spacedock` skill injection, and some old assertions (especially "immediate wait after dispatch") contradict the r0 Codex contract where mailbox notification is the completion signal and `wait_agent` is only a sparing critical-path accelerator.

## Proposed approach

Rework the Go live harness under `internal/ensigncycle` so common journeys are scenario functions with runtime-specific drivers. The Codex suite must include at least the three old shared Codex/Claude overlaps:

- gate guardrail: drive a gated workflow to the gate and assert the FO waits for approval instead of self-approving or archiving;
- rejection flow: exercise validation rejection and feedback routing back to implementation, using Codex notification/reuse semantics from r0;
- merge hook guardrail: exercise the PR merge hook guard, ensuring terminalization cannot bypass the registered merge hook/mod-block flow.

The harness should launch real Codex headlessly via `codex exec --json` (or a proven headless equivalent) with isolated `HOME` and `CODEX_HOME`, load the current checkout's Spacedock plugin/skills, preserve JSONL/transcript artifacts when CI sets a test-artifact root, and allow local subscription auth by copying only `~/.codex/auth.json` into the isolated `CODEX_HOME` while keeping CI strict on `OPENAI_API_KEY`.

Stack f4 on `spacedock-ensign/codex-runtime-adapter` before validation. The `codex-live` job in `.github/workflows/runtime-live-e2e.yml` should depend on the offline Go suite, use the approval-gated `CI-E2E-CODEX` environment, read only `OPENAI_API_KEY`, install/log in to Codex, build the local `spacedock` binary, prove Codex plugin readiness against the current checkout, run the shared Codex scenario suite, and upload artifacts.

Update `docs/dev/README.md` (the evergreen development doc referenced by `AGENTS.md`) with the Codex live test setup and principles: live behavior is the proof; static grep tests are not a substitute; shared runtime scenarios are preferred over one-off host-only smokes; the job must test the PR/current checkout, not `next`; serial cheap preflight before heavier cases; environment approval protects spend/secrets; artifacts are part of the debugging contract.

## Riskiest unknown

The current-checkout Codex plugin load has already been proven by the first f4 implementation: a generated local marketplace can point `spacedock@spacedock` at the PR checkout and `codex plugin list` resolves that temp local path. The remaining riskiest unknown is whether the old shared Codex scenarios can be expressed in the current Go harness without recreating stale Python semantics. The implementation should port the scenario intent, not line-for-line Python code.

Record any scenario that cannot be faithfully ported and why. A port that reduces the suite back to a single gate smoke does not satisfy this task.

If a full scenario is too expensive for every CI run, it may be split into a cheap required live preflight plus a separately gated extended Codex journey, but the acceptance criteria must still include the three shared journeys and validation must run them locally or in approved CI before passing.

## Out of scope

- Do not revive the old Python static workflow-grep tests.
- Do not port stale Python syntax or runtime internals line-for-line.
- Do not reintroduce `test_codex_dispatch_immediate_wait.py` semantics; immediate wait-after-spawn is stale after r0.
- Do not add a model/cost matrix yet. Cost ledger integration belongs to the separate journey-cost-ledger task.

## Acceptance criteria

**AC-1 — f4 is stacked on r0 and therefore exercises the Codex runtime adapter.**
Verified by: `git merge-base --is-ancestor origin/spacedock-ensign/codex-runtime-adapter HEAD` passes from `.worktrees/spacedock-ensign-codex-live-ci`, and the live Codex artifacts show `skills/first-officer/references/codex-first-officer-runtime.md` was available from the current checkout plugin cache.

**AC-2 — The Go live harness has shared runtime scenarios for the old Claude/Codex overlap.**
Verified by: the code exposes and tests scenario definitions for gate guardrail, rejection flow, and merge hook guardrail, and the Codex live suite runs those same scenario definitions rather than a Codex-only one-off test.

**AC-3 — A real Codex runtime run executes the shared scenarios against the current checkout.**
Verified by: `OPENAI_API_KEY=... go test -tags live -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v` (or the implemented equivalent) launches real Codex headlessly, runs the gate guardrail, rejection flow, and merge hook guardrail scenarios, and fails on self-approval, missing rejection feedback routing, or merge-hook bypass. Without live Codex auth, the test may skip locally; in CI the secret-bearing job must fail clearly when `OPENAI_API_KEY` is absent after environment approval.

**AC-4 — GitHub Actions has an approval-gated `codex-live` lane that runs the shared scenario suite against the current checkout.**
Verified by: `.github/workflows/runtime-live-e2e.yml` includes a `codex-live` job that needs the offline Go gate, uses environment `CI-E2E-CODEX`, reads `OPENAI_API_KEY` but not `ANTHROPIC_API_KEY`, installs/logs in to Codex, builds the local `spacedock` binary, loads or installs the current checkout's plugin/skills, runs the shared Codex scenario suite, and uploads artifacts. Validation should cite a successful Actions run when available; local validation must at least run the same job commands through the live shared suite on a machine with Codex auth.

**AC-5 — Codex plugin readiness is proven with real Codex CLI commands, not assumed from docs.**
Verified by: the live job or its setup script runs the current Codex plugin install/load path and `go test ./internal/cli -run TestCodexResolveManifestAgainstInstalledHost -v` (or a stronger current-checkout equivalent) after installation. The evidence must show the resolver sees `spacedock@spacedock` from the intended checkout/cache; a run that only tests remote `next` does not satisfy this AC.

**AC-6 — Evergreen docs explain Codex live setup and testing principles in `docs/dev/README.md`.**
Verified by: `docs/dev/README.md` contains a Codex live CI/testing section covering required local tools/secrets, GitHub environment/secret names, how to run the local shared suite, why static grep tests do not prove runtime behavior, why shared runtime scenarios are preferred over one-off host-only smokes, why CI must test the current checkout rather than `next`, and how artifacts are preserved for debugging.

**AC-7 — The implementation does not restore tautological static tests as the primary proof.**
Verified by: validation evidence centers on live command execution and resulting state/output. Syntax or parser checks are allowed as supporting guardrails only; they cannot be the only proof for AC-1 through AC-5.

## Test plan

- Cheap setup check: run `codex --version`, `codex login status` or equivalent, and a current-checkout plugin install/load spike before writing the live test.
- Go gates: `go test ./...`, `go test ./... -race`, and focused tests for any new live harness helpers.
- Stack check: from the f4 worktree, `git merge-base --is-ancestor origin/spacedock-ensign/codex-runtime-adapter HEAD`.
- Codex CLI smoke: `go test ./internal/cli -run TestCodexResolveManifestAgainstInstalledHost -v` after plugin setup.
- Shared live runtime suite: run the new `go test -tags live -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v` locally when Codex auth is available; otherwise document the skip and rely on the approval-gated CI run.
- CI evidence: after the PR opens, approve `CI-E2E-CODEX` and cite the `codex-live` job result and uploaded artifacts in validation.

## Stage Report: ideation

- DONE: Scoped the task to live Codex behavior, not static workflow grep tests
  The task revives the old Python-era `codex-live` setup pattern but explicitly rejects the tautological static assertions as primary proof.
- DONE: Named the current-checkout plugin-load risk as the first spike
  CI must prove the PR checkout's plugin/skills are loaded; testing `spacedock-dev/spacedock --ref next` would false-pass against already-published code.
- DONE: Added the evergreen docs requirement to the ACs
  `docs/dev/README.md`, which `AGENTS.md` references as the development workflow definition, is the target documentation surface.

### Summary

Ready for implementation: add a Go live Codex gate-hold smoke, wire an approval-gated `codex-live` Actions job around it, prove current-checkout plugin loading, and document the setup/principles in `docs/dev/README.md`.

## Stage Report: implementation

- DONE: Current-checkout Codex plugin load is proven and cannot false-pass against next
  Generated `.agents/plugins/marketplace.json` with `source: local`; isolated `codex plugin list` showed a temp local path, resolver passed, and the cache contained new `codex_live_test.go`.
- DONE: Live Go Codex gate-hold smoke and codex-live CI lane are implemented
  Code commit `11381e15397336626144d7ee75964fd6d65c5afd` adds `TestLiveCodexGateGuardrail`; `.github/workflows/runtime-live-e2e.yml` adds `codex-live` on `CI-E2E-CODEX`.
- DONE: docs/dev/README.md documents Codex live setup and proof principles
  `docs/dev/README.md` now covers local setup, `CI-E2E-CODEX`/`OPENAI_API_KEY`, current-checkout marketplace loading, live proof, and artifacts.
- SKIPPED: Local live Codex execution with real OPENAI_API_KEY
  No `OPENAI_API_KEY` is present locally; `go test -count=1 -tags live -run TestLiveCodexGateGuardrail ./internal/ensigncycle -v` compiled and skipped, and the CI-required missing-secret path failed clearly.

### Summary

Implemented a build-tagged Codex live gate-hold smoke, helper tests for local marketplace/auth/gate assertions, and a `codex-live` Actions lane that installs the current checkout as a local Codex marketplace before running the live smoke. Verified with uncached `go test -count=1 ./...`, uncached `go test -count=1 ./... -race`, focused Codex helper tests, YAML parsing, and an isolated Codex plugin install proving the cache came from the current checkout rather than remote `next`.

## Stage Report: validation

- DONE: AC-1 - Go live Codex smoke proves FO gate-hold behavior through real Codex runtime
  `go test -count=1 -tags live -run TestLiveCodexGateGuardrail ./internal/ensigncycle -v` compiled and skipped without `OPENAI_API_KEY`; required CI mode failed clearly with the approval-gated key message.
- DONE: AC-2 - GitHub Actions has approval-gated `codex-live` lane against current checkout
  `.github/workflows/runtime-live-e2e.yml` has `codex-live` behind `CI-E2E-CODEX`, needs `offline`, reads only `OPENAI_API_KEY`, installs/logs in to Codex, builds local `spacedock`, installs the current checkout marketplace, runs resolver/live tests, and uploads artifacts.
- DONE: AC-3 - Codex plugin readiness is proven with real Codex CLI commands/current-checkout mechanism, not remote `next`
  Isolated temp `HOME`/`CODEX_HOME` ran `codex plugin marketplace add`, `codex plugin add spacedock@spacedock`, `codex plugin list`, and `go test -count=1 ./internal/cli -run TestCodexResolveManifestAgainstInstalledHost -v`; listing pointed at the temp local checkout path.
- DONE: AC-4 - `docs/dev/README.md` documents Codex live setup and proof principles
  The Codex Live CI section covers local prerequisites, `CI-E2E-CODEX`/`OPENAI_API_KEY`, current-checkout marketplace loading, runtime proof over static grep, live smoke command, and artifacts.
- DONE: AC-5 - Primary proof is live command execution/resulting state/output, not tautological static tests
  Evidence is the live-tagged `codex exec --json` harness, state/output assertions, real Codex plugin commands, and Go gates; no stale Python/static workflow-grep proof was revived.
- DONE: Local gates are rerun from the code worktree, including focused helper tests, `go test ./...`, and `go test ./... -race`
  From `.worktrees/spacedock-ensign-codex-live-ci`: `go test -count=1 ./internal/ensigncycle`, `go test -count=1 ./...`, and `go test -count=1 ./... -race` all passed.
- DONE: The validation report recommends PASSED or REJECTED and flags any remaining CI/live-auth caveats
  Recommendation: PASSED; caveat: no local `OPENAI_API_KEY` or approved Actions run was available, so full live-auth runtime execution remains to be observed in `codex-live`.

### Summary

Validation recommends PASSED for commit `11381e15397336626144d7ee75964fd6d65c5afd`. The implementation provides a real Codex CLI/current-checkout setup path plus a live-tagged FO gate-hold smoke; local validation proved the non-auth skip and CI-required missing-secret failure paths, but did not observe a paid live Codex run.

## Stage Report: implementation follow-up

- DONE: Restored the old Python harness's local subscription-auth ergonomics for the Go Codex live smoke
  Code commit `e4f3b568` adds a local Codex auth path: when `OPENAI_API_KEY` is absent and `~/.codex/auth.json` exists, the test copies only `auth.json` into a temporary `CODEX_HOME`, preserving isolated plugin marketplace state while using the operator's existing Codex login.
- DONE: Kept CI approval-gated auth strict
  `SPACEDOCK_CODEX_LIVE_REQUIRED=1` still requires `OPENAI_API_KEY`; local subscription auth is not accepted as a CI substitute.
- DONE: Proved the live smoke locally through real Codex without `OPENAI_API_KEY`
  `env -u OPENAI_API_KEY SPACEDOCK_BIN=/tmp/spacedock-codex-live SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/spacedock-codex-live-artifacts go test -count=1 -tags live -run TestLiveCodexGateGuardrail ./internal/ensigncycle -v` passed in 24.4s.
- DONE: Reran the normal Go gates after the follow-up
  `go test ./...` passed 747 tests in 12 packages, and `go test ./... -race` passed 747 tests in 12 packages.

### Summary

Follow-up validation is required because the previous validation report is now stale: the local live Codex smoke no longer skips when the operator has an existing Codex subscription login.

## Stage Report: validation follow-up

- DONE: Re-check the entity acceptance criteria against the full implementation plus the follow-up commit.
  AC-1 through AC-5 remain satisfied at code commit `e4f3b568ff123192a89f04ad6313e790e099c70d`; the checks below prove behavior outside the task body.
- DONE: AC-1 - Go live Codex smoke proves FO gate-hold behavior through a real Codex runtime.
  Built `/tmp/spacedock-codex-live-validation` from the code worktree, then `env -u OPENAI_API_KEY ... go test -count=1 -tags live -run TestLiveCodexGateGuardrail ./internal/ensigncycle -v` passed in 40.93s.
- DONE: AC-2 - GitHub Actions has an approval-gated `codex-live` lane that runs the live smoke against the current checkout.
  `.github/workflows/runtime-live-e2e.yml` still gates `codex-live` on `CI-E2E-CODEX`, requires only `OPENAI_API_KEY`, installs/logs in to Codex, builds local `spacedock`, installs the current checkout plugin, runs resolver/live tests, and uploads artifacts.
- DONE: AC-3 - Codex plugin readiness is proven with real Codex CLI commands, not assumed from docs.
  Live artifacts show `codex plugin list` installed `spacedock@spacedock` from the temp local marketplace path, and `go test -count=1 ./internal/cli -run TestCodexResolveManifestAgainstInstalledHost -v` passed.
- DONE: AC-4 - Evergreen docs explain Codex live setup and testing principles in `docs/dev/README.md`.
  The Codex Live CI docs now describe `codex login`, local `~/.codex/auth.json` copying into temp `CODEX_HOME`, `OPENAI_API_KEY`, CI strictness, current-checkout plugin loading, live proof, and artifacts.
- DONE: AC-5 - The implementation does not restore tautological static tests as the primary proof.
  Validation evidence is real `codex exec --json` execution, temp workflow state assertions, Codex plugin commands, focused helper tests, and Go gates; no static workflow-grep proof was revived.
- DONE: Specifically verify the old Python harness issue: local Codex subscription auth should work in an isolated home without `OPENAI_API_KEY`, by copying only `~/.codex/auth.json` into temp `CODEX_HOME`; CI-required mode must still require `OPENAI_API_KEY`.
  Local `~/.codex/auth.json` was nonempty; the live smoke passed with `OPENAI_API_KEY` unset and `codex login status` reported ChatGPT auth, while `SPACEDOCK_CODEX_LIVE_REQUIRED=1` failed clearly without `OPENAI_API_KEY`.
- DONE: Run the smallest command set that proves the follow-up.
  Focused auth-helper tests passed, `go test -count=1 ./...` passed 747 tests in 12 packages, and `go test -count=1 ./... -race` passed 747 tests in 12 packages.
- SKIPPED: Approved GitHub Actions `codex-live` run.
  No approved Actions run was available in this local validation; the local live command exercised the same current-checkout plugin and live Codex gate-smoke path.
- SKIPPED: `gofmt -w ./cmd ./internal`.
  This worker cannot edit product code; read-only `gofmt -l ./cmd ./internal` reported unrelated `internal/status/enum_scope_test.go`, while the codex-live touched Go files were gofmt-clean.

### Summary

Recommendation: PASSED for commit `e4f3b568ff123192a89f04ad6313e790e099c70d`. The follow-up restores local subscription-auth ergonomics without weakening CI: local isolated `CODEX_HOME` auth passed a real live Codex smoke with no `OPENAI_API_KEY`, and CI-required mode still hard-fails without `OPENAI_API_KEY`.

## Stage Report: PR prep

- DONE: Rebased the code branch onto current `origin/next`
  The branch now carries two commits: `96dd51bc` (`ci: add codex live gate smoke`) and `5e4a15f7` (`test: allow codex live smoke to use local auth`).
- DONE: Reran post-rebase verification
  `gofmt -l ./cmd ./internal` was clean, `go test -count=1 ./...` passed 745 tests in 12 packages, and `go test -count=1 ./... -race` passed 745 tests in 12 packages.
- DONE: Reran the local subscription-auth live Codex smoke post-rebase
  Built `/tmp/spacedock-codex-live-pr`, then `env -u OPENAI_API_KEY SPACEDOCK_BIN=/tmp/spacedock-codex-live-pr SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/spacedock-codex-live-pr-artifacts go test -count=1 -tags live -run TestLiveCodexGateGuardrail ./internal/ensigncycle -v` passed in 33.19s.

### Summary

The PR branch is rebased and ready to push after captain approval of the draft body.

## Stage Report: scope reframe

- DONE: Pulled f4 back from stale PR-prep
  Cleared `mod-block: merge:pr-merge` and moved f4 from `validation` back to `implementation` because the previous one-off Codex gate smoke was only a harness spike.
- DONE: Reframed the deliverable around existing shared runtime scenarios
  The task body now requires the old Python-era Claude/Codex overlap: gate guardrail, rejection flow, and merge hook guardrail, ported as shared Go scenario definitions rather than Codex-only smoke logic.
- DONE: Made r0 stacking an explicit acceptance criterion
  AC-1 now requires f4 to be based on `origin/spacedock-ensign/codex-runtime-adapter` so the live Codex run exercises r0's runtime adapter files and Codex-shaped dispatch contract.

### Summary

The prior f4 validation evidence is stale. The next implementation pass must stack f4 on r0 and turn the Codex live lane into a shared scenario suite that reuses the existing Claude journey shape.

## Stage Report: implementation follow-up (cycle 2)

- DONE: Reviewed and salvaged the staged WIP from the crashed worker
  Code commit `1aeb7a96` keeps the current-checkout marketplace/local-auth harness and replaces the one-off Codex smoke target with shared scenario execution.
- DONE: Completed coherent shared scenario definitions under `internal/ensigncycle`
  `codexSharedScenarios()` now defines `gate-guardrail`, `rejection-flow`, and `merge-hook-guardrail` with old Python provenance, behavior intent, and live timeouts.
- DONE: Ensured the live Codex suite is named `TestLiveCodexSharedScenarios` and runs all three shared scenarios
  Final live run passed all subtests: gate guardrail 40.56s, rejection flow 170.55s, merge-hook guardrail 32.53s.
- DONE: Preserved current-checkout Codex plugin marketplace loading and local auth behavior
  The suite still installs `spacedock@spacedock` from a generated local marketplace, checks the r0 Codex runtime adapter in the plugin cache, and copies only `~/.codex/auth.json` into isolated `CODEX_HOME`.
- DONE: Updated CI and evergreen docs for the shared suite
  `.github/workflows/runtime-live-e2e.yml` now runs `TestLiveCodexSharedScenarios`; `docs/dev/README.md` documents the three scenarios, local auth, current-checkout plugin loading, and why static grep/host-only smoke is insufficient.
- DONE: Ran required verification on the final rebased code branch
  `go test ./internal/ensigncycle`, `go test ./...`, and `go test ./... -race` passed after dropping the unrelated `internal/status` churn.
- DONE: Checked the required formatting command after removing unrelated status churn
  `gofmt -w ./cmd ./internal` exits 0 but dirties only `internal/status/fence_helpers_test.go` and `internal/status/mutate.go`; those upstream-only formatting diffs were restored so f4 stays clean.
- SKIPPED: Approved GitHub Actions `codex-live` run
  No push or CI approval was performed because this dispatch explicitly said not to push either branch.

### Summary

The f4 code branch is now rebased on current `origin/next` and carries the shared Codex live suite at commit `1aeb7a96` without unrelated `internal/status` churn. Local live verification used existing subscription auth with `OPENAI_API_KEY` unset and proved all three required shared scenarios through real `codex exec --json`; the follow-up cleanup reran focused/full Go gates and documented the upstream gofmt dirtiness separately.

## Stage Report: implementation (cycle 3)

- DONE: Rebase f4 onto current origin/next without unrelated changes
  Rebasing onto `origin/next` succeeded cleanly; f4 now carries `27c1ba5b`, `4c04a346`, and `407b82ef` on top of `1dfac9a5`, and the code worktree is clean.
- DONE: Preserve shared Codex live scenario suite and Phase 0.A contract changes
  `TestLiveCodexSharedScenarios`, the three shared scenarios, CI/docs wiring, and current-checkout plugin checks remain present; `git merge-base --is-ancestor 528ffa17 HEAD` passed after the named r0 remote branch had been pruned.
- DONE: Rerun required gates and live shared suite when local auth allows
  `gofmt -w ./cmd ./internal` was run, unrelated current-next formatting drift was restored, `go test ./internal/ensigncycle`, `go test ./...`, `go test ./... -race`, the resolver smoke, and `TestLiveCodexSharedScenarios` all passed.
- DONE: Local live shared-suite artifacts prove current-checkout Codex setup
  Local ChatGPT auth was available; the live run passed all three scenarios in 253.86s and artifacts show the generated marketplace plugin plus r0 Codex FO adapter in the plugin cache.
- SKIPPED: Approved GitHub Actions `codex-live` run
  No approved Actions environment run was available from this local worker; local live artifacts show `spacedock@spacedock` installed from a generated current-checkout marketplace and the r0 Codex FO adapter present in the plugin cache.

### Summary

The f4 code branch is rebased onto current `origin/next` with no new code edits beyond the rebase, now ending at `407b82ef`. Local verification proved the shared Codex live suite through real `codex exec --json` using isolated local ChatGPT auth, while preserving the Phase 0.A/r0 contract content and keeping unrelated formatting drift out of the branch.

## Stage Report: validation (cycle 2)

- DONE: AC-1 — f4 is stacked on r0 and therefore exercises the Codex runtime adapter.
  The exact `origin/spacedock-ensign/codex-runtime-adapter` ref is deleted upstream, but `git merge-base --is-ancestor 528ffa17 HEAD` passed and live artifacts recorded `skills/first-officer/references/codex-first-officer-runtime.md` in the current-checkout plugin cache.
- DONE: AC-2 — The Go live harness has shared runtime scenarios for the old Claude/Codex overlap.
  `go test -count=1 ./internal/ensigncycle` passed and `TestCodexSharedScenarioDefinitions` requires `gate-guardrail`, `rejection-flow`, and `merge-hook-guardrail` with old Python provenance.
- DONE: AC-3 — A real Codex runtime run executes the shared scenarios against the current checkout.
  `env -u OPENAI_API_KEY ... go test -count=1 -tags live -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v` passed all three scenarios in 280.62s using isolated local ChatGPT auth.
- DONE: AC-4 — GitHub Actions has an approval-gated `codex-live` lane that runs the shared scenario suite against the current checkout.
  `.github/workflows/runtime-live-e2e.yml` gates `codex-live` on `CI-E2E-CODEX`, uses only `OPENAI_API_KEY`, installs/logs in to Codex, builds local `spacedock`, installs the current checkout plugin, runs resolver/live tests, and uploads artifacts.
- DONE: AC-5 — Codex plugin readiness is proven with real Codex CLI commands, not assumed from docs.
  The live setup ran `codex plugin marketplace add`, `codex plugin add spacedock@spacedock`, `codex plugin list`, and the resolver smoke passed; artifacts show `spacedock@spacedock` installed from a generated local marketplace path, not remote `next`.
- DONE: AC-6 — Evergreen docs explain Codex live setup and testing principles in `docs/dev/README.md`.
  The Codex Live CI section covers local/CI auth, `CI-E2E-CODEX`, `OPENAI_API_KEY`, current-checkout marketplace loading, shared scenarios, live proof over static grep, and artifact preservation.
- DONE: AC-7 — The implementation does not restore tautological static tests as the primary proof.
  Evidence centers on real `codex exec --json`, resulting workflow state/output assertions, Codex plugin commands, focused Go tests, full Go gates, and race gates; no Python static workflow-grep proof was revived.
- DONE: Run Go gates and required formatter command.
  `gofmt -w ./cmd ./internal` ran, then gofmt-only drift in unrelated current-next files was restored; `go test -count=1 ./...` and `go test -count=1 ./... -race` each passed 833 tests in 12 packages.
- DONE: Run detached adversarial audit for the CI/live-runtime surface.
  In a throwaway checkout, deleting `merge-hook-guardrail` from `codexSharedScenarios()` made `TestCodexSharedScenarioDefinitions` fail with `[gate-guardrail rejection-flow], want ...`; no material test-strength hole found.
- SKIPPED: Approved GitHub Actions `codex-live` run.
  No approved Actions environment run was available locally; validation used the same live shared suite and current-checkout plugin path, but CI artifact evidence is still pending.

### Summary

Recommendation: PASSED for code commit `407b82efb836`. Caveats: the named r0 remote branch has been pruned so AC-1 used the r0 commit `528ffa17` ancestor check, the branch is now behind latest `origin/next` by three release-fix commits after fetch, and approved GitHub `codex-live` evidence is not yet available.
