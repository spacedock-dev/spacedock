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

### Mechanism 1 -- implicit state sync in the gate close/consume verbs

`gate record` (close sources: `--decision`, `--room`) and `gate consume` end by running the exact `state commit` machinery -- path-scoped `add`+`commit` restricted to the entity unit, then the shared statesync publish (push, on-reject pull --rebase, re-push) -- for the entity they just wrote. The scope is deliberately this narrow: these two verbs are the only ones inside AC-1's measured window (the ceremony starts at the open gate, at `gate record`), and `--consume`'s recoverability ordering (Resolution durable before the consume attempt) requires the sync to live in record's own path, not only in the composite flag. `gate prepare` and `gate record --round` are NOT in scope -- see Alternatives rejected 6; the prepare-written room rides the close commit regardless (the entity commit unit already includes the flat entity's companion dir -- observed in the spike fixtures), and the post-prepare `«state.commit»` step is unchanged.

- Split-root: one machine-parseable line `sync=pushed|local-only|no-op|failed|halted phase=record|consume` appended to the verb's stdout in EVERY outcome, including failures -- callers branch on the final sync line plus the exit code, never by inferring from which prose lines printed. Inline workflows: no state checkout, no sync line, output byte-identical to today.
- Commit messages name the verb: `gate: record <slug> <stage> <decision>`, `gate: consume <slug> -> <target-stage>`, replacing the generic `state: update <slug>`.
- Failure semantics preserve "close/commit failure halts": when sync fails, the gate write is already durable -- `sync=failed` + stderr detail + exit 1; a same-entity rebase conflict emits `sync=halted` and propagates the existing `state commit` HALT rendering and exit 3 (spike-demonstrated, see the spike record). No force-push, no auto-resolve; statesync already owns that discipline.
- Consume's stale path (pending -> superseded is a durable write on a nonzero-exit consume) commits before the exit-1 propagates -- spike-verified ordering.
- The standalone `state commit` verb remains for the flows that keep it: `gate prepare` binding commits, `gate record --round` publications, entity seeding after `new`, FO body/report edits, merge-archive publication.

### Mechanism 2 -- `gate record ... --consume`

```
spacedock gate record <entity> --decision approve --actor ID [--reason TEXT] --consume [--workflow-dir DIR]
spacedock gate record <entity> --room PATH --consume [--workflow-dir DIR]
```

After a successful approve-close, run consume in the same invocation -- literally the existing record handler then the existing consume handler, each with its mechanism-1 sync, so the Resolution is committed before the consume attempt and the recoverability ordering the lifecycle depends on is preserved mechanically. Output: the record line, its `sync=... phase=record` line, the consume line (with `consumed=`/`target-stage=`/`route=`), its `sync=... phase=consume` line -- the FO's post-consume `status --read` verification collapses into reading this output.

- Chat source: `--consume` with `--decision revise|hold` is a usage error (exit 2, before any write). The flag never softens a non-approve decision.
- Room source: the decision lives inside the room; on a revise/hold close, report the close and skip consume (`consume=skipped`, exit 0).
- Terminal target: consume spends nothing, reports `route=approved-awaiting-merge`, exit 0; merge guard remains the sole terminal consumer.
- Ineligible/stale/blocked after close: the close stays durable; exit nonzero with the standalone consume diagnostics.
- `--actor` stays required with no default: the captain-vs-FO actor distinction is untouched.

One `--consume` invocation has multiple distinct durable landing positions; the discriminator is the FINAL `sync=`/`phase=` line plus the exit code -- machine-parseable, never inferred from prose:

| Landing position | Final line | Exit |
|---|---|---|
| close refused (validation, nothing durable) | no `recorded` line, no sync line | 1 or 2 |
| close durable, record-sync failed | `sync=failed phase=record` | 1 |
| close durable, record-sync rebase conflict | `sync=halted phase=record` + HALT stderr | 3 |
| close durable+synced, consume refused (ineligible/blocked, no write) | `sync=... phase=record` + consume diagnostic | 1 |
| close durable+synced, consume stale (superseded write, synced) | `sync=... phase=consume` | 1 |
| advanced, consume-sync failed | `sync=failed phase=consume` | 1 |
| advanced, consume-sync rebase conflict | `sync=halted phase=consume` + HALT stderr | 3 |
| advanced and synced (or terminal `route=`) | `sync=pushed\|local-only\|no-op phase=consume` | 0 |

