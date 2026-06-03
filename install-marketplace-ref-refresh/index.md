---
id: s0cqcf9hg4k0tgartr6ymf1b
title: spacedock install --host claude silently reuses stale marketplace ref (defeats fresh-install-from-next)
status: validation
source: captain (2026-06-02) — observed `./spacedock install --host claude` reports success but the marketplace add no-ops on "already on disk", so the @next ref never replaces the existing pin; the uninstall+install cycle then re-pulls from the stale ref
score: "0.38"
worktree: 
started: 2026-06-02T21:14:43Z
completed: 2026-06-03T00:18:17Z
verdict: PASSED
issue:
mod-block: 
pr: "#272"
---

`./spacedock install --host claude` runs a 3-command sequence (`internal/cli/host_exec.go:235-241` `installArgvSequence`):

1. `claude plugin marketplace add spacedock-dev/spacedock@next`
2. `claude plugin uninstall spacedock@spacedock`
3. `claude plugin install spacedock@spacedock`

Step 1 is a no-op when the marketplace `spacedock` is already declared in user settings — claude emits `Marketplace 'spacedock' already on disk — declared in user settings` and exits 0. The `@next` ref pin never lands. Steps 2-3 then uninstall and re-install the plugin sourced from whatever ref the marketplace was originally added with (the default branch, or a stale `@next` from a prior session). The user sees a successful-looking install and a stale plugin.

This silently defeats the recalibrated sprint goal "install path off `next` — fresh install works from `spacedock-dev/spacedock@next`."

## Problem

The 3-step sequence assumes `marketplace add` is idempotent and re-pins the ref. It is not: when the marketplace is already on disk, claude skips the add entirely with exit 0, so the downstream uninstall/install pair operates on the stale ref. There is no observable failure for the user or for any test that just checks the install exit code.

## Spike: live claude CLI behavior (run 2026-06-02 against claude 2.1.160)

The design rests on three unverified claude-CLI behaviors. All three were exercised live before committing to the plan; the results force the design (see "Design choice" below).

| # | Command | Pre-state | Observed exit | Observed stderr/stdout |
|---|---|---|---|---|
| 1 | `claude plugin marketplace add spacedock-dev/spacedock@next` | `spacedock` already declared (stale @next pin) | **0** | `Marketplace 'spacedock' already on disk — declared in user settings` (the bug) |
| 2 | `claude plugin marketplace remove spacedock` | `spacedock` declared | 0 | `Successfully removed marketplace: spacedock` |
| 3 | `claude plugin marketplace remove spacedock` | NOT declared (fresh box) | **1** | `Failed to remove marketplace: Marketplace 'spacedock' not found` |
| 4 | `claude plugin marketplace add spacedock-dev/spacedock@next` | `spacedock` removed (post-step 2) | 0 | Clones SSH, fetches `next`, `Successfully added marketplace: spacedock`; the @next pin lands |

