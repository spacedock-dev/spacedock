---
title: Pi install-managed skill placement — ship Spacedock as a pi package; make spacedock install --host pi actually install
status: implementation
source: "Captain (2026-06-19): supersedes pi-ensign-skill-injection (k8t, archived REJECTED) and pi-launcher-repo-resolution (2m1, archived REJECTED). Both picked clone-bound workarounds (repo symlink; cwd-fallback record) for the fact that spacedock install --host pi is check-only (writes nothing). Verified against the obra/superpowers reference (.pi/extensions/superpowers.ts + package.json pi.skills) and pi-subagents source: the correct mechanism is install-managed package placement. pi install git:github.com/spacedock-dev/spacedock is the install source."
score:
started: 2026-06-19T22:41:57Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-pi-install-managed-skill-placement
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

## Ideation finalization (2026-06-19)

### Mechanism confirmation (all verified against live code)

1. **`spacedock install --host pi` is check-only** — `internal/cli/pi.go:runInitWithPi` (line 100) runs `checkPiRuntime` and prints a report; no `ops.Install` call, no file writes. Confirmed: `grep -n "ops.Install" internal/cli/pi.go` returns nothing.
2. **The cwd fallback** — `piRuntimeConfigFromEnv` (pi.go:214) resolves `repo` as `--plugin-dir` → `SPACEDOCK_REPO_ROOT` → `dir` (cwd, fallback at ~line 226). The `--skill` flags (pi.go:87-89) pass `cfg.firstOfficerDir()` and `cfg.ensignDir()` to the parent pi session only.
3. **pi-subagents children use their own discovery** — `buildSkillPaths` (skills.ts:318) is a filesystem scan; children do NOT inherit the parent's `--skill` flags. The repo's `skills/` is in none of the scan paths.
4. **`pi install git:` is a first-class source** — `pi install --help` confirms `git:github.com/user/repo` is valid; it registers in `settings.json` `packages` (currently a string array: `["npm:pi-subagents", "npm:pi-intercom", ...]`).
5. **`collectSettingsPackageSkillPaths` (skills.ts:288) is the key discovery path for installed packages** — it reads `settings.json` `packages` entries and resolves each via `resolveSettingsPackageRoot` (skills.ts:267), which handles `git:` (→ `{agentDir}/git/{host}/{repoPath}`), `npm:` (→ `{agentDir}/npm/node_modules/{name}`), `file:`, `~/`, absolute, and relative paths. Then `extractSkillPathsFromPackageRoot` reads `package.json -> pi.skills` from the resolved package root.
6. **The superpowers reference** — `obra/superpowers` ships BOTH a pi extension (`.pi/extensions/superpowers.ts` implementing `resources_discover` → `{ skillPaths: [skillsDir] }`) AND `package.json` `pi.skills: ["./skills"]`. The extension handles parent-session discovery; the `pi.skills` entry handles child discovery via `collectSettingsPackageSkillPaths`.

### Spike evidence (live, 2026-06-19) — RISKIEST MECHANISM PROVEN

**Question:** Does `pi install` of a minimal package with `package.json pi.skills` make `discoverAvailableSkills` list the skill, from a non-repo cwd?

**Pre-install** (no spike package): `discoverAvailableSkills(cwd)` → total 2 (`find-skills`, `pi-intercom`). Spike skill ABSENT.

**Spike package** (`/tmp/pi-skill-spike/`): minimal `package.json` with `{ "name": "pi-skill-spike", "type": "module", "pi": { "skills": ["./skills"] } }` + `skills/spike-skill/SKILL.md`.

**Install:** `pi install /tmp/pi-skill-spike` → "Installed /tmp/pi-skill-spike". Registered in `settings.json packages` as `"../../../../tmp/pi-skill-spike"` (relative to agentDir).

**Post-install from the repo cwd:** `discoverAvailableSkills(cwd)` → total 3, including `spike-skill ( user-package )`. PRESENT.

**Post-install from `/tmp` (non-repo cwd):** `discoverAvailableSkills("/tmp")` → total 3, including `spike-skill ( user-package )`. **PRESENT from a non-repo cwd — no cwd dependency.**

**Cleanup:** `pi remove /tmp/pi-skill-spike` → cleanly unregistered. `settings.json packages` restored to original 3 entries.

