---
id: 7ea4knxzvf3s4vve2zvr4ka0
title: Determine + document the intended team-vs-bare dispatch mode for headless `-p` runs
status: backlog
source: "0203-T3 surfaced (2026-06-14): TestLiveEnsignCycle flaked on a sonnet team-vs-bare coin-flip; captain steer: \"the important thing is determining the expected and intended behavior and document\""
started:
completed:
verdict:
score:
worktree:
issue:
---

The FO's team-vs-bare dispatch-mode choice for a **headless `-p` (non-interactive)** run is under-specified in the contract, so models coin-flip — which makes `TestLiveEnsignCycle` (the legacy full-cycle live smoke) intermittently fail. The deliverable is NOT to paper over it by accepting either mode; it is to **determine the intended behavior and document it** so the FO is deterministic and the test asserts the intended mode.

## Problem

The contract triggers single-entity **bare** mode on "`-p` **AND** the prompt names a specific entity" (`first-officer-shared-core.md` Single-Entity Mode), with bare's rationale being premature-session-termination safety in `-p` (`claude-first-officer-runtime.md`: the Agent tool without `team_name` blocks until completion). But:

- The TRIGGER ("names a specific entity") is a narrow proxy. A `-p` run with a single entity that the prompt does NOT name (exactly `TestLiveEnsignCycle`) falls in a gap: the trigger implies team, the safety rationale implies bare.
- Team-mode-IN-`-p` is genuinely fragile: the upstream claude-code premature-teardown bug forces the `antiShutdownOverride` hack (`internal/ensigncycle/claude_live_runner_test.go:23`, used at `live_test.go:122`) just to make a team survive a `-p` run. That hack is the evidence the team path is not robust under `-p`.
- So the two governing concerns pull apart: **premature-termination safety** favors bare in `-p`; **concurrency** (team's whole purpose) is only safe in an interactive session.

Result: the FO coin-flips, and the smoke gates on the coin (`isTeamCreate`) instead of on its real invariant (the dispatch→done cycle completes).

## Evidence (empirical, 2026-06-14)

Same T3 refs, `TestLiveEnsignCycle`:
- CI sonnet → **bare** (no TeamCreate) → FAIL (`live_test.go:186`).
- local sonnet ×2 → **team** → PASS. local opus ×1 → **team** → PASS.

So it is a low-frequency model coin-flip, NOT a deterministic regression. The smoke's end-state assertions (entity archived, `verdict`, path-scoped commit) held in EVERY run — the cycle completed regardless of mode.

## Proposed determination (FO analysis — confirm/refine at ideation gate)

Key the decision on **session interactivity + concurrency need**, not entity-naming:
- **Interactive session → team mode** (concurrency available; no premature-termination risk; the captain keeps the session alive).
- **Headless `-p` (non-interactive) → bare mode by default** (the upstream premature-teardown bug makes team fragile; a single-entity sequential drive needs no concurrency). This generalises the current "single-entity" articulation to its actual intent.

Consequence for the live smoke: a test that wants to exercise the **team-mode** cycle (TeamCreate → dispatch → bounded teardown) must FORCE team mode explicitly (a prompt cue), because the default `-p` behaviour is bare. Otherwise the smoke should assert the **bare** cycle (the `-p` default) and team-mode coverage lives in a dedicated, team-forced scenario.

## Proposed approach (to refine at ideation)

1. **Document the intent** unambiguously in the FO contract (the Single-Entity Mode / team-vs-bare prose): `-p` non-interactive → bare default; interactive → team; the rule keys on interactivity, not entity-naming.
2. **Prefer a code gate over prose** (this workflow's discipline): the launcher / `spacedock claude` already knows it is `-p` (non-interactive) — have it signal the dispatch mode deterministically so the FO does not coin-flip on a prose reading. A model-interpreted prose rule has a ceiling of "wording present"; a code-driven mode selection is the real guarantee.
3. **Align `TestLiveEnsignCycle`** to the documented intent — assert the intended mode's cycle completion (and/or split a team-forced scenario for team coverage). De-flakes the smoke.

## Acceptance criteria (sketch — ideation fleshes out, all external-proven)

- **AC-1** — `TestLiveEnsignCycle` is deterministic across repeated runs on both models (no team-vs-bare coin-flip failure). Verified by repeated live runs on the chosen mode.
- **AC-2** — The team-vs-bare decision is enforced by a code gate (or a real test), not a prose-only rule. Verified by a Go test that the dispatch mode resolves deterministically from the run context (`-p` vs interactive).
- **AC-3** — Team-mode coverage (TeamCreate → dispatch → bounded teardown) is still exercised by SOME live scenario (the one that forces team).

## Out of scope

The other live scenarios' bare/team handling beyond what this determination touches; the `comm-officer` polish-over-reach guard (separate); the merge-ref mod-block-section consolidation (separate T3 flag).

## Notes

Fast-follow surfaced by T3 (`fo-contract-prose-audit`); not a v0.20.3 blocker (T3's own behavior-preservation ACs passed on opus+codex). Captain may pull into a sprint or keep as fast-follow.

## Design determination (captain, 2026-06-14)

The intended behavior is settled — it is a **driving-mode** question, not just team-vs-bare. Collapse to two modes keyed off the single `-p`/interactive signal; gate-resolution is an explicit opt-in, NOT a property of the mode.

- **Interactive (no `-p`):** boot → greet → **STOP for input** (a human steers). Unchanged.
- **Headless `-p` (default):** boot → **drive all dispatchable work** → stop at the **first gate OR terminal** → **exit, reporting gate status.** No greet-stop. The FO does **not** decide gates (a gate is a human-owned decision; with no decision-maker present it is the natural stop boundary). This dissolves the auto-resolve-vs-blocked contradiction: default `-p` never resolves gates.
- **Headless `-p --auto-approve` (opt-in):** the FO is **given the conn** — resolves gates from the report verdict (PASS→advance; REJECT-with-`feedback-to`→bounce within the 3-cycle cap; REJECT-without-`feedback-to`/ambiguity/escalation→still stop+report) and drives to **end state (terminal).** This is the path the live-e2e harness sets to exercise a full feedback cycle, and a deliberate "drive it to done" operator choice.

**What this removes:** the greet-vs-drive coin-flip (in `-p` the FO always drives — deterministic), the `antiShutdownOverride` band-aid (it exists only to fight the greet-and-stop default under `-p`), and the fuzzy single-entity-mode special case (it is just `-p` scoped to a named entity; `--auto-approve` for the full-cycle test path). The contract boot step shrinks to one rule: interactive greets-and-stops; headless drives-to-gate-and-exits; `--auto-approve` drives-to-terminal.

**Killing the flake:** `TestLiveEnsignCycle`'s "FO subprocess exited (code=0) before TeamCreate matched" dies because `-p` now requires driving to first dispatch (no coin-flip). Separately, raise the **1-minute no-progress quiet budget** on the dispatch-close step (`live_test.go`) — a legit live ensign turn exceeds 1m (the second, independent flake signature: "dispatch close did not close within 1m0s").

**Team-vs-bare under this model:** orthogonal to the driving mode. `-p` drives regardless of team/bare; the team-survival fragility under `-p` (the upstream premature-teardown bug) is handled by the dispatch mechanism, not by greet-stopping. If team mode stays too fragile under `-p`, headless can dispatch bare — but that is a robustness choice, not the mode determination.

**ACs (ideation pins; behavioral, not prose-grep):** a live drive proving (a) `-p` with no `--auto-approve` drives to a gate and exits with gate status (no greet-stop, no gate decision); (b) `-p --auto-approve` drives a full PASS/REJECT cycle to terminal; (c) interactive greets-and-stops; plus the contract simplification (band-aid + special-case removed) and the quiet-budget bump. The live harness sets `--auto-approve` and asserts the dispatch→terminal cycle, not the `isTeamCreate` coin.

### Refinement (captain, 2026-06-14): auto-approve stays PROSE-based

`--auto-approve` is NOT a new launcher flag. It is a **prose mode** — the FO is *told / given the conn* to resolve gates via the prompt/contract, consistent with what already exists (`skills/commission/SKILL.md:22`: "if the user says to skip confirmation or auto-approve gates, proceed"). Keep the driving mode prose/prompt-determined — that is in line with the FO being a prose-driven contract agent; do NOT flag-ify it.

The flake's real cause is therefore narrower: the contract is **silent/ambiguous on `-p`**, so the model coin-flips greet-vs-drive. The fix is **deterministic contract prose**, not a code gate: state the two modes unambiguously so a `-p` FO always drives. Wherever the determination above says `-p --auto-approve`, read it as "`-p` + given-the-conn-to-auto-approve (prose)", not a parsed flag. The contract-prose change is the deliverable; its behavioral proof is the live drive (FO deterministically drives under `-p`), not a prose-grep over the clause.
