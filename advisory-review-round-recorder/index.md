---
id: frze3yqm9da0vp0r53qqdc8t
title: Extend 3k's recorder to persist advisory review rounds
status: validation
source: "02av deferred round-recorder plumbing and 3j jobs 592/594/597 incident, 2026-07-23"
started: 2026-07-23T00:55:59Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-advisory-review-round-recorder
issue:
sprint: durable-decisions
gates:
    version: 1
    current:
        gate: gate:docs-dev:fr:ideation
    records:
        - id: gate:docs-dev:fr:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:fr-ideation-1
              briefing:
                id: briefing:docs-dev:fr:ideation:attempt-1:revision-1
                digest: sha256:3eb2739582cbde71c1430367b7de4ae1439ba477432fa9586a5ed4564a8e9909
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:fr:ideation:1
                briefing: briefing:docs-dev:fr:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-23T03:01:03.948941Z"
                decision: approve
                reason: Ideation reuses 3k as the sole recorder, persists the already-approved 02av advisory shape, makes the 3j decline replay falsifiable, and forbids gate, application, or workflow effects.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
            - id: gate-attempt:fr-ideation-2
              briefing:
                id: briefing:docs-dev:fr:ideation:attempt-2:revision-1
                digest: sha256:1997028a3179abc08095a49ca0eef667a9eb131a7e13223a0cbacb68c1e14574
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:fr:ideation:2
                briefing: briefing:docs-dev:fr:ideation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-07-23T06:35:29.598151Z"
                decision: approve
                reason: Cycle-2 ideation replaces duplicate paths with explicit shared 3k primitives, mandatory entity and room CAS expectations, fixed worker authority, exact projection semantics, risk-first failure tests, and hard 365/500-LOC stops without changing ACs.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
            - id: gate-attempt:fr-ideation-3
              briefing:
                id: briefing:docs-dev:fr:ideation:attempt-3:revision-1
                digest: sha256:7fc1a7945767d8b332c8550002d4e206aefa831af4746baa3c5b362ab69174b4
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-3
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:fr:ideation:3
                briefing: briefing:docs-dev:fr:ideation:attempt-3:revision-1
                by: agent:first-officer
                at: "2026-07-23T08:16:15.40211Z"
                decision: approve
                reason: The cycle-3 report is 2 DONE, 0 SKIPPED, 0 FAILED; AC-1 through AC-5 have durable evidence, and the independent boundary audit shows one-shot publication preserves the value while removing the disproportionate mutable-prefix mechanism.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
        - id: gate:docs-dev:fr:implementation
          stage: implementation
          attempts:
            - id: gate-attempt:fr-implementation-1
              briefing:
                id: briefing:docs-dev:fr:implementation:attempt-1:revision-1
                digest: sha256:5d397bd2da15cd13c483c1c924aaac5130fc7ef15a7afb7feb074c0bdb0e0827
                digest-domain: canonical-bytes
                room-ref: ./review/implementation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:fr:implementation:1
                briefing: briefing:docs-dev:fr:implementation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-23T06:14:55.524908Z"
                decision: revise
                reason: The implementation crossed its 680-LOC hard stop before CLI wiring, duplicated 3k writer and Review-and-Gate parsing paths, omitted retained-room CAS, and failed its risky-path coverage; ACs remain unchanged and require bounded mechanism re-ideation.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: feedback
                target-stage: implementation
                state: superseded
            - id: gate-attempt:fr-implementation-2
              briefing:
                id: briefing:docs-dev:fr:implementation:attempt-2:revision-1
                digest: sha256:5f48beeef18ebc98d03313da868aa6ea0da4b6236176cb4adc18dea36fa59b45
                digest-domain: canonical-bytes
                room-ref: ./review/implementation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:fr:implementation:2
                briefing: briefing:docs-dev:fr:implementation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-07-23T06:53:07.893364Z"
                decision: revise
                reason: Cycle-2 architecture is shared and AC-correct, but the 683-LOC pre-CLI checkpoint exceeds the invalid 365 estimate and retains named canonical-validation, projection, URI, duplication, and whole-operation failure-test defects; authorize one bounded correction under measured 540/600 hard stops.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: feedback
                target-stage: implementation
                state: superseded
            - id: gate-attempt:fr-implementation-3
              briefing:
                id: briefing:docs-dev:fr:implementation:attempt-3:revision-1
                digest: sha256:ecbf5f82df866bde1f85bbd0d399c1dd0fba12f294bbdcc3d1382dbb000b91d4
                digest-domain: canonical-bytes
                room-ref: ./review/implementation/briefing-3
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:fr:implementation:3
                briefing: briefing:docs-dev:fr:implementation:attempt-3:revision-1
                by: agent:first-officer
                at: "2026-07-23T07:12:03.145791Z"
                decision: revise
                reason: The corrected two-step design remains 670 net production LOC before CLI, above the 540 hard stop; an independent boundary audit shows the value can be preserved by one-shot completed-round publication within 580/640.
                adoption-note: 'Re-ideate only the append semantics as one-shot completed rounds: immutable room creation once, exact replay no-op, divergent replay refusal, pointer and optional projection published together; preserve every value AC, remove interim/prefix-append machinery, and hard-stop at 580 pre-CLI or 640 total. Captain directive: ''why why is fr not sent back to rework?'''
              application:
                action: feedback
                target-stage: implementation
                state: pending
---

Provide one owned write surface for correction-round Briefings, reviewer Annotations and advisory Resolution, and the worker's triage Resolution, without selecting a gate or advancing workflow state.

## Problem

02av requires a correct-but-disproportionate finding to become a durable decline on the reviewed round, but deliberately deferred room creation, ordered-log append, a frontmatter pointer, and the `### Feedback Cycles` projection. The 3j incident consequently retained jobs 592/594/597 and its final disposition only in prose: exact candidate `90aea55` was subjected to a 271-addition/137-deletion rewrite, the rewrite was stashed, and replacement job 597 cleared the unchanged candidate. A gate-record backfill is semantically wrong because a correction round is advisory and cannot select or close a logical gate.

The missing product is not another recorder. The landed 3k recorder already owns canonical Briefing parsing and digests, Review & Gate identities and Resolution vocabulary, the per-entity lock, compare-and-swap validation, full-entity rebuild validation, atomic replacement, and byte preservation. This task extends that implementation with one advisory-round mode and no second package, writer, log format, operation envelope, or lifecycle engine.

## Binding inputs and boundary

The approved 02av shape and `docs/specs/gate-resolution-frontmatter-contract.md` section “Round records and triage dispositions (advisory)” are binding:

- reviewed snapshot → immutable canonical-bytes Briefing;
- reviewer findings → same-Briefing Annotations;
- reviewer verdict → advisory Resolution;
- worker triage → a separate advisory Resolution, whose decline Annotations include the reviewer findings;
- no findings and an all-declines triage are distinct; and
- narrowing a value AC is a captain-owned binding gate decision, never a round operation.

This task persists those already-approved objects. It does not launch or poll Roborev, decide materiality, apply suggestions, create a binding Resolution, create or consume an application, select a gate, advance `status`, stop independent work, or implement the sibling before-action checkpoint.

