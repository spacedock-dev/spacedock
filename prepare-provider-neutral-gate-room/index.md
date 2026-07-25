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
                state: superseded
                blockers: []
            - id: gate-attempt:s4-ideation-2
              briefing:
                id: briefing:docs-dev:s4:ideation:attempt-2:revision-1
                digest: sha256:706374c6491bbf8b3a43a6469aa85eefbeb513d16f1dcfad1297b6eff97bb949
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:s4:ideation:2
                briefing: briefing:docs-dev:s4:ideation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-07-25T05:51:29.330095Z"
                decision: approve
                reason: Frozen room-owned sources remove checkout-topology coupling, the independent reopen spike falsifies the old locator, and staff review found no material issue; implementation remains dependency-held behind 6y.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
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
for recorder binding and validation. It is not presentation-ready for current Subspace
package mode; `git-root-review-v1-materialization`
(`rqh46ey33aqq4rt72b4w1m2q`) owns that required pre-v0.27 bridge. This task adds no
provider discovery, probe, executable, invocation, materialization, or retained-provider
simulation.

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

The original preparation design also confused a byte pin with a durable locator.
6y's retained gate-contract Artifact uses
`../../../../../../../docs/specs/gate-resolution-frontmatter-contract.md`; that URI
reaches the main checkout only because the state checkout happens to be nested under
`docs/dev`. Its `rev` verifies bytes after lookup but cannot identify the main Git
history or reconstruct a missing file. Main and state are separate histories, so a
cross-checkout `..` chain is not a portable source model.

The unreleased-v1 ruling is direct alignment. A prepared request has one required
Briefing locator, id, and digest. No old request shape, basename fallback, provider
version match, caller-selected output vector, compatibility wrapper, or
`association.json` is introduced.

Recorder-ready and presentation-ready are now explicit separate states. Current
Subspace package mode receives the canonical Briefing unchanged, opens Artifact URIs as
filesystem paths, and rejects Reference URIs containing `://`. It has no logical-root
map and never calls Spacedock's Git resolver. Therefore a `git-root` Briefing can be
bound, validated, and associated by Spacedock but cannot be claimed as renderable by
Subspace.

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
failures, and eventual room-to-provider invocation; Spacedock owns recorder-ready room
preparation and recording. At `27b32eb`, room/request composition is still future q0
work, but q0 transport alone cannot resolve Git-root content. This task therefore
documents only a provider-neutral recorder boundary and neither implements nor
simulates presentation. The filed sibling
`git-root-review-v1-materialization`/`rqh46ey33aqq4rt72b4w1m2q` is the explicit
cross-repository presentation owner and blocks the 0.27 pre-release until its actual
Subspace end-to-end proof passes.

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

### Cycle-5 independent-checkout spike

The correction spike used 6y's retained validation Briefing and its exact gate-contract
bytes. In the current nested checkout, the seven-level `..` URI resolved to raw SHA
`6425a361d95b0f5fc64d712aa9da8bd1b377ca9579ec6da5c6e3523fc8039003`.
The spike then made an ordinary state-repository commit containing the Briefing and a
room-local frozen copy, cloned that state history under an unrelated temporary root
with no main checkout at the implied topology, and reopened from the cloned Briefing
directory. The old URI was not found. The candidate
`sources/artifact-gate-contract-6425a361d95b` URI resolved, its raw SHA was still
`6425…9003`, and `cmp` found its bytes exact.

This is the smallest falsifier because it removes only the hidden nesting assumption.
It proved that a room-local copy can survive, but Cycle 6 rejects that result because it
duplicates an object already owned by one split-root Git history. The spike remains the
counterexample to `..`, not the selected implementation.

### Cycle-6 Git-object reopen spike

The replacement spike pinned the main gate contract at commit
`cc51e518a3420b01fd4b455e9710d38803dc6d3e`, path
`docs/specs/gate-resolution-frontmatter-contract.md`, raw SHA `6425…9003`, and the
state entity at commit `6a2ea132e2ef125a8996378c1306009b25c979fd`, path
`prepare-provider-neutral-gate-room/index.md`, raw SHA `f019…235f`. It cloned the main
and state histories separately, made later commits that moved both selected paths out
of their worktrees, and then moved the two checkouts again to unrelated directory
depths. With an explicit logical root map, `git show <commit>:<path>` in each moved
checkout recovered the pinned raw SHA exactly; neither current worktree path existed,
and the old seven-level `..` URI was not found.

The selected locator is therefore:

```text
git-root://main/cc51e518a3420b01fd4b455e9710d38803dc6d3e/docs/specs/gate-resolution-frontmatter-contract.md
git-root://state/6a2ea132e2ef125a8996378c1306009b25c979fd/prepare-provider-neutral-gate-room/index.md
```

The URI locates an immutable commit/path object in an already-present local Git object
database; the Review v1 `rev` independently verifies its raw bytes with SHA-256.
Checkout paths are runtime root-map values, not durable metadata. The implementation
fixture repeats the move/later-worktree-commit case and fails if resolution consults
current worktree bytes, needs the old nesting, fetches a remote, or returns bytes that
do not match `rev`.

