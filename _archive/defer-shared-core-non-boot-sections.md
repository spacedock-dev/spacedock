---
id: s6qamkh7efky9zh5jh6ba6xq
title: Defer non-boot-reachable shared-core sections to their trigger points
source: boot-forensics profile of FO session afd74765, 2026-07-25
status: backlog
started:
completed:
verdict:
score:
worktree:
issue:
---

The first-officer shared core loads whole at boot (~10.1k tok of a 40.4k-tok cold boot). Roughly a third of it cannot fire until a trigger that already owns a deferred load point, so a greet-and-stop session - the most common interactive shape - pays for completion, gate, feedback, and hook prose it can never reach. Move those sections behind their triggers, measured by boot occupancy in tokens.

## Problem

Measured on session afd74765 (2026-07-25, Claude Code, `spacedock 0.26.0+dev`):

- Cold baseline 23,721 tok; occupancy after a greet-and-stop boot 40,450 tok (x1.71 growth).
- `first-officer-shared-core.md` is 26,092 bytes ~= 10,120 tok - 58% of the 12,826-tok turn that loaded the contract, and the largest recurring static load in the session.
- The deferral machinery already works: the dispatch, write, merge, status, and dispatch-recovery cores all stayed unread at a greet.

The gap is that the boot-resident core is not itself subject to that discipline. Sections unreachable at a greet-and-stop boot, with measured byte cost:

Riding the existing first-dispatch trigger (5,331 bytes ~= 2,035 tok) - none can fire before a worker completes, and `fo-dispatch-core.md` is loaded by then, which the shared core already asserts for its reuse-or-fresh prose:

1. `## Completion and Gates` - 2,678
2. `«gate.ac-cross-check»` - 1,046
3. `«gate.assemble-verdict»` - 844
4. `«feedback.route»` - 763

Needing an «engage»-time trigger (2,572 bytes ~= 980 tok) - first reachable when «engage» runs `state ready` and the startup hooks:

1. `## State Management` - 692
2. `«hooks.run»` - 640
3. `«halt.rebase-conflict»` - 617
4. `## Mod Hook Convention` - 450
5. `«state.commit»` - 173

Narrow trigger (1,202 bytes ~= 460 tok):

1. `## Compaction continuity` - 698, mid-session only
2. `## Single-Entity Scope` - 374, headless scoped runs only
3. `## Working Directory` - 130, dispatch/worktree only

Candidate span totals 9,105 of 26,092 bytes - 35% of the boot-resident core.

Open question for ideation, not a recommendation: `## Working Principles` is 6,301 bytes (24% of the file), the largest single section. It shapes judgment from the first decision, so it has a real claim to residency; the counter-argument is that its gate/merge/triage bar is only exercised at those boundaries. Ideation should rule on it explicitly rather than leave it unexamined.

Two honesty notes carried from the profile. This is a deferral, not a reduction: a session that dispatches, gates, and merges pays the same total tokens, just later - the win is confined to the greet-only, question-only, and gate-only shapes, and the task should say which shapes it is optimizing. And the profile's first pass claimed a greet uses "maybe a quarter" of the core; the measured answer is that about two-thirds is boot-reachable, so the ceiling here is ~35%, not ~75%.

## Proposed approach

{Ideation fills this in. Starting hypothesis: append the first-dispatch group to the existing dispatch-core read so it costs zero extra round trips, and give the «engage» group one new trigger; every moved block leaves a one-line pointer at its trigger in `## Deferred load points`. Pointer overhead and any added Read round trip count against the measured delta rather than being excluded from it.}

## Out of scope

- Any change to what the FO does. This is load timing only; every section keeps its current semantics and trigger conditions.
- The harness-injected `deferred_tools_delta` load (~640 tok of MCP tool names) - not ours to defer.
- Orphan-worktree cleanup, which would shrink the boot JSON by ~940 tok. Tracked separately.
- Reducing the 23,721-tok cold baseline (system prompt, CLAUDE.md files, tool schemas).

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - A greet-and-stop FO boot occupies measurably fewer context tokens than the recorded 40,450-tok baseline.**
Verified by: re-running the boot-forensics per-turn occupancy extraction on a fresh greet-and-stop session after the change, reported in tokens against the 40,450 baseline. The number lives outside every file this task edits and can move the wrong way - added pointer prose or one new boot-time Read would raise it. Falsifying change: any split whose pointers restate the moved content, or that lands an extra boot-time read, leaving occupancy flat or higher.

**AC-2 - No section is moved behind a trigger it can be reached before.**
Verified by: {ideation designs this. It must not be a prose-grep over the contract - per the Proof policy a match over an instruction file asserts only that we wrote what we wrote. `internal/contractlint` structural absence plus reference closure is the sanctioned static form; the behavioral form is a live-lane run driving an FO through boot, engage, dispatch, completion, and gate with no missing-section failure.}

**AC-3 - {Mechanism AC, paired to AC-1: the candidate sections load at their trigger instead of at boot. Ideation states its falsifying change and the value AC it serves.}**

## Test plan

{Ideation fills this in. Must name the required host live lanes: the diff touches `skills/**/references/**`, so per the Proof policy this is the shipped-contract high-stakes surface - the detached adversarial audit applies, and every host lane whose adapter or host-neutral core is touched is required green rather than optional.}

## Superseded

Folded into `de-lecture-and-defer-fo-contract` (wn3dg7txnrte0jrcxf56b859) by captain order, 2026-08-24 ("file and fold s6q"). That task carries this body's deferral groups, measured profile, and AC-1 verbatim, plus the contract de-lecture cut list from the same day's FO audit. Verdict left empty per the supersede convention.
