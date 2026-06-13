---
id: pt0mz5stt4c1ve7ynz24r3yv
title: docs site — address reviewer feedback (simplify, dedupe, fix rendering)
status: validation
source: "PR #343 intake (captain 2026-06-12). Addresses reviewer (Karen) feedback on the docs site — headline asks: the docs are wordy / Spacedock feels complex, and apply the Recce doc-writing principle. Intaken directly at implementation for captain-steered interactive revision on the existing PR branch."
started: 2026-06-12T17:34:58Z
completed:
verdict:
score:
worktree: .worktrees/docs-site-feedback
pr: "#343"
issue:
mod-block:
---

Intake of **PR #343** (`docs-site-feedback` → `main`) for continued, captain-steered revision. The PR already does two things; this entity tracks refining it to done.

## Problem

Reviewer feedback on the docs site: the writing is wordy and Spacedock reads as complex, and the docs should follow the Recce doc-writing principle (structure + simplicity). Three issues were visible in the rendered site (e.g. a dark-on-dark header in slate mode; the Home page leading with the "multi-agent orchestrator" claim and ~7 terms up front).

## Approach (in-flight on the PR; revision is captain-steered)

PR #343 already:
1. Adds `docs/site/CLAUDE.md` — a standing authoring directive (adapted from Recce's `doc.md`) governing structure + simplicity for everything under `docs/site/`, auto-loaded when editing docs, deferring to `voice-and-tone.md` for voice, and excluded from the published build via `exclude_docs`.
2. Acts on each feedback item — fixes that needed no decision plus ones gated on a maintainer call (resolved with the maintainer).

Remaining work is **interactive**: the captain steers the doc voice/content revisions directly with the ensign. The ensign applies that steering on the `docs-site-feedback` branch (so #343 updates), keeps the rendered docs correct, and does not regress the build.

## Out of scope

- Opening a new PR or new branch — revisions land on the existing `docs-site-feedback` branch.
- The docs-deploy mechanism (the landing site owns deploy; the mkdocs Pages job is descoped) — this is content/voice, not deploy.

## Acceptance criteria

(Captain-steered; validation confirms. Each verified by the rendered docs / build, not a prose-grep.)

**AC-1 (sketch) — the reviewer feedback items are addressed in the rendered docs.** Verified by: the rendered site shows the simplified Home (problem-first, fewer terms) and the resolved rendering issues (e.g. header legible in slate mode) — observed in the built site / before-after, not the source prose alone.

**AC-2 (sketch) — the authoring directive exists and is excluded from the published build.** Verified by: `docs/site/CLAUDE.md` present and absent from the built `site/` output (`exclude_docs` honored) — observed in the build output.

**AC-3 (sketch) — the docs build stays green.** Verified by: `mkdocs build --strict` exits 0 after the revisions.

## Notes

Interactive intake — the captain steers revision in the ensign's pane directly. SSH push is currently down; the ensign commits locally on `docs-site-feedback`, and the FO pushes the updated PR branch (via the gh-HTTPS route or once the key is restored) when the captain is satisfied.

## Progress (session 1, 2026-06-12)

Phase 1 complete and committed on `docs-site-feedback` (19 commits ahead of main, working tree clean): structural restructure + tree-wide tone sweep + all of Karen's 2026-06-12 feedback. See the handoff in the next-session prompt for the full done-list and the Phase-2 plan.

Key protocol notes for the next session:
- `mkdocs build --strict` (AC-3) is **deliberately unenforced during iteration** per captain; build with plain `mkdocs build`, re-enable `--strict` at the end and fix the link punch-list (several inbound links now point at GitHub URLs / changed anchors after the Contributing/Advanced/Reference restructure).
- comm-officer polish is **best-effort and non-blocking** (2-min timeout, then proceed un-polished). Never block on it. `SendMessage` takes ONLY `to`/`summary`/`message` — adding `type`/`recipient`/`content` is what corrupted earlier sends.
- Phase 2 (per-page paragraph polish) NOT yet started on disk. comm-officer returned polish for `index.md` and `install.md` that is not yet applied — re-derive or re-request in the fresh session.

## Progress (session 2, 2026-06-12): Phase 2 per-page polish

Phase 2 complete across the full page order. 22 commits this session on `docs-site-feedback` (43 ahead of main, working tree clean); each page got an editorial pass against `docs/site/AGENTS.md` plus a comm-officer in-place polish (all replies landed within the window; none timed out). `mkdocs build --strict` exits 0 — the anticipated inbound-link punch-list from the Phase-1 restructure turned out to be already clean, so AC-3 needed no link fixes.

### Per-page line-change summary (for the captain's second-pass walk)

Diffstat is session-2 only (`cc473605..HEAD`, net -24 lines). index.md and install.md were polished last session; advanced/ and reference/ passed clean with no edits.

- **get-started/first-launch.md** (+5/-5): tightened lede (setup sentence deduped, safehouse moved to grammar section only), dense survey-report sentence split, redundant repeat-link dropped.
- **get-started/first-workflow.md** (+19/-24): front-loaded 3-term glossary block + roles paragraph replaced by definitions at first use (captain in lede, entity in the design-phase question, gate at the gates section); "This page walks…" announce sentence cut.
- **concepts/operating-model.md** (+5/-7): page-announce lede sentence cut; next-links parallelized; reverted one comm-officer inversion ("standing job entire").
- **concepts/workflows-and-entities.md** (+12/-18): entity section leads with file forms instead of restating the lede; line-oriented-parser note deduped; per-stage flag bullets collapsed to one sentence + link (stage-lifecycle now owns the property table); duplicate `--next` code block dropped.
- **concepts/stage-lifecycle.md** (+7/-11): "This page uses…" opener merged into content; duplicate status code block merged; next-links activated.
- **concepts/gates-and-decisions.md** (+13/-8): announce lede cut, three redundant tails trimmed, internal test-scenario reference cut; gained a "Where to go next" block (was the only concept page without one).
- **concepts/worked-example.md** (+15/-15): opens with the entity instead of "This page traces…"; sentence fragment fixed; "AC" glossed on first use; internal codenames removed (Commander, DoD).
- **running-workflows/commission.md** (+6/-6): "pleonastic"→"redundant"; tighten-the-README leverage point deduped (the dedicated section owns it).
- **running-workflows/survey.md** (+4/-4): duplicated read-only claim cut; "Gemini is a deferred follow-up" roadmap aside cut.
- **running-workflows/operating.md** (+3/-3): announce lede cut; gate-call names aligned with gates-and-decisions (approve / redo with feedback / reject).
- **running-workflows/debrief-and-refit.md** (+3/-3): comm-officer polish only (page passed my directive check clean).

Open captain decision still deferred from session 1: index.md para-3 bar-sharpening vs the "bar sharpens as you use it" bullet read as partly redundant; both left in, raise when walking the page.

## Stage Report: implementation

- DONE: Phase-2 per-page polish completed across the remaining page order (starting at get-started/first-launch.md; index.md and install.md already applied), with a per-page line-change summary for the captain's second-pass walk.
  11 pages edited over 22 commits (`cc473605..HEAD`, net -24 lines); summary recorded above; advanced/ + reference/ light pass found no edits needed.
- DONE: AC-3 re-enabled at the very end only: a clean `.venv-docs/bin/mkdocs build --strict` after fixing the inbound-link punch-list from the Contributing/Advanced/Reference restructure.
  `.venv-docs/bin/mkdocs build --strict` exits 0; the punch-list was empty (no warnings to fix). AC-2 also re-verified: built `site/` contains no `AGENTS.md` or `contributing/`; `llms.txt` + `llms-full.txt` present.
- DONE: All revisions committed locally on the existing `docs-site-feedback` branch so PR #343 updates — no push (FO coordinates), no new branch or PR.
  Branch `docs-site-feedback`, 43 commits ahead of main, working tree clean; no push attempted on the code branch.

### Summary

Phase 2 per-page polish is complete: every page in the reader-journey order got a directive pass (announce-sentence cuts, terms defined at first use, cross-page dedup with stage-lifecycle now owning the stage-property table) plus comm-officer polish, committed page by page. The strict build is green with no link fixes required, and AC-2 exclusions re-verified in the build output. The branch is ready for the captain's second-pass page walk and the FO-coordinated push to PR #343.

## Progress (session 2, captain-steered walk)

After the Phase-2 stage report, the captain walked the site live and drove a major revision wave (~40 commits): Get started fully rebuilt (install ~120 words; "Your first launch" renamed to get-started/survey.md; first-workflow example-led with real commission/gate samples), Concepts to the new standard (worked-example.md KILLED; mermaid stage diagram added; gates page leads with delegation + the conn prompt), reference/sandbox.md added, and AGENTS.md extended in three codification rounds (now including the per-paragraph revision loop). Strict build green throughout; branch ~105 ahead of main; no pushes (SSH down).

## Next-session handoff (Phase 3) — for the FO to dispatch a fresh ensign

# HANDOFF — docs-site-feedback (PR #343), Phase 3 (Running workflows + close-out)

## Role & setup
You are the ensign on PR #343 (`docs-site-feedback` → `main`), continuing a captain-steered interactive revision of the Spacedock mkdocs site. On start, run `Skill(skill="spacedock:ensign")`.

- Code worktree: `/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/docs-site-feedback`
- Branch (code): `docs-site-feedback` — commit here so PR #343 updates. Do NOT create a new branch/PR. Do NOT push (SSH down; the FO coordinates pushes).
- Entity (split-root state): `/Users/clkao/git/spacedock-research/spacedock-v1/docs/dev/.spacedock-state/docs-site-feedback-revision.md` — path-scoped state commits only, never bare `git add -A`.
- Build with `.venv-docs/bin/mkdocs build --strict` after every edit; it is green now and must stay green.

## Read FIRST
`docs/site/AGENTS.md` — the single authoring directive, heavily extended this session; it is the bar. Load-bearing rules: reader's seat, never the maintainer's; agent-facing surfaces (`spacedock status`, `spacedock dispatch`) off user pages; the agent is the interface (describe capabilities, never field-by-field key-in tables); durable value, not output inventory; one real sample from the product's real templates beats description; built-ins must read as built in; verify behavior claims against the owning skill before writing them; link by payoff; headings take the reader's angle.

## The per-paragraph revision loop (captain-mandated; also codified in AGENTS.md)
For EVERY paragraph ask: (1) why would the reader care, (2) do they have the context at this point in the reader journey, (3) does it make them want to read on. Propose the revision YOURSELF, send your proposal with its three-question rationale to comm-officer (text passthrough; say "reply with notes, do NOT edit files"), then TRIAGE the return: accept real improvements; reject em-dash insertions (its recurring offender) and anything that alters captain-set wording; re-check the file on disk before acting on late-arriving replies. comm-officer is best-effort and non-blocking (2-minute timeout, then proceed). `SendMessage` takes ONLY `to`/`summary`/`message`.

## Captain interaction
Interactive and terse; expects motion. Established pattern: when asked "what's wrong with this page?", deliver findings and STOP; apply on "do it". Captain-dictated wordings are verbatim-sacred (e.g. the conn prompt on the gates page, "send it back unless this now needs reframing").

## Done — do not rework
- Get started (Welcome, Install, Survey your project, Your first workflow): rebuilt under direct captain steering.
- Concepts (operating model, workflows & entities, stage lifecycle, gates & decisions): current standard. `worked-example.md` is deliberately KILLED (mixed sprint machinery into concepts); do not resurrect it.
- `reference/sandbox.md` (title + table only), nav updated, all GitHub links point at `main`.

## Remaining — your job, in order
1. **running-workflows/operating.md** — the biggest job. It still teaches `spacedock status` flags (`--where`, `--validate`, `--resolve`) to humans: an agent-facing violation. Recast as touchpoints: the typical session loop (align with first-workflow's "Session to session" close), gates arriving, delegation. Verify behavior claims against the first-officer and feedback-rejection-flow skills.
2. **running-workflows/commission.md** — dedupe against the new first-workflow (which now owns the design-summary, pilot, and gate samples); run "The four things you name" through the capabilities-not-fields lens; ID-style detail is a candidate for "ask the agent".
3. **running-workflows/debrief-and-refit.md** — re-check under "the agent is the interface" and durable-value rules; the refit per-file strategy bullets may be output inventory.
4. **reference/command-reference.md** — reference pages may document commands, but check the Workflow group is framed as what the agents run, where true.
5. After any page changes: re-read touched pages in journey order (term-before-introduction check), strict build, AC-2 exclusions hold (no AGENTS.md or contributing/ in `site/`; `llms.txt` + `llms-full.txt` present).
6. Update the entity stage report; path-scoped state commit; the state push will fail while SSH is down — report it, never force.

## Open follow-ups to FILE as backlog seeds, not to do
- commission skill template prints "Entity identity will use `sd-b32`." — reads internal in user-facing output; fix belongs in `skills/commission`, not the docs. Captain: "let's worry about sd-b32 later."
- A new modality for examples (replacing the killed worked-example page) — captain decision pending.
- Eyeball the stage-lifecycle mermaid diagram on the slate theme (first mermaid on the site; gold `#e3b04b` gate borders).

## Stage Report: implementation (cycle 2)

- DONE: The four remaining handoff pages (operating, commission, debrief-and-refit, command-reference) reworked with behavior claims verified against the owning skills.
  operating/commission/debrief-and-refit were already reworked by predecessor c2 (re-verified clean this session: operating recast as touchpoints, commission deduped, debrief split from refit by cadence). command-reference was the genuinely-remaining page — fully reworked this session (below).
- DONE: command-reference Workflow group framed as what the agents run, verified against the owning source.
  Split into per-section tables (dropped the monolithic top table), Launch made prose-only (DRY), Workflow group framed as agent-driven (the first officer runs status/new/dispatch/state/completion). Every command verified against the cobra defs in `internal/cli/cli.go` — caught and fixed a real error: the `state` row only documented `state init`, corrected to cover both `state init` (resumes) and `state new` (births). `8017d918`, `e5cf2f9c`.
- DONE: Frontmatter contract ported from the v0 branch into the repo; mdschemas are SSOT; the site links to it.
  Captain-directed addition beyond the handoff. `docs/schema/{entity,workflow-readme}.mdschema.yml` ported from the unmerged v0 branch `spacedock-ensign/spacedock-frontmatter-contract-spec` (entity `known_corpus_failures` emptied — v0 files). The site `reference/frontmatter-contract.md` is a reader-seat stub: lede + two subsections (Workflow README, Entity) linking the schemas on `main`; field tables and the line-parser prose dropped (the Go parser uses `yaml.Unmarshal`, verified — the old "flat/nested-ignored" claim was false). `3fba5836`..`90db73df`, `1f17ed3f`.
- DONE: After page changes — journey re-read in nav order (term-before-introduction), strict build green, AC-2 exclusions hold.
  All 19 pages re-read in nav order: term introductions ordered, cross-page references consistent (operating's "no status commands to learn" now agrees with command-reference; frontmatter-contract inbound links match the stub). `.venv-docs/bin/mkdocs build --strict` exits 0, no warnings; `site/` has no AGENTS.md or contributing/, `llms.txt` + `llms-full.txt` present; `docs/schema/*.mdschema.yml` not leaked into `site/`.
- DONE: Open follow-ups reported as backlog-seed candidates; all revisions committed locally on docs-site-feedback — no push, no new branch/PR.
  New backlog seed filed to the FO (team-lead): full mdschema conformance validator — `spacedock status --validate` enforces only a subset of the schemas; either extend Go `--validate` or add a conformance checker (the v0 Python validator was NOT ported). Prior follow-ups above still stand. Branch `docs-site-feedback`, 117 commits ahead of main, working tree clean (only a pre-existing untracked `nohup.out`); no push attempted.

### Summary

This session (cycle 2) was captain-steered and interactive. The handoff's four pages were either already done by c2 (re-verified) or reworked here (command-reference). Beyond the handoff, the captain directed porting the v0 frontmatter/state-machine contract into the repo: the two mdschemas land in `docs/schema/` as the field SSOT, and the site reference page became a reader-seat stub linking them. Reworking command-reference surfaced and fixed a real doc error (missing `state new`) and confirmed the Go parser no longer enforces flat-only frontmatter. Full journey re-read clean, strict build green, AC-2 holds. State push will fail while SSH is down — reported, not forced.

## Stage Report: validation

- DONE: Each AC-N reproduced with evidence from OUTSIDE the task body; reconciled against as-shipped reality (CLAUDE.md -> AGENTS.md rename), flagging any AC lacking external evidence.
  AC-1: rendered `index.md` leads problem-first (gate idea, "you are the captain"), no "multi-agent orchestrator"/7-term opener; slate-header fix is concrete — `docs/site/stylesheets/brand.css:11` sets `color: var(--fg)` so the wordmark/tabs stay legible under slate, present in built `site/.../brand.css` (`grep -c var(--fg)` = 1). AC-2: source `docs/site/AGENTS.md` (19.7K, the renamed directive) + `contributing/` present, `exclude_docs` lists `AGENTS.md` + `contributing/`, both absent from built `site/`. AC-3: `.venv-docs/bin/mkdocs build --strict` exits 0. All three ACs proven by build output / on-disk state, none self-referential.
- DONE: `.venv-docs/bin/mkdocs build --strict` re-run independently exits 0; built `site/` excludes AGENTS.md and contributing/, includes llms.txt + llms-full.txt; `docs/schema/*.mdschema.yml` not leaked into `site/`.
  Independent clean rebuild EXIT=0, zero warnings/errors. `find site -iname 'AGENTS*'/-iname 'CLAUDE*'` empty; no `site/contributing/`; `site/llms.txt` (2.1K) + `site/llms-full.txt` (38K) present; no `*.mdschema.yml` under `site/`.
- DONE: PASSED/REJECTED recommendation confirming the docs reflect the latest captain feedback — reject if prose is stale, over-specified, or contradicts as-shipped pages.
  PASSED. Spot-checks against source-of-truth confirm currency: `worked-example.md` killed (absent from source + nav); `operating.md` no longer teaches `spacedock status` flags (`--where/--validate/--resolve` absent — agent-facing surface off the user page); command-reference `state init`/`state new` row matches the cobra `Short: "init resumes, new births"` in `internal/cli/cli.go:341`; the false "flat/nested-ignored" line-parser prose is gone everywhere in `docs/site/` and the live parser uses `yaml.Unmarshal` (`internal/status/frontmatter.go:99`); frontmatter-contract is a reader-seat stub linking the two schemas on `main`. 20 rendered pages = home + 19 source pages, no orphans.

### Summary

PASSED. Fresh independent validation of PR #343 `docs-site-feedback` (branch at 117 ahead of main, tree clean apart from the pre-existing untracked `nohup.out`). The strict build re-runs green with zero warnings; all three build-output assertions hold (AGENTS.md + contributing/ excluded, llms.txt/llms-full.txt included, schemas not leaked). The AC sketches predate the rework, so I reconciled against as-shipped state: AC-2's CLAUDE.md is now AGENTS.md and `exclude_docs` correctly names it plus contributing/. Spot-checks confirm the docs reflect the latest captain feedback and contradict nothing as-shipped — the command-reference and frontmatter-contract reworks match the Go source-of-truth, and the killed worked-example / recast operating page are real on disk. No re-edits made (read-only validator); no push attempted.
