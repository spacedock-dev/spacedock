---
title: docs.yml — drop the GitHub Pages deploy-API job, push the built site directly to gh-pages
status: validation
source: "captain (2026-06-15) — the docs workflow's Deploy-to-GitHub-Pages step (actions/deploy-pages@v4) 404s on every main push (Pages source not configured as Actions); build step is green. Drop the Pages-API deploy and push the built site directly to the gh-pages branch (classic branch-served Pages, no Actions-source dependency)."
score: 0.30
sprint: 0203-fo-efficiency
id: j1ys5wh6wprhy8t8jd98z7g6
worktree: .worktrees/spacedock-ensign-docs-pages-direct-push
started: 2026-06-15T03:17:20Z
---

The `docs` workflow's `deploy` job uses `actions/deploy-pages@v4`, which fails `HttpError: Not Found (404)` on every push to main because the repo's Pages source is not set to GitHub Actions. The strict `build` step passes; only the Pages-API publish fails. Replace the deploy-API path with a direct push of the built site to `gh-pages`.

## Stage Report: implementation

- DONE: REMOVE the Pages-API publish path — delete the `deploy` job AND the `Upload Pages artifact` step; keep `mkdocs build --strict` as-is
  Commit cfbf7d32: deploy job, github-pages environment, pages:write/id-token:write perms, and upload-pages-artifact step all gone; build job's `mkdocs build --strict` unchanged.
- DONE: On push to main only, push built `site/` DIRECTLY to `gh-pages` via peaceiris/actions-gh-pages@v4 (github_token, publish_dir ./site, publish_branch gh-pages); update header comment; keep `concurrency: group: pages`
  New "Publish to gh-pages" step gated on `if: github.ref == 'refs/heads/main'`; header comment rewritten to describe classic branch-served Pages; concurrency guard retained. Added `permissions: contents: write` to the build job (peaceiris pushes the branch via GITHUB_TOKEN).
- DONE: Verify strict build clean, deploy-API refs gone, yaml valid
  `mkdocs build --strict` → exit 0 (pinned docs/requirements.txt in throwaway venv). `grep -rn 'deploy-pages\|upload-pages-artifact\|github-pages' .github/workflows/docs.yml` → no matches (exit 1). yaml parsed valid via yq and ruby (actionlint not installed locally).

### Summary

Replaced the failing GitHub Pages deploy-API path in `.github/workflows/docs.yml` with a direct push of the built `site/` to the `gh-pages` branch using `peaceiris/actions-gh-pages@v4`. The strict-build PR gate is untouched and still passes. Added `contents: write` to the build job so the action can push the branch via `GITHUB_TOKEN`.

OPERATOR FOLLOW-UP (captain-owned, one-time): the repo Pages source must be set to "Deploy from a branch: gh-pages / (root)" in Settings > Pages for the pushed branch to be served. The workflow push itself succeeds regardless; only serving depends on this setting.

## Stage Report: validation

- DONE: `mkdocs build --strict` (pinned docs/requirements.txt) exits 0 — the PR gate is intact
  Clean venv from docs/requirements.txt, `mkdocs build --strict` → exit 0; site/index.html produced.
- DONE: Pages deploy-API path is GONE — grep deploy-pages|upload-pages-artifact|github-pages → no matches; standalone `deploy` job removed
  `grep -nE 'deploy-pages|upload-pages-artifact|github-pages' .github/workflows/docs.yml` → exit 1 (no matches); single `build` job only.
- DONE: New publish path correct — single "Publish to gh-pages" step using peaceiris/actions-gh-pages@v4, github_token GITHUB_TOKEN, publish_dir ./site, publish_branch gh-pages, gated if main; permissions contents: write on build job; concurrency group: pages retained; PRs never publish
  docs.yml:46-52 step matches spec; build job has `permissions: contents: write` (l27-28); `concurrency: group: pages` (l18-20); publish gated `if: github.ref == 'refs/heads/main'`.
- DONE: Structural validity — yaml parses; no dangling needs: build / steps.deployment refs survive
  `yq '.'` → OK; grep for `needs: build|steps.deployment|id-token|pages: write|environment:` → exit 1 (none). actionlint not installed locally (sandbox-denied).
- DONE: Operator follow-up recorded — repo Pages source must be set to "Deploy from a branch: gh-pages"
  Captured in implementation report (entity body, OPERATOR FOLLOW-UP note).

### Summary

Validation PASSED. The failing Pages deploy-API path is fully removed and replaced by a single peaceiris/actions-gh-pages@v4 publish step gated to main; the strict-build PR gate still exits 0 and PRs never publish. actionlint was unavailable in this sandbox, but yq parse plus grep for dangling deploy-job references confirm structural validity. The one-time operator action (set repo Pages source to "Deploy from a branch: gh-pages") is recorded for the captain.
