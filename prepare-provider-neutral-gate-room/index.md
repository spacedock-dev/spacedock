---
id: s4ykctf21g60dvfgdd6cy9ny
title: Prepare provider-neutral gate rooms and align canonical Briefing recording
status: implementation
source: "Durable-decisions cross-repo dogfood ruling after xb and Subspace em review, 2026-07-24"
started: 2026-07-24T14:54:10Z
completed:
verdict:
score: "1.0"
worktree: .worktrees/spacedock-ensign-prepare-provider-neutral-gate-room
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
                state: consumed
                blockers: []
            - id: gate-attempt:s4-ideation-3
              briefing:
                id: briefing:docs-dev:s4:ideation:attempt-3:revision-1
                digest: sha256:d80b23af1136b1caffb1786878d98b1933799f38f25fbcb99dee36466cea3469
                digest-domain: canonical-bytes
                room-ref: ./review/ideation/briefing-3
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:s4:ideation:3
                briefing: briefing:docs-dev:s4:ideation:attempt-3:revision-1
                by: agent:first-officer
                at: "2026-07-25T17:34:49.721839Z"
                decision: approve
                reason: Cycle 11 provides one mechanical provider-neutral prepare command, collision-free folder and flat rooms, exact Git-root and Briefing authority, no copied sources or stored association, and the sole room-only provider handoff; independent re-review found no remaining issue.
              application:
                action: advance
                target-stage: implementation
                state: consumed
                blockers: []
review-round:
    id: round:s4ykctf21g60dvfgdd6cy9ny:implementation:13
    stage: implementation
    cycle: 13
    briefing:
        id: briefing:s4ykctf21g60dvfgdd6cy9ny:implementation:round-13
        digest: sha256:d25e3d5cf7887d1d30a12572b9f44de34c385ab7bb7b37657828ed7382c79d5d
        digest-domain: canonical-bytes
        room-ref: ./review/implementation/round-13
---

## Cycle 11 governing design

This section carries forward the independently approved Cycle 10 room, schema,
Git-root, summary, and association design with the Cycle 11 planning and current-main
owner corrections. It supersedes all cycles 1–10 design, surface table, acceptance
criterion, test plan, documentation proposal, and ownership statement retained later
in this file. Those sections remain only as the correction record. The feedback
history and prior stage reports remain authoritative history, but they do not
authorize implementation.

Make gate preparation one mechanical, provider-neutral operation. The First Officer
supplies a question, one Markdown gate-review Artifact, its concise summary, and any
supporting files. Spacedock derives the entity-owned room, gate/attempt/Briefing/source
identities, immutable Git-root locators, raw SHA-256 revisions, request authority,
canonical digests, and open binding. The prepared room is recorder-ready and contains
exactly two authoritative metadata files at preparation time; selected source payloads
remain singular Git objects. A later provider is allowed to add its owned evidence
subtree without falsifying that prepare-time invariant.

### Current baseline and authorization reset

This design is based on current Spacedock `main`
`4ff98d8cd97ebcf17b6a583070ce69234e24fc87`, with the recorded First Officer lifecycle
landed at `deac7f8a`. In that composition:

- `skills/fo-gate-lifecycle/SKILL.md` owns capability preflight, retaining/binding a
  gate, durable state commits, semantic recording, consumption, and routing;
- `skills/present-gate/SKILL.md` owns chat rendering and presentation-channel
  overrides, including the public `/subspace:r gate <room>` form;
- `internal/gates` owns request/Briefing validation, room-backed Result recording,
  derived association, provider-evidence pins, and read-only validation; and
- `spacedock state commit` already scopes folder entities as one directory but scopes a
  flat entity as only its Markdown file.

The only existing ideation approval is attempt 2 for the rejected frozen-copy
architecture. It remains immutable history and is never reinterpreted, reopened, or
consumed again. After this correction and an independent review, the First Officer
must bind a fresh attempt-3 Briefing whose canonical digest covers this governing
design and the new review. That later gate operation is First Officer frontmatter
authority, not ideation-worker work.

### Problem and value

Current `gate record --briefing` still makes the First Officer handcraft canonical
JSON, ids, revisions, and a room, and it rejects a valid canonical Briefing unless its
basename is `briefing.json`. Room-backed recording later appends that basename again.
The same assumption reaches validation and application eligibility. A prepared
request already has enough authority to name an arbitrary canonical Briefing; every
operation should use that frozen locator.

The old source layout also crosses from state to main with `..` paths. That works only
while one checkout happens to be nested under the other. The Captain selected the
smaller durable identity instead: logical root (`main` or `state`), full immutable
commit, repository-relative path, and independent raw SHA-256. No source copy is
permitted.

The end value is observable: one command turns judgment plus file choices into one
validated open room with zero caller-authored metadata files and zero copied source
payloads, for either supported entity form. The room remains the sole public handoff
to future provider presentation.

### Current-main owner map

| Boundary | Single owner | Required behavior |
|---|---|---|
| Lifecycle capability preflight | First Officer via `fo-gate-lifecycle` | Resolve the launcher once and run exactly one fresh `gate --help`; require `prepare`, `record`, `validate`, `eligibility`, `consume`, and every flag used by prepare: `--question`, `--artifact`, `--summary`, `--reference`, and `--workflow-dir`. Stale help halts before source commits, preparation, or state effects. |
| Question, primary summary, and selected files | First Officer via `fo-gate-lifecycle` | After the capability preflight, author human judgment, commit newly authored selections, call `gate prepare`, and never author JSON or reconstruct ids/paths. |
| Entity resolution and room placement | `status.ResolveActivePath` plus `internal/gates` | Resolve the canonical folder/flat entity once and derive its collision-free room under lock. |
| Git source identity and bytes | new `internal/gitsource` package | Classify `main`/`state`, address `<full-commit>:<repo-path>`, compare selected bytes, and reopen only local objects. |
| Canonical request, Briefing, ids, digests, and binding | `internal/gates` | Publish the two metadata files, validate them through production readers, and bind the attempt atomically on handled outcomes. |
| Flat entity commit/archive durability | existing state/status machinery | Commit and archive the flat Markdown plus its exact companion directory without sweeping siblings. |
| Chat versus override routing | `present-gate` | Render chat when no override is selected. After prepare succeeds and the bound room is committed, a selected override performs exactly `/subspace:r gate <room>` with the emitted room and performs no provider probe or fallback selection. |
| Git-root materialization and provider invocation | `git-root-review-v1-materialization` (`rqh46ey33aqq4rt72b4w1m2q`) | Its fixed room-only entry derives every internal argument, owns provider discovery/capability/failure semantics, resolves sources, invokes Subspace, and retains provider evidence. |
| Provider Result recording | existing `gate record --room` | Re-read the room authority, derive association in memory, and pin exact provider Result/inventory bytes. |

There is no hybrid lifecycle owner. s4 changes `fo-gate-lifecycle` for capability
preflight and preparation, and `present-gate` only where its current-main channel
contract must stop naming `briefing.json`, probing a selected provider, selecting a
probe-driven chat fallback, or reconstructing authority. s4 does not launch, probe,
fake, or otherwise invoke a provider.

### Agent-facing command

The new command is:

```text
spacedock gate prepare ENTITY --question TEXT --artifact REVIEW.md \
  --summary TEXT [--reference FILE ...] [--workflow-dir DIR]
```

Before committing a newly authored selection or causing any preparation/state effect,
the First Officer performs the one fresh help preflight above. Missing `prepare` or
any of `--question`, `--artifact`, `--summary`, `--reference`, or `--workflow-dir` is
a stale launcher failure, not permission to continue with an older command or manual
frontmatter.

`--artifact`, `--question`, and `--summary` are each required exactly once.
`--reference` may repeat and preserves caller order. The Artifact filename ends in
`.md` or `.markdown` case-insensitively and receives `text/markdown`. References use
the closed, host-independent extension table:

| Extension | mediaType |
|---|---|
| `.md`, `.markdown` | `text/markdown` |
| `.json` | `application/json` |
| `.yaml`, `.yml` | `application/yaml` |
| `.txt`, `.log` | `text/plain` |
| every other extension | `application/octet-stream` |

Relative selected paths resolve against the invocation directory through the existing
CLI normalization seam before `internal/gates` runs. Repeating the exact normalized
input path is an error. Distinct paths with identical bytes remain distinct selected
sources.

The CLI rejects repeated `--summary` occurrences before Git calls or mutation, whether
their values match or differ, with:

```text
gate prepare accepts --summary exactly once
```

Invalid UTF-8 is rejected before Git calls or mutation with:

```text
--summary must be valid UTF-8
```

On success stdout is exactly:

```text
room=/clean/absolute/entity-room/path
briefing=briefing:<derived-id>
digest=sha256:<64 lowercase hex>
state=open
```

The cleaned absolute `room=` value is authoritative. No caller searches directories or
reconstructs it from the attempt number, Briefing id, entity, workflow path, status
output, or digest.

### Collision-free room placement

Both entity forms converge on the same entity-owned room home:

```text
# Folder entity
<state-root>/<slug>/index.md
<state-root>/<slug>/review/<stage>/briefing-<attempt-number>/

# Flat entity
<state-root>/<slug>.md
<state-root>/<slug>/review/<stage>/briefing-<attempt-number>/
```

For a flat entity, `<state-root>/<slug>/` is its companion artifact directory and must
not contain `index.md`; discovery therefore continues to resolve `<slug>.md`. A real
`<slug>/index.md` means the folder form wins under the existing conflict rule, so
preparation never ambiguously writes for both forms. Stage names and attempt numbers
come from the validated current gate, not caller path text. An occupied divergent room
is refused.

Folder state commits already scope `<slug>/`. Flat state commits must scope exactly the
two literal pathspecs `<slug>.md` and `<slug>/`, including tracked deletions, without
including any sibling. Flat archive/finalize moves both to
`_archive/<slug>.md` and `_archive/<slug>/`; the room reference
`./<slug>/review/...` remains valid from the archived Markdown's new parent. Archive
rollback and its path-scoped Git commit cover both paths as one operation. This is the
minimum durability extension that makes flat preparation real rather than an
uncommittable or orphaned preview.

Inline workflows use the same two layouts beside their README, but their ordinary
main-branch commit policy remains unchanged.

### Prepare-time authoritative files and later evidence

Under the entity lock, preparation builds a sibling temporary candidate, validates it
with the production request/Briefing/source readers, publishes it, and binds the open
attempt:

```text
<room>/
  gate-briefing.json
  request.json
```

Exactly these **2 regular files** exist immediately after successful preparation,
regardless of selected-source count. There are **0 copied selected-source payloads**,
no `association.json`, and no provider directory. `gate-briefing.json` is merely the
generator's chosen locator, never a canonical basename rule.

After presentation begins, the provider/materialization owner may add:

```text
<room>/provider/
  ... retained Result, inventory, log, diagnostics, and ephemeral-source lifecycle ...
```