That is a local address/reopen spike only. It did not invoke `gate prepare`, did not
exercise a production Spacedock resolver, and did not send a Git-root Briefing through
Subspace. It proves the commit/path identity can recover bytes from moved local object
databases; it is not presentation evidence. Current Subspace's filesystem-path/`://`
rules are the retained negative boundary, and sibling `rqh46ey33aqq4rt72b4w1m2q` owns
the first real resolved-byte presentation proof.

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
non-symlink regular file owned by the workflow's main or state Git history. Ownership
is determined by matching the selected checkout's Git common directory to the
definition or entity root; a linked implementation worktree therefore remains `main`.
Preparation derives the full `HEAD` commit and repository-relative path from that
selected checkout, reads `<commit>:<path>` from its local object database, and requires
those bytes to equal the selected worktree file. Untracked, dirty, missing-object, or
third-repository selections fail before mutation with an instruction to commit the
selected source. The full commit must be reachable from at least one ordinary local,
remote-tracking, or tag ref reported by `git for-each-ref --contains`; reflog-only and
detached unreachable commits are rejected. For newly authored state content, the First
Officer makes the path-scoped source commit before prepare, so the later room/binding
commit retains it as an ancestor. Main/worktree content likewise remains on a retained
branch through the gate lifecycle.

Each generated Artifact or Reference records a closed locator of the form
`git-root://<root>/<full-commit>/<repo-path>` plus its full raw SHA-256 `rev`. V1 root
names are `main` for the definition Git history and `state` for a distinct entity Git
history; an inline workflow emits only `main`. The commit is the owning repository's
full lowercase object id; that repository's local Git object database determines its
hash format, so no separate format segment is serialized absent a demonstrated
cross-root ambiguity. Each UTF-8 repository path segment uses canonical RFC 3986
escaping while `/` remains the segment separator; empty/dot segments, `..`,
backslashes, query, fragment, userinfo, ports, unknown roots, abbreviated object ids,
and noncanonical escapes fail. References receive a deterministic media type from a
closed, case-insensitive extension table: `.md`/`.markdown` =
`text/markdown`, `.json` = `application/json`, `.yaml`/`.yml` =
`application/yaml`, `.txt`/`.log` = `text/plain`, and every other extension =
`application/octet-stream`. No host MIME database participates.

A blob OID alone was rejected: it can retrieve bytes but omits the Captain-required
repository path and revision context. A branch/ref name was rejected because it moves.
A remote URL was rejected because object acquisition is transport, not identity. The
full commit plus path is the smallest immutable history address, and raw SHA-256 remains
the independent byte check.

This is one Review v1 URI profile, not a generic URI framework. Existing clean relative
URIs keep their existing meaning for files genuinely owned by the retained room: they
resolve from the Briefing directory with containment checks. `gate prepare` emits
Git-root locators for its selected repository files and does not copy them. The
request's Briefing locator remains a clean room-relative URI because
`gate-briefing.json` is room-owned.

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
Git-root URIs have identical full bytes and basename slug, append stable `-1`, `-2`, ...
suffixes in locator lexical order while preserving caller order in the Briefing; an
exact repeated normalized input path is rejected. Authoritative revisions always retain
the full digest. The Briefing id, gate id, attempt id, JCS Briefing digest, request JCS
digest, room reference, and Captain actor/approver are binary-owned.

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

Cycle 6 adds no request field and no command flag. The schema delta is the one closed
`git-root` Artifact/Reference URI profile; raw-byte SHA remains Review v1's existing
`rev`. The CLI constructs the current `main`/`state` root map from the same resolved
definition and entity roots already used for recorder operations and passes it to
`internal/gitsource`; it never serializes checkout paths. No current provider handoff
consumes that map. This task stops at recorder readiness; sibling
`git-root-review-v1-materialization` owns the actual resolved-byte consumer contract
instead of assuming q0 transport will carry roots.

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

Every operation receives the current workflow root map: `main` is the Git history
containing the definition root and `state` is the distinct history containing the
resolved entity root. Physical checkout moves are handled by reopening the workflow at
its current definition/entity locations; no absolute checkout path enters the room.
Resolution always reads `git cat-file blob <commit>:<path>` from the named root's
existing local object database and verifies the Artifact/Reference raw SHA. It does not
read current worktree bytes, search neighboring directories, or fetch a remote. A
missing root/object/path or digest mismatch fails before mutation. A fresh clone works
when its local object database already contains the pinned commit; object acquisition
and remote identity remain outside this task.

That local-object requirement is a lifecycle precondition, not hidden recovery.
Recorder bind/close/validate each recheck that the full commit and addressed blob exist
in the named root. A shallow or partial clone is supported only when both objects are
already local; Spacedock never deepens, fetches, hydrates, or substitutes current
worktree bytes. The owning ref must remain reachable until the gate and all read-only
validation are complete. Deleting/rewriting its last containing ref can make Git prune
the object and makes the room fail closed; s4 creates no retention ref or cache. Cross-
machine object acquisition and any provider-side retention/materialization belong to
`rqh46ey33aqq4rt72b4w1m2q`.

