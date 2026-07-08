---
id: e6j9adxnn5hgv4hd7g5edr3t
title: "spacedock state ready resolves init paths relative to cwd and requires an origin remote"
status: validation
source: "GitHub issue #484 (spacedock-dev/spacedock#484), filed by clkao 2026-07-07, from a live split-root dogfooding session on this exact workflow shape (5 stages, ~6 entities, sd-b32 ids)."
started: 2026-07-07T22:59:51Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-state-ready-cwd-path-resolution
issue: spacedock-dev/spacedock#484
mod-block: merge:pr-merge
---

`spacedock state ready` has two cwd/remote-sensitivity defects observed while recovering a deleted split-root state checkout. (1) When invoked from inside an agent worktree (`.worktrees/<worker>-<entity>/`), its manual-fallback hint proposes re-adding the state checkout at a path nested UNDER that worktree instead of the project root — the command should resolve paths against the project root (or the workflow dir's declared location) independent of cwd, since hooks/scripts may invoke it from anywhere. (2) When the repo has no `origin` remote (deliberate local-only project), the command hard-fails on `git fetch origin <state-branch>` even though the state branch exists locally and a plain `git worktree add` from it succeeds — it should fall back to the local branch when the fetch fails or no remote is configured, matching `state commit`'s documented local-only behavior. Full repro, expected behavior, and the manual workaround used are in the linked issue.

## Problem statement

`spacedock state ready` is the boot gate that must converge a split-root state checkout from wherever it is invoked. Issue #484 (live dogfooding on 0.24.0-pre2) showed it failing to recover a deleted state checkout for three compounding reasons:

1. **cwd-relative anchoring.** State-checkout paths are computed as `filepath.Join(workflowDir, relPath)` (`internal/cli/state.go:50`, `internal/cli/state_sync.go:380`), and with no `--workflow-dir` the workflow dir is discovered from cwd. Invoked from inside an agent worktree, discovery lands on the worktree's copy of the workflow README (`git rev-parse --show-toplevel` inside a linked worktree is the worktree root, so the downward-scan fallback anchors there), and the resume path — plus the manual-fallback hint — targets `.worktrees/<w>/docs/<wf>/.spacedock-state`. That path is never correct: the split-root checkout is a singleton anchored at the main worktree, and the dir is gitignored so it never exists in worktree copies.
2. **Hard origin requirement.** The fresh-resume path fatals when `git fetch origin <state-branch>` fails (`state.go:64`), even when the orphan branch exists locally and `git worktree add` from it succeeds. A deliberately local-only repo (no `origin`) can never resume, contradicting `state commit`'s local-only carve-out and `state new`'s best-effort push.
3. **Stale registration blocks re-add.** In the reported repro the checkout directory was deleted but its worktree registration survived; `git worktree add` then fatals exit 128 ("missing but already registered worktree"), so fixes 1+2 alone still leave the reported scenario unrecoverable.

## Proposed approach

1. **Main-worktree anchoring for split-root checkouts.** One shared resolver maps (workflowDir, README `state:` relPath) → checkout path, routed through by the four CLI state verbs the #484 boot-gate exercises: `ready`/`commit` (replacing the join in `resolveStateCheckout`, `state_sync.go:380`) and `init`/`new` (replacing their own `filepath.Join(workflowDir, relPath)` at `state.go:50`/`135`). `state sweep` is deliberately NOT routed through it — it resolves via the separate `splitRootStateCheckout` (`internal/dispatch/helpers.go`) shared with `dispatch build`/`dispatch reconcile`, which stays out of scope (see below). When workflowDir sits inside a linked worktree (`git rev-parse --git-dir` ≠ `--git-common-dir`), the checkout resolves to `<main-worktree-root>/<repo-relative-prefix>/<relPath>` (prefix via `git rev-parse --show-prefix`; main root via the absolute common dir or the first `git worktree list --porcelain` entry — both spiked below). Otherwise — main worktree, or workflowDir not in a git repo at all (the standalone test fixtures) — today's join is unchanged. The re-anchor applies regardless of whether workflowDir came from discovery or an explicit `--workflow-dir`: the flag picks WHICH workflow, not where the singleton checkout lives; registration-path comparisons are realpath-normalized (macOS `/var` → `/private/var`).
2. **Resume fallback ladder** in `state init` (inherited by `state ready`'s absent-checkout resume): origin present → fetch; fetch OK → `git worktree add` (unchanged). No origin, or fetch failed → if the local state branch exists, `worktree add` from it and exit 0 — the message says local-only (no origin case) or warns the fetch failed and peers' state may be missing (stdout warning, matching `state new`'s push-warning style). Neither fetchable nor local → exit 1; the manual-fallback hint names the main-anchored path, and when the branch exists neither locally nor on origin, points at `spacedock state new` instead.
3. **Targeted stale-registration repair**, in the absent-checkout resume path only: if `git worktree list --porcelain` carries a registration for the (main-anchored, realpath-normalized) state path — necessarily stale, since the directory is absent — run `git worktree remove --force <statePath>` before the add and report the repair on stdout. Never `git worktree prune` (would touch unrelated registrations); never touch a registration whose directory exists (that remains the present-checkout no-op/refresh path).
4. **Doc diff** (below) recording cwd-independence and local-only resume in the split-root docs.

## Out of scope

- `state sweep` and `status --boot` STATE_BACKEND / dispatch-build state resolution from inside a worktree. These route through `internal/dispatch/helpers.go`'s `splitRootStateCheckout` (shared by `dispatch build`, `dispatch reconcile`, and `state sweep`) and `internal/status/roots.go` — an FO-invoked surface normally run from the main context, not the agent-worktree cwd the #484 boot-gate repro exercises. Re-anchoring `splitRootStateCheckout` cannot target `sweep` alone (the three callers share the one function), so it would drag `dispatch build` (every ensign dispatch) and `dispatch reconcile` into scope — materially larger than the reported defect; candidate follow-up if the cwd defect is observed there. `state sweep` is also read-only (no `worktree add`/`remove`, no commit), so AC-5's stale-registration repair is structurally N/A to it regardless.
- The present-checkout `state ready` pull failure with origin configured stays exit 1 (`state_sync.go:183`); only the absent-checkout resume gains the local fallback.
- README `state:`/`state-branch:` divergence between an agent branch and main: the resolved README's declared values are used as-is; editing those fields on a worktree branch is unsupported.
- General workflow-discovery semantics (walk-up / downward scan) are unchanged; only split-root checkout anchoring changes.

**Scope note for the gate:** fix 3 (stale-registration repair) is not one of the issue title's two defects, but the reported repro cannot converge without it — the issue's own workaround was exactly `worktree remove --force` + re-add. If struck at the gate, AC-1's setup must drop the stale registration.

## Acceptance criteria

- **AC-1 (value measure, independent baseline).** The issue-484 repro converges: in a repo whose split-root state checkout directory was deleted (orphan state branch and stale worktree registration intact) and which has no `origin` remote, `spacedock state ready` invoked with cwd inside an agent worktree exits 0 and restores the checkout at the main-root state path with its entity files present — where 0.24.0-pre2 exits non-zero and proposes a worktree-nested path. *Verified by:* a real-git e2e test reproducing the full setup, asserting exit code 0, entity file readable at the main-anchored path, and no `.spacedock-state` anywhere under `.worktrees/`.
- **AC-2.** The split-root checkout resolution used by `state ready`/`commit`/`init`/`new` is cwd-independent: the same main-worktree-anchored checkout resolves from the repo root, the workflow dir, an agent-worktree root, and the agent-worktree's workflow-dir copy; failure-path hints name the main-anchored path. (`state sweep` / `dispatch` resolution is explicitly out of scope — see Out of scope.) *Verified by:* e2e tests running `state ready` (and one `state commit`) from each cwd, asserting the main checkout is the one read/mutated and no nested checkout appears; one test forces the manual-fallback hint from a worktree cwd and asserts the hinted path.
- **AC-3.** A repo with no `origin` remote resumes an absent state checkout from the local state branch: `state ready`/`state init` exit 0, the checkout exists with its entities, and the output states local-only — the same carve-out family as `state commit`. *Verified by:* a no-origin e2e test asserting exit 0, on-disk entities, and the local-only wording.
- **AC-4.** With `origin` configured but unreachable, resume falls back to the local state branch (exit 0) and warns that the fetch failed; when neither a fetchable origin branch nor a local branch exists, the command exits non-zero with the main-anchored hint (and names `spacedock state new` when the branch was never birthed). *Verified by:* e2e tests with origin pointed at a nonexistent path, with and without the local branch.
- **AC-5.** A stale worktree registration for the state path (registered, directory missing) is repaired automatically during resume; a registration whose directory exists is never removed or re-added. *Verified by:* the AC-1 repro test plus a regression test asserting a present checkout keeps its worktree registration and HEAD across `state ready`.
- **AC-6.** The split-root docs state that the checkout resolves against the repository's main worktree regardless of invocation directory and that a no-origin repo resumes from the local state branch. *Verified by:* the doc diff below applied verbatim; static grep in review.

## Test plan

New real-git e2e tests in `internal/cli` beside `state_from_root_test.go` / `state_ready_test.go`, reusing the existing helpers (`twoHostStateWorkflow`, `commissionSplitWorkflow`, `writeSplitReadmeRepo`, `git(t, ...)`); the `run(...)` entry takes cwd explicitly, so worktree-cwd cases need no chdir:

1. `state ready` from an agent-worktree cwd with the main checkout present → resolves the main checkout, no nested dir (AC-2).
2. `state ready` from an agent-worktree cwd with the checkout absent → resumes at the main-anchored path (AC-2).
3. Full #484 repro: no origin + stale registration + worktree cwd → exit 0, checkout restored, entity present (AC-1, AC-5).
4. No-origin resume from repo root → exit 0, local-only message (AC-3).
5. Unreachable-origin fallback with local branch → exit 0 + warning; without local branch → exit non-zero, main-anchored hint (AC-4).
6. Never-birthed branch → hint names `spacedock state new` (AC-4).
7. Present-checkout regression: registration + HEAD untouched by `state ready` (AC-5).
8. `state commit` from a worktree cwd lands in the main checkout (AC-2 — exercises the shared resolver from a second verb, confirming it is not `state ready`-only; `state sweep` is not tested here, it is out of scope).

Cost: each test is an in-process real-git fixture like the existing suite (sub-second); no live workflow tests needed — every claim is command-level. Regression: full `go test ./...` must stay green, notably `TestStateCommitExplicitWorkflowDirOverridesDiscovery` (explicit flag still picks the workflow) and the standalone non-git-root fixtures (no re-anchor outside a repo).

## Spike record

Riskiest mechanisms exercised first in a throwaway repo (this session, git 2.39.5 / Apple Git-154); these seed tests 1-3:

- `git worktree add <path> <local-branch>` with **no origin remote** succeeds; checkout carries the branch's files.
- Deleted checkout dir with registration intact: `git worktree list --porcelain` marks it `prunable (gitdir file points to non-existent location)`; `git worktree add` at that path fatals exit 128 "missing but already registered worktree"; `git worktree remove --force <path>` exits 0 and a re-add converges.
- From `<main>/.worktrees/w-e/docs/dev`: `--git-dir` = `<main>/.git/worktrees/w-e` ≠ `--git-common-dir` = `<main>/.git` (both print `.git` from the main worktree — the linked-worktree discriminator); `--show-prefix` = `docs/dev/`; `git worktree list --porcelain` lists the main worktree first even from inside the linked worktree; `--show-toplevel` is the worktree root (confirming the defect's anchor). Caveat: bare `--git-common-dir` prints relative `.git` in the main worktree but absolute in a linked one — normalize via `--path-format=absolute` (supported on 2.39.5) or join before use.

## Documentation diff

`docs/site/advanced/split-root-state.md`, second paragraph — before:

> On a fresh clone the state checkout is absent; run `spacedock state init` to restore it.

After:

> On a fresh clone the state checkout is absent; run `spacedock state init` to restore it. The checkout always lands at the workflow's declared path under the repository's main worktree, wherever the command runs from; a repo with no `origin` remote resumes from the local state branch.

## Stage Report: ideation

- DONE: Problem statement, proposed approach, and out-of-scope boundary written into the entity body
  Three-part root-cause analysis (cwd anchoring at state.go:50/state_sync.go:380, hard fetch at state.go:64, stale registration), four-part approach, four out-of-scope items plus an explicit gate note on the third fix's scope.
- DONE: Entity-level acceptance criteria (each with how it's verified) plus a test plan added
  Six ACs, each with a verification clause; AC-1 measures the end-value (the #484 repro converges, exit 0 + on-disk state) against the 0.24.0-pre2 baseline that exits non-zero. Test plan names 8 real-git e2e tests beside the existing state suite with cost and regression scope.
- DONE: If the design rests on an unverified mechanism, the riskiest path is spiked first and recorded (or "no spike needed" with the proven mechanisms)
  Spiked in a throwaway repo before writing ACs: no-origin `worktree add`, stale-registration fatal/repair/re-add, and linked-worktree main-root detection (`--git-dir` vs `--git-common-dir`, `--show-prefix`, porcelain main-first) — recorded under "Spike record" with the relative-vs-absolute `--git-common-dir` caveat.

### Summary

Fleshed out issue #484 into a three-fix design: a shared main-worktree-anchored resolver for split-root checkout paths, a local-branch fallback ladder for resume, and targeted stale-registration repair — the third being required for the reported repro to converge (flagged for the gate since it exceeds the issue title's two defects). All risky git mechanisms were spiked first on git 2.39.5; a concrete doc diff for split-root-state.md is included per the user-visible-behavior rule.

### Feedback Cycles

**Cycle 1 (ideation gate, captain reject).** The Proposed approach's opening line claims "one shared resolver... used by every split-root state verb (`ready`/`commit`/`sweep` via `resolveStateCheckout`...)". This is factually wrong for `sweep`: `runStateSweep` (`internal/cli/state_sync.go:197-202`) never calls `resolveStateCheckout` — it delegates to `dispatch.Sweep`, which resolves its checkout path through a *third*, independent implementation, `internal/dispatch/helpers.go`'s `splitRootStateCheckout` (also used by `dispatch reconcile`). This directly contradicts the entity's own "Out of scope" line, which already excludes `internal/dispatch/helpers.go` as "a different surface, not part of the #484 repro." As scoped, the fix would leave `state sweep` cwd-broken from a worktree, so AC-2's "cwd-independent... every split-root state verb" claim would not actually hold post-implementation.

Routed back to ideation: either (a) extend the fix's scope to also cover `internal/dispatch/helpers.go`'s `splitRootStateCheckout`, or (b) correct the "every split-root state verb" claim to exclude `sweep` (consistent with the existing Out-of-scope line) and adjust AC-2 / AC-5 / the Test plan so no verified-by claim overclaims `sweep` coverage. Either resolution is acceptable; the gate blocks on the self-contradiction, not on a specific choice between them.

## Stage Report: ideation (cycle 2)

- DONE: Resolve the sweep / splitRootStateCheckout self-contradiction: either extend the fix's scope to cover internal/dispatch/helpers.go, or correct the "every split-root state verb" claim and out-of-scope sweep explicitly
  Chose (b). Verified in-tree that `runStateSweep` (state_sync.go:202) delegates to `dispatch.Sweep`, which resolves via `splitRootStateCheckout` (dispatch/helpers.go:123) — the SAME function backing `dispatch build` (build.go:476) and `dispatch reconcile` (reconcile.go:234); the three share one function, so `sweep` cannot be re-anchored alone. Proposed-approach point 1 rewritten to name the four verbs actually covered (`ready`/`commit` via resolveStateCheckout, `init`/`new` via their own joins) and state `sweep` is deliberately excluded; the Out-of-scope item now names `state sweep` with the shared-function rationale.
- DONE: Re-check AC-2, AC-5, and the Test plan so no verified-by claim overclaims sweep coverage under whichever resolution is chosen
  AC-2 statement now scopes cwd-independence to `state ready`/`commit`/`init`/`new` with a parenthetical excluding `sweep`/`dispatch`; its verified-by already tested only `ready`+`commit`. AC-5 is scoped to the resume path (`ready`/`init`) and its verified-by cites only the AC-1 repro + present-checkout regression — no sweep; confirmed sweep is read-only so registration repair is structurally N/A (recorded in Out-of-scope). Test plan item 8 tightened; no test item claims `sweep`.

### Summary

Resolved the cycle-1 gate rejection via option (b): `state sweep` is now explicitly out of scope, consistent with the pre-existing out-of-scope exclusion of `internal/dispatch/helpers.go`. Confirmed by reading the code that `sweep`/`dispatch build`/`dispatch reconcile` all resolve through the one shared `splitRootStateCheckout`, so re-anchoring it (option a) would expand the blast radius far beyond the #484 boot-gate defect — YAGNI favors scoping out. No verified-by claim now overclaims sweep; the #484 repro (AC-1) and its cwd/no-origin/stale-registration fixes are unchanged.

## Stage Report: implementation

- DONE: Shared main-worktree-anchored resolver implemented for ready/commit/init/new (sweep explicitly excluded per Out of scope) per the Proposed approach
  `internal/cli/state_checkout.go` adds `resolveSplitRootCheckout`, `isLinkedWorktree` (`--path-format=absolute --git-dir` vs `--git-common-dir`, realpath-normalized), and `mainWorktreeRoot` (`git worktree list --porcelain` first entry). Routed through by `resolveStateCheckout` (state_sync.go, backs `ready`/`commit`) and both `statePath` computations in state.go (`init`/`new`). Mechanisms independently spiked in a throwaway repo (git 2.39.5) before coding, confirming the ideation spike record's git-dir/common-dir/porcelain-ordering claims.
- DONE: Resume fallback ladder (no-origin / fetch-failed to local branch) and targeted stale-registration repair implemented
  New `resumeAbsentSplitRootCheckout` (state.go) — shared by `state init` and `state ready`'s absent-checkout path — repairs a stale `git worktree` registration for the absent statePath, then tries fetch→local-branch→fail, warning on fetch-failure fallback and disambiguating "never birthed" (names `state new`) from "origin unreachable, indeterminate" (generic manual-fallback hint) via a new `remoteBranchStatus` helper (`ls-remote` reachability vs. found-ness). `state ready` also gates its own post-resume `pull --rebase` on whether the resume actually reached origin, fixing a second-order bug the AC-4 test caught: an unreachable-but-configured origin was making `state ready` fail a redundant pull immediately after a successful local-branch fallback.
- DONE: All named e2e tests added and green, including the AC-2 worktree-cwd manual-fallback-hint test and disambiguation of which cwd variant each test covers (both flagged at the ideation gate)
  `internal/cli/state_checkout_test.go`, 9 tests covering test-plan items 1-8 plus the flagged hint test: present/absent checkout from a worktree cwd (bare-discovery variant), the full #484 repro, no-origin resume from root, unreachable-origin with/without a local branch, never-birthed-from-worktree-cwd (explicit-`--workflow-dir`-into-worktree variant, satisfying the flagged AC-2 hint test), present-checkout registration/HEAD regression, and `state commit` from a worktree cwd. Each test's doc comment names its cwd variant. Mutation-tested two: disabling the stale-registration repair reproduces git's exact exit-128 fatal (confirming the test is load-bearing, not vacuous); forcing `neverBirthed=true` flips the unreachable-origin test's hint to the wrong wording (confirming the disambiguation is exercised). Full `go test ./...` and `golangci-lint` stay green; two lint findings in the new code (ineffassign, unchecked `filepath.Walk` error) were fixed, not suppressed.

### Summary

Implemented the three-fix design from ideation: a shared main-worktree-anchoring resolver, a no-origin/unreachable-origin local-branch fallback ladder, and stale-worktree-registration repair, plus the documented split-root-state.md diff. Found and fixed one gap the design didn't anticipate — `state ready`'s post-resume pull needed to know whether the resume reached origin, or it would immediately re-fail the pull after a successful local fallback — via a small refactor extracting the resume logic into `resumeAbsentSplitRootCheckout` shared by `state init` and `state ready`. All 8 test-plan items plus the ideation-flagged worktree-cwd hint test are green; full suite and lint are clean.

## Stage Report: validation

- DONE: Run go test ./... and golangci-lint; reproduce the reported test results and the mutation-tested claims yourself rather than trusting the stage report
  `go test ./... -count=1` green across all 15 packages (internal/cli 88.4s, full repo ~4.5min). All 9 new tests in `internal/cli/state_checkout_test.go` independently re-run and pass. `golangci-lint run ./...` panics ("file requires newer Go version go1.26, application built with go1.25") — reproduced identically on unmodified `main` (570de040), so it is a pre-existing local toolchain/binary mismatch, not caused by this diff. Per-package lint (`./internal/cli/...`, `./internal/status/...`, and a loop over every `go list` package individually) runs clean of panics on both `main` and the worktree; the finding sets are byte-identical between the two (22 pre-existing errcheck/staticcheck issues in `internal/cli`, 37 combined with `internal/status`), and none touch `state.go`, `state_checkout.go`, `state_checkout_test.go`, or `status/path.go` — confirming this diff introduces zero new lint findings.
- DONE: Verify AC-1 through AC-6 with reproduced evidence, including the flagged AC-2 worktree-cwd hint test and cwd-variant disambiguation
  AC-1 (`TestStateReadyIssue484Repro`): exit 0, entity restored, no `.spacedock-state` under `.worktrees/` — passed. AC-2: cwd-independence across bare-discovery-from-worktree-cwd (`TestStateReadyFromWorktreeCwd{Present,Absent}CheckoutResolvesMain`, `TestStateCommitFromWorktreeCwdLandsInMainCheckout`) and the flagged explicit-`--workflow-dir`-into-worktree variant (`TestStateReadyNeverBirthedFromWorktreeCwdHintsStateNewAtMainPath`) — two distinct cwd shapes, both passed, both assert no nested checkout. AC-3 (`TestStateReadyNoOriginResumeFromRoot`): exit 0, "local-only" in stdout — passed. AC-4: unreachable-origin-with-local-branch warns+exit 0 (`TestStateReadyUnreachableOriginFallsBackToLocalBranch`), unreachable-without-local hints main-anchored path without "state new" (`TestStateReadyUnreachableOriginNoLocalBranchHintsMainAnchoredPath`), never-birthed hints "state new" (`TestStateReadyNeverBirthedFromWorktreeCwdHintsStateNewAtMainPath`) — all passed. AC-5: repro test plus `TestStateReadyPresentCheckoutRegressionUntouched` (worktree registration + HEAD untouched on a present checkout) — passed. AC-6: read `docs/site/advanced/split-root-state.md` directly — the shipped sentence matches the entity's Documentation diff "After" text verbatim.
- DONE: Detached adversarial audit on a throwaway checkout (never the implementation worktree)
  Cloned the worktree's branch to a scratch dir outside the repo tree. Mutation 1 (disable `repairStaleWorktreeRegistration` in `resumeAbsentSplitRootCheckout`): reproduced git's exact exit-128 "missing but already registered worktree" fatal; 3 tests went red (`TestStateReadyIssue484Repro`, `TestStateReadyNoOriginResumeFromRoot`, `TestStateReadyUnreachableOriginFallsBackToLocalBranch` — the latter two hit it too because their fixtures simulate checkout loss via `os.RemoveAll`, which itself leaves a stale registration). Mutation 2 (drop the `!originReached` gate in `runStateReady` so the post-resume pull always runs): `TestStateReadyUnreachableOriginFallsBackToLocalBranch` went red with the exact redundant-pull-after-fallback failure the implementation report described. Mutation 3 (force `neverBirthed := true` unconditionally): `TestStateReadyUnreachableOriginNoLocalBranchHintsMainAnchoredPath` went red, catching the wrong "state new" wording on an indeterminate (not never-birthed) origin. All three confirm the suite is load-bearing, not vacuous. Reverted all mutations and deleted the scratch checkout.

### Summary

Reproduced every claim in the implementation report independently rather than trusting it: full test suite green, all 9 new tests re-run and passed, all 6 ACs verified against their cited tests plus a verbatim doc-diff read, and three adversarial mutations on a throwaway checkout (stale-registration repair, origin-reachability gate, never-birthed disambiguation) each drove a distinct, predicted test red. The only deviation from the report is `golangci-lint run ./...`, which panics in this environment on a pre-existing go1.25/go1.26 toolchain mismatch reproduced identically on unmodified `main` — per-package lint confirms zero new findings from this diff. PASSED.
