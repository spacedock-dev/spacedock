---
title: Require Codex worker spawn after approved gate advance
status: ideation
source: "KD validation cycle 4, Codex keep-moving-posture at b60d1c8"
started: 2026-07-30T23:39:35Z
completed:
verdict:
score: 0.95
worktree:
issue:
milestone: 0.27.0
id: f02j6dbnd4jakwczahv1tg2h
gates:
    version: 1
    current:
        gate: gate:f02j6dbnd4jakwczahv1tg2h:backlog
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
---

A supported headless Codex keep-moving journey advanced `approved-gate`, built its implementation dispatch artifact, never invoked `spawn_agent` / `worker.spawn`, then read a report and terminalized the task. This violates the shipped dispatch boundary and silently bypasses worker execution and write-scope authority.

The exact evidence is retained at `/tmp/spacedock-validation-cycle4-b60d1c8.qPdIXO/codex/codex-shared-scenarios/keep-moving-posture/codex-exec.jsonl`: item 17 advances and builds; items 26–27 read and finalize; no worker identity or spawn exists. The live harness reported: `the FO advanced the approved entity "approved-gate" but did not dispatch its next stage`.

Ideation must find the smallest fix that makes the supported Codex journey actually spawn and await a worker after gate approval. It must not repair the observer, parse transcript/provider dialects, or let a dispatch commit or generated artifact substitute for an observed worker spawn.

## Acceptance criteria

**AC-1 — Approved work executes under worker authority.** In the supported Codex keep-moving journey, every ready implementation task, including one made ready by gate approval, has exactly one completed `spawn_agent` with a non-empty receiver handle after its build and before any report read, completion credit, or terminalization. Test with the existing live keep-moving journey and its positive fixture; either moving the spawn after the report read or removing its receiver handle must fail.

**AC-2 — Missing spawn remains a hard failure.** Codex rejects `bare_mode=true` before writing a dispatch artifact, while a successful named Codex build that is not forwarded to `spawn_agent` earns no report/completion credit and cannot validate as terminal. Test both controls: a dispatch-build unit test expects exit 2 and no artifact for Codex bare mode; the existing keep-moving negative fixture retains successful build, durable report, and merge output but omits spawn and must fail on missing dispatch.

**AC-3 — The correction stays runtime- and observer-coherent.** The implementation uses the existing Codex build-to-`CodexMultiAgentV2SpawnInput`/`spawn_agent` binding and the existing completed-collaboration-item observation; it adds no transcript grammar, provider-event parser, receipt, stored schema, retry loop, or second observer protocol. Test with the existing Codex adapter/contractlint suite and a diff review that rejects any new event or persistence representation; this mechanism serves AC-1.

**AC-4 — Supported host semantics do not regress.** Claude bare dispatch remains anonymous, blocking, and sequential; Pi keeps forwarding its assignment through its native substrate; named Codex dispatch still emits the unchanged `name` and `prompt` fields. Test the existing Claude merged/bare, Pi host, Codex host, and dispatch parity suites.

## Problem and root cause

The failure is not a missing sentence. The shared core and Codex adapter already say a successful build is not a dispatch and require `spawn_agent`. In the recorded run, however, the Codex FO invoked:

`dispatch build ... --stage implementation ... --bare-mode`

`bare_mode` is declared as the anonymous sequential shape: no `name`, no `team_name`, and no `run_in_background`. `internal/dispatch/build.go` applies that omission generically, even on host `codex`, so item 17 returned a Codex assignment pointer with no `name`. The Codex adapter's required mapping has no `task_name` to forward; the FO skipped the native spawn and later treated a durable report as retroactive dispatch evidence. The existing keep-moving observer then compounded the authority breach by crediting build + report/merge as dispatch.

This is a host-shape mismatch. Claude can realize bare mode with an unnamed blocking `Agent(...)`. Pi can forward an unnamed build artifact into `subagent(...)`/its team substrate and receives a native run handle. Codex has only `spawn_agent(task_name,message,fork_turns)` for fresh dispatch: it requires a name, is async, and returns the addressable handle. Codex therefore has no faithful anonymous blocking realization of `bare_mode`.

## Spike result

The cheapest boundary was exercised against current `main` with the same entity and `host:"codex"`. With `bare_mode:true`, `dispatch build` exited 0 and emitted `prompt` but no `name`; with only `bare_mode:false` changed, it exited 0 and emitted the same prompt plus `name:"spacedock-ensign-codex-approved-gate-worker-spawn-ideation"`. The focused Codex host-shape tests passed.

The live native spawn path is also proven by this ideation dispatch: the named build artifact was forwarded through `spawn_agent` and produced the addressable worker `/root/spacedock_ensign_f02j6dbnd4_ideation`, which is writing this report. The risky mechanism is therefore not hypothetical; the defect is the unsupported bare request being allowed to bypass it.

## Proposed approach

