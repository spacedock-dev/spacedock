---
name: present-gate
description: "First-officer gate-presentation rendering — the captain-facing gate-review template plus the nine assembly rules for filling it (lede-first/decision-last spine, chosen-direction prose, report citation, reviewer-finding tiers, single recommendation, concrete bounce-back asks, length budget). Invoke at the gate point after the FO has decided a stage must be presented."
user-invocable: false
---

# Present Gate

This skill carries the first-officer's captain-facing gate-presentation rendering: the gate-review format template and the assembly rules for filling it. The decide-to-gate and AC-cross-check policy stays always-on in the FO contract; this skill loads at the gate point to render the decision the FO has already made.

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

## Emit the gate to Bridge (host-neutral)

After rendering the gate-review to the captain in-session, the FO ALSO pushes the same gate to Bridge so a remote captain can decide it from the command-center UI. This is host-neutral — it lives here in `present-gate` (loaded by every host), never in a Claude-only hook. It is a **direct file append**, not a verb: add one line to `_bridge/fo-initiate.jsonl` (relative to the fleet root), using the `docs/seam-contract.md` §2.7 shape verbatim:

```json
{"schema":1,"id":"gate-<entity>-<stage>","ts":"<now RFC3339>","kind":"gate-review","workflow":"<slug>","entity":"<entity-slug>","ship_id":"<slug>/<entity-slug>","host":"<claude|codex|pi>","session_id":"<this session's id>","headline":"<the gate lede — the one-liner the captain sees>","body":"<optional supporting prose>","status":"open"}
```

- `schema` MUST be `1` and `id` non-empty and `kind` `gate-review` — Bridge drops the line otherwise.
- `id` is the fold key AND (defaulting `request_id` to it) the loop-closure correlator. It MUST be a STABLE function of `(entity, stage)` — e.g. `gate-<entity>-<stage>` — so re-emitting the same gate on each drain tick folds to ONE card instead of stacking duplicates. Never derive it from a timestamp or random value.
- `status` is ALWAYS `open`; Bridge overlays the resolution itself from the matching inbox `decision` intent (`approve`→approved, `reject`→rejected) by `request_id`. Do not write a resolved status.
- `host`/`session_id` carry attribution so the card shows which host raised the gate; pass the current host and session id.
- `headline` ≤ 240 chars (the gate lede); keep `body` ≤ 2000.

### Channel boundary — a gate goes to fo-initiate ONLY

A gate-review is a decidable interrupt. It goes to `_bridge/fo-initiate.jsonl` and NOWHERE ELSE:

- NEVER `fo-feed.jsonl` — that stream is ambient git narration (dispatch/advance), not a decidable ask; a gate rendered there has no Approve/Reject affordance.
- NEVER `fo-replies.jsonl` — that stream requires an `in_reply_to` correlator to a captain intent; an FO-initiated gate has no such parent and would be silently dropped.

Approve/Reject on the fo-initiate card close the loop back through the inbox by `request_id` (a `decision` intent the FO drains per the `bridge-seam` mod); that is why the id must be stable and shared with the gate the captain sees.
