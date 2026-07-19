---
id: 2690fpqe9pkn917am6bt6eqs
title: Make dispatch-build help match its input-mode parser
status: validation
source: "FO dogfood, 2026-07-19: 0.26.0-pre0 help advertised stdin JSON plus --advance but the same invocation selected flag/file mode and rejected stdin."
started: 2026-07-19T05:04:28Z
completed:
verdict:
score: "0.65"
worktree: .worktrees/spacedock-ensign-dispatch-build-help-input-mode-drift
issue:
milestone: 0.26.0
group: binary-ux
---

Make `spacedock dispatch build --help` describe every supported input form, its
mode-selection rules, and a complete reuse-advance invocation that the parser
actually accepts.

## Problem

The current help prints only the stdin-JSON usage, omits the accepted
`--entity-path`, `--stage`, `--checklist-file`, `--scope-notes-file`, and
`--feedback-context-file` flags, then lists `--advance` without explaining that
it participates in request-flag detection. A caller following that surface can
pipe a complete JSON request and add `--advance`, only to receive:

```text
error: flag/file input requires --entity-path, --stage, and --checklist-file
```

The implementation selects flag/file parsing whenever `hasRequestFlags()` is
true, so stdin is ignored in that case. The help neither exposes that boundary
nor supplies the complete flag/file form required by the First Officer
contract. This is a recurring dispatch-boundary round trip.

The related active task `dispatch-build-flag-form-version-skew` and closed
GitHub issue #313 concern an older accepted binary versus newer plugin
instructions. This task concerns one current binary contradicting its own help.

## Acceptance criteria

- **AC-1:** `dispatch build --help` documents the complete stdin-JSON and
  flag/file forms, including all accepted request flags and the exact rule that
  selects between them. Verified by running the rendered examples against the
  real parser, not by grepping source prose.
- **AC-2:** The supported reuse-advance form is unambiguous: either stdin JSON
  plus `--advance` works, or help explicitly rejects that combination and gives
  a complete accepted flag/file example. Verified by paired success/failure
  command tests asserting stdout, stderr, and exit status.
- **AC-3:** A behavioral help-example test fails whenever an example printed by
  `--help` no longer parses successfully with a minimal workflow/entity fixture.

## Scope

Keep the fix to command help, input-mode selection if necessary, and
load-bearing behavioral tests. Do not redesign the dispatch envelope or add a
third request form.

## Proposed approach

**Direction (a): HELP-ONLY.** Rewrite `printBuildUsage`
(`internal/dispatch/dispatch.go:290`) so `dispatch build --help` documents both
input modes, every request flag, the exact mode-selection rule, and an accepted
reuse-advance (`--advance`) invocation in flag/file form. Do NOT touch the
parser or input-mode selection.

Why (a), not (b): the spike (below) proves the parser already accepts the
reuse-advance form the First Officer contract actually uses — flag/file with
`--advance` (`--entity-path`/`--stage`/`--checklist-file`). The FO never pipes
stdin+`--advance`: `claude-fo-dispatch.md` (Reuse-advance handle) and
`fo-dispatch-core.md` (`## Reuse and Fresh Dispatch`) both build advance as
`dispatch build --advance` with `--workflow-dir --entity-path --stage
--checklist-file`. So the parser is correct and the help is the only thing
wrong. Direction (b) — reclassify `--advance`/`--bare-mode`/`--feedback-reflow`
as mode-neutral modifiers and merge `opts.Advance` into stdin-parsed fields so
stdin+`--advance` is accepted — is strictly larger (ripples through
`hasRequestFlags`, `isBuildRequestFlag`, and every request-flag test), blurs the
input-mode boundary the codebase deliberately keeps disjoint, and satisfies a
combination no caller needs (YAGNI). (a) removes the contradiction with a
one-function text change and matches the shipped FO contract.

**Exact mode-selection rule to document** (read from `parseBuildOptions` +
`loadBuildFields`): presence of ANY request flag selects flag/file mode and
stdin is ignored; otherwise stdin JSON is read. Request flags are
`--entity-path`, `--stage`, `--checklist-file`, `--scope-notes-file`,
`--feedback-context-file`, `--team-name`, `--bare-mode`, `--feedback-reflow`,
`--advance`. Flag/file mode requires `--entity-path`, `--stage`, and
`--checklist-file` (else exit 2, `flag/file input requires --entity-path,
--stage, and --checklist-file`). `--workflow-dir` and `--host` apply to both
modes; `--print-schema` and `--validate-only FILE` are separate modes checked
before the request-flag test.

