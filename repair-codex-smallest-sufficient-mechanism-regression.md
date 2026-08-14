---
title: Repair Codex smallest-sufficient mechanism regression
status: ideation
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
                state: pending
started: 2026-08-13T22:06:12Z
worktree: .worktrees/spacedock-ensign-repair-codex-smallest-sufficient-mechanism-regression
---

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
