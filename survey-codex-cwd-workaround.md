---
id: 69rk6t1vbehsd4fwxnsnjwma
title: survey — work around agentsview not persisting Codex cwd (scoping fallback + hint), do NOT fix agentsview
status: done
source: "captain (2026-06-08) — agentsview ingests Codex sessions and derives project but does NOT persist Codex cwd into sessions.cwd, so the survey's cwd-scoped repo-identity query misses all Codex sessions. Decision: work around it in OUR survey skill with a fallback/hint, not fix agentsview."
score: "0.28"
started: 2026-06-08T15:29:12Z
completed:
verdict: superseded
worktree:
issue:
sprint: 0198-pre-flip-hardening
group: survey
sprint-readiness: defer
superseded-by: survey-skill-correctness-pass
---

The survey's repo-identity scoping (the `git rev-parse --git-common-dir` coalesce, xn) scopes by `sessions.cwd`. agentsview does not persist Codex `cwd`, so cwd-scoped queries miss all Codex sessions even if the Claude-only `agent` filter were broadened. Work around it survey-side with a presence count and a hint; do not fix agentsview.

## Problem

- agentsview ingests Codex sessions, derives `project`, but leaves `sessions.cwd` blank for Codex.
- Every query in `references/queries.sql` scopes by `cwd = :repo_root OR cwd LIKE :repo_root || '/%'`, so a blank-cwd Codex session can never match the repo scope — it is invisible to the survey on the cwd axis. (Separately, the queries also filter `agent='claude'`; surfacing Codex *decision/scaffold/work-by-area* signals is a deferred follow-up — SKILL.md §13, §96, §209 — and is OUT of scope here.)
- Matching Codex by `project` instead catches them but loses repo-root identity: `project` is a basename-style key that demonstrably collides across unrelated repos (spike below), so a silent `project` union would fold a *different* repo's Codex sessions into this repo's survey.

## Spike: ground the premise (DONE — see Stage Report)

Ran the riskiest unknown first against real agentsview v0.32.1 data, driven through the binary (raw `~/.agentsview/sessions.db` reads are TCC-denied to this process, exactly as SKILL.md §21 warns). Ingested real Codex rollouts from `~/.codex/sessions` into a process-readable `AGENTSVIEW_DATA_DIR` via `agentsview session sync <path>`, then queried that copy. Findings:

1. **Codex `cwd` is blank, `project` is populated.** 13/13 ingested Codex sessions: `cwd = ''` (empty string, NOT NULL — schema is `cwd TEXT NOT NULL DEFAULT ''`), `project` set (e.g. `spacedock_v1`, `orb_backtest`, `spacedock_gym`). The repo root appears in the session's first-message *text* but agentsview does not lift it into `cwd`. Premise confirmed.
2. **Claude `cwd` is populated.** 10115/10170 Claude sessions (99.5%) carry a non-blank `cwd`. So this is Codex-specific, not a general blank-cwd problem — the cwd-scoped queries work for Claude and only drop Codex.
3. **`project` collides across repos (bounds the fallback risk).** Distinct repo roots sharing one `project` key, from real data: `workspace` → 6 different roots (e.g. six `dataagentbench/_runs/*/…/workspace`), `spacedock` → spans `/Users/clkao/git/spacedock`, `/Users/clkao/git/spacedock-prompt/spacedock`, and `/Users/clkao/git/spacedock-prompt`. A `project`-only Codex union for a repo whose key is `workspace` or `spacedock` would pull in unrelated repos' Codex sessions. This is the quantified basename-collision risk.
4. **agentsview derives `project` from the git-root basename, not cwd basename.** For this repo, the root cwd, every `.worktrees/*` cwd, and the `docs/dev/.spacedock-state` subdir all key to ONE `project='spacedock_v1'`. (This refines SKILL.md §64 / the current fixture, which assume per-cwd-basename divergent keys like `_spacedock_state`. Not this task's deliverable, but the build worker should not assume divergent keys for the Codex fallback — Codex's `project` equals the repo-root basename, so it matches the same `project` Claude derives.)

