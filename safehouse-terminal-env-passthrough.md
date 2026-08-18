---
id: 89bdt7vhk362yw630kt9s12k
title: Safehouse passes terminal-host env for Zellij only, and passes names that do not exist
status: backlog
source: "Captain CL, 2026-08-18, from a live observation: `env | grep TMUX` returned nothing inside a safehouse-wrapped session. Traced to `internal/safehouse/safehouse.go:68-73`. Captain's direction: take the env set from subspace's own probe, and add a name to `--env-pass` only when that variable exists."
started:
completed:
verdict:
score:
worktree:
issue:
---

A sandboxed session can detect Zellij and nothing else. Five other terminal hosts lose their identity at the sandbox boundary, and the allowance lists variable names whether or not they exist.

## Problem

`terminalTargetingEnvArgs` (`internal/safehouse/safehouse.go:68-73`) is the whole allowance:

```go
func terminalTargetingEnvArgs() []string {
	if _, present := os.LookupEnv("ZELLIJ"); !present {
		return nil
	}
	return []string{"--env-pass=ZELLIJ,ZELLIJ_PANE_ID,ZELLIJ_SESSION_NAME"}
}
```

Two separate faults.

**Fault 1 — the host set is Zellij-only.** When `ZELLIJ` is absent the function returns `nil`, so nothing is passed. A tmux session therefore reaches the child with no `TMUX` and no `TMUX_PANE`. Confirmed live by the captain: `env | grep TMUX` is empty inside a safehouse-wrapped session.

The frontdoor is not the cause. `launchEnv` (`internal/cli/frontdoor.go:73-79`) forwards the entire parent environment and only swaps `SPACEDOCK_BIN`. The loss is at the sandbox boundary.

**Fault 2 — presence is checked per host, not per variable.** The gate reads one sentinel (`ZELLIJ`), then passes three names unconditionally. `ZELLIJ_PANE_ID` and `ZELLIJ_SESSION_NAME` are listed whether or not they are set. The gate variable is also not one of the two the consumer actually probes, so the sentinel and the payload disagree about what identifies a Zellij pane.

## The consumer already defines the correct set

The `subspace:r` skill documents a probe of nine signals across six hosts (`plugins/subspace/skills/r/SKILL.md:71-91`), in its own resolution order:

| Host | Signals |
|---|---|
| Zellij | `ZELLIJ_SESSION_NAME`, `ZELLIJ_PANE_ID` |
| tmux | `TMUX`, `TMUX_PANE` |
| Herdr | `HERDR_ENV`, `HERDR_PANE_ID` |
| CMUX | `CMUX_WORKSPACE_ID`, `CMUX_SURFACE_ID` |
| Ghostty | `TERM_PROGRAM=ghostty` |
| Apple Terminal | `TERM_PROGRAM=Apple_Terminal` |

Safehouse passes three names. One of them, `ZELLIJ`, is not in this probe at all. Six probed signals never cross the boundary.

The consequence is concrete: inside safehouse, `subspace:r` resolves Zellij or falls through to none. A captain in tmux, Herdr, CMUX, Ghostty, or Apple Terminal gets the fallback rather than their real pane host.

Note that the probe treats an empty value as absent. So a name passed for an unset variable is not merely useless — it risks presenting an empty variable where the consumer expects nothing.

## Proposed approach

{Ideation fills this in. Two changes: widen the set to the consumer's nine signals, and build `--env-pass` from the variables that are actually present rather than from a fixed list gated on one sentinel. Decide where the list should live — duplicating a sibling repository's probe invites drift, so name how the two stay agreed, or why duplication is acceptable here.}

## Out of scope

Changing `launchEnv` or the frontdoor. Owning an operator configuration surface for env passthrough; the existing comment states that safehouse deliberately adds a default without parsing caller arguments, and that boundary stays.

## Expected surface and tolerance

Estimate net LOC change: +40 across 2 files. Report insertions and deletions separately. Do not declare a gross tolerance. Semantics changed: the `--env-pass` argv safehouse composes.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - A terminal host's identifying variables survive the sandbox boundary for every host the consumer probes.**
This is the measuring AC: the count of probed signals that reach the child, out of the nine, must equal the count set in the parent. Verified by running the consumer's own probe script inside a safehouse-wrapped session under at least tmux and Zellij, and comparing present/absent for each signal against the same probe run outside the sandbox. Fails on the current code, where a tmux parent yields `TMUX=absent` in the child.

**AC-2 - No variable name appears in `--env-pass` unless that variable is set in the parent.**
Verified by a unit test over the argv composer with a stubbed environment: a parent holding only `TMUX` and `TMUX_PANE` produces an allowance naming those two and nothing else; an empty parent produces no allowance at all. Fails if the composer emits a fixed list, or emits names gated on a different variable than the one being passed.
