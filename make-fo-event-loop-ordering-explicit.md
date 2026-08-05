---
title: Make FO event-loop ordering and idle wait explicit
status: validation
source: "Captain follow-up after the 2026-08-03 durable-decisions execution-gap diagnosis."
started: 2026-08-03T16:00:51Z
completed:
verdict:
score: 0.98
worktree: .worktrees/spacedock-ensign-make-fo-event-loop-ordering-explicit
issue:
sprint: durable-decisions
group: fo-contract
milestone: 0.27.0
id: ej9kwkvw94w6rh6n5ek7qrbf
gates:
    version: 1
    records:
        - id: gate:ej9kwkvw94w6rh6n5ek7qrbf:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:ej9kwkvw94w6rh6n5ek7qrbf-backlog-1
              briefing:
                id: briefing:ej9kwkvw94w6rh6n5ek7qrbf:backlog:attempt-1:revision-1
                digest: sha256:4dbaa852a23e16b0c2c732a3fac6956ef5f95db8c603f428b1636a5983d3064b
                request-digest: sha256:2716adf5a2f9014015e76146c012c4feb29a71001469f747748059142ecf26cb
                room-ref: ./make-fo-event-loop-ordering-explicit/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ej9kwkvw94w6rh6n5ek7qrbf:backlog:1
                briefing: briefing:ej9kwkvw94w6rh6n5ek7qrbf:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T15:55:01.310359Z"
                decision: approve
                reason: Captain conn approves ideation. SO concurs with parallel ideation for EJ and G3; preserve EJ-before-G3 implementation landing because both define the shared FO dispatch/merge contract.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:ej9kwkvw94w6rh6n5ek7qrbf:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:ej9kwkvw94w6rh6n5ek7qrbf-ideation-1
              briefing:
                id: briefing:ej9kwkvw94w6rh6n5ek7qrbf:ideation:attempt-1:revision-1
                digest: sha256:271f27759b8b21a8d58843bccd6b95f56e4de8663df3c539cec422d3631917c4
                request-digest: sha256:bdfe10a7c6aeadf9d18ef3aa6e1e98d7dabe404d33e874c1d240bb46d7b3264c
                room-ref: ./make-fo-event-loop-ordering-explicit/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ej9kwkvw94w6rh6n5ek7qrbf:ideation:1
                briefing: briefing:ej9kwkvw94w6rh6n5ek7qrbf:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-03T16:11:57.592946Z"
                decision: approve
                reason: Captain conn authorizes approval when SO concurs. The ideation report supplies the false-stop trace, ordered mod/PR and gate handling, idle/reconcile retry, unresolved-worker wait matrix, and keep-moving evidence. SO concurs with EJ-before-G3 implementation landing after SK validation.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:ej9kwkvw94w6rh6n5ek7qrbf:validation
          stage: validation
          attempts:
            - id: gate-attempt:ej9kwkvw94w6rh6n5ek7qrbf-validation-1
              briefing:
                id: briefing:ej9kwkvw94w6rh6n5ek7qrbf:validation:attempt-1:revision-1
                digest: sha256:3596fec21319ec0d39fc21bdf1b95eb824ac4e2420eebe2a32a6a55601341469
                request-digest: sha256:fbd48aae4643c6411a54eac46ce31f59d5aabd650bedf3de739a748bba1c1682
                room-ref: ./make-fo-event-loop-ordering-explicit/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ej9kwkvw94w6rh6n5ek7qrbf:validation:1
                briefing: briefing:ej9kwkvw94w6rh6n5ek7qrbf:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-05T16:16:17.040889Z"
                decision: revise
                reason: Reconcile candidate 5fa79d480 with current origin/main in skills/first-officer/references/fo-dispatch-core.md. Preserve EJ event-loop and wait semantics. Produce a new exact head, then run fresh validation and gate evidence.
            - id: gate-attempt:ej9kwkvw94w6rh6n5ek7qrbf-validation-2
              briefing:
                id: briefing:ej9kwkvw94w6rh6n5ek7qrbf:validation:attempt-2:revision-1
                digest: sha256:a08a4971ba483ac6128e650baa268bb2cc3dd3bcae035dd80461ef5563fbb7f6
                request-digest: sha256:e5e0e15301749b7356232f2d2188be9beb90175f3c8e56a0a98bbc87ef5fbf34
                room-ref: ./make-fo-event-loop-ordering-explicit/review/validation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:ej9kwkvw94w6rh6n5ek7qrbf:validation:2
                briefing: briefing:ej9kwkvw94w6rh6n5ek7qrbf:validation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-05T17:00:48.821361Z"
                decision: approve
                reason: Captain accepted exact candidate c7f3bc0d8 after clean current-main reconciliation and fresh cycle-4 validation.
              application:
                target-stage: done
                state: pending