The two-file claim is therefore a prepare-time authoritative-metadata invariant and
the no-override baseline, not a permanent whole-room count. Provider evidence cannot
replace or rewrite `request.json`, the located canonical Briefing, or any Git source.
The separate `minimize-gate-room-retention` task may narrow provider evidence by
outcome; s4 does not pre-decide that policy.

Handled validation, Git, bind, publish, or stale-comparison failure removes only the
new candidate and leaves the entity and all retained rooms byte-identical. Exact replay
is a no-op. A process/power loss between directory rename and entity replacement can
leave an unbound occupied room; s4 adds no journal or automatic deletion policy and a
later prepare refuses divergence.

### Local Git-object contract

Every selected file must be a readable, non-symlink regular file owned by exactly one
workflow root:

- `main` is the Git common directory that owns the workflow definition, including a
  linked implementation worktree;
- `state` is the distinct Git common directory that owns the resolved entity root; and
- an inline workflow emits only `main`.

Preparation asks the owning repository for its full `HEAD` object id and
repository-relative path, reads `<full-commit>:<repo-path>` from the existing local
object database, and requires those bytes to equal the selected worktree file.
Untracked, dirty, third-repository, symlink, absent-object, or mismatched selections
fail before room/entity mutation with an instruction to commit the exact source.

The closed locator is:

```text
git-root://<main|state>/<full-commit>/<escaped-repository-path>
```

Review v1's existing `rev` is the independent full raw-byte pin:

```text
sha256:<64 lowercase hex>
```

Repository path segments use canonical RFC 3986 escaping while `/` remains the
separator. Empty/dot segments, `..`, backslashes, query, fragment, userinfo, ports,
unknown roots, abbreviated commit ids, and noncanonical escapes fail. The owning
repository's full object-id width is accepted; s4 serializes no redundant object-format
tag, remote URL, ref name, checkout path, or provider path.

There is deliberately **no ordinary-ref containment test and no reflog policy**. A
clean detached checkout is valid when its full commit/path object exists locally and
the selected bytes match. Observing a containing branch cannot guarantee future
retention and refusing a detached object does not improve addressability. Every later
bind/record/validate/materialize read fails closed if the named root, commit, path,
blob, or raw SHA is unavailable or changed. Spacedock never fetches, deepens, hydrates,
searches neighboring paths, substitutes worktree bytes, creates a retention ref, or
uses a remote.

The prior independently moved main/state spike remains the risk proof: after both
worktrees moved and later commits removed the current paths, `git show
<commit>:<path>` reopened the exact pinned bytes while the old `..` locator failed.
No new spike is needed for Cycle 11: current entity discovery already distinguishes
flat `<slug>.md` from folder `<slug>/index.md`, and the implementation's first flat
test owns the new commit/archive behavior. The future room-to-provider consumer is not
an s4 premise; rqh owns and must exercise it end to end.

### Canonical Briefing, request, and exact summary

The generated Briefing is Review v1. Its binary-derived id is scoped by entity, stage,
attempt, and revision. Artifact/Reference ids are scoped by that Briefing and their
one-based canonical order, so they are deterministic and collision-free without a
digest-prefix extension algorithm. The required Artifact is first; References follow
caller order. Every item records its Git-root URI, media type, and full raw `rev`.

The primary Artifact alone carries:

```json
"summary": "  Résumé — validates Git-root presentation exactly.  "
```

The value must be a JSON string, valid UTF-8, and nonblank after a validation-only trim.
The canonical value is the exact caller-supplied logical string: Spacedock does not
trim, normalize Unicode, collapse whitespace, extract Markdown, summarize source
bytes, or synthesize prose. References receive no summary. Identity inventory,
association, Result, and the resolved-source manifest never duplicate it.

This mandatory-summary profile is selected only by an s4-prepared request. Existing
request-less `gate record --briefing PATH` and advisory-round Briefings remain valid
without a summary; their bytes are neither rewritten nor migrated. No compatibility
request shape or prepared-room basename fallback is added.

The exact string may contain any valid UTF-8 value representable in process argv (the
OS boundary excludes NUL), including whitespace, printable Unicode, and control
characters that the JSON encoder escapes. s4 preserves it; it does not render to a
terminal. The downstream rqh/Subspace owner must display printable Unicode and ordinary
spaces losslessly while rendering terminal-control/format code points through a
reversible visible escape.
It must never emit caller-authored controls into ANSI, a terminal title, a format
string, or another control channel. Safety changes only the view, never the canonical
Briefing string.

`request.json` has one closed v1 shape:

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

The `briefing.digest` is the canonical SHA-256 digest of the located Briefing. The
request itself is canonically digested and frozen on the binding. The binary, never
the caller, derives gate, attempt, Briefing id/digest, actor, and approver.

The locator is a clean, nonempty, slash-relative room path. Absolute paths,
empty/dot paths, `..`, backslashes, and symlink escapes fail before mutation. Nested
locators are valid. One shared resolver handles:

- request-backed bindings by validating `request.json`, resolving its exact locator,
  and checking id/digest; and
- request-less bindings by retaining and resolving the exact Briefing file reference.

Binding, room-backed recording, validation, and application eligibility all call that
resolver. None appends `briefing.json`.

### JSON authority and derived association

One s4 recursive token-stream reader rejects duplicate object member names before
typed decode or canonicalization in the request, located Briefing, Result, and
presented inventory. Detection covers every object depth, including extra fields;
Go's last-member-wins behavior is never authority.

Provider recording derives one in-memory `spacedock-result-association` from the exact
request, located canonical Briefing, Result, and presented inventory. It writes no
association file. The binding retains request and Briefing canonical digests; provider
evidence retains Result and inventory raw digests. Recording and read-only validation
recompute all four and fail without mutation on absence, substitution, duplicate
members, raw-byte drift, identity drift, or incomplete recursive inventory coverage.

### Room-only rq/Subspace handoff

The future public agent-facing provider command remains exactly:

```text
/subspace:r gate <room>
```

`<room>` is the exact `room=` emitted by `gate prepare`. The model supplies no entity,
workflow directory, Briefing path/id/digest, actor, approver, destination, provider
package, manifest, executable, capability, Result path, inventory path, or terminal
mechanics. The Subspace skill and fixed entry derive every internal value by validating
the bound room/request and using the current workflow discovery/root map. Possessing an
arbitrary directory is insufficient: the derived entity, current gate attempt, room
reference, request digest, Briefing locator/id/digest, actor, and approver must all
match before provider side effects.

The First Officer reaches this handoff only after successful preparation and the
state commit that durably binds the emitted room. It invokes the exact command once,
reconstructs no coordinate, and performs no availability, version, capability, or
fallback probe. The rqh fixed entry owns provider discovery, capability refusal,
launch failure, and retained-failure semantics after it receives the room.

rqh's closed resolved-source manifest must bind the canonical Briefing as well as its
source tuples:

```json
"briefing": {
  "id": "briefing:task:validation:attempt-1:revision-1",
  "digest": "sha256:<canonical Briefing digest>"
}
```

That pair is copied from validated room authority, then independently checked against
the unchanged canonical Briefing bytes. A manifest with correct source items but a
different Briefing id or digest fails before presentation. The manifest carries no
summary; Subspace reads the exact summary from the canonical Briefing. rqh's closed
manifest parser rejects duplicate object member names recursively before display.
rqh owns all materialization, capability probing, provider allocation/invocation,
cleanup, terminal-safe display, retained evidence, and the real moved-root
presentation E2E. s4 owns no provider call, fake, selected-provider success lane, or
copied provider fixture.

### Mechanism-to-value ledger

| Mechanism | Serves | Simpler alternative considered | Why insufficient |
|---|---|---|---|
| One `gate prepare` command | AC-1 value: zero caller metadata and one valid open room | Keep asking the FO to author JSON then call `gate record` | Preserves the manual error source and partial-room failure. |
| Folder/flat derived room home plus flat commit/archive unit | AC-1 value: both supported entity forms retain their room | Support folders only | Contradicts Captain authority and leaves ordinary flat tasks unable to use gates. |
| Exact request locator | AC-3 | Join `briefing.json` | Reproduces the observed arbitrary-basename failure. |
| Root/full-commit/path plus raw SHA | AC-2 | Raw SHA alone or room copies | Raw SHA cannot locate bytes; copies violate the Captain's no-duplication ruling. |
| Local-object failure without ref policy | AC-2 | Require an ordinary containing ref | The ref check rejects resolvable objects and cannot enforce later retention. |
| Required caller-authored primary summary | AC-3, AC-4 | Extract Markdown prose | Makes Spacedock invent or reinterpret human judgment. |
| Recursive duplicate rejection | AC-3, AC-5 | Typed `encoding/json` decode alone | Conflicting authority can otherwise be accepted last-wins. |
| Room-only `/subspace:r gate` and Briefing-bound manifest | AC-4 | Let the model repeat entity/workflow/identity argv | Reconstructed values can drift from the one frozen request authority. |
| In-memory association and four existing pins | AC-5 | Persist `association.json` | Adds a second durable truth that can diverge. |

## Advisory expected surface and semantic reconciliation

The named table is an advisory planning estimate against current `main` `4ff98d8c`,
not the pre-rebase 6y branch. It currently names **26 files, approximately
+1,717/-187 lines = 1,904 changed LOC**, but neither the per-file deltas nor those
totals authorize, reject, or reset an implementation:

| File | Expected delta | Purpose |
|---|---:|---|
| `internal/cli/cli.go` | `+80/-8` | Route/help/validate `gate prepare`, exact flag cardinality, path normalization, and four-line output. |
| `internal/cli/gate_test.go` | `+175/-18` | Folder/flat CLI, exact summary/stdout/errors, arbitrary locator, byte-clean command cases, and the complete `gate --help` prepare surface. |
| `internal/cli/state_sync.go` | `+34/-10` | Treat flat Markdown plus its exact companion directory as one literal commit unit. |
| `internal/cli/state_commit_test.go` | `+90/-4` | Prove flat room inclusion, deletions, sibling isolation, and replay. |
| `internal/gates/prepare.go` (new) | `+230/-0` | Derived layouts/ids/media, request/Briefing construction, validation, publication, and bind. |
| `internal/gates/prepare_test.go` (new) | `+250/-0` | Two-file/zero-copy layout, both forms, exact summary, ids, replay, occupancy, and atomic failures. |
| `internal/gitsource/source.go` (new) | `+125/-0` | Closed Git-root grammar, common-history classification, committed-byte comparison, and local object reads. |
| `internal/gitsource/source_test.go` (new) | `+145/-0` | Moved independent roots, linked/detached worktrees, escaping, missing object, raw SHA, and no-fetch controls. |
| `internal/gates/operation.go` | `+95/-30` | Request locator, exact-summary profile, duplicate-safe parse, and shared bound-Briefing resolver. |
| `internal/gates/application.go` | `+8/-4` | Eligibility uses the same resolver without basename inference. |
| `internal/gates/io.go` | `+20/-6` | Recompute request/Briefing/Result/inventory pins from exact retained inputs. |
| `internal/gates/json.go` (new) | `+60/-0` | Recursive duplicate-member refusal shared by authority documents. |
| `internal/gates/testdata/gate-room/request.json` | `+1/-0` | Add the required canonical Briefing locator. |
| `internal/gates/testdata/gate-room/briefing.json` | `+1/-0` | Add the request-backed primary Artifact summary. |
| `internal/status/mutate.go` | `+25/-7` | Move a flat entity and companion room subtree together on archive. |
| `internal/status/native_mutation_test.go` | `+50/-4` | Active/archive flat-room path and sibling controls. |
| `internal/status/merge.go` | `+40/-16` | Snapshot, rollback, stage, and commit both flat archive paths. |
| `internal/status/merge_guard_test.go` | `+75/-4` | Finalize/rollback/durability proof for flat rooms. |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | `+28/-14` | Reuse the real journey with one fresh complete help preflight, prepare, bind commit, and presentation; stale prepare help halts before effects and an agent-run provider probe is forbidden without adding a provider lane. |
| `internal/contractlint/fo_function_reference_invariant_test.go` | `+20/-6` | Pin the complete preflight, post-bind room-only handoff, current-main owner split, and forbidden provider probe/composition. |
| `docs/specs/gate-resolution-frontmatter-contract.md` | `+90/-26` | Normative prepare/layout/locator/summary/two-state/provider-boundary contract. |
| `docs/site/reference/command-reference.md` | `+22/-5` | Complete fresh-help preflight, new command, room layouts, exact output, arbitrary Briefing names, and local-object failures. |
| `docs/site/reference/frontmatter-contract.md` | `+3/-3` | Remove the `briefing.json` canonical-basename claim. |
| `docs/site/concepts/gates-and-decisions.md` | `+14/-5` | Mechanical prepare, exact summary, recorder-ready state, and room-only handoff. |
| `skills/fo-gate-lifecycle/SKILL.md` | `+26/-12` | Require complete fresh help before effects; own source commit, exact-summary preparation, bind commit, and room forwarding. |
| `skills/present-gate/SKILL.md` | `+10/-5` | Replace provider probe/fallback with the exact post-bind `/subspace:r gate <room>` handoff and prohibit caller-reconstructed authority. |

Before the first product edit, the implementation worker must declare the intended
paths, expected additions/deletions, and every known deviation from this table. That
declaration makes review drift visible; it is not a numerical budget or permission to
expand scope. Before implementation review, the worker must reconcile the path-scoped
`git diff --stat` and `git diff --numstat` against the declaration, explaining every
added, omitted, renamed, or materially different path in semantic owner/behavior
terms. The reviewer judges those explanations against this design and AC coverage,
never against a total or percentage.

No test, fixture, lint, worker instruction, or gate may assert aggregate file totals,
aggregate changed-LOC totals, percentage variance, or any other implementation-surface
count as a pass/fail condition. A numerical deviation alone neither authorizes
expansion nor returns the work to ideation.

Only semantic expansion resets the design: another public argument, caller-selected
authority, a third authoritative prepared metadata file, copied sources, remote
acquisition, stored association, compatibility behavior, weaker Briefing binding, or
a changed provider/validator owner. rqh's manifest and Subspace files remain outside
s4 surface.

## Acceptance criteria

**AC-1 (VALUE) — One command produces one durable recorder-ready room for either
entity form with zero caller-authored metadata and zero copied sources.** From a
split-root fixture containing only committed selected files and a folder or flat entity,
`gate prepare` creates/binds the derived open attempt, prints the exact four lines, and
leaves exactly request + located canonical Briefing in the room. Selected-source count
does not change the prepare-time file count of 2. The flat form uses
`<root>/<slug>/review/...`; `state commit <slug>` commits its Markdown and companion
room only, and archive/finalize moves both so the room reference remains readable.
*Test:* one table exercises both forms through real CLI/Git commits, replay, divergent
occupancy, sibling dirt, archive, and archive rollback. It counts files and payload
sentinels on disk, reopens the bound Briefing after archive, and fails if any metadata
was caller-authored, any source was copied, any sibling was committed, or the flat room
is orphaned.

**AC-2 — Every selected source is an exact, movable, local Git object without an
invented ref-retention rule.** Each URI names `main|state`, full commit, escaped
repository path, and a matching raw SHA-256 `rev`. Independently moved main/state
checkouts and later worktree path deletion still resolve the old bytes. A clean
detached checkout with the object present succeeds even when no ordinary ref contains
the commit. Dirty, untracked, third-repository, symlink, unknown-root, bad escape,
abbreviated-commit, missing-object, and raw-SHA mismatch cases fail without fetch,
fallback, room, or entity mutation. *Test:* a real Git fixture records argv and removes
the old topology; deleting any identity component, consulting current worktree bytes,
adding `for-each-ref --contains`, or issuing a remote/object-acquisition command makes a
positive/control case fail.

**AC-3 — The prepared request freezes an arbitrary canonical Briefing and the exact
request-backed primary summary.** Nested clean Briefing locators bind, room-record,
validate, and remain eligible without a `briefing.json` leaf. The primary summary is
the exact caller string, including the two-space/non-ASCII sentinel; References remain
summary-free. Repeated flags and invalid UTF-8 fail deterministically before Git.
Request-less and advisory summary-free Briefings keep their prior behavior and bytes.
Duplicate members at every depth of request, Briefing, Result, and inventory fail
before mutation. *Test:* command/adversarial tables continue an arbitrary locator
through close and eligibility, inject normalization-sensitive summary values and
last-wins duplicates, and assert whole-tree bytes plus lock absence on failure.

**AC-4 — The only future provider-facing authority is the prepared room.** Current-main
`fo-gate-lifecycle` first runs one fresh help preflight requiring `prepare` and all
prepare flags, then prepares and commits. Missing help capability fails before source
commit, preparation, or state effects. `present-gate` keeps chat as the unselected
default; after successful prepare and bind commit, a selected override performs exactly
`/subspace:r gate <room>`. It neither probes a provider nor selects a fallback, and it
reconstructs no room coordinate or authority. No s4 Go code or skill constructs a
provider executable, probe, materializer argv, or caller-reconstructed entity/workflow/
Briefing/actor/approver/destination. The rqh fixed entry owns provider capability and
failure behavior. Its manifest must bind the exact canonical Briefing id plus canonical
SHA-256 digest, preserve the canonical Briefing unchanged, exclude summary duplication,
reject manifest duplicates recursively, and render the exact string through a
reversible terminal-safe view. *Test:* extend the existing contract and lifecycle
fixtures so deleting `prepare` or any prepare flag from fresh help halts before every
source/preparation/state effect, and adding an agent-run availability/version/capability
probe fails. The no-override host journey observes prepare before bind commit and chat
presentation; the selected-override contract observes one exact handoff only after the
bind commit, without adding a provider lane. rqh's separate public-entry E2E must fail
when only its manifest Briefing id or digest changes, when its manifest has a recursive
duplicate, or when its caller surface adds any authority field.

**AC-5 — Recording derives one association from four recomputed pins and permits
provider-owned evidence only after preparation.** A request-backed room records fixed
Result/inventory evidence below `provider/`, validates, and closes with no
`association.json`. Mutating/deleting request, located Briefing, Result, inventory,
source locator component, raw SHA, primary summary, or recursive inventory coverage
fails recording/read-only validation without entity mutation. The room is two files
before provider work and may contain the owned evidence subtree afterward. *Test:* a
real CLI fixture hashes each phase, asserts all four durable digests and the exact tree,
then mutates each authority input independently; persisting an association or treating
the post-provider tree as a two-file invariant fails.

## Test plan and proof order

1. **Baseline and first reds:** pin `4ff98d8c`; add folder/flat `gate prepare` CLI reds,
   the arbitrary `decision-material.data` locator differential, exact/repeated/invalid
   summary cases, complete `gate --help` command/flag coverage, and the flat
   state-commit/archive tests before implementation. Extend the existing lifecycle
   preflight fixture first: removing `prepare` or any prepare flag from fresh help must
   halt before source commits, preparation, or state effects.
2. **Local objects:** implement the closed Git-root package against real independent
   main/state repositories. Prove linked and detached classification, moved roots,
   later worktree deletion, canonical escaping, full object ids, raw SHA, missing
   local objects, and recorded absence of fetch/deepen/hydration/worktree fallback.
   Do not add ref/reflog/prune policy lanes.
3. **Preparation and authority:** red/green exact two-file publication, ordinal ids,
   media types, replay/occupancy/error cleanup, exact summary, nested locator, and
   recursive duplicate rejection. Continue the locator positive case through
   `gate record --room`, `gate validate`, and eligibility; keep request-less/advisory
   controls summary-free.
4. **Flat durability:** run real path-scoped commit, ordinary archive, terminal
   finalize, and forced commit-failure rollback for flat companion rooms. Hash sibling
   dirt and prove the archived room reference still resolves. Folder behavior remains
   the unchanged control.
5. **Current-main lifecycle:** update the existing shared recorded-gate observation in
   place to one complete fresh help preflight, source commit, prepare, bind commit,
   present, record, and consume. Run its existing Claude, Codex, and Pi lanes at final
   tip for the no-override behavior s4 changes. Extend the existing contract fixture
   with the selected-override sequence: after the bind commit it emits exactly
   `/subspace:r gate <room>`, and a mutation that adds any agent-run availability,
   version, or capability probe fails. Add no provider lane, fake, new host, or
   prose-token proof.
6. **Repository gates:** `gofmt -w ./cmd ./internal`, `go test ./...`,
   `go test ./... -race`, strict docs build, and `git diff --check`. Verify no s4 Go
   dependency or process path names Subspace; only the local Git helper executes a
   process and its argv begins with literal `git` and contains no acquisition command.
   Verify no test, lint, or gate treats aggregate file/LOC totals or percentage variance
   as pass/fail.
7. **Independent re-review:** move both roots again, use a detached source commit,
   inject duplicate authority members, test flat archive rollback, and try to find a
   caller-reconstructed provider field. Reconcile path-scoped `git diff --stat` and
   `git diff --numstat` semantically against the worker's pre-edit declaration; totals
   carry no verdict. Bind fresh ideation attempt 3 only after this corrected design and
   its staff review are committed. The 0.27 pre-release remains separately blocked on
   rqh's real room-only Subspace E2E.

## Documentation change proposal

Implementation applies these concrete semantics:

```diff
--- docs/site/reference/command-reference.md
+++ docs/site/reference/command-reference.md
@@
-For source-checkout or retained development launchers ... It must list `record`, `validate`, `eligibility`, `consume`, and the semantic record flags.
+For source-checkout or retained development launchers, immediately before every gate lifecycle resolve the launcher and run `spacedock gate --help` exactly once. It must list `prepare`, `record`, `validate`, `eligibility`, `consume`, and the prepare flags `--question`, `--artifact`, `--summary`, `--reference`, and `--workflow-dir`, as well as the existing semantic record flags. Missing capability halts before committing selected sources, preparing a room, or mutating state.
@@
-| `spacedock gate record <entity> --briefing PATH/briefing.json` | Bind a complete retained package manifest whose basename is exactly `briefing.json`. |
+| `spacedock gate prepare <entity> --question TEXT --artifact REVIEW.md --summary TEXT [--reference FILE ...]` | Derive and bind a recorder-ready room. Folder `<slug>/index.md` and flat `<slug>.md` entities both use `<slug>/review/<stage>/briefing-N`. Immediately after preparation the room contains exactly `request.json` plus its located canonical Briefing and no copied sources. Selected files are exact local `git-root://<main|state>/<full-commit>/<path>` objects with raw SHA-256 revisions; no containing ref, fetch, or worktree fallback is required. Success prints `room`, `briefing`, `digest`, and `state`. |
+| `spacedock gate record <entity> --briefing PATH` | Bind the exact canonical Briefing file; no basename is canonical. Request-less/advisory Briefings do not require or synthesize a summary. Prepared rooms instead freeze a clean relative locator, id, digest, and exact caller-authored primary Artifact summary. |

--- docs/site/concepts/gates-and-decisions.md
+++ docs/site/concepts/gates-and-decisions.md
@@
-Before the first officer shows a gate, it binds the exact retained Briefing and commits that package.
+Before the first officer shows a gate, it commits newly authored selected sources and calls `gate prepare` with its question, primary Markdown review, exact concise summary, and References. Spacedock authors and binds the two-file recorder-ready room; the first officer commits that entity-owned room. Chat remains the default. A provider override receives only `/subspace:r gate <room>` and derives all request, Briefing, root, authority, materialization, and evidence mechanics from that bound room.

--- skills/fo-gate-lifecycle/SKILL.md
+++ skills/fo-gate-lifecycle/SKILL.md
@@
-**Capability preflight.** Per lifecycle ... Require `record`, `validate`, `eligibility`, `consume`, ...
+**Capability preflight.** Per lifecycle, resolve `${SPACEDOCK_BIN:-spacedock}` and run exactly one fresh `gate --help`. Require `prepare`, `record`, `validate`, `eligibility`, `consume`, `--question`, `--artifact`, `--summary`, `--reference`, `--workflow-dir`, and the existing semantic record flags. Missing capability halts before selected-source commits, preparation, or state mutation.
@@
-**Retain and bind.** Assemble `ROOM/briefing.json` ... then run `gate record ENTITY --briefing BRIEFING`.
+**Prepare and bind.** Select one Markdown gate-review Artifact and any References, author one concise exact primary-Artifact summary, commit every newly authored selection in its owning main/state history, then run `${SPACEDOCK_BIN:-spacedock} gate prepare ENTITY --question QUESTION --artifact REVIEW --summary SUMMARY [--reference FILE ...] --workflow-dir WORKFLOW_DIR`. Require the four output lines and `state=open`; commit the entity's folder/flat room unit before presentation. Supply judgment and paths, never JSON or derived authority.

--- skills/present-gate/SKILL.md
+++ skills/present-gate/SKILL.md
@@
-Probe before side effects ... If the presenter is missing or mismatched ... use chat.
-Pass one prepared room ... Invoke a Subspace override as `/subspace:r gate <gate-room>`.
+For a selected override, wait until `gate prepare` succeeds and `«state.commit»` commits the bound room, then invoke exactly `/subspace:r gate <gate-room>` once using the emitted `room=` value. Do not probe provider availability, version, or capability; do not select a probe-driven chat fallback; and do not pass or reconstruct entity, workflow, Briefing, actor, approver, destination, provider, manifest, or output fields. The rqh fixed entry owns provider capability and failure behavior from the room-only handoff.
```

The normative spec adds the exact folder/flat commit/archive contract, prepare-time
two-file versus later provider-evidence states, closed Git-root grammar, detached
local-object acceptance, arbitrary locator, request-backed exact summary, recursive
duplicate refusal, four recomputed pins, room-only future handoff, and the rqh manifest's
Briefing id/digest binding. The frontmatter reference removes only the basename claim;
the v1 `gates` schema and request version remain unchanged.

## Out of scope

- Provider discovery, capability probing, executable selection, materialization,
  package allocation, terminal transport, provider invocation, cleanup, Result
  production, or a provider success fixture in s4.
- Remote object acquisition, ref retention, reflog policy, copied source payloads,
  source caches, generic URI/root registries, or provider-specific Git access.
- Compatibility request parsing, an absent-locator or `briefing.json` fallback, a
  second request version, caller-authored ids/digests/authority, or migration wrappers.
- `association.json`, caller-selected provider paths, duplicate summary copies,
  Reference summaries, Markdown summary extraction, NLP, or Unicode normalization.
- A permanent two-file room claim after provider work, universal provider-evidence
  retention, crash journals, automatic deletion of unbound crash residue, or a current
  resume protocol.
- Implementing or approving rqh, consuming attempt 2, binding attempt 3 before
  independent review, or changing gate/status frontmatter during this ideation stage.

## Historical design record (cycles 1–9; superseded)

Make gate-room preparation one mechanical operation. The First Officer supplies the
decision question, one concise primary-Artifact summary, and selected files; Spacedock
derives room placement, portable ids, locators, revisions, canonical digests, request
authority, and the open gate attempt.
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

## Superseded proposed command and room contract

The new command is:

```text
spacedock gate prepare ENTITY --question TEXT --artifact GATE-REVIEW.md \
  --summary TEXT [--reference FILE ...] [--workflow-dir DIR]
