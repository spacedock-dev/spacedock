---
title: Trim per-adapter dispatch operational prose toward cumulative-negative (preserve await/reuse/guardrail contract)
status: backlog
sprint: 0221-layered-fo
group: binary-ux
id: adk755xqeb4a9dxhhgtjwawh
---

After the b2 capabilities-as-«fn» reframe (pi-back-channel-dispatch), the capability layer is consolidated: the body calls «fn»s and each `→` line carries the per-host realization. The reframe netted −10 vs the reviewed state but is still +162 CUMULATIVE vs origin/main (origin/main had no capability layer + no AC-1 test). The b2 ensign correctly HELD scope rather than gut load-bearing prose to hit a number.

This task is the careful follow-up cut: now that the «fn» `→` lines consolidate the present/absent logic, trim the now-redundant per-adapter OPERATIONAL prose — candidates: `claude-fo-dispatch.md` `## Worker Back-Channel` / `## Spawn Call` machinery, codex `### Foreground wait` / idle-wake observations — WITHOUT gutting load-bearing contract.

## Hard constraint (the SB2-class risk this task must avoid)
Do NOT cut await/reuse/guardrail contract to hit a line count. The `## Awaiting Completion` premature-reap ban, the idle guardrail, and any reuse-advance/await semantics not captured by the terse `→` lines are load-bearing and STAY byte-intact. If a candidate block carries unique contract, KEEP it.

## Acceptance criteria
- **AC-1** — cumulative line delta of the capability-surface files vs origin/main is NEGATIVE (a real net cut), reported in the stage report.
- **AC-2** — the await/reuse/guardrail contract is preserved: `## Awaiting Completion` (premature-reap ban, idle guardrail) byte-intact; no reuse/await semantic dropped. Verified by diffing those sections (unchanged) + the contractlint host-neutrality/guardrail tests green.
- **AC-3** — `go build ./...` + `go test ./internal/contractlint/` green; scoped to the dispatch adapter prose.

Follows b2 (pi-back-channel-dispatch) — operates on the SAME adapter files; sequence AFTER b2 lands to avoid conflict.
