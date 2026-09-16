---
title: Use semantic branch and worker names
status: validation
source: Captain request 2026-09-15; GitHub issue 624
started: 2026-09-15T18:27:49Z
completed:
verdict:
score: 0.8
worktree: .worktrees/spacedock-ensign-semantic-dispatch-names
issue: spacedock-dev/spacedock#624
pr: "#802"
mod-block: merge:pr-merge
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
        - id: gate:6es505tn1zz2597hetvnqn7y:validation
          stage: validation
          attempts:
            - id: gate-attempt:6es505tn1zz2597hetvnqn7y-validation-1
              briefing:
                id: briefing:6es505tn1zz2597hetvnqn7y:validation:attempt-1:revision-1
                digest: sha256:d6cc8df5fdeec4debb3be23db3579531d776a172da56b35e292d658489b6bb12
                room-ref: '@review/validation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:6es505tn1zz2597hetvnqn7y:validation:1
                briefing: briefing:6es505tn1zz2597hetvnqn7y:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-09-16T02:13:00.181176Z"
                decision: approve
                reason: Captain approved validation in Subspace resolution:binding-1789524348826001000; accepts direction and presented evidence.
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:6es505tn1zz2597hetvnqn7y-validation-2
              briefing:
                id: briefing:6es505tn1zz2597hetvnqn7y:validation:attempt-2:revision-1
                digest: sha256:87f64cf46c42a95fc6be16f609be0b7d3c58c9a4aa3d5838595e11e94dc93615
                room-ref: '@review/validation/briefing-2'
              resolution:
                type: Resolution
                id: resolution:spacedock:6es505tn1zz2597hetvnqn7y:validation:2
                briefing: briefing:6es505tn1zz2597hetvnqn7y:validation:attempt-2:revision-1
                by: person:captain
                at: "2026-09-16T20:44:39.48535Z"
                decision: approve
                reason: Captain binding resolution binding-1789591360273000000 approves all three named corrected validation snapshots in /tmp/stack-corrections-review.md for stack publication and final-tip CI. Restack preserves all patches and tree. No merge authority or live failure waiver.
              application:
                target-stage: done
                state: pending
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


## Stage Report: implementation

- DONE: Implement semantic branch and worker names while preserving registered active branches and legacy lifecycle identity.
  Code commits `23b8732dc` and `1b3b70508`: fresh slug branches and readable SHA-256-shortened workers; registered custom/legacy branches persist, ambiguous identities refuse, mixed-generation cohorts retain original handles.
- DONE: Prove collision, shortening, retry/recovery and actual native dispatch behavior through existing proof owners.
  Focused dispatch/CLI/feedback/registry/skill tests pass; literal legacy identities, stage-boundary/cross-generation collisions, equal-cycle ties, branch collisions, detached paths and sync-failure retries fail independently if their guards regress.
- DONE: Commit minimal candidate and complete checks/report within approved scope for independent validation.
  Scheduling base `4ce49f1ea` to `1b3b70508`: exactly 46 files, +615/-185 (+430 net), within the amended 46-file/450-net ceiling; captain resolution `resolution:binding-1789516131152632000` (artifact SHA-256 `ab235343c7cb33a1fd870b527728f761ecd606c194e3245910cb8a659d17c932`).
- DONE: Exercise native Codex semantic dispatch and same-handle advance.
  `TestLiveSemanticNamesCodex` passed in 109.67s: spawn `ci_duration_hints_ideation`, follow-up `/root/ci_duration_hints_ideation`, body-only fixture commits `5c9057416c344afc9eb904d01ea5554d7342303a` and `21354d633b08a59979d36dce5c891e9e062f5ffd`, unchanged frontmatter and clean state.
- DONE: Preserve falsifiable native and instruction evidence.
  `/tmp/semantic-native-evidence-v2/codex-shared-scenarios/semantic-names-codex/` retains native lifecycle, entity and state log; offline replay rejects altered prefix/handle with markers retained. Final live-tag compilation passes; replay does not claim another live run. Skill checks prove canonical envelope/slot compatibility, not live Claude compliance.
- DONE: Run detached adversarial checks and formatting.
  Detached edits restoring the prefix, accepting ambiguous ownership, reconstructing a branch, restoring the recovery template, or changing the legacy cap each fail intended assertions; `/tmp/semantic-detached-audit.log`, `/tmp/semantic-legacy-cap-audit.log`, `/tmp/semantic-native-negative-controls.log`. `gofmt -w ./cmd ./internal` ran; unrelated baseline formatting was kept outside this diff.
- SKIPPED: Code push and CI execution.
  Explicit FO restriction; both code commits remain local. The targeted proof registry entry uses the category authorized by #792 and must be combined with that policy during restacking.

