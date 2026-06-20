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

## Stage Report: implementation

- DONE: Export `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` into the claude child env when the parent does not set it; preserve an explicit parent value (even =0).
  `internal/cli/frontdoor.go` runClaude launch site now wraps `withAgentTeams(launchEnv(os.Environ()))`; `withAgentTeams` appends `=1` only when `hasEnv` is false. Codex/pi launch sites unchanged. Commit 330ae64b.
- DONE: TDD — failing launcher test FIRST asserting the unset-parent and set-parent paths; then implement to green.
  Added `TestClaudeFrontDoorEnablesAgentTeamsWhenParentUnset` (AC-1), `TestClaudeFrontDoorPreservesExplicitAgentTeams` (AC-2), plus a codex scope-guard test; confirmed compile-fail (undefined `agentTeamsEnv`) before adding the constant/helper, then 3/3 PASS.
- DONE: Scoped + green: `go build ./...` and the launcher test package pass; diff touches ONLY the `spacedock claude` launch path.
  `go build ./...` clean; `go test ./internal/cli/` ok (41.7s); `git diff origin/main...HEAD` = only `frontdoor.go` + `frontdoor_test.go` (no codex/pi launch change, no contract/skills/references change).

### Summary
`spacedock claude` now enables claude's worker↔FO back-channel by construction: the launcher injects `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` into the exec'd claude child env unless the parent already set it (explicit operator value, even =0, is preserved). The fix is a 2-function helper (`withAgentTeams`/`hasEnv`) applied at the runClaude `ops.Launch` site only; the shared `launchEnv` and the codex/pi launch paths are untouched. Test-driven: tests written and confirmed red (compile-fail) before the constant/helper existed, then green.

## Stage Report: validation

- DONE: Reproduce AC-1/AC-2 from the worktree (`go test ./internal/cli/ -run AgentTeams -v`); launcher injects the flag when parent unset, preserves explicit parent value, codex does NOT inject.
  3/3 PASS: `TestClaudeFrontDoorEnablesAgentTeamsWhenParentUnset` (AC-1), `TestClaudeFrontDoorPreservesExplicitAgentTeams` (AC-2, parent=0 stays 0), `TestCodexFrontDoorDoesNotEnableAgentTeams` (scope guard).
- DONE: Independently confirm the AC-1 test is load-bearing (not vacuous), per the detached-adversarial-audit discipline.
  Mutation test: reverted the production line to plain `launchEnv(os.Environ())`; AC-1 went RED (env shows flag absent), AC-2/codex stayed green; restored. Test genuinely depends on the impl.
- DONE: AC-3 scope + green — `go build ./...`, `go test ./internal/cli/`, diff scope.
  Build clean; `internal/cli` test package passes; `git diff --name-only origin/main...HEAD` = exactly `internal/cli/frontdoor.go` + `internal/cli/frontdoor_test.go`. No codex/pi launch change (codex site line 557 unwrapped; `runPi` in pi.go untouched), no contract/skills/references.
- FAILED: LIGHT LOCAL LIVE E2E — prove the env flag actually enables the back-channel end-to-end (before=absent / after=present).
  The flag does NOT discriminate back-channel behavior in Claude Code 2.1.183 (the exact version the entity cites). Built launcher from worktree; `printenv CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` in the launched child = `1` (injection works), but the required before/after contrast collapses: SendMessage is callable and the spawn→address→reply round-trip succeeds with the flag UNSET (plain `claude -p`) AND with it forced `=0` in `--agent spacedock:first-officer` context — `ADDRESSABLE: yes` in every state. No "exists but not enabled in this context" gating observed.

### Summary
RECOMMENDATION: REJECTED. The code change is correct and well-scoped — AC-1/AC-2/AC-3 unit-level pass, the AC-1 test is load-bearing under mutation, the launcher provably injects `=1` into the live child env, and the diff is exactly the two intended files. But the entity's premise is falsified by live evidence on the cited version (Claude Code 2.1.183): the worker↔FO back-channel works WITHOUT the flag. With the flag unset, and even with it explicitly `=0` in FO context, I spawned a named background agent and addressed it via `SendMessage`, receiving a reply — the full round-trip. The required E2E proof ("the env flag actually enables the back-channel end-to-end, not just that the env map carries it") therefore cannot be satisfied: there is no before/after delta to observe. This is the validation-stage reject criterion "tests pass but prove an obsolete/wrong target behavior" — the unit tests verify env-map plumbing that is now a runtime no-op for the stated goal. Recommend the captain decide whether to (a) drop the change as no-longer-needed, (b) re-establish the gating premise (identify the exact session shape where flag-off truly disables addressing — the entity's original forensics may predate a Claude Code behavior change), or (c) keep the harmless belt-and-suspenders injection but rewrite the ACs to claim only "flag is set in child env," not "enables the back-channel."

## Feedback Cycles

### Cycle 1 — validation REJECTED → rewrite ACs (captain-directed)
Validation found (live e2e) that `env -u CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS claude -p` STILL exposes SendMessage — the env var is not the runtime gate on this machine. The real enabler is `~/.claude/settings.json` (`env.CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS: "1"` + `teammateMode: "auto"`), which re-applies the flag to every child regardless of shell env. So the launcher's env-export does not independently prove "enables the back-channel." Captain's call: KEEP the harmless export, REWRITE the ACs to claim only what the code does.

Rework:
1. Rewrite `## Acceptance criteria` to claim ONLY: the launcher sets `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` in the claude child env when the parent is unset, preserves an explicit parent value, scoped to the claude launch path. DROP the "enables the back-channel end-to-end" / live-e2e claim.
2. Add a short note (entity body + one-line code comment near `withAgentTeams`): the AUTHORITATIVE enabler is `~/.claude/settings.json` (env flag + `teammateMode: auto`); the launcher export is best-effort for users without that setting and a no-op when settings already enable it — it does NOT independently enable the feature.
3. Keep `frontdoor.go` + the 3 tests (already green for the narrowed ACs). Confirm `go build ./...` + `go test ./internal/cli/` green; diff still scoped to the claude launch path.
