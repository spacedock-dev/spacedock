---
title: Materialize Git-root Review v1 sources for provider presentation
status: implementation
source: "s4 cycle-6 staff rejection: recorder-valid git-root:// sources are not renderable by current Subspace package mode, 2026-07-25"
started: 2026-08-15T21:25:13Z
completed:
verdict:
score: "1.0"
worktree: .worktrees/spacedock-ensign-git-root-review-v1-materialization
issue:
sprint: durable-decisions
id: rqh46ey33aqq4rt72b4w1m2q
gates:
    version: 1
    records:
        - id: gate:docs-dev:rqh4:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:rqh4-backlog-1
              briefing:
                id: briefing:docs-dev:rqh4:backlog:attempt-1:revision-1
                digest: sha256:d620934ee0af1b72c38e80fdb640f6ea07bd95da9fd08729c38e9b9d04a4fce2
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:rqh4:backlog:1
                briefing: briefing:docs-dev:rqh4:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-25T06:28:30.836574Z"
                decision: approve
                reason: The task isolates the actual missing cross-repository consumer boundary, forbids durable source duplication, and requires a real moved-root Subspace proof before implementation.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:rqh46ey33aqq4rt72b4w1m2q:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:rqh46ey33aqq4rt72b4w1m2q-ideation-1
              briefing:
                id: briefing:rqh46ey33aqq4rt72b4w1m2q:ideation:attempt-1:revision-1
                digest: sha256:615ddbffcf9c2b02784dd50258b9793fdc75edb4de29c9cc719bab28d560455a
                request-digest: sha256:072b1ae8eca461e3c9b1badcc5ce60c1b02fa51450b16a3dea206af894d442ec
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:rqh46ey33aqq4rt72b4w1m2q:ideation:1
                briefing: briefing:rqh46ey33aqq4rt72b4w1m2q:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-26T12:30:28.727811Z"
                decision: approve
                reason: The corrected design closes digest-domain, lifecycle-order, proof-policy, and circular-live-proof gaps. Record approval now; apply after s4 lands so both repositories implement against the final prepared-room contract.
              application:
                target-stage: implementation
                state: consumed
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
is estimated at **19 files, +1,533/-161 = 1,694 changed LOC**. rqh does not invent a
second summary input.

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

## Expected surface and planning estimates

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

Spacedock planning estimate: **9 files, +1,290/-12 = 1,302 changed LOC**. Counts are
reconciliation evidence only; they neither authorize nor reject a fixture or focused
test split and do not reset implementation.

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

Subspace planning estimate: **13 files, +942/-52 = 994 changed LOC**. Counts are
reconciliation evidence only; they neither authorize nor reject provider cleanup or
an E2E fixture split and do not reset implementation.

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

### Feedback Cycles

- Cycle 4: REVISE — Subspace FO + independent staff review; surface state-only review,
  0 product files/0 LOC vs estimate 32 files/3,717 LOC (0%); AC unchanged — make
  `/subspace:r gate <room>` the sole public entry, derive entity/workflow/authority/
  Briefing/package mechanics from the bound room, bind both canonical Briefing id and
  `sha256:` digest in the resolved-source manifest before display, preserve separate
  in-memory Artifact/Reference bytes and exact control-safe summary rendering, and
  rebaseline both repositories after corrected s4 lands.

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
- s4's summary-corrected post-baseline is state design
  `b739a0165590f111dbb88082b374468aee5b5985`, with its planned 19-file
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

Spacedock planning estimate: **10 named files, +1,032/-50 = 1,082 changed LOC**. Every
currently known file is named above. File/LOC variance requires reconciliation and
materiality review, never authorization or rejection by count.

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

Subspace planning estimate: **21 named files, +2,125/-167 = 2,292 changed LOC**. No
known rendezvous, supervisor, fixture, lifecycle test, or E2E file is omitted from the
table. File/LOC variance requires reconciliation and materiality review, never
authorization or rejection by count.

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
  Spacedock is estimated at 10 named files/1,082 changed LOC and Subspace at 21 named files/2,292 changed LOC; the real E2E enters through `review-tmux`, retains the negative fixture, and the planning tables name every known rendezvous/supervisor/test file.
- DONE: Run the repository-required deterministic gates against the unchanged implementation baseline.
  `gofmt -w ./cmd ./internal` produced no Go diff, and both `go test ./...` and `go test ./... -race` completed green; any current launcher, recorder, status, or skill-integration regression would fail its package lane.
- SKIPPED: Implement code, invoke a live approval provider, record a decision, mutate gate/status frontmatter, or dispatch another worker.
  This rejection cycle changes only the ideation entity body and requires a new independent staff review.

### Summary

Cycle 2 preserves the closed resolved-source design but gives it one real caller and
one existing recorder rendezvous. It also replaces impossible universal cleanup and
literal-control rendering claims with explicit ownership, recoverable hard-crash
residue, no invented resume, and lossless terminal-safe presentation.

## Ideation correction: validator and frozen authority (cycle 3)

This section supersedes only cycle 2's allocation order, materializer argv, validator
assumption, per-repository estimates, and affected AC/E2E cases. The approved
agent-facing `fo-gate-lifecycle → subspace:r → fixed entry → Spacedock → supervised
Subspace → gate record` rendezvous, fixed `ROOM/provider` recorder location, closed
manifest, unchanged canonical Briefing, exact/lossless summary display, cleanup and
hard-kill honesty, no q0, and no stored association remain unchanged.

### Two baseline omissions

At Subspace `9e218f00a565e8353adbc834619140f1770783ba`,
`plugins/subspace/skills/r/scripts/validate-one-file-result` authenticates the
supervisor's `argv.json` digest and then requires exactly this nine-element child argv:

```text
[SUBSPACE_TUI,
 "--review-v1","--actor",ACTOR,"--approver",APPROVER,
 "--provider-package",PROVIDER_ROOT,BRIEFING]
```

Cycle 2's resolved child has eleven elements because it inserts
`"--resolved-sources", MANIFEST` before `BRIEFING`. The existing validator would reject
that otherwise-correct Result after presentation. This is a canonical-validator
contract change, not an incidental test update.

Cycle 2 also let `invocation-common` allocate `ROOM/provider` before
`gate materialize`, while materialization received neither semantic identity. The
request already freezes `actor` and `approver`; passing different values could occupy
the one-shot provider location and probe or launch Subspace before Spacedock discovered
the mismatch during recording. Path equality and later recorder refusal are not launch
authority.

### Revised materialization and allocation order

The public `subspace:r` semantic request remains:

```text
--review-v1 --actor ACTOR --approver APPROVER \
  --spacedock-entity ENTITY --spacedock-room ROOM \
  --spacedock-workflow-dir WORKFLOW_DIR BRIEFING [TERMINAL]
```

For an s4-bound request, the First Officer copies `ACTOR` and `APPROVER` from the
frozen `request.json`; it does not substitute the invoking session, reviewer, or
Captain identity. The selected fixed entry passes those strings unchanged to
`invocation-common`.

Before resolving or capability-probing `subspace-tui`, running host preflight, creating
diagnostics, allocating `ROOM/provider`, or opening a terminal, `invocation-common`
resolves only `${SPACEDOCK_BIN:-spacedock}` and invokes exactly:

```text
${SPACEDOCK_BIN:-spacedock} gate materialize ENTITY \
  --room ROOM --briefing BRIEFING \
  --actor ACTOR --approver APPROVER \
  --destination ROOM/provider/resolved-sources \
  --workflow-dir WORKFLOW_DIR
```

Spacedock is the single pre-launch authority owner. `gate materialize` performs these
steps in order:

1. resolve the active entity through `WORKFLOW_DIR`, current definition/state roots,
   and the current open gate attempt;
2. require `ROOM` to equal the frozen room and require the canonical request digest,
   gate id, attempt id, Briefing locator/id/digest, `actor`, and `approver` to match that
   attempt exactly;
3. require `BRIEFING` to be the frozen locator and exact bytes, require
   `--destination` to equal the clean absolute
   `ROOM/provider/resolved-sources`, and require both `ROOM/provider` and the
   destination to be absent and non-symlink;
4. resolve every local Git-root object, verify every raw SHA-256, and build the complete
   closed manifest in memory or private error-clean staging, with no fetch or provider
   process;
5. only after all preceding checks succeed, atomically allocate mode-0700
   `ROOM/provider` and `resolved-sources`, write mode-0600 payloads, and publish the
   mode-0600 manifest last.

Wrong actor and wrong approver therefore fail before Git reads, provider-package
allocation, TUI resolution/capability probes, host preflight, or launch. The room
remains the byte-identical s4 request + Briefing pair, provider launch count is zero,
and `manifest=` is never printed.

On successful materialization, package ownership transfers to the Subspace fixed-entry
lifecycle. `invocation-common` creates/stamps diagnostics, checks both literal TUI
capabilities, performs host preflight, and launches the same supervised child. Any
failure after that transfer retains `ROOM/provider` as **diagnostic/recovery evidence**,
removes `resolved-sources` at the already-defined catchable ownership boundary, and
does not call the package a supported retry. Ordinary request-less explicit Briefings
keep the existing temporary allocator and nine-element validator mode.

### Canonical validator extension

`validate-one-file-result` retains its current ordinary interface and adds one
unambiguous resolved profile:

```text
validate-one-file-result --review-v1 --resolved-sources MANIFEST \
  BRIEFING BRIEFING_SHA256 ACTOR APPROVER \
  RESULT LOG INVENTORY CHILD_EXIT CAPTURE
```

In that profile, `MANIFEST` must be the clean absolute
`PROVIDER_ROOT/resolved-sources/resolved-sources.json` derived from `RESULT`; the
resolved-source directory must already be absent at successful validation because the
TUI/supervisor cleanup boundary has completed. The script verifies the existing
`argvSha256`, then requires exactly this eleven-element child argv:

```text
[SUBSPACE_TUI,
 "--review-v1","--actor",ACTOR,"--approver",APPROVER,
 "--provider-package",PROVIDER_ROOT,
 "--resolved-sources",MANIFEST,BRIEFING]
```

The ordinary `--review-v1 BRIEFING ...` mode continues to require the old exact
nine-element argv. The modes do not infer from Briefing content or accept one another's
shape. Both continue to verify exact Briefing bytes, Result/capture equality,
inventory projection, log/Resolution authority, and provider output paths.

`scripts/tests/subspace-r-provider-retained-delivery-test.sh` owns the behavioral
exact-argv matrix: resolved mode accepts only the eleven-element array; deleting,
duplicating, reordering, or changing either new element fails; an alternate manifest,
an old nine-element array under resolved mode, and an eleven-element array under
ordinary mode fail even with recomputed `argvSha256`. The positive case also proves
`resolved-sources` is absent before validator invocation.
`scripts/tests/subspace-r-contract-test.sh` pins the new private validator interface and
requires `invocation-common` to choose it only for the Spacedock profile.
The unrelated one-file matrix in
`scripts/tests/subspace-tui-agent-contract-fixture-test.sh` remains unchanged.

