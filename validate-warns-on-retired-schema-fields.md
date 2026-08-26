---
id: wyzg6knr8whtmg79mxkb78jg
title: status --validate warns on retired gate-schema fields and exits zero
status: ideation
source: "email-triage field report 2026-08-26: historical batches carry retired gate-schema fields (current, digest-domain — the latter absent from today's source), so --validate exits non-zero forever on any workflow with history; captain ruling: no migration, just warn"
started: 2026-08-26T22:00:13Z
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:wyzg6knr8whtmg79mxkb78jg:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:wyzg6knr8whtmg79mxkb78jg-backlog-1
              briefing:
                id: briefing:wyzg6knr8whtmg79mxkb78jg:backlog:attempt-1:revision-1
                digest: sha256:ae90a8befc3588f499f77645840d99414c911362ce1184fe1bb180afe01575a1
                room-ref: ./validate-warns-on-retired-schema-fields/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:wyzg6knr8whtmg79mxkb78jg:backlog:1
                briefing: briefing:wyzg6knr8whtmg79mxkb78jg:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-26T21:59:45.808528Z"
                decision: approve
              application:
                target-stage: ideation
                state: consumed
---

A workflow with history cannot pass `status --validate` because entities written under an old gate schema carry fields the current schema retired. A check that is always red cannot fail meaningfully. Captain ruling: do not build a migration. Downgrade retired-field findings to warnings and exit zero when they are the only findings.

## Problem

The canonical gate model deliberately uses `yaml.Decoder.KnownFields(true)`. That
keeps authority fail-closed, but it also rejects two fields that Spacedock itself wrote
before commit `f566f821b` simplified v1: top-level `gates.current` and
`gates.records[*].attempts[*].briefing.digest-domain`. A workflow that retains those
otherwise-valid records therefore gets a permanent `Error: invalid gates` and exit 1
from `status --validate`. The operator can no longer distinguish real current-schema
corruption from clean immutable history.

The fix must restore that validation signal without making an old selector authoritative,
accepting arbitrary unknown keys, or rewriting stored entities. This is a read-only
compatibility boundary: current workflow `status` plus ordered records/attempts remain the
only gate selector, and every canonical writer continues to emit today's schema.

## Proposed approach

Extend the existing clone-before-decode compatibility pass in `internal/gates/io.go`.
Classify compatibility by exact YAML location and historical shape, remove accepted
retired nodes only from the validation clone, and return typed diagnostics while keeping
the original node intact for compare-and-swap and byte preservation.

The explicit classification is:

| Class | Exact stored location and accepted historical shape | Diagnostic |
| --- | --- | --- |
| retired-warning | `gates.current`: a mapping containing only scalar `gate`, whose nonblank value names a record in the same document | `Warning: retired gate field 'current' at gates.current: ...` |
| retired-warning | `gates.records[*].attempts[*].briefing.digest-domain`: scalar value exactly `canonical-bytes` | `Warning: retired gate field 'digest-domain' at gates.records[N].attempts[N].briefing.digest-domain: ...` |
| retired-silent (existing, unchanged) | `gates.records[*].attempts[*].provider-evidence` | no diagnostic; it retains its existing frozen-provider compatibility behavior |
| application-extension (existing, unchanged) | unknown keys only inside an exact `application` mapping | existing unknown-application warning |
| strict | every other unknown field, wrong location, duplicate retired key, or malformed retired shape/value | `Error: invalid gates: ...`, exit 1 |

The retired-warning set is a closed code-owned table, not a key-name substring match and
not a permissive YAML tag. Adding or changing a gate field does not enter the table
automatically: the strict decoder rejects it until a deliberate code, fixture, and
contract-doc change adds its exact location and historical shape. That is how evolution
fails closed.

`internal/status/validate.go` renders the new diagnostic class with the existing entity
evidence fields and deterministic path order. On explicit `status --validate`, an active
entity with only retired findings writes warnings to stderr, writes `VALID` to stdout,
and exits 0. With `--json`, stdout remains exactly
`{"command":"validate","valid":"true"}` and the same warnings remain on stderr.
Ordinary status/read commands accept and ignore the retired nodes without printing them.
The existing publish-only archive policy also remains: archived retired nodes are
readable but do not produce archived-scope warnings. No command grammar changes.

The retired nodes never populate `gates.Document`, never affect readiness or mutation
authority, and are never emitted by a canonical gates writer. `status --validate` is
read-only and preserves the complete entity bytes. There is no scan, migration, schema
version, or background rewrite.

Smallest alternative considered: catch `yaml.v3` decoder error strings containing
`current` or `digest-domain` and downgrade them in `status`. It is insufficient because
the decoder stops at the first unknown field, the error text does not provide a safe
exact-location classifier, and the approach could downgrade a same-named field in the
wrong subtree while still leaving gate reads unusable. Reusing the existing filtered-clone
seam is the smallest mechanism that serves AC-1 through AC-4 without weakening strict
decode.

## Risk evidence

The 2026-08-26 email-triage field report observed the permanent nonzero validation over
historical batches. Repository history independently identifies the pair: commit
`f566f821b` removed `Document.Current` and `Briefing.DigestDomain` together while making
status plus record order authoritative. Current `TestPrototypeAndUnknownGateShapesFailClosed`
and `TestDigestDomainFieldFailsClosed` exercise the red baseline.

