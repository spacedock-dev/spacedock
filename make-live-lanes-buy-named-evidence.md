---
title: Make every live-lane minute buy named evidence
status: implementation
source: "Captain recarve of live-test-truth, 2026-08-03. Absorbs 5p5, b91, v5w, b8, dv, and supersedes 36."
score: 1.0
sprint: live-test-truth
group: named-live-evidence
sprint-readiness: ready
id: 15ec08nz1ypn0dzs8b8xznr7
gates:
    version: 1
    records:
        - id: gate:15ec08nz1ypn0dzs8b8xznr7:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:15ec08nz1ypn0dzs8b8xznr7-backlog-1
              briefing:
                id: briefing:15ec08nz1ypn0dzs8b8xznr7:backlog:attempt-1:revision-1
                digest: sha256:21d99a321493d7d1a6f8699d6f36913a7e9774e952adc96574d482409e4f46d1
                request-digest: sha256:97c4053e9255b7e7a298ec8cc026e0ea331cf318f5e26605d74b2d41b66c8da4
                room-ref: ./make-live-lanes-buy-named-evidence/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:15ec08nz1ypn0dzs8b8xznr7:backlog:1
                briefing: briefing:15ec08nz1ypn0dzs8b8xznr7:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T12:18:54.272717Z"
                decision: approve
                reason: Captain explicitly approved the outcome-shaped recarve and directed immediate redispatch.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:15ec08nz1ypn0dzs8b8xznr7:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:15ec08nz1ypn0dzs8b8xznr7-ideation-1
              briefing:
                id: briefing:15ec08nz1ypn0dzs8b8xznr7:ideation:attempt-1:revision-1
                digest: sha256:e714241f573d403ab40e249008ec231862cbecf6923caa0702215e3d0ab16bea
                request-digest: sha256:750737199c2a6462a3f41bcb547930871600e3b38b7301c4fc159a1aae11fa9f
                room-ref: ./make-live-lanes-buy-named-evidence/review/ideation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-03T14:22:50.455856Z"
                reason: Preflight staff review found cross-member ownership and sequencing defects, and the shared sprint index changed. Withdraw this stale binding before the authorized fold.
            - id: gate-attempt:15ec08nz1ypn0dzs8b8xznr7-ideation-2
              briefing:
                id: briefing:15ec08nz1ypn0dzs8b8xznr7:ideation:attempt-2:revision-1
                digest: sha256:9879538006cec25ad0e845f7dd890d7baf28e41e5f6508197e2b47797198019b
                request-digest: sha256:e8d7a319473e108617c8570b5c37056679c0f0f692d3dba9e15424d49c028d1c
                room-ref: ./make-live-lanes-buy-named-evidence/review/ideation/briefing-2
              withdrawal:
                by: agent:first-officer
                at: "2026-08-03T14:33:56.141728Z"
                reason: Preflight review closure changed the shared sprint index after attempt 2 was frozen. Replace it with a final package that binds the review artifact.
            - id: gate-attempt:15ec08nz1ypn0dzs8b8xznr7-ideation-3
              briefing:
                id: briefing:15ec08nz1ypn0dzs8b8xznr7:ideation:attempt-3:revision-1
                digest: sha256:845b4e5d389bff504eef32b1a2743226a282a0c2e6ac678c92ca2a818a3bc1b3
                request-digest: sha256:8533fcc684e1456ca323b073d760c687adee5ad882287db8881fb2961e93cabe
                room-ref: ./make-live-lanes-buy-named-evidence/review/ideation/briefing-3
              resolution:
                type: Resolution
                id: resolution:spacedock:15ec08nz1ypn0dzs8b8xznr7:ideation:3
                briefing: briefing:15ec08nz1ypn0dzs8b8xznr7:ideation:attempt-3:revision-1
                by: person:captain
                at: "2026-08-03T14:45:26.155175Z"
                decision: approve
                reason: Approved after staff review. Land named-evidence lane hygiene after 3d.
              application:
                target-stage: implementation
                state: consumed
started: 2026-08-03T12:19:37Z
worktree: .worktrees/spacedock-ensign-make-live-lanes-buy-named-evidence
---

## Outcome

A test operator can explain the evidence from each live CI selection, control, artifact, and minute of model use. Dead and duplicate surfaces do not run.