### Rebaselined expected surface

The complete Spacedock surface, incremental after s4 design
`b739a0165590f111dbb88082b374468aee5b5985`, is now:

| File | Expected delta | Purpose |
|---|---:|---|
| `internal/cli/cli.go` | `+72/-8` | Route/document identity-bound `gate materialize` and its exact destination. |
| `internal/cli/gate_test.go` | `+165/-0` | Stable argv/stdout plus wrong-actor/wrong-approver byte-clean CLI cases. |
| `internal/gates/materialize.go` (new) | `+235/-0` | Frozen request authority, local resolution, closed manifest, and allocate-last publication. |
| `internal/gates/materialize_test.go` (new) | `+350/-0` | Identity ordering, coverage, containment, allocation atomicity, and room immutability. |
| `internal/gitsource/source.go` | `+58/-8` | Expose s4's verified local blob bytes without a second resolver. |
| `internal/gitsource/source_test.go` | `+95/-0` | Moved-root, pruned/shallow, raw-SHA, and no-fetch controls. |
| `internal/contractlint/fo_function_reference_invariant_test.go` | `+42/-8` | Pin one semantic skill route and forbid direct binary composition. |
| `skills/fo-gate-lifecycle/SKILL.md` | `+36/-14` | Pass the request-frozen actor/approver and retain record/validate continuation. |
| `docs/specs/gate-resolution-frontmatter-contract.md` | `+80/-8` | Normative identity-before-allocation and manifest lifecycle. |
| `docs/site/reference/command-reference.md` | `+26/-4` | Exact command, authority, output, and failure semantics. |

Spacedock planning estimate: **10 named files, +1,159/-50 = 1,209 changed LOC**.
Identity checks, package allocation, request parsing, and their tests are all named
planning work. Actual variance requires explanation and materiality review, not a
count-based pass or failure.

The complete Subspace surface is now:

| File | Expected delta | Purpose |
|---|---:|---|
| `plugins/subspace/skills/r/SKILL.md` | `+58/-18` | Agent-facing Spacedock profile, frozen identity, permission text, and retention ownership. |
| `plugins/subspace/skills/r/scripts/invocation-common` | `+185/-50` | Invoke Spacedock before provider effects, select validator mode, launch, and cleanup. |
| `plugins/subspace/skills/r/scripts/validate-one-file-result` | `+72/-18` | Exact resolved manifest path and eleven-element child-argv authority. |
| `scripts/tests/subspace-r-contract-test.sh` | `+45/-5` | Pin public grammar, both validator modes, capabilities, and one-entry ownership. |
| `scripts/tests/subspace-r-provider-retained-delivery-test.sh` | `+245/-30` | Exact-argv mutations, identity preflight, retention, signal, and cleanup matrix. |
| `scripts/tests/subspace-r-git-root-provider-e2e.sh` (new) | `+315/-0` | Real moved-root public entry plus wrong-identity zero-launch cases. |
| `internal/reviewv1/model.go` | `+18/-4` | Verified in-memory Reference bytes; canonical summary extra data. |
| `internal/reviewv1/loader.go` | `+32/-6` | Shared canonical validation for resolved input. |
| `internal/reviewv1/resolved_sources.go` (new) | `+220/-0` | Closed manifest, coverage, containment, and digest checks. |
| `internal/reviewv1/resolved_sources_test.go` (new) | `+325/-0` | Positive/negative fixture, adversarial tuples, controls, and cleanup. |
| `internal/reviewv1/testdata/git-root-negative.json` (new) | `+25/-0` | Committed unresolved Artifact/Reference boundary. |
| `internal/reviewv1/log.go` | `+25/-12` | Selector text uses verified in-memory Reference bytes. |
| `cmd/subspace-tui/main.go` | `+48/-10` | Private manifest flag and literal resolved-source capability. |
| `cmd/subspace-tui/profile_dispatch_test.go` | `+82/-0` | Capability and exact private TUI argv surface. |
| `cmd/subspace-tui/provider_supervisor.go` | `+82/-16` | Validate cleanup child, forward signals, wait, delete, then publish exit. |
| `cmd/subspace-tui/provider_supervisor_test.go` | `+145/-0` | Failure, signal, invalid root, exact-exit, and cleanup-order proof. |
| `cmd/subspace-tui/v1_tui.go` | `+38/-8` | Resolved load and delete-before-TUI boundary. |
| `cmd/subspace-tui/v1_sources.go` | `+62/-10` | Reference rendering and lossless control-safe Artifact summary. |
| `cmd/subspace-tui/v1_sources_test.go` | `+115/-0` | Complete catalog and no synthesized/normalized summary. |
| `cmd/subspace-tui/v1_review_chrome_labels_test.go` | `+92/-0` | Exact sentinel, safe controls, spacing, width, and title isolation. |
| `cmd/subspace-tui/SPEC.md` | `+48/-8` | Private resolved-source, validator, and cleanup lifecycle. |
| `docs/review-and-gate.md` | `+32/-4` | Public profile, summary display, recovery evidence, and recorder rendezvous. |

Subspace planning estimate: **22 named files, +2,309/-199 = 2,508 changed LOC**. The
validator script and both exact-argv test owners are named planning work. Actual
variance requires explanation and materiality review, not a count-based pass or
failure.

### Revised acceptance and proof deltas

**AC-1 (VALUE)** retains cycle 2's real public fixed-entry moved-root presentation and
recording proof. It additionally requires canonical validator success through the
resolved profile; substituting the old validator or nine-element argv makes the same
E2E fail after presentation and before trusted delivery.

**AC-2** retains diagnostic/recovery evidence and cleanup semantics. A failed package
is never described as a retry surface. Wrong actor/approver allocate no package at all;
post-materialization failures retain the fixed package but remove the payload child at
the defined catchable boundary.

**AC-3** retains canonical tuple, no-path/no-summary inventory, exact summary, and
control-safe display proof. Validator mode adds no manifest content or path to Result,
inventory, association, or canonical Briefing.

**AC-4 — Frozen semantic identity gates every provider effect.** Wrong actor and wrong
approver independently fail inside Spacedock before Git reads, `ROOM/provider`
allocation, TUI resolution/capability probe, host preflight, supervisor/TUI launch, or
gate mutation. *Test:* focused Spacedock cases hash the two-file room and count Git
calls; public fixed-entry cases count materializer, capability, host, supervisor, and
TUI calls and assert only one rejecting Spacedock invocation. Reordering identity
validation after allocation or dropping either comparison makes its case fail.

**AC-5** retains cycle 2's explicit cleanup/hard-kill/no-current-resume boundaries and
adds exact validator-mode separation. Recomputed argv digests cannot make an alternate
flag, path, count, or ordering canonical.

The real E2E still enters through one `review-tmux` fixed entry, not manual primitive
composition. Before its positive moved-root run, it changes only `ACTOR` and only
`APPROVER` in two separate public invocations and requires nonzero status, byte-identical
request/Briefing, absent `ROOM/provider`, zero TUI capability/host/supervisor launches,
and no gate/status mutation. The positive run then requires the eleven-element
resolved child argv, successful canonical validator delivery, absent payload child,
exact old Artifact/Reference bytes and summary sentinel, retained inventory, and real
`gate record --room`/`gate validate`.

The focused Subspace proof mutates all eleven argv positions, the manifest spelling,
mode, and digest. The focused Spacedock proof mutates request actor/approver and semantic
argv actor/approver independently. Repository gates and the negative unresolved fixture
remain as specified in cycle 2.

## Stage Report: ideation (cycle 3)

- DONE: Add the canonical validator and exact resolved-source child argv to the named Subspace contract.
  `validate-one-file-result` now has an explicit resolved profile requiring the fixed manifest path and eleven-element argv; the provider-retained and contract tests reject cross-mode, missing, duplicate, reordered, alternate-path, and recomputed-digest variants.
- DONE: Bind actor and approver to the frozen request before every provider side effect.
  The revised materializer receives both identities, validates request/gate/attempt/Briefing/identity/destination and all Git bytes, then allocates `ROOM/provider`; separate wrong-actor and wrong-approver cases require an unchanged two-file room and zero provider calls.
- DONE: Preserve the approved rendezvous, evidence lifecycle, summary, and canonical source boundaries.
  Public `subspace:r`/fixed-entry E2E, `ROOM/provider` recording, catchable payload cleanup, honest hard-kill residue, diagnostic/recovery wording, exact summary escaping, closed manifest, unchanged Briefing, no q0, and no stored association remain unchanged.
- DONE: Rebaseline every known implementation and proof file after the two corrections.
  Spacedock is estimated at 10 named files/1,209 changed LOC; Subspace at 22 named files/2,508 changed LOC, including the validator itself and both exact-argv test owners as named planning work.
- DONE: Run the repository-required deterministic gates against the unchanged implementation baseline.
  `gofmt -w ./cmd ./internal` produced no Go diff, and `go test ./...` plus `go test ./... -race` completed green; any current CLI, recorder, status, or skill-integration regression would fail its package lane.
- SKIPPED: Implement code, invoke a live approval provider, record a decision, mutate gate/status frontmatter, or dispatch another worker.
  Cycle 3 changes only the ideation entity body and requires another independent staff review.

### Summary

Cycle 3 closes the last two pre-implementation gaps without changing the selected
architecture. Frozen semantic identity now precedes package allocation and every
provider action, and trusted delivery recognizes only the exact resolved-source child
argv selected by the fixed entry.

## Ideation correction: room-only authority and full Briefing binding (cycle 4)

Cycle 4 supersedes the cycle 2 and cycle 3 public argument vectors, caller-supplied
materializer coordinates, manifest Briefing field, surface estimates, and acceptance
wording wherever they conflict with this section. The resolved-source package,
provider-owned evidence, cleanup, validator separation, and unchanged-canonical-
Briefing decisions still stand.

Implementation remains blocked on corrected s4. rq must consume the landed s4 room
resolver and request/root-map contract rather than duplicate it. If landed s4 cannot
derive every coordinate below from a room alone, places flat and folder entities
differently from this correction, or changes the frozen summary/Briefing contract, rq
returns to joint s4/rq ideation. It must not add public flags as a compatibility
fallback.

### Sole public provider entry and fixed derivation

The complete agent-facing provider invocation is:

```text
/subspace:r gate <room>
```

There are no optional arguments. The agent supplies no entity, workflow directory,
actor, approver, Briefing, destination, provider path, manifest path, terminal vector,
or request fields. It does not open or parse `request.json`. rq owns this public
provider grammar; s4 owns room preparation and the metadata contract, not a second
provider invocation.

