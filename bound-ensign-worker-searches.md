---
title: The ensign discipline does not bound a search, so a worker can search the whole machine
status: backlog
score: 0.7
source: "FO session 2026-10-03: worker run 53f7f450 ran an unbounded find across /tmp and the home directory. The captain killed it by hand."
id: tq7py5k801apdnhhv3xwarzh
---

## Problem

A worker that does not know where something lives can search the whole machine.
Run `53f7f450` did exactly that: an unbounded `find /tmp /Users/clkao` that ran
for minutes, and the captain had to kill it. Nothing in the discipline stops it.

`skills/ensign/references/ensign-shared-core.md` has a `## Rules` section and a
`## Background Bash Discipline` section. Neither one bounds a search. So every
worker is free to search the home directory, `/tmp`, or the filesystem root.

The first officer caused it too. The dispatch named an outcome — make the grade
observe two things — and did not name the file that holds the evidence. The
worker filled that gap with a filesystem hunt.

## Visible value

A worker keeps its searches inside the repo, the named directory, or one session
directory. A wide search costs seconds instead of minutes, and the operator never
has to kill a command by hand.

## Out of scope

- The first officer's dispatch content. Naming the location, when the first
  officer knows it, stays the first officer's job.
- Any tool change, runtime change, or new check.
- Any new lane.

## Expected surface and tolerance

`skills/ensign/references/ensign-shared-core.md`: one rule added under
`## Rules`. Estimate four lines, one file. Tolerance: net 0 to +10 lines, one
file.

## Acceptance criteria

**AC-1 — The discipline bounds every search.**
The file tells the worker to keep `find`, recursive `grep` and recursive reads
inside the repo, the named directory, or a single session directory, and never to
search the home directory, `/tmp`, or the filesystem root.
Verified by: the check that already loads this file,
`internal/contractlint/boundary_guard_control_test.go`, extended to require the
rule. Falsifier: remove the rule and that check must fail.

**AC-2 — The discipline says to escalate instead of hunting.**
When the worker does not know where something is, and the dispatch did not name
it, the worker asks the first officer instead of searching wide.
Verified by: the same check requiring this sentence. Falsifier: remove the
sentence and the check must fail.

**AC-3 (no-regression) —** The existing rules and the background-bash polling
discipline do not change, and the existing skill checks stay green.

## Test plan

Extend the existing reader in `internal/contractlint/boundary_guard_control_test.go`
to require both rules by their content. Do not assert a single bare word.
Deterministic checks only. No live run.
