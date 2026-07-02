---
title: "state ready/sweep/commit auto-discover the lone workflow (drop the --workflow-dir requirement)"
status: implementation
source: "FO boot-ergonomics friction 2026-07-01 (Claude Commander session hit it twice; Pi FO Stage B report shows the '(x2 with --workflow-dir)' annotation). `spacedock state ready`/`sweep`/`commit` fail with 'no workflow here — pass --workflow-dir' when run from the repo root, while `status --boot`/`status --discover`/`spacedock new` already auto-discover the lone commissioned workflow. The inconsistency costs a failed-then-retried call at boot (doubling Stage-B state-read tokens) and is a recurring papercut."
group: tooling
id: 82edd88rq11q2f05z5nhfhj8
started: 2026-07-02T01:23:59Z
worktree: .worktrees/spacedock-ensign-state-subcommand-workflow-autodiscovery
---

## Problem
`spacedock state ready`, `state sweep`, and `state commit <slug>` require an explicit `--workflow-dir` when run from the project root — they do NOT auto-discover the single commissioned workflow the way `status --boot`, `status --discover`, and `spacedock new` already do. At FO boot this surfaces as a `no workflow here — pass --workflow-dir` error on the first attempt, then a retry with the flag: wasted round-trips and boot-context tokens (the Pi FO Stage B report's `(x2 with --workflow-dir)` annotation is this friction).

## Root cause
The state verbs already auto-discover — but with only HALF the shared discovery. When `--workflow-dir` is empty they call `status.DiscoverWorkflowDir` (walk-UP to the enclosing commissioned README, `internal/status/discover_walkup.go:16`) and error on a miss:
- `resolveWorkflowDir` (`internal/cli/state_sync.go:417`), used by `ready`/`sweep`/`commit`;
- the same logic duplicated inline in `parseStateInitArgs` (`internal/cli/state.go:96`) and `parseStateNewArgs` (`internal/cli/state.go:356`).

The path `status --boot`/`--discover`/`new` share (`internal/status/native_runner.go:75-83`) — also used by `merge guard` (`internal/status/merge.go:73-82`) — tries the walk-up FIRST, then falls back to `discoverWorkflowDownward` (`internal/status/native_runner.go:565`): scan downward from the git toplevel via `discoverWorkflows`; exactly 1 match resolves, 0 is a named error, >1 refuses listing the candidate dirs with a `--workflow-dir` instruction. From a repo root whose toplevel README is not commissioned, the walk-up misses and only the downward fallback finds the nested workflow — so the state verbs fail where `status`/`new` succeed.

Reproduced at HEAD 59251f18: from this repo's root, bare `spacedock state sweep` exits 1 with `no workflow here — pass --workflow-dir`, while `spacedock status --discover` from the same cwd resolves `docs/dev` and bare `spacedock status` renders the queue (exit 0).

**Gap audit:** all five state verbs share the gap — `init`, `new`, `ready`, `sweep`, and `commit` (all resolve via walk-up only). `merge guard` already carries the walk-up + downward combo and is not affected.

## Approach
Reuse the existing shared discovery — no parallel implementation:

1. Extract the walk-up + downward-fallback combo (the exact block at `native_runner.go:75-83` / `merge.go:73-82`) into one exported resolver in `internal/status` (walk up via `DiscoverWorkflowDir`; on miss, `discoverWorkflowDownward` with its existing diagnostics).
2. Call it from the three current sites: `dispatch` (native_runner), `MergeGuard`, and the CLI's `resolveWorkflowDir` empty-flag branch — deduplicating the combo that today appears twice in `internal/status`.
3. Route `parseStateInitArgs`/`parseStateNewArgs` through `resolveWorkflowDir`, removing their duplicated inline copies, so all five state verbs resolve identically.

Behavior notes:
- The zero-workflow and ambiguity diagnostics for bare state verbs become the shared ones (`Error: no commissioned Spacedock workflow found in …` / `multiple commissioned Spacedock workflows found under …; pass --workflow-dir <dir> to pick one:` + candidate list) — exact parity with `spacedock new`. No test pins the old `no workflow here` message (grep over `*_test.go` and `fixtures/` is empty).
- An explicit `--workflow-dir` (or `PIPELINE_DIR` on the status path) still skips discovery entirely — unchanged.

## Acceptance criteria
- **AC-1 (Measures the end-value)** — From the repo toplevel of a single-workflow project (toplevel README not commissioned, lone workflow nested), bare `spacedock state ready`, `state sweep`, and `state commit <slug>` each complete in **one invocation**: exit 0 with the operation observably performed (ready: peer state integrated / absent checkout resumed; sweep: the merged-PR listing emitted; commit: the path-scoped commit present in the state checkout's `git log`). Independent baseline that can move the wrong way: at 59251f18 the same invocations exit 1 — two calls per verb (fail, retry with flag). This is the boot retry-double-call the entity exists to eliminate: 1 call, not 2, per boot verb. *Test:* new real-git CLI tests driving `run()` with hostDir = repo toplevel and no `--workflow-dir`, asserting exit 0 plus the on-disk/output effect per verb.
- **AC-2** — With two commissioned workflows under the toplevel, the bare verbs refuse (non-zero exit) and stderr names every candidate dir and instructs `--workflow-dir` — the same refusal `spacedock new` emits. *Test:* two-workflow fixture per the assertions in `TestNewFromRootMultiWorkflowRefuses` (`internal/status/native_new_from_root_test.go:45`).
- **AC-3** — The explicit `--workflow-dir` path does not regress: the existing state-verb suites (`state_ready_test.go`, `state_sweep_test.go`, `state_commit_test.go`, `state_init_test.go`, `state_new_test.go`, `state_sync_test.go` — every invocation passes `--workflow-dir`) pass unmodified. *Test:* the unmodified suites themselves.
- **AC-4** — Workflow-dir resolution flows through one shared resolver: the walk-up + downward combo exists once in `internal/status`, with `status` dispatch, `merge guard`, and all five `state` verbs as callers; `state init`/`new` no longer carry inline copies. *Test:* existing `native_new_from_root` + merge-guard suites stay green after the refactor (behavioral proof the shared path still serves its original callers); reviewer confirms no second implementation in the diff.

## Test plan
- **New CLI tests** (one file, e.g. `internal/cli/state_from_root_test.go`), real-git fixtures shaped like `native_new_from_root_test.go` (plain toplevel README + nested commissioned workflow) combined with the split-root helpers already in the cli suite (`twoHostStateWorkflow`, `writeEntity`):
  - from-root `state ready` on a split-root fixture: exit 0, checkout present/integrated (AC-1);
  - from-root `state sweep` with zero pr-pending entities: exit 0, offline — per `TestStateSweepIsReadOnly` (AC-1);
  - from-root `state commit <slug>`: exit 0, commit landed path-scoped in the state checkout log (AC-1);
  - from-root `state init` (absent checkout) and `state new` (unbirthed README): resolve and act, exit 0 (AC-1 extension to the audited verbs);
  - two-workflow ambiguity refusal for the bare verbs (AC-2).
- **Existing suites unmodified** — AC-3/AC-4 evidence.
- Estimated cost: small. One exported resolver + four call-site edits + two dedups in product code; ~1 new test file (~150–250 lines). Command-level behavior is the claim, so CLI tests suffice — no live workflow smoke test needed.

## Doc impact
Assessed — **no doc diff required.** No docs-site page or skill text names a `--workflow-dir` requirement for the state verbs; they already prescribe the bare invocations this task makes work from the repo root:
- FO skill: `skills/first-officer/references/first-officer-shared-core.md:121` (`spacedock state ready`), `:129` (`spacedock state sweep`), `:136` (`spacedock state commit <slug>`) — all bare.
- Docs site: `docs/site/advanced/split-root-state.md:5` (`run spacedock state init`, bare) and `docs/site/reference/command-reference.md:49` (no flag mention).
- Help text: the `state` command's Use line already marks the flag optional (`[--workflow-dir DIR]`, `internal/cli/cli.go`); still accurate after the change.

## Spike determination
No spike needed: the design relies only on proven mechanisms. (a) The gap was reproduced live at HEAD 59251f18 (bare `state sweep` from repo root exits 1; `status --discover` and bare `status` from the same cwd resolve `docs/dev`). (b) The downward fallback being reused is shipping code already exercised in production by `status`/`new`/`merge guard` and pinned by `TestNewFromRootSingleWorkflowSucceeds`/`TestNewFromRootMultiWorkflowRefuses`. Nothing unverified remains on the riskiest path.

## Related
- Sibling boot-ergonomics item: `status --read --json --frontmatter` boot-lean projection.
- Discovery precedent to reuse: `status --boot`/`status --discover`/`spacedock new` auto-discovery.

## Stage Report: ideation

- DONE: The rough sketch is tightened into measured ACs with test plans: single-workflow auto-discovery succeeds with no flag (exit 0, operation performed — the boot retry-double-call vanishes), >1 workflow reports candidates and requires --workflow-dir (parity with spacedock new), and the explicit --workflow-dir path does not regress.
  AC-1 measures 1-call-vs-2 against the 59251f18 baseline; AC-2 mirrors TestNewFromRootMultiWorkflowRefuses; AC-3 is the unmodified existing suites; each AC names its test.
- DONE: The approach names the existing discovery code path (the one status --boot/--discover/new share) and reuses it for state ready/sweep/commit — no parallel discovery implementation — and audits whether other state subcommands share the gap.
  Path named: walk-up (discover_walkup.go:16) + downward fallback (native_runner.go:565), combined at native_runner.go:75-83 and merge.go:73-82; plan extracts one shared resolver. Audit: all five state verbs share the gap (init/new via inline copies at state.go:96/356, ready/sweep/commit via state_sync.go:417); merge guard already has the fallback.
- DONE: Doc impact is assessed: if the docs site or skill text names the --workflow-dir requirement for state subcommands, the concrete before/after doc diff is recorded in the task body per the ideation stage definition.
  Audited FO skill, docs site, and help text — none names the requirement; all already prescribe the bare invocations, so no doc diff is required (files and line numbers recorded in the Doc impact section).

### Summary

Root-caused the gap: state verbs resolve an empty --workflow-dir with only the walk-up half of discovery, while status/new/merge guard add the downward-scan fallback; reproduced live at 59251f18 (bare `state sweep` from repo root exits 1 while `status --discover` resolves docs/dev from the same cwd). The approach extracts the existing walk-up+downward combo into one exported resolver reused by all callers, fixing all five state verbs and deduplicating two inline copies. No spike needed (mechanisms shipping and test-pinned) and no doc diff needed (docs already show bare invocations).
