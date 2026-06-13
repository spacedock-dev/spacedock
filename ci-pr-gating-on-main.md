---
id: ea9kke1e8q0wyhx0wjv4cyzr
title: CI PR gating on main (post-flip trunk)
status: ideation
source: captain (2026-06-13, during the 8p push turn) — pre-cut audit main-gating item pulled forward
started: 2026-06-13T07:01:21Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0201-post-flip-release-model
group: release-model
sprint-readiness: ready
---

Post-flip, `main` is the release trunk and PRs target it, but `install-e2e.yml` and `runtime-live-e2e.yml` trigger `pull_request` on `[next]` only — so a PR to `main` runs neither the offline `go test ./...` gate nor install-e2e (only `docs.yml`, which already targets `main`, runs). Captain directive (2026-06-13): ensure CI gates fire on `main`-PRs before the sprint relies on the trunk model. This is the pre-cut antipattern-audit's known main-PR-gating item, pulled forward. Ideation pins the exact trigger change and the live-lanes design call.