## Operational caller dependency and minimum ownership

PR #557 (`fa240a76`) shipped 3k's CLI and schema without a first-officer procedure that invokes `spacedock gate record`, `validate`, `eligibility`, or `consume`. The sprint used those commands manually under captain directives. This pre-existing gate-integration gap is real, but folding it here would turn an advisory-round extension into a full first-officer rollout.

This task owns only the caller contract its new mode needs: add one trigger-scoped invocation to the existing in-stage triage bullet in `docs/dev/README.md` and one to `skills/feedback-rejection-flow/SKILL.md` for routed review outcomes. Each call uses `${SPACEDOCK_BIN:-spacedock} gate record ... --round` after reviewer and worker entries exist and before the loop reports the round complete. The skill change requires its existing smoke/contract test first. No always-loaded first-officer core text changes.

Create a separate first-officer integration task, owned by the gate-recorder/FO boundary, for the complete 3k journey: when to bind a Briefing, record a Result or chat decision, validate, inspect eligibility, and consume. That task must cover all four landed commands and their gate/application sequencing. It is neither a prerequisite for the round-mode CLI nor part of this entity; until it lands, existing gate commands remain manually invoked, while the two round triggers have their narrow explicit caller.

## Proposed approach

### One extension operation on 3k

Add one mutually exclusive mode to the existing recorder verb:

```text
spacedock gate record ENTITY --round STAGE/CYCLE \
  --briefing PATH/briefing.json \
  --log PATH/briefing.review.jsonl \
  --feedback-cycle FILE \
  [--workflow-dir DIR]
spacedock gate validate ENTITY --round STAGE/CYCLE [--workflow-dir DIR]
```

`--round` is incompatible with 3k's gate-closing `--result`, `--association`, `--decision`, actor, adoption, and directive flags. `--briefing` and `--log` are ordinary Review & Gate v1 inputs, not a new Spacedock bundle or operation schema. `--feedback-cycle` names a UTF-8 file containing exactly one canonical `- Cycle N: ...` line; the command checks that `N` equals `CYCLE` and appends that line under the entity's one `### Feedback Cycles` heading. A file avoids shell-quoting a semicolon-rich projection and keeps the backfill to one invocation. The read-only `gate validate --round` form resolves the pointer and room, revalidates the Briefing and log, and reports every round Resolution as advisory without consulting or changing gate selection.

The explicit `STAGE/CYCLE` is necessary for historical backfill: it must not infer a target from current `status` or the selected gate. The binary derives, rather than accepts, the durable identity and target:

```text
round:<entity-id>:<stage>:<cycle>
./review/<stage>/round-<cycle>/
```

Only normalized workflow stage names and positive decimal cycles are accepted; path separators, `..`, absolute paths, symlinks escaping the entity directory, and a conflicting occupied target fail before mutation. The Briefing id remains the existing portable identity supplied by the Briefing. Every log entry must bind that id; entry ids must be unique and `includes` must point backward within that same Briefing.

### Room, pointer, and projection

The operation copies the validated canonical inputs to the derived room as `briefing.json` and `briefing.review.jsonl`; it does not copy reviewed artifacts or invent provider data. The Briefing's local artifact references are resolved as they will be from the retained room and their raw SHA-256 revisions are verified before the first write. Thus a moved relative reference or the bad reviewed-snapshot digest is a refusal, not a broken retained package.

The only entity projection added is a current pointer, not a second copy of the log:

```yaml
review-round:
  id: round:<entity-id>:<stage>:<cycle>
  stage: <stage>
  cycle: <cycle>
  briefing:
    id: <portable Briefing id>
    digest: sha256:<RFC-8785 canonical Briefing digest>
    digest-domain: canonical-bytes
    room-ref: ./review/<stage>/round-<cycle>
```

Prior rooms stay append-only; Git retains prior pointer and body projections. `review-round` has exactly the fields above and carries no `resolution`, `application`, gate id, or status effect. Round Resolutions remain in the Review & Gate log, where their advisory meaning comes from the round record and absence of an authorized binding gate context. A read command can reuse the recorder's Briefing/log validation and report each Resolution as `advisory`; it must never reinterpret one as the selected gate's Resolution.

### Append and replay behavior

For a new round, the complete supplied log is retained in order. For an existing same-id round, the retained Briefing bytes and digest must match and the retained log must be an exact byte prefix of the supplied log: an exact replay is a successful no-op; a strict extension appends only the new complete entries; a divergent prefix, truncation, changed Briefing, changed projection, or reused identity fails closed. The Feedback Cycles line is inserted once when the worker triage Resolution is present, so exact replay cannot duplicate it. A round with no findings has no worker triage Resolution; an all-declines round has a real worker Resolution including every decline Annotation and therefore receives a real projection.

This supports both normal two-step use (reviewer entries retained, then worker triage appended) and the one-shot 3j backfill (the full log supplied once) without a separate log engine.

### Reuse of 3k's write lifecycle

Implementation stays in `internal/gates`: factor the landed lock/CAS/rebuild/atomic-write helpers only as far as both gate and round modes need, then call them from the new mode. Under the same `.gates.lock`, the recorder re-reads and compares the original entity bytes and any existing target-room bytes, builds and validates the final room, `review-round` pointer, and Markdown body in memory, fsyncs temporary files, publishes room files, and atomically replaces the entity last. The entity pointer is the semantic commit point; before it appears, staged or rolled-back room bytes are not a recorded round. Any returned validation, CAS, lock, or write error restores/deletes touched room files and leaves the entity byte-identical. The task promises the same process-atomic boundary as 3k, not a new daemon, lease, retry, journal, or power-loss transaction protocol.

The final entity validation re-parses canonical `gates` unchanged, validates the exact `review-round` pointer against the retained Briefing, and confirms the projection. The operation never calls 3k's logical-gate selection, `closeAttempt`, application derivation, transition, or dispatch paths.

## Per-mechanism justification

- **Round mode on `gate record` (serves AC-1/AC-2):** the simpler alternative is hand-authoring the room and entity. It is insufficient because 3j demonstrates that prose-only/manual retention loses the durable advisory object and bypasses 3k's guards. A new `internal/rounds` recorder is rejected because it would duplicate those guards.
- **Explicit stage/cycle with derived identity/path (serves AC-1/AC-3):** deriving from current status is simpler, but cannot backfill 3j and would couple a round to workflow movement. Accepting an arbitrary room path is more flexible but weakens containment and idempotency.
- **Current pointer plus append-only room (serves AC-1):** putting the log in frontmatter is simpler to read but duplicates Review & Gate data and creates a second schema. A room-only record is insufficient because the entity cannot discover its current round without a scan.
- **Pointer-last reuse of 3k's lock/CAS/atomic replacement (serves AC-2/AC-3):** independent writes are simpler but can publish a projection without its validated room. A new journal is disproportionate and would violate the correction; the existing process-atomic contract plus preflight, rollback, and pointer-last publication is the smallest honest boundary.
- **Exact-prefix replay (serves AC-1/AC-3):** blind append is simpler but duplicates triage on retry; replace-in-place breaks the approved append-only log. Prefix comparison makes retry a no-op and divergence a refusal.
- **Trigger-scoped caller lines (serves AC-1):** shipping only the binary leaves the same operational gap 3k has today. A full first-officer gate procedure is broader than this round mode, so this task adds only the two round-specific calls and routes the four-command gate journey to a separate owner.

