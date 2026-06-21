---
title: Pi dev-override skill loading — runPi must pass the Spacedock extension + skills when --plugin-dir is set
status: ideation
source: "Captain (2026-06-20): eq (#406) retired the --skill cfg.firstOfficerDir() / --skill cfg.ensignDir() flags (D4) and moved skill discovery to the .pi/extensions/spacedock.ts extension's resources_discover. But runPi never passes the Spacedock extension to pi — it passes only --extension <pi-subagents> + --skill <pi-subagents>. The installed path works by accident (pi auto-loads registered extensions from settings.json packages); the dev path (--plugin-dir .) is BROKEN — the Spacedock skills are not loaded. pi does not auto-discover .pi/extensions/ from cwd (verified empirically). This is a regression from eq, and the pi parity gap with z2t's --plugin-dir caveat (z2t AC-6 documents that --plugin-dir bypasses installed-plugin resolution; pi's --plugin-dir bypasses skill loading entirely)."
score:
started: 2026-06-21T04:22:56Z
completed:
verdict:
worktree:
issue:
sprint:
sprint-readiness: ready
id: ev8gzecy33zm84spxrj239md
---

# Pi dev-override skill loading

## End value

`./spacedock pi --plugin-dir .` loads the Spacedock first-officer + ensign skills from the local checkout — the FO boots with its contract, not skill-less. The dev-override path is self-sufficient: it doesn't depend on the package being separately `pi install`'d into `settings.json packages`. Parity with how Codex's `--plugin-dir` loads the local plugin checkout.

## Problem — root cause already determined

