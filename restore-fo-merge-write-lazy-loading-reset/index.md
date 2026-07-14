---
title: Restore lazy loading for first-officer merge and write cores (clean reset)
status: validation
source: Clean reset from rejected 1k implementation, captain direction 2026-07-14
started: 2026-07-14T14:12:36Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-restore-fo-merge-write-lazy-loading-reset
issue:
milestone: 0.25.0
id: gk7ceyrs4496jgp535w3awfp
---

Restore the intended first-officer loading boundary: boot reads the shared core and active runtime adapter, while write authority loads at the first FO-authored mutation and merge handling loads only at a terminal or merge-mod recovery boundary.

## Problem

The earlier fix for delayed-reference filesystem hunting overcorrected by eagerly importing `fo-write-core.md` and `fo-merge-core.md` from the first-officer entry skill. A mutation-free interactive run that stops at a gate now pays for mutation and terminal ceremony it never uses. The defect was nondeterministic discovery, not lazy loading: deferred reads must use the loader-supplied first-officer base plus one literal canonical suffix, never cwd, wrapper-skill discovery, alternate paths, or search.

The prior implementation branch and its proof harness were rejected. This task intentionally carries no prior stage reports, feedback cycles, parser design, or implementation plan.

## Required outcome

- `skills/first-officer/SKILL.md` eagerly imports only the shared core; the active runtime adapter remains a boot read.
- A mutation-free gate hold reads neither write nor merge core.
- The first FO mutation reads the exact write core before the mutation. The first terminal or `mod-block=merge:*` recovery reads the exact merge core before merge-owned action; when both trigger together, write precedes merge.
- Merged-PR discovery is not required on boot; it may happen on `engage`.
- Do not change mutation resolution, merge resolution, routing, or terminal semantics. If a resolution is broken, that is a separate defect shared by eager and lazy loading.

## Mechanism/value trace

- Served value: avoid at least 8,000 bytes of unused cold contract while preserving supported-host behavior.
- Simplest route: change the import/load cues and observe real host tool-call order plus durable workflow outcomes through the existing scenario runners.
- Rejected heavier route: no new command parser, operation-language interpreter, PTY/session controller, lifecycle state, daemon, lease, or recovery protocol. If existing host events cannot prove a claim, stop for architecture review instead of implementing a second runtime.

## Acceptance criteria

- **AC-1 (VALUE):** Fresh Claude and Codex mutation-free gate journeys load shared core plus the selected adapter, read neither deferred core through the gate, and save at least 8,000 cold bytes against the eager baseline.
- **AC-2:** Existing filing and terminal/recovery journeys observe exact write/merge reads before their real supported commands, with no broad search, wrapper invocation, alternate-path retry, or earlier owned action; their durable success/refusal/archive outcomes remain unchanged.
- **AC-3:** Structural coverage proves one eager canonical import and two resolvable delayed files only. It does not infer runtime order from instruction prose.
- **AC-4:** Focused tests, `go test ./...`, `go test ./... -race`, and the relevant exact-head local Codex live journeys are green before Roborev is requested. The harness orchestrates and observes the supported runtime; it does not model arbitrary shell semantics.

## Stage Report: implementation

- DONE: Deliver the must-have deferred write/merge loading boundaries while preserving real gate, mutation, and terminal outcomes.
  Commit `09288711` leaves only the shared core eager, defers exact write/merge references, and keeps gate hold, atomic filing, and terminal refusal outcomes green on both hosts.
- DONE: Use existing supported-host event streams and scenario runners; do not implement a shell parser, runtime controller, or parallel lifecycle, and obtain exact-head local Codex live green before Roborev.
  Existing gate, filing, and merge-hook runners now grade emitted Claude/Codex path and action events; final exact-head Codex was 3/3 green before any Roborev request.
- DONE: Keep the change narrow; if another controller/lifecycle layer appears or no visible product proof emerges within 90 minutes, stop and request a design reset.
  The change is limited to load cues, structural guards, and one bounded event-field observer; it adds no runtime product code, controller, lifecycle, or shell grammar.
- DONE: AC-1 (VALUE)
  Fresh exact-head Claude and Codex gate journeys read shared core plus their adapter and neither deferred core; deferring the unchanged write+merge files saves `5,843 + 2,830 = 8,673` cold bytes.
- DONE: AC-2
  Exact-head Claude and Codex were each 3/3 green: filing reads write before `new`, terminal reads write then merge before `status --set`, and durable success/refusal/no-archive outcomes remain enforced.
- DONE: AC-3
  Contractlint proves the sole eager import, the two canonical loader-base suffixes, reference closure, ceremony anchors, and absence of wrapper core skills; runtime order comes only from host events.
- DONE: AC-4
  Focused trace controls, `go test ./...`, `go test ./... -race`, exact-head Codex 3/3, and exact-head Claude 3/3 all passed on commit `09288711`.

### Summary

Restored the intended cold boundary without changing mutation, merge, routing, or terminal resolution semantics. The existing supported-host journeys now make the load order visible and falsifiable while retaining their durable outcome assertions.

### Roborev

- Exact-head review requested for `09288711f8c9069eb7afefe5057c7c737236f4d2` after the recorded Codex and Claude live-green runs. No prior review existed for that SHA.
- Request: Roborev job `1423`, agent `codex`, reasoning `thorough`; status `done`, verdict `F` (finished 2026-07-14T23:27:09+08:00).
- Medium finding at `internal/ensigncycle/fo_deferred_load_trace_test.go:39`: `observeLoad` can count any successful event containing a matching path suffix as a completed read. Roborev recommends requiring the exact loader-supplied base and host-native read operation, correlating a genuine full-file Codex read, and adding negative cases for path-only mentions, partial reads, and same-suffix alternate roots.
- Medium finding at `internal/ensigncycle/fo_deferred_load_trace_test.go:172`: the checks cover atomic filing and `status=done`, but not the broader set of FO-authored mutations or every existing mutating live scenario. Roborev recommends detecting the first general FO mutation across those scenarios while retaining the specialized terminal merge-order assertion.
- No code was changed in response; both concrete findings await first-officer feedback routing.
