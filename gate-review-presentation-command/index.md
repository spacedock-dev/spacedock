---
title: Gate presentation as an overridable channel with atomic result retention
status: validation
source: "Split from the gate-recorder task (3k), captain-approved 2026-07-21. The subspace-coupled presentation half; 3k cycles 11-12 are its banked design history."
id: xbatj4hxtxw9t83vvmfem27f
sprint: durable-decisions
group: recorder
started: 2026-07-21T01:43:36Z
worktree: .worktrees/spacedock-ensign-gate-review-presentation-command
gates:
    version: 1
    current:
        gate: gate:docs-dev:xb:ideation
    records:
        - id: gate:docs-dev:xb:validation
          stage: validation
          attempts:
            - id: gate-attempt:xb-validation-1
              briefing:
                id: briefing:docs-dev:xb:validation:attempt-1:revision-1
                digest: sha256:772a856dcd3dd7d5a1bcfb589854b4b7f5b70bb26393a7e1e90aa2605daf0911
                digest-domain: canonical-bytes
                room-ref: ./review/validation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:xb:validation:1
                briefing: briefing:docs-dev:xb:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-23T23:49:14.728783Z"
                decision: approve
                reason: Spacedock 612b72fc and provider 198f7623 have 6/6 ACs evidenced, retained-delivery and association suites green, zero binary coupling, and no material finding.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: advance
                target-stage: done
                state: pending
                blockers: []
        - id: gate:docs-dev:xb:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:xb-ideation-1
              briefing:
                id: briefing:docs-dev:xb:ideation:attempt-1:revision-1
                digest: sha256:a552c2b7978d9fb642beddba360b926bcf3c334072dba60d31744bba18cae552
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-4
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:xb:ideation:1
                briefing: briefing:docs-dev:xb:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-24T12:18:27.93418Z"
                decision: approve
                reason: Cycle 12 aligns all six ACs to the approved prepared-room boundary; independent staff re-review APPROVES; only +42/-0 test lines are authorized, projecting 1,297/1,300 changed LOC.
                adoption-note: Captain granted the First Officer the conn to approve sprint gates and required suspicion of drift; preserve exact candidate 98ebb458 and add only the three recorded test variants.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
review-round:
    id: round:xbatj4hxtxw9t83vvmfem27f:implementation:13
    stage: implementation
    cycle: 13
    briefing:
        id: briefing:gate-review-presentation-command:implementation:round-13
        digest: sha256:74c9675995b6a1deeb0cebc1925545b420719d99d4bf1c699b8dac7ea04c2fe5
        digest-domain: canonical-bytes
        room-ref: ./review/implementation/round-13
---

Gate presentation remains an overridable channel of the present-gate skill, not a
`spacedock` presentation verb. xb's product boundary is narrower: Spacedock consumes one
scaffold-prepared gate room through
`spacedock gate record <entity> --room <gate-room>`. The room fixes the request,
canonical Briefing, attempt identity, and provider output locations before recording.
The recorder derives all association data itself, accepts only direct Captain authority,
and closes the attempt atomically while retaining exact byte identities for later
validation.

Room construction and provider transport remain outside xb. The candidate has no
Subspace dependency, presentation launcher, provider lifecycle, or second association
input. This keeps both end values: a decision becomes durable only when it is exactly
the Captain-authorized decision presented for the bound attempt, and Spacedock's binary
does not couple its release train to a presentation provider.

## Expected surface + tolerance

The approved exact-tip candidate is `98ebb458`: 17 files,
`+1086/-169` = **1,255 changed lines**, below Cycle 9's hard cap of **1,300**.
Its surface is the provider-neutral recorder path, focused fixtures and behavioral tests,
the gate contract and user documentation, and the declarative present-gate seam. It
contains no provider transport or modeled provider harness.

Cycle 11 adds an exact estimated **`+42/-0` test-only lines** in the existing
`internal/cli/gate_test.go`: 12 lines for the post-bind Briefing-question mutation, 18
lines for the two compact presented-inventory matrix rows, and 12 lines for deleting
each retained file through the existing validation loop. The projected surface is
still 17 files and **1,297 changed lines**, three lines under the approved 1,300-line
cap. No cap reset, production edit, fixture file, provider harness, or documentation
change is proposed. Any larger implementation delta must first justify why the existing
CLI matrices and retained-file loop cannot carry the proof.

## Problem

The implementation and the binding design prose disagree. Validation cycle 4 found no
material defect at `98ebb458`, but the entity still required a retired provider-side
override script, pane/session and launcher behavior, presenter fallback, and
provider-minted-ID normalization. Those obligations test a transport architecture that
Cycles 4-9 explicitly removed.

Leaving them binding would either reject a correct recorder for work outside its
ownership or cause implementation to rebuild excluded provider machinery. The task
therefore needs a design-only reset around the supported boundary: Spacedock records
exactly one already-prepared room, validates every authority and presentation claim
before mutation, and preserves enough immutable identity to detect later byte changes.

## Prepared-room recorder contract

1. **One input boundary.** Provider recording accepts only
   `spacedock gate record <entity> --room <gate-room>`. The fixed room contains
   `request.json`, `briefing.json`, `provider/result.json`, and
   `provider/presented-inventory.json`. Callers cannot submit provider paths, argv,
   actor/approver flags, adoption text, or an association.
2. **Frozen request, Briefing, attempt, and authority.** `request.json` names the
   derived logical gate and current open attempt, the exact canonical Briefing id and
   JCS digest, and identical Captain actor/approver authority. Binding freezes the
   request digest beside the Briefing digest and room reference. Malformed, stale,
   moved, or rebound rooms fail before any entity write.
3. **Direct Captain Result only.** The fixed Result must name the canonical Briefing
   exactly and carry a wrapper-free binding `Resolution` whose `by` equals the frozen
   Captain authority. Advisory status/binding/identity wrappers, different authority,
   unknown authority fields, nested adoption provenance, and provider-minted Briefing
   ids cannot close the gate. Delegated First Officer approval remains the separate
   chat `--decision ... --directive` path.
4. **Complete derived presentation inventory.** The recorder derives the canonical
   inventory from every Briefing Artifact and every recursively reached Reference. The
   fixed presented inventory must map every derived id and revision exactly once, and
   the Result primary must be an Artifact. No caller-authored or stored duplicate
   association is accepted.
5. **Atomic one-use close.** The recorder verifies the whole room before replacing the
   entity. Any failed check leaves the entity byte-identical. A valid room closes only
   the current open attempt once; closed history, its application, and its provider
   evidence remain frozen, and stale replay cannot rebind or append through the old
   room.
6. **Exact retained-byte identity.** A provider close stores only the raw SHA-256
   digests of `provider/result.json` and
   `provider/presented-inventory.json`. `gate validate` recomputes both fixed room files
   and fails if either is missing or byte-changed. The association remains derived and
   unstored.
7. **Provider-neutral boundary.** The binary has no presentation verb and no Subspace
   dependency. The repository may state the declarative room handoff, but it does not
   construct rooms, launch or poll a provider, model pane/session lifecycle, or own
   provider retention behavior.

## Proof basis: no new spike needed

No mechanism in the reset is unverified. Exact-tip behavior tests already exercise the
prepared-room success path, adjacent authority and identity mutations, recursive
inventory rejection, failure atomicity, frozen history, and retained-file byte
tampering. Cycle 11 identifies three narrower falsifiability gaps in that same supported
path: changed Briefing question bytes after binding, same-cardinality inventory
substitutions, and missing retained files. They extend the existing CLI mutation
matrices and retained-file loop; they do not require a spike, new harness, interface, or
production behavior. The complete repository and race suites, strict documentation
build, formatting and dependency checks passed at `98ebb458`. Final-tip Roborev job
**2028** reviewed `dd6bd114..98ebb458` and returned **P — No issues found**.

## Acceptance criteria

