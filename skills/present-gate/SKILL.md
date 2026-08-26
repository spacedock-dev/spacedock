---
name: present-gate
description: "First-officer gate-presentation rendering — the captain-facing gate-review template and assembly rules, including workflow-owned finding labels. Invoke at the gate point after the FO has decided a stage must be presented."
user-invocable: false
---

# Present Gate

This skill carries the first-officer's captain-facing gate-presentation rendering: the gate-review format template and the assembly rules for filling it. The decide-to-gate and AC-cross-check policy stays always-on in the FO contract; this skill loads at the gate point to render the decision the FO has already made.

## Presentation channels

Stable v1 permits chat or Subspace to present the committed gate. Both channels return semantic decision and reason input to the First Officer. After the prepared room and binding commit, render the gate template exactly once in the selected channel before the next decision-mutation tool call. The presentation must name the entity and stage, a compact bound Briefing snapshot, one recommendation, and the decision ask; exact headings are not required. A short digest prefix can identify the snapshot in prose; the canonical Briefing and recorder validation retain the exact full digest authority. Prior narration, child/tool output, or a later summary does not complete presentation. Delegated authority does not waive presentation. Only then hand the semantic decision and reason to `${SPACEDOCK_BIN:-spacedock} gate record <entity> --decision approve|revise|hold --actor ID [--reason TEXT]`.

A qualifying review is semantic: wording may vary while those five facts remain explicit.
The presenter emits once; the lifecycle, prompt, and host adapter do not duplicate its review.

Subspace is a presentation interface, not a second recorder. It returns no Result or inventory for Spacedock to ingest and supplies no recorder authority or output paths. Provider-specific room-backed Result recording is explicitly outside v1; a failed or partial presentation never becomes an approval.

## Gate Presentation

Present gate reviews in this format:

```
Gate review: {entity title} — {stage}
Recommend {approve | reject: {one-line reason}}.
Reviewed snapshot: {bound Briefing identity and compact digest prefix}

{Render `Evidence for {stage}:` with only the rows selected by Gate content or the bounded fallback. Omit the heading when no decision-relevant evidence exists.}

{If reviewer findings exist, render them under `Reviewer findings` in the active workflow's declared category order and exact labels. Preserve each recorded category and omit empty categories. If the workflow declares none, use one neutral `Findings:` list. If no reviewer ran, omit this whole block.}

Decision: {one-line decision prompt naming what approval/rejection does in concrete terms — e.g., "approve to enter implementation in worktree `.worktrees/...`" or "reject to bounce back to {feedback-to target} with the authorized findings above"}.
```

### Captain-facing assembly rules

- Fetch the current stage with `${SPACEDOCK_BIN:-spacedock} dispatch show-stage-def --workflow-dir {workflow_dir} --stage {stage}` before selecting evidence. Never infer content from the stage name.
- Treat a `Gate content` instruction as the workflow's authoritative presentation preference and override. When it is absent, show a concise, meaningful, decision-relevant subset of the current stage definition, selected Artifact and References, current stage report, checklist and AC evidence, and findings. Do not fabricate facts or dump every source.
- Generic hints apply only inside that fallback: before a direction is selected, show what is proposed and what proof the decision needs; after selection, show the direction, risks, surface, and proposed proof; after execution, show actual results, checks, acceptance evidence, and readiness. The declared stage definition decides which hint fits.
- Omit missing evidence, empty result classes, empty finding categories, zero-result rows, and placeholders such as `None` or `N/A`. A negative summary such as `no material findings` means the finding group is empty; omit it. Do not print an aggregate count that names a zero class. Preserve the workflow's finding labels and order; presentation does not classify findings.
- Name the task and stage, Briefing identity and digest, one recommendation, and one concrete decision effect. Keep `${SPACEDOCK_BIN:-spacedock} gate record --decision` as the sole recorder; presentation adds no authority.
- Keep the decision visible and the review concise. Use the workflow's declared entity label in authored prose, name concrete reasons for rejection, and mention a worktree only when the decision changes worktree state.
