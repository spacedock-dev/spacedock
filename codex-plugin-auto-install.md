---
id: z9pqhvtj21smtxka8p88j23r
title: Codex plugin auto-install on `spacedock codex` (mirror #311; codex 0.137.0 supports marketplace + plugin add)
status: ideation
source: "captain (2026-06-08) — codex-cli 0.137.0 adds `codex plugin marketplace add` + `codex plugin add` (0.132.0 lacked them). The front-door auto-install is currently Claude-only (frontdoor.go:315 'codex has nothing to auto-install'); now the Codex binary-first path can ensure the plugin exists too, like #311 did for Claude."
score: "0.32"
started: 2026-06-08T15:48:51Z
completed:
verdict:
worktree:
issue:
sprint: 0198-pre-flip-hardening
group: binary-ux
sprint-readiness: ready
---

Bring the front-door plugin auto-install to the Codex binary-first path, now that codex-cli 0.137.0 supports adding a marketplace + plugin non-interactively.

## Problem

- `spacedock claude` already auto-installs a missing plugin (`NoPluginFound → ops.Install → launch`, frontdoor.go:177-185 — the #311 pattern).
- `spacedock codex` does NOT: `frontdoor.go:315` records the auto-install as Claude-only ("codex has nothing to auto-install"), because earlier codex (0.132.0) could not add a marketplace/plugin from the CLI.
- codex-cli **0.137.0** now has `codex plugin marketplace add` / `plugin add` / `list` — so the Codex path can ensure the plugin exists on first launch, the same way Claude does.

## Proposed approach (ideation firms)

Extend `execHost.Install` (host_exec.go:271, currently Claude-only) + `installArgvSequence` to a Codex sequence (`codex plugin marketplace add <source>` + `codex plugin add spacedock@spacedock`), and wire the codex path's `NoPluginFound → auto-install → launch` mirroring the Claude branch. Respect `--no-install`. Reuse the existing `codexEntryInstalled` listing check for idempotency.

## Riskiest unknown (spike first)

Whether the 0.137.0 commands run NON-INTERACTIVELY (no prompt) for an unattended auto-install, and the exact marketplace+add sequence + flags (source/ref). codex 0.137.0 has the commands (confirmed); the spike confirms the non-interactive sequence + the front-door wiring end-to-end.

## Acceptance criteria (sketch)

- `spacedock codex` with no plugin installed auto-installs it (marketplace add + plugin add) and launches — verified by a host-fixture/behavior test (fake codex host) asserting the install sequence runs on `NoPluginFound`, plus the existing gate semantics (mismatch still fails fast; `--no-install` opts out).

## Notes

High-stakes: this is the **front-door launcher** — validation gets a detached adversarial audit. Codex analog of `rbjk` (#311). 0198 binary-ux; couples with `qa` (the install/launch journey).