**AC-1 (VALUE) — Only the exact authorized room can make a Captain decision durable.**
The canonical fixture closes its bound attempt exactly once; every adjacent
mutation of request identity, Briefing identity/digest, room identity, authority,
Result authority, or presentation inventory is rejected, with **zero entity-byte
changes across the rejection matrix**. The independently mutated room files are the
baseline that can move the wrong way: any accepted mutation or any rejection-side
entity delta fails the value. *Test:* `TestGateRecordConsumesDirectBindingResultFromPreparedRoom`,
including a valid Briefing whose question bytes change only after its attempt binds,
`TestGateRoomRejectsRequestAuthorityRebindingWithoutMutation`,
`TestGateRoomRejectsRequestBackedRoomMoveWithoutMutation`,
`TestGateRoomRejectsBriefingOnlyAttemptWithoutFrozenRequest`,
`TestGateRoomRejectsAdvisoryEvidenceAndUnauthorizedResolutionWithoutMutation`, and
`TestGateRoomRejectsDistinctAuthorityAndDifferentRoomWithoutMutation`.

**AC-2 — Request, Briefing, attempt, and authority identity are exact and frozen before recording.**
The request must bind the derived current-stage gate and open attempt, the
canonical Briefing id/digest, the bound room, and equal Captain actor/approver authority;
post-bind changes fail before mutation. *Test:* bind a valid Briefing, mutate only its
question, then require room recording to fail with a byte-identical entity; retain the
request-authority, room-move, briefing-only, malformed-Reference, and distinct-identity
tests named under AC-1, plus `TestGateBriefingRejectsMalformedReferenceBeforeBinding`.

**AC-3 — A provider close contains one direct Captain Result and a complete recursively derived presentation.**
Wrapper/advisory fields, adoption provenance, wrong authority,
provider-minted identity, incomplete or mistyped inventory, duplicate mapping, or a
Reference primary fail closed. The valid fixture records only the portable nested
Resolution. *Test:* `TestGateRecordConsumesDirectBindingResultFromPreparedRoom`,
`TestGateRoomRejectsAdvisoryEvidenceAndUnauthorizedResolutionWithoutMutation`, and
`TestExactCanonicalBriefingIsIndependentAssociationInventory`; the room matrix includes
same-cardinality variants that duplicate one presented id and that replace one correct
revision with a wrong revision, each rejected without entity mutation.

**AC-4 — Close is atomic and one-use, and closed history cannot be rebound.** Every room
validation completes before entity replacement; a valid room closes only the last open
current-stage attempt, repeat/stale use does not create another closure, and later
application remains a separate one-use transition. *Test:*
`TestCanonicalLifecycleRebindCloseFreezeAndSupersede`,
`TestWriterCASValidationAtomicityAndLock`,
`TestGateRecordConsumesDirectBindingResultFromPreparedRoom`, and
`TestGateEligibilityAndConsumeAuthorizeOnce`.

**AC-5 — The closed attempt identifies the exact retained Result and inventory bytes.**
Successful close freezes both raw digests; appending even one byte to either retained
file, or deleting either file, makes `gate validate` fail; restoring the exact bytes
makes it pass. No copied artifact or stored association is required. *Test:* extend
`TestGateRecordConsumesDirectBindingResultFromPreparedRoom`'s existing retained-file
loop so each Result/inventory path is first byte-mutated and then deleted. A present
but byte-changed file must name the fixed path and `frozen digest`; a missing file must
name the fixed path plus its read/not-found error. Restoring the exact bytes must pass
validation.

**AC-6 (VALUE) — Provider coupling in the Spacedock binary remains zero.** The CLI
exposes no `gate review` presentation verb and `go list -deps ./cmd/spacedock` contains
zero Subspace dependencies. The only provider form is room consumption; reintroducing
either a verb or dependency makes the measured count non-zero. *Test:*
`TestGateReviewVerbIsAbsentAndSideEffectFree`, command-surface fixtures, and the exact-tip
dependency check.

## New mechanisms (value AC served / simplest alternative / why insufficient)

- **One prepared-room input** — serves AC-1 and AC-2. *Alternative:* accept Result,
  association, authority, output paths, and provider arguments separately.
  *Insufficient:* caller-supplied duplicate facts can disagree with the attempt and make
  the recorder certify an association it did not derive.
- **Direct nested Captain Resolution** — serves AC-1 and AC-3. *Alternative:* accept
  advisory output plus an adoption note or normalize provider identity.
  *Insufficient:* it obscures who made the binding decision and permits evidence created
  for another Briefing to acquire local authority.
- **Recursive inventory derivation** — serves AC-1 and AC-3. *Alternative:* trust a
  caller-built association or only check the primary Artifact. *Insufficient:* either
  lets supporting References disappear or be mistyped while still claiming complete
  presentation.
- **Atomic close on the existing canonical writer** — serves AC-1 and AC-4.
  *Alternative:* add a provider transaction log or modeled provider harness.
  *Insufficient:* the canonical compare-and-swap writer already proves byte-clean
  rejection and frozen history; a second model duplicates the implementation without
  strengthening the supported boundary.
- **Two raw retained-file digests** — serves AC-5. *Alternative:* copy provider files
  into entity state, store a full evidence registry, or retain only a derived
  association. *Insufficient:* copies and registries expand ownership, while a derived
  association cannot identify which exact provider bytes were accepted.
- **No presentation verb or provider dependency** — serves AC-6. *Alternative:* own
  provider launch and lifecycle inside Spacedock. *Insufficient:* it couples release
  trains without improving room verification.

## Behavioral test plan

Behavior-first Go fixtures drive the real CLI against copied room files and compare exact
entity bytes. No live provider, transport smoke, provider script, or modeled provider
interpreter is needed. Estimated cost is low because the complete plan already passes at
the exact tip.

1. **Canonical close and mutation matrix (AC-1, AC-2, AC-3).** Bind the fixture
   Briefing, record its room, and assert the portable Captain Resolution. Independently
   mutate every request/Briefing/attempt/authority field, including the question of a
   valid Briefing after binding; mutate advisory wrappers, adoption fields, Result
   binding, inventory entry/type/revision/coverage, and room path. The inventory matrix
   includes a same-cardinality duplicate id and a same-cardinality wrong revision.
   Assert non-zero exit and byte-identical entity for every rejection.
2. **Lifecycle and atomicity (AC-4).** Exercise open, identical rebind, changed rebind,
   close, stale successor, frozen closure, compare-and-swap failure, and one-use
   application. Assert one closure and no mutation on every rejected replay.
3. **Exact-byte validation (AC-5).** Close a valid room, then use the existing
   retained-file loop to append one byte and separately delete each Result/inventory
   file. For a present-but-changed file, require `gate validate` to name the fixed path
   and `frozen digest`. For a missing file, require the fixed path plus the
   read/not-found error; do not require a digest diagnostic for bytes that cannot be
   read. Restore the original bytes and require success.
4. **Zero-coupling boundary (AC-6).** Require `gate review` to remain absent and
   side-effect-free, inspect the real command surface, and require zero Subspace entries
   from `go list -deps ./cmd/spacedock`.
5. **Repository confidence.** Run focused CLI/gates/contractlint tests, `go test ./...`,
   `go test ./... -race`, strict MkDocs, `gofmt -w ./cmd ./internal`,
   `git diff --check`, and the surface/cleanliness audit. Roborev 2028 is the existing
   detached exact-tip review.

## Documentation change proposal

The exact-tip candidate already carries the required user-visible change; this reset
does not propose additional product documentation. The approved before/after is:

```diff
-spacedock gate record ENTITY --result PATH --association PATH [provider identity flags]
+spacedock gate record ENTITY --room PATH [--workflow-dir DIR]
+
+The scaffold prepares one room binding request authority, gate attempt, canonical
+Briefing, and fixed provider outputs. The recorder derives the complete recursive
+Artifact/Reference association, accepts only a direct Captain Resolution, and stores
+the exact Result and presented-inventory digests for later validation.
```

`docs/specs/gate-resolution-frontmatter-contract.md` owns the normative room contract;
`docs/site/concepts/gates-and-decisions.md` explains it; the command reference exposes
only `--room`; and `skills/present-gate/SKILL.md` hands the already-prepared room to the
provider and then to the recorder.

## Out of scope

- Gate-room construction or a room-preparation command.
- Provider transport, override scripts, launch/poll/wait behavior, pane/session state,
  launcher-death handling, fallback selection, retained-file creation, or provider
  diagnostics.
