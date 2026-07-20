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

## Out of scope

The 4 assertion-free/mirror tests (owned by `tautological-test-fixes`); the contractlint prose-backstop concern (owned by `contractlint-prose-function-backstop-retirement` + `contractlint-mixed-structural-boundary`). NOTE: a `.md`-prose-grep sweep the same session found the non-contractlint packages CLEAN, so the third-shape debt is this CLI-output-grep set, not a codebase-wide markdown problem.

## Acceptance criteria

- **AC-1** — each of the 8 is removed or narrowed so it no longer asserts redundant/prose output wording; verified by: for a narrowed test, its remaining assertions still fail under a real behavior mutation the sibling also catches; for a removed test, the named sibling behavioral test still covers the behavior (go test green, diff review). No new output-wording grep is introduced.
