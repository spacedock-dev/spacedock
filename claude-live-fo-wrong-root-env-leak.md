---
title: "Live claude-live FO boots the wrong root — a CI env leak lures the FO off its launch cwd"
status: validation
source: "0240/0qt PR #446 CI (2026-06-30): claude-live (sonnet) failed TestLiveClaudeSharedScenarios twice at claude_live_runner_test.go:122 — 'FO booted the wrong root: expected the fixture root .../001, but the boot command targets /home/user/spacedock-workflow (run 1) / /tmp (run 2) — a CI env leak likely lured the FO off its launch cwd.' claude-live OPUS PASSED the identical suite both runs, so it is sonnet-lane env/model behavior, not the scenarios or a product bug. Recurred 2/2; forced an admin-override merge of 0qt."
id: 2wbxv8hdq5m754h45ehfd75j
started: 2026-07-02T01:23:59Z
worktree: .worktrees/spacedock-ensign-claude-live-fo-wrong-root-env-leak
mod-block: merge:pr-merge
pr: pr-merge:462
---
The live `claude-live` shared-scenarios harness (`internal/ensigncycle`, `claude_live_runner_test.go`) expects the dispatched FO to boot in the per-scenario fixture root (`/tmp/TestLiveClaudeSharedScenarios…/001`). On the sonnet lane the FO's first boot command instead `cd`s to a leaked path OUTSIDE the fixture — `/home/user/spacedock-workflow` (run 1), `/tmp` (run 2) — and `claude_live_runner_test.go:122` flags "a CI env leak likely lured the FO off its launch cwd."

It recurred on both PR #446 runs while the opus lane passed the identical suite, so it is sonnet-lane environment/model behavior, not a scenario or product defect — and distinct from the handoff's named "sonnet headless-gate" flake.

Impact: it blocks the sonnet live lane for ANY PR (it forced the 0qt admin-override merge). Likely fix: have the harness pin/assert the FO's launch cwd (or scrub the leaking env var) so a leaked cwd cannot lure the FO off the fixture root. Surfaced during the 0240 sprint's 0qt merge; filed at captain request.

## Problem statement (reframed with evidence)

The filed premise — a CI env leak lures the sonnet FO off its launch cwd — is disproved by the archived evidence. Both failures are harness-detector false positives on a benign, model-invented probe path; the FO never operated outside the fixture, and both runs completed the scenario correctly.

