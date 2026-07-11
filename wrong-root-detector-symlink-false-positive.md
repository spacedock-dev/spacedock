---
title: "Wrong-root boot detector false-positives on macOS: observed path not symlink-resolved before comparison"
status: validation
source: "Found 2026-07-09 while investigating one of sc5's (live-runner-boot-preamble-hardening) three anomalous PR #490 merge-gate failures locally, per captain direction to verify live behavior locally rather than only via GitHub Actions CI (which runs on ubuntu-latest and never hits this macOS-only symlink form). Confirmed via a local live re-run of TestLiveClaudeSharedScenarios/self-evidence-merge-triage on this machine (macOS): the test failed with 'FO booted the wrong root: expected .../private/var/folders/.../001, but it read the workflow README at .../var/folders/.../001/README.md' — the two paths are IDENTICAL directories (macOS /var is a symlink to /private/var); the FO did not actually wander."
started: 2026-07-09T11:49:49Z
completed:
verdict:
score: 0.5
worktree: .worktrees/spacedock-ensign-wrong-root-detector-symlink-false-positive
issue:
id: 5qae7c01tnytacehaphrda4s
mod-block: merge:pr-merge
pr:
---

`internal/ensigncycle/wrong_root_detect_impl_test.go`'s `detectWrongRootBoot` resolves the fixture root via `filepath.EvalSymlinks` before comparing (lines 39-46, with an explicit comment acknowledging the macOS `/var` vs `/private/var` symlink case) — but the two call sites that extract the FO's OBSERVED path from its tool-call stream do not apply the same resolution to that observed path before calling `isUnder`:

- `wanderWorkflowReadme(filePath, fixtureRoot string)` (line 153-165): only `filepath.Clean(filePath)`, no `EvalSymlinks`.
- `wanderTarget(command, fixtureRoot string)` (line 178-198): only `filepath.Clean(arg.path)` on each Bash-command path argument, no `EvalSymlinks`.

When Go's `t.TempDir()` returns an unresolved `/var/folders/...` path (the fixture root) and the FO's own tool calls report that same unresolved form (because that literally is its cwd/argument), while the CALLER'S `fixtureRoot` argument has already been resolved to `/private/var/folders/...`, `isUnder(p, fixtureRoot)` sees two paths with no common prefix and false-flags a wander that never happened — the FO read/operated on exactly the right directory, just spelled two ways that resolve to the same inode.

This is macOS-specific (Linux/CI has no `/var` → `/private/var` symlink, so `ubuntu-latest` CI runners never hit it) but it directly undermines the value of running these live scenario tests locally on a macOS development machine — any scenario whose observed Read/Bash path happens to preserve the unresolved `/var/...` spelling can spuriously fail with a "wrong-root" diagnosis that is actually correct behavior.

## Proposed approach

Centralize into ONE shared canonicalizer that resolves symlinks for the wrong-root comparison, and route ALL THREE sites through it: the fixture-root resolution currently inline at `detectWrongRootBoot` (lines 45-47), and the two observed-path sites — `wanderWorkflowReadme`'s `filepath.Clean(filePath)` (line 157) and `wanderTarget`'s `filepath.Clean(arg.path)` (lines 185 and 197). Centralizing is not cosmetic: `isUnder` only preserves the under-relationship when BOTH operands are canonicalized by the SAME function, so a single helper is what guarantees the fixture side and the observed side can never drift into inconsistent symlink spellings (the exact defect today).

The canonicalizer is a **deepest-existing-ancestor** resolve, NOT the plain "`EvalSymlinks`, else use the `Clean`ed unresolved form" fallback the draft proposed. The spike (below) disproves the plain fallback: for a stream-extracted path that does not exist at check-time, `EvalSymlinks` errors and the unresolved `/var/...` form survives — leaving the very mismatch we are trying to erase. The helper instead:

1. `p = filepath.Clean(p)`; if `EvalSymlinks(p)` succeeds, return it.
2. Otherwise walk up to the deepest EXISTING ancestor, `EvalSymlinks` that ancestor, and re-join the non-existent tail onto the resolved ancestor.
3. If nothing up to the filesystem root resolves, return `Clean(p)` (true last-resort fallback).

Because the macOS divergence is a system-symlink PREFIX (`/var` → `/private/var`, and likewise `/tmp` → `/private/tmp`), and `/var`/`/var/folders` always exist, step 2 canonicalizes even a fully non-existent `/var/folders/.../README.md` to its `/private/var/...` form. The transform is prefix-preserving, so applying it to both operands of every existing `isUnder`/`==` comparison leaves all current results unchanged (parent/child/sibling relationships are invariant), which is why existing detector tests stay green while the symlink blind spot closes.

The helper lives beside the detector in `internal/ensigncycle/wrong_root_detect_impl_test.go` (the detector is test-infrastructure — all of it is in `_test.go` files; there is no shipped product-code change).

