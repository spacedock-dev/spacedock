---
title: 'No-guidance-control micro-test harness for FO-contract prose cuts + restructures'
status: backlog
source: 'superpowers v6 prose-efficacy lessons + CL direction (2026-06-16). The token-cleanup proposal and the binary-simplification dehydration lever (Phase 0.A, wg) verify prose cuts by adversarial REASONING only — docs-confidence, never empirical; nobody has run an FO with vs without a clause and observed the difference. The 0205 prose-function restructure will rewrite the highest-residency boot-resident prose with no cheap verifier. This is the missing empirical layer. Shares the bare-claude launch + durable-state grade harness with w4 haiku-loop-spike.'
score: 0.5
sprint: 0205-layered-fo
sprint-readiness: defer
issue:
id: y2r7ew51xqs6q3avsb6mcaka
---

The missing empirical layer over the FO-contract prose-reduction work.

## Problem

Cut decisions today rest on adversarial REASONING, not measurement. `fo-contract-token-cleanup-2026-06-15.md` classifies ~638 boot-resident tokens as safe-cut / cut-with-care / keep, each verified by "construct a break attempt and judge whether the FO would misbehave" — a strong prior, but docs-confidence, not measurement. The `binary-simplification-roadmap` dehydration lever (Phase 0.A, wg) has the same gap, and 0205's prose-function restructure rewrites the most re-read prose in the system with no cheap verifier.

## Approach (thin prototype — prove the method, do NOT build a framework)

A no-guidance-control micro-test, reusing the proven internal/ensigncycle bare-`claude --model haiku` launch + durable-state grade that haiku-loop-spike (w4) proved (AC-1). For a candidate clause: sample the FO's behavior N>=5 on the smallest realistic exercise that exposes the clause's job, in two arms — WITH the clause and WITHOUT it (the control) — in the real surrounding contract; read behavior by hand (a string match is not a behavior); variance-as-signal. Decision rule, both directions: control already satisfies -> confirmed dead weight (delete, do not merely compress); absence changes behavior across N -> confirmed load-bearing (keep). Prove it discriminates on three candidates spanning the verdict space: SC-13 (a keep — "keep dispatching ready entities when one blocks"), SC-5 (cut-with-care — dispatch deferred-module to a pointer), SC-3 (safe-cut — event-loop status parenthetical). Then stop.

## Why it matters

Converts the proposal's strongest layer (reasoned classification) into its missing layer (empirical evidence); can promote a safe-cut OR overturn a keep. Same harness 0205's spike needs. Sits upstream of the expensive 4-live-scenario behavior-preservation oracle (~$0.15-0.30/sample vs ~$12/run).

## Out of scope

- Applying the 638-token cuts (that is the token-cleanup proposal, post-method, through the dispatched-worker path).
- A general eval framework (YAGNI until the cut batch justifies it).
- The judgment-audit generalization (Haiku-vs-opus baseline over all routing-table judgment points) — that rides w4 / 0205, not this.

Source: superpowers v6 writing-skills "Micro-Test Wording Before Full Scenarios", positive-instruction-redesign-design, strict-cost-sdd-design.
