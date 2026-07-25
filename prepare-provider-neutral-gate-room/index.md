---
id: s4ykctf21g60dvfgdd6cy9ny
title: Prepare provider-neutral gate rooms and align canonical Briefing recording
status: ideation
source: "Durable-decisions cross-repo dogfood ruling after xb and Subspace em review, 2026-07-24"
started: 2026-07-24T14:54:10Z
completed:
verdict:
score: "1.0"
worktree:
issue:
sprint: durable-decisions
gates:
    version: 1
    current:
        gate: gate:docs-dev:s4:ideation
    records:
        - id: gate:docs-dev:s4:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:s4-backlog-1
              briefing:
                id: briefing:docs-dev:s4:backlog:attempt-1:revision-1
                digest: sha256:8d6888f2f9d067835f24c8845d703547638ff919f71f709c681e856551cfb80f
                digest-domain: canonical-bytes
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:s4:backlog:1
                briefing: briefing:docs-dev:s4:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-07-24T14:53:41.564011Z"
                decision: approve
                reason: Captain approved filing the narrow post-em alignment task and dispatching it through the sprint.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:docs-dev:s4:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:s4-ideation-1
              briefing:
                id: briefing:docs-dev:s4:ideation:attempt-1:revision-1
                digest: sha256:2185e46203b2e4747d8a7db557f14737492cf519c9c7d40b632dfe83c51a0074
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:s4:ideation:1
                briefing: briefing:docs-dev:s4:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-24T15:41:06.723937Z"
                decision: approve
                reason: The corrected design removes manual room metadata and basename fiction, preserves provider neutrality, reuses existing proof surfaces, and passed independent staff review; implementation remains pending until 6y lands on xb-rebased main.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge.
              application:
                action: advance
                target-stage: implementation
                state: pending
                blockers: []
---

Make gate-room preparation one mechanical operation. The First Officer supplies the
decision question and selected files; Spacedock derives room placement, portable ids,
locators, revisions, canonical digests, request authority, and the open gate attempt.
Chat remains the complete no-override path. The command emits one frozen room suitable
for a future provider handoff, but this task adds no provider discovery, probe,
executable, invocation, or retained-provider simulation.

## Problem

The xb recorder verifies a prepared room, but no command prepares one. A First Officer
still handcrafts `request.json`, Briefing ids and digests, artifact revisions, and a
room layout. The recorder then contradicts Subspace's current provider contract by
requiring the canonical Briefing leaf to be `briefing.json`.

That contradiction is observable, not hypothetical. It caused this task's first
relative-path record to fail before mutation with `Rel: can't make
docs/dev/.spacedock-state/... relative to ...`; the identical absolute path succeeded.
That path-base defect is already owned by 6y's CLI normalization seam. This task extends
that one launch-directory policy to its new inputs and `--room`; it does not add a
second resolver inside `internal/gates`.

The unreleased-v1 ruling is direct alignment. A prepared request has one required
Briefing locator, id, and digest. No old request shape, basename fallback, provider
version match, caller-selected output vector, compatibility wrapper, or
`association.json` is introduced.

## Cross-repository contract observed

At filing, Subspace `main` was `5ce887c` and the last reported active `em` commit was
`df75026`. Read-only inspection on 2026-07-24 found `main` still at `5ce887c`, while
`spacedock-ensign/remove-invented-briefing-review-preconditions` had advanced to
`27b32eb` through:

- `f3466bd`: package mode accepts one readable Briefing with any filename; Subspace
  allocates and owns the retained provider package and its Result, log, inventory, and
  diagnostics. Caller-supplied output paths are gone.
- `27b32eb`: package allocation precedes presenter discovery and capability preflight,
  so their failures retain and report the provider package plus diagnostics.
- Package eligibility is the literal
  `subspace-tui --supports review-v1-provider-package-v1` capability, not an exact
  version. The one-file Markdown profile keeps its unrelated version rule.

