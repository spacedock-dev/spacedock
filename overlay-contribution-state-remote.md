---
id: dqc125b36qaaezmcdjr7sfhj
title: Support a foreign state remote so an overlay workflow can create and resume split-root state
status: backlog
source: FO research session 2026-08-18 - overlay-contribution probe (workflow run wf_40ea0f6e-aa8)
started:
completed:
verdict:
score:
worktree:
issue:
---

`spacedock state new` and `state init` assume the split-root state branch lives in the code repository and is reachable through that repository's `origin`. That assumption blocks the overlay-contribution shape, where a contributor overlays an untracked workflow directory inside a clone of a repository they do not own and keeps entity state in a separate private repository.

## Problem

Three behaviors, all observed in the probe run:

- `state new` (`internal/cli/state.go:99-166`) accepts only `--workflow-dir`. It births the orphan branch in the code repository, pushes it to the code `origin`, and appends `{dir}/{state}/` to the code repository's **tracked** `.gitignore` (`state.go:155` via `appendGitignoreEntry`, implemented `state.go:203-236`), then prints `Next: commit the .gitignore edit on your code branch`. Following spacedock's own printed instruction puts a spacedock line in an upstream pull request. Observed ` M .gitignore` after `state new`; with no tracked `.gitignore` present the result is `?? .gitignore` instead. There is no variant that leaves the code repository clean.
- `appendGitignoreEntry` never consults `git check-ignore`, so it writes the entry even when `.git/info/exclude` already covers the path (observed: `git check-ignore -v spacedock/flow/.spacedock-state` returns `.git/info/exclude:7:spacedock/`).
- `state init` (`internal/cli/state.go:64`) runs `git fetch origin <branch>` in the workflow directory, i.e. against the code repository's origin, so a state branch hosted in another repository cannot be resumed on a fresh clone. `state ready` - the first call in the FO startup sequence - inherits the failure through `internal/cli/state_sync.go:268`. The printed manual fallback is also wrong for this shape: it suggests `git worktree add` off the code repo.

The rest of the system already accepts a foreign state checkout. A hand-built `git clone -b <branch> <private-repo> <statePath>` ran the full loop at exit 0, because `statesync.exactCheckoutRoot` (`internal/statesync/publish.go:149-163`) only requires the exact Git toplevel and `Preflight` (`:55-61`) only requires `HEAD == state-branch`. Only creation and resume are missing.

One candidate approach: an optional `state-remote: <url>` README frontmatter field honored by `runStateInit` (clone instead of fetch+worktree-add) and `runStateNew` (skip the gitignore append, clone-and-seed the foreign repo). Ideation owns the choice.

## Out of scope

The repository-qualifier defect in the PR probes, the worktree branch-name prefix, and the `merge guard` arm-vs-run mismatch are filed as separate tasks.

## Acceptance criteria

**AC-1 (VALUE) - After `state new` in an overlay workflow, the code repository's working tree is clean with no manual repair.**
Verified by: a behavior fixture that runs `state new` against a workflow overlaid at `<repo>/spacedock/flow` with `spacedock/` in `.git/info/exclude`, then asserts `git status --porcelain --untracked-files=all` is empty. Fails today with ` M .gitignore`. The falsifying change is any code path that writes a tracked `.gitignore` entry.

**AC-2 (VALUE) - A state branch hosted in a repository other than the code repository can be created and then resumed on a fresh clone.**
Verified by: a fixture with three local bare repos (upstream, fork, private state) that creates the state checkout through a spacedock command, discards the working clone, re-clones the fork, runs `state ready`, and asserts exit 0 with the entity history intact. Fails today with `fatal: couldn't find remote ref`.

**AC-3 - `appendGitignoreEntry` does not write an entry that `git check-ignore` already resolves.**
Verified by: a unit test with a repo whose `.git/info/exclude` covers the state path, asserting the tracked `.gitignore` is byte-unchanged. The falsifying change is removing the `check-ignore` guard.

## Test plan

Go unit tests for the `check-ignore` guard and for parsing whatever field ideation selects. Command-level behavior fixtures driving the binary against local bare repos for AC-1 and AC-2; the probe run established that local bare repos exercise remotes, pushes, and fork topology adequately, so no network or forge account is needed. A live lane is required only if the change edits the shipped FO contract's `state ready` text.
