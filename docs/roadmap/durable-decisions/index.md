# Durable decisions

**Sprint:** the entities matching `sprint: durable-decisions` — list with `spacedock status --workflow-dir docs/dev --where sprint=durable-decisions`. Membership and per-task state are the query, never enumerated here.
**Target train:** stable **0.27.0** — bound at scope-lock 2026-07-21; movable by captain decision without re-carve (this line is the only place the train lives).
**Theme:** the decision the captain makes becomes the state the system holds. 0260 shipped the gate-record notation as a hand-run convention and collected its failure evidence in one day of production use: a self-conflicting attempt pointer, stale applications, digests no committed tree reproduces, `--set` re-serialization breaking hand edits, results destroyed on the untested path. This sprint mechanizes the proven shape — a recorder that owns every gates write, a presentation command that cannot lose a result, eligibility only where live need exists, and the finding-decline reframed onto the same record model.
**Scope-lock (captain, 2026-07-21; amended 2026-07-22, 2026-08-01, and 2026-08-03):** the recorder group ONLY in the original build. Captain amendment added `vn`, the folder-form state-commit persistence boundary, after live 3k dogfood proved that a recorded gate package otherwise requires manual Git commits. The 2026-08-01 pre-stable necessity audit adds only the cuts required to make the unreleased v1 surface smaller and internally operable: terminal authority is spent only with delivery proof; preparation is limited to actionable gates; speculative application and gate-state fields are removed or justified; the workflow-specific round recorder is cut; provider-backed closure is proved on the exact candidate or cut; minor help and docs reconcile last. The 2026-08-03 amendment adds `sk`, `ej9`, and `g3`. These tasks make completed gates discoverable, make event-loop stops truthful, and define moving-target conflict ownership. These amendments add no compatibility or migration obligation. The panel, provenance, router, mining, and estate candidates remain held out.
**Evidence:** `docs/dev/.spacedock-state/durable-gate-approval-pending-blockers/production-evidence-2026-07-20-fo-dry-run.md` (8 findings), the 0260 closure-pass findings (advisory-digest hole, pointer conflict, uncommitted seat), and `_debriefs/2026-07-20-01-0260-shaping.md` float findings 1-13.

## Goal (success criterion)

A captain's gate decision, once made, is recorded by machinery that cannot mis-file it, presented by machinery that cannot lose it, and consumed exactly once — proven by checks that can fail:

- The recorder accepts and emits one canonical unreleased-v1 gates projection; pilot encodings fail closed without mutation and have no compatibility or migration path. Canonical lifecycle, frozen closure/application, cross-logical-gate re-entry, complete Briefing-derived Result association, CAS/locking/atomic replacement, and byte-exact preservation outside `gates` are independently exercised (3k). The second-pending-application refusal follows the application layer to h1.
- A recorded briefing digest is mechanically verifiable: the drift check reproduces the digest from a committed snapshot, closing the advisory-digest hole (3k).
- A folder-form gate record is durably one entity commit unit: `state commit <slug>` includes the index and every new or changed room artifact beneath that entity, excludes sibling dirt, and cannot report a clean no-op while package files remain untracked (vn).
- Gate presentation is an overridable channel of the present-gate skill — default chat, override the hardened float script — with THIS repo's deliverable being the channel prose, the recorder-side validation verbs, and the subspace-free binary criterion (xb). The override script and its committed drive suite are a NAMED CROSS-REPO RELEASE CONDITION: the 0.27.0 pre-cut gates on a pinned subspace revision carrying that suite, owned by the subspace workflow — captain ruling 2026-07-21; this repo neither absorbs nor silently depends on it.
- Blockers/eligibility ship only against a demonstrated live consumer; absent one, h1 closes as a recorded decline, not a build (h1 — its own body carries this condition). "Consumed exactly once" means the AUTHORIZATION: the application is spent atomically with the status transition, provably once; the dispatch EFFECT is the dispatch machinery's, documented at-least-once retryable (receipts stay declined per the boundary alignment) — captain ruling 2026-07-21 from the codex seat's exactly-once fork; the two crash windows get authorization-side fixtures that surface rather than double-fire.
- The DoD line moved from 0260: a seeded correct-but-disproportionate finding produces a recorded decline — as the ensign's advisory resolution — and a zero-line diff in live replay (02av reframe). Honest mechanism statement (captain-approved restatement, 2026-07-21, closing the codex seat's fifth material finding): this proves the TRIAGE SEMANTICS on convention-recorded rounds — the record is hand-authored into the room per the contract's explicit interim; MACHINE recording of advisory rounds is the deferred rounds-generalization follow-up, not this train's claim.
- Stable v1 exposes only machinery whose necessity has a demonstrated journey. Chat approval works without a presentation provider; a terminal approval remains pending until delivery proves the spend; invalid-stage preparation is byte-clean; open stale preparation has a truthful withdrawal; unused application and gate-state fields are removed or independently justified; workflow-specific advisory-round policy does not ship as a generic gate primitive; provider-backed closure remains only with an exact-candidate proof.

