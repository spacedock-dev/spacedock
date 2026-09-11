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
pr: "#784"
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
        - id: gate:09ptxgq8zma6qnx0h1w8a6wp:validation
          stage: validation
          attempts:
            - id: gate-attempt:09ptxgq8zma6qnx0h1w8a6wp-validation-1
              briefing:
                id: briefing:09ptxgq8zma6qnx0h1w8a6wp:validation:attempt-1:revision-1
                digest: sha256:546174f116f05e7785399a05614381173ada0877ee40733b8928a9db3121c5c9
                room-ref: '@review/validation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:09ptxgq8zma6qnx0h1w8a6wp:validation:1
                briefing: briefing:09ptxgq8zma6qnx0h1w8a6wp:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-09-11T00:30:13.874529Z"
                decision: revise
                reason: 'Captain expands this task: permit FO-authored scope instructions appended directly to the helper pointer sent to the ensign; unchanged prompt forwarding is not a requirement. Preserve the pointer optimization and checklist verification. Remove mandatory scope-notes scratch-file usage for ordinary dispatch, retain existing optional file compatibility and feedback transport, update contracts and exercise actual appended-instruction behavior. Re-review the combined implementation before delivery.'
            - id: gate-attempt:09ptxgq8zma6qnx0h1w8a6wp-validation-2
              briefing:
                id: briefing:09ptxgq8zma6qnx0h1w8a6wp:validation:attempt-2:revision-1
                digest: sha256:9a04c6a82dfa68d070a90229f8fbb32f282544503ab42d923d12f378e2364b00
                room-ref: '@review/validation/briefing-2'
              resolution:
                type: Resolution
                id: resolution:spacedock:09ptxgq8zma6qnx0h1w8a6wp:validation:2
                briefing: briefing:09ptxgq8zma6qnx0h1w8a6wp:validation:attempt-2:revision-1
                by: person:captain
                at: "2026-09-11T05:52:31.315756Z"
                decision: approve
                reason: 'Captain: approve PR and ci env. Approves validation attempt 2 for PR delivery and required CI environment approval; does not waive required checks.'
              application:
                target-stage: done
                state: pending
review-round:
    id: round:09ptxgq8zma6qnx0h1w8a6wp:validation:1
    stage: validation
    cycle: 1
    briefing:
        id: briefing:09ptxgq8zma6qnx0h1w8a6wp:validation:attempt-1:revision-1
        digest: sha256:546174f116f05e7785399a05614381173ada0877ee40733b8928a9db3121c5c9
        room-ref: '@review/validation/round-1'
mod-block: merge:pr-merge
---

Simplify normal FO dispatch to one CLI flag/file interface. The checklist can arrive on stdin through `--checklist-file -`; JSON requests and their schema/validation controls are retired. This design supersedes cycle 1 and includes the captain-directed appended-scope amendment; the combined implementation requires independent validation before delivery.

## Problem and captain decision

The helper only needs checklist text during assembly, yet the normal FO contract requires a separate checklist file. Keeping both JSON and flag/file requests creates mode selection, duplicate control locations, schema validation, and migration work for callers without improving this ordinary transaction.

Captain reset: “why do we want those argv in json as well? wouldn't it be simplify for checklist to be from stdin?”; “yeah and we can retire the json stdin code path”; “ok. send it back”. This explicitly replaces the prior JSON-compatibility requirement. Preserve supported flag/file callers and stamp/advance semantics, not the retired JSON API. Removing one scratch file is a small simplification, not a latency claim.

Captain-directed amendment after validation attempt 1 (rejected briefing digest `546174f116f05e7785399a05614381173ada0877ee40733b8928a9db3121c5c9`): preserve the emitted pointer but allow ordinary FO scope instructions appended directly to the worker message, without a scope-notes scratch file or helper roundtrip. The artifact is the standard assignment, while supplemental notes live in the conversation. Preserve transport identity/model and lifecycle authority; retain the dispatch checklist for report verification without rereading input files. This is an explicit scope change, not a discovered product defect, and retains the prior AC promises.

## Selected contract

The helper has exactly one input path. Required assignment identifiers and invocation controls are CLI flags. Checklist source is the existing file option with one conventional sentinel:

```text
spacedock dispatch build --workflow-dir DIR --entity-path FILE --stage STAGE
  --checklist-file FILE|- [--scope-notes-file FILE] [--feedback-context-file FILE]
  [--host claude|codex|pi] [--bare-mode] [--feedback-reflow] [--stamp | --advance]
```

`--checklist-file -` reads the command's supplied stdin reader until EOF; it does not open a file named `-` and does not read global os.Stdin behind the injectable reader. A literal file named `-` remains addressable as `./-`. Any other checklist path retains ordinary file behavior and does not consume stdin. Existing optional scope/feedback file options remain compatible; `-` has no new special meaning for them. Ordinary FO scope instructions may instead be appended directly after the intact emitted pointer in the worker message, with no helper roundtrip. Feedback context remains opaque file transport. There is no autodetection, terminal probe, JSON request, new mode selector, or second assembly pipeline.

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

**AC-1 — Normal FO dispatch needs zero checklist scratch files.** The shared FO command prefers --checklist-file - and feeds literal checklist text to stdin. A real-command fixture executes that documented shape and gets the helper-owned artifact/envelope with zero checklist input files, compared with the current file-mode fixture's one required checklist file. The existing live full-ensign-cycle journey reaches durable done state and commits using the revised skill; its captured dispatch evidence shows the stdin checklist source without a separately written checklist file. Optional scope files and opaque feedback files remain supported. Ordinary scope instructions may be appended directly after the intact helper pointer; no scope-notes scratch file or helper roundtrip is required. The existing live journey must show that such an appended literal constraint reaches the worker and affects its durable report, while retaining the prior stdin, done and commit proof. The FO retains its checklist to verify the report without rereading original input files. Requiring a checklist or ordinary scope scratch file again, or dropping the appended constraint, fails this value check.

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
| AC-1 live value | internal/ensigncycle/shared_live_runner_test.go TestLiveCommonFullEnsignCycle and shared_promoted_live_test.go runFullEnsignCycleJourney: revised FO core, existing transcript artifacts plus durable done/report/path-scoped commit proof; check stdin source and absence of a separate checklist write, plus actual appended constraint receipt and literal durable worker output in existing captured evidence | FO writes a checklist/scope scratch file, omits the appended instruction, or the worker loses the literal constraint or completion | One existing Codex lane journey, roughly 5–15 minutes; no new lane/driver |
| AC-2 | build_json_ergonomics_test.go's literal-checklist test (retain filename or rename only if useful), build_hazards_test.go: file/stdin table seeded by the spike, exact normalized section bytes; scope/feedback file byte assertions; execute a quoted-heredoc hazard once through existing command fixture | Strip indentation/trailing space; Scanner rejects 70 KB; unquote heredoc delimiter | Small deterministic tables and one shell transport exercise |
| AC-3 CLI | build_input_mode_test.go and build_errors_test.go: retired switches including equals, no-trio JSON caller, missing args, empty input, reader failure, regular path with a poison stdin reader; explicit no envelope/no mutation on input errors | Keep legacy JSON fallback, read stdin in file mode, accept retired switch | Existing parser/command owner, seconds |
| AC-3 lifecycle | build_stamp_test.go and build_advance_test.go: drive existing success/mismatch, incompatibility, retry-sync, inline-before-worktree and fresh/advance envelope assertions using both checklist sources | Stamp before input failure, skip status guard/sync, emit spawn fields on advance | Existing deterministic Git fixtures; no replacement lifecycle harness |
| AC-3 retained callers | Existing dispatch golden/host/hazard/state fixtures, ensigncycle cycle/feedback callers and skill integration fixtures move to explicit CLI/stdin/files while keeping meaningful output/state assertions | Change an artifact's host bootstrap, feedback body, model resolution, state commit, or worktree path | Existing suites; broad mechanical migration, no new runtime model |

Retire only tests whose claimed behavior is itself removed: JSON syntax/types/version requirements, JSON-vs-CLI host conflict precedence, print-schema and validate-only success. Replace their entry-point expectations with retirement diagnostics. Preserve non-JSON validation claims (missing entity/stage, bad host, feedback prerequisites, worktree errors) in the surviving CLI tests. Arrays used by old tests become newline checklist text under the newly declared semantics; multiline feedback/scope become files. Keep independent expected artifact bytes rather than deriving expectations with the production parser.

