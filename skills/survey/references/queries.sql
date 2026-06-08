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

-- name: codex-scoped
-- #321 — Codex attributed to THIS repo by exec_command.$.workdir prefix (DISTINCT from
-- codex-presence). agentsview persists no Codex session cwd, but Codex's `exec_command`
-- tool calls DO carry `$.workdir` — the absolute working directory of the shell command.
-- Scope a Codex session to this repo when it has at least one `exec_command` whose
-- `$.workdir` is `:repo_root` or strictly under it — the workdir analogue of the Claude
-- cwd-prefix scope. This ADMITS this repo's Codex and EXCLUDES a same-basename sibling
-- (whose workdirs fall under a different absolute prefix), so the count is sibling-free
-- where codex-presence (project NAME only) is not. The two stay distinct: codex-presence
-- is the name-matched superset (may include a sibling); codex-scoped is the workdir-attributed
-- subset. A Codex session whose only workdirs are outside the prefix (a /tmp run) or that
-- has no usable workdir (a near-empty session) is conservatively excluded — same discipline
-- as the Claude blank-cwd exclusion.
SELECT COUNT(*) AS codex_scoped_sessions
FROM sessions s
WHERE s.agent = 'codex'
  AND EXISTS (
    SELECT 1 FROM tool_calls t
    WHERE t.session_id = s.id
      AND t.tool_name = 'exec_command'
      AND (json_extract(t.input_json, '$.workdir') = :repo_root
           OR json_extract(t.input_json, '$.workdir') LIKE :repo_root || '/%')
  );

