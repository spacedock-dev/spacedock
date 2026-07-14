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
mod-block:
pr: "#503"
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

**Cycle 2 (post-validation Roborev review, job 487).** The independent
thorough branch review rejected exact head
`9c24747050250c60b5807bcd7bb89bbbfae28f28` with five material findings:

- `internal/cli/state_checkout.go:22`: main-worktree anchoring reaches only the
  selected state subcommands, so `status` and `state sweep` can still target a
  wrong checkout.
- `internal/cli/state.go:457`: concurrent resume can race through destructive
  `git worktree remove --force` behavior.
- `internal/cli/state.go:217`: `state new` from a linked worktree edits the wrong
  `.gitignore`, leaving the main checkout dirty.
- `internal/cli/state_checkout.go:23`: metadata failure falls back to a
  worktree-local path instead of failing closed.
- `internal/cli/state_sync.go:157`: `state ready --json` can emit prose before
  the JSON document.

Routed back to implementation: close each repository-boundary, concurrency,
cleanliness, fail-closed, and output-atomicity defect with an adversarial
behavioral regression. Re-run full and race gates, then replace Roborev job 487
at the corrected exact head before the validation gate can be presented again.

**Cycle 3 (replacement validation, Roborev job 707; captain escalation).** The
five cycle-2 repairs, focused tests, full suite, race suite, and detached
lock-removal audit passed at exact PR head
`3d50bd9af2dd31e4b0d5ad6bab09df0f7459e535`. Authoritative corrected-guideline
Roborev job `707` nevertheless rejected the exact
`557f8df3e6a62d34987edda70533375fc48ba8f6..3d50bd9af2dd31e4b0d5ad6bab09df0f7459e535`
range:

- HIGH: line-delimited porcelain parsing breaks on a real main-worktree path
  containing a newline; use NUL-delimited complete-record parsing.
- HIGH: the first porcelain record can be a bare repository; fail closed unless
  a usable primary worktree is identified.
- MEDIUM: the resume lock ends before origin-backed `pull --rebase`, leaving
  reachable/unreachable-origin concurrency non-atomic.
- MEDIUM: `state init` fetches before restoring an existing local branch but can
  return without integrating the fetched remote state.

Required live Claude Sonnet, Claude Opus, Codex, and Pi jobs also remain WAITING
and are not green. Per the three-cycle feedback rule, this cycle is escalated to
the captain instead of being automatically routed to another implementation
worker. PR #503 must not merge in its current state.

Captain decision: send cycle 3 back for another implementation round. Fix and
test locally without pushing a new head. The next fresh validator must run the
corrected exact-range Roborev review before any push, CI trigger, deployment
approval, or live-lane engagement. Only a Roborev PASS may release the exact
head to the PR and CI; any deliberately unhandled failure mode must first be
explicitly scoped out in the task and validation evidence.

**Cycle 4 (Roborev-first validation, job 775).** Authoritative exact-range
review rejected local head `d346aba4dffd706918e8b0f746ca1be78ba7c9c6`
before tests, push, or CI:

- MEDIUM: waiters can return success after the creator leaves a post-creation
  convergence failure; clean up or propagate the failed state and test the
  creator-failure concurrency path.
- MEDIUM: the non-Unix resume lock is process-local, so cross-process repair and
  creation remain racy; provide a real cross-process guarantee or fail
  explicitly with behavioral proof.
- LOW: shared exclude is mutated for main-worktree births; limit the persistent
  ignore rule to linked-worktree placement.

Routed back to implementation under the captain's standing instruction: no push
or CI before a fresh corrected exact-range Roborev PASS. Any finding not fixed
must be explicitly scoped out with evidence before the next review.

**Cycle 5 (Roborev-first validation, job 867).** Authoritative thorough review
stored exact
`557f8df3e6a62d34987edda70533375fc48ba8f6..677215842a369850fcefbf239851da69a3389d2b`,
included all four corrected guideline sections, and rejected before tests, push,
PR update, or CI. Earlier job `862` is non-authoritative because its positional
range included the base commit.

- HIGH: failed-resume cleanup force-removes the published checkout path even
  though `state new`, state writers, and direct Git operations do not honor the
  resume lock. A concurrent valid checkout or uncommitted state can be deleted.
  Converge privately before publication, or prove exclusive ownership and
  unchanged state before cleanup, with cross-command writer coverage.