Implementation must re-check the landed `em` tip. These deltas strengthen, rather than
change, the boundary here: Subspace owns provider outputs, probing, retained preflight
failures, and eventual room-to-provider invocation; Spacedock owns room preparation and
recording. At `27b32eb`, room/request composition is still future q0 work. This task
therefore documents only a provider-neutral one-room handoff and neither implements nor
simulates a selected-provider transport.

## Riskiest-first spike

The spike built the current `cmd/spacedock` and used only
`internal/gates/testdata/gate-room/*` plus a minimal temporary workflow. Its room was
otherwise valid and named the exact fixture Briefing bytes `decision-material.data`.

1. `gate record ... --briefing decision-material.data` exited 1 before mutation:
   `--briefing must name a canonical package manifest named briefing.json`.
2. Copying those identical bytes to `briefing.json` made the same bind exit 0 and open
   `gate-attempt:3k-validation-1`.
3. Removing only `briefing.json` while retaining `decision-material.data` made
   `gate record ... --room ROOM` exit 1 at
   `read bound canonical Briefing: .../briefing.json: no such file or directory`.
4. Restoring the identical `briefing.json` bytes made the same room close successfully
   with `resolution:captain-validation-1`.

The minimum post-em seam is therefore exact: require `request.briefing.locator`,
resolve that frozen local locator once, and use the resulting path in binding, room
recording, validation, and eligibility wherever code currently checks a basename or
appends `briefing.json`. The implementation's first red test is this same differential
fixture with no `briefing.json` present.

## Proposed command and room contract

The new command is:

```text
spacedock gate prepare ENTITY --question TEXT --artifact GATE-REVIEW.md \
  [--reference FILE ...] [--workflow-dir DIR]
```

`--artifact` is required exactly once and is the concise gate review. Its filename must
end in `.md` or `.markdown` case-insensitively, and the generated Artifact always carries
`"mediaType": "text/markdown"`. `--reference` selects zero or more existing supporting
files in caller order; selecting the same normalized absolute path twice is an error.
The First Officer owns the question and file choices; it never supplies an id, digest,
room, attempt, locator, provider path, actor, or approver.

The CLI normalizes relative `--artifact`, `--reference`, `--briefing`, and `--room`
values against the invocation directory through 6y's retained-input normalization
helper before entering `internal/gates`. Each selected file must be a readable,
non-symlink regular file. For each generated Artifact or Reference, Spacedock records
the slash-relative path from the generated Briefing's directory to the selected source
and its full raw SHA-256 revision; these source URIs may contain `..` and are distinct
from the request locator's contained-path rules. References receive a deterministic
media type from a closed, case-insensitive extension table: `.md`/`.markdown` =
`text/markdown`, `.json` = `application/json`, `.yaml`/`.yml` =
`application/yaml`, `.txt`/`.log` = `text/plain`, and every other extension =
`application/octet-stream`. No host MIME database participates.

Under the entity lock, preparation derives the current-stage gate and next open attempt,
then publishes this fresh room as one operation:

```text
<entity-dir>/review/<stage>/briefing-<attempt-number>/
  gate-briefing.json
  request.json
```

`gate-briefing.json` is the generator's chosen filename, not a canonicality condition.
Artifact and Reference ids are derived from type, a normalized basename slug, and the
shortest unique raw-digest prefix. Prefixes start at 12 lowercase hex characters and
extend all colliding ids together by four characters through the full 64. If distinct
source URIs have identical full bytes and basename slug, append stable `-1`, `-2`, ...
suffixes in slash-relative URI lexical order while preserving caller order in the
Briefing; an exact repeated source path is rejected. Authoritative revisions always
retain the full digest. The Briefing id, gate id, attempt id, JCS Briefing digest,
request JCS digest, room reference, and Captain actor/approver are binary-owned.

`request.json` has the closed v1 shape:

```json
{
  "type": "spacedock-gate-presentation-request",
  "version": "1",
  "gate": "gate:task:validation",
  "attempt": "gate-attempt:task-validation-1",
  "briefing": {
    "locator": "gate-briefing.json",
    "id": "briefing:task:validation:attempt-1:revision-1",
    "digest": "sha256:..."
  },
  "actor": "person:captain",
  "approver": "person:captain"
}
```

