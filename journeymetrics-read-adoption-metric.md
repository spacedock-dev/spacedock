---
title: Measure actual `status --read` adoption in FO/ensign journeys via journeymetrics
status: ideation
source: "captain (2026-06-16, 0204 sprint) — e6a's status --read adoption was proven WORKING once (AC6 single live drive). No ongoing verification that real FO/ensign sessions call --read; the six contract sites are wording, and wording-present is not behavior. journeymetrics already parses every journey transcript but records ToolCallsByName by tool NAME only — `status --read` is a Bash subcommand, invisible today as a generic Bash call."
score:
sprint: 0204-structured-reads
sprint-readiness: ready
issue:
id: hf4jmbksapyg2d9s0zj85wca
started: 2026-06-16T20:02:50Z
---

## Problem
e6a's `status --read` adoption is proven to work exactly once (AC6's single live FO drive during e6a validation). There is no standing, behavioral check that real FO/ensign sessions actually invoke `--read`. The instruction lives in six contract sites; that is wording, not observed behavior.

`internal/journeymetrics` parses every journey transcript and records `ToolCallsByName` — but keyed by tool NAME only (`toolCallsByName[block.Name]++` in `claude.go:75`; the codex side does the same in `codex.go:52`). A `status --read` invocation is a Bash subcommand, so it is counted as a generic `Bash` call and the `--read` flag is invisible. A scoped `Read` (one carrying `offset`/`lines`) is counted identically to a full-file `Read`.

## Spike result
**No spike needed: the Bash `tool_use` command STRING is already on disk in the transcripts journeymetrics parses; the parser merely discards it.**

The riskiest assumption was whether the Bash command args (not just the tool name) are reachable in the stream-json transcript. They are, proven against committed real captures:

- Real ensigncycle capture `internal/ensigncycle/testdata/shallow-boot-greet.stream.jsonl` carries `{"type":"tool_use",...,"name":"Bash","input":{"command":"spacedock status --boot --json"}}`. Other real captures carry `spacedock status --discover`, `spacedock status --workflow-dir … --boot`, etc. The `input.command` string is present and complete.
- Real `Read` tool_use blocks carry `input.file_path` (e.g. `internal/ensigncycle/testdata/sonnet_teamdelete_marker_continues.stream.jsonl`). A scoped `Read` carries `offset`/`lines` alongside it in the same `input` object — the same object the parser drops.
- The committed journeymetrics fixtures (`testdata/claude_terminal_split.stream.jsonl`, `testdata/claude_no_terminal.stream.jsonl`) already carry `"input":{"command":"…"}` on their Bash blocks, confirming the format the parser consumes.

The reason `--read` is invisible today is purely that `claudeContent` (`claude.go:214-218`) unmarshals only `Type`/`ID`/`Name` and never reads `input`. The fix is to unmarshal `input` and inspect it — no new data source, no runtime handoff, no on-disk format change. Codex is symmetric: `tool_call.started` carries `arguments.cmd` (e.g. `{"cmd":"go test ./..."}` in `testdata/codex-exec.jsonl`), currently discarded.

## Proposed approach
Extend the claude transcript parser (and mirror on codex) to surface a per-journey `status --read` adoption count, and a scoped-`Read` count, on the journey `Record`.

**Detection (claude.go):**
1. Extend `claudeContent` to unmarshal the tool-use `input` as `json.RawMessage` (kept raw so the struct stays a thin decode of the block; interpretation lives in small helpers).
2. For a `tool_use` whose `Name == "Bash"`, decode `input.command` and increment a `StatusReadCalls` counter when the command string contains a `status --read` invocation. Detection matches the launcher-agnostic forms — `spacedock status --read`, `sd status --read`, `${SPACEDOCK_BIN:-spacedock} status --read`, and a bare `status --read` — by testing for the `status` subcommand followed (in any flag order) by a `--read` token in the command string. The detector is a single helper `commandInvokesStatusRead(cmd string) bool` so the matching rule is unit-testable in isolation and shared by both runtimes.
3. For a `tool_use` whose `Name == "Read"`, decode `input` and increment a `ScopedReadCalls` counter when `offset` OR `lines` is present and non-zero (a full-file Read carries neither).
4. Dedup both counters on the same `toolID` key the existing `toolCallsByName` dedup uses (`claude.go:67-75`), so a repeated stream delta does not double-count — matching the existing tool-call dedup discipline.

