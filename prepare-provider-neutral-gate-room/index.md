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
        gate: gate:docs-dev:s4:backlog
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
---

Make gate-room preparation one mechanical operation. The First Officer supplies the
decision question and selected files; Spacedock derives room placement, portable ids,
locators, revisions, canonical digests, request authority, and the open gate attempt.
Chat remains the default when no override is selected. A selected provider consumes
the same frozen package exactly once and owns every preflight and retained output.

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
change, the boundary here: Subspace owns provider outputs and probing; Spacedock owns
room preparation and recording. Subspace q0 later owns `/subspace:r gate <gate-room>`;
this task neither implements nor simulates that transport.

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

`--artifact` is required exactly once and is the concise gate review. `--reference`
selects zero or more existing supporting files in caller order. The First Officer owns
the question and file choices; it never supplies an id, digest, room, attempt, locator,
provider path, actor, or approver.

The CLI normalizes relative `--artifact`, `--reference`, `--briefing`, and `--room`
values against the invocation directory through 6y's retained-input normalization
helper before entering `internal/gates`. Each selected file must be a readable,
non-symlink regular file. Spacedock records a slash-relative URI and raw SHA-256
revision without copying reproducible source bytes.

Under the entity lock, preparation derives the current-stage gate and next open attempt,
then publishes this fresh room as one operation:

```text
<entity-dir>/review/<stage>/briefing-<attempt-number>/
  gate-briefing.json
  request.json
```

`gate-briefing.json` is the generator's chosen filename, not a canonicality condition.
Artifact and Reference ids are derived from type, normalized basename, and a digest
prefix, so same-basename files remain distinct without caller ids. The Briefing id,
gate id, attempt id, JCS Briefing digest, request JCS digest, room reference, and
Captain actor/approver are binary-owned.

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

The locator is a clean, non-empty, slash-relative path contained by the room. Absolute
paths, empty/dot paths, `..`, backslashes, and any symlink escape fail before mutation.
The resolved target must be a readable regular file. Nested clean locators are valid;
the generated room uses one leaf. A request-backed bind stores the room reference and
request digest as today. A Briefing-only bind without `request.json` stores the exact
file reference, so arbitrary filenames work without inventing a request fallback.

Preparation stages a complete candidate room in a sibling temporary directory, validates
the same bytes through the production request/Briefing loader, atomically publishes the
fresh room, and binds the attempt under the same entity lock. Exact replay is a no-op.
An occupied divergent room, stale entity comparison, or bind failure removes only the
new temporary room and leaves entity and retained rooms byte-identical. The existing
immutable-room publication pattern is the model; there is no daemon or recovery schema.

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

## Presentation selection

Chat is the default only when no override is selected. That arm prepares/binds the
canonical room, renders the gate review in the current conversation, and records the
semantic chat decision. It does not invoke or probe a provider.

Once a declared override is selected, the First Officer prepares the room and invokes
that provider exactly once with the room. It does not independently discover the
presenter, run a version check, run the literal capability probe, construct provider
argv or output paths, retry the invocation, or select chat after any provider failure.

The invoked provider owns presenter discovery and the literal
`review-v1-provider-package-v1` probe. Package-mode eligibility depends on that
capability's zero exit, never release text. The provider allocates its retained package
before discovery/probing, records capability stdout, stderr, and exit there, and reports
the package on a nonzero preflight. Such a failure opens no host surface but is still
the one completed provider invocation: the gate remains open with diagnostic evidence,
and chat is not invoked. A false positive capability claim fails later at its concrete
parser, loader, publication, or validation boundary with the same no-retry/no-chat rule.

The skill must not claim the future Subspace q0 transport exists before it lands. It
states only the provider-neutral one-room/one-invocation contract that q0 will implement.

## New mechanisms and rejected alternatives

| Mechanism | Value AC | Simplest alternative | Why insufficient |
|---|---:|---|---|
| One `gate prepare` operation | AC-1 | Tell the FO to write two JSON files and call `gate record` | Preserves the manual ids/digests and partial-room failure that caused the task. |
| Frozen local Briefing locator | AC-2, AC-5 | Keep joining `briefing.json` | Fails the reproduced valid room and contradicts the provider contract. |
| One recursive duplicate-member reader | AC-3, AC-5 | Rely on `encoding/json` plus typed structs | Go accepts conflicting duplicates last-wins; the detached counterexample can close under the wrong authority. |
| Provider-owned literal capability preflight | AC-4 | Pre-probe or match an exact Subspace version in Spacedock | Duplicates provider mechanics, loses retained preflight evidence, and couples Spacedock to another release train. |
| In-memory derived association | AC-5 | Persist `association.json` | Creates a second durable truth that can diverge from the four frozen inputs. |

