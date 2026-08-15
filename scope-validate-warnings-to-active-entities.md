---
id: f0zn4sr7nz7xsxmyw6aw6bsm
title: Scope validate warn channels to active entities
status: backlog
source: "0.27 cut audit (2026-08-14), adversarially verified; captain directed filing"
started:
completed:
verdict:
score:
worktree:
issue:
---

Stop emitting warn-tier findings for archived entities in `status --validate`. Today 125 of 126 report lines are archived-scope warnings (121 unknown-gate-application-field, 4 verdict-enum), the report still ends VALID, and archived scope is publish-only, so no tool-mediated fix can ever silence them. The alarm fires identically forever and carries no information. The pre-commit hook echoes the full report on every state commit, so every commit dumps 51KB of noise.

Scope: filter the two warn channels (internal/status/validate.go:228-244 gate diagnostics; internal/status/field_conformance.go:88-108 enum conformance) to active entities. Keep structural and ID validation scope-inclusive. Keep the read tolerance itself - it is load-bearing for every gates read over the legacy corpus. Precedent: 6c45fd59c fixed the identical pathology for verdict case with the same rationale.

## Problem

{Ideation fills this in.}

## Proposed approach

{Ideation fills this in.}

## Out of scope

The pre-commit hook echo (kept - it has a demonstrated consumer). The read tolerance. Active-scope warnings.

## Expected surface and tolerance

Estimate net LOC change: small, across ~2 files. Observable semantics change: validate output over archived entities.

## Acceptance criteria

**AC-1 - A clean validate run over docs/dev emits zero archived-scope warn lines and still ends VALID.**
Verified by: one-off run recorded in the report, baseline 125 warnings before, 0 after.

**AC-2 - An active entity with an unknown application field or an out-of-enum verdict still warns.**
Verified by: the existing active-fixture tests stay green (gate_application_warning_test.go, field_conformance_test.go).

**AC-3 - The suite stays green.**
Verified by: go test ./... and go test ./... -race pass.

## Test plan

One archived-fixture test asserting warn absence, existing active-fixture tests as the regression floor.