```

`--artifact` is required exactly once and is the concise gate review. Its filename must
end in `.md` or `.markdown` case-insensitively, and the generated Artifact always carries
`"mediaType": "text/markdown"`. `--summary` is required exactly once and is the First
Officer's concise human-readable summary of that primary Artifact. The value must be a
JSON-string-compatible UTF-8 string that is nonblank after trimming, but Spacedock
stores the exact supplied string: it does not trim, rewrite, summarize, inspect a
Markdown heading, or derive prose from Artifact bytes. `--reference` selects zero or
more existing supporting files in caller order; selecting the same normalized absolute
path twice is an error. References receive no generated `summary`. The First Officer
owns the question, summary, and file choices; it never supplies JSON, an id, digest,
room, attempt, locator, provider path, actor, or approver.

The CLI counts `--summary` occurrences before path normalization or Git work. Repeating
the flag—whether the values match or differ—fails with
`gate prepare accepts --summary exactly once`; an argument containing invalid UTF-8
fails with `--summary must be valid UTF-8`. Both are deterministic handled errors
before lock acquisition, filesystem mutation, or Git process execution.

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

The generated canonical Briefing carries the supplied prose on the primary Artifact
itself, not only inside the Markdown payload:

```json
{
  "type": "Briefing",
  "version": "1",
  "id": "briefing:task:validation:attempt-1:revision-1",
  "question": "Should this candidate advance?",
  "artifacts": [
    {
      "id": "artifact:gate-review-0123456789ab",
      "uri": "git-root://state/0123456789abcdef0123456789abcdef01234567/review/validation/gate-review.md",
      "mediaType": "text/markdown",
      "rev": "sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
      "summary": "  Résumé — validates Git-root presentation exactly.  "
    }
  ],
  "context": [
    {
      "type": "Reference",
      "id": "reference:staff-review-fedcba987654",
      "uri": "git-root://state/fedcba9876543210fedcba9876543210fedcba98/review/validation/staff-review.md",
      "mediaType": "text/markdown",
      "rev": "sha256:fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210"
    }
  ]
}
```

`summary` uses Review v1's existing Artifact extra-field model; s4 does not add a
Briefing version, a generic extension registry, a summary object, or a Reference
variant. For an s4-prepared/request-backed Briefing, the duplicate-safe raw reader
validates the primary `artifacts[0].summary` as a nonblank valid-UTF-8 string. Canonical
Briefing bytes retain it, while identity-only inventory and association projections
remain id/URI/media type/revision and do not duplicate the prose.

That mandatory-summary profile is not retroactive. A request-less
`gate record --briefing PATH` binding and an advisory-round Briefing keep their existing
Review v1 behavior: no primary summary is required, an absent summary is not synthesized,
and their canonical bytes are neither rewritten nor migrated. The request boundary,
not the basename or a global Briefing parser switch, selects the stricter s4 profile.

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

Cycle 8 adds `--summary` but no request field. The schema delta remains one closed
`git-root` Artifact/Reference URI profile plus the required primary-Artifact `summary`
string within Review v1's existing extra-field model; raw-byte SHA remains Review v1's
existing `rev`. The CLI constructs the current `main`/`state` root map from the same
resolved definition and entity roots already used for recorder operations and passes
it to `internal/gitsource`; it never serializes checkout paths. No current provider
handoff consumes that map. This task stops at recorder readiness; sibling
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

## Superseded recorder and JSON authority

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

## Superseded First Officer ownership and provider boundary

This design rebases on the latest 6y tip inspected at `60adfc1f` (the lifecycle work is
in `e9415a17`) rather than current `main`'s transitional presentation-channel prose.
After 6y lands, `skills/fo-gate-lifecycle/SKILL.md` owns preparation, gate capability
preflight, mutation, presentation routing, recording, and consumption.
`skills/present-gate/SKILL.md` is rendering-only and is not changed by this task.
Composed with current main's xb recorder, 6y's pre-xb `--result --association` wording
and capability anchor do not survive: provider-backed recording uses the retained room
surface and derives the association in memory. The lifecycle's Spacedock-only
`gate --help` preflight requires `prepare` and `--room`, not `--association`.

For the no-override path, `fo-gate-lifecycle` writes one concise sentence summarizing
the primary gate-review Artifact, commits any newly authored selected source in its
owning Git history, then passes that exact sentence through `--summary` to
`gate prepare`. It commits the entity folder containing the two-file room and binding,
passes the gate-review Artifact to the rendering-only `present-gate`, and records the
captain's semantic chat decision. The First Officer supplies prose but never edits
Briefing JSON; a Markdown summary section is neither read nor substituted. The source
commit makes each selected worktree byte an immutable reachable object; the later
folder-scoped state commit includes only the generated metadata room and binding
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

## Superseded mechanisms and rejected alternatives

| Mechanism | Value AC | Simplest alternative | Why insufficient |
|---|---:|---|---|
| One `gate prepare` operation | AC-1 | Tell the FO to write two JSON files and call `gate record` | Preserves the manual ids/digests and partial-room failure that caused the task. |
| Required caller-authored primary Artifact `summary` | AC-1, AC-4 | Extract a heading or summary section from the Markdown | A buried section is not canonical Briefing metadata and extraction makes Spacedock invent or reinterpret prose; the provider needs one directly exposed human description. |
| Closed Git-root locator and local-object resolver | AC-1, AC-5 | Copy selected bytes into the room | Copies survive the reopen but duplicate objects the Captain requires to remain singular; path plus raw SHA alone cannot locate bytes after checkout movement. |
| Explicit presentation dependency `rqh46ey33aqq4rt72b4w1m2q` | AC-4 | Say q0 will carry logical roots later | Transport does not turn Git objects into filesystem bytes current Subspace accepts; only a real consumer/materialization API plus end-to-end presentation can close that gap. |
| Frozen local Briefing locator | AC-2, AC-5 | Keep joining `briefing.json` | Fails the reproduced valid room and contradicts the provider contract. |
| One recursive duplicate-member reader | AC-3, AC-5 | Rely on `encoding/json` plus typed structs | Go accepts conflicting duplicates last-wins; the detached counterexample can close under the wrong authority. |
| Stable room/identity stdout handoff | AC-1 | Omit the room or make callers reconstruct it from ids/directory layout | Hides the published artifact and can select the wrong attempt under retries. |
| In-memory derived association | AC-5 | Persist `association.json` | Creates a second durable truth that can diverge from the four frozen inputs. |

## Superseded expected surface and tolerance

Baseline assumption: latest 6y (`60adfc1f`, including lifecycle owner `e9415a17`) lands
first. Relative retained-input normalization is then available in `internal/cli`, the
existing recorded-gate journey targets `fo-gate-lifecycle`, and `present-gate` contains
rendering only. Against that composition, the smallest expected implementation is these
19 files and about `+1,533/-161` lines (**1,694 changed LOC**):

The inspected 6y tip is still pre-xb-rebase, so implementation must not start until
6y's final xb rebase lands. Re-read that landed tip before creating the worktree; if it
changes lifecycle ownership, recorder commands, shared live assertions, exact-authority
capture, the prepared request boundary, or any declared file/delta below, return to
ideation for a surface reset rather than implementing against this provisional
composition.

| File | Expected delta | Purpose |
|---|---:|---|
| `internal/cli/cli.go` | `+68/-6` | Route `gate prepare`, require/pass the primary summary, normalize inputs, and print the stable four-line result. |
| `internal/cli/gate_test.go` | `+215/-25` | Reuse CLI fixtures for summary/source preparation, stdout, committed-source refusal, arbitrary-locator eligibility, and byte-clean failures. |
| `internal/gates/prepare.go` (new) | `+232/-0` | Derivation, exact primary summary, recorder-ready Git-root locators, ids/media types, and error-atomic room publication. |
| `internal/gates/prepare_test.go` (new) | `+235/-0` | Focused exact-summary/no-derived-prose, replay, collision, locator selection, two-file room, and handled-error tests. |
| `internal/gitsource/source.go` (new) | `+175/-0` | Closed root/commit/path URI grammar, root map, ref-retention checks, common-history classification, and local `git cat-file` resolution. |
| `internal/gitsource/source_test.go` (new) | `+170/-0` | Independent moved-checkout, linked-worktree classification, later-worktree, shallow/pruned/missing-object, escaping, and raw-SHA controls. |
| `internal/gates/operation.go` | `+94/-35` | Closed request locator, primary-summary validation without identity duplication, and the one exact Briefing resolver. |
| `internal/gates/application.go` | `+12/-4` | Route reviewed-input eligibility through that resolver instead of `briefing.json`. |
| `internal/gates/io.go` | `+30/-8` | Recompute the four retained provider inputs through duplicate-safe reads. |
| `internal/gates/json.go` (new) | `+75/-0` | Recursive duplicate-member rejection. |
| `internal/gates/testdata/gate-room/request.json` | `+1/-0` | Add the locator to the canonical fixture. |
| `internal/gates/testdata/gate-room/briefing.json` | `+1/-0` | Give the canonical primary Artifact its required human summary. |
| `internal/ensigncycle/recorded_gate_lifecycle_test.go` | `+33/-16` | Add the source commit and caller-authored summary before prepare to the shared no-override observation; add no lane. |
| `internal/contractlint/fo_function_reference_invariant_test.go` | `+22/-8` | Pin lifecycle ownership, Git-only process boundary, forbidden provider mechanics, and rendering-only `present-gate`. |
| `docs/specs/gate-resolution-frontmatter-contract.md` | `+98/-30` | Normative summary/prepare/request/Git-root/retention/resolver/atomicity contract. |
| `docs/site/reference/command-reference.md` | `+25/-6` | New verb, primary summary, committed-source locator, stdout, and arbitrary-name recording. |
| `docs/site/reference/frontmatter-contract.md` | `+3/-3` | Remove the manifest-basename claim. |
| `docs/site/concepts/gates-and-decisions.md` | `+14/-6` | Caller-authored primary summary, committed-source preparation, recorder-ready boundary, and presentation dependency. |
| `skills/fo-gate-lifecycle/SKILL.md` | `+30/-14` | Replace hand bind/association wording with source commit, summary, prepare, recorder-ready recording, and provider halt. |

Cycle 9 changes no expected file or LOC allocation. Repeated-flag/invalid-UTF-8 cases
fit the existing `internal/cli/gate_test.go` budget; request-backed versus request-less/
advisory profile cases fit the existing `internal/gates/operation.go`, gate fixture, and
CLI test allocations; the non-ASCII whitespace sentinel replaces, rather than adds, an
ordinary exact-summary value.

Tolerance is **+2 files and +25% changed LOC** (hard cap 21 files / 2,118 changed
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

## Superseded acceptance criteria

**AC-1 (VALUE) — One command turns judgment and file choices into a validated open
gate room with zero caller-authored metadata files.** Starting baseline: the fixture has
one gate-review Markdown file, selected References, and no room. After `gate prepare`,
the derived room contains exactly the request and located canonical Briefing—**2**
regular files regardless of selected-source count, with **0** duplicated source
payloads. The attempt is open and `gate validate` succeeds. The four stdout lines expose
the exact cleaned absolute room, Briefing id, digest, and open state; the required
Artifact is `text/markdown`, its `summary` is the exact caller-supplied nonblank string,
no Reference gains a generated summary, Reference media types follow the closed table,
and every selected source URI names its `main`/`state` root, full commit, and
repository-relative path while `rev` matches the Git object's raw bytes. *Test:* a real
split-root CLI/Git fixture commits selected main/state files and gives the Markdown a
conflicting `## Summary` section. It invokes prepare with a different sentinel
`--summary '  Résumé — validates Git-root presentation exactly.  '`, asserts
pre-command metadata count 0, post-command count 2, that whitespace-bearing/non-ASCII
sentinel exactly at `artifacts[0].summary`, zero Reference summaries, and exact room/
stdout/locators/media types, then makes later commits that move both worktree paths. It
relocates main and state checkouts independently, supplies their new logical root map,
and revalidates the exact old bytes from both local object databases. Missing or
whitespace-only `--summary` exits nonzero before mutation. Two deterministic CLI cases
also run before any Git call or mutation: repeated same and different values both exit
1 with exactly `gate prepare accepts --summary exactly once\n`; an in-process argument
`string([]byte{0xff})` exits 1 with exactly
`--summary must be valid UTF-8\n`. Removing the cardinality check permits last-value
wins and fails the repeated case; allowing `encoding/json` replacement of invalid bytes
fails the invalid-UTF-8 case. Removing the CLI value and extracting the Markdown section
makes the positive sentinel assertion fail. The test also fails if the fixture supplies
metadata, the room copies either source, a locator omits any identity component, output
changes under a launch directory containing spaces, or resolution needs the old `..`
topology.
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
positive case fail. That request-backed fixture contains the mandatory primary summary.
A paired request-less `gate record --briefing` fixture and existing advisory-round
fixture deliberately omit it; both retain their prior success behavior and exact input
Briefing bytes, with no synthesized summary. Applying the request-backed summary check
globally makes those two positive compatibility controls fail.

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
an actual Spacedock-to-Subspace presentation proof before the 0.27 pre-release; its
consumer must preserve the canonical Briefing unchanged and expose the exact primary
Artifact `summary` rather than copying or deriving it in a resolved-source manifest.
Its E2E uses `  Résumé — validates Git-root presentation exactly.  ` and asserts the
Artifact chrome retains both leading spaces, both trailing spaces, and the exact
non-ASCII code points. No
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
inventory are each deleted and byte-mutated independently. The primary summary is
deleted, changed to whitespace, and byte-mutated independently; every Reference remains
summary-free. Each selected Git locator is also changed one component at a time: root,
full commit, path, and raw
SHA. Every variant fails recording or read-only validation without changing the entity,
while a later worktree edit/move remains green because it does not change the pinned
object. *Test:* real CLI end-to-end fixture asserts the four provider digest pins, each
Git object/raw-SHA pair, and the exact room tree; adding an association input/file,
copying a source payload, omitting a local object, or changing one locator component
fails.

## Superseded test plan and proof order

0. **Baseline gate, before implementation:** require 6y's final xb rebase and
   Cycle-31 authority-capture design to be landed, record its tip, and compare lifecycle
   ownership, recorder commands, exact-authority capture, prepared-request fields,
   shared assertions, and the expected-surface table. Any mismatch returns to ideation
   for reset; a Captain conn is never passed as an ordinary frozen Reference.
1. **Focused red/green, low cost:** add the arbitrary-name spike as the first command
   test using the existing gate-room fixture, then add focused `prepare_test.go` cases
   for exact stdout data, exact primary summary/no Reference summary, repeated
   `--summary` with stable stderr, invalid UTF-8 constructed in-process with stable
   stderr, Git-root URIs/media types/raw digests,
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
   level and nested authority-bearing objects. Delete, blank, type-change, duplicate,
   and byte-change `artifacts[0].summary`; keep a conflicting Markdown `## Summary` and
   prove it is never substituted. Assert entity bytes and lock state, not only error
   substrings. The arbitrary-locator positive case continues through provider room
   closure and CLI eligibility, so `application.go` cannot silently retain its basename
   join. Run request-less bind and advisory-round controls with summary-free Briefings;
   they must keep their existing behavior and exact input bytes.
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
   actual provider-neutral resolved-byte API and Spacedock-to-Subspace E2E, including
   exact Artifact-chrome display of
   `  Résumé — validates Git-root presentation exactly.  `.
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

## Superseded documentation change proposal

The implementation applies these concrete semantics (line wrapping may follow the
target file):

```diff
--- docs/site/reference/command-reference.md
+++ docs/site/reference/command-reference.md
@@
-| `spacedock gate record <entity> --briefing PATH/briefing.json` | Bind a complete retained package manifest whose basename is exactly `briefing.json`. Other basenames fail before mutation. |
+| `spacedock gate prepare <entity> --question TEXT --artifact REVIEW.md --summary TEXT [--reference FILE ...]` | Derive and bind one recorder-ready two-file room. `--summary` is required caller-authored prose stored unchanged as the primary Artifact's Review v1 `summary`; Spacedock does not extract it from Markdown, and References receive none. Selected files must be committed objects in the workflow's main or state Git history and retained by an ordinary ref; the generated Briefing records `git-root://<root>/<full-commit>/<path>` and raw SHA-256 `rev` without copying payloads. The required objects must already be local; no fetch, deepen, hydration, worktree fallback, or retention ref is created. Success prints exactly `room`, `briefing`, `digest`, and `state` key/value lines. Current Subspace cannot present these locators; `rqh46ey33aqq4rt72b4w1m2q` owns that pre-release dependency and must display the canonical summary. |
+| `spacedock gate record <entity> --briefing PATH` | Bind any readable canonical Briefing by its exact path. This request-less/advisory-compatible path does not require or synthesize a primary summary. A prepared request-backed room instead requires the primary Artifact summary and freezes its Briefing locator, id, and digest in `request.json`; every later operation resolves that locator rather than a canonical basename. |