The task owns the lane experience. Selector edits, PTY retirement, offline tags, metrics, and controls are implementation steps.

## Design inputs

- `consolidate-pi-front-door-smoke-evidence` (`5p5`): one Pi smoke must cover the front door, child dispatch, durable output, and boot contract.
- `retire-legacy-pty-live-lane` (`b91`): the current Claude host cannot exercise the legacy PTY regime.
- `wire-claude-dispatch-recovery-live-proofs` (`v5w`): bare and break-glass recovery are required Claude substrate evidence.
- `remove-dead-live-lane-controls-and-metrics` (`b8`): some controls and metrics have no runtime effect or consumer.
- `untag-offline-live-tests` (`dv`): 12 deterministic tests do not run because they have the live build tag.
- `wire-pi-subagent-smoke` (`36`): superseded because a second Pi smoke duplicates spend instead of consolidating the oracle.

## Acceptance criteria

**AC-1 (VALUE) — Each selected live proof has one named, non-duplicate evidence claim.**
Verified by: the lane inventory maps every selector to its unique claim and reports no unreachable selected test.

**AC-2 — The Pi lane spends one smoke run for its required substrate evidence.**
Verified by: one selected smoke grades the front door, child dispatch, durable output, and boot contract.

**AC-3 — The Claude lane runs current recovery evidence and no obsolete PTY regime.**
Verified by: bare and break-glass proofs run, merged-agent evidence remains, and legacy PTY setup and selectors are absent.

**AC-4 — Offline tests, controls, and metrics tell the truth.**
Verified by: the 12 deterministic tests run without the live tag. Each retained control changes a runtime setting, and each retained metric has a producer and consumer.

**AC-5 — The lane cost is visible.**
Verified by: the plan and final evidence record selected test names, durations, model cost where available, and removed duplicate spend.

## Ideation requirements

Start with the current selector-to-claim and producer-to-consumer maps. Choose deletion for an unowned surface unless a cheaper real connection has visible value.

Use `$simple-english` in pragmatic mode for the complete plan.

## Problem

The live workflow mixes release evidence, skipped probes, offline tests, and decorative controls.

The current Pi selector spends one model run but omits its available boot-contract grader. A second direct Pi smoke repeats the same fixture.

The current Claude host supports merged Agent dispatch. The workflow still installs tmux and selects two tests that always skip on this host.

Two required Claude recovery proofs exist but no CI selector runs them. Twelve deterministic tests also remain hidden behind the live build tag.

The workflow publishes three metrics paths. The Pi path has no producer, and the current PR comparison reads only Claude metrics.

## Current selector-to-claim map

This map uses the workflow at commit `2997da9a2`. A selector is accountable only when it has a named claim and an observable result.

| Lane and current selector | Named evidence claim | Disposition |
|---|---|---|
| Claude `TestLiveEnsignCycle` | `full-ensign-cycle` | Keep until the common-suite task replaces its top-level identity. |
| Claude `TestLiveDefaultHeadlessStopsAtGate/default-headless-recorded-gate-stop` | `default-headless-gate-stop` | Keep until common-suite promotion. |
| Claude `TestLiveDefaultHeadlessStopsAtGate/default-headless-withdrawn-gate-recovery` | `withdrawn-gate-recovery` | Keep until common-suite promotion. |
| Claude `TestLiveZeroDiscoverReportsAndStops` | `zero-discovery` | Keep until common-suite promotion. |
| Claude `TestLiveClaudeSharedScenarios/<journey-id>` | One common journey with the same stable ID | Preserve the common-suite owner. Do not create a second runtime-specific proof. |
| Claude `TestLivePtyStandingResidencyInjectsCommOfficer` | Legacy TeamCreate roster evidence | Delete. The installed host cannot produce this evidence. |
| Claude `TestLivePtyEnsignCycleTeamTeardown` | Legacy TeamCreate lifecycle evidence | Delete. The installed host cannot produce this evidence. |
| Claude `TestLiveMergedTeamModeDispatch` | `claude-merged-agent-dispatch` | Keep as the current Claude substrate proof. |
| Codex `TestLiveCodexSharedScenarios/<journey-id>` | One common journey with the same stable ID | Preserve the common-suite owner. |
| Codex `TestCodexResolveManifestAgainstInstalledHost` | Installed current-checkout manifest resolution | Keep as a model-free host-install proof. |
| Pi `TestSharedScenarioRunnerCoverage` | Scenario-table and runner-map parity | Keep until the common-suite task replaces the coverage map. |
| Pi `TestPiSharedScenarioCoverage` | Honest Pi live, codified, or gap classification | Keep until executable Pi common runners replace the map. |
| Pi `TestLivePiFrontDoorSmoke` | `pi-front-door-subagent-dispatch` | Keep and add the boot-contract grader to this run. |
| Pi `TestLivePiRecordedGateLifecycle` | No claim because the first statement always skips | Remove the selector and delete the quarantined wrapper. |

