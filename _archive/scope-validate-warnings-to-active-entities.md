---
id: f0zn4sr7nz7xsxmyw6aw6bsm
title: Scope validate warn channels to active entities
status: done
source: "0.27 cut audit (2026-08-14), adversarially verified; captain directed filing"
started: 2026-08-15T02:55:41Z
completed: 2026-08-15T19:56:05Z
verdict: PASSED
score:
worktree: .worktrees/spacedock-ensign-scope-validate-warnings-to-active-entities
issue:
gates:
    version: 1
    records:
        - id: gate:f0zn4sr7nz7xsxmyw6aw6bsm:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:f0zn4sr7nz7xsxmyw6aw6bsm-backlog-1
              briefing:
                id: briefing:f0zn4sr7nz7xsxmyw6aw6bsm:backlog:attempt-1:revision-1
                digest: sha256:fb33a5d317815de37d2dd3915e99879cee89c15c0c8e151f382e99e8ab1d4e31
                request-digest: sha256:ae895696854c9670e468ff95993a9192b1c7136794226febf95723e25e7b7ad6
                room-ref: ./scope-validate-warnings-to-active-entities/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:f0zn4sr7nz7xsxmyw6aw6bsm:backlog:1
                briefing: briefing:f0zn4sr7nz7xsxmyw6aw6bsm:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T02:53:49.94516Z"
                decision: approve
                reason: 'Captain ruling 2026-08-14 (dispatch them): approved into ideation'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:f0zn4sr7nz7xsxmyw6aw6bsm:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:f0zn4sr7nz7xsxmyw6aw6bsm-ideation-1
              briefing:
                id: briefing:f0zn4sr7nz7xsxmyw6aw6bsm:ideation:attempt-1:revision-1
                digest: sha256:3e809f22e80af54aa94b0e207d3fa13a7728ac23492fbf0fd8792a9b8b641faf
                request-digest: sha256:968622d940e98e4aeab52ee72d9cae24a193e8093827257a331dde170ca53a23
                room-ref: ./scope-validate-warnings-to-active-entities/review/ideation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-15T03:23:23.90906Z"
                reason: Entity amended after prepare (d3590a590, citation and honesty fix); withdrawing stale attempt to re-prepare against current bytes
            - id: gate-attempt:f0zn4sr7nz7xsxmyw6aw6bsm-ideation-2
              briefing:
                id: briefing:f0zn4sr7nz7xsxmyw6aw6bsm:ideation:attempt-2:revision-1
                digest: sha256:1248c71674b60a0e6ef2885ea0e84be72207ad3ebd49c73e1921f85c3925785b
                request-digest: sha256:d3703fc33ace4b63951ea1f2e708635ffc763680d1cd0856012531659c80cd0d
                room-ref: ./scope-validate-warnings-to-active-entities/review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:f0zn4sr7nz7xsxmyw6aw6bsm:ideation:2
                briefing: briefing:f0zn4sr7nz7xsxmyw6aw6bsm:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-15T03:54:51.36043Z"
                decision: approve
                reason: 'Captain ruling 2026-08-15 (approve all except x8): approved into implementation'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:f0zn4sr7nz7xsxmyw6aw6bsm:validation
          stage: validation
          attempts:
            - id: gate-attempt:f0zn4sr7nz7xsxmyw6aw6bsm-validation-1
              briefing:
                id: briefing:f0zn4sr7nz7xsxmyw6aw6bsm:validation:attempt-1:revision-1
                digest: sha256:72d2b0a83c2ea97a1046cf76fe6ff17fdd0b3de522aa6566ae897e56dee41b71
                request-digest: sha256:e4114c1a92548d5dd4b6480a03bab228549251ad667ccbedc49fd6b17002f7da
                room-ref: ./scope-validate-warnings-to-active-entities/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:f0zn4sr7nz7xsxmyw6aw6bsm:validation:1
                briefing: briefing:f0zn4sr7nz7xsxmyw6aw6bsm:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T16:08:52.025657Z"
                decision: approve
                reason: 'Captain batch approval 2026-08-15: validation PASSED; land via releng-27 train'
              application:
                target-stage: done
                state: consumed