mod-block: merge:pr-merge
pr: pr-merge:628
---

Make the First Officer event loop mechanically explicit so a dispatch-only empty result cannot hide merge recovery, ready gates, or a required idle/reconcile pass. The task preserves the existing state and runtime boundaries; it makes the ordering observable and testable.

## Problem

The current contract has the required branches, but an FO can run a filtered `status --next` query, see `dispatchable: []`, and report idle without processing `mod-block`/PR recovery, the separate `ready_gates` surface, the one-shot idle hook and reconcile retry, or the unresolved-worker wait rule. This loses the distinction between no dispatchable worker, a pending merge action, a gate awaiting Captain input, and a genuinely idle session.

## Proposed approach

Specify one event-loop trace and bind every runtime adapter to it. Each iteration drains worker messages when supported, processes every `mod-block` and PR action before querying dispatchables, handles `ready_gates` as gate or merge actions rather than worker dispatches, and then runs `status --next`. After the first empty result, run the idle hook exactly once, reconcile the roster, and query `status --next` once more. Install async wait only when the worker set is unresolved and there is no dispatchable entity, presentable gate, mod action, or other state work; a completed, errored, or absent worker set never qualifies. Preserve keep-moving for unrelated ready tasks and do not add a new state field or scheduler command.

## Out of scope

- Changing the commissioned workflow definition or stage taxonomy.
- Adding a new status command, daemon, watchdog, or generic resolver worker.
- Auto-resolving PR conflicts, force-pushing, or changing merge/approval authority.
- Rewriting task-specific findings or changing the Codex, Claude, or Pi transport APIs.

## Acceptance criteria

**AC-1 (VALUE) — A ready action cannot be hidden by a dispatch-only empty result.**
Verified by: a fixture-backed event-loop trace with one `mod-block`/PR row, one `awaiting-captain` or `approved-awaiting-merge` gate, and one genuinely dispatchable task; the recorded order shows mod/PR and gate handling before `status --next`, and the unrelated ready task still dispatches.

**AC-2 — Empty dispatch results trigger one idle/reconcile retry and a truthful stop.**
Verified by: a command-log fixture that returns an empty first `status --next`, releases a task from the idle hook or reconcile, and asserts one idle hook, one reconcile, one second `status --next`, and dispatch of the newly unblocked task; the unchanged-empty variant asserts the explicit no-dispatchable stop reason.

**AC-3 — Async wait is reserved for unresolved workers.**
Verified by: Codex runtime fixtures covering active unresolved, completed, errored, and absent worker sets; only the active unresolved case with no other work may emit `wait_agent(timeout_ms: 300000)`, and no other case emits a wait.

**AC-4 — Gate readiness is not conflated with worker dispatchability.**
Verified by: boot/status fixtures that expose `awaiting-captain`, `approved-awaiting-advance`, and `approved-awaiting-merge` rows, and assert presentation, consume/advance, and merge routing respectively without a worker spawn or a false idle report.

**AC-5 — One blocked task does not park unrelated work.**
Verified by: a multi-entity fixture with one conflict/mod-blocked task and two independent ready tasks; the trace shows the blocked task held with its evidence while both independent tasks reach their declared dispatch or gate action.

## Expected surface and semantic boundaries

Expected change is limited to the First Officer dispatch/event-loop contract and its behavior fixtures: `skills/first-officer/references/first-officer-shared-core.md`, `skills/first-officer/references/fo-dispatch-core.md`, the host runtime adapter notes needed to bind the wait predicate, and focused event-loop/recorded-workflow tests. Estimate 4–8 files and +120/-20 lines, with a 2x tolerance. Allowed semantic change: ordering and stop-reason observability only. Stored state, command grammar, gate authority, worker transport, and stage behavior remain unchanged.

