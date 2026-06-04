---
title: FO should auto-continue from completed non-gated stages to next-stage dispatch
status: implementation
source: captain (2026-06-04) — FO stopped after implementation reporting instead of immediately advancing to validation; AI-engineer review found the current contract implies but does not enforce this lifecycle invariant
score: "0.32"
started: 2026-06-04T15:05:37Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-fo-auto-continues-after-stage-completion
issue:
id: wmn2x3k7j0fjshvdz126ray3
---

The first officer should not stop after an implementation worker completes and the implementation stage report is filed. For a non-gated, non-terminal stage, the FO should continue the workflow lifecycle: verify the stage report, advance to the next stage, and dispatch the next worker — all before ending its turn. In the dev workflow, that means implementation completion immediately advances to validation and dispatches an independent fresh validator unless a gate, terminal merge ceremony, blocker, or captain decision interrupts.

## Problem

The existing contract *implies* this behavior but never forbids the failure case, so a confused FO can stop early without violating any explicit clause. The two relevant spans both assume continuation rather than mandating it:

- `first-officer-shared-core.md` `## Completion and Gates` describes what happens when a worker completes (verify report, check gate, reuse-or-fresh) and says "**If fresh dispatch:** … Run `status --next` and dispatch the next stage." But it never says the FO MUST do this *before ending the turn*, and never names a completion-only captain-facing stop as forbidden. The non-gated branch (`If not gated: terminal → merge; else decide reuse-or-fresh.`) routes the FO into dispatch logic but does not close the door on stopping.
- `claude-first-officer-runtime.md` `## Event Loop` says "Repeat from step 1 after each agent completion until the captain ends the session." "Repeat from step 1" reads as *next iteration*, which a model can interpret as *next turn* — file the report, surface a status to the captain, and end the turn, expecting to resume later. Nothing pins the advance-and-dispatch as same-turn-mandatory.

Concretely, this let the FO report implementation completion for the Pi entity `pi-stage-dispatch-uses-build-artifact` and stop, instead of immediately advancing to validation and dispatching a fresh validator. The failure surfaced under Pi (`pi-subagents`), where the subagent result returns to the parent and the parent must itself drive the next lifecycle step — there is no separate event-loop turn to "resume" into.

This violates the workflow's own Working Principle ("prefer a code gate over a prose-only rule"): the lifecycle invariant lives only in implied prose, so its ceiling is wording-is-present, and there is no failable proof that an FO actually advances.

## Out of scope

- A binary-level state-machine command that emits the next mandatory action (e.g. a `dispatch next-action` helper). The captain's seed and the AI-engineer note floated this; it is a larger mechanism than the gap needs and would itself need its own ideation. This entity's deliverable is contract wording on three surfaces + a failable live regression scenario. If a helper is wanted later, it is a separate task.
- Changing the gated-stage behavior. Gates, terminal merge ceremony, blockers, and captain-decision interrupts are unchanged — they remain the explicit exceptions that DO stop the auto-continuation.
- The `claude-first-officer-runtime.md` `## Event Loop` and `codex-first-officer-runtime.md` are NOT amended here (the assignment scopes the three surfaces below). The host-neutral clause lands in the shared core; the Claude/Codex event loops already inherit it by reference. Only Pi gets a runtime-specific clause because its `pi-subagents` lifecycle is where the gap concretely bit.

## Proposed approach

Make the lifecycle invariant explicit on three surfaces, each backed by a failable check, plus one failable live regression scenario. The shared-core wording is host-neutral (it must not name a Claude-only event-loop step or helper — `prose_neutrality_test.go` polices the generic core).

### Surface 1 — `skills/first-officer/references/first-officer-shared-core.md`, `## Completion and Gates`

Insert a new clause immediately after the non-gated routing line. Current text (line 129):

> If not gated: terminal → merge; else decide reuse-or-fresh.

Add, on the line directly below it:

> **A completed non-gated, non-terminal stage is not a stopping point.** After verifying the report, the FO MUST advance the entity to the next stage and dispatch it (reuse-or-fresh per below) BEFORE ending its turn. It does not file a completion-only status and stop, waiting for the captain or a later turn to resume — advancing is the FO's own next action, not the captain's. The only spans that legitimately halt the turn here are: the next stage is `gate: true` (present the gate and wait), the entity is terminal (run the merge/cleanup ceremony), an explicit blocker (a rebase-conflict halt, an unmet clarification), or a captain decision the contract requires. Absent one of those, stopping after a completion-only report is a contract violation.

This sits inside `## Completion and Gates` alongside t3's and a9's edits — see the serialization note below.

### Surface 2 — `docs/dev/README.md`, `### implementation` stage definition

The dev workflow README must state that filing the implementation report is not a stopping point. Current `### implementation` `Outputs:` bullet (line 97):

> - **Outputs:** The deliverable committed to the relevant repo or state checkout, with a summary of what was produced and where

Add a following bullet under the same `### implementation` Outputs list:

> - Implementation completion is not a stopping point: once the deliverable is committed and the stage report filed, the entity routes immediately to independent `validation` dispatch — a fresh validator, since `validation` is `fresh: true` — unless a gate, blocker, terminal ceremony, or captain decision intervenes. The FO does not park a completed implementation and wait.

This is dev-workflow-specific (it names `validation` and `fresh: true` by name), which the shared-core clause cannot do (the shared core is workflow-agnostic).

### Surface 3 — `skills/first-officer/references/pi-first-officer-runtime.md`, `## Awaiting Completion`

Pi's `pi-subagents` lifecycle is where the gap bit. After the parent observes the subagent result and verifies the stage report, it must continue the shared lifecycle in the same turn — there is no separate event-loop iteration to defer into. Current `## Awaiting Completion` text (after the line about reading the entity file and verifying the stage report):

> For `pi-subagents`, the completion signal is the subagent result returned to the parent. After the result arrives, read the entity file and verify the stage report exactly as the shared core requires. Do not advance state based only on a cheerful worker summary.

Add a following sentence/paragraph in that subsection:

> Verifying the stage report is not the end of the parent's turn. Once the report is verified for a non-gated, non-terminal stage, the parent MUST continue the shared `## Completion and Gates` lifecycle in the same turn — advance the entity and dispatch the next stage (a fresh subagent when the next stage is `fresh: true`) — and only then return its final response. It does not return a completion-only result for the captain to resume from, unless the shared core's halt spans apply (next stage gated, terminal, blocked, or awaiting a captain decision).

### Serialization with t3 and a9 (same section, must not land concurrently)

All three of these in-flight entities edit `first-officer-shared-core.md` `## Completion and Gates`:

- **t3 — `gate-presentation-skill-extraction`** (id `t3w3`, implementation): extracts `## Gate Presentation` and rewrites the gated branch's "present the stage report per `## Gate Presentation` below" reference (around line 152 of the section).
- **a9 — `feedback-rejection-flow-skill-extraction`** (id `a9nt`, ideation): extracts the Feedback Rejection Flow that the gated branch's auto-bounce lines reference (around lines 154-155).
- **this entity — `wmn2`**: adds the non-gated-branch auto-continuation clause (around line 129-130).

The three edits touch different lines of the same section, so a concurrent worktree merge will conflict (the section is one contiguous block in one file). The implementation stage for this entity MUST serialize behind t3 (already in implementation) and a9: this entity's worker rebases onto whichever of t3/a9 has landed and re-anchors its insertion against the then-current `## Completion and Gates` text, rather than landing in parallel. The FO should not dispatch this entity's implementation worktree while t3's or a9's implementation worktree is open against the same file. (The shared-core split-root contract makes worktree edits land via merge; disjoint *lines* still conflict because they are the same `## Completion and Gates` block.)

