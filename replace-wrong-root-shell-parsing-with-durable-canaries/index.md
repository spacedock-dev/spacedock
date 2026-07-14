---
title: Replace wrong-root shell parsing with durable canaries
status: backlog
source: "Test-infrastructure audit 2026-07-14."
started:
completed:
verdict:
score:
worktree:
issue:
milestone: 0.26.0
id: jr7ap76zj69xr25gkvv30zbc
---

## Problem

The wrong-root detector has grown into a partial shell interpreter that tokenizes selected Bash shapes, recognizes workflow arguments, correlates tool results, and canonicalizes paths. It cannot generally establish execution scope and has accumulated false-positive repairs.

## Required outcome

Prove isolation through observable state: launch with explicit cwd/workflow targeting, plant distinct fixture and real-checkout canaries, assert intended fixture mutation, and assert the real checkout remains unchanged. Use structured argv/cwd capture only when diagnosis requires it.

## Acceptance criteria

- **AC-1 (VALUE):** A live or fixture-backed journey proves the disposable workflow changes as expected while the real checkout's canaries and Git state remain unchanged.
- **AC-2:** The launch boundary records or fixes cwd and workflow targeting through structured fields rather than inferred shell narration.
- **AC-3:** Claim-breaking wrong-root mutations fail the test, including a command shape the former partial tokenizer could misclassify.
- **AC-4:** The shell grammar stops growing and is deleted or retained only as non-authoritative diagnostics.

## Mechanism/value trace

Simplest route: before/after durable canaries plus structured cwd/argv. If a shell parser cannot prove the value for arbitrary commands, do not repair its grammar.