Preparation stages the two-file candidate room in a sibling temporary directory,
validates its request, Briefing, and every Git-root/room-relative source through the
production loaders, atomically publishes the fresh room, and binds the attempt under the
same entity lock. Exact replay is a no-op. An occupied divergent room, stale entity
comparison, changed selected commit/bytes on replay, or handled resolve/bind/write
failure removes only the new candidate and leaves the entity and retained rooms
byte-identical. That is error atomicity while the entity lock is held, not cross-file
crash atomicity: a process or power loss after room rename but before entity replacement
can leave an unbound room. This task adds no journal or recovery schema; a later prepare
refuses divergent occupancy, and crash recovery remains an explicit operator concern.

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

For the no-override path, `fo-gate-lifecycle` first commits any newly authored selected
source in its owning Git history, then runs `gate prepare`, commits the entity folder
containing the two-file room and binding, passes the gate-review Artifact to the
rendering-only `present-gate`, and records the captain's semantic chat decision. The
source commit makes each selected worktree byte an immutable reachable object; the
later folder-scoped state commit includes only the generated metadata room and binding
without taking its stdout path as an argument.

6y's Cycle-31 authority-capture reset remains a separate authority boundary. This task
records only locators for Artifact and Reference content explicitly selected for
presentation; it does not turn `--reference` into a delegated-conn channel, add an
authority flag, or let the First Officer author request actor/approver metadata. 6y owns
the smallest recorder/scaffold reference that captures exact Captain bytes without
agent retyping. Before s4 implementation, the landed 6y design must be compared with
this closed prepared request. If 6y requires a new request field, a generated authority
object, or different actor/approver derivation, s4 returns to ideation for a joint
request-surface reset. It must not smuggle authority through a content Reference.

The lifecycle text must call the output **recorder-ready**, not presentation-ready.
No-override chat remains complete because the First Officer already has the selected
filesystem paths and renders the gate review through `present-gate`. A selected override
must halt as unavailable before provider invocation until
`git-root-review-v1-materialization` lands; neither the room nor stdout is a promise that
current Subspace can render Git-root content. This repository does not discover,
version-check, capability-probe, launch, or materialize for a provider. There is no
selected-override execution arm in this task. Subspace q0 still owns transport and
retained preflight, but sibling `rqh46ey33aqq4rt72b4w1m2q` owns resolved-byte
presentation and the real cross-repository proof. Spacedock tests do not fake either
future event.

## New mechanisms and rejected alternatives

| Mechanism | Value AC | Simplest alternative | Why insufficient |
|---|---:|---|---|
| One `gate prepare` operation | AC-1 | Tell the FO to write two JSON files and call `gate record` | Preserves the manual ids/digests and partial-room failure that caused the task. |
| Closed Git-root locator and local-object resolver | AC-1, AC-5 | Copy selected bytes into the room | Copies survive the reopen but duplicate objects the Captain requires to remain singular; path plus raw SHA alone cannot locate bytes after checkout movement. |
| Explicit presentation dependency `rqh46ey33aqq4rt72b4w1m2q` | AC-4 | Say q0 will carry logical roots later | Transport does not turn Git objects into filesystem bytes current Subspace accepts; only a real consumer/materialization API plus end-to-end presentation can close that gap. |
| Frozen local Briefing locator | AC-2, AC-5 | Keep joining `briefing.json` | Fails the reproduced valid room and contradicts the provider contract. |
| One recursive duplicate-member reader | AC-3, AC-5 | Rely on `encoding/json` plus typed structs | Go accepts conflicting duplicates last-wins; the detached counterexample can close under the wrong authority. |
| Stable room/identity stdout handoff | AC-1 | Omit the room or make callers reconstruct it from ids/directory layout | Hides the published artifact and can select the wrong attempt under retries. |
| In-memory derived association | AC-5 | Persist `association.json` | Creates a second durable truth that can diverge from the four frozen inputs. |

## Expected surface and tolerance

Baseline assumption: latest 6y (`60adfc1f`, including lifecycle owner `e9415a17`) lands
first. Relative retained-input normalization is then available in `internal/cli`, the
existing recorded-gate journey targets `fo-gate-lifecycle`, and `present-gate` contains
rendering only. Against that composition, the smallest expected implementation is these
18 files and about `+1,413/-161` lines (**1,574 changed LOC**):

The inspected 6y tip is still pre-xb-rebase, so implementation must not start until
6y's final xb rebase lands. Re-read that landed tip before creating the worktree; if it
changes lifecycle ownership, recorder commands, shared live assertions, exact-authority
capture, the prepared request boundary, or any declared file/delta below, return to
ideation for a surface reset rather than implementing against this provisional
composition.

