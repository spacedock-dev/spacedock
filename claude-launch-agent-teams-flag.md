---
title: spacedock claude exports CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1 so FO sessions get the worker back-channel
status: implementation
sprint: 0221-layered-fo
group: binary-ux
id: 662sh1n92mkf33rzgwxy8zcd
worktree: .worktrees/spacedock-ensign-claude-launch-agent-teams-flag
started: 2026-06-20T04:47:59Z
---

The worker↔FO back-channel (`SendMessage`/`TeamCreate`) is gated behind the experimental env flag `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS`, OFF by default. A `spacedock claude` FO session launched without it can spawn NAMED background agents (the `Agent` tool) but cannot ADDRESS them — `SendMessage` reports "exists but is not enabled in this context" — so reuse-advance, mid-run steering, and cooperative shutdown silently do not work, and the FO contract's back-channel assumption is false. Confirmed this session (Claude Code 2.1.183, flag unset): `SendMessage` was absent until the flag was set; transcript forensics show same-`agentSetting` FO sessions split purely on the flag (flag-on used `SendMessage` 20–55×, flag-off zero). "named ≠ addressable."

## Fix (test-only-on-launcher, root-cause)
`spacedock claude` exports `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` into the environment of the `claude` child process it launches (the launcher's claude exec path — `internal/cli`), UNLESS the parent env already sets it (respect an explicit operator override either way). Makes the back-channel reliable by construction rather than ambient-env-dependent. Scope: the `spacedock claude` launch path only; do NOT change `spacedock codex` / `spacedock pi`.

## Acceptance criteria
- **AC-1** — `spacedock claude` sets `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` in the launched `claude` process environment when the parent env does not set it. Verified by a launcher test asserting the child env (the exec/command env the launcher builds) carries the flag.
- **AC-2** — an explicit parent value is PRESERVED, not overridden (e.g. parent `=0` stays `=0`). Verified by a test exercising the parent-set path.
- **AC-3** — scoped + green: `go build ./...` and the launcher test package pass; `git diff` touches only the `spacedock claude` launch path (no codex/pi launch change, no contract/skills change).
