---
title: Align Claude break-glass recovery with the selected dispatch mode
status: validation
score: 0.96
source: "PRs #627, #628, #629, and #631 fail TestLiveBreakGlassShimRecovery after PR #626 selected it for required CI. The worker completes through bare blocking dispatch, but the oracle requires a named background worker. History: named recovery template 8e66ead, blanket single-task bare rule ecffced, live selection 4cc0d8."
sprint: durable-decisions
started: 2026-08-07T04:35:54Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-align-claude-break-glass-dispatch-oracle
pr:
issue:
mod-block:
id: 824ecawn5jttbykcgx82nbf4
gates:
    version: 1
    records:
        - id: gate:824ecawn5jttbykcgx82nbf4:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:824ecawn5jttbykcgx82nbf4-backlog-1
              briefing:
                id: briefing:824ecawn5jttbykcgx82nbf4:backlog:attempt-1:revision-1
                digest: sha256:57468f2b244c40ece6ed7b5b8ff8ff84e6ad9407b392af8728d373bc078f419c
                request-digest: sha256:cdfd50d2fd911a26787373b9b3a66a43ec80afc9220a94a9f6ea48c55007b3c7
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:824ecawn5jttbykcgx82nbf4:backlog:1
                briefing: briefing:824ecawn5jttbykcgx82nbf4:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-07T04:35:21.493364Z"
                decision: approve
                reason: The task corrects a proven contract-oracle contradiction, preserves strict worker-completion evidence, and removes repeated unrelated release waivers.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:824ecawn5jttbykcgx82nbf4:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:824ecawn5jttbykcgx82nbf4-ideation-1
              briefing:
                id: briefing:824ecawn5jttbykcgx82nbf4:ideation:attempt-1:revision-1
                digest: sha256:a1fbe668fa324ddfc943f744fbddb639dec7038ba135b65c326f05e7152b0e74
                request-digest: sha256:d1a591247f18d01459a76fc3ed4824c4f0f80234f070f268b4878516f3a7e33e
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:824ecawn5jttbykcgx82nbf4:ideation:1
                briefing: briefing:824ecawn5jttbykcgx82nbf4:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-07T04:42:25.670697Z"
                decision: approve
                reason: The reconciled design preserves dispatch mode, retains strict durable completion proof, fits the approved surface, and isolates landed A7 changes.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:824ecawn5jttbykcgx82nbf4:validation
          stage: validation
          attempts:
            - id: gate-attempt:824ecawn5jttbykcgx82nbf4-validation-1
              briefing:
                id: briefing:824ecawn5jttbykcgx82nbf4:validation:attempt-1:revision-1
                digest: sha256:2b385a727780caeacc6ff9d7afd35e3110b85dd36589385a2157bd008efed8b5
                request-digest: sha256:3a57e6b4c069fb81bcbe977eff1cf9cb91ffe482a57beeaa7a409b322595d1a2
                room-ref: ./review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:824ecawn5jttbykcgx82nbf4:validation:1
                briefing: briefing:824ecawn5jttbykcgx82nbf4:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-07T23:51:09.763957Z"
                decision: revise
                reason: 'Captain approves the design correction: blocking bare mode accepts absent or explicit false run_in_background while still rejecting true, names, teams, wrong identity, extra workers, and incomplete durable reports. Fix the two narrow oracle false positives too.'
            - id: gate-attempt:824ecawn5jttbykcgx82nbf4-validation-2
              briefing:
                id: briefing:824ecawn5jttbykcgx82nbf4:validation:attempt-2:revision-1
                digest: sha256:c306f003a56ee182551aa40a72ed542410848bff04b53f4478edf9731232f2fa
                request-digest: sha256:712e3c0f873aa6b5a1539d92b60397fb38102f8b27127c6a1534cfac01a3984a
                room-ref: ./review/validation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:824ecawn5jttbykcgx82nbf4:validation:2
                briefing: briefing:824ecawn5jttbykcgx82nbf4:validation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-08T07:05:34.321236Z"
                decision: approve
                reason: 'Exact head 43fd2e79d satisfies all four ACs: Sonnet and Opus pass both recovery modes, crossed and malformed supported shapes fail, explicit wrong identity fails, and durable report evidence is exact-section and Git scoped.'
              application:
                target-stage: done
                state: superseded
---

Stable CI must evaluate the supported Claude break-glass behavior instead of rejecting a successful worker because two contracts disagree.

## Problem and value

