---
id: kjne6h1jkft6b6ek3p3068jh
title: One-command split-root onboarding — `spacedock state new` (birth the orphan state branch + linked worktree for an existing repo)
status: backlog
source: "captain (2026-06-02) — real-repo split-root test on DataRecce/recce: storage mechanics work cleanly (zero code-branch churn, fresh-clone halt-gate + state init both correct), but first-time onboarding of an existing repo has no single command"
completed:
verdict:
score: "0.24"
worktree:
issue:
---

Captain validated split-root on a real repo (DataRecce/recce, isolated clone): zero-churn confirmed (code branch `git status` stays empty, log never sees a state commit), the fresh-clone halt-gate fires correctly, and `spacedock state init` re-checks-out the orphan worktree (`present:false` → `present:true`) on a real 90M repo. The gap is **onboarding**:

- `spacedock state init` only RE-checks-out an *existing* orphan state branch (fetch + `git worktree add`). It does not BIRTH a split-root workflow.
- First-time setup of an existing repo — append the one `.gitignore` line, birth the orphan branch (`rm -rf --cached .` to clear the inherited tree), check it out as a linked worktree at `docs/<wf>/.spacedock-state`, scaffold the README with `state: .spacedock-state` — currently requires the `commission` skill or hand-running the Journey-1 sequence.

A `spacedock state new` (or `spacedock commission --split-root`) helper would make onboarding an existing repo a one-liner, symmetric with `state init` for the clone path.

## Acceptance criteria (provisional — harden at ideation)

**AC-1 — one command births a split-root workflow on an existing repo.** A single invocation creates the orphan state branch, the `.gitignore` entry, the linked state worktree at the declared `state:` path, and the README scaffold — leaving the code branch's `git status` clean.

**AC-2 — symmetric with `state init`.** After `state new` on repo A and a push, a fresh clone of A runs `state init` and lands in the same working state (boot `present:true`, entities render).

## Notes
- Captain flagged (no action): the orphan `spacedock-state/*` branch must be pushed to be shareable, so it shows in the repo's branch list / PR UI — expected, not a defect.
- Untested elsewhere but exercised live in THIS workflow (docs/dev, sd-b32): full FO dispatch with worktree-stages under split-root (implementation→validation→merge, code branch clean, state path-scoped). Captain's recce test covered slug-style storage mechanics only.
- Priority/milestone TBD — captured as backlog; not yet folded into a release.