## Pre-stable necessity cut (captain, 2026-08-01)

The membership query remains authoritative. These names record responsibility boundaries, not a second roster:

| Boundary | Owner | Required outcome |
|---|---|---|
| terminal approval and delivery | `resolution-consume-terminal-before-delivery` (`1w6`) | terminal consume routes without spending; delivery success spends and terminalizes atomically; delivery rework supersedes and routes |
| stale open preparation | `withdraw-stale-open-gate-attempt` (`0m6`) | withdraw without a fabricated Resolution, then prepare a normal successor |
| invalid-stage preparation | `reject-gate-prepare-outside-actionable-stage` (`hq3`) | refuse before entity or room mutation unless the current stage is an actionable gate |
| application schema | `minimize-v1-gate-application-schema` (`nth`) | retain only state with a supported producer and consumer; no compatibility decoder |
| gate-state schema | `simplify-gate-state-v1-schema` (`jcc`) | remove prototype digest compatibility and derive or minimally justify stored gate selection |
| completed-gate discovery | `gate-agent-ergonomics` (`sk`) | derive the gate-preparation action from durable evidence; add no stored scheduler state |
| event-loop order | `make-fo-event-loop-ordering-explicit` (`ej9`) | process merge, gate, dispatch, and wait actions in one explicit order; stop only after the truthful idle check |
| moving-target conflicts | `define-fo-moving-target-conflict-ownership` (`g3`) | preserve authority and exact evidence; require Captain-owned reconciliation and fresh validation |
| advisory rounds | `cut-workflow-specific-round-recorder-from-v1` (`wjk`) | remove the development-policy round surface from stable v1; redesign later as workflow-neutral |
| provider closure | `prove-or-cut-provider-backed-gate-closure` (`a73`) | one pinned exact-candidate transaction or no public provider-backed closure in v1 |
| help and docs | `polish-v1-gate-command-surface` (`f6c`) | reconcile discoverability and prose only after the semantic owners land; execute the re-homed Contract landing pass — strip owner tags and genericize example ids in the contract spec, with a render re-check (captain re-homing, 2026-08-03) |

No cleanup ticket may absorb a semantic change from the rows above. No semantic owner may preserve a prototype shape merely because current pilot metadata contains it; a one-off state transformation is cheaper and more truthful than a compatibility layer before v1 release.

## Remainder exit contract (captain, 2026-08-03)

The remainder of this sprint is the pre-stable necessity cut, and the cut table above is its definition of done. The sprint's theme applied to itself: the captain's cut decisions become the state the spec holds.

- **Done when the table closes.** Every boundary row ends in exactly one of two states: its required outcome landed and verified by merge evidence, or its surface cut. Row state lives in the membership query and merge history, never in this file; a row is never closed by assumption or an inherited label.
- **The spec is the surviving contract.** At close, `docs/specs/gate-resolution-frontmatter-contract.md` describes exactly the machinery that shipped: cut surfaces recorded only under "Explicitly outside v1", no pending sections, no shaping scaffolding. Every semantic owner's PR carries its own spec-section edit (the one-spec rule); the scaffolding strip belongs to `f6c` via the re-homed Contract landing pass.
- **Anything not a row is not in the sprint.** A member carrying `sprint: durable-decisions` without a boundary row either gains one by captain decision or is cut from the sprint; `bv` and `47g` are decided under this rule now, not queued (`47g` was already flagged cut-or-reshape). New discoveries during the remainder file as next-train seeds, never as members.
- **Cross-sprint release edge:** the live-test-truth sprint's deferred `rm` (restore-quarantined-common-live-journeys) and its `v0.27.0` DoD wait on this sprint's remaining landings; a slip here blocks that cut or moves the shared train — a captain decision visible from both indexes.

## Streamlined common journeys

### Chat approval into a nonterminal stage

1. The First Officer runs `gate prepare`, then `state commit` so the canonical Briefing and request are durable.
2. The First Officer presents the Briefing in chat and records the decision with `gate record --decision ... --actor ...`, then commits it. The actor is the renderer identity required by the recording-identity ruling: `agent:first-officer` for an FO-rendered delegated close, or `person:captain` for a Captain-rendered decision over content the Captain saw.
3. On approval, `gate consume` atomically spends the pending application with the next-stage transition; the First Officer commits and dispatches that stage.

The agent constructs no JSON and supplies no output paths. Chat is the default and requires no Subspace installation.

### Chat approval into the terminal stage