pr: pr-merge:700
archived: 2026-08-15T19:56:05Z
---

Stop emitting warn-tier findings for archived entities in `status --validate`. Today 125 of 126 report lines are archived-scope warnings (121 unknown-gate-application-field, 4 verdict-enum), the report still ends VALID, and archived scope is publish-only, so no tool-mediated fix can ever silence them. The alarm fires identically forever and carries no information. The pre-commit hook echoes the full report on every state commit, so every commit dumps 51KB of noise.

Scope: filter the two warn channels (internal/status/validate.go:228-244 gate diagnostics; internal/status/field_conformance.go:88-108 enum conformance) to active entities. Keep structural and ID validation scope-inclusive. Keep the read tolerance itself - it is load-bearing for every gates read over the legacy corpus. Precedent: 6c45fd59c fixed the identical pathology for verdict case with the same rationale.

## Problem

`status --validate` sweeps active *and* archived entities through two warn-tier channels. Measured at HEAD (4d1912a69) against `docs/dev`:

- 126 report lines total, of which 125 are `Warning:` lines and every one carries `scope=archived`.
- 121 are `unknown gate application field` (validate.go:238-241), and only three distinct keys account for all of them: 60 `action`, 60 `blockers`, 1 `execution-hold`.
- 4 are `field 'verdict' value "superseded" is not one of [PASSED REJECTED]` (field_conformance.go:100-104) — every one the same token.
- The run exits 0 and stdout is `VALID`. Captured output is 51457 bytes.
- Active-scope warn lines: zero.

The value chain is broken at the fix step. Archived scope is publish-only and terminal, so no tool-mediated write can ever clear these findings — the alarm fires identically on every invocation, forever, and a signal that cannot change carries no information.

The seed described the 4 verdict warnings as the same pathology `6c45fd59c` fixed. Checking rather than assuming corrected that: they are all `verdict: superseded`, a token the schema's conventional set does not declare at all, not the case variance that commit addressed. `6c45fd59c` is nonetheless the right precedent, and the corpus shows why — the state checkout holds 178 entities with lowercase `verdict: passed`/`rejected` and *none* of them warn, because that commit made the enum check case-insensitive precisely so archived terminal entities would not have to be hand-edited to silence a diagnostic about a value the tool itself wrote. This task applies the same ruling one level up: the 121 application-field keys are legacy extensions the current writer never emits and the read tolerance deliberately accepts, and the 4 `superseded` verdicts sit on terminal entities. Neither is fixable where it is reported.

Whether `superseded` belongs in the verdict enum is a real question this task deliberately does not answer — see Out of scope.

The cost is concrete and recurring. The state-checkout pre-commit hook (`.git/hooks/pre-commit:35-37`) captures the report with `2>&1` and echoes it to stderr on every state commit, so each commit dumps all 51457 bytes. During a multi-ensign program that is once per entity write, per ensign. The one line an operator needs — `VALID`, or a real `Error:` — is buried under 125 lines that will never differ.

## Proposed approach

Skip non-active entities at the two warn emission points. Structural and ID validation stay scope-inclusive.

`internal/status/validate.go`, inside `gateValidationDiagnostics`, after the `err != nil` branch and before the compatibility loop:

```go
		// Warn tier only: archived scope is publish-only, so an application-field
		// extension there can never be cleared by a tool-mediated write. The
		// structural error above stays scope-inclusive.
		if e.scope != "active" {
			continue
		}
```

`internal/status/field_conformance.go`, as the first statement of the `fieldConformanceWarnings` entity loop:

```go
		if e.scope != "active" {
			continue
		}
```

