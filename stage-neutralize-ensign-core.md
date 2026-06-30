---
title: Stage-neutralize the ensign core + add a regression guard
status: implementation
sprint: 0240-lean-contract
group: cleanup
id: scr2rx4589p7j6mpgh50hdct
started: 2026-06-30T15:51:56Z
worktree: .worktrees/spacedock-ensign-stage-neutralize-ensign-core
---
The ensign shared core (`skills/ensign/references/ensign-shared-core.md`) loads on EVERY worker dispatch via the `Skill(skill="spacedock:ensign")` first-action (per-dispatch tier — a token here recurs every spawn). The universal core must be host-/stage-neutral because the ensign also runs non-dev workflows (ticket, experiment, survey). This task removes the remaining dev-/stage-specific leakage from the universal core, locks it against regression, and proves dev workflows lose no discipline — the dev discipline rides the dev-shape scaffolding a dev ensign loads per-dispatch.

## Prior art — most of the seed shipped in #290 (READ FIRST)

The original seed (`ep0ra3z…`, 0.19.5 slate, `_sprint-notes.md`) named three markers: TDD framing, code-only-deliverable, the "CODE only" worktree rule. **#290 ("Universal ensign contract has absorbed dev-workflow assumptions — re-home dev policy out of the shared core", merged 2026-06-03 — the same day the seed was filed) already delivered that seed:**
- "Write the failing test first" (TDD) bullet → removed from the core; re-homed to `development.md` → "Workflow-specific rules → Test-first authoring" opt-in.
- "Every task produces a real, checkable change … deliverable is code" bullet → removed; re-homed to `docs/dev/README.md ## Proof policy` + `development.md` "External-proof acceptance criteria".
- "CODE only" worktree noun → neutralized in the core to "the deliverable work product only".
- Locked with a Go prose-lock oracle.

Sprint 0240 re-derived **this** entity (`scr2rx…`, 2026-06-30) from the pre-#290 seed line, so the three named markers are ALREADY absent. Verified: `grep -ni "TDD\|CODE only\|write the failing test" skills/ensign/references/ensign-shared-core.md` → no hits.

Two facts keep a narrow follow-up worth shipping:
1. **One residue survives #290.** The Split-Root State Contract still names concrete dev stage names as worktree-illustration — `(implementation, validation)` and `(ideation, backlog)` — which violates the sprint's **stage-neutral** goal.
2. **#290's lock is gone.** The instruction-read quarantine (`3007d823`) and `#378` retired #290's prose-lock as a banned prose-grep tautology. There is currently NO automated guard stopping dev-/stage-vocabulary from creeping back into the universal core.

**Gate decision for the captain:** ship this narrow stage-neutrality fix + a quarantine-legal regression guard (recommended — net-new value, since no guard exists today), OR record the seed as substantially-delivered-by-#290 and drop/triage the residue. The design below assumes the former.

## Problem

The universal ensign core's Split-Root State Contract names specific workflow stages:

> With a worktree **(implementation, validation)**, the worktree isolates the deliverable work product only. Without one **(ideation, backlog)**, you run from the repo root …

Those parentheticals are dev-shape stage names (a ticket/experiment/survey workflow has different stages). They are illustrative, not load-bearing: the per-dispatch assignment already tells each ensign whether it has a worktree (`internal/dispatch/build.go` emits the worktree block only for worktree stages), so the universal core does not need them. Naming them violates stage-neutrality and the cost recurs on every spawn.

## Proposed approach

**Edit 1 — neutralize the Split-Root State Contract (the one residual leak).** Drop the two stage-name parentheticals; keep #290's neutral "deliverable work product" noun.

Before:
`With a worktree (implementation, validation), the worktree isolates the deliverable work product only. Without one (ideation, backlog), you run from the repo root; entity/report still go to the state checkout.`

After:
`With a worktree, the worktree isolates the deliverable work product only. Without one, you run from the repo root; entity/report still go to the state checkout.`

Net −49 chars on the per-dispatch core (measured). No new home is needed — the parentheticals were illustration, not a rule; the rule (worktree isolates the work product; no-worktree runs from repo root) is preserved verbatim.