## Doc diff: `dispatch build --help`

The only user-visible surface is `printBuildUsage`
(`internal/dispatch/dispatch.go:290-313`). No docs-site page mirrors it (the two
`docs/` hits are historical roadmap/debrief archives). The top-level
`printUsage` already shows both forms; this task fixes `build --help`, which
today advertises only stdin.

**BEFORE** (current `printBuildUsage` body, dispatch.go:291-312):

```text
Usage:
  spacedock dispatch build --workflow-dir DIR

Build an ensign dispatch artifact from stdin JSON and write the JSON envelope to stdout.

Flags:
  --workflow-dir DIR   Workflow definition directory containing README.md.
  --host HOST          Override the runtime host (claude|codex|pi). Defaults to the detected runtime.
  --team-name NAME     Select the legacy TeamCreate-registry dispatch shape. On host=claude, auto-team is the default — omit this unless you mean legacy team mode.
  --bare-mode          Emit the bare sequential shape (no name, no team_name, no run_in_background).
  --advance            Emit a reuse-advance pointer message for a live worker instead of a spawn envelope. Incompatible with --bare-mode.

Stdin JSON fields:
  schema_version  Dispatch schema version. The current supported value is 2.
  entity_path     Path to the entity file for this dispatch.
  workflow_dir    Workflow directory for the dispatch request.
  stage           Stage name to dispatch.
  checklist       Array of checklist strings for the dispatched worker.

Example:
  {"schema_version":2,"entity_path":"thing.md","workflow_dir":".","stage":"implementation","checklist":["DONE: run tests"]}
```

**AFTER** (proposed `printBuildUsage` body):

```text
Usage:
  spacedock dispatch build --workflow-dir DIR < request.json                                                   (stdin mode)
  spacedock dispatch build --workflow-dir DIR --entity-path FILE --stage STAGE --checklist-file FILE [flags]   (flag/file mode)
  spacedock dispatch build --print-schema
  spacedock dispatch build --validate-only FILE

Build an ensign dispatch artifact and write the JSON envelope to stdout. The
request comes from a JSON object on stdin OR from flags/files. The two are
selected by the rule below and never merged.

Input mode selection:
  If ANY request flag is present the request is read from flags/files and stdin
  is IGNORED (flag/file mode); otherwise the request is read as a JSON object on
  stdin (stdin JSON mode). Request flags:
    --entity-path  --stage  --checklist-file  --scope-notes-file
    --feedback-context-file  --team-name  --bare-mode  --feedback-reflow  --advance
  Flag/file mode requires --entity-path, --stage, and --checklist-file; any
  request flag with one of the three missing fails:
    error: flag/file input requires --entity-path, --stage, and --checklist-file
  Because --advance is a request flag, piping JSON on stdin together with
  --advance is NOT accepted -- it selects flag/file mode and ignores the piped
  JSON. Pass a reuse-advance request in flag/file form (see the --advance
  example below).

Flags:
  --workflow-dir DIR            Workflow definition directory containing README.md (both modes).
  --host HOST                   Override the runtime host (claude|codex|pi). Defaults to the detected runtime (both modes).
  --entity-path FILE            Entity file for this dispatch (flag/file mode).
  --stage STAGE                 Stage name to dispatch (flag/file mode).
  --checklist-file FILE         File of checklist lines, one per line (flag/file mode).
  --scope-notes-file FILE       Optional scope-notes file (flag/file mode).
  --feedback-context-file FILE  Optional feedback-context file; required with --feedback-reflow (flag/file mode).
  --team-name NAME              Select the legacy TeamCreate-registry dispatch shape. On host=claude, auto-team is the default — omit this unless you mean legacy team mode.
  --bare-mode                   Emit the bare sequential shape (no name, no team_name, no run_in_background).
  --feedback-reflow             Route a rejection back to its feedback-to target stage; requires --feedback-context-file.
  --advance                     Emit a reuse-advance pointer message for a live worker instead of a spawn envelope. Incompatible with --bare-mode.
  --print-schema                Print the stdin request JSON schema and exit.
  --validate-only FILE          Validate a request JSON file without writing a dispatch; exit 0 on success.

Stdin JSON request fields (stdin JSON mode):
  schema_version  Dispatch schema version. The current supported value is 2.
  entity_path     Path to the entity file for this dispatch.
  workflow_dir    Workflow directory for the dispatch request.
  stage           Stage name to dispatch.
  checklist       Array of checklist strings for the dispatched worker.
  (optional: team_name, scope_notes, feedback_context, bare_mode, is_feedback_reflow, advance, host)

Examples:
  stdin JSON mode:
  {"schema_version":2,"entity_path":"thing.md","workflow_dir":".","stage":"implementation","checklist":["DONE: run tests"]}

  flag/file mode:
  spacedock dispatch build --workflow-dir . --entity-path thing.md --stage implementation --checklist-file impl.checklist

  reuse-advance (flag/file mode):
  spacedock dispatch build --workflow-dir . --entity-path thing.md --stage validation --checklist-file validation.checklist --advance
```

