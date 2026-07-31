---
id: 1w62z8c5fq5g5cmhzf5sd79w
title: "Resolution-consume semantic hole: gate approval spends authority into terminal status before delivery proof, leaving no send-back path when delivery later fails"
status: ideation
source: "FO self-incident, 2026-07-31, in session fielding the fo-boot-install-hint-linux-direct-sandbox PR-merge ceremony. Captain caught the model–reality mismatch: status=done + verdict/completed unset + mod-block=merge:pr-merge, all after the validation approval was consumed into done. Only the credential delay kept the failure mode from being ratified."
started: 2026-07-31T14:27:47Z
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    current:
        gate: gate:1w62z8c5fq5g5cmhzf5sd79w:backlog
    records:
        - id: gate:1w62z8c5fq5g5cmhzf5sd79w:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:1w62z8c5fq5g5cmhzf5sd79w-backlog-1
              briefing:
                id: briefing:1w62z8c5fq5g5cmhzf5sd79w:backlog:attempt-1:revision-1
                digest: sha256:45c7f88e4a02a5eea0d3febe7431f09b3425a4e963fe0eff7bef7d9e5398bb84
                digest-domain: canonical-bytes
                request-digest: sha256:f808d3e34848fa57e0df5e17b9e278ccce23b1c7d8a77548f944bda23ce7a343
                room-ref: ./resolution-consume-terminal-before-delivery/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:1w62z8c5fq5g5cmhzf5sd79w:backlog:1
                briefing: briefing:1w62z8c5fq5g5cmhzf5sd79w:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-07-31T14:25:14.857596Z"
                decision: approve
                reason: 'Captain''s explicit order in this session (2026-07-31): dispatch ideation for 1w62 and discuss the proposed solution with him — accepts the direction that the consume-into-done semantic hole (caught at z3''s merge boundary) must be designed before further terminal-consume ceremonies ratify the pattern.'
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

A gate approval on a terminal-target stage (`validation → done`) is *immediately consumable into terminal status*: the binding approval's authority is marked spent (`consumed=true`) and the entity's status flips to `done` at consume time, while delivery (PR push, CI, merge) and workflow terminal fields (`verdict`, `completed`) remain pending. Three desyncs follow:

1. **Authority is spent before the spend is justified.** The consumed approval is irrevocable by design ("never double-applied"), yet the deliverable it ratified is only *pending* — there is no delivery proof at the moment the approval is spent.
2. **Two status classes per entity.** `status=done` (decision monoid) with empty `verdict`/`completed`/`pr` (delivery incomplete). The mod-block guard withholds terminal fields deliberately, but the status field itself already claims terminal. CL feedback: "right now if it's done it can't be sent back."
3. **No regression verb exists.** If the PR fails (CI red, closed-unmerged, a new defect uncovered), the workflow's own pre-done `feedback-to: implementation` wiring is *inert* (it's scoped to pre-done rejections). The pr-merge mod's answer to a closed PR is literally "report to the captain and wait"; the only physical reverse routes are frontmatter surgery or a second gate cycle that would have to dance around the consumed record. Prior live precedent (`worktree-clear-guard-misfire-on-backward-route`) shows even bookkeeping of a backward route trips the same merge-hook guard.

The pure-`merge: local` workflow hides this; any PR-gated workflow exposes it. Three review rounds missed it because they scoped to the API of gates, not the failure mode of delivery-after-consume.

## Problem

The consumption semantic treats "terminal-stage approval" and "terminal delivery" as one event; they are two, and the second can fail. Concretely, today:

1. **One atomic write, one irreversible spend.** `gates.Consume` (`internal/gates/application.go`) co-writes `application.state: consumed` and `status: done` in a single entity replacement for a terminal-target approval. The authority ledger and the workflow status both claim completion at the same instant, while the pr-merge mod deliberately withholds `verdict`/`completed` behind the merge-hook guard. Authority is spent with zero delivery proof in hand.
2. **Two status classes, one field.** A `status=done` entity is either *delivered* (verdict/completed set, archived or archivable) or *done-but-undelivered* (PR open/unpushed/CI unknown). Nothing in frontmatter distinguishes them without composing status + mod-block + pr + verdict — and CL's operational rule is "if it's done it can't be sent back," so the second class is a dead end.
3. **No recorded reverse route.** `applicationForDecision` wires `revise` to `feedback-to` only for pre-done gate rejections. The application's durable state enum (`pending/consumed/superseded/not-applicable` in `model.go::validateApplication`) has no value for "delivery failed after approval," `ValidateTransition` freezes closed attempts, and the only physical routes backward are frontmatter surgery or a second gate cycle that must dance around the consumed record (both rejected in the filing). The pr-merge mod's CLOSED-PR answer — "report to the captain and wait" — has no verb behind it.