- MEDIUM: one repository-wide durable outcome file lets workflow B overwrite
  workflow A's result before A's waiter reads it. Key outcomes by canonical
  checkout path and cover concurrent resumes of two workflows.

Routed back to implementation under the captain's standing instruction. These
are concurrency and data-loss defects inside the existing task scope; no
reframing is needed. No push or CI is permitted before a fresh corrected
exact-range Roborev PASS. A deliberately unhandled case must be explicitly
scoped out with evidence before re-review.

**Cycle 6 (Roborev-first validation, job 904).** Authoritative thorough review
stored exact
`557f8df3e6a62d34987edda70533375fc48ba8f6..4cdb53c27ed8165ee425e2bca36af316caff3892`,
included all four corrected guideline sections, and rejected before tests, PR
update, or CI:

- HIGH: existing state directories are trusted without proving they are the
  exact registered worktree on the expected state branch, allowing Git to climb
  and mutate the main code branch.
- HIGH: resume publishes the final checkout before origin convergence, so
  non-locking state writers can lose concurrent entity changes.
- MEDIUM: an external symlink into a linked worktree is misclassified as a
  standalone repository.
- MEDIUM: `strings.TrimSpace` corrupts valid leading or trailing whitespace in
  Git-reported paths and prefixes.

Routed back to implementation under the captain's standing instruction. These
are existing-scope repository identity, publication-race, symlink, and exact-path
defects; no reframing is needed. No PR update or CI is permitted before a fresh
corrected exact-range Roborev PASS.

**Cycle 7 (Roborev-first validation, job 951).** Authoritative thorough review
stored exact
`557f8df3e6a62d34987edda70533375fc48ba8f6..f57b0924523ee75ef0a2c08f9c999d9c6b0756a1`,
included all four corrected guideline sections, and rejected before tests, PR
update, or CI:

- MEDIUM: repository-wide `git worktree prune --expire now` can delete unrelated
  temporarily missing worktree registrations while repairing one state path.
- MEDIUM: a waiter can observe the published checkout before lock release, skip
  the completed resume outcome, and issue a redundant pull that turns a
  successful local fallback into failure when origin is unreachable.
- MEDIUM: checkout-path-only durable outcomes survive deletion and independent
  recreation, allowing a new checkout generation to inherit stale `ready`.

Routed back to implementation under the captain's standing instruction. These
are existing-scope targeted cleanup, publication ordering, and generation
identity defects; no reframing is needed. No PR update or CI is permitted before
a fresh corrected exact-range Roborev PASS.

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

## Stage Report: implementation (cycle 2)

- DONE: Close all five Roborev job 487 defects with focused adversarial regressions covering wrong-root reach, concurrent resume, linked-worktree .gitignore cleanliness, metadata failure, and atomic JSON output.
  Commit `3d50bd9a` moves repository placement into `internal/status` for CLI/status/dispatch, adds a repository-scoped resume lock with an in-lock existence recheck, protects the main physical checkout through the shared Git exclude, fails closed on corrupt worktree metadata, and suppresses resume prose under `--json`; linked-root read/write/sweep, 12-caller race, exact registration/HEAD, three JSON resume variants, and cleanup-state regressions pass.
- DONE: Preserve compatibility-first output and split-root semantics while proving failures occur before destructive or cross-checkout mutation.
  Existing prose mode and state paths remain stable; the metadata-failure regression asserts exit 1 before nested-path creation or main-state HEAD mutation, the race asserts one registration and preserved entity bytes/HEAD, and the linked-worktree birth regression proves the main checkout stays clean while the invoking branch receives the tracked `.gitignore` edit.
- DONE: Run gofmt plus focused, full, and race gates; commit the repair on the existing branch and leave PR #503 ready for a replacement exact-head Roborev review.
  Ran `gofmt -w ./cmd ./internal`, focused adversarial tests, `go test ./... -count=1`, and `go test ./... -race -count=1` successfully; pushed exact head `3d50bd9a` to `origin/spacedock-ensign/state-ready-cwd-path-resolution`, updating PR #503 without any local merge.