The project-key normalization observed: repo basename `spacedock-v1` → key `spacedock_v1` (non-alphanumerics, here `-`, become `_`).

## Decision: project-fallback presence count + hint (NOT a silent scope union)

The blend, decided from the spike:

- **DO add a dedicated `codex-presence` query** (a new labeled query in `references/queries.sql`) that counts Codex sessions whose `project` equals the repo-root project key, AND reports `blank_cwd` among them. It is a SEPARATE count, not a union into the Claude `scoping`/decision/scaffold sets — so it cannot silently inflate or contaminate the existing Claude-scoped numbers. The repo-root project key is derived in SKILL.md prose from the resolved `REPO_ROOT` basename with the same `-`/non-alphanumeric → `_` normalization agentsview applies (bound as a new `:repo_project` parameter).
- **DO surface a hint in the report** whenever the `codex-presence` count is non-zero: a line stating that N Codex sessions matched this repo *by project name only* (cwd unrecorded by agentsview), so the reader knows (a) Codex work exists here and the Claude-only survey body does not cover it, and (b) the match is basename-grade, not repo-root-grade — it may include a same-named sibling repo. The hint is the honest accounting SKILL.md §98/§208 already demands for blank-cwd sessions, extended to Codex.
- **DO NOT union Codex `project` rows into the cwd-scoped Claude query set.** The collision evidence (spike finding 3) makes a silent union unsafe: it would mis-attribute another repo's Codex sessions as this repo's. Presence-count-plus-hint gives the reader the signal without the false precision.

Rationale: the task asks the survey to *account for* Codex sessions when cwd is absent. A flagged presence count does exactly that, honestly, without reintroducing the same-basename-collision bug the cwd-scoping was built to avoid. Promoting Codex into the full survey body (decisions, scaffold, work-by-area) stays the documented deferred follow-up.

## Acceptance criteria

- **AC-1 — `codex-presence` query exists and counts project-matched Codex sessions, reporting blank-cwd.** The reference file `skills/survey/references/queries.sql` contains a labeled `-- name: codex-presence` query that, bound to a `:repo_project` parameter, returns the count of `agent='codex'` sessions whose `project = :repo_project` and the count of those with blank `cwd`. Verified by: the query-smoke test (AC-2) extracts and runs this labeled query and asserts the counts against the fixture — fails if the query is missing, renamed, or returns the wrong shape.

- **AC-2 — query-smoke fixture proves project-matched Codex rows are counted and a same-`project` sibling under a different root is NOT silently folded into the Claude scope.** `skills/integration/testdata/survey/fixture-sessions.sql` carries Codex rows with the production-accurate shape (`project` set, `cwd = ''`), at least one matching `:repo_project` for the in-repo case. `survey_queries_test.go` asserts: (a) `codex-presence` counts the matching Codex rows with blank_cwd > 0; (b) the existing `scoping` query's `sessions` count is UNCHANGED by the added Codex rows (Codex stays out of the Claude-scoped count — no silent union). Expected values come from the fixture rows (an independent source), not from skill prose. Verified by: `go test ./skills/integration/ -run TestSurveyQuerySmoke` passes.

- **AC-3 — the survey report renders a Codex-presence hint when the count is non-zero, and omits it when zero.** SKILL.md step 2/4 prose derives `:repo_project` from `REPO_ROOT`, runs `codex-presence`, and the report template gains a hint line shown only when the count > 0, stating the count and that the match is by project name only (cwd unrecorded). Verified by: a live drive of the survey skill on a repo with project-matched blank-cwd Codex sessions present in the survey DB, observing the rendered hint line in the report; and a drive where no such Codex rows exist, observing the hint absent. (A string/grep over SKILL.md does NOT satisfy this — only invoking the skill and observing the rendered report does.)

## Test plan

