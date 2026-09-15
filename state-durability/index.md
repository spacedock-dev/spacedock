---
title: Make state transitions durable
status: ideation
source: Captain-approved combined durability work, 2026-09-15
started: 2026-09-15T20:07:01Z
completed:
verdict:
score: 0.95
worktree:
issue: spacedock-dev/spacedock#689
pr:
mod-block:
id: 3tzfv0rrctdb066xt942zjbd
gates:
    version: 1
    records:
        - id: gate:3tzfv0rrctdb066xt942zjbd:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:3tzfv0rrctdb066xt942zjbd-ideation-1
              briefing:
                id: briefing:3tzfv0rrctdb066xt942zjbd:ideation:attempt-1:revision-1
                digest: sha256:fb08cb93e8003d48d0b6521799735e48cba11cb33e6b547559d938e83ade3e8d
                room-ref: '@review/ideation/briefing-1'
---

Prevent false completion and silent loss of workflow state. One task covers GitHub issues #689, #790 and #630 with separate behavioral acceptance criteria.

## Problem

On main `438053493838dc70c9478b3d991309d566783e85`, a passed local merge guard without a registered hook consumes terminal approval and archives before its printed manual merge instruction runs. A later conflict leaves false completion. Explicit retirement moves state without committing it. Missing split-root storage returns the same empty JSON as initialized empty storage and accepts filing into an ignored directory.

The captain chose one combined task. Keep retirement distinct from delivered completion. Stack this work above the scheduling improvement at the bottom of the next PR stack, after in-flight changes are ready; ideation does not push code or run CI.

## Risk evidence

