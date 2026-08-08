---
title: "Legacy verdict tokens flood every state-checkout validation commit"
status: validation
source: "Durable-decisions Commander dogfood, 2026-07-22: a normal `spacedock state commit` ran the state-checkout pre-commit validator and printed 117 pre-existing verdict-enum warnings before succeeding, burying any new warning attributable to the entity being committed. Dedupe found the schema validator and pre-commit-hook tasks, but no entity owns legacy-token migration or bounded warning output."
sprint:
id: mr9k7c0g35jhrrdv4zqyjctw
started: 2026-08-08T00:16:07Z
worktree: .worktrees/spacedock-ensign-legacy-verdict-warning-flood
pr: "#634"
---

One path-scoped state mutation currently emits 117 `Warning: field 'verdict' ... is not one of [PASSED REJECTED]` lines from the checkout-wide pre-commit validation: 104 lowercase `passed`, 9 lowercase `rejected`, and four other legacy/superseded values. The command exits successfully, but the historical backlog swamps the warning signal for the current change.

## Problem

The schema validator correctly treats conventional-field violations as warnings, and the pre-commit hook correctly runs checkout-wide validation. The development state checkout predates the canonical verdict enum, however, so every state commit repeats a large fixed corpus of legacy warnings. A newly introduced warning is difficult to distinguish in the flood, and routine command output becomes disproportionately noisy.

## Boundary for ideation

Determine the smallest safe ownership boundary between a one-time legacy data migration and bounded/delta-aware pre-commit diagnostics. Preserve schema-driven validation, warning severity, historical decision meaning, path-scoped state commits, and hard-error blocking. Do not silently suppress novel warnings or rewrite archived prose bodies; only frontmatter tokens whose semantics can be proven equivalent are candidates for migration.

## Acceptance sketch

- A clean state mutation produces bounded, actionable validation output instead of replaying the current 117-line legacy corpus.
- A newly introduced invalid verdict still appears and remains distinguishable; structural validation errors still block the commit.
- Any migrated `passed`/`rejected` tokens preserve their semantic value as canonical `PASSED`/`REJECTED`, with an auditable count and no unrelated entity-body changes.

## Evidence

`spacedock status --workflow-dir docs/dev --validate` exits 0 and currently emits 117 stderr lines, all verdict-enum warnings: 104 contain `value "passed"`, 9 contain `value "rejected"`, and the remainder are other legacy/superseded tokens. The installed pre-commit hook captures combined validation output and echoes it wholesale on every state-checkout commit.

## Out of scope

- Changing the durable-decisions sprint or its release criteria.
- Weakening warning-tier schema conformance or pre-commit hard-error enforcement.
- Bulk rewriting entity bodies, reports, review rooms, or unrelated frontmatter.

## Stage Report: implementation

- DONE: Review PR #634 with Roborev before changing the imported candidate, and preserve every finding for workflow disposition.
  Roborev job 1037 reviewed untouched head `6c45fd59`; its exact finding and captain-decline comment remain available via `roborev show --job 1037`, with advertised log path `/Users/clkao/.roborev/logs/jobs/1037.log`.
- DONE: Prove that canonical writes plus legacy-compatible reads remove the verdict-warning flood while novel invalid values remain visible.
  Focused status tests fail if writes stop storing `PASSED`/`REJECTED`, lowercase legacy reads warn, novel `needs-work` is rewritten/silenced, or structural errors stop exiting 1; the live checkout now emits 0 lowercase legacy verdict warnings and exactly four `superseded` warnings instead of 117 verdict warnings.
- DONE: Deliver the adopted candidate with focused, full, race, formatting, and exact surface evidence, or stop unchanged on an unresolved material finding.
  Adopted PR commit `6c45fd59c7377eadfb2c2013d048bb77fa004c69` is clean; `gofmt -w ./cmd ./internal`, focused tests, `go test ./...`, and `go test ./... -race` passed, and the PR surface is 16 files, 234 insertions, 31 deletions.

### Review-finding disposition

