---
id: 6b8k79kbmzn8n48g2amf6q4m
title: Terminal-guard verdict=rejected consistency — align contract prose + the --set/--archive asymmetry
status: backlog
source: sprint-end antipattern reviews (2026-06-01) — 0.19.3 minor-findings bucket
started:
completed:
verdict:
score: "0.26"
worktree:
issue:
---

Two related inconsistencies in how the terminal-transition guard treats `verdict=rejected`:

1. **Contract prose vs code drift.** The merge-hook terminal guard already EXEMPTS `verdict=rejected` in code (`internal/status/handlers.go:163`; Python oracle `:2648`), so a rejected entity can terminalize without a PR/mod-block. But the FO operating-contract enumerations of the guarded terminal fields (`first-officer-shared-core.md`, Mod-Block Enforcement) never state the exemption — a reader who trusts the prose over-blocks rejected entities.
2. **`--set` vs `--archive` asymmetry.** `--set` exempts `verdict=rejected` but `--archive` (`internal/status/mutate.go:228`) does NOT, so reject-then-archive half-passes (the `--set` goes through) / half-blocks (the archive refuses) — an entity can stall mid-terminalization.

## Acceptance criteria (provisional — harden at ideation)

**AC-1 — `--set` and `--archive` agree on verdict=rejected.** A rejected entity in a merge-hook workflow with empty `pr`/`mod-block` either terminalizes-and-archives cleanly OR is blocked consistently by both surfaces — never half-and-half. Ideation decides the direction (exempt both, or guard both).
Verified by: a behavioral test driving reject → terminal `--set` → `--archive` on a merge-hook fixture, asserting both surfaces behave identically.

**AC-2 — Contract prose states the verdict=rejected exemption.** The FO terminal-guard enumerations name the exemption so prose matches code.
Verified by: the contract reference text states it (prose-presence is legitimate here — the contract text IS the deliverable, per the workflow's legitimacy boundary).

## Notes
- `internal/status` serialized lane (handlers.go/mutate.go) — coordinate with #251 and any other status-lane implementation (one impl worktree at a time). The contract-prose edit is protected scaffolding → dispatched worker only, authored in the vendored copy.
