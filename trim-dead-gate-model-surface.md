---
id: ebgwr177kjjs6w5thhywz408
title: Trim dead gate model and projection surface
status: ideation
source: "0.27 cut audit (2026-08-14), adversarially verified; captain directed filing"
sprint: durable-decisions
started: 2026-08-15T02:55:35Z
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:ebgwr177kjjs6w5thhywz408:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:ebgwr177kjjs6w5thhywz408-backlog-1
              briefing:
                id: briefing:ebgwr177kjjs6w5thhywz408:backlog:attempt-1:revision-1
                digest: sha256:fc72271948d2273a9e0ede89b05604c4741457bd98fadaf50b9ebd1e7457b14e
                request-digest: sha256:5f016522677589c77f88f557e4e5f1d306c9f507dbce176995941f0789e7da34
                room-ref: ./trim-dead-gate-model-surface/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ebgwr177kjjs6w5thhywz408:backlog:1
                briefing: briefing:ebgwr177kjjs6w5thhywz408:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T02:53:37.741016Z"
                decision: approve
                reason: 'Captain ruling 2026-08-14 (dispatch them): approved into ideation'
              application:
                target-stage: ideation
                state: consumed
---
Remove four verified-dead pieces of the gate model surface. Verified against HEAD `4d1912a69`.

1. The nine projected `gate-*` status columns (`internal/status/discover.go:217-228`). No skill, agent, CI lane, live doc, or live-history query names one. Keep the `gates.Read` call at `discover.go:216` — it feeds the load-bearing gate-readiness chain. Also reconcile the `docs/site/reference/frontmatter-contract.md:13` sentence that names three of the nine.
2. `Summary.Condition` and `Summary.Eligible` (`internal/gates/model.go:95-96`). Debris from the eligibility cut (`013c8729e`). No writer, no reader.
3. `Annotation.Target`, `.Kind`, `.Body` (`internal/gates/review.go:18-20`). Decoded, never validated, projected, or read. The retained raw review log preserves the bytes.
4. The `ReadWithWarnings` and `Diagnostic` aliases (`internal/gates/io.go`). Zero callers, and an internal package permits no external caller.

Captain override recorded here: the archived decision `remove-standalone-gate-eligibility.md:173` said to keep `gate-state`, `gate-application`, and `gate-target-stage`. That retention named no consumer, and the audit re-confirmed none exists at HEAD — the three columns are read by no skill, agent, doc, or CI lane, and only by tests that exist to assert the projection itself. Approval of this entity overrides that sentence. The override is deliberate and narrow: it retires the three named columns only, and changes nothing about the gate model those columns projected.

## Problem

Each of the four pieces has a broken value chain. Verified item by item at HEAD:

**1. The nine columns have no consumer and one user-visible surface.** `newEntity` writes `gate`, `gate-attempt`, `gate-state`, `gate-briefing`, `gate-resolution`, `gate-decision`, `gate-application`, `gate-application-state`, and `gate-target-stage` into every scanned entity's field map. None is in `defaultEntityKeys`, so none appears in the default table. They surface two ways: when named explicitly via `--fields`, and — because `--all-fields` scans the keys present on each entity — in every `--all-fields` run. Outside tests, nothing names them: no file under `skills/`, `agents/`, or `docs/site/` names any of the nine except the one `frontmatter-contract.md` sentence this entity reconciles.

**2. `Summary.Condition` and `Summary.Eligible` are unreachable.** `CurrentSummary` (`model.go:363-384`) is the sole constructor of a populated `Summary`, and it sets neither field. Every `Summary` in the program therefore carries the zero value. Note the near-miss: the separate `Eligibility` struct has fields of the same two names, and those are live and load-bearing (`application.go`, `delivery.go`, `merge.go:573-582`, `gate_ceremony.go:76-77`). A textual grep for `.Condition` or `.Eligible` cannot tell the two apart, so the compiler — not grep — is the verifier for this item.

**3. The three `Annotation` prose fields are decoded and dropped.** `parseReviewLog` decodes review-log JSONL into `reviewEntry`, and `validateAnnotation` checks only `Type`, `ID`, `Briefing`, `By`, `At`, and `Includes`. The one consumer of the parsed entries, `ValidateRoundFile` (`round.go:131-136`), projects only `Type`, `ID`, `Decision`, and `Advisory`. The other decode site, `providerResult.Annotations`, is read only by `verifyProviderResolution` (`operation.go:593-609`), which touches `Type`, `ID`, and `Briefing`. No path re-marshals an `Annotation`: the only `json.Marshal` in the package is `prepare.go:727`, and every digest over retained review or provider bytes is taken on the raw file bytes (`RawDigest`, `CanonicalDigest`). Removing the three fields therefore cannot move a digest.

