---
id: vtgw5b1ar42s6x807nfdzwb9
title: Three shared live-lane failures that are not model faults
status: implementation
source: "Captain, 2026-10-27: Haiku analysis after PR #853 (runs 37703335512, 40448057053, 40447244422) found four of nine Haiku journey failures are grader bugs or cross-model failures. Captain: file the targeted fixes and get them passing locally."
started: 2026-10-27T04:56:52Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-live-lane-shared-failure-fixes
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:vtgw5b1ar42s6x807nfdzwb9:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:vtgw5b1ar42s6x807nfdzwb9-backlog-1
              briefing:
                id: briefing:vtgw5b1ar42s6x807nfdzwb9:backlog:attempt-1:revision-1
                digest: sha256:481cb8d061369714414c8a8464c4d84fb6dd9483d2bbb40976d16f250c72edff
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:vtgw5b1ar42s6x807nfdzwb9:backlog:1
                briefing: briefing:vtgw5b1ar42s6x807nfdzwb9:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-10-27T04:56:38.965708Z"
                decision: approve
                reason: 'Captain directed in chat: ''File a task with those targeted fix, and get them passed locally'' — approves starting design on the three fixes.'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:vtgw5b1ar42s6x807nfdzwb9:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:vtgw5b1ar42s6x807nfdzwb9-ideation-1
              briefing:
                id: briefing:vtgw5b1ar42s6x807nfdzwb9:ideation:attempt-1:revision-1
                digest: sha256:af9e708d3f7c3b5a59fc82cced1e0f7c816395c612b60fc487dcc99c0bb31cf8
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:vtgw5b1ar42s6x807nfdzwb9:ideation:1
                briefing: briefing:vtgw5b1ar42s6x807nfdzwb9:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-10-27T05:38:20.47712Z"
                decision: approve
                reason: 'Captain approves the ideation design: fixes 1, 2a+2b, 3a+3b with net +75 across 37 files (tolerance +40..+120 / 32..42 files), keeping the binary fix 2a.'
              application:
                target-stage: implementation
                state: consumed
---

Fix three live-lane failures that hit Sonnet and Haiku alike, so live results show real model behavior and not harness or contract faults.

## Problem