### Summary

Closed all five findings from Roborev job 487 on the existing PR branch. Repository-boundary resolution is now one fail-closed source shared by status, dispatch, and state verbs; resume is serialized across processes; JSON stdout is atomic; and the expanded real-Git suite checks exact roots, refs, output bytes, registrations, failure ordering, and cleanup state before replacement review.

## Stage Report: implementation (cycle 3)

- DONE: Replace line-delimited worktree parsing with complete NUL-delimited records that handle newline paths and reject bare-only or unusable primary-worktree candidates.
  Local commit `d346aba4` requests `git worktree list --porcelain -z` and shares a complete-record parser across placement, stale-registration repair, and boot orphan scanning; byte-preservation, real newline-root, truncated-record, synthetic bare-primary, and real bare-repository-with-linked-worktree regressions pass, while placement also rejects non-absolute, missing, non-worktree, prunable, and common-Git-dir candidates.
- DONE: Keep resume synchronization atomic through origin-backed pull and make restored local branches integrate fetched remote state, with reachable/unreachable concurrency regressions.
  `state init` and `state ready` now hold the repository-scoped lock from the absent/present observation recheck through fetch, add, pull/rebase or fallback, output, and conflict cleanup; an 8-caller boundary barrier proves reachable callers all return only after the peer HEAD is present and unreachable-origin waiters do not re-pull, while a separate `state init` regression verifies exact peer bytes and HEAD after restoring a stale local branch.
- DONE: Run focused, full, and race tests and commit locally without pushing; leave an exact clean head for fresh Roborev-first validation, with any excluded failure mode made explicit.
  Focused adversarial tests pass; `go test ./... -count=1 -p 1` and `go test ./... -race -count=1 -p 1` pass at clean head `d346aba4dffd706918e8b0f746ca1be78ba7c9c6`. Package-parallel full/race runs intermittently hit the unrelated pre-existing `TestSonnetTeamDeleteHangReplay` replay-order flake, which passed 10/10 non-race and on isolated race rerun. The code branch remains one local commit ahead of remote head `3d50bd9a`; it was not pushed and CI was not triggered or approved.

### Summary

Closed all four findings from corrected-guideline Roborev job 707 without excluding a failure mode. Worktree identity now survives arbitrary path bytes and rejects non-working-tree primaries, while resume creation and remote convergence form one serialized transaction for both `state init` and `state ready`; the exact clean local head is ready for the required fresh Roborev-first validation before any push or CI action.

## Stage Report: validation (cycle 2)

- DONE: Independently reproduce all five Roborev job 487 repairs and every AC against exact head 3d50bd9a; verify failures precede cross-checkout or destructive mutation and JSON output remains atomic.
  Exact head `3d50bd9a` was clean and matched PR #503. Focused real-Git tests passed for wrong-root status/sweep, 12-caller resume, linked-worktree `.gitignore` cleanliness, fail-closed metadata, all three absent-checkout JSON variants, and AC-1 through AC-5; AC-6's shipped documentation sentence matches the specified diff. The failure-order tests assert no nested checkout or main-state HEAD mutation, exact registration/HEAD/entity preservation, clean main status, and one decodable JSON document.
- DONE: Run gofmt verification, focused adversarial tests, go test ./..., go test ./... -race, and the required detached adversarial audit without modifying the implementation worktree.
  `gofmt -l ./cmd ./internal` returned empty, `git diff --check` passed, the focused suite passed, and `go test ./... -count=1` plus `go test ./... -race -count=1` passed. In detached scratch checkout `/tmp/e6j-detached-audit.hZbeYP/repo`, bypassing the repository lock made `TestConcurrentStateReadySerializesRepairAndCreation` fail in all 5 runs during competing stale-registration removal; the implementation worktree remained clean.
- FAILED: Run a replacement thorough Roborev review for exact base 557f8df3 and head 3d50bd9a under the corrected repository guideline, then report PR #503 head/CI state and a PASSED or REJECTED recommendation.
  Authoritative Roborev job `707` stored exact range `557f8df3e6a62d34987edda70533375fc48ba8f6..3d50bd9af2dd31e4b0d5ad6bab09df0f7459e535`, thorough reasoning, corrected Compatibility/Trust/Behavioral-proof/Review-focus guideline, and verdict FAIL. PR #503 remains OPEN/MERGEABLE at that exact head; offline, docs build, and Ubuntu/macOS install checks pass, but Claude Sonnet, Claude Opus, Codex, and Pi live jobs remain WAITING and are not green.

