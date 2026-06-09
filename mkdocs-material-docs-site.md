---
id: wvyqyybd2vvknehb1a8ak9kr
title: MkDocs Material docs site + GitHub Pages publish
status: validation
source: "captain (2026-06-06) - install-journey should be part of a complete public-facing docs site; organize docs so they actually build the site. Fast-follow to nb (readme-main-flip-reconciliation), which ships the reader-first README + reworked install content."
started: 2026-06-09T02:55:57Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-mkdocs-material-docs-site
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
- **A machine-readable surface alongside the human nav (on-brand).** Beyond the human nav tree, the build emits an agent-readable surface — an `llms.txt` LLM index (via `mkdocs-llmstxt`) and a reachable `AGENTS.md` — reached through a small "For agents" theme link, NOT a human nav section. This is deliberate for a docs site whose readers include a user's first officer (itself an AI agent); the agent-readability is then measured by dogfooding a14y against the deployed site (see Setup sketch + AC-8/AC-9/AC-10).
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
- **Content root.** A dedicated `docs/site/` as MkDocs' `docs_dir` (set in `mkdocs.yml`), NOT the repo's mixed `docs/`. This keeps dev/process docs (`docs/dev/`, `docs/specs/`, `docs/roadmap/`) out of the site by construction — the site sees only `docs/site/`. `docs/site/` is the single content root, populated two ways (see the decided note below): pure-content pages are MOVED in (their canonical home becomes `docs/site/`); functional files that the tool reads at fixed paths are SYMLINKED/included in (canonical path unchanged).
- **`mkdocs.yml` shape.** `site_name`, `theme: { name: material }`, `docs_dir: docs/site`, and an explicit `nav:` that mirrors the approved outline above (top-level sections → their pages). `--strict` is the gate, so every nav entry must resolve to a file and every internal link must resolve.
- **Build + deploy.** A GitHub Pages publish workflow (`.github/workflows/`, alongside the existing `release.yml` / `runtime-live-e2e.yml`): on push to the stable branch (`main`, per `docs/releasing.md`), `pip install -r docs/requirements.txt`, `mkdocs build --strict`, then deploy the built `site/` to GitHub Pages (the `actions/deploy-pages` flow or `mkdocs gh-deploy`). The strict build also runs as a PR check so a broken nav/link fails before merge.
- **Agent-readable docs (the on-brand differentiator).** Spacedock's docs are read by AI agents — a user's first officer parsing the docs IS an agent — so the site is built dual-audience: humans navigate via the nav; machines read a curated "for agents" surface. Borrows a14y.dev's pattern.
  - Add the `mkdocs-llmstxt` plugin (pinned in `docs/requirements.txt`, wired in `mkdocs.yml` `plugins:`) so `mkdocs build --strict` generates an **`llms.txt`** — a curated LLM index of the docs — into `site/`.
  - Surface **`AGENTS.md`** as a reachable machine artifact in the build output (a page that includes/links it, or a copy under `site/`), so an agent can fetch the repo's agent instructions from the docs site.
  - Add a small **"For agents"** footer/header link (a theme-level link, NOT a human nav page) pointing at `llms.txt` + `AGENTS.md`.
  - This closes the on-brand loop: docs for an agent tool, made agent-readable, measured by an agent-readability spec — Spacedock dogfooding its own thesis.
- **Material built-ins (low cost, mostly theme config).** Lean on the theme instead of custom work. In `mkdocs.yml` `theme:`:
  - `palette` with a **light/dark toggle** (two `palette` entries with a `scheme: default`/`slate` toggle).
  - `features:` including `content.code.copy` (copy-to-clipboard code blocks), `search.suggest` + `search.highlight`, `navigation.sections` + `navigation.instant`, and `content.tabs.link`.
  - Use Material **admonitions/callouts** and **content tabs** in pages where they earn their place — e.g. per-host install tabs (Homebrew / Linux / Codex/Pi) on the Install page, warning/note callouts in the runtime and releasing pages.

**Decided (captain): `docs/site/` is the single content root.** Two population rules, by whether the file is also consumed by tooling:

