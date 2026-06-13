---
id: 5wvrtfjvjz78fy9xg55p6pjg
title: Survey output redesign — one coherent "value & numbers first" rewrite folding all six feedback seeds
status: ideation
source: "Coalesces the survey-improvement cluster (six user-sourced feedback seeds, 2026-06-08→06-12) into one redesign of the spacedock:survey output. Captain call 2026-06-13: six seeds all rewrite one skill = one coherent change, not six worktrees colliding on SKILL.md."
started: 2026-06-13T04:53:59Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0202-survey-improvements
group: survey
sprint-readiness: ready
---

One coherent rewrite of the `spacedock:survey` output into the captain-locked **"value & numbers first"** structure (canonical mock in `docs/roadmap/0202-survey-improvements/index.md`). The six feedback seeds below are the requirement record; this task is their single execution. All proof is the survey's rendered output over constructed fixtures — never a prose-grep of SKILL.md (the survey discipline). Do NOT regress the validated inference accuracy or the decision-frontier triage (three+ users confirmed these — the keep-signal).

## Requirement bands (each cites its source seed; all render INTO the locked structure)

- **R1 — value-prop legibility (`9h` survey-value-prop-legibility).** Plain "what this gives you" lede + concrete numbers from the user's own data; rename **"mechanical" → "manual"** for repetitive-substantive tracks (reserve "mechanical" for trivial); a **"threads to pull"** actionable section (unresolved/hanging + proactive suggestions), not decision-history narration.
- **R2 — honest lens (`5x` survey-lens-honesty).** Recent-window snapshot framing (surface evolution where the signal exists, else say "recent window, not whole history") + ONE clear partial-lens caveat naming what the corpus can't see.
- **R3 — subagent-dispatch fact (`za` survey-report-subagent-dispatch-fact).** Report the FACT that a session dispatched subagents so an orchestrated repo isn't read as idle (count/presence only; subagent CONTENT stays excluded).
- **R4 — output hygiene (`h5` survey-output-polish).** Collapse the always-empty Codex `workstreams: (unlabeled) — N` breakdown to one honest line; strip the model scratch-reasoning preamble (`I have everything I need. Let me…`) — cross-check findings stay as report content.
- **R5 — knowledge-work archetype + Codex count (`zb` survey-knowledge-work-archetype).** Name a third archetype (knowledge-work loop) beside mechanical/manual-drive and exploration-steering; lead the Codex section with the workdir-attributed count (name-match demoted to caveat) — **verify against the current skill; zb#1 may already be shipped.**
- **R6 — mode-aware framing (`zw` survey-iteration-framing-and-branch-attribution).** Exploration mode leads with iterate/steer (not gates-as-the-whole-model); manual/mechanical mode keeps gate-drive framing. Branch-aware work-by-area as a **detect-branch-and-merge + caveat** (conditional — it did not bite two of the runs).

