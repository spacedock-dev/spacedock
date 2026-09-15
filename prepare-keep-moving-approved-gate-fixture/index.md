---
title: Prepare the approved gate in the keep-moving live fixture
status: validation
source: "Captain filing from PR #784 runtime CI attempt 2, Codex job 103174402234"
id: wvdv5c5mkkd23yhmr9p16czd
gates:
    version: 1
    records:
        - id: gate:wvdv5c5mkkd23yhmr9p16czd:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:wvdv5c5mkkd23yhmr9p16czd-backlog-1
              briefing:
                id: briefing:wvdv5c5mkkd23yhmr9p16czd:backlog:attempt-1:revision-1
                digest: sha256:4267b265607a4a9ccf99397d80bb759eafbd76382dfeb268c97d9cce9ad97e4c
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:wvdv5c5mkkd23yhmr9p16czd:backlog:1
                briefing: briefing:wvdv5c5mkkd23yhmr9p16czd:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-09-15T04:01:30.288629Z"
                decision: approve
                reason: Captain requested dispatch of the Codex live failure tasks, targeted local verification, and stacked PRs.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:wvdv5c5mkkd23yhmr9p16czd:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:wvdv5c5mkkd23yhmr9p16czd-ideation-1
              briefing:
                id: briefing:wvdv5c5mkkd23yhmr9p16czd:ideation:attempt-1:revision-1
                digest: sha256:10a5618ff09104122df3c972d5d2a5df112cc6128e553c98aa7df024b8b36bb1
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:wvdv5c5mkkd23yhmr9p16czd:ideation:1
                briefing: briefing:wvdv5c5mkkd23yhmr9p16czd:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-09-15T04:07:04.691315Z"
                decision: approve
                reason: Minimal real gate setup is proven; bounded fixture correction retains authority and dispatch proof. Proceed under captain scope through targeted verification and PR stack.
                conn:
                    quote: dispatch codex live test failure tasks, verify locally for targeted failure, and open PR as stack, then trigger codex ci on stack tip.
                    source: active thread goal supplied by captain
              application:
                target-stage: implementation
                state: consumed
        - id: gate:wvdv5c5mkkd23yhmr9p16czd:validation
          stage: validation
          attempts:
            - id: gate-attempt:wvdv5c5mkkd23yhmr9p16czd-validation-1
              briefing:
                id: briefing:wvdv5c5mkkd23yhmr9p16czd:validation:attempt-1:revision-1
                digest: sha256:7fb53e4a5d6e91e2a7c1d270fe3045c97254b47c4d84fbe5c014a3c729a52615
                room-ref: '@review/validation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:wvdv5c5mkkd23yhmr9p16czd:validation:1
                briefing: briefing:wvdv5c5mkkd23yhmr9p16czd:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-09-15T04:33:19.834871Z"
                decision: approve
                reason: Independent targeted local Codex and real gate controls satisfy all layer ACs at7ca67fdda. Captain authorized opening the verified stack; full-stack CI and final merge remain separate.
                conn:
                    quote: dispatch codex live test failure tasks, verify locally for targeted failure, and open PR as stack, then trigger codex ci on stack tip.
                    source: active thread goal supplied by captain
              application:
                target-stage: done
                state: pending
started: 2026-09-15T04:01:46Z
worktree: .worktrees/spacedock-ensign-prepare-keep-moving-approved-gate-fixture
mod-block: merge:pr-merge
pr: pr-merge:794
---

## Problem statement

The host-neutral keep-moving fixture says the captain just approved `approved-gate` at review, but seeds no gate record or bound briefing. A conforming `gate record --decision approve --consume` therefore fails before the continuation behavior under test can begin. The test grades the missing implementation dispatch as a keep-moving violation.

This task belongs to the development workflow because its deliverable is a corrected executable live-test fixture and regression proof in `internal/ensigncycle`. Keep the runtime's refusal of unprepared gates intact. This is separate from the existing Codex smallest-mechanism grader repair and from Pi or Sonnet posture defects.

## Evidence

