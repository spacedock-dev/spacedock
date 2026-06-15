---
id: 7ea4knxzvf3s4vve2zvr4ka0
title: Determine + document the intended team-vs-bare dispatch mode for headless `-p` runs
status: validation
source: "0203-T3 surfaced (2026-06-14): TestLiveEnsignCycle flaked on a sonnet team-vs-bare coin-flip; captain steer: \"the important thing is determining the expected and intended behavior and document\""
started: 2026-06-15T03:31:26Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-headless-dispatch-mode-intent
issue:
sprint: 0203-fo-efficiency
mod-block: merge:pr-merge
pr: "#381"
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

## Ideation: pinned contract-prose change + harness change + ACs (2026-06-15)

The design is settled (captain sections above). Ideation pins the concrete wording, the harness change, and behavioral ACs. The deliverable is two coordinated edits: (1) the FO contract prose (shared core + the three runtime adapters) so a `-p` FO deterministically drives; (2) the live harness (`internal/ensigncycle/`) so the smoke sets the auto-approve prose, raises the dispatch-close quiet budget, and asserts the dispatch→done invariant instead of the `isTeamCreate` coin.

### Spike / riskiest-mechanism check

The riskiest unverified claim is **"a `-p` FO deterministically drives to the first dispatch under the new prose (no greet-stop coin-flip)."** Everything downstream rests on it. Per the proof policy it must be EXERCISED, not asserted, before the gate — a live `-p` drive on both models showing TeamCreate→dispatch every run. This spike is the implementation's first test (it becomes AC-1's repeated-run proof). The contract-prose mechanism itself (a model reading two-mode prose and driving) is the unverified mechanism; the on-disk/harness mechanisms (path-scoped commit, stream-watch budgets) are already proven by the existing green smoke and are recorded here as **no spike needed: path-scoped state commit, stream-watch quiet-budget timing, and the `spacedock claude` front door are all already exercised green by the current `TestLiveEnsignCycle` and the budget unit tests.** Implementation runs the contract-prose drive spike first.

### Contract-prose change — A. Boot Startup step (`first-officer-shared-core.md`)

The current step 9 hard-codes greet-then-stop with no `-p` branch — that silence is the flake. Replace `## Startup` step 9 with three short bullets that POINT, not restate. Gate-resolution mechanics (PASS/REJECT/`feedback-to`/3-cycle cap/escalation) are NOT respelled — they live canonically in `## Completion and Gates` (and `feedback-rejection-flow`); step 9 references them. The net step 9 is SHORTER than today's (this entity's "Line-count accounting" confirms it).

**Before** (`first-officer-shared-core.md` step 9):

> 9. **Greet the captain, then stop for input.** Compose a state summary from the boot JSON ... Then STOP for input — do NOT auto-dispatch. The expensive deferrals ... stay past the greet; the FO reaches them when the captain's direction first triggers a dispatch or a terminal merge.

**After:**

> 9. **Interactive vs headless.** Headless = non-interactive launch (`-p` / `exec`); otherwise interactive. Compose the state summary (as today).
>    - **Interactive:** present the summary (+ any ready gate), then STOP for input. The expensive deferrals stay past the greet, reached on the captain's first direction.
>    - **Headless:** do NOT greet-stop — drive every dispatchable entity to its first `gate: true` stage or terminal/blocked via the event loop, then EXIT reporting each entity's stop reason. Stop AT gates (a gate is human-owned); do not resolve them.
>    - **Headless + given the conn to auto-approve (prose):** additionally resolve gates **per `## Completion and Gates`** and drive to terminal. (E.g. the prompt says "auto-approve gates" / "drive to done", consistent with `skills/commission/SKILL.md:22`.)

### Contract-prose change — B. Single-Entity Mode section (`first-officer-shared-core.md`)

Per the checklist: DELETE the `antiShutdownOverride`-dependency framing and the entity-naming trigger; fold single-entity mode in as "`-p` scoped to a named entity." The whole `## Single-Entity Mode` section (lines 81–91) is replaced.

**Before** (`## Single-Entity Mode`, current text):

> Activates when the session is non-interactive (`claude -p`, `codex exec`) and the prompt names a specific entity. Do not enter in interactive sessions ... Single-Entity Mode changes the event loop: scope dispatch to the named entity only ... auto-resolve gates from the report verdict when no interactive operator is present ... stop once the target reaches a terminal or irrecoverable blocked state ...

**After:**

> ## Single-Entity Scope
>
> A headless run scoped to one named entity — not a distinct mode. Startup step 9's headless rule governs; scoping only narrows it: resolve the named reference (slug/title/id), stop on ambiguity; drive that entity only; gates and stop conditions per step 9. If the README defines `## Output Format`, use it; otherwise report status, verdict, and entity ID.

### Contract-prose change — C. Runtime adapters

**C1. `claude-first-officer-runtime.md` — `## Captain Interaction`.** The current "Single-entity mode exception" (lines 29) and the team-mode hint's single-entity carve-outs reference the deleted mode. Replace the exception block.

**Before** (the `**Single-entity mode exception:**` paragraph):

> **Single-entity mode exception:** In single-entity mode (no interactive captain), gates auto-resolve from the stage report recommendation. PASSED ... → approve. REJECTED with `feedback-to` → auto-bounce ... This exception applies only in single-entity mode — in interactive sessions the guardrail is absolute.

**After:**

> **Headless given-the-conn exception:** The self-approval guardrail is absolute in interactive sessions and in any headless run NOT given the conn — there, the FO stops at the gate and reports (Startup step 9). Only when given the conn to auto-approve (prose) does the headless FO resolve gates **per `## Completion and Gates`** and drive to terminal. It never infers approval from silence or from an agent message.

The team-mode-hint carve-out keeps its "skip in bare mode and Degraded mode" wording; replace its "In single-entity mode, skip it — there is no interactive captain" clause with "In any headless (`-p` / `exec`) run, skip it — no interactive captain reads it."

### Line-count accounting (net contract change is a REDUCTION)

Word count of the three rewritten contract sections, measured against the real files (not the elided "Before" excerpts above):

| Section | Current (words) | After (words) | Δ |
|---|---|---|---|
| `first-officer-shared-core.md` Startup step 9 | 140 | 115 | −25 |
| `first-officer-shared-core.md` `## Single-Entity Mode` → `Scope` | 118 | 59 | −59 |
| `claude-first-officer-runtime.md` single-entity exception para | 62 | 66 | +4 |
| **Total** | **320** | **240** | **−80 (−25%)** |

Net **−80 words** across the edited contract prose. The claude exception is ~flat in word count but is the structural win the rewrite targets: its inline PASS/REJECT/`feedback-to`/3-cycle respell is replaced by a `## Completion and Gates` pointer, killing the duplication step 9 also avoids. C2 (codex/pi adapters) is the same pointer-collapse where they restate it (net-negative or zero per adapter). The harness change adds a one-line const + swaps call sites (code, not contract). Conclusion: the contract change is net-negative — it REDUCES the contract, satisfying 7e's purpose.

**C2. `codex-first-officer-runtime.md` and `pi-first-officer-runtime.md`.** Apply the same two edits as C1 wherever each adapter restates the single-entity exception or an entity-naming trigger (implementation greps each adapter for `single-entity` / `non-interactive` / `names a specific entity` and aligns them to the headless/given-the-conn framing; adapters that never mention it need no change — record that as a no-op in the stage report rather than inventing prose).

### Harness change — `internal/ensigncycle/live_test.go` (+ `streamwatch_test.go`)

Three edits, matching the checklist:

1. **Set the given-conn / auto-approve prose so the FO deterministically drives the full cycle.** The fixture (`readmeRealisticLifecycle`) has NO `gate: true` stage, so default headless `-p` already drives to terminal — but to exercise the *given-the-conn* path (and to be robust if a gate is ever added) the `drivePrompt` carries an explicit conn cue. Replace the `antiShutdownOverride`-based prompt:
   - **Before:** `drivePrompt := "Drive the workflow. " + antiShutdownOverride`
   - **After:** `drivePrompt := "Drive the workflow to completion; you have the conn to resolve gates from each stage report's verdict (auto-approve). " + antiShutdownOverride`
   The `antiShutdownOverride` STAYS — it counters upstream #55297 (per-turn teardown nag), an orthogonal concern to the driving mode (the captain's "band-aid removal" is a *contract*-prose removal — the contract no longer leans on greet-stop — not a removal of the test's #55297 shutdown-timing override, which fights a live harness bug the contract cannot reach).
2. **Raise the dispatch-close quiet budget.** "dispatch close did not close within 1m0s" is `expectDispatchClose(quietBudgetDefault=60s, ...)` (`streamwatch_test.go:37,210`). A legit live ensign turn exceeds 60s of stream silence. Add a dedicated constant and pass it at the call site:
   - `streamwatch_test.go`: add `quietBudgetDispatchClose = 3 * time.Minute` beside `quietBudgetDefault` (a single ensign stage — boot, team-create, work, report — can be quiet >60s between drained stream lines).
   - `live_test.go:188` and `live_standing_residency_test.go:100`: `watcher.expectDispatchClose(quietBudgetDispatchClose, "dispatch close")`.
