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
--     counts the in-repo cwds and excludes the out-of-repo session. agentsview keys
--     `project` by the git-root basename, so all three in-repo Claude checkouts share ONE
--     `project` key (`proj`); the cwd-prefix's job is to EXCLUDE a same-basename sibling.
--   - tool_calls.skill_name carries namespaced (`superpowers:*`), bare
--     (`running-research-spikes`), and `spacedock:*` self rows — so the #319 family
--     tally reports a `superpowers` family and EXCLUDES the dominant `spacedock` self
--     rows. A `superpowers` family with no files on disk was invoked but not checked in.
--   - Edit/Write input_json.$.file_path rows: worktree `src/` edits + a main-checkout
--     `src/` edit (the #317.2 WORK-BY-AREA query strips `.worktrees/<wt>/` and buckets
--     them all as `src`), `internal`/`docs` product areas, `.claude`/`.beads` config
--     (demoted via the `kind` partition, still counted), and an external-sibling path
--     OUTSIDE the prefix (flagged `<external>`, also config-demoted).
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
-- Codex sessions land cwd='' in production (agentsview does not persist Codex cwd), so
-- they cannot be cwd-prefix-scoped and the #69 codex-presence query matches them by
-- `project` (the git-root basename) ALONE. Two blank-cwd Codex rows carry `project='proj'`:
-- one is THIS repo's Codex history, one is a SAME-BASENAME sibling repo whose Codex
-- sessions key identically (the documented collision) — both are counted by codex-presence
-- (which is why the report states the match is by project NAME only), and NEITHER leaks
-- into the cwd-scoped Claude counts (Claude scope filters agent='claude').
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
-- The repo-root prefix the smoke binds is /repo/proj/ . Three Claude sessions sit under
-- it (repo root, a subdir, a worktree-style path), all sharing the ONE git-root-basename
-- `project` key (`proj`); one Claude session sits OUTSIDE it; one Claude session has a
-- blank cwd; two blank-cwd Codex sessions key to `proj` (one this repo, one a same-basename
-- sibling) and are out of Claude scope.
-- ============================================================================

-- Claude session A — repo root. project is the git-root basename (`proj`).
-- Carries the answered decision and an Edit under internal/.
INSERT INTO sessions VALUES
  ('claude:aaaaaaaa-1111-2222-3333-444444444444', 'proj', 'claude',
   '/repo/proj', 'main',
   '/u/.claude/projects/-repo-proj/aaaaaaaa.jsonl',
   '2026-06-05', '2026-06-05', 'Pick up the parser refactor and ship it.', 8, 3);

-- Claude session B — a SUBDIR checkout (the split-root state dir). agentsview keys it by
-- the git-root basename, so its `project` is `proj`, the SAME key as the root — the
-- cwd-prefix is what places it in scope. Carries the rejected (OPEN) decision + a veto.
INSERT INTO sessions VALUES
  ('claude:bbbbbbbb-5555-6666-7777-888888888888', 'proj', 'claude',
   '/repo/proj/docs/dev/.spacedock-state', 'main',
   '/u/.claude/projects/-repo-proj-docs-dev-_spacedock_state/bbbbbbbb.jsonl',
   '2026-06-06', '2026-06-06', 'Now wire up the regression suite.', 6, 2);