### Reviewer findings

- HIGH: `internal/status/repository_placement.go:64` parses line-delimited porcelain paths; a real main-worktree path containing a newline split across records, proving it can anchor state under the wrong root. Use `--porcelain -z` and NUL parsing.
- HIGH: `internal/status/repository_placement.go:68` accepts the first worktree record without its `bare` marker; a real bare repo listed the bare Git directory first, proving state can be placed in administration storage. Parse complete records and fail closed without a usable primary worktree.
- MEDIUM: `internal/cli/state_sync.go:164` releases the resume lock before origin-backed `pull --rebase`; waiters also treat configured origin as reached. Add reachable- and unreachable-origin concurrency coverage and keep synchronization atomic.
- MEDIUM: `internal/cli/state.go:103` fetches before `state init` restores an existing local branch but does not integrate the fetched remote state, so successful init can return a stale checkout.

### Summary

The five job-487 repairs and their stated AC proofs pass at exact head `3d50bd9a`, and the race test is independently mutation-proven. Replacement review exposed four adjacent repository-path and synchronization defects, while required live CI remains unapproved, so the validation recommendation is **REJECTED**; do not merge PR #503.

## Stage Report: validation (cycle 3)

- FAILED: Run corrected-guideline Roborev first on exact range 557f8df3..d346aba4 and inspect stored range/prompt; on any finding stop REJECTED without pushing PR #503 or engaging CI.
  Authoritative thorough Roborev job `775` stored exact range `557f8df3e6a62d34987edda70533375fc48ba8f6..d346aba4dffd706918e8b0f746ca1be78ba7c9c6`, all four corrected guideline sections, and verdict FAIL. It found a false-success waiter after failed checkout convergence, non-Unix cross-process locking reduced to an in-process mutex, and an over-broad shared-exclude mutation.
- SKIPPED: Only after Roborev PASS independently reproduce newline/bare-worktree parsing and atomic resume convergence, verify every AC, gofmt, focused/full/race suites, and a detached adversarial mutation audit.
  Roborev failed the first gate, so the assignment required validation to stop before these downstream checks.
- SKIPPED: Only after all local gates pass push exact head d346aba4 to PR #503 and then engage required CI; report actual check outcomes, treating waits/skips/flakes as non-green until this-run evidence resolves them.
  No code was pushed and no CI job was triggered or approved; PR #503 remains on remote head `3d50bd9a` while exact reviewed head `d346aba4` remains local-only.

### Reviewer findings

- MEDIUM: a waiter that observed an absent checkout can return success after a creator leaves a post-creation `pull --rebase` failure behind; remove the failed checkout/registration or propagate convergence failure, and add creator-failure concurrency coverage.
- MEDIUM: `internal/cli/state_resume_lock_other.go` provides only an in-process mutex on `!unix`, so separate processes can race repair and creation; provide a cross-process lock or fail explicitly and test separate processes.
- LOW: shared `.git/info/exclude` is mutated even for main-worktree births; gate it on linked-worktree placement so normal births do not leave a persistent hidden ignore rule.

### Summary

Validation is **REJECTED** at the Roborev-first gate. Exact-range job `775` found three unscoped failure modes, so no tests beyond review, PR update, or CI engagement were performed; these findings must return to implementation before another fresh validation.

## Stage Report: implementation (cycle 4)

- DONE: Make creator post-creation convergence failure atomic for all waiters: clean up or propagate the failure, with deterministic concurrent creator/waiter regressions.
  Local commit `67721584` adds a repository-common-dir resume outcome (`pending`/`failed`/`ready`) keyed to the real state path. Creators publish `ready` only after full convergence; on post-creation failure they publish `failed` and remove the checkout/registration while preserving the local branch ref. `TestConcurrentStateReadyCreatorPullFailureNeverLooksReady` releases six callers from the same absent-checkout boundary into a deterministic pull/rebase conflict and proves every caller exits 1, the failed checkout and worktree registration are absent, and the local branch HEAD is preserved.
