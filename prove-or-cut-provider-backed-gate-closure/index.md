---
title: Prove or cut provider-backed gate closure before stable v1
status: backlog
source: "Pre-0.27 gate-machinery necessity audit, 2026-08-01: provider evidence exists in pilot state, but the stable surface lacks one pinned exact-candidate end-to-end proof while the chat journey is the primary value path."
started:
completed:
verdict:
score: "0.85"
worktree:
issue:
pr:
sprint: durable-decisions
id: a732sahay8wzgqrd2yr0xxr7
---

Provider-backed closure is conditional v1 scope. Keep it only if the exact release candidate can complete the same gate transaction as chat through one pinned provider package; otherwise remove the provider-only public surface from 0.27 rather than ship an unproven second path.

## Problem

The core chat journey already provides value without Subspace. Provider rooms add a second authority path, retained evidence package, capability handshake, Git-root materialization, and request schema. Archived dogfood proves pieces of that path, not one exact-candidate transaction. Carrying the path without that proof expands the trust boundary immediately before stable v1.

## Proposed approach

Run one pinned release-candidate exercise using a real prepared room and one provider package. The public agent experience stays room-only: prepare and commit the room, hand the room to the provider, then record and consume through Spacedock. The agent never reconstructs request fields or output paths. If the proof fails for a product defect that cannot be corrected within this ticket's small boundary, cut provider-backed recording and its public skill wiring from v1; do not add compatibility or fallback machinery. If retained, collapse duplicated request authority (`actor` and `approver`) to the one value the recorder actually enforces.

## Streamlined common scenarios

- Default chat: `gate prepare` -> state commit -> present in chat -> `gate record --decision ...` -> state commit -> `gate consume` -> state commit. No provider installation is needed.
- Provider override: `gate prepare` -> state commit -> `/subspace:r gate <room>` -> `gate record --room <room>` -> state commit -> the same `gate consume` path. Spacedock validates the room; the agent supplies no JSON fields or output paths.
- Provider unavailable or incapable: use the default chat path before presentation. Do not partially launch, retry through another provider, or translate a failed provider package into a chat approval.

## Out of scope

Do not add provider-specific dependencies to the Spacedock binary, a multiplexer contract, multi-provider fallback, exact-version matching, or new presentation semantics. Multiple Artifact support remains separately owned.

## Acceptance criteria

**AC-1 - Stable v1 has exactly one proven provider-backed transaction or no public provider-backed closure surface.**
Verified by: either a retained exact-revision live evidence package covering prepare, capability probe, presentation, record, validate, and consume on the candidate; or CLI/skill tests proving the provider-only entry points are absent while chat remains green.

**AC-2 - A retained provider path remains room-only for the agent.**
Verified by: the live invocation accepts only the room locator and provider selection; requiring the agent to parse `request.json`, name outputs, or reconstruct actor/approver makes the test fail.

**AC-3 - A retained request carries one unambiguous authority value.**
Verified by: request-schema and recorder tests with one authority field, plus a malformed duplicate/conflicting-authority refusal. This criterion is not applicable if the provider path is cut.

## Test plan

Use the existing provider-room fixture and one pinned real provider package. Do not create a new multiplexer harness. The proof must run on the exact candidate; archived evidence is context, not the pass condition.
