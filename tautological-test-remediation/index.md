---
id: hwk58jy8akxhwzdydq8ztzrc
title: Remediate the 54 tautological tests — mutation-verify, then convert / re-bind / demote
status: ideation
source: "captain (2026-06-04) — the tautological-test sweep (Workflow w71il5awf, this session) found 54 of 61 instruction-file-assertion tests are tautological (banned as behavioral proof per the proof-policy f8b257cf). Concentrated in the contract-decomposition extraction tests + the hostneutrality/prose-lock suite. Captain: file and dispatch."
score: "0.32"
started: 2026-06-05T01:14:49Z
completed:
verdict:
worktree:
issue:
---

The tautological-test sweep found **54 of 61** instruction-file-assertion tests are tautological: they match a substring/absence over an instruction file the model ingests (the FO/ensign contracts, a workflow README, a skill) and so cannot fail meaningfully — a negation/paraphrase that keeps the substring passes. Per the proof-policy (`f8b257cf`), none can stand as behavioral proof.

## Problem (where they cluster)

- **Contract-decomposition extraction tests** (`skills/integration/present_gate_test.go`, `feedback_rejection_flow_test.go`, `using_claude_team_test.go`) — ~20 tests. The `…InvokesSkill` ones are worst: a behavioral claim ("the FO invokes present-gate at gate time") proven by a substring in the contract.
- **Hostneutrality / prose-lock suite** (`internal/hostneutrality/*` — dev-leakage, prose-neutrality, prose-inflator, split-root-sync, codex-runtime-contract) — ~14 tests. `TestFOHaltGateProse` / `TestFOSyncProse` / `TestCommissionJourneyProse` assert the contract SAYS the FO halts/syncs, never that it does.
- A handful in `skill_text_test.go` / `skill_surface_test.go`.

The contract-decomposition line shipped its "proof" on these; the only real proof was always the AC-3 live drives.

## Proposed approach

1. **Mutation-verify** the 54 (the sweep was static classification): for each, construct the meaning-inverting / rename / drift edit and confirm it stays green — confirm the count before acting.
2. **Triage each:**
   - **convert-to-behavioral-drive** — the `…InvokesSkill` / `…HaltGate` / `…SyncProse` cluster (they claim behavior → replace with a live/behavioral drive that RUNS it). Highest priority.
   - **re-bind-to-independent-source** — the banned-token / vocab absence checks → bind the token list to a code constant/manifest so it can fail on drift (becomes a legit invariant).
   - **demote-to-non-AC** — the "prose moved out of core / present in skill" extraction checks → keep as a text-consistency sanity check, NOT counted as AC satisfaction.
3. **Prioritize** the ones that currently stand as the SOLE proof of a behavioral AC.

## Acceptance criteria (ideation formalizes; per the policy each converted test's proof is RUNNING it, not another presence check)

**AC-1 (behavioral) — the converted `…Invokes…` / `…HaltGate` / `…Sync` tests are behavioral drives that red on the broken BEHAVIOR (not broken text).**
Verified by: each new drive demonstrated to red on a real behavior-break (a live or offline behavioral negative), mutation-controlled.

**AC-2 (invariant) — the re-bound tests bind to an independent source and red on drift.**
Verified by: each re-bound test demonstrated to red when the code-side source diverges from the file.

**AC-3 — no test remains as the sole proof of a behavioral AC via a presence/absence match over an ingested file.**
Verified by: a RE-RUN of the tautological-test sweep showing the count drops to the demoted-only set (each remaining presence check is an explicit non-AC sanity check, not claimed as behavioral proof).

## Test plan

- Mutation-verify pass (offline, per-test) → convert / re-bind / demote → the converted behavioral drives (some live-gated) → re-run the sweep as the AC-3 check.
- High-stakes (CI / shipped scaffolding) → detached adversarial audit before merge.
- Pairs with `eykb` (port the proof-policy to shipped scaffolding) — same policy, the contract-prose half; this is the test-suite half.

## Notes

Provenance: tautological-test sweep (this session, Workflow w71il5awf): 54/61 tautological, 5 invariant, 2 behavioral, 2 invariants flagged for review. Caveat: not all 54 are deletions — many DEMOTE to non-AC sanity checks; the goal is "no presence check stands as behavioral proof," not "no presence check exists." Likely phased by cluster (the `…Invokes…`/`…HaltGate` cluster first). Sibling: `eykb`, the dev-README fix `f8b257cf`.