Current active caller inventory: 22 internal/dispatch test files reference mergeStdin/buildHostStdin; additional direct JSON cases live in build_errors_test.go. Real non-dispatch callers exist in internal/ensigncycle/cycle_test.go, feedback_test.go, pi_live_runner_test.go, and the generated operating prompt in haiku_loop_spike_live_test.go; skills/integration/dispatch_test.go directly calls RunWithLauncher with JSON. Migrate these call sites explicitly, including host/control fields becoming argv. Update native parity tests to provide CLI inputs; historical oracle inputs/recorded transcripts can remain JSON where they describe the old oracle/history, but must never be translated into a hidden native input compatibility path. Do not edit historical evidence merely to erase a search hit. The shared reviewer-reuse table's intentionally malformed --json examples remain negative evidence, labeled accordingly.

Necessity: the helper's only new input behavior is selecting stdin for checklist path '-', serving AC-1. Always writing the existing checklist file costs the very extra artifact the captain removed. JSON removal serves AC-3 and reduces mode/schema handling; a new --assignment-stdin selector or JSON-to-flags bridge would preserve that complexity. The existing line parser serves AC-2; a new parser or protocol is unnecessary. The FO default and direct appended scope instructions serve AC-1: the existing worker message already carries supplemental instructions, so no file, receipt or metadata mechanism is necessary. Existing skill integration and live journey own that proof; no second enforcement subsystem.

Implementation runs go test ./..., go test ./... -race, and gofmt -w ./cmd ./internal in its isolated worktree. Cycle-1 whole-repo attempts encountered unrelated pre-existing untracked tmp/dispatch-027/stash-recovery-copy/untracked-files compile errors; these are not a reason to delete coverage or to repair unrelated files in this task.

## Concrete documentation and FO wording

In internal/dispatch/dispatch.go, replace the stdin JSON/file-mode usages with the Selected contract's single grammar. Delete Input mode selection, required JSON fields, JSON examples, and successful --print-schema/--validate-only descriptions. Keep current file and advance examples; add this executable paired example (shell heredoc is checklist text, not JSON):

```sh
spacedock dispatch build --workflow-dir . --entity-path thing.md --stage validation --checklist-file - --advance <<'CHECKLIST'
DONE: verify
CHECKLIST
```

Replace --checklist-file's help with: “FILE contains one checklist item per non-empty line; use - to read the same format from stdin until EOF. Blank lines are ignored; CRLF and LF are accepted.” Add the retirement diagnostic and migration paragraph: “JSON requests, --print-schema and --validate-only are retired. Put entity, stage, host and controls in flags. Pipe checklist lines to --checklist-file -; append ordinary scope instructions after the intact helper pointer sent to the worker. --scope-notes-file remains optional; use --feedback-context-file for opaque feedback. Stdout remains the dispatch JSON envelope.”

In skills/first-officer/references/fo-dispatch-core.md, replace:

> guard: write fragile inputs (checklist, scope notes, feedback context) to files first — one checklist item per non-empty line — so Markdown/backticks/shell-vars survive shell quoting.

with:

> guard: feed checklist text to stdin with --checklist-file - — one item per non-empty line. Use a quoted heredoc delimiter that does not occur as a complete input line, or pass literal stdin bytes through the tool. Append ordinary scope instructions directly after the intact helper pointer; --scope-notes-file remains optional, and feedback stays in --feedback-context-file. Preserve Markdown, backticks and shell variables literally.

Change its normal effect example's --checklist-file {checklist_file} to --checklist-file - and show `<<'CHECKLIST'`, `{one_checklist_item_per_nonempty_line}`, and terminating `CHECKLIST` around the existing command. Keep CLI options and lifecycle clauses; preserve the emitted pointer and transport, identity, description and model fields while allowing ordinary FO scope instructions appended directly. The durable artifact is the standard assignment, not a claim to contain all conversation context. A checklist file remains supported for an existing supplied file; do not create one merely to dispatch. Update the one reuse-advance reference in skills/first-officer/references/claude-fo-dispatch.md to say “--checklist-file - with literal checklist lines on stdin; an existing checklist file is also accepted”; feedback context stays a file. Update shared FO/ensign and Claude, Codex and Pi runtime forwarding/bootstrap clauses consistently for fresh and advance messages. This captain amendment supersedes the former unchanged-forwarding and mandatory scope-file wording; no new envelope or metadata channel is introduced.

