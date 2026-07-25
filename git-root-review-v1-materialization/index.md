---
title: Materialize Git-root Review v1 sources for provider presentation
status: ideation
source: "s4 cycle-6 staff rejection: recorder-valid git-root:// sources are not renderable by current Subspace package mode, 2026-07-25"
started:
completed:
verdict:
score: "1.0"
worktree:
issue:
sprint: durable-decisions
id: rqh46ey33aqq4rt72b4w1m2q
gates:
    version: 1
    current:
        gate: gate:docs-dev:rqh4:backlog
    records:
        - id: gate:docs-dev:rqh4:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:rqh4-backlog-1
              briefing:
                id: briefing:docs-dev:rqh4:backlog:attempt-1:revision-1
                digest: sha256:d620934ee0af1b72c38e80fdb640f6ea07bd95da9fd08729c38e9b9d04a4fce2
                digest-domain: canonical-bytes
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:rqh4:backlog:1
                briefing: briefing:docs-dev:rqh4:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-25T06:28:30.836574Z"
                decision: approve
                reason: The task isolates the actual missing cross-repository consumer boundary, forbids durable source duplication, and requires a real moved-root Subspace proof before implementation.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge. Use your judgement. make sure the implementation declares their intended change (loc etc) and be suspicious about any drift.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

Bridge recorder-ready Git-addressed Briefings to actual provider presentation without
duplicating selected sources in durable gate state.

## Problem and exact revisions

s4's narrowed source identity is
`git-root://<main|state>/<full-commit>/<repository-relative-path>` plus Review v1's
raw SHA-256 `rev`. Spacedock maps `main` and `state` from the currently opened
workflow, reads the addressed local Git objects, and verifies their raw bytes after
the checkouts move independently. That is recorder resolution, not byte delivery.

This design was prepared against:

- Spacedock `cc51e518a3420b01fd4b455e9710d38803dc6d3e`, with the superseding
  summary-corrected s4 design at state commit
  `9a7b7e2fe1042257985562aa3270e811e57a3b2c`;
- released Subspace `eeea57550c082757314c8a22e20c19c0af7bd9db`
  (`0.10.0-beta.4`);
- adjacent Subspace `6280ec446b0c6eff75dd0b2486f25c5e68c2f770`.

The Subspace loader and Reference-selector boundary are byte-identical between beta.4
and the adjacent tip. The later `v1_tui.go` changes are reader-navigation work and do
not change `reviewv1.Load`, `Log.targetText`, provider inventory, CLI arguments, or
capability preflight.

## Riskiest-first room-to-provider spike

The throwaway spike used two real committed sources:

```text
git-root://main/cc51e518a3420b01fd4b455e9710d38803dc6d3e/docs/specs/gate-resolution-frontmatter-contract.md
rev sha256:6425a361d95b0f5fc64d712aa9da8bd1b377ca9579ec6da5c6e3523fc8039003

git-root://state/33b7b48fa0ba9c509a93bb70238d322a734de9f3/prepare-provider-neutral-gate-room/index.md
rev sha256:fcd90ea063597c2f49312538e267364f699461bc49c06dab603a4d8435df4fce
```

The first was the Briefing's Artifact and the second a Reference. Current Spacedock
recorded and validated the exact Briefing bytes:

```text
record_exit=0
validate_exit=0
state=open
briefing=briefing:spike:git-root-materialization
```

Both released beta.4 and the adjacent Subspace tip then passed
`--supports review-v1-provider-package-v1` with exit 0 and received that same
Spacedock-accepted Briefing through real package mode. Both exited 2 before opening
the TUI:

```text
artifact "artifact:gate-contract": resolve "git-root://main/cc51e518a3420b01fd4b455e9710d38803dc6d3e/docs/specs/gate-resolution-frontmatter-contract.md":
open .../briefing-1/git-root:/main/cc51e518a3420b01fd4b455e9710d38803dc6d3e/docs/specs/gate-resolution-frontmatter-contract.md:
no such file or directory
```

The provider had already allocated only an empty mode-0600 `review.jsonl`; it emitted
no Result or presented inventory. Replacing only the Artifact URI with a temporary
relative file containing the verified Git blob advanced through canonical loading and
failed later only because the noninteractive spike had no `/dev/tty`. A focused test
against beta.4's exact `Log.targetText` then preserved the second boundary:

```text
suggestion reference "reference:s4-design" is not locally resolvable
```

That test passes only because the Reference URI contains `://`. The spike therefore
separates all three facts: Spacedock accepts the canonical identities; local Git can
recover their bytes; current Subspace receives neither Artifact nor Reference bytes.
q0 transport has no implicit mechanism between those facts.

The spike used a candidate two-file room because `gate prepare` is still s4 ideation,
not a command at `cc51e518`. Its Briefing is production-recorder accepted and its
request uses s4's closed shape, but no claim is made that current main generated it.
Implementation's real E2E must replace this candidate with s4's landed command output.

## Decision

Select a provider-neutral resolved-source manifest plus provider-owned ephemeral
payloads. Spacedock remains the sole logical-root/Git resolver. Subspace gains one
opt-in loader for already-resolved sources but no knowledge of workflow layout,
checkout discovery, Git commands, or remote transport.

The alternatives are smaller only in one dimension:

- **Native Git-root resolution in Subspace is rejected.** It would require root-map
  transport anyway, duplicate Spacedock's Git/object/retention rules in each provider,
  and make a Review v1 presenter understand Spacedock split-root topology.
- **A byte stream or inherited-file-descriptor bundle is deferred.** It avoids named
  temporary files but adds framing, backpressure, terminal-stdin separation, descriptor
  inheritance across every host, and crash/retry semantics. No current provider
  boundary consumes it.
- **A rewritten compatibility Briefing is rejected.** Replacing canonical URIs with
  temporary paths changes the review identity and makes Result/inventory association
  depend on a second Briefing. No wrapper, `association.json`, or durable source copy
  is introduced.
- **q0 alone is insufficient.** q0 may eventually transport the room and this manifest
  to one provider invocation; it does not map roots, read objects, verify bytes, or
  define cleanup.

The manifest is the smallest mechanism that fits the provider's existing filesystem
package boundary while keeping the canonical Briefing unchanged.

## Cross-repository contract

After s4 lands, the provider-facing Spacedock command is:

```text
spacedock gate materialize ENTITY --room ROOM --destination DIR \
  [--workflow-dir DIR]
```

`ENTITY` and `--workflow-dir` give Spacedock the same resolved definition/entity roots
used by recorder operations. `--room` must be the exact bound s4 room. `--destination`
must be an existing, empty, clean absolute, non-symlink directory with mode 0700. The
provider integration creates and owns it; Spacedock never chooses a provider cache or
terminal.

Spacedock performs these actions before publishing a manifest:

1. validate the bound request, canonical Briefing locator/id/digest, and current
   main/state logical-root map;
2. walk Artifacts in order and recursively reached References in canonical context
   order;
3. for each `git-root` source, validate the closed URI, invoke the production
   `internal/gitsource` local-object reader, and verify the full raw SHA-256 `rev`;
4. write mode-0600 payloads under a private staging directory, publish the manifest
   last, and remove the whole candidate on any failure.

Success prints exactly:

```text
manifest=/clean/absolute/provider-owned/dir/resolved-sources.json
sources=<decimal-count>
```

`resolved-sources.json` has this closed v1 shape:

```json
{
  "type": "spacedock-resolved-sources",
  "version": "1",
  "briefing": "briefing:docs-dev:s4:ideation:attempt-3:revision-1",
  "items": [
    {
      "type": "Artifact",
      "id": "artifact:gate-review",
      "uri": "git-root://main/<full-commit>/path/to/gate-review.md",
      "mediaType": "text/markdown",
      "rev": "sha256:<64-lowercase-hex>",
      "path": "payload/0001"
    },
    {
      "type": "Reference",
      "id": "reference:staff-review",
      "uri": "git-root://state/<full-commit>/path/to/staff-review.md",
      "mediaType": "text/markdown",
      "rev": "sha256:<64-lowercase-hex>",
      "path": "payload/0002"
    }
  ]
}
```

Items cover every Git-root Artifact and recursively reached Git-root Reference exactly
once. Item order is canonical presentation order. `path` is clean, manifest-relative,
contained, and names a non-symlink regular file; it is deliberately unrelated to the
repository path. An exact canonical tuple is
`type/id/uri/mediaType/rev`. Duplicate, missing, extra, reordered, unknown-field, or
tuple-mismatched items fail closed.

Room-relative sources retain Review v1's existing behavior and do not appear in the
resolved manifest. s4-generated selected repository sources are all Git-root sources,
so its two-source room produces two manifest items.

Subspace advertises a separate literal capability:

```text
subspace-tui --supports review-v1-resolved-sources-v1
```

The private package entry becomes:

```text
subspace-tui --review-v1 --actor ACTOR --approver APPROVER \
  --provider-package PROVIDER_ROOT \
  --resolved-sources /absolute/path/resolved-sources.json \
  CANONICAL_BRIEFING
```

`review-v1-provider-package-v1` remains valid for ordinary filesystem Briefings.
`review-v1-resolved-sources-v1` opts into this additive loader; no version comparison
or fallback path is permitted. For Subspace, the manifest path must equal
`<provider-package>/resolved-sources/resolved-sources.json`; the
`resolved-sources` child must be a mode-0700, non-symlink directory. This makes cleanup
an exact provider-owned child deletion rather than permission to remove an arbitrary
manifest parent.

Subspace reads the canonical Briefing unchanged, validates the manifest's exact
coverage against it, reads each contained payload, and recomputes every raw SHA-256.
Artifacts enter the existing in-memory Artifact bytes. Resolved References enter the
in-memory target/source catalog, render by their canonical `mediaType`, and supply
selector text without rereading the URI. The ordinary loader remains unchanged when
the private flag is absent.

The presented inventory is still derived only from the canonical Briefing, never from
temporary paths or manifest order. It therefore emits every original
Briefing/source id, Git-root URI, media type, and raw revision exactly once. Spacedock's
recorder still derives the association from request + Briefing + Result + inventory
and stores no `association.json`.

## Artifact summary clarification

s4's exact producer surface is one required caller-authored `--summary TEXT` value,
accepted exactly once. It requires an exact UTF-8 string nonblank after trimming but
stores the supplied string unchanged at `artifacts[0].summary`. The resolved-source
manifest does not copy, transform, or generate it. Subspace reads
`Artifact.Extra["summary"]` from the canonical Briefing and makes the exact value visible
in focused Artifact presentation chrome without injecting it into source bytes or
selector coordinates.

References have no summary field in this design. Neither Spacedock nor Subspace invokes
NLP, derives a summary from bytes, or promotes Reference labels into summaries. The
summary is omitted from presented inventory because it is presentation content, not
source identity; the retained canonical Briefing remains its authority.

This clarification adds no rqh Spacedock file because s4 owns generation. In Subspace
it adds approximately `+35/-0` LOC inside the already-planned source/chrome files;
`model.go` is already touched for in-memory Reference bytes. The coordinated s4 surface
is now **19 files, +1,533/-161 = 1,694 changed LOC**, capped at **21 files / 2,118
changed LOC**. rqh does not invent a second summary input.

## Ephemeral ownership and cleanup

The invoking provider integration owns `--destination` from creation through cleanup.
For Subspace it is exactly `<provider-package>/resolved-sources`. Spacedock may write
only inside it and removes partial candidates on every handled materialization
failure. Once Subspace validates that exact contained path, it loads and verifies all
bytes into memory, then removes only the known `resolved-sources` child before opening
the TUI. The parent provider supervisor also installs an idempotent cleanup trap for
normal exit, loader failure, signal, or child crash.

The resolver directory is never included in Result, inventory, diagnostics, or
delivered output. Provider-package resume removes the exact well-known child before
re-materializing. A hard power loss can therefore leave only explicitly recoverable
provider-private residue, never a copied payload in the retained Spacedock room or its
Git history.

Completion, open review, provider validation failure, source-loader failure, signal,
and resume tests each assert that `resolved-sources/` is absent. The retained room
continues to contain exactly request + canonical Briefing.

## Missing, shallow, partial, and pruned objects

Only Spacedock maps logical roots and reads Git:

- missing/unknown root, commit, path, object, or digest mismatch fails before manifest
  publication and before Subspace launch;
- a shallow or partial clone succeeds only when the addressed commit and blob are
  already present locally;
- no fetch, deepen, lazy-object hydration, neighboring-directory search, current
  worktree fallback, or remote URL is allowed;
