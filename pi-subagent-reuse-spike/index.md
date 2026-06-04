---
title: Pi subagent reuse and feedback-to resume semantics
status: ideation
source: captain (2026-06-04) — Pi runtime can technically resume subagents, but Spacedock has not proven reusable worker semantics for feedback-to loops or token/cost tradeoffs
score: "0.33"
started: 2026-06-04T00:00:00Z
completed:
verdict:
worktree:
issue:
id: 7qcp0ttfhzq01cetbx85vcjx
---

# Pi subagent reuse and feedback-to resume semantics

Explore whether Spacedock's Pi runtime should support reusing/resuming an existing Pi subagent worker instead of always dispatching a fresh worker, especially for validation `feedback-to: implementation` loops.

## Problem

Spacedock's Pi runtime must preserve independent stage verification. A validation rejection is not an acceptance-finalization turn inside the same implementation agent; it is a separate workflow stage with external evidence, entity stage reports, state/product commits, and a `feedback-to: implementation` route when validation rejects implementation output.

`pi-subagents` can technically continue from a previous child session: its README documents `subagent({ action: "resume", id, message })`, and explains that completed runs are revived from a persisted child `.jsonl` session file rather than by restarting the same OS process. The Pi runtime notes already propose minimal worker metadata and a completion epoch so stale completions cannot satisfy later follow-ups. The unproven questions are operational and semantic:

- whether a resumed/forked Pi worker retains enough useful implementation context to make validation feedback cheaper or clearer;
- whether that retained context becomes stale, over-confident, or contradictory to the validator's new evidence;
- whether any token/time savings offset extra registry state, epoch handling, and lifecycle cleanup;
- whether direct `subagent({ action: "resume" })` is reliably available in the non-interactive Spacedock runtime path, or whether only lower-level session forking is available;
- how to prevent a previous completion/result from being mistaken for the new feedback turn's completion.

## Spike findings

### Local setup and safety determination

Live testing was safe enough to attempt with a temporary workflow only. The local machine had:

- `pi` CLI at `/Users/clkao/.local/state/fnm_multishells/28390_1780511451715/bin/pi`, version `0.78.0`;
- installed user packages `npm:pi-subagents` and `npm:pi-intercom` from `~/.pi/agent/npm/node_modules/`;
- Pi runtime docs present in `.worktrees/spacedock-ensign-pi-runtime-no-subagent-acceptance/`, not in the main checkout;
- a temp workflow at `/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/sd-pi-reuse-XXXXXX.1cMXZb8GWC/workflow` with state under `.spacedock-state/`;
- a temp Pi session directory at `/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/sd-pi-reuse-XXXXXX.1cMXZb8GWC/pi-sessions`.

No product code or active workflow state was intentionally mutated for the spike.

### Live/manual experiment

The initial run asked a parent Pi session to launch one worker against the temp workflow and append an `initial-marker` section to `.spacedock-state/reuse-marker.md`. The final parent report recorded:

- parent session: `/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/sd-pi-reuse-XXXXXX.1cMXZb8GWC/pi-sessions/2026-06-04T05-09-44-073Z_019e9109-9209-7427-ad14-beaab2105ad9.jsonl`;
- reported worker run id: `bcf031b1-6f0d-41fc-ab64-3b128857b28e`;
- reported worker session label: `subagent-worker-bcf031b1-6f0d-41fc-ab64-3b128857b28e-1`;
- child session file: `/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/sd-pi-reuse-XXXXXX.1cMXZb8GWC/pi-sessions/2026-06-04T05-10-46-339Z_019e910a-8543-7d16-8ca8-64e119d5af28.jsonl`;
- artifact pointer observed by parent: `/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/pi-subagents-uid-501/artifacts/bcf031b1-6f0d-41fc-ab64-3b128857b28e_worker_input.md`;
- durable marker evidence: `## initial-marker` with `UTC time: 2026-06-04T05:10:53Z`.

A direct natural-language request for `subagent({ action: "resume", id, message })` was not reliably executed as a visible `subagent` tool call in this non-interactive harness. The parent Pi process instead used the underlying reachable mechanism documented by `pi-subagents`: revive/fork from the stored child session file. That resume-equivalent follow-up appended a second durable marker:

- revived/forked session file: `/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/sd-pi-reuse-XXXXXX.1cMXZb8GWC/pi-sessions/2026-06-04T05-12-54-790Z_019e910c-7b06-7cc5-9af8-1f2fd1a439a3.jsonl`;
- durable marker evidence: `## resumed-feedback-marker` with `UTC time: 2026-06-04T05:12:59Z`;
- child-level elapsed evidence reported by the parent command: `9.896s`;
- outer parent-turn wall clock measured by the shell wrapper: `79s`.

For comparison, a fresh worker-style follow-up appended a third marker:

- fresh child session file: `/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/sd-pi-reuse-XXXXXX.1cMXZb8GWC/pi-sessions/2026-06-04T05-13-38-643Z_019e910d-2653-7f66-ae35-efe0c7ee4ece.jsonl`;
- durable marker evidence: `## fresh-followup-marker` with `utc: 2026-06-04T05:13:53Z`;
- child-level elapsed evidence reported by the parent command: `24.617s`;
- outer parent-turn wall clock measured by the shell wrapper: `70s`.

Token/cost telemetry was available in Pi JSON session messages, but the simple live run did not expose a clean child-only ledger for the exact resume-vs-fresh comparison. The parent session reported approximate turn costs such as `0.063992` for the resume report turn and `0.105097` for the fresh comparison report turn, while the fresh child session itself reported smaller per-message costs culminating around `0.009332` for its final response. These numbers are useful only as evidence that telemetry exists; they are not strong proof that reuse saves money because parent orchestration/cache effects dominate the recorded totals.

### Interpretation

The spike proves that a persisted Pi child session can be revived/forked and can produce a second durable artifact after a feedback-style follow-up. It does not prove that Spacedock should route normal `feedback-to` loops to the previous implementation worker by default.

Reasons to avoid first-class default reuse now:

- Direct `subagent({ action: "resume" })` was documented, but not proven as a clean, visible tool invocation in the exact non-interactive harness path used here.
- The resume-equivalent path depends on durable child session files and may create a new child process from old context, so Spacedock still needs registry, epoch, state, and cleanup logic; reuse is not a free in-memory continuation.
- Validation feedback should be interpreted against current validator evidence and current checkout state. A reused implementation worker may retain stale assumptions, previous success framing, or pre-feedback task boundaries.
- Token/cost evidence is inconclusive. The child-level elapsed sample favored revive/fork (`9.896s` vs `24.617s`), but parent wall-clock overhead was comparable and the telemetry was not isolated enough to justify workflow complexity.
- Fresh dispatch is simpler to reason about, easier to validate independently, and matches the existing Spacedock model where validation is a separate agent/stage rather than same-agent acceptance finalization.

## Proposed approach

Recommendation: **defer first-class reuse and keep fresh redispatch as the workflow default. Support resume only as manual/debug or explicitly opt-in experimental tooling until a later live harness proves clean direct resume, durable registry behavior, and token/cost benefit.**

For normal Pi `feedback-to: implementation` loops:

1. Validation rejection writes durable findings in the entity validation report.
2. The first officer routes a fresh implementation dispatch that includes the validator's concrete evidence and required fixes.
3. The fresh implementation worker treats the feedback as a new cycle and writes a new implementation stage report/commit.
4. The registry may record previous Pi worker handles for observability and manual debugging, but the scheduler must not choose reuse automatically.

For manual/debug resume only:

- require an explicit operator/debug flag or command, never an implicit `feedback-to` default;
- require an existing worker registry record with worker label, substrate, run/session handle, entity slug, stage, state, and completion epoch;
- increment the epoch before sending the follow-up;
- mark the worker active for the new epoch;
- accept completion only when the returned evidence matches the current worker label, run/session handle, and epoch, and when the entity/stage report or file evidence exists on disk;
- record that resume was used in the stage report or run ledger so validation can distinguish fresh vs reused context.

This preserves compatibility-first behavior while keeping the useful part of the spike: Pi resume/fork can be a debugging affordance and future optimization candidate, not a semantic dependency of Spacedock validation.

## Failure modes and lifecycle model

- **Stale completion:** an old child result arrives after a follow-up starts. Mitigation: completion epochs; old epoch cannot complete new epoch.
- **Stale context:** resumed worker remembers pre-validation assumptions or old filesystem state. Mitigation: fresh dispatch by default; if manual resume is used, prompt must include current entity path, stage, feedback, checkout/worktree, and evidence to reread.
- **Handle drift:** run id exists but child session file is missing, cleaned up, or belongs to a different entity/stage. Mitigation: registry must store entity slug/stage/state and verify the session file exists before resume.
- **State collision:** a resumed worker writes to an old worktree or state checkout. Mitigation: current workflow dir, entity path, state checkout, and worktree path must be in the follow-up prompt and verified in completion evidence.
- **Completion ambiguity:** `resume` may revive a new child process with a different resulting session/run identity. Mitigation: registry stores both original handle and current epoch's observed completion evidence; completion is not just transcript text.
- **Cleanup:** temp child sessions/artifacts can accumulate. Mitigation: lifecycle policy may retain only recent handles for debugging and mark records closed after terminal validation/done; never require historical handles for normal workflow progress.
- **Budget surprise:** forked/resumed context can include large old transcript context. Mitigation: keep fresh dispatch default and require live token/cost evidence before any automatic reuse mode.

## Out of scope

- Implementing production reuse behavior in this stage.
- Changing existing Pi frontdoor/support code before direct resume semantics are proven in a harness.
- Reusing Claude/Codex team semantics blindly; Pi reuse must be Pi-native.
- Treating `pi-subagents` acceptance contracts as Spacedock validation. Spacedock uses independent validation, not same-agent acceptance finalization.
- Treating transcript text as proof without durable entity/file evidence.