In docs/site/reference/command-reference.md, append to the dispatch row: “Normal dispatch passes checklist lines on stdin with --checklist-file -; all other assignment identifiers and controls remain flags, with ordinary scope instructions optionally appended directly after the intact helper pointer. Scope files remain optional, and feedback keeps its opaque file option. The checklist format is one item per non-empty line, identical for file and stdin. JSON request input, --print-schema and --validate-only are retired; stdout is still the same JSON envelope. Existing --checklist-file FILE calls remain supported.” Keep stamp wording unchanged. In docs/runtime-support.md, replace “accept host: <host>” for dispatch requests with “accept --host <host>”. No installer/host lifecycle changes.

## Expected surface and semantic boundaries

Estimate net LOC change: -500, across 40 files. Approximately 650 insertions and 1,150 deletions. Tolerance: net -900 to -150; 34–48 files. Revisit the design gate if deletion does not outweigh growth or the file tolerance is exceeded; do not retain dead JSON helpers just to fit the estimate. This wider file count reflects retiring actual callers and preserving their proof, not adding product mechanisms.

Expected product files: internal/dispatch/build.go and dispatch.go; the 22 dispatch tests using mergeStdin/buildHostStdin (build_statecommit, build_hazards, build_namecap, build_input_mode, build_codex_host, build_merged_mode, native_subcommands_routing, reconcile_namecap, build_advance, build_team_name_retired, build_advance_contract_parity, self_contained_assignment, build_stage_discipline_delivery, build_state_no_origin, standing_parity, build_parity, codex_bootstrap, cycle2_parity, build_advance_measurement, build_advisory_probe, build_json_ergonomics, build_pi_host; all suffixed _test.go); build_errors_test.go, build_stamp_test.go, parity_harness_test.go and help_test.go if its fixture requires migration; internal/ensigncycle/cycle_test.go, feedback_test.go, pi_live_runner_test.go, haiku_loop_spike_live_test.go and shared_promoted_live_test.go; skills/integration/dispatch_test.go; the shared FO/ensign references and their Claude, Codex and Pi runtime forwarding/bootstrap clauses; docs/site/reference/command-reference.md and docs/runtime-support.md. Up to two adjacent helper/golden files may adjust as callers are migrated. The state entity/report is one separate checkout file outside this product estimate.

Permitted semantic changes: CLI JSON request and schema/validation retirement; required-flag migration diagnostics; stdin checklist source and the FO default transport; ordinary scope instructions appended directly after the intact helper pointer with no required scratch file or helper roundtrip; no JSON host/flag precedence; checklist normalization follows the existing file format. Unchanged: output JSON schema, pointer/artifact assembly, stored workflow/entity formats, state authority, host detection, stamp/advance lifecycle and failure order after input loading, worker transport identity/model and spawning authority, feedback content, and supported file callers. No new runtime lane, serializer, broad harness, compatibility bridge, model schema, or worker protocol.

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


## Stage Report: validation

- DONE: Independently assess AC-1, AC-2, and AC-3 against the approved revision and implementation; verify CLI retirement, stdin/file semantics, preserved callers and stamp/advance authority using the report and appropriate existing proof.
  Validated candidate `ceaf05b6d6b48a1cf10d712cd65e3ffc3e31d3fb` against approved ideation attempt 2; AC citations are this entity's lines 106, 108 and 110. Cycle-1 JSON compatibility is superseded.
- DONE: AC-1 — Normal FO dispatch needs zero checklist scratch files.
  Independent `TestBuildChecklistStdin`, `TestBuildQuotedHeredocTransport`, `TestSplitRootFolderWorktreeDispatch` and `TestFlatEntitySlugUnchanged` passed: helper artifact/envelope, literal stdin and zero extra inputs versus file baseline. Opening dash as a file or requiring scratch input fails these checks.