**4. The two aliases have zero callers.** `ReadWithWarnings` forwards to `ReadDiagnostics`; `Diagnostic` aliases `Warning`. Both are unreferenced, and `internal/gates` admits no caller outside this module.

**Correction to the seed.** The seed declared "No observable semantics change: no consumer exists for the removed surface." The first half is false and the second half is what actually holds. `status --all-fields` output does change — nine keys disappear. No consumer depends on them, but the output is user-visible and the change must be declared, not waved off. The seed's AC-2 was also insufficient for item 2, for the `Eligibility` name-collision reason above.

## Proposed approach

Delete the four surfaces, keep every boundary named below, retire the tests that exist only to assert the removed projection, repair the tests that used it incidentally, and reconcile the one live doc sentence.

### Keep-boundaries

Each of these stays, and the removal is only correct if it does:

- **`gates.Read` at `discover.go:216`.** The call and its `doc` / `gateErr` results stay. They populate `entity.gateDoc` and `entity.gateInvalid`, which feed `materializeGateReadiness` and `CurrentStageReadinessWithReport`. Only the nine `fields[...]` assignments and their enclosing `if gateErr == nil` block go.
- **The `gate-readiness` column and the whole readiness chain.** Unchanged in code and in output.
- **`Eligibility.Condition` and `Eligibility.Eligible`.** Live; only the `Summary` twins of the same name go.
- **`ReadDiagnostics`, `Warning`, `FormatWarning`, `SummaryFileDiagnosticsAt`,** and the `status --validate` / `gate validate` warning surface they serve.
- **`filterApplicationMappings`** and its bounded unknown-key read tolerance.
- **`Annotation.Type`, `.ID`, `.Briefing`, `.By`, `.At`, `.Includes`** — the fields `validateAnnotation` actually checks.
- **Every retained byte.** Review logs, provider Results, presented inventories, and their frozen digests are untouched.
- **Archived records that quote the old columns**, e.g. `docs/roadmap/durable-decisions/staff-review-sprint-close.md:493-495`. These are historical records of a past review, not live surface; they keep their text.
- **The `internal/ensigncycle` `gate-state: preserve-me` fixtures** (`shared_fixtures_test.go:161`, `shared_round_recording_test.go:171,264`). These pin byte-preservation of an author-written frontmatter key that merely shares a name with a projected column. They are not consumers and do not change. (Removing the projection incidentally makes such a user-authored `gate-state` key surface truthfully instead of being overwritten.)

### Test dispositions

Six tests reference the removed columns. Two exist only to assert the projection and go with it; four use it incidentally and keep their own claim:

- DELETE `TestStatusTextAndJSONProjectApprovedPendingApplication` (`gates_coexist_test.go`) — its entire subject is the removed projection.
- DELETE `TestStatusTextAndJSONProjectAllRecordedResolutionStates` — same.
- EDIT `TestStatusProjectsSharedGateReadinessReducer` — narrow `--fields` to `id,gate-readiness` and drop the two `2n`/`wd` column-detail assertions. Its readiness claim, which covers all five readiness values across seven entities, is untouched and becomes the proof for AC-4.
- EDIT `TestUnrelatedSetPreservesGatesAndStatusProjectsResolution` → rename `TestUnrelatedSetPreservesGates`; keep the byte-preservation half (an unrelated `--set` must not perturb the `gates` subtree) and drop the trailing projection assertion.
- EDIT `TestStatusValidateReportsGateApplicationExtensionsAsWarnings` — swap the ordinary-status probe from `--fields id,gate-application,gate-target-stage` to `--fields id,gate-readiness`, asserting the truncated `approved-awaiting-a` (the convention already used at `gates_coexist_test.go:112`). This keeps the load-bearing claim — ordinary status does not print the `--validate` warnings — and strengthens it, since the entity now proves a document with five unknown application keys still yields a usable readiness.
- EDIT `TestRoutedTerminalApprovalSurfacesExistingDisplay` (`internal/cli/terminal_consume_test.go`) — drop `gate-application` from the fields and the assertion; keep the `gate-readiness == approved-awaiting-merge` and status-unchanged claims, and update the comment, which currently names the "pending-application display" that no longer exists.

