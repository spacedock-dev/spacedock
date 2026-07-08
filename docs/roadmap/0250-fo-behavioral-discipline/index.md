# 0250 — FO behavioral discipline (0.25.0)

**Sprint:** the entities matching `sprint: 0250-fo-behavioral-discipline` — list current members with `spacedock status --workflow-dir docs/dev --where sprint=0250-fo-behavioral-discipline`. Membership and per-task state are the query, never enumerated or tracked in this doc.
**Theme:** a captain can leave the FO to run itself between deliberate touches — boot into a light greet, trigger work with an explicit `engage` verb, and trust the FO to self-triage correctly while unattended. The flagship (`k74g`) makes the boot light and gives the captain the `engage` verb; the discipline cluster (`z25`/`zm`/`vcm`) is what makes `engage` SAFE to lean on — an FO that holds its own evidence bar, does the right-sized thing, and keeps moving without stalling or busywork.

## Goal (success criterion)

An interactive FO boot is LIGHT and ceremonial-free, and the captain drives with an explicit `engage` verb:

- The greet lists the managed workflow(s) and hints "Use engage" — it renders NO gate review, and when multiple workflows are discovered it does NOT ask the captain to pick one.
- `«engage»(workflow)` is a captain-invoked FO interaction verb whose effect is running the EXISTING `«dispatch.next-action»()` event-loop skeleton (dispatch a ready entity / present a ready gate / advance a completed non-gated stage) to its stopping condition for the named workflow — no new binary command, no new dispatch mechanism.
- The launcher bootstrap prompt no longer ends with the throwaway "Engage." flourish that collides with the verb.

And the FO is trustworthy to leave running between engages, because it holds its own decisions to the bar it imposes on workers:

- It holds its own gate/merge/triage decisions to "required verification follows from what changed" — no lane waived as "unrelated" by intuition, no red labeled "known flake" without reading this run (`z25`).
- It reaches for the smallest sufficient mechanism — refusing BOTH over-orchestration (a workflow/worker/PR where a direct edit/commit is right) AND gate-time re-verification busywork (re-running stage-owned verification the validator already produced) (`zm`).
- It keeps moving — a gate approval triggers the next action, independent work dispatches in parallel, an async launch doesn't end the turn while independent work remains, and a captain correction narrows scope without halting the session (`vcm`).

Every behavioral claim is proven by a LIVE FO drive observing the behavior against a baseline that moved the wrong way — never a prose-grep over the contract the change writes.

## Why

**The boot is heavy AND ceremonial.** This session's interactive boot ran ~5 minutes: the contract (Startup step 8) directs the greet to render a full `present-gate` review — stage-report read + AC cross-check — for every ready gate before stopping, so greet cost scales with the number of ready gates; step 3 asks the captain to pick when multiple workflows are discovered; and the launcher bootstrap prompt (`internal/cli/frontdoor.go`) ends with a throwaway "Engage." flourish. A captain who just wants to see what's managed and then explicitly trigger a sweep can't — the boot front-loads assembly work before there is any per-entity direction, and "Engage." as a flourish now collides with the real verb. `k74g` makes the greet a light manifest + an `engage` hint, and defines `«engage»(workflow)` as a captain-invoked entry into the event-loop logic the contract ALREADY defines (`«dispatch.next-action»()`, fo-dispatch-core.md — a `→ prose` skeleton, driver binary descoped to 0222) — a contract-prose addition, no new code surface.

**The discipline cluster is what makes `engage` safe to lean on — not a separate theme.** An explicit `engage` verb only pays off if the FO can be trusted to run the sweep correctly while the captain is not watching each step. The three discipline entities are exactly the failures that make an unattended FO untrustworthy, each observed in production:

- `z25` — the FO self-exempts its OWN decisions from the evidence bar it imposes on workers (the ezf/hf incident: merged a Claude-adapter change on deterministic lanes while leaving `claude-live` unapproved; labeled a live-CI red "the known flake" without reading the actual failing test).
- `zm` — the FO climbs to a heavier mechanism than the task needs (a workflow + dispatched worker to apply edits it already held; a PR for a direct-commit roadmap doc), AND — the friction-#1 fold from this session — re-runs stage-owned verification inline (a full `go test ./...` in the dispatcher's turn to "double-check" a validation entity the validator already verified).
- `vcm` — the FO false-stops (0223): pauses after an approval to ask permission for reversible work, files independent followups serially instead of in parallel, ends its turn on an async launch with work remaining, treats a captain question as "stop the session."

An FO that self-triages correctly, does the right-sized thing, and doesn't stall is the precondition for `engage` being a verb the captain can fire and walk away from. That is why they ship together, not as separate tracks.