| File | Expected delta | Purpose |
|---|---:|---|
| `internal/cli/cli.go` | `+60/-6` | Route `gate prepare`, normalize inputs, and print the stable four-line result. |
| `internal/cli/gate_test.go` | `+190/-25` | Reuse CLI fixtures for preparation, stdout, committed-source refusal, arbitrary-locator eligibility, and byte-clean failures. |
| `internal/gates/prepare.go` (new) | `+220/-0` | Derivation, recorder-ready Git-root locators, ids/media types, and error-atomic room publication. |
| `internal/gates/prepare_test.go` (new) | `+210/-0` | Focused replay, collision, locator selection, two-file room, and handled-error tests. |
| `internal/gitsource/source.go` (new) | `+175/-0` | Closed root/commit/path URI grammar, root map, ref-retention checks, common-history classification, and local `git cat-file` resolution. |
| `internal/gitsource/source_test.go` (new) | `+170/-0` | Independent moved-checkout, linked-worktree classification, later-worktree, shallow/pruned/missing-object, escaping, and raw-SHA controls. |
| `internal/gates/operation.go` | `+80/-35` | Closed request locator and the one exact Briefing resolver. |
| `internal/gates/application.go` | `+12/-4` | Route reviewed-input eligibility through that resolver instead of `briefing.json`. |
| `internal/gates/io.go` | `+30/-8` | Recompute the four retained provider inputs through duplicate-safe reads. |
| `internal/gates/json.go` (new) | `+75/-0` | Recursive duplicate-member rejection. |
| `internal/gates/testdata/gate-room/request.json` | `+1/-0` | Add the locator to the canonical fixture. |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | `+28/-16` | Add the selected-source commit before prepare to the shared no-override observation; add no lane. |
| `internal/contractlint/fo_function_reference_invariant_test.go` | `+22/-8` | Pin lifecycle ownership, Git-only process boundary, forbidden provider mechanics, and rendering-only `present-gate`. |
| `docs/specs/gate-resolution-frontmatter-contract.md` | `+85/-30` | Normative recorder-ready prepare/request/Git-root/retention/resolver/atomicity contract. |
| `docs/site/reference/command-reference.md` | `+20/-6` | New verb, committed-source locator, stdout, and arbitrary-name recording. |
| `docs/site/reference/frontmatter-contract.md` | `+3/-3` | Remove the manifest-basename claim. |
| `docs/site/concepts/gates-and-decisions.md` | `+10/-6` | Committed-source preparation, recorder-ready boundary, and presentation dependency. |
| `skills/fo-gate-lifecycle/SKILL.md` | `+22/-14` | Replace hand bind/association wording with source commit, prepare, recorder-ready recording, and provider halt. |

Tolerance is **+2 files and +25% changed LOC** (hard cap 20 files / 1,968 changed
LOC), for a focused resolver test or fixture split only. A change to
`skills/present-gate/SKILL.md`, any schema field in `gates`, new dependency, provider
executable/probe/transport, selected-provider harness, compatibility request shape,
caller output path, association artifact, or broader lifecycle routing requires a
design reset.

This surface deliberately excludes every file in
`git-root-review-v1-materialization`. That sibling owns any resolved-source manifest,
ephemeral materialization, provider API, Subspace change, and real provider E2E. If s4
implementation needs one of those to make its own recorder tests pass, the boundary is
wrong and returns to ideation rather than borrowing sibling scope.

## Acceptance criteria

**AC-1 (VALUE) — One command turns judgment and file choices into a validated open
gate room with zero caller-authored metadata files.** Starting baseline: the fixture has
one gate-review Markdown file, selected References, and no room. After `gate prepare`,
the derived room contains exactly the request and located canonical Briefing—**2**
regular files regardless of selected-source count, with **0** duplicated source
payloads. The attempt is open and `gate validate` succeeds. The four stdout lines expose
the exact cleaned absolute room, Briefing id, digest, and open state; the required
Artifact is `text/markdown`, Reference media types follow the closed table, and every
selected source URI names its `main`/`state` root, full commit, and
repository-relative path while `rev` matches the Git object's raw bytes. *Test:* a real
split-root CLI/Git fixture commits selected main/state files, asserts pre-command
metadata count 0, post-command count 2, exact two-file room/stdout/locators/media types,
and then makes later commits that move both worktree paths. It relocates main and state
checkouts independently, supplies their new logical root map, and revalidates the exact
old bytes from both local object databases. It fails if the fixture supplies metadata,
the room copies either source, a locator omits any identity component, output changes
under a launch directory containing spaces, or resolution needs the old `..` topology.
Dirty, untracked, third-repository, abbreviated-object, and missing-local-object cases
exit nonzero with byte-identical entity/room state. A retained ordinary ref containing
the source commit remains green after worktree movement. Reflog-only/detached
unreachable commits, a shallow clone missing the commit or blob, and a fixture that
deletes the last containing ref and prunes the object all fail closed without fetch,
deepen, hydration, worktree fallback, or state mutation.

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

