---
title: Classify a stalled Pi live run from the artifacts the runner already archives
status: backlog
score: 0.7
source: "pi-ux carve review, 2026-10-03: stall classification has no Pi implementation owner once the Claude-only capture task stays out of the sprint."
id: z6eb1krpyfmr01dwyb7703hx
sprint: pi-ux
group: tooling
sprint-readiness: ready
gates:
    version: 1
    records:
        - id: gate:z6eb1krpyfmr01dwyb7703hx:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:z6eb1krpyfmr01dwyb7703hx-backlog-1
              briefing:
                id: briefing:z6eb1krpyfmr01dwyb7703hx:backlog:attempt-1:revision-1
                digest: sha256:df46fdbb84739f3ffc3e48e5b644510c514be915e658607aebeb62c54152990d
                room-ref: '@review/backlog/briefing-1'
---

A stalled Pi live run leaves artifacts that nobody classifies. The operator cannot
tell transient weather from a defect, so each stall costs a re-run and teaches
nothing.

## Problem

The Pi live runner already archives `pi-stdout.txt`, `pi-stderr.txt`, and the root
session JSONL. Task `live-ci-api-error-log-capture` adds an always-on API debug file
for the **Claude lane only**, and its body states that Pi needs no capture change.

Capture is therefore not the gap. The gap is that no artifact turns a Pi stall into a
classification. Existing stderr capture alone does not establish whether a stalled run
was transient weather or a CLI retry or hang defect.

## Out of scope

The Claude-lane capture change, owned by `live-ci-api-error-log-capture`. Adding a Pi
debug flag, unless evidence shows the archived artifacts cannot classify a stall.

## Expected surface and tolerance

Estimate net LOC change: +40, across 2 files (the Pi live runner and its test).
Insertions ~+55, deletions ~-15. Tolerance: +/-20 net LOC, +/-1 file.

Declared semantic changes: a stalled Pi run gains a recorded classification and its
stated limits. This task must NOT change the journey set, the grading rules, or any
XFAIL binding.

## Acceptance criteria

**AC-1 - A stalled Pi run produces a classification, within stated limits.**
Verified by: a fixture that feeds the archived stall artifacts to the classifier and
asserts the recorded result, including the limit text for the case the artifacts cannot
settle. Independent baseline that can move the wrong way: the current tree records no
classification for the same fixture.
Falsifying edit: make the classifier return a verdict for every input; the
"cannot settle" assertion must turn RED.

**AC-2 - The classifier reads the artifacts, not the test's expectations.**
Verified by: the fixture writes the artifacts and the assertion derives its expectation
from those written files. Falsifying edit: delete the artifact write; the test must fail.

**AC-3 (no-regression) -** `go test ./internal/ensigncycle/...` is green, and no existing
journey assertion or XFAIL binding changes.

## Test plan

Primary proof owner: the Pi live runner test beside `pi_shared_live_runner_test.go`.
Deterministic Go tests over fixture artifacts. No live lane is needed for the classifier;
one live lane run confirms the classification appears on a real stall.
