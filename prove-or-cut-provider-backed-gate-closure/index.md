---
title: Prove or cut provider-backed gate closure before stable v1
status: validation
source: "Pre-0.27 gate-machinery necessity audit, 2026-08-01: provider evidence exists in pilot state, but the stable surface lacks one pinned exact-candidate end-to-end proof while the chat journey is the primary value path."
started: 2026-08-03T23:54:07Z
completed:
verdict:
score: "0.85"
worktree: .worktrees/spacedock-ensign-prove-or-cut-provider-backed-gate-closure
issue:
pr:
sprint: durable-decisions
id: a732sahay8wzgqrd2yr0xxr7
gates:
    version: 1
    records:
        - id: gate:a732sahay8wzgqrd2yr0xxr7:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:a732sahay8wzgqrd2yr0xxr7-backlog-1
              briefing:
                id: briefing:a732sahay8wzgqrd2yr0xxr7:backlog:attempt-1:revision-1
                digest: sha256:018e08a55eb205b626903dd53e5d99b610662dc63e1f22adaff859f2b5381970
                request-digest: sha256:dbd10def580302ef507bb0f7e07338bad509c91e76a3910101935c3f81834dc5
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:a732sahay8wzgqrd2yr0xxr7:backlog:1
                briefing: briefing:a732sahay8wzgqrd2yr0xxr7:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-03T23:53:47.120708Z"
                decision: approve
                reason: Science Officer concurs with the parallel BV/A7 ideation branch. Captain directed the next ideation dispatch. This gate only authorizes shaping the exact-candidate provider proof or cut; it does not authorize provider implementation or fallback machinery.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:a732sahay8wzgqrd2yr0xxr7:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:a732sahay8wzgqrd2yr0xxr7-ideation-1
              briefing:
                id: briefing:a732sahay8wzgqrd2yr0xxr7:ideation:attempt-1:revision-1
                digest: sha256:5400cfb247a4df7e3328e64bce41f9ca925cd3ef34fbecc19d9d9dd9b1e831f0
                request-digest: sha256:071f26a65624aa7817646d1fbad7e0d435e76876f6df4704c7af845671749a37
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:a732sahay8wzgqrd2yr0xxr7:ideation:1
                briefing: briefing:a732sahay8wzgqrd2yr0xxr7:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-04T00:13:19.389451Z"
                decision: hold
                reason: 'Science Officer advises hold: the design is coherent and q0 passed, but the retained-delivery fixture ended at noninteractive EOF and is not evidence. Obtain an interactive retained proof or the permitted chat-only clean-cut evidence with required full/race/live proof, then re-present; do not dispatch implementation or add fallback machinery.'
            - id: gate-attempt:a732sahay8wzgqrd2yr0xxr7-ideation-2
              briefing:
                id: briefing:a732sahay8wzgqrd2yr0xxr7:ideation:attempt-2:revision-1
                digest: sha256:a68f79a61ebc60a3e12631e2d68c4780c492362009b129a35f2c2fdedb56ba56
                request-digest: sha256:7f9fbea50c61cdd31f64e8a5b4f5007891590723f7606a977e67cb5cc02fa6ec
                room-ref: ./review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:a732sahay8wzgqrd2yr0xxr7:ideation:2
                briefing: briefing:a732sahay8wzgqrd2yr0xxr7:ideation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-07T01:19:00.836649Z"
                decision: approve
                reason: The clean-cut branch delivers one stable closure path. Three ideation items are done, the provider path lacks exact-candidate proof, and the design preserves the chat transaction without new machinery.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:a732sahay8wzgqrd2yr0xxr7:validation
          stage: validation
          attempts:
            - id: gate-attempt:a732sahay8wzgqrd2yr0xxr7-validation-1
              briefing:
                id: briefing:a732sahay8wzgqrd2yr0xxr7:validation:attempt-1:revision-1
                digest: sha256:6f17ce9e13d86077f5c1b0bc3924abe855adac0d453320ce4539116c2f72581e
                request-digest: sha256:89844d44acfb1f60f457715cc97916884770b242e360ae81109d0a719c35705d
                room-ref: ./review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:a732sahay8wzgqrd2yr0xxr7:validation:1
                briefing: briefing:a732sahay8wzgqrd2yr0xxr7:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-07T03:04:23.824737Z"
                decision: revise
                reason: Captain clarified that A7 owns one semantic recorder, not one presentation interface. Preserve Subspace presentation and semantic decision return through the existing --decision recorder. Remove only --room Result and inventory ingestion.
            - id: gate-attempt:a732sahay8wzgqrd2yr0xxr7-validation-2
              briefing:
                id: briefing:a732sahay8wzgqrd2yr0xxr7:validation:attempt-2:revision-1
                digest: sha256:7ed4480dfb984b8b79f1131e72ecb19c90348ec5e6de5e3854be7a606cc00555
                request-digest: sha256:3d0bf891b8233ef222cacbe31099d54d2208ae4b378f6f1649ce527fcb50ec63
                room-ref: ./review/validation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:a732sahay8wzgqrd2yr0xxr7:validation:2
                briefing: briefing:a732sahay8wzgqrd2yr0xxr7:validation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-07T03:42:13.553459Z"
                decision: revise
                reason: A7-V2 is material and owned by A7. Remove the obsolete room Result and inventory proof clauses from the canonical stable specification. Keep the one-recorder, multiple-presentation-interface boundary and change no CLI or recorder code.
            - id: gate-attempt:a732sahay8wzgqrd2yr0xxr7-validation-3
              briefing:
                id: briefing:a732sahay8wzgqrd2yr0xxr7:validation:attempt-3:revision-1
                digest: sha256:b443692bd31448b369bfd0eb9e4246d5cf03d977eba157742bd9367a3de4505a
                request-digest: sha256:d1b84c51278f95eaf91a17da4fef7e597a02b51b97aa886f1311e58ae8165cc0
                room-ref: ./review/validation/briefing-3
              resolution:
                type: Resolution
                id: resolution:spacedock:a732sahay8wzgqrd2yr0xxr7:validation:3
                briefing: briefing:a732sahay8wzgqrd2yr0xxr7:validation:attempt-3:revision-1
                by: agent:first-officer
                at: "2026-08-07T04:08:40.464434Z"
                decision: approve
                reason: Final validation at a2286c249 passes A7 acceptance criteria. Chat and Subspace presentation converge on one semantic recorder, provider-room ingestion is absent, and detached mutants prove the protected boundary.
              application:
                target-stage: done
                state: pending
