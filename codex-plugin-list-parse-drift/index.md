---
id: 5pb3h1ewewvrjayhgyabx5y4
title: spacedock codex can't detect an installed codex plugin — codex plugin list parse drift
status: validation
source: captain (2026-06-01) — reproduced LIVE: `spacedock codex` reports no plugin while `codex plugin list` shows spacedock@spacedock installed, enabled 0.19.2
started: 2026-06-02T02:04:04Z
completed:
verdict:
score: "0.40"
worktree: .worktrees/spacedock-ensign-codex-plugin-list-parse-drift
issue:
mod-block: merge:pr-merge
pr: "#258"
---

`spacedock codex` refuses to launch — "no installed codex plugin found. Run `spacedock install --host codex`" — even when the codex plugin IS installed and enabled. Reproduced live against the just-released 0.19.2: `codex plugin list` shows `spacedock@spacedock  installed, enabled  0.19.2   https://github.com/spacedock-dev/spacedock.git, ref next`, but the codex front door can't see it. This is a launcher-replacement regression (Goal #1) shipped in 0.19.2 — it warrants a fast 0.19.3 patch.

## Root cause (verified in code + live)
`codexEntryInstalled` (`internal/cli/host_exec.go:96`) greps for the literal PAREN form:
```go
if strings.Contains(line, id+" (installed") { return true }
```
i.e. it expects `spacedock@spacedock (installed`. But the installed codex renders the COMMA form with NO parens. Captured live from `codex plugin list` (codex-cli 0.136.0 on this box; same comma format the 0.19.2 report saw — the format, not the version, is what drifted):
```
PLUGIN               STATUS              VERSION  PATH
spacedock@spacedock  installed, enabled  0.12.1   /Users/.../marketplaces/spacedock/plugins/spacedock
```
The status token is `installed,` (comma, NO parens), and the row sits in a column-aligned TABLE with a `PLUGIN/STATUS/VERSION/PATH` header. So the `Contains(id+" (installed")` never matches → `resolveCodexManifest` returns "" (not installed) → `spacedock codex` (and `spacedock doctor --host codex`) report no-plugin-found. The doc-comment (`:93-95`) asserts the format is `<id> (installed[, enabled]) | (not installed)` — wrong for current codex.

**Why the drift shipped — the real coverage gap:** `codexEntryInstalled` IS unit-tested (`TestCodexEntryInstalled`, `codex_resolver_unit_test.go:96`), and it passes — but every case tests the PAREN form (`(installed)`, `(installed, enabled)`, `(not installed)`). The real comma form is absent from the test table, so the test green-lit a predicate that fails on the format codex actually emits. The gap is "tests the wrong format," not "no test." Closing it means adding captured-real-format cases that fail against today's predicate.

## Spike — riskiest mechanism exercised first (throwaway, result recorded)
The design rests on one unverified mechanism: does a field-split parse-tolerant predicate correctly discriminate over REAL codex table output (column-aligned, header row, a marketplace PATH cell that contains the bare word `spacedock`)? Exercised before locking the plan:
1. Confirmed the current predicate returns `false` for the real comma form → drift reproduced.
2. Prototyped the proposed predicate (`strings.Fields` → find field == id → `strings.Trim(next, "(),")` → `== "installed"`) against the live-captured table block and four discriminators. All held:
   - real `spacedock@spacedock  installed, enabled  0.12.1  <path>` → ✓ matched
   - `chrome@openai-bundled  not installed` → ✗ rejected (next token is `not`)
   - legacy `spacedock@spacedock (installed, enabled)` → ✓ matched (codex-revert tolerance)
   - foreign id (`browser@openai-bundled`) against the spacedock block → ✗ rejected
   - the marketplace PATH line `/Users/.../marketplaces/spacedock/plugins/spacedock` (contains bare `spacedock`, not the id as a field) → ✗ no false match

Outcome: the field-split + trim approach is sound against real output. The implementer's first unit test seeds from this exact captured block. Refinement over the seed fix-direction: `strings.Trim(tok, "(),")` collapses the "strip leading `(`, trailing `,`/`)`" steps into one, and the field-equality guard (not substring) is what rejects the PATH-cell `spacedock`.

