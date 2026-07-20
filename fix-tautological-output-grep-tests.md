---
id: hthnpaag41m1jaxb6mxwj3y2
title: Fix 8 confirmed tautological output-grep tests (the third shape)
status: ideation
source: "Verified sweep 2026-07-19 (25 candidate files → 17 triaged → adversarial verify → 8 confirmed, 9 refuted). The 'third shape': a test with a real t.Fatalf/Errorf sink AND hand-written literals asserting rendered command/help/doc OUTPUT WORDING that no machine consumer parses — distinct from the assertion-free + mirror shapes in tautological-test-fixes. Triggered by a brittle help-output grep shipped on PR #516... wait, PR #526 (dispatch build --help) that passed 4 lenient reviews; already fixed there. Distinguishing rule: does a machine consumer parse the string, and would a real behavior change (not a rewording) flip it?"
started: 2026-07-20T03:29:35Z
completed:
verdict:
score:
worktree:
sprint: 0260-proportionality
group: test-cleanups
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

## Out of scope

The 4 assertion-free/mirror tests (owned by `tautological-test-fixes`); the contractlint prose-backstop concern (owned by `contractlint-prose-function-backstop-retirement` + `contractlint-mixed-structural-boundary`). NOTE: a `.md`-prose-grep sweep the same session found the non-contractlint packages CLEAN, so the third-shape debt is this CLI-output-grep set, not a codebase-wide markdown problem.

## Acceptance criteria

**AC-1 (VALUE) — No real regression coverage is lost by the cleanup.**
For every deleted or narrowed test, a real BEHAVIOR mutation (not a rewording) that the removed assertions nominally guarded still turns a RETAINED test RED; a mutation that stays GREEN means coverage was lost and fails this AC. This is the independent baseline that can move the wrong way: silently over-deleting shows up as a mutation the suite no longer catches.
Verified by: the per-test mutation matrix in the Test plan, exercised once in implementation — for each delete, the named sibling(s) go RED under the behavior mutation; for each narrow, the retained dynamic assertion goes RED.

**AC-2 — The 8 named tests no longer assert redundant/prose output wording, and nothing else regresses.**
Verified by: the six delete + two narrow edits land exactly as described in Proposed approach; a diff/`git grep` over the touched files shows NO new output-wording `strings.Contains` grep introduced; `go test ./internal/cli/... ./internal/dispatch/... ./internal/ensigncycle/...` is GREEN and `go test ./...` is green in CI.

## Test plan

- Surface: Go unit tests only. No fixture, CLI, or live-workflow additions. Cost: low — edits to 5 `_test.go` files plus one `go test` run; the mutation matrix is a handful of throwaway one-line source edits, each reverted after observing RED.
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
