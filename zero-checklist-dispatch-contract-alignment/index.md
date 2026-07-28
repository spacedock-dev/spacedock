---
title: "Align zero-signal dispatch checklists between the FO contract and dispatch build"
status: backlog
source: "se0 local Codex keep-moving proof, 2026-07-28: a stage with no Outputs or naturally advancing AC produced the contract-valid empty checklist, while dispatch build rejected it and forced break-glass dispatch."
score: 0.72
sprint: durable-decisions
group: dispatch
sprint-readiness: ready
issue:
id: xxqk1kq7v8h2am9cvm6y8gyw
---

## Problem

The First Officer dispatch contract says `dispatch.checklist` may produce zero to three linchpin signals and must not pad. The shipped `spacedock dispatch build` contract rejects an empty checklist with `error: checklist must not be empty`. Both cannot be authoritative.

The mismatch occurred on a supported path during the `se0` local Codex keep-moving proof: the fixture's implementation stage declared no `Outputs` and no naturally advancing AC, so the FO correctly produced an empty checklist. All three helper calls failed and the FO entered break-glass dispatch. The keep-moving fixture now declares its real output because that journey is not a zero-checklist test; this task owns the product/contract mismatch rather than hiding it in that unrelated live scenario.

This is not a request for checklist padding. Ideation must choose one authority: either `dispatch build` accepts a well-formed zero-item checklist and emits an ergonomic assignment, or the FO contract requires at least one genuine linchpin and defines what happens when the workflow declares none. Preserve the agent-first rule that invented work does not become a ceremonial checklist item.

## Acceptance criteria

**AC-1 (VALUE) - A First Officer can dispatch a valid stage without inventing checklist work or falling into break-glass because two shipped contracts disagree.**
Verified by: an end-to-end dispatch-build fixture for a stage with no `Outputs` or naturally advancing AC follows the chosen contract and produces one unambiguous result.

**AC-2 - The FO contract, dispatch-build validation, generated assignment, CLI diagnostic, and tests state one rule for a zero-item checklist.**
Verified by: focused contract/CLI tests cover the chosen positive and the malformed-input negative without compatibility behavior.

**AC-3 - Non-empty checklist behavior and the three-item maximum remain unchanged.**
Verified by: existing dispatch golden/ergonomics tests plus one boundary case remain green.

## Boundary

Do not change `se0` product semantics. The `se0` keep-moving scenario retains a truthful implementation output and continues to test motion. This follow-up resolves the general zero-signal dispatch contract independently, with no backward-compatibility layer because v1 is unreleased.