## Acceptance criteria

- **AC-1 — Parse-tolerant FIELD matching pinned to real codex output.** `codexEntryInstalled` reports a plugin installed iff a whitespace-delimited field on some line equals the id exactly AND the next field, after stripping surrounding `()` and a trailing `,`, equals `installed` exactly. It matches the current comma form `<id>  installed, enabled  <ver>  <path>`, matches the legacy paren form `<id> (installed[, enabled])`, and REJECTS `<id>  not installed` (next field is `not`), a foreign-id line, and the marketplace PATH line that contains the bare marketplace word but not the id as a field. The bare substring `installed` is never matched (it is contained in `not installed`).
  - **Verified by:** Go unit test `TestCodexEntryInstalledRealFormats` (extend/replace the table in `internal/cli/codex_resolver_unit_test.go:96`) whose cases include the live-captured comma block above (header + data rows, multi-space-aligned), the `not installed` row, the legacy paren forms, a foreign-id line, and the marketplace PATH line. The new comma-form and PATH-line cases MUST fail against the current paren-literal predicate (run them red first) and pass after the fix. `go test ./internal/cli/` green.

- **AC-2 — End-to-end `spacedock codex` resolve path verified, not just the predicate.** After `codexEntryInstalled` passes, `resolveCodexManifest` resolves a cached manifest dir (`latestVersionDir`/`codexHome`) and `spacedock codex` reaches launch (does not abort with "no installed codex plugin found"). The predicate-only fix is insufficient on its own — the cache-resolution tail must still land on a real `plugin.json`.
  - **Verified by:** a fixture-cache Go test that drives `resolveCodexManifest`'s cache tail — builds a temp `CODEX_HOME` with `plugins/cache/spacedock/spacedock/<ver>/.codex-plugin/plugin.json`, and asserts the resolved manifest path is that file (this is the half the unit predicate test does not cover). PLUS a live confirmation recorded in the validation/implementation stage report: against the box's installed codex (0.136.0, spacedock@spacedock installed, enabled), `spacedock codex` no longer reports no-plugin-found and proceeds to launch. The fixture test is the gate; the live run is the real-world proof captured in the report.

- **AC-3 — Scope boundary stated and decided (no ambiguity).** Exactly two things are IN: the `codexEntryInstalled` parse-fix and its captured-real-format unit tests (AC-1), and the end-to-end resolve verification (AC-2). The `latestVersionDir` semver-vs-lexical concern flagged in the seed is decided **OUT — already fixed, nothing to ship**: `latestVersionDir` (`host_exec.go:123`) already orders via the integer-per-component `compareVersion` (`host_exec.go:147`) and is already guarded by `TestLatestVersionDirSemverOrder` / `TestLatestVersionDirSemverNotLexical` / `TestLatestVersionDirSingleVersion` (all passing). The seed's "latent lexical bug" note was stale. Also OUT of this 0.19.3 item: the codex resolver no longer needs the cache-only layout fallback now that `codex plugin list` carries a `PATH` column (a separate follow-up — the resolver could read the listed PATH directly; not changed here to keep the patch minimal and the front door fixed fast).
  - **Verified by:** the diff touches only `codexEntryInstalled` (+ its doc-comment) in `host_exec.go` and the test files; `latestVersionDir`/`compareVersion` are unchanged (git diff shows no edit to lines 118-176). The three `TestLatestVersionDir*` tests still pass unmodified, confirming the OUT decision rests on existing green coverage, not a promise.

