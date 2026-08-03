---
id: wqnh3jcxzrwej8ry1ryhzbm7
title: Live test registry — single source of truth for live-tagged tests (fixtures, purpose, entry point, enabling lanes)
status: backlog
source: "Captain directed, 2026-08-03, after two audits (live E2E suite composition; unwired-test archaeology + guard design) found 33 live-tagged test functions in internal/ensigncycle with only 12 actually reachable from any CI live lane's -run selector -- 9 real orphaned-from-creation live scenarios plus 12 offline-nature tests wrongly tagged live. Captain explicitly declined the audit's proposed code guard (an AST-diffing contractlint mechanism with a ratcheted ceiling) in favor of a documented, human-reconciled registry: 'no guard, just a document SSOT for test registry. what fixtures, what purpose, entry point, what live lanes enabled this. anything not listed is to be culled or need to be in the list.'"
started:
completed:
verdict:
score: 0.7
worktree:
issue:
---

A single checked-in document enumerating every live-tagged test function in `internal/ensigncycle`, serving as the reconciliation checklist the captain asked for: anything not listed gets culled (deleted) or added to the list (with justification). No code guard -- this is a document, maintained by discipline/incremental review, not enforced by a compiled check.

## Why a document, not a guard

The captain considered and explicitly declined a proposed `internal/contractlint` AST-diffing guard mechanism (already designed and prototyped by a prior audit) in favor of a simpler, human-legible artifact. Do not build the guard as part of this entity.

## Required content per test (from the captain's exact spec)

For every `Test...` function in `internal/ensigncycle/*_live_test.go` (anything gated by `//go:build live`), the registry must record:
1. **Entry point** -- exact `file:line` and function name.
2. **Purpose** -- one or two sentences on what the test actually verifies (not a restatement of its name).
3. **Fixtures** -- what fixture(s)/helper(s) it depends on (e.g. shared scenario definitions, testdata directories, specific assertion helpers).
4. **Enabling lane(s)** -- which CI job(s) in `.github/workflows/runtime-live-e2e.yml` actually select it via `-run` (claude-sonnet / claude-opus / codex / pi), or explicitly "none" if unwired, or "skips at runtime (TODO id X)" if it's registered but self-skips.

## Raw material already gathered (do not re-derive from scratch)

Two prior audits this session already did the archaeology; reuse their findings as the registry's first draft rather than re-investigating from zero:

**33 total live-tagged test functions in `internal/ensigncycle`. 12 actually wired (reachable from some lane's `-run` selector). 21 unwired**, of two distinct kinds:

### A. Real live scenarios that run somewhere (wired, at least one lane)
`feedback-3-cycle-escalation`, `merge-hook-guardrail`, `filing`, `shallow-boot`, `self-evidence-merge-triage`, `smallest-sufficient-mechanism` (partially -- TODO-skipped everywhere, see below), plus `TestLiveMergedTeamModeDispatch`, `TestLivePiFrontDoorSmoke`. Exact file:line and per-lane matrix already tabulated in the live-suite-composition audit's table (reproduce/expand it, don't redo the CI-selector cross-reference from scratch).

### B. TODO-skip tests (wired/selected, but self-skip at runtime with a reason string)
Five `TODO(<24-char-id>)` guards, each already traced to an entity:
- `9adv48yhye5s2vkhwd7ge52d` -- guards `smallest-sufficient-mechanism`, `keep-moving-posture` (both hosts) -- entity `repair-entered-stage-dispatch-and-post-gate-terminalization`, backlog.
- `w5bfnrvpcphw857nzz93340c` -- guards `recorded-gate-lifecycle` (sonnet), `gate-guardrail` (opus) -- entity `reliable-exact-digest-in-gate-review`, backlog.
- `3zzpdw704df1g8pg1x9thzmw` -- guards `gate-guardrail` (sonnet) -- entity `sonnet-gate-guardrail-no-authority`, in flight this session.
- `zbcj98qfwtax61vxdzrf615e` -- guards `rejection-flow` (all hosts, effectively unconditional) -- entity `bind-post-rework-briefing-at-rejection-regate`, **archived 2026-08-03 while carrying verdict: rejected** -- flag this exact anomaly prominently in the registry entry; do not silently normalize it.
- `9w59t6m1qc46hccd54p04z2j` -- guards pi's `recorded-gate-lifecycle` -- entity `pi-delegated-gate-continuation-reliability`, backlog.

