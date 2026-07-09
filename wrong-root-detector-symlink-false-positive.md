---
title: "Wrong-root boot detector false-positives on macOS: observed path not symlink-resolved before comparison"
status: ideation
source: "Found 2026-07-09 while investigating one of sc5's (live-runner-boot-preamble-hardening) three anomalous PR #490 merge-gate failures locally, per captain direction to verify live behavior locally rather than only via GitHub Actions CI (which runs on ubuntu-latest and never hits this macOS-only symlink form). Confirmed via a local live re-run of TestLiveClaudeSharedScenarios/self-evidence-merge-triage on this machine (macOS): the test failed with 'FO booted the wrong root: expected .../private/var/folders/.../001, but it read the workflow README at .../var/folders/.../001/README.md' — the two paths are IDENTICAL directories (macOS /var is a symlink to /private/var); the FO did not actually wander."
started: 2026-07-09T11:49:49Z
completed:
verdict:
score: 0.5
worktree:
issue:
id: 5qae7c01tnytacehaphrda4s
---

`internal/ensigncycle/wrong_root_detect_impl_test.go`'s `detectWrongRootBoot` resolves the fixture root via `filepath.EvalSymlinks` before comparing (lines 39-46, with an explicit comment acknowledging the macOS `/var` vs `/private/var` symlink case) — but the two call sites that extract the FO's OBSERVED path from its tool-call stream do not apply the same resolution to that observed path before calling `isUnder`:

- `wanderWorkflowReadme(filePath, fixtureRoot string)` (line 153-165): only `filepath.Clean(filePath)`, no `EvalSymlinks`.
- `wanderTarget(command, fixtureRoot string)` (line 178-198): only `filepath.Clean(arg.path)` on each Bash-command path argument, no `EvalSymlinks`.

When Go's `t.TempDir()` returns an unresolved `/var/folders/...` path (the fixture root) and the FO's own tool calls report that same unresolved form (because that literally is its cwd/argument), while the CALLER'S `fixtureRoot` argument has already been resolved to `/private/var/folders/...`, `isUnder(p, fixtureRoot)` sees two paths with no common prefix and false-flags a wander that never happened — the FO read/operated on exactly the right directory, just spelled two ways that resolve to the same inode.

This is macOS-specific (Linux/CI has no `/var` → `/private/var` symlink, so `ubuntu-latest` CI runners never hit it) but it directly undermines the value of running these live scenario tests locally on a macOS development machine — any scenario whose observed Read/Bash path happens to preserve the unresolved `/var/...` spelling can spuriously fail with a "wrong-root" diagnosis that is actually correct behavior.

## Proposed direction (not yet fleshed — for ideation)

Resolve symlinks on the observed path (with a fallback to the unresolved form when `EvalSymlinks` errors, e.g. a path that does not exist at check-time) at both call sites before the `isUnder` comparison — matching the pattern `detectWrongRootBoot` already uses for the fixture-root side. Consider centralizing into one shared "canonicalize for wrong-root comparison" helper both call sites and the existing fixture-root resolution route through, so the fix cannot drift out of sync between the three sites in the future.

## Acceptance criteria (draft — ideation to firm up)

- A red-control: a synthetic stream where the FO's observed Read/Bash path uses the UNRESOLVED `/var/...` form while the fixture root is passed in already-resolved `/private/var/...` form (or vice versa) must NOT be flagged as a wander.
- Existing wrong-root-detection positive cases (a genuine wander to an unrelated path) must still be caught — the fix must not weaken the detector, only correct its symlink blind spot.
- Ideally verified against a REAL local live run reproducing today's false positive (self-evidence-merge-triage on this machine), not just synthetic unit fixtures, since the synthetic case is easy to get subtly wrong (e.g. resolving the wrong side, or resolving both sides but breaking the genuine-wander case).
