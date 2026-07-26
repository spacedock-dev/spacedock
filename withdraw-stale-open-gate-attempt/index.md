---
title: Withdraw a stale open gate attempt without fabricating a decision
status: ideation
source: "Observed by the Subspace Shaping FO on 2026-07-26 after a legitimate sprint re-scope left a frozen request-backed attempt open with no truthful exit."
started: 2026-07-26T10:50:55Z
completed:
verdict:
score: 1.0
worktree:
issue:
sprint: durable-decisions
id: 0m6vtrw4qh9w4x6bn06x5hen
gates:
    version: 1
    current:
        gate: gate:0m6vtrw4qh9w4x6bn06x5hen:backlog
    records:
        - id: gate:0m6vtrw4qh9w4x6bn06x5hen:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:0m6vtrw4qh9w4x6bn06x5hen-backlog-1
              briefing:
                id: briefing:0m6vtrw4qh9w4x6bn06x5hen:backlog:attempt-1:revision-1
                digest: sha256:9bfedeb38906e04bae528cedfdb96f101efaa1d63c819b44922a8ee6e5db60f6
                digest-domain: canonical-bytes
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
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

A prepared request-backed gate attempt is correctly frozen, but a legitimate re-scope currently has no truthful operation that retires the open attempt. The operator must either record `hold` or another Resolution against a Briefing it already knows is stale, or hand-edit gate state. The first path fabricates a decision; the second bypasses the only canonical writer.

## Problem and boundary

Add one narrow recorder action for a selected, current-stage, request-backed open attempt. The action records who withdrew it, when, and why; it records no Resolution, provider evidence, application, status change, or dispatch. The retained Briefing and request stay digest-bound in their original folder. A later ordinary `--briefing` bind appends attempt N+1 instead of rebinding attempt N.

This task does not add compatibility or migration, generic event/history machinery, a cancel/supersede vocabulary, a new writer, automatic room construction, withdrawal of closed attempts, or withdrawal by Captains/ensigns. Closed attempts and chat-only open attempts keep their existing behavior.

## Exercised spike

On 2026-07-26 a throwaway implementation exercised `TestSpikeWithdrawPreparedRoomColdBootAndReplace` through the real CLI on a folder-form split-root fixture. Attempt 1 bound canonical `briefing.json` plus request-digested `request.json`; `gate record --withdraw` produced `state=withdrawn`; a fresh `status --boot --identify --json` produced `readiness=withdrawn-awaiting-prepare`; binding room 2 appended attempt 2; and `gate record --room` recorded its actual Captain approval. The first attempt retained no Resolution or application, the second alone carried the approval, and attempt-1 Briefing/request bytes were identical before and after. The focused Go test passed in 0.10s; the throwaway code was removed.

The spike settled the riskiest mechanisms: an explicit third attempt state round-trips through canonical YAML, the existing ordered-attempt rule supplies successor identity without a lineage field, and the existing boot scheduling index can carry the one recovery action.

## Proposed approach

### Command and truthful authority

Extend the existing semantic recorder grammar:

```text
spacedock gate record ENTITY --withdraw --actor agent:first-officer --reason TEXT [--workflow-dir DIR]
```

`--withdraw` is exactly one semantic source alongside `--briefing`, `--room`, and `--decision`; combining them refuses. `agent:first-officer` is the only accepted actor because withdrawal is an FO lifecycle mutation, not a Captain verdict. The reason is required after whitespace trimming. The binary supplies the RFC3339Nano UTC timestamp and derives the selected gate and attempt; callers supply no timestamp, id, target, or operation envelope.

The implementation stays inside `RecordSemantic`, `lockEntity`, subtree compare-and-swap validation, `atomicWrite`, and later `state commit`. No second command writer or direct frontmatter path is introduced.

### Minimum durable shape

Attempt N gains only this optional mapping:

```yaml
withdrawal:
  by: agent:first-officer
  at: "2026-07-26T11:30:00.123456Z"
  reason: Sprint re-scope replaced the reviewed candidate.
```