## Test plan
- **Predicate unit test (AC-1)** — extend `TestCodexEntryInstalled` (or add `TestCodexEntryInstalledRealFormats`) in `internal/cli/codex_resolver_unit_test.go`. Cases: live-captured comma table block (installed), `not installed` row, `(installed)`/`(installed, enabled)` legacy, foreign-id line, marketplace PATH line. Cost: trivial (table-driven, no I/O). Authoring order: add the comma + PATH cases, run red against current code, then fix the predicate, run green.
- **Cache-tail fixture test (AC-2)** — temp `CODEX_HOME` with the `plugins/cache/spacedock/spacedock/<ver>/.codex-plugin/plugin.json` layout; assert `resolveCodexManifest`'s post-predicate resolution returns that path. Note: `resolveCodexManifest` shells real `codex plugin list`, so the fixture either targets the cache-resolution helpers directly or the test box must have codex+spacedock installed; the implementer picks the seam — record which. Cost: low (filesystem fixture under `t.TempDir()`).
- **Live e2e (AC-2 proof, recorded not gated)** — on a box with codex installed and spacedock@spacedock enabled, run `spacedock codex` and confirm it no longer reports no-plugin-found; capture the outcome in the stage report. Cost: low (single command); availability-dependent, so it is recorded evidence, not the CI gate.
- No new live-workflow smoke test needed beyond the above; the change is a single-predicate parse fix with a fixture-backed cache tail.