---

Provider-backed closure is conditional v1 scope. Keep it only if one exact Spacedock
release candidate completes the same gate transaction as chat through one pinned
Subspace package. Otherwise cut the provider-only closure surface before the stable
tag; chat remains the only presentation path. There is no fallback, multiplexer, or
compatibility layer in either outcome.

## Problem and value

Chat already delivers the value transaction without a second runtime: prepare a
review, show it, record one decision, and consume it once. A provider adds a second
trust boundary (room request, capability probes, Git-root materialization, retained
Result/inventory, and recorder handoff). The repository has fixture and pilot pieces,
but no one exact-candidate proof from preparation through consumption. Shipping that
surface on archived fragments would make an unproven provider path part of v1.

The value to measure is not "a provider binary started". It is one captain decision
that closes the same prepared attempt, survives validation, and is consumed exactly
once without the agent authoring authority or output coordinates.

## Candidate pin and proof boundary

The ideation candidate is immutable and explicit:

```text
SPACEDOCK_CANDIDATE=a97a2b4cf9a1b708c6750d9758f3f5a67ef34d58
SUBSPACE_CANDIDATE=e6006dcaf0b64daebd29e9d4e4d226ce5e6716fa
```

`SPACEDOCK_CANDIDATE` is the full `main` tip at this dispatch. `SUBSPACE_CANDIDATE`
is the q0 `subspace:r gate <room>` branch tip; it includes the gate-room grammar,
provider-retention lifecycle, and canonical five-key binding Result. The pair is a
test input, not a moving branch. The implementation/validation worker must build
both from these exact commits, record each binary's SHA-256 and the provider skill
checkout SHA, and stop if either commit is unavailable. Archived runs at
Spacedock `59dc90a21cabb923dae4d03f2e3682f2c487da11` and earlier q0/s4 tips are
context only; they do not satisfy this candidate lock. If a prerequisite merge
changes the Spacedock tip, this entity must receive a superseding candidate pin
before the proof starts.