- Provider-minted-ID normalization, advisory adoption notes, or any other mechanism
  that converts non-Captain evidence into a binding Result.
- A modeled provider interpreter/harness, caller-authored association, stored duplicate
  association, provider transaction log, evidence registry, or artifact copies.
- A `spacedock gate review` verb, Subspace import/shell-out, or Subspace product change.
- Blocker evaluation, execution holds, eligibility policy, or application effects beyond
  preserving the recorder's existing one-use boundary.
- Room archival, deletion, or reformatting; v1 validates the retained fixed files in
  place.

## Stage Report: ideation

- DONE: A concrete command design against the banked 3k cycles 11-12 lifecycle: blocking presenter, atomic result/log/diagnostics retention on success AND failure (the destroyed hold-path result and the blank-float EOF are the red fixtures), and the provider id-mapping implemented as SPECIFIED in 3k's gate-resolution-frontmatter-contract.md.
  `## Required capability` designs the six-step lifecycle (validate → announce → blocking child → owned retention dir → completion-is-exit+validated+retained → no frontmatter write) matching cycle 11-12's addressable-blocking-presenter and non-terminal-pane/timeout discipline; `## Provider id-mapping` implements the contract's normalize-after-digest-validation rule; AC-1/AC-2 carry both red fixtures as fixtures.
- DONE: The cross-repo dependency declared honestly: what needs subspace-tui surfaces vs what ships spacedock-side now; the working-copy-skill ritual as the measured interim baseline it must beat.
  `## What ships spacedock-side now vs. the subspace-tui dependency` splits the two: the full validate/launch/retain/id-map lifecycle ships spacedock-side (testable against a stub child); the briefing-package + `--result` coexistence and the non-EOF transport are the declared subspace-tui gaps. The `review-local-zellij` interim is named as the baseline, and AC-1 measures the command beating it (N/N vs 0/N retention).
- DONE: Expected surface + tolerance declared; riskiest unverified mechanism spiked end-to-end first or an auditable no-spike-needed recorded.
  `## Expected surface + tolerance` declares ~350-550 prod LOC ≈ equal test LOC, ~2-3 `internal/` files + 1 `cmd/` verb, tolerance 2×, with a hard self-check fencing off 3k/h1/subspace scope. `## Spike` records the end-to-end retention spike (7/7, `scratchpad/retention-spike.sh`) plus the exact-tip binary probe proving the provider retains nothing on the blank-float EOF path.

### Summary

Designed `spacedock gate review` as the presentation half split from 3k: a blocking presenter that owns a never-deleted retention directory and atomically retains result/log/diagnostics on every exit path, reproducing and beating the interim ritual's two destruction defects. The riskiest mechanism (atomic retention on success AND failure) was spiked end-to-end against the real tip binary before design lock — the provider writes nothing to `--result` on the blank-float EOF path, so the command owns retention. Scope holds strictly to the resolution/briefing side: no frontmatter writes (3k), no application/eligibility (h1), no subspace-tui product work.

## Stage Report: ideation (cycle 2)

Preflight fold applied (the first decline: the destroyed-approval fixture).

- DONE: Extend the evidence base to findings 1-15.
  Bumped both evidence lines (overview + Problem) to findings 1-15, cited 3k's attempt-7 resolution provenance note (verified at the source: index.md line 164), and marked finding 14 as the deliberate live-session numbering skip.
