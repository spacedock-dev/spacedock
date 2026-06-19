---
id: 3e2bb3r432vbgp8nnw54wnjk
title: Wrong-root boot guard mis-parses a compound `cd X; …` boot (strip the trailing `;`)
status: implementation
source: 'v0.22.0 cut (2026-06-19): the opus claude-live lane on PR #390 went red on TestLiveClaudeSharedScenarios/rejection-flow — a FALSE POSITIVE. The cut shipped on the sonnet proof with an e2e-gate waiver; this is the recorded follow-up to make the guard robust so the live lane is clean-green (no waiver) next time.'
started: 2026-06-19T07:45:59Z
completed:
verdict:
score: 0.4
worktree: .worktrees/spacedock-ensign-wrong-root-guard-compound-boot
issue:
sprint: 0221-layered-fo
group: cleanup
sprint-readiness: ready
---

Harden the wrong-root boot guard so a compound `cd X; …` boot command is parsed correctly, eliminating a false-positive that latest-opus trips.

## Problem

`bootPathArgs` (`internal/ensigncycle/wrong_root_detect_impl_test.go:106`) does `strings.Fields(command)` and takes the token after `cd` as the boot target. When the FO issues a compound boot in one bash call — `cd /tmp/…/001; ls -la; echo "===README==="; sed -n '1,60p' README.md` — `strings.Fields` yields `["cd", "/tmp/…/001;", …]`, so the cd target comes out as `/tmp/…/001;` with the trailing `;` glued on (no space before it). `isUnder("/tmp/…/001;", "/tmp/…/001")` is then false → `detectWrongRootBoot` fires "FO booted the wrong root" though the FO booted the CORRECT root and merely chained exploration after it.

The parser's own comment concedes the assumption: *"It splits on whitespace (the boot commands are simple `cd …` forms; quoting is not exercised by the real boot stream)."* That held for pinned 2.1.177 and for sonnet, but **latest-opus composes its boot as a chained `cd X; ls; cat README` command**, which violates the assumption. The guard's "a CI env leak likely lured the FO off its cwd" message is a misleading built-in hypothesis — the real cause is the unstripped `;`.

Evidence: v0.22.0 release run — claude-live (opus, CI-E2E-OPUS) red on this; claude-live (sonnet) + the merged-team-mode lane + all deterministic lanes green. The 0.22.0 e2e-gate was waived for exactly this false-positive.

## Proposed approach

Strip a trailing shell-separator token (`;`, `&&`, `|`) from each path argument `bootPathArgs` extracts (the `cd` target and each `--workflow-dir` value), so a separator glued directly to the path with no preceding space — `/tmp/…/001;` — resolves to the clean path `/tmp/…/001` before `wanderTarget`'s `isUnder` check. This is the minimal fix: it leaves the whitespace split intact (the space-separated `cd X && ls` form already parses correctly — `&&` lands in its own field), and touches only the three extraction points in `bootPathArgs`. Full command-splitting on `;`/`&&`/`|` boundaries is unnecessary scope — the only failure mode is a separator glued to a path token, which a per-token trim covers.

Keep it a pure parsing fix in the test harness (`bootPathArgs` lives only in `wrong_root_detect_impl_test.go`, a `_test.go` file invisible to product code) — no product code changes, no contract changes.

### Spike (mechanism verified, ideation 2026-06-19)

Reproduced the bug and confirmed the fix offline against copies of `bootPathArgs`/`wanderTarget`/`isUnder` before committing to the design:

- **Cause confirmed:** the opus compound boot `cd <fixtureRoot>; ls -la; echo …; sed …README.md` makes `strings.Fields` yield the path token `<fixtureRoot>;` (trailing `;` glued on). `filepath.Clean` keeps the `;` (a valid path char), `filepath.IsAbs` is still true, and `isUnder("<fixtureRoot>;", "<fixtureRoot>")` is FALSE → `wanderTarget` reports a wander → false-positive wrong-root error, though the FO booted the CORRECT fixture root.
- **Fix confirmed BOTH directions:** with a trailing-separator trim, all eight probe cases pass — the two false-positive forms (`cd <fixtureRoot>;` and the `&&`-glued `cd <fixtureRoot>&& ls`) no longer flag, while the two genuine wrong-root forms (`cd <realRepo>; …` compound and `cd <realRepo> && …` spaced) still flag, and the existing `--workflow-dir` inside/outside and no-cd cases are unchanged.

No remaining unverified mechanism: the fix relies only on `strings.TrimSuffix` and the existing `filepath`/`isUnder` semantics, all exercised in the spike.

## Out of scope

- The broad-search overstep guard (`tq0`) — a sibling boot-guard hardening, distinct mechanism.
- Any product/FO-contract change — this is a test-harness parser fix.

## Acceptance criteria (proof = behavior, never prose-grep)