## Riskiest-mechanism spike (DONE — result on record)

Spiked `filepath.EvalSymlinks` behavior on a non-existent observed path BEFORE committing to the fallback design (checklist item 3), plus the live-vs-replay reproduction conditions. Standalone Go spike, scratchpad, three findings:

- **`EvalSymlinks` errors (IsNotExist) on a non-existent path.** So the fallback branch is genuinely reached for stream-extracted paths that no longer exist at check-time — this is not theoretical.
- **The draft's plain-`Clean` fallback is WRONG.** For a ghost `/var/folders/.../does/not/exist/README.md`, `EvalSymlinks` errors and `Clean` leaves it in `/var/...` form → still mismatches a resolved `/private/var/...` fixture root → STILL false-flags. The deepest-existing-ancestor resolve returns `/private/var/folders/.../does/not/exist/README.md` (prefix resolved via the existing `/var` symlink, ghost tail preserved) → matches. This finding is why the fix is ancestor-walk, not plain fallback.
- **Reproduction requires the fixture root to physically EXIST at check-time.** With the fixture root present (the live case), current code false-flags — bug reproduced. With the fixture root GONE (a clean machine replaying an archived stream), the fixture-side `EvalSymlinks` no-ops, both sides stay unresolved `/var/...`, and there is NO flag. Consequence for the test plan: an archived-stream replay canNOT serve as a deterministic RED control off the origin machine; the RED control must build real on-disk symlink divergence itself.

## Acceptance criteria

- **AC-1 (RED control — symlink mismatch no longer false-flags; deterministic, cross-platform).** An offline test builds a REAL on-disk symlink divergence in `t.TempDir()` (a real dir plus a symlink to it — deterministic on macOS AND on Linux CI, not reliant on the `/var` symlink), passes the LINK-spelling root to `detectWrongRootBoot`, and feeds a stream whose observed Read/Bash path uses the same LINK spelling. Two variants: (a) the observed README EXISTS under the link; (b) a GHOST variant where the observed path does not exist on disk. Both must return `nil` (no wander). Test: on the pre-fix detector both variants must ERROR (the false positive is reproduced); on the post-fix detector both return `nil`. The ghost variant is the discriminating case — it fails under the rejected plain-`Clean` fallback and passes only under the ancestor-walk helper, locking in the correct design.
- **AC-2 (positive control — a genuine wander is still caught).** With the same symlinked setup, an observed path pointing at a genuinely unrelated directory (outside the resolved root) must still be flagged as a wander. Additionally, every existing case in `TestDetectWrongRootBoot`, `TestDetectWrongRootBootRealPR446Commands`, and `TestDetectWrongRootBootRealPR446Streams` must stay green — the fix corrects the symlink blind spot without weakening detection. Test: `go test ./internal/ensigncycle/ -run TestDetectWrongRootBoot` stays fully green with the new red/positive cases added.
- **AC-3 (end-value — the real local false positive is gone; measured against a baseline that can move the wrong way).** Re-run the real live scenario `TestLiveClaudeSharedScenarios/self-evidence-merge-triage` locally on this macOS machine. Baseline = the count of scenarios false-flagged as wrong-root on a local macOS live pass: currently ≥1 (self-evidence-merge-triage, per the source note), target 0, with genuine-wander detections unchanged (AC-2 guards the other direction). Test: pre-fix the scenario fails with the wrong-root false positive; post-fix the wrong-root check no longer fires (the scenario passes or fails only on unrelated grounds). This is the "reproducible against the real local false positive, not just synthetic fixtures" requirement; the live re-run is the independent baseline, AC-1/AC-2 are the fast deterministic guards.

## Test plan

- **AC-1 / AC-2 — offline Go unit tests** in `internal/ensigncycle/wrong_root_detect_test.go`, alongside the existing `TestDetectWrongRootBoot` table. Uses `os.Symlink` under `t.TempDir()` and the existing `streamLine`/`bashToolLine`/`toolResultLine` helpers to synthesize the boot stream. Cost: low (sub-second, no network, no live model); complexity: low. These run in CI on Linux too because the divergence is self-made, so they guard against regressing the fix on every push — a real gain, since the underlying bug is macOS-only and Linux CI would otherwise never exercise it.
- **AC-3 — live workflow smoke test**, macOS only, run by hand: `go test ./internal/ensigncycle/ -run 'TestLiveClaudeSharedScenarios/self-evidence-merge-triage'`. Cost: high (real Claude run, minutes, spends live budget); complexity: medium (requires the live harness/credentials). Run once pre-fix to capture the false positive and once post-fix to confirm it is gone; not part of the default CI gate (CI is ubuntu-latest and never hits the symlink form).
- No fixture-file (`testdata/*.stream.jsonl`) archived-replay case is added: the spike shows a replayed stream cannot be a deterministic RED control off the origin machine (the fixture root must exist at check-time), so it would be a green-only no-op rather than a real guard.

