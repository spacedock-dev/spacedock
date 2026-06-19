---
id: 5xst4dae1m2acx76czpe857v
title: FO behavioral-robustness — live opus-4-8 FO intermittently deviates from contract guidance
status: backlog
source: '0221 cut (2026-06-19): the live claude-live OPUS lane flaked on essentially every PR with DIFFERENT non-deterministic FO-behavior deviations, while SONNET stayed compliant. Diagnosed diff-unrelated across 3e #398 (sonnet broad-search) and rgq #399 (two reruns: no spacedock-new, then merged-team-mode exit-before-terminalize) and accepted non-blocking under the ci-env conn to land 0.22.1. This is the recorded root-cause follow-up.'
started:
completed:
verdict:
score: 0.6
worktree:
issue:
group: cleanup
---

The live opus-4-8 first officer intermittently fails to follow contract guidance precisely — surfaced by the live-e2e lanes doing their job — making the opus claude-live lane unreliable on every PR. This is the strong-model echo of the weak-model concern the layered-FO sprint targets: even opus deviates, so the deferred tier/delegation mechanism (72) and report-and-stop hardening matter.

## Problem

Across the 0.22.1 cut the opus claude-live lane failed on three DISTINCT non-deterministic FO behaviors, none related to the PR under test, while sonnet/codex/pi passed:

1. **No `spacedock new`** — the FO hand-filed a seed instead of the atomic-create path (`claude_live_runner_test.go:122`).
2. **Broad-search after zero-discover** — ran `find <root> -type f` instead of report-and-stop (`zero_discover_live_test.go`, the `tq0` sibling).
3. **Exit before terminalize** — the FO subprocess exited (code 0) before driving the entity to terminal in merged-team-mode (`merged_team_mode_live_test.go:176`).

The lane won't converge by re-running (a different deviation each run), so it blocks every PR's merge.

## Proposed direction (ideation to refine)

Harden the FO contract's behavioral robustness so a strong-model FO reliably complies — and/or right-size any over-strict live assertions where a capable FO legitimately varies. Distinguish, per failing scenario, "contract guidance too weak" vs "test too strict." Relates to `tq0` (broad-search) and to the deferred `72` (tier/delegation forces verdict/judgment compliance structurally).

## Out of scope

- The deferred `72` tier mechanism itself (separate rebuild on the cleaned foundation).