## Acceptance criteria

**AC-1 - Pi feedback loops preserve independent validation semantics.**
Verified by: a future live Pi workflow test that starts with a validation rejection, dispatches the routed implementation follow-up without a `pi-subagents` acceptance contract, and verifies the resulting on-disk entity stage reports and commits rather than a child self-acceptance result.

**AC-2 - Fresh redispatch remains the default for Pi `feedback-to: implementation`.**
Verified by: Pi runtime instruction invariant tests plus a live rejection-flow smoke whose run ledger/session evidence shows a new implementation dispatch was created by default after validation feedback.

**AC-3 - Any opt-in resume path is epoch-safe and cannot accept stale completions.**
Verified by: registry unit tests that mark a worker completed, reactivate it for a follow-up epoch, reject epoch-0 completion evidence, and accept only matching epoch/run evidence after the follow-up writes durable output.

**AC-4 - Any opt-in resume path records durable, auditable worker handles.**
Verified by: a test or smoke fixture that writes a registry record containing substrate, worker label, run/session handle, entity slug, stage, state, completion epoch, and session file path, then reopens the registry and verifies the record before resume.

**AC-5 - Resume is not promoted to workflow-default without reproducible budget evidence.**
Verified by: a live resume smoke that records child-level elapsed time, token usage, and cost for resumed and fresh follow-ups in a machine-readable artifact; the implementation gate fails if telemetry is missing or only parent transcript prose is available.

## Test plan

- **Registry epoch tests:** extend `internal/piruntime` tests to cover active-again transitions, stale completion rejection, missing worker/session errors, entity/stage mismatch, and persistence across reopen.
- **Instruction invariant tests:** assert Pi first-officer and ensign runtime references keep fresh dispatch as default, forbid `pi-subagents` acceptance contracts for Spacedock stages, and require treating each follow-up assignment as a fresh cycle unless explicit manual/debug resume is requested.
- **Dispatch/ledger tests:** if a worker registry is added to production code, test that the first officer records only minimal durable metadata and that normal validation `feedback-to` routes do not consume old worker handles automatically.
- **Live resume smoke:** create a temp split-root workflow; run an initial Pi worker that writes a durable marker/report; resume or revive from the stored child handle/session with a feedback-style follow-up that writes a second marker/report; verify both markers on disk and registry epoch evidence. This smoke should fail loudly if direct `subagent({ action: "resume" })` is unavailable in the active non-interactive harness.
- **Fresh comparison smoke:** run the same follow-up through fresh dispatch and capture child-level elapsed/token/cost telemetry in a parseable artifact; use it as evidence only when child telemetry is isolated from parent orchestration/cache costs.
- **Lifecycle cleanup tests:** verify closed/terminal worker records are ignored for scheduling, missing session files fail safely, and cleanup does not affect normal fresh dispatch.

## Stage Report: ideation

- DONE: Read the entity and `docs/dev/README.md` workflow rules; kept frontmatter status as `ideation`.
- DONE: Inspected Pi runtime docs/code from `.worktrees/spacedock-ensign-pi-runtime-no-subagent-acceptance/` because the main checkout did not contain `skills/first-officer/references/pi-first-officer-runtime.md`, `skills/ensign/references/pi-ensign-runtime.md`, or `internal/piruntime/*`.
- DONE: Inspected installed `pi-subagents` README and skill file under `/Users/clkao/.pi/agent/npm/node_modules/pi-subagents/`.
- DONE: Confirmed local live setup had `pi` 0.78.0, `pi-subagents`, and `pi-intercom` installed; used only temp workflow/state paths for the live experiment.
- DONE: Ran a smallest live/manual experiment that wrote `initial-marker`, then a resume-equivalent revived/forked follow-up that wrote `resumed-feedback-marker`, then a fresh comparison follow-up that wrote `fresh-followup-marker`.
- DONE: Recorded available run/session/artifact evidence and elapsed observations: reported run id `bcf031b1-6f0d-41fc-ab64-3b128857b28e`, child session `2026-06-04T05-10-46-339Z_019e910a-8543-7d16-8ca8-64e119d5af28.jsonl`, resume-equivalent elapsed `9.896s`, fresh elapsed `24.617s`.
- DONE: Modeled stale completion/epoch safety, durable handle storage, context drift, budget tradeoffs, failure modes, cleanup/lifecycle, and workflow default vs manual/debug use.
- SKIPPED: Did not claim direct `subagent({ action: "resume" })` was proven; the non-interactive harness did not expose a clean visible `subagent` tool call, so the live evidence is classified as persisted-session revive/fork proof rather than full direct tool-resume proof.
- SKIPPED: Did not use `pi-subagents` acceptance contracts; Spacedock validation remains independent.
- SUMMARY: Pi persisted child context can be revived/forked and can produce durable second evidence, but direct tool-resume and child-isolated budget savings are not proven enough to make reuse a default workflow semantic. Keep fresh dispatch as default; reserve resume for manual/debug or later opt-in experiments guarded by registry epochs and live telemetry.
