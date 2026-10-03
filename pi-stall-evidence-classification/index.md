---
title: Classify a stalled Pi live run from the artifacts the runner already archives
status: ideation
score: 0.7
source: "pi-ux carve review, 2026-10-03: stall classification has no Pi implementation owner once the Claude-only capture task stays out of the sprint."
id: z6eb1krpyfmr01dwyb7703hx
sprint:
group: tooling
sprint-readiness: defer
gates:
    version: 1
    records:
        - id: gate:z6eb1krpyfmr01dwyb7703hx:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:z6eb1krpyfmr01dwyb7703hx-backlog-1
              briefing:
                id: briefing:z6eb1krpyfmr01dwyb7703hx:backlog:attempt-1:revision-1
                digest: sha256:df46fdbb84739f3ffc3e48e5b644510c514be915e658607aebeb62c54152990d
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:z6eb1krpyfmr01dwyb7703hx:backlog:1
                briefing: briefing:z6eb1krpyfmr01dwyb7703hx:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-10-03T04:04:28.754736Z"
                decision: approve
                reason: Covers classifying a Pi stall from already-archived artifacts; capture is explicitly not the gap.
                conn:
                    quote: i already said dispatch to ideation, but don't present the ideation gate until staff review finishes
                    source: captain instruction, this session, 2026-10-03
              application:
                target-stage: ideation
                state: consumed
started: 2026-10-03T04:04:59Z
---

A stalled Pi live run leaves artifacts that nobody classifies. The operator cannot
tell transient weather from a defect, so each stall costs a re-run and teaches
nothing. This task records what the existing evidence supports, not a root-cause
verdict that the evidence cannot justify.

## Problem

The common Pi live runner already archives `pi-stdout.txt`, `pi-stderr.txt`, root
session JSONL under `sessions/`, and `process-status.txt`. In
`internal/ensigncycle/pi_shared_live_runner_test.go`, the timeout branch calls
`t.Fatalf` after writing process status but before loading the root session.
No stall classification is currently written. Task `live-ci-api-error-log-capture`
changes the **Claude lane only**; Pi capture is not the gap.

## Proposed approach

For a common Pi journey whose existing per-run cap expires, write
`<artifactDir>/pi-stall-classification.json` after the existing archive writes and
before the existing timeout failure. Here `artifactDir` is the same scenario
folder that already contains `pi-stdout.txt` and `process-status.txt`. The timeout
still fails the journey; this diagnostic never changes grading or retry behavior.
Successful runs and non-timeout failures keep their current behavior.

Use one pure classifier over the already-archived stdout, stderr, root session
JSONL, and process status, including explicit missing/read-error information.
Keep it outside the `live` build tag so ordinary Go tests exercise it. The runner
loads those artifacts and writes the result; the classifier neither controls a
process nor reads external state. Do not use the success-path `onePiSession`
helper to fatal on a missing timeout transcript: incomplete evidence is a normal
classification input.

The output contains `classification`, `evidence`, and `limits`. Evidence names
artifact-relative paths and JSONL line/field or tool-call IDs that support the
result, or names missing/invalid inputs. Limits are nonempty and readable without
opening another document. Do not duplicate raw transcripts or credentials in the
summary. Only three classifications are supported:

- **`assistant-error-observed`**: process status records `timeout=true`, and an
  unambiguous root session's last assistant message has `stopReason:"error"`
  and a nonempty `errorMessage`. Cite that record. This proves an assistant error
  was recorded before timeout, not that a provider outage caused the stall, that
  the error was transient, or that Pi retried incorrectly.
- **`tool-call-pending`**: process status records `timeout=true`, and the root
  session's final assistant tool-call turn contains a `type:"toolCall"` with an
  ID for which no later root `role:"toolResult"` message has the corresponding
  `toolCallId`. Cite the call and any unmatched IDs. A matching error result
  still completes the call. This proves only that completion is absent from the
  archive; it does not prove the tool, worker, or CLI hung or was still running.
