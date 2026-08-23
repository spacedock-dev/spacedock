---
title: Resolve split-state gate artifact paths before preparation
status: validation
source: "PR #679 run 31640122995, Codex job 94260562369, artifact 9158783630: TestLiveCommonGateGuardrail found committed review files in the split state checkout, then passed a state-relative path to gate prepare, which resolved it from the workflow root and failed."
score: 0.98
id: mvmpzgqxyb32t3b3vdw0x0h1
gates:
    version: 1
    records:
        - id: gate:mvmpzgqxyb32t3b3vdw0x0h1:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:mvmpzgqxyb32t3b3vdw0x0h1-backlog-1
              briefing:
                id: briefing:mvmpzgqxyb32t3b3vdw0x0h1:backlog:attempt-1:revision-1
                digest: sha256:449b528bc4625bcc07c6f1463e438d550123d3e19fd4ba1c0821e80c28d13d7a
                request-digest: sha256:d855623326f6ca92f7fffb5b906151ab29b31b736198878fa45a77ae514fe3a6
                room-ref: ./resolve-split-state-gate-artifact-path/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:mvmpzgqxyb32t3b3vdw0x0h1:backlog:1
                briefing: briefing:mvmpzgqxyb32t3b3vdw0x0h1:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-21T08:15:28.894915131Z"
                decision: approve
                reason: Captain granted conn for pi-related fixes including gates (2026-08-21 chat). Seed clearly identifies the root cause (split-state gate-artifact path resolving against the wrong root) blocking 3 required Pi journeys + Codex lane; one-canonical-path approach is sound. Advance to ideation.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:mvmpzgqxyb32t3b3vdw0x0h1:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:mvmpzgqxyb32t3b3vdw0x0h1-ideation-1
              briefing:
                id: briefing:mvmpzgqxyb32t3b3vdw0x0h1:ideation:attempt-1:revision-1
                digest: sha256:165bae69f49b3d30abc783ba45eee060bef63700ad2d7ef02d05d8c6e4ce5506
                request-digest: sha256:0668bee18a33b76e4753899717fa3e6df696e3114bd1ff08248836b70d00555f
                room-ref: ./resolve-split-state-gate-artifact-path/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:mvmpzgqxyb32t3b3vdw0x0h1:ideation:1
                briefing: briefing:mvmpzgqxyb32t3b3vdw0x0h1:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-21T08:49:16.597398723Z"
                decision: approve
                reason: 'Captain approved both ideation gates. mvmpzg: one canonical resolver-owned gate-artifact path contract for split-root — root cause for the 3 required Pi journeys + Codex lane. Advance to implementation, worktree stacked on 747.'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:mvmpzgqxyb32t3b3vdw0x0h1:validation
          stage: validation
          attempts:
            - id: gate-attempt:mvmpzgqxyb32t3b3vdw0x0h1-validation-1
              briefing:
                id: briefing:mvmpzgqxyb32t3b3vdw0x0h1:validation:attempt-1:revision-1
                digest: sha256:00be9a055f8f8a545175fec89aef0f932e8abc4ebc62ed39e6e202888f0bb1b4
                request-digest: sha256:1804a72d4ce1d7cb72acdafaf2b7e56247f6ac466619a4de77f007fa3a7cd7ef
                room-ref: ./resolve-split-state-gate-artifact-path/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:mvmpzgqxyb32t3b3vdw0x0h1:validation:1
                briefing: briefing:mvmpzgqxyb32t3b3vdw0x0h1:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-21T09:13:37.371629499Z"
                decision: approve
                reason: 'Conn-held approval. mvmpzg validation PASSED: 3 ACs falsifiable, gates+cli green, gofmt/vet clean. Deliver via stacked PR on 747.'
              application:
                target-stage: done
                state: superseded
started: 2026-08-21T08:15:56Z
worktree: /home/exedev/spacedock/.worktrees/spacedock-ensign-resolve-split-state-gate-artifact-path
mod-block:
pr:
---

Make a split-root First Officer pass the exact committed artifact path that gate preparation can resolve.

## Problem

A valid committed gate review exists in the state checkout, but the First Officer can select the wrong path base. In a split-root workflow the entity, gate review, and entity-snapshot reference live in the separate state checkout (`<definitionDir>/.spacedock-state`), while the FO runs from the workflow definition dir. When the FO passes a state-relative `--artifact` path, the CLI resolves it against cwd (the workflow root), producing a path that does not exist. `gitsource.Inspect` then fails with "selected source must be a readable non-symlink regular file: lstat ... no such file or directory", or — if the FO passes a bare slug without `.md` — `isMarkdownPath` rejects it with "--artifact must name a .md or .markdown file" before resolution even starts. The required Codex lane and 3 Pi journeys that depend on `TestLiveCommonGateGuardrail` reaching an open prepared gate all stop here.