## Documentation

No doc diff needed: this changes an internal test-harness detector (`internal/ensigncycle`, all `_test.go`), not a user-visible surface — no CLI output, command surface, startup banner, or docs-site behavior changes. The only user-observable string is a test-failure message, which is not a product doc surface.

### Feedback Cycles

- Cycle 1 (validation): The first validation report rejected AC-3 only because the dispatch imposed a filesystem-safehouse prerequisite that is not part of this entity's acceptance criteria or documented local-live procedure. The captain confirmed `benchmark-token` is available and explicitly authorized necessary local live runs. No implementation change is requested; rerun the same independent validator against the canonical fresh-`HOME` live harness and report AC-3's actual outcome.

## Stage Report: ideation

- DONE: Concrete proposed fix — resolve symlinks on the observed path at both call sites before `isUnder`, with a documented fallback, matching detectWrongRootBoot's fixture-root resolution; decide whether to centralize
  See "## Proposed approach": ONE deepest-existing-ancestor canonicalizer routed through all three sites (fixture-root lines 45-47, `wanderWorkflowReadme` line 157, `wanderTarget` lines 185/197); centralization justified as the invariant that keeps both `isUnder` operands consistently spelled.
- DONE: Acceptance criteria with red-control + positive-control, reproducible against the real local false positive, not just synthetic fixtures
  See "## Acceptance criteria": AC-1 (deterministic self-made-symlink red control, existing+ghost variants), AC-2 (genuine-wander positive control + existing tests green), AC-3 (real live self-evidence-merge-triage re-run, false-flag count 1→0 as the value measurement).
- DONE: Riskiest-mechanism spike — confirm EvalSymlinks behavior when the observed path doesn't exist at check-time, before committing to the fallback design
  See "## Riskiest-mechanism spike": ran Go spikes — EvalSymlinks errors (IsNotExist) on ghost paths; the draft's plain-`Clean` fallback leaves `/var/...` unresolved and STILL false-flags; deepest-existing-ancestor resolve fixes it; reproduction requires the fixture root to physically exist, so archived-stream replay is not a deterministic RED control.

### Summary

Firmed the fix to a single deepest-existing-ancestor canonicalizer shared by the fixture-root side and both observed-path sites, rejecting the draft's plain-`Clean`-fallback after the spike proved it still false-flags ghost (non-existent) stream paths. The spike also established that reproduction requires the fixture root to exist at check-time, which redirected the test plan away from an archived-stream replay toward a self-made on-disk symlink red control (deterministic and cross-platform, so it also guards the macOS-only bug on Linux CI) plus a real live self-evidence-merge-triage re-run as the end-value measurement. No doc diff needed — the detector is internal test infrastructure with no user-visible surface.

## Stage Report: implementation

- DONE: Complete and commit one deepest-existing-ancestor canonicalizer used consistently for fixture and observed paths, with TDD red/green evidence.
  Commit `49c57685a62ea501000ded8846eb13274222e3c6` routes the fixture root, observed README path, and both observed Bash-path passes through `canonicalizeWrongRootPath`.
- DONE: Preserve genuine wrong-root detection and prove existing-path, ghost-path, and unrelated-directory controls with focused tests.
  The pre-fix focused run failed the three symlink-equivalent observed-path cases while its genuine-wander control passed; the restored fix passed all 5 symlink controls and all 25 `TestDetectWrongRootBoot` cases.
- DONE: Append a complete implementation Stage Report naming commit, changed files, and exact clean verification commands; leave the live macOS end-value run to validation.
  This report names the commit, both changed files, the exact commands and results, and explicitly reserves AC-3's live end-value run for the fresh validation stage.
- SKIPPED: Re-run the real live macOS `TestLiveClaudeSharedScenarios/self-evidence-merge-triage` scenario.
  Intentionally left to independent validation as required by the implementation dispatch; no live model budget was spent in this stage.

### Summary

Commit `49c57685a62ea501000ded8846eb13274222e3c6` changes `internal/ensigncycle/wrong_root_detect_impl_test.go` and `internal/ensigncycle/wrong_root_detect_test.go`: one deepest-existing-ancestor canonicalizer now preserves symlink equivalence for existing and ghost paths without weakening unrelated-directory detection. TDD evidence used `go test ./internal/ensigncycle -run '^TestDetectWrongRootBootSymlinkMismatch$' -count=1` (RED: 1 control passed, 4 test nodes failed for the expected false positives; GREEN: 5 passed), followed by `go test ./internal/ensigncycle/ -run TestDetectWrongRootBoot -count=1` (25 passed). Clean verification used `gofmt -w ./cmd ./internal`, `gofmt -d internal/ensigncycle/wrong_root_detect_impl_test.go internal/ensigncycle/wrong_root_detect_test.go` (empty), `git diff --check 821e5b08742696b1bde6e6304024086bce03b5fb..49c57685a62ea501000ded8846eb13274222e3c6` (empty), `go test ./...` (exit 0), `go test ./... -race` (exit 0), and `git status --short --branch` (clean).