Single-task Claude dispatch selects bare blocking mode. Break-glass recovery currently hard-codes a named background worker. The required live oracle enforces the older recovery shape.

This contradiction makes unrelated PRs red even when the recovery worker completes and commits its report. Release decisions then require repeated manual waivers.

Current-trunk evidence is direct: `skills/first-officer/references/claude-fo-dispatch.md` selects bare dispatch for a single entity, while `skills/fo-dispatch-recovery/SKILL.md` has only a named `run_in_background=true` manual template. `TestLiveBreakGlassShimRecovery` uses one entity, but `assertBreakGlassObservables` accepts only the named background shape.

## Proposed recovery contract

The dispatch mode is selected before `dispatch build` runs and remains authoritative if assembly fails. Break-glass reports the helper failure and loads `spacedock:fo-dispatch-recovery` as today, then uses exactly one of two explicit manual `Agent` arms:

- **Selected bare mode:** call `Agent` synchronously with `subagent_type`, required `description`, optional effective `model`, and the manual prompt. The call omits `name`, `team_name`, and `run_in_background`; the observable Claude stream may preserve the omission or normalize it to explicit `run_in_background=false`, and both represent the same blocking bare transport. Explicit `true` remains invalid. Omit the prompt's `SendMessage` completion block because the blocking return is the completion signal.
- **Selected team mode:** call the current merged-Claude shape with `subagent_type`, required `description`, capped `name`, `run_in_background=true`, optional effective `model`, and the manual prompt. Omit `team_name`. Retain the prompt's `SendMessage(to="team-lead", ...)` completion block.

Both arms carry the same ensign skill invocation, verbatim stage definition, entity path, checklist, stage-report requirement, and summary slot. Recovery must not probe another transport, retry in the other mode, or convert a bare selection into a named worker merely to satisfy the oracle. The existing stamp/sync and rebase-conflict exclusions remain upstream and unchanged: those failures never enter break-glass.

This contract serves AC-1 and AC-2. The simplest alternative, loosening the oracle to accept the current run, is insufficient because the manual contract would still order the wrong shape. Forcing every recovery into team mode is also insufficient because it breaks single-entity blocking semantics and can end a headless run before durable completion.

No spike is needed. Existing fixtures already prove the two native shapes (`TestBuildBareModeUnchanged` and the merged-mode build tests), the required live failure proves Claude reaches the recovery skill and executes the worker, and the live harness already provides bounded execution plus Git-backed fixture helpers. The unverified work is contract/oracle alignment, which the tests below exercise first.

## Acceptance criteria

**AC-1 (VALUE) — Required Claude live CI gives one correct result for successful break-glass recovery.**

The still-required `TestLiveBreakGlassShimRecovery` passes its bare and team cases when the selected mode runs exactly one worker, returns within the harness bound, and leaves a committed implementation report containing the fixture marker. It fails if no worker runs, more than one worker runs, the marker/report is absent, the entity change is uncommitted, or the selected mode changes.

**AC-2 — Break-glass recovery preserves the dispatch mode selected before helper failure.**

For a single-task bare dispatch, recovery has no `name` or `team_name`, reports absent or explicit-false `run_in_background`, and blocks for completion. For a team dispatch, recovery has a capped `name`, `run_in_background=true`, no `team_name`, and the `team-lead` completion signal.

**AC-3 — The contract, manual template, fixture, and live oracle describe the same behavior.**

An offline table accepts selected-bare/bare-call with absent or explicit-false `run_in_background` and selected-team/team-call, and rejects `true` in bare mode plus both crossed pairs. It also rejects names or teams in bare mode, a missing recovery-skill load, a helper report after the first `Agent`, a malformed prompt, zero workers, or multiple workers.

Claude's live stream may omit the defaulted `subagent_type` from a successful named-background `Agent` input. The required merged-team oracle therefore recognizes the transport from `Agent` plus nonempty `name`, `run_in_background=true`, and absent `team_name`; it proves ensign identity independently from the prompt's `Skill(skill="spacedock:ensign")` and the on-disk member meta's `agentType`. It must reject a bare call, legacy `team_name`, missing ensign prompt/meta, or missing durable result.

**AC-4 — The correction does not weaken worker-completion evidence.**

The proof still requires the ensign skill, verbatim stage definition, one worker execution, fixture marker, exact `## Stage Report: implementation` shape with `DONE` and `Summary`, a path-scoped commit containing that result, a clean entity worktree, and bounded stop at `status: implementation`.