Notes on the AFTER text: it preserves the substrings the existing
`TestDispatchBuildHelpBeforeRequiredFlags` asserts (`Usage:`, `spacedock
dispatch build --workflow-dir DIR`, `--workflow-dir`, `schema_version`,
`entity_path`, `workflow_dir`, `stage`, `checklist`) except the bare `stdin
JSON` adjacency — that test's assertions are extended in the test plan to cover
the new content.

## Test plan

All tests live in `internal/dispatch/*_test.go` and drive the real parser
in-process via the existing `runNative(stdin, args...)` harness
(`parity_harness_test.go`) with `readmeWorktree`/`entityFM` fixtures — no
source-prose grep. Estimated cost: low; three Go test functions, no new binary
build, seconds to run.

- **AC-1 (documents both forms + all request flags + the selection rule).**
  Extend `TestDispatchBuildHelpBeforeRequiredFlags` (help_test.go) to assert the
  new help text contains: each of the nine request-flag names, the literal
  selection-rule sentence fragments (`stdin is IGNORED`, `requires
  --entity-path, --stage, and --checklist-file`, `piping JSON on stdin together
  with --advance is NOT accepted`), and the flag/file + reuse-advance example
  command lines. This is the documentation-presence half. The behavioral half is
  the AC-3 help-example test (shared mechanism below) — it runs the rendered
  examples against the real parser, which is what makes AC-1 behavioral rather
  than a prose assertion.

- **AC-2 (reuse-advance form unambiguous: paired success/failure asserting
  stdout+stderr+exit).** New `TestDispatchBuildAdvanceInputMode` with two
  sub-cases against one fixture:
  - success: flag/file `--advance` (`--entity-path/--stage/--checklist-file
    --advance`) exits 0, stdout is a JSON advance envelope
    (`schema_version`/`prompt`/`model`, no `subagent_type`/`name`), stderr empty.
  - failure: stdin JSON piped with `--advance` and no flag/file trio exits 2 with
    stderr exactly `error: flag/file input requires --entity-path, --stage, and
    --checklist-file`. This is the exact contradiction the help now documents as
    rejected, pinned as a behavioral test.

- **AC-3 (help-example test fails when any printed example no longer parses).**
  New `TestDispatchBuildHelpExamplesParse`: render `dispatch build --help`,
  extract every positive example (a line beginning `{` = stdin JSON example; a
  line beginning `spacedock dispatch build` = flag/file example), materialize a
  minimal fixture (README with `implementation`+`validation` non-worktree stages,
  a `thing.md` entity, `impl.checklist`/`validation.checklist`), rewrite ONLY the
  leaf path values to the fixture's absolute paths (flag names, field names, JSON
  shape, stage names, and `--advance` presence stay verbatim — the parser
  absolutizes paths, so this preserves every parse concern), run each through
  `runNative`, and assert exit 0. The test fails if any printed example drops a
  required field, renames a flag, or (the shipped bug) advertises a
  stdin+`--advance` form the parser rejects.

### Mechanism justification

