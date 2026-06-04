---
id: xae5tx4hhyce916x034y3q9x
title: Evaluate promoting feedback-rejection guarantees from prose to binary-enforced gates
status: ideation
source: "captain (2026-06-04) — the a9 detached audit's dominant finding is that FO-behavioral guarantees living as contract prose (3-cycle escalation, budget-probe fail-safe) are immune to static tests and only sparsely covered by live scenarios. Where a guarantee is mechanizable, enforcing it in the binary (a guard / a tracked field) is the stronger fix than a live scenario — the 'code gate over prose-only rule' working principle. A third decomposition lever next to ceremony->binary (p2) and judgment->skill (t3/a9)."
score: "0.24"
started: 2026-06-04T20:04:39Z
completed:
verdict:
worktree:
issue:
---

The audit exposed that some FO-behavioral guarantees are inherently un-static-testable (a text oracle can't prove the FO obeys "escalate on cycle 3"), and live-scenario coverage for them is expensive. But several of these guarantees are **mechanizable** — the binary could track and enforce them deterministically, the way `status --set`/`--archive` already refuse a terminal transition without `pr`/`mod-block`. That eliminates the body-vs-label gap entirely (the guarantee is no longer prose).

This is the **guarantee → code-gate** lever — a third token-efficiency/robustness lever alongside the 0.19.6 line's ceremony→binary (p2) and judgment→lazy-skill (t3/a9).

## Candidates to evaluate

- **3-cycle escalation.** Could be a binary-tracked feedback-cycle counter that the system enforces — refuse a 4th auto-bounce / force an escalation marker on the 3rd rejection. Eliminates the infinite-loop regression class. Analogous to the existing `mod-block` terminal guard.
- **Budget-probe fail-safe.** The probe is already a binary command (`spacedock dispatch context-budget`). Question: should the system *force* the consult (e.g. the dispatch/reuse path refuses to reuse an over-budget member) rather than instructing the FO in prose to consult it?

## Scope of this task

Ideation/decision: for each candidate, decide whether it is genuinely mechanizable (a deterministic guard over real state) vs inherently a model judgment (stays prose + live scenario). Where mechanizable, scope the binary change; where not, defer to `feedback-nonhappy-live-coverage`. Output is a decision + (for the mechanizable ones) a concrete guard design — not docs-only; per dev policy, a decision with nothing shipped belongs in the roadmap, so this either produces a guard worth implementing or records the determination that these stay prose-guarded.

## Grounding: what the existing guards actually are

Two real reference guards anchor the "mechanizable" bar:

- **The terminal-transition guard** (`internal/status/handlers.go` `runSet`, `internal/status/mutate.go` `runArchive`). A `--set`/`--archive` reads real frontmatter fields (`mod-block`, `pr`, `verdict`) and, for the live-run case, parses a SECTION OF THE ENTITY BODY (`classifyLiveACs` in `live_proof.go` walks `**AC-N` blocks), then REFUSES the transition with a specific stderr message + exit 1, with `--force` as the override. The decisive property: the guard sits on an **irreversible on-disk transition** (finalization / archival), and the signal it gates on is **durable on-disk state** (a field value, a body section), not a model intention.
- **The context-budget probe** (`internal/claudeteam/contextbudget.go`, surfaced as `spacedock dispatch context-budget --name X`). This is a READ-ONLY binary command that reports a member's transcript budget. It produces a reading; it does not gate anything.

A candidate is mechanizable only if (a) there is a durable on-disk transition to sit in front of, AND (b) the guard's decision is computable from durable on-disk state — not from a model's intent that leaves no on-disk trace.

## Decision

### Candidate 1 — 3-cycle escalation: MECHANIZABLE (partially), via a dedicated cycle-record command — NOT a `--set status` guard

**The cycle count is durable on-disk state.** The FO records each rejection round in the `### Feedback Cycles` body section (FO Write Scope, first-officer-shared-core.md:263). A deterministic, section-scoped counter over that section is trivial — exercised below (the spike counts entries and isolates the section in ~25 lines). So half the guard (the count) is real state.

**A `--set status={feedback-to-target}` guard would FALSE-FIRE — proven by spike.** The naive design (refuse a 4th `--set status=implementation` when validation declares `feedback-to: implementation`) is wrong: the binary cannot distinguish a feedback *bounce* from a legitimate *forward re-entry* into the same stage. A `--set status=implementation` carries only the target stage name; the disambiguating "this is a rejection bounce" signal is `is_feedback_reflow`, which lives on the **dispatch-build** input path (`internal/dispatch/build.go:157`), NOT as a `--set` field and NOT as durable on-disk state. The spike confirms `--set status=validation→implementation` succeeds with exit 0 and no bounce signal. Gating cycle-count on the status transition would refuse legitimate non-bounce re-entries.

**The correct hook is the cycle-record write itself.** Recording a feedback cycle is, unlike a status transition, *unambiguously* a bounce event. If the binary OWNS the `### Feedback Cycles` append via a dedicated command — `spacedock status --record-feedback-cycle {slug}` (working name) — it can count existing entries and deterministically refuse the 4th append, or stamp an escalation marker on the 3rd, over durable on-disk state with no ambiguity. This mirrors the live-run guard's "parse a body section, decide three-way, refuse" shape, but on a body-write command rather than `--set status`. The append moves from FO prose-discipline ("Track cycles in `### Feedback Cycles`") into binary-owned state, which simultaneously closes the body-vs-label gap AND makes the count tamper-evident.

The mechanizable scope for the downstream implementation task: a new `status --record-feedback-cycle {slug}` subcommand that (1) appends a timestamped entry to `### Feedback Cycles` (creating the section if absent), (2) computes the post-append cycle number from existing entries, (3) on reaching the escalation threshold writes a durable escalation marker (a frontmatter field, e.g. `feedback-escalate: cycle-3`, so it is queryable by `status` like `mod-block`) and refuses any further auto-bounce record, exit 1, `--force` override in the established idiom. The FO's prose changes from "on cycle 3, escalate" to "invoke `status --record-feedback-cycle`; on its refusal/escalation-marker, escalate" — the same prose→binary promotion `mod-block` already models.

### Candidate 2 — Budget-probe fail-safe: NOT MECHANIZABLE as a hard binary refusal — stays prose + live scenario

The probe is already a binary command, but the fail-safe is a **reuse-vs-fresh-dispatch DECISION**, and there is no irreversible on-disk transition for a guard to sit in front of. Reusing an over-budget agent vs fresh-dispatching produces **identical durable on-disk state** (`status={next_stage}`, same frontmatter) — the difference is purely which live runtime handle the FO messages, which leaves no on-disk trace a guard could read. Condition (a) of the mechanizability bar fails: no on-disk transition to gate. The guarantee is "the FO consults the probe and decides correctly," which is exactly the un-static-testable FO-behavioral class the a9 audit flagged. It stays prose (reuse condition 0, first-officer-shared-core.md:134 + 165) and defers to the live-scenario lever in `feedback-nonhappy-live-coverage`. (This resolves the open question that sibling's body raised: "Decide whether the budget-probe fail-safe is better served by the binary-gate task" — answer: no, it is not mechanizable.)

## Riskiest-mechanism spike (done first, per ideation policy)

The design's soundness rested on one unverified claim: *can a binary guard fire at the feedback-bounce moment over durable state?* Exercised before committing the design (throwaway spikes, results recorded):

1. **`--set status=X` carries no bounce signal** — staged an entity at `validation` (which declares `feedback-to: implementation`) with two prior `### Feedback Cycles` entries; ran `status --set status=implementation`; got `status: validation -> implementation`, exit 0, no signal distinguishing bounce from re-entry. This is what KILLED the `--set status` guard design and redirected the hook to the cycle-record write.
2. **Section-scoped cycle counting is deterministic** — a ~25-line counter isolates the `### Feedback Cycles` section and returns the highest `- Cycle N` entry (2 for a 2-cycle body, 0 for a body with no section, ignoring a `Cycle 9` line in a sibling section). Confirms the count half is trivially mechanizable and tamper-evident.

## Acceptance criteria

- **AC-1** — The roadmap/queue carries the determination that Candidate 1 is mechanizable via a `status --record-feedback-cycle` command (NOT a `--set status` guard, for the false-fire reason), and Candidate 2 is non-mechanizable and stays prose + live-scenario.
  - **Verified by:** this entity's body decision section is the recorded determination; its on-disk presence (the `## Decision` section with both per-candidate verdicts) is the checkable artifact a separate read confirms. This is a text claim ("the determination is recorded with both verdicts and the named boundary"), proven at its own level. The downstream guard implementation is a SEPARATE queued task — this ideation task ships the decision, not the guard.
- **AC-2** — The boundary with `feedback-nonhappy-live-coverage` is named: Candidate 2 (budget-probe) defers to the live-scenario sibling; Candidate 1's *escalation behavior* (does the FO actually escalate when the binary refuses?) stays a live concern even after the guard ships, because the guard enforces the count, not the FO's response to a refusal.
  - **Verified by:** the `## Boundary` section names both, and the sibling entity's open "decide whether budget-probe is better served by binary-gate" question is answered in this body. Checkable by reading both entities' bodies for the cross-reference.
- **AC-3** — The exact contract prose a guard would replace is quoted, so the downstream implementer knows the before/after.
  - **Verified by:** the `## Contract prose a guard would replace` section quotes the verbatim first-officer-shared-core.md lines (162-164 for the cycle record/escalation; 134 + 165 for the budget fail-safe) with their before→after.

## Test plan (for the downstream implementation task this decision authorizes)

This ideation task ships the decision only; the guard is a separate queued task. Its test plan, scoped here:

- **Unit (Go, `internal/status`):** drive `status --record-feedback-cycle {slug}` against a temp entity. Assert: (1) first invoke appends `### Feedback Cycles` + `- Cycle 1`; (2) Nth invoke appends `- Cycle N`; (3) at the escalation threshold the command stamps the durable escalation marker (frontmatter field) AND exits non-zero on the next auto-bounce attempt; (4) `--force` overrides with a warning in the `mod-block` idiom; (5) the counter is section-scoped (a `Cycle 9` line in another section does not inflate the count). Fixture-level, byte-observable over the resulting on-disk file — same proof altitude as the existing `archive_guard_test.go` / live-run guard tests. Cost: low (no network, no live runtime).
- **Negative control:** strip the escalation-marker write in production and prove the refusal-on-Nth assertion goes RED — the same mutation-proves-the-test discipline `feedback_test.go` NEG-A uses.
- **NOT covered by this guard (stays live):** that the FO *obeys* the refusal by escalating to the human is FO-LLM behavior with no in-process Go seam — that remains the `feedback-nonhappy-live-coverage` `feedback-3-cycle-escalation` scenario's job. The guard makes the count tamper-evident; the live scenario proves the FO's response. No spike needed for the guard's own mechanism beyond the two above — it composes the already-proven body-parse + frontmatter-mutate + terminal-guard machinery (`live_proof.go`, `mutate.go`, `handlers.go`).

## Boundary with the live-scenario sibling

`feedback-nonhappy-live-coverage` (gq) owns the inherently-behavioral guarantees. This task (xa) resolves which of the two shared candidates cross the mechanizability bar:

- **Candidate 1 (3-cycle):** SPLIT. The COUNT mechanizes here (binary-owned `### Feedback Cycles` + escalation marker); the FO's RESPONSE to a refusal (actually escalating to the human) stays live — the sibling's `feedback-3-cycle-escalation` scenario. So both tasks touch Candidate 1, at different layers: this one makes the count un-forgeable, the sibling proves the FO acts on it.
- **Candidate 2 (budget-probe):** ENTIRELY the sibling's. Non-mechanizable (no on-disk transition), defers in full to live coverage. This answers the open question the sibling's body posed.

## Contract prose a guard would replace

Verbatim from `skills/first-officer/references/first-officer-shared-core.md`:

- **Cycle record + escalation (mechanizable — Candidate 1)**, lines 162-164:
  > 2. Track cycles in `### Feedback Cycles` in the entity body.
  > 3. On cycle 3, escalate to the human instead of another round.

  Before→after: the FO is *instructed in prose* to track and to escalate. After the guard: the FO *invokes* `status --record-feedback-cycle {slug}`; the binary owns the append and the count, and refuses/marks on the threshold — the "track" and "on cycle 3" clauses become a tool invocation whose enforcement is binary-side. (Line 263's "The FO owns `### Feedback Cycles`" write-scope clause narrows: the FO triggers the write via the command rather than hand-editing the section.)

- **Budget fail-safe (NOT mechanizable — Candidate 2, quoted for the record)**, line 134 (reuse condition 0) + line 165 (its restatement in the rejection flow):
  > 0. Consult the runtime adapter's context-budget probe. If it reports the worker over budget OR the probe source is unavailable, dispatch fresh (fail-safe — never silent-reuse on an absent reading)…
  > 4. Consult the budget probe (reuse condition 0). If the old ensign is over budget or the source is unavailable, shut down and fresh-dispatch…

  These STAY prose — there is no on-disk transition to gate, so no guard replaces them. They defer to the live scenario.

## Notes

Sibling: `feedback-nonhappy-live-coverage` (the live-scenario lever for the inherently-behavioral guarantees). Provenance: a9 (`feedback-rejection-flow-skill-extraction`) detached audit, 2026-06-04.

## Stage Report: ideation

- DONE: For each candidate guarantee, decide whether it is genuinely MECHANIZABLE as a deterministic binary guard over real state vs inherently a model judgment, grounded against the existing guards.
  Candidate 1 (3-cycle): MECHANIZABLE via a dedicated cycle-record command (NOT a `--set status` guard); Candidate 2 (budget-probe): NOT mechanizable. Grounded against `handlers.go`/`mutate.go` terminal guards, `live_proof.go` body-parse, and `contextbudget.go`. See `## Decision`.
- DONE: For the mechanizable ones, scope the concrete binary change; for the non-mechanizable ones, record the determination that they stay prose + live-scenario (defer to `feedback-nonhappy-live-coverage`).
  Candidate 1 scope: `status --record-feedback-cycle {slug}` owning the `### Feedback Cycles` append + escalation marker. Candidate 2: recorded as prose + deferred to sibling. AC-1..3 + test plan in body.
- DONE: Output a decision with something SHIPPABLE per dev policy. Name the boundary with the live-scenario sibling and the exact contract prose a guard would replace.
  Shipped: the recorded decision (the checkable artifact for an ideation task per dev policy — guard is a separately-queued task). `## Boundary` names the sibling split; `## Contract prose a guard would replace` quotes first-officer-shared-core.md:162-164 and :134/:165 verbatim with before→after.

### Summary

Two spikes ran first (per ideation riskiest-mechanism policy): (1) `--set status=X` carries no feedback-bounce signal — KILLED the naive `--set status` guard and redirected the hook to a dedicated cycle-record command; (2) section-scoped cycle counting is deterministic. Determination: Candidate 1's count mechanizes (a `status --record-feedback-cycle` command, mirroring the `mod-block` prose→binary promotion), but the FO's RESPONSE to a refusal stays live; Candidate 2 (budget-probe) is non-mechanizable (no on-disk transition to gate) and stays prose + the live-scenario sibling. Pre-existing unrelated failure noted: `TestMigrationCheckFixturesParseConsistently` fails on `_debriefs/*.md` fixtures on a clean tree (no edits by me) — flagged to team-lead, out of scope for this task.
