---
id: 7fhzvvk8d5smj858bp47xbjq
title: Collapse the gate-approval ceremony from ~16 tool calls to 1-2
status: ideation
source: "Ships-counselor friction rollup, 2026-08-02, theme 1 (gate-approval ceremony): measured on sonnet-gate-guardrail-no-authority's ideation->implementation gate -- 16 discrete FO tool calls and 156s wall clock to apply one captain word ('approve'), with 0 additional captain turns needed. Recurs at every nonterminal gate for every entity, forever. Captain directed: file and dispatch (ideation, via a fable-model ensign)."
started: 2026-08-02T06:25:27Z
completed:
verdict:
score: 0.6
worktree:
issue:
gates:
    version: 1
    current:
        gate: gate:7fhzvvk8d5smj858bp47xbjq:backlog
    records:
        - id: gate:7fhzvvk8d5smj858bp47xbjq:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:7fhzvvk8d5smj858bp47xbjq-backlog-1
              briefing:
                id: briefing:7fhzvvk8d5smj858bp47xbjq:backlog:attempt-1:revision-1
                digest: sha256:2cb1660e6fb16e9c6303564fc34d957a34155eda05ed8646a06b770a92ed221c
                digest-domain: canonical-bytes
                request-digest: sha256:a9c5b97c54d3372732e38b8ae7157af41a6aaca2dae2cdecece420ab2fd0572a
                room-ref: ./collapse-gate-approval-ceremony/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:7fhzvvk8d5smj858bp47xbjq:backlog:1
                briefing: briefing:7fhzvvk8d5smj858bp47xbjq:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-02T06:25:11.852955Z"
                decision: approve
                reason: 'Captain directed in chat: file the gate-approval-ceremony friction as a task and dispatch it to ideation via a fable-model ensign.'
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

Applying one captain gate decision (approve/revise/hold) currently costs ~16 raw FO tool calls and ~2.5 minutes of wall clock, most of it mechanical: re-invoking `state commit` after nearly every binary call, then a separate set of frontmatter stamps and worktree creation before the next dispatch even begins. This entity is about collapsing that ceremony toward 1-2 calls per gate decision without weakening any of the authority/integrity checks the ceremony exists to enforce.

## Problem

Measured this session (`sonnet-gate-guardrail-no-authority`, ideation -> implementation gate, exact tool-call sequence from the transcript): `status --set` (retitle) -> `gate record --decision approve` -> `state commit` -> `gate consume` -> `state commit` -> `status --read` (verify) -> `status --set` (worktree stamp) -> `state commit` -> `git status` -> `git worktree add` -> 2x `Write` (checklist/scope-notes scratch files) -> `ToolSearch` -> `SendMessage` (prior-worker shutdown) -> `dispatch build` -> `Agent` (the actual spawn). 16 calls, 156s, 0 additional captain turns, to apply the word "approve".

3-4 of the 16 are the same `state commit` re-invoked idempotently right after `gate record` and again after `gate consume` -- each of those binary commands already knows it just mutated durable state. The remaining bulk is the post-consume "ordinary dispatch" procedure (frontmatter stamps, worktree creation) that always follows a nonterminal consume in the same fixed shape.

This is not a one-off cost: every entity, at every nonterminal `gate: true` stage, pays this in full. A workflow that runs 50 entities through 2 gates each pays it 100 times.

## Proposed approach

Two collapsed commands replace the ten binary/git commands of the measured ceremony. All three seed candidates land, two of them reshaped: candidate 1 as sketched; candidate 2 as a `--consume` flag on `gate record` rather than a `gate approve` verb; candidate 3 as a `--stamp` flag on `dispatch build` rather than a `dispatch advance` verb.

Collapsed ceremony (nonterminal captain approve, fresh dispatch to a worktree stage -- the measured shape):