## Test plan

Start with the smallest falsifying fixture: replay the current false-stop trace and assert that a filtered `status --next=[]` cannot terminate the iteration while a mod-block, ready gate, or unresolved worker remains. Add deterministic command-log tests for each branch and one runtime-backed Codex wait test. Run the focused status/dispatch/ensigncycle suites, then `go test ./...`, `go test ./... -race`, formatting, and the required host live lanes when contract/runtime adapter files change.

## Stage-specific test gates

- Ideation must replay the falsifying fixture and record the exact ordered trace before selecting a landing shape.
- Implementation must keep the trace assertions independent of contract prose and preserve the existing state and command bytes.
- Validation must cover the empty-first-query retry, gate/merge separation, unresolved-worker wait predicate, and keep-moving matrix, then run full, race, formatting, and relevant live evidence.

## Stage Report: ideation

### False-stop replay and ordered trace

The current checkout provides a direct falsifier for treating an empty dispatch query as idle (run 2026-08-04):

```text
status --where 'mod-block !=' --json --fields id,slug,mod-block
  26n/headless-recorded-gate-stop-stage-coherence  merge:pr-merge
  32/boot-identify-multi-workflow-llm-retry-friction merge:pr-merge
  93/document-dispatch-entity-path-base              merge:pr-merge
status --next --json --fields id,slug
  {"dispatchable":[]}
status --boot --identify --json
  {"dispatchable":[],"ready_gates":[],"pr_state":{"status":"local",...},"team_state":{"present":true,...}}
```

The first query exposes three pending PR/merge actions; the second query alone therefore cannot be a stop decision. The boot-identify status and recorded fixtures provide the complementary gate case: `internal/status/boot_identify_test.go` passes with `ready_gates` rows `mf` (`awaiting-captain`), `r4` (`awaiting-captain`), `wd` (`withdrawn-awaiting-prepare`), and `2n` (`approved-awaiting-merge`) alongside dispatchable `dispatch-me`. The terminal-gate test keeps `approved-awaiting-merge` in `ready_gates` through `gate consume` until merge delivery, proving that consume is not worker dispatch. The Codex wait fixture `TestCodexWaitAgentSteeringEvidence` also passes: after spawn and active-scope exhaustion it records `wait_agent(timeout_ms:300000)`, retries after captain steering, then verifies the durable report after `final_status`.

The implementation-gate trace must be explicit about both the existing baseline and the required decision order. For each iteration:

1. If a roster is addressable, reconcile it; then drain addressable inbound worker messages. If the runtime has no roster/message primitive, record that adapter boundary rather than inferring a clean roster.
2. Query `status --where "mod-block !="`; process every blocking mod/PR before dispatch. A merged PR records its sentinel and runs merge guard; an open PR remains held; a closed PR asks the Captain; an unavailable check is skipped without clearing the block. No new worker is dispatched by this branch.
3. Read the boot/status `ready_gates` projection and route each gate as state work: present `awaiting-captain`, prepare a successor for `withdrawn-awaiting-prepare`, consume/advance `approved-awaiting-advance` (with successor and observed spawn), and deliver/merge `approved-awaiting-merge` (never spawn). Any unresolved gate prevents an idle claim.
4. Run `status --next`; dispatch every returned row, requiring worker-spawn evidence (a successful dispatch build is not spawn evidence). Keep moving independent ready rows when one row is blocked.
5. Only after the first empty `status --next`, run the idle hook exactly once, reconcile when the adapter supports it, and run one second `status --next`. Dispatch newly released work; if the second result is still empty, report an explicit `no-dispatchable` stop only when there is no mod/PR action, gate action, or other state work.
6. In Codex, emit `wait_agent(timeout_ms:300000)` only when an active unresolved worker remains and the post-retry scope has no dispatchable, gate, mod, or state work. A completed, errored, or absent worker never qualifies. A timeout retries the same epoch; a completion notification starts durable report/state verification and another iteration. Pi has no `wait_agent`; Claude's reconcile and idle behavior follows its bound team identity.

