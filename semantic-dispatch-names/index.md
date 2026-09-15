---
title: Use semantic branch and worker names
status: implementation
source: Captain request 2026-09-15; GitHub issue 624
started: 2026-09-15T18:27:49Z
completed:
verdict:
score: 0.8
worktree:
issue: spacedock-dev/spacedock#624
pr:
mod-block:
id: 6es505tn1zz2597hetvnqn7y
gates:
    version: 1
    records:
        - id: gate:6es505tn1zz2597hetvnqn7y:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:6es505tn1zz2597hetvnqn7y-ideation-1
              briefing:
                id: briefing:6es505tn1zz2597hetvnqn7y:ideation:attempt-1:revision-1
                digest: sha256:9b52132f4ba2803e8e48052ffbdd3f114fb56bd5e42e3db6d6c394d5699bd72f
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:6es505tn1zz2597hetvnqn7y:ideation:1
                briefing: briefing:6es505tn1zz2597hetvnqn7y:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-09-15T23:04:13.349623Z"
                decision: approve
                reason: Captain binding resolution binding-1789513304916529000 explicitly approved naming implementation, estimate340/24 cap450/27, preserving legacy identity and required runtime proof; reviewed document sha256:3e4c7f21b5d53b1abfcc5796d19b7bbe7d90ff2f5c5398415f1e11b245fe9c81.
              application:
                target-stage: implementation
                state: consumed
---

Use short, readable task names for public branches and dispatched workers. Keep descriptive detail in titles.

## Problem

Issue #624 and the captain's request concern names people read in PR lists and worker namespaces. At origin/main `438053493838dc70c9478b3d991309d566783e85`, `thing` produces branch `spacedock-ensign/thing` and worker `spacedock-ensign-thing-implementation`. Long names can replace all meaningful words with an entity ID. Removing only the dispatch prefix also makes reconciliation ignore the worker and leaves stamp/prompt branch assumptions inconsistent.

## Proposed approach

1. New branches use the existing entity slug verbatim: `ci-duration-hints`. New workers use `ci-duration-hints-ideation`; the existing Codex adapter maps this to `ci_duration_hints_ideation`. Agent type remains separate dispatch metadata. No workflow naming option, template engine, alias table, new stored identity, or bulk rename.
2. Keep existing worktree directory conventions to avoid unrelated path churn. When `worktree:` already names a registered Git worktree, preserve both the recorded path (including a custom or absolute path unlike the generated default) and its actual branch, and use both in the assignment. Resolve through Git worktree registration, not a path basename. Detached or unregistered existing paths refuse with a diagnostic; never silently claim a different branch. A fresh branch collision refuses before stamps, commits, or worktree creation and names the occupied ref; the operator can choose a different task slug. Do not auto-checkout an unrelated same-name branch.
3. Use one small naming helper shared by build and reconciliation. Names fit a 56-character base budget, reserving eight characters for existing retry/cycle suffixes under the 64-character host ceiling. This covers `-retry` and cycle suffixes through eight characters, not unbounded cycle numbers; reject a longer requested suffix before spawn. Short names pass unchanged. Overlong names keep a readable slug head plus eight lowercase hex characters from SHA-256 of the full slug, then the complete stage suffix. This works uniformly for sd-b32, sequential, and slug identities without substituting an opaque ID for all text. Require at least eight readable slug characters in a shortened name; reject a stage that leaves insufficient room. Names and stages must satisfy the existing lowercase kebab-case grammar; do not introduce lossy arbitrary normalization. Codex's existing hyphen-to-underscore mapping is injective on that grammar.
4. Reconciliation resolves a worker against actual active/archived entities and known stages: enumerate recomputed semantic and legacy candidates, including sd-b32 ID caps, sequential ID caps, and old slug truncation. Compare unsuffixed entity/stage candidates before interpreting retry/cycle/numeric suffixes; collect all valid interpretations and decline ambiguity rather than assigning arbitrary precedence. A real stage `review-2` must not silently become `review` cycle 2. Resolve identity before grouping superseded workers, so legacy and semantic names for the same entity/stage form one cohort. Never classify an ambiguous candidate or use a readable prefix alone as identity. Generation checks every interpretation generation and entity/stage boundary among the current workflow's active/archived entities and refuses before mutation. Examples: semantic slug `spacedock-ensign-thing` aliases legacy slug `thing`; `alpha-review` at `done` aliases `alpha` at `review-done`. Digest and cross-generation collisions fail safely rather than routing to a sibling. Preserve every original roster handle while grouping; equal-cycle ties stay unclassified for superseded shutdown, never overwrite a handle or select an arbitrary loser. Existing slug uniqueness and runtime parent namespace remain the boundaries; no global cross-repository allocator.
5. Reuse addresses the existing host handle/task path, including legacy names; `--advance` supplies a new assignment rather than renaming that handle. Fresh replacement workers get the new naming form. Both fresh-retry/recovery instruction owners consume the helper-emitted canonical base name and append only the bounded existing `-retry` or `-cycleN` suffix; they never reproduce truncation or hash logic in prose. When build is unavailable, selected named-mode recovery may use a retained, validated semantic envelope for the exact same entity/stage; without that envelope, hold that entity and report the unavailable name rather than inventing a name or switching transport. Legacy handles remain usable for reuse, not as an algorithm for fresh names. Keep completion validation, shutdown authority, state formats, retry count and ledger ownership unchanged.