- if refs are rewritten and Git prunes the object before materialization or a later
  resume, materialization fails closed even when `rev` remains known;
- pruning after manifest publication does not change the current presentation because
  the verified bytes are already provider-owned; the next resume must re-resolve and
  can fail.

This lifecycle makes the local-object retention precondition observable rather than
pretending the raw SHA can recover absent bytes.

## Expected surface and tolerance

The rqh implementation begins only after s4 lands its exact command, URI, summary, and
root-map contract. Planned Spacedock surface:

| File | Expected delta | Purpose |
|---|---:|---|
| `internal/cli/cli.go` | `+55/-5` | Route and document `gate materialize`. |
| `internal/cli/gate_test.go` | `+100/-0` | Stable output, argument, and byte-clean failure cases. |
| `internal/gates/materialize.go` (new) | `+135/-0` | Resolve bound room and publish the closed manifest. |
| `internal/gates/materialize_test.go` (new) | `+260/-0` | Exact coverage, atomicity, and room immutability. |
| `internal/gitsource/materialize.go` (new) | `+190/-0` | Private payload staging around s4's object reader. |
| `internal/gitsource/materialize_test.go` (new) | `+280/-0` | Moved roots, missing/pruned object, raw-SHA failures. |
| `scripts/tests/git-root-subspace-e2e.sh` (new) | `+180/-0` | Real two-repository binary/TUI composition. |
| `docs/specs/gate-resolution-frontmatter-contract.md` | `+70/-5` | Normative resolved-source/lifecycle contract. |
| `docs/site/reference/command-reference.md` | `+20/-2` | Provider-facing command and failure semantics. |

Spacedock baseline: **9 files, +1,290/-12 = 1,302 changed LOC**. Tolerance is +2 files
and +25% changed LOC, hard cap **11 files / 1,628 changed LOC**, only for a fixture or
focused test split. A new request field, provider executable, remote fetch, durable
cache, or generic URI resolver resets ideation.

Planned Subspace surface:

| File | Expected delta | Purpose |
|---|---:|---|
| `internal/reviewv1/model.go` | `+12/-4` | In-memory resolved Reference bytes; Artifact summary stays in `Extra`. |
| `internal/reviewv1/loader.go` | `+35/-6` | Share canonical validation with resolved input. |
| `internal/reviewv1/resolved_sources.go` (new) | `+205/-0` | Closed manifest, containment, coverage, and digest checks. |
| `internal/reviewv1/resolved_sources_test.go` (new) | `+300/-0` | Tuple, payload, nesting, digest, and cleanup cases. |
| `internal/reviewv1/log.go` | `+20/-12` | Selector text uses verified in-memory Reference bytes. |
| `cmd/subspace-tui/main.go` | `+30/-8` | Private flag and literal capability. |
| `cmd/subspace-tui/profile_dispatch_test.go` | `+65/-0` | Capability and exact argv surface. |
| `cmd/subspace-tui/v1_tui.go` | `+25/-5` | Select resolved loader and cleanup before TUI. |
| `cmd/subspace-tui/v1_sources.go` | `+45/-8` | Render resolved Reference bytes and Artifact summary. |
| `cmd/subspace-tui/v1_sources_test.go` | `+90/-0` | Complete source catalog and no synthesized summary. |
| `cmd/subspace-tui/v1_review_chrome_labels_test.go` | `+55/-0` | Exact visible summary without selector drift. |
| `cmd/subspace-tui/SPEC.md` | `+35/-5` | Private package/materialization lifecycle. |
| `docs/review-and-gate.md` | `+25/-4` | Resolved-source profile and Artifact summary semantics. |

Subspace baseline: **13 files, +942/-52 = 994 changed LOC**. Tolerance is +2 files
and +25% changed LOC, hard cap **15 files / 1,243 changed LOC**, only for provider
cleanup or one E2E fixture split. A Git dependency/command, root map, terminal wrapper,
second canonical Briefing, Reference summary, or retained source cache resets ideation.

## Acceptance criteria

**AC-1 (VALUE) — One real Subspace presentation renders 2/2 canonical Git-root
sources after independent checkout movement.** An s4-generated room retains exactly
two regular files and zero selected-source payloads. After main/state move to unrelated
depths and later worktree commits move both current paths, materialization supplies the
old Artifact and Reference bytes; the actual Subspace package capability presents both,
shows the exact primary Artifact summary, and emits a two-item canonical inventory.
*Test:* the cross-repository E2E drives the production binaries and one TUI process,
asserts visible sentinel bytes from both old objects and the exact summary, then checks
Result/inventory ids, URIs, media types, revs, room count=2, and payload count=0.

**AC-2 — Durable gate state remains free of copied selected-source payloads.** The
room and its Git commit contain request + canonical Briefing only. Normal completion,
open review, validation failure, loader failure, signal, child crash, and resume leave
no resolved payload under the retained provider package. *Test:* lifecycle table hashes
the room tree before/after and asserts `resolved-sources/` absent after each outcome;
removing any cleanup edge makes its case retain a sentinel payload.

**AC-3 — Provider association remains canonical despite ephemeral paths.** Every
Git-root Artifact and recursively reached Reference maps exactly once by
`type/id/uri/mediaType/rev`; temporary `path` never enters Result or presented
inventory. The Artifact summary comes only from the canonical Briefing and References
have none. *Test:* adversarial manifest cases duplicate, omit, reorder, add, and mutate
each tuple/path/payload; all fail before TUI. The positive inventory equals the
Briefing-derived tuple list byte-for-byte and contains no temporary path or summary.

**AC-4 — Missing local retention fails before presentation and never fetches.** Unknown
root, absent commit/path/blob, shallow/partial object absence, and raw-SHA mismatch
leave destination empty, room bytes unchanged, provider launch count 0, and no binding
Result. *Test:* the moved-root Git fixture prunes the last ref/object for one case and
uses a missing-object shallow clone for another; a fake `git` records argv and rejects
any `fetch`, `clone`, `pull`, or network helper.

**AC-5 — The cross-repository contract is exercised through production boundaries.**
The committed E2E uses landed `spacedock gate prepare`, production
`gate materialize`, both literal Subspace capabilities, canonical package mode, the
resolved Reference selector path, and Spacedock `gate record --room`/`gate validate`.
*Test:* replacing materialization with `git cat-file` in the harness, bypassing package
mode, rewriting the Briefing, or substituting a provider fixture causes an explicit
precondition assertion to fail.

