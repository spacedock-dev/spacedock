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
| `internal/cli/gate.go` | `+55/-5` | Route and document `gate materialize`. |
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