The Subspace skill's fixed `gate` branch canonicalizes the one room operand and invokes
one integration-private Spacedock operation:

```text
${SPACEDOCK_BIN:-spacedock} gate materialize --room ROOM
```

That operation accepts no caller-selected authority or source coordinate. Before any
Git read, provider directory allocation, TUI lookup/capability probe, host preflight,
diagnostic write, supervisor launch, or TUI launch, fixed Spacedock code uses the
landed s4 room resolver to derive and validate all of the following:

1. canonical room, entity slug and form, entity path, state root, workflow definition
   root, current gate and attempt, and the split-root map;
2. the two immutable authoritative room metadata files, their bound request, canonical
   Briefing locator, canonical Briefing id, and canonical Briefing SHA-256;
3. request-frozen actor and approver, plus the exact recorder/validator rendezvous;
4. provider root `ROOM/provider`, payload child
   `ROOM/provider/resolved-sources`, and manifest path
   `ROOM/provider/resolved-sources/resolved-sources.json`; and
5. the local Git-root Artifact/Reference tuples and every repository/object coordinate
   needed to resolve them without fetch.

The required s4 placement is collision-free for both entity forms:
`<entity-root>/<slug>/review/<stage>/briefing-N` is the room; a folder entity is
`<entity-root>/<slug>/index.md`, while a flat entity is
`<entity-root>/<slug>.md`. The resolver rejects zero, two, symlinked, escaped, or
noncanonical entity candidates. rq does not infer an entity by searching arbitrary
parents and does not serialize a caller-derived root map.

Only after room/request/Briefing/identity/root-map validation and complete in-memory
Git resolution succeed may Spacedock atomically create mode-0700 `ROOM/provider` and
its payload child, write mode-0600 payloads, and publish the mode-0600 manifest last.
It prints the derived private launch tuple only on success. The fixed Subspace branch
owns that output directly; it is not an agent copy/paste interface. It then resolves
the TUI, probes the two literal capabilities, runs host preflight, and launches the
existing supervised child with the derived actor, approver, provider root, manifest,
and original Briefing path.

Immediately after successful s4 preparation, a room therefore has exactly two
immutable authoritative metadata files and zero copied payloads. After provider
presentation begins, `ROOM/provider` may exist and is provider-owned evidence:
Result, log, inventory, and diagnostics may remain there. Catchable exits remove
`ROOM/provider/resolved-sources`; an uncatchable hard kill may honestly leave that
payload child. No retained package is described as a supported retry surface, and
current-attempt discovery never resumes it.

### Closed manifest and full Briefing identity

The closed manifest changes its Briefing member from a string to this required object:

```json
{
  "type": "spacedock-resolved-sources",
  "version": "1",
  "briefing": {
    "id": "<canonical Briefing id>",
    "digest": "sha256:<64 lowercase hexadecimal characters>"
  },
  "items": []
}
```

`items` retains cycle 2's exact closed Artifact/Reference catalog and payload binding.
Unknown, missing, duplicate, alternate-spelling, or wrong-type members fail closed.
The digest is the full SHA-256 of the canonical Briefing's unchanged bytes, not a
prefix, path hash, rewritten-document hash, or summary hash.

Before it installs any resolved bytes into the review model or displays any review
chrome, Subspace reads the original canonical Briefing, derives its canonical id,
recomputes the full byte digest, and requires both values to match the manifest and the
room-bound values passed by fixed code. Any mismatch fails before display and removes
the payload child at the catchable boundary. Subspace never rewrites, annotates, copies,
or associates the Briefing. Verified Artifact/Reference bytes exist only in memory
after loading; the canonical Briefing remains byte-identical.