The hash suffix serves AC-2: plain truncation conflates long sibling slugs; retaining only ID tokens loses AC-1's readable meaning. Candidate recomputation serves AC-2/3 and needs no identity ledger because the workflow already holds the entity slug and stage set. Git registration serves AC-3: deriving a branch from today's naming rule would rename or misdescribe existing PR heads. Collision refusal is sufficient; automatic alternative branches and a configurable naming policy add no value needed by this task.

## Delivery dependency

Scheduling is the base PR of the next stack; semantic naming is the PR above it. Implementation starts in an isolated worktree based on that scheduling PR head (or its merged main successor), records the exact base SHA, and keeps its diff scoped to naming. This delivery order does not authorize changes to active branch names, worker identities, or remote metadata. The pinned main SHA below remains the historical spike baseline.

## Risk evidence

A reproducible throwaway spike is committed beside this entity: `spike/run.py`, `spike/semantic_spike_test.go.txt`, and `spike/prototype.patch`. Run `python3 docs/dev/.spacedock-state/semantic-dispatch-names/spike/run.py` from a source checkout containing the base commit. Requirements: Python 3 with tar extraction filters, Git, Go. It archives the pinned base to a temporary directory, uses fixture-owned temporary Git repositories and team configuration, and never edits active branches or runtime identities.

Observed: the unchanged base passes legacy dispatch/stamp, repeat stamp, and terminal reconciliation. The independent semantic expectation fails with `spacedock-ensign-thing-implementation` instead of `thing-implementation`. A four-site prototype then passes: the registered branch is `thing`, the emitted worker is `thing-implementation`, repeat stamp succeeds, and both semantic and legacy roster names yield exactly one `lingering` record for `thing`. This exercises real command routing, on-disk stamps/commits/worktree creation, and the production reconcile entry point with a fixture roster; it is not a native runtime claim.

The prototype establishes the risky prefix-free dispatch/reconcile path, not the finished implementation. Long-name hashing, ambiguous matches, mixed-name superseded cohorts, and preserving a pre-existing legacy branch remain first-test obligations below. The prototype intentionally leaves the old name-cap logic intact and must not be shipped verbatim.

## Out of scope

Renaming active tasks/workers/public PR heads, changing stable entity IDs or stored frontmatter schemas, new runtime providers, scheduler changes, automatic PR creation, automatic branch disambiguation, global naming registries, unrelated contract cleanup, and remote metadata changes.

## Expected surface and tolerance

Estimate net LOC change: +340, across 24 files. Expected insertions: 560; deletions: 220. Forecast band: net +150 through +450; a correct smaller change is acceptable and must not be padded. Hard tolerance: at most +450 net LOC and 27 changed source/documentation files relative to the recorded scheduling-stack base. Exceeding either maximum requires returning to the gate. Ideation state artifacts do not count toward the implementation diff.

