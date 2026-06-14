---
id: 5ew2jxagk11mr0fzd0rtpdp0
title: In-module restructure of the FO contract refs (consolidate duplicated obligations + collapse redundant sections)
status: backlog
source: "T3 (fo-contract-prose-audit) deferred this (2026-06-14): T3 shipped the mechanical-safe subset (4 dead-ref cuts + comm-officer concision); the substantive restructure was scoped out (duplicated obligations marked KEEP; the merge mod-block section-collapse FLAGGED out-of-scope). Captain: the audit \"would imagine a bigger cleanup and in-module restructure.\""
started:
completed:
verdict:
score:
worktree:
issue:
---

The substantive in-module cleanup of the post-j9 FO contract refs that T3 deferred. T3 did the safe mechanical subset (dead-ref repairs) and a low-value within-line concision polish (~75 of 79 changed lines, which over-reached on meaning 3× and cost two amend cycles + a rejection). The real value — consolidating genuinely-duplicated obligations and collapsing redundant sections within the modules — was explicitly punted as "judgment-call restructure, not a behavior-preserving mechanical cut."

## Scope (the deferred restructure)

1. **Collapse the merge ref's two mod-block sections.** `skills/first-officer/references/claude-fo-merge.md` carries `## Mod-Block Enforcement` AND `## Mod-Block Enforcement at Terminal Transitions` — they restate the same mechanism-level invariant + the session-resume rule (j9 added the first adjacent to the pre-existing second). Collapse to one canonical section + the guard, behavior-preserving. (This was T3's explicit FLAGGED-out-of-scope item.)
2. **Consolidate genuinely-duplicated obligations** the boot-resident↔deferred split created. T3's survey found duplications (C1 MODS-REPORT vs RUN-STARTUP-HOOKS; C5 gate-spine boot-resident vs reuse-conditions deferred; the concurrency-safe-commit rule; worktree-ownership) and marked them KEEP ("already canonical+pointer"). Re-examine each: where the two phrasings can collapse to ONE canonical statement + a pointer (vs a genuine kept duplication at two lifecycle moments), do so.
3. **Module-level coherence pass** (optional): re-organize within a module where the cross-cycle prose accretion left it incoherent — NOT a comm-officer concision pass (T3 already did that; it's low-value + risky on the contract).

## Why this is risky (and how to prove it)

A section-collapse or obligation-consolidation can DROP or INVERT an obligation that the live scenarios don't exercise — exactly the class the detached adversarial audit exists for (it caught T3's dropped NEVER-qualifier). So:
- **AC (behavioral, live):** the existing live shared scenarios (`gate-guardrail` / `rejection-flow` / `feedback-3-cycle-escalation` / `merge-hook-guardrail`, Claude + Codex) stay green after the restructure. The mod-block collapse specifically must keep `merge-hook-guardrail` green.
- **AC (high-stakes detached audit):** a word-level diff of every collapsed/consolidated obligation against the pre-restructure baseline, confirming no MUST/MUST-NOT/qualifier dropped or inverted, and the `TERMINAL_TEARDOWN_BOUNDED` marker byte-intact.
- **AC (structural):** `internal/contractlint` reference-closure + the offline gate (`go test ./...`) stay green.

## Out of scope

- The team-vs-bare dispatch-mode determination (separate task `7e` / `headless-dispatch-mode-intent`).
- A comm-officer concision polish (T3 did it; it is low-value + meaning-change-risky on the contract — do NOT repeat it here; if comm-officer is used at all, harden its guard per the xf note first).

## Notes

Fast-follow, not a v0.20.3 blocker (T3's behavior-preservation shipped). A `comm-officer` polish-over-reach guard (it changed contract meaning 3× under "light-touch") should be folded into `xf` (which moves comm-officer usage prose into its mod) — a hard "never touch MUST/MUST-NOT/qualifiers in contract prose" rule.
