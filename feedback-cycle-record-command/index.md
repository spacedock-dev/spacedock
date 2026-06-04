---
id: bwr6j6edkmfx5sbz73cr2952
title: spacedock status --record-feedback-cycle — binary-owned feedback-cycle count + 3-cycle escalation guard
status: backlog
source: "captain (2026-06-04) — forked from xa (feedback-guarantee-binary-gate) per the roadmap-the-decision + separate-build-task call. xa's ideation determined Candidate 1 (3-cycle escalation) is mechanizable via a dedicated cycle-record command (a spike disproved a --set status guard) and Candidate 2 (budget-probe) is not. This task SHIPS the Candidate-1 guard; xa closed as a roadmap decision."
score: "0.30"
started:
completed:
verdict:
worktree:
issue:
---

Promote the feedback-rejection 3-cycle-escalation guarantee from FO contract prose into a binary-enforced gate over durable on-disk state. The FO currently *tracks* feedback cycles in the `### Feedback Cycles` body section and is *instructed in prose* to escalate to the human on the 3rd rejection instead of auto-bouncing a 4th time — a prose-only guarantee whose ceiling is "the wording is present" and whose drift mode is an infinite reject→re-implement→reject loop burning tokens. This task makes the count binary-owned and tamper-evident, the same prose→binary promotion `mod-block` already models.

This is the **guarantee → code-gate** lever (the third token-efficiency/robustness lever alongside ceremony→binary and judgment→lazy-skill). The decision + grounding live in xa (`feedback-guarantee-binary-gate`, archived) and the binary-simplification roadmap.

## Problem (from xa's determination)

- The cycle count IS durable on-disk state (`### Feedback Cycles` entries), so a section-scoped counter is deterministic and tamper-evident (xa spike confirmed: ~25 lines, ignores a `Cycle N` line in a sibling section).
- A `--set status={feedback-to-target}` guard would FALSE-FIRE (xa spike confirmed): `--set status=implementation` carries no bounce signal — the disambiguating `is_feedback_reflow` lives on the dispatch-build input path, not as a `--set` field or durable state — so gating the status transition would refuse legitimate forward re-entry.
- The correct hook is the cycle-record WRITE itself: recording a feedback cycle is unambiguously a bounce event, so a dedicated command can own the append + count + escalation with no ambiguity.

## Proposed approach (to be formalized at ideation — grounding from xa)

A new `spacedock status --record-feedback-cycle {slug}` subcommand in `internal/status` that:
1. Appends a timestamped entry to `### Feedback Cycles` (creating the section if absent).
2. Computes the post-append cycle number from existing section-scoped entries.
3. On reaching the escalation threshold, stamps a durable escalation marker as a queryable frontmatter field (e.g. `feedback-escalate: cycle-3`, like `mod-block`) AND refuses any further auto-bounce record (exit 1), with a `--force` override in the established idiom.

The FO's feedback-rejection prose then changes from "Track cycles in `### Feedback Cycles`" / "On cycle 3, escalate" to "invoke `status --record-feedback-cycle`; on its refusal/escalation-marker, escalate." The `feedback-rejection-flow` skill (shipped via a9 #297) is the prose surface to update; the shared-core write-scope clause narrows (the FO triggers the write via the command rather than hand-editing the section).

**Out of scope:** the budget-probe fail-safe (xa Candidate 2 — non-mechanizable, stays prose + gq's live coverage); and whether the FO actually OBEYS a refusal by escalating (FO-LLM behavior with no in-process Go seam — that is gq's `feedback-3-cycle-escalation` live scenario, not this guard).

## Acceptance criteria (preliminary — ideation formalizes)

**AC-1 — `status --record-feedback-cycle {slug}` owns the `### Feedback Cycles` append and a section-scoped count.**
Verified by: a Go unit test in `internal/status` driving the command against a temp entity — first invoke appends `### Feedback Cycles` + a cycle-1 entry; Nth invoke appends a cycle-N entry; a `Cycle N` line in another section does not inflate the count.

**AC-2 — On the escalation threshold the command stamps a durable escalation marker and refuses a further auto-bounce.**
Verified by: a Go unit test asserting the threshold invoke stamps the queryable frontmatter marker AND the next auto-bounce attempt exits non-zero (with `--force` overriding), plus a negative control (strip the marker write in production → the refusal-on-Nth assertion goes red), mirroring `feedback_test.go` NEG-A.

**AC-3 — The `feedback-rejection-flow` skill prose invokes the command rather than instructing prose-tracking/prose-escalation.**
Verified by: a section-scoped presence oracle over the `feedback-rejection-flow` skill asserting the cycle-record/escalation steps invoke `status --record-feedback-cycle` (text claim, proven at its own level); the behavioral half — that the FO acts on a refusal — is gq's live scenario, not this AC.

## Test plan (from xa, scoped here)

- **Unit (Go, `internal/status`):** drive the command against temp entities; assert append/count/threshold-marker/refuse/`--force`/section-scoping. Byte-observable over the resulting on-disk file — same altitude as `archive_guard_test.go` / the live-run guard tests. Cost: low (no network, no live runtime).
- **Negative control:** strip the escalation-marker write and prove the refusal assertion goes red (mutation-proves-the-test).
- **Presence oracle (offline):** the skill-prose AC-3 check rides the existing `skill_text_test.go` `sectionAfter` pattern.
- **No spike needed:** xa already ran the two riskiest spikes (no `--set status` bounce signal; deterministic section-scoped count). This composes the already-proven body-parse + frontmatter-mutate + terminal-guard machinery (`live_proof.go`, `mutate.go`, `handlers.go`).
- **High-stakes → detached adversarial audit before merge:** this edits the `status` mutation/guard paths (a named high-stakes surface).

## Notes

Siblings: xa (`feedback-guarantee-binary-gate`, archived — the determination + grounding) and gq (`feedback-nonhappy-live-coverage` — the live half: proves the FO escalates when the guard refuses). Provenance: a9 (`feedback-rejection-flow-skill-extraction`) detached audit, 2026-06-04.