Expected source owners: `internal/dispatch/build.go`, `stamp.go`, `reconcile.go`, and new `names.go` (all under `internal/dispatch`). Expected proof owners there: new `names_test.go`; existing `build_namecap_test.go`, `build_stamp_test.go`, `build_codex_host_test.go`, `build_pi_host_test.go`, `build_input_mode_test.go`, `build_errors_test.go`, `build_parity_test.go`, `build_hazards_test.go`, `reconcile_decompose_test.go`, `reconcile_test.go`, `reconcile_de_safety_test.go`, `build_advance_test.go`, and `build_advance_contract_parity_test.go`. One targeted runtime proof file: `internal/ensigncycle/semantic_names_live_test.go`. Documentation: `docs/site/concepts/workflows-and-entities.md` and `docs/site/reference/command-reference.md`. The review adds `skills/first-officer/references/claude-fo-dispatch.md`, `skills/fo-dispatch-recovery/SKILL.md`, and the existing smoke owner `skills/integration/dispatch_test.go`, consuming three of the original six contingency files. The remaining three-file allowance covers genuinely affected fixture/golden owners discovered by the full suite, not new mechanisms.

Allowed semantic changes: fresh branch/worker output names, deterministic readable shortening, and conservative naming/branch-collision refusal, and holding selected named-mode manual recovery when no validated semantic envelope is available; reconciliation recognizes both naming generations. Command flags/grammar, JSON keys, entity IDs, state authority, runtime dispatch tool bindings, and existing handle/branch identities do not change. Worktree directory convention remains unchanged.

## Acceptance criteria

