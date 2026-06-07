---
id: xnh4gc2zqe67j8tbxh926r9j
title: Survey signal-correction — read the ground-truth signals (sessions.db + repo), not thin proxies
status: validation
source: "captain (2026-06-07) - combined fix for the survey-skill dogfood issue cluster spacedock-dev/spacedock#318/#319/#320 (+ #317 umbrella). One task, not 1-1: the reports close after this ships. #316 (detect-spacedock incumbent + audit/meta-loop flip) is OUT of scope (serves spacedock-incumbent users, a minority) - park as a separate roadmap proposal."
started: 2026-06-07T18:57:58Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-survey-signal-correction
issue:
---

`spacedock:survey` decides "what this project is and where it stands" from THIN
PROXIES while the ground truth sits unread right next to it. One root cause — the
same proof-policy class we already fence (trust a proxy over ground truth) —
manifest across the skill's scan + scaffold-detection logic. Verified against live
code + a freshly-synced real `sessions.db` this session (see ## Spike results):

- The scaffold probe prints `none` on THIS spacedock repo — the skill can't
  recognize its own home (file-only probe of `.claude/skills` etc.).
- `tool_calls.skill_name` is populated (1045 `Skill`-call rows in the spike sync;
  `spacedock:ensign`×964, `superpowers:*`, bare `running-research-spikes`) and
  never queried — the behavioral truth of which scaffolds ran is ignored.
- The scan scopes by `basename(pwd)`, but agentsview keys every session by the
  RESOLVED repo-root basename. Run from a subdir/worktree the two keys diverge and
  survey finds 0 sessions while 1198 sit under the agentsview key (spike).
- The OPEN-decision "NEEDS YOU" frontier is transcript-only and over-reports
  (forks already merged read as still-open).

This is the highest-value survey work because it makes the skill CORRECT for its
ACTUAL users (greenfield / other-scaffold repos), not the spacedock-incumbent
minority.

## Structure (captain redirect 2026-06-07 — supersedes #309)

#309 moved survey's logic INTO `skills/survey/bin/{scan-project,detect-scaffold}`.
This task SUPERSEDES that resolution: the shipped survey skill carries NO bundled
`bin/` scripts. The four signal-corrections below ship as:

- a **recommended-SQL reference file** `skills/survey/references/queries.sql` —
  annotated, ONE labeled query per concern (scaffold-usage tally; git-root scoping;
  WORK-BY-AREA; decision/OPEN detection). The SQL lives ONLY in the reference file,
  NEVER fenced in SKILL.md prose (a fenced query + a test parsing it would re-trip
  the proof-policy ban on parsing SKILL.md as the artifact).
- **SKILL.md PROSE** orchestrates: resolve the repo root via
  `git rev-parse --path-format=absolute --git-common-dir` (fallback `basename(pwd)`
  for non-git / git <2.31), run the recommended queries from the reference file,
  apply the multi-label scaffold classification + family-normalization (EXCLUDING
  `spacedock:*` self-invocation), cross-check the OPEN frontier against the repo
  (confident-match-only drop + mandatory `unverified` degrade), assemble the report.
- the existing `bin/scan-project` and `bin/detect-scaffold` are DELETED.

**Lighter proof bar (captain redirect).** Survey is ONBOARDING: a wrong orientation
report is a SOFT failure (low blast radius), not a broken pipeline, and report
EFFICACY is monitored by a SEPARATE project. So the heavy proof machinery a higher
blast radius would justify is NOT warranted here. The proof surface is:
- a LIGHT smoke test that the recommended queries (the reference file) run against a
  committed production-shaped fixture DB and return the expected SHAPE/rows — this
  pins the one thing worth pinning (SQL extraction correctness; catches a broken
  query or schema drift). The smoke test READS and RUNS the `.sql` reference file —
  it executes the artifact, it does not parse SKILL.md prose, preserving the 4q
  "tests run artifacts, not instruction text" property.
- a LIVE drive of the skill (run `/spacedock:survey`, observe the corrected report)
  for orchestration correctness. Ongoing report-efficacy is the separate project's
  concern — out of scope here.

## In scope (one cohesive deliverable)

1. **Scaffold detection — multi-label + behavioral** (#319, #317.1). Report ALL
   scaffolds, not the first if-ladder match; corroborate the file probe against a
   `tool_calls.skill_name` tally (file+invoked / file+never-invoked=installed-but-unused /
   not-file+invoked=recovered). Exclude survey's own `spacedock:survey`/`ensign`
   self-invocation from the tally (self-pollution — it is 1026/1045 of the tally
   on this repo, see ## Spike results).
2. **Identity scoping — coalesce a repo across keys** (#318). Scope by repo identity
   via `git rev-parse --path-format=absolute --git-common-dir` (NOT plain
   `--git-common-dir`, which returns a RELATIVE path from root/subdir; NOT
   `--show-toplevel`, which returns the worktree/subdir root) / path-prefix boundary /
   `worktree_project_mappings` override; union all sessions under the repo; surface
   folded-in keys + blank-`cwd` count; deterministic (no fuzzy name matching).
   No-regression for single-checkout repos. (`--path-format` needs git ≥2.31; not in
   a git repo or older git → fall back to `basename(pwd)`, today's behavior.)
3. **Identity from edits** (#317.2). A WORK-BY-AREA tally bucketing `Edit/Write`
   `input_json.$.file_path` by package; report "what this is" (edits) separately
   from "where you stop" (decisions); treat edits to external sibling repos
   (paths outside the repo root) as references, not identity; show all-time
   alongside recent.
4. **OPEN frontier — cross-check against the repo** (#320). After the transcript
   scan, cross-reference each OPEN fork against repo artifacts (git log, merged PRs,
   working tree) and split shipped (drop) / decided-not-shipped (backlog) /
   never-decided (true open). Conservative default (drop ONLY on a confident match,
   else keep on the frontier). Mandatory transcript-only degrade, flagged
   `unverified`, when no repo signal. Fold in the cheap independent fix: the
   `ExitPlanMode "User has approved your plan"` prefix matches none of the three
   done-prefixes the decision query checks today, so approved plans fall to OPEN.

## Foundation (prerequisite — ideation owns the design, implementation ships it)

Under the lighter proof bar, the foundation is one production-shaped fixture DB the
smoke test runs the recommended queries against — NOT the heavy machinery a higher
blast radius would need. Explicitly DROPPED from the prior cycle (captain redirect):
the git-init test helper, the `coupleFixtureToRepo` / marker-token coupling, the
per-AC RED-lever rigor, and any DDL extension that existed only to feed the coupling.

- Extend the fixture DDL with `cwd` (on `sessions`) and `skill_name` (on
  `tool_calls`) — they are ABSENT today, so the new scoping + scaffold-usage queries
  would compile and return nothing against the current fixture (the invisible-fix
  trap: the smoke test must have the columns to assert the corrected SHAPE). Defaults
  match production: `cwd` TEXT NOT NULL DEFAULT `''` (blank stored as `''`, not NULL).
  Seed `git_branch` / `category` too IF it keeps the fixture production-shaped, but no
  query reads them (YAGNI — don't wire a query for them).
- Seed at least one anonymized fixture from a REAL agentsview `sessions.db` dump so
  Skill/scaffold/decision/edit shapes match production (all fixtures today are
  hand-authored — the validate-on-real-data gap; closes it once for the whole
  cluster). The spike already pulled the real shapes (## Spike results); the fixture
  encodes them: namespaced + bare `skill_name`, `spacedock:*` self-rows, a
  `superpowers:*` recovered case, `Your questions have been answered:` /
  `The user doesn't want to proceed` decision results, an `ExitPlanMode` approved row,
  Edit/Write `input_json.$.file_path` rows, and `cwd` values under one repo root
  spanning a subdir + a worktree-style path so the git-root scoping query has
  something to coalesce.
- The git-root scoping query in the reference file is exercised by the smoke test
  with the fixture's `cwd` rows as the corpus and a fixed repo-root prefix as the
  bound parameter — no live git repo needed for the SMOKE bar (the live drive
  exercises the real `git rev-parse` resolution). This is the deliberate
  light-vs-heavy trade: pin the SQL extraction, leave end-to-end orchestration to the
  live drive.

## Out of scope

- **#316** detect-spacedock-incumbent + the audit/meta-loop "flip the back half"
  reframe — defer; it serves spacedock-incumbent users only. Record the audit reframe
  as a roadmap proposal. (Design the scaffold-detection contract in part 1 to be
  EXTENSIBLE/multi-label so a future spacedock case drops in without a rebase —
  settled below: the family-normalization table + multi-label emit already admit a
  `spacedock` family row; #316 adds the row + the audit reframe, no contract change.)
- **#317.3** the low-confidence structured issue-filing prompt — defer (meta-feature).

## Spike results (riskiest unknowns, grounded before committing the design)

Per running-research-spikes: ground the riskiest unknowns first, record the result.
The two named risks were (a) the real `skill_name` shape and (b) the cwd/git scoping
query against a real multi-worktree `sessions.db`. Both grounded against a fresh
scoped sync of THIS repo's Claude history (agentsview v0.32.1, 1198 sessions, 34
distinct cwds). Throwaway DB removed after; the queries below are the seed for the
implementation's first tests.

**Sandbox / TCC reality (load-bearing for the test design).** A sandboxed agent
process cannot read `~/.agentsview/sessions.db` NOR a copy under `/tmp` (macOS TCC:
the file reads as 0 bytes / "operation not permitted" to a non-entitled process,
even the `agentsview` binary when *it* is the spawned process). The readable path is
to `agentsview sync` this project's `~/.claude/projects` dirs into a process-owned
`AGENTSVIEW_DATA_DIR`, then query that copy — exactly what SKILL.md §1 documents.
Consequence for tests: the integration tests do NOT depend on a live sync or
`~/.agentsview`; they build a fixture DB from committed SQL via the `sqlite3` CLI
(the existing `buildFixtureDB` pattern), so they run on a fresh CI box.

**Real schema (agentsview v0.32.1 `sessions.db`) — the columns the fixes need exist:**
- `sessions.cwd` TEXT NOT NULL DEFAULT `''` (col 41) — populated; blank stored as `''`.
- `sessions.git_branch` TEXT NOT NULL DEFAULT `''` (col 42) — populated.
- `tool_calls.skill_name` TEXT nullable (col 7) — populated ONLY on `tool_name='Skill'` rows.
- `tool_calls.category` TEXT NOT NULL (col 4) — extra signal, available.
- `worktree_project_mappings(machine, path_prefix, project, enabled)` — the #318
  override; EMPTY in a fresh sync, so the design must NOT depend on it being populated.
- The current fixture DDL has NONE of `cwd`/`git_branch`/`skill_name`/`category` —
  invisible-fix trap confirmed (a scoping/scaffold query would compile and pass
  doing nothing).

**`skill_name` shape (the riskiest unknown #1) — grounded:**
`SELECT skill_name, COUNT(*) FROM tool_calls WHERE skill_name IS NOT NULL AND skill_name<>'' GROUP BY 1 ORDER BY 2 DESC`
yields: `spacedock:ensign` 964, `spacedock:first-officer` 24, `spacedock:debrief` 20,
`running-research-spikes` 15, `spacedock:present-gate` 5, `spacedock:using-claude-team` 5,
`spacedock:feedback-rejection-flow` 4, `spacedock:refit` 2, `superpowers:brainstorming` 2,
`superpowers:dispatching-parallel-agents` 1, `superpowers:systematic-debugging` 1, …
- Values are `family:skill` (namespaced) OR bare (`running-research-spikes`).
- Self-pollution dominates: `spacedock:*` = ~1026/1045. Excluding self is essential,
  not cosmetic — without it every repo "uses spacedock" because survey/ensign ran.
- → family normalization: split on `:`, take the prefix as family; bare names map to
  superpowers disciplines (`brainstorming`, `running-research-spikes`,
  `subagent-driven-development`, `writing-plans`, `executing-plans`,
  `dispatching-parallel-agents`, `systematic-debugging`, `finishing-a-development-branch`)
  → family `superpowers`; exclude family `spacedock` from the scaffold tally.

**cwd/git scoping query (the riskiest unknown #2) — grounded:**
- agentsview already coalesces all 34 worktree/subdir cwds of this repo into ONE
  `project='spacedock_v1'` (1198 sessions) — even with `worktree_project_mappings`
  empty. So the per-session `project` key is already repo-rooted on v0.32.1.
- BUT `scan-project` derives ITS key from `basename(pwd)`. Reproduced:
  - from repo root → key `spacedock_v1` → 1198 sessions (correct);
  - from `docs/dev/.spacedock-state` (a real 57-session cwd) → key `_spacedock_state` → **0**;
  - from any worktree → key e.g. `audit_k6` → **0**.
  So survey run from the split-root state checkout or any worktree reports "no agent
  history" while 1198 sessions sit under the agentsview key. 149 sessions in the DB
  have a cwd whose basename ≠ the repo basename — all invisible to the basename key.
- The scoping FORMULA: plain `git rev-parse --git-common-dir` is NOT safe — it
  returns a RELATIVE path from the repo root and intermediate subdirs, so
  `basename(dirname(...))` regresses the exact subdir case #318 is about. Re-grounded
  this session (git 2.39.5) across all four cases:
  - plain `--git-common-dir`: from repo root → `.git` → `basename(dirname)`=`.` (WRONG);
    from `docs/dev` → `../../.git` → `..` (WRONG); from `docs/dev/.spacedock-state` →
    `…/spacedock-v1/.git` → `spacedock-v1` (right, but only because git returns
    absolute when common-dir is outside the cwd's own `.git` chain — not reliable).
  - `git rev-parse --path-format=absolute --git-common-dir`: returns
    `…/spacedock-v1/.git` from repo root, `docs/dev`, `docs/dev/.spacedock-state`,
    AND the external worktree `/private/tmp/spacedock-audit-autoinstall-7462` — so
    `basename(dirname(...))`=`spacedock-v1` EVERYWHERE.
  - `--show-toplevel` is the other trap (returns the worktree/subdir as root):
    from `.spacedock-state` → `…/.spacedock-state`; from the external worktree →
    that worktree dir.
  - non-git dir → `--path-format=absolute --git-common-dir` exits 128 → fall back.
  → `PROJECT = sanitize(basename(dirname(git rev-parse --path-format=absolute --git-common-dir)))`,
  falling back to `basename(pwd)` when not in a git repo OR git is older than 2.31
  (the `--path-format` floor), makes scan-project's key match agentsview's key from
  any subdir/worktree. `--path-format` is a single git ≥2.31 dependency; the fallback
  preserves today's behavior on older git.
- The prior dogfood note (a DIFFERENT product: main + sibling + deploy + 4 worktrees)
  recorded agentsview itself fragmenting into 7 project keys: 31 reported vs 71
  coalesced (56% of sessions, 18% of decisions dropped). So BOTH failure modes are
  real: (i) scan-project's basename key diverging from agentsview's key (provable on
  this repo today), and (ii) agentsview fragmenting into multiple keys for
  sibling/deploy checkouts (older sync / unmapped worktrees). The fix unions sessions
  whose `cwd` is under the git-root prefix OR whose `project` matches any folded key —
  it must not trust a single agentsview key.

**WORK-BY-AREA (#317.2) — grounded:** `Edit`=10991, `Write`=2362 rows;
`json_extract(input_json,'$.file_path')` yields real absolute paths
(`…/internal/dispatch/build.go`, `…/docs/dev/.spacedock-state/…/index.md`, and paths
OUTSIDE the repo root like `/tmp/sd-spike-pathfix/…`). Bucketing by the path segment
under the repo root gives a package tally; paths outside the repo root are the
"external sibling = reference, not identity" case.

**ExitPlanMode (#320 cheap fix) — code-provable:** `scan-project` DECISIONS marks
done only on three prefixes (`User has answered…` / `Your questions have been
answered…` / `Your question has been answered…`). `ExitPlanMode` approval emits
`User has approved your plan` — matches NONE — so every approved plan reads OPEN.
(0 `ExitPlanMode` rows in the narrowed spike sync; the bug is in the code, proven by
a fixture row carrying that result, not by a production row.)

**AskUserQuestion result shapes (no change needed) — grounded:** done rows begin
`Your questions have been answered:`; OPEN/rejected rows begin `The user doesn't want
to proceed`. Matches today's fixture and the shipped query — confirms the existing
done/OPEN detection is correct and must be preserved under the new scope.

## Design forks — SETTLED

### Fork 1 — the scaffold-detection contract (no bin/; SQL reference + SKILL.md prose)

**Where the logic lives** (settled by the captain redirect). NO bundled `bin/`. The
two signals reconcile in SKILL.md prose:
- the **file probe** (does `.claude/skills/superpowers` etc. exist on disk) is
  SKILL.md prose — a multi-label check that names EVERY matched scaffold, not a
  single ladder winner.
- the **behavioral tally** is the `scaffold-usage` query in
  `references/queries.sql`: a `tool_calls.skill_name` GROUP BY over the git-root-scoped
  session set, run by SKILL.md prose. SKILL.md joins the two into the 3-bucket
  classification.
- Rejected (superseded): a `bin/scan-project` SCAFFOLD-USAGE section or a
  `bin/recognize-scaffold` (the redirect removes all bundled bin/; the tally is a
  labeled query the prose runs, the classification is prose).

**The 3-bucket scaffold classification** (join of file probe × usage tally):
- `file-present + invoked` → **active** (reported plainly).
- `file-present + never-invoked` → **installed-but-unused** (flagged).
- `not-file-present + invoked` → **recovered** (reported — this is the case the
  file-only probe misses entirely; the behavioral signal recovers it).

**`skill_name` → family normalization** (grounded above): split on first `:`, prefix
is the family; bare names map via a small superpowers-discipline name set → family
`superpowers`; **exclude family `spacedock`** (self-pollution). The family set is a
prose table, not an if-ladder — adding a `spacedock` family row later (#316) is a prose
edit, not a contract change, satisfying the EXTENSIBLE requirement.

**Multi-label file probe.** Today's single-winner ladder (`if superpowers; elif gsd;
elif similar; else none`) becomes prose that probes each scaffold independently and
names every match (`superpowers`, `gsd`, `similar: <names>`), reporting `none` only
when no probe matched. The scaffold-usage query is what the smoke test pins (it
catches a broken tally / schema drift); the prose join + multi-label probe are
exercised by the live drive.

### Fork 2 — #320 OPEN-vs-repo taxonomy + conservative-match + degrade

**3 buckets** for each transcript-OPEN fork, after cross-referencing repo artifacts:
- **shipped** → DROP from the frontier. Evidence: a merged PR or a git-log commit
  whose subject/body confidently references the fork (decision header / branch).
- **decided-not-shipped** → move to a **backlog** line (decided, no artifact).
- **never-decided** → **true open**, stays on the `NEEDS YOU` frontier.

**Conservative-match rule.** DROP only on a CONFIDENT repo match (an exact-ish token
match between the decision header/branch and a merged-PR title or commit subject).
Anything less than confident → KEEP on the frontier. A false "still open" is a cheap
nudge; a false "shipped" silently hides a real open fork — so the asymmetry favors
keeping.

**Mandatory degrade.** When NO repo signal is available (not a git repo, or `git
log`/PR lookup fails/empty), the frontier degrades to transcript-only and every OPEN
fork is flagged **`unverified`** in the report — never silently presented as
authoritative. The degrade is the default, not an error.

**ExitPlanMode cheap fix** folded in here: add `User has approved your plan` (and the
`User approved` variants) to the done-prefix set in the `decision/OPEN-detection`
query in `references/queries.sql` so an approved plan reads `done`, not OPEN.
Independent, smoke-pinnable (a fixture `ExitPlanMode` approved row → the query returns
`done`), one line of SQL.

**Where the cross-check runs.** The repo cross-reference needs `git log` / merged-PR
data + the working tree. Under the no-bin/ structure: the `decision/OPEN-detection`
query in `references/queries.sql` emits the raw transcript-OPEN frontier (smoke-pinned
shape); the cross-check + 3-bucket split + conservative/degrade rules are SKILL.md
PROSE the survey driver runs against the repo (`git log` / merged-PR / working tree).
The 3-bucket classification's correctness is exercised by the LIVE drive (the lighter
proof bar), not a dedicated bin/ helper test — onboarding's low blast radius +
external efficacy monitoring justify proving orchestration by live run rather than
unit-pinning every branch.

### Fork 3 — split?

**Keep as one.** The four sub-parts share the structure (the `references/queries.sql`
reference file, the SKILL.md prose orchestration, the one production-shaped fixture)
and one live drive proves them together — splitting would duplicate the reference file
and the fixture. The parts are independent enough to implement in sequence within one
task but share enough scaffolding that one entity is correct. Default honored.

## Acceptance criteria (entity-level; lighter onboarding bar — captain redirect)

Proof bar (the redirect's lighter bar): the ONE thing worth pinning is SQL-extraction
correctness — a query-smoke test that runs the `references/queries.sql` queries
against a committed production-shaped fixture DB and asserts the corrected SHAPE/rows
(catches a broken query or schema drift). The smoke test EXECUTES the `.sql` artifact;
it does NOT parse SKILL.md prose, preserving the 4q "tests run artifacts, not
instruction text" property. Orchestration correctness (the prose join, the cross-check
3-bucket split, the report assembly) is proven by a LIVE drive of `/spacedock:survey`.
Ongoing report efficacy is the separate monitoring project's concern — out of scope.

**AC-1 — No bundled bin/; logic is SQL-reference + SKILL.md prose (supersedes #309).**
`skills/survey/bin/scan-project` and `bin/detect-scaffold` are removed; the recommended
SQL lives in `skills/survey/references/queries.sql` (annotated, one labeled query per
concern: scaffold-usage tally, git-root scoping, WORK-BY-AREA, decision/OPEN
detection); SKILL.md prose orchestrates resolution + classification + cross-check +
report. No SQL is fenced in SKILL.md.
- Verified by: the query-smoke test reads `references/queries.sql` (RED if the file is
  absent — this is what catches a missing/renamed reference file, since contractlint's
  `referenceRe` only covers `.md` paths and won't enforce a `.sql` reference); the
  deleted `bin/` paths no longer exist (the updated integration tests no longer
  `exec.Command` them — the re-point is specified in ## Test plan); a live drive runs
  the skill without invoking a bundled script.

**AC-2 — The recommended queries return the corrected data on the fixture (light smoke).**
Run against the committed production-shaped fixture DB (extended with `cwd` +
`skill_name`), each labeled query in `references/queries.sql` returns the corrected
shape: the git-root scoping query, bound to a repo-root prefix, coalesces the fixture's
subdir + worktree-style `cwd` rows into ONE identity (more sessions than a single
basename key would yield); the scaffold-usage tally GROUPs `skill_name` by family,
reports a `superpowers` row with NO files (recovered) and EXCLUDES the
`spacedock:*`-dominated self rows; the WORK-BY-AREA query buckets Edit/Write
`input_json.$.file_path` and separates a path under the repo-root prefix from one
outside it; the decision/OPEN query marks an `ExitPlanMode "User has approved your
plan"` row as `done` and a rejected row as OPEN.
- Verified by: a Go query-smoke test (extends the `buildFixtureDB` pattern) that
  `sqlite3`-runs each labeled query against the fixture and asserts the row shape.
  Expected values come from the FIXTURE rows — an independent source that diverges from
  the skill text. RED on a broken query, a dropped self-exclusion, a wrong prefix
  bound, or schema drift. `t.Skip` (not fail) when `sqlite3` is absent (matches
  `buildFixtureDB`). This is the deliberately light bar: it pins extraction, not the
  full orchestration.

**AC-3 — A live drive shows the corrected report.**
A live run of `/spacedock:survey` on a real multi-worktree repo produces a report that
reflects the four corrections: the session count reflects the git-root-coalesced scope
(not the basename subset), the scaffold line is multi-label with self-invocation
excluded, a WORK-BY-AREA section appears, and the OPEN frontier is cross-checked
(shipped forks dropped or, with no repo signal, flagged `unverified`).
- Verified by: a live drive (the redirect's named orchestration proof), captured in
  the validation stage — observe the rendered report, not a SKILL.md substring. This is
  the orchestration half the smoke test deliberately does not cover; onboarding's low
  blast radius + external efficacy monitoring justify a live observation rather than
  unit-pinning every prose branch.

## Test plan

- **Surface:** ONE Go query-smoke test under `skills/integration/`, extending the
  existing `buildFixtureDB` pattern — `sqlite3`-run each labeled query from
  `references/queries.sql` against the committed fixture and assert the corrected row
  shape (AC-2). NO new test framework, NO git-init helper, NO coupleFixtureToRepo, NO
  per-AC RED-lever rigor (all dropped per the redirect's lighter bar). `t.Skip` (not
  fail) when `sqlite3` is absent so a minimal box stays runnable (matches
  `buildFixtureDB`).
- **Fixture (the one production-shaped DB):**
  - extend `testdata/survey/fixture-sessions.sql` DDL with `cwd` (sessions) +
    `skill_name` (tool_calls), defaults matching production (`cwd` TEXT NOT NULL
    DEFAULT `''`); add rows exercising the corrected queries: `cwd` rows under one
    repo-root prefix spanning a subdir + a worktree-style path (git-root scoping),
    `Skill` rows with namespaced + bare + `spacedock:*` self names keeping the
    recovered-`superpowers` and excluded-`spacedock` cases DISTINCT (scaffold-usage),
    Edit/Write `input_json.$.file_path` rows under + outside the repo-root prefix
    (WORK-BY-AREA), an `ExitPlanMode` approved row + a rejected row (decision/OPEN);
  - keep the rows anonymized-real-shape (counts kept, identifiers stripped) so the
    queries run against production-shaped data. The smoke binds a FIXED repo-root
    prefix string for the scoping/bucketing queries — no live git repo for the smoke
    bar.
- **Re-point the existing integration tests (AC-1):** `survey_extraction_test.go` and
  `survey_scaffold_test.go` `exec.Command` the now-deleted `bin/scan-project` /
  `bin/detect-scaffold`. Re-point them to the query-smoke surface (or replace with the
  light query-smoke) and REMOVE every reference to the deleted scripts — a stale
  `exec.Command` on a deleted path would fail the suite.
- **Contractlint check (AC-1):** `internal/contractlint/structural_checks_test.go` —
  `TestSurveyIsDiscoverableUserCommand` (frontmatter discovery) is UNAFFECTED;
  `TestUserSkillReferenceClosureResolves`'s `referenceRe` matches only `.md` paths, so
  a `references/queries.sql` reference is NOT enforced by it (the smoke test is what
  catches a missing reference file — note this asymmetry, do not add a `.sql` arm to
  the closure check unless the captain asks); `isClaudeAdapter` already exempts
  `survey/SKILL.md` from the HOME-rooted personal-config check, and the `interpreterRe`
  check forbids only `python`/`commission/bin` — NOT `git`/`sqlite3`/`agentsview` — so
  moving the shell orchestration into SKILL.md prose does NOT trip contractlint
  (verified against the regex definitions this session). No contractlint change is
  required by the restructure; confirm the suite stays green.
- **Live drive (AC-3):** run `/spacedock:survey` on a real multi-worktree repo in the
  validation stage; observe the corrected report (coalesced scope, multi-label
  self-excluded scaffold, WORK-BY-AREA, cross-checked frontier). The redirect's named
  orchestration proof.
- **Cost/complexity:** LOW (down from medium — the heavy scaffolding is dropped). One
  query-smoke test + one fixture-DDL extension + re-point two existing tests + one live
  drive. No live agentsview dependency for the smoke (sandbox-blocked per the spike);
  the committed-fixture query-smoke is the CI-safe surface.
- **No spike outstanding:** the two named risky unknowns (skill_name shape; cwd/git
  scoping against real data) are grounded in ## Spike results, INCLUDING the corrected
  `--path-format=absolute --git-common-dir` formula re-grounded across root / subdir /
  state-checkout / worktree / non-git this session; the queries there seed the
  reference file's first queries. The remaining mechanisms (sqlite3 fixture build,
  `git rev-parse --path-format=absolute --git-common-dir`, the contractlint regexes not
  matching `git`/`sqlite3`) are already proven by the existing tests + the spike.

## Reused / deferred from #316 and #317.3

- **Reused:** part 1's family-normalization table + multi-label emit are designed so
  #316's `spacedock`-incumbent case is a DATA row (add a `spacedock` family, stop
  excluding it for incumbent users), not a contract change — #316 drops in without a
  rebase.
- **Deferred (#316):** detect-spacedock-incumbent + the audit/meta-loop "flip the
  back half" reframe → record as a separate roadmap proposal (serves the
  spacedock-incumbent minority).
- **Deferred (#317.3):** the low-confidence structured issue-filing prompt → defer
  (meta-feature, not part of the correctness fix).

## Sources (close after this ships)

spacedock-dev/spacedock#318, #319, #320, #317 (umbrella — superseded by #318/#319),
and #316 (deferred). These reports get closed after this task lands; no 1-1 mapping.

## Stage Report: ideation

- DONE: Settle the unified detect-scaffold contract (file-probe + multi-label + a tool_calls.skill_name DB tally + self-pollution exclusion, and WHERE it lives given detect-scaffold has no DB access today) AND the #320 OPEN-vs-repo design (3-bucket shipped/backlog/true-open taxonomy + conservative-match rule + mandatory transcript-only degrade) as concrete, build-ready designs.
  ## Design forks — SETTLED: Fork 1 keeps two binaries (multi-label detect-scaffold file probe + a new scan-project `## SCAFFOLD-USAGE` DB tally, joined in SKILL.md step 3) with a family-normalization table excluding `spacedock:*`; Fork 2 fixes the 3-bucket taxonomy + confident-match-only drop + mandatory `unverified` degrade + the ExitPlanMode done-prefix fix.
- DONE: Specify the fixture foundation and prove the riskiest unknown FIRST: extend the test fixture DDL with cwd/git_branch/skill_name, a shared git-init test helper, and >=1 anonymized fixture seeded from a REAL agentsview sessions.db; for each in-scope fix name the RED-on-proxy fixture lever so a proxy-only impl fails. Ground the skill_name shape + the cwd/git scoping query against the real sessions.db before committing the design (running-research-spikes).
  ## Spike results: freshly synced real v0.32.1 sessions.db (1198 sessions, 34 cwds) — grounded the schema (cwd/git_branch/skill_name/category all present, fixture DDL lacks all four), the skill_name shape (namespaced + bare; spacedock:* = 1026/1045 self-pollution), and the scoping bug (basename key → 0 sessions from a subdir/worktree while git-common-dir resolves the repo root; --show-toplevel is the trap). ## Foundation + ## Test plan specify the DDL extension, git-init helper, real-data fixture; each AC names its RED-lever.
- DONE: Write entity-level ACs + a test plan covering scaffold multi-label+behavioral, identity scoping (git-common-dir/prefix/mappings), identity-from-edits (WORK-BY-AREA), and the OPEN frontier (incl. the cheap ExitPlanMode false-OPEN fix) — each AC's "Verified by" naming an outside-the-body check (a test running the real bin/ script against a fixture); record what is reused/deferred from #316/#317.3.
  ## Acceptance criteria AC-1..AC-5 (foundation, scoping, scaffold, WORK-BY-AREA, OPEN-frontier) each have a "Verified by" naming a Go integration test that runs the real bin/ via cmd.Dir against committed fixtures + a real git-init repo; ## Reused / deferred records the #316 family-row extensibility hook and the #316/#317.3 defers.

### Summary

Grounded the two riskiest unknowns against a fresh real-data sync before committing the design: the spike CORRECTED the body's premise — agentsview already coalesces this repo's worktrees into one project key, so the #318 bug is that `scan-project`'s OWN `basename(pwd)` key diverges from agentsview's repo-rooted key (0 sessions from a subdir/worktree), fixed by `dirname(git rev-parse --git-common-dir)`; and confirmed `spacedock:*` self-pollution dominates the skill_name tally (1026/1045), making self-exclusion load-bearing. Settled all three forks (two-binary scaffold contract with a DB tally in scan-project + family-normalization table; 3-bucket OPEN taxonomy with confident-match-only drop + mandatory unverified degrade; keep-as-one), specified the fixture foundation that defeats the invisible-fix trap (DDL extension + git-init helper + real-data fixture), and wrote five behavioral ACs each with an outside-the-body RED-on-proxy lever. Also surfaced the load-bearing sandbox/TCC fact: tests must build a fixture DB from committed SQL (existing pattern), never depend on a live agentsview read.

## Stage Report: ideation (cycle 2)

Revision against the independent staff review's two material defects + polish.

- DONE: M1 — the #318 scoping FORMULA was broken (plain `git rev-parse --git-common-dir` returns a RELATIVE path from root/subdir, so `basename(dirname(...))` regresses the exact subdir case #318 targets). Re-grounded the corrected `git rev-parse --path-format=absolute --git-common-dir` across root / `docs/dev` / `.spacedock-state` / external-worktree / non-git this session (git 2.39.5); fixed the formula in ## In scope item 2, ## Spike results scoping section (now shows broken-vs-fixed per-case), and AC-2; stated the git ≥2.31 floor + the `basename(pwd)` fallback.
  Evidence: spike output — plain form gives `.`/`..` from root/subdir (WRONG); `--path-format=absolute` gives `spacedock-v1` from all four locations; non-git exits 128 → fallback.
- DONE: M2 — the proof-harness coupling for AC-2/AC-4/AC-5 was unspecified (`buildFixtureDB` ingests static SQL verbatim, but the union/bucketing match the git-init helper's RUNTIME `t.TempDir()` prefix). Specified the mechanism in ## Foundation, AC-1, AC-2, AC-4, AC-5, and ## Test plan: static SQL carries `__REPO_ROOT__`/`__SUBDIR__`/`__WORKTREE__`/`__EXTERNAL__` marker tokens in `sessions.cwd` + Edit/Write `input_json.$.file_path`; a `coupleFixtureToRepo(t, db, repo)` helper rewrites them to the helper's real dirs via path-scoped `UPDATE` after the static build, so the scan sees real on-disk paths.
  Evidence: AC-1 now has a coupling Verified-by (coupled cwd equals the helper's subdir path, RED if un-substituted); AC-5 shares one source of truth between the fixture OPEN header and the helper's commit subject.
- DONE: Polish — (a) `git_branch`/`category` flagged seed-only/YAGNI (no AC queries them) in AC-1 + ## Foundation; (b) AC-3 recovered-`superpowers` and excluded-`spacedock` cases made DISTINCT rows for diagnosable family-table failures; (c) git-init helper specified to `t.Skip` (not `t.Fatal`) when `git worktree add` is unavailable, matching `buildFixtureDB`.

### Summary

The corrected scoping formula is `PROJECT = sanitize(basename(dirname(git rev-parse --path-format=absolute --git-common-dir)))` with a `basename(pwd)` fallback for non-git / git <2.31 — re-grounded across all four location cases before writing it, so the body no longer ships a formula that regresses the subdir case #318 exists to fix. The M2 proof-harness coupling is now a named, realizable mechanism (marker tokens in static SQL + a post-build `coupleFixtureToRepo` UPDATE keyed to the git-init helper's runtime dirs), so AC-2/AC-4/AC-5's "real on-disk path, never a fixture boolean" RED-levers can actually fire. Fork settlements, taxonomy, proof philosophy, and the rest of the spike are unchanged — confirmed sound by the review.

## Stage Report: ideation (cycle 3)

Captain redirect (2026-06-07): two design changes — no bundled bin/ (supersedes #309),
lighter onboarding proof bar. Substance of the four fixes kept; structure + proof bar
changed.

- DONE: Change 1 — NO bundled bin/ scripts (supersedes #309). Added ## Structure: the bundled `bin/scan-project` + `bin/detect-scaffold` are DELETED; their SQL moves to a recommended-SQL REFERENCE FILE `skills/survey/references/queries.sql` (one labeled query per concern), their scoping / file-probe / classification / report-assembly move to SKILL.md PROSE; SQL is NEVER fenced in SKILL.md (would re-trip the 4q ban). Rewrote Fork 1 (no-bin/ contract), Fork 2 (cross-check is prose; ExitPlanMode fix is one SQL line in the reference), Fork 3 (shares the reference + fixture).
  Evidence: ## Structure section + Fork 1/2/3 rewrites; AC-1 asserts the bin/ paths are gone + the reference file is read.
- DONE: Change 2 — LIGHTER testing (onboarding = low blast radius; efficacy monitored by a separate project). DROPPED the git-init helper, coupleFixtureToRepo / marker-token coupling, per-AC RED-lever rigor, and the coupling-only DDL. New proof = a query-smoke test (recommended queries run against a committed production-shaped fixture, assert corrected shape) + a live drive of `/spacedock:survey`. Rewrote ## Foundation, ## Acceptance criteria (5 heavy ACs → 3 light: AC-1 structure, AC-2 query-smoke, AC-3 live drive), ## Test plan.
  Evidence: ## Acceptance criteria + ## Test plan now specify ONE query-smoke test + one fixture-DDL extension (cwd + skill_name only) + a live drive; heavy machinery explicitly listed as DROPPED.
- DONE: Update existing tests + contractlint. Test plan re-points `survey_extraction_test.go` + `survey_scaffold_test.go` off the deleted `exec.Command(bin/...)` to the query-smoke surface, removing all deleted-script references. Checked `internal/contractlint/structural_checks_test.go` this session: discovery check unaffected; `referenceRe` matches only `.md` (so a `.sql` reference is not enforced by the closure check — the smoke catches a missing reference file); `isClaudeAdapter` already exempts survey/SKILL.md from the HOME-rooted check, and `interpreterRe` forbids only `python`/`commission/bin` — NOT `git`/`sqlite3`/`agentsview` — so moving shell orchestration into SKILL.md prose does NOT trip contractlint. No contractlint change required; confirm green.
  Evidence: ## Test plan "Re-point the existing integration tests" + "Contractlint check" bullets, grounded against the regex definitions (`interpreterRe = python3?|commission/bin`; `referenceRe` requires `\.md`).

### Summary

The four signal-corrections are intact (#318 git-root scoping via `--path-format=absolute --git-common-dir`; #319 multi-label + behavioral scaffold detection with `spacedock:*` excluded; #317.2 WORK-BY-AREA; #320 OPEN-vs-repo + the ExitPlanMode done-prefix fix), now expressed as one annotated `references/queries.sql` (one labeled query per concern) + SKILL.md prose orchestration, with NO bundled bin/ — this supersedes #309's bin/-script resolution. The proof bar drops to onboarding-appropriate light: a query-smoke that EXECUTES the `.sql` reference against a production-shaped fixture (pins SQL-extraction correctness, preserves the 4q "tests run artifacts, not SKILL.md prose" property) plus a live `/spacedock:survey` drive for orchestration; the heavy git-init/coupling/RED-lever machinery is dropped because survey is low-blast-radius onboarding whose efficacy is monitored elsewhere. Verified against the contractlint regexes that the prose-shell restructure (git/sqlite3/agentsview) does not trip any structural check, so no contractlint change is required.

## Stage Report: implementation

Worktree branch `spacedock-ensign/survey-signal-correction`, commit 19d03da9 (8 files, +598/-468).

- DONE: Delete skills/survey/bin/scan-project + bin/detect-scaffold; move their SQL into skills/survey/references/queries.sql (one labeled query per concern, SQL NEVER fenced in SKILL.md) and their scoping/file-probe/classification/report-assembly into SKILL.md prose; apply the four fixes (#318 git-root scoping; #319 multi-label + skill_name family tally with spacedock:* excluded; #317.2 WORK-BY-AREA edit-path bucketing; #320 OPEN-vs-repo confident-match drop + unverified degrade + ExitPlanMode done-prefix).
  Both bin scripts `git rm`'d; `references/queries.sql` carries the 4 labeled queries (scoping / scaffold-usage / work-by-area / decision-open); SKILL.md §2 resolves the repo root via `git rev-parse --path-format=absolute --git-common-dir` (pwd fallback) + runs the queries, §3 is the multi-label probe + 3-bucket join, §4 is the OPEN cross-check + WORK-BY-AREA section. `grep` confirms 0 fenced SQL and 0 forbidden interpreter tokens in SKILL.md.
- DONE: TDD the LIGHT query-smoke: extend fixture DDL with cwd + skill_name, add rows, write the Go query-smoke RED first then green; re-point survey_extraction_test.go + survey_scaffold_test.go off the deleted bin/ exec.Command.
  Fixture extended with `cwd`/`git_branch` (sessions) + `category`/`skill_name` (tool_calls); the smoke (`survey_queries_test.go`, renamed from survey_extraction_test.go) sqlite3-runs each labeled query from `references/queries.sql` and asserts the corrected shape — 5/5 pass. RED proven by reverting each fix in turn (drop spacedock exclusion → spacedock:4 leaks; revert scoping to project='proj' → 1 session not 3; drop ExitPlanMode prefix → approved plan reads OPEN), all RED, restored green. survey_scaffold_test.go deleted (probe is now prose, live-drive proven); no `exec.Command(bin/...)` references remain.
- DONE: Confirm internal/contractlint stays green; go test ./... + gofmt clean; commit on the worktree branch; leave AC-3 live drive for validation.
  `go test ./...` 1140/1140 across 16 packages; `internal/contractlint` 20/20; `gofmt -l` clean repo-wide; `go vet` clean. Committed 19d03da9. AC-3 live drive NOT run (validation's job).

### Summary

Shipped the four signal-corrections as `skills/survey/references/queries.sql` (4 labeled queries) + SKILL.md prose orchestration with NO bundled bin/ (supersedes #309). The query-smoke executes the `.sql` artifact against a production-shaped fixture (extended with `cwd` + `skill_name`) and was proven to have teeth by reverting each fix to RED before greening. Grounded two mechanisms beyond the smoke: the SKILL.md `run_query` awk helper runs the labeled queries against the fixture to the same corrected shapes, and `git rev-parse --path-format=absolute --git-common-dir` resolves to the SAME repo root from the worktree, a subdir, main, and the split-root state checkout (the #318 coalesce), falling back to pwd for non-git. Full suite + contractlint + gofmt green. One observation for validation: `skills/integration/testdata/survey/scaffolds/` is now orphaned (no test consumes it after the bin-detector test's removal) — it documents the multi-label probe cases the AC-3 live drive exercises, so I left it rather than silently deleting; the captain/validator can decide whether to prune it.

## Stage Report: validation

Validated commit `19d03da9` on worktree branch `spacedock-ensign/survey-signal-correction` against the lighter onboarding bar (NORMAL validation — no detached adversarial audit; survey is low-blast-radius onboarding whose report efficacy is monitored by a separate project, so it is NOT one of the four high-stakes surfaces).

- DONE: Reproduce the query-smoke (skills/integration/survey_queries_test.go) — confirm it passes AND has teeth.
  `go test ./skills/integration/ -run TestSurveyQuerySmoke -v` → all 4 subtests PASS (scoping, scaffold-usage, work-by-area, decision-open). TEETH proven by reverting all THREE claimed levers in `references/queries.sql` and re-running: (1) drop `WHERE family <> 'spacedock'` → `spacedock:4` leaks into the tally, RED (`survey_queries_test.go:169,175`); (2) revert prefix union to `project = 'proj'` → `sessions=1` not 3, RED (`:144,147`); (3) drop the `User has approved your plan%` done-prefix → approved plan reads `PLAN:OPEN`, RED (`:226`). Restored to byte-identical (git diff clean), green again.
- DONE: Confirm the structure — bin/scan-project + bin/detect-scaffold gone, references/queries.sql carries the 4 labeled queries, NO fenced SQL in SKILL.md, no exec.Command(bin/...) remains; go test ./... + internal/contractlint + gofmt all green.
  `skills/survey/bin/` does not exist; `references/queries.sql` has 4 `-- name:` blocks (scoping/scaffold-usage/work-by-area/decision-open); `grep '```sql' SKILL.md` → 0; no `scan-project`/`detect-scaffold`/`exec.Command(bin/...)` refs anywhere in skills/integration (the only `exec.Command` calls are to `sqlite3`). `go test ./...` 1140/1140 across 16 packages; `internal/contractlint` 20/20; `gofmt -l` clean; `go vet` no issues.
- DONE: AC cross-check — AC-1, AC-2 verified from outside the body; AC-3 flagged as captain-supplied live evidence pending at the gate.
  AC-1 (no bundled bin/; SQL-reference + SKILL.md prose; structure): verified above — the deleted bin/ paths are gone, the smoke RED-on-missing-reference-file lever is real (the test reads `references/queries.sql` at `survey_queries_test.go:57`, t.Fatal if absent), no SQL fenced in SKILL.md. AC-2 (labeled queries return corrected data on the fixture): verified — the smoke executes the `.sql` artifact (NOT SKILL.md prose) and asserts FIXTURE-derived shapes (3 coalesced sessions, superpowers=3 / spacedock excluded, internal=2 / skills=1 / <external>=1, ExitPlanMode done + rejection OPEN); RED-on-revert confirms the expected values diverge from a broken query. AC-3 (live drive): NOT attempted by me — the spike proved a sandboxed process is TCC-blocked from the agentsview sessions.db. Instead confirmed SKILL.md prose orchestration is internally consistent with `references/queries.sql`: the 4 `run_query <name>` calls in §2 match the 4 labeled queries exactly; the §2 awk extractor correctly pulls the `scoping` block; and the live `git rev-parse --path-format=absolute --git-common-dir` mechanism resolves to the SAME repo root from the worktree root, a worktree subdir, main, and a main subdir (the #318 coalesce). AC-3 remains captain-supplied live evidence pending at the gate.
- DONE: Decide the orphaned skills/integration/testdata/survey/scaffolds/ — recommend prune or keep-with-reason.
  RECOMMEND PRUNE. Confirmed no test consumes it (`grep -rn scaffolds skills/integration/` → 0 after `survey_scaffold_test.go` was deleted). The multi-label file probe is now SKILL.md prose exercised by the AC-3 live drive, which runs against REAL repos, not committed fixtures — so the scaffolds dir has no automated consumer and the live drive does not read it. Non-blocking either way (dead fixture, not a correctness issue); flagging for the captain to drop in a cleanup commit. Not a gate condition.

### Honesty note

The query-smoke's expected values come from the FIXTURE rows (`testdata/survey/fixture-sessions.sql`), an independent source — NOT from SKILL.md prose. The test loads and EXECUTES `references/queries.sql` via sqlite3 (`survey_queries_test.go:55-87, 91-111`); it never parses SKILL.md. This preserves the 4q "tests run artifacts, not instruction text" property. The RED-on-revert sweep is the proof the asserted values can actually diverge from a broken query — a static spelling-check could not produce `spacedock:4 leaks` / `sessions=1` / `PLAN:OPEN`.

### Recommendation

**PASSED** against the lighter onboarding bar. AC-1 and AC-2 are verified from outside the entity body with a teeth-proven smoke (3/3 levers RED-on-revert); structure + full hygiene (1140 tests + contractlint + gofmt + vet) are green; the SKILL.md prose orchestration is internally consistent with the SQL reference and the #318 git-root mechanism works live. AC-3 (the full `/spacedock:survey` live drive) is the captain's manual re-test from other projects via `--plugin-dir` to this worktree — pending captain-supplied live evidence at the gate, per the redirect. One non-blocking cleanup recommendation: prune the orphaned `testdata/survey/scaffolds/` dir.

### Summary

Reproduced the AC-2 query-smoke (4/4 subtests pass) and proved its teeth by reverting all three claimed fixes in `references/queries.sql` to RED in turn, then restoring byte-identical. Verified AC-1 structure (no bin/, 4 labeled queries, no fenced SQL, no stale exec.Command) and full hygiene (1140 tests, contractlint 20/20, gofmt + vet clean) from outside the body. AC-3's full live drive is sandbox/TCC-blocked, so confirmed instead that the SKILL.md prose is internally consistent with the SQL reference (run_query names ↔ labeled queries, awk extraction works) and that the #318 git-root coalesce resolves identically from four checkout locations live — flagging the full live drive as captain-supplied evidence pending at the gate. Recommend PRUNE for the orphaned scaffolds fixture. Verdict: PASSED.
