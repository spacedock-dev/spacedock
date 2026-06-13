# Gates & decisions

A gate is the decision point at the end of a stage where nothing advances without your vote. When a stage is marked `gate: true` in the workflow README, the first officer stops after the worker completes, renders a gate review, and waits for you. It never self-approves.

## What a gate carries

A gate review has a fixed spine. The first three lines and the last line carry the decision; everything between is supporting evidence. If you stop reading after line three, you can still vote.

A gate review looks like this:

```
Gate review: {entity title} — {stage}
Chosen direction: {one-line summary of the worker's approach, or n/a}
Recommend {approve | reject: {one-line reason}}.

Checklist (from ## Stage Report in {entity_file_path} lines {start}-{end}):
- DONE: {≤10-word gist}
- SKIPPED: {gist} — {reason}
- FAILED: {gist} — {reason}

Reviewer findings
  Material: {fact-corrections, contract violations, missing AC evidence}
  Polish:   {wording, format drift, non-blocking suggestions}

Assessment: {N} done, {N} skipped, {N} failed.

Decision: {what approval/rejection does in concrete terms}.
```

Reading the fields:

- **`Chosen direction` names what the worker picked**, so you do not have to open the entity file to learn it. Ideation picks an approach; validation picks PASS or REJECTED. Stages with no choice to make show `n/a`.
- **`Recommend` is the first officer's verdict, stated exactly once.**
- **`Checklist` is a gist roll-up, not the report.** The full `## Stage Report` is cited by file path and line range; open it when you want the detail.
- **Reviewer findings split into `Material` and `Polish`.** Material items (fact-corrections, contract violations, missing acceptance-criterion evidence, claims the codebase contradicts) are the ones that should move your vote. Polish is non-blocking. An empty tier is dropped.
- **`Decision` names what your vote does in concrete terms.** For example, "approve to enter implementation" or "reject to bounce back to design".

Every acceptance criterion is also cross-checked before the review reaches you; a criterion left without cited evidence is named in the review rather than passed over.

## The three calls

You answer a gate with one of three calls.

- **Approve.** The entity advances to the next stage and the first officer dispatches it. If the next stage opens or closes a worktree, the `Decision` line told you so. Approving the terminal stage runs the merge and cleanup ceremony.
- **Redo with feedback.** You approve the direction but send concrete fixes back. Name the specific asks ("tighten the AC-2 substring assertion, correct the file path claim"), not "address the reviewer's notes". The first officer routes your asks back to the stage that owns the work, the worker re-does it, and the gate is re-presented.
- **Reject.** At a stage with a `feedback-to` target, rejecting bounces the work back to that target stage to be fixed and re-validated; this is the feedback cycle below. At a stage without `feedback-to`, rejection is terminal for that path.

A redo and a reject at a `feedback-to` stage run the same routing machinery. The difference is whether you are correcting a direction you accept or sending it back. Both name the concrete fix asks so the next worker has something to act on.

**A terminal close needs a recorded verdict.** Approving the terminal stage finalizes the entity, and Spacedock refuses to finalize or archive a terminal entity that carries no `verdict`: the outcome must be on the record. See the [frontmatter contract](../reference/frontmatter-contract.md#entity-fields) for the mechanics.

## Feedback cycles and the loop cap

When a feedback stage recommends `REJECTED`, or you reject at a `feedback-to` stage, the work routes back to the stage named in `feedback-to`: the stage that owns the fix, not the reviewer that flagged it. In the dev workflow, `validation` has `feedback-to: implementation`, so a rejected validation sends the deliverable back to implementation, not back to the validator.

The first officer tracks each round in a `### Feedback Cycles` section in the entity body. Your concrete findings route to the `feedback-to` target, the reviewer re-runs after the fix, and the gate comes back to you as a fresh review.

**The loop caps at three.** On cycle 3 the first officer escalates to you instead of bouncing a fourth time: the same fix has failed twice, so the call returns to a human.

## The detached adversarial audit

For high-stakes surfaces, a passing validation is necessary but not sufficient. Before merging, the first officer also runs a read-only adversarial audit. The audit catches the hole that validation cannot see itself: a test that passes today but would also pass on a broken future edit.

A workflow names its own high-stakes surfaces; routine, low-blast-radius changes do not need an audit, a normal validation suffices.

The audit is read-only and cannot touch the deliverable. It tries to refute the validation: it constructs an adversarial edit that the deliverable's own tests should catch and confirms they do. A test that stays green under an edit that breaks the claim is a hole. "Refuted nothing material" is a valid recorded outcome.

Material findings route back through the normal feedback flow, and the gate is not presented as clean until they are closed. A clean audit is noted in the review.

## Where to go next

- [A worked example](worked-example.md) traces one real entity through every gate to a recorded verdict.
- [Operating a workflow](../running-workflows/operating.md) covers answering gates in the day-to-day loop.