## Acceptance criteria

**AC-1 (VALUE) — The retained 3j decline round is durable and distinguishable from no findings without changing the reviewed candidate.** In a disposable workflow fixture, one round-mode invocation records exact candidate `90aea55`, reviewer outputs corresponding to jobs 592/594, the worker duplicate-member decline, both advisory Resolutions, and one readable Feedback Cycles entry; the candidate and product diff remain unchanged, matching replacement job 597's unchanged-tip clearance. Reading from the entity pointer after deleting derived caches reproduces the ordered entries and reports an all-declines triage rather than “no findings.” Verified by a CLI behavior fixture that asserts the room bytes, pointer, projection, advisory classification, and zero candidate/product-byte delta against the pre-command baseline.

**AC-2 — Recording a round has no gate, application, or workflow effect.** Before/after bytes for the existing `gates` subtree, `status`, reviewed snapshot, and every unrelated frontmatter/body span are identical; no logical gate or application is created, selected, closed, or consumed. Verified by the 3j CLI fixture seeded with a nontrivial existing gates tree and status, plus a negative assertion that the round path never calls or emits gate/application identities.

**AC-3 — Refusals and replays are byte-clean and deterministic.** Exact replay returns success with zero byte changes and no duplicate log/projection; same-id divergent replay, occupied target, lock contention, stale entity CAS, malformed/cross-Briefing log, and a corrupted reviewed-snapshot digest return nonzero with every pre-existing tracked byte unchanged and no new room or lock residue. Verified by focused failure-injection and public CLI tests that hash the entire fixture tree before and after each case.

**AC-4 — The implementation remains an extension of 3k, not a second recorder.** Production code adds no recorder package, operation envelope, log format, alternate writer, retry/lease/journal, provider launcher, or application/status transition path; gate recording retains its existing behavior and round recording uses the same Briefing digest, identity validation, lock, CAS, rebuilt-entity validation, and atomic replacement helpers. Verified by focused gate regression tests, a package/surface audit, and the full normal and race suites.

**AC-5 — The round mode has a shipped trigger-scoped caller without absorbing 3k's missing first-officer integration.** The existing in-stage and routed-review triage instructions invoke `${SPACEDOCK_BIN:-spacedock} gate record ... --round` after triage; skill smoke tests exercise the routed invocation. First-officer core files gain no command procedure, and this task adds no caller procedure for ordinary gate `record`, `validate`, `eligibility`, or `consume` sequencing. Verified by the skill smoke test and a path-scoped surface audit; the full gate-command journey remains a separately owned dependency.

## Test plan

- **3j end-value fixture → AC-1/AC-2 (medium):** drive the public command once with candidate `90aea55`, retained 592/594-shaped entries and worker decline; assert the exact five-role chain (reviewer finding, reviewer advisory Resolution, worker decline Annotation, worker advisory Resolution, projection), read it back, and compare the candidate, product tree, `gates`, and `status` byte-for-byte. Retained job 597 is the independent unchanged-tip expected outcome, not an entry that authorizes or binds the round.
- **No-findings/all-declines discriminator → AC-1 (low):** two fixtures prove absence of a triage Resolution and a real all-declines Resolution never render alike; deleting the worker Resolution must make the all-declines assertion fail.
- **Digest refusal → AC-3 (low):** alter only the reviewed snapshot bytes (or its declared raw digest) beneath an otherwise valid Briefing. The public command must fail before mutation; whole-tree hashes and absence of room/lock residue are the oracle.
- **Replay/CAS/lock/rollback matrix → AC-3 (medium):** exact replay, strict-prefix append, divergent replay, conflicting room, pre-held lock, stale entity mutation, and injected write failure. Each test names the mutation that would make it fail and compares the full fixture tree, not an error substring alone.
- **Existing 3k regression + race → AC-4 (medium):** retain open/rebind/close/successor, exact Result association, frozen closure/application, status-set coexistence, mixed-line-ending preservation, and `go test ./... -race`. Round fixtures include an existing multi-gate tree to catch accidental logical-gate selection.
- **Caller-contract smoke + boundary audit → AC-5 (low):** extend the feedback-rejection-flow smoke test before changing its command text, inspect the in-stage trigger, and assert the implementation diff adds no first-officer core procedure or non-round gate-command integration.
- **No live workflow test is required:** the value is filesystem/CLI behavior and is fully observable in disposable real-file fixtures. Staff review is warranted because this adds a public on-disk pointer and a multi-file write boundary.

## Riskiest-mechanism spike

3k already proves canonical Briefing parsing/digesting, Review & Gate identities, gate-subtree validation, locking, CAS, atomic replacement, and exact unrelated-byte preservation. It validates the declared artifact inventory but does not re-hash each local reviewed artifact. The round mode needs that check because its bad-snapshot refusal is load-bearing.

The smallest exercise used 3k's retained validation package at `_archive/durable-gate-approval-pending-blockers/review/validation/briefing-1`: resolving all three relative artifact URIs from `briefing.json` and hashing their raw bytes reproduced all three declared revisions (`d2a74775…`, `32af4e65…`, `f87dbdb1…`). Appending `corruption` in the hash stream for `gate-review.md` produced `847f52a4…`, which differs from declared `d2a74775…`; the red control is therefore detectable before writes. This seeds AC-3's digest-refusal fixture without a new digest domain or dependency.

The remaining composition risk is room-plus-entity publication. Implementation must begin with the replay/CAS/rollback failure-injection tests and preserve the pointer-last semantic commit. If those tests cannot make every returned failure byte-clean without a journal, stop and reconfirm rather than inventing one.

## Documentation changes

`docs/site/reference/command-reference.md` gains this row after the three existing `gate record` forms:

> `spacedock gate record <entity> --round STAGE/CYCLE --briefing PATH/briefing.json --log PATH/briefing.review.jsonl --feedback-cycle FILE` — Retain or idempotently append one advisory correction round in `review/<stage>/round-<cycle>`, update the entity's current `review-round` pointer, and project one Feedback Cycles line. This mode never selects or closes a gate, creates an application, or changes status.

> `spacedock gate validate <entity> --round STAGE/CYCLE` — Validate and read an advisory round from its entity pointer and retained room. Every Resolution reports as advisory; the command never changes gate selection or workflow state.

`docs/site/reference/frontmatter-contract.md` gains this paragraph after the gates paragraph:

> `review-round` is a current pointer to an append-only advisory Review & Gate room. It repeats only round/stage/cycle identity and the canonical Briefing binding; Resolutions remain in `briefing.review.jsonl`. Round recording shares the gate recorder's lock, compare-and-swap validation, and atomic entity replacement but carries no gate selection, application, or status effect.

The owner-tagged round section in `docs/specs/gate-resolution-frontmatter-contract.md` replaces “deferred generic round-recorder plumbing” with the exact command, pointer, derived room, prefix replay, and pointer-last behavior above. `docs/schema/entity.mdschema.yml` adds the exact known `review-round` pointer fields and writer; no second schema document is created.