### C. Unwired real live scenarios (never in any CI selector, ever, per git history)
`TestLiveReanchorGateRejectsMeansOnlyRegressed` (`ac2_reanchor_live_test.go:21` -- already-known tautology, see `ac2-reanchor-live-scenario-repair` entity; registry should note "do not wire without repair first"), `TestLiveAutoContinueAfterImplementation` (`auto_continue_live_test.go:39`), `TestLivePiAutoContinueAfterImplementation` (`auto_continue_pi_live_test.go:32`), `TestLiveBareReachable` (`dispatch_recovery_live_test.go:57`), `TestLiveBreakGlassShimRecovery` (`dispatch_recovery_live_test.go:83`), `TestLiveHaikuLoopSpike` (`haiku_loop_spike_live_test.go:414` -- intentional one-shot spike, durable evidence already checked into `testdata/haiku_loop_spike_runs/`), `TestLiveHaikuLoopSpikeN` (`haiku_loop_spike_n_live_test.go:28` -- same), `TestLivePiSubagentEnsignSmoke` (`pi_live_runner_test.go:21` -- carries the AC-1 pi boot-contract grader; recommend wiring this one into pi-live's selector as part of registry cleanup, since a debrief incorrectly claims it already runs there), `TestLivePrimitiveRunsAgainstClaudeAdapter` (`livescenario_adapter_live_test.go:78`).

### D. Offline-nature tests wrongly tagged `live` (12, no CI selector needed at all -- just an incorrect build tag)
`TestClaudeTODOModelScope`, `TestClaudeRejectionFlowTODOModelScope`, `TestClaudeSonnetGateGuardrailTODOModelScope`, `TestCleanupKeepMovingRootRetainsOnlyFailures`, `TestCodexLiveRunnerExecArgvEnablesMultiAgentV2`, `TestCodexLiveRunnerUsesSpacedockFrontDoorBeforeHostArgs`, `TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle`, `TestShallowBootFixtureContainsOnlyHeldGate`, `TestPiLiveSmokePromptRequiresExactStageReportHeading`, `TestPiLiveEnvDropsForeignRuntimeMarkers`, `TestPiLiveEnvScrubsAmbientPiSubagentMarkers`, `TestPiIntercomPackageRootDefaultsBesideSubagents`. Registry should list these separately with a recommendation: drop the `live` tag so `go test ./...` runs them (no live spend needed) -- this is a candidate immediate cull-into-compliance, not a thing to leave "unlisted."

## The "cull or list" discipline going forward

The registry's own preamble should state the rule plainly: any live-tagged test function not present in this document is out of compliance -- either add it (with entry point, purpose, fixtures, enabling lane) or delete it. This entity's own implementation should leave the registry in a state where, as of landing, every currently-existing live-tagged test is accounted for (100% coverage of the 33), even if several entries explicitly say "unwired, no lane, kept as a documented spike" or similar.

## Acceptance criteria

**AC-1 (VALUE)** -- Every live-tagged test function that exists in `internal/ensigncycle` at merge time appears in the registry with all four required fields (entry point, purpose, fixtures, enabling lane(s)). Verified by cross-checking the registry's test-name list against `go doc` / `go/build`-based enumeration (reuse the prior audit's proven enumeration approach) -- 1:1, no gaps.

**AC-2** -- The registry explicitly flags the known anomalies rather than silently listing them as ordinary rows: the `zbcj98qfwtax61vxdzrf615e` archived-while-rejected entity, the ac2 re-anchor tautology, and the 12 wrongly-live-tagged tests.

**AC-3** -- The registry states its own maintenance discipline (the cull-or-list rule) in its own preamble, in plain terms a future contributor can follow without re-deriving it.

## Out of scope

No code guard, no CI wiring changes, no test deletions or tag changes -- this entity is documentation only. Follow-up entities (already identified) can act on the registry's findings (untagging the 12 offline tests, wiring `TestLivePiSubagentEnsignSmoke`, resolving the `zbcj...` anomaly).

## Test plan

None in the code-test sense -- this is a documentation deliverable. Verification is the cross-check described in AC-1 (an ensign can do this with a throwaway enumeration script during implementation; it does not need to be a committed test).