## Expected surface and tolerance

Baseline assumption: 6y lands first and this branch extends its CLI path-normalization
helper. The expected implementation is exactly these 15 files and about
`+1,125/-165` lines (**1,290 changed LOC**):

| File | Expected delta | Purpose |
|---|---:|---|
| `internal/cli/cli.go` | `+40/-5` | Route `gate prepare`; extend 6y's path normalization. |
| `internal/cli/gate_test.go` | `+220/-35` | Command preparation, arbitrary locator, byte-clean adversarial matrix. |
| `internal/gates/prepare.go` (new) | `+230/-0` | Derivation and atomic room publication. |
| `internal/gates/prepare_test.go` (new) | `+230/-0` | Pure preparation/replay/collision/file-selection behavior. |
| `internal/gates/operation.go` | `+90/-45` | Request locator and shared exact Briefing resolution. |
| `internal/gates/io.go` | `+35/-8` | Recompute all retained provider inputs. |
| `internal/gates/json.go` (new) | `+80/-0` | Recursive duplicate-member rejection. |
| `internal/gates/testdata/gate-room/request.json` | `+1/-0` | Add the locator to the canonical fixture. |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | `+70/-10` | Augment 6y's provider-neutral shared journey. |
| `internal/contractlint/fo_function_reference_invariant_test.go` | `+15/-10` | Skill order and truthful capability anchors. |
| `docs/specs/gate-resolution-frontmatter-contract.md` | `+60/-35` | Normative prepare/request/recorder contract. |
| `docs/site/reference/command-reference.md` | `+10/-6` | New verb and arbitrary-name recording. |
| `docs/site/reference/frontmatter-contract.md` | `+4/-4` | Remove manifest-basename claim. |
| `docs/site/concepts/gates-and-decisions.md` | `+14/-9` | Chat default and capability behavior. |
| `skills/present-gate/SKILL.md` | `+28/-28` | Mechanical prepare plus provider-neutral selection. |

Tolerance is **+2 files and +25% changed LOC** (hard cap 17 files / 1,613 changed
LOC), for a focused helper or fixture split only. Any schema field in `gates`, new
dependency, provider executable/transport, compatibility request shape, caller output
path, association artifact, or broader FO lifecycle edit requires a design reset.

## Acceptance criteria

**AC-1 (VALUE) — One command turns judgment and file choices into a validated open
gate room with zero caller-authored metadata files.** Starting baseline: the fixture has
one gate-review Markdown file, selected References, and no room. After `gate prepare`,
the derived room contains exactly the request and located canonical Briefing, the attempt
is open, and `gate validate` succeeds; source artifacts remain in place. *Test:* real CLI
fixture asserts pre-command metadata count 0, post-command count 2, exact paths/digests,
and no copied source. It fails if the fixture must supply an id/digest/room or prewrite
JSON.

**AC-2 — Every request-backed operation uses the frozen readable Briefing locator,
independent of basename.** A clean nested `decision-material.data` locator binds,
records, validates, and remains eligible; traversal, absolute/backslash locator,
symlink escape, substitution, id mismatch, and digest mismatch all exit nonzero with
byte-identical entity state. *Test:* command-level table seeded by the spike; deleting
the old `briefing.json` assumption must not affect the positive case.

**AC-3 — Conflicting duplicate members in every authority-bearing room document fail
before mutation.** Request, located Briefing (including nested context), Result
(including nested Resolution), and presented inventory cases each inject one
last-wins counterexample. *Test:* detached adversarial table requires nonzero exit,
diagnostic naming the duplicate member, unchanged whole entity bytes, and no lock
residue. Removing the recursive reader makes at least one case bind or close.

**AC-4 — Presentation preserves the selected channel and delegates compatibility to
one provider invocation.** With no override, only chat prepares, presents, and records.
With an override selected, the provider receives exactly one prepared room and internally
runs exactly one literal provider-package capability probe; varied `--version` text does
not change eligibility. A missing presenter or nonzero capability retains one reported
provider package with inputs and capability stdout/stderr/exit, opens no host surface,
leaves the gate open, and invokes neither chat nor the recorder. A false capability
claim retains its concrete downstream failure with no retry or chat. *Test:* augment
6y's shared recorded-gate journey and host bindings with the exact event ledgers below;
any Spacedock-side probe, version comparison, second provider invocation, chat fallback,
or closed gate fails the sequence.

**AC-5 — Provider recording has one recomputed association and no parallel durable
artifact.** The full fixture prepares, receives fixed Result/inventory outputs, closes,
and validates with no `association.json`. Request, located Briefing, Result, and
inventory are each deleted and byte-mutated independently; every variant fails
recording or read-only validation without changing the entity. *Test:* real CLI
end-to-end fixture asserts the four frozen digest pins and exact room tree; adding an
association input/file or omitting one frozen input fails.