**Conclusion:** The mechanism works end-to-end. `pi install` + `package.json pi.skills` makes pi-subagents children discover the skill via `collectSettingsPackageSkillPaths` with NO cwd dependency, NO symlink, NO clone requirement. The superpowers pattern (extension + `pi.skills`) is proven to work for Spacedock's repo shape.

### Approach (finalized)

Four deliverables, all in implementation's scope (ideation designs only):

**D1 — Root `package.json`.** Add at the repo root:
```json
{
  "name": "spacedock",
  "type": "module",
  "pi": {
    "extensions": ["./.pi/extensions/spacedock.ts"],
    "skills": ["./skills"]
  }
}
```
The `pi.skills` entry is what pi-subagents' `collectSettingsPackageSkillPaths` reads after `pi install` registers the package. The `pi.extensions` entry is what pi loads for the parent-session `resources_discover` hook.

**D2 — `.pi/extensions/spacedock.ts`.** A pi extension implementing `resources_discover` returning `{ skillPaths: [<repo>/skills] }`, with the path resolved relative to the extension's own location (`dirname(fileURLToPath(import.meta.url))` → `../..` → `skills/`), exactly like superpowers. This makes the main pi session discover the skills, replacing the launcher's `--skill` flags. The extension is the parent-side discovery mechanism; `pi.skills` is the child-side. Both are needed.

**D3 — `spacedock install --host pi` rewrite.** Stop being check-only. Run `pi install git:github.com/spacedock-dev/spacedock` (published source) or `pi install <local --plugin-dir path>` (dev override). This registers the package in `settings.json packages` and places the repo in pi's package store (`{agentDir}/git/github.com/spacedock-dev/spacedock` for git:, or the local path for dev). The existing doctor check stays as a post-install verification step. The `--plugin-dir` flag for install (currently rejected at pi.go:199) is lifted to specify a local dev source.

**D4 — Retirement of `--skill` flags + cwd fallback.** In `piRuntimeConfigFromEnv`, retire the `--skill cfg.firstOfficerDir()` and `--skill cfg.ensignDir()` flags from the `runPi` launch args (pi.go:87-89) — the extension's `resources_discover` handles parent discovery. Retire the cwd fallback (pi.go:226-228) — the installed package is found by `collectSettingsPackageSkillPaths` regardless of cwd. For the dev-override transition (a developer working in a clone without `pi install`), `--plugin-dir` / `SPACEDOCK_REPO_ROOT` still resolve the repo for the extension to find skills; the cwd fallback is demoted to a last resort with a warning, or removed entirely if the dev-override path is sufficient.

### Finalized ACs (behavior-bound, never prose-grep)

**AC-1 — `spacedock install --host pi` actually installs the package.**
Verified by: running `spacedock install --host pi` and confirming (a) `~/.pi/agent/settings.json` `packages` includes the Spacedock entry (`git:github.com/spacedock-dev/spacedock` or a local path), and (b) the repo is present in pi's package store (`{agentDir}/git/github.com/spacedock-dev/spacedock` for git:, or the local path). A live run, not a prose claim.

**AC-2 — Both the main pi session and pi-subagents children discover the Spacedock skills post-install, with no clone/cwd/symlink dependency.**
Verified by: (a) `subagents-doctor` (or `discoverAvailableSkills(cwd)`) lists `ensign` (and other Spacedock skills) as `user-package` source post-install; (b) a probe `subagent(... skill:["ensign"])` dispatch whose run meta carries NO `skillsWarning` AND whose child exhibits ensign-contract behavior (the inverse of the `0637e2ed` failure — ideation probe = design in entity body, not product edits); (c) the probe is run from a NON-REPO cwd, proving no cwd dependency. The spike proved this mechanism works (see Spike evidence).

**AC-3 — The launcher's `--skill` flags and cwd fallback are retired (or demoted below the install record with a documented transition).**
Verified by: a Go test that `piRuntimeConfigFromEnv` no longer emits `--skill cfg.firstOfficerDir()` / `--skill cfg.ensignDir()` in the launch args when the package is installed (the extension handles parent discovery), and that the cwd fallback is either removed or demoted below the install-record resolution. Existing frontdoor tests stay green. The retirement is behavior-tested, not just documented.