Evidence, from run 28466995641 (PR #446, attempts 1 and 2; artifact `runtime-live-e2e-claude-live-sonnet`, IDs 7990484876 / 7991648615, per-scenario `claude-stream.jsonl` under `live-artifacts/claude/sonnet/claude-shared-scenarios/feedback-3-cycle-escalation/`):

- Both attempts failed the same subtest, `TestLiveClaudeSharedScenarios/feedback-3-cycle-escalation`, fataled by `detectWrongRootBoot` at `claude_live_runner_test.go:122` — no other scenario failed.
- **No env leak.** Zero `GITHUB_*`/`RUNNER_*` strings appear anywhere in either captured stream — the PR #365 `cleanEnviron` scrub held. The stream's system-init event shows the FO's cwd WAS the fixture root in both runs. The claude-live job runs directly on `ubuntu-latest` (HOME is `/home/runner`), so the flagged `/home/user/spacedock-workflow` cannot exist there; it appears exactly once in the whole stream — inside the FO's own command — and nowhere in the repo, fixture, prompt, or any tool result the FO read. It is model-invented, not leaked.
- **Attempt 1** flagged command: `cd /home/user/spacedock-workflow 2>/dev/null || cd .; ${SPACEDOCK_BIN:-spacedock} --version; echo "---"; pwd; git rev-parse --show-toplevel 2>&1`. The `|| cd .` fallback held: the same command's tool result prints the fixture root for both `pwd` and `git rev-parse --show-toplevel`.
- **Attempt 2** flagged command: `cd /tmp && (echo "--version:"; ${SPACEDOCK_BIN:-spacedock} --version) 2>&1; echo "---"; git rev-parse --show-toplevel 2>&1`. The cd succeeded transiently (result: exit 128, `fatal: not a git repository` — from /tmp), but did not persist: the very next Bash call's `pwd` printed the fixture root and the FO noted "We're at the project root already (not /tmp)".
- **Both runs finished the scenario correctly inside the fixture**: correct cycle-3 escalation recorded, entity edits at absolute fixture paths, final stream `result` event `subtype=success`. Only the detector fataled them (attempt 1 at 65.42s, attempt 2 at 70.67s).

Root cause: sonnet decorates the FO contract's Startup step 1 version probe with a speculative repo-root sniff — `cd <guessed path>` plus `git rev-parse --show-toplevel` — before settling on its launch cwd. `detectWrongRootBoot` (`wrong_root_detect_impl_test.go`) treats ANY absolute `cd` token outside the fixture as a wander, so the harmless transient probe kills a run that would have passed every scenario assertion. Opus does not emit the probe; hence opus green / sonnet red, 2/2.

## Proposed approach

Precision fix inside `detectWrongRootBoot` — no env-scrub change (nothing leaked), no cwd pinning (the launch cwd is already correct), no FO-contract prose change (the probe is harmless model style; a contract ban would spend tokens across all runtimes to patch a harness-only misjudgment). All seven live call sites (claude shared runner, gate-stop, live cycle, merged team-mode, pty team-mode) get the fix for free since it lands inside the detector.

1. **Gate the bare-`cd` signature on same-command corroboration.** A `cd <outside-fixture>` reds only when the SAME command also carries a workflow-operative token (`--workflow-dir`, `--boot`, `--discover`, `status --read`, `state commit`, `new`, or a README/entity path). Pure probes (`--version`, `pwd`, `git rev-parse`) are tolerated.
2. **Standalone operative signatures stay first-and-fatal unchanged:** a `--workflow-dir <outside>` and a workflow-README `Read` outside the fixture.
3. **Reword the failure text** so it stops asserting a cause the evidence disproved. Before: `— a CI env leak likely lured the FO off its launch cwd`. After: `— the FO's boot operated outside its launch cwd`.
4. **Update the offline fixtures.** The synthetic `cd_away_from_fixture_root_reds` case (`cd <realRepo> && spacedock --version`) flips to the corroborated shape (`cd <realRepo> && spacedock status --discover`) — the probe-only shape is benign per the observed cwd-reset behavior. Add the two REAL captured sonnet boot commands verbatim as must-pass cases, and the full captured streams as replay testdata.

### Spike (riskiest unverified mechanism, exercised first)

The riskiest mechanism was whether the corroboration rule actually separates the two real streams from genuine wanders, and whether the archived artifact streams even parse in the detector's shapes. A throwaway spike test (written, run, deleted; results below) replayed both full ~135KB captured streams:

- Current detector on real streams: reproduced BOTH exact CI failures, byte-identical flagged commands and paths (RED reproduction confirmed — the artifact stream-json parses via `streamEntry`/`toolUseBlock`).
- Corroboration-gated variant: GREEN on both real streams; RED on all four corroborated-wander shapes (cd+discover same command, standalone `--workflow-dir` escape, off-fixture README read, compound `cd <realRepo>; ls; cat README.md`). `go test -run TestSpike` passed 2/2 subtests.

Residual risk, on the record: if a future claude version persists Bash cwd across calls, a probe-only cd followed by a relative `status --discover` in a LATER command could boot off-fixture without tripping any signature. The scenario end-state assertions remain the backstop for that case (an FO operating in the wrong root cannot satisfy the fixture-path assertions), and the two captured streams show current claude resets cwd per call.

## Acceptance criteria

- **AC-1 (value, measured against a baseline that can move the wrong way):** Replaying the two archived PR #446 sonnet streams (checked in as testdata) through `detectWrongRootBoot` returns nil — RED today (both currently return the wrong-root error; reproduced in the spike), GREEN with the fix — WHILE all 4 corroborated-wander fixtures still return an error (4/4). The wrong-way direction is real: deleting the cd heuristic outright also greens the real streams but drops the wander fixtures to 3/4 (it passes `cd <realRepo> && spacedock status --discover`). Tested by: `go test ./internal/ensigncycle/ -run TestDetectWrongRootBoot` (offline, milliseconds, no model spend).
- **AC-2 (field confirmation):** One approved sonnet-lane `runtime-live-e2e` run completes with zero `FO booted the wrong root` failures whose flagged command lacks a workflow-operative token. Tested by: the run's gotestsum step log plus the archived `claude-shared-scenarios-detail.jsonl`. (Live and model-stochastic: sonnet emitted the probe 2/2 on #446, but a green run alone cannot distinguish "fixed" from "probe not emitted" — AC-1's deterministic replay is the gating proof; this is the field check.)
- **AC-3 (diagnostic preserved):** An operative wander still fails loud, naming both the expected fixture root and the wandered-to path. Tested by: the existing error-text assertions in `wrong_root_detect_test.go`, kept green over the corroborated fixtures.
- **AC-4:** The detector's failure text no longer asserts a CI env leak as the cause. Tested by: a unit assertion on the error string of a corroborated-wander case.

