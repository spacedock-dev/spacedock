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

- `application.state: delivering` — the terminal-target approval is *bound* to delivery: authority committed, spend deferred, and the binding names its hook (`mod-block: merge:{hook}` written in the same atomic consume). The entity is parked in-delivery. Eligibility treats it fail-closed (not consumable again), exactly as unknown states already fail closed.
- `application.state: reversed` + **required** `application.reversal: {by, at, reason, ref}` — delivery failed after binding; the send-back is recorded without erasing or re-buying the approval. The subtree is required nonblank on `reversed` and forbidden on every other state (`validateApplication` enforces both directions). The reversed record is terminal history; `gates.current` is never rewound.
- Both states extend the guarded mutation path (`validateApplicationMutation`), generalized ~15 lines to permit the reversal subtree in the same swap. Frozen-attempt invariants otherwise unchanged.

### Forward direction (consume path)

Classification is by **actual hook registration** (the same `scanMods` mod-scan `merge guard` performs), never by merge-policy name: `merge: local` with a registered hook mod is delivery-bound, not hook-free.

- `gate consume` is unchanged for non-terminal targets and for terminal targets when **no merge hook is registered at all** (truly hookless): consume co-writes `consumed` + terminal status exactly as today (AC-3 pinned). Bare terminal entry without verdict/completed on this path remains the deliberate semantics `internal/status/handlers.go:210-218` already documents.
- When the target is terminal **and a merge hook is registered** (`merge: pr`, or `merge: local` with a registered hook mod): consume writes `application.state: delivering` **and** `mod-block: merge:{hook}` in the SAME atomic write, and leaves frontmatter `status` at the gated stage. Arming in the consume write binds the selected hook identity durably: restart-safety comes from the recorded `mod-block`, not from a live rescan, and there is no crash-before-arm window. If **more than one** merge hook is registered, consume refuses closed (no silent selection). On restart, if the recorded hook mod is missing or renamed, `merge guard` fails closed naming the expected hook from `mod-block`. The consume learns registration through the same scan passed in from the CLI layer rather than a new scanner in gates.
- `merge guard <slug>` confirms the arm from the delivering state (arm condition moves from "terminal status + empty mod-block" to "delivering application + `mod-block` naming a registered hook"; the confirm is idempotent because consume already armed). The blocked phase is unchanged. The pr-merge mod's existing "the entity stays at its current stage with `pr` set until the PR is merged" becomes the normative shape the consume side finally honors.
- Boot/readiness: `delivering` gets an explicit readiness/display classification (a `delivering-awaiting-merge` row added to the `computeReadyGates` recognized set and the status table), so a parked entity is visible as awaiting delivery, never silently invisible.
- On the merge sentinel, `merge guard --verdict passed` finalize consumes `delivering→consumed` **and** writes `status: done` + `verdict` + `completed` + mod-block clear in the finalize envelope. For every delivery-bound application, `status=<terminal>` without verdict/completed is unrepresentable: the canonical "is this entity really finished?" answer is a single frontmatter check, not a two-class composition (AC-2).

### Backward direction (delivery-failure path)

- New verb: `spacedock merge regress <slug> [--to <stage>] --reason <text> --by <actor> [--pr <ref>]`. Under the shared entity lock (below), one guarded replacement: `delivering→reversed` (+ reversal subtree), `status` := `--to` or the README `feedback-to` of the delivering record's stage (dev workflow: `validation → implementation`), `mod-block` cleared, `pr` cleared (PR ref preserved in `reversal.ref`).
- **Recorded provenance (reversal contract):** `--reason` nonblank required; `--by` nonblank required (actor ref, e.g. `person:captain` — the CLI currently has no actor argument, so it gains one); `at` stamped by the verb. `reversal{by,at,reason,ref}` is **required nonblank iff `application.state == reversed` and forbidden on every other state** — `validateApplication` enforces both directions. `ref` is derived from the entity's current `pr` field; an optional `--pr` may supply it only when `pr` is empty, and a `--pr` that disagrees with the current `pr` is an error.
- **Route guard:** the default target is exactly the record stage's configured `feedback-to`. A `--to` override must name a defined, **nonterminal, non-gated working stage** (a gated target would suppress normal dispatch and defeat the send-back; same-stage and terminal/undefined targets are refused). Overrides beyond the configured `feedback-to` are authorized only by captain direction, recorded via `--by`; the stage-transition guard validates the target transition as it would any status write. Archived entities are read-only through the existing resolver (`internal/status/native_runner.go:530-556`), so a finalized-and-archived entity can never be regressed. Refuses closed: no delivering application (and no legacy-delivering equivalent, below) → error.
- Normal lifecycle resumes from the routed stage: implementation work continues with the existing worktree pointer intact (regression preserves the worktree: the routed stage is a working stage whose dispatch expects the pointer); re-entry into the gated stage appends a successor attempt (`...-2`) on the SAME gate record via the proven `prepareTarget` successor path — no dancing around a consumed record, no second gate cycle.
- The pr-merge mod's CLOSED-unmerged and declined-draft prose gains `merge regress` as the structured route (the captain still decides; the verb is the mechanism once chosen).

