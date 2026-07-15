---
title: Generalize the Cowork binary bootstrap from survey into the shared launcher gate
status: backlog
source: captain request, live Cowork dogfood 2026-07-13
started:
completed:
verdict:
score:
worktree:
issue:
id: 8v5w1fk28m5bssm4147ad0ae
---

Make every spacedock skill front door work on a fresh Cowork project, not just survey. Today only survey (per survey-claude-cowork-runtime-detection, id eqn) owns the Cowork consent/install lifecycle for `.spacedock/bin/spacedock`; any other entry point that reaches the binary version gate on a Cowork VM aborts with host-wrong remediation.

## Problem

The FO/commission startup contract's version-gate abort for a missing binary says `brew install spacedock-dev/homebrew-tap/spacedock` and "launch with `spacedock claude`" — both wrong inside the Cowork VM (no brew, no host shell, no launcher). A fresh-project Cowork user who goes straight to `/spacedock:commission` gets a designed workflow whose first FO boot then dead-ends at that gate. The working bootstrap (positive Cowork capability probe → selected-folder two-surface binding → consent/Full-network relaunch → checksum-verified session install → exclusive-create of the exact mounted `.spacedock/bin/spacedock` → dual-surface verify) exists but is specified only inside survey.

## Proposed approach

Extract eqn's proven bootstrap into a shared reference the skills' binary version gate points at:

1. The version gate's binary-absent class gains a Cowork arm: when the Cowork capability probe is positive, route to the shared consent/install flow instead of the brew/launcher hint; when negative, keep today's host remediation unchanged.
2. `SPACEDOCK_BIN` resolution learns the exact mounted path: on Cowork, resolve `<selected-folder>/.spacedock/bin/spacedock` first (existing-state classes exactly as eqn: working reuses with zero network, broken stops without overwrite, absent enters consent).
3. Survey keeps only its evidence-adapter specifics; the bootstrap section becomes a pointer to the shared flow — one lifecycle, one prompt text, one set of failure classes (including COWORK_SHELL_UNAVAILABLE), no drift between front doors.

## Out of scope

- Changing eqn's design or re-proving the bootstrap (run 0/1 already live-proven; run 2 tracked by eqn).
- Git credential/push handling from the VM, and the state-checkout gitdir portability (filed separately as state-checkout-portable-gitdir).
- Non-Cowork runtimes' version-gate behavior.

## Acceptance criteria

**AC-1 - A fresh Cowork project reaching the binary version gate from any front door (first-officer, commission) is routed to the consent-gated bootstrap and never shown brew/`spacedock claude` remediation.**
Verified by: fixture replay of the gate's binary-absent class with a positive Cowork probe asserting the shared-flow routing and prompt text; a negative-probe control asserting today's host hint unchanged.

**AC-2 - After bootstrap, all front doors resolve the same exact mounted binary with zero re-download.**
Verified by: a second front-door invocation in the same replay observing reuse of `.spacedock/bin/spacedock` (zero network/install events), matching eqn's existing-state classes.

**AC-3 - The bootstrap lifecycle is specified once.**
Verified by: survey and the launcher gate both referencing the single shared flow document/section; a grep-level duplication check in contract lint or the skill smoke suite.

## Test plan

Fixture-level routing replays for the gate arms (cheap); reuse eqn's committed Cowork event fixtures where applicable; one live Cowork smoke rides along with eqn's run-2 session rather than a separate lane.