The prove-or-cut boundary is after the provider returns its retained Result and
before the first `gate record --room` mutation. Retention is proven only when one
run on the pinned pair observes every step below and ends with a consumed
application. Any refusal, provider launch failure, malformed/flattened Result,
authority disagreement, evidence drift, chat fallback, or second consume is a cut
signal—not permission to add a retry or translation.

## Common transaction trace

The default chat and retained-provider journeys share preparation and consumption:

1. The First Officer commits selected Artifact/References, runs `gate prepare`, and
   accepts only the emitted `room=`, `briefing=`, `digest=`, and `state=open` lines.
   Preparation writes `request.json` plus the canonical Briefing and no provider
   files. `state commit <slug>` durably binds that exact room before presentation.
2. Chat renders the bound Briefing once, then records with
   `gate record <slug> --decision approve|revise|hold --actor ...`; the close commit
   precedes any consume. No provider is installed or consulted.
3. If the retained branch is exercised, the agent passes only the opaque room to
   `/subspace:r gate <room>` after the bind commit. The agent does not read
   `request.json`, name a Result/log/inventory path, choose a provider executable,
   or repeat actor/approver/workflow coordinates.
4. The provider validates the frozen request and canonical Briefing, probes the
   literal capabilities `spacedock-gate-room-v1`, `review-v1-provider-package-v1`,
   and `review-v1-resolution-mode-v1`, materializes Git-root sources ephemerally,
   and opens one real decision surface. It owns `provider/result.json`,
   `provider/presented-inventory.json`, `provider/review.jsonl`, and diagnostics.
   The returned binding Result has exactly `type`, `briefing`, `artifact`,
   `annotations`, and `resolution`; authority appears once as `resolution.by`.
5. Spacedock runs `gate record <slug> --room <room>`. It re-reads and validates the
   request, arbitrary canonical Briefing locator, Result, and complete inventory;
   recomputes request/Briefing/Result/inventory pins; and writes only the frozen
   provider-evidence digests plus the Resolution. It derives the association in
   memory and writes no `association.json`. Any missing, advisory, duplicate,
   malformed, or drifted input leaves entity and room bytes unchanged.
6. Commit the closed attempt and provider evidence. Run `gate validate` and require
   a clean result, then `gate eligibility` and require the same approved pending
   application that chat produces. Run `gate consume` once and commit the successor
   status/application; terminal targets still obey the delivery/merge boundary.
   A second consume must refuse and preserve the entity byte-for-byte.

Provider unavailable or incapable is not a third journey. The agent stops with the
provider diagnostic; it does not launch a second provider, retry, or convert a
failed provider package into a chat approval. The ordinary chat path is selected
before presentation when no provider is intentionally retained.

## Retained design (only if the proof passes)

The smallest retained provider surface is one room-only handoff over the existing
recorder. The request has one authority value: `approver` is the sole value used to
authorize `resolution.by`; `actor` is removed from the provider-facing request and
public skill contract rather than kept as a second equal field. The recorder's chat
actor (`person:captain` or `agent:first-officer` with a reason) remains a separate
chat recording input and is not reconstructed by the provider. A duplicate or
conflicting authority member is rejected before JSON decode and before mutation.

No provider-specific dependency enters the Spacedock binary. The provider owns
capability discovery, source materialization, terminal selection, retained output,
and failure diagnostics; Spacedock owns room preparation, recorder validation,
provider evidence pins, and consumption. No exact-version check is added: the three
literal capability probes are the provider compatibility authority.

### Expected surface — retained branch

The estimate is 8 tracked Spacedock files, about +180/-120 lines, tolerance ±40%:

- `internal/gates/prepare.go`, `internal/gates/operation.go`, and their focused tests:
  emit/read one `approver` authority value, reject duplicate/conflicting authority,
  and keep the existing four digest checks (+70/-45).
- `internal/gates/testdata/gate-room/request.json` and provider fixtures: update the
  one authority field and canonical five-key Result (+20/-10).
- `skills/present-gate/SKILL.md` and `skills/fo-gate-lifecycle/SKILL.md`: state the
  opaque-room-only provider handoff and no caller output paths (+25/-35).
- `docs/specs/gate-resolution-frontmatter-contract.md` and
  `docs/site/concepts/gates-and-decisions.md`: document the retained boundary and
  exact authority/Result shape (+65/-30).

The q0 provider repository is an input pinned above, not a second Spacedock writer;
this task does not edit that repository. Any provider-side change needed to make the
candidate pass is outside this small boundary and therefore forces the cut branch.