The pure-`merge: local` workflow hides the hole (delivery is local and synchronous inside the terminal ceremony, so the window between spend and proof is small). Any PR-gated workflow exposes it. The model must make approval-authority, stage status, and terminal fields cohere in BOTH directions of the consume/delivery boundary: forward (spend lands only with delivery proof) and backward (a failed delivery is a recorded state transition, never surgery or re-bought authority).

## Proposed approach

**Delivery-bound consume + recorded regression verb.** Two new durable application states and one new verb; no README stage-grammar change.

### Durable gate-record vocabulary (stored-format change)

- `application.state: delivering` — the terminal-target approval is *bound* to delivery: authority committed, spend deferred. The entity is parked in-delivery. Eligibility treats it fail-closed (not consumable again), exactly as unknown states already fail closed.
- `application.state: reversed` + optional `application.reversal: {by, at, reason, ref}` — delivery failed after binding; the send-back is recorded without erasing or re-buying the approval. The reversed record is terminal history; `gates.current` is never rewound.
- Both states extend the guarded mutation path (`validateApplicationMutation`), generalized ~15 lines to permit the reversal subtree in the same swap. Frozen-attempt invariants otherwise unchanged.

### Forward direction (consume path)

- `gate consume` is unchanged for non-terminal targets and for terminal targets when **no merge hook is registered** (`merge: local`, hook-free): consume co-writes `consumed` + terminal status exactly as today (AC-3 pinned).
- When the target is terminal **and a merge hook is registered** (`merge: pr`): consume writes `application.state: delivering` atomically and leaves frontmatter `status` at the gated stage. The consume learns hook registration through the same mod-scan `merge guard` already performs, passed in from the CLI layer rather than a new scanner in gates.
- `merge guard <slug>` arms `mod-block` from the delivering state (arm condition moves from "terminal status + empty mod-block" to "delivering application + empty mod-block"); the blocked phase is unchanged. The pr-merge mod's existing "the entity stays at its current stage with `pr` set until the PR is merged" becomes the normative shape the consume side finally honors.
- On the merge sentinel, `merge guard --verdict passed` finalize consumes `delivering→consumed` **and** writes `status: done` + `verdict` + `completed` + mod-block clear in the finalize envelope. After this change, `status=<terminal>` without verdict/completed is unrepresentable: the canonical "is this entity really finished?" answer is a single frontmatter check, not a two-class composition (AC-2).

### Backward direction (delivery-failure path)

- New verb: `spacedock merge regress <slug> [--to <stage>] --reason <text> [--pr <ref>]`. Under the entity lock, one guarded replacement: `delivering→reversed` (+ reversal subtree), `status` := `--to` or the README `feedback-to` of the delivering record's stage (dev workflow: `validation → implementation`), `mod-block` cleared, `pr` cleared (PR ref preserved in `reversal.ref`). Refuses closed: no delivering application → error; `--to` names a terminal or undefined stage → error.
- Normal lifecycle resumes from the routed stage: implementation work continues with the existing worktree pointer intact (regression clears nothing, so the `worktree-clear-guard-misfire-on-backward-route` shape cannot arise); re-entry into the gated stage appends a successor attempt (`...-2`) on the SAME gate record via the proven `prepareTarget` successor path — no dancing around a consumed record, no second gate cycle.
- The pr-merge mod's CLOSED-unmerged and declined-draft prose gains `merge regress` as the structured route (the captain still decides; the verb is the mechanism once chosen).

### Boundary coherence (both directions)

| boundary event | authority (`application.state`) | status | terminal fields |
|---|---|---|---|
| consume, hookless terminal or non-terminal target | `pending→consumed` | → target (as today) | unchanged |
| consume, hooked terminal target | `pending→delivering` | stays at gated stage | unset |
| merge-sentinel finalize | `delivering→consumed` | → terminal | verdict+completed, same envelope |
| `merge regress` | `delivering→reversed` (+ reversal record) | → feedback-to target | unset (never written pre-delivery) |

### Surfaces that change

