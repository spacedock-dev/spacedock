# ABOUTME: MkDocs hook that enriches the per-page .md mirrors emitted by
# ABOUTME: mkdocs-llmstxt with YAML frontmatter and a ## Sitemap section.
"""Enrich the llmstxt per-page markdown mirrors for agent-readability.

mkdocs-llmstxt writes a markdown mirror of every page at
``<site_dir>/<page.url>index.md`` in its own ``on_post_build``. Those mirrors
carry no frontmatter and no sitemap section, and the plugin offers no config to
inject them. This hook collects each page's metadata during the build, then in
its own ``on_post_build`` (which MkDocs runs after the plugin's, because hooks
register after the ``plugins`` list) prepends a YAML frontmatter block
(``title``, ``description``, ``doc_version``, ``last_updated``) and appends a
short ``## Sitemap`` section of sibling-page links.

Clears the agent-readability ``markdown.frontmatter`` and
``markdown.sitemap-section`` checks.
"""

import json
import os
from pathlib import Path

# Per-page metadata gathered during the build, keyed by the page's output URL
# (e.g. "concepts/workflows-and-entities/", or "" for the home page). Populated
# in on_page_content; consumed in on_post_build.
_PAGES = {}

# Site-level config captured in on_config for use in on_post_build.
_CONFIG = {}


def on_config(config, **kwargs):
    _CONFIG["site_description"] = config.get("site_description", "")
    _CONFIG["doc_version"] = (config.get("extra", {}) or {}).get("doc_version", "")
    _CONFIG["site_dir"] = config["site_dir"]
    return config


def on_page_content(html, page, config, files, **kwargs):
    """Capture each page's frontmatter inputs and its sibling links.

    Runs after the page is rendered, so page.meta (including the git-revision
    date) and page.parent are fully populated.
    """
    meta = page.meta or {}

    title = (page.title or "").strip()
    description = (meta.get("description") or _CONFIG.get("site_description") or "").strip()

    last_updated = (
        meta.get("git_revision_date_localized_raw_iso_datetime")
        or meta.get("git_revision_date_localized_raw_datetime")
        or ""
    )

    # Sibling pages within the same nav section, linked to their .md mirrors.
    siblings = []
    parent = page.parent
    if parent is not None and getattr(parent, "children", None):
        for child in parent.children:
            if getattr(child, "is_page", False) and child.url and child.url != page.url:
                child_title = (child.title or "").strip()
                # Mirror lives at <child.url>index.md, relative to this mirror.
                href = _relative_mirror_href(page.url, child.url)
                siblings.append((child_title, href))

    _PAGES[page.url] = {
        "title": title,
        "description": description,
        "last_updated": last_updated,
        "siblings": siblings,
    }
    return html


def on_post_build(config, **kwargs):
    """Rewrite each emitted .md mirror with frontmatter + a Sitemap section."""
    site_dir = Path(_CONFIG["site_dir"])
    doc_version = _CONFIG.get("doc_version", "")

    for url, data in _PAGES.items():
        mirror = site_dir / url / "index.md"
        if not mirror.is_file():
            continue

        body = mirror.read_text(encoding="utf-8")

        frontmatter = _build_frontmatter(
            title=data["title"],
            description=data["description"],
            doc_version=doc_version,
            last_updated=data["last_updated"],
        )
        sitemap = _build_sitemap(data["siblings"])

        mirror.write_text(frontmatter + body.rstrip("\n") + "\n" + sitemap, encoding="utf-8")


def _build_frontmatter(title, description, doc_version, last_updated):
    """A YAML frontmatter block with the four agent-readability keys.

    Values are JSON-encoded so colons, quotes, and unicode stay valid YAML
    (JSON is a YAML subset for scalar strings).
    """
    lines = ["---"]
    lines.append("title: " + json.dumps(title, ensure_ascii=False))
    lines.append("description: " + json.dumps(description, ensure_ascii=False))
    lines.append("doc_version: " + json.dumps(str(doc_version), ensure_ascii=False))
    lines.append("last_updated: " + json.dumps(str(last_updated), ensure_ascii=False))
    lines.append("---")
    lines.append("")
    return "\n".join(lines) + "\n"


def _build_sitemap(siblings):
    """A ## Sitemap section linking sibling pages' markdown mirrors."""
    lines = ["", "## Sitemap", ""]
    if siblings:
        for title, href in siblings:
            label = title or href
            lines.append(f"- [{label}]({href})")
    else:
        lines.append("- [Documentation home](/docs/)")
    lines.append("")
    return "\n".join(lines)


def _relative_mirror_href(from_url, to_url):
    """Relative href from one page's mirror to another page's .md mirror.

    Both mirrors sit at <url>index.md. From "a/b/" to "c/d/" the relative path
    from a/b/index.md to c/d/index.md is computed against the page directories.
    """
    from_dir = from_url.rstrip("/")
    to_dir = to_url.rstrip("/")
    target = (to_dir + "/index.md") if to_dir else "index.md"
    rel = os.path.relpath(target, start=from_dir if from_dir else ".")
    return rel
