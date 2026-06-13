# Prompt — lift the docs site's agent-readability (a14y scorecard 0.2.0)

**Repo:** `spacedock-dev/spacedock` (the docs product, MkDocs Material, served at `https://spacedock.md/docs/`).
**Why:** Under the new a14y scorecard `0.2.0` (`flat-pool-v1`, per-page), the landing scores 94/100 but every `/docs/*` page scores ~43–68 — dragging the site-wide score to 53/100. The landing passes these checks because its Astro layout injects OG tags, JSON-LD, and a markdown alternate-link into every page `<head>`; the MkDocs docs never got the equivalent. This task brings the docs to parity.

**Governance (load-bearing):** This repo is the product. Do NOT push to `main`. Make the changes on a branch and open a PR (`gh pr create`); a maintainer merges. The downstream landing build clones this repo's `main` HEAD, so nothing reaches production until the PR merges.

**How to verify any fix locally:** build the docs (`mkdocs build --strict`) and run the scorer against a built page, or against the deploy preview once the PR has one:
```
a14y check <page-url> --mode page          # e.g. a built/previewed concepts/workflows-and-entities page
```
Target: clear the per-page checks listed below. Re-run after each root-cause fix — most are template-level and clear on every page at once.

---

## The failing per-page checks and their fixes

### A. `overrides/main.html` — ONE template edit clears 6 checks on every docs page

The override already exists and extends `base.html` with a `{% block extrahead %}` (it currently injects the Bunny font links). **Extend that same block** — after `{{ super() }}` and the font links — to inject, for every page:

1. **Open Graph** (`html.og-title`, `html.og-description`):
   ```html
   <meta property="og:title" content="{{ page.title | striptags }} - {{ config.site_name }}">
   <meta property="og:description" content="{{ page.meta.description | default(config.site_description) }}">
   <meta property="og:type" content="article">
   <meta property="og:url" content="{{ page.canonical_url }}">
   ```

2. **JSON-LD** (`html.json-ld`, `html.json-ld.date-modified`, `html.json-ld.breadcrumb`) — one `<script type="application/ld+json">` carrying a `TechArticle` (or `Article`) node with `headline`, `description`, `url`, and `dateModified`, plus a `BreadcrumbList` built from the page's nav ancestors. Material exposes `page.ancestors` (reversed = root→current) and `page.canonical_url`. Sketch:
   ```html
   <script type="application/ld+json">
   {
     "@context": "https://schema.org",
     "@graph": [
       {
         "@type": "TechArticle",
         "headline": {{ (page.title | striptags) | tojson }},
         "description": {{ (page.meta.description | default(config.site_description)) | tojson }},
         "url": {{ page.canonical_url | tojson }},
         "dateModified": "{{ page.meta.git_revision_date_localized_raw_iso_datetime | default('') }}",
         "isPartOf": { "@type": "WebSite", "name": {{ config.site_name | tojson }}, "url": {{ config.site_url | tojson }} }
       },
       {
         "@type": "BreadcrumbList",
         "itemListElement": [
           {% for anc in page.ancestors | reverse %}
           { "@type": "ListItem", "position": {{ loop.index }}, "name": {{ anc.title | striptags | tojson }} }{% if not loop.last or page %},{% endif %}
           {% endfor %}
           { "@type": "ListItem", "position": {{ (page.ancestors | length) + 1 }}, "name": {{ page.title | striptags | tojson }}, "item": {{ page.canonical_url | tojson }} }
         ]
       }
     ]
   }
   </script>
   ```
   Verify the exact Jinja variable for the modified date against the installed plugin (see B); if absent, fall back to `{{ build_date_utc }}`. Validate the emitted JSON parses (the check requires *parseable* JSON-LD).

3. **Markdown alternate-link** (`markdown.alternate-link`):
   ```html
   <link rel="alternate" type="text/markdown" href="{{ page.canonical_url }}index.md">
   ```
   (Confirm the mirror URL shape — the scorer found the mirror at `<page-url>index.md`, e.g. `…/workflows-and-entities/index.md`. Match whatever your `.md` emission actually produces.)

### B. `mkdocs.yml` — add the git-revision-date plugin (feeds `dateModified` above)

Add `mkdocs-git-revision-date-localized-plugin` to `docs/requirements.txt` (pin it) and enable it in `plugins:`:
```yaml
  - git-revision-date-localized:
      enable_creation_date: false
      type: iso_datetime
```
It populates `page.meta.git_revision_date_localized_raw_iso_datetime` (or similar — confirm the variable name for the pinned version) which the JSON-LD `dateModified` reads. This also indirectly supports the `.md` mirror `last_updated` frontmatter (C).

### C. The per-page `.md` mirror — add frontmatter + a Sitemap section (`markdown.frontmatter`, `markdown.sitemap-section`)

The `.md` mirror is emitted by `mkdocs-llmstxt` (0.3.0). The scorer wants each mirror to carry YAML frontmatter with **`title`, `description`, `doc_version`, `last_updated`** and a **`## Sitemap`** section. Inspect how your build currently writes the per-page `.md` (the llmstxt plugin config, or a `hooks:` script). Add the frontmatter (title/description from page meta, `doc_version` from the docs/release version, `last_updated` from the git date) and append a short `## Sitemap` section (a few links to sibling/section pages, or a link to `/docs/sitemap.md`). If the plugin can't inject this directly, a small MkDocs `on_post_build`/`on_page_content` hook over the emitted mirrors is the clean mechanism.

### D. Glossary (`html.glossary-link`)

Add a terminology/glossary page (e.g. `docs/site/reference/glossary.md`) defining the core terms (workflow, entity, stage, gate, first-officer, ensign, mod, standing teammate, split-root) and link to it from a site-wide location so EVERY page satisfies the check — the cleanest is a link in `overrides/partials/footer.html` (or a nav entry under Reference). The check looks for a link whose text/href signals "glossary"/"terminology".

### E. Code-fence languages (`code.language-tags`)

Fenced code blocks render `<code>` without a `language-*` class when the source omits the language hint. Sweep `docs/site/**/*.md` for bare ```` ``` ```` fences and add the right language (`bash`, `yaml`, `text`, `console`, etc.). This is per-source-content but mechanical.

---

## Explicitly NOT in this PR (handled in the landing repo)

- `markdown.canonical-header` — a `Link: <…>; rel="canonical"` HTTP header on `.md` responses → lives in the landing repo's `netlify.toml`.
- `markdown.content-negotiation` — serve markdown for `Accept: text/markdown` → landing-repo CDN/edge concern (the landing page fails this too; it's a known hard item).
- `html.text-ratio` (14.1%) — driven by Material's nav/search chrome inflating the HTML. Structural to the theme; the `.md` mirror is the real agent-readable path, so treat as accept-skip unless an easy chrome trim helps. Do not chase it at the cost of the human-facing docs UX.

## Definition of done

- `overrides/main.html` injects OG + JSON-LD (TechArticle + BreadcrumbList + dateModified) + the markdown alternate-link on every page; the JSON-LD parses.
- git-date plugin enabled and feeding `dateModified`/`last_updated`.
- `.md` mirrors carry the required frontmatter + a `## Sitemap` section.
- A glossary page exists and is linked site-wide.
- Source code fences declare languages.
- `mkdocs build --strict` is green.
- `a14y check <built-or-preview docs page> --mode page` shows checks A–E cleared (expect the per-page docs score to jump from ~68 into the 90s, matching the landing).
- Changes on a branch, PR opened, NOT pushed to `main`.