The field name supplies the event type and the owning attempt supplies identity, so `type`, `id`, `state`, successor, and lineage fields would be redundant. Attempt states are mutually exclusive:

- open: no `withdrawal` and no `resolution`;
- withdrawn: `withdrawal` present, no `resolution`, `provider-evidence`, or `application`;
- closed: `resolution` present, no `withdrawal`, with the existing application rules.

Validation accepts no other shape. Existing canonical-v1 records need no rewrite; unknown fields still fail closed.

### Withdrawal, replacement, and frozen history

Under the existing entity lock, withdrawal succeeds only when all of these are true:

1. the workflow status resolves to the selected current logical gate;
2. its last ordered attempt is open and has no prior withdrawal;
3. the attempt has a nonempty `request-digest`;
4. retained `briefing.json` still matches its canonical digest/id;
5. retained `request.json` still matches its digest and exact gate, attempt, Briefing, actor, approver, and Captain authority;
6. actor is exactly `agent:first-officer` and reason is nonblank.

Only after every check does one atomic entity replacement add `withdrawal`. The command never writes room files. Every nonzero result—including wrong actor, blank reason, stale selection, chat-only/open/closed state, changed Briefing/request bytes, lock contention, and repeat withdrawal—leaves entity and room bytes unchanged and leaves no lock residue.

`ValidateTransition` freezes withdrawn attempts exactly as it freezes closed attempts. `gate validate` rechecks the retained Briefing/request digests for every withdrawn attempt, including historical attempt N after N+1 exists. `record --briefing` treats a withdrawn last attempt as retired: it derives and validates the ordinary successor attempt id, appends N+1, and never changes N. The existing state-commit folder scope commits `index.md` and both room trees while excluding dirty sibling entities.

### Status, boot, and First Officer recovery

`gate validate` and `gate-state` report `withdrawn` for a selected withdrawn attempt while `gate-resolution`, `gate-decision`, and application fields remain empty. `gate-readiness` reports `withdrawn-awaiting-prepare`; identify boot includes that actionable row in `ready_gates`.

The FO handles that row in one way: read the entity, build a replacement room for derived attempt N+1, bind it with the existing `gate record --briefing`, commit, and present it. After binding, the normal projection is `state=open` / `awaiting-captain`. A restart therefore never confuses a withdrawal with validation still in progress, a Captain decision, or the replacement attempt.

### Mechanism choices

| Mechanism | Value AC | Simplest alternative | Why insufficient |
| --- | --- | --- | --- |
| Attempt-local `withdrawal` mapping | AC-1, AC-2 | Record `hold` | `hold` asserts a Captain decision and creates an application. |
| `gate record --withdraw` semantic source | AC-1, AC-2 | Add a `gate withdraw` writer or hand edit | Both split canonical write ownership; hand edit also bypasses lock/CAS. |
| Separate withdraw, then existing `--briefing` prepare | AC-1, AC-3 | Require replacement Briefing in the withdrawal command | A legitimate re-scope may pause before the replacement exists; atomic coupling makes cold-boot recovery impossible. |
| Ordered successor after a withdrawn attempt | AC-1, AC-2 | Add successor/lineage pointers | Existing ordered attempts and derived ids already identify N+1. |
| `withdrawn-awaiting-prepare` ready row | AC-3 | Reuse `validating` or omit the row | Both make restart ambiguous and can strand the retired attempt. |
| Digest revalidation plus frozen-transition guard | AC-2 | Trust Git history alone | The recorder must refuse against current corrupt bytes before a later commit can preserve them. |

## Acceptance criteria

**AC-1 (VALUE)** One stale prepared attempt can be retired and replaced with exactly zero false Resolutions: the resulting logical gate has attempt N with one withdrawal and no Resolution/application, plus attempt N+1 carrying only its actual recorded decision.

