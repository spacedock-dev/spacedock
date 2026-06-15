---
title: docs.yml — drop the GitHub Pages deploy-API job, push the built site directly to gh-pages
status: backlog
source: "captain (2026-06-15) — the docs workflow's Deploy-to-GitHub-Pages step (actions/deploy-pages@v4) 404s on every main push (Pages source not configured as Actions); build step is green. Drop the Pages-API deploy and push the built site directly to the gh-pages branch (classic branch-served Pages, no Actions-source dependency)."
score: 0.30
sprint: 0203-fo-efficiency
id: j1ys5wh6wprhy8t8jd98z7g6
---

The `docs` workflow's `deploy` job uses `actions/deploy-pages@v4`, which fails `HttpError: Not Found (404)` on every push to main because the repo's Pages source is not set to GitHub Actions. The strict `build` step passes; only the Pages-API publish fails. Replace the deploy-API path with a direct push of the built site to `gh-pages`.
