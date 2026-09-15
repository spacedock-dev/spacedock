---
title: Repair Codex smallest-sufficient mechanism regression
status: validation
source: PR #679 run 31728107636, Codex job 94541783359
sprint: test-behavior-completeness
sprint-readiness: ready
score: 0.95
id: bfmczd31ydpp4stqjstf6xwx
gates:
    version: 1
    records:
        - id: gate:bfmczd31ydpp4stqjstf6xwx:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:bfmczd31ydpp4stqjstf6xwx-backlog-1
              briefing:
                id: briefing:bfmczd31ydpp4stqjstf6xwx:backlog:attempt-1:revision-1
                digest: sha256:fbe0c77dc6680f348a746c335c9bb00b50dcf4bac9c42f0d6731a50433a01316
                request-digest: sha256:56165dd9b5981ef952d2a5020d59ce18bbcaff0320dfe6946db4028d4964232c
                room-ref: ./repair-codex-smallest-sufficient-mechanism-regression/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:bfmczd31ydpp4stqjstf6xwx:backlog:1
                briefing: briefing:bfmczd31ydpp4stqjstf6xwx:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-13T22:05:22.190377Z"
                decision: approve
                reason: Captain approved the scoped direction for ideation.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:bfmczd31ydpp4stqjstf6xwx:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:bfmczd31ydpp4stqjstf6xwx-ideation-1
              briefing:
                id: briefing:bfmczd31ydpp4stqjstf6xwx:ideation:attempt-1:revision-1
                digest: sha256:f8cbfa8045523d26368153d05edeff806433a94d74556c933638951b43200058
                request-digest: sha256:f6e1e151ecbae24ed89d7c631e3df2e7dd13ef48e25b45163ea102761a8dcbc5
                room-ref: ./repair-codex-smallest-sufficient-mechanism-regression/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:bfmczd31ydpp4stqjstf6xwx:ideation:1
                briefing: briefing:bfmczd31ydpp4stqjstf6xwx:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-13T22:24:49.414799Z"
                decision: approve
                reason: Captain approved the evidence-backed ordered-correlation repair for implementation.
              application:
                target-stage: implementation
                state: consumed
            - id: gate-attempt:bfmczd31ydpp4stqjstf6xwx-ideation-2
              briefing:
                id: briefing:bfmczd31ydpp4stqjstf6xwx:ideation:attempt-2:revision-1
                digest: sha256:b3a8dca8eacd7c324f43a05f91a63f82133685fcdd9d447e6709a11a1d11aa73
                request-digest: sha256:d7de4a9880e22be3746a1c5abf7389e167088c4702147e69611dd481b4227562
                room-ref: ./repair-codex-smallest-sufficient-mechanism-regression/review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:bfmczd31ydpp4stqjstf6xwx:ideation:2
                briefing: briefing:bfmczd31ydpp4stqjstf6xwx:ideation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-14T05:13:59.771801Z"
                decision: approve
                reason: 'Native lifecycle plus durable-state identity removes shell-command false engages and has a complete fail-closed proof plan; implementation waits for filing PR #686.'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:bfmczd31ydpp4stqjstf6xwx:validation
          stage: validation
          attempts:
            - id: gate-attempt:bfmczd31ydpp4stqjstf6xwx-validation-1
              briefing:
                id: briefing:bfmczd31ydpp4stqjstf6xwx:validation:attempt-1:revision-1
                digest: sha256:4910a1560360b6f864931523b771a1a6b47a700d135499f41b20465269e1088e
                room-ref: '@review/validation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:bfmczd31ydpp4stqjstf6xwx:validation:1
                briefing: briefing:bfmczd31ydpp4stqjstf6xwx:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-09-15T05:10:02.847126Z"
                decision: approve
                reason: Independent replay and exact e52ece67a targeted Codex live pass current four release-attribution criteria. Full normal/race completed with only identical independently reproduced baseline resolver defect; other packages pass and no race findings. Superseded historical criteria do not redefine current approved scope.
                conn:
                    quote: dispatch codex live test failure tasks, verify locally for targeted failure, and open PR as stack, then trigger codex ci on stack tip.
                    source: active thread goal supplied by captain
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:bfmczd31ydpp4stqjstf6xwx-validation-2
              briefing:
                id: briefing:bfmczd31ydpp4stqjstf6xwx:validation:attempt-2:revision-1
                digest: sha256:777c4122a8a45d76afd0e70744c0e30e9215c0dde2584c3146a2b9e3ac6d2a75
                room-ref: '@review/validation/briefing-2'
started: 2026-08-13T22:06:12Z
worktree: .worktrees/spacedock-ensign-repair-codex-smallest-sufficient-mechanism-regression
mod-block:
pr:
---

## Current captain-approved direction — release Python-write attribution (2026-09-14)

This section supersedes the historical candidate, scope estimates, acceptance criteria,
and implementation ladder below for the current repair. Prior reports remain as history;
the rejected commits are not part of this candidate. The captain authorized the current
release false-negative repair through the first officer's implementation dispatch.

- Released workflow: the smallest-sufficient-mechanism live Codex scenario directly
  resolves the two ladder notes, adds the strategy doc, commits, then dispatches ready work.
- Observable harm: pre3 and 0.27.3 reject successful Python parent edits because
  `ssmCodexEditVerbRe` recognizes shell word patterns instead of successful file changes.
- Authority: captain-ruling[2026-09-14] — repair successful parent attribution with real
  content evidence; preserve native worker proof and reject read-only, failed, or delegated edits.
- Trigger evidence: `/tmp/spacedock-pre3-codex-live.zip` and
  `/tmp/spacedock-0273-codex-live.zip`, smallest-mechanism public events `item_10`/`item_8`,
  successful original receipts `809aebc`/`4ea22cd` (all three requested files committed).
  Exact two public events are now durable fixture
  `internal/ensigncycle/testdata/codex_smallest_mechanism_python.jsonl`, SHA-256
  `f7a5642736eb926b6b3205516e3f5605d20a869abca0f4c24342f33d22687be8`.

Current acceptance criteria:

1. Both captured Python command shapes receive in-house edit credit from successful
   parent execution and independently checked exact note changes in real Git history.
   The receipt commit must be strictly after the fixture HEAD captured immediately before the run.
2. Read-only mentions, echoed real receipts, failed/in-progress commands, missing native
   receipts, delegation before the parent commit, wrong/unchanged paths or content, and
   incorrect final file bytes cannot create shell-edit credit.
3. Existing native lifecycle and durable commissioned-journey proof stays unchanged;
   no product instrumentation, new observer framework, or owner binding changes.
