---
title: "Verify complete break-glass worker reports before accepting completion"
status: "backlog"
source: "Captain filing request; pre3 and 0.27.3 live release failure triage"
started: ""
completed: ""
verdict: ""
score: "0.9"
worktree: ""
issue: ""
pr: ""
mod-block: ""
id: 4bnd9tb6gqj4592dt9w8nwvn
---

Verify complete break-glass worker reports before accepting completion.

## Problem

The 0.27.3 Claude substrate run accepted a worker result without the required Summary. The worker ran successfully; this is completion verification, not missing team infrastructure.

Evidence: https://github.com/spacedock-dev/spacedock/actions/runs/34923156550, job 104235722775, artifact 10379746970, scenario break-glass-shim-selected-team. widget-task.md contains DISPATCH-RECOVERY-WORKER-RAN, the exact Stage Report: implementation heading, and two DONE rows, but no ### Summary. The FO read it and declared completion. assertCompleteRecoveryReport correctly rejected it before committed-durability checks. pre3 substrate passed.

Related prior work: align-claude-break-glass-dispatch-oracle (completed); align-claude-break-glass-agent-proof covers Agent topology, not this incomplete-report acceptance.

## Proposed approach

Apply the existing report-completeness requirement at worker return. Return a missing section to the same addressable worker and verify its repaired committed report before accepting completion. Include the complete canonical report format in the fallback assignment. Rename selected-team / dispatchModeTeam and the recovery prompt's team-mode wording to named background worker; preserve Agent name + run_in_background and the prohibition on team_name. Do not restore TeamCreate or weaken the existing report oracle.

## Risk evidence

The captured release artifacts above establish the failure trigger. Replay them before implementation; a live runtime claim requires the targeted live AC below, not a prose-presence check.

## Out of scope

Unrelated release failures and broader test-harness redesign.

## Expected surface and tolerance

Existing recovery contract, fallback assignment, and adjacent recovery tests; estimate net +30 to +80 LOC across 3–5 files. No new lifecycle metadata, hooks, CI lane, or runtime protocol.

## Acceptance criteria

**AC-1 — An incomplete fallback report is not accepted.**
Verified by: a behavioral recovery exercise using the captured missing-Summary report observes a correction to the same worker, then acceptance only after a complete committed report. Removing the return-time check must fail this exercise.

**AC-2 — Dispatch remains one named background worker without team infrastructure.**
Verified by: existing mode/cardinality tests accept the supported Agent shape and reject team_name or a duplicate worker; retain bare-mode coverage.

**AC-3 — The affected Claude recovery journey passes live.**
Verified by: a targeted local TestLiveBreakGlassShimRecovery run exercises both dispatch variants and passes the unchanged completeness and durability oracle.

## Test plan

Add focused failing behavioral cases first. Reuse the existing test owner and captured artifacts, then run the targeted live AC. Before completion run go test ./..., go test ./... -race, and gofmt -w ./cmd ./internal. Report unavailable live execution separately from passing evidence. No prose-grep proof.

### Feedback Cycles

