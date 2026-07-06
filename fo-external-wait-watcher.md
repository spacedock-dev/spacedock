---
title: "Dispatch core names the external-wait watcher pattern (no turn-end on a signal with no sender)"
status: backlog
source: "FO Commander session 2026-07-06/07 (0250 drive): after dispatching PR #475's live lanes, the FO ended its turn 'holding' for CI results — but GitHub CI emits no teammate event, so nothing could wake it; the captain prodded twice ('well?', 'were you idling?'). Root cause: claude-fo-dispatch.md's Awaiting Completion states 'the runtime will wake you when a real event arrives' as if universal; it holds for ensign SendMessage/task notifications, not for CI lanes / PR merges / any external async process. The correct mechanism (a run_in_background watcher, e.g. gh run watch, whose exit re-invokes the FO) existed but no contract line names the pattern. Posture half (never end turn with actionable async work) is 0250 vcm's approved scope; this task is the Claude-adapter mechanism half."
started:
completed:
verdict:
score: 0.4
worktree:
issue:
id: 3jbyvy0j94bx2nvxt9j9vs17
---

The FO contract's waiting guidance covers exactly one async channel — dispatched workers, whose completion genuinely wakes the FO — and its "end the turn empty; the runtime will wake you" rule gets over-generalized to channels with no event source (CI lanes, PR merge state, external commands). Proposed direction: a clause in the runtime-neutral dispatch core (fo-dispatch-core.md, beside the event-loop/awaiting guidance) naming the class "async external signal with no runtime event source" and its general mechanism: launch a bounded background terminal command as a watcher (e.g. `gh run watch <id> --exit-status`, `gh pr checks --watch`) whose process exit re-invokes the FO, instead of ending the turn to wait or sleep-polling inline. Every host has background terminal commands — the pattern is not host-specific; per-host invocation shapes ride the adapters' `→` lines per the contract's existing host-binding convention. Acceptance sketch: value — a live FO drive with a pending CI check reaches the check's resolution with zero captain prods (baseline: the 2026-07-06 incident, two prods); mechanism — the core clause ships and a live drive observes the watcher being armed. Coordinate with 0250 vcm (posture) to avoid duplicated resident bytes.
