---
title: Nine runtime tokens named in the host adapters have no owner — close or delete them
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

Nine runtime tokens named in the Codex and Pi host adapters have no owner anywhere in the repo — no Go emitter, no build fixture, no live-lane assertion. Their prose checks were deleted (not retained) under the captain ruling of 2026-07-20, because a phrase check for a token nothing exercises proves only that we wrote the word. 0.26.0 ships these gaps knowingly; this entity tracks closing them.

## Problem

The contractlint runtime-semantics retirement routed each runtime-meaning claim to an independent source. Six tokens had no such source. They are recorded inline as `UNCOVERED` annotations on `codexToolTokens` / `piSubstrateNativeTokens` in `internal/contractlint/runtime_binding_block_test.go`; containment still bounds WHERE each name may appear, but nothing proves WHAT it does.

Set corrected at validation cycle 1 by a mechanical enumeration (three prior hand-audits each missed a different token). NINE tokens have no owner:

| Token | Adapter | State |
| --- | --- | --- |
| `send_message` | codex-first-officer-runtime.md | no Go string literal names it |
| `list_agents` | codex-first-officer-runtime.md | no Go string literal names it |
| `interrupt_agent` | codex-first-officer-runtime.md | no Go string literal names it |
| `timeout_ms` | codex-first-officer-runtime.md | arg of `wait_agent(timeout_ms)`; the tool is owned, the arg is not |
| `path_prefix` | codex-first-officer-runtime.md | arg of `list_agents(path_prefix?)`; neither is owned |
| `member_spawn` | pi-first-officer-runtime.md | absent from `teams.go` AND all of `internal/ensigncycle` |
| `cwd: <resolved repo root>` | pi-first-officer-runtime.md | no Go string literal names it |
| `contact_supervisor` | pi-ensign-runtime.md | no Go string literal names it |
| `need_decision` | pi-ensign-runtime.md | the `reason:` value paired with `contact_supervisor` |

`subagent` and `intercom` are NAMED by Go source but not asserted by it — the hits are the `pi-subagents`/`pi-intercom` package names, the unrelated `subagent_type` field, and dispatch prose. Named is not owned, but they are a weaker case than the nine above.

## Unowned retired SEMANTIC claims (second class, not yet enumerated)

The token table above is complete. A SECOND class is not: retired *sentences* with no owner. Roughly a dozen exist; enumerating them is part of this entity's work. Known members:

- "Pi's model-space binding is provider/model strings" — `TestBuildPiHostIgnoresModelWithNote` proves a claude-enum model is dropped with a note, not what Pi's model space IS.
- "file verification remains the completion gate" and "non-fresh resume is only an explicit manual/debug exception" — document-only instructions. ("Fresh redispatch remains the default" IS bound, via `SubagentStageDispatch`'s `context: "fresh"`.)

## Naming and AC residue carried from the parent

The parent's bindings are doc-to-Go-DECLARATION agreement, not wire evidence: neither `CodexMultiAgentV2Spawn.ToolArgs()` nor any of the five `piruntime` constructors has a production caller. The code comments now say so, but three identifiers (`piEmittedRuntimeToken`, `piEmittedRuntimeTokens`, `TestPiEmittedRuntimeTokensBindGoSource`) and the parent's AC-2/approach prose still say "emits". Correct the naming and the AC wording here.

## Proposed approach

For each token, either add an assertion that records the actual tool call and resulting workflow state, or delete the token from the adapter if Spacedock neither emits nor depends on it. Do NOT re-add a prose grep — that is the failure mode the parent entity removed.

`subagent(` and `intercom(` depend on the pi live lane, which is red for a pre-existing environment reason (`Cannot find module '@earendil-works/pi-coding-agent'`) tracked by entity 268. Sequence after it.

## Acceptance criteria

**AC-1 - Each of the six tokens is either exercised by a named assertion that reds when the behavior regresses, or removed from the adapter.**
Verified by: for each token, a planted divergence (the tool call absent or renamed) reds the named assertion; or the token no longer appears in any adapter.

## Out of scope

Re-adding phrase checks. Changing Codex or Pi runtime policy.