### Retained documentation diff

In `skills/present-gate/SKILL.md`, replace the current override rule:

```diff
- A workflow or session may declare one presentation override. Apply this contract:
- 1. Pass only the emitted room ...
- 2. Reconstruct no authority ...
- 3. Record only prepared-room authority ...
+ A retained provider is one room-only presentation channel. After the bind commit,
+ pass exactly the emitted room to its declared `gate <room>` interface. The agent
+ supplies no request, authority, Result, inventory, log, or output path. The provider
+ returns its retained Result/inventory; Spacedock validates and records them with
+ `gate record --room`. A failed handoff leaves the attempt open; do not retry,
+ select another provider, or translate it into chat.
```

The contract and concepts pages receive the same observable statement: chat is the
default; a retained provider receives only the committed room; recorder and consume
ownership stay in Spacedock. No provider executable, package path, or multiplexer
name is added to those docs.

## Clean-cut design (if the proof fails)

Cut the second closure path, not the chat transaction. Remove the room-backed
provider recording option from the public CLI/help and the provider-only skill
wiring; keep `gate prepare`, chat `gate record --decision`, `gate validate`, and
`gate consume` unchanged. Remove provider-room positive fixtures and leave only a
negative command test proving `--room` is refused before entity/room mutation. Do
not leave a documented-but-unreachable `/subspace:r gate <room>` promise.

### Expected surface — clean-cut branch

The cut is a separate 10-file, approximately -260/+35-line change, tolerance ±40%:

- `internal/cli/cli.go` and gate help tests: remove `gate record --room` grammar and
  reject it as an unknown provider-only option before locking (+8/-18).
- `internal/gates/operation.go`, room-specific tests, and provider fixture files:
  delete the provider close arm and provider-evidence positive fixtures while leaving
  chat recorder validation intact (+15/-190).
- `skills/present-gate/SKILL.md`, `docs/specs/gate-resolution-frontmatter-contract.md`,
  and `docs/site/concepts/gates-and-decisions.md`: state chat-only v1 and move the
  provider room handoff under **Explicitly outside v1** (+12/-45).
- One skill/CLI smoke fixture proves no `gate <room>`, `gate record --room`, provider
  package path, or fallback wording is public, while the chat command sequence still
  passes (+0/+35).

No authority field is changed on the cut branch: AC-3 is not applicable because no
provider request is shipped. The room generated by `gate prepare` remains useful for
chat's bound Briefing and ordinary recorder validation; only provider-backed closure
is absent.

### Clean-cut documentation diff

Replace the public override paragraph in `skills/present-gate/SKILL.md` with:

```diff
- A workflow or session may declare one presentation override ... `gate record --room`.
+ Stable v1 presents gates in chat only. The First Officer binds the room, renders the
+ canonical Briefing once, and records the semantic decision with `gate record
+ --decision`. Provider presentation and room-backed Result recording are explicitly
+ outside v1; a failed or partial external experiment never becomes a chat approval.
```

In the contract/spec and concepts pages, remove the override branch from the lifecycle
diagram and add one sentence under **Explicitly outside v1**: "Provider-backed
presentation, `gate record --room`, retained provider evidence, and provider package
selection are not stable-v1 surfaces; the chat transaction remains supported." The
roadmap's provider row records the cut and names no external fallback.

## Acceptance criteria

**AC-1 (VALUE) — The exact candidate has one closure outcome: a complete retained provider transaction or no provider closure surface.**

Verified either by one live run at the pinned Spacedock/Subspace commits that records
the exact candidate tuple, capability exits, one room-only invocation, Result and
inventory digests, `gate record --room`, `gate validate`, one successful `gate consume`,
and a byte-clean second consume; or by the clean-cut CLI/skill smoke plus the ordinary
chat journey. A provider run that stops before recording, uses archived commits, or
needs a fallback does not pass.

**AC-2 — A retained provider is room-only and preserves the chat value transaction.**

Verified by the real invocation argv and retained diagnostics: the agent supplies only
the emitted room; no `request.json` parsing, actor/approver reconstruction, output
paths, provider executable, or multiplexer arguments appear. The same candidate's
chat journey remains green with prepare, presentation, decision record, validation,
and one-use consume.

**AC-3 — A retained request has one unambiguous authority value.**

