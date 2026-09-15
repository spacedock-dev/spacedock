---
title: Keep the FO running across handoffs and status questions
status: implementation
source: Captain-approved continuation correction, 2026-09-15
started: 2026-09-15T22:39:21Z
completed:
verdict:
score: 0.95
worktree: .worktrees/spacedock-ensign-fo-continuation
issue: spacedock-dev/spacedock#735
pr:
mod-block:
id: ad242krer2tckgx7150cyb7d
gates:
    version: 1
    records:
        - id: gate:ad242krer2tckgx7150cyb7d:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:ad242krer2tckgx7150cyb7d-ideation-1
              briefing:
                id: briefing:ad242krer2tckgx7150cyb7d:ideation:attempt-1:revision-1
                digest: sha256:4bcd7ba842757abf40653c51fcc102aa8b043d60e9c7306557ff1f3e6de48a18
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:ad242krer2tckgx7150cyb7d:ideation:1
                briefing: briefing:ad242krer2tckgx7150cyb7d:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-09-15T23:10:07.962105Z"
                decision: approve
                reason: 'Prepared-room binding resolution binding-1789513783661268000 approves focused continuation change cap20net/four files. Captain asks why pre0 was read: session catalog pins cachedpre0, while implementation baseline is currentmain4380534 including798. Native behavior proof remains required.'
              application:
                target-stage: implementation
                state: consumed
---

Make the existing FO continue authorized work without a user reminder or /goal. Put the next-action rule at dispatch, completion, revision-routing and user-interruption boundaries.

## Problem

In this session the FO repeatedly ended its turn after spawning workers or sending revisions back. It also answered status questions with final responses while authorized async work remained. The contract already requires monitoring and immediate completion routing; these were principally compliance failures, not missing authorization.

The cached pre0 dispatch contract has a contradictory gated-successor stop clause. PR #798 merged at main 438053493838dc70c9478b3d991309d566783e85 and fixes that separate successor issue. Inspect current source before editing; do not reimplement its fix.

## Proposed approach

Use existing owners only:
- fo-dispatch-core.md: put the required next action beside dispatch and completion/revision handoffs. Handoff acknowledgment is commentary, not a final response.
- codex-first-officer-runtime.md: answer status/report/why in commentary and resume required monitoring in the same turn. Preserve explicit captain pause/stop/cancel/replacement.
- first-officer-shared-core.md: define final-response reasons once and reference them: named captain decision, concrete unmet dependency, or requested outcome complete. An active worker alone calls for monitoring, not a final answer.

Keep the edit small. No user or worker reminder protocol, new watcher, ledger, /goal requirement, scheduler or standing enforcement process. Workers keep their normal completion signals. Do not lengthen several files with duplicate versions of the same rule.

## Risk evidence

Existing GitHub #735 describes the same status-interruption failure. Independent diagnosis: /tmp/fo-premature-stop-diagnosis.md. Representative bad final responses in this session were 'Dispatched a read-only triage', 'Sent back to the owner ... before implementation', and 'worker is completing ... validation comes next', while work remained active. A real prepared approval review with no independent work is a valid captain stop; ordinary handoffs are not.

Inspect the refreshed source at origin/main, not root HEAD or only the installed cached skill. Identify exactly what remains after #798.

## Out of scope

Tool/harness changes, additional worker notification machinery, new goals, broad skill rewrites, automatic captain approval, and changing valid stopping boundaries.

## Current-source diagnosis and proposed wording