**AC-1 — Fresh public names expose the task meaning with less irrelevant text.** For independently specified slug `ci-duration-hints` and stage `ideation`, dispatch emits worker `ci-duration-hints-ideation` (26 characters versus the baseline's 43) and creates branch `ci-duration-hints` (17 versus 34); neither includes agent type. Registered worktree, assignment branch text, and envelope agree. Proof owner: build/namecap and stamp integration tests comparing exact expected strings and actual Git registration. Restoring the prefix fails the value measurement.

**AC-2 — Long and ambiguous names remain readable, bounded, and distinguishable.** Every supported id style retains readable slug words, a complete stage suffix, and cycle headroom; long same-prefix siblings produce different names and reconcile to the correct entity. Invalid input, insufficient stage budget, an equal generated candidate, and a pre-existing fresh branch produce refusal with no state/artifact/worktree mutation. Proof owner: `names_test.go`, existing namecap/error/stamp tests. Removing the digest, accepting an ambiguous candidate, or moving validation after stamping each has a distinct failing case. Hash collision behavior uses crafted candidate data at the resolver seam; do not require brute-forcing SHA-256. Fixtures include cross-generation aliases, `alpha-review`/`review-done` stage-boundary ambiguity, a real `review-2` stage, and over-budget suffix refusal. Unsuffixed matching must not erase genuine ambiguity from other interpretations.

**AC-3 — Legacy and new workers retain correct lifecycle ownership.** Existing registered worktrees keep their recorded paths and original branches through build/stamp and advance, including a custom recorded path unlike the generated default. Legacy prefixed, sd-b32-capped, sequential-capped, slug-truncated, semantic, and semantic-shortened workers resolve uniquely or remain conservatively unclassified; mixed-generation retry/cycle cohorts produce the correct superseded/lingering/owned outcomes. Equal-cycle mixed-generation cohorts retain both original handles and cause no arbitrary superseded shutdown. Reuse keeps the existing host handle and cleanup targets the recorded worktree. Fresh retries consume the canonical generated name and unavailable-build recovery follows the bounded cached-envelope-or-hold rule. Proof owner: stamp, advance, reconcile decomposition, five-class, and ownership-safety fixtures. Rejecting legacy input, grouping raw name tokens before resolving identity, or reconstructing existing branches from new defaults fails independent cases.

**AC-4 — A native Codex dispatch and reuse complete under the semantic namespace.** One isolated temporary split-root journey uses a generated `ci-duration-hints-ideation` dispatch, spawns `ci_duration_hints_ideation`, then sends one `--advance` assignment to the same returned handle. Both turns append distinct marker-bearing stage reports and commit only the fixture entity body; no frontmatter changes by workers. Captured native tool arguments show the same handle on reuse and no agent prefix. Durable state records both new commits, exact markers, and a clean entity path. Proof owner: targeted live harness using existing Codex auth/home/process helpers; envelope-only checks cannot satisfy this AC.

## Test plan

Implementation writes focused tests before changing each mechanism. Extend existing owners rather than duplicate naming/branch assertions in static skill tests. Before editing FO instruction text, extend `skills/integration/dispatch_test.go` with skill smoke cases that inspect the shipped retry/recovery naming slots and exercise their canonical-name source through real dispatch output: fresh retry, cached same-entity/stage recovery, and no-envelope hold. Restoring the prefixed manual template or its ID/truncation recipe must fail the instruction-contract check; changing binary naming while keeping stale instruction expectations must fail the output check. These checks establish instruction and envelope compatibility, not live Claude compliance. Existing `internal/ensigncycle/dispatch_recovery_assert_unit_test.go` offline retry-bound controls remain regression coverage without adding a Claude live lane. The spike supplies the first integration fixture. Unit naming cases take milliseconds; dispatch/stamp/reconcile Git fixtures should take under a minute together. Existing advance tests own assignment transport; only extend their name/handle compatibility assertions.

The targeted native proof is `go test -tags live ./internal/ensigncycle -run '^TestLiveSemanticNamesCodex$' -count=1 -v`: one parent session, one fresh worker, one follow-up, maximum ten minutes. Reuse the existing isolated Codex home/auth and runner helpers, add no runtime driver or provider. The parent uses the built candidate binary and generated pointer; it records exact spawn/follow-up arguments and inspects the fixture state Git log. The worker only writes and commits the fixture entity. A name override back to the legacy prefix or follow-up to a new handle must fail evidence validation even if marker text exists. This targeted proof is scheduled for implementation/validation, not claimed by the offline ideation spike. No remote push, CI matrix, or modification of live identities is needed.

After implementation run `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`, plus the workflow's applicable detached audit. Compare actual net LOC and file count to the approved bounds. Preserve existing legacy-input fixtures; adjust only expectations whose fresh-generation contract changes.

## Proposed retry/recovery instruction diff

In `skills/first-officer/references/claude-fo-dispatch.md`, replace the retry naming clause:

```diff
- carrying a distinct `-retry` suffix on the `{worker_key}-{slug}-{stage}` name
+ carrying a distinct `-retry` suffix on the canonical base `name` emitted by dispatch build for this entity/stage; retain existing handles for nudges
```

In `skills/fo-dispatch-recovery/SKILL.md`, replace the named manual template slot and old cap recipe:

```diff
-    name="{worker_key}-{slug}-{stage}",  // old ID-first or slug-truncation overflow recipe
+    name="{canonical_name}",  // helper-emitted semantic base name for this exact entity/stage
```

Add before the selected named-mode template:

```diff
+Use the canonical semantic name from a successful dispatch envelope for this entity/stage. When build is unavailable, use a retained validated envelope for that same assignment; if none exists, hold this entity and report the missing canonical name. Do not derive a name from the old worker handle, duplicate shortening logic, or change dispatch mode. A replacement appends only its existing `-retry` or `-cycleN` suffix, at most eight characters; refuse over-budget or already-occupied final names before spawn.
```

Replace the dead-ensign suffix instruction with:

```diff
- Fresh-dispatch under a `-cycleN` suffix when replacing a zombie ensign.
+ Fresh-dispatch using the canonical name above plus a bounded `-cycleN` suffix when replacing a zombie ensign; preserve the old handle only for bookkeeping.
```

These instruction changes serve AC-1/3. Consuming an existing helper result is simpler than a second prose naming algorithm; holding only the affected named assignment is the bounded fallback when that result is unavailable.

## Proposed documentation diff

In `docs/site/concepts/workflows-and-entities.md`, after the paragraph beginning “An entity lives as a flat file”:

```diff
+Choose a short, readable slug, such as `ci-duration-hints`; put descriptive detail in the title.
+New code branches use the slug. New worker names use `<slug>-<stage>`; Codex displays underscores in place of hyphens.
+Long worker names retain a readable slug prefix and a deterministic suffix. Existing branches and worker handles keep their names.
```

In `docs/site/reference/command-reference.md`, immediately after the command summary table:

```diff
+`dispatch build` names fresh workers `<slug>-<stage>`. Names reserve room for retry/cycle suffixes;
+long slugs keep a readable prefix plus a deterministic eight-character suffix.
+`dispatch build --stamp` creates new code branches named `<slug>` and preserves the branch of an existing registered worktree.
+An occupied new branch or ambiguous generated worker name is refused before mutation. Choose a distinct task slug to resolve a fresh-name collision.
+Legacy worker names remain valid for reconciliation and reuse.
```

### Feedback Cycles


## Stage Report: ideation

- DONE: Prove readable branch and worker naming through real dispatch and reconciliation for both new and legacy names.
  `spike/run.py` exits 0: legacy baseline passes, semantic expectation fails unchanged baseline, prototype creates branch `thing` and reconciles both worker generations to terminal `thing`; restoring the prefix breaks the exact-name assertion.
- DONE: Identify the smallest compatibility design, collision handling, and exact file/net-line tolerance.
  Proposed approach and expected surface: no naming option/ledger; preserve registered branches and existing handles, refuse collisions; +300 net across 21 files, tolerance +150..450 net and maximum 27 files.
- DONE: Record acceptance evidence owners and a bounded targeted runtime proof without changing active branches or identities.
  AC-1 through AC-4 name existing proof owners and independent failures; one Codex fresh/reuse journey capped at ten minutes is planned. All spike mutations used temporary Git fixtures.
- SKIPPED: Full normal/race suites and source formatting.
  Ideation changes only the state design and retained throwaway evidence; production implementation and its required full-suite/format checks remain later-stage obligations.
- SKIPPED: State branch remote push.
  First officer explicitly instructed no remote changes; state work is committed locally with entity-scoped paths.

### Summary

Prefix-free names passed a reproducible dispatch/stamp/reconcile spike while legacy worker resolution remained usable. The design keeps readable shortening and compatibility in the existing naming/Git lifecycle boundaries, with concrete documentation wording, acceptance proofs, and net-line/file limits ready for review.


## Stage Report: ideation (cycle 2)

- DONE: Prove readable branch and worker naming through real dispatch and reconciliation for both new and legacy names.
  Prior committed `spike/run.py` evidence remains applicable; this correction changes the reviewed design only and does not claim new runtime or product proof.
- DONE: Identify the smallest compatibility design, collision handling, and exact file/net-line tolerance.
  Independent review caveats are explicit fixtures: cross-generation/stage-boundary aliases, stage `review-2`, legacy sequential caps and truncation, custom recorded paths, pre-mutation refusals, and equal-cycle handle ties; estimate +340 net (560 insert/220 delete), 24 files, unchanged +450/27 ceilings.
- DONE: Record acceptance evidence owners and a bounded targeted runtime proof without changing active branches or identities.
  Both fresh-retry/recovery instruction owners now have concrete wording and existing `skills/integration/dispatch_test.go` smoke ownership before text edits; the single ten-minute Codex journey remains unchanged. Scheduling is explicitly the base PR and naming its successor.
- SKIPPED: Product implementation and new test execution.
  FO FIX is limited to the ideation body; no source, skill, active identity, or remote changed. Existing spike evidence was not rerun for this prose-only amendment.

### Summary

The targeted review gap is covered without adding naming configuration, a ledger, a new runtime driver, or another live lane. Recovery consumes a canonical semantic envelope or holds the affected named assignment when none is available; the design and ACs now make all compatibility edge cases explicit.