Verified by request/recorder tests that use only `approver` and `resolution.by`, plus
recursive duplicate and conflicting-authority inputs that fail closed with unchanged
entity and room bytes. This criterion is explicitly not applicable on the clean-cut
branch, which ships no provider request.

**AC-4 — The candidate does not gain fallback, multiplexer, or compatibility machinery.**

Verified by diff inspection and negative tests: no second provider selector, exact
version comparison, provider dependency in `go list -deps ./cmd/spacedock`, alternate
Result envelope, or chat fallback exists. Removing any one of these negative checks
must make the test red.

## Test plan and spike record

The existing provider-room fixture remains the offline control. It must exercise a
complete Result/inventory, truncated inventory, wrong type/id/revision, duplicate
authority, and retained-byte drift. The q0 candidate's own contract fixture passed
at `e6006dcaf0b64daebd29e9d4e4d226ce5e6716fa`; its full retained-delivery fixture was
started from a temporary archive but stopped in the noninteractive EOF child path,
so it is not live/provider evidence. That is the riskiest unresolved mechanism.

Implementation must run the real provider package on the exact pair in a disposable
split-root workflow. Use one pinned terminal/provider package, not a new multiplexer
harness. Capture: candidate SHAs and binary hashes; capability stdout/stderr/exit
files; emitted room and state-commit SHAs; provider `result.json`, inventory, log, and
diagnostics; recorder/validate/eligibility/consume exits; and before/after entity
hashes. A complete run is the retained AC-1 evidence. A refusal must be classified
as a product defect only if it is inside this task's small Spacedock surface; otherwise
execute the clean cut.

Run the required repository checks for whichever branch lands:
`go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`. For the
retained branch also run the provider repository's full/race/contract suite and the
real headed room journey; a green fake fixture alone is insufficient. For the cut
branch, the chat live lane and the negative `--room`/skill smoke must pass. No archived
provider evidence is promoted to candidate proof.

## Stage Report: ideation

- DONE: Pin the exact candidate and trace the room-only provider transaction from preparation through presentation, record, validation, and consume; identify the prove-or-cut boundary.
  Locked `SPACEDOCK_CANDIDATE=a97a2b4cf9a1b708c6750d9758f3f5a67ef34d58` and `SUBSPACE_CANDIDATE=e6006dcaf0b64daebd29e9d4e4d226ce5e6716fa`; the numbered trace names every handoff and the boundary before `gate record --room`. Archived q0/s4 evidence is explicitly non-proof.
- DONE: Define the minimal retained-provider or clean-cut design, including one authority value, expected surface, documentation changes, and no fallback or multiplexer machinery.
  Retention uses one `approver`/`resolution.by` authority and an 8-file bounded Spacedock surface; failure uses a separate 10-file chat-only cut with explicit CLI/skill/spec docs and no provider fallback.
- DONE: Specify falsifiable live/provider and cut-path evidence using the existing fixture, with exact candidate, refusal, chat-preservation, full/race, formatting, and required-lane requirements.
  AC-1, AC-2, AC-3, and AC-4 name candidate/binary hashes, capability/refusal and room-only argv evidence, chat green controls, provider full/race/contract plus Spacedock full/race/format checks, and clean-cut negative smoke requirements.

### Summary

The provider path now has a single falsifiable release test: exact Spacedock and
Subspace commits, one opaque room handoff, one canonical Result authority, recorder
validation, and one-use consume. The q0 contract fixture is green but its real
retained-delivery run is not evidence because the noninteractive child stalled. The
implementation therefore has one honest decision: complete the live proof on the
pinned pair, or land the separately scoped chat-only cut. No fallback, multiplexer,
exact-version, or provider dependency can be smuggled into stable v1.

## Stage Report: implementation

- DONE: Ship one stable chat closure path and remove the unproved provider-only closure surface.
  Commits `c3c864e0b` and `4ff999250` remove the public `--room` grammar, provider close arm, positive provider fixtures, lifecycle branches, and provider instructions; `TestGateRecordRejectsProviderRoomBeforeMutation` fails if the option is accepted or reaches the entity lock.
- DONE: Preserve prepare, presentation, decision record, validation, and one-use consume with focused and repository tests.
  `TestRecordedGateLifecycleRealCLIReplay` exercises the real CLI through prepare, chat presentation trace, decision record, validation, eligibility, one consume, and byte-clean repeat refusal; `TestStableV1GateSkillsExposeOnlyChatClosure` fails if chat commands disappear or provider-only skill wiring returns; formatting and focused checks completed cleanly.