## Test plan

1. Add offline mode-aware transcript fixtures and mutation controls before changing the skill. A mode-preserving bare stream passes with absent or explicit-false `run_in_background`; deleting the only `Agent`, adding a second `Agent`, setting it true, adding bare names/teams, or swapping bare/team fields fails AC-1/AC-2. Cost: small table tests, no model.
2. Add a Git-backed durable-result assertion over the existing recovery fixture. A committed marker plus complete report passes; independently remove the marker, report heading, `DONE`, `Summary`, path-scoped commit, or clean-worktree condition and observe failure. This serves AC-1/AC-4; checking only final prose was considered and rejected because it cannot distinguish an uncommitted worker result.
3. Rewrite the recovery skill as the two explicit mode arms above, then make `assertBreakGlassObservables` accept the selected mode and require exactly one matching call. This serves AC-2/AC-3; a single permissive template was considered and rejected because omitted and present transport fields are the mode contract.
4. Keep `TestLiveBreakGlassShimRecovery` selected in `claude-live` and run two subtests through the real Claude front door and failing helper shim: the existing single-entity prompt selects bare; a mode-only cue selects team. Each runs the stream oracle and durable-result assertion. Do not mark either case optional or TODO. The harness timeout is the bounded-stop proof.
5. Correct the required merged-team oracle's stream predicate and parser to tolerate only an omitted `subagent_type`, then retain its stronger independent grades: named/background/no-team transport, ensign skill in the dispatch artifact, on-disk `agentType=spacedock:ensign`, terminal entity, and path-scoped commit. An offline mutation with omitted `subagent_type` must pass; removing any independent identity or completion evidence must fail. This serves AC-3/AC-4; accepting every `Agent` as an ensign was considered and rejected because a standing or unrelated worker could false-positive.
6. Run focused offline recovery and merged-oracle tests, both required Claude substrate tests on Sonnet and Opus, `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`.

## Documentation diff

In `docs/runtime-live-ci.md`, keep the command and required selector unchanged. Change the substrate-table description from “Merged, bare, and break-glass dispatch” to “Merged and bare dispatch, plus break-glass recovery that preserves the selected bare/team mode and commits the worker report.” No CLI reference changes are needed because command grammar is unchanged.

## Expected surface and semantics

Expected implementation surface: 8 files, about 220 insertions and 55 deletions. Gate tolerance: at most 9 files and 290 insertions; exceeding either returns to ideation.

- `skills/fo-dispatch-recovery/SKILL.md` — two mode-preserving manual arms (~30 insertions, ~10 deletions).
- `internal/ensigncycle/dispatch_recovery_assert_test.go` — mode-aware stream and durable-result oracle (~60 insertions, ~20 deletions).
- `internal/ensigncycle/dispatch_recovery_assert_unit_test.go` — positive/crossed-mode/durability mutation controls (~55 insertions, ~10 deletions).
- `internal/ensigncycle/dispatch_recovery_live_test.go` — bare/team live cases and durable grading (~25 insertions, ~5 deletions).
- `internal/ensigncycle/dispatch_recovery_fixtures_test.go` — team-selection prompt fixture only (~12 insertions).
- `internal/ensigncycle/live_test.go` — recognize a live named-background dispatch when Claude omits its defaulted `subagent_type`, without weakening later completion grades (~12 insertions, ~5 deletions).
- `internal/ensigncycle/merged_team_mode_test.go` — parse and mutation-test that same observed stream shape while retaining the independent member-meta identity check (~18 insertions, ~5 deletions).
- `docs/runtime-live-ci.md` — the concrete wording above (~3 insertions).

Observable semantic change: only Claude runtime behavior after a non-stamp `dispatch build` assembly failure, plus the oracle that grades it. Command grammar, helper JSON, stored entity/frontmatter formats, state authority, retry bounds, normal dispatch, and other hosts do not change.

## A7 / PR #632 reconciliation boundary

This design is reconciled against landed A7 merge `e24e8234a0eff20fe1d5efa3bb600eb51811a532`. A7 does not touch any of the eight expected files. It does touch the same `internal/ensigncycle` package (`recorded_gate_lifecycle_test.go`, `shared_scenarios_test.go`) and shared First Officer gate prose, but not Claude dispatch selection, recovery, or the merged dispatch oracle.