1. Preparation, presentation, recording, and state commits are identical to the nonterminal journey.
2. `gate consume` leaves the terminal application pending and returns `approved-awaiting-merge`; it does not write terminal status.
3. The merge guard proves delivery. Success atomically writes pending to consumed together with terminal status, verdict, and completion. Retryable trouble changes none of them. Rework atomically writes pending to superseded, clears delivery state, and routes to the declared feedback stage.
4. A later validation uses a fresh gate attempt and fresh approval; superseded authority is never re-spent.

### Review input changes before a decision

1. The First Officer runs `gate withdraw --reason ...` against the open prepared attempt and commits the withdrawal.
2. `gate prepare` creates the ordinary successor attempt and room.
3. No `hold` Resolution is fabricated for content that was never decided.

### Invalid preparation request

At an ungated or terminal current stage, `gate prepare` exits nonzero before allocating an attempt or writing a room. The operator advances or completes the real workflow stage; there is nothing to withdraw or repair.

### Hold and revise

A hold records the Resolution and stops. A revise records the Resolution and routes through the workflow's declared feedback behavior. Neither case carries application metadata unless a demonstrated consumer actually applies it.

### Provider recording cut

Stable v1 permits chat or Subspace to present the committed gate. Both return semantic
decision and reason input to the First Officer, who records it through the same
`gate record --decision` path. Stable v1 does not ship `gate record --room`, Result or
inventory ingestion, or another recorder. The prepared room continues to bind the
canonical Briefing used by validation and one-use consumption.

## Constraints

