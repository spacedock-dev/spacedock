---
id: ev9h1phd37rb7r5cfrzpgenh
title: Docs site agent-readability (a14y) parity with the landing page
status: ideation
source: vendor prompt (/tmp/docs-a14y-vendor-prompt.md), captain intake 2026-06-13 — standalone doc improvement (NOT a 0.20.3 sprint item)
started: 2026-06-13T19:36:46Z
completed:
verdict:
score:
worktree:
issue:
---

Bring the MkDocs `/docs/*` pages to agent-readability (a14y scorecard 0.2.0) parity with the landing page. Standalone doc-product improvement — not part of the 0.20.3 FO-efficiency sprint.

## Problem

Under the a14y scorecard 0.2.0 (flat-pool-v1, per-page), the landing scores 94/100 but every `/docs/*` page scores ~43–68, dragging the site to 53/100 — the landing's Astro layout injects OG / JSON-LD / markdown-alternate into every `<head>`; the MkDocs docs never got the equivalent. Five root-cause fixes (A–E), most template-level (clear on every page at once). Full spec preserved in `vendor-prompt.md` beside this file.

## Governance (load-bearing)

The docs repo IS the product served to production. **Do NOT push to `main`** — changes on a branch, `gh pr create`, a maintainer merges; the downstream landing build clones `main` HEAD, so nothing reaches production until merge. Open question for ideation: this workflow's `pr-merge` mod bases PRs on `next` — confirm the correct merge target for a docs-product change (next vs main) before the merge boundary.

## For ideation

The a14y scorer is the external proof oracle: AC proof is `a14y check <built/preview page> --mode page` clearing checks A–E + `mkdocs build --strict` green — never a grep over the templates. Confirm the exact Jinja variable names (git-revision-date plugin, mkdocs-llmstxt mirror) against the installed plugin versions before relying on them. Firm approach + ACs + test plan from `vendor-prompt.md`.