No spike needed: the required mechanism is already proven by
`TestRetiredProviderEvidenceKeyReadsSilentlyWithoutWideningAttemptTolerance`,
`TestReadDiagnosticsFiltersOnlyExactApplicationMappings`, and
`TestFieldConformanceWarnsDoNotGateExit`. A focused 2026-08-26 run of those seams passed;
the implementation changes their classification inputs, not parser, exit-code, or storage
machinery.

## Out of scope

Any migration or rewrite of historical entities; a second gate-schema version; tolerance
for `current-attempt`, `sequence`, `previous-attempt`, explicit attempt `state`, arbitrary
unknown fields, or malformed legacy values; changing the existing silent
`provider-evidence`, application-extension, field-conformance, flat-room, or archived
warning policies; and changing gate selection, readiness, or writer authority.

## Expected surface and tolerance

Estimate net LOC change: +145, across 6 files. Expected gross composition is about +160
insertions and -15 deletions: compatibility classification in `internal/gates/io.go`,
diagnostic rendering in `internal/status/validate.go`, focused gate and CLI fixtures in
`internal/gates/gates_test.go` plus one new `internal/status/*_test.go`, and contract wording
in `docs/specs/gate-resolution-frontmatter-contract.md` and
`docs/site/reference/frontmatter-contract.md`.

Tolerance: net LOC may vary by +/-45 and the surface by +1 file. Test helper extraction or
renaming the compatibility pass is within tolerance; a migration command, new schema
version, changed CLI grammar, stored-format writer change, archive-policy change, or a
generic unknown-field allowlist is outside it regardless of LOC.

Observable semantics: command grammar is unchanged; stored formats are unchanged and no
bytes are written by validation; gate authority is unchanged; runtime behavior changes
only so exact legacy nodes are non-authoritative readable compatibility, active explicit
validation reports them as warnings, and retired-only validation exits 0.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - A workflow containing both exact historical retired fields changes from the main-branch baseline's exit 1 to exit 0, while retaining two field-and-path-specific warnings and `VALID`/`valid:"true"` output.**
Verified by: a CLI fixture runs text and JSON `status --validate` against an active entity
with valid current-schema gate data plus `gates.current.gate` and
`briefing.digest-domain: canonical-bytes`; it asserts exit 0, exact warning count/classes,
stdout grammar, and no `Error:`. Removing either classifier or routing warnings into the
error slice makes the fixture fail.

**AC-2 (BOUNDARY) - A current-schema defect remains fatal even when the same entity carries the retired pair.**
Verified by: a sibling CLI fixture changes the canonical `briefing.digest` to a non-SHA-256
value and asserts exit 1, `valid:"false"` in JSON, and `Error: invalid gates`. Broadly
dropping the Briefing or treating all gate diagnostics as warnings makes it fail.

**AC-3 (BOUNDARY) - Compatibility is limited to the two exact historical encodings; unrelated unknown fields and corrupted legacy lookalikes still fail closed.**
Verified by: a gate-reader table covers a same-named key at the wrong location, typo
`digest-domains`, non-`canonical-bytes` digest domain, `current` with an extra key, and a
`current.gate` that names no record. Each must return an error, while the exact pair yields
typed sorted warnings and a canonical document with status/order authority. Replacing the
classifier with name-only filtering makes at least one negative turn green and the test
fail.

**AC-4 (STORAGE/AUTHORITY) - Validation and ordinary reads never rewrite retired bytes or use them as gate authority.**
Verified by: the reader fixture hashes the entity before and after diagnostic and ordinary
reads, asserts byte equality, and plants a well-formed but status-stale `current.gate`; the
projected readiness must follow entity `status` and ordered attempts. Filtering the source
node or restoring pointer selection makes it fail.

**AC-5 (DOCUMENTATION) - The public and internal frontmatter contracts distinguish exact retired compatibility from arbitrary unknown-field rejection.**
Verified by: review the concrete diff below against the implementation and run the docs
build already used by the repository; the wording must name both locations, warn-only
active `--validate`, non-authority, no rewrite, and strict fallback.

## Test plan

1. Add gate-package table tests for exact retired classification, deterministic paths,
   preserved source nodes/bytes, status-based authority, and the AC-3 strict negatives.
   Cost: medium; in-memory/temp-file Go fixtures, no live workflow.
2. Add status-package command fixtures for retired-only text and JSON branches plus the
   mixed invalid-digest control. Assert output bytes, stderr class/count, and exit codes.
   Cost: medium; native CLI fixture, no external runtime.
3. Run `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` as required.
   The full suite guards unknown application compatibility, archived warning scope,
   corruption, mutation atomicity, and golden output beyond the focused fixtures.
4. Run `mkdocs build --strict` after applying the approved wording, matching the docs CI
   command in `.github/workflows/docs.yml`. Cost: low; docs-only check, no live workflow.

Concrete documentation change for `docs/specs/gate-resolution-frontmatter-contract.md`
and the corresponding paragraph in `docs/site/reference/frontmatter-contract.md`:

```diff
-Unknown or prototype fields inside binary-owned `gates` fail closed except unknown keys
-inside an exact `application` mapping; those keys warn on explicit `status --validate`.
+Unknown or prototype fields inside binary-owned `gates` fail closed. Read-only
+compatibility accepts only the historical `gates.current: {gate: <existing-record-id>}`
+mapping and `records[*].attempts[*].briefing.digest-domain: canonical-bytes` scalar.
+They never supply authority or appear in canonical writes. Explicit `status --validate`
+warns for them on active entities but remains valid and exits zero; ordinary reads and
+archived warning scope stay silent, and validation never rewrites stored bytes. Unknown
+keys inside an exact `application` mapping retain their existing warning-only behavior.
```

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}
