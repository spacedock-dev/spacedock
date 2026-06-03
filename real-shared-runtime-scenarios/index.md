---
id: 8y7yten220npj6g4kj4680p2
title: Extract real shared runtime scenarios for Claude and Codex live CI
status: implementation
source: "captain (2026-06-03) - follow-up from f4 codex-live-ci: make the shared scenario reuse real, not Codex-only"
started: 2026-06-03T07:09:59Z
completed:
verdict:
score: "0.68"
worktree: .worktrees/spacedock-ensign-real-shared-runtime-scenarios
issue:
mod-block:
---

f4 ports the old Python Claude/Codex journey overlap into the Codex live lane, but the Go implementation is still Codex-shaped: `codexSharedScenarios()` defines the scenario list and the current Claude live CI continues to run `TestLiveEnsignCycle` rather than consuming the same scenario table.

This task makes the sharing real. The scenario definitions should be host-neutral, with Claude and Codex runners implementing the same scenario IDs through their runtime-specific launch, auth, plugin, artifact, and transcript mechanisms.

## Problem

Runtime regressions should be caught once per user journey and then exercised by each supported host. Today the old shared journey intent is preserved for Codex, but Claude and Codex do not share a single Go scenario contract. That allows parity drift: a scenario can exist for Codex only, the Claude lane can keep testing a different full-cycle shape, and CI can appear broad while the actual common journey set is split.

## Proposed approach

Extract the f4 scenario list into a host-neutral table, such as `sharedRuntimeScenarios()`, that contains only runtime-neutral facts: scenario ID, old Python provenance, behavior intent, fixture/journey shape, timeout/cost class, and expected durable outcome.

Move host-specific details behind runner adapters:

- Codex runner: `codex exec --json`, isolated `CODEX_HOME`, local Codex marketplace/plugin install, Codex artifacts.
- Claude runner: `spacedock claude`, isolated Claude config/auth, team/session artifacts, current checkout plugin path.

Both runners should implement the same scenario IDs. Meta tests should fail if either runtime is missing a runner for a shared scenario. Assertions should prefer durable workflow state and output over transcript phrasing: frontmatter, archive/no-archive state, stage reports, exact fix markers, merge-hook refusal, and only durable user-facing final-message obligations such as a gate review and decision prompt.

## Out of scope

- Do not add a large new scenario matrix in this task.
- Do not add token/cost ledger serialization; that belongs to journey-cost-ledger work.
- Do not remove the existing Claude full-cycle smoke unless the shared scenario suite fully replaces its coverage and validation proves the replacement.

## Design refinement (host-neutral table + runner-adapter boundary)

