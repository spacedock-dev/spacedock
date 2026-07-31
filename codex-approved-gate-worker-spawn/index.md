---
title: Require Codex worker spawn after approved gate advance
status: validation
source: "KD validation cycle 4, Codex keep-moving-posture at b60d1c8"
started: 2026-07-30T23:39:35Z
completed:
verdict:
score: 0.95
worktree: .worktrees/spacedock-ensign-codex-approved-gate-worker-spawn
issue:
milestone: 0.27.0
id: f02j6dbnd4jakwczahv1tg2h
gates:
    version: 1
    current:
        gate: gate:f02j6dbnd4jakwczahv1tg2h:ideation
    records:
        - id: gate:f02j6dbnd4jakwczahv1tg2h:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:f02j6dbnd4jakwczahv1tg2h-backlog-1
              briefing:
                id: briefing:f02j6dbnd4jakwczahv1tg2h:backlog:attempt-1:revision-1
                digest: sha256:7deff3c5b9e3616c2f7619faf23da6168456c71d8cd968d5c68a3a73c08306f7
                digest-domain: canonical-bytes
                request-digest: sha256:bd94ed0c5d9047c300db27b7014a1289f1bedd84f037fa6ecf938bd72bf53150
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:f02j6dbnd4jakwczahv1tg2h:backlog:1
                briefing: briefing:f02j6dbnd4jakwczahv1tg2h:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-30T23:39:29.432938Z"
                decision: approve
                reason: Captain conn approves immediate ideation because the exact supported Codex run already proves a Material worker-authority breach; the task owns diagnosis and the smallest runtime/contract correction, while KD remains frozen.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:f02j6dbnd4jakwczahv1tg2h:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:f02j6dbnd4jakwczahv1tg2h-ideation-1
              briefing:
                id: briefing:f02j6dbnd4jakwczahv1tg2h:ideation:attempt-1:revision-1
                digest: sha256:82e64986e9e8a9be0dff8d59e0a6d59bb9ac21759ba9878ab6b3027a007e6244
                digest-domain: canonical-bytes
                request-digest: sha256:cdcfda518e8129d2f9efe37b76f885b6735371118b21d14679a68f4b110b3b21
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:f02j6dbnd4jakwczahv1tg2h:ideation:1
                briefing: briefing:f02j6dbnd4jakwczahv1tg2h:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-31T00:01:28.193473Z"
                decision: approve
                reason: 'Captain conn approves the reproduced seven-file correction: Codex bare mode cannot supply the required task_name, so the binary must reject it before artifact creation and the adapter must use the existing named spawn shape without adding an observer dialect.'
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
---

A supported headless Codex keep-moving journey advanced `approved-gate`, built its implementation dispatch artifact, never invoked `spawn_agent` / `worker.spawn`, then read a report and terminalized the task. This violates the shipped dispatch boundary and silently bypasses worker execution and write-scope authority.

The exact evidence is retained at `/tmp/spacedock-validation-cycle4-b60d1c8.qPdIXO/codex/codex-shared-scenarios/keep-moving-posture/codex-exec.jsonl`: item 17 advances and builds; items 26–27 read and finalize; no worker identity or spawn exists. The live harness reported: `the FO advanced the approved entity "approved-gate" but did not dispatch its next stage`.

Ideation must find the smallest fix that makes the supported Codex journey actually spawn and await a worker after gate approval. It must not repair the observer, parse transcript/provider dialects, or let a dispatch commit or generated artifact substitute for an observed worker spawn.

## Acceptance criteria

**AC-1 — Approved work executes under worker authority.** In a supported Codex live validation, every ready implementation task, including one made ready by gate approval, is present in the live runtime roster as an addressable worker after its named build and before the first report read; the worker then produces the durable stage report/state used for completion. Test once on the exact candidate by observing the runtime roster/addressable handle in-session, then verifying the worker-authored report, commit, and resulting workflow state; retain that run as validation evidence, not as a committed event dialect.

**AC-2 — Missing spawn remains a hard failure.** Codex rejects `bare_mode=true` before writing a dispatch artifact, while a one-off control that successfully builds the named artifact but deliberately omits the native spawn has no matching live roster worker and no worker-authored report/state, so validation stops before report reads or completion credit. Test the mechanical arm with a dispatch-build unit test expecting exit 2 and no artifact; test the successful-build/no-spawn arm during the retained live validation by checking the roster first and refusing to proceed when the addressable worker is absent.

