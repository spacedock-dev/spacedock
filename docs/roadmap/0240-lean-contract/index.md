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
- **The per-dispatch ensign core is host-/stage-neutral.** No dev-workflow-specific discipline (TDD framing, "CODE only" worktree, code-only deliverable) remains in the universal ensign core; it is re-homed into the dev-shape scaffolding, measured net-NEGATIVE in the universal core, and a dispatch-build test on a dev fixture confirms a dev-workflow ensign still receives it.
- **The redundant `--read` adoption guidance is trimmed to its residue without regressing adoption.** Grep is named as the primary section-locator consistently across the FO + ensign cores (including the gate-verdict read, FO `:100` / Site D); the dispatch goldens stay byte-identical (a negative control — the dispatch-prompt hint was already removed in #392); and a journeymetrics before/after shows `--read`+scoped-`Read` adoption does not regress (a behavioral measurement, not a prose-grep).
- **The read-adoption metric measures the ensign, not just the FO.** journeymetrics folds the dispatched-ensign sub-agent transcript (`subagents/agent-*.jsonl`) into the per-journey `--read`+scoped-`Read` counts, so the adoption-not-regress measurement above is a real before/after instead of a vacuous `0==0`. Verified by a journeymetrics unit test with a perturbation control. *(This is `f5`, pulled in as `82k`'s predecessor at the 2026-06-30 ideation gate.)*
- **The `«fn»` binding layer is refined net-neutrally.** The z4 refinements land with every FO contract file **≤ its v0.22.0 baseline** (the promotions trade prose for structure); `go build ./...` + `go test ./internal/contractlint/ ./internal/ensigncycle/` green; `docs/runtime-support.md` reflects the final binding shapes; and the new `«fn»`-registry block indexes the deferred modules (Dispatch, Merge, and the new Status-Viewer / Write-Scope references).
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
- **Wave 1 kept** (decision 5): the marginal Wave-1 trims stay in 0240; the sprint is not re-weighted toward the deferrals.

## Out of scope

- The smaller third deferral sibling (Probe ~347, Single-Entity ~107, Mod Hook ~147 tok) — a follow-on after these land, not this sprint.
- 0.25.0 candidates: `fo-opus-behavioral-robustness`, `minor-version-compat-coupling`.
- Adjacent contract-quality items not pulled in unless the captain adds them: `proof-policy-shipped-scaffolding`, `codex-foreground-wait-phrase-check-retire`, `ac2-reanchor-scenario-falsifiable`.
- The two unfiled tooling gaps (`status --validate` near-dup ids; `status --resolve` archived-scope).
- No new binary behavior beyond the contractlint guard-loosening z4 AC-3 requires.
