---
id: 2hdjgcy3g0y118hyymb1gwgw
title: Land the fo-gate-lifecycle frozen-binding/replay-idempotency wording within its byte cap
status: backlog
source: "Deferred from sonnet-gate-guardrail-no-authority (id 3zzpdw704df1g8pg1x9thzmw) implementation, 2026-08-02: ideation's mechanism 2 (two skills/fo-gate-lifecycle/SKILL.md sentence corrections) does not fit the existing 6600-byte component cap even at the tightest verbatim-preserving trim tried during implementation. Captain chose to defer rather than fold an unscoped prose trim or a cap change into that task as a side effect."
started:
completed:
verdict:
score: 0.5
worktree:
issue:
---

Make `skills/fo-gate-lifecycle/SKILL.md` explicit that (a) a rejected `--artifact`/`--reference` selection should be corrected and re-prepared once rather than escalated by probing with throwaway content against a real entity, and (b) a frozen gate-room binding is terminal — present or halt, never unfreeze by reverting or deleting binary-owned state — within the file's existing 6600-byte component cap, or via an explicit separate proposal to raise that cap.

## Problem

The parent task's live-Sonnet diagnosis (`sonnet-gate-guardrail-no-authority`, ideation stage report) found a real gate-lifecycle escape-hatch pattern: when `gate prepare` rejected an uncommitted selected source without naming which one, the FO bisected by binding throwaway `--question "test" --summary "test"` content to the real gate room, then — when a divergent rebind was correctly refused as frozen — escaped by `git checkout --` and `rm -rf`-ing the binary-owned `gates:` frontmatter and room, repeating the cycle before finally preparing correctly. `SKILL.md:27` already forbids hand-editing/deleting binary-owned authority, but doesn't yet say a frozen binding is terminal (present-or-halt) or forbid the throwaway-probe pattern by name; `SKILL.md:65`'s replay-idempotency line doesn't distinguish cross-session resume from a same-gate-entry re-read, which is what made the second probe look safe.

Mechanism 1 of the parent task (attributing the rejected path in `internal/gates/prepare.go`) already removes the ambiguity that triggered this specific incident and is landed. This follow-up is the reinforcement layer: closing the contract-text gap so the same escape-hatch pattern is explicitly forbidden even if some future ambiguity (not this exact one) arises.

## Proposed approach

Not yet ideated to a specific landing shape — that's this task's job. Starting material, carried over from the parent task's implementation attempt:

**Original (main, current):**
- Line 27: "The real capability check is the `gate prepare` invocation. Nonzero prepare halts; surface its exact error and refresh or rebuild the selected version-gated bundle when the command is unavailable. Never hand-edit `gates:` or delete/revert/replace binary-owned entity or room authority."
- Line 65: "Exact prepare replay is idempotent; a divergent open room is frozen—surface the refusal and stop."

**Ideation's proposed (verbatim, +uncounted bytes over cap as literal text):**
- Line 27 after: "The real capability check is the `gate prepare` invocation. When a nonzero prepare names a rejected `--artifact` or `--reference`, correct that one selection — commit the exact source, or drop it — and prepare once more; every other nonzero halts, surfacing its exact error and refreshing or rebuilding the selected version-gated bundle when the command is unavailable. Never diagnose by preparing throwaway `--question`/`--summary` content against a real entity: a successful prepare binds the room. Never hand-edit `gates:` or delete/revert/replace binary-owned entity or room authority; a frozen-binding refusal means the existing binding stands — present it or halt, never unfreeze it by reverting or deleting."
- Line 65 after: "Exact prepare replay is idempotent across sessions resuming an open gate; within one gate entry issue at most one successful prepare and treat its emitted lines as the binding — never replay to re-read them. A divergent open room is frozen—surface the refusal and stop."

**Implementation's tightest trim attempt (still +65B over cap after also cutting the pre-existing bundle-refresh clause — did not close the gap):**
- Line 27: "The real capability check is the `gate prepare` invocation. Rejected `--artifact`/`--reference`: fix it, reprepare once; other nonzero halts—surface it, refresh/rebuild if unavailable; never probe with throwaway content on a real entity (success binds the room); never hand-edit `gates:` or delete/revert/replace binary-owned authority — a frozen binding stands, present or halt, never unfreeze."
- Line 65: "**Resume.** Use boot/entity state and prior result; `gate validate`/`gate eligibility` are optional diagnostics, never positive-path requirements. One successful prepare per gate entry; reuse its emitted lines, never replay. A divergent open room is frozen—surface the refusal and stop. Require the exact Resolution commit before routing closed state. Pending approval → consume; revise/hold → route/stop; consumed → dispatch only if nonterminal, else merge; stale → supersede then replace. Surface nonzero command, exit, remedy; never repair frontmatter."

Baseline file was 6592 bytes against the 6600-byte cap (`internal/contractlint/fo_function_reference_invariant_test.go:14`) — 8 bytes of headroom for both edits combined. Ideation here should either (a) find wording that actually fits within 8 bytes net, (b) identify specific unrelated bytes elsewhere in the same file that are safe to cut to make room, with each cut justified on its own, or (c) write the separate evidence-backed proposal `docs/roadmap/durable-decisions/staff-review-sprint-close.md:358-362` reserves for changing a per-component cap, and let the captain decide there.

## Out of scope

Reopening mechanism 1 (`internal/gates/prepare.go` rejected-source attribution) or mechanism 3 (`internal/ensigncycle/livescenario_adapter_live_test.go` oracle-condition naming) — both already landed on the parent task. Any change to `gate prepare`/`record`/`consume` binary behavior. Raising or adding a cap for any file other than `skills/fo-gate-lifecycle/SKILL.md`. Weakening `assertRecordedGateHoldLog` or its negative fixtures.

## Acceptance criteria

**AC-1 (VALUE) - `skills/fo-gate-lifecycle/SKILL.md` explicitly forbids the throwaway-probe/destroy-room escape pattern and states a frozen binding is terminal, within its component cap (or a captain-approved raised cap).**
Verified by: the landed file text stating (a) correct-one-rejected-selection-then-reprepare-once, (b) never probe with throwaway content against a real entity, (c) a frozen binding is present-or-halt, never unfrozen by revert/delete — and `internal/contractlint/fo_function_reference_invariant_test.go` (`TestFOInstructionComponentCaps`) passing with the cap either unchanged or raised only via an explicit recorded captain decision cited in this entity.

**AC-2 - The replay-idempotency line distinguishes cross-session resume from same-gate-entry re-read.**
Verified by: the landed `SKILL.md` text at the former line 65 location stating idempotent replay applies across sessions resuming an open gate, and that within one gate entry only one successful prepare should occur, its emitted lines treated as the binding rather than replayed to re-read.

**AC-3 - Existing guardrails and unrelated lanes are unchanged.**
Verified by: `go test ./...`, `go test ./... -race`, `gofmt -l ./cmd ./internal` clean; no change to `assertRecordedGateHoldLog`'s accept/reject set or to any other capped component's byte count.

## Test plan

Offline only — no live run needed, this is contract text. Reproduce the byte-budget math first (current file size vs cap) before proposing wording, since that constraint is what defeated the first attempt. Run `go test ./...` (including `internal/contractlint`) plus `-race` plus `gofmt` after landing.