and its doc comment changes "for every active + archived entity" to "for every active entity" with the publish-only rationale.

**The placement is the whole design, and it is not the obvious one.** `gateValidationDiagnostics` returns `errs` and `warns` from a single loop over a single entity slice. Filtering the slice at the `validateWorkflow` call site — the smaller-looking edit, one line instead of two guards — also silently drops archived *structural* gate errors. Spike 2 below proves that error is fatal (exit 1) today. The guard therefore has to sit after the error branch and before the warning loop, per entity, not at the call site.

**No new mechanism.** The predicate `e.scope != "active"` is already the codebase's scope filter, used verbatim at validate.go:211 (external-proof sub-check), validate.go:326-332 (`anyActive`), and discover.go:257. Every entity reaching these functions is constructed by `newEntity` with scope explicitly `"active"` (discover.go:158,207) or overwritten to `"archived"` (discover.go:304); no path yields empty scope, so the predicate is total.

### Alternatives considered

1. Filter the entity slice at the `validateWorkflow` call site. Rejected on evidence, not taste: spike 2 shows it drops the archived `Error: invalid gates:` and turns exit 1 into exit 0. This is the alternative that looks simpler and is wrong.
2. Hand-migrate the 125 archived entities. Rejected: `6c45fd59c` already ruled that silencing a diagnostic by hand-editing archived terminal entities is the wrong direction, and archived scope is publish-only.
3. Trim or silence the pre-commit hook echo. Rejected and out of scope: the echo has a demonstrated consumer — it surfaces real errors and the hook blocks on exit 1. Muting the transport would hide genuine failures to fix a noise problem in the payload.
4. Add a `--quiet` / `--no-warn` flag. Rejected: adds command grammar to work around a default that is itself information-free, and every caller (hook, first officer, CI) would have to opt in. The default would stay wrong.
5. Delete both warn channels outright. Rejected: for active entities these findings are actionable, and that is the live consumer this task preserves.

### Coordination: `remove-gate-validate-subcommand`

Same program, and the two rationales were circular. That entity's seed justifies removal partly because the warning class "prints at state publish" — which is this sweep. My first draft justified this scoping by noting `gate validate archive:<slug>` still reports on demand (verified live below). Each entity was citing the surface the other removes.

Checked rather than assumed: `internal/cli/cli.go:337` is the only non-test consumer of `gates.SummaryFileDiagnosticsAt`, and there is no Go `publish` command. If both land, the application-extension warning is unreachable through any CLI for archived entities.

That outcome is acceptable, but only on the standalone argument, so this task rests on it alone: the finding is unactionable in archived scope by construction, therefore not emitting it loses no actionable information. The `gate validate` reachability is recorded below as evidence about today's behavior, not as this task's justification. Active-scope coverage is unaffected either way.

**Resolved 2026-08-15.** The sibling ensign independently re-verified all three claims against HEAD, agreed, and committed the rewrite as `c52bc0326`: their entity now disclaims leaning on this one, states the archived-scope reachability loss as a declared semantic change rather than a silent consequence, and rests on unactionability for the same reason this one does. Neither entity now cites the other's surface. Their substantive grounds — zero usage across 424 attempts, retained-authority faults refused by every gate write path, the round check duplicated by `gate record --round` — do not depend on this task, and mine does not depend on theirs. Diffs stay file-disjoint (`internal/status/` here; `internal/cli/`, `internal/gates/io.go`, tests and docs there), so there is no merge or ordering dependency in either direction. They also confirmed that this task's `e.scope != "active"` skip leaves their AC-4 verifier intact, since `gateApplicationWarningFixture` plants its entity at the workflow root in active scope.