### Legacy adoption (guarded one-way migration)

A live entity already carries the legacy shape (`status: done`, blank verdict/completed, terminal application `state: consumed`, `mod-block: merge:pr-merge`, `pr: "#586"` — `docs/dev/.spacedock-state/fo-boot-install-hint-linux-direct-sandbox.md:4-8,125-137`). The new finalize and regress accept that shape as a **delivering-equivalent**, so no hand surgery:

- Shape test: `application.state==consumed` ∧ `status==terminal` ∧ `verdict==""` ∧ `mod-block` matches `merge:*`. Only this exact conjunction adopts.
- **Finalize** on a legacy-shaped entity writes `status` (already terminal) + `verdict` + `completed` + mod-block clear, then the shared resume phase; the gate record is untouched (authority was already spent as `consumed` — history preserved byte-for-byte).
- **Regress** on a legacy-shaped entity adopts `consumed→reversed` with the full reversal subtree (`ref` from the current `pr`), guarded through the same parameterized mutation validator with authority history (resolution, briefing bindings, consume fact) preserved; only the application state and subtree move.
- Tests drive fixtures copied from the live fo-boot entity shape in three delivery outcomes: PR open (finalize not yet possible — refuses), PR closed-unmerged (regress), merged sentinel `pr-merge:{N}` (finalize + archive-resume).

### Shared lock/CAS and crash resume

- **One lock, one compare-and-swap.** `lockEntity` (`internal/gates/operation.go:382-397`) generalizes into the shared entity lock honored by every writer on an entity carrying gate records: `status --set` sentinel writes (`pr=pr-merge:{N}` — which today rewrite frontmatter without that lock, `internal/status/mutate.go:71-193`), `merge regress`, and `merge guard` finalize. All three perform the full-document compare-and-swap the gates writer already does (`internal/gates/io.go:188-226`); a CAS loser fails closed with a stale-document error, so exactly one coherent winner exists per race. Finalize and regress reuse the status layer's terminal/proof/worktree validators rather than re-implementing them (one validation boundary), so no finalize writer can bypass those guards.
- **Archive-resume phase.** An entity in the post-crash shape `application.state==consumed` ∧ `status==terminal` ∧ verdict/completed set ∧ still **active** (not archived) is treated as an idempotent resume point: finalize re-runs archive-move, commit, and publication without rewriting terminal fields. Existing archived recovery handles the archived half; this covers the active-but-finalized half the current finalize leaves behind on a crash between its separate steps (mutation, archive, commit, publish — `internal/status/merge.go:374-421`). Injected-crash tests cut power after the atomic replacement, after the archive move, after the commit, and after publication.

### Boundary coherence (both directions)

| boundary event | authority (`application.state`) | status | terminal fields |
|---|---|---|---|
| consume, truly hookless terminal or non-terminal target | `pending→consumed` | → target (as today) | unchanged |
| consume, hooked terminal target (any policy with a registered hook) | `pending→delivering` + `mod-block: merge:{hook}` in one write | stays at gated stage | unset |
| merge-sentinel finalize | `delivering→consumed` | → terminal | verdict+completed, same envelope |
| `merge regress` | `delivering→reversed` (+ required reversal record) | → configured `feedback-to` target | unset (never written pre-delivery) |
| legacy finalize (adoption) | `consumed` (unchanged) | → terminal (already) | verdict+completed written |
| legacy regress (adoption) | `consumed→reversed` (+ reversal record) | → configured `feedback-to` target | cleared/unset |
| crash-resume finalize | `consumed` (unchanged) | terminal (already) | verdict+completed (already) — resume archive/commit/publish |

### Surfaces that change