--- docs/site/concepts/gates-and-decisions.md
+++ docs/site/concepts/gates-and-decisions.md
@@
-Before the First Officer shows a gate, it binds the exact retained Briefing and commits that package.
+Before the First Officer shows a no-override gate, it writes one concise primary-Artifact summary, commits newly authored selected sources, and passes the summary and file choices to `gate prepare`; Spacedock writes the canonical JSON and binds the two-file room mechanically. It then commits the entity folder containing that room and renders in chat. Git-root locators reopen exact committed objects through the current main/state root map after checkout movement. This is recorder-ready, not presentation-ready: a selected override halts before invocation and names `git-root-review-v1-materialization`/`rqh46ey33aqq4rt72b4w1m2q`, which owns the resolved-byte consumer, exact canonical-summary display, and actual provider E2E. Spacedock does not fetch objects or discover, probe, launch, or materialize for a provider.

--- skills/fo-gate-lifecycle/SKILL.md
+++ skills/fo-gate-lifecycle/SKILL.md
@@
-**Retain and bind.** Assemble `ROOM/briefing.json` ... then run `gate record ENTITY --briefing BRIEFING`.
+**Prepare and bind.** Select one Markdown gate-review Artifact and any References, and write one concise human-readable sentence summarizing the primary Artifact. Do not substitute a Markdown `Summary` section. Commit every newly authored selection in its owning main/state Git history, then run `${SPACEDOCK_BIN:-spacedock} gate prepare ENTITY --question QUESTION --artifact REVIEW --summary SUMMARY [--reference FILE ...] --workflow-dir WORKFLOW_DIR`. Require the generated primary Artifact's `summary` to equal `SUMMARY`, require References to carry no generated summary, require the four stable output lines and `state=open`, then commit the entity folder containing the generated two-file room and binding before presentation. The FO supplies prose and choices, never JSON. The emitted `room=` value is the sole recorder/diagnostic locator and must never be reconstructed or searched for. Do not copy selected repository objects into the room.
+**Presentation boundary.** Call the result recorder-ready. With no override, render the generated review through `present-gate` and record chat. With a selected override, halt before invocation and name `git-root-review-v1-materialization`/`rqh46ey33aqq4rt72b4w1m2q`; current Subspace cannot consume the Git-root Artifacts/References. This lifecycle does not promise a room/root handoff, name/discover/version-check/capability-probe/launch a provider, materialize sources, or construct or mutate provider output paths.
```

The normative spec makes the same substitutions, scopes the required exact
primary-Artifact `summary`/summary-free Reference behavior to s4-prepared/request-backed
Briefings without a generic extension schema, preserves request-less/advisory Briefings
unchanged, defines the closed request shape and room publication/error-atomicity
behavior, defines the closed
`git-root://<root>/<full-commit>/<repo-path>` URI grammar, ordinary-ref source-commit
rule, and local-only shallow/partial/prune failure behavior, and states explicitly that
association is recomputed and unstored. It labels the output recorder-ready and names
the separate presentation dependency, which preserves the canonical Briefing and
exposes its exact summary without manifest duplication. The frontmatter reference
removes only its exact-basename claim; no `gates` schema or request field changes.
`skills/present-gate/SKILL.md` remains unchanged and rendering-only after 6y.

## Superseded Git-root and authority boundaries

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

## Superseded release dependency and readiness language

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
for that gate. The sibling must pass the canonical Briefing unchanged through its
resolved-source flow, keep `summary` out of the duplicated resolved-source manifest,
and prove Subspace reads Review v1 Artifact `Extra["summary"]` and renders the exact
`  Résumé — validates Git-root presentation exactly.  ` string in Artifact chrome,
including its two leading/trailing spaces and non-ASCII code points.

## Superseded out of scope

- Subspace q0, room-to-provider invocation, Git-root consumer/materialization,
  terminal transport, provider discovery or capability probing, provider output
  allocation/mutation, retained-preflight proof, provider retention implementation, or
  presentation E2E; these presentation concerns are owned by
  `rqh46ey33aqq4rt72b4w1m2q`.
- Compatibility request parsing, a `briefing.json` fallback for prepared requests, or
  migration wrappers.
- `association.json`, caller-selected Result/log/inventory/diagnostic paths, or provider
  argv.
- Summary extraction/NLP, Markdown-heading conventions, summary objects, Reference
  summaries, generic Artifact-extension machinery, or copying the primary summary into
  identity inventory, association, or a resolved-source manifest.
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

- **Cycle 8 — canonical primary-Artifact summary (2026-07-25).** The canonical
  Briefing must expose a concise human-readable summary on its primary gate-review
  Artifact; a summary section buried inside the Markdown payload is insufficient.
  Spacedock must not invent prose and the First Officer must not handcraft JSON.
  Preserve the caller-authored summary through Review v1's existing Artifact
  extra-field model, generate no Reference summaries, and require the downstream
  materialization/presentation dependency to retain and display the canonical value
  without copying it into a parallel manifest.

- **Cycle 9 — approved staff polish (2026-07-25).** Staff returned APPROVE with no
  material findings. Before the superseding gate, make repeated `--summary` and invalid
  UTF-8 refusal deterministic, scope mandatory summary validation to
  s4-prepared/request-backed Briefings so request-less/advisory inputs remain unchanged,
  and use a whitespace-bearing non-ASCII sentinel in the downstream exact-display E2E.
  These are proof/scope clarifications, not an interface or architecture reset.

- Cycle 10: REVISE — Subspace FO + independent staff review; surface state-only review,
  0 product files/0 LOC vs estimate 19 files/1,694 LOC (0%); AC unchanged — bind a
  fresh attempt over current main, define collision-free folder/flat room placement,
  scope the two-file invariant to preparation before provider evidence exists, remove
  the invented ordinary-ref retention policy, align rq on room-only invocation and
  full Briefing-digest verification, and recompute the smallest owner/file/LOC surface.

- Cycle 11: REVISE — independent staff review at `1132eecb`; surface state-only review,
  0 product files/0 LOC vs advisory estimate 26 files/1,904 LOC (0%); AC unchanged —
  remove numerical surface authority, require complete fresh prepare help before
  effects, replace the agent-owned provider probe/fallback with the exact post-bind
  room handoff, and keep manifest duplicate parsing with rq.

- Cycle 12: REJECTED — fresh validation and detached adversarial audit; surface 37 files and 4,013 changed LOC vs estimate 26 files and 1,904 changed LOC (211%); AC unchanged

  Under the delegated sprint conn, the design-reset decision reconfirms the approved
  interface and architecture because the reconciled expansion is support and proof
  coverage, not semantic drift. Correct only the literal flat-archive pathspec outcome
  defect, the masked same-byte/different-path identity detector, and the
  selected-override observed-behavior proof defect; add no protocol, controller,
  compatibility layer, public argument, or narrowed AC.

- Cycle 13: REJECTED — correction revalidation and detached adjacent-trace audit; surface 41 files and 4,555 changed LOC vs estimate 26 files and 1,904 changed LOC (239%); AC unchanged

  Under the delegated sprint conn, the cycle-limit decision reconfirms a final bounded
  evidence-only correction: count every help attempt, reject Agent provider probes and
  chat gate presentation anywhere in the selected-override trajectory, change no
  prompt or product behavior, and add no new harness mechanism.

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

## Stage Report: ideation (cycle 8)

- DONE: Add the smallest honest canonical primary-Artifact summary interface.
  `gate prepare` now requires one `--summary TEXT`; Spacedock rejects missing/blank/invalid UTF-8, stores the exact caller-authored string at `artifacts[0].summary`, never extracts Markdown prose, and generates no Reference summaries or caller-authored JSON.
- DONE: Reconcile Review v1 preservation with the downstream materialization/presentation owner.
  `summary` uses the existing Artifact extra-field model; canonical Briefing bytes retain it, identity projections and the resolved-source manifest do not duplicate it, and `rqh46ey33aqq4rt72b4w1m2q` must prove Subspace renders the exact canonical value in Artifact chrome.
- DONE: Update the command/room example, acceptance criteria, tests, docs, and LOC delta.
  The example shows primary `summary` beside Git-root identity and a summary-free Reference; AC-1/AC-4/AC-5 and the proof order falsify Markdown extraction, missing/blank/type-changed/duplicate summaries, prose duplication, and provider-side derivation. The surface is 19 files at +1,533/-161 (1,694 changed LOC), capped at 21 files/2,118 LOC.
- DONE: Update the retained preview room while keeping provider evidence immutable.
  The historical frozen-source preview's primary Artifact now has one concise summary and `request.json` pins recomputed canonical digest `sha256:5898364167a01112a01ab4fbededf06aad4f2c01b844fa7d703565b060336fd5`; supporting entries gained no summaries.
- SKIPPED: Consume/apply the obsolete ideation approval, dispatch implementation, invoke a provider, relaunch Subspace, materialize Git-root inputs, or mutate provider outputs.
  Cycle 8 changes only s4's state design and retained local preview metadata; the external provider package and provider-owned Result/log/inventory/diagnostics remain untouched.

### Summary

Cycle 8 makes the human description a first-class part of the canonical Briefing
without turning Spacedock into a prose generator or making the First Officer author
JSON. The exact summary now has one source of truth from prepare through the filed
materialization dependency, with a rebaselined recorder-only implementation surface.

## Stage Report: ideation (cycle 9)

- DONE: Make repeated `--summary` and invalid-UTF-8 refusal explicit deterministic tests.
  Same-value and different-value repetition both require exit 1 and exact `gate prepare accepts --summary exactly once\n`; an in-process `string([]byte{0xff})` requires exit 1 and exact `--summary must be valid UTF-8\n`, with zero Git calls and byte-identical state.
- DONE: Scope mandatory primary summary to s4-prepared/request-backed Briefings.
  The request boundary selects the stricter profile; paired request-less bind and advisory-round fixtures remain summary-free, retain prior behavior and exact Briefing bytes, and fail if the new validation is applied globally.
- DONE: Hand the downstream exact-display proof a whitespace-bearing, non-ASCII sentinel.
  The coordinated E2E value is `  Résumé — validates Git-root presentation exactly.  ` and must retain both leading spaces, both trailing spaces, `é`, and the em dash in Subspace Artifact chrome without manifest duplication or normalization.
- DONE: Preserve the approved interface, architecture, docs surface, and LOC plan.
  These cases fit the existing 19-file +1,533/-161 (1,694 changed LOC) allocation and 21-file/2,118-LOC cap; no command, request, provider, or ownership boundary changed.
- SKIPPED: Implement, consume the obsolete approval, invoke a provider, or mutate provider outputs.
  Cycle 9 is a state-only pre-gate clarification and leaves code, gate application, Subspace, and retained provider evidence untouched.

### Summary