They raised one adjacent finding, verified here rather than taken on report. In the main tree `SummaryFileDiagnosticsAt` has exactly two consumers: `internal/cli/cli.go:337` and `SummaryFileAt` (io.go:199), which discards the warnings slice with `_`. `SummaryFileAt` in turn is reached only by `SummaryFile` (io.go:195) and a test; `SummaryFile` is reached only by `application_test.go:365`. So once cli.go:337 goes, the whole `SummaryFile`/`SummaryFileAt`/`SummaryFileDiagnosticsAt` trio is production-unreachable, reached from tests only — which undercuts the approved keep-boundary that retains the reader because `SummaryFileAt` consumes it. That belongs to neither entity's scope; both of us recommend filing it separately. The reasons to keep it for now are that it is a general read surface rather than gate-validate machinery, and that removing it means re-pointing `prepare_test.go:956` and `application_test.go:365`, which currently prove retained-authority refusal through it. Notably it is *not* that removal would orphan `nearestWorkflowDir`: that helper has five further callers (`operation.go:212`/`:277`/`:345`, `application.go:72`/`:109`) and would be unaffected — a reason I asserted to the sibling ensign and they correctly refuted.

### Spike record

Both spikes ran against HEAD as throwaway tests in `internal/status/`, then were deleted; they seed the two new tests in the test plan.

1. **Does an `_archive/`-placed fixture reach both warn channels?** Yes. A `_archive/task.md` entity carrying an unknown `application` field and `verdict: MAYBE` produced `Warning: unknown gate application field 'action' ... scope=archived` and `Warning: field 'verdict' value "MAYBE" ... scope=archived`, exit 0, stdout `VALID`. This confirms `validationFixture`-style subpath planting yields archived scope through the validate path, so the new test has a genuine red-before-green.
2. **Is an archived structural gate error fatal today?** Yes. An `_archive/`-placed entity with an unrelated key inside the gates block exits **1** with `Error: invalid gates: ... field unrelated not found in type gates.Attempt ... scope=archived`. This is the keep-boundary the call-site filter would have broken.

3. **Does the proposed patch actually reach 0, and does the suite survive it?** Yes to both, measured rather than predicted. The two guards above were applied to a scratch build, run against the live `docs/dev` corpus, and reverted (`git checkout --`; working tree returned clean, ideation ships no code). Result: 126 lines / 51457 bytes / 125 archived warns before, **1 line / 6 bytes / 0 archived warns after**, exit 0 and stdout `VALID` in both. `go test ./internal/status/` passed with the guards applied (237s) — the package holding every validate, field-conformance, golden-envelope, and independent-parity test — so AC-2 and the golden fixtures are already known to survive.

Verified live on the real corpus (evidence about current behavior, not load-bearing justification per the coordination note): `gate validate archive:ac2-reanchor-live-scenario-repair --workflow-dir docs/dev` prints both application-field warnings through `SummaryFileDiagnosticsAt`, a path this change does not touch.

Nothing else is unverified: the entity enumeration, scope assignment, and warn-vs-error split are all read directly from HEAD and exercised above. Implementation is a re-application of a proven patch plus the two tests, which is why the estimate carries a tight tolerance.

### Documentation

`status --validate` stderr is user-visible, and one contract sentence becomes imprecise. `docs/site/reference/frontmatter-contract.md` line 13:

Before:
> those keys produce warnings only on explicit `status --validate` or `gate validate`, are ignored for authority, and are never written.

After:
> those keys produce warnings only on explicit `status --validate` over active entities or `gate validate`, are ignored for authority, and are never written.

No change to `docs/site/reference/command-reference.md`: its `status` row does not describe the warn tier, and its `gate validate` row (line 96) describes a path this task leaves alone. If `remove-gate-validate-subcommand` lands first, that entity owns dropping the `or gate validate` clause; this diff stays correct either way.

## Out of scope

The pre-commit hook echo (kept - it has a demonstrated consumer). The read tolerance - load-bearing for every gates read over the legacy corpus. Active-scope warnings. The `gate validate` surface. Structural gate errors, flat/folder conflict detection, stage-name validation, and id validation, all of which stay scope-inclusive.