4. The narrow layer is committed with focused offline proof. The FO serializes targeted
   local Codex live execution after stacking, then runs full normal/race tests on the
   combined tip. This worker must not spend those runs independently.

Current corrected approach uses the fixture HEAD captured immediately before the run,
a successful public parent command, a matching native receipt before the first spawn,
and a receipt commit strictly descending from that baseline. Exact commit parent/current
note bytes and final files must agree. Under the isolated-fixture/no-prior-worker boundary,
this proves a new parent-segment transaction without recognizing shell invocation text.
The previous candidate at `9029decfc` omitted the baseline; its rejected live evidence
remains tied to that unchanged SHA. Correction `e52ece67a` supplies the baseline and
normalizes the two observed native output shapes: legacy result strings and typed
custom-tool input_text execution results with explicit exit code 0 and a complete
receipt line in the decoded output field. Prose and metadata cannot supply receipts.
Structured file changes must be completed successfully and name the exact file.
The executable captured-output test now explicitly names its scratch branch `main`
to match the immutable native receipt, instead of inheriting the host Git default.
Strict full-receipt matching and all acceptance criteria remain unchanged.
The corrected smallest-mechanism layer is six test/fixture files, +302/-23 lines
against predecessor `5ff85f00` before the later two-line portability correction; no CLI or stored formats change.
The rejected history is preserved at `archive/ssm-rejected-0b809073b-20260914` ->
`0b809073b`; the owned branch was prepared cleanly from `origin/main` `2a7b87198`.

## Historical design and reports (superseded where conflicting)

## Problem statement

PR #679 run `31728107636`, Codex job `94541783359`, reported
`smallest-mechanism-violation`. The seed diagnosis blamed Codex for creating
`roadmap-strategy.md`, but exact source truth disproves that diagnosis: the
`mechanism-choice/mixed-authority` prompt explicitly requires Codex to create that
one-line file and commit it directly. Codex did so, made both deterministic edits,
opened no PR, dispatched both commissioned entities, and left a clean tree.

The regression is in the Codex trace classifier at exact `origin/main`
`177eb454011001a296f1f09bc2889d3436df0a54`. Its public `codex exec --json` stream
does not expose a target on the wait records, so the old classifier did not recognize
the two real commissioned worker lifecycles. The first repair candidate then treated
shell `command_execution` text containing `dispatch build --entity-path` as engage
evidence. Captain rejected that design: an echo, quotation, pseudo command, or failed
command can contain identical bytes without spawning a child.

Engage must therefore be proved from native Codex lifecycle records, then paired with
the durable entity journey. Shell command text is not execution evidence for engage.
The held candidate commits `c8108260b` and `0b809073b` remain historical validation
evidence only; no candidate code or binding changes are authorized before a new
ideation gate.

## Risk evidence and dependency

The exact public session JSONL from the named run/job was downloaded read-only and
replayed through the pinned parser in a throwaway checkout. The trace was:

- both ladder-note edits observed in-house;
- no edit worker dispatch and no PR observed;
- the direct commit observed;
- no engage recognized from either real `dispatch build` command; and
- `ready-one` falsely marked as carrying a gate justification.

The focused replay failed with `the FO did not dispatch commissioned entity
"ready-one"`; the live durable assertion independently proves both commissioned
journeys. This establishes the value gap without accepting command strings as a
repair.

No new mechanism spike is needed before the gate. Existing
`TestCodexIsolatedHomeCollaborationLifecycle` exercises the required public host
records end to end in an isolated Codex home: one public `thread.started`, the unique
parent rollout, structured `spawn_agent` calls, call-ID-matched handle outputs,
`session_meta.source.subagent.thread_spawn` parent/agent identity, waits, assistant
completion, and `task_complete`. A read-only 2026-08-13 session inspection confirmed
the same fields on the current host, including `SubAgentActivity` with
`agent_thread_id` and `agent_path`.

Implementation has one explicit dependency: the filing repair candidate
`ee53f53d2` (or its approved merged successor) supplies
`codexCorrelatedParentRollout`, which binds exactly one public parent thread ID to
exactly one isolated-home parent rollout and errors on missing or ambiguous matches.
Do not cherry-pick, duplicate, or assume that candidate is merged. Implementation
waits for it to land, rebases on its merged API, and requests a new surface decision if
the API differs materially.

## Proposed approach

Make one test-only correction with no product instrumentation:

1. Use the filing observer dependency to resolve the one public parent thread to the
   one parent rollout. Load sibling session JSONL files from the already isolated
   Codex home; do not copy them into a product state store or add an archive hook.
2. For each expected `(entity slug, stage)`, derive the existing dispatch worker name
   (`spacedock-ensign-<slug>-<stage>`). In the parent rollout require exactly one
   `response_item/function_call` in the `agents` namespace whose name is
   `spawn_agent` and whose structured `task_name` equals that expected name.
3. Bind the spawn atomically: its unique `call_id` must have exactly one
   `SubAgentActivity(kind=started)` carrying non-empty `agent_thread_id` and canonical
   `agent_path`, and exactly one `function_call_output` for the same call ID whose
   returned task handle equals that path. Duplicate call IDs, names, handles, started
   records, or child candidates are observation errors, never XFAIL evidence.
4. Bind exactly one child session by all three typed fields: session ID equals
   `agent_thread_id`, `session_meta` parent ID equals the public parent thread, and
   `agent_path` equals the returned handle. Require one initial child task and one
   `task_complete` whose terminal assistant message names the expected entity path and
   stage.
5. Bind targetless waits by exclusivity and order. Between a child's spawn and its
   completion, at most one expected child may be outstanding. Require a structured
   `wait_agent` call after spawn and its matched non-timeout output after child
   `task_complete`; then require the next entity's spawn only after that lifecycle
   closes. Any overlap makes the wait ambiguous and fails closed.
6. Pair that typed lifecycle with the existing durable proof for the same slug/stage:
   a scoped dispatch transition, later path-scoped worker stage report, later terminal
   fields, canonical archive, and clean active path. Lifecycle alone cannot credit
   work; durable state alone cannot credit a native spawn.
7. Attribute a per-entity gate justification only from parent assistant commentary in
   the same typed Codex `turn_id` as that entity's `spawn_agent` call. Earlier prose in
   another turn cannot attach; a true same-turn justification still violates the
   contract.