1. At `internal/dispatch/build.go`, after resolving `host` and `bare_mode` but before advisory checks, dispatch-body assembly, or file write, reject `host == "codex" && bare_mode` with exit 2 and an exact diagnostic: `bare_mode is unsupported on host codex; Codex worker.spawn requires a named spawn_agent task`. This is the mechanical correction serving AC-1/AC-2.
2. In `skills/first-officer/references/codex-first-officer-runtime.md`, bind Codex fresh dispatch explicitly to the named shape: omit `--bare-mode` and `--team-name`; a zero-exit build must carry `name` and `prompt`, which map through the existing `CodexMultiAgentV2SpawnInput` to `spawn_agent(..., fork_turns="none")`. This explains the host choice but is not the enforcement.
3. In the existing keep-moving behavioral assertion, stop granting Codex dispatch credit from build + durable report/status. Reuse the already parsed completed `spawn_agent` collaboration item and its `receiver_thread_ids`; require that identity for the approved entity and independent ready entities. Flip the current no-`spawn_agent` “standing-loop dialect” fixture from positive to the missing-spawn negative. No new provider grammar is introduced.
4. Rerun the supported Codex keep-moving live journey. Its success proof is the actual native spawn plus durable worker result; the no-spawn control must fail even when build, report, and merge outputs are present.

## Why Codex bare mode is rejected, not reinterpreted

Making Codex `bare_mode=true` still emit `name` is smaller in bytes but wrong in semantics. The documented flag promises no name and a blocking sequential shape. Emitting a name would silently turn it into the ordinary named async Codex shape, make `bare_mode` lie to callers, and diverge from the valid Claude and Pi meanings. A host-scoped refusal is honest: Codex cannot realize that shape, so the invalid combination fails before an artifact can be mistaken for dispatch.

The adapter omission alone is also insufficient: the recorded prompt already required spawn and still selected bare mode. The binary refusal makes the exact bad request impossible; the successful-build/no-spawn behavior control independently proves that a caller cannot regain completion credit merely by omitting the native tool call.

## Expected surface and semantic budget

- `internal/dispatch/build.go`: 6–10 insertions for the host-scoped guard; command grammar unchanged, but `host=codex + bare_mode=true` changes from exit 0/anonymous artifact to exit 2/no artifact.
- `internal/dispatch/build_codex_host_test.go`: 25–40 insertions for bare rejection/no-file and named-output controls.
- `skills/first-officer/references/codex-first-officer-runtime.md`: 2–5 insertions binding Codex to the named shape.
- `internal/ensigncycle/shared_keep_moving_test.go`: 8–20 insertions and 8–20 deletions to require the existing completed spawn/receiver identity and remove report-as-dispatch credit for this journey.
- `internal/ensigncycle/shared_keep_moving_negative_test.go`: 15–30 changed lines to make successful-build/no-spawn a failing control.
- `docs/runtime-support.md` and `docs/site/contributing/adding-a-runtime.md`: 2–4 insertions each documenting that an unrepresentable host shape must be rejected rather than silently reinterpreted.

Expected total: 7 files, about 60–90 insertions and 15–35 deletions. Tolerance: at most 9 files, 120 insertions, and 60 deletions. No command addition, stored-format change, worker receipt, observer schema, transcript dialect, provider parser, controller, daemon, or retry mechanism is authorized. Authority changes only by forbidding anonymous Codex dispatch and requiring the already-native addressable spawn. KD candidate `b60d1c8` is read-only evidence and is outside the edit boundary.

## Documentation diff

In both runtime-support documents, add after the dispatch host-mode guidance:

Before:

`Teach spacedock dispatch build to accept host: "<host>" when the assignment shape differs by host.`

After:

`Teach spacedock dispatch build to accept host: "<host>" when the assignment shape differs by host. If a generic dispatch mode has no faithful native call shape, reject that host/mode combination before artifact creation; do not silently reinterpret it.`

No command-reference change is needed: the public reference does not document `--bare-mode`; the changed surface is the runtime-integration contract.

## Test plan

1. Dispatch unit test, low cost: `host=codex, bare_mode=true` exits 2, names the unsupported combination, and creates no dispatch file; changing only to false exits 0 with unchanged prompt and a non-empty name.
2. Existing adapter/contract tests, low cost: `CodexMultiAgentV2SpawnInput` maps the successful build to exactly `task_name`, `message`, and `fork_turns=none`; Claude bare/merged and Pi host tests retain their native shapes.
3. Keep-moving fixture, medium cost: a completed `spawn_agent` with a receiver handle before durable completion passes. The same successful build/report/merge stream with no spawn fails specifically on approved dispatch; a spawn without a receiver handle and a spawn placed after completion also fail.
4. Live Codex journey, high cost and required because runtime authority is the claim: gate approval must lead to a native addressable spawn before the report read and merge; verify the worker's durable report/commit and final state. Use the existing harness and structured collaboration item already consumed by it, not a new event dialect.
5. Repository verification: `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.

## Rejected alternatives

- Prose-only strengthening: rejected because the current shared core and adapter already contain the obligation; the same model prompt skipped it.
- Observer-only tightening: necessary as proof but not a product correction; without the binary host-mode guard, the runtime still emits the unusable anonymous Codex artifact.
- Emit `name` in Codex bare mode: rejected because it violates the flag's anonymous/blocking semantics and hides an unsupported host/mode combination.
- Credit a dispatch build, commit, report, wait, status, or merge as spawn: rejected because none creates worker identity or delegates write authority.
- Parse the diagnostic JSONL or add a receipt/registry/controller: rejected as a second protocol and unnecessary; the existing Codex spawn binding and completed collaboration item already expose the required identity.
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
