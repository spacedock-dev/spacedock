---
title: Warn on legacy gate application fields
status: validation
source: "Captain directive 2026-08-03: unknown keys in the exact gate application mapping must warn, not fail the state read."
score: "1.0"
sprint: durable-decisions
id: jympnaf11wg4qmd4z85a3ayv
gates:
    version: 1
    records:
        - id: gate:jympnaf11wg4qmd4z85a3ayv:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:jympnaf11wg4qmd4z85a3ayv-backlog-1
              briefing:
                id: briefing:jympnaf11wg4qmd4z85a3ayv:backlog:attempt-1:revision-1
                digest: sha256:899a41fbeb26148c4e30be8a8fbe36c8c18fde57e87823e09a3f772a843a0a68
                request-digest: sha256:612e30cdf011456f2204d8c05c092a8213f074722a427806282324bb35215dd8
                room-ref: ./warn-on-legacy-gate-application-fields/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:jympnaf11wg4qmd4z85a3ayv:backlog:1
                briefing: briefing:jympnaf11wg4qmd4z85a3ayv:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T13:53:01.677814Z"
                decision: approve
                reason: Captain directed dispatch of the bounded warning-only compatibility implementation; preserve strict errors for all other fields and reconcile shared io.go before merge.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:jympnaf11wg4qmd4z85a3ayv:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:jympnaf11wg4qmd4z85a3ayv-ideation-1
              briefing:
                id: briefing:jympnaf11wg4qmd4z85a3ayv:ideation:attempt-1:revision-1
                digest: sha256:9861a730590c25991e7af1d7b0812e4cac142d703a976ade09c1ba8899c4fc39
                request-digest: sha256:4467110f4d270d10859812828ea2ecfc7588e71d48d433ca964ef0dfd8ac784b
                room-ref: ./warn-on-legacy-gate-application-fields/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:jympnaf11wg4qmd4z85a3ayv:ideation:1
                briefing: briefing:jympnaf11wg4qmd4z85a3ayv:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T14:10:57.122806Z"
                decision: revise
                reason: 'Captain broadened the compatibility boundary: every unknown key in the exact records[*].attempts[*].application mapping must warn deterministically, not only action/blockers. Keep non-application unknowns, malformed/non-mapping applications, missing or invalid canonical fields, duplicate canonical keys, and bad bindings fatal. Update the task body, ACs, docs diff, and test matrix; do not authorize implementation under the old narrow brief.'
            - id: gate-attempt:jympnaf11wg4qmd4z85a3ayv-ideation-2
              briefing:
                id: briefing:jympnaf11wg4qmd4z85a3ayv:ideation:attempt-2:revision-1
                digest: sha256:d3b6425224611ec75fa9fa4d167d77d63bd5e88538d9ccd1ead278dbf57102ca
                request-digest: sha256:8a4af690385527958cbea4d249f2c8566acbfc5b311cbb5b345aee4650454007
                room-ref: ./warn-on-legacy-gate-application-fields/review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:jympnaf11wg4qmd4z85a3ayv:ideation:2
                briefing: briefing:jympnaf11wg4qmd4z85a3ayv:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-03T14:25:30.564741Z"
                decision: approve
                reason: 'Captain approved the revised exact-application warning boundary: every unknown application key warns and is ignored for canonical decode/eligibility; all non-application unknowns, malformed shapes, invalid canonical data, duplicates, and bad bindings remain fatal. Implement after WJ PR #610 merges.'
              application:
                target-stage: implementation
                state: consumed
started: 2026-08-03T13:53:46Z
worktree: .worktrees/spacedock-ensign-warn-on-legacy-gate-application-fields
---

Legacy or extension state can contain keys that are not part of the v1 application
schema. The canonical gate reader currently rejects that state before it can expose
valid `target-stage` and `state` fields. The Captain wants those unknown keys to be
warning-only when, and only when, they occur in the exact application mapping owned by
`gates.records[*].attempts[*].application`.

## Problem

`internal/gates/io.go` uses strict YAML known-field decoding for the complete gates
tree. Removed or extension application keys therefore produce fatal
`field <name> not found in type gates.Application` errors and block status,
eligibility, and gate recovery. Unknown application keys are compatibility findings,
not authority; unknown keys outside that exact mapping and invalid canonical values
remain fatal.

## Proposed approach

