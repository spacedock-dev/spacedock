---
title: Select the actionable Codex default-headless task
status: validation
score: "0.90"
source: Repeated Codex wrong/queued target entry; DVD run 31432758302 artifact 9080028678; filing run 31434160297 artifact 9080564383, 2026-08-10
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-product
id: 272j6s25f9mry6nxbf4yjxvt
gates:
    version: 1
    records:
        - id: gate:272j6s25f9mry6nxbf4yjxvt:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:272j6s25f9mry6nxbf4yjxvt-backlog-1
              briefing:
                id: briefing:272j6s25f9mry6nxbf4yjxvt:backlog:attempt-1:revision-1
                digest: sha256:d276a71ce0905bc7454b3df2bf92eef7fbe8e938fc575908afa97d493ba37c21
                request-digest: sha256:38daadfeab1726b731607f66dd5ce060bd32fd2eb3bb1ef583818480aa604e69
                room-ref: ./select-actionable-codex-default-headless-task/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:272j6s25f9mry6nxbf4yjxvt:backlog:1
                briefing: briefing:272j6s25f9mry6nxbf4yjxvt:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-11T05:24:15.819247Z"
                decision: approve
                reason: The repaired seed isolates one product-only Codex target-selection outcome, excludes prohibited mechanisms and other hosts, and defines the required bound-to-unbound proof ladder.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:272j6s25f9mry6nxbf4yjxvt:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:272j6s25f9mry6nxbf4yjxvt-ideation-1
              briefing:
                id: briefing:272j6s25f9mry6nxbf4yjxvt:ideation:attempt-1:revision-1
                digest: sha256:a1f9e70c1b9d4caf0dae8e8ae24d7d81ba9892c9997b254233044947c1946e01
                request-digest: sha256:5a7226a4a0ebddb2c7914d1386c17d8515eee86a684139bbb7a18fa1560dd40c
                room-ref: ./select-actionable-codex-default-headless-task/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:272j6s25f9mry6nxbf4yjxvt:ideation:1
                briefing: briefing:272j6s25f9mry6nxbf4yjxvt:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-11T05:45:24.861062Z"
                decision: approve
                reason: The exact bound Codex XFAIL proves the missing native implementation spawn. The selected one-file declarative correction and binding-only removal ladder are the smallest product-only design within the approved caps and exclusions.
              application:
                target-stage: implementation
                state: consumed
started: 2026-08-11T05:24:54Z
worktree: .worktrees/spacedock-ensign-select-actionable-codex-default-headless-task
---
## Problem

Codex default-headless selects the correct implementation task from durable boot
and `status --next`, builds its dispatch envelope, but can skip `spawn_agent`.
It can then use an empty wait and a self-authored report as completion evidence,
advance to validation, and present a clean gate for work no native worker did.

## Value

A headless Codex user gets exactly one native implementation worker before the
task enters validation. The First Officer accepts only that worker's completion
signal and durable report, then commits and presents the clean validation gate
without consuming authority or dispatching a successor.

## Scope

- Change only the Codex First Officer runtime instruction at the fresh-spawn
  boundary. Keep the host-neutral dispatch and gate lifecycle unchanged.
- Use the existing public `spawn_agent` binding, returned handle, completion
  signal, stage report, and live journey. Add no new runtime mechanism.
- Remove only the Codex `default-headless-gate-stop` XFAIL row and its mirrored
  reconciliation row, and only after the bound target reports XPASS.
- Exclude hooks, observer references, temporary workflow state,
  transcript-driven product logic, instrumentation, and new standing lint or CI.
- Exclude Pi, Opus, Sonnet, target-only CI, PR work, command grammar, stored
  formats, gate authority, and unrelated host bindings.

## Exact bound baseline

The authorized local subscription run used clean detached source
`07ce3ddd30e644b289deda98d3a589ec18e57e41`, a binary built from that source,
`OPENAI_API_KEY` unset, isolated ChatGPT OAuth, and artifact root
`/tmp/272j-bound-codex-07ce.6AuoNv`.