- **Go/CLI:** `internal/gates/model.go` (state enum, `Reversal` type, validation, phase/summary mapping incl. new display states/readiness row), `internal/gates/application.go` (delivering-consume branch with atomic mod-block arm, generalized mutation validator, hook-registered input), `internal/gates/operation.go` (lockEntity generalized to the shared entity lock), `internal/status/mutate.go` (`status --set` honors the shared lock + full-document CAS), `internal/status/merge.go` (arm-confirm from delivering, finalize spends application, archive-resume phase, `merge regress` subcommand, shared validation boundary), `internal/status/format.go` (delivering/reversed readiness classification), `internal/cli/{cli,help}.go` (verb routing + usage). Dispatch selectors must not dispatch delivering entities — audit `internal/dispatch` selection during implementation (expected +0–60 LOC).
- **Spec:** `docs/specs/gate-resolution-frontmatter-contract.md` (new states, required reversal record, mutation rules, legacy-adoption shape test).
- **Entity schema:** `docs/schema/entity.mdschema.yml` (gates invariants: closed attempts frozen *except* the enumerated guarded application-state transitions; regression/adoption writers named).
- **Public docs:** `docs/site/reference/command-reference.md` (`gate consume` no longer always advances on hooked terminal targets; `merge regress`), `docs/site/reference/frontmatter-contract.md` (delivering/reversed application lifecycle), `docs/site/concepts/gates-and-decisions.md` (terminal-target consume binds, finalize spends).
- **Workflow mod + plugin mod (distinct files):** `docs/dev/_mods/pr-merge.md` (regress routes; the two-step finalize note becomes the delivering+sentinel+one-envelope path), `mods/pr-merge.md` (shipped generic mod text).
- **Skills:** `skills/first-officer/references/fo-merge-core.md` (armed/finalize/regress semantics, delivering resume class), `skills/first-officer/references/first-officer-shared-core.md` (delivering application — not status=done — prompts the merge ceremony), `skills/fo-gate-lifecycle/SKILL.md` (delivering consume outcome). **Ensign contract: no change** — stage work is identical.

### Alternatives considered

1. **Stage split (`done` reserved for delivered; merge ceremony in a new non-terminal stage):** rejected — README stage-grammar change for every workflow, entity migration, and application target-stage machinery churn for the same invariant.
2. **Spend-deferred alone (approval stays `approved-pending` until the merge hook proves delivery):** insufficient — with status still at the gated stage, the application would remain re-consumable (`pending` + status==record.Stage is the eligibility precondition), breaking one-consume semantics; the `delivering` state is the fail-closed form of the same idea.
3. **Regression via `status --set` backward or a fresh gate cycle:** rejected in the filing — surgery has no record; a new gate cycle dances around the consumed record and invites double-application.

## Out of scope

Changing the pr-merge hook's PR body creation mechanics; changing the consume path of workflows with **no registered merge hook at all** (truly hookless — unchanged by design); `merge: local` *with a registered hook* is deliberately in scope (it is delivery-bound, not hook-free); shipping changes to the fo-boot-install-hint deliverable itself (its legacy-shaped live entity becomes the guarded adoption target and fixture source, not a rework target); changing CI-check polling mechanics in the pr-merge mod (the regression verb is the route regardless of how a failure is discovered).

## Acceptance criteria

Each AC names a property of the finished entity.

**AC-1 — A PR-gated task whose delivery fails after validation-gate approval has a structured, recorded path back to a working stage: the reversal is on the durable gate record, authority is never re-bought, and the same entity can later reach a true terminal state through the recorded path.** *(value-measuring AC: exercises the delivery-failure-after-approval scenario end to end)*
Verified by: `internal/cli/merge_regress_test.go::TestMergeRegressDeliveryFailureScenario` — a CLI-driven behavior fixture **starting from a real pending approval**: README workflow with a registered pr-merge hook mod, binding approval recorded by `gate record`, then CLI `gate consume` → asserts `application.state=delivering` + `mod-block=merge:pr-merge` armed in the consume write + status still at validation. Then PR reported closed-unmerged: `merge regress --reason … --by person:captain` → asserts on disk: `status=implementation`, `application.state=reversed` with `reversal.{by,at,reason,ref=#57}`, `mod-block`/`pr` cleared, `gates.current` unchanged, exit 0. Then continues the scenario: successor gate attempt-2 on the same record, fresh approval, delivering consume, merge sentinel, `merge guard --verdict passed` → asserts `status=done` + `verdict=passed` + `completed` set + `application.state=consumed`. Independently movable baseline: the same fixture run against the pre-change binary (verb unknown, consume flips straight to done — the hole, visibly reproduced) and against any variant whose finalize writes terminal status without verdict/completed or whose regression leaves `pr` set — the assertions move the wrong way. Fixture sibling `TestMergeRegressPreservesWorktree` asserts the worktree pointer survives regression untouched (the `worktree-clear-guard-misfire-on-backward-route` shape: no worktree clear is issued, no terminal guard is tripped).

