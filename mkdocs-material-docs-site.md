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
│   ├─ Debrief & refit                         source: debrief skill + refit skill                 OP
│   └─ ▶ Example: a real workflow  [e6 SLOT]   source: e6 (readme-real-workflow-example-link)      OP/NU
│         (the post-flip-stable link to one concrete in-repo workflow plugs in HERE)
│
├─ Advanced topics                                                                               audience: OP/CO
│   ├─ Sprints & roadmap                       source: docs/roadmap/README                         OP/CO
│   │     (FO↔Commander operating model; the sprint lifecycle checklist; roadmap-as-strategy-
│   │      layer; the convention-only dry run — prose + frontmatter + native --where, NO binary)
│   ├─ Mods & standing teammates              source: first-officer skill (Mod Hook Convention +  OP/CO
│   │     Standing Teammates) + docs/dev/_mods/
│   │     (startup/idle/merge lifecycle hooks via status --boot; the comm-officer prose-polisher)
│   ├─ External-tracker bridge                source: docs/specs/state-behavior-extension.md      OP/CO
│   │     (issue/source frontmatter to kata/Linear/GitHub Issues; one-way unless declared)
│   └─ Multi-workflow & split-root state      source: docs/specs/ + docs/dev/README state model   OP/CO
│         (the .spacedock-state checkout; definition_dir vs state_dir; concurrency-safe commits)
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
    ├─ Voice & tone                            source: README voice + comm-officer mod prose       CO
    │     discipline (docs/dev/_mods/comm-officer.md)
    │     (the writing-style guide for public docs; the comm-officer defers to THIS page)
    └─ Architecture notes                      source: docs/specs/ + AGENTS.md                     CO
          (split-root state; scenario-testing principles; project shape)
