---
title: "Make Pi delegated gate approval continue reliably through successor dispatch"
status: validation
source: "se0 exact-tip Pi recorded-gate proof on 2026-07-28: one run invented the Briefing digest before reading it; the single retry presented the correct bound gate but stopped before recording or consuming the delegated approval."
score: 0.9
sprint: pi-ux
group: gate
sprint-readiness: ready
issue:
id: 9w59t6m1qc46hccd54p04z2j
gates:
    version: 1
    records:
        - id: gate:9w59t6m1qc46hccd54p04z2j:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:9w59t6m1qc46hccd54p04z2j-backlog-1
              briefing:
                id: briefing:9w59t6m1qc46hccd54p04z2j:backlog:attempt-1:revision-1
                digest: sha256:7e4248e69210ba2f494ada2ed8e817d47694c8e53f6ddc9013493fcb20ff9bf6
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:9w59t6m1qc46hccd54p04z2j:backlog:1
                briefing: briefing:9w59t6m1qc46hccd54p04z2j:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-10-03T04:04:11.958238Z"
                decision: approve
                reason: Covers pi delegated gate continuation through successor dispatch; the retained exact-tip run proves a real reliability defect with a stated proof path.
                conn:
                    quote: i already said dispatch to ideation, but don't present the ideation gate until staff review finishes
                    source: captain instruction, this session, 2026-10-03
              application:
                target-stage: ideation
                state: consumed
        - id: gate:9w59t6m1qc46hccd54p04z2j:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:9w59t6m1qc46hccd54p04z2j-ideation-1
              briefing:
                id: briefing:9w59t6m1qc46hccd54p04z2j:ideation:attempt-1:revision-1
                digest: sha256:0ce8fcbb90a800c492647523f030f156c7595b4d83a2820ade0962b05ab2e5d7
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:9w59t6m1qc46hccd54p04z2j:ideation:1
                briefing: briefing:9w59t6m1qc46hccd54p04z2j:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-10-03T18:48:34.03925Z"
                decision: approve
                reason: 'Approve the ideation baseline and advance to implementation. This task measures a fault it does not repair: it restores the Pi recorded-gate grade and demonstrates repeated presentation-through-successor dispatch at one exact tip. The conduct repair belongs to gcmfwfjd9735b58sbzw7xsb8 (repair-pi-recorded-gate-lifecycle), outside this sprint; the captain directive is to build the measurement on 271''s branch.'
                conn:
                    quote: make this work
                    source: captain, pi session 2026-10-03
              application:
                target-stage: implementation
                state: consumed
started: 2026-10-03T04:04:50Z
worktree: .worktrees/spacedock-ensign-pi-delegated-gate-continuation-reliability
---

## Problem

The delegated Pi gate journey can look finished while violating the reviewed
snapshot or stopping before approval application. The retained exact-tip
`ce4ac943` first run presented an invented digest before reading its source,
then completed a successor; the authorized retry displayed the canonical
snapshot but ended at presentation with no child session. Neither is success.

At inspected source tip `cdfa462d1`, the old
`TestLivePiRecordedGateLifecycle` has become the active shared
`TestLiveCommonRecordedGateLifecycle`. Its binding has `nil` gaps: the seed's
linked TODO is already gone. More importantly,
`piSharedLiveDriver.prepareRecordedGate` returns `noLiveGrade`. The shared
scenario grades durable state and command ordering, but does not establish
Pi root-review provenance/read order, requested worker model, or observed
completion followed by the FO's durable-report verification. A green result
from that path alone cannot answer the original reliability question.

## Proposed scope and authority restatements

These changes to the seed are **proposed for captain approval at ideation**, per
supervisor ruling after inspecting the current source; they are not an approved
product change or a claim that the end-value has been delivered:

- AC-1 attaches exactness to the canonical Briefing ID/digest the FO reads and
  binds before decision. `present-gate` permits a compact digest prefix in prose;
  no displayed digest may be invented. Do not change shared `present-gate`.
- AC-4 verifies the already-active shared registry rather than removing an
  absent TODO. No Pi TODO/XFAIL is introduced.
- This task restores the **Pi proof owner only**, not a second gate protocol or
  a conduct remedy. The Pi recorded-gate conduct repair belongs to
  `repair-pi-recorded-gate-lifecycle` (`gcmfwfjd9735b58sbzw7xsb8`, gc/gcm), in
  **pi-live-completeness**, outside this sprint. A green grade here does not
  discharge gc. If the restored grade exposes failed conduct, retain the red
  result and route it there; do not expand into skill/lifecycle repair.