**Surfacing on the record:** add two integer fields to `Observation` and `Record` (`StatusReadCalls`, `ScopedReadCalls`) with `omitempty` JSON tags (`status_read_calls`, `scoped_read_calls`). `BuildRecord` copies them from the observation (`record.go:53-55`), exactly as `ToolCallsByName` flows today. No ledger change is required: `AggregateLedger` preserves each `Record` verbatim inside `ScenarioLedgerEntry.Observations`, so the per-release ledger surfaces the counts per observation, and successive release ledgers give the over-time series. "Over time" is the sequence of per-release ledgers, not a new aggregate field — this is the smallest change that satisfies the requirement.

**Codex (codex.go):** mirror the Bash-arg detection against `tool_call.started`'s `arguments.cmd`, reusing the shared `commandInvokesStatusRead` helper, and surface `StatusReadCalls` on `CodexCharacterization`. Codex has no first-class `Read` tool (file reads go through the shell), so `ScopedReadCalls` is claude-only; the codex characterization carries the Bash-derived `status --read` count only. Note that FO/ensign `--read` adoption is principally a Claude-runtime concern (the FO/ensign runtimes drive Claude Code), so the codex mirror is for completeness and symmetry, not the primary signal.

### Per-journey ledger field decision (checklist item 1)
The exact fields that surface `--read` adoption over time:
- `Record.status_read_calls` (int, omitempty) — count of Bash tool calls whose command invokes `status --read`. This is the primary adoption signal: it directly counts the behavior the six contract sites are trying to induce.
- `Record.scoped_read_calls` (int, omitempty) — count of `Read` tool calls carrying `offset`/`lines`. This is the secondary/corroborating signal: the ensign-shared-core contract steers agents to `status --read` *in order to* then issue a scoped `Read`, so scoped-Read volume is the downstream evidence that the hint changed read behavior, not just that the flag was typed.

Both are surfaced; `status_read_calls` is the headline metric for the redundant-instruction-site trim decision.

## Why
The durable behavioral proof that the e6a adoption sticks — and the evidence that should drive the redundant-instruction-site trim (see read-hint-adoption-bloat-trim). Measure before trimming.

## Documentation impact
None. This change adds two `omitempty` integer fields to the emitted journey-metric records and the release ledger JSON. No CLI surface, startup banner, command output, or docs-site-described behavior changes — the journey-metrics record/ledger schema is an internal release artifact, not a documented user surface. No doc diff is owed.

## Acceptance criteria
- **AC-1** — A journey record surfaces a count of `status --read` Bash invocations (`status_read_calls`) and a count of scoped `Read` calls carrying `offset`/`lines` (`scoped_read_calls`). Tested by a journeymetrics unit test over a fixture stream-json transcript that contains one `status --read` Bash `tool_use` and one scoped `Read` `tool_use`, asserting both counts are 1; and a perturbation fixture with a plain `spacedock status` Bash call and a full-file `Read` (no offset/lines), asserting both counts are 0. Non-tautological by perturbation: deleting the detection in `claude.go` makes the first fixture's assertion fail.
- **AC-2** — The `status --read` detector matches the launcher-agnostic command forms and rejects non-`--read` `status` calls. Tested by a table-driven unit test on `commandInvokesStatusRead`: positives include `spacedock status --read PATH`, `sd status --read PATH`, `${SPACEDOCK_BIN:-spacedock} status --read PATH`, and `spacedock status --json --read PATH` (flag-order independence); negatives include `spacedock status`, `spacedock status --boot`, `spacedock dispatch show-stage-def` (a non-status command), and `echo 'status --read'` is accepted-or-rejected per the chosen rule and asserted either way so the boundary is pinned.
- **AC-3** — The counts dedup across repeated stream deltas. Tested by a unit test feeding a multi-delta transcript that repeats the same `status --read` Bash `tool_use` block (same tool id) across two `assistant` rows — mirroring the real multi-delta shape already covered by `TestParseClaudeTurnsMergesToolUseAcrossDeltas` — asserting `status_read_calls == 1`, not 2.
- **AC-4** — The codex characterization surfaces `status_read_calls` from `tool_call.started` `arguments.cmd`. Tested by a unit test over a codex fixture containing a `tool_call.started` whose `arguments.cmd` invokes `status --read`, asserting the count is captured, and zero when the command is a plain `status`.