Implementation must start from or rebase onto `e24e8234`, then run focused recovery/merged-oracle tests plus full/race suites on that reconciled tree. Task 824 must not edit or revert A7-owned gate lifecycle, command, spec, site, or roadmap changes.

The two reported oracle assumptions both belong in 824, narrowly. Opus proves the recovery oracle incorrectly substitutes team shape for a successful selected-bare call. Sonnet proves the team half cannot use serialized `subagent_type` as its sole dispatch discriminator. Correcting only the first would make the new team recovery proof depend on an assumption the required merged-team lane has already falsified. This is not weaker completion evidence: transport shape, ensign prompt, on-disk member identity, durable report/terminal entity, path-scoped commit, and bounded stop remain jointly required.

## Stage Report: ideation

- DONE: Define one recovery contract that preserves the selected bare or team dispatch mode.
  The proposed contract has explicit bare blocking and merged-team background arms and forbids recovery-time mode changes.
- DONE: Define falsifiable positive and negative proof for recovery and durable worker completion.
  The test plan names crossed-mode, missing/extra-worker, malformed-prompt, missing-report-part, uncommitted-result, and timeout failures.
- DONE: Set the smallest expected file and LOC surface, including the later A7 reconciliation boundary.
  The reconciled baseline is 8 files/~220 insertions with a 9-file/290-insertion gate; landed A7 remains outside scope at merge `e24e8234`.

### Summary

Ideation aligns manual recovery with the mode selected before helper failure and keeps both required Claude substrate checks load-bearing. It adds durable Git-backed completion evidence, admits Claude's observed omission of defaulted `subagent_type` only where independent identity proof remains, and reconciles cleanly against landed A7.

## Stage Report: implementation

- DONE: Implement explicit bare and team recovery arms without changing the selected dispatch mode.
  Commit `bbe3d7a05` gives bare recovery a blocking field-omitting arm and team recovery the named-background, no-`team_name`, `team-lead` arm; crossed-mode and multiple-worker mutations fail.
- DONE: Correct only the two proven Claude oracle assumptions while retaining independent identity and durable completion checks.
  Mode-aware stream tests require one exact transport, while Git-backed mutations reject missing marker/report parts, non-path-scoped or uncommitted results, and dirty entity state; merged transport tolerates omitted `subagent_type` only alongside ensign artifact and member-meta proof.
- DONE: Stay within the approved 9-file and 290-insertion limit and keep every required Claude check enabled.
  The A7-reconciled diff is 9 files and 285 insertions; selectors remain required, with both break-glass modes in `TestLiveBreakGlassShimRecovery` and no A7-owned or unrelated roadmap files changed.

### Evidence

- Focused offline recovery and merged-oracle tests passed; their mutations reject crossed modes, zero/multiple workers, malformed prompts, missing durable report parts, uncommitted or dirty entity results, and missing merged identity/completion evidence.
- `go test ./...` failed only in `internal/gates/TestV1PilotManifestReadsAndValidates`: `codex-launch-multi-agent-v2.md` and `gate-agent-ergonomics.md` were absent from the shared state checkout. All other packages passed.
- `go test ./... -race` failed on the same two absent shared-state manifests in `internal/gates/TestV1PilotManifestReadsAndValidates`. All other packages, including `internal/ensigncycle`, passed under race.
- Sonnet `TestLiveBreakGlassShimRecovery` passed both `selected-bare` and `selected-team`, including strict stream shape and durable committed-result grading.
- Sonnet `TestLiveMergedTeamModeDispatch` failed before FO work with HTTP 401: `OAuth access token has been revoked.`
- Opus `TestLiveBreakGlassShimRecovery/selected-team` passed. Its first `selected-bare` run committed the complete durable result but failed the strict stream oracle because the call serialized `run_in_background:false` instead of omitting the key.
- After the manual contract was clarified to prohibit passing false, the isolated Opus `selected-bare` rerun reached a committed blocking worker result, then failed the harness's one-minute no-progress quiet bound before the stream/durable oracle ran.
- Opus `TestLiveMergedTeamModeDispatch` failed before FO work with HTTP 401: `OAuth access token has been revoked.`

### Summary

Implemented mode-preserving Claude break-glass recovery, strict selected-mode and one-worker stream grading, durable Git-backed worker-result grading, and the narrowly tolerant merged-team transport parser without weakening independent ensign identity or completion checks. The candidate is committed directly atop A7 at `bbe3d7a05`; remaining red observations are recorded above without disabling or relaxing any required check.

## Stage Report: validation

