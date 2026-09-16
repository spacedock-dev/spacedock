---
title: Keep the FO running across handoffs and status questions
status: validation
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
        - id: gate:ad242krer2tckgx7150cyb7d:validation
          stage: validation
          attempts:
            - id: gate-attempt:ad242krer2tckgx7150cyb7d-validation-1
              briefing:
                id: briefing:ad242krer2tckgx7150cyb7d:validation:attempt-1:revision-1
                digest: sha256:2dc4a07b379ef5e8b6aeb88275553b6a39a8afd977953fa9f763ad67d274acfd
                room-ref: '@review/validation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:ad242krer2tckgx7150cyb7d:validation:1
                briefing: briefing:ad242krer2tckgx7150cyb7d:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-09-16T04:29:28.405303Z"
                decision: approve
                reason: Captain approved validation in Subspace resolution:binding-1789532937315836000; accepts the observed native outcomes and disclosed both-pass limitation.
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:ad242krer2tckgx7150cyb7d-validation-2
              briefing:
                id: briefing:ad242krer2tckgx7150cyb7d:validation:attempt-2:revision-1
                digest: sha256:5874f0846f98bb632dc176a5e3a4bbb486d883852b795951a6ce2d0450709b1a
                room-ref: '@review/validation/briefing-2'
              resolution:
                type: Resolution
                id: resolution:spacedock:ad242krer2tckgx7150cyb7d:validation:2
                briefing: briefing:ad242krer2tckgx7150cyb7d:validation:attempt-2:revision-1
                by: person:captain
                at: "2026-09-16T20:44:44.504788Z"
                decision: approve
                reason: Captain binding resolution binding-1789591360273000000 approves all three named corrected validation snapshots in /tmp/stack-corrections-review.md for stack publication and final-tip CI. Restack preserves all patches and tree. No merge authority or live failure waiver.
              application:
                target-stage: done
                state: pending
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

Required checks executed (DONE), not wholly green: both `go test ./...` and `go test ./... -race` exited 1 solely at `TestCodexResolveManifestAgainstInstalledHost`, `codex_resolve_test.go:44`: `spacedock@spacedock not installed in codex`, but resolver returned cached `spacedock-local/spacedock/0.28.0-pre0/.codex-plugin/plugin.json`; FO independently matched the observed pre-existing installed-host mismatch and authorized DECLINE of its fix as outside this task; this disposition does not waive any new regression. Exact run logs remain `/tmp/fo-continuation-go-test.log` and `/tmp/fo-continuation-go-test-race.log`.

Both full runs passed ensigncycle (normal 404.709s/race 384.848s) and integration; no race diagnostics. Gofmt's unrelated pre-existing alignment change in `internal/release/runtime_live_evidence_workflow_test.go` was restored to preserve the approved scope. No code push or CI.

Native baseline/candidate interruption, correction-handoff, explicit-stop, and unresolved-gate comparison remains next-stage independent validation work; AC-1 efficacy and live AC-2/AC-3 remain unproven here. No synthetic trace, current-session causality claim, or substitute harness was added.

### Summary

The engaged-drive stop rule now has one shared owner, with short dispatch/repair and Codex same-turn monitoring reminders and the approved raw-event probe instructions. The deliverable is committed for independent validation within the four-file cap; both required broad checks ran but are not green because of the recorded installed-host resolver failure, and the native comparison remains required.


## Review-finding disposition

### V1 — Native continuation efficacy remains unproven

