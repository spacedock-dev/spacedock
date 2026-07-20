---
title: "Replace reuse/advance prose pointers with dispatch-build behavior coverage"
status: backlog
source: "Contractlint antipattern sweep, 2026-07-11: internal/contractlint/reuse_advance_pointer_test.go treats instruction call-shape prose as runtime proof."
score: 0.29
id: 1wd14erf3dm632qx6yrzvfqw
sprint: 0260-proportionality
group: contract-cleanups
sprint-readiness: defer
---

## Problem

`reuse_advance_pointer_test.go` asserts that FO, Claude-route, and ensign instruction files name an `advance` helper. Those phrases can be paraphrased or retained while the real `spacedock dispatch build --advance` path stops producing a reusable next-stage dispatch.

## Proposed approach

Replace the semantic phrase assertions with behavior coverage for `spacedock dispatch build --advance`: a fixture-backed request should produce the correct next-stage dispatch artifact and host route, including the reuse-versus-fresh decision inputs the binary owns. Preserve only reference-closure or host-file locality checks that are genuinely structural.

## Out of scope

Changing the workflow's reuse policy or adding new runtime hosts.

## Acceptance criteria

**AC-1 (VALUE) - A generated advance dispatch artifact is executable evidence of the intended next-stage route rather than a prose pointer.**
Verified by: focused `dispatch build --advance` behavior tests with fixture inputs and observable emitted output/state.

**AC-2 - The three semantic instruction-file assertions no longer stand in as proof of advance behavior.**
Verified by: focused contractlint test results plus a source-level test inventory that leaves only structural checks in this area.

## Test plan

Add focused CLI/package fixtures first, including a negative route or malformed-stage case. Run the relevant dispatch tests, `go test ./internal/contractlint/...`, and `go test ./...`.
