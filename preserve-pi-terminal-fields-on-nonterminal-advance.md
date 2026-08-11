---
title: Preserve Pi terminal fields on nonterminal advance
status: backlog
score: "1.0"
source: Captain recovery directive; fh6 commit 4a98f40b4, 2026-08-11
sprint: test-behavior-completeness
sprint-readiness: ready
group: pi-product
id: kqdnfzjh921ryad7n6h82m1a
---

Pi must not erase legitimate `completed` or `verdict` fields during a nonterminal advance.

## Problem

Commit `4a98f40b4` added a Pi First Officer clause that clears `completed` and `verdict` before a nonterminal status advance. Those fields can contain legitimate durable workflow state, so the automatic cleanup can destroy valid information.

## Proposed approach

Revert only the shipped terminal-field-clearing clause. Retain the useful fh6 oracle improvements. Add or adjust the smallest behavioral test that proves a nonterminal Pi advance leaves legitimate terminal fields unchanged.

## Out of scope

No test-only product mechanism, hook, protocol, state store, parser loop, XFAIL mutation, live or Pi run, CI change, or unrelated fh6 change.

## Acceptance criteria

**AC-1 (VALUE) - A nonterminal Pi advance preserves legitimate `completed` and `verdict` fields byte-for-byte.**
Verified by: a focused behavioral test that starts with nonempty legitimate values, performs the supported nonterminal-advance behavior, and asserts both values remain unchanged. Clearing either field must make the test fail.

**AC-2 - The fh6 terminal-field-clearing instruction is absent while unrelated Pi runtime instructions remain unchanged.**
Verified by: the exact source diff against main and focused instruction/runtime contract checks. Restoring the clearing clause or changing another Pi instruction must fail the scope check.

**AC-3 - Useful fh6 oracle improvements and all XFAIL bindings, assertions, reconciliation rows, and owners remain unchanged.**
Verified by: exact diff inspection plus focused oracle, registry-reconciliation, and active-owner checks. Any oracle or binding change must fail the comparison.

**AC-4 - Repository behavior remains green after the narrow reversal.**
Verified by: focused tests, `go test ./...`, `go test ./... -race`, gofmt, and `git diff --check` on one immutable candidate.

## Test plan

Ideation must locate the exact shipped clause and the nearest existing behavioral test boundary before implementation. Use the smallest unit or fixture-backed test that can falsify field preservation. No live, Pi, or CI run is permitted. The validator batches one complete adversarial matrix before one verdict. One authorized correction pass is allowed; a second candidate-owned rejection requires design reset or HOLD.
