---
id: wvyqyybd2vvknehb1a8ak9kr
title: MkDocs Material docs site + GitHub Pages publish
status: ideation
source: "captain (2026-06-06) - install-journey should be part of a complete public-facing docs site; organize docs so they actually build the site. Fast-follow to nb (readme-main-flip-reconciliation), which ships the reader-first README + reworked install content."
started: 2026-06-09T02:55:57Z
completed:
verdict:
score:
worktree:
issue:
---

Stand up a public-facing documentation site for Spacedock using **MkDocs + Material**, consuming the user-facing docs reorganized by `readme-main-flip-reconciliation` (nb). nb ships the reader-first `README.md` and a clean install guide (`docs/install-journey.md`); this task makes those (and the other user-facing docs) build into a published site.

## Problem

The repo's `docs/` is a flat, mixed-audience pile: user-facing docs (`install-journey.md`, `releasing.md`, `runtime-support.md`) sit next to dev/process docs (`docs/dev/`, `docs/specs/`, `docs/roadmap/`). There is no site generator, no information architecture, and no published site. Users have no public-facing docs home — the README carries everything.

## Proposed approach

Ideation fills this in. Sketch:

- Add MkDocs + Material config (`mkdocs.yml`) with an information architecture that includes ONLY user-facing pages (Home, Install, Usage, Concepts/Workflows, ...) and excludes dev/process docs.
- Reuse nb's reader-first README as the site home and the reworked install guide as the Install page.
- A GitHub Pages publish workflow (build the site, deploy to Pages on push to the stable branch).
- Decide the docs content root: MkDocs defaults to `docs/`, which currently mixes audiences — choose nav-include vs. a dedicated site content dir so dev/process docs stay out of the public site.

## Out of scope

- README/install CONTENT rewrites (owned by nb).
- The main-flip release mechanics and branch transition.

## Acceptance criteria

Ideation fills these in. Sketch:

- `mkdocs build --strict` produces a site with no broken links or unresolved refs.
- The published site (or a CI artifact of the build) renders Home + Install from the reader-first docs.
- Dev/process docs are excluded from the public navigation.

## Test plan

Ideation fills this in. `mkdocs build --strict` in CI; a link check; a Pages deploy smoke.
</content>
</invoke>