Cycle 9 closes the approved reviewer’s three nonblocking proof gaps without reopening
the design. Cardinality/encoding failures are deterministic, compatibility scope is
explicit, and downstream exact display now has a normalization-sensitive sentinel.

## Stage Report: ideation (cycle 10)

- DONE: The corrected design binds a fresh ideation attempt over current main and names one coherent owner map from preparation through the rq/Subspace room-only handoff.
  The governing design is pinned to main `4ff98d8c`, assigns preparation to `fo-gate-lifecycle`, channel routing to `present-gate`, and requires a new attempt-3 Briefing after independent review; attempts 1–2 remain immutable history (AC-1, AC-4).
- DONE: Folder and flat entities receive collision-free rooms; the prepare-time two-metadata-file invariant coexists explicitly with a later provider-owned evidence subtree and zero copied source payloads.
  Both forms converge on `<state-root>/<slug>/review/...`; the flat Markdown plus companion directory form one exact commit/archive unit, while provider evidence is permitted only after the two-file recorder-ready state (AC-1, AC-5).
- DONE: The local-object policy, exact-summary boundary, Briefing-digest binding, files/LOC estimate, and falsifiable tests are rebaselined without invented retention or compatibility obligations.
  Detached local objects are accepted without ordinary-ref/reflog policy; request-backed summary bytes stay exact and downstream-safe, the rqh manifest binds Briefing id plus digest, and the current estimate is 26 files / 1,904 changed LOC with behavior-changing controls for AC-1–AC-5.
- SKIPPED: Implement code, mutate gate frontmatter, bind attempt 3, or invoke a provider.
  Cycle 10 is a state-only ideation correction; the First Officer owns the fresh gate after review, and rqh owns provider materialization/invocation.

### Summary

Cycle 10 replaces the stale pre-rebase hybrid with one current-main, room-authoritative
contract. It makes flat rooms durably committable/archivable, removes the invented ref
policy, narrows “two files” to preparation, and hands rqh one room-only interface with
canonical Briefing identity and digest intact.

## Stage Report: ideation (cycle 11)

- DONE: Remove numerical implementation authority while retaining the named table as an advisory estimate.
  The design has no file/LOC/percentage tolerance or cap. The implementation worker instead declares intended paths and deltas before editing, reconciles path-scoped `git diff --stat` and `git diff --numstat` semantically before review, and no test, lint, worker instruction, or gate may decide by aggregate totals.
- DONE: Correct the current-main capability and provider-handoff owner map.
  One fresh `fo-gate-lifecycle` help preflight requires `prepare` and every prepare flag before source, preparation, or state effects; after preparation and the bind commit, `present-gate` emits exactly `/subspace:r gate <room>` without reconstructing coordinates, probing a provider, or choosing a fallback. The existing command, contract, and lifecycle fixture plans fail on stale help or an agent-run probe (AC-4).
- DONE: Keep recursive duplicate parsing with the document owner.
  s4's shared reader covers only request, located Briefing, Result, and presented inventory; rq's closed manifest parser independently rejects recursive duplicates before display (AC-3, AC-4, AC-5).
- SKIPPED: Implement product changes, bind ideation attempt 3, mutate gate frontmatter, invoke a provider, or edit rq.
  Cycle 11 changes only s4's governing ideation body and leaves the approved room/schema/Git-root/summary/association design intact.

### Summary

Cycle 11 removes numerical planning gates and aligns the current-main skills with the
actual room-only ownership boundary. The approved preparation contract is unchanged;
stale launcher help now fails before effects, and provider capability/failure remains
entirely behind rq's fixed post-bind room handoff.

## Implementation Intended Surface Declaration

Before the first product edit, implementation intends to follow the Cycle 11 advisory
surface exactly: `internal/cli/{cli.go,gate_test.go,state_sync.go,state_commit_test.go}`;
new `internal/gates/{prepare.go,prepare_test.go,json.go}`; new
`internal/gitsource/{source.go,source_test.go}`; existing
`internal/gates/{operation.go,application.go,io.go}` and the two `gate-room` JSON
fixtures; `internal/status/{mutate.go,native_mutation_test.go,merge.go,merge_guard_test.go}`;
`internal/ensigncycle/recorded_gate_lifecycle_test.go`;
`internal/contractlint/fo_function_reference_invariant_test.go`;
`docs/specs/gate-resolution-frontmatter-contract.md`;
`docs/site/reference/{command-reference.md,frontmatter-contract.md}`;
`docs/site/concepts/gates-and-decisions.md`; and
`skills/{fo-gate-lifecycle,present-gate}/SKILL.md`.

The pre-edit estimate is the advisory table's per-path total of approximately
`+1,717/-187`. There are no known path, ownership, command-surface, schema, provider,
validator, or compatibility deviations at declaration time. Actual path and line
deltas will be reconciled semantically before review; counts are planning evidence,
not acceptance authority.

## Implementation Surface Reconciliation

The final implementation diff against `4ff98d8c` contains 37 files at
`+3,683/-330`. Every declared path is present. The larger delta comes chiefly from
complete production Git-source resolution, atomic flat-companion failure proof, and
the shared host-neutral lifecycle oracle rather than a change to the approved command,
schema, provider, authority, or compatibility boundaries.

Eleven supporting paths were added beyond the advisory declaration:
`internal/cli/state.go`; `internal/gates/{gates_test.go,model.go}`; and
`internal/ensigncycle/{claude_live_runner_test.go,codex_live_runner_test.go,gate_assert_impl_test.go,gate_assert_test.go,livescenario_adapter_live_test.go,recorded_gate_lifecycle_pi_live_test.go,shared_fixtures_test.go,shared_scenarios_negative_test.go}`.
They isolate machine-readable Git stdout, reject authority-document duplicate members,
and route every existing host lane through one prepared recorded-gate fixture and
shared oracle. No declared path was omitted, and no provider implementation,
materializer, retained provider output, or sibling workflow entity was changed.

Final-panel rounds found and closed unsafe stage placement, exact live identity,
entity-selected replay, hidden duplicate-locator authority loss, slug-style identity,
and stderr-corrupted atomic flat commits. The final detector correction scopes exact
authority counts to YAML frontmatter while leaving Stage Report proof body-scoped.
The final panel's locator-less legacy-room and historical-object recommendations remain
declined: this is the approved unreleased v1 contract, whose retained authority is
fail-closed and whose prepared request is frozen. The valid resume wording finding was
fixed by distinguishing mutable request-less binds from frozen request-backed prepared
binds; no additional panel cycle was opened.

## Stage Report: implementation

- DONE: One `gate prepare` command produces one durable recorder-ready room for folder and flat entities with zero caller-authored metadata and zero copied selected-source payloads.
  Evidence: `TestPrepareCreatesOneTwoFileRecorderRoomForFolderAndFlatEntities`, slug-style folder/flat coverage, atomic publication rollback, and flat commit/archive tests including injected successful Git stderr. Falsifier: any third prepare-time file, copied source, sibling sweep, partial companion commit, or leaked room after a handled failure.
- DONE: Every selected file is an exact movable local Git object bound by logical root, full commit, repository path, and raw SHA-256, with no fetch, worktree fallback, or ref-retention precondition.
  Evidence: `internal/gitsource` moved-root, linked-worktree, detached-object, literal-path, dirty-source, and missing-object tests. Falsifier: resolution from current worktree bytes, a remote fetch, path/root/commit drift, or raw digest mismatch.
- DONE: The request freezes an arbitrary canonical Briefing locator, id, full digest, and exact request-backed primary summary; request-less/advisory inputs remain unchanged and authority documents reject duplicate members.
  Evidence: arbitrary-locator bind/record/eligibility tests, exact whitespace-bearing summary tests, request-less/advisory controls, recursive depth/duplicate tests, and the hidden-final-locator mutation. Falsifier: basename reconstruction, normalized summary bytes, global summary enforcement, last-member-wins parsing, or byte mutation on refusal.
- DONE: The First Officer performs one complete fresh help preflight before effects, and the only selected-provider handoff is post-bind `/subspace:r gate <room>` with no agent probe, fallback, or reconstructed coordinate.
  Evidence: contractlint and integration smoke tests, the shared deterministic lifecycle, final-product-tip Codex PASS, and two Claude durable-flow completions whose string-scope false positives are now covered structurally. Falsifier: missing/duplicate help, any pre-bind effect, provider probe/fallback, or a handoff coordinate not copied from emitted `room=`.
- DONE: Room-backed recording derives association from recomputed request, Briefing, Result, and inventory pins, writes no association.json, and permits provider evidence only after preparation.
  Evidence: `TestGateRecordConsumesDirectBindingResultFromPreparedRoom`, retained-authority mutation matrices, inventory checks, and the two-file preparation oracle. Falsifier: any pin drift accepted, unknown authority used, `association.json` written, or provider evidence present at prepare time.
- FAILED: Run the final-product-tip Pi live recorded-gate lane.
  The repository scenario did not start: the installed global `pi-subagents` extension requires unavailable `@earendil-works/pi-coding-agent`. Earlier local auth was also externally stale; no product command or repository assertion failed.

### Summary

Implementation is complete at `f99af3df`. Full tests, race tests, formatting,
`git diff --check`, strict documentation build, focused skill smoke tests, Codex live,
and durable Claude lifecycle evidence are complete; Pi remains an external runtime
installation/authentication blocker rather than an s4 behavior failure.

## Stage Report: validation

- FAILED: Independently verify AC-1 folder and flat gate prepare, exact two-file preparation, replay, atomic state commit/archive, sibling isolation, rollback, and readable archived room.
  Ordinary-form suites passed, but detached `TestValidationFlatArchiveTreatsValidSlugAsLiteralPathspec` failed: `020-[x]` finalize committed `020-x.md` and `_archive/020-x.md` and omitted the literal live deletion; material AC-1 outcome defect in nonliteral archive pathspecs.
- FAILED: Independently verify AC-2 exact git-root main/state identities across moved roots, detached local objects, literal paths, dirty/missing/mismatch failures, and absence of fetch/ref/worktree fallback.
  Positive/moved/detached/literal/failure suites passed, but deleting repository-path equality from `SameLogicalRevision` still passed `TestSameLogicalRevisionIgnoresUnrelatedCommitButNotPathOrBytes`; material AC-2 evidence defect because raw-byte inequality masks the promised path detector.
- DONE: Independently verify AC-3 arbitrary Briefing locator/id/digest, exact summary bytes, request-less/advisory controls, recursive duplicate-member rejection, and byte-clean failure behavior.
  Arbitrary-locator close/eligibility, compatibility controls, duplicate-depth, invalid/cardinality, and cleanup tests passed; trimming the whitespace/non-ASCII summary made `TestPrepareCreatesOneTwoFileRecorderRoomForFolderAndFlatEntities` fail.
- FAILED: Independently verify AC-4 one fresh help preflight before effects and the sole selected-provider handoff /subspace:r gate <room>, with no provider probe, fallback, or caller-reconstructed authority.
  Fresh-help and no-override lifecycle tests passed, but adding `subspace --version` plus chat fallback to `present-gate` still passed both contract tests; selected override has only instruction-file substring proof, a material AC-4 evidence/mechanism defect requiring an observed-behavior harness reset.
