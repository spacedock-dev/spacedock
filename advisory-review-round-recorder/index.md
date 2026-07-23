---
id: frze3yqm9da0vp0r53qqdc8t
title: Extend 3k's recorder to persist advisory review rounds
status: implementation
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
        gate: gate:docs-dev:fr:implementation
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
