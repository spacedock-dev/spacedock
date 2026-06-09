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

Stand up a public-facing documentation site for Spacedock using **MkDocs + Material**, published to GitHub Pages. The headline deliverable of this ideation is the **information architecture (the nav tree)** — concrete and complete enough for the captain to approve the IA before any site is built. The mkdocs-material setup sketch and the acceptance criteria support that outline. This is ideation: no site is built and no host is launched here.

The site serves three audiences, not just new users: a **new user** (install → first launch), an **operator/captain** (running workflows, the stage lifecycle, gates, the FO↔Commander model), and a **contributor** (the workflow, the contract, proof policy, releasing, the multi-host runtime). The outline below names, for every section, its source material in the repo and its primary audience, so the IA serves what Spacedock actually is rather than a generic template.

## Problem

The repo's documentation is real but has no public home and no information architecture. The reader-first `README.md` carries the whole pitch + install + quick start; the install walkthrough lives in `docs/install-journey.md`; the operating model and conventions are spread across nine skills under `skills/`, the workflow READMEs (`docs/dev/README.md`, `docs/roadmap/README.md`), and the design/process docs (`docs/releasing.md`, `docs/runtime-support.md`, `docs/specs/`). A newcomer can install, but there is no navigable site that takes a user from "what is this" to "running my first workflow", an operator to "how gates and the lifecycle work", or a contributor to "the contract and how to add a runtime". There is no site generator, no nav tree, and nothing published.

This task does not rewrite that content (nb owns the README/install rewrites). It gives the content a **structure that builds into a published site**, with a nav tree grounded in the real surfaces.

## Survey of the real surfaces (what the IA is grounded in)

Every nav section below cites one of these. Read before outlining; cited inline so the IA is auditable.

- **Root `README.md`** — the pitch ("nothing ships without a decision"), what's different, install (brew / curl|sh / source), quick start (`/spacedock:commission`), the three roles table (Captain / First Officer / Ensign), how-it-works, the usage grammar.
- **`docs/install-journey.md`** — the full first-run walkthrough: Homebrew, Linux `curl | sh`, Codex/Pi, build-from-source, `spacedock doctor`, the command grammar (`spacedock claude "task" [--safehouse…] [-- host-flags…]`).
- **`docs/runtime-support.md`** — the multi-host runtime model: runtime layers (skill adapters → dispatch host mode → registries → launch/install UX → live runner), the acceptance checklist, the "manifesting from void" operating prompt, the Pi live-smoke mechanism.
- **`docs/releasing.md`** — cutting a stable release from `main`, what the tag push fires (goreleaser, Homebrew tap, manifest stamps, marketplace), dev-only `next` publishing.
- **`docs/dev/README.md`** — the canonical workflow README: schema/frontmatter, the stage lifecycle (backlog → ideation → implementation → validation → done), gate semantics, the detached adversarial audit, the runtime live-CI model, the proof policy (no prose-grep ACs). The single best contributor-facing artifact in the repo.
- **`docs/roadmap/README.md`** — the FO↔Commander operating model, the sprint lifecycle checklist, the sprint sequence, how sprint membership works via `--where` queries.
- **`AGENTS.md`** — repo build instructions, project shape (`cmd/`, `internal/cli`, `internal/status`), state-branch bootstrap rules, skill-development rules, the expected `go test ./...` gate.
- **`skills/` (nine skills)** — the operator/contributor surface:
  - `first-officer` — runs/resumes a workflow, dispatches workers, manages gates, advances state.
  - `ensign` — executes one stage as a dispatched worker.
  - `commission` — interactively design and generate a workflow (`/spacedock:commission`).
  - `survey` — read a brownfield project's agent history, report the implicit workflow + open decisions, offer to commission (`/spacedock:survey`).
  - `debrief` — capture a session into a record the next session resumes from (`/spacedock:debrief`).
  - `refit` — upgrade an existing workflow's scaffolding to the current version (`/spacedock:refit`).
  - `present-gate` — the captain-facing gate-review rendering (FO-internal).
  - `feedback-rejection-flow` — rejection routing back to the feedback-to stage (FO-internal).
  - `using-claude-team` — the Claude Code team-harness discipline (FO-internal).
