---
title: Six runtime tokens named in the host adapters have no owner — close or delete them
status: backlog
source: "roborev job 328 finding 4b on entity 8413fc05vp (contractlint runtime-semantics retirement); captain ruling 2026-07-20: delete the prose check and record the gap rather than retain a phrase check for a token nothing exercises."
started:
completed:
verdict:
score:
worktree:
issue:
id: y7deh2nsk5hh3a0zx1mf9j06
---

Six runtime tokens named in the Codex and Pi host adapters have no owner anywhere in the repo — no Go emitter, no build fixture, no live-lane assertion. Their prose checks were deleted (not retained) under the captain ruling of 2026-07-20, because a phrase check for a token nothing exercises proves only that we wrote the word. 0.26.0 ships these gaps knowingly; this entity tracks closing them.

## Problem

The contractlint runtime-semantics retirement routed each runtime-meaning claim to an independent source. Six tokens had no such source. They are recorded inline as `UNCOVERED` annotations on `codexToolTokens` / `piSubstrateNativeTokens` in `internal/contractlint/runtime_binding_block_test.go`; containment still bounds WHERE each name may appear, but nothing proves WHAT it does.

| Token | Adapter | State |
| --- | --- | --- |
| `send_message` | codex-first-officer-runtime.md | no assertion anywhere |
| `list_agents` | codex-first-officer-runtime.md | no assertion anywhere |
| `interrupt_agent` | codex-first-officer-runtime.md | no assertion anywhere |
| `member_spawn` | pi-first-officer-runtime.md | absent from `teams.go` AND all of `internal/ensigncycle` |
| `intercom(` | pi-first-officer-runtime.md | no assertion anywhere |
| `subagent(` | pi-first-officer-runtime.md | appears only as PROMPT text in the pi live lanes, never asserted |

Two agent-facing Pi instructions are also unowned: "non-fresh resume is only an explicit manual/debug exception" and "file verification remains the completion gate" are document-only. ("Fresh redispatch remains the default" IS bound, via `SubagentStageDispatch`'s `context: "fresh"`.)

## Proposed approach

For each token, either add an assertion that records the actual tool call and resulting workflow state, or delete the token from the adapter if Spacedock neither emits nor depends on it. Do NOT re-add a prose grep — that is the failure mode the parent entity removed.

`subagent(` and `intercom(` depend on the pi live lane, which is red for a pre-existing environment reason (`Cannot find module '@earendil-works/pi-coding-agent'`) tracked by entity 268. Sequence after it.

## Acceptance criteria

**AC-1 - Each of the six tokens is either exercised by a named assertion that reds when the behavior regresses, or removed from the adapter.**
Verified by: for each token, a planted divergence (the tool call absent or renamed) reds the named assertion; or the token no longer appears in any adapter.

## Out of scope

Re-adding phrase checks. Changing Codex or Pi runtime policy.