## Test plan

- **Tier 1 — offline unit (fixture; the RED→GREEN gate; milliseconds, no model):** Extend `wrong_root_detect_test.go`: two real-stream replay cases from testdata (fetch: `gh api repos/spacedock-dev/spacedock/actions/artifacts/7990484876/zip` and `.../7991648615/zip`, path `spacedock/spacedock/live-artifacts/claude/sonnet/claude-shared-scenarios/feedback-3-cycle-escalation/claude-stream.jsonl`, ~135KB each), the 4 corroborated-wander cases, the two real boot commands verbatim as must-pass command-level cases, and the error-text assertions (AC-3/AC-4).
- **Tier 2 — offline package sweep (seconds, no model):** `go test ./internal/ensigncycle/` (default tags) fully green — guards the negative-scenario suites that share these fixtures.
- **Tier 3 — live confirmation (expensive: CI-E2E approval, ~30m wall, real model spend):** one sonnet-lane `runtime-live-e2e` `workflow_dispatch`; grade per AC-2. Confirmation only, not the gating proof.

No doc diff needed: the change is harness-internal (live-test detector precision); no user-visible CLI output, command surface, or docs-site behavior changes.

## Stage Report: ideation

- DONE: The lure mechanism is identified with evidence, not asserted: which leaked env var (or cwd influence) steers the sonnet FO to /home/user/spacedock-workflow or /tmp — spike the riskiest unverified mechanism first and record the result in the task body.
  Disproved the leak with the archived run-28466995641 streams (zero GITHUB_*/RUNNER_* strings; init cwd = fixture root; /home/user/... exists nowhere but the FO's own command): the mechanism is a model-invented probe path plus detectWrongRootBoot's over-broad bare-cd rule. Spike (run first, throwaway, deleted) reproduced both exact CI REDs on the real streams and showed a corroboration-gated variant greens them while all 4 wander fixtures stay red.
- DONE: Acceptance criteria include at least one VALUE-measuring AC against a baseline that can move the wrong way (e.g. sonnet-lane TestLiveClaudeSharedScenarios boots the fixture root, RED today on the leak reproduction, GREEN with the fix), each AC paired with how it is tested.
  AC-1 measures nil-on-real-streams (RED today, spike-reproduced) against the 4/4 corroborated-wander fixtures that a naive heuristic deletion drops to 3/4; AC-2/3/4 each name their test.
- DONE: The test plan separates a cheap harness-level reproduction/spot-check from the expensive live sonnet-lane run, and states which fixture/CLI/live tiers are needed.
  Tier 1 offline replay-fixture unit tests are the RED→GREEN gate (no model spend); Tier 2 package sweep; Tier 3 one approved sonnet-lane live run as field confirmation only.

### Summary

Ideation reframed the task: the filed "CI env leak" premise is disproved by the archived PR #446 artifacts — both failures were detector false positives on sonnet's speculative version-probe cd (fallback held / cd non-persistent; both runs finished the scenario correctly with result=success). Proposed fix is precision inside detectWrongRootBoot (same-command corroboration for bare cd, operative signatures unchanged, reworded error text), proven by a spike replaying the real streams. Cheap deterministic proof is the testdata replay; the live sonnet lane is confirmation only.

## Stage Report: implementation

- DONE: AC-1 lands RED then GREEN — the two archived PR #446 sonnet streams checked in as testdata, replaying through detectWrongRootBoot returns nil, while all 4 corroborated-wander fixtures still return an error (4/4).
  Fetched artifacts 7990484876/7991648615 via `gh api .../actions/artifacts/{id}/zip`; checked in as `internal/ensigncycle/testdata/claude_live_wrong_root_pr446_attempt{1,2}.stream.jsonl`. Manually reverted the impl to pre-fix and reran `TestDetectWrongRootBootRealPR446Streams` — reproduced both exact byte-identical CI failures (old "CI env leak" text) — then restored the fix and reran green. Commit 62d8c6d6.
- DONE: The detector change matches the approved design — bare cd reds only with same-command corroboration; standalone operative signatures (--workflow-dir, README read) stay first-and-fatal; failure text reworded off the disproved env-leak cause with unit assertions on the error string.
  `hasWorkflowOperativeSignature` in wrong_root_detect_impl_test.go checks --workflow-dir/--boot/--discover/status --read/state commit/new/README-or-.md-token; `wanderTarget` gates only `cd`-origin path args on it via a new `bootPathArg{path,cd}` tag, leaving `--workflow-dir` args ungated. Both error strings reworded to "the FO's boot operated outside its launch cwd"; `!strings.Contains(err.Error(), "CI env leak")` assertions added to 3 red subtests (AC-3/AC-4).
- DONE: Tier 1 and Tier 2 offline suites are green — extended wrong_root_detect_test.go (real-stream replays, the two real boot commands verbatim as must-pass cases, error-text assertions) and go test ./internal/ensigncycle/ — no live run this stage.
  `go test ./internal/ensigncycle/ -run TestDetectWrongRootBoot -v`: 11/11 + 2/2 + 2/2 subtests pass. `go test ./internal/ensigncycle/...`: full package green (7.2s). `go vet ./internal/ensigncycle/...` and `go build ./...` clean; gofmt clean. Confirmed no other caller references the changed `bootPathArgs`/`wanderTarget` signatures (grep, single file).

### Summary

Fixed detectWrongRootBoot to gate its bare-`cd`-escapes-fixture signature on same-command workflow-operative corroboration, leaving the standalone `--workflow-dir`/README-read signatures first-and-fatal. Checked in the two archived PR #446 sonnet streams as testdata and added replay tests (verified RED on the pre-fix detector by manual revert, byte-identical to the real CI failures, then GREEN restored) alongside command-level verbatim cases and the 4 corroborated-wander fixtures (still 4/4 red). Full offline package suite (Tier 1 + Tier 2) is green; no live run performed this stage per the stage boundary (AC-2 belongs to validation).

## Stage Report: validation

- DONE: AC-1 reproduced, not trusted: run the extended wrong_root_detect suite yourself; confirm the two archived PR #446 stream replays return nil AND all 4 corroborated-wander fixtures return an error (4/4); then exercise the wrong-way baseline — locally delete the bare-cd corroboration heuristic and confirm the wander fixtures drop below 4/4 (proving the tests can fail); revert. Reject self-referential evidence.
  Fresh `-count=1` run: both real-stream replays nil, all 4 wander fixtures red (15/15 subtests). Non-self-referential: re-fetched artifacts 7990484876/7991648615 from GitHub and diffed — testdata byte-identical to the live CI captures, and each stream's init cwd equals the test's fixtureRoot verbatim. Wrong-way both directions: removing the corroboration gate re-reds both real streams with the exact CI-flagged commands/paths; deleting the cd signature outright greens the streams but drops the wander fixtures to 2/4 (cd_away and compound_cd_away both fail — stricter than the entity's predicted 3/4 because the implementation-added compound-cd fixture also catches the deletion). Both mutations reverted; tree clean; suite green post-revert.
- DONE: AC-3 and AC-4 verified on the error path: an operative wander still fails loud naming the expected fixture root and the wandered-to path, and the failure text no longer asserts a CI env leak — confirm via the unit assertions AND by reading the shipped error string.
  The 4 red subtests assert the error names the fixture root (and the wandered-to path where applicable) and 3 assert `!Contains("CI env leak")`; observed the live error output during the mutation run naming both paths with the reworded suffix "the FO's boot operated outside its launch cwd". Both shipped format strings (wrong_root_detect_impl_test.go:64,73) carry no leak claim; a repo grep finds "CI env leak" only in comments and the test assertions, and the runner call site (claude_live_runner_test.go:448) relays the detector error verbatim.
- DONE: Tier 2 sweep green (go test ./internal/ensigncycle/ plus applicable Testing Resources); AC-2 (live sonnet-lane field confirmation) is explicitly deferred — record it as pending-captain-approval in your report, not SKIPPED silently; close with a PASSED/REJECTED recommendation citing evidence per AC.
  `go test ./internal/ensigncycle/ -count=1` ok (8.0s); repo-wide `go test ./... -count=1` 15/15 packages ok; `go vet` and `gofmt -l` clean. AC-2 is PENDING-CAPTAIN-APPROVAL: one sonnet-lane `runtime-live-e2e` workflow_dispatch (~30m, real model spend) remains as field confirmation only — AC-1's deterministic replay is the gating proof.

### Summary

Recommendation: PASSED (AC-2 field confirmation pending captain approval of a live sonnet-lane run). AC-1 verified by independent artifact re-fetch (byte-identical testdata) plus two-direction mutation testing at commit 62d8c6d6: the pre-fix detector reproduces both CI failures on the real streams, the naive cd-deletion baseline drops the wander fixtures below 4/4, and the shipped detector separates them cleanly. AC-3/AC-4 verified on the exercised error path and the shipped format strings. No code changes were needed; the implementation worktree is clean at 62d8c6d6.