Test: a fresh-binary folder-form CLI fixture binds room N, withdraws, prepares N+1, records N+1's room-backed Result, and asserts the two exact attempt shapes. The test fails if withdrawal maps to `hold`, rewrites N, reuses N's id, or places the actual decision on N.

**AC-2 (INTEGRITY)** A withdrawn attempt remains immutable, digest-verifiable historical evidence, and every refused withdrawal changes zero entity or room bytes.

Test: the lifecycle fixture compares attempt-N YAML and whole-room bytes after withdrawal, replacement, validation, and decision; `gate validate` re-verifies historical Briefing/request digests. A table mutates actor, blank reason, current selection, Briefing bytes, request bytes, closed/chat-only state, and repeat-withdraw inputs, then asserts nonzero exit, identical pre/post tree digest, and no lock file. A direct transition test mutates/deletes a withdrawn attempt and must fail.

**AC-3 (RECOVERY)** A cold boot after withdrawal exposes exactly one `ready_gates` row with `withdrawn-awaiting-prepare`; after N+1 binds, that same entity projects the N+1 attempt as `awaiting-captain`, with no dispatch or workflow-status change in either state.

Test: the real split-root lifecycle restarts through `status --boot --identify --json` on both sides of replacement and asserts exact rows, attempt ids, status bytes, and empty dispatchable output. It fails if withdrawal is omitted, projected as `validating`/decision-ready, or if the old attempt remains current.

## Behavior-first test plan

1. Add focused model/transition tables in `internal/gates/gates_test.go` for the three exclusive states, minimum withdrawal validation, frozen history, readiness, and “withdrawn means append successor.” Estimated <2s; adversarial edit: change the open test back to `Resolution == nil`, which must make the successor assertion fail.
2. Extend `internal/ensigncycle/recorded_gate_lifecycle_test.go` with the fresh-binary folder-form path from AC-1 through AC-3, the byte-clean refusal matrix, cold boot, room-backed replacement decision, and real `state commit` Git assertions. Estimated 5–15s; fixture includes dirty sibling state and asserts commits contain `index.md` plus both rooms but exclude the sibling.
3. Extend the FO skill contract test before changing command text. It must require `--withdraw`, `withdrawn-awaiting-prepare`, FO-only actor/reason, commit-after-withdraw, and re-prepare rather than decision/consume. Estimated <1s; deleting any one lifecycle instruction must fail the test.
4. Run `go test ./internal/gates ./internal/ensigncycle ./internal/contractlint`, then the repository-required `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`. No live host test is needed: the claim is binary state behavior, boot projection, Git durability, and shipped FO instruction text, all covered by fresh-binary/real-Git fixtures.

## Expected surface

Baseline: 14 existing files, approximately +579/-55 LOC (net +524). Tolerance is ±20% per test/doc file and ±10 LOC per production file; any new file, dependency, production package, or unlisted file requires design re-entry.

| File | Expected LOC | Purpose |
| --- | ---: | --- |
| `internal/gates/model.go` | +38/-6 | Withdrawal model, exclusive-state validation, summary/readiness state |
| `internal/gates/operation.go` | +68/-10 | Semantic source, authority/refusals, withdraw write, successor and freeze rules |
| `internal/gates/io.go` | +16/-0 | Historical withdrawn-room digest verification |
| `internal/cli/cli.go` | +28/-8 | `--withdraw` grammar, help, source validation, stable output |
| `internal/status/format.go` | +1/-1 | Include the recovery readiness in `ready_gates` |
| `internal/gates/gates_test.go` | +105/-0 | State, transition, refusal, and successor unit tests |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | +190/-0 | Fresh-binary, cold-boot, real-Git, byte-clean lifecycle test |
| `internal/contractlint/fo_function_reference_invariant_test.go` | +15/-3 | FO lifecycle text smoke test added before skill text changes |
| `skills/fo-gate-lifecycle/SKILL.md` | +18/-6 | Preflight, withdraw/commit/re-prepare, and resume behavior |
| `docs/specs/gate-resolution-frontmatter-contract.md` | +42/-10 | Canonical schema and lifecycle contract |
| `docs/site/concepts/gates-and-decisions.md` | +20/-4 | Operator-facing stale-open recovery explanation |
| `docs/site/reference/frontmatter-contract.md` | +7/-3 | Compact field/invariant reference |
| `docs/site/reference/command-reference.md` | +10/-3 | Public grammar and readiness capability |
| `docs/schema/entity.mdschema.yml` | +6/-4 | Machine-readable writer and state invariants |