On success, stdout is exactly four newline-terminated `key=value` lines:

```text
room=/clean/absolute/entity/path/review/<stage>/briefing-<n>
briefing=briefing:<derived-id>
digest=sha256:<64 lowercase hex>
state=open
```

`room` is the cleaned absolute path of the published, entity-derived room; no symlink
lookup or directory scan is required. `briefing` and `digest` are the exact frozen
binding read back from the successful operation. When a future handoff or diagnostic
needs the room, this emitted value is authoritative; callers must not rediscover it from
ids, attempt numbers, status output, or directory contents. The no-override chat path
has no following command that consumes `room=`.

The locator is a clean, non-empty, slash-relative path contained by the room. Absolute
paths, empty/dot paths, `..`, backslashes, and any symlink escape fail before mutation.
The resolved target must be a readable regular file. Nested clean locators are valid;
the generated room uses one leaf. A request-backed bind stores the room reference and
request digest as today. A Briefing-only bind stores the exact file reference, so
arbitrary filenames work without inventing an adjacent-request or basename fallback.
One shared resolver handles both cases: a request-backed binding validates
`request.json` and resolves its locator within the bound room; a request-less binding
resolves the exact stored file reference. Binding, room recording, validation, and
application eligibility all call that resolver rather than appending a filename.

Preparation stages a complete candidate room in a sibling temporary directory, validates
the same bytes through the production request/Briefing loader, atomically publishes the
fresh room, and binds the attempt under the same entity lock. Exact replay is a no-op.
An occupied divergent room, stale entity comparison, or handled bind/write failure
removes only the new candidate and leaves the entity and retained rooms byte-identical.
That is error atomicity while the entity lock is held, not cross-file crash atomicity:
a process or power loss after room rename but before entity replacement can leave an
unbound room. This task adds no journal or recovery schema; a later prepare refuses
divergent occupancy, and crash recovery remains an explicit operator concern.

## Recorder and JSON authority

One recursive token-stream JSON reader rejects duplicate object member names before
canonicalization or typed decoding. The request, located Briefing, exact Result, and
presented inventory all pass through it. Detection applies at every object depth, not
only to currently known authority fields; therefore a later field cannot silently
inherit Go's last-member-wins behavior.

After duplicate rejection, existing closed typed validation remains authoritative.
Unknown Result authority fields still fail, Result `Resolution.by` still must equal the
frozen approver, and inventory still must cover every Artifact and recursively reached
Reference exactly once.

The recorder derives one in-memory `spacedock-result-association` from the validated
request (including locator, gate, attempt, and authority), located Briefing, exact raw
Result, and exact inventory. It writes no association file. Durable state already has
the four required pins: request and Briefing digests on the binding, plus raw Result and
inventory digests in provider evidence. `gate validate` recomputes all four for a
provider-backed attempt; recording and validation fail if any input is missing,
substituted, duplicated, or changed.

## First Officer ownership and provider boundary

This design rebases on the latest 6y tip inspected at `60adfc1f` (the lifecycle work is
in `e9415a17`) rather than current `main`'s transitional presentation-channel prose.
After 6y lands, `skills/fo-gate-lifecycle/SKILL.md` owns preparation, gate capability
preflight, mutation, presentation routing, recording, and consumption.
`skills/present-gate/SKILL.md` is rendering-only and is not changed by this task.
Composed with current main's xb recorder, 6y's pre-xb `--result --association` wording
and capability anchor do not survive: provider-backed recording uses the retained room
surface and derives the association in memory. The lifecycle's Spacedock-only
`gate --help` preflight requires `prepare` and `--room`, not `--association`.

For the no-override path, `fo-gate-lifecycle` runs `gate prepare`, commits the entity
folder containing the generated room and binding, passes the gate-review Artifact to
the rendering-only `present-gate`, and records the captain's semantic chat decision.
Folder-scoped state commit includes the generated room without taking its stdout path
as an argument.

