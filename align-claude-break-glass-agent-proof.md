---
title: Align Claude break-glass recovery proof with the current Agent tool
status: backlog
source: "PR #680 run 31640122346, Sonnet job 94260537662, artifact 9159859622: TestLiveBreakGlassShimRecovery rejected selected-bare after a successful worker completion because the Agent call lacked the expected bare marker, and rejected selected-team after observing two Agent calls."
score: 0.90
id: s2atxdv146qknetdjvx0xer6
gates:
    version: 1
    records:
        - id: gate:s2atxdv146qknetdjvx0xer6:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:s2atxdv146qknetdjvx0xer6-backlog-1
              briefing:
                id: briefing:s2atxdv146qknetdjvx0xer6:backlog:attempt-1:revision-1
                digest: sha256:0d8a3a8721fadf9acfb0dce2f6399acbcb123bafc947928ea80a03cd21299c62
                room-ref: '@review/backlog/briefing-1'
---

Restore a truthful break-glass substrate proof for the current Claude Agent tool surface.

## Problem

Both recovery variants completed useful work, but their expected Agent call topology no longer matches the current host surface.

## Proposed approach

Diagnose the current supported Agent call semantics and the duplicate-call origin. Correct the product contract only if user behavior is wrong; otherwise correct the oracle. Do not add hooks or instrumentation.

## Out of scope

Do not change common live journeys, Codex behavior, or normal dispatch.

## Acceptance criteria

**AC-1 - Bare recovery proves exactly one supported worker dispatch and completion.**
Verified by: the retained selected-bare artifact or a fixture-backed equivalent passes without inventing an obsolete mode marker.

**AC-2 - Team recovery proves one logical worker dispatch.**
Verified by: the retained selected-team shape is classified correctly, while a real duplicate worker dispatch fails.

## Test plan

First replay both retained substrate artifacts. Add focused positive and duplicate-dispatch negatives. Run one final local substrate target only if fixtures cannot establish the host boundary.