- **Go/CLI:** `internal/gates/model.go` (state enum, `Reversal` type, validation, phase/summary mapping incl. new display states), `internal/gates/application.go` (delivering-consume branch, generalized mutation validator, hook-registered input), `internal/status/merge.go` (arm-from-delivering, finalize spends application, `merge regress` subcommand), `internal/cli/{cli,help}.go` (verb routing + usage). Dispatch selectors must not dispatch delivering entities — audit `internal/dispatch` selection during implementation (expected +0–60 LOC).
- **Spec:** `docs/specs/gate-resolution-frontmatter-contract.md` (new states, reversal record, mutation rules).
- **Mod/skills:** `mods/pr-merge.md` (regress routes), `skills/first-officer/references/fo-merge-core.md` (armed/finalize/regress semantics, delivering resume class), `skills/first-officer/references/first-officer-shared-core.md` (delivering application — not status=done — prompts the merge ceremony), `skills/fo-gate-lifecycle/SKILL.md` (delivering consume outcome). **Ensign contract: no change** — stage work is identical.

### Alternatives considered

1. **Stage split (`done` reserved for delivered; merge ceremony in a new non-terminal stage):** rejected — README stage-grammar change for every workflow, entity migration, and application target-stage machinery churn for the same invariant.
2. **Spend-deferred alone (approval stays `approved-pending` until the merge hook proves delivery):** insufficient — with status still at the gated stage, the application would remain re-consumable (`pending` + status==record.Stage is the eligibility precondition), breaking one-consume semantics; the `delivering` state is the fail-closed form of the same idea.
3. **Regression via `status --set` backward or a fresh gate cycle:** rejected in the filing — surgery has no record; a new gate cycle dances around the consumed record and invites double-application.

## Out of scope

Changing the pr-merge hook's PR body creation mechanics; changing `merge: local` workflow semantics (which hide the hole); revisiting the shipped fo-boot-install-hint work (its fallout is the incident, not the target of this task); changing CI-check polling mechanics in the pr-merge mod (the regression verb is the route regardless of how a failure is discovered).

## Acceptance criteria

Each AC names a property of the finished entity.

**AC-1 — A PR-gated task whose delivery fails after validation-gate approval has a structured, recorded path back to a working stage: the reversal is on the durable gate record, authority is never re-bought, and the same entity can later reach a true terminal state through the recorded path.** *(value-measuring AC: exercises the delivery-failure-after-approval scenario end to end)*
Verified by: `internal/cli/merge_regress_test.go::TestMergeRegressDeliveryFailureScenario` — a CLI-driven behavior fixture: entity with a `delivering` terminal-target application, `mod-block=merge:pr-merge`, `pr=#57` (closed unmerged). Runs `merge regress`, then asserts on disk: `status=implementation`, `application.state=reversed` with `reversal.{by,at,reason,ref=#57}`, `mod-block`/`pr` cleared, `gates.current` unchanged, exit 0. Then continues the scenario: successor gate attempt-2 on the same record, fresh approval, delivering consume, merge sentinel, `merge guard --verdict passed` → asserts `status=done` + `verdict=passed` + `completed` set + `application.state=consumed`. Independently movable baseline: the same fixture run against the pre-change binary (verb unknown) and against any variant whose finalize writes terminal status without verdict/completed or whose regression leaves `pr` set — the assertions move the wrong way. Fixture sibling `TestMergeRegressPreservesWorktree` asserts the worktree pointer survives regression untouched (the `worktree-clear-guard-misfire-on-backward-route` shape: no worktree clear is issued, no terminal guard is tripped).

**AC-2 — A single canonical "is this entity really finished?" answer is derivable from frontmatter for all four terminal-target hook classes (no hook, `merge: local`, `merge: pr`, a hypothetically-named future hook), and the done-but-undelivered shape is unrepresentable at every lifecycle moment.**
Verified by: `internal/status/finalize_status_guard_test.go::TestTerminalFinishedClassificationMatrix` — golden status-output fixtures over hook-class × lifecycle-moment (pre-consume, delivering, finalized, reversed) asserting the invariant `status==terminal ⟺ verdict≠"" ∧ completed≠""` on every row, one finished classification per row. The fixture generator refuses the current `status=done` + `mod-block=merge:*` + empty-verdict shape, so the test fails on the pre-change world.

**AC-3 — The hook-free terminal consume path is unchanged: a `merge: local` / hook-free task reaches terminal status exactly when its last gate approval is consumed, with no new state values written.**
Verified by: `internal/gates/application_test.go::TestConsumeHooklessTerminalUnchanged` (consume co-writes consumed + terminal status when no merge hook is registered) and a merge-guard regression fixture pinning the existing auto-arm/finalize behavior for the hookless path.

Supporting (mechanism) criteria, subordinate to AC-1..3:
- `internal/gates/gates_test.go` new cases: `delivering` and `reversed` validate; unknown states still fail closed; at most one pending OR delivering application per record; `delivering`/`reversed` are never re-consumable (`EvaluateEligibility` stays fail-closed).
- `internal/cli/merge_test.go::TestMergeRegressRefusals` — refuses with no delivering application and with a terminal/undefined `--to`.
- Dispatch-side test (named after the dispatch audit): an entity with a delivering application is not selected as stage work.

