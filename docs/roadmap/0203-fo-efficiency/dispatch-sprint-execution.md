# 0203 (0.20.3) — FO efficiency — Commander dispatch (cold-boot)

## Boot

Sprint = the entities matching `sprint: 0203-fo-efficiency` (query, not a list): **j9** (lazy-teamcreate-shallow-boot), **#344** (context-budget-spurious-warnings), **T3** (fo-contract-prose-audit). Boot the first officer (`spacedock claude`), `status --boot`, and read each entity body for its gate-approved design + ACs. Readiness: `staff-review.md` (verdict READY). Goal/DoD: `index.md`. Evidence: `boot-analysis.md`.

## Deliverable & DoD

**0.20.3** = the FO-efficiency restructure + the context-budget probe fix. Done when merged to `main` — see `index.md` Definition of Done. Headline: a live drive **measures** boot reaching < ~60k with the ~89k team-mode re-cache absent before the greet.

## Drive order — ⚠️ coordination

1. **j9 first** (the backbone; T3 keys off it). Per the operating principle — *begin with the end; do the hardest first, de-risk when it's cheap* — within j9 land the cheap, biggest lever first: **lazy-TeamCreate + shallow-greet + the AC-6 measured-saving drive**, proving the <60k/89k saving *before* the full contract split. Then **Phase-1** (the split + the offline-gate-assertion retarget). If the cheap measure already clears <60k, the split is the contract-cleanup / residual-~16k play, not load-bearing for the saving — decide then whether it stays in 0.20.3.
2. **T3 after j9 Phase-1 lands** — the slimmed refs must exist before the audit. Step-0 survey decides the collapse fork (cut vs recorded decision).
3. **#344** — already validated (`46224f5f`); merge with the batch. No ordering constraint (zero overlap).

## Per-member build notes

### j9 — lazy-teamcreate-shallow-boot · shipped-scaffolding surface · ⚠️ HIGH-STAKES
The 3-phase restructure (full spec in the entity body, AC-1..AC-6). It rewrites the very contract the FO + ensigns run under — test the new contract in isolation (the live `shallow-boot` scenario + the `contractlint` closure test) before merge. **Retarget the two offline-gate assertions** (`TestNoUnexpectedModHookOrPRMergeIntroduced` allowlist; `TestGradeMarkerMatchesContract` source) to the post-split layout as an explicit subtask, and keep `go test ./...` green (AC-5). Clean the 3 Polish residuals from `staff-review.md` (stale AC-count lines; AC-1(b) attribution; AC-6 89k-soft-spot note).

### #344 — context-budget-spurious-warnings · dispatch-path · DONE (held pre-merge)
Implemented + validated on `spacedock-ensign/context-budget-spurious-warnings` @ `46224f5f` (the `<synthetic>` census skip + the `[1m]`-suffix window promotion). 5 ACs green; golden parity zero-churn; detached audit confirmed the over-suppression guards load-bearing. Just merge with the batch.

### T3 — fo-contract-prose-audit · contract-cleanup · BLOCKED on j9 Phase-1
The 4-step audit method (survey → mechanical cut → comm-officer polish) against the slimmed refs. Step-0 survey is the collapse fork: non-empty inventory → code change; empty/trivial → a recorded roadmap decision (AC-4). Steer the survey to KEEP the budget-probe reuse-condition-0 prose (a deliberate cross-host abstraction split, not collapsible duplication).

## Detached adversarial audit (before merge)

High-stakes surface: **j9** (shipped contract/scaffolding). Run a read-only detached audit on a throwaway checkout of the merge result before merging j9 — refute that the live scenarios + `contractlint` guards would catch a broken edit. **#344** already had its detached audit. **T3** is behavior-preserving (live scenarios) — routine.

## Pre-cut antipattern audit (⚠️ before the v0.20.3 tag)

All merged, tag not yet fired → an INDEPENDENT staff-eng reviewer over the assembled sprint. **Critically: confirm AC-6 actually measured the <60k/89k saving in the live run** — the sprint's whole point. Verify main-PR CI gating. Ship-blockers fixed pre-cut; non-blockers seed the next sprint.

## Cut

Fire `v0.20.3` once the three are merged and the pre-cut audit is clean.

## Out of scope (deferred)

p2/vc (0.20.4 binary-simplification line); xp (cross-session FO↔Commander comms — the coordination gap this sprint hit live); ey (proof-policy port to shipped scaffolding).
