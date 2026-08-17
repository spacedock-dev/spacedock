---
title: Codex filing journey files one variant and skips the second
status: backlog
source: "Run 31991864922 attempt 1 codex lane, filing journey: variant 001 filed, no spacedock new for variant 002 anywhere in the stream; captain approved filing owners for tolerated residual modes at the 0.27 composite-green ruling, 2026-08-17"
id: 3ptrzm0c11egac021a54dbd8
---

## Problem

The codex FO drove the filing journey but created only the first of the fixture's two variants — the stream carries no `spacedock new` for variant 002 at all, so the grade (`assertFilingCommands`) is an honest miss on an untouched journey. Observed once; the same journey went green on the two following codex runs (31991864922 attempt 2, 31996696789). One occurrence crossed the owner bar when the 0.27 composite-green ruling tolerated it as a known noise mode.

## Proposed approach

Evidence before mechanism, because one data point cannot separate model variance from fixture underdetermination. First, re-read the filing fixture's variant enumeration with the Cycle-line lesson in mind: does it pin "file BOTH variants" as explicitly as the escalation fixture pins its target file, or can a literal reader believe one filing satisfies it? If the wording is ambiguous, pin it (the 9b1c86a60 pattern: name each required artifact, one sentence). If the wording is already determined, run a short targeted codex filing loop to measure the miss rate on shipped bytes; a recurrent rate makes this a model-adherence owner tracked with the metrics instrument, a zero rate over the loop archives it as the one-off it currently looks like.

## Out of scope

- The filing grader (honest; grades `spacedock new` invocations, the right observable).
- Claude and pi lanes (no observed occurrence).