- FAILED: Verify all four acceptance criteria at exact head, including crossed-mode and durable-result negative controls.
  AC-1/AC-2 fail live on Sonnet and Opus because successful blocking recovery serializes `run_in_background:false`; AC-3/AC-4 fail the adversarial identity and report-structure controls below.
- FAILED: Audit the 9-file and 285-insertion surface for A7 isolation and no weakened identity or completion evidence.
  Surface and A7 isolation are exact and clean, but `mergedEnsignDispatches` accepts an explicit non-ensign type and `assertCompleteRecoveryReport` accepts report tokens outside the report section.
- DONE: Run focused, full, race, and required live evidence; classify every failure from this candidate.
  Focused tests pass; full/race fail only on two absent shared-state manifests; both-model live results and merged-lane auth failures are classified below.

### Findings

- Material evidence defect / Needs decision: the bare live oracle cannot distinguish an omitted field from Claude's normalized `run_in_background:false`.
  Released user/normal workflow: required Sonnet and Opus `TestLiveBreakGlassShimRecovery/selected-bare`; harm: false-red CI after exactly one blocking worker leaves a complete path-scoped committed result; authority: value-ac[AC-1] requires one correct live result; trigger: both 2026-08-07 runs emitted no name/team/completion signal but `run-present=true, run_in_background=false`. Proposed ownership: this task; disposition: design reset because accepting normalized false conflicts AC-2's specified absence boundary.
- Material evidence defect / narrow task-owned fix: explicit wrong worker identity passes the merged parser.
  Released user/normal workflow: required merged Claude oracle; harm: an `Agent(subagent_type="general-purpose", name=..., run_in_background=true)` can be graded as ensign; authority: value-ac[AC-4] requires independent ensign identity without weakening; trigger: overlay adversarial test returned one dispatch for that exact stream at `merged_team_mode_test.go:57-60`. Proposed disposition: fix after FO authorization.
- Material evidence defect / narrow task-owned fix: durable report parts are not atomically scoped.
  Released user/normal workflow: required break-glass durable grading; harm: unrelated prior `- DONE:` and `### Summary` text satisfies an empty implementation report; authority: value-ac[AC-4] requires the exact implementation report shape; trigger: overlay adversarial test passed scattered tokens through `assertCompleteRecoveryReport` at `dispatch_recovery_assert_test.go:289-295`. Proposed disposition: fix after FO authorization.

### Evidence

- Focused offline tests passed crossed bare/team pairs, zero/multiple workers, malformed prompt, durability mutations, and merged identity fixtures; the two added overlay counterexamples failed as described and were removed without changing HEAD.
- `go test ./...` and `go test ./... -race` passed `internal/ensigncycle` and failed only `internal/gates/TestV1PilotManifestReadsAndValidates` because `codex-launch-multi-agent-v2.md` and `gate-agent-ergonomics.md` are absent from the shared state checkout; external-state failures, not candidate failures.
- Sonnet and Opus each passed `selected-team`; each `selected-bare` produced the marker, complete report, path-scoped commit, clean entity, and blocking completion, then failed only the strict false-vs-absent stream check.
- Sonnet and Opus `TestLiveMergedTeamModeDispatch` each failed before FO work with HTTP 401 `OAuth access token has been revoked`; harness/auth evidence failure, not candidate behavior evidence.
- `gofmt -d ./cmd ./internal`, `git diff --check`, and worktree status were clean; exact diff is 9 files/285 insertions, directly atop A7 `e24e8234`, with no A7-owned file touched.

### Recommendation

REJECTED. The supported live bare path false-reds on both required models, and the new merged identity and durable-report oracles admit concrete false positives; no candidate mutation was authorized or made.

### Summary

Validation confirms the mode-preserving worker behavior and team recovery, but the observable omission requirement is incompatible with current Claude stream normalization. Two additional adversarial controls show weakened identity and completion evidence, so all four ACs lack valid end-to-end proof at this head.

### Feedback Cycles

- Cycle 1: REJECTED — detached validation / normalized false plus two oracle false positives; surface 9 files/+285 vs estimate 8 files/about +220/-55; AC narrowed: bare blocking accepts absent or explicit false

## Stage Report: implementation (cycle 2)

- DONE: Implement the Captain-approved bare-mode rule: accept absent or explicit false run_in_background while rejecting true, names, teams, crossed modes, zero workers, and multiple workers.
  Commit `43fd2e79d` adds absent/false positives and true/name/team negative controls; the existing crossed-mode and worker-cardinality controls remain green.
