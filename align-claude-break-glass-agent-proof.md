---
title: Align Claude break-glass recovery proof with the current Agent tool
status: ideation
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
              resolution:
                type: Resolution
                id: resolution:spacedock:s2atxdv146qknetdjvx0xer6:backlog:1
                briefing: briefing:s2atxdv146qknetdjvx0xer6:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-11-03T18:27:44.543208Z"
                decision: approve
                reason: 'Captain directed in chat on 2026-11-03: ''for those 2, do we know how to fix it? if not dispatch local probe on top of the worktree.'''
              application:
                target-stage: ideation
                state: consumed
        - id: gate:s2atxdv146qknetdjvx0xer6:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:s2atxdv146qknetdjvx0xer6-ideation-1
              briefing:
                id: briefing:s2atxdv146qknetdjvx0xer6:ideation:attempt-1:revision-1
                digest: sha256:78fef8b36b99112c5b01923a6afb00b91b3d16a91612e3ecd1ce58826af7f869
                room-ref: '@review/ideation/briefing-1'
started: 2026-11-03T18:29:08Z
---

Restore a truthful break-glass substrate proof for the current Claude Agent tool surface.

## Problem

`TestLiveBreakGlassShimRecovery` grades two recovery arms. Replay of six retained Sonnet streams (PR #680 run 31640122346; run 31894934349 attempts 1 and 2) through the current oracle reproduces every CI verdict. It shows two different causes. Neither cause is an obsolete Agent tool shape.

1. The `selected-bare` arm never selects bare. It reuses `breakGlassShimPrompt()`. Commit 4cc0d890d removed the entity name from that prompt on purpose, so that the then team-shaped break-glass test would not select bare by accident. `TestDispatchRecoveryPromptsSelectIntendedDispatchModes` still forbids the name. Commit bbe3d7a05 then split the test into two arms and reused this prompt for the bare arm. With no named entity and no `/spacedock bare`, the Claude contract selects named background dispatch when `SendMessage` is present. No retained bare run passed `--bare-mode` to `dispatch build`. Sonnet chose by itself. PR #680 and attempt 1 chose team ("team mode, since SendMessage/Agent are available in this session") and failed. Attempt 2 chose blocking and passed. This is the same root cause as the selected-bare half of `repair-sonnet-live-flakes` (060xp). No selected mode is lost at recovery, so a fix to the adapter's bare-mode preservation text cannot help.
2. The bare template and the bare oracle depend on a host claim that is now false. `fo-dispatch-recovery` says an omitted `run_in_background` and `false` both mean blocking. The oracle accepts an omitted key as bare. Live probe on Claude Code 2.1.303 (Sonnet `claude-sonnet-5-5`): a name-less `Agent` call without `run_in_background` returned "Async agent launched" with `subagent_stats` `requested.unset=1, started_in_background=1`. With explicit `run_in_background=false` it returned the worker hand-back with `requested.foreground=1, started_in_background=0`. A literal template follower therefore gets a background worker, and the oracle grades it as bare.

The `selected-team` oracle is correct. In PR #680 the second `Agent` call was `Agent(name="spacedock-ensign-widget-task-implementation", description="Shut down implementation worker", prompt="{\"type\":\"shutdown_request\"}")`. The host spawned a second agent (`ListAgents`: "spacedock-ensign-widget-task-implementation [f3a40f] · general-purpose · running"). Each `Agent` call is a real spawn. The rule "exactly one Agent() call" counts this extra spawn correctly. Both later team runs delivered the shutdown with `SendMessage` and passed. That FO slip happened once in three runs, on an older contract. This task does not change the contract for it.

## Proposed approach

1. Fixture: add `breakGlassShimBarePrompt()` = `"The captain has issued the command: /spacedock bare\n\n" + breakGlassShimPrompt()`. The bare arm uses it. The team arm keeps its explicit team sentence. The base prompt keeps its "no entity name" guard. This uses the same selector as `TestLiveBareReachable`, which selected `--bare-mode` with a blocking `Agent` in all three retained runs.
2. Oracle (bare arm only): `modeShape = !hasName && !hasTeamName && hasRunInBackground && !runInBackground`. Remove the "omitted key means blocking" comment on `streamContentBlock`. The team arm and the one-spawn count do not change.
3. Contract, `skills/fo-dispatch-recovery/SKILL.md` bare arm. Before: "Omit `name`, `team_name`, and `run_in_background` entirely; Claude's observable tool stream may preserve that omission or normalize it to `run_in_background=false`, and both mean blocking bare dispatch. Never pass or accept `run_in_background=true`." After: "Omit `name` and `team_name`, and pass `run_in_background=false`: the Agent tool starts a call without that key in the background." Add `run_in_background=false,` after `description=` in the bare template block.

Mechanisms considered:
- Grade the host `subagent_stats` result counters. Not chosen: 2.1.229 and 2.1.233 do not emit them, and the input key is already proven against them. They stay as spike evidence only.
- Name the entity to enter single-entity mode. Not chosen: the existing guard forbids the name, and `/spacedock bare` is the explicit captain selector that mirrors the team arm.
- Keep the omitted-key acceptance. Rejected: on the current host that shape runs in the background, so the oracle would pass a mode violation.

No hooks, no instrumentation, no new oracle mechanism.

## Expected surface

Estimate net LOC change: +12, across 5 files (insertions about 27, deletions about 15). Tolerance: net +0 to +25, 5 files.

- `internal/ensigncycle/dispatch_recovery_fixtures_test.go` (+6/-0): bare prompt builder.
- `internal/ensigncycle/dispatch_recovery_live_test.go` (+1/-1): the bare arm uses the bare prompt.
- `internal/ensigncycle/dispatch_recovery_assert_test.go` (+3/-4): bare predicate and comment.
- `internal/ensigncycle/dispatch_recovery_assert_unit_test.go` (+15/-9): base bare fixture carries `false`; omitted key becomes a negative; team negatives and positives; guard extension.
- `skills/fo-dispatch-recovery/SKILL.md` (+2/-1).

Semantic changes: FO runtime behavior on the break-glass bare arm only (pass explicit `false`). Live-test fixture prompt for one substrate arm. No change to command grammar, stored formats, authority, normal `dispatch build` output, common live journeys, or Codex/Pi. No documentation-site text describes this template (`grep` of `docs/` for bare `run_in_background` omission finds nothing), so there is no doc diff.

Coordination: this task does not remove the `breakGlassSelectedBareGap` XFAIL bindings. The Sonnet binding is owned by 060xp (its AC-3 needs unbound CI PASS). The Haiku binding has no owner and no analysable artifact, because the `/tmp/sd-f7` art2/art7 trees contain zero files. The FO should either fold 060xp's selected-bare half into this task or let 060xp retire the binding after this lands.

Follow-ups, out of scope (normal dispatch): `dispatch build --bare-mode` omits `run_in_background`, `claude-fo-dispatch.md` says "omit when field absent", and `assertBareReachableObservables` treats omitted as false. On 2.1.303 a verbatim mapping starts a background worker. Sonnet added `false` itself in all three retained bare-reachable runs. Also, the ensign runtime still targets `team-lead`, which is unreachable. In every retained run the ensign retried to `main`.

## Acceptance criteria

**AC-1 - Selected-bare break-glass recovery dispatches exactly one blocking worker on the current Claude Agent tool.**
Verified by: live Sonnet `TestLiveBreakGlassShimRecovery/selected-bare` on the implementation head passes its observable and durable checks, and the stream's host `subagent_stats` show `foreground=1, started_in_background=0`. Baseline: 1 of 3 retained CI runs passed (by chance). Ideation spike on the proposed bytes: 2 of 2 passed with `--bare-mode` and a host-reported foreground spawn.

**AC-2 - The bare oracle accepts only the call shape the current host runs in the foreground.**
Verified by: an offline fixture with `run_in_background:false` and no name passes. Each of these negatives fails: key omitted, `true`, `name` present, `team_name` present, two calls. Reverting the predicate to accept an omitted key makes the omitted-key negative pass, so the test goes red.

**AC-3 - Team recovery counts every Agent() call as a spawn.**
Verified by: an offline positive (one named background call followed by a `SendMessage` shutdown_request) passes. A PR #680-shaped negative (worker call plus `Agent(name=<same>, prompt=shutdown_request)` without `subagent_type`) fails with "observed 2 Agent() calls". Replay of the six retained streams gives the CI verdicts: PR #680 bare FAIL, team FAIL; 31894934349 attempt 1 bare FAIL, team PASS; attempt 2 bare PASS, team PASS.

**AC-4 - Each break-glass arm selects its mode before the helper fails.**
Verified by: `TestDispatchRecoveryPromptsSelectIntendedDispatchModes` asserts that the bare arm's prompt carries `/spacedock bare`, that the team arm's prompt does not, and that neither names `widget-task`. Removing the selector from the bare prompt makes the test go red. The AC-1 live stream shows `dispatch build ... --bare-mode` before recovery.

**AC-5 - The recovery skill's bare template tells the FO to pass `run_in_background=false`.**
Verified by: AC-1's live stream shows the recovery `Agent` call with `run_in_background:false`. This AC is a means, paired with AC-1.

## Test plan

- Primary proof owner: `dispatch_recovery_assert_unit_test.go` (offline, deterministic, no spend) for AC-2, AC-3 and AC-4. Falsifying edits: revert the predicate (AC-2); exempt non-ensign `Agent` calls from the count (AC-3); drop the bare selector (AC-4).
- Live: one Sonnet `TestLiveBreakGlassShimRecovery` run (both arms, about 3.5 minutes) on the implementation head for AC-1 and AC-5. Use the commands in `docs/runtime-live-ci.md`. Note: `~/.claude/benchmark-token` is revoked ("OAuth token revoked"). The runner's macOS keychain stored-login seed authenticated the spike.
- `go test ./...`, `go test ./... -race`, `gofmt`.

## Spike record

Probe base: detached `live-lane-shared-failure-fixes` at 64274cef5ba2b0c678f5d2882ed588e6ba8e583c, with throwaway edits of items 1 to 3 above. Results:
- Host probe on 2.1.303: omitted key gives background; `false` gives foreground.
- Replay with a throwaway Go test calling `assertBreakGlassObservables`: the six verdicts above, unchanged under the new predicate.
- Live selected-bare: 2 of 2 PASS (XPASS against the 060xp binding).
- With the new predicate, only `TestAssertBreakGlassObservablesOffline/selected_bare_omits_run_in_background` went red, as expected. That case becomes a negative.

## Stage Report: ideation

- DONE: With local Sonnet probes on a throwaway checkout of the live-lane-shared-failure-fixes branch head, record the actual Agent tool call shapes in both break-glass recovery variants (selected-bare and selected-team), including any duplicate call, and compare them with what the break-glass oracle in dispatch_recovery_live_test.go accepts.
  Probe base is 64274cef5. Two local Sonnet selected-bare runs each made one `Agent(subagent_type, description, run_in_background:false)`, host foreground=1. Retained team calls and the PR #680 duplicate (Agent-delivered shutdown_request spawning a general-purpose agent) are in Problem. The oracle is in `dispatch_recovery_assert_test.go`.
- DONE: Decide whether the defect is in user behavior (contract) or in the oracle, and propose the smallest correction with focused positive and duplicate-dispatch negative cases and a net LOC estimate.
  Oracle side: the bare-arm fixture never selects bare. Contract and oracle: both claim an omitted `run_in_background` is blocking, which 2.1.303 disproves. The team oracle is correct. Net +12 across 5 files, with positives and negatives in AC-2 and AC-3.
- DONE: Replay the retained artifacts (PR #680 run 31640122346 if still downloadable, and /tmp/sd-f7 art2/art3/art7 break-glass streams) as fixture-backed evidence before relying on new live runs.
  Replayed PR #680 and both attempts of run 31894934349 (060xp source) through `assertBreakGlassObservables`. All six CI verdicts reproduce. The /tmp/sd-f7 art2/art3/art7 trees hold 0 files (directories only), so they could not be replayed.

### Summary

Selected-bare fails because its prompt deliberately carries no bare selector (since 4cc0d890d). Sonnet picks team or bare by chance (1 of 3 retained passes). 060xp has the same root cause. Separately, the recovery template and the bare oracle treat an omitted `run_in_background` as blocking. A Claude Code 2.1.303 probe shows that shape runs in the background. Proposal: an explicit `/spacedock bare` prompt for the bare arm, an explicit `false` in the template and the oracle, and no change to the team oracle. The PR #680 second Agent call was a real extra spawn. Spike: 2 of 2 live passes. The benchmark token is revoked, so the keychain seed was used.
