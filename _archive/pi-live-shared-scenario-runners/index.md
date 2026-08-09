---
title: Pi live shared scenario runners — enable shallow-boot + filing + gate-guardrail on the pi-live CI lane
status: backlog
source: Superseded on 2026-08-09 by restore-optional-manual-pi-common-live-ci. The original task assumed that Pi lacked shared runners; current main has the common runner and now lacks only the optional manual CI job.
score:
started:
completed:
verdict:
worktree:
issue:
sprint:
sprint-readiness:
id: asc4j2vs2qqgtg04hp87j2tq
archived: 2026-08-09T14:30:08Z
---

# Pi live shared scenario runners

## End value

The `pi-live` CI lane runs `TestLivePiSharedScenarios` covering **shallow-boot**, **filing**, and **gate-guardrail** — the three shared runtime scenarios pi can now meaningfully exercise after 0223. A regression in the pi FO boot path, the `spacedock new` atomic-id-mint path, or the gate-presentation contract fails CI instead of surfacing in a live session. (The blank-id `pi-devoverride-package-ok` boot failure this sprint would have been caught by `shallow-boot`.)

## Problem — what exists vs. what's needed

**Codex (the pattern to mirror) has:**
- `codexLiveRunner` struct + `newCodexLiveRunner` (auth/HOME isolation, codex install).
- `TestLiveCodexSharedScenarios` (iterates `sharedRuntimeScenarios()`, runs each).
- `codexScenarioRunners()` map — 6 entries.
- 6 `runCodex<Scenario>Scenario` funcs — the actual scenario implementations.

**Pi has:**
- `pi_live_runner_test.go` — ONLY `TestLivePiFrontDoorSmoke` (trivial dispatch smoke) + a prompt-shape test. **No `piLiveRunner`, no `piScenarioRunners` map, no per-scenario runners.**
- `pi_shared_coverage_test.go` — the coverage meta-test. The `mode: "gap"` entries are **documentation, not a switch**: the meta-test only validates the mode is one of `live`/`codified`/`gap` + has a reason string. Flipping `"gap"`→`"live"` changes nothing about what runs — only the coverage bookkeeping. No test would execute.

**Reusable from the smoke + codex runner:** the shared scenario definitions (`sharedRuntimeScenarios()`), shared assertions (`shared_assertions_test.go`), the smoke's auth/env/fixture helpers (`newPiLiveSmokeFixture`, `seedPiLiveAuth`, `piLiveEnv`, `runPiLiveCommand`), the CI lane's pi install + substrate setup (`PI_SUBAGENTS_PACKAGE_ROOT`/`PI_INTERCOM_PACKAGE_ROOT`).

**Net-new:** the `piLiveRunner` consolidation, `TestLivePiSharedScenarios` driver, `piScenarioRunners` map, and 3 per-scenario runner funcs (the substantive work).

## Scope — 3 scenarios, in 2 tiers

### Tier 1 (no dispatch — pi-ready now, no 0223 dependency)

- **`shallow-boot`** — FO boots, greets, advances a merged PR before-greet (S7b), **no team created, no worker dispatched**, then stops for input. Exercises the boot path that just broke (the blank-id `pi-devoverride-package-ok` failure surfaced at boot). Highest value, lowest risk. No dispatch → no dependency on skill injection / model stamping / back-channel.
- **`filing`** — FO files a new seed entity via the atomic `spacedock new <slug>` path (not `--next-id` + hand-write). Exercises the id-minting the Commander bypassed (hand-wrote the blank-id file). No dispatch.

### Tier 2 (single dispatch — after capstone #409 merges)

