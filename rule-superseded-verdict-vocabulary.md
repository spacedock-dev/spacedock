---
id: x2ezetxr82pztr4pqt1g4dhx
title: Rule the superseded verdict into or out of the schema vocabulary
status: implementation
source: "scope-validate-warnings ideation, 2026-08-15: 4 archived entities carry verdict superseded, a token the conventional enum [PASSED REJECTED] never admitted"
started:
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:x2ezetxr82pztr4pqt1g4dhx:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:x2ezetxr82pztr4pqt1g4dhx-backlog-1
              briefing:
                id: briefing:x2ezetxr82pztr4pqt1g4dhx:backlog:attempt-1:revision-1
                digest: sha256:d361d63a1849a995f40b359b0a666f56642c0c060b80186f813f2d6c235f23f4
                request-digest: sha256:3b81644f231329651f19ea0d696ce81ad85c558246a98dea5bdc85ad77b93f85
                room-ref: ./rule-superseded-verdict-vocabulary/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:x2ezetxr82pztr4pqt1g4dhx:backlog:1
                briefing: briefing:x2ezetxr82pztr4pqt1g4dhx:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T21:24:45.264332Z"
                decision: approve
                reason: 'Captain directive 2026-08-15: dispatch all five onto the stack tip'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:x2ezetxr82pztr4pqt1g4dhx:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:x2ezetxr82pztr4pqt1g4dhx-ideation-1
              briefing:
                id: briefing:x2ezetxr82pztr4pqt1g4dhx:ideation:attempt-1:revision-1
                digest: sha256:1cbd2f29086978446c2b292eb966aa96fe310c819eaf0a673d3fff578a16ebd1
                request-digest: sha256:1205c00481df70ef3f4acf7f2cf898fbbd8defb65d456206248c99413e62c2c1
                room-ref: ./rule-superseded-verdict-vocabulary/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:x2ezetxr82pztr4pqt1g4dhx:ideation:1
                briefing: briefing:x2ezetxr82pztr4pqt1g4dhx:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T23:03:59.955713Z"
                decision: approve
                reason: 'Captain batch approval 2026-08-15 (approve all): into implementation as stack layers'
              application:
                target-stage: implementation
                state: consumed
---

Writers intentionally emit verdict: superseded for superseded entities, but the schema enum admits only PASSED and REJECTED. Archived-scope warnings are silenced now, so the bite is forward-looking: the next active entity superseded on purpose warns as invalid, and any tool trusting the enum misreads the four archived records. Decide: admit superseded (and define its semantics) or route supersede through a different field and stop the writer.

## Problem

`verdict: superseded` is a dead June-2026 convention with a live bite.

Four archived entities carry it: `survey-scaffold-state-the-fact`,
`survey-agentsview-detect-under-sandbox`, `survey-codex-cwd-workaround` (all
completed 2026-06-11) and `read-hint-adoption-bloat-trim` (2026-06-16). Nothing
has written it since. The 13 entities superseded afterwards — `e7426315b` (12
live-test-truth members absorbed, 2026-08-03) and `bf592d776`
(`restore-quarantined-common-live-journeys`, 2026-08-09) — were moved to
`_archive` with `verdict:` left EMPTY, the reason carried in body prose and the
commit message. The convention did not survive; only the four records did.

The value chain is broken at every link:

- **No consumer.** Nothing branches on the token. The only verdict predicate is
  `isRejectedVerdict` (`internal/status/verdict.go:58`), a merge-ceremony
  carve-out with three call sites (`mutate.go:356`, `mutate.go:378`,
  `handlers.go:292`). `internal/journeymetrics` never reads `verdict` at all.
- **No writer.** The tool-supported terminal writer refuses it: `merge guard
  --verdict` takes only `passed|rejected` (`internal/status/merge.go:60`). The
  four records were hand-written through `status --set`, which validates nothing.
- **No verifier.** The schema enum is `[PASSED REJECTED]`
  (`docs/schema/entity.mdschema.yml`), so the token warns on every read.

Writing it on an active entity is worse than a warning — it strands the entity.
Measured at the stack tip (spike record below): the token produces one
`status --validate` warning AND makes `--archive` exit 1, because the archive
guard (`mutate.go:378`) reads every non-empty, non-`rejected` verdict as
approval-style and demands terminal status or the merge ceremony. The field
written to say "this was superseded" is exactly what blocks the superseding.

The gap that produced the four records is still open, and it is not
superseded-specific: `status --set {slug} verdict=banana` is accepted, exit 0.