The Claude and Codex shared selectors contain these current journey IDs:

- `gate-guardrail`
- `recorded-gate-lifecycle`
- `rejection-flow`
- `feedback-3-cycle-escalation`
- `merge-hook-guardrail`
- `filing`
- `shallow-boot`
- `self-evidence-merge-triage`
- `smallest-sufficient-mechanism`
- `keep-moving-posture`

`15e` lands after `3d` and consumes its Codex progress-aware liveness behavior. It then updates the current selectors, workflow guards, and live CI guide.

`ys` lands last. It owns the final migration of those shared surfaces to `TestLiveSharedScenarios`.

## Current control map

| Control | Producer | Runtime consumer | Disposition |
|---|---|---|---|
| `claude_version` | `workflow_dispatch` input | Claude installer | Keep. It changes the installed host version. |
| `codex_version` | `workflow_dispatch` input | npm Codex installer | Keep. It changes the installed host version. |
| `effort` | `workflow_dispatch` input | Step-summary text only | Delete. No host receives the value. |
| Claude matrix `model` | Workflow matrix | `SPACEDOCK_LIVE_MODEL` and Claude launch | Keep. It changes every selected Claude model run. |
| `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` | Claude job environment | Legacy PTY path only | Delete with the obsolete PTY regime. The merged path is flag-free. |
| `SPACEDOCK_CODEX_LIVE_REQUIRED` | Codex job environment | Codex authentication gate | Keep. It turns missing approved credentials into failure. |
| `SPACEDOCK_PI_LIVE_REQUIRED` | Pi job environment | Pi authentication gate | Keep. It turns missing approved credentials into failure. |
| `SPACEDOCK_PI_LIVE_CHILD_MODEL` | Pi job environment | Pi root and child launch | Keep. It changes the provider-qualified model. |
| `PI_OFFLINE` | Pi job and isolated Pi environment | Pi package-discovery behavior | Keep. The Pi process receives the value. |
| `SPACEDOCK_LIVE_ARTIFACT_DIR` | Each live job | Host runners and artifact upload | Keep. Every lane writes and uploads this path. |

## Current metric and artifact map

| Surface | Producer | Consumer | Disposition |
|---|---|---|---|
| Claude journey JSON | `emitClaudeScenarioMetrics` | PR delta, release ledger, and operator | Keep. Recovery proofs will use the same producer. |
| Codex journey JSON | `emitCodexScenarioMetrics` | Release ledger only | Connect to the PR delta job and keep the release consumer. |
| Pi journey directory | None | Wildcard upload only | Delete the environment value, upload path, and guard expectation. |
| Pi root and child session JSONL | Pi and pi-subagents | Pi grader and uploaded diagnostic artifact | Keep. Add measured duration and model cost to the grade artifact. |
| `*-detail.jsonl` | One `gotestsum` run per selector group | Test operator after a named failure | Keep. Rename the Claude substrate file after PTY retirement. |
| Candidate binary provenance | Candidate-build steps | Test operator and artifact upload | Keep. It proves the tested revision. |
| Claude and Codex host transcripts | Runtime adapters | Assertions first, operator diagnostics second | Keep. Assertions use durable results where available. |
| `pty-team-mode-detail.jsonl` and tmux session data | Skipped PTY step | No evidence consumer on the current host | Delete. |

The PR delta job will depend on both Claude and Codex jobs. It will download both metrics artifacts before it runs the comparison.

This connection has visible value. A pull request will show Codex duration and token changes before release, not only in a later release ledger.

## Measured baseline