## Minimum real E2E and proof order

1. **First red:** preserve the ideation spike as a fixture assertion: the exact
   Spacedock-accepted Git-root Briefing makes beta.4/current package mode exit 2 at
   `room/git-root:/...`, while local `git cat-file` succeeds. This fails if the test
   stops at local resolution.
2. **Focused Spacedock:** red/green the closed manifest, destination containment,
   exact recursive order, atomic manifest-last publication, moved-root reads,
   shallow/pruned absence, and raw-SHA mismatch. Run
   `go test ./internal/gitsource ./internal/gates ./internal/cli -count=1`.
3. **Focused Subspace:** red/green the private flag/capability, exact manifest coverage,
   payload containment, Artifact + Reference byte loading, selector behavior, exact
   summary presentation, inventory identity, and cleanup. Run
   `go test ./internal/reviewv1 ./cmd/subspace-tui -count=1`.
4. **Real cross-repository E2E:** create independent temporary main/state Git
   repositories, commit two sentinel files, run landed s4 `gate prepare`, commit the
   two-file room, relocate both checkouts, and make later commits that remove the
   worktree paths. Run `gate materialize`; capability-probe the exact built Subspace
   binary; drive one actual package-mode TUI in a private tmux session through Artifact,
   Reference, selector, and binding decision. Record and validate its retained outputs.
   Assert the canonical summary is visible, both old sentinel bodies were presented,
   inventory has exactly two original tuples, and all ephemeral bytes are gone.
5. **Retention negative:** repeat after deleting the final containing ref, running
   `git reflog expire --expire=now --all`, and running `git gc --prune=now`; also use a
   shallow clone that lacks the commit. Require nonzero materialization, zero provider
   launches, empty destination, and byte-identical room.
6. **Repository gates:** Spacedock runs `gofmt -w ./cmd ./internal`,
   `go test ./...`, and `go test ./... -race`; Subspace runs its Go and shell suites,
   `go test ./...`, and `go test ./... -race`; both run `git diff --check`. The detached
   audit verifies only Spacedock invokes `git`, no changed command contains remote
   acquisition, and no durable room/provider output contains payload sentinels.

## Documentation change proposal

```diff
--- docs/site/reference/command-reference.md
+++ docs/site/reference/command-reference.md
@@
+| `spacedock gate materialize ENTITY --room ROOM --destination DIR` | Resolve a bound recorder-ready room into one provider-owned ephemeral manifest. Reads only existing local `main`/`state` Git objects, verifies every raw SHA-256, publishes the manifest last, and never fetches or changes the retained room. |

--- cmd/subspace-tui/SPEC.md
+++ cmd/subspace-tui/SPEC.md
@@
+Package mode may additionally receive `--resolved-sources ABSOLUTE-MANIFEST` when
+`--supports review-v1-resolved-sources-v1` succeeds. The canonical Briefing remains
+unchanged. The manifest supplies verified ephemeral bytes for exact Git-root
+Artifact/Reference tuples; Subspace re-verifies every digest, derives inventory from
+the Briefing, and removes the resolved-source root before opening the TUI.

--- docs/review-and-gate.md
+++ docs/review-and-gate.md
@@
+A primary Artifact may carry a concise caller-authored `summary`. Providers preserve
+and present that exact value separately from source bytes. References do not acquire
+summaries, and providers do not synthesize them. A resolved-source manifest is an
+ephemeral byte-delivery profile, not a second Briefing or source identity.
```

## Out of scope

- q0 room-to-terminal routing, a second terminal launch, terminal discovery, or host
  adapter changes;
- compatibility wrappers, rewritten Briefings, generic URI registries, provider-side
  Git resolution, remote fetch/deepen/hydration, or cross-machine object acquisition;
- durable source caches, source payloads in gate/provider evidence, `association.json`,
  caller-selected Result/log/inventory paths, or a new request field;
- NLP summary generation, Reference summaries, authority capture, source mutation,
  suggestion application, or broader Review v1 schema work.

## Stage Report: ideation

- DONE: Exercise the actual current Spacedock Git-root Briefing to Subspace package-mode boundary and preserve the concrete failure that separates local object resolution from presentation.
  Spacedock `cc51e518` recorded and validated the exact two-repository Briefing; Subspace beta.4 `eeea575` and adjacent `6280ec4` capability-passed, then both exited 2 at filesystem `git-root:/main/...` resolution, while the Reference selector retained its explicit `://` rejection.
- DONE: Choose the smallest provider-neutral resolved-byte/materialization contract, assign exact ownership and files/LOC in each repository, and keep durable gate state free of copied source payloads.
  The design selects Spacedock-owned resolution plus a closed manifest and Subspace-owned private ephemeral payloads: Spacedock is 9 files/1,302 changed LOC, Subspace is 13 files/994 changed LOC, and the retained room remains exactly two metadata files.
- DONE: Define falsifiable moved-main/state end-to-end, inventory association, object-retention failure, and cleanup evidence sufficient for a new ideation gate.
  AC-1 through AC-5 require landed s4 output, independently moved roots, one real package-mode TUI, exact Artifact summary and source inventory, pruned/shallow failure before launch, and cleanup across completion, failure, signal, crash, and resume.
- DONE: Run the repository-required deterministic gates.
  `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` completed cleanly; any existing command, recorder, status, or integration regression would fail its package lane.

### Summary

The gate design proves q0 transport has no byte-delivery semantics and chooses one
provider-neutral manifest over provider-side Git or a rewritten Briefing. It preserves
canonical identities and the exact caller-authored Artifact summary, assigns cleanup to
the private provider package, and makes the first implementation proof a real
Spacedock-to-Subspace moved-root E2E.

## Ideation correction: production rendezvous (cycle 2)

This section supersedes the earlier command surface, cleanup claims, file/LOC tables,
acceptance criteria, E2E composition, and documentation proposal. The low-level
decision remains unchanged: the canonical Briefing is never rewritten; Spacedock owns
Git-root resolution and raw-SHA verification; a closed resolved-source manifest carries
verified ephemeral bytes; Subspace keeps Git and workflow topology out of the provider;
and no `association.json`, q0 dependency, generic resolver, Reference summary, or
generated prose is added.

