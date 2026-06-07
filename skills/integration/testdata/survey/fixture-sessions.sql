-- ABOUTME: Committed agentsview-shaped sessions.db fixture for the survey skill's
-- ABOUTME: query-smoke — the corrected queries (references/queries.sql) run against it.
--
-- These rows mirror what `agentsview sync` produces for a multi-checkout repo, with the
-- columns the corrected queries touch. The smoke test (skills/integration) runs each
-- labeled query from skills/survey/references/queries.sql against this fixture, binding
-- a FIXED repo-root prefix string (no live git repo needed), and asserts the corrected
-- SHAPE — so a broken query, a dropped self-exclusion, a wrong prefix bound, or schema
-- drift reds the smoke.
--
-- Anonymized-real-shape: counts and result/skill_name shapes match production
-- (agentsview v0.32.1), identifiers stripped. The shapes the smoke pins:
--   - sessions.cwd spans the repo root, a subdir, and a worktree-style path under ONE
--     repo-root prefix, plus a session OUTSIDE the prefix — so the #318 scoping query
--     coalesces the in-repo cwds across divergent basename-derived `project` keys and
--     excludes the out-of-repo session.
--   - tool_calls.skill_name carries namespaced (`superpowers:*`), bare
--     (`running-research-spikes`), and `spacedock:*` self rows — so the #319 family
--     tally reports a `superpowers` family and EXCLUDES the dominant `spacedock` self
--     rows. A `superpowers` family with no files on disk is the `recovered` case.
--   - Edit/Write input_json.$.file_path rows under the repo root (an `internal` and a
--     `skills` bucket) and one OUTSIDE it — so the #317.2 WORK-BY-AREA query buckets by
--     package and flags the external-sibling path as `<external>`.
--   - decision rows: an answered AskUserQuestion (done), a rejected one (OPEN), and an
--     ExitPlanMode "User has approved your plan" approval — so the #320 query marks the
--     approved plan `done` (the cheap done-prefix fix) and the rejection OPEN.
--
-- CRITICAL — decision results match the PRODUCTION shape. agentsview v0.32.1 emits a
-- NON-EMPTY result_content for EVERY decision: answered → "Your questions have been
-- answered…"; rejected → "The user doesn't want to proceed…" (questions left "(No
-- answer provided)"); ExitPlanMode approval → "User has approved your plan". A
-- NULL-result OPEN row would be a shape production never emits — a tautology against a
-- fiction. OPEN-detection keys on the ABSENCE of an answered/approved confirmation.
--
-- A sibling NON-Claude (codex) session inside the repo prefix proves the Claude-only
-- scope: it must NOT leak into the Claude counts. (Surfacing non-Claude agents is a
-- deferred follow-up.)
--
-- Defaults match production: cwd TEXT NOT NULL DEFAULT '' (blank stored as '', not
-- NULL); skill_name TEXT nullable, populated only on tool_name='Skill' rows; git_branch
-- and category are seeded production-shaped but NO query reads them (YAGNI).

CREATE TABLE sessions (
  id TEXT PRIMARY KEY,
  project TEXT,
  agent TEXT,
  cwd TEXT NOT NULL DEFAULT '',
  git_branch TEXT NOT NULL DEFAULT '',
  file_path TEXT,
  started_at TEXT,
  ended_at TEXT,
  first_message TEXT,
  message_count INTEGER,
  user_message_count INTEGER
);

CREATE TABLE tool_calls (
  id INTEGER PRIMARY KEY,
  session_id TEXT,
  tool_name TEXT,
  category TEXT NOT NULL DEFAULT '',
  skill_name TEXT,
  input_json TEXT,
  result_content TEXT
);

CREATE TABLE messages (
  id INTEGER PRIMARY KEY,
  session_id TEXT,
  role TEXT,
  content TEXT
);

-- ============================================================================
-- The repo-root prefix the smoke binds is /repo/proj/ . Three Claude sessions sit
-- under it with DIVERGENT basename-derived `project` keys (repo root, a subdir, a
-- worktree-style path); one Claude session sits OUTSIDE it; one Claude session has a
-- blank cwd; a codex session sits inside but is out of Claude scope.
-- ============================================================================

