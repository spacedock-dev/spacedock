---
id: 0tmv5bry1wbkww2758y88pay
title: Remove the gate validate subcommand
status: validation
source: "Captain directive, 2026-08-14: value review found near-zero value; every fault it reports also surfaces elsewhere"
started: 2026-08-15T02:55:29Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-remove-gate-validate-subcommand
issue:
sprint: durable-decisions
gates:
    version: 1
    records:
        - id: gate:0tmv5bry1wbkww2758y88pay:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:0tmv5bry1wbkww2758y88pay-backlog-1
              briefing:
                id: briefing:0tmv5bry1wbkww2758y88pay:backlog:attempt-1:revision-1
                digest: sha256:30b8fe0e5bc71046ebb679fe9cc6bab32c2bae8fc294c99e68b2307502b675a3
                request-digest: sha256:c67727e7678ea741b41d589a7b94161862beb12ae8bfa7c2554546a0c4c5157b
                room-ref: ./remove-gate-validate-subcommand/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:0tmv5bry1wbkww2758y88pay:backlog:1
                briefing: briefing:0tmv5bry1wbkww2758y88pay:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T02:53:25.738011Z"
                decision: approve
                reason: 'Captain ruling 2026-08-14 (dispatch them): approved into ideation'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:0tmv5bry1wbkww2758y88pay:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:0tmv5bry1wbkww2758y88pay-ideation-1
              briefing:
                id: briefing:0tmv5bry1wbkww2758y88pay:ideation:attempt-1:revision-1
                digest: sha256:784c8e9beef27ebdf2ef98b3f6c467b6c171ef13573b0fb36791a5da980350b4
                request-digest: sha256:890193901130940a42765c8b31f530a8e8c43821266cb69020d820fc4f9d60ef
                room-ref: ./remove-gate-validate-subcommand/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:0tmv5bry1wbkww2758y88pay:ideation:1
                briefing: briefing:0tmv5bry1wbkww2758y88pay:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T03:17:29.897072Z"
                decision: approve
                reason: Captain approved ideation 2026-08-15; implementation on opus; amendment reviewed as strict tightening
              application:
                target-stage: implementation
                state: consumed
        - id: gate:0tmv5bry1wbkww2758y88pay:validation
          stage: validation
          attempts:
            - id: gate-attempt:0tmv5bry1wbkww2758y88pay-validation-1
              briefing:
                id: briefing:0tmv5bry1wbkww2758y88pay:validation:attempt-1:revision-1
                digest: sha256:46e965126bca1650a31471ea5394e164141b00bcd664125dd2d194b2fa70da86
                request-digest: sha256:7cb9fd97f694687348d02133172eb3ababbe4ff4429b9538f5d67cf6cc42e992
                room-ref: ./remove-gate-validate-subcommand/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:0tmv5bry1wbkww2758y88pay:validation:1
                briefing: briefing:0tmv5bry1wbkww2758y88pay:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T16:07:48.648737Z"
                decision: approve
                reason: 'Captain batch approval 2026-08-15: validation PASSED; land via releng-27 train'
              application:
                target-stage: done
                state: pending
---

Remove the `gate validate` CLI subcommand. It has zero live uses across 424 recorded attempts. Fatal faults surface identically on every gate command through the shared reader. The warning class prints at state publish. The read-only digest sweep and the round check have no recorded consumer.

Scope: the two CLI branches in internal/cli/cli.go (plain and --round), the usage text, tests of the subcommand, and skill prose that names `gate validate`. Keep `ValidateRoundFile` (live consumer: `gate record --round`) and `SummaryFileDiagnosticsAt` (live consumer: `SummaryFileAt`).

## Problem

`spacedock gate validate` is a read-only diagnostic with no consumer. The 2026-08-14
usage audit found zero invocations across 424 recorded gate attempts, and the FO skill
never calls it — no `skills/` file names it at HEAD.

Each of the three faults it reports already has another verifier, confirmed against HEAD:

1. Retained-authority faults (drifted `request.json`, frozen-digest mismatch, bad
   retained binding) come from `validateRetainedAuthority`, which every gate *write*
   path calls before mutating: `prepare.go:115`, `operation.go:214` and `:279`,
   `application.go:84` and `:122`, `delivery.go:173`. A drifted room is refused
   byte-clean whether or not `gate validate` exists.
2. The unknown-application-field warning class is emitted independently by
   `status --validate` through `gateValidationDiagnostics` (`internal/status/validate.go:228`),
   with its own formatting and its own test
   (`internal/status/gate_application_warning_test.go`, which asserts an exact warning
   count of 5). **Scope matters here, and this rationale does not lean on any in-flight
   work:** for *active* entities that sweep covers the class today and continues to,
   whatever `scope-validate-warnings-to-active-entities` decides. For *archived* entities
   this removal does eliminate the last CLI path to the warning — `gate validate
   archive:<slug>` works today and is the only other printer (verified: it exits 0 and
   prints both warnings for `archive:ac2-reanchor-live-scenario-repair`, resolving via
   `internal/status/resolve.go:48`). That loss is acceptable on its own terms, not because
   another task covers it: archived scope is publish-only and terminal, so no
   tool-mediated write can ever clear an application-field extension there. The finding is
   unactionable by construction, and an alarm that can never be acted on is not a
   diagnostic worth a command surface. See Coordination.
3. The `--round` read calls `gates.ValidateRoundFile`, which `gate record --round`
   already calls on publication (`internal/cli/cli.go:369`) and prints via `printRound`.

Three corrections to the seed description, found by checking HEAD:

- The seed says the warning class "prints at state publish". It does not; the surviving
  surface is `status --validate`.
- The seed scopes "skill prose that names `gate validate`". There is none —
  `git grep "gate validate" -- skills/` is empty. (An older worktree copy of
  `skills/fo-gate-lifecycle/SKILL.md` names it; the tracked file no longer does.)
- The seed estimates ~4 files. The real surface is 8 files, because three docs and two
  extra test files carry the verb.

## Proposed approach

Delete the subcommand and the one helper that dies with it, re-point the tests that used
it as a probe for other behavior, and reconcile the three docs.

**1. Command surface** (`internal/cli/cli.go`). Remove the `validate` branch (lines
322-347) and drop `validate` from four strings: the `Use:` line (175), the `Short:` line
(176, "inspect, " goes), the published usage text (181), and the unknown-subcommand
message (185). `printRound` stays — `gate record --round` still uses it.