- DONE: Provide a real non-Unix cross-process resume guarantee or an explicit fail-closed unsupported result, proven across separate processes rather than a mutex-only test.
  Replaced the `!unix` process-local mutex with explicit `errStateResumeLockUnsupported`; its callback is never invoked. New child-process tests prove the exact unsupported wrapper fails closed in two separate processes and prove the Unix `flock` path serializes two separate processes through a gate plus an `O_EXCL` active-section sentinel (green for three repetitions). `GOOS=windows GOARCH=amd64 go test -c -o /tmp/spacedock-cli-windows.test.exe ./internal/cli` also passes, proving the non-Unix path compiles.
- DONE: Limit shared exclude mutation to linked-worktree placement, run focused/full/race tests, and commit locally without push for another Roborev-first validation.
  Shared `.git/info/exclude` mutation is now gated on `placement.Linked`; `TestStateNewFromMainDoesNotPersistSharedExclude` proves a normal main-worktree birth leaves it untouched while the existing linked-worktree regression remains green. Focused adversarial tests, `git diff --check`, Windows cross-compilation, `go test ./... -count=1 -p 1`, and `go test ./... -race -count=1 -p 1` pass at clean head `677215842a369850fcefbf239851da69a3389d2b`. A package-parallel full run hit the pre-existing `TestSonnetTeamDeleteHangReplay` replay-order flake; this run's serialized full and race suites both executed that package and passed. The code branch remains two local commits ahead of remote `3d50bd9a`; no code push or CI trigger/approval occurred.

### Summary

Closed all three Roborev job 775 findings without excluding their failure modes: waiter success now requires a durable matching convergence result, unsupported locking fails closed across processes, and shared exclude mutation is linked-worktree-only. The exact clean local head is fully verified and remains unpushed for the required fresh Roborev-first validation.

## Stage Report: validation (cycle 4)

- FAILED: Run corrected-guideline Roborev first on exact range 557f8df3..67721584 and inspect stored range/prompt; on any finding stop REJECTED without push or CI.
  Authoritative thorough Roborev job `867` stored exact range `557f8df3e6a62d34987edda70533375fc48ba8f6..677215842a369850fcefbf239851da69a3389d2b`, all four corrected guideline sections, and verdict FAIL; earlier job `862` is non-authoritative because positional range semantics included the base commit.
- SKIPPED: Only after Roborev PASS independently reproduce convergence-failure fan-out, cross-process Unix serialization and non-Unix fail-closed behavior, linked-only exclude mutation, every AC, gofmt, focused/full/race suites, and detached audit.
  Roborev failed the first gate, so the assignment required validation to stop before every downstream local check.
- SKIPPED: Only after all local gates pass push exact head 67721584 to PR #503 and then engage required CI; report actual outcomes and block on waiting, skipped, failed, or unapproved evidence.
  No code was pushed and no CI was triggered or approved; PR #503 remains on remote head `3d50bd9a` while exact reviewed head `67721584` remains local-only.

### Reviewer findings

- HIGH: failed-resume cleanup force-removes the final checkout path although `state new`, state writers, and direct Git operations do not honor the resume lock; a concurrent valid checkout or uncommitted state can be deleted. Converge privately before publication, or prove exclusive ownership and unchanged state before cleanup, with cross-command writer coverage.
- MEDIUM: one repository-wide outcome file lets workflow B overwrite workflow A's result before A's waiter reads it, causing A to reject a successfully restored checkout. Key durable outcomes by canonical checkout path and cover concurrent resumes of two workflows.

### Summary

Validation is **REJECTED** at the Roborev-first gate. Exact-range job `867` found two unscoped concurrency and data-loss defects, so local tests, PR update, and CI engagement were not performed; return the findings to implementation before another fresh validation.

## Stage Report: implementation (cycle 5)

- DONE: Prevent failed-resume cleanup from deleting a concurrently valid checkout or uncommitted state, using private convergence or proven ownership plus cross-command writer coverage.
  Local commit `4cdb53c2` removes every failed-resume deletion of the public checkout and keeps failure fan-out in the per-resume outcome; stale registrations now recover with Git's same-path `worktree add --force`, which replaces absent metadata but refuses an occupied checkout. Deterministic `state ready` and `state init` races preserve a concurrent `state new` birth, a different direct-Git worktree plus exact uncommitted bytes, and an uncommitted `status --set` mutation after pull/rebase failure; temporarily restoring forced cleanup drove all three ready-path tests red.
