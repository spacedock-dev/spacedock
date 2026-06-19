---
title: Pi ensign skill injection — make the Spacedock ensign skill discoverable by pi-subagents
status: ideation
source: "Captain (2026-06-19): dispatched pi-back-channel-dispatch ideation with skill:[\"ensign\"]; run meta carried \"skillsWarning: Skills not found: ensign\". The worker ran skill-less and defaulted to implementation behavior in an ideation (gate, read-only) stage. Root-cause the injection failure and fix it so dispatched Pi ensigns actually load the Spacedock ensign contract."
score:
started: 2026-06-19T18:00:51Z
completed:
verdict:
worktree:
issue:
sprint: 0223-pi-dispatch-contract
sprint-readiness: ready
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

## Problem — the parent/child discovery split (reframed)

The captain's constraint (load-bearing): the skill should NOT be discoverable merely because pi was launched from this repo's cwd. The bug is a **parent/child discovery split**, not a cwd-accident on the parent side:

- **Parent pi session**: `spacedock pi` passes `--skill <repo>/skills/ensign/SKILL.md` explicitly to the `pi` binary (`internal/cli/pi.go:87-89`). The parent finds ensign because the launcher registers it — explicit, launch-managed.
- **pi-subagents children**: use their OWN discovery (`discoverAvailableSkills(cwd)` → `buildSkillPaths` in `node_modules/pi-subagents/src/agents/skills.ts:318`) and do NOT inherit the parent's `--skill` flags. The search paths are `.pi/skills/`, `.agents/skills/`, `~/.pi/agent/skills/`, package roots, settings — none of which contain the Spacedock ensign skill at `skills/ensign/SKILL.md`.

So the parent sees ensign (via `--skill`), the child does not (via filesystem discovery). The child-side fix is a **project-declared skill** — a committed repo artifact under `.pi/skills/` that pi-subagents' `buildSkillPaths` discovers. This is distinct from the parent's `--skill` flag and distinct from the parent-side cwd-fallback fix (sibling task `pi-launcher-repo-resolution`).

Note: `{cwd}/.pi/skills` is cwd-keyed, so it works only when the child's cwd is the repo. This is acceptable BECAUSE it is a committed project artifact (not cwd-luck): the FO dispatches with cwd = project root, so the child inherits the repo as its cwd, and the committed `.pi/skills/ensign` symlink is found. The project declares the skill; the child discovers it because the project declared it.

## Approach — mechanism decision: (a) project-declared symlink

### Spike evidence (live, 2026-06-19)

**Riskiest mechanism**: does a `.pi/skills/ensign -> ../../skills/ensign` symlink actually make `discoverAvailableSkills(cwd)` list ensign? **SPIKED — PASSED.**

Pre-fix (no `.pi/` dir): ran the real `discoverAvailableSkills(cwd)` from `pi-subagents/src/agents/skills.ts` via `npx tsx`:
```
total: 2
  find-skills -> (user)
  pi-intercom -> (user-package)
```
Ensign ABSENT — matches the `0637e2ed` failure.

Created `.pi/skills/ensign -> ../../skills/ensign` (throwaway spike). Post-fix:
```
total: 3
  ensign -> (project)
  find-skills -> (user)
  pi-intercom -> (user-package)
ensign found: YES
  source: project
```
Ensign PRESENT, classified `project` (highest priority). Symlink removed after — implementation creates it properly.

### Mechanism decision: (a) `.pi/skills/ensign -> ../../skills/ensign` symlink

The symlink is a **project-declared skill for child discovery** — committed, travels with the clone, source of truth stays at `skills/ensign/`. pi-subagents' `buildSkillPaths` resolves it through the symlink and classifies it `project` (highest priority).

