---
title: "Retire catastrophe-clause prose backstops that masquerade as safety tests"
status: backlog
source: "Contractlint antipattern sweep, 2026-07-11: internal/contractlint/prose_function_backstop_test.go checks rebase and reuse-diagnostic wording rather than a safety outcome."
score: 0.27
id: h6fxjt306rxwdzjsrv1cchrs
sprint:
group:
sprint-readiness:
---

## Problem

`prose_function_backstop_test.go` makes a literal clause about rebase safety and a reuse diagnostic appear enforced. The instruction text is written by the same change being checked; a text match cannot show that a conflicting rebase is aborted or that runtime diagnostics appear when needed.

## Proposed approach

For each clause, identify the actual enforceable surface. Add a behavioral test if a binary/adapter guard owns it; otherwise remove the test's enforcement claim and retain prose as non-testable operator guidance. Do not introduce synonym lists or broader grep rules.

## Out of scope

Automating human conflict resolution or changing the conflict-resolution policy.

## Acceptance criteria

**AC-1 (VALUE) - Safety evidence distinguishes real guard behavior from operator guidance.**
Verified by: a negative/fixture behavior test for each executable guard, or an explicit focused removal of the prose-only assertion where no such guard exists.

**AC-2 - No catastrophe-clause literal matcher remains presented as behavioral proof.**
Verified by: focused contractlint tests and a review of the affected test's remaining assertion kinds.

## Test plan

Exercise the smallest applicable conflict or diagnostic fixture before changing assertions. Run focused contractlint tests and the owning command/runtime test package, then `go test ./...`.
