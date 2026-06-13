# Gates & decisions

A gate is the decision point at the end of a stage where nothing advances without your vote. When a stage declares a gate, the first officer stops after the worker completes, presents a review, and waits for you. It never self-approves.

Each call you make sharpens the bar, and the destination is delegation. When you are sufficiently confident in the workflow and the bar you set, hand over the conn and let the first officer drive multiple tasks with auto-approval:

```
you have the conn to drive toward your sprint goal, authorized to approve and
merge PR on CI green. use your judgement.
```

Until then, the gates are yours.

## What you see at a gate

A gate review has a fixed spine: the first three lines and the last line carry the decision; everything between is supporting evidence. If you stop reading after line three, you can still vote.

```
Gate review: Fix the flaky login test — review
Chosen direction: replace sleep-based waits with event polling
Recommend reject: the AC-2 retry scenario has no covering test.

Checklist (from ## Stage Report in docs/ship-features/fix-the-flaky-login-test.md):
- DONE: login test stable across 50 consecutive runs
- FAILED: retry scenario unproven — no test exercises it

Reviewer findings
  Material: AC-2 cites a test file that does not exist
  Polish:   stage report wording drifts from the template

Assessment: 1 done, 0 skipped, 1 failed.

Decision: approve to close; reject to bounce back to implementation.
```

Material findings are the ones that should move your vote; Polish never blocks. The Decision line tells you concretely what your vote does. Every acceptance criterion is cross-checked before the review reaches you; a criterion without cited evidence is named rather than passed over.

## The three calls

- **Approve.** The work advances to the next stage. Approving the terminal stage merges and closes it.
- **Redo with feedback.** You accept the direction but send concrete fixes back. Name the specific asks ("tighten the AC-2 substring assertion, correct the file path claim"), not "address the reviewer's notes".
- **Reject.** The work bounces back to the stage that owns the fix, carrying your findings.

Redo and reject differ only in whether you accept the direction; both carry your concrete asks so the next worker has something to act on. Nothing closes without its verdict on the record.

## When work is rejected

Rejections bounce automatically: the findings route back, the work is redone, and the reviewer re-runs, with no stop at your desk. The gate reaches you only when the work passes review, or after **three failed rounds**, when the call returns to a human instead of bouncing again. Every round is on the record in the item's file.

A useful rejection to type at a gate: "send it back unless this now needs reframing".

## The detached adversarial audit

For high-stakes surfaces, a passing validation is necessary but not sufficient. Before merging, the first officer also runs a read-only adversarial audit. The audit catches the hole validation cannot see on its own: a test that passes today but would also pass on a broken future edit.

A workflow names its own high-stakes surfaces. Routine, low-blast-radius changes do not need an audit; a normal validation suffices.

The audit is read-only and cannot touch the deliverable. It tries to refute the validation: it constructs an adversarial edit that the deliverable's own tests should catch and confirms they do. A test that stays green under an edit that breaks the claim is a hole. "Refuted nothing material" is a valid recorded outcome.

Material findings route back through the normal feedback flow, and the gate is not presented as clean until they are closed. A clean audit is noted in the review.

## Where to go next

- [A worked example](worked-example.md) traces one real entity through every gate to a recorded verdict.
- [Operating a workflow](../running-workflows/operating.md) covers answering gates in the day-to-day loop.