- **The command surface** (`internal/cli/help.go`, grouped help) — Launch: `claude` / `codex` / `pi`; Setup: `install` / `doctor`; Workflow: `status` / `new` / `state init` / `completion` / `dispatch`; plus `--version`.
- **`.claude-plugin/` and `.codex-plugin/`** (`plugin.json`, `marketplace.json`) — the plugin manifest, `requires-contract: ">=1,<2"`, the marketplace entry. Grounds the install/compat and contributor-release sections.
- **`docs/specs/`** — `state-behavior-extension.md` (split-root state semantics, external-tracker bridge) and `scenario-testing-principles.md` (cross-runtime scenario model). Contributor reference.

## Proposed approach — the OUTLINE (information architecture / nav tree)

The headline deliverable. The site is **content-only**: every page is an existing or lightly-adapted Markdown file under a dedicated site content root (`docs/site/`, see the setup sketch), so dev/process docs are excluded by *what is in the nav*, not by exclusion globs. Each section names its **source** and **primary audience** (NU = new user, OP = operator/captain, CO = contributor).

```
Home                                         source: README.md (pitch/what's-different)         audience: NU
│
├─ Get started                                                                                  audience: NU
│   ├─ Install                                source: docs/install-journey.md                    NU
│   │     (Homebrew · Linux curl|sh · Codex/Pi · build-from-source · doctor · command grammar)
│   ├─ Your first launch                      source: README "Install"/"Quick start" + survey    NU
│   │     (spacedock claude "/spacedock:survey"; what the front-door banner means)
│   └─ Your first workflow                    source: README "Quick start" + commission skill    NU
│         (/spacedock:commission; the design/review gates)
│
├─ Concepts                                                                                      audience: OP
│   ├─ The operating model                    source: README "How it works" + roadmap README     OP
│   │     (Captain / First Officer / Ensign; the FO-shapes / Commander-drives split)
│   ├─ Workflows & entities                    source: docs/dev/README "Schema"/"File Naming"      OP
│   │     (a workflow = a dir of markdown + a README; entity frontmatter; the state checkout)
│   ├─ The stage lifecycle                     source: docs/dev/README "Stages"                    OP
│   │     (backlog → ideation → implementation → validation → done; worktree vs inline)
│   └─ Gates & decisions                       source: docs/dev/README + present-gate skill        OP
│         (gate review, approve/redo/reject, feedback cycles, the detached adversarial audit)
│
├─ Running workflows                                                                             audience: OP
│   ├─ Commission a workflow                   source: commission skill                            OP
│   ├─ Survey an existing project              source: survey skill                                OP
│   ├─ Operating a workflow                    source: first-officer skill + README "How it works" OP
│   │     (driving stages, dispatch, gates; spacedock status / --next / --where)
│   ├─ Sprints & roadmap                       source: docs/roadmap/README                         OP
│   │     (FO↔Commander, the sprint lifecycle, membership via --where)
│   ├─ Debrief & refit                         source: debrief skill + refit skill                 OP
│   └─ ▶ Example: a real workflow  [e6 SLOT]   source: e6 (readme-real-workflow-example-link)      OP/NU
│         (the post-flip-stable link to one concrete in-repo workflow plugs in HERE)
│
├─ Reference                                                                                     audience: OP/CO
│   ├─ Command reference                       source: internal/cli/help.go grouped help           OP/CO
│   │     (claude/codex/pi · install/doctor · status/new/state/completion/dispatch · --version)
│   └─ Multi-host support                       source: docs/runtime-support.md                     OP/CO
│         (Claude tier-1; Codex/Pi experimental; what "supported" means)
│
└─ Contributing                                                                                  audience: CO
    ├─ The development workflow               source: docs/dev/README                             CO
    │     (the workflow that builds Spacedock; the stage gates it runs under)
    ├─ Proof policy                            source: docs/dev/README (ideation/validation)       CO
    │     (no prose-grep ACs; behavior proven by exercising; the detached adversarial audit)
    ├─ Adding a runtime                        source: docs/runtime-support.md                     CO
    │     (the runtime layers; the "manifesting from void" prompt; the Pi live-smoke mechanism)
    ├─ Releasing                               source: docs/releasing.md                           CO
    │     (cut from main; what the tag fires; next is dev-only)
    └─ Architecture notes                      source: docs/specs/ + AGENTS.md                     CO
          (split-root state; scenario-testing principles; project shape)
```