This order makes the false-stop trace observable: the current status-only path is `mod-block? -> status --next=[] -> idle/stop`, while the required path is `message drain -> mod/PR -> ready_gates -> status --next -> (idle once + reconcile + status --next retry) -> explicit stop or guarded wait`.

### Acceptance-criteria and boundary reconciliation

The ACs and existing boundaries are consistent, with the following implementation fixtures required:

| AC | Fixture assertion | Existing boundary preserved |
| --- | --- | --- |
| AC-1 | Mixed mod/PR + gate + dispatchable command log proves mod and gate actions precede `status --next`, while the independent row dispatches. | `status` is read-only; PR merge guard and gate lifecycle own mutation. |
| AC-2 | First empty result emits one idle hook, one reconcile, one retry; released work dispatches; unchanged empty reports `no-dispatchable`. | `hooks.run` remains caller-owned and exact-once; no new scheduler command. |
| AC-3 | Active-unresolved/completed/errored/absent matrix permits wait only for active-unresolved with an otherwise empty scope. | Codex uses `list_agents`/`wait_agent`; completion attribution is mailbox/task-path plus durable state, never a handle alone. |
| AC-4 | `awaiting-captain`, `approved-awaiting-advance`, and `approved-awaiting-merge` route to presentation, consume/advance, and merge respectively, with no false worker spawn/idle. | `ready_gates` and gate bytes remain authoritative; terminal consume remains pending until merge delivery. |
| AC-5 | One blocked row does not prevent two independent ready rows from reaching dispatch/gate actions. | Keep-moving semantics and worker write scopes remain unchanged. |

Focused status, gate-lifecycle, and Codex wait tests already pass. Add a host-neutral command-log/event-loop fixture first, then adapter-specific assertions; do not invent a shipped `dispatch next-action` command. Run the focused suites, `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`; run host live lanes only if runtime adapter notes change.

### Unresolved risk and recommendation

The separate unhandled-completion-ledger risk remains: a final-status/mailbox completion can arrive while the FO is handling another request, and no durable completion feed currently records that it was handled. Ordering the available drain and reconcile operations reduces the window but cannot recover a lost completion or make a reconcile without session identity authoritative; the live `dispatch reconcile --workflow-dir docs/dev` probe reported git/filesystem drift only because no team identity was resolved. Keep that ledger/reclaim work separate, and make wait fail-safe when roster identity is absent.

Recommendation for implementation-gate entry: **APPROVE ideation**. Land a contract-first numbered event-loop trace plus a small host-neutral command-log reducer/fixture and the Codex wait matrix. Make the mixed falsifying fixture red against the current `status --next=[]` stop before updating the contract and adapter notes. Preserve all current state fields, command bytes, transport APIs, gate authority, and stage behavior; the intended semantic change is ordering and stop-reason observability only.

### Ideation completion checklist

- DONE — replayed the false-stop trace and recorded the ordered mod/PR, gate, dispatch, idle/reconcile, and wait decisions.
- DONE — reconciled ACs, test plan, and Codex wait predicate against existing state/transport boundaries and named unresolved risks.
- DONE — committed the ideation report with falsifying evidence and the implementation-gate recommendation.

## Stage Report: implementation

- DONE: Land the contract-first ordered event-loop trace and a falsifying host-neutral command-log fixture proving mod/PR and ready-gate actions precede status --next while unrelated ready work still dispatches.
  Commit `e2d947287` adds the ordered contract and `TestFirstOfficerEventLoopCommandLog`; moving mod/gate events after `next`, dropping either independent dispatch, or accepting the legacy false-stop trace makes it fail.
- DONE: Implement the one-shot idle/reconcile/retry and unresolved-worker wait predicate across the declared host boundaries without adding state fields, command grammar, gate authority, or a shipped scheduler command.
  The shared, Claude, Codex, and Pi bindings now agree; retry action reordering/repetition and active/completed/errored/absent wait-matrix mutations fail focused tests, while Codex, Claude shallow-boot, and Pi front-door live retries passed.
