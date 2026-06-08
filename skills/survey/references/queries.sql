-- ABOUTME: Recommended read-only survey queries over agentsview's sessions.db —
-- ABOUTME: one labeled query per concern, run by SKILL.md prose against the synced copy.
--
-- These are the runnable surface for the survey skill: SKILL.md prose resolves the
-- repo root, binds the named parameters below, runs each labeled query against the
-- agentsview sessions.db copy, and assembles the report from the rows. The SQL lives
-- ONLY here — never fenced in SKILL.md — so the smoke test EXECUTES this file (it runs
-- the artifact) rather than parsing the instruction prose.
--
-- Named parameters the caller (SKILL.md prose, or the smoke test) must bind before
-- running a query:
--   :repo_root  — the resolved repo-root path, NO trailing slash, from
--                 `dirname(git rev-parse --path-format=absolute --git-common-dir)`
--                 (fallback pwd when not in a git repo / git <2.31). agentsview derives
--                 `project` from the git-root basename, so EVERY checkout of one repo —
--                 the root, a subdir, a worktree, the split-root state dir — already
--                 shares ONE `project` key; the risk is the INVERSE, that this basename
--                 key COLLIDES with a same-basename sibling repo elsewhere on disk. The
--                 boundary match is `cwd = :repo_root OR cwd LIKE :repo_root || '/%'`:
--                 it admits the cwd AT the repo root and every path strictly under it,
--                 while a same-basename sibling like `/other/proj` does NOT match (a bare
--                 `LIKE :repo_root || '%'` would wrongly fold the sibling in). So the
--                 absolute cwd-prefix is the ONLY column that keeps a same-basename
--                 sibling's sessions out of this repo's scope.
--   :repo_project — the git-root basename agentsview keys `project` by, normalized
--                 (non-alphanumerics → `_`, e.g. `spacedock-v1` → `spacedock_v1`). Used
--                 ONLY by the codex-presence query: Codex sessions land cwd='' (agentsview
--                 does not persist Codex cwd), so they cannot be cwd-prefix-scoped; they
--                 are matched by `project = :repo_project` alone, which also catches a
--                 same-basename sibling's Codex sessions (hence a flagged presence count,
--                 never a union into the cwd-scoped Claude sets).
--
-- The session-scoping subquery is repeated inline in each query (rather than a shared
-- view) so each labeled query is independently runnable: the smoke test extracts ONE
-- query block and runs it in isolation.

-- name: scoping
-- #318 — git-root identity scoping. Count the Claude sessions whose cwd is under the
-- repo-root prefix. agentsview keys `project` by the git-root basename, so every in-repo
-- checkout (root, subdir, worktree, state dir) already shares ONE key — the cwd-prefix's
-- job is the INVERSE of a merge: it EXCLUDES a same-basename SIBLING repo whose sessions
-- carry the identical basename-derived `project`. The absolute cwd-prefix is the only
-- column that distinguishes this repo's checkouts (all under the repo-root path) from a
-- same-named repo elsewhere on disk. Surfaces the blank-cwd count (sessions agentsview
-- captured no cwd for, which the repo-root scope cannot place).
SELECT
  COUNT(*) AS sessions,
  COALESCE(SUM(cwd = ''), 0) AS blank_cwd,
  substr(MIN(started_at), 1, 10) || ' .. ' || substr(MAX(ended_at), 1, 10) AS span
FROM sessions
WHERE agent = 'claude'
  AND file_path NOT LIKE '%/subagents/%'
  AND (cwd = :repo_root OR cwd LIKE :repo_root || '/%');

-- name: codex-presence
-- #69 — flagged Codex presence (NOT a union into the Claude scope). agentsview does not
-- persist Codex cwd (every Codex session lands cwd=''), so the cwd-prefix scope above
-- misses them entirely. Match by `project = :repo_project` alone — the git-root basename.
-- That basename can COLLIDE with a same-basename sibling repo (the same reason the Claude
-- scope needs the absolute cwd-prefix), so this count MAY include a same-named sibling's
-- Codex sessions; the report states the match is by project NAME only. Returns the matched
-- Codex count and the blank-cwd sum among them; the report renders a hint ONLY when the
-- count is > 0. This is a presence flag — these sessions are NEVER folded into the
-- cwd-scoped Claude scoping/scaffold/work-by-area sets (collision risk).
SELECT
  COUNT(*) AS codex_sessions,
  COALESCE(SUM(cwd = ''), 0) AS blank_cwd
FROM sessions
WHERE agent = 'codex'
  AND project = :repo_project;