Add a diagnostic read path beside the existing `Read` contract. Read the entity bytes
once, parse the frontmatter, and retain the original `gates` node for compare-and-swap
and all later writes. Deep-clone only that node for validation. At exactly
`gates.records[*].attempts[*].application`, require a mapping, retain only the
canonical `target-stage` and `state` entries in the validation clone, and remove every
other mapping key. Collect one `{path, field}` warning for each unknown key, including
extensions such as `action`, `blockers`, `execution-hold`, and `feedback`. Sort and
de-duplicate warnings by path and field, then decode the clone with
`yaml.Decoder.KnownFields(true)` and run the unchanged `Validate` function. Unknown
application values are ignored for authority and are never silently honored.

Expose `ReadDiagnostics` (or the equivalent named warning-returning reader) and keep
`Read` as a compatibility wrapper that discards warnings. `status --validate` and
`gate validate` render stable `Warning:` lines that include the entity path and the
exact unknown field. Ordinary status, eligibility, and consume reads accept unknown
application keys without changing their existing output or authority behavior unless
a caller explicitly requests diagnostics. A non-mapping, null, or sequence
`application`, malformed YAML, an unrelated unknown key, an invalid canonical value,
a bad binding, a duplicate canonical key, or an unknown key at another path remains a
hard error.

The writer never receives the filtered clone. A successful compatibility read therefore
does not rewrite or normalize bytes; only an existing explicit state mutation can
change the entity, and that mutation still serializes the canonical two-field
application.

## Sequencing and dependency boundary

Science Officer concurrence is recorded: ideation may run while WJ validation is open,
but implementation and validation must wait for WJ's exact candidate
`d1aac2e035f4e52c94141a2deb218ab5872d0fa5` to be approved and merged. The safe sequence
is `WJ validation/gate -> jy implementation -> jy validation -> SK rebase and fresh
validation`; `jy` and WJ both touch `internal/gates/io.go`, so no jy merge or rebase
may discard WJ's exact-head evidence. No WJ or SK worktree is changed by ideation.

## Acceptance criteria

**AC-1 (VALUE) — Valid applications with arbitrary extensions remain readable.**
On a fixture set with valid canonical applications carrying arbitrary unknown keys
(`action`, `blockers`, `execution-hold`, `feedback`, and at least one nested or
non-string value), every read succeeds and retains the canonical
`target-stage`/`state` values. Each distinct unknown `{path, field}` produces exactly
one structured warning. `status --validate` and `gate validate` exit zero and report
the warnings. Unknown values are ignored; they do not change routing, eligibility,
authority, or stored canonical state. Against the current strict-reader baseline
(0/5 compatibility fixtures readable), the revised reader must read 5/5; an approve
Resolution with a valid next-stage target and `state: pending` remains
`approved-pending` eligible when its binding is current.

**AC-2 — Strict authority and shape checks remain fail-closed.**
Fixtures with an unrelated application/gates key, a non-mapping/null/sequence
application, malformed YAML, a missing target stage, an invalid application state, a
duplicate canonical key, a bad binding, or an unknown key outside the exact application
path exit nonzero and perform no write. A mapping with valid canonical fields plus any
number or shape of unknown extension keys is accepted; no other location is relaxed.

**AC-3 — Diagnostics are deterministic and operator-visible.**
Repeated reads of the same entity produce the same sorted, de-duplicated warning list,
independent of source mapping order. The explicit validation surfaces print one stable
`Warning:` line per `{path, field}` and keep the normal success/exit result. Ordinary
read output remains byte-compatible unless diagnostics are requested.

**AC-4 — Compatibility reads preserve source bytes and canonical writes.**
A successful diagnostic read leaves the original entity bytes unchanged, including
unrelated frontmatter and application formatting. A later approved state mutation
writes only the canonical `target-stage` and `state` application shape through the
existing locked writer.

## Semantic risk and controls

The deliberate risk is that an unknown application extension may have represented a
constraint in an older producer. Under this contract the reader ignores that key for
decode, eligibility, and authority, and reports it as a warning; it never silently
honors `execution-hold`, `feedback`, or another extension. Canonical behavior is
unchanged: an approve Resolution with a valid next-stage `target-stage` and
`state: pending` remains eligible, while missing or invalid canonical values remain
fatal. A bad binding outside the application mapping remains fatal. The exact path
boundary, warning output, strict validation of the filtered clone, and byte
preservation limit this compatibility exception.

## Expected surface and semantic boundaries

Expected changes are limited to `internal/gates/io.go` and focused gates tests,
`internal/status/discover.go`/`validate.go` and status fixtures,
`internal/cli/cli.go` and gate-validation golden tests, plus the gate-resolution
contract and frontmatter/command reference wording. Estimate 6–10 files and
`+180/-45` lines, with a 2x tolerance. Allowed semantic change: read-time tolerance
and explicit warning diagnostics for every unknown key at one exact application path.
Stored canonical schema, `Application` fields, `Validate`, authority spending,
CAS/write behavior, command grammar, workflow taxonomy, and worker transport remain
unchanged.