- Run: https://github.com/spacedock-dev/spacedock/actions/runs/34568401666/attempts/2 ; Codex job `103174402234`; candidate `09f123d33b1d7a5b8d3d2dbc85e536739ba7ea1b`.
- Artifact `10188347210`, created 2026-09-11T07:03:37Z; member `live-artifacts/codex/codex-shared-scenarios/keep-moving-posture/codex-exec.jsonl`. Local archive `/tmp/pr784-codex-attempt2.zip` is only a temporary convenience.
- Test `TestLiveCommonKeepMovingPosture` failed after 158.91s with `durable keep-moving journeys = 2/3: map[approved-gate:missing path-scoped dispatch entry with stage and started]`.
- Parent trace shows `gate record approved-gate --decision approve --actor person:captain --reason "Captain approved review in the current request." --consume --workflow-dir <fixture-root>` exiting 1 with `Error: entity has no gates record`.
- The immediately preceding file read proves `approved-gate.md` contains id/title/status=review and no gates record. Codex continues the two independent tasks to archive and reports the approval blocker. This is not an auth error or a provider timeout.
- Source: `writeKeepMovingWorkflow`, `keepMovingApprovedEntity`, and `keepMovingPrompt` in `internal/ensigncycle/shared_fixtures_test.go`; grading in `assertDurableKeepMoving` in `internal/ensigncycle/shared_keep_moving_durable_test.go`.

## Scope and proposed direction

Construct the real prepared review attempt and committed briefing through the supported gate API before presenting the captain's approval to the FO. Keep approval consumption and successor dispatch as observed agent actions, not fixture-precompleted work. Confirm the exact preparation requirements before choosing the smallest fixture setup. Do not weaken approval guards, force the transition, or classify a missing precondition as an agent continuation failure.

The observed setup defect is established; whether additional agent failures remain after repairing it needs the targeted live run. This filing does not claim PR #784 is cleared for merge.

## Acceptance criteria

**AC-1 — The fixture's stated approval can be consumed through the supported gate lifecycle.**
Verified by: a deterministic setup exercise prepares the actual gate, records and consumes approval successfully, and observes implementation as the successor; omit the prepared record and the same consume must fail.

**AC-2 — The live scenario still measures continuation and independent work.**
Verified by: the targeted Codex keep-moving run observes the approved task's implementation dispatch, worker report, and terminal journey, alongside the independent tasks and nonterminal questioned task. Existing durable journey checks remain active; suppressing the approved task's dispatch must still fail the grader.

**AC-3 — Approval authority is not bypassed by fixture preparation.**
Verified by: inspect and exercise the initialized fixture before the captain decision; it has a prepared review attempt but no consumed approval or implementation dispatch. A missing approval cannot advance it. Do not use forced status changes as preparation.

## Stage-specific test gates

- Ideation: reproduce the seed/consume mismatch with the installed candidate and demonstrate the smallest valid gate setup in a throwaway fixture. Name the existing primary proof owners and expected file/LOC scope.
- Implementation: add the focused regression before the fixture change; run the relevant deterministic tests and format changed Go files. Per captain-approved stack scheduling, full Go/race suites and repository formatting run once at the combined stack tip.
- Validation: independently verify setup and refusal behavior, then run the targeted Codex scenario. Any remaining live failure must be diagnosed separately from the repaired setup.

## Implementation plan and expected surface

Use the existing single-root `writeKeepMovingWorkflow` fixture. Add one short review Markdown artifact before `gitInit`; after initialization run the real `gate prepare approved-gate --question "Advance to implementation?" --artifact <root>/gate-review.md --summary "Ready to proceed to implementation." --workflow-dir <root>`, then `state commit approved-gate --workflow-dir <root>`. Reuse `buildRecordedGateBinary` and `mustRecordedGate`; no new framework or split-root fixture is needed. The agent still owns recording the captain decision, consuming it, dispatching implementation, and completing the durable journey.

This preparation serves AC-1 and AC-3. Bare status seeding is insufficient because it has no gate record; hand-authored gate YAML would duplicate digest/binding machinery. The supported prepare and state-commit commands are the smallest exercised route to a real committed review package.

Estimate net LOC change: +100, across 3 files (100 insertions, 0 deletions), tolerance +50 net LOC and no additional files. Expected files: `internal/ensigncycle/shared_fixtures_test.go` (about 10 lines), new `internal/ensigncycle/shared_keep_moving_fixture_test.go` (about 80), and `internal/ensigncycle/shared_keep_moving_durable_test.go` (about 10 for an approved-task missing-dispatch control). Changes affect only test fixture initial state and regression coverage; CLI grammar, stored formats, approval authority, production runtime, skills, and durable grader semantics stay outside scope. No user-visible documentation diff is needed.

### Spike evidence and primary proof

At `origin/main` 2a7b8719843e40b79545f0bb4def6609cdd9ebbf, the installed pre3 binary reproduced unprepared `gate record --decision approve --actor person:captain --consume` exiting 1 with `entity has no gates record`. The throwaway `TestKeepMovingIdeationSpike` then proved prepare plus state commit retains `approved-gate/review/review/briefing-1/index.json` in HEAD; the prepared entity remains at review with no started, resolution, or application. Missing-decision consume exits 1, reports ineligible/consumed=false, and leaves entity bytes unchanged. Recording captain approval with consume exits 0 and writes status=implementation/application.state=consumed.