The lifecycle text may state the provider-neutral boundary: a landed presentation
override receives the one prepared room, never caller-built provider argv or output
paths; after that external transport returns, the same retained room is the
provider-backed recorder input. It must also state that this repository does not
discover, version-check, capability-probe, or launch a provider. There is no
selected-override execution arm in this task. Subspace q0 owns room-to-provider
invocation and its cross-repository proof; Subspace `27b32eb` already owns and proves
package allocation plus retained discovery and capability-preflight failures.
Spacedock tests do not fake those future events.

## New mechanisms and rejected alternatives

| Mechanism | Value AC | Simplest alternative | Why insufficient |
|---|---:|---|---|
| One `gate prepare` operation | AC-1 | Tell the FO to write two JSON files and call `gate record` | Preserves the manual ids/digests and partial-room failure that caused the task. |
| Frozen local Briefing locator | AC-2, AC-5 | Keep joining `briefing.json` | Fails the reproduced valid room and contradicts the provider contract. |
| One recursive duplicate-member reader | AC-3, AC-5 | Rely on `encoding/json` plus typed structs | Go accepts conflicting duplicates last-wins; the detached counterexample can close under the wrong authority. |
| Stable room/identity stdout handoff | AC-1 | Omit the room or make callers reconstruct it from ids/directory layout | Hides the published artifact and can select the wrong attempt under retries. |
| In-memory derived association | AC-5 | Persist `association.json` | Creates a second durable truth that can diverge from the four frozen inputs. |

## Expected surface and tolerance

Baseline assumption: latest 6y (`60adfc1f`, including lifecycle owner `e9415a17`) lands
first. Relative retained-input normalization is then available in `internal/cli`, the
existing recorded-gate journey targets `fo-gate-lifecycle`, and `present-gate` contains
rendering only. Against that composition, the smallest expected implementation is these
16 files and about `+986/-161` lines (**1,147 changed LOC**):

The inspected 6y tip is still pre-xb-rebase, so implementation must not start until
6y's final xb rebase lands. Re-read that landed tip before creating the worktree; if it
changes lifecycle ownership, recorder commands, shared live assertions, or any declared
file/delta below, return to ideation for a surface reset rather than implementing
against this provisional composition.

| File | Expected delta | Purpose |
|---|---:|---|
| `internal/cli/cli.go` | `+60/-6` | Route `gate prepare`, normalize inputs, and print the stable four-line result. |
| `internal/cli/gate_test.go` | `+175/-25` | Reuse CLI fixtures for preparation, stdout, arbitrary-locator eligibility, and byte-clean failures. |
| `internal/gates/prepare.go` (new) | `+220/-0` | Derivation, ids/media types, and error-atomic room publication. |
| `internal/gates/prepare_test.go` (new) | `+190/-0` | Focused replay, collision, URI, selection, and handled-error tests. |
| `internal/gates/operation.go` | `+80/-35` | Closed request locator and the one exact Briefing resolver. |
| `internal/gates/application.go` | `+12/-4` | Route reviewed-input eligibility through that resolver instead of `briefing.json`. |
| `internal/gates/io.go` | `+30/-8` | Recompute the four retained provider inputs through duplicate-safe reads. |
| `internal/gates/json.go` (new) | `+75/-0` | Recursive duplicate-member rejection. |
| `internal/gates/testdata/gate-room/request.json` | `+1/-0` | Add the locator to the canonical fixture. |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | `+24/-16` | Update the shared no-override observation used by existing Claude/Codex/Pi lanes; add no lane. |
| `internal/contractlint/fo_function_reference_invariant_test.go` | `+18/-8` | Pin lifecycle ownership, forbidden provider mechanics, and rendering-only `present-gate`. |
| `docs/specs/gate-resolution-frontmatter-contract.md` | `+60/-30` | Normative prepare/request/resolver/atomicity contract. |
| `docs/site/reference/command-reference.md` | `+14/-6` | New verb, stdout, and arbitrary-name recording. |
| `docs/site/reference/frontmatter-contract.md` | `+3/-3` | Remove the manifest-basename claim. |
| `docs/site/concepts/gates-and-decisions.md` | `+6/-6` | Mechanical no-override preparation and future handoff boundary. |
| `skills/fo-gate-lifecycle/SKILL.md` | `+18/-14` | Replace hand bind/association wording with prepare and provider-neutral room recording. |