### Riskiest unknown — AC-5's failable regression scenario (resolved at ideation, see Spike)

The design's soundness rests on one unverified question: *can a failable test catch an FO that stops after implementation, and what is the right substrate?* The Spike section below resolves the substrate choice and records the smallest probe that would invalidate it. The three contract-text ACs (AC-1/AC-2/AC-3) rest only on the already-proven `skill_text_test.go` `sectionAfter` presence-oracle pattern — proof at the text's own level, the legitimate case (the claim IS the text). AC-4 is the concrete scenario definition. AC-5 is the failable live regression.

## Spike (riskiest unknown — RESOLVED at ideation)

**The genuinely-novel question is the AC-5 substrate, not the contract wording.** AC-1/AC-2/AC-3 are presence oracles over instruction text — the `skill_text_test.go` `sectionAfter` pattern is already proven (it backs the live binary-absent / split-root / contract-range locks today), and the claim there IS the text, so proof at the text's own level is legitimate. The unexercised mechanism is AC-5: a test that starts from an implementation-ready dev entity, runs a real FO, and FAILS if the FO stops after the implementation report instead of advancing to validation and dispatching a fresh validator.

### Candidate substrates considered

| Substrate | What it would take | Verdict |
|---|---|---|
| **`internal/livescenario` primitive** (p4, `live-verification-gate`) — `Scenario{Runbook, Setup, Assert}` + `Runner` seam, graded on durable entity state before→after + observed output | author ONE scenario standalone, importable, run against the existing Claude live adapter (`claudeRunnerAdapter` in `livescenario_adapter_live_test.go`); offline negative case is a pure `Assert` over a hand-built broken end-state (no model spend) | **CHOSEN** |
| **Shared-scenario table** (`sharedRuntimeScenarios()` in `internal/ensigncycle`) — the gate-guardrail / rejection-flow / merge-hook-guardrail journeys | add a 4th entry to the host-neutral table + BOTH a Claude and a Codex runner | **REJECTED — too heavy** |
| **Fixture/golden or runtime assertion** (no real agent) | a frontmatter/state fixture asserting post-conditions | **REJECTED — cannot decide the claim** |

### Why livescenario, not the shared table

The shared-scenario table is parity-locked: `TestSharedRuntimeScenarioDefinitions` pins the EXACT name set `{gate-guardrail, rejection-flow, merge-hook-guardrail}` and reflects the type, so a 4th entry reds it until both the table AND both per-host runner maps are updated; `TestSharedScenarioRunnerCoverage` then requires a Claude AND a Codex runner for the new ID. That is the right machinery for a regression every host must independently exercise — but this entity's regression is the *generic* implementation→validation auto-continuation, and the p4 primitive was promoted EXACTLY so a scenario can be authored standalone (outside the ensigncycle parity machinery) and run against the existing Claude adapter. Using livescenario keeps the change to one new file plus the contract edits, with no edit to the pinned shared table. (A future task can promote it INTO the shared table if Codex parity is wanted; that is a separate, additive step, not a blocker here.)

### Why not a fixture/runtime assertion

A pure fixture can assert "given this end-state, the post-conditions hold" but cannot decide whether a *real FO actually advances* — the whole point of the regression is the live producer's behavior, which only a real agent run reveals (the same lesson `live-verification-gate` encoded: a recording proves the watcher, not the producer; a contract-text check proves the words, not the behavior). So AC-5 is a runtime-observable claim and its proof is a live run, gradable by the livescenario `Assert`.

### The smallest probe that would invalidate the choice (run at ideation)

The choice rests on two already-proven mechanisms; I ran the smallest exercise of each rather than assume:

