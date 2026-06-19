---
title: Pi launcher repo resolution — stop resolving the Spacedock repo from cwd; record it at install
status: backlog
source: "Captain (2026-06-19, 0223-pi-dispatch-contract sprint scope-lock): `spacedock pi` resolves repoRoot as --plugin-dir -> SPACEDOCK_REPO_ROOT -> working directory (cwd). The cwd fallback means running the homebrew binary from any directory points `--skill <cwd>/skills/ensign/SKILL.md` wherever you happen to stand — the skill is registered by cwd-luck, not explicit install. Captain: 'it should not use the current directory's skill without me explicitly running pi with --skills.'"
score:
started:
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

## Approach (candidate fixes — ideation confirms and picks)

- **(a) Install records the repo path durably.** `spacedock install --host pi` writes the resolved repo path to a known location (e.g. `~/.pi/agent/spacedock-repo` or a spacedock-managed config), and `spacedock pi` reads it as a resolution source ahead of the cwd fallback. Install-time explicit; survives launch from any cwd.
- **(b) Remove the cwd fallback; require `--plugin-dir` or `SPACEDOCK_REPO_ROOT`.** Strictest — launch from elsewhere fails loudly with a clear "set --plugin-dir or SPACEDOCK_REPO_ROOT" message. May break existing operator muscle memory of running `spacedock pi` from the repo root.
- **(c) Keep the cwd fallback but WARN when `<cwd>/skills/ensign/SKILL.md` is absent.** Softest — preserves convenience but surfaces the cwd-luck. Weakest; a warning is a prose-only gate (banned by proof policy as the sole assurance).

Ideation picks one (recommend (a) — install-managed, explicit, non-breaking), records the decision, and plans verification. Consider how claude/codex hosts resolve the repo (do they have the same cwd fallback? parity matters).

## Acceptance criteria (provisional — ideation finalizes; proof = behavior, never prose-grep)

**AC-1 — `spacedock pi` launched from a cwd that is NOT the Spacedock repo still resolves the correct skill paths.**
Verified by: a live launch from a temp/unrelated cwd (after install) where the `--skill` flags point at the real Spacedock repo, confirmed by doctor output or the resolved config — NOT by prose-grep over the launcher.

**AC-2 — The repo path is recorded at install time, not derived from cwd at launch.**
Verified by: install writes a durable record; launch reads it; the cwd fallback is either removed or demoted below the install record. A behavior test that launches from a non-repo cwd and asserts the skill path resolves to the install-recorded repo.

**AC-3 — Parity with claude/codex repo resolution is explicit.**
Verified by: ideation records whether claude/codex have the same cwd fallback and whether this fix should apply to all hosts or pi-only; the decision is behavior-tested or documented with a concrete reason.

## Out of scope

- Child-side skill discovery (pi-subagents) — that is `pi-ensign-skill-injection`. This task is the launcher's repo-resolution only.
- The back-channel / named-capability hardening — that is `pi-back-channel-dispatch`.
- Model stamping — that is `pi-dispatch-model-stamping`.

## Test plan

- Live `spacedock pi` launch from a non-repo cwd (AC-1) — assert `--skill` paths resolve to the real repo via doctor or config dump.
- Install-record behavior test (AC-2) — install writes the record; a fresh launch from elsewhere reads it. Bounded Go test or behavior fixture.
- Parity check (AC-3) — read claude/codex launcher resolution; record the decision.

## Related

- `pi-ensign-skill-injection` (`k8tbnmcbyqc5kkhj0m9vewq4`) — child-side skill discovery; composes with this task.
- `internal/cli/pi.go` (`piRuntimeConfigFromEnv`, `runPi`) — the resolution source of truth.
- `0223-pi-dispatch-contract` sprint index.
