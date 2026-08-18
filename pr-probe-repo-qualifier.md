---
id: 12c4caspvg5z75t9846jn7c0
title: Honor the repository qualifier in every PR-state probe
status: backlog
source: FO research session 2026-08-18 - overlay-contribution probe (workflow run wf_40ea0f6e-aa8)
started:
completed:
verdict:
score:
worktree:
issue:
sprint: overlay-contribution
group: pr-tracking
---

Spacedock's PR-state probes pass an entity's `pr` field to `gh pr view` without `--repo`, so a qualified `owner/repo#N` reference errors and a bare `#N` reference silently resolves against whatever repository the launch directory happens to point at.

## Problem

Observed live against gh 2.68.1 and real GitHub repositories:

- `internal/status/boot.go:121-131` runs `exec.Command("gh","pr","view",pr,"--json","state","--jq",".state")` with no `--repo` and no `cmd.Dir`, after trimming only a leading `#`.
- `internal/dispatch/reconcile.go:99-106` duplicates the invocation; its own comment at `:93-98` states it "mirrors boot.go's PR_STATE probe invocation exactly".
- `gh pr view` has no `owner/repo#N` positional form (`USAGE: gh pr view [<number> | <url> | <branch>]`). A qualified reference falls through to the branch-name branch and returns `no pull requests found for branch "spacedock-dev/spacedock#42"`. Boot renders it `ERROR`. This is exactly the form `mods/pr-merge.md:15` prescribes.
- A bare `#42` resolves from the launch directory. Same entity file, same `--workflow-dir`, two launch directories: `MERGED` from a clone of `cli/cli`, `CLOSED` from a clone of `nodejs/node`.
- With both an `origin` and a second remote literally named `upstream`, gh selects `upstream`. In the conventional fork layout the probe is therefore accidentally correct, riding an undeclared naming convention that nothing declares, documents, or checks.
- `state sweep` inherits the defect and fails silently: `internal/dispatch/reconcile.go:637-640` discards the gh error and `continue`s, so a qualified `pr` field is invisible to merged-detection with no diagnostic at all. The FO runs `state sweep` at every «engage».
- `mods/pr-merge.md:120` records the `pr` field unqualified (`#57`), contradicting the qualified handling its own `:15` prescribes.

Blast radius today is bounded but real: `boot.go:108-114` short-circuits identify mode before the gh loop, so the FO's own `status --boot --identify --json` never reaches it, but `status --boot --json` does, and `state sweep` does at every engage.

## Out of scope

Opening a fork-to-upstream pull request, and any second-remote or fork concept in the create path. This task covers only reading PR state correctly for a reference the entity already carries.

## Acceptance criteria

**AC-1 (VALUE) - A `pr` field carrying `owner/repo#N` resolves to the PR's real state instead of ERROR.**
Verified by: a behavior fixture with a qualified `pr` field and a recorded gh invocation, asserting the probe passes `--repo owner/repo` and that the rendered state matches the PR's actual state. Fails today (renders ERROR).

**AC-2 (VALUE) - The same entity reports the same PR state regardless of the launch directory.**
Verified by: a fixture running identical `status --boot --json --workflow-dir X` from two directories whose git remotes differ, asserting byte-identical PR rows. Fails today (MERGED versus CLOSED for the same entity).

**AC-3 - A gh probe failure during `state sweep` produces a diagnostic instead of a silent skip.**
Verified by: a fixture with an unresolvable `pr` reference, asserting a non-empty warning in the sweep output and that the entity is not silently absent from the report. The falsifying change is restoring the bare `continue` at `reconcile.go:637-640`.

**AC-4 - The shipped pr-merge mod records `pr` in the same qualified form the state probes parse.**
Verified by: driving the mod's record step in a fixture and asserting the written `pr` value round-trips through the probe path without ERROR.

## Test plan

Go unit tests for the reference parser (bare, `#N`, `owner/repo#N`) and for the probe argument construction. Command-level behavior fixtures for AC-1 through AC-3, with gh stubbed on PATH so the assertion is on the arguments spacedock passes rather than on live GitHub. AC-4 drives the mod step against the same stub. The two probe sites must be corrected together or deduplicated; a fix at one site only will leave `state sweep` wrong.
