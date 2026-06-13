---
id: 7d72hqzg9tpnr6tdjfa34nxt
title: dispatch context-budget — suppress spurious model warnings on healthy team members
status: backlog
source: github#344 (captain intake 2026-06-13)
started:
completed:
verdict:
score:
worktree:
issue: "#344"
---

`spacedock dispatch context-budget` emits two warnings that read as faults but are environmental noise on a healthy reused team member, eroding trust in the reuse-condition-0 budget signal. Intook to the 0.20.3 (0203-fo-efficiency) sprint as FO dispatch-path quality.

## Problem

(github#344) On a healthy reused team member, `spacedock dispatch context-budget --name {member}` emits:
- `config_drift_warning` — team config records the captain session model string (e.g. `claude-fable-5[1m]`) while the runtime jsonl reports the canonical id (`claude-fable-5`); these never match, so it fires every probe.
- `mixed_models_warning` — jsonl `<synthetic>` (harness-injected) entries mix with the real model → "multiple models seen — using smallest context window"; with `<synthetic>` unrecognized the fallback window is opaque.

Expected: `[1m]`-suffixed config values normalize before comparison; `<synthetic>` entries excluded from the model census. Version: binary 0.20.0 (contract 1), plugin 0.19.9.

## Notes

Ideation fills approach + acceptance criteria — each proven by a Go test over the probe's actual output (warnings present/absent), never a prose check.
