---
title: Pi launcher repo resolution — stop resolving the Spacedock repo from cwd; record it at install
status: ideation
source: "Captain (2026-06-19, 0223-pi-dispatch-contract sprint scope-lock): `spacedock pi` resolves repoRoot as --plugin-dir -> SPACEDOCK_REPO_ROOT -> working directory (cwd). The cwd fallback means running the homebrew binary from any directory points `--skill <cwd>/skills/ensign/SKILL.md` wherever you happen to stand — the skill is registered by cwd-luck, not explicit install. Captain: 'it should not use the current directory's skill without me explicitly running pi with --skills.'"
score:
started: 2026-06-19T21:53:58Z
completed:
verdict:
worktree:
issue:
sprint: 0223-pi-dispatch-contract
sprint-readiness: ready
id: 2m1cgn22ygmwtxe43z2hx7xw
---

# Pi launcher repo resolution

## End value

`spacedock pi` (and `spacedock install --host pi`) resolve the Spacedock skill repo **explicitly** — via `--plugin-dir`, `SPACEDOCK_REPO_ROOT`, or a path **recorded at install time** — never by falling back to the working directory. Running the homebrew `spacedock` binary from an unrelated directory no longer silently points `--skill` at a nonexistent or wrong path.

## Problem — root cause already determined

`internal/cli/pi.go` `piRuntimeConfigFromEnv` (line 214) resolves `repo` in priority order:

1. `--plugin-dir` flag
2. `SPACEDOCK_REPO_ROOT` env
3. **`dir` (the working directory)** — fallback

The fallback feeds `cfg.firstOfficer` (`<repo>/skills/first-officer/SKILL.md`) and `cfg.ensign` (`<repo>/skills/ensign/SKILL.md`), which `spacedock pi` passes as `--skill` flags to the `pi` binary (pi.go:87-89). So when launched without `--plugin-dir` or `SPACEDOCK_REPO_ROOT`, the skill path is `<cwd>/skills/...` — correct only by cwd-luck when cwd happens to be the Spacedock repo. `install --host pi`'s "Spacedock skills: <repo>" output is this resolved `repoRoot`; in the FO session that triggered this task it was the cwd.

This is the parent-side (launcher) half of the cwd-accident concern. The child-side half (pi-subagents not discovering the skill) is `pi-ensign-skill-injection`; this task is the launcher half — they compose: the launcher registers the skill explicitly for the parent, and the project-declared skill makes it discoverable for children.

## Parity analysis — claude/codex vs pi (confirmed in live code)

**Claude/Codex do NOT have this problem.** They delegate skill discovery to the host's native plugin system:

- **Install** (`internal/cli/init.go`): `ops.Install(host, marketplaceSource, devBranch)` — drives the host's plugin install verb (`claude plugin add` / `codex plugin add`) from the marketplace. The host durably records the plugin location in its own plugin registry. No spacedock-side repo path resolution.
- **Launch** (`internal/cli/frontdoor.go`): `--plugin-dir` is a **passthrough** flag to the host binary — it goes into `fd.passthrough` and is forwarded to `claude`/`codex` verbatim. Spacedock does NOT resolve it to a repo path or pass `--skill` flags. The host's own plugin system discovers skills from the plugin dir or the installed marketplace plugin.
- `repoRoot` in `frontdoor.go` (line 164) is used ONLY for **workflow detection** (the launch banner), not for skill resolution.
- There is **no cwd fallback** for skill resolution in claude/codex — skills are either installed via the marketplace (host-native discovery) or loaded via `--plugin-dir` passthrough (host-native plugin-dir support).

**Pi is unique:** Pi has NO native plugin mechanism. `spacedock pi` must resolve the repo path ITSELF and pass `--skill` flags explicitly (pi.go:87-89). `spacedock install --host pi` is a **check-only no-op** — it does NOT write anything (no `ops.Install` call, no file writes; confirmed: `grep -n 'ops.Install\|WriteFile' internal/cli/pi.go` returns empty). It just runs `checkPiRuntime` and reports readiness. Additionally, `spacedock install --host pi` **rejects `--plugin-dir`** (pi.go:199: "is not supported; use SPACEDOCK_REPO_ROOT or run from the Spacedock checkout"), so you can't even tell install which repo to record.

**Parity decision (AC-3): the fix is pi-only.** Claude/codex don't need it — their skill discovery is host-native via the plugin marketplace. Applying a repo-path record to all hosts would add a redundant config layer to hosts that already handle this. The pi-specificity is structurally justified: Pi is the only host where spacedock owns skill-path resolution.

## Mechanism decision — (a) install records the repo path durably

**Picked: (a)** — `spacedock install --host pi` writes the resolved repo path to a durable record, and `spacedock pi` reads it as a resolution source ahead of the cwd fallback.

**Rationale over alternatives:**
- vs **(b) remove cwd fallback entirely**: non-breaking. Operators who run `spacedock pi` from the repo root (the common dev workflow) still work — the cwd fallback remains as a last resort, but the install record takes priority. (b) would break muscle memory and is unnecessarily hostile.
- vs **(c) warn-only**: a warning is a prose-only gate — banned by proof policy as the sole assurance. The install record is a durable behavioral change (the path IS resolved from it), not a label.