```
spacedock gate record <slug> --decision approve --actor person:captain --reason "..." --consume --workflow-dir DIR
spacedock dispatch build --stamp --workflow-dir DIR --entity-path ... --stage ... --checklist-file ... [--scope-notes-file ...] [--team-name ...]
```

then the host spawn (`Agent(...)`) from the emitted envelope, exactly as today. Post-gate reuse-advance is also 2 commands: `gate record --consume` + `dispatch build --advance` -- consume already advanced status, the worktree is unchanged, `started` is already stamped, so no stamps are needed there at all (and `--stamp` + `--advance` is rejected by design, see mechanism 3).

### Mechanism 1 -- implicit state sync in the mutating gate verbs

`gate prepare`, `gate record` (all three semantic sources, including `--round`), and `gate consume` end by running the exact `state commit` machinery -- path-scoped `add`+`commit` restricted to the entity unit, then the shared statesync publish (push, on-reject pull --rebase, re-push) -- for the entity they just wrote.

- Split-root: one machine-parseable line `sync=pushed|local-only|no-op` appended to the verb's stdout. Inline workflows: no state checkout, no sync line, output byte-identical to today.
- Commit messages name the verb: `gate: prepare <slug> <stage>`, `gate: record <slug> <stage> <decision>`, `gate: consume <slug> -> <target-stage>`, replacing the generic `state: update <slug>`.
- Failure semantics preserve "close/commit failure halts": when sync fails, the gate write is already durable -- stderr detail + exit 1; a same-entity rebase conflict propagates the existing `state commit` HALT rendering and exit 3. No force-push, no auto-resolve; statesync already owns that discipline.
- Consume's stale path (pending -> superseded is a durable write on a nonzero-exit consume) commits before the exit-1 propagates -- spike-verified ordering.
- The standalone `state commit` verb remains for the flows that still need it: entity seeding after `new`, FO body/report edits, merge-archive publication.

### Mechanism 2 -- `gate record ... --consume`

```
spacedock gate record <entity> --decision approve --actor ID [--reason TEXT] --consume [--workflow-dir DIR]
spacedock gate record <entity> --room PATH --consume [--workflow-dir DIR]
```

After a successful approve-close, run consume in the same invocation -- literally the existing record handler then the existing consume handler, each with its mechanism-1 sync, so the Resolution is committed before the consume attempt and the recoverability ordering the lifecycle depends on is preserved mechanically. Output: the record line, its sync line, the consume line (with `consumed=`/`target-stage=`/`route=`), its sync line -- the FO's post-consume `status --read` verification collapses into reading this output.

- Chat source: `--consume` with `--decision revise|hold` is a usage error (exit 2, before any write). The flag never softens a non-approve decision.
- Room source: the decision lives inside the room; on a revise/hold close, report the close and skip consume (`consume=skipped`, exit 0).
- Terminal target: consume spends nothing, reports `route=approved-awaiting-merge`, exit 0; merge guard remains the sole terminal consumer.
- Ineligible/stale/blocked after close: the close stays durable; exit nonzero with the standalone consume diagnostics.
- `--actor` stays required with no default: the captain-vs-FO actor distinction is untouched.

### Mechanism 3 -- `dispatch build --stamp`

`--stamp` folds fo-dispatch-core dispatch steps 5-7 (frontmatter stamps, state commit, worktree creation) into the build invocation, executed before artifact assembly:

1. Refuse unless entity `status` == `--stage` (nonzero exit, no mutation). `--stamp` NEVER changes `status:` -- status advancement stays owned by `gate consume` (gated) or `status --set` (non-gated), which mechanically preserves "Never use `status --set` to advance a gate" and this entity's authority out-of-scope.
2. Stamp `started` (skipped if already set) and, for a worktree-declaring stage, `worktree=.worktrees/{worker_key}-{slug}`, through the native `status --set` machinery -- inheriting its mutation guards (mod-block, stage membership, away-status refusals) rather than adding a parallel frontmatter writer.
3. Path-scoped state commit + publish (the mechanism-1 seam), message `dispatch: <slug> entering <stage>`.
4. `git worktree add -b {worker_key}/{slug} .worktrees/{worker_key}-{slug}` at the main repo root when the stage declares a worktree and the path is absent; an existing worktree is a skip, not an error.
5. Assemble and emit the envelope exactly as today.

