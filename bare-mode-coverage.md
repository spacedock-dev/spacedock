---
id: 2ehx77vn94s20k8cqs2p8vhw
title: Bare-mode dispatch coverage — explicit forced-bare scenario + `-p`-assumes-bare audit
status: backlog
source: 'Shaping-FO 2026-06-18 — verified live on claude 2.1.181 that a plain headless `claude -p` (no CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS flag, no --teammate-mode, no tmux) exposes the merged team channel flag-free (SendMessage in tools[], Agent background teammate, token round-trip). On unpinned claude `-p` is team-enabled by default, so the old "-p implies bare/sequential" equation is dead and bare loses its incidental coverage at 2.1.177-pin-retirement.'
started:
completed:
verdict:
score:
worktree:
issue:
---

Seed: keep bare mode (the sequential dispatch fallback) provably covered once unpinned claude makes headless `-p` team-enabled by default.

## Problem

Historically headless `claude -p` exposed no team tools, so spacedock's `-p` dispatch fell into bare/sequential mode — and every existing `-p` live scenario (ensign cycle, shared scenarios, default-headless-stops-at-gate) implicitly covered bare because the live-e2e job is pinned to claude 2.1.177. Verified this session on 2.1.181: a plain `-p` (flag-free, no `--teammate-mode`, no tmux) carries the merged team channel — `Agent(name, run_in_background=true)` + `SendMessage`. So when the 2.1.177 pin is retired (gated behind m4's legacy lane per 9243's deprecation trigger), those `-p` scenarios flip from bare to merged, and bare mode — retained per 9243 Change 4 for a genuinely degraded host (Agent/SendMessage absent) or explicit operator command — loses its incidental coverage and risks becoming untested/dead.

## Proposed approach (seed — ideation to flesh out)

- A dedicated, EXPLICITLY-triggered forced-bare live scenario so the sequential fallback stays proven independent of host team-capability.
- An audit of every place that assumes `-p` ⇒ bare (the FO contract headless handling, the existing live shared scenarios, the `internal/dispatch/build.go` bare branch reachability) before the pin is dropped, each with a disposition.

## First ideation question

Is `claude --bare` (claude's "minimal mode: skip hooks, LSP, plugin") the clean forced-bare trigger — i.e. does it ALSO drop the merged team channel (no `SendMessage`/`Agent`-team in `tools[]`)? Verify empirically (a flag-free `-p --bare` probe, inspect `tools[]`) before designing the scenario. If `--bare` does not suppress the team channel, find the trigger that does (stub Agent absent / explicit operator flag / a spacedock-side force).

## Out of scope

- The merged-lane work (m4) and the 9243 contract redefinition (Change 4 — already done).
- The 2.1.177 pin-drop itself (separate pin-retirement work).

## Sequencing

AFTER m4 (legacy + merged lanes) and tied to 2.1.177 pin-retirement (gated by 9243's deprecation trigger: legacy retires only when no live lane drives it). NOT a blocker for m4 — it is a downstream consequence of the merged model becoming the `-p` default.

## Acceptance criteria (provisional — ideation finalizes; proof = behavior, never prose-grep)

**AC-1 — Forced-bare sequential dispatch is proven on a team-capable host.**
Verified by: a live scenario that forces bare mode (team channel suppressed) and observes sequential one-at-a-time dispatch, green on current/unpinned claude.

**AC-2 — No `-p`-assumes-bare assumption survives unaudited into pin-retirement.**
Verified by: an audit artifact naming every `-p` scenario / contract clause that assumed bare/sequential, each with a disposition (re-pinned, converted to merged, or explicitly forced-bare).

## Test plan

A live forced-bare scenario (the AC-1 oracle) plus a read-only audit pass (AC-2). Cost is small if `--bare` is the trigger; larger if a stub-absent-Agent harness is needed. Ideation resolves the trigger first (the first ideation question) before sizing.

## Related

- `m4` live-team-mode-terminal-harness (m40mphxan8phr3t3tp03gk89, PR #390) — owns the legacy + merged lanes; sequenced before this.
- `9243` using-claude-team-merged-model-support — Change 4 redefined bare mode (no-TeamCreate ≠ bare); owns the `build.go` bare branch.
- The 2.1.177 pin keystone (#395) — its retirement is what opens this coverage hole.
