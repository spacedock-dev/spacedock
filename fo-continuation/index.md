---
title: Keep the FO running across handoffs and status questions
status: ideation
source: Captain-approved continuation correction, 2026-09-15
started:
completed:
verdict:
score: 0.95
worktree:
issue: spacedock-dev/spacedock#735
pr:
mod-block:
id: ad242krer2tckgx7150cyb7d
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

## Expected surface and tolerance

Ideation must propose a small net-line/file estimate after reading the current owners and existing behavioral proof. Reuse existing tests and at most one focused behavior comparison; explain any additional mechanism before adding it.

## Acceptance criteria

**AC-1 — A status question does not terminate an authorized async workflow.**
Verified by: a supported runtime exercise starts a worker, receives a captain status question, answers, reinstalls monitoring in the same turn, observes completion and dispatches the ready successor without another captain nudge. Oracle is observed runtime/tool sequence and durable state, not source wording. Removing continuation must fail.

**AC-2 — Dispatch and correction handoffs continue to the declared stopping condition.**
Verified by: existing dispatch/rejection journey or the same focused exercise shows monitoring after spawn and revise dispatch, then completion verification and next action. A final answer immediately after the handoff fails.

**AC-3 — Explicit stop and captain-decision boundaries remain effective.**
Verified by: paired pause/stop and unresolved approval cases halt without unauthorized dispatch or approval. The same continuation rule must not override captain control.

## Test plan

Read the current contracts and issue735, map claims to existing journey owners, and spike the smallest supported before/after behavior comparison. Do not build a transcript simulator that merely encodes the desired answer or a prose-grep test. Reuse existing skill smoke tests before changing command text. Required normal/race/format checks and independent validation follow implementation; do not start broad baselines during ideation. No expensive CI until local validation and stack readiness.

This task will land above the scheduling PR. No code push or CI until the stack is ready and authorized.

### Feedback Cycles