```

Notes on the IA:

- **Audience spine.** "Get started" is the new-user funnel (NU); "Concepts" + "Running workflows" + "Reference" serve the operator/captain (OP); "Advanced topics" is the power-user layer for constructs built ON TOP of the base primitives (OP/CO); "Contributing" serves the contributor (CO). A reader self-selects from the top-level nav.
- **Running workflows vs. Advanced topics.** "Running workflows" stays about operating a SINGLE workflow (commission, survey, operate, debrief/refit). "Advanced topics" is the scale/extensibility layer: multi-sprint strategy (roadmap/Commander), lifecycle mods + standing teammates, the external-tracker bridge, and multi-workflow/split-root state — patterns built on the base primitives, not the primitives themselves. A construct too thin for its own page folds into a sibling rather than padding the nav.
- **The e6 slot is explicit.** "Running workflows → Example: a real workflow" is the named slot where `e6` (readme-real-workflow-example-link) plugs in. e6's own AC already says its link is checked "or the docs-site build, once `mkdocs-material-docs-site` lands" — so once both land, `mkdocs build --strict` is the link-check e6 relies on, and the slot is the page that link lands on. Until e6 lands, the page is a placeholder stub (a one-line "example coming" page) so the nav entry exists and `--strict` is happy; e6 fills it.
- **Excluded from the public site (by omission from nav).** The active workflow STATE (`docs/dev/.spacedock-state/`), the per-sprint dispatch packages, internal proposals/mods, and the FO-internal skills' *implementation detail* (`present-gate`/`feedback-rejection-flow`/`using-claude-team` are referenced conceptually under Gates but not published as standalone how-to pages). These are dev/process artifacts, not public docs.
- **Content reuse, not duplication.** Pages whose source is a single existing doc (Install, Releasing, Multi-host, the dev workflow) are that doc, moved/linked under the site root. Pages synthesized from multiple sources (the operating model, first-launch, gates) are new short pages that link out — kept thin to avoid drift with the README/skills they summarize.
- **The Voice & tone page is drafted here, not stubbed.** Its full content is authored in the "Voice & tone (drafted page content)" section below; implementation lifts it into `docs/site/contributing/voice-and-tone.md` largely verbatim. It is grounded in the project's real voice (the README) and the comm-officer mod's prose discipline, and it BECOMES the project voice guide the comm-officer defers to.

## Voice & tone (drafted page content)

This is the actual content for the **Contributing → Voice & tone** page (above), drafted in ideation rather than stubbed. Implementation lifts it into `docs/site/contributing/voice-and-tone.md`. It is grounded in two real voice signals, not invented:

- **The root `README.md`** — Spacedock's actual public voice: direct, anti-hype, claim-first ("Spacedock is a multi-agent orchestrator where nothing ships without a decision"), concrete over abstract, second person addressed to the reader.
- **The `comm-officer` mod** (`docs/dev/_mods/comm-officer.md`) — the project's prose discipline. It applies the `elements-of-style:writing-clearly-and-concisely` skill (Strunk) and is **light-touch by default**: "Preserve the caller's voice, rhythm, and technical vocabulary. Cut empty words, tighten sentences, fix clear grammar errors." Critically, it **defers to a project voice guide when one exists** (comm-officer.md: "If a voice guide applies to this project … load it on first use and defer to it when it conflicts with Strunk"). This page IS that guide — so the comm-officer and any doc-writer apply it.

### Voice

- **Precise, honest, technical.** State what is true and what a thing does. Spacedock's own pitch leads with the claim, not the adjective.
- **No marketing or hype adjectives.** Avoid "powerful", "seamless", "revolutionary", "effortless", "blazing-fast", "game-changing". If a sentence still carries its meaning with the adjective removed, remove it. The product ethos — evidence over assertion — is the writing ethos.
- **Concrete over abstract.** Name the command, the file, the outcome. Prefer "`spacedock status --next` lists the items ready to dispatch" over "Spacedock surfaces actionable work."
- **Claim-first, then support.** Lead a section with the load-bearing sentence; follow with the detail. Mirrors the README's "what's different" bullets (bold claim, then the mechanism).

### Tone and register per audience

- **New-user pages (Get started): welcoming and encouraging, still precise.** Assume no prior Spacedock knowledge; define a term on first use; show the command and the output to expect. Confidence-building, never breezy.
- **Operator pages (Concepts, Running workflows): direct and operational.** The reader is doing the work; tell them the steps and the decision points plainly.
- **Contributor pages (Contributing, Reference): exact and unembellished.** Precision outranks warmth. Name the contract, the test, the failure mode.
- **Person and tense.** Second person ("you run", "you approve") and present tense for how-to and instructions — the README's register. Imperative for steps ("Run `spacedock doctor`."). Describe the system in the present tense ("the first officer dispatches an ensign"), not the future.

### Canonical terminology and capitalization

Pinned from how the README and skills actually use these forms — not an imposed guess. Use them consistently.

| Term | Form | Notes (grounded in real usage) |
|------|------|--------------------------------|
| Spacedock | `Spacedock` | The product. Always capitalized. `spacedock` (lowercase, code font) only as the literal command/binary. |
| Captain | `Captain` (role) / `the captain` (prose) | README roles table capitalizes the role name; running prose uses lowercase "the captain" (skills use `{captain}`). The human operator. |
| First Officer | `First Officer` (role) / `the first officer` (prose) | Same pattern: Title Case in the roles table / when naming the role; lowercase in running prose ("the first officer reads the README"). The orchestrator agent. |
| Ensign | `Ensign` (role) / `the ensign`, `ensigns` (prose) | Same pattern. The worker agent that moves one item through one stage. |
| workflow | `workflow` | Common noun, lowercase. A directory of markdown entities + a README. |
| entity | `entity` | Common noun, lowercase. One work item (a markdown file or folder). The README also says "work item" — prefer "entity" in docs, gloss it as "work item" on first use for new users. |
| stage | `stage` | Common noun, lowercase. backlog → ideation → implementation → validation → done. |
| gate | `gate` | Common noun, lowercase. The decision point at the end of a stage. |
| sprint | `sprint` | Common noun, lowercase. A grouped set of entities driven to a deliverable. |
| worktree, mod, safehouse | lowercase | Common nouns. `safehouse` is also the sandbox profile filename `.safehouse`. |

Rule of thumb: capitalize a **role** when you name it as a role (the roles table, a definition); use lowercase for the same word in ordinary running prose; never capitalize the common-noun primitives (workflow, entity, stage, gate, sprint).

### Formatting conventions

- **Commands and code.** Inline code font for commands, flags, filenames, and identifiers: `spacedock claude`, `--strict`, `mkdocs.yml`. Multi-line commands and config in fenced code blocks with a language tag (` ```bash `, ` ```yaml `).
- **Show expected output.** For a command a reader runs, name the result, as `docs/install-journey.md` does ("Prints the installed version, e.g. `spacedock 0.20.0`").
- **Headings.** Sentence case ("Get started", "Your first workflow"), not Title Case. One `#` h1 per page (the page title); section headings start at `##`.
- **Links.** Descriptive link text, never "click here" / "this link". Internal links are relative so `mkdocs build --strict` can resolve them.
- **Lists and emphasis.** Bullets for parallel items; bold for the load-bearing claim of a bullet (the README pattern). Use emphasis sparingly — if everything is bold, nothing is.

