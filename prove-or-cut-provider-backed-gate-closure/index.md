---
title: Prove or cut provider-backed gate closure before stable v1
status: ideation
source: "Pre-0.27 gate-machinery necessity audit, 2026-08-01: provider evidence exists in pilot state, but the stable surface lacks one pinned exact-candidate end-to-end proof while the chat journey is the primary value path."
started: 2026-08-03T23:54:07Z
completed:
verdict:
score: "0.85"
worktree:
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