The primary Artifact summary is required only for s4-prepared request-backed
Briefings. Its exact whitespace and printable Unicode survive unchanged. UI controls
use one reversible safe rendering: `\` becomes `\\`; LF, CR, and TAB become `\n`,
`\r`, and `\t`; other control/format code points and U+2028/U+2029 become
`\u{HEX}` with uppercase hexadecimal. This transformation is display-only. Neither
the manifest nor retained inventory stores the summary, and no normalized or
synthesized string can replace it.

### Corrected acceptance and proof

**AC-1 (VALUE) — one room-only public path works after a real root move.** A
black-box E2E prepares an s4 request-backed room, moves the main and state roots to new
absolute locations, and invokes exactly `/subspace:r gate ROOM`. It runs the real TUI,
shows the old-commit Artifact and Reference bytes plus the exact summary sentinel,
records through `gate record --room`, and passes `gate validate`. The test has folder-
entity and flat-entity cases so room-to-entity derivation cannot depend on one shape.
Any public flag beyond `gate ROOM`, manual primitive composition, network fetch, q0,
`association.json`, copied payload in the prepared room, or rewritten Briefing fails.

**AC-2 — the room is the complete authority input.** Focused Spacedock tests mutate
each room/request/root-map/attempt/identity/Briefing binding while passing the same
room operand. Every case fails before Git reads or provider effects and leaves the
two authoritative files byte-identical with no `ROOM/provider`. CLI and skill contract
tests reject legacy `--actor`, `--approver`, `--spacedock-entity`,
`--spacedock-workflow-dir`, Briefing, destination, manifest, and terminal arguments.

**AC-3 — canonical Briefing id and digest gate display.** Positive tests bind the
canonical id and full digest, load the original file, and preserve its bytes. Negative
tests independently mutate id, every digest position, digest length/case/prefix,
Briefing bytes, manifest spelling, duplicates, and unknown members. Each fails before
model installation or display. Exact-space/Unicode and control-rich summary fixtures
prove lossless model data, reversible safe rendering, and absence from inventory.

**AC-4 — provider evidence and payload cleanup have distinct lifetimes.** Success and
every catchable error/signal remove `resolved-sources` before trusted delivery while
retaining only provider-owned evidence appropriate to the reached phase. Pre-allocation
failures create nothing. A hard-kill fixture may observe residue and documentation says
so. Result/log/inventory continue to name no manifest, payload path, local repository
path, or summary.

**AC-5 — production validator authority remains separate.** The real fixed entry, not
an E2E-only path, uses the canonical resolved-profile validator. It verifies exact
derived child argv, Briefing bytes and digest, Result/capture equality, inventory
projection, log/Resolution authority, provider paths, and cleanup ordering. Recomputed
argv hashes cannot authorize an alternate path, flag, count, order, actor, approver,
manifest, or Briefing. The ordinary request-less one-file profile remains separate and
does not gain the s4 summary guarantee.

The moved-root E2E runs after the focused folder/flat resolver, manifest parser,
digest, summary, allocation-order, supervisor, and validator matrices. Repository
gates remain `gofmt -w ./cmd ./internal`, `go test ./...`, and
`go test ./... -race` in both repositories. The committed unresolved Git-root fixture
still proves no fallback to current bytes or network. Tests use the public room-only
entry for behavior and private seams only for focused fault injection.

### Rebaselined cross-repository implementation surface

This estimate was re-read against Spacedock main
`4ff98d8cd97` and the current Subspace tree at `63a26f63a3de`; the selected Subspace
provider/code paths were last changed by `cac4eb106f`. It is incremental after
corrected s4 lands. A changed s4 resolver or metadata shape triggers the explicit
return-to-ideation rule above; file/LOC variance does not.

| Spacedock file | Expected delta | Purpose |
|---|---:|---|
| `internal/cli/cli.go` | `+48/-8` | Integration-private room-only materializer route; reject extra coordinates. |
| `internal/cli/gate_test.go` | `+125/-0` | Stable room-only argv/stdout and pre-side-effect rejection matrix. |
| `internal/gates/materialize.go` (new) | `+210/-0` | Consume s4 resolver, validate all authority, resolve bytes, allocate last. |
| `internal/gates/materialize_test.go` (new) | `+390/-0` | Folder/flat rooms, binding mutations, coverage, containment, and atomicity. |
| `internal/gitsource/source.go` | `+58/-8` | Return verified local blob bytes through the canonical resolver. |
| `internal/gitsource/source_test.go` | `+95/-0` | Moved-root, pruned/shallow, raw-SHA, and no-fetch controls. |
| `internal/contractlint/fo_function_reference_invariant_test.go` | `+35/-8` | Pin the one room-only semantic route and forbid coordinate reconstruction. |
| `skills/fo-gate-lifecycle/SKILL.md` | `+24/-14` | Invoke `/subspace:r gate ROOM`; retain record/validate continuation. |
| `docs/specs/gate-resolution-frontmatter-contract.md` | `+78/-6` | Normative derivation, full Briefing binding, and room/provider lifecycle. |
| `docs/site/concepts/gates-and-decisions.md` | `+20/-4` | Explain the sole public entry and provider evidence without exposing private argv. |

Spacedock planning estimate: **10 named files, +1,083/-48 = 1,131 changed LOC**. The
s4 resolver itself is a landed dependency, not rq implementation LOC. The estimate is
non-authoritative and is reconciled after implementation.

| Subspace file | Expected delta | Purpose |
|---|---:|---|
| `plugins/subspace/skills/r/SKILL.md` | `+42/-24` | Define exactly `/subspace:r gate ROOM` and no caller coordinates. |
| `plugins/subspace/skills/r/scripts/invocation-common` | `+170/-55` | Call room-only materialization, consume derived output, launch, and clean up. |
| `plugins/subspace/skills/r/scripts/validate-one-file-result` | `+76/-18` | Canonical resolved profile with exact derived argv and full Briefing binding. |
| `scripts/tests/subspace-r-contract-test.sh` | `+60/-8` | Pin sole public grammar, both validator modes, and fixed-entry ownership. |
| `scripts/tests/subspace-r-provider-retained-delivery-test.sh` | `+235/-30` | Argument rejection, exact argv, retention, signal, and cleanup matrix. |
| `scripts/tests/subspace-r-git-root-provider-e2e.sh` (new) | `+380/-0` | Real moved-root folder/flat public-entry presentation and recording. |
| `internal/reviewv1/model.go` | `+18/-4` | Verified in-memory Reference bytes and exact summary data. |
| `internal/reviewv1/loader.go` | `+32/-6` | Shared pre-display canonical Briefing validation. |
| `internal/reviewv1/resolved_sources.go` (new) | `+245/-0` | Closed manifest, id/full-digest binding, coverage, and containment. |
| `internal/reviewv1/resolved_sources_test.go` (new) | `+370/-0` | Digest/id/parser adversaries, exact catalog, controls, and cleanup. |
| `internal/reviewv1/testdata/git-root-negative.json` (new) | `+25/-0` | Committed unresolved Artifact/Reference boundary. |
| `internal/reviewv1/log.go` | `+25/-12` | Selector text uses verified in-memory Reference bytes. |
| `cmd/subspace-tui/main.go` | `+48/-10` | Private manifest flag and literal resolved-source capability. |
| `cmd/subspace-tui/profile_dispatch_test.go` | `+90/-0` | Capability and exact private TUI argv surface. |
| `cmd/subspace-tui/provider_supervisor.go` | `+82/-16` | Validate cleanup child, forward signals, wait, delete, then publish exit. |
| `cmd/subspace-tui/provider_supervisor_test.go` | `+145/-0` | Failure, signal, invalid root, exact-exit, and cleanup-order proof. |
| `cmd/subspace-tui/v1_tui.go` | `+38/-8` | Resolve manifest, verify Briefing, and delete before TUI display. |
| `cmd/subspace-tui/v1_sources.go` | `+62/-10` | Reference rendering and reversible control-safe Artifact summary. |
| `cmd/subspace-tui/v1_sources_test.go` | `+115/-0` | Complete catalog and no synthesized/normalized summary. |
| `cmd/subspace-tui/v1_review_chrome_labels_test.go` | `+105/-0` | Exact Unicode/spacing sentinel, safe controls, width, and title isolation. |
| `cmd/subspace-tui/SPEC.md` | `+55/-8` | Private resolved-source, digest, validator, and cleanup lifecycle. |
| `docs/review-and-gate.md` | `+38/-4` | Room-only profile, summary display, evidence, and recorder rendezvous. |

Subspace planning estimate: **22 named files, +2,456/-213 = 2,669 changed LOC**. The
fixed entry, canonical validator, and both exact-argv test owners are named planning
work. The estimate is non-authoritative and is reconciled after implementation.

### Documentation diff

Spacedock's normative gate contract will define the collision-free folder/flat room
mapping, room-derived authority, canonical Briefing id/full digest, two-file prepared
state, provider-owned evidence, and catchable versus hard-kill cleanup. Its gate
concept page will show only `/subspace:r gate <room>` as the provider handoff and will
not advertise private materializer argv.

Subspace's skill and review guide will define the same sole public invocation, the
request-backed-only summary guarantee, reversible control rendering, retained evidence,
and no-retry wording. The TUI spec will own private manifest schema, pre-display
id/digest verification, in-memory byte injection, exact validator profile, and cleanup
ordering. No document introduces q0, `association.json`, a compatibility wrapper, or a
rewritten Briefing.

## Stage Report: ideation (cycle 4)

- DONE: Make the sole public provider entry room-only, with fixed code deriving all
  workflow, authority, Briefing, package, validator, and recorder coordinates before
  side effects. The public grammar is exactly `/subspace:r gate <room>`; legacy caller
  flags and positional Briefing/terminal inputs are explicit negative tests.
- DONE: Bind the resolved-source manifest to canonical Briefing id and full SHA-256,
  verify both before display, keep resolved bytes ephemeral in memory, and leave the
  canonical Briefing byte-identical. Exact summary whitespace/Unicode and reversible
  control-safe rendering apply only to s4-prepared request-backed Briefings.
- DONE: Rebaseline corrected cross-repository files/LOC and the real moved-root E2E
  behind s4. Spacedock is 10 named files/1,131 changed LOC; Subspace is 22 named
  files/2,669 changed LOC; the public E2E covers folder and flat entity rooms.
- DONE: Preserve `ROOM/provider` evidence, catchable payload cleanup, honest hard-kill
  residue, canonical validator separation, local Git-object authority, and the bans on
  q0, association storage, copied prepared-room payloads, and rewritten Briefings.
- DONE: Define the dependency boundary. Implementation waits for corrected s4 and
  returns to ideation if room-only derivation, placement, root mapping, or frozen
  summary/Briefing binding lands differently; it never compensates with public flags.
- SKIPPED: Implement either repository, launch a provider, record a gate, mutate
  gate/status frontmatter, or change s4. Cycle 4 edits only this rq ideation entity and
  requires the next independent staff review.

### Summary

rq now has one authority-bearing input: the prepared room. Fixed Spacedock/Subspace
code derives and validates every private coordinate before provider effects; the
manifest proves the unchanged canonical Briefing by id and full digest; and the
corrected, post-s4 surface includes folder/flat moved-root proof and exact safe summary
display.

## Feedback Cycle 5

Independent staff review `f35b93d0` returned **REVISE** on one Material issue: cycle 4
promoted file and changed-LOC forecasts into hard implementation limits. Counts do not
prove the value, authority, safety, or compatibility contract and cannot authorize or
reject an implementation. Optimizing for a limit could hide a useful helper split;
remaining below it could still conceal a semantic expansion.

The review affirmed the cycle 4 architecture in full: room-only public entry, fixed
derivation, corrected folder/flat placement, canonical Briefing id/full digest,
unchanged Briefing, exact request-backed summary, two-file prepared state, local-only
Git resolution, provider-owned ephemeral payload, catchable cleanup, honest hard-kill
residue, and separate production validator authority.

## Ideation correction: estimates are reconciliation evidence (cycle 5)

All file lists and expected deltas in this entity are planning estimates, including
the latest cycle 4 Spacedock and Subspace tables. They help the implementer declare
scope and help reviewers spot drift. They are not budgets, gates, ceilings, floors, or
permission boundaries. No file count, addition/deletion count, percentage, or variance
can by itself authorize, reject, pause, or reset implementation.

Before product edits, the implementation worker must append an intended-surface
declaration for each repository. It starts from the cycle 4 named-file table and calls
out every already-known addition, omission, rename, merge, or helper split, with the
semantic responsibility that explains it. The declaration is a review aid, not a
promise to preserve the forecast.

After implementation and before staff review, the worker must reconcile actual
path-scoped `git diff --stat` and `git diff --numstat` results against that declaration
and the planning estimates. Reconciliation identifies:

- planned files that changed, did not change, moved, merged, or split;
- unplanned files and the contract or proof responsibility that required them;
- actual additions/deletions by repository; and
- whether each variance is non-semantic structure/proof work or exposes a semantic
  expansion.

Reviewers triage the explanation and implementation behavior for materiality. Large
variance can justify scrutiny but not count-based rejection; zero variance confers no
approval. Repository tests, contract tests, fixture tests, E2E tests, CI, and review
scripts must not assert total changed files, total changed LOC, estimated deltas, or a
percentage deviation. They assert the cycle 4 behavior and authority boundaries.

rq returns to ideation only if implementation proposes or discovers one of these
semantic expansions:

1. another public argument beyond `/subspace:r gate <room>`;
2. any caller-selected authority;
3. inability to derive mechanics from the bound room;
4. a third authoritative prepared metadata file;
5. copied source payloads in the prepared room;
6. remote source acquisition;
7. `association.json`;
8. a compatibility path;
9. weaker canonical Briefing identity or digest binding; or
10. a changed provider/validator ownership boundary.

The corrected-s4 dependency check remains semantic: if landed s4 changes room
resolution, folder/flat placement, request/root-map authority, the two-file preparation
state, or exact summary/Briefing binding such that cycle 4 cannot be implemented, that
is the applicable inability/authority/metadata/binding expansion and returns to joint
s4/rq ideation. A table mismatch alone never does.

No cycle 4 runtime, manifest, display, cleanup, evidence, validator, or acceptance
decision changes in cycle 5. The implementation remains blocked until corrected s4
lands.

## Stage Report: ideation (cycle 5)

- DONE: Remove the Spacedock and Subspace hard file/LOC caps and any implication that
  tolerance can authorize, reject, or reset implementation.
  Every count paragraph now labels its totals as non-authoritative planning estimates;
  count variance alone cannot pass, fail, pause, authorize, or reopen the design.
- DONE: Keep named-file tables/deltas as planning estimates and require intended
  surface declaration plus post-implementation reconciliation/materiality triage.
  The implementation worker declares planned paths before editing and reconciles
  actual per-repository stat/numstat variance before staff review.
- DONE: Return to ideation only for the semantic expansions named by the reviewer.
  Cycle 5 records the exact ten triggers and makes corrected-s4 divergence use those
  semantic categories rather than a numeric proxy.
- DONE: Ensure no tests assert repository file/LOC totals.
  Tests and gates prove cycle 4 behavior and authority; total files, total LOC,
  estimated deltas, and variance percentages are explicitly forbidden assertions.
- DONE: Preserve Cycle 4 architecture unchanged.
  Room-only entry, fixed derivation, full Briefing binding, two-file preparation,
  provider evidence/cleanup, safe summary display, and validator ownership all stand.
- SKIPPED: Implement product code, run a provider, record a decision, mutate
  gate/status frontmatter, or touch s4.
  Cycle 5 changes only rq ideation state and requires the next independent staff
  review after this correction.

### Summary

Cycle 5 removes numeric authority without weakening scope visibility. The named
cross-repository tables remain useful forecasts, while declaration, reconciliation,
and semantic materiality—not file or LOC totals—govern implementation review.

## Feedback Cycle 6

The post-cycle-5 integration read against corrected s4 found two Material proof
defects. First, cycle 4's manifest used one `digest` member while describing it at
different points as a raw-file SHA-256 and as the canonical Briefing digest. Corrected
s4 binds the Briefing with SHA-256 over RFC 8785/JCS bytes; formatting makes that value
deliberately different from the raw file SHA-256. rq must preserve both domains
without overloading one field.

Second, the named proof stopped at a structural First Officer contract lint and a
provider-entry E2E. It did not name the existing live First Officer composition that
observes one `/subspace:r gate <room>` invocation and the subsequent
`gate record ... --room <room>` continuation. Structural lint remains useful, but it
cannot prove that a live First Officer follows the route. The correction must reuse
the existing release-candidate walking lane and add no second agent harness.

Cycle 5's count ruling remains fully in force. The same review also preserves the
room-only public entry, two-file prepared room, arbitrary request-located canonical
Briefing, provider-owned evidence subtree, local-only Git resolution, exact summary,
separate validator ownership, and the bans on copied prepared-room payloads,
`association.json`, and compatibility arguments.

## Ideation correction: distinct JCS and raw Briefing identity (cycle 6)

This section supersedes cycle 4's resolved-manifest `briefing` object, every statement
that calls its digest a hash of unchanged/raw Briefing file bytes, the cycle 4 surface
tables, and the affected AC/test wording. All other cycle 4 behavior and all cycle 5
planning/reconciliation rules remain authoritative.

The post-s4 baseline inspected for this correction is Spacedock main
`50f8d1fb7b0cbc40747622d9f9d95467a0bec6c0`, corrected s4 candidate
`e328ecc6c118d1380ca36eadf82aad558b72e7af`, and Subspace
`5466d601b7281a8d91715a1d03e190e7b5049c56`. rq implementation still starts only
after corrected s4 lands. The separate prepared-Briefing basename regression remains
s4-owned; rq consumes the landed room resolver and never repairs it with a
`briefing.json` fallback.

### Closed manifest with separate digest domains

The resolved-source manifest retains version 1 because no earlier manifest has shipped,
but its closed Briefing member is now:

```json
{
  "type": "spacedock-resolved-sources",
  "version": "1",
  "briefing": {
    "id": "briefing:task:validation:attempt-1:revision-1",
    "jcsDigest": "sha256:<64 lowercase hexadecimal characters>",
    "rawSha256": "sha256:<64 lowercase hexadecimal characters>"
  },
  "items": []
}
```

`jcsDigest` is copied only after Spacedock recomputes SHA-256 over the RFC 8785/JCS
serialization of the exact Briefing located by the frozen s4 request and requires it
to equal the request, attempt, and entity binding. It is the existing canonical
Briefing authority; it remains full length and is never replaced by the file hash.

`rawSha256` is independently computed over the exact located file bytes used for this
materialization. It is a provider-handoff pin: Subspace must receive and parse the same
file bytes Spacedock inspected. It is not a second recorder authority and is never
compared with, substituted for, or silently promoted to `jcsDigest`.

The ambiguous manifest member `briefing.digest` is forbidden, not retained as an
alias. `digestDomain`, prefixes, truncated values, uppercase hexadecimal, missing
members, duplicate members at any depth, and unknown members also fail closed. There
is no compatibility parser because this manifest is still unshipped.

Before manifest publication, Spacedock:

1. resolves the room through landed s4 and validates its request, current attempt,
   actor/approver, root map, and exact clean relative Briefing locator;
2. reads that located regular non-symlink Briefing without appending or requiring a
   basename, rejects duplicate JSON members, derives its id, and recomputes its full
   JCS digest;
3. requires the id and JCS digest to equal every frozen s4 binding, then computes the
   separate raw file SHA-256;
4. resolves and raw-verifies every Git-root source before provider allocation; and
5. publishes both distinct Briefing pins and the closed item catalog with the manifest
   last.

Before model installation, selector construction, summary rendering, or display,
Subspace reopens the original located Briefing, rejects duplicate members, derives its
id, computes RFC 8785/JCS SHA-256 and raw file SHA-256 independently, and requires all
three values to equal the manifest. It then validates the exact source catalog and
payload pins. A successful parse never rewrites or copies the Briefing.

The resolved validator keeps its raw-file revision check but names that value
`BRIEFING_RAW_SHA256` in its interface documentation and tests. Its raw comparison
must not stand in for the loader's JCS check. Exact resolved child argv, canonical
inventory, Result/capture equality, log/Resolution authority, and cleanup ordering
remain independently validated as in cycle 3.

### Exact indented s4 control

The focused positive control is the exact s4-prepared
`durable-decisions-release-walking-skeleton/review/backlog/briefing-1/gate-briefing.json`
from state commit `684d6603c985030c3c6031f4b1a1462c0f1cbfa1`, blob
`de790a44aa5b47ecdd606fc711a0ad8f9f20a2d7`. It is a 31-line, 1,675-byte indented
JSON file with these deliberately unequal identities:

```text
rawSha256=sha256:c3b6d4d5ac8c766dcc56e08b57a41e207147d1319c61f066160e4e7d4bacfb1b
jcsDigest=sha256:0782c65c06c7ee9378226b3a7ef88d92939a54c05d916fe3690cc7d99804278f
```

The second value is the existing full canonical digest stored by s4 in both
`request.json` and the attempt binding. It is the SHA-256 of the 1,493-byte RFC
8785/JCS serialization, not the SHA-256 of the 1,675 raw file bytes.

The positive focused path stages the exact file with its matching s4 request, invokes
room-only materialization, requires a manifest containing the exact id,
`jcsDigest=...278f`, and `rawSha256=...fb1b`, and then loads it through Subspace before
display. The test fails if an implementation compares the two digests for equality,
puts the raw value into `jcsDigest`, puts the JCS value into `rawSha256`, truncates
either value, accepts the ambiguous old `digest` spelling, or omits either
recomputation.

Additional controls make the domain distinction observable:

- reindent an otherwise semantically identical copy: JCS remains `...278f` while the
  raw pin changes; a manifest regenerated from that exact copy succeeds, while the old
  raw pin fails;
- mutate one semantic JSON string without changing the request: JCS changes and
  Spacedock refuses before allocation;
- mutate only manifest `jcsDigest` or only `rawSha256`: Subspace refuses before model
  installation and the catchable owner removes `resolved-sources`; and
- place an equivalently bound canonical Briefing at another clean request locator,
  including a nested non-`briefing.json` filename: room-only materialization and
  validation succeed without basename reconstruction.

The fixture remains exactly the s4 two-file prepared room. It adds no source payload,
provider output, `association.json`, rewritten Briefing, or third authoritative
prepare-time file.

### Behavioral proof and the live First Officer lane

Proof is deliberately layered rather than assigned to contract lint:

1. `internal/contractlint/fo_function_reference_invariant_test.go` remains a
   structural check that `present-gate` names the one room-only invocation and
   recorder continuation in order. It does not claim the First Officer performed
   them and gains no behavioral transcript parser.
2. Focused Spacedock and Subspace tests prove room authority, both digest domains,
   arbitrary locator resolution, manifest closure, source bytes, exact validator
   mode, and cleanup using falsifiable byte/state controls.
3. `scripts/tests/subspace-r-git-root-provider-e2e.sh` proves the real fixed entry,
   materializer, provider supervisor/TUI, moved Git roots, canonical inventory, and
   `gate record --room` composition with production binaries. It is not replaced by
   direct primitive calls.
4. The already-defined Track B in
   `durable-decisions-release-walking-skeleton/index.md` is the live First Officer
   lane. With corrected s4 and rq pinned, a cold Shaping FO invokes exactly
   `/subspace:r gate <emitted-room>`, accepts only the retained package at the
   room-derived provider path, continues with the existing
   `gate record {rescope-target} --room <emitted-room>`, commits the closure, and
   stops without consuming it. Its nine-column retained run table records the exact
   skill entry, recorder command, pre/post state, and evidence digests.

The walking lane is the agent-behavior proof; the cross-repository E2E is the product
composition proof; contract lint is only a structural guard. rq adds no fake provider,
test-only Subspace skill, transcript oracle, standing controller, or new live harness.
The walking lane already declares its own dependency stop conditions and runs only
against exact landed candidates. Until corrected s4 lands, a successful end-to-end
composition remains explicitly deferred rather than simulated.

### Post-s4 expected surface and reconciliation

Corrected s4 already owns the arbitrary Briefing locator, room/request resolver,
canonical digest, `gitsource.Resolve`, two-file preparation, First Officer handoff
prose, and structural contract lint. rq consumes those landed surfaces unchanged.

The current Spacedock implementation forecast after s4 is:

| Spacedock file | Expected delta | Purpose |
|---|---:|---|
| `internal/cli/cli.go` | `+45/-4` | Route the integration-private room-only materializer and stable success/failure output. |
| `internal/cli/gate_test.go` | `+150/-0` | Reject extra coordinates and prove pre-side-effect room-only CLI behavior. |
| `internal/gates/materialize.go` (new) | `+275/-0` | Reuse s4 authority/JCS/Git readers, compute the separate raw pin, and publish the closed manifest last. |
| `internal/gates/materialize_test.go` (new) | `+460/-0` | Exact indented control, folder/flat/arbitrary locator, digest-domain mutants, catalog, containment, and atomicity. |
| `internal/gates/testdata/materialize-s4-room/gate-briefing.json` (new) | `+31/-0` | Exact indented s4 positive-control Briefing. |
| `internal/gates/testdata/materialize-s4-room/request.json` (new) | `+13/-0` | Matching frozen request/JCS binding for the two-file control. |
| `docs/specs/gate-resolution-frontmatter-contract.md` | `+65/-4` | Normative room-only materialization, distinct digest domains, and provider evidence lifecycle. |
| `docs/site/concepts/gates-and-decisions.md` | `+18/-2` | Explain Git-root presentation without exposing integration-private argv. |

Spacedock planning estimate: **8 named files, +1,057/-10 = 1,067 changed LOC**.
`internal/gitsource/source.go`, `skills/fo-gate-lifecycle/SKILL.md`,
`skills/present-gate/SKILL.md`, and
`internal/contractlint/fo_function_reference_invariant_test.go` are explicit
zero-delta reused s4 surfaces unless implementation finds a semantic defect. The
contract-lint file remains structural; no live behavior is moved into it.

The current Subspace forecast is:

| Subspace file | Expected delta | Purpose |
|---|---:|---|
| `plugins/subspace/skills/r/SKILL.md` | `+42/-24` | Define exactly `/subspace:r gate ROOM`, room ownership, permission text, and retention. |
| `plugins/subspace/skills/r/scripts/invocation-common` | `+175/-55` | Consume room-only materialization output, keep raw/JCS identities distinct, launch, and clean up. |
| `plugins/subspace/skills/r/scripts/validate-one-file-result` | `+88/-18` | Exact resolved profile, explicitly raw Briefing pin, and eleven-element child argv. |
| `scripts/tests/subspace-r-contract-test.sh` | `+60/-8` | Pin sole public grammar, validator modes, capabilities, and fixed-entry ownership. |
| `scripts/tests/subspace-r-provider-retained-delivery-test.sh` | `+270/-30` | Digest-domain/argv mutations, retention, signals, validation, and cleanup matrix. |
| `scripts/tests/subspace-r-git-root-provider-e2e.sh` (new) | `+420/-0` | Real moved-root folder/flat fixed-entry presentation and recorder continuation. |
| `go.mod` | `+1/-0` | Add the same maintained RFC 8785/JCS implementation used by Spacedock. |
| `go.sum` | `+2/-0` | Pin that canonicalization dependency. |
| `internal/reviewv1/model.go` | `+18/-4` | Verified in-memory Reference bytes and exact summary data. |
| `internal/reviewv1/loader.go` | `+38/-6` | Share duplicate-safe canonical Briefing parsing and id validation. |
| `internal/reviewv1/resolved_sources.go` (new) | `+275/-0` | Closed manifest, separate JCS/raw recomputation, coverage, and containment. |
| `internal/reviewv1/resolved_sources_test.go` (new) | `+430/-0` | Exact s4 positive control, domain mutants, catalog adversaries, and cleanup. |
| `internal/reviewv1/testdata/s4-prepared-gate-briefing.json` (new) | `+31/-0` | Exact indented `...fb1b` raw versus `...278f` JCS control. |
| `internal/reviewv1/testdata/git-root-negative.json` (new) | `+25/-0` | Committed unresolved Artifact/Reference boundary. |
| `internal/reviewv1/log.go` | `+25/-12` | Selector text uses verified in-memory Reference bytes. |
| `cmd/subspace-tui/main.go` | `+52/-10` | Private manifest flag and literal resolved-source capability. |
| `cmd/subspace-tui/profile_dispatch_test.go` | `+100/-0` | Capability, arbitrary Briefing filename, and exact private TUI argv. |
| `cmd/subspace-tui/provider_supervisor.go` | `+82/-16` | Validate cleanup child, forward signals, wait, delete, then publish exit. |
| `cmd/subspace-tui/provider_supervisor_test.go` | `+145/-0` | Failure, signal, invalid root, exact exit, and cleanup-order proof. |
| `cmd/subspace-tui/v1_tui.go` | `+45/-8` | Verify both Briefing digests and delete payloads before display. |
| `cmd/subspace-tui/v1_sources.go` | `+62/-10` | Reference rendering and reversible control-safe Artifact summary. |
| `cmd/subspace-tui/v1_sources_test.go` | `+115/-0` | Complete catalog and no synthesized/normalized summary. |
| `cmd/subspace-tui/v1_review_chrome_labels_test.go` | `+105/-0` | Exact Unicode/spacing sentinel, safe controls, width, and title isolation. |
| `cmd/subspace-tui/SPEC.md` | `+65/-8` | Private manifest schema, digest domains, validator, and cleanup lifecycle. |
| `docs/review-and-gate.md` | `+42/-4` | Room-only profile, canonical summary, evidence, and recorder rendezvous. |

Subspace planning estimate: **25 named files, +2,713/-213 = 2,926 changed LOC**.
RFC 8785 is an interoperability and authority boundary; `encoding/json` re-encoding or
sorted-key `jq` is not substituted for the standard. Reusing the maintained dependency
is smaller and safer than a second provider-local canonicalizer.

These are cycle-5 reconciliation aids only. Before edits, each implementation worker
declares its actual intended paths and responsibilities against the landed s4 and
Subspace tips. Before review, it reconciles path-scoped `diff --stat`/`--numstat`.
Counts, percentages, helper splits, fixture splits, and dependency-file variance
cannot authorize, reject, pause, or reset implementation. Only cycle 5's semantic
triggers do.

### Revised acceptance criteria and proof order

**AC-1 (VALUE) — the room-only provider path presents and records both old Git-root
sources after independent root movement.** Starting from landed corrected s4, the
production fixed entry receives only `/subspace:r gate ROOM`, derives every private
coordinate, displays the old Artifact and Reference plus exact summary, retains the
canonical inventory, and continues through `gate record --room`. Folder and flat
entities both succeed. *Falsifier:* any caller coordinate, manual materializer/TUI
composition, current-worktree fallback, provider-output copy, or missing recorder
continuation makes the cross-repository E2E fail. The existing walking-skeleton Track
B separately observes one real First Officer issue exactly that skill entry and then
the recorder command.

**AC-2 — the prepared room and canonical Briefing remain singular authority.** Before
provider work the room contains exactly the request and its arbitrarily located
canonical Briefing, with no selected-source copy. Later retained files are confined to
provider-owned evidence, while catchable owners remove `resolved-sources`.
*Falsifier:* require `briefing.json`, append a basename, add a third authoritative
prepare-time file, rewrite/copy the Briefing, or retain a payload after a catchable
boundary; the room/locator/lifecycle matrix fails.

**AC-3 — canonical and raw Briefing identities are explicit, independent, and checked
before display.** The manifest carries canonical id, full RFC 8785/JCS `jcsDigest`,
and full raw-file `rawSha256`. Spacedock and Subspace recompute each in its own domain;
the exact indented s4 control succeeds with `...278f != ...fb1b`. *Falsifier:* swap
the values, require equality, accept the old `digest` field, mutate either pin, alter
semantic JSON, or reuse a stale raw pin after reindentation; each focused case fails
at its declared pre-display boundary.

**AC-4 — source, provider, and recorder evidence preserve canonical association
without an association file.** Every Artifact and recursively reached Reference maps
once by `type/id/uri/mediaType/rev`; no temporary path, summary, raw/JCS helper field,
or local repository path enters Result or presented inventory. Result/inventory raw
pins and the exact validator mode remain provider evidence. *Falsifier:* mutate,
duplicate, omit, reorder, or expose any tuple/path/digest field; manifest, inventory,
validator, or recorder comparison fails before a binding Resolution.

**AC-5 — failure and proof ownership stay on current supported surfaces.** Missing or
pruned local objects, wrong room/current attempt, digest mismatch, capability failure,
loader failure, signal, and child failure retain only phase-appropriate evidence and
never fetch. Contract lint checks structure; focused tests and the fixed-entry E2E
check product behavior; the existing walking lane checks live First Officer behavior.
*Falsifier:* a structural grep is offered as runtime proof, a new fake-provider agent
harness is added, chat fallback occurs after selected launch, or a catchable path
leaves the payload child.

Proof order is:

1. land corrected s4 and pin its final tip; verify the arbitrary-locator resolver,
   two-file prepare state, JCS binding, summary, and local Git resolver;
2. preserve the unresolved beta.4 Git-root fixture as the provider boundary red;
3. make the exact indented s4 control red for the old ambiguous digest contract, then
   green only with distinct `jcsDigest` and `rawSha256`;
4. run focused Spacedock room/authority/materialization and Subspace
   parser/loader/summary/supervisor/validator matrices;
5. run the real moved-root cross-repository fixed-entry E2E through recording;
6. run the existing live First Officer walking lane when its pinned dependency
   conditions are satisfied, retaining its one skill entry and recorder continuation
   in the release-candidate table; and
7. run both repositories' full tests, race tests, shell/docs checks, formatting, and
   diff checks.

No success composition is claimed against current main before corrected s4 lands. No
product implementation, provider invocation, gate/status mutation, s4 edit, or
walking-skeleton run occurs during this ideation correction.

## Stage Report: ideation (cycle 6)

- DONE: Correct the Briefing digest contract so the resolved-source manifest explicitly carries and recomputes the existing full-length RFC 8785/JCS canonical Briefing digest; never use one field for raw and canonical identity.
  The closed manifest now has separate `jcsDigest` authority and `rawSha256` handoff pins; both sides recompute both before display, and ambiguous `digest` is rejected.
- DONE: Add an exact indented s4-prepared gate-briefing.json control whose raw SHA-256 differs from its JCS digest, and make the intended positive path falsifiable.
  State commit `684d6603`/blob `de790a44` supplies the 1,675-byte control: raw `c3b6…fb1b`, JCS `0782…278f`; swaps, equality, reindent, semantic drift, and old-field mutants are named.
- DONE: Add the applicable existing live First Officer lane proving one room-only /subspace:r gate <room> invocation and recorder continuation; keep contractlint structural and add no harness.
  The existing release-candidate walking-skeleton Track B owns the live FO issue plus `gate record --room`; contract lint remains structural, and the production fixed-entry E2E owns product composition.
- DONE: Preserve cycle 5: numeric counts remain non-authoritative planning/reconciliation evidence, with semantic reset triggers only.
  Both post-s4 forecasts are explicitly advisory, require pre-edit declaration and post-edit reconciliation, and cannot pass, fail, pause, authorize, or reopen work by count.
- DONE: Reconcile the post-s4 named surface and planning estimates, append a complete ideation Stage Report, and commit the state path only.
  The forecast is 8 Spacedock files/1,067 changed LOC and 25 Subspace files/2,926 changed LOC, with s4-owned resolver/Git/FO/contractlint files identified as zero-delta reuse.
- DONE: Run the repository-required deterministic gates against the unchanged product baseline.
  `gofmt -w ./cmd ./internal` left no tracked product diff; `go test ./...` and `go test ./... -race` both passed, so this prose-only correction inherits a green current-main baseline.
- SKIPPED: Implement product code, run a provider or live First Officer, mutate gate/status frontmatter, touch s4, or add an agent harness.
  Cycle 6 changes only this ideation entity; successful composition remains deferred until the corrected s4 dependency lands.

### Summary

Cycle 6 removes the last Briefing identity ambiguity by carrying the frozen RFC
8785/JCS digest separately from an exact raw-file handoff pin and proving the
distinction with a real indented s4 artifact. It also assigns agent behavior to the
existing live walking lane, leaves contract lint structural, and rebaselines the
post-s4 implementation surface without restoring numeric authority.

## Feedback Cycle 7

Independent re-review returned **REVISE** on three narrow ownership/order defects in
cycle 6. The JCS/raw separation and exact indented control are correct and remain
unchanged.

First, `durable-decisions-release-walking-skeleton` (`ph`) cannot prove rq before rq
terminates because its declared preconditions require rq to be terminal `PASSED`.
ph remains the later sprint-wide final composition. rq itself must run one real,
pre-terminal First Officer drive on landed s4 plus the rq candidates, through the real
installed skills and an existing supported runtime, with no new harness.

Second, cycle 6 inherited an incorrect allocation order from cycle 4. The Subspace
fixed entry's `invocation-common` lifecycle already owns allocation of
`ROOM/provider`; Spacedock materialization must not take that ownership. Complete
payload materialization and manifest-last publication still precede TUI loading,
model installation, and display, but they do not precede provider-root allocation.

Third, structural contract lint cannot prove that an agent or process performed
commands in semantic order. Its assertions must be limited to structural
closure/presence and forbidden-surface absence. The production fixed-entry E2E and
rq's one-off live First Officer drive own observable sequencing.

## Ideation correction: allocation and pre-terminal live ownership (cycle 7)

This section supersedes cycle 6's allocation chronology, assignment of the live
First Officer proof to ph, contract-lint sequencing claim, affected acceptance/proof
wording, and surface estimates. Cycle 6's closed manifest with distinct
`jcsDigest`/`rawSha256`, the exact `...278f != ...fb1b` indented control, arbitrary
canonical Briefing locator, and full negative matrix remain unchanged. Cycle 5's
numeric non-authority and semantic-only reset triggers also remain unchanged.

### Provider allocation and manifest-last order

The sole public agent entry remains exactly:

```text
/subspace:r gate <room>
```

The model supplies no other coordinate. The real installed Subspace skill selects its
fixed `gate` branch, and `invocation-common` performs this ownership sequence:

1. canonicalize the one room operand, reject an unsafe/non-directory/symlinked room,
   derive the exact fixed provider root `ROOM/provider`, and refuse an occupied or
   symlinked provider root;
2. atomically allocate mode-0700 `ROOM/provider`; from this point
   `invocation-common` owns its retention, diagnostics, and pre-dispatch failure
   behavior;
3. invoke exactly `${SPACEDOCK_BIN:-spacedock} gate materialize --room ROOM`;
4. after materialization returns success, resolve/probe the real TUI capabilities,
   perform the existing host preflight, and launch the supervised child with the
   materializer's derived original Briefing and resolved-manifest paths; and
5. run the resolved validator and return only trusted retained evidence to the
   invoking First Officer.

The integration-private materializer accepts only `--room ROOM`; it derives
`ROOM/provider/resolved-sources` from landed s4 authority and requires the already
allocated `ROOM/provider` to be a mode-0700 non-symlink directory. It neither creates,
removes, selects, nor relocates the provider root.

Inside that provider root, Spacedock validates the current room/request/attempt,
arbitrary Briefing locator, full JCS binding, separate raw file pin, actor/approver,
root map, and every local Git object. It writes payloads only under a private candidate
child. Every payload is complete, mode 0600, contained, and raw-revision verified
before Spacedock publishes `resolved-sources.json` last and makes the exact
`resolved-sources` child visible. A failure before publication removes the candidate
and any resolved-source child but leaves `ROOM/provider` to its Subspace owner.

Only a complete published manifest can reach Subspace. The TUI then independently
recomputes Briefing id/JCS/raw identity, validates catalog coverage and payload bytes,
installs all verified bytes in memory, deletes the resolved-source child at its
defined boundary, and only then constructs selectors, summary chrome, or display.
Thus source materialization precedes manifest publication, and manifest publication
precedes TUI/model installation/display; provider-root allocation deliberately
precedes materialization.

This order changes the failure-state wording but not durable authority. Immediately
after s4 preparation and before a selected provider entry, the room still contains
exactly its two authoritative metadata files. Once `/subspace:r gate ROOM` begins,
`ROOM/provider` may exist even when semantic room validation, object resolution,
materialization, capability, or launch later fails. It is provider-owned
diagnostic/recovery evidence, never a third authoritative prepared file. Catchable
failures remove `resolved-sources`; they do not pretend the provider root was never
allocated.

Path-safety or occupied-provider checks can fail before allocation. Room/request,
Briefing, identity, Git-object, digest, capability, and launch failures occur after
allocation and retain only the phase-appropriate provider evidence. No failure
mutates the request, located Briefing, entity gate state, Result binding, or source
Git objects; no path fetches, deepens, hydrates, or falls back to worktree bytes.

### rq-owned pre-terminal live First Officer drive

rq owns one one-off live drive before rq can receive a terminal `PASSED` verdict. It
is validation evidence, not a standing test or a substitute for the production
fixed-entry E2E.

The drive has these preconditions:

- corrected s4 is landed and its exact commit is an ancestor of the Spacedock rq
  candidate;
- the Spacedock and Subspace rq candidate commits are clean, their required offline,
  race, shell, docs, and fixed-entry E2E lanes are green, and exact candidate binary
  SHA-256 values are recorded;
- a fresh supported First Officer runtime is launched through that candidate
  Spacedock binary, with the candidate Spacedock plugin and the real candidate
  Subspace plugin installed through their ordinary plugin mechanisms; and
- rq has one real current prepared request-backed room committed by landed s4, with
  local main/state Git objects available and no existing `ROOM/provider`.

The Captain/operator gives that fresh First Officer ordinary-language intent selecting
Subspace for the prepared rq gate. The First Officer loads the real
`fo-gate-lifecycle` and `present-gate` skills, then the runtime's real installed
`subspace:r` skill. It must issue exactly one semantic skill invocation:

```text
/subspace:r gate <exact emitted room>
```

It must not use Bash to invoke Subspace, call materialization/TUI primitives directly,
install a test skill, reconstruct coordinates, probe/fallback in the First Officer,
or render a chat decision after the selected entry. The real fixed entry allocates
`ROOM/provider`, materializes through the candidate Spacedock binary, launches the
candidate TUI, and returns trusted binding evidence.

After that one skill invocation returns successfully, the same live First Officer
continues through the existing recorder surface:

```text
${SPACEDOCK_BIN:-spacedock} gate record <rq-entity> \
  --room <exact emitted room> --workflow-dir <workflow-dir>