**AC-1 — a compound `cd X; …` boot is NOT flagged wrong-root when X is the fixture root (false-positive eliminated).**
The finished `detectWrongRootBoot` returns nil for a chained boot whose `cd` target is the fixture root with a separator glued on (`cd <fixtureRoot>; ls -la; echo "===README==="; sed -n '1,60p' README.md`), and likewise for the `&&`-glued form (`cd <fixtureRoot>&& ls`).
Verified by: new `TestDetectWrongRootBoot` subcases — `compound_cd_into_fixture_passes` (the `;`-glued opus boot) and `compound_cd_into_fixture_amp_passes` (the `&&`-glued form) — asserting `detectWrongRootBoot` returns nil. Written failing-first against the current parser (they red with the false-positive wrong-root error), then green after the trim.

**AC-2 — a genuine wrong-root compound boot STILL trips the guard (no false-negative regression).**
The finished `detectWrongRootBoot` still returns a wrong-root error for an off-fixture compound boot (`cd <realRepo>; ls -la; cat README.md`), with the error naming both the expected fixture root and the actual wandered-to path.
Verified by: a new `TestDetectWrongRootBoot` subcase `compound_cd_away_from_fixture_reds` asserting a non-nil error containing both `<fixtureRoot>` and `<realRepo>`. This case passes on the current parser too (the `;`-glued off-fixture path is still outside the fixture) — it is the guard against the trim over-stripping and silently disabling detection. Run the full subtest set after the trim to confirm every prior direction (`--workflow-dir` inside/outside, workflow-README, contract-skill read, plain cd, empty stream) stays green.

**AC-3 — the live opus claude-live lane goes green on the shared scenarios without a waiver.**
Verified by: a Runtime Live E2E run (opus) with `TestLiveClaudeSharedScenarios` green (no wrong-root false positive), so the next cut needs no e2e-gate waiver for this cause. Integration confirmation of AC-1/AC-2 in the real boot stream; naturally exercised by the next live-e2e run, not a separate code change.

## Test plan

Failing-first offline table test in `wrong_root_detect_test.go` (`go test ./internal/ensigncycle/ -run TestDetectWrongRootBoot`):

1. Add the AC-1 subcases (the `;`-glued and `&&`-glued cd-into-fixture booleans) and run — they MUST red against the current `bootPathArgs` with the false-positive wrong-root error, proving the test exercises the bug.
2. Add the AC-2 subcase (compound cd-away-from-fixture reds) — green even pre-fix; it pins the no-false-negative direction.
3. Apply the trailing-separator trim to the three extraction points in `bootPathArgs` (cd target, `--workflow-dir` spaced value, `--workflow-dir=` value).
4. Re-run the full `TestDetectWrongRootBoot` set — all subcases (new + the six existing) MUST be green.

The deterministic offline test is the primary proof (cost: minutes). The live opus lane green (AC-3) is the integration confirmation, exercised by the next live-e2e run. No fixture or live-workflow test is needed for AC-1/AC-2 — the parser bug is fully reproducible from synthetic boot-command strings, as the spike showed.

## Documentation impact

None. `bootPathArgs`/`detectWrongRootBoot` are test-harness code (`wrong_root_detect_impl_test.go`, a `_test.go` file): zero product-code callers (`grep detectWrongRootBoot --include="*.go"` outside `_test.go` is empty), zero contract/skill text references, no CLI output, banner, or docs-site surface touched. The only doc-tree hits are this entity and its own debrief/archive records. So there is no doc diff for the ideation gate to review.

## Related

- `tq0` zero-discover-broad-search-hardening — sibling boot-guard flakiness (broad-search detector).
- v0.22.0 / PR #390 — where this surfaced; the e2e-gate waiver recorded the cause.

## Stage Report: ideation

- DONE: AC + test plan pinning BOTH directions — a compound `cd X; ls; cat README` boot parses correctly (no false-positive wrong-root) AND a genuine wrong-root boot still trips the guard (no false-negative), reproduced by a failing test first
  AC-1 pins the false-positive elimination (`;`-glued and `&&`-glued cd-into-fixture pass), written failing-first against the current parser; AC-2 pins the no-false-negative direction (off-fixture compound boot still reds, naming both roots); test plan steps the failing-first → trim → full-green sequence. Mechanism verified offline before committing the design — see the Spike section (the `;`-glued path token breaks `isUnder`; the trailing-separator trim fixes all 8 probe cases both directions).
- DONE: Confirm the change is test-harness-only (bootPathArgs in wrong_root_detect_impl_test.go) with zero product/contract impact, so the doc-diff is none
  `grep detectWrongRootBoot --include="*.go"` outside `_test.go` is empty; `bootPathArgs` lives only in the `_test.go` file (main + two worktree copies); no contract/skill/docs-site references. Recorded in the new "Documentation impact: None" section.

### Summary

Fleshed out the seed into a behavior-first spec for a pure test-harness parser fix: strip a trailing shell-separator (`;`/`&&`/`|`) from each path token `bootPathArgs` extracts, so an opus compound boot `cd <fixtureRoot>; ls; …` resolves to the clean fixture path instead of `<fixtureRoot>;` (which fails the `isUnder` check and false-positives). Spiked the riskiest path first — copies of `bootPathArgs`/`wanderTarget`/`isUnder` confirm the cause and the fix in both directions (false-positive gone, genuine wrong-root still flagged) — so no unverified mechanism remains for implementation. Scoped to the minimal per-token trim (rejected full command-splitting as unnecessary); doc-diff is none.