Baseline: `origin/main` at `438053493838dc70c9478b3d991309d566783e85` (#798). Its dispatch core already orders a gated successor's dispatch before gate preparation; preserve that fix. Issue #735 confirms the missing status-interruption clarification. The cached pre0 contract is not the implementation baseline.

Remaining gaps are local reminders at three existing boundaries, not missing monitoring machinery:

1. Shared core's cadence still says “unless that stage is a gate” and “Yield only when blocked on the async result with no other work”. Replace those two cadence bullets with:
   > A gate approval triggers advancement and successor dispatch under the dispatch core's Gate successor guard. Continue authorized work until the existing loop reaches a declared stop; an unresolved worker calls for the runtime's monitoring, not a final response.
   > Final responses require an explicit captain pause/stop/cancel/replacement, a named captain decision or concrete unmet dependency with no independent authorized work left, completion of the requested scope, or the loop's existing post-retry `no-dispatchable` stop. Status, report, and explanation questions do not change scope or create a stop.
   This removes the surviving successor exception without reopening #798's implementation or proofs. Keep the independent-entities bullet unchanged. The final-response rule is scoped to an engaged drive; preserve greet-and-stop boot.
2. Dispatch core step 7 currently says “Await the worker result per `«async-dispatch»`…”. Append:
   > After recording the handle, continue `«dispatch.next-action»()`; handoff narration follows the shared final-response rule.
   Append to its existing completion/report-repair boundary:
   > A report-repair or revision handoff returns to the same loop; sending it is not completion.
   These are references beside actual transitions, not copies of the stop taxonomy. Feedback routing still owns authorization, correction records, reviewer identity, and re-gating.
3. Codex wait notes currently end the interruption instruction with “When the FO becomes idle again, it MUST resume monitoring unresolved workers.” Replace that sentence with:
   > Answer status/report/why questions in commentary, route ready authorized work, and resume unresolved-worker monitoring in the same turn once idle. Apply the shared final-response rule; answering the question alone is not a stop.
   Retain existing completion attribution, timeout behavior, and captain control. Host wait duration follows higher-priority session limits (60 seconds here); this task does not redesign timeout policy.

## Expected surface and tolerance

Estimate net LOC change: +12, across 4 files; anticipated +19 insertions / -7 deletions. Tolerance: net +8 lines (maximum +20 net), no extra production files. Files: the three contract owners above and `docs/dev/codex-idle-notification-probe.md`. Reports and raw one-off validation evidence remain task artifacts in the state checkout, not additional production machinery. If implementation needs more, explain why before expanding.

Observable semantics: engaged FO runtime continuation and Codex response-channel choice become explicit. No command grammar, stored format, authority, worker signaling, gate consent, scheduler, launch behavior, or runtime capability changes. Preserve no-work and greet stops. Land above scheduling; no code push or CI until stack authorization.

### Proposed documentation diff

In `docs/dev/codex-idle-notification-probe.md`, async comparison step 4 currently says “If captain input resumes the FO's active loop, record the worker as unchanged and continue useful active-scope work. When the FO becomes idle again, resume monitoring the same unresolved worker.” Replace with:

> If a status/report/why question resumes the FO, record its commentary answer and same-turn return to monitoring after ready work. An explicit pause/stop/cancel/replacement instead ends or redirects the authorized drive. Record raw host events and durable state; a later captain nudge is not autonomous continuation.

This is the existing documentation owner of the changed host behavior; no new site page or public command documentation is needed.

## Acceptance criteria

**AC-1 — A status question does not terminate an authorized async workflow.**
Verified by: a supported runtime exercise starts a worker, receives a captain status question, answers, reinstalls monitoring in the same turn, observes completion and dispatches the ready successor without another captain nudge. Oracle is observed runtime/tool sequence and durable state, not source wording. Removing continuation must fail.

**AC-2 — Dispatch and correction handoffs continue to the declared stopping condition.**
Verified by: existing dispatch/rejection journey or the same focused exercise shows monitoring after spawn and revise dispatch, then completion verification and next action. A final answer immediately after the handoff fails.

**AC-3 — Explicit stop and captain-decision boundaries remain effective.**
Verified by: paired pause/stop and unresolved approval cases halt without unauthorized dispatch or approval. The same continuation rule must not override captain control.

## Test plan

Existing primary owners:

- `internal/ensigncycle/codex_wait_agent_steering_test.go` owns the interruption/unchanged-worker/resumed-wait oracle and negative mutants. Its `docs/dev/_evidence/codex-wait-agent-steering-semantics/2026-07-23-dogfood.json` is explicitly a reduced trace with a correlated completion tail, not a raw before/after status-question proof.
- `TestLiveCommonKeepMovingPosture` and `TestLiveCommonRejectionFlow` in `internal/ensigncycle/shared_live_runner_test.go`, with durable history and rejection topology assertions, own dispatch, correction/reviewer routing, and gate outcomes across supported hosts. #798 already extended these. Do not add duplicate journey fixtures.
- `TestLiveCodexWaitMatrixFromShippedAdapter` owns the active/completed/errored/absent wait choice, but intentionally exits after one wait; it cannot prove AC-1. `skills/integration/codex_idle_notification_test.go` validates evidence classification, not actual persistence.

No spike needed for a new mechanism: none is proposed. Native `spawn_agent`, `followup_task`, `wait_agent`, and mailbox delivery are bound in this session; the existing durable journeys and steering evidence establish the underlying path. Instruction efficacy is unproven and must be checked during validation, not inferred from these reads. The assignment permits specifying the concrete exercise at ideation; no behavior improvement is claimed here.

One focused, manual native-Codex before/after comparison serves AC-1 and the handoff portion of AC-2. Use two isolated copies of the existing keep-moving/rejection workflow fixtures, baseline pinned to 4380534 and candidate differing only in this contract/doc change. Reuse their stage definitions and binary-generated entities/dispatches; do not fabricate worker reports or a transcript. Start each drive with identical neutral scope: “Engage this workflow and drive the named task to its next captain decision.” Let real ensigns do the declared work. When the worker is visibly running, send the actual captain question “status?” once; no prompt may instruct the FO to continue or mention the expected rule. Route one normal authorized correction through the existing rejection scenario, then observe the resulting handoff too. No extra user nudge after either event.

Capture raw host conversation/tool events (including commentary/final channel and worker handle/epoch), source SHA and loaded contract paths, and state-checkout commits. Positive oracle: spawn/revise -> monitoring; captain question -> commentary -> same-turn monitoring; matching worker completion -> durable report read -> successor dispatch -> prepared unresolved captain gate. Count extra captain nudges: candidate must be zero. A final answer between handoff/question and a real stop fails, even if later activity flushes completion. If both variants pass, report no observed difference; do not claim causality from one sample. If baseline does not encounter an active worker at injection, rerun only that invalid sample.

Within this same bounded comparison, repeat the interruption with “Stop driving this task; do not dispatch another stage.” Candidate must cease automatic successor/revision dispatch and monitoring for that scope; an already-running worker may finish, which is not unauthorized new dispatch. Also leave the final approval unresolved: no gate consume, self-approval, or successor dispatch. These negative controls serve AC-3. Removing the interruption sentence must be capable of exposing the old premature-final behavior; removing stop handling must fail the negative control. Actual trace/state outcomes, not a hand-scripted simulator or wording grep, decide.

Cost: one baseline and one candidate workflow drive plus the short explicit-stop control in each, roughly 15–30 minutes total depending on worker latency; single-sample evidence, no recurring CI lane. Use native interactive steering, because the current headless `codex exec --json` runner takes a single prompt and supplies no proven mid-turn captain-input injector. Do not build an injector or invoke the unrelated behavior-diff duo runner (its worker-first sequential shape does not exercise an active async wait). If native steering cannot be exercised in validation, mark AC-1 unmet rather than substituting synthetic evidence.

Before implementation edits, run the existing focused steering negative tests and applicable skill integration smoke tests; no new command text is planned. After implementation, run required `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`; independent validation owns the one focused live comparison and reuses existing green journey evidence where its claims already match. No broad baseline during ideation, no new reminders/watchers/ledger/goals, and no expensive CI before stack readiness.

### Feedback Cycles


## Stage Report: ideation

- DONE: Identify the exact remaining continuation gaps after #798 and propose the smallest nonduplicative contract edit.
  Read origin/main 4380534 and issue #735; retained #798, replaced shared cadence ambiguity, and specified local dispatch/Codex reminders with one stop-rule owner.
- DONE: Exercise or specify a concrete supported before/after interruption-handoff proof with explicit-stop negative control.
  Specified one native Codex comparison with real workers, raw channel/tool events, durable state, no extra captain nudge, stop control, and unresolved approval; no live result claimed.
- DONE: Record exact surface/tolerance and existing proof ownership without new reminder or monitoring infrastructure.
  Estimated +12 net (+19/-7), four files, +8 net tolerance; named steering, keep-moving, rejection, wait-matrix and evidence-schema owners and their limits.

### Summary

The remaining correction is a small contract/doc change at existing transition boundaries. Ideation specifies the missing live proof without building a new harness or repeating #798's successor work; implementation and validation remain outstanding.

## Stage Report: implementation

- DONE: Apply the approved small continuation reminders to current source, preserving798 and one shared stopping taxonomy.
  Commit `184feb67b` atop scheduling base `4ce49f1ea` changes only the three approved contract owners and existing probe documentation; #798 successor guard and independent-entities cadence remain intact.
- DONE: Use existing focused proof owners and retain honest native interruption/stop evidence requirements without adding infrastructure.
  Before edits, `go test ./internal/ensigncycle -run '^TestCodexWaitAgentSteering' -count=1` and `go test ./skills/integration -count=1` passed; steering rerun passed after edits. Reduced-trace mutants reject cancellation/replacement/missing resumption, while integration rejects false idle-wake classification; neither proves candidate instruction efficacy.
- DONE: Commit within cap20net/four files with required checks and report for independent validation.
  Task-only diff from `4ce49f1ea`: four files, +11/-10, net +1; clean committed worktree and `git diff --check` pass. Required normal/race suites and `gofmt -w ./cmd ./internal` ran; suite failures are recorded below.
- FAILED: Required full-suite green result.
  Both `go test ./...` and `go test ./... -race` exited 1 solely at `TestCodexResolveManifestAgainstInstalledHost`, `codex_resolve_test.go:44`: `spacedock@spacedock not installed in codex`, but resolver returned cached `spacedock-local/spacedock/0.28.0-pre0/.codex-plugin/plugin.json`; FO reports an independently matching baseline. No candidate fix authorized or made.
- DONE: Preserve scope and verification boundaries.
  Both full runs passed ensigncycle (normal 404.709s/race 384.848s) and integration; no race diagnostics. Gofmt's unrelated pre-existing alignment change in `internal/release/runtime_live_evidence_workflow_test.go` was restored to preserve the approved scope. No code push or CI.
- SKIPPED: Native baseline/candidate interruption, correction-handoff, explicit-stop, and unresolved-gate comparison during implementation.
  Independent validation owns the approved native exercise; AC-1 efficacy and live AC-2/AC-3 remain unproven here. No synthetic trace, current-session causality claim, or substitute harness was added.

### Summary

The engaged-drive stop rule now has one shared owner, with short dispatch/repair and Codex same-turn monitoring reminders and the approved raw-event probe instructions. The deliverable is committed for independent validation within the four-file cap; both required broad checks ran but are not green because of the recorded installed-host resolver failure, and the native comparison remains required.