- DONE: Keep the change self-contained: add no parser, fallback, multiplexer, compatibility path, or provider dependency.
  `go list -deps ./cmd/spacedock` contains no Subspace/provider dependency, the diff adds no dependency or alternate envelope, and provider-positive fixtures were deleted rather than translated.

### Summary

Stable v1 now has one gate-closure journey: prepare and bind the canonical Briefing,
present it in chat, record the semantic decision, validate eligibility, and consume
approval once. The First Officer accepted the 16-file `+216/-1280` footprint for
validation because the estimate delta consists of deleting the provider close arm and
its positive tests/fixtures, with no replacement machinery. `gofmt -w ./cmd ./internal`,
`git diff --check`, focused tests, and the real chat replay passed. Both `go test ./...`
and `go test ./... -race` ran all packages; every package except `internal/gates`
passed, where only `TestV1PilotManifestReadsAndValidates` failed because
`codex-launch-multi-agent-v2.md` and `gate-agent-ergonomics.md` now live under
`_archive/` rather than the active paths pinned by the exact BV candidate test.

## Review-finding disposition

### A7-V1 — Candidate cuts the permitted Subspace presentation channel

- Observation: candidate `4ff999250` says stable v1 presentation is chat-only and
  places provider-backed presentation outside v1 in `present-gate`, the gate spec,
  concepts documentation, and roadmap; `TestStableV1GateSkillsExposeOnlyChatClosure`
  requires that over-cut wording.
- Evidence: the Captain clarified that the stable boundary is one semantic recorder:
  Subspace may present the committed gate and return semantic decision/reason input,
  after which the FO uses existing `gate record --decision`; only `gate record --room`
  and Result/inventory ingestion must be absent.
- Defect kind: outcome defect plus evidence defect, affecting AC-1/AC-2 at the exact
  presentation-to-recorder boundary. The shipped contract prohibits intended value,
  and its contract test rejects corrected channel-neutral wording.
- Release scope: material. The trigger is the supported Subspace presentation path,
  and the candidate explicitly declares it unsupported across shipped surfaces.
- Proposed ownership/disposition: current task, narrow correction. Keep all CLI and
  recorder deletions; revise only affected skill/docs wording and the contract test to
  permit Subspace presentation feeding semantic decision/reason into `--decision`,
  without adding a parser, authority reconstruction, Result envelope, or room close.

### Feedback Cycles

- Cycle 1: REJECTED — validation feedback; surface 16 files/1498 LOC vs estimate 10 files/295 LOC (508%); AC unchanged
- Cycle 2: REJECTED — validation feedback; surface 16 files/1521 LOC vs estimate 10 files/295 LOC (516%); AC unchanged

## Stage Report: validation

- DONE: Prove that the exact candidate exposes one complete chat closure and rejects provider-room recording before mutation.
  At exact `4ff999250`, focused CLI/real-binary tests passed; `--room` exits 2 with exact unknown-flag output before lock creation and preserves entity/room bytes, while chat prepare/record/consume is exercised. A detached accepted-`--room` mutant and no-op-consume mutant made their relevant tests fail.
- DONE: Prove that the removal adds no parser, fallback, multiplexer, compatibility path, alternate authority, or provider dependency.
  Diff inspection from `698867babe7d57eb309dca476ae91187e92a3a57`, `go list -deps ./cmd/spacedock`, and module inspection found deletion rather than replacement machinery and no Subspace/provider dependency.
- DONE: Run focused, full, race, and detached adversarial checks; classify the pilot-manifest failures against current durable state.
  Focused checks passed; full and race runs failed only `TestV1PilotManifestReadsAndValidates` because two unchanged manifest paths were renamed 99%-similar into `_archive/` by durable-state commits `e3aa5de` and `3f85cb0`, so they are stale cross-state evidence, not A7 regressions. `gofmt` and `git diff --check` were clean.
- FAILED: Recommend PASSED with no material findings.
  Recommend REJECTED: A7-V1 materially over-cuts the supported presentation channel, and the detached channel-neutral skill mutant proves the new contract test pins that obsolete target.

### Summary

