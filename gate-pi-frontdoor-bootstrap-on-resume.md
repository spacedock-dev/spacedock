---
id: 4avk4msa3ktyk1fdt6j5ktw1
title: Gate the Pi front door bootstrap on resume
status: ideation
source: "Captain CL, 2026-08-25: 'spacedock pi --resume didn't avoid loading the spacedock initial contract.' The Pi front door (internal/cli/pi.go) appends piBootstrapPrompt unconditionally — no containsResume gate, unlike the Claude/Codex front door (internal/cli/frontdoor.go:428,447) which suppresses its bootstrap prompt on --resume/-r/--continue/-c. The spacedock .pi extension session_start handler also re-injects FO_BOOTSTRAP_TEXT with no resume detection. CL hypothesis 'compaction hook leaked into general startup' checked and disproven: session_compact is correctly scoped to the compaction event (PR #738 / force-boot-at-compaction-boundary); the leak is the front door + session_start, neither resume-aware."
gates:
    version: 1
    records:
        - id: gate:4avk4msa3ktyk1fdt6j5ktw1:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:4avk4msa3ktyk1fdt6j5ktw1-backlog-1
              briefing:
                id: briefing:4avk4msa3ktyk1fdt6j5ktw1:backlog:attempt-1:revision-1
                digest: sha256:b3f27b5850f0b44aea082a43333494efec5ca22b1f64f3e0561a0372fe956e40
                request-digest: sha256:c7df81974554a473c4a131ead3f0471be0554004d090c66c69582482a99e48c0
                room-ref: ./gate-pi-frontdoor-bootstrap-on-resume/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:4avk4msa3ktyk1fdt6j5ktw1:backlog:1
                briefing: briefing:4avk4msa3ktyk1fdt6j5ktw1:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-27T00:11:22.702651Z"
                decision: approve
                reason: 'Captain approve: enter ideation to flesh out the approach and test plan'
              application:
                target-stage: ideation
                state: consumed
started: 2026-08-27T00:11:57Z
---

`spacedock pi --resume` must not load the Spacedock first-officer contract as if starting fresh; a resume carries its own session intent and the contract survives in the system prompt via resources_discover. Today it loads the contract via two unconditional paths.

## Problem

`internal/cli/pi.go:271` appends `piBootstrapPrompt` ("Use $spacedock:first-officer for this whole Pi session.") to the Pi argv with no resume gate — `pi.go` has zero `containsResume` calls; `frontdoor.go` (Claude/Codex) has 3 and gates its bootstrap on `if !resume` (frontdoor.go:428,447). So `spacedock pi -- --resume` still injects the launch prompt that tells the session to load the FO contract.

Second path: `.pi/extensions/spacedock.ts` `session_start` handler (line 80) sets `injectBootstrap=true` unconditionally — no resume detection. On a resume, `session_start` still fires, and the `context` hook re-injects `FO_BOOTSTRAP_TEXT` (the `hasStructuralBootstrap` idempotency guard checks `event.messages` for a prior injection, but the injection is a transient context-message modification, not a durable logged message, so a resumed transcript won't match it). The extension has compaction-awareness (PR #738: `session_compact` switches to the lighter `injectBootRecord`) but no resume-awareness.

## Proposed approach

Gate the Pi front door bootstrap on `!resume`, mirroring the Claude/Codex front door: add `containsResume(fd.passthrough)` to `pi.go` and suppress `piBootstrapPrompt` when resume is forwarded. Decide in ideation whether the extension's `session_start` should also detect resume (if Pi's session_start event carries a resume marker) and skip `injectBootstrap`, or whether gating the launch prompt alone suffices — the extension's re-injection is redundant on resume but load-bearing for compaction survival, so touching it risks the PR #738 compaction boundary. Lean: fix the front door only; leave the extension's session_start alone unless a resume marker is cheaply available.

## Risk evidence

Backlog: the gap is verified by code reading — `pi.go` has no `containsResume` call; `frontdoor.go` gates on `!resume`; the extension's `session_start` is unconditional. CL's compaction-hook-leak hypothesis checked and disproven: `session_compact` → `injectBootRecord` is a separate event scoped to compaction; it does not fire on startup/resume. Riskiest unverified mechanism: whether gating the launch prompt alone stops the contract loading on resume, or whether the extension's `session_start` re-injection independently loads it (needs a live `spacedock pi --resume` probe to settle).

## Out of scope

Changing the compaction-boundary behavior (PR #738 / `force-boot-at-compaction-boundary` owns that). The ensign child-session bootstrap exemption (`PI_SUBAGENT_CHILD=1`, owned by `pin-ensign-contract-entry-point`). The Claude/Codex front doors (already gated on resume).

## Expected surface and tolerance

Estimate net LOC change: ~+15, across ~2 files (`internal/cli/pi.go` add the `containsResume` gate + a test in `internal/cli/pi_frontdoor_test.go` or `pi_launch_test.go`; possibly `.pi/extensions/spacedock.ts` if ideation chooses to make `session_start` resume-aware). Tolerance: ±50%.

## Acceptance criteria

**AC-1 (value) — `spacedock pi --resume` does not inject the fresh-start FO contract.**
Verified by: a fixture in `internal/cli` that launches the Pi front door with `--resume` in the passthrough and asserts the assembled argv does NOT contain `piBootstrapPrompt`, while a launch without `--resume` still contains it. Fails today (pi.go appends unconditionally).

**AC-2 (serves AC-1) — the resume gate matches the Claude/Codex front door's resume token set.**
Verified by: a test asserting `--resume`, `--resume=<id>`, `-r`, `--continue`, `-c` all suppress the prompt (mirroring frontdoor.go:556-559), and non-resume passthrough does not. Fails today (no gate exists).

**AC-3 (serves AC-1) — a live `spacedock pi --resume` does not re-load the contract as a fresh start.**
Verified by: a live probe (or a fixture exercising the extension's `session_start`) confirming a resumed session does not receive `FO_BOOTSTRAP_TEXT` as a fresh injection; if the extension's `session_start` cannot be made resume-aware cheaply, record the residual and downgrade AC-1's "Verified by" accordingly so the AC matches what the front-door gate actually proves.

## Test plan

Fixture in `internal/cli`: extend `pi_frontdoor_test.go` / `pi_launch_test.go` with a resume-suppression test mirroring `frontdoor_test.go`'s codex/claude resume tests (frontdoor_test.go:848-915). Live (optional, gated on whether the extension needs a change): a `spacedock pi --resume` probe confirming no fresh `FO_BOOTSTRAP_TEXT` injection. Cost: low — the front-door gate is a few lines + a test.

### Feedback Cycles