1. **Filing grader false negative.** `commandFilesViaNew` in `internal/ensigncycle/shared_filing_test.go` accepts only a space, tab, quote, newline, or end of line after the slug. A correct `${SPACEDOCK_BIN:-spacedock} new wire-the-thing; echo "exit=$?"` fails with `filing-command-not-observed`. Both Haiku runs did exactly this (10-07 and 10-27), and the entity landed with id 001. Sonnet passed only because its command ended with `--workflow-dir . ;`.
2. **keep-moving-posture: merge guard refuses the approved-gate entity.** After the captain-approved `approved-gate` entity finishes implementation and reaches the ungated terminal stage, `spacedock merge guard` reports it "ineligible" (no pending terminal approval), so it is never archived. The Sonnet and Haiku streams show the same refusal. This is the failure behind Sonnet's XFAIL (owner `060xp69y61yhrww23g3wvwqy`) and Haiku's blank-owner XFAIL. It is not yet known whether the defect is in the binary's merge guard (an entity that was gated earlier, then reached an ungated terminal) or in the FO contract's terminal path.
3. **Break-glass `selected-team` fails on every model, and the `team-lead` completion target is unreachable.** `TestLiveBreakGlassShimRecovery/selected-team` fails on Sonnet (PR #853 run 40447244422 and the `kent/compass-poc` run 39608837699) and on Haiku. Every worker-dispatching scenario, for both models, logs `No agent named 'team-lead' is reachable`. The completion target is pinned to `team-lead`, but the current Claude Code resolves only `main`. The suspected link between the two is unproven.

## Diagnosis

Evidence: CI streams in `/tmp/sd-f7/art{2,3,7}` (runs 40448057053, 40447244422, 37703335512) and 14 local Sonnet live runs on this checkout (`661c607ae`), listed under Risk evidence.

**Item 1 (grader, test code).** Confirmed. Fed through `commandFilesViaNew`, the recorded Haiku command `… | ${SPACEDOCK_BIN:-spacedock} new wire-the-thing; echo "exit=$?"` returns false. The ` && …` and ` --workflow-dir . ;` forms return true. The direct-launcher regex in `internal/ensigncycle/shared_filing_test.go:53` ends the slug with `[ \t'\n]|$` only, so a `;`, `&`, `|`, `)` or `"` that directly follows the slug fails. The var-capture branch already splits on separators and is not affected.

**Item 2 (binary, plus FO contract).** Two shared defects and one residual model choice:

1. Binary defect, `internal/status/merge.go` `consumedNonterminalHistory`. It treats a consumed nonterminal approval as history only when `currentStatus == app.TargetStage`. `status --set <slug> status=done verdict=PASSED completed` is allowed on that entity (it classifies the pre-update status `implementation`). Then `merge guard` classifies the post-update status `done` and refuses with `condition "ineligible"`. A gate-less entity in the same state finalizes. Reproduced deterministically with the existing `consumedNonterminalWorkflow` fixture: `--set` exits 0, then the guard exits 1 with the exact CI message. The same guard run from `implementation` exits 0. Both Haiku CI streams (10-07, 10-27) and 3 of 3 local Sonnet baseline runs took this path. The predicate also refuses an entity that moved on through a second ungated stage, which is the same class of defect.
2. FO contract defect, `skills/first-officer/references/fo-merge-core.md:18`. The only shipped invocation is `spacedock merge guard <slug>`. No FO contract file names `--verdict`, and the binary requires it. Every keep-moving stream that ran the guard (2 Haiku CI, 5 local) first got `merge guard requires --verdict {passed|rejected} or --rework`. One local run then gave up and left all three entities unarchived.
3. Contract gap, `skills/first-officer/references/fo-dispatch-core.md:53` (`If not gated: terminal → merge`). The Sonnet CI 10-27 stream and one local run never invoked the guard. They used `status --set … status=done` and `state commit`, saying the stage had "no code deliverable". The backlog claim that the Sonnet stream showed the same refusal is not correct for run 40447244422. That stream has no `merge guard` call at all.

**Item 3 (two contract defects plus a stale pin in binary and contract against a harness change; the suspected link is disproved).**

1. `team-lead` is unreachable on the current Claude Code harness. 40 of 40 `SendMessage(to="team-lead")` calls in the 10-27 CI streams failed, and 23 of 23 decisive `SendMessage(to="main")` calls succeeded, in headless `-p`. The same call from this ensign session also failed (`No agent named 'team-lead' is reachable`), and the harness's SendMessage tool lists `"main"` as the main-conversation address. The pin lives in the binary (`internal/dispatch/build.go` `completionSignalBlock`, with the comment "AC-6") and in the contract (`skills/ensign/references/claude-ensign-runtime.md`, `skills/first-officer/references/claude-fo-dispatch.md:33` and `:43`, `skills/fo-dispatch-recovery/SKILL.md`). `claude-fo-dispatch.md:7` already says workers use `SendMessage(to="main")`, so the contract contradicts itself. The pin dates from `1d691b457` (2026-06-18, Claude Code .178). The harness changed after that.
2. The link is disproved. Haiku 10-07 `selected-team` PASSED with `team-lead` unreachable, and the break-glass template, assertion and runner did not change between `29525bea` and `861a6e13a`. In every failing stream the worker skipped its commit before it called `SendMessage`. Completion still came through the task notification.
3. The real `selected-team` cause is in the contract, `skills/fo-dispatch-recovery/SKILL.md`. The break-glass worker prompt has no commit step. The normal artifact says "after all commits" and the reuse template says "Commit before sending your completion message", but the break-glass prompt says neither. Failing streams: Sonnet CI and local baseline `entity has uncommitted changes: M widget-task.md`, and Sonnet `selected-bare` XFAIL `break-glass-durable-result` from the same cause. The template asks for DONE/SKIPPED/FAILED entries but never asks for the `### Summary` the durable check requires. Haiku CI and local spike run 1 failed on that.

## Proposed approach

**Fix 1.** In `shared_filing_test.go:53`, change the slug terminator from `(?:[ \t'\n]|$)` to `(?:[\s'";&|)]|$)`. Add a table test to `shared_filing_pr399_test.go`. The recorded Haiku command and `…new wire-the-thing;echo` give true. `…new wire-the-thing-two; echo` and `…new wire-the-thingy` give false, which is the boundary guard. Remove the now-obsolete `liveXFail("claude-haiku", "")` binding on `filing`, and its registry row.

**Fix 2a (binary).** In `consumedNonterminalHistory`, a consumed approval whose target is the record stage's immediate, non-terminal successor is history. The current status must be that successor, a later stage reached through ungated stages only, or the terminal stage at the end of that path. A gated stage on the path, a terminal target, or a pending, stale or superseded authority still refuses. This was spiked in a scratch copy:

```diff
-	if app == nil || app.State != "consumed" || app.TargetStage != currentStatus {
+	if app == nil || app.State != "consumed" {
 		return false, nil
 	}
 	stages := parseStagesBlock(filepath.Join(workflowDir, "README.md"))
 	for i, stage := range stages {
-		if stage.Name != record.Stage || i+1 >= len(stages) {
+		if stage.Name != record.Stage || i+1 >= len(stages) || stages[i+1].Name != app.TargetStage {
 			continue
 		}
-		target := stages[i+1]
-		return target.Name == currentStatus && !target.gate && !target.terminal, nil
+		if stages[i+1].terminal {
+			return false, nil
+		}
+		// The spent approval entered an ordinary successor; the entity may have moved
+		// on through ungated stages or already carry the terminal status. A gated
+		// stage on that path is not history this approval covers.
+		for _, later := range stages[i+1:] {
+			if later.gate || (later.terminal && later.Name != currentStatus) {
+				return false, nil
+			}
+			if later.Name == currentStatus {
+				return true, nil
+			}
+		}
+		return false, nil
 	}
```

**Fix 2b (contract).** Apply this before/after text.
- `fo-merge-core.md:18`. Before: ``- → **shipped**: `` `spacedock merge guard <slug>` `` — invoke it directly per phase``. After: ``- → **shipped**: `` `spacedock merge guard <slug> --verdict passed|rejected` `` — `--verdict` is required on every phase; invoke it directly per phase``.
- `fo-dispatch-core.md:53`. Before: `If not gated: terminal → merge; else decide reuse-or-fresh.` After: ``If not gated and the next stage is terminal, run `${SPACEDOCK_BIN:-spacedock} merge guard {slug} --verdict passed|rejected`: it writes the terminal status, verdict and `completed`, then archives. A `status --set` to the terminal stage does not archive and never replaces the guard. Otherwise decide reuse-or-fresh.``

Simplest alternative considered: contract only (2b) without 2a. With 2b, 3 of 3 local runs called the guard directly. 2a is still the fix for a real binary inconsistency: `--set` accepts the terminal write, and then the guard refuses the same entity. Haiku pre-set the status in 2 of 2 CI runs. The gate can drop 2a if it judges 2b enough.

Remove `liveXFail("claude-sonnet", "060xp…")` and `liveXFail("claude-haiku", "")` from `keep-moving-posture`, keep the pi binding, and remove the two registry rows. Reason: 3 of 3 Sonnet local runs XPASS with the fixes.

**Fix 3a (completion target).** Pin `main` instead of `team-lead` in `build.go` `completionSignalBlock` (3 strings and the comment). Make the same change in `claude-ensign-runtime.md` (lines 11, 15, 22, 27: `team-lead` becomes `main`, and the prose says "the first officer"), in `claude-fo-dispatch.md:33` (``pinned to the single name **`main`** — the build helper emits `SendMessage(to="main", …)` … Pin one name.``), in `claude-fo-dispatch.md:43`, and in the recovery team template and its line 25. Regenerate the 15 `internal/dispatch/testdata/golden/*` files. Update `build_merged_mode_test.go` (the assertion flips: `main` is required and `team-lead` is banned), `internal/ensigncycle/cycle_test.go` (the regex) and the banned lists in the codex, pi and json tests. The `team-lead` member name in `reconcile*` team-config fixtures is a TeamCreate roster name and does not change.

**Fix 3b (break-glass template).** In both prompts in `skills/fo-dispatch-recovery/SKILL.md`, make two edits. After the stage-report paragraph, change "…with a `- DONE:` / `- SKIPPED:` / `- FAILED:` entry." to "…entry, then end the section with a `### Summary` subsection." Then append ``\n\n### Commit\n\nBefore you finish, commit your work path-scoped: `git -C {workflow_dir} add -- {entity_file_path} && git -C {workflow_dir} commit -m \"{entity title}: {stage} report\" -- {entity_file_path}`.``. In the team prompt this goes before `### Completion Signal`. Simplest alternative considered: rely on the ensign contract's "MUST commit". That failed in 3 of 4 break-glass streams on 10-27 and in 2 of 2 local baseline subtests. Remove `liveXFail("claude-sonnet", "060xp…")` from `breakGlassSelectedBareGap` and its map test, because Sonnet selected-bare XPASSed in 3 of 3 spike runs. The Haiku blank binding stays. After both removals, owner `060xp69y61yhrww23g3wvwqy` has no binding left. Archiving it is the FO's call.

**Docs (user-visible: merge guard accepts more shapes).**
- `docs/site/reference/command-reference.md:107`, merge guard row. Append: "It also finalizes an entity whose only gate history is a consumed nonterminal approval, after the entity has moved on through ungated stages or a terminal `status --set` has already landed."
- `docs/site/concepts/gates-and-decisions.md:117-120`. Change "The worker can write atomic terminal fields without `--force` after its report is durable." to "The worker can write atomic terminal fields without `--force` after its report is durable, and `merge guard --verdict` finalizes and archives it from any later ungated stage or after that terminal write." Keep the fail-closed sentences.
- `docs/specs/gate-resolution-frontmatter-contract.md:254-257`. Change "…can complete that successor without `--force`." to "…can complete that successor, or any later stage reached only through ungated stages, without `--force`; a gated stage on that path is not covered."
- `docs/runtime-live-ci-registry.md`. Remove the keep-moving Sonnet and Haiku rows and the filing Haiku row.
- The CLI output format and the dispatch artifact layout do not change. Only the completion target string changes.

## Risk evidence

Riskiest mechanisms: whether `main` is the reachable lead name in headless `-p`, and whether the live journeys flip. Both were exercised locally before the gate. Auth: `~/.claude/benchmark-token` inside agent-safehouse, Sonnet `claude-sonnet-5-5`. All commands followed `docs/runtime-live-ci.md`. The runner builds its own binary from the test's checkout and ignores `SPACEDOCK_BIN`. Binary fixes were therefore run from a scratch source copy, and skill fixes through `SPACEDOCK_REPO_ROOT` pointing at a staged plugin copy.

- `TestLiveBreakGlassShimRecovery` baseline on main: `selected-team` FAIL `entity has uncommitted changes: M widget-task.md`, `selected-bare` XFAIL `break-glass-durable-result`, and `team-lead` unreachable. This matches CI.
- Spike run 1 (commit step and `main`, no Summary line): commits landed, zero `team-lead` errors, `selected-bare` XPASS. `selected-team` FAIL: the Summary was missing (the Haiku CI shape).
- Spike runs 2 and 3 (all of 3a for the skills, plus 3b): both subtests PASS in both runs, `selected-bare` XPASS, zero `team-lead` errors.
- `TestLiveCommonKeepMovingPosture` with the main binary and main skills, 3 runs: 3 of 3 hit `ineligible` on `approved-gate` (one FO then reverted to `implementation` and re-ran, which XPASSed).
- Fix 2a alone, 2 runs: one XPASS (the guard archived `approved-gate` after a pre-set `done`). One XFAIL (the FO gave up after the `--verdict` usage error).
- 2a plus 2b `fo-merge-core` line, 2 runs: one run archived all three on the first try, then hung after `end_turn` (no-progress kill). That hang signature also appears in Haiku CI `recorded-gate-lifecycle`, so it is pre-existing. One XFAIL (the FO skipped the guard, the Sonnet CI shape).
- Combined (2a, 2b both lines, 3a binary and skills, 3b), 3 runs: 3 of 3 XPASS, 0 usage errors, 0 `ineligible`, 0 `team-lead` errors.
- Deterministic: the merge-guard reproduction and the spiked fix pass in the scratch copy with the existing `TestConsumedNonterminal*` and `TestTerminal*`. The spiked fix also refuses an entity hand-moved to an unprepared gated stage. `go test ./internal/status ./internal/gates` pass there, except for tests that need `docs/dev` and `install.sh`, which are absent from the scratch copy.

## Out of scope

- Haiku-specific behavior failures: the rejection-flow worker ordering, the named bare `Agent()` call, the literal FAILED worker report, and Haiku skipping `Skill(spacedock:ensign)`.
- Changing which model or cadence the lane runs, or the blank-owner Haiku bindings, except where this task removes a binding its fix makes obsolete (filing, keep-moving).
- The `parseLiveGap` wiring and the two docs that disagree on whether `pi-live` must pass.
- The pre-existing headless hang after `end_turn` (no-progress kill). It was seen once locally and once in Haiku CI.

## Expected surface and tolerance

Estimate net LOC change: +75, across 37 files (about +150 insertions, -75 deletions). Tolerance: net +40 to +120, files 32 to 42. The 15 goldens and most skill edits replace lines one for one.

Semantics this task may change:
- Runtime behavior: `merge guard` and terminal `status --set` now accept consumed nonterminal history on later ungated stages and after a terminal write.
- Dispatch artifact text: the completion target becomes `main`, and the break-glass prompts gain a Commit step and a Summary instruction.
- FO contract terminal routing names `merge guard … --verdict`.

Command grammar, stored formats and gate authority are unchanged. Pending, stale, superseded, terminal-target and gated-path authority still fail closed.

## Acceptance criteria

**AC-1 - A correct filing command with a slug directly followed by `;` (or `&`, `|`, `)`, `"`) grades as filed, and a longer slug sharing the prefix does not.**
Verified by: a table test over the recorded Haiku command. It is red on `661c607ae` and green after the fix, and it fails if the terminator class loses `;`. The negative row `wire-the-thing-two` fails if the boundary is dropped. A local live `TestLiveCommonFiling` run passes on Sonnet.

**AC-2 - A keep-moving drive leaves all three independent entities archived terminal, including the gate-consumed `approved-gate`, measured against the baseline of 3 of 3 refusals.**
Verified by: (a) a `internal/cli/terminal_consume_test.go` test. `--set status=done verdict=PASSED completed` followed by `merge guard --verdict passed` exits 0 and archives. It is red on `661c607ae`. (b) A companion test: an entity hand-moved to an unprepared gated stage is still refused. (c) The existing `TestConsumedNonterminal*` and `TestTerminal*` stay green. (d) Local live `TestLiveCommonKeepMovingPosture` on Sonnet, with the Sonnet and Haiku bindings removed, PASSES in at least 2 of 3 runs, and no stream contains `condition "ineligible"` or `merge guard requires --verdict`.

**AC-3 - Break-glass recovery leaves a committed, complete worker report in both modes, and worker completion messages reach the FO.**
Verified by: (a) local live `TestLiveBreakGlassShimRecovery` on Sonnet. Both subtests PASS in 2 of 2 runs, and `selected-bare` passes with the Sonnet binding removed. (b) No stream from those runs or the AC-2 runs contains `No agent named 'team-lead' is reachable`. (c) `TestBuildMergedModeCompletionSignal` requires `SendMessage(to="main"` and bans `team-lead`. (d) A `skills/integration` check that the recovery team template, `claude-ensign-runtime.md` and the build's completion block name the same target. This fails if one of them drifts back to `team-lead`.

## Test plan

- Primary proof owners: `shared_filing_pr399_test.go` (AC-1), `internal/cli/terminal_consume_test.go` (AC-2, deterministic), `TestLiveCommonKeepMovingPosture` (AC-2, live), `TestLiveBreakGlassShimRecovery` and `TestBuildMergedModeCompletionSignal` (AC-3).
- The new deterministic tests are cheap and go in existing files. Write them red first. The merge-guard reproduction from this ideation seeds the first AC-2 test.
- Live runs, local and Sonnet, through the commands in `docs/runtime-live-ci.md`: filing once, keep-moving 3 times, break-glass 2 times. Cost is about 1 minute and well under $1 per run.
- Before claiming completion, run `go test ./...`, `go test ./... -race` and `gofmt`. The golden regeneration uses the package's existing update path.

### Feedback Cycles

## Stage Report: ideation

- DONE: Diagnose items 2 (keep-moving merge guard "ineligible") and 3 (selected-team plus team-lead unreachable) to a proven root cause from the recorded streams and code, naming whether each defect is in the binary, the FO/ensign contract, or the harness.
  Item 2: a binary defect (`consumedNonterminalHistory` compares the post-`--set` status), reproduced deterministically, plus two FO contract gaps (`merge guard <slug>` without the required `--verdict`, and terminal routing that does not name the guard). Item 3: the break-glass contract template has no commit step and no Summary instruction. The `team-lead` pin in the binary and the contract is stale against the harness, which resolves `main` (40 of 40 failed vs 23 of 23 succeeded in CI, and it fails from this session too). The suspected link is disproved: 10-07 passed with `team-lead` unreachable. See `## Diagnosis`.
- DONE: Propose the smallest fix per item with a falsifiable AC each, an exact proposed diff or file list with net LOC estimate and tolerance, and doc changes if user-visible behavior changes.
  `## Proposed approach` has the merge.go diff, before/after contract wording, binding removals and the doc diffs. AC-1 to AC-3 each name a red-first test. Estimate: net +75 across 37 files, tolerance +40 to +120 and 32 to 42 files.
- DONE: Exercise the riskiest mechanism before the gate, at minimum a local live run (or the cheapest local reproduction) that shows the failure reproduces locally, so implementation can prove each fix passes locally.
  Local live Sonnet runs inside agent-safehouse (benchmark-token auth). Baselines reproduce both failures: break-glass `selected-team` uncommitted, and keep-moving `ineligible` in 3 of 3. Combined spike: break-glass both subtests PASS in 2 of 2 and keep-moving XPASS in 3 of 3, with zero `team-lead` errors. Item 1 was reproduced with a matcher probe (red on the recorded command).

### Summary

All three failures are shared defects, not model faults. Item 1 is a grader regex. Item 2 is a binary predicate that refuses the gate-consumed entity once its terminal status is set, made worse by FO contract text that omits the required `--verdict`. Item 3 is a break-glass worker prompt with no commit or Summary step, alongside a `team-lead` completion pin that no longer resolves on the current harness, which addresses the lead as `main`. The backlog's suspected `team-lead` link to the `selected-team` failure is disproved. All fixes were spiked through local live runs before the gate, outside the main checkout (scratch source copy and staged plugin copy). One residual risk: a pre-existing headless hang after `end_turn` was seen once and is out of scope.

## Stage Report: implementation

- DONE: Implement fixes 1, 2a, 2b, 3a and 3b exactly as approved in Proposed approach, with red-first deterministic tests for AC-1, AC-2 (including the gated-path refusal) and AC-3 (main required, team-lead banned, and the three target names kept in agreement), the obsolete XFAIL bindings and registry rows removed, and the proposed doc edits applied.
  Commits on `live-lane-shared-failure-fixes`: `4426074f7` (fix 1), `0841d0165` (2a, 2b, docs, filing and keep-moving bindings and rows), `f5e10d1bf` (3a, 3b, 15 goldens, Sonnet selected-bare binding), `3393288ea` (test compaction). The merge.go and skill text match the ideation spike byte for byte.
  AC-1: `TestCommandFilesViaNewSlugTerminators` was red on base for 6 rows, including the recorded Haiku `…new wire-the-thing; echo "exit=$?"`. Removing `;` from the terminator class reds that row. Dropping the boundary reds the `wire-the-thing-two` and `wire-the-thingy` rows.
  AC-2: `TestMergeGuardFinalizesConsumedHistoryAfterTerminalSet` (`--set status=done…`, then `merge guard --verdict passed` exits 0 and writes `_archive/task.md`) was red on base with the exact CI `condition "ineligible"` refusal. `TestMergeGuardRefusesConsumedHistoryPastUnpreparedGate` goes red if I delete the `later.gate ||` clause (the guard then archives from the unprepared gated stage). `TestConsumedNonterminal*` and `TestTerminal*` stay green.
  AC-3: `TestBuildMergedModeCompletionSignal` requires `SendMessage(to="main"` and bans `team-lead`. It was red before the build.go change. `TestCompletionTargetAgreesAcrossBuildAndContract` (skills/integration) requires exactly one target, `main`, in the build body, `claude-ensign-runtime.md` and the recovery SKILL.md. Drifting SKILL.md or the ensign runtime back to `team-lead` reds it.
- DONE: Prove each fix passes locally: live Sonnet filing once, keep-moving 3 runs (at least 2 PASS, no "ineligible", no --verdict usage error) and break-glass 2 runs (both subtests PASS), all with zero team-lead-unreachable errors, run from this worktree's checkout.
  `SPACEDOCK_REPO_ROOT`=worktree, `SPACEDOCK_LIVE_MODEL=claude-sonnet-5-5`, benchmark-token auth, commands from `docs/runtime-live-ci.md`. Results: `TestLiveCommonFiling` PASS. `TestLiveCommonKeepMovingPosture` PASS 3 of 3 (bindings removed, so these are real PASS, not XPASS). `TestLiveBreakGlassShimRecovery` PASS 2 of 2, with both `selected-bare` and `selected-team` passing each time. Across all 8 streams: 0 `No agent named`, 0 `condition "ineligible"` or `refusing to finalize`, 0 `merge guard requires --verdict`. The FO sent `"to":"main"` in every keep-moving and break-glass stream. The only bare `ineligible` hits are FO contract text echoed in tool results. Logs and streams are in `/tmp/sd-ll-impl/`.
- FAILED: Measured diff within the approved tolerance (net +40..+120, 32..42 files); gofmt, go test ./..., and go test ./... -race pass.
  `git diff --numstat $(git merge-base main HEAD)..HEAD`: 37 files, +193 -81, net +112. The first measure was +138, over tolerance, so I made the two longest tests smaller (`3393288ea`) and did not change scope. gofmt is clean on every file I touched. One gofmt drift is already on main in `internal/release/runtime_live_evidence_workflow_test.go`; I did not touch it. `go test ./...` passes on all 23 packages: 18 in the full run, and the 5 that hit ENOSPC from other writers (cli, dispatch, ensigncycle, gates, status) passed when rerun one at a time at `-p 1`. The `-race` run is blocked by disk. Free space was 197Mi, and the per-package `-race -p 2` script stopped at its 200Mi floor before the first package with tests. I asked the FO to free space or authorize deletions.

### Summary

All five approved fixes are in the worktree branch, each with a red-first deterministic test that a specific code change makes fail. The local live Sonnet journeys run from this worktree all pass: filing 1/1, keep-moving 3/3 and break-glass 2/2 (both subtests). They show none of the three failure signatures. The diff measured +138 at first, so I made the tests smaller to reach net +112 across 37 files, inside tolerance. The `go test ./... -race` item is still open because the disk is full (under 200Mi free). It is not caused by a code fault. My SendMessage to `team-lead` from this session failed and the send to `main` worked, which confirms fix 3a.