f4 (PR #278) already factored the scenario surface into four layers; three of them are
ALREADY host-neutral and only the fourth is Codex-shaped. The rename is small because
the seam already exists:

1. **Scenario table** — `codexSharedScenarios()` returns `[]codexSharedScenario{name,
   oldPythonTest, intent, timeout}`. Every field is runtime-neutral already. Rename the
   type to `sharedRuntimeScenario` and the constructor to `sharedRuntimeScenarios()`;
   `oldPythonTest` keeps the `--runtime codex` provenance string but the table carries no
   launch/auth/plugin field, so a host-neutral table is a rename plus a meta-test that
   the struct exposes no Claude-only/Codex-only field.
2. **Fixtures** — `writeCodexGateWorkflow` / `codexGateReadme` / `codexGateEntity` (and the
   rejection + merge-hook equivalents) write plain spacedock workflow files (README,
   entity, `_mods/`). Nothing in them is host-specific — the spike reused them verbatim to
   drive a Claude run. Drop the `Codex` infix at implementation; they are shared fixtures.
3. **Prompts** — `codexGatePrompt()` et al. already say `Use $spacedock:first-officer`,
   which both hosts honor. The spike fed `codexGatePrompt()` unchanged to `spacedock
   claude` and the Claude FO behaved correctly. Shared, modulo the rename.
4. **Assertions** — `assertCodexGateHeld(before, after, observed)`,
   `assertCodexRejectionFlow(entity, observed)`, `assertCodexMergeHookGuardHeld(before,
   after, observed)` consume only entity-state strings + an `observed` output string.
   They are already host-neutral functions; the spike ran `assertCodexGateHeld` UNCHANGED
   on Claude-produced state+output and it passed.

The host-specific surface is exactly ONE thing: the **runner adapter** that turns a shared
scenario into a real launch and returns `(before, after, observed)`:

| Concern            | Codex runner (`codexLiveRunner`)                          | Claude runner (new)                                              |
|--------------------|-----------------------------------------------------------|-----------------------------------------------------------------|
| Auth / HOME isol.  | isolated `CODEX_HOME` + copied `auth.json` / `OPENAI_API_KEY` | clean `HOME` + OAuth benchmark-token / `ANTHROPIC_API_KEY` (`isolatedClaudeEnv`) |
| Plugin install     | local Codex marketplace symlink + `codex plugin add`      | `spacedock claude --plugin-dir <checkout> --skip-contract-check` |
| Launch             | `codex exec --json --output-last-message <file>`          | `spacedock claude -- -p <prompt> --output-format stream-json`    |
| `observed` extract | read `--output-last-message` file (+ jsonl)               | extract the `result`/`success` event's `result` text from the stream (front-door analog of `--output-last-message`) |
| Artifacts          | jsonl / final-message / stderr                            | session jsonl / stream transcript                               |

The shared coverage meta-test (AC-2/AC-3) iterates `sharedRuntimeScenarios()` and fails
if EITHER host's runner map lacks a runner for a shared scenario ID — that is the parity
guard against drift.

## Spike result (riskiest unknown, exercised first)

**Riskiest unknown:** the Codex runner gets a clean final-message from
`codex exec --output-last-message`; the Claude front door is stream-watched and does NOT
extract any final-message string today (`streamEntry` doesn't even parse the `result`
event). So: can a Claude runner produce an `observed` string the SHARED assertions accept,
against the SAME fixture — i.e. is the sharing real or only Codex-shaped?

**Exercise (throwaway):** staged the existing `gate-guardrail` fixture
(`codexGateReadme`/`codexGateEntity`), drove it through the real
`spacedock claude --plugin-dir <checkout> --skip-contract-check -- -p <codexGatePrompt()>
--output-format stream-json --model sonnet` front door, extracted the final message from
the stream's `{"type":"result","subtype":"success","result":"…"}` event, and ran the
UNCHANGED `assertCodexGateHeld(before, after, finalMessage)`.

**Result: PASS.** The Claude FO held the gate — entity unmutated, still `status: review`,
no `completed`/`verdict` set, `gate-check.md` not archived — and its final message carried
both `Gate review:` and `Decision:`. The host-neutral assertion accepted Claude-produced
state+output verbatim. This proves (a) the `result`/`success` event is the front-door
analog of Codex `--output-last-message` and is the correct Claude `observed` source, and
(b) the table/fixtures/prompts/assertions are genuinely host-neutral, not Codex-shaped.

**Seed for the implementation's first test:** the Claude runner's `observed` extractor must
prefer the `result` event's `result` field and fall back to the last assistant `text`
block; the spike's `extractClaudeFinalMessage` is the shape to port. One known machine
gotcha surfaced and is recorded so it is not silent: `isolatedClaudeEnv` reads the OAuth
benchmark-token, which can be expired — a stale token yields an `is_error:true,
api_error_status:401` `result` event and a non-zero exit BEFORE any FO work, so the Claude
runner must fail loudly on a 401 result (distinct from a gate-assertion failure) rather
than feeding the 401 text into the assertion.

## Acceptance criteria

**AC-1 - Shared runtime scenarios are defined once in host-neutral code.**
Verified by: focused Go tests under `internal/ensigncycle` assert the shared scenario table includes at least `gate-guardrail`, `rejection-flow`, and `merge-hook-guardrail`, and that the table type does not encode Claude-only or Codex-only runner fields.

**AC-2 - Codex live tests consume the shared table.**
Verified by: `go test -tags live -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v` iterates the host-neutral table and fails if any shared scenario lacks a Codex runner.

**AC-3 - Claude live tests consume the same shared table.**
Verified by: `go test -tags live -run TestLiveClaudeSharedScenarios ./internal/ensigncycle -v` or the implemented equivalent iterates the same host-neutral table and fails if any shared scenario lacks a Claude runner.

**AC-4 - CI proves both runtime lanes run the shared scenario suite.**
Verified by: `.github/workflows/runtime-live-e2e.yml` runs the shared Claude scenario suite in `claude-live` and the shared Codex scenario suite in `codex-live`, with environment-gated secrets and artifact upload preserved.

**AC-5 - Scenario assertions are behavior/state oriented, not transcript-shape tautologies.**
Verified by: unit or fixture tests exercise at least one negative case per shared scenario where a broken state transition, missing rejection route, or merge-hook bypass makes the assertion fail.

**AC-6 - Evergreen docs explain the shared-scenario contract.**
Verified by: `docs/dev/README.md` documents how to add a shared runtime scenario, what belongs in the host-neutral definition, what belongs in each runner, and the commands for local Claude/Codex live execution.

**AC-7 - The Claude runner's final-message extraction is the proven front-door analog of Codex `--output-last-message`.**
Verified by: an OFFLINE (default-tag) Go unit test feeds a synthetic `spacedock claude` stream-json transcript to the Claude `observed` extractor and asserts it returns the `result`/`success` event's `result` text; a second case with no `result` event asserts the last assistant `text` block fallback; a third case asserts an `is_error:true, api_error_status:401` `result` event is surfaced as a launch failure (non-zero/error), NOT fed into a scenario assertion. (This is offline because the extractor is pure string parsing — no live spend — so the riskiest mechanism gets a fast guard.)

## Test plan

- Focused unit/meta tests for the shared scenario table and runner coverage.
- Focused negative assertion tests for gate hold, rejection routing, and merge-hook guard behavior.
- `go test ./internal/ensigncycle`
- `go test ./...`
- `go test ./... -race`
- Local live run for Codex when Codex auth is available.
- Local live run for Claude when Claude auth is available, or approved `CI-E2E`/`CI-E2E-OPUS` evidence if local auth is unavailable.
- Approved `CI-E2E-CODEX` evidence for the Codex lane before validation passes.

## Stage Report: ideation

- DONE: Design the host-neutral shared scenario table (e.g. sharedRuntimeScenarios()) + the per-host runner-adapter boundary, building directly on the just-merged f4 codex scenario list.
  See "Design refinement" section: the four f4 layers (table/fixtures/prompts/assertions) are already host-neutral and need only a rename; the sole host-specific surface is the runner adapter, captured as a Codex-vs-Claude concern table (auth, plugin, launch, observed-extract, artifacts) with the coverage meta-test as the parity guard.
- DONE: Exercise the riskiest unknown first: a Claude runner consuming a shared scenario ID end-to-end (the real sharing mechanism), before committing the design.
  See "Spike result" section: drove the existing `gate-guardrail` fixture + `codexGatePrompt()` through the real `spacedock claude` front door, extracted the final message from the stream `result`/`success` event, and the UNCHANGED `assertCodexGateHeld(before, after, finalMessage)` PASSED on Claude-produced state+output (gate held, entity unmutated, not archived, final message carried `Gate review:`/`Decision:`). Throwaway spike artifacts removed.

### Summary

The sharing is real, not Codex-shaped: a live Claude FO consumed the shared gate-guardrail scenario ID end-to-end and the host-neutral assertion accepted its output verbatim, proving the table/fixtures/prompts/assertions are already runtime-neutral and only the runner adapter is host-specific. The spike pinned the Claude `observed` source to the stream's `result`/`success` event (the front-door analog of Codex `--output-last-message`) and surfaced a machine gotcha — a stale OAuth benchmark-token returns a 401 `result` and must be a loud launch failure, distinct from an assertion failure — both seeded into AC-7 and the implementation's first test. Added AC-7 for an offline extractor unit test so the riskiest mechanism gets a fast, no-spend guard.

## Stage Report: implementation

- DONE: Extract f4's codexSharedScenarios into a host-neutral sharedRuntimeScenarios() table + per-host runner adapters (Codex + Claude), with the coverage meta-test that fails if a shared scenario lacks a runner for either host (AC-1/AC-2/AC-3).
  `sharedRuntimeScenario`/`sharedRuntimeScenarios()` in shared_scenarios_test.go (commit 80a171e1); `TestSharedRuntimeScenarioDefinitions` reflects over the type and rejects any host-named field (AC-1); `codexScenarioRunners()`/`claudeScenarioRunners()` maps + `TestSharedScenarioRunnerCoverage` fail if either host misses a scenario (AC-2/AC-3, commit 590f90bd). No-spend guards: `go test -tags live -run 'TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions'` → 2/2 pass.
- DONE: Add the Claude runner driving `spacedock claude` and consuming the same scenario IDs; wire .github/workflows/runtime-live-e2e.yml to run both the claude-live and codex-live shared suites (AC-4). The ideation spike already proved the host-neutral assertion accepts Claude output verbatim.
  `claudeLiveRunner` + `TestLiveClaudeSharedScenarios` in claude_live_runner_test.go drives the spike-proven `spacedock claude -- -p <prompt> --output-format stream-json` front door against the SHARED fixtures/prompts/assertions (commit 590f90bd); claude-live now runs `TestLiveEnsignCycle` AND `TestLiveClaudeSharedScenarios`, codex-live runs `TestLiveCodexSharedScenarios` (commit be9eae84). Both live suites SKIP (not fatal) with no auth — verified.
- DONE: Land the behavior/state negative-case tests (AC-5), the AC-7 Claude final-message extraction (stream result/success event; a stale-OAuth-401 must be a LOUD launch failure, distinct from an assertion failure), and the AC-6 README docs for adding a shared scenario.
  AC-5: shared_scenarios_negative_test.go builds the broken end-state per scenario (gate advanced/self-approved, rejection un-routed/partial, merge-hook terminalized) from the real fixtures and proves each assertion reds (commit 3ca05537). AC-7: `extractClaudeFinalMessage` prefers the result/success event, falls back to last assistant text, and surfaces is_error/401 as `errClaudeLaunchFailed` — 5 offline cases pass (commit 80a171e1). AC-6: docs/dev/README.md "Runtime Live CI" section + `TestSharedScenarioDocsContract` (reds when a clause is dropped — verified) (commit c6e5a164).

### Summary

Made the f4 Codex-only scenario surface genuinely shared: the scenario table, fixtures, prompts, and assertions are renamed host-neutral, and the sole host-specific layer is a per-host runner adapter. The Claude runner drives the real `spacedock claude` front door over the SAME scenario IDs and assertions the Codex runner uses, with `extractClaudeFinalMessage` as the proven front-door analog of Codex `--output-last-message` (401 = loud launch failure). A cross-host coverage meta-test guards against one-host-only drift, per-scenario negative cases prove the assertions are behavior/state-oriented (not transcript tautologies), and the README documents the contract. `go test ./...` 843 pass (race-clean); live build/vet clean; both live shared suites skip cleanly without auth. Left streamwatch.go (held am rename) and yy's teardown-hang fixtures untouched. High-stakes CI/live machinery — a detached adversarial audit should follow validation before merge.

## Stage Report: validation

- DONE: Reproduce evidence for every AC: AC-1 (host-neutral table rejects host-named fields), AC-2/3 (runner-coverage meta-test fails if a host misses a scenario), AC-4 (both CI lanes wired, SKIP-not-fatal without auth), AC-5 (negative cases red), AC-6 (README), AC-7 (401 → loud launch failure). Run the offline guards + go test ./...; flag any AC whose only true oracle is a live auth'd run as by-construction-pending-live.
  Offline focused tests 18/18 pass; `go test ./...` all 10 packages `ok` (race not re-run, prior race-clean claim untouched); live-tagged compile+vet clean. AC-1 verified (Mutant 1: a `codexAuth` field reds `TestSharedRuntimeScenarioDefinitions`). AC-2/3 structural guard `TestSharedScenarioRunnerCoverage` passes no-spend; Mutant 2 (drop Claude merge-hook runner) reds it with "no Claude runner". AC-4 workflow runs `TestLiveClaudeSharedScenarios`+`TestLiveEnsignCycle` (claude-live) and `TestLiveCodexSharedScenarios` (codex-live), env-gated + artifacts; both live suites SKIP (not fatal) with empty HOME/no key (verified). AC-7 offline 5/5 pass; Mutant 3 (neuter the is_error guard) reds the two 401 cases. LIVE-PROVEN beyond skip: a real `spacedock claude` gate-guardrail launch fired against the local OAuth token — token is expired, so the runner surfaced `is_error/api_error_status:401` as `errClaudeLaunchFailed` ("claude launch failed before any FO work…401"), NOT fed into an assertion — AC-7's live path confirmed. The real 401 event carries `subtype:"success",is_error:true` (not `error_during_execution`); the extractor keys on `IsError` so it still catches it (verified with a throwaway case). By-construction-pending-live: AC-2 (Codex gate/rejection/merge-hook behavior) and AC-3's gate-hold assertion (needs a valid Claude credential — the local token is expired; CI `CI-E2E`/`CI-E2E-CODEX` approval is the oracle).
- FAILED: Confirm the Claude runner consumes the SAME shared scenario table as Codex (not a parallel copy) — the coverage meta-test must fail if either host lacks a runner for a shared scenario.
  The table IS single-source: both `codexScenarioRunners()` and `claudeScenarioRunners()` are iterated against the one `sharedRuntimeScenarios()` table, and `TestSharedScenarioRunnerCoverage` reds in BOTH directions (missing runner / orphan runner) — Mutant 2 confirmed. BUT the detached adversarial audit found a Material test-strength hole in AC-5's merge-hook negative coverage: `assertMergeHookGuardHeld` has TWO observed-output checks — line 52 (`"merge hook"`/`"merge-hook"` mention) and line 55 (`"cannot advance to terminal"`, the check that proves the guard FIRED). Mutant 6 removed ONLY the line-55 guard-error check and the WHOLE shipped suite stayed GREEN. The shipped negative case (`shared_scenarios_negative_test.go:106`) uses observed `"terminalized merge-check to done"`, which fails the line-52 mention first, so it never isolates line 55. I confirmed the hole with a throwaway adversarial test (observed = "…registers a merge hook [local-merge]. Proceeding.": passes line 52, lacks line 55) — it passes on shipped code and reds when line 55 is removed, proving line 55 is load-bearing and uncovered. This is exactly the AC-5 promise ("not transcript-shape tautologies") for the merge-hook scenario's most behavior-specific clause.

### Summary

Validation of high-stakes live-CI machinery. Six of seven ACs reproduce cleanly: AC-1/AC-4/AC-6/AC-7 fully verified offline, and AC-7's 401-launch-failure path is additionally PROVEN LIVE — a real `spacedock claude` launch hit an expired local token and the runner surfaced the 401 as a loud launch failure distinct from an assertion, the exact spike-predicted gotcha. AC-2/AC-3 are single-source-shared (one table, two runner maps, bidirectional coverage guard) with live behavior by-construction-pending a valid-auth CI run. RECOMMENDATION: REJECTED for one Material finding from the detached adversarial audit — the merge-hook `cannot advance to terminal` guard-error check (`shared_assertions_impl_test.go:55`), the single assertion proving the merge guard actually fired, is not isolated by any shipped negative test: removing it leaves the entire suite green. Fix is a one-line negative case in `shared_scenarios_negative_test.go` whose observed string mentions a merge hook but omits the terminal-guard refusal (route back through implementation feedback per detached-audit policy). All other audit mutants (host-named field, dropped runner, neutered 401 guard, marker-only rejection tautology) were correctly caught. Worktrees left pristine; audit checkout removed.

## Feedback Cycles

### Cycle 1 — detached adversarial audit (in validation, 2026-06-03) — MATERIAL

Validation verified AC-1/4/6/7 (AC-7's 401-launch-failure proven LIVE against a real `spacedock claude`); AC-2/3 are single-source-shared with live behavior by-construction-pending a valid-auth CI run. The detached audit found one Material test-strength hole:

- **MATERIAL — the merge-hook guard assertion is not isolated by a negative test.** The `cannot advance to terminal` guard-error check (`shared_assertions_impl_test.go:55`) — the only assertion proving the merge-hook guard actually FIRED — has no shipped negative case: removing it (Mutant 6) leaves the whole suite GREEN, so a regression that broke the merge-guard scenario's firing would go uncaught.
  **Fix:** add a negative case in `shared_scenarios_negative_test.go` that isolates the merge-hook guard assertion — reproduce Mutant 6 (remove/break the `cannot advance to terminal` check) and confirm the new negative case goes RED (and GREEN on the shipped assertion).

Re-run the targeted negative suite + `go test ./...`, append a cycle-2 Stage Report, commit on your worktree branch, then SendMessage(to="team-lead", message="Done: ...").

## Stage Report: implementation (cycle 2)

- DONE: Close the Material finding — isolate the merge-hook `cannot advance to terminal` guard-error check (shared_assertions_impl_test.go:55) with a negative case, reproduce Mutant 6, confirm RED on mutant / GREEN on shipped.
  Added a negative case to `TestMergeHookGuardrailNegativeBypass` (shared_scenarios_negative_test.go) with observed = "…registers a merge hook [local-merge]. Proceeding without terminalization." — it passes the line-52 merge-hook mention check and so fails ONLY on the line-55 guard-error check, the clause no prior case reached. Verified: GREEN on the shipped assertion; under Mutant 6 (line-55 check removed via perl, then reverted byte-identical — `git diff` clean) the case reds at shared_scenarios_negative_test.go:121 AND the whole ensigncycle suite reds (previously stayed fully green). Commit 77e1258e on worktree branch spacedock-ensign/real-shared-runtime-scenarios.

### Summary

Cycle-1 detached-audit Material finding closed: the merge-hook terminal-guard check is now independently covered by an isolating negative case that reds exactly when that check is removed (Mutant 6 caught — the whole suite now goes red where it previously stayed green). The shipped assertion is untouched (only the test gained a case). `go test ./...` 843 pass, race-clean; live build/vet clean. No other files changed; streamwatch.go and yy's fixtures remain untouched.