**AC-2 — For every delivery-bound application (any policy with a registered merge hook), a single canonical "is this entity really finished?" answer is derivable from frontmatter, and the done-but-undelivered shape is unrepresentable on the hooked paths.** The invariant `status==terminal ⟺ verdict≠"" ∧ completed≠""` is asserted on entity shapes reachable from deliver-bound applications; the truly hookless path deliberately keeps bare terminal entry (`internal/status/handlers.go:210-218`) and is excluded from the invariant by classification, not overlooked.
Verified by: `internal/status/finalize_status_guard_test.go::TestTerminalFinishedClassificationMatrix` — golden status-output fixtures over hook-registered class (`merge: pr`, `merge: local` with registered hook, a hypothetically-named future hook mod) × lifecycle-moment (pre-consume, delivering, finalized, reversed, legacy-adopted), asserting the invariant on every hooked row and one finished classification per row. **Every supported terminal-transition producer is exercised** — hookless consume, hooked consume, sentinel finalize, legacy adoption finalize, and regress — so the matrix covers production reachability, not fixture policy. Adversarial discriminator: leaving the old always-consume branch intact (or authoring a `status=done` + `mod-block=merge:*` + empty-verdict row) turns these tests red.

**AC-3 — The truly hookless terminal consume path (no merge hook registered at all) is unchanged: such a task reaches terminal status exactly when its last gate approval is consumed, with no new state values written; and a `merge: local` workflow WITH a registered local hook is delivery-bound, not hook-free.**
Verified by: `internal/gates/application_test.go::TestConsumeHooklessTerminalUnchanged` (consume co-writes consumed + terminal status when no merge hook is registered), a merge-guard regression fixture pinning the existing auto-arm/finalize behavior for the hookless path, and a **registered-local consume test expecting `delivering`** (classification is by actual hook registration, refuting the "`merge: local` is hookless" shortcut `internal/status/merge_guard_test.go:69-88` shows to be false).

Supporting (mechanism) criteria, subordinate to AC-1..3:
- `internal/gates/gates_test.go` new cases: `delivering` and `reversed` validate; unknown states still fail closed; at most one pending OR delivering application per record; `delivering`/`reversed` are never re-consumable (`EvaluateEligibility` stays fail-closed); `reversal` required-iff-reversed both directions.
- `internal/cli/merge_test.go::TestMergeRegressRefusals` — refuses with no delivering application, with a terminal/undefined/gated/same-stage `--to`, with blank `--reason`/`--by`, and with a `--pr` that disagrees with the current `pr`.
- Dispatch-side test made discriminating (named after the dispatch audit): a delivering entity surfaces the `delivering-awaiting-merge` readiness classification and is never emitted as stage work by `dispatch build` — discriminating because the pre-change world has no such classification row and the entity's consume would have flipped status to the terminal stage instead.
- Race tests under the shared lock: sentinel-write-versus-regress and finalize-versus-regress each run with the two operations interleaved; one coherent winner or a clean stale-document refusal — never a torn state.
- Crash-injection tests: finalize interrupted after the atomic replacement, after the archive move, after the commit, and after publication; re-run resumes idempotently to one correct end state.
- Legacy-adoption tests from the fo-boot entity shape: open PR (refuse), closed-unmerged (regress), merged sentinel (finalize + resume).

## Test plan

Fixture/CLI-first, per the incident: the delivery-failure-after-consume path is a first-class state, exercised through the binary, not grepped in prose.