The target completed in 312.42 seconds with status `XFAIL`, owner
`kky8pg7wc8xgb985epwss092`, and the single semantic code
`implementation-worker-not-dispatched`. The setup artifact records the exact
source SHA and `Logged in using ChatGPT`.

The command log proves that boot and `status --next --json` selected
`recorded-gate-task` correctly. Codex entered implementation and built its
implementation envelope. The next collaboration event was `wait` with empty
`receiver_thread_ids` and empty `agents_states`, not `spawn_agent`. Codex then
read a report, advanced to validation, and prepared and committed the open gate.
The current defect is therefore native dispatch/completion, not target
selection or gate preparation.

The earlier artifact `/tmp/272j-bound-codex.9vVlhU` is retained separately as
invalid infrastructure evidence. It ran stale source `ff9bb4506`, emitted no
Codex stream event, and timed out after 60 seconds. It is not a semantic
baseline.

## Selected approach

Tighten the existing `«worker.spawn»` bullet in the Codex runtime adapter. The
replacement will say, in the Codex host vocabulary, that after a zero-exit
fresh build the next host action is `spawn_agent`, its returned handle must
exist before any wait, edit, report read, or stage advance, and an unavailable
or handle-less spawn stops with a concrete blocker. It will explicitly reject
an empty wait or self-authored report as a substitute.

This serves AC-1 and AC-2. It introduces no new mechanism: the Codex adapter
already binds `spawn_agent`, and the host-neutral core already defines the
spawn, handle, completion, and report gates. The exact live baseline exercised
the riskiest path and proved that generic host-neutral wording alone is
insufficient at this Codex boundary.

The simplest alternative is no change because the generic core already names
the invariant. The exact run disproves that alternative. A binary receipt or
host-event observer is larger and would require a prohibited host observation
or transcript mechanism, so it is outside this task.

The approved product wording is exact:

```diff
-  A zero-exit build must carry `name` and `prompt`; map them through `CodexMultiAgentV2SpawnInput` to `spawn_agent(task_name,message,fork_turns="none")`.
+  A zero-exit build must carry `name` and `prompt`; map them through `CodexMultiAgentV2SpawnInput` to `spawn_agent(task_name,message,fork_turns="none")`. Its next Codex host action MUST be that `spawn_agent` call.
+  Record the returned handle before any wait, direct file change, report read, or stage advance.
+  If `spawn_agent` is unavailable or returns no handle, stop with the concrete blocker; an empty wait or self-authored report is not completion.
```

No documentation-site change is needed. The adapter is the shipped host
integration contract, and this task changes no command, output, format, or
documented runtime-live outcome.

## Expected surface

- `skills/first-officer/references/codex-first-officer-runtime.md`: replace the
  fresh-spawn paragraph near lines 19-20; estimate 3 insertions and 1 deletion,
  4 gross lines.
- `internal/ensigncycle/shared_live_runner_test.go`: remove only the Codex XFAIL
  binding after XPASS; estimate 1 insertion and 1 deletion.
- `internal/contractlint/live_registry_reconciliation_test.go`: mirror only that
  Codex binding removal; estimate 1 insertion and 1 deletion.
- Expected total: 3 files, 8 gross lines, 2 net lines. Hard cap: 3 files and 12
  gross lines; product cap: 1 file and 6 gross lines.

The declared semantic change is Codex First Officer runtime behavior only:
fresh dispatch cannot advance without a returned native worker handle and its
completion signal. Command grammar, stored formats, authority, other hosts,
and gate behavior do not change. Crossing any undeclared semantic boundary or
the file/line cap requires a new gate decision.

## Acceptance criteria

**AC-1 (VALUE): The exact Codex default-headless journey records one native implementation `spawn_agent`, its completion before validation, and one DONE implementation report.** The existing `assertImplementationWorkerLifecycle`
oracle measures this end value. Removing the spawn or moving validation before
completion restores `implementation-worker-not-dispatched`.

**AC-2: The journey ends at one committed, clean, open validation gate, with no decision, consumed authority, terminal fields, or successor dispatch.** The
existing `assertGateHeld` and `assertRecordedGateHoldLog` checks test the
resulting state and command order. A gate decision, post-prepare status change,
or successor build fails them.

