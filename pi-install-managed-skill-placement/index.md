---
title: Pi install-managed skill placement — ship Spacedock as a pi package; make spacedock install --host pi actually install
status: ideation
source: "Captain (2026-06-19): supersedes pi-ensign-skill-injection (k8t, archived REJECTED) and pi-launcher-repo-resolution (2m1, archived REJECTED). Both picked clone-bound workarounds (repo symlink; cwd-fallback record) for the fact that spacedock install --host pi is check-only (writes nothing). Verified against the obra/superpowers reference (.pi/extensions/superpowers.ts + package.json pi.skills) and pi-subagents source: the correct mechanism is install-managed package placement. pi install git:github.com/spacedock-dev/spacedock is the install source."
score:
started: 2026-06-19T22:41:57Z
completed:
verdict:
worktree:
issue:
sprint: 0223-pi-dispatch-contract
sprint-readiness: ready
id: eqrcrxcyye56nfwm997bj33d
---

# Pi install-managed skill placement

## End value

A user runs `spacedock install --host pi` (or the dev equivalent `pi install ./local/path` / `--plugin-dir`) and Spacedock is **actually installed** — the skills (`ensign`, `first-officer`, etc.) become discoverable by **both** the main pi session **and** pi-subagents children, with no clone, no cwd dependency, no symlink, and no launcher `--skill` flag. This mirrors how claude/codex already work (`ops.Install` issues the host's plugin-install) and retires Pi's status as the one host where install is check-only.

## Problem — root cause already determined

- `spacedock install --host pi` is **check-only** — it runs `checkPiRuntime` and prints a report; no `ops.Install`, no file writes (`internal/cli/pi.go:runInitWithPi`). Contrast claude/codex (`frontdoor.go:297,480`, `init.go:34`) which call `ops.Install(host, marketplaceSource, devBranch)` and durably register the plugin.
- Today the parent pi session discovers the Spacedock skills only because `spacedock pi` passes `--skill <repo>/skills/{ensign,first-officer}/SKILL.md` flags (`pi.go:87-89`), with `repo` resolved as `--plugin-dir` → `SPACEDOCK_REPO_ROOT` → **cwd** (fallback). pi-subagents **children** do NOT inherit those flags — they use `discoverAvailableSkills(cwd)` → `buildSkillPaths` (`pi-subagents/src/agents/skills.ts:318`), a filesystem scan of `.pi/skills/`, `.agents/skills/`, `~/.pi/agent/skills/`, package roots, settings. The repo's `skills/` is in none of them → `skill:["ensign"]` → "Skills not found: ensign".
- So the defect is a **parent/child discovery split** born of install being check-only. The prior tasks' mechanisms (a `.pi/skills/ensign` repo symlink; a cwd-fallback demotion + custom install-record file) were clone-bound workarounds that don't survive the actual install path.

## The verified mechanism — pi-installable package

Reference: `obra/superpowers` ships as a pi package with BOTH a pi extension (`.pi/extensions/superpowers.ts` implementing `resources_discover` → `{ skillPaths: [skillsDir] }`, resolved relative to the extension's own location) AND a `package.json` `pi.skills` entry (`"./skills"`). Verified in pi-subagents source: `collectInstalledPackageSkillPaths` (skills.ts:136) scans `{agentDir}/npm/node_modules/` and extracts `package.json -> pi.skills` per package — so a `pi install`'d package's skills are discovered by children with no cwd dependency. `pi install git:github.com/...` is a first-class source (`pi install --help` confirms); it registers the source in `settings.json` `packages`.

The design:
1. **Add a root `package.json`** declaring `name: "spacedock"`, `type: "module"`, and `pi: { extensions: ["./.pi/extensions/spacedock.ts"], skills: ["./skills"] }`.
2. **Add `.pi/extensions/spacedock.ts`** — a pi extension implementing `resources_discover` returning `{ skillPaths: [<repo>/skills] }` (path resolved relative to the extension's own location, like superpowers). This makes the **main pi session** discover the skills, replacing the launcher's `--skill` flags.
3. **`spacedock install --host pi`** stops being check-only — it runs `pi install git:github.com/spacedock-dev/spacedock` (or `pi install <local --plugin-dir path>` for dev), cloning/linking the repo into pi's package store and registering it in `settings.json` `packages`. Skills travel with the package.
4. **Both parent and child discover** — parent via the extension's `resources_discover`; child via pi-subagents' package-root scan reading `package.json -> pi.skills`. Unified, no `--skill` flag, no cwd dependency, no clone requirement, no symlink.
5. The launcher's `--skill` flags and the cwd fallback in `piRuntimeConfigFromEnv` are **retired** (or demoted below the install record for a dev-override transition).

## Scope

In scope:
- Root `package.json` with the `pi` block.
- `.pi/extensions/spacedock.ts` (the `resources_discover` extension; resolve skills dir relative to own location).
- `spacedock install --host pi` rewrite to run `pi install` (published `git:` source + local `--plugin-dir` dev override).
- Retirement of the `--skill` flags + cwd fallback in `piRuntimeConfigFromEnv` (or documented demotion).
- Verify both parent (extension) and child (package-root scan) discover the skills.

Out of scope:
- The back-channel / named-capability hardening — `pi-back-channel-dispatch` (capstone).
- Model stamping — `pi-dispatch-model-stamping`.
- npm publishing — `pi install git:...` is the source; no npm publish needed.
- First-officer skill injection specifically — `ensign` is the load-bearing one for FO dispatch, but the package ships ALL skills (`./skills`); per-skill injection is not separate work.

## Composition with the capstone (load-bearing)

This task **removes the staff-review gap-1 child-cwd seam**. Once install-managed placement lands, pi-subagents children discover skills via the package-root scan (not cwd-keyed symlink), so the capstone's `cwd: <resolved repo root>` wiring in `async-dispatch` and the AC-2 non-repo-launch pin become **unnecessary for skill discovery**. The capstone's ideation should be re-checked and the gap-1 fold-in revised (the `cwd:<repo>` wiring may still be useful for the child's working directory — ensigns work on the repo — but it's no longer the thing that makes ensign load). A parallel re-ideation of the capstone's fold-in accompanies this task.

## Acceptance criteria (provisional — ideation finalizes; proof = behavior, never prose-grep)

**AC-1 — `spacedock install --host pi` actually installs the package.**
Verified by: running `spacedock install --host pi` (or the dev `--plugin-dir` equivalent) and confirming the Spacedock package is registered in `~/.pi/agent/settings.json` `packages` and present in pi's package store. A live run, not a prose claim.

**AC-2 — Both the main pi session and pi-subagents children discover the Spacedock skills post-install, with no clone/cwd/symlink dependency.**
Verified by: (a) `subagents-doctor` lists `ensign` (and other Spacedock skills) as `user-package` source post-install; (b) a probe `subagent(... skill:["ensign"])` dispatch whose run meta carries NO `skillsWarning` AND whose child exhibits ensign-contract behavior (the inverse of the `0637e2ed` failure); (c) the probe is run from a NON-repo cwd, proving no cwd dependency.

**AC-3 — The launcher's `--skill` flags and cwd fallback are retired (or demoted below the install record with a documented transition).**
Verified by: a Go test that `piRuntimeConfigFromEnv` no longer relies on the cwd fallback when the package is installed; existing frontdoor tests stay green; the retirement is behavior-tested, not just documented.

**AC-4 — The dev override (`--plugin-dir` / `pi install ./local/path`) points at a local checkout so changes are picked up without reinstall.**
Verified by: a dev-mode install from a local path that discovers the skills from that checkout; a change to a skill file is visible without reinstall.

## Test plan

- Live `spacedock install --host pi` (AC-1) — `settings.json packages` + package store present.
- `subagents-doctor` + probe dispatch from non-repo cwd (AC-2) — the load-bearing behavioral proof.
- Go test for `piRuntimeConfigFromEnv` (AC-3) — cwd fallback retired/demoted.
- Dev-override install (AC-4) — local path discovery + live edit visibility.
- `pi-live` lane green (this task touches `internal/cli/pi.go` + adds `package.json` + `.pi/extensions/` — a pi-only surface; per the path→lane mapping, `pi-live` required).

## Related

- Supersedes `pi-ensign-skill-injection` (`k8tbnmcbyqc5kkhj0m9vewq4`, archived REJECTED) and `pi-launcher-repo-resolution` (`2m1cgn22ygmwtxe43z2hx7xw`, archived REJECTED) — both archived 2026-06-19 with supersession notes.
- Reference: `obra/superpowers` (`.pi/extensions/superpowers.ts` + `package.json pi.skills`) — the proven pattern.
- `pi-subagents/src/agents/skills.ts` (`buildSkillPaths`, `collectInstalledPackageSkillPaths`) — the child-discovery source of truth.
- `internal/cli/pi.go` (`runInitWithPi`, `piRuntimeConfigFromEnv`, `runPi`) — the install + resolution source of truth.
- `pi-back-channel-dispatch` (`b23y61pgk93ph44pz506m2wy`) — capstone; this task removes its staff-review gap-1 child-cwd seam. A parallel re-ideation of the capstone's fold-in accompanies this task.
- `0223-pi-dispatch-contract` sprint index.