- **Pure-content docs MOVE into `docs/site/`** — `install-journey.md` and the synthesized pages (operating model, first-launch, gates, voice & tone, etc.). Their canonical home becomes `docs/site/`.
- **Functional files are SYMLINKED/included into `docs/site/`, NOT moved** — moving them would break the consumers that read them at fixed paths. These stay canonical where they are:
  - `docs/dev/README.md` — the LIVE workflow definition `spacedock status --workflow-dir docs/dev` parses; MUST stay at `docs/dev/`.
  - `AGENTS.md` — root-required (agents/build read it there).
  - `docs/releasing.md` — read by the release process; keep canonical.
  - `docs/specs/*` — referenced by the workflow README; keep canonical.

So the site is one content root; the functional contributor docs appear in it via symlink or a MkDocs include, with their canonical path unchanged. No functional file physically relocates.

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

**AC-4 — Process/state artifacts are absent from the public site; only the explicitly-surfaced docs appear.**
Verified by: `docs_dir` is `docs/site/`, so no source dir tree (`docs/dev/`, `docs/specs/`, `docs/roadmap/`) is itself part of the build input — the site sees only `docs/site/`. The active workflow state (`docs/dev/.spacedock-state/`), the roadmap, and the sprint packages never appear in the `site/` output tree. The functional contributor files that ARE published (e.g. `docs/dev/README.md`, selected `docs/specs/*`, `AGENTS.md`) appear ONLY because they are symlinked/included into `docs/site/` per the decided content-root rule — not because their original dir trees are in the build input, and not by physically relocating them. Confirmed by the `site/` output tree after a build: the surfaced pages present, the state/roadmap/sprint artifacts absent. (On-disk output state, not a prose check.)

**AC-5 — The e6 example slot exists in the nav as a resolvable page (a stub until e6 fills it).**
Verified by: the "Running workflows → Example: a real workflow" `nav:` entry resolves to a committed file, so `mkdocs build --strict` passes with the slot present; e6 later replaces the stub's content and its link resolves under the same strict build.

**AC-6 — A GitHub Pages publish workflow builds with `--strict` and deploys on push to the stable branch, and the strict build also gates PRs.**
Verified by: the workflow file under `.github/workflows/` exists and its build step is `mkdocs build --strict`; the deploy job targets GitHub Pages. (Workflow correctness is proven by the strict build passing in CI; the deploy smoke — the Pages URL serving the built site — is the implementation/validation live check, recorded once the workflow runs.)

**AC-7 — A "Contributing → Voice & tone" page exists in the nav and renders under `mkdocs build --strict`.**
Verified by: the `nav:` "Voice & tone" entry resolves to a committed file (`docs/site/contributing/voice-and-tone.md`) and `mkdocs build --strict` passes with it present (exit 0; the rendered `site/contributing/voice-and-tone/index.html` exists) — checked by the build and the output file, not a prose-grep of the page. The page's content is the guide drafted in the "Voice & tone (drafted page content)" section of this entity, lifted largely verbatim. Linkage: this page becomes the project voice guide the `comm-officer` mod (`docs/dev/_mods/comm-officer.md`) defers to over plain Strunk; that the comm-officer obeys it is a behavioral property of the mod, not this build's claim — the build's claim is only that the page exists and renders strictly.

