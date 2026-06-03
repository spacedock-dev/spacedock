---
id: "020"
title: Real external proof
status: implementation
score: "0.50"
source: roadmap
---

# Real external proof

A terminal --set on this entity must succeed under
`require-external-proof: true`.

## Acceptance criteria

**AC-1 — The classifier compiles and runs.**
Verified by: a Go test `TestClassifierBasic` in `internal/status/external_proof_test.go`
asserts the classifier returns a single flag for a constructed self-ref body.