All shell `command_execution` items and native `exec` inputs are ignored for engage,
whether successful, failed, quoted, echoed, or shaped like a real command. The simplest
alternative, hardening the command regex or checking exit status, cannot distinguish a
successful `echo spacedock dispatch build ...` from a spawn. Using durable state alone
cannot prove the required native dispatch boundary. The correlated lifecycle plus
durable identity is the smallest existing public-behavior proof serving AC-1 and AC-2.

## Acceptance criteria

- **AC-1 (VALUE):** One supported Codex run produces exactly two complete, unambiguous
  lifecycle-plus-durable bundles for `ready-one/ready` and `ready-two/ready`, zero
  shell-derived engages, and zero per-entity justifications, while retaining both
  direct edits, the direct strategy-doc commit, no edit-worker climb, and no PR.
  **Test:** the atomic lifecycle grader returns two identities and the existing durable
  grader returns the same two identities; set equality, count `2`, and zero
  justification codes are asserted together.
- **AC-2:** Failed, quoted, echoed, and pseudo dispatch command text can never create an
  engage, and every incomplete or ambiguous native lifecycle fails closed.
  **Test:** a table removes or duplicates each spawn, call ID, handle, started record,
  child session, wait, wait output, task completion, and durable journey; it also adds
  wrong entity/stage identity, timed-out wait, overlapping children, reordered events,
  and command-text decoys. Every row must fail for its isolated reason.
- **AC-3:** Only same-turn structured narration can justify a commissioned spawn.
  **Test:** the exact PR #679 separated-turn narration passes; moving the same gate
  vocabulary and entity name into the spawn's `turn_id` produces
  `smallest-mechanism-violation` without changing lifecycle or durable bytes.
- **AC-4:** Binding state changes grading only: one immutable observation bundle is a
  green XPASS while the Codex target binding exists and a normal PASS when evaluated
  unbound.
  **Test:** hash the public stream, parent rollout, two child sessions, and durable Git
  state; feed the same digests/bytes to bound and unbound grade calls and require
  `xpass` then `pass` with no semantic codes.
- **AC-5:** The final candidate preserves every other owner and passes focused, full,
  race, and one exact local Codex target.
  **Test:** run focused lifecycle/negative tests, `go test ./...`,
  `go test ./... -race`, then the one allowed isolated-home live target after all
  structured checks pass. Diff the live binding and reconciliation registries to prove
  only this entity's Codex row changed; Sonnet, filing, Pi, and Opus rows are identical.

## Expected surface and tolerance

- `internal/ensigncycle/codex_live_runner_test.go`: about `+35/-5` lines to pass the
  isolated home into the smallest-mechanism observer and surface correlation errors.
- `internal/ensigncycle/shared_smallest_mechanism_test.go`: about `+170/-70` lines to
  replace command-text engage recognition with atomic lifecycle correlation.
- `internal/ensigncycle/shared_smallest_mechanism_negative_test.go`: about `+185/-25`
  lines for the fail-closed matrix, same-turn discriminator, and immutable grade pair.
- `internal/ensigncycle/testdata/codex_smallest_mechanism_lifecycle/public.jsonl`:
  about 8 inserted JSONL records.
- `internal/ensigncycle/testdata/codex_smallest_mechanism_lifecycle/parent-rollout.jsonl`:
  about 12 inserted JSONL records.
- `internal/ensigncycle/testdata/codex_smallest_mechanism_lifecycle/child-ready-one.jsonl`
  and `child-ready-two.jsonl`: about 12 inserted JSONL records total.
- `internal/ensigncycle/shared_live_runner_test.go`: `+1/-1`, removing only this
  entity's Codex binding after XPASS.
- `internal/contractlint/live_registry_reconciliation_test.go`: `+1/-1`, changing only
  the matching reconciliation row after XPASS.

Signed estimate against exact `origin/main`: **+320 net LOC** (approximately 424
insertions and 104 deletions) across exactly nine files. Separate tolerance: **±60 net
LOC**, no more than 500 gross insertions, and no tenth file. The merged filing observer
dependency is not part of this task's diff. Any need to modify its helper, add product
instrumentation, or exceed a tolerance requires a new design gate.

## Declared semantic changes

- **Test classification:** Codex engage exists only when one native spawn/handle/child/
  wait/completion lifecycle and one durable entity journey agree on identity. Shell
  command text never creates engage. Same-turn typed narration alone can attach a gate
  justification. The supported run changes from false XFAIL to XPASS/PASS.
- **Command grammar:** unchanged.
- **Stored formats:** unchanged.
- **Authority and write scope:** unchanged.
- **Runtime/product behavior:** unchanged; Codex already performed the requested work.
- **Documentation:** no site diff. No user-visible command or runtime behavior changes;
  the existing runtime-live registry already states the correct required outcome.

## Verification ladder and correction budget

1. Wait for the filing observer dependency to merge, rebase from exact approved main,
   and confirm all other target rows equal `origin/main`. Do not restore or edit a
   binding as a separate preparatory action.
2. Add the structured happy fixture and every AC-2/AC-3 negative first. Confirm the
   happy lifecycle is unsupported on pinned source truth while every decoy fails.
3. Apply the single allowed correction: replace command-text engage recognition with
   the lifecycle-plus-durable grader. Run focused tests, `gofmt -w ./cmd ./internal`,
   `go test ./...`, and `go test ./... -race`. If any fail-closed discriminator fails,
   stop; do not patch the parser again.
4. With the safety binding still present from main, run the exact local target once:
   `SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 40m -run
   '^TestLiveCommonSmallestSufficientMechanism$' ./internal/ensigncycle -v`, with
   `SPACEDOCK_CODEX_LIVE_REQUIRED` unset. This is the sole live revalidation and must
   emit green XPASS after two structured/durable bundles pass.
5. Before changing source, hash the public stream, correlated parent/child sessions,
   and durable repository state. Evaluate those same in-memory bytes through the
   unbound grade path and require normal PASS with identical digests. Only then remove
   this entity's Codex binding and matching reconciliation row; do not rerun the model
   or touch Sonnet, filing, Pi, Opus, or any other owner row.

If parent correlation is unavailable, lifecycle events are missing or ambiguous, the
bound run is not XPASS, or unchanged bytes are not unbound PASS, stop and return to
design. The correction budget permits no second parser-only repair.

## Stage Report: backlog

- DONE: Record the released workflow and observable harm.
- DONE: Preserve exact run, job, code, and artifact identity.
- DONE: Classify this as a product repair.
- DONE: Define the XFAIL-or-fix decision boundary.

### Summary

Restore smallest-sufficient Codex behavior and keep any temporary expected failure explicit.

## Stage Report: ideation