-- Claude session A — repo root. project basename matches the repo basename (`proj`).
-- Carries the answered decision and an Edit under internal/.
INSERT INTO sessions VALUES
  ('claude:aaaaaaaa-1111-2222-3333-444444444444', 'proj', 'claude',
   '/repo/proj', 'main',
   '/u/.claude/projects/-repo-proj/aaaaaaaa.jsonl',
   '2026-06-05', '2026-06-05', 'Pick up the parser refactor and ship it.', 8, 3);

-- Claude session B — a SUBDIR checkout. Its cwd basename (`.spacedock-state`) keys as
-- `_spacedock_state`, diverging from the repo basename — invisible to a basename-only
-- scope, recovered by the prefix union. Carries the rejected (OPEN) decision + a veto.
INSERT INTO sessions VALUES
  ('claude:bbbbbbbb-5555-6666-7777-888888888888', '_spacedock_state', 'claude',
   '/repo/proj/docs/dev/.spacedock-state', 'main',
   '/u/.claude/projects/-repo-proj-docs-dev-_spacedock_state/bbbbbbbb.jsonl',
   '2026-06-06', '2026-06-06', 'Now wire up the regression suite.', 6, 2);

-- Claude session C — a WORKTREE-style checkout. cwd basename (`feature-x`) keys as
-- `feature_x`, again divergent. Carries the ExitPlanMode approval + a skills/ Write +
-- the superpowers / bare / spacedock Skill rows for the scaffold tally.
INSERT INTO sessions VALUES
  ('claude:cccccccc-9999-aaaa-bbbb-cccccccccccc', 'feature_x', 'claude',
   '/repo/proj/.worktrees/feature-x', 'feature-x',
   '/u/.claude/projects/-repo-proj-.worktrees-feature-x/cccccccc.jsonl',
   '2026-06-06', '2026-06-06', 'Build the feature behind a worktree.', 5, 2);

-- Claude session D — a BLANK cwd (production stores blank as ''). Under no prefix; it
-- must NOT count toward the repo scope, and the blank-cwd tally surfaces it.
INSERT INTO sessions VALUES
  ('claude:dddddddd-0000-1111-2222-333333333333', '', 'claude',
   '', '',
   '/u/.claude/projects/unknown/dddddddd.jsonl',
   '2026-06-03', '2026-06-03', 'A session whose cwd agentsview never captured.', 3, 1);

-- Claude session E — OUTSIDE the repo prefix entirely (a different project on the same
-- machine). The scoping query MUST exclude it; its Skill row MUST NOT inflate the tally.
INSERT INTO sessions VALUES
  ('claude:eeeeeeee-4444-5555-6666-777777777777', 'otherproj', 'claude',
   '/elsewhere/otherproj', 'main',
   '/u/.claude/projects/-elsewhere-otherproj/eeeeeeee.jsonl',
   '2026-06-02', '2026-06-02', 'Unrelated project that shares the machine.', 4, 1);

-- Out-of-scope sibling: a Codex session INSIDE the repo prefix. Claude scope excludes
-- it from every count.
INSERT INTO sessions VALUES
  ('codex:ffffffff-8888-9999-aaaa-bbbbbbbbbbbb', 'proj', 'codex',
   '/repo/proj', 'main',
   '/u/.codex/sessions/rollout-ffffffff.jsonl',
   '2026-06-04', '2026-06-04', 'Out-of-scope codex session under the repo root.', 4, 1);