The recorder cut itself is sound: stable v1 rejects provider-room ingestion before
mutation, preserves the semantic chat closure, and adds no replacement machinery.
The candidate is nevertheless not releasable because its skill, specification,
concepts, roadmap, and test prohibit Subspace presentation rather than only removing
provider-specific recording. Apply the narrow wording/test correction above and rerun
focused, full, race, and detached checks; do not restore `--room` or provider evidence.

## Stage Report: implementation (cycle 2)

- DONE: Preserve Subspace as a permitted gate-presentation channel that returns semantic decision and reason input.
  Commits `bff542ae7` and `a2ed9d654` require and document chat or Subspace presentation converging on semantic decision/reason input; a detached chat-only skill mutant makes `TestStableV1GateSkillsUseOneRecorderAcrossPresentationChannels` fail.
- DONE: Keep one standard semantic recorder and keep provider-specific room, Result, and inventory ingestion removed.
  The focused `TestGateRecordRejectsProviderRoomBeforeMutation` passes, `go list -deps ./cmd/spacedock` names no Subspace/provider dependency, and the correction changes no CLI, recorder, Result, inventory, authority, or dependency code.
- DONE: Replace the chat-only contract assertion with a channel-neutral test and rerun the required evidence.
  The channel-neutral contract test and real chat transaction pass; `gofmt -w ./cmd ./internal` and `git diff --check` are clean. Full/race runs still fail the two stale active paths in `TestV1PilotManifestReadsAndValidates`; race also hit `TestCodexProcessActivityResetsQuietBudget` once, whose immediate focused race rerun passed.

### Summary

The correction preserves Subspace as a presentation interface while keeping
`gate record --decision` as the sole recorder and retaining every provider-ingestion
deletion. Required focused, live, full, race, and detached evidence was rerun; the
remaining exact failures are the previously classified archived pilot paths and one
non-reproducing quiet-budget race timing failure.

### A7-V2 — Canonical spec still requires the removed room-Result mechanism

- Observation: corrected candidate `a2ed9d654` declares room-backed Result/inventory
  ingestion outside v1 at spec lines 260-266, but its current Behavioral proof item 7
  still requires release tests to verify a room Result, request authority, inventory,
  and presentation mapping before close; line 107 also promises Result validation.
- Released user and normal workflow: stable-v1 implementers and release validators use
  `docs/specs/gate-resolution-frontmatter-contract.md` as the canonical shipped contract.
- Observable harm: the contract simultaneously cuts and requires the same recorder
  mechanism, leaving release proof impossible without restoring forbidden ingestion.
- Authority: value-ac[AC-1] requires either a retained complete provider transaction or
  no provider closure surface; the clean cut cannot retain a room-Result proof promise.
- Trigger evidence: exact `a2ed9d654`, spec lines 104-107 and 260-283; the CLI refusal
  test proves `--room` is absent, so this is a live contradiction rather than history.
- Defect kind: outcome and evidence defect at AC-1's shipped-contract/proof boundary.
- Release scope: material; the contradiction is on the normal stable release checklist.
- Proposed ownership/disposition: current task, narrow correction. Remove the obsolete
  Result-validation clause and Behavioral proof item 7; reconcile remaining current
  command/roadmap references to semantic decisions without restoring any machinery.

## Stage Report: validation (cycle 2)

- DONE: Prove that the exact candidate exposes one complete chat closure and rejects provider-room recording before mutation.
  At corrected head `a2ed9d654`, the focused CLI and real-binary lifecycle tests pass; `--room` remains an exact exit-2 refusal before lock creation with entity and room bytes unchanged.
- DONE: Prove that the removal adds no parser, fallback, multiplexer, compatibility path, alternate authority, or provider dependency.
  The five-file correction from `4ff999250` changes only skill/docs/contract-test prose; full diff inspection and dependency listings show no CLI, recorder, authority, envelope, or provider dependency addition.
- DONE: Run focused, full, race, and detached adversarial checks; classify the pilot-manifest failures against current durable state.
  Focused tests passed and the detached chat-only presentation mutant made `TestStableV1GateSkillsUseOneRecorderAcrossPresentationChannels` fail. Full and race runs failed only the unchanged two manifest paths, which current state and rename commits prove now live under `_archive/`; the prior quiet-budget race failure did not recur.