- DONE: Pass focused event-loop/gate/wait tests, go test ./..., go test ./... -race, gofmt for Go changes, and the required live lanes for any runtime-adapter contract changes.
  Focused suites, contractlint, gofmt, both full suites with immutable `SPACEDOCK_STATE_ROOT` commit `73f41e2a2232ebb561710bce568641ec976d5f3d`, Codex isolated live, Claude isolated shallow-boot, and Pi front-door passed; mutable live state still fails the pilot manifest after eight concurrent moves under `_archive`, and the broad Claude run had four transient 60s no-progress exits before the isolated pass.

### Summary

Committed a seven-file, +200/-10 implementation that makes the FO loop drain and route mod/PR and gate work before dispatch projection, retries an empty projection once, and chooses explicit stop or guarded unresolved-worker wait. The fixture is contract-prose-independent and preserves stored state, command grammar, gate authority, transport APIs, and the unshipped scheduler boundary; immutable verification is green and the distinct mutable-state/live-host drift is recorded above without false-greening it.

## Review-finding disposition

- Finding: AC-3's purported Codex runtime matrix is disconnected from the shipped adapter. `TestCodexWaitPredicate` calls only a test-local `shouldWaitForWorker` formula and never loads or exercises the Codex runtime contract or observes a `wait_agent` call.
- Released user and normal workflow: a Codex First Officer choosing whether to monitor an active unresolved worker, or to stop for a completed, errored, or absent worker, after the one-shot retry and with no other work.
- Observable harm: a wrong shipped predicate can stop while a worker is unresolved or wait on terminal work while the suite remains green.
- Authority: `contract[docs/dev/README.md#validation]` requires reproducing each AC's cited evidence and rejecting evidence that cannot establish the delivered behavior; this leaves AC-3 unproved.
- Trigger evidence: in throwaway checkout `/tmp/spacedock-ej-audit.6dMqaV`, inverting the candidate adapter to wait on completed workers and reject active unresolved workers still left both focused tests green.
- Proposal: Material evidence defect owned by EJ validation; preserve candidate bytes and correct only the proof. Replace the local predicate oracle with an exact-head, real-adapter-driven Codex matrix that observes exact `wait_agent(timeout_ms: 300000)` emission for active unresolved and absence for completed, errored, and absent workers; do not add another scheduler/controller.
- Recommendation: REJECTED until AC-3 has valid adapter-bound evidence. Mutable pilot-manifest drift and the Pi account failure below are separate and are not candidate findings.

### Feedback Cycles

- Cycle 1: REJECTED — validation adapter-mutation audit; surface 7 files/210 LOC vs estimate 140 (150%); AC unchanged
- Cycle 2: REJECTED — validation rollout-path audit; surface 8 files/250 LOC vs estimate 140 (179%); AC unchanged
- Cycle 3: REJECTED — Captain current-main reconciliation; surface 9 files/253 LOC vs estimate 140 (181%); AC unchanged

## Stage Report: validation

- DONE: Reproduce the mixed mod/PR, ready-gate, and independent-dispatch command-log behavior and verify AC-1, AC-4, and AC-5 with exact order and spawn evidence.
  AC-1: `mixed` is reconcile, drain, mod-hold, three gate routes, then next and independent-a/b dispatch; AC-4: gate-present/advance/merge precede next and only independent-a/b spawn; AC-5: both independent tasks dispatch while blocked-pr remains held. Reordering or an extra gate spawn fails the fixture oracle.
- FAILED: Adversarially verify the one-shot idle/reconcile/retry path and Codex unresolved-worker wait matrix for AC-2 and AC-3 across active, completed, errored, and absent workers.
  AC-2 passes exact retry and unchanged-empty traces: one idle, one reconcile, one second next, then released dispatch or `no-dispatchable`. AC-3 fails its proof boundary because an inverted shipped adapter still passes the test-local matrix.
- DONE: Confirm no stored-state, command-grammar, gate-authority, transport, or shipped-scheduler drift; run focused, full, race, formatting, and required exact-head live evidence with AC-1 through AC-5 citations.
  Candidate `e2d947287` is exactly seven files at +200/-10; no state/CLI implementation changed, `git diff --check` and `gofmt -d` are clean, and the scheduler remains explicitly prose-only. Immutable `SPACEDOCK_STATE_ROOT` commit `73f41e2a` passes full and race; Codex and Claude exact-head shallow-boot pass.
