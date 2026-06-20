---
title: Pi dev-override package-OK — honor --plugin-dir/SPACEDOCK_REPO_ROOT in spacedockPackageOK (regression fix for eq D5b)
status: implementation
source: "Commander (2026-06-20): pi-live lane on qn PR #407 (TestLivePiFrontDoorSmoke) failed — eq's D5b spacedockPackageOK only checks settings.json packages, not the --plugin-dir dev override, so a dev-override launch bails 'Pi runtime is not ready'. Fix on main (post-eq-merge); unblocks qn #407."
score:
started: 2026-06-20T01:30:00Z
completed: 2026-06-20T01:34:15Z
verdict: PASSED
worktree:
issue:
sprint: 0223-pi-dispatch-contract
sprint-readiness: ready
id: 
mod-block:
pr: "#408"
archived: 2026-06-20T01:34:15Z
---

# Pi dev-override package-OK

Regression fix: `checkPiRuntime`'s `spacedockPackageOK` (eq D5b) only reads `settings.json packages`, ignoring `cfg.repoRoot` (`--plugin-dir`/`SPACEDOCK_REPO_ROOT` dev override). A dev-override launch against a fresh pi-home bails before the ensign. Fix: honor a valid dev-override checkout (`skills/ensign/SKILL.md` present) as satisfying the gate. Unblocks qn #407's pi-live lane.

## Acceptance criteria

- AC-1: `checkPiRuntime` with `cfg.repoRoot` set to a checkout containing `skills/ensign/SKILL.md` + no `settings.json` package yields `spacedockPackageOK = true` + `piRuntimeLaunchReady = true` (Go test).
- AC-2: `repoRoot` empty + no package → `spacedockPackageOK = false` (unchanged; Go test).
- AC-3: `TestPiRuntimeConfigRetiresSkillFlagsAndCwdFallback` still PASS (no regression).
- AC-4: `pi-live` lane `TestLivePiFrontDoorSmoke` green in CI (the live proof; runs in CI, not worktree).

## Stage Report: implementation

(pending ensign)