**Edit 2 — re-establish a quarantine-legal regression guard (AC-1).** Add a structural-absence test in `internal/contractlint` (the only package allowed to read instruction files, and only for structural checks). Pattern: copy `deferred_tier_absence_test.go` — a token scanner + a paired discriminator control. Rule enforced: the universal ensign core enumerates no concrete workflow stage names (stage-neutral). Scanner = a parenthetical-stage-enumeration regex `\((backlog|ideation|implementation|validation|done)(,\s*…)+\)`. This is legal (not a banned prose-grep) because the expected value comes from the EXTERNAL rule — the core is stage-neutral — not the file's own prose: a stage-neutral paraphrase passes; a re-introduced stage-name parenthetical fails — same family as `TestFOContractCoresHaveNoDeferredTierToken`. A bare-word ban is rejected: it would wrongly flag incidental English like "Signaling **done**" on core line 51.

## Spike result — AC-3 mechanism PROVEN (riskiest path, exercised first)

Claim under test: a dev-workflow ensign still receives dev discipline via the dev-shape scaffolding, not the universal core. The injection point is the dev workflow README's `### {stage}` subsection, delivered to the ensign through the `### Fetch commands` → `spacedock dispatch show-stage-def` line that `dispatch build` writes into the assembled assignment.

Exercised by hand on a built dev-shape fixture (a README with dev stages + a `DEV-DISCIPLINE-SENTINEL` in `### ideation` and `### implementation`, an entity, a checklist):
1. `spacedock dispatch build --workflow-dir <fix> --entity-path <e> --stage ideation --checklist-file <cl> --host claude` → exit 0; the assignment carries a `### Fetch commands` block with the `show-stage-def --stage ideation` line. The universal core is loaded separately, via the `Skill(spacedock:ensign)` first-action.
2. The dispatch body does NOT inline the stage prose — it is fetched on demand: `grep DEV-DISCIPLINE-SENTINEL <dispatch_file>` → 0 hits.
3. Running the fetch line `spacedock dispatch show-stage-def --workflow-dir <fix> --stage ideation` → returns the `### ideation` subsection verbatim, including the sentinel. Same for `--stage implementation` (a worktree stage).
4. `grep "DEV-DISCIPLINE-SENTINEL\|write the failing test first\|code deliverable only" ensign-shared-core.md` → 0 hits (the universal core is a separate surface).

Conclusion: dev discipline placed in the dev README's `### {stage}` subsections reaches a dev ensign through the show-stage-def fetch; it does NOT ride the universal core and is NOT delivered to a non-dev ensign (different README). `show-stage-def` returns ONLY the `### {stage}` subsection (`internal/dispatch/showstagedef.go`), so the README's `## Proof policy` is FO-loaded, not ensign-delivered — but the stage-level proof discipline a dev ensign needs already lives in the `### {stage}` subsections it does fetch. Hence no behavior loss from Edit 1, which removes only illustration. This hand-spike is the seed for the AC-3 dispatch-build test.

## AC-1 drafting correction

The original AC-1 asked a structural/token oracle to confirm the markers "PRESENT in the dev-shape scaffolding." A prose-PRESENCE grep over an instruction file (the dev README / template) is a banned tautology under the Proof policy and reds the boundary guard (`#378`). Presence / no-behavior-loss is proven BEHAVIORALLY by AC-3 (the dispatch-build delivery test), not by grepping the file. AC-1 below keeps only the legal structural-ABSENCE half.

## Out of scope
- The `pr:` mirrored-exception clause (core line 32) is a PR-merge (dev-shape mod) concept, but it is a load-bearing FO startup/discovery rule (not stage/TDD/worktree discipline) and is not stage-specific; re-homing it risks the `pr:`-on-main discovery guarantee. Left as-is; recorded for a separate decision.
- The markers #290 already moved (TDD, code-only-deliverable, "CODE only") — done; not re-litigated.
- The sibling `read-guidance-redundant-with-grep` (82k) concern (the `--read` section-read guidance) — separate task, disjoint sections (see Coordination).