1. **The livescenario primitive is live and green** in this tree (the offline grade path + the Claude adapter compile under `-tags live`). Ran `go test ./internal/livescenario/ -v` → `3 passed` (the `Scenario/Run/Runner` triple + `TestScenarioGradesBrokenOutcomeNegative` — a broken durable outcome reds the grade). This is the offline-negative-case mechanism AC-5 rides.
2. **A real FO dispatches a follow-up worker from a mid-lifecycle non-gated state, and the durable second stage report appears** — the structural twin of AC-5. The `rejection-flow` shared scenario already proves this for the symmetric direction (validation→implementation): the FO, given a mid-lifecycle report, dispatches the next worker and a SECOND `## Stage Report: implementation` plus a routed `status:` appears in durable state. Ran the offline negative guard for it: `go test ./internal/ensigncycle/ -run TestRejectionFlowNegativeMissingRoute` → PASS (the un-routed / partial-route end-states red the assertion). AC-5 is the mirror image: start with an implementation report, assert a validation stage report + `status: validation` (+ a fresh worktree) appears.

**Result:** no new mechanism. AC-5 composes (a) p4's livescenario primitive {runbook, setup, durable-outcome grade} run against (b) the existing Claude live adapter, with (c) the rejection-flow-style durable-state assertion. The probe confirmed both substrates are green in this tree, so the design only composes already-proven behavior. The one design decision the probe settled: grade on DURABLE STATE (a `## Stage Report: validation` appears AND `status:` advanced past implementation AND a fresh validation worktree exists), not on transcript phrasing — exactly the durable-vs-transcript discipline p4 and the shared assertions already use, so the grade greens only when the FO truly advanced and reds on a completion-only stop even if the transcript narrates intent to continue.

## Acceptance criteria

Each AC names a property of the finished entity and a "Verified by" that names something outside this task body that can fail.

**AC-1 — The FO shared contract's `## Completion and Gates` section forbids a completion-only stop after a non-gated, non-terminal stage and requires advancing/dispatching the next stage before ending the turn.**
Verified by: a Go presence oracle in `skills/integration` over `first-officer-shared-core.md`, scoped to the `## Completion and Gates` section via the existing `sectionAfter` helper (the proven `skill_text_test.go` pattern), asserting the section's non-gated branch carries the not-a-stopping-point clause AND the "before ending its turn" obligation AND names the halt exceptions (gate / terminal / blocker / captain decision). The check is scoped to that section, so the clause cannot be satisfied by unrelated prose elsewhere; deleting or gutting the clause reds the test. (Proof at the claim's own level: the claim IS the contract text. The behavioral half is AC-5's live regression.)

**AC-2 — The dev workflow README states implementation report filing is not a stopping point and routes immediately to fresh validation dispatch.**
Verified by: a Go presence oracle over `docs/dev/README.md` scoped to the `### implementation` stage section, asserting the Outputs carry a clause that implementation completion routes immediately to independent `validation` dispatch (naming `validation` and the `fresh` validator) unless gated/blocked/terminal. Scoped to the `### implementation` section so a stray mention elsewhere cannot satisfy it; gutting the clause reds the test.

**AC-3 — Pi runtime guidance preserves the same lifecycle after a `pi-subagents` completion.**
Verified by: a Go presence oracle over `pi-first-officer-runtime.md` scoped to the `## Awaiting Completion` section, asserting that after verifying the subagent result/stage report the parent must continue the shared `## Completion and Gates` lifecycle in the same turn (advance + dispatch, fresh subagent when `fresh: true`) and not return a completion-only result unless gated / terminal / blocked / awaiting a captain decision. Scoped to that section; gutting the clause reds the test.

**AC-4 — Ideation defines a concrete failable runtime scenario.**
Verified by: this entity's Spike + Test plan sections containing the AC-5 scenario definition — fixture shape (a split-root-free dev-shaped workflow: backlog→implementation→validation(`worktree: true, fresh: true, gate: true`)→done, with the entity parked at an implementation-ready state carrying a filed `## Stage Report: implementation`), the FO prompt (the neutral `Use $spacedock:first-officer` runbook, no "drive to done" coaching), the expected durable end-state (a `## Stage Report: validation` section appears OR `status:` advanced past implementation with a fresh validation worktree present and the validation gate presented), and the negative case (an end-state still at `status: implementation` with no validation stage report reds the grade even if the transcript narrates intent to continue). AC-4's own proof is AC-5 existing and running — it is not satisfied by this prose alone but by the AC-5 test below, which encodes this definition.