- DONE: Separate mutable fixture and live-host failures from candidate behavior.
  Mutable `TestV1PilotManifestReadsAndValidates` alone misses eight paths moved under `_archive`; Pi reaches the exact-head front door but OpenRouter rejects before behavior with HTTP 402 credit/max-token limits. Neither touches AC-1 through AC-5 candidate semantics.

### Summary

Validation recommends REJECTED for one Material evidence defect: AC-3's test duplicates the intended predicate locally and cannot fail when the shipped Codex adapter is wrong. AC-1, AC-2, AC-4, and AC-5 have reproduced exact trace evidence, all immutable offline suites are green, and the narrow correction is an adapter-driven Codex wait matrix rather than a new controller.

## Stage Report: implementation (cycle 2)

- DONE: Replace the disconnected test-local wait predicate with an exact-head, real-adapter-driven Codex matrix.
  Commit `51eb5da25` removes the local predicate oracle and adds a live matrix that installs the current checkout through the real Codex launcher, verifies the installed source head, and reads the resulting Codex rollout function calls.
- DONE: Prove exact `wait_agent(timeout_ms: 300000)` behavior for active unresolved workers only, and demonstrate that an inverted shipped adapter fails.
  The exact-head matrix passed active-unresolved/completed/errored/absent: the active case emitted exactly one `wait_agent` call with exact arguments `{"timeout_ms":300000}`, while every terminal or absent case emitted none. In detached throwaway checkout `/tmp/spacedock-ej-invert.FK2i4R/worktree`, reversing only the shipped adapter rule made the active-unresolved subtest fail with observed calls `0/0`, wanted `1/1`; this proves the matrix is falsifiable through the delivered adapter rather than a duplicated formula.
- DONE: Preserve the ordered-loop behavior and authorized surface, then rerun focused, full, race, formatting, and relevant live verification.
  The candidate remains exactly eight files at +240/-10. Focused event-loop/gate/wait and contractlint suites passed; `gofmt -w ./cmd ./internal` left the worktree clean; `go test ./...` and `go test ./... -race` passed with immutable `SPACEDOCK_STATE_ROOT` commit `73f41e2a2232ebb561710bce568641ec976d5f3d`; and the exact-head Codex live matrix passed all four cases. The previously recorded mutable live-state pilot-manifest drift remains separate from candidate behavior and is not represented as green.

### Summary

Cycle 2 closes the validation finding without changing the event-loop contract: AC-3 now observes the real installed adapter and exact Codex rollout call, rejects the inverted adapter, and retains the active-only wait rule. The implementation and its adapter-bound evidence are ready for validation re-entry.

## Review-finding disposition (validation cycle 2)

- Finding: the AC-3 live observer assumes the isolated Codex home is under `artifactRoot/_codex-home`, but the shared runner deliberately uses the user-cache or repo-adjacent fallback when artifacts are temporary or inside the checkout. The supported default local and CI layouts therefore cannot find any rollout file.
- Released user and normal workflow: maintainers running the required exact-head Codex lane locally with no artifact override, and CI storing artifacts inside the checkout.
- Observable harm: all four AC-3 cases fail before their wait cardinality is graded, so the required evidence cannot gate the shipped adapter in supported validation environments.
- Authority: `contract[docs/dev/README.md#validation]` requires runnable proof at the delivered observation boundary; an environment-specific pass that fails in the normal local and CI layouts leaves AC-3 unproved.
- Trigger evidence: exact candidate `51eb5da25` fails active, completed, errored, and absent 4/4 with `session rollout ... = [], want one` under the default layout. With an external artifact root the same candidate passes 4/4, proving path selection—not adapter behavior—is the boundary failure.
- Proposal: Material evidence defect owned by EJ. Plumb `newCodexLiveIsolatedHome`'s actual `codexHome` onto `codexLiveRunner` and pass that path to `observedCodexWaitCalls`; do not add a scheduler, controller, command, state field, or completion ledger.
- Recommendation: REJECTED until the adapter-bound matrix passes in the supported default layout. The immutable suite, mutable pilot drift, and Pi account limitation remain separate and unchanged.