- **Mechanism already de-risked (this ideation's spike).** The on-disk shape the design rests on — Codex `cwd` blank, `project` populated, `project` derived from git-root basename, `project` collides across repos — is confirmed against real agentsview v0.32.1 (Stage Report below). No further spike needed before build; the build's first test is the fixture row AC-2 seeds from this shape.
- **AC-1 + AC-2 — query-smoke (fixture, cheap, ~seconds).** Add Codex rows (`project` set, `cwd=''`) to the committed fixture and two assertions to `TestSurveyQuerySmoke`: `codex-presence` returns the expected matched/blank counts, and `scoping.sessions` is unchanged versus the pre-Codex-row baseline. Reuses the existing sqlite3-driven harness. Production-accurate shape matters: the current fixture's lone Codex row has `cwd` SET (seeded only to prove Claude-scope exclusion) — the build must add blank-cwd Codex rows rather than rely on that one, and should keep an out-of-repo / colliding-`project` consideration in mind so the count is `project`-scoped, not global.
- **AC-3 — live workflow drive (the only proof for the rendered hint).** Drive the survey skill end-to-end against a survey DB seeded with project-matched blank-cwd Codex sessions (the spike's `agentsview session sync` into a readable `AGENTSVIEW_DATA_DIR` is the recipe) and observe the hint in the rendered report; repeat against a Codex-free DB and observe its absence. Cost: minutes; needs the agentsview binary (present) and the readable-data-dir workaround (TCC-safe path proven in the spike).
- Estimated total cost/complexity: LOW. One new labeled query, a few fixture rows, two test assertions, a prose+template hint, one live drive. No binary/Go changes beyond the test file; no agentsview changes (explicitly out of scope).

## Out of scope

- Fixing agentsview to persist Codex `cwd` (upstream `kenn-io/agentsview`) — explicitly NOT this task.
- Surfacing Codex *decision*, *scaffold-usage*, or *work-by-area* signals (broadening the `agent='claude'` filter across the body) — the documented deferred follow-up; this task only adds a presence count + hint.

## Notes

Survey-skill workaround; same area as xn. Root cause is upstream agentsview (`sessions.cwd` blank for Codex); we deliberately work around it rather than block on the upstream fix. The spike also surfaced that agentsview v0.32.1 keys `project` by git-root basename (not cwd basename), which diverges from the current fixture/SKILL.md model — flagged for the build worker and a possible separate cleanup, but not changed here. Candidate sprint: 0198-pre-flip-hardening (current) / survey-followup.

## Stage Report: ideation

- DONE: Ground the premise (spike): confirm agentsview leaves sessions.cwd null/empty for Codex sessions — inspect a real sessions.db or a representative fixture; record the finding.
  Ingested 13 real Codex rollouts from `~/.codex/sessions` via `agentsview session sync` (v0.32.1) into a readable AGENTSVIEW_DATA_DIR (raw `~/.agentsview` read is TCC-denied); 13/13 had `cwd=''` with `project` populated, while 10115/10170 Claude sessions had cwd set. Premise confirmed; recorded in "Spike" section.
- DONE: Decide the workaround: pick the blend (project-fallback union and/or hint) and bound the basename-collision risk the fallback reintroduces.
  Chose a dedicated `codex-presence` count + report hint, explicitly NOT a silent scope union — bounded by real collision data (`workspace`→6 roots, `spacedock`→3 roots). Recorded in "Decision" section.
- DONE: Produce build-ready ACs + test plan: a query-smoke fixture with Codex rows (project set, cwd null) asserting they are included/flagged, plus the live-drive hint when applicable.
  AC-1/AC-2 (query-smoke over production-accurate `cwd=''` Codex fixture rows; `scoping` count unchanged) + AC-3 (live-drive rendered hint). Recorded in "Acceptance criteria" and "Test plan".

### Summary

Spiked the premise against real agentsview v0.32.1 data driven through the binary (TCC blocks raw DB reads): Codex sessions land with blank `cwd` and a populated, git-root-basename `project` key that demonstrably collides across unrelated repos. That collision evidence drove the decision against a silent `project` union — the design is a separate flagged `codex-presence` count plus a report hint, so Codex presence is surfaced honestly without contaminating the Claude-scoped numbers. A side finding (agentsview keys `project` by git-root basename, not cwd basename — diverging from the current fixture/SKILL.md model) is flagged for the build worker but left unchanged here.