Notes on the IA:

- **Audience spine.** "Get started" is the new-user funnel; "Concepts" + "Running workflows" + "Reference" serve the operator/captain; "Contributing" serves the contributor. A reader self-selects from the top-level nav.
- **The e6 slot is explicit.** "Running workflows → Example: a real workflow" is the named slot where `e6` (readme-real-workflow-example-link) plugs in. e6's own AC already says its link is checked "or the docs-site build, once `mkdocs-material-docs-site` lands" — so once both land, `mkdocs build --strict` is the link-check e6 relies on, and the slot is the page that link lands on. Until e6 lands, the page is a placeholder stub (a one-line "example coming" page) so the nav entry exists and `--strict` is happy; e6 fills it.
- **Excluded from the public site (by omission from nav).** The active workflow STATE (`docs/dev/.spacedock-state/`), the per-sprint dispatch packages, internal proposals/mods, and the FO-internal skills' *implementation detail* (`present-gate`/`feedback-rejection-flow`/`using-claude-team` are referenced conceptually under Gates but not published as standalone how-to pages). These are dev/process artifacts, not public docs.
- **Content reuse, not duplication.** Pages whose source is a single existing doc (Install, Releasing, Multi-host, the dev workflow) are that doc, moved/linked under the site root. Pages synthesized from multiple sources (the operating model, first-launch, gates) are new short pages that link out — kept thin to avoid drift with the README/skills they summarize.

## Setup sketch (NOT a build — implementation owns the full setup)

One section, since the captain wants the OUTLINE first. This sketches the shape; implementation builds it.

- **Generator + theme.** MkDocs with the Material theme (`mkdocs` + `mkdocs-material`, pinned in a `docs/requirements.txt`). Material gives the nav, search, and section grouping the IA above assumes.
- **Content root.** A dedicated `docs/site/` as MkDocs' `docs_dir` (set in `mkdocs.yml`), NOT the repo's mixed `docs/`. This keeps dev/process docs (`docs/dev/`, `docs/specs/`, `docs/roadmap/`) out of the site by construction — the site sees only `docs/site/`. Pages are real Markdown files placed (moved or symlinked from their canonical location) under `docs/site/` matching the nav tree.
- **`mkdocs.yml` shape.** `site_name`, `theme: { name: material }`, `docs_dir: docs/site`, and an explicit `nav:` that mirrors the approved outline above (top-level sections → their pages). `--strict` is the gate, so every nav entry must resolve to a file and every internal link must resolve.
- **Build + deploy.** A GitHub Pages publish workflow (`.github/workflows/`, alongside the existing `release.yml` / `runtime-live-e2e.yml`): on push to the stable branch (`main`, per `docs/releasing.md`), `pip install -r docs/requirements.txt`, `mkdocs build --strict`, then deploy the built `site/` to GitHub Pages (the `actions/deploy-pages` flow or `mkdocs gh-deploy`). The strict build also runs as a PR check so a broken nav/link fails before merge.

Open question for the gate (does not block IA approval): whether site pages are **moved** into `docs/site/` (single source of truth, the canonical doc location changes) or **symlinked/included** from their current paths (canonical location unchanged, site is a view). Recommendation: move the purely-user-facing docs (install-journey) under `docs/site/`; keep contributor docs (releasing, runtime-support, specs, dev README) canonical where they are and surface them via MkDocs include or a thin linking page, since AGENTS.md and the workflow reference them by their current paths.

