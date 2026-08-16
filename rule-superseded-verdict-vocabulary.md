---
id: x2ezetxr82pztr4pqt1g4dhx
title: Rule the superseded verdict into or out of the schema vocabulary
status: validation
source: "scope-validate-warnings ideation, 2026-08-15: 4 archived entities carry verdict superseded, a token the conventional enum [PASSED REJECTED] never admitted"
started: 2026-08-15T23:04:19Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-rule-superseded-verdict-vocabulary
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
        - id: gate:x2ezetxr82pztr4pqt1g4dhx:validation
          stage: validation
          attempts:
            - id: gate-attempt:x2ezetxr82pztr4pqt1g4dhx-validation-1
              briefing:
                id: briefing:x2ezetxr82pztr4pqt1g4dhx:validation:attempt-1:revision-1
                digest: sha256:0fc807705d9edbd49705b920dd987f30dbca59b0f533fd5972ade83ff1209bfb
                request-digest: sha256:5b4a630193bd5a992262a70ab7942157629f11c3890b27bee7a5665a3ff63e7d
                room-ref: ./rule-superseded-verdict-vocabulary/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:x2ezetxr82pztr4pqt1g4dhx:validation:1
                briefing: briefing:x2ezetxr82pztr4pqt1g4dhx:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-16T02:22:29.49511Z"
                decision: approve
                reason: 'Captain 2026-08-16 (approve all to the stack): validation PASSED; surface overrun accepted as AC-bearing test density; lands as a stack layer'
              application:
                target-stage: done
                state: pending
pr: "#716"
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

## Stage Report: implementation

- DONE: Execute the gated design exactly: superseded stays OUT; --set gains the schema-lookup admission check; archived records untouched
  `64c037b56`: enum unchanged; `conventionalViolation` (verdict.go) reuses `canonicalConventional`'s schema lookup for admission; guard in `runSet` immediately before the write, `--force` bypasses. No archived record edited — `git status` on the state checkout's `_archive/` is empty.
- DONE: Value AC both numbers: active-entity supersede warning 1 to 0 AND --archive exit 1 to 0; AC-3 negative control intact (severity-hardening must not silence the check)
  Both numbers moved, exercised end-to-end with the built binary: unsupported token → warnings 1, archive exit 1; supported path (clear `verdict`, then `--archive`) → warnings 0, archive exit 0.
- DONE: Base is the stack top ab112b6d3; land as the next available layer per the ratified Stacked mode
  `64c037b56` is `ab112b6d3` + 1 on `spacedock-ensign/rule-superseded-verdict-vocabulary` — the same one-layer-off-the-tip shape all four peer branches hold. Layer commit left unpushed, matching peers (only `tune-dev-template-for-gated-stages` has a remote, still at `ab112b6d3`); the FO owns merge sequencing.

Tests — `internal/status/verdict_vocabulary_test.go`, 7 tests, grouped by claim:

- **AC-2, the ruling is enforced where verdicts are written.** `TestSetRefusesNonConventionalToken` (3 rows: `superseded`, `banana`, `SUPERSEDED`) asserts exit non-zero AND the entity bytes unchanged; `TestSetForceWritesNonConventionalToken`; `TestSetClearAndCaseFoldStillPass` (5 rows). Fails if the guard is absent, mutates before refusing, loses the `--force` escape, swallows `verdict=`, breaks `passed`→`PASSED`, or reaches a field with no `conventional` list. The byte assertion is what a bare exit-code check would miss.
- **AC-1, both numbers.** `TestSupersededVerdictBaselineWarnsAndBlocks` pins the wrong-way pair (warning count exactly 1, archive exit exactly 1); `TestSupersedePathIsWarningFree` pins 0 and 0 through the supported path. Fails if either number stops being what it is — which is what keeps 1→0 an end-value rather than a tautology.
- **AC-3, negative control.** `TestNonConventionalVerdictStillWarns` requires exactly one warning naming `[PASSED REJECTED]`. Setting `verdict.invalid_severity: error` (which makes `isWarnSeverity` false and skips the field entirely) or deleting the `conventional` list drops the count to 0 and fails it — the spike's proven cheat is closed.
- **AC-4, archived untouched and silent.** `TestArchivedScopeStaysSilent` requires zero warnings from an archived entity carrying the token. State half verified directly: exactly the four named records still carry `verdict: superseded`, no active entity does.
- **AC-5.** `go test ./...` and `go test ./... -race` both green on the final tree; `gofmt -l` and `go vet` clean.

### Summary

The ruling is OUT, implemented as declared: the enum is unchanged and the only behavior change is that a schema `conventional` list, already consulted on write to fold case, is now also consulted for admission — advisory on read, closed on write, clearing always passes, `--force` bypasses.

Two deviations to weigh at the gate. First, surface overran: **8 files, +365/-12 against the declared 6 files / ~+125/-2 (tolerance ±2 files, ±45 lines)**. Files are at the tolerance edge; lines are ~195 past it. The overrun is concentrated in one place — the new test file is +281 against the ~+90 estimated — because the estimate implied ~13 lines per test while this package's actual per-test density (see `live_guard_test.go`, `verdict_case_test.go`) runs 25-40 lines including the doc comment. Production code landed near estimate (verdict.go +39 vs ~+18, handlers.go +28 vs ~+14, docs on target). I did not cut AC-bearing tests to fit the number; recalibrate the estimate or tell me which claim to drop.