- The 0260 captain rulings bind unchanged: declared expected surface + tolerance (default 2× unless the entity declares otherwise); no minted terminology, bare ordinals; a new check/enforcement process needs explicit captain approval and normally its own entity; prose-greps are one-off validation evidence, never committed tests.
- The 0260 operating directives bind every FO session driving this sprint: assume the members' intended behavior in FO conduct; present gates via the Subspace float with the design in the presented artifact (probe-first ritual per the shaping debrief); record resolutions in gates frontmatter — by hand only until 3k lands, then recorder-owned.
- A scripted fan-out declares expected agent count, tolerance, and economic reasonableness before launch (z7's authoring-time amendment, in force for FO conduct now).
- This sprint eats its own output at the first opportunity: once the recorder can record, this sprint's own remaining gates use it.
- **Digest domains (captain ruling, 2026-07-21, from the codex seat):** the contract names two explicit v1 digest domains — the canonical-bytes Briefing digest the recorder emits and uses for Result association, and the raw-file pin. A marked raw-file pin is never silently reinterpreted; one formatting-only fixture proves the domains diverge.
- **Recording identity (Captain-approved Cycle-31 provenance retirement; amended 2026-07-25):** preserve renderer identity. An FO-rendered delegated close records `agent:first-officer` with a nonblank FO evidence judgment as reason, never quoted grant/directive text; chat grant text is neither authenticated nor retained. `person:captain` remains only for a Captain-rendered decision over content the Captain saw. The [Gate Resolution frontmatter contract](../../specs/gate-resolution-frontmatter-contract.md) is canonical; prior records remain honest history.

## Sequencing

1. Land the already-proven core corrections first: `0m6` for truthful withdrawal, `1w6` for terminal authority/delivery, `zbc` for post-rework binding freshness, and `kd` for the dispatch-pinned launcher used by ensign-owned workflow calls.
2. Ideate the independent v1 cuts in parallel: `hq3` invalid-stage preparation, `nth` application schema, `jcc` gate-state schema, and `wjk` removal of workflow-specific round recording. Their implementations may overlap in `internal/gates`; merge order must be declared at their ideation gates rather than discovered through conflict.
3. Land the dependencies of `sk`, then implement `sk`. Its approved design requires landed `s4`, `gqs`, and `0m6` behavior.
4. Ideate `ej9` and `g3` in parallel. Implement `ej9` after `sk` fixes completed-gate discovery. `g3` owns conflict recovery and fresh-evidence rules.
5. Re-run the chat-default journey after those cuts. It is the minimum stable value and does not wait for a provider.
6. Run `a73` only on that settled candidate. Retain provider-backed closure on one pinned end-to-end pass. Otherwise, cut it without delaying chat.
7. Run `f6c` last. It reconciles help and documentation with landed behavior and must not invent semantics.
8. Run the pre-cut audit on the exact candidate. Run the full tests, race tests, formatting, state check, and smallest real chat journey. Provider evidence is required only if `a73` retained that surface.

## Out of scope

- The panel member (7wv), provenance enforcement, the AGENTS.md router + stale Priorities, sprint-close mining, and the estate items (1az, 2kz, the four contractlint retirements) — all remain unlabeled backlog; held out at scope-lock as too heavy for this train.
- bw's deferred record command — dissolved: the advisory-records direction supersedes it; 02av's reframe decides the plumbing boundary.
- The stakes member — parked, does not revive here.
- The staff-review non-dispatching clarification (subspace shame log ask) — surfaced at lock, not ratified; stays with its own repo's report.
- Direct edits to sibling repos; subspace-tui product work (xb consumes its surface, does not build it).
- Reintroducing a generic advisory-round recorder in this cut. `workflow-neutral-advisory-round-recorder` retains that later design problem after `wjk` removes the development-policy implementation.
- Multiple-Artifact gate preparation and Subspace convenience work beyond the minimal room-only provider proof. Those remain independently filed and cannot block the chat-default stable value.

## Lifecycle checklist

**Shape — Shaping FO**
- [x] **Scope-lock** with the captain — locked 2026-07-21 to the recorder group; amended 2026-07-22 to add vn's persistence boundary after live dogfood
- [x] **Carve** — 5 members stamped; membership remains the drivable query above
- [ ] **Ideate** — the original four completed 2026-07-21; vn's existing ideation is ready for its sprint gate
- [x] **⚠️ Preflight staff review (sprint-wide)** — fable seat 2026-07-21, READY AFTER FOLDS (6 material, 1 needs-decision, 5 recorded declines) → `staff-review.md`; all folds applied and verified, both captain rulings given; closure complete
- [x] **Present ideation gates** — all four closed-approve with pending advances, every briefing byte-verifiable in its room; two design questions answered mid-gate (the no-subspace fallback; the decoupling reframe)
- [x] **Package** — `dispatch-sprint-execution.md`, written after the preflight resolved (the 0260 sequencing lesson, kept)

**Drive — Commander (separate cold-booted session)**
- [ ] Implementation → validation → done per member; detached adversarial audit on shipped-contract surfaces
- [ ] **Contract landing pass** (captain placement ruling, 2026-07-21): strip the spec's shaping scaffolding — owner-tag lines, diagram task-id prefixes (converted to component words with a render re-check via a float), example ids genericized — so the landed spec speaks only component terms. Owner: the Commander, as the recorder member's final step before its merge. **Re-homed 2026-08-03:** the recorder merged without the pass; ownership moves to `f6c` (its cut-table row carries it) and this checkbox closes with `f6c`.
- [x] **⚠️ Pre-cut necessity audit** — 2026-08-01 found the responsibility boundaries above; stable readiness remains false until their sprint owners close or their surfaces are cut
- [ ] **⚠️ Exact-candidate verification** before the tag — chat journey required; provider journey conditional on `a73` retention
- [ ] **Cut 0.27.0** — `go test ./...` + `-race` green, `gofmt` clean, then `docs/releasing.md` *(captain authorizes)*

**Close — Shaping FO**
- [ ] Sprint-close mining pass (by hand, second manual run; the recipe is banked in the 0260 forensics) → seeds the next train

## Responsibility boundary (captain-aligned, 2026-07-21)

| Task | Owns | Boundary line |
|---|---|---|
| 3k recorder | what the decision IS: gate → attempt → briefing binding (replaceable open, frozen closed) → exact resolution; record invariants; snapshot digests; status surfacing of recorded state; the contract doc (sections owner-tagged) | recording never advances status, never dispatches, never computes eligibility |
| vn folder-form state commit | how one entity's index, Briefing rooms, reports, exact Results, and associations become one durable state commit without sweeping sibling dirt | commits workflow evidence only; never interprets gates or changes recorder semantics |
| h1 applications + blockers + eligibility | what the decision DOES: the one-use application (pending → consumed exactly-once via the existing transition path; superseded on drift; not-applicable on hold), declared blockers (fail-closed), execution holds, the eligibility computation | an application never exists without a closed binding approval; a resolution stands alone fine |
| xb presentation command | how the decision is OBTAINED: package validation, blocking TUI child, atomic retention on success AND failure, presenter lifecycle, probe ritual, provider id-mapping (specified in 3k's contract) | hands the validated resolution to the recorder; never writes gates itself |
| 02av triage on advisory records | what a non-gate decision MEANS: triage rule, decline as the ensign's advisory resolution, AC-narrowing graduating to a binding gate attempt | round records carry no application, ever — advisory cannot advance anything by construction |

Rules: (1) 3k records, h1 authorizes — resolutions stand alone, applications depend; (2) only binding resolutions may carry applications; (3) one write surface — the recorder binary, extended by h1's verbs, called by xb, never a second writer; (4) one spec — 3k's contract doc, sections owner-tagged, each task edits only its sections; (5) inside h1, applications have a proven consumer (the 0260 Commander) while blockers/holds must be justified or declined separately.

Dogfooding change protocol: frictions in your own spec section — amend and record the round in your entity; frictions in another owner's section — route through the FO to the owner, consumers re-anchor on landed text; amendment after a closed gate — a superseding attempt on the owner's gate record; every friction is a finding triaged under the standing taxonomy (declines legal). The first dogfood channel is this sprint's own remaining gates.