- **`gate-guardrail`** — FO halts at a human gate, presents the review, no self-approval/mutation/archival. Needs a dispatched ensign producing a stage report at a gate stage, then the FO presents. `eq` (#406) + `bdt` (#405) make a pi ensign dispatch correct; the capstone's back-channel isn't strictly required (the ensign completes and returns; the FO presents — no mid-run `need_decision`). Runnable now, but safer to land after #409 so the dispatch path is the final 0223 shape.

### Deferred (out of scope — recorded for a later task)

- **`rejection-flow`** + **`feedback-3-cycle-escalation`** — depend on reviewer **reuse** (reuse-condition-1: reusable reviewer handle). Pi's default is fresh-redispatch (friction 9, deferred in 0223). The reuse assertion would fail on pi. Defer until reuse-advance lands.
- **`merge-hook-guardrail`** — FO-side guardrail (the `pr-merge` mod-block enforcement). Doesn't strictly need dispatch but needs a PR open + the merge ceremony — CI-environment-complex. Defer until Tier 1/2 green.

## Approach (mirror the codex runner)

1. **`piLiveRunner`** — consolidate the smoke's `newPiLiveSmokeFixture`/`seedPiLiveAuth`/`piLiveEnv`/`runPiLiveCommand` into a `piLiveRunner` struct (bin, env, artifact root) + `newPiLiveRunner` constructor, mirroring `codexLiveRunner`/`newCodexLiveRunner`. Pi-specific: copy `~/.pi/agent/auth.json` or `OPENAI_API_KEY`, pi install, `PI_SUBAGENTS_PACKAGE_ROOT`/`PI_INTERCOM_PACKAGE_ROOT` setup (the CI lane already does the install; the runner consumes it).
2. **`TestLivePiSharedScenarios`** — mirror `TestLiveCodexSharedScenarios` (iterate `sharedRuntimeScenarios()`, `t.Run` each).
3. **`piScenarioRunners()`** map — 3 entries: `shallow-boot`, `filing`, `gate-guardrail`.
4. **3 per-scenario runner funcs** — `runPiShallowBootScenario`, `runPiFilingScenario`, `runPiGateGuardrailScenario`. Each sets up the scenario-specific workflow fixture (reuse `shared_fixtures_test.go` where possible), runs `spacedock pi` with the FO prompt, asserts the shared post-conditions (reuse `shared_assertions_test.go`). Pi's dispatch model (subagent) differs from codex (thread), so these are adaptations of the codex runners, not copies.
5. **CI lane** — add a `gotestsum --jsonfile pi-shared-scenarios-detail.jsonl --format pkgname -- -tags live -count=1 -timeout 40m -run TestLivePiSharedScenarios ./internal/ensigncycle` step to the `pi-live` job (after the existing smoke step). Upload the detail jsonl as an artifact.
6. **Coverage map** — flip `shallow-boot`, `filing`, `gate-guardrail` from `"gap"` to `"live"` in `pi_shared_coverage_test.go` (bookkeeping; the meta-test then documents them as covered).

## Acceptance criteria (provisional — ideation finalizes; proof = behavior)

**AC-1 — `TestLivePiSharedScenarios` runs shallow-boot + filing + gate-guardrail on the pi-live lane and they pass.**
Verified by: the pi-live CI lane (or a local `-tags live` run with auth) executing `TestLivePiSharedScenarios` and all 3 subtests passing; the detail jsonl artifact captured. Each scenario asserts its shared post-conditions (boot-greet-no-dispatch for shallow-boot; atomic-id-mint entity filed for filing; gate-halt-no-self-approval for gate-guardrail).

**AC-2 — `shallow-boot` catches a boot-path regression (the blank-id class).**
Verified by: the shallow-boot scenario's assertion that a booted FO greets + reports accurate state from a clean state checkout — a blank-id archived entity (injected as a negative fixture) fails the boot and the test reds. (The scenario's value proposition: it would have caught `pi-devoverride-package-ok`.)

**AC-3 — The `piLiveRunner` reuses the smoke's auth/env/fixture helpers (no duplication).**
Verified by: a structural review that `piLiveRunner` consolidates — not duplicates — `newPiLiveSmokeFixture`/`seedPiLiveAuth`/`piLiveEnv`/`runPiLiveCommand`; the smoke test still passes against the consolidated runner.

**AC-4 — The 3 coverage-map entries flip `gap`→`live` and the meta-test passes.**
Verified by: `go test -run TestPiSharedScenarioCoverage ./internal/ensigncycle` green with the 3 entries as `live` (honest reason updated from "gap" rationale to "covered by TestLivePiSharedScenarios").

**AC-5 — Deferred scenarios stay `gap` with their honest reasons.**
Verified by: `rejection-flow`, `feedback-3-cycle-escalation`, `merge-hook-guardrail` remain `gap` in the coverage map with their reuse/complexity reasons; the meta-test passes.

## Out of scope

- `rejection-flow`, `feedback-3-cycle-escalation`, `merge-hook-guardrail` runners (deferred — see Scope).
- A pi `TestLiveEnsignCycle` equivalent (the claude lane's core live drive) — related but separate; file as a sibling if not absorbed into this task's ideation.
- The `pi-launch-fnm-multishell-race` fix (`j7n…`) and the `status-validate-determinism` fix (`a1a…`) — separate followups; this task consumes their fixes if they've landed but doesn't depend on them (shallow-boot/filing don't dispatch; gate-guardrail works without the fnm fix on CI's non-fnm runner).

## Test plan

- Local `-tags live` run with pi auth: `TestLivePiSharedScenarios` 3/3 green.
- CI `pi-live` lane: the new `gotestsum -run TestLivePiSharedScenarios` step green; detail jsonl uploaded.
- `go test ./...` (non-live) green — the runner scaffolding + coverage map flip don't break the offline suite.
- AC-2 negative fixture: inject a blank-id archived entity, assert shallow-boot reds.
- This is a CI/test-surface change (the `pi-live` lane + `internal/ensigncycle/*_test.go`). High-stakes per the proof policy (CI/release machinery). Detached adversarial audit at validation.

## Related

- `internal/ensigncycle/codex_live_runner_test.go` — the runner pattern to mirror.
- `internal/ensigncycle/pi_live_runner_test.go` — the existing smoke + helpers to consolidate.
- `internal/ensigncycle/pi_shared_coverage_test.go` — the coverage map to flip 3 entries.
- `internal/ensigncycle/shared_scenarios_test.go` (`sharedRuntimeScenarios()`), `shared_assertions_test.go`, `shared_fixtures_test.go` — shared, reusable.
- `.github/workflows/runtime-live-e2e.yml` (`pi-live` job, line ~442) — the lane to extend.
- 0223 members: `eq` (#406 merged, skill injection), `bdt` (#405 merged, model stamping), `b2` (#409 pending, back-channel capstone) — gate-guardrail depends on eq+bdt; safer after b2.
- `pi-devoverride-package-ok` (archived) — the blank-id boot failure shallow-boot would have caught.
