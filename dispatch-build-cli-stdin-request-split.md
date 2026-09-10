---
title: Separate dispatch invocation controls from stdin assignment data
status: validation
source: Captain discussion of checklist transport and dispatch assembly, 2026-09-10
started: 2026-09-10T20:26:39Z
completed:
verdict:
score: 0.5
worktree: .worktrees/spacedock-ensign-dispatch-build-cli-stdin-request-split
issue:
pr:
id: 09ptxgq8zma6qnx0h1w8a6wp
gates:
    version: 1
    records:
        - id: gate:09ptxgq8zma6qnx0h1w8a6wp:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:09ptxgq8zma6qnx0h1w8a6wp-backlog-1
              briefing:
                id: briefing:09ptxgq8zma6qnx0h1w8a6wp:backlog:attempt-1:revision-1
                digest: sha256:cb3ae42ccb7d64430aef08321f5cb896c7ad8566738a23db978598647dde5adc
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:09ptxgq8zma6qnx0h1w8a6wp:backlog:1
                briefing: briefing:09ptxgq8zma6qnx0h1w8a6wp:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-09-10T20:26:00.677523Z"
                decision: approve
                reason: 'Captain requested: assume this behavior already works in this session. now dispatch it. Authorization is to begin ideation, not implementation.'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:09ptxgq8zma6qnx0h1w8a6wp:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:09ptxgq8zma6qnx0h1w8a6wp-ideation-1
              briefing:
                id: briefing:09ptxgq8zma6qnx0h1w8a6wp:ideation:attempt-1:revision-1
                digest: sha256:8eb17db161d3dd3b9454cf66e8b8f3a5babb78d424bcdfe452ff76192f15d666
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:09ptxgq8zma6qnx0h1w8a6wp:ideation:1
                briefing: briefing:09ptxgq8zma6qnx0h1w8a6wp:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-09-10T20:44:53.720381Z"
                decision: revise
                reason: 'Captain approved a narrower replacement: keep parameters in CLI flags; support --checklist-file - for newline-delimited checklist stdin; preserve scope/feedback file options; retire legacy JSON-stdin parsing and request-schema controls; migrate JSON callers/tests. This explicitly supersedes the legacy JSON compatibility requirement and assignment-stdin proposal. Return revised ideation before implementation.'
            - id: gate-attempt:09ptxgq8zma6qnx0h1w8a6wp-ideation-2
              briefing:
                id: briefing:09ptxgq8zma6qnx0h1w8a6wp:ideation:attempt-2:revision-1
                digest: sha256:d4e1106ddcaaddb3ace907389610f5b27a861fbf68acb9de46ace9932d391dec
                room-ref: '@review/ideation/briefing-2'
              resolution:
                type: Resolution
                id: resolution:spacedock:09ptxgq8zma6qnx0h1w8a6wp:ideation:2
                briefing: briefing:09ptxgq8zma6qnx0h1w8a6wp:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-09-10T21:03:13.639109Z"
                decision: approve
                reason: 'Captain approved ideation attempt 2: checklist stdin via --checklist-file -, retirement of JSON request input and schema controls, caller/test migration, and FO default update within the stated scope.'
              application:
                target-stage: implementation
                state: consumed
---

Simplify normal FO dispatch to one CLI flag/file interface. The checklist can arrive on stdin through `--checklist-file -`; JSON requests and their schema/validation controls are retired. This cycle supersedes the cycle-1 design; implementation remains gated.

## Problem and captain decision

The helper only needs checklist text during assembly, yet the normal FO contract requires a separate checklist file. Keeping both JSON and flag/file requests creates mode selection, duplicate control locations, schema validation, and migration work for callers without improving this ordinary transaction.

Captain reset: “why do we want those argv in json as well? wouldn't it be simplify for checklist to be from stdin?”; “yeah and we can retire the json stdin code path”; “ok. send it back”. This explicitly replaces the prior JSON-compatibility requirement. Preserve supported flag/file callers and stamp/advance semantics, not the retired JSON API. Removing one scratch file is a small simplification, not a latency claim.

## Selected contract

There is exactly one input path. Required assignment identifiers and invocation controls are CLI flags. Checklist source is the existing file option with one conventional sentinel:

```text
spacedock dispatch build --workflow-dir DIR --entity-path FILE --stage STAGE
  --checklist-file FILE|- [--scope-notes-file FILE] [--feedback-context-file FILE]
  [--host claude|codex|pi] [--bare-mode] [--feedback-reflow] [--stamp | --advance]
```

`--checklist-file -` reads the command's supplied stdin reader until EOF; it does not open a file named `-` and does not read global os.Stdin behind the injectable reader. A literal file named `-` remains addressable as `./-`. Any other checklist path retains ordinary file behavior and does not consume stdin. Scope/feedback remain file options; `-` has no new special meaning for them. There is no autodetection, terminal probe, JSON request, new mode selector, or second assembly pipeline.

Both checklist sources call the same byte-to-lines parser: split on LF, remove one terminal CR from each line, ignore lines whose TrimSpace is empty, preserve all other bytes including indentation and trailing spaces. Each retained line is one checklist item. Final LF and no final LF yield the same last item. Blank lines and checklist line-ending style are intentionally not preserved. Scope/feedback keep existing whole-file opaque transport for multiline Markdown and trailing newlines. Do not introduce bufio.Scanner's default token limit: a 70 KB line must work as it does today.

The smallest code change reads bytes from the selected source, then uses the existing line parser and `fieldsFromBuildFlags`/stamp/assembly flow. Delete JSON input loading, its mode-selection aggregate, request-schema emission, request-only validation branches and unreachable request-version checks. The output JSON envelope and `schema_version: 2` remain unchanged. A private normalized field map may remain if used by stamp/assembly; do not expand this task into a typed-model rewrite merely because it currently uses json.RawMessage internally. Runtime host comes from --host or existing detection; JSON host/control precedence disappears with the request API.

Intentional breaking-change diagnostics:

- `--print-schema` and `--validate-only` (including equals forms) exit 2 before stdin/file access: `error: dispatch build JSON request input and --print-schema/--validate-only are retired; use --entity-path, --stage, and --checklist-file FILE|-`.
- No complete required assignment flag trio: preserve the existing first stderr line `error: flag/file input requires --entity-path, --stage, and --checklist-file`, then append `JSON request input is retired; pass checklist lines on stdin with --checklist-file -`. Do not read/sniff stdin to distinguish old JSON from an incomplete CLI call. Missing --workflow-dir retains its existing earlier diagnostic.
- With a complete CLI and --checklist-file -, JSON-looking text is ordinary checklist text. Do not add JSON detection to reject a valid literal checklist line. A former JSON-only invocation has no required assignment trio and is rejected before read.
- Stdin read error exits 1 with `error: failed to read checklist stdin: <cause>` and no envelope or stamp. File-read diagnostics remain unchanged. Empty or whitespace-only checklist preserves the current exit 1 and `error: missing required field 'checklist'` diagnostic (the existing empty parser result is nil). Do not silently turn an empty input into a dispatch.

`--help` stays side-effect free. Existing CLI exclusions (--advance with --bare-mode, --stamp with --advance) remain before checklist reads; required workflow/trio checks precede source reads. Retired-switch refusal occurs in option parsing. After successful loading, preserve all existing downstream validation and stamp ordering: do not promise broad new validation-before-stamp behavior.

`--stamp` still requires status already equal to --stage, stamps started/worktree, commits/synchronizes state, creates the stage worktree when needed, then emits the helper-owned artifact/envelope. It never advances stage or spawns a worker. Existing host restrictions, split-root state ownership, opaque feedback transport, pointer bootstrap and fresh/advance output shapes remain unchanged.

## Acceptance criteria

**AC-1 — Normal FO dispatch needs zero checklist scratch files.** The shared FO command prefers --checklist-file - and feeds literal checklist text to stdin. A real-command fixture executes that documented shape and gets the helper-owned artifact/envelope with zero checklist input files, compared with the current file-mode fixture's one required checklist file. The existing live full-ensign-cycle journey reaches durable done state and commits using the revised skill; its captured dispatch evidence shows the stdin checklist source without a separately written checklist file. Scope/feedback files are still permitted. Requiring a checklist scratch file again fails this value check.