## Stage Report: validation (cycle 2)

- FAILED: Re-review exact candidate 51eb5da25 and prove AC-3 with the real shipped Codex adapter: active unresolved emits exact wait_agent(timeout_ms: 300000), while completed, errored, and absent emit none.
  The external-root run proves the adapter result exactly (active `1/1` exact-300000; completed/errored/absent zero), but the required default-layout run fails all four before observation because its rollout glob uses the wrong Codex-home root.
- DONE: Invert the shipped adapter in a detached checkout and require the adapter-bound matrix to fail; regress AC-1, AC-2, AC-4, and AC-5 ordered-loop evidence.
  Detached `/tmp/spacedock-ej51-audit.XuxXAd` inversion fails active with `0/0` wanted `1/1` and completed with one forbidden wait; AC-1/AC-4 mod+three-gate-before-next and spawn separation, AC-2 single retry/truthful stop, and AC-5 independent-a/b dispatch all remain green.
- DONE: Confirm the eight-file +240/-10 candidate stays within approved tolerance and passes focused, contractlint, formatting, immutable full/race, and relevant exact-head live evidence.
  Exact HEAD and surface match; `git diff --check`, `gofmt -d`, focused suites, contractlint, and immutable `SPACEDOCK_STATE_ROOT=73f41e2a` full/race pass. Exact-head Codex live is behavior-green only with the external-root workaround and fails in the supported default layout.

### Summary

Validation cycle 2 recommends REJECTED for a narrow Material AC-3 evidence defect in rollout-path attribution. The shipped adapter behaves correctly and the live oracle detects inversion, but the candidate must carry the runner's actual isolated Codex-home path into the observer so the same proof works by default and in CI.

## Stage Report: implementation (cycle 3)

- DONE: Plumb the actual isolated codexHome path from newCodexLiveIsolatedHome through codexLiveRunner into observedCodexWaitCalls.
  Commit `5fa79d480` retains the created home on the live runner and makes rollout discovery start at that exact path; reverting either handoff reproduces `session rollout ... = [], want one` under default temporary artifacts.
- DONE: Make the real-adapter wait matrix pass all four cases under default local/CI layout, retain the external-root control, and keep detached adapter inversion failures.
  Exact-head default and external-root runs both pass active `1/1` exact-300000 plus zero completed/errored/absent calls; detached `/tmp/spacedock-ej-pathfix-invert.ndZf4W/worktree` fails active `0/0` wanted `1/1` and completed with one forbidden wait when only the shipped adapter rule is inverted.
- DONE: Preserve adapter behavior and approved scope; rerun focused, contractlint, formatting, immutable full/race, and relevant exact-head live evidence.
  The nine-file +242/-11 candidate changes only two live-test harness paths in this cycle; focused ordered-loop/gate/wait/home tests, contractlint, gofmt, diff-check, and immutable `SPACEDOCK_STATE_ROOT=73f41e2a2232ebb561710bce568641ec976d5f3d` full/race all pass, with no adapter, command, state, scheduler, controller, or completion-ledger change.

### Summary

Cycle 3 fixes only rollout attribution by carrying the actual isolated Codex home through the existing runner. AC-3 now runs and grades the real adapter in both supported default and external artifact layouts, while the unchanged negative control still rejects both inverted decisions.

## Stage Report: validation (cycle 3)

- DONE: Re-review exact candidate 5fa79d480 and prove the shipped Codex adapter matrix passes active/completed/errored/absent under both default local/CI and external artifact layouts using the runner's actual codexHome.
  AC-3 passes twice at exact HEAD: default cache-backed and external artifact-root runs each observe active `1/1` exact `wait_agent(timeout_ms: 300000)` and zero waits for completed, errored, and absent; reverting the codexHome handoff would make rollout discovery empty.
- DONE: Require detached adapter inversion to fail the active and completed cases; regress AC-1, AC-2, AC-4, and AC-5 ordered-loop evidence.
  Detached `/tmp/spacedock-ej5fa-audit.f5Peg0` inversion fails active with `0/0` wanted `1/1` and completed with one forbidden wait; AC-1/AC-4 retain mod plus all three gate routes before next with no gate spawn, AC-2 retains one idle/reconcile/retry and truthful stop, and AC-5 dispatches both independent tasks.
