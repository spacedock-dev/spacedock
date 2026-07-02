---
title: "state ready/sweep/commit auto-discover the lone workflow (drop the --workflow-dir requirement)"
status: validation
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

## Stage Report: implementation

- DONE: One exported resolver in internal/status carries the walk-up + downward-fallback combo, with status dispatch (native_runner.go), MergeGuard, and the CLI's resolveWorkflowDir empty-flag branch as callers, and the inline copies in parseStateInitArgs/parseStateNewArgs removed — all five state verbs resolve identically (AC-4).
  `status.ResolveWorkflowDir` added at internal/status/discover_walkup.go:38; native_runner.go:75-79 and merge.go:73-79 now call it; cli/state_sync.go:417-424 `resolveWorkflowDir` calls `status.ResolveWorkflowDir` on empty workflowDir; cli/state.go's `parseStateInitArgs`/`parseStateNewArgs` now end with `return resolveWorkflowDir(workflowDir, dir, stderr)`, no inline copy. Commit c2635be4 on branch spacedock-ensign/state-subcommand-workflow-autodiscovery.
- DONE: New from-root CLI tests land per the task's Test plan: bare state ready/sweep/commit/init/new each exit 0 with the operation observably performed on the single-nested-workflow fixture, and the two-workflow fixture gets the candidate-listing refusal (AC-1, AC-2).
  New file internal/cli/state_from_root_test.go (6 tests): TestStateReadyFromRootResolves, TestStateSweepFromRootResolves, TestStateCommitFromRootResolves, TestStateInitFromRootResolves, TestStateNewFromRootResolves, TestStateVerbsFromRootRefuseAmbiguousWorkflows (drives ready/sweep/commit against a two-workflow fixture). All pass (`go test ./internal/cli/... -run FromRoot`: 6/6 passed).
- DONE: The existing state-verb and discovery suites pass unmodified (state_ready/sweep/commit/init/new/sync tests, native_new_from_root, merge-guard) — the explicit --workflow-dir path and existing callers do not regress (AC-3).
  Full repo test suite green: `go test ./...` all packages ok, including internal/cli and internal/status.

### Summary

Extracted the shared walk-up + downward-fallback discovery block (previously duplicated in native_runner.go's dispatch and merge.go's MergeGuard) into one exported `status.ResolveWorkflowDir`, and routed the CLI's `resolveWorkflowDir` plus `parseStateInitArgs`/`parseStateNewArgs` through it, removing their inline walk-up-only copies. All five state verbs now auto-discover the lone commissioned workflow with the same zero/ambiguity diagnostics as `spacedock new`. Live-verified: bare `spacedock state sweep` from this repo's toplevel now exits 0 (previously exited 1 requiring `--workflow-dir`). Added internal/cli/state_from_root_test.go covering the from-root success and ambiguity-refusal paths; all existing suites pass unmodified.

## Stage Report: validation

- DONE: AC-1's end-value is reproduced against the failing baseline, not trusted: run the new from-root CLI tests yourself AND exercise the built binary from the repo toplevel of a single-workflow fixture — bare state ready/sweep/commit each succeed in ONE invocation (exit 0, effect observed), where the pre-change behavior was fail-then-retry; reject self-referential evidence.
  Ran `go test ./internal/cli/... -run FromRoot -count=1` fresh: 6/6 pass. Built spacedock-baseline from `git archive 59251f18` and spacedock-fixed from worktree HEAD c2635be4; on a scratch single-workflow split-root fixture (plain toplevel, workflow at docs/dev), baseline bare ready/sweep/commit each exit 1 (`no workflow here — pass --workflow-dir`) while `status --discover` resolves from the same cwd; fixed binary in one invocation each: ready exit 0 and integrated a peer clone's pushed entity (peer-task.md appeared in the checkout), sweep exit 0 emitting the merged-PR listing (`0 entity(ies) merged but not yet terminalized.`), commit exit 0 landing path-scoped (git show name-only = seed-task.md only) and pushed to origin.