## Stage Report: validation

- DONE: Independently review the rebased diff against origin/main and verify the canonicalizer is applied consistently without weakening genuine-wander detection.
  Audited clean rebased commit `5c754fb0b52d4bd5775fec2526bc2675dc67875a` against current `origin/main` `2137832bcbb7fe78e38d054c3e4b9c0e1ac32186`: fixture root, README observation, and both Bash-path passes share `canonicalizeWrongRootPath`; no weakening finding.
- DONE: Reproduce the focused symlink existing/ghost/unrelated controls, then run go test ./... and go test ./... -race on the rebased commit.
  Focused symlink test passed all 5 nodes (existing README, ghost README, Bash workflow-dir, genuine unrelated wander, parent); detector family passed 25 nodes; both full gates passed 2,101 tests in 17 packages.
- SKIPPED: Measure AC-3 live only inside a proven filesystem-isolated safehouse; if isolation is unavailable, do not run unsandboxed and report the exact safety blocker instead of rendering a false PASS.
  Filesystem isolation is unavailable: `command -v safehouse` exited 1, invoking the discovered local binary failed with `Operation not permitted`, no `.safehouse` profile exists, and this worker has unrestricted filesystem access; the live model was not launched unsandboxed.
- DONE: Verify AC-1 (symlink mismatch no longer false-flags; deterministic, cross-platform).
  `go test ./internal/ensigncycle -run '^TestDetectWrongRootBootSymlinkMismatch$' -count=1 -v` passed the existing-path and deepest-existing-ancestor ghost-path variants.
- DONE: Verify AC-2 (a genuine wander is still caught).
  The same focused run passed `genuine_wander_still_reds`, and `go test ./internal/ensigncycle -run 'TestDetectWrongRootBoot' -count=1` passed all 25 detector nodes including archived PR #446 command/stream controls.
- FAILED: Verify AC-3 (the real local false positive is gone).
  No safe AC-3 evidence was produced because the required filesystem-isolated safehouse is unavailable; this criterion remains unverified rather than falsely passed.

### Summary

The rebased implementation is internally consistent and all deterministic acceptance and repository-wide tests pass: `go test ./...` and `go test ./... -race` each completed with 2,101 passing tests across 17 packages. Recommendation: REJECTED pending a safe, filesystem-isolated AC-3 live run; this is an evidence blocker, not an implementation defect found in AC-1 or AC-2.

## Stage Report: validation (cycle 2)

- DONE: Re-review the updated entity and apply the captain's canonical local-live clarification.
  The cycle-1 safehouse prerequisite was not an AC-3 requirement; this cycle used the repository's documented fresh-isolated-`HOME` Claude harness with the available `~/.claude/benchmark-token`.
- DONE: Run exactly the post-fix macOS `TestLiveClaudeSharedScenarios/self-evidence-merge-triage` scenario required by AC-3.
  With `SPACEDOCK_BIN=$PWD/spacedock` and `SPACEDOCK_REPO_ROOT=$PWD`, `go test ./internal/ensigncycle/ -tags live -count=1 -timeout 40m -run 'TestLiveClaudeSharedScenarios/self-evidence-merge-triage' -v` passed in 207.016s (subtest 206.57s, exit 0).
- DONE: Record whether wrong-root detection fires and distinguish any unrelated failure.
  Wrong-root false-flag count was 0: the detector did not fire; the scenario itself passed, so there was no unrelated failure to qualify or mask the AC-3 result.
- DONE: Reconfirm the worktree remains clean and retain the prior deterministic/full/race evidence.
  The generated launcher was removed and `git status --short --branch` is clean; cycle 1's 5-node symlink controls, 25-node detector family, and both 2,101-test full/race gates remain the recorded evidence on unchanged commit `5c754fb0b52d4bd5775fec2526bc2675dc67875a`.
- DONE: Verify all acceptance criteria and make a terminal recommendation.
  AC-1 and AC-2 retain their independently reproduced deterministic evidence; AC-3 now has the required live macOS evidence, moving the observed wrong-root false-flag baseline from at least 1 to 0.

### Summary

Validation cycle 2 closes the only evidence gap: the real post-fix macOS `self-evidence-merge-triage` live scenario passed under the canonical fresh-isolated-`HOME` harness without any wrong-root detection or unrelated failure. Recommendation: PASSED; no implementation code changed, and the worktree is clean.