### Executed-checks limitation (resolved disposition)

Both `go test ./...` and `go test ./... -race` exit 1 solely at unchanged `TestCodexResolveManifestAgainstInstalledHost`, `internal/cli/codex_resolve_test.go:44`: resolver returns `/Users/clkao/.codex/plugins/cache/spacedock-local/spacedock/0.28.0-pre0/.codex-plugin/plugin.json` although `spacedock@spacedock` is not installed. Disposition: FO-authorized DECLINE, resolved; repair of this unchanged baseline is outside this task. The suites are not wholly green. Every other package passes and race emits no data-race diagnostic. This disposition does not waive any new regression. Logs: `/tmp/semantic-full-normal.log`, `/tmp/semantic-full-race.log`.

### Summary

Semantic public names now preserve task meaning while existing registered branches and worker handles retain their identities. Native Codex spawn/reuse, conservative collision behavior and literal historical-name compatibility are exercised, with the final candidate inside the amended scope. Required broad checks were executed and are not wholly green: the unchanged resolver baseline failure has a resolved FO-authorized DECLINE disposition.

## Stage Report: validation

- DONE: Assess every AC against exact committed naming candidate, retained raw native spawn/reuse evidence and legacy compatibility; reuse owned green checks.
  Reviewed `1b3b70508` against scheduling base `4ce49f1ea`; independently confirmed 46 files, +615/-185 (+430 net), within the captain-amended 46/450 ceiling.
- DONE: Perform an independent detached semantic audit of identity, collision, shortening and retry/recovery claims, reporting proposed findings before candidate changes.
  Detached audit passed in 4.179s; artifacts under `validation-audit/` include runnable test source, output and failing prefix/ambiguous-owner mutation controls. Candidate bytes and HEAD stayed unchanged; no new material finding.
- DONE: Report PASSED or REJECTED with per-AC evidence, exact scope and resolver baseline limitation; preserve pending stack policy dependency and no code push or CI.
  Recommendation: PASSED for the naming candidate, subject to the recorded delivery dependencies; broad-suite baseline remains an explicit failure, not a green claim.
- DONE: AC-1 — Fresh public names expose the task meaning with less irrelevant text.
  Independent `TestValidationAuditExactPublic` executes ideation dispatch/stamp and checks exact 26-character worker, actual 17-character Git branch and assignment branch agreement; restoring the prefix fails.
- DONE: AC-2 — Long and ambiguous names remain readable, bounded, and distinguishable.
  Independent matrix covers short/long/single-character slugs, three stages, five suffixes, invalid/Unicode input, over-budget suffix and 1,000 entities; existing semantic tests cover literal legacy caps, sibling digests, collision seam and no-mutation stamp refusal. Accepting ambiguous ownership fails the negative control.
- DONE: AC-3 — Legacy and new workers retain correct lifecycle ownership.
  Reused owned advance/reconcile/ownership/stamp checks plus detached semantic cohort tests: registered legacy/custom paths survive; literal historical caps resolve; mixed-generation ties retain handles. Recovery smoke proves canonical output and instruction slots only, not live Claude execution; stale template/digest controls remain independently falsifiable.
- DONE: AC-4 — A native Codex dispatch and reuse complete under the semantic namespace.
  Raw native call IDs independently link one `ci_duration_hints_ideation` spawn to `/root/ci_duration_hints_ideation` and one followup to that same handle. Retained entity/log show two marker commits touching only the fixture entity; owned live run passed in 109.67s. Encrypted message bodies limit independent verbatim-prompt inspection.
- DONE: Reuse required broad checks and formatting evidence without competing for the broad-suite slot.
  Normal/race logs fail solely at unchanged `TestCodexResolveManifestAgainstInstalledHost`, line 44 (FO-authorized DECLINE); all other packages pass, no race diagnostic. Implementation ran gofmt; detached `git diff --check` passes.
- SKIPPED: Code push and CI execution.
  Explicit restriction. Targeted registry categorization still requires #792 policy during future restacking; scheduling remains the base PR. No normal-CI claim.

### Summary

PASSED: exact public names, conservative ownership, legacy compatibility and native same-handle reuse satisfy AC-1 through AC-4 on the committed candidate. Independent detached checks found no new material or deferred finding; the known resolver baseline failure and future policy-stack dependency remain explicit limitations. Native call identity is independently attributable, while encrypted prompt contents rely on the owned live proof rather than independent byte-level inspection.


## Stage Report: implementation

- DONE: Reconcile naming branch onto the same-stage CI promotion while preserving lower coverage and the strict checker.
  Tested head `844fad458` descends from `bf64ebeca`; FO clean-restacked it onto `342edfa89` as current head `7535aad706`. Artifact `8d7d0cc5e` records the preserved six-variant coverage and explicit experiments.