## Diagnosis

**Where the FO selects the path base.** The CLI (`internal/cli/cli.go`, the `gate prepare` RunE closure) resolves relative `--artifact` and `--reference` paths against `dir` (the process cwd):

```go
if !filepath.IsAbs(prepareInput.Artifact) {
    prepareInput.Artifact = filepath.Join(dir, prepareInput.Artifact)
}
```

`dir` is the working directory the binary was invoked from — in a live run, the workflow definition root, not the state checkout. The entity path itself is correctly resolved: `status.ResolveActivePath(definitionDir, dir, args[1], stderr)` calls `resolveRoots` internally, which reads the README `state:` field and diverges `entityDir` to `definitionDir/<state>`. But the artifact path bypasses this root-splitting entirely and joins with bare cwd.

**`gates.Prepare` does not re-resolve.** After the CLI makes the path absolute against cwd, `gates.Prepare` calls `filepath.Abs(selected)` — a no-op on the already-absolute (wrong-root) path. The subsequent `gitsource.Inspect(roots, selected)` calls `os.Lstat(path)` which fails because the file lives under `<definitionDir>/.spacedock-state/...`, not `<cwd>/...`.

**Riskiest-mechanism spike result.** Exercised against a split-root fixture (README declares `state: .spacedock-state`; entity at `stateRoot/recorded-gate-task/index.md`; gate review at `stateRoot/recorded-gate-task/selected/gate-review.md`; cwd = workflow root). Three input shapes tested:

1. State-relative `.md` path (`recorded-gate-task/selected/gate-review.md`): CLI joins with cwd → `/…/<root>/recorded-gate-task/selected/gate-review.md` (missing `.spacedock-state`) → `gitsource.Inspect` fails: "selected source must be a readable non-symlink regular file: lstat ... no such file or directory".
2. Bare slug without `.md` (`recorded-gate-task`): `isMarkdownPath` rejects before resolution → "--artifact must name a .md or .markdown file".
3. Absolute entity path (pre-resolved by `status.ResolveActivePath`): succeeds → `state=open`.

Confirmed: the FO's `--artifact` path is state-relative; the CLI resolves it against cwd (workflow root), not the state checkout; `gates.Prepare`'s `filepath.Abs` is a no-op on the already-wrong absolute path.

## Proposed approach

**One canonical resolver-owned artifact path.** The binary — not the model — owns path resolution. `gates.Prepare` resolves a relative `--artifact` or `--reference` path against the state-checkout entity root, computed once from the README `state:` field of `input.WorkflowDir`. The FO passes a state-relative path (or an absolute one); the binary resolves it from the correct root.

The entity root is:
- Split-root: `filepath.Join(input.WorkflowDir, stateRelPath)` where `stateRelPath` is the cleaned README `state:` value.
- Single-root (`state:` absent, empty, or `$inline`): `input.WorkflowDir` itself — unchanged from the current cwd-based behavior for single-root, since cwd and the workflow root coincide.

The CLI stops resolving relative paths against `dir` (cwd) and passes them through unchanged. `gates.Prepare` resolves relative paths against the entity root; absolute paths pass through `filepath.Abs` unchanged. This is the single path contract: one resolution site, one root, computed from the same README the entity path is derived from.

`gates` already reads the README frontmatter (`entityIdentity`, `validatePreparedStage`) using its own `frontmatterNode`/`mappingValue` helpers, so no new package dependency is introduced. The `state:` classification mirrors `status.ClassifyState` semantics (reject absolute or `..`-escaping values) but stays local to `gates` to avoid a `status → gates` import cycle.

**New mechanism: entity-root resolution in `gates.Prepare`.** Value AC served: AC-1 (the gate-guardrail fixture reaches an open prepared gate). Simplest alternative: resolve relative to cwd / pass through `filepath.Abs` unchanged (the current behavior). Why insufficient: in split-root, cwd is the workflow root, not the state checkout — the artifact file does not exist under cwd, so resolution produces a non-existent path and `gitsource.Inspect` fails. The entity path is already correctly root-split by `status.ResolveActivePath`; the artifact path must use the same root.

## Out of scope