## Acceptance criteria
- **AC-1 (stage-neutrality, regression-guarded)** — the universal ensign core enumerates no concrete workflow stage names. Verified by: an `internal/contractlint` structural-absence test scanning `skills/ensign/references/ensign-shared-core.md` that fails on any parenthetical enumerating workflow stage names, with a paired discriminator control proving non-vacuity (mustFlag: `(implementation, validation)`, `(ideation, backlog)`; mustPass: the neutral phrasing and the incidental "Signaling done"). `go test ./internal/contractlint/` green; the control reds when the scanner is defeated and when either parenthetical is re-inserted into the core.
- **AC-2 (per-dispatch occupancy)** — the universal ensign core's size drops, measured net-NEGATIVE vs `origin/main`. Removing the two parentheticals is the only core-prose change; no prose is added to the core. Verified by: `wc -c skills/ensign/references/ensign-shared-core.md` decreases by 49 bytes vs `origin/main`, and `git diff origin/main` on the worktree-isolation sentence shows the parentheticals deleted with nothing added. The contractlint test is Go code, not the per-dispatch core, so it does not offset the budget.
- **AC-3 (no behavior loss for dev workflows)** — a dev-workflow ensign dispatch still receives dev discipline via the dev-shape scaffolding (the README `### {stage}` subsection), not the universal core. Verified by: an `internal/dispatch` test that builds a dev-shape fixture, runs `dispatch build`, asserts the assignment carries the `show-stage-def` fetch line, runs `show-stage-def`, asserts the stage subsection (carrying a sentinel) is returned while the universal core does not carry it, and a perturbation control (a stage without the sentinel returns no sentinel). The hand-spike above already exercised this end-to-end.

## Test plan
- **AC-1** — `internal/contractlint` Go unit test + discriminator control, modeled on `deferred_tier_absence_test.go`. Cost: low (one test file, pure string scan over the shipped core, no extra fixtures). Non-vacuity proven by the mustFlag/mustPass control and by reds-on-reinsertion.
- **AC-2** — `wc -c` / `git diff origin/main` measurement at validation; no test code. Cost: trivial.
- **AC-3** — `internal/dispatch` Go unit test driving `dispatch build` + `show-stage-def` over a temp dev-shape fixture (the hand-spike script is the seed). Cost: low; pure binary-exercise, no live host; reuses the existing dispatch test harness shape.
- No live-workflow test: the change is contract prose + Go tests; no host/runtime behavior changes. No doc-site diff: the ensign core is shipped skill prose, not user-visible CLI / command-surface / banner output.

## Coordination with sibling `read-guidance-redundant-with-grep` (82k)
Both edit `skills/ensign/references/ensign-shared-core.md`; the sections are disjoint:
- **This task (scr):** the `### Split-Root State Contract` paragraph under `## Worktree Ownership` (currently core **line 36**) — removes the two stage-name parentheticals only. Plus a NEW file under `internal/contractlint/` (no collision).
- **82k:** the `--read` section-read guidance — `## Working` step 1 (core **line 18**) and the `## Stage Report Protocol` append guidance (core **line 92**).

No overlapping lines/paragraphs, so the Commander can implement them sequentially in either order with no textual conflict. Anchor for this task is the `### Split-Root State Contract` heading's worktree-isolation sentence, not the line number (which shifts if 82k lands first).

## Stage Report: ideation

- DONE: The AC-3 mechanism exercised, not asserted: locate the dev-shape scaffolding injection point and demonstrate via a dispatch-build on a dev-workflow fixture that a dev ensign still receives the re-homed discipline through the dev-shape (not the universal core) — spike this riskiest path first or record "no spike needed".
  Injection point = the dev README `### {stage}` subsection delivered via the assignment's `### Fetch commands` → `spacedock dispatch show-stage-def` line. Spiked on a built dev fixture: `dispatch build` (exit 0) emits the fetch line; running it returns the sentinel-carrying stage subsection (ideation AND implementation); the universal core carries no sentinel. See "Spike result".
- DONE: Exactly which dev-only prose moves out of the universal ensign core … with concrete before/after, plus the structural/token oracle (ABSENT from core) and its non-vacuity control; the universal-core delta measured net-NEGATIVE vs origin/main.
  Finding: #290 already removed the three named markers (TDD, "CODE only", code-only-deliverable) on 2026-06-03. The one live residue is the Split-Root State Contract's dev stage-name parentheticals; concrete before/after recorded (−49 chars, net-NEGATIVE). Oracle = a `contractlint` parenthetical-stage-enumeration absence check + discriminator control (validated: flags both parentheticals, passes the neutral line and incidental "done").