**AC-2 — Checklist and supporting prose retain the surviving format's semantics.** File and stdin checklist sources produce identical artifact sections for quotes, backticks, dollar expressions, Unicode, indentation, trailing spaces, CRLF, blank lines, EOF without LF, and long lines. Expected checklist bytes are independently specified after the documented line normalization. Scope/feedback files preserve multiline Markdown and trailing newlines. Shell sentinel remains absent. Trimming content, expanding shell expressions, adding a scanner token limit, or treating stdin lines differently fails this AC. Exact JSON-array element boundaries and empty checklist lines are intentionally not promised by the retired API.

**AC-3 — The single supported input interface replaces JSON without losing dispatch behavior.** All active in-tree callers use flags plus checklist file/stdin; retired JSON/schema-validation entry points fail with actionable diagnostics and no reads/mutation. Existing supported flag/file dispatches retain success behavior, host resolution, artifacts, and stamp/advance state effects. Existing byte-hazard, split-root/worktree, durable commit/sync ordering, retry, and advance-envelope assertions remain exercised through the surviving input path. The implementation removes the JSON reader/selector/schema machinery rather than hiding it behind an unused adapter.

## Bounded seam evidence (cycle 2)

A temporary Git workflow with one non-worktree backlog entity was driven through installed 0.28.0-pre2. Six input cases each ran once with a regular --checklist-file path and once with --checklist-file /dev/stdin, the inherited descriptor fed through Python subprocess input bytes. Twelve actual builds confirmed identical file/stdin outcomes: EOF without LF, LF, CRLF plus blank lines, empty, whitespace-only, and a 70,000-character line. Nonempty cases compared complete generated artifact bodies and independently expected checklist sections; empty cases compared exit 1 and exact stderr. The hazard line was:

```python
h = '  Keep "quotes", `ticks`, $(touch SHOULD_NOT_EXIST), $HOME, 雪 & <tag>  '
```

The CRLF case was `'\r\n'+h+'\r\n\t \r\nsecond\r\n'`, expected `h+'\nsecond'`. The shell sentinel was absent. Temporary inputs/artifacts were removed. Initial empty-input expectation (`checklist must not be empty`) was falsified; actual behavior is `missing required field 'checklist'`, now reflected in the contract. This is evidence for the shared read-to-EOF/line-parser seam, not evidence that dash routing is already implemented. /dev/stdin was a declared one-off macOS probe aid; implementation and durable Go tests use io.Reader and never depend on that path.

Supporting check passed in both packages: `go test ./internal/dispatch ./skills/integration -run 'TestBuildFlagFileInputModePreservesLiteralChecklist|TestDispatchBuildAdvanceInputMode|TestSplitRootFolderWorktreeDispatch|TestFlatEntitySlugUnchanged' -count=1`. These are current-behavior baselines, not implementation proof. Existing focused dispatch and skill integration tests are the baseline: TestBuildFlagFileInputModePreservesLiteralChecklist, TestDispatchBuildAdvanceInputMode, TestSplitRootFolderWorktreeDispatch, and TestFlatEntitySlugUnchanged. The old JSON-plus-advance rejection assertion must be deliberately rewritten around JSON retirement and positive stdin-checklist-plus-advance, not preserved as a second mode.

## Migration and proof plan

Add focused failing tests before implementation and before changing FO command wording. No new broad harness, hidden JSON conversion in runNative, or generic compatibility shim.

