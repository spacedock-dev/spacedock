# Durable decisions

**Sprint:** the entities matching `sprint: durable-decisions` — list with `spacedock status --workflow-dir docs/dev --where sprint=durable-decisions`. Membership and per-task state are the query, never enumerated here.
**Target train:** stable **0.27.0** — bound at scope-lock 2026-07-21; movable by captain decision without re-carve (this line is the only place the train lives).
**Theme:** the decision the captain makes becomes the state the system holds. 0260 shipped the gate-record notation as a hand-run convention and collected its failure evidence in one day of production use: a self-conflicting attempt pointer, stale applications, digests no committed tree reproduces, `--set` re-serialization breaking hand edits, results destroyed on the untested path. This sprint mechanizes the proven shape — a recorder that owns every gates write, a presentation command that cannot lose a result, eligibility only where live need exists, and the finding-decline reframed onto the same record model.
**Scope-lock (captain, 2026-07-21):** the recorder group ONLY. The panel, provenance, router, mining, and estate candidates surfaced at lock were all held out — this train is one cluster, sequenced by one dependency.
**Evidence:** `docs/dev/.spacedock-state/durable-gate-approval-pending-blockers/production-evidence-2026-07-20-fo-dry-run.md` (8 findings), the 0260 closure-pass findings (advisory-digest hole, pointer conflict, uncommitted seat), and `_debriefs/2026-07-20-01-0260-shaping.md` float findings 1-13.

## Goal (success criterion)

A captain's gate decision, once made, is recorded by machinery that cannot mis-file it, presented by machinery that cannot lose it, and consumed exactly once — proven by checks that can fail:

- Replaying 0260's real hand-authored gate history through the recorder reproduces every record; the recorder REFUSES the record-integrity failure fixtures 0260 produced live (a current-attempt pointer conflict — z7's real one is the red fixture; a mutation of a frozen closure), and `--set` on unrelated fields no longer touches the gates block (3k). The second-pending-application refusal follows the application layer to h1 (resolution-first split, captain 2026-07-21).
- A recorded briefing digest is mechanically verifiable: the drift check reproduces the digest from a committed snapshot, closing the advisory-digest hole (3k).
- Gate presentation is an overridable channel of the present-gate skill — default chat, override the hardened float script — with THIS repo's deliverable being the channel prose, the recorder-side validation verbs, and the subspace-free binary criterion (xb). The override script and its committed drive suite are a NAMED CROSS-REPO RELEASE CONDITION: the 0.27.0 pre-cut gates on a pinned subspace revision carrying that suite, owned by the subspace workflow — captain ruling 2026-07-21; this repo neither absorbs nor silently depends on it.
- Blockers/eligibility ship only against a demonstrated live consumer; absent one, h1 closes as a recorded decline, not a build (h1 — its own body carries this condition). "Consumed exactly once" means the AUTHORIZATION: the application is spent atomically with the status transition, provably once; the dispatch EFFECT is the dispatch machinery's, documented at-least-once retryable (receipts stay declined per the boundary alignment) — captain ruling 2026-07-21 from the codex seat's exactly-once fork; the two crash windows get authorization-side fixtures that surface rather than double-fire.
- The DoD line moved from 0260: a seeded correct-but-disproportionate finding produces a recorded decline — as the ensign's advisory resolution — and a zero-line diff in live replay (02av reframe). Honest mechanism statement (captain-approved restatement, 2026-07-21, closing the codex seat's fifth material finding): this proves the TRIAGE SEMANTICS on convention-recorded rounds — the record is hand-authored into the room per the contract's explicit interim; MACHINE recording of advisory rounds is the deferred rounds-generalization follow-up, not this train's claim.

## Constraints