**The leanness constraint (supporting rationale, load-bearing).** These four all ADD text to `first-officer-shared-core.md`, and 0.23/0.24 were leanness cuts that deliberately clawed contract bytes BACK (`z25` was explicitly deferred from 0230: "Adds contract text — same opposition to the [leanness] gate"). 0.25 lands the cluster now that the lean baseline exists — but under a byte-discipline constraint: the additions land as the smallest coherent form, preferring lazy-loaded references over boot-resident lines (`zm`'s own AC-5), held to a measured ceiling vs the 0.24.0 baseline. The leanness DISCIPLINE becomes a constraint ON this sprint, not a reason to defer the cluster a third time.

## Definition of Done

`v0.25.0` ships when, merged to `main` and proven by a live interactive FO drive (not a checklist of "entity X merged"):

- **The greet is light and ceremonial-free.** A live interactive boot lists the managed workflow(s) + an "Use engage" hint, renders zero gate reviews, and does not force a workflow pick — measured against the current contract's greet, which renders a `present-gate` review per ready gate and asks which workflow when multiple are discovered.
- **`engage` works as an FO verb.** A live `«engage»(workflow)` invocation runs the `«dispatch.next-action»()` skeleton to its stopping condition for the named workflow — observed dispatching a ready entity / presenting a ready gate / advancing a completed non-gated stage — proving the verb wraps the existing logic with no new mechanism. Headless/single-entity mode is unchanged.
- **The bootstrap flourish is gone without breaking the launch.** `frontdoor.go`'s `bootstrapPrompt` and `codexBootstrapPrompt` no longer end with "Engage.", and a launched FO still boots into role (a launcher test/observation, measured against the current constants that carry the flourish). `pi.go`'s `piBootstrapPrompt` is untouched — it has no flourish.
- **The three disciplines are observed live.** A live FO drive (or drives) shows the FO holding its own gate/merge/triage to the required-verification bar (z25), choosing the smallest sufficient mechanism and refusing gate-time re-verification (zm), and keeping moving after an approval / in parallel / past an async launch / narrowing-not-halting on a correction (vcm) — each proven behaviorally, against the named production baselines, never a prose-grep.
- **The leanness gate holds.** The combined Startup + Working-Principles additions are `wc -c` within an agreed ceiling vs the 0.24.0 baseline of the affected resident files, with deferrable content in lazy references. Pin the baseline + measurement command BEFORE implementing (riskiest-first de-risk, per 0230).
- **The cut is clean.** `go test ./...` green from root; detached adversarial audits for the high-stakes surface (the shipped FO contract — every member touches it); the pre-cut antipattern audit clean; the `7v` coordination resolved (k74g's bootstrap edit sequenced before 7v, or 7v re-derives its target from the live constant); then the release ritual per [`docs/releasing.md`](../../releasing.md).

## Sequencing (the load-bearing operational constraint)

Two waves, because the code half and the contract-prose half have different collision profiles.

- **Wave 0 — the bootstrap-prompt code half, FIRST and parallel-safe.** `k74g`'s frontdoor.go edit (remove "Engage." from `bootstrapPrompt` + `codexBootstrapPrompt`) touches only Go source, not the contract-prose files — implementable in parallel with everything and independent of the Working-Principles cluster. It goes first because it unblocks the `7v` sequencing coordination (7v's byte-identity target is the codex constant this edit mutates). `pi.go`'s `piBootstrapPrompt` carries NO "Engage." flourish — untouched here; pi parity is 7v's separate job.
- **Wave 1 — the contract-prose cluster, STRICT SERIAL** (0240's "same 2–3 contract files → never parallel" precedent). All four remaining edits touch `first-officer-shared-core.md`; ideation may fan out in parallel (each writes only its own entity body), but implementation serializes to avoid merge conflicts in the same file — and, the real reason, because z25/zm/vcm compose into ONE coherent Working-Principles section that must not contradict itself. Proposed order + why:
  1. **`k74g` contract-prose** (Startup steps 3+8 + the new `«engage»` function block). It edits Startup and adds a new function block — DIFFERENT sections from the Working-Principles trio, so it is collision-light against them; land the flagship's headline value early (riskiest/most-valuable-first), and defining `«engage»` first gives the discipline principles a concrete behavioral context to bind to (the disciplines govern how the FO behaves DURING the engage loop `«engage»` invokes).
  2. **`z25`** (the evidence-bar principle in Working Principles + the `present-gate` evidence-surfacing rule). Land the foundation first: "required verification follows from what changed" is the bar zm and vcm both lean on.
  3. **`zm`** (the smallest-sufficient-mechanism gate; prefer a lazy reference per its AC-5, carrying the friction-#1 gate-time-re-verification fold). A specific application of the right-sized-effort bar z25 establishes.
  4. **`vcm`** (the four keep-moving strengthenings across Working Principles + Clarification). LAST because it is the counterweight that must be reconciled AGAINST z25's "verify what's required" and zm's "justify before climbing" — the composition tension (keep-moving vs justify-before-climbing vs verify-what-changed) resolves in vcm's final wording, so it lands into the shape the other two established.
- **Cross-member coherence pass (the sprint preflight staff review — independent, never a self-review).** Before the Wave-1 gates lock, one independent reviewer reconciles the COMBINED Working-Principles section: do z25 + zm + vcm compose without contradiction, and is the union within the leanness ceiling? This is the sprint's reason for being a sprint — driven individually, the three are colliding, potentially contradictory edits to the same paragraphs. The order above is provisional; the coherence pass may re-order.

## Leanness baseline (pinned at preflight, 2026-07-06)

Baseline (`v0.24.0` tag, confirmed byte-identical to HEAD for both files — `git diff v0.24.0..HEAD` touches neither):
- `skills/first-officer/references/first-officer-shared-core.md`: 21,663 bytes
- `skills/present-gate/SKILL.md`: 5,337 bytes

Measured resident additions as ideated (independent staff review, 2026-07-05): k74g engage block 1,433 · z25 S1 bullet 1,192 (+~350 to present-gate) · zm 733 · vcm S1–S4 2,242 → sum ≈ 5,600 bytes, +25.8% of the shared-core baseline.

**Ceiling: resident additions to `first-officer-shared-core.md` + `present-gate/SKILL.md` combined MUST NOT exceed the as-ideated total above (≈5,600 bytes / +25.8%).** This is a not-to-exceed cap, not a target to spend up to. Three trims were requested of the respective entities post-preflight as recoverable margin (not required to hit the cap, which already includes them unapplied): z25's duplicated closing sentence (~230 bytes, folds into its own lazy-loaded S2 rule), k74g's `scope:` bullet rationale (~150 bytes), and a wording-economy pass on vcm's S1–S4 (heaviest single addition at 2,242 bytes) that preserves its guard clauses — those clauses are the substance that resolved the z25/zm/vcm composition seams, not filler to cut.

Measurement command for the implementing Commander: `wc -c skills/first-officer/references/first-officer-shared-core.md skills/present-gate/SKILL.md` post-implementation; delta = post − the baseline above. Lazy-reference files (e.g. `references/fo-smallest-sufficient-mechanism.md`) are excluded from the resident count by design (zm AC-5).

## Out of scope

- **Multi-workflow `«engage»` sweep** (engaging all managed workflows at once). The `«engage»(workflow)` signature is forward-compatible by construction — a later `«engage»(w1, w2, …)` or all-managed default extends the signature — but sweeping multiple workflows is a NAMED future extension, not designed here; this entity's design must not preclude it.
- **The FO native-tooling cluster** (`3t` status --where robust/discoverable, `fk` status --read boot-lean projection, `tv` status --checklist/--ac-scan default-stage). Independent `internal/status` code fixes with no shared contract surface or value gate — driven individually through `docs/dev`, not sprint members. (They came from the same boot-friction session but cleave to a different surface.)
- **pi bootstrap-prompt parity** (`7v`) — out of this sprint; coordinate the sequencing only (Wave 0 before 7v).
- **A binary `engage` command / a `dispatch next-action` driver binary** — `«engage»` is prose wrapping the existing `→ prose` skeleton; the driver binary stays descoped (roadmap 0222). No new code surface this sprint.
- **Reopening the 0.23/0.24 leanness deferrals** (notarize, next-independent-release-line, etc.) unless they serve the discipline goal.

## Sprint lifecycle checklist (owner-tagged)

**Shape — Shaping FO**
- [x] **Scope-lock** with the captain — k74g flagship + z25/zm/vcm; native-tooling cluster driven separately *(captain, 2026-07-04)*
- [x] **Carve** — stamp `sprint: 0250-fo-behavioral-discipline` on the four members; this `index.md` created (this doc)
- [x] **Ideate** each member (k74g contract-prose + `«engage»` definition; z25/zm/vcm), riskiest mechanism first — for k74g, the live-boot/engage observability spike; check banked ideation first (z25/zm/vcm carry substantial bodies already)
- [x] **⚠️ Preflight staff review (sprint-wide, independent)** — the Working-Principles composition/coherence + leanness-ceiling review above *(2026-07-05: NEEDS MINOR REWORK → READY as a set once zm's one clause lands; see Leanness baseline section)*
- [x] **Present ideation gates** — per member; never self-approve *(captain decides)* — *(captain approved all four 2026-07-06; zm's blocking rework landed and was verified before this checkbox closed)*
- [x] **Package** — `dispatch-sprint-execution.md` (cold-boot Commander recipe)

**Drive — Commander (separate cold-booted session)**
- [ ] Implementation → validation → done per member; detached adversarial audit at validation for the shipped-contract surface
- [ ] Merge each to `main`; state commits concurrency-safe
- [ ] **⚠️ Pre-cut antipattern audit** (independent, before the tag fires)
- [ ] **Cut** `v0.25.0` per `docs/releasing.md` *(captain authorizes)*
