---
title: "Collapse the FO Startup recipe to <=4 prose steps behind one shipped boot call"
status: backlog
source: "Captain, 0250 Commander session 2026-07-07, post k7 fo-boot-engage-split merge: 'i am allergic to the > 4 steps recipes. can we now further clean it up?' Startup in first-officer-shared-core.md is an 8-step prose recipe (version gate, project root, discovery, taxonomy read, state.boot, state.ensure-ready, state.sweep-merged, greet/headless). Steps 2-7 are deterministic orchestration the binary can own — the workflow's own prefer-code-gate-over-prose principle applied to its boot."
started:
completed:
verdict:
score: 0.5
worktree:
issue:
id: 1y4ynffdxcgxn5eqcgw1mps3
---

The FO Startup procedure spends six of its eight prose steps orchestrating deterministic reads the binary already implements piecemeal (--version, --discover, --read README, --boot, state ready, state sweep). Direction: a single `spacedock boot` verb that runs the whole pre-greet sequence internally and emits one JSON boot record (taxonomy + boot sections + ready/sweep outcomes), preserving EVERY per-class abort behavior (binary absent/wrong version, zero/multi discovery, split-root halt classes) with remediation on stderr — the collapse must not blur failure UX. Startup prose shrinks to: resolve launcher + boot call; greet/drive from the record. Acceptance sketch: value — the resident Startup section's numbered steps drop 8 → <=4 AND the shared-core resident byte delta is NEGATIVE, measured against the pre-change file; the per-class aborts are each observed live (scenario per class, not prose-grep); mechanism — the boot verb ships with the recipe rewrite pointing at it. High-stakes surfaces (shipped contract + status/state paths): detached adversarial audit + full lane treatment. Sequencing: touches first-officer-shared-core.md Startup section — respects the 0250 Wave-1 strict-serial chain; earliest safe slot is after vcm's merge.