## Proposed documentation and skill wording

`skills/fo-gate-lifecycle/SKILL.md` preflight changes from requiring the existing flags to: “Require `record`, `validate`, `eligibility`, `consume`, `--briefing`, `--room`, `--decision`, `--withdraw`, `--actor`, `--reason`.” Add: “For `withdrawn-awaiting-prepare`, commit the withdrawal if needed, derive and retain attempt N+1, bind it with `--briefing`, commit, and present; never record or consume the withdrawn attempt.”

`docs/specs/gate-resolution-frontmatter-contract.md` changes “Resolution absence means open; Resolution presence means closed” to: “Neither withdrawal nor Resolution means open; withdrawal alone means withdrawn; Resolution alone means closed. Withdrawn and closed attempts are frozen.” Add the exact YAML mapping and command grammar above, FO-only authority, request-backed/current-attempt checks, byte-clean refusals, and successor-append rule.

`docs/site/concepts/gates-and-decisions.md` adds: “If a prepared room becomes stale before any decision, the first officer records a reasoned withdrawal. This is not approve, revise, or hold: it writes no Resolution or application, preserves the old room, and boot asks the first officer to prepare attempt N+1.”

`docs/site/reference/frontmatter-contract.md` changes “Resolution absence means open and presence means closed” to: “An attempt is open when both `withdrawal` and `resolution` are absent, withdrawn when only `withdrawal` is present, and closed when only `resolution` is present. Withdrawn and closed attempts are frozen; withdrawal has FO attribution, time, and a required reason.”

`docs/site/reference/command-reference.md` adds this table row: “`spacedock gate record <entity> --withdraw --actor agent:first-officer --reason TEXT` — Retire the selected request-backed open attempt without a Resolution or application; the next Briefing bind appends a successor.” Its capability-preflight sentence also names `--withdraw`.

`docs/schema/entity.mdschema.yml` changes the gate invariant to: “neither withdrawal nor Resolution means open; withdrawal alone means withdrawn; Resolution alone means closed; withdrawn and closed attempts are frozen,” and keeps the writer exactly `spacedock gate record`.

## Stage Report: ideation

- DONE: Exercise the smallest real stale-open withdrawal and replacement path before choosing the durable shape.
  A throwaway real-CLI test bound request-backed room 1, withdrew it, cold-booted, appended room 2, recorded its provider Result, and failed if room-1 bytes or decision ownership changed.
- DONE: Define truthful withdrawal authority, frozen-history invariants, cold-boot projection, and byte-clean refusal semantics.
  The proposed approach fixes FO-only authority, the three mutually exclusive attempt states, retired-attempt freezing, `withdrawn-awaiting-prepare`, digest checks, and the complete nonzero/no-byte-change matrix for AC-1, AC-2, and AC-3.
- DONE: Declare exact files and LOC, documentation wording, and behavior-first tests without compatibility or generic history machinery.
  The gated baseline names 14 existing files at +579/-55 LOC with tolerance, concrete doc/skill text, falsifiable fresh-binary and real-Git tests, and explicit exclusions.

### Summary

Ideation defines a minimal attempt-local withdrawal recorded by the existing gate recorder, followed by the existing Briefing bind to append N+1. The exercised spike proved the state shape, cold-boot recovery, frozen room bytes, and truthful placement of the later Captain decision without adding compatibility or generic history machinery.
