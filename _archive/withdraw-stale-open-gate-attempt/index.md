---
title: Withdraw a stale open gate attempt without fabricating a decision
status: done
source: "Observed by the Subspace Shaping FO on 2026-07-26 after a legitimate sprint re-scope left a frozen request-backed attempt open with no truthful exit."
started: 2026-07-26T10:50:55Z
completed: 2026-08-02T08:04:26Z
verdict: passed
score: 1.0
worktree: .worktrees/spacedock-ensign-withdraw-stale-open-gate-attempt
issue:
sprint: durable-decisions
id: 0m6vtrw4qh9w4x6bn06x5hen
gates:
    version: 1
    records:
        - id: gate:0m6vtrw4qh9w4x6bn06x5hen:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:0m6vtrw4qh9w4x6bn06x5hen-backlog-1
              briefing:
                id: briefing:0m6vtrw4qh9w4x6bn06x5hen:backlog:attempt-1:revision-1
                digest: sha256:9bfedeb38906e04bae528cedfdb96f101efaa1d63c819b44922a8ee6e5db60f6
                request-digest: sha256:0526bd61039ea579f7595d23e5ccfd8bd3d3f18ee7ce5b64b211456df37f8524
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:0m6vtrw4qh9w4x6bn06x5hen:backlog:1
                briefing: briefing:0m6vtrw4qh9w4x6bn06x5hen:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-26T10:48:05.853485Z"
                decision: approve
                reason: The observed supported re-scope cannot be represented truthfully today; ideation is authorized to define the minimum withdrawal semantics without weakening frozen room integrity.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:0m6vtrw4qh9w4x6bn06x5hen:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:0m6vtrw4qh9w4x6bn06x5hen-ideation-1
              briefing:
                id: briefing:0m6vtrw4qh9w4x6bn06x5hen:ideation:attempt-1:revision-1
                digest: sha256:8319f9b10709cca9e5ae3f6c547b1fd2db6abef666f020013f38bdf845da5c0d
                request-digest: sha256:d0bc5693618ca77595374a51ae615385e9c3cb0fd4f62f740be8263ccdb6c2a7
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:0m6vtrw4qh9w4x6bn06x5hen:ideation:1
                briefing: briefing:0m6vtrw4qh9w4x6bn06x5hen:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-26T12:16:38.514969Z"
                decision: approve
                reason: AC-1 through AC-3 are evidenced and staff corrections are closed. Record approval now; apply only after s4 lands so implementation rebases on the final prepare and retained-authority surface.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:0m6vtrw4qh9w4x6bn06x5hen:validation
          stage: validation
          attempts:
            - id: gate-attempt:0m6vtrw4qh9w4x6bn06x5hen-validation-1
              briefing:
                id: briefing:0m6vtrw4qh9w4x6bn06x5hen:validation:attempt-1:revision-1
                digest: sha256:2811a3f1ae4343575536d65c3cea66a0d644e97182cf1814476d23f9e0b9f61c
                request-digest: sha256:c43512b5f252776d536e89b7065dbcf5784b7fc96d7e84aa5e8fff63e5cd89af
                room-ref: ./review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:0m6vtrw4qh9w4x6bn06x5hen:validation:1
                briefing: briefing:0m6vtrw4qh9w4x6bn06x5hen:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-30T23:33:58.696832Z"
                decision: approve
                reason: 'Captain conn accepts candidate 5a8be3220: deterministic lifecycle and adversarial evidence prove truthful withdrawal, immutable attempt N, and provider-backed successor N+1; the skipped real-Claude repeat is externally quota-blocked before FO work and is not a product defect.'
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:0m6vtrw4qh9w4x6bn06x5hen-validation-2
              briefing:
                id: briefing:0m6vtrw4qh9w4x6bn06x5hen:validation:attempt-2:revision-1
                digest: sha256:658468c64c65eedc3f5f02a11f329084ef6acde73f071f9238f3078c44dcf88a
                request-digest: sha256:4d9578638c01f1b590fbfa3e2971c6b73b921cd729a49a7a43294b80bade36be
                room-ref: ./review/validation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:0m6vtrw4qh9w4x6bn06x5hen:validation:2
                briefing: briefing:0m6vtrw4qh9w4x6bn06x5hen:validation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-02T08:02:48.262099Z"
                decision: approve
                reason: Captain authorized merge on green and sprint advancement on Science advisory approval; fresh cycle-2 report, AC-1/2/3 scan, exact-head CI, and Science Officer APPROVE/PASSED all support this candidate. The prior 5a8be3220 approval is stale.
              application:
                target-stage: done
                state: consumed
mod-block:
pr: pr-merge:580
archived: 2026-08-02T08:04:26Z
---