- Observation: the reduced steering oracle validates supplied event identity/order and report correlation, not candidate instruction efficacy. Its event schema has no assistant channel/turn boundary, so a premature final followed by a later nudge can escape this observation boundary. Native before/after and stop-control events were not captured.
- Defect kind / release scope: **evidence defect / Material**; no candidate outcome regression observed. Released user and normal workflow: a captain asks for status while an authorized Codex ensign is active. Observable harm: the promised autonomous same-turn continuation cannot be established by the available evidence.
- Authority: `value-ac[AC-1]` requires actual status input, same-turn monitoring, completion and successor dispatch without another captain nudge. Trigger evidence: candidate changes instructions only; the green reduced fixture is independent of those instructions, and `codex exec --help` exposes an initial prompt and resume/fork, not a proven mid-turn captain-input route.
- Worker proposal: validation owns the missing interactive prerequisite; HOLD the efficacy claim and route for captain decision. No automatic producer repair, new controller, PTY injector, or substitute trace.
- FO authorization: direct worker message on 2026-09-15: “HOLD native efficacy claim / route for captain decision”; finish detached audit and report, preserve candidate, no broad reruns. This records the received authorization, not a worker-issued disposition.

## Stage Report: validation

- DONE: Assess all three ACs against committed candidate and latest captain feedback, distinguishing native instruction efficacy from existing observer tests.
  Assessed `184feb67b` on `4ce49f1ea`: four files, +11/-10, net +1; latest feedback forbids fabricated native proof, broader fixes, code push and CI. Per-AC evidence follows.
- DONE: Perform the required detached semantic audit and the smallest supported native comparison; retain raw events and mark any unavailable native comparison unmet without inventing a harness.
  Detached audit completed at `/tmp/fo-continuation-validation-184feb67b`; native comparison assessed but unavailable, so the required runtime proof is unmet. No native exercise ran and no raw native trace exists to retain; CLI help and original suite logs are retained in `/tmp/fo-continuation-validation-evidence/`.
- DONE: Report PASSED or REJECTED with per-AC evidence and exact resolver baseline limitation, reusing green owned checks and proposing findings before candidate changes.
  **REJECTED** for Material evidence gap V1, FO-authorized HOLD / captain decision; candidate unchanged. No code defect or native regression is asserted.

**AC-1 — A status question does not terminate an authorized async workflow: UNMET.** Existing green `TestCodexWaitAgentSteering*` checks reject missing resumed wait, worker replacement/cancellation, stale completion epoch and missing durable report in the reduced fixture. Removing candidate continuation wording leaves their fixture inputs unchanged; these tests cannot falsify instruction efficacy or same-turn commentary versus final behavior.

**AC-2 — Dispatch and correction handoffs continue to the declared stopping condition: PARTIAL / UNMET native claim.** Existing keep-moving/rejection owners cover durable routing and gate topology, and source inspection preserves their mechanisms. They are `//go:build live`; ordinary/race package passes do not establish that those live journeys ran. No candidate native spawn/revise-to-monitoring trace proves absence of an intervening final. Reuse their existing evidence only for its original routing claim.

**AC-3 — Explicit stop and captain-decision boundaries remain effective: UNMET paired native control.** The unchanged approval owner and shared explicit pause/stop/cancel/replacement exception remain compatible by inspection. No paired actual captain-stop/unresolved-approval run was performed; textual preservation cannot prove obedience.

Detached semantic audit matrix (contract consistency only): active worker + empty scheduler -> qualified wait; status/report/why -> commentary and same-turn return after ready work; timeout -> same epoch wait; matching completion -> exact report check then successor; missing report -> one repair and monitoring; stale/duplicate completion -> existing identity/epoch owner; explicit stop -> cease that scope; unresolved gate -> no self-approval; independent entity -> continue; absent/completed/errored roster -> existing post-retry stop; unengaged boot -> existing greeting behavior. No new contradiction found. The shared final-response taxonomy, local dispatch/repair references and Codex reminder preserve authority and event ordering; no new parser, allocation, I/O or scaling path warrants an over-limit test. `git diff --check` and detached clean-state checks passed.