- **Go unit tests** (`internal/gates`): enum/validation matrix incl. reversal required-iff-reversed, guarded mutations (`pending→delivering` + atomic arm, `delivering→consumed`, `delivering→reversed`, legacy `consumed→reversed`), eligibility fail-closed. Cost: low; in-memory plus YAML round-trips.
- **Behavior fixtures driving the binary** (`internal/cli`, `internal/status`): the AC-1 scenario (from a real pending approval) and AC-2 matrix above; `merge regress` refusal matrix; hookless and registered-local classification pinning (AC-3); race tests and crash-injection tests above; legacy-adoption fixtures. Cost: medium — new fixture scaffolding mirroring the existing `merge_guard_test.go` harness.
- **Golden fixtures**: status table/JSON classification for delivering/reversed display states and the `delivering-awaiting-merge` readiness row.
- **Retained implementation-opening spike (committed, not reverted):** before feature implementation, spike the full cross-package transaction — CLI `gate consume` → delivering+arm → sentinel `status --set` → atomic finalize, plus the regress race paths — and commit the test/patch with captured output so the evidence is independently reviewable and seeds the first tests.
- **Detached adversarial audit (required — high-stakes surface):** the diff touches the status mutation/guard paths and the shipped FO contract, two of the four audit-triggering surfaces (`docs/dev/README.md:82-83`); a read-only throwaway-checkout audit attempts adversarial edits the deliverable's tests must catch.
- **Live lanes + live drive (required):** the diff touches shipped FO contract/skill files, so the applicable live lanes per the path→lane mapping must be green; and as a contract/skill change, it is PASSED only after a durable live workflow scenario observed covering consume→deliver AND consume→regress→re-enter→deliver (`docs/dev/README.md:198`).
- Estimated complexity: the gates/state changes are small; most cost is fixture scaffolding in cli/status, the shared-lock races, and the live pass.

### Expected surface (gated baseline)

Files + insertions, tolerance ±30% (deletions small):

| file | ±LOC | changeable semantics |
|---|---|---|
| `internal/gates/model.go`, `application.go`, `io.go`, `operation.go` | +260 | stored formats: `application.state` enum gains `delivering`, `reversed`; new required `application.reversal` subtree (additive — existing entity files remain valid); shared entity lock |
| `internal/gates/*_test.go`, `testdata/` | +350 | — |
| `internal/status/merge.go`, `mutate.go`, `format.go` | +260 | command grammar: new `merge regress` verb (with `--by`/`--reason` required); authority: armed/finalize semantics keyed to delivering; runtime: finalize atomically writes terminal status+verdict+completed+consumed; `status --set` honors shared lock + CAS; readiness classification |
| `internal/status/merge_guard_test.go`, `merge_regress_test.go`, race/crash/legacy fixtures | +560 | — |
| `internal/cli/cli.go`, `help.go` (+ tests) | +120 | command grammar: `merge regress` usage/output; `gate consume` outcome classification under `merge: pr` |
| `internal/dispatch/` (audit-driven) | +0–60 | runtime behavior: delivering entities not dispatched as stage work |
| `docs/specs/gate-resolution-frontmatter-contract.md` | +70 | spec: states, required reversal record, mutation rules, adoption shape test |
| `docs/schema/entity.mdschema.yml` | ±12 | schema invariants: guarded application transitions, regression/adoption writers |
| `docs/site/reference/command-reference.md`, `docs/site/reference/frontmatter-contract.md`, `docs/site/concepts/gates-and-decisions.md` | ±50 | user-visible semantics: consume binds under a registered hook; regress verb; delivering/reversed lifecycle |
| `docs/dev/_mods/pr-merge.md` | ±25 | workflow-local mod: regress routes; delivering/sentinel ceremony supersedes the two-step/CLOSED prose |
| `mods/pr-merge.md` | ±30 | shipped mod text: CLOSED/declined routes name `merge regress` |
| `skills/first-officer/references/fo-merge-core.md` | ±25 | FO contract: armed/finalize/regress semantics; delivering resume class |
| `skills/first-officer/references/first-officer-shared-core.md` | ±8 | FO contract: ceremony prompt is the delivering application |
| `skills/fo-gate-lifecycle/SKILL.md` | ±10 | consume outcome classification |

**Changeable semantics declared:** command grammar (new `merge regress` with required `--by`/`--reason`; `gate consume` result classification under a registered hook); stored formats (application state enum + required reversal subtree — additive); authority (terminal-target approval under a registered merge hook is *bound* with its hook identity recorded, not spent, until delivery proof); runtime behavior (consume under a registered hook no longer flips status to terminal at consume time — `merge: local` with a registered hook included; FO merges from the delivering state; dispatch/readiness exposes delivering entities as awaiting delivery, not stage work; `status --set` on gated entities is lock+CAS-guarded).

### Doc diff proposal (user-visible surface)

This repo's user-visible surface for these verbs is the CLI help and the FO contract/skills above. Proposed wording deltas (implementation applies verbatim, gate reviews):

