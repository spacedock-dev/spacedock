---
id: n3md630s32x3x8zey698xkn6
title: Frontmatter writer doesn't quote YAML-ambiguous values (leading-#, internal colon) — breaks pr-mod + GitHub render
status: validation
source: "FO (2026-06-02) — pr=#N eaten by s0's comment-strip (live: jg #260, #258, #259) AND captain-observed GitHub YAML parse errors on colon-bearing source values"
started: 2026-06-03T00:51:44Z
completed:
verdict:
score: "0.34"
worktree: .worktrees/spacedock-ensign-frontmatter-hash-quoting
issue:
mod-block: merge:pr-merge
---

The frontmatter writer used to emit `pr: #260` and `source: a: b` bare, so strict YAML readers (the new yaml.v3 reader landed in zj, GitHub's frontmatter renderer) read them as a comment-only line and a nested mapping respectively. zj (PR #274, just shipped at HEAD `75cc2b81`) closes the writer side of both classes at `internal/status/mutate.go` (`needsExplicitQuoting` + `setScalarValue`): a value with a leading `#`, an internal `: ` (colon-space), or a ` #`/`\t#` inline-comment trigger is double-quoted on every write. New writes are safe; the live in-session smoke confirmed it (this session set `pr=#273` PLAIN on qs and the writer round-tripped the value correctly).

What's left for n3, post-zj, is the legacy-data side: 13 archived entities still carry bare `pr: #N` values written before zj's writer policy. They are terminal (archived, never re-written by `--set`), but the yaml.v3 reader returns their `pr` field as the empty string and GitHub renders them broken. A one-time re-quote pass fixes them; no writer or contract change is needed.

## Current state (audited)

- **Writer-quoting shipped (zj, post-zj baseline).** `needsExplicitQuoting` at `mutate.go:218-239` returns true for: a leading `#`, an internal `: `, a ` #`/`\t#` inline-comment trigger. `setScalarValue` at `mutate.go:199-216` applies `yaml.DoubleQuotedStyle` when it returns true. Confirmed by reading `internal/status/mutate.go` at HEAD and verified live: this session's `spacedock status --set qs pr=#273` PLAIN auto-quoted to `pr: "#273"` and round-trips.
- **FO contract has no `pr='"#N"'` workaround to remove.** The pr-merge mod at `docs/dev/_mods/pr-merge.md:64` writes the value PLAIN today: `spacedock status --workflow-dir docs/dev --set {slug} pr=#{N}`. The first-officer shared core writes the local-merge sentinel PLAIN: `pr=local-merge:{short-sha}` (single colon, no following space → not a hazard). A repo-wide grep of `skills/`, `docs/dev/_mods/`, `mods/` for `pr='?"#`, `pr=\"#`, `pr='#` finds zero hits. The `pr="#260"` form the original entity body cited as a "workaround applied to live state" was a captain-typed escape on individual entity files before zj shipped, not contract text — there is nothing to delete.
- **Legacy values that still misread, surveyed exhaustively.** A raw-vs-yaml.v3 audit (mimicking the production reader; throwaway scanner, deleted) over all 78 frontmatter blocks in `docs/dev/.spacedock-state/` (21 active + 57 archived) returns EXACTLY 13 hazards, all of class "bare `pr: #N`", all in archived entities — listed in the migration table below. The four other classes (colon-space `: `, `#`-comment inline trigger ` #`/`\t#`, other YAML indicators `! & * | > % @ \`` , leading/trailing whitespace) had ZERO matches anywhere. The session's earlier hand-fix of `_archive/front-door-plugin-dir/index.md`'s colon-bearing `source:` value (zj's spike's single live-parse hazard) is in the audit's clean baseline.

## Migration table (the 13 affected archived entities)

All status `done`, all hand-fixed by re-writing the line `pr: #N` → `pr: "#N"` in the state checkout. No `--set` invocation — these are terminal, the file is the artifact, the write is a single sed-equivalent line edit.

| slug | bare value | quoted value |
|---|---|---|
| architecture-review-cleanups | `pr: #247` | `pr: "#247"` |
| claude-runtime-segregation | `pr: #246` | `pr: "#246"` |
| cli-cobra-redesign | `pr: #241` | `pr: "#241"` |
| cli-ergonomics | `pr: #242` | `pr: "#242"` |
| host-neutrality-seam | `pr: #244` | `pr: "#244"` |
| init-upgrade-and-contract-remedy | `pr: #240` | `pr: "#240"` |
| live-e2e-pr-trigger-ergonomics | `pr: #245` | `pr: "#245"` |
| merge-ceremony-ergonomics | `pr: #255` | `pr: "#255"` |
| mods-definition-dir-location | `pr: #252` | `pr: "#252"` |
| no-hidden-machine-dependencies | `pr: #249` | `pr: "#249"` |
| ship-working-principles-in-contract | `pr: #248` | `pr: "#248"` |
| split-root-collaboration-remote-resume | `pr: #256` | `pr: "#256"` |
| status-enumeration-and-validation | `pr: #254` | `pr: "#254"` |

The set is small enough to hand-fix in one path-scoped state commit per file (or one path-scoped sweep) without writing a migration helper. A helper would be more code than the data it migrates, and zj's writer-quoting prevents the class from re-introducing itself.

## Acceptance criteria

**AC-1 — every active + archived entity in `docs/dev/.spacedock-state/` round-trips its frontmatter `pr` field through a strict yaml.v3 read.** The 13 archived `pr: #N` entries above are re-quoted to `pr: "#N"`; no remaining `pr:` line decodes to the empty string when the file is parsed.
Verified by: `spacedock status --workflow-dir docs/dev --include-archive --where "pr !="` returns a non-empty result (today it returns zero entities), and the row count matches the count of files that carry a `pr:` line. A throwaway raw-vs-yaml.v3 scan over the state-checkout tree, run from the captain's terminal, reports zero `pr`-empty-but-line-present mismatches.

**AC-2 — the regression class is pinned by a Go test in `internal/status`.** A test asserts that for every `*/index.md` under `internal/status/testdata/` (existing fixtures), `parseFrontmatter` returns a non-empty value for any frontmatter key whose raw line is non-empty after the `key:` prefix and not a comment-only line. The same test catches a fresh-future leading-`#` regression in the writer policy.
Verified by: a new `internal/status/no_yaml_silent_drop_test.go` adds `TestNoSilentYAMLValueDrop` that walks `testdata/`, parses each entity through `ParseFrontmatter`, and asserts no key with a raw non-empty value decodes to `""`. Adding a fixture with bare `pr: #99` makes the test fail; quoting it makes it pass. Test is green at HEAD against existing fixtures.

**AC-3 — the legacy-data correction is committed to the state branch and visible on GitHub.** The 13 re-quoted entities land in `spacedock-state/dev` (path-scoped, one commit covering all 13 files is fine — they are terminal archive entities, no concurrent writer). GitHub's frontmatter rendering of any of the 13 entity URLs no longer shows an empty `pr:` row.
Verified by: a `git log -- docs/dev/.spacedock-state/_archive/architecture-review-cleanups/index.md` shows the re-quote commit; opening one of the 13 entity files on GitHub renders the `pr` field as the quoted scalar string (`#247`), not as an empty value or a YAML error. The captain spot-checks one URL.

## Out of scope

- Re-implementing what zj already shipped (writer-quoting for `#`-leading + `: ` + ` #`/`\t#`).
- Modifying yaml.v3 reader behavior (which is now the production reader; loud parse-failure on malformed input is the divergence #4 contract zj documented).
- Other YAML-indicator classes not observed in any real value (`!`, `&`, `*`, `|`, `>`, `%`, `@`, backtick, leading/trailing space) — zero hits across 78 entities; deferred until a real hit appears, at which point zj's `needsExplicitQuoting` is the one place to extend.
- The captain-mentioned colon-bearing `source:` GitHub-render issue beyond the single hand-fixed `_archive/front-door-plugin-dir/index.md` — that file already carries the quoted form in the current state checkout, and no other entity has the hazard.
- A migration helper / walk-and-quote tool — 13 hand edits is cheaper than the helper + its idempotency test, and zj's writer prevents the class from re-introducing itself.

## Test plan

- **AC-1 verification — throwaway audit script.** A short Go program (or a `go run` snippet against `gopkg.in/yaml.v3`) walks the state checkout, parses each FM, and asserts zero `pr`-empty-but-line-present cases. Cost: trivial; runs in under a second over 78 files. Lives in `/tmp/` and is deleted after the implementer confirms 0 hits — it is the audit's run, not a test that ships.
- **AC-2 verification — Go test in `internal/status`.** `TestNoSilentYAMLValueDrop` in a new `no_yaml_silent_drop_test.go`. Uses `filepath.WalkDir` over `testdata/`, opens each `index.md`-shaped file, isolates the FM via the existing `frontmatterSlice` helper, parses via `ParseFrontmatter`, and asserts that for every key whose raw `key:` line has a non-empty post-`:` substring (after trimming whitespace and stripping a leading-`#` inline-comment-only line), the decoded value is non-empty. To validate the test catches a real regression, the implementer adds a temporary `testdata/no-silent-drop-regression/index.md` fixture with `pr: #99` (bare), runs `go test ./internal/status -run TestNoSilentYAMLValueDrop` and watches it FAIL, then changes the fixture to `pr: "#99"` and watches it PASS. Cost: low (one new file, ~40 lines, no new dependencies). Altitude is right: the claim is about the reader-vs-writer contract, and the test exercises the actual reader on actual fixture files.
- **AC-3 verification — git log + GitHub spot check.** No test framework needed; the verification is the commit landing and a captain-driven URL spot check on one of the 13 entities. Cost: trivial.

**Total cost: LOW.** No new behavior in the binary; one Go test (~40 lines) plus 13 one-line hand edits in archived entity files. No fixture, CLI, or live-workflow tests beyond AC-2's unit test.

## No spike needed

zj already exercised the load-bearing mechanism (a yaml.v3 round-trip on real entities, including the leading-`#` writer-quoting and the colon-space writer-quoting) at its ideation. The throwaway audit script I ran for AC-1 (78 files scanned, 0 parse errors, 13 known-class mismatches) IS the post-zj exercise of the only remaining unknown: "are there hazard classes zj's spike missed?" Answer: no — the same hazard zj found (and resolved at the writer) is the entire legacy-data residue. The implementation composes proven behavior: the existing reader, the existing writer, and 13 one-line edits.

## Notes
- `internal/status` lane (the regression-guard test + the legacy-data edits). No coordination needed with other status-lane items; the writer-policy lane (zj) is done.
- 0.19.4-class per FO recommendation — but the work has shrunk to 13 one-line edits + 1 test, so a fold into 0.19.3 is also reasonable if the captain wants.

## Stage Report: ideation

- DONE: The scope is sharpened against the CURRENT post-zj baseline.
  Read `internal/status/mutate.go` at HEAD `75cc2b81`: `needsExplicitQuoting` (lines 218-239) returns true for leading-`#`, internal `: ` (colon-space), AND ` #`/`\t#` inline-comment trigger; `setScalarValue` (lines 199-216) applies `DoubleQuotedStyle` when true. Confirmed in the body's "Current state (audited)" section. Repo-wide grep of `skills/`, `docs/dev/_mods/`, `mods/` for `pr='?"#`, `pr=\"#`, `pr='#` returns zero hits — the FO contract has no workaround to remove (the pr-merge mod and the FO core both write `pr=#{N}` and `pr=local-merge:{short-sha}` PLAIN, relying on zj's writer). A throwaway raw-vs-yaml.v3 audit over all 78 frontmatter blocks in `docs/dev/.spacedock-state/` returns EXACTLY 13 hazards, all class "bare `pr: #N`", all archived — listed in the body's migration table. Zero other hazard classes anywhere (colon-space, ` #` inline trigger, other YAML indicators).
- DONE: The ACs are entity-level and each cites a runnable check outside the entity body.
  AC-1 (`docs/dev/.spacedock-state/` is hazard-free) is verified by `spacedock status --workflow-dir docs/dev --include-archive --where "pr !="` returning a non-empty result (currently zero — the bug) plus a raw-vs-yaml audit reporting zero mismatches. AC-2 (regression-guard test) is verified by a new `TestNoSilentYAMLValueDrop` in `internal/status/no_yaml_silent_drop_test.go` that walks `testdata/`, asserts no key with a raw non-empty value decodes to `""`, and is exercise-validated by a temporary `pr: #99` fixture making it FAIL then `pr: "#99"` making it PASS. AC-3 (the 13 archived entities re-quoted and visible on GitHub) is verified by `git log` showing the commit and a captain spot-check on one of the 13 entity URLs.
- DONE: A test plan that names Go-test verifications.
  AC-2 names `TestNoSilentYAMLValueDrop` in `internal/status/no_yaml_silent_drop_test.go`. AC-1's audit is a throwaway `go run` script that does not ship. AC-3 is a git log + URL spot check. "No spike needed" is recorded with the proven mechanism: zj's spike already exercised yaml.v3 round-trip + writer-quoting on real entities; the throwaway audit I ran for this ideation is the post-zj exercise of "are there hazard classes zj missed?" — answer: no. Total cost LOW: 13 one-line edits + ~40 lines of Go test.

### Summary

Sharpened against the post-zj baseline: zj's writer-quoting at HEAD `75cc2b81` already covers leading-`#`, colon-space, and ` #`/`\t#`; the FO contract already writes `pr=#{N}` PLAIN with zero workaround sites to remove (confirmed by a repo-wide grep). What remains for n3 is the legacy-data residue: a raw-vs-yaml.v3 audit over all 78 frontmatter blocks in `docs/dev/.spacedock-state/` returns exactly 13 archived entities with bare `pr: #N` (yaml.v3 reads them as empty), and zero other hazard classes anywhere. ACs land at AC-1 (hazard-free state checkout, verified by a raw-vs-yaml scan + a `--where "pr !="` query that today returns zero), AC-2 (a new `TestNoSilentYAMLValueDrop` in `internal/status` pins the regression class on `testdata/`), and AC-3 (the 13 entity files re-quoted and visible on GitHub, verified by git log + spot check). No migration helper, no contract edit, no writer change — total work shrinks to ~40 lines of test + 13 one-line edits, scope honestly reflects that.

## Stage Report: implementation

- DONE: Re-quote the 13 archived entities listed in the entity body's migration table (lines 30-42): rewrite `pr: #N` → `pr: "#N"` in each `_archive/{slug}/index.md`. Path-scoped commit per the FO write-scope discipline.
  All 13 archived files edited (bare `pr: #N` → `pr: "#N"`) and landed in one path-scoped commit on `spacedock-state/dev` (`fix(archive): re-quote bare \`pr: #N\` in 13 terminal entities`), staging exactly the 13 listed paths so no concurrent writer's staged work is cross-attributed.
- DONE: Add `internal/status/no_yaml_silent_drop_test.go` with `TestNoSilentYAMLValueDrop`: walks `internal/status/testdata/`, isolates each FM via the existing `frontmatterSlice` helper, parses via `ParseFrontmatter`, and asserts no key with a raw non-empty value decodes to `""`. Exercise-validate.
  Test added on the CODE BRANCH (`spacedock-ensign/frontmatter-hash-quoting`, commit `cee6d05a`). Exercise-validation per dispatch: a temporary `testdata/no-silent-drop-regression/index.md` with bare `pr: #99` made the test FAIL (one assertion at no_yaml_silent_drop_test.go:51); changing it to `pr: "#99"` made the test PASS; the temp fixture was removed before committing.
- DONE: AC-1 verification: throwaway raw-vs-yaml.v3 scan over `docs/dev/.spacedock-state/` reports zero `pr`-empty-but-line-present cases. AC-2: `go test ./internal/status -run TestNoSilentYAMLValueDrop -v` PASS. AC-3: 13 edits land in `spacedock-state/dev`. Full repo `go test ./...` green.
  AC-1: throwaway `go run /tmp/yaml_audit.go` walked 68 frontmatter files in the state checkout and reported `Scanned 68 files, 0 hazards` (the scanner was deleted after). AC-2: `TestNoSilentYAMLValueDrop` PASS at HEAD with the temp fixture removed. AC-3: 13 archived files committed on `spacedock-state/dev`; push pending (the ensign push step follows this report). Full repo `go test ./...` returns all packages `ok` (cli 8.124s, status 6.602s, dispatch 5.827s, ensigncycle 3.158s, others sub-second; no FAIL anywhere).

### Summary

Implementation is two artifacts: (1) one path-scoped state-branch commit re-quoting `pr: #N` → `pr: "#N"` in the 13 terminal archived entities listed in the body's migration table, and (2) one code-branch commit adding `internal/status/no_yaml_silent_drop_test.go` with `TestNoSilentYAMLValueDrop`, which walks `testdata/` and asserts no key whose raw post-`:` value is non-empty decodes to `""`. The exercise-validation per the dispatch ran inline (temp fixture with bare `pr: #99` → FAIL; quoted → PASS; temp fixture removed). All three ACs are satisfied at this revision: AC-1's audit reports 0 hazards over 68 files, AC-2's new test PASSES at HEAD, AC-3 lands the 13 re-quotes on `spacedock-state/dev` and the full repo `go test ./...` is green.

## Stage Report: validation

- DONE: AC-1 verification at THIS revision — fresh raw-vs-yaml.v3 audit (yaml.Node-based, so it correctly distinguishes silent-drop null scalars from legit non-string decoded values like time.Time) over all `index.md` files under `docs/dev/.spacedock-state/`.
  Throwaway scanner (`/tmp/yaml_audit_validation.go`) reported `Scanned 68 files, 0 hazards` against state-checkout HEAD `92970066`; scanner deleted after.
- DONE: AC-2 verification — `go test ./internal/status -run TestNoSilentYAMLValueDrop -v` PASS at code-branch HEAD `cee6d05a` (the new test at `internal/status/no_yaml_silent_drop_test.go`).
  PASS confirmed; exercise-validated end-to-end: a temporary `testdata/no-silent-drop-regression-validation/index.md` with bare `pr: #99` made the test FAIL with `key "pr" has raw value "#99" but decodes to "" (yaml.v3 silently dropped it)` at `no_yaml_silent_drop_test.go:51`; quoting to `pr: "#99"` made it PASS; temp fixture removed; `git status` clean.
- DONE: AC-3 verification — the 13 listed archived entities on `spacedock-state/dev` carry `pr: "#N"` (quoted) and the re-quote commit is in the state log.
  `grep -l 'pr: "#'` over the 13 listed `_archive/{slug}/index.md` paths returned all 13; the three spot-checked lines are literally `pr: "#247"`, `pr: "#256"`, `pr: "#254"`; a yaml.v3 decode of two of those files returns the strings `"#247"` and `"#254"` (not empty). `git log --oneline | head` on state shows `ee18b424 fix(archive): re-quote bare \`pr: #N\` in 13 terminal entities`.
- DONE: Full repo `go test ./...` green at code-branch HEAD `cee6d05a`.
  `Go test: 808 passed in 12 packages` — no FAIL anywhere; the new test does not regress any existing test.

### Summary

Validation PASSED at code-branch HEAD `cee6d05a` and state-checkout HEAD `92970066`. All three ACs are satisfied with proof outside the entity body: AC-1 by a fresh re-run of the raw-vs-yaml.v3 audit (0 hazards over 68 files, with a yaml.Node-based classifier so non-string decoded scalars like time.Time are not mis-flagged), AC-2 by `TestNoSilentYAMLValueDrop` passing and exercise-validated against a temporary `pr: #99` fixture (FAIL → flip to quoted → PASS → fixture removed), and AC-3 by all 13 listed archived entities carrying the quoted `pr: "#N"` form plus the `ee18b424` re-quote commit visible in the state log. Recommendation: **PASSED**.