The verdict enum itself. `verdict: superseded` appears on 4 archived entities and is not in the schema's conventional set, so the warning is arguably correct rather than spurious — it is simply unactionable where it fires. This task does not widen the enum, because whether `superseded` is a legitimate terminal verdict or a misuse of the field is a decision with consequences for the terminal ceremony, not a cleanup. After this change the check still fires for active entities, so if an FO sets `verdict: superseded` on a live entity the diagnostic still reaches someone who can act on it. That is the correct place for the question to surface, and it is worth filing separately.

## Expected surface and tolerance

| File | Change | Lines |
|---|---|---|
| `internal/status/validate.go` | scope guard + comment in `gateValidationDiagnostics` | +5 |
| `internal/status/field_conformance.go` | scope guard + doc-comment correction | +6 |
| `internal/status/archived_warn_scope_test.go` | new, two tests | +60 |
| `docs/site/reference/frontmatter-contract.md` | one sentence | ~1 |

Estimate: ~72 insertions across 4 files. Tolerance: ±25 lines, ±1 file. Net LOC is positive and intended to be — the noise removal is runtime output, not source. The production change is 11 lines; the rest is the test that keeps it honest.

Observable semantics declared:

- **Command grammar:** unchanged. No flag, subcommand, or argument added or removed.
- **Stored formats:** unchanged. No frontmatter, gates, or schema write.
- **Authority:** unchanged. The warn tier never carried authority and never gated an exit code.
- **Runtime behavior:** exactly one change — `status --validate` stderr no longer carries warn lines for entities in archived scope. Exit codes are unchanged in every case, including archived structural errors. `gate validate`, the read-path validation pre-check (`native_runner.go:486`, which already opts out via `includeExternalProof=false`), and `--validate --json` stdout are all untouched.

## Acceptance criteria

**AC-1 (value) - A validate run over `docs/dev` emits zero warn lines carrying `scope=archived`, down from a recorded baseline of 125, and still exits 0 with stdout `VALID`.**
Verified by: the same built binary run against `docs/dev` immediately before and after the change, both recorded in the stage report — counting `Warning:` lines by scope and the captured byte size. Baseline at ideation: 125 archived-scope warn lines, 51457 bytes; target: 0 and 6 bytes. Re-measure the "before" number in the same session rather than trusting 125, since the corpus is under concurrent write by this program. Moves the wrong way if the guard is too narrow (count stays above 0) or too broad (exit flips to 1, or `VALID` disappears).

**AC-2 (preservation) - An active entity with an unknown gate application field or an out-of-enum verdict still produces its `Warning:` line under `--validate`.**
Verified by: `TestStatusValidateReportsGateApplicationExtensionsAsWarnings` (asserts all five extension fields warn and the count is exactly 5) and `TestFieldConformanceWarnsSurface` (asserts a `Warning: field '<name>'` line per schema violation) stay green unmodified. Both fixtures are top-level placed, so they fail if the guard's predicate is inverted or catches active scope.

**AC-3 (keep-boundary) - An archived entity whose gates block is structurally invalid still fails validation with `Error: invalid gates:` and exit 1.**
Verified by: a new test planting an `_archive/`-placed entity with an unrelated key inside its gates block, asserting exit 1 and the `Error: invalid gates:` line. Fails if the scope filter is applied to the entity slice at the call site instead of at the two warn emission points — the specific wrong implementation spike 2 identified.

**AC-4 - The suite stays green.**
Verified by: `go test ./...` and `go test ./... -race` pass.

## Test plan

One new file, `internal/status/archived_warn_scope_test.go`, fixture-level through `runNative` — no live workflow needed, matching how both existing warn-channel tests are written. Roughly 60 lines, low cost.