Tolerance is **+2 files and +25% changed LOC** (hard cap 18 files / 1,434 changed
LOC), for a focused resolver test or fixture split only. A change to
`skills/present-gate/SKILL.md`, any schema field in `gates`, new dependency, provider
executable/probe/transport, selected-provider harness, compatibility request shape,
caller output path, association artifact, or broader lifecycle routing requires a
design reset.

## Acceptance criteria

**AC-1 (VALUE) — One command turns judgment and file choices into a validated open
gate room with zero caller-authored metadata files.** Starting baseline: the fixture has
one gate-review Markdown file, selected References, and no room. After `gate prepare`,
the derived room contains exactly the request and located canonical Briefing, the attempt
is open, and `gate validate` succeeds; source artifacts remain in place. The four stdout
lines expose the exact cleaned absolute room, Briefing id, digest, and open state; the
required Artifact is `text/markdown`, Reference media types follow the closed table, and
every source URI is slash-relative to the generated Briefing directory. *Test:* real CLI
fixture asserts pre-command metadata count 0, post-command count 2, exact stdout,
paths/media types/full digests/relative URIs, and no copied source. It fails if the
fixture must supply metadata, if stdout omits or misstates the published room, or if
output changes under a launch directory containing spaces.

**AC-2 — Every request-backed operation uses the frozen readable Briefing locator,
independent of basename.** A clean nested `decision-material.data` locator binds,
records, validates, and leaves an approved provider-backed application eligible;
traversal, absolute/backslash locator, symlink escape, substitution, id mismatch, and
digest mismatch all exit nonzero with byte-identical entity state. *Test:* the
command-level table seeded by the spike uses the existing gate-room fixture, closes
through `gate record --room`, then calls CLI `gate eligibility` and requires
`approved-pending eligible=true`. Deleting any one call to the exact resolver in bind,
room record, validation, or `internal/gates/application.go` makes the arbitrary-locator
positive case fail.

**AC-3 — Conflicting duplicate members in every authority-bearing room document fail
before mutation.** Request, located Briefing (including nested context), Result
(including nested Resolution), and presented inventory cases each inject one
last-wins counterexample. *Test:* detached adversarial table requires nonzero exit,
diagnostic naming the duplicate member, unchanged whole entity bytes, and no lock
residue. Removing the recursive reader makes at least one case bind or close.

**AC-4 — Spacedock owns a truthful one-room lifecycle boundary without implementing a
provider transport.** With no override, 6y's `fo-gate-lifecycle` runs prepare once,
commits the entity folder containing the generated room and binding, renders through the
unchanged rendering-only `present-gate`, and records chat. Its future-override wording
promises only that a landed override receives that one room. No Spacedock Go or skill
change names a provider executable, runs a provider availability, version, or capability
probe, allocates provider outputs, or simulates invocation. *Test:* the shared
recorded-gate observation requires `gate-help → prepare → state-commit → chat-render →
decision-record`, followed by 6y's unchanged
decision-commit/consume/consumed-commit barriers. At final tip, the existing Claude,
Codex, and Pi live lanes must each satisfy that same observation. Legitimate structural
checks require `fo-gate-lifecycle` ownership, keep `present-gate` rendering-only, and
reject `subspace-tui`, `/subspace:r`, `--supports`, `--version`, or a new process-launch
import in the changed gate/lifecycle surface. No prose-derived room-consumption mutant
or fake override lane is added; q0 owns room-to-provider and retained-preflight proof.