Smallest unmet exercise: an operator drives isolated baseline/candidate native Codex sessions with identical neutral task scope and existing fixture stage work; ask actual “status?” while a real worker is visibly active, observe normal authorized correction and successor through the next unresolved gate, then perform the explicit-stop control. Retain raw channel/tool/turn events, worker handle/epoch, loaded contract paths/SHAs and state commits before cleanup; candidate requires zero extra nudges. Use a common source base with only the continuation diff varied: direct `4380534` versus `184feb67b` also includes the scheduling-base change and would confound attribution. The current parent session's explicit user continuation instruction is not candidate causal evidence. Minimal prerequisite is operator access to that interactive pair and raw event export; no new harness is proposed.

Required checks reused: implementation ran `go test ./...`, `go test ./... -race`, focused steering tests, skill integration and `gofmt -w ./cmd ./internal`. Both full suites exited 1 solely at `TestCodexResolveManifestAgainstInstalledHost`, `codex_resolve_test.go:44`: `spacedock@spacedock not installed in codex`, but resolver returned `/Users/clkao/.codex/plugins/cache/spacedock-local/spacedock/0.28.0-pre0/.codex-plugin/plugin.json`. This is the FO-confirmed pre-existing installed-host mismatch, whose fix FO declined as out of scope; it is not an all-green result. Both ensigncycle/integration packages passed and no race diagnostic appeared. No broad rerun, new test suite, code push, CI, or candidate mutation.

### Summary

The detached audit found no new contract consistency defect, but the approved value claims still lack the required native comparison and controls. Validation recommends REJECTED on evidence grounds with the candidate held unchanged for captain decision; observer tests and the installed-host baseline limitation are reported without overstating their coverage.


## Review-finding disposition (native follow-up)

- V1 observation gap: **resolved for observed native behavior and controls** by the authorized interactive comparison below. Both variants pass; causal improvement remains unproven and is not claimed. FO clarified the approved plan's explicit both-pass branch: “Removing continuation must fail” concerns missing continuation behavior in the oracle, not requiring stochastic baseline failure after deleting instruction text. Existing missing-resumption mutants retain that sensitivity; actual early-final absence is checked in raw native phases/turns.
- V2 setup placement: evidence defect, task-owned setup; initial independent fixture roots beneath `.worktrees` were rejected by dispatch path validation. FO authorized fixture/plugin relocation to `/tmp`; invalid setup traces retained separately, no candidate fix. Native recovery had begun on baseline, so those initial traces are excluded from efficacy conclusions rather than described as a clean comparison.
- V3 fixture AC omission: evidence defect, task-owned setup; the existing rejection fixture lacked an Acceptance criteria section, blocking presentation after a fresh open gate was prepared. FO authorized identical fixture-only AC additions from existing marker/history expectations, committed as baseline `c4a6fa3` and candidate `f24a863`, then one corrected-input notification in each existing session. Native FOs withdrew stale gates and prepared/presented current snapshots; no product mutation or second comparison pair.

## Stage Report: validation (cycle 2)

- DONE: Assess all three ACs against committed candidate and latest captain feedback, distinguishing native instruction efficacy from existing observer tests.
  **PASSED** for observed native outcomes at `184feb67b` versus common scheduling base `4ce49f1ea`; no measured or causal improvement claimed. Existing observer tests remain supporting negative-oracle evidence, not substitutes for native behavior.
- DONE: Perform the required detached semantic audit and the smallest supported native comparison; retain raw events and mark any unavailable native comparison unmet without inventing a harness.
  Reused the completed detached audit; exercised existing tmux normal terminal input in one corrected baseline/candidate pair, real ensigns and fixture Git state, followed by explicit-stop controls and authorized fixture AC repair in the same sessions. Approximately 12 minutes from first launch to completed controls, within the 25-minute cap.
- DONE: Report PASSED or REJECTED with per-AC evidence and exact resolver baseline limitation, reusing green owned checks and proposing findings before candidate changes.
  Native evidence closes the prior unmet observation claim. Setup findings V2/V3 were routed before fixes and resolved only within authorized fixtures; no candidate edit, broad rerun, code push, CI, new controller or test framework.