**AC-4 — The dev override (`--plugin-dir` / `pi install ./local/path`) points at a local checkout so changes are picked up without reinstall.**
Verified by: a dev-mode install from a local path (`pi install ./local/path` or `spacedock install --host pi --plugin-dir <local>`) that discovers the skills from that checkout; a change to a skill file (e.g. appending a comment to `skills/ensign/SKILL.md`) is visible in `discoverAvailableSkills` without reinstall — the local path is read live, not copied.

### Test plan

- **Live `spacedock install --host pi`** (AC-1) — `settings.json packages` + package store present. Bounded live run.
- **`subagents-doctor` + probe dispatch from non-repo cwd** (AC-2) — the load-bearing behavioral proof. `discoverAvailableSkills("/tmp")` lists ensign; a `subagent(... skill:["ensign"])` dispatch from `/tmp` has no `skillsWarning` and the child behaves as an ensign. The spike proved the discovery mechanism; implementation reproduces it for the real Spacedock package.
- **Go test for `piRuntimeConfigFromEnv`** (AC-3) — assert the `--skill` flags are absent from launch args when installed; assert the cwd fallback is removed/demoted. Existing frontdoor tests green (non-regression).
- **Dev-override install** (AC-4) — `pi install ./local/path` discovers skills; edit a skill file; `discoverAvailableSkills` reflects the edit without reinstall.
- **`pi-live` lane green** — this task touches `internal/cli/pi.go` (high-stakes: the front-door launcher) + adds `package.json` + `.pi/extensions/`. Per the path→lane mapping, `pi-live` required. The `package.json` + `.pi/extensions/` are pi-convention artifacts; `internal/cli/pi.go` is the launcher change.

### No further spike needed

The riskiest mechanism — `pi install` + `package.json pi.skills` making pi-subagents discover the skill from a non-repo cwd — is PROVEN by the spike above. The extension (`resources_discover`) pattern is proven by the superpowers reference. The `resolveSettingsPackageRoot` path resolution for `git:`/local sources is confirmed in source. No unverified mechanism remains.

## Stage Report: ideation

- DONE: Mechanism confirmed against live code — `spacedock install --host pi` is check-only (no `ops.Install`); cwd fallback at pi.go:226; `--skill` flags at 87-89; `collectSettingsPackageSkillPaths` (skills.ts:288) reads `settings.json packages` via `resolveSettingsPackageRoot` (skills.ts:267) which handles `git:`/`npm:`/`file:`/local paths.
- DONE: Riskiest mechanism SPIKED PASSED — `pi install /tmp/pi-skill-spike` (minimal package with `pi.skills`) → `discoverAvailableSkills` lists `spike-skill` as `user-package` from both repo cwd AND `/tmp` (non-repo). No cwd dependency. Cleaned up with `pi remove`.
- DONE: Approach finalized — 4 deliverables (root `package.json` with `pi` block; `.pi/extensions/spacedock.ts` with `resources_discover`; `spacedock install --host pi` rewrite to run `pi install`; retirement of `--skill` flags + cwd fallback).
- DONE: ACs finalized (4, behavior-bound) — install actually installs; both parent+child discover from non-repo cwd; `--skill`/cwd retired or demoted; dev override works with live edits.
- DONE: Test plan finalized — live install, subagents-doctor + probe from non-repo, Go test for config resolution, dev-override live edit, `pi-live` lane.
- DONE: Composition with capstone recorded — this task removes the staff-review gap-1 child-cwd seam; the capstone's `cwd:<repo>` wiring is no longer needed for skill discovery (may remain for working-directory concern).

### Summary

Ideation complete. Mechanism confirmed against live code + superpowers reference. Riskiest mechanism spiked PASSED (pi install + pi.skills → discoverable from non-repo cwd). Four deliverables designed, four behavior-bound ACs finalized, test plan covers live install + probe + Go test + dev override + pi-live lane. No product files edited (ideation = design only). Entity body + stage report committed to state checkout.

## Staff review #2 fold-in (2026-06-19)

### Gap 1 — `repoRoot` source post-install undefined (spec gap, closed)

**The gap (verified against live code).** The doctor at `internal/cli/pi.go:293-294` computes `firstOfficerOK`/`ensignOK` via `ops.Stat(cfg.firstOfficer)` / `ops.Stat(cfg.ensign)`, and `cfg.firstOfficer`/`cfg.ensign` derive from `repoRoot` (`pi.go:263-264`: `filepath.Join(repo, "skills", "{first-officer,ensign}", "SKILL.md")`). `piRuntimeConfigFromEnv` (`pi.go:220-228`) resolves `repo` as `--plugin-dir` → `SPACEDOCK_REPO_ROOT` → **cwd** (fallback). `piRuntimeLaunchReady` (`pi.go:299`) gates launch readiness on `firstOfficerOK && ensignOK`.

