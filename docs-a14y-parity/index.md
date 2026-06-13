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

## Governance (load-bearing) — RESOLVED: PR bases on `main`, not `next`

The docs repo IS the product served to production. **Do NOT push to `main`** — changes on a branch, `gh pr create`, a maintainer merges.

**Merge-target resolution (this was the open ideation question):** the PR MUST base on **`main`**, which is an explicit exception to the `pr-merge` mod's default `next` base (`docs/dev/_mods/pr-merge.md` line 62 hardwires `gh pr create --base next`). Three independent facts force `main`:
- `.github/workflows/docs.yml` builds `mkdocs build --strict` on `pull_request → main` and deploys to GitHub Pages **only on push to `main`** (lines 8–13, 42–49). The docs site is served from `main`.
- The downstream landing build clones this repo's `main` HEAD (vendor spec); nothing reaches production until merge to `main`.
- The docs source files this PR edits (`overrides/main.html`, `mkdocs.yml`, all of `docs/site/`) **exist on `main` but NOT on `next`** (verified: `git cat-file -e origin/next:overrides/main.html` → absent; present on `origin/main`; the two branches have diverged, `next` is behind on the docs). A PR based on `next` would target files that do not exist there.

Implementation must base the branch on `origin/main` and `gh pr create --base main`. Record this as the pr-merge boundary override for this entity.

## Approach (firmed from vendor-prompt.md — five root-cause fixes A–E)

All five clear specific per-page a14y `0.2.0` checks. A and B are template/config-level (clear on every page at once); C is a build hook; D is one page + a site-wide link; E is config + a source sweep.

### A. `overrides/main.html` — extend the existing `{% block extrahead %}`
After `{{ super() }}` and the existing Bunny font links, inject (for every page) OG tags, one `<script type="application/ld+json">` carrying a `TechArticle` (`headline`, `description`, `url`, `dateModified`) + a `BreadcrumbList` from `page.ancestors | reverse`, and a markdown `<link rel="alternate" type="text/markdown" href="{{ page.canonical_url }}index.md">`. Clears `html.og-title`, `html.og-description`, `html.json-ld`, `html.json-ld.date-modified`, `html.json-ld.breadcrumb`, `markdown.alternate-link`.

**CORRECTION to the vendor sketch (proven in spike, load-bearing):** the entire injected block MUST be wrapped in `{% if page %} … {% endif %}`. The Material `404.html` renders `extrahead` with `page = None`; without the guard `mkdocs build --strict` dies with `'None' has no attribute 'meta'` on the first `page.meta`/`page.title`/`page.canonical_url`/`page.ancestors` reference. The vendor sketch omits this guard.

### B. `mkdocs.yml` + `docs/requirements.txt` — git-revision-date plugin
Pin `mkdocs-git-revision-date-localized-plugin` in `docs/requirements.txt` and enable in `plugins:`:
```yaml
  - git-revision-date-localized:
      enable_creation_date: false
      type: iso_datetime
```
**Verified against installed plugin (1.5.3):** with `type: iso_datetime` it populates `page.meta.git_revision_date_localized_raw_iso_datetime` — exactly the variable the JSON-LD `dateModified` reads (source: plugin.py line 345 `git_revision_date_localized_raw_{date_type}`). In the spike it resolved to a real value (`2026-06-13 03:34:50`). Note the format is `YYYY-MM-DD HH:MM:SS` (space-separated, no `T`/offset); the a14y `html.json-ld.date-modified` check accepted it. **Sandbox note (machine-specific, not a CI concern):** under this sandbox the plugin's default parallel processing fails with `PermissionError` on semaphore creation; add `enable_parallel_processing: false` for local sandboxed builds. CI (`ubuntu-latest`) does not need it — leave the committed config at the plugin default unless a local build is required, or set it and accept the (negligible) serial cost.