## Test plan

Fixture/CLI-first, per the incident: the delivery-failure-after-consume path is a first-class state, exercised through the binary, not grepped in prose.

- **Go unit tests** (`internal/gates`): enum/validation matrix, guarded mutations (`pending→delivering`, `delivering→consumed`, `delivering→reversed` + reversal subtree), eligibility fail-closed. Cost: low; in-memory plus YAML round-trips.
- **Behavior fixtures driving the binary** (`internal/cli`, `internal/status`): the AC-1 scenario and AC-2 matrix above; `merge regress` refusal matrix; hookless regression pinning (AC-3). Cost: medium — new fixture scaffolding mirroring the existing `merge_guard_test.go` harness.
- **Golden fixtures**: status table/JSON classification for delivering/reversed display states.
- **Live workflow smoke**: not required for the core claim; one manual confirmation pass of the pr-merge mod's CLOSED-PR prose route in the dev workflow post-merge.
- Estimated complexity: the gates/state changes are small and spike-proven; most cost is fixture authoring in cli/status.

### Expected surface (gated baseline)

Files + insertions, tolerance ±30% (deletions small):

| file | ±LOC | changeable semantics |
|---|---|---|
| `internal/gates/model.go`, `application.go`, `operation.go` | +220 | stored formats: `application.state` enum gains `delivering`, `reversed`; new `application.reversal` subtree (additive — existing entity files remain valid) |
| `internal/gates/*_test.go`, `testdata/` | +350 | — |
| `internal/status/merge.go` | +180 | command grammar: new `merge regress` verb; authority: armed/finalize semantics keyed to delivering; runtime: finalize atomically writes terminal status+verdict+completed+consumed |
| `internal/status/merge_guard_test.go`, `merge_regress_test.go` | +450 | — |
| `internal/cli/cli.go`, `help.go` (+ tests) | +120 | command grammar: `merge regress` usage/output; `gate consume` outcome classification under `merge: pr` |
| `internal/dispatch/` (audit-driven) | +0–60 | runtime behavior: delivering entities not dispatched as stage work |
| `docs/specs/gate-resolution-frontmatter-contract.md` | +60 | spec: states, reversal record, mutation rules |
| `mods/pr-merge.md` | ±30 | host-integration text: CLOSED/declined routes name `merge regress` |
| `skills/first-officer/references/fo-merge-core.md` | ±25 | FO contract: armed/finalize/regress semantics; delivering resume class |
| `skills/first-officer/references/first-officer-shared-core.md` | ±8 | FO contract: ceremony prompt is the delivering application |
| `skills/fo-gate-lifecycle/SKILL.md` | ±10 | consume outcome classification |

**Changeable semantics declared:** command grammar (new `merge regress`; `gate consume` result classification); stored formats (application state enum + reversal subtree — additive); authority (terminal-target approval under a registered merge hook is *bound*, not spent, until delivery proof); runtime behavior (consume under `merge: pr` no longer flips status to terminal at consume time; FO merges from the delivering state; dispatch must not pick up delivering entities).

### Doc diff proposal (user-visible surface)

This repo's user-visible surface for these verbs is the CLI help and the FO contract/skills above. Proposed wording deltas (implementation applies verbatim, gate reviews):

- `internal/cli/help.go`, merge section — add: `merge regress <slug> [--to <stage>] --reason <text> [--pr <ref>]` — Record a delivery failure on a delivery-bound terminal approval: reverses the application into `feedback-to`, clearing mod-block and pr. The approval is recorded as reversed, never re-bought.
- `skills/fo-gate-lifecycle/SKILL.md` consume paragraph — after "a terminal target enters the existing merge guard/hook": "Under a registered merge hook the consume returns `delivering` rather than flipping status; the entity then enters the merge guard/hook from its current stage with `status` unchanged."
- `mods/pr-merge.md`, startup/idle CLOSED branch — after "Wait for the captain's direction before taking action.": "If the captain directs a send-back, run `merge regress <slug> --reason ... --pr {N}` to route the entity through its configured `feedback-to`; do not clear `pr`/`mod-block` by hand."
- `docs/specs/gate-resolution-frontmatter-contract.md` — add `delivering`/`reversed` to the application-state grammar with the forward/backward transition table from this entity body.

### Spike result (riskiest unverified mechanism — exercised)

