---
title: Select the actionable Codex default-headless task
status: ideation
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
started: 2026-08-11T05:24:54Z
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

## Acceptance criteria and test plan

**AC-1 (VALUE): The exact Codex default-headless journey records one native
implementation `spawn_agent`, its completion before validation, and one DONE
implementation report.** The existing `assertImplementationWorkerLifecycle`
oracle measures this end value. Removing the spawn or moving validation before
completion restores `implementation-worker-not-dispatched`.

**AC-2: The journey ends at one committed, clean, open validation gate, with no
decision, consumed authority, terminal fields, or successor dispatch.** The
existing `assertGateHeld` and `assertRecordedGateHoldLog` checks test the
resulting state and command order. A gate decision, post-prepare status change,
or successor build fails them.

**AC-3: The correction remains one Codex adapter change within the declared
product cap and preserves all other runtime semantics.** `git diff --stat`,
`git diff --check`, and review against the expected surface test this boundary.
AC-3 serves AC-1; it is not satisfied if AC-1 is not green.

**AC-4: The Codex target progresses from bound XFAIL to bound XPASS, then to
unbound PASS after only its two binding rows change.** Retain the exact bound
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