### Mechanism 3 -- `dispatch build --stamp`

`--stamp` folds fo-dispatch-core dispatch steps 5-7 (frontmatter stamps, state commit, worktree creation) into the build invocation, executed before artifact assembly:

1. Refuse unless entity `status` == `--stage` (nonzero exit, no mutation). `--stamp` NEVER changes `status:` -- status advancement stays owned by `gate consume` (gated) or `status --set` (non-gated), which mechanically preserves "Never use `status --set` to advance a gate" and this entity's authority out-of-scope.
2. Stamp `started` (skipped if already set) and, for a worktree-declaring stage, `worktree=.worktrees/{worker_key}-{slug}`, through the native `status --set` machinery -- inheriting its mutation guards (mod-block, stage membership, away-status refusals) rather than adding a parallel frontmatter writer.
3. Path-scoped state commit + publish (the mechanism-1 seam), message `dispatch: <slug> entering <stage>`.
4. `git worktree add -b {worker_key}/{slug} .worktrees/{worker_key}-{slug}` at the main repo root when the stage declares a worktree and the path is absent; an existing worktree is a skip, not an error.
5. Assemble and emit the envelope exactly as today.

`--stamp` is opt-in: bare `dispatch build` stays a pure artifact assembler, so the existing build test corpus, break-glass retries, and `--advance`/`--feedback-reflow` paths are untouched. `--stamp` with `--advance` is a usage error (the post-gate reuse path needs no stamps). Inline-workflow `--stamp` stamps frontmatter but leaves the main-repo commit to the FO as today -- the measured ceremony is split-root; folding inline main-repo commits is a separate decision this entity does not take.

Failure discrimination (stdout stays a pure JSON envelope, so build's discriminator lives on stderr + exit code, unlike the gate verbs' stdout sync line): every `--stamp`-phase failure writes a `dispatch build --stamp:`-prefixed stderr diagnostic and emits NO envelope; a stamp-sync rebase conflict propagates the state-sync HALT rendering with exit 3. This matters because fo-dispatch-core's existing block clause reads any nonzero `dispatch build` as the Break-Glass Manual Dispatch trigger -- but manually dispatching against a HALTed or unsynced state tree is exactly what `first-officer-shared-core.md`'s HALT clause forbids. The doc diff below disambiguates: break-glass remains the remedy only for artifact-assembly failures (no `--stamp:` stderr prefix); a `--stamp`-phase failure burned no authority and emitted nothing -- remedy the named stamp/sync problem and rerun the same build; exit 3 is `«halt.rebase-conflict»` -- halt dispatch entirely, never break-glass.

### Alternatives rejected

1. **`gate approve <slug>` convenience verb** (the seed sketch): duplicates record's source/actor parsing in a second command surface; cannot cover the room-backed approve path without re-implementing source selection; and invites a defaulted `--actor person:captain`, eroding the explicit actor distinction. The `--consume` flag reaches the same one-call shape with strictly less mechanism, covers both semantic sources, and makes AC-2's "same underlying operations" proof structural -- the flag sequences the two existing handlers. No `gate revise`/`gate hold` siblings are needed either: with mechanism 1, a revise/hold close is already a single call.
2. **`dispatch advance <slug>` verb** (the seed sketch): would always be immediately followed by `dispatch build`, never issued alone -- re-creating a two-call ceremony where one suffices. And `advance` already means live-worker reuse-advance in this namespace (`dispatch build --advance`); a same-namespace homonym with different semantics is an operator hazard for LLM FOs.
3. **Automatic (flagless) stamping in `dispatch build`**: repurposes a pure assembler into an unconditional mutator -- breaks the existing build test corpus's assumptions, makes failed-spawn retries and break-glass replays mutate state, and entangles the `--advance`/`--feedback-reflow` modes that need different (or no) stamps.
4. **Batch/deferred commit** (one commit at ceremony end): violates the recoverability ordering (the Resolution must be durable before any consume attempt); saves no FO calls once commits are implicit.
5. **Auto-committing `status --set` generally**: touches every FO and ensign flow for a ceremony that no longer issues a bare `status --set`; unjustified blast radius for AC-1.
6. **Implicit sync in `gate prepare` (and `gate record --round`)**: an earlier draft included prepare "for consistency." Dropped on staff review: prepare is outside AC-1's measured window (the ceremony starts at the open gate, at `gate record`), so folding push/rebase into it buys zero measured calls while adding a new network-failure surface to a verb whose nonzero-exit remedy prose ("refresh or rebuild the version-gated bundle", fo-gate-lifecycle line 27) would then misroute sync failures -- new failure class, new doc surface, no value transaction. Consistency alone is not necessity. Same ruling for `--round` publications. Both keep today's explicit `«state.commit»`, and prepare's block clause stays valid as written.