### C. The per-page `.md` mirror — frontmatter + `## Sitemap` via a post-build hook
**Verified mechanism:** `mkdocs-llmstxt` (0.3.0) writes each per-page mirror at `<page-path>/index.md` in its `on_post_build` (plugin.py line 185); the emitted mirror has **no frontmatter and no Sitemap** and the plugin offers no config to inject them. The clean mechanism is a MkDocs `hooks:` script with its own `on_post_build` that runs after the plugin, re-reads each emitted `<page>/index.md`, prepends YAML frontmatter (`title` from the H1, `description` from page/site meta, `doc_version` from the docs/release version, `last_updated` from the git date) and appends a short `## Sitemap` section. Proven end-to-end in spike: the enriched mirror cleared both `markdown.frontmatter` and `markdown.sitemap-section` with `mkdocs build --strict` green. Mirror URL shape confirmed: `…/concepts/workflows-and-entities/index.md` (matches the alternate-link in A).

### D. Glossary page + site-wide link
Add `docs/site/reference/glossary.md` defining the core terms (workflow, entity, stage, gate, first-officer, ensign, mod, standing teammate, split-root), add it to `nav:` under Reference, and add a site-wide link in `overrides/partials/footer.html` mirroring the existing `sd-footer-llms` "For agents: llms.txt" link (footer.html line 49) so EVERY page satisfies `html.glossary-link`. Confirmed still-failing in spike (`no glossary/terminology link`) — the check is real.

### E. Code-fence languages — `use_pygments: false` + source sweep (CORRECTION to vendor framing)
**Load-bearing spike finding:** the vendor spec frames E as "sweep bare ` ``` ` fences and add the language hint." That is necessary but **NOT sufficient.** With the current `pymdownx.highlight` (Pygments active), even fences that DO declare a language render as `<div class="highlight"><pre><span></span><code>…` with **no `language-*` class** — Pygments highlights inline via `<span class="…">` and the `language_prefix` only applies to *non-Pygments* blocks. So `code.language-tags` (which looks for `language-*` on `<code>`) fails on correctly-labeled fences. Proven fix: set `use_pygments: false` under `pymdownx.highlight`; the spike then emitted `<code class="language-yaml">` and the check PASSED (`2 blocks`), score 77→83. **This is a real UX decision for the captain:** `use_pygments: false` removes server-side syntax coloring (Material then needs a client-side highlighter for colored code, or code renders monochrome). Source fences must ALSO declare languages — a sweep of `docs/site/**/*.md` for bare fences is still required (found 4 bare opening fences across 3 files: `concepts/gates-and-decisions.md` ×2, `get-started/first-workflow.md`, `running-workflows/commission.md`). **Open question to raise at the gate:** accept the `use_pygments: false` UX tradeoff, OR keep Pygments and treat `code.language-tags` as accept-skip like the other theme-structural checks. Recommend deciding at the ideation gate before implementation.

### Explicitly NOT in this PR (vendor scope, confirmed by the scorer)
The remaining failing checks in the spike are exactly the vendor's out-of-scope items: `markdown.canonical-header` (landing repo `netlify.toml` Link header), `markdown.content-negotiation` (landing CDN/edge `Accept: text/markdown`), `html.text-ratio` (13.5% — Material nav/search chrome; the `.md` mirror is the agent-readable path; accept-skip, do not chase at the cost of human UX), `http.content-type-html` (the HTML page is text/html by design).

## Riskiest-first spike record (throwaway, executed before firming the plan)

Built a disposable overlay (`mkdocs-spike.yml` + a spike `overrides/main.html` + a spike post-build hook), ran `mkdocs build --strict`, served `site/` locally, and scored a built page with the real a14y CLI. All spike files removed; the real `overrides/main.html`/`mkdocs.yml`/`docs/requirements.txt` were never modified.

- **Proof tool runnable here:** `a14y` is the npm package `a14y` (CLI 0.4.21), runnable via `npx --yes a14y check <url> --mode page -o json`, defaulting to scorecard `0.2.0` (matches the vendor spec). It fetches a **served URL** (not a file path), so the local proof path is `mkdocs build` → `python3 -m http.server` on `site/` → `a14y check http://localhost:PORT/<page>/`.
- **Scores (real scorer, `concepts/workflows-and-entities/` page):** baseline current docs **58** → with A+C **77** → with A+C+E(`use_pygments:false`) **83**. The only in-scope check still failing at 83 is D (glossary); D + the per-page fence sweep close the remaining gap into the 90s.
- **Verified against installed pinned versions:** git-date variable `page.meta.git_revision_date_localized_raw_iso_datetime` (✓ resolves to a real date); `page.ancestors | reverse` produces the breadcrumb (`['Concepts','Workflows & entities']`); `page.canonical_url` resolves; emitted JSON-LD parses (`json.loads` succeeded); mirror URL shape `…/index.md` (✓). Corrections found: the `{% if page %}` 404 guard (A) and the `use_pygments: false` requirement (E) — both absent from the vendor sketch.