## Mechanism choices and rejected alternatives

- **Clone, strip, then strict-decode** serves AC-1/AC-2. Turning off `KnownFields` or
  decoding into `map[string]any` would also accept unrelated prototype fields and
  weaken fail-closed authority checks.
- **Retain the original node and bytes** serves AC-4. Decoding and re-marshalling the
  filtered document as the write source would normalize opaque formatting and could
  bypass the existing compare-and-swap expectation.
- **A warning-returning reader plus an old `Read` wrapper** serves AC-3 while keeping
  ordinary read callers stable. Printing warnings from every status/eligibility read
  would change scripts and scheduler output without an explicit diagnostic request.
- **Exact path with canonical-field allow-list** serves all four ACs. A generic
  unknown-field compatibility mode outside the application mapping or a migration
  rewrite is broader than the value and would conceal future schema errors. Within the
  exact mapping, warning every non-canonical key is intentional so future extensions
  do not create a new fatal-read outage; those extensions are not authority.

## Documentation diff proposed by ideation

In `docs/specs/gate-resolution-frontmatter-contract.md`, replace “The binary-owned
model is closed: unsupported fields inside `gates` fail closed” with: “The
binary-owned model is closed for canonical validation and writes. A reader tolerates
unknown keys only under each `records[*].attempts[*].application` mapping, reports
them as warnings on explicit `status --validate` or `gate validate`, ignores them for
authority, and never writes them. All other unknown or malformed fields fail closed.”
Keep the canonical example unchanged and state that `target-stage` and `state` remain
the only canonical application fields.

In `docs/site/reference/frontmatter-contract.md`, replace “Unknown or prototype fields
inside binary-owned `gates` fail closed” with the same bounded exception and add:
“`status --validate` and `gate validate` print the entity path and unknown field;
normal reads preserve their existing output. Unknown application extensions are
ignored and are never authority.” Replace the following “contains exactly
`target-stage` and `state`” sentence with “Canonical writes contain exactly
`target-stage` and `state`; read-only compatibility may observe arbitrary unknown
application keys as warnings.”

In `docs/site/reference/command-reference.md`, extend the `gate validate` row with:
“It also reports warning-only unknown keys under the exact application path (including
legacy `action`/`blockers` and extensions such as `execution-hold`/`feedback`);
unrelated unknown fields remain errors.”

## Test plan

Start with the reader fixture matrix and the decode-copy probe before implementation:
valid `action`, valid `blockers`, both together, arbitrary `execution-hold`, arbitrary
`feedback`, a nested extension value, repeated reads, reversed mapping order, repeated
unknown names at separate application paths, unrelated key,
non-mapping/null/sequence application, malformed YAML,
missing target, invalid state, duplicate canonical key, bad binding, and unknown keys
at a non-application path. Assert warning order/content and de-duplication, canonical
projection, strict errors, ignored extension values, and unchanged source bytes. Add
status and gate CLI goldens for explicit warning output (including arbitrary extension
names) and an ordinary-read compatibility case.