### How this guide is applied

The `comm-officer` mod loads a project voice guide on first use and defers to it over plain Strunk. Once this page ships, it is that guide: doc contributions and comm-officer polish follow these rules. When this guide is silent on a question, fall back to Strunk via `elements-of-style:writing-clearly-and-concisely`.

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

**AC-2 — The rendered site's navigation matches the approved outline: the top-level sections (Home, Get started, Concepts, Running workflows, Advanced topics, Reference, Contributing) and their pages, exactly as gated.**
Verified by: the `nav:` block in the committed `mkdocs.yml` is the build's nav source, and `mkdocs build --strict` fails if any `nav:` entry points at a missing file — so a nav that drifts from the approved tree (a renamed/missing page) fails the build. Approval of the IA at the gate fixes the expected `nav:` shape; the strict build enforces that every entry resolves. (The check's independent source is the set of files on disk vs. the `nav:` declaration — they can diverge, which is what makes it able to fail.)

**AC-3 — The site renders Home and the Get-started pages from the reader-first docs (README + install-journey), and the published/CI-artifact build contains those rendered HTML pages.**
Verified by: after `mkdocs build --strict`, the expected output files exist under `site/` (e.g. `site/index.html`, `site/get-started/install/index.html`); checked by file existence on the build output, not by grepping page prose.

**AC-4 — Dev/process docs are absent from the public site.**
Verified by: `docs_dir` is `docs/site/`, so `docs/dev/`, `docs/specs/`, and `docs/roadmap/` are not part of the build input; confirmed by their absence from the `site/` output tree after a build. (On-disk output state, not a prose check.)

**AC-5 — The e6 example slot exists in the nav as a resolvable page (a stub until e6 fills it).**
Verified by: the "Running workflows → Example: a real workflow" `nav:` entry resolves to a committed file, so `mkdocs build --strict` passes with the slot present; e6 later replaces the stub's content and its link resolves under the same strict build.

**AC-6 — A GitHub Pages publish workflow builds with `--strict` and deploys on push to the stable branch, and the strict build also gates PRs.**
Verified by: the workflow file under `.github/workflows/` exists and its build step is `mkdocs build --strict`; the deploy job targets GitHub Pages. (Workflow correctness is proven by the strict build passing in CI; the deploy smoke — the Pages URL serving the built site — is the implementation/validation live check, recorded once the workflow runs.)

**AC-7 — A "Contributing → Voice & tone" page exists in the nav and renders under `mkdocs build --strict`.**
Verified by: the `nav:` "Voice & tone" entry resolves to a committed file (`docs/site/contributing/voice-and-tone.md`) and `mkdocs build --strict` passes with it present (exit 0; the rendered `site/contributing/voice-and-tone/index.html` exists) — checked by the build and the output file, not a prose-grep of the page. The page's content is the guide drafted in the "Voice & tone (drafted page content)" section of this entity, lifted largely verbatim. Linkage: this page becomes the project voice guide the `comm-officer` mod (`docs/dev/_mods/comm-officer.md`) defers to over plain Strunk; that the comm-officer obeys it is a behavioral property of the mod, not this build's claim — the build's claim is only that the page exists and renders strictly.

## Test plan