The original end-value remains the measured target: repeated, clean delegated
Pi journeys finish presentation through successor verification. Restoring its
measurement is necessary but is not itself proof of reliable conduct. AC-1
remains unproven until the live threshold below is met; this proposal does not
allow a grade-only implementation to label an uncompleted journey successful.

## Chosen approach

Replace Pi's existing `prepareRecordedGate` no-op callback with a Pi-only
semantic grade of the retained root session and its observed tool results,
then leave the existing shared durable lifecycle assertion intact. Helpers and
deterministic controls stay in `internal/ensigncycle`, not `.pi/extensions/`.
The existing callback receives `liveResult`, including the artifact directory;
Pi's driver already retains the root session there. No shared runner or registry
change is needed or authorized.

The grade observes, rather than scripts, the transaction:

1. Resolve the selected canonical bound ID/full digest from successful gate
   outputs and retained authority, not assistant claims or fixed test constants.
   Require committed binding and a successful canonical read before review.
2. Require one root assistant review before the decision tool call, with entity,
   stage, bound snapshot, recommendation, and actionable decision effect. A
   rendered prefix must agree with the canonical digest. Tool output, child
   text, reasoning, later summaries, and a second corrective review cannot stand
   in for that first captain-facing review.
3. Under the fixture's existing explicit conn, observe successful approval and
   consume, then exactly one successor worker on the requested model. Observe
   that worker's completion and a subsequent FO read verifying its durable
   report, with the existing shared grade still checking commits, state, and
   exactly one durable successor effect. Do not equate dispatch build or a
   final summary with worker completion. Keep existing corrective-build policy
   unchanged; this task adds no retry mechanism.
4. Fail closed on absent or ambiguous evidence. Keep historic obsolete command
   syntax only as retained evidence, never as supported compatibility behavior.

The simplest alternative is leaving the no-op callback and trusting durable
state. The retained first run falsifies it: a successor can exist after a false
review. A prose-only warning cannot make that failure visible to the oracle.
Extending the existing Pi proof owner serves AC-1/AC-3 without new production
mechanisms, CLI verbs, transports, schemas, or runtime orchestration.

## Acceptance criteria

**AC-1 (VALUE, proposed authority/proof-owner restatement).** In a declared batch
of three fresh Pi delegated recorded-gate journeys at one exact candidate tip,
**3/3** complete the full transaction: canonical ID/full digest read and bound
before decision; no invented rendered snapshot; root review before approval;
approval recorded and consumed; exactly one successor on the requested model;
completion observed and durable report verified before final yielding. Compact
canonical prefixes remain permitted. Count every run, including failures and
skips; no selective retries or green-only evidence. The independent fixed
baseline is **0/2 valid complete journeys** from the retained exact-tip first run
and retry below (one invented digest despite completed lifecycle, one early
stop). This measures end conduct, not presence of grading code.

Proof owner: restored Pi callback in live
`TestLiveCommonRecordedGateLifecycle`, plus its unchanged shared durable grade.
**Named falsifying edit `invent-before-read`:** replace the root snapshot digest
with a different digest before the canonical read while leaving successful
consume/successor evidence intact; the run must fail. Live failures keep this AC
unmet and route to gc rather than prompting an in-scope conduct patch.

**AC-2.** The canonical retained package and existing provider-neutral lifecycle
remain authoritative. No hardcoded fixture values in the grader/agent, weakened
oracle, added retries, Pi-only gate protocol, or Claude/Codex grading changes
exist. Shared `present-gate`, shared scenario runner, and registry are unchanged.

Verified by exact diff audit, existing real-CLI/provenance tests, and deterministic
replay against different canonical values read from retained input.
**Named falsifying edit `trust-presenter`:** derive the expected digest from the
assistant's review instead of canonical successful output; the retained wrong
review must still fail, exposing this mutant. Detached adversarial validation
also rejects touching the shared runner or changing another runtime's grade.

**AC-3.** Retained wrong-digest and early-stop runs remain red for the relevant
semantic failure, and independent negative controls reject child/tool-only
presentation, decision-before-presentation, wrong successor model, missing
completion, and missing post-completion report verification. Correct canonical
compact presentation is not falsely rejected for being compact.