Update the old `TestRemovedApplicationShapesFailClosedWithoutMutation` rows so every
unknown key in the exact application mapping is warning-only; retain fail-closed rows
for `not-applicable` and all unrelated paths. Run focused gates/status/CLI tests, then
`gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.

## Riskiest mechanism probe

The throwaway `go run /tmp/jy_ideation_revision_probe.go` exercised a YAML-node deep
clone, exact-path filtering, deterministic sorting, strict decode, validation, and
source-node preservation. It produced sorted warnings for `action`, `blockers`,
`execution-hold`, `feedback`, and a nested extension. The same run kept an unrelated
key, malformed YAML, invalid state, missing target, non-mapping application, duplicate
canonical key, and a non-application binding-shaped unknown field as errors. An
unknown binding-shaped key inside the application mapping is covered by the warning
path. Implementation tests must make the evidence durable; the temporary probe is not
production code.

## Stage Report: ideation

- DONE: Validate the revised warning-only compatibility design against every acceptance criterion
  AC-1 through AC-4 define measurable success, strict-negative controls, diagnostics, and byte-preserving writes.
- DONE: Exercise the riskiest decode-copy, deterministic-warning, byte-preservation, and strict-negative paths
  `go run /tmp/jy_ideation_revision_probe.go` observed arbitrary extension warnings, unchanged source node, unrelated-key failure, malformed-shape failures, invalid-state failure, duplicate-canonical failure, and missing-target failure.
- DONE: Record the chosen direction, rejected alternatives, exact probes, and a complete ideation revision report
  The proposed reader API, canonical allow-list at the exact application path, sequencing dependency, semantic risk, doc diff, alternatives, and test matrix are recorded above.
- AC-1 evidence (see AC-1 above): `go run /tmp/jy_ideation_revision_probe.go` retained
  canonical `target-stage: implementation` and `state: pending` while emitting sorted
  warnings for `action`, `binding`, `blockers`, `execution-hold`, `feedback`, and
  `nested`; strict canonical decode succeeded.
- AC-2 evidence (see AC-2 above): the same probe returned errors for an unrelated
  top-level key, malformed YAML, non-mapping/null/sequence applications, a missing
  target, invalid state, duplicate canonical key, and a non-application binding-shaped
  field. The exact application mapping alone was filtered.
- AC-3 evidence (see AC-3 above): the probe sorted warnings by `{path, field}` and
  printed a stable warning list; the test plan adds reversed-order and repeated-read
  fixtures plus explicit `status --validate`/`gate validate` warning goldens.
- AC-4 evidence (see AC-4 above): the probe compared the original cloned source node
  after filtering and confirmed `source-preservation=ok`; the proposed writer keeps
  that original node and the test plan adds raw-byte preservation and canonical-write
  assertions.

### Summary

The ideation revision broadens the compatibility seam from two named legacy fields to
every unknown key in the exact
`gates.records[*].attempts[*].application` mapping. The reader clones and filters all
non-canonical keys, emits deterministic warnings, retains strict canonical validation,
and keeps the original node for compare-and-swap and writes. Unknown keys outside that
mapping, malformed shapes, duplicate canonical keys, invalid canonical values, and bad
bindings remain fatal. Extensions such as `execution-hold` and `feedback` are ignored,
not honored. Science Officer concurrence keeps implementation after WJ's exact-head
validation and before SK's rebase/fresh validation; no production code or sibling
worktree was changed in this revision.

## Stage Report: implementation

- DONE: Implement the approved exact-application warning boundary
  Commit `0ad6c36d08d36d0e066a60386a0a1e7783cc47f8` adds `ReadDiagnostics`, clones only the gates node for strict validation, strips every non-canonical key under `gates.records[*].attempts[*].application`, emits sorted/de-duplicated `{path, field}` warnings, and keeps the original node for CAS and writes. Canonical `target-stage` and `state` remain authoritative; malformed application shapes, non-application unknowns, invalid canonical values, duplicate canonical keys, and bad bindings remain fatal.
- DONE: Surface diagnostics and preserve compatibility semantics
  `status --validate` and `gate validate` print stable `Warning:` lines with entity paths; ordinary status, eligibility, and consume reads stay warning-free. Explicit gate validation still checks retained request/provider authority. Canonical state mutation removes ignored extensions from the written application shape.
- DONE: Update fixtures, authority controls, and contract wording
  Commit `c5bd1225b9434763b9a76fd238d9699a2e564bef` updates boot readiness controls to prove `blockers` and `execution-hold` do not gate a valid approval, and aligns `docs/schema/entity.mdschema.yml`, frontmatter contracts, and command reference text with the exact warning boundary.
- DONE: Verify focused behavior and formatting
  `gofmt -w ./cmd ./internal` passed. Focused `go test` passed for gates diagnostics, status validation warnings, gate validation warnings/retained-binding rejection, boot readiness extension controls, and the recorded-gate consume lifecycle. The focused test matrix covers arbitrary nested values, warning order and de-duplication, source-byte preservation, malformed/null/sequence applications, missing/invalid canonical fields, duplicate canonical keys, unknown fields outside application, and bad retained request binding.
- DONE: Record implementation size and full-suite evidence
  Diff from WJ tip `ea6723acb630c7eac64712f5b30e08f35a5a44e6`: 13 files, `+466/-47` lines (net `+419`), including `+310/-13` production/docs and `+156/-34` focused test updates. Full `go test ./...` was run; the remaining failure is the pre-existing pilot-manifest fixture checkout omission (six listed state paths absent), while the formerly stale removed-application lifecycle and boot controls now pass.
- DONE: Run the required race suite
  `go test ./... -race` completed with the same six pre-existing pilot-manifest state-path omissions and no race or implementation failures; every other package passed under the race detector.

### Summary

Implementation is complete on branch `spacedock-ensign/warn-on-legacy-gate-application-fields` at `c5bd1225b9434763b9a76fd238d9699a2e564bef`, directly based on WJ tip `ea6723acb630c7eac64712f5b30e08f35a5a44e6`. Unknown application extensions warn only at explicit validation surfaces and cannot alter eligibility or authority; all unrelated and malformed gate state remains strict. Both full and race suites are green outside the pre-existing six missing pilot-manifest paths.
