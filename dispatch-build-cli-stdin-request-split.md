---
title: Separate dispatch invocation controls from stdin assignment data
status: ideation
source: Captain discussion of checklist transport and dispatch assembly, 2026-09-10
started: 2026-09-10T20:26:39Z
completed:
verdict:
score: 0.5
worktree:
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
---

Make stdin assignment data usable alongside explicit CLI invocation controls without changing existing dispatch callers. This ideation records the implementation contract; product implementation has not started.

## Problem

FO dispatch currently writes checklist text to a temporary file and passes --checklist-file. The helper reads it only during assembly and embeds its contents in the generated assignment. No downstream reader needs the original checklist file.

Stdin JSON already exists, but request flags such as --stamp and --advance select flag/file mode. That prevents a straightforward combination of CLI invocation controls and structured assignment data. Avoid creating a second prompt assembler: the helper must continue producing the assignment and pointer envelope.

The original stdin interface came from the Python helper and was preserved by native commit 39454b644. Commit cfa3b671c added flag/file mode to avoid fragile shell quoting of prose. Any new default must preserve safe transport of arbitrary checklist, Markdown, scope, and feedback content. Removing a scratch file is a small simplification, not a claimed large latency improvement: writing the file and invoking build can already share one tool call.

## Selected contract

Add one boolean opt-in, `--assignment-stdin`. Keep the existing three entry points (legacy stdin, flag/file, schema/validation) intact. Do not infer input mode from whether stdin is a terminal, has bytes, or contains particular JSON keys.

Exact new grammar:

```text
spacedock dispatch build --assignment-stdin --workflow-dir DIR
  [--host claude|codex|pi] [--stamp | --advance]
  [--bare-mode] [--feedback-reflow]
```

Stdin is exactly one JSON object. Required fields are `schema_version` (2), `entity_path` (string), `stage` (string), and `checklist` (nonempty array of strings). Optional `scope_notes` and `feedback_context` are strings or null, retaining existing optional-text behavior. No prose is trimmed, split into lines, shell-evaluated, or reconstructed. Checklist elements retain their internal and trailing newlines; the existing assembler joins elements with one newline and adds its existing section delimiters.

`schema_version` stays in JSON as format metadata, not a new CLI switch. The helper inserts `workflow_dir`, `bare_mode`, `advance`, and `is_feedback_reflow` from CLI options into its existing internal request map, then uses the existing stamp/assembly path. Host remains the CLI flag or existing runtime detection; no new default or host resolution rule. Stage/model resolution remains workflow-owned; there is no new assignment model override.

In the new mode reject these JSON keys even when they equal CLI values or are null: `workflow_dir`, `host`, `bare_mode`, `advance`, `stamp`, `is_feedback_reflow`. Reject other unknown keys as well, so an attempted alternate control cannot be silently ignored. This is a new-mode allowlist, not a new generic JSON parser; ordinary JSON decoding semantics remain unchanged. No precedence merge is needed because accepted assignment keys and controls do not overlap.

Reject `--entity-path`, `--stage`, `--checklist-file`, `--scope-notes-file`, and `--feedback-context-file` whenever `--assignment-stdin` is present, including explicitly empty `--flag=` forms. Track their presence independently of the existing request-flag aggregate. Existing control flags may still set that aggregate; the opt-in takes priority over that legacy selection rule.

`--assignment-stdin` cannot combine with `--print-schema` or `--validate-only FILE`: reject with exit 2 before stdin read. Those two existing operations continue to describe/validate legacy schema-v2 requests; do not add another schema or validation subcommand in this task. Help explicitly states this limit. `--help` remains side-effect free. Grammar conflicts and forbidden assignment keys exit 2 with empty stdout; malformed/non-object JSON and existing field-validation failures retain existing exit behavior. New diagnostics should name the conflicting flag/key and the permitted location.

Mode selection table:

| Invocation | Source and behavior |
| --- | --- |
| `--assignment-stdin` plus allowed controls | Read stdin once as assignment data, inject controls, use existing helper pipeline |
| Opt-in plus assignment flags/files or schema/validation | Exit 2 before stdin read or mutation |
| No opt-in, any existing request flag | Existing flag/file mode; stdin still ignored; required trio unchanged |
| No opt-in, no request flags | Existing full schema-v2 stdin request, including legacy JSON control keys |
| Existing schema/validation operations without opt-in | Existing precedence and behavior unchanged |

Existing CLI incompatibility checks (`--stamp` with `--advance`, `--advance` with `--bare-mode`) remain before input loading. New input-shape checks happen during loading, before stamp. After normalization, do not move any existing downstream check across stamp or change its failure order. `--feedback-reflow` gets its required feedback text from `feedback_context` in the new mode. `--bare-mode` remains unsupported on Codex.

`--stamp` requires entity status already matching the assignment's `stage`; it stamps started/worktree, commits and synchronizes state, creates the stage worktree when needed, then emits the assembled pointer envelope. It neither advances status nor spawns a worker. The existing helper is still the only prompt assembler and emits exactly the same envelope/artifact shapes. No additional preflight or lifecycle abstraction.

## Acceptance criteria

**AC-1 — FO dispatch can supply assignment text without scratch input files.** A fixture invokes the real command path with `--assignment-stdin --workflow-dir DIR --stamp` and serialized assignment JSON; stdout supplies a readable helper-owned artifact containing checklist, scope, and feedback. Input scratch-file count is zero, compared with the current flag/file fixture's required checklist file (plus optional scope/feedback files). Changing the new path to require any of those files makes this test fail. No latency claim.

**AC-2 — Assignment content survives transport unchanged.** Multiline Markdown, quotes, backticks, dollar expressions, Unicode, angle brackets, and trailing newlines arrive in the artifact as independently specified byte sequences, with the existing section delimiters around them. A shell-command sentinel is not created. An implementation that trims, line-splits, HTML-escapes, or shell-expands values fails the test.

**AC-3 — Existing callers and dispatch effects remain compatible.** Legacy stdin-only and flag/file fixtures keep their current exits, diagnostics, artifacts, and host behavior. New opt-in cases reject duplicate control locations and file/assignment mixtures; identical duplicate values are still rejected. Existing stamp status-match refusal, stamp/advance exclusion, commit/sync/worktree ordering, retry behavior, and fresh-versus-advance envelopes hold when driven by the new input adapter. Build alone returns an envelope and does not spawn a worker.

## Proof and bounded spike

Observed baseline supplied by FO: installed 0.28.0-pre2, `dispatch build --workflow-dir docs/dev --stamp` with assignment JSON on stdin exits 2 with `flag/file input requires --entity-path, --stage, and --checklist-file`. This is evidence of the current mode boundary, not evidence the new switch exists.

A one-off Python standard-library fixture created a temporary Git workflow with a non-worktree backlog stage and one entity, serialized the following value with `json.dumps(..., ensure_ascii=False)`, passed bytes through `subprocess.run(..., input=...)` to installed 0.28.0-pre2 in legacy stdin mode, decoded its output envelope, and read the generated artifact:

```python
hazard = 'Keep "quotes", `backticks`, $(touch SHOULD_NOT_EXIST), $HOME, 雪 & <tag>\n```md\n# literal\n```\n\n'
```

The value occupied one checklist element, scope_notes, and feedback_context. Result: exit 0; artifact contained all three exact copies including final newlines; sentinel absent; workflow input files were exactly README.md and thing.md, with no checklist/scope/feedback scratch files. The temporary fixture and generated artifact were removed. This proves serializer-to-assembler transport on the retained mechanism, not proposed CLI integration. Reproduce in Go using the existing fixture helpers and `encoding/json.Marshal`; keep the expected payload as an independent literal, not as output from the production normalizer.

Bounded existing tests ran successfully: `go test ./internal/dispatch -run 'TestDispatchBuildAdvanceInputMode|TestBuildFlagFileInputModePreservesLiteralChecklist|TestBuildNoHTMLEscape|TestDispatchBuildHelpExamplesParse' -count=1`. They exercise legacy mode selection, literal file input, output escaping, and advertised example parsing. The opt-in integration is intentionally unimplemented and must begin with failing tests.

## Test plan and mechanism necessity

| Value | Existing primary proof owner and focused extension | Falsifying change | Cost/type |
| --- | --- | --- | --- |
| AC-1 | `build_input_mode_test.go` drives opt-in stamped build with `stampFixture`; assert artifact/envelope and zero scratch inputs | Restoring required checklist-file in new mode | Small deterministic Go fixture, seconds |
| AC-2 | `build_json_ergonomics_test.go` literal transport test, informed by `build_hazards_test.go` escaping coverage; add the above stdin payload and exact section bytes | Trim a trailing newline or re-escape payload before assembly | One deterministic fixture, seconds |
| AC-3 input routing | `build_input_mode_test.go` table for opt-in, legacy, forbidden flags including empty equals forms, forbidden JSON controls equal/different/null, unknown keys, schema/validation exclusion, malformed JSON | Select file mode merely because --stamp is set, or accept JSON host | Small table; reader that fails if touched proves grammar errors occur before read |
| AC-3 lifecycle | `build_stamp_test.go`: reuse status-mismatch, success, incompatibility, retry-sync and inline-before-worktree assertions for both input adapters; `build_advance_test.go` retains advance envelope authority | Assemble before stamp completes, skip status match, or emit spawn-only fields on advance | Existing deterministic Git fixtures; no new live runtime harness |
| Documentation accuracy | `TestDispatchBuildHelpExamplesParse` in `build_input_mode_test.go`: extend its existing example representation to pair serialized stdin with explicit argv | Publish a flag or JSON key combination the parser refuses | Small extension to existing owner |

The only new mechanism is the explicit adapter/selector, serving AC-1 and AC-3. Merely removing `--stamp`/`--advance` from `hasRequestFlags` would silently change legacy invocations and let JSON controls compete with flags, so it is insufficient. The new-mode key allowlist and incompatible-input checks serve AC-3; precedence merging is shorter locally but violates one authoritative source. Existing JSON serialization serves AC-2; hand-built heredoc JSON is insufficient because shell quoting does not escape JSON. No second assembler, generic transport layer, duplicate JSON-member parser, new broad harness, or host dispatch rewrite.

Implementation adds the focused failing tests before changing behavior, then runs `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`. This task adds CLI capability and documentation; it does not migrate FO skill commands, so no new live worker journey or skill prose smoke suite is owed. A later skill migration must bring its existing smoke and live-journey proof owners.

Required repository-wide checks were attempted during ideation: both `go test ./...` and `go test ./... -race` report undefined `decodeProviderResult` and `verifyProviderResolution` in the pre-existing untracked `tmp/dispatch-027/stash-recovery-copy/untracked-files/archived_annotation_decode_test.go`. No unrelated files were repaired. `gofmt -w ./cmd ./internal` ran; its unrelated two-line formatting change was restored to keep this stage state-only.

## Concrete documentation changes

In `internal/dispatch/dispatch.go`, add this usage line to both relevant help surfaces:

```text
spacedock dispatch build --assignment-stdin --workflow-dir DIR [--stamp | --advance] [--host claude|codex|pi]
```

Replace “The request comes from a JSON object on stdin OR from flags/files. The two are selected by the rule below and never merged.” with: “Use --assignment-stdin to keep invocation controls in CLI flags and assignment data in stdin JSON. Without this opt-in, the existing stdin and flag/file selection rules below are unchanged.” Prefix the existing request-flag mode paragraph with “Without --assignment-stdin,”. Replace the current absolute rejection wording for stdin plus --advance with: “Without --assignment-stdin, --advance selects flag/file mode and ignores stdin. With --assignment-stdin, --advance reads assignment data from stdin and emits the existing reuse pointer.”

Add flag text: “--assignment-stdin Read assignment JSON from stdin; controls stay in CLI flags. Requires --workflow-dir. Cannot combine with assignment/file flags, --print-schema, or --validate-only.” Add assignment fields and forbidden-key rules from Selected contract. Keep legacy schema fields under their current heading. Change --workflow-dir and --host labels from “both modes” to “all input modes”; describe --feedback-reflow's text source as --feedback-context-file in file mode or feedback_context in assignment mode. Replace stamp's “equals --stage” with “equals the requested stage”.

Add a paired help example with argv `spacedock dispatch build --assignment-stdin --workflow-dir . --advance` and stdin `{"schema_version":2,"entity_path":"thing.md","stage":"validation","checklist":["DONE: verify"]}`. Explain that the JSON object is stdin data, not another command. Keep existing examples intact.

In `docs/site/reference/command-reference.md`, append to the dispatch row: “Use `dispatch build --assignment-stdin --workflow-dir DIR --stamp` with a JSON object on stdin containing `schema_version: 2`, `entity_path`, `stage`, and a `checklist` array; optional `scope_notes` and `feedback_context` also stay in that object. Host, stamp/advance, bare mode, and feedback-reflow controls stay in CLI flags. This form needs no checklist, scope, or feedback input files. Serialize the object with a JSON encoder and pass its bytes to stdin; a quoted heredoc alone does not escape arbitrary JSON content. Legacy stdin and flag/file calls retain their existing behavior. `--print-schema` and `--validate-only` continue to describe legacy full requests and cannot combine with `--assignment-stdin`.” Replace the row's “equals `--stage`” with “equals the requested stage”.

## Expected surface and semantic boundaries

Estimate net LOC change: +240, across 7 files. Approximately 290 insertions and 50 deletions. Tolerance: net +160 to +330 and 6–8 files; exceed either bound only after revisiting the approved design.

Expected product files: `internal/dispatch/dispatch.go` (flag and help), `internal/dispatch/build.go` (small input adapter), `internal/dispatch/build_input_mode_test.go` (routing and help examples), `internal/dispatch/build_json_ergonomics_test.go` (byte transport), `internal/dispatch/build_stamp_test.go` (adapter reuse), `internal/dispatch/build_advance_test.go` (new entry route), and `docs/site/reference/command-reference.md`. The ideation body/report is one separate state-checkout file, excluded from product LOC estimates. No persistent fixture assets are needed.

Permitted semantics: one opt-in CLI input grammar and its errors; CLI help/docs explaining it. Stored formats, output envelope schema/version, generated prompt/assignment assembly, state authority, host rules, stamp lifecycle and failure ordering, worker spawning, legacy input selection, and legacy diagnostics remain unchanged. First officer owns state; ensigns own assigned deliverables/reports. No changes under agents/ or references/ and no FO skill migration.

## Stage Report: ideation

- DONE: Recommend the smallest CLI-control and stdin-assignment split with exact grammar, unambiguous mode selection, and legacy compatibility.
  Selected contract specifies --assignment-stdin, disjoint field ownership, exact mixed-input rejection, and unchanged legacy routing.
- DONE: Exercise the riskiest transport and mode-selection assumptions in a bounded fixture; tie AC-1, AC-2, and AC-3 to existing behavioral proof owners.
  Actual serialized stdin build preserved three hazard payload copies; existing mode/transport/help tests passed; proof table names the edit that would falsify each claim.
- DONE: Commit the ideation report with concrete contract and documentation wording, expected files and net LOC with tolerance, and semantic boundaries; stop before implementation.
  This path-scoped state commit contains the contract, proposed help/site wording, +240 net LOC across 7 files baseline, and explicit unchanged lifecycle boundaries.

### Summary

Recommended an explicit assignment-stdin opt-in so CLI controls can accompany structured assignment data without reinterpreting existing callers. Transport and legacy selection were exercised; proposed CLI integration remains for implementation after the ideation gate.
