---
id: vva363ty3afdgvxcacrfx6mc
title: Keep the spacedock launcher resident — spawn the host instead of exec'ing away
status: ideation
source: captain (2026-06-10) — zellij restart legibility today, sidecar launches later
started: 2026-06-29T22:28:07Z
completed:
verdict:
score:
worktree:
issue:
---

`spacedock claude` currently replaces itself with the host process via `syscall.Exec` (`internal/cli/host_exec.go:227`), so the launcher vanishes from the process tree the moment the host starts. The captain wants the launcher to stay resident as the parent process for two reasons:

1. **Future sidecars** — a resident launcher can spawn companion processes (sidecar services, watchers) alongside the host later; an exec'd-away launcher cannot.
2. **Zellij legibility now** — with exec, zellij's restart view shows the raw underlying command (`Waiting to run: bash …/safehouse --trust-workdir-config -- claude --dangerously-skip-permissions --agent spacedock:first-officer --plugin-dir … --resume <session-id>`) instead of a legible `spacedock claude …` line, making session restarts hard to read.

The change: spawn the host as a child (wait, forward signals, propagate exit code, pass through the TTY) instead of `syscall.Exec`.

## Problem

{What is broken or missing, and why it matters. Ideation fills this in — seed context above.}

## Proposed approach

{How the task intends to solve the problem. Ideation fills this in. Candidate shape: replace the `syscall.Exec` call site with `exec.Command` + `cmd.Wait`, with signal forwarding (SIGINT/SIGTERM/SIGWINCH), exit-code propagation, and stdin/stdout/stderr TTY passthrough. Consider whether all hosts (claude/codex/pi) or only `spacedock claude` change, and Windows semantics where exec already emulates spawn-wait.}

## Out of scope

{Ideation confirms. Likely: actual sidecar features — this task only establishes the resident-parent process model that makes them possible later.}

## Acceptance criteria

Seed sketches — ideation must firm each with an end-state property and an independently checkable "Verified by".

**AC-1 (seed) - While the host runs, the spacedock launcher remains its parent process.**
Verified by: {ideation — e.g. a behavior test launching a stub host and asserting the process tree / PPID.}

**AC-2 (seed) - The launcher exits with the host's exit code.**
Verified by: {ideation — stub host exiting non-zero; assert launcher exit code matches.}

**AC-3 (seed) - Terminal signals reach the host (Ctrl-C interrupts the host, not just the launcher) and interactive TTY behavior is preserved.**
Verified by: {ideation — signal-forwarding test against a stub host.}

## Test plan

{Ideation fills this in. Likely Go unit/behavior tests around the launch path with a stub host binary; a manual zellij observation is supporting evidence, not the proof.}