## Out of scope

- README/install CONTENT rewrites — owned by nb (`readme-main-flip-reconciliation`).
- The main-flip release mechanics and branch transition.
- Writing the e6 example link itself (this task leaves the SLOT; e6 fills it post-flip).
- Building or deploying the site, or any live host run — this is ideation. Implementation builds; this approves the IA.
- Versioned docs (mike), API/reference auto-generation, or a docs search backend beyond Material's built-in.

## Acceptance criteria

Each AC names a property of the finished entity (the built site) and how it is verified. **No AC is satisfied by a prose-grep of page content** — the checkable behavior is the BUILT site, per the workflow's proof policy.

**AC-1 — `mkdocs build --strict` succeeds against the committed `mkdocs.yml` + `docs/site/`, producing a `site/` with no broken nav entries, internal links, or unresolved references.**
Verified by: running `mkdocs build --strict` (in CI and reproducibly locally after `pip install -r docs/requirements.txt`); exit code 0 and a populated `site/` directory. `--strict` turns any unresolved nav/link/reference into a non-zero exit, so this is a real failure surface, not a prose check.

**AC-2 — The rendered site's navigation matches the approved outline: the top-level sections (Home, Get started, Concepts, Running workflows, Reference, Contributing) and their pages, exactly as gated.**
Verified by: the `nav:` block in the committed `mkdocs.yml` is the build's nav source, and `mkdocs build --strict` fails if any `nav:` entry points at a missing file — so a nav that drifts from the approved tree (a renamed/missing page) fails the build. Approval of the IA at the gate fixes the expected `nav:` shape; the strict build enforces that every entry resolves. (The check's independent source is the set of files on disk vs. the `nav:` declaration — they can diverge, which is what makes it able to fail.)

**AC-3 — The site renders Home and the Get-started pages from the reader-first docs (README + install-journey), and the published/CI-artifact build contains those rendered HTML pages.**
Verified by: after `mkdocs build --strict`, the expected output files exist under `site/` (e.g. `site/index.html`, `site/get-started/install/index.html`); checked by file existence on the build output, not by grepping page prose.

**AC-4 — Dev/process docs are absent from the public site.**
Verified by: `docs_dir` is `docs/site/`, so `docs/dev/`, `docs/specs/`, and `docs/roadmap/` are not part of the build input; confirmed by their absence from the `site/` output tree after a build. (On-disk output state, not a prose check.)

**AC-5 — The e6 example slot exists in the nav as a resolvable page (a stub until e6 fills it).**
Verified by: the "Running workflows → Example: a real workflow" `nav:` entry resolves to a committed file, so `mkdocs build --strict` passes with the slot present; e6 later replaces the stub's content and its link resolves under the same strict build.

**AC-6 — A GitHub Pages publish workflow builds with `--strict` and deploys on push to the stable branch, and the strict build also gates PRs.**
Verified by: the workflow file under `.github/workflows/` exists and its build step is `mkdocs build --strict`; the deploy job targets GitHub Pages. (Workflow correctness is proven by the strict build passing in CI; the deploy smoke — the Pages URL serving the built site — is the implementation/validation live check, recorded once the workflow runs.)

## Test plan

