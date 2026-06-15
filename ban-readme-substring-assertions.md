---
id: 5h0chdcad99dq0z50qzwf1za
title: Replace README-substring (prose-grep) test assertions — prove the seam, not the prose
status: ideation
source: "captain + nb (readme-main-flip-reconciliation) reconciliation 2026-06-07 — PR #315 edits tests/test_codex_plugin_packaging.py to assert README content by substring (assert \"spacedock codex\" in readme). Not on next today; when #315's content lands on main it must be replaced. Same proof-policy class as #309/4q and the survey signal-correction work."
started: 2026-06-15T01:55:42Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0203-fo-efficiency
group: proof-policy
sprint-readiness:
---

A test that asserts README content by substring (`assert "spacedock codex" in readme`) proves only that the text is present — the banned prose-grep tautology. When PR #315's content lands on `main`, replace that assertion with a behavioral one.

## Problem

PR #315's `tests/test_codex_plugin_packaging.py` greps the README for substrings to "test" Codex install/launch. A substring match over a doc the implementer wrote can't fail on a valid paraphrase and passes on an inverted clause — it tracks spelling, not behavior. This is exactly the proof-policy class the project fences (the 4q reframe, #309, the survey signal-correction `references/queries.sql` work).

## Proposed approach

When #315's content reaches `main` (or at the flip), replace the README-substring assertion: prove Codex install/launch by exercising the SEAM (the packaging/launch behavior — does `spacedock codex` resolve + launch), not by grepping README prose. If the only real invariant is "the README documents the command," that is not a behavioral test and should be dropped, not kept as a grep.

## Out of scope

