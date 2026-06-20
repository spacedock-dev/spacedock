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

## Folded-in contract-hygiene findings (captain review 2026-06-20)

These are the token reduction the binary-migration program OWES — the END (a measurably leaner contract), not deferrable polish. AC-1's "cumulative line delta NEGATIVE" is the right gate; hold every item below to it. All are contract-file (scaffolding) edits → dispatched worker, contractlint-guarded; none behavioral.

- **Decouple step-number coupling + add the missing capabilities (PR #414 #4/#5).** `claude-first-officer-runtime.md:13`, `claude-fo-dispatch.md:159 & 171`, and `fo-merge-core.md` couple teardown to "Merge-and-Cleanup step 10" — banned by `docs/runtime-support.md:13`. Root cause: the core (post-#409) defines six capability `«fn»`s but NO `«worker.shutdown»` (nor `«worker.spawn»`), though `runtime-support.md:9` names them. Add `«worker.shutdown»` (and `«worker.spawn»`) as core capability `«fn»`s with per-host `→` bindings, then replace the step-10 references with the capability name.
- **Evergreen the contract — strip temporal markers.** The operating contract must carry NO sprint/version/roadmap temporality: `→ **shipped** (this sprint):` (×4 — `fo-merge-core.md:14`, `first-officer-shared-core.md:166/174/181`) → `→ **shipped**:`; "tracked for 0.20.6" (`claude-fo-dispatch.md:159`); "(0222) … descoped to roadmap 0222" (`fo-dispatch-core.md:150`). These annotations encode means-thinking (what shipped this sprint), not value.
- **Name capabilities concretely — drop generic "the verb".** "the verb" (×7 — shared-core ×3, fo-merge-core ×4) is confusing; say `«merge.guard»` / `spacedock merge guard` (or the specific capability).
- **#6 — design decision, NOT an auto-cut.** Host tool bindings in the core's `→` lines are the contractlint-guarded transition form `runtime-support.md:15` permits; the stricter shape (distribute bindings to per-runtime files) is optional. Record the decision before any cut.

Root insight (see entity `p2` / `pr-complete-binary-command`): the contract GREW (`fo-merge-core.md` 41→49 this sprint) because `merge guard` is a THIN finalize-helper — the prose carries "when X signals, the FO does Y" choreography on top of the still-FO-owned orchestration. The verb is the MEANS; the leaner contract is the END, and it was never harvested. The real collapse is the FAT `spacedock pr complete` (p2) absorbing the orchestration so the choreography prose disappears.