- DONE: Confirm the nine-file +242/-11 candidate stays within approved tolerance and passes focused, contractlint, formatting, immutable full/race, and relevant exact-head live evidence.
  Exact `5fa79d480`, `git diff --check`, `gofmt -d`, focused event-loop/status/dispatch/contractlint suites, and immutable `SPACEDOCK_STATE_ROOT=73f41e2a` full and race pass. One unrelated empty-body dispatch parity result passed five isolated reruns and the full rerun; mutable pilot drift and Pi account limits remain separate prior evidence.

### Summary

Validation cycle 3 recommends PASSED with no material, deferred-risk, or polish findings. AC-1 through AC-5 have falsifiable command, state, or exact-head live evidence; the candidate stays within approved scope and adds no scheduler, controller, command, state field, or completion ledger.

## Stage Report: implementation (cycle 4)

- DONE: Commit a minimal current-main reconciliation that preserves EJ event-loop ordering and wait semantics.
  Merge commit `c7f3bc0d8` resolves the sole `fo-dispatch-core.md` overlap by retaining main's canonical `dispatchable+ready_gates` retry rule and EJ's ordered all-actions stop/wait rule; `origin/main` is an ancestor and no conflict remains.
- DONE: Keep the approved nine-file EJ surface; do not add scheduler, command, state, or authority behavior.
  Against current `origin/main` (`be0e8453e`), the candidate remains the approved nine files at +242/-11; the reconciliation adds no candidate file or semantic surface.
- DONE: Run focused checks plus gofmt, full, and race suites; record the exact new candidate head and any residual conflict.
  Exact head `c7f3bc0d890c6fc380ea710d133b6ac0505c278b` passes focused event-loop/wait, status, and contractlint tests, `gofmt -w ./cmd ./internal`, `git diff --check`, full and race suites with immutable state `a0169cc2d8a5e4912ed33f75ca8422a767e71c9e`, and the default-layout exact-head Codex live wait matrix; residual conflict: none.

### Summary

Reconciled EJ onto current main with one manual conflict resolution that composes main's ready-gate scheduler envelope with EJ's explicit event-loop ordering and unresolved-worker wait semantics. The exact candidate remains nine files at +242/-11 and is ready for fresh validation; an initial full-suite run against obsolete state snapshot `73f41e2a` exposed only expected pilot-manifest drift and was replaced by the compatible immutable snapshot cited above.

## Stage Report: validation (cycle 4)

- DONE: Re-review exact candidate c7f3bc0d890c6fc380ea710d133b6ac0505c278b and prove it merges cleanly with current origin/main while preserving ordered event-loop and unresolved-worker wait semantics.
  Fresh fetch leaves `origin/main=be0e8453e`, which is an ancestor of exact HEAD; `git merge-tree --write-tree origin/main HEAD` exits 0. The mixed trace retains reconcile/drain, mod and all three gate routes before next, both independent dispatches, one retry, and truthful stop.
- DONE: Verify every EJ acceptance criterion with independent evidence, including active, completed, errored, and absent Codex worker states across supported artifact layouts.
  AC-1/4/5 pass exact ordered/spawn and real status/CLI gate-route checks; AC-2 passes exact retry/released and unchanged-empty traces; AC-3 passes default and external-root live runs, with active waiting exactly once at 300000 ms and completed/errored/absent waiting zero times.
- DONE: Confirm the candidate remains nine files at +242/-11 and run focused, contractlint, formatting, full, race, and applicable exact-head live checks.
  Exact HEAD remains clean at nine files +242/-11 with unchanged SHA-256 hashes; focused ensigncycle/status/CLI/contractlint, `git diff --check`, `gofmt -d`, immutable-state full/race, and both Codex live layouts pass.

### Summary

Validation cycle 4 recommends PASSED with no material, deferred-risk, or polish findings. The current-main reconciliation is conflict-free, preserves all AC-1 through AC-5 behavior, and leaves the exact candidate bytes unchanged.