**AC-3 — The correction stays runtime-coherent without owning an observer.** The implementation uses the existing Codex build-to-`CodexMultiAgentV2SpawnInput`/`spawn_agent` binding and adds no keep-moving observer edit, transcript grammar, provider-event parser, receipt, stored schema, retry loop, or second protocol. Test with the existing Codex adapter/contractlint suite plus a changed-files review; if durable automated attribution remains required, it is explicitly blocked on 8b (`codex-keep-moving-durable-evidence-attribution-flake`) and is not implemented here. This mechanism serves AC-1.

**AC-4 — Supported host semantics do not regress.** Claude bare dispatch remains anonymous, blocking, and sequential; Pi keeps forwarding its assignment through its native substrate; named Codex dispatch still emits the unchanged `name` and `prompt` fields. Test the existing Claude merged/bare, Pi host, Codex host, and dispatch parity suites.

## Problem and root cause

The failure is not a missing sentence. The shared core and Codex adapter already say a successful build is not a dispatch and require `spawn_agent`. In the recorded run, however, the Codex FO invoked:

`dispatch build ... --stage implementation ... --bare-mode`

`bare_mode` is declared as the anonymous sequential shape: no `name`, no `team_name`, and no `run_in_background`. `internal/dispatch/build.go` applies that omission generically, even on host `codex`, so item 17 returned a Codex assignment pointer with no `name`. The Codex adapter's required mapping had no `task_name` to forward, and the FO skipped the native spawn. Any separate durable-attribution/observer defect exposed by the same run belongs to 8b, not this task.

This is a host-shape mismatch. Claude can realize bare mode with an unnamed blocking `Agent(...)`. Pi can forward an unnamed build artifact into `subagent(...)`/its team substrate and receives a native run handle. Codex has only `spawn_agent(task_name,message,fork_turns)` for fresh dispatch: it requires a name, is async, and returns the addressable handle. Codex therefore has no faithful anonymous blocking realization of `bare_mode`.

## Spike result

The cheapest boundary was exercised against current `main` with the same entity and `host:"codex"`. With `bare_mode:true`, `dispatch build` exited 0 and emitted `prompt` but no `name`; with only `bare_mode:false` changed, it exited 0 and emitted the same prompt plus `name:"spacedock-ensign-codex-approved-gate-worker-spawn-ideation"`. The focused Codex host-shape tests passed.

The live native spawn path is also proven by this ideation dispatch: the named build artifact was forwarded through `spawn_agent` and produced the addressable worker `/root/spacedock_ensign_f02j6dbnd4_ideation`, which is writing this report. The risky mechanism is therefore not hypothetical; the defect is the unsupported bare request being allowed to bypass it.

## Proposed approach

1. At `internal/dispatch/build.go`, after resolving `host` and `bare_mode` but before advisory checks, dispatch-body assembly, or file write, reject `host == "codex" && bare_mode` with exit 2 and an exact diagnostic: `bare_mode is unsupported on host codex; Codex worker.spawn requires a named spawn_agent task`. This is the mechanical correction serving AC-1/AC-2.
2. In `skills/first-officer/references/codex-first-officer-runtime.md`, bind Codex fresh dispatch explicitly to the named shape: omit `--bare-mode` and `--team-name`; a zero-exit build must carry `name` and `prompt`, which map through the existing `CodexMultiAgentV2SpawnInput` to `spawn_agent(..., fork_turns="none")`. This explains the host choice but is not the enforcement.
3. Update dispatch help and the two runtime-support documents to state the same host-shape rule: an unrepresentable host/mode combination fails before artifact creation and is never silently reinterpreted.
4. Run one supported exact-candidate Codex validation. For the positive arm, build the named artifact, invoke the native spawn, observe its addressable task path/handle in the live roster before any report read, wait for completion, and verify the worker-authored report, commit, and workflow state. For the negative arm, build the same named artifact but omit spawn; the empty matching roster check stops validation before reading or crediting any report. Retain the validation record with the entity/review evidence rather than committing a provider-event parser.

No `internal/ensigncycle/shared_keep_moving*`, Codex JSONL evidence, or other standing observer file is changed. If CI still requires durable automated attribution after the runtime correction, 8b (`codex-keep-moving-durable-evidence-attribution-flake`, `8bnkrtq4rw46xkbez5zrbmmj`) owns the redesigned host-neutral oracle and is an explicit dependency rather than hidden scope here.