**AC-4 — Spacedock truthfully produces a recorder-ready room and does not claim current
provider presentation.** With no override, 6y's `fo-gate-lifecycle` runs prepare once,
after committing every newly authored selected source, commits the entity folder
containing the generated two-file room and binding, renders through the unchanged
rendering-only `present-gate`, and records chat. With a selected override, the
lifecycle halts before provider invocation with an unavailable diagnostic naming
`git-root-review-v1-materialization`/`rqh46ey33aqq4rt72b4w1m2q`; it makes no room/root
handoff promise because current Subspace cannot consume `git-root://` Artifacts or
References. That sibling owns the resolved-byte consumer/materialization contract and
an actual Spacedock-to-Subspace presentation proof before the 0.27 pre-release. No
Spacedock Go or skill change in s4 names a provider executable, runs a provider
availability/version/capability probe, allocates or mutates provider outputs, fetches
source objects, materializes presentation files, or simulates invocation. *Test:* the
shared no-override recorded-gate observation requires
`gate-help → selected-source-commit → prepare → room/binding-commit → chat-render →
decision-record`, followed by 6y's unchanged
decision-commit/consume/consumed-commit barriers. At final tip, the existing Claude,
Codex, and Pi live lanes must each satisfy that same observation. Legitimate structural
checks require `fo-gate-lifecycle` ownership, keep `present-gate` rendering-only, confine
process execution to `internal/gitsource`'s literal `git` object reads, and reject
`subspace-tui`, `/subspace:r`, `--supports`, `--version`, remote fetch, or process-launch
imports elsewhere in the changed gate/lifecycle surface. A selected-override command
fixture requires zero provider invocations, byte-identical state, and the dependency id
in its diagnostic; deleting that halt must fail the fixture. No fake success lane or
prose-derived room-consumption mutant is added. The sibling's acceptance gate, not s4,
owns actual presentation evidence.

**AC-5 — Provider recording has one recomputed association and no parallel durable
artifact.** The full fixture prepares, receives fixed Result/inventory outputs, closes,
and validates with no `association.json`. Request, located Briefing, Result, and
inventory are each deleted and byte-mutated independently. Each selected Git locator is
also changed one component at a time: root, full commit, path, and raw
SHA. Every variant fails recording or read-only validation without changing the entity,
while a later worktree edit/move remains green because it does not change the pinned
object. *Test:* real CLI end-to-end fixture asserts the four provider digest pins, each
Git object/raw-SHA pair, and the exact room tree; adding an association input/file,
copying a source payload, omitting a local object, or changing one locator component
fails.

## Test plan and proof order

0. **Baseline gate, before implementation:** require 6y's final xb rebase and
   Cycle-31 authority-capture design to be landed, record its tip, and compare lifecycle
   ownership, recorder commands, exact-authority capture, prepared-request fields,
   shared assertions, and the expected-surface table. Any mismatch returns to ideation
   for reset; a Captain conn is never passed as an ordinary frozen Reference.
1. **Focused red/green, low cost:** add the arbitrary-name spike as the first command
   test using the existing gate-room fixture, then add focused `prepare_test.go` cases
   for exact stdout data, Git-root URIs/media types/raw digests,
   12-to-64-character digest-prefix extension/full-digest suffixes, two-file replay,
   occupied room, and handled-error cleanup. `internal/gitsource` first gets the
   split-root Git fixture that commits both sources, advances/moves their worktree paths,
   relocates both checkouts independently, and resolves the old commits through the new
   root map with the original layout unavailable. Add ordinary-ref reachability,
   reflog-only/detached-unreachable, shallow/partial missing-object, and simulated
   last-ref deletion plus prune controls; every missing-object lane must prove there
   was no fetch, deepen, hydration, or worktree fallback. Run
   `go test ./internal/gitsource ./internal/gates ./internal/cli -count=1`.
2. **Adversarial JSON, medium cost:** mutate each of the four room documents at top
   level and nested authority-bearing objects. Assert entity bytes and lock state, not
   only error substrings. The arbitrary-locator positive case continues through provider
   room closure and CLI eligibility, so `application.go` cannot silently retain its
   basename join.
3. **Existing no-override FO journey only, high cost at final tip:** update 6y's shared
   recorded-gate observation in place to assert
   `gate-help → selected-source-commit → prepare → room/binding-commit → chat-render →
   decision-record` before its unchanged decision-commit/consume/consumed-commit
   barriers. Run the existing
   `TestLiveClaudeSharedScenarios` and `TestLiveCodexSharedScenarios`
   `recorded-gate-lifecycle` cases plus `TestLivePiRecordedGateLifecycle` against the
   final implementation tip; all three must observe the revised sequence. Add no host
   lane, harness, provider fake, provider-capability ledger, prose-derived room mutant,
   selected-provider success lane, or cross-repo invocation test. A selected override
   halts before invocation and names `rqh46ey33aqq4rt72b4w1m2q`; that sibling owns the
   actual provider-neutral resolved-byte API and Spacedock-to-Subspace E2E.
4. **Repository gates:** `gofmt -w ./cmd ./internal`, `go test ./...`,
   `go test ./... -race`, strict docs build, `git diff --check`, and verify
   `go list -deps ./cmd/spacedock` contains no Subspace package, only
   `internal/gitsource` invokes a process and its argv begins with literal `git`, no
   `fetch`/`clone` command exists in that package, and lifecycle command text contains
   none of `subspace-tui`, `/subspace:r`, `--supports`, or `--version`.
5. **High-stakes detached audit:** independently inject conflicting duplicate
   `by`, locator, digest, id, and inventory members and try to refute the byte-clean
   claim. Move independently cloned main/state roots again, advance their worktrees past
   moved/deleted source paths, and reopen the pinned commits through only the new root
   map. Any worktree lookup, neighboring-directory search, remote fetch, copied payload,
   missing identity component, or changed raw byte is a rejection. Re-check landed
   Subspace `em` only to confirm the ownership boundary; do not launch Subspace,
   materialize sources, or treat the local Git helper result as presentation evidence.
