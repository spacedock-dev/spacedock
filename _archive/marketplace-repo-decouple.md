---
id: w6bhzvezybbrarkk56zemndd
title: Separate marketplace repo — decouple the manifest from the plugin branches
status: done
source: captain aside (2026-06-09) during the 0.20.0 flip — identified as the root cause of the branch-local source.ref friction.
started:
completed: 2026-06-21T04:26:53Z
verdict: REJECTED
score:
worktree:
issue:
sprint: 0201-post-flip-release-model
sprint-readiness: defer
archived: 2026-06-21T04:27:09Z
---

The marketplace is **self-referential**: `.claude-plugin/marketplace.json` lives in the plugin repo, on each branch, with `source.ref` pointing back into the same repo. This forces `main` and `next` to permanently differ on `source.ref` (`main` vs `next`), so every stable release must re-apply a `main`-only `source.ref: main` settle and `next → main` is never a clean fast-forward.

A standalone marketplace repo holding ONE manifest with explicit per-channel refs (`@main` stable, `@next` edge) would decouple the manifest from plugin content: the plugin branches carry no manifest, `main` and `next` stop diverging, and `next → main` at release becomes a clean fast-forward with no per-release settle. Evaluate the install-flow change (the binary pins a marketplace repo + selects a channel entry, instead of pinning a branch of the plugin repo), the resolver behaviour (verified during the flip: the host clones the plugin from the manifest's `source.url@source.ref`), and the migration for existing pinned installs.

## Superseded

SUPERSEDED by `marketplace-repo-and-pinned-channels` (`gp`, PR #352), which folded this task's standalone-marketplace decouple into the 0201 Model B release work. `gp` landed the marketplace-source repoint, entry-name channel selection, removal of the in-branch marketplace manifest, and the live decoupling proof; the remaining edge-install/source-override work is tracked separately by `spacedock-marketplace-source-env` (`z2t`).