1. `TestValidateSkipsArchivedWarnChannels` (AC-1's mechanism, AC-2's inverse). Plants one `_archive/`-placed entity carrying both defect kinds at once — an unknown `application` field and `verdict: MAYBE` — and asserts stderr contains no `scope=archived` substring, with exit 0 and stdout `VALID`. Asserting on `scope=archived` rather than on each message text makes one assertion cover both channels and any warn channel added later. Spike 1 proved this fixture currently emits three such lines, so the test is red before the change.
2. `TestValidateKeepsArchivedGateErrorsFatal` (AC-3). Plants an `_archive/`-placed entity with an unrelated key inside the gates block; asserts exit 1 and `Error: invalid gates:`. Spike 2 proved this passes at HEAD, so it is a guard against the regression the change could introduce rather than a red-to-green test.

Active-scope preservation needs no new test: `gate_application_warning_test.go` and `field_conformance_test.go` already pin both channels on top-level fixtures and serve as the regression floor unmodified. AC-1's corpus measurement is a recorded one-off run, not a test, since it depends on the live `docs/dev` corpus.

## Stage Report: ideation

- DONE: Only the two warn-tier channels filter to active scope; structural and ID validation stay scope-inclusive
  Guards specified at the two warn emission points (validate.go `gateValidationDiagnostics` after the error branch, field_conformance.go `fieldConformanceWarnings` loop head); spike 2 proves the alternative call-site filter would have dropped the archived `Error: invalid gates:` and flipped exit 1 to 0, which is why placement is per-emission not per-slice.
- DONE: An archived-fixture test asserts warn absence; existing active-fixture tests stay green
  Both tests specified in the test plan with their falsifying conditions; ideation ships no code, so the permanent file is implementation's to write. The fixture mechanism is proven, not assumed: spike 1 shows an `_archive/`-placed entity currently emits three `scope=archived` warn lines, so the new test is genuinely red before the change. Existing tests confirmed green under the real patch (spike 3).
- DONE: Value AC measured on docs/dev: 125 archived-scope warnings before, 0 after, report still VALID
  Measured, not predicted (spike 3): scratch build with both guards run against live `docs/dev` gave 126 lines / 51457 bytes / 125 archived warns before and 1 line / 6 bytes / 0 archived warns after, exit 0 and stdout `VALID` in both. Source reverted with `git checkout --`; `git status --porcelain internal/status/` clean.

### Summary

The removal scope is confirmed against HEAD (4d1912a69) and every keep-boundary named: structural gate errors, flat/folder and stage-name checks, all id-style validation, the gates read tolerance, the pre-commit hook echo, and the `gate validate` path all stay as they are. One seed claim was corrected against the corpus: the 4 verdict warnings are `verdict: superseded`, a token absent from the schema enum, not the case-variance residue `6c45fd59c` fixed — 178 lowercase verdicts exist and stay silent because of that commit's case-insensitive read. Whether `superseded` should join the enum is flagged out of scope and worth filing. The load-bearing finding is that `gateValidationDiagnostics` returns errors and warnings from one loop, so the smaller-looking edit — filtering the entity slice at the call site — silently drops archived structural gate errors; spike 2 measured that error as fatal today, and AC-3 was added to pin it. Three spikes ran and were reverted, including a full apply-measure-revert of the proposed patch, so implementation is re-applying a proven 11-line change plus two tests.

One cross-entity finding needs the gate's attention: `remove-gate-validate-subcommand` and this entity had circular justifications — its seed cites the `status --validate` sweep as the surviving surface for the warning class, while this entity's first draft cited `gate validate` as the surviving surface for archived scope. `internal/cli/cli.go:337` is the only non-test consumer of `gates.SummaryFileDiagnosticsAt` and there is no Go `publish` command, so if both land the warning becomes unreachable by CLI for archived entities. That is acceptable, but only on the standalone argument that the finding is unactionable in publish-only terminal scope; the rationale here was rewritten to rest on that alone, and the sibling ensign was notified. Active-scope coverage is unaffected either way, and the two diffs are file-disjoint.