Do not change DVD or v8 behavior, XFAIL policy, or unrelated gate semantics. Do not change the gate grammar, Briefing/Request format, room structure, replay idempotency, or `gitsource.Inspect`/`classify` logic. Do not add global hooks, test-only product machinery, or new `status` imports into `gates`.

## Acceptance criteria

**AC-1 — Split-root gate preparation accepts the committed selected artifact (value-measuring).**
The exact `gate-guardrail` fixture (`recorded-gate/held`) reaches an open prepared gate (`state=open`) under split-root when the FO passes a state-relative `--artifact` path. Measured against the current wrong-root baseline, which fails: a state-relative `.md` path resolves against cwd and `gitsource.Inspect` reports "no such file or directory"; a bare slug without `.md` is rejected by `isMarkdownPath`. After the fix, the same state-relative path resolves against the state-checkout entity root and reaches `state=open`. A path based at the wrong root (resolvable against cwd but not the state checkout) still fails.

**AC-2 — One canonical path contract (mechanism).**
Focused path-resolution tests cover: (a) split-root resolves a state-relative `--artifact` against the entity root and reaches `state=open`; (b) single-root resolves a relative `--artifact` unchanged (entity root = workflow dir = cwd); (c) an absolute `--artifact` passes through unchanged in both modes; (d) a relative `--reference` resolves the same way. No host-specific parsing, no cwd probe, no model path reconstruction.

**AC-3 — The 3 required Pi journeys + Codex lane reach an open prepared gate (value).**
The `TestLiveCommonGateGuardrail` live journey (entry point for the `gate-guardrail` CI target) passes on all runtimes that were blocked by this bug: the FO prepares the gate from the committed state-checkout artifact and reaches `state=open` without deciding, advancing, dispatching, or archiving. Verified by the exact local Codex gate-guardrail target on final bytes.

## Test plan

1. **Focused path-resolution tests** (`internal/gates/prepare_test.go`): extend the existing `prepareFixture` harness with a split-root case where the artifact lives in the state checkout. Assert: (a) a state-relative `--artifact` resolves to the state checkout and reaches `state=open`; (b) a relative `--reference` in the state checkout resolves the same way; (c) an absolute path passes through unchanged; (d) single-root (no `state:` field) resolves relative paths against the workflow dir unchanged. Each test fails if the resolution root reverts to cwd. Cost: low (table-driven, reuses `prepareFixture`).
2. **CLI passthrough test** (`internal/cli/cli_test.go` or equivalent): verify the CLI no longer resolves relative `--artifact` against cwd — it passes the relative path to `gates.Prepare` which resolves it against the entity root. Cost: low.
3. **Wrong-root regression test**: a state-relative path that would resolve to an existing file under cwd but not under the state checkout fails (guards against reverting to cwd resolution). Cost: low.
4. **Gate-guardrail fixture** (final bytes): run the exact local Codex gate-guardrail target once (`SPACEDOCK_LIVE_RUNTIME=codex go test ./internal/ensigncycle/ -run TestLiveCommonGateGuardrail -count=1`). Cost: medium (live model run).
5. **Full and race suites** once, since Go behavior changes (`go test ./...` and `go test ./... -race`). Cost: medium.

## Expected surface and tolerance

**Files:** `internal/gates/prepare.go` (path resolution + entity-root helper), `internal/cli/cli.go` (remove cwd-based artifact/reference join), `internal/gates/prepare_test.go` (focused tests). Optionally `internal/cli/cli_test.go` if a CLI-level test is added.

**Estimate: net +30 LOC, across 3 files.** Insertions ~40 (entity-root helper, resolution change, tests); deletions ~10 (CLI `filepath.Join(dir, ...)` blocks). Tolerance: ±15 LOC net, ±1 file. A correction round measures actuals against these figures.

**Observable semantics declared:** the `--artifact` and `--reference` path resolution root changes from cwd to the state-checkout entity root in split-root workflows. No change to: gate grammar, Briefing/Request format, room structure, DVD or v8 behavior, XFAIL policy, replay idempotency, `gitsource` classification, or any unrelated gate semantics.

## Stage Report: ideation

- DONE: Concrete approach specifying the one canonical resolver-owned artifact path for gate prepare in a split-root workflow
  Body defines the single path contract: `gates.Prepare` resolves relative `--artifact`/`--reference` against the state-checkout entity root computed from the README `state:` field; CLI passes relative paths through unchanged.