**New resolution priority in `piRuntimeConfigFromEnv`:**
1. `--plugin-dir` flag (unchanged)
2. `SPACEDOCK_REPO_ROOT` env (unchanged)
3. **install-recorded path** (NEW — read from a durable file written by `spacedock install --host pi`)
4. `dir` / cwd (demoted to last resort)

**Install action changes:**
- `spacedock install --host pi` currently rejects `--plugin-dir` (pi.go:199). The fix should **accept `--plugin-dir`** (or `SPACEDOCK_REPO_ROOT`) at install time so the operator can specify which repo to record, then write the resolved path to the durable record. When neither is given, it uses cwd (same as today) but NOW records it — so subsequent launches from elsewhere resolve correctly.
- The install record file: a standalone file containing the absolute repo path (no JSON parsing, easy read/write). Candidate location: `~/.pi/agent/spacedock-repo-path` (consistent with existing pi state under `~/.pi/agent/`). Implementation picks the exact path; ideation recommends this location for consistency with `auth.json`, `sessions/`, `npm/` siblings.
- `spacedock install --host pi`'s success output adds the recorded path: `Spacedock skills: <path> (recorded)`.

**What stays the same:**
- `--plugin-dir` on `spacedock pi` launch still works and still wins (highest priority).
- `SPACEDOCK_REPO_ROOT` still works and still wins over the install record.
- The cwd fallback stays as a last resort — non-breaking, but no longer the default when an install record exists.
- `spacedock doctor --host pi` reads the same resolution path and reports the source (install-record vs cwd) so the operator can see where the repo path came from.

## Acceptance criteria (finalized — proof = behavior, never prose-grep)