`eq` (#406, commit `a3f7efa8`) shipped the install-managed skill placement:
- **D1:** root `package.json` with `pi.extensions: ["./.pi/extensions/spacedock.ts"]` + `pi.skills: ["./skills"]`.
- **D2:** `.pi/extensions/spacedock.ts` — a pi extension implementing `resources_discover` returning `{ skillPaths: [<repo>/skills] }`.
- **D3:** `spacedock install --host pi` runs `pi install git:github.com/spacedock-dev/spacedock` → registers the package in `settings.json packages` → pi auto-loads the extension at startup → `resources_discover` fires → skills announced.
- **D4:** retired the `--skill cfg.firstOfficerDir()` / `--skill cfg.ensignDir()` flags from `runPi`.

**The gap:** D4 retired the dev-path skill loading (the `--skill` flags) before confirming the extension-based replacement covers the dev path. `runPi` currently passes:
```go
argv := []string{
    "pi",
    "--extension", cfg.extensionPath,        // pi-subagents extension (NOT Spacedock)
    "--skill", cfg.subagentsSkill,            // pi-subagents skill (NOT Spacedock)
}
```
The Spacedock extension (`.pi/extensions/spacedock.ts`) is **NOT passed**. pi does NOT auto-discover `.pi/extensions/` from cwd (verified empirically — a test extension in `/tmp/pi-ext-test/.pi/extensions/test-ext.ts` did NOT fire when `pi` was launched from that cwd). Extensions are only loaded if (a) registered via `pi install` in `settings.json packages`, OR (b) passed explicitly via `--extension <path>`.

**The installed path works by accident:** the Spacedock package is registered in `settings.json packages` (from a prior `pi install`), so pi auto-loads the extension → skills announced. But this is independent of `--plugin-dir` — it works because the package was separately installed, not because the dev-override does anything for skill loading.

**The dev path (`--plugin-dir .`) is broken:** no `pi install` in the dev path → the Spacedock extension isn't registered → not loaded → `resources_discover` never fires → **the Spacedock FO/ensign skills are not loaded**. The FO boots without its contract.

## Parity with z2t

`z2t` (`spacedock-marketplace-source-env`) is the Codex equivalent: it ensures `spacedock install --host codex` works programmatically, the resolver is channel-aware, and `--plugin-dir` is documented as "bypasses installed-plugin resolution." The pi parity gap: `z2t` AC-6 documents that Codex's `--plugin-dir` bypasses installed-plugin resolution; **pi's `--plugin-dir` bypasses skill loading entirely** — it's not just a resolution bypass, it's a skill-loading regression.

## Approach

`runPi` must pass the Spacedock extension + skills explicitly when the dev-override is active (`--plugin-dir` or `SPACEDOCK_REPO_ROOT` set, i.e. `cfg.repoRoot != ""`). This is the dev-override equivalent of what `pi install` registers for the installed path:

```go
// After the existing argv (pi, --extension <pi-subagents>, --skill <pi-subagents>),
// add the Spacedock extension + skills from the dev-override checkout.
if cfg.repoRoot != "" {
    spacedockExt := filepath.Join(cfg.repoRoot, ".pi", "extensions", "spacedock.ts")
    spacedockSkills := filepath.Join(cfg.repoRoot, "skills")
    if _, err := os.Stat(spacedockExt); err == nil {
        argv = append(argv, "--extension", spacedockExt, "--skill", spacedockSkills)
    }
}
```

This makes the dev path self-sufficient: `pi` loads the Spacedock extension via `--extension`, the extension fires `resources_discover`, and the skills are announced to the parent session. pi-subagents children discover the same skills via `collectSettingsPackageSkillPaths` (if the package is also installed) OR via the parent's `--skill` forwarding (pi-subagents inherits `--skill` flags? — verify; if not, the dev-path child discovery needs the same `--skill` on the subagent spawn, which is the Q3/cwd concern from 0223).

### Installed path (no --plugin-dir)

When `cfg.repoRoot == ""` (installed path, no dev override), the Spacedock extension is loaded by pi's registered-package mechanism (`settings.json packages` → pi auto-loads). `runPi` does NOT need to pass the extension explicitly — the installed path works. But verify this: does `runPi` need to pass `--extension <installed-package-path>/.pi/extensions/spacedock.ts` even in the installed case, or does pi auto-load registered extensions at startup? If pi auto-loads, the installed path is fine as-is. If not, `runPi` must resolve the installed package's extension path from `settings.json packages` and pass it.

## Scope

In scope:
- `runPi` passes `--extension <repo>/.pi/extensions/spacedock.ts` + `--skill <repo>/skills` when `cfg.repoRoot != ""` (dev-override path).
- Verify the installed path (no `--plugin-dir`) — does pi auto-load registered extensions, or does `runPi` need to pass the extension there too?
- A Go test asserting the dev-override argv includes the Spacedock extension + skills paths.
- A live/manual smoke: `./spacedock pi --plugin-dir . -- --version` (or a full FO boot) confirms the Spacedock skills load (the FO greets with the workflow state summary, proving the FO contract loaded).

Out of scope:
- The installed-path `pi install` mechanism (D3, already shipped by `eq`).
- The `z2t` Codex parity (separate task).
- The `t0g` binding-block restructure (separate task).
- The fnm-multishell race (already shipped #416).
- The prompt-warmth stopgap (already merged #417's sibling).

## Acceptance criteria (provisional — ideation finalizes; proof = behavior)

**AC-1 — `./spacedock pi --plugin-dir .` loads the Spacedock FO/ensign skills.**
Verified by: a live/manual smoke — launch `spacedock pi --plugin-dir .` from the repo root and confirm the FO boots with its contract (greets with the workflow state summary, not a skill-less prompt). The FO contract loading is the behavioral proof.

**AC-2 — `runPi` passes `--extension <repo>/.pi/extensions/spacedock.ts` + `--skill <repo>/skills` when `cfg.repoRoot != ""`.**
Verified by: a Go test with a fake `ops` capturing the argv; assert the dev-override argv includes the Spacedock extension + skills paths.

**AC-3 — The installed path (no `--plugin-dir`) still works.**
Verified by: a Go test or live smoke confirming the installed path (no `cfg.repoRoot`) does NOT add the dev-override extension/skill flags (the installed mechanism handles it), OR does add them from the resolved package path (if pi doesn't auto-load registered extensions — verify + record).

## Test plan

- Go unit test (AC-2): fake `ops` capturing argv; dev-override argv includes Spacedock extension + skills.
- Live/manual smoke (AC-1): `spacedock pi --plugin-dir .` boots the FO with the contract.
- `pi-live` lane (touches `internal/cli/pi.go` — pi-only surface).

## Related

- `eq` `pi-install-managed-skill-placement` (#406, merged) — retired the `--skill` flags (D4) that this task restores for the dev path.
- `z2t` `spacedock-marketplace-source-env` — the Codex parity task; AC-6 documents the `--plugin-dir` caveat this task closes for pi.
- `internal/cli/pi.go` (`runPi`, `piRuntimeConfigFromEnv`) — the source of truth.
- `.pi/extensions/spacedock.ts` — the extension that fires `resources_discover` (shipped by `eq`).
- `docs/runtime-support.md` — the runtime guide (may need updating: the `--plugin-dir` dev path now loads the extension explicitly).