6. **Separate release dependency:** s4 may merge after its recorder-ready criteria
   pass, but the 0.27 pre-release remains blocked until
   `git-root-review-v1-materialization`/`rqh46ey33aqq4rt72b4w1m2q` demonstrates the
   actual consumer/materialization path and provider presentation E2E. Do not import
   that sibling implementation into s4 to manufacture a presentation claim.

## Documentation change proposal

The implementation applies these concrete semantics (line wrapping may follow the
target file):

```diff
--- docs/site/reference/command-reference.md
+++ docs/site/reference/command-reference.md
@@
-| `spacedock gate record <entity> --briefing PATH/briefing.json` | Bind a complete retained package manifest whose basename is exactly `briefing.json`. Other basenames fail before mutation. |
+| `spacedock gate prepare <entity> --question TEXT --artifact REVIEW.md [--reference FILE ...]` | Derive and bind one recorder-ready two-file room. Selected files must be committed objects in the workflow's main or state Git history and retained by an ordinary ref; the generated Briefing records `git-root://<root>/<full-commit>/<path>` and raw SHA-256 `rev` without copying payloads. The required objects must already be local; no fetch, deepen, hydration, worktree fallback, or retention ref is created. Success prints exactly `room`, `briefing`, `digest`, and `state` key/value lines. Current Subspace cannot present these locators; `rqh46ey33aqq4rt72b4w1m2q` owns that pre-release dependency. |
+| `spacedock gate record <entity> --briefing PATH` | Bind any readable canonical Briefing by its exact path. A prepared room instead freezes its Briefing locator, id, and digest in `request.json`; every later operation resolves that locator rather than a canonical basename. |

--- docs/site/concepts/gates-and-decisions.md
+++ docs/site/concepts/gates-and-decisions.md
@@
-Before the First Officer shows a gate, it binds the exact retained Briefing and commits that package.
+Before the First Officer shows a no-override gate, it commits newly authored selected sources, prepares and binds the two-file room mechanically, commits the entity folder containing that room, then renders in chat. Git-root locators reopen exact committed objects through the current main/state root map after checkout movement. This is recorder-ready, not presentation-ready: a selected override halts before invocation and names `git-root-review-v1-materialization`/`rqh46ey33aqq4rt72b4w1m2q`, which owns the resolved-byte consumer and actual provider E2E. Spacedock does not fetch objects or discover, probe, launch, or materialize for a provider.

--- skills/fo-gate-lifecycle/SKILL.md
+++ skills/fo-gate-lifecycle/SKILL.md
@@
-**Retain and bind.** Assemble `ROOM/briefing.json` ... then run `gate record ENTITY --briefing BRIEFING`.
+**Prepare and bind.** Select one Markdown gate-review Artifact and any References; commit every newly authored selection in its owning main/state Git history, then run `${SPACEDOCK_BIN:-spacedock} gate prepare ENTITY --question QUESTION --artifact REVIEW [--reference FILE ...] --workflow-dir WORKFLOW_DIR`. Require the four stable output lines and `state=open`; commit the entity folder containing the generated two-file room and binding before presentation. The emitted `room=` value is the sole recorder/diagnostic locator and must never be reconstructed or searched for. Do not copy selected repository objects into the room.
+**Presentation boundary.** Call the result recorder-ready. With no override, render the generated review through `present-gate` and record chat. With a selected override, halt before invocation and name `git-root-review-v1-materialization`/`rqh46ey33aqq4rt72b4w1m2q`; current Subspace cannot consume the Git-root Artifacts/References. This lifecycle does not promise a room/root handoff, name/discover/version-check/capability-probe/launch a provider, materialize sources, or construct or mutate provider output paths.
```

The normative spec makes the same substitutions, defines the closed request shape and
room publication/error-atomicity behavior, defines the closed
`git-root://<root>/<full-commit>/<repo-path>` URI grammar, ordinary-ref source-commit
rule, and local-only shallow/partial/prune failure behavior, and states explicitly that
association is recomputed and unstored. It labels the output recorder-ready and names
the separate presentation dependency. The frontmatter reference removes only its
exact-basename claim; no `gates` schema or request field changes.
`skills/present-gate/SKILL.md` remains unchanged and rendering-only after 6y.

## Git-root and authority boundaries

Preparation records an immutable reference to explicitly selected presentation content;
it does not snapshot or duplicate that content. The logical root and commit/path locate
the existing Git object, and the full raw `rev` verifies the bytes returned from the
local object database. Current worktree paths and bytes are non-authoritative after
preparation. This is deliberately smaller than a repository transport: there is no
remote URL, ref name, clone/fetch, checkout id, provider path, or fallback search.

This object reference does not solve 6y's delegated-authority problem by accident.
Exact Captain authority is not caller-authored metadata and is not a presentation
Reference. 6y's recorder/scaffold owns that capture. If its landed correction composes
with the current closed request, s4 remains unchanged; if it needs the prepared request
to name an authority object, that is an explicit pre-implementation design reset with
joint acceptance evidence.

## Release dependency and readiness language