- DONE: AC-1 — Existing live journey reaches durable done with the revised FO default.
  Reused and inspected `/tmp/dispatch-cli-stdin-candidate-live.log` (165.52s PASS) and its full-ensign-cycle `codex-exec.jsonl`: successful quoted-heredoc `--checklist-file - --stamp`, candidate-path fetch command, schema-v2 pointer, archived done/PASSED and clean Git log containing completion commit `0664c38`; no checklist write in captured commands. Reverting the default or losing durable completion invalidates this proof.
- DONE: AC-2 — Checklist and supporting prose retain the surviving format's semantics.
  Independent `TestChecklistSourcesHaveIdenticalLiteralSections` passed both sources across CRLF/blanks, EOF/LF, Unicode/shell hazards, JSON-looking text, 70 KB and doubled-CR boundaries; independent expected sections and whole-artifact equality preserve opaque scope/feedback bytes. Heredoc test checks the shell sentinel stays absent.
- DONE: AC-3 — The single supported input interface replaces JSON without losing dispatch behavior.
  Independent retirement/no-read/no-mutation, literal `./-` versus stdin `-`, advance-envelope, status-mismatch, retry-sync and inline-before-worktree checks passed. JSON reader/schema/selector/version helpers are deleted; migrated callers pass direct flags/bytes, with historical oracle JSON retained as history. Restoring JSON fallback, consuming poison stdin or bypassing state/sync guards falsifies these checks.
- DONE: Reproduce the bounded command, input and lifecycle evidence without repeating broad green suites.
  In detached candidate checkout: `go test ./internal/dispatch ./skills/integration -run 'TestBuildChecklistStdin|TestBuildQuotedHeredocTransport|TestChecklistSourcesHaveIdenticalLiteralSections|TestBuildRequestControlsRetired|TestBuildInputFailuresDoNotReadOrStamp|TestChecklistFileDoesNotReadStdin|TestDispatchBuildAdvanceInputMode|TestStampRefusesStatusStageMismatchWithoutMutation|TestStampRetriesSyncOnRetryEvenWhenAlreadyStamped|TestStampCommitsInlineBeforeWorktreeCreation|TestSplitRootFolderWorktreeDispatch|TestFlatEntitySlugUnchanged' -count=1`; PASS 10.183s/0.919s, `/tmp/dispatch-validation-focused.log`.
- DONE: Perform the required detached adversarial audit on a throwaway checkout for changed CLI/contract behavior; attack a distinct claim rather than rerunning already-green suites, and route material findings before any candidate edits.
  Created detached `.validation-audit` beneath the assigned worktree at `ceaf05b6d`; changed only `checklistLines` from `strings.TrimSuffix(line, "\r")` to `strings.TrimRight(line, "\r")`, attacking exact one-CR removal while leaving common LF/CRLF results intact. Existing test caught the wrong bytes; no new harness or candidate changes.
- DONE: Demonstrate the distinct falsifying edit and restoration.
  `go test ./internal/dispatch -run '^TestChecklistSourcesHaveIdenticalLiteralSections$/one_terminal_CR$' -count=1` exited 1 at `build_json_ergonomics_test.go:296` (`checklist section bytes changed`); restored code exited 0 (0.362s). Logs: `/tmp/dispatch-validation-mutant.log`, `/tmp/dispatch-validation-restored.log`; detached checkout removed after cleanup.
- DONE: Audit necessity, soundness, representations and scaling.
  Dash source selection is the smallest mechanism delivering AC-1; deleting JSON serves AC-3 without an adapter. Both sources use one linear byte-to-lines parser and existing normalized fields/assembly; no Scanner cap or new lifecycle controller. Read-to-EOF allocation remains proportional to input, matching the supported file behavior; 70 KB coverage exercises the relevant limit. Source-read and CLI refusal precede stamp, while downstream validation/stamp order remains as approved.