## Stage Report: implementation

- DONE: Re-apply the proven per-emission guards exactly as gated, never the call-site filter; archived structural gate errors must still fail closed
  Commit `c275dee24` on `spacedock-ensign/scope-validate-warnings-to-active-entities`: `e.scope != "active"` guard in `validate.go` `gateValidationDiagnostics` sits after the `err != nil` error branch and before the compatibility loop, and at the head of the `field_conformance.go` `fieldConformanceWarnings` entity loop. Neither `validateWorkflow` call site was touched. `TestValidateKeepsArchivedGateErrorsFatal` pins the boundary: an `_archive/`-placed entity with an unrelated key inside its gates block still exits 1 with `Error: invalid gates:`; it fails if the guard is ever moved to the entity slice at the call site.
- DONE: Archived-fixture test proven red before the change and green after; existing active-fixture tests untouched and green
  New `internal/status/archived_warn_scope_test.go`. Red before the guards with the exact three `scope=archived` lines spike 1 predicted (application fields `action` + `blockers`, plus `verdict "MAYBE"`); green after. `TestValidateSkipsArchivedWarnChannels` asserts stderr carries no `scope=archived` substring with exit 0 and stdout `VALID`, so it fails if either guard is dropped or its predicate inverted. AC-2's preservation tests were not edited — `git show --stat c275dee24` lists one new test file and no modified test file — and `TestStatusValidateReportsGateApplicationExtensionsAsWarnings` (5 extension fields warn, count exactly 5), `TestStatusValidateKeepsUnrelatedGateUnknownFieldsFatal`, and all `TestFieldConformance*` subtests pass; they fail if the guard catches active scope.
- DONE: One-off measure on docs/dev recorded in the report: 125 archived warns to 0, VALID, both streams compared
  Same source built before and after the guards, both binaries run back to back against the live `docs/dev`. Before: exit 0, stdout `VALID` (6 bytes), stderr 51451 bytes / 125 lines, 125 `scope=archived` warns (121 unknown gate application field — 60 `action`, 60 `blockers`, 1 `execution-hold` — and 4 `verdict "superseded"`), 0 active-scope warns, 0 errors. After: exit 0, stdout `VALID` (6 bytes, byte-identical), stderr 0 bytes / 0 lines, 0 archived warns, 0 errors. Both streams captured separately rather than merged with `2>&1`.

### Summary

Implementation is the ideation patch re-applied unchanged plus the two specified tests and the one-sentence contract correction in `docs/site/reference/frontmatter-contract.md`: 74 insertions across 4 files against an estimate of ~72 across 4, inside the declared ±25-line / ±1-file tolerance. AC-1, AC-2, and AC-3 are all satisfied by direct exercise rather than inspection.

AC-4 needs a caveat, and the caveat is not this change. `go test ./...` at Go's default per-package timeout fails in `internal/ensigncycle` — but so does the unmodified merge-base `4d1912a69`, checked in a throwaway worktree, which times out at the same 10m in the same package (`TestDurableQuestionedRejectsTerminalHistory` mid-run there, `TestDurableTaskJourneys` mid-run here; the package simply needs more than 10m on this machine). Given an adequate budget, `go test ./... -timeout 40m` and `go test ./... -race -timeout 60m` both leave `internal/status` green (245s and 107s) and `internal/ensigncycle` green under race. Two failures survive, neither reachable from this diff. `internal/cli` `TestCodexResolveManifestAgainstInstalledHost` fails in both runs and identically at merge-base with the machine-level cause in its own output — `Failed to read config file /Users/clkao/.codex/config.toml: Operation not permitted (os error 1)`. `internal/ensigncycle` `TestCodexProcessActivityResetsQuietBudget` failed only in the loaded plain run, holds a 250ms quiet budget over a helper subprocess, and passes 5/5 unloaded and under race. Both are reported rather than dispositioned: no candidate byte was changed for either, and neither is owned by this task.