Also remove the now-unused `gates` import and the `decodeSingleStatusEntity` helper from `gates_coexist_test.go` (the helper's only callers are in the deleted test; the stale import is a hard compile error).

### One new test, and why

`TestArchivedProviderResultDecodesWithoutAnnotationProseFields` (~25 lines, `internal/gates`) drives a real retained provider Result copied from `_archive/fo-boot-install-hint-linux-direct-sandbox/review/ideation/briefing-1/provider/result.json` through `decodeProviderResult` and `verifyProviderResolution`.

- **Value AC it serves:** AC-5. Frozen archived authority must still decode after modeled fields are removed.
- **Simplest alternative considered:** rely on the existing `advisory-round` fixture, whose `briefing.review.jsonl` already carries `target`, `kind`, and `body`.
- **Why it is insufficient:** that fixture covers the review-log path only. The provider path uses a different decoder entry (`decodeAuthorityJSON`) against different retained bytes, and it is the exact path the concurrent `remove-provider-evidence-fields` entity proved is sensitive to decode strictness. The retained bytes also already carry a `selectors` key the struct has never modeled, so the test pins the tolerance itself. It fails if the authority decode is ever made strict.
- **Joint ownership:** if `remove-provider-evidence-fields` lands first, `providerResult` is gone and this test goes with it; the review-log coverage then stands alone.

### Doc diff

`docs/site/reference/frontmatter-contract.md`, final sentence of line 13:

    - `spacedock status --fields gate-state,gate-decision,gate-resolution` surfaces the current state without changing the default status table.
    + `spacedock status --fields gate-readiness` surfaces the current gate readiness without changing the default status table.

Replacement rather than deletion: the sentence answers "how do I see gate state from the CLI", and `gate-readiness` is the surviving answer. `gate-readiness` is named in no live doc today, so this also gives the load-bearing column its first documentation. Leaving the sentence unchanged would ship a documented command that renders three empty columns.

## Out of scope

The gate-readiness chain, `gates.Read`, `filterApplicationMappings` read tolerance, `ReadDiagnostics` and its `status --validate` consumer, resolution includes, mediaType.

## Expected surface and tolerance

Measured on the spike (full working change, all four removals plus test and doc reconciliation): **8 files, +8 / -132, net -124 lines.**

| File | + | - | net |
| --- | --- | --- | --- |
| `internal/status/gates_coexist_test.go` | 2 | 99 | -97 |
| `internal/status/discover.go` | 0 | 12 | -12 |
| `internal/gates/io.go` | 0 | 9 | -9 |
| `internal/gates/review.go` | 0 | 3 | -3 |
| `internal/gates/model.go` | 0 | 2 | -2 |
| `internal/cli/terminal_consume_test.go` | 3 | 4 | -1 |
| `internal/status/gate_application_warning_test.go` | 2 | 2 | 0 |
| `docs/site/reference/frontmatter-contract.md` | 1 | 1 | 0 |

Tolerance: ±2 files and ±40 lines, plus the ~25-line new test if it ships. The band is wide on the upside because two concurrent cut entities touch `internal/gates/model.go` and `internal/gates/io.go`, so a rebase can shift counts without changing this entity's scope.

### Declared semantic changes

Files and lines measure cost; these are the observable semantics this entity changes:

1. **`status --all-fields` output shrinks.** Text and JSON both lose the nine `gate-*` keys. `gate-readiness` is unaffected. Measured before/after on the same fixture — before: `gate`, `gate-application`, `gate-application-state`, `gate-attempt`, `gate-briefing`, `gate-decision`, `gate-readiness`, `gate-resolution`, `gate-state`, `gate-target-stage`; after: `gate-readiness` alone.
2. **`status --fields <removed-name>` renders an empty column.** The names stay accepted — `--fields` has no allowlist — so the header prints and the value is blank. No error, no exit-code change.
3. **Command grammar: unchanged.** No flag, subcommand, or argument shape changes.
4. **Stored formats and authority: unchanged.** No entity bytes, no gate document, no digest, no eligibility or readiness decision changes.
5. **Documentation:** the one sentence above.

## Acceptance criteria

**AC-1 — The change removes substantially more lines than it adds.**
Verified by: cumulative line delta against `origin/main` is negative, and no shallower than -80. Spike measured -124. This is the end-value measure; it moves the wrong way if the removal grows compensating scaffolding.

**AC-2 — No source references the removed columns, fields, or aliases.**
Verified by two checks, because neither alone is sufficient: grep across `cmd`, `internal`, and `skills` for the nine column names, `ReadWithWarnings`, and `gates.Diagnostic` returns no matches; and `go build ./... && go vet ./...` compile clean, which is the only verifier for `Summary.Condition` / `Summary.Eligible` since `.Condition` and `.Eligible` remain live on `Eligibility`.

**AC-3 — The suite stays green.**
Verified by: `go test ./...` and `go test ./... -race` pass, with one carve-out that predates this entity. `TestCodexResolveManifestAgainstInstalledHost` (`internal/cli/codex_resolve_test.go:33`) shells out to the installed `codex` binary and fails identically on untouched HEAD in this sandbox (`Failed to read config file ~/.codex/config.toml: Operation not permitted`). It is environment-dependent, not a regression, and implementation must re-confirm it fails the same way before and after. The slow packages run near the default 10-minute per-package timeout under load (`internal/cli` 479s, `internal/ensigncycle` 514s, and one cold concurrent run tripped 600s on `internal/cli`), so pass `-timeout 40m` to keep the run deterministic.

**AC-4 — The gate-readiness chain is unchanged end-to-end.**
Verified by: `TestStatusProjectsSharedGateReadinessReducer`, which asserts all five readiness values — `validating`, `awaiting-captain`, `withdrawn-awaiting-prepare`, `approved-awaiting-merge`, `approved-awaiting-advance` — across a seven-entity split-root corpus. It fails if the removal disturbs `gates.Read`, the reducer, or the readiness projection.

**AC-5 — Retained gate bytes still decode and verify.**
Verified by: `TestArchivedProviderResultDecodesWithoutAnnotationProseFields` against a real archived provider Result carrying annotation `target`, `kind`, `body`, and the never-modeled `selectors`; plus the existing `advisory-round` fixture for the review-log path. Fails if removing modeled fields turns a tolerated unknown key into a decode or verification error.

**AC-6 — The live doc names no removed column, and the command it does name is true.**
Verified by: grep of `docs/site/` for the nine names returns nothing, and running the command the reconciled sentence names prints a populated `GATE-READINESS` column.

## Test plan

Deletion plus the existing suite, two small test repairs, and one new ~25-line test. Tests that exist only to assert the removed projections are deleted with them; no test is weakened to keep it passing. Cost is low: no new fixture infrastructure, one copied archived fixture file. `internal/gates` and `internal/status` run in roughly 2.5 and 4 minutes respectively; `internal/cli` needs the raised timeout noted in AC-3.

## Spike record

The riskiest unverified mechanism was whether removing modeled fields from `Annotation` could break retained-byte decoding or move a frozen digest. It was exercised first, on an isolated `git archive` copy of HEAD, never in the shared repo root. Results:

- **All four removals compile clean.** `go build ./...` and `go vet ./...` (which compiles tests) report nothing. This is the compiler proof that nothing reads `Summary.Condition` or `Summary.Eligible`.
- **`internal/gates` stays fully green with no test changes at all** (`ok`, 138s), including the `advisory-round` fixture whose review log carries `target`, `kind`, and `body`. The `Annotation`, `Summary`, and alias removals break nothing.
- **The nine-column removal breaks exactly five tests**, all in `internal/status`, all on the removed-projection assertion. Baseline on the same copy was green (`internal/gates` ok 158s, `internal/status` ok 258s). After the test dispositions above, all five pass again.
- **A real archived provider Result decodes and verifies** under the reduced `Annotation`. The retained bytes carry `target`, `kind`, `body`, and `selectors`; `selectors` is already an unknown key today, so the removal moves the other three into a tolerance the production decoder already exercises. `decodeAuthorityJSON` uses plain `json.Unmarshal` — lenient — and the sole strict decode in the package, `KnownFields(true)` at `io.go:98`, applies to the YAML `gates` frontmatter, which models neither `Summary` nor `Annotation`.
- **The observable change was measured, not assumed**, by building both binaries and diffing `--all-fields --json` on one fixture. This is what falsified the seed's "no observable semantics change".
- **The whole suite was run on the finished change**, not just the touched packages: `go test ./...` is green across every package except the pre-existing `TestCodexResolveManifestAgainstInstalledHost` environment failure described in AC-3, which reproduces identically on the untouched baseline copy.
- **The reconciled doc sentence was proven true and the old one proven false** by running both commands against the patched binary: `--fields gate-readiness` prints `approved-awaiting-a…`; `--fields gate-state,gate-decision,gate-resolution` prints three empty columns.

## Coordination with concurrent 0.27 cut entities

- **`remove-gate-validate-subcommand` — direct textual conflict.** Both entities edit the same paragraph, `docs/site/reference/frontmatter-contract.md:13`. That paragraph also contains "warnings only on explicit `status --validate` or `gate validate`", which their removal makes false, while this entity replaces the paragraph's final sentence. The two edits must be sequenced or hand-merged; they do not conflict semantically. Their declared keeps (`SummaryFileDiagnosticsAt`, `ValidateRoundFile`) match this entity's keep-boundaries exactly, and the `status --validate` warning test this entity repairs survives their cut untouched.
- **`remove-provider-evidence-fields` — same files, disjoint regions, one shared trap.** Both touch `internal/gates/model.go` and `internal/gates/io.go` in non-overlapping places. Their spike found that `KnownFields(true)` makes removing a *YAML-modeled* gate field break six archived records. That trap was checked against this entity's removals and does not apply: `Summary` is a derived projection with no YAML tags, and `Annotation` is JSON-decoded leniently. If their entity lands first, delete `TestArchivedProviderResultDecodesWithoutAnnotationProseFields` with `providerResult`.
- **`remove-startup-capability-probe` and `retire-requires-contract-sentinel`** — no file, symbol, or doc overlap with this entity.

## Stage Report: ideation

- DONE: Removal keeps gates.Read and the gate-readiness chain green; only the nine projections and dead fields go
  Spiked the full removal on an isolated `git archive` copy of HEAD: `internal/gates` stayed green with zero test edits (ok, 138s), and the nine-column cut broke exactly five `internal/status` tests, every one on the removed-projection assertion. `TestStatusProjectsSharedGateReadinessReducer` still returned all five readiness values (its failure output shows `gate-readiness` correct while only the removed columns went empty). Re-pointing that test to `--fields id,gate-readiness` makes it AC-4's proof; it fails if `gates.Read` or the reducer is disturbed.
- DONE: The captain override of the eligibility-cut retention sentence is recorded in the design
  Recorded in the body preamble with the re-verified basis: `remove-standalone-gate-eligibility.md:173` retained `gate-state`, `gate-application`, `gate-target-stage`; the audit re-confirmed at HEAD that no skill, agent, doc, or CI lane reads them, and the override is scoped to those three columns only, changing nothing about the gate model beneath them.
- DONE: Frontmatter-contract doc sentence reconciled with the removal
  Concrete before/after diff recorded for `docs/site/reference/frontmatter-contract.md:13`, replacing the `--fields gate-state,gate-decision,gate-resolution` sentence with a `--fields gate-readiness` one. Both halves proven against the patched binary: the new command prints a populated `GATE-READINESS` column; the old one prints three empty columns, so leaving the sentence would ship documented-but-dead output.

### Summary

Confirmed all four removals against HEAD `4d1912a69` and spiked the complete change end to end (net -124 lines across 8 files) without touching the shared repo root. Two seed claims were corrected: the removal *does* change observable output — `status --all-fields` loses nine keys, measured by diffing two built binaries — and the seed's grep-based AC could not verify `Summary.Condition`/`Summary.Eligible`, because `Eligibility` carries live fields of the same names, so the compiler is the verifier. The riskiest mechanism, whether dropping `Annotation` fields breaks retained bytes or moves a frozen digest, was exercised against a real archived provider Result carrying `target`, `kind`, `body`, and the never-modeled `selectors`; it decodes and verifies unchanged, because the authority decoder is lenient and the sole strict decode covers YAML gate frontmatter that models neither struct. Every keep-boundary is named, all six affected tests have a delete-or-repair disposition with rationale, and the two real overlaps with concurrent cut entities are recorded — including a direct textual conflict with `remove-gate-validate-subcommand` on the same doc paragraph.