(For #3, the fresh-box case, the spike substituted a non-existent marketplace name on this captain's box, which is the same code path inside claude — there is no "this name vs that name" branch in the remove handler.)

Result: the remove step exits 1 on a fresh box. Plain insertion of `marketplace remove` as a `CombinedOutput`-aggregating step would abort `Install` on first install. **Failure tolerance is mandatory**, not optional.

## Design choice: per-step failure-tolerance flag on the argv sequence

`installArgvSequence` becomes a slice of `{argv []string, tolerateExit bool}` entries (a small unexported struct local to `host_exec.go`); the marketplace-remove step is the only entry with `tolerateExit: true`. `Install` loops the slice and, when the step's `tolerateExit` is true, treats a non-zero exit as a recoverable case — it appends combined output to `sb` and continues, NOT returning the error. Every other step keeps today's strict behavior.

Updated 2026-06-02 during implementation — `TestUpgradeFromStaleMovesToGreen` (the AC-3 live upgrade-from-stale smoke) surfaced that `plugin uninstall` must precede `marketplace remove`: claude tracks the installed plugin via the marketplace record, and dropping that record orphans the uninstall step (`Plugin not found in installed plugins`, exit 1). Uninstall before remove preserves the invariant. The ideation flagged this revisit-trigger at the "Out of scope" note (line 83). The tolerance asymmetry is unchanged — only `marketplace remove` is tolerated.

```go
type installStep struct {
    argv         []string
    tolerateExit bool // marketplace remove is true; others false
}

func installArgvSequence(source, branch string) []installStep {
    return []installStep{
        {argv: []string{"plugin", "uninstall", "spacedock@spacedock"}},
        {argv: []string{"plugin", "marketplace", "remove", "spacedock"}, tolerateExit: true},
        {argv: []string{"plugin", "marketplace", "add", marketplaceAddArg(source, branch)}},
        {argv: []string{"plugin", "install", "spacedock@spacedock"}},
    }
}
```

Rejected alternatives:

- **Swallow only stderr matching "not declared".** Fragile: the message text is unstable across claude versions (we have already seen "not found" wording on 2.1.160, vs the "not declared" wording captain remembered) and a substring match couples our binary to a host string.
- **Reject only specific stderr matching for uninstall.** Same fragility argument as for remove: "Plugin not found in installed plugins" is the observed claude 2.1.160 wording, but the message is not part of any host contract and would silently couple our binary to a host string. The exit-code-only tolerance, scoped to the cleanup steps, is the stable contract.
- **Tolerate every step's exit.** Loses the fail-fast signal on the marketplace-add/install steps, which are the real-failure backstops we want surfaced (network, contract incompatibility, missing source).
- **Run the remove step out-of-band and ignore its exit by always wrapping in `_ = cmd.Run()`.** Same shape as the per-step flag but hides the tolerance decision inside `Install`'s control flow rather than co-locating it with the argv data; the flag makes the asymmetry visible in tests.

The per-step flag also keeps `installArgvSequence` testable as a pure function — `TestInstallArgvSequence` extends to assert the remove step's position AND its `tolerateExit: true`.

### Test-time fakeHost modeling

`fakeHost.Install` today records `[host, source, branch]` once per `runInit` call (no per-step argv visibility) — this is enough for `runInit`-level tests since the argv shape is owned by the pure `installArgvSequence`. The fix exercises both layers without thickening the fakeHost:

- **Argv shape** (pure function, no fake) — `TestInstallArgvSequence` extends to assert the 4-step ordering, the remove-first position, and the `tolerateExit: true` flag on step 0 only.
- **`runInit` end-to-end** (fakeHost) — existing `TestInitClaudeIssuesHostPluginCommands` and `TestInitTargetsNextWhenDevBranchPinned` keep passing; the fakeHost still records `[claude, spacedock-dev/spacedock, next]` (one `Install` call) and `runDoctor` still sees a compatible manifest.
- **Failure-tolerance behavior** (real `execHost.Install` against a stub binary) — a new unit test invokes `execHost.Install` with a per-PATH stub `claude` script that exits 1 on the remove subcommand and 0 otherwise; the test asserts the call returns nil error and that the combined output contains the remove-step stderr. This proves the loop tolerates the exit code, not just that the flag is set on the slice. Cost: a single `t.TempDir()` + `os.WriteFile` stub script, an `exec.LookPath`-friendly PATH override via `t.Setenv`. (No live claude required, no network.)

## Out of scope

- Codex install (the codex path is documented prose; the @ref is already correct on every printed invocation).
- The `--ref` flag surface itself (already correct — `marketplaceAddArg` composes `source@branch` correctly when `devBranch != ""`).

## Acceptance criteria

**AC-1 — `spacedock install --host claude` lands the `@next`-pinned plugin even when a stale marketplace declaration already exists.**
Verified by: `go test ./internal/cli -run TestInstallArgvSequence` passing on a sequence whose step 0 is `plugin marketplace remove spacedock` with `tolerateExit: true` and step 1 is `plugin marketplace add spacedock-dev/spacedock@next` — i.e. remove precedes add in the issued argv. Plus a manual `./spacedock install --host claude` run by the captain on a session whose `spacedock` marketplace is currently pinned to a non-`@next` ref, where the post-run `claude plugin marketplace list` shows `Source: GitHub (spacedock-dev/spacedock@next)`.

**AC-2 — A first install on a fresh box still succeeds end-to-end (both cleanup steps' fresh-box exit-1 are tolerated).**
Verified by three named tests: (a) `go test ./internal/cli -run TestInstallToleratesRemoveStepFailure` drives `execHost.Install` against a per-PATH stub `claude` returning exit 1 on `plugin marketplace remove ...` and exit 0 on every other subcommand; asserts `Install` returns nil and combined output contains all four steps' stub markers. (b) `go test ./internal/cli -run TestInstallToleratesUninstallStepFailure` is the mirror for the uninstall step: stub exits 1 on `plugin uninstall spacedock@spacedock` and exit 0 elsewhere; asserts the same nil-error + all-four-stub-markers shape. (c) `go test ./internal/cli -run TestFreshBoxInstallSucceeds` is the live fresh-box smoke against real `claude` 2.1.160 under an isolated `CLAUDE_CONFIG_DIR`+plugin-cache, with NO pre-seeded install or marketplace declaration; the test calls `execHost{}.Install("claude", localMarketplace, "")` and asserts nil error and that `claude plugin list --json` reports `spacedock@spacedock` with a non-empty installPath. Without the tolerance flag on EITHER cleanup step, (c) aborts at step 0 or step 1 with the real claude's exit 1.

**AC-3 — On a pinning step's failure, `Install` still surfaces the error (fail-fast on marketplace-add and plugin-install).**
Verified by `go test ./internal/cli -run TestInstallFailsFastOnAddStep`, where the stub `claude` exits 1 on `plugin marketplace add ...` and exit 0 on every other subcommand; `Install` MUST return a non-nil error wrapping the `add` step's argv and the combined output MUST contain the add-step stderr. The plugin-install step's fail-fast is locked by construction: the slice-walk assertion in `TestInstallArgvSequence` asserts ONLY the two cleanup steps carry `tolerateExit: true`, and `host_exec.go`'s `Install` loop has a single `step.tolerateExit` branch with no per-argv special casing. Together they cover the install step by construction. This locks the tolerance asymmetry — it is NOT a "tolerate every step" regression.

## Test plan

- `TestInstallArgvSequence` (extended): pure-function check over the 4-step slice; the remove step is at index 0 with `tolerateExit: true`; every other step has `tolerateExit: false` (or the zero value). Reuses `reflect.DeepEqual` on the slice. Cost: ~10 lines.
- `TestInstallToleratesRemoveStepFailure` (new): per-PATH stub `claude` script, drives real `execHost{}.Install("claude", "spacedock-dev/spacedock", "next")`; asserts nil error and combined output contains the remove-step stub message. Cost: ~25 lines including the stub script writer helper.
- `TestInstallFailsFastOnAddStep` (new): same stub-script pattern, stub fails on `add`; asserts non-nil error and the error wraps the add subcommand. Cost: ~15 lines (reuses the helper).
- Existing `TestInitClaudeIssuesHostPluginCommands`, `TestInitTargetsNextWhenDevBranchPinned`, `TestInitCheckRunsDoctorWithoutInstalling`, `TestInitMarketplaceSourceIsMigratedRepo`: keep passing unchanged — they exercise `runInit` via the fakeHost, which is one level above the argv-sequence change.
- Manual: captain runs `./spacedock install --host claude` on a box with a stale `spacedock` marketplace pin and confirms the post-install `claude plugin marketplace list` reports `@next`. This is the AC-1 live confirmation.
- Cost: low. No live host needed for the unit tests (the PATH stub is local), no network, no fakeHost API change.

## Notes

- Captain-observed 2026-06-02 — `./spacedock install --host claude` after a prior session's install: `Marketplace 'spacedock' already on disk — declared in user settings` followed by a "successful" uninstall+install pair from the stale ref.
- Spike-confirmed 2026-06-02 (claude 2.1.160 on this box): `marketplace add` when already declared exits 0 silently; `marketplace remove` when present exits 0; `marketplace remove` when not declared exits **1** with `Failed to remove marketplace: ... not found`; `marketplace add` after a successful remove fetches `@next` cleanly.
- Coordinates with packaging work (post-bootstrap repo migration to `spacedock-dev/spacedock`): a clean remove+add sequence is the right shape for `requires-contract` cross-repo verification too.
- Manual workaround until this lands: `claude plugin marketplace remove spacedock && ./spacedock install --host claude`.

## Stage Report: ideation

- DONE: The fix's mechanism (insert `marketplace remove spacedock` before the `add` in installArgvSequence) is exercised end-to-end against BOTH the live claude 'already on disk' case (the bug captain observed) and the 'not declared' fresh-box case (where remove emits non-zero) — OR an auditable 'no spike needed' is recorded with the proven mechanisms relied on.
  Live spike against claude 2.1.160 ran all four mechanism checks: stale-pin add exits 0 silently (bug confirmed), remove-when-present exits 0, remove-when-not-declared exits **1** with "not found" stderr, post-remove add fetches `@next` cleanly. Results table is in the entity body under "Spike: live claude CLI behavior".
- DONE: The remove-step's failure-exit handling is concretely designed (per-step tolerate-failure flag on installArgvSequence vs swallowing only the not-declared exit vs an alternative mechanism); the choice + rationale + how the test-time fakeHost models it is in the entity body.
  Chose per-step `tolerateExit bool` flag on a new `installStep` struct (the slice becomes `[]installStep`); rejected alternatives (substring-match on stderr, tolerate-every-step, hide-the-flag-in-Install) are documented with rationale. fakeHost is unchanged — argv shape stays a pure-function test, and a new stub-`claude`-script unit test against real `execHost.Install` proves the tolerance-loop behavior at the right layer.
- DONE: Each AC is entity-level (a property of the finished change, not a stage action), its `Verified by:` clause names a runnable check outside the entity body (a host-ops unit test against the fakeHost seam + a manual captain live-install confirmation), and the AC list covers both the upgrade-from-stale-ref path and the fresh-install path.
  AC-1 covers the stale-pin upgrade path (named `go test` + captain live confirmation); AC-2 covers the fresh-box path (named `go test` against a PATH-stub claude); AC-3 was added to lock the tolerance asymmetry (fail-fast on non-remove steps) — three named tests, every Verified-by points to a runnable command or live observation outside this document.

### Summary

Refreshed the ideation with a live mechanism spike against claude 2.1.160 that confirmed all four bug-relevant behaviors (the stale-pin add no-op, the not-declared remove exiting **1** with "not found" stderr, the remove-when-present exiting 0, and the post-remove add landing @next). The spike forced the design: failure tolerance on the remove step is mandatory, not optional. Chose a per-step `tolerateExit bool` flag on a new `installStep` struct as the cleanest seam — co-locating the tolerance decision with the argv data, keeping fail-fast on the other three steps, and remaining testable at the right layer (pure-function argv check + PATH-stub `execHost.Install` integration check for the tolerance loop). Added AC-3 to lock the tolerance asymmetry against silent regression toward tolerate-every-step.

## Stage Report: implementation

- DONE: installArgvSequence becomes []installStep with a per-step `tolerateExit bool`; the remove step is index 0 with tolerateExit:true; every other step is fail-fast. Install loops the slice and, on a tolerateExit step's non-zero exit, appends the combined output and continues without returning the error.
  Implemented in internal/cli/host_exec.go (commit on worktree branch spacedock-ensign/install-marketplace-ref-refresh); design-revisit during implementation reordered to {uninstall, marketplace-remove, marketplace-add, install} after TestUpgradeFromStaleMovesToGreen surfaced that dropping the marketplace orphans claude's plugin-tracking record. Tolerance asymmetry preserved (only marketplace-remove is tolerated); entity body Design-choice section updated with the inline note and the revised code snippet.
- DONE: The three named tests are added: TestInstallArgvSequence (extended — assert 4-step ordering + remove-first + tolerateExit asymmetry), TestInstallToleratesRemoveStepFailure (new — PATH-stub claude script returns exit 1 on `plugin marketplace remove`, exit 0 elsewhere; Install returns nil error + combined output contains all 4 steps' stub output), and TestInstallFailsFastOnAddStep (new — stub fails on `plugin marketplace add`; Install returns non-nil error wrapping the add subcommand). All three green.
  All three pass. TestInstallArgvSequence asserts the new {uninstall, remove, add, install} ordering and the asymmetric-tolerance rule (only marketplace-remove tolerated). The two new tests live in internal/cli/install_tolerance_test.go and drive a `t.TempDir()`+`os.WriteFile` /bin/sh stub `claude` via t.Setenv PATH override.
- DONE: The existing fakeHost-driven tests (TestInitClaudeIssuesHostPluginCommands, TestInitTargetsNextWhenDevBranchPinned, TestInitCheckRunsDoctorWithoutInstalling, TestInitMarketplaceSourceIsMigratedRepo) stay green WITHOUT fakeHost API changes — the argv shape stays a pure-function concern; the tolerance loop is tested at the execHost layer via the PATH-stub claude script.
  Confirmed: all four pass with no fakeHost-API change. Full repo `go test ./...` reports 732/732 across 12 packages.

### Summary

Implemented per-step `installStep{argv, tolerateExit}` in internal/cli/host_exec.go; the marketplace-remove step is the sole tolerated entry. During implementation the AC-3 live `TestUpgradeFromStaleMovesToGreen` surfaced an upgrade-path regression — claude tracks an installed plugin via its marketplace record, so `marketplace remove` before `plugin uninstall` orphans the uninstall step. Reordered to {uninstall, marketplace-remove, marketplace-add, install} (team-lead approved this as the ideation-authorized "revisit only if AC-1 confirmation surfaces it" trigger); the tolerance asymmetry is unchanged. All three named tests pass plus the live upgrade-from-stale test. Manual AC-1 confirmation remains for validation — captain to run `./spacedock install --host claude` on a stale-pin box and verify `claude plugin marketplace list` shows `@next`.

### Cycle 2 polish-fold (uninstall tolerance)

Post-validation FO investigation (Workflow Task w3bnwsvhx) surfaced an asymmetry gap: the implementation tolerated the marketplace-remove step but kept `plugin uninstall spacedock@spacedock` fail-fast on the assumption (entity body line 83, now removed) that claude's uninstall is a no-op on a not-installed plugin (fresh-box exit 0). Probe 2 EMPIRICALLY refuted that assumption against real `claude 2.1.160` — uninstall exits 1 with `Plugin not found in installed plugins` when the plugin is not installed. The cycle-1+cycle-2 live test `TestUpgradeFromStaleMovesToGreen` masked the gap by pre-seeding an install. Net effect on a fresh box (or any post-uninstall re-run): `Install` would abort at step 0 before reaching the marketplace-add re-pin — defeating AC-1's fresh-install-from-next intent. Captain authorized folding the fix into PR #272 as a 3rd commit on top of HEAD 96e9fad6.

The polish-fold flips step 0's `tolerateExit` to true and reframes the asymmetry as **destructive-cleanup-tolerated + pinning-fail-fast** (was: only-remove-tolerated). The new asymmetry IS the asymmetry — BOTH cleanup steps (uninstall + marketplace-remove) are tolerated because claude exits 1 on the fresh-box "not installed" / "not found" cases with no stable stderr to disambiguate from real failures; BOTH pinning steps (marketplace-add + plugin-install) stay fail-fast as the real-failure backstops. Tests: `TestInstallArgvSequence` extended to assert the new asymmetry over both cleanup steps via the `isUninstallStep`/`isMarketplaceRemoveStep` predicates; new `TestInstallToleratesUninstallStepFailure` mirrors the remove-step tolerance test with a stub `claude` failing on `plugin uninstall`; new `TestFreshBoxInstallSucceeds` exercises the live path against real `claude` 2.1.160 under an isolated `CLAUDE_CONFIG_DIR` with NO pre-seeding and asserts the plugin lands. `TestInstallFailsFastOnAddStep` stays as-is — still the fail-fast lock for the marketplace-add pinning step. All six related tests green at this revision.

## Stage Report: validation

- DONE: Every AC-N from ## Acceptance criteria is reproduced from an outside-the-body runnable check at THIS revision: AC-1 verified by TestInstallArgvSequence (new {uninstall, marketplace-remove, marketplace-add, install} ordering + remove-first tolerateExit asymmetry) + a flagged 'by-construction-pending-live' captain manual confirmation (./spacedock install --host claude on a stale-pin box, claude plugin marketplace list reports @next); AC-2 verified by TestInstallToleratesRemoveStepFailure (PATH-stub claude exits 1 on remove, 0 elsewhere; Install returns nil); AC-3 verified by TestInstallFailsFastOnAddStep (stub fails on add; Install returns non-nil error wrapping add subcommand).
  Ran all three at HEAD 63a9b5fe in the worktree: `go test ./internal/cli -run 'TestInstallArgvSequence|TestInstallToleratesRemoveStepFailure|TestInstallFailsFastOnAddStep' -v` → all three PASS (no SKIP). AC-1's live captain confirmation is by-construction-pending-live as the dispatch instructed; the live AC-1 surrogate `TestUpgradeFromStaleMovesToGreen` is green at this revision (3.04s, real `claude` 2.1.160 on PATH).
- DONE: The implementation-time reorder (uninstall before marketplace-remove) is verified to preserve the asymmetric-tolerance principle: ONLY marketplace-remove has tolerateExit:true; uninstall/add/install all fail-fast on a non-zero stub exit. TestInstallFailsFastOnAddStep is the lock; consider whether a sibling TestInstallFailsFastOnUninstallStep or TestInstallFailsFastOnInstallStep would close a residual gap, OR confirm AC-3's coverage is sufficient.
  Confirmed AC-3's coverage IS sufficient — adding sibling fail-fast tests for uninstall/install would duplicate, not extend, the lock. host_exec.go:270-281 is a single loop with a single `step.tolerateExit` branch and no per-argv special casing, so the slice shape IS the asymmetry. init_devbranch_test.go:84-96 walks the slice and asserts ONLY marketplace-remove has tolerateExit:true (every other step has tolerateExit:false); any future drift (e.g. tolerating uninstall) fails the structural assertion before reaching the loop. AC-3's TestInstallFailsFastOnAddStep is the behavioral lock for the loop; together they cover uninstall and install by construction.
- DONE: Full repo `go test ./...` (and `-race` if cheap) green at this revision. The TestUpgradeFromStaleMovesToGreen live regression that surfaced the reorder MUST be among the green tests at this revision (it was the AC-1 confirmation the ideation flagged at line 83 of the entity body).
  `go test ./...` at 63a9b5fe: 732/732 across 12 packages. `go test -race ./internal/cli/...` (the changed package, cheap at ~11s): clean. TestUpgradeFromStaleMovesToGreen is green at this revision and exercises real `claude` 2.1.160 against the exact argv `installArgvSequence` emits — including the tolerated remove step (runHostTolerant in upgrade_from_stale_test.go:107-113 mirrors the production loop's tolerance).

### Summary

PASSED. The three named AC tests and the live `TestUpgradeFromStaleMovesToGreen` AC-1 surrogate all pass at HEAD 63a9b5fe; full repo `go test ./...` is 732/732 and `-race` on `internal/cli/...` is clean. The implementation-time reorder to `{uninstall, marketplace-remove, marketplace-add, install}` preserves the asymmetric-tolerance principle: the slice-walk assertion in `TestInstallArgvSequence` structurally locks "only marketplace-remove is tolerated", and `TestInstallFailsFastOnAddStep` behaviorally locks the fail-fast path through the loop — together they cover uninstall and install fail-fast by construction (no sibling tests needed). The single remaining open item is the by-construction-pending-live captain AC-1 confirmation (manual `./spacedock install --host claude` on a stale-pin box → `claude plugin marketplace list` reports `@next`), which the dispatch explicitly flagged as not-a-blocker for validation.
