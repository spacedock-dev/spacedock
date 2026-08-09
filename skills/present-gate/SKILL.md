---
name: present-gate
description: "First-officer gate-presentation rendering — the captain-facing gate-review template and assembly rules, including workflow-owned finding labels. Invoke at the gate point after the FO has decided a stage must be presented."
user-invocable: false
---

# Present Gate

This skill carries the first-officer's captain-facing gate-presentation rendering: the gate-review format template and the assembly rules for filling it. The decide-to-gate and AC-cross-check policy stays always-on in the FO contract; this skill loads at the gate point to render the decision the FO has already made.

## Presentation channels

Stable v1 permits chat or Subspace to present the committed gate. Both channels return semantic decision and reason input to the First Officer. After the prepared room and binding commit, render the gate template exactly once in the selected channel before the next decision-mutation tool call. The presentation must name the entity and stage, a compact bound Briefing snapshot, one recommendation, and the decision ask; exact headings are not required. A short digest prefix may identify the snapshot in prose; `request.json`, the canonical Briefing, and recorder validation retain the exact full digest authority. Prior narration, child/tool output, or a later summary does not complete presentation. Delegated authority does not waive presentation. Only then hand the semantic decision and reason to `${SPACEDOCK_BIN:-spacedock} gate record <entity> --decision approve|revise|hold --actor ID [--reason TEXT]`.

A qualifying review is semantic: wording may vary while those five facts remain explicit.
The presenter emits once; the lifecycle, prompt, and host adapter do not duplicate its review.

Subspace is a presentation interface, not a second recorder. It returns no Result or inventory for Spacedock to ingest and supplies no recorder authority or output paths. Provider-specific room-backed Result recording is explicitly outside v1; a failed or partial presentation never becomes an approval.

## Gate Presentation

Present gate reviews in this format:

```
Gate review: {entity title} — {stage}
Recommend {approve | reject: {one-line reason}}.
Reviewed snapshot: {bound Briefing identity and compact digest prefix}

Stage evidence ({stage}):
{For backlog: Outcome: {seed outcome}. Scope boundary: {included work and explicit cuts}. Proof readiness: {proposed observable proof and known unknowns}.}
{For ideation: Chosen direction: {selected approach}. Risk evidence: {riskiest-mechanism spike result or recorded no-spike basis}. Expected surface: {files, LOC, tolerance, and semantic-change declaration}. Acceptance proofs: {each AC and its proposed observable proof}.}
{For validation, cite ## Stage Report in {entity_file_path} lines {start}-{end}, then render only non-empty classes:
- DONE: {≤10-word gist of item}
- SKIPPED: {gist} — {one-line reason}
- FAILED: {gist} — {one-line reason}
Assessment: {N} done, {N} skipped, {N} failed.
Checks executed: {checks actually run and their outcomes}. Acceptance evidence: {actual evidence by AC}. Delivery readiness: {ready or concrete blocker}.}

{If reviewer findings exist, render them under `Reviewer findings` in the active workflow's declared category order and exact labels. Preserve each recorded category and omit empty categories. If the workflow declares none, use one neutral `Findings:` list. If no reviewer ran, omit this whole block.}

Decision: {one-line decision prompt naming what approval/rejection does in concrete terms — e.g., "approve to enter implementation in worktree `.worktrees/...`" or "reject to bounce back to {feedback-to target} with the authorized findings above"}.
```

### Captain-facing assembly rules

The template is the floor, not the ceiling. The FO MUST hold to the following discipline when filling it:

1. **Lede first, decision last, nothing between them buried.** The title, recommendation, reviewed snapshot, and final decision line are the common spine. Everything else is supporting evidence; keep the vote visible without repeating the recommendation.
2. **Select evidence by gate stage.** Backlog shows the seed outcome, scope boundary, and proof readiness; it does not claim a Stage Report, executed check, or delivery. Ideation shows the chosen direction, risk evidence, expected surface, semantic-change declaration, and each AC's proposed observable proof; it does not invent execution results. Validation shows actual report results, executed checks, AC evidence, and delivery readiness. For ideation and validation, state the selected direction (the approach or PASS/REJECTED) in the stage evidence rather than making the captain infer it.
3. **Omit unavailable and empty evidence.** Backlog and ideation have no validation Assessment or result classes. At validation, keep the numeric `Assessment:` totals, but render a DONE, SKIPPED, or FAILED heading and its bullets only when that class has at least one item. Never substitute `None`, `N/A`, a zero-result row, or an empty findings block. Cite the Stage Report and paraphrase each rendered item as a verb-noun gist (≤10 words, no new facts); append the one-line reason to SKIPPED/FAILED. If a decision-relevant reviewer finding directly questions an item's evidence, inline that evidence paragraph under the finding. Otherwise no Stage Report content appears.
4. **Reviewer findings preserve workflow classification.** Render findings in the active workflow's declared category order and exact labels, preserving each recorded category and dropping empty categories. If the workflow declares none, use a neutral `Findings:` list. Presentation never classifies.
5. **Recommendation appears exactly once.** The `Recommend {approve | reject: {reason}}` line is the only place the FO states its verdict. Do not duplicate it elsewhere or re-explain it in an enumerated list.
6. **Bounce-back recommendations name the concrete asks.** If recommending reject, the reason line names the specific concerns by content, not by reference. Bad: "address the reviewer's five concrete notes." Good: "tighten AC-2 substring assertion; correct the file X claim; cut the format-pedantry aside."
7. **No format-pedantry asides.** Format drift (`1./2./3./4.` instead of `**AC-N**`, missing trailing period) is not load-bearing for a gate decision. Surface it only when the active workflow's policy makes it gate-blocking, under that policy's exact category label.
8. **One sentence of worktree heads-up when approval changes worktree state.** When approving opens or closes a worktree, the Decision line names it: "approve to enter implementation in worktree `.worktrees/{worker_key}-{slug}`". One sentence, not a section.
9. **Target length: 15-25 lines of FO-authored prose.** The full gate message should fit in 15-25 lines. If it exceeds 25, the FO is over-narrating; cut.
10. **FO-authored prose speaks the workflow's declared label.** Where the gate-summary prose the FO writes — stage evidence or the `Decision:` line — names the kind of thing under review, use the workflow's declared `entity-label` / `entity-label-plural` from `«state.boot»()`, not the generic "entity". A `ticket` workflow's Decision line says "approve to enter implementation on this ticket"; an `experiment` workflow says "experiment". The `{entity title}` placeholder and structural headings stay generic — only the FO-authored noun localizes.
11. **Surface verification state as evidence, not as a label.** When the gate turns on checks that ran outside this presentation (CI lanes, a validation report), hold them to the shared core's self-evidence bar (`## Working Principles`): state which relevant checks actually ran and passed, and read any failure from this run's evidence (the failing test/assertion), never from an inherited "known flake" label. The captain votes on which checks are green and why a red is red.