| AC | Existing proof owner and intended change | Distinct falsifying edit | Cost/type |
| --- | --- | --- | --- |
| AC-1 command | internal/dispatch/build_input_mode_test.go: execute the documented stdin checklist example with the real parser and count checklist scratch inputs (zero vs file baseline one); keep help examples executable | Make --checklist-file - open a path or require a scratch file | Small deterministic fixture, seconds |
| AC-1 skill | skills/integration/dispatch_test.go, especially TestSplitRootFolderWorktreeDispatch and TestFlatEntitySlugUnchanged: migrate actual calls and add the FO stdin transport assertion before prose edits | Retain JSON in the integrated caller, or misroute its state/worktree path | Existing skill smoke owner, seconds |
| AC-1 live value | internal/ensigncycle/shared_live_runner_test.go TestLiveCommonFullEnsignCycle and shared_promoted_live_test.go runFullEnsignCycleJourney: revised FO core, existing transcript artifacts plus durable done/report/path-scoped commit proof; check stdin source and absence of a separate checklist write in existing captured evidence | FO ignores new default and writes checklist scratch file, or fails worker completion | One existing Codex lane journey, roughly 5–15 minutes; no new lane/driver |
| AC-2 | build_json_ergonomics_test.go's literal-checklist test (retain filename or rename only if useful), build_hazards_test.go: file/stdin table seeded by the spike, exact normalized section bytes; scope/feedback file byte assertions; execute a quoted-heredoc hazard once through existing command fixture | Strip indentation/trailing space; Scanner rejects 70 KB; unquote heredoc delimiter | Small deterministic tables and one shell transport exercise |
| AC-3 CLI | build_input_mode_test.go and build_errors_test.go: retired switches including equals, no-trio JSON caller, missing args, empty input, reader failure, regular path with a poison stdin reader; explicit no envelope/no mutation on input errors | Keep legacy JSON fallback, read stdin in file mode, accept retired switch | Existing parser/command owner, seconds |
| AC-3 lifecycle | build_stamp_test.go and build_advance_test.go: drive existing success/mismatch, incompatibility, retry-sync, inline-before-worktree and fresh/advance envelope assertions using both checklist sources | Stamp before input failure, skip status guard/sync, emit spawn fields on advance | Existing deterministic Git fixtures; no replacement lifecycle harness |
| AC-3 retained callers | Existing dispatch golden/host/hazard/state fixtures, ensigncycle cycle/feedback callers and skill integration fixtures move to explicit CLI/stdin/files while keeping meaningful output/state assertions | Change an artifact's host bootstrap, feedback body, model resolution, state commit, or worktree path | Existing suites; broad mechanical migration, no new runtime model |

Retire only tests whose claimed behavior is itself removed: JSON syntax/types/version requirements, JSON-vs-CLI host conflict precedence, print-schema and validate-only success. Replace their entry-point expectations with retirement diagnostics. Preserve non-JSON validation claims (missing entity/stage, bad host, feedback prerequisites, worktree errors) in the surviving CLI tests. Arrays used by old tests become newline checklist text under the newly declared semantics; multiline feedback/scope become files. Keep independent expected artifact bytes rather than deriving expectations with the production parser.

Current active caller inventory: 22 internal/dispatch test files reference mergeStdin/buildHostStdin; additional direct JSON cases live in build_errors_test.go. Real non-dispatch callers exist in internal/ensigncycle/cycle_test.go, feedback_test.go, pi_live_runner_test.go, and the generated operating prompt in haiku_loop_spike_live_test.go; skills/integration/dispatch_test.go directly calls RunWithLauncher with JSON. Migrate these call sites explicitly, including host/control fields becoming argv. Update native parity tests to provide CLI inputs; historical oracle inputs/recorded transcripts can remain JSON where they describe the old oracle/history, but must never be translated into a hidden native input compatibility path. Do not edit historical evidence merely to erase a search hit. The shared reviewer-reuse table's intentionally malformed --json examples remain negative evidence, labeled accordingly.

Necessity: the only new transport behavior is selecting stdin for checklist path '-', serving AC-1. Always writing the existing checklist file costs the very extra artifact the captain removed. JSON removal serves AC-3 and reduces mode/schema handling; a new --assignment-stdin selector or JSON-to-flags bridge would preserve that complexity. The existing line parser serves AC-2; a new parser or protocol is unnecessary. The FO default change serves AC-1 so this is exercised normal behavior, not an unused option. Existing skill integration and live journey own that proof; no second enforcement subsystem.

Implementation runs go test ./..., go test ./... -race, and gofmt -w ./cmd ./internal in its isolated worktree. Cycle-1 whole-repo attempts encountered unrelated pre-existing untracked tmp/dispatch-027/stash-recovery-copy/untracked-files compile errors; these are not a reason to delete coverage or to repair unrelated files in this task.

## Concrete documentation and FO wording

In internal/dispatch/dispatch.go, replace the stdin JSON/file-mode usages with the Selected contract's single grammar. Delete Input mode selection, required JSON fields, JSON examples, and successful --print-schema/--validate-only descriptions. Keep current file and advance examples; add this executable paired example (shell heredoc is checklist text, not JSON):

