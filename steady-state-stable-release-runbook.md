---
id: qpfmdxy6438fsndp9nw4c89e
title: Pin the steady-state stable-release runbook (advance main from a green next tip); reconcile docs/releasing.md
status: backlog
source: FO+captain (2026-06-09) post-0.20.0-flip — the recurring release process is not pinned; docs/releasing.md still describes a main-integration flow.
started:
completed:
verdict:
score:
worktree:
issue:
sprint-readiness: defer
---

The 0.20.0 flip was a one-time divergent replacement. The RECURRING stable release (0.20.1 / 0.21.0) needs a pinned runbook consistent with the post-flip model:

- Dev integrates on `next` (worktrees branch off `next`, PRs to `next`).
- At a cut: take a **green `next` tip** (live e2e green), advance `main` to it via `--force-with-lease` (a straight `next → main` fast-forward would revert `source.ref` to `next`), **re-apply the `main`-only settle** (`source.ref: main` + calendar bump), then tag the green-run commit (tag-the-green-tip + the `e2e-gate`).

`docs/releasing.md` currently describes a MAIN-integration flow ("ensure all release content is merged to `main`", worktree off `main`) that does NOT match the next-integration + advance-`main`-at-release reality — reconcile it (or supersede it). Note: `marketplace-repo-decouple` (w6) would simplify this to a clean fast-forward with no per-release settle; sequence accordingly.

## Folded into releasing-doc-pre-stamp-drift (0230)

This task's runbook substance is folded into the consolidated `releasing-doc-pre-stamp-drift` member as its approach step 3 (fold the steady-state runbook into `docs/releasing.md`). This entity is no longer a separate 0230 member. Verify-first caveat carried by the lead: `docs/releasing.md:3` already says stable is "cut from `main`" and origin/stable advances from `main`, which contradicts this entity's "dev integrates on `next`" premise — so the fold is likely a DELETION of stale next-framing, not the addition of a next->main step. Do not document a next-line the live machinery doesn't have.