**AC-3: The correction remains one Codex adapter change within the declared product cap and preserves all other runtime semantics.** `git diff --stat`,
`git diff --check`, and review against the expected surface test this boundary.
AC-3 serves AC-1; it is not satisfied if AC-1 is not green.

**AC-4: The Codex target progresses from bound XFAIL to bound XPASS, then to unbound PASS after only its two binding rows change.** Retain the exact bound
XPASS artifact before editing the rows. Run reconciliation after the edit, then
retain a fresh unbound normal-PASS artifact. Any semantic code, infrastructure
failure, or extra binding change fails this criterion.

**AC-5: The candidate passes one focused, full, and race ladder.** Run exactly:

```text
go test ./internal/contractlint -run '^(TestInitialWorkerSpawnGuardPrecedesCompletionAndValidation|TestCodexSpawnSignatureBindsToolArgs)$' -count=1
go test ./...
go test ./... -race
```

After these checks, run the exact bound local subscription target once with a
fresh artifact root and `OPENAI_API_KEY` unset. On XPASS, remove only the two
Codex binding rows, run `TestRuntimeLiveRegistryReconciliation`, and run the
same target once more for unbound PASS. Independent validation inspects the
exact diff and retained XFAIL, XPASS, and PASS artifacts; this task opens no PR
and starts no CI.

## Stage Report: backlog

- DONE: Seed end value
  Codex selects the actionable implementation task from durable boot and next state. It reaches a clean, open validation gate.
- DONE: Included scope
  The scope permits the smallest declarative Codex First Officer or public binary behavior correction. It uses existing public-behavior tests. The exact local subscription moves from bound XFAIL to XPASS, binding-only removal, and unbound PASS.
- DONE: Excluded scope
  The scope excludes global or host hooks, observer references, temporary state, transcript-driven product mechanisms, product instrumentation, and new standing CI or lint. It excludes target-only CI and work for Pi, Opus, or Sonnet.
- DONE: Proof plan
  Record the current bound local Codex XFAIL baseline. Run one focused, full, and race ladder, then record the exact XPASS. Remove only the Codex binding and reconciliation. Record the exact unbound PASS, independent validation, and required PR CI.

### Summary

The Captain's current directive supersedes the seed's historical n28 acknowledgment wording. This report defines the current boundary without changes to the seed body.

## Stage Report: ideation

- DONE: Prove the current bound local-subscription Codex default-headless target XFAIL before product edits, retaining exact artifact, status, and semantic codes.
  Exact source `07ce3ddd` produced XFAIL in `/tmp/272j-bound-codex-07ce.6AuoNv` with only `implementation-worker-not-dispatched`; no product bytes changed.
- DONE: Select the smallest declarative First Officer or public-binary behavior correction that follows durable boot/next task selection, with exact files, line estimate, tolerance, and prohibited-mechanism exclusions.
  The design sharpens one Codex adapter spawn boundary, then removes two evidence bindings within a 3-file/12-gross cap; every prohibited mechanism and other host is excluded.
- DONE: Rewrite the stale n28-era value and acceptance criteria into the Captain-approved product-only outcome, and define the single focused/full/race plus XPASS-to-unbound-PASS verification ladder.
  The body now measures native spawn/completion before validation, clean open-gate state, exact surface, focused/full/race checks, bound XPASS, binding-only removal, reconciliation, and unbound PASS.

### Summary

The exact baseline proves correct task selection but no native implementation spawn.
The selected correction is a bounded Codex adapter clarification backed by the existing live behavior oracle and an XPASS-to-PASS removal ladder.

## Stage Report: implementation

- DONE: Apply only the approved declarative Codex First Officer spawn-boundary correction within the one-file, six-gross product cap.
  Candidate `c716f2c524a438083dcb581b1c9efbe91e4b10b5` has sole parent/base `07ce3ddd30e644b289deda98d3a589ec18e57e41`; its exact diff is one file, `skills/first-officer/references/codex-first-officer-runtime.md`, with 3 insertions and 1 deletion (4 gross), and it remains unpushed/unmerged under HOLD.