```sh
spacedock dispatch build --workflow-dir . --entity-path thing.md --stage validation --checklist-file - --advance <<'CHECKLIST'
DONE: verify
CHECKLIST
```

Replace --checklist-file's help with: “FILE contains one checklist item per non-empty line; use - to read the same format from stdin until EOF. Blank lines are ignored; CRLF and LF are accepted.” Add the retirement diagnostic and migration paragraph: “JSON requests, --print-schema and --validate-only are retired. Put entity, stage, host and controls in flags. Pipe checklist lines to --checklist-file -; use --scope-notes-file and --feedback-context-file for supporting prose. Stdout remains the dispatch JSON envelope.”

In skills/first-officer/references/fo-dispatch-core.md, replace:

> guard: write fragile inputs (checklist, scope notes, feedback context) to files first — one checklist item per non-empty line — so Markdown/backticks/shell-vars survive shell quoting.

with:

> guard: feed checklist text to stdin with --checklist-file - — one item per non-empty line. Use a quoted heredoc delimiter that does not occur as a complete input line, or pass literal stdin bytes through the tool. Keep scope notes and feedback context in files. Preserve Markdown, backticks and shell variables literally.

Change its normal effect example's --checklist-file {checklist_file} to --checklist-file - and show `<<'CHECKLIST'`, `{one_checklist_item_per_nonempty_line}`, and terminating `CHECKLIST` around the existing command. Keep all other options, lifecycle clauses, and verbatim helper-output forwarding unchanged. A checklist file remains supported for an existing supplied file; do not create one merely to dispatch. Update the one reuse-advance reference in skills/first-officer/references/claude-fo-dispatch.md to say “--checklist-file - with literal checklist lines on stdin; an existing checklist file is also accepted”; feedback context stays a file. This narrowly authorized skill change supersedes the cycle-1 no-skill-migration boundary.

In docs/site/reference/command-reference.md, append to the dispatch row: “Normal dispatch passes checklist lines on stdin with --checklist-file -; all other assignment identifiers and controls remain flags, with scope/feedback supplied by their existing file options. The checklist format is one item per non-empty line, identical for file and stdin. JSON request input, --print-schema and --validate-only are retired; stdout is still the same JSON envelope. Existing --checklist-file FILE calls remain supported.” Keep stamp wording unchanged. In docs/runtime-support.md, replace “accept host: <host>” for dispatch requests with “accept --host <host>”. No installer/host lifecycle changes.

## Expected surface and semantic boundaries

Estimate net LOC change: -500, across 40 files. Approximately 650 insertions and 1,150 deletions. Tolerance: net -900 to -150; 34–48 files. Revisit the design gate if deletion does not outweigh growth or the file tolerance is exceeded; do not retain dead JSON helpers just to fit the estimate. This wider file count reflects retiring actual callers and preserving their proof, not adding product mechanisms.

Expected product files: internal/dispatch/build.go and dispatch.go; the 22 dispatch tests using mergeStdin/buildHostStdin (build_statecommit, build_hazards, build_namecap, build_input_mode, build_codex_host, build_merged_mode, native_subcommands_routing, reconcile_namecap, build_advance, build_team_name_retired, build_advance_contract_parity, self_contained_assignment, build_stage_discipline_delivery, build_state_no_origin, standing_parity, build_parity, codex_bootstrap, cycle2_parity, build_advance_measurement, build_advisory_probe, build_json_ergonomics, build_pi_host; all suffixed _test.go); build_errors_test.go, build_stamp_test.go, parity_harness_test.go and help_test.go if its fixture requires migration; internal/ensigncycle/cycle_test.go, feedback_test.go, pi_live_runner_test.go, haiku_loop_spike_live_test.go and shared_promoted_live_test.go; skills/integration/dispatch_test.go; the two FO reference files; docs/site/reference/command-reference.md and docs/runtime-support.md. Up to two adjacent helper/golden files may adjust as callers are migrated. The state entity/report is one separate checkout file outside this product estimate.