- DONE: Independently verify AC-5 association derivation from recomputed request/Briefing/Result/inventory pins, no association.json, and provider evidence only after preparation.
  Prepared-room record/validate mutation matrices passed and fail on missing/drifted request, Briefing, Result, inventory, locator, rev, summary, coverage, persisted association, or prepare-time provider evidence.
- DONE: Perform the required semantic adversarial matrix and detached audit for changed high-stakes CLI/skill surfaces; reject tests that pass while observable authority, identity, order, bytes, or cleanup is wrong.
  Detached audit caught the AC-1 sibling sweep and the AC-2/AC-4 surviving mutants; CLI four-line order and exact-summary normalization mutants were correctly rejected.
- DONE: Reconcile implementation against the latest captain rulings and approved design; report any compatibility preservation, inferred precondition, extra public argument, or provider dependency as drift.
  All declared paths plus eleven explained support paths stay within approved owners; no extra public argument, copied source, stored association, compatibility fallback, or s4 provider dependency was found.
- DONE: Run applicable focused tests, go test ./..., go test ./... -race, gofmt check, git diff --check, strict docs build, skill smoke/contract checks, and a cheapest-first Pi infrastructure probe; report exact commands and results.
  PASS: focused AC regex suite; `go test ./...`; `go test ./... -race`; `test -z "$(gofmt -l ./cmd ./internal)"`; `git diff --check 4ff98d8c..HEAD`; `uv run --with-requirements docs/requirements.txt mkdocs build --strict`; `go test ./skills/integration ./internal/contractlint -count=1`; Pi coverage guard. Pi live failed pre-product because installed `pi-subagents` lacks `@earendil-works/pi-coding-agent`.
- FAILED: Report 8 done, 0 skipped, 0 failed only if all promised ACs have valid evidence and no material finding remains; separately classify deferred risks/polish and recommend PASSED or REJECTED.
  Substantive result is 5 done, 0 skipped, 3 failed; three material findings remain, no deferred-risk or polish finding blocks, Pi is external infrastructure, and the recommendation is REJECTED.

### Summary

Validation reproduced the ordinary-path implementation evidence but found one material
product defect and two material proof defects under detached mutation. Recommend
REJECTED; fix AC-1 narrowly, strengthen AC-2 with same-byte distinct paths, and reset
AC-4 selected-override proof to observed post-bind behavior before re-entry.

## Implementation Correction Surface Reconciliation

The correction advances implementation from `f99af3df` to `6b18f6de`. Against
baseline `4ff98d8c`, the final surface is 41 files at `+4,213/-342`; relative to the
prior 37-file implementation, the only new paths are the AC-1 literal-path regression
and three AC-4 selected-override harness files. AC-2 strengthens its existing test,
and the final-review correction changes only existing lifecycle oracle files.

No command, schema, provider, materializer, source locator, public argument, or
compatibility boundary changed. The product correction is the literal Git pathspec
used to commit archive moves; all other correction work strengthens independent
evidence for already-approved behavior.

## Stage Report: implementation (cycle 2)

- DONE: Triage and fix the material AC-1 literal flat-archive pathspec defect with a real Git regression covering wildcard-like slug plus sibling isolation and deletion.
  Evidence: `TestMergeGuardFlatArchiveTreatsSlugAsLiteralPathspec` uses `020-[x]`,
  proves the literal live deletion and archive destination in HEAD, and proves both
  matching live/archive siblings retain committed bytes while remaining dirty.
  Falsifier: removing `:(literal)` from archive move pathspecs commits a matching
  sibling and omits the literal source deletion.
- DONE: Strengthen AC-2 proof with same raw bytes at distinct repository paths so removing path equality fails independently.
  Evidence: `TestSameLogicalRevisionIgnoresUnrelatedCommitButNotPathOrBytes` inspects
  `review.md` and `same-bytes.md` with identical `review\n` bytes and rejects them as
  different logical revisions.
  Falsifier: the path-equality mutant compiles but fails at the new same-byte,
  distinct-path assertion rather than at a raw-byte inequality.
- DONE: Replace AC-4 instruction-text proof with observed behavior for exactly one post-bind room-only selected override and absence of probe, fallback, or reconstructed authority, using the smallest existing runtime fixture.
  Evidence: the real Claude run passed in 119.93 seconds with one exact root
  `Skill(skill="subspace:r", args="gate <emitted-room>")` after prepare and bind,
  no Bash/provider probe, Agent detour, chat fallback, record, or consume.
  The native Skill-call sink owns no provider behavior; state HEAD advanced, HEAD
  contains the entity plus both room files byte-exactly, and the entity unit is clean.
  Falsifier: unit mutants reject pre-bind/duplicate/reconstructed calls, probes,
  fallback, post-handoff close, missing help, failed prepare, and no-op durability.
- DONE: Run focused tests, full tests, race tests, formatting/diff/docs/skill checks affected by the correction, obtain final Roborev review, update the durable report, and stop.
  Evidence: focused AC-1/AC-2/AC-4 and contract tests, `go test ./...`,
  `go test ./... -race`, `gofmt -w ./cmd ./internal`, `git diff --check`, and strict
  MkDocs all pass at `6b18f6de`; Roborev branch-final job 2440 completed.
  Falsifier: either final material finding originally survived—dynamic-bound identity
  and committed-binding durability—but both now have real-CLI/live regression proof.

### Final review triage

Roborev's High dynamic-digest finding was material and fixed: the hold oracle now
reads the durable Briefing identity/digest, requires the review to name both exactly,
and a real prepare proves a non-legacy digest against canonical room bytes. Its Medium
durability finding was material and fixed with Git HEAD, committed-byte, and scoped
cleanliness checks. Low performance, historical-status diagnostics, recovery prose,
identity-cardinality, usage-polish, unreachable-guard, and test-naming notes are
declined as correct-but-disproportionate for this validation correction; promote when
a representative latency run is unacceptable, historical drift is observed, an
operator recovery is required, or a Low note masks wrong observable behavior.

### Summary

Cycle 2 closes all three rejected acceptance-criterion findings and both material
final-review findings without expanding the approved product boundary. The corrected
implementation and its provider-neutral observed-behavior evidence are ready for
validation re-entry.

## Stage Report: validation (cycle 2)

- DONE: Reproduce the AC-1 valid flat-slug literal archive regression and prove exact source deletion/archive destination plus sibling isolation in Git HEAD and dirty worktree.
  `TestMergeGuardFlatArchiveTreatsSlugAsLiteralPathspec` passed for `020-[x]`; removing `:(literal)` made it fail after committing dirty sibling bytes, while the corrected run deletes the live source, commits the archive destination, and leaves both matching siblings dirty at their prior HEAD bytes.
- DONE: Reproduce the AC-2 same-byte/different-path mutant and confirm repository-path identity is independently detected.
  The equal-byte `review.md`/`same-bytes.md` control passed at `6b18f6de`; deleting path equality made `TestSameLogicalRevisionIgnoresUnrelatedCommitButNotPathOrBytes` fail with `same=true`, independently of raw revision.
- FAILED: Reproduce AC-4 observed selected-override behavior: exact dynamic Briefing identity, committed entity plus two room files, one post-bind room-only Skill call, and no probe/fallback/reconstructed authority/provider behavior.
  Dynamic prepare/state tests and `TestLiveClaudeSelectedGateOverride` passed in 113.03s, but detached adjacent traces with an extra failed help call, a pre-prepare Agent provider probe, or chat presentation before the override all graded PASS; material AC-4 evidence defect at the oracle boundary.
- DONE: Recheck AC-3 and AC-5 plus applicable focused/full/race/format/diff/docs/skill gates, reconcile correction scope to the approved unreleased-v1 design, and report PASSED or REJECTED with material/deferred/polish classification.
  PASS: focused AC-3/AC-5 suites; `go test ./...`; `go test ./... -race`; `test -z "$(gofmt -l ./cmd ./internal)"`; `git diff --check f99af3df..HEAD`; `uv run --with-requirements docs/requirements.txt mkdocs build --strict`; and `go test ./skills/integration ./internal/contractlint -count=1`. Scope remains the approved v1; one material evidence defect remains, no deferred risk was found, earlier Low notes remain nonblocking polish, and recommendation is REJECTED.

### Summary

Correction revalidation closes the prior AC-1 outcome defect and AC-2 detector hole,
and its positive selected-override live run is clean. The gate remains REJECTED because
AC-4's new oracle still accepts three promised wrong traces; the fix stays inside the
existing observation harness.

## Stage Report: implementation (cycle 3)

- DONE: Make the selected-override oracle reject extra failed or successful gate-help attempts while preserving exactly one successful fresh preflight.
  `countCommands` now requires one total help exit and one successful help exit; the
  failed-help-before-success and duplicate-success mutants both fail the oracle.
- DONE: Make the oracle reject Agent provider probes anywhere in the selected-override trajectory, including before prepare.
  Root Agent actions now fail independent of `prepareAt`; the adjacent pre-prepare
  Agent mutant fails, as does the existing post-bind Agent mutant.
- DONE: Make the oracle reject a complete chat gate review anywhere in a selected-override trajectory, including before the Skill handoff.
  Semantic root chat now fails independent of `overrideAt`; the adjacent pre-handoff
  review mutant fails, as does the existing post-handoff fallback mutant.
- DONE: Add the three adjacent-event mutants, run focused/full/race/format/diff checks and final Roborev, update the durable correction report, and stop without product or prompt changes.
  Commit `2d7ee074` changes only the observer and its unit matrix (`+29/-9`);
  focused tests, `go test ./...`, race, `gofmt`, and `git diff --check` pass, and
  final Roborev job 2443 completed.

### Final review triage

The final panel reported no defect in the two-file correction. Its broad branch
findings do not trigger another pass under the Captain-bound Cycle 13 evidence-only
scope.

The legacy request-profile finding is correct-but-disproportionate: no released v1
user owns the prototype state, the observable harm is an upgrade stall, compatibility
would reopen the explicitly excluded v1 boundary, and the trigger requires pre-v1
persisted request bytes. Promote when the Captain adds prototype migration support.

The crafted companion-symlink finding is a deferred risk: an operator selecting a
crafted repository could redirect writes outside state, which would violate path
containment, but the trigger is adversarial and outside the supported flow. Promote
when a released user reaches that layout through an operator-selected repository.

Historical-object revalidation repeats the approved fail-closed local-object policy;
promote when ordinary merge/prune behavior strands a released workflow. The candidate
cleanup-test note is polish because
`TestPrepareRollsBackPublishedRoomAfterBindingWriteFailure` exercises post-publication
rollback. The remaining status, malformed-request, subprocess, diagnostic, and weak
assertion notes are pre-existing deferred risks or polish, not defects in Cycle 13.

### Summary

Cycle 3 closes the three adjacent-event oracle gaps without changing product behavior,
prompt text, or the live harness. The complete selected-override trajectory now enforces
one fresh help attempt and excludes every Agent probe and semantic chat presentation.