## Proposed approach

**Ruling: `superseded` stays OUT of the verdict vocabulary.** The enum remains
`[PASSED REJECTED]`. The supported way to supersede an entity is what the August
practice already does and the binary already accepts: leave `verdict` empty and
run `spacedock status --archive {slug}`, recording why in the body and the commit
message. Proven at the tip — it archives clean and validates clean.

One thing changes. Close the write gap at the single normalisation point every
frontmatter write already passes through. `canonicalConventional`
(`internal/status/verdict.go:18`) consults the schema's conventional list on
write today, but only to fold case (`verdict=passed` → `PASSED`). Extend that
same lookup to admission:

> A conventional list is advisory ON READ and closed ON WRITE. History keeps
> whatever token it carries — no archived record is ever edited — but a NEW
> write must pick a token from the list.

Concretely, `status --set {slug} {field}={token}` refuses, byte-clean and
non-zero, when the field declares a `conventional` list and the token matches
none of it case-insensitively. The refusal names the field, the token, and the
allowed set. Clearing (`{field}=`) always passes — that is the supported
supersede shape. `--force` bypasses, the uniform escape the other five `--set`
guards already honor (`handlers.go` lines 181, 223, 246, 282, 292).

The guard belongs in `runSet` beside those guards, not inside
`updateFrontmatter`: `force` is in scope there, the refusal lands before any
mutation (byte-clean by construction), and the shared write engine stays a pure
normaliser for finalize and archive, which pass already-canonical values.

### Keep boundaries

- **`application.state: superseded` in `internal/gates` is a different field**
  with its own vocabulary (`pending|consumed|superseded`, `gates/model.go:334`)
  and its own guarded transition. This ruling does not touch it. The name
  collision is the main way to get this task wrong.
- **The four archived records are not edited.** Archived scope is publish-only,
  and at the stack tip field conformance already skips it
  (`field_conformance.go:94`, from `scope-validate-warnings-to-active-entities`).
  They keep `verdict: superseded` and emit nothing.
- **`verdict` is the only field carrying a `conventional` list today**, so the
  new refusal's blast radius is exactly the field being ruled on. The mechanism
  is field-generic because the schema is the SSOT, not because other fields are
  in scope.
- **`merge guard --verdict` is unchanged.** It already refuses non-conventional
  tokens; this task adds the matching refusal on the other writer.

### Alternatives considered

- **Admit `superseded` into the enum.** Insufficient, and costlier than it looks.
  Admitting the token without also carving it out of the archive and `--set`
  guards leaves the stranding bug in place — the entity still cannot be archived.
  The honest version of this arm therefore generalises `isRejectedVerdict` into a
  "did not ship" predicate across its three call sites, extends `merge guard
  --verdict`, and adds a third semantic branch, all to encode a distinction no
  code reads and that prose already carries. It buys vocabulary, not behavior.
- **Document the ruling and enforce nothing.** Rejected by AC-2 and by the
  record: the four entities exist *because* the writer accepts anything. A rule
  the writer does not hold is the state we are already in.
- **Harden the schema instead — set `verdict.invalid_severity: error`.** Actively
  harmful, and a trap this task must not fall into. Exercised at the tip:
  flipping `warn`→`error` makes `isWarnSeverity`
  (`field_conformance.go:116`) return false, the field is skipped entirely, and
  validate prints `VALID` with no warning. No error path for field conformance
  exists, so the check does not harden — it vanishes. That would satisfy a naive
  reading of AC-1 ("zero warnings") while making the surface strictly worse,
  which is why AC-3 pins the negative control.
- **Delete the `conventional` list.** The same silencing, stated openly rather
  than by accident. It also discards the case-folding the write path depends on.

### Spike record

Exercised against a binary built from the stack tip `origin/main` 697ffbf43
(local `main` is 10 commits behind and does NOT carry the archived-scope skip).
Fixture: `internal/status/testdata/enum-scope-workflow`.

1. **Baseline warning.** Active entity carrying `verdict: superseded` →
   `status --validate` prints exactly one line,
   `Warning: field 'verdict' value "superseded" is not one of [PASSED REJECTED] … scope=active`,
   exit 0 (`VALID`). The archived twin carrying the same token prints nothing.
2. **Baseline stranding.** The same active entity at non-terminal stage
   `backlog` → `status --archive` exits **1**: `cannot be archived — verdict
   'superseded' is set but status 'backlog' is not the terminal stage`. With
   `verdict:` empty, the same archive exits 0 and stamps `archived:`.