- **Primary gate — the strict build (fixture/CLI-level, low cost).** `mkdocs build --strict` is the single load-bearing check: it proves the nav resolves, internal links resolve, and references resolve. It runs locally (after `pip install -r docs/requirements.txt`) and in CI on every PR. This is the check that can fail; it is not a prose-grep.
- **Output-tree assertions (on-disk state).** After a strict build, assert the expected `site/` files exist (Home, the Get-started pages, the e6 slot, the Voice & tone page) and that no `docs/dev` / `docs/specs` / `docs/roadmap` pages appear in `site/`. These check rendered output, not source prose.
- **Pages deploy smoke (live, implementation/validation).** Once the workflow runs on the stable branch, confirm the Pages URL serves the built site. This is the only step that needs CI to run; no host/agent launch is involved.
- **No live host run needed.** Building docs does not launch `spacedock claude/codex/pi`. The riskiest mechanism here is `mkdocs build --strict` resolving the `nav:` against `docs/site/` and Material rendering the section grouping — a standard, well-proven MkDocs path, so **no spike needed: the mechanism is MkDocs' documented `nav:` + `docs_dir` + `--strict` behavior plus the Material theme, all already-proven tooling.** The one declared dependency a fresh setup needs is `mkdocs` + `mkdocs-material` (pinned in `docs/requirements.txt`); they are not installed in this ideation environment, which is expected — implementation installs them and CI pins them. No hidden machine dependency: the requirements file is the declaration.
- **Cost: low.** Config (`mkdocs.yml`, `docs/requirements.txt`) + content scaffolding under `docs/site/` (mostly moves/links of existing docs + a few thin synthesized pages) + one CI workflow. No binary/Go changes.

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

### Captain IA fold

- DONE: Add a new top-level nav section `Advanced topics` after "Running workflows", before "Reference" (audience OP/CO power-user).
  New nav block; audience-spine note updated to name it the on-top/power-user layer.
- DONE: Relocate "Running workflows → Sprints & roadmap" INTO Advanced topics and expand it.
  Moved out of Running workflows; expanded to FO↔Commander operating model + sprint lifecycle + roadmap-as-strategy-layer + the convention-only dry run (prose + frontmatter + native `--where`, NO binary). Source: docs/roadmap/README.md.
- DONE: Populate Advanced topics with the other real on-top constructs (survey, don't invent).
  Added Mods & standing teammates (first-officer Mod Hook Convention §262 + Standing Teammates §283 + docs/dev/_mods/ — comm-officer.md, pr-merge.md; startup/idle/merge hooks via `status --boot`, the comm-officer prose-polisher); External-tracker bridge (state-behavior-extension.md — `issue`/`source` to kata/Linear/GitHub Issues); Multi-workflow & split-root state (docs/specs/ + docs/dev/README state model — `.spacedock-state`, definition_dir vs state_dir). Each tagged source + audience.
- DONE: Update IA notes (audience spine) and AC-2 (top-level list includes Advanced topics).
  Added a "Running workflows vs. Advanced topics" note clarifying single-workflow vs scale/extensibility; AC-2's section list now reads Home, Get started, Concepts, Running workflows, Advanced topics, Reference, Contributing.
- DONE: Clean the stray `</content>` / `</invoke>` tool-fragment between Test plan and Stage Report.
  Removed both lines.

The tree is now seven top-level sections. Advanced topics groups the four on-top constructs that exist in the repo today; no construct was invented. The strict-build ACs are unchanged in kind — the new section's pages resolve under the same `mkdocs build --strict` gate.

- DONE: Add a "Contributing → Voice & tone" page — drafted in ideation (not stubbed), grounded in real voice signals, with the comm-officer-defers-to-it linkage and an AC.
  New nav entry under Contributing; full page content authored in the "Voice & tone (drafted page content)" section (voice, tone/register per audience, canonical terminology+capitalization pinned from real README/skills usage, formatting conventions, the application note); AC-7 verifies it builds under `mkdocs build --strict` (not a prose-grep); grounded in README voice + comm-officer.md (which applies elements-of-style/Strunk and defers to a project voice guide — this page IS that guide).

The tree is now seven top-level sections with Voice & tone as a Contributing sub-page (page content drafted in-body for implementation to lift). Voice signals surveyed, not invented: the README's actual register and the comm-officer mod's prose discipline (comm-officer.md line 148 — defers to a project voice guide; line 140 — light-touch, preserve voice).