- DONE: Promote native naming into the Codex lane, move deterministic identity proof offline, and fix exact branch membership with rejecting negative controls.
  Artifact `8d7d0cc5e` records default-suite discovery/execution, missing Codex scheduling-row rejection, actual handoff-grader red-to-green replay, and wrong/missing/extra branch negatives; actual scope is 55 files and +463 net lines.
- DONE: Run applicable local checks and record exact candidate, scope and evidence, leaving native verification to final-tip CI.
  On `844fad458`, focused checks and live-tag compilation passed; normal/race each exited 1 solely at the verified resolver baseline, with no data-race diagnostic. Artifact `8d7d0cc5e` retains logs; rebased-head full checks and final-tip CI are not claimed.

### Summary

The tested implementation is `844fad458` on parent `bf64ebeca`. FO subsequently clean-restacked it onto coverage-checker repair `342edfa89`; the current candidate is `7535aad706e1d0ae88cad8bd45f5ebf4eadf96b8`. No new full-suite result is claimed for that rebased head. Actual naming scope is 55 files, +817/-354 (+463 net), reflecting the captain's amended coverage requirement and authorized branch-inventory correction.

The existing Codex scheduled lane now selects the native semantic-name/same-handle proof. The stamped-owner identity check runs in the default offline suite using a freshly built checkout binary. The earlier targeted-only exemption is superseded by the captain's coverage correction; the three explicit experiments remain unchanged. Legacy identity claims retain their offline proof, and earlier raw native evidence is preserved.

The actual handoff grader rejected `["conflict-owner" "main"]` before correction and passed the deterministic full-fixture replay afterward. Exact membership retains the two-branch limit and rejects wrong, missing and extra branches. Default test discovery/execution, legacy and collision controls, live-tag compilation, and registry checks pass on the tested head. Removing the Codex scheduling row makes its checker fail. Detailed evidence is committed in `8d7d0cc5e`, [routine-coverage-and-branch-inventory.md](validation-audit/routine-coverage-and-branch-inventory.md).

On `844fad458`, both `go test ./...` and `go test ./... -race` exited 1 solely at `TestCodexResolveManifestAgainstInstalledHost`, `internal/cli/codex_resolve_test.go:44`: `spacedock@spacedock` is not installed, but the resolver returned `/Users/clkao/.codex/plugins/cache/spacedock-local/spacedock/0.28.0-pre0/.codex-plugin/plugin.json`. All other packages passed; no data-race diagnostic appeared. Logs: `/tmp/semantic-promotion-full-normal.log` and `/tmp/semantic-promotion-full-race.log`. Required checks were executed and are not wholly green. The existing FO-authorized baseline DECLINE is resolved and does not waive new failures. Independent validation and final-tip CI remain pending; no local model run or new CI result is claimed.

## Stage Report: validation (cycle 2)

- DONE: Verify native naming is actually selected by the Codex lane and deterministic owner identity runs offline, with explicit exemptions preserved.
  On `7535aad706e1d0ae88cad8bd45f5ebf4eadf96b8`, the actual `342edfa89` registry checker passes with executed scheduler variants; default-tag `TestConflictOwnerStampedIdentity` passes. Three explicit experiments and original native identity assertions remain unchanged; `validation-audit/cycle2/` retains evidence.
- DONE: Independently reproduce exact branch membership behavior and rejecting controls using the real handoff grader; preserve semantic and legacy identity assertions.
  Complete grader replay passes semantic and legacy owner branches and rejects missing/wrong/extra inventories; restoring lexical ordering fails on `["conflict-owner" "main"]`. Current semantic/Codex checks pass; prior AC-1/2 matrix and AC-4 raw native evidence remain preserved, with per-AC mapping in the artifact.
- DONE: Assess rebased candidate, scope and exact local evidence against captain requirements, with final combined and host CI explicitly pending.
  Local recommendation PASSED, no new finding: exact scope 55 files, +817/-354 (+463 net). Prior full normal/race results belong to `844fad458`/`bf64ebeca` and fail solely at the resolved resolver baseline; no full run is claimed for this rebased head. Final combined suites and both-host live CI remain pending; no model run, code change, push or CI action occurred.

### Summary

The rebased naming candidate satisfies the amended routine-selection requirement and preserves strict ownership checks; the independent actual-grader replay confirms that the branch-order correction fixes the observed false failure without relaxing branch membership. Local validation is PASSED, while final combined-tip checks and native host CI remain explicitly pending. Earlier native spawn/reuse evidence is retained as historical evidence and is not presented as a new run.
