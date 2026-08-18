---
id: j7jhntfa2ve8g6jwhatktrrv
title: Under a delegated conn, record the FO as approver and cite the grant
status: backlog
source: "Captain CL, 2026-08-18, reframing the live-lane inventory: the auto-continue journey tests the conn and the FO behaved correctly under it; the defect is the approver label, not the approval. Corroborated by the in-tree audit note at internal/ensigncycle/shared_live_runner_test.go:139 — 'finding 9 — approval-actor alternation under a delegated conn: recording person:captain for a decision no captain made in-session grades green.'"
started:
completed:
verdict:
score:
worktree:
issue:
---

When the FO resolves a gate under a delegated conn, the durable record should name the FO as approver and cite where the conn was granted. Today it can name the captain for a decision the captain never made in-session.

## Problem

A headless FO with a conn grant is authorized to resolve gates. The contract says so: the grant is a phrase quoted from the prompt, and with it the FO may "resolve gates per `## Completion and Gates` and drive to terminal."

So approving is correct behavior. What is wrong is the signature.

Observed twice on the claude-sonnet lane (runs 31996696789 and 32092321763 attempt 2): the FO ran `gate record <task> --decision approve --actor person:captain --consume`, durably closing the gate and attributing the decision to a captain who was not in the session. The FO then recited the exact conn-grant rule in its final message, which shows the prose was known and did not determine the record.

The system has no convention for who signs a delegated approval, so the FO invented one. The in-tree audit records the same gap from the other side: `internal/ensigncycle/shared_live_runner_test.go:139` notes that recording `person:captain` for a decision no captain made **grades green**.

The cost is an audit trail that misattributes authority. A reader of the state checkout cannot distinguish a decision the captain made from one the FO made under a grant, which is exactly the distinction the gate exists to preserve.

## Proposed approach

{Ideation fills this in. The shape the captain named: under a conn, the approver is the FO, and the record cites where the grant was given. That implies a durable representation of the grant — today it exists only as prose in a prompt. Decide the smallest form that `gate record` can check, and whether a missing citation refuses or warns.}

## Out of scope

Changing whether a conn authorizes gate resolution; it does. Changing the grant phrases. The `human-gate-bypassed` grader's own logic beyond what this attribution change requires.

## Expected surface and tolerance

Estimate net LOC change: +70 across 3 files. Report insertions and deletions separately. Do not declare a gross tolerance. Semantics changed: the actor recorded on a conn-delegated approval, and the grammar `gate record` accepts.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - A reader of durable state can tell a captain decision from an FO decision made under a grant.**
This is the measuring AC: over a recorded gate resolved under a conn, the count of approvals attributed to a human actor with no in-session human decision must be ZERO. Verified by resolving a gate in a headless run with a conn grant and reading the durable record: the approver is the FO and the grant is cited. Fails on today's behavior, which records `person:captain` and is indistinguishable from a real captain approval.

**AC-2 - The grader's expectation matches the corrected behavior.**
Verified by the auto-continue journey grading GREEN on an FO that approves under a conn with FO attribution and a citation, and RED on one that attributes to `person:captain` with no in-session captain — the inverse of today, where the audit note records that the false attribution grades green. Fails if the journey still treats a correct conn-delegated approval as a bypass.