- DONE: Reject explicit non-ensign merged dispatches and require DONE and Summary tokens inside the exact implementation Stage Report section.
  Focused mutations now fail on `subagent_type="general-purpose"` and on DONE/Summary tokens placed before an empty exact implementation report; both pass only after the strict predicates in `43fd2e79d`.
- DONE: Update the task contract and focused mutations, stay within the approved surface, and rerun focused, full, race, formatting, diff, and applicable live evidence at the exact head.
  The same 9 implementation files are touched at 350 insertions/62 deletions; focused tests and formatting/diff checks pass, full/race fail only on two absent shared-state manifests, and Sonnet/Opus live runs stop before FO work on OAuth 401.

### Summary

The correction accepts Claude's explicit-false normalization only for blocking bare recovery and closes the two reviewer-proven identity and report-scoping false positives. The candidate is committed at `43fd2e79d`; all offline task-owned evidence passes, while live model evidence remains externally blocked by expired or revoked Claude OAuth credentials.

## Stage Report: validation (cycle 2)

- DONE: Verify absent-or-false blocking bare semantics and all crossed-mode, true, name, team, zero-worker, and multiple-worker negative controls at exact head 43fd2e79d.
  Focused tests pass both bare encodings, team mode, and every named mutation; fresh Sonnet and Opus live runs pass selected-bare and selected-team end to end at `43fd2e79d`.
- DONE: Verify explicit non-ensign merged dispatch and scattered report-token counterexamples fail while durable identity, exact Stage Report, path-scoped commit, clean worktree, and bounded-stop evidence remain strict.
  Focused plus temporary overlay controls reject explicit `general-purpose` identity and DONE/Summary before or after the exact section; Git-backed marker, report, commit, dirty-path, and completion-bound mutations remain falsifiable.
- DONE: Audit the fixed 9-file surface against A7 and F6C boundaries, then run focused, full, race, formatting, diff, and applicable live evidence with external failures classified from this run.
  The exact A7-based surface is nine files at +350/-62 with no A7-owned file or F6C command/spec/help change; F6C overlap is confined to 824-owned recovery oracles/skill and its runtime-live description.

### Acceptance evidence

- AC-1: Sonnet passed both break-glass modes in 278.90s and Opus passed both in 499.83s; each live case requires exactly one worker, bounded return, marker, exact report, path-scoped commit, and clean entity path.
- AC-2: Offline absent/false positives and true/name/team/crossed-mode negatives pass; both models independently passed selected-bare and selected-team through the real failing-helper front door.
- AC-3: Offline tables reject crossed pairs, malformed prompts, zero/multiple workers, explicit non-ensign types, and scattered report tokens; contract and fixture remain aligned with the observed live recovery shapes.
- AC-4: Durable-result mutations independently reject missing marker, heading, DONE, Summary, commit, or clean state; merged artifact and member-meta mutations retain independent ensign identity.

### External evidence limitations

- `go test ./...` and `go test ./... -race` pass `internal/ensigncycle` and fail only `internal/gates/TestV1PilotManifestReadsAndValidates` because two shared-state manifests are absent; this is external checkout state, not candidate behavior.
- Fresh Sonnet and Opus `TestLiveMergedTeamModeDispatch` each stop before First Officer work with HTTP 401 `OAuth access token has been revoked`; offline merged transport, artifact, and member-meta controls pass, but this run adds no new merged live observation.

### Deferred risks

- The stream oracle treats present non-boolean `run_in_background` values (`null`, `"false"`, or `0`) as false because decode errors are ignored. Claude tool input canonically supplies a JSON boolean, so no supported trigger or value AC fails; promote to Material if a supported host emits an untyped value or fixture ingestion admits noncanonical tool input.

### Recommendation

PASSED. No material finding remains; the malformed-boolean parser weakness is deferred with a concrete promotion condition, and the two external failures do not contradict candidate behavior.

### Summary

Validation at exact head `43fd2e79d` reproduces all four acceptance criteria, including fresh two-model live proof of both recovery modes and adversarial proof of the two corrected false positives. The candidate stays within the approved nine-file boundary, formatting and diff checks are clean, and only classified external-state limitations plus one unsupported-input deferred risk remain.

## Stage Report: implementation (cycle 3)