- DONE: Key durable resume outcomes by canonical checkout path and prove simultaneous resumes for two workflows cannot overwrite or misread one another.
  Outcome filenames now use the SHA-256 of the realpath-normalized checkout under the repository common Git dir while retaining the canonical path inside the JSON record. Four callers simultaneously resume two absent workflows in one repository and all return ready with distinct durable results; overwriting workflow A's outcome leaves B ready, and mutating the implementation back to one repository-wide file makes the regression fail.
- DONE: Run focused/full/race gates and commit locally only; a fresh validator must run corrected exact-range Roborev before any PR update or CI engagement.
  Changed-file gofmt and `git show --check` are clean; Windows cross-compilation produced a PE32+ test executable; focused lifecycle tests, full `internal/cli`, `go test ./... -count=1`, and `go test ./... -race -count=1` pass at clean head `4cdb53c27ed8165ee425e2bca36af316caff3892`. The branch is three local commits ahead of unchanged remote `3d50bd9a`; no code push, PR update, CI trigger, or CI approval occurred.

### Summary

Closed both job-867 findings without deleting or overwriting any public checkout: failed convergence is durable and fail-closed, stale metadata recovery is non-destructive, and every non-locking adjacent writer path is exercised for both resume verbs. Outcomes are independently keyed per canonical checkout, and the exact clean local head is ready for corrected exact-range Roborev-first validation.

## Stage Report: validation (cycle 5)

- FAILED: Run corrected-guideline Roborev first on exact range 557f8df3..4cdb53c2 and inspect the stored range and all four guideline sections; on any finding stop REJECTED without tests, PR update, or CI.
  Authoritative thorough Roborev job `904` stored exact range `557f8df3e6a62d34987edda70533375fc48ba8f6..4cdb53c27ed8165ee425e2bca36af316caff3892`, all four corrected guideline sections, and verdict FAIL.
- SKIPPED: Only after Roborev PASS independently reproduce non-destructive failed-resume fan-out, adjacent writer races, per-checkout durable outcomes, every AC, gofmt, focused/full/race suites, and detached audit evidence.
  Roborev failed the first gate, so no local validation or detached audit ran.
- SKIPPED: Only after every local gate passes push exact head 4cdb53c2 to PR #503 and engage required CI; report actual outcomes and block on waiting, skipped, unapproved, or failed evidence.
  No code was pushed and no CI was triggered or approved; PR #503 remains on remote head `3d50bd9a` while exact reviewed head `4cdb53c2` remains local-only.

### Reviewer findings

- HIGH: ready/init/commit trust any existing state directory without proving it is the registered worktree rooted exactly there on the expected state branch; Git can climb to the main worktree and mutate its code branch.
- HIGH: resume publishes the final checkout before origin convergence; non-locking state writers can mutate entities while add or pull/rebase is still operating, losing concurrent state.
- MEDIUM: lexical `.git` discovery misclassifies an external symlink into a linked worktree as standalone and resolves state under the linked worktree instead of the main worktree.
- MEDIUM: `strings.TrimSpace` corrupts legitimate leading or trailing whitespace in Git-reported paths and prefixes, redirecting workflow state.

### Summary

Validation is **REJECTED** at the Roborev-first gate. Exact-range job `904` found four unscoped repository-identity, publication-race, symlink, and path-byte defects; local tests, PR update, and CI engagement were not performed.

## Stage Report: implementation (cycle 6)

- DONE: Prove every existing state directory is the exact registered worktree on the expected state branch before ready/init/commit can mutate, preventing Git parent discovery into the code branch.
  Local commit `f57b0924` adds exact existing-checkout validation before mutation. Repository-backed state must match one canonical, non-bare, non-prunable `git worktree list --porcelain -z` record on the expected `refs/heads/<state branch>`, with an exact top-level and common-Git-dir match; standalone compatibility accepts only an exact repository root. A two-checkout-by-three-verb adversarial matrix covers a plain directory and a wrong-branch linked worktree across `state ready`, `state init`, and `state commit`, asserting the main HEAD and porcelain remain byte-for-byte unchanged.