Each mechanism serves AC-1 directly: mechanism 1 removes the 3 idempotent `state commit` re-invocations (and the class of forgot-to-commit nudge loops); mechanism 2 merges record+consume+verify-read into one call; mechanism 3 removes the 4-call post-consume block (`status --set`, `state commit`, `git status`, `git worktree add`). The simplest alternative for each -- keep hand-running the calls -- is exactly the measured 10-command ceremony.

## Expected surface

~12 files, ~800 insertions / ~60 deletions, tolerance +/-40%:

- Go production (~250 lines): `internal/cli/cli.go` (gate sync seam + `--consume`), `internal/cli/state_sync.go` (seam extraction from `runStateCommit`), `internal/dispatch/build.go` plus a new `stamp.go` (`--stamp`).
- Tests (~450 lines): `internal/cli/gate_test.go` (fixture state-branch alignment + `--consume`/sync cases), new `internal/cli/gate_ceremony_count_test.go` (the AC-1 harness), `internal/dispatch/build_stamp_test.go`.
- Prose/docs (~120 lines churn): `skills/fo-gate-lifecycle/SKILL.md`, `skills/first-officer/references/fo-dispatch-core.md`, `skills/first-officer/references/first-officer-shared-core.md`, `docs/site/reference/command-reference.md`, `docs/site/concepts/gates-and-decisions.md`.

Declared observable-semantics changes: **command grammar** (`gate record` gains `--consume`; `dispatch build` gains `--stamp`); **command output** (the gate close/consume verbs append a `sync=... phase=...` line in split-root mode; new nonzero exits after a durable write when sync fails, including the existing exit-3 HALT shape; `--stamp` failures carry a `dispatch build --stamp:` stderr prefix); **runtime behavior** (the gate close/consume verbs and `--stamp` builds now commit and push the state checkout -- network I/O inside previously local verbs). NOT changed: frontmatter and stored formats, gate authority and actor rules, digest/Briefing integrity, `gate prepare` semantics and output, `gate record --round`, the room-backed Result path, terminal merge-guard routing.

## Doc diff (load-bearing lines, before -> after)

`skills/fo-gate-lifecycle/SKILL.md`:

- Lines 27/29 (prepare capability check and bind commit): UNCHANGED. Prepare gains no sync (Alternatives rejected 6), so its existing nonzero-exit remedy ("refresh or rebuild the selected version-gated bundle") stays the correct and only reading of a prepare failure, and the post-bind `«state.commit»(slug)` sentence stays as written.
- Line 35 (record block): add as first spelling: "`${SPACEDOCK_BIN:-spacedock} gate record ENTITY --decision approve --actor person:captain [--reason REASON] --consume --workflow-dir WORKFLOW_DIR` -- the captain-approve fast path: close, sync, consume, sync in one call."
- Line 48 (close), before: "After every successful close, `«state.commit»(slug)` must commit that exact Resolution before approve, revise, hold, or any consume attempt. Close/commit failure halts." After: "Record commits and syncs its own Resolution before returning (and before its `--consume` consume attempt). Branch on the FINAL `sync=... phase=record|consume` line plus the exit code, never on which lines printed: `phase=record` nonzero means the close is durable but the advance did not happen (sync failure, or a consume refusal after a clean record sync -- read the consume diagnostic); `phase=consume` nonzero means consume wrote (advance or supersede) but its sync failed or HALTed; `sync=halted` exit 3 is `«halt.rebase-conflict»` -- HALT. Close/sync failure halts."
- Line 56 (consume), before: "it atomically writes successor status plus consumed state -- commit through `«state.commit»(slug)`, then ordinary dispatch." After: "it atomically writes successor status plus consumed state and commits+syncs it (`sync=... phase=consume` line), then ordinary dispatch via `dispatch build --stamp`."

`skills/first-officer/references/first-officer-shared-core.md`:

- Line 117, before: "Every state write is one call: `«state.commit»(slug)`." After: "Every state write is one call: `«state.commit»(slug)` -- and the gate close/consume verbs plus `dispatch build --stamp` make that call themselves; their `sync=` line (or `--stamp`'s stderr/exit) is the call's result. Do not re-run `«state.commit»` after them."
- Line 133 block clause: unchanged in text; the gate verbs' and `--stamp`'s exit-3 HALT is the same `«state.commit»` exit 3 it already covers -- HALT per stderr, never force, never dispatch against the halted tree.

`skills/first-officer/references/fo-dispatch-core.md`:

- Dispatch steps 5-7, before: the three-step `status --set` stamp / commit-on-main / worktree-create sequence. After: one step -- "For a gate-consumed entry, `status` is already advanced: run `«dispatch.build» --stamp`, which stamps `started`/`worktree=`, commits+syncs state, and creates the declared worktree before emitting the envelope. For a non-gated entry, first advance with `${SPACEDOCK_BIN:-spacedock} status --workflow-dir {workflow_dir} --set {slug} status={next_stage}`, then the same `--stamp` build." The `«dispatch.build»` effect block gains the `--stamp` flag line and its `--advance` incompatibility.
- Line 152 (`«dispatch.build»` block clause), before: "on non-zero exit (or missing binary) ONLY -- read stderr, report the helper failure to the captain, then use the adapter's Break-Glass Manual Dispatch template ...". After: "on non-zero exit (or missing binary), read stderr FIRST: a `dispatch build --stamp:`-prefixed diagnostic is a stamp/sync failure -- no envelope was emitted and no authority burned; remedy the named problem and rerun the same build, never break-glass. Exit 3 is `«halt.rebase-conflict»` -- HALT dispatch entirely (shared-core HALT clause); never manually dispatch against a halted or unsynced state tree. Only an assembly failure (nonzero WITHOUT the `--stamp:` prefix, or missing binary) triggers the adapter's Break-Glass Manual Dispatch template." This removes the collision between the break-glass trigger and the shared-core HALT rule.

`docs/site/reference/command-reference.md`: the `gate record` rows gain the `--consume` composition sentence; the `gate consume` row and a new `dispatch build --stamp` note gain the sync/stamp behavior; the `state commit` section notes gate verbs now sync themselves.

`docs/site/concepts/gates-and-decisions.md` (~line 81), before: "After an authorized decision, it records and commits the Resolution before every route." After: "After an authorized decision, the recorder itself commits and syncs the Resolution before every route (`--consume` folds the approve's consume into the same call)."

## Ideation spike record

Riskiest unverified mechanism (per the seed test plan): folding state commit into `gate record`/`gate consume` without breaking existing gate-lifecycle tests. Spiked 2026-08-02 in an isolated clone at 23ed415bb:

- ~20-line spike seam: after record/consume success (and consume's stale write), invoke the existing `runStateCommit` machinery with the slug derived from the resolved entity path -- pure composition, no gates-package changes.
- Full `go test ./internal/cli` (all 15 gate test functions plus the rest of the package, 60s) green after exactly ONE fixture line: `gatePrepareCLIFixture`'s state checkout needed `git branch -M spacedock-state/dev` -- the fixture sat on `main` while real checkouts are born on the state branch by `state init`/`new`; statesync preflight rightly refuses the mismatch. `internal/gates` and `internal/statesync` untouched and green.
- End-to-end on a split-root fixture: `gate record --decision approve` -> Resolution + path-scoped commit (`state: update task`), no-origin local-only handled; `gate consume` -> status ideation->implementation + consumed mark + second commit; porcelain clean after both.

Second spike round (staff review: "the real new risk -- push reject -> pull --rebase -> conflict HALT, now inside a gate verb -- went unexercised"). Same seam, two-writer origin fixtures mirroring `state_commit_test.go`'s `twoHostStateWorkflow` shape:

- Peer-integration path: peer pushed a DISJOINT entity's commit first; `gate record --decision approve` inside the spiked binary hit the push rejection, ran pull --rebase, replayed its record commit atop the peer's, re-pushed -- exit 0, stdout `Committed task, integrated peers' state, and pushed to spacedock-state/dev.`, state log showing the record commit atop the peer commit.
- Conflict HALT path: peer rewrote the SAME entity's frontmatter and pushed first; `gate record --decision approve` -> record durable, push rejected, rebase CONFLICT -> the verb exited 3 (observed, not asserted) with the full existing HALT rendering on stderr (`HALT -- same-entity rebase conflict on spacedock-state/dev.` / `Conflicting path(s): task.md` / peer commit named / never-force guidance), rebase aborted (no rebase dir, porcelain clean), the Resolution still durable in the local HEAD commit (`decision: approve` present in `git show HEAD:task.md`), and the peer's edit untouched on origin.
- Conclusion: the fold is safe, including the genuinely new risk surface; exit-3 propagation through a gate verb is demonstrated; the one-line fixture alignment is the entire observed breakage surface, and it must be declared under AC-2.

No spike needed for the other two mechanisms -- proven code paths: `--consume` sequences the two handlers the spike exercised in-process; `--stamp` writes through the native `status --set` handlers (existing corpus), the spiked commit seam, and `git worktree add` driven from Go as already shipped in `internal/cli/state.go` (`state init`/`new`).

## Out of scope

Any weakening of gate authority, digest/Briefing integrity checks, frozen-binding refusals, or the captain-vs-FO actor distinction. Any change to `gate prepare`'s semantics or the room-backed Result path. The mechanism-2 byte-cap follow-up (`gate-lifecycle-hardening-byte-budget`) -- unrelated. Batch/bulk gate operations across multiple entities -- this is about one entity's one gate decision.

## Acceptance criteria

**AC-1 (VALUE) - The binary/git command count for one captain "approve" at a nonterminal `gate: true` stage, from open gate through emitted dispatch envelope, drops 10 -> 2 against the transcript baseline.**
Independent baseline: the 16 FO tool calls in `## Problem` (sonnet-gate-guardrail-no-authority transcript). Of those 16, ten are binary/git commands once the incidental retitle `--set` is excluded (`gate record`, `state commit`, `gate consume`, `state commit`, `status --read`, `status --set`, `state commit`, `git status`, `git worktree add`, `dispatch build`); the other five (two checklist/scope-note Writes, ToolSearch, prior-worker shutdown SendMessage, the Agent spawn) are host-side calls outside the binary's reach, so the post-change ceiling on total FO tool calls is <=7.
Verified by: `internal/cli/gate_ceremony_count_test.go`, a fixture harness that drives the compiled command surface twice over the same split-root fixture (gated stage -> worktree-declaring successor): once through the 10-command before-list transcribed from the transcript, once through the after-list pinned to exactly the two documented FO commands (`gate record --decision approve --actor person:captain --consume`, `dispatch build --stamp`) and nothing else. The harness asserts (a) every command exits 0 and (b) both runs converge to equivalent end-state predicates -- Resolution closed and application consumed, status advanced to the successor, `started`/`worktree=` stamped, state-checkout log containing the ceremony commits with a clean tree, worktree present on its branch, spawn envelope emitted. The end-state parity is the falsifiable instrument: if the design regresses so that a required step is missing from the two-command sequence, the after-run's end-state predicates fail (an unstamped worktree, an unsynced commit, a missing envelope), and closing that gap forces a visible third command into the pinned list. The 10 -> 2 delta itself is a property of the two committed command lists the test executes verbatim, not a separate self-referential length assertion.

**AC-2 - No loss of gate authority or integrity guarantees.**
Verified by: the existing gate-lifecycle test corpus (frozen-binding refusal, digest mismatch rejection, actor/authority validation, room-backed vs chat-decision paths, terminal routing) passes with zero production-assertion changes and exactly one declared fixture change -- `gatePrepareCLIFixture`'s state checkout moves to the declared state branch (one line, spike-verified). New negative tests pin the preserved boundaries: `--consume` with `--decision revise|hold` exits 2 before any write; `--stamp` with entity status != `--stage` refuses without mutation; `--stamp --advance` exits 2. Composition is structural, not parallel: `--consume` sequences the existing record and consume handlers; `--stamp` writes through the native `status --set` machinery (inheriting its mutation guards) and the same `state commit` seam.

**AC-3 - The FO contract and user docs carry the collapsed call shape as the primary path.**
Verified by: the `## Doc diff` above is applied to `skills/fo-gate-lifecycle/SKILL.md`, `skills/first-officer/references/fo-dispatch-core.md`, `skills/first-officer/references/first-officer-shared-core.md`, `docs/site/reference/command-reference.md`, and `docs/site/concepts/gates-and-decisions.md`; grep over fo-gate-lifecycle finds no remaining post-close/post-consume `«state.commit»` instruction (remaining `«state.commit»` references live only where the standalone verb is still real: the post-prepare bind commit, `--round` publications, entity seeding, body/report edits, merge archive); fo-dispatch-core's dispatch steps no longer instruct separate stamp/commit/worktree-create commands as the primary path, and its `«dispatch.build»` block clause distinguishes stamp/sync failures and the exit-3 HALT from the break-glass trigger.

## Test plan

- `internal/cli/gate_ceremony_count_test.go` (AC-1): the before/after harness above -- temp git fixtures (main repo + split-root state checkout on `spacedock-state/dev`), no live host, ~250 lines. This reproduces the baseline as a runnable count first (the before-list), then measures the after-list on the same fixture shape, so before/after is a real number.
- Behavior tests (AC-2, mechanisms): implicit sync line and exit semantics per verb -- pushed/local-only/no-op, the peer-integration re-push path, and the exit-3 HALT inside `gate record` via the same two-writer origin/conflict fixtures `state_commit_test.go` already uses (`twoHostStateWorkflow`; both paths seeded by this ideation's spike, see the spike record); each `--consume` landing position from the mechanism-2 table asserted on its final `sync=`/`phase=` line and exit code (happy path, usage-error, room-source revise skip, terminal `route=approved-awaiting-merge`, stale-supersede commit-before-exit, sync-failed and HALT at both phases); `--stamp` stamps+commit+worktree-add, status-mismatch refusal, idempotent re-run (existing worktree, already-set `started`), non-worktree stage (`started` only), `--advance` incompatibility, and the stderr discrimination contract (stamp failures carry the `dispatch build --stamp:` prefix and emit no envelope; assembly failures do not carry it).
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

## Stage Report: ideation (cycle 2)

Staff review returned "approve with changes" (4 fixes + one AC reframe); all five addressed:

- DONE: Narrow mechanism 1 -- drop `gate prepare` from the implicit-sync scope.
  Scope is now `gate record` close sources + `gate consume` only; prepare (and `--round`) recorded as Alternatives rejected 6 with the zero-AC-1-benefit / new-failure-surface rationale; prepare's line-27 block clause thereby stays valid unchanged.
- DONE: Add the missing block-clause and shared-core doc updates.
  Doc diff now covers `first-officer-shared-core.md:117` ("Every state write is one call" amended for self-syncing verbs), the `fo-dispatch-core.md:152` break-glass trigger disambiguation (stamp/sync failures and exit-3 HALT are never break-glass; assembly failures only), and notes shared-core:133's HALT clause binds unchanged. The fo-gate-lifecycle:27 misroute is resolved by dropping prepare sync rather than re-documenting it.
- DONE: Give `--consume` a machine-parseable discriminator.
  The sync line becomes `sync=pushed|local-only|no-op|failed|halted phase=record|consume`, emitted in every split-root outcome; an 8-row landing-position table (final line + exit code) added to mechanism 2; the previously wrong SKILL.md line-48 after-wording rewritten to branch on the final phase/sync line (consume-refusal-after-clean-record-sync now correctly reads as `phase=record` + consume diagnostic).
- DONE: Spike the actual risky path (push reject -> pull --rebase -> conflict HALT inside a gate verb).
  Second spike round, two-writer origin fixtures mirroring `twoHostStateWorkflow`: disjoint-peer path re-pushed after rebase integration (exit 0, "integrated peers' state"); same-entity conflict HALTed with exit 3 observed, full HALT stderr (conflicting path, peer commit, never-force), rebase aborted clean, Resolution durable in local HEAD, peer edit preserved on origin. Recorded in the spike record.
- DONE: Reframe AC-1's verification without the self-referential length assertion.
  Assertion (c) dropped; verification now rests on end-state parity between the pinned 2-command after-list and the 10-command transcript before-list -- a missing required step fails the parity predicates and forces a visible third command into the pinned list.

### Summary

Revised per staff review: implicit sync narrowed to the two verbs inside AC-1's measured window, failure routing disambiguated in all three prose surfaces the review named, `--consume` given a branchable `sync=`/`phase=` discriminator with a full landing-position table, and the genuinely new risk (rebase-conflict HALT inside a gate verb) demonstrated with exit-3 propagation shown rather than asserted. The `gate approve`/`dispatch advance` rejections stand as reviewed.