-- Claude session C — a WORKTREE-style checkout. Same git-root basename, so `project` is
-- again `proj` — placed in scope by the cwd-prefix, not a distinct key. Carries the
-- ExitPlanMode approval + a skills/ Write + the superpowers / bare / spacedock Skill rows
-- for the scaffold tally.
INSERT INTO sessions VALUES
  ('claude:cccccccc-9999-aaaa-bbbb-cccccccccccc', 'proj', 'claude',
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

-- ============================================================================
-- CODEX SESSIONS — four attributed to THIS repo (F, F2, F3, F4) + one same-basename
-- SIBLING (G). All carry cwd='' (agentsview persists no Codex session cwd). Attribution
-- to this repo is by exec_command.$.workdir prefix (the codex-scoped query, #321),
-- NOT by project name. The four F* sessions get an exec_command row whose $.workdir is
-- under /repo/proj (codex-scoped counts them = 4); G's is under /sibling/proj (excluded).
-- All five carry project='proj' so codex-presence (name-only, #69) counts 5 — proving the
-- two signals differ (presence 5 ⊃ scoped 4). The four F* first_messages are real-shape so
-- the codex-workstreams clustering rule (#322) has a dispatch-pattern, a task/entity-pattern,
-- an unlabeled, and a SECOND distinct dispatch task to cluster. (AC-1's "1 vs 2" describes
-- the MECHANISM minimally; the clustering AC needs 4 attributed sessions, so the fixture
-- scales to scoped=4 / presence=5 — the binding asserts are sibling-excluded + the
-- prefix-load-bearing flip, which hold at 4/5 exactly as at 1/2.)
-- ============================================================================

-- Codex F — DISPATCH pattern → workstream `journey-cost-ledger` (stage suffix stripped).
INSERT INTO sessions VALUES
  ('codex:ffffffff-8888-9999-aaaa-bbbbbbbbbbbb', 'proj', 'codex',
   '', '',
   '/u/.codex/sessions/rollout-ffffffff.jsonl',
   '2026-06-04', '2026-06-04',
   'Read /tmp/spacedock-dispatch/spacedock-ensign-journey-cost-ledger-implementation.md and treat its content as your assignment.',
   4, 1);

-- Codex F2 — TASK/ENTITY backtick pattern → workstream `orient-workflow-discovery`. A
-- leading reviewer-label backtick precedes the keyword so the rule must anchor on the
-- `Spacedock entity` token, not the first backtick globally.
INSERT INTO sessions VALUES
  ('codex:f2f2f2f2-0000-1111-2222-333333333333', 'proj', 'codex',
   '', '',
   '/u/.codex/sessions/rollout-f2f2f2f2.jsonl',
   '2026-06-04', '2026-06-04',
   'You are `142-validation/Ensign`, a fresh validation worker for Spacedock entity 142 `orient-workflow-discovery`. Working directory: /repo/proj.',
   4, 1);

-- Codex F3 — UNLABELED. An encouragement/meta first_message carries no task → (unlabeled).
INSERT INTO sessions VALUES
  ('codex:f3f3f3f3-4444-5555-6666-777777777777', 'proj', 'codex',
   '', '',
   '/u/.codex/sessions/rollout-f3f3f3f3.jsonl',
   '2026-06-04', '2026-06-04',
   'You totally got this. Take your time. Captain asked me to tell subagents they are appreciated.',
   3, 1);

-- Codex F4 — a SECOND distinct DISPATCH task → workstream `codex-live-ci`. Proves the
-- cluster key is the extracted {TASK}, not a constant: F4 must NOT merge with F.
INSERT INTO sessions VALUES
  ('codex:f4f4f4f4-8888-9999-aaaa-bbbbbbbbbbbb', 'proj', 'codex',
   '', '',
   '/u/.codex/sessions/rollout-f4f4f4f4.jsonl',
   '2026-06-04', '2026-06-04',
   'Read /tmp/spacedock-dispatch/spacedock-ensign-codex-live-ci-validation.md and treat its content as your assignment.',
   4, 1);

-- Codex session G — a SAME-BASENAME SIBLING repo's Codex history. Its git-root basename is
-- also `proj`, so it keys to the identical `project` and codex-presence CANNOT distinguish
-- it from the F* sessions (the documented collision — the report states "match by project
-- NAME only"). Blank cwd, like all Codex sessions. Its exec_command $.workdir is under
-- /sibling/proj (OUTSIDE the repo prefix) so codex-scoped EXCLUDES it. Out of Claude scope.
INSERT INTO sessions VALUES
  ('codex:11111111-2222-3333-4444-555555555555', 'proj', 'codex',
   '', '',
   '/u/.codex/sessions/rollout-11111111.jsonl',
   '2026-06-04', '2026-06-04',
   'Read /tmp/spacedock-dispatch/spacedock-ensign-sibling-task-implementation.md and treat its content as your assignment.',
   3, 1);

-- ----------------------------------------------------------------------------
-- CODEX exec_command rows carry $.workdir — the attribution signal (#321 codex-scoped) and
-- the per-session activity signal (#323 codex-activity). The four F* sessions' workdirs are
-- under /repo/proj (one is a worktree path, proving the prefix admits worktrees); G's is
-- under /sibling/proj. update_plan + spawn_agent rows exercise the activity tally.
-- ----------------------------------------------------------------------------
INSERT INTO tool_calls (id, session_id, tool_name, input_json, result_content) VALUES
  (40, 'codex:ffffffff-8888-9999-aaaa-bbbbbbbbbbbb', 'exec_command',
   '{"command":"go test ./...","workdir":"/repo/proj/.worktrees/journey-cost-ledger"}', NULL),
  (41, 'codex:ffffffff-8888-9999-aaaa-bbbbbbbbbbbb', 'update_plan',
   '{"plan":[{"step":"explore","status":"completed"}]}', NULL),
  (42, 'codex:f2f2f2f2-0000-1111-2222-333333333333', 'exec_command',
   '{"command":"rg orient","workdir":"/repo/proj"}', NULL),
  (43, 'codex:f3f3f3f3-4444-5555-6666-777777777777', 'exec_command',
   '{"command":"ls","workdir":"/repo/proj/internal"}', NULL),
  (44, 'codex:f4f4f4f4-8888-9999-aaaa-bbbbbbbbbbbb', 'exec_command',
   '{"command":"git status","workdir":"/repo/proj"}', NULL),
  (45, 'codex:f4f4f4f4-8888-9999-aaaa-bbbbbbbbbbbb', 'spawn_agent',
   '{"task":"sub"}', NULL),
  -- G (sibling): exec_command workdir is OUTSIDE the repo prefix → codex-scoped excludes it.
  (46, 'codex:11111111-2222-3333-4444-555555555555', 'exec_command',
   '{"command":"go build","workdir":"/sibling/proj"}', NULL);

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
-- MUST be excluded; superpowers (namespaced + bare) survives, invoked-but-not-on-disk. (session C)
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

-- ============================================================================
-- WORK-BY-AREA WORKTREE-ATTRIBUTION (#317.2, F-corrected / AC-7a). A worktree-based
-- project (torahmap's `work-on-issue.sh` shape) drives the agent IN a worktree per issue,
-- so PRODUCT code lands under `.worktrees/<wt>/…`. The corrected query strips the physical
-- `.worktrees/<wt>/` prefix and buckets by the LOGICAL area, so a worktree `src/` edit
-- counts as `src` ALONGSIDE a main-checkout `src/` edit — NOT as `.worktrees`/`<external>`.
-- A `kind` partition demotes genuine config (`.claude`/`.beads`/`.git`/`<external>`) to a
-- footnote (still counted). Session WT's cwd IS a worktree (the torahmap shape).
-- ============================================================================

-- Claude session WT — a WORKTREE-cwd session (the torahmap work-on-issue shape). git_branch
-- `issue-42` is its track key. Edits two worktree `src/` files; carries the MECHANICAL
-- signature for mode-classification (gate-pass decision + worktree loop markers + code, no veto).
INSERT INTO sessions VALUES
  ('claude:77777777-aaaa-bbbb-cccc-dddddddddddd', 'proj', 'claude',
   '/repo/proj/.worktrees/issue-42', 'issue-42',
   '/u/.claude/projects/-repo-proj-.worktrees-issue-42/77777777.jsonl',
   '2026-06-07', '2026-06-07', 'Run the work-on-issue loop for issue 42 in its worktree.', 6, 2);

-- Worktree-attribution Edit/Write rows: two worktree `src/` edits (strip to `src`), a
-- main-checkout `src/` edit (also `src` — all three bucket together), a `docs/` product
-- edit, a `.claude` (config-demote), a `.beads` (config-demote). The existing #266 rows add
-- an `internal` product bucket + an `<external>` sibling (config-demote). A `.claude/worktrees/`
-- edit proves THAT prefix strips too (→ `internal`).
INSERT INTO tool_calls (id, session_id, tool_name, input_json, result_content) VALUES
  (50, 'claude:77777777-aaaa-bbbb-cccc-dddddddddddd', 'Edit',
   '{"file_path":"/repo/proj/.worktrees/issue-42/src/render.ts"}', NULL),
  (51, 'claude:77777777-aaaa-bbbb-cccc-dddddddddddd', 'Edit',
   '{"file_path":"/repo/proj/.worktrees/issue-42/src/palette.ts"}', NULL),
  (52, 'claude:77777777-aaaa-bbbb-cccc-dddddddddddd', 'Write',
   '{"file_path":"/repo/proj/src/main.ts"}', NULL),
  (53, 'claude:77777777-aaaa-bbbb-cccc-dddddddddddd', 'Write',
   '{"file_path":"/repo/proj/docs/spec.md"}', NULL),
  (54, 'claude:77777777-aaaa-bbbb-cccc-dddddddddddd', 'Edit',
   '{"file_path":"/repo/proj/.claude/memory.md"}', NULL),
  (55, 'claude:77777777-aaaa-bbbb-cccc-dddddddddddd', 'Edit',
   '{"file_path":"/repo/proj/.beads/tracker.db"}', NULL),
  (56, 'claude:77777777-aaaa-bbbb-cccc-dddddddddddd', 'Edit',
   '{"file_path":"/repo/proj/.claude/worktrees/wt9/internal/codex.go"}', NULL);

-- WT's mechanical-signature decision (gate-pass) + an exploration-mode contrast follows.
INSERT INTO tool_calls (id, session_id, tool_name, input_json, result_content) VALUES
  (57, 'claude:77777777-aaaa-bbbb-cccc-dddddddddddd', 'ExitPlanMode',
   '{"plan":"Implement issue 42 in the worktree."}',
   'User has approved your plan. You can now start coding.');

-- ============================================================================
-- MODE-CLASSIFICATION (#324, G / AC-8a). Two designed tracks plus an explicit
-- neither-dominant track. The classifier groups by git_branch and scores the signatures.
--   Track `issue-feed`   — MECHANICAL: gate-pass decision, worktree/work-on-issue loop
--                          markers, code edits, ZERO veto.
--   Track `landing-copy` — EXPLORATION: multiple [Request interrupted / doesn't-want vetoes,
--                          a rejected/cancelled decision, `.md` content edits.
--   Track `mixed-bag`    — NEITHER dominant: one veto + one passed decision + one `.md` edit,
--                          so neither score wins by the margin → `unlabeled` (generic
--                          book-keeping, never a guessed automation pitch).
-- (WT's `issue-42` is ALSO mechanical — a second mechanical track — but the G assertion
-- pins the three designed tracks explicitly.)
-- ============================================================================

-- MECHANICAL track `issue-feed` — two sessions.
INSERT INTO sessions VALUES
  ('claude:91111111-1111-1111-1111-111111111111', 'proj', 'claude',
   '/repo/proj', 'issue-feed',
   '/u/.claude/projects/-repo-proj/91111111.jsonl',
   '2026-06-07', '2026-06-07', 'Drive the issue-feed renderer via the work-on-issue loop.', 5, 2),
  ('claude:92222222-2222-2222-2222-222222222222', 'proj', 'claude',
   '/repo/proj/.worktrees/issue-feed', 'issue-feed',
   '/u/.claude/projects/-repo-proj-.worktrees-issue-feed/92222222.jsonl',
   '2026-06-07', '2026-06-07', 'Continue the issue-feed worktree implementation.', 4, 2);
INSERT INTO tool_calls (id, session_id, tool_name, input_json, result_content) VALUES
  (60, 'claude:91111111-1111-1111-1111-111111111111', 'AskUserQuestion',
   '{"questions":[{"header":"Reindex strategy","question":"Incremental vs full reindex?"}]}',
   'Your questions have been answered: "Incremental vs full reindex?"="incremental"'),
  (61, 'claude:91111111-1111-1111-1111-111111111111', 'Edit',
   '{"file_path":"/repo/proj/src/feed.ts"}', NULL),
  (62, 'claude:92222222-2222-2222-2222-222222222222', 'Edit',
   '{"file_path":"/repo/proj/internal/feed/index.go"}', NULL);

-- EXPLORATION track `landing-copy` — two sessions, prose edits + vetoes + a rejected path.
INSERT INTO sessions VALUES
  ('claude:a3333333-3333-3333-3333-333333333333', 'proj', 'claude',
   '/repo/proj', 'landing-copy',
   '/u/.claude/projects/-repo-proj/a3333333.jsonl',
   '2026-06-07', '2026-06-07', 'Draft the landing hero copy; try a few framings.', 7, 4),
  ('claude:a4444444-4444-4444-4444-444444444444', 'proj', 'claude',
   '/repo/proj', 'landing-copy',
   '/u/.claude/projects/-repo-proj/a4444444.jsonl',
   '2026-06-07', '2026-06-07', 'Rework the story section; the last direction was wrong.', 6, 3);
INSERT INTO tool_calls (id, session_id, tool_name, input_json, result_content) VALUES
  (63, 'claude:a3333333-3333-3333-3333-333333333333', 'AskUserQuestion',
   '{"questions":[{"header":"Hero framing","question":"Hero-vs-story framing?"}]}',
   'The user doesn''t want to proceed with this tool use. The tool use was rejected.'),
  (64, 'claude:a3333333-3333-3333-3333-333333333333', 'Write',
   '{"file_path":"/repo/proj/content/hero.md"}', NULL),
  (65, 'claude:a4444444-4444-4444-4444-444444444444', 'Write',
   '{"file_path":"/repo/proj/content/story.md"}', NULL);

-- NEITHER-DOMINANT track `mixed-bag` — one session, balanced signals → unlabeled.
INSERT INTO sessions VALUES
  ('claude:b5555555-5555-5555-5555-555555555555', 'proj', 'claude',
   '/repo/proj', 'mixed-bag',
   '/u/.claude/projects/-repo-proj/b5555555.jsonl',
   '2026-06-07', '2026-06-07', 'Some odds and ends across the repo.', 4, 2);
INSERT INTO tool_calls (id, session_id, tool_name, input_json, result_content) VALUES
  (66, 'claude:b5555555-5555-5555-5555-555555555555', 'AskUserQuestion',
   '{"questions":[{"header":"Odds and ends","question":"Which loose end first?"}]}',
   'Your questions have been answered: "Which loose end first?"="the docs"'),
  (67, 'claude:b5555555-5555-5555-5555-555555555555', 'Write',
   '{"file_path":"/repo/proj/docs/notes.md"}', NULL);

-- ----------------------------------------------------------------------------
-- Veto + loop markers in the message stream (interruption + mechanical-loop signals,
-- prose-read). Session B carries the original veto; the G tracks carry their signatures.
-- ----------------------------------------------------------------------------
INSERT INTO messages VALUES
  (1, 'claude:aaaaaaaa-1111-2222-3333-444444444444', 'user', 'Pick up the parser refactor and ship it.'),
  (2, 'claude:bbbbbbbb-5555-6666-7777-888888888888', 'user', 'Now wire up the regression suite.'),
  (3, 'claude:bbbbbbbb-5555-6666-7777-888888888888', 'user', '[Request interrupted by user]'),
  (4, 'claude:cccccccc-9999-aaaa-bbbb-cccccccccccc', 'user', 'Build the feature behind a worktree.'),
  (5, 'codex:ffffffff-8888-9999-aaaa-bbbbbbbbbbbb', 'user', 'A codex session in this repo; cwd unrecorded.'),
  -- WT (issue-42): worktree loop marker, no veto → reinforces mechanical.
  (6, 'claude:77777777-aaaa-bbbb-cccc-dddddddddddd', 'user', 'Run the work-on-issue loop for issue 42 in its worktree.'),
  -- issue-feed (mechanical): worktree/work-on-issue loop markers, no veto.
  (7, 'claude:91111111-1111-1111-1111-111111111111', 'user', 'Drive the work-on-issue loop in the worktree.'),
  (8, 'claude:92222222-2222-2222-2222-222222222222', 'user', 'Continue the worktree implementation.'),
  -- landing-copy (exploration): repeated vetoes / doesn't-want-to-proceed steering.
  (9, 'claude:a3333333-3333-3333-3333-333333333333', 'user', '[Request interrupted by user]'),
  (10, 'claude:a3333333-3333-3333-3333-333333333333', 'user', 'doesn''t want to proceed — try a warmer tone'),
  (11, 'claude:a4444444-4444-4444-4444-444444444444', 'user', '[Request interrupted by user] rethink the framing'),
  -- mixed-bag (neither dominant): a single veto, balancing its one passed decision + one .md edit.
  (12, 'claude:b5555555-5555-5555-5555-555555555555', 'user', '[Request interrupted by user]');