**AC-5 — Provider recording has one recomputed association and no parallel durable
artifact.** The full fixture prepares, receives fixed Result/inventory outputs, closes,
and validates with no `association.json`. Request, located Briefing, Result, and
inventory are each deleted and byte-mutated independently; every variant fails
recording or read-only validation without changing the entity. *Test:* real CLI
end-to-end fixture asserts the four frozen digest pins and exact room tree; adding an
association input/file or omitting one frozen input fails.

## Test plan and proof order

0. **Baseline gate, before implementation:** require 6y's final xb rebase to be landed,
   record its tip, and compare lifecycle ownership, recorder commands, shared assertions,
   and the expected-surface table. Any mismatch returns to ideation for reset.
1. **Focused red/green, low cost:** add the arbitrary-name spike as the first command
   test using the existing gate-room fixture, then add focused `prepare_test.go` cases
   for exact stdout data, relative URIs/media types, 12-to-64-character digest-prefix
   extension/full-digest suffixes, replay, occupied room, and handled-error cleanup.
   Run `go test ./internal/gates ./internal/cli -count=1`.
2. **Adversarial JSON, medium cost:** mutate each of the four room documents at top
   level and nested authority-bearing objects. Assert entity bytes and lock state, not
   only error substrings. The arbitrary-locator positive case continues through provider
   room closure and CLI eligibility, so `application.go` cannot silently retain its
   basename join.
3. **Existing FO journey only, high cost at final tip:** update 6y's shared
   recorded-gate observation in place to assert
   `gate-help → prepare → state-commit → chat-render → decision-record` before its
   unchanged decision-commit/consume/consumed-commit barriers. Run the existing
   `TestLiveClaudeSharedScenarios` and `TestLiveCodexSharedScenarios`
   `recorded-gate-lifecycle` cases plus `TestLivePiRecordedGateLifecycle` against the
   final implementation tip; all three must observe the revised sequence. Add no host
   lane, harness, provider fake, provider-capability ledger, prose-derived room mutant,
   or cross-repo invocation test. q0 owns those proofs after its room command exists.
4. **Repository gates:** `gofmt -w ./cmd ./internal`, `go test ./...`,
   `go test ./... -race`, strict docs build, `git diff --check`, and verify
   `go list -deps ./cmd/spacedock` contains no Subspace package, the changed gate code
   imports no process launcher, and lifecycle command text contains none of
   `subspace-tui`, `/subspace:r`, `--supports`, or `--version`.
5. **High-stakes detached audit:** independently inject conflicting duplicate
   `by`, locator, digest, id, and inventory members and try to refute the byte-clean
   claim. Re-check landed Subspace `em` only to confirm the ownership boundary; do not
   run its future room transport as Spacedock evidence.

## Documentation change proposal

The implementation applies these concrete semantics (line wrapping may follow the
target file):

```diff
--- docs/site/reference/command-reference.md
+++ docs/site/reference/command-reference.md
@@
-| `spacedock gate record <entity> --briefing PATH/briefing.json` | Bind a complete retained package manifest whose basename is exactly `briefing.json`. Other basenames fail before mutation. |
+| `spacedock gate prepare <entity> --question TEXT --artifact REVIEW.md [--reference FILE ...]` | Derive and bind one recorder-ready room. Success prints exactly `room`, `briefing`, `digest`, and `state` key/value lines; `room` is the clean absolute path the first officer uses directly. The required Artifact is Markdown; generated source URIs are relative to the generated Briefing directory. |
+| `spacedock gate record <entity> --briefing PATH` | Bind any readable canonical Briefing by its exact path. A prepared room instead freezes its Briefing locator, id, and digest in `request.json`; every later operation resolves that locator rather than a canonical basename. |

--- docs/site/concepts/gates-and-decisions.md
+++ docs/site/concepts/gates-and-decisions.md
@@
-Before the First Officer shows a gate, it binds the exact retained Briefing and commits that package.
+Before the First Officer shows a no-override gate, it prepares and binds the room mechanically, commits the entity folder containing that room, then renders in chat. A future presentation override receives the command's authoritative `room=` value through its own landed transport; Spacedock does not discover, probe, or launch a provider.

--- skills/fo-gate-lifecycle/SKILL.md
+++ skills/fo-gate-lifecycle/SKILL.md
@@
-**Retain and bind.** Assemble `ROOM/briefing.json` ... then run `gate record ENTITY --briefing BRIEFING`.
+**Prepare and bind.** Select one Markdown gate-review Artifact and any References, then run `${SPACEDOCK_BIN:-spacedock} gate prepare ENTITY --question QUESTION --artifact REVIEW [--reference FILE ...] --workflow-dir WORKFLOW_DIR`. Require the four stable output lines and `state=open`; commit the entity folder containing the generated room and binding before presentation. The emitted `room=` value is the sole future override/diagnostic locator and must never be reconstructed or searched for.
+**Presentation boundary.** With no override, render the generated review through `present-gate` and record chat. A future landed override receives the one prepared room and returns that retained room for `gate record --room`; this lifecycle does not name, discover, version-check, capability-probe, or launch a provider, and does not construct provider output paths or an association.
```