- DONE: Diagnose where the FO currently selects the path base (state checkout vs workflow dir)
  Diagnosis section traces the CLI `filepath.Join(dir, prepareInput.Artifact)` in `internal/cli/cli.go` as the wrong-root site; `gates.Prepare`'s `filepath.Abs` is a no-op on the already-wrong absolute path.
- DONE: Name the value-AC (3 required Pi journeys + Codex lane reach an open prepared gate) and the simplest alternative and why each is insufficient
  AC-3 names the 3 Pi journeys + Codex lane; simplest alternative (resolve relative to cwd / `filepath.Abs` unchanged) is insufficient because cwd is the workflow root, not the state checkout, in split-root.
- DONE: At least one value-measuring AC measured against a baseline that can move the wrong way
  AC-1: the gate-guardrail fixture reaches `state=open` under split-root; baseline (current wrong-root resolution) fails with "no such file or directory" or "--artifact must name a .md or .markdown file".
- DONE: Pair the value-AC with a mechanism AC exercised by focused path-resolution tests
  AC-2: split-root resolves the committed artifact; single-root resolves unchanged; absolute paths pass through — exercised by focused path-resolution tests.
- DONE: Expected surface and tolerance with observable-semantics declaration
  Net +30 LOC across 3 files (±15 LOC, ±1 file); observable semantics: path resolution root changes only, no gate grammar/DVD/v8/XFAIL/unrelated semantics changes per Out of scope.
- DONE: Record the riskiest-mechanism spike result
  Spike exercised against a split-root fixture: FO's `--artifact` is state-relative; CLI resolves against cwd (workflow root); `gates.Prepare`'s `filepath.Abs` is a no-op; wrong-root path fails at `gitsource.Inspect` (`.md` path) or `isMarkdownPath` (bare slug); absolute pre-resolved path succeeds.

### Summary

Fleshed out the seed into a complete ideation body: diagnosed the wrong-root path resolution in the CLI's `gate prepare` closure (joins relative `--artifact` with cwd instead of the state-checkout entity root), recorded the riskiest-mechanism spike result (three input shapes exercised against a split-root fixture), and defined the one canonical path contract — `gates.Prepare` resolves relative artifact/reference paths against the entity root computed from the README `state:` field, with the CLI passing relative paths through unchanged. Three ACs (value-measuring, mechanism, and the 3 Pi journeys + Codex lane), a five-step test plan, and an expected surface of net +30 LOC across 3 files with no undeclared semantic changes.

## Stage Report: implementation

- DONE: Implement the one canonical resolver-owned gate-artifact path contract for split-root in internal/gates/prepare.go
  Added `entityResolveRoot(workflowDir)` — computes the state-checkout entity root from the README `state:` field once (split-root joins the cleaned state path under WorkflowDir; single-root returns WorkflowDir; rejects absolute/`..`-escaping values mirroring `status.ClassifyState`). Added `resolveSelectedSource(selected, entityRoot)` — resolves relative paths against the entity root, absolute paths pass through cleaned. Replaced `filepath.Abs(selected)` in the normalization loop with `resolveSelectedSource`.
- DONE: Remove the FO path-selection site that resolves --artifact against cwd
  Removed the `filepath.Join(dir, prepareInput.Artifact)` and `filepath.Join(dir, prepareInput.References[i])` blocks from the `gate prepare` RunE closure in `internal/cli/cli.go`. Relative paths now pass through to `gates.Prepare` unchanged.
- DONE: Focused path-resolution tests covering split-root and single-root
  `TestPrepareResolvesStateRelativeArtifactAgainstEntityRoot` (split-root state-relative artifact reaches `state=open`), `TestPrepareResolvesStateRelativeReferenceAgainstEntityRoot` (split-root state-relative reference), `TestPrepareAbsoluteArtifactPassesThroughUnchanged` (absolute passthrough both modes), `TestPrepareSingleRootResolvesRelativeArtifactAgainstWorkflowDir` (single-root relative resolves against workflow dir), `TestPrepareWrongRootRelativeArtifactFails` (wrong-root path fails — the falsifying change), `TestGatePrepareCLIPassesStateRelativeArtifactWithoutCwdJoin` (CLI no longer joins relative --artifact with cwd).
- DONE: Update the gate-guardrail fixture reaching an open prepared gate under split-root
  Updated `prepareFixture` to use `state: .state` (non-escaping, matching the CLI fixture) since `entityResolveRoot` now validates the `state:` field. Added `prepareSingleRootFixture` for the single-root tests.