-- ----------------------------------------------------------------------------
-- DECISIONS (#320): answered (done), rejected (OPEN), ExitPlanMode approval (done).
-- ----------------------------------------------------------------------------

-- Answered AskUserQuestion → done. (session A)
INSERT INTO tool_calls (id, session_id, tool_name, skill_name, input_json, result_content) VALUES
  (1, 'claude:aaaaaaaa-1111-2222-3333-444444444444', 'AskUserQuestion', NULL,
   '{"questions":[{"header":"Refactor scope","question":"Which scope should the parser refactor cover?"}]}',
   'Your questions have been answered: "Which scope should the parser refactor cover?"="tokenizer + entrypoint"');

-- Rejected AskUserQuestion → OPEN, production shape (non-empty rejection result). (session B)
INSERT INTO tool_calls (id, session_id, tool_name, skill_name, input_json, result_content) VALUES
  (2, 'claude:bbbbbbbb-5555-6666-7777-888888888888', 'AskUserQuestion', NULL,
   '{"questions":[{"header":"Test framework","question":"Which test framework should the regression suite use?"}]}',
   'The user doesn''t want to proceed with this tool use. The tool use was rejected. Questions asked:
- "Which test framework should the regression suite use?"
  (No answer provided)');

-- ExitPlanMode approval → done via the "User has approved your plan" prefix (the #320
-- cheap fix). Before the fix this matched NO done-prefix and fell to OPEN. (session C)
INSERT INTO tool_calls (id, session_id, tool_name, skill_name, input_json, result_content) VALUES
  (3, 'claude:cccccccc-9999-aaaa-bbbb-cccccccccccc', 'ExitPlanMode', NULL,
   '{"plan":"Implement the feature behind the worktree, then add tests."}',
   'User has approved your plan. You can now start coding.');

-- ----------------------------------------------------------------------------
-- SCAFFOLD-USAGE (#319): Skill rows across families. spacedock:* dominates (self) and
-- MUST be excluded; superpowers (namespaced + bare) survives as `recovered`. (session C)
-- ----------------------------------------------------------------------------
INSERT INTO tool_calls (id, session_id, tool_name, skill_name, input_json, result_content) VALUES
  (10, 'claude:cccccccc-9999-aaaa-bbbb-cccccccccccc', 'Skill', 'spacedock:ensign', NULL, NULL),
  (11, 'claude:cccccccc-9999-aaaa-bbbb-cccccccccccc', 'Skill', 'spacedock:ensign', NULL, NULL),
  (12, 'claude:cccccccc-9999-aaaa-bbbb-cccccccccccc', 'Skill', 'spacedock:first-officer', NULL, NULL),
  (13, 'claude:cccccccc-9999-aaaa-bbbb-cccccccccccc', 'Skill', 'spacedock:debrief', NULL, NULL),
  (14, 'claude:cccccccc-9999-aaaa-bbbb-cccccccccccc', 'Skill', 'superpowers:brainstorming', NULL, NULL),
  (15, 'claude:cccccccc-9999-aaaa-bbbb-cccccccccccc', 'Skill', 'superpowers:systematic-debugging', NULL, NULL),
  (16, 'claude:cccccccc-9999-aaaa-bbbb-cccccccccccc', 'Skill', 'running-research-spikes', NULL, NULL);

-- A Skill row in the OUT-OF-REPO session E — it must NOT inflate the in-repo tally.
INSERT INTO tool_calls (id, session_id, tool_name, skill_name, input_json, result_content) VALUES
  (20, 'claude:eeeeeeee-4444-5555-6666-777777777777', 'Skill', 'superpowers:writing-plans', NULL, NULL);

-- ----------------------------------------------------------------------------
-- WORK-BY-AREA (#317.2): Edit/Write file_path rows — two under the repo root (an
-- `internal` and a `skills` bucket) and one OUTSIDE it (an external sibling reference).
-- ----------------------------------------------------------------------------
INSERT INTO tool_calls (id, session_id, tool_name, input_json, result_content) VALUES
  (30, 'claude:aaaaaaaa-1111-2222-3333-444444444444', 'Edit',
   '{"file_path":"/repo/proj/internal/dispatch/build.go"}', NULL),
  (31, 'claude:aaaaaaaa-1111-2222-3333-444444444444', 'Edit',
   '{"file_path":"/repo/proj/internal/dispatch/parse.go"}', NULL),
  (32, 'claude:cccccccc-9999-aaaa-bbbb-cccccccccccc', 'Write',
   '{"file_path":"/repo/proj/skills/survey/references/queries.sql"}', NULL),
  (33, 'claude:cccccccc-9999-aaaa-bbbb-cccccccccccc', 'Edit',
   '{"file_path":"/sibling/otherlib/src/util.go"}', NULL);

-- ----------------------------------------------------------------------------
-- Veto marker in session B's message stream (interruption signal, prose-read).
-- ----------------------------------------------------------------------------
INSERT INTO messages VALUES
  (1, 'claude:aaaaaaaa-1111-2222-3333-444444444444', 'user', 'Pick up the parser refactor and ship it.'),
  (2, 'claude:bbbbbbbb-5555-6666-7777-888888888888', 'user', 'Now wire up the regression suite.'),
  (3, 'claude:bbbbbbbb-5555-6666-7777-888888888888', 'user', '[Request interrupted by user]'),
  (4, 'claude:cccccccc-9999-aaaa-bbbb-cccccccccccc', 'user', 'Build the feature behind a worktree.'),
  (5, 'codex:ffffffff-8888-9999-aaaa-bbbbbbbbbbbb', 'user', 'Out-of-scope codex session under the repo root.');