**AC-1 — `spacedock pi` launched from a cwd that is NOT the Spacedock repo resolves the correct skill paths after install.**
Verified by: a Go behavior test that (1) runs `spacedock install --host pi` from a temp repo dir, (2) then resolves `piRuntimeConfigFromEnv` from a DIFFERENT cwd (not the repo), and (3) asserts `cfg.firstOfficer` and `cfg.ensign` point at the install-recorded repo, not the cwd. The test does NOT shell out to `spacedock pi` (that's a live-lane concern); it exercises the config resolution function directly. Additionally, a `pi-live` lane run confirms `spacedock doctor --host pi` from a non-repo cwd reports the install-recorded path.

**AC-2 — The repo path is recorded at install time and read at launch; the cwd fallback is demoted below the install record.**
Verified by: a Go test that (1) calls the install path (writing the record), (2) verifies the record file exists and contains the resolved repo path, (3) calls `piRuntimeConfigFromEnv` with no `--plugin-dir`, no `SPACEDOCK_REPO_ROOT`, a non-repo cwd, and the install record present, and (4) asserts `cfg.repoRoot == <install-recorded path>` and `cfg.pluginDirSource == "install record"` (or equivalent). A second test confirms the cwd fallback still works when NO install record exists (non-breaking).

**AC-3 — The fix is pi-only; claude/codex are unaffected.**
Verified by: the change is confined to `internal/cli/pi.go` (and any new install-record helper) — no edits to `internal/cli/frontdoor.go` or `internal/cli/init.go`'s claude/codex paths. A Go test confirms claude/codex launch paths are unchanged (existing frontdoor tests stay green). The parity decision is documented above with the structural reason: Pi is the only host where spacedock owns skill-path resolution.

**AC-4 — `spacedock install --host pi` accepts `--plugin-dir` to specify the repo to record.**
Verified by: a Go test that runs install with `--plugin-dir <temp-dir>` and asserts the record contains that path (not cwd). Today install rejects `--plugin-dir` (pi.go:199); the fix lifts that restriction for the install command.

## Out of scope

- Child-side skill discovery (pi-subagents) — that is `pi-ensign-skill-injection`. This task is the launcher's repo-resolution only.
- The back-channel / named-capability hardening — that is `pi-back-channel-dispatch`.
- Model stamping — that is `pi-dispatch-model-stamping`.
- Removing the cwd fallback entirely (option b) — non-breaking is the priority; cwd stays as last resort.
- Cross-host repo-path record (applying to claude/codex) — they don't need it (host-native plugin discovery; AC-3).

## Test plan

**Go unit tests (AC-1, AC-2, AC-4):**
- `TestPiRuntimeConfigResolvesFromInstallRecord`: install writes record → config from non-repo cwd resolves to install-recorded path. Follows the existing `TestPiRuntimeConfigResolvesEnvPathsForSubagentsIntercomAuthAndSessions` pattern.
- `TestPiRuntimeConfigCwdFallbackWhenNoInstallRecord`: no install record → cwd fallback still works (non-breaking). Follows the existing `TestPiRuntimeConfigDefaultsIntercomAndAuthPathsUnderHome` pattern.
- `TestPiInstallAcceptsPluginDirAndRecordsPath`: install with `--plugin-dir <temp>` writes that path to the record. Supersedes `TestPiInstallRejectsPluginDir` (which asserts the current rejection).
- `TestPiInstallRecordsRepoPath`: install from a temp repo dir writes the resolved path; record file exists and contains the absolute path.

**`pi-live` lane (AC-1 live confirmation):**
- `spacedock doctor --host pi` from a non-repo cwd reports the install-recorded path in the "Spacedock skills:" line (the live observable).

**Non-regression (AC-3):**
- Existing `internal/cli/frontdoor_test.go` tests stay green — no changes to claude/codex launch paths.

## Riskiest mechanism — no spike needed

The mechanism is straightforward Go file I/O (write a path file at install, read it at config resolution) and a priority-order change in `piRuntimeConfigFromEnv`. No unverified runtime mechanism, no format round-trip, no external tool handoff. The riskiest path is the install-write itself — but `os.WriteFile` of a single path string is proven. Record: no spike needed — the mechanism relies on proven stdlib file I/O and an existing priority-resolution pattern already in the function.

## Related

- `pi-ensign-skill-injection` (`k8tbnmcbyqc5kkhj0m9vewq4`) — child-side skill discovery; composes with this task.
- `internal/cli/pi.go` (`piRuntimeConfigFromEnv`, `runPi`, `runInitWithPi`) — the resolution + install source of truth.
- `internal/cli/init.go` — claude/codex install (for parity understanding; NOT modified).
- `internal/cli/frontdoor.go` — claude/codex launch (for parity understanding; NOT modified).
- `0223-pi-dispatch-contract` sprint index.

## Related

- `pi-ensign-skill-injection` (`k8tbnmcbyqc5kkhj0m9vewq4`) — child-side skill discovery; composes with this task.
- `internal/cli/pi.go` (`piRuntimeConfigFromEnv`, `runPi`) — the resolution source of truth.
- `0223-pi-dispatch-contract` sprint index.

## Stage Report: ideation

- DONE: Confirmed cwd-fallback root cause against live `internal/cli/pi.go` `piRuntimeConfigFromEnv` (line 214): resolution priority is `--plugin-dir` → `SPACEDOCK_REPO_ROOT` → `dir` (cwd). The cwd fallback feeds `cfg.firstOfficer` and `cfg.ensign` which are passed as `--skill` flags to `pi` (pi.go:87-89).
- DONE: Parity analysis — claude/codex do NOT have this problem. They delegate skill discovery to the host's native plugin system (marketplace install + `--plugin-dir` passthrough). Pi is unique: no native plugin mechanism, so `spacedock pi` must resolve the repo path itself. `spacedock install --host pi` is check-only (no `ops.Install` call, no file writes) and rejects `--plugin-dir` (pi.go:199). Fix is pi-only (AC-3).
- DONE: Mechanism decision — (a) install records the repo path durably. New priority: `--plugin-dir` → `SPACEDOCK_REPO_ROOT` → install-recorded path (NEW) → cwd (demoted). Non-breaking; cwd stays as last resort. Install also accepts `--plugin-dir` to specify which repo to record (lifts the current rejection at pi.go:199).
- DONE: Finalized ACs (AC-1 through AC-4) — all behavior-bound, verified by Go unit tests + pi-live lane confirmation. No prose-grep.
- DONE: Test plan — 4 Go unit tests following existing patterns (`TestPiRuntimeConfigResolvesEnvPaths...`, `TestPiRuntimeConfigDefaults...`), pi-live lane doctor check, non-regression on frontdoor tests.
- DONE: Riskiest mechanism — no spike needed (proven stdlib file I/O + existing priority-resolution pattern).

### Summary

Ideation complete. Root cause confirmed in live code. Parity analysis determines the fix is pi-only (claude/codex use host-native plugin discovery). Mechanism (a) picked: install records the repo path durably, launch reads it ahead of cwd fallback, non-breaking. Four finalized ACs, four Go unit tests + pi-live confirmation. No product files edited — design only in entity body.


## Superseded (2026-06-19)

SUPERSEDED by the merged task `pi-install-managed-skill-placement` (filed 2026-06-19). Captain review uncovered that both this task's mechanism (the ` .pi/skills/ensign` repo symlink) and `pi-launcher-repo-resolution`'s (the cwd-fallback demotion + install-record file) are clone-bound workarounds for the fact that `spacedock install --host pi` is check-only (writes nothing). The correct mechanism — verified against the `obra/superpowers` reference and pi-subagents source — is install-managed package placement: ship Spacedock as a pi package (`package.json` with `pi.extensions` + `pi.skills`), make `spacedock install --host pi` run `pi install git:github.com/spacedock-dev/spacedock`, and let both the parent (via the extension's `resources_discover`) and pi-subagents children (via the package-root scan of `package.json -> pi.skills`) discover the skills. No clone, no cwd, no symlink. The merged task absorbs both this task's and `pi-launcher-repo-resolution`'s scope; the staff-review gap-1 `cwd:<repo>` wiring in the capstone becomes unnecessary once the install mechanism lands (child no longer cwd-keyed). Archived REJECTED as superseded — no deliverable merged (ideation-only).