## Why Codex bare mode is rejected, not reinterpreted

Making Codex `bare_mode=true` still emit `name` is smaller in bytes but wrong in semantics. The documented flag promises no name and a blocking sequential shape. Emitting a name would silently turn it into the ordinary named async Codex shape, make `bare_mode` lie to callers, and diverge from the valid Claude and Pi meanings. A host-scoped refusal is honest: Codex cannot realize that shape, so the invalid combination fails before an artifact can be mistaken for dispatch.

The adapter omission alone is also insufficient: the recorded prompt already required spawn and still selected bare mode. The binary refusal makes the exact bad request impossible; the one-off successful-build/no-spawn roster control independently proves the validation stops when the native tool call is omitted, without shipping an observer dialect.

## Expected surface and semantic budget

- `internal/dispatch/build.go`: 6–10 insertions for the host-scoped guard; command grammar unchanged, but `host=codex + bare_mode=true` changes from exit 0/anonymous artifact to exit 2/no artifact.
- `internal/dispatch/build_codex_host_test.go`: 25–40 insertions for bare rejection/no-file and named-output controls.
- `skills/first-officer/references/codex-first-officer-runtime.md`: 2–5 insertions binding Codex to the named shape.
- `internal/dispatch/dispatch.go` and `internal/dispatch/help_test.go`: 2–8 insertions documenting and pinning the Codex exception in `dispatch build --help`.
- `docs/runtime-support.md` and `docs/site/contributing/adding-a-runtime.md`: 2–4 insertions each documenting that an unrepresentable host shape must be rejected rather than silently reinterpreted.

Expected total: 7 files, about 45–70 insertions and 0–10 deletions. Tolerance: at most 8 files, 95 insertions, and 20 deletions; the extra file is allowed only for an existing dispatch parity fixture that must change to compile. No `internal/ensigncycle` edit is authorized. No command addition, stored-format change, worker receipt, observer schema, transcript dialect, provider parser, controller, daemon, or retry mechanism is authorized. Authority changes only by forbidding anonymous Codex dispatch and requiring the already-native addressable spawn. KD candidate `b60d1c8` is read-only evidence and is outside the edit boundary.

## Documentation diff

In both runtime-support documents, add after the dispatch host-mode guidance:

Before:

`Teach spacedock dispatch build to accept host: "<host>" when the assignment shape differs by host.`

After:

`Teach spacedock dispatch build to accept host: "<host>" when the assignment shape differs by host. If a generic dispatch mode has no faithful native call shape, reject that host/mode combination before artifact creation; do not silently reinterpret it. Codex fresh dispatch is always named, so host=codex rejects bare_mode.`

Also change the existing `dispatch build --help` line from:

`--bare-mode  Emit the bare sequential shape (no name, no team_name, no run_in_background).`

to:

`--bare-mode  Emit the bare sequential shape (no name, no team_name, no run_in_background); unsupported on host=codex.`

No site command-reference change is needed: it does not enumerate dispatch-build flags.

## Test plan