- `internal/cli/help.go`, merge section — add: `merge regress <slug> [--to <stage>] --reason <text> --by <actor> [--pr <ref>]` — Record a delivery failure on a delivery-bound terminal approval: reverses the application into the record stage's configured `feedback-to`, clearing mod-block and pr. The approval is recorded as reversed (with by/at/reason/ref), never re-bought.
- `skills/fo-gate-lifecycle/SKILL.md` consume paragraph — after "a terminal target enters the existing merge guard/hook": "Under a registered merge hook the consume returns `delivering` rather than flipping status; the entity then enters the merge guard/hook from its current stage with `status` unchanged."
- `mods/pr-merge.md`, startup/idle CLOSED branch — after "Wait for the captain's direction before taking action.": "If the captain directs a send-back, run `merge regress <slug> --reason ... --by person:captain` to route the entity through its configured `feedback-to`; do not clear `pr`/`mod-block` by hand." (The parallel `docs/dev/_mods/pr-merge.md` delta is listed above.)
- `docs/specs/gate-resolution-frontmatter-contract.md` — add `delivering`/`reversed` to the application-state grammar with the forward/backward transition table from this entity body.
- `docs/site/reference/command-reference.md:93` — `gate consume` row: "Atomically advance status and spend an eligible pending approval once" → "Atomically advance status and spend an eligible pending approval once; under a terminal target with a registered merge hook the approval is *bound* (`delivering`) and spent only at the merge-sentinel finalize." Add a `merge regress` row: "Record a delivery failure on a delivery-bound approval: reverse the application into the record stage's `feedback-to`, clearing mod-block and pr; the approval is recorded as reversed, never re-bought."
- `docs/site/reference/frontmatter-contract.md:17` — application lifecycle sentence gains: "`delivering` when a registered merge hook defers delivery (bound, not yet spent), `reversed` with a required `reversal{by,at,reason,ref}` record when delivery failed after binding."
- `docs/site/concepts/gates-and-decisions.md:64-75` — "Approval to a terminal target is consumed before the existing merge and terminalization path begins" → "Approval to a terminal target is *bound* (state `delivering`, hook identity recorded) when a merge hook is registered; the merge-sentinel finalize performs the consume; a failed delivery routes back through `merge regress` as a recorded reversal."
- `docs/schema/entity.mdschema.yml:107-118` — gates `writer`/invariants: the gates writer gains the enumerated guarded application transitions (`pending→delivering`, `delivering→consumed`, `delivering→reversed`, legacy `consumed→reversed` adoption); "closed attempts are frozen" becomes "closed attempts are frozen except through the enumerated guarded application-state transitions."
- `docs/dev/_mods/pr-merge.md:18-25` — the two-step finalize prose is superseded: sentinel `--set` then `merge guard` finalize (which under delivering consumes the binding and terminalizes in one envelope); the CLOSED branch's "Wait for the captain's direction" gains "If the captain directs a send-back, run `spacedock merge regress {slug} --reason ... --by person:captain` to route the entity through its configured `feedback-to`; do not clear `pr`/`mod-block` by hand."
- Consistency: a docs-consistency check exercising `spacedock gate consume --help` / `merge regress --help` output against the documented verbs (output-driven, not prose-grep).

### Spike result (riskiest unverified mechanism — exercised)

Spiked 2026-07-31, throwaway, reverted after the run (never committed): `TestSpikeDeliveryRegressionRoundTrip` in `internal/gates` with a temporary enum extension. **Result: PASS.** Candidate question: can a consume/regression verb pair round-trip a consumed gate record under lock without corrupting `gates.current`? Findings:

1. `pending→delivering` round-trips through the existing `validateApplicationMutation` unchanged (it is already parameterized by from/to) — full-doc validation, byte-exact mutation guard, `gates.current` intact.
2. `delivering` is fail-closed under `EvaluateEligibility`: not eligible, not consumed — no re-consumability introduced.
3. `delivering→reversed` with a `reversal{by,at,reason,ref}` subtree works with a ~15-line generalization of the same guarded writer; durable YAML marshal/unmarshal re-validates.
4. `prepareTarget` on re-entry appends successor attempt `...-2` on the SAME record; the frozen reversed attempt-1 is undisturbed; `gates.current` is never rewound — the regression route needs no successor machinery of its own.
5. `ValidateTransition` still rejects smuggling `reversed→consumed` through the ordinary transition guard (the one spike failure was a test-double aliasing bug — a shallow struct copy shares the Records backing array; deep-copy via YAML for mutation probes — not a mechanism violation).