## Test plan
- **What verifies it:** Go unit tests in `internal/journeymetrics` (`claude_test.go`, `codex_test.go`), driving the existing exported parsers `ParseClaudeJSONL` and `CharacterizeCodexExecJSONL` over new committed `testdata` fixtures, plus a table-driven test of the `commandInvokesStatusRead` helper. No new test framework; extends the established fixture-driven pattern.
- **Fixtures needed:** two new claude stream-json fixtures (one WITH a `status --read` Bash call + a scoped `Read`, one WITHOUT — plain `status` + full-file `Read`) and one codex fixture (or an extension of `codex-exec.jsonl`) with a `status --read` shell call. Fixtures are hand-authored in the proven format confirmed by the spike; the multi-delta dedup fixture follows the shape committed in `claude_test.go`.
- **Cost/complexity:** Low. Pure parser extension over in-repo fixtures; sub-second `go test ./internal/journeymetrics/...`. No live workflow run is needed — the mechanism (command string reachable in the transcript) is already proven by committed real captures, so a fixture exercise is sufficient proof of the parse behavior. A live FO/ensign drive is the *consumer* of this metric (future adoption measurement), not a precondition for proving the parser.
- **Out of scope:** This task measures adoption; it does not trim the six contract sites (that is read-hint-adoption-bloat-trim, gated on this metric's evidence). It does not add a CLI command to display the counts — the counts live in the release ledger JSON, read by `cmd/spacedock-release`.

## Stage Report: ideation

- DONE: Pin the detection design in `internal/journeymetrics` (claude.go + codex.go): detect `status --read` in the Bash `tool_use` command string AND/OR count scoped `Read` calls carrying `offset`/`lines`; decide the exact per-journey ledger field that surfaces `--read` adoption over time.
  Design recorded in `## Proposed approach` and `### Per-journey ledger field decision`: extend `claudeContent` to decode `input` (raw), add `commandInvokesStatusRead` helper shared by both runtimes, surface `Record.status_read_calls` (primary) + `Record.scoped_read_calls` (corroborating) flowing through `BuildRecord` and preserved verbatim per-observation by `AggregateLedger` — over-time = the per-release ledger sequence, no new aggregate field.
- DONE: AC-1's journeymetrics unit test must be NON-TAUTOLOGICAL by perturbation: a fixture transcript WITH a `status --read` Bash call + a scoped `Read` asserts the count is captured; a fixture WITHOUT them asserts zero; removing the detection makes the test fail.
  AC-1 specifies exactly this WITH/WITHOUT fixture pair (counts 1/1 vs 0/0) plus the explicit perturbation clause; AC-2/AC-3/AC-4 pin the matcher boundary, delta-dedup, and codex mirror.
- DONE: Verify the riskiest assumption and record the spike result: confirm the Bash `tool_use` command STRING is actually reachable in the transcript journeymetrics parses.
  `## Spike result` records "no spike needed" with proof from committed real captures (`shallow-boot-greet.stream.jsonl` carries `spacedock status --boot --json` in `input.command`; codex `arguments.cmd` carries the shell command) — the string is on disk; the parser only discards it because `claudeContent` omits `input`. Additionally exercised the end-to-end detection in a throwaway Go program over real-capture-shaped lines: `status_read_calls=2, scoped_read_calls=1` as designed (then removed the throwaway).

### Summary

Fleshed out the ideation: the riskiest assumption (is the Bash command STRING reachable in the parsed transcript?) is proven YES against committed real captures and an exercised throwaway detection run, so the spike resolves to "no spike needed." Design surfaces two per-journey `Record` fields — `status_read_calls` (headline adoption signal, Bash-arg detection via a shared `commandInvokesStatusRead` helper) and `scoped_read_calls` (corroborating, `Read` with offset/lines) — carried through `BuildRecord` and preserved per-observation by the existing ledger, so "over time" needs no new aggregate. Four ACs pin non-tautological WITH/WITHOUT fixtures, the matcher boundary, delta-dedup, and the codex mirror; test plan is low-cost fixture-driven Go unit tests, no live run needed to prove the parser.