The normative spec makes the same substitutions, defines the closed request shape and
room publication/error-atomicity behavior, and states explicitly that association is
recomputed and unstored. The frontmatter reference removes only its exact-basename
claim; no `gates` schema change is proposed. `skills/present-gate/SKILL.md` remains
unchanged and rendering-only after 6y.

## Deferred mutable-source movement

Preparation pins selected source bytes by full raw digest but does not copy them; their
URIs remain relative references from the generated Briefing directory. Snapshotting or
freezing mutable sources is deliberately deferred. Promote that policy into this task
and reset the design before implementation only if a fixture or landed q0 contract
requires a prepared provider to reopen those URIs after a normal state commit and shows
that an allowed move, deletion, or byte change can occur before presentation. Without
that evidence, source stability through the decision is the lifecycle precondition and
the room stays the two generated metadata files.

## Out of scope

- Subspace q0, room-to-provider invocation, terminal transport, provider discovery or
  capability probing, provider output allocation, retained-preflight proof, or provider
  retention implementation.
- Compatibility request parsing, a `briefing.json` fallback for prepared requests, or
  migration wrappers.
- `association.json`, caller-selected Result/log/inventory/diagnostic paths, or provider
  argv.
- Broader lifecycle-next-action prose, advisory-round preparation, readiness projection,
  crash-atomic multi-file transactions, artifact copying, or generic JSON framework work.

### Feedback Cycles

- **Cycle 5 — cross-checkout source-locator correction (2026-07-25).** Do not consume
  the current ideation approval. A Briefing-relative URI is sound for room-owned files,
  but a `..` chain that escapes the state checkout into the main checkout depends on
  one local nesting layout. A raw SHA verifies bytes only after lookup; it cannot locate
  or reconstruct them. Revise the design so Git-owned inputs identify the checkout
  history, exact commit, repository-relative path, and raw-byte digest, or are frozen in
  the retained room. Prove the choice by reopening the package when main and state are
  independently located. Preserve Review v1's Briefing-relative rule for room-owned
  files and avoid adding provider-specific transport.

## Stage Report: ideation

- DONE: Reproduce xb's arbitrary-Briefing-basename failure with the smallest valid prepared room, then identify the minimum post-em Spacedock seam that makes it pass.
  Current fixture bytes fail only at the pre-lock basename guard and the later `filepath.Join(..., "briefing.json")`; identical bytes under that leaf bind and close, so the minimum seam is the frozen request locator used by bind/record/validate/eligibility.
- DONE: Produce one coherent unreleased-v1 design with exact expected files/LOC and no compatibility layer, association.json, Subspace transport, or broader ergonomics machinery.
  The design defines one `gate prepare` command, a closed locator-bearing request, 15 files / 1,290 changed LOC with 17 / 1,613 tolerance, and explicit reset triggers for every excluded mechanism.
- DONE: Define falsifiable command/fixture/live evidence for mechanical preparation, locator binding, all authority-bearing duplicate-member rejection, derived association, and capability-based presentation.
  AC-1–AC-5 each name a behavior-changing counterexample, byte/on-disk assertions, proof cost, and the 6y shared journey to augment rather than a second lifecycle harness.

