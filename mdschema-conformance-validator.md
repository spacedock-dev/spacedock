---
id: tdpnhct3kqk99e5fj447c1xm
title: Full mdschema conformance validator (status --validate enforces a subset)
status: validation
source: captain-approved, surfaced by pt0 docs-site port (PR #343, 2026-06-13)
started: 2026-06-13T05:52:37Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-mdschema-conformance-validator
issue:
sprint: 0202-survey-improvements
group: cleanup
sprint-readiness: ready
---

`docs/schema/entity.mdschema.yml` + `workflow-readme.mdschema.yml` are now the SSOT for the frontmatter contract (ported from the v0 branch into the v1 repo during the #343 docs pass). But `spacedock status --validate` (`internal/status/validate.go`) enforces only a SUBSET of that schema; nothing checks full conformance against the mdschema files.

## Problem

`status --validate` currently checks: entity-form (flat/folder) conflicts, stage-name regex, per-id-style id presence/uniqueness, and the opt-in external-proof policy. It does NOT check per-field types/patterns (e.g. `verdict` in {PASSED, REJECTED}, `mod-block` pattern, `score` numeric coercion) or the schema invariants end to end. So the mdschema files can drift from what the binary actually enforces, and malformed frontmatter the schema would reject can pass `--validate`.

## Decision: (a) extend Go `status --validate`

**Chosen: (a) — extend the Go `--validate` path to cover the per-field types/patterns the entity mdschema declares.** Reject (b), the standalone Python/script checker.

Rationale:
- **The binary already owns this surface.** `--validate` (`internal/status/validate.go`, dispatched from `handlers.go:370`) already enforces the *structural* subset: flat/folder conflicts, stage-name regex, id presence/uniqueness per id-style, and the opt-in external-proof policy. Per-field type/pattern checks are the same kind of check on the same parsed entities — a new sibling sub-check, not a new tool.
- **Single source of enforcement.** A second checker (Python or otherwise) would re-implement frontmatter parsing (`ParseFrontmatter`, fence finder, last-key-wins) that already lives in Go, and would drift from it. The v0 `scripts/validate_frontmatter_contract.py` was deliberately NOT ported for exactly this reason: Python in the Go repo needs a CI hook to stay honest, and we have no such hook. One enforcement path, in the binary that ships, run by the FO's existing `--validate` invocation.
- **No new dependency.** The spike (below) proved `gopkg.in/yaml.v3` — already in `go.sum` — parses the mdschema files directly. JSON is a YAML subset, so the `.yml`-that-is-JSON schema files load into typed Go structs with no extra library.

**Schema-driven, not hand-transcribed.** The implementation reads `docs/schema/entity.mdschema.yml` and drives field checks from the parsed `fields` map (pattern, conventional-enum, type, severity), rather than re-spelling each pattern in Go. This is what closes the documented drift gap: the binary enforces *what the schema says*, so the schema cannot silently diverge from the binary. The schema file ships in-repo and is read at validate time (embedded via `go:embed` so there is no runtime filesystem dependency on `docs/`).

### Severity is load-bearing — warns must not gate reads

The riskiest *design* finding (not a mechanism risk): **every per-field check in the entity schema is `invalid_severity: warn`** — `verdict`, `mod-block`, `score`, `status`-unknown, and legacy-numeric ids all carry `warn`, not error. The current `--validate` is binary (any error → exit 1, else `VALID` exit 0); there is no warn tier today.

Two consequences the implementation MUST honor:
- **Field-conformance findings are warnings, not errors.** They print to stderr but do NOT flip the exit code, mirroring the schema's own declared severity. Exit 1 stays reserved for the existing structural-error classes plus any future schema field whose severity is explicitly `error`.
- **The read-path pre-check (`failOnValidationErrors`, cwd-gated `status`/`--next`/`--boot`) must NOT fire on warn-tier field findings.** `validate.go:140-145` already documents the lockout hazard: a read path that fails on a flagged AC locks the FO out of the listing that shows the broken entity. Warn-tier field conformance has the same hazard and must stay out of the read-path gate — surfaced only by the explicit `--validate` command.

This makes "full conformance" mean: the binary *checks and surfaces* every field the schema declares, at the severity the schema declares — not "every malformed field becomes a hard exit-1 error." A schema field marked `error` (none today) would gate; `warn` fields advise.

### Riskiest mechanism — exercised (spike PASSED)

The design's soundness rested on one unverified mechanism: **can Go parse the mdschema `.yml` (which is actually JSON) and check a known-bad fixture's frontmatter against the schema-derived pattern/enum end-to-end?** Exercised before committing to the plan:

- Parsed the real `docs/schema/entity.mdschema.yml` with `yaml.v3` into a typed struct (`fields[].pattern`, `.conventional`, `.invalid_severity`).
- Extracted `mod-block` pattern `^[^:]+:[^:]+$` (severity `warn`) and `verdict` conventional enum `[PASSED REJECTED]` (severity `warn`).
- Fed known-bad values: `mod-block: "noColonHere"` correctly FAILS the pattern; `verdict: MAYBE` correctly out of the enum.
- All assertions passed; the schema's own bytes are the source of truth for the expected pattern/enum, not a transcription.

Proven mechanisms the rest of the plan composes: `yaml.v3` JSON-subset parse of the schema file (spike), the existing `ParseFrontmatter` reader and `activeAndArchivedEntities` enumeration (shipped, covered by `frontmatter_test.go` / parity tests), and the existing `validateWorkflow` → exit-code wiring (shipped, `native_validate_test.go`). No further spike needed for those.

## Out of scope

- The #343 doc pass deliberately did schemas + light prose only (captain decision = option (a) territory, deferred). This task is the cleanup that closes the enforce-what-you-document gap.
- **Body-structure invariants** (`required_opening` problem-statement paragraph, `recognized_sections`, stage-report shape) — frontmatter field conformance is the closeable gap; body parsing is a larger separate task and not in this scope.
- **The workflow-README schema** (`workflow-readme.mdschema.yml`) — README stage-name regex is already enforced; full README field/invariant conformance is a follow-up, not this task.
- **Promoting any `warn` field to `error`.** This task enforces the schema's declared severities as-is. Changing a severity is a separate, captain-visible decision.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified. The independent source of truth in every behavioral AC is a **fixture entity file** (real frontmatter bytes) and the **schema file itself** — never a grep over instruction prose.

**AC-1 — `--validate` surfaces every per-field schema violation the entity schema declares.**
For each field carrying a `pattern` or `conventional` enum in `entity.mdschema.yml` (`mod-block`, `verdict`, `id` per id-style, `score`, ISO-8601 dates), a fixture entity whose frontmatter violates that field produces a matching diagnostic on stderr naming the field and the violated rule.
Verified by: a Go test in `internal/status` that builds fixture workflows with one known-bad field each and asserts the run's stderr contains the field-named diagnostic. The expected pattern/enum is read from the schema, and the bad value is in the fixture — both outside the test's own assertions.

**AC-2 — Field-conformance findings are warn-tier: they print but do not change the exit code.**
A workflow whose only defects are warn-tier field violations (e.g. `verdict: MAYBE`, `mod-block: noColon`) exits 0 from `spacedock status --validate` while still printing the warnings to stderr. A structural error (dup id, flat/folder conflict) still exits 1.
Verified by: a Go test asserting exit code 0 with non-empty warning stderr for a warn-only fixture, and exit 1 for a structural-error fixture — exit codes are observable command output, the fixtures are the independent source.

**AC-3 — Warn-tier field findings do not gate the read path.**
A workflow with a warn-tier field violation (and no structural error) still serves the default `status` table / `--next` / `--boot` without exit 1; the FO is never locked out by a field warning.
Verified by: a Go test running an enumerate op (default table) over a warn-violation fixture and asserting exit 0 with the table on stdout — mirrors the existing `TestNativeValidationGatesReads` shape inverted (warn must NOT gate).

**AC-4 — Enforcement is schema-driven: editing the schema changes what the binary checks, with no Go-source edit.**
The field checks read their patterns/enums from the embedded `entity.mdschema.yml`; the schema file is the source of the expected rules.
Verified by: a Go test that loads the embedded schema bytes, parses out a field's pattern (e.g. `mod-block`), and asserts the same pattern the validator uses to flag a fixture — proving the validator's rule and the schema file agree (two independently-readable values that can diverge), not a hardcoded Go regex literal.

## Test plan

**What verifies the implementation:** Go unit/behavior tests in `internal/status`, following the established `validationFixture` + `runNative` + golden-envelope pattern already used by `native_validate_test.go`. No Python, no new dependency (`yaml.v3` already vendored; schema embedded via `go:embed` so tests have no `docs/` filesystem dependency).

**Fixtures (independent source of truth):** small in-test workflow fixtures, one known-bad frontmatter field per case:
- `mod-block: noColonHere` (violates `^[^:]+:[^:]+$`)
- `verdict: MAYBE` (out of `[PASSED, REJECTED]`)
- `score: notanumber` (non-numeric where schema says numeric_string)
- a malformed-ISO `started`/`completed`
- a structural-error control (dup id) to confirm exit 1 still fires
Each fixture is the independent source — the test asserts the binary's diagnostic against the fixture's planted defect, never against a grep of any prose file.

**Cost / complexity:** Moderate. The schema-parse + field-loop is small (the spike is ~40 lines and already runs green). The bulk is wiring the warn tier into `validateWorkflow`'s return so warns are separable from errors at the exit-code boundary, plus the `go:embed` of the schema. Estimate: one implementation stage, ~4-6 focused tests. No live-workflow drive needed — all claims are command-level (exit code + stderr) and fixture-driven, provable by `go test ./internal/status/...`.

**Test tiers:** fixture + CLI-level (drive the native binary via `runNative`, assert exit + stderr). Golden envelopes for the diagnostic wording where stable. No live workflow smoke test required; runtime behavior is not the claim — frontmatter conformance is.

## Stage Report: ideation

- DONE: Decide and record the approach: (a) extend Go `status --validate` to full mdschema coverage, or (b) a standalone conformance checker — with rationale.
  Chose (a) — recorded under "## Decision: (a) extend Go `status --validate`": binary already owns the surface, single enforcement path, no new dependency; (b)/Python rejected (no CI hook, would re-implement frontmatter parsing).
- DONE: Riskiest mechanism FIRST: parse the mdschema YAML and check a known-bad fixture's frontmatter against it end-to-end — exercise it.
  Spike PASSED: `yaml.v3` (already in go.sum) parsed the real `entity.mdschema.yml` (JSON-in-.yml), extracted `mod-block` pattern + `verdict` enum + `invalid_severity`, and correctly failed `mod-block: noColonHere` and `verdict: MAYBE`. Result recorded under "### Riskiest mechanism — exercised". Surfaced the load-bearing finding that all per-field checks are `warn` severity.
- DONE: ACs: feed known-bad frontmatter fixtures, assert non-zero exit + the right diagnostic (the fixture is the independent source, never a prose-grep); one test plan.
  AC-1..AC-4 written, each grounded in a fixture entity + the schema file as independent sources; AC-2 corrects the naive "non-zero exit" assumption — warn-tier fields exit 0, only structural errors exit 1. Test plan names fixtures, tiers (fixture+CLI via `runNative`), no Python/new-dep, `go:embed` for the schema.

### Summary

Decided option (a): extend Go `--validate` to drive per-field type/pattern/enum checks from the embedded `entity.mdschema.yml`, closing the documented enforce-what-you-document drift gap with one enforcement path and no new dependency. The spike proved the schema (JSON-in-.yml) parses with the already-vendored `yaml.v3` and that a known-bad fixture is correctly rejected end-to-end. The key design finding: every per-field check is `invalid_severity: warn`, so "full conformance" means *surface at the declared severity* — field findings are warnings that print but must NOT flip the exit code or gate the read path (the documented FO-lockout hazard), with exit 1 reserved for existing structural errors. ACs and test plan are fixture-driven and command-level, no live drive needed.

## Stage Report: implementation

- DONE: Schema-driven per-field validation: extend internal/status/validate.go to drive field checks from the EMBEDDED entity.mdschema.yml (go:embed, yaml.v3) — mod-block pattern, verdict enum, score numeric, ISO-8601 started/completed — reading patterns/enums FROM the schema so editing the schema changes what the binary checks (AC-4).
  `schema_embed.go` (module root) embeds `docs/schema/entity.mdschema.yml`; `internal/status/field_conformance.go` parses it once (`loadEntitySchema`) and drives `fieldViolation` from the parsed `pattern`/`conventional`/`type`. AC-4 proven live: removed `PASSED` from the schema's verdict enum, rebuilt, binary warned on `verdict: PASSED` with zero Go edits; reverted (commit 76a924a7). `id`-per-id-style stays in the existing structural path (already enforced, gating) — the warn tier covers only the four fields with no prior enforcement; matches the test plan's fixtures (4 warn cases + dup-id structural control).
- DONE: Wire the WARN tier so field-conformance findings PRINT to stderr but do NOT flip the exit code (AC-2) AND do NOT gate the read path (AC-3); STAFF WATCH (P3): warns must NOT perturb the stdout=="VALID" golden.
  `validateWorkflow` now returns `(errs, warns)`; warns computed only on the explicit-`--validate` opt-in (same flag as external-proof), printed to stderr in handlers.go, exit code driven only by `errs`. Read-path `failOnValidationErrors` drops warns (passes `includeExternalProof=false`). Warns carry a `Warning:` prefix (vs structural `Error:`) via shared `entityEvidenceLine`; `entityEvidence` keeps its oracle-parity `Error:` shape untouched. STAFF WATCH held: `native_validate_test.go` golden (stdout=="VALID") still green — warns go to stderr only.
- DONE: Tests in internal/status over known-bad fixtures (mod-block noColonHere, verdict MAYBE, score notanumber, malformed-ISO started, plus dup-id structural control) — AC-1: field-named diagnostic on stderr; exit 0 warn-only / exit 1 structural; AC-4: a test parses a schema field pattern and asserts the validator agrees. `go test ./internal/status/` green over the WHOLE package.
  `internal/status/field_conformance_test.go`: 4 tests / 7 subtests, all PASS. Whole `internal/status` package green; whole module (`go test ./...`) green; `go vet` clean; `gofmt` clean.

### Summary

Added schema-driven warn-tier field conformance to `spacedock status --validate`: per-field type/pattern/enum checks (mod-block, verdict, score, ISO-8601 dates) are driven from the embedded `entity.mdschema.yml` (`go:embed` at module root — Go forbids `..` in embed patterns, so the canonical SSOT is embedded directly, no drift-prone copy). `validateWorkflow` returns `(errs, warns)`; warns are opt-in to the explicit `--validate` command, print to stderr with a `Warning:` prefix, and never flip the exit code or gate the read path — exit 1 stays reserved for the existing structural-error classes. The single re-enumeration bug (double "preferring folder form" warning) was caught by the existing parity goldens and fixed by computing warns inside `validateWorkflow`'s existing single enumeration. AC-1..AC-4 all verified by `field_conformance_test.go` and exercised against the real built binary (warn-only → exit 0 + VALID; default table → not gated; dup-id+warn → exit 1; schema-edit → behavior change with no Go edit). Notable decision: `id`-per-id-style remains in the existing structural (gating, hand-coded) path rather than the schema-driven warn tier, to avoid double-reporting and avoid disturbing the certified id-parity goldens — out of scope per "minimal changes".

## Stage Report: validation

- DONE: Reproduce AC-1..AC-4 with external evidence on the BUILT binary; `go test ./internal/status/` and `go test ./...` green; id validation still fires via the existing structural path.
  All on the built binary (`go build -o /tmp/sd-validate-binary ./cmd/spacedock`):
  AC-1 — `--validate` over a fixture with `mod-block: noColonHere`, `score: notanumber`, `verdict: MAYBE`, `started: not-a-date` emitted a field-named `Warning:` line per defect on stderr (`Warning: field 'mod-block' value "noColonHere" does not match pattern ^[^:]+:[^:]+$ …`, etc.).
  AC-2 — warn-only fixture: exit 0, stdout `VALID`, warns on stderr; dup-id+warn fixture: exit 1 with both the `Warning:` line and `Error: duplicate id` lines.
  AC-3 — default table / `--next` / `--boot` over a warn-violation fixture: all exit 0 with output on stdout and EMPTY stderr (read path drops warns entirely — no FO lockout).
  AC-4 — edited ONLY the embedded schema (dropped `PASSED` from `verdict.conventional`), rebuilt with `go-files-changed=0`, and `verdict: PASSED` flipped from accepted (0 warns) to `Warning: field 'verdict' value "PASSED" is not one of [REJECTED]`; restored the schema, worktree clean, behavior reverted.
  `go test ./internal/status/ -run Conformance` = 4 tests / 7 subtests PASS; whole `internal/status` package green; `go test ./...` green (16 pkgs ok, 2 no-test).
  id-per-id-style still fires via the structural (gating) path: non-numeric id → `Error: non-numeric sequential id` exit 1 (gates `--validate` AND the read path); missing id → `Error: missing required id` exit 1. AC-1's id requirement is satisfied; the warn tier covers only the four schema fields with no prior enforcement.
- DONE: STAFF P3 — warn-tier stderr does NOT perturb the exact-match `stdout=="VALID"` golden in `native_validate_test.go`.
  On the detached audit checkout, `TestNativeValidationParity/valid` (the `stdout==VALID` golden) is green; the `native-validate-valid` golden envelope is exit 0 / stdout `VALID` / EMPTY stderr — warns are stderr-only and do not leak into the structural validate golden.
- FAILED: Detached adversarial audit (HIGH-STAKES status guard surface) on a SEPARATE throwaway checkout — refute (a)..(d).
  Throwaway: `git worktree add --detach /tmp/sd-audit-checkout <impl-HEAD 76a924a7>` (never the impl worktree; removed after). (a) warn flips exit code → `TestFieldConformanceWarnsDoNotGateExit/warn-only-exits-0` + all `TestFieldConformanceWarnsSurface` RED — refuted. (b) warns gate the read path (`failOnValidationErrors` includes warns) → `TestFieldConformanceWarnsDoNotGateReads` RED — refuted. (c) hardcode a Go regex for `mod-block` instead of reading the schema → `TestFieldConformanceSchemaDriven` RED (`"^.+:.+$" != "^[^:]+:[^:]+$"`) — refuted. (d) drop the `Warning:`/`Error:` prefix split (warns emit `Error:`) → WHOLE PACKAGE STAYS GREEN — NOT refuted. This is a `Material:` test-strength hole.

### Material/Polish audit findings

- **Material: the `Warning:`/`Error:` prefix split on field-conformance lines is untested.** Adversarial edit (d) — `entityEvidenceLine("Warning", …)` → `entityEvidenceLine("Error", …)` in `internal/status/field_conformance.go:104` — makes a warn-tier field finding print with the `Error:` prefix (the FO's gating-vs-advisory signal) while the whole `internal/status` package stays green. The warn-only tests assert only the field NAME (`strings.Contains(nErr, tc.wantField)`) and non-empty stderr, never the `Warning:` prefix; no golden captures a field-conformance warn line. The design treats the prefix as the surfaced-severity mechanism ("surface at the declared severity"), so an `Error:`-prefixed warn silently misreports a non-gating advisory as a structural error — a regression the suite cannot see. The checklist explicitly named edit (d) as one a test must red; it does not. Fix (route through implementation): in `TestFieldConformanceWarnsSurface` assert each line starts with `Warning:` (e.g. `strings.Contains(nErr, "Warning: field '"+tc.wantField)`), or add a golden envelope over a warn-only `--validate` whose stderr pins the `Warning:` prefix. No production-code change needed — the behavior is correct; the test is missing.

### Summary

REJECTED. AC-1..AC-4 all reproduce with external evidence on the BUILT binary, the full package and whole module are green, id validation still fires via the structural path, and the STAFF P3 golden holds (warns are stderr-only, `stdout==VALID` unperturbed). Three of four adversarial refutations land — (a) exit-code, (b) read-path gating, (c) schema-driven agreement each turn a test red. But refutation (d) — drop the `Warning:`/`Error:` prefix split — leaves the whole package green, a `Material:` test-strength hole the checklist explicitly required a test to catch. The fix is test-only (assert the `Warning:` prefix on field-conformance lines); production behavior is correct. Routing back through the validation→implementation feedback flow to close the gap before the gate is presented clean.

### Feedback Cycles

- **Cycle 1 — validation REJECTED → implementation (2026-06-13).** Material (detached adversarial audit, refutation (d)): the `Warning:`/`Error:` prefix split on field-conformance lines is untested — flipping warns to `Error:` (`internal/status/field_conformance.go:104`) leaves the whole `internal/status` package green, so a non-gating advisory misreported as a structural `Error:` would ship unseen. Fix is TEST-ONLY (production behavior is correct): in `TestFieldConformanceWarnsSurface` assert each field-conformance line starts with `Warning: field '<name>` (e.g. `strings.Contains(nErr, "Warning: field '"+tc.wantField)`), or add a warn-only golden envelope pinning the `Warning:` prefix; then confirm the prefix-drop edit (d) now reds. Routed to the live implementation worker (reuse_ok, 11.3% ctx); re-validate after.

## Stage Report: implementation (cycle 1 — Warning: prefix test)

- DONE: Assert each field-conformance warn line carries the `Warning:` prefix bound to the field, and prove edit (d) reds it.
  `TestFieldConformanceWarnsSurface` now asserts `strings.Contains(nErr, "Warning: field '"+tc.wantField)` AND rejects `"Error: field '"+tc.wantField`. Proven to bite: applied edit (d) (`entityEvidenceLine("Warning", …)` → `("Error", …)` at field_conformance.go:104) — all four subtests RED with `stderr missing warn line "Warning: field '…"` (got the `Error:` line); restored production code → green. Test-only change; production code byte-identical to cycle-0 commit 76a924a7.
- DONE: `go test ./internal/status/` green over the whole package; `go vet` + `gofmt` clean; committed on worktree branch.
  Whole `internal/status` package green after restore; `go vet` clean; `gofmt` clean. Test-only commit 66d67387 on `spacedock-ensign/mdschema-conformance-validator`.

### Summary

Closed the test-only gap from the detached audit's refutation (d): the warn-only tests asserted only the field name, so flipping field-conformance warns to a gating `Error:` prefix left the package green. `TestFieldConformanceWarnsSurface` now pins the `Warning:` prefix bound to each field and rejects an `Error:` line for that field — verified to red under edit (d) and green when restored. No production change (the `Warning:`/`Error:` split was already correct); commit 66d67387.

## Stage Report: validation (cycle 2)

- DONE: PRIMARY — re-run detached refutation (d) and confirm `TestFieldConformanceWarnsSurface` now REDS (the gate to closing the REJECTED).
  Throwaway `git worktree add --detach /tmp/sd-audit-c2 <cycle-1 HEAD 66d67387>` (removed after). Baseline green; applied edit (d) (`entityEvidenceLine("Warning", …)` → `("Error", …)` at field_conformance.go:104) → all four subtests RED (`stderr missing warn line "Warning: field 'mod-block"`, …; and the new line-71 guard rejects an `Error: field '…` line). The previously-unrefuted edit now bites. Restored → green.
- DONE: Re-confirm AC-1..AC-4 still reproduce + STAFF P3 golden holds; spot-check refutations (a)(b)(c) still red.
  Cycle-1 fix is TEST-ONLY: `git diff 76a924a7 66d67387 -- ':!*_test.go'` is EMPTY (production byte-identical to what cycle 1 validated; only field_conformance_test.go +10/-2). Built `/tmp/sd-c2`: AC-1+AC-2 warn-only → exit 0, stdout VALID, 4 `Warning: field` lines; AC-2 structural dup+warn → exit 1 with 2 `Error: duplicate id`; AC-3 default table over warn fixture → exit 0, table on stdout, EMPTY stderr. AC-4 `TestFieldConformanceSchemaDriven` green (live schema-edit→rebuild behavior change proven in cycle 1; production unchanged). P3: `TestNativeValidationParity` green, `native-validate-valid` golden = exit 0 / stdout VALID / empty stderr (warns stderr-only, unperturbed). Spot-checks: (a) warn flips exit → AC-2 RED; (b) warns gate read path → AC-3 RED; (c) hardcoded Go regex → AC-4 RED. All four breaking edits now bite.
- DONE: `go test ./internal/status/` green over the whole package; `go test ./...` green.
  Whole `internal/status` package green; `go test ./...` green (16 pkgs ok, 2 no-test). Audit checkout removed; implementation worktree pristine (`git diff --quiet` exit 0).

### Material/Polish audit findings (cycle 2)

- Cycle-1 Material finding CLOSED: the `Warning:`/`Error:` prefix split is now pinned — `TestFieldConformanceWarnsSurface` asserts each line carries `Warning: field '<name>` AND rejects `Error: field '<name>`, so a non-gating advisory misreported as a structural `Error:` reds (verified: edit (d) now reds all four subtests). No new material findings; refuted nothing further material.

### Summary

PASSED. The cycle-1 fix closes the sole Material finding: edit (d) — the previously-unrefuted "drop the Warning:/Error: prefix split" — now reds all four subtests of `TestFieldConformanceWarnsSurface` on a detached checkout. The fix is test-only (`git diff … -- ':!*_test.go'` empty), so AC-1..AC-4 and the STAFF P3 golden carry forward and were re-confirmed on the freshly-built cycle-2 binary; all four adversarial edits (a)(b)(c)(d) now bite. Whole package and whole module green. Gate is clean.

## Stage Report: implementation (cycle 1 — rebase onto main)

- DONE: Rebase `spacedock-ensign/mdschema-conformance-validator` onto current `origin/main` (main moved +16; #350 added sandbox-posture tests).
  `git fetch origin main` + `git rebase origin/main` was CLEAN — both my commits replayed (now 3887419e, 1ae74cf5). Main's #350 touched `internal/status/harness_test.go` and added `boot_sandbox_test.go`, neither of which I touched, so no conflict; working tree clean, no new commit needed (a clean rebase replays existing commits).
- DONE: Re-test against the integrated package — AC-1..AC-4 + STAFF P3 `stdout==VALID` golden hold alongside main's new tests; whole module green.
  `go test ./internal/status/` green over the WHOLE package (incl. `boot_sandbox_test.go`); `TestFieldConformanceWarnsSurface/...DoNotGateExit/...DoNotGateReads/...SchemaDriven` and `TestNativeValidationParity` (the VALID golden) all PASS; `go test ./...` green; `go vet` + `gofmt` clean.

### Summary

Merge-prep rebase onto current `origin/main` (+16, including #350's sandbox-posture test additions) was clean — no overlap with my files, both commits replayed without conflict. Re-tested against the integrated `internal/status` package: AC-1..AC-4, the STAFF P3 `stdout=="VALID"` golden, and main's new `boot_sandbox_test.go` all pass together; whole module green; vet + gofmt clean. Branch is ready for the PR to `main` (not pushed — no force-push needed).
