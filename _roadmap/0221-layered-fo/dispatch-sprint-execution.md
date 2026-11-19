# 0221 (0.22.1) — layered FO — Commander dispatch (cold-boot)

## Boot

Sprint = the entities matching `sprint: 0221-layered-fo` (query, not a list). Drivable set:
`spacedock status --workflow-dir docs/dev --where sprint=0221-layered-fo --where 'sprint-readiness != defer'`
→ **6re** (gate-extract-verbs), **rgq** (state-verbs), **mz** (merge-finalize), **72** (fo-tier-delegation), **czw** (prose-function-restructure), **3e** (wrong-root-guard-compound-boot). Deferred (NOT in the drive): **kt** (haiku-drive-validation — the gate proof, carve LAST) and **y2** (fo-contract-token-cut).

Boot the first officer (`spacedock claude`), run `status --boot`, and read each member body for its gate-approved design + ACs. Readiness: `staff-review.md` (2026-06-19 sprint-wide, verdict GAPS-TO-CLOSE — the four holes are resolved below). Goal/DoD: `index.md`. Foundation: this sprint rides the **0.22.0 merged floor** (m4 + CI-unpin shipped; `TeamCreate` gone on `.178+`, named-background-`Agent` back-channel is the norm).

## Deliverable & DoD

**0.22.1** = the Haiku-operable layered FO: the verb core (`rgq`/`mz`) + gate-extract (`6re`) + tier-delegation (`72`) + the prose-function restructure (`czw`), all on top of the 0.22.0 floor. Done when merged to `main` — see `index.md` Definition of Done. **DoD amended (captain 2026-06-19):** the `next-action` *binary* descoped to `0222-dispatch-driver`; this sprint ships its Haiku-operable `«dispatch.next-action»` *prose-function* form in `czw` (w4 found `«next»` HOLDs). The sprint's terminal proof — a live Haiku FO driving one entity dispatch→gate→merge→terminal with the verdict routed to a standing L3 — is **`kt`, deferred until the drivable set ships** and re-grounded to the merged floor (below).

## Drive order — ⚠️ coordination

The three contract-editors touch the same boot-resident regions (`## Startup`, `## Completion and Gates`, `## FO Write Scope`). **RATIFIED landing order (captain 2026-06-19): `6re` → `72` → restructure (`czw`).** The roster shows all three `ready` with no ordering field — the Commander enforces the order from here and `index.md`. Landing `czw` early forces a rewrite collision and risks silently dropping `72`'s fail-safe property.

1. **`6re` first** — built and rebased onto current `main` (branch `spacedock-ensign/gate-extract-verbs`, 4 ahead / 0 behind, merge-ready pending its validation gate). Drive it through validation → merge; it lands the `## Startup` step-4 rewrite the restructure builds on.
2. **`rgq` + `mz` (verb core)** — gate-approved at ideation, 2y-unblocked; build in parallel. They ship the `state ready|sweep|commit` and `merge guard` verbs whose bodies `czw` flips guillemet→backtick. `rgq` rebase-HALT and `mz` mod-block atomicity are enforced by the *verb*, not FO discipline — verify that at validation.
3. **`72`** — gate-approved; lands the `## Startup` step-1.5 (`«fo.tier»` arming) + the gated-stage verdict-route bullet. Its prior two material gaps are **folded** (fail-safe default; the `SPACEDOCK_FO_MODEL` launcher surface — verified genuinely-new against `frontdoor.go`).
4. **`czw` (restructure) LAST** — re-expresses the merged `6re`+`72` edits into prose-function form, folding `«fo.tier»()` in as a sixth declared prose-function. **RATIFIED: `«gate.assemble-verdict»` stays prose** (verdict is judgment; body names `6re` extract modes + the L3 route). Its non-grep proof is the `contractlint` backtick↔cobra-command-tree binding.
5. **`3e`** — independent test-harness parser fix; no ordering constraint (zero contract overlap). Land anytime; clears the v0.22.0 e2e-gate waiver so the next live-e2e is clean-green.

## Per-member build notes

### 6re — gate-extract-verbs · new `spacedock gate` CLI (read-only over files) · built, merge-ready
The three extract modes feed L3 structured input and never adjudicate. Prior material gaps (AC-1 stage-scoped selection; Verb-3 retarget to `### Feedback Cycles`; AC-4 flags) — confirm closed at the validation gate before merge.