Verified deterministically from retained Pi events, preserving role, tool-call
identity/result, and order, rather than authoring a model transcript. Additional
controls make one named edit to retained events and assert the specific failed
obligation, so the historical run's other errors cannot mask the tested claim.
**Named falsifying edits:** `accept-wrong-digest`, `accept-early-stop`,
`accept-child-review`, `accept-tool-review`, `accept-late-review`,
`ignore-model`, `ignore-completion`, and `ignore-report-read` each remove the
corresponding rejection. Each must fail its source-named control. Obtain a real
current passing trace for end-to-end positive control; do not hand-author a
passing transcript or “repair” the historical failure into a claimed live pass.

**AC-4 (proposed registry restatement).** The registered Pi journey remains
active without a Pi TODO/XFAIL binding. Exact-tip evidence includes the focused
three-run batch, the registered Pi common live package, and Pi front-door smoke;
failures are retained and attributed, never relabeled as proof. Existing
Claude/Codex lanes and deterministic gate tests stay active and unchanged.

Verified by registry reconciliation, candidate SHA, command results, and retained
artifacts. **Named falsifying edit `quarantine-pi`:** insert a Pi TODO/XFAIL on
this journey; the binding audit must reject it. A skipped focused journey is not
one of AC-1's three completed runs.

## Expected surface (proposed added proof-owner surface)

Estimate net LOC change: **+500, across 8 files**. Estimate **510 insertions,
10 deletions**. Tolerance: **net +350 to +650**, at most **9 files**, confined to
Pi test support, retained fixtures, and its documentation. Fixture JSONL is
counted in net LOC; the two raw roots are 465,293 bytes/121 lines, so also budget
at most 650 KB retained fixtures to keep line counts from hiding payload growth.
Exceeding these bounds requires renewed design approval, not silent expansion.

Expected files:

- `internal/ensigncycle/pi_shared_live_runner_test.go`: wire its existing callback.
- `internal/ensigncycle/pi_recorded_gate_grade_test.go`: Pi-only observation grade.
- `internal/ensigncycle/pi_recorded_gate_grade_test_test.go`: deterministic controls.
- `internal/ensigncycle/testdata/pi-recorded-gate/first-root.jsonl`.
- `internal/ensigncycle/testdata/pi-recorded-gate/retry-root.jsonl`.
- `internal/ensigncycle/testdata/pi-recorded-gate/first-child.jsonl`.
- `internal/ensigncycle/testdata/pi-recorded-gate/README.md`: provenance, checksums,
  source coordinates, any necessary redactions, and reproduction commands.
- `docs/runtime-support.md`: proof boundary paragraph below.

The ninth-file tolerance permits a retained **real** positive trace, not another
runtime or production seam. Entity body/report are separate state-checkout
planning artifacts, not code-surface LOC. No production command grammar,
stored format, decision authority, or runtime behavior changes are authorized.
Observable test semantics change: Pi falsely green incomplete/misattributed
journeys become red. No edits to shared runner, shared registry, Claude/Codex,
`present-gate`, FO scaffolding, dispatch scheduler, or fixture prompts.

## Test plan and risk evidence

**Primary proof owner: live**, `TestLiveCommonRecordedGateLifecycle` with
`SPACEDOCK_LIVE_RUNTIME=pi`. Deterministic replay is the oracle's proof, not a
replacement for agent conduct. Add focused replay tests before wiring the
callback. Use the retained first run/retry and parameterized single-obligation
mutants, plus a real current positive when available. Never teach the agent the
fixture's transcript, values, or expected phrasing.

**Spike performed (deterministic trace inspection, not a live success).** Parsed
both actual root JSONL files with Python's standard JSON parser, retaining event
roles and array order. The first root's line 39 contains the wrong captain-facing
digest; its canonical review read arrives at line 42. It later records/consumes,
spawns at 70, receives completed wait at 73, and reads the durable entity at 76.
The retry ends at line 42's review, with no child file. This proves the necessary
counterexample observations exist in retained Pi events. The current driver
retains root JSONL and passes artifactDir into its existing grading callback;
there is no proposed unverified host API. Semantic-grade correctness and current
live success remain implementation/validation obligations, not spike results.

**Commands for implementation/validation** (from candidate checkout; Pi auth and
package prerequisites are those declared in `docs/runtime-support.md` and CI):