- DONE: Compare approved surface and preserve candidate ownership.
  `af70297dd..ceaf05b6d`: 39 files, +1050/-1636, net -586, within 34–48 files and -900 to -150. Candidate HEAD and bytes unchanged; `git diff --check` clean. `gofmt -w ./cmd ./internal` in audit exposed only pre-existing alignment in unchanged `internal/release/runtime_live_evidence_workflow_test.go`; discarded that audit-only formatting.
- DONE: Verify recorded required-suite limitations and authorized dispositions.
  Inspected normal/race logs: unchanged installed-manifest failure matches `/tmp/dispatch-baseline-resolver.log`; normal run also hit disk exhaustion, as did both broad race attempts. The prohibited prose-read test was removed; current behavioral proof replaces it. Reused final input race (2.305s), cycle/feedback race (4.140s) and reported full dispatch race; these do not make whole-repository suites green.
- SKIPPED: Repeat whole-repository normal/race suites or launch provider operations.
  FO authorized no third broad retry or unrelated cache cleanup; only 1.2 GiB was available. Existing `/tmp/dispatch-all-tests.log`, `/tmp/dispatch-all-race.log` and `/tmp/dispatch-all-race-serial.log` retain failures. Provider operations are outside this assignment.
- DONE: Commit a canonical validation report with exact AC citations, approved-surface comparison, honest baseline/resource limitations, and PASSED or REJECTED recommendation; do not approve, merge, publish a PR, or modify implementation.
  Recommendation: PASSED for local AC validation; detached audit refuted nothing material. No new material, deferred-risk or polish findings. This recommendation is not merge authorization: all host live CI lanes required by `docs/dev/README.md#proof-policy` remain pending before merge, acknowledged by FO; none is waived or claimed green.

### Summary

AC-1, AC-2 and AC-3 have independent command/state evidence for the approved interface, and the detached mutation demonstrated a real byte-preservation failure boundary. Recommend PASSED for validation, with whole-suite baseline/resource failures and required host live CI before merge explicitly outstanding; candidate code was not modified.


## Stage Report: implementation (cycle 2)

- DONE: Implement the captain-approved amendment allowing direct appended scope instructions with the helper pointer; update the task and affected contracts consistently while preserving pointer bootstrap and checklist verification.
  `4db65d817` updates shared FO/ensign contracts and Claude, Codex and Pi forwarding/bootstrap clauses: intact pointer plus ordinary conversation instructions, preserved transport identity/model/lifecycle authority, retained checklist verification, optional scope files and opaque feedback files. Current design/AC wording now reflects the captain amendment without narrowing prior promises.
- DONE: Exercise actual worker receipt and use of appended instructions with existing behavioral proof; preserve existing dispatch lifecycle and optional file compatibility without a new harness.
  Before skill edits, extended the existing full-cycle journey and passed its missing/expanded-note negative controls (0.295s, `/tmp/dispatch-amendment-proof.log`); no instruction-file assertions or new harness. The live scenario alone excludes scope-file routing because it explicitly requests appended notes; optional file compatibility remains exercised separately.
- DONE: AC-1 — appended literal scope affects durable worker output without helper input or another scratch file.
  Candidate-bound `TestLiveCommonFullEnsignCycle` passed in 124.26s (`/tmp/dispatch-amendment-live.log`), retaining stdin checklist, durable done/report and path-scoped commit assertions. The report records exactly “The worker preserved $HOME and `uname` as literal scope text.” Missing or expanded text makes the assertion fail; routing it through the helper fails this scenario's command check.
- DONE: Preserve the live separation evidence and state its visibility limit.
  `/tmp/dispatch-amendment-live/codex-shared-scenarios/full-ensign-cycle/` holds `codex-exec.jsonl`, `codex-final-message.txt`, `codex-process-result.txt` and `dispatch-assignment.md`. The preserved artifact matches the accepted run's fetch command and excludes the literal note; captured helper input also excludes it, while the worker report includes it. The FO reads the report and counts 1 done/0 skipped/0 failed. The public stream omits the spawn payload, so worker receipt is established by that separation and durable effect, not a captured spawn-argument assertion.