- **`inconclusive`**: the default when neither rule applies, required input is
  missing/unreadable, root-session selection is ambiguous, the transcript is
  malformed/truncated or has an unsupported shape/branch, or evidence supports
  conflicting classifications. In particular, silence, a successful-looking
  message, repeated error words in user/tool text, or stderr alone do not prove
  a retry loop or hang. Name the evidence gap rather than invent a verdict.

Both positive classifications require a readable, uniquely selected root session
and consistent timeout evidence. Stdout and stderr remain available as context;
keyword matches in them never override missing structured evidence. An earlier
error followed by a later assistant turn is not the final assistant error.
Missing newer error fields must not crash parsing or be synthesized from prose;
older transcripts can qualify for `tool-call-pending` if their tool-call/result
shape suffices, otherwise the result is `inconclusive` with the shape limit.
Every result states: "Archived Pi diagnostics do not establish whether this stall
was transient weather, a CLI retry-loop defect, or a hang."

This single classifier plus summary serves AC-1. The simplest alternative,
leaving operators to read the archives manually, is today's zero-recorded-result
baseline. A stderr-only keyword label is smaller but mistakes quoted errors for
runtime diagnostics and cannot identify an unmatched tool call. No second
supervisor, state store, background monitor, or capture mechanism is needed.

### Operator artifact contract (M3 fold; proposed for the captain's gate)

These are proposed AC-1/AC-2 interface and edge-case clarifications, not new
classifications or causal claims. Minimal input: an immutable snapshot of the
four archived sources (stdout, stderr, process status, and the immediate
`sessions/*.jsonl` candidates), each carrying artifact-relative path and either
bytes or an explicit missing/read-error state; directory-listing errors are
inputs too. The pure classifier returns one result, never a parse/read exception.
The runner alone performs I/O. No live process inspection or external lookup is
part of the contract. Process status must contain one unambiguous `timeout=true`
flag in the existing `error=... timeout=...` record; missing, false, malformed, or
contradictory timeout evidence makes the result `inconclusive`. Process error text
is not a classification signal. Stdout/stderr are context only: unavailable
context is cited as a gap, never synthesized into structured evidence.

Concrete output (illustrative archived root filename and line):

```json
{
  "classification": "assistant-error-observed",
  "evidence": [
    {"path": "process-status.txt", "line": 1, "field": "timeout", "detail": "timeout=true"},
    {"path": "sessions/root.jsonl", "line": 6, "field": "message.stopReason", "detail": "final assistant stopReason=error; nonempty message.errorMessage"}
  ],
  "limits": [
    "Archived Pi diagnostics do not establish whether this stall was transient weather, a CLI retry-loop defect, or a hang."
  ]
}
```

Field domains: `classification` is exactly one of the three labels above;
`evidence` is a nonempty array of objects with required nonempty strings `path`
and `detail`, optional positive integer `line` (1-based physical JSONL/text line),
optional nonempty string `field` (source field path), and optional nonempty string
`toolCallId` (the exact unmatched ID). Paths are relative to `artifactDir`, never
absolute or traversing `..`; `sessions/` denotes a listing/selection gap. Omit
inapplicable optional fields, rather than emitting nulls or invented locations.
`detail` is a bounded explanation of the observation or missing/invalid input,
not copied error prose, raw content, or OS errors containing host paths/secrets.
`limits` is a nonempty array of nonempty strings, always containing the exact
causal-limit sentence above plus any applicable shape/evidence limitations.
Every unmatched call ID is cited. The result remains the same for the same input;
use artifact-path/line order for evidence and fixed ordering for limits.

Root selection and supported shapes:

- Inspect only immediate regular `sessions/*.jsonl` files. Do not follow symlinks,
  recurse into child archives, pick newest/largest, or merge files. Zero candidates,
  multiple candidates, a listing/read error, or a non-regular candidate is
  `inconclusive`, citing the gap or all candidate paths. A lone candidate must
  begin with a valid `type:"session"` header and must not declare `parentSession`;
  a declared fork/child is unsupported even when it is the only file.