`--stamp` is opt-in: bare `dispatch build` stays a pure artifact assembler, so the existing build test corpus, break-glass retries, and `--advance`/`--feedback-reflow` paths are untouched. `--stamp` with `--advance` is a usage error (the post-gate reuse path needs no stamps). Inline-workflow `--stamp` stamps frontmatter but leaves the main-repo commit to the FO as today -- the measured ceremony is split-root; folding inline main-repo commits is a separate decision this entity does not take.

### Alternatives rejected

1. **`gate approve <slug>` convenience verb** (the seed sketch): duplicates record's source/actor parsing in a second command surface; cannot cover the room-backed approve path without re-implementing source selection; and invites a defaulted `--actor person:captain`, eroding the explicit actor distinction. The `--consume` flag reaches the same one-call shape with strictly less mechanism, covers both semantic sources, and makes AC-2's "same underlying operations" proof structural -- the flag sequences the two existing handlers. No `gate revise`/`gate hold` siblings are needed either: with mechanism 1, a revise/hold close is already a single call.
2. **`dispatch advance <slug>` verb** (the seed sketch): would always be immediately followed by `dispatch build`, never issued alone -- re-creating a two-call ceremony where one suffices. And `advance` already means live-worker reuse-advance in this namespace (`dispatch build --advance`); a same-namespace homonym with different semantics is an operator hazard for LLM FOs.
3. **Automatic (flagless) stamping in `dispatch build`**: repurposes a pure assembler into an unconditional mutator -- breaks the existing build test corpus's assumptions, makes failed-spawn retries and break-glass replays mutate state, and entangles the `--advance`/`--feedback-reflow` modes that need different (or no) stamps.
4. **Batch/deferred commit** (one commit at ceremony end): violates the recoverability ordering (the Resolution must be durable before any consume attempt); saves no FO calls once commits are implicit.
5. **Auto-committing `status --set` generally**: touches every FO and ensign flow for a ceremony that no longer issues a bare `status --set`; unjustified blast radius for AC-1.

Each mechanism serves AC-1 directly: mechanism 1 removes the 3 idempotent `state commit` re-invocations (and the class of forgot-to-commit nudge loops); mechanism 2 merges record+consume+verify-read into one call; mechanism 3 removes the 4-call post-consume block (`status --set`, `state commit`, `git status`, `git worktree add`). The simplest alternative for each -- keep hand-running the calls -- is exactly the measured 10-command ceremony.

## Expected surface

~12 files, ~800 insertions / ~60 deletions, tolerance +/-40%:

- Go production (~250 lines): `internal/cli/cli.go` (gate sync seam + `--consume`), `internal/cli/state_sync.go` (seam extraction from `runStateCommit`), `internal/dispatch/build.go` plus a new `stamp.go` (`--stamp`).
- Tests (~450 lines): `internal/cli/gate_test.go` (fixture state-branch alignment + `--consume`/sync cases), new `internal/cli/gate_ceremony_count_test.go` (the AC-1 harness), `internal/dispatch/build_stamp_test.go`.
- Prose/docs (~120 lines churn): `skills/fo-gate-lifecycle/SKILL.md`, `skills/first-officer/references/fo-dispatch-core.md`, `skills/first-officer/references/first-officer-shared-core.md`, `docs/site/reference/command-reference.md`, `docs/site/concepts/gates-and-decisions.md`.

