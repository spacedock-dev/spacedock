---
title: "Split Codex runtime semantics from contractlint phrase checks into live or fixture behavior tests"
status: backlog
source: "Contractlint antipattern sweep, 2026-07-11: codex_multi_agent_v2_contract_test.go and Codex portions of runtime_binding_block_test.go assert runtime meaning from host-adapter prose."
score: 0.34
id: 8413fc05vpp8116k54x8br15
---

## Problem

The Codex runtime contract tests encode semantic requirements such as capability binding, wait/reuse meaning, and validation routing as expected prose. They can pass after a behavior regression if the wording remains, and fail after a safe paraphrase.

## Proposed approach

Map each Codex semantic claim to the narrowest independent proof: launcher fixture, dispatch-build output, or live shared scenario where host behavior is the claim. Retain contractlint only for real structural properties such as reference closure, bounded block shape, and host-local file placement.

## Out of scope

Changing Codex multi-agent policy, model defaults, or the set of Codex tools.

## Acceptance criteria

**AC-1 (VALUE) - Codex capability, wait, reuse, and validation-route claims are checked through an observable runtime or fixture outcome.**
Verified by: focused launcher/dispatch tests and, where necessary, a live shared Codex scenario that records the actual tool call and resulting workflow state.

**AC-2 - Codex semantic phrase checks are removed or reduced to purely structural assertions.**
Verified by: focused contractlint inventory/test results showing no runtime meaning derives from the adapter's literal wording.

## Test plan

Inventory assertions by claim before edits. Add the smallest behavior test per claim, run focused Codex live tests serially when a live claim changes, then run `go test ./...` and `go test ./... -race`.