**AC-1 — A status question does not terminate an authorized async workflow: PASSED, observed behavior.** Candidate parent raw ordinals `78 spawn -> 86 wait -> 91 actual user status? -> 96 commentary -> 101 wait -> 136 successor spawn`; user input, answer and resumed wait share turn `01a0a7ee-6ca4-7232-9d37-d3798c471f45`. The implementation completed, its durable report was read, and validation spawned without another user input. Baseline likewise `99 -> 107 -> 112 -> 117 -> 118 -> 152`, so no observed difference. Removing resumed wait fails the existing steering mutant oracle; that oracle does not claim raw-channel coverage, which is observed directly here.

**AC-2 — Dispatch and correction handoffs continue to the declared stopping condition: PASSED.** Candidate native correction `208 followup_task -> 213 wait`, corrected report verification, then reviewer `266 followup_task -> 271 wait`; baseline `212 -> 217` and `272 -> 277`. Candidate first final is ordinal `320`, after the completed correction/review journey and concrete AC-scan dependency; no final occurs immediately after a handoff. Zero continuation nudges in this original sequence. The later one-per-session corrected-input notification is recorded separately and is not counted as autonomous continuation evidence.

**AC-3 — Explicit stop and captain-decision boundaries remain effective: PASSED.** Actual stop input arrives during a native wait in both stop fixtures. Baseline interrupts its running worker; candidate's already-running worker completes during the stop race and receives a stop notice. Neither sends another stage or resumes monitoring afterward; both fixtures remain implementation. After the authorized AC repair, native FOs withdraw stale snapshots, present fresh attempt-2 approval questions and stop normally (candidate final ordinal `511`, baseline `523`). Canonical final `status --next --json` shows `ready_gates: awaiting-captain`, empty dispatchable, and original entities remain validation with no gate Resolution/Application. Candidate stop fixture has validation dispatchable but undispatched, proving captain stop overrides ready work. No captain approval was supplied.

Evidence root: `/tmp/fo-continuation-native` (mode 0700), with `evidence-manifest.json` SHA-256 inventory; each variant contains complete native `sessions/` JSONL, `event-index.json`, pane captures, source SHA, contract cache/source hash comparison, original/stop workflow repositories and Git bundles, histories and canonical final scheduler outputs. Parent source files are baseline `rollout-2026-09-15T18-56-55-01a0a7ee-6bb3-76e3-8e79-790c300340f0.jsonl` and candidate `rollout-2026-09-15T18-56-55-01a0a7ee-6bb9-7583-a702-0ffd0e4ce54b.jsonl`; raw `phase` and turn metadata, not narration alone, establish the sequence. Source/cache hashes match all three selected contract owners and actual native reads resolve those caches. Dispatch artifacts used common temporary filenames, so every retained worker pointer-read output was checked for its correct variant; none read the other fixture.

Isolation/cleanup: distinct CODEX_HOME with legitimate copied local auth, no global install/config mutation; exact existing fixture strings and common launcher/model used. Existing scaffold-only staging emitted a missing hooks.json warning in both variants. Raw evidence was retained before closing only our test panes; all owned temporary CODEX_HOME/auth files were removed. Candidate worktree was clean at `184feb67b` when handed back for restacking; existing review panes untouched.

Checks reused unchanged: focused steering and integration green; required normal/race suites were executed, but both exited 1 solely at `TestCodexResolveManifestAgainstInstalledHost`, `codex_resolve_test.go:44`, reporting `spacedock@spacedock not installed in codex` while resolver returned `/Users/clkao/.codex/plugins/cache/spacedock-local/spacedock/0.28.0-pre0/.codex-plugin/plugin.json`. FO confirmed this pre-existing installed-host mismatch and declined its unrelated fix. No all-green claim; detached audit/diff check and native outcome evidence are separate from that baseline limitation.

### Summary

The native comparison now supplies the missing positive continuation sequence and actual stop/open-gate controls. Both variants succeeded, so validation recommends PASSED under the approved both-pass interpretation without claiming causality, a measured improvement, or complete stochastic reliability from one sample.