This task's D3 (run `pi install`) writes NO repo-path install record — it registers a package in `settings.json` `packages` and places the repo in pi's package store. This task's D4 retires the `--skill` flags and the cwd fallback. But AC-3 references an "install-record resolution" that D3 does not produce.

**Consequence (verified):** under the headline scenario (parent launched from a non-repo cwd, package installed, no `--plugin-dir`/`SPACEDOCK_REPO_ROOT`), `repoRoot` resolves to cwd (the demoted/retired fallback) → the doctor's `Stat(cfg.firstOfficer/ensign)` checks a NONEXISTENT path → `firstOfficerOK`/`ensignOK` report FALSE → `piRuntimeLaunchReady` returns FALSE. The doctor would report the runtime NOT ready and the skills broken, even though the skills ARE discoverable via the package-root scan (`collectSettingsPackageSkillPaths`). This is an internal inconsistency in this task's own AC-3 (it names a producer its design does not define).

### Resolution — (a) retire the repo-path Stat checks in favor of package-registration + package-root-scan verification

**Decision: resolution (a).** The doctor RETIRES its repo-path `Stat(cfg.firstOfficer/ensign)` checks (pi.go:293-294, 324-325) in favor of confirming the package is registered AND discoverable via the real mechanism this task ships.

**Rationale:** D4 retires the `--skill cfg.firstOfficerDir()` / `--skill cfg.ensignDir()` launch flags (pi.go:87-89) — the parent-side discovery mechanism those Stat checks test. Once the flags are gone, `cfg.firstOfficer`/`cfg.ensign` (the paths the flags pointed at) are no longer the discovery mechanism for EITHER the parent (the extension's `resources_discover` handles parent discovery) OR the child (the package-root scan handles child discovery). The Stat-based checks are testing a retired mechanism. The real discovery mechanism is: (parent) the extension's `resources_discover` returning the skills dir; (child) `collectSettingsPackageSkillPaths` reading `settings.json packages` → `package.json pi.skills`. The doctor should test THAT.

Rejected (b) (launcher reads the package store path from `settings.json packages` and sets `repoRoot` to the installed package's root, keeping the Stat checks repointed at the installed package). (b) preserves the doctor's existing shape but adds plumbing to resolve the package store path into `repoRoot`, and it tests a path (`Stat` of a file inside the installed package) that is an implementation detail of the package store, not the actual discovery contract. (a) is cleaner: the doctor verifies the contract (package registered + discoverable), not a filesystem coincidence.

**New D5 — Doctor skill-check retirement + package-registration verification.** Append to the approach (after D4):

- **D5a — Retire `firstOfficerOK`/`ensignOK` as repo-path Stat checks.** Remove the `res.firstOfficerOK = ops.Stat(cfg.firstOfficer) == nil` / `res.ensignOK = ops.Stat(cfg.ensign) == nil` lines (pi.go:293-294). Remove the corresponding `printPiCheck` lines (pi.go:324-325). Remove `firstOfficerOK`/`ensignOK` from `piRuntimeLaunchReady` (pi.go:299). The `cfg.firstOfficer`/`cfg.ensign` fields and their derivation (pi.go:263-264) are also retired — they exist only to feed the retired flags/checks.
- **D5b — Add `spacedockPackageOK` check.** Replace the retired skill checks with a single check: is the Spacedock package registered in `settings.json packages` AND discoverable via the package-root scan? Concretely: read `~/.pi/agent/settings.json` `packages` (the same source `collectSettingsPackageSkillPaths` reads); confirm a Spacedock entry is present; confirm `discoverAvailableSkills(cwd)` (or the equivalent internal call) lists `ensign` (and `first-officer`) as `user-package` source. `spacedockPackageOK` replaces `firstOfficerOK && ensignOK` in `piRuntimeLaunchReady`. The doctor report prints `spacedockPackageOK` with the registered source and the discovered skills, with the remedy "run `spacedock install --host pi`" on failure.
- **D5c — `repoRoot` resolution.** With the Stat checks and `--skill` flags retired, `repoRoot` is no longer needed for skill discovery (the package-root scan is cwd-independent). `repoRoot` is retained ONLY for the dev-override path: when `--plugin-dir` or `SPACEDOCK_REPO_ROOT` is set (a developer working in a clone without `pi install`), `repoRoot` resolves to that path and the extension's `resources_discover` uses it. When NEITHER is set AND the package is installed, `repoRoot` is empty/unused — the package store is the source. The cwd fallback is removed entirely (not demoted) — it served only the retired `--skill` flags, and leaving it would silently re-introduce the "skills found at `<cwd>/skills/`" false-positive/negative the package mechanism eliminates. AC-3's "demoted below the install record" option is withdrawn; the cwd fallback is removed.

**Composition consumer — capstone `cwd:<repo>` source.** The capstone (`pi-back-channel-dispatch`) sources its `cwd:<repo>` working-directory argument from "the same install-recorded / explicitly-resolved repo path the launcher records." With this fold-in, that source is now: `--plugin-dir` / `SPACEDOCK_REPO_ROOT` (dev override), OR the installed package's resolved root (read from `settings.json packages` via `resolveSettingsPackageRoot` — the same call `collectSettingsPackageSkillPaths` uses). The capstone's gap-1 re-check claim that the "source is unchanged" is TRUE under this fold-in: the source is the resolved package root (install case) or the explicit override (dev case), not `2m1`'s retired install-record file. The capstone should source `cwd:<repo>` from this resolved path.

### AC-3 revision

**AC-3 (revised) — The launcher's `--skill` flags, the cwd fallback, and the repo-path Stat skill checks are retired; the doctor verifies the package is registered and discoverable.**
Verified by: a Go test that (a) `piRuntimeConfigFromEnv` no longer emits `--skill cfg.firstOfficerDir()` / `--skill cfg.ensignDir()` in the launch args (the extension handles parent discovery); (b) the cwd fallback is removed (not demoted) — `repoRoot` is empty when no `--plugin-dir`/`SPACEDOCK_REPO_ROOT` is set and the package is installed; (c) `checkPiRuntime` no longer computes `firstOfficerOK`/`ensignOK` via `Stat(cfg.firstOfficer/ensign)` — instead it computes `spacedockPackageOK` (package registered in `settings.json packages` AND `ensign` discoverable as `user-package` source); (d) `piRuntimeLaunchReady` gates on `spacedockPackageOK`, not on the retired `firstOfficerOK`/`ensignOK`; (e) the doctor reports `spacedockPackageOK` TRUE from a NON-REPO cwd when the package is installed (the headline scenario that was broken before this fold-in). Existing frontdoor tests stay green (non-regression). The retirement and the new check are behavior-tested, not just documented.

### What this fold-in does NOT change

- D1 (root `package.json`), D2 (`.pi/extensions/spacedock.ts`), D3 (`spacedock install --host pi` runs `pi install`) are unchanged.
- AC-1, AC-2, AC-4 are unchanged.
- The spike evidence (package-root scan discovers skills from non-repo cwd) is unchanged and still authoritative.
- The capstone's gap-1 re-check (cwd:<repo> reframed as working-directory concern) is unchanged in stance; this fold-in specifies the `repoRoot` SOURCE the re-check left undefined.

## Stage Report: ideation (staff review #2 fold-in 2026-06-19)

- DONE: Verified gap 1 against live code — doctor's `firstOfficerOK`/`ensignOK` (pi.go:293-294) are `Stat(cfg.firstOfficer/ensign)` from `repoRoot` (pi.go:263-264); `piRuntimeLaunchReady` (pi.go:299) gates on them; `repoRoot` falls back to cwd (pi.go:220-228). Under the non-repo-cwd installed scenario, the doctor reports NOT ready despite skills being discoverable. Real spec gap.
- DONE: Picked resolution (a) — retire the repo-path Stat checks (D4 retires the `--skill` flags they test; the real mechanism is package-registration + package-root scan). Rejected (b) (repoint Stat at installed package) — tests a filesystem coincidence, not the discovery contract.
- DONE: Appended D5 (retire Stat checks + `cfg.firstOfficer/ensign`; add `spacedockPackageOK` checking package registration + package-root-scan discovery; remove cwd fallback entirely, not demote; `repoRoot` retained only for dev-override path).
- DONE: Revised AC-3 — `--skill` flags + cwd fallback + repo-path Stat checks all retired; doctor verifies package registered + discoverable; `piRuntimeLaunchReady` gates on `spacedockPackageOK`; doctor reports TRUE from non-repo cwd when installed. Behavior-bound Go test.
- DONE: Specified the capstone's `cwd:<repo>` source — `--plugin-dir`/`SPACEDOCK_REPO_ROOT` (dev) or the installed package's resolved root (via `resolveSettingsPackageRoot`). The capstone's re-check "source is unchanged" claim is now true.

### Summary

Gap 1 (repoRoot source post-install) closed. Resolution (a): the doctor retires its repo-path Stat skill checks (which test the retired `--skill` flags) in favor of `spacedockPackageOK` (package registered + discoverable via package-root scan). The cwd fallback is removed entirely (not demoted). AC-3 revised to behavior-test the retirement + the new check from a non-repo cwd. The capstone's `cwd:<repo>` source is specified. Append-only; prior approach/ACs/test-plan preserved (AC-3 superseded by this fold-in's revision). No product files edited (ideation = design only).

## Stage Report: implementation

- DONE (AC-1 — install actually installs): `spacedock install --host pi --plugin-dir <worktree>` runs `pi install <worktree>` (the dev-override local path). Live run output: "Installing <worktree>... Installed <worktree>". Registered entry in `~/.pi/agent/settings.json` `packages`:
  ```
  "../../git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-pi-install-managed-skill-placement"
  ```
  (relative to agentDir `~/.pi/agent`; resolves to the absolute worktree path). The resolved package root contains `package.json` with `name: "spacedock"` and `pi.skills: ["./skills"]`; `skills/ensign/SKILL.md` exists there. For a local-path dev install the repo IS the package source (read live, not copied); for the published `git:github.com/spacedock-dev/spacedock` source pi clones into `{agentDir}/git/github.com/spacedock-dev/spacedock`. The post-install doctor reports `OK Spacedock package: <worktree>` and `Pi runtime ready.` — a live run, not a prose claim.
- DONE (AC-2a — discoverable as user-package from non-repo cwd): cleaned the stale `~/.pi/agent/skills/{ensign,first-officer,...}` symlinks left by the rejected k8t task (they pointed at a different worktree and would have false-positive'd the `user` source). Post-install, pi-subagents' own `discoverAvailableSkills("/tmp")` (non-repo cwd) lists `ensign ( user-package )` and `first-officer ( user-package )` — plus all other Spacedock skills (commission, debrief, feedback-rejection-flow, present-gate, refit, survey, using-legacy-claude-team) as `user-package`. No cwd dependency (run from /tmp), no symlink, no clone requirement. This reproduces the spiked mechanism for the real Spacedock package.
- DONE (AC-2b — probe dispatch carries no skillsWarning): replicated pi-subagents' `skillsWarning(cwd, skills)` EXACT logic (agent-management.ts:155-160: `missing = skills.filter(s => !new Set(discoverAvailableSkills(cwd).map(s=>s.name)).has(s))`) from `/tmp`:
  - `skill:["ensign"]` → `undefined` (NO warning)
  - `skill:["ensign","first-officer"]` → `undefined` (NO warning)
  - control `skill:["nonexistent-skill-xyz"]` → `"Warning: skills not found: nonexistent-skill-xyz."` (proves the probe logic is sound)
  This is the precise inverse of the `0637e2ed` failure ("skills not found: ensign").
  NESTING CONSTRAINT (recorded honestly per scope note): I am myself a dispatched child worker; my tool surface does NOT include the pi-subagents `subagent` tool (it is provided to the PARENT session by the pi-subagents extension), so I could not dispatch a nested `subagent(... skill:["ensign"])` from within this child to capture the child's run-meta + ensign-contract behavior directly. I proved the load-bearing precondition the probe depends on — the exact `skillsWarning` computation returns no warning for `skill:["ensign"]` from a non-repo cwd — which is the condition that was false at `0637e2ed`. The FO (parent) can run the actual `subagent(... skill:["ensign"])` probe dispatch from the parent level at validation if the live run-meta capture is required.
- DONE (AC-3 — retirement + new check, behavior-tested): Go test `TestPiRuntimeConfigRetiresSkillFlagsAndCwdFallback` asserts (a) `piRuntimeConfigFromEnv` with no `--plugin-dir`/`SPACEDOCK_REPO_ROOT` yields `repoRoot == ""` from a non-repo cwd (cwd fallback REMOVED, not demoted); (b) the dev override is retained (`--plugin-dir`/`SPACEDOCK_REPO_ROOT`); (c) `runPi` launch args contain exactly ONE `--skill` flag (pi-subagents only) — the retired `--skill <repo>/skills/{first-officer,ensign}` flags are absent; (d) `checkPiRuntime` computes `spacedockPackageOK` (via `ops.SpacedockPackageStatus`) instead of the retired `firstOfficerOK`/`ensignOK` Stat checks, and `piRuntimeLaunchReady` gates on it; (e) the doctor reports `MISSING Spacedock package` (non-zero) when not installed and `OK Spacedock package` (zero) when installed, both from a non-repo cwd. `go test ./internal/cli/ ./internal/piruntime/ -race` green. `go test ./...` green except a PRE-EXISTING unrelated failure in `internal/status` (`TestMigrationCheckFixturesParseConsistently` — a `docs/dev/_debriefs/2026-06-19-01.md` frontmatter `session-date` fixture issue) that fails identically on the clean baseline (verified by stashing my changes) and touches none of my files.
- DONE (AC-4 — dev override live edit visibility): with the worktree installed via the local-path dev override, adding a marker skill `skills/ac4-live-probe/SKILL.md` to the worktree makes `discoverAvailableSkills("/tmp")` list `ac4-live-probe ( user-package )` WITHOUT reinstall (fresh process, no cache); removing it makes it disappear. The local path is read live, not copied.
- DONE (AC-3e live): the built binary `spacedock doctor --host pi` run from `/tmp` (non-repo cwd) reports `OK Spacedock package: <worktree>` and exits 0 — the real `spacedockPackageOK` check reading the actual `~/.pi/agent/settings.json` via the replicated `resolveSettingsPackageRoot` + `extractSkillPathsFromPackageRoot` logic.

### Implementation notes

- `piRuntimeOps` gained two methods: `PiInstall(source)` (the install seam; real impl execs `pi install <source>`) and `SpacedockPackageStatus(agentDir, home)` (the discovery check; real impl replicates `collectSettingsPackageSkillPaths` → `resolveSettingsPackageRoot` + `package.json pi.skills` + `ensign/SKILL.md`). Tests inject canned values via the fake.
- `parsePiSetupArgs` now accepts `--plugin-dir` for install (the pi dev-override source); non-pi install rejects it in `runInitWithPi` (preserving the claude/codex `--plugin-dir is not supported` contract); non-pi doctor still rejects via the `runDoctor` re-parse.
- The `--skill cfg.subagentsSkill` (pi-subagents) flag is RETAINED in `runPi` — only the Spacedock first-officer/ensign flags were retired (the extension's `resources_discover` handles parent discovery of those).
- `pi-live` lane: the `//go:build live` pi runner/ensign tests construct prompts that reference skill paths directly and pass `--plugin-dir`; they are out of the baseline `go test ./...` gate. Retiring the launcher `--skill` flags changes the parent's skill-registry path but the live tests' mechanics rely on explicit path references in prompts, not the flags. The `pi-live` lane should be re-run at validation; if the parent's `$spacedock:first-officer` skill reference needs the package installed in the live-test env, that is the install-managed contract this task ships.
- Pre-existing dirty files `internal/cli/prose_function_routing_test.go` and `internal/status/section_read.go` (smart-quote comment corruption) were present on arrival and are NOT part of this commit; only the 6 task files were committed.

### Summary

Implementation complete. `spacedock install --host pi` now actually installs (runs `pi install`), shipping Spacedock as a pi package via root `package.json` (`pi.skills` + `pi.extensions`) and `.pi/extensions/spacedock.ts` (`resources_discover`). The launcher's `--skill` first-officer/ensign flags and the cwd fallback are retired; the doctor gates on `spacedockPackageOK` (package registered + ensign discoverable via the package-root scan), proven live from a non-repo cwd. AC-1/AC-2(a+b-precondition)/AC-3/AC-4 proven with live evidence; AC-2b's nested-probe-dispatch run-meta capture is left to the FO at validation due to the child-has-no-subagent-tool nesting constraint (the load-bearing `skillsWarning` precondition is proven). Go tests green (race-clean on touched packages); one pre-existing unrelated status-fixture failure noted.
