---
title: "Split Pi runtime semantics from contractlint phrase checks into independent behavior tests"
status: backlog
source: "Contractlint antipattern sweep, 2026-07-11: pi_runtime_negative_contrast_test.go and Pi portions of runtime_binding_block_test.go infer live semantics from instruction prose."
score: 0.28
id: sgepgtm1qtjz527pyg0bee7n
sprint: 0260-proportionality
group: contract-cleanups
---

## Problem

Pi runtime semantic and negative-contrast checks search the adapter wording for required or banned meanings. They do not observe Pi's tool surface, generated dispatch, or resulting workflow state, so they are neither robust to paraphrase nor evidence of host behavior.

## Proposed approach

Keep only independently meaningful structural placement/closure checks in contractlint. Replace semantic claims with fixture-backed Pi dispatch behavior or live Pi scenarios when the host is available; where a runtime claim has no executable owner, demote the prose to guidance rather than asserting it is enforced.

## Out of scope

Building a general Pi host emulator or changing Pi's runtime support contract beyond proof coverage.

## Acceptance criteria

**AC-1 (VALUE) - Pi runtime-route and capability claims have independent observable proof or are honestly treated as guidance.**
Verified by: focused Pi fixture/live behavior tests or a focused removal of unenforceable semantic assertions.

**AC-2 - Negative-contrast phrase lists no longer serve as the proof of Pi behavior.**
Verified by: focused contractlint test results plus the new independent test evidence.

## Test plan

Classify each current assertion as structural, fixture-testable, live-testable, or guidance-only. Add/run the smallest independent cases, then run the affected Go packages and `go test ./...`.
