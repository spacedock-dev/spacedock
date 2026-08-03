---
id: hthnpaag41m1jaxb6mxwj3y2
title: Fix 8 confirmed tautological output-grep tests (the third shape)
status: done
source: "Verified sweep 2026-07-19 (25 candidate files → 17 triaged → adversarial verify → 8 confirmed, 9 refuted). The 'third shape': a test with a real t.Fatalf/Errorf sink AND hand-written literals asserting rendered command/help/doc OUTPUT WORDING that no machine consumer parses — distinct from the assertion-free + mirror shapes in tautological-test-fixes. Triggered by a brittle help-output grep shipped on PR #516... wait, PR #526 (dispatch build --help) that passed 4 lenient reviews; already fixed there. Distinguishing rule: does a machine consumer parse the string, and would a real behavior change (not a rewording) flip it?"
started: 2026-07-20T03:29:35Z
completed: 2026-07-20T13:49:49Z
verdict: passed
score:
worktree: .worktrees/spacedock-ensign-fix-tautological-output-grep-tests
sprint: 0260-proportionality
group: test-cleanups
mod-block:
pr: pr-merge:535
archived: 2026-07-20T13:49:49Z
---

Fix (narrow to observed behavior, or remove relying on the sibling behavioral test) the 8 confirmed tautological output-grep tests. Each is verified redundant/prose by adversarial refutation — the behavioral coverage it claims is already provided by a named sibling.

## Problem

Eight tests assert rendered OUTPUT WORDING via string-contains but prove no behavior a real production change could break (they fail only on a rewording, or are strict subsets of a sibling). Confirmed by adversarial verify:

- **cli/cli_test.go:27 `TestTopLevelHelpGroupedJargonFree`** (medium, prose) — asserts the `topLevelHelp` const verbatim; redundant with `TestTopLevelHelpFormsAreIdentical` + `TestLeadingFlagBareStillHelp`.
- **cli/pi_frontdoor_test.go:571 `TestRuntimeSupportDocsKeepPiDoctorVsLiveTalkbackBoundary`** (medium, prose) — greps the dev doc `docs/runtime-support.md`, invokes zero production code; redundant with `TestPiDoctorReportsMissingAndHealthyRuntime`.
- **dispatch/build_team_name_advisory_test.go:109 `TestBuildHelpDocumentsShapeFlags`** (medium, prose) — greps `build --help` for flag names + prose phrases; redundant with `TestDispatchBuildHelpBeforeRequiredFlags` + `TestBuildLegacyTeamNameAdvisory` (which exercises the flags' real effects).
- **cli/cli_test.go:107 `TestVersionContractToken`** (low, redundant) — strict subset of `TestVersion`.
- **cli/state_commit_test.go:154 `TestStateCommitHaltStderrCarriesRemediationAndPeerCommit`** (low, prose) — ONLY line 154 (the static advisory sentence copied from state_sync.go); the rest of the test (peer-commit at :148) is genuinely dynamic — narrow, don't delete.
- **cli/verbs_test.go:94 `TestHelpListsNewVerbs`** (low, prose) — greps `--help` for verb descriptions copied from the const; the verbs' behavior is covered by `TestNewVerbMintsInDiscoveredWorkflow` + `TestCompletionShells`.
- **ensigncycle/filing_readme_template_test.go:28 `TestFilingReadmeTaskTemplateRoundTripsThroughNew`** (low, redundant) — lines 28-29 (title/status) redundant with the entity-check sub-test that also drives `new`.
- **ensigncycle/filing_readme_template_test.go:30 (same test)** (low, mirror) — five heading assertions over the test's own written fixture.

## Proposed approach

Per-test fix: six deletions (the whole test is covered by a named behavioral sibling) and two narrowings (part of the test is genuinely dynamic and stays). Test-only — no production code changes and no user-visible CLI output changes, so no doc diff.

| # | Test (file) | Fix | Retained behavioral coverage |
|---|-------------|-----|------------------------------|
| 1 | `TestTopLevelHelpGroupedJargonFree` (`cli/cli_test.go`) | delete | `TestTopLevelHelpFormsAreIdentical` (all help forms byte-identical, exit 0, empty stderr) + `TestLeadingFlagBareStillHelp` (bare `spacedock` → help + exit 0 + empty stderr). The dropped assertions (group headers/order, command-name rows, footer, jargon-absence) respond only to rewordings — no machine parses them and no behavior change flips them. |
| 2 | `TestRuntimeSupportDocsKeepPiDoctorVsLiveTalkbackBoundary` (`cli/pi_frontdoor_test.go`) | delete | Pure grep over `docs/runtime-support.md`, invokes zero production code (a banned prose-grep). `TestPiDoctorReportsMissingAndHealthyRuntime` drives `doctor --host pi` and asserts the missing/healthy output including `live child talkback` / `durable marker probe` / intercom-bridge lines. |
| 3 | `TestBuildHelpDocumentsShapeFlags` (`dispatch/build_team_name_advisory_test.go`) | delete | `TestDispatchBuildHelpBeforeRequiredFlags` (`build --help` exit 0, empty stderr, usage + flags) + `TestBuildLegacyTeamNameAdvisory` (the flags' real effect: advisory fires on the legacy arm via the `legacy TeamCreate-registry dispatch shape` marker, silent on the merged arm, envelope shape unchanged). |
| 4 | `TestVersionContractToken` (`cli/cli_test.go`) | delete | Strict subset of `TestVersion`, whose line-1 exact-equality (`"spacedock " + displayVersion() + " " + frozenContractToken`) already carries `frozenContractToken`. |
| 5 | `TestStateCommitHaltStderrCarriesRemediationAndPeerCommit` (`cli/state_commit_test.go`) | narrow | Remove the two `strings.Contains` greps for the `Next: HALT dispatch…` and `Never ` + "`" + `git push --force`… sentences — both copied verbatim from `state_sync.go:228-229` (prose-to-code mirror). Keep the exit-3 check and the dynamic `"Peer commit: "+peerSHA` assertion (peerSHA from `git rev-parse`), which itself proves the HALT emits a populated, computed stderr diagnostic rather than only an exit code. |
| 6 | `TestHelpListsNewVerbs` (`cli/verbs_test.go`) | delete | `TestNewVerbMintsInDiscoveredWorkflow` (drives `new`, mints id, `--validate` clean) + `TestCompletionShells` (drives `completion`, exit-code + script-content contract). |
| 7+8 | `TestFilingReadmeTaskTemplateRoundTripsThroughNew` required-list loop (`ensigncycle/filing_readme_template_test.go`) | narrow | Remove the loop's tautological entries: `title:` / `status:` (test 7 — re-verified by the `literal template creates entity and mints id` sub-test at :47-67, which drives `new` and asserts the created entity preserves them) and the five `##` headings (test 8 — a mirror grep over the test's OWN `filingReadme()` fixture at `shared_fixtures_test.go:318`). The loop's only non-tautological member `---\n` duplicates the adjacent `HasPrefix(template, "---\n")`, so the whole loop is removed; the structural `HasPrefix` + no-`id:` checks (:40-45) and both sub-tests stay. |

**Expected surface:** 6 `_test.go` files (`cli/cli_test.go`, `cli/pi_frontdoor_test.go`, `cli/state_commit_test.go`, `cli/verbs_test.go`, `dispatch/build_team_name_advisory_test.go`, `ensigncycle/filing_readme_template_test.go`), six test deletions + two narrowings, roughly −135 / +0 LOC (removal-only; the narrowings delete assertions and add none), no production code and no non-test files.

## Out of scope

The 4 assertion-free/mirror tests (owned by `tautological-test-fixes`); the contractlint prose-backstop concern (owned by `contractlint-prose-function-backstop-retirement` + `contractlint-mixed-structural-boundary`). NOTE: a `.md`-prose-grep sweep the same session found the non-contractlint packages CLEAN, so the third-shape debt is this CLI-output-grep set, not a codebase-wide markdown problem.

## Acceptance criteria

**AC-1 (VALUE) — No real regression coverage is lost by the cleanup.**
For every deleted or narrowed test, a real BEHAVIOR mutation (not a rewording) that the removed assertions nominally guarded still turns a RETAINED test RED; a mutation that stays GREEN means coverage was lost and fails this AC. This is the independent baseline that can move the wrong way: silently over-deleting shows up as a mutation the suite no longer catches.
Verified by: the per-test mutation matrix in the Test plan, exercised once in implementation — for each delete, the named sibling(s) go RED under the behavior mutation; for each narrow, the retained dynamic assertion goes RED.

**AC-2 — The 8 named tests no longer assert redundant/prose output wording, and nothing else regresses.**
Verified by: the six delete + two narrow edits land exactly as described in Proposed approach; a diff/`git grep` over the touched files shows NO new output-wording `strings.Contains` grep introduced; `go test ./internal/cli/... ./internal/dispatch/... ./internal/ensigncycle/...` is GREEN and `go test ./...` is green in CI.

## Test plan

- Surface: Go unit tests only. No fixture, CLI, or live-workflow additions. Cost: low — edits to 6 `_test.go` files plus one `go test` run; the mutation matrix is a handful of throwaway one-line source edits, each reverted after observing RED.
- Mutation matrix (AC-1 proof, run once in implementation; each mutation reverted immediately):
  1. bare `spacedock` prints bytes ≠ `--help` → `TestTopLevelHelpFormsAreIdentical` RED.
  2. drop a `live child talkback` / `durable marker probe` line from `doctor --host pi` → `TestPiDoctorReportsMissingAndHealthyRuntime` RED.
  3. make `build --help` exit 2 → `TestDispatchBuildHelpBeforeRequiredFlags` RED; suppress the legacy advisory → `TestBuildLegacyTeamNameAdvisory` RED.
  4. change `--version` line 1 to omit `frozenContractToken` → `TestVersion` RED.
  5. narrowed HALT test: emit a bogus peer sha on the HALT path → the retained `"Peer commit: "+peerSHA` assertion RED.
  6. break the `new` alias → `TestNewVerbMintsInDiscoveredWorkflow` RED; break `completion` → `TestCompletionShells` RED.
  7+8. mutate `filingReadme()` so the template can't round-trip through `new` (e.g., drop `title:`) → the `literal template creates entity and mints id` sub-test RED.
- Regression: `go test ./...` green after the edits.
- No spike needed: the fix relies only on removing/narrowing Go `_test.go` assertions and re-running `go test` — a proven mechanism. Each named sibling was confirmed to exist and assert the behavioral kernel by reading it (`cli_test.go:56,156,79`; `pi_frontdoor_test.go:584`; `dispatch/help_test.go:10` + `build_team_name_advisory_test.go:23`; `verbs_test.go:22,65`; `filing_readme_template_test.go:47-67`). No parser round-trip, on-disk format, runtime handoff, or tool-flag support is introduced or in question.

## Stage Report: ideation

- DONE: Confirm-the-banked-spec ideation, thin by design: the 8 named tests get behavior-observing replacements or deletion-with-reason, no scope growth beyond the 8, smallest test surface per fix
  Proposed approach fixes all 8 — 6 deletions each covered by a named behavioral sibling, 2 narrowings that keep the genuinely dynamic assertions (HALT peer-commit sha; filing `new` round-trip) and drop only source-prose/fixture-mirror greps; per-test coverage table + AC-1 mutation matrix.
- DONE: Record "no spike needed" with the proven mechanisms relied on, per the ideation stage rule - this entity is entirely dev-specific; no contract or template content
  Test plan records "no spike needed": test-only Go edits + `go test`; every sibling confirmed to exist and assert the kernel by reading (line cites). No contract, template, or user-visible CLI-output change, so no doc diff.

### Summary

Confirmed the banked 8-test list by reading each test and its named sibling(s): six deletes are fully covered by existing behavioral siblings; two narrows (state-commit HALT, filing template loop) keep the dynamic assertions and drop the tautological ones. The two HALT remediation greps are verbatim copies of `state_sync.go:228-229`; the filing heading greps mirror the test's own `filingReadme()` fixture. Value AC is a per-test mutation-survival matrix — a real behavior change still turns a retained test RED — and no spike is needed since the fix is test-only edits on a proven mechanism.

## Stage Report: implementation

- DONE: All 8 confirmed tautological output-grep tests are fixed — six deletions and two narrowings — across exactly the 6 named _test.go files, with no production code and no non-test file touched
  Commit `ff3be305` on `spacedock-ensign/fix-tautological-output-grep-tests`. Deletes `TestTopLevelHelpGroupedJargonFree` + `TestVersionContractToken` (cli/cli_test.go), `TestRuntimeSupportDocsKeepPiDoctorVsLiveTalkbackBoundary` (cli/pi_frontdoor_test.go), `TestHelpListsNewVerbs` (cli/verbs_test.go), `TestBuildHelpDocumentsShapeFlags` (dispatch/build_team_name_advisory_test.go); narrows the two source-prose greps out of the state-commit HALT test and removes the required-literal loop from the filing template test. `git diff --numstat main...HEAD` lists exactly the 6 declared `_test.go` paths and nothing else.
- DONE: The mutation matrix is RUN, not asserted: for each retained sibling that carries the behavioral coverage, a seeded one-line production break is observed RED and the revert observed GREEN, with the observed output recorded in the stage report
  10/10 mutants RED, 10/10 reverts GREEN, tree clean after (`git status --short` empty). Observed failures below.
- DONE: The stage report states the ACTUAL LOC delta and file count against the declared 6 files / six deletions + two narrowings / ~-135/+0 LOC
  See "Surface against declaration" below.

### Mutation matrix (observed)

| # | Seeded one-line production break | Retained sibling | Observed RED |
|---|---|---|---|
| 1 | `cli.go:125` bare-help path prints an extra line | `TestTopLevelHelpFormsAreIdentical` | `cli_test.go:16: help form [] output differs from --help` |
| 2 | `pi.go:636` drops the `durable marker probe` line | `TestPiDoctorReportsMissingAndHealthyRuntime` | `pi_frontdoor_test.go:622: healthy doctor output missing "durable marker probe"` |
| 3a | `dispatch.go:30` build `--help` returns 2 | `TestDispatchBuildHelpBeforeRequiredFlags` | `help_test.go:15: dispatch build --help exit=2, want 0` (both `--help` and `-h` subtests) |
| 3b | `build.go:752` advisory written to `io.Discard` | `TestBuildLegacyTeamNameAdvisory` | `build_team_name_advisory_test.go:62: must emit the legacy advisory; stderr=""` |
| 4 | `cli.go:584` version line 1 omits `frozenContractToken` | `TestVersion` | `cli_test.go:47: version first line = "spacedock 0.25.1+dev", want "spacedock 0.25.1+dev (contract 3)"` |
| 5 | `state_sync.go:225` HALT emits a bogus peer sha | narrowed `TestStateCommitHaltStderrCarriesRemediationAndPeerCommit` | `state_commit_test.go:148: HALT stderr should name the peer commit "29ef2b4"` |
| 6a | `cli.go:328` `new` alias forwards `--nope` | `TestNewVerbMintsInDiscoveredWorkflow` | `verbs_test.go:33: new stdout = "ID SLUG …", want minted-id narration` |
| 6b | `cli.go:605` `completion bash` returns 2 | `TestCompletionShells` | `verbs_test.go:70: completion bash exit=2` |
| 7 | `status/new.go:77` `new` stops stamping the minted id | narrowed filing test, `literal template creates entity and mints id` | `filing_readme_template_test.go:51: created entity has no minted id` |
| 7b | Task Template fixture loses `title:` (the plan's named mutation for the removed `title:`/`status:` greps) | same sub-test | `filing_readme_template_test.go:48: created entity did not preserve template frontmatter` |

Reproduce (no machine-specific tooling; each row is one edit): from the worktree at `ff3be305`, apply the row's break at the named `file:line`, run `go test <pkg> -run '^<SiblingTest>$' -count=1` and observe the quoted failure, then `git checkout -- <file>` and re-run to observe GREEN. Mutation 7b is the one non-production seed: the removed assertion was over the test's own fixture, so the fixture is the only place a "regression" can be seeded — recorded as such rather than dressed up as a production break.

### Surface against declaration

Declared: 6 `_test.go` files, six deletions + two narrowings, ~−135/+0 LOC, no production code, no non-test files. **Actual: 6 `_test.go` files, six deletions + two narrowings, −141/+4 LOC**, no production code, no non-test files. Well inside the default 2x tolerance; no reconfirm/re-scope/park/escalate needed. The +4 is comments only, and is a deviation from the declared removal-only shape: two doc-comment lines on the narrowed HALT test and the two `ABOUTME:` lines on `build_team_name_advisory_test.go`, both of which named the now-deleted assertions and would otherwise have been left actively false. No new assertion, helper, test, check, gate, or lint was added; `git diff main...HEAD | grep '^+[^+]'` returns those four comment lines and nothing else.

### Summary

Six deletions and two narrowings landed across exactly the 6 declared `_test.go` files at −141/+4 (comments only), commit `ff3be305`. AC-1 was proved by running the matrix rather than asserting it: ten seeded one-line breaks each turned the named retained sibling RED with the failure text recorded above, and every revert returned GREEN with a clean tree — no mutation survived, so no real regression coverage was dropped. `go test ./...` is green across all 17 packages. No replacement test was written and no enforcement mechanism was added, per the sprint constraints; the one honest caveat is mutation 7b, which necessarily seeds a test fixture rather than production because the assertion it retires was a mirror of that fixture.

### Review findings (roborev job 325, panel branch_final)

Three findings. One accepted and fixed; two declined with a zero-line diff.

**ACCEPTED — Low: test name overstates its remaining coverage** (`internal/cli/state_commit_test.go:125`). Renamed `TestStateCommitHaltStderrCarriesRemediationAndPeerCommit` -> `TestStateCommitHaltStderrCarriesPeerCommit`, doc-comment reflowed to match; no assertion changed. Commit `db34be27`. This is the same false-signal defect as the stale doc-comments already fixed in `ff3be305`, in a louder place: a name promising `CarriesRemediation` when remediation is no longer asserted locally is exactly the overstated-coverage shape this sprint removes.

**DECLINED — Medium: "help contracts are no longer fully asserted; add golden fixtures or structured assertions covering the complete required help output."** Grounds: (a) already adjudicated at the approved gate — the dropped assertions (group headers/order, command-name rows, footer, jargon-absence) respond only to rewordings; the sprint's own preflight staff review raised this independently and recorded a DECLINE (`docs/roadmap/0260-proportionality/staff-review-fable-delta.md:58`). (b) The proposed remedy re-creates the anti-pattern at larger scale: an implementer-written golden fixture of help output has no independent source that can diverge, so it fails on any rewording and passes on any behavior change — the brittleness that triggered this task. (c) The finding's specific claim that `--host`, `--team-name`, `--bare-mode` are left unverified is factually wrong; I checked rather than assumed. `grep -rln 'bare_mode\|--bare-mode\|BareMode' internal --include='*_test.go'` returns 24 test files driving the flag behaviorally (`build_advance_test.go`, `build_merged_mode_test.go`, `ensigncycle/cycle_test.go`, others); `build_pi_host_test.go` drives `host: pi`; `TestBuildLegacyTeamNameAdvisory` drives `team_name` (mutation 3b above). Delete a flag and those go RED — the help grep was never the binding. Promote-to-material condition: an observed help-content regression that no behavioral test catches.

**DECLINED — Medium: "the state-conflict test no longer verifies remediation guidance; assert the complete stderr contract, preferably with an exact fixture."** Grounds, verified specifically because a safety boundary (never force-push / never auto-resolve) deserves more than a reflexive decline: (a) the guidance has ONE production emitter, `internal/cli/state_sync.go:228-229`. (b) Both literals are still asserted at `internal/cli/state_ready_test.go:115` (the `Next: HALT dispatch` line) and `:118` (the `Never git push --force` safety line) — if production stops emitting either, that test goes RED, so the safety guidance cannot silently disappear. (c) The narrowing kept the exit-3 check and the dynamic computed peer-commit assertion, which proves the HALT emits a populated, computed diagnostic rather than only an exit code (mutation 5 above). Promote-to-material condition: the remediation emitter gains a second call site, or `state_ready_test.go:115`'s assertion is removed — either would leave this path genuinely unasserted.

**OBSERVATION (not fixed — outside the declared surface).** `internal/cli/state_ready_test.go:115` and `:118` are themselves the third shape: literals copied verbatim from `state_sync.go:228-229`. They are load-bearing today only because they are the last assertion of that emitter, which is why decline 2 stands. They are not among the declared 8 and are not this entity's to touch — logged as a next-train candidate so the debt is not lost.

### Final surface

**6 `_test.go` files, +8 / -145 (net -137)** against the declared 6 files / six deletions + two narrowings / ~-135/+0. Inside the default 2x tolerance; no reconfirm/re-scope/park/escalate. The only added non-comment line in the whole branch is the renamed `func TestStateCommitHaltStderrCarriesPeerCommit(t *testing.T) {`, which replaces the removed old signature — `git diff main -- . | grep '^+[^+]' | grep -v '^+\s*//'` returns that one line. Still no production code, no non-test file, no new assertion, test, check, gate, or lint. `go test ./...` green across all 17 packages after the rename.

## Stage Report: validation

- DONE: Every AC-N is verified with evidence you REPRODUCE yourself. Re-run the mutation matrix independently — for each of the 10 rows, seed the named one-line break at the named file:line, observe the named sibling go RED with the quoted failure, revert, observe GREEN
  10/10 reproduced from a clean `db34be27` tree by an independently written driver (seed → `go test <pkg> -run '^<Sibling>$' -count=1` → `git checkout --` → re-run): 10/10 mutants RED, 10/10 reverts GREEN, `git status --short` empty after. Every named `file:line` matched the expected source text (the driver hard-fails on seed mismatch), so the cited locations are accurate. AC-1 holds: no behavior mutation survives.
- DONE: Run the offline lane (`go test ./...` and `go test ./... -race`) and report the result; also confirm the FO prompt-surface ratchet test still passes
  `go test ./... -count=1` green (15 packages with tests + 2 without); `go test ./... -race -count=1` green, same 17. No package broke on a helper or fixture the deleted tests carried. `TestFOFunctionPromptSurfaceShrinks` + `TestFOFunctionReferenceInvariant` + `TestFOFunctionReferenceCheckpointMetrics` PASS (run explicitly, not inferred). AC-2's targeted lane `./internal/cli/... ./internal/dispatch/... ./internal/ensigncycle/...` green. Captain's pi-lane-red waiver cited but not needed — nothing red anywhere.
- DONE: Check that BOTH roborev declines are properly RECORDED with grounds and promote-to-material conditions, and that the accepted rename left no stale reference to the old test name
  Both recorded with lettered grounds and an explicit promote-to-material condition. Decline 1's factual claim reproduces exactly: `grep -rl 'bare_mode\|--bare-mode\|BareMode' internal --include='*_test.go'` → 24 files, `build_pi_host_test.go` drives `host: pi`. Decline 2's reproduces exactly: `state_sync.go:228-229` are the SOLE production emitters and `state_ready_test.go:115,118` the SOLE surviving assertions. Rename clean — no reference to the old name survives; the two `CarriesRemediationAndPeerCommit` hits are `TestStateReadyHaltStderrCarriesRemediationAndPeerCommit`, a different test whose name is still accurate. The `state_ready_test.go:115` next-train observation is recorded and untouched.

### Surface and claim checks

`git diff --numstat main...HEAD` → exactly the 6 declared `_test.go` paths, 8 insertions / 145 deletions; the only added non-comment line is the renamed func signature. No production code, no non-test file, no added assertion. Two entity claims independently confirmed rather than taken on trust: the deleted `TestRuntimeSupportDocsKeep…` was a pure `os.ReadFile` grep over `docs/runtime-support.md` invoking zero production code, and the filing-template README parsed by the narrowed test is the test's own `filingReadme()` fixture written into `t.TempDir()` — so the removed heading assertions had zero production coverage.

### Deferred risk (does not block)

**Top-level help CONTENT has no surviving guard.** Measured, not assumed: five seeded edits to `topLevelHelp` (`internal/cli/help.go`) — drop the `new` row, drop the `completion` row, drop the `Setup` group header, reorder the group headers, inject the banned jargon `META` — all SURVIVE a full `go test ./... -count=1` on the branch. Confirmed as a real delta: restoring `main`'s `cli_test.go`/`verbs_test.go` against the same seed turns `TestTopLevelHelpGroupedJargonFree` RED. Counter-datum: dropping the `merge guard` row IS caught, by the surviving `internal/cli/merge_test.go:106`. Why deferred, not material: AC-1 scopes to behavior mutations, and none of the five changes any behavior — the commands keep working, and matrix rows 6a/6b prove `new` and `completion` are themselves pinned behaviorally; no machine consumer parses this help. This is the measured form of the already-adjudicated roborev decline 1, attached as numbers rather than reversed. Promote-to-material: an observed help-content regression that reaches a user, or a machine consumer begins parsing top-level help.

### Polish (does not block)

- Matrix row 6a's recorded failure quote is not byte-reproducible. The seed "`cli.go:328` `new` alias forwards `--nope`" does not pin where `--nope` goes; my placement yields `verbs_test.go:30: new exit=1 stderr="Error: --new requires a slug argument"` rather than the recorded `verbs_test.go:33: new stdout = …`. Same sibling, same RED, under-determined transcript.
- Two retained siblings are themselves the third shape: `merge_test.go:106` (help-row grep) and `dispatch/help_test.go:10` (help-content grep). Correctly outside the declared 8 — noting for the next-train ledger beside the already-recorded `state_ready_test.go:115`.

### Recommendation: PASSED

Material findings: none.

### Summary

Re-ran the AC-1 mutation matrix independently rather than accepting the recorded table: all 10 rows reproduce, mutant RED and revert GREEN, from a clean tree that stayed clean. No behavior mutation survives the deletions, so no real regression coverage was dropped and the central question this stage exists to answer is settled in the implementation's favor. Both offline lanes (`./...` and `./... -race`) are green across all 17 packages and the FO prompt-surface ratchet still passes; both roborev declines are honestly recorded with grounds and promote conditions, and each decline's load-bearing factual claim reproduces exactly. The one substantive residual — top-level help content is now unguarded, proven by five surviving seeded edits — is a wording-drift deferred risk under AC-1's own behavior/rewording line, already adjudicated at the approved gate; recorded with its exact trigger and promote condition rather than reopened.