- DONE: Converge origin privately before publishing the final checkout so non-locking state writers cannot lose concurrent entity changes, with adversarial writer coverage.
  Resume now creates an adjacent private worktree, fetches and rebases there, and only then atomically renames it to the final path, repairs metadata, validates exact identity, and records `ready`. A pre-publication origin test proves the private branch already equals the peer HEAD while the final path is absent. Deterministic ready/init races cover `status --set`, direct Git worktree creation with uncommitted bytes, and `state new`; publication fails closed and preserves every concurrent writer. Pull/rebase failure leaves no published checkout and preserves the stale prunable registration and local branch.
- DONE: Classify external symlinks into linked worktrees correctly and preserve exact whitespace path bytes; run focused/full/race gates and commit locally only for fresh Roborev-first validation.
  Repository placement now probes canonical Git metadata before lexical fallback, so a real external symlink into a linked worktree anchors to the primary worktree. Git path framing removes exactly one terminal LF instead of trimming path whitespace; unit and real-repository tests preserve spaces, tabs, and newline-bearing roots/prefixes. Changed files are gofmt-clean, `git show --check` and Windows cross-compilation pass, focused lifecycle/race suites pass, `go test ./... -count=1` passes, and `go test ./... -race -count=1` passes (CLI 325.677s, status 94.903s). The code branch is clean and four local commits ahead of unchanged remote `3d50bd9a`; no code push, PR update, CI trigger, or CI approval occurred.

### Summary

Closed all four job-904 findings with exact checkout identity checks, private origin convergence plus atomic publication, canonical linked-worktree classification, and byte-preserving Git path parsing. Exact clean head `f57b0924523ee75ef0a2c08f9c999d9c6b0756a1` is locally verified and intentionally unpushed for a fresh corrected exact-range Roborev-first validation.

## Stage Report: validation (cycle 6)

- FAILED: Run corrected-guideline Roborev first on exact range 557f8df3..f57b0924 and inspect the stored range and all four guideline sections; on any finding stop REJECTED without tests, PR update, or CI.
  Authoritative thorough Roborev job `951` stored exact range `557f8df3e6a62d34987edda70533375fc48ba8f6..f57b0924523ee75ef0a2c08f9c999d9c6b0756a1`, Compatibility posture, Trust boundaries, Behavioral proof, and Review focus, and verdict FAIL.
- SKIPPED: Only after Roborev PASS independently reproduce exact registered-worktree identity, private pre-publication origin convergence and writer races, external-linked-worktree symlink classification, whitespace-byte preservation, every AC, gofmt, focused/full/race suites, and detached audit evidence.
  Roborev failed the mandatory first gate, so no local validation suites or detached audit ran.
- SKIPPED: Only after every local gate passes push exact head f57b0924 to PR #503 and engage required CI; report actual outcomes and block on waiting, skipped, unapproved, or failed evidence.
  No code was pushed, PR #503 was not updated, and no CI was triggered or approved; exact reviewed head `f57b0924` remains local-only.

### Reviewer findings

- MEDIUM: `internal/cli/state.go:620` runs repository-wide `git worktree prune --expire now`, which can delete unrelated temporarily missing worktree registrations instead of repairing only `statePath`.
- MEDIUM: `internal/cli/state_sync.go:166` can observe the published checkout before the publisher releases the lock, then ignore the completed resume outcome and issue a redundant pull; an unreachable origin can turn successful local fallback into failure.
- MEDIUM: `internal/cli/state_resume_outcome.go:17` persists outcomes indefinitely keyed only by checkout path, so a deleted and independently recreated checkout can inherit a stale `ready` outcome without current remote integration.

### Summary

Validation is **REJECTED** at the Roborev-first gate. Exact-range job `951` found three remaining scope and generation-identity defects; local tests, PR publication, and CI engagement were correctly not performed.

## Stage Report: implementation (cycle 7)

