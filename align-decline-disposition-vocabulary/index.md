---
title: Accept all declared non-material finding classes
status: backlog
source: "Captain-approved local task from debrief issue 2, 2026-08-02"
score: 0.7
worktree:
issue:
id: shra0x0r2bf7ka0q1m4ft79a
---

The gate parser must accept the non-material finding classes that the workflow contract declares.

## Problem

The workflow defines `deferred-risk` and `polish`, but the parser accepts only `correct-but-disproportionate` for a decline. A valid documented disposition therefore fails as invalid.

## Proposed approach

Align the parser with the declared class vocabulary. Keep materiality rules in the workflow. Add command and parser tests for every declared class and for malformed dispositions.

## Out of scope

Do not change finding ownership, materiality policy, or the review-round storage model.

## Acceptance criteria

**AC-1 - Every declared non-material class parses as a valid decline disposition.**
Verified by: parser and CLI tests that submit `deferred-risk`, `polish`, and `correct-but-disproportionate` records and assert successful validation.

**AC-2 - Invalid class names remain rejected.**
Verified by: a negative parser test that changes one class name to an undeclared value and asserts the existing invalid-disposition error.

## Test plan

Run focused parser and CLI tests. Run `go test ./...` for the full suite.