The spike also established that single-root `gate prepare` alone leaves its room uncommitted; `state commit` is required for the committed-fixture requirement. Retained local evidence: `/tmp/keep-moving-ideation-wvdv/internal/ensigncycle/keep_moving_spike_test.go` and `/tmp/keep-moving-gate-spike-wvdv4` (throwaway paths, not dependencies of shipped tests). Final exercised command: `go test ./internal/ensigncycle -run '^TestKeepMovingIdeationSpike$' -count=1 -v` (1.446 seconds, pass).

Primary proof owners are the existing real command helpers in `recorded_gate_lifecycle_test.go`, `TestDurableTaskJourneys`, and `TestLiveCommonKeepMovingPosture` with `assertDurableKeepMoving`. Add the focused `TestKeepMovingPreparedGate` before the fixture edit: verify committed open preparation, unchanged missing-decision refusal, successful captain approval consumption, and the unprepared refusal control. Removing preparation, preconsuming approval, or omitting the state commit must each fail it. Use on-disk parsed state, command exit codes, Git objects, and byte equality; no prose-grep proof.

AC-2 stays with the existing durable journey proof and one missing-approved-dispatch case in its existing test table: suppressing that dispatch must fail while preserving the other two journeys. The live test must still observe all three dispatch/report/terminal journeys and the questioned task's nonterminal correction. Deterministic checks cost seconds plus one binary build; the single Codex live run costs minutes and is scheduled serially by the FO after implementation.

Targeted deterministic command: `go test ./internal/ensigncycle -run '^(TestKeepMovingPreparedGate|TestDurableTaskJourneys|TestDurableKeepMoving.*|TestDurableQuestioned.*|TestRetainedAtomicWorkerJourney)$' -count=1`.

Exact targeted Codex command, from the prepared stack checkout: `SPACEDOCK_BIN=/opt/homebrew/Caskroom/spacedock@next/0.28.0-pre3/spacedock SPACEDOCK_LIVE_RUNTIME=codex SPACEDOCK_CODEX_LIVE_REQUIRED=1 go test -tags live ./internal/ensigncycle -run '^TestLiveCommonKeepMovingPosture$' -count=1 -v -timeout 15m`. This requires Codex on PATH and the existing supported local auth or CI auth input; the harness handles isolation. The candidate changes only Go test fixtures, so the specified pre3 executable is sufficient. FO runs the full normal/race/CI checks once at the combined four-fix stack tip.

## Stage Report: ideation

- DONE: Prove the smallest real prepared-gate fixture setup and preserve the missing-approval refusal.
  Installed pre3 spike passed: unprepared refusal, real prepare + state commit, committed briefing, unchanged missing-decision refusal, then approval consumption to implementation; dropping state commit failed the Git-object assertion.
- DONE: Record a concise implementation plan, exact targeted Codex command, primary tests, and bounded surface for the captain-approved stack.
  Plan above bounds fixture/regression work to +100 net LOC across 3 files (tolerance +50); existing durable checks and serialized targeted Codex run own continuation proof.

### Summary

The fixture needs a real prepared review room and its state commit before the prompt grants approval. The approval decision and continuation remain agent actions; no production authority or grader relaxation is proposed. Live execution and combined-stack full suites remain scheduled after implementation.

## Stage Report: implementation

- DONE: Implement the real prepared and committed review gate fixture with focused positive and missing-authority regressions, preserving agent-owned consume and continuation.
  Commit `7ca67fdda5433654b6f67276560b3eb19c0684f6` adds real prepare + state commit; the initial regression failed with `entity has no gates record` before the fixture change.
- DONE: Commit the bounded layer from current origin/main and report focused test evidence and the exact targeted local Codex command.
  Branch starts at `2a7b8719843e40b79545f0bb4def6609cdd9ebbf`; 3 files, 74 additions/1 deletion (+73 net), within the approved +100/+50 estimate.
- DONE: Verify preparation, authority refusal, and successful consumption through real commands and durable state.
  `TestKeepMovingPreparedGate` asserts an open bound attempt, committed briefing/entity, no successor work in history, unchanged missing-decision and unprepared refusals, and captain approval consumed to implementation; removing preparation/commit or preconsuming approval makes it fail.