The rejected design stopped between two real owners. Current Subspace
`invocation-common` allocates the provider package only inside a selected fixed entry,
while the proposed materializer expected a package path from an undefined caller. A
manual test that ran Spacedock and Subspace as independent binaries could not prove the
agent-facing gate route. The corrected design makes the existing `subspace:r` fixed
entry lifecycle that caller and uses the recorder's already-supported
`ROOM/provider/{result.json,review.jsonl,presented-inventory.json}` rendezvous.

### Exact baselines and post-s4 assumption

- Spacedock implementation baseline remains
  `cc51e518a3420b01fd4b455e9710d38803dc6d3e`.
- Subspace implementation baseline is
  `9e218f00a565e8353adbc834619140f1770783ba`.
- s4's summary-corrected post-baseline is its planned 19-file
  `+1,533/-161` surface. rqh estimates below are incremental after that surface lands;
  overlapping files are existing post-s4 files, not recreated files.
- The required s4 sentinel is exactly
  `  Résumé — validates Git-root presentation exactly.  `: two leading and two
  trailing ASCII spaces, composed `é`, and an em dash. Preparation and every retained
  canonical Briefing preserve those UTF-8 bytes without trimming or normalization.

At the Subspace baseline, `invocation_prepare_provider` in
`plugins/subspace/skills/r/scripts/invocation-common` allocates
`subspace-r-provider.XXXXXX`, and its private child receives that path through
`--provider-package`. At the Spacedock baseline, `recordRoomLocked` in
`internal/gates/operation.go` reads only
`ROOM/provider/result.json` and `ROOM/provider/presented-inventory.json`. Neither
repository currently moves, adopts, or discovers the other's package. The corrected
profile replaces the temporary allocator only for a Spacedock-bound room; ordinary
explicit Briefings keep their current temporary retained-package behavior.

### One supported agent-facing path

The selected-provider branch in post-s4
`skills/fo-gate-lifecycle/SKILL.md` invokes the installed `subspace:r` skill once with
this semantic request:

```text
--review-v1 --actor ACTOR --approver APPROVER \
  --spacedock-entity ENTITY \
  --spacedock-room ROOM \
  --spacedock-workflow-dir WORKFLOW_DIR \
  BRIEFING [TERMINAL]
```

The three `--spacedock-*` values are all-or-none and are valid only with
`--review-v1`. `ROOM` is the exact `room=` returned by landed s4; `BRIEFING` is the
exact bound Briefing locator in that room; `WORKFLOW_DIR` is the definition checkout
used to resolve the linked state root; and `ENTITY` is resolved by Spacedock within
that workflow. The ordinary explicit-Briefing form remains available for filesystem
Briefings and does not infer this profile from JSON or a URI. The optional terminal
retains the existing skill-level selection grammar and is never forwarded as semantic
gate data.

The skill selects one existing `review-<terminal>` entry and calls it once with the
same semantic argv. The entry sources
`plugins/subspace/skills/r/scripts/invocation-common`; no model-authored materializer,
provider-package path, TUI command, validator, recorder command, cleanup command, or
second entry is permitted.

Inside `invocation-common`, the Spacedock profile performs this exact sequence:

1. Parse the closed argv; canonicalize `WORKFLOW_DIR`, `ROOM`, and `BRIEFING`; reject
   a pre-existing `ROOM/provider`; allocate exactly that path as a new mode-0700
   non-symlink directory; and stamp the existing private diagnostics. Allocation stays
   before executable/capability/host preflight so the current explicit-Briefing rule
   continues to retain and report a package for every failure after semantic entry.
2. Resolve `${SPACEDOCK_BIN:-spacedock}` and `subspace-tui` to exact executable regular
   files. Probe both literal TUI capabilities,
   `review-v1-provider-package-v1` and
   `review-v1-resolved-sources-v1`, then run the selected host's existing structural
   preflight. No terminal opens yet.
3. Allocate the package's `resolved-sources` child as mode 0700. The caller never
   chooses Result, log, inventory, diagnostics, source payload, or manifest paths.
4. Invoke exactly:

   ```text
   ${SPACEDOCK_BIN:-spacedock} gate materialize ENTITY \
     --room ROOM --briefing BRIEFING \
     --destination ROOM/provider/resolved-sources \
     --workflow-dir WORKFLOW_DIR
   ```

   The command re-resolves `ENTITY`, requires `ROOM` to equal the current attempt's
   bound room, requires `BRIEFING` to equal the request locator and digest, derives the
   definition/state logical-root map, and requires `--destination` to be exactly the
   new private child of that bound room. Thus path possession alone is insufficient
   authority. It writes mode-0600 payloads, publishes
   `resolved-sources.json` last, and prints exactly `manifest=...` and `sources=N`.
   The entry captures stdout, stderr, exit status, and exact Spacedock argv under
   diagnostics and accepts only the expected manifest path and canonical source count.
5. Launch one existing Go provider supervisor through the selected host. Its child is:

   ```text
   subspace-tui --review-v1 --actor ACTOR --approver APPROVER \
     --provider-package ROOM/provider \
     --resolved-sources ROOM/provider/resolved-sources/resolved-sources.json \
     BRIEFING
   ```

   The supervisor also receives the exact
   `ROOM/provider/resolved-sources` cleanup root through private
   `--ephemeral-root` wiring. It validates that this is the one expected child of the
   provider package before starting the TUI.
6. Subspace reads the canonical Briefing unchanged, validates exact manifest coverage,
   reads and re-hashes every payload, installs Artifact and Reference bytes in memory,
   and removes the exact resolved-source child before setting the terminal title or
   opening Bubble Tea. Result and presented inventory continue to derive solely from
   the canonical Briefing.
7. The existing common lifecycle waits for exact child exit, calls
   `validate-one-file-result` once, returns its trusted bytes, and reports the retained
   provider package. Because that package is already the recorder's fixed
   `ROOM/provider`, the First Officer runs existing
   `spacedock gate record ENTITY --room ROOM --workflow-dir WORKFLOW_DIR` and
   `gate validate`; no adoption, copy, association file, or external-output flag is
   required.

This is one production path from the intended agent-facing gate choice through package
allocation, resolution, presentation, retained evidence, recording, and validation.
The public skill owns semantic routing; `invocation-common` owns allocation and
pre-dispatch orchestration; Spacedock owns room/root/object authority; the Go supervisor
and TUI own post-dispatch cleanup; and the existing recorder owns durable association.