## Acceptance criteria (proven by the EXTERNAL a14y scorer + the build — never a grep over templates/source)

The proof oracle is `a14y check <built-or-preview page> --mode page` (scorecard `0.2.0`) plus `mkdocs build --strict`. Local reproduction: `mkdocs build` → serve `site/` over HTTP → `a14y check http://localhost:PORT/<page>/ --mode page -o json`. The deploy-preview / CI path is the canonical reproduction once the PR is open; the local served-build path is the equivalent and was exercised in the spike.

- **AC-1 — Template injects OG + JSON-LD + alternate-link on every page, page-None safe.** On a built docs page, `a14y check … --mode page` reports `pass` for `html.og-title`, `html.og-description`, `html.json-ld`, `html.json-ld.breadcrumb`, and `markdown.alternate-link`; `mkdocs build --strict` exits 0 (the build proves the 404 page's `page=None` path does not throw). *Verified by:* a14y JSON output (those check `status` = `pass`) + `mkdocs build --strict` exit 0.
- **AC-2 — git-date plugin feeds a real `dateModified`.** `a14y check` reports `html.json-ld.date-modified` = `pass` with a non-empty ISO date in the emitted JSON-LD. *Verified by:* a14y JSON output (`html.json-ld.date-modified` status `pass`, message carrying the date) + `mkdocs build --strict` exit 0 with the plugin enabled.
- **AC-3 — `.md` mirrors carry the required frontmatter + a `## Sitemap` section.** `a14y check` reports `markdown.frontmatter` = `pass` and `markdown.sitemap-section` = `pass`; the emitted `<page>/index.md` begins with a YAML frontmatter block (`title`, `description`, `doc_version`, `last_updated`) and contains a `## Sitemap` section. *Verified by:* a14y JSON output (both checks `pass`) — the emitted mirror file on disk is the artifact the check parses.
- **AC-4 — Glossary page exists and is linked site-wide.** `a14y check` reports `html.glossary-link` = `pass` on a non-glossary page (e.g. the concepts page), proving the link is present on every page, not just the glossary itself. *Verified by:* a14y JSON output (`html.glossary-link` status `pass`).
- **AC-5 — Source code fences declare languages and render with `language-*` classes.** `a14y check` reports `code.language-tags` = `pass` on a page containing code blocks. *Verified by:* a14y JSON output (`code.language-tags` status `pass`). (Mechanism gated on the gate decision in fix E; if the captain elects accept-skip, AC-5 is removed and `code.language-tags` is documented as an accepted exception alongside the text-ratio / canonical-header / content-negotiation skips.)
- **AC-6 — Build stays green and the change ships via a PR to `main`.** `mkdocs build --strict` exits 0 with all fixes applied; the change lands on a branch with a PR opened against `main` (NOT `next`, NOT a push to `main`). *Verified by:* `mkdocs build --strict` exit 0 + the open PR's base branch = `main` (`gh pr view --json baseRefName`).

## Test plan

- **Cost/complexity:** low-to-moderate. No Go code; all changes are in the docs product (templates, mkdocs config, one Python build hook, one markdown page, a fence sweep). The mechanism risk was the unknown and is now retired by the spike.
- **What verifies it:** the a14y scorer (`0.2.0`) is the external oracle for every AC; `mkdocs build --strict` is the build gate. No fixture/Go tests apply — this is a docs-product change. The canonical CI proof is the `docs` workflow's strict build on the PR plus the deploy-preview/served-build a14y run; the spike proved the identical local served-build path works (`npx a14y check http://localhost:PORT/<page>/ --mode page`).
- **Reproduction recipe (for implementation + validation):** install `docs/requirements.txt` (+ the pinned git-date plugin) into a venv → `mkdocs build` (add `enable_parallel_processing: false` to the git-date plugin if building under a sandbox) → `python3 -m http.server` on `site/` → `npx --yes a14y check http://localhost:PORT/concepts/workflows-and-entities/ --mode page -o json` → assert the AC check IDs report `status: pass`. Re-run after each fix; A+B+C+D+E together should push the in-scope checks to `pass` and the page score into the 90s.
- **No new mechanism left unexercised:** the design composes already-proven behavior — every fix mechanism (template injection with 404 guard, git-date variable, post-build mirror enrichment, `use_pygments` language class, glossary site-wide link) was driven end-to-end against the real scorer in the spike. The one deferred decision (fix E UX tradeoff) is a captain choice, not an unverified mechanism.

## Stage Report: ideation

- DONE: Firm the approach from vendor-prompt.md into the entity body (the 5 fixes A–E + out-of-scope), and exercise the riskiest unknowns FIRST: confirm the proof tool `a14y check --mode page` is actually runnable here, and confirm the assumed Jinja/plugin variables exist against the INSTALLED versions; record what was verified; if a variable/tool is absent, adjust the fix and say so.
  All five fixes A–E + out-of-scope firmed in the body. Riskiest-first spike (disposable overlay, since removed) drove `mkdocs build --strict` + the real a14y scorer end-to-end: tool runnable via `npx a14y check <served-url> --mode page` (scorecard 0.2.0); git-date var `page.meta.git_revision_date_localized_raw_iso_datetime` ✓, `page.ancestors`/`page.canonical_url` ✓, mirror URL `…/index.md` ✓; JSON-LD parses. Two corrections recorded: `{% if page %}` 404 guard (fix A) and `use_pygments: false` for `language-*` classes (fix E) — both absent from the vendor sketch.
- DONE: Acceptance criteria proven by the EXTERNAL a14y scorer + the build: each AC names `a14y check <built/preview page> --mode page` clearing its specific checks (A–E) and `mkdocs build --strict` green — never a grep over the templates or source. If `a14y` is not runnable in this environment, state the feasible proof path (deploy-preview / CI) explicitly so the AC is reproducible.
  AC-1..AC-6 each cite the specific a14y check ID(s) reporting `status: pass` from `-o json` output + `mkdocs build --strict` exit 0. a14y IS runnable here (npm `a14y` 0.4.21); the local served-build path (`mkdocs build` → `http.server` → `a14y check http://localhost:PORT/<page>/`) was exercised and is the equivalent of the deploy-preview/CI proof. No AC is a grep over templates/source.
- DONE: Resolve the governance merge-target: a docs-product change is PR-only (never push `main`); confirm whether the PR bases on `next` or `main` (the downstream landing clones `main` HEAD) and record it, so the pr-merge boundary targets the right branch.
  Resolved to **`main`** (explicit override of the pr-merge mod's default `next` base, `_mods/pr-merge.md` line 62). Three forcing facts recorded: `docs.yml` deploys Pages only on push to `main`; landing clones `main` HEAD; the edited docs source files exist on `main` but NOT on `next` (verified `git cat-file -e origin/next:overrides/main.html` → absent). A PR to `next` would target nonexistent files.

### Summary

Firmed the vendor spec's five fixes (A: template OG/JSON-LD/alternate-link; B: git-date plugin; C: post-build mirror frontmatter+Sitemap hook; D: glossary + site-wide link; E: code-fence languages) into the entity with entity-level ACs whose only proof is the external a14y scorer (0.2.0) + `mkdocs build --strict`. The riskiest-first spike retired all mechanism risk against the installed pinned versions and produced before/after scores (58 → 77 with A+C → 83 with A+C+E). Two load-bearing corrections to the vendor sketch surfaced and are recorded: the `{% if page %}` 404 guard (A) and that Pygments strips `language-*` classes so E needs `use_pygments: false` (a UX tradeoff to decide at the gate). Governance merge-target resolved to `main`, overriding the pr-merge mod's `next` default, with the three forcing facts on record.
