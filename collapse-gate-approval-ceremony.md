---
id: 7fhzvvk8d5smj858bp47xbjq
title: Collapse the gate-approval ceremony from ~16 tool calls to 1-2
status: ideation
source: "Ships-counselor friction rollup, 2026-08-02, theme 1 (gate-approval ceremony): measured on sonnet-gate-guardrail-no-authority's ideation->implementation gate -- 16 discrete FO tool calls and 156s wall clock to apply one captain word ('approve'), with 0 additional captain turns needed. Recurs at every nonterminal gate for every entity, forever. Captain directed: file and dispatch (ideation, via a fable-model ensign)."
started:
completed:
verdict:
score: 0.6
worktree:
issue:
gates:
    version: 1
    current:
        gate: gate:7fhzvvk8d5smj858bp47xbjq:backlog
    records:
        - id: gate:7fhzvvk8d5smj858bp47xbjq:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:7fhzvvk8d5smj858bp47xbjq-backlog-1
              briefing:
                id: briefing:7fhzvvk8d5smj858bp47xbjq:backlog:attempt-1:revision-1
                digest: sha256:2cb1660e6fb16e9c6303564fc34d957a34155eda05ed8646a06b770a92ed221c
                digest-domain: canonical-bytes
                request-digest: sha256:a9c5b97c54d3372732e38b8ae7157af41a6aaca2dae2cdecece420ab2fd0572a
                room-ref: ./collapse-gate-approval-ceremony/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:7fhzvvk8d5smj858bp47xbjq:backlog:1
                briefing: briefing:7fhzvvk8d5smj858bp47xbjq:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-02T06:25:11.852955Z"
                decision: approve
                reason: 'Captain directed in chat: file the gate-approval-ceremony friction as a task and dispatch it to ideation via a fable-model ensign.'
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

Applying one captain gate decision (approve/revise/hold) currently costs ~16 raw FO tool calls and ~2.5 minutes of wall clock, most of it mechanical: re-invoking `state commit` after nearly every binary call, then a separate set of frontmatter stamps and worktree creation before the next dispatch even begins. This entity is about collapsing that ceremony toward 1-2 calls per gate decision without weakening any of the authority/integrity checks the ceremony exists to enforce.

## Problem

Measured this session (`sonnet-gate-guardrail-no-authority`, ideation -> implementation gate, exact tool-call sequence from the transcript): `status --set` (retitle) -> `gate record --decision approve` -> `state commit` -> `gate consume` -> `state commit` -> `status --read` (verify) -> `status --set` (worktree stamp) -> `state commit` -> `git status` -> `git worktree add` -> 2x `Write` (checklist/scope-notes scratch files) -> `ToolSearch` -> `SendMessage` (prior-worker shutdown) -> `dispatch build` -> `Agent` (the actual spawn). 16 calls, 156s, 0 additional captain turns, to apply the word "approve".

3-4 of the 16 are the same `state commit` re-invoked idempotently right after `gate record` and again after `gate consume` -- each of those binary commands already knows it just mutated durable state. The remaining bulk is the post-consume "ordinary dispatch" procedure (frontmatter stamps, worktree creation) that always follows a nonterminal consume in the same fixed shape.

This is not a one-off cost: every entity, at every nonterminal `gate: true` stage, pays this in full. A workflow that runs 50 entities through 2 gates each pays it 100 times.

## Proposed approach

{Ideation fills this in -- explore and choose among (or combine) these candidates, in ascending order of mechanism weight; do not assume the first is sufficient without checking:}

1. **Make `state commit` implicit inside the mutating gate verbs.** `gate record` and `gate consume` already know they just durably mutated `gates:`/`status` frontmatter -- have each of them commit+push atomically as part of the same command, removing the separate `state commit` call after each. Cheapest, most mechanical, smallest surface change.
2. **A `spacedock gate approve <slug> --reason TEXT --workflow-dir DIR` convenience verb** that composes record(decision=approve, actor=person:captain) + consume + the commits from (1) into one call, for the common "captain approved in chat" path. Needs a matching `revise`/`hold` shape or an explicit `--decision` flag; keep the existing `record`/`consume` primitives for the room-backed and delegated-actor paths that don't fit this shortcut.
3. **A `spacedock dispatch advance <slug> --workflow-dir DIR` verb** folding the post-consume "ordinary dispatch" bookkeeping -- the `started`/`worktree=` frontmatter stamp and `git worktree add` for a worktree-declaring next stage -- into one call, so the FO doesn't hand-run four more commands between a successful consume and the actual `dispatch build`/spawn.

Consider whether (2) and (3) should compose into a single `gate approve --and-dispatch` verb, or stay separate primitives the FO chains -- this is a real design choice ideation should make explicit, not default to "more verbs."

## Out of scope

Any weakening of gate authority, digest/Briefing integrity checks, frozen-binding refusals, or the captain-vs-FO actor distinction. Any change to `gate prepare`'s semantics or the room-backed Result path. The mechanism-2 byte-cap follow-up (`gate-lifecycle-hardening-byte-budget`) -- unrelated. Batch/bulk gate operations across multiple entities -- this is about one entity's one gate decision.

## Acceptance criteria

**AC-1 (VALUE) - A captain "approve" decision at a nonterminal `gate: true` stage collapses from the measured baseline toward 1-2 FO-issued commands.**
Independent baseline: 16 (this entity's own `## Problem` section, sourced from the sonnet-gate-guardrail-no-authority session transcript). Verified by: a fixture or live-workflow test that counts the FO-issued commands (not internal binary subprocess calls) needed to take one entity from an open gate through a successfully dispatched next stage, before and after this change, on the same fixture entity/gate shape.

**AC-2 - No loss of gate authority or integrity guarantees.**
Verified by: existing gate-lifecycle tests (frozen-binding refusal, digest mismatch rejection, actor/authority validation, room-backed vs chat-decision paths) pass unchanged; any new convenience verb is proven to compose the exact same underlying operations (record+consume, or the dispatch-advance stamps) rather than a parallel shortcut implementation.

**AC-3 - The FO contract (skills/first-officer, skills/fo-gate-lifecycle) is updated to use the new call shape wherever it currently documents the longer sequence.**
Verified by: `references/fo-gate-lifecycle` (or wherever the gate lifecycle prose lives after this change) and `references/fo-dispatch-core.md` reflect the collapsed call sequence; no stale prose still instructs the old multi-call ceremony as the primary path.

## Test plan

Reproduce the 16-call baseline as a fixture-driven count first (a scripted FO-simulation harness or a recorded transcript replay), so before/after has a real number, not an estimate. Go unit/fixture tests for the new command(s). No live workflow run needed unless the FO contract prose change is judged high-risk enough to warrant one.