GitHub run `30378538074` at commit `55b3b1331` supplies the current lane timing sample.

| Surface | Sonnet or single-lane time | Opus time | Model cost |
|---|---:|---:|---:|
| Install tmux | 8 seconds | 9 seconds | $0 |
| Both PTY selectors | less than 1 second, both skipped | less than 1 second, both skipped | $0 |
| Merged Claude proof | 127 seconds | 144 seconds | Not available |
| Selected Pi front-door step | 133 seconds | Not applicable | Not in CI summary |

A retained local Pi front-door run lasted 104.808 seconds. Its root and child sessions cost $0.277493 together.

That same retained run provides the Pi consolidation spike. The child read `skills/ensign/SKILL.md` first and the Pi adapter second.

The first five child reads contain no first-officer path. Thus, the selected front-door run already produces the facts that `assertPiEnsignBootContract` grades.

The superseded direct smoke was measured at 172.5 seconds. Selecting it beside the front-door smoke adds one duplicate root and child model run.

The new Claude recovery proofs have no CI timing because CI never selects them. The new substrate step gets a 20-minute total backstop per matrix leg.

Implementation must record each recovery test duration and emitted journey metrics. A total longer than 20 minutes requires a new gate decision.

## Target lane experience

### Claude substrate lane

Use one Claude substrate step for these three proofs:

1. `TestLiveMergedTeamModeDispatch` proves current named background Agent dispatch.
2. `TestLiveBareReachable` proves explicit bare dispatch without recovery loading.
3. `TestLiveBreakGlassShimRecovery` proves a real helper failure routes through break-glass dispatch.

Keep each test name because each name owns one registry claim. Use one `gotestsum` command and one detail artifact.

Delete the tmux install, PTY step, PTY upload path, PTY driver, PTY tests, and PTY-only helpers.

Retain only the generic nested-session environment helper that the merged proof uses. Keep the merged-host version discriminator.

### Pi substrate lane

Keep `TestLivePiFrontDoorSmoke` as the only runtime-specific Pi smoke. Add `assertPiEnsignBootContract` after its durable-state grader.

Delete `TestLivePiSubagentEnsignSmoke`. Its direct Pi launch bypasses the Spacedock front door and repeats the same fixture.

Delete `TestLivePiRecordedGateLifecycle` and remove it from the selector. The wrapper cannot produce evidence because it always skips.

Write these fields to `pi-ensign-boot-grade.json`:

- the stable proof ID
- the root and child models
- the root and child durations
- the root and child costs where the session records them
- the existing spawn and boot-contract grade

The raw Pi sessions remain the source. The grade artifact gives the operator one concise consumer.

### Offline tests

Move these 12 tests to untagged files without changing their assertions:

1. `TestClaudeTODOModelScope`
2. `TestClaudeRejectionFlowTODOModelScope`
3. `TestClaudeSonnetGateGuardrailTODOModelScope`
4. `TestCleanupKeepMovingRootRetainsOnlyFailures`
5. `TestCodexLiveRunnerExecArgvEnablesMultiAgentV2`
6. `TestCodexLiveRunnerUsesSpacedockFrontDoorBeforeHostArgs`
7. `TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle`
8. `TestShallowBootFixtureContainsOnlyHeldGate`
9. `TestPiLiveSmokePromptRequiresExactStageReportHeading`
10. `TestPiLiveEnvDropsForeignRuntimeMarkers`
11. `TestPiLiveEnvScrubsAmbientPiSubagentMarkers`
12. `TestPiIntercomPackageRootDefaultsBesideSubagents`

Move the smallest required helper code with each test. Do not expose a live host launch to the default suite.

## Implementation plan

1. Add workflow guard tests for the post-`3d`, pre-`ys` selectors, controls, artifacts, and consumers.
2. Add an offline command that names all 12 deterministic tests without `-tags live`.
3. Add Pi grade tests with retained session fixtures and wrong boot-contract controls.
4. Make the selected Pi smoke run both durable-state graders.
5. Delete the duplicate and always-skipped Pi wrappers.
6. Add the two Claude recovery selectors beside the merged proof.
7. Delete the legacy PTY workflow and source surfaces.
8. Move the 12 deterministic tests and their required helpers to untagged files.
9. Delete `effort` and the dead Pi metrics path.
10. Connect Codex metrics to the PR delta job.
11. Update the operating guide and exact local commands.
12. Run the focused, full, race, format, and approved live gates.

