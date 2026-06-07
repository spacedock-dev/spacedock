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
PROXIES while the ground truth sits unread right next to it. One root cause, the
same proof-policy class we already fence (trust a proxy over ground truth),
manifest across the skill's two artifacts (`skills/survey/bin/detect-scaffold`,
`skills/survey/bin/scan-project`). Verified against live code/data this session:

- `detect-scaffold` prints `none` on THIS spacedock repo — the skill can't
  recognize its own home (file-only probe of `.claude/skills` etc.).
- `tool_calls.skill_name` is populated (438 rows incl. `spacedock:ensign`×117)
  and never queried — the behavioral truth of which scaffolds ran is ignored.
- `basename(cwd)` scoping drops ~67% of sessions on a real multi-worktree repo.
- The OPEN-decision "NEEDS YOU" frontier is transcript-only — on a real run every
  OPEN fork was already merged (all false positives).

This is the highest-value survey work because it makes the skill CORRECT for its
ACTUAL users (greenfield / other-scaffold repos), not the spacedock-incumbent
minority.

## In scope (one cohesive deliverable)

1. **Scaffold detection — multi-label + behavioral** (#319, #317.1). Report ALL
   scaffolds, not the first if-ladder match; corroborate the file probe against a
   `tool_calls.skill_name` tally (file+invoked / file+never-invoked=installed-but-unused /
   not-file+invoked=recovered). Exclude survey's own `spacedock:survey`/`ensign`
   self-invocation from the tally (self-pollution).
2. **Identity scoping — coalesce a repo across keys** (#318). Scope by repo identity
   via `git rev-parse --git-common-dir` (NOT `--show-toplevel`, which returns the
   worktree root) / path-prefix boundary / `worktree_project_mappings` override;
   union all sessions under the repo; surface folded-in keys + blank-`cwd` count;
   deterministic (no fuzzy name matching). No-regression for single-checkout repos.
3. **Identity from edits** (#317.2). A WORK-BY-AREA tally bucketing `Edit/Write`
   paths by package; report "what this is" (edits) separately from "where you stop"
   (decisions); treat edits to external sibling repos as references, not identity;
   show all-time alongside recent.
4. **OPEN frontier — cross-check against the repo** (#320). After the transcript
   scan, cross-reference each OPEN fork against repo artifacts (git log, merged PRs,
   working tree) and split shipped (drop) / decided-not-shipped (backlog) /
   never-decided (true open). Conservative default (drop ONLY on a confident match,
   else keep on the frontier). Mandatory transcript-only degrade, flagged
   `unverified`, when no repo signal. Fold in the cheap independent fix: the
   `ExitPlanMode "User has approved your plan"` prefix matches none of the three
   done-prefixes today (`scan-project`), so approved plans fall to OPEN.

## Foundation (prerequisite — ideation owns the design)

- Extend the test fixture DDL with `cwd` / `git_branch` / `skill_name` — they are
  ABSENT today, so a naive scoping/scaffold fix passes green doing nothing
  (the invisible-fix trap).
- Build ONE shared git-init test helper (three sub-parts need a real git repo/worktree
  fixture; today fixtures are static file trees).
- Seed at least one anonymized fixture from a REAL agentsview `sessions.db` dump so
  Skill/Edit/decision shapes match production (all fixtures today are hand-authored —
  the validate-on-real-data gap; closes it once for the whole cluster).

## Out of scope

- **#316** detect-spacedock-incumbent + the audit/meta-loop "flip the back half"
  reframe — defer; it serves spacedock-incumbent users only. Record the audit reframe
  as a roadmap proposal. (Design the scaffold-detection contract in part 1 to be
  EXTENSIBLE/multi-label so a future spacedock case drops in without a rebase.)
- **#317.3** the low-confidence structured issue-filing prompt — defer (meta-feature).

## Design forks for ideation to settle (then gate to captain)

- The unified `detect-scaffold` contract: where the file-probe and the DB tally
  reconcile (`detect-scaffold` has no DB access today — new `bin/recognize-scaffold`
  vs fold into `scan-project` vs a SKILL.md join), and the `skill_name`→family
  normalization (strip `superpowers:*` prefixes, attribute bare names, exclude self).
- The #320 3-bucket taxonomy + the conservative-match rule + the degrade contract.
- Whether any sub-part is large enough to split (default: keep as one).

## Proof discipline (non-negotiable — proof-policy class)

Every AC reads the ground-truth signal that already exists AND ships a fixture
carrying that signal so a proxy-only impl goes RED. NO SKILL.md substring tests —
run the real `bin/` script via `cmd.Dir` against committed fixtures (the existing
`survey_*_test.go` pattern). RED-levers: scaffold marker at a non-`cwd` path; a
worktree row under a different project key with a unique decision header; DB carries
`superpowers:*` rows while disk has none ("recovered"); "shipped" derived from a real
on-disk git repo, never a fixture-provided boolean.

## Sources (close after this ships)

spacedock-dev/spacedock#318, #319, #320, #317 (umbrella — superseded by #318/#319),
and #316 (deferred). These reports get closed after this task lands; no 1-1 mapping.
