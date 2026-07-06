---
id: tvstbznw83y5vwgc0jhr6ss8
title: "status --read --checklist/--ac-scan default --stage to the entity's current stage (drop the required-flag round-trip)"
status: implementation
sprint: 0250-fo-behavioral-discipline
source: "FO session 2026-07-04, boot-friction #3: status --read <entity> --checklist and --ac-scan both error 'requires --stage <stage>' when --stage is omitted (reproduced against entity 72, current stage validation), forcing a round-trip to learn a stage the entity's own status field already names. Distinct code path from 3t (--where robustness) and fk (--read frontmatter projection). Drafted by the session's science-officer teammate; sibling in the 0203-fo-efficiency sprint."
started: 2026-07-06T16:17:13Z
completed:
verdict:
score: 0.3
worktree: .worktrees/spacedock-ensign-status-checklist-acscan-default-stage
issue:
---

`status --read <entity> --checklist` and `--read <entity> --ac-scan` hard-require an explicit `--stage` even though the entity's own `status` frontmatter field already names its current stage, forcing a two-call sequence at gate assembly — one of the hottest FO read paths.

## Problem

`spacedock status --read <entity> --checklist` and `--read <entity> --ac-scan` exit 1 with `--checklist requires --stage <stage>` / `--ac-scan requires --stage <stage>` when `--stage` is omitted, even though the entity's frontmatter `status:` field already names the stage the caller almost always wants — the current one. At gate assembly this forces a two-call round-trip: first `status --read <entity> --fields status` to learn the stage, then re-issue with `--stage <that>`.

Reproduced 2026-07-07 against `docs/dev/.spacedock-state/fo-self-evidence-bar.md` (current `status: ideation`): bare `--checklist` → `Error: --checklist requires --stage <stage>` (exit 1); `--fields status --json` → `"status":"ideation"`; `--stage ideation --checklist` → the checklist. Two calls to read the current stage's roll-up.

## Approach

When `--stage` is omitted for `--checklist` / `--ac-scan`, default it to the resolved entity's frontmatter `status` field (its current stage). Explicit `--stage <name>` is unchanged and still selects any non-current stage's report/ACs. This is a read-only projection change in `internal/status` (the `--read` gate path in `gate_extract.go` / `native_runner.go`); no mutation path and no `--set`/finalize/archive guard path is touched.

### Where the default resolves

`runReadGate` (gate_extract.go) already resolves the `--read` ref to a path and reads its bytes into `data`. After that read, when `stage == ""`, set `stage = parseFrontmatterContent(data)["status"]` — the existing in-memory frontmatter parser (frontmatter.go), reused so the file is read once. A non-empty status flows into the existing `runChecklist` / `runACScan` unchanged. The default is computed once, before mode dispatch, so it applies to `--checklist`, `--ac-scan`, text output, and `--json` uniformly.

### Fail-loud semantics

Two distinct empty cases, both loud (exit 1, no silent emit):

1. **Current stage has no matching Stage Report.** The defaulted status flows into the existing `selectStageReport`, which returns `ok=false` → `no ## Stage Report for stage %q in this file` (exit 1). Because the caller never typed the stage, the diagnostic when the stage was defaulted names it as the current status so the defaulting is transparent, e.g. `no ## Stage Report for stage "ideation" (current status; --stage omitted) in this file`. The loud path itself is proven — verified 2026-07-07: `--stage validation --checklist` on an entity without a validation report → `Error: no ## Stage Report for stage "validation" in this file` (exit 1).
2. **Target has no `status` field to default from.** The resolved `--read` target may be a file with no `status` frontmatter — the workflow README, or a plain report file (both are valid `--read` targets). The default is then empty and the command exits 1 with a diagnostic that names the situation and points at `--stage`, e.g. `--stage omitted and <path> has no status frontmatter to default from; pass --stage <stage>`. It must NOT reuse the bare `requires --stage` wording (misleading now that the flag defaults) and must NOT silently emit. Verified 2026-07-07: `--read docs/dev/README.md --fields status --json` → `"status":""`, so the empty-default case is reachable and detectable.

## Out of scope