- DONE: AC-2/AC-3 — retain literal file compatibility, refusal safety and integration proof.
  Focused dispatch/integration/contractlint checks passed normally (0.804s/0.619s/0.794s) and under race (1.886s/1.659s/2.277s): `/tmp/dispatch-amendment-focused.log`, `/tmp/dispatch-amendment-focused-race.log`. They exercise independent expected literal sections including scope/feedback files, false dash filename and no-mutation refusal behavior, split-root/flat dispatch, help routing and instruction-read boundaries; breaking those behaviors fails the checks.
- DONE: Commit the combined change and an honest stage report with AC evidence, surface comparison, and validation limitations; preserve unrelated edits and do not merge or publish.
  Combined `af70297dd..4db65d817`: 45 files, +1105/-1657, net -552 (approved 34–48 files/net -900 to -150). Amendment alone: 11 existing files, +58/-24; no new file or runtime mechanism. Code worktree clean; required gofmt ran, unrelated baseline formatting restored, diff check clean. State changes are body-only and committed path-scoped; no PR, merge or provider publication.
- DONE: Preserve prior failed-suite observations and authorized dispositions without claiming whole-suite green.
  The baseline installed-manifest mismatch remains declined outside this task (`/tmp/dispatch-all-tests.log`, `/tmp/dispatch-baseline-resolver.log`); both failed broad race attempts retain disk-exhaustion evidence (`/tmp/dispatch-all-race.log`, `/tmp/dispatch-all-race-serial.log`). No broad suites were repeated for this amendment. An extra live capture run had just started when FO declined an evidence-only repeat; canceled at 17.430s with `signal: interrupt`, exit 1 (`/tmp/dispatch-amendment-receipt.log`), and excluded from acceptance evidence.

### Summary

Implemented the captain-directed scope change from rejected validation attempt 1, briefing digest `546174f116f05e7785399a05614381173ada0877ee40733b8928a9db3121c5c9`; it is not classified as a discovered product defect. The combined candidate preserves the original stdin/JSON-retirement ACs and now exercises direct appended instructions through durable worker output. Ready for independent validation with the prior full-suite limitations and public spawn-payload visibility limit explicitly retained.

### Feedback Cycles

- Cycle 1: REJECTED — captain-directed scope amendment / validation re-review; surface 45 files/net -552 LOC vs estimate 40 files/net -500 LOC (12.5% files, 10.4% LOC; within declared tolerance); AC expanded: permit direct appended scope instructions with the intact helper pointer, retain checklist verification and optional file compatibility. Captain authorized folding this into the same task; implementation completed at 4db65d817, pending independent re-review.


## Stage Report: validation (cycle 2)

- DONE: Re-review the combined candidate and amended ACs: direct appended scope instructions must affect durable worker output without a scope scratch file, preserving pointer bootstrap, checklist verification, file compatibility and lifecycle authority.
  Reviewed `4db65d817` against captain-expanded AC-1 and retained AC-2/AC-3 (this entity, Acceptance criteria, lines 132/134/136). This supersedes the prior validation recommendation following the captain-directed scope amendment; it is not correction of a discovered implementation defect.
- DONE: AC-1 — Normal FO dispatch needs zero checklist scratch files, and appended scope affects durable worker output.
  Reused `/tmp/dispatch-amendment-live.log`: candidate-bound `TestLiveCommonFullEnsignCycle` PASS in 124.26s. Inspected the successful stdin-heredoc dispatch command, schema-v2 pointer, preserved `dispatch-assignment.md`, worker report read, archived done/PASSED and clean working tree in `/tmp/dispatch-amendment-live/codex-shared-scenarios/full-ensign-cycle/codex-exec.jsonl`; the existing live grader also requires a path-scoped completion commit.
- DONE: AC-1 — Distinguish supplemental conversation scope from helper/artifact contents and verify the checklist.
  The preserved artifact's fixture path/fetch command matches the accepted run; neither it nor the helper command contains “The worker preserved $HOME and `uname` as literal scope text.” The durable implementation report contains that exact note; FO reads the report and states 1 done/0 skipped/0 failed without rereading checklist/scope inputs. Captured commands show one dispatch and no scope/checklist scratch-file write. Dropping or shell-expanding the note invalidates the observed result.
- DONE: State the worker-receipt observation boundary precisely.
  The public live stream omits spawn arguments. Report/artifact/helper-command separation and the durable effect support conversation delivery; this is not an assertion that exact spawn arguments were captured. The note grader alone checks entity text, so report placement and absence from helper input are established by the inspected accepted-run evidence, not that predicate alone.
