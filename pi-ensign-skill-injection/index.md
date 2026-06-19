---
title: Pi ensign skill injection — make the Spacedock ensign skill discoverable by pi-subagents
status: backlog
source: "Captain (2026-06-19): dispatched pi-back-channel-dispatch ideation with skill:[\"ensign\"]; run meta carried \"skillsWarning: Skills not found: ensign\". The worker ran skill-less and defaulted to implementation behavior in an ideation (gate, read-only) stage. Root-cause the injection failure and fix it so dispatched Pi ensigns actually load the Spacedock ensign contract."
score:
started:
completed:
verdict:
worktree:
issue:
sprint:
sprint-readiness:
id: k8tbnmcbyqc5kkhj0m9vewq4
---

# Pi ensign skill injection

## End value

A dispatched Pi worker invoked with `skill:["ensign"]` (or the FO's default ensign-dispatch) actually loads the Spacedock ensign contract (`skills/ensign/SKILL.md` + its references) into its system prompt, so stage dispatches enforce the ensign's stage-output discipline (ideation = design in the entity body, not product edits; implementation = worktree deliverable; etc.) instead of falling back to the builtin `worker`'s implementation-biased prompt.

## Problem — root cause already determined

pi-subagents resolves injectable skills via its **own** discovery (`discoverAvailableSkills(cwd)` → `buildSkillPaths` in `node_modules/pi-subagents/src/agents/skills.ts:318`). The search paths are, in priority order:

1. `{cwd}/.pi/skills/` (project)
2. `{cwd}/.agents/skills/` (project)
3. `{agentDir}/skills/` (user — `~/.pi/agent/skills/`)
4. `{homedir}/.agents/skills/` (user)
5. installed package skill roots (`package.json -> pi.skills`)
6. settings-package skill roots
7. `{cwd}` package root (`package.json -> pi.skills`)
8. settings.json skill entries

The Spacedock ensign skill lives at `skills/ensign/SKILL.md` — **not** in any of these. The main pi agent discovers it (it appears in the FO session's `<available_skills>`) via the project's plugin/manifest skill discovery, but **pi-subagents does not use that discovery** — it uses the path list above. `subagents-doctor` confirms: `skills: total 2 (user 1, user-package 1)` — ensign is absent.

Consequence: `skill: ["ensign"]` on a `subagent(...)` call emits `Warning: skills not found: ensign` and the child runs as a bare `worker`. The builtin `worker` prompt is implementation-biased ("make the edits," "do not return a success summary without making edits"), so a skill-less worker in an ideation (gate, read-only, no-worktree) stage silently behaves as an implementer and risks contaminating `main` product files. In a worktree implementation stage the same failure would be invisible (it would look like normal work). **The ensign skill injection is the gate that enforces the stage-output contract; without it, the contract is unenforced.**

## Evidence (live, 2026-06-19)

- Run `0637e2ed` (pi-back-channel-dispatch ideation dispatch): meta at `~/.pi/agent/sessions/.../subagent-artifacts/0637e2ed_worker_0_meta.json` carries `"skillsWarning": "Skills not found: ensign"`.
- `subagents-doctor` output: `skills: total 2 (user 1, user-package 1)` — ensign not listed.
- The worker then made contract-doc edits during ideation (contained — reverted, nothing landed on `main`; only a state-checkout stage-report commit `0a1e2787`).

## Approach (candidate fixes — ideation confirms and picks)

Make the ensign skill discoverable by pi-subagents. Candidate mechanisms, all project-scoped (travel with the clone, no per-machine setup):

- **(a) Symlink**: `.pi/skills/ensign -> ../../skills/ensign`. Simplest; the skill source of truth stays at `skills/ensign/`; `.pi/skills/ensign/SKILL.md` resolves through the symlink. Source classified `project` (highest priority).
- **(b) `.pi/settings.json` skills entry**: declare the path explicitly. More config, but no symlink (some platforms/checkouts handle symlinks poorly).
- **(c) `package.json -> pi.skills` array**: declare `["skills/ensign/SKILL.md"]` (or the dir). Works only if a `package.json` is present at the project root (there is one — `package-lock.json` is untracked in `git status`, so confirm).
- **(d) Move the skill**: relocate `skills/ensign/` to `.pi/skills/ensign/`. Rejected — breaks the existing main-agent discovery and the skill's relationship to `skills/` siblings.

Ideation picks one (recommend (a) symlink — lowest drift, source of truth unchanged), records the decision, and plans verification. Note: the first-officer runtime adapter (`skills/first-officer/references/pi-first-officer-runtime.md`) says "The child must load the Spacedock ensign skill and Pi ensign runtime adapter before working" — the fix makes that instruction actually achievable on pi-subagents; update the adapter if the mechanism needs documenting.

## Acceptance criteria (provisional — ideation finalizes; proof = behavior, never prose-grep)

**AC-1 — `subagents-doctor` lists the ensign skill as discoverable.**
Verified by: `subagents-doctor` output showing ensign in the skills count/list (a live run, not a static claim).

**AC-2 — A `subagent(... skill:["ensign"])` dispatch loads the ensign contract with no "skills not found" warning.**
Verified by: a probe dispatch whose run meta carries no `skillsWarning` AND whose child exhibits ensign-contract behavior (e.g., in an ideation probe, produces a design in the entity body rather than editing product files — the inverse of the `0637e2ed` failure).

**AC-3 — The discovery mechanism is project-scoped and travels with the clone.**
Verified by: the fix lives in the repo (symlink or committed config), not in `~/.pi/agent/` user config; a fresh clone on another machine discovers ensign with no manual setup.

## Out of scope

- The model-stamping friction (null model → settings default, not parent model) — separate concern, tracked alongside `pi-back-channel-dispatch`.
- Changing pi-subagents' search paths upstream (the fix is project-side discovery wiring, not a pi-subagents change).
- The ensign skill's *content* — this task is about *injectability*, not the contract text.

## Test plan

- `subagents-doctor` (AC-1) — live, cheap.
- A probe `subagent(... skill:["ensign"])` dispatch (AC-2) — run meta has no `skillsWarning`; child behavior matches ensign contract. Bounded probe, not a full stage drive.
- Fresh-clone or `git stash`-and-verify reasoning for AC-3 (the fix is a committed repo artifact; confirm it is version-controlled).

## Related

- `pi-back-channel-dispatch` (`b23y61pgk93ph44pz506m2wy`) — sibling Pi-dispatch-friction task; this one unblocks correct ensign dispatch for that and every subsequent Pi stage dispatch.
- Run `0637e2ed` — the failure instance.
- `node_modules/pi-subagents/src/agents/skills.ts:318` (`buildSkillPaths`) — the search-path source of truth.