- **New help-example execution test** serves AC-3 (and AC-1's behavioral half).
  Simplest alternative: keep substring/`assertContainsAll` grep over the help
  text (today's `help_test.go` style). Insufficient — grep proves the words
  appear, not that the example parses; that is exactly the drift that shipped
  (help said stdin+`--advance` works while the parser rejected it, and a grep for
  `stdin JSON` passes on the broken help). Only executing the rendered example
  against the real parser catches it.
- **No parser/input-mode change** — the enabling mechanism direction (b) would
  add is unnecessary: the spike proves the value AC-2 targets (an accepted
  reuse-advance form) already exists in flag/file mode.

## Spike: parser reuse-advance form (riskiest unknown, exercised)

Built the real binary (`go build ./cmd/spacedock`) and ran both reuse-advance
forms plus all three proposed help examples against a minimal fixture
(README with non-worktree `implementation`/`validation` stages, `thing.md`
entity), `CLAUDECODE=1`:

- stdin JSON + `--advance` (the advertised form) → **exit 2**, stderr `error:
  flag/file input requires --entity-path, --stage, and --checklist-file`. The
  contradiction, reproduced.
- flag/file + `--advance` (the FO-contract form) → **exit 0**, valid advance
  envelope: no `subagent_type`/`name`/`team_name`, `model: null`, dispatch file
  suffixed `-advance`, `prompt` the reuse-advance pointer.
- All three proposed help examples (stdin JSON `implementation`; flag/file
  `implementation`; reuse-advance flag/file `validation`) → **exit 0**.

FO contract cross-check: `skills/first-officer/references/claude-fo-dispatch.md`
(Reuse-advance handle) and `fo-dispatch-core.md` (`## Reuse and Fresh Dispatch`)
both invoke advance as flag/file `--advance` — no stdin pipe. The direction the
chosen fix documents is the direction the contract exercises. Spike settled:
direction (a), parser unchanged.

## Stage Report: ideation

- DONE: A Proposed approach that resolves the help↔parser contradiction and names WHICH fix — direction (a) HELP-ONLY.
  `## Proposed approach` names (a), documents the exact `hasRequestFlags()` mode-selection rule + all request flags, and justifies rejecting (b) as larger/contract-mismatched.
- DONE: A Test plan that proves AC-1/AC-2/AC-3 behaviorally.
  `## Test plan`: AC-1 help-text presence + shared example-execution; AC-2 paired flag/file-advance success vs stdin+advance failure asserting stdout/stderr/exit; AC-3 help-example test that runs every rendered example through the real parser.
- DONE: Spike the riskiest unknown before the gate — what the parser accepts for reuse-advance and what the FO contract requires.
  `## Spike`: built binary, stdin+`--advance` → exit 2, flag/file `--advance` → exit 0; FO refs invoke flag/file advance. Settled direction (a).

### Summary

Chose direction (a) HELP-ONLY: the parser already accepts the reuse-advance form the FO contract uses (flag/file `--advance`), so only `printBuildUsage` is wrong. Recorded a concrete before/after doc diff for `dispatch build --help` documenting both input modes, all nine request flags, the exact stdin-vs-flag/file selection rule, and an accepted `--advance` example. The spike reproduced the shipped contradiction (stdin+`--advance` → exit 2) and confirmed all three proposed examples parse (exit 0), de-risking the AC-3 help-example test. No parser change; ACs left as fixed.

## Stage Report: implementation

- DONE: printBuildUsage (internal/dispatch/dispatch.go) rewritten so `dispatch build --help` documents BOTH input modes, ALL request flags, the EXACT stdin-vs-flag/file selection rule, and a COMPLETE reuse-advance (--advance) invocation in flag/file form. Parser / input-mode selection UNCHANGED (direction a).
  Commit 6b4d0156. BEFORE: only the stdin-JSON usage + a single stdin `Example:` (omitted the flag/file trio, `--scope-notes-file`/`--feedback-context-file`/`--feedback-reflow`, and any selection rule). AFTER: two usage lines (stdin mode / flag/file mode), an "Input mode selection" block (any request flag → flag/file, `stdin is IGNORED`; trio-required + exact `error: flag/file input requires ...`; `--advance is NOT accepted` on stdin), a Flags block listing all nine request flags, and three labeled examples. Only `printBuildUsage` changed; `parseBuildOptions`/`loadBuildFields`/`hasRequestFlags`/`isBuildRequestFlag` untouched (git diff scope = dispatch.go help func + two test files).
- DONE: Behavioral tests — AC-2 paired success/failure asserting stdout+stderr+exit; AC-3 help-example test that runs EVERY rendered example through the REAL parser and FAILS if any no longer parses.
  `TestDispatchBuildAdvanceInputMode` (build_input_mode_test.go): flag/file `--advance` → exit 0, advance envelope (schema_version/prompt/model, no subagent_type/name), empty stderr; stdin JSON + `--advance` → exit 2, empty stdout, stderr trimmed == `error: flag/file input requires --entity-path, --stage, and --checklist-file`. `TestDispatchBuildHelpExamplesParse`: extracts every example from the rendered Examples section, rewrites ONLY leaf path values to a minimal non-worktree impl+validation fixture, runs each through `runNative` (the real parser), asserts exit 0. AC-1 extended `TestDispatchBuildHelpBeforeRequiredFlags` (help_test.go) to assert the nine flag names, the selection-rule fragments, and the two example command lines.
- DONE: Verify — gofmt/vet clean; tests + -race green; AC-3 genuinely executes the rendered examples (would have caught the shipped drift).
  `gofmt -l` clean on all three files; `go vet ./internal/dispatch/` clean; `go test -count=1 ./internal/dispatch/` and `go test -race -count=1 ./internal/dispatch/` green. Drift proof: temporarily replacing the reuse-advance example with the broken stdin+advance form (`dispatch build --workflow-dir . --advance`) made AC-3 FAIL with exit 2 / `error: flag/file input requires --entity-path, --stage, and --checklist-file`, then reverted — the shipped bug is now caught. `go test ./...` is green except `internal/cli`'s `TestCodexResolveManifestAgainstInstalledHost`, a pre-existing environmental failure (the installed Codex host cannot read `~/.codex/config.toml`, `Operation not permitted`; fails identically with the sandbox disabled and touches no dispatch code).

### Summary

Direction (a) HELP-ONLY, as ideated: rewrote only `printBuildUsage` so `dispatch build --help` documents both input modes, all nine request flags, the exact mode-selection rule, and an accepted flag/file reuse-advance example; the parser and input-mode selection are unchanged. Added AC-2 (paired advance success/failure asserting stdout+stderr+exit) and AC-3 (a help-example test that executes every rendered example against the real parser), and extended AC-1's help-presence test. The AC-3 test was proven to catch the shipped drift by a temporary break. All dispatch tests and -race are green; the sole `./...` red is a pre-existing environmental Codex-host config-read failure unrelated to this change.

## Stage Report: validation

- DONE: AC-1 — help documents BOTH input modes, all nine request flags, and the EXACT selection rule; presence-tested in help_test.go, behavior proven by AC-3.
  `TestDispatchBuildHelpBeforeRequiredFlags` PASS: asserts the nine flag names, the three rule fragments (`is IGNORED (flag/file mode)`, `requires --entity-path, --stage, and --checklist-file`, `--advance is NOT accepted`), and both example command lines against RENDERED stdout (`runNative("","build","--help")` → exit 0, empty stderr). Not self-referential: AC-3 executes those same rendered examples against the real parser (independent oracle).
- DONE: AC-2 — reuse-advance form unambiguous; paired success/failure asserting stdout+stderr+exit.
  `TestDispatchBuildAdvanceInputMode` PASS: flag/file `--advance` → exit 0, envelope has schema_version/prompt/model, no subagent_type/name, empty stderr; stdin JSON + `--advance` → exit 2, empty stdout, stderr trimmed == `error: flag/file input requires --entity-path, --stage, and --checklist-file`.
- DONE: AC-3 + direction-(a) integrity — help-example test EXECUTES every rendered example through the real parser; drift-catch reproduced independently; parser UNCHANGED.
  `TestDispatchBuildHelpExamplesParse` PASS (renders help, extracts 3 examples incl. ≥1 `--advance`, remaps only leaf paths, runs each via `runNative`, asserts exit 0). Independent drift-catch: replaced the reuse-advance example with the broken stdin+advance form `spacedock dispatch build --workflow-dir . --advance` → test FAILED exit=2, stderr `error: flag/file input requires --entity-path, --stage, and --checklist-file`; restored via `git checkout` (empty diff vs 6b4d0156, clean tree). Direction (a): dispatch.go diff is a single hunk entirely inside `printBuildUsage`; `parseBuildOptions` (dispatch.go:131), `isBuildRequestFlag` (:185), `hasRequestFlags` (:127), `loadBuildFields` (build.go:125) all outside the hunk and byte-unchanged (build.go absent from the commit stat; each name appears 0× in the diff). No third request form, no envelope redesign. `gofmt -l` clean, `go vet` clean; `go test` + `-race` green on internal/dispatch; `go test ./...` green except the named environmental `internal/cli` TestCodexResolveManifestAgainstInstalledHost.

### Summary

VERDICT: PASSED. Independently reviewed commit 6b4d0156 against AC-1/AC-2/AC-3. All three ACs have valid, reproduced evidence. AC-3 is genuinely load-bearing: I broke a rendered example and confirmed the help-example test fails with exit 2 and the exact trio-required error, then restored. Direction-(a) integrity holds — only `printBuildUsage` changed; the four parser/input-mode functions are byte-unchanged, no third request form, no envelope redesign. gofmt/vet clean; dispatch tests + `-race` green; `go test ./...` green except the declared pre-existing environmental `internal/cli` Codex-host config-read failure (denied read of `~/.codex/config.toml`, touches no dispatch code). No material findings; no deferred risks.