- DONE: Reconcile the candidate with current main and the YS live-test inventory.
  Merge commit `6c50fc82e` incorporates canonical `origin/main`/PR #633 and keeps the `TestLiveCommon` registry, selectors, fixture/proof annotations, and substrate artifact bundle intact.
- DONE: Resolve both live-test conflicts without weakening the selected dispatch-mode oracle.
  The guide retains YS common-inventory rows plus 824's strict substrate claim; `TestLiveBreakGlassShimRecovery` retains both bare/team subtests, one-worker mode grading, and durable path-scoped committed-report grading.
- DONE: Update required inventory coverage and prove focused, full, race, Sonnet, and Opus behavior.
  Registry reconciliation and focused mutations pass; full/race fail only on the two absent shared-state pilot manifests; Sonnet passes bare/team in 389s and Opus passes bare/team in 539s at `6c50fc82e`.

### Summary

The candidate is reconciled onto the canonical YS inventory without resurrecting removed common-runner surfaces or absorbing the separate journey-delta-comment defect. The runtime registry now requires selected bare/team preservation plus a committed complete worker report, and both required Claude models prove that behavior at exact head `6c50fc82e`.

## Stage Report: validation (cycle 3)

- DONE: Verify exact head 6c50fc82e preserves the selected bare and team dispatch-mode contract and all acceptance criteria.
  Focused mode/cardinality/identity/report mutations pass at `6c50fc82e`; fresh Sonnet and Opus runs each pass selected-bare and selected-team with strict stream and durable Git grading.
- DONE: Verify YS registry, annotations, selectors, artifacts, and strict durable-report grading remain consistent.
  Registry reconciliation proves 16 registry/source bindings, fixture/proof annotations, three exact common selectors, and retired-symbol absence; workflow mutation guards prove named evidence and substrate artifact selection.
- DONE: Re-run focused, full, race, Sonnet, and Opus evidence; classify each current failure from its own output.
  Focused and both-model live evidence pass; full/race fail only on the two currently absent shared-state pilot manifests, reproduced independently in both commands below.

### Acceptance evidence

- AC-1: Sonnet passed both recovery modes in 417.14s and Opus passed both in 487.66s; every subtest requires exactly one worker, bounded return, the fixture marker, exact report, path-scoped commit, and a clean entity path.
- AC-2: Absent and explicit-false bare inputs pass, while true/name/team/crossed-mode inputs fail; the named background team shape and completion signal pass offline and live on both models.
- AC-3: Contract, fixture, oracle, live-proof annotations, and registry outcome agree; missing skill/report ordering/stage definition, zero/multiple workers, explicit non-ensign identity, and scattered report tokens all fail focused controls.
- AC-4: Git-backed mutations reject every missing report component, uncommitted result, dirty entity, and non-path-scoped commit; YS reconciliation retains the substrate detail artifact and exact selector.

### Reconciliation boundaries

- Relative to current-main parent `d3e70e958`, the candidate is exactly the approved nine 824 files; the merge resolution adds only recovery semantics to the YS registry, guide, and live test.
- `TestRuntimeLiveRegistryReconciliation` rejects the retired common-runner symbols and passes. The candidate does not modify `.github/workflows/runtime-live-e2e.yml` relative to current main, so it neither owns nor absorbs `journey-delta-comment`; its separate workflow guards pass unchanged.

### Current external failures

- `go test ./...` fails only `internal/gates/TestV1PilotManifestReadsAndValidates` because `codex-launch-multi-agent-v2.md` and `gate-agent-ergonomics.md` are absent from the shared state checkout; all other packages pass.
- `go test ./... -race` independently fails only the same two missing-manifest cases; `internal/ensigncycle`, `contractlint`, `release`, and `status` pass under race.

### Deferred risk

- The previously recorded noncanonical JSON-boolean risk remains unchanged: a present `null`, string, or numeric `run_in_background` can decode as false. Supported Claude tool input emits a boolean; promote if a supported stream emits an untyped value or ingestion begins accepting noncanonical tool input.

### Recommendation

PASSED. No material finding remains, the YS reconciliation preserves all four ACs, and the current red results are external shared-state availability failures rather than candidate defects.

### Summary

Successor validation at exact head `6c50fc82e` replaces the stale pre-YS proof with fresh focused, full, race, Sonnet, and Opus evidence. The selected-mode recovery contract and durable completion proof remain strict inside the canonical YS inventory, without restoring retired runners or taking ownership of the journey-delta job.

## Stage Report: implementation (cycle 4)