- DONE: AC-2 — Checklist and supporting prose retain the surviving format's semantics.
  Reused prior independent literal-source matrix and cycle-2 focused normal/race results (`/tmp/dispatch-amendment-focused.log`, `/tmp/dispatch-amendment-focused-race.log`): ordinary optional scope/feedback files still preserve bytes; file/stdin/false-dash behavior remains covered. Trimming, shell expansion, a Scanner cap or special scope-file dash routing would fail the existing checks; the amendment changes no parser code.
- DONE: AC-3 — The single supported input interface replaces JSON without losing dispatch behavior.
  Reused prior retirement/no-read/no-mutation, host, split-root, stamp/retry-sync and advance-envelope evidence. The amendment changes instruction forwarding and the existing live grader only; it preserves emitted pointer, transport identity/model and FO lifecycle authority. Shared and Claude/Codex/Pi instructions consistently permit appended notes for their supported dispatch paths while retaining optional scope files and opaque feedback files.
- DONE: Evaluate amendment evidence and an appropriately bounded detached adversarial audit, reusing prior green proof and recording precise limitations without redundant live or broad-suite reruns.
  Created detached `.validation-audit-scope` beneath the assigned worktree at `4db65d817`. The existing note negative-control test passed pristine; changed its grader to accept the prefix “The worker preserved” instead of the exact literal note. This would accept shell-expanded text despite successful lifecycle; the existing test rejected the mutation. No new harness or provider call.
- DONE: Record a distinct falsification and restored result.
  `go test -tags live ./internal/ensigncycle -run '^TestFullCycleSupplementalNoteRejectsMissingWorkerOutput$' -count=1`: pristine exit 0 (0.290s), mutant exit 1 at `shared_promoted_live_test.go:29` (“missing or shell-expanded supplemental note must fail despite completed lifecycle”), restored exit 0. Logs: `/tmp/dispatch-validation-scope-pristine.log`, `/tmp/dispatch-validation-scope-mutant.log`, `/tmp/dispatch-validation-scope-restored.log`.
- DONE: Challenge necessity and soundness across the amended representations.
  Appending conversation text delivers the requested value through the existing worker message; no serializer, scratch-file lifecycle or helper roundtrip is needed. Traced intact pointer → standard artifact plus supplemental conversation → exact durable report → retained-checklist count → done state. Missing/expanded/exact note controls and existing file/stdin/advance cases cover adjacent variants; no new allocation, parser, blocking I/O or multiplicative hot path was introduced. Detached audit refuted nothing material.
- DONE: Compare scope and preserve ownership.
  Actual `af70297dd..4db65d817`: 45 files, +1105/-1657, net -552, inside 34–48 files/net -900 to -150; amendment is 11 existing files, +58/-24. Candidate HEAD/bytes unchanged, `git diff --check` clean. Restored and removed the owned detached checkout; no candidate repair or unrelated change.
- SKIPPED: Repeat accepted live evidence, whole-repository suites or provider operations.
  Prior normal-suite installed-manifest baseline and normal/race disk failures remain recorded and unwaived. FO declined redundant live capture; `/tmp/dispatch-amendment-receipt.log` ends `signal: interrupt` at 17.430s and is excluded. Disk availability improved to 7.9 GiB, but no changed code or new failure justified another broad run. Required all-host live CI remains pending before merge.
- DONE: Commit the canonical validation report with AC evidence, scope comparison, findings and PASSED or REJECTED recommendation; do not modify candidate code, approve, merge or publish.
  Recommendation: PASSED for amended local AC validation. No new material, deferred-risk or polish findings; preserve the explicit public-stream visibility limit and outstanding baseline/resource/CI obligations. This report provides no merge authorization and does not mark pending host lanes green.

### Summary

The combined candidate satisfies the expanded scope through literal supplemental instructions that affect the durable worker report, alongside the original stdin and JSON-retirement evidence. Recommend PASSED for validation, with exact spawn arguments unobserved and all required host live CI still pending before merge; no candidate code was changed.
