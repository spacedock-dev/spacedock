---
id: twq68r4y8qg0wetztajtmmzz
title: Make `next` its own release line (0.21.0+) decoupled from the stable v* tag
status: backlog
source: FO+captain (2026-06-09) 0.20.0 main-flip drive — the deferred "(b)" the captain chose but held to get main done first.
started:
completed:
verdict:
score:
worktree:
issue:
sprint: 0201-post-flip-release-model
sprint-readiness: defer
---

The edge channel currently rides the stable `v*` tag: one `goreleaser` run builds BOTH casks (`spacedock` + `spacedock@next`) at the same version, and `next`'s `plugin.json` is not stamped. So post-flip the edge install resolves the `next` plugin under `next`'s own version (observed: `0.19.9`), trailing the stable `0.20.0` — functionally correct (edge gets next content) but version-incoherent.

Give `next` an independent version line (`0.21.0`+) with its own publish trigger so the edge channel has coherent versioning. Touches the release design: `release.yml` is tag-push-only / single-target `main`; `next-publish.yml` is the edge calendar-bump path. Decide the mechanism — a `next`-side version-tag scheme vs. extending `next-publish.yml` to stamp + publish an edge cask version — without re-coupling the two channels.

## Superseded

Superseded into `patch-release-line-support` (d0g21c517b5nvga1ybwckapk) by captain order, 2026-08-25. This body described the retired next-branch model; edge gained its own version line via the auto-pre0 tags, and the surviving work - old-line patch support and the automated preversion bump - lives in the successor with the 2026-08-25 v0.27.0-cut incidents as its spec. Verdict left empty per the supersede convention.
