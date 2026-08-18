---
id: 9gaphcqfw3t5ayknr7fpvkjd
title: Make the rejection-worker-topology red diagnosable, then judge the conduct
status: backlog
source: "Captain CL, 2026-08-18, in chat (\"reconfirm and file\"), after the 0.27 readout named this the only live-lane red with no owner. Failing journey: rejection-worker-topology on codex-live. Prior disposition: live-lane-gap-inventory.md finding 7 said \"new entity\"; none of ca9/6ht/j7j/vka covers it."
started:
completed:
verdict:
score:
worktree:
issue:
---

`rejection-worker-topology` is red on codex-live and cannot be diagnosed, because the evidence that would explain it is deleted before anyone can read it.

The grader checks the shape of the workers the rejection journey produces — how many, of what kind, in what relationship. When it reds, the reason lives in the codex native rollout. That rollout is written under an isolated `CODEX_HOME` which `t.Cleanup` removes when the test ends. So the failure is repeatable and its cause is unrecoverable.

This is the last of the four failing live journeys with no owner. `ca9` and `6ht` fix two graders that redden compliant agents. `vka` closes a real ceremony gap. Nothing touches this one.

## Why this is one entity and not two

The obvious split — "persist the evidence" then "fix the conduct" — does not survive contact, because nobody can say yet whether there IS a conduct defect. The open question recorded in the inventory is whether a repeated `followup_task` to the same live worker within one round is a violation or is fine. That question is unanswerable while the payload is deleted.

So the deliverable is: make the failure observable, read it, and then either fix the product or correct the grader. The first half is worthless alone — persisted artifacts nobody reads change nothing — and the second half is unreachable without it.

## Ordering: hardest thing first

Do NOT begin by designing a fix for conduct that has not been observed. Begin by persisting the rollout and followup payloads into run artifacts, re-running the journey, and reading what actually happened. The diagnosis decides the rest of the work, and it may show the grader is wrong rather than the agent.

## Out of scope

The other three live-lane reds, each owned by an entity already in flight. Changing what the rejection journey tests. Any broader artifact-retention scheme for other journeys — persist what this diagnosis needs and no more.

## Ideation must settle

1. Whether the payload can be persisted without changing what the journey exercises. Copying out of an isolated `CODEX_HOME` before cleanup must not alter the run it is observing.
2. What the retained evidence must contain to answer the actual question — repeated `followup_task` to one live worker in a single round: violation, or sanctioned?
3. Whether this is a product defect or a grader defect. State the answer as a finding, not a preference. "The grader is wrong" is a legitimate and complete outcome here, exactly as it was for `6ht`.