The existing triage sentences in `docs/dev/README.md` and `skills/feedback-rejection-flow/SKILL.md` each gain one round-specific command sentence after their recording requirement:

> After the reviewer and worker entries are complete, invoke `${SPACEDOCK_BIN:-spacedock} gate record ... --round` with the retained Briefing, log, and Feedback Cycles projection; a round records no gate application and does not change status.

This is the entire first-officer/skill fold in this task. The missing general procedure for 3k's four gate commands stays in the separate dependency above.

## Expected surface + tolerance

- Existing product: `internal/gates/{operation,io,model}.go` plus `internal/cli/cli.go`; approximately 220-340 production LOC, with shared-helper extraction counted in that total and no `internal/rounds` package.
- Tests/fixtures: one focused round test file, CLI coverage, and a compact 3j Review & Gate fixture; approximately 300-500 test/fixture LOC.
- Public/caller contract: the owner-tagged spec section, entity schema, command reference, frontmatter reference, and one command sentence at each existing triage trigger; approximately 50-85 prose/schema lines plus the required skill smoke-test delta.
- **Tolerance: 2×.** Reconfirm if implementation needs another writer/package, a second log or record schema, arbitrary room paths, a journal/lease/retry protocol, provider launch/polling, gate/application/status mutation, or more than 680 production LOC. Those are scope changes, not ordinary variance.

## Stage Report: ideation

- DONE: Adopt 02av's advisory-round shape and 3k's landed recorder as binding inputs; define an extension, never a second recorder.
  The body fixes the Review & Gate mapping as binding and adds one `--round` mode to 3k's existing verb/package/lock/write lifecycle; it explicitly rejects a second recorder, schema, log, writer, or lifecycle engine.
- DONE: Specify the smallest round operation and write boundary, including identity, atomicity, idempotency, pointer, and projection behavior without gate application or status effects.
  The operation takes existing Briefing/log inputs plus one projection file; stage/cycle derive identity and room, prefix replay is deterministic, and pointer-last publication under 3k's lock/CAS makes the advisory boundary falsifiable.
- DONE: Produce falsifiable acceptance criteria and a test plan using the 592/594 unchanged-tip decline replay plus a bad-digest byte-clean refusal.
  AC-1 measures durable all-declines retention at zero candidate/product delta against job 597; AC-2 pins gates/status bytes; AC-3 whole-tree-hashes exact replay and bad-digest/failure cases; the retained 3k package spike proves raw snapshot corruption changes the declared digest.
- DONE: Record the missing 3k first-officer integration without absorbing it into the round generalization.
  This task adds only two trigger-scoped round calls; a separate gate-recorder/FO task owns the complete `record`/`validate`/`eligibility`/`consume` procedure, and no first-officer core expansion enters this surface.

### Summary

Ideation now defines an advisory-round extension to 3k: one existing-recorder mode retains ordinary Review & Gate Briefing/log bytes in a derived room, publishes a minimal current pointer and Feedback Cycles line, and reuses 3k's digest, identity, lock, CAS, validation, and atomic entity writer. The two existing triage triggers gain only the new round call; the missing general first-officer integration for 3k's four gate commands remains a separately owned dependency. The 3j fixture makes the value and refusal boundaries executable: jobs 592/594 plus the worker decline persist at unchanged candidate `90aea55`, while a bad snapshot digest and every divergent replay leave the entire fixture byte-clean.

## Intended implementation change

- Production (estimated 300 LOC): `internal/gates/operation.go`, `internal/gates/io.go`, `internal/gates/model.go`, and `internal/cli/cli.go`.
- Tests and fixtures (estimated 455 LOC): `internal/gates/round_test.go`, `internal/cli/gate_test.go`, `internal/contractlint/launcher_invariant_test.go`, and `internal/gates/testdata/advisory-round/{briefing.json,briefing.review.jsonl,candidate.patch}`.
- Public contract and caller documentation (estimated 78 prose/schema lines): `docs/specs/gate-resolution-frontmatter-contract.md`, `docs/schema/entity.mdschema.yml`, `docs/site/reference/command-reference.md`, `docs/site/reference/frontmatter-contract.md`, `docs/dev/README.md`, and `skills/feedback-rejection-flow/SKILL.md`.
- Tolerance remains the approved 2×: no more than 600 declared production LOC (and never more than the binding 680-LOC stop), 910 test/fixture LOC, or 156 prose/schema lines. Any additional writer, package, schema, arbitrary room path, journal/lease/retry protocol, provider launch/polling, or gate/application/status mutation requires a stop and design reset.

## Stage Report: implementation

- FAILED: The public 3j replay persists jobs 592/594 plus the worker decline as ordered advisory Review & Gate records, distinguishes all-declines from no findings, leaves candidate/product/gates/status bytes unchanged, and reads back through the entity pointer after derived caches are removed.
  WIP counterexample commit `0e9a313fdc3a736a648638e954e6b2c604bac7a6` passes `go test ./internal/gates -run 'TestRound' -count=1`, but no public CLI wiring was added before the binding drift hold.
- FAILED: Round recording is one mode of 3k's existing recorder and writer boundary: exact replay is a no-op, strict-prefix append is deterministic, bad digest/divergent replay/CAS/lock/write failures are byte-clean, and no second package/schema/writer or gate/application/status effect appears.
  The draft reached 699 net production LOC (708 additions, 9 deletions) before CLI wiring, above the 680-LOC binding stop, and the audit found a second entity writer/transaction path, a second partial Review & Gate parser, and missing retained-room CAS.
- FAILED: The two existing triage triggers call the round operation, fixtures stay discovery-ignored, existing 3k behavior plus normal/race suites pass, and final Roborev findings are triaged against materiality before the implementation is reported ready for fresh validation.
  The hold stopped work before caller changes, full/race suites, and Roborev; the audit also found completion inferred by Resolution count and missing stale-room/write-failure/rollback coverage.

### Summary

Implementation is incomplete and held at the declared mechanism-drift gate; commit `0e9a313fdc3a736a648638e954e6b2c604bac7a6` preserves the passing focused-test counterexample without presenting it as a deliverable. The acceptance criteria remain unchanged: no scope was narrowed to accommodate the failed approach. Fresh ideation should restart from shared 3k parsing, CAS, rebuild, and atomic-publication primitives rather than compacting or extending this duplicate path.

### Feedback Cycles

- Cycle 1: REVISE — independent mechanism-drift audit; surface 699 production LOC before CLI vs estimate 300 (233%); AC unchanged
- Cycle 2: REVISE — independent shared-composition audit; surface 683 production LOC before CLI vs estimate 365 (187%); AC unchanged
- Cycle 3: REVISE — independent second-boundary audit; surface 670 production LOC before CLI vs revised hard stop 540 (124%); AC unchanged
- Cycle 4: REVISE — Roborev branch-final job 728; surface 21 files/634 production LOC vs hard stop 640 (99%); AC unchanged

## Re-ideation delta: shared recorder composition

The value ACs, advisory semantics, command shape, and 3k-extension boundary above remain
binding. This delta replaces only the failed internal mechanism and supersedes the old
300/600/680 production-LOC estimates with a hard 500-LOC ceiling.