-- name: codex-workstreams
-- #322 — cluster the codex-scoped sessions into ensign-task WORKSTREAMS from each session's
-- first_message. This is the runnable artifact for the clustering rule (NOT prose): three
-- ordered cases, in order —
--   1. Dispatch pattern — first_message reads a dispatch file
--      (`…/spacedock-dispatch/spacedock-ensign-{TASK}-{stage}.md`); the workstream is {TASK},
--      with the trailing `-ideation`/`-implementation`/`-validation` stage stripped and the
--      label anchored on the `.md` boundary so trailing prose does not leak in.
--   2. Task/entity pattern — first_message names a `Spacedock task`/`Spacedock entity`; the
--      label is the backtick-quoted {TASK} that FOLLOWS that keyword (anchor on the keyword,
--      not the first backtick globally, so a leading reviewer-label backtick does not leak).
--   3. Unlabeled — a null first_message or an encouragement/meta message ("Captain asked me
--      to tell subagents…", "You totally got this…") carries no task; it falls to an honest
--      `(unlabeled)` bucket — counted, never invented a name for.
-- Same `EXISTS exec_command workdir-under-prefix` scope as codex-scoped (the clustering only
-- runs over the attributed set). `(unlabeled)` sorts last so the named tracks lead.
WITH scoped AS (
  SELECT s.id, s.first_message AS fm
  FROM sessions s
  WHERE s.agent = 'codex'
    AND EXISTS (
      SELECT 1 FROM tool_calls t
      WHERE t.session_id = s.id
        AND t.tool_name = 'exec_command'
        AND (json_extract(t.input_json, '$.workdir') = :repo_root
             OR json_extract(t.input_json, '$.workdir') LIKE :repo_root || '/%')
    )
),
dispatch_tail AS (  -- the text AFTER the dispatch-file marker, for the case-1 sessions
  SELECT id, fm,
    substr(fm, instr(fm, 'spacedock-dispatch/spacedock-ensign-')
               + length('spacedock-dispatch/spacedock-ensign-')) AS after_marker,
    substr(fm, instr(fm, 'Spacedock ')) AS after_kw  -- the text from `Spacedock task/entity …`
  FROM scoped
),
labeled AS (
  SELECT id,
    CASE
      WHEN fm LIKE '%spacedock-dispatch/spacedock-ensign-%.md%' THEN
        replace(replace(replace(
          substr(after_marker, 1, instr(after_marker, '.md') - 1),
        '-ideation', ''), '-implementation', ''), '-validation', '')
      WHEN fm LIKE '%Spacedock task `%' OR fm LIKE '%Spacedock entity %`%' THEN
        substr(
          substr(after_kw, instr(after_kw, '`') + 1),
          1,
          instr(substr(after_kw, instr(after_kw, '`') + 1), '`') - 1
        )
      ELSE '(unlabeled)'
    END AS workstream
  FROM dispatch_tail
)
SELECT workstream, COUNT(*) AS sessions
FROM labeled
GROUP BY workstream
ORDER BY (workstream = '(unlabeled)') ASC, sessions DESC, workstream ASC;

-- name: codex-activity
-- #323 — per-tool ACTIVITY tally over the codex-scoped set. Codex's tool surface is its own
-- shape: `exec_command` (shell commands — reads dominate; this is also the $.workdir signal),
-- `update_plan` (STRUCTURED plan-step progression — a real decision/workstream signal),
-- `spawn_agent` (sub-agent fan-out — Codex multi-agent). The report renders these as the
-- Codex activity summary (a coherent DB signal, NOT raw-rollout parsing). Same workdir-prefix
-- scope as codex-scoped, so the tally covers only this repo's attributed Codex.
SELECT t.tool_name AS tool, COUNT(*) AS calls
FROM tool_calls t
JOIN sessions s ON t.session_id = s.id
WHERE s.agent = 'codex'
  AND t.tool_name IN ('exec_command', 'update_plan', 'spawn_agent')
  AND EXISTS (
    SELECT 1 FROM tool_calls w
    WHERE w.session_id = s.id
      AND w.tool_name = 'exec_command'
      AND (json_extract(w.input_json, '$.workdir') = :repo_root
           OR json_extract(w.input_json, '$.workdir') LIKE :repo_root || '/%')
  )
GROUP BY t.tool_name
ORDER BY calls DESC;

-- name: scaffold-usage
-- #319 / #317.1 — behavioral scaffold tally. GROUP the tool_calls.skill_name signal by
-- FAMILY over the repo-scoped Claude session set: split on the first ':' (namespaced
-- names like `superpowers:brainstorming` → family `superpowers`); bare names that are
-- superpowers disciplines (`brainstorming`, `running-research-spikes`, …) also map to
-- `superpowers`; anything else keeps its own name as the family. The `spacedock` family
-- is EXCLUDED — survey/ensign self-invocation dominates the tally (1026/1045 on this
-- repo) and would make every surveyed repo read as "uses spacedock". This is the
-- behavioral half the file-only probe misses: a family that appears here with no files
-- on disk was invoked but not checked in (the report states the usage + presence fact).
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
-- #317.2 — WORK-BY-AREA: identity from edits, by LOGICAL area regardless of physical
-- location. Worktree-based projects (e.g. a `work-on-issue.sh` that drives the agent IN a
-- worktree per issue) land their PRODUCT code under `.worktrees/<wt>/…` — so a naive
-- first-segment bucket would mis-file a `.worktrees/<wt>/src/x` edit as `.worktrees` and
-- HIDE the real work. The fix: strip a leading `.worktrees/<wt>/` (and `.claude/worktrees/
-- <wt>/`) physical prefix FIRST, then bucket by the next segment — so a worktree `src/`
-- edit and a main-checkout `src/` edit BOTH count as `src`. The result reflects logical
-- areas (`src`/`internal`/`docs`/…) wherever the edit physically landed.
--
-- A `kind` partition then demotes genuine NON-code paths WITHOUT filtering them (still
-- counted, honest full accounting): `config` = `.claude` (memory/config/plans), `.beads`
-- (issue tracker), `.git`, and `<external>` (a sibling repo path outside the repo root —
-- a reference, not this project). `product` = every logical area under the repo root.
-- `ORDER BY (kind='config') ASC, edits DESC` puts product areas FIRST even when a config
-- bucket out-counts them; the report leads with product and footnotes the config sum.
SELECT
  area,
  CASE WHEN area IN ('.claude', '.beads', '.git', '<external>') THEN 'config' ELSE 'product' END AS kind,
  COUNT(*) AS edits
FROM (
  SELECT
    -- area: the logical first segment after stripping any worktree physical prefix.
    CASE
      WHEN logical IS NULL THEN '<external>'
      WHEN instr(logical, '/') > 0 THEN substr(logical, 1, instr(logical, '/') - 1)
      ELSE logical
    END AS area
  FROM (
    SELECT
      CASE
        WHEN fp NOT LIKE :repo_root || '/%' THEN NULL  -- external sibling: bucket <external>
        -- strip a leading `.worktrees/<wt>/` (drop the first two segments of the relative path)
        WHEN rel LIKE '.worktrees/%/%' THEN
          substr(rel, instr(substr(rel, length('.worktrees/') + 1), '/') + length('.worktrees/') + 1)
        -- strip a leading `.claude/worktrees/<wt>/` (drop the first three segments)
        WHEN rel LIKE '.claude/worktrees/%/%' THEN
          substr(rel, instr(substr(rel, length('.claude/worktrees/') + 1), '/') + length('.claude/worktrees/') + 1)
        ELSE rel
      END AS logical
    FROM (
      SELECT fp,
        CASE WHEN fp LIKE :repo_root || '/%'
             THEN substr(fp, length(:repo_root) + 2)  -- path relative to the repo root
             ELSE NULL END AS rel
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
    )
  )
)
GROUP BY area
ORDER BY
  (area IN ('.claude', '.beads', '.git', '<external>')) ASC,  -- config demotes, product leads
  edits DESC;

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

-- name: mode-classification
-- #324 (G) — classify each TRACK into a work MODE, so the report can make the RIGHT
-- commission offer per track (automation for mechanical, book-keeping for exploration)
-- instead of one undifferentiated pitch. The track key is the session's `git_branch`:
-- worktree-based projects (a `work-on-issue.sh` that branches per issue) carry one branch
-- per track, and content/design exploration runs on its own branch(es) too — so branch is
-- the per-track key the survey already has. Sessions with a blank branch are not a track
-- and drop out (they fold into the generic report, never a guessed mode).
--
-- Per track, tally the signatures the survey already reads (all repo-scoped, subagent-free):
--   veto     — `[Request interrupted` / `doesn't want to proceed` markers in the messages
--   loop     — `worktree` / `work-on-issue` markers (the mechanical issue→worktree→PR loop)
--   passed   — answered/approved AskUserQuestion/ExitPlanMode decisions (gate-pass)
--   rejected — the user-doesn't-want-to-proceed decisions (gate-fail / cancelled path)
--   code     — Edit/Write to a code file (`.go`/`.ts`/`.py`/`.rs`/`.js`/`.tsx`/`.go`…)
--   prose    — Edit/Write to a `.md` content/doc file
-- Score the two signatures and label by the DOMINANT one with a MARGIN guard:
--   mechanical signature = loop present + gate-pass-dominant + zero veto + code-heavy
--   exploration signature = veto present + a rejected/cancelled path + prose-heavy
-- A label is assigned ONLY when one score beats the other by >= 2 (a clear dominance); a
-- track with neither clearly dominant stays `unlabeled` and the report gives it the generic
-- book-keeping offer — NEVER a guessed automation pitch (the asymmetry favors not
-- mis-offering: a missed automation offer is a cheap omission; a wrong automation pitch at
-- creative work is the misread to avoid). The report reads `mode` per track to pick the offer.
WITH track_sessions AS (
  SELECT s.id, s.git_branch AS track
  FROM sessions s
  WHERE s.agent = 'claude'
    AND s.git_branch <> ''
    AND s.file_path NOT LIKE '%/subagents/%'
    AND (s.cwd = :repo_root OR s.cwd LIKE :repo_root || '/%')
),
vetoes AS (
  SELECT ts.track, COUNT(*) AS n
  FROM track_sessions ts JOIN messages m ON m.session_id = ts.id
  WHERE m.content LIKE '%[Request interrupted%'
     OR m.content LIKE '%doesn''t want to proceed%'
  GROUP BY ts.track
),
loops AS (
  SELECT ts.track, COUNT(*) AS n
  FROM track_sessions ts JOIN messages m ON m.session_id = ts.id
  WHERE m.content LIKE '%worktree%' OR m.content LIKE '%work-on-issue%'
  GROUP BY ts.track
),
passed AS (
  SELECT ts.track, COUNT(*) AS n
  FROM track_sessions ts JOIN tool_calls t ON t.session_id = ts.id
  WHERE t.tool_name IN ('AskUserQuestion', 'ExitPlanMode')
    AND (t.result_content LIKE 'Your questions have been answered%'
         OR t.result_content LIKE 'Your question has been answered%'
         OR t.result_content LIKE 'User has answered%'
         OR t.result_content LIKE 'User has approved your plan%'
         OR t.result_content LIKE 'User approved%')
  GROUP BY ts.track
),
rejected AS (
  SELECT ts.track, COUNT(*) AS n
  FROM track_sessions ts JOIN tool_calls t ON t.session_id = ts.id
  WHERE t.tool_name IN ('AskUserQuestion', 'ExitPlanMode')
    AND t.result_content LIKE '%doesn''t want to proceed%'
  GROUP BY ts.track
),
edits AS (
  SELECT ts.track,
    SUM(CASE WHEN fp LIKE '%.md' THEN 1 ELSE 0 END) AS prose,
    SUM(CASE WHEN fp LIKE '%.go' OR fp LIKE '%.ts' OR fp LIKE '%.tsx'
              OR fp LIKE '%.py' OR fp LIKE '%.rs' OR fp LIKE '%.js' THEN 1 ELSE 0 END) AS code
  FROM track_sessions ts
  JOIN (
    SELECT session_id, json_extract(input_json, '$.file_path') AS fp
    FROM tool_calls WHERE tool_name IN ('Edit', 'Write')
  ) e ON e.session_id = ts.id
  GROUP BY ts.track
),
sig AS (
  SELECT t.track,
    COALESCE(v.n, 0) AS veto, COALESCE(l.n, 0) AS loop,
    COALESCE(p.n, 0) AS passed, COALESCE(r.n, 0) AS rejected,
    COALESCE(e.code, 0) AS code, COALESCE(e.prose, 0) AS prose
  FROM (SELECT DISTINCT track FROM track_sessions) t
  LEFT JOIN vetoes v ON v.track = t.track
  LEFT JOIN loops l ON l.track = t.track
  LEFT JOIN passed p ON p.track = t.track
  LEFT JOIN rejected r ON r.track = t.track
  LEFT JOIN edits e ON e.track = t.track
),
scored AS (
  SELECT *,
    ((loop > 0) + (passed > rejected) + (veto = 0) + (code > prose)) AS mech,
    ((veto > 0) + (rejected > 0) + (prose > code)) AS expl
  FROM sig
)
SELECT track,
  CASE
    WHEN mech - expl >= 2 THEN 'mechanical'
    WHEN expl - mech >= 2 THEN 'exploration'
    ELSE 'unlabeled'
  END AS mode
FROM scored
ORDER BY track;