- DONE: Rebase the post-YS 824 candidate onto actual landed 0Y tip 9021cbf37 without discarding either side.
  FO-directed no-force reconciliation uses merge `d42b6f40d` with parents `6c50fc82e` and exact landed 0Y `9021cbf37`; task-owned Sonnet proof restoration is committed at pushed exact head `4ecc0440d`.
- DONE: Preserve selected bare/team recovery grading, YS registry integration, and pilot-manifest decoupling.
  Focused recovery mutations reject crossed modes, wrong cardinality/identity, and incomplete or uncommitted reports; registry/cadence guards pass, the retired runners remain absent, and 0Y's deleted pilot-manifest test is not restored.
- DONE: Run focused, inventory, full, race, Sonnet, and Opus evidence at the new exact head.
  Focused, inventory, release-cadence, `go test ./...`, and `go test ./... -race` pass; Sonnet and Opus each run both recovery cases but stop before FO work because local Claude OAuth is expired and cannot refresh.

### Summary

Reconciled 824 with landed 0Y through the smallest fast-forwardable, non-force topology and pushed exact head `4ecc0440d` on only the named 824 branch. The selected-mode and durable-report oracle remains strict inside the YS inventory, while 0Y's Sonnet-on-PR, Opus-pre-release, and Codex cadence plus its pilot-manifest deletion remain intact; fresh live model execution is externally blocked by expired OAuth.

## Stage Report: validation (cycle 4)

- DONE: Verify exact head 4ecc0440d preserves selected bare/team recovery, worker identity, cardinality, and durable report grading.
  Focused controls accept bare absent/false and the named-background team shape, and reject crossed modes, true/name/team bare fields, zero/multiple workers, explicit non-ensign identity, incomplete/scattered reports, uncommitted results, dirty state, and non-path-scoped commits.
- DONE: Verify YS live inventory, 0Y pilot-manifest decoupling, and Sonnet/Codex/Opus cadence remain intact.
  Registry reconciliation proves all 16 common journeys and exact selectors; release guards prove one Sonnet-PR/Opus-pre-release Claude row plus the Codex approval lane, while full/race pass with the retired runners and pilot-manifest gate absent.
- DONE: Run focused, inventory, full, race, formatting, and available Sonnet/Opus evidence; classify unavailable live evidence precisely.
  Focused, registry, cadence, `go test ./...`, `go test ./... -race`, gofmt listing, `git diff --check`, and cleanliness checks pass; both fresh model probes fail authentication before any First Officer work.

### Acceptance evidence

- AC-1: The exact-head offline and Git-backed graders require one selected-mode worker, bounded completion, marker, exact report, path-scoped commit, and clean entity; prior Sonnet/Opus proof at `6c50fc82e` exercised both modes, and `4ecc0440d` restores the Sonnet selector after the no-force 0Y merge without changing those graders.
- AC-2: The adjacent mode matrix accepts omitted/false bare background state and selected team, while rejecting true, bare names/teams, crossed pairs, and wrong cardinality.
- AC-3: Recovery prose, fixtures, stream oracle, registry annotations, and workflow selectors agree; mutations fail on missing skill load, late helper report, malformed prompt, wrong identity, and retired-runner restoration.
- AC-4: Complete-result mutations independently reject every missing marker/report component, scattered tokens, uncommitted or dirty results, and non-path-scoped commits; merged transport still requires prompt and member-meta ensign identity.

### External evidence limitation

- Fresh `claude-sonnet-5` and `claude-opus-4-8` runs each attempted both recovery cases with Claude Code 2.1.220, but every case returned `OAuth session expired and could not be refreshed` before First Officer work. These are credential failures and provide no new live candidate evidence.

### Deferred risk

- The existing unsupported-input risk remains: non-boolean `run_in_background` JSON can be treated as false. Supported Claude tool input is boolean; promote if a supported stream emits an untyped value or ingestion begins accepting noncanonical input.

### Recommendation

PASSED. Exact-head deterministic evidence shows no selected-mode, identity, cardinality, durability, YS inventory, 0Y decoupling, or cadence regression; fresh live evidence is unavailable solely at the pre-work credential boundary.

### Summary

Validation at pushed head `4ecc0440d` reproduces the task-owned acceptance evidence and the YS/0Y reconciliation boundaries through focused adversarial, inventory, full, race, and formatting checks. The candidate remains unchanged and clean; local Sonnet and Opus evidence is precisely classified as unavailable before First Officer execution.