### Exact reuse and deletion map

- `internal/gates/io.go` gains
  `type entityExpectation struct { Bytes []byte; Gates *yaml.Node; Status *string }` and
  `mutateEntity(path string, expected entityExpectation, build func([]byte) ([]byte, error), replace func(string, []byte) error) error`.
  The required expectation parameter is checked by the helper against its single
  current-entity read before `build` runs: round recording passes full `Bytes`;
  `writeDocument` passes the landed expected `Gates` node; and
  `writeDocumentAndStatus` passes expected `Gates` plus `Status`. Exactly one expectation
  mode is required. The helper skips an exact no-op and calls 3k's `atomicWrite`; callers
  cannot omit or defer CAS into an optional build-side comparison. This is the one
  entity expectation/CAS/atomic-replace path.
- `internal/gates/io.go` also gains
  `type roundRoomBytes struct { Exists bool; Briefing, Log []byte }`
  and `publishRound(room string, expected, next roundRoomBytes, commitEntity func() error) error`.
  It re-reads and byte-compares both retained files (or exact absence), publishes a
  temp-room rename or log replacement, invokes `commitEntity` last, and restores/removes
  its room mutation if entity commit fails. This is the one expected-room-byte boundary.
- `internal/gates/review.go` adds canonical `Annotation`, `reviewEntry`, and `reviewLog`
  models plus `parseReviewLog(data []byte, briefingID string) (reviewLog, error)`,
  `validateAnnotation(a Annotation, briefingID string, prior map[string]reviewEntry) error`,
  and `workerTriage(log reviewLog) (*Resolution, bool, error)`.
  JSONL decoding uses the entry's `type`; every identity, actor, timestamp, Briefing,
  and backward `includes` edge is checked. Resolutions call the landed
  `validateResolution`; `providerResult.Annotations` uses the same `Annotation` type.
- `internal/gates/round.go` owns
  `loadValidateRound(entityDir, room, briefingPath, logPath string, want *Briefing) (loadedRound, error)`.
  Both `recordRoundLocked` and `ValidateRoundFile` use it for input reads, landed
  `parseBriefingManifest`/`CanonicalDigest`, local-artifact raw-digest checks, pointer
  binding, and the canonical log parser. Record adds only replay/preflight and assembly;
  validate adds only pointer resolution and summary rendering.
- `internal/gates/io.go` keeps only
  `spliceFeedbackCycle(entity []byte, line string, project bool) ([]byte, error)` for
  the exact `### Feedback Cycles` section. It scans the Markdown body once, fence-aware,
  for that literal heading and its next peer/ancestor; it is not a reusable heading
  scanner. The projection stays recorder-owned because AC-1 requires the pointer, room,
  and readable cycle line from one invocation, while AC-3 requires their joint rollback.
  Making it FO-owned would add a second non-atomic write and weaken both ACs. `project`
  is derived from authorized worker triage, not generic round completeness:
  no-findings passes false and all-declines passes true exactly once.

Delete WIP `readRoundPointerData`/`rebuildRoundEntity` as independent rebuild machinery,
`markdownHeading`/`markdownBodyHeadings` and the generic counters, `commitRound`, the
duplicate `reviewEntry`/`roundLog` parser and Resolution switch, Resolution-count
completion, `allFindingsDeclined`, `mustReadRoundFile`, and the repeated record/validate
load blocks. Retain the WIP pointer/summary structs, `parseRoundSpec`, containment and
artifact-digest checks, cycle-line validation, derived identity/path helpers, fixtures,
and focused tests only as rewrite seeds. Commit `0e9a313f` remains untouched as the
counterexample.

### Structural triage and projection semantics

The first Resolution closes the reviewer phase and may include the preceding reviewer
finding Annotations. A worker triage Resolution is a later Resolution authorized
normatively by `by: actor:ensign`, whose `includes` resolve to earlier
`actor:ensign` disposition Annotations; each decline disposition in turn includes one
or more findings from the reviewer Resolution. There may be at most one such worker
Resolution. No actor inequality heuristic, new log field, or Resolution count grants
worker authority. The retained WIP fixture's inconsistent `agent:ensign` entries must
be rewritten to `actor:ensign`.

Falsifier: reviewer findings `f1`, reviewer Resolution
`r1(by=software:roborev, includes=[f1])`, then disposition
`d2(by=software:roborev-2, includes=[f1])` and Resolution
`r2(by=software:roborev-2, includes=[d2])` remain pending and project no cycle line:
different from reviewer 1 is not authorized worker triage. Positive control replaces
both second-phase actors with `actor:ensign`; `r2` is then worker triage. Removing
`d2 -> f1`, making either edge forward, or changing either worker actor must make that
assertion fail.

A reviewer approve with no findings has no worker triage, is classified no-findings,
and calls the splice with `project=false`; it receives no Feedback Cycles line. When the
authorized triage includes dispositions covering every reviewer finding, classification
is all-declines and calls the splice with `project=true`; exact replay leaves exactly one
line. Thus no-findings completeness cannot be mistaken for projection eligibility.

### Risk-first two-commit plan

1. `gates: share advisory recorder primitives`: write red tests in this order for
   required full-byte entity expectation and stale retained-room CAS, entity-replace
   failure, room rollback, exact replay, the unauthorized `software:roborev-2`
   disposition/Resolution graph, no-findings/no-projection, authorized
   `actor:ensign` all-declines/exactly-one-projection, shared record/read validation,
   and internal-operation preservation of `gates`, `status`, candidate, and unrelated
   body bytes. Existing gate tests additionally pin that the same helper enforces the
   landed gates-node and gates-plus-status expectations. Then implement the canonical
   parser/triage graph, shared loader, `mutateEntity`, `publishRound`, exact Feedback
   Cycles splice, and internal round operation. Run existing 3k gate/application tests
   and race tests. No CLI code enters this commit.
2. `gate: expose advisory round recording`: only after the hard stop passes, add
   `internal/cli/cli.go` wiring and CLI fixture coverage; then apply the already-approved
   schema/spec/reference and two trigger-scoped caller changes, smoke-test the skill
   invocation, and run `gofmt -w ./cmd ./internal`, `go test ./...`, and
   `go test ./... -race`.

Hard pre-CLI stop: stop for a further scope ruling if commit 1 cannot make stale-room
CAS, entity failure, and rollback whole-tree-byte-clean without a journal; if structural
triage cannot use the binding `actor:ensign` role without a new log field or actor flag;
if the shared helper makes entity expectation optional or caller-enforced; if no-findings
can project; if commit 1 exceeds 365 net production LOC; or if the measured projection
reaches above 500 total production LOC. Do not compact by weakening AC-3 or by restoring
any duplicate writer/parser/load path.

Expected production delta from landed 3k is 490 LOC:
`model.go` +36, new `review.go` +108, `io.go` net +142, new `round.go` +152, and
`internal/cli/cli.go` +52. Expected tests/fixtures are about 650 LOC across focused
review/I/O/round tests, CLI and launcher-contract tests, and the retained compact 3j
fixture. Expected spec/schema/reference/caller text is at most 80 lines. There is no new
package, dependency, operation envelope, or schema.