- Rewriting #315 itself (it's a separate PR; this task handles the assertion once it lands).
- The general prose-grep ban already shipped for skills (#309/4q) — this is the README/packaging-test instance.

## Acceptance criteria

Ideation/implementation fills in. Sketch: no test asserts README content by substring/regex; Codex install/launch is proven by exercising the packaging/launch seam (output/exit/observable behavior), with a fixture or live check that can actually fail on broken behavior.

## Test plan

The replacement test goes RED on broken Codex install/launch behavior (not on a paraphrase of the README).

## FO note — concrete trigger (2026-06-14)

The dev-README slim (commit `a9e669ae`, FO-direct, unpushed local `main`) relocated the Runtime-Live-CI / Codex-watchdog content out of `docs/dev/README.md` into `docs/runtime-live-ci.md`, which EXPOSED two prose-grep doc-contract guards this task targets:

- `internal/ensigncycle/shared_scenarios_docs_test.go` — `TestSharedScenarioDocsContract`
- `internal/ensigncycle/codex_collab_wait_watchdog_test.go` — `TestCodexForegroundWaitWatchdogDocsContract`

Both assert literal clause strings are present in `docs/dev/README.md` (e.g. `### Shared runtime scenarios`, `Codex foreground-wait watchdog`, `codexScenarioRunners()`); they now fail (clause count 0) because the prose moved. Same README-substring anti-pattern this task removes — surfaced by the 0.20.4 `e6a` implementation, which found them red at its branch base.

Blocking relationship: this breakage blocks BOTH the README-slim push AND any 0.20.4 `e6a` merge (e6a's base includes `a9e669ae`). Resolution (this task's call): retarget the guards to `docs/runtime-live-ci.md`, or convert them to bind an independent source per the proof policy — not a relocated prose-grep.

## Ideation (2026-06-15)

### Call: drop, do not retarget

The FO note offered two resolutions for the README guards — retarget to `docs/runtime-live-ci.md`, or rebind to an independent source. Ideation's verdict: **drop both README guards.** Retargeting just relocates the prose-grep tautology to a new file; it cures the red without curing the anti-pattern. Rebinding to an independent source is the right move only when a real seam is left UNPROVEN by the rest of the suite — and here it is not. Both README guards are redundant atop behavioral tests that already exist and red on real breakage (proven below). The seed's explicit guidance: where an assertion's only invariant is "the doc documents X" with no behavior, DROP it rather than keep it as a grep.

### The prose-grep assertions to replace, and the real seam each greps past

The named triggers, plus the siblings a sweep surfaced. For each: the prose it greps vs. the real behavioral seam, and the disposition.

1. **`internal/ensigncycle/shared_scenarios_docs_test.go` :: `TestSharedScenarioDocsContract`** — reads `docs/dev/README.md`, asserts ~13 clauses present (`## Runtime Live CI`, `### Shared runtime scenarios`, `### Local live execution`, `sharedRuntimeScenario`, `runner adapter`, `codexScenarioRunners()`, `claudeScenarioRunners()`, `pi_shared_coverage_test.go`, `To add a shared runtime scenario`, `TestSharedScenarioRunnerCoverage`, and three `go test -tags live …` command strings).
   - **Real seam:** "every shared scenario has a Codex runner, a Claude runner, and a Pi coverage entry, both directions." That seam is ALREADY proven, executably, by `internal/ensigncycle/shared_coverage_meta_test.go` :: `TestSharedScenarioRunnerCoverage` (binds `sharedRuntimeScenarios()` ↔ `codexScenarioRunners()`/`claudeScenarioRunners()`/`piSharedScenarioCoverageMap()`, reds on drift either way) and by `TestSeedScenariosDocLock` (binds the doc's `<!-- seed-scenarios -->` IDs ↔ the code table). The README clauses that name code symbols (`codexScenarioRunners()` etc.) are proven by those tests; the rest (section anchors, the how-to sentence, the live-command strings) are pure doc-presence with no behavior behind them.
   - **Disposition: DROP.** No behavioral coverage is lost — `TestSharedScenarioRunnerCoverage` + `TestSeedScenariosDocLock` remain as the real binding. The live-command strings document an operator runbook; their correctness is a doc-quality concern, not a behavioral invariant a unit test can prove (the strings can be present and wrong, or absent and the suite still runs).

2. **`internal/ensigncycle/codex_collab_wait_watchdog_test.go` :: `TestCodexForegroundWaitWatchdogDocsContract`** — reads `docs/dev/README.md`, asserts 4 clauses (`Codex foreground-wait watchdog`, `` `collab:wait` / `wait_agent` ``, `durable workflow-state progress`, `typed stall`).
   - **Real seam:** "a repeated/silent foreground `collab:wait`/`wait_agent` stall fires a typed stall and kills the proc; durable workflow-state progress disarms it." That seam is ALREADY proven, in the SAME FILE, by `TestCodexCollabWaitWatchdogRepeatedWaitStall`, `…SilentAfterWaitStall`, `…DurableProgressBeforeBudgetClearsSilentWait`, `…DurableProgressDisarmsRepeatedWaitEpoch`, and `…PositiveControls` — fixtures that drive the watchdog and assert the typed `codexCollabWaitStallError` arm, `proc.wasKilled()`, and the disarm paths. The doc-grep proves only that the README spells the feature's name.
   - **Disposition: DROP.** The five behavioral tests are the proof; the prose-grep adds nothing.

3. **`internal/contractlint/codex_foreground_wait_shape_test.go` :: `TestCodexForegroundWaitSectionCarriesOperatorInterruptionShape`** (reads `skills/first-officer/references/codex-first-officer-runtime.md`) and **`:: TestCodexIdleProbeForegroundWaitInterruptionIsNonTerminal`** (reads `docs/dev/codex-idle-notification-probe.md`) — each calls `requireSectionContains(...)` for ~6 bare prose phrases (`Esc`, `operator interruption`, `returns control`, `same handle`, `non-terminal`, …).
   - **Real seam vs. prose:** The same file already has the RIGHT model: `foregroundWaitLifecycleClaimError(section)` with positive controls (`TestForegroundWaitLifecycleClaimRejectsTerminalMutations`) and negated controls (`…AcceptsNegatedClaims`) — a semantic oracle that reds when the clause's MEANING inverts (worker-is-terminalized vs. worker-is-not), which a bare `strings.Contains("operator interruption")` cannot do. The bare `requireSectionContains` presence asserts are pure prose-grep: a meaning-inverting paraphrase keeps every grepped token.
   - **Disposition: DROP the bare `requireSectionContains` presence asserts.** The lifecycle-claim semantic oracle (`foregroundWaitLifecycleClaimError` + its two control tests) stays — it is the legitimate meaning-bearing check and is the model the README drops point at. (These live inside the `internal/contractlint` quarantine, which the boundary guard PERMITS to read instruction prose; but the package doc bans prose-grep "here and everywhere," so the bare presence asserts violate the guard's own stated contract even though the structural sweep does not yet catch in-quarantine prose-grep.)

4. **`internal/ensigncycle/teardown_marker_consistency_test.go` :: `TestGradeMarkerMatchesContract`** — reads `skills/first-officer/references/claude-fo-merge.md` and `skills/using-claude-team/SKILL.md`, asserts the CODE const `terminalTeardownMarker` appears verbatim.
   - **Judgment: NOT a pure prose-grep; OUT OF SCOPE.** This binds a code symbol (a const the grade matches) to its verbatim copies in the contract — a cross-source equality (code↔doc), the same shape as `TestSeedScenariosDocLock`, not a "doc says X" tautology. The boundary guard's package doc calls the code↔prose shape "CODE-BOUND-AS-BEHAVIOR-SUBSTITUTE" and notes the BEHAVIOR must be proven by a test that RUNS it — but auditing whether the live FO actually emits the marker is the teardown-discipline track's concern, not this task's. This task removes prose-as-behavioral-proof; it does not re-litigate the marker doc-lock. Listed here so the determination is on the record. (If the gate disagrees, fold it into a follow-up; do not expand this task's scope silently.)

### Why these slip past the existing boundary guard

`internal/contractlint/boundary_guard_test.go` :: `TestNoInstructionReadsOutsideQuarantine` already bans instruction-file reads outside the quarantine, and its package doc declares prose-grep "BANNED, here and everywhere." Two gaps let the targets through:

- **Scope gap:** the detector's `instructionPathSegments` set (`instruction_read_detector_test.go`) recognizes `skills`/`references`/`agents`/… but NOT `docs/dev/README.md` or `docs/runtime-live-ci.md`. So a doc-site README is not classified as an instruction surface, and `directlyReadsInstructionFile` returns false for the two README guards — they read `docs/dev/README.md`, which carries none of the recognized segments. That is precisely why they were never caught.
- **In-quarantine gap:** the contractlint quarantine is exempt from the read-location sweep, so the bare `requireSectionContains` prose-greps inside `codex_foreground_wait_shape_test.go` are structurally permitted even though the package doc bans prose-grep.

The structural fix is to extend the detector to recognize the doc-site README surfaces, so the boundary guard reds if a README prose-grep is ever re-introduced anywhere. That makes the "no test asserts doc content by substring as behavioral proof" AC a re-runnable code fact, not a one-time cleanup.

### Proposed approach

1. **Delete** `TestSharedScenarioDocsContract` (from `shared_scenarios_docs_test.go`) and `TestCodexForegroundWaitWatchdogDocsContract` (from `codex_collab_wait_watchdog_test.go`). Leave `TestSeedScenariosDocLock` and the five `TestCodexCollabWaitWatchdog*` behavioral tests untouched — they are the real proof. This unblocks the README-slim push and the 0.20.4 `e6a` merge (both reds disappear because the guards that depended on relocated prose are gone, not relocated).
2. **Delete** the bare `requireSectionContains(...)` presence asserts in `TestCodexForegroundWaitSectionCarriesOperatorInterruptionShape` and `TestCodexIdleProbeForegroundWaitInterruptionIsNonTerminal`. Keep `requireForegroundWaitLifecycleClaim` / `foregroundWaitLifecycleClaimError` and its two control tests. If a test function is left with only the lifecycle-claim call, keep it (it still drives the semantic oracle); if a function would be left empty, remove it.
3. **Extend the boundary guard** so a README/doc-site prose-grep cannot return. Add the doc-site README surfaces to the detector's recognized instruction surfaces (e.g. add `docs/dev/README.md` and `docs/runtime-live-ci.md`, or a `docs/` + `README.md` recognition, to `instructionPathSegments`/`isInstructionPathLiteral`) so `directlyReadsInstructionFile` flags a test that reads them. Extend `TestBoundaryGuardDetectsAPlantedInstructionRead` (`boundary_guard_control_test.go`) with a planted README-prose-grep fixture proving the extended detector REDS on it and stays GREEN on a non-prose README read (e.g. a fixture that copies a README as test input without `strings.Contains` over its prose). This is the independent oracle for AC-3.

### Doc changes

None. This task touches only test code (`*_test.go`) and the contractlint detector. No user-visible CLI output, banner, or docs-site content changes. The README prose that the dropped guards greped is unaffected — it simply stops being asserted by a test.

### Spike — riskiest unverified mechanism

The riskiest assumption is "the behavioral tests that make the README guards redundant actually exist and red on real breakage" — if they did not, dropping the guards would lose coverage. **Exercised during ideation** (not deferred): confirmed by direct read that `TestSharedScenarioRunnerCoverage` (binds scenario table ↔ both runner maps + Pi coverage, both directions) and the five `TestCodexCollabWaitWatchdog*` tests (drive the watchdog through fixtures, assert typed stall + `proc.wasKilled()` + disarm) are present and behavioral. The detector scope gap was also confirmed by reading `instructionPathSegments` — `docs/dev/README.md` carries no recognized segment, which is why the guards slipped through. No runtime/parser/on-disk-format mechanism is newly introduced; the boundary-guard extension reuses the existing AST detector and its mutation-control harness. **No further spike needed:** the implementation rests on (a) Go test deletion, (b) extending an existing AST-literal predicate, (c) the existing planted-fixture mutation control — all proven mechanisms in this package.

## Acceptance criteria

- **AC-1 (README guards gone; behavioral proof intact).** Neither `TestSharedScenarioDocsContract` nor `TestCodexForegroundWaitWatchdogDocsContract` exists in the test tree. `TestSeedScenariosDocLock`, `TestSharedScenarioRunnerCoverage`, and the five `TestCodexCollabWaitWatchdog*` tests still exist and pass. The shared-scenario parity seam and the watchdog-stall seam remain proven by tests that red on broken behavior, not on relocated prose.
  - *Test:* `go test ./internal/ensigncycle` passes; `grep -rn 'TestSharedScenarioDocsContract\|TestCodexForegroundWaitWatchdogDocsContract' internal/` returns no hits; `TestSharedScenarioRunnerCoverage` reds when a runner is removed from one host map (existing mutation behavior — verify by a temporary local deletion, not a committed test).
- **AC-2 (prose-grep presence asserts dropped; semantic oracle kept).** The bare `requireSectionContains` presence asserts in `codex_foreground_wait_shape_test.go` are removed; `foregroundWaitLifecycleClaimError` and its positive/negated control tests remain and pass. The retained foreground-wait check reds on a meaning-inverting paraphrase (terminal vs. non-terminal), which a substring presence check could not.
  - *Test:* `go test ./internal/contractlint` passes; `TestForegroundWaitLifecycleClaimRejectsTerminalMutations` / `…AcceptsNegatedClaims` still present and green; the section no longer asserts bare prose tokens by presence.
- **AC-3 (no doc/README substring used as behavioral proof; structurally enforced).** After the change, no `*_test.go` asserts `docs/dev/README.md` (or `docs/runtime-live-ci.md`) content by `strings.Contains`/regexp as a behavioral proof, and the boundary guard's detector recognizes those doc-site README surfaces so a re-introduction reds the guard.
  - *Test:* the extended `TestNoInstructionReadsOutsideQuarantine` passes (zero offenders) and `TestBoundaryGuardDetectsAPlantedInstructionRead` reds on a planted README-prose-grep fixture and stays green on a non-prose README read — an independent oracle (the AST detector over real parsed source), never a doc grep. Backstop grep: `grep -rln 'docs/dev/README.md\|docs/runtime-live-ci.md' --include='*_test.go' internal | xargs grep -l 'strings.Contains'` returns no behavioral-grep reader.

## Test plan

- **Cost/complexity:** low. Pure Go test-tree edit plus one AST-detector extension; no live model, no fixtures driving the binary. Runs in the standard offline `go test` pass.
- **What verifies it:**
  - `go test ./internal/ensigncycle ./internal/contractlint` — confirms the surviving behavioral tests pass and the boundary guard (extended) passes (AC-1, AC-2, AC-3 green path).
  - The extended `boundary_guard_control_test.go` planted-fixture pair — confirms the detector REDS on a README prose-grep and stays GREEN on a non-prose README read (AC-3 oracle reds on the banned shape; not a doc grep).
  - Temporary local mutation (not committed): remove a runner from one host map and confirm `TestSharedScenarioRunnerCoverage` reds; this demonstrates the seam the dropped README guard never proved is still independently proven (AC-1).
- **Fixture/CLI/live:** Go unit + AST-fixture only. No CLI behavior fixture, no live workflow — the claim is "tests prove behavior, not prose," verifiable entirely offline.

## Stage Report: ideation

- DONE: Enumerate every README/doc-substring (prose-grep) test assertion to replace, naming the real seam vs. the prose.
  Four candidates classified in "The prose-grep assertions to replace": the two named README guards (`TestSharedScenarioDocsContract`, `TestCodexForegroundWaitWatchdogDocsContract`) + two `requireSectionContains` sites in `codex_foreground_wait_shape_test.go`; `TestGradeMarkerMatchesContract` examined and ruled OUT OF SCOPE (code↔doc binding, not pure prose-grep).
- DONE: Design the behavioral replacements — prove the seam, or DROP where the only invariant is "doc documents X".
  Verdict: DROP both README guards (seams already proven by `TestSharedScenarioRunnerCoverage`/`TestSeedScenariosDocLock` and the five `TestCodexCollabWaitWatchdog*` tests — confirmed present by read); DROP the bare presence asserts, keep the `foregroundWaitLifecycleClaimError` semantic oracle.
- DONE: ACs bound to replacement tests reding on broken behavior; AC that no test asserts doc content by substring as behavioral proof (structural code fact).
  AC-1/AC-2 bound to the surviving behavioral tests + a mutation check; AC-3 extends the boundary-guard AST detector to recognize `docs/dev/README.md`/`docs/runtime-live-ci.md`, proven by an extended planted-fixture mutation control — a code fact, not a doc grep.

### Summary

Investigated the named triggers and swept the test tree. The two named README doc-contract guards are pure prose-grep over `docs/dev/README.md` whose real seams are ALREADY proven by existing behavioral tests in the same packages (the scenario-runner parity guard and the five watchdog-stall fixtures, confirmed by read) — so the call is DROP, not retarget (retargeting just relocates the tautology). Two sibling `requireSectionContains` presence asserts drop in favor of the existing `foregroundWaitLifecycleClaimError` semantic oracle. The doc-substring guards slipped past the existing boundary guard because its detector does not classify `docs/dev/README.md` as an instruction surface; AC-3 closes that gap and makes "no doc-substring behavioral proof" a re-runnable code fact. `TestGradeMarkerMatchesContract` is a code↔doc binding, not pure prose-grep, and is ruled out of scope. No doc changes; offline Go-only test plan; riskiest assumption (the redundant-coverage tests exist and red on real breakage) exercised during ideation.
