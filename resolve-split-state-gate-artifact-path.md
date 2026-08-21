---
title: Resolve split-state gate artifact paths before preparation
status: ideation
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
started: 2026-08-21T08:15:56Z
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