- Deterministic: `go test ./internal/ensigncycle -run 'TestPiRecordedGate|TestRecordedGateLifecycle' -count=1`.
- Registry: `go test ./internal/contractlint -run '^TestRuntimeLiveRegistryReconciliation$' -count=1`.
- Focused live: `SPACEDOCK_LIVE_RUNTIME=pi go test -tags live ./internal/ensigncycle -run '^TestLiveCommonRecordedGateLifecycle$' -count=3 -v -timeout 45m`.
- Registered lane: `SPACEDOCK_LIVE_RUNTIME=pi go test -tags live ./internal/ensigncycle -run '^TestLiveCommon' -count=1 -parallel 4 -timeout 50m`.
- Front door: `go test -tags live ./internal/ensigncycle -run '^TestLivePiFrontDoorSmoke$' -count=1 -timeout 15m`.
- Regression: `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, and `go vet -tags live ./internal/ensigncycle`.

Set `SPACEDOCK_LIVE_ARTIFACT_DIR` to a distinct retained directory per batch;
record SHA, Pi/model/package versions, full commands, exits, root and child
sessions, and durable state/commit evidence. Credential/setup skips are not
passes. Estimated work: medium (one Pi grading seam and adversarial controls);
deterministic tests take seconds to about a minute, focused live up to 36 minutes
at the default 12-minute cap/run, registered package up to 50 minutes. Do not
launch extra attempts to replace a failed sample. A detached adversarial audit
at validation challenges necessity, false positives, role/order spoofing,
fixture dependence, and scope; no autonomous review fanout is part of this task.

Ideation checks actually run: focused real-CLI replay, provenance mutants, and
missing-event controls passed (27.379s); registry reconciliation passed (0.528s).
No live lane was run during ideation. These checks validate available mechanisms
and the proposed test owner; they do not establish AC-1.

## Proposed documentation diff

No product/skill instruction wording changes: this proposal changes proof only.
In `docs/runtime-support.md`, after the existing durable-outcome list under
“Pi live-smoke mechanism”, add this paragraph (before: no corresponding paragraph):

> The Pi execution of `TestLiveCommonRecordedGateLifecycle` additionally grades
> the root-session review against canonical bound authority before decision,
> then observes approval application, the requested successor model, completion,
> and the first officer's durable-report verification. A correct compact digest
> prefix is allowed in review prose; an invented value, child/tool-only review,
> or presentation followed by an early stop is not. Retained negative traces
> validate the grader; only repeated clean live journeys establish conduct.

Implementation applies this wording only when the restored grade actually owns
these checks. Do not advertise the current no-op path as providing them.

## Boundaries and routed findings

- `w5bfnrvpcphw857nzz93340c` owns exact-digest presentation reliability.
  Current `present-gate` explicitly allows compact prefixes; no wording repair
  to that shared surface is proposed here.
- `gqsw81ghf48hr2n3jg6k7nx8` is already archived and owns entered-stage dispatch.
  Reuse its scheduler/real-CLI coverage; do not reopen its implementation.
- Route the following **conduct finding to gc**, without editing scaffolding:
  Pi adapter's async re-review paragraph says
  “`gate prepare` → `state commit` → `status --read` → `present-gate` → stop”,
  while shared `fo-gate-lifecycle` says “An existing conn requires present once,
  immediate record, then the route below.” That unconditional adapter stop
  conflicts with delegated continuation. It is not established as the cause of
  the historical prepared-gate failure and is not grounds to expand this
  grade-only proposal. Any later scaffolding remedy needs its own approved scope
  and detached adversarial validation.
- If Pi-only grading cannot satisfy these obligations without shared changes,
  stop and report the proof gap; do not change other runtimes' oracles.

## External dependency handoff (M2, proposed for captain's gate)

This staff-review fold supplements, rather than replaces, the proof-owner and
AC restatements in `a95ecf5dd` and `5436dde75`. **Proposed for captain approval:**
this task cannot repair conduct. AC-1's unchanged **3/3 at one exact candidate
tip** depends on conduct repair owned by `repair-pi-recorded-gate-lifecycle`
(`gcmfwfjd9735b58sbzw7xsb8`, gc), in **pi-live-completeness**, outside this sprint.
Coordinate that repair against this task's restored grade; retain red evidence
and route failures to gc, never attempt the repair here or hide it with a binding.

The presentation dependency on `reliable-exact-digest-in-gate-review`
(`w5bfnrvpcphw857nzz93340c`, w5) is discharged for this task by evidence under
the approved compact-prefix contract: the FO reads and binds the exact canonical
Briefing ID/full digest before decision, while the captain-facing prose may
render a matching compact digest prefix. The restored grade and its controls
must accept that rendering and reject invented or vague snapshot references;
the unchanged live batch must prove the complete transaction. This does not
import w5's older full-displayed-digest requirement, claim w5 is closed, or
commission w5 inside this sprint. Any remaining presentation conduct failure
needs an external owner handoff, not an in-scope repair.

Before Commander acceptance and **before any merge requiring a green Pi lane**
(including earlier launcher merges, not just this task), the external handoff
must record accepted code SHAs, the responsible repair owner, candidate and owner
availability, and retrievable exact-tip closure evidence. That evidence must
show the unbound Pi recorded-gate failure repaired and presentation satisfied
under the compact-prefix distinction above, coordinated with the restored grade;
required green-lane evidence and this task's declared three-run batch remain due.
No accepted repair SHAs, availability, or closure evidence are supplied by this
fold: the dependency remains open, not implicitly commissioned or satisfied.
If the handoff is unavailable, stop and escalate for the captain to revise scope
or the promised outcome; do not accept grade-only work as reliability or add
replacement runs to the declared batch.

## Retained evidence and portability

Artifacts are **absent from this checkout's root `live-artifacts/` path**, but
were found and actually inspected in the local old worktree beneath:
`/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-live-lanes-red-on-every-branch/live-artifacts/local-proof/`.

- First: `final-ce4ac943-pi-recorded-gate-pinned/pi-recorded-gate-lifecycle/run/sessions/2026-07-28T15-57-34-288Z_019fa972-16d0-7881-8a94-421ffb1695d3.jsonl`.
  SHA256 `4ecc18637c62143b8cbae7fbf584fa3093145aa3f6d8d46daa31f4343dfe85de`;
  79 lines/249,019 bytes. Review line 39 names `sha256:ec12b8…`; the actual
  canonical value is `sha256:0a54f1…`. Same event acknowledges invention before
  reading. Its go-test log fails “gate review omits its decision facts”.
- Retry: `final-ce4ac943-pi-recorded-gate-pinned-retry/pi-recorded-gate-lifecycle/run/sessions/2026-07-28T16-03-39-912Z_019fa977-ab08-773b-ac42-334e8babde6b.jsonl`.
  SHA256 `b25a248167a200b7a5ac80ac7034fa2fbdb1377c236d52fd32389ddcea5dfe72`;
  42 lines/216,274 bytes. It reads canonical sources and ends at the review;
  its go-test log fails “Pi child sessions=0, want exactly one”.
- First child: under the first run's same `sessions/` directory,
  `2026-07-28T15-57-34-288Z_019fa972-16d0-7881-8a94-421ffb1695d3/8a89cecb/run-0/session.jsonl`
  (21 lines). Retry has no corresponding child.
- Both stderr files contain launcher banners, not the session trace; only the
  first adds an unacknowledged async intercom advisory. Do not use stderr alone
  as root-review proof.

Before tests depend on these inputs, implementation must retain them under the
listed repository testdata paths with provenance and checksums (and document any
redaction). Do not depend on this local worktree or ephemeral `/var/folders/`
paths during validation. If inputs cannot be retained safely, report a blocker;
no fabricated replacement trace. Historical commands are negative evidence,
not modern CLI compatibility requirements.

## Stage Report: ideation

- DONE: Task body states the problem, the chosen approach, and criteria under the exact `## Acceptance criteria` heading, with every AC paired to a named falsifying edit.
  Proposed Pi-only grade restoration at the existing callback; AC-1–4 each name proof and mutants, with no production/scaffolding changes.
  AC-2 proof plan: exact diff audit, existing real-CLI/provenance tests, and replay with different retained canonical values preserve provider-neutral authority; `trust-presenter` must fail, and detached audit rejects shared-runner or other-runtime grade changes.
  AC-3 proof plan: retained wrong-digest/early-stop runs and single-obligation role/order, model, completion, and report-read controls must reject their named mutants; a real current passing trace supplies the positive control, including canonical compact presentation.
  AC-4 proof plan: registry reconciliation rejects `quarantine-pi`; retain exact candidate SHA, all three focused live runs, registered Pi common package and front-door smoke results/artifacts, while Claude/Codex lanes and deterministic gate tests remain active and unchanged.