- The 0260 captain rulings bind unchanged: declared expected surface + tolerance (default 2× unless the entity declares otherwise); no minted terminology, bare ordinals; a new check/enforcement process needs explicit captain approval and normally its own entity; prose-greps are one-off validation evidence, never committed tests.
- The 0260 operating directives bind every FO session driving this sprint: assume the members' intended behavior in FO conduct; present gates via the Subspace float with the design in the presented artifact (probe-first ritual per the shaping debrief); record resolutions in gates frontmatter — by hand only until 3k lands, then recorder-owned.
- A scripted fan-out declares expected agent count, tolerance, and economic reasonableness before launch (z7's authoring-time amendment, in force for FO conduct now).
- This sprint eats its own output at the first opportunity: once the recorder can record, this sprint's own remaining gates use it.
- **Digest domains (captain ruling, 2026-07-21, from the codex seat):** the contract names two digest domains — the canonical-bytes briefing digest the recorder emits and validates going forward, and the raw-file pin. Every gate record written during shaping is EXPLICITLY a raw-file pin (honest history, no rewrite, no silent reinterpretation); the recorder accepts the marked legacy domain and one formatting-only fixture proves the domains diverge.
- **Recording identity (captain ruling, 2026-07-21):** a resolution is recorded under the identity that actually rendered it. A chat-directed fold closes under the FO's identity acting on delegated authority, with the reason quoting the captain's directive; the captain's identity appears only on resolutions the captain rendered over content the captain saw (a float resolved, presented content approved in chat). One sentence lands in the contract's lifecycle rules at implementation; prior records stand as honest history.

## Sequencing

1. **3k's gate is the opening event** — rebind the open attempt's briefing to the post-cut content, float, close. Everything designs against its approved contract.
2. **3k implementation leads the build**; xb and h1 ideations may run in parallel once 3k's gate closes (they design against the contract, not the binary).
3. **02av's reframe ideation dispatches only after 3k's gate closes** — its brief, boundary question, and design inputs are recorded in the entity (`## Reframe brief`). It must answer the plumbing-boundary question rather than absorb it.
4. **xb is cross-repo sequenced** with the subspace-tui briefing-package surface; the working-copy-skill ritual is its interim and its evidence base.

## Out of scope

- The panel member (7wv), provenance enforcement, the AGENTS.md router + stale Priorities, sprint-close mining, and the estate items (1az, 2kz, the four contractlint retirements) — all remain unlabeled backlog; held out at scope-lock as too heavy for this train.
- bw's deferred record command — dissolved: the advisory-records direction supersedes it; 02av's reframe decides the plumbing boundary.
- The stakes member — parked, does not revive here.
- The staff-review non-dispatching clarification (subspace shame log ask) — surfaced at lock, not ratified; stays with its own repo's report.
- Direct edits to sibling repos; subspace-tui product work (xb consumes its surface, does not build it).

## Lifecycle checklist

**Shape — Shaping FO**
- [x] **Scope-lock** with the captain — locked 2026-07-21: recorder group only (3k, xb, h1, 02av)
- [x] **Carve** — 4 members stamped, index written; drivable query above
- [x] **Ideate** — complete 2026-07-21: the recorder through eight attempts (two captain scope cuts, the fold, the render war); xb through three (the fallback question, the decoupling reframe); h1 with its recorded decline; 02av's reframe in one round
- [x] **⚠️ Preflight staff review (sprint-wide)** — fable seat 2026-07-21, READY AFTER FOLDS (6 material, 1 needs-decision, 5 recorded declines) → `staff-review.md`; all folds applied and verified, both captain rulings given; closure complete
- [x] **Present ideation gates** — all four closed-approve with pending advances, every briefing byte-verifiable in its room; two design questions answered mid-gate (the no-subspace fallback; the decoupling reframe)
- [x] **Package** — `dispatch-sprint-execution.md`, written after the preflight resolved (the 0260 sequencing lesson, kept)

**Drive — Commander (separate cold-booted session)**
- [ ] Implementation → validation → done per member; detached adversarial audit on shipped-contract surfaces
- [ ] **Contract landing pass** (captain placement ruling, 2026-07-21): strip the spec's shaping scaffolding — owner-tag lines, diagram task-id prefixes (converted to component words with a render re-check via a float), example ids genericized — so the landed spec speaks only component terms. Owner: the Commander, as the recorder member's final step before its merge.
- [ ] **⚠️ Pre-cut audit** before the tag
- [ ] **Cut 0.27.0** — `go test ./...` + `-race` green, `gofmt` clean, then `docs/releasing.md` *(captain authorizes)*

**Close — Shaping FO**
- [ ] Sprint-close mining pass (by hand, second manual run; the recipe is banked in the 0260 forensics) → seeds the next train

## Responsibility boundary (captain-aligned, 2026-07-21)

| Task | Owns | Boundary line |
|---|---|---|
| 3k recorder | what the decision IS: gate → attempt → briefing binding (replaceable open, frozen closed) → exact resolution; record invariants; snapshot digests; status surfacing of recorded state; the contract doc (sections owner-tagged) | recording never advances status, never dispatches, never computes eligibility |
| h1 applications + blockers + eligibility | what the decision DOES: the one-use application (pending → consumed exactly-once via the existing transition path; superseded on drift; not-applicable on hold), declared blockers (fail-closed), execution holds, the eligibility computation | an application never exists without a closed binding approval; a resolution stands alone fine |
| xb presentation command | how the decision is OBTAINED: package validation, blocking TUI child, atomic retention on success AND failure, presenter lifecycle, probe ritual, provider id-mapping (specified in 3k's contract) | hands the validated resolution to the recorder; never writes gates itself |
| 02av triage on advisory records | what a non-gate decision MEANS: triage rule, decline as the ensign's advisory resolution, AC-narrowing graduating to a binding gate attempt | round records carry no application, ever — advisory cannot advance anything by construction |

Rules: (1) 3k records, h1 authorizes — resolutions stand alone, applications depend; (2) only binding resolutions may carry applications; (3) one write surface — the recorder binary, extended by h1's verbs, called by xb, never a second writer; (4) one spec — 3k's contract doc, sections owner-tagged, each task edits only its sections; (5) inside h1, applications have a proven consumer (the 0260 Commander) while blockers/holds must be justified or declined separately.

Dogfooding change protocol: frictions in your own spec section — amend and record the round in your entity; frictions in another owner's section — route through the FO to the owner, consumers re-anchor on landed text; amendment after a closed gate — a superseding attempt on the owner's gate record; every friction is a finding triaged under the standing taxonomy (declines legal). The first dogfood channel is this sprint's own remaining gates.