-- name: scaffold-usage
-- #319 / #317.1 — behavioral scaffold tally. GROUP the tool_calls.skill_name signal by
-- FAMILY over the repo-scoped Claude session set: split on the first ':' (namespaced
-- names like `superpowers:brainstorming` → family `superpowers`); bare names that are
-- superpowers disciplines (`brainstorming`, `running-research-spikes`, …) also map to
-- `superpowers`; anything else keeps its own name as the family. The `spacedock` family
-- is EXCLUDED — survey/ensign self-invocation dominates the tally (1026/1045 on this
-- repo) and would make every surveyed repo read as "uses spacedock". This is the
-- behavioral half the file-only probe misses: a family that appears here with no files
-- on disk is a `recovered` scaffold.
SELECT family, COUNT(*) AS invocations
FROM (
  SELECT
    CASE
      WHEN instr(skill_name, ':') > 0 THEN substr(skill_name, 1, instr(skill_name, ':') - 1)
      WHEN skill_name IN (
        'brainstorming', 'writing-plans', 'executing-plans',
        'subagent-driven-development', 'running-research-spikes',
        'dispatching-parallel-agents', 'systematic-debugging',
        'finishing-a-development-branch', 'verification-before-completion',
        'receiving-code-review', 'requesting-code-review', 'using-superpowers'
      ) THEN 'superpowers'
      ELSE skill_name
    END AS family
  FROM tool_calls t
  JOIN sessions s ON t.session_id = s.id
  WHERE s.agent = 'claude'
    AND s.file_path NOT LIKE '%/subagents/%'
    AND (s.cwd = :repo_root OR s.cwd LIKE :repo_root || '/%')
    AND t.tool_name = 'Skill'
    AND t.skill_name IS NOT NULL
    AND t.skill_name <> ''
)
WHERE family <> 'spacedock'
GROUP BY family
ORDER BY invocations DESC;

-- name: work-by-area
-- #317.2 — WORK-BY-AREA: identity from edits. Bucket Edit/Write target file_paths over
-- the repo-scoped Claude session set by the FIRST path segment under the repo root
-- (`internal`, `skills`, `docs`, …). A file_path that is NOT under the repo-root prefix
-- is an edit to an external sibling repo — a reference, not this project's identity — so
-- it buckets as `<external>`. "What this is" (where edits land) is reported separately
-- from "where you stop" (the decision frontier).
SELECT
  CASE
    WHEN fp LIKE :repo_root || '/%' THEN
      CASE
        WHEN instr(substr(fp, length(:repo_root) + 2), '/') > 0
          THEN substr(substr(fp, length(:repo_root) + 2), 1, instr(substr(fp, length(:repo_root) + 2), '/') - 1)
        ELSE substr(fp, length(:repo_root) + 2)
      END
    ELSE '<external>'
  END AS area,
  COUNT(*) AS edits
FROM (
  SELECT json_extract(t.input_json, '$.file_path') AS fp
  FROM tool_calls t
  JOIN sessions s ON t.session_id = s.id
  WHERE s.agent = 'claude'
    AND s.file_path NOT LIKE '%/subagents/%'
    AND (s.cwd = :repo_root OR s.cwd LIKE :repo_root || '/%')
    AND t.tool_name IN ('Edit', 'Write')
    AND json_extract(t.input_json, '$.file_path') IS NOT NULL
)
GROUP BY area
ORDER BY edits DESC;

-- name: decision-open
-- #320 — the decision/OPEN frontier (raw transcript scan; the repo cross-check is
-- SKILL.md prose). A decision is `done` ONLY when the result is an answered/approved
-- confirmation; every other result shape (rejection, error, prompt echo) is OPEN. The
-- ExitPlanMode approval result is `User has approved your plan` — folded into the
-- done-prefix set here so an approved plan reads `done`, not OPEN (the cheap #320 fix;
-- without it every approved plan fell to the frontier). OPEN sorts before done so the
-- load-bearing frontier cannot be truncated by the recency LIMIT.
SELECT
  COALESCE(json_extract(input_json, '$.questions[0].header'), 'PLAN') AS header,
  CASE
    WHEN result_content LIKE 'User has answered%'
      OR result_content LIKE 'Your questions have been answered%'
      OR result_content LIKE 'Your question has been answered%'
      OR result_content LIKE 'User has approved your plan%'
      OR result_content LIKE 'User approved%'
    THEN 'done'
    ELSE 'OPEN'
  END AS status,
  substr(replace(COALESCE(json_extract(input_json, '$.questions[0].question'), input_json), char(10), ' '), 1, 110) AS question
FROM tool_calls t
JOIN sessions s ON t.session_id = s.id
WHERE s.agent = 'claude'
  AND s.file_path NOT LIKE '%/subagents/%'
  AND (s.cwd = :repo_root OR s.cwd LIKE :repo_root || '/%')
  AND t.tool_name IN ('AskUserQuestion', 'ExitPlanMode')
ORDER BY status ASC, t.id DESC
LIMIT 20;