- DONE: A value AC measures the end outcome - a repeated Pi recorded-gate journey completes presentation through successor dispatch - against an independent baseline that can move the wrong way (the retained exact-tip run whose digest was invented before it was read).
  AC-1 requires 3/3 complete live transactions against the inspected 0/2 retained baseline; line 39's invented digest precedes canonical read at 42, and retry ends with zero children.
- DONE: Expected surface declares net LOC change, files, and tolerance; the test plan names the primary proof owner and states whether that proof is a live lane or deterministic.
  Proposed +500 net LOC (510 insertions/10 deletions), eight files, +350..+650 and nine-file ceiling; primary proof is live, replay controls deterministic.
- DONE: Inspect current contracts, proof owner, and retained evidence before selecting the approach.
  At cdfa462d1 Pi prepareRecordedGate is a no-op; actual retained JSONL exposes review/read order, successful wait, and post-completion report read; focused replay/provenance/missing-event and registry tests passed.
- DONE: Record authority, registry, and proof-owner restatements as proposed for captain approval.
  Exactness attaches to canonical authority, compact prose stays allowed, absent TODO is not recreated, and gc owns out-of-sprint conduct repair; a green grade here does not discharge gc.
- DONE: Record concrete documentation change and cross-task boundaries.
  docs/runtime-support.md addition is quoted; shared present-gate and runner remain untouched; unconditional Pi adapter stop is routed as a gc finding, not an in-scope remedy.