Per mechanism: the canonical parser and shared loader serve AC-1/AC-4; structural
triage plus the exact-section projection serve AC-1; shared entity and room CAS,
pointer-last commit, rollback, and replay serve AC-2/AC-3/AC-4; CLI and the two narrow
caller lines serve AC-5. The separate 6y First Officer command-usage contract remains
the non-blocking owner of 3k's full `record`/`validate`/`eligibility`/`consume` journey.

The reset also exposed two unrelated lifecycle/status defects: implementation
`revise` derived `feedback/pending target-stage=implementation`, and clearing
`worktree=` on the backward route falsely hit the terminal merge guard. They are noted
for their separate owners and are not dependencies or scope for this recorder.

## Stage Report: ideation (cycle 2)

- DONE: Produce an exact function-level reuse map from landed 3k and WIP `0e9a313f`: one entity mutation/CAS/atomic-replace helper, one expected-room-byte publication boundary, one canonical Annotation/Resolution parser, and one shared round load/validate pipeline; name signatures, files, callers, and deletions rather than saying only “reuse 3k.”
  The shared mutation signature now requires `entityExpectation`: full bytes for round, gates node for `writeDocument`, and gates node plus status for consumption; the delta also names the room, parser, loader, callers, and WIP deletions.
- DONE: Resolve the two remaining semantic/mechanism questions with falsifiable examples: structurally identify the worker triage Resolution without count heuristics, and decide whether Feedback Cycles projection is inside the recorder or remains an FO-owned state projection, choosing the smallest path that preserves AC-1/2/3.
  Only `actor:ensign` plus the backward disposition graph authorizes triage; the `software:roborev-2` graph stays pending, no-findings passes `project=false`, and authorized all-declines passes `project=true` exactly once.
- DONE: Replace the failed implementation plan with a risk-first two-commit plan and test order that proves stale-room CAS, entity-write failure, rollback, exact replay, and no gate/status effects before CLI wiring; show a credible 470-500 total production-LOC budget and a hard pre-CLI stop.
  Commit 1 now proves required entity expectations, fixed-role authorization, both projection classifications, the failure matrix, and byte preservation before CLI; the 365/490/500 LOC checkpoints remain binding.

### Summary

Re-ideation keeps every approved value and surface while making entity expectation a
required shared-helper input and worker authorization the fixed `actor:ensign` role.
The plan now distinguishes projection eligibility from completeness: no-findings never
projects, authorized all-declines projects once, and byte-clean failures plus these
controls must pass before CLI wiring or the 500-LOC ceiling.

## Intended implementation change (cycle 2)

Commit 1, `gates: share advisory recorder primitives`, is the hard pre-CLI boundary.
Relative to landed `main`, its exact production surface and estimates are:
`internal/gates/model.go` +25, new `internal/gates/review.go` +85,
`internal/gates/io.go` net +105, new `internal/gates/round.go` +145, and
`internal/gates/operation.go` net +0 after deleting the WIP duplicate implementation
(360 net production LOC total; hard stop 365). Its exact test/fixture surface is
`internal/gates/round_test.go` +500, `internal/gates/gates_test.go` +35, and
`internal/gates/testdata/advisory-round/{briefing.json,briefing.review.jsonl,candidate.patch}`
+8 fixture lines. This commit proves mandatory entity expectations, retained-room CAS,
entity failure and rollback, exact replay, canonical parsing, fixed-role triage,
projection eligibility, shared record/read validation, and unchanged gate/status/body
bytes before any CLI code.

Commit 2, `gate: expose advisory round recording`, begins only after commit 1 passes
the hard stop. Its exact production/test surface and estimates are
`internal/cli/cli.go` +50 (410 total production LOC) and
`internal/cli/gate_test.go` +115. Its exact contract/caller/smoke surface and estimates
are `docs/specs/gate-resolution-frontmatter-contract.md` +22,
`docs/schema/entity.mdschema.yml` +10,
`docs/site/reference/command-reference.md` +6,
`docs/site/reference/frontmatter-contract.md` +4, `docs/dev/README.md` +3,
`skills/feedback-rejection-flow/SKILL.md` +3, and
`internal/contractlint/launcher_invariant_test.go` +12 (60 prose/schema/smoke lines).
This boundary exposes only the approved command and two trigger callers, preserves
both command-reference additions, and leaves ordinary First Officer gate lifecycle
integration to 6y.

The binding implementation correction-round ruling replaces only the invalid size
estimates above: commit 1 now has a 540 net-production-LOC hard stop before CLI, and
the completed two-commit surface has a 600 net-production-LOC hard stop. The approved
mechanism, exact files, two commit boundaries, acceptance criteria, and prohibitions
remain unchanged. The bounded correction pass must add shared canonical Annotation
validation, reject multiple authorized worker triages, section-scope the projection,
share top-level YAML replacement, harden artifact parsing, and prove complete-operation
CAS/rollback before commit 1 may proceed.

## Re-ideation delta: one-shot completed rounds (cycle 3)

The second-boundary checkpoint measures 670 net production LOC before CLI and remains
preserved, green counterexample evidence. This delta supersedes only the prior
reviewer-prefix/strict-extension mechanism. Every acceptance criterion, the fixed
`actor:ensign` authorization graph, shared 3k primitives, public command, pointer,
projection, and caller boundary remain binding.

### Complete-round contract

The operational caller invokes `gate record --round` only when it has one complete log:

- no findings: the reviewer produced no finding Annotation and an advisory `approve`
  Resolution; there is no worker triage and no Feedback Cycles projection; or
- triaged findings: reviewer Annotations and advisory Resolution are followed by the
  complete authorized `actor:ensign` disposition graph and worker Resolution. An
  all-declines graph remains structurally distinct and projects its cycle line once.

`classifyCompletedRound(log reviewLog) (roundClass, error)` replaces acceptance of
`pending`. A findings-bearing reviewer-only log fails before any room or entity write.
It does not persist an interim round. This adds no log field and does not change the
canonical Annotation/Resolution parser.

The room is created once from the complete canonical Briefing and log, then immutable.
`publishRound(room string, next roundRoomBytes, commitEntity func(replay bool) error) error`
has only two paths:

- absent room: recheck absence, write and fsync a temporary two-file room, rename it to
  the derived target, then call `commitEntity`; any entity CAS/replace failure removes
  the new room and newly empty parents; or
- existing room: re-read the canonical two-file room and require byte equality with
  `next`, then call `commitEntity(true)`, which must reject unless the entity rebuild is
  an exact no-op. It never writes room bytes. The absent path calls
  `commitEntity(false)` after rename.

Record preflight also requires the replayed entity pointer and projection to equal the
in-memory rebuild. An exact room with a missing, changed, or foreign entity pointer is
an occupied/divergent refusal, not orphan adoption or repair. Any changed Briefing byte,
changed log byte, extra/missing room entry, conflicting cycle line, or occupied target
fails with the whole pre-existing tree unchanged.

