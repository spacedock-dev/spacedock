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

The regression is in the Codex trace classifier. At exact `origin/main`
`177eb454011001a296f1f09bc2889d3436df0a54`, any `agent_message` that contains both
smallest-sufficient vocabulary and a commissioned slug is treated as a per-entity
dispatch justification, even when the message explicitly discusses separate in-house
work. The run's early message says the boot view exposes `ready-one` while the
captain-specified edits are the smallest sufficient mechanism and "outside worker
dispatch". The real `dispatch build` for `ready-one` occurs 29 JSONL events later.
The classifier correlates those unrelated events and creates the false red.

## Read-only spike

The exact public session JSONL from the named run/job was downloaded read-only and
replayed through the pinned parser in a throwaway checkout. The trace was:

- both ladder-note edits observed in-house;
- no edit worker dispatch and no PR observed;
- the direct commit observed;
- no engage recognized from either real `dispatch build` command; and
- `ready-one` falsely marked as carrying a gate justification.

The focused replay failed with `the FO did not dispatch commissioned entity
"ready-one"`; the live durable assertion independently proves both commissioned
journeys and fills engage, after which the retained false justification produces the
reported semantic code. This is the riskiest public-behavior spike: it demonstrates
that the host behavior is correct and that the failure is an event-correlation defect.

## Proposed approach

Make one parser correction, with no product hook or new observability surface:

1. Recognize the existing Codex `spacedock dispatch build ... --entity-path
   .../<slug>.md` command as the commissioned engage event. Preserve the existing
   legacy `status --set ... status=done` and native `spawn_agent` paths.
2. Replace the global `agent_message` substring attribution with ordered correlation.
   Each new agent message replaces the pending narration. A dispatch of a commissioned
   slug evaluates only the pending narration for that same slug, then consumes it.
   Thus the run's later neutral "entering the standing workflow loop" message clears
   the earlier edit explanation before dispatch, while a justification immediately
   framing a dispatch still fails. A `spawn_agent` prompt continues to be graded
   directly because prompt and dispatch are one event.
3. Add a focused regression sequence derived from the exact JSONL and a discriminator
   in which the justification actually frames the dispatch. Replay the exact public
   artifact during validation; do not add a session hook, state store, protocol,
   lifecycle guard, global hook, or replacement compaction hook.
4. Remove only the Codex target binding for this journey after the corrected bytes
   produce XPASS. Keep the Pi binding and every other target/reconciliation row intact.

The simplest alternative, adding a negation regex for "outside worker dispatch", is
insufficient because equivalent truthful wording would regress again. Ignoring all
Codex `agent_message` evidence is also insufficient because it would make a real
per-entity justification invisible. Ordered correlation is the smallest mechanism
that preserves both sides of the behavioral assertion. It serves AC-1 and AC-2.

## Acceptance criteria

- **AC-1 (VALUE):** The unchanged PR #679 Codex session bytes classify as a correct
  smallest-sufficient run: two direct edits, one direct strategy-doc commit, no
  worker/PR climb, two commissioned engages, and zero per-entity justifications.
  **Test:** replay the downloaded JSONL through `codexMechanismTrace` and
  `gradeSmallestSufficientMechanism`; any missing behavior or false justification
  makes the test fail.
- **AC-2:** A real Codex per-entity justification still produces
  `smallest-mechanism-violation`, while unrelated earlier narration that names the same
  entity does not.
  **Test:** paired parser fixtures differ only in whether neutral narration intervenes
  before the same `dispatch build`; the bound case must fail and the separated case
  must pass.
- **AC-3:** The target binding changes only grading, never the observed session bytes:
  corrected semantics are a green XPASS while bound and a normal PASS after unbinding.
  **Test:** feed one immutable replay payload to the bound and unbound grade paths and
  assert `xpass` then `pass`, with identical digest/bytes and no semantic codes.
- **AC-4:** The final unbound targeted local Codex journey passes normally and leaves
  the required one-line strategy doc, both resolved notes, both commissioned durable
  journeys, and a clean tree.
  **Test:** run only `TestLiveCommonSmallestSufficientMechanism` with the isolated-home
  Codex harness and verify normal PASS plus the existing durable assertions.
- **AC-5:** Existing host-neutral and Codex negative cases, the offline suite, and the
  race suite remain green.
  **Test:** focused parser tests, `go test ./...`, and `go test ./... -race`; changing
  true-violation attribution, command grammar, or shared-host behavior makes one of
  these checks fail.

## Expected surface and tolerance

- `internal/ensigncycle/shared_smallest_mechanism_test.go`: about 25 insertions and 8
  replacements for ordered narration/dispatch correlation and `dispatch build`
  recognition.
- `internal/ensigncycle/shared_smallest_mechanism_negative_test.go`: about 35
  insertions for the exact-shape regression, true-violation discriminator, and
  unchanged-byte grade proof.
- `internal/ensigncycle/shared_live_runner_test.go`: one binding removal after XPASS.

Tolerance: these three files only, at most 90 gross insertions and 60 net new lines.
No committed transcript fixture is expected; the public artifact is validation input
and the focused fixture keeps only its behaviorally relevant events. Crossing the
file or line tolerance requires redesign before another correction.

## Declared semantic changes

- **Test classification:** a Codex narration is attributed to a commissioned dispatch
  only through an ordered same-entity dispatch event; real `dispatch build` commands
  count as engage evidence. The exact run changes from false FAIL/XFAIL to XPASS/PASS.
- **Command grammar:** unchanged.
- **Stored formats:** unchanged.
- **Authority and write scope:** unchanged.
- **Runtime/product behavior:** unchanged; Codex already performed the requested work.
- **Documentation:** no site diff. No user-visible command or runtime behavior changes;
  the existing runtime-live registry already states the correct required outcome.

## Verification ladder and correction budget

1. Preserve the baseline spike above and add the paired focused test before the parser
   edit; confirm the separated exact-run shape fails on pinned source truth.
2. Apply the single parser correction. Run the focused positive and negative parser
   tests, then `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.
3. Replay the same captured JSONL bytes with the Codex XFAIL still bound. Require a
   green XPASS alert and zero semantic codes.
4. Remove only this Codex binding and replay those unchanged bytes. Require normal
   PASS; a digest check proves the observation did not change between grading modes.
5. Spend the one allowed revalidation on the exact unbound local command:
   `SPACEDOCK_LIVE_RUNTIME=codex go test -tags live -count=1 -timeout 40m -run
   '^TestLiveCommonSmallestSufficientMechanism$' ./internal/ensigncycle -v`, with
   `SPACEDOCK_CODEX_LIVE_REQUIRED` unset so the existing isolated Codex home is used.

If the discriminator weakens, the unchanged-byte replay remains red, or the one live
revalidation fails semantically, stop. Do not make a second parser-only correction;
redesign the assertion around a different public behavior boundary.

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