**AC-5 — A failable live regression scenario catches an FO that stops after implementation.**
Verified by: a live scenario authored via the `internal/livescenario` primitive {runbook, setup, durable-outcome `Assert`}, run by a real agent through the existing Claude live adapter, that starts from the AC-4 implementation-ready fixture and FAILS unless the FO advances to validation and dispatches/runs a fresh validator (or presents the validation gate). Plus an OFFLINE negative case (a pure `Assert` over a hand-built end-state still at `status: implementation` with no `## Stage Report: validation`) that reds the grade with no model spend — proving the assertion is behavior/state-oriented, not a transcript-shape tautology. The live half's `Verified by:` carries a `live <ci-run:|session:>` citation at terminal (this is a runtime-observable AC subject to p4's live-run gate under `require-external-proof`).

## Test plan

- **Spike DONE at ideation** (see Spike): both substrate mechanisms exercised green in this tree (`go test ./internal/livescenario/` 3/3; `go test ./internal/ensigncycle/ -run TestRejectionFlowNegativeMissingRoute` PASS). No new mechanism — AC-5 composes p4's livescenario primitive + the Claude live adapter + the rejection-flow-style durable-state assertion.
- **AC-1/AC-2/AC-3 (Go presence oracles, `skills/integration` + `internal/hostneutrality`):** one `sectionAfter`-scoped presence test per surface. Cost: trivial. Each reds when its clause is deleted or gutted; scoped to the named section so unrelated prose cannot satisfy it. Adversarial guard: gut each clause and confirm its test reds before restoring (the `skill_text_test.go` discipline). These are legitimate text-claim presence checks, NOT behavioral substring proxies — the behavioral half is AC-5.
- **AC-5 (live-scenario exercise + offline negative):** author one scenario via the promoted `internal/livescenario` primitive: a `Setup` that stages the implementation-ready dev fixture (backlog→implementation→validation(worktree, fresh, gate)→done), a `Runbook` that is the neutral `Use $spacedock:first-officer` prompt, and an `Assert` grading the durable end-state (validation stage report present / `status:` advanced / fresh validation worktree / gate presented). Run it against a real agent through the Claude adapter (`//go:build live`, like p4's `TestLivePrimitiveRunsAgainstClaudeAdapter`). Add an OFFLINE negative `Assert` case over a hand-built `status: implementation` end-state with no validation report → reds the grade (no model spend), proving the assertion is state-oriented not transcript-shaped. Cost: medium — needs a live credential + minutes of agent wallclock (live-gated). The offline negative + the contract presence tests run in the secret-free `go test ./...` lane; only the live half spends a credential.
- **High-stakes → detached adversarial audit BEFORE merge.** Trigger: this is the FO's OWN operating contract on the shipped scaffolding surface (the dev template names shipped contract/scaffolding as high-stakes). The audit runs read-only on a detached checkout of the merge result and tries to REFUTE: construct a contract edit that gutted the auto-continuation clause but kept a near-synonym, and confirm the AC-1/2/3 presence oracles red on it (catching a presence check too loose to notice a meaning-inverting paraphrase); and confirm the AC-5 offline negative reds on a transcript that narrates "advancing to validation" while the durable state stayed at `status: implementation` (catching a grade that trusts the transcript over the state).
- **Serialization gate (process, not a test):** the implementation worker for this entity rebases onto and re-anchors against the then-current `## Completion and Gates` after t3 and a9 land; the FO does not open this entity's implementation worktree concurrently with t3's or a9's against `first-officer-shared-core.md`.

## Notes

An AI-engineer review found this should not be docs-only. The solution is contract wording on three surfaces (AC-1/AC-2/AC-3, each backed by a scoped presence oracle) PLUS a failable live regression (AC-5) that proves a real FO actually advances — the behavioral half the presence checks cannot supply. The binary-helper "next-action command" the seed floated is deliberately out of scope (a larger, separate mechanism). The serialization overlap with t3 (`gate-presentation-skill-extraction`) and a9 (`feedback-rejection-flow-skill-extraction`) on `## Completion and Gates` is named so implementation does not land concurrently against the same section.

## Stage Report: ideation

- DONE: Pin the exact contract wording and insertion points, grounded against the live files, across all three surfaces and name the t3/a9 serialization overlap
  `## Proposed approach` pins before/after wording + insertion anchors for `first-officer-shared-core.md` `## Completion and Gates` (after line 129's non-gated routing line), `docs/dev/README.md` `### implementation` Outputs (after line 97), and `pi-first-officer-runtime.md` `## Awaiting Completion` (after the verify-the-report sentence). The serialization note names t3 (`gate-presentation-skill-extraction`, id t3w3, implementation, ~line 152) and a9 (`feedback-rejection-flow-skill-extraction`, id a9nt, ideation, ~lines 154-155) editing the same `## Completion and Gates` block; verified both entities' status/ids from the live state checkout.
- DONE: Resolve the riskiest unknown — AC-5's failable regression substrate — and run the smallest probe that would invalidate the choice
  `## Spike` chooses the `internal/livescenario` primitive over the parity-locked shared-scenario table (which `TestSharedRuntimeScenarioDefinitions` + `TestSharedScenarioRunnerCoverage` would force a 4th table entry + both host runners) and over a non-live fixture (cannot decide the producer's behavior). Probe: `go test ./internal/livescenario/` 3/3 green (offline grade + negative case) and `go test ./internal/ensigncycle/ -run TestRejectionFlowNegativeMissingRoute` PASS (the symmetric mid-lifecycle dispatch + durable-state assertion). Recorded: no new mechanism — AC-5 composes proven parts; grade on durable state (validation report + advanced status + fresh worktree), not transcript.
- DONE: Write entity-level ACs each backed by a code/test gate, every "Verified by" naming something outside the entity body that can fail
  AC-1/AC-2/AC-3 are `sectionAfter`-scoped presence oracles (the proven `skill_text_test.go` pattern) over the three contract surfaces — legitimate text-claim proof at the claim's own level. AC-4 is the concrete AC-5 scenario definition (fixture shape, FO prompt, durable end-state, negative case). AC-5 is the failable live regression via the `internal/livescenario` primitive run through the Claude live adapter, plus an offline state-oriented negative case, with a `live <ci-run:|session:>` citation subject to p4's live-run gate. Added an `## Out of scope` (the binary next-action helper is deferred) and a detached-audit trigger (this edits the FO's own shipped contract).

### Summary

Fleshed the ideation for the FO's own operating contract — a HIGH-STAKES edit. The gap is that `## Completion and Gates` and the Claude `## Event Loop` both ASSUME continuation without forbidding a completion-only stop; the fix lands an explicit, host-neutral not-a-stopping-point clause on the shared core, a dev-specific routes-to-fresh-validation clause on `docs/dev/README.md`, and a same-turn-continuation clause on the Pi `## Awaiting Completion` (where the regression bit). Each contract surface is backed by a section-scoped presence oracle; the behavioral half is a failable live regression (AC-5) authored via the already-promoted `internal/livescenario` primitive and graded on durable state, not transcript. The riskiest unknown (AC-5's substrate) is resolved with a probe confirming both proven mechanisms green in this tree — no new mechanism. Named the t3/a9 serialization overlap so implementation does not land concurrently against the same `## Completion and Gates` block.
