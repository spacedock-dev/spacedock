---
title: Trim remaining dispatch adapter prose into capability bindings
status: validation
sprint: 0230-stable-finalization
group: binary-ux
id: adk755xqeb4a9dxhhgtjwawh
started: 2026-06-21T06:05:14Z
worktree: .worktrees/spacedock-ensign-trim-dispatch-adapter-prose
---

After the b2 capabilities-as-«fn» reframe and #418's Codex/Pi binding-block cleanup, the target is clear: per-runtime files should be capability bindings, not lifecycle narration. The shared core owns when each `«fn»` runs; adapters bind how the host realizes it.

This task is the careful follow-up cut for remaining dispatch/adapter prose, especially Claude and shared dispatch residue. Do not redo #418's Codex/Pi binding-block rewrite or the first-class `«worker.spawn»` / `«worker.shutdown»` core-heading work it landed. Instead, collapse remaining adapter narration into `«fn»` binding lines where the semantics are already represented by capability names.

## Hard constraint (the SB2-class risk this task must avoid)
Do NOT cut await/reuse/guardrail contract to hit a line count. The `## Awaiting Completion` premature-reap ban, the idle guardrail, and any reuse-advance/await semantics not captured by the terse `→` lines are load-bearing and STAY byte-intact. If a candidate block carries unique contract, KEEP it.

