---
title: Warn on legacy gate application fields
status: ideation
source: "Captain directive 2026-08-03: legacy application.action and application.blockers must warn, not fail the state read."
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
started: 2026-08-03T13:53:46Z
---

Legacy state written before the v1 application-schema cut can contain
`application.action` and `application.blockers`. The canonical gate reader currently
rejects that state before it can expose the valid `target-stage` and `state` fields.

## Problem

`internal/gates/io.go` uses strict YAML known-field decoding for the complete gates
tree. The removed application keys therefore produce a fatal
`field action not found in type gates.Application` error and block status, eligibility,
and gate recovery. These two legacy keys are compatibility findings, not authority;
unrelated unknown keys and invalid canonical values remain fatal.

## Proposed approach

Add a diagnostic read path beside the existing `Read` contract. Read the entity bytes
once, parse the frontmatter, and retain the original `gates` node for compare-and-swap
and all later writes. Deep-clone only that node for validation. At exactly
`gates.records[*].attempts[*].application`, remove mapping keys named `action` or
`blockers` from the clone and collect `{path, field}` warnings. Sort and de-duplicate
the warnings by path and field, then decode the clone with `yaml.Decoder.KnownFields(true)`
and run the unchanged `Validate` function.

Expose `ReadDiagnostics` (or the equivalent named warning-returning reader) and keep
`Read` as a compatibility wrapper that discards warnings. `status --validate` and
`gate validate` render stable `Warning:` lines that include the entity path and the
exact legacy field. Ordinary status, eligibility, and consume reads accept the legacy
keys without changing their existing output or authority behavior unless a caller
explicitly requests diagnostics. A non-mapping `application`, malformed YAML, an
unrelated unknown key, an invalid canonical value, a bad binding, or any legacy key at
another path remains a hard error.

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

**AC-1 (VALUE) — Valid legacy applications remain readable.**
On a fixture set with four valid canonical applications carrying the legacy keys (two
`action`, two `blockers`), all four reads succeed and retain the canonical
`target-stage`/`state` values, with exactly four structured warnings. The current
strict reader is the moving baseline: it succeeds on 0/4. `status --validate` and
`gate validate` exit zero for those entities.

**AC-2 — Strict authority and shape checks remain fail-closed.**
Fixtures with an unrelated application/gates key, a non-mapping application, malformed
YAML, a missing target stage, an invalid application state, a bad binding, or a legacy
key outside the exact application path exit nonzero and perform no write. A mapping
that contains only `action`/`blockers` plus valid canonical fields is the sole
compatibility exception.

**AC-3 — Diagnostics are deterministic and operator-visible.**
Repeated reads of the same entity produce the same sorted, de-duplicated warning list;
the explicit validation surfaces print one stable `Warning:` line per `{path, field}`
and keep the normal success/exit result. Ordinary read output remains byte-compatible
unless diagnostics are requested.

**AC-4 — Compatibility reads preserve source bytes and canonical writes.**
A successful diagnostic read leaves the original entity bytes unchanged, including
unrelated frontmatter and application formatting. A later approved state mutation
writes only the canonical `target-stage` and `state` application shape through the
existing locked writer.

## Expected surface and semantic boundaries

Expected changes are limited to `internal/gates/io.go` and focused gates tests,
`internal/status/discover.go`/`validate.go` and status fixtures,
`internal/cli/cli.go` and gate-validation golden tests, plus the gate-resolution
contract and frontmatter/command reference wording. Estimate 6–10 files and
`+180/-45` lines, with a 2x tolerance. Allowed semantic change: read-time tolerance
and explicit warning diagnostics for exactly two historical keys at one path. Stored
canonical schema, `Application` fields, `Validate`, authority spending, CAS/write
behavior, command grammar, workflow taxonomy, and worker transport remain unchanged.

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
- **Exact path and field allow-list** serves all four ACs. A generic unknown-field
  compatibility mode or a migration rewrite is broader than the value and would
  conceal future schema errors.

## Documentation diff proposed by ideation

In `docs/specs/gate-resolution-frontmatter-contract.md`, replace “The binary-owned
model is closed: unsupported fields inside `gates` fail closed” with: “The binary-owned
model is closed for canonical validation and writes. A reader tolerates only legacy
`action` and `blockers` keys under each `records[*].attempts[*].application`, reports
them as warnings on explicit `status --validate` or `gate validate`, and never writes
them. All other unknown or malformed fields fail closed.” Keep the canonical example
unchanged.

In `docs/site/reference/frontmatter-contract.md`, replace “Unknown or prototype fields
inside binary-owned `gates` fail closed” with the same bounded exception and add:
“`status --validate` and `gate validate` print the entity path and legacy field; normal
reads preserve their existing output.” Replace the following “contains exactly
`target-stage` and `state`” sentence with “Canonical writes contain exactly ...;
read-only compatibility may observe the two named legacy keys as warnings.”

In `docs/site/reference/command-reference.md`, extend the `gate validate` row with:
“It also reports warning-only legacy `application.action`/`application.blockers` at
the exact application path; unrelated unknown fields remain errors.”

## Test plan

Start with the reader fixture matrix and the decode-copy probe before implementation:
valid `action`, valid `blockers`, both together, repeated reads, reversed mapping order,
unrelated key, non-mapping application, malformed YAML, missing target, invalid state,
bad binding, and legacy keys at a non-application path. Assert warning order/content,
canonical projection, strict errors, and unchanged source bytes. Add status and gate
CLI goldens for explicit warning output and an ordinary-read compatibility case.

Update the old `TestRemovedApplicationShapesFailClosedWithoutMutation` rows only for
`action` and `blockers`; retain fail-closed rows for `execution-hold`, `feedback`,
`not-applicable`, and all unrelated fields. Run focused gates/status/CLI tests, then
`gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.

## Riskiest mechanism probe

The throwaway `go run ./tmp/jy_ideation_probe.go` exercised a YAML-node deep clone,
exact-path stripping, deterministic sorting, strict decode, validation, and source-node
preservation. It produced four sorted warnings (`records[0]` action/blockers followed
by `records[1]` action/blockers); the unrelated key remained a strict unknown-field
error; `invalid-state` and `missing-target` remained `Validate` errors. The probe was
removed after the result was recorded; implementation tests must make the evidence
durable.

## Stage Report: ideation

- DONE: Validate the smallest warning-only compatibility design against every acceptance criterion
  AC-1 through AC-4 define measurable success, strict-negative controls, diagnostics, and byte-preserving writes.
- DONE: Exercise the riskiest decode-copy, deterministic-warning, byte-preservation, and strict-negative paths
  `go run ./tmp/jy_ideation_probe.go` observed sorted warnings, unchanged source node, unrelated-key failure, invalid-state failure, and missing-target failure.
- DONE: Record the chosen direction, rejected alternatives, exact probes, and a complete ideation report
  The proposed reader API, exact allow-list/path, sequencing dependency, expected surface, doc diff, alternatives, and test matrix are recorded above.

### Summary

Ideation selects a narrow diagnostic-read compatibility seam: clone and filter only
legacy `action`/`blockers`, then retain strict canonical validation and the original
write node. Science Officer concurrence places implementation after WJ's exact-head
validation and before SK's rebase/fresh validation; no production code or sibling
worktree was changed in this stage.
