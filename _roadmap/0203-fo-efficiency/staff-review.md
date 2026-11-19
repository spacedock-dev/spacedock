# 0203 FO-efficiency — preflight staff review

**Verdict: READY** (after one j9 rework cycle). Sprint-wide preflight over the pooled sprint — per-item + cross-cutting. Shaping-FO session, 2026-06-13.

## Per-item readiness

- **j9** (lazy-teamcreate-shallow-boot) — **READY.** The contract restructure (boot-resident/deferred split → lazy-TeamCreate → shallow-boot-then-greet). Reviewed in depth: a 4-lens design panel + adversarial re-verify (6 + 3 findings closed), then this sprint preflight surfaced 3 batch-blockers — all fixed in cycle 4 and re-verified. AC set AC-1..AC-6, all external-proven, no prose-grep.
- **#344** (context-budget-spurious-warnings) — **READY** (validated, held pre-merge). Narrow Go fix; 5 ACs green including golden parity zero-churn; a detached adversarial audit confirmed the over-suppression guards load-bearing. Zero file overlap / behavioral interaction with j9/T3. Merges with the batch.
- **T3** (fo-contract-prose-audit) — **READY.** A 4-step audit method against j9's *planned* modules; ACs external (live-scenario preservation + `wc` size-floor + collapse-decision record); correctly sequenced behind j9 Phase-1; collapses to a roadmap decision if the split leaves nothing to cut.

## Cross-cutting coherence

- **Zero file overlap.** #344 touches `internal/claudeteam` Go; j9/T3 touch the FO markdown refs + `contractlint`/`ensigncycle` tests. No merge conflict is possible.
- **#344 ↔ j9** meet only at the stable `reuse_ok` surface — the contract references the budget probe as a black box, and neither warning key appears anywhere in the refs. Clean.
- **T3 → j9 Phase-1** dependency is sound + non-circular: T3 designs its method (not a cut-list) against j9's planned modules, with an explicit collapse fork.

## Material findings — all CLOSED in j9 cycle 4

- **M1** — the Phase-1 split breaks two hard-coded offline-gate tests (`TestNoUnexpectedModHookOrPRMergeIntroduced` via the `## Hook:` allowlist; `TestGradeMarkerMatchesContract` via the `TERMINAL_TEARDOWN_BOUNDED` marker). → j9 owns an explicit retarget subtask + **AC-5** (offline gate exits 0 post-split).
- **M2** — AC-4's structural proof was unbuildable (cited a test that reads only `SKILL.md`). → rewritten as a real new `contractlint` `os.Stat`-oracle test.
- **M3** — the sprint's headline goal (<~60k / 89k saving) had no AC. → **AC-6** measured-saving live drive (greet-turn context < ceiling + no pre-greet ~89k spike, off `claude-stream.jsonl`).

## Residuals — Polish, fold into implementation (non-blocking)

- j9: three stale historical lines say "AC-1..AC-4" (live set is AC-1..AC-6) — mark superseded.
- j9: AC-1(b) credits the team `config.json` path to a "comm-officer hook" that doesn't ship — loose attribution; the path itself is real.
- j9 AC-6 soft spot: the ~89k is asserted (no team was created in the forensics run); the negative control rides an eager-team fixture — the present/absent cache-creation-spike signal stays falsifiable.
- T3: AC-4's collapse-decision-record cites `README.md`, now the conventional `index.md` — repoint the path at T3 implementation.

## Provenance

Boot forensics: `boot-analysis.md`. Per-item ideation: the entity bodies under `docs/dev/.spacedock-state/`. Preflight (4-lens) + j9 re-verify: shaping-FO session, 2026-06-13.
