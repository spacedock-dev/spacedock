---
title: Journey metrics record unbound failures as passes
status: backlog
source: "Named for separate filing in repair-codex-rejection-round-recording; captain approved filing"
id: 3rdq29x1wn907j5363e9gn4r
---

## Problem

`scenarioBehaviorResult` (internal/ensigncycle/journey_metrics_live_test.go:181) sets `Passed: true` and only attaches the real grade when the journey carries an XFAIL binding. A red run on an unbound journey therefore ships a metrics record claiming a pass — both red rejection-flow repair-loop runs did. Downstream consumers of journey metrics see passes that were failures: the lying-label class.

## Proposed approach

Record the actual grade for every run, bound or unbound; `Passed` reflects the real outcome. Add a unit case where an unbound failing journey must produce a failing metrics record, proven by falsifying edit (reverting the fix reds it).

## Out of scope

- Rejection-flow behavior and journey grading (other entities own those).