- DONE: Produce a behavior-first design with a concrete approach, independently testable acceptance criteria, and expected files/insertions/tolerance plus all allowed semantic changes.
  The design declares three files, a 90-insertion tolerance, five falsifiable ACs, and the sole classification semantic change.
- DONE: Diagnose the PR #679 smallest-mechanism violation read-only against exact origin/main 177eb454011001a296f1f09bc2889d3436df0a54, recording the riskiest public-behavior spike or an auditable no-spike basis.
  Exact public JSONL replay reproduced false cross-event justification attribution; the prompt itself proves `roadmap-strategy.md` was required.
- DONE: Define one verification ladder that preserves the test-product firewall, permits only the declared correction/revalidation, and proves bound XFAIL/XPASS then unchanged-byte unbound normal PASS when applicable.
  The five-rung ladder uses existing isolated-home JSONL, one parser correction, immutable-byte grading, and one final live revalidation.

### Summary

The exact run is product-correct; a Codex trace parser correlated unrelated narration with a much later dispatch. The design narrows attribution to ordered same-entity dispatch events, preserves true negatives, removes only the owning XFAIL after XPASS, and forbids a second parser tweak if revalidation fails.

## Stage Report: implementation

- DONE: Deliver the approved ordered same-entity narration/dispatch correlation and real dispatch-build engage recognition within the declared three-file surface and tolerance.
  Commit c8108260b changes two approved files with 54 gross insertions and 41 net new lines; command engages require `dispatch build --entity-path .../<slug>.md` and consume only the latest narration.
- DONE: Add paired falsifiable parser coverage proving the exact separated PR #679 shape passes while a true dispatch-framing justification still fails, without product instrumentation.
  `TestCodexMechanismNarrationCorrelatesWithDispatch` passes the separated sequence and fails if its neutral intervening narration is removed; the focused parser group passes.
- DONE: Commit the candidate and report focused results plus the frozen verification status without spending the one exact live revalidation reserved for the final candidate.
  Candidate c8108260b is committed; `go test ./... -race` passed, the isolated timed-out durable subtest passed on rerun, and the exact live Codex command was not run.

### Summary

The Codex classifier now recognizes real dispatch-build engages and attributes gate narration only when the latest same-entity narration frames that dispatch. The candidate stays inside the approved test-only surface; validation retains the immutable artifact XPASS/unbind ladder and the sole live revalidation.

## Stage Report: validation

- DONE: Independently reproduce AC-1 through AC-5, including the exact immutable PR #679 replay, paired true-violation control, and the bound XPASS then unchanged-byte unbound normal PASS sequence before changing only the owning Codex binding.
  PASSED: run 31728107636/job 94541783359/artifact 9194350789 replay SHA-256 62404d67af406e0e3fc3f364dad10a94dfe6f2648b0621828385d7f1cda4137a (112102 bytes) graded XPASS/codes=[] then PASS/codes=[].
- DONE: Run the required frozen verification once: focused checks, full suite, race suite, detached adversarial audit for the live-test surface, and one exact targeted local Codex smallest-sufficient live run on final bytes.
  Focused and `go test ./...` passed; package-only `go test -race -timeout 20m ./internal/ensigncycle` passed; exact live target passed normally in 405.20s on 0b809073b.
- DONE: Report PASSED or REJECTED with every failure classified from this run, exact evidence identifiers, final diff/tolerance, and confirmation that no other owner binding or reconciliation row changed.
  PASSED: four approved test files, 56 gross insertions/15 deletions, +41 net versus ceilings 90 gross/60 net; only bfmczd31ydpp4stqjstf6xwx's Codex binding and matching reconciliation row changed.

### Acceptance evidence

- AC-1: The immutable replay trace contains both in-house edits, direct commit, no edit dispatch/PR, ready-one and ready-two engages, and zero justifications; changing any field fails the atomic assertion.
- AC-2: The separated narration passes, while removing the neutral intervening narration and repeated later true justification both produce the commissioned-dispatch violation.
- AC-3: The same digest and byte count produced bound XPASS/codes=[] before unbinding and unbound PASS/codes=[] afterward; no observation bytes changed.
- AC-4: `SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 40m -run '^TestLiveCommonSmallestSufficientMechanism$' ./internal/ensigncycle -v` passed normally in 405.20s.
- AC-5: Focused host-neutral/Codex negatives and `go test ./...` passed; the race-instrumented changed package passed in 1185.709s with its three 250ms timing controls also passing focused under `-race`.

### Reviewer findings

- Deferred risk (evidence defect): a non-executing command string such as `echo spacedock dispatch build ...` can counterfeit engage evidence because command success is not parsed.
  Trigger is unobserved and outside the supported replay; the exact public and local supported paths satisfy AC-1/AC-4. Promote if a supported Codex trace emits a failed, quoted, or pseudo dispatch command.
- Infrastructure classification: exact `go test ./... -race` exceeded the default 10m package timeout; the 20m all-package run completed every other package but three unrelated 250ms subprocess controls timed out under load.
  Those three controls passed focused with `-race`, and the isolated changed package passed fully with `-race -timeout 20m`; no race detector finding or changed-surface failure occurred.
- Clean detached audit: separated narration, Unicode/EOF, out-of-order dispatch, missing dispatch, repeated dispatch, and true-framing variants preserved the intended invariant apart from the deferred trigger above.

### Summary

PASSED. Exact public bytes proved XPASS before the owner-only unbind and normal PASS afterward, and the single final live Codex target passed on commit 0b809073b; no material finding remains.

## Stage Report: ideation (cycle 2)

- DONE: Replace shell-command-text engage proof with a behavior-first design based on correlated native Codex lifecycle events and durable entity state.
  The design requires atomic spawn/call-ID/handle/child/wait/completion evidence joined to the existing slug/stage durable journey; command text is never engage evidence.
- DONE: Define fail-closed evidence for spawn, handle, wait, child completion, entity identity, failed commands, quoted commands, pseudo commands, missing events, and ambiguous correlation.
  AC-2 and the seven-step approach enumerate isolated negatives, typed identity, exclusivity, ordering, duplicate rejection, durable completion, and command-text decoys.
- DONE: Revise the expected surface, signed net LOC estimate, acceptance criteria, and one-correction verification ladder for captain review without changing candidate bytes.
  The revision declares nine files, +320 net LOC with ±60 tolerance, five ACs, a filing-observer dependency, bound XPASS/immutable unbound PASS, and one live run; worktree HEAD remains 0b809073b.

### Summary