### Summary

The design removes both manual room metadata and the recorder's canonical-basename
fiction while preserving one provider-neutral, request-frozen authority boundary.
Subspace's active `em` delta is recorded at `27b32eb`; implementation must re-check its
landed form and reuse 6y's CLI path normalization before final validation.

## Stage Report: ideation (cycle 2)

- DONE: Repair the provider-selection contract to match active Subspace em at `27b32eb`.
  The design now distinguishes no-override chat from a selected override's single invocation; Subspace owns package allocation, discovery, and the literal capability probe, while preflight failure retains evidence, leaves the gate open, and invokes no chat or recorder.
- DONE: Make AC-4 and its integration proof falsify pre-probing, version matching, retry, fallback, and lost preflight evidence.
  Four exact event ledgers cover chat, capable override, non-capable override, and stale capability claims; forbidden Spacedock probes, second invocations, host launch after nonzero capability, chat fallback, recorder calls, or gate closure fail the sequence.
- DONE: Align the proposed user and skill wording without expanding Spacedock into provider transport.
  The concrete documentation diff now says one room/one provider invocation, provider-owned retained preflight, literal capability rather than version, and no chat after a selected override fails; q0 remains explicitly future Subspace work.

### Summary

Cycle 2 removes the material cross-repository mismatch: chat is a default channel, not
a fallback from a selected provider. A selected provider is invoked once with the
prepared room and owns retained compatibility evidence from its first preflight byte.

## Stage Report: ideation (cycle 3)

- DONE: Remove future-q0 transport simulation from AC-4 and the Spacedock harness.
  AC-4 now proves only no-override chat, one-room handoff wording, and absence of provider executables/probes; q0 owns invocation and Subspace retains its existing preflight proof.
- DONE: Rebase the design surface on 6y's actual lifecycle ownership.
  Inspection of latest 6y `60adfc1f`/`e9415a17` retargets behavior and checks to `fo-gate-lifecycle`, leaves `present-gate` rendering-only, and recalculates 16 files / 1,147 changed LOC.
- DONE: Add eligibility to the exact locator-resolution seam.
  `internal/gates/application.go` is declared, and the existing gate-room/CLI fixture must close and remain `approved-pending eligible=true` with only an arbitrary nested locator.
- DONE: Complete the ergonomic prepare command contract.
  Stable stdout returns absolute room/id/digest/open state, URIs are Briefing-directory-relative, Markdown/media types and collision rules are deterministic, and the FO is forbidden to rediscover the room.
- DONE: Bound proof and atomicity to landed behavior.
  The plan reuses 6y's chat journey without a provider lane, defers mutable-source snapshots behind an evidence trigger, and distinguishes handled-error cleanup under lock from unclaimed crash atomicity.

### Summary

Cycle 3 supersedes cycle 2's selected-provider event ledger. The accepted mechanical
prepare, arbitrary locator, recursive duplicate rejection, four pins, and unstored
association remain, now within the smallest post-6y Spacedock-owned boundary.

## Stage Report: ideation (cycle 4)

- DONE: Remove unobservable direct `room=` consumption from AC-4.
  Stable room/id/digest stdout now serves AC-1 only; AC-4 proves the observable no-override sequence and legitimate provider-absence/ownership checks without a prose-derived room mutant.
- DONE: Require all existing host-neutral recorded-gate live lanes at final tip.
  The plan names Claude and Codex shared `recorded-gate-lifecycle` cases plus Pi's existing live test, all using one revised shared observation and unchanged commit/consume barriers.
- DONE: Preserve the 6y final-xb-rebase implementation dependency.
  The design blocks worktree creation until that rebase lands and requires an ideation reset if its final ownership, commands, assertions, or expected surface differ.

### Summary

Cycle 4 makes proof match observable behavior: no no-override command consumes the
reported room, while every existing live host must prove the revised lifecycle. The
16-file / 1,147-LOC surface remains provisional until 6y's final xb rebase lands.