**2. Orphaned helper** (`internal/gates/io.go`). `gates.FormatWarning` (lines 37-42) has
exactly one caller, `cli.go:343`, inside the branch being deleted. It is gate-validate
specific by its own doc comment ("the stable operator-facing form used by gate
validation") and `status --validate` formats warnings itself. It is removed as part of
the removal, not as separate cleanup — leaving it would create exactly the dead machinery
this program exists to retire.

**3. Test re-pointing.** Five call sites use `gate validate` as a probe. Two are
load-bearing refusal assertions and were spiked (below); three are incidental read checks
inside larger flows that already assert the real outcome one line later via `gate consume`.

- `internal/ensigncycle/recorded_gate_lifecycle_test.go:411` (withdrawn retained-authority
  drift) and `:580` (forced-close validation mismatch): re-point the probe at
  `gate consume`, keeping the same assertion strings.
- `:563-571` (`validate-read` subtest, "validate read changed no workflow bytes"): delete.
  It tests only the removed subcommand.
- `internal/cli/gate_test.go:430` and `:489`: delete the `state=closed` probe lines. Both
  are immediately followed by a `gate consume` assertion over the same authority path.
- `internal/cli/gate_test.go:227`: delete the `--round` validate assertion; the
  `gate record --round` assertion three lines above already covers the round summary.
- `internal/cli/gate_application_warning_test.go`: delete
  `TestGateValidateReportsApplicationExtensionWarnings` (the warning class keeps its
  verifier in `internal/status/gate_application_warning_test.go`); re-point
  `TestGateValidateRejectsBadRetainedRequestBinding` at `gate consume`, which refuses on
  the same `validateRetainedAuthority` call.

**4. Usage-error assertion.** `internal/cli/gate_test.go:316`
(`TestRemovedGateVerbsAreAbsentAndSideEffectFree`) is the established venue: it already
covers the previously-removed `review` and `eligibility` verbs, asserting exit 2, the
unknown-subcommand stderr, and an unchanged working directory. Add `"validate"` to its
verb table and update the expected message. This adds no new mechanism.

*Simplest alternative considered:* rely on the exact-text help fixture alone
(`gate_test.go:112`). Insufficient — the help fixture proves the text changed, not that
the grammar rejects the verb. A surviving branch would still accept `gate validate` while
the help omitted it. The usage-error assertion is what makes AC-1 falsifiable.

**5. Help fixture.** `internal/cli/gate_test.go:112` is a deliberate exact-text contract
("Grammar additions/removals must change this fixture openly"). Drop its
`gate validate` line.

### Spike: do the two refusal probes survive re-pointing?

The riskiest unverified claim was that `gate consume` refuses identically where
`gate validate` did. `ConsumeAt` calls `validateRetainedAuthority` immediately after
`Read`, before any status or eligibility logic (`internal/gates/application.go`), so it
should refuse first — but that needed exercising, not reading.

A throwaway test in `internal/ensigncycle` replayed both scenarios end-to-end against a
freshly built binary with the probe swapped to `gate consume`. Both passed:

- Withdrawn retained-authority drift: `exit=1`, stderr
  `Error: attempt gate-attempt:recorded-gate-task-validation-1 retained request.json does not match its frozen digest`
  — contains `frozen digest`, entity bytes unchanged.
- Forced-close validation mismatch: `exit=1`, stderr
  `Error: attempt gate-attempt:recorded-gate-task-validation-1: only approve resolutions may carry an application`
  — contains `application`, byte-clean, no `.gates.lock` residue.

Both existing assertion strings hold unchanged. The spike file was deleted; the
implementation's first tests are the two re-pointed probes.

### Doc diff

`docs/site/reference/command-reference.md` — delete both `gate validate` rows (96, 97).

`docs/site/reference/frontmatter-contract.md:13`

> Before: … those keys produce warnings only on explicit `status --validate` or `gate validate`, are ignored for authority, and are never written.
>
> After: … those keys produce warnings only on explicit `status --validate`, are ignored for authority, and are never written.

`docs/specs/gate-resolution-frontmatter-contract.md:89-90`

> Before: … reports them as warnings on explicit `status --validate` or `gate validate`, ignores them for authority, and never writes them.
>
> After: … reports them as warnings on explicit `status --validate`, ignores them for authority, and never writes them.

`docs/specs/gate-resolution-frontmatter-contract.md:193-194`

> Before: … the recorder preserves that body byte-for-byte. `spacedock gate validate <entity> --round STAGE/CYCLE` replays the pointer and reports every Resolution as advisory structural evidence.
>
> After: … the recorder preserves that body byte-for-byte. The published round is the durable evidence; `gate record --round` reports every Resolution as advisory structural evidence on publication.

`docs/specs/gate-resolution-frontmatter-contract.md:254` — delete the line
`spacedock gate validate ENTITY [--workflow-dir DIR]` from the command-surface block.

`docs/specs/gate-resolution-frontmatter-contract.md:260`

> Before: … callers cannot submit an operation envelope or candidate identities. `gate validate` is read-only and reports the selected record's last attempt.
>
> After: … callers cannot submit an operation envelope or candidate identities.

## Keep-boundary

`ValidateRoundFile` **stays** — confirmed live at HEAD: `internal/cli/cli.go:369`
(`gate record --round`).

`SummaryFileDiagnosticsAt` **stays** per the approved boundary, but the boundary's stated
reason does not hold at HEAD and the gate should see this. Its production caller today is
the `gate validate` branch (`cli.go:337`). The other caller, `SummaryFileAt`
(`io.go:199`), is itself reached only from `SummaryFile` (`io.go:195`) and from
`internal/gates/prepare_test.go:956`; `SummaryFile`'s only caller is
`internal/gates/application_test.go:365`. So after this removal the whole
`SummaryFile` / `SummaryFileAt` / `SummaryFileDiagnosticsAt` trio is production-unreachable
— reached from tests only. Two further points, both confirmed independently by the
`scope-validate-warnings-to-active-entities` ideation (see Coordination): nothing prints
the warnings slice `SummaryFileDiagnosticsAt` returns, since `SummaryFileAt` discards it
with `_`; and the `newSummaryFile` hits in `cmd/spacedock-release/e2e_gate_test.go` are an
unrelated local helper, not this function, so a naive `git grep SummaryFile` over-reports
consumers (as do the stale `.worktrees/` copies).

**Why it stays anyway: scope discipline, not liveness.** This is the honest answer if the
gate asks why the keep-boundary survives the evidence above. The trio is a general read
surface rather than gate-validate machinery, its retirement is a separate value-chain
question, and removing it here would mean re-pointing two `internal/gates` tests
(`prepare_test.go:956`, `application_test.go:365`) that currently prove retained-authority
refusal through it. Neither this task's value AC nor the peer's needs it gone. Recommend
filing it as its own audit entity — but the keep-boundary should not be read as a claim
that the reader is still used, because after this cut it is not.

One blast-radius claim raised in coordination does **not** hold, checked rather than
accepted: retiring the trio would not orphan `nearestWorkflowDir`. That helper has five
other callers — `application.go:72` and `:109`, `operation.go:212`, `:277`, and `:345` —
so `io.go:195` is one of six. The reasons to defer are the two above, not this one.

The doc comment on `SummaryFileDiagnosticsAt` ("is the gate-validate read surface") does
become stale and is corrected in place, since it names a command that no longer exists.

## Coordination

Overlap found with `scope-validate-warnings-to-active-entities` (in flight, uncommitted at
`internal/status/validate.go` and `internal/status/field_conformance.go`). That task narrows
the unknown-gate-application-field warn channel to active-scope entities — the same channel
this task relies on as the surviving verifier for AC-4.

**Rationale independence (raised by that entity's ideation, 2026-08-14, and agreed).** The
two removals were at risk of justifying each other in a loop: that task's draft cited
`gate validate archive:<slug>` as the surviving on-demand surface for archived entities,
while this task cited the `status --validate` sweep. Each entity was pointing at the
surface the other removes. Verified against HEAD rather than taken on report:
`gate validate archive:ac2-reanchor-live-scenario-repair --workflow-dir docs/dev` exits 0
and prints both warnings today; `internal/cli/cli.go:337` is the sole non-test consumer of
`gates.SummaryFileDiagnosticsAt`; and there is no `publish` command (`state` is
`init|new|ready|sweep|commit`). So if both land as scoped, the application-extension
warning becomes unreachable through any CLI for archived entities.

Both entities now stand on unactionability alone for the archived case, and neither cites
the other's surface. This task's rationale is unaffected in substance: its load-bearing
grounds — zero usage across 424 attempts, retained-authority faults refused by every gate
write path, and the round check duplicated by `gate record --round` — are all independent
of that task. Only the warning-class point needed its scope stated precisely, which is now
done in Problem item 2 and declared as a semantic change in Expected surface.

No textual collision (that task edits `internal/status/`; this one edits `internal/cli/`,
`internal/gates/io.go`, three test files, and three docs). The AC-4 verifier still holds —
`gateApplicationWarningFixture` writes an active-scope entity at the workflow root, so the
new `e.scope != "active"` skip does not suppress it. No sequencing constraint in either
direction; either may land first.

Also checked, no overlap: `remove-startup-capability-probe` (skills/, install docs) and
`retire-requires-contract-sentinel`, which explicitly lists "`gate validate` demotion" as
out of its scope, deferring to this entity.

## Out of scope

The rest of the gate decision surface. `ValidateRoundFile` and the shared readers.
Retiring `SummaryFile` / `SummaryFileAt` / `SummaryFileDiagnosticsAt` (see Keep-boundary).
`gates.ReadWithWarnings` and the `gates.Diagnostic` alias, which already have zero callers
at HEAD — pre-existing dead code this change neither creates nor touches.

## Expected surface and tolerance

Estimate net LOC change: **-75, tolerance ±30**, across **8 files**:
`internal/cli/cli.go`, `internal/gates/io.go`, `internal/cli/gate_test.go`,
`internal/cli/gate_application_warning_test.go`,
`internal/ensigncycle/recorded_gate_lifecycle_test.go`,
`docs/site/reference/command-reference.md`, `docs/site/reference/frontmatter-contract.md`,
`docs/specs/gate-resolution-frontmatter-contract.md`.

Roughly 81 deletions against roughly 10 insertions (re-worded doc lines, one verb-table
entry, two re-pointed probe lines).

**Observable semantics changed:**

- **Command grammar (changed).** `spacedock gate validate <entity>` moves from a
  supported read-only subcommand (exit 0 on a clean read, 1 on a fault) to an unknown
  subcommand (exit 2, usage error, no side effects). `gate validate --round` goes with it.
- **Published help text (changed).** The `gate --help` exact-text contract loses one
  grammar line; `Use:` and `Short:` lose `validate` / `inspect,`.
- **Diagnostics reachability (changed, no enforcement lost).** The unknown-application-field
  warning class is reachable only through `status --validate` afterward. For active
  entities that is unchanged coverage. For archived entities this removes the last CLI
  printer: `gate validate archive:<slug>` is the only surface that reports the class for
  archived scope today, and `gates.SummaryFileDiagnosticsAt` — the reader behind it — is
  left with no printer for its warnings slice at all. Declared deliberately, on the ground
  that archived scope is publish-only and terminal so the finding is unactionable.
  Retained-authority faults keep failing every gate write path closed; only the read-only
  report of them goes.
- **Unchanged:** stored formats (`gates` frontmatter, room bytes, round pointers),
  authority and eligibility rules, exit codes of every surviving command, and all runtime
  or dispatch behavior.

## Acceptance criteria

**AC-1 - The command grammar no longer carries the subcommand.**
`spacedock gate validate <entity>` exits 2 with the unknown-subcommand usage error and
changes nothing on disk, and the published `gate --help` contract no longer lists it.
Verified by: `"validate"` added to the verb table in
`TestRemovedGateVerbsAreAbsentAndSideEffectFree` (`internal/cli/gate_test.go:316`), which
asserts exit 2, the unknown-subcommand stderr, and an unchanged working directory; plus
the exact-text help fixture at `gate_test.go:112`. Fails if the branch survives, if the
verb is silently accepted, or if the help line is left behind.

**AC-2 - The change removes more lines than it adds.**
Cumulative line delta of the delivery diff against `origin/main` is negative.
Verified by: `git diff --shortstat origin/main` shows deletions exceeding insertions. This
is the value measure — a "removal" that nets positive has not removed anything.

**AC-3 - The removal leaves no orphaned machinery behind it.**
No symbol loses its last production consumer without being removed or named in the
Keep-boundary section above.
Verified by: `git grep FormatWarning` returns no hits outside history; the only symbols
left consumer-thin are the three named in Keep-boundary. Fails if the validate branch is
deleted while `FormatWarning` is left dangling.

**AC-4 - Every behavior the subcommand was a probe for keeps a verifier.**
Retained-authority drift and forced-close mismatch are still refused byte-clean with the
same operator-facing strings, and the unknown-application-field warning class is still
asserted.
Verified by: the two re-pointed `gate consume` probes in
`internal/ensigncycle/recorded_gate_lifecycle_test.go` still assert `frozen digest` and
`application` with entity bytes unchanged and no lock residue (spiked green, see above);
`internal/status/gate_application_warning_test.go` still asserts the exact warning count
of 5 through `status --validate`. Fails if a refusal degrades to exit 0, changes message,
or mutates the entity.

**AC-5 - The suite stays green.**
Verified by: `go test ./...` and `go test ./... -race` pass.

## Test plan

Deletion-driven, no new harness. Cost: low — one afternoon, no new fixtures.

- **Re-pointed first**, before the deletion, so the coverage move is proven independently
  of the removal: the two `internal/ensigncycle` refusal probes swapped to `gate consume`.
  Already exercised end-to-end against a freshly built binary in the ideation spike; both
  refusals reproduce with the existing assertion strings.
- **Then the deletion**, with the verb-table entry and help fixture updated in the same
  commit so AC-1 fails loudly if either half is forgotten.
- **Existing suite as the regression net.** `go test ./...` and `-race`. The
  `internal/ensigncycle` recorded-gate tests build and drive the real binary, so the
  command-grammar change is proven at the CLI boundary, not just in unit scope.
- **No live workflow test needed.** Nothing here touches runtime dispatch or state sync;
  the claim is command-grammar and test-coverage shape, both offline-checkable.

## Stage Report: ideation

- DONE: Keep-boundary confirmed: ValidateRoundFile (gate record --round) and SummaryFileDiagnosticsAt (SummaryFileAt) stay
  `ValidateRoundFile` confirmed live at `internal/cli/cli.go:369`. `SummaryFileDiagnosticsAt`
  kept as instructed, but confirmation found its stated reason does not hold: `SummaryFileAt`
  reaches only tests, so the trio becomes test-only after removal — recorded in Keep-boundary
  for the gate to rule on rather than silently expanded into.
- DONE: Usage-error assertion covers the removed subcommand; help text and skill prose reconciled
  Assertion venue is the existing `TestRemovedGateVerbsAreAbsentAndSideEffectFree`
  (`internal/cli/gate_test.go:316`), which already covers removed `review`/`eligibility`;
  adding `"validate"` needs no new mechanism. Help text is the exact-text fixture at
  `gate_test.go:112` plus four strings in `cli.go`. Skill prose needs no change —
  `git grep "gate validate" -- skills/` is empty at HEAD, correcting the seed's scope.
- DONE: Command-grammar semantic change declared in the estimate
  "Expected surface and tolerance" declares the grammar change (supported read-only verb →
  unknown subcommand, exit 2), the help-text change, and the diagnostics-reachability
  change, plus what is explicitly unchanged (stored formats, authority, exit codes,
  runtime behavior). Estimate: -75 LOC ±30 across 8 files, correcting the seed's ~4.

### Summary

Confirmed the removal scope against HEAD and found the seed wrong on three points: the
warning class survives at `status --validate` (not "state publish"), no skill prose names
the verb, and the surface is 8 files rather than 4. The one genuine consequence beyond
deletion is that `gates.FormatWarning` loses its only caller, so it is removed with the
branch. Spiked the riskiest claim — that the two load-bearing refusal probes keep refusing
when re-pointed at `gate consume` — end-to-end against a freshly built binary; both
reproduce their existing assertion strings byte-clean, so no coverage is lost. Flagged for
the gate: the approved keep-boundary for `SummaryFileDiagnosticsAt` rests on a chain that
dead-ends in tests, which I kept rather than expanded, and recommend as its own entity.

Post-report coordination with `scope-validate-warnings-to-active-entities` removed a
circular justification between the two entities — each had cited the surface the other
removes. Verified the peer's claims against HEAD (`gate validate archive:<slug>` works and
is the only archived-scope printer; no `publish` command exists) and rewrote this task's
warning-class rationale to stand on unactionability of archived findings alone. The removal
plan and estimate are unchanged; only the rationale's precision and one declared semantic
change moved. See Coordination.

## Stage Report: implementation

- DONE: Re-point the two refusal probes to gate consume BEFORE the deletion, per the approved test plan
  Commit `18f851eb8` (re-point) precedes `d1d23b68a` (deletion), and both probes were run green
  while `gate validate` still existed — so the coverage move is proven against the surviving
  command rather than by the removal. `TestRecordedGateLifecycleWithdrawColdBootReplaceAndConsume`
  asserts exit != 0, stderr containing `frozen digest`, and entity bytes unchanged;
  `AC5RefusalMatrix/forced-close-validation-mismatch` asserts exit != 0, output containing
  `application`, and no `.gates.lock` residue. Either fails if `gate consume` refused later than
  `validateRetainedAuthority`, refused for a different reason (an eligibility message would miss
  both substrings), or exited 0.
- DONE: Deletion lands within -75 LOC plus or minus 30 across the 8 declared files; verb-table entry and help fixture change in the same commit
  `git diff --shortstat` against the merge-base: 8 files changed, 20 insertions, 102 deletions —
  net **-82**, inside the declared -75 +/- 30, across exactly the 8 declared files and no others.
  AC-2 (deletions exceed insertions) holds. The verb-table entry and the exact-text help fixture
  are both in `d1d23b68a` with the branch deletion.
  AC-1 exercised on a freshly built binary, not just in unit scope: `gate validate task` and
  `gate validate task --round ideation/1` each exit 2 with
  `unknown subcommand (want: prepare|withdraw|record|consume)`, the working directory is
  byte-identical by shasum before/after, and `gate --help` no longer prints the grammar line.
  `TestRemovedGateVerbsAreAbsentAndSideEffectFree` now covers `validate` beside `review` and
  `eligibility`; it fails if the branch survives or the verb is silently accepted. The help
  fixture is a two-sided exact-text contract — changing only `cli.go` or only the fixture fails it.
- DONE: Suite green including -race; FormatWarning goes with its last caller; keep-boundary trio untouched
  `go test ./...` and `go test ./... -race` are green in every package except one pre-existing
  environmental failure, `TestCodexResolveManifestAgainstInstalledHost`: the real `codex` CLI
  cannot read `/Users/clkao/.codex/config.toml` ("Operation not permitted") under this session's
  sandbox. Verified not mine — the identical failure reproduces on the unmodified merge-base
  `4d1912a69` in a throwaway worktree. No data races reported.
  AC-3: `git grep FormatWarning` returns nothing outside history; it was removed in the same
  commit as its only caller. `printRound` and `ValidateRoundFile` stay live at `cli.go:343-344`
  (`gate record --round`), so the keep-boundary for `ValidateRoundFile` is intact.
  AC-4: `internal/status/gate_application_warning_test.go` still asserts all five named extension
  fields and an exact warning count of 5 through `status --validate`; it fails if the surviving
  warning channel drops or double-counts the class.
  Keep-boundary trio untouched as approved: `SummaryFile` / `SummaryFileAt` /
  `SummaryFileDiagnosticsAt` still compile and keep their two test consumers
  (`prepare_test.go:956`, `application_test.go:365`). Confirmed they are now production-unreachable,
  exactly as ideation predicted; only the stale doc comment naming `gate validate` was corrected.

### Summary

Removed the `gate validate` subcommand in two commits: the probe re-pointing first, proven green
against the still-present command, then the deletion with the verb-table entry and help fixture in
one commit. Surface came in at -82 LOC across exactly the 8 declared files, inside the -75 +/- 30
tolerance. `gates.FormatWarning` went with its last caller.

Two judgment calls stayed inside the approved surface rather than expanding it, both flagged for the
gate. `TestGateValidateRejectsBadRetainedRequestBinding` was renamed to
`TestGateConsumeRejectsBadRetainedRequestBinding` since its name otherwise cited a removed command;
its file is still `gate_application_warning_test.go` though it no longer holds a warning test, and
`TestGateRequestLocatorCarriesArbitraryBriefingNameThroughRecordValidateAndEligibility` keeps a name
citing two now-removed verbs. Both are naming rot, not dead machinery, so renaming them is churn
outside this task's declared 8 files.

The one red test in the run is environmental and pre-existing, reproduced on unmodified `main`
before being attributed here.

## Stage Report: validation

- DONE: Independently re-exercise every AC against the worktree binary and merge-base, never by reading the report: usage-error exit 2 both arms, shortstat within -75 plus or minus 30 across the 8 declared files, FormatWarning grep empty, warning-count-5 test, race suite
  AC-1: freshly built worktree binary, `gate validate task` and `gate validate task --round ideation/1`
  both exit 2 with `unknown subcommand (want: prepare|withdraw|record|consume)`, zero stdout, no files
  created in an empty dir, and sha256-identical bytes across a full copy of a real workflow tree;
  `gate --help` contains zero `validate` lines. Mutation check in a throwaway checkout: restoring the
  merge-base grammar makes `TestRemovedGateVerbsAreAbsentAndSideEffectFree/validate` fail (exit 1 — the
  branch accepts the verb) and fails the exact-text help fixture, so both AC-1 tests are two-sided;
  restoring merge-base `cli.go` alone breaks the build on `undefined: gates.FormatWarning`, proving the
  helper died with its last caller. AC-2: `git diff --shortstat 4d1912a69..HEAD` = 8 files, +20/-102,
  net -82, inside -75±30; numstat file list matches the 8 declared files exactly. AC-3:
  `git grep FormatWarning` empty; `ValidateRoundFile` live at `internal/cli/cli.go:343`; keep-boundary
  trio untouched. AC-4: `TestRecordedGateLifecycleWithdrawColdBootReplaceAndConsume` PASS solo (asserts
  exit!=0, `frozen digest`, entity bytes unchanged); `AC5RefusalMatrix/forced-close-validation-mismatch`
  PASS (asserts `application`, byte-clean tree digest, no `.gates.lock` residue);
  `TestStatusValidateReportsGateApplicationExtensionsAsWarnings` PASS (all five named fields, exact
  count 5, ordinary status silent); `TestGateConsumeRejectsBadRetainedRequestBinding` PASS. AC-5: plain
  suite green in every package except the environmental Codex test; `go test ./... -race -timeout 25m`
  green with the same single exception. No `gate validate` mention survives in any tracked file; the
  three doc diffs match the approved doc-diff plan line for line; gofmt clean.
- DONE: Calibrate suite verdicts to the merge-base baseline: reproduce any failure on unmodified 4d1912a69 in a clean control worktree before attributing it; TestCodexResolveManifestAgainstInstalledHost is known environmental
  Reproduced the Codex failure myself on a clean control worktree at unmodified `4d1912a69`: identical
  `Failed to read config file /Users/clkao/.codex/config.toml: Operation not permitted`. The first
  contended full run also hit 10m default package timeouts in `internal/cli` and `internal/ensigncycle`
  (several sibling ensign suites were running); solo reruns with `-timeout 25m` finished green in 814s
  (cli, codex-only failure — the mid-panic state-commit test passes) and 790s (ensigncycle, fully ok),
  so both timeouts are load, not the candidate.
- DONE: Verdict PASSED or REJECTED with per-AC evidence citations; the flagged naming-rot items are out of scope, report-only
  Verdict: PASSED. No material findings. Polish (report-only, out of scope per checklist): naming rot in
  `TestGateRoundRecordAndValidateCLI` and
  `TestGateRequestLocatorCarriesArbitraryBriefingNameThroughRecordValidateAndEligibility` (names cite
  removed verbs), and `TestGateConsumeRejectsBadRetainedRequestBinding` living in
  `gate_application_warning_test.go` which no longer holds a warning test. Deferred risk: `internal/cli`
  and `internal/ensigncycle` exceed go test's default 10m package timeout under concurrent ensign load
  on this machine; promotes to material if CI runs suites on similarly shared hardware without a raised
  `-timeout`.

### Summary

Every AC re-exercised independently against the worktree binary and the merge-base, including a
mutation check proving the AC-1 tests fail when the old grammar is restored. The diff is -82 LOC across
exactly the 8 declared files, all re-pointed probes and the warning-count-5 verifier pass, and plain
plus race suites are green with the single calibrated-environmental Codex exception reproduced on
unmodified `4d1912a69`. Verdict PASSED; only polish-level naming rot and one environment-shaped
deferred risk noted.