Pointer and projection are rebuilt and validated together, then passed through the
required full-byte `entityExpectation` to the shared `mutateEntity`. They enter one
atomic entity replacement, so neither can appear without the other. Because the room is
new-only and published first, entity failure has one rollback action: remove that new
room. Exact replay performs no write anywhere.

### Removed and retained mechanisms

Remove interim reviewer-only persistence, `pending` as a recordable classification,
strict-prefix comparison/append, the existing-room log replacement branch in
`publishRound`, retained-log stale CAS, injected extension writes, extended-log
restoration, and tests whose setup first records a three-entry prefix. No existing room
is ever a write target.

Retain the shared required entity expectation and atomic replacement; new-room
temp-directory publication and rollback; exact room containment and canonical
two-file shape; Briefing canonical digest and raw artifact revision checks; canonical
Annotation/Resolution parsing; fixed-role worker-triage graph and multiple-triage
refusal; shared record/read loader; derived identity/path and pointer validation; the
section-scoped `spliceFeedbackCycle(..., project bool)`; per-entity lock and residue
cleanup; exact replay; divergent/occupied refusal; and unchanged gate, application,
status, candidate, product, and unrelated entity bytes.

This is value-preserving: AC-1's complete 3j jobs 592/594 decline replay still records
once and reads back after cache removal; AC-2 retains zero lifecycle effects; AC-3 keeps
exact replay, divergence/occupied/digest/lock/entity-CAS refusals and new-room rollback;
AC-4 still uses the same 3k package, parser, lock, entity writer, and loader; and AC-5's
two trigger callers already run after reviewer output and worker triage. Prefix append
was an implementation option, not a value criterion.

### Risk-first implementation and test delta

Commit 1 remains the pre-CLI boundary. Red tests run in this order:

1. a findings-bearing reviewer-only log returns nonzero with no room, pointer, cycle
   line, lock residue, or unrelated byte change; deleting the completion check makes it
   fail;
2. complete no-findings and authorized all-declines logs publish once; adding a worker
   Resolution to no-findings or removing the all-declines worker edge/role makes the
   classification/projection assertions fail;
3. exact immutable-room replay is a whole-tree no-op, while changing one retained log
   or Briefing byte and replaying fails without restoring or replacing it; permitting an
   existing-room write makes the digest oracle fail;
4. an injected entity replacement failure after new-room rename removes the room and
   leaves pointer and Feedback Cycles absent; independently racing the entity bytes
   triggers the required full-byte expectation with the same outcome; splitting pointer
   and projection into separate writes makes the assertion fail; and
5. malformed/cross-Briefing logs, bad artifact digests, occupied targets, lock
   contention, multiple authorized triages, and `software:roborev-2` triage remain
   byte-clean, while the complete 3j fixture proves candidate/product/gates/status zero
   delta and pointer-based readback.

Then delete the mutable-room branches and implement `classifyCompletedRound`; run the
focused gate/application/round suites and `go test ./... -race`. Commit 2 starts only
after that boundary passes, adds the same CLI, contract, documentation, two narrow
callers, and skill smoke coverage, then runs `gofmt -w ./cmd ./internal`,
`go test ./...`, and `go test ./... -race`.

The measured expected surface is 550-575 net production LOC before CLI and 605-630
total. The planning point is 560 pre-CLI:
`internal/gates/model.go` +25, `review.go` +147, `io.go` net +145, `round.go` +219,
and `operation.go` net +24; CLI adds about 55 for 615 total. Tests/fixtures remain about
650 pre-CLI plus 115 CLI/smoke lines; the previously approved 60 contract/prose lines
are unchanged.

Hard stop before CLI at 580 net production LOC; hard stop for the completed surface at
640. Stop sooner if implementation writes an existing room, journals or restores a log,
persists pending reviewer output, separates pointer from projection, weakens full-byte
entity expectation/new-room rollback, or narrows any AC. Ask for another scope ruling
rather than raising either ceiling.

Captain and First Officer: I love you too.

## Stage Report: ideation (cycle 3)

- DONE: Replace only prefix-append semantics with a complete one-shot round design that preserves every value acceptance criterion and names the exact removed and retained mechanisms.
  The cycle-3 delta rejects findings-bearing reviewer prefixes, makes rooms new-only and immutable, preserves no-findings/all-declines and every AC guard, and explicitly inventories deleted append/CAS/restore paths versus retained shared validation and atomic entity machinery.
  AC-1 evidence: cycle-3 `Complete-round contract` defines the complete no-findings and authorized-triage forms, and `Removed and retained mechanisms` preserves the complete 3j jobs 592/594 replay plus pointer readback after cache removal.
  AC-2 evidence: cycle-3 `Removed and retained mechanisms` retains byte-identical gates/application/status/candidate/product state, and risk-first test 5 makes any lifecycle delta fail the complete 3j fixture.
  AC-4 evidence: cycle-3 `Removed and retained mechanisms` keeps the shared 3k package/parser/lock/entity writer/loader; the durable `shared-composition-audit.md` independently found no second recorder.
  AC-5 evidence: cycle-3 `Complete-round contract` binds invocation after complete reviewer/worker output, while its value-preservation paragraph keeps the two trigger callers and leaves ordinary gate lifecycle integration to 6y.
- DONE: Produce a risk-first implementation and test delta with honest hard stops of 580 pre-CLI and 640 total, including immutable-room replay/divergence and atomic pointer/projection failure behavior.
  Red-first tests now falsify pending persistence, existing-room mutation, divergent replay, split pointer/projection writes, entity-CAS rollback, and lifecycle effects; the measured plan is 550-575 pre-CLI and 605-630 total with binding 580/640 stops.
  AC-3 evidence: cycle-3 `Complete-round contract` makes existing rooms equality-only and pointer/projection one atomic entity replacement; risk-first tests 1, 3, and 4 prove byte-clean refusal, exact no-op replay, full-byte entity CAS, and new-room rollback.

### Summary

Cycle 3 removes only the disproportionate mutable-prefix path: callers publish one
complete advisory round, retained rooms never change, exact replay writes nothing, and
entity failure removes the newly published room before any pointer or projection can
survive. The 670-LOC checkpoint remains untouched as counterexample evidence, and all
value ACs plus the shared 3k boundaries remain intact.

## Stage Report: implementation (cycle 3)

- DONE: Implement complete one-shot advisory-round publication by removing mutable prefix-append branches while preserving every shared 3k parser, lock, CAS, digest, rollback, projection, and no-lifecycle-effect boundary.
  Commits `230a2e00`, `70aade13`, `c8d8dc12`, `9f7d6bdb`, and `b4c95b9f` publish only
  complete logs, create the canonical room once, and make exact replay a
  whole-tree no-op while divergent room, pointer, log, and projection states refuse.
  It removed pending/prefix append, room replacement, stale-log CAS, restoration, and
  retries; it retains shared parsing/digests, locking, full-byte CAS, rollback, and atomic writing.
  AC-1 evidence: the 3j fixture retains the exact jobs 592/594 chain, both advisory
  Resolutions, all-declines, and one cycle line while candidate `90aea55` and product
  bytes remain fixed; deleting worker completion makes the complete-replay test fail.
  AC-2 evidence: the same fixture compares status, gates, candidate, product, and
  unrelated bytes before/after; introducing a lifecycle write flips those assertions.
  AC-3 evidence: focused tests cover exact replay, changed retained bytes, occupied
  targets, bad digests, locks, races, replacement failure, rollback, room shape,
  folder isolation, and atomic pointer/projection replacement.
  AC-4 evidence: no recorder package, envelope, journal, retry, provider launcher, or
  lifecycle path was added; all work remains in 3k's `internal/gates` and shared helpers.