- FAILED: Recommend PASSED with no material findings.
  Recommend REJECTED: A7-V1 is fixed, but A7-V2 leaves the canonical stable contract requiring the provider-specific room Result/inventory proof that this candidate intentionally removed.

### Summary

The authorized cycle-1 correction works: chat and Subspace presentation now converge
on one `gate record --decision` recorder, the chat transaction remains green, and the
detached regression control is falsifiable. The candidate still cannot pass because
its canonical Behavioral proof section promises the removed room-backed mechanism.
Delete that stale normative proof and reconcile the few remaining current references;
do not change the already-correct CLI or recorder cut.

## Stage Report: implementation (cycle 3)

- DONE: Remove the obsolete room Result and inventory proof requirement from the canonical stable specification.
  Commits `f3e7d9953` and `a2286c249` remove the Result-validation clause and Behavioral proof item 7, replace the stale presentation-mapping proof, and make `TestStableV1CanonicalContractDoesNotRequireRoomResultProof` fail if those requirements return.
- DONE: Reconcile remaining normative references with one semantic recorder across chat and Subspace presentation.
  The roadmap now removes Result association, exact-Result persistence, provider mapping, and pinned provider-closure requirements while preserving chat or Subspace presentation converging on `gate record --decision`.
- DONE: Run focused, full, race, and detached evidence without changing CLI or recorder code.
  Focused contract/skill/CLI tests and the real chat transaction pass; the detached chat-only mutant fails at channel convergence. Full and race runs fail only `TestV1PilotManifestReadsAndValidates` because `codex-launch-multi-agent-v2.md` and `gate-agent-ergonomics.md` now live under `_archive/`; the cycle diff contains only specification, roadmap, and contract-test files.

### Summary

The canonical stable contract no longer requires the provider-specific mechanism that
the candidate removed. Chat and Subspace remain presentation interfaces over one
semantic recorder, with the unchanged stale pilot-manifest paths as the only full/race
failure and no CLI, recorder, authority, skill, or provider-code change in this cycle.

### A7-V2 validation disposition

Corrected at `f3e7d9953` and `a2286c249`: the canonical spec no longer promises
Result validation or a room-Result release proof, and the current roadmap no longer
requires Result association, exact-Result persistence, provider mapping, or pinned
provider closure. No material A7 finding remains.

## Stage Report: validation (cycle 3)

- DONE: Prove that the exact candidate exposes one complete chat closure and rejects provider-room recording before mutation.
  At corrected head `a2286c249`, `TestRecordedGateLifecycleRealCLIReplay` passes the real prepare/chat-record/consume journey and `TestGateRecordRejectsProviderRoomBeforeMutation` proves exact exit 2 before lock creation with unchanged entity/room bytes. This satisfies AC-1's clean-cut branch; AC-2 preserves chat and Subspace presentation converging on semantic input, and AC-3 is not applicable because no provider request is shipped.
- DONE: Prove that the removal adds no parser, fallback, multiplexer, compatibility path, alternate authority, or provider dependency.
  Cycle 3 changes only the spec, roadmap, and contract test; the full BV-to-head diff adds no alternate recorder or envelope, while `go list -deps ./cmd/spacedock` and `go list -m all` contain no Subspace/provider dependency. This reproduces AC-4.
- DONE: Run focused, full, race, and detached adversarial checks; classify the pilot-manifest failures against current durable state.
  Focused tests, formatting, and diff checks pass. Detached chat-only and restored-Result-validation mutants make the two contract tests fail. Full and race runs fail only `TestV1PilotManifestReadsAndValidates`: A7 leaves its manifest byte-unchanged, while state commits `3f85cb081` and `e3aa5b753` prove the two missing active paths were moved to `_archive/` by merge guard.
- DONE: Recommend PASSED with no material A7 findings.
  A7-V1 and A7-V2 are corrected and no new material or deferred A7 risk remains. The stale pilot manifest is an observed external release-check failure and must be refreshed before the stable tag; it is not evidence against this candidate's ACs.

### Summary

Recommend PASSED for A7 at exact corrected candidate `a2286c249`. Stable v1 now has
one semantic recorder, permits chat or Subspace presentation, rejects `--room` before
mutation, and carries no Result/inventory ingestion or replacement machinery. The
repository-wide full and race commands are not green solely because an unchanged
pilot fixture points at two entities that current durable state has archived; that
separate release prerequisite remains visible rather than being treated as a pass.