- DONE: Run the exact focused contract checks, one full suite, and one race suite on the immutable product candidate.
  Focused `TestInitialWorkerSpawnGuardPrecedesCompletionAndValidation|TestCodexSpawnSignatureBindsToolArgs`, `go test ./...`, and `go test ./... -race` all passed on clean immutable head `c716f2c5`; removing either generic spawn/signature contract or introducing a suite/race regression would fail this ladder.
- FAILED: Run one bound local-subscription Codex default-headless target and retain exact XPASS evidence before any binding edit.
  The bound run used exact source `c716f2c5`, `OPENAI_API_KEY` unset, and ChatGPT auth, but ended XFAIL after 418.61s with `[gate-hold-violation implementation-worker-not-dispatched]`; retained artifact `/tmp/272j-bound-xpass-c716.E4N8Jj` would instead have reported XPASS only if the native worker lifecycle and gate hold were both observed.
- SKIPPED: Remove only the Codex default-headless XFAIL row and its mirrored reconciliation row, then run reconciliation and one fresh unbound local Codex target requiring normal PASS.
  No binding row was edited and no reconciliation or unbound run was performed because the prerequisite bound XPASS failed and the Captain directed HOLD/design reset.
- DONE: Write a complete implementation Stage Report with exact SHAs, diff, checks, retained artifacts, and explicit AC-1 through AC-5 evidence.
  This report records the candidate/base, exact capped diff, green offline ladder, failed bound proof, unchanged bindings, retained artifacts, and HOLD recommendation.

- FAILED: AC-1 — one native implementation `spawn_agent`, completion before validation, and one DONE implementation report.
  The live stream read the new MUST-spawn text verbatim but emitted no `spawn_agent`; after the zero-exit implementation build it claimed dispatch, issued two empty waits, directly read the pre-existing report, and advanced to validation.
- FAILED: AC-2 — one committed, clean, open validation gate without authority consumption or successor dispatch.
  The final gate was open, but `gate-hold-violation` proves the full required gate-hold command/state boundary was not satisfied, and validation followed false implementation completion evidence.
- FAILED: AC-3 — one bounded Codex adapter correction that preserves other semantics and serves green AC-1.
  The one-file/four-gross surface and `git diff --check` passed, but AC-3 is unmet by its own value re-anchor because AC-1 remained red.
- FAILED: AC-4 — bound XFAIL to XPASS, binding-only removal, reconciliation, then unbound PASS.
  The candidate remained bound XFAIL, so both Codex rows remain intact and the post-XPASS steps were correctly not run.
- DONE: AC-5 — focused, full, and race ladder.
  All three exact commands passed on immutable `c716f2c5`; any contract, package, or data-race failure would invalidate this evidence.

- FAILED: Missing supported primitive.
  The product lacks an existing public host/binary dispatch boundary that actually invokes native spawn or fail-closes stage advance when no native worker handle exists; declarative text alone is not that primitive, and prohibited hooks or observer mechanisms are not authorized.
- DONE: Retained evidence and disposition.
  Baseline `/tmp/272j-bound-codex-07ce.6AuoNv`, invalid-infrastructure baseline `/tmp/272j-bound-codex.9vVlhU`, and failed-candidate `/tmp/272j-bound-xpass-c716.E4N8Jj` remain preserved; recommendation is HOLD/design reset and candidate `c716f2c5` must not be pushed or merged.

### Summary

The approved declarative correction stayed inside its surface cap and passed every offline check, but the bound Codex journey read it and still skipped native dispatch before advancing on false evidence. AC-1, AC-2, AC-3, and AC-4 remain unmet; the Captain-directed outcome is HOLD/design reset until a supported fail-closed native dispatch primitive is approved.

## Stage Report: implementation (cycle 2)

- DONE: Verify the exact clean candidate and push its branch without changing candidate bytes.
  Local and remote `spacedock-ensign/grade-codex-session-lifecycle` both resolve to `229f5227ff127dbf2c813e7b239daa5a4cdd452d`, whose sole parent/base is `be0fd17c64081313ba7f5feeb9ed42bfb72f7049`.
