---
title: "Replace contractlint's FO write-classifier policy interpreter with independent enforcement proof"
status: backlog
source: "Contractlint antipattern sweep, 2026-07-11: internal/contractlint/fo_write_core_mutation_gate_test.go executes policy parsed from the instruction file it polices."
score: 0.31
id: fw2sd7w5aq6wzjgchvpm1vns
sprint: 0260-proportionality
group: contract-cleanups
---

## Problem

`fo_write_core_mutation_gate_test.go` parses the `FO-WRITE-CLASSIFIER` Markdown table and runs a test-only classifier over its rows. The test can confirm that the prose table spells out a policy, but no production consumer executes that table, so it cannot prove that an FO mutation is guarded in a real workflow.

## Proposed approach

In ideation, choose the smallest honest boundary: move a genuinely executable policy to a production-owned representation and test its observable allow/deny behavior, or retire the semantic interpreter while retaining only structural marker/closure checks. Do not replace it with another instruction-text matcher.

## Out of scope

Changing FO write authority or broadening direct-FO product-edit permission.

## Acceptance criteria

**AC-1 (VALUE) - A forbidden FO mutation is rejected by an independently executed guard, or the repository explicitly stops claiming that prose parsing enforces it.**
Verified by: a negative behavior test against the production command/guard, or a focused removal test showing no semantic policy interpreter remains in contractlint.

**AC-2 - Contractlint retains only structural checks for the classifier block if the block remains documentation.**
Verified by: focused contractlint test coverage that checks block boundaries/closure without deriving allow-or-deny policy from its prose.

## Test plan

Start with focused tests for `internal/contractlint/fo_write_core_mutation_gate_test.go`. If an executable guard is introduced, add command- or package-level positive and negative cases; otherwise run the focused contractlint package plus `go test ./...`.