- Reviewer observation (Roborev job 1037, exact): `Severity: Medium`; `Location: internal/status/mutate.go:154`; `Problem: Canonicalizing before change reporting alters the stable status --set stdout from the caller-provided passed/rejected to PASSED/REJECTED, as shown by the modified golden fixtures. This unnecessary public-output regression can break consumers that parse the existing output.` `Fix: Canonicalize only the persisted frontmatter value while retaining the requested spelling for CLI change reporting, and keep the existing lowercase golden output.`
- Released user and normal workflow: `spacedock status --set <slug> ... verdict=passed` is the normal mutation path and reports its resolved change.
- Observable harm proposed by the worker: three stdout goldens changed from lowercase to uppercase, potentially affecting exact-output consumers.
- Affected boundary proposed by the worker: `contract[AGENTS.md#Priorities]` requires stable, fixture-tested command output.
- Trigger evidence: `handlers.go` reports the value returned from `updateFrontmatter`; canonicalization therefore appears in text, quiet, and JSON narration as well as storage.
- Worker proposal: Material, task-owned, fix; initial FO authorization: fix.
- Final authorization and outcome: `captain-ruling[2026-08-07]: status --set output need not remain lowercase; canonical stored-value narration is accepted.` Disposition DECLINED; the in-progress correction was removed and byte equality to imported head was verified before further checks.
- Rerun: Roborev job 1040 reports `No issues found`; canonical artifact is `roborev show --job 1040`, with advertised log path `/Users/clkao/.roborev/logs/jobs/1040.log`.

### Verification notes

- The first live-checkout full/race attempts exposed an unrelated stale pilot manifest whose two named entities have since moved under `_archive/`; both exact suites passed with `SPACEDOCK_STATE_ROOT` pinned to immutable state commit `a0169cc2d8a5e4912ed33f75ca8422a767e71c9e`, where the manifest paths exist.
- Live candidate validation exited 0 with `VALID`; verdict warnings were bounded to the four deliberately novel `superseded` tokens, while all 104 `passed` and 9 `rejected` legacy tokens were accepted case-compatibly.

### Summary

Adopted GitHub PR #634 unchanged: conventional verdict writes use the schema spelling, legacy case variants read compatibly, and unconventional values retain warning visibility without a bulk state migration. The only Roborev finding was preserved and declined by explicit captain ruling; the unchanged candidate then passed focused, full, race, formatting, live-state, and final Roborev checks.

## Stage Report: validation

- DONE: Independently verify that canonical verdict writes and legacy-compatible reads remove lowercase warning noise while genuinely unconventional values still warn.
  Focused status/CLI tests fail if either writer stores noncanonical case, legacy `passed`/`rejected` warns, `needs-work` is rewritten, unconventional verdicts stop warning, or structural errors stop exiting 1; live validation exited 0/`VALID` with 0 lowercase and 4 `superseded` verdict warnings.
- DONE: Confirm the unchanged PR #634 candidate satisfies the task value without bulk state migration, and audit both Roborev dispositions against the captain ruling.
  HEAD remained exact candidate `6c45fd59c7377eadfb2c2013d048bb77fa004c69` with a clean worktree; job 1037's stdout finding is DECLINED by `captain-ruling[2026-08-07]`, and rerun job 1040 reports no issues.
- DONE: Reproduce the proportionate focused evidence, inspect the exact diff and public output change, and recommend PASSED or REJECTED with any residual risk.
  PASSED: the 16-file/234-addition/31-deletion diff, uppercase text goldens, `git diff --check`, formatting, focused tests, and pinned-state race tests for status/CLI/gates passed; an unpinned whole-repo race only reproduced the documented unrelated stale pilot-manifest drift.

### Summary

Recommend PASSED with no material candidate finding: canonical writes and case-compatible reads remove the 113 lowercase verdict warnings without migrating state, while four genuinely unconventional values remain visible. Residual operational risk is adjacent rather than candidate-owned: the current live checkout also emits 125 unknown gate-application-field warnings, so total pre-commit output remains noisy even though the legacy verdict corpus is gone; revisit if bounded output is broadened beyond verdict diagnostics.
