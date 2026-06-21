---
title: Trim remaining dispatch adapter prose into capability bindings
status: backlog
sprint: 0221-layered-fo
group: binary-ux
id: adk755xqeb4a9dxhhgtjwawh
---

After the b2 capabilities-as-«fn» reframe and #418's Codex/Pi binding-block cleanup, the target is clear: per-runtime files should be capability bindings, not lifecycle narration. The shared core owns when each `«fn»` runs; adapters bind how the host realizes it.

This task is the careful follow-up cut for remaining dispatch/adapter prose, especially Claude and shared dispatch residue. Do not redo #418's Codex/Pi binding-block rewrite or the first-class `«worker.spawn»` / `«worker.shutdown»` core-heading work it landed. Instead, collapse remaining adapter narration into `«fn»` binding lines where the semantics are already represented by capability names.

## Hard constraint (the SB2-class risk this task must avoid)
Do NOT cut await/reuse/guardrail contract to hit a line count. The `## Awaiting Completion` premature-reap ban, the idle guardrail, and any reuse-advance/await semantics not captured by the terse `→` lines are load-bearing and STAY byte-intact. If a candidate block carries unique contract, KEEP it.

## Acceptance criteria
- **AC-1** — cumulative line delta of the remaining capability-surface files vs current main is NEGATIVE (a real net cut), reported in the stage report.
- **AC-2** — the await/reuse/guardrail contract is preserved: `## Awaiting Completion` (premature-reap ban, idle guardrail) byte-intact; no reuse/await semantic dropped. Verified by diffing those sections (unchanged) + the contractlint host-neutrality/guardrail tests green.
- **AC-3** — adapter/core text follows the pseudo-code contract principle in `docs/runtime-support.md`: capability names and binding lines first; prose only for fuzzy judgment, host quirks proven by probes, or guardrails that do not fit one capability.
- **AC-4** — `go build ./...` + `go test ./internal/contractlint/` green; scoped to the dispatch adapter prose.

Follows b2 (pi-back-channel-dispatch) — operates on the SAME adapter files; sequence AFTER b2 lands to avoid conflict.

## Folded-in contract-hygiene findings (captain review 2026-06-20)

These are the token reduction the binary-migration program OWES — the END (a measurably leaner contract), not deferrable polish. AC-1's "cumulative line delta NEGATIVE" is the right gate; hold every item below to it. All are contract-file (scaffolding) edits → dispatched worker, contractlint-guarded; none behavioral.

- **Decouple remaining step-number coupling.** #418 added first-class `«worker.spawn»` and `«worker.shutdown»` headings and converted Codex/Pi adapters to binding blocks. This task owns any remaining mutable step-number references in Claude/shared dispatch or merge prose; replace them with capability names instead of duplicating teardown sequence text.
- **Evergreen the contract — strip temporal markers.** The operating contract must carry NO sprint/version/roadmap temporality: `→ **shipped** (this sprint):` (×4 — `fo-merge-core.md:14`, `first-officer-shared-core.md:166/174/181`) → `→ **shipped**:`; "tracked for 0.20.6" (`claude-fo-dispatch.md:159`); "(0222) … descoped to roadmap 0222" (`fo-dispatch-core.md:150`). These annotations encode means-thinking (what shipped this sprint), not value.
- **Name capabilities concretely — drop generic "the verb".** "the verb" (×7 — shared-core ×3, fo-merge-core ×4) is confusing; say `«merge.guard»` / `spacedock merge guard` (or the specific capability).
- **Binding-location decision is resolved for Codex/Pi.** #418 made per-runtime binding blocks the preferred shape for those adapters. For remaining runtimes, follow `docs/runtime-support.md`: put concrete host tools in runtime binding sections; keep shared/core text on capability names.

Root insight (see entity `p2` / `pr-complete-binary-command`): the contract GREW (`fo-merge-core.md` 41→49 this sprint) because `merge guard` is a THIN finalize-helper — the prose carries "when X signals, the FO does Y" choreography on top of the still-FO-owned orchestration. The verb is the MEANS; the leaner contract is the END, and it was never harvested. The real collapse is the FAT `spacedock pr complete` (p2) absorbing the orchestration so the choreography prose disappears.