- Support a single linear history: when entry IDs/`parentId` links are present,
  IDs must be unique and each entry must link to the preceding entry (first
  entry's parent is null). Reject broken/forward links, cycles, repeated IDs,
  multiple roots, branches, and mixed linked/unlinked entries. Legacy wholly
  unlinked records can use physical order if their message/tool shapes suffice;
  absent newer error fields still do not become fabricated assistant errors.
- `branch_summary`, `compaction`, or other history-changing/unknown record shapes
  are unsupported, not flattened into a guessed active branch. Known metadata
  (`model_change`, `thinking_level_change`, `session_info`, `label`, `custom`)
  may be non-message context but still participates in link validation. Unknown
  message roles or malformed classification-relevant fields are inconclusive.
  No branch reconstruction or summary interpretation is added.
- Validate every nonblank JSONL record, including records after an apparent
  positive signal. A malformed/truncated record makes the entire result
  `inconclusive` and cites its physical line; never skip it and classify the
  remaining records. Blank lines may be ignored without renumbering. A complete
  final JSON object without a trailing newline is valid; an incomplete object is
  not. An empty/header-only file lacks qualifying evidence. Reader/size-limit
  failures, if encountered, are explicit gaps, not silent partial parses.

Proposed AC-1/AC-3 recording-failure clarification: make one best-effort summary
write before the existing timeout fatal. If serialization, create, write, or
close fails, report a sanitized diagnostic-write warning naming the summary path,
then execute the original timeout failure, with its existing cap and artifact
location. Do not fatal/return early from the summary helper, claim a summary was
recorded, retry the journey, or turn a failed write into a passing/inconclusive
journey outcome. An unwritable destination is an explicit exception to AC-1's
persisted-file guarantee, not an exception to the timeout failure guarantee.
Existing archive-write behavior is unchanged; this clarification covers the new
diagnostic write, not an archive-layer rewrite.

## Risk evidence and dependency

A bounded, local-only format spike on Pi **1.0.0** exercised `pi --print` against
a Python standard-library loopback HTTP server returning HTTP 400 with
`{"error":{"message":"stall-spike-invalid-request","type":"invalid_request_error"}}`.
It used a fresh HOME/PI_CODING_AGENT_DIR, dummy local API key, a custom
`openai-completions` model pointed at that server, retry disabled, `--session-dir`,
and extensions/skills/context/tools disabled. No real provider or credentials
were used. Observed: exit 1; stderr contained the 400 diagnostic; root JSONL line
6 was `type:"message"`, `message.role:"assistant"`,
`message.stopReason:"error"`, with nonempty `message.errorMessage` containing the
same diagnostic. This exercises persistence, not stall classification: it is an
error-format probe, not a stalled lane run or evidence of transient weather.
The sanitized shape and local-server recipe above are the durable spike evidence;
implementation seeds its first fixture from them, without depending on temp files.
Existing `pi_rejection_extractors_test.go` and its tests already exercise Pi root
`toolCall` / `toolResult` matching. No new runtime event stream is assumed.

**Dependency:** `pi-live-lane-pin-refresh` owns moving the pinned family
(pi-coding-agent **0.85.1**, pi-subagents **0.67.0**) forward. The spike does not
validate that older family. The live wiring cannot be validated on the pinned
family; schedule the confirming live run after that pin-refresh task lands.
The classifier must nevertheless tolerate older/missing fields as described
above. Unsupported older diagnostics are a recorded limit, not a reason to add
a Pi debug flag or silently upgrade the lane in this task.

## Expected surface and tolerance

**Proposed for the captain's ideation gate**, replacing the seed's placeholder:
Estimate net LOC change: +220, across 4 files. Insertions ~+225, deletions ~-5.
Tolerance: +/-80 net LOC, +/-1 file.

- `internal/ensigncycle/pi_shared_live_runner_test.go`: timeout-path loading and
  summary write before the existing failure.
- `internal/ensigncycle/pi_stall_classification_impl_test.go`: non-live-tagged pure
  classifier, following this package's existing test-harness helper pattern.
- `internal/ensigncycle/pi_stall_classification_test.go`: deterministic artifact
  fixtures and negative cases; a live-tagged wiring check can live in the runner.
- `docs/runtime-live-ci.md`: the small operator-facing documentation addition below.

**M3 proposed surface clarification for the captain's gate:** also touch
`internal/ensigncycle/pi_live_runner_test.go` solely to attach the deterministic
wiring subtest to its existing registered owner, as specified below. This brings
the planned surface to 5 files (the existing +1-file tolerance); retain the
original +220 net LOC estimate and +/-80 tolerance as the proposed baseline.
No registry or CI-selector change is proposed. Re-estimate at implementation if
the strict parsing/fixtures cannot fit; do not silently widen the tolerance.

Declared observable changes: one diagnostic JSON artifact on common Pi per-run
timeouts. No CLI grammar, entity format, authority, runtime lifecycle, journey
set, grading rule, timeout budget, or XFAIL binding changes. No additional runtime
capture. The non-live split is testability of this one mechanism, not a new service.

## Acceptance criteria

**AC-1 (value) - A capped Pi run leaves one evidence-bounded classification.**
An induced or real common-runner per-run timeout leaves exactly one readable
`<artifactDir>/pi-stall-classification.json`, containing one of the three labels
above, artifact evidence references, and the explicit causal limit. The original
timeout still fails. Measure recorded classifications for the same stall case:
current baseline at `cdfa462d1d426febb2391733507c7fde2143fb0a` is 0; changed runner
must produce 1. A missing file, wrong label, unsupported certainty, or changed
failure outcome fails this criterion. Proof: deterministic runner-wiring fixture
with an induced timeout, followed by one live lane confirmation on a real stall
(after pin refresh; a deliberately capped real Pi run qualifies). Removing the
pre-fatal write must make the fixture fail. No new journey or production timeout
mechanism is required for the induction.

**AC-2 - Each classification is justified by its archived evidence, and gaps stay inconclusive.**
The recorded artifact uses `assistant-error-observed` only for timeout evidence
plus the final structured assistant error (`stopReason:"error"`, nonempty
`errorMessage`); `tool-call-pending` only for timeout evidence plus an unmatched
root tool-call ID in the final assistant tool-call turn; and `inconclusive` for
missing, ambiguous, unsupported, conflicting, or otherwise insufficient evidence.
Every result says archived Pi diagnostics cannot establish transient weather,
a CLI retry-loop defect, or a hang, and cites the supporting artifact or gap.
The detailed evidence rules above apply. Table-driven fixtures write literal
archived inputs and assert the resulting label, cited source, and limit using fixed
expectations, not expectations computed by the classifier. Cover final structured
assistant error; unmatched root tool call; matching success/error tool results;
earlier error followed by progress; missing/empty/malformed/truncated/ambiguous
session; missing process status; unknown/older shapes; conflicting signals; and
error-like words only in stdout, stderr, user content, or tool output. Deleting a
qualifying artifact or replacing it with irrelevant text must change a positive
case to `inconclusive`; restoring it must restore the positive result. A hardcoded
positive label, unchecked substring match, ignored tool result, or fabricated
new-format fields must each turn its corresponding negative case RED.

**AC-3 (no-regression) - Diagnostic recording does not alter journey semantics.**
`go test ./internal/ensigncycle/...` remains green; existing journey assertions,
XFAIL bindings, process-status capture, and timeout/failure outcomes are unchanged.
Successful and non-timeout runs produce no stall artifact. Proof: focused negative
wiring cases plus suite and diff review. Removing the timeout guard must fail the
no-artifact cases; swallowing the original timeout must fail AC-1's exit assertion.
The operator documentation contains the minimal classification contract below.

**M3 proposed AC clarifications for the captain's gate:** AC-1/AC-2 include the
operator artifact contract, root selection, and strict parsing rules above;
AC-1's file guarantee applies when the diagnostic destination is writable.
AC-3 includes preserving the original timeout even when diagnostic writing
fails. All original labels, evidence limits, baseline, and live-proof dependency
remain unchanged. These clarifications await captain approval.

## Test plan

Primary proof owner: the existing common Pi runner
`internal/ensigncycle/pi_shared_live_runner_test.go`; add the non-live fixture suite
beside it. Plain `go test ./internal/ensigncycle/...` proves classification rules
without Pi, auth, or a network. A deterministic `-tags live` wiring check with a
stub child (selected by test name, not the full paid lane) drives the real runner's
timeout/archive/failure path and success/non-timeout controls. Use the existing
run-timeout override; do not introduce a second production deadline mechanism.
`SPACEDOCK_PI_LIVE_TIMEOUT_MINUTES=1` is the existing minimum (positive integer),
so budget roughly 60 seconds for the deterministic timeout wiring case and run
its expected `t.Fatalf` in an isolated test subprocess. Cost: small Go fixture
suite, moderate wiring coverage, no new dependency.

**M3 proposed proof placement for the captain's gate:** place the wiring helper
in `pi_shared_live_runner_test.go`, invoked by a `stall-classification-wiring`
subtest of the existing registered `TestLivePiFrontDoorSmoke` in
`pi_live_runner_test.go`. Place the existing paid smoke body in a sibling subtest,
with Pi/auth/package setup inside that sibling, so the deterministic selector
`go test -tags live ./internal/ensigncycle -run '^TestLivePiFrontDoorSmoke$/^stall-classification-wiring$' -count=1 -timeout=3m`
needs no Pi installation, auth, or provider. Keep the existing smoke assertions
and full-test selection behavior. The wiring helper drives `piSharedLiveDriver.run`
with the stub binary and an isolated expected-failure test subprocess; helper
functions are not new exported `Test...` entries. The registry reconciler in
`internal/contractlint/live_registry_reconciliation_test.go` rejects **every**
new live-tagged top-level `Test...` declaration without registration, including
deterministic tests. Reusing the `pi-front-door-subagent-dispatch` owner in
`docs/runtime-live-ci-registry.md` avoids an unaccounted registry file or new
journey. Run `go test ./internal/contractlint/...` to check that invariant.

Proposed additional fixtures under AC-2 cover candidate enumeration (including
child directories, symlinks, and forks), linked/legacy linear records, branching,
unknown records, and malformed JSONL before/after a qualifying record. Positive
controls with irrelevant metadata and a complete final line lacking a newline
prevent blanket rejection from passing. Fixed expected artifact fields assert
schema domains and sanitized details; accepting a skipped malformed line or
selecting one of multiple roots must turn a negative case RED. Under AC-1/AC-3,
make the diagnostic destination a directory for a deterministic write failure
(no permission/root-user dependence); assert the diagnostic warning AND original
timeout fatal, not merely nonzero exit. Returning early/fataling in the summary
writer must fail that assertion. Keep the existing readable-write 0-to-1 proof.

One live lane run confirms the classification on a **real stall**, after
`pi-live-lane-pin-refresh`; if a normal lane never stalls, deliberately shorten the
existing per-run cap for one actual Pi journey and retain its failed-run archive.
Record run/commit identity, process status, summary, and cited session evidence.
A green run with no stall is not this proof, and an induced cap proves only the
recorded evidence class, not a naturally occurring defect. Budget one selected
journey and the existing cap, not repeated paid reruns to force a positive label;
a justified `inconclusive` is a valid live result. The deterministic fixture owns
the 0-to-1 baseline comparison and every label/negative case (AC-1/AC-2); live proof
owns real-process wiring (AC-1); unchanged harness behavior and suite own AC-3.
All additional checks serve these ACs; their distinct falsifying edits are above.

## Proposed documentation diff

In `docs/runtime-live-ci.md`, immediately after the existing paragraph beginning
"For a custom slow `:max`-thinking model", insert this paragraph (no current text
removed):

> A common Pi journey that hits its per-run cap writes
> `pi-stall-classification.json` beside `pi-stdout.txt` and `pi-stderr.txt`.
> `assistant-error-observed` requires a final structured assistant error in the
> root session; `tool-call-pending` requires a root tool call with no matching
> recorded result. Both require timeout evidence. `inconclusive` means evidence
> is missing, unsupported, ambiguous, or insufficient. The summary cites the
> evidence and its limits; none of these labels establishes transient weather,
> a CLI retry-loop defect, or a hang. Older Pi diagnostics may only support
> `inconclusive`. The original timeout still fails the journey.

**M3 proposed documentation addition for the captain's gate:** append to that
insertion (and include the concrete JSON example above immediately after it):

> The JSON object has `classification` (one of those three strings), `evidence`
> (a nonempty array of objects with artifact-relative `path`, readable `detail`,
> and optional 1-based `line`, source `field`, and unmatched `toolCallId`), and
> `limits` (a nonempty string array including the explicit causal limit).
> Evidence details do not copy raw transcripts or error messages. Classification
> reads the existing stdout, stderr, process status, and immediate
> `sessions/*.jsonl` candidates; stdout/stderr words are never decisive. A unique,
> readable, non-fork root and consistent `timeout=true` are required for a
> positive label. Multiple roots, unsupported branches/history shapes, or any
> malformed nonblank JSONL record yield `inconclusive`, with a cited gap; records
> are not silently skipped. Failure to write the summary emits a diagnostic
> warning and preserves the original timeout failure; it does not promise a
> readable summary when its destination is unwritable.

## Out of scope

Claude-lane capture (`live-ci-api-error-log-capture`), Pi debug flags or capture
changes, pin refresh itself, new journeys, grading/XFAIL changes, automatic retries,
process-tree monitoring/control, child-session reconstruction, and causal diagnosis
beyond the root diagnostics. Only the common Pi runner's existing per-run cap is
covered; an outer `go test`/CI kill before the runner archives is not promised a
summary.

## Stage Report: ideation

- DONE: Task body defines the artifact location, the supported stall classifications, the evidence each requires, and an explicit inconclusive result, all under the exact `## Acceptance criteria` heading.
  AC-2 binds the location and three evidence rules above; missing/ambiguous/older diagnostics default to inconclusive, never a causal weather/retry/hang verdict.
- DONE: A value AC measures an induced or real stall producing a recorded classification within those stated limits, against an independent baseline that can move the wrong way (today no classification is recorded).
  AC-1 measures 0 recorded summaries at cdfa462d1d426febb2391733507c7fde2143fb0a versus 1 after an induced runner timeout, retaining failure; removing the pre-fatal write falsifies it.
- DONE: The test plan names the primary proof owner and states that one live lane run confirms the classification on a real stall; the body promises no causal certainty beyond what Pi's diagnostics support.
  Primary owner is pi_shared_live_runner_test.go; one real-Pi capped journey follows pi-live-lane-pin-refresh, with the causal limit in every result.
- DONE: AC-1 proof plan.
  Deterministic real-runner timeout wiring fixture asserts the persisted artifact and failure; one post-pin-refresh live stall confirms real-process wiring, not causal diagnosis.
- DONE: AC-2 proof plan and diagnostic-format spike.
  Literal positive/negative artifact fixtures falsify fabricated certainty; local Pi 1.0.0 + loopback HTTP 400 produced exit 1 and persisted assistant stopReason:error/errorMessage at root JSONL line 6.
- DONE: AC-3 proof plan.
  Success/non-timeout no-artifact controls, unchanged timeout exit, suite and diff review protect journey/XFAIL semantics; removing the timeout guard falsifies the controls.
- DONE: Proposed surface and operator documentation diff.
  FO approved proposing +220 net LOC (+225/-5), 4 files, +/-80 net LOC and +/-1 file for the captain's gate; body includes the exact docs insertion.
- FAILED: Repository Go validation attempts completed within their local budgets.
  go test ./... and go test ./... -race each timed out at 240s; go test ./internal/ensigncycle/... timed out at 120s. No green suite result is claimed; gofmt -w ./cmd ./internal completed.
- SKIPPED: Implementation fixtures and live stall proof.
  This dispatch is ideation only; classifier/tests remain planned, and live wiring awaits pi-live-lane-pin-refresh from the older 0.85.1/0.67.0 pinned family.

### Summary

Fleshed out one evidence-bounded, pure artifact classifier and its timeout-path summary, with three conservative classifications and a measurable 0-to-1 recorded-result acceptance criterion. A local diagnostic-format spike proved the new Pi error fields; the body explicitly separates that evidence from the pending older-version compatibility and real-stall proof, and proposes the revised surface for the captain's gate.

### Acceptance scan evidence

`spacedock status --read <entity> --ac-scan --json --workflow-dir docs/dev` (selected dispatch binary, exit 0):

```json
{"command":"read","stage":"ideation","acs":[{"id":"AC-1","line":"156","unevidenced":"false","citations":[{"line":"253","text":"  AC-1 measures 0 recorded summaries at cdfa462d1d426febb2391733507c7fde2143fb0a versus 1 after an induced runner timeout, retaining failure; removing the pre-fatal write falsifies it."},{"line":"256","text":"- DONE: AC-1 proof plan."}]},{"id":"AC-2","line":"169","unevidenced":"false","citations":[{"line":"251","text":"  AC-2 binds the location and three evidence rules above; missing/ambiguous/older diagnostics default to inconclusive, never a causal weather/retry/hang verdict."},{"line":"258","text":"- DONE: AC-2 proof plan and diagnostic-format spike."}]},{"id":"AC-3","line":"189","unevidenced":"false","citations":[{"line":"260","text":"- DONE: AC-3 proof plan."}]}]}
```

## Stage Report: ideation (cycle 2)

- DONE: Give the classification artifact one concrete output example with field domains, plus a minimal input/result contract, because it is an operator-facing interface and not a private helper.
  Proposed M3 operator contract defines snapshot inputs, JSON fields, sanitized evidence, and deterministic result; Python parsed the example and checked its domains and exact causal limit.
- DONE: Define root-selection rules, unsupported-branch rules, malformed-JSONL handling that does not silently skip records, and the behavior when writing the diagnostic fails; preserve the original timeout failure rather than letting a summary failure disguise it.
  Proposed AC-1/AC-2/AC-3 clarifications define immediate unique root selection, linear/legacy shapes, strict all-record parsing, and warning-before-original-timeout behavior, including the unwritable-destination exception.
- DONE: Name where the deterministic live-tagged wiring check resides, noting that the registry rejects every new live-tagged test without registration, and use an existing registered owner or explicitly account for that surface.
  Proposed helper in pi_shared_live_runner_test.go is selected through TestLivePiFrontDoorSmoke/stall-classification-wiring; pi_live_runner_test.go is the explicit fifth file within the existing tolerance, with no registry/CI change.
- DONE: Preserve the existing design and mark acceptance-criteria or scope changes as proposed for the captain's gate.
  Original body/report lines were verified preserved in order and frontmatter byte-equivalent; classifications, causal limits, baseline, and pin-refresh dependency are unchanged from the retained 9b3d1ec4b design.
- DONE: Validate registry assumptions and scanner compatibility.
  go test ./internal/contractlint/... -run '^Test.*Registry.*$' -count=1 -timeout=60s -v passed, including TestRuntimeLiveRegistryReconciliation (unregistered live Test declarations fail); selected-binary --ac-scan reports all three ACs unevidenced=false.
- DONE: Run formatting and diff checks without broadening the ideation edit.
  gofmt -w ./cmd ./internal completed; its pre-existing formatting-only change to internal/release/runtime_live_evidence_workflow_test.go was restored. git diff --check passed; no code changes retained.
- FAILED: Full repository validation completed within the local budgets.
  go test ./... and go test ./... -race each exceeded a 120s subprocess budget and were terminated; neither is claimed green. The prior report's timeout evidence remains intact.
- SKIPPED: Implement and execute classifier/wiring fixtures or a live stall.
  Ideation-only M3 fold: fixtures are proposed, not implemented; real-Pi stall proof still follows pi-live-lane-pin-refresh. No new runtime-support claim or spike is made.

### Summary

Added only the staff-review contract and failure/selection details missing from the retained design, plus explicit registered wiring ownership and a concrete documentation extension. All acceptance and surface clarifications remain proposed for the captain's gate; classifier behavior, causal limits, and the deferred real-stall proof are preserved.
