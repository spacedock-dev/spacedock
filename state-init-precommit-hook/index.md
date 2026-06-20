---
title: state init installs pre-commit validate hook — spacedock state init auto-installs the state-checkout pre-commit hook
status: backlog
source: "Shaping FO (2026-06-20): the state checkout (docs/dev/.spacedock-state) is a git worktree sharing the main repo's .git; a pre-commit hook at .git/hooks/pre-commit fires for both main-repo and state-checkout commits. A hand-installed hook now runs `spacedock status --validate` on state commits and blocks on hard errors, but it was hand-installed, not auto-managed. `spacedock state init` should install it so every fresh state checkout gets the guard without manual setup. Depends on status-validate-determinism landing first (the hook is only reliable once validate is deterministic)."
score:
started:
completed:
verdict:
worktree:
issue:
sprint:
sprint-readiness:
id: mv92pdw9pxqd2v5esw9xarx8
---

# state init installs pre-commit validate hook

## End value

`spacedock state init` (the command that initializes a split-root state checkout) auto-installs a pre-commit hook into the state checkout's shared git dir (`.git/hooks/pre-commit`) that runs `spacedock status --validate` on state-checkout commits and blocks on hard errors (missing/duplicate id). Every fresh state checkout — and every `spacedock state init` re-run — gets the guard without manual setup. The hook is state-checkout-aware (fires only for commits touching the state checkout, not every main-repo commit).

## Problem — root cause already determined

- The state checkout (`docs/dev/.spacedock-state`) is a **git worktree** of the main repo, sharing the main `.git`. Hooks live in the main repo's `.git/hooks/` and fire for all worktrees.
- A hand-installed `.git/hooks/pre-commit` now exists (this session) running `spacedock status --validate` on state commits. It works for deterministic catches (bad enum values on tracked files, duplicate ids) but is **unreliable for blank-id untracked entities** until `status-validate-determinism` lands.
- `spacedock state init` does not install any hook today — the guard is manual, so fresh state checkouts (new clones, new contributors, CI) ship without it. That's how the Commander's blank-id `pi-devoverride-package-ok` escaped: no hook, and validate's non-determinism meant even ad-hoc validate calls missed it.

## Approach (ideation confirms)

- `spacedock state init` (and `spacedock state init` re-runs) installs/refreshes the pre-commit hook at the shared `.git/hooks/pre-commit`.
- The hook is idempotent: if an existing hook is recognized as the spacedock-managed one (a marker comment, e.g. `# spacedock-managed: state-validate`), refresh it; if a different hook exists, refuse to overwrite (surface the conflict, don't clobber).
- The hook is state-checkout-aware: detect state-checkout commits (worktree toplevel basename `.spacedock-state` OR staged paths under `docs/dev/.spacedock-state`) and run validate only then; exit 0 for unrelated main-repo commits.
- The hook unsets git's hook env vars (`GIT_DIR`, `GIT_INDEX_FILE`, etc.) before running validate so the subprocess sees a clean env.
- The hook's content is owned by the binary (embedded or generated), not by the operator — operators don't hand-edit it.

## Acceptance criteria (provisional — ideation finalizes; proof = behavior)

**AC-1 — `spacedock state init` installs the pre-commit hook into the shared git dir.**
Verified by: a Go test that runs `spacedock state init` against a fresh state checkout and asserts `.git/hooks/pre-commit` exists, is executable, and contains the spacedock-managed marker + the validate invocation.

**AC-2 — The hook blocks a state commit with a corrupt entity (after status-validate-determinism lands).**
Verified by: a Go/behavior test that stages a blank-id entity in the state checkout, attempts a commit, and asserts the commit is BLOCKED with the validate error. (Depends on `status-validate-determinism` for reliability.)

**AC-3 — The hook does NOT fire (or exits 0) for unrelated main-repo commits.**
Verified by: a Go test that commits a non-state file to the main repo and asserts the hook runs validate only when staged paths touch the state checkout (or skips entirely for main-repo-only commits).

**AC-4 — Re-running `spacedock state init` refreshes an existing spacedock-managed hook without clobbering a custom one.**
Verified by: a Go test that (a) re-running init refreshes a hook with the spacedock marker; (b) init refuses to overwrite a hook without the marker, surfacing the conflict.

## Out of scope

- The validate determinism fix — `status-validate-determinism` (sibling; this task depends on it).
- Cross-worktree hook isolation (a main-repo-only commit skipping the hook is the current design; deeper isolation is not needed).
- Hook content for claude/codex hosts — the state checkout is host-neutral; one hook serves all.

## Test plan

- Go tests (AC-1/AC-3/AC-4) in `internal/cli/` or `internal/state/`: init installs the hook; hook fires only for state commits; init refreshes/clobbers correctly.
- AC-2 depends on `status-validate-determinism`; test it after that lands (or mark as integration).
- This is a `state init` / state-checkout-init change — high-stakes surface (state init runs at bootstrap). Detached adversarial audit at validation.

## Related

- `status-validate-determinism` (sibling followup) — MUST land first; the hook's blank-id blocking depends on it.
- `spacedock state init` (the command this extends) — `internal/cli/state*.go` / `internal/state/`.
- The hand-installed hook at `.git/hooks/pre-commit` (this session) — the reference content; the binary owns a managed version.
- `pi-devoverride-package-ok` — the escaped-blank-id instance that motivated both followups.