3. **Writer gap is general.** `--set verdict=superseded` → exit 0, written.
   `--set verdict=banana` → exit 0, written. `--set verdict=passed` → exit 0,
   stored as `PASSED`, confirming the write path already consults the list.
4. **Silencing trap is real.** With `invalid_severity: error` on `verdict` and a
   rebuild, the same active fixture validates `VALID` with no warning at all.

Nothing in the proposed approach rests on an unexercised mechanism: the schema
lookup, the write path, and the `--force` escape are all existing, exercised
code. The spike seeds the implementation's first test (AC-2).

### Documentation

The refusal is user-visible CLI behavior, so two doc sites change.

`docs/dev/README.md:72` — before:

    | `verdict` | enum | PASSED or REJECTED - set at final stage |

after:

    | `verdict` | enum | PASSED or REJECTED - set at final stage. Closed on write: `--set` refuses any other token (`--force` bypasses). To supersede an entity, leave `verdict` empty and `--archive` it, recording why in the body. |

`skills/fo-status-viewer/SKILL.md:31-34` — before:

    The `--set` flag updates entity frontmatter fields:
    - `--set {slug} field=value` sets a field
    - `--set {slug} field=` clears a field
    - `--set {slug} started` or `completed` auto-fills a UTC ISO 8601 timestamp (skipped if already set)

after:

    The `--set` flag updates entity frontmatter fields:
    - `--set {slug} field=value` sets a field
    - `--set {slug} field=` clears a field
    - `--set {slug} started` or `completed` auto-fills a UTC ISO 8601 timestamp (skipped if already set)
    - a field with a schema `conventional` list (today: `verdict`, `[PASSED REJECTED]`) is closed on write — any other token is refused byte-clean; `--force` bypasses. Clearing always passes. To supersede an entity, clear `verdict` and `--archive` it.

`internal/cli/help.go` is deliberately NOT changed: the `--set` help does not
enumerate field values, and the refusal message names the allowed set itself.

## Out of scope

- Hand-editing archived entities (publish-only, per 6c45fd59c precedent).
- `application.state: superseded` in `internal/gates` — a different field (see
  Keep boundaries).
- Hardening any field other than `verdict`. The guard is schema-driven, so a
  future field that declares a `conventional` list inherits it; adding such a
  field is not this task.
- Changing `merge guard --verdict`, which already refuses.

## Expected surface and tolerance

Baseline for the diff is the stack tip `origin/main` 697ffbf43, not local `main`.

| File | Expected change |
| --- | --- |
| `internal/status/verdict.go` | +~18 — `conventionalViolation(field, value)` helper beside `canonicalConventional`, same schema lookup |
| `internal/status/handlers.go` | +~14 — the guard in `runSet`, before any mutation, `--force` bypass |
| `docs/schema/entity.mdschema.yml` | +1 — a `semantics` key on `verdict` recording advisory-on-read / closed-on-write |
| `docs/dev/README.md` | ~1 changed line (verdict row) |
| `skills/fo-status-viewer/SKILL.md` | +1 bullet |
| `internal/status/verdict_vocabulary_test.go` | +~90 — new file, AC-2/3/4 tests |

**Estimate: 6 files, ~+125 / -2 lines.** Tolerance: ±2 files, ±45 lines. A
correction round that pushes past this recalibrates at the gate rather than
absorbing the overrun.

### Declared semantic changes

- **Runtime behavior — CHANGED (the one deliberate change).** `status --set
  {slug} {field}={token}` now exits non-zero and writes nothing when the field
  declares a `conventional` list and the token is outside it. `--force`
  preserves the old behavior. Today this affects `verdict` only.
- **Command grammar — unchanged.** No new flags, no new subcommands.
- **Stored formats — unchanged.** The enum is not extended; no entity, active or
  archived, is rewritten.
- **Authority — unchanged.** No gate, approval, or consumption path is touched.

## Acceptance criteria

**AC-1 (value) - An intentional supersede of an active entity completes with zero
validate warnings and zero archive refusals.**
Verified by: a fixture active entity superseded through the supported path
(`verdict` cleared, then `--archive`) — `status --validate` emits zero warnings
and `--archive` exits 0. Independent baseline that moves the wrong way, measured
at the tip: the same entity carrying `verdict: superseded` emits one warning and
`--archive` exits 1. The pair (warning count 1→0, archive exit 1→0) is the
end-value; either number regressing fails the AC.

