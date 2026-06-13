---
id: gf038f54jj76dw8fkgke9ek9
title: Split-root state sync should degrade when the state checkout has no origin remote
status: backlog
source: "FO dogfood (2026-06-06) - split-root state instructions require push/pull against origin, but a local state checkout may have no origin remote; workers can commit locally but remote sync is impossible."
score: "0.25"
worktree: ""
issue:
sprint: 0202-survey-improvements
group: cleanup
sprint-readiness: ready
---

Split-root state sync currently assumes the state checkout has an `origin`
remote and a shared state branch. The FO and ensign contracts tell writers to
push after path-scoped commits and to `pull --rebase origin {state_branch}` on
rejection. That is correct for the shared `docs/dev/.spacedock-state` checkout,
but it is wrong for a local-only state checkout with no `origin` remote: the
agent can commit valid local state, but any required push/pull command is
impossible noise.

The runtime should know the difference. When a split-root state checkout has no
remote, remote sync should become an explicit local-only mode: keep path-scoped
commits, skip remote push/pull, and surface a clear "state not remotely synced"
status instead of treating the missing remote as a workflow failure.

## Acceptance criteria

**AC-1 - Boot/status exposes state remote availability.**
Verified by a fixture-backed `status --boot` or equivalent state-inspection test
that distinguishes a split-root state checkout with `origin` from one without
any remote.

**AC-2 - Dispatch instructions do not require impossible remote sync.**
Verified by a dispatch-build fixture where the state checkout has no `origin`;
the emitted FO/ensign state-commit guidance keeps path-scoped local commits but
does not instruct `git push origin` or `git pull --rebase origin`.

**AC-3 - Shared-state behavior remains unchanged when origin exists.**
Verified by the existing split-root sync tests plus a focused assertion that the
normal remote-backed checkout still emits/uses push and pull-rebase guidance.

**AC-4 - Missing remote is visible, not silent.**
Verified by command output or prompt text that names the local-only state mode so
operators know state will not survive on a shared remote until a remote is
configured.

## Stage test gates

- Ideation should decide whether this belongs in `status --boot`,
  `dispatch build`, a future `state sync` helper, or all three.
- Implementation should use real-git fixtures with and without an `origin`
  remote, not string-only instruction checks.
- Validation should run the focused state/dispatch tests plus `go test ./...`.