**Cross-revision rule (from `h5`'s 2nd round):** verify each ask against the CURRENT skill before building — some may be partly shipped (zb#1) or need narrowing (zw#1/#2). **Pending input:** the author-followup reply (with the captain) folds into R5/R6 — ideate-and-revise, do not block.

## Verification against the current skill (ideation, 2026-06-13)

Each band was checked against the live `skills/survey/SKILL.md` + `references/queries.sql` (HEAD on `next`, after #335) before specifying, per the cross-revision rule:

- **R1 (value/legibility) — NOT shipped.** Current step-4 body opens with the synthesis-fence headline (`Found {N} sessions … here's the lay of the land:`) then `PROJECT:` — no plain value lede, no `BY THE NUMBERS` block. "mechanical" is the live mode label and offer vocabulary throughout. "Threads to pull" framing absent (the frontier renders as `NEEDS YOU`/`BACKLOG`/`RECENT DECISIONS`, decision-history-shaped). All of R1 is new work.
- **R2 (honest lens) — NOT shipped; ALREADY FIRMED in `5x`.** `5x`'s ideation stage is complete with a concrete doc diff (title+subhead snapshot framing + one consolidated `WHAT THIS CAN'T SEE` block replacing scattered asides) and a "no spike needed" determination reusing `scoping.span` + `scoping.blank_cwd` + the `codex-presence − codex-scoped` gap. R2 here ADOPTS `5x`'s firmed AC-1/AC-2 and doc diff verbatim — not re-designed.
- **R3 (subagent-dispatch fact) — NOT shipped; ALREADY FIRMED in `za`.** `za`'s ideation is complete with the `dispatch-fact` query, a live-DB spike (the parent-join marker; `subagent_session_id` is a dead end in agentsview v0.32.1), full ACs, and the `BY THE NUMBERS` render line. R3 ADOPTS `za`'s firmed query + ACs verbatim. Re-exercised on this repo's live data during this ideation: `dispatch-fact` returns **31 orchestrating parents / 1675 subagents** (matches `za`'s 31/1673 spike + 2 newer sessions) — proven.
- **R4 (output hygiene) — NOT shipped.** No `ad-hoc`/`unclassified`/`scratch`/`preamble` text in the current skill (grep clean). The Codex `workstreams:`/`activity:` lines render unconditionally (no collapse-when-all-unlabeled). No instruction suppresses the model's `I have everything I need. Let me…` scratch preamble. Both R4 asks are new work.
- **R5 (knowledge-work archetype + Codex count) — zb#1 SHIPPED, zb#2 NOT.** zb#1 (lead with the workdir-attributed Codex count, demote name-match to caveat) shipped in #335 — verified at SKILL.md step-4 `CODEX` block (`{codex_scoped_sessions} … attributed to this repo by exec_command working dir` + the `by project NAME only` caveat). **R5 narrows to JUST the knowledge-work archetype** (zb#2): `mode-classification` has only `mechanical`/`exploration`/`unlabeled` — no knowledge-work mode. Re-specing the Codex-count-first behavior is banned (already shipped).
- **R6 (mode-aware framing + branch attribution) — NOT shipped; narrowed per h5.** Current offer leads with gates for mechanical AND the headline `gated automation for the mechanical tracks` for both. zw#1 narrows to **exploration-mode-specific** (mechanical KEEPS gate framing — do not strip). zw#2 (branch-aware work-by-area) is **conditional** (detect-branch-and-merge + caveat; did not bite two of h5's runs) — lowest blast radius, lowest priority.

## Riskiest-mechanism exercise (DONE — run on this repo's live agentsview DB)

The assignment names the riskiest unknown: *are the lede's concrete figures computable from the existing session data?* The four `BY THE NUMBERS` figures in the captain-locked mock were each exercised against this repo's synced agentsview DB (v0.32.1, 66 in-repo Claude parent sessions, 2026-05-30 .. 06-13):

| Lede figure (mock) | Data source | Result on live data | Verdict |
|---|---|---|---|
| `hand-steering interruptions` | `decision-open` decisions (`AskUserQuestion`/`ExitPlanMode`) + veto markers (`[Request interrupted` / `doesn't want to proceed`) over the repo scope — the EXISTING `{V}` interruption total | 114 decisions + 19 vetoes | **PROVEN — already-computed `{V}`, no new query** |
| `hanging threads (started, never closed)` | the **count** of `decision-open` OPEN rows that survive the step-4 repo cross-check (the `NEEDS YOU` set) | 10 OPEN / 104 done (pre-cross-check) | **PROVEN — derived count of an existing query+cross-check, no new query** |
| `decisions you made with no follow-up action` | the **count** of cross-check `decided-not-shipped` forks (the `BACKLOG` set) | subset of the 10 OPEN after cross-check | **PROVEN — derived count of the existing cross-check, no new query** |
| `dispatched-subagent count` | `za`'s `dispatch-fact` (parent-join marker) | 31 parents / 1675 subagents | **PROVEN — `za`'s spiked query** |

**Key finding:** the lede's three decision-frontier figures are NOT new queries — they are COUNTS of the sections the existing `decision-open` query + the step-4 repo cross-check already produce. The redesign surfaces those counts in the lede; it does not compute new signal. The ONLY new query the whole redesign needs is `za`'s already-spiked `dispatch-fact`. Per the ideation gate's spike rule, this is recorded as **"no spike needed beyond `za`'s `dispatch-fact` spike"** — every other lede figure composes already-proven reads (`decision-open`, the veto markers, `scoping.span`/`blank_cwd`, the `codex-presence − codex-scoped` gap), all re-exercised above.

## Chosen design — the consolidated output (ideation, 2026-06-13)

One step-4 report-template rewrite renders the captain-locked **value-&-numbers-first** spine. The structure is captain-locked (`0202-survey-improvements/index.md`); this design maps each band onto it, it does NOT redesign it. The rendered report, top to bottom:

```
SpaceDock survey — your last {N} days                          ← R2 (title carries window)
(recent-window snapshot · agent logs only{· {blank_cwd} sessions had no working dir, not placed})   ← R2 subhead

WHAT THIS GIVES YOU                                             ← R1 value lede (plain language)
  {plain "you steer your agents by hand ~{interruptions} times a {window}; about {repeated} are
   the same few moves. A SpaceDock workflow can run those for you and stop only where you'd want a say."}

BY THE NUMBERS                                                  ← R1 + R3 (concrete figures, fixture-derived)
  {interruptions}  hand-steering interruptions
  {hanging}  hanging threads (started, never closed)            ← count of NEEDS YOU (post-cross-check)
  {no_followup}  decisions you made with no follow-up action    ← count of BACKLOG (decided-not-shipped)
  {sessions}  sessions read ({codex line if codex-scoped>0; name-match caveat})
  {if dispatch>0: {orchestrated}  sessions dispatched subagents ({dispatched} dispatched — their work isn't shown here)}   ← R3

HOW YOU WORK                                                    ← R1 (de-jargoned) + R5 (archetype) + R6 (mode label)
  {the inferred loop as an arrow chain} — {one honest line naming the dominant mode in PLAIN terms:
   "Mostly manual, repetitive tracks (not trivial — they take real work)." for manual;
   "Mostly exploratory — you steer an iterating agent." for exploration;
   "A knowledge-work loop: intake → process → file → log → close." for knowledge-work}

  ↓ full analysis: modes, work-by-area, what we can't see       ← R2 pointer

═══ (everything below the fold is the demoted detail section — the existing body, de-jargoned) ═══

THREADS TO PULL   (only if any OPEN fork survives the cross-check)   ← R1#4 reframe of NEEDS YOU
  {steady-state + the open/hanging forks + proactive "have you thought about …" suggestions —
   NOT decision-history narration; exploration forks framed as held-not-lost, mechanical as backlog}

BACKLOG   (decided-not-shipped)                                 ← unchanged section, feeds the no-follow-up count
RECENT DECISIONS                                                ← unchanged
HOW YOU WORK (detail): WORKSTREAMS · WORK BY AREA · CODEX · SCAFFOLD · INTERRUPTIONS   ← existing body, manual/mechanical vocab
WHAT THIS CAN'T SEE   (the agent-log corpus is a partial lens)  ← R2 consolidated block (from 5x's doc diff)

then: the discovery → commission bridge, mode-keyed (R6: exploration leads iterate/steer; manual/mechanical keeps gates)
```

**Vocabulary rename (R1#3, spans the whole render AND the query).** `mode-classification` emits the literal label string. The rename `mechanical → manual` for repetitive-substantive tracks is a one-word change to the query's `CASE … THEN 'mechanical'` → `'manual'` PLUS every render reference (`HOW YOU WORK` line, `THREADS TO PULL` sub-headers, the commission offer, the synthesis-guidance prose). "mechanical" is reserved for genuinely-trivial tracks; since the classifier has no trivial-vs-substantive split today, ALL current `mechanical` labels become `manual` (the substantive case is the only one the classifier detects). This keeps the query output and the rendered vocabulary in lockstep — the AC asserts the rendered label, derived from the fixture's track signals.

## Acceptance criteria (consolidated — one per band, R1–R6)

Each AC is an entity-level end-state property, verified by a **survey render over a constructed fixture** with the expected value DERIVED FROM the fixture's session rows (an independent source that diverges from SKILL.md prose); a prose-grep of SKILL.md never satisfies any AC. The query-level ACs (R3, R5, R6#2) are pinned by the `TestSurveyQuerySmoke` harness (`skills/integration/survey_queries_test.go` against `testdata/survey/fixture-sessions.sql`); the render-level ACs (R1, R2, R4, R6#1) are pinned by a live survey drive over the fixture/corpus, mirroring the existing survey live-drive ACs.

**AC-1 (R1) — the output opens with a plain value lede + a `BY THE NUMBERS` block whose figures derive from the fixture, and repetitive-substantive tracks read "manual" not "mechanical".**
Verified by: a survey render over the fixture shows (a) a `WHAT THIS GIVES YOU` plain-language lede before any mode/track vocabulary, (b) a `BY THE NUMBERS` block whose `hand-steering interruptions` equals the fixture's `decision-open` decisions + veto markers, `hanging threads` equals the fixture's post-cross-check OPEN count, and `sessions read` equals `scoping.sessions` — each number matching the fixture rows, not skill prose; (c) the repetitive-substantive fixture track rendered as `manual` (and `mode-classification` emits `manual`, asserted in the query smoke). Non-vacuous: a fixture mutation that adds an interruption flips the `hand-steering interruptions` count.

**AC-2 (R2) — recent-window snapshot framing + ONE consolidated partial-lens statement** (adopts `5x`'s firmed AC-1/AC-2).
Verified by: a date-spanning fixture renders the window in the title + the `recent-window snapshot` subhead clause with `{N}` matching `scoping.span`; an off-corpus-gap fixture (blank-cwd Claude session + name-only Codex superset) renders ONE `WHAT THIS CAN'T SEE` block naming `{blank_cwd}` and the `codex-presence − codex-scoped` gap, with NO second scattered copy of those caveats. Both numbers derive from fixture rows. (Test fixtures + doc diff are `5x`'s, adopted here.)

**AC-3 (R3) — the rendered survey names the dispatch fact in `BY THE NUMBERS`** (adopts `za`'s firmed AC-1/AC-2/AC-3).
Verified by: the `dispatch-fact` query, run over a fixture seeded with dispatched-subagent sessions, returns `sessions_that_orchestrated | subagents_dispatched` counting only subagents whose parent is an in-repo non-subagent Claude session (a subagent of an out-of-repo parent does NOT count; a non-vacuous re-point flips the counts) — pinned in the query smoke; AND a live survey drive over a fixture/corpus with dispatched subagents renders the dispatch line in `BY THE NUMBERS` with the count equal to the query, while a non-orchestrated corpus renders NO dispatch line. (Query + ACs are `za`'s, adopted verbatim.)

**AC-4 (R4) — the Codex breakdown never presents an all-unlabeled cluster list as a breakdown, and no model scratch-reasoning precedes the report.**
Verified by: (a) a survey render over a fixture whose Codex sessions are ALL `(unlabeled)` shows a single honest unclassified line (e.g. `{codex_scoped} Codex sessions, unclassified (ad-hoc shell-driven)`), NOT a `workstreams: (unlabeled) — N` breakdown row — the count from the fixture's codex-scoped rows; AND a fixture with ≥2 distinct named Codex workstreams STILL renders the breakdown (proving the collapse is conditional on all-unlabeled, not a blanket removal). (b) The rendered report begins at the `SpaceDock survey —` title with no `I have everything I need` / `Let me …` preamble — observed in the live-drive output; the cross-check FINDINGS still appear as report content (in `THREADS TO PULL` / `BACKLOG`). The all-unlabeled-collapse is enforced by render logic over the `codex-workstreams` rows (a query-smoke assertion on the all-unlabeled fixture: every row is `(unlabeled)`); the scratch-suppression is a live-drive assertion (the rendered first line is the title).

**AC-5 (R5) — a knowledge-work repo is classified and named, not left "generic book-keeping".**
Verified by: a survey render over a knowledge-work fixture (a process/file/log loop, content/ops edit profile, NO issue→worktree→PR signature, NO veto-heavy creative signature) classifies the track `knowledge-work` (asserted in the query smoke: `mode-classification` emits `knowledge-work` for that track) and the `HOW YOU WORK` line + the commission offer NAME the knowledge-work loop with an honest book-keeping offer, rather than the unlabeled fallback. Non-vacuous: a fixture mutation that removes the process/file/log loop markers drops the track back to `unlabeled`. (zb#1 is SHIPPED — NOT re-spec'd; this AC covers only zb#2, the archetype.)

**AC-6 (R6) — exploration offers lead with iterate/steer; manual/mechanical offers keep gates; work-by-area is branch-aware (detect + caveat).**
Verified by: (a) a survey render over an exploration-mode fixture frames the offer as agent-iterates/you-steer and does NOT lead with "explicit approval gates"; the SAME render over a manual-mode fixture KEEPS the gate-and-drive framing (both observed; the two offers differ — pinned by the live drive over a both-modes fixture). (b) [conditional] a fixture where product code lands via merged PRs while scaffolding is edited directly on a branch renders a work-by-area that detects the branch-and-merge workflow and caveats the edit-count signal so it does NOT falsely conclude "scaffolding > product" — observed in the render. R6#2 is the lowest-priority arm (did not bite two of h5's runs); if a branch-merge signal proves uncomputable from the agentsview corpus during implementation, it degrades to the caveat-only path (detect branch-and-merge → flag the work-by-area as branch-skewed) rather than a full diffstat re-attribution, recorded as the fallback.

## Test plan (consolidated — one unified fixture suite + one doc diff)

**One fixture suite.** Extend the existing survey golden-fixture harness (`skills/integration/survey_queries_test.go` + `testdata/survey/fixture-sessions.sql`) rather than building a new one — it is the established survey-feedback proof shape, runs under `go test ./skills/integration/ -run TestSurveyQuerySmoke`, sqlite3-driven, no live deps. Fixture additions, each band's rows seeded so the expected value is an independent oracle:

- **R3 (dispatch-fact):** `za`'s fixture work — add `parent_session_id` + `relationship_type` columns (production-shaped), seed two in-repo parents dispatching subagents (one →2, one →1; distinct=2, total=3), one subagent of the out-of-repo parent E (must NOT count), one `%/subagents/%` file_path row. Non-vacuous re-point. Add `dispatch-fact` to the required-queries list.
- **R5 (knowledge-work track):** seed a track with a process/file/log loop + content/ops `.md`-and-data edits + no issue→PR signature + no veto-heavy path → `mode-classification` emits `knowledge-work`. Non-vacuous: strip the loop markers → drops to `unlabeled`.
- **R6#1 vocab + mode:** the existing `issue-feed` track's `mode-classification` output asserts `manual` (was `mechanical`); the `landing-copy` track stays `exploration`. (Render-side offer framing is the live drive.)
- **R4 all-unlabeled Codex:** seed a Codex-scoped set whose `first_message`s are all encouragement/meta → `codex-workstreams` returns only `(unlabeled)` rows; plus a contrast fixture with ≥2 named clusters (already present: `journey-cost-ledger`, `orient-workflow-discovery`, `codex-live-ci`) so the collapse is conditional.
- **R2 (5x's fixtures):** the date-spanning fixture (span → `{N}` days) + the off-corpus-gap fixture (blank-cwd + name-only Codex superset) from `5x`'s firmed test plan.

**Render-level ACs (R1, R2 framing, R4 scratch, R6 offers)** are proven by a **live survey drive** over the fixture corpus (and a smoke over this repo, which qualifies: 66 in-repo sessions, 31 orchestrating parents, multiple tracks), mirroring the sprint's existing survey live-drive ACs — the rendered output carries the lede/`BY THE NUMBERS`/`THREADS TO PULL`/`WHAT THIS CAN'T SEE` with fixture-derived numbers and no scratch preamble. Per the survey discipline, a grep over SKILL.md never satisfies any AC; every expected value comes from the fixture/corpus session rows.

**One doc diff (all SKILL.md + queries.sql edits, applied by implementation):**
1. `references/queries.sql`: append `za`'s `dispatch-fact` query (verbatim from `za`'s doc diff); rename `mode-classification`'s `'mechanical'` → `'manual'` label; add the `knowledge-work` classification branch (process/file/log loop + content-edit signature, scored with the existing margin guard, beside mechanical/exploration).
2. `SKILL.md` step 2: add `run_query dispatch-fact` + its prose paragraph (from `za`); add the knowledge-work archetype to the `mode-classification` prose; update the `manual` vocabulary in the signal-accounting prose.
3. `SKILL.md` step 4: rewrite the report template to the value-&-numbers-first spine above — `WHAT THIS GIVES YOU` lede, `BY THE NUMBERS` block (interruptions, hanging-threads count, no-follow-up count, sessions-read, `za` dispatch line), `HOW YOU WORK` de-jargoned line, the `↓ full analysis` pointer, then the demoted detail (`THREADS TO PULL` reframe of the frontier, `BACKLOG`, `RECENT DECISIONS`, the existing body sections, `WHAT THIS CAN'T SEE` from `5x`); the conditional all-unlabeled Codex collapse; the explicit "render the report directly, no `I have everything I need` scratch preamble" instruction; the mode-aware commission offer (exploration → iterate/steer, manual/mechanical → keep gates).
4. `SKILL.md` synthesis guidance: update `mechanical → manual` vocabulary; add the knowledge-work mode; the branch-aware work-by-area caveat (R6#2 conditional).

Cost/complexity: medium. The query changes are low (one append verbatim + one label rename + one classification branch). The bulk is the step-4 template rewrite (the value-&-numbers-first render), proven by the live drive. No new agentsview ingestion, no new mechanism beyond `za`'s already-spiked query.

## Out of scope

- Landing-page CTA + bundling the agent-conversation-index dependency (separate, non-survey tasks — see 0202 cross-cutting).
- Generalizing the decision-frontier / partial-lens beyond the survey (0.21.x decision-abstraction).
- Reproducing any surveyed corpus content (anonymization discipline).

(The firmed acceptance criteria and test plan are above, in `## Acceptance criteria (consolidated — one per band, R1–R6)` and `## Test plan (consolidated …)`. The original sketch ACs/test plan were superseded by ideation.)

## Stage Report: ideation

- DONE: Produce ONE coherent design covering all six bands (R1–R6) rendered into the captain-locked value-&-numbers-first structure, verifying each band against the CURRENT skill first
  `## Verification against the current skill` records the per-band shipped/not-shipped finding (zb#1 SHIPPED #335 → R5 narrows to the archetype; zw#1 narrowed exploration-only; zw#2 conditional; R2/R3 adopt 5x/za's firmed ACs+doc diffs verbatim; R1/R4 all-new). `## Chosen design` maps every band onto the locked spine without redesigning it.
- DONE: Run the riskiest mechanism FIRST — the lede's concrete figures are computable from existing session data
  `## Riskiest-mechanism exercise` ran all four BY THE NUMBERS figures on this repo's live agentsview DB (66 in-repo sessions): interruptions 114 decisions + 19 vetoes, hanging-threads = decision-open OPEN count (10), no-follow-up = BACKLOG subset, dispatch 31/1675. Key finding: the three frontier figures are COUNTS of what decision-open + the cross-check already produce — no new query; the only new query is za's already-spiked dispatch-fact. Recorded "no spike needed beyond za's dispatch-fact spike".
- DONE: Produce the consolidated AC set (one per band), each proven by a fixture render with the expected value derived from fixture rows; one unified test plan; one doc diff
  Six consolidated ACs (AC-1…AC-6) each bound to an independent fixture-row oracle (never a SKILL.md prose-grep); query-level ACs pinned by TestSurveyQuerySmoke, render-level by a live drive. One fixture suite extends the existing survey golden harness; one 4-part doc diff (queries.sql + SKILL.md steps 2/4/synthesis).

### Summary
Folded the six feedback seeds into one value-&-numbers-first redesign of the step-4 survey report template. The load-bearing finding from the riskiest-mechanism exercise: the lede's three decision-frontier figures (hanging threads, no-follow-up decisions, interruptions) are derived COUNTS of sections the existing `decision-open` query + step-4 cross-check already produce — the redesign surfaces them, it computes no new signal, and the only new query is `za`'s already-spiked `dispatch-fact` (re-confirmed 31/1675 on live data). R2 and R3 adopt `5x`'s and `za`'s already-firmed ACs and doc diffs verbatim; R5 narrows to just the knowledge-work archetype (zb#1's Codex-count-first is SHIPPED in #335); R6#2 (branch-aware work-by-area) is the lowest-priority conditional arm with a caveat-only fallback.