- DONE: Keep the correction test-only and record its exact surface.
  The four files are `internal/ensigncycle/claude_live_runner_test.go`, `internal/ensigncycle/claude_runtime_helpers_test.go`, and the two `testdata/codex_native_lifecycle/{public,parent-rollout}.jsonl` fixtures; the diff is +93/-17.
- DONE: Preserve focused, full, and race verification evidence without rerunning it.
  Focused lifecycle/lookup falsifiers, `go test ./...`, and `go test ./... -race` passed on `229f5227f`; removing handle correlation, accepting ambiguous lookup, or introducing package/race failure makes this ladder red.
- DONE: Prove the native Codex primitive from isolated plain-host evidence.
  `/tmp/codex-v2-spawn-probe.hXx1Qu` maps each public parent thread ID to one internal rollout and records spawn < returned handle < child completion at 13<15<24 and 13<15<23; public stdout omits those authoritative rows.
- FAILED: Prove the corrected bound Spacedock target as XPASS.
  Exact candidate `229f5227f` remained bound XFAIL after 282.48s with `implementation-worker-not-dispatched`; `/tmp/272-session-grade-live.UuY9E5` retains the result, and diagnostic internal-session evidence showed the exact Spacedock run genuinely lacked `spawn_agent`.
- DONE: Preserve every product, hook, skill, state-protocol, and binding byte.
  The candidate changes only test harness and fixture files; both Codex default-headless XFAIL rows remain unchanged.
- SKIPPED: Remove bindings, run an unbound target, rerun live/full/race, open a PR, or advance workflow state.
  The bound prerequisite stayed red and the assignment explicitly forbids these actions.

- FAILED: AC-1 — exactly one native implementation worker completes before validation with one DONE report.
  The corrected evidence source no longer mistakes public-stream omission for absence, but the exact bound Spacedock run still had no native spawn, so AC-1 remains red.
- DONE: AC-2 — one committed, clean, open validation gate with no consumed authority or successor dispatch.
  The 282.48s corrected bound grade reported only `implementation-worker-not-dispatched`; the gate-held assertions emitted no AC-2 semantic code and would fail on decision, consume, post-prepare mutation, or successor build.
- FAILED: AC-3 — the bounded correction serves green AC-1 while preserving runtime semantics.
  Runtime semantics are preserved by a test-only diff, but AC-3 remains unmet by its value re-anchor because AC-1 is red.
- FAILED: AC-4 — bound XPASS, binding-only removal, reconciliation, then unbound PASS.
  The first step remained XFAIL, so both rows remain and every post-XPASS step was correctly skipped.
- DONE: AC-5 — focused, full, and race ladder.
  All three owned offline checks passed on `229f5227f`; the assignment forbids repeating them at this completion boundary.
- DONE: Write and push this complete cycle-2 implementation report without advancing state.
  This report records the pushed code SHA, exact test-only diff, falsifiable checks, retained probes, corrected XFAIL, unchanged bindings, and revised classification.

### Summary

The test observability defect is fixed: Codex lifecycle grading now uses the isolated parent rollout and correlates spawn to its returned handle before accepting child completion. The prior missing-native-primitive conclusion is withdrawn, but AC-1, AC-3, and AC-4 remain unmet because the exact bound Spacedock run still genuinely lacked spawn; checklist count is 8 DONE, 1 SKIPPED, and 4 FAILED.

## Stage Report: implementation (cycle 3)

- DONE: Verify the exact base and push the final candidate without changing product bytes.
  Base is `be0fd17c64081313ba7f5feeb9ed42bfb72f7049`; local and remote `spacedock-ensign/grade-codex-session-lifecycle` both resolve to `2b660fb6fb6ca6566abe086c6d7ee3adcee30eb0`.
- DONE: Correct Codex lifecycle grading from the isolated parent session log and retain only the first correlated implementation completion.
  Commit `d7b0991810c52f8f1a52a59226670b07688d0534` adds the `completed < 0` guard; the fixture passes with implementation `Done:` before validation plus a later same-handle validation `Done:`, while no completion and after-validation-only completion stay red.