1. Dispatch unit test, low cost: `host=codex, bare_mode=true` exits 2, names the unsupported combination, and creates no dispatch file; changing only to false exits 0 with unchanged prompt and a non-empty name.
2. Existing adapter/contract tests, low cost: `CodexMultiAgentV2SpawnInput` maps the successful build to exactly `task_name`, `message`, and `fork_turns=none`; Claude bare/merged and Pi host tests retain their native shapes.
3. Dispatch help/parity tests, low cost: help names the Codex exception; existing Claude bare, Pi bare/native wrapper, Codex named, JSON/flag input, and parity cases stay green.
4. One-off live Codex validation, high cost and required because runtime authority is the claim: after approved-gate advances, observe the named addressable worker in the runtime roster before any report read, then verify its durable report/commit/state. Run a deliberate named-build/no-spawn control first; absence from the roster must stop the procedure before report inspection or completion credit. Retain the positive and negative observations as validation evidence only.
5. If automated durable attribution is required by the final gate, declare 8b incomplete/blocking and wait for its host-neutral oracle; do not add that mechanism to f02.
6. Repository verification: `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.

## Rejected alternatives

- Prose-only strengthening: rejected because the current shared core and adapter already contain the obligation; the same model prompt skipped it.
- Observer-only tightening: rejected from this scope; without the binary host-mode guard, the runtime still emits the unusable anonymous Codex artifact, and durable automated attribution belongs to 8b.
- Emit `name` in Codex bare mode: rejected because it violates the flag's anonymous/blocking semantics and hides an unsupported host/mode combination.
- Credit a dispatch build, commit, report, wait, status, or merge as spawn: rejected because none creates worker identity or delegates write authority.
- Reuse or extend an existing JSONL/provider-event parser: rejected; reusing a standing dialect is still the forbidden mechanism. The one-off live roster observation proves this candidate, while 8b owns any durable host-neutral replacement.
- Add a receipt/registry/controller: rejected as a second protocol and unnecessary for the host-mode correction.
- Invoke Codex tools from the Go binary: rejected because `spawn_agent` is a host-native model tool, not a subprocess/API available to `spacedock`.

## Stage Report: ideation

- DONE: Reproduce and explain why the supported Codex FO advanced and terminalized approved-gate without an observed worker spawn.
  Reproduced `host=codex + bare_mode=true` returning an anonymous successful artifact; item 17 used that shape, and existing report-as-dispatch credit let items 26–27 pass without identity.
- DONE: Spike the cheapest correction at the actual Codex dispatch boundary, with a missing-spawn control that fails before report reads or completion credit.
  The current named Codex shape emits the required task name and this addressable worker proves the native spawn; the design rejects Codex bare mode mechanically and flips successful-build/no-spawn into a failing control.
- DONE: Define value ACs, exact semantic/file/LOC boundary, test plan, and rejected alternatives without transcript/provider dialect parsing or a second observer protocol.
  AC-1 measures native addressable worker execution; the seven-file, bounded plan reuses the existing adapter and collaboration item and explicitly forbids new dialects, receipts, controllers, and KD edits.

### Summary

The defect is a host-shape mismatch, not missing prose: Codex accepted a bare artifact that omitted the task name its only native spawn requires. The smallest honest correction rejects that unsupported combination before artifact creation, binds Codex to named dispatch, and proves the value by requiring the existing native spawn identity even when build/report/merge evidence otherwise looks complete.

## Stage Report: ideation (cycle 2)

- DONE: Reproduce and explain why the supported Codex FO advanced and terminalized approved-gate without an observed worker spawn.
  Root cause remains the reproduced successful `host=codex + bare_mode=true` artifact with no `name`; durable-attribution concerns from the same incident are explicitly routed to 8b.
- DONE: Spike the cheapest correction at the actual Codex dispatch boundary, with a missing-spawn control that fails before report reads or completion credit.
  The design mechanically rejects Codex bare mode before artifact write and uses a one-off named-build/no-spawn runtime-roster control that halts validation before report inspection.
- DONE: Define value ACs, exact semantic/file/LOC boundary, test plan, and rejected alternatives without transcript/provider dialect parsing or a second observer protocol.
  Revised ACs and the seven-file budget authorize no `internal/ensigncycle` edit; actual spawn is retained live evidence, and any durable automated oracle is an explicit dependency on 8b.

### Summary

Cycle 2 removes every standing event-parser and keep-moving observer change from f02. The task now owns only the mechanical Codex host-mode correction, its adapter/help/docs coverage, and one-off live roster-plus-durable-state validation; 8b remains the sole owner of any redesigned automated attribution oracle.

## Stage Report: implementation

- DONE: Reject `host=codex + bare_mode=true` before dispatch artifact creation while preserving named Codex output and Claude/Pi host shapes.
  Commit `959d6d850` adds the pre-assembly exit-2 guard; `TestBuildCodexHostRejectsBareModeBeforeArtifactCreation` fails if bare Codex writes its deterministic artifact or named Codex loses `name`/`prompt`, while the host-parity selection fails on Claude bare/merged or Pi native-shape drift.
- DONE: Update the Codex adapter, dispatch help, and runtime host-shape docs within the approved seven-file, 45-70 insertion design.
  Commit `959d6d850` changes exactly 7 files at +55/-2; help pins the exception, both runtime guides require rejecting unrepresentable modes, and contractlint fails if the adapter diverges from `task_name,message,fork_turns="none"`.
- DONE: Run focused/parity/full/race checks, then retain one positive live roster-to-durable-worker journey and one named-build/no-spawn stop control.
  Focused dispatch/adapter/help tests, host parity, `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and `git diff --check` pass; the named-build/no-spawn roster was empty before report reads, while spawned task `/root/spacedock_ensign_codex_approved_gate_worker_spawn_implementation/spacedock_ensign_approved_gate_implementation` appeared running before reads and authored clean temp commit `852a12a50fb3914ef8d253c8c5719efd65ca1f34`, report, proof, and `status=done`.