- DONE: Replace repository-wide worktree pruning with state-path-scoped repair that cannot delete unrelated registrations, with an adversarial temporarily-missing sibling worktree proof.
  Commit `ff5bb036` replaces `git worktree prune --expire now` with Git's target-specific `git worktree remove <statePath>` while the final path is proven absent. `TestStateReadyRepairsOnlyTargetRegistration` makes both the state checkout and an unrelated sibling registration prunable, resumes the state checkout, and proves the sibling remains exactly registered and prunable.
- DONE: Make waiters consume the completed resume outcome after publication/lock release before any redundant pull, preserving successful local fallback when origin is unreachable.
  Ready/init callers now remember a generation-neutral `pending` outcome observed beside a newly published path and, after acquiring the repository lock, consume its generation-bound `ready` result before remote work. `TestStateReadyWaiterObservingPublishedCheckoutConsumesOutcome` makes origin unreachable after private convergence but before lock release and proves the waiter exits ready without a redundant pull; the later-peer control proves a completed outcome never suppresses a future ordinary pull.
- DONE: Bind durable resume outcomes to the current checkout generation so deletion and independent recreation cannot inherit stale ready; run focused/full/race gates and keep the branch local for fresh Roborev-first validation.
  Each `ready` outcome now carries a cryptographically random generation token stored in that linked worktree's private Git administration directory; reads require an exact token match, and outcome records publish atomically. The delete/recreate regression proves a same-path independent worktree cannot inherit stale ready. Focused and race-focused matrices pass; complete CLI passes (229.436s), `go test ./... -count=1` passes, `go test ./... -race -count=1` passes (CLI 311.640s, status 82.529s), broad gofmt ran with unrelated pre-existing drift restored, and Windows cross-compilation produced a PE32+ executable. The code branch is clean at `ff5bb0363a130a23228e6f5bf14679eaacb390e6`, five local commits ahead of unchanged remote `3d50bd9a`; no code push, PR update, CI trigger, or CI approval occurred.

### Summary

Closed all three job-951 findings with target-scoped registration repair, publication-aware waiter outcome consumption, and checkout-generation-bound durable outcomes. The exact clean local head is fully verified and intentionally unpushed for fresh corrected exact-range Roborev-first validation.

## Stage Report: validation (cycle 7)

- FAILED: Run corrected-guideline Roborev first on exact range 557f8df3..ff5bb036 and inspect the stored range and all four guideline sections; on any finding stop REJECTED without tests, PR update, or CI.
  Authoritative thorough Roborev job `1016` stored exact range `557f8df3e6a62d34987edda70533375fc48ba8f6..ff5bb0363a130a23228e6f5bf14679eaacb390e6`, Compatibility posture, Trust boundaries, Behavioral proof, and Review focus, and verdict FAIL.
- SKIPPED: Only after Roborev PASS independently reproduce target-only registration repair, pending-publication waiter outcome consumption, checkout-generation freshness across delete/recreate and later peers, every AC, gofmt, focused/full/race suites, Windows cross-compilation, and detached audit evidence.
  Roborev failed the mandatory first gate, so no local validation suites, cross-compilation, or detached audit ran.
- SKIPPED: Only after every local gate passes push exact head ff5bb036 to PR #503 and engage required CI; report actual outcomes and block on waiting, skipped, unapproved, or failed evidence.
  No code was pushed, PR #503 was not updated, and no CI was triggered or approved; exact reviewed head `ff5bb036` remains local-only.

### Reviewer findings

- MEDIUM: `internal/cli/state.go:624` checks destination absence before cleanup and rename, leaving a TOCTOU window that can remove or replace a checkout created after the check. Publication needs an atomic no-replace primitive and a deterministic interleaving test.
- MEDIUM: `internal/cli/state.go:591` treats any registration resolving to `statePath` as stale without proving expected branch, prunable status, and uniqueness; wrong-branch, active, or duplicate registrations can be removed.
- MEDIUM: `internal/cli/state_sync.go:166` recognizes a concurrent resume only when it observed `pending`; a caller seeing the published checkout after `ready` but before unlock can wait and then redundantly pull, turning a successful fallback into failure during an outage.

### Summary

Validation is **REJECTED** at the Roborev-first gate. Exact-range job `1016` found three remaining publication-order and registration-identity defects; local tests, PR publication, and CI engagement were correctly not performed.
