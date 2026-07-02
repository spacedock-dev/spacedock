---
title: "State/merge error-path remediation moves from resident contract prose into the binary's own failure output"
status: ideation
group: tooling
source: "fable-token-trim-scout analysis 2026-07-02: the five state «fn» bodies (first-officer-shared-core.md:106-144, ~1,015 tok boot-resident) and much of fo-merge-core.md (~1,824 tok at terminal) re-describe shipped commands plus failure remedies the binary could emit at fire time — e.g. state commit exit 3 already means 'same-entity rebase conflict aborted' and its stderr can carry the halt/surface/never-force instructions; merge guard phases can name the FO's next step in their output. Est. ~600-700 off boot, ~700-900 off terminal. Fresh angle: instruction-space -> tool-output placement, not prose relocation."
id: s058nyrecqwtegn36rx6yew1
started: 2026-07-02T03:02:55Z
---

## Problem
Contract prose pays resident tokens to restate what shipped commands do and how to recover when they fail, even though that guidance is only needed at the moment the command actually fails or signals — when the binary is already talking to the FO through stdout/stderr.

## Desired direction (for ideation to refine)
`state ready/sweep/commit` failure modes and `merge guard` phase signals emit their own next-action guidance (e.g. exit-3 stderr carries the halt-and-surface instructions; `armed`/`blocked`/`finalized` lines name the FO's next step). The contract keeps guard one-liners and pre-invocation prohibitions (never-force stays resident AND stays a binary refusal — belt and suspenders is the current state; ideation decides what prose is safe to drop given the binary refusal already exists). Shared-core «fn» bodies and fo-merge-core shrink accordingly.

## Rough acceptance sketch (ideation tightens into measured ACs + a test plan)
- Binary failure/phase output carries the remediation text, pinned by tests on exit codes + stderr/stdout content.
- Boot-resident shared-core and terminal fo-merge-core shrink by a measured token delta; the dropped prose is provably covered by binary output (no guidance lost).
- Pre-invocation prohibitions remain resident; contractlint stays green (touches skills/**, so claude-live gates the merge).