**AC-8 — The strict build produces a populated `site/llms.txt` (the `mkdocs-llmstxt` plugin's output).**
Verified by: `mkdocs build --strict` (with `mkdocs-llmstxt` pinned in `docs/requirements.txt` and wired in `mkdocs.yml`) exits 0 and the file `site/llms.txt` exists and is non-empty — an on-disk build-output check, not a prose-grep of page content. The plugin failing to run, or producing no index, fails the assertion.

**AC-9 — The "for agents" machine surface (`llms.txt` + `AGENTS.md`) is reachable in the build output.**
Verified by: after the strict build, `site/llms.txt` exists (AC-8) and `AGENTS.md` is reachable from the built site (a `site/agents/…` page that includes/links it, or a copied `site/AGENTS.md`), and the "For agents" theme link resolves under `--strict`. Checked by build-output file existence + the strict link resolution, not by grepping prose.

**AC-10 — The agent-readability score from dogfooding a14y against the deployed docs is recorded.**
Verified by: at validation/implementation time, running `npx a14y <deployed-docs-url>` against the live Pages site and recording the returned score in the entity/PR (target: a high score; the a14y badge can be surfaced once scored). This needs the deployed URL, so it is an implementation/validation RECORDED METRIC, not a build-time gate — the cited evidence is the a14y command output, an external check that can fail, not the docs' own prose.

**AC-11 — The built site has a light/dark palette toggle and copy-to-clipboard code blocks enabled.**
Verified by: the committed `mkdocs.yml` `theme:` block (two `palette` entries with the scheme toggle; `features:` including `content.code.copy`) drives the build, and the rendered `site/` output contains the toggle control and the per-code-block copy affordance (Material emits these from the config). Checked by the build output produced from the config, not a prose-grep — a `theme:` that drops the toggle/feature produces output that fails the assertion.

## Test plan

- **Primary gate — the strict build (fixture/CLI-level, low cost).** `mkdocs build --strict` is the single load-bearing check: it proves the nav resolves, internal links resolve, and references resolve. It runs locally (after `pip install -r docs/requirements.txt`) and in CI on every PR. This is the check that can fail; it is not a prose-grep.
- **Output-tree assertions (on-disk state).** After a strict build, assert the expected `site/` files exist (Home, the Get-started pages, the e6 slot, the Voice & tone page, **`site/llms.txt`, and the reachable `AGENTS.md` surface**) and that the process/state artifacts are absent — no active workflow state (`docs/dev/.spacedock-state/`), roadmap, or sprint-package pages in `site/`. Only the explicitly-surfaced contributor files (e.g. `docs/dev/README.md`, selected `docs/specs/*`) appear, via the `docs/site/` symlink/include — not their whole source dir trees. Also assert the `theme:`-driven output carries the light/dark toggle and code-copy affordance. These check rendered output, not source prose.
- **Agent-readability dogfood (live, validation-recorded metric).** After deploy, run `npx a14y <deployed-docs-url>` and record the agent-readability score in the entity/PR. This is the on-brand loop — an agent-tool's docs measured by an agent-readability spec — and needs the live URL, so it is recorded at validation, NOT a build-time gate. The evidence is the a14y output, an external check that can fail.
- **Pages deploy smoke (live, implementation/validation).** Once the workflow runs on the stable branch, confirm the Pages URL serves the built site. This is the only step that needs CI to run; no host/agent launch is involved.
- **No live host run needed.** Building docs does not launch `spacedock claude/codex/pi`. The riskiest mechanism here is `mkdocs build --strict` resolving the `nav:` against `docs/site/`, Material rendering the section grouping + theme features, and `mkdocs-llmstxt` emitting `site/llms.txt` — all standard, well-proven MkDocs/plugin paths, so **no spike needed: the mechanism is MkDocs' documented `nav:` + `docs_dir` + `--strict` behavior, the Material theme's documented `palette`/`features`, and the `mkdocs-llmstxt` plugin's documented `site/llms.txt` output — already-proven tooling.** The declared dependencies a fresh setup needs are `mkdocs` + `mkdocs-material` + `mkdocs-llmstxt` (pinned in `docs/requirements.txt`); they are not installed in this ideation environment, which is expected — implementation installs them and CI pins them. The a14y dogfood needs `npx a14y` and the deployed URL at validation. No hidden machine dependency: the requirements file + the recorded a14y command are the declarations.
- **Cost: low.** Config (`mkdocs.yml`, `docs/requirements.txt`) + content scaffolding under `docs/site/` (mostly moves/links of existing docs + a few thin synthesized pages) + the `mkdocs-llmstxt` plugin and Material `theme:` config (no custom code) + one CI workflow. The a14y dogfood is a one-line validation step. No binary/Go changes.

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

- DONE: Fold #1 (agent-readable docs) + #5 (Material built-ins) from the a14y.dev review — setup sketch + IA notes + ACs.
  #1: setup sketch adds `mkdocs-llmstxt` → `site/llms.txt`, an `AGENTS.md` machine surface, and a small "For agents" theme link (not a human nav page); an IA note on the machine-readable surface alongside the human nav; the a14y dogfood as a validation-recorded metric. ACs: AC-8 (llms.txt produced), AC-9 (for-agents surface reachable), AC-10 (a14y score recorded against the deployed site — recorded metric, not build gate), all build/artifact-verified, no prose-grep. On-brand loop noted (agent-tool docs, agent-readable, measured by an agent-readability spec).
  #5: setup sketch adds Material `theme:` palette light/dark toggle, `features:` (content.code.copy, search.suggest/highlight, navigation.sections/instant, content.tabs.link), and admonitions/content-tabs in pages that benefit (per-host install tabs). AC-11 verifies the toggle + code-copy from the committed `theme:` config driving the build output, not a prose-grep.
  Test plan updated: output-tree assertions now include site/llms.txt + the AGENTS.md surface + the theme toggle/code-copy; the a14y dogfood added as a live validation-recorded step; dependencies extended to mkdocs-llmstxt (pinned) + npx a14y; cost remains low (config + plugin + theme, no custom code).

- DONE: Resolve the move-vs-symlink open question (captain decision) — record in the entity body, ensure no AC implies functional files relocate.
  Replaced the "Open question for the gate" paragraph with a decided note: `docs/site/` is the single content root; pure-content docs (install-journey + synthesized pages) MOVE in; functional files the tool reads at fixed paths (docs/dev/README.md — the live workflow def `spacedock status --workflow-dir docs/dev` parses; AGENTS.md — root-required; docs/releasing.md — release process; docs/specs/* — referenced by the workflow README) are SYMLINKED/included, canonical path unchanged. Tightened the Content-root setup bullet to match. Tightened AC-4 (and the matching test-plan output-tree assertion) so it asserts process/state artifacts absent + only explicitly-surfaced files present via the docs/site/ symlink — no AC now implies a functional file physically relocates.

## Deferred follow-ups (recorded at implementation)

- **Content-sourcing / no-duplication policy — DEFERRED (captain: "worry about it later").** No pymdownx-snippets/include machinery was built. The thin synthesized pages (Home, the operating model, gates, proof policy, etc.) carry light overlap with the README and the skills they summarize; this is acceptable for v1. When drift becomes a maintenance cost, revisit with a snippet/include mechanism or by collapsing a synthesized page into its single source.
- **`docs/specs/frontmatter-contract.md` standalone spec — FOLLOW-UP (captain, this turn).** This task delivered the `Reference → Frontmatter contract` nav slot as a page surfacing the dev README's "Schema / Field Reference" table. The standalone, workflow-neutral `docs/specs/frontmatter-contract.md` is a later follow-up, not this task.
- **The e6 example slot is a stub (AC-5).** `running-workflows/example.md` is an honest "example coming" placeholder so the nav entry resolves under `--strict`. `e6` (readme-real-workflow-example-link) fills it post-flip; once it lands, this same strict build is the link-check e6 relies on.
- **AC-10 (a14y dogfood score) is a validation-time recorded metric, not a build gate.** It needs the deployed Pages URL: run `npx a14y <deployed-docs-url>` after the Pages deploy and record the score. Not satisfiable at implementation time (no live URL yet).
- **AC-6 deploy smoke is a live check.** The workflow file exists and the `--strict` build step is in place; the Pages URL serving the built site is confirmed once `docs.yml` runs on `main` (the stable branch). Pages must be enabled for the repo with the "GitHub Actions" source for the deploy job to publish.

## Stage Report: implementation

- DONE: Stand up the mkdocs-material docs site so it builds GREEN under `mkdocs build --strict` — `mkdocs.yml` with the approved 7-section nav + theme config; `docs/site/` content root with pure-content MOVED in and functional files SYMLINKED; `docs/requirements.txt` pinning mkdocs + mkdocs-material + mkdocs-llmstxt; every nav entry and internal link resolves.
  `mkdocs build --strict` exits 0 (verified from a fresh `rm -rf site` build). install-journey.md `git mv`'d to docs/site/get-started/install.md; 5 functional files symlinked as git mode-120000 (docs/dev/README.md, AGENTS.md, docs/releasing.md, docs/specs/state-behavior-extension.md, docs/runtime-support.md). requirements pins mkdocs==1.6.1, mkdocs-material==9.6.21, mkdocs-llmstxt==0.3.0 (all installed clean). Worktree commit 14cab445.
- DONE: Wire the agent-readable surface + publish CI — mkdocs-llmstxt generates `site/llms.txt`; the "For agents" link + AGENTS.md reachable in the build (AC-8/AC-9); a GitHub Pages workflow builds with `--strict` as a PR check and deploys on push to the stable branch.
  `site/llms.txt` generated, 3178 bytes / 54 lines, sectioned index of all pages (AC-8). AGENTS.md reachable via symlinked `agents/index.html` (AC-9). "For agents" banner (docs/overrides/main.html) links llms.txt + agents/ with depth-correct `| url` paths. `.github/workflows/docs.yml`: build job runs `mkdocs build --strict` on PR + push to main; deploy job (`actions/deploy-pages@v4`) gated to `refs/heads/main` (AC-6).
- DONE: Author the net-new pages to a FIRST-VERSION bar — synthesized Concepts/Running/Advanced pages, the Voice & tone page (lifted verbatim from the ideation draft), the e6-slot stub, and a `Reference → Frontmatter contract` slot; nav complete and `--strict` green; content-sourcing policy DEFERRED, no snippet machinery, README↔site overlap recorded.
  All 7 sections populated; Voice & tone lifted from the "Voice & tone (drafted page content)" body section; example.md is an honest e6 stub; frontmatter-contract.md surfaces the dev README schema table per the docs/specs convention. No pymdownx-snippets machinery built; overlap + standalone-spec follow-ups recorded above.
- DONE: AC-11 theme features verified from build output.
  Palette toggle server-rendered (two `name="__palette"` inputs, `data-md-color-scheme` default+slate); `content.code.copy` present in the page feature-config JSON that drives Material's runtime copy-button injection. Both come from the committed `theme:` config.
- DONE: Keep the deliverable disjoint from live Commander 0200 work and the `go test ./...` baseline gate green.
  Touched only docs/site/, mkdocs.yml, docs/requirements.txt, docs/overrides/, .github/workflows/docs.yml, .gitignore, and the three move-induced fixes (README link, install.sh hint, install_doc_test.go path). Did NOT touch release.yml/frontdoor.go/.goreleaser.yaml. `go test ./...` green: 1247 passed in 16 packages.
- SKIPPED: AC-10 (a14y agent-readability score against the deployed docs).
  Recorded metric, not a build-time gate — needs the live Pages URL, which exists only after the workflow deploys. Deferred to validation/post-deploy (see Deferred follow-ups).

### Summary

Stood up the MkDocs Material docs site to the approved 7-section IA, green under `mkdocs build --strict` (the load-bearing gate; verified from a clean build, exit 0). `docs/site/` is the single content root: install-journey MOVED in via `git mv`, five functional files SYMLINKED (git mode-120000) so their canonical consumers keep reading them at fixed paths. Wired mkdocs-llmstxt (`site/llms.txt`), a reachable AGENTS.md page, the "For agents" banner, the Material palette toggle + code-copy, and a GitHub Pages workflow (strict PR check + deploy-on-main). Authored the net-new synthesized pages, the Voice & tone page (lifted verbatim from the ideation draft), the e6 stub, and the Reference → Frontmatter contract slot. Move-induced breakage fixed (README link, install.sh hint, install_doc_test.go path); `go test ./...` green (1247 passed). Deferred per captain: the no-duplication policy (no snippet machinery; light README↔site overlap accepted) and the standalone frontmatter-contract spec. AC-10 (a14y) and AC-6 deploy smoke are post-deploy live checks for validation. Code on branch `spacedock-ensign/mkdocs-material-docs-site`, commit 14cab445.
