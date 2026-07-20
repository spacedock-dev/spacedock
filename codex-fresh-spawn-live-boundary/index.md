---
title: Enforce Codex fresh-spawn isolation at the live FO boundary
status: backlog
source: "Captain-directed v0.25.2 follow-up, 2026-07-20; live FO escape after v0.25.1 / archived rt8 / PR #532"
score: "1.0"
milestone: 0.25.2
started:
completed:
verdict:
worktree:
issue:
id: 6cc3rvfd44y6x3352hh21v8b
---

v0.25.1 hardened `internal/dispatch/codex_v2_adapter.go` so its generated map always carries `fork_turns: "none"`, and aligned the Codex FO prose. That did not close the live boundary: in a 2026-07-20 first-officer session, the FO directly invoked `spawn_agent` with `fork_turns: "all"`. The call was rejected only because it also carried an incompatible explicit agent type; the Spacedock architecture itself did not make inherited-turn spawning impossible.

This is the escaped value defect from archived `codex-fresh-dispatch-context-isolation` (`rt8`, PR #532), not a reopening of its historical release record. Its first validation correctly reported that adapter output and a one-off host probe did not prove the instruction-driven FO invocation. The task later narrowed the acceptance boundary to the adapter map and shipped v0.25.1, leaving the original live claim unenforced.

The 0.25.2 fix must bind the actual worker-creation seam used by an assumed Spacedock FO. In current architecture, `fork_turns` is never a selectable continuity mechanism: every fresh spawn is isolated, and deliberate continuity exclusively uses `followup_task` on the existing handle. `"all"`, numeric forks, omission that defaults to full history, and helper/runtime override channels are invalid Spacedock behavior.

Ideation must identify the smallest executable enforcement seam. If Spacedock cannot enforce this at the live tool boundary, stop with a proven upstream/runtime blocker rather than relabeling adapter or prose evidence as the fix. Coordinate with `per-host-stage-model-override` (`e3g`) so model/effort work cannot restore or conditionalize the isolation invariant.

## Acceptance criteria

**AC-1 (VALUE) — An instruction-driven Codex First Officer cannot issue a fresh Spacedock worker spawn that inherits parent turns.** Verified by an integrated live observation joining the generated dispatch artifact, the exact FO-issued `spawn_agent` arguments, and child-visible context. A seeded temptation to use `"all"`, a numeric fork, or omission must still produce exact `fork_turns: "none"`; a mutation/baseline that permits inherited context must turn the proof red.

**AC-2 — Deliberate continuity remains exclusively `followup_task` on the existing worker handle.** Verified in the same live journey: a fresh child lacks the parent canary, then a follow-up reaches that exact child and recalls its own prior-turn marker.

**AC-3 — Adapter-only maps, contract wording, and one-off raw host probes are supporting evidence, never substitutes for AC-1.** The original live value criterion may not be narrowed after rejection. If no executable Spacedock seam can own AC-1, the gate is REJECTED or BLOCKED with the upstream dependency named; it is not PASSED.

**AC-4 — v0.25.2 ships the fix on the stable line without rewinding `next`.** Verified by the exact release candidate SHA: required Go/full/race gates, integrated Codex live evidence for AC-1/2, annotated `v0.25.2` from `main`, and durable proof that `next` retains the invariant through the documented propagation path.

## Scope guard

Do not build a new general agent harness or fork-mode configuration surface. Reuse the narrowest existing live capture path that can observe the real FO tool call. Do not expose `fork_turns` to workflow authors. Do not change stage reuse policy. Do not fold unrelated model/effort routing into this patch.