## Notes / related
- This is `internal/cli` (the codex resolver) — NOT the serialized `internal/status` lane, so it does not collide with status work; but it IS the host front door (high-stakes launch path) → **validation gets a detached adversarial audit** (not this stage's job; noted for the record).
- Ships via a normal PR onto `next`; cut 0.19.3 after it lands (the codex front door is broken in the shipped 0.19.2).
- Stale-assumption cleanup in scope of AC-1's doc-comment edit: the `:93-95` comment claiming the `<id> (installed…)` paren format must be corrected to the real comma/table format so the next reader is not re-misled.

## Stage Report: ideation

- DONE: AC pins parse-tolerant FIELD matching to REAL codex output: matches the comma form `<id>  installed, enabled  <ver>` AND the legacy paren form `<id> (installed`, matches `installed` exactly, and REJECTS `not installed` plus foreign-id lines — with a captured-real-format Go unit test named as the proof (this closes the untested gap that let the drift ship in 0.19.2).
  AC-1 written; field-equality + `Trim(next,"(),")=="installed"` proven against live-captured codex 0.136.0 table in the spike; proof test named `TestCodexEntryInstalledRealFormats`. Corrected the seed's "untested" claim: a test exists but covers only paren forms — the real gap is the missing comma/PATH-line cases (must run red first).
- DONE: Test plan verifies the END-TO-END `spacedock codex` resolve path, not just the predicate: after codexEntryInstalled passes, latestVersionDir/codexHome must resolve a cached manifest dir and the front door must actually launch — name the verification (live codex resolve and/or a fixture-cache test).
  AC-2 names both: a temp-`CODEX_HOME` fixture-cache Go test (the gate) asserting `resolveCodexManifest`'s cache tail lands on a real `plugin.json`, PLUS a live `spacedock codex` run on the codex-installed box recorded in the impl/validation report (the real-world proof).
- DONE: Scope boundary stated explicitly: codexEntryInstalled parse-fix + tests is IN; the latent latestVersionDir lexical-vs-semver bug (host_exec.go ~:121) is decided IN (fold the fix + a guarding test) or OUT (separate 0.19.3 item) — not left ambiguous.
  AC-3 decides it OUT — already fixed: `latestVersionDir`→`compareVersion` is integer-per-component semver and already guarded by 3 passing `TestLatestVersionDir*` tests; the seed note was stale. The OUT decision is backed by existing green coverage, verified by the diff leaving lines 118-176 untouched.

### Summary
Hardened the AC + test plan and locked scope behavior-first. Exercised the riskiest unknown first: confirmed the current paren-literal predicate returns false on real codex output (drift reproduced) and that the proposed field-split + `Trim(.,"(),")` predicate correctly discriminates real comma form / not-installed / legacy paren / foreign-id / marketplace-PATH-line against a LIVE-captured codex 0.136.0 table block — the implementer's first test seeds from that exact block. Two seed claims were stale and corrected on the record: `codexEntryInstalled` is in fact tested (but only the wrong paren format, which is the real gap), and the `latestVersionDir` lexical bug is already fixed and guarded — so AC-3 decides it OUT with existing-coverage backing rather than re-opening it.

## Stage Report: implementation

- DONE: RED-FIRST proof — the new comma-form table case must FAIL against the current paren-literal predicate before the fix, then PASS after.
  `TestCodexEntryInstalledRealFormats/real-comma-table` ran RED (`= false, want true`) against the paren-literal predicate, then PASS after the fix (full subtest run captured both states). The marketplace-PATH-line case already passed pre-fix (current code also rejects it) — its role is a regression guard on the new field-split predicate, not a re-green; the comma-form failure is what proves the test catches the real drift. This is noted rather than contrived to fail (forcing it red would mean writing a deliberately-wrong predicate first).
- DONE: Minimal fix per the spike — field-equality + `strings.Trim(next,"(),")=="installed"` in codexEntryInstalled, doc-comment corrected from the paren claim to the real comma/table format; `go test ./internal/cli/` green.
  `host_exec.go` commit 12a6193f: `strings.Fields(line)` → field == id → `strings.Trim(fields[i+1],"(),")=="installed"`; `:92-103` doc-comment rewritten to the table/comma format (paren form noted as tolerated legacy). `go test ./internal/cli/` green; full repo `go test ./...` green; `go vet` clean. latestVersionDir/compareVersion untouched.
- DONE: AC-2 end-to-end — fixture-cache test green AND a live `spacedock codex` run recorded showing the front door no longer reports 'no installed codex plugin found'; note which seam the fixture targets.
  Seam: extracted the post-predicate cache tail into `codexCacheManifest()` and drove it directly with a temp `CODEX_HOME` holding `plugins/cache/spacedock/spacedock/0.12.1/.codex-plugin/plugin.json` (`TestCodexCacheManifestResolvesCachedPluginJSON` asserts the resolved path == that file; `TestCodexCacheManifestAbsentCache` covers degradation) — this targets the cache-resolution helper, not the codex-shelling `resolveCodexManifest`, so it needs no installed host. Live front-door run (built to /tmp, gitignored-equivalent, removed after): BEFORE fix `spacedock codex -- --version` → "Spacedock: no installed codex plugin found." (exit 1); AFTER fix → "plugin-predates-contract: your installed Spacedock plugin is out of date" (exit 1) — the no-plugin gate is gone; the front door now resolves the cached manifest and reaches the contract verdict (expected, since the box's installed 0.12.1 predates the dev binary's contract). The live RUN-verified `TestCodexResolveManifestAgainstInstalledHost` also exercised the installed host and PASSED (returned the real `/Users/clkao/.codex/.../0.12.1/.codex-plugin/plugin.json`), after its own paren-literal install-check was routed through the fixed predicate.

### Summary
Replaced codex's paren-literal install detection with field-based matching pinned to the real comma/table format codex 0.136.0 emits (header row, multi-space columns, marketplace PATH cell). The comma-form unit case failed RED against the old predicate then passed after the fix, closing the "tested the wrong format" gap. Extracted `codexCacheManifest()` to make the post-predicate cache tail fixture-testable without shelling codex. Live proof: the front door went from "no installed codex plugin found" to a normal contract verdict, confirming the resolver now sees the installed plugin. AC-3 scope held — `latestVersionDir`/`compareVersion` are untouched. One in-scope cleanup beyond the AC: the live RUN-verified resolver test duplicated the same paren-literal drift to decide "installed", so it was routed through `codexEntryInstalled`; without it that test asserted the old broken no-plugin behavior and failed once the resolver was correct.

## Stage Report: validation

- DONE: Reproduce AC-1's red-first proof yourself: confirm TestCodexEntryInstalledRealFormats/real-comma-table FAILS against the pre-fix paren-literal predicate and PASSES after — i.e. the test genuinely catches the real drift, not just greens. Also confirm the predicate rejects `not installed`, a foreign id, and the marketplace PATH line.
  Reverted predicate to pre-fix `Contains(id+" (installed")` → `real-comma-table` FAILED red (`= false, want true`); restored fix → all 6 subtests PASS. The comma-form red is the load-bearing proof; `not-installed-comma`/`marketplace-path-line`/`foreign-id-installed` all return false (field-equality + next-field read, not substring). Also ran the production predicate over the UNEDITED 199-line live `codex plugin list` capture: detected spacedock@spacedock installed, rejected a `not installed` plugin and an absent id.
- DONE: Reproduce AC-2 end-to-end: run the codexCacheManifest fixture test (temp CODEX_HOME → real plugin.json) AND independently confirm the live front door — `spacedock codex` no longer reports 'no installed codex plugin found' against the installed codex. Verify the cache seam is actually on resolveCodexManifest's real path, not a decoy.
  `TestCodexCacheManifest{ResolvesCachedPluginJSON,AbsentCache}` green. Seam confirmed real: `resolveCodexManifest` (host_exec.go:80) `return codexCacheManifest()` — the exact fn the fixture drives. Live before/after dev builds against the box's installed codex 0.136.0 (spacedock@spacedock installed, enabled 0.12.1): BEFORE `spacedock codex -- --version` → "no installed codex plugin found" (exit 1); AFTER → "plugin-predates-contract" (exit 1) — no-plugin gate gone, resolver sees the plugin and reaches the contract verdict (exit-1 contract mismatch expected: installed 0.12.1 predates the dev binary's contract).