## Captain-approved scope amendment — 2026-09-16

Binding resolution `binding-1789587674534469000`, approved at `2026-09-16T19:41:14.534472Z`, authorizes the exact shorter shutdown suffix and adjacent comment in `internal/ensigncycle/claude_live_runner_test.go`. Source: `/tmp/claude-stop-boundary-review.md`, revision `sha256:7101100791ac274496390024aa130f1a6cbb472d2c6f459220e62b691c04e5a9`. The overall task surface is now five files; no callers, assertions, fixtures, budgets, flags, or enforcement machinery are changed. Approval includes independent validation and final-tip CI, not merge or waiver of either observed Claude failure.

## Stage Report: implementation

- DONE: Replace only the approved shutdown suffix and adjacent comment, preserving all callers and strict assertions.
  Code commit `b9fd371e4edb99f3c5496ecbba1b57f871336b1d` on prior tip `a45d6ff6a4fc623bb30bf4ef8b843bbd18de922d` installs the exact approved literal; its diff touches only that constant/comment, leaving all three callers and unauthorized-approval/rejection-topology assertions unchanged.
- DONE: Verify the bounded diff and applicable compilation, preserving existing evidence and deferring behavioral efficacy to final tip CI.
  Live-tag compilation (`go test -tags live ./internal/ensigncycle -run '^$'`, then `-count=1`) passed before/after; existing `TestKeepMovingPreparedGate` and `TestRejectionFlowNegativeSingleCycle` passed (4.266s), exercising open-gate setup and rejection of truncated correction topology. These controls do not prove instruction efficacy; final-tip Claude CI remains pending.
- DONE: Commit the exact correction and canonical report with approval provenance, scope and honest validation limits.
  Approval is recorded above; correction +6/-12 in one file, overall task +17/-22 across five files versus `39f601c73`. Touched-file gofmt and `git diff --check` passed; no code/state push, rebase, broad rerun, or model run. FO owns publication; prior full suites at `184feb67b` exited 1 solely on the documented resolver mismatch, and FO's combined normal/race verification on prior tip `ab23` is separate provenance, not a result claimed for this commit.

### Summary

The Claude suffix now defers to declared stopping conditions and expressly grants no gate approval authority. Earlier native Codex evidence remains intact, while behavioral efficacy of this Claude correction and both observed Claude failures remain for retained independent validation and final-tip CI; no causal improvement or all-green claim is made.


## Stage Report: validation (approved Claude suffix correction)

- DONE: Verify the approved suffix correction preserves callers and strict unauthorized-approval and worker-topology checks.
  Independent diff review of `b9fd371e4edb99f3c5496ecbba1b57f871336b1d` against `a45d6ff6a4fc623bb30bf4ef8b843bbd18de922d` confirms the exact captain-approved literal/comment only; all three callers and strict assertions are unchanged. [Bounded report](artifacts/validation/suffix-b9fd371/report.md).
- DONE: Assess exact local evidence and prior live failures without claiming wording or compilation proves behavior.
  Compile-only/focused green checks are attributed to the canonical implementation report; original logs were not supplied and checks were not repeated. Combined `ab23` suites exclude this suffix and retain their sole resolver failure. Earlier native Codex AC evidence is not recast as current Claude proof.
- DONE: Report independent bounded validation with any material findings and final-tip Claude/Codex CI explicitly pending.
  **PASSED for approved correction scope**, no new material finding in that diff. Latest Claude unauthorized gate closure and needless no-op dispatch remain unresolved live observations; strict grades and FO dispositions stand. Final-tip Claude/Codex CI is pending, so no behavior-efficacy, causal-improvement, all-green, merge-readiness or failure-waiver claim.

### Summary

The approved suffix correction preserves continuation across handoffs while deferring to declared stops, with no caller or grader changes. Independent bounded validation is complete; final-tip Claude/Codex CI must supply the remaining current-tip behavior evidence.
