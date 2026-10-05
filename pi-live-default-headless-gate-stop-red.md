---
title: Repair the Pi default-headless-gate-stop red
sprint: pi-live-completeness
source: "Run at chain top c94323204: FAIL /default-headless-gate-stop owner= observed=[gate-hold-violation gate-not-held implementation-worker-not-dispatched], after the gate-prepare selected-source anchor was corrected."
id: penfp034pt9s3cgwp7wg3ykk
status: backlog
---

## Problem

The Pi `default-headless-gate-stop` journey fails live with the missing-prepare reds plus the worker
observation red, after the gate-prepare anchor fix.

## Value

This journey guards the gate-stop posture and is red for reasons no task owns.

## Acceptance criteria

**AC-1** The live journey passes with no XFAIL binding.
**AC-2** Removal of its binding requires an independent task-level fix plus five consecutive passes.

## Verification

Run the Pi lane with `live_cadence=pi` and read the register.