Spiked 2026-07-31, throwaway, reverted after the run (never committed): `TestSpikeDeliveryRegressionRoundTrip` in `internal/gates` with a temporary enum extension. **Result: PASS.** Candidate question: can a consume/regression verb pair round-trip a consumed gate record under lock without corrupting `gates.current`? Findings:

1. `pending→delivering` round-trips through the existing `validateApplicationMutation` unchanged (it is already parameterized by from/to) — full-doc validation, byte-exact mutation guard, `gates.current` intact.
2. `delivering` is fail-closed under `EvaluateEligibility`: not eligible, not consumed — no re-consumability introduced.
3. `delivering→reversed` with a `reversal{by,at,reason,ref}` subtree works with a ~15-line generalization of the same guarded writer; durable YAML marshal/unmarshal re-validates.
4. `prepareTarget` on re-entry appends successor attempt `...-2` on the SAME record; the frozen reversed attempt-1 is undisturbed; `gates.current` is never rewound — the regression route needs no successor machinery of its own.
5. `ValidateTransition` still rejects smuggling `reversed→consumed` through the ordinary transition guard (the one spike failure was a test-double aliasing bug — a shallow struct copy shares the Records backing array; deep-copy via YAML for mutation probes — not a mechanism violation).

The spike test body is the seed for implementation's first gates tests. Remaining unspiked mechanisms are proven ones already shipping: status/gates co-write under lock (`writeDocumentAndStatus` in `application.go`), mod registration scan (`scanMods` in `merge.go`), finalize envelope (`merge.go`).

### Risk notes

- The consume-side hook-registration lookup crosses a package boundary (`gates` learning about `_mods/`); implementation passes the fact in from `internal/cli` rather than growing a scanner in gates. Not spiked — the scanMods mechanism is proven in `merge.go`.
- `internal/dispatch` selection behavior for delivering entities is audited, not assumed; sized at +0–60 LOC within tolerance.

## Stage Report: ideation

- DONE: Design names both directions of the consume/delivery boundary: where approval authority, status, and terminal fields land on the consume path AND on the delivery-failure path, with a concrete regression route that is neither frontmatter surgery nor a second gate cycle dancing around a consumed record — and names which shipped surfaces change (status/consume verbs, gate records, pr-merge mod, FO/ensign contracts).
  "Delivery-bound consume + recorded regression verb" — boundary-coherence table names authority/status/terminal-fields per event both directions; regression route is `merge regress` (durable `delivering→reversed` + `application.reversal` record + feedback-to status write under lock); surfaces table names gates/application/merge/cli/dispatch, contract spec, pr-merge mod, FO core/merge-core/gate-lifecycle, and records ensign contract unchanged.
- DONE: Entity-level ACs with per-AC test names; at least one AC MEASURES the end value by exercising a delivery-failure-after-approval scenario (fixture/CLI-driven, not prose-grep) against an independently movable baseline; expected surface (files + LOC ±tolerance + changeable semantics: command grammar, stored formats, authority, runtime behavior) declared per README ideation rules.
  AC-1 value-measures via `internal/cli/merge_regress_test.go::TestMergeRegressDeliveryFailureScenario` (CLI fixture: regress then re-deliver to true-done; baseline moves if finalize splits status/verdict or regression leaves pr); AC-2 `TestTerminalFinishedClassificationMatrix`; AC-3 `TestConsumeHooklessTerminalUnchanged`; surface table with ±30% LOC tolerance and declared semantics (grammar/formats/authority/runtime).
- DONE: Riskiest unverified mechanism spiked with result recorded (candidate: whether a consume/regression verb pair can round-trip a consumed gate record under lock without corrupting gates.current), or explicit "no spike needed: {proven mechanisms}" on the record.
  Spiked 2026-07-31 (throwaway, reverted uncommitted): `TestSpikeDeliveryRegressionRoundTrip` PASS — pending→delivering uses existing guarded mutation unchanged; delivering fail-closed; delivering→reversed+reversal subtree needs ~15-line generalization; prepareTarget appends successor attempt on same record with gates.current intact; ValidateTransition still rejects reversed→consumed smuggling. One initial failure was a test aliasing bug, not a mechanism violation.

### Summary

Ideation delivered a full design for the consume-into-terminal semantic hole: approval authority is bound (new `delivering` state), not spent, when a merge hook defers delivery; terminal status, verdict, and completed land atomically only at merge-sentinel finalize; delivery failure routes back through a recorded `merge regress` verb into `feedback-to` without surgery or a fiction of re-bought authority. The riskiest mechanism (guarded gate-record round-trip without corrupting `gates.current`) was exercised in a PASS spike; remaining mechanisms are proven shipping code. The design awaits the ideation gate and the captain's verdict on the direction per his 2026-07-31 resolution.