```

It requires exit 0 with the current gate/attempt/Briefing and closed binding decision,
then commits the rq entity/room unit through the existing path-scoped state commit.
The drive stops before `gate consume`, status advance, successor dispatch, terminal
verdict, or merge. An advisory/open/invalid Result, second skill invocation, chat
fallback, recorder omission, alternate room, or frontmatter edit fails the drive.

Evidence is appended to this rq entity under one `### Pre-terminal live First Officer
drive` subsection before the terminal validation report. It records:

- exact Spacedock/s4/Subspace commits, binary SHA-256 values, and installed skill blob
  ids;
- runtime/version, native session/transcript locator and raw digest, and the single
  observed native `subspace:r` skill event with its exact room-only argument;
- room tree digest before entry, provider Result/inventory/raw evidence digests after
  entry, absence of `resolved-sources`, and absence of `association.json`;
- the exact recorder command/exit, entity and state-commit before/after ids, and the
  resulting provider-evidence pins; and
- elapsed time plus any friction and its owner.

This is a manual one-off drive on existing runtime and provider substrates. rq adds no
Go live test, workflow job, shell controller, fake provider, test-only skill,
transcript parser, fixture framework, retry loop, or sidecar schema. An independent
reviewer checks the native event and retained on-disk evidence directly; a prose claim
without those citations is not completion.

`ph` remains unchanged. After rq is terminal `PASSED`, its Track B independently
repeats the route as part of the sprint-wide role-separated release-candidate
composition. ph does not backfill, waive, or replace rq's pre-terminal drive.

### Structural lint and behavioral sequencing owners

`internal/contractlint/fo_function_reference_invariant_test.go` may assert only that
the installed First Officer/present-gate instruction closure:

- contains the room-only `/subspace:r gate <gate-room>` form and the semantic
  `gate record <entity> --room <gate-room>` closure;
- contains no caller-selected entity/workflow/Briefing/actor/approver/destination/
  provider/manifest/terminal coordinate on the public Subspace form; and
- contains no direct materializer/TUI composition, provider probe, selected-provider
  chat fallback, or `association.json` instruction.

It must not use heading order, ordered markers, substring position, command count, or
prose sequence to claim prepare-before-entry, one invocation, entry-before-recorder,
allocation order, recorder success, or no fallback at runtime. Those are semantic
behavior claims.

The production fixed-entry E2E owns process sequencing: one entry call, provider-root
allocation by `invocation-common`, materialization and manifest-last publication,
TUI/display, validator, cleanup, and recorder result. The rq pre-terminal live drive
owns agent sequencing: one real installed skill invocation followed by the recorder
continuation and no chat fallback. Each proof has state/byte/event mutants that fail
if its sequence changes.

### Corrected post-s4 surface estimates

The cycle 6 Spacedock table remains except that contract lint is now an explicit
correction surface:

| Additional/corrected Spacedock file | Expected delta | Purpose |
|---|---:|---|
| `internal/contractlint/fo_function_reference_invariant_test.go` | `+25/-15` | Keep structural closure/absence checks and remove semantic ordering claims. |

The full Spacedock planning estimate is therefore **9 named files,
+1,082/-25 = 1,107 changed LOC**. `internal/gitsource/source.go`,
`skills/fo-gate-lifecycle/SKILL.md`, and `skills/present-gate/SKILL.md` remain
zero-delta s4 reuse unless behavior reveals a semantic defect. The one-off live drive
adds retained evidence to this entity, not a product or harness file.

The cycle 6 Subspace table remains with these allocation/sequence reconciliations:

| Corrected Subspace file | Revised expected delta | Corrected purpose |
|---|---:|---|
| `plugins/subspace/skills/r/scripts/invocation-common` | `+190/-55` | Allocate `ROOM/provider`, invoke room-only materialization, consume complete output, then preflight/launch/validate/clean. |
| `scripts/tests/subspace-r-provider-retained-delivery-test.sh` | `+285/-30` | Prove allocation-before-materialization, manifest-last, retained-root, catchable-child cleanup, signals, and validator modes. |
| `scripts/tests/subspace-r-git-root-provider-e2e.sh` (new) | `+440/-0` | Prove one production fixed entry through allocation, moved-root materialization, TUI, cleanup, and recorder continuation. |

The full Subspace planning estimate becomes **25 named files,
+2,763/-213 = 2,976 changed LOC**. The real installed-skill First Officer drive is
one-off validation evidence and adds no Subspace harness file.

As in cycle 5, these numbers are planning and reconciliation evidence only. They are
not budgets, gates, tolerances, permission boundaries, pass/fail signals, or reset
conditions. Pre-edit declaration, path-scoped post-edit reconciliation, and reviewer
semantic materiality triage remain required; numeric variance alone does nothing.

### Acceptance and proof corrections

Cycle 6 AC-1 remains the value criterion, but its two behavioral layers are now
non-circular: the fixed-entry E2E proves production composition, and rq's own
pre-terminal live drive proves one real First Officer invocation plus recorder
continuation before rq can pass. ph is later sprint-wide composition only.

Cycle 6 AC-2 is corrected so the two-file count ends when selected provider entry
begins. `invocation-common` may allocate `ROOM/provider` before semantic
materialization validation. Every catchable failure after that allocation retains
only phase-appropriate provider evidence and removes the resolved-source child; no
provider root is reclassified as prepared authority.

Cycle 6 AC-3 and AC-4 remain unchanged: the exact indented positive has distinct full
JCS/raw pins, both are recomputed before display, arbitrary canonical filenames work,
canonical source tuples survive, and no identity helper or `association.json` enters
Result/inventory.

Cycle 6 AC-5 is corrected so structural lint proves only instruction closure and
forbidden-surface absence. The fixed-entry E2E proves product order; the pre-terminal
live drive proves agent order. A contract-lint test that passes or fails based on
semantic command order is itself a failing implementation.

The corrected proof order is:

1. land corrected s4 and pin the candidate Spacedock/Subspace tips;
2. red/green the exact indented JCS/raw control and unresolved Git-root boundary;
3. run focused authority, allocation, manifest-last, parser, summary, supervisor,
   validator, retention, and cleanup matrices;
4. run the production moved-root fixed-entry E2E through recorder continuation;
5. run both repositories' complete deterministic/race/shell/docs/format gates;
6. perform rq's one-off pre-terminal live First Officer drive on the exact candidates,
   append its native-event/on-disk evidence to this entity, and independently review
   it; and
7. only then allow rq's terminal validation/verdict path. ph remains blocked until rq
   and its other named prerequisites are terminal `PASSED`.

No live drive, provider invocation, product edit, s4 edit, gate/status/frontmatter
mutation, or ph change occurs in this ideation correction.

## Stage Report: ideation (cycle 7)

- DONE: rq owns a pre-terminal one-off live FO drive on landed s4 + rq candidate, using real installed skill and existing runtime substrate, proving exactly `/subspace:r gate <room>` then `gate record --room`; no harness.
  The drive now blocks rq terminal validation, uses candidate binaries/plugins in a fresh supported runtime, and retains native skill-event plus recorder/on-disk evidence in this entity.
- DONE: ph remains final composition evidence and keeps rq terminal prerequisite.
  ph is explicitly unchanged and downstream; it neither supplies nor waives rq's pre-terminal live evidence.
- DONE: Preserve invocation-common allocation of ROOM/provider; materialization must complete before resolved-source manifest publication and before TUI/model installation/display, not before provider allocation.
  `invocation-common` allocates/owns the provider root first; Spacedock completes verified payloads and publishes the manifest last before any TUI load or display.
- DONE: Remove semantic command-order claims from contractlint; keep only structural closure/absence, with fixed-entry E2E + rq live drive owning sequencing.
  Contract lint is limited to required/forbidden instruction surfaces; product order belongs to the real E2E and agent order to the one-off live drive.
- DONE: Preserve the JCS/raw fix and numeric non-authority.
  Cycle 6's distinct full pins and exact unequal control remain unchanged, while all revised file/LOC totals remain advisory reconciliation evidence.
- DONE: Run the repository-required deterministic gates against the unchanged product baseline.
  `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` are rerun before the path-scoped state commit; any failure will replace this DONE line with FAILED evidence.
- SKIPPED: Implement product code, run the live drive/provider, mutate gate/status frontmatter, touch s4/ph, or add a harness.
  Cycle 7 changes only this rq ideation entity; implementation and pre-terminal live evidence wait for landed corrected s4.

### Summary

Cycle 7 makes the proof non-circular and restores the existing Subspace package
ownership boundary. rq must prove its own real First Officer handoff before passing;
`invocation-common` allocates provider state before manifest-last materialization; and
contract lint returns to structural checks while behavioral sequencing stays with
exercised production and live paths.

## Stage Report: ideation (final AC evidence addendum)

- DONE: AC-1 has a concrete production-boundary proof and an independent value
  falsifier.
  Verified by: Spacedock candidate `66a937defa34725e2259e4d7f835870f6d680d1f`
  already supplies `TestPrepareCreatesOneTwoFileRecorderRoomForFolderAndFlatEntities`
  and `TestInspectAndResolveUseExactLocalGitObjectsAcrossMovedRoots`; the named
  implementation proof is `scripts/tests/subspace-r-git-root-provider-e2e.sh`, which
  must fail if either old sentinel is absent, the exact summary is not displayed, the
  inventory is not 2/2, or `gate record --room` does not close the same binding.
- DONE: AC-2 preserves Subspace allocation ownership, manifest-last publication, and
  the retained-package/no-retained-payload boundary.
  Verified by: Subspace tip `5466d601b7281a8d91715a1d03e190e7b5049c56`,
  blob `a3a49d2c9b1a63afe9f126f9090525ebe8fa505a`, shows
  `invocation_prepare_provider` allocating a private package under `umask 077` and
  `invocation_main` retaining that package after trusted delivery. The corrected
  implementation keeps that owner but fixes the allocation at mode-0700
  `ROOM/provider`; `internal/gates/materialize_test.go` must prove complete mode-0600
  payloads before manifest-last publication, while
  `scripts/tests/subspace-r-provider-retained-delivery-test.sh` must prove catchable
  paths retain provider evidence and remove `resolved-sources`.
- DONE: AC-3 binds the complete Briefing independently in JCS and raw-byte domains,
  with a real unequal positive control.
  Verified by: state commit `684d6603c985030c3c6031f4b1a1462c0f1cbfa1`
  contains blob `de790a44aa5b47ecdd606fc711a0ad8f9f20a2d7` (1,675
  raw bytes); its `rawSha256` is
  `sha256:c3b6d4d5ac8c766dcc56e08b57a41e207147d1319c61f066160e4e7d4bacfb1b`
  and its RFC 8785/JCS `jcsDigest` is
  `sha256:0782c65c06c7ee9378226b3a7ef88d92939a54c05d916fe3690cc7d99804278f`.
  `internal/gates/materialize_test.go` and
  `internal/reviewv1/resolved_sources_test.go` must accept that unequal pair and reject
  equality, swaps, the ambiguous `digest` member, and mutation of either full pin
  before model installation or display.
- DONE: AC-4 preserves the exact local Git-object materialization and refusal surface.
  Verified by: candidate `66a937defa34725e2259e4d7f835870f6d680d1f` blobs
  `674058e94ae7dc47c90d6b24aeb7f451e1c14965` (`source.go`) and
  `f94c5dbde4a00b4f2bb8a31cb49852806ecff97b` (`source_test.go`) resolve only
  `rev-parse --verify <full-commit>^{commit}` plus
  `cat-file blob <commit>:<path>`, then verify raw SHA-256.
  `TestInspectAndResolveUseExactLocalGitObjectsAcrossMovedRoots` proves moved-root
  reads, and `TestResolveRejectsNoncanonicalOrUnverifiableGitRootCoordinates` proves
  malformed coordinates, raw mismatch, and missing local objects refuse; the planned
  materializer composes this surface without fetch, clone, pull, worktree fallback, or
  source mutation.
- DONE: AC-5 separates the room-only entry, product sequence, live First Officer
  sequence, and structural lint proof owners.
  Verified by: present-gate blob
  `454c0522b6b0b0542725ec0d19104e68dcb26a00` contains exactly
  `/subspace:r gate <gate-room>` and subsequent
  `gate record <entity> --room <gate-room>`. The fixed-entry E2E owns observable
  allocation/materialization/manifest/TUI/validator/cleanup/record order; rq's
  pre-terminal live drive must cite one native installed-skill event followed by the
  recorder event before rq can pass. Current contractlint blob
  `38c85f51f8b4051efaba0957094957f307a25048` exposes the ordered `sequence(1, 7)`
  check that this implementation must remove; its replacement may assert only required
  room/record closure and forbidden-coordinate/fallback absence, never semantic order.
- DONE: The final ideation report now gives every scanned AC an external repository
  anchor and a named implementation-time falsifier.
  Verified by: an AC provenance scan can resolve AC-1 through AC-5 to the five entries
  above without treating planned implementation success as already executed evidence.
- SKIPPED: Change the design, checklist, product, gate/status/frontmatter, s4/ph, or
  run the provider/live lane.
  This addendum repairs report provenance only; all behavioral proofs remain assigned
  to their declared implementation or pre-terminal validation owner.

### Summary

Every AC now has reproducible ideation evidence outside the claim itself and a
concrete test or live-event condition that can fail later. The addendum preserves the
cycle 7 architecture and makes no implementation or workflow-state claim.

## Stage Report: implementation

- FAILED: Execute the approved ideation design exactly against the stack tip; every AC keeps its evidence discipline including AC-1's real Subspace presentation rendering 2/2 canonical Git-root sources
  The Spacedock half is complete at `579f12a00`; the Subspace half and AC-1's real presentation are not done. Two blockers, both raised to the FO and neither self-resolvable: the ~2,976 LOC Subspace surface is outside this dispatch's declared worktree, and AC-1's `gate record --room` continuation does not exist at the stack tip (see the material finding below).
- DONE: Declared surface within its gated tolerance; durable gate state stays free of copied selected-source payloads per AC-2
  Declared 9 Spacedock files / +1,082/-25; actual 9 files / +1,291/-0 (+19% additions, zero deletions because the one deletion surface was already satisfied upstream). One path deviation: the CLI route and its tests are `internal/cli/gate_materialize.go` and `internal/cli/gate_materialize_test.go` rather than inline in `cli.go`/`gate_test.go`, per the repo rule against a single large CLI file. AC-2 is proven by `TestMaterializeKeepsPreparedRoomAtTwoAuthoritativeFiles`, which fails if materialization adds a third room file, rewrites the canonical Briefing bytes, or writes `association.json`.
- DONE: Suites per the entity test plan; environmental failures reproduced on a clean control before attribution
  `gofmt -w ./cmd ./internal` left no residue, `go test ./...` and `go test ./... -race` both exit 0 across 20 packages with no data race. No environmental failure was observed, so no attribution was needed.
- SKIPPED: Remove the ordered `sequence(1, 7)` semantic-order claim from `internal/contractlint/fo_function_reference_invariant_test.go` (cycle 7 AC-5)
  Already satisfied at the stack tip: commit `723028f01` ("Retire banned prose-grep contract pins") removed it before this dispatch. Re-adding the replacement prose-grep closure assertion cycle 7 asks for would contradict that just-landed retirement, so it is held for a captain ruling rather than implemented.

### Material finding (worker proposal; held for FO/captain authorization)

AC-1 and AC-5 depend on a `gate record --room` recorder continuation and a
`/subspace:r gate <room>` instruction surface that no longer exist. Ideation was written
against the corrected-s4 ensign-branch candidate (`e328ecc6`/`66a937de`); what landed on
main is the provider-neutral cut, which removed both.

1. Released user and normal workflow: a First Officer presenting a gate through Subspace on the 0.27 pre-ship line.
2. Observable harm: the chain stops at the published manifest. No installed instruction routes an agent to the materializer, and no provider Result can be recorded.
3. Affected value AC: `value-ac[AC-1]` — the room-only path must present both Git-root sources and continue through `gate record --room`.
4. Trigger evidence at stack tip `d1d8f745`: `internal/cli/gate_test.go:692` asserts exit 2 with stderr exactly `Error: unknown gate flag: --room`; `docs/specs/gate-resolution-frontmatter-contract.md` lists `gate record --room` under "Explicitly outside v1"; `grep -rn subspace skills/` returns zero matches; commit `cb267d09e` removed `/subspace:r gate` from present-gate.

Proposal — materiality: Material. Ownership: not rq's; the remedy reverses a landed
captain-visible scope ruling. Disposition: Needs decision. No candidate edit was made to
`present-gate`, `gate record`, or that test.

### Summary

`gate materialize --room ROOM` is implemented, documented, and green: room-only entry
with every other coordinate derived by fixed code, arbitrary canonical Briefing locator,
separate `jcsDigest`/`rawSha256` domains pinned by the exact 1,675-byte s4 control, closed
manifest published last by rename, mode-0600 payloads, and provider-root ownership left to
Subspace. Falsifiability was checked by mutation rather than asserted: swapping the two
digest domains, restoring the ambiguous `digest` spelling, and dropping the mode-0700 check
each turn the suite red. The value chain cannot close inside this dispatch — the Subspace
repository half is out of worktree scope and AC-1's recorder continuation was cut from v1.