- The `--fields` frontmatter projection path (`fk` sibling) and `--where` robustness (`3t` sibling) — distinct code paths, unchanged here.
- Any mutation / `--set` / finalize / archive behavior.
- The checklist / ac-scan extraction and AC-citation logic itself.

## Acceptance criteria

**AC-1 (value — the round-trip collapses from two `status --read` calls to one, measured against the origin/main baseline).** Obtaining an entity's current-stage checklist (or ac-scan) costs ONE `status --read` call. Metric: the number of `status --read` invocations to get the current stage's roll-up. Baseline (`origin/main`): the bare `status --read <entity> --checklist` (no `--stage`) exits 1 (`--checklist requires --stage <stage>`), so the caller must first run `status --read <entity> --fields status` to learn the stage, then re-issue with `--stage <that>` — 2 calls. After: `status --read <entity> --checklist` alone exits 0 and emits the current stage's checklist — 1 call. The count moves the wrong way under regression: reintroducing the required-flag error pushes it back to 2.
- *Verified by:* a Go/behavior test on a fixture entity whose current `status` names a stage that has a `## Stage Report`: (a) the single no-`--stage` `--checklist` call exits 0 and its stdout is the current stage's checklist; (b) a pinned assertion that this single call does NOT emit `requires --stage`. The origin/main side of the baseline is the behavior reproduced above (bare call exits 1); the test encodes the after-state as the regression guard.

**AC-2 (default output equals explicit-current output, both flags, both modes).** For an entity whose current `status` names a stage with a matching `## Stage Report`, `--checklist` and `--ac-scan` with `--stage` omitted produce byte-identical stdout to the same command with `--stage <current-status>` given explicitly, in both text and `--json` modes; the `--json` envelope's `stage` field equals the current status.
- *Verified by:* a fixture test comparing omitted-vs-explicit stdout for {`--checklist`, `--ac-scan`} × {text, `--json`}.

**AC-3 (explicit non-current stage still selects it).** `--stage <name>` with a name other than the current `status` still selects that stage's report/ACs; the default applies only when `--stage` is omitted.
- *Verified by:* a fixture entity carrying Stage Reports for two stages (a current one and a prior one); `--stage <prior> --checklist` returns the prior stage's checklist, unaffected by the current-status default.

**AC-4 (fail-loud: current stage has no Stage Report).** When `--stage` is omitted and the entity's current `status` stage has no matching `## Stage Report`, the command exits 1 with a diagnostic naming the stage as the current/defaulted one — never a silent empty emit.
- *Verified by:* a fixture entity in a stage with no report; bare `--checklist` and bare `--ac-scan` each exit 1, stderr names the stage.

**AC-5 (fail-loud: no status field to default from).** When `--stage` is omitted and the resolved `--read` target has no `status` frontmatter field (e.g. the workflow README or a plain report file), the command exits 1 with a diagnostic that the stage could not be defaulted and `--stage` must be passed — not the bare `requires --stage` wording, and not a silent emit.
- *Verified by:* a fixture test running bare `--checklist` on a no-status file (the workflow README or a frontmatter-less fixture), asserting exit 1 and the named diagnostic.

## Spike / proof determination

No spike needed. The design composes two already-proven mechanisms, both exercised 2026-07-07 during ideation:

1. **Frontmatter `status` read** — `ParseFrontmatter` / `parseFrontmatterContent` is the same parser used for `commissioned-by` (native_runner.go:174) and every `--fields` projection. Exercised: `--read <entity> --fields status --json` returns `"status":"ideation"` for an entity and `"status":""` for the README. The `status` field IS the current-stage field (entity frontmatter carries `status: ideation`).
2. **`selectStageReport` fail-loud** — the existing gate path returns `ok=false` → `no ## Stage Report for stage %q` (exit 1). Exercised: `--stage validation --checklist` on an entity without a validation report → exit 1 with that message.

The new code is a ~5-line read-only projection in `runReadGate`: substitute the current `status` when `--stage` is omitted, plus one empty-status fail-loud branch. No mutation path, no `--set`/finalize/archive guard-path edit. The reproduction above seeds the implementation's first tests (AC-1 / AC-4 / AC-5 fixtures).