- SKIPPED: Execute the candidate live batch and registered Pi package during ideation.
  This stage defines the proof; no candidate grade exists yet, and offline checks do not establish live conduct or satisfy AC-1.

### Summary

Shaped the task into a proposed Pi-only proof restoration after finding the current journey's no-op grading seam and inspecting both retained counterexamples. The original end-value remains a three-run live requirement, not a claim that tests or prose repaired conduct; captain approval of the restated ACs and scope, followed by independent staff review, is still required.

### Report repair verification

`spacedock status --read docs/dev/.spacedock-state/pi-delegated-gate-continuation-reliability.md --ac-scan --json --workflow-dir docs/dev` (exit 0):

```json
{"command":"read","stage":"ideation","acs":[{"id":"AC-1","line":"119","unevidenced":"false","citations":[{"line":"320","text":"  Proposed Pi-only grade restoration at the existing callback; AC-1–4 each name proof and mutants, with no production/scaffolding changes."},{"line":"325","text":"  AC-1 requires 3/3 complete live transactions against the inspected 0/2 retained baseline; line 39's invented digest precedes canonical read at 42, and retry ends with zero children."},{"line":"335","text":"  This stage defines the proof; no candidate grade exists yet, and offline checks do not establish live conduct or satisfy AC-1."}]},{"id":"AC-2","line":"138","unevidenced":"false","citations":[{"line":"321","text":"  AC-2 proof plan: exact diff audit, existing real-CLI/provenance tests, and replay with different retained canonical values preserve provider-neutral authority; `trust-presenter` must fail, and detached audit rejects shared-runner or other-runtime grade changes."}]},{"id":"AC-3","line":"150","unevidenced":"false","citations":[{"line":"322","text":"  AC-3 proof plan: retained wrong-digest/early-stop runs and single-obligation role/order, model, completion, and report-read controls must reject their named mutants; a real current passing trace supplies the positive control, including canonical compact presentation."}]},{"id":"AC-4","line":"167","unevidenced":"false","citations":[{"line":"323","text":"  AC-4 proof plan: registry reconciliation rejects `quarantine-pi`; retain exact candidate SHA, all three focused live runs, registered Pi common package and front-door smoke results/artifacts, while Claude/Codex lanes and deterministic gate tests remain active and unchanged."}]}]}
```


## Stage Report: ideation (cycle 2)

- DONE: Record that this task cannot repair conduct: its 3/3 requirement depends on conduct repair owned outside the sprint, and the task must not attempt it.
  Proposed M2 handoff section names gc in pi-live-completeness as external conduct owner; existing proof-only scope, ACs, and three-run batch are unchanged.
- DONE: Record how the presentation dependency is discharged under the approved compact-prefix contract, distinguishing canonical authority exactness from rendered prose.
  M2 section requires exact canonical ID/full digest authority and matching compact-prefix rendering, with restored-grade controls and live proof; it neither imports w5's older display rule nor claims w5 closure.
- DONE: Name the external handoff requirement in the body: accepted code SHAs, responsible owner, availability, and closure evidence before any merge that requires a green Pi lane.
  M2 section requires all four before Commander acceptance and any green-Pi-lane merge, including earlier launcher merges; absent handoff remains an explicit open dependency for captain escalation.
- SKIPPED: Repair conduct or execute the candidate live batch.
  This is only the M2 ideation fold; no conduct implementation is authorized and no live result is claimed.
- DONE: Check the planning artifact without claiming conduct proof.
  Reader/AC scan checks AC-1, AC-2, AC-3, and AC-4; exact comparison preserves all original text and frontmatter, allowing only the M2 section and this report; git diff --check passes.
  Standing proof plans from 5436dde75: AC-1 remains the 3/3 live threshold; AC-2 retains the trust-presenter/diff audit; AC-3 retains semantic negative controls and a real positive trace; AC-4 retains quarantine-pi and exact-tip lane evidence. None is newly proven by this fold.
