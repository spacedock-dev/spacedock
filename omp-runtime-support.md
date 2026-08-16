---
title: Add native Oh My Pi runtime support
status: backlog
source: "Captain request, 2026-08-14: OMP launch and skill-loading compatibility"
started:
completed:
verdict:
score:
worktree:
issue:
id: 4xd9bxqek1m4yhtwr4m5pd9z
---

Add a first-class Spacedock runtime lane for Oh My Pi (OMP), so `spacedock omp` can launch an OMP first officer with the correct plugin/skill loading and native worker lifecycle.

## Problem

Spacedock currently recognizes Claude, Codex, and Pi, but not OMP. `spacedock omp` exits 2 with `unknown command: omp`. The current Pi lane cannot be treated as an alias: it launches `pi`, installs with `pi install`, passes Pi-specific `--skill` flags, bootstraps through the Pi extension API, and binds ensign dispatch to `pi-subagents`.

OMP exposes `OMPCODE=1`, a plugin manager, skills under plugin `skills/<name>/SKILL.md`, persistent local linking through `omp plugin link`, and per-invocation development loading through `omp --plugin-dir <path>`. Its native worker surface is OMP's `task`/agent mechanism, which must be probed and bound rather than assumed equivalent to Pi.

The current repository package has a `package.json.pi` manifest and `.pi/extensions/spacedock.ts`; OMP documents `package.json.pi` fallback support, but extension API and worker behavior still require runtime evidence.

## Proposed approach

Add an explicit `omp` launch/setup lane after a focused runtime probe. Detect `OMPCODE` as the OMP host marker and prevent inherited `CLAUDECODE=1` from silently misclassifying an OMP session; preserve ambiguity safety when both markers genuinely identify active hosts. Launch `omp` with its native `--plugin-dir`/plugin discovery contract instead of Pi's `--skill` argv. Add OMP install and doctor checks using `omp plugin install` or `omp plugin link` semantics as appropriate.

Add first-officer and ensign runtime bindings for OMP's actual extension, skill, task-agent, completion, messaging, and shutdown surfaces. Keep shared workflow capabilities host-neutral. Prove the path with fixture-backed CLI tests and a live-gated split-root smoke that verifies process exit, skill loading, one native worker dispatch, durable entity mutation, path-scoped commit evidence, and clean state.

## Out of scope

- Modifying OMP itself or its upstream plugin manager.
- Making OMP a silent alias of the Pi launcher.
- Copying skills into undocumented OMP directories instead of using its plugin discovery contract.
- Claiming runtime support from environment-marker or help-text inspection alone.
- Reworking Claude, Codex, or existing Pi behavior except for shared detector tests required to preserve their contracts.

## Expected surface and tolerance

Estimate net LOC change: +500, across 15 files.

## Acceptance criteria

**AC-1 — `spacedock omp` launches OMP with the Spacedock package and task prompt under the native OMP loading contract.**
Verified by: CLI launch test capturing exact argv, including `omp` and the supported local `--plugin-dir`/resource arguments, plus a live smoke process that starts and exits successfully.

**AC-2 — OMP sessions resolve as OMP even when `OMPCODE=1` is present, and inherited Claude markers cannot silently claim the host.**
Verified by: runtimehost table tests covering OMP-only, Claude-only, both-marker ambiguity, and no-marker environments; `spacedock --version` output and dispatch policy match each case.

**AC-3 — A local Spacedock checkout is discoverable by OMP without copying skill files.**
Verified by: isolated-home fixture using `omp --plugin-dir` or `omp plugin link`, asserting discovery of `skills/first-officer/SKILL.md` and `skills/ensign/SKILL.md` from the checkout.

**AC-4 — OMP installation and doctor checks use OMP's real plugin state and report actionable readiness.**
Verified by: isolated plugin-root fixture or live-gated command output for install/link, package manifest, skill discovery, extension readiness, and auth prerequisites; no Pi-specific install command is executed.

**AC-5 — OMP first-officer dispatch uses the native worker mechanism and verifies durable completion.**
Verified by: runtime adapter tests for task-agent payload and completion routing, plus a live-gated split-root smoke that proves one worker report, entity content, state-checkout commit, and clean entity path.

**AC-6 — Existing host lanes retain their current contracts.**
Verified by: focused Claude, Codex, and Pi launch/setup/runtime tests and the repository's required Go test and race commands.

## Test plan

First probe OMP's installed CLI and source-level documentation for plugin manifest fallback, `--plugin-dir`, skill discovery, task-agent payloads, completion signals, and shutdown. Add focused unit tests before implementation for command routing, argv construction, marker precedence, setup checks, and native dispatch shape. Add an isolated fixture test for plugin/skill discovery. Finish with a live-gated OMP split-root smoke using copied credentials and temporary OMP state, then run `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`.