The workflow guard serves AC-1, AC-3, and AC-4. A prose inventory alone cannot stop a dead selector from returning.

The Pi grade connection serves AC-2 and AC-5. A second live smoke is more expensive and proves no additional substrate boundary.

The existing recovery assertions serve AC-3. A new recovery framework is not necessary because both live tests already emit named metrics.

The PR metrics connection serves AC-4 and AC-5. Deleting Codex metrics removes an active release-ledger producer.

## Expected surface

The baseline is 22 files, about 350 insertions, and about 1,750 deletions. The corrected tolerance is 27 files, 1,000 insertions, and 2,400 deletions.

The delegated Captain conn authorized this estimate correction on 2026-08-05 after the complete compacted candidate measured 27 files, 993 insertions, and 2,381 deletions. The correction covers only the named surfaces below. It adds no outcome or path class; the original estimate undercounted moved deterministic controls, the new falsifiable Pi grade and lane guard, and the actual size of the approved PTY and quarantined-Pi deletions.

| Surface | Expected change |
|---|---:|
| `.github/workflows/runtime-live-e2e.yml` | About 45 insertions and 55 deletions |
| `docs/runtime-live-ci.md` | About 100 insertions and 45 deletions |
| `internal/ensigncycle/pi_live_runner_test.go` and one untagged Pi helper file | About 150 moved or new lines and 55 deletions |
| `internal/ensigncycle/dispatch_recovery_live_test.go` | Comments or registry binding only, at most 8 insertions |
| `internal/ensigncycle/claude_live_runner_test.go` and one untagged unit file | About 75 moved lines and 25 deletions |
| `internal/ensigncycle/codex_live_runner_test.go` and one untagged unit file | About 110 moved lines and 100 deletions |
| `internal/ensigncycle/livescenario_adapter_live_test.go` and one untagged unit file | About 55 moved lines |
| `internal/ensigncycle/shallow_boot_fixture_live_test.go` | Remove the build tag |
| Four PTY source and test files | Delete about 1,500 lines |
| Two Claude host-capability files | Delete legacy-only branches, about 60 lines |
| `internal/ensigncycle/recorded_gate_lifecycle_pi_live_test.go` | Delete about 100 lines |
| `internal/release/journey_workflow_test.go` | About 45 insertions and 20 deletions |
| `internal/release/workflow_exec_guard_test.go` | About 35 insertions and 20 deletions |

The four PTY files are `pty_live_driver_test.go`, `pty_team_mode_live_test.go`, `pty_pane_test.go`, and `pty_session_test.go`.

Move the shared nested-session helper from `pty_session_test.go` before deletion. The merged Claude proof remains its consumer.

No product command, stored workflow format, or authority rule changes. CI selectors, controls, artifacts, test build tags, and live runtime behavior change.

## Acceptance criteria

**AC-1 (VALUE) - Every selected live proof has one named, reachable claim.**

Proved by: the workflow inventory maps every exact selector to a registry claim or a named zero-cost guard.

The inventory reports zero selected top-level tests that always skip. Restoring either PTY selector or the Pi quarantine makes the guard fail.

**AC-2 - The Pi lane spends one runtime-specific smoke run.**

Proved by: `TestLivePiFrontDoorSmoke` uses `spacedock pi` and grades front-door loading, child dispatch, durable output, and the boot contract.

The selector contains no second Pi substrate smoke. Removing any one grader makes a focused negative test fail.

**AC-3 - The Claude lane proves the supported substrate and recovery paths.**

Proved by: merged, bare, and break-glass tests run by exact name on both Claude matrix roles.

The workflow contains no tmux install, PTY selector, PTY artifact, or legacy team environment flag.

**AC-4 - Offline tests, controls, and metrics tell the truth.**

Proved by: all 12 tests pass without the live tag. Each retained control changes a named consumer.

Claude and Codex metrics reach the PR comparison and release ledger. No Pi metrics directory exists without a producer.

**AC-5 - The lane cost is visible and duplicate spend decreases.**

Proved by: the final evidence lists each selected test, duration, model, and model cost when the host records it.