### Summary

Codex bare mode now fails honestly at the dispatch boundary with the approved diagnostic, while normal named Codex, Claude, and Pi semantics remain intact. The implementation stays inside the seven-file mechanical surface and proves the value through a one-off native roster observation followed by worker-authored durable state, without adding an observer or provider dialect.

## Review-finding disposition

### Validation evidence defect (closed)

- Reviewer observation: the first validator-positive fixture used the exact candidate for named build/spawn but inherited the older repository launcher for its terminal status write while calling that launcher exact.
- Worker proposal: evidence defect; validation-fixture owned; Material to AC-1 proof only; rerun one fresh positive arm without candidate mutation.
- Released user and normal workflow: exact-candidate live validation is the required supported Codex validation workflow.
- Observable harm: commit `7441ce12` could not prove exact-candidate attribution for the terminal state write; no candidate product harm was observed.
- Authority: `value-ac[AC-1]` requires one exact-candidate live roster-to-worker-authored-durable-state proof.
- Trigger evidence: the worker report named `/Users/clkao/git/spacedock-research/spacedock-v1/spacedock` (`0.26.0+dev`) while the build/spawn binary was the `959d6d850` build (`0.27.0-pre2+dev`).
- First Officer authorization: rerun only the positive arm with the absolute candidate launcher pinned and preflighted; no negative rerun, candidate mutation, observer/parser, fallback, or third attempt.
- Resolution: task `.../spacedock_ensign_validator_positive_exact_implementation` was addressable before report read, preflighted the exact launcher, and authored clean commit `c7ab171f`; the candidate stayed unchanged.

## Stage Report: validation

- DONE: Reproduce AC-1 through AC-4 at exact candidate 959d6d850, including fail-before-artifact Codex bare mode and unchanged named/Claude/Pi shapes.
  AC-1: exact-candidate task `/root/spacedock_ensign_codex_approved_gate_worker_spawn_validation/spacedock_ensign_validator_positive_exact_implementation` was live/addressable before report read, then authored clean commit `c7ab171f`, proof, report, and `status=done`.
  AC-2: JSON, CLI-flag, and CODEX_THREAD_ID host-resolution arms each exited 2 with the approved diagnostic and no artifact; the named build retained name/prompt.
  AC-3: `TestCodexSpawnSignatureBindsToolArgs` and adapter tests bind exactly `task_name,message,fork_turns=none`; changed-path/diff review found no observer, event parser, receipt, schema, controller, daemon, retry, or `internal/ensigncycle` edit.
  AC-4: focused Claude bare, Pi wrapper, Codex host, help, dispatch parity, full, and race suites passed without supported-host shape drift.
- DONE: Independently run the one-off named-build/no-spawn stop control and positive addressable-worker-to-durable-report journey without a standing event parser.
  Control `spacedock-ensign-validator-no-spawn-implementation` had no matching roster worker and stopped before report inspection; the authorized positive rerun produced exact-launcher proof at commit `c7ab171f`.
- DONE: Verify the exact seven-file +55/-2 surface, focused/parity/full/race/format checks, docs/adapter coherence, and absence of observer/schema/controller expansion.
  Candidate `959d6d850` remained clean at 7 files and +55/-2; `go test ./...`, `go test ./... -race`, focused suites, empty `gofmt -d`, and `git diff --check` passed.
- DONE: Perform the semantic adversarial pass and detached audit.
  The matrix covered JSON/flag/env host inputs, bare/named Codex, Claude/Pi bare, fail cleanup, exact diagnostic, name/prompt cardinality, adapter argument identity, and no new scaling path; detached checkout `/tmp/spacedock-f02-detached-audit.gBQk0m` was clean and green at `959d6d850`.
- DONE: Classify and close the validation evidence finding before rerun.
  The First Officer authorized the fixture-owned AC-1 evidence correction; no candidate mutation or product finding resulted, and the exact-candidate rerun closed the proof gap.

### Summary

Validation recommends **PASSED**: AC-1 through AC-4 reproduce on exact candidate `959d6d850`, with the no-spawn arm stopping at the live roster boundary and the positive arm proving addressable worker authority through durable commit/state. The only finding was a validator-fixture attribution error; it was classified, authorized, and closed by one exact-launcher positive rerun without changing the candidate. No material finding or deferred risk remains.