### rgq — state-verbs · ⚠️ HIGH-STAKES (status mutation/guard path)
`state ready|sweep|commit <slug>`. The rebase-conflict HALT must be enforced by the commit verb refusing to return (exit non-zero, state left rebase-aborted), not by prose. Detached adversarial audit at validation.

### mz — merge-finalize · ⚠️ HIGH-STAKES (merge/terminal path)
`merge guard <slug>` — atomic mod-block set→invoke→clear→terminalize. The status tool's refusal of terminal-with-mod-block-set is the backstop; confirm it still holds against the 2y-extracted merge core. Detached audit at validation.

### 72 — fo-tier-delegation · ⚠️ HIGH-STAKES (shipped contract + launcher front-door)
Fail-safe tier default (unset/unresolvable → level-2-only, route gate verdicts); the `SPACEDOCK_FO_MODEL` launcher export + `--env-pass` plumbing. Prove the var reaches a launched FO. Detached audit at validation.

### czw — prose-function-restructure · ⚠️ HIGH-STAKES (shipped contract restructure) · lands LAST
Re-express the cores as `«fn»` invocations; the `contractlint` notation↔command-tree check (with a RED control) is the non-grep proof. Behavior-preserving — the existing closure/ceremony suites stay green. Detached audit at validation.

### 3e — wrong-root-guard-compound-boot · test-harness only · routine
Trailing-separator trim in `bootPathArgs`. Zero product/contract impact (doc-diff none) — routine, no detached audit.

## Detached adversarial audit (before merge)

High-stakes surfaces per the proof policy: **rgq** (status mutation/guard), **mz** (merge/terminal), **72** (shipped contract + front-door launcher), **czw** (shipped contract). Run a read-only detached audit on a throwaway checkout of each merge result before merge — refute that the deliverable's own tests would catch a claim-breaking edit. **6re** is a read-only CLI over files (low blast radius — routine). **3e** is test-harness-only — routine.

## Pre-cut antipattern audit (⚠️ before the v0.22.1 tag)

All members merged, tag not yet fired → an INDEPENDENT staff-eng reviewer over the assembled sprint. Confirm the boot-region landing order did not drop `72`'s fail-safe property in `czw`'s re-expression, and that the `contractlint` notation check is genuinely non-grep (binds against the compiled command tree, not the prose). Ship-blockers fixed pre-cut; non-blockers seed the next sprint.

## Cut

Fire `v0.22.1` once the drivable set is merged and the pre-cut audit is clean. Follow `docs/releasing.md`.

## Deferred — carve before kt drives (gate-proof preconditions)

`kt` is the sprint gate proof, sequenced LAST and currently `sprint-readiness: defer`. Before carving it out of defer, close two items the 2026-06-19 staff review flagged:

1. **Re-ground `kt` to the merged floor.** Its body still assumes the retired 2.1.177 pin + interactive-tmux-only team tools; 0.22.0 unpinned CI and `m4` proved headless `-p` is team-enabled flag-free on the merged floor. Rewrite the launch-shape premise before carving.
2. **The `reconcile.go` headless-meta roster-discovery fix — DECISION DEFERRED to kt-carve (captain 2026-06-19).** On a headless `-p` merged host, `reconcile` degrades to git-only (the loop step the gate proof exercises). At carve, either file the fix (cleanup group, like `3e`) or record an explicit "kt drives without it" decision.

## Out of scope (deferred)

- **`0222-dispatch-driver`** — the `next-action` binary (the fully-qualified `{action, team_action, …}` return schema) and `vp` (dispatch-build-request-file). 0.22.1 ships only the prose-function form.
- **`kt` + the `reconcile.go` decision** — deferred to kt-carve (above).
- **`y2` fo-contract-token-cut** — deferred; its premise is superseded — re-ground to the revised `docs/dev/_proposals/fo-contract-token-cleanup.md` (~420 default-path tokens, RT-4/RT-2 retired, team cuts legacy-only) before it drives.
- **Codex / Pi Haiku operability** — the live drive validates the Claude substrate only; other hosts are a follow-up.