- DONE: Keep the final correction test-only and record its exact surface.
  The four lifecycle files are `internal/ensigncycle/{claude_live_runner_test.go,claude_runtime_helpers_test.go}` and `internal/ensigncycle/testdata/codex_native_lifecycle/{public,parent-rollout}.jsonl`; only `internal/ensigncycle/shared_live_runner_test.go` and `internal/contractlint/live_registry_reconciliation_test.go` add the two row changes, for six files and `+105/-20` from base.
- DONE: Run one focused lifecycle ladder, one full suite, and one race suite on the immutable correction bytes.
  The correlated-session/lookup/negative-control focus passed in 1.71s, `go test ./...` passed in 280.95s, and `go test ./... -race` passed in 332.01s; removing first-completion retention, handle correlation, package correctness, or race safety makes the corresponding rung red.
- DONE: Run one bound local Codex default-headless target and retain XPASS evidence before binding edits.
  Exact head `d7b099181` reported `XPASS ALERT` with no semantic codes in 307.25s; internal evidence records one spawn at parent line 77, its correlated handle at 79, implementation `Done:` at 93, and validation at 103.
- DONE: Remove only the Codex default-headless binding and mirrored reconciliation row, then verify registry ownership and formatting.
  Commit `2b660fb6f` changes only those two rows; registry/binding-format passed in 1.59s, the mutable active-owner join passed in 2.01s, and gofmt plus `git diff --check` passed.
- DONE: Run one fresh unbound local Codex default-headless target requiring normal PASS.
  Exact final head `2b660fb6f` reported normal PASS with no binding annotation in 724.50s; its independent parent records one spawn at line 52, handle at 54, implementation `Done:` at 66, and validation at 76.
- DONE: Retain the exact offline, bound, and unbound artifacts.
  Offline logs are `/tmp/272-codex-transition.TcjLIT/artifacts/correction/{focused,full,race,unbind-registry-format,unbind-owner}.log`; bound evidence is under `/tmp/272-codex-transition.TcjLIT/artifacts/bound-xpass` and unbound evidence under `/tmp/272-codex-transition.TcjLIT/artifacts/unbound-pass`, including each `go-test.log`, `internal-lifecycle-verdict.txt`, `thread_spawn_edges.json`, `internal/sessions/`, `internal/state_5.sqlite`, and live process artifacts.
- DONE: Preserve product, skill, hook, command, state-protocol, authority, and unrelated runtime bytes.
  The final branch diff contains only the six named test files; no CI or PR was started, and the earlier adapter-text hypothesis is superseded by authoritative session evidence rather than shipped as a product change.

- DONE: AC-1 — exactly one native implementation worker completes before validation with one DONE report.
  Both preserved lifecycle verdicts independently prove one spawn, its returned handle, one matching implementation `Done:` before validation, and the durable report oracle; deleting spawn/completion or moving completion after validation restores red.
- DONE: AC-2 — one committed, clean, open validation gate without authority consumption or successor dispatch.
  Bound XPASS and unbound PASS both emitted an empty semantic set, so `assertGateHeld` and `assertRecordedGateHoldLog` accepted the gate; any decision, consume, terminal field, post-prepare mutation, or successor build would emit a gate-hold failure.
- DONE: AC-3 — the bounded correction preserves all other runtime semantics and serves green AC-1.
  The final correction is stricter than the declared product cap: zero product files and only the six test surfaces; AC-1 is green on two live runs, so the preservation/value re-anchor is satisfied without the rejected adapter prose.
- DONE: AC-4 — bound XPASS, binding-only removal, reconciliation, then unbound PASS.
  Preserved `bound-xpass/internal-lifecycle-verdict.txt` reports XPASS before commit `2b660fb6f`; only the two declared rows then changed, focused reconciliation passed, and `unbound-pass/internal-lifecycle-verdict.txt` reports normal PASS.
- DONE: AC-5 — focused, full, and race ladder.
  All three owned offline rungs passed once on final correction bytes before the bound run; no check or live target was retried after a failed rung because none failed.

### Summary

Cycle 3 closes the observability defect and the live transition without changing product semantics: the grader reads the parent session, correlates the returned worker handle, and freezes the first implementation completion. The exact bound target XPASSed, only its two bindings were removed, and a fresh unbound target passed; checklist count is 14 DONE, 0 SKIPPED, and 0 FAILED.
