---
id: xnh4gc2zqe67j8tbxh926r9j
title: Survey signal-correction — read the ground-truth signals (sessions.db + repo), not thin proxies
status: ideation
source: "captain (2026-06-07) - combined fix for the survey-skill dogfood issue cluster spacedock-dev/spacedock#318/#319/#320 (+ #317 umbrella). One task, not 1-1: the reports close after this ships. #316 (detect-spacedock incumbent + audit/meta-loop flip) is OUT of scope (serves spacedock-incumbent users, a minority) - park as a separate roadmap proposal."
started: 2026-06-07T18:57:58Z
completed:
verdict:
score:
worktree:
issue:
---

`spacedock:survey` decides "what this project is and where it stands" from THIN
PROXIES while the ground truth sits unread right next to it. One root cause — the
same proof-policy class we already fence (trust a proxy over ground truth) —
manifest across the skill's two artifacts (`skills/survey/bin/detect-scaffold`,
`skills/survey/bin/scan-project`). Verified against live code + a freshly-synced
real `sessions.db` this session (see ## Spike results):

- `detect-scaffold` prints `none` on THIS spacedock repo — the skill can't
  recognize its own home (file-only probe of `.claude/skills` etc.).
- `tool_calls.skill_name` is populated (1045 `Skill`-call rows in the spike sync;
  `spacedock:ensign`×964, `superpowers:*`, bare `running-research-spikes`) and
  never queried — the behavioral truth of which scaffolds ran is ignored.
- `scan-project` scopes by `basename(pwd)`, but agentsview keys every session by
  the RESOLVED repo-root basename. Run from a subdir/worktree the two keys diverge
  and survey finds 0 sessions while 1198 sit under the agentsview key (spike).
- The OPEN-decision "NEEDS YOU" frontier is transcript-only and over-reports
  (forks already merged read as still-open).

This is the highest-value survey work because it makes the skill CORRECT for its
ACTUAL users (greenfield / other-scaffold repos), not the spacedock-incumbent
minority.

## In scope (one cohesive deliverable)

1. **Scaffold detection — multi-label + behavioral** (#319, #317.1). Report ALL
   scaffolds, not the first if-ladder match; corroborate the file probe against a
   `tool_calls.skill_name` tally (file+invoked / file+never-invoked=installed-but-unused /
   not-file+invoked=recovered). Exclude survey's own `spacedock:survey`/`ensign`
   self-invocation from the tally (self-pollution — it is 1026/1045 of the tally
   on this repo, see ## Spike results).
2. **Identity scoping — coalesce a repo across keys** (#318). Scope by repo identity
   via `git rev-parse --git-common-dir` (NOT `--show-toplevel`, which returns the
   worktree/subdir root) / path-prefix boundary / `worktree_project_mappings` override;
   union all sessions under the repo; surface folded-in keys + blank-`cwd` count;
   deterministic (no fuzzy name matching). No-regression for single-checkout repos.
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
   done-prefixes today (`scan-project`), so approved plans fall to OPEN.

## Foundation (prerequisite — ideation owns the design, implementation ships it)

- Extend the test fixture DDL with `cwd` / `git_branch` (on `sessions`) and
  `skill_name` / `category` (on `tool_calls`) — they are ABSENT today, so a naive
  scoping/scaffold fix passes green doing nothing (the invisible-fix trap).
  Defaults match production: `cwd`/`git_branch` default `''` (NOT NULL), blank rows
  store `''` not NULL — so blank-`cwd` handling is exercised, not skipped.
- Build ONE shared git-init test helper (parts 2 and 4 need a real git repo +
  worktree; today fixtures are static file trees). The helper `git init`s a temp
  repo, adds a worktree, and returns paths so a test can `cmd.Dir` into a subdir or
  worktree and exercise the real `git rev-parse --git-common-dir` resolution.
- Seed at least one anonymized fixture from a REAL agentsview `sessions.db` dump so
  Skill/Edit/decision shapes match production (all fixtures today are hand-authored —
  the validate-on-real-data gap; closes it once for the whole cluster). The spike
  already pulled the real shapes (## Spike results); the fixture encodes them.

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
- `git rev-parse --git-common-dir` resolves the repo root from any subdir/worktree;
  `--show-toplevel` does NOT (returns the worktree/subdir as root):
  - from `docs/dev/.spacedock-state`: `--show-toplevel`=`…/.spacedock-state` (TRAP),
    `dirname(--git-common-dir)`=`…/spacedock-v1` (correct);
  - from external worktree `/private/tmp/spacedock-audit-autoinstall-7462`:
    `--show-toplevel`=that worktree (TRAP), `dirname(--git-common-dir)`=`…/spacedock-v1`.
  → `PROJECT = sanitize(basename(dirname(git rev-parse --git-common-dir)))`, falling
  back to `basename(pwd)` when not in a git repo, makes scan-project's key match
  agentsview's key from any subdir/worktree.
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

### Fork 1 — the unified detect-scaffold contract

**Where the DB tally lives.** `detect-scaffold` is file-only and has NO `$DB`.
`scan-project` already opens `$DB` and already owns the repo scope (and gains the
git-root scope from part 2). Settled: **keep two binaries; add the behavioral tally
to `scan-project` as a new `## SCAFFOLD-USAGE` section** (it has the DB + the scoped
session set), keep `detect-scaffold` as the FILE probe but make it **multi-label**
(emit every matched scaffold, one per line, not a single ladder winner). SKILL.md
step 3 joins the two outputs into the 3-bucket classification. Rationale: each binary
keeps a single responsibility; no DB dependency is forced onto the file probe; the
join is cheap and already where the comparative-benefit logic lives.
- Rejected: a new `bin/recognize-scaffold` doing both (duplicates scan-project's DB
  open + scope; two places to keep the scope correct).
- Rejected: folding the file probe into scan-project (couples two independent
  signals; breaks the existing `detect-scaffold` test surface and `cmd.Dir` pattern).

**The 3-bucket scaffold classification** (join of file probe × usage tally):
- `file-present + invoked` → **active** (reported plainly).
- `file-present + never-invoked` → **installed-but-unused** (flagged).
- `not-file-present + invoked` → **recovered** (reported — this is the case the
  file-only probe misses entirely; the behavioral signal recovers it).

**`skill_name` → family normalization** (grounded above): split on first `:`, prefix
is the family; bare names map via a small superpowers-discipline name set → family
`superpowers`; **exclude family `spacedock`** (self-pollution). The family set is a
table, not an if-ladder — adding a `spacedock` family row later (#316) is data, not a
contract change, satisfying the EXTENSIBLE requirement.

**`detect-scaffold` multi-label.** Today's ladder (`if superpowers; elif gsd; elif
similar; else none`) becomes: probe each scaffold independently, print every match
(`superpowers`, `gsd`, `similar: <names>` as appropriate), print `none` only when no
probe matched. The existing per-fixture test (`survey_scaffold_test.go`) still passes
for single-scaffold fixtures; a new two-scaffold fixture asserts BOTH labels emit.

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
`User approved` variants) to the done-prefix set in `scan-project` DECISIONS so an
approved plan reads `done`, not OPEN. Independent, fixture-provable, one-line.

**Where the cross-check runs.** The repo cross-reference needs `git log` / merged-PR
data + the working tree — that is the survey driver's job (SKILL.md step 4 already
narrates the report). Settled: `scan-project` emits the raw transcript-OPEN frontier
(as today); the cross-check + 3-bucket split is a new step the SKILL.md step-4 report
performs against the repo, with the conservative/degrade rules above. The
behavioral, testable half is a small helper the SKILL invokes — a `bin/` cross-check
that takes the OPEN headers + the repo root and returns shipped/backlog/open — so the
classification is exercised by a test, not left as SKILL prose. (Implementation may
fold it into scan-project or a sibling `bin/`; the AC is the behavior, not the file.)

### Fork 3 — split?

**Keep as one.** The four sub-parts share the fixture foundation (the DDL extension,
the git-init helper, the real-data fixture) — splitting would duplicate or serialize
that foundation. The parts are independent enough to implement in sequence within one
task but share enough scaffolding that one entity is correct. Default honored.

## Acceptance criteria (entity-level; each "Verified by" names an outside check)

Proof discipline (non-negotiable — proof-policy class): every AC reads the
ground-truth signal that already exists AND a committed fixture carries that signal so
a proxy-only impl goes RED. NO SKILL.md substring tests — every check runs the real
`bin/` script via `cmd.Dir` against committed fixtures (the existing
`survey_*_test.go` pattern). The expected values come from the FIXTURE rows / the
real on-disk git repo — an independent source that can diverge from the skill text.

**AC-1 — Fixture foundation carries the ground-truth signals.**
The committed fixture DB has `sessions.cwd`, `sessions.git_branch`,
`tool_calls.skill_name`, `tool_calls.category` (defaults matching production: `cwd`/
`git_branch` NOT NULL default `''`), at least one fixture is seeded from anonymized
REAL `sessions.db` shapes (namespaced + bare `skill_name`, `Your questions have been
answered:` / `The user doesn't want to proceed` decision results), and a shared
git-init test helper materializes a real temp git repo + worktree.
- Verified by: the fixture DB builds (`buildFixtureDB` succeeds) and a test asserts
  the four columns exist and the git-init helper returns a path where
  `git rev-parse --git-common-dir` resolves — RED if any column or the helper is
  missing. (Foundation AC: it is the prerequisite the other ACs' RED-levers stand on.)

**AC-2 — Identity scoping coalesces a repo across keys (#318).**
Run against a fixture DB whose sessions span the repo root cwd, a subdir cwd, and a
worktree cwd (distinct agentsview `project` keys), with the scan invoked from a
SUBDIR/worktree via the git-init helper, survey reports ONE identity whose session
count / date range / decision + interruption totals equal the git-root-coalesced
totals — not the single basename-key subset; it lists the folded-in keys and a
blank-`cwd` count; a populated `worktree_project_mappings` row overrides path-prefix
inference; a single-checkout fixture is unchanged (no regression); coalescing is
deterministic (git-root / prefix / mapping only).
- Verified by: an integration test that `cmd.Dir`s into the helper's subdir/worktree
  and asserts the coalesced count. RED-lever: a worktree-cwd row under a DIFFERENT
  project key carrying a UNIQUE decision header — a basename-only impl drops it and
  the header is absent from output; the git-root impl unions it in. The expected
  count comes from the fixture rows, not the skill text.

**AC-3 — Scaffold detection is multi-label + behavioral (#319, #317.1).**
`detect-scaffold` emits ALL matched scaffolds (not the first ladder winner); the
`scan-project` `## SCAFFOLD-USAGE` tally classifies each scaffold as active /
installed-but-unused / recovered from the join of file-probe × `skill_name` usage,
with `family:skill` + bare-name normalization and `spacedock:*` self-invocation
EXCLUDED.
- Verified by: (a) a two-scaffold fixture repo → `detect-scaffold` test asserts BOTH
  labels (extends `survey_scaffold_test.go`); (b) a fixture DB whose `tool_calls` has
  `superpowers:*` rows but whose disk tree has NO superpowers files → SCAFFOLD-USAGE
  reports superpowers as **recovered** (RED-lever: a file-only impl reports `none`);
  (c) a fixture DB dominated by `spacedock:ensign` rows → the tally does NOT report
  spacedock (self-pollution excluded; RED-lever: a no-exclusion impl reports
  spacedock). Expected labels come from the fixture file tree + DB rows.

**AC-4 — Identity from edits: WORK-BY-AREA (#317.2).**
`scan-project` emits a `## WORK-BY-AREA` section bucketing `Edit`/`Write`
`input_json.$.file_path` by package under the repo root, all-time alongside recent,
reported separately from decisions; paths OUTSIDE the repo root are bucketed as
external references, not repo identity.
- Verified by: a fixture DB with Edit/Write rows whose `file_path` spans two repo
  packages + one external-sibling path → the test asserts both internal packages
  appear with counts and the external path is segregated as a reference. RED-lever:
  the external path's package name must NOT appear among the repo-identity buckets.
  Expected buckets come from the fixture rows.

**AC-5 — OPEN frontier is cross-checked against the repo (#320).**
The OPEN frontier is split shipped (drop) / decided-not-shipped (backlog) /
never-decided (true open) by cross-referencing each OPEN fork against the on-disk git
repo (merged PRs / git log / working tree); DROP happens ONLY on a confident match;
with no repo signal the frontier degrades to transcript-only with every OPEN flagged
`unverified`; and the `ExitPlanMode "User has approved your plan"` result reads
`done`, not OPEN.
- Verified by: (a) an integration test using the git-init helper to build a real
  repo whose git log confidently references one OPEN fork's header → that fork is
  classified shipped (dropped) while an unreferenced OPEN fork stays on the frontier
  (RED-lever: "shipped" derived from the REAL on-disk git log, never a
  fixture-provided boolean — a transcript-only impl keeps the merged fork on the
  frontier); (b) a no-git-repo case → the frontier is emitted transcript-only and
  flagged `unverified`; (c) a fixture `ExitPlanMode` row with
  `User has approved your plan` → DECISIONS marks it `done` (RED-lever: today's
  three-prefix impl marks it OPEN). Expected classifications come from the real git
  repo + fixture rows.

## Test plan

- **Surface:** Go integration tests under `skills/integration/`, extending the
  existing `survey_extraction_test.go` / `survey_scaffold_test.go` pattern — run the
  real `bin/scan-project` and `bin/detect-scaffold` via `exec.Command` + `cmd.Dir`
  against committed fixtures. NO new test framework. `sqlite3` + `git` + `bash` are
  the executors; tests `t.Skip` (not fail) when a tool is absent so a minimal box
  stays runnable without a false pass (matches `buildFixtureDB`).
- **Fixtures (the foundation, AC-1):**
  - extend `testdata/survey/fixture-sessions.sql` DDL with the four columns +
    defaults, and add rows exercising the new scope (a subdir-cwd session, a
    worktree-cwd session under a different project key, `Skill` rows with namespaced +
    bare + `spacedock:*` self skill_names, Edit/Write rows with internal + external
    file_paths, an `ExitPlanMode` approved row);
  - one anonymized real-shape fixture (counts kept, identifiers stripped) so the
    Skill/Edit/decision shapes match production;
  - a shared `gitInitHelper(t)` that `git init`s a temp repo, adds a worktree, and
    returns repo-root / subdir / worktree paths.
- **No-fixture-boolean rule:** "shipped" (AC-5) and "coalesced from a worktree"
  (AC-2) are derived from a REAL on-disk git repo built by the helper, never a
  fixture-provided flag — that is what makes the proxy-only impl RED.
- **Cost/complexity:** medium. Fixture DDL extension + git-init helper are the bulk
  of the new scaffolding (shared across AC-2/AC-4/AC-5). The four behavioral ACs are
  ~1 focused test each. No live workflow / live agentsview dependency — the spike
  proved the live read is sandbox-blocked, so the committed-fixture path is the only
  CI-safe surface and it is sufficient.
- **No spike outstanding:** the two named risky unknowns (skill_name shape; cwd/git
  scoping against real data) are grounded in ## Spike results; the queries there seed
  the implementation's first tests. The remaining mechanisms (sqlite3 fixture build,
  `cmd.Dir` exec, `git rev-parse --git-common-dir`) are already proven by the
  existing tests + the spike.

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