The spike test body is the seed for implementation's first gates tests. Remaining mechanisms claimed proven: status/gates co-write under lock (`writeDocumentAndStatus` in `application.go`), mod registration scan (`scanMods` in `merge.go`), finalize envelope (`merge.go`) — but the **cross-package composition** of those mechanisms (CLI consume → deliver+arm → sentinel write → one-envelope atomic finalize, and the regress race paths) is new and unspiked; it is owed as the retained, committed implementation-opening spike named in the test plan, with output captured for independent review.

### Staff Review

Round 1 answers to the FO's staff reviewer. Provenance: builtin `reviewer` subagent on model `openai-codex/gpt-5.6-sol:xhigh`, run `5da3cbb1-41ac-467d-b024-6d29f676bac7` (OVERALL: reject). Each entry: compressed disagreement, severity, answer, resolution.

1. **Done-but-undelivered unrepresentable spans hookless / `merge: local` wording (blocker). Answer: accept-and-revise.** Verified the citations: consume writes only `consumed`+terminal status (`internal/gates/application.go:176-180`); bare terminal entry without verdict/completed is deliberate (`internal/status/handlers.go:210-218`); registered local hooks arm today (`internal/status/merge_guard_test.go:69-88`). Resolution: classification is by actual hook registration only; AC-2's invariant narrowed to delivery-bound applications; AC-3 rewritten (hookless = no registered hook; registered-local expects `delivering`); all "`merge: local` / hook-free" wording rewritten across Forward direction, Out of scope, ACs, surface table.
2. **Stored-format change needs a migration story (blocker). Answer: accept-and-revise.** Verified the live legacy shape in `fo-boot-install-hint-linux-direct-sandbox.md:4-8,125-137`. Resolution: new "Legacy adoption" section — a guarded shape test adopts `consumed`+terminal+blank-verdict+`mod-block: merge:*` as delivering-equivalent for finalize (record untouched, authority history preserved) and regress (`consumed→reversed` with full reversal subtree), with fixtures for open/closed/merged-sentinel variants; AC-2, test plan, expected surface updated.
3. **Merge-sentinel finalize not re-entrant past the frontmatter write (material). Answer: accept-and-revise.** Verified finalize's separate steps (`internal/status/merge.go:374-421`). Resolution: `consumed + terminal + verdict/completed + active` defined as an idempotent archive-resume phase re-running archive-move/commit/publication; injected-crash tests at all four boundaries named in supporting criteria and test plan.
4. **Entity lock does not establish concurrency safety (blocker). Answer: accept-and-revise.** Verified the lock is gates-writers-only (`internal/gates/operation.go:382-397`) and `status --set` rewrites without it (`internal/status/mutate.go:71-193`). Resolution: one shared entity lock + full-document compare-and-swap (the gates-writer CAS, `internal/gates/io.go:188-226`) honored by sentinel writes, regress, and finalize; one shared status-layer validation boundary; sentinel-vs-regress and finalize-vs-regress race tests require one coherent winner or a clean refusal.
5. **`delivering` alone is not a restart-safe binding (material). Answer: accept-and-revise.** Verified hook discovery is a live rescan (`internal/status/mutate.go:514-551`) and boot recognizes three readiness strings (`internal/status/format.go:78-88`). Resolution: consume writes `delivering` + `mod-block: merge:{hook}` atomically (arm inside the consume write — no crash-before-arm window); multi-hook registration refuses closed; missing/renamed hook fails closed naming the recorded hook; new `delivering-awaiting-merge` readiness classification mapped in boot/status.
6. **`reversed` does not guarantee a recorded regression (material). Answer: accept-and-revise.** Resolution: reversal subtree required nonblank iff `state==reversed`, forbidden otherwise; `--reason` and `--by` required flags (the CLI gains an actor argument); `ref` derived from current `pr`, mismatched `--pr` is an error; `validateApplication` enforces.
7. **Any nonterminal `--to` is not safe (material). Answer: accept-and-revise.** Verified `--to backlog` (a gate stage) suppresses dispatch (`internal/status/format.go:214-221`). Resolution: default restricted to the record stage's configured `feedback-to`; overrides limited to defined nonterminal non-gated working stages, authorized by captain direction recorded via `--by`; same-stage refused; archived entities pinned read-only through the existing resolver (`internal/status/native_runner.go:530-556`).
8. **AC-1/AC-2 fixtures did not prove the hole closed (material). Answer: accept-and-revise.** Resolution: AC-1 starts from a real pending approval and drives CLI `gate consume` with an actually registered hook; AC-2 exercises every supported terminal-transition producer; both name adversarial discriminators (old consume branch left intact turns them red); the dispatch test is replaced with a discriminating readiness-classification assertion.
9. **Live smoke wrongly waived (material). Answer: accept-and-revise.** Verified `docs/dev/README.md:82-83,198` requirements. Resolution: test plan gains the detached adversarial audit (status mutation/guard + shipped contract triggers), the applicable live lanes per path→lane mapping, and a durable live scenario covering consume→deliver and consume→regress→re-enter→deliver.
10. **Named documentation surfaces incomplete (material). Answer: accept-and-revise.** Verified `mods/pr-merge.md` (shipped) and `docs/dev/_mods/pr-merge.md:18-25` (workflow-local, stale) are distinct files, plus the cited site pages and schema. Resolution: `docs/dev/_mods/pr-merge.md`, `docs/site/reference/command-reference.md`, `docs/site/reference/frontmatter-contract.md`, `docs/site/concepts/gates-and-decisions.md`, and `docs/schema/entity.mdschema.yml:107-118` added to Surfaces, expected surface, and the doc diff proposal with concrete wording deltas.
11. **Spike did not exercise the riskiest integration, and left no artifact (material). Answer: partially accept.** Accept: the round-trip spike covered guarded record mutation only, and reverting it uncommitted left no reviewable artifact — the design now owes a retained, committed implementation-opening spike of the full consume→arm→sentinel→atomic-finalize + regress-race transaction (test plan, spike section). Defense (bounded): the original "no spike needed for the remaining mechanisms" claim was scoped too broadly but not baseless — the cited mechanisms (co-write, scanMods, finalize envelope) do ship; the correction is that their novel cross-package *composition* is what needed spiking, now scheduled.

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

