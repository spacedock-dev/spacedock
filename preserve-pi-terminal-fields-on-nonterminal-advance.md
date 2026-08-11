---
title: Preserve Pi terminal fields on nonterminal advance
status: backlog
score: "1.0"
source: Captain recovery directive; fh6 commit 4a98f40b4, 2026-08-11
sprint: test-behavior-completeness
sprint-readiness: ready
group: pi-product
id: kqdnfzjh921ryad7n6h82m1a
gates:
    version: 1
    records:
        - id: gate:kqdnfzjh921ryad7n6h82m1a:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:kqdnfzjh921ryad7n6h82m1a-backlog-1
              briefing:
                id: briefing:kqdnfzjh921ryad7n6h82m1a:backlog:attempt-1:revision-1
                digest: sha256:32e187c414902c2a7644a6db6ab005ed6909d8a31c409de008224bfa005637cf
                request-digest: sha256:592276c34d915145bb9c1d69da1e820b377379c81dd858f964a5d593e95e14b9
                room-ref: ./preserve-pi-terminal-fields-on-nonterminal-advance/review/backlog/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-11T00:41:41.736257Z"
                reason: seed lacks required backlog Stage Report
            - id: gate-attempt:kqdnfzjh921ryad7n6h82m1a-backlog-2
              briefing:
                id: briefing:kqdnfzjh921ryad7n6h82m1a:backlog:attempt-2:revision-1
                digest: sha256:7a96f5c08af39757fb5ce74b61f8c733a2ce4029204f66286cd12c0893e72f4e
                request-digest: sha256:f6e73e87b5350b673a1ff7f0a87aea5a506dd6e91ffdce3c7ce5aac6b28e2790
                room-ref: ./preserve-pi-terminal-fields-on-nonterminal-advance/review/backlog/briefing-2
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

## Stage Report: backlog

- DONE: A Pi nonterminal advance preserves legitimate `completed` and `verdict` bytes.
- DONE: The scope reverts only the terminal-field-clearing clause from fh6 commit `4a98f40b4`. It retains the oracle improvements and adds the smallest behavioral test.
- DONE: The exclusions are mechanisms, hooks, protocols, state stores, parser loops, XFAIL changes, live runs, Pi runs, and CI changes.
- DONE: AC-1 uses a focused behavioral test that proves byte-for-byte preservation of both terminal fields.
- DONE: AC-2 uses the exact source diff and focused contract checks to prove the narrow instruction reversal.
- DONE: AC-3 uses exact diff inspection and focused ownership checks to prove that the oracle improvements and XFAIL records are unchanged.
- DONE: AC-4 uses focused tests, full tests, race tests, gofmt, and `git diff --check` on one immutable candidate.