### Closed manifest and summary behavior

The earlier `spacedock-resolved-sources` version 1 manifest shape and exact
`type/id/uri/mediaType/rev/path` rules remain unchanged. It covers every Git-root
Artifact and recursively reached Git-root Reference in canonical order, contains no
summary, and permits only clean contained regular payload paths. Subspace's inventory
contains canonical identity tuples only and never exposes a provider path.

Mandatory summary validation is scoped to s4-prepared, request-backed rooms. Arbitrary
request-less or advisory Briefings keep current behavior: Subspace displays a summary
when a primary Artifact supplies a JSON string but does not newly require one. The
canonical string remains unchanged in `Artifact.Extra["summary"]`; neither manifest,
inventory, Result, Reference, nor source bytes duplicate it.

Terminal rendering is lossless but never emits caller-authored terminal controls.
Printable Unicode code points and ordinary ASCII spaces, including the sentinel's
leading/trailing spaces, render unchanged. Backslash renders as `\\`; LF, CR, and TAB
render as `\n`, `\r`, and `\t`; every remaining Unicode control or format code point
(`Cc`/`Cf`, including ESC, BEL, DEL, bidi controls, and zero-width format controls) and
the `U+2028`/`U+2029` line/paragraph separators renders as uppercase `\u{HEX}`. No
normalization, trimming, ANSI interpretation, terminal-title interpolation, or
format-string use occurs. The reversible display encoding decodes to the exact
canonical UTF-8 string, while the bytes sent to the terminal contain no control from
the summary. A focused test uses spaces, composed and decomposed Unicode, backslash,
LF, TAB, ESC, BEL, RLO, and the two separators to prove both round-trip identity and
control-free output.

### Retention, cleanup, signals, crashes, and continuation

`ROOM/provider` is retained evidence/recovery state, not wholly transient scratch.
Existing explicit-Briefing behavior is preserved: after allocation, preflight,
materialization, host-launch, signal, child, validation, and delivery failures report
the recoverable provider-package path and retain diagnostics plus any provider evidence
already produced. Normal, feedback, binding, and open outcomes also retain Result, log,
inventory, child-exit, argv, capability, materializer, and stderr evidence for the
recorder. The earlier claim that every outcome leaves a two-file room was wrong; two
files is the pre-provider and no-override baseline only.

Only `ROOM/provider/resolved-sources` is always intended for deletion:

- before dispatch, `invocation-common` removes that exact child on materializer,
  capability, host-preflight, or launch failure, then preserves the rest of the package;
- after dispatch, the TUI removes it immediately after successful in-memory loading,
  before interactive presentation;
- if loading or the child fails earlier, the Go supervisor removes it after exact child
  exit and before atomically publishing child-exit evidence;
- on HUP, INT, or TERM, the supervisor forwards the signal to the exact child, waits,
  removes the child, then records the status. If only the invoking shell is interrupted,
  it reports recovery and does not race deletion against the still-running supervisor;
  cleanup occurs when that supervisor reaches its child-exit boundary;
- an uncatchable process-group kill or power loss can leave the private child. This is
  reported honestly as recoverable residue inside the already-named provider package;
  neither Git state nor another directory is scanned or deleted automatically.

Current Review v1 package mode has no supported resume command. A later user-directed
retry is a fresh one-entry invocation and is blocked while `ROOM/provider` exists; it
does not silently delete or reuse failure evidence. A future continuation feature must
take the emitted package explicitly, remove only its validated `resolved-sources`
child after proving no child is alive, re-resolve the canonical Git objects, and
re-verify all revisions. rqh does not claim that absent resume surface or hide hard-kill
residue behind “resume cleanup.”

### Negative boundary fixture

The beta.4 room-to-provider spike remains warranted as a committed negative boundary,
not as the success E2E. Subspace adds
`internal/reviewv1/testdata/git-root-negative.json` with one Git-root Artifact and one
Git-root Reference. The resolved-source tests prove ordinary
`review-v1-provider-package-v1` still fails that fixture at Artifact filesystem
resolution, the Reference selector still refuses unresolved `://`, and only the
literal resolved-source capability plus complete manifest advances. Removing the new
capability/manifest from the positive path must restore the recorded failure.

### Rebaselined expected surface

The Spacedock table is incremental after s4:

| File | Expected delta | Purpose |
|---|---:|---|
| `internal/cli/cli.go` | `+62/-8` | Route/document `gate materialize` and its exact room/Briefing/destination authority. |
| `internal/cli/gate_test.go` | `+125/-0` | Stable argv/stdout, bound-room, bound-Briefing, and exact-destination cases. |
| `internal/gates/materialize.go` (new) | `+215/-0` | Closed recursive manifest and atomic manifest-last publication. |
| `internal/gates/materialize_test.go` (new) | `+305/-0` | Coverage, room immutability, destination containment, and failure atomicity. |
| `internal/gitsource/source.go` | `+58/-8` | Expose s4's verified local blob bytes to materialization without a second resolver. |
| `internal/gitsource/source_test.go` | `+95/-0` | Moved-root, pruned/shallow, raw-SHA, and no-fetch materialization controls. |
| `internal/contractlint/fo_function_reference_invariant_test.go` | `+42/-8` | Pin one `subspace:r` semantic route and forbid direct TUI/materializer composition in FO prose. |
| `skills/fo-gate-lifecycle/SKILL.md` | `+36/-14` | Replace the selected-override halt with the one skill invocation and record/validate continuation. |
| `docs/specs/gate-resolution-frontmatter-contract.md` | `+72/-8` | Normative materialization authority, manifest, and retained-package boundary. |
| `docs/site/reference/command-reference.md` | `+22/-4` | Exact provider-facing command and failure semantics. |

Spacedock baseline: **10 named files, +1,032/-50 = 1,082 changed LOC**. Tolerance is
at most +2 genuinely new files and +20% changed LOC, hard cap **12 files / 1,298
changed LOC**. Every currently known file is named above; the tolerance cannot absorb
the public rendezvous, a second resolver, recorder adoption, provider executable,
request field, remote acquisition, or cache.

The Subspace table is likewise fully named:

| File | Expected delta | Purpose |
|---|---:|---|
| `plugins/subspace/skills/r/SKILL.md` | `+58/-18` | Agent-facing Spacedock profile, permission text, one-entry rule, and retention ownership. |
| `plugins/subspace/skills/r/scripts/invocation-common` | `+165/-38` | Parse authority, allocate `ROOM/provider`, invoke materialization, build exact child, and clean pre-dispatch payloads. |
| `scripts/tests/subspace-r-contract-test.sh` | `+38/-5` | Pin public grammar, fixed-entry ownership, capabilities, and forbidden manual composition. |
| `scripts/tests/subspace-r-provider-retained-delivery-test.sh` | `+190/-28` | Preflight/materialize/signal/child/validation retention and payload-cleanup matrix. |
| `scripts/tests/subspace-r-git-root-provider-e2e.sh` (new) | `+285/-0` | Real moved-root prepare-to-entry-to-record/validate journey. |
| `internal/reviewv1/model.go` | `+18/-4` | Verified in-memory Reference bytes; summary remains canonical extra data. |
| `internal/reviewv1/loader.go` | `+32/-6` | Share canonical validation with resolved input. |
| `internal/reviewv1/resolved_sources.go` (new) | `+220/-0` | Closed manifest, coverage, containment, and digest checks. |
| `internal/reviewv1/resolved_sources_test.go` (new) | `+325/-0` | Positive/negative fixture, adversarial tuples, controls, and cleanup. |
| `internal/reviewv1/testdata/git-root-negative.json` (new) | `+25/-0` | Committed unresolved Artifact/Reference boundary from the spike. |
| `internal/reviewv1/log.go` | `+25/-12` | Selector text uses verified in-memory Reference bytes. |
| `cmd/subspace-tui/main.go` | `+48/-10` | Private manifest flag and literal resolved-source capability. |
| `cmd/subspace-tui/profile_dispatch_test.go` | `+82/-0` | Capability and exact private argv surface. |
| `cmd/subspace-tui/provider_supervisor.go` | `+82/-16` | Validate one cleanup child, forward signals, wait, delete, then publish exit. |
| `cmd/subspace-tui/provider_supervisor_test.go` | `+145/-0` | Child failure, signal, invalid root, exact-exit, and cleanup-order proof. |
| `cmd/subspace-tui/v1_tui.go` | `+38/-8` | Resolved load and delete-before-TUI boundary. |
| `cmd/subspace-tui/v1_sources.go` | `+62/-10` | Reference rendering and lossless control-safe Artifact summary. |
| `cmd/subspace-tui/v1_sources_test.go` | `+115/-0` | Complete source catalog and no synthesized/normalized summary. |
| `cmd/subspace-tui/v1_review_chrome_labels_test.go` | `+92/-0` | Exact sentinel, safe controls, spacing, width, and title isolation. |
| `cmd/subspace-tui/SPEC.md` | `+48/-8` | Private resolved-source and cleanup lifecycle. |
| `docs/review-and-gate.md` | `+32/-4` | Public profile, summary display, recovery, and recorder rendezvous. |

Subspace baseline: **21 named files, +2,125/-167 = 2,292 changed LOC**. Tolerance is
at most +2 genuinely new files and +20% changed LOC, hard cap **23 files / 2,750
changed LOC**. No named rendezvous, supervisor, fixture, lifecycle test, or E2E file is
deferred into tolerance. Git commands, workflow discovery, a second package allocator,
generic transport, rewritten Briefing, or retained source cache reset ideation.

### Revised acceptance criteria

**AC-1 (VALUE) — The intended selected-provider gate route presents and records both
old Git-root sources after independent checkout movement.** Starting from landed s4,
the First Officer's declared Subspace route invokes one `subspace:r` semantic request.
The selected fixed entry allocates the bound `ROOM/provider`, materializes through
Spacedock, launches one real Subspace TUI, returns trusted evidence, and existing
`gate record --room`/`gate validate` close the attempt. The Artifact and Reference old
sentinels, exact `  Résumé — validates Git-root presentation exactly.  ` summary, and
two canonical inventory tuples are visible/retained after both current worktree paths
move. *Falsifier:* bypass the fixed entry, precreate the provider package manually,
call either binary directly, omit any room/root argument, or move provider outputs
from another package; the E2E's entry/argv/allocation/recorder assertions fail.

**AC-2 — Durable state contains provider evidence but no selected-source copy.** Before
provider invocation the room has exactly request + canonical Briefing. After a
recordable outcome it additionally has the existing provider Result/log/inventory and
diagnostics, while `resolved-sources` is absent. Catchable preflight, load, child,
signal, validation, and delivery failures retain recovery evidence but not the payload
child; an uncatchable group kill may retain only that named private residue and never
mutates gate frontmatter. *Falsifier:* remove each cleanup owner in turn; its lifecycle
case retains a sentinel payload after the owner reaches its defined boundary.

**AC-3 — Canonical association and summary survive the rendezvous without path or
control injection.** Manifest and presented inventory match every original
`type/id/uri/mediaType/rev` exactly once; no temporary path, summary, or association
file enters identity evidence. Printable summary text and spaces display unchanged;
unsafe code points use the reversible escape form, never raw terminal controls or
title bytes. *Falsifier:* mutate any tuple, payload, summary byte, spacing, Unicode
composition, escape, or temporary path and require failure or a changed expected
lossless view.

**AC-4 — Room/root/object authority fails before host launch and never acquires
objects.** Wrong entity, unbound room, mismatched Briefing locator/digest, non-exact
destination, unknown root, absent/pruned/shallow object, or raw-SHA mismatch yields
zero host launches, no manifest, no gate mutation, and no fetch/deepen/hydration.
*Falsifier:* a fake `git` records every argv and fails the test on remote acquisition;
destination and room hashes detect partial writes outside the owned package.

**AC-5 — Cleanup and recovery claims match current supported surfaces.** The shell owns
pre-dispatch cleanup, the TUI owns delete-after-load, and the Go supervisor owns
post-dispatch child/signal cleanup. Provider evidence remains on every existing
explicit-Briefing preservation path. The test explicitly records that SIGKILL can leave
private residue and that current Review v1 has no resume command; a fresh invocation
refuses to reuse the occupied package. *Falsifier:* any test that silently deletes
retained diagnostics, claims a nonexistent resume, scans other temp packages, or races
a live child is rejected.

### Real E2E and proof order