## Doc-diff determination

User-visible CLI behavior changes (`--checklist` / `--ac-scan` no longer require `--stage`), so per the ideation stage def the docs-site command reference is updated. One cell in `docs/site/reference/command-reference.md` (the `spacedock status` row) documents these flags.

Before:

> … with `--stage X --checklist` / `--stage X --ac-scan` it extracts a stage report's checklist items with line ranges and per-AC evidence citations for the first officer's gate prep) |

After:

> … with `--checklist` / `--ac-scan` it extracts a stage report's checklist items with line ranges and per-AC evidence citations for the first officer's gate prep; `--stage` defaults to the entity's current `status` when omitted (so a bare `--read <entity> --checklist` reads the current stage's report), and `--stage X` reads a non-current stage) |

## Test plan

- **Level:** Go unit / behavior tests in `internal/status`, driving the native runner / `dispatch` over fixture files under `internal/status/testdata` (the pattern used by `cycle3_extract_test.go`, `section_read_test.go`, and the golden read tests). No live-workflow test — the behavior is deterministic file projection.
- **Fixtures:**
  1. Entity with current `status: <S>` and a matching `## Stage Report: <S>` plus a `## Acceptance criteria` section (AC-1, AC-2).
  2. Reuse (1) with a second `## Stage Report: <prior>` and status set to the later stage (AC-3).
  3. Entity with current `status: <S>` and NO `## Stage Report: <S>` (AC-4).
  4. A no-status file — the workflow README or a frontmatter-less fixture (AC-5).
- **Assertions:** exit codes, exact stderr diagnostics for the two fail-loud cases, and stdout parity between omitted and explicit-current for both flags and both output modes.
- **Cost:** low. Pure Go table tests over in-repo fixtures; no network, no binary build beyond `go test ./internal/status`.

## Related

- Siblings `fk` (frontmatter projection) and `3t` (`--where` robustness) — distinct `--read` code paths.
- Sprint `0250-fo-behavioral-discipline`.

## Stage Report: ideation

- DONE: Design: --checklist and --ac-scan default --stage to the entity's current status frontmatter when omitted, explicit --stage still selects a non-current stage; at least one AC MEASURES the value against a baseline that moves the wrong way (the reproduced two-call round-trip / bare 'requires --stage' error collapses to one call) and fail-loud semantics are defined for a current stage with no matching Stage Report
  Approach + 5 ACs written. AC-1 is the value AC: invocation count to read the current-stage roll-up, 2 (origin/main, bare call exits 1) → 1, regressible. AC-4/AC-5 define the two fail-loud cases (current stage has no report; target has no status field).
- DONE: Spike-or-record per proof policy: exercise the riskiest mechanism first or record 'no spike needed' with the proven mechanisms it relies on (status field read + existing selectStageReport fail-loud path); note the change is read-only projection in internal/status, no mutation/guard-path edits
  Recorded "no spike needed" in the Spike / proof determination section — both mechanisms exercised 2026-07-07: `--read <entity> --fields status --json` returns `"status":"ideation"` (and `""` for README); `--stage validation --checklist` on a report-less stage → exit 1 `no ## Stage Report for stage "validation"`. Change is a ~5-line read-only projection in runReadGate; no mutation/guard-path edit.
- DONE: Doc diff proposed in the body per the ideation stage def: concrete before/after for the docs-site command-reference text documenting the two flags' --stage behavior
  Doc-diff determination section carries the exact before/after for the `spacedock status` cell in docs/site/reference/command-reference.md.

### Summary

Fleshed the task into Problem / Approach / Acceptance criteria / Test plan. The design defaults `--stage` to the entity's frontmatter `status` for `--checklist` / `--ac-scan` when omitted, resolved once in `runReadGate` via the existing `parseFrontmatterContent`; explicit `--stage` is unchanged. Both fail-loud cases (no Stage Report for the current stage; no `status` field to default from) are defined and exit 1. No spike needed: the status-field read and the `selectStageReport` fail-loud path were both exercised against live entities during ideation, and the reproduction seeds the implementation's first tests.
