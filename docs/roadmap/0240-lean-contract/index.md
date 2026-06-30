# 0240 — lean contract (0.24.0)

**Sprint:** the entities matching `sprint: 0240-lean-contract` — list current members with `spacedock status --workflow-dir docs/dev --where sprint=0240-lean-contract`. Membership and per-task state are the query, never enumerated or tracked in this doc.
**Theme:** context occupancy + `«fn»` structure — defer phase-specific reference out of the boot-resident FO core and the per-dispatch ensign core, and refine the `«fn»` binding layer net-neutrally.

## Goal (success criterion)

A first officer's boot and a dispatched ensign's spawn each load only what that phase needs. The boot-resident FO core greets without carrying the write-time, status-query, and filing reference; the per-dispatch ensign core carries no dev-workflow-specific discipline. The deferred reference still resolves at its trigger (a write, a status query, a new-entity filing, a dispatch). The `«fn»` binding layer is tightened where the prose/binding boundary sits on the wrong side, net-neutral on size (legibility, not occupancy). Every reduction is proven by a **measured file delta vs `origin/main`**, never a prose-grep over the instruction files.

## Why

The contract is already lean — the 0203/0204 token roadmap is mostly shipped, and a 41-agent efficiency review found only ~925 raw tokens of prose trims left, no order-of-magnitude move. The remaining prize is **architecture, not prose-trimming**: ~32% of the ~7,170-token boot core is phase-specific reference (Status Viewer, FO Write Scope, ID Styles, Probe, Single-Entity, Mod Hook) needed at dispatch/status/write time, not at greet. Deferring it — the same name-a-pointer / defer-the-body pattern the Dispatch and Merge modules already use — is the real win. The ensign core compounds it at a different tier: it loads on *every* worker spawn, so dev-only discipline baked into the universal core recurs per dispatch and is also a correctness leak (the ensign runs non-dev workflows too). The `«fn»` refinements are quality polish deferred past 0.23.0 by the 0230 binding-effectiveness audit.

## Definition of Done

`v0.24.0` ships when, merged to `main`:

- **Boot-core occupancy drops, measured net-NEGATIVE vs `origin/main`.** The Status Viewer + Issue Filing reference (~675 tok) and the FO Write Scope + ID Styles reference (~1,027 tok) are deferred out of `first-officer-shared-core.md` into lazily-loaded references reached by a boot-core pointer at their triggers.
- **The greet/boot path is independent of the deferred references — proven as a real gate, not a prose claim.** A greet-and-stop boot completes without loading any deferred reference; a dispatching boot loads the write reference before its first frontmatter write or `spacedock new`; and the boot's own `«state.sweep-merged»` write path does **not** depend on the deferred FO Write Scope reference. Verified on fixtures.
- **No write-permission guard or reachability is lost.** Every FO Write Scope rule, ID-style minting path, Status Viewer invocation, and Issue Filing rule resolves at its trigger via the boot-core pointer.
- **The universal ensign core is stage-neutral and regression-guarded.** The dev stage-name parentheticals (`(implementation, validation)` / `(ideation, backlog)`) are removed from the Split-Root State Contract (−49 chars, net-negative), and a NEW `contractlint` structural-absence test (with a discriminator control) locks the core against re-introducing concrete stage names — #290's prose-lock having been retired by the instruction-read quarantine. *(The TDD / "CODE only" / code-only-deliverable markers were already re-homed by **#290** — PASSED/archived 2026-06-03 — prior art, not a 0240 deliverable; a dispatch-build test confirms that re-homing still delivers the discipline to a dev ensign.)*
- **The redundant `--read` adoption guidance is trimmed to its residue without regressing adoption.** Grep is named as the primary section-locator consistently across the FO + ensign cores (including the gate-verdict read, FO `:100` / Site D); the dispatch goldens stay byte-identical (a negative control — the dispatch-prompt hint was already removed in #392); and a journeymetrics before/after shows `--read`+scoped-`Read` adoption does not regress (a behavioral measurement, not a prose-grep).
- **The read-adoption metric measures the ensign, not just the FO.** journeymetrics folds the dispatched-ensign sub-agent transcript (`subagents/agent-*.jsonl`) into the per-journey `--read`+scoped-`Read` counts, so the adoption-not-regress measurement above is a real before/after instead of a vacuous `0==0`. Verified by a journeymetrics unit test with a perturbation control. *(This is `f5`, pulled in as `82k`'s predecessor at the 2026-06-30 ideation gate.)*
- **The gate's AC-scan sees value-annotated ACs.** `status --ac-scan` enumerates `**AC-N (annotation)**` headers (not just bare `**AC-N**`), so the deterministic `«gate.ac-cross-check»` extract no longer silently drops the value AC. Verified by a RED-then-GREEN fixture test + a no-over-match control. *(This is `48g`, pulled into 0240 2026-06-30 — surfaced live: the FO's `--ac-scan` rollups returned empty on every 0240 entity this session.)*
- **The `«fn»` binding layer is refined net-neutrally.** The z4 refinements trade prose for structure: every FO contract file z4 edits is net-neutral-or-NEGATIVE vs `origin/main` (z4's own-delta guardrail), and stays ≤ its v0.22.0 baseline for files z4 *materially restructures* — `claude-fo-dispatch.md` is **exempt** from the absolute ceiling (already +508 B over v0.22.0 from prior work; z4 only collapses one ref there, net ~0). `go build ./...` + `go test ./internal/contractlint/ ./internal/ensigncycle/` green; `docs/runtime-support.md` reflects the final binding shapes; and a **four-entry** `«fn»`-registry block indexes all four deferred modules (Dispatch, Merge, Status-Viewer, Write-Scope), with every one of the four load-points + greet-guards proven to survive the consolidation.
- **Debriefs route correctly for split-root workflows.** The debrief skill reads/writes/commits to `{state_checkout}/_debriefs/` for a `state:`-declared workflow (single-root unchanged), with continuous numbering and no cross-home collisions. Verified on split-root and single-root fixtures.
- **The cut is clean and contract-consistent.** `go test ./...` green from the root; the cut is consistent with the contract-2 binary (`origin/main` carries #443, the contract 1→2 bump); the **pre-cut antipattern audit** is clean; then the release ritual per [`docs/releasing.md`](../../releasing.md).

## Sequencing (the load-bearing operational constraint)

The members edit the **same 2–3 contract files**, so implementation dispatches run in **strict sequence, never in parallel** (parallel = merge conflicts in the same paragraphs). Ideation may fan out in parallel — it writes only to each member's own body. The implementation order:

- **Predecessor (different surface, parallel-safe):** `journeymetrics-ensign-read-adoption` (f5) — pulled into 0240 at `82k`'s ideation gate (2026-06-30, captain) as `82k`'s measurement predecessor: it folds the dispatched-ensign sub-agent transcript into the journeymetrics `--read` adoption metric that `82k`'s AC-1 measures (without it the before/after is a vacuous `0==0`). f5 edits `internal/journeymetrics` (Go) — no collision with the contract-file waves, implementable in parallel — but it MUST land (implementation + validation) before `82k`'s AC-1 can be measured.
- **Wave 1 — trim the cores:** `read-guidance-redundant-with-grep` (82k) and `ensign-contract-dev-leakage` (scr) — both edit the ensign core; coordinate their edits. (`82k` ideation gate, 2026-06-30: FO `:100` / Site D confirmed in scope; `f5` above is its measurement predecessor.)
- **Wave 2 — defer reference out of the boot core:** `entity-status` (84) then `defer-write-scope-id-styles` (k4) — both edit `first-officer-shared-core.md` (different sections).
- **Wave 3 — last:** `fn-binding-refinements` (z4) — its #6 registry block *indexes* the now-deferred structure, so it must follow Wave 2.
- **Independent (any time, different file):** `debrief-split-root-state-home` (7d4).

## Gate decisions (ideation)

Captain decisions at each member's ideation gate (the *design* gate; gates lock after the sprint-wide preflight staff review):

- **`82k`** — APPROVED (2026-06-30). `f5` pulled into 0240 as the measurement predecessor (1a); FO `:100` / Site D confirmed in scope (2).
- **`f5`** — APPROVED (2026-06-30). AC-2 (widen `commandInvokesStatusRead` launcher recognition) confirmed IN scope, so the `status_read` before/after is non-vacuous on real data (the fold alone left it ~0 — only ~8% of real `status --read` launcher forms were recognized).
- **`scr`** — APPROVED (2026-06-30), scope **(a)**: ship the stage-neutrality fix (−49 chars) + the quarantine-legal regression guard. The original seed (`ep0ra3z`) was delivered by **#290** (PASSED/archived); `scr` re-scoped to the residue. Packaging follow-ups: re-title (current title is a misnomer) + rename the slug to clear the `ensign-contract-dev-leakage` collision with archived `ep0ra3z`.
- **`7d4`** — APPROVED (design, 2026-06-30). Complexity-review (split-root + skill integration) routed into the consolidated independent staff review.
- **`k4`** — APPROVED (design, 2026-06-30). Sweep-independence accepted: the boot sweep is binary-self-authorizing via `state sweep` (read-only detect) + the pr-merge mod's `status --set` / `status --archive` finalize — **NOT `merge guard`** (that's the deferred terminal ceremony, never a boot path; the spike + AC-2 test-pin are re-pointed to the `runSet` / `runArchive` guards in the fold). Complexity-review folded into the staff review.
- **`z4`** — APPROVED (design, 2026-06-30), two ratified AC corrections: **AC-6** needs ONE live re-run (the shipping `TestLiveReanchorGateRejectsMeansOnlyRegressed` scenario — z4 is not purely offline); **AC-4** re-scoped to two parts — (a) z4's own per-file delta net-neutral-or-negative [z4's gate]; (b) absolute ≤ v0.22.0 only for files z4 materially restructures (`claude-fo-dispatch.md`'s pre-existing +508 B overage is not z4's charge).
- **`84`** — APPROVED (design, 2026-06-30). Arrow-taxonomy for the adapter-less deferred refs (`→ reference` vs `→ runtime-binding`, shared with `k4`) deferred to `z4` (owns `runtime-support.md` + the registry); staff review to confirm consistency.
- **Wave 1 kept** (decision 5): the marginal Wave-1 trims stay in 0240; the sprint is not re-weighted toward the deferrals.
- **Gates LOCK after the preflight staff review** (multi-lens independent Workflow, launched 2026-06-30): all 7 design-approvals are provisional until the review's Material findings are folded.
- **Late additions (captain, 2026-06-30):** `48g` (ac-scan-value-annotation-skip) pulled in as the **8th member** (gate-tooling; parallel-safe like `f5`) — ideated after the cohort review launched, so it gets a **focused follow-on review** rather than the cohort pass. `0qt` (dispatch-build-team-name-advisory) is a standalone usability **mustfix driven SEPARATELY** from the cohort (its own review track) — **not** a 0240 member.

## Out of scope

- The smaller third deferral sibling (Probe ~347, Single-Entity ~107, Mod Hook ~147 tok) — a follow-on after these land, not this sprint.
- 0.25.0 candidates: `fo-opus-behavioral-robustness`, `minor-version-compat-coupling`.
- Adjacent contract-quality items not pulled in unless the captain adds them: `proof-policy-shipped-scaffolding`, `codex-foreground-wait-phrase-check-retire`, `ac2-reanchor-scenario-falsifiable`.
- The two unfiled tooling gaps (`status --validate` near-dup ids; `status --resolve` archived-scope).
- No new binary behavior beyond the contractlint guard-loosening z4 AC-3 requires.
