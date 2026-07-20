---
title: Retire non-authoritative test oracles
status: backlog
source: "Test-infrastructure audit 2026-07-14."
started:
completed:
verdict:
score:
worktree:
issue:
milestone: 0.26.0
id: xazfg0zej10pzthyaj14v9kq
sprint: 0260-proportionality
group: test-cleanups
sprint-readiness: defer
---

## Problem

Three active test clusters no longer authoritatively prove product behavior: the inferred dispatch-close state machine has no live caller, archived Haiku spike tests can log behavioral breaks without failing, and the product-edit guard oracle is consumed only by its own synthetic tests. Keeping them active creates maintenance and false-confidence costs.

## Required outcome

Delete obsolete inferred lifecycle and self-testing oracles while retaining the simple stream capture, process bounds, cleanup, and durable historical records that still serve current tests.

## Acceptance criteria

- **AC-1 (VALUE):** The active suite contains no test that can observe its protected behavior fail and still report green by design.
- **AC-2:** The unused dispatch-close reconstruction and its self-tests are removed without removing basic stream capture, quiet timeout, or process cleanup.
- **AC-3:** Executable Haiku spike tests are retired while their archived research evidence remains available.
- **AC-4:** The orphan product-edit oracle and synthetic-only tests are removed; any replacement must directly observe a disposable repository and native worker-dispatch evidence.
- **AC-5:** `go test ./...` and `go test ./... -race` pass, and the cleanup has a materially negative test-infrastructure line delta.

## Mechanism/value trace

Simplest route: delete code with no authoritative live consumer. Do not replace dead state machines unless a current value AC first demonstrates a missing proof boundary.