- DONE: AC-3 scope held: latestVersionDir/compareVersion untouched; the second drift-site fix in codex_resolve_test.go is a legitimate same-bug fix (not masking a regression). Full `go test ./...` green + `go vet` clean. Then a PASSED/REJECTED recommendation.
  Diff hunks span host_exec.go:77-125; `latestVersionDir` (line 145) and `compareVersion` (line 169) are below the hunks — definitions untouched. `codex_resolve_test.go:29` change (paren-literal `Contains` → `codexEntryInstalled`) fixes a duplicate of the SAME drift: the old line would have made the RUN-verified test assert the broken no-plugin path and fail once the resolver was correct — not masking a regression. `TestCodexResolveManifestAgainstInstalledHost` RUN-verified (0.08s, not skipped). Full `go test ./...` EXIT=0 (10 ok, 2 no-test, 0 FAIL); `go vet ./...` EXIT=0.

### Summary
PASSED. All three ACs independently reproduced by exercising behavior, not re-reading. AC-1's red-first holds: the comma-form case genuinely fails against the pre-fix paren-literal predicate and passes after — and the production predicate correctly discriminates over the unedited live codex capture. AC-2's cache seam is on `resolveCodexManifest`'s real path (not a decoy), the fixture test gates it, and the live front door transitioned from no-plugin-found to a contract verdict on the installed host. AC-3 scope held — `latestVersionDir`/`compareVersion` untouched, the second drift-site fix is the same bug at a second call site (not a regression mask), full suite green and vet clean. Worktree left clean after all temporary red-first reverts were restored (empty `git diff`).

## Stage Report: implementation (cycle 2 — post-validation doc touch-up)

- DONE: Doc-only completion within AC-1's doc-correction scope, requested after validation PASSED — the adversarial audit found stale sibling references the host_exec.go fix did not cover.
  Commit 5705e4f7 on worktree branch `spacedock-ensign/codex-plugin-list-parse-drift` (NOT pushed — pr-mod pushes on captain approval). Three doc-comment edits, no logic/test-assertion changes: (1) `codex_resolver_unit_test.go` `TestCodexEntryInstalled` doc-comment rewritten — it claimed the paren form is THE format; now describes the current comma/table form as primary (cross-refs `TestCodexEntryInstalledRealFormats`) and notes paren as tolerated legacy. (2)+(3) the two stale "Codex 0.132.0" version refs in `host_exec.go` (`ResolveManifest`, `resolveCodexManifest` doc-comments) refreshed to the observed 0.136.0. Also found and refreshed a third same-stale ref the FO's note scoped to host_exec.go did not name: `codex_resolve_test.go:14` ("codex 0.132.0 rejects --json") → 0.136.0 — same misleading claim in a sibling doc-comment, left would re-mislead the next reader. `go test ./internal/cli/` green after the edits (comments-only, behavior unchanged).

### Summary
Comments/doc-only follow-up after validation PASSED: corrected the `TestCodexEntryInstalled` doc-comment that still presented the paren form as THE codex list format, and refreshed all three stale "codex 0.132.0" version references (two in host_exec.go per the audit, plus one sibling in codex_resolve_test.go) to the observed 0.136.0. No logic or test-assertion changes; `go test ./internal/cli/` stays green. Committed 5705e4f7 on the worktree branch, unpushed.