Second, the declared semantic change had two consequences the ideation did not anticipate, both in the existing test corpus and neither altering an assertion: `live_guard_test.go` used `verdict=accepted` at four sites purely as a non-empty token while exercising the live-run guard (retargeted to `verdict=passed`; those tests assert `status: done`, never the verdict value), and `verdict_case_test.go`'s `needs-work` case already carried `--force` so it still passes — its comment gained the closed-on-write half. A corpus sweep found no other non-conventional `verdict=` writer.

One incidental finding worth recording: `--set` refuses archived entities outright (`archived entity is read-only`), so AC-4's "never edited" boundary has a second, independent layer beneath this guard.

## Stage Report: validation

- DONE: Independently re-exercise the ACs against worktree commit 64c037b56, never by reading the report: active-entity supersede warning 1-to-0 AND archive exit 1-to-0; the --set admission check refuses out-of-enum tokens; AC-3 negative control (severity-hardening must not silence the check); archived records untouched
  All reproduced with a binary built at 64c037b56 driving scratch copies of the enum-scope-workflow fixture; per-AC citations below.
- DONE: Rule on the declared surface overrun: 8 files +365/-12 vs estimate 6 files ~+125/-2 - the overrun is test-density accounting per the implementer; judge whether the tests are AC-bearing or padding
  Ruling: AC-bearing, not padding. The new file runs 40 lines/test, inside the package norm (live_guard 36, verdict_case 28, field_conformance 41); the ~+90 estimate implied ~13 lines/test, which no test file in this package achieves. Each of the 7 tests pins a distinct falsifiable claim — three falsified live (below) — and no two tests prove the same claim. Recommend recalibrating the estimate at the gate, not trimming tests. AC unchanged, not narrowed.
- DONE: Suites green plain and -race; verdict PASSED or REJECTED with per-AC citations
  `go test ./...` exit 0 AND `go test ./... -race` exit 0 at 64c037b56; `gofmt -l` empty, `go vet` clean. Verdict: PASSED.

Per-AC evidence, all independently reproduced:

- AC-1 PASS: active fixture with `verdict: superseded` → exactly 1 validate warning, `--archive` exit 1; after `--set verdict=` → 0 warnings, `--archive` exit 0, file lands in `_archive/`. Both numbers moved 1→0.
- AC-2 PASS: `--set verdict=superseded|banana|SUPERSEDED` each exit 1 with the entity byte-identical (cmp against a pre-copy), refusal naming field, token, and [PASSED REJECTED]; multi-field `status=done verdict=banana` refuses byte-clean too (no half-write); `verdict=` clears; `passed`→`PASSED`, `rejected`→`REJECTED` fold intact; `--force` writes the token verbatim; a field with no conventional list (`source`) is unguarded.
- AC-3 PASS: active `verdict: banana` → exactly 1 warning naming field, token, allowed set, slug. Falsified on a throwaway export of 64c037b56 (never the worktree): `invalid_severity: error` → TestNonConventionalVerdictStillWarns fails at 0 warnings; deleting the conventional list → that test plus both fold rows fail; reverting the runSet guard to ab112b6d3 → all three refusal rows fail. The tests cannot pass while the behavior is wrong.
- AC-4 PASS: exactly the four named `_archive/` records carry `verdict: superseded`; last commits touching them are the original June archives (1426ca688 2026-06-11, ca8e7706e 2026-06-16) and `git status` over `_archive/` is clean; an archived fixture carrying the token validates silent; independent second layer confirmed — `--set` on an archived entity refuses (`archived entity is read-only`).
- AC-5 PASS: full plain and -race suites, exit 0 each, zero FAIL lines.

New finding for Review-finding disposition (recommendation only; FO authorizes): whitespace-only `--set "verdict=   "` passes the admission guard (trimmed to empty, treated as a clear) but stores literal `verdict: '   '`, which then warns on read as non-conventional. No stranding: `--archive` still exits 0, and a padded real token (`" passed "`) trims and folds to `PASSED` cleanly, so the inconsistency is confined to whitespace-only values. Recommended class: Deferred risk. Evidence fields: (1) released user/normal workflow — no supported writer emits a quoted whitespace-only verdict; the documented clear is bare `verdict=`; (2) observable harm — one spurious validate warning, no archive block, no data loss; (3) authority — none: no value AC fails (AC-1's supported path verified passing), so Material cannot be established; (4) trigger evidence — this session at 64c037b56: set exit 0, stored `'   '`, warning count 1, archive exit 0. Promote-to-material condition: a supported scripted writer that can interpolate whitespace-only values into `verdict=`, or the archive guard stops treating whitespace as empty (which would recreate the stranding AC-1 removes).

### Summary

PASSED, all five ACs independently reproduced against a binary built at worktree commit 64c037b56 — never from the implementation report — with the AC-1 value pair (warning 1→0, archive exit 1→0) observed end-to-end and AC-2/AC-3 falsified three ways on a throwaway checkout to confirm the tests have teeth. The surface overrun is ruled AC-bearing test density, not padding or scope creep: production code and docs land near estimate, the delta is per-test line count at the package's own norm, and the two unanticipated test-corpus touches (live_guard retarget `accepted`→`passed`, verdict_case `--force` comment) are consequences of the declared semantic change with no assertion altered. One deferred risk recorded above (whitespace-only verdict value); no material findings; recommend the gate proceed and recalibrate the estimate.
