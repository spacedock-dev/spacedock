---
id: 4ctz4sfybfk0mfsbnrjcc7bv
title: Force one boot at the compaction boundary — a compacted FO resumes on stale bindings
status: backlog
source: "Captain CL, 2026-08-18, in chat: 'file the compaction improvement, detail the problem diagnosis and proposed solution.' Raised after the FO opened three PRs by hand without reading the pr-merge mod, and diagnosed the root cause as never having run Startup in a compaction-resumed session."
started:
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:4ctz4sfybfk0mfsbnrjcc7bv:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:4ctz4sfybfk0mfsbnrjcc7bv-backlog-1
              briefing:
                id: briefing:4ctz4sfybfk0mfsbnrjcc7bv:backlog:attempt-1:revision-1
                digest: sha256:8c097b36dc8addb17335f90ea018c35bb0f1b083f31bce9d390c80df9e543608
                request-digest: sha256:79ea2c142d7fb3cbbd6ec4afb6b981efee79195296b7750c8b3b4166550a49f9
                room-ref: ./force-boot-at-compaction-boundary/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:4ctz4sfybfk0mfsbnrjcc7bv:backlog:1
                briefing: briefing:4ctz4sfybfk0mfsbnrjcc7bv:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-18T23:14:57.630325Z"
                decision: approve
                reason: 'Captain approved in chat: ''ok dispatch 4c and nr.'' Accepts the seed into ideation; the captain independently identified this as the highest-value item on the board because it makes every future session more reliable rather than fixing one test.'
              application:
                target-stage: ideation
                state: pending
---

A compacted First Officer session inherits the previous session's **narrative** but not its **bindings**. The transcript summary preserves what the FO was discussing and drops what it was standing on: which binary, which contract version, which mods are registered, which workers are alive, and what durable state actually says.

## Problem — diagnosis

The FO contract's Startup runs a binary gate, one `«state.boot»()`, and `«engage»`, which runs `state ready`, `state sweep`, and `«hooks.run»("startup")` exactly once. A compaction-resumed session does not know it is a new session. It believes it is a continuation, so none of that runs.

Six failures in one session (2026-08-18), all traceable to that single omission:

1. **Binary gate never ran.** The first command failed against `/opt/homebrew/Caskroom/spacedock@next/0.27.0-pre7/spacedock`, deleted by the pre8 cask upgrade. Startup step 1 resolves the launcher before anything else.
2. **The `mods` map was never read**, because no boot record was obtained. The FO then opened three GitHub PRs by hand, bypassing the `pr-merge` mod completely: no gate recorded, no `pr:` field, no `mod-block`, no merge guard armed. Three PRs existed that the workflow could not see. Repaired only when the captain asked the FO to check the mod.
3. **`state ready` / `state sweep` never ran**, so the session never converged with the remote state checkout.
4. **`«hooks.run»("startup")` never fired**, so `pr-merge`'s PR-pending scan and `comm-officer` injection never happened.
5. **The session ran the `0.27.0-pre1` contract while `pre8` was installed**, so `fo-gate-lifecycle` was absent from its loaded set. Four gates were "presented" in chat and never `gate prepare`d. Durable readiness stayed `needs-preparation` for hours while the FO reported them to the captain as awaiting decision.
6. **An orphaned worker kept running** — the ideation ensign for a withdrawn entity. A roster reconcile at startup would have reclaimed it.

Item 5 is the sharpest: the FO's prose and durable state disagreed, and only the binary knew. That is the same defect class as two entities already in flight — `conn-delegated-approver-attribution` (a grader that trusted the agent's self-description) and `merge-guard-requires-preceding-report` (an FO that trusted a completion message over durable state) — one level up. Here the FO trusts **its own summary** over the binary.

## Why a prose rule will not fix this

The gate-attribution convention has lived in `skills/fo-gate-lifecycle/SKILL.md` since `aa04e95d8` (2026-07-25). Live First Officers recited it in their own final messages and still stamped `person:captain` — recorded in this repo's finding-9 audit note. Prose in the contract demonstrably does not hold against a behavior the session is already inclined toward.

Worse, the compaction summary carries an explicit counter-instruction: *"Resume directly — do not acknowledge the summary, do not recap, do not preface."* That is correct for conversational continuity and wrong for a stateful dispatcher. Any mechanism here must survive an instruction telling the session to keep going.

## Proposed approach

Bind exactly **one `«state.boot»()`** to the compaction boundary, before any other action. Nothing else.

The boot record already carries `mods`, `ready_gates`, `dispatchable`, `pr_state`, `team_state`, `state_backend`, `definition_dir`/`entity_dir`, and the binary version gate. Every one of the six failures above is answered by that single call. The contract does not need re-injecting; the state needs re-reading.

### Relationship to issue #595

#595 is this same problem filed for Codex, and it reaches the same diagnosis in its own words: a teammate must retain contract and authority "without trusting a lossy transcript summary." It records that the Codex surface exposes no compaction callback, and proposes a digest-pinned bundle of contract, adapter, README, entity state, worker identity, gates, PR, and CI.

This entity is the Claude-side instance and deliberately proposes far less. If Claude exposes a usable boundary, one boot call ships without any bundle. #595's bundle remains the right shape only where no callback exists at all.

## Ideation must settle

1. **Whether Claude Code exposes a usable compaction boundary, proven by a captured event rather than documentation.** This repo's own `hook-events.jsonl` contains exactly one `PreToolUse` record and no compaction evidence, so nothing local proves it today. Prove the boundary first; design second.
2. **Whether a hook can force an action before the next tool call, or only inject text.** Injected text is prose and fails for the reason above. If text is all the host offers, say so and record the boundary as unsupported — do not ship a prose rule dressed as a mechanism.
3. **Where the mechanism belongs** — a new mod hook point, the runtime adapter, or the binary. Prefer the smallest that holds.
4. **Whether a host-independent fail-closed tell can substitute.** If the FO must be able to prove it booted — for example by holding a boot record whose identity matches the current session — then a session that cannot prove it booted refuses gate and merge actions until it does. This needs no host callback and would cover Codex and Pi as well. Weigh it against the hook approach on necessity, not preference.

## Out of scope

The Codex boundary — #595 owns it. Any context-reinjection bundle. Post-compaction worker-identity recovery. Changing what Startup does; this entity changes only whether it is forced to run.

## Value

The First Officer mutates durable state, resolves gates, and drives merges. A compacted session doing that work on stale bindings is the highest-authority actor in the system operating on a summary it wrote about itself. The value is that the FO cannot act on workflow state it has not re-read.