**Why over the alternatives:**
- **(b) `.pi/settings.json` skills entry**: extra config file for no benefit — the symlink achieves the same discovery with zero config drift. Symlink-phobia (some platforms handle symlinks poorly) is not a real concern here: macOS/Linux dev environments, and the symlink target is within the same repo. Rejected.
- **(c) `package.json -> pi.skills`**: CONFIRMED no root `package.json` exists today (`ls package.json` → not found). Would require adding a `package.json` solely for this purpose — more surface area, more drift. Rejected.
- **(d) Move the skill** to `.pi/skills/ensign/`: breaks the existing main-agent discovery (the launcher's `--skill <repo>/skills/ensign/SKILL.md` path) and the skill's relationship to `skills/` siblings. Rejected.

### What implementation ships

1. `.pi/skills/ensign -> ../../skills/ensign` (symlink, committed to the repo).
2. If the first-officer runtime adapter (`skills/first-officer/references/pi-first-officer-runtime.md`) needs documenting that the child loads the ensign skill via this project-declared path, propose the doc edit in the implementation stage report (do not edit in ideation).
3. Consider whether `first-officer` and other Spacedock skills (`commission`, `debrief`, `refit`) should also be symlinked under `.pi/skills/` — they are NOT injected by the FO dispatch today (only `ensign` is), so out of scope for this task. Note as a follow-up consideration.

### Composition with siblings

- `pi-launcher-repo-resolution` (sibling, same sprint) fixes the **parent-side** cwd-fallback. Together: parent gets the skill via `--skill` (explicit, install-recorded repo), child gets it via `.pi/skills/ensign` (project-declared). Both sides explicit, no cwd-luck.
- `pi-back-channel-dispatch` (capstone, same sprint) depends on this task — its `pi-live` drive requires a dispatched ensign that actually loaded the ensign contract.

## Acceptance criteria (provisional — ideation finalizes; proof = behavior, never prose-grep)

**AC-1 — `discoverAvailableSkills(cwd)` lists the ensign skill as discoverable post-fix.**
Verified by: running the real `discoverAvailableSkills(cwd)` from `pi-subagents/src/agents/skills.ts` (or `subagents-doctor` if available) post-fix and observing ensign in the list with source `project`. The spike already proved this (see Spike evidence); implementation reproduces it as a committed artifact.

**AC-2 — A `subagent(... skill:["ensign"])` dispatch loads the ensign contract with no "skills not found" warning, and the child exhibits ensign-contract behavior.**
Verified by: a probe dispatch (FO-side — the child worker has no `subagent` tool) whose run meta carries NO `skillsWarning` AND whose child exhibits ensign-contract behavior (in an ideation probe, produces a design in the entity body rather than editing product files — the inverse of the `0637e2ed` failure). This AC is FO-run: implementation hands the probe to the FO.

**AC-3 — The discovery mechanism is a committed project artifact that travels with the clone.**
Verified by: `.pi/skills/ensign` is a symlink committed to the repo (not `~/.pi/agent/` user config); `git log` shows it as a committed file; a fresh clone discovers ensign with no manual setup (the symlink resolves through the repo's own `skills/ensign/`).

## Out of scope

- The model-stamping friction (null model → settings default, not parent model) — separate concern, tracked alongside `pi-back-channel-dispatch`.
- Changing pi-subagents' search paths upstream (the fix is project-side discovery wiring, not a pi-subagents change).
- The ensign skill's *content* — this task is about *injectability*, not the contract text.

## Test plan

- **AC-1 — `discoverAvailableSkills(cwd)` (or `subagents-doctor`) post-fix**: live, cheap. The spike already proved the mechanism; implementation reproduces it as a committed artifact and re-runs the discovery to confirm. Cost: trivial (one `npx tsx` call or `subagents-doctor`).
- **AC-2 — Probe `subagent(... skill:["ensign"])` dispatch**: FO-side (child has no `subagent` tool). Run a bounded ideation-style probe dispatch; check run meta for no `skillsWarning`; verify child produces a design in the entity body rather than editing product files. Cost: one probe dispatch + verification. The probe is NOT a full stage drive — it tests injection + behavior, not stage completion.
- **AC-3 — Committed artifact check**: `git log -- .pi/skills/ensign` shows it as a committed symlink; verify the symlink resolves (`readlink .pi/skills/ensign` → `../../skills/ensign`, `ls .pi/skills/ensign/SKILL.md` → resolves). A fresh-clone reasoning suffices (the symlink target is within the repo, so any clone has it). Cost: trivial.

## No spike needed (beyond the one already run)

The riskiest mechanism — does a `.pi/skills/ensign` symlink actually make `discoverAvailableSkills(cwd)` list ensign? — was **spiked and PASSED** (see Approach → Spike evidence). The mechanism relies on proven stdlib symlink resolution + pi-subagents' existing `buildSkillPaths` filesystem scan. No unverified mechanism remains.

## Related

- `pi-back-channel-dispatch` (`b23y61pgk93ph44pz506m2wy`) — sibling Pi-dispatch-friction task; this one unblocks correct ensign dispatch for that and every subsequent Pi stage dispatch.
- `pi-launcher-repo-resolution` (`2m1cgn22ygmwtxe43z2hx7xw`) — sibling, fixes the parent-side cwd-fallback. Composes: parent gets the skill via `--skill`, child gets it via `.pi/skills/ensign`.
- Run `0637e2ed` — the failure instance.
- `node_modules/pi-subagents/src/agents/skills.ts:318` (`buildSkillPaths`) — the search-path source of truth.

## Stage Report: ideation

- DONE: Root cause reframed against the captain's constraint — the bug is a parent/child discovery split (parent uses `--skill` flags from the launcher, child uses its own `discoverAvailableSkills(cwd)` filesystem scan), not a cwd-accident. The child-side fix is a project-declared skill, distinct from the parent-side cwd-fallback fix (sibling task `pi-launcher-repo-resolution`).
- DONE: Riskiest mechanism SPIKED — created `.pi/skills/ensign -> ../../skills/ensign` symlink, ran the real `discoverAvailableSkills(cwd)` from `pi-subagents/src/agents/skills.ts` via `npx tsx`. Pre-fix: 2 skills, ensign ABSENT. Post-fix: 3 skills, ensign PRESENT, source `project` (highest priority). Symlink removed after spike (implementation creates it properly). Durable evidence recorded in the Approach → Spike evidence section.
- DONE: Mechanism decision finalized — (a) `.pi/skills/ensign -> ../../skills/ensign` symlink. Project-declared skill for child discovery; committed, travels with the clone, source of truth stays at `skills/ensign/`. Rationale recorded over (b) `.pi/settings.json` (extra config, no benefit), (c) `package.json pi.skills` (no root `package.json` exists — confirmed), (d) move (breaks parent-side discovery).
- DONE: ACs finalized — AC-1 `discoverAvailableSkills` lists ensign post-fix (live); AC-2 probe dispatch loads ensign with no `skillsWarning` AND child exhibits ensign-contract behavior (FO-side — child has no `subagent` tool); AC-3 committed project artifact, travels with clone.
- DONE: Test plan finalized — AC-1 live discovery re-run; AC-2 FO-side probe dispatch; AC-3 `git log` + `readlink` verification.
- DONE: No product files edited (ideation = design only). Main clean. No `.pi/` wiring created (spike cleaned up).

### Summary

Ideation complete for `pi-ensign-skill-injection`. The root cause is a parent/child skill-discovery split: the parent pi session gets the ensign skill via the launcher's `--skill` flags, but pi-subagents children use their own `discoverAvailableSkills(cwd)` filesystem scan that doesn't include `skills/ensign/SKILL.md`. The fix is a project-declared `.pi/skills/ensign -> ../../skills/ensign` symlink — a committed repo artifact that pi-subagents' `buildSkillPaths` discovers with `project` source priority. The mechanism was spiked live and PASSED (ensign goes from absent to present, source `project`). ACs and test plan are finalized and behavior-bound. No product files were edited; the spike symlink was cleaned up.

## Superseded (2026-06-19)

SUPERSEDED by the merged task `pi-install-managed-skill-placement` (filed 2026-06-19). Captain review uncovered that both this task's mechanism (the ` .pi/skills/ensign` repo symlink) and `pi-launcher-repo-resolution`'s (the cwd-fallback demotion + install-record file) are clone-bound workarounds for the fact that `spacedock install --host pi` is check-only (writes nothing). The correct mechanism — verified against the `obra/superpowers` reference and pi-subagents source — is install-managed package placement: ship Spacedock as a pi package (`package.json` with `pi.extensions` + `pi.skills`), make `spacedock install --host pi` run `pi install git:github.com/spacedock-dev/spacedock`, and let both the parent (via the extension's `resources_discover`) and pi-subagents children (via the package-root scan of `package.json -> pi.skills`) discover the skills. No clone, no cwd, no symlink. The merged task absorbs both this task's and `pi-launcher-repo-resolution`'s scope; the staff-review gap-1 `cwd:<repo>` wiring in the capstone becomes unnecessary once the install mechanism lands (child no longer cwd-keyed). Archived REJECTED as superseded — no deliverable merged (ideation-only).