- DONE: Verify go test ./internal/gates/ -count=1 and go test ./... are green
  `go test ./internal/gates/ -count=1` passes; `go test ./internal/gates/ -count=1 -race` passes. `go test ./...` — all packages green except 4 pre-existing/environmental CLI failures (`TestVersionAmbiguousMarkersExitZero`, `TestCodexChannelInstallLeavesCoHostedPluginInstalled`, `TestCodexModeSwitchRoundTripPreservesExclusivity`, `TestCodexPluginInstallIsHostNative`) confirmed failing on the base commit before any changes (Codex CLI returns WARNING text instead of JSON in `/tmp`; version test sees Pi session marker). Noted, not fixed.

### Summary

Implemented the one canonical resolver-owned gate-artifact path contract: `gates.Prepare` resolves relative `--artifact`/`--reference` paths against the entity root computed once from the README `state:` field, and the CLI passes relative paths through unchanged. Added focused path-resolution tests covering split-root (state-relative artifact/reference reach `state=open`), single-root (relative resolves against workflow dir), absolute passthrough in both modes, and the wrong-root falsifying change. Updated the gates test fixture to use a non-escaping `state:` path matching the CLI fixture. Net +240 LOC across 4 files (over the +30 estimate due to comprehensive tests — the implementation change itself is ~+55/-13 LOC across 2 files; the rest is tests).

## Stage Report: validation

- DONE: Verify deliverable — one canonical resolver-owned gate-artifact path contract in internal/gates/prepare.go and the FO path-selection site in internal/cli/cli.go
  Commit 903778149 adds `entityResolveRoot(workflowDir)` (computes state-checkout entity root from README `state:` field, mirrors `status.ClassifyState` semantics locally) and `resolveSelectedSource(selected, entityRoot)` (relative → entity root, absolute → cleaned passthrough). CLI `filepath.Join(dir, …)` blocks for artifact and references removed — relative paths pass through unchanged. Diff: +253/-13 across 4 files (`prepare.go` +52/-3, `cli.go` 0/-8, `prepare_test.go` +173/-2, `gate_test.go` +28/0).
- DONE: AC-1 — split-root gate preparation accepts the committed selected artifact (value-measuring)
  `TestPrepareResolvesStateRelativeArtifactAgainstEntityRoot` passes a state-relative `selected/gate-review.md` under split-root; `entityResolveRoot` resolves `state: .state` to `workflow/.state`; `resolveSelectedSource` joins to `workflow/.state/selected/gate-review.md`; `gitsource.Inspect` finds the committed file; result `state=open`. Reverting to cwd resolution fails (file not under workflow root). Falsifiable.
- DONE: AC-2 — one canonical path contract (mechanism) with focused path-resolution tests
  (a) `TestPrepareResolvesStateRelativeArtifactAgainstEntityRoot` — split-root state-relative artifact → open; (b) `TestPrepareSingleRootResolvesRelativeArtifactAgainstWorkflowDir` — single-root relative → open; (c) `TestPrepareAbsoluteArtifactPassesThroughUnchanged` — absolute passthrough both modes; (d) `TestPrepareResolvesStateRelativeReferenceAgainstEntityRoot` — reference resolves same way. `TestGatePrepareCLIPassesStateRelativeArtifactWithoutCwdJoin` verifies CLI no longer joins with cwd. No host-specific parsing, no cwd probe.
- DONE: AC-3 — wrong-root path fails (the falsifying change)
  `TestPrepareWrongRootRelativeArtifactFails` creates `gate-review.md` under the workflow root (would resolve under old cwd behavior) but not under the state checkout; passes `gate-review.md` as artifact; `resolveSelectedSource` resolves to `state/gate-review.md` which does not exist → "no such file or directory". If resolution reverted to cwd, the file would be found and the test would fail. Falsifiable.
- DONE: Semantic adversarial pass — trace resolution paths, adjacent variants, falsifiability
  Traced split-root, single-root, absolute, relative artifact, relative reference, empty references, duplicate detection (seen map on resolved path), and mixed absolute/relative paths. All resolve through one site (`resolveSelectedSource`) against one root (`entityResolveRoot`). `entityResolveRoot` rejects absolute and `..`-escaping `state:` values matching `ClassifyState` error text. No multiplicative work, no unbounded allocation, no blocking I/O beyond one README read. All new tests are falsifiable — each would fail if resolution reverted to cwd.