- DONE: Preserve approved-task dispatch proof.
  `TestDurableTaskJourneys/missing_approved_dispatch` suppresses only that task's dispatch and requires exactly 2/3 journeys and one approved-task failure; accepting its absent dispatch makes the control fail.
- DONE: Run the relevant deterministic tests and format changed Go files.
  `go test ./internal/ensigncycle -run '^(TestKeepMovingPreparedGate|TestDurableTaskJourneys|TestDurableKeepMoving.*|TestDurableQuestioned.*|TestRetainedAtomicWorkerJourney)$' -count=1` passed (162.906s); changed Go files formatted and `git diff --check` passed.
- SKIPPED: Run targeted live Codex and full normal/race suites in this worker.
  Captain-approved scheduling assigns serialized live validation and combined-stack full suites to the FO; no live run, code push, or PR was initiated here.

### Summary

The fixture now supplies a real committed review package while approval consumption, implementation dispatch, and completion remain agent actions. Production gate guards and durable grader semantics are unchanged; this is stack layer 1, ready for independent validation.

Exact targeted local command from the stack checkout: `SPACEDOCK_BIN=/opt/homebrew/Caskroom/spacedock@next/0.28.0-pre3/spacedock SPACEDOCK_LIVE_RUNTIME=codex SPACEDOCK_CODEX_LIVE_REQUIRED=1 go test -tags live ./internal/ensigncycle -run '^TestLiveCommonKeepMovingPosture$' -count=1 -v -timeout 15m`.

## Review-finding disposition

### Validation observation: local Codex command-host negotiation

- Released user and normal workflow: captain-requested local Codex Luna/max keep-moving validation through the supported isolated OAuth harness, exact candidate `7ca67fdda5433654b6f67276560b3eb19c0684f6`.
- Observable harm: all three executable tool calls timed out before reading instructions or workflow state; no journey ran, so this attempt cannot establish continuation behavior.
- Authority: value-ac[AC-2] targeted live evidence must observe all three durable journeys and questioned-task correction.
- Trigger evidence: Codex CLI `0.154.0`, thread `01a0a347-39a3-7241-b187-af7fec8cce5b`; `codex-exec.stderr.txt` records negotiation timeouts at 04:16:29, 04:17:01, and 04:17:37 UTC on 2026-09-15; test failed 0/3 after 120.546 seconds.
- Validator proposal: material evidence defect; ownership is local runtime infrastructure, outside this fixture-only candidate; hold AC-2 evidence and recover the host through an isolated tool smoke before a separately authorized rerun. No product or grader change is proposed.
- FO consultation: FO agreed this is runtime infrastructure and authorized non-mutating diagnosis, with no retry while the run remained live. The run is now terminal; recovery/rerun disposition remains with FO.

## Stage Report: validation

- DONE: Independently verify real gate setup and no-authority boundaries at layer1 commit 7ca67fdda; inspect focused proof and attack unowned claims without repeating already-green tests gratuitously.
  Clean exact HEAD; `go test ./internal/ensigncycle -run '^TestKeepMovingPreparedGate$' -count=1 -v` passed (3.878s), exercising committed preparation, unchanged refusal, and real approval consumption.
- FAILED: Run the targeted local Codex keep-moving journey with the CI model configuration, preserve artifacts and exact SHA, and assess every AC without changing candidate code.
  Session 79626 terminated with test exit 1, no skip; three command-host negotiation timeouts yielded 0/3 journeys before any gate command. Artifacts: `/tmp/spacedock-stack-keep-moving-live`.
- DONE: AC-1 — The fixture's stated approval can be consumed through the supported gate lifecycle.
  Independent focused execution consumed captain approval to implementation; unprepared control refused unchanged. Removing preparation fails the same test before consumption.
- FAILED: AC-2 — The live scenario still measures continuation and independent work.
  Existing missing-approved-dispatch control retains exact 2/3 grading; unchanged live grader still requires three ordered journeys, overlap, and meaningful nonterminal correction, but this runtime attempt could not execute them.
- DONE: AC-3 — Approval authority is not bypassed by fixture preparation.
  Focused execution proved one open review attempt, committed briefing/entity, no started/successor history, and byte-identical missing-approval refusal; preconsumption or omitting state commit would fail.
- DONE: Semantic adversarial pass and scope review.
  Matrix covered unprepared, prepared/open, missing approval, consumed approval, and missing approved dispatch; production guards and grader unchanged. Real command parsing, Git objects, byte equality, exact actor/target, and ordered history provide independent assertions. No changed hot-path scaling risk.