Sources: [#689](https://github.com/spacedock-dev/spacedock/issues/689), [#790](https://github.com/spacedock-dev/spacedock/issues/790), [#630](https://github.com/spacedock-dev/spacedock/issues/630). The reproducible Python standard-library harness is `artifacts/ideation/reproduce.py`; build the supplied main revision and pass the absolute binary path as its sole argument. It uses only temporary repositories, a local bare origin, Git, and the real CLI. `artifacts/ideation/transcript.json` records commands, exit codes and outputs.

- #689: real gate prepare/record/consume leaves `pending`; no-hook local `merge guard --verdict passed` exits 0, writes `consumed` and `done`, archives, and prints the manual merge instruction. A subsequent real conflicting `git merge --no-ff task` exits 1; ancestry check confirms task HEAD never reached main. The nonconflicting control delivers only after the premature archive.
- #790: all four flat/folder × local-only/bare-origin cases move the entity and companion artifact, leave HEAD unchanged, and refuse `state commit`. **Already fixed:** refusal now exits 1, not the historical 0. Clean archived publication recovery already exists; preserve it. `status --help` still omits `--archive`.
- #630: missing and initialized-empty state return byte-identical `status --json` entity/pagination objects and exit 0. Missing and non-checkout directories accept `new`; the missing case becomes ignored on-disk state before `state commit` exits 1. `state new` creates valid storage and subsequent filing/commit succeeds. Existing merge guard and `--validate` already diagnose a missing directory; strengthen ordinary status/new rather than claiming those guards are absent.
- Mechanism spike: after an actual successful merge, the existing `status --set task pr=local-merge:<full merge SHA>` and merge guard path finalize a real bound approval. The retry cannot create another archive commit. This seeds delivery-first regression tests; no new approval writer or delivery ledger is needed.

## Proposed approach

1. **AC-1, delivery proof before mutation:** in `internal/status/merge.go`, the `merge: local`, no-registered-hook, passed-verdict path refuses with exit 1 until `pr` holds the existing `local-merge:<SHA>` sentinel. Before calling `finalize`, resolve the configured trunk (`trunk`, default `main`), require that the sentinel names a commit reachable from that trunk, and, when `worktree` is declared, require that worktree HEAD is reachable from the same trunk. Missing/unreadable worktrees or refs refuse, preserving approval and active bytes. A state-only entity without a worktree records the current trunk commit explicitly. The guard prints delivery-first instructions; the operator performs the merge, records the sentinel through the existing setter, and retries. A retry after archive may retain today's archived-read-only refusal; exactly once means no second spend/commit. Registered hooks, PR proof, rejected verdicts and rework keep existing semantics.
   The simpler alternative is only reversing prose or accepting any SHA-shaped sentinel. Neither makes premature archive impossible: a caller can still run guard first or record an undelivered task SHA. Ancestry checks reuse Git, with no new persisted format, coordination or automatic merge engine.
2. **AC-2, explicit split-root retirement owns its archive commit:** route the explicit `status --archive` entry through a small status-package wrapper around existing snapshot, `runArchive`, literal-path-scoped `commitArchiveMove`, rollback and publication helpers. Preflight before movement; keep the raw archive primitive used by merge guard to avoid nesting transactions. A commit failure restores active state; a publication failure retains the local archive commit and exits 1 (rebase HALT remains 3). `state commit <slug>` resumes publication of the clean archive. Report local-only or pushed durability in text/JSON; never call retirement delivered or consume approval. Keep inline archive semantics unchanged.
   The simpler alternative, teaching ordinary archived `state commit` to stage arbitrary dirt, weakens the intentionally publish-only recovery boundary. A special move detector or new flag adds another path where the existing archive snapshot already has the source/destination information. Reuse the transaction at the command that owns the move.
3. **AC-3, storage check before read/filing:** extract a read-only exact-checkout-root/expected-branch check from `internal/statesync/publish.go`, reuse it for entity-dependent status and `new` before filesystem creation. Missing path yields exit 1 with `state-checkout-missing`; an existing ordinary directory, a checkout subdirectory, detached HEAD or wrong branch yields `state-checkout-invalid`. Include the state path and expected branch; JSON gets a single error object instead of `entities: []`. Valid empty state retains today's successful empty object. Preserve `--discover`, README/explicit-file `--read`, and initialization so recovery remains reachable; cross-workflow entity reads validate each contributing workflow. Boot may fail with the same structured error rather than present healthy emptiness.
   Do not invoke `Preflight` for reads: it can abort a rebase. Its existing mutation/sync behavior remains owned there. The simpler directory-exists check permits an ignored subdirectory of the code checkout; this is the observed loss mode. No doctor command, holder registry or automatic checkout creation is needed.

### Independent design review constraints

Independent necessity/soundness review approved this direction; see `artifacts/ideation/design-review.md`. Implementation must check local delivery before both the default finalization path and the existing merged-sentinel shortcut, resolving trunk in the definition's code repository. Test a reachable old trunk sentinel with an undelivered worktree HEAD to prove the second ancestry check independently. Extraction of storage checks must preserve `Preflight`'s rebase handling before HEAD-branch validation, since a rebase can detach HEAD; do not precede archive mutation preflight with that pure branch check. Shared archive helpers need caller-appropriate commit messages and diagnostics so retirement never claims merge-guard delivery. Rollback promises restored active bytes and untouched sibling index entries, not preservation of preexisting staged versions of the retired entity itself. These are constraints within the chosen owners and tolerance; no additional architecture or broad proof is requested. Await the captain's binding ideation approval before implementation.

## Out of scope

External entity path repair, worker fencing, new approval records, additional workflow stages, new state services, runtime adapters, bulk refits, automatic local Git merging, and broad guard policy changes. Raw forced/manual state edits are not a new authority boundary. No code push or CI during ideation.

## Expected surface and tolerance

Estimate net LOC change: **+560, across 15 files** (about 740 insertions, 180 deletions; tests and documentation included). Net tolerance: **+350 to +800**, at most **18 files**. Stop for renewed design review before exceeding either cap or changing an undeclared semantic.

Production owners: `internal/status/merge.go`, `internal/status/native_runner.go`, `internal/status/roots.go`, a small `internal/status/archive_transaction.go`, `internal/statesync/publish.go`, `internal/cli/help.go` (six). Primary proof owners: `internal/cli/terminal_consume_test.go`, `internal/cli/state_commit_test.go`, a new `internal/cli/state_storage_test.go`, `internal/statesync/publish_test.go`, `internal/cli/status_help_test.go`, `internal/status/merge_guard_test.go`, `internal/status/merge_policy_guard_test.go` (seven, including compatibility expectation updates only where the command changed). Docs: `docs/site/reference/command-reference.md` and `docs/site/advanced/split-root-state.md` (two). Helper relocation within these owners is allowed; no new framework/package.

Observable semantics: no new grammar or stored fields. Passed no-hook local completion now requires checked existing delivery proof; retirement in split-root becomes a durable command; missing/invalid storage changes ordinary status/new exit and error output. Archive text/JSON acquires explicit durability, including failure after local commit. Gate authority, PR-hook behavior, archive recovery's publish-only contract and valid-empty/inline status success stay unchanged.

## Acceptance criteria

**AC-1 — Failed local delivery leaves the task active with terminal approval unspent.**
Verified by real temporary repositories and real gate prepare/record/consume. Run guard before delivery and after a real merge conflict; both refuse with unchanged entity bytes, pending approval, active location, no archive commit and task HEAD absent from trunk. Wrong/unreachable sentinel and missing or unreadable worktree refuse. Recovery merges successfully, records the actual sentinel, then consumes and archives exactly once; task HEAD and sentinel are ancestors of the configured trunk. Repeat does not spend or commit twice. Preserve registered-hook/PR and rejected-verdict controls.

**AC-2 — Retirement has a documented command path that commits the full archive move.**
Verified by flat and folder entities, flat companion artifacts and staged dirty siblings, crossed with local-only and bare-origin split-root state. Explicit archive removes tracked active paths, tracks every target artifact, leaves unrelated index/files untouched, preserves the retired status and gate history, and reports the publication result. Inject commit failure to prove full active restoration; inject push failure to prove local archive durability plus nonzero exit; recover with `state commit` and inspect fresh remote state. Same-entity rebase conflict retains existing HALT semantics. An uncommitted move reported as durable success fails.

**AC-3 — Missing state storage cannot masquerade as an empty initialized workflow or accept unsafe filing.**
Verified by missing, ordinary directory, wrong branch, detached HEAD, and valid-empty checkout fixtures. Default/JSON/boot status has a distinct nonzero storage diagnosis for invalid cases; `new` and `status --new` refuse with byte-clean filesystem/index/ref snapshots. Read-only inspection does not abort an existing rebase. Successful `state new` for a fresh workflow and `state init` for a clone restore normal empty status and durable filing. Definition reads and discovery remain usable. A silently created ignored directory fails.

## Test plan

Write focused regressions in the named primary owners before each production change. Deterministic real CLI/Git tests are the main proof; the harness supplies the independent current-main baseline. No live agent lane is needed unless the actual diff expands into runtime/skill behavior.

- Delivery: extend terminal approval CLI fixtures with real branch/worktree merges. Falsifier: restore default direct-finalize, remove reachability checks, or spend before proof. Failed-delivery assertions must fail independently of the success assertions.
- Retirement: extend archive/state publication fixtures. Falsifier: omit the source deletion/companion from staging, commit with a bare index, suppress publication failure, or remove rollback. Fresh clone/`ls-tree` and preserved sibling hashes establish value beyond command output.
- Storage: test through CLI plus pure-check tests in statesync. Falsifier: skip the guard, accept parent-repo Git discovery, or call mutating preflight during a read. Snapshot bytes, Git index and refs before refusals. Include supported init recovery rather than only error strings.
- Output: use existing help/JSON fixture conventions. Falsifier: drop the archive help row or omit typed storage/durability output. These support the behavioral tests rather than replace them.
- Estimated work: 1–2 focused implementation days; local real-Git tests take minutes and need no network service. Run `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`; docs apply their existing build check. Independent validation uses a detached checkout of the proposed stack head and its own conflicting merge, archive failure, and missing-checkout fixtures, not just this harness. Challenge necessity before adding checks; distinguish failure and recovery results in the validation report.

### Ideation baseline check limits

The real 146-command reproduction and two focused existing checks passed. `go test ./...` on the isolated main snapshot exited 1: `internal/cli` and `internal/ensigncycle` hit cumulative 10-minute package timeouts; the installed Codex manifest assertion and two Codex quiet-budget assertions also failed. `internal/status`, `internal/statesync`, and `internal/gates` passed in that run. Preserve `artifacts/ideation/baseline-normal.log` and `baseline-focused.log`; these are baseline observations, not failures introduced by a production diff. `gofmt -w ./cmd ./internal` completed in the isolated snapshot. `go test ./... -race` also exited 1: cumulative CLI/ensigncycle timeouts, the installed-host assertion, and dispatch `TestBuildModsParity` output mismatch; see `baseline-race.log`. No data-race report appeared. Per FO direction, no broad retries or implementation gate are added during ideation.

## Proposed documentation diff

`docs/site/reference/command-reference.md`, add a status archive row:

> `spacedock status --archive <slug>` retires an entity without marking delivery complete. In split-root workflows it commits the full move and publishes state. A commit failure restores the active entity. If publication fails, the local archive commit remains; retry with `spacedock state commit <slug>`.

Append to the existing merge-guard row:

> With `merge: local` and no merge hook, merge onto the configured trunk first. Record `pr=local-merge:<merge SHA>` with `status --set`, then rerun merge guard. It checks that the trunk contains the recorded commit and the declared worktree's HEAD before consuming approval.

Replace the opening status/archive recovery clarification after the existing `state commit` paragraph (the publish-only rule remains):

> To retire an active split-root entity, use `status --archive <slug>`. Use `state commit <slug>` to retry publication of an already committed archive.

`docs/site/advanced/split-root-state.md`, replace “On a fresh clone the state checkout is absent; run `spacedock state init` to restore it.” with:

> On a fresh clone, `status` reports the missing state checkout and filing refuses. Run `spacedock state init` to restore it. For a newly declared workflow with no state branch, use `spacedock state new`. Status distinguishes missing or invalid storage from a valid workflow with no entities.

### Feedback Cycles

## Stage Report: ideation

- DONE: Reproduce all three durability failures against current main with real CLI/Git fixtures and distinguish already-fixed behavior.
  `artifacts/ideation/reproduce.py` and `transcript.json`: 146 actual binary/Git commands at main 4380534 demonstrate premature approval consumption before conflict, all four retirement/storage combinations, unsafe filing, and delivery-first/init recovery; #790's rc=0 is already fixed to rc=1.
- DONE: Design the smallest combined repair in existing owners with exact scope/tolerance and separate failure/recovery acceptance evidence.
  AC-1, AC-2, and AC-3 name independent on-disk/ancestry/remote assertions and failure injections; six production owners reuse existing proof, archive and storage primitives, +560 net LOC / 15 files with +350..+800 / 18-file tolerance; concrete documentation wording is included.
- DONE: Record a bounded implementation and independent-validation plan without new coordination infrastructure.
  The test plan names primary proof owners, distinct falsifying edits, local deterministic cost, and a detached independent audit; baseline `TestTerminalDeliveryFailureReworkRoundTrip` and `TestStateCommitRefusesDirtyArchivedEntityBeforePublication` passed (58.301s), preserving hook recovery and the already-fixed nonzero refusal.

### Summary

The combined repair remains justified on current main, with the dirty-archive exit-code claim corrected. The design makes delivery proof precede approval consumption, makes split-root retirement own its existing archive transaction, and rejects unsafe storage before reads or filing; the implementation plan stays above the scheduling PR in the requested stack.