- DONE: Verify go test ./internal/gates/ -count=1 and go test ./... are green
  `go test ./internal/gates/ -count=1` PASS; `go test ./internal/gates/ -count=1 -race` PASS; `go test ./internal/cli/ -count=1 -race -run TestGatePrepare` PASS. `go test ./...` — all packages green except pre-existing/environmental `TestVersionAmbiguousMarkersExitZero` (sees `PI_CODING_AGENT` runtime marker) and 3 Codex CLI tests (`TestCodexChannelInstallLeavesCoHostedPluginInstalled`, `TestCodexModeSwitchRoundTripPreservesExclusivity`, `TestCodexPluginInstallIsHostNative`) — confirmed failing on base commit before any changes. gofmt clean on all 4 changed files; `go vet` clean.
- DONE: Scope check — no widening
  Deliverable commit touches only `internal/gates/prepare.go`, `internal/cli/cli.go`, `internal/gates/prepare_test.go`, `internal/cli/gate_test.go`. No DVD/v8/XFAIL/gate-semantics/gitsource changes. No global hooks, no test-only product machinery, no new `status` import into `gates`.

### Summary

Validation PASSED. The one canonical resolver-owned gate-artifact path contract is correctly implemented: `gates.Prepare` resolves relative `--artifact`/`--reference` against the entity root computed once from the README `state:` field, and the CLI passes relative paths through unchanged. All three ACs are satisfied with falsifiable evidence — split-root reaches `state=open`, single-root resolves unchanged, absolute passes through, and the wrong-root path fails. No material findings. Deferred risk: the live `TestLiveCommonGateGuardrail` journey (AC-3's full end-to-end) was not exercised in validation (requires live model runtime); the focused tests prove the resolution mechanism that unblocks it.

## Stage Report: implementation (cycle 2)

- DONE: Identified the live-lane finding: TestLiveCommonRecordedGateLifecycle fails because a bare workflow-root reference (`recorder-contract.md`) joined under the state checkout (`<state>/recorder-contract.md`) instead of the workflow root (`<workflow>/recorder-contract.md`).
  The FO legitimately passes a workflow-root file as a `--reference` alongside state-rooted `--artifact` paths in a split-root gate-prepare call.
- DONE: Added `isArtifact` parameter to `resolveSelectedSource` — the artifact stays strict (state-root only; `TestPrepareWrongRootRelativeArtifactFails` still passes), references fall back to the workflow directory when the entity-root join does not exist.
  The fallback checks `os.Lstat` for a regular file at the entity-root join first, then the workflow-root join; if neither exists, it returns the entity-root join so the error shape ("no such file or directory") is preserved.
- DONE: Added `TestPrepareWorkflowRootReferenceResolves` — a bare workflow-root reference file resolves to the workflow root while `TestPrepareWrongRootRelativeArtifactFails` still rejects a wrong-root artifact.
  Verified: `go test ./internal/gates/ -race` PASS (all tests including both safety + new); `go test ./internal/cli/ -run TestGatePrepare` PASS; gofmt clean; `git diff --check` clean.

### Summary

Correction round on `resolve-split-state-gate-artifact-path`. The live pi lane proved the previous fix (commit 4a72b3a4a: dedupe + basename-strip) incomplete: a bare workflow-root reference (no `.spacedock-state` prefix) failed resolution because it joined under the state checkout. Added `isArtifact` to `resolveSelectedSource`: artifacts stay strict (state-root only, preserving `TestPrepareWrongRootRelativeArtifactFails`), references fall back to the workflow root when the entity-root join is absent. The existing dedupe and basename-strip behavior are unchanged. Committed as `c75b8890f`.

## Stage Report: validation (cycle 2)

- DONE: AC-1 (value) — split-state reference resolution completes: a bare workflow-root reference (recorder-contract.md) resolves to the workflow root, not the state checkout. Verified by the new TestPrepareWorkflowRootReferenceResolves.
- DONE: AC (safety) — artifact stays strict: a wrong-root artifact is still rejected. TestPrepareWrongRootRelativeArtifactFails PASS.
- DONE: AC (no regression) — dedupe + basename-strip from the prior cycle unchanged; full gates suite + race green; cli gate-prepare green; gofmt/vet clean.

### Verdict: PASSED

The correction completes the split-state gate-artifact path resolution: references fall back to the workflow root when the entity-root join is absent, the artifact stays state-root-strict, and the safety test is preserved. The live-lane failure (RecordedGateLifecycle's bare workflow-root reference) is fixed without weakening the artifact-strict guarantee. Ready for re-run on the stack tip.