## Acceptance criteria
- **AC-1 (VALUE — measured against the v0.22.0 baseline that can move the wrong way)** — cumulative byte delta of the two cut-surface files vs the v0.22.0 git tag is NEGATIVE and both hit the gate: `git show v0.22.0:…/fo-dispatch-core.md | wc -c` = 17488 and `…/fo-merge-core.md | wc -c` = 8059; after the cut, working-tree fo-dispatch-core.md <= 17488 B AND fo-merge-core.md <= 8059 B. Report both `wc -c` figures and the signed cumulative delta (now 22929+8597=31526 -> target <= 25547, a cut of >= 5979 B) in the stage report. This is the DoD slice the value gate measures.
- **AC-2 (SB2-class byte-intact assertion — names the ACTUAL load-bearing text in the cut-surface files)** — in fo-dispatch-core.md the `## Reuse and Fresh Dispatch` reuse-condition list (conditions 0-4), the model-mismatch diagnostic line containing the verbatim anchor `does not match next stage effective_model`, the supersede-shutdown paragraph, and the `«async-dispatch»`/`«completion-signal»`/`«context-budget»` `→`/`block` lines are SEMANTICALLY UNCHANGED (no reuse/await/guardrail clause dropped); in fo-merge-core.md the `«merge.guard»` armed/blocked/finalized phase semantics and the never-finalize-on-pr-presence-alone rule are unchanged. Verified by diffing those spans and confirming only narration around them moved. Separately assert the out-of-scope SB2 sections in claude-fo-dispatch.md (`## Awaiting Completion`, premature-reap ban, DISPATCH IDLE GUARDRAIL, lines 58-95) are BYTE-IDENTICAL to pre-cut — this member must not touch that file; `git diff` shows zero changes to claude-fo-dispatch.md.
- **AC-3 (mechanism — paired with AC-1's value gate)** — the trimmed core text follows the pseudo-code contract principle in `docs/runtime-support.md`: capability names + compact `guard`/`effect`/`done-when`/`block`/`→` binding lines first; prose retained ONLY for fuzzy judgment, probe-backed host quirks, or a guardrail that does not fit one capability. The `spacedock dispatch build` block is collapsed to a `«dispatch.build»`-shaped body with the verbatim-forward rule and a single break-glass `block:` line, not the multi-step procedural paragraph.
- **AC-4 (mechanism guard)** — the contract is evergreen — no sprint/version/roadmap temporal marker survives in either cut file (`grep -nE '\(0222\)|this sprint|roadmap 022|0\.2[0-9]\.[0-9]' fo-dispatch-core.md fo-merge-core.md` returns nothing).
- **AC-5 (regression gate)** — `go build ./...` and `go test ./internal/contractlint/` are green — specifically the host-neutrality/capability-binding/boot-resident-closure/prose-function-backstop tests (capability_binding_test.go, boot_resident_closure_test.go, structural_checks_test.go, prose_function_backstop_test.go, layering_restore_test.go) pass, proving the trimmed cores still satisfy the lazy-load naming and capability-surface contract. Scoped to the dispatch/merge core prose; no behavioral code change.
- **AC-6 (docs/runtime-support.md alignment)** — `docs/runtime-support.md` reflects the final `«dispatch.build»`/capability-binding shape this cut lands (the `## Dispatch Adapter` / binding-block conventions), so the authority doc and the trimmed fo-dispatch-core/fo-merge-core files do not drift. Cross-check the documented conventions against the shipped shape.

Sequencing: b2 (pi-back-channel-dispatch, #409) is MERGED — the same-file conflict hazard is cleared and the b2 hold is moot. The only remaining sequencing constraint is that M2 must NOT edit shared-core's `→ **shipped** (this sprint)` markers (those are M1 / `fo-contract-token-cut`'s gate-file territory); M2 harvests only fo-merge-core.md's shipped marker and fo-dispatch-core.md's `(0222)` marker.

## Folded-in contract-hygiene findings (captain review 2026-06-20)

These are the token reduction the binary-migration program OWES — the END (a measurably leaner contract), not deferrable polish. AC-1's "cumulative line delta NEGATIVE" is the right gate; hold every item below to it. All are contract-file (scaffolding) edits → dispatched worker, contractlint-guarded; none behavioral.

- **Decouple remaining step-number coupling.** #418 added first-class `«worker.spawn»` and `«worker.shutdown»` headings and converted Codex/Pi adapters to binding blocks. This task owns any remaining mutable step-number references in Claude/shared dispatch or merge prose; replace them with capability names instead of duplicating teardown sequence text.
- **Evergreen the contract — strip temporal markers (M2's two cut files only).** The operating contract must carry NO sprint/version/roadmap temporality. This member harvests only the markers in ITS gate files: `fo-merge-core.md`'s `→ **shipped** (this sprint):` → `→ **shipped**:`, and `fo-dispatch-core.md`'s "(0222) … descoped to roadmap 0222". The shared-core `→ **shipped** (this sprint)` lines (`first-officer-shared-core.md:166/174/181`) are M1 / `fo-contract-token-cut`'s gate-file territory — M2 must NOT edit them. `claude-fo-dispatch.md`'s "tracked for 0.20.6" marker is in the unowned do-not-touch file (AC-2 asserts zero diff there). These annotations encode means-thinking (what shipped this sprint), not value.
- **Name capabilities concretely — drop generic "the verb".** "the verb" (×7 — shared-core ×3, fo-merge-core ×4) is confusing; say `«merge.guard»` / `spacedock merge guard` (or the specific capability).
- **Binding-location decision is resolved for Codex/Pi.** #418 made per-runtime binding blocks the preferred shape for those adapters. For remaining runtimes, follow `docs/runtime-support.md`: put concrete host tools in runtime binding sections; keep shared/core text on capability names.

Root insight (see entity `p2` / `pr-complete-binary-command`): the contract GREW (`fo-merge-core.md` 41→49 this sprint) because `merge guard` is a THIN finalize-helper — the prose carries "when X signals, the FO does Y" choreography on top of the still-FO-owned orchestration. The verb is the MEANS; the leaner contract is the END, and it was never harvested. The real collapse is the FAT `spacedock pr complete` (p2) absorbing the orchestration so the choreography prose disappears.

## Stage Report: implementation

- DONE: AC-1 (the value gate): fo-dispatch-core.md ≤ 17488 B AND fo-merge-core.md ≤ 8059 B
  fo-dispatch-core.md 22929→17446 B; fo-merge-core.md 8597→7436 B; cumulative 31526→24882 B, signed delta −6644 B (cut > 5979 target). Commit 2171cf77.
- DONE: SB2 byte-intact protection (in-surface reuse/await text + «merge.guard» phase semantics + never-finalize rule stay byte-intact; git-diff-ZERO to claude-fo-dispatch.md; do NOT touch shared-core's `→ shipped` markers)
  Verified byte-intact vs HEAD: reuse conditions 0-4, the verbatim `does not match next stage effective_model` diagnostic, supersede-shutdown, the «async-dispatch»/«completion-signal»/«context-budget» →/block lines, the «merge.guard» effect (armed/blocked/finalized) + `never finalize on pr-presence alone`. `git diff HEAD` shows zero changes to claude-fo-dispatch.md and first-officer-shared-core.md.
- DONE: go build + contractlint green; live gate/rejection/merge-hook scenarios stay green; docs/runtime-support.md alignment for any binding shape this cut changes
  `go build ./...` and full `go test ./...` green (including internal/cli's prose-function routing binder and the full contractlint suite — capability_binding, boot_resident_closure, structural_checks, prose_function_backstop, layering_restore). docs/runtime-support.md already documents the `→ shipped`/adapter-owns-host-tools conventions this cut lands; no doc edit needed (the cut moved the core toward the doc's principle, no drift). Live `-tags live` scenarios not run in this worker (validation/captain territory).

### Summary

Collapsed the multi-step `spacedock dispatch build` procedure into a compact `«dispatch.build»()` guard/effect/done-when/block/→shipped body (AC-3), deferred the duplicated per-host tool detail on the `«addressable-worker»`/`«worker-identity»`/`«roster-reconcile»` `→` lines to the adapters' `## Runtime implementation` blocks (where #418 already binds them), and removed the `fo-merge-core.md` armed/blocked/finalized narration bullets that restated the `«merge.guard»` definition. Evergreened the two temporal markers (the `(0222)`/`roadmap 0222` next-action note and the already-harvested merge `→ shipped` marker) without dropping the prose-function routing test's verb anchor. AC-2-protected spans and the two do-not-touch files are byte-intact.
