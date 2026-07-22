---
name: present-gate
description: "First-officer gate-presentation rendering — the captain-facing gate-review template plus the nine assembly rules for filling it (lede-first/decision-last spine, chosen-direction prose, report citation, reviewer-finding tiers, single recommendation, concrete bounce-back asks, length budget). Invoke at the gate point after the FO has decided a stage must be presented."
user-invocable: false
---

# Present Gate

This skill carries the first-officer's captain-facing gate-presentation rendering: the gate-review format template and the assembly rules for filling it. The decide-to-gate and AC-cross-check policy stays always-on in the FO contract; this skill loads at the gate point to render the decision the FO has already made.

## Presentation channels

Chat is the default channel. Render the gate template in the current conversation, then hand the captain's semantic decision to `${SPACEDOCK_BIN:-spacedock} gate record ... --decision`. An override changes only where the captain sees the review; gate policy and recorder ownership stay unchanged.

A workflow or session may declare one presentation override. Apply this contract:

1. **Probe before side effects.** Run the override's read-only availability and version probe before creating a package or retention path. If the presenter is missing or mismatched, launch nothing, create nothing, mutate nothing, and emit one line naming both the install or upgrade remedy and the chat fallback. Then use chat.
2. **Present the complete canonical Briefing.** Give the override the retained `briefing.json`, its exact question, and every Artifact and supporting Reference at the revision the manifest names. Derive the presentation title from that Briefing. A one-file float is advisory evidence only; the manifest's mention of unpresented artifacts never turns it into a complete-package review.
3. **Block through retention.** Keep the gate-attempt ensign addressable for the whole invocation. Pane or session creation and a `wait_agent` timeout are not completion. Re-wait after a timeout; complete only after the presenter exits, the Result validates, and retention finishes.
4. **Retain from the first byte.** Before launch, the caller chooses a room-resident result path outside launcher-owned scratch. The override writes the exact Result there and retains the review log and diagnostics, including argv, exit code, and stderr. It never deletes the retention directory on success, hold, validation failure, child failure, or launcher death.
5. **Associate only what the reviewer saw.** Build `spacedock-result-association` from the artifacts actually presented, not from the manifest alone. Hand the exact retained Result and association to `${SPACEDOCK_BIN:-spacedock} gate record ... --result FILE --association FILE --actor ID`. The recorder verifies the frozen Briefing digest, complete presentation mapping, actor, and exact Result before it normalizes the provider Briefing id.
6. **Retain failures without inventing a decision.** After launch, a missing or invalid Result leaves the gate open with its package and diagnostics recoverable. Do not replace that attempted review with a chat decision or write entity frontmatter directly.

The hardened override script and its committed CI drive suite belong to the provider. A Subspace override is release-eligible only at the pinned provider revision that carries both; Spacedock's binary and this skill contain no Subspace transport implementation.

## Gate Presentation

Present gate reviews in this format:

```
Gate review: {entity title} — {stage}
Chosen direction: {one-line summary of the ensign's chosen approach, or `n/a` for stages without a chosen-direction concept (e.g., simple work stages, merge)}
Recommend {approve | reject: {one-line reason}}.

Checklist (from ## Stage Report in {entity_file_path} lines {start}-{end}):
- DONE: {≤10-word gist of item}
- SKIPPED: {gist} — {one-line reason}
- FAILED: {gist} — {one-line reason}

{If reviewer findings exist, render them under a `Reviewer findings` heading in two tiers — `Material:` (fact-corrections, contract violations, missing AC evidence, broken claims) and `Polish:` (wording, format drift, non-blocking suggestions). Drop the tier entirely if it has no items. If no reviewer ran, omit this whole block.}

Assessment: {N} done, {N} skipped, {N} failed.

Decision: {one-line decision prompt naming what approval/rejection does in concrete terms — e.g., "approve to enter implementation in worktree `.worktrees/...`" or "reject to bounce back to {feedback-to target} with the material findings above"}.
```

### Captain-facing assembly rules

The template is the floor, not the ceiling. The FO MUST hold to the following discipline when filling it:

1. **Lede first, decision last, nothing between them buried.** The first three lines (title, chosen direction, recommend) and the final line (decision) are the spine. Everything else is supporting evidence; if the captain stops reading after line three, they can still vote.
2. **Chosen direction is required as FO prose.** When the stage selected among options (ideation picks an approach, validation picks PASS/REJECTED), name it on the `Chosen direction:` line; don't make the captain infer from the Checklist gist or open the entity file. For stages without a chosen direction, use `n/a`.
3. **Cite the Stage Report; render a one-line gist roll-up.** Do not paste it into the gate message. Under `Checklist:`, render one bullet per DONE/SKIPPED/FAILED item as a verb-noun gist (≤10 words, FO paraphrase, no new facts). For SKIPPED/FAILED, append `— {one-line reason}`. Cite the full report by file path and line range. If a reviewer Material finding directly questions a checklist item's evidence, inline that item's evidence paragraph under the finding so the captain can decide without opening the file. Otherwise no Stage Report content appears.
4. **Reviewer findings render in priority tiers.** Group into `Material:` (fact-corrections, contract violations, missing AC evidence, claims contradicted by the codebase) and `Polish:` (wording, format drift, non-blocking suggestions). Drop empty tiers. Do not flat-bullet material next to polish.
5. **Recommendation appears exactly once.** The `Recommend {approve | reject: {reason}}` line is the only place the FO states its verdict. Do not duplicate it elsewhere or re-explain it in an enumerated list.
6. **Bounce-back recommendations name the concrete asks.** If recommending reject, the reason line names the specific concerns by content, not by reference. Bad: "address the reviewer's five concrete notes." Good: "tighten AC-2 substring assertion; correct the file X claim; cut the format-pedantry aside."
7. **No format-pedantry asides.** Format drift (`1./2./3./4.` instead of `**AC-N**`, missing trailing period) is not load-bearing for a gate decision. Surface only if it blocks the gate; if it does, it is a Material finding, not a separate paragraph.
8. **One sentence of worktree heads-up when approval changes worktree state.** When approving opens or closes a worktree, the Decision line names it: "approve to enter implementation in worktree `.worktrees/{worker_key}-{slug}`". One sentence, not a section.
9. **Target length: 15-25 lines of FO-authored prose.** The full gate message should fit in 15-25 lines. If it exceeds 25, the FO is over-narrating; cut.
10. **FO-authored prose speaks the workflow's declared label.** Where the gate-summary prose the FO writes — the `Chosen direction:` line, the `Checklist:` gist roll-up, the `Decision:` line — names the kind of thing under review, use the workflow's declared `entity-label` / `entity-label-plural` from `«state.boot»()`, not the generic "entity". A `ticket` workflow's Decision line says "approve to enter implementation on this ticket"; an `experiment` workflow says "experiment". The `{entity title}` placeholder and the structural headings (`Gate review:`, `Checklist:`, `Decision:`) stay generic — only the FO-authored noun localizes.
11. **Surface verification state as evidence, not as a label.** When the gate turns on checks that ran outside this presentation (CI lanes, a validation report), hold them to the shared core's self-evidence bar (`## Working Principles`): state which relevant checks actually ran and passed, and read any failure from this run's evidence (the failing test/assertion), never from an inherited "known flake" label. The captain votes on which checks are green and why a red is red.