- FAILED: Complete repository-wide regression checks during this documentation-only fold.
  go test ./... and go test ./... -race each exceeded the 240-second tool budget; no pass claimed. gofmt completed; its unrelated pre-existing formatting delta was restored.

### Summary

Added only the proposed M2 external-dependency, compact-prefix, and handoff clarification for the captain's gate. The committed proof-owner restoration and AC restatements stand unchanged; external repair acceptance and closure evidence remain outstanding.

## Stage Report: implementation

- DONE: Finish the restoration of the Pi recorded-gate grade in the Pi driver only; keep the shared scenario runner/grader, shared present-gate, registry, and Claude/Codex grades untouched.
  Commit d70bd32d6 wires `piSharedLiveDriver.prepareRecordedGate` (`internal/ensigncycle/pi_shared_live_runner_test.go`) to `piRecordedGateFindings`; `git diff --cached --name-only` is exactly the declared 8 files. No shared runner/grader, `present-gate`, registry, or Claude/Codex grade file is in the set.
- DONE: LOC-surface decision (bounded expansion): keep the full AC-3 proof and record actual net LOC, stating plainly it exceeds the +500 estimate.
  Actual 845 insertions, 2 deletions, **net +843 across the declared 8 files**. This exceeds the +500 ideation estimate and the +650 net tolerance: the retained raw session traces (first-root 79/249,019; retry 42/216,274; child 21/21,619 bytes = 486,912 bytes JSONL, under the 650 KB cap) plus the full eight-mutant single-obligation table are counted in net LOC. No controls were cut to fit the estimate; the overage is estimate error, not scope creep (semantics stay proof-only; no production grammar/format/authority/runtime change).
- DONE: Verify the deterministic proof.
  `go test ./internal/ensigncycle -run 'TestPiRecordedGate|TestRecordedGateLifecycle' -count=1` ok 42.7s; `go test ./internal/contractlint -run '^TestRuntimeLiveRegistryReconciliation$' -count=1` ok 0.33s; `gofmt -l` on the three changed Go files empty; `go vet -tags live ./internal/ensigncycle` clean; `go build ./...` clean; `go test ./...` ok (ensigncycle 282.8s, all packages pass). No live lane run (CI Pi OAuth blocked, refresh_token_reused). A pre-existing `gofmt -l` delta in `internal/release/runtime_live_evidence_workflow_test.go` is outside this surface and was left untouched.
- DONE: Write this `## Stage Report: implementation` with exact changed files, actual net LOC, and AC-1..AC-4; commit code on the branch and push both.
  Code commit d70bd32d6 pushed to `origin/spacedock-ensign/pi-delegated-gate-continuation-reliability`; this state report committed path-scoped and pushed to `spacedock-state/dev`.

AC citations:

- **AC-1 (VALUE, live 3/3) — NOT met.** No live lane was available (Pi auth blocked), so the declared three-run batch did not run. The proof owner is restored at the existing callback and the deterministic controls pass, but only repeated clean live journeys establish conduct; this remains red and routes to gc (`gcmfwfjd9735b58sbzw7xsb8`), not an in-scope repair.
- **AC-2 — satisfied deterministically.** `TestPiRecordedGateFollowsCanonicalValues` fails if the grader hardcodes the retained digest (it rewrites canonical id+digest and expects no mismatch); `TestPiRecordedGateRetainedRootsStayRed` fails a `trust-presenter` grader that derives the expected digest from the review; the exact diff audit shows shared `present-gate`, runner, and registry unchanged.
- **AC-3 — satisfied deterministically.** `TestPiRecordedGateRetainedRootsStayRed` asserts first-root red only for invented-digest/before-authority and retry red only for approval-not-recorded. `TestPiRecordedGateControls` runs the eight named mutants (`accept-wrong-digest`, `accept-early-stop`, `accept-child-review`, `accept-tool-review`, `accept-late-review`, `ignore-model`, `ignore-completion`, `ignore-report-read`): each asserts the full grade names its obligation and the single-obligation mutant no longer reports it. `TestPiRecordedGateCompactDigestPrefixAccepted` fails a full-digest-only grader. `TestPiRecordedGateCleanObservationPasses` is a unit-level positive control at the judge seam (no live passing trace exists in this checkout). `TestPiRecordedGateArtifactLoader` pins live session discovery.
- **AC-4 — satisfied deterministically.** `TestRuntimeLiveRegistryReconciliation` passes; no Pi TODO/XFAIL binding added. The focused live batch, registered Pi common package, and front-door smoke remain due and are NOT claimed.


