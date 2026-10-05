---
title: Repair the Pi auto-continue dispatch red
sprint: pi-live-completeness
source: "Run at chain top c94323204: FAIL /auto-continue-after-implementation--auto-continue/single-root owner= observed=[validation-worker-not-dispatched], after the shared assert began crediting bg_wait and the native subagent-notify completion."
id: s0gq9p69nztejw8xp3by4k7f
---

## Problem

The Pi `auto-continue-after-implementation` journey fails live with `validation-worker-not-dispatched` in both
variants, even after the shared assert was changed to credit `bg_wait` and the native `subagent-notify` completion.

## Value

Until this is repaired, the Pi lane cannot tell a real Pi dispatch failure from this known one.

## Acceptance criteria

**AC-1** The live journey passes in both variants with no XFAIL binding.
**AC-2** Its binding is restored first, and removal requires five consecutive passes.

## Verification

Run the Pi lane with `live_cadence=pi` and read the register.