The evidence compares the final run with run `30378538074`. It records the removed 17-second tmux setup and avoided 172.5-second Pi duplicate.

## Test plan

Run the workflow guards first:

```bash
go test ./internal/release -run 'TestRuntimeLiveWorkflow|TestWorkflowsPreserveAndPublishJourneyCosts' -count=1
```

Mutate each retained selector out of a temporary workflow. Each mutation must fail with the missing claim name.

Restore the PTY selector, the Pi quarantine, `effort`, or the Pi metrics path. Each restoration must fail the lane-shape guard.

Run all 12 deterministic tests without the live tag:

```bash
go test ./internal/ensigncycle -run 'TestClaudeTODOModelScope|TestClaudeRejectionFlowTODOModelScope|TestClaudeSonnetGateGuardrailTODOModelScope|TestCleanupKeepMovingRootRetainsOnlyFailures|TestCodexLiveRunnerExecArgvEnablesMultiAgentV2|TestCodexLiveRunnerUsesSpacedockFrontDoorBeforeHostArgs|TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle|TestShallowBootFixtureContainsOnlyHeldGate|TestPiLiveSmokePromptRequiresExactStageReportHeading|TestPiLiveEnvDropsForeignRuntimeMarkers|TestPiLiveEnvScrubsAmbientPiSubagentMarkers|TestPiIntercomPackageRootDefaultsBesideSubagents' -count=1
```

Run the Pi grader against retained fixtures. Remove the ensign read or add a first-officer read to make the negative cases fail.

Run one approved Pi live proof:

```bash
go test -tags live -count=1 -timeout 15m -run '^TestLivePiFrontDoorSmoke$' ./internal/ensigncycle -v
```

The result must contain one Pi root session, one child session, one state commit, and one passing grade artifact.

Run the Claude substrate proofs on Sonnet and Opus through the workflow step:

```bash
go test -tags live -count=1 -timeout 20m -run '^(TestLiveMergedTeamModeDispatch|TestLiveBareReachable|TestLiveBreakGlassShimRecovery)$' ./internal/ensigncycle -v
```

The bare negative controls must reject recovery loading and a missing bare Agent. The break-glass controls must reject wrong ordering and wrong Agent shape.

Run the required repository gates:

```bash
gofmt -w ./cmd ./internal
go test ./...
go test ./... -race
```

Inspect the workflow run, detail JSON, journey JSON, and Pi grade. Record the final producer-to-consumer map and measured cost table.

## Documentation change

Replace the Claude PTY description with this text:

> The Claude substrate step runs merged Agent dispatch, explicit bare dispatch, and dispatch-build break-glass recovery. Each test owns one host-specific claim.

Replace the Pi quarantine description with this text:

> `TestLivePiFrontDoorSmoke` is the only Pi substrate smoke. One run proves the current-checkout front door, child dispatch, durable output, and the ensign boot contract.

Add the selector-to-claim, control, metric, artifact, and measured-cost tables from this plan to `docs/runtime-live-ci.md`.

Remove local commands for PTY and the quarantined Pi gate wrapper. Add the exact Claude substrate command from the test plan.

## Spike record

No new model-backed ideation spike is necessary. The risky Pi join is visible in the retained front-door sessions.

The child read the ensign skill at rank 1 and the Pi adapter at rank 2. It read no first-officer path in its first five reads.

The Claude recovery tests already use the current front door and have offline negative controls. CI selection and measured live execution remain implementation evidence.

## Stage Report: ideation

- DONE: Map every current live selector to one evidence claim and every control or metric to its producer and consumer.
  The body maps 14 selector surfaces, 10 controls, and eight metric or artifact surfaces with final dispositions.
- DONE: Resolve Pi duplication, obsolete Claude PTY work, required recovery proofs, offline tags, and dead surfaces as one lane experience.
  The plan keeps one Pi smoke, three Claude substrate proofs, 12 always-on tests, and no selected always-skip wrapper.
- DONE: Produce one plan with measured lane cost, removed spend, expected files, and falsifiable acceptance evidence.
  The plan records run `30378538074`, a $0.277493 Pi sample, a 22-file estimate, five criteria, and mutation controls.

### Summary

The plan gives every retained live minute a named claim and a visible artifact. It deletes duplicate, skipped, and unowned surfaces instead of preserving decorative machinery.