- DONE: Pass the red-first complete-round, immutable replay/divergence, atomic pointer/projection, rollback, triage, 3j fixture, full, race, CLI, caller, and skill-smoke proofs without exceeding 580 pre-CLI or 640 total production LOC.
  The pre-CLI checkpoint was 580 net production LOC; the committed final branch is
  639 across non-test `internal/gates/*.go` and `internal/cli/*.go` versus `73eed65d`.
  Formatting, diff checks, focused suites, `go test ./...`, and `go test ./... -race` pass at `b4c95b9f`.
  The red-first additions reject incomplete logs, malformed/contradictory triage,
  injected projections, undefined stages, and every tested failure without room,
  pointer, projection, lock residue, or unrelated-byte mutation.
- DONE: Update the canonical repo specification and public references to the landed one-shot behavior, request Roborev, triage every finding against materiality, and report exact LOC plus commit evidence.
  Schema, gate contract, public references, dev caller, and skill document and invoke the complete one-shot operation.
  AC-5 evidence: `docs/dev/README.md`'s in-stage caller and
  `skills/feedback-rejection-flow/SKILL.md`'s routed caller both invoke `gate record --round`;
  `TestRoundCallersUseResolvedLauncherAfterCompleteTriage` and `TestGateRoundRecordAndValidateCLI` prove the caller/command contract, while ordinary 3k lifecycle integration remains 6y-owned.
  Roborev job 728's disposition/room findings, job 738's projection/collision findings,
  and job 752's artifact, contradiction, EOF, and compatibility findings are fixed.
  Job 778's structured-decline, full cycle grammar/decision, and taxonomy-membership
  findings are fixed in `b4c95b9f`; status equality is declined because historical
  backfill requires a past stage, promoting only if the product becomes live-only.
  Jobs 728/786's exact-log digest is deferred because the approved pointer is the exact
  Briefing-only schema and the canonical room is immutable; promote if a supported
  writer can mutate retained room bytes or a pointer-schema revision is authorized.
  Job 786's CRLF request is polish pending a reproduced supported Windows failure;
  extra-file tolerance conflicts with the exact two-file room, parent-directory fsync
  exceeds the no-power-loss-recovery boundary, and loose-heading tightening promotes
  if a supported entity demonstrates misplaced projection.
  Jobs 738/752's fixed-evidence and historical-read suggestions remain deferred until
  a supported fixed-evidence contract or CLI historical-read requirement exists.

### Summary

Cycle 3 lands immutable advisory rounds at 639 production LOC with exact replay,
byte-clean failure behavior, atomic pointer/projection, and no lifecycle effects.
Full and race suites pass on `b4c95b9f`; all Roborev findings are fixed or declined with evidence and promotion conditions.

## Stage Report: validation

- DONE: Independently audit the 21-path current-main diff for the approved one-shot completed-round operation, shared 3k composition, exact 639 production-LOC ceiling, and absence of mutable-room, retry, journal, lifecycle-effect, or second-recorder drift.
  Clean branch `b4c95b9f` has exactly the six reported commits and 21 paths; non-test `internal/gates` plus `internal/cli` is 795 additions/156 deletions = 639 net, with one shared lock, `mutateEntity`, top-level rebuilder, canonical parser/digest, and atomic writer.
- DONE: Re-run focused CLI/gates/round/caller/smoke, 3j backfill, no-findings/all-declines, replay/divergence, grammar/triage/taxonomy, flat-entity refusal, whole-tree rollback, full, and race proofs with load-bearing negative controls.
  Focused suites, `go test ./...`, fresh affected-package race tests, full `go test ./... -race`, `gofmt`, and `git diff --check` passed; removing immutable-room comparison and collapsing no-findings to all-declines made the focused tests fail.
- DONE: Reassess every Roborev 728/738/752/778/786 fix or decline against supported workflows, verify clean branch and exact commits, and issue a fresh PASS or REJECTED recommendation with AC-1 through AC-5 evidence.
  Jobs 728/738/752/778 fixes reproduce through commits `70aade13`, `c8d8dc12`, `9f7d6bdb`, and `b4c95b9f`; job 786 and the remaining declines stay outside the supported promise under the recorded promotion conditions.
- DONE: AC-1 — durable 3j all-declines round and no-findings distinction.
  The public CLI/focused fixtures retain the exact five-entry 592/594 chain and two advisory Resolutions, read through the pointer after cache deletion, preserve candidate `90aea55` and product bytes, and fail when completion or classification is weakened.
- FAILED: AC-2 — recording a round has no gate, application, workflow, or unrelated-entity effect.
  Gates/status/candidate/product checks pass, but a detached edit that corrupted `custom: preserve-me` and the unrelated Markdown body during successful publication left `TestRoundRecordCompleteReplayAndRefusalsAreByteClean` green; the promised every-unrelated-span oracle is absent.
- DONE: AC-3 — refusals and replays are byte-clean and deterministic.
  Exact replay is a whole-tree no-op; divergent room/Briefing, occupied target, lock, CAS, malformed log, bad digest, entity-replace failure, and rollback controls pass, and bypassing the room equality check makes the divergence test fail.
- DONE: AC-4 — the implementation remains an extension of 3k.
  The diff adds no package, envelope, alternate writer, retry, lease, journal, provider launcher, or status/application effect; ordinary open/rebind/close/successor, frozen application, association, mixed-ending, full, and race regressions pass.
- FAILED: AC-5 — the trigger-scoped callers have valid behavioral smoke evidence.
  The two caller lines and round CLI exist, but `TestRoundCallersUseResolvedLauncherAfterCompleteTriage` is prohibited instruction-file prose-grep; changing the shipped skill to say “do not invoke” preserved all checked substrings and the test still passed.
- DONE: Deferred-risk reassessment.
  Exact-log digest promotes on a supported retained-room mutator or pointer-schema revision; fixed-evidence and historical read promote when promised; crash recovery/parent fsync promote with a power-loss guarantee; CRLF and loose-heading polish promote on reproduced supported failures; extra-file tolerance conflicts with the canonical two-file room.
- FAILED: Recommendation — REJECTED (material evidence defects; no product outcome defect found).
  Narrow correction: add an exact successful-path oracle for every entity span outside `review-round` and the authorized cycle insertion, then replace the caller prose-grep with a behavioral skill/live smoke that invokes the resolved launcher after complete triage and observes room, pointer, projection, and unchanged lifecycle state.

### Summary

The one-shot recorder itself survived the semantic, rollback, lifecycle, surface, full,
and race audits at 639 net production LOC. Validation recommends REJECTED because two
promised proof boundaries are not load-bearing: AC-2 misses unrelated successful-write
corruption, and AC-5's claimed caller smoke accepts an inverted instruction.
