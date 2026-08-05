---
id: d4j5cp6g3ghvjjap7djqa9vw
title: Track shipped-surface containment closure
status: backlog
source: Captain-authorized filing from the shipped-surface containment audit, 2026-08-06
started:
completed:
verdict:
score: 0.90
group: shipped-surface-containment
worktree:
issue:
pr:
mod-block:
---

Track closure of the four user-facing containment outcomes found by the 2026-08-05 shipped-surface audit.

This is a coordination record. It does not authorize a fifth implementation or a general containment framework.

## Owned outcomes

- `make-first-officer-consumer-repository-portable` (`sbkske5gyadwfxs9vmd61zb3`)
- `commission-portable-supported-runtime-workflows` (`w9h7bzax0qdmtmp8r5eepk0h`)
- `make-debrief-consumer-repository-safe` (`9ezthwqxnff2tfchvg5z65sm`)
- `make-survey-installed-package-portable` (`f5yp12pqbp646dhrz6ayks4z`)

EP2 retains ownership of `mods/pr-merge.md`. The four outcomes do not duplicate that repair.

## Problem

The audit records release-age containment defects across four shipped user journeys. Separate tasks own the repairs. A closure record is necessary to prevent partial repair from being mistaken for complete shipped-surface containment.

## Closure rule

Close this record only after all four owned tasks reach `done` with `PASSED`. Then repeat the read-only shipped-surface audit against the release candidate and record its exact commit.

The repeat audit is evidence for closure. It is not a new committed lint, framework, or implementation task.

## Out of scope

This record owns no product code, contract rewrite, consumer harness, runtime lane, or standing enforcement. It does not join the live-test-truth sprint.

## Acceptance criteria

**AC-1 (VALUE) - Every rejected shipped surface has a passed user-facing repair or an explicit external owner.**
Verified by: the four owned tasks are `done/PASSED`, EP2 records the `pr-merge` repair, and the repeat audit names no unowned rejected surface. Reopening or removing one owner makes closure fail.

**AC-2 - Closure does not add a fifth implementation mechanism.**
Verified by: this record has no implementation worktree or product diff. Any proposed shared lint, framework, or consumer harness routes through the value task that needs it or a separate Captain decision.

## Test plan

Query the four task records and EP2's final state. Repeat the audit on the exact release-candidate commit. Record the commit, owned outcomes, remaining findings, and disposition.
