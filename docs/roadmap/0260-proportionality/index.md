# 0260 — Proportionality (0.26.0)

**Sprint:** the entities matching `sprint: 0260-proportionality` — list current members with `spacedock status --workflow-dir docs/dev --where sprint=0260-proportionality`. Membership and per-task state are the query, never enumerated or tracked in this doc.
**Theme:** the contract prices over-engineering, not just under-verification. A 2-week forensic audit of agent sessions across spacedock_v1/zaphod/spacedock_subspace (15 confirmed incidents, 4 HIGH-severity runaway review loops) found the agents mostly *complied* with the contract: every shipped code mechanism enforces evidence-production or momentum, every anti-over-engineering guard is prose-only, and nothing anywhere can express "this is a prototype." 0260 makes rigor proportional to declared stakes, makes evidence falsifiable, and cleans the fabricated rigor (tautological tests, prose-phrase lint checks) out of the estate.
**Evidence:** `docs/dev/.spacedock-state/_evidence/0260-agent-derail-forensics/` — 18 incident records with session:ordinal citations + adversarial verify verdicts, and 3 remedy-coverage analyses (FO contract, workflow READMEs/template, roborev configs). Entity Problem sections cite into it.
**Re-lock (2026-07-20):** the captain applied the sprint's own thesis to the roster. Essence restated: (1) the feedback loop learns to reframe instead of repair; (2) evidence must be able to fail; (3) disciplines that already exist gain reach — into packets, reviewers, and the template. The stakes member reduces from a new ontology to a read-through of existing declared posture; the design-reset routing merged into the cycle-record member; the smallest-sufficient sharpening folded into the ladder member; estate cleanups not killing an incident class carry `sprint-readiness: defer`. Drivable set: `spacedock status --workflow-dir docs/dev --where sprint=0260-proportionality --where 'sprint-readiness != defer'`.
**Release:** stable **0.26.0** (current line: 0.26.0-pre1).

## Goal (success criterion)

An ensign, reviewer, or FO holding a finding, an AC, or an urge to build infrastructure has a declared answer to "how much rigor does this project want?" — and the loop that couldn't reframe now can:

- ~~The stakes/posture member~~ — **parked** (captain hold, 2026-07-20, after the re-lock's essence test): its would-be consumers are already served — validators anchor materiality on per-entity value ACs + the committed finding-triage, and roborev injects its config posture line. The read-through's riskiest mechanisms are spike-proven and banked in the parked entity for revival if a later member surfaces a genuinely unreached consumer. Direction trail preserved in the entity; the residual README posture-consolidation edit rides the `template` group.
- A rejection whose findings indict the mechanism's architecture — or a repair cycle that grows the diff — halts at a design-reset decision instead of dispatching repair (`reframe` group).
- An ensign triages review findings against stakes before fixing; a correct-but-disproportionate finding gets a recorded decline, not a dutiful fix (`triage` group).
- The falsifiability ladder replaces "prefer a code gate over a prose-only rule": shipped system guards → existing mechanical checks → falsifiable exercise → captain judgment → build new machinery (last, consent-gated). New enforcement surfaces are not "obvious reversible work"; investigations hit a fan-out checkpoint; identifier minting is reserved to the system, ad-hoc itemization uses bare ordinals (`ladder` group).
- The tautological-test estate is fixed and gated against recurrence; contractlint's prose-phrase checks retire in favor of behavior tests (`test-cleanups`, `contract-cleanups` groups).
- The dev template ships the scar tissue (Stakes scaffold, materiality taxonomy, offline/interactive AC split, small-change fast path, fixed Verified-by example) and refit propagates content, so the next commissioned workflow starts protected (`template` group).

## Layer map

| Layer | What lands there |
|---|---|
| Spacedock contract (skills + binary) | reframe routing + dispatch-refusal gate; stakes read-through (boot, dispatch build); ensign finding-consumption rule + decline disposition; falsifiability ladder + infra consent + fan-out checkpoint + bare-ordinal rule; "5/5 passed is sufficient" removal |
| Dev template + refit | `## Stakes` scaffold + commission rigor question; materiality taxonomy; AC split; small-change fast path; size-gated semantic adversarial pass; Verified-by example fix; refit content propagation; project AGENTS.md router scaffold + maintenance |
| Project AGENTS.md / CLAUDE.md (router layer) | one-line stakes digest + pointer to the workflow README as process authority; auto-loaded into ad-hoc sessions the README never reaches. Source of truth stays the README `## Stakes` — one source, three channels (AGENTS.md digest, dispatch packet, roborev config); never a fourth divergent copy. AGENTS.md-first for codex coverage; Claude ingestion verified by canary, not assumed |
| Repo-local (this repo) | testlint/contractlint changes; the 8 tautological-test fixes; roborev config alignment; refresh the stale AGENTS.md `## Priorities` (bootstrap-era "do not add PR or mod behavior" still shipping to every agent) |
| CL's personal CLAUDE.md | out of sprint scope — separate audit in flight; carve-out drafted for CL to apply by hand |

## Definition of Done

`v0.26.0` ships when, merged to `main` and proven by checks that can fail:

- Replaying the archived e6j cycle history (26 files / +3,373 on a 2-defect fix, 10 cycles) against the new feedback flow halts at a design-reset decision by cycle 2; the dispatch-refusal guard has a unit test fed that fixture shape.
- A seeded correct-but-disproportionate finding against a low-stakes fixture entity produces a recorded decline and a zero-line diff in live replay.
- A replayed ideation brief mandating a process-control harness (the 7h PTY shape) trips the consent stop before dispatch.
- testlint fails red on the reverted 11-phrase contract-presence test and passes green on `main`; the 8 confirmed tautological output-grep tests are fixed; gate review reads new-test assertion content, not pass counts.
- The runtime-semantics contractlint retirements land (prose-phrase checks replaced by live/fixture behavior tests); the remaining four retirements carry `sprint-readiness: defer` for the next train.
- A scratch workflow commissioned from the updated template contains the materiality taxonomy and the fixed Verified-by example; a refit dry-run against a commissioned README shows the content delta arriving.

## Constraints

- **Leanness (inherited from 0250):** net contract-byte delta vs the 0.25 baseline is measured; additions prefer lazy-loaded references over boot-resident lines.
- **Live-drive proof rule (inherited from 0250, verbatim):** every behavioral claim is proven by a live drive observing the behavior — never a prose-grep over the contract the change writes. The incident records in `_evidence/` are the replay fixtures.
- **Bare-ordinal itemization:** no minted reference schemes in sprint docs or entity bodies; sanctioned identifiers only (entity ids/slugs, AC-N, session:ordinal).
- **This sprint practices its own thesis.** Shaping and driving decisions are held to the ladder: a proposed new check names the rung it sits on; a disproportionate-but-correct idea gets a recorded decline.

## Gate approvals — 3k notation (Commander: read this)

Ideation approvals from shaping are recorded durably in each entity's `gates:` frontmatter per the contract banked in `docs/dev/.spacedock-state/durable-gate-approval-pending-blockers/gate-resolution-frontmatter-contract.md` (a convention dry-run; the 3k binary recorder is NOT in this sprint):

- A closed attempt with `resolution.decision: approve` and `application: {action: advance, target-stage: implementation, state: pending}` IS the captain's approval. Do not re-ask the captain for a closed gate.
- Apply a pending application only when its declared `blockers` are satisfied and the entity content still matches `briefing.digest`; on drift, mark it `superseded` and open a new attempt.
- Consumption (`state: consumed`) happens only through the normal transition/dispatch path, once.

## Sequencing

- `stakes` is parked (see Goal) — no member blocks on it; the `triage` rule cites the committed finding-triage taxonomy directly.
- `reframe` leads: `bw` carries the merged design-reset scope (kills all four HIGH incidents).
- `triage`, `ladder`, `template` are unblocked and parallel; `z7` carries the folded smallest-sufficient sharpening plus four recorded live examples of language-minting for its identifier/abstraction clause.
- `test-cleanups` and the driving `contract-cleanups` pair are independent and fully parallel; the banked `0qe` ideation is not re-ideated — its overlap with `w0` is judged at its gate.

## Out of scope

- CL's personal CLAUDE.md (separate, personal).
- Direct edits to zaphod / spacedock-subspace repos — their READMEs receive the delta via refit under their own workflows after `template` lands.
- Implementing 3k itself (notation used as convention only).
- `mzk` (shared-GitHub-state ownership tracing) — real, different theme, stays in backlog.

## Lifecycle checklist

**Shape — Shaping FO**
- [x] **Scope-lock** with the captain — locked 2026-07-20 full-roster, then **re-locked same day** under the sprint's own essence test: two merges (`ve`→`bw`, `1p9`→`z7`), the stakes reduction, and seven `sprint-readiness: defer` stamps (`fw` `1w` `h6` `b7` `3a` `xaz` `cy`)
- [x] **Carve** — 23 members stamped (`sprint` / `group`), 5 new entities filed with evidence citations, index.md written
- [ ] **Ideate** each gated member — riskiest mechanism first (stakes read-through spike leads)
- [ ] **⚠️ Preflight staff review (sprint-wide)** — independent reviewer, refute the sprint as a whole → `staff-review.md`
- [ ] **Present ideation gates** — AC cross-check per member; approvals recorded in 3k notation *(captain decides)*
- [ ] **Package** — `dispatch-sprint-execution.md` (cold-boot Commander package)

**Drive — Commander (separate cold-booted session)**
- [ ] Implementation → validation → done per member; detached adversarial audit for high-stakes surfaces
- [ ] Merge each to `main`; state commits concurrency-safe
- [ ] **⚠️ Pre-cut antipattern audit** — independent staff-eng reviewer over the assembled sprint, before the tag
- [ ] **Cut 0.26.0** — `go test ./...` green, then `docs/releasing.md` *(captain authorizes)*

**Close — Shaping FO**
- [ ] Seed the next sprint from deferred findings + post-cut release verification
