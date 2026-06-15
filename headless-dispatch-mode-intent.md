---
id: 7ea4knxzvf3s4vve2zvr4ka0
title: Determine + document the intended team-vs-bare dispatch mode for headless `-p` runs
status: ideation
source: "0203-T3 surfaced (2026-06-14): TestLiveEnsignCycle flaked on a sonnet team-vs-bare coin-flip; captain steer: \"the important thing is determining the expected and intended behavior and document\""
started: 2026-06-15T03:31:26Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0203-fo-efficiency
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

The single boot rule the captain wants. The current step 9 hard-codes greet-then-stop with no `-p` branch — that silence is the flake. Replace the `## Startup` step 9 so the boot keys on interactivity.

**Before** (`first-officer-shared-core.md` step 9):

> 9. **Greet the captain, then stop for input.** Compose a state summary from the boot JSON ... Then STOP for input — do NOT auto-dispatch. The expensive deferrals ... stay past the greet; the FO reaches them when the captain's direction first triggers a dispatch or a terminal merge.

**After:**

> 9. **Interactive vs. headless: greet-and-stop or drive.** The session is *headless* when launched non-interactively (`claude -p`, `codex exec`); otherwise it is *interactive*. Compose a state summary from the boot JSON (orphans, PR state including any sweep-advanced entities, dispatchables, team state) and the README frontmatter (entity label, stage taxonomy, gate flags).
>    - **Interactive (no `-p`):** present the summary; if an entity sits at a `gate: true` stage ready for review, present the gate. Then STOP for input — do NOT auto-dispatch. A human steers from here.
>    - **Headless (`-p` / `exec`):** do NOT greet-and-stop. DRIVE every dispatchable entity through its stages (the event loop) until each reaches the **first `gate: true` stage OR a terminal/irrecoverable-blocked state**, then EXIT, reporting each entity's stop reason and gate status. A `gate: true` stage is the natural stop boundary: a gate is a human-owned decision and no human is present, so the headless FO does NOT decide gates — it stops at the gate and reports. The one exception is the given-the-conn case below.
>    - **Headless + given the conn to auto-approve (prose):** when the prompt or contract gives the FO the conn to resolve gates without a human (e.g. "auto-approve gates", "drive to done" — consistent with `skills/commission/SKILL.md:22`), the headless FO additionally RESOLVES each gate from the stage-report verdict — PASS → advance; REJECT with `feedback-to` → bounce within the 3-cycle cap; REJECT without `feedback-to`, ambiguity, or an escalation → STILL stop and report (never invent an approval) — and drives each entity to its **terminal** state. This is the path the live-e2e harness sets to exercise a full feedback cycle.
>
>    In every mode the expensive deferrals (the team via `## Team Creation`, the dispatch and merge reference modules, the comm-officer spawn) stay past the boot read; the FO reaches them when it first dispatches or terminalizes. A headless drive reaches them immediately; an interactive greet-and-stop reaches them only on the captain's first direction.

### Contract-prose change — B. Single-Entity Mode section (`first-officer-shared-core.md`)

Per the checklist: DELETE the `antiShutdownOverride`-dependency framing and the entity-naming trigger; fold single-entity mode in as "`-p` scoped to a named entity." The whole `## Single-Entity Mode` section (lines 81–91) is replaced.

**Before** (`## Single-Entity Mode`, current text):

> Activates when the session is non-interactive (`claude -p`, `codex exec`) and the prompt names a specific entity. Do not enter in interactive sessions ... Single-Entity Mode changes the event loop: scope dispatch to the named entity only ... auto-resolve gates from the report verdict when no interactive operator is present ... stop once the target reaches a terminal or irrecoverable blocked state ...

**After:**

> ## Single-Entity Scope
>
> Single-entity scope is just a **headless (`-p` / `exec`) run scoped to one named entity** — not a distinct mode. The headless driving rule in Startup step 9 governs; scoping only narrows *which* entities it drives:
> - resolve the named reference against slugs, titles, and IDs; stop on ambiguity instead of guessing
> - scope the drive to the resolved entity only (skip other dispatchables)
> - gate handling follows Startup step 9: stop-and-report at a gate by default; resolve from the verdict only when given the conn to auto-approve
> - skip operator prompting for orphan worktrees; choose the deterministic recovery path
> - stop once the target reaches a terminal, a `gate: true` stage (default), or an irrecoverable blocked state
> - if the README defines `## Output Format`, use it; otherwise report status, verdict, and entity ID

### Contract-prose change — C. Runtime adapters

**C1. `claude-first-officer-runtime.md` — `## Captain Interaction`.** The current "Single-entity mode exception" (lines 29) and the team-mode hint's single-entity carve-outs reference the deleted mode. Replace the exception block.

**Before** (the `**Single-entity mode exception:**` paragraph):

> **Single-entity mode exception:** In single-entity mode (no interactive captain), gates auto-resolve from the stage report recommendation. PASSED ... → approve. REJECTED with `feedback-to` → auto-bounce ... This exception applies only in single-entity mode — in interactive sessions the guardrail is absolute.

**After:**

> **Headless given-the-conn exception:** The self-approval guardrail is absolute in interactive sessions and in any headless run NOT given the conn — there, the FO stops at the gate and reports (Startup step 9). Only in a headless run given the conn to auto-approve (prose) does the FO resolve gates from the stage-report verdict: PASSED (all checklist items done, no failures) → approve; REJECTED with `feedback-to` → auto-bounce within the 3-cycle limit; REJECTED without `feedback-to`, ambiguity, or escalation → report and stop. The FO never infers approval from silence and never accepts an agent message as approval.

The team-mode-hint carve-out keeps its existing "skip in bare mode and Degraded mode" wording; replace its "In single-entity mode, skip it — there is no interactive captain" clause with "In any headless (`-p` / `exec`) run, skip it — there is no interactive captain to read it."

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