`room` publication, `state=open`, and a successful `gate validate` mean
**recorder-ready** only. s4 adds no `presentation-ready` flag, output, or implication.
The no-override chat path is complete because it renders the already selected review
file; a selected provider override is unavailable and halts before invocation.

The backlog entity `git-root-review-v1-materialization`
(`rqh46ey33aqq4rt72b4w1m2q`, durable-decisions) is the named release dependency. It owns
the provider-neutral API that turns the frozen Git address into bytes current provider
packages can consume, any required ephemeral materialization and retention policy, and
the first actual Spacedock-to-Subspace presentation E2E. s4 may merge when its
recorder-ready acceptance criteria pass, but 0.27 pre-release cannot proceed until that
sibling's presentation gate passes. A root map merely passed through q0 is not evidence
for that gate.

## Out of scope

- Subspace q0, room-to-provider invocation, Git-root consumer/materialization,
  terminal transport, provider discovery or capability probing, provider output
  allocation/mutation, retained-preflight proof, provider retention implementation, or
  presentation E2E; these presentation concerns are owned by
  `rqh46ey33aqq4rt72b4w1m2q`.
- Compatibility request parsing, a `briefing.json` fallback for prepared requests, or
  migration wrappers.
- `association.json`, caller-selected Result/log/inventory/diagnostic paths, or provider
  argv.
- Broader lifecycle-next-action prose, advisory-round preparation, readiness projection,
  crash-atomic multi-file transactions, copied selected sources, remote object
  acquisition, configurable root registries, authority transport through `--reference`,
  or a generic URI/JSON framework.

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

- **Cycle 6 — Captain Subspace rejection of frozen copies (2026-07-25).** The
  presentation prototype retained one arbitrary `gate-briefing.json`, a closed
  locator-bearing request, and frozen gate-review/staff-review/entity sources. Subspace
  returned `revise`: “no, i thought we'd reference the object in one of the split-root
  repo instead of copying. they should be addressible and we just need a
  name/path/sha.” Reason: “i don't like duplicating files.” Replace copied sources with
  one provider-neutral split-root object reference that names the owning root
  (`main`/`state` or the workflow-defined equivalent), repository-relative path, and
  immutable Git object/revision plus raw-byte SHA. Define how a resolver reopens the
  exact bytes when checkouts move independently, without depending on a `..` topology,
  remote URL, provider transport, or compatibility layer. The exact provider package
  is retained at
  `/private/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/subspace-r-provider.FAUeX8`.
  Do not consume the pending approval; a corrected Briefing must supersede it.

- **Cycle 6 staff rejection — recorder-ready is not presentation-ready
  (2026-07-25).** `git-root://root/commit/path` plus raw `rev` can identify and verify
  the historical blob for Spacedock, but current Subspace package mode receives the
  Briefing unchanged, treats Artifact URIs as filesystem paths, and rejects Reference
  URIs containing `://`. The room carries no logical-root map, and Spacedock's
  `git cat-file` resolver is not in the provider presentation path. Decide and name the
  cross-repository contract: either narrow s4 honestly to a recorder-ready
  Git-addressed room and file the resolved-byte/provider work as a release dependency,
  or define the actual consumer/materialization API and exercise it end to end before
  claiming presentation readiness. Rebaseline the surface around that owner. Drop
  `<object-format>` unless a concrete ambiguity proves it necessary; root + full commit
  + repository-relative path plus raw SHA is the smaller Captain-directed form. Keep
  local-object retention, dirty/untracked/third-repository rejection, and production
  moved-checkout resolution explicit.

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

## Stage Report: ideation (cycle 5)

- DONE: Replace layout-dependent cross-checkout `..` source URIs with the smallest durable Git-addressed or room-frozen source model while retaining Briefing-relative URIs for room-owned files.
  Every selected Artifact/Reference is now frozen below `sources/`, addressed by one clean contained Briefing-relative URI, and verified by its full raw SHA; the design adds no Git locator or generic URI schema.
- DONE: Add a falsifiable independent-main/state-checkout reopen case that verifies exact source resolution and bytes after ordinary state commit, without provider-specific transport.
  The spike made 6y's concrete `../../../../../../../docs/specs/gate-resolution-frontmatter-contract.md` disappear after an unrelated-path state clone while the frozen copy reopened at the same `6425…9003` SHA with exact bytes; the implementation fixture fails on any escape, absence, or byte drift.
- DONE: Reconcile the correction with 6y's authority-capture boundary and publish a revised surface/LOC/test plan ready for a new ideation gate.
  s4 freezes presentation content only and cannot carry Captain conn through `--reference`; landed 6y exact-authority capture is a pre-implementation request-boundary check, with a reset required for any request-field/authority-file change. The provisional surface remains 16 files at +1,090/-161 (1,251 changed LOC), with an 18-file/1,564-LOC cap.
- SKIPPED: Consume the existing ideation approval or implement code.
  Cycle 5 is a state-only design correction; it requires a fresh ideation gate and no provider, code worktree, gate application, or implementation test was launched.

### Summary

The revised design removes checkout topology from Review v1 resolution by making every
selected source a retained room-owned file. The independent reopen spike falsified the
old URI and preserved exact bytes under the frozen model; the authority and provider
boundaries remain explicit, and the corrected surface is ready for a new ideation gate.