**AC-2 - The ruling is enforced where verdicts are written, not only documented.**
Verified by: `status --set {slug} verdict=superseded` exits non-zero and leaves
the entity file byte-identical; `--set verdict=` (clear) and `--set
verdict=passed` (stored `PASSED`) still succeed; `--force` writes the
non-conventional token. Fails if the guard is absent, mutates before refusing, or
breaks the existing case-fold.

**AC-3 (negative control) - A non-conventional verdict on an active entity still
warns.**
Verified by: an active fixture carrying `verdict: banana` produces exactly one
warning naming `[PASSED REJECTED]`. This fails if AC-1 is "achieved" by flipping
`invalid_severity` to `error`, by deleting the `conventional` list, or by
otherwise silencing field conformance — the cheat the spike proved is available.

**AC-4 - The four archived records are untouched and silent.**
Verified by: the four archived entities still carry `verdict: superseded`
byte-for-byte in the state checkout, and an archived fixture carrying the token
produces zero warnings.

**AC-5 - The suite stays green.**
Verified by: `go test ./...` and `go test ./... -race`.

## Test plan

New file `internal/status/verdict_vocabulary_test.go`, driving the existing
`enum-scope-workflow` fixture (extended with the entities each case needs).
Table-driven over the `--set` cases so a new conventional field is one row.

| Test | Claim | Fails when |
| --- | --- | --- |
| `TestSetRefusesNonConventionalToken` | AC-2 | the guard is missing, or it mutates then errors (asserts file bytes unchanged, not just exit code) |
| `TestSetForceWritesNonConventionalToken` | AC-2 | `--force` stops being the uniform escape |
| `TestSetClearAndCaseFoldStillPass` | AC-2 | the guard swallows `verdict=` or breaks `passed`→`PASSED` |
| `TestSupersedePathIsWarningFree` | AC-1 | the supported path regains a warning, or `--archive` starts refusing an empty-verdict entity |
| `TestSupersededVerdictBaselineWarnsAndBlocks` | AC-1 baseline | the baseline stops being measurable — pins warning count 1 and archive exit 1 for the bad token |
| `TestNonConventionalVerdictStillWarns` | AC-3 | field conformance is silenced by severity flip or list deletion |
| `TestArchivedScopeStaysSilent` | AC-4 | archived scope starts warning again |

Cost: one new Go test file, fixture-level, no live workflow run — the claims are
all command-level behavior (exit code, stdout/stderr bytes, on-disk state), which
the existing `internal/status` harness already drives. A separate one-line check
that the four archived records still read `verdict: superseded` covers AC-4's
state half.

## Stage Report: ideation

- DONE: Rule superseded into or out of the verdict vocabulary with writer-side enforcement, archived records untouched
  Ruled OUT: enum stays `[PASSED REJECTED]`; enforcement is a schema-driven admission check on the `--set` write path (`--force` bypasses); the four archived records are never edited and are already silent at the tip.
- DONE: Value AC: an intentional active-entity supersede produces zero validate warnings under the ruling; baseline today warns
  AC-1 states the value as a two-number delta measured at the tip — warning count 1→0 AND archive exit 1→0 — after the spike found the token also *blocks* `--archive`, not just warns. AC-3 adds the negative control that stops AC-1 being met by silencing the checker.
- DONE: Design against the stack tip
  Built and exercised `origin/main` 697ffbf43; local `main` is 10 commits behind and lacks the archived-scope skip (`13faf8f55`) the "archived untouched" boundary rests on. The surface estimate is baselined against the tip.

### Summary

The ruling is OUT, on necessity: `verdict: superseded` has no consumer (nothing
branches on it; `journeymetrics` never reads `verdict`), no current writer
(`merge guard --verdict` refuses it, and the August supersede practice archives
with `verdict` empty), and no verifier (the enum warns). The only real change is
closing the writer gap that produced the four records — `--set` accepts any
token, `verdict=banana` included — by extending the schema lookup the write path
already performs for case-folding into an admission check.

Two spike findings shaped the design. First, `verdict: superseded` on an active
entity does not merely warn: `--archive` exits 1, because the archive guard reads
any non-empty non-`rejected` verdict as approval-style. The field written to
record a supersede is what blocks it. Second, setting `invalid_severity: error`
makes the conformance check vanish rather than harden — `isWarnSeverity` returns
false and the field is skipped — so a naive "zero warnings" AC is satisfiable by
silencing the checker. AC-3 exists to close that hole.

The main way to get this task wrong is the name collision with
`application.state: superseded` in `internal/gates`, an unrelated field with its
own vocabulary; it is called out in Keep boundaries and Out of scope.