Permitted semantic changes: CLI JSON request and schema/validation retirement; required-flag migration diagnostics; stdin checklist source and the FO default transport; no JSON host/flag precedence; checklist normalization follows the existing file format. Unchanged: output JSON schema, pointer/artifact assembly, stored workflow/entity formats, state authority, host detection, stamp/advance lifecycle and failure order after input loading, worker spawning, feedback content, and supported file callers. No new runtime lane, serializer, broad harness, compatibility bridge, model schema, or worker protocol.

## Stage Report: ideation

- DONE: Recommend the smallest CLI-control and stdin-assignment split with exact grammar, unambiguous mode selection, and legacy compatibility.
  Selected contract specifies --assignment-stdin, disjoint field ownership, exact mixed-input rejection, and unchanged legacy routing.
- DONE: Exercise the riskiest transport and mode-selection assumptions in a bounded fixture; tie AC-1, AC-2, and AC-3 to existing behavioral proof owners.
  Actual serialized stdin build preserved three hazard payload copies; existing mode/transport/help tests passed; proof table names the edit that would falsify each claim.
- DONE: Commit the ideation report with concrete contract and documentation wording, expected files and net LOC with tolerance, and semantic boundaries; stop before implementation.
  This path-scoped state commit contains the contract, proposed help/site wording, +240 net LOC across 7 files baseline, and explicit unchanged lifecycle boundaries.

### Summary

Recommended an explicit assignment-stdin opt-in so CLI controls can accompany structured assignment data without reinterpreting existing callers. Transport and legacy selection were exercised; proposed CLI integration remains for implementation after the ideation gate.

## Stage Report: ideation (cycle 2)

- DONE: One CLI flag/file path with --checklist-file - and explicit JSON retirement.
  Selected contract records the captain reset, exact grammar/EOF behavior and retirement diagnostics; AC-1 measures zero checklist scratch files against the existing one-file baseline.
- DONE: Migrated ACs and proof owners preserving literal-text/lifecycle coverage.
  AC-2's 12-call file/stdin-descriptor spike and four existing dispatch/integration baseline tests passed; AC-3 maps active JSON caller migration to existing host, artifact, feedback and stamp assertions rather than deleting them.
- DONE: Revised docs, FO default, estimate and committed cycle-2 report.
  Concrete help/site/FO wording prefers stdin; existing skill smoke and full-ensign-cycle journey own value proof; baseline is -500 net LOC across 40 files with declared tolerance.

### Summary

The revised design removes JSON request input and uses the existing checklist-file option's dash convention for stdin. It changes normal FO dispatch, preserves file-format and lifecycle behavior, and accounts for real caller/test migration; no product implementation was performed.

## Stage Report: implementation

- DONE: Implement the approved single flag/file interface with --checklist-file -; retire JSON request input and request-schema controls while preserving literal checklist, scope, feedback, envelope and lifecycle behavior.
  Code commit `ceaf05b6d` removes the JSON reader/selector/schema and JSON-host precedence; the existing normalized assembly map and schema-v2 output remain.
- DONE: Migrate real in-tree callers and existing tests without losing behavioral coverage; add focused failing tests before changes, update the approved FO/help/docs wording, and demonstrate AC-1, AC-2, and AC-3.
  Direct CLI/stdin/file calls replace active JSON callers; existing golden, host, state, bootstrap and feedback assertions remain. Initial stdin and retired-switch tests failed before the implementation, then passed.
- DONE: Commit the implementation and canonical report with actual checks, AC evidence, net LOC/files versus tolerance, and residual blockers; preserve unrelated work and stop for independent validation.
  `ceaf05b6d`: 39 product files, +1050/-1636, net -586 (approved 34–48 files and -900 to -150); code worktree clean; this report is committed path-scoped in the state checkout. No PR or merge performed.
- DONE: AC-1 — zero checklist scratch inputs in normal dispatch.
  `TestBuildChecklistStdin` compares zero new inputs with the one-file baseline; split-root/flat integration checks preserve directory entries and exact checklist bytes. Requiring a scratch file makes these fail.