- **Primary gate — the strict build (fixture/CLI-level, low cost).** `mkdocs build --strict` is the single load-bearing check: it proves the nav resolves, internal links resolve, and references resolve. It runs locally (after `pip install -r docs/requirements.txt`) and in CI on every PR. This is the check that can fail; it is not a prose-grep.
- **Output-tree assertions (on-disk state).** After a strict build, assert the expected `site/` files exist (Home, the Get-started pages, the e6 slot) and that no `docs/dev` / `docs/specs` / `docs/roadmap` pages appear in `site/`. These check rendered output, not source prose.
- **Pages deploy smoke (live, implementation/validation).** Once the workflow runs on the stable branch, confirm the Pages URL serves the built site. This is the only step that needs CI to run; no host/agent launch is involved.
- **No live host run needed.** Building docs does not launch `spacedock claude/codex/pi`. The riskiest mechanism here is `mkdocs build --strict` resolving the `nav:` against `docs/site/` and Material rendering the section grouping — a standard, well-proven MkDocs path, so **no spike needed: the mechanism is MkDocs' documented `nav:` + `docs_dir` + `--strict` behavior plus the Material theme, all already-proven tooling.** The one declared dependency a fresh setup needs is `mkdocs` + `mkdocs-material` (pinned in `docs/requirements.txt`); they are not installed in this ideation environment, which is expected — implementation installs them and CI pins them. No hidden machine dependency: the requirements file is the declaration.
- **Cost: low.** Config (`mkdocs.yml`, `docs/requirements.txt`) + content scaffolding under `docs/site/` (mostly moves/links of existing docs + a few thin synthesized pages) + one CI workflow. No binary/Go changes.
</content>
</invoke>

## Stage Report: ideation

- DONE: Produce the docs-site OUTLINE as the headline deliverable — the mkdocs-material nav tree (top-level sections + their pages), grounded in what Spacedock ACTUALLY is and does, with per-section source + audience.
  The nav tree under "Proposed approach — the OUTLINE" names six top-level sections (Home, Get started, Concepts, Running workflows, Reference, Contributing), each page tagged with its repo source and audience (NU/OP/CO).
- DONE: Survey the real surfaces FIRST and cite them, do not template (root README, the nine skills, the command surface, docs/, .claude-plugin/.codex-plugin).
  Read README.md, install-journey.md, runtime-support.md, releasing.md, docs/dev/README.md, docs/roadmap/README.md, AGENTS.md, all nine SKILL.md files, internal/cli/help.go (command surface), and both plugin manifests; captured in the "Survey of the real surfaces" section, each cited inline in the IA.
- DONE: For each nav section name its source material and its primary audience (new user / operator-captain / contributor).
  Every nav line carries `source:` + audience tag; the audience spine is summarized in the IA notes.
- DONE: Leave an explicit slot for e6 (readme-real-workflow-example-link) to plug into.
  "Running workflows → Example: a real workflow [e6 SLOT]"; AC-5 makes the slot a resolvable stub page; cross-checked against e6's own AC (its link-check IS this site's `--strict` build).
- DONE: Sketch (not build) the mkdocs-material setup — mkdocs.yml shape, theme + nav, where doc sources live (docs/site/), build + deploy via GitHub Pages CI.
  "Setup sketch" section: mkdocs + mkdocs-material pinned in docs/requirements.txt, docs_dir: docs/site, explicit nav: mirroring the outline, Pages workflow building with `--strict` on push to main + as a PR check.
- DONE: Give ACs + a test plan whose checkable behavior is the BUILT site, never a prose-grep of content; state cost (low) and that no live host run is needed.
  Six ACs (AC-1..AC-6) all verified by the strict build / output-tree state / workflow file, never by grepping page prose; test plan names `mkdocs build --strict` as the load-bearing check, records "no spike needed" (MkDocs is proven tooling), declares the mkdocs-material dependency in docs/requirements.txt, states cost low and no host launch.

### Summary

Delivered the headline IA: a six-section mkdocs-material nav tree grounded in a first-hand survey of the real surfaces (README, nine skills, command surface in internal/cli/help.go, docs/, plugin manifests), with every page tagged by source + audience across the three audiences (new user / operator-captain / contributor). Left an explicit, AC-backed e6 slot under "Running workflows" and confirmed e6's own link-check is this site's `--strict` build. Key decisions: a dedicated `docs/site/` content root (excludes dev/process docs by construction, not by globs), content reuse over duplication for thin synthesized pages, and `mkdocs build --strict` as the single load-bearing proof — no prose-grep ACs, no spike needed, no live host run. One open question recorded for the gate (move vs. symlink contributor docs into the site root); it does not block IA approval.