### Deliverable change (cycle 6): 9w removes a divergent Pi judge

**9w no longer restores a Pi grade. It removes a divergent one. Its
presentation-reliability acceptance stays OPEN.** Code commit `fb4428e9c` reverts
the judge commit `d70bd32d6`.

- Deleted `internal/ensigncycle/pi_recorded_gate_grade_test.go` and its companion
  `pi_recorded_gate_grade_test_test.go` (the companion only served this judge),
  the fixtures it alone consumed (`internal/ensigncycle/testdata/pi-recorded-gate/`),
  and the `docs/runtime-support.md` paragraph that advertised it.
- Restored `piSharedLiveDriver.prepareRecordedGate` in
  `internal/ensigncycle/pi_shared_live_runner_test.go` to the no-op
  (`return d, noLiveGrade`) so the shared assertion runs.
- The shared assertion itself is unchanged.

Why the judge is divergent (captain's evidence):

- it requires the legacy `gate record --briefing` shape, while the shared
  contract at `internal/ensigncycle/recorded_gate_lifecycle_test.go:19-47`
  explicitly accepts the collapsed `gate record --decision approve --consume`;
- its digest detector searches presentation prose instead of the prepare
  output's own fields;
- its review pattern rejects the presentation style that `present-gate` /
  `SKILL.md:13-16` allows;
- it compares the spawn argument against the root model instead of the child's
  effective model, and the child used the expected model;
- it calls `t.Fatalf` before the XFAIL handling runs, so no binding can defer it.

### Live measurement history (produced by the now-removed judge)

Recorded for the history; the recorded-gate red they show was produced by the
divergent judge, so it is not a product verdict. With the judge removed, the
shared assertion runs; whether it passes is unknown until the next lane run.

- Tip lane (single-use CI rotation): run `37165203650`, ref
  `spacedock-ensign/pi-delegated-gate-continuation-reliability` at `d70bd32d6`,
  `-f live_cadence=pi`. Offline job green; `pi-live` ran 17 tests, 14 passed,
  3 failed. The `TestLiveCommonRecordedGateLifecycle` red was the removed
  judge's `pi-gate-binding-missing` / `pi-gate-canonical-authority-missing` /
  `pi-gate-review-missing` / `pi-gate-successor-model-mismatch`.
- Local focused batch on the captain's own Pi credentials (no CI rotation
  consumed): the entity's `-count=3` form is contaminated by one shared artifact
  dir (attempt 1 valid; attempts 2-3 invalid on `Pi root sessions=N, want exactly
  one`); three separate invocations with distinct dirs were the valid attempts,
  all rejected by the removed judge. Kept as history only.

### Deferred reds and bindings (rule: defer everything that is not a product defect)

- `auto-continue-after-implementation` keeps its binding to
  `mk72bnt1b5hsp9sfv83979xs` (`repair-pi-worker-lifecycle-observation`). Record:
  **evidence defect** — the assert recognizes only `subagent` status and
  `subagent_wait`, while this host supplies `bg_wait` and a native completion
  notice.
- `default-headless-gate-stop` keeps its binding to
  `gcmfwfjd9735b58sbzw7xsb8` (`repair-pi-recorded-gate-lifecycle`), with a
  **corrected record**: the reference was not missing. The officer mis-transcribed
  the path — it passed `/tmp/TestLiveCommonDefaultGateStop732520555/003/README.md`
  while the workflow root is `/tmp/TestLiveCommonDefaultHeadlessGateStop732520555/003`
  (same random suffix, one word dropped while reproducing a path the officer could
  see). The binary failed closed and the officer stopped per its own no-retry
  rule. Recorded as a **fragility red**; the deferred fix is to compose the
  reference from the observed workflow root instead of reproducing it. Stays
  deferred: no new task, and gc is not pulled in.
- `recorded-gate-lifecycle` stays unbound for now. With the judge gone the shared
  assertion runs; whether it passes is unknown until the next lane run.

### Summary

Cycle 6 removes the divergent Pi-only recorded-gate judge that cycle 3 added:
9w no longer restores a Pi grade, it removes one, and its presentation-reliability
acceptance stays OPEN. The `prepareRecordedGate` hook is back to a no-op so the
shared assertion runs, and the shared assertion is unchanged. The two
non-product-defect reds are deferred with bindings and corrected records; no new
task was created and gc was not pulled in.