## Test plan and proof order

1. **Focused red/green, low cost:** add the arbitrary-name spike as the first command
   test, then preparation/replay/collision unit tests. Run
   `go test ./internal/gates ./internal/cli -count=1`.
2. **Adversarial JSON, medium cost:** mutate each of the four room documents at top
   level and nested authority-bearing objects. Assert entity bytes and lock state, not
   only error substrings.
3. **Shared FO journey, medium/high cost:** extend 6y's landed recorded-gate fixture
   rather than building a second lifecycle harness. Assert these exact provider-neutral
   event ledgers:
   - no override: `prepare → chat-present → chat-record`, with no provider event;
   - capable override: `prepare → provider-invoke → provider-package →
     capability-probe(0) → host-launch → result-retained → room-record`, with no chat or
     version probe;
   - missing/non-capable override: `prepare → provider-invoke → provider-package →
     capability-probe(nonzero) → provider-failure`, with no host launch, chat, recorder,
     retry, or gate closure;
   - stale capability claim: one provider invocation and package, followed by its
     concrete retained downstream failure, with no retry or chat.
   Host adapters reuse the shared scenario; run only live lanes required by the final
   skill/runtime diff.
4. **Repository gates:** `gofmt -w ./cmd ./internal`, `go test ./...`,
   `go test ./... -race`, strict docs build, `git diff --check`, and verify
   `go list -deps ./cmd/spacedock` contains no Subspace package.
5. **High-stakes detached audit:** independently inject conflicting duplicate
   `by`, locator, digest, id, and inventory members and try to refute the byte-clean
   claim. Re-check landed Subspace `em` before final validation.

## Documentation change proposal

The implementation applies these concrete semantics (line wrapping may follow the
target file):

```diff
--- docs/site/reference/command-reference.md
+++ docs/site/reference/command-reference.md
@@
-| `spacedock gate record <entity> --briefing PATH/briefing.json` | Bind a complete retained package manifest whose basename is exactly `briefing.json`. Other basenames fail before mutation. |
+| `spacedock gate prepare <entity> --question TEXT --artifact FILE [--reference FILE ...]` | Derive and bind one recorder-ready room; callers choose judgment and files, not ids, digests, locators, authority, or output paths. |
+| `spacedock gate record <entity> --briefing PATH` | Bind any readable canonical Briefing. A prepared request must bind that exact file by local locator, id, and digest. |

--- docs/site/concepts/gates-and-decisions.md
+++ docs/site/concepts/gates-and-decisions.md
@@
-If the provider is missing or has the wrong version, the first officer names the remedy and returns to chat.
+Chat is the default only when no override is selected. Once selected, an override receives one prepared room and owns presenter discovery plus the literal provider-package capability probe; exact provider versions are not eligibility. Preflight failure retains the reported provider package and diagnostics, leaves the gate open, and never retries or invokes chat.

--- skills/present-gate/SKILL.md
+++ skills/present-gate/SKILL.md
@@
-1. **Probe before side effects.** Run the override's read-only availability and version probe...
-2. **Pass one prepared room.** The scaffold owns `request.json`, the canonical `briefing.json`...
+1. **Select once.** Chat is the default only when no override is selected. A selected override is invoked exactly once and never falls back to chat.
+2. **Prepare mechanically.** Run `${SPACEDOCK_BIN:-spacedock} gate prepare ...`; pass the resulting room, never provider argv or output paths.
+3. **Leave preflight with the provider.** The provider allocates retained state, discovers its presenter, and runs the literal provider-package capability probe; Spacedock does not pre-probe or compare a version.
+4. **Retain failed invocations.** A provider preflight or downstream failure reports its retained package, leaves the gate open, and triggers no retry, recorder call, or chat presentation.
+5. **Resolve the frozen Briefing.** The request's local locator, id, and digest identify the canonical Briefing; no filename is canonical.
```

The normative spec makes the same substitutions, defines the closed request shape and
room publication behavior, and states explicitly that association is recomputed and
unstored. The frontmatter reference removes only its exact-basename claim; no `gates`
schema change is proposed.

## Out of scope

- Subspace q0, `/subspace:r gate`, terminal transport, provider output allocation, or
  provider retention implementation.
- Compatibility request parsing, a `briefing.json` fallback for prepared requests, or
  migration wrappers.
- `association.json`, caller-selected Result/log/inventory/diagnostic paths, or provider
  argv.
- Broader lifecycle-next-action prose, advisory-round preparation, readiness projection,
  artifact copying, mutable-source snapshot policy, or generic JSON framework work.

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