A provider-neutral prepared attempt can be frozen open after its reviewed candidate
becomes stale. Today there is no truthful supported exit: recording `hold` fabricates a
Captain decision, while editing `gates:` bypasses the canonical writer. This design adds
one FO lifecycle operation that retires the stale attempt without a Resolution, keeps
its retained authority immutable, and lets the existing `gate prepare` path create the
ordinary successor.

This cycle-2 design is rebased on s4's implementation tip `e328ecc6`
(`prepare-provider-neutral-gate-room`, PR #570). It supersedes the cycle-1 command,
room-basename, replacement, and proof choices preserved in the first Stage Report
below. Implementation must recheck the s4 tip after it lands; semantic drift in
preparation, retained authority, or provider recording requires design re-entry.

## Problem and boundary

The operation applies only to the selected current-stage attempt when that attempt is
open, request-backed, prepared by `gate prepare`, and has no provider output. It adds
one attempt-local withdrawal marker. It writes no Resolution, provider evidence,
application, workflow status, replacement room, or dispatch.

This task adds no compatibility or migration, generic event/history model,
cancel/supersede vocabulary, lineage pointers, alternate writer, Captain/ensign
withdrawal, closed-attempt withdrawal, chat-only withdrawal, or automatic replacement.
The legacy Briefing recorder may remain for its existing request-less behavior, but it
is neither the stale-room recovery path nor a supported positive-path proof.

## Spike evidence and corrected proof boundary

The cycle-1 throwaway spike established only that an attempt-local third state can
round-trip, survive a cold boot, preserve attempt-N bytes, and allow a later attempt to
own the real decision. Its `gate record --withdraw`, `briefing.json`, and subsequent
`record --briefing` calls predated the binding s4 surface and are not acceptance
evidence for this design.

The supported behavior proof is now exactly:

```text
gate prepare N
state commit
gate withdraw N
state commit
cold boot
gate prepare N+1
state commit
provider close N+1 with gate record --room
state commit
gate consume N+1
state commit
```

The implementation test must drive that sequence through a freshly built real binary
on a folder-form split-root fixture. It must not substitute direct YAML mutation,
`record --briefing`, or chat `--decision` for any step.

## Governing design

### Public command and implicit authority

The public grammar is exactly:

```text
spacedock gate withdraw ENTITY --reason TEXT [--workflow-dir DIR]
```

`withdraw` is a sibling of `prepare`, `record`, `validate`, `eligibility`, and
`consume`, not a semantic source of `gate record`. It accepts only one nonblank UTF-8
reason and the ordinary workflow directory. There is no `--actor`, timestamp, attempt,
gate, successor, or room flag. Passing `--actor` or any record/prepare flag is a usage
error before mutation.

The binary always records `by: agent:first-officer`; that fixed attribution is the
command's authority contract, not caller-supplied provenance. It supplies the
RFC3339Nano UTC timestamp and derives the current gate and attempt under the existing
entity lock. Successful stdout is stable:

```text
withdrawn gate=<gate> attempt=<attempt> state=withdrawn briefing=<briefing>
```

### Minimum durable shape and explicit attempt states

Attempt N gains only this optional mapping:

```yaml
withdrawal:
  by: agent:first-officer
  at: "2026-07-26T11:30:00.123456Z"
  reason: Sprint re-scope replaced the reviewed candidate.
```

The field name and owning attempt supply type and identity, so no event id, state,
successor, or lineage field is added. `Withdrawal` validation requires the exact fixed
actor, a parseable timestamp, and a nonblank reason. A withdrawn attempt must also
retain a nonempty valid `request-digest`.

Every attempt is classified before lifecycle logic runs:

- open: both `withdrawal` and `resolution` are absent;
- withdrawn: `withdrawal` is present and `resolution`, `provider-evidence`, and
  `application` are absent;
- closed: `resolution` is present and `withdrawal` is absent, with all existing
  provider-evidence and application rules unchanged.

No other combination validates. `Resolution == nil` alone never means open.
Withdrawn and closed are the two terminal/retired attempt states; “terminal” here is
attempt-local and does not change the entity's workflow stage.

### One writer, retained authority, and byte-clean refusal

`gates.Withdraw` uses the same `lockEntity`, decoded old-node compare-and-swap,
`Validate`, `ValidateTransition`, `writeDocument`, and later path-scoped
`state commit` sequence as the other gate mutations. A new CLI subcommand does not
create a second state writer.

While holding that lock, it first calls s4's existing `validateRetainedAuthority` on
the complete document, then resolves the selected current-stage attempt. That
validator deliberately walks every request-backed historical attempt, reopens its
canonical `gate-briefing.json` and `request.json`, resolves retained Git-root source
objects, and—where provider evidence exists—revalidates the exact provider Result,
presented inventory, and durable Resolution equality. Withdrawal does not skip the
candidate attempt or weaken the provider-evidence branch.

The full validator is O(history), including Git-source reopening. That cost is accepted
for correctness because no measured fixture or live trace currently breaches a
declared latency bound. Caching, indexing, incremental validation, and pruning are
deferred until a separate task supplies both a reproducible real-binary latency
measurement and an explicit maximum for a stated attempt count.

After full retained-authority validation, withdrawal requires all of these:

1. entity status names a gate stage and resolves the same record as
   `gates.current.gate`;
2. the target is that record's last ordered attempt and its explicit state is open;
3. its `request-digest` is present and valid;
4. the retained prepared room still contains exactly s4's
   `gate-briefing.json` and `request.json`, with no provider output;
5. the trimmed reason is nonblank.

The exact two-entry room check prevents withdrawing an already-returned provider
decision merely because it has not yet been recorded. Such a room must follow the
existing provider recording path.

Only after every check succeeds does one atomic entity replacement add `withdrawal`.
The command never writes the room. Blank reason, unknown flags, wrong selection or
stage, chat-only/open-without-request, closed, already withdrawn, corrupt retained
authority, provider output, lock contention, validation failure, and CAS failure all
return nonzero with identical entity and room bytes and no residual lock file.

### Frozen history and unchanged provider-evidence rules

`ValidateTransition` freezes both retired states. A previously withdrawn attempt
cannot be edited, deleted, closed later, rebound, given provider evidence/application,
or changed when N+1 is appended. Its retained Briefing/request and Git-source
authority remain subject to every later `gate validate`, `gate prepare`,
`gate record`, `gate eligibility`, and `gate consume` read that already invokes the
full retained-authority validator.

The existing closed prepared-attempt rules do not change: if durable
`provider-evidence` exists, its Result and presented-inventory digests must verify and
the retained Result Resolution must exactly equal the durable Resolution. Withdrawal
creates no provider evidence and introduces no validation branch that permits provider
evidence without a Resolution.

### s4 preparation and mutation audit

The implementation must replace resolution-null shorthand with explicit state
predicates at every s4 preparation seam:

- `prepareTarget`: an open prepared attempt replays N; a withdrawn or closed attempt
  derives N+1.
- Pre-publish frozen-binding check: only an explicit open request-backed attempt may
  replay, and a changed binding still refuses.
- `preparedEntityReplaySource` and `preparedReplay`: only an explicit open prepared
  attempt is replayable; withdrawn N is never a replay source.
- Post-publication entity mutation: only open N can take an idempotent same binding;
  withdrawn or closed N appends the derived successor. A withdrawn attempt has no
  pending application to supersede.
- `publishPreparedRoom` retains its existing clean publish/rollback behavior and exact
  two-file room. A failure after publishing N+1 remains byte-clean under s4's current
  rollback/CAS contract and can never rewrite N's room.

The focused preparation tests must exercise both target selection and the
post-publication mutation branch. Merely unit-testing `prepareTarget` is insufficient.

### Recorder, application, summary, and readiness audit

All gate mutation guards classify the attempt explicitly:

- room-backed and chat recording accept only open attempts and refuse withdrawn ones
  before adding provider evidence, Resolution, or application;
- the retained legacy Briefing bind, if still present, may mutate only an explicit
  open attempt and must append after either retired state;
- application eligibility and consumption require the explicit closed state, so a
  withdrawn attempt remains read-only and ineligible;
- transition validation freezes withdrawn and closed attempts, while retaining the
  one existing pending-to-superseded application exception only for closed attempts.

`CurrentSummary` and `gate validate` report `state=withdrawn` with empty resolution,
decision, and application fields. `CurrentStageReadiness` reports
`withdrawn-awaiting-prepare`. `computeReadyGates` includes that state in the
identify-boot `ready_gates` array, yielding exactly one row for the entity:

```json
{"id":"<id>","slug":"<slug>","current":"<stage>","readiness":"withdrawn-awaiting-prepare"}
```

No dispatch becomes eligible and workflow status is unchanged. After `gate prepare`
appends N+1, summary/readiness return to `open` / `awaiting-captain` for N+1.

### First Officer cold-boot recovery

The FO capability preflight requires the `withdraw` subcommand and `--reason` flag.
When an ordinary presentation becomes stale before any provider decision, the FO runs
`gate withdraw`, commits the selected entity, and stops unless it already has the
inputs needed for a new preparation. On a fresh boot,
`withdrawn-awaiting-prepare` has exactly one recovery action: run `gate prepare` for
the derived successor, commit its entity plus new two-file room, present that emitted
room, and stop at the Captain boundary. It never records, consumes, presents, or
dispatches the withdrawn attempt, and never uses `record --briefing` for recovery.

## Mechanism choices

| Mechanism | Value AC | Simplest alternative | Why insufficient |
| --- | --- | --- | --- |
| Attempt-local `withdrawal` | AC-1, AC-2 | Record `hold` | `hold` asserts a Captain decision and creates an application. |
| Separate `gate withdraw` using the shared writer | AC-1, AC-2 | Add `--withdraw` to `gate record` | The binding public surface is a distinct FO lifecycle action; provenance is implicit, not a record source. |
| Existing `gate prepare` after withdrawal | AC-1, AC-3 | Rebind with `record --briefing` | s4 owns provider-neutral creation, target derivation, publication rollback, and retained Git-source authority. |
| Explicit open/withdrawn/closed classifier | AC-1, AC-2 | Test only `Resolution == nil` | Both open and withdrawn lack a Resolution but only open may replay or mutate. |
| Full retained-authority validation | AC-2 | Validate only N | Corrupt historical evidence or provider association must fail before any new mutation. |
| `withdrawn-awaiting-prepare` ready row | AC-3 | Omit the retired state | Cold boot would strand N or mistake it for a Captain-ready room. |

## Acceptance criteria

**AC-1 (VALUE)** A stale prepared attempt is retired and replaced with exactly zero
false Resolutions: attempt N has one FO-attributed withdrawal and no Resolution,
provider evidence, or application; attempt N+1 alone carries its actual room-backed
Captain Resolution and consumed application.

Test: the fresh-binary split-root fixture runs the exact supported sequence above.
It fails if withdrawal routes through `record`, writes caller-selected attribution,
uses chat decision recording, reuses N, changes workflow status, or places the
provider Resolution on N.

**AC-2 (INTEGRITY)** Withdrawn N is immutable, retains verifiable s4 authority, and
every refusal is byte-clean without weakening provider-evidence validation.

Test: before/after tree digests cover `index.md`, room N, room N+1, dirty sibling
state, and lock paths across blank reason, extra/`--actor` flags, stale selection,
chat-only attempt, closed attempt, repeat withdrawal, corrupt Briefing/request/Git
source, provider output, and lock contention. Historical validation after N+1 close
and consume must still detect mutations to N and must still detect Result,
presented-inventory, or Resolution mismatch on any closed provider attempt.

**AC-3 (RECOVERY)** Cold boot after withdrawal emits exactly one
`withdrawn-awaiting-prepare` ready row; preparation appends N+1 and returns the entity
to `awaiting-captain` without dispatch or status change.

Test: the real-binary lifecycle asserts exact boot JSON before and after replacement.
The registered live FO gate-stop lane starts from a durably withdrawn prepared attempt
and proves that the shipped FO uses `gate prepare`, commits N+1, presents the emitted
room, and stops open without `record --briefing`, decision, consume, or dispatch.

### Proof ownership

Model-read Markdown is not a behavioral oracle. This task deletes the committed s4
contractlint assertions that search the FO lifecycle or presenter prose for
`gate prepare`, emitted-room text, lifecycle headings, or prepare/record/consume
ordering. It adds no contractlint assertion for `gate withdraw`, `state commit`,
`gate prepare`, their arguments, output, or relative order.

Contractlint may retain only structural closure, skill/capability/reference wiring,
frontmatter validity, absence/dedup invariants, and byte caps. The fresh-binary
split-root lifecycle owns command behavior and durable order; the existing registered
live FO lane owns whether the model follows the recovery sequence. No prose substring
test may substitute for either proof.

## Behavior-first test plan

1. Add model and transition tests first in `internal/gates/gates_test.go`: the three
   exclusive states, fixed withdrawal attribution/time/reason, request-backed
   requirement, frozen retired history, summary, and readiness. Mutating the classifier
   back to `Resolution == nil` must make withdrawn replay/readiness tests fail.
2. Add s4 regression tests first in `internal/gates/prepare_test.go`: prepare N,
   install a valid withdrawal through the production writer seam, prepare N+1, and
   verify target derivation, replay exclusion, post-publish append, exact old-room
   bytes, rollback, and retained-authority refusal. A mutant that takes the old
   `previous.Resolution == nil` branch must fail.
3. Add CLI tests first in `internal/cli/gate_test.go` for the exact sibling-command
   grammar, implicit actor, stable stdout, and usage/semantic refusal split. Add status
   projection tests in `internal/status/gates_coexist_test.go` and
   `internal/status/boot_identify_test.go`.
4. Extend `internal/ensigncycle/recorded_gate_lifecycle_test.go` with the exact
   prepare-N → withdraw → cold-boot → prepare-N+1 → provider-room-close → consume
   sequence using a freshly built binary, real split-root commits, dirty sibling, and
   whole-tree byte oracles. This is the behavioral acceptance proof.
5. Extend the already registered `TestLiveDefaultHeadlessStopsAtGate` lane in
   `internal/ensigncycle/live_gate_stop_test.go` and its existing Claude runner
   orchestration with a withdrawn-start variant. The live oracle requires exactly one
   successful N+1 prepare and commit, the emitted-room presentation, and no legacy
   Briefing bind/close/consume/dispatch.
6. Delete from `internal/contractlint/fo_function_reference_invariant_test.go` every
   model-read assertion of prepare/record/consume command text, emitted output,
   presenter instructions, and lifecycle ordering. Add no
   withdraw/commit/prepare-content oracle. Preserve only the file's allowed structural
   closure, wiring, frontmatter, absence/dedup, and byte-cap checks.
7. Run focused packages and the registered applicable live lane when credentials are
   available, then `gofmt -w ./cmd ./internal`, `go test ./...`, and
   `go test ./... -race`.

## Expected implementation surface

Estimate baseline: s4 tip `e328ecc6`. The expected delta is 21 existing files,
approximately +980/-160 LOC. These are planning estimates, not line-count claims:
test/doc files may vary ±25% and production files ±15 lines as s4 lands. A new
production package, dependency, persistence shape, writer, or proof substitution
requires design re-entry; ordinary estimate variance is reported in the next Stage
Report.

| File | Estimated changed LOC | Purpose |
| --- | ---: | --- |
| `internal/gates/model.go` | +52/-18 | Withdrawal type, explicit classifier, validation, summary/readiness |
| `internal/gates/operation.go` | +86/-20 | Shared-lock withdrawal and explicit recorder/transition guards |
| `internal/gates/prepare.go` | +31/-13 | Target, replay, pre-publish, and post-publish explicit-state branches |
| `internal/gates/application.go` | +6/-4 | Explicit closed-state eligibility/consume guard |
| `internal/cli/cli.go` | +43/-12 | Exact sibling grammar, help, parsing, stable output |
| `internal/status/format.go` | +2/-1 | Schedule withdrawn recovery in `ready_gates` |
| `internal/gates/gates_test.go` | +130/-0 | Three-state validation, freeze, summary/readiness |
| `internal/gates/prepare_test.go` | +170/-0 | Successor publication, replay exclusion, rollback, retained authority |
| `internal/cli/gate_test.go` | +82/-0 | Grammar, stdout, implicit actor, byte-clean refusals |
| `internal/status/gates_coexist_test.go` | +24/-0 | Withdrawn field projection |
| `internal/status/boot_identify_test.go` | +32/-0 | Exact cold-boot ready row |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | +220/-8 | Fresh-binary supported sequence and Git/byte oracles |
| `internal/ensigncycle/live_gate_stop_test.go` | +42/-2 | Registered withdrawn-start live variant |
| `internal/ensigncycle/claude_live_runner_test.go` | +56/-6 | Existing lane fixture/oracle orchestration |
| `internal/contractlint/fo_function_reference_invariant_test.go` | +0/-60 | Remove model-read command/content/order assertions; retain allowed structural checks |
| `skills/fo-gate-lifecycle/SKILL.md` | +22/-8 | Preflight, withdrawal, commit, cold-boot prepare |
| `docs/specs/gate-resolution-frontmatter-contract.md` | +48/-18 | Canonical third state, authority, retained validation |
| `docs/site/concepts/gates-and-decisions.md` | +20/-5 | Operator-facing truthful recovery |
| `docs/site/reference/frontmatter-contract.md` | +14/-6 | Compact state and frozen-history invariants |
| `docs/site/reference/command-reference.md` | +14/-5 | Exact command and output |
| `docs/schema/entity.mdschema.yml` | +10/-5 | Machine-readable shape and state invariants |

No `internal/gates/io.go` change is expected: withdrawal reuses s4's existing
`validateRetainedAuthority` rather than adding or weakening an authority reader.

## Proposed documentation and skill wording

`skills/fo-gate-lifecycle/SKILL.md`:

> Require `prepare`, `withdraw`, `record`, `validate`, `eligibility`, and `consume`
> plus `--reason` in the single lifecycle capability preflight. When the current
> prepared room is stale and has no provider output, run
> `spacedock gate withdraw ENTITY --reason TEXT`, commit the entity, and do not
> present or close that attempt. On `withdrawn-awaiting-prepare`, run `gate prepare`
> for the successor, commit its emitted room and binding, present that room, and stop
> open. Never recover with `record --briefing`.

`docs/specs/gate-resolution-frontmatter-contract.md`:

> An attempt is open only when both `withdrawal` and `resolution` are absent,
> withdrawn when only `withdrawal` is present, and closed when only `resolution` is
> present. Withdrawn and closed attempts are retired and frozen. Withdrawal is
> available only through `spacedock gate withdraw ENTITY --reason TEXT`; the writer
> records `by: agent:first-officer`, UTC time, and the required reason after validating
> the complete retained attempt history. It writes no Resolution, provider evidence,
> application, room, or workflow status.

`docs/site/concepts/gates-and-decisions.md`:

> If a prepared room becomes stale before any provider decision, the first officer
> records a reasoned withdrawal. Withdrawal is not approve, revise, or hold: it
> preserves the old room without a Resolution or application. Cold boot then asks the
> first officer to prepare attempt N+1.

`docs/site/reference/frontmatter-contract.md`:

> Open means neither withdrawal nor Resolution; withdrawn means withdrawal alone;
> closed means Resolution alone. Withdrawn and closed attempts are frozen. A
> withdrawal is FO-attributed, timestamped, reasoned, request-backed, Resolution-free,
> provider-evidence-free, and application-free.

`docs/site/reference/command-reference.md`:

> `spacedock gate withdraw ENTITY --reason TEXT [--workflow-dir DIR]` — Retire the
> selected current-stage prepared attempt without a Resolution or application. The
> command records implicit `agent:first-officer` attribution; the next `gate prepare`
> appends a successor.

`docs/schema/entity.mdschema.yml`:

> neither withdrawal nor Resolution means open; withdrawal alone means withdrawn;
> Resolution alone means closed; withdrawn and closed attempts are frozen; withdrawal
> is request-backed and has fixed FO attribution, time, and nonblank reason.

## Stage Report: ideation

- DONE: Exercise the smallest real stale-open withdrawal and replacement path before choosing the durable shape.
  A throwaway real-CLI test bound request-backed room 1, withdrew it, cold-booted, appended room 2, recorded its provider Result, and failed if room-1 bytes or decision ownership changed.
- DONE: Define truthful withdrawal authority, frozen-history invariants, cold-boot projection, and byte-clean refusal semantics.
  The proposed approach fixes FO-only authority, the three mutually exclusive attempt states, retired-attempt freezing, `withdrawn-awaiting-prepare`, digest checks, and the complete nonzero/no-byte-change matrix for AC-1, AC-2, and AC-3.
- DONE: Declare exact files and LOC, documentation wording, and behavior-first tests without compatibility or generic history machinery.
  The gated baseline names 14 existing files at +579/-55 LOC with tolerance, concrete doc/skill text, falsifiable fresh-binary and real-Git tests, and explicit exclusions.

### Summary

Ideation defines a minimal attempt-local withdrawal recorded by the existing gate recorder, followed by the existing Briefing bind to append N+1. The exercised spike proved the state shape, cold-boot recovery, frozen room bytes, and truthful placement of the later Captain decision without adding compatibility or generic history machinery.

## Stage Report: ideation (cycle 2)

- DONE: Rebase the command, authority, and lifecycle on the current s4 preparation surface.
  The governing design now specifies `spacedock gate withdraw ENTITY --reason TEXT`,
  implicit durable `by: agent:first-officer`, the shared lock/CAS writer, and the exact
  prepare-N through provider-close/consume real-binary proof. The obsolete cycle-1
  `record --withdraw`, `briefing.json`, and recovery `record --briefing` paths are
  explicitly non-contractual.
- DONE: Audit every resolution-null branch that could mistake a withdrawn attempt for an open one.
  The design names `prepareTarget`, replay selection, pre-publish and post-publish
  mutation, room/chat/legacy record guards, application consumption, readiness,
  summary, validation, and transition freezing. Withdrawn is a Resolution-free
  terminal attempt state; existing provider-evidence rules remain unchanged.
- DONE: Define retained-authority cost, refusal integrity, behavioral proof, and the corrected implementation surface.
  Withdrawal reuses s4's complete O(history) validator before mutation, defers
  optimization until measured latency exceeds a declared bound, specifies byte-clean
  refusal and frozen-room oracles, separates structural contractlint from real-binary
  and registered live-FO evidence, and estimates 21 existing files at approximately
  +1,000/-105 LOC with no new package, dependency, writer, compatibility, or generic
  history machinery.

### Summary

Cycle 2 replaces the obsolete recorder-source design with a narrow FO-owned
`gate withdraw` command and rebases successor recovery on `gate prepare`. The
withdrawn attempt remains Resolution-free, provider-evidence-free, application-free,
fully retained-authority-validated, and frozen while N+1 follows the ordinary s4
prepare, provider record, and consume lifecycle.

## Stage Report: ideation (cycle 3)

- DONE: Remove every committed contractlint assertion of withdraw/commit/prepare command content or ordering in model-read prose.
  The governing proof policy now deletes s4's `gate prepare`, emitted-room,
  presenter-text, and lifecycle-order substring assertions and forbids replacement
  withdraw/commit/prepare prose oracles.
- DONE: Keep contractlint within its structural proof boundary.
  The retained allowance is limited to structural closure,
  skill/capability/reference wiring, frontmatter validity, absence/dedup invariants,
  and byte caps; the expected surface records the contractlint change as deletion-only.
- DONE: Preserve the accepted withdrawal, preparation, retained-authority, and behavioral proof design.
  The exact grammar, implicit FO attribution, explicit three-state audit, frozen
  history, unchanged provider-evidence rules, O(history) validator, fresh-binary
  lifecycle, and registered live FO lane remain unchanged.

### Summary

Cycle 3 removes model-read command content and ordering from contractlint's authority.
Real-binary split-root tests own the durable command sequence, and the existing live FO
lane owns model behavior; contractlint remains structural only.

## Stage Report: ideation (cycle 4 — AC evidence addendum)

- DONE: Remove every committed contractlint assertion of withdraw/commit/prepare command content or ordering in model-read prose.
  Cycle 3 keeps the implementation surface deletion-only for those assertions and assigns behavior/order exclusively to fresh-binary and live-FO proof.
- DONE: Keep contractlint within its structural proof boundary.
  Cycle 3 retains only structural closure, skill/capability/reference wiring, frontmatter validity, absence/dedup invariants, and byte caps.
- DONE: Preserve the accepted withdrawal, preparation, retained-authority, and behavioral proof design.
  Cycles 2–3 leave the exact command, three-state model, O(history) validator, frozen history, provider rules, and supported lifecycle unchanged.
- DONE: Cite concrete ideation evidence for zero-false-Resolution replacement.
  AC-1 — The throwaway real-CLI spike observed withdrawn N with no Resolution/application and provider Result ownership only on N+1; the s4 `prepareTarget`/replay/post-publish audit rebases that observation into the exact prepare N → withdraw → cold boot → prepare N+1 → `record --room` → consume proof, which fails on any decision placed on N.
- DONE: Cite concrete ideation evidence for canonical withdrawal integrity.
  AC-2 — The spike kept N's retained Briefing/request bytes identical, while the s4 audit anchors implicit FO withdrawal to `lockEntity`/`writeDocument`, full retained-authority validation, and `ValidateTransition`; the refusal matrix compares the whole entity/room tree and residual lock path so any refused write or later mutation of N fails.
- DONE: Cite concrete ideation evidence for cold recovery and retained authority.
  AC-3 — The spike observed fresh boot projection `withdrawn-awaiting-prepare`, and the s4 audit found `SummaryFileAt`, `Prepare`, `EligibilityFileAt`, and `ConsumeAt` all enter `validateRetainedAuthority`; the real-binary oracle requires withdrawn validate success, byte-clean ineligible eligibility/consume before N+1 closes, then consumption of N+1 while historical N remains revalidated.

### Summary

This addendum changes no design or product surface. It preserves the prior checklist
accounting and gives the ideation gate concrete, falsifiable citations for AC-1,
AC-2, and AC-3.

## Stage Report: implementation

- DONE: Implement the approved request-backed open-attempt withdrawal state and exact `gate withdraw ENTITY --reason TEXT` command through the existing locked writer, with fixed FO attribution and byte-clean refusal.
  Commit `5a8be3220` adds the third state and shared-lock writer; model, CLI, and refusal tests fail on false attribution, conflicting state, room/provider drift, repeat/closed/request-less targets, stale selection, or changed bytes.
- DONE: Integrate explicit open/withdrawn/closed semantics through prepare, record, transition, eligibility, summary, and ready-gate projection; prove prepare N → withdraw → cold boot → prepare N+1 → provider close → consume with immutable retained authority.
  `TestRecordedGateLifecycleWithdrawColdBootReplaceAndConsume` drives a freshly built binary and fails if N is reused/rewritten, N owns a Resolution/application, N+1 lacks provider authority, boot readiness drifts, consumption fails, or dirty sibling state is swept.
- FAILED: Delete the identified prose-command contractlint proof, keep the 21-file +980/-160 boundary unless a declared reset trigger occurs, and run focused/full/race plus the applicable behavior-first live lane only after cheap proofs pass.
  Base `070f36ae0` already removed the prose oracle; the authorized candidate is 20 existing files at +824/-60 with no semantic narrowing. Focused, full, race, and live-tag compile passed, but both registered live attempts reached the real Claude adapter and failed before FO work with HTTP 429 weekly quota, resetting 1pm Asia/Taipei; the corrected withdrawn-start fixture passed local setup and reached that same infrastructure boundary.

### Summary

Commit `5a8be3220` implements truthful FO withdrawal, retained frozen authority, successor preparation, provider close/consume, cold-boot recovery, operator documentation, and the registered durable-state/invocation-log live variant without transcript grammar. The production delta is smaller through reuse of s4 seams; the First Officer authorized the exact 20-file +824/-60 variance because contractlint was already structurally clean and acceptance scope is unchanged.

## Stage Report: validation

- DONE: Reproduce AC-1 through AC-3 at exact candidate 5a8be3220, including open/withdrawn/closed authority and byte-clean refusal.
  Focused gates/CLI/status/ensigncycle tests passed at `5a8be3220`; they fail if the command fabricates closure, exposes caller attribution, reuses N, mutates refusal bytes, or projects the wrong boot state.
- DONE: Exercise prepare N, withdraw, cold boot, prepare N+1, provider close, eligibility, and consume while proving N remains immutable and undecided.
  `TestRecordedGateLifecycleWithdrawColdBootReplaceAndConsume` passed through a freshly built binary; it fails if N changes or owns Resolution/application, N+1 lacks provider authority, boot is not singular, or consume does not advance.
- DONE: Verify the 20-file +824/-60 surface, focused/full/race evidence, docs/skill coherence, and the registered live journey when infrastructure permits.
  Diff audit found exactly 20 existing files at +824/-60; `go test -count=1 ./...`, `go test -count=1 -race ./...`, formatting check, and live-tag compile passed, and the withdrawn variant remains registered in `runtime-live-e2e.yml`.
- DONE: Perform the semantic adversarial pass over explicit open/withdrawn/closed lifecycle authority.
  A detached `5a8be3220` audit proved withdrawn N rejects room/chat record and consume byte-cleanly, remains ineligible, preserves its room through N+1 consume, and rejects later presented-inventory drift.
- SKIPPED: Run another real Claude withdrawn-start journey.
  At 2026-07-31 07:28 Asia/Taipei the local OAuth lane was still inside the already-observed HTTP 429 weekly-quota window before its 13:00 reset; per dispatch, no redundant live attempt was burned.
- DONE: Recommend PASSED with deferred risks separated from material findings.
  AC-1 through AC-3 have command/state evidence, no material finding or deferred product risk emerged, and candidate mutation or reviewer rerun was neither needed nor authorized.

### Summary

Validation recommends PASSED for exact candidate `5a8be3220`. Deterministic, repository-wide, race, live-registration, and detached adversarial evidence all preserve the approved behavior; only the registered real-Claude observation remains externally quota-blocked rather than product-red.

## Stage Report: validation (cycle 2)

- DONE: Verify the rebased PR candidate against origin/main at the exact pushed head and reconcile the stale prior validation authority.
  PR #580 is `MERGEABLE/CLEAN` at head `4dd5322e9a438d498a03e2192a8397c4d76c01e2` and base/merge-base `23ed415bb3f16393f7b5a0f6c19c9f259b6c4617`; this cycle is the only current validation authority.
- DONE: Run focused, full, race, format, diff, and applicable live/CI checks, triaging any failure from the new head rather than the old report.
  Focused gates/CLI/status/ensigncycle, `go test -count=1 ./...`, `go test -count=1 -race ./...`, live-tag compile, `gofmt -l ./cmd ./internal`, and `git diff --check` passed on `4dd5322e9`.
- DONE: Reconfirm the withdrawn-attempt lifecycle, exact 20-file/+824/-60 scope, and no unrelated rebase drift before a fresh validation gate.
  The approved pre-rebase delta was 20/+824/-60; the exact rebased PR is 20/+830/-66 because overlapping base edits change diff arithmetic, while `range-diff` marks both immediate pre/post-rebase candidate commits `=` patch-equivalent with no unrelated file or semantic drift.
- DONE: Reproduce AC-1 through AC-3 from their command/state evidence on the rebased candidate.
  `TestRecordedGateLifecycleWithdrawColdBootReplaceAndConsume` freshly built the binary and drove prepare N → withdraw → cold boot → prepare N+1 → room close → consume; it fails if N changes, owns closure, or boot/readiness cardinality drifts.
- DONE: Reproduce AC-2's byte-clean refusal, immutable withdrawn history, retained s4 authority, and unchanged closed-provider validation on exact head `4dd5322e9`.
  `TestWithdrawRefusalsLeaveEntityRoomAndLockBytesClean` covered blank reason, provider output, corrupt request, repeat/closed, lock, request-less, and stale-selection tree digests; `TestRecordedGateLifecycleWithdrawColdBootReplaceAndConsume` kept N byte-identical through N+1 consume and rejected later N-request and N+1-Result drift without entity mutation.
- DONE: Record every exact-head GitHub lane and recommend PASSED only after required live CI completed green.
  Run `30736599933`: offline, Claude Sonnet, Claude Opus, Codex, Pi, and journey-delta all completed SUCCESS; docs build and Ubuntu/macOS install checks also completed SUCCESS on `4dd5322e9`.
- DONE: Recommend PASSED with deferred risks separated from material findings.
  AC-1 through AC-3 have fresh exact-head behavioral evidence; no material finding, evidence defect, deferred product risk, or polish finding remains.

### Summary

Validation recommends PASSED for rebased candidate `4dd5322e9` on PR #580. The fresh lifecycle, focused/full/race/format/diff suites, patch-equivalent rebase audit, and every required exact-head live lane are green; the old `5a8be3220` report and its +824/-60 comparison are retained only as stale history.