- DONE: Add finding 15 as the controller/launcher-death red fixture, with the AC-2 clause.
  Added a Problem bullet naming finding 15 (a dying launcher unlinked the captain's own approval) as the primary red fixture; extended AC-2 with the launcher-death clause (the caller-owned `--result` path survives because the launcher owns no cleanup over it) and its test. Exercised it: spike fixture C kills a real launcher process after it writes the result — interim scratch is unlinked (finding 15 reproduced), caller-owned result survives. Spike now 9/9.

### Summary

The fold strengthens the evidence base with finding 15 — the strongest red fixture for this exact mechanism, since it is the caller-owned retention directory's whole reason to exist. Verified the provenance at 3k's attempt-7 note rather than taking it on faith, then proved the fixture by exercising a real killed launcher process (spike 9/9), not by asserting it.

## Stage Report: ideation (cycle 3)

Captain gate feedback (attempt 1, revise): "what happens when the user does not have subspace installed? what's the fallback?"

- DONE: Design the no-subspace fallback — detection, fallback, shape.
  Added `## No-subspace fallback`: detection resolves + version-gates the TUI binary before any side effect and, on absence/mismatch, exits non-zero naming both the install remedy and the chat fallback with zero side effects; fallback is chat presentation recorded through the recorder exactly as every 0260 gate was, presentation-agnostic, under the recording-identity ruling (verified at `docs/roadmap/durable-decisions/index.md` Constraints line 25); shape mirrors the teams-unavailable-selects-bare-dispatch ruling (verified at `dispatch-failure-retry-rung.md` — an unavailable capability is not a mode, just the ordinary condition selecting the alternate path). Made detection step 1 of Required capability, before any side effect.
- DONE: Make it checkable — one AC clause + test-plan line, exercised.
  Strengthened AC-5 to require non-zero exit + a message naming both remedy and fallback + zero launch + no retention directory created (entity byte-unchanged), with the matching test-plan line; exercised it as spike fixture D (absent presenter → exit 3, message names remedy+fallback, no dir created). Spike now 12/12. Added the fallback sentence to the doc diff.

### Summary

The no-subspace answer names an existing practice rather than inventing one: detection fails clean and names the fallback, chat presentation records through the same recorder under the recording-identity ruling, and a missing presenter is an ordinary channel selection — not a mode, never blocking the gate. Grounded both cited rulings at their sources and proved the checkable clause by exercising fixture D (12/12).

## Stage Report: ideation (cycle 4)

Captain gate feedback (attempt 2, revise): "i am wondering if this should be left as an overridable of the current present-gate skill, so that they are not coupled."

- DONE: Assess the decoupling honestly, then reframe the architecture around it.
  Read the present-gate skill (its default and only behavior is chat presentation) and confirmed the reframe fits: presentation becomes an overridable present-gate channel (default chat, override the hardened float script), the spacedock binary carries zero Subspace code, and the cross-repo release-train coupling — this design's biggest declared liability — dissolves to the opt-in override script. Rewrote the overview, Expected surface, Problem close, the channel section (replacing the binary lifecycle), and the id-mapping section.
- DONE: Split the duties by their honest home; move the id-mapping implementation and amend the owner tag.
  Added `## Splitting the duties honestly` (table): retention → override script + committed drive suite (the 12 spike fixtures); result validation + id-normalization + record handoff → recorder-side (already 3k's parse-and-verify duty); detection/fallback → the skill's channel selection. Proposed the owner-tag amendment (id-mapping specified AND implemented recorder-side) to 3k via the change protocol; did not edit 3k's surface.
- DONE: Re-estimate the surface with tolerance; state the honest counter-case.
  Spacedock Go surface for presentation → ~zero. New surface declared: present-gate skill prose (~25-40 lines) + one hardened override script (~80-140 LOC + the 12-fixture drive suite, homed in the subspace repo) + a small recorder-verb ask; tolerance 2× against that surface. `## Honest assessment` records the one load-bearing condition — the override script MUST carry a committed CI-run drive suite or it repeats `review-local-zellij`'s untested-script defect (the exact class this task exists to kill); that condition fixes the script's home (subspace repo, already `go test`-covered). Found no piece that genuinely cannot be a skill override without losing a guarantee.
- DONE: Keep the no-subspace fallback, probe ritual, and every red fixture; make the decoupling itself checkable.
  Fallback is now the default channel (unchanged behavior); probe ritual and fixtures A-D retained. Added AC-6 (VALUE): the spacedock binary depends on no Subspace package and exposes no presentation verb — a build/dependency assertion (not a prose-grep) whose count regresses the wrong way if a coupling returns. Reframed AC-1/2/3/5, the test plan, and the doc diff (present-gate channel, no `gate review` verb).

### Summary

The captain's question is well-founded and the answer is yes: the binary was the wrong vehicle. Presentation moves to an overridable present-gate channel, the spacedock binary ends Subspace-free (AC-6, measurable), and every guarantee relocates without loss — retention to a testable override script (the spike already proved the contract in bash), validation and id-mapping to the recorder where 3k's contract already puts result verification. The single load-bearing condition, surfaced for the captain rather than buried: the override script must carry a committed drive suite (the 12 fixtures), or the reframe reintroduces the untested-script defect it exists to remove.

### Feedback Cycles

- Cycle 1: REJECTED — Roborev job 541; surface 4 files/49 added, 3 removed vs estimate skill 25-40 lines + tests 20-60 lines + docs 8-16 lines (within 2×); AC unchanged.
- Cycle 2: PASSED — Roborev job 542; surface 4 files/49 added, 3 removed vs estimate skill 25-40 lines + tests 20-60 lines + docs 8-16 lines (within 2×); AC unchanged.
- Cycle 3: REJECTED — fresh validation and Roborev job 1955; surface 4 files/49 added, 3 removed vs estimate skill 25-40 lines + tests 20-60 lines + docs 8-16 lines (within 2×); AC unchanged.
- Cycle 4: REJECTED — Captain design reset after validation cycle 3: replace advisory-adoption plumbing with the gate-room/minimal-binding Result boundary and derive provider association behind the scaffold; surface 4 files/49 added, 3 removed vs estimate skill 25-40 lines + tests 20-60 lines + docs 8-16 lines (within 2×); AC unchanged.
- Cycle 5: REJECTED — Roborev job 1974; surface 15 files/708 LOC vs pre-edit declaration 9 source/docs files plus fixture set/~585 changed lines (121%); AC unchanged.
- Cycle 6: DESIGN RESET — First Officer scope ruling under the Captain's sprint conn; surface 17 files/1,410 changed LOC vs rework declaration ~585 (241%, beyond 2×); AC unchanged. Remove the 270-line modeled provider smoke/contract interpreter rather than repair it into a second implementation. Spacedock's existing CLI fixtures own prepared-room integrity and failure atomicity; Subspace owns executable `/subspace:r gate <room>` transport behavior; retain only a minimal structural assertion of the declarative seam if needed.
- Cycle 7: REJECTED — Roborev job 2019 and First Officer authority triage; surface 16 files/1,178 changed LOC after the Cycle-6 harness cut vs rework declaration ~585 (201%); AC unchanged. The `agent:first-officer == agent:first-officer` gate-room path is Material because `request.json` carries no independently verifiable Captain delegation. v1 provider rooms are Captain decision rooms; delegated FO approval remains the existing `--decision ... --directive` path. Pre-bind malformed/stale request validation is Material and must be byte-clean. The request for an in-repo executable provider orchestrator is declined by the Cycle-6 ownership ruling.
- Cycle 8: REJECTED — Roborev job 2022; surface 16 files/1,159 changed LOC vs rework declaration ~585 (198%, within 2×); AC unchanged. Nested provider `Resolution.adoption-note` is Material because a direct binding Captain room must not persist advisory-adoption provenance. Reject it before mutation and retain the existing separate delegated chat-decision path.
- Cycle 9: REJECTED — Roborev job 2025 found that a closed provider attempt discarded the only durable identity of the exact Result and presented inventory it accepted, so later replacement of either room file escaped `gate validate`. This is Material to AC-4's exact-attempt binding: Git retains prior bytes but the attempt cannot identify which revision was authoritative. The approved correction persists only `result-digest` and `presented-inventory-digest`, freezes them with the closed attempt, and makes `gate validate` recompute both fixed room files. Association remains derived and unstored; transaction logs, evidence registries, artifact copies, provider transport, and a second harness remain excluded. Pre-edit declaration: `internal/gates/model.go` +16, `operation.go` +4, `io.go` +35, `internal/cli/gate_test.go` +45/-9, and the gate contract +12/-2; expected +112/-11 from the 1,170-line current surface, with a new hard cap of 1,300 changed lines and no file-count increase.
- Cycle 10: REJECTED — validation cycle 4 design/spec reset; surface 17 files/1,255 changed LOC vs Cycle-9 cap 1,300 (97%); AC reset required. The candidate has no material prepared-room defect, but AC-1 through AC-5 still require the retired override script, provider lifecycle, presenter fallback, and provider-minted-ID normalization. Decision: re-scope ideation only to replace those obsolete obligations with exact room authority/identity, derived recursive inventory, failure atomicity, one-use history, and retained-byte tamper detection while preserving the zero-coupling and decision-integrity end values; do not send implementation back to rebuild excluded provider machinery.
- Cycle 11: REJECTED — independent ideation staff review; surface 17 files/1,255 changed LOC vs cap 1,300 (97%); AC unchanged. Three supported-path evidence gaps remain: no cited post-bind Briefing-byte mutation proves AC-2, the AC-3 inventory matrix lacks same-cardinality duplicate-id and wrong-revision cases, and AC-5 lacks deletion cases for both retained provider files. These are test-only falsifiability gaps, not new interface or provider obligations. Revise the test plan and intended surface around compact additions to the existing CLI matrices; do not add production behavior, a new harness, or provider transport.
- Cycle 12: REJECTED — independent ideation staff re-review; surface unchanged at 17 files/1,255 changed LOC vs cap 1,300 (97%); AC unchanged. The Cycle-11 AC-5 test prose accidentally required missing retained files to emit “frozen digest,” but current supported behavior names the fixed path plus the read/not-found error; only byte mismatches use the frozen-digest diagnostic. Remove that invented message obligation, preserve the distinct assertions, and keep the test-only +42/-0 estimate unchanged.
- Cycle 13: REJECTED — Roborev job 2039 implementation round 1; surface 17 files/1,293 changed LOC vs estimate 1,300 changed LOC (99%); AC unchanged

## Stage Report: implementation

- DONE: Ship the present-gate overridable-channel contract with chat fallback, complete canonical Briefing presentation, exact retained Result handoff, and honest association to what the reviewer actually saw.
  Commits `cf6008fd` and `612b72fc` add default chat, the six-part override contract, exact retained Result handoff to `gate record`, and one-to-one provider Artifact/Reference mapping to the complete canonical inventory.
- DONE: Prove missing/mismatched presenter falls back without side effects and the full-package journey never promotes a single-file float into a complete-package association.
  The skill makes availability/version probing precede package and retention creation; the provider's 12-fixture drive suite remains the required pinned cross-repo proof and is not claimed as a local run, while `TestGateRecordConsumesExactResultOnlyWithCompleteAssociation` rejects the primary-only fixture and proves the entity stays byte-unchanged.
- DONE: Keep the Spacedock binary Subspace-free and preserve the provider-owned hardened override script plus committed drive suite as the named pinned cross-repo release condition.
  `go list -deps ./cmd/spacedock` returns no Subspace dependency; `TestGatePresentationRemainsOutsideBinary` fails if a `gate review` verb appears or if rejecting it changes the working directory; the skill names the pinned provider script and CI suite as release eligibility.
- DONE: Reconcile implementation surface against the pre-edit declaration.
  Actual surface is 15 skill lines, 8 docs lines, and 26 test-line additions with 3 replacements across 4 files, below the declared prose estimate and within the 2× test/docs tolerance; no Go production, recorder, provider, gate-state, or frontmatter surface changed.
- DONE: Run repository and documentation verification.
  PASS (exit 0): `go test ./...`; PASS (exit 0): `go test ./... -race`; PASS (exit 0): `mkdocs build --strict`; PASS (exit 0): `gofmt -w ./cmd ./internal`; PASS (exit 0): `git diff --check`; final code-worktree `git status --short` output was empty (clean).
- DONE: Request and triage Roborev on the completed commits.
  Job 541 found one material AC-5 ambiguity about provider References; `612b72fc` made every Artifact/Reference mapping explicit without expanding recorder scope, and job 542 returned `No issues found.`

### Summary

The present-gate skill now owns an overridable rendering channel: chat remains the side-effect-free default, and a provider override must present the complete canonical package, retain exact results and diagnostics, and associate only content the reviewer saw. The implementation leaves the binary Subspace-free and keeps the hardened provider script plus its committed drive suite as the pinned cross-repo release condition.

## Stage Report: validation

- DONE: Reproduce the complete canonical Briefing journey: exact question and every Artifact/Reference presented, exact retained Result handed off, and association limited to content the reviewer actually saw.
  `TestExactCanonicalBriefingIsIndependentAssociationInventory` binds the exact question and three-Artifact inventory at digest `sha256:0a54f1ba...`; `TestGateRecordConsumesExactResultOnlyWithCompleteAssociation` maps `artifact:primary`, `reference:entity-snapshot`, and `reference:recorder-contract` one-to-one, then records exact Result digest `sha256:46096103...` and decision `revise`.
- SKIPPED: Attack missing/mismatched presenter fallback and primary-only/single-file promotion paths, proving zero side effects on fallback and fail-closed no-mutation association rejection.
  The provider fallback fixtures are absent from this repository and were not claimed; the in-repo primary-only association attack returned exit 1 with `complete presentation mapping` and left the entity byte-identical.
- DONE: Verify the binary stays Subspace-free, the provider script/committed-drive-suite remains an honestly named pinned release condition, and all ACs/tests/surface claims survive detached adversarial review.
  `go list -deps ./cmd/spacedock` found zero Subspace dependencies; `gate review` remains absent; the skill makes the provider script plus committed CI suite a release condition, and the detached audit found no in-repo material defect.
- SKIPPED: AC-1 (VALUE) — No presented decision is lost on any exit path.
  This proof belongs to the provider-owned 12-fixture drive suite; no provider repository or pinned revision was supplied, so validation did not claim a local run.
- SKIPPED: AC-2 — Retention survives every failure class, including launcher/controller death.
  The launcher-death, retention-write, crash, EOF, hold, and validation-failure fixtures remain the same provider-owned release condition.
- SKIPPED: AC-3 — Pane/session creation and wait-timeout are never completion.
  The blocking-child and return-on-pane mutants require the provider script and committed drive suite, which are outside this checkout.
- DONE: AC-4 — The recorded result is keyed to the attempt briefing id, only after digest validation (recorder-homed, proposed).
  The matching fixture normalized both provider ids to the attempt id; an independent mismatch run returned exit 1 and preserved the entity and Result bytes. Removing the revision check made that audit fail by closing the gate.
- SKIPPED: AC-5 — The override channel validates the briefing and derives the title before any launch; an absent or version-mismatched presenter falls back to chat with zero side effects.
  The shipped skill states the probe-first, zero-side-effect fallback contract, but only the provider-owned suite can prove launch count, title derivation, and retention-directory absence.
- DONE: AC-6 (VALUE) — Presentation adds zero Subspace coupling to the spacedock binary, and no channel mutates entity frontmatter.
  Dependency count was zero; the absent-verb test returned exit 2 and left its working directory unchanged; all recorder rejection controls preserved the entity bytes.
- DONE: Audit the declared 4-file/15-skill-line/8-doc-line/26-test-line surface.
  Diff from `fa240a76` is exactly 4 files: 15 skill additions, 8 docs additions, 21 CLI-test additions, and 5 additions/3 replacements in contractlint; no Go production, provider, recorder, or frontmatter file changed.
- DONE: Run detached adversarial controls.
  Mutants accepting `gate review`, trusting association-declared inventory, skipping exact-Result digest binding, changing the exact question, and skipping canonical-revision binding each broke the named focused test or independent audit.
- DONE: Verify Roborev jobs 541/542 closed all material findings without ownership crossover.
  Job 541's sole Reference-association finding is corrected by `612b72fc`; job 542 reports `No issues found.` No provider, recorder, or gate-state ownership crossed. Roborev metadata still marks both review records `closed:false`, an administrative state rather than an unresolved code finding.
- DONE: Run normal, race, documentation, formatting, and cleanliness checks.
  PASS: `go test ./...`, `go test ./... -race`, pinned-env `mkdocs build --strict`, `gofmt -w ./cmd ./internal` with no diff, `git diff --check`, and clean committed implementation worktree at `612b72fc`.
- DONE: Recommend PASSED for the in-repo deliverable, with the cross-repo release condition explicit.
  No material in-repo finding remains; release eligibility still requires a pinned provider revision carrying the hardened override script and its committed 12-fixture CI drive suite.

### Summary

Fresh detached validation passed the four-file Spacedock deliverable and five claim-breaking controls. The provider transport remains deliberately outside this repository: its pinned script and 12-fixture CI suite are an unmet cross-repo release condition, not local test evidence.

## Stage Report: validation (cycle 2)

- DONE: Run and pin the provider-owned retained-delivery drive suite at sibling Subspace revision 198f76238aeb74ff38900e17b751f0460d0c55ee, mapping results explicitly to previously skipped AC-1, AC-2, AC-3, and AC-5.
  `scripts/tests/subspace-r-provider-retained-delivery-test.sh` passed at exact `198f7623`; its approve/revise/hold/open, blank/EOF, crash, invalid-result, retention-write, launcher-death, alive-child, missing/mismatched-presenter, complete-package, and title rows fail on deletion, early delivery, relaunch, or inventory drift.
- DONE: Replay the complete canonical Briefing, Result, and presented-inventory association at Spacedock candidate 612b72fc; distinguish provider defects from the current Codex/Safehouse Zellij transport limitation and make no sibling-repository edits.
  Subspace's frozen Result `sha256:46096103...` and association `sha256:95ca15ab...` are byte-identical to Spacedock fixtures; recursive inventory checks and the exact recorder test passed, while the primary-only map failed without mutation. No provider defect appeared; no headed captain float is claimed because this Codex turn does not expose the agreed `/subspace:r gate <gate-room>` surface, and the private 4p vector is not a substitute.
- DONE: Reissue a fresh exact-tip PASSED or REJECTED recommendation for all six ACs, preserving the zero-Subspace binary boundary and treating the verbose internal 4p vector as implementation plumbing rather than xb's agent-facing interface.
  **PASSED** at Spacedock `612b72fc` with provider `198f7623`: all six ACs have executable or retained-state evidence, zero material findings remain, and `/subspace:r <file.md>` plus `/subspace:r gate <gate-room>` remains the public shape.
- DONE: AC-1 (VALUE) — No presented decision is lost on any exit path.
  The provider matrix retained complete bundles for approve, revise, hold, and open; blank/EOF, crash, validation failure, retention failure, and launcher death retained every produced result/log/inventory/diagnostic byte, so removing any case artifact fails the suite.
- DONE: AC-2 — Retention survives every failure class, including launcher/controller death.
  Child exit `42` and launcher exit `43` propagated once; launcher-death kept the non-empty Result, log, inventory, argv/stderr diagnostics, and death marker, while retention-write failure kept the Result plus its error.
- DONE: AC-3 — Pane/session creation and wait-timeout are never completion.
  The alive-child fixture published a pane marker and Result while holding the exact child alive; entry return, validation, child-exit publication, and delivery stayed absent until release, then occurred exactly once.
- DONE: AC-4 — The recorded result is keyed to the attempt briefing id, only after digest validation (recorder-homed, proposed).
  The exact Spacedock recorder replay normalized the provider envelope only through the complete digest/revision-bound association; missing/changed nested Reference and canonical id/revision mutations all failed closed.
- DONE: AC-5 — The override channel validates the briefing and derives the title before any launch; an absent or version-mismatched presenter falls back to chat with zero side effects.
  The complete provider row derived `Subspace — Ship the complete package?` and ordered 2 Artifacts plus 2 recursively reached References; missing/mismatched presenter rows returned `127`/`2` before host preflight or launch, leaving chat selection to the probe-first Spacedock channel contract.
- DONE: AC-6 (VALUE) — Presentation adds zero Subspace coupling to the spacedock binary, and no channel mutates entity frontmatter.
  `go list -deps ./cmd/spacedock` found zero Subspace dependencies; the absent `gate review` command and recorder rejection controls changed no working directory or entity bytes.
- DONE: Run the required repository verification at the exact Spacedock tip.
  PASS: focused cross-repo tests, `go test ./...`, `go test ./... -race`, `uv run --with-requirements docs/requirements.txt mkdocs build --strict`, `gofmt -w ./cmd ./internal` with no diff, `git diff --check`, and clean code worktree.

### Summary

The previously skipped provider-owned evidence is now pinned and reproduced, and the exact provider Result/inventory/association crosses into Spacedock unchanged and fails closed under independent mutations. Validation recommends **PASSED** for all six ACs at the exact tips; the provider's verbose deterministic vector remains private plumbing, not a new public Spacedock or agent-facing interface.

## Stage Report: implementation (rework)

- DONE: Rebase `spacedock-ensign/gate-review-presentation-command` onto current `origin/main` and resolve only the two known semantic conflicts in `internal/cli/gate_test.go` and `internal/contractlint/fo_function_reference_invariant_test.go`.
  Rebasing `612b72fc` onto `dd6bd114` produced only those conflicts; candidate `4779fff8` preserves main's gate tests and adds the presentation refusal beside them, while host baselines equal main plus xb's measured 2,872-byte load.
- DONE: Preserve xb's approved behavior and intended four-file delta: presentation stays an overridable `present-gate` channel; complete association remains required; the binary remains Subspace-free.
  The final diff against `origin/main` is exactly the original four files and `+49/-3`; complete provider Artifact/Reference mapping remains in the skill/docs, the recorder rejection test passes, and `go list -deps ./cmd/spacedock` names no Subspace dependency.
- DONE: Reconcile with merged gate-lifecycle/advisory-round behavior from commits `03b1a7fc`, `f06cce04`, `c355fbe4`, and `c9dfc491`; do not invent compatibility constraints for prototype behavior.
  All four commits are ancestors of `origin/main`; round record/validate, owning-workflow discovery, eligibility/consume-once, and presentation-refusal tests pass, with the latter expecting main's real `record|validate|eligibility|consume` command surface.
- DONE: Before committing, report actual changed files/LOC versus xb's pre-rebase 4 files, +49/-3, and treat unexplained product-surface drift as a blocker.
  Pre-commit audit reported docs `+8`, CLI test `+21`, load-ratchet test `+5/-3`, and present-gate skill `+15`; no extra file, production Go surface, provider/recorder change, or unexplained drift appeared.
- DONE: Run focused affected tests, gofmt as applicable, `go test ./...`, `go test ./... -race`, mkdocs strict, and `git diff --check`.
  Focused CLI/contractlint suites, `gofmt -w ./cmd ./internal`, normal/race suites, `uv run --with-requirements docs/requirements.txt mkdocs build --strict`, and diff checks all passed; removing a merged verb, complete association, or per-host load increment breaks its named test.
- DONE: Commit the rebased integration candidate and write an implementation/rework report, but do not self-validate or mutate workflow frontmatter/gates.
  Rebased commits are `f8fd5c2a` and `425812bc`; integration reconciliation is `4779fff8`. This append-only body report makes no validation recommendation and leaves frontmatter/gates untouched for a fresh validator.

### Summary

The rebased xb candidate preserves the approved four-file `+49/-3` boundary while composing with merged gate lifecycle and folder-form advisory rounds. Candidate `4779fff8` is clean and fully implementation-tested; independent validation remains the next gate.

## Stage Report: validation (cycle 3)

- DONE: Independently verify rebased candidate 4779fff8 preserves xb’s six ACs, exact four-file +49/-3 boundary, and composes semantically with the merged gate lifecycle/advisory-round tests without compatibility inventions.
  The diff from `dd6bd114` is exactly docs `+8`, CLI test `+21`, contractlint `+5/-3`, and skill `+15`; round record/validate, owning-workflow, eligibility/consume-once, absent-presentation-verb, association, and host-load tests pass against the real `record|validate|eligibility|consume` surface.
- DONE: Reproduce the pinned Subspace 198f762 retained-delivery suite and complete Briefing/Result/association path against candidate 4779fff8; explicitly establish that the old 612b72fc approval cannot authorize this new candidate and identify the replacement Briefing inputs.
  Exact provider commit `198f76238aeb74ff38900e17b751f0460d0c55ee` passed its 12-fixture suite; the exact Result/full-association recorder path passes, while an advisory Result with that association but no adoption note fails with `advisory Result requires --adoption-note`.
- DONE: Establish approval freshness and replacement Briefing inputs.
  The closed attempt’s immutable question and approval reason name only `612b72fc`; although `gate eligibility` reports its application pending/eligible, it cannot authorize `4779fff8`. A successor Briefing must freeze the `4779fff8`/`dd6bd114` four-file diff, current entity plus this validation, merged lifecycle/advisory evidence, provider `198f7623` suite evidence, and the corrected handoff; only its fresh exact Result and complete association can authorize landing.
- DONE: AC-1 (VALUE) — No presented decision is lost on any exit path.
  Provider rows for approve, revise, hold/open, blank/EOF, crash, invalid result, retention failure, and launcher death retain their promised bytes; deleting any retained artifact breaks the pinned suite.
- DONE: AC-2 — Retention survives every failure class, including launcher/controller death.
  The pinned suite preserves Result/log/inventory/diagnostics across child and launcher failures and proves nonzero status propagation without relaunch.
- DONE: AC-3 — Pane/session creation and wait-timeout are never completion.
  The alive-child row withholds delivery and validation until the blocking child exits, then publishes exactly once.
- FAILED: AC-4 — The recorded result is keyed to the attempt briefing id, only after digest validation (recorder-homed, proposed).
  The recorder’s complete association and digest normalization work, but the shipped skill’s prescribed command omits the required `--adoption-note`; Subspace’s advisory Result therefore fails at the supported handoff boundary. This is a material outcome defect and a narrow same-layer fix, not a design reset.
- DONE: AC-5 — The override channel validates the briefing and derives the title before any launch; an absent or version-mismatched presenter falls back to chat with zero side effects.
  Complete-package/title and missing/mismatched-presenter rows pass at the pinned provider revision before host launch or retention creation.
- DONE: AC-6 (VALUE) — Presentation adds zero Subspace coupling to the spacedock binary, and no channel mutates entity frontmatter.
  `go list -deps ./cmd/spacedock` reports zero Subspace dependencies; absent-verb and recorder rejection controls preserve the working directory/entity.
- DONE: Run focused tests, gofmt, go test ./..., go test ./... -race, mkdocs strict, diff/cleanliness checks, detached adversarial review and Roborev as applicable; classify every finding by materiality and issue a fresh exact-tip PASSED or REJECTED report.
  Final uncontaminated normal/race suites, strict MkDocs, gofmt-no-diff, diff check, and clean code worktree pass; Roborev panel job 1955 is REJECTED at exact tip.
- FAILED: Material outcome finding — advisory Result handoff omits adoption authority.
  Supported trigger: provider returns `status:"advisory"`; `gate record` exits 1 without `--adoption-note`. Fix AC-4 by branching on Result authority, supplying captain-authorized adoption text, and behaviorally testing the exact prescribed invocation.
- SKIPPED: Deferred evidence risk — provider revision is not pinned by a code-repository release check.
  Current gate evidence pins and reproduces `198f7623`, so no present AC lacks proof; promote if a stable release consumes the skill without a fresh exact provider pin and suite result.
- SKIPPED: Deferred evidence risk — absent-verb test compares directory entry count, and host-load ratchets enforce ceilings rather than exact byte equality.
  Current early-return path creates nothing and all hosts equal main plus 2,872 bytes; promote the first if the early-return guard moves, and the second if ratchets become exact accounting rather than upper bounds.
- FAILED: Fresh exact-tip recommendation.
  **REJECTED** at Spacedock `4779fff8` with provider `198f7623`: AC-4’s documented supported handoff cannot record the provider’s advisory Result; no old `612b72fc` Resolution authorizes this candidate.

### Summary

The rebase, four-file boundary, provider retention suite, merged gate behavior, and five ACs validate cleanly. One material handoff defect remains: the present-gate override contract must carry captain-authorized adoption for advisory Results and prove that exact invocation before a fresh candidate-bound Briefing can be approved.

## Stage Report: implementation (cycle 2)

- DONE: Replace the rejected advisory-adoption handoff with one prepared-room recorder boundary.
  Commit `d933db27` makes `spacedock gate record <entity> --room <gate-room>` the only provider form; callers no longer supply Result, association, actor, approver, output paths, adoption text, or provider argv.
- DONE: Keep request/package origination in the scaffold and derive every provider input behind the room.
  The recorder consumes fixed `request.json`, canonical `briefing.json`, `provider/result.json`, and `provider/presented-inventory.json`; no room-preparation command, provider launcher, Subspace import, or transport adapter was added.
- DONE: Freeze the complete request-backed attempt before provider recording.
  Commits `c49be15f`, `c0f8f45b`, `6e50dba5`, and `290b9249` bind the JCS request digest, reject open/successor rebinding, validate derived gate/attempt/Briefing identity before every write, and require `actor == approver == person:captain`.
- DONE: Accept only a wrapper-free direct Captain Result and preserve honest provenance.
  Independent status/binding/actor/approver/resolutionId cases, nested `adoption-note`, distinct authority, changed request authority, and advisory output all fail before mutation; delegated FO approval remains the separate `--decision ... --directive` chat path.
- DONE: Require complete Artifact plus recursive Reference presentation and an Artifact-typed Result primary.
  The recorder derives the canonical inventory recursively, rejects missing/mistyped/duplicate mappings and a Reference used as the primary Artifact, and never accepts a caller-built association.
- DONE: Bind exact retained provider bytes to immutable durable state.
  Commit `98ebb458` stores only raw Result and presented-inventory digests on the closed attempt; `gate validate` recomputes both fixed files and the focused test proves either post-close byte change fails until restored.
- DONE: Preserve one-use/frozen history and failure atomicity.
  Every malformed request, invalid inventory, unauthorized Result, stale successor, room move, advisory wrapper, and retained-file tamper test asserts rejected entity bytes are unchanged; whole-attempt transition validation freezes the stored evidence digests.
- DONE: Keep executable provider transport at the ruled ownership boundary.
  Cycle 6 removed the modeled 270-line interpreter; `/subspace:r gate <gate-room>` remains Subspace-owned executable evidence, while one contractlint assertion pins probe-before-side-effects, prepared-room invocation, room-only recording, and no post-launch chat fallback.
- DONE: Reconcile actual surface against the declarations and approved resets.
  Final diff from `main` is 17 files, `+1086/-169` = 1,255 changed LOC; this exceeds the original ~585-line rework declaration but stays below Cycle 9's approved 1,300-line cap. The authorized extra file is `internal/gates/io.go`; the proposed integration harness is absent.
- DONE: Run focused, repository, race, documentation, formatting, dependency, and cleanliness checks.
  PASS: focused CLI/gates/contractlint controls, `go test ./...`, `go test ./... -race`, `uv run --with-requirements docs/requirements.txt mkdocs build --strict`, `gofmt -w ./cmd ./internal`, `git diff --check`, clean worktree, and zero Subspace entries in `go list -deps ./cmd/spacedock`.
- DONE: Request final-tip adversarial review and triage every correction round.
  Jobs 1974, 1984, 1987, 1993, 2003, 2014, 2019, 2022, and 2025 drove the recorded authority, binding, inventory, provenance, scope, and retention corrections; final-tip job 2028 at `98ebb458` returned `No issues found.`
- DONE: Commit a reviewable exact-tip candidate without mutating workflow frontmatter or gate state.
  The implementation branch is clean at `98ebb458`; this append-only report changes only the entity body and leaves `status`, `gates`, and application state untouched for independent validation.

### Summary

The corrected implementation exposes one room-only provider boundary, validates and freezes Captain authority before mutation, derives the complete presentation association internally, and binds the exact retained Result/inventory bytes to the closed attempt. It preserves zero Subspace binary coupling and the Subspace-owned transport boundary, passes all required checks, and cleared final-tip Roborev job 2028 with no findings.

## Stage Report: validation (cycle 4)

- FAILED: Reconcile all six acceptance criteria with the Captain-approved prepared-room recorder boundary; reject stale provider-side obligations rather than proving the wrong target.
  The binding AC section still specifies the retired override script, provider-minted-id normalization, presenter fallback, and provider drive suite; Cycle 4-9 instead excludes provider transport and makes `gate record <entity> --room <gate-room>` the shipped boundary. This is a design/spec reset, not an implementation feedback cycle.
- SKIPPED: AC-1 (VALUE) — No presented decision is lost on any exit path.
  Its cited override-script retention matrix is provider-owned and intentionally absent after Cycle 6; it cannot prove the current Spacedock recorder deliverable.
- SKIPPED: AC-2 — Retention survives every failure class, including launcher/controller death.
  Launcher/controller lifecycle is outside the approved repository scope; Spacedock now consumes already-retained fixed room files and owns their post-close digest verification.
- SKIPPED: AC-3 — Pane/session creation and wait-timeout are never completion.
  Pane, session, and provider wait behavior belongs to `/subspace:r gate <gate-room>` and present-gate orchestration, not the recorder command.
- SKIPPED: AC-4 — The recorded result is keyed to the attempt briefing id, only after digest validation (recorder-homed, proposed).
  Provider-minted-id normalization is obsolete: the approved recorder requires the Result to name the exact canonical Briefing, then verifies frozen request/Briefing identity and raw Result bytes; the replacement invariant passes the focused matrix.
- SKIPPED: AC-5 — The override channel validates the briefing and derives the title before any launch; an absent or version-mismatched presenter falls back to chat with zero side effects.
  Title derivation, provider probing, and fallback are provider/orchestrator behavior excluded from the prepared-room recorder; the repository carries only the declarative seam.
- DONE: AC-6 (VALUE) — Presentation adds zero Subspace coupling to the spacedock binary, and no channel mutates entity frontmatter.
  `go list -deps ./cmd/spacedock` reports zero Subspace dependencies; the absent `gate review` verb and every recorder rejection case leave entity bytes unchanged.
- DONE: Reproduce exact-tip authority, identity, recursive inventory, failure atomicity, one-use history, and post-close Result/inventory tamper detection across adjacent variants.
  Uncached CLI tests at `98ebb458` reject distinct/rebound authority, wrong gate/attempt/Briefing/digest/room, advisory wrappers, nested adoption provenance, incomplete/mistyped recursive inventory, Reference primaries, stale successors, and both retained-file mutations without entity changes; valid close and consume occur once.
- DONE: Audit the 17-file/1,255-line surface against Cycle 9's 1,300-line cap and exclusions; verify focused/full/race/docs/dependency checks plus final-tip Roborev 2028.
  Diff from `main` is exactly 17 files and `+1086/-169` (1,255 changed lines), with no provider transport or second harness; focused tests, `go test ./...`, `go test ./... -race`, strict MkDocs, gofmt, diff/cleanliness, and zero-Subspace dependency checks pass. Roborev job 2028 covers `dd6bd114..98ebb458`, verdict P, `No issues found.`
- SKIPPED: Deferred risk — closed provider attempts require retained room files forever.
  Trigger is later archival, deletion, or reformatting of either fixed provider file; that is outside the current retain-in-place contract, whose supported path validates cleanly. Promote if room cleanup or archival becomes supported.
- FAILED: Fresh exact-tip recommendation.
  **REJECTED — DESIGN/SPEC RESET.** The implementation at `98ebb458` has no material prepared-room-recorder defect, but the six binding ACs must be replaced with criteria for exact room authority/identity, derived recursive inventory, failure atomicity, one-use history, and retained-byte tamper detection before this validation gate can pass.

### Summary

The Captain-approved prepared-room recorder passes its current behavioral, surface, repository, and final-tip review checks with no material implementation finding. Validation rejects only because the entity still binds acceptance to the superseded provider-script architecture; ideation must reset the ACs rather than implementation rebuilding intentionally excluded transport and retention machinery.

## Stage Report: ideation (cycle 5)

- DONE: Replace obsolete provider-side obligations with the prepared-room recorder contract.
  Rewrote the problem, approach, acceptance criteria, mechanisms, test plan, documentation proposal, and exclusions around one scaffold-prepared room consumed by `spacedock gate record <entity> --room <gate-room>`. The binding design now requires exact request/Briefing/attempt identity and Captain authority, a direct wrapper-free Result, complete recursively derived presentation inventory, atomic one-use close, frozen raw Result/inventory digests, and later tamper detection. Override-script, pane/session, launcher-death, fallback, provider-minted-ID normalization, adoption-note, provider-harness, and duplicate-association obligations are explicitly outside xb rather than criteria implementation must satisfy.
- DONE: Map every rewritten criterion to exact-tip evidence and explicit remaining gaps.
  AC-1 through AC-5 cite the existing CLI/gates tests for the canonical close, adjacent authority and identity mutations, recursive inventory, byte-clean rejection, frozen lifecycle, exact retained-byte digests, and byte tampering; Cycle 11 now names the three missing negative cases without claiming they already ran. AC-6 cites the absent-presentation-verb test and zero-Subspace dependency check. The proof basis retains final-tip Roborev job 2028's `P — No issues found`; no replacement provider harness or recording machinery is proposed.
- DONE: Rebaseline the approved surface and decide whether product work remains.
  The clean candidate remains exactly 17 files and `+1086/-169` = 1,255 changed lines. Cycle 11 budgets an exact estimated `+42/-0` in existing `internal/cli/gate_test.go`, projecting 17 files/1,297 changed lines under the approved 1,300-line cap. Only these three focused test gaps remain; no production, fixture-file, interface, provider, harness, or documentation edit is justified.
- DONE: Fold the Cycle 11 staff-review evidence gaps into the existing proof surfaces.
  The ACs and behavioral plan now require a post-bind Briefing-question mutation, same-cardinality duplicate-id and wrong-revision inventory variants, and deletion of each retained provider file after close. Each addition reuses the existing CLI mutation matrices or retained-file loop and asserts byte-clean rejection or read-only validation failure; no new mechanism or harness is introduced.
- DONE: Fold Cycle 12's supported diagnostic split into AC-5 without changing scope.
  The retained-file loop now expects a present-but-byte-changed file to name its fixed path plus `frozen digest`, while deletion expects the fixed path plus the read/not-found error. It requires no invented digest message when bytes cannot be read and leaves the `+42/-0` estimate unchanged.
- DONE: AC-1 — Specify proof that only the exact authorized room makes a decision durable.
  At exact tip `98ebb458`, `TestGateRecordConsumesDirectBindingResultFromPreparedRoom` closes the canonical fixture once, while the existing authority, room, briefing-only, advisory, and distinct-identity matrices reject adjacent mutations with byte-identical entity state. Cycle 11 adds the missing valid-Briefing question mutation after binding to that same byte-clean room-record matrix.
- DONE: AC-2 — Specify proof that request, Briefing, attempt, and authority identity freeze.
  At exact tip `98ebb458`, the request-authority, room-move, briefing-only, malformed-Reference, and distinct-identity tests reject malformed, stale, moved, and rebound identity before mutation. Cycle 11 adds an explicit post-bind question-byte change and requires room recording to reject with the entity byte-identical.
- DONE: AC-3 — Specify proof of direct Captain Result and complete recursive presentation.
  At exact tip `98ebb458`, `TestGateRecordConsumesDirectBindingResultFromPreparedRoom`, `TestGateRoomRejectsAdvisoryEvidenceAndUnauthorizedResolutionWithoutMutation`, and `TestExactCanonicalBriefingIsIndependentAssociationInventory` cover direct authority and incomplete/mistyped presentation. Cycle 11 adds same-cardinality duplicate-id and wrong-revision rows to the existing inventory mutation matrix so cardinality alone cannot satisfy completeness.
- DONE: AC-4 — Prove atomic one-use close and frozen history.
  At exact tip `98ebb458`, `TestCanonicalLifecycleRebindCloseFreezeAndSupersede`, `TestWriterCASValidationAtomicityAndLock`, `TestGateRecordConsumesDirectBindingResultFromPreparedRoom`, and `TestGateEligibilityAndConsumeAuthorizeOnce` cover close-before-replace validation, byte-clean rejection, stale successor refusal, frozen closure, and one-use application.
- DONE: AC-5 — Specify proof of exact retained Result and inventory byte identity.
  At exact tip `98ebb458`, `TestGateRecordConsumesDirectBindingResultFromPreparedRoom` freezes both raw digests, makes `gate validate` reject a one-byte change to either fixed provider file with its path and `frozen digest`, and passes after exact restoration. The planned deletion rows require only the missing fixed path plus the read/not-found error before restoration.
- DONE: AC-6 — Prove zero provider coupling in the binary.
  At exact tip `98ebb458`, `TestGateReviewVerbIsAbsentAndSideEffectFree` proves the presentation verb is absent without side effects, and `go list -deps ./cmd/spacedock` returns zero Subspace dependencies; either regression makes the measured count non-zero.

### Summary

Reset the binding design to the Captain-approved prepared-room recorder boundary and folded Cycle 11's three test-only falsifiability gaps into its existing CLI matrices. Cycle 12 keeps the deletion proof honest: missing files assert their fixed path and read/not-found error, while only present-but-changed files assert `frozen digest`. Candidate `98ebb458` remains untouched at 17 files/1,255 changed lines; the next implementation pass is limited to an exact estimated `+42/-0` in `internal/cli/gate_test.go`, projecting 1,297 under the 1,300 cap. No production behavior, interface, provider transport, harness, fixture file, or documentation edit is proposed.

## Stage Report: implementation (cycle 3)

- DONE: Add only the three Cycle 11 proof variants to the existing CLI matrices.
  Commit `944e928e` changes only `internal/cli/gate_test.go`, adding `+38/-0`: a post-bind canonical Briefing-question mutation; same-cardinality duplicate-id and wrong-revision presented inventories; and deletion of each retained Result/inventory file after close.
- DONE: Prove every rejected room-record variant remains failure-atomic.
  The post-bind Briefing edit rejects with `frozen digest`; both same-cardinality inventory corruptions reject with `complete presentation mapping`; all three reuse the existing byte-identity assertions and leave the entity unchanged.
- DONE: Preserve the supported retained-file diagnostic split.
  Present-but-byte-changed Result and inventory files make `gate validate` name the fixed path plus `frozen digest`; deleting either file makes validation name the fixed path plus its read/not-found error. Exact restoration makes validation pass again.
- DONE: Stay inside the authorized surface and ownership boundary.
  Final diff from `main` is 17 files and `+1124/-169` = 1,293 changed lines, below the Cycle 9 hard cap of 1,300. No production, documentation, fixture-file, provider, interface, harness, or diagnostic bytes changed; executable provider transport remains Subspace-owned.
- DONE: Run focused, repository, race, documentation, formatting, dependency, and cleanliness checks.
  PASS: focused CLI tests, `go test ./...`, `go test ./... -race`, `uv run --with-requirements docs/requirements.txt mkdocs build --strict`, `gofmt -w ./cmd ./internal`, `git diff --check`, clean code worktree, and zero Subspace entries in `go list -deps ./cmd/spacedock`.
- DONE: Request final-tip adversarial review and canonically record its triage.
  Roborev job 2039 reviewed exact tip `944e928e`. Canonical implementation round 13 records 0 material findings fixed and three correct-but-disproportionate declines with promotion conditions: global retained-file mutation preconditions promote if the contract makes retained files an application precondition or archival needs a lifecycle rule; an in-repo provider smoke promotes if executable orchestration ownership moves here; wording polish promotes on demonstrated user confusion or a production reset of that error path.
- DONE: Leave a clean, reviewable exact-tip candidate without changing workflow control state.
  The implementation branch is clean at `944e928e`; this append-only report changes only the entity body and leaves frontmatter, gates, and application state to the first officer.

### Summary

The Cycle 11 falsifiability gaps are closed with 38 test-only lines at exact tip `944e928e`. The candidate now explicitly proves post-bind Briefing immutability, content-complete inventory matching despite equal cardinality, and distinct changed-versus-missing retained-file diagnostics, while preserving byte-clean rejection, the Subspace transport boundary, and the 1,300-line cap.