- DONE (staff review round 1): All 11 findings answered — 10 accept-and-revise applied to the design, 1 partial accept (finding 11) with a bounded defense — each logged in the `### Staff Review` section with disagreement, severity, answer, resolution, and reviewer provenance (openai-codex/gpt-5.6-sol:xhigh, run 5da3cbb1).
  Reviewer mechanism citations verified against the code before answering (`application.go:176-180`, `handlers.go:210-218`, `merge_guard_test.go:69-88`, `mutate.go:71-193,514-551`, `merge.go:374-421`, `operation.go:382-397`, `io.go:188-226`, `format.go:78-88,214-221`, `native_runner.go:530-556`, live legacy entity shape, both pr-merge mod files).
- DONE (staff review round 1): Blockers resolved in the design: (1) hook-registration classification replaces merge-policy/hook-free wording with AC-2 narrowed and AC-3 rewritten (+ registered-local `delivering` test); (2) guarded legacy-adoption path shapes the fo-boot entity (open/closed/merged-sentinel fixtures, authority history preserved); (4) shared entity lock + full-document CAS boundary across sentinel writes, regress, and finalize with race tests.
  Design deltas applied in Forward direction (atomic consume-side `mod-block` arm, multi-hook refusal, `delivering-awaiting-merge` readiness), Legacy adoption + Shared lock/CAS sections, Backward direction (required `--by`/`--reason`, reversal contract, restricted `--to`), Boundary table rows, Surfaces/expected surface/doc diff rows (schema, site docs, workflow-local mod).
- DONE (staff review round 1): Material findings folded into ACs and test plan: AC-1 starts from a real pending approval via CLI consume with an actual registered hook; AC-2 exercises every terminal-transition producer with adversarial discriminators; discriminating dispatch test; crash-injection resume tests; reversal provenance contract; route guards; detached adversarial audit + applicable live lanes + durable consume→deliver and consume→regress→re-enter→deliver live scenario added per `docs/dev/README.md:82-83,198`; retained committed implementation-opening cross-package spike scheduled.
  Body coherence re-checked after edits (no orphaned `hook-free`/`merge: local` shortcuts, optional-reversal or legacy-free wording remaining).

### Summary

Ideation delivered a full design for the consume-into-terminal semantic hole: approval authority is bound (new `delivering` state), not spent, when a merge hook defers delivery; terminal status, verdict, and completed land atomically only at merge-sentinel finalize; delivery failure routes back through a recorded `merge regress` verb into `feedback-to` without surgery or a fiction of re-bought authority. The riskiest mechanism (guarded gate-record round-trip without corrupting `gates.current`) was exercised in a PASS spike; remaining mechanisms are proven shipping code. The design awaits the ideation gate and the captain's verdict on the direction per his 2026-07-31 resolution.