1. Land s4 and record its final implementation tip. Verify its post-baseline files and
   exact summary contract before rqh red tests.
2. Preserve the committed unresolved fixture as the first red: ordinary package mode
   fails Git-root Artifact/Reference resolution, while the resolved capability is
   absent from the invocation. This is a boundary fixture, not success evidence.
3. Red/green focused Spacedock tests for exact bound room/Briefing/destination
   authority, closed recursive manifest, moved roots, local-only object failures,
   atomic publication, and stable output.
4. Red/green focused Subspace tests for both literal capabilities, manifest coverage,
   Artifact/Reference bytes, summary safety, provider-supervisor cleanup ordering, and
   the current retention matrix.
5. Run `scripts/tests/subspace-r-git-root-provider-e2e.sh` with checked-out Spacedock
   and Subspace repositories supplied explicitly by the test job. It builds both exact
   binaries, creates independent main/state repositories, runs landed `gate prepare`,
   commits the two-file room, relocates both roots, deletes current worktree paths,
   and invokes `review-tmux` once using the public Spacedock profile. A private real
   tmux session drives Artifact, Reference, selector, and binding decision. The test
   asserts one entry, one package allocation at bound `ROOM/provider`, one materializer
   argv, one supervisor/TUI, visible old sentinels and exact printable summary, complete
   inventory, absent resolved child, then real `gate record --room` and `gate validate`.
   The script fails if replaced by manual `git cat-file`, direct `subspace-tui`,
   separate materialize/TUI calls, rewritten Briefing, copied provider outputs, or a
   provider fixture.
6. Run the retention negative after final-ref deletion, reflog expiry, and prune, plus
   a shallow clone missing the commit. Require failure before host launch, zero remote
   Git argv, no manifest, unchanged gate binding, and only the explicitly retained
   preflight diagnostics allowed by the lifecycle.
7. Run both repositories' full Go, race, shell, docs, formatting, and diff checks. A
   detached audit verifies only Spacedock invokes Git, the FO skill names only
   `subspace:r`, the Subspace skill calls exactly one fixed entry, summary bytes never
   enter a terminal title, and no room/provider tree retains selected-source payload
   after a catchable owner boundary.

### Corrected documentation changes

The implementation updates these exact public contracts:

```diff
--- skills/fo-gate-lifecycle/SKILL.md
+++ skills/fo-gate-lifecycle/SKILL.md
@@
-With a selected provider override, halt and name rqh.
+With a selected Subspace override, invoke `subspace:r` once with actor, approver,
+the exact s4 entity/room/workflow/Briefing values, and optional terminal. Accept only
+the returned trusted package at `ROOM/provider`; then run existing `gate record
+ENTITY --room ROOM` and `gate validate`. Do not allocate paths or compose binaries.

--- docs/site/reference/command-reference.md
+++ docs/site/reference/command-reference.md
@@
+| `spacedock gate materialize ENTITY --room ROOM --briefing BRIEFING --destination ROOM/provider/resolved-sources --workflow-dir DIR` | Resolve the current bound s4 room through local main/state Git objects into one closed provider-private manifest. All authority is revalidated; no fetch occurs; the canonical Briefing and gate state remain unchanged. |

--- plugins/subspace/skills/r/SKILL.md
+++ plugins/subspace/skills/r/SKILL.md
@@
+The all-or-none Spacedock Review v1 profile carries entity, bound room, workflow
+directory, and canonical Briefing to the same one-entry lifecycle. The entry allocates
+only the fixed `ROOM/provider`, calls Spacedock materialization, launches one
+supervised TUI, preserves ordinary provider evidence/recovery semantics, and removes
+only the resolved-source child at its defined ownership boundaries.

--- docs/review-and-gate.md
+++ docs/review-and-gate.md
@@
+Git-root presentation requires `review-v1-resolved-sources-v1` and the public
+Spacedock profile. Artifact summaries remain canonical strings: printable text and
+spaces render unchanged, while terminal-control and format code points use a
+reversible visible escape and never enter terminal control channels.
```

The normative Spacedock spec and Subspace `SPEC.md` carry the same allocation,
authority, cleanup, hard-crash, and no-current-resume boundaries. The existing
`gate record --room` contract is unchanged because the provider package now lands at
its already-supported location.

## Stage Report: ideation (cycle 2)

- DONE: Close the missing production rendezvous from the intended agent-facing gate route through recording.
  The corrected route is `fo-gate-lifecycle → subspace:r → one fixed entry → invocation-common → gate materialize → provider supervisor/TUI → validate-one-file-result → gate record --room/validate`; `ROOM/provider` is both the private package and the recorder's existing evidence location.
- DONE: Assign exact package allocation, room/root authority, launch, retention, and cleanup owners.
  `invocation-common` allocates the bound package and owns pre-dispatch cleanup, Spacedock revalidates entity/room/Briefing/root/destination and resolves bytes, the TUI deletes after load, and the supervisor deletes after child/signal failure while preserving all non-source recovery evidence.
- DONE: Correct crash/resume and Artifact-summary claims without weakening the accepted low-level contract.
  Catchable paths remove only `resolved-sources`; hard kill may leave named private residue, current Review v1 has no resume surface, and exact canonical summary bytes use lossless control-safe display with the coordinated two-space `Résumé` sentinel.
- DONE: Name and rebaseline every post-s4 implementation and proof file.
  Spacedock is 10 named files/1,082 changed LOC and Subspace is 21 named files/2,292 changed LOC; the real E2E enters through `review-tmux`, retains the negative fixture, and no rendezvous/supervisor/test file is hidden under tolerance.
- DONE: Run the repository-required deterministic gates against the unchanged implementation baseline.
  `gofmt -w ./cmd ./internal` produced no Go diff, and both `go test ./...` and `go test ./... -race` completed green; any current launcher, recorder, status, or skill-integration regression would fail its package lane.
- SKIPPED: Implement code, invoke a live approval provider, record a decision, mutate gate/status frontmatter, or dispatch another worker.
  This rejection cycle changes only the ideation entity body and requires a new independent staff review.

### Summary

Cycle 2 preserves the closed resolved-source design but gives it one real caller and
one existing recorder rendezvous. It also replaces impossible universal cleanup and
literal-control rendering claims with explicit ownership, recoverable hard-crash
residue, no invented resume, and lossless terminal-safe presentation.
