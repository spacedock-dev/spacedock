---
title: status --validate non-deterministic on untracked entity files (root cause of escaped blank-id)
status: backlog
source: "Shaping FO (2026-06-20): investigating why a Commander-filed entity (pi-devoverride-package-ok) shipped with a blank id through its whole lifecycle (create → pr-pending → mod-block → terminalize → archive) and only surfaced when a fresh boot ran --validate. A pre-commit hook running `spacedock status --validate` could not reliably block the same blank-id entity: the same file on disk returns `Error: missing required id` on some invocations and `VALID` on others, across both clean-shell and git-hook (GIT_DIR/GIT_INDEX_FILE set) environments. The non-determinism is in the status tool's entity scan, not the hook."
score:
started:
completed:
verdict:
worktree:
issue:
sprint:
sprint-readiness:
id: a1a9k6tqnpj7292z80xxwkf2
---

# status --validate non-deterministic on untracked entity files

## End value

`spacedock status --validate` deterministically scans every on-disk entity file (tracked or not, staged or not, committed or not) and reliably reports a blank/missing id as a hard error every time. A pre-commit hook (or any caller) can then trust validate to block corrupt entities before they land on the state branch. This is the root-cause fix that makes the state-checkout pre-commit hook reliable.

## Problem — root cause already determined

The Commander filed `pi-devoverride-package-ok` by hand-writing the file (bypassing `spacedock new`'s id-minting), so it shipped with `id:` blank from birth (commit `74af536f`) through terminalize (`b0e759e8`) and archive (`2784813b`) — 5 commits, never caught. It only surfaced when a fresh boot ran `--validate`. Investigating why a pre-commit hook couldn't catch it, I reproduced:

- A blank-id entity file (`zprobe-hook-test/index.md` with `id:`) placed on disk in the state checkout.
- `spacedock status --workflow-dir docs/dev --validate` run repeatedly against the SAME on-disk file.
- **Non-deterministic results:** some invocations return `Error: missing required id: ... slug=zprobe-hook-test`; others return `VALID`. Observed across:
  - clean shell (`env -u GIT_DIR -u GIT_INDEX_FILE -u GIT_WORK_TREE`)
  - git-hook env (`GIT_DIR=...worktrees/-spacedock-state`, `GIT_INDEX_FILE=.../index` set, mimicking git's hook invocation)
  - both `${SPACEDOCK_BIN}` and bare `spacedock` (both 0.22.0)

The validate path: `validateWorkflow` (internal/status/validate.go:154) → `activeAndArchivedEntities` → `scanEntitiesActive` (discover.go:193) → `discoverEntityFiles` → `os.ReadDir`-based. The entity is on disk; `os.ReadDir` should find it every time. The non-determinism must be in `discoverEntityFiles` (a filter that inconsistently includes/excludes), or in `loadActiveEntityFields`'s worktree-mirror overlay (discover.go:166) interacting with `FindGitRoot` under varying git context — a worktree-backed entity overlay could mask the pipeline-dir copy's fields on some code paths.

Because validate is non-deterministic, the Commander's blank-id file scanned as VALID at every lifecycle gate that ran validate (if any did), and the pre-commit hook cannot reliably block the same defect.

## Approach (ideation confirms)

- Reproduce deterministically: find the minimal invocation + on-disk state that flips validate between Error and VALID for the same file. Likely a race or a git-context-dependent path in `discoverEntityFiles` / `loadActiveEntityFields`.
- Inspect `discoverEntityFiles` and `loadActiveEntityFields` (discover.go) for a filter or overlay that depends on git state (tracked status, worktree presence, index state) rather than pure filesystem.
- Fix so validate scans every on-disk `.md`/`index.md` entity deterministically, regardless of git tracking/staging/worktree state. The id check (`validateSDB32`, validate.go:261) already errors on `id == ""` — the bug is upstream, in the entity enumeration.
- Add a regression test: a blank-id untracked entity file on disk → validate errors EVERY time (run validate N times, assert all error).

## Acceptance criteria (provisional — ideation finalizes; proof = behavior)

**AC-1 — `spacedock status --validate` deterministically errors on a blank-id on-disk entity.**
Verified by: a Go test that writes a blank-id entity file to the state checkout (untracked, not staged), runs `--validate` N times (e.g. 10), and asserts EVERY run returns the `missing required id` error (no VALID results). The inverse of the current non-determinism.

**AC-2 — The determinism holds across git-hook env (GIT_DIR/GIT_INDEX_FILE set).**
Verified by: the same test run with the worktree git env vars set (mimicking a pre-commit hook invocation); still errors every time.

**AC-3 — Tracked/staged/committed entities are unaffected (no regression).**
Verified by: existing validate tests stay green; a newly-`spacedock new`'d entity (tracked, staged, committed) validates cleanly.

## Out of scope

- The pre-commit hook auto-install at `spacedock state init` — separate task (`state-init-precommit-hook`); depends on this task landing first (the hook is only reliable once validate is deterministic).
- The `status: implementation` (not `done`) latent inconsistency on archived `pi-devoverride-package-ok` — separate data-quality cleanup; not the root cause.
- Changing `validate`'s error/warning classification — only the entity-enumeration determinism.

## Test plan

- Go test (AC-1/AC-2/AC-3) in `internal/status/`: blank-id untracked entity → validate errors N/N times; same under git-hook env; tracked-entity regression.
- Reproduction script (for ideation): the exact sequence I ran (place blank-id file, run validate 3× in clean env, 3× in git-hook env) — record which invocation flips and why.
- This is a `status` mutation/guard-path change — high-stakes surface per the dev-workflow proof policy. Detached adversarial audit at validation.

## Related

- `pi-devoverride-package-ok` (archived, `id` now `18s156v0txff12fv61yqptk2`) — the escaped-blank-id instance that surfaced this.
- `internal/status/validate.go` (`validateWorkflow`, `validateSDB32`), `internal/status/discover.go` (`scanEntitiesActive`, `discoverEntityFiles`, `loadActiveEntityFields`) — the scan path.
- `state-init-precommit-hook` (sibling followup) — the pre-commit hook auto-install; depends on this task.
- The pre-commit hook currently installed at `.git/hooks/pre-commit` (worktree-shared) — documents this limitation in its header.