- DONE: AC-1 — revised FO default exercised through the existing live journey.
  Candidate-bound `TestLiveCommonFullEnsignCycle` passed in 165.52s: durable done/report/path-scoped commit assertions; captured successful `--checklist-file - --stamp` quoted heredoc, with no separate checklist write. Evidence: `/tmp/dispatch-cli-stdin-candidate-live/codex-shared-scenarios/full-ensign-cycle/codex-exec.jsonl` and `codex-final-message.txt`.
- DONE: AC-2 — literal transport and whole-file supporting prose.
  `TestChecklistSourcesHaveIdenticalLiteralSections` compares independent expected sections for indentation, trailing spaces, CRLF/blank lines, EOF, one-terminal-CR removal, Unicode, JSON-looking text and 70 KB lines; `TestBuildQuotedHeredocTransport` runs a built CLI and checks the expansion sentinel stays absent. Trimming, Scanner limits or an unquoted heredoc falsify these checks.
- DONE: AC-3 — retirement, input refusal and literal dash-file behavior.
  Final input race checks passed (2.305s): retired no-flag JSON, empty input and reader/flag failures preserve artifact/entity/HEAD; `./-` never reads poison stdin, while `-` reads stdin even with a file named `-`. Restoring JSON fallback, reading the false file sentinel or mutating on refusal fails them.
- DONE: AC-3 — retained lifecycle and caller proof.
  Existing stamp success/mismatch/retry-sync/inline-before-worktree fixtures now run both sources; advance retains pointer-only envelopes. Full dispatch race passed (74.840s); existing cycle/feedback success and broken-output controls passed under race (4.140s). Skipping state guards/sync or changing artifact/report behavior falsifies these checks.
- DONE: Correct the owned Material proof-policy defect with FO authorization.
  The first full run caught the added instruction-text test via `TestNoInstructionReadsOutsideQuarantine`; contributors could not pass required checks (`contract[AGENTS.md#expected-commands]`). FO authorized replacing prohibited prose-grep with existing integration behavior; boundary/integration/input checks passed (1.265s/1.438s/2.246s). The live journey remains the FO-default proof owner.
- DONE: Run the required whole-repository tests and investigate failures.
  `go test ./...` failed on the above corrected lint issue and pre-existing `TestCodexResolveManifestAgainstInstalledHost`: `spacedock@spacedock not installed in codex, but resolver returned "/Users/clkao/.codex/plugins/cache/spacedock-local/spacedock/0.28.0-pre0/.codex-plugin/plugin.json"`. Log: `/tmp/dispatch-all-tests.log`; normal dispatch/integration/ensigncycle suites separately passed (41.337s/7.333s/247.462s).
- DONE: Establish baseline dependence and record the authorized disposition of the installed-manifest fixture mismatch.
  Identical failure reproduced from unchanged baseline `af70297ddae6ec64444849e8e3fcf57484bc16e1` (`/tmp/dispatch-baseline-resolver.log`): the test expects stable identity while production deliberately accepts the supported local install. Material verification issue, outside task ownership; FO authorized decline within this task, not a green-suite claim.
- DONE: Attempt the required whole-repository race check and record the authorized environmental limitation.
  `go test ./... -race` exhausted temporary disk space; `GOFLAGS=-p=1 go test ./... -race` repeated `no space left on device` in CLI and broad ensigncycle fixtures, then completed with exit 1. Remaining packages, including dispatch/contractlint/integration, passed; relevant cycle/feedback race tests passed separately. Logs: `/tmp/dispatch-all-race.log`, `/tmp/dispatch-all-race-serial.log`, `/tmp/dispatch-cycle-feedback-race.log`. FO authorized stopping after these failed attempts: no third broad retry or unrelated cleanup; the environmental limitation remains for independent validation.
- DONE: Required formatting and live setup accounting.
  Ran `gofmt -w ./cmd ./internal`; changed Go files and `git diff --check` are clean. Excluded an initial CI-auth-flag refusal and an installed-pre2 live run that correctly failed the new source assertion; the accepted live run explicitly used the built candidate with the existing local-auth path, without a new harness.

### Summary

Implemented the approved CLI/file interface, migrated callers, and proved the checklist value through exact command fixtures and the existing live Codex journey. The implementation is committed and ready for independent validation; full-suite green is explicitly limited by the reproduced baseline resolver test and repeated temporary-disk exhaustion, with no remaining task-owned failure identified.