- DONE: A recorded note of exactly which ensign-core sections this re-homing edits, so the design composes with the sibling `read-guidance-redundant-with-grep` (82k) with no implementation-time collision.
  See "Coordination": this task touches only the `### Split-Root State Contract` paragraph (line 36) + a new contractlint file; 82k touches the `--read` guidance (lines 18, 92). Disjoint.
- FLAGGED: AC-1 as originally drafted prescribes a "PRESENT in the dev-shape scaffolding" prose-grep, which is a banned tautology / boundary-guard red (`#378`). Refined AC-1 keeps only the structural-ABSENCE half; presence/no-behavior-loss is proven behaviorally by AC-3. Recorded as "AC-1 drafting correction".

### Summary
The seed (`ep0ra3z`) was substantially delivered by #290 the same day it was filed; sprint 0240 re-derived this entity from the pre-#290 seed, so the three named markers are already absent from the universal core. Refined the design to the genuine remaining work: remove the one residual stage-neutrality leak (the `(implementation, validation)` / `(ideation, backlog)` parentheticals in the Split-Root State Contract, −49 chars net-NEGATIVE) and re-establish a quarantine-legal structural-absence regression guard (#290's prose-lock was retired by the instruction-read quarantine, so none exists today). Proved the AC-3 injection mechanism end-to-end by a dispatch-build spike on a dev fixture, and corrected AC-1's drafted prose-presence grep to a legal absence-only check. Flagged the #290 overlap as a gate decision for the captain: ship the narrow fix + guard (recommended) or close as substantially-done-by-#290.

## Stage Report: implementation

- DONE: Edit 1 — remove ONLY the two stage-name parentheticals `(implementation, validation)` and `(ideation, backlog)` from the Split-Root State Contract worktree-isolation sentence, keeping the rule text verbatim. AC-2: `wc -c` decreases by exactly 49 bytes vs origin/main and `git diff origin/main` shows only the parentheticals deleted, nothing added.
  Verified: delta = 49 bytes (8941 → 8892); diff is a single line replacement, rule prose unchanged. Commit ca73544d.
- DONE: AC-1 guard — internal/contractlint structural-absence test (modeled on deferred_tier_absence_test.go) that fails on any parenthetical enumerating workflow stage names in the ensign core, with a paired discriminator control (mustFlag the two parentheticals; mustPass the neutral phrasing AND the incidental "Signaling done"), proving non-vacuity.
  `internal/contractlint/ensign_core_stage_absence_test.go`: scanner + discriminator + re-insertion control. Exercised non-vacuity live: defeating the regex reds both controls; re-inserting a parenthetical into the real core reds the absence check. `go test ./internal/contractlint/` green.
- DONE: AC-3 behavior-loss control — internal/dispatch test that builds a dev-shape fixture, runs dispatch build (assignment carries the show-stage-def fetch line, not inlined stage prose), runs show-stage-def asserting the `### {stage}` subsection carrying a sentinel is returned while the universal core does not carry it, plus a perturbation control (a stage without the sentinel returns no sentinel).
  `internal/dispatch/build_stage_discipline_delivery_test.go`: build emits Skill-load + fetch line, sentinel absent from the assignment body, present in show-stage-def ideation output, absent for validation. `go test ./internal/dispatch/` green.

### Summary

Removed the one residual stage-neutrality leak from the universal ensign core (the `(implementation, validation)` / `(ideation, backlog)` parentheticals in the Split-Root State Contract, −49 bytes net-negative) and re-established a quarantine-legal regression guard, since #290's prose-lock was retired and none existed. The contractlint structural-absence guard scans the core for a stage-name-enumeration parenthetical with paired discriminator + re-insertion controls; non-vacuity was exercised live (defeat-the-scanner reds the controls, re-insert-into-core reds the absence check). The dispatch behavior-loss control proves dev stage discipline rides the dev-shape show-stage-def fetch, not the core, so the edit loses no discipline. Full repo `go test ./...` green.
