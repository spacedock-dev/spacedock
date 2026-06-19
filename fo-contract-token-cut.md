---
title: Recover ~638 boot-resident FO-contract tokens (micro-test-verified cuts)
status: backlog
source: The 2026-06-15 fo-contract-token-cleanup proposal classified ~638 boot-resident tokens (safe-cut/cut-with-care/keep) by adversarial reasoning but filed no task to APPLY them; wg/95b were prior sweeps, now archived/done. This task owns the actual recovery, gated by the no-guidance-control micro-test (superpowers v6) that converts the reasoned verdicts into empirical ones. CL direction 2026-06-16.
score: 0.5
sprint: 0205-layered-fo
sprint-readiness: defer
issue:
id: y2r7ew51xqs6q3avsb6mcaka
group: cleanup
---

Apply the ~638-token cut list in `docs/dev/_proposals/fo-contract-token-cleanup.md` to the shipped FO-contract files — each cut empirically confirmed (not just reasoned) by a no-guidance-control micro-test. The deliverable is the trimmed contract files; the micro-test is the method that makes the cut safe.

## Problem

The proposal classifies ~638 boot-resident tokens across four files (shared-core, the Claude runtime adapter, using-claude-team, present-gate) as safe-cut / cut-with-care / keep — every verdict reached by ADVERSARIAL REASONING ("construct a break attempt; would the FO misbehave?"). That is a strong prior but docs-confidence, not measurement: nobody has run an FO with and without a clause and observed the difference. So the cuts sit unapplied, and the 13 conservative "keep" verdicts may be over-cautious. The boot-resident contract is re-carried every turn of every FO session (~13.4K tokens of permanent occupancy, measured this session), so the recovery compounds across the whole session.

## Proposed approach

Apply the cuts, gated by the empirical layer the proposal now specifies (its "Cut criteria — no-guidance control" section):

1. **Validating-first-step (pay the small bill first).** Prove the no-guidance-control micro-test discriminates, on three candidates spanning the verdict space — SC-13 (a `keep`), SC-5 (`cut-with-care`), SC-3 (`safe-cut`). Reuse the bare-`claude` launch + durable-state grade the haiku-loop-spike (w4) proved (AC-1). For a candidate clause: sample FO behavior N>=5 on the smallest realistic exercise that exposes the clause's job, in two arms — WITH the clause and WITHOUT it (the control) — in the real surrounding contract; read behavior by hand (a string match is not a behavior); variance-as-signal. ~$0.15-0.30/sample.
2. **Apply across the full list.** Control already satisfies the behavior -> confirmed dead weight (delete, do not merely compress). Absence changes behavior across N -> confirmed load-bearing (keep the load-bearing fragment, cut only the inert remainder). RE-TEST the 13 keeps — the control can overturn a `keep` whose feared misbehavior does not actually occur.
3. **Ship through the dispatched-worker path.** The contract files are shipped scaffolding (not FO-direct-editable), so the trimmed files land via a dispatched worker in a worktree, with the four live shared scenarios (gate-guardrail, rejection-flow, feedback-3-cycle-escalation, merge-hook-guardrail) green on Claude AND Codex as the behavior-preservation oracle, plus a `wc` size-delta and a detached adversarial audit of the diff (per the proposal's Closing note).

## Deliverable

The trimmed FO-contract files committed, with: a measured token-delta (target ~638, actual reported), the four live scenarios green on both hosts, and a per-cut record of the micro-test verdict (applied / kept / keep-overturned). The no-guidance-control method, proven here as a byproduct, is reused later by 0205's prose-function restructure — but THIS task ships the cut, not a method.

## Out of scope

- A general eval framework (the micro-test is a thin per-clause check, not a product — YAGNI until a future batch justifies one).
- 0205's prose-function restructure (it reuses the method; different task).
- The judgment-audit generalization over the full routing table (rides w4 / 0205).
- Net-new cut hunting beyond the proposal's verified list (re-testing the 13 keeps IS in scope; finding fresh cuts is not).

Source: the 2026-06-15 token-cleanup proposal + superpowers v6 (writing-skills "Micro-Test Wording Before Full Scenarios", positive-instruction-redesign-design, strict-cost-sdd-design).