- SKIPPED: Full normal/race suites and repository-wide formatting.
  Captain-approved scheduling assigns these once to the combined stack tip; candidate remains unchanged at 3 files, +73 net LOC.

### Summary

Recommendation: REJECTED for incomplete AC-2 live evidence, held on local runtime infrastructure; AC-1 and AC-3 pass independently, and no fixture defect was found. The live run used a freshly built candidate binary, CI Luna/max shim, isolated local OAuth with inherited API key/required/repo overrides unset, and retained exact source-head, stderr, final message, and process-result artifacts. The harness removed its isolated auth/config and rollout during cleanup; no shared credentials/configuration were changed, and no retry, candidate edit, PR, or merge occurred.

## Stage Report: validation (cycle 2 — authorized local-host recovery)

- DONE: Independently verify real gate setup and no-authority boundaries at layer1 commit 7ca67fdda; inspect focused proof and attack unowned claims without repeating already-green tests gratuitously.
  Prior independent `TestKeepMovingPreparedGate` passed (3.878s); unchanged candidate and current clean HEAD remain `7ca67fdda5433654b6f67276560b3eb19c0684f6`.
- DONE: Run the targeted local Codex keep-moving journey with the CI model configuration, preserve artifacts and exact SHA, and assess every AC without changing candidate code.
  Authorized retry passed `TestLiveCommonKeepMovingPosture` (306.97s; package 307.454s), session 31189 exit 0; `/tmp/spacedock-stack-keep-moving-live/retry-1/test.log` and `_setup/keep-moving-posture/source-head.txt` retain result and exact candidate.
- DONE: AC-1 — The fixture's stated approval can be consumed through the supported gate lifecycle.
  Independent deterministic positive/unprepared controls passed; retry also executed real captain `gate record ... --consume` successfully before implementation dispatch.
- DONE: AC-2 — The live scenario still measures continuation and independent work.
  Unchanged live grader passed all three ordered dispatch/report/terminal/archive journeys plus overlap and meaningful historically nonterminal questioned correction; removing the approved dispatch still has the existing exact 2/3 negative control.
- DONE: AC-3 — Approval authority is not bypassed by fixture preparation.
  Focused proof verifies committed open preparation, no successor history and byte-identical missing-approval refusal; live approval was recorded by the first officer after fixture initialization.
- DONE: Resolve the local-host evidence defect without candidate changes.
  FO authorized one isolated tool smoke and one retry on smoke success. Current CLI 0.154.0 Luna/max smoke executed `printf HOST-SMOKE-OK` exit 0 (PID 59772), then the targeted retry passed with the same cleaned outer-session environment.
- DONE: Preserve diagnosis and recovery evidence.
  Original failure remains under `/tmp/spacedock-stack-keep-moving-live`; `host-smoke` retains environment-summary, stdout, stderr, result, and rollout; `retry-1` retains live artifacts plus sampled rollouts/workflow history. Isolated auth/config were cleaned; shared auth/config were untouched.
- SKIPPED: Full normal/race suites and repository-wide formatting.
  Captain scheduling keeps these at the combined stack tip; no candidate edits, additional live tests, PR, or merge in this worker.

### Summary

Recommendation: PASSED; all three ACs have behavioral evidence, and the prior local-host evidence hold is resolved by the authorized successful retry. No material candidate findings or deferred candidate risks remain. The recovery removes inherited runtime markers as a set; it does not establish which individual marker caused the original negotiation timeout.

Reproduction environment: remove `OPENAI_API_KEY`, `CODEX_AUTH_JSON`, `CODEX_CI`, `CODEX_THREAD_ID`, `CODEX_SESSION_ID`, `NODE_REPL_TRUSTED_BROWSER_CLIENT_SHA256S`, `SPACEDOCK_CODEX_LIVE_REQUIRED`, and `SPACEDOCK_REPO_ROOT`; use `PATH=/tmp/spacedock-codex-fixes-ci-shim:/Users/clkao/go/bin:$PATH`, `SPACEDOCK_CODEX_REAL_BIN=/opt/homebrew/bin/codex`, `SPACEDOCK_BIN=/tmp/spacedock-stack-keep-moving-live/bin/spacedock` (built from this exact candidate), `SPACEDOCK_LIVE_RUNTIME=codex`, and `SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/spacedock-stack-keep-moving-live/retry-1`. Command: `go test -tags live ./internal/ensigncycle -run '^TestLiveCommonKeepMovingPosture$' -count=1 -v -timeout 15m`. The shim is the verbatim origin/main CI Luna/max wrapper, with its source and local location explicitly supplied by the FO.