Cycle 2 removes the rejected command-text architecture and uses only existing public Codex lifecycle and durable workflow evidence. Implementation is gated on the filing observer merge and captain approval; the held code and binding commits were not changed.


### Additional captain-filed evidence — PR #784 Codex attempt 2

Captain requested this evidence be added to this existing repair task; this records a new failure mode without changing the approved design, stage, or candidate.

- Run: https://github.com/spacedock-dev/spacedock/actions/runs/34568401666/attempts/2 ; Codex job 103174402234; candidate `09f123d33b1d7a5b8d3d2dbc85e536739ba7ea1b`.
- Artifact ID `10188347210`, created 2026-09-11T07:03:37Z. Evidence member: `live-artifacts/codex/codex-shared-scenarios/smallest-sufficient-mechanism/codex-exec.jsonl`. Local downloaded archive: `/tmp/pr784-codex-attempt2.zip` (temporary convenience, not the durable source).
- Exact finding: `the FO did not apply the deterministic edit to "ladder-note-alpha.md" in-house — the edit whose content it already held was not made with an in-house Edit`.
- The parent command_execution instead runs Python over `ladder-note-alpha.md` and `ladder-note-beta.md`, asserts the old line occurs once, and calls `p.write_text(s.replace(old, 'Status: RESOLVED'))`. It also writes `roadmap-strategy.md`, then stages all three files and commits. Observed command output: `[main c411281] Resolve ladder notes and add roadmap strategy`, `3 files changed, 3 insertions(+), 2 deletions(-)`; command succeeds.
- At the tested candidate, `internal/ensigncycle/shared_smallest_mechanism_test.go` uses `ssmCodexEditVerbRe` to recognize apply_patch, redirects, tee, or sed -i. It does not recognize this successful Python write. The durable commissioned-journey check passed before the edit-recognition assertion failed. Both ready tasks completed; this is a different evidence gap from the older worker-lifecycle correlation problem above.
- Diagnosis: false negative in recognition of a legitimate in-house edit. Do not repair it by merely adding more command-name regexes: command text alone is not proof of a successful write. Preserve successful result/content attribution and negatives for read-only mentions, failed writes, and delegated edits when evaluating the repair.
- Follow-up proof: replay the captured event against the existing smallest-mechanism proof owner, with independently verified expected file bytes and parent ownership; falsify by removing or failing the actual edit while retaining the command words. No implementation change or test rerun was performed for this evidence filing.

## Stage Report: implementation (cycle 3)

- DONE: Reproduce the release Python-write false negative against current origin/main and repair parent edit attribution with successful execution plus durable content evidence, preserving negative controls and existing native worker proof.
  Commit `e8d39dcd1`; both exact release commands fail against `2a7b87198` and pass when replayed with real fixture commits and matching native parent receipts. Original artifact receipts are preserved in the fixture; regression replay substitutes the freshly created Git receipt, not fabricated commit bytes.
- DONE: Commit the narrow current fix and focused offline regression evidence as a stack layer; prepare the exact targeted Codex smallest-mechanism local run for FO scheduling.
  Four files +183/-23; the current approval/evidence section replaces stale directions explicitly. No owner registry/binding, product implementation, or commissioned lifecycle code changed.
- DONE: Preserve rejected historical work before preparing the current candidate.
  Backup `archive/ssm-rejected-0b809073b-20260914` points to `0b809073b` including `c8108260b`; owned branch reset to `2a7b87198` before this new layer, without force-push.
- DONE: Exercise the focused falsifiable regression controls.
  `go test ./internal/ensigncycle -run 'Test(.*Smallest.*|CodexNativeLifecycle.*|ImplementationLifecycleAndObserverNegativeControls)$' -count=1` passes: omitting the byte comparison admits wrong content; ignoring successful completion admits failed writes; ignoring native spawn order admits delegated-before-parent commits; command-only credit admits read-only/echo receipts. Structured failed/started/wrong-path edits also reject.
- DONE: Verify the live adapter compiles without launching a model.
  `go test -tags live ./internal/ensigncycle -run 'Test(CodexSmallestMechanism.*|AssertCodexSmallestSufficientMechanism|SmallestMechanismTraceSelectsCodexDialect)$' -count=1` passes. `gofmt -w ./cmd ./internal` ran; unrelated pre-existing formatting drift was excluded; `git diff --check` passes.
- SKIPPED: Full normal/race suites and targeted live execution.
  FO explicitly owns serialized live validation and full normal/race at the combined tip. Reserved exact target: `SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 40m -run '^TestLiveCommonSmallestSufficientMechanism$' ./internal/ensigncycle -v`; use FO-provided CI model shim and predecessor tip.

### Summary

Current release Python edits now require a completed exit-0 parent command, a matching native parent receipt before any spawn, and the exact requested committed/final note bytes; command text alone cannot supply edit credit. The focused layer is committed and ready for the FO's stacked live/full validation; this report does not claim those deferred checks passed.

### Stack verification follow-up — 2026-09-14

- DONE: Rebase only the smallest-mechanism layer onto exact predecessor `5ff85f00bac795ad298cddaca09fc51f27484716`.
  Clean new stack tip `9029decfcce628ae6ec82cba9461e7589c1ad421`; original layer `e8d39dcd1` preserved in history/report. Lower layers unchanged; no force-push.