## Stage Report: ideation (cycle 6)

- DONE: Reference the existing object in one split-root repository instead of copying selected Git-owned sources.
  Prepared rooms return to exactly two metadata files; selected Artifact/Reference payloads remain singular committed objects in the workflow's main or state Git history.
- DONE: Define the smallest provider-neutral locator using owning root name, repository-relative path, immutable Git object/revision, and raw-byte SHA.
  The closed `git-root://<root>/<object-format>/<full-commit>/<repo-path>` URI carries history identity while Review v1's existing `rev` carries the independent raw SHA-256 byte pin.
- DONE: Define exact recovery and independent-checkout movement semantics.
  Resolution uses `git cat-file blob <commit>:<path>` in the current logical root's existing local object database, never current worktree bytes, neighboring paths, remote fetch, or a serialized checkout location.
- DONE: Preserve Briefing-relative URIs for genuinely room-owned artifacts without duplicating selected entity/spec files.
  Clean relative URIs still resolve from the Briefing with containment; generated request-to-Briefing locator stays relative, while selected repository files receive Git-root URIs and no payload copy.
- DONE: Preserve arbitrary Briefing location, duplicate-member rejection, derived association, and provider-neutral ownership.
  The request locator, recursive JSON reader, four recomputed provider pins, unstored association, no-override lifecycle, and q0/Subspace transport boundary remain unchanged.
- DONE: Spike an independent-checkout reopen without `..` topology.
  Separate main/state clones moved to unrelated depths and later worktree commits moved both selected paths; the pinned objects still resolved to `6425…9003` and `f019…235f`, while the old `..` URI did not resolve.
- DONE: Recalculate schema, command, LOC, and test surface.
  No request field or CLI flag is added; the URI profile and committed-source precondition produce an 18-file +1,439/-161 plan (1,600 changed LOC), capped at 20 files/2,000 LOC, with focused Git-root and moved-checkout proof.
- SKIPPED: Implement, consume the pending ideation approval, launch Subspace again, or modify provider outputs.
  Cycle 6 changes only this state design and requires a superseding ideation Briefing; the retained provider package and all provider-owned outputs were left untouched.

### Summary

Cycle 6 replaces the rejected copy model with one logical-root, commit/path, raw-SHA
reference to each existing Git object. The moved-checkout spike exercises local-object
recovery without topology or provider transport, and the revised design is ready for a
superseding ideation gate.

## Stage Report: ideation (cycle 7)

- DONE: Narrow s4 from presentation-ready to recorder-ready after checking current Subspace package semantics.
  Current Subspace receives the Briefing unchanged, treats Artifact URIs as filesystem paths, and rejects Reference URIs containing `://`; s4 therefore proves preparation, binding, recording, and validation only, while a selected override halts before invocation.
- DONE: File and name the actual provider-presentation dependency before the 0.27 pre-release.
  `git-root-review-v1-materialization` (`rqh46ey33aqq4rt72b4w1m2q`) owns the resolved-byte consumer/materialization contract, provider-side retention as needed, and the first actual Spacedock-to-Subspace E2E; merely carrying logical roots through q0 cannot satisfy it.
- DONE: Reduce the durable source address to the Captain-directed form and make source-lifecycle preconditions falsifiable.
  The v1 locator is `git-root://<root>/<full-commit>/<repo-path>` plus raw SHA-256 `rev`; source bytes must match committed history, the commit must remain on an ordinary local/remote-tracking/tag ref, and missing shallow/partial/pruned objects fail closed without fetch, deepen, hydration, cache, retention ref, or worktree fallback.
- DONE: Rebaseline acceptance criteria, docs, and implementation tolerance around the recorder/provider ownership split.
  AC-4 and the FO/docs wording now call the result recorder-ready, require zero provider invocation for selected overrides, and name the dependency diagnostic. The implementation remains 18 files at +1,413/-161 (1,574 changed LOC), capped at 20 files/1,968 LOC, and excludes every sibling materialization/provider file.
- DONE: Distinguish the independent Git helper spike from production and presentation evidence.
  The spike proves only that a root/full-commit/path address can recover exact bytes from moved local object databases; it did not exercise `gate prepare`, the production resolver, Subspace, or an end-to-end presentation.
- DONE: Preserve the accepted recorder authority boundaries.
  Arbitrary frozen Briefing locators, recursive duplicate-member rejection, four recomputed provider pins, unstored association, two-file rooms, no-override chat, and 6y's separate exact-Captain-authority seam remain in scope.
- SKIPPED: Implement code, consume or apply a gate, launch Subspace, materialize sources, invoke a provider, or mutate provider outputs.
  This cycle is a state-only ideation correction; the retained provider package remains read-only and a superseding ideation gate is required.

### Summary

Cycle 7 makes the readiness claim match the current consumer: s4 produces a durable
Git-addressed recorder room, not a currently presentable Subspace package. The filed
dependency now owns the missing consumer and real E2E, while s4 retains a smaller
root/full-commit/path locator, explicit local-object retention rules, and a rebaselined
recorder-only proof surface.
