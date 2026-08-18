---
id: nr75fq7ha3nmvsegbd22cgqa
title: Finish the legacy Claude TeamCreate retirement — the contract retired it, the binary and a live proof did not
status: backlog
source: "Captain CL, 2026-08-18, in chat ('file a proper legacy claude team retirement task'), after TestLiveBreakGlassShimRecovery/selected-team failed the claude-live lane on the ca9/6ht/j7j stack by demanding a team_name the shipped contract tells the FO to omit."
started:
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:nr75fq7ha3nmvsegbd22cgqa:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:nr75fq7ha3nmvsegbd22cgqa-backlog-1
              briefing:
                id: briefing:nr75fq7ha3nmvsegbd22cgqa:backlog:attempt-1:revision-1
                digest: sha256:0bb3e4f258b019130d5c0bd47aad6f12bf8064cf9b30f3e9080524bf89a00f14
                request-digest: sha256:0fd53f31d5636406dca65acd9596814d0da6fb9c15c45a6fb2e57882b3c8827d
                room-ref: ./finish-legacy-claude-team-retirement/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:nr75fq7ha3nmvsegbd22cgqa:backlog:1
                briefing: briefing:nr75fq7ha3nmvsegbd22cgqa:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-18T23:15:14.692985Z"
                decision: approve
                reason: 'Captain approved in chat: ''ok dispatch 4c and nr.'' The captain raised the retirement question himself after the break-glass red, and asked for this filing; approving it into ideation closes a nondeterministic merge blocker.'
              application:
                target-stage: ideation
                state: pending
---

`ecffcedef` (#549) retired the legacy TeamCreate dispatch path and deleted `skills/using-legacy-claude-team/`. The retirement stopped at the skill boundary. The binary still carries the flag and the field, and a live proof still fails a First Officer for obeying the current contract.

## Problem

The shipped Claude adapter now states the post-retirement shape plainly (`skills/first-officer/references/claude-fo-dispatch.md:7`):

> a worker is `Agent(name=…, run_in_background=true)` **(no `team_name`)** … `spacedock dispatch build` emits this shape (`name` present, `team_name` absent, `run_in_background` true)

and at `:26` it marks the field `// absent in the normal shape; map it verbatim if the build emits it`.

Three things still contradict that:

1. **The binary keeps the legacy surface.** `internal/dispatch/build.go` still declares `TeamName *string \`json:"team_name,omitempty"\`` and `teamNamePattern`, and `dispatch build --help` still advertises `--team-name NAME  Select the legacy TeamCreate-registry dispatch shape`.
2. **A live proof demands the retired shape.** `TestLiveBreakGlassShimRecovery/selected-team` (`internal/ensigncycle/dispatch_recovery_live_test.go:116`) fails with *"the only Agent() call did not preserve selected team mode"* when the FO emits `team_name-present=false` — which is exactly what `claude-fo-dispatch.md` instructs. Its sibling `selected-bare` is already `liveXFail("claude-sonnet", …)` against `repair-sonnet-live-flakes`; `selected-team` carries no registration and hard-fails through `t.Fatalf`.
3. **It blocks merges nondeterministically.** That test reddened the `claude-live` lane on the `ca9`/`6ht`/`j7j` stack (run 32189720211) while the same commit passed the same test locally, and `main` passed it both in CI and locally. The captain merged the stack over it.

## Why this matters beyond one flaky test

This is the third live-lane test in one day that reddened a **compliant** agent. `ca9` fixed a filing recognizer that missed a newline-terminated command. `6ht` fixed a build-count bar that punished a corrective rebuild. This one fails an FO for emitting the shape its own shipped contract prescribes.

A grader that punishes the documented behavior is worse than no grader: it trains the reader to treat live reds as noise, which is exactly the habit that lets a real red through.

## Ideation must settle

1. **Whether legacy team mode is retired in fact or only in the skills.** Decide one way. If retired, `--team-name`, `TeamName`, `teamNamePattern`, and the `selected-team` live case go together, and `dispatch build --help` stops advertising a mode we do not ship. If it is a supported fallback, say what supports it and where that is written, because `claude-fo-dispatch.md` currently says the opposite.
2. **What `selected-team` should become.** Delete with the mode, re-anchor to assert the post-retirement shape (break-glass preserves `name` and `run_in_background` and emits no `team_name`), or register it under `repair-sonnet-live-flakes` beside its sibling. Deleting a live proof needs the captain's explicit approval; recommend, do not assume.
3. **Whether the break-glass path has a real defect underneath.** The failure is nondeterministic across otherwise identical runs, so before retiring the assertion, establish whether the FO's shape choice under break-glass is stable at all. A retirement that hides a real instability is not a retirement.
4. **What else `#549` left behind.** Sweep for other survivors of the same retirement — help text, fixtures, mod prose, adapter references, `dispatch reconcile --team-name`.

## Out of scope

`repair-sonnet-live-flakes` and its existing xfail registrations. The auto-team dispatch shape itself. Any change to break-glass recovery behavior beyond what settling item 3 requires.

## Value

A live lane is only worth its cost if a red means something. Today one of its tests fails agents for following the contract we ship, and it has already forced a merge decision over a red. Either the mode exists and the contract is wrong, or the mode is gone and the test is. Both are cheap; the ambiguity is not.