- DONE: Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` sequentially once at the stack tip.
  Both test commands exit 1 solely on `TestCodexResolveManifestAgainstInstalledHost`; all other packages pass (ensigncycle normal 507.994s, race 443.321s), with no race detector finding. The suites are not green.
- DONE: Diagnose and establish the baseline for the host-dependent failure without global plugin/auth/config mutation.
  At exact `2a7b87198` in `/tmp/spacedock-codex-stack-full/baseline-2a7b87198`, the sole focused resolver test also fails with the identical assertion: stable plugin absent but resolver returns installed `spacedock-local/0.28.0-pre0` manifest. Test checks only stable ID; production explicitly prioritizes local ID. FO disposition: keep stack unchanged and record the baseline defect outside this layer.
- DONE: Preserve exact verification environment and logs.
  `/tmp/spacedock-codex-stack-full/summary.md`, source-head files, normal/race/baseline logs, exit files, and SHA-256 manifest; PATH prepends `/Users/clkao/go/bin`, `GOFLAGS=-p=2`, stale `SPACEDOCK_BIN`/`SPACEDOCK_REPO_ROOT` unset. Inherited unrelated gofmt field-alignment drift was reported and discarded with FO approval; diff retained as evidence.
- SKIPPED: Targeted Codex live validation.
  FO explicitly withheld live execution while the local host-negotiation issue is investigated. No lower-layer product fixes or host configuration changes were made.


### Review finding and authorized correction — parent transaction baseline

- Reviewer observation: at exact `9029decfc`, the retained
  `TestCodexSmallestMechanismParentCommit/heredoc_receipt` control prints a `git commit`
  line inside a `cat` heredoc while retaining a real existing commit receipt and file
  bytes. Both edits receive credit. Evidence:
  `/tmp/spacedock-stack-smallest-live-9029dec-validation/adversarial.log` and
  `adversarial.diff`; validator owns the detached reproduction.
- Released user/workflow: smallest-sufficient-mechanism Codex live grading must prove
  the parent made the deterministic edits during its run.
- Observable harm: command text in nonexecuted heredoc data counterfeits invocation;
  the same existing Git commit/receipt can be credited without a current edit.
- Affected authority: captain-ruling[2026-09-14] — current criterion 2 requires that
  echoed real receipts cannot create in-house edit credit.
- Trigger evidence: the detached control fails with both edits `credited=true, want false`.
- Worker proposal: Material, owned grader evidence defect. Snapshot fixture HEAD before
  the run; require a strict descendant receipt commit, exact bytes, and native receipt
  before first spawn. Remove shell invocation matching. The isolated fixture excludes
  external writers; an existing commit is input state, not a new parent transaction.
- First-officer authorization: FIX only this baseline/strict-descendant correction,
  update current approach/proof with no end-value narrowing, retain all native/content
  boundaries, and add executable positive-after-baseline and old-receipt echo/heredoc
  controls. No Trace2, framework, controller, or product instrumentation. Expected core
  change is about 20–35 lines plus focused test changes; report actual surface.
- Execution hold: candidate bytes remain unchanged until live session `33223` exits and
  validator releases the worktree. Then focused red-first tests, fix, and commit only;
  full/live reruns remain FO-scheduled. This entry records authorization, not completion.


## Stage Report: validation (cycle 3)

- DONE: Independently verify the exact smallest-mechanism stack layer against current ACs and adversarial execution/receipt/content controls.
  Reviewed clean candidate `9029decfcce628ae6ec82cba9461e7589c1ad421`, four-file layer +183/-23; focused smallest/native-lifecycle/observer negative group PASS (3.277s). Detached heredoc control fails on both credited note edits.
- DONE: Validate targeted Codex live behavior at exact candidate using the serialized FO-provided lane; preserve failures and request disposition before any changes or retries.
  Exactly one `TestLiveCommonSmallestSufficientMechanism` run at `9029dec`, Luna/max shim and isolated existing auth, FAIL in 205.26s (exit 1), no skip; no retry. FO notified before producer mutation.
- DONE: Report every AC with evidence, exact candidate and test results; retain the proven baseline resolver failure without treating suites as green.
  Current AC-1 and AC-2 have material evidence failures below; AC-3 unchanged native/durable journey code and focused controls pass; AC-4 layer/focused/live/full evidence exists but full suites are not green.
- FAILED: AC-1 — successful Python parent execution receives valid edit credit.
  Captured release fixture positives pass, but actual live parent commit `45a1610` receives no alpha credit: native receipt uses `custom_tool_call_output` with an array of text blocks; decoder accepts only `function_call_output` with string output. Both committed/final notes independently match exact expected bytes.
- FAILED: AC-2 — echoed real receipts cannot create shell-edit credit.
  Detached `TestCodexSmallestMechanismParentCommit/heredoc_receipt` retains real Git/public/native receipt and puts the apparent commit in a `cat` heredoc; both notes incorrectly receive credit. Removing shell-word invocation inference and establishing a pre-run baseline is FO-authorized for the producer.
- DONE: AC-3 — preserve native lifecycle and durable commissioned-journey proof.
  Layer changes no lifecycle/commissioned-journey implementation or owner binding; focused lifecycle and observer negatives pass. Original live advances both ready tasks through worker commits and archive; its sole finding is parent-edit attribution.
- DONE: AC-4 — committed focused proof and FO-serialized verification.
  Candidate binary built; focused group PASS. FO full normal/race both exit 1 solely on baseline-identical `TestCodexResolveManifestAgainstInstalledHost`, evidence `/tmp/spacedock-codex-stack-full/summary.md`; full suites deliberately not repeated.
- DONE: Preserve exact failure evidence without auth or candidate edits.
  `/tmp/spacedock-stack-smallest-live-9029dec-validation/`: focused/live/adversarial logs and exits, adversarial.diff plus detached audit checkout, three preserved-native sessions, workflow-latest.tar.gz, final-workflow, live-diagnosis.txt and evidence.sha256. Parent session `01a0a35e-d7b2-7c92-bf06-001491cba28c`.

### Review-finding disposition

- Material evidence defect, AC-1: released Codex Python-write workflow falsely fails despite exact committed/final content; authority `value-ac[AC-1]` requires successful parent attribution. Trigger: live public commit `45a1610`, native parent zero-based line 38 `custom_tool_call_output` array, first spawn line 76. Decoder drops the supported receipt shape. Owner: this layer; FO consulted, candidate correction authorization belongs to FO.
- Material evidence defect, AC-2: promised echoed-real-receipt control permits heredoc command text to counterfeit invocation; authority `value-ac[AC-2]` explicitly forbids echo credit. Trigger: retained adversarial.diff/log and real fixture receipt. This is a synthetic observation control like the existing echo test, not an observed malicious live action. FO authorized producer pre-run HEAD/strict-descendant repair; no validator candidate changes or reruns.
- Baseline defect, outside layer: installed local plugin makes the existing resolver expectation fail identically at `2a7b87198`; FO disposition retains unchanged stack and honest non-green suite status. No race detector finding was reported.

### Summary

REJECTED at exact `9029decfcce628ae6ec82cba9461e7589c1ad421`: the real live workflow completes its edits and workers but the receipt decoder rejects its supported native shape, and a detached heredoc control admits false edit credit. Both findings were preserved and routed before any candidate mutation or retry; this report evaluates only the original candidate and does not certify the producer's subsequent correction.

Validation disposition follow-up: FO authorized the producer to normalize the observed native `custom_tool_call_output` / `input_text` execution-result shape for AC-1, alongside the pre-run HEAD / strict-descendant correction for AC-2. Authorization excludes a broad parser, other-host changes, or another live run; validator awaits a new candidate SHA and explicit review/live grant. Original `9029dec` remains REJECTED.


## Stage Report: implementation (cycle 4)

- DONE: Repair the owned AC-2 evidence defect with the authorized parent transaction baseline.
  `e52ece67abdb282ab061c41ee897a8a341b22954` snapshots fixture HEAD before launch, requires a strictly later descendant receipt commit, and removes shell invocation/name matching. Executable echo/heredoc controls reuse real pre-run commits and fail; removing the baseline makes the heredoc and missing-baseline tests fail again.
- DONE: Repair the separately authorized AC-1 supported native-output false negative.
  The exact retained custom output envelope now decodes input_text execution-result JSON, requires explicit exit 0, and matches a full receipt line in output. Legacy string output stays supported. Failed/missing-exit, wrong-block/event, metadata-only, and partial-receipt controls reject; reverting normalization breaks the captured-envelope positive.
- DONE: Preserve exact byte evidence and native parent ordering.
  Both release command shapes execute real Python edits/commits after the snapshot. Wrong/unchanged/wrong-path/dirty content and a native spawn before the receipt still reject; native lifecycle controls pass unchanged.
- DONE: Run focused red-first proof, then focused validation and live-adapter compilation without launching a model.
  `/tmp/ssm-baseline-normalization-red.log` records custom-envelope false negative and heredoc/missing-baseline false positives. Focused smallest/native tests pass in 7.319s; live-tag focused compile/tests pass in 7.207s. Green logs are `/tmp/ssm-baseline-normalization-green.log` and `-live-compile.log`.
- DONE: Commit bounded correction and current report for independent validation.
  Correction: five files +160/-41, of which core/carrier/runner +64/-14, executable tests +95/-27, exact native fixture +1. Entire smallest layer versus `5ff85f00`: six files +302/-23. `gofmt` and diff checks complete; previously approved unrelated inherited formatting delta excluded again. Lower layers unchanged.
- SKIPPED: Full normal/race and targeted live reruns.
  FO reserves these for scheduling after independent review. The earlier normal/race evidence remains explicitly tied to `9029decfc` and its known baseline resolver failure; the rejected 205.26s live is not presented as evidence for this correction.

### Summary

The isolated pre-run Git boundary now supplies transaction attribution, eliminating shell parsing and old-receipt credit. Typed normalization repairs the actual native output shape from the rejected live run while keeping explicit execution success, full receipts, exact bytes, and pre-spawn parent attribution; candidate `e52ece67a` is ready for independent validation.


## Stage Report: validation (cycle 4)

- DONE: Verify both authorized AC1/AC2 corrections at e52ece67a preserve exact bytes, baseline ancestry, native receipt success and pre-spawn order; inspect executable negative controls.
  Clean exact candidate `e52ece67abdb282ab061c41ee897a8a341b22954`; independent review confirms pre-launch HEAD capture, strict descendant/current-history checks, exact committed/final bytes, typed successful custom result receipt, and receipt before first native spawn. No candidate edits.
- DONE: After independent review finds no material defect, run the exact targeted Codex smallest-mechanism live test once with the proven isolated environment and retain evidence.
  `TestLiveCommonSmallestSufficientMechanism` PASS normally in 217.15s (package 217.373s), exit 0, no skip/XFAIL/retry. Fresh binary at exact candidate, proven Luna/max shim, cleaned inherited env, isolated existing auth; parent direct commit `1a1f2ec`, both commissioned workers completed.
- DONE: Report current criteria individually and exact candidate checks honestly; prior full/race results belong to 9029decfc, not this candidate.
  AC evidence follows. Producer owns fresh full normal/race at `e52ece67a`, logs `/tmp/spacedock-codex-stack-e52-full`; these remain pending at report time. The old `9029decfc` suite failures are not current-tip coverage.
- DONE: AC-1 — successful Python parent attribution with independently checked exact bytes.
  New detached `TestValidatorFrozenSmallestRun/actual_original_boundary` replays original failed public/native/state bytes and now credits both notes; removing typed custom output normalization makes this fail. New live passes the complete mechanism and durable journey grader.
- DONE: AC-2 — reject old or echoed receipt credit while preserving success/content/order controls.
  The same frozen live bytes reject with baseline equal to receipt commit and baseline later than receipt (`receipt_already_in_input`, `receipt_older_than_input`); removing strict baseline checks makes these fail. Reviewed producer's executed echo/heredoc negatives, failed/missing-exit/wrong-kind/partial/metadata controls, exact-byte and pre-spawn negatives; no repeated producer suite.
- DONE: AC-3 — native lifecycle and durable commissioned-journey proof unchanged.
  Correction changes only test observer/carrier/runner, executable controls and native fixture; owner bindings and lifecycle/journey implementation unchanged. New live passes both commissioned worker journeys and parent mechanism assertion.
- DONE: AC-4 — committed focused proof and serialized current-tip live proof.
  Producer focused/compile logs record green checks at `e52ece67a`; validator added an independent three-case frozen-artifact replay PASS (0.660s) and one granted exact live PASS. Full-suite ownership remains with FO/producer.
- SKIPPED: Repeat full normal/race suites or already-owned green checks.
  Dispatch explicitly reserves fresh full suites to producer and prohibits redundant checks absent new concern. Release completeness still requires recording their current-tip results; this is not a claim that all suites passed.
- DONE: Preserve evidence and exact execution identity without auth.
  `/tmp/spacedock-stack-smallest-live-e52ece67-validation/`: frozen-replay.log/exit, standalone validator_frozen_ssm_test.go, detached audit, live.log/exit, setup source-head, three preserved-native sessions, workflow-latest.tar.gz, and evidence.sha256. Live session 50497 and preservation session 8212 exited.

### Reviewer findings

Both prior Material evidence findings are resolved at `e52ece67a`: the actual custom result is typed and success-checked, and old-receipt credit is excluded by the pre-run transaction boundary. Detached replay and current live provide independent behavior evidence; no new material finding or deferred risk was identified within the approved isolated-fixture/no-external-writer boundary.

### Summary

PASSED for the assigned correction re-review and targeted live validation at exact `e52ece67abdb282ab061c41ee897a8a341b22954`. The single current live run passed normally, and the immutable prior failure now passes only with its legitimate starting boundary; fresh full normal/race results remain separately owned and pending.

### Final corrected-tip full-suite evidence — e52ece67a

- DONE: Run full normal then race suites sequentially once at exact `e52ece67abdb282ab061c41ee897a8a341b22954`.
  `go test ./...` session `21174` exited 1; `go test ./... -race` session `10974` exited 1. Both failed solely on `TestCodexResolveManifestAgainstInstalledHost` at `codex_resolve_test.go:44`; every other package passed. Ensigncycle passed normal in 259.635s and race in 265.386s. No race detector findings.
- DONE: Preserve exact current-tip evidence and distinguish the known baseline failure.
  `/tmp/spacedock-codex-stack-e52-full/summary.md`, source-head file, both logs/exit files, and log SHA-256 manifest. The failure is identical to the already reproduced `2a7b87198` baseline: test checks stable installation only while resolver returns installed `spacedock-local/0.28.0-pre0`. No new failures; baseline was not rerun.
- DONE: Maintain verification scope and candidate integrity.
  PATH prepended `/Users/clkao/go/bin`, `GOFLAGS=-p=2`, stale launcher/repo overrides unset. Candidate stayed clean and unchanged; no host configuration edits, broad reruns, or live tests by this worker. Gofmt had already run on this candidate, and unrelated inherited formatting drift was not reintroduced. Independent live PASS remains separately documented in validation report `b8e0dc086`.
- FAILED: Fully green all-package normal/race commands.
  Both exact commands remain exit 1 solely because of the established baseline resolver expectation defect. This report does not reinterpret those exits as green.


## Stage Report: implementation (cycle 5)

- DONE: Investigate the new offline CI failure without changing the candidate before authorization.
  Published tip `474cdb8c5`, run `34984126222`, job `104431734992`, failed only `TestCodexSmallestMechanismParentCommit/custom_output` (both note edits false). The owned test/grader/fixture bytes match `e52ece67a`; job log retained at `/tmp/ssm-ci-34984126222-failed.log`.
- DONE: Establish and record the owned Material fixture-portability finding and separate FO FIX authorization.
  Supported workflow: offline CI replays the exact native-output positive. Harm: host default branch changes the real receipt to `master` while captured receipt remains `main`. Authority: value-ac[AC-1] — the captured supported-output positive must pass independently of host Git defaults. Trigger: per-process `init.defaultBranch=master` reproduces the exact failure, while `main` passes on unchanged code. FO authorized only explicit scratch branch `main` plus comment; no grader/runtime/contract/assertion change or criteria narrowing.
- DONE: Apply the smallest authorized test setup correction and commit it.
  `0fb8dd634d397e59ba169ec81d5d602594b02800` atop `e52ece67a`: one test file, two added lines directly after `gitInit`; full receipt comparison stays strict. Removing the pin reproduces the master-default failure.
- DONE: Verify this layer locally before any CI retry.
  Full parent matrix under process-local master passes in 8.964s; all focused smallest/native controls under main pass in 8.860s; the same focused controls with `-race` under master pass in 10.367s with no race findings. `gofmt` on the only changed file and `git diff --check` pass. Red differential logs: `/tmp/ssm-ci-default-branch-{master,main}.log`; green logs/exits: `/tmp/ssm-ci-branch-fix/`.
- SKIPPED: Push, CI retry, live test, and broad full-suite reruns.
  FO explicitly withheld delivery/CI/live pending independent validation. This two-line test setup correction changes no runtime or grader behavior. Code and this state report are committed locally for the FO; no push performed by this worker.

### Summary

The captured native envelope legitimately names `main`; the scratch test repository now makes that premise explicit instead of depending on developer Git configuration. Both main/master environment controls and focused race checks pass on `0fb8dd634`, with strict receipt semantics and every existing negative unchanged.


## Stage Report: validation (cycle 5)

- DONE: Independently verify the exact two-line scratch-branch correction0fb8dd634 against CI failure34984126222 and process-local main/master differential evidence.
  Clean exact `0fb8dd634d397e59ba169ec81d5d602594b02800` adds only a comment and scratch `git branch -m main` after `gitInit`; retained CI log and pre-fix master-default reproduction show the identical two false edit assertions, while pre-fix main passes.
- DONE: Confirm runtime, grader, assertions and lower layers unchanged; assess existing passing live evidence applicability without another live run.
  Complete diff versus `e52ece67abdb282ab061c41ee897a8a341b22954` is one offline test file, +2/-0; strict full-receipt comparison and every negative remain unchanged. The prior 217.15s live PASS applies to unchanged runtime/grader behavior, but remains an `e52ece67a` run, not a new live run at `0fb8dd634`.
- DONE: Record PASSED or material findings with exact candidate and current criterion evidence; no code edits, CI or broad reruns.
  PASSED. Reviewed all retained differential/green logs and exit files; no unowned claim required another branch-default probe. Validator ran only read-only source/evidence inspection and diff checks, then wrote this state report.
- DONE: AC-1 — captured successful parent edits receive credit independently of host Git defaults.
  The immutable native receipt names main; explicit scratch branch main makes the real transaction receipt agree without weakening matching. Current parent matrix under master passes (8.964s), focused main passes (8.860s), and focused master race passes (10.367s); deleting the pin restores the recorded CI-equivalent false negative.
- DONE: AC-2 — unsuccessful, old/echoed, delegated, or wrong-content evidence remains rejected.
  Producer's passing full parent matrix includes every existing execution, receipt, baseline, content, path and pre-spawn negative. Only setup branch identity changed; old-receipt and strict native receipt assertions are unchanged, so the pin cannot turn an invalid receipt into credit by weakening the grader.
- DONE: AC-3 — native lifecycle and durable commissioned proof unchanged.
  Exact one-file diff excludes runtime, observer/grader, live runner, lifecycle, owner binding and lower-layer changes. Existing `e52ece67a` live proof remains applicable to those identical bytes.
- DONE: AC-4 — narrow committed layer with focused proof and correctly attributed prior verification.
  Current-tip main/master normal and master race focused checks pass. Existing `e52ece67a` full normal/race each exit 1 solely on the established baseline resolver expectation; no new all-package or live run is claimed at `0fb8dd634`.
- SKIPPED: Duplicate branch probes, live/full suites, code push, and CI retry.
  FO explicitly requested bounded review of already-owned deterministic evidence; no new concern justified repetition. State report alone is path-scoped committed and synced.

### Reviewer findings

The owned Material AC-1 fixture-portability evidence defect is resolved by the authorized scratch-branch pin. CI run `34984126222` / job `104431734992` is matched by `/tmp/ssm-ci-default-branch-master.log`, contrasted with `-main.log`; current green logs and exits are in `/tmp/ssm-ci-branch-fix/`. No criteria narrowing, assertion weakening, new material finding, or deferred risk was identified.

### Summary

PASSED at exact `0fb8dd634d397e59ba169ec81d5d602594b02800`. The two-line fixture correction removes host branch-default dependence while preserving all runtime and evidence semantics; prior live/full-suite results retain their original SHA and outcome labels.
