---
id: zajryfzgmnbzb5vjw5hven0z
title: Survey reports the FACT of subagent dispatch (so orchestrated repos don't look idle)
status: backlog
source: "captain (2026-06-08) — follow-up surfaced during 47rx's F spike. The survey scopes to the Claude session set and EXCLUDES dispatched-subagent sessions, so in a spacedock/agent-orchestrated repo most work is invisible (e.g. 7468/7680 internal edits here landed in excluded subagent sessions) → the survey under-reports and the repo looks idle. Captain scoping: we do NOT need subagent CONTENT for now — just the FACT that the session dispatched subagents."
started: 2026-06-13T04:49:28Z
completed:
verdict:
score:
worktree:
issue:
group: survey
sprint-readiness: defer
sprint: 0202-survey-improvements
---

When a session dispatched subagents, the survey should say so — surface the fact of orchestration so a spacedock/agent-orchestrated repo isn't reported as having done little.

## Problem

The survey's body (work-by-area, workstreams, interruptions) is built from the parent Claude sessions and EXCLUDES dispatched-subagent sessions. In a repo where the real work happens inside dispatched subagents (a spacedock-orchestrated repo), that work is invisible — the survey under-reports and the project reads as idle. 47rx's F spike measured this directly here: 7468/7680 internal edits were in excluded subagent sessions.

The fix is NOT to fold in subagent CONTENT (their edits/work-by-area) — that's deliberately out of scope. It is to report the FACT: this session dispatched subagents (orchestration happened), so the reader understands the parent-session view is only part of the picture.

## Proposed approach

A new labeled query `dispatch-fact` in `skills/survey/references/queries.sql` counts, over the same repo-scope the body already uses, two numbers from the session rows agentsview already has:

- `sessions_that_orchestrated` — the count of DISTINCT in-repo parent sessions that dispatched at least one subagent;
- `subagents_dispatched` — the total subagent sessions those parents dispatched.

The query joins each subagent row (`relationship_type = 'subagent'`) back to its parent via `parent_session_id`, and keeps the join only where the PARENT is an in-repo, non-subagent Claude session (the body's exact scope: `agent='claude' AND file_path NOT LIKE '%/subagents/%' AND (cwd = :repo_root OR cwd LIKE :repo_root || '/%')`). The report renders one BY THE NUMBERS line from these two figures, demoting the subagent-content exclusion into the partial-lens caveat. Subagent CONTENT (their edits / work-by-area / decisions) stays excluded everywhere — this surfaces the dispatch FACT only.

### Marker decision (the riskiest mechanism — spiked, see below)

The dispatch marker is the **parent-join** (`relationship_type='subagent'` + `parent_session_id` → repo-scoped parent), NOT the `Task` tool-call's `subagent_session_id`. Spiking the live DB showed `subagent_session_id` is unpopulated on parent `Task` rows in agentsview v0.32.1 (zero matches), while all 1673 subagent rows resolve to a present parent via `parent_session_id`. The parent-join also yields the distinct-orchestrating-parent count the assignment names, and does not depend on the subagent carrying a usable `cwd` (agentsview does not guarantee one — Codex subagents land `cwd=''`), so it scopes correctly through the parent regardless.

## Riskiest-mechanism spike (DONE — exercised on the live DB)

The riskiest unknown was: *is "this parent session dispatched subagents" detectable from the session data the survey already reads?* Exercised against this repo's synced agentsview DB (v0.32.1), scoped to `/Users/clkao/git/spacedock-research/spacedock-v1`:

- `relationship_type` distinct values: `'' → 67`, `'subagent' → 1673`. The 1673 subagent rows are exactly the set `file_path LIKE '%/subagents/%'` already filters out everywhere in `queries.sql`.
- All 1673 subagent rows resolve to a PRESENT parent row via `parent_session_id` (0 missing).
- Joined to repo-scoped non-subagent Claude parents: **31 distinct orchestrating parents, 1673 subagents dispatched** — out of 66 total in-repo parent sessions (the body's set). So ~47% of in-repo sessions orchestrated, yet the body shows none of that work.
- The `Task` tool-call's `subagent_session_id` column is UNPOPULATED on parent rows here (0 matches) — that candidate marker is a dead end in this agentsview version.
- Subagents in this corpus DO carry a `cwd` (1673/1673), so a direct subagent-cwd-prefix scope and the parent-join scope AGREE here (both 1673). The parent-join is chosen as primary because it (a) yields the distinct-parent count, (b) ties to the body's exact parent scope, and (c) survives a corpus where subagent `cwd` is blank (which agentsview permits).

Marker: `sessions.relationship_type = 'subagent'` + `sessions.parent_session_id`. No further spike needed.

## Render location (banked "value & numbers first" structure)

The dispatch fact is a value-and-numbers figure, so it renders inside **BY THE NUMBERS** (docs/roadmap/0202-survey-improvements/index.md mock), not in the demoted detail section. One added line, sourced from the two query figures, distinct from the parent-session body counts:

```
 31  sessions dispatched subagents (1673 dispatched — their work isn't shown here)
```

The "their work isn't shown" clause is the lens-completeness pointer that pairs with `5x`'s single partial-lens statement; it keeps the subagent-CONTENT exclusion honest without folding content in. The line is dropped entirely when `sessions_that_orchestrated = 0` (a non-orchestrated repo renders no dispatch line — same fill-or-drop rule as the rest of the report).

## Out of scope

- Subagent session CONTENT (edits, work-by-area, decisions) — explicitly deferred; this surfaces only the dispatch fact.
- Changing the Claude-session scoping of the body itself.
- The `Task`-tool `subagent_session_id` marker — unpopulated in agentsview v0.32.1 (see spike); not used.

## Acceptance criteria

- **AC-1 — the `dispatch-fact` query counts orchestration over the body's scope.** `skills/survey/references/queries.sql` carries a `-- name: dispatch-fact` labeled query that, run over a fixture seeded with dispatched-subagent sessions, returns `sessions_that_orchestrated | subagents_dispatched` counting only subagents whose PARENT is an in-repo, non-subagent Claude session.
  Verified by: `skills/integration/survey_queries_test.go` — a new `dispatch-fact` sub-test runs the labeled query against `testdata/survey/fixture-sessions.sql` and asserts both counts. The expected values are DERIVED FROM the fixture's seeded session rows (an independent source that diverges from SKILL.md prose); a broken join, a dropped parent-scope filter, or schema drift reds the check.
- **AC-2 — a subagent of an OUT-of-repo parent does not inflate the count.** A subagent whose parent session is outside the repo prefix is excluded from both counts.
  Verified by: the fixture seeds one such subagent; the sub-test asserts it is not counted, plus a non-vacuous mutation (re-point that subagent's parent to an in-repo parent) that FLIPS the counts — proving the parent-scope is load-bearing, not a constant.
- **AC-3 — the rendered survey names the dispatch fact in BY THE NUMBERS.** A live survey drive over a corpus/fixture with dispatched-subagent sessions renders the dispatch-fact line in the value-and-numbers block, with the count matching the query, distinct from the parent-session body counts.
  Verified by: a live-drive of the survey (the AC-3 live drive, mirroring the existing survey live-drive ACs) — the rendered output carries the dispatch line and its number equals the `dispatch-fact` query's `sessions_that_orchestrated` over the same corpus; a non-orchestrated corpus renders NO dispatch line.

## Test plan

- **Query-smoke (AC-1, AC-2)** — extend `skills/integration/survey_queries_test.go` with a `dispatch-fact` sub-test (the harness already extracts a labeled query from `queries.sql` and runs it against `testdata/survey/fixture-sessions.sql`). Fixture work: add `parent_session_id` + `relationship_type` columns to the fixture `sessions` table (production-shaped — currently trimmed), and seed (i) two in-repo parents that dispatch subagents (one dispatches 2, one dispatches 1 → distinct-parents=2, subagents=3), (ii) one subagent whose parent is the OUT-of-repo session E (must NOT count), (iii) a subagent file_path under `%/subagents/%` matching production. Add the `dispatch-fact` name to the existing required-queries list. Non-vacuous mutation: re-point the out-of-repo subagent's `parent_session_id` to an in-repo parent and assert the counts climb (2→3 subagents, distinct-parents 2→3 if a new parent). Cost: low — pure fixture + one Go sub-test, no live deps; runs under `go test ./skills/integration/ -run TestSurveyQuerySmoke`.
- **Live-drive (AC-3)** — run `spacedock:survey` over a real orchestrated corpus (this repo qualifies: 31 orchestrating parents) and confirm the rendered BY THE NUMBERS block names the dispatch line with the count from the query; confirm a non-orchestrated corpus renders no dispatch line. Cost: medium (manual live run), mirrors the sprint's existing live-drive ACs.
- Per the survey discipline, a grep over SKILL.md never satisfies any AC here — every expected number comes from the fixture/corpus session rows, an independent source.

## Doc diff (applied at implementation)

**`skills/survey/references/queries.sql`** — append a labeled query (after `mode-classification`):

```sql
-- name: dispatch-fact
-- za — the FACT of subagent dispatch, so an orchestrated repo isn't read as idle. The
-- body EXCLUDES subagent sessions everywhere (`file_path NOT LIKE '%/subagents/%'`), so a
-- repo whose real work lands in dispatched subagents shows an almost-empty body. This
-- surfaces only the FACT (a count), never subagent CONTENT. Marker: a subagent row is
-- `relationship_type = 'subagent'` and links to its parent via `parent_session_id`
-- (agentsview v0.32.1; the Task-tool `subagent_session_id` is unpopulated). Scope is the
-- body's exact parent scope — count a subagent ONLY when its PARENT is an in-repo,
-- non-subagent Claude session — so a subagent of an out-of-repo parent stays out.
-- `sessions_that_orchestrated` is the DISTINCT in-repo parent count (the orchestration
-- fact); `subagents_dispatched` is the total. The report renders one BY THE NUMBERS line
-- from these and drops it when sessions_that_orchestrated = 0.
SELECT
  COUNT(DISTINCT p.id) AS sessions_that_orchestrated,
  COUNT(*)             AS subagents_dispatched
FROM sessions sub
JOIN sessions p ON p.id = sub.parent_session_id
WHERE sub.relationship_type = 'subagent'
  AND p.agent = 'claude'
  AND p.file_path NOT LIKE '%/subagents/%'
  AND (p.cwd = :repo_root OR p.cwd LIKE :repo_root || '/%');
```

**`skills/survey/SKILL.md`** — two edits:

1. In step 2 (the `run_query` list, after `mode-classification`), add:
   ```
   run_query dispatch-fact      # za — orchestration FACT: distinct in-repo parents that dispatched subagents + total dispatched
   ```
   and a short prose paragraph (sibling to "Track modes" / "Honest signal accounting"):
   > **Dispatch fact (`dispatch-fact`).** The body EXCLUDES dispatched-subagent sessions, so an orchestrated repo (most work inside subagents) reads as nearly idle. `dispatch-fact` counts the DISTINCT in-repo parent sessions that dispatched subagents and the total dispatched, joining each subagent (`relationship_type='subagent'`) to its parent and keeping only in-repo parents. It surfaces the FACT of orchestration, never subagent CONTENT.

2. In step 4's `BY THE NUMBERS` block of the rendered template (the value-and-numbers fence), add a fill line that renders only when `sessions_that_orchestrated > 0`:
   ```
   {if sessions_that_orchestrated>0: {sessions_that_orchestrated}  sessions dispatched subagents ({subagents_dispatched} dispatched — their work isn't shown here)}
   ```

The before/after wording above is the concrete diff; implementation applies it. Exact line placement is left to implementation against the then-current files, but the added query, the `run_query` line, the prose paragraph, and the conditional BY THE NUMBERS fill line are the four required changes.

## Stage Report: ideation

- DONE: Firm where the FACT of subagent dispatch renders in the banked "value & numbers first" structure so an orchestrated repo isn't read as idle — a count of dispatched subagents / sessions-that-orchestrated, distinct from the parent-session body; subagent CONTENT stays excluded.
  Render location firmed as a BY THE NUMBERS line (see "Render location" section): `31  sessions dispatched subagents (1673 dispatched — their work isn't shown here)`, dropped when count=0; subagent CONTENT remains excluded everywhere.
- DONE: Riskiest-mechanism check FIRST — that "this parent session dispatched subagents" is detectable from the session data already read (a dispatch marker), exercised on a real session.
  Spiked against this repo's live agentsview DB (v0.32.1): marker is `relationship_type='subagent'` + `parent_session_id`; 31 distinct in-repo orchestrating parents / 1673 subagents over the body's scope. The `Task`-tool `subagent_session_id` candidate is unpopulated (dead end). See "Riskiest-mechanism spike".
- DONE: Propose the doc diff and a test plan whose AC proof is a survey run over a fixture with dispatched-subagent sessions rendering the dispatch fact, the count DERIVED FROM the fixture session rows (independent source) — never a prose-grep.
  Doc diff (queries.sql `dispatch-fact` query + SKILL.md run_query line, prose paragraph, conditional BY THE NUMBERS line) recorded; AC-1/AC-2 proven by a new `dispatch-fact` sub-test in the existing `skills/integration/survey_queries_test.go` harness over `testdata/survey/fixture-sessions.sql` (counts derived from seeded rows + a non-vacuous parent-repoint mutation); AC-3 by a live survey drive.

### Summary

Spiked the dispatch marker on the live DB and proved the FACT of orchestration is detectable from data the survey already reads: 31 in-repo parent sessions dispatched 1673 subagents here, all reachable via `relationship_type='subagent'` + `parent_session_id`, while the `Task`-tool `subagent_session_id` column is unpopulated in agentsview v0.32.1. Firmed the render as a single BY THE NUMBERS line (value-and-numbers first, drop-when-zero), chose the parent-join scope over a direct subagent-cwd scope because it yields the distinct-parent count and survives blank-subagent-cwd corpora, and wrote the concrete doc diff (a new `dispatch-fact` query + three SKILL.md edits) plus ACs proven by the existing fixture-driven query-smoke harness — expected counts derived from seeded fixture rows, never a SKILL.md grep.