Declared observable-semantics changes: **command grammar** (`gate record` gains `--consume`; `dispatch build` gains `--stamp`); **command output** (gate verbs append a `sync=...` line in split-root mode; new nonzero exits after a durable write when sync fails, including the existing exit-3 HALT shape); **runtime behavior** (gate verbs now commit and push the state checkout -- network I/O inside previously local verbs). NOT changed: frontmatter and stored formats, gate authority and actor rules, digest/Briefing integrity, `gate prepare` semantics, the room-backed Result path, terminal merge-guard routing.

## Doc diff (load-bearing lines, before -> after)

`skills/fo-gate-lifecycle/SKILL.md`:

- Line 29 (prepare), before: "Preparation binds two recorder-ready files without source copies. `«state.commit»(slug)` commits the folder entity or flat Markdown-plus-companion room unit." After: "Preparation binds two recorder-ready files without source copies and commits+syncs them itself (`sync=` line; sync failure or HALT is nonzero)."
- Line 35 (record block): add as first spelling: "`${SPACEDOCK_BIN:-spacedock} gate record ENTITY --decision approve --actor person:captain [--reason REASON] --consume --workflow-dir WORKFLOW_DIR` -- the captain-approve fast path: close, sync, consume, sync in one call."
- Line 48 (close), before: "After every successful close, `«state.commit»(slug)` must commit that exact Resolution before approve, revise, hold, or any consume attempt. Close/commit failure halts." After: "Record commits and syncs its own Resolution before returning (and before its `--consume` consume attempt); require the `sync=` line. A nonzero exit after the recorded line means the close is durable but unsynced -- halt."
- Line 56 (consume), before: "it atomically writes successor status plus consumed state -- commit through `«state.commit»(slug)`, then ordinary dispatch." After: "it atomically writes successor status plus consumed state and commits+syncs it (`sync=` line), then ordinary dispatch via `dispatch build --stamp`."

`skills/first-officer/references/fo-dispatch-core.md` dispatch steps 5-7, before: the three-step `status --set` stamp / commit-on-main / worktree-create sequence. After: one step -- "For a gate-consumed entry, `status` is already advanced: run `«dispatch.build» --stamp`, which stamps `started`/`worktree=`, commits+syncs state, and creates the declared worktree before emitting the envelope. For a non-gated entry, first advance with `${SPACEDOCK_BIN:-spacedock} status --workflow-dir {workflow_dir} --set {slug} status={next_stage}`, then the same `--stamp` build." The `«dispatch.build»` effect block gains the `--stamp` flag line and its `--advance` incompatibility.

`docs/site/reference/command-reference.md`: the `gate record` rows gain the `--consume` composition sentence; the `gate consume` row and a new `dispatch build --stamp` note gain the sync/stamp behavior; the `state commit` section notes gate verbs now sync themselves.

`docs/site/concepts/gates-and-decisions.md` (~line 81), before: "After an authorized decision, it records and commits the Resolution before every route." After: "After an authorized decision, the recorder itself commits and syncs the Resolution before every route (`--consume` folds the approve's consume into the same call)."

## Ideation spike record

Riskiest unverified mechanism (per the seed test plan): folding state commit into `gate record`/`gate consume` without breaking existing gate-lifecycle tests. Spiked 2026-08-02 in an isolated clone at 23ed415bb:

- ~20-line spike seam: after record/consume success (and consume's stale write), invoke the existing `runStateCommit` machinery with the slug derived from the resolved entity path -- pure composition, no gates-package changes.
- Full `go test ./internal/cli` (all 15 gate test functions plus the rest of the package, 60s) green after exactly ONE fixture line: `gatePrepareCLIFixture`'s state checkout needed `git branch -M spacedock-state/dev` -- the fixture sat on `main` while real checkouts are born on the state branch by `state init`/`new`; statesync preflight rightly refuses the mismatch. `internal/gates` and `internal/statesync` untouched and green.
- End-to-end on a split-root fixture: `gate record --decision approve` -> Resolution + path-scoped commit (`state: update task`), no-origin local-only handled; `gate consume` -> status ideation->implementation + consumed mark + second commit; porcelain clean after both.
- Conclusion: the fold is safe; the one-line fixture alignment is the entire observed breakage surface, and it must be declared under AC-2.

No spike needed for the other two mechanisms -- proven code paths: `--consume` sequences the two handlers the spike exercised in-process; `--stamp` writes through the native `status --set` handlers (existing corpus), the spiked commit seam, and `git worktree add` driven from Go as already shipped in `internal/cli/state.go` (`state init`/`new`).

## Out of scope

Any weakening of gate authority, digest/Briefing integrity checks, frozen-binding refusals, or the captain-vs-FO actor distinction. Any change to `gate prepare`'s semantics or the room-backed Result path. The mechanism-2 byte-cap follow-up (`gate-lifecycle-hardening-byte-budget`) -- unrelated. Batch/bulk gate operations across multiple entities -- this is about one entity's one gate decision.

## Acceptance criteria

**AC-1 (VALUE) - The binary/git command count for one captain "approve" at a nonterminal `gate: true` stage, from open gate through emitted dispatch envelope, drops 10 -> 2 against the transcript baseline.**
Independent baseline: the 16 FO tool calls in `## Problem` (sonnet-gate-guardrail-no-authority transcript). Of those 16, ten are binary/git commands once the incidental retitle `--set` is excluded (`gate record`, `state commit`, `gate consume`, `state commit`, `status --read`, `status --set`, `state commit`, `git status`, `git worktree add`, `dispatch build`); the other five (two checklist/scope-note Writes, ToolSearch, prior-worker shutdown SendMessage, the Agent spawn) are host-side calls outside the binary's reach, so the post-change ceiling on total FO tool calls is <=7.
Verified by: `internal/cli/gate_ceremony_count_test.go`, a fixture harness that drives the compiled command surface twice over the same split-root fixture (gated stage -> worktree-declaring successor): once through the 10-command before-list transcribed from the transcript, once through the 2-command after-list (`gate record --decision approve --actor person:captain --consume`, `dispatch build --stamp`). The harness asserts (a) every command exits 0, (b) both runs converge to equivalent end-state predicates -- Resolution closed and application consumed, status advanced to the successor, `started`/`worktree=` stamped, state-checkout log containing the ceremony commits with a clean tree, worktree present on its branch, spawn envelope emitted -- and (c) the after-list length is 2 while the before-list length is 10. The counts are literal command-list lengths: a regression that adds a required step grows the after-list and fails the test, so the number can move the wrong way.

**AC-2 - No loss of gate authority or integrity guarantees.**
Verified by: the existing gate-lifecycle test corpus (frozen-binding refusal, digest mismatch rejection, actor/authority validation, room-backed vs chat-decision paths, terminal routing) passes with zero production-assertion changes and exactly one declared fixture change -- `gatePrepareCLIFixture`'s state checkout moves to the declared state branch (one line, spike-verified). New negative tests pin the preserved boundaries: `--consume` with `--decision revise|hold` exits 2 before any write; `--stamp` with entity status != `--stage` refuses without mutation; `--stamp --advance` exits 2. Composition is structural, not parallel: `--consume` sequences the existing record and consume handlers; `--stamp` writes through the native `status --set` machinery (inheriting its mutation guards) and the same `state commit` seam.

**AC-3 - The FO contract and user docs carry the collapsed call shape as the primary path.**
Verified by: the `## Doc diff` above is applied to `skills/fo-gate-lifecycle/SKILL.md`, `skills/first-officer/references/fo-dispatch-core.md`, `skills/first-officer/references/first-officer-shared-core.md`, `docs/site/reference/command-reference.md`, and `docs/site/concepts/gates-and-decisions.md`; grep over fo-gate-lifecycle finds no remaining post-close/post-consume `«state.commit»` instruction (remaining `«state.commit»` references live only where the standalone verb is still real: entity seeding, body/report edits, merge archive); fo-dispatch-core's dispatch steps no longer instruct separate stamp/commit/worktree-create commands as the primary path.

## Test plan

- `internal/cli/gate_ceremony_count_test.go` (AC-1): the before/after harness above -- temp git fixtures (main repo + split-root state checkout on `spacedock-state/dev`), no live host, ~250 lines. This reproduces the baseline as a runnable count first (the before-list), then measures the after-list on the same fixture shape, so before/after is a real number.
- Behavior tests (AC-2, mechanisms): implicit sync line and exit semantics per verb -- pushed/local-only/no-op, and the exit-3 HALT via the same origin/conflict fixture patterns `state_commit_test.go` already uses; `--consume` happy path, usage-error path, room-source revise skip, terminal `route=approved-awaiting-merge`, stale-supersede commit-before-exit; `--stamp` stamps+commit+worktree-add, status-mismatch refusal, idempotent re-run (existing worktree, already-set `started`), non-worktree stage (`started` only), `--advance` incompatibility.
- Prose (AC-3): grep-level absence/presence checks plus the doc diff applied. Static checks are acceptable here because the behavioral half is exercised by the harness; instruction text has no runtime of its own.
- No live workflow run: the prose change swaps call shapes over identical underlying operations and the harness drives the real binary, so live-runtime risk is judged low. If the gate review disagrees, the fallback is a single live gate on a scratch entity.

## Stage Report: ideation

- DONE: Choose and design the concrete collapse mechanism(s) among the three candidates (implicit state commit in gate record/consume; a gate approve convenience verb; a dispatch advance verb) or a considered combination, with exact command signatures and behavior, and record why the alternative(s) not chosen are insufficient.
  All three candidates land, two reshaped: implicit sync in gate verbs; `--consume` on `gate record` (approve-verb rejected: duplicate parse surface, no room-source coverage, actor-default hazard); `--stamp` on `dispatch build` (advance-verb rejected: always paired with build, `--advance` homonym). Signatures and five rejected alternatives in `## Proposed approach`.
- DONE: Refine Acceptance criteria and Test plan: in particular, specify exactly how AC-1's FO-issued-command count is measured (fixture harness or transcript replay) against the independent baseline of 16, and confirm AC-2/AC-3 are concretely verifiable.
  AC-1: `gate_ceremony_count_test.go` fixture harness, before-list (10 binary/git commands transcribed from the 16-call transcript) vs after-list (2), identical end-state predicates; host-side remainder documented (total <=7 of 16). AC-2 names the one declared fixture change + new negative tests; AC-3 enumerates the five prose/doc files with a written doc diff.
- DONE: Spike the riskiest unverified mechanism first per the Proof policy (e.g. whether state commit can be safely folded into gate record/consume without breaking existing gate-lifecycle tests) and record the result, or record "no spike needed" with the proven mechanisms relied on.
  Spiked in an isolated clone at 23ed415bb: ~20-line seam reusing `runStateCommit`; full `internal/cli` suite green after one fixture line (state checkout branch `-M spacedock-state/dev`); end-to-end record->commit and consume->advance+commit exercised, porcelain clean. "No spike needed" recorded for `--consume` and `--stamp` with the proven mechanisms named. Details in `## Ideation spike record`.

### Summary

Designed the ceremony collapse as two flags on existing verbs plus implicit state sync in the mutating gate verbs, taking the measured 10 binary/git commands per captain-approve to 2 (`gate record --decision approve --actor person:captain --consume`, then `dispatch build --stamp`) without moving any authority boundary: `--stamp` can never change `status:`, `--actor` keeps no default, and composition reuses the existing handlers and mutation guards. The riskiest mechanism (folding state commit into record/consume) was spiked against the full CLI test suite -- safe, with a one-line fixture alignment as the entire breakage surface. AC-1 now has a concrete falsifiable harness design measuring 10 -> 2 against the transcript baseline.