## Stage Report: validation

- DONE: Independently re-exercise the ACs against worktree commit c275dee24, never by reading the report: build the binary and run status --validate over docs/dev - zero archived-scope warn lines, still VALID, both streams compared to the recorded 125-line baseline
  Binaries built from clean 4d1912a69 and c275dee24, run back to back over live docs/dev, streams captured separately: base 125 warn lines / 51451 stderr bytes, all scope=archived (60 `action`, 60 `blockers`, 1 `execution-hold`, 4 `verdict "superseded"` — matches the recorded baseline exactly); branch 0 warn lines / 0 stderr bytes; both exit 0 with byte-identical 6-byte stdout `VALID`. `--validate --json` stdout also byte-identical across binaries with empty branch stderr.
- DONE: The guards must be per-emission, not call-site: verify an archived entity with structurally invalid gates still fails closed (the spike-2 hazard), and the archived-fixture warn-absence test is red on merge-base and green on branch
  Diff inspection: guards sit inside `gateValidationDiagnostics` after the invalid-gates error branch and at the `fieldConformanceWarnings` loop head; no `validateWorkflow` call-site edit; nothing error-tier follows either guard. An independently constructed `_archive/` fixture with an unrelated key inside gates exits 1 with `Error: invalid gates: ... scope=archived` on the branch binary, byte-identical to base. `archived_warn_scope_test.go` copied onto a throwaway 4d1912a69 checkout: `TestValidateSkipsArchivedWarnChannels` fails there with exactly the three predicted `scope=archived` lines; `TestValidateKeepsArchivedGateErrorsFatal` passes there (regression guard, as designed); both green on branch.
- DONE: Existing active-fixture tests untouched and green (gate_application_warning_test.go count of 5, field_conformance active enum case); suite differential against clean 4d1912a69; verdict PASSED or REJECTED with per-AC citations
  `git show --stat c275dee24`: one new test file, no existing test modified. On branch: `TestStatusValidateReportsGateApplicationExtensionsAsWarnings` (count exactly 5 asserted at gate_application_warning_test.go:36 — fails if the guard catches active scope), `TestStatusValidateKeepsUnrelatedGateUnknownFieldsFatal`, and all `TestFieldConformance*` subtests including verdict-out-of-enum pass. `go test ./... -timeout 40m` and `go test ./... -race -timeout 60m` on branch each show one failure — `internal/cli` `TestCodexResolveManifestAgainstInstalledHost` — reproduced identically at clean 4d1912a69 (machine cannot read `~/.codex/config.toml`, "Operation not permitted"); `internal/status` and `internal/ensigncycle` green in both runs. No failure attributable to the diff.

### Summary

Verdict: PASSED, all four ACs verified by independent exercise. AC-1: 125 archived-scope warns to 0 over docs/dev, `VALID` stdout and exit 0 unchanged, measured with separately built base/branch binaries. AC-2: both active warn channels preserved by unmodified fixture tests whose assertions fail if the predicate inverts. AC-3: archived structural gate errors still fail closed, proven live and by the new test; the spike-2 call-site hazard is absent by diff inspection and red/green evidence. AC-4: plain and race suites green except one pre-existing environmental failure identical at clean merge-base, reported not dispositioned, owned by neither this task nor this machine's diff.

Adversarial pass clean: warn lines never reach exit codes or `--json` stdout (handlers.go routes warns to stderr before the errs/exit decision), the scope predicate is total (discovery assigns only "active"/"archived"), `fieldConformanceWarnings` is warn-only by construction (`isWarnSeverity` filter), the read-path pre-check already opts out and discards warns, and `gate validate` still reports archived application-field warnings live. gofmt clean. Surface 74 insertions / 4 files against estimate ~72 / 4, inside declared tolerance. No material findings and no deferred risks introduced; nothing new enters Review-finding disposition.
