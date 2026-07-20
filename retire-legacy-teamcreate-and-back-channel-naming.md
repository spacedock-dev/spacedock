---
id: q4pw3xb4nf4cwfdjtwbn17mz
title: Retire legacy TeamCreate path and rename back-channel to inter-agent communication
status: backlog
source: captain (CL), 2026-07-20 session
started:
completed:
verdict:
score:
worktree:
issue:
---

Retire the legacy TeamCreate machinery from the shipped FO contract and stop naming the worker/FO messaging surface "back-channel".

## Problem

Two contract-surface cleanups, captain-directed:

1. **Legacy TeamCreate path.** `claude-fo-dispatch.md` carries a legacy-override line that makes every session's first dispatch probe `ToolSearch(select:TeamCreate)`; when it matches, the FO loads `skills/using-legacy-claude-team`. Current runtimes no longer expose TeamCreate (this session's probe: no match), so the probe is dead weight on the first-dispatch path, and `using-legacy-claude-team/SKILL.md` is one of the 13 FO-surface measured files — retired, it shrinks the measured prompt surface. The claude runtime adapter and the terminal-teardown prose also carry legacy-override references that ride only on this trigger.

2. **"Back-channel" naming.** The contract names inter-agent messaging "back-channel" throughout (`## Worker Back-Channel` in claude-fo-dispatch.md, Codex "mailbox back-channel", Pi "contact_supervisor/intercom back-channel", «addressable-worker» prose). Captain direction: do not name it "back-channel" — name it for what it is, e.g. "inter-agent communication" (final term chosen at ideation).

## Proposed approach

{Ideation fills this in. Expected shape: delete the legacy-override line and the using-legacy-claude-team skill; update the FO-surface file list and reference-closure expectations in contractlint; rename the concept across claude-fo-dispatch.md, fo-dispatch-core.md capability prose, and the codex/pi adapter bullets. Both ratchets must account for the surface shrink honestly.}

## Out of scope

{Ideation fills this in. Likely: no behavior change to dispatch or messaging itself — naming and dead-path removal only.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - No TeamCreate probe or legacy-team path remains in the shipped contract; first dispatch has no legacy branch.**
Verified by: contractlint reference-closure/structural checks after the skill's removal; the FO-surface measured-file list no longer contains `using-legacy-claude-team/SKILL.md`; review-time grep cited as evidence.

**AC-2 (VALUE) - The FO prompt surface shrinks by at least the retired skill's byte size.**
Verified by: `TestFOFunctionPromptSurfaceShrinks` / checkpoint-metrics byte counts before vs after, an independent number that can move the wrong way.

**AC-3 - The shipped contract nowhere names inter-agent messaging "back-channel"; the concept carries a descriptive name.**
Verified by: review-time grep over the FO surface cited as evidence (no committed prose-grep, per proof policy).

**AC-4 - The FO-surface ratchet baseline is re-tightened to the post-retirement measured value (zero-slack convention).**
This task carries the F1 mitigation from `codex-post-compaction-contract-reload` (archived, validation cycle 7): with 204 B of slack at baseline 122634, the scratch-proven inert-heading escape (122584 B) passes every mechanical gate, leaving the review-time grep as the only guard for inert-prose reintroduction. Re-tightening here closes it mechanically.
Verified by: `foFunctionReferenceBaselineBytes` equals measured+1 after the retirement lands; `TestFOFunctionPromptSurfaceShrinks` green at zero slack.

## Test plan

{Ideation fills this in. Expected: contractlint structural suites + surface-metrics deltas; claude-live lane required (shipped contract change); no committed grep tests.}