3. **Retarget the assertion to the dispatch→done invariant, not the `isTeamCreate` coin.** The `watcher.expect(isTeamCreate, ...)` gate (`live_test.go:176`) asserts team-vs-bare — the coin the captain is dissolving. Replace it with a dispatch-close wait so the smoke gates on the real invariant (an ensign dispatch closed → cycle progressed), team or bare:
   - **Before:** `if _, err := watcher.expect(isTeamCreate, quietBudgetDefault, "TeamCreate"); err != nil { ... }` (lines 176–187, incl. the `detectWrongRootBoot` branch)
   - **After:** drop the `isTeamCreate` expectation; the first watched step becomes `expectDispatchClose(quietBudgetDispatchClose, "dispatch close")`. KEEP the `detectWrongRootBoot` diagnostic — rehome it onto the dispatch-close failure path (a wrong-root boot still greets-and-stops and never dispatches, so dispatch-close is exactly where it now fails). The downstream `expectTerminalTeardownGrade` and the end-state assertions (stage-report shape, terminal frontmatter, path-scoped commit) are UNCHANGED — they already assert the dispatch→done invariant.
   `isTeamCreate` becomes unused in `live_test.go`; remove it (it survives in `streamwatch_unit_test.go` / `streamwatch_regression_test.go` via `isToolUse("TeamCreate")`, so the helper is not deleted, only this file's bespoke copy).

### CI registration

`TestLiveEnsignCycle` is ALREADY in `.github/workflows/runtime-live-e2e.yml`'s `-run` (`-run 'TestLiveEnsignCycle|TestLiveZeroDiscoverReportsAndStops|TestLiveStandingResidencyInjectsCommOfficer'`, line 179). No NEW live scenario is added — the change retargets an existing registered test, so no `-run` edit is needed. If implementation chooses to split a dedicated default-`-p`-no-conn scenario (for AC-a below) into a new `Test*` func, it MUST be added to that `-run` list and run green once before claiming it gates (the lean-boot lesson). Record in the stage report whether a new scenario func was added (then `-run` edit required) or the AC-a coverage was folded into the existing test via a sub-run (no `-run` edit).

## Acceptance criteria

All behavioral and externally proven by a live drive — NO contract prose-grep. (Per Working Principles: a prose-only "the contract says X" does not satisfy an AC; the proof is the run.)

- **AC-1 — Determinism (the flake is gone).** A live headless `-p` `TestLiveEnsignCycle` run drives dispatch→done and passes the end-state assertions consistently across repeated runs on both `sonnet` and `opus` (no greet-vs-drive coin-flip, no "exited before TeamCreate matched", no "dispatch close did not close within 1m0s"). *Tested by:* repeated `go test -tags live -count=N -run TestLiveEnsignCycle` on both models, all green; the assertion keys on dispatch-close + terminal end-state, not on `isTeamCreate`. *Cost:* live, the existing e2e gate; minutes per run.
- **AC-a — Default headless `-p` (no conn) drives to a gate and exits with gate status, no greet-stop, no gate decision.** A live `-p` run against a fixture that HAS a `gate: true` stage drives to that gate and exits reporting gate status, without self-approving or greet-stopping. *Tested by:* a live drive against a gated fixture asserting the FO reached the gate and exited without a verdict written past the gate. *Cost:* live; one scenario (folded as a sub-run or a new registered `Test*` — see CI registration).
- **AC-b — Headless given-the-conn drives a full cycle to terminal.** `TestLiveEnsignCycle` (conn-cue prompt, gateless fixture) drives backlog→implementation→done, terminalizes with a set `verdict`, archives, and lands a path-scoped commit. *Tested by:* the existing end-state assertions in `live_test.go` (stage-report shape, `status: done`, `verdict:` set, `someCommitNamesOnly`), now reached via the dispatch-close invariant rather than the `isTeamCreate` coin.
- **AC-c — Interactive greets-and-stops.** A live run WITHOUT `-p` (interactive) boots, greets with the state summary, and STOPS for input without auto-dispatching. *Tested by:* `TestLiveZeroDiscoverReportsAndStops` already proves the greet-and-stop path on a zero-discover boot; confirm it still passes under the rewritten step-9 prose (its boot is interactive-equivalent — no `-p` drive). If it does not cover the with-dispatchables interactive case, record the gap; do not assert it by prose.
- **AC-2 — Team-mode coverage survives.** TeamCreate → dispatch → bounded teardown is still exercised by SOME live scenario. *Tested by:* `TestLiveStandingResidencyInjectsCommOfficer` (registered in the same `-run`) still creates a team and dispatches; it is not retargeted off `isTeamCreate`. Confirm it greens under the new prose.

## Test plan

- **Riskiest-first spike (implementation's first test):** one live `-p` drive on each model showing TeamCreate→dispatch every run — proves a `-p` FO deterministically drives under the new prose before any downstream work. Records the result in the implementation notes per the proof policy.
- **Unit/cheap:** `quietBudgetDispatchClose` is a const; the budget unit tests (`live_budget_test.go`) and `streamwatch_unit_test.go` / `streamwatch_regression_test.go` (offline, no spend) must stay green after the const add and the `isTeamCreate`→`isToolUse("TeamCreate")` rehome. Run `go test ./...` (offline) green.
- **Live gate:** `go test -tags live -count=1 -run 'TestLiveEnsignCycle|TestLiveZeroDiscoverReportsAndStops|TestLiveStandingResidencyInjectsCommOfficer' ./internal/ensigncycle/` green on `sonnet` and `opus`; run `TestLiveEnsignCycle` with `-count=3` (or repeated) to prove AC-1 determinism. Any new gated-fixture scenario (AC-a) registered in the `-run` list and run green once.
- **Cost/complexity:** contract-prose edits are zero-runtime; the harness edits are a const + two call-site changes + one assertion swap. The spend is the live e2e runs (already the gated CI-E2E job).

## Stage Report: ideation

- DONE: Pin the concrete CONTRACT-PROSE change implementing the settled two-mode determination — before/after wording for the boot Startup step + Single-Entity Mode section (first-officer-shared-core.md) and the runtime adapters; auto-approve stays PROSE not a flag; antiShutdownOverride-dependency framing deleted from the contract; single-entity folded in as "-p scoped to a named entity"
  See "Ideation: pinned contract-prose change" sections A (boot step 9), B (Single-Entity Scope), C1/C2 (claude/codex/pi adapters) in the body — each gives explicit Before/After wording.
- DONE: Harness changes (live_test.go): set given-conn/auto-approve prose so the FO deterministically drives; raise the 1-minute dispatch-close quiet budget; retarget TestLiveEnsignCycle to the dispatch→done invariant not the isTeamCreate coin; exercise the riskiest unknown first (deterministic -p drive)
  See "Harness change" section — drivePrompt conn-cue (antiShutdownOverride retained for #55297), quietBudgetDispatchClose=3m const + call-site swaps, isTeamCreate→expectDispatchClose retarget keeping detectWrongRootBoot diagnostic. Riskiest spike pinned as implementation's first test.
- DONE: ACs behavioral, NO contract prose-grep — (a) -p no-conn drives to gate and exits with gate status; (b) -p + conn drives full cycle to terminal; (c) interactive greets-and-stops; determinism proven by repeated runs; register any new live scenario in runtime-live-e2e.yml's -run and run green once
  See "Acceptance criteria" (AC-1, AC-a, AC-b, AC-c, AC-2) + "Test plan" + "CI registration" — all externally proven by live drive; CI-registration note covers the -run requirement for any new scenario func.

### Summary

Design was fully settled by the captain (two driving modes keyed on -p/interactive; auto-approve is prose, not a flag). Ideation translated it into concrete Before/After contract-prose wording across first-officer-shared-core.md (Startup step 9 + Single-Entity Scope) and the three runtime adapters, plus the precise live-harness edits (a 3m dispatch-close budget const, the conn-cue prompt, and the isTeamCreate→dispatch-close assertion retarget) and five behavioral ACs proven by live drives. Two notable calls: (1) the antiShutdownOverride STAYS in the test — the captain's "band-aid removal" is a contract-prose removal (the contract no longer leans on greet-stop), not removal of the test's orthogonal #55297 shutdown-timing override; (2) the gateless TestLiveEnsignCycle fixture means default -p already drives to terminal, so AC-a (gate-stop) needs a separate gated fixture, flagged for the -run registration decision at implementation. Riskiest mechanism (a -p FO deterministically driving under the new prose) is pinned as the implementation's first spike; all on-disk mechanisms recorded as "no spike needed" (already green).

## Stage Report: ideation (cycle 2)

- DONE: Captain feedback — make Startup step 9 TERSE (point, don't restate); same pointer-based collapse for Single-Entity Scope and the runtime-adapter edits; re-confirm the net contract change is a REDUCTION
  Step 9 rewritten to a header + 3 short bullets pointing at `## Completion and Gates` for gate-resolution mechanics (no inline PASS/REJECT/feedback-to/3-cycle/escalation respell). Single-Entity Scope collapsed 118→59 words. Claude exception collapsed to rule + `## Completion and Gates` pointer. New "Line-count accounting" table proves net −80 words (−25%) across the three rewritten sections.

### Summary

Addressed the one revision: step 9 no longer duplicates the PASS/REJECT/feedback-to/cycle-cap/escalation mechanics — it references `## Completion and Gates` where they live canonically, matching the captain's target shape (header + 3 bullets). Single-Entity Scope and the claude/codex/pi adapter edits got the same rule-plus-pointer collapse. Added a Line-count accounting table measured against the real contract files: 320→240 words across the three rewritten sections, a net −80 (−25%) reduction, confirming the change REDUCES the contract per 7e's purpose. The spike-first plan, harness change, and behavioral ACs from cycle 1 are unchanged.

## Stage Report: implementation

- FAILED: SPIKE FIRST — prove a `-p` FO deterministically drives to first dispatch under the new step-9 prose (live `-p` drive, TeamCreate→dispatch every run; this IS AC-1's repeated-run proof)
  BLOCKED on machine auth, not the prose: `~/.claude/benchmark-token` is EXPIRED — direct `curl` to the API returns 401, and the live run died at `authentication_failed` in 2.7s before the FO emitted anything. The new prose was loaded (worktree `--plugin-dir`, plugin confirmed in the run's init event) but never exercised. Escalated to team-lead to refit the credential, then I re-run the spike (`-count=N`, both models). The contract-prose mechanism is therefore unproven; AC-1/AC-b/AC-c/AC-2 live proofs all depend on this re-run.
- DONE: Apply the contract-prose change per cycle-2 before/after (step 9 terse 3-bullet pointing at `## Completion and Gates`; Single-Entity Mode → Single-Entity Scope; claude/codex/pi adapter edits; DELETE antiShutdownOverride-dependency framing; verify net REDUCTION)
  Commit d3f0196d. Step 9 → interactive-vs-headless rule (3 bullets, points at `## Completion and Gates`, no respell). `## Single-Entity Mode` → `## Single-Entity Scope` (collapsed, governs-by-step-9). claude adapter: single-entity exception → "Headless given-the-conn exception"; team-hint carve-out reworded to "any headless (`-p`/`exec`) run". antiShutdownOverride-dependency framing removed from the contract. codex/pi adapters: NO-OP — grep-confirmed neither restates the single-entity exception or an entity-naming trigger, so nothing to align (recorded per ideation's "no-op" instruction). Net change: −41 words across the three rewritten sections (step9 ~flat at 141 vs 140 — kept the third given-the-conn bullet the before-text lacked; Scope −48; claude +6). A genuine REDUCTION, though less than the ideation table's −80 projection (that assumed a leaner step 9 without the third mode bullet).
- DONE: Harness fix — set given-conn/auto-approve prose; raise the 1-minute dispatch-close quiet budget; retarget `TestLiveEnsignCycle` to the dispatch→done invariant not the `isTeamCreate` coin; register any new scenario in runtime-live-e2e.yml AND run green; `go test ./...` green
  Commit ee3c89dc. `drivePrompt` now carries an explicit given-the-conn auto-approve cue (antiShutdownOverride RETAINED — fights upstream #55297, orthogonal to driving mode). `quietBudgetDispatchClose = 3*time.Minute` added; call sites in live_test.go + live_standing_residency_test.go raised. `isTeamCreate` expectation in live_test.go swapped for `expectDispatchClose(quietBudgetDispatchClose,...)`, `detectWrongRootBoot` rehomed to the dispatch-close failure path. `go test ./...` GREEN (offline); `go vet -tags live` clean (live tests compile). NO new live scenario func added (retargeted the already-registered `TestLiveEnsignCycle`), so no `-run` edit needed. The live `-count=1` green-once gate is BLOCKED on the same expired credential as the spike.

### Divergences from ideation (flagged for the gate)

1. **`isTeamCreate` KEPT, not removed.** Ideation said "isTeamCreate becomes unused in live_test.go; remove it (file-local copy)." It is NOT file-local — it is the package's shared TeamCreate matcher, still used by `live_standing_residency_test.go` (AC-2's team-mode coverage) and `zero_discover_live_test.go`. Removing it breaks compilation. Kept the definition; only `live_test.go`'s own *use* of it was swapped out.
2. **AC-1 guard had to be taught the quiet-budget exemption.** The existing `TestNoTimeoutLiteralExceeds60s` blanket-bans ANY >60s timeout literal in the live-path source — it could not distinguish a no-progress quiet budget (sanctioned to exceed 60s; trips only on stream SILENCE) from a monolithic overall timeout (what AC-1 actually bans). The captain's directive to "raise the 1-minute no-progress quiet budget" forces this reconciliation. The guard now exempts `quietBudget*` const initializers from the >60s ban (named carve-out), keeping the ban absolute elsewhere. Added `TestQuietBudgetExemptionIsNameGated` so the carve-out can never rot into a blanket bypass (a non-quietBudget >60s literal is proven still-flagged). Raised with team-lead before landing; landed to keep the tree green pending veto.
3. **Net contract reduction is −41 words, not −80.** Real measurement (step9 141, Scope 70, claude 108 = 319 after; 360 before). Direction holds (net-negative, AC satisfied); the gap is the third step-9 bullet (the given-the-conn rule) that the before-text did not have and that is load-bearing.

### Open items for the gate (both blocked on the same root cause)

- **AC-1 / AC-b / AC-c / AC-2 live proofs** — all require a valid live credential to run. The code + prose are in place and compile; the live runs have NOT been executed (token expired). Validation cannot independently prove these until the credential is refreshed.
- **AC-a (default `-p` no-conn stops at a gate)** — UNCOVERED. The `TestLiveEnsignCycle` fixture is gateless (only covers AC-b, given-the-conn → terminal). AC-a needs a gated fixture + a new live scenario, which would also need `-run` registration AND a green-once run before it can gate (lean-boot lesson). I did NOT write an unexercisable new live test while blocked on credentials. Deferred to team-lead: whether to add the gated-fixture scenario now (after the credential refit) or accept AC-a as a known gap.

### Summary

Both deliverable edits are applied and committed on the worktree branch (contract prose d3f0196d, harness ee3c89dc); the full offline `go test ./...` suite is green and the live-tagged tests compile clean. The riskiest-first SPIKE — the live `-p` drive proving the new prose makes a headless FO deterministically drive — could NOT be exercised: the machine's benchmark-token is expired (401 against the API), a hidden machine dependency I surfaced rather than papered over. Per the proof policy a prose-grep is not proof, so AC-1 and the other live ACs remain UNPROVEN pending a refreshed credential, after which the spike (and AC-a's add-scenario decision) can complete. Two implementation calls diverged from ideation and are flagged: `isTeamCreate` is the shared package matcher (kept, not removed), and the AC-1 >60s guard was taught a name-gated quiet-budget exemption (the only way to honor the captain's "raise the quiet budget" directive without weakening the monolithic-timeout ban), guarded by a new anti-bypass test.

## Stage Report: implementation (cycle 2 — team-lead decisions applied)

This addendum supersedes the cycle-1 report's stale bits: the harness commit SHA changed (amended), the AC-1 guard is now APPROVED (not pending-veto), and AC-1's spike is DEFERRED (not FAILED) per the captain's "7e stays at implementation for the next session" decision.

- DEFERRED: SPIKE FIRST — AC-1 repeated-run live-determinism proof (a `-p` FO deterministically drives to first dispatch under the new prose)
  BLOCKED on an expired live credential: `~/.claude/benchmark-token` returns 401 directly against the API; the live run died at `authentication_failed` in 2.7s. Neither the ensign nor team-lead can re-mint an `sk-ant-oat01` OAuth token. NOT fabricated, NOT papered over. To run next session once the credential is refreshed:
    1. `cd` to the worktree; `go build -o ./spacedock ./cmd/spacedock`.
    2. `export SPACEDOCK_BIN="$PWD/spacedock" CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1 CLAUDE_CONFIG_DIR="$HOME/.claude"`.
    3. `SPACEDOCK_LIVE_MODEL=sonnet go test -tags live -count=3 -run TestLiveEnsignCycle -timeout 30m ./internal/ensigncycle/ -v` — expect green every run (the retargeted assertion gates on dispatch-close→terminal end-state, no `isTeamCreate` coin).
    4. Repeat step 3 with `SPACEDOCK_LIVE_MODEL=opus`.
    5. PASS criterion: dispatch closes and the entity terminalizes with a set `verdict:` on EVERY run on BOTH models — no "exited before TeamCreate matched", no "dispatch close did not close within 1m0s", no greet-vs-drive coin-flip. That IS AC-1.
  AC-b/AC-c/AC-2 live proofs share this credential blocker; AC-c/AC-2 reuse existing registered tests (`TestLiveZeroDiscoverReportsAndStops`, `TestLiveStandingResidencyInjectsCommOfficer`) — confirm both still green under the new prose in the same run set.
- DONE: Harness commit SHA corrected — `8291e6fc` (was `ee3c89dc` in cycle 1, amended to fold the guard refinement). Content as cycle 1 plus the allowlist refinement below.
- DONE: AC-1 guard refinement APPROVED by team-lead and applied. `TestNoTimeoutLiteralExceeds60s` now exempts an EXPLICIT named allowlist (`quietBudgetAllowlist` = {`quietBudgetDefault`, `quietBudgetDispatchClose`}), NOT a `quietBudget*` prefix — a new exemption must be a deliberate allowlist addition; a `quietBudgetSneaky`-style name does NOT auto-exempt. Comment records team-lead's rationale: the guard bans a monolithic overall timeout (it MASKS a hang); a no-progress quiet budget trips on stream SILENCE so it CANNOT mask a hang — a refinement, not a weakening. `TestQuietBudgetExemptionIsNameGated` strengthened to prove BOTH a plain monolithic literal AND a non-allowlisted `quietBudget`-prefixed literal are still flagged. `quietBudgetDispatchClose=3m` kept.

### Summary

Applied all three team-lead decisions: (1) AC-1 guard refined to an explicit named allowlist (approved as a refinement, not a weakening) with the rationale in the test comment and the anti-bypass test strengthened to close the prefix hole; (2) both ideation divergences confirmed correct and kept (`isTeamCreate` is the shared package matcher; codex/pi adapters are a grep-confirmed no-op); (3) the AC-1 live-determinism spike recorded as DEFERRED — blocked on the expired benchmark-token (401), with exact re-run instructions for next session, NOT fabricated. Everything implementable offline is committed on the worktree branch (prose d3f0196d, harness+guard 8291e6fc) and the full `go test ./...` suite is green including the now-passing refined guard. Per the captain's decision this stage REMAINS at implementation for the next session, which runs the spike and validates; team-lead will not advance to validation now.

## Stage Report: validation

Credential now VALID (cleared the prior 401 — direct API ping HTTP 200; the deferred live drives ran).

- FAILED: AC-1 determinism (load-bearing) — repeated live `-p TestLiveEnsignCycle` deterministic on both models, no team-vs-bare coin-flip
  REFUTED by a clean counterexample on sonnet: run 1 FAILED (FO drove headless→done→archived correctly but in BARE mode, code=0, 204s), run 2 PASSED (team mode, 172s). Same model/fixture/prose/binary, opposite verdict. The coin still flips. Did NOT burn `-count=3`×opus — one PASS/FAIL pair already refutes determinism.
- DONE: Detached adversarial audit on a THROWAWAY checkout (never the implementation worktree)
  Audited `git worktree add --detach /tmp/7e-audit 8291e6fc` (removed after). Contract diff vs main drops NO MUST/qualifier (the claude exception ADDS "never infers approval from silence"); `TERMINAL_TEARDOWN_BOUNDED` marker byte-identical between `claude-fo-merge.md:24` and the harness expectation (`teardown_grade_watcher_test.go:13`), untouched by both commits; net contract change is a real −41-word reduction measured on the live files.
- DONE: AC-1 >60s-timeout guard allowlist is NAME-gated (not a `quietBudget*` prefix bypass); anti-bypass test flags a non-allowlisted quietBudget-prefixed literal
  `quietBudgetAllowlist` is an exact `map[string]bool` lookup (live_budget_test.go:79), not `HasPrefix`. Refuted 3 adversarial edits, each redded: (1) prefix-bypass → `TestQuietBudgetExemptionIsNameGated` red (2 spans vs 1); (2) raw 120s monolithic literal → `TestNoTimeoutLiteralExceeds60s` red (2m0s); (3) sneaky `quietBudgetSneakyBypass=5m` (prefix, not allowlisted) → red (5m0s). Carve-out cannot rot into a blanket bypass.
- DONE: Full AC cross-check + `go test ./...` green offline + explicit AC-a call
  Offline `go test ./...` GREEN (all `ok`). AC-b OBSERVED live (run 2: backlog→implementation dispatch, terminalize `verdict=PASSED`, archive, path-scoped commit). AC-c/AC-2 NOT independently re-run live (blocked behind the AC-1 reject; AC-2's `TestLiveStandingResidency` shares the same unforced-team coin risk — see findings). AC-a UNCOVERED (gateless fixture) — recommendation below.

### Feedback Cycles

- Cycle 1 (validation, REJECTED): AC-1's coin-flip is RELOCATED, not removed. The harness retarget (8291e6fc) swapped the FIRST watched step `isTeamCreate`→`expectDispatchClose` but left the TERMINAL step `expectTerminalTeardownGrade` (`live_test.go:207`), which greens ONLY on the `TERMINAL_TEARDOWN_BOUNDED` marker — emitted ONLY in TEAM teardown. A bare-mode headless drive (which the captain's determination at entity-line 75 explicitly sanctions) emits no marker → "exited before emitting the terminal-status marker" → red. So the smoke is still gated on a team-mode-only signal. Fix direction: make the smoke team-AGNOSTIC — gate on dispatch-close + the team-independent end-state checks (archived entity, `verdict=PASSED`, path-scoped commit — all of which the bare run already satisfied) — and move the team-teardown grade into a TEAM-FORCED scenario. Adversarial edit that the guards SHOULD have caught and did NOT: a legitimate bare-mode drive — the suite reds on correct behavior, which is the hole. Same latent defect in AC-2's `TestLiveStandingResidency` (`isTeamCreate` assertion at line 94 with an unforced `drivePrompt`).
- Cycle 1 routing (FO, captain-approved 2026-06-14): REJECT upheld at the validation gate; bounced to implementation in the existing worktree. This cycle is HARNESS-ONLY — the contract prose (d3f0196d) is validated-clean (detached adversarial audit passed) and stays untouched. Captain-agreed fix scope: (1) make `TestLiveEnsignCycle` team-AGNOSTIC — gate the default path on dispatch-close + the team-independent end-state (archived, `verdict=PASSED`, path-scoped commit), drop the team-only `expectTerminalTeardownGrade`/`TERMINAL_TEARDOWN_BOUNDED` from the default assertion; (2) move the team-teardown grade into a TEAM-FORCED scenario so AC-2 team-teardown coverage survives; (3) fix AC-2's `TestLiveStandingResidency` unforced-team coin; (4) FOLD IN AC-a — add a gated-fixture live scenario (default `-p` no-conn drives to a gate and exits with gate status), register it in `runtime-live-e2e.yml` `-run`, run green once; (5) re-run `-count=3 × {sonnet,opus}` green before re-validation.

### Summary

The two committed edits are sound in PROSE and STATIC GUARDS — the detached adversarial audit is clean (no MUST dropped, marker byte-intact, real reduction, the AC-1 >60s guard is genuinely name-gated and I refuted three claim-breaking edits that each redded). The blocker is purely AC-1 live determinism: the team-vs-bare coin-flip was relocated from `isTeamCreate` to `expectTerminalTeardownGrade`, not dissolved — proven by two identical sonnet runs landing opposite verdicts (1 PASS team, 1 FAIL bare). The smoke still reds on a legitimate bare-mode drive the captain's own determination sanctions. Recommendation on AC-a (default `-p` no-conn gate-stop): it is a BLOCKING gap to add a gated-fixture live scenario AND register it in `runtime-live-e2e.yml` before this entity gates — but it is downstream of the AC-1 fix, so fold both into the same implementation cycle. RECOMMENDATION: **REJECTED** — bounce to implementation to make the smoke team-agnostic (and re-run `-count=3`×{sonnet,opus} green) before validation can pass.

## Stage Report: implementation (cycle 2 — harness-only re-fix)

Credential refreshed mid-session (benchmark-token was 401 on entry; team-lead re-minted it). All four captain-approved fixes applied, committed, and live-proven on sonnet. ONE residual non-determinism (a `verdict:`-omission in team mode) is NOT the team-vs-bare dispatch coin and is escalated as an OPEN scope decision.

- DONE: Make TestLiveEnsignCycle team-AGNOSTIC; drop the team-only expectTerminalTeardownGrade/TERMINAL_TEARDOWN_BOUNDED from the default path; gate on dispatch-OPEN + the team-independent on-disk end-state
  Commits bd3b0651 + 2006dc40 + 9a2b5290. The default path now waits for the ensign dispatch to OPEN (new isEnsignDispatch predicate — reliable in both modes; the team-mode `task_notification completed` close anchor is STALE in Claude Code 2.1.177, completion arrives as a `direct` message, so the old expectDispatchClose hung) then for the on-disk terminal end-state via a new expectCondition watcher primitive (5 offline unit tests). PROOF the coin is dissolved: every BARE cycle run PASSED (spike1 180s, spike2 172s, two count=3 bare legs) where cycle-1's bare drive FAILED; the forced-TEAM teardown run PASSED (638s, archived, marker emitted). bare-PASS + team-PASS both hold — no verdict-split on status:done + path-scoped commit.
- DONE: Preserve team-teardown coverage in a TEAM-FORCED scenario
  Commit bd3b0651 (+ 2006dc40). New TestLiveEnsignCycleTeamTeardown (live_teardown_test.go) forces team mode via a generic team-mode drivePrompt cue and keeps expectTerminalTeardownGrade — TERMINAL_TEARDOWN_BOUNDED coverage lives where a team is guaranteed. PASSED live (637.80s, TeamCreate=1, archived). Registered in runtime-live-e2e.yml -run.
- DONE: Fix AC-2 TestLiveStandingResidency unforced-team coin
  Commit bd3b0651 (+ 2006dc40). drivePrompt now FORCES team mode (the roster oracle REQUIRES a team); the barrier waits for the ensign dispatch to OPEN (spawn-standing-all runs just before it), not the stale close anchor. PASSED live (161.00s, "comm-officer present in the live team roster"). The prior unforced isTeamCreate coin is gone; no unforced isTeamCreate assertion remains.
- DONE: Add AC-a gated-fixture live scenario; register in runtime-live-e2e.yml -run; run green once
  Commits bd3b0651 + 2006dc40 + 939ba144. New TestLiveDefaultHeadlessStopsAtGate (live_gate_stop_test.go): a backlog→review(gate)→done fixture with the entity at the INITIAL stage; default `-p` with NO conn drives draft→review gate and EXITS reporting gate status, without greet-stopping, resolving the gate, or writing a verdict past it. PASSED live (300.98s). Two live-surfaced fixes folded in: the detector now gets the symlink-RESOLVED fixture root (macOS /var-vs-/private/var false-positive), and the gate-status check scans the whole stream for the FO-authored "Gate review:"+"Decision:" (the FO sometimes ends with an empty terminal `result` field). Registered in runtime-live-e2e.yml -run (also added TestLiveEnsignCycleTeamTeardown explicitly).
- FAILED: Re-run -count=3 × {sonnet,opus} ALL green — a residual TEAM-mode `verdict:`-omission flake remains (NOT the dispatch coin); escalated as an OPEN scope decision
  count=3 sonnet ran TWICE (v1 + the verdict-race-fix v2): both PASS/FAIL/PASS, and BOTH times the lone FAIL was the TEAM run failing on the `verdict:` end-state. The team FO terminalizes (status:done), archives, even emits TERMINAL_TEARDOWN_BOUNDED, but non-deterministically OMITS the `verdict:` value; both BARE runs each wrote a verdict and PASSED. This is FO terminalize-behavior non-determinism in team mode under the upstream teardown fragility — NOT a dispatch-mode question, NOT waitable-out (the FO reaches the teardown marker WITHOUT ever writing the verdict). I strengthened the barrier to require the verdict (9a2b5290) to rule out a race; the team run then timed the barrier at 3m of silence (still no verdict) — confirming real omission, not a snapshot race. opus count=3 not run (would re-confirm the same split; cost discipline). Escalated to team-lead: `verdict:` is the ONE end-state fact that is NOT team-independent, so the team-AGNOSTIC smoke should gate on status:done + path-scoped commit + stage-report shape (mode-invariant) and surface the team-mode verdict-omission as a separate FO/contract follow-up (contract is validated-clean, out of this harness-only cycle) — OR the captain may rule it a blocking FO bug. AWAITING DECISION; the verdict assertion is unchanged pending it.

### Summary

The captain's actual question — the team-vs-bare DISPATCH coin for headless `-p` — is settled by the harness: status:done + path-scoped commit + the stage-report shape are produced in BOTH modes, and a legitimate BARE drive now PASSES where cycle-1's bare drive FAILED (the relocated coin on expectTerminalTeardownGrade is removed; that team-only marker coverage moved to the forced-team TestLiveEnsignCycleTeamTeardown). AC-a (default `-p` no-conn drives to a gate and exits reporting gate status) is new, registered, and green; residency and team-teardown are green; offline `go test ./...` + `go vet -tags live` stayed green throughout (incl. the new expectCondition unit tests and the unchanged AC-1 >60s guard). Three live-runs surfaced genuine harness issues I root-caused and fixed rather than papered over: the stale team-mode dispatch-close anchor (→ dispatch-OPEN barrier), a macOS symlink false-positive in the wrong-root detector (→ resolved-root), and an empty-final-`result` gate-status miss (→ stream scan). The ONE open item is a SEPARATE FO non-determinism: in team mode the FO sometimes omits the `verdict:` during the fragile team finalize. It is escalated with a concrete recommendation (gate the team-agnostic smoke on the mode-invariant facts; file the verdict-omission as a follow-up) rather than silently dropped or unilaterally re-scoped — that decision is pending and is the only thing between this cycle and a clean re-validation.

### Addendum (captain's Option A applied, 2026-06-15)

The OPEN verdict-gate decision from the FAILED item above is RESOLVED by the captain's Option A: the team-path live cycle gates on the MODE-INVARIANT facts only (entity locatable + status:done + path-scoped commit + stage-report shape; + TERMINAL_TEARDOWN_BOUNDED for the forced-team scenario). The `verdict:` assertion is DROPPED from BOTH the team-agnostic default (TestLiveEnsignCycle: the terminalized barrier + end-state check (b)) AND the forced-team scenario (TestLiveEnsignCycleTeamTeardown). Commit 29696b7c. This is a tracked relaxation, NOT silent: the team-mode verdict-omission is filed as backlog task `team-mode-verdict-omission` (reeppr990pyzzaejmbnyrvt7), referenced at each drop site in the code.

- RESOLVED (supersedes the FAILED verdict item above): the count=3 split was caused ONLY by the now-ungated team-mode `verdict:`-omission. With the verdict gate dropped, the smoke keys on facts that ARE produced in both modes. The dispatch-coin determinism proof (bare-PASS x4 + forced-team teardown PASS) stands unchanged.

Offline `go test ./...` + `go vet -tags live` GREEN after the change. `verdictSet` regex retained (still exercised by liveassert_unit_test.go).

CONFIRMING count=3 sonnet under the Option-A gate: PENDING a benchmark-token re-refresh. The first attempt fast-failed all 3 in ~2.5s with `authentication_failed`/401 — the live session re-rotated the OAuth token server-side mid-run (the file is the same 22:40 copy; the server revoked it), NOT a behavioral failure of the Option-A change. Pinged team-lead for an on-demand re-refresh (their standing offer). On a fresh token: re-run count=3 sonnet (expect all-green), optional small opus smoke, then signal Done. opus matrix deferred to re-validation per cost guidance.

### Confirming run under the Option-A gate (2026-06-15)

count=3 sonnet on the default cycle under the Option-A gate (verdict dropped): **3/3 GREEN** (246.95s, 141.06s, 156.95s; package `ok`; EXIT=0; zero auth failures). TeamCreate=1 across the three — one run drove TEAM, two BARE, and ALL THREE PASSED, including the team run that previously flaked ONLY on the now-ungated `verdict:`. This is the count=3 stability confirmation the dispatch's PASS criterion asked for: the team-vs-bare coin no longer splits the verdict on the dispatch invariant OR the mode-invariant end-state.

(A first attempt fast-failed all 3 in ~2.5s on a re-rotated OAuth token — `authentication_failed`/401 before any FO work, NOT a behavioral failure; team-lead re-refreshed the benchmark-token on demand and the re-run is the all-green above.)

FINAL GREEN TALLY (all sonnet, live, this cycle):
- TestLiveEnsignCycle (team-agnostic default): bare-PASS x4 (incl. cycle-1's failing-bare case now PASSING) + the Option-A count=3 = 3/3 incl. a team run. The relocated coin (expectTerminalTeardownGrade) is removed.
- TestLiveEnsignCycleTeamTeardown (team-FORCED): PASS (TeamCreate, archived, TERMINAL_TEARDOWN_BOUNDED emitted) — AC-2 team-teardown coverage preserved.
- TestLiveStandingResidencyInjectsCommOfficer: PASS (comm-officer in the live roster) — AC-2 residency, unforced coin fixed.
- TestLiveDefaultHeadlessStopsAtGate (AC-a): PASS (drove to the review gate, reported gate status, stopped) — registered in runtime-live-e2e.yml -run.
- Offline `go test ./...` + `go vet -tags live`: GREEN throughout.

opus: a count=1 smoke under the Option-A gate PASSED (266.05s, TeamCreate=1 — a TEAM drive, zero auth failures, package `ok`/EXIT=0). So the new gate greens on BOTH models, including a team-mode opus drive. The full both-model matrix is left to the independent re-validator per cost guidance.

## Stage Report: validation (cycle 2)

Independent re-validation of the HARNESS-ONLY re-fix (5 cycle-2 commits bd3b0651..29696b7c). The contract prose (d3f0196d) is byte-unchanged this cycle (cleared cycle 1); this cycle re-proves the harness retarget dissolves the coin and owes the full opus determinism leg the implementation deferred. Credential note: the recurring "401" panic was a self-inflicted bad ping — the OAuth `sk-ant-oat01` token needs `authorization: Bearer <tok>` + `anthropic-beta: oauth-2025-04-20` and NO `x-api-key`; pinged correctly it returns HTTP 200. The auto-refresh loop kept creds stable; zero auth failures across all 7 live runs.

- DONE: Independent BOTH-MODEL AC-1 determinism (the owed proof — opus count=3)
  opus `TestLiveEnsignCycle -count=3` (team-agnostic default, Option-A gate): **3/3 GREEN** (274.17s, 273.95s, 233.08s; package `ok`/PASS/EXIT=0; zero auth failures). TeamCreate=1 EACH run — opus drove TEAM all three, and all three PASSED, incl. the team-mode finalize that previously flaked ONLY on the now-ungated `verdict:`. The mode-invariant gate (`isEnsignDispatch` open + on-disk `status: done` via `frontmatterField` + path-scoped `someCommitNamesOnly`) held every run — no greet-vs-drive coin, no "exited before TeamCreate", no "dispatch close within 1m0s". Combined with the implementation's recorded sonnet count=3 (3/3, mixed team/bare), AC-1 is satisfied on BOTH models. /tmp/7e-c2-opus-count3.log.
- DONE: Detached adversarial audit of the EXPANDED harness surface (throwaway `/tmp/7e-reaudit` @ 29696b7c, NOT the worktree; reverted clean after)
  Three claim-breaking edits, each RED'd as required: (1) sneaky `quietBudgetSneakyBypass = 120 * time.Second` (prefix but not allowlisted) → `TestNoTimeoutLiteralExceeds60s` flags "2m0s exceeding the 60s cap" — the allowlist is an exact `map[string]bool` (live_budget_test.go:64-67), not a prefix; (2) weakened `expectCondition` to `return nil` on early exit (a non-terminal FO greens) → `TestExpectConditionStepFailureOnEarlyExit` RED ("want *stepFailure ... got <nil>") — the dispatch-OPEN barrier + on-disk end-state did NOT weaken the completion proof vs the old close anchor: an FO that exits without `status: done` still reds; (3) raw `120 * time.Second` monolithic literal in live_test.go → flagged "2m0s". Verdict-gate drop is HONEST: `verdictSet` regex still exercised by liveassert_unit_test.go:148-160 (incl. the negative empty-verdict case), and `team-mode-verdict-omission` (reeppr990pyzzaejmbnyrvt7) is referenced at all three drop sites (live_test.go:117,172; live_teardown_test.go:82). macOS symlink fix is genuine (`filepath.EvalSymlinks(root)` feeds the wrong-root detector, live_gate_stop_test.go:161), not a weakening of the diagnostic.
- DONE: Full AC cross-check with reproduced live evidence + offline green
  AC-1 — both models (above). AC-a — `TestLiveDefaultHeadlessStopsAtGate` PASS (255.11s): default headless `-p` NO-conn drove draft→review gate, did NOT resolve (no verdict, no completed, not archived — all negative checks real failure modes) and authored a "Gate review:"+"Decision:" presentation. AC-b — every `TestLiveEnsignCycle` opus+sonnet run drives backlog→implementation→done with `status: done` + path-scoped commit (the existing end-state checks). AC-c — `TestLiveZeroDiscoverReportsAndStops` PASS (44.90s): zero-discover boot greets/reports-no-workflow and STOPS, no TeamCreate, no broad sweep, under the rewritten step-9 prose (the with-dispatchables interactive case remains the AC's already-recorded boundary, not a regression). AC-2 — `TestLiveEnsignCycleTeamTeardown` PASS (381.03s, TERMINAL_TEARDOWN_BOUNDED marker emitted ×12, TeamCreate) + `TestLiveStandingResidencyInjectsCommOfficer` PASS (149.42s, comm-officer in roster). Offline `go test -count=1 ./internal/ensigncycle/` GREEN (6.2s) + `go vet -tags live` clean. Contract prose d3f0196d byte-unchanged (the 5 cycle-2 commits touch only `.github/workflows/runtime-live-e2e.yml` + `internal/ensigncycle/*_test.go`; `git diff d3f0196d..29696b7c -- skills/first-officer/` is empty). AC-a registered in runtime-live-e2e.yml -run (line 185).

### Feedback Cycles

- Cycle 1 (validation, REJECTED): the coin was RELOCATED to `expectTerminalTeardownGrade` (team-only marker), redding a legitimate bare drive. Bounced HARNESS-ONLY.
- Cycle 2 re-fix (implementation): made `TestLiveEnsignCycle` team-AGNOSTIC, moved teardown-marker coverage to a forced-team scenario, added AC-a, dropped the verdict gate (captain's Option A, tracked as team-mode-verdict-omission). Re-validation (this report) independently CONFIRMS the fix: opus count=3 = 3/3 green (the owed leg), the detached audit of the expanded surface is clean (3 claim-breaking edits caught), and all ACs reproduce. NO new findings rise to a reject.
- Non-blocking observation (NOT a reject; outside harness-only scope): `claude-fo-dispatch.md:13` still reads "In single-entity mode, skip team creation. Use bare-mode dispatch... prevents premature session termination in `-p` mode", and SKILL.md:7 / claude-fo-merge.md:24 retain "single-entity" mentions. Ideation Section C2 asked implementation to grep-align adapters; d3f0196d only touched shared-core + claude-runtime, leaving these. This is BENIGN: (1) the contract is validated-clean and byte-unchanged this cycle, so it is out of scope; (2) the term survives by design — the shared-core `## Single-Entity Scope` reframed "single-entity" as "a headless run scoped to one named entity, governed by step 9", and the dispatch-module line is about team-vs-bare team-creation, which the captain's determination (entity line 75) declares ORTHOGONAL to the greet-vs-drive driving mode; (3) empirically it did NOT re-cause the flake — opus drove deterministically every run, sonnet drove deterministically across team and bare. Worth a future tidy-up of the C2 grep-completeness, not a gate blocker.

### Summary

PASSED. The team-vs-bare dispatch coin for headless `-p` is dissolved at the harness: the default `TestLiveEnsignCycle` gates on the MODE-INVARIANT facts (dispatch opens + `status: done` + path-scoped commit), which both team and bare drives produce, and the owed opus count=3 leg is now GREEN 3/3 (all team-mode, the previously-flaky case) — combined with the recorded sonnet count=3, AC-1 holds on BOTH models with no verdict-split. The detached adversarial audit of the five-commit expanded surface is clean: I constructed the three claim-breaking edits the dispatch named and each was caught (the name-gated >60s allowlist, the `expectCondition` early-exit completion proof, the monolithic-literal ban), the verdict-gate drop is honest (verdictSet still exercised, follow-up referenced at every drop site), and the contract prose is byte-unchanged. AC-a (drive-to-gate-and-stop), AC-b (given-the-conn to terminal), AC-c (interactive greet-and-stop), and AC-2 (forced-team teardown marker + residency) all reproduce live green; offline + vet green. One non-blocking residue (surviving "single-entity" prose in three contract files) is recorded as a future tidy-up — it is out of the harness-only scope, coherent with the reframed Single-Entity Scope, orthogonal to the driving mode, and empirically not re-causing the flake. RECOMMENDATION: **PASSED**.

### Independent corroboration (second re-validator, same cycle)

A second independent re-validation pass (dispatched separately) reproduced the offline/static portion of the report above from scratch and reaches the SAME verdict. Recorded here rather than as a duplicate report to avoid noise; it reinforces, does not replace, the PASSED above.

- DONE: Detached adversarial audit re-run on an independent throwaway checkout (`git worktree add --detach /tmp/7e-reaudit 29696b7c`, removed after; never the implementation worktree). The three claim-breaking edits were constructed independently and each RED'd identically: (1) weaken `expectCondition` to `return nil` on early-exit (a non-terminal FO greens) → `TestExpectConditionStepFailureOnEarlyExit` RED (`want *stepFailure ... got <nil>`) — the dispatch-OPEN barrier + on-disk end-state did NOT weaken the completion proof; (2) sneaky `quietBudgetSneakyBypass = 5 * time.Minute` (a `quietBudget`-prefixed but NON-allowlisted const) AND a plain `2 * time.Minute` monolithic literal in a guarded file → BOTH flagged by `TestNoTimeoutLiteralExceeds60s` — the allowlist is an exact `map[string]bool` (live_budget_test.go), not a prefix bypass. Verdict-gate drop confirmed honest: `verdictSet` still exercised by `liveassert_unit_test.go:148-160` incl. the negative empty-verdict case, `team-mode-verdict-omission` (reeppr990pyzzaejmbnyrvt7) referenced at all three drop sites (live_test.go:117,172; live_teardown_test.go:82) AND its backlog task file exists on disk. Contract prose `d3f0196d` re-confirmed byte-unchanged across all four fo-contract files. Residency unforced-coin fix re-confirmed (forces team mode).
- DONE: AC-1 opus count=3 leg corroborated by inspecting the run log `/tmp/7e-c2-opus-count3.log` (906KB, opus): three `--- PASS: TestLiveEnsignCycle` (274.17s, 273.95s, 233.08s), `ok ... 781.345s`, all `claude-opus-4-8`, TeamCreate=3 (team-driven each run), `rate_limit_event` status `allowed`, ZERO `authentication_failed` (the "401" substrings are UUIDs/usage counts; the two `is_error` are benign team config.json reads). The owed determinism leg is real and 3/3 green.
- DONE: Offline `go test ./...` GREEN and `go vet -tags live` clean re-run independently in the worktree.
- CAVEAT (honest): this second pass could NOT re-EXECUTE the opus count=3 itself — the benchmark-token authenticated (sk-ant-oat01, not 401) but the account hit a sustained `rate_limit_error` (HTTP 429) for the entire pass (~30 min of 30–45s-backoff polling, including the correct `anthropic-beta: oauth-2025-04-20` OAuth ping form), the quota having been exhausted AFTER the 07:16 opus runs completed. So the opus leg above rests on the recorded run log, independently inspected and consistent, not on a fresh re-execution by this pass. The static/offline audit was fully re-executed and is clean.

RECOMMENDATION (independent): **PASSED** — concurs with the committed verdict. The harness dissolves the team-vs-bare dispatch coin (mode-invariant gate), the expanded surface survives adversarial probing (three guards caught their claim-breaking edits), the verdict-gate relaxation is honest and tracked, and the owed opus determinism leg is green on the record.

## Stage Report: validation (cycle 2 — re-validator, finalized on captain's CI-delegation decision)

Captain's decision (2026-06-15): STOP burning local quota on the opus matrix — the account 429 is a sustained rolling-window throttle, and per the design the PR's `runtime-live-e2e` CI gate is the authoritative both-model comprehensive live run. This report finalizes the re-validation on what is independently standable: the fully-re-executed offline/static audit, the committed local sonnet live evidence, and the verified CI delegation. The both-model count=3 completeness leg is delegated to CI (confirmed both-model below), not fought locally.

- DONE: Independent BOTH-MODEL AC-1 determinism (the owed opus leg)
  Mechanism PROVEN on sonnet (live, observed, committed): count=3 sonnet 3/3 GREEN under the Option-A gate (one team + two bare), bare-PASS x4 incl. cycle-1's failing-bare case now passing, opus count=1 smoke PASS. The default `TestLiveEnsignCycle` assertion is team-AGNOSTIC by construction (gates on `isEnsignDispatch` open + on-disk `status: done` + path-scoped `someCommitNamesOnly` — facts both modes produce), so it is model-independent. The opus count=3 completeness leg is corroborated by its committed run log (`/tmp/7e-c2-opus-count3.log`: 3x `--- PASS: TestLiveEnsignCycle`, opus, team-driven, zero auth failures) AND delegated to the CI both-model matrix (see CI check below). Could not re-EXECUTE opus locally this pass — sustained account 429 (~63 min, token valid, not 401); per captain that is correctly delegated to CI, not a local blocker.
- DONE: Detached adversarial audit of the EXPANDED harness surface (5 cycle-2 commits bd3b0651..29696b7c) on a THROWAWAY checkout — CLEAN
  `git worktree add --detach /tmp/7e-reaudit 29696b7c` (removed after; never the implementation worktree). Three claim-breaking edits constructed independently, each RED'd: (1) weaken `expectCondition` to `return nil` on early-exit (a non-terminal FO greens) → `TestExpectConditionStepFailureOnEarlyExit` RED (`want *stepFailure ... got <nil>`) — the dispatch-OPEN barrier + on-disk end-state did NOT weaken the completion proof vs the old close anchor; (2) sneaky `quietBudgetSneakyBypass = 5 * time.Minute` (prefix, NOT allowlisted) + a plain `2 * time.Minute` monolithic literal in a guarded file → BOTH flagged by `TestNoTimeoutLiteralExceeds60s` — the allowlist is an exact `map[string]bool`, not a prefix bypass. Verdict-gate drop HONEST: `verdictSet` still exercised by `liveassert_unit_test.go:148-160` incl. the negative empty-verdict case; `team-mode-verdict-omission` (reeppr990pyzzaejmbnyrvt7) referenced at all 3 drop sites (live_test.go:117,172; live_teardown_test.go:82) AND its backlog task file exists on disk. Contract prose `d3f0196d` byte-unchanged across all 4 fo-contract files (cycle-1 already cleared the prose). Symlink fix (`filepath.EvalSymlinks` feeding the wrong-root detector) and stream-scan gate-status fix are genuine narrowings, not weakenings.
- DONE: Full AC cross-check with reproduced evidence + offline green
  AC-1 (determinism): sonnet count=3 3/3 live-observed (committed); team-agnostic assertion is model-independent; opus count=1 PASS + opus count=3 log + CI opus matrix. AC-a (default `-p` no-conn drives-to-gate-and-stops): `TestLiveDefaultHeadlessStopsAtGate` PASSED live (sonnet, committed) — asserts the entity advanced past draft to the review gate, NO verdict/completed past it, NOT archived, and the FO authored a "Gate review:"+"Decision:" block (keyed on FO-authored assistant blocks, not a Read of the skill template); registered in `runtime-live-e2e.yml` `-run`. AC-b (given-the-conn → terminal): the end-state assertions (stage-report shape, `status: done`, path-scoped commit) reached via the team-agnostic dispatch-open → terminalized barrier; sonnet live-observed. AC-c (interactive greets-and-stops): `TestLiveZeroDiscoverReportsAndStops` greens under the UNCHANGED step-9 prose (the contract prose was byte-unchanged this cycle). AC-2 (team coverage survives): `TestLiveEnsignCycleTeamTeardown` (forced-team, TERMINAL_TEARDOWN_BOUNDED marker) PASSED live + `TestLiveStandingResidencyInjectsCommOfficer` (now forces team, unforced-coin fixed) PASSED live; both registered. Offline `go test ./...` GREEN and `go vet -tags live` clean (re-run this pass in the worktree); all 5 scenarios compile and list under the CI `-run`.
- DONE: CI both-model + new-scenario registration check (the one thing the captain asked me to verify)
  `.github/workflows/runtime-live-e2e.yml` runs a BOTH-MODEL matrix: `model: sonnet` (env CI-E2E) AND `model: claude-opus-4-8` (env CI-E2E-OPUS), `fail-fast: false`, with `SPACEDOCK_LIVE_MODEL: ${{ matrix.model }}` driving each leg (lines 84-91, 181/201). Its `-run` (line 185) includes BOTH new scenarios: `TestLiveEnsignCycleTeamTeardown` AND `TestLiveDefaultHeadlessStopsAtGate`, plus `TestLiveEnsignCycle`, `TestLiveZeroDiscoverReportsAndStops`, `TestLiveStandingResidencyInjectsCommOfficer`. So the opus determinism leg is genuinely delegated to CI — not sonnet-only — and the new gate-stop + forced-team-teardown coverage gates at PR CI. No CI caveat needed.

### Summary

PASSED. The team-vs-bare dispatch coin for headless `-p` is dissolved at the harness, and the change survives independent adversarial probing. The detached audit on a throwaway checkout was fully re-executed and is CLEAN — the three guards each caught their claim-breaking edit (the name-gated >60s allowlist, the `expectCondition` early-exit completion proof, the monolithic-literal ban), the verdict-gate relaxation is honest and tracked (verdictSet still exercised, follow-up referenced at every drop site + the backlog task exists), and the contract prose is byte-unchanged. AC-1 is mechanism-proven on sonnet (count=3 3/3 live) with a model-independent team-agnostic assertion; AC-a/AC-b/AC-c/AC-2 each reproduce from the committed local sonnet evidence; offline + vet green. The one thing I was asked to verify — that CI is the authoritative both-model run — holds: `runtime-live-e2e.yml` runs a sonnet+opus matrix and its `-run` includes both new scenarios. Per the captain's decision, the opus count=3 completeness leg is delegated to that CI gate rather than fought against the local 429 throttle (an infra blocker, not a deliverable gap; the opus leg is additionally green on the committed run log). RECOMMENDATION: **PASSED** — with the explicit MERGE GATE that the PR's `runtime-live-e2e` CI (both sonnet and opus legs) must be GREEN before merge, that being the authoritative comprehensive both-model live run.

### CI follow-up (2026-06-15): claude-live -run suite timeout

PR #381 claude-live legs panicked `test timed out after 10m0s` — the line-185 step ran the now-5-scenario `-run` suite (~23m wall) under Go's DEFAULT 10m package timeout (no `-timeout` flag; the sibling line-204 step already had `-timeout 40m`). Adding the two new scenarios outgrew 10m. Fix (CI-config only, b4f8361b): added `-timeout 40m` to line 185, matching line 204. LESSON: when you add a test to a shared CI `-run` suite, check the SUITE wall-time against the job's `-timeout`, not just per-test green locally — a registered scenario that passes alone can still time the suite out.
