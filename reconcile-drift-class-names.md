---
title: Reconcile drift output carries descriptive class names instead of A-E letters
status: backlog
source: captain (2026-06-14, this session) — `spacedock dispatch reconcile` emits opaque single-letter drift classes (A-E) when the descriptive names already exist in the FO dispatch contract event-loop step-0 mapping (A=lingering, B=superseded, C=un-advanced-pr, D=stale-branch, E=local-main-drift). The letter is indirection over a name the `reason` field already states in English. Decided this session — emit the descriptive names and drop the letters; a string enum is as stable to branch on, with no machine-stability cost.
started:
completed:
verdict:
score: 0.3
worktree:
issue:
id: pd7fqh4f8yzf9dacbbbamfg7
---

The `spacedock dispatch reconcile` output tags each drift entry with a bare letter (`"class": "A"` .. `"E"`). The descriptive names already exist — the FO dispatch contract event-loop step-0 spells every class out parenthetically — but the helper discards them and emits the letter, so a reader (FO or human) must carry the A-E mapping to read the output. Decided this session: emit the descriptive name and drop the letter.

## Problem

`spacedock dispatch reconcile` emits `drift[].class` as a single letter A-E whose meaning lives only in the FO dispatch contract event-loop step-0 mapping:
- A = lingering (roster member, no live work)
- B = superseded
- C = un-advanced-pr
- D = stale-branch
- E = local-main-drift

The `reason` field already states the condition in English, so the letter adds an indirection the reader must resolve against the contract. The FO one-line drift summary (`A={N} B={N} C={N} D={N} E={N}`) is equally opaque.

## Proposed approach

Behavior decided this session — ideation formalizes the ACs, test plan, and doc-diff, not the design:
- Emit the descriptive name as `drift[].class`: `lingering` | `superseded` | `un-advanced-pr` | `stale-branch` | `local-main-drift`. Remove the letters.
- Update the FO dispatch contract event-loop step-0 — the per-class action mapping and the one-line summary format — to reference the names; the FO branches on the string name.
- A string enum is as stable to branch on as a letter, so there is no machine-stability cost; the gain is self-documenting JSON and a readable summary line (`stale-branch=3` instead of `D=3`).

## Out of scope

The drift-class set itself (no new classes, no semantics change); the reconcile detection logic; any other `dispatch` subcommand's output.

## Acceptance criteria

**AC-1 — `dispatch reconcile` output carries the descriptive class name, never a letter.** Each `drift[].class` is one of `lingering`/`superseded`/`un-advanced-pr`/`stale-branch`/`local-main-drift`.
Verified by: {a behavioral test over real reconcile output asserting the named classes and the absence of any bare A-E letter; ideation pins it.}

**AC-2 — The FO dispatch contract event-loop step-0 mapping and one-line summary reference the descriptive names, consistent with the helper emitted values.** The contract and the helper name the same class set — two independent sources that can diverge — not a prose-grep over the contract.
Verified by: {ideation pins the check — an independent-source consistency binding of the contract class set to the helper emitted enum, never a substring match over the contract.}

## Test plan

{Ideation fills. Likely a Go test over the reconcile helper output (the existing reconcile test surface) asserting the descriptive names, plus the contract/helper enum-consistency check. Record the contract step-0 doc-diff (before/after wording) in the body.}
