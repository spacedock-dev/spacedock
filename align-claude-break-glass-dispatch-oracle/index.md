---
title: Align Claude break-glass recovery with the selected dispatch mode
status: implementation
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
---

Stable CI must evaluate the supported Claude break-glass behavior instead of rejecting a successful worker because two contracts disagree.

## Problem and value

Single-task Claude dispatch selects bare blocking mode. Break-glass recovery currently hard-codes a named background worker. The required live oracle enforces the older recovery shape.

This contradiction makes unrelated PRs red even when the recovery worker completes and commits its report. Release decisions then require repeated manual waivers.

Current-trunk evidence is direct: `skills/first-officer/references/claude-fo-dispatch.md` selects bare dispatch for a single entity, while `skills/fo-dispatch-recovery/SKILL.md` has only a named `run_in_background=true` manual template. `TestLiveBreakGlassShimRecovery` uses one entity, but `assertBreakGlassObservables` accepts only the named background shape.

## Proposed recovery contract

The dispatch mode is selected before `dispatch build` runs and remains authoritative if assembly fails. Break-glass reports the helper failure and loads `spacedock:fo-dispatch-recovery` as today, then uses exactly one of two explicit manual `Agent` arms:

- **Selected bare mode:** call `Agent` synchronously with `subagent_type`, required `description`, optional effective `model`, and the manual prompt. Omit `name`, `team_name`, and `run_in_background`. Omit the prompt's `SendMessage` completion block because the blocking return is the completion signal.
- **Selected team mode:** call the current merged-Claude shape with `subagent_type`, required `description`, capped `name`, `run_in_background=true`, optional effective `model`, and the manual prompt. Omit `team_name`. Retain the prompt's `SendMessage(to="team-lead", ...)` completion block.

Both arms carry the same ensign skill invocation, verbatim stage definition, entity path, checklist, stage-report requirement, and summary slot. Recovery must not probe another transport, retry in the other mode, or convert a bare selection into a named worker merely to satisfy the oracle. The existing stamp/sync and rebase-conflict exclusions remain upstream and unchanged: those failures never enter break-glass.

This contract serves AC-1 and AC-2. The simplest alternative, loosening the oracle to accept the current run, is insufficient because the manual contract would still order the wrong shape. Forcing every recovery into team mode is also insufficient because it breaks single-entity blocking semantics and can end a headless run before durable completion.

No spike is needed. Existing fixtures already prove the two native shapes (`TestBuildBareModeUnchanged` and the merged-mode build tests), the required live failure proves Claude reaches the recovery skill and executes the worker, and the live harness already provides bounded execution plus Git-backed fixture helpers. The unverified work is contract/oracle alignment, which the tests below exercise first.

## Acceptance criteria

**AC-1 (VALUE) — Required Claude live CI gives one correct result for successful break-glass recovery.**

The still-required `TestLiveBreakGlassShimRecovery` passes its bare and team cases when the selected mode runs exactly one worker, returns within the harness bound, and leaves a committed implementation report containing the fixture marker. It fails if no worker runs, more than one worker runs, the marker/report is absent, the entity change is uncommitted, or the selected mode changes.

**AC-2 — Break-glass recovery preserves the dispatch mode selected before helper failure.**

For a single-task bare dispatch, recovery has no `name`, `team_name`, or `run_in_background` and blocks for completion. For a team dispatch, recovery has a capped `name`, `run_in_background=true`, no `team_name`, and the `team-lead` completion signal.

**AC-3 — The contract, manual template, fixture, and live oracle describe the same behavior.**

An offline table accepts selected-bare/bare-call and selected-team/team-call, and rejects both crossed pairs. It also rejects a missing recovery-skill load, a helper report after the first `Agent`, a malformed prompt, zero workers, or multiple workers.

Claude's live stream may omit the defaulted `subagent_type` from a successful named-background `Agent` input. The required merged-team oracle therefore recognizes the transport from `Agent` plus nonempty `name`, `run_in_background=true`, and absent `team_name`; it proves ensign identity independently from the prompt's `Skill(skill="spacedock:ensign")` and the on-disk member meta's `agentType`. It must reject a bare call, legacy `team_name`, missing ensign prompt/meta, or missing durable result.

**AC-4 — The correction does not weaken worker-completion evidence.**

The proof still requires the ensign skill, verbatim stage definition, one worker execution, fixture marker, exact `## Stage Report: implementation` shape with `DONE` and `Summary`, a path-scoped commit containing that result, a clean entity worktree, and bounded stop at `status: implementation`.

## Test plan

1. Add offline mode-aware transcript fixtures and mutation controls before changing the skill. A mode-preserving stream passes; deleting the only `Agent`, adding a second `Agent`, or swapping bare/team fields fails AC-1/AC-2. Cost: small table tests, no model.
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