- DONE: AC-2 and AC-4 verified: the two-workflow fixture refusal names every candidate and instructs --workflow-dir (parity with spacedock new), and the resolver dedup is real — one shared walk-up+downward implementation in internal/status with status dispatch, merge guard, and all five state verbs as callers, no inline copies left in state.go.
  Two-workflow scratch fixture: bare ready/sweep/commit each exit 1 with `Error: multiple commissioned Spacedock workflows found under …; pass --workflow-dir <dir> to pick one:` listing docs/dev and docs/ops — byte-identical to `spacedock new`'s refusal from the same fixture. Grep over non-test product code: `discoverWorkflowDownward` is called only from `ResolveWorkflowDir` (discover_walkup.go:45); callers are native_runner.go:76 (status dispatch), merge.go:75 (merge guard), state_sync.go:420 (CLI resolveWorkflowDir, which parseStateInitArgs/parseStateNewArgs/parseStateCommitArgs/parseWorkflowJSONArgs all route through); the diff removes both inline copies from state.go. frontdoor.go:213 calls only the walk-up for banner labeling, not verb resolution.
- DONE: AC-3 plus the full applicable suites are green (existing state-verb, native_new_from_root, merge-guard tests unmodified; go test over the touched packages and the Testing Resources the stage names); close with a PASSED/REJECTED recommendation citing evidence per AC.
  `git diff main --name-only -- '*_test.go'` shows only the new state_from_root_test.go — existing suites unmodified. Fresh `go test ./... -count=1`: all 16 packages ok. `go test -race -count=1 ./internal/cli/... ./internal/status/...`: ok. Named suites: 22 TestState{Ready,Sweep,Commit,Init,New,Sync}* pass; full merge-guard suite and TestNewFromRoot{SingleWorkflowSucceeds,MultiWorkflowRefuses} pass. Testing Resources smoke: --help exit 0, --version exit 0, `status --workflow-dir docs/dev --validate` from the main checkout prints VALID exit 0 (it exits 1 inside the linked worktree for baseline and fixed alike — the state checkout only exists in the main checkout, environmental, not a regression). Zero-workflow diagnostic is the shared named error, exit 1.

### Summary

PASSED. Every AC reproduced independently: the 59251f18 baseline binary fails all three bare verbs from a fixture toplevel where the fixed binary completes each in one invocation with the effect observed on disk; the two-workflow refusal matches `spacedock new` byte-for-byte; the resolver exists once in internal/status with all six call paths routed through it; and the full suite, race-enabled touched packages, and all named existing suites pass with zero existing test files modified. Detached adversarial audit not run — this is a routine, low-blast-radius CLI resolver dedup with behavioral tests, not a contract/skill/gate surface.

Recommendation: PASSED.

### Feedback Cycles

**Cycle 1 (2026-07-02, detached adversarial audit, pre-gate).** The validator skipped the audit as low-blast-radius; the FO overrode (workflow-dir resolution decides where state mutations land — the proof policy's "status mutation and guard paths" surface). Audit `audit-state-workflow-autodiscovery` on a throwaway clone: shipped code behaved correctly under all six mutants; two material TEST-STRENGTH holes routed to implementation: (M1) explicit `--workflow-dir` precedence unpinned on the state-verb path — a mutant in `resolveWorkflowDir` (internal/cli/state_sync.go:418) letting a successful discovery silently override the explicit flag stays GREEN suite-wide; fix: two-workflow fixture, run a state verb from inside/above workflow A with `--workflow-dir` naming B, assert the mutation lands in B. (M2) `TestStateSweepFromRootResolves` (internal/cli/state_from_root_test.go:42) passes even when single-match resolution returns the git toplevel — sweep on a wrong-but-existing dir exits 0 with an empty valid envelope; fix: assert the JSON `state_branch` field (populated only on real workflow resolution). Four other mutants (ambiguity fail-open, wrong-dir on the other four verbs, walk-up removal, commission bypass, zero-match fail-open) were killed by existing tests. Secondary notes recorded: CLI state-verb suite has no walk-up or zero-workflow case of its own and the ambiguity loop skips init/new — acceptable while all five verbs share one resolver.
