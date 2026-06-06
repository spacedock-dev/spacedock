-- ABOUTME: Committed agentsview-shaped sessions.db fixture for the survey skill's
-- ABOUTME: Claude extraction test — two Claude sessions under one project key, with decisions + a veto.
--
-- These rows mirror what `agentsview sync` produces for this repo's Claude history:
-- the project key is the cwd basename with non-alphanumerics replaced by '_'
-- (survey_fixture_proj here). Survey filters by that `project` column AND agent='claude'.
-- The fixture seeds a known set of decision/interruption signals so the skill's step-2
-- block, run verbatim, must surface them; a regression in the queries (wrong project
-- filter, dropped agent scope, dropped OPEN-first ordering, broken json_extract) drops
-- a known row and reds the test.
--
-- CRITICAL — the decision results match the PRODUCTION shape. agentsview v0.32.1 emits a
-- NON-EMPTY result_content for EVERY decision: an answered one begins "Your questions have
-- been answered…"/"User has answered…", a rejected/abandoned one begins "The user doesn't
-- want to proceed… was rejected" (questions left "(No answer provided)"). OPEN-detection
-- must key on the ABSENCE of an answered-confirmation — a NULL-result OPEN row would be a
-- shape production never emits, making the test a tautology against a fiction.
--
-- A sibling NON-Claude (codex) session under the SAME project key is included to prove
-- the Claude-only scope: it must NOT leak into the Claude counts. (Surfacing non-Claude
-- agents' decision signals is a deferred follow-up — out of scope for now.)
--
-- The 20 answered Claude decisions with ids newer than the OPEN row are intentional:
-- the shipped query's recency tail is `ORDER BY status ASC, t.id DESC LIMIT 20`.
-- If a regression removes `status ASC`, those 20 answered rows fill the LIMIT window
-- and truncate the older OPEN frontier row out of DECISIONS.
--
-- The columns are the subset of agentsview's v1 schema the skill's queries touch.

CREATE TABLE sessions (
  id TEXT PRIMARY KEY,
  project TEXT,
  agent TEXT,
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
  input_json TEXT,
  result_content TEXT
);

CREATE TABLE messages (
  id INTEGER PRIMARY KEY,
  session_id TEXT,
  role TEXT,
  content TEXT
);

-- Claude session 1: one ANSWERED decision + a normal user prompt. The answered result
-- mirrors the PRODUCTION shape: real agentsview v0.32.1 records an answered decision with
-- result_content beginning "Your questions have been answered: …" (or "User has answered …").
INSERT INTO sessions VALUES
  ('claude:11111111-2222-3333-4444-555555555555', 'survey_fixture_proj', 'claude',
   '/u/.claude/projects/-tmp-survey_fixture_proj/11111111.jsonl',
   '2026-06-05', '2026-06-05', 'Pick up the parser refactor and ship it.', 8, 3);

INSERT INTO tool_calls VALUES
  (1, 'claude:11111111-2222-3333-4444-555555555555', 'AskUserQuestion',
   '{"questions":[{"header":"Refactor scope","question":"Which scope should the parser refactor cover?"}]}',
   'Your questions have been answered: "Which scope should the parser refactor cover?"="tokenizer + entrypoint"');

INSERT INTO messages VALUES
  (1, 'claude:11111111-2222-3333-4444-555555555555', 'user', 'Pick up the parser refactor and ship it.');

-- Claude session 2: one OPEN decision matching the PRODUCTION shape — a REJECTED decision.
-- Real agentsview NEVER emits a NULL result for a decision; a rejected/abandoned fork still
-- carries a non-empty result_content: "The user doesn't want to proceed… The tool use was
-- rejected", with the questions left "(No answer provided)". OPEN-detection MUST key on the
-- ABSENCE of an answered-confirmation, not on a NULL that production never produces. (Plus a
-- user-veto interrupt marker in the message stream.)
INSERT INTO sessions VALUES
  ('claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'survey_fixture_proj', 'claude',
   '/u/.claude/projects/-tmp-survey_fixture_proj/66666666.jsonl',
   '2026-06-06', '2026-06-06', 'Now wire up the regression suite.', 6, 2);

INSERT INTO tool_calls VALUES
  (2, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Test framework","question":"Which test framework should the regression suite use?"}]}',
   'The user doesn''t want to proceed with this tool use. The tool use was rejected. Questions asked:
- "Which test framework should the regression suite use?"
  (No answer provided)');

INSERT INTO tool_calls VALUES
  (4, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Recent answered 01","question":"Which recent answered decision 01 should be kept?"}]}',
   'Your questions have been answered: "Which recent answered decision 01 should be kept?"="answer 01"'),
  (5, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Recent answered 02","question":"Which recent answered decision 02 should be kept?"}]}',
   'Your questions have been answered: "Which recent answered decision 02 should be kept?"="answer 02"'),
  (6, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Recent answered 03","question":"Which recent answered decision 03 should be kept?"}]}',
   'Your questions have been answered: "Which recent answered decision 03 should be kept?"="answer 03"'),
  (7, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Recent answered 04","question":"Which recent answered decision 04 should be kept?"}]}',
   'Your questions have been answered: "Which recent answered decision 04 should be kept?"="answer 04"'),
  (8, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Recent answered 05","question":"Which recent answered decision 05 should be kept?"}]}',
   'Your questions have been answered: "Which recent answered decision 05 should be kept?"="answer 05"'),
  (9, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Recent answered 06","question":"Which recent answered decision 06 should be kept?"}]}',
   'Your questions have been answered: "Which recent answered decision 06 should be kept?"="answer 06"'),
  (10, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Recent answered 07","question":"Which recent answered decision 07 should be kept?"}]}',
   'Your questions have been answered: "Which recent answered decision 07 should be kept?"="answer 07"'),
  (11, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Recent answered 08","question":"Which recent answered decision 08 should be kept?"}]}',
   'Your questions have been answered: "Which recent answered decision 08 should be kept?"="answer 08"'),
  (12, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Recent answered 09","question":"Which recent answered decision 09 should be kept?"}]}',
   'Your questions have been answered: "Which recent answered decision 09 should be kept?"="answer 09"'),
  (13, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Recent answered 10","question":"Which recent answered decision 10 should be kept?"}]}',
   'Your questions have been answered: "Which recent answered decision 10 should be kept?"="answer 10"'),
  (14, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Recent answered 11","question":"Which recent answered decision 11 should be kept?"}]}',
   'Your questions have been answered: "Which recent answered decision 11 should be kept?"="answer 11"'),
  (15, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Recent answered 12","question":"Which recent answered decision 12 should be kept?"}]}',
   'Your questions have been answered: "Which recent answered decision 12 should be kept?"="answer 12"'),
  (16, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Recent answered 13","question":"Which recent answered decision 13 should be kept?"}]}',
   'Your questions have been answered: "Which recent answered decision 13 should be kept?"="answer 13"'),
  (17, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Recent answered 14","question":"Which recent answered decision 14 should be kept?"}]}',
   'Your questions have been answered: "Which recent answered decision 14 should be kept?"="answer 14"'),
  (18, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Recent answered 15","question":"Which recent answered decision 15 should be kept?"}]}',
   'Your questions have been answered: "Which recent answered decision 15 should be kept?"="answer 15"'),
  (19, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Recent answered 16","question":"Which recent answered decision 16 should be kept?"}]}',
   'Your questions have been answered: "Which recent answered decision 16 should be kept?"="answer 16"'),
  (20, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Recent answered 17","question":"Which recent answered decision 17 should be kept?"}]}',
   'Your questions have been answered: "Which recent answered decision 17 should be kept?"="answer 17"'),
  (21, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Recent answered 18","question":"Which recent answered decision 18 should be kept?"}]}',
   'Your questions have been answered: "Which recent answered decision 18 should be kept?"="answer 18"'),
  (22, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Recent answered 19","question":"Which recent answered decision 19 should be kept?"}]}',
   'Your questions have been answered: "Which recent answered decision 19 should be kept?"="answer 19"'),
  (23, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'AskUserQuestion',
   '{"questions":[{"header":"Recent answered 20","question":"Which recent answered decision 20 should be kept?"}]}',
   'Your questions have been answered: "Which recent answered decision 20 should be kept?"="answer 20"');

INSERT INTO messages VALUES
  (2, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'user', 'Now wire up the regression suite.'),
  (3, 'claude:66666666-7777-8888-9999-aaaaaaaaaaaa', 'user', '[Request interrupted by user]');

-- Out-of-scope sibling: a Codex session under the SAME project key. Survey's Claude
-- scope must EXCLUDE it — it must not appear in the session count, the decisions, or the
-- interruption totals.
INSERT INTO sessions VALUES
  ('codex:cccccccc-dddd-eeee-ffff-000000000000', 'survey_fixture_proj', 'codex',
   '/u/.codex/sessions/rollout-cccccccc.jsonl',
   '2026-06-04', '2026-06-04', 'Out-of-scope codex session under the same project key.', 4, 1);

INSERT INTO tool_calls VALUES
  (3, 'codex:cccccccc-dddd-eeee-ffff-000000000000', 'update_plan',
   '{"explanation":"x","plan":[{"step":"A codex-only step that must not surface","status":"in_progress"}]}',
   NULL);

INSERT INTO messages VALUES
  (4, 'codex:cccccccc-dddd-eeee-ffff-000000000000', 'user', 'Out-of-scope codex session under the same project key.');
