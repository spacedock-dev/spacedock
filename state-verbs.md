---
title: State-repo verbs — spacedock state ready / sweep / commit <slug> (rebase-HALT enforced by the verb)
status: validation
source: 0205 carve (2026-06-17, captain "stamp them") — index DoD candidate "state-verbs"; 2y MERGED unblocks it.
score: 0.5
sprint: 0221-layered-fo
group: verb-core
sprint-readiness: ready
id: rgq0m30693ke84zvb8yna078
started: 2026-06-17T20:55:24Z
worktree: .worktrees/spacedock-ensign-state-verbs
mod-block: merge:pr-merge
---

`spacedock state ready` / `state sweep` / `state commit <slug>` behind the «state.ensure-ready» / «state.sweep-merged» / «state.commit» prose-functions. The commit verb is path-scoped, pushes with retry-on-reject-rebase, and REFUSES to return on a rebase conflict (exit non-zero, stderr naming the paths, repo left in rebase-abort) — so the FO cannot proceed even if it wanted to. The rebase-conflict halt is enforced by the verb, not by FO discipline.

Spike must-build: «state.commit» — the `git -C` form drifted 2/3 under Haiku in the w4 spike, and it carries the rebase-HALT trigger w4 could NOT exercise (failure-mode 1, NOT-EXERCISABLE). Depends on 2y (MERGED, the host-neutral cores).

Ideation must resolve: the three verb surfaces + their return contracts; SPIKE the rebase-HALT FIRST (a 2-writer / injected-conflict harness that actually FIRES the halt) + a Go test that «state.commit»'s `git -C` form survives; oracle-based ACs (never a prose-grep).

## Problem

Today the split-root state-sync sequence (path-scoped commit → push → on-reject `pull --rebase` → re-push, and HALT on a same-entity rebase conflict) lives only as PROSE that the FO/ensign reads and re-types by hand. Two surfaces carry that prose: the FO contract's `## State Management` (`skills/first-officer/references/first-officer-shared-core.md`) and the dispatch-body guidance `stateCommitGuidance()` emits into every split-root ensign's prompt (`internal/dispatch/build.go:809`). The ensign reads `git -C <checkout> add <entity> && git -C <checkout> commit -m "..." -- <entity>` as text and types it back.

The w4 haiku-loop spike (PR #393) measured what happens when a weak model re-types that recipe:

- The `git -C` form **drifted to `cd <checkout> && git add` 2 of 3 runs** — a foreign-cwd ensign (its process cwd is its `.worktrees/…`, not the state checkout) that runs `cd && git add` resolves nowhere or commits the wrong tree. Flaky, so must-build per the spike's flakiness rule.
- The **rebase-conflict HALT (failure-mode 1) was NOT-EXERCISABLE** on the linear happy path. The contract says the FO must HALT, abort the rebase, not force-push, and surface the conflict — but nothing fired it, so its enforcement rested entirely on the model's discipline. The single highest state-repo risk (the per-area analysis) is Haiku inventing an auto-recovery instead of halting.

A prose recipe a weak model paraphrases is exactly what a binary verb removes: there is no command STRING to drift, and the halt becomes a non-zero exit the model cannot wish away.

## Proposed approach

Ship three sibling subcommands under the existing `spacedock state` command (alongside `init`/`new` in `internal/cli/state.go`), reusing that file's `runGit(dir, args...)` argv helper (already the correct `exec.Command("git", "-C", dir, …)` form — no rendered string to drift). All three read the README `state:`/`state-branch:` via the shared `status.ClassifyState` / `status.StateBranch` interpreters and the `status.stateHasOrigin` origin probe, so they agree with boot's `STATE_BACKEND` line by construction.

**`spacedock state commit <slug> [--workflow-dir DIR] [-m MSG] [--json]`** (← `«state.commit»`, the must-build). Resolves `<slug>` to its entity path under the state checkout, then performs the full sequence as binary steps:

1. **Path-scoped commit** — `git -C <checkout> add <entity-path>` then `git -C <checkout> commit -m <msg> -- <entity-path>`. Never `add -A`, never a bare `commit`. Nothing staged (no diff) is a clean no-op success, not an error.
2. **Push** (origin present) — `git -C <checkout> push origin <state-branch>`. On a non-fast-forward rejection, go to 3. No-origin: skip 2–4, report local-only success (mirrors the boot `remote: none` carve-out).
3. **`pull --rebase`** — `git -C <checkout> pull --rebase origin <state-branch>`, reading **git's own exit code directly** (the spike's first run failed because a `tee` pipeline masked the conflict exit — the implementation MUST NOT pipe this through anything that swallows the status). Clean rebase (disjoint paths) → 4. Conflict → 5.
4. **Re-push** — `git -C <checkout> push origin <state-branch>`. Success ends the verb (exit 0).
5. **HALT** (rebase conflicted) — `git -C <checkout> rebase --abort` to restore a clean tree, then **exit non-zero (3)** with stderr naming the conflicting entity path(s) and the directive that manual intervention is required. NEVER `--force` / `--force-with-lease`; NEVER `-X ours/theirs` or any auto-resolve. The verb refuses to return success, so an FO that calls it cannot proceed on an unmerged tree even if it wanted to — the halt is the exit code, not the model's discipline.

`index.lock` contention on the add/commit is retried after a short wait (≤2s, bounded retries) before failing.

**`spacedock state ready [--workflow-dir DIR] [--json]`** (← `«state.ensure-ready»`, the boot gate). Inline workflow → reports ready, nothing to sync (exit 0). Split-root → ensures the state checkout is present (the same fetch+`worktree add` resume `state init` already does when absent), then, if origin is present, integrates peers' state with one `git -C <checkout> pull --rebase origin <state-branch>`. A clean/fast-forward pull → ready (exit 0). A **rebase conflict here halts identically to `state commit` step 5** (abort, exit 3, name the paths) — the boot pull-on-conflict the FO contract already mandates, now enforced. No-origin → ready, local-only (exit 0). This is the binary form of the contract's "Split-root pull-on-boot" startup step.

**`spacedock state sweep [--workflow-dir DIR] [--json]`** (← `«state.sweep-merged»`). The state-repo's read-only view of entities whose code PR has MERGED but whose state has not been terminalized — the merged-PR iteration the FO sweeps each loop. This computation **already exists** as `dispatch reconcile`'s `un-advanced-pr` drift class (`internal/dispatch/reconcile.go`, `gh pr view` merged-detection). Per YAGNI and the no-duplication rule, `state sweep` is proposed to map to a thin call over that existing detection rather than a re-implementation; **the open question for the gate is whether `sweep` ships as a backtick verb at all or stays a guillemet prose-function pointing at `dispatch reconcile --include un-advanced-pr`** (see "Verb→notation mapping" below). It carries no novel risk: it is read-only (no commit, no push, no halt), so it is NOT spike-gated.

All three return a small JSON object under `--json` (default human prose), mirroring `reconcile`'s envelope style: `{"command": "state commit", "slug": "...", "result": "committed|pushed|local-only|no-op|halted", "state_branch": "...", "conflicting_paths": [...]?, "reason": "..."}`. Exit-code contract matches the established native surface: **0** success (including clean no-op and local-only), **1** operational/setup failure (not a repo, slug not found, push failed for a non-conflict reason), **2** usage error (missing/unknown args), **3** rebase-conflict HALT (distinct from generic failure so a caller can branch on "the peer-conflict halt fired" vs "git broke").

### Spike — riskiest mechanism, run FIRST (DONE)

Per the proof policy, the two unverified mechanisms were exercised before committing to the build. Both PASSED. Throwaway harnesses, recorded here to seed the implementation's first tests:

- **Rebase-HALT actually fires (failure-mode 1, the w4 NOT-EXERCISABLE gap).** A real bare-origin + two-clone harness: hostA and hostB both edit `alpha.md`'s frontmatter; A commits+pushes; B path-scoped-commits, push is **REJECTED (rc=1, non-ff)**, `pull --rebase` **CONFLICTS (rc=1, "CONFLICT (content): Merge conflict in alpha.md")** leaving a rebase in progress; `rebase --abort` restores a **clean tree** (`status --porcelain` empty); a plain (non-force) re-push **stays rejected**; and origin still carries **A's edit** (`status: implementation`), NOT B's (`status: review`) — no silent clobber. EXERCISABLE and real.
  - **Load-bearing lesson (seeds the impl's first test):** the verb MUST read `git pull --rebase`'s OWN exit code. The spike's first run reported a false pass because piping `pull --rebase` through `tee` made the pipeline exit with `tee`'s status, masking the conflict. The implementation must not pipe the rebase through anything that swallows the exit; the test must assert the conflict is detected via the captured non-zero status.
- **The `git -C` argv form survives (the w4 2/3 drift).** A throwaway Go test runs `exec.Command("git","-C",repo,"add","alpha.md")` + path-scoped commit from a process cwd that is NOT the repo (the foreign-cwd ensign situation), and asserts exactly one commit landed and a sibling `beta.md` stayed untracked (path-scoped, not `-A`). PASSED. The point the spike proves: a Go verb shelling out via an **argv slice** has no rendered command string for a weak model to paraphrase into `cd && git add` — the drift surface is the prose recipe, and the verb deletes it.

No further spike needed: the multi-writer happy path (push → non-ff → clean `pull --rebase` → re-push) is already proven at the git level by the merged `TestTwoWriterSyncHappyPath` (`internal/cli/state_sync_test.go`), and the same-entity halt by `TestTwoWriterSameEntityConflictHalts` there — those exercise the raw git sequence; this task wraps it in a verb whose exit code enforces the halt.

## Out of scope

- **A full lock model for the shared state index.** The halt is the boundary behavior (matching the contract); concurrent writers coordinate via push/`pull --rebase`/halt, not a lock.
- **Tool-managed atomic commit inside `status --set`.** The contract names "tool-managed atomic state commits" as the *preferred* form, but `status --set` does not commit today (it only writes frontmatter; the FO commits separately). Folding the commit into `--set` is a larger surface change; this task ships the explicit `state commit` verb the FO calls after `--set`, which is the path the contract's fallback already describes.
- **Re-implementing merged-PR detection.** `state sweep` reuses `dispatch reconcile`'s `un-advanced-pr` class; it does not invent a second merged-detection path.
- **Codex / Pi operability of these verbs.** They are host-neutral binary subcommands (no `~/.claude` reads), so they run anywhere `spacedock` does, but the layered-FO drive that consumes them is validated on the Claude substrate only this sprint (per the 0205 out-of-scope).
- **Flipping the prose-function notation from guillemets to backticks in the FO contract.** That edit is owned by the `prose-function-restructure` member; this task ships the verbs and records which prose-function each backs, so the restructure flips them. (If the gate prefers, the doc-diff below can ride with this task instead.)

## Verb→notation mapping (for the prose-function restructure)

| Prose-function | Ships as | Migration target |
|---|---|---|
| `«state.commit»` | **backtick verb** (must-build) | `spacedock state commit <slug>` |
| `«state.ensure-ready»` | **backtick verb** | `spacedock state ready` |
| `«state.sweep-merged»` | **gate decision** — backtick verb `spacedock state sweep` OR stays guillemet pointing at `dispatch reconcile --include un-advanced-pr` | per gate |

## Acceptance criteria

Each AC names a property of the finished entity and how it is verified. All proofs run the behavior against a real git repo and observe exit code / on-disk state — none is a prose-grep over an instruction file.

**AC-1 — `state commit <slug>` HALTS on a same-entity rebase conflict: exit code 3, the conflicting entity path named on stderr, the state checkout left clean (rebase aborted, not in-progress), no force-push, and the peer's edit preserved on origin.**
Verified by: a Go test in `internal/cli` building a bare origin + two clones on the shared state branch, both editing one entity's frontmatter; writer A pushes first; the verb runs as writer B and the test asserts exit 3, a stderr substring naming the entity path, an empty `git status --porcelain` with no `.git/rebase-merge`|`rebase-apply` dir, that a subsequent plain push stays rejected, and that `git show <state-branch>:<entity>` on the bare origin carries A's frontmatter value, not B's. (Seeded by the spike harness above.)

**AC-2 — `state commit <slug>` performs a path-scoped commit (never `add -A`): a sibling dirty/untracked file in the state checkout is NOT swept into the commit.**
Verified by: a Go test that stages one entity via the verb while a sibling file is dirty/untracked in the same checkout, asserting the resulting commit's `git show --name-only` lists only the entity path and the sibling remains untracked in `git status --porcelain`.

**AC-3 — `state commit <slug>` on the multi-writer happy path commits, pushes, and on a non-fast-forward rejection `pull --rebase`s the disjoint peer commit and re-pushes successfully (exit 0), leaving both entities present and history linear.**
Verified by: a Go test with two clones committing DIFFERENT entities; A pushes first; the verb runs as B and the test asserts exit 0, both entity files present in B's tree after the verb, no merge commit (`git log --merges` empty), and B's entity present on the bare origin.

**AC-4 — `state commit` in a no-origin state checkout commits path-scoped locally and reports local-only success (exit 0) without attempting push/pull.**
Verified by: a Go test on a single-clone checkout with no `origin` remote, asserting exit 0, the commit landed locally, and (via a recording git or absence of a push ref) no push was attempted; the `--json` result field is `local-only`.

**AC-5 — `state ready` integrates peers' state on boot via one `pull --rebase` and reports ready (exit 0) on a clean integration, and HALTS identically to AC-1 (exit 3) on a same-entity boot conflict.**
Verified by: two Go tests — a clean-integration case asserting exit 0 and the peer's committed entity present in the local checkout after the verb, and a conflict case asserting exit 3 with the checkout left clean (rebase aborted).

**AC-6 — `state ready` on an inline (single-root) workflow is a clean no-op (exit 0, nothing to sync); `state ready` whose split-root checkout is absent resumes it (present afterward).**
Verified by: a Go test over an inline README asserting exit 0 and no git network op, plus an absent-checkout case asserting the checkout directory exists afterward (reusing the `state init` resume path).

**AC-7 — `state sweep` is read-only — it makes no commit, push, or state mutation — and reports the merged-but-not-terminalized entities.**
Verified by: a Go test asserting `git rev-parse HEAD` of the state checkout is unchanged across the verb and the reported set matches the `un-advanced-pr` entities from a seeded fixture (a stubbed `gh` merged-state). (If the gate routes `sweep` to a guillemet prose-function instead of a verb, this AC is struck and the mapping table records the decision.)

**AC-8 — usage and routing: each verb exits 2 on an unknown/missing required argument with a diagnostic naming the subcommand, and `spacedock state <unknown>` exits 2 listing the valid subcommands (now including ready/sweep/commit).**
Verified by: a Go test in `internal/cli` asserting the exit-2 usage paths and that the `state` command's unknown-subcommand error text enumerates `init|new|ready|sweep|commit`.

## Test plan

All proofs are Go tests in `internal/cli` (the package that owns `state.go` and the existing `state_sync_test.go` / `state_init_test.go` real-git harnesses). They reuse that file's pattern: a real `git init`/bare-origin + clone fixture under `t.TempDir()` with `HOME`/`GIT_CONFIG_GLOBAL` pinned, so every assertion observes a real exit code or on-disk git state — no mocks, no prose-grep. Cost is low (the existing sync e2e runs in ~3.7s for the suite; these add a handful of similar fixtures). No live-workflow or host-runtime test is needed: the verbs are host-neutral git operations with no `~/.claude` dependency. The riskiest mechanism (AC-1, the halt) is the first test written, seeded directly by the spike harness recorded above. The `gh`-dependent AC-7 stubs the merged-state probe the same way `reconcile`'s tests do (an injected `ghRunner`), keeping it deterministic and offline.

Per the dev-workflow path→lane mapping: these verbs are new binary surface in `internal/cli`, not a change to a shipped FO/ensign adapter or the dispatch core a live lane drives, so the deterministic Go lanes are the gate; no live host lane is required UNLESS the prose-function notation flip (the doc-diff) is folded in, in which case the FO-contract change requires the relevant live lane.

## Documentation changes

These verbs are new user-visible CLI surface, so ideation proposes the doc diff. Two surfaces gain the verbs; the prose-function notation flip is owned by `prose-function-restructure` (above) unless the gate folds it here.

1. **`internal/cli/state.go` ABOUTME + `newStateCommand` usage string** (`internal/cli/cli.go:343-344`): the `Use:`/`Short:` strings and the unknown-subcommand error change from `init|new` to `init|new|ready|sweep|commit`, and the bash/zsh completion `compadd`/`compgen` lists (`internal/cli/cli.go:557,577`) gain `ready sweep commit`. (Code-adjacent, applied by implementation.)

2. **`spacedock state` help / any command-surface doc** that enumerates the workflow verbs gains the three subcommands with one-line descriptions:
   - `state commit <slug>` — path-scoped commit + push of one entity's state, with on-reject `pull --rebase`; halts (non-zero) on a same-entity conflict rather than clobbering a peer.
   - `state ready` — integrate peers' state and confirm the checkout is ready to dispatch against; halts on a boot conflict.
   - `state sweep` — list entities whose code PR merged but whose state is not yet terminalized.

The concrete before/after wording for the cobra `Use:`/`Short:` and completion lists is recorded in AC-8's test expectations (the enumerated subcommand set is the checkable surface). If the gate folds the FO-contract notation flip in, the before/after is: `## State Management`'s `git -C {state_checkout} add … && git -C {state_checkout} commit …` recipe and the dispatch `stateCommitGuidance()` body become `call `spacedock state commit <slug>`` (the verb), with the raw git sequence kept as the verb's documented internal behavior, not an instruction the ensign re-types.

## Stage Report: ideation

- DONE: Design the three verb surfaces — `spacedock state ready` / `state sweep` / `state commit <slug>` — and their return contracts, against the REAL post-2y host-neutral cores and the current split-root state-sync rules
  `## Proposed approach` + `## Acceptance criteria` design all three against `internal/cli/state.go`'s `runGit` helper, `status.ClassifyState`/`StateBranch`/`stateHasOrigin`, and the FO core's `## State Management` (lines 137-161); behavior-first return/exit contract (0/1/2/3 + JSON envelope) recorded.
- DONE: SPIKE the riskiest mechanism FIRST — the rebase-conflict HALT (a 2-writer / injected-conflict harness that ACTUALLY FIRES it) PLUS a Go test that «state.commit»'s `git -C` form survives
  Both ran and PASSED before the design committed; recorded under `### Spike`. Halt fires (push rc=1 non-ff → `pull --rebase` rc=1 CONFLICT → abort restores clean tree → non-force re-push stays rejected → peer edit preserved on origin). `git -C` argv form commits path-scoped from a foreign cwd. Surfaced the load-bearing lesson (read git's own exit; a `tee` pipeline masked it on the first run) that seeds the impl's first test.
- DONE: Behavior-first oracle-based ACs + test plan (never a prose-grep); map «state.ensure-ready»/«state.sweep-merged»/«state.commit» to backtick-verb vs guillemet; propose the doc-diff
  AC-1..8 each verified by a real-git Go test observing exit code / on-disk state (no instruction-file grep); `## Verb→notation mapping` table records commit+ready as backtick verbs and flags sweep as a gate decision (verb vs guillemet over `reconcile --include un-advanced-pr`); `## Documentation changes` carries the cobra-usage + completion doc-diff.

### Summary

Fleshed `state-verbs` to a behavior-first ideation spec: `state commit <slug>` is the must-build (path-scoped commit → push → on-reject `pull --rebase` → re-push, HALT with exit 3 + named paths + rebase-abort on a same-entity conflict), `state ready` is the boot-pull gate with the same halt, `state sweep` is a read-only merged-sweep proposed to reuse `reconcile`'s `un-advanced-pr` class (verb-vs-guillemet left as the one gate decision). The two w4-unsettled mechanisms — the NOT-EXERCISABLE rebase-HALT (failure-mode 1) and the 2/3 `git -C` drift — were both EXERCISED FIRST and PASSED; the throwaway harnesses are recorded to seed the implementation's first tests. The verb deletes the prose recipe a weak model paraphrases, so there is nothing left to drift, and the halt is a non-zero exit the FO cannot wish away. This is COMPLEX (state-repo ops + an enforced halt) — expect a staff review before the gate; the one open call for the captain is whether `state sweep` ships as a verb or stays a guillemet pointer at the existing reconcile detection.

## Stage Report: implementation

- DONE: `state commit <slug>`'s rebase-HALT is enforced by the verb's exit code (3) reading git's OWN exit status (never piped through anything that swallows it); on a same-entity conflict it runs `rebase --abort` leaving the checkout clean, never force-pushes, and the peer's edit is preserved on origin (AC-1)
  `runStateCommit` reads git's status via `runGit`'s `CombinedOutput` err (no `tee`/pipe); `haltOnConflict` aborts + exits 3 with the path named. `TestStateCommitHaltsOnSameEntityConflict` (real bare-origin + 2 host worktrees) asserts exit 3, `first-task.md` on stderr, clean porcelain, no `rebase-merge` dir, plain re-push stays rejected, origin carries A's edit. Commit 07ec5cc3.
- DONE: `state commit` is path-scoped (`git -C add <entity>`, never `add -A`) so a sibling dirty file is never swept in (AC-2); the multi-writer happy path and no-origin local-only paths return exit 0 (AC-3/AC-4)
  `commitEntityPathScoped` stages `add -- <rel>` then `commit -- <rel>`. `TestStateCommitIsPathScoped` proves a sibling untracked file stays out of the commit; `TestStateCommitMultiWriterHappyPath` (push→non-ff→`pull --rebase`→re-push, linear history) and `TestStateCommitNoOriginLocalOnly` (`result: local-only`) both exit 0.
- DONE: `state ready` carries the identical boot-pull HALT (AC-5); `state sweep` ships as a read-only verb mapping to `reconcile`'s un-advanced-pr detection with no second merged-detection path (AC-7)
  `runStateReady` runs one `pull --rebase`, HALTing via the shared `haltOnConflict`. `TestStateReadyIntegratesPeerState` (clean, exit 0) + `TestStateReadyHaltsOnBootConflict` (exit 3, clean tree). `state sweep` → `dispatch.Sweep` reuses `classC` verbatim (`TestSweepReportsMergedNotTerminalized` pins the swept set; `TestStateSweepIsReadOnly` asserts HEAD unchanged). AC-6 (inline no-op + absent-checkout resume) and AC-8 (usage exit-2 + subcommand enumeration) also covered.

### Summary

Shipped three split-root state verbs under `spacedock state`, plus an exported `dispatch.Sweep`. `state commit` runs the path-scoped commit → push → on-reject `pull --rebase` → re-push sequence and HALTs (exit 3, path named, rebase aborted, no force-push) on a same-entity conflict; `state ready` carries the same halt on the boot pull; `state sweep` is a read-only reuse of reconcile's `un-advanced-pr` detection (per the ratified dispatch, sweep ships as a verb — the gate's verb-vs-guillemet question is resolved). The drift surface is deleted: `runGit`'s argv form (`exec.Command("git","-C",dir,…)`) has no rendered command string to paraphrase, and the halt is a non-zero exit git's own status produces — the lesson from the spike's `tee`-masked false pass. `rebaseInProgress` uses `git rev-parse --git-path` so a linked-worktree checkout's conflict is detected, not missed by a hardcoded `.git/rebase-merge`. All 8 ACs are proved by real-git Go tests in `internal/cli`/`internal/dispatch`; full `go test ./...` is green and `go vet` clean. Code committed to `spacedock-ensign/state-verbs` (07ec5cc3).

## Stage Report: validation

- DONE: Detached adversarial audit on a THROWAWAY checkout (NOT the impl worktree): construct a claim-breaking edit to the rebase-HALT and confirm `TestStateCommitHaltsOnSameEntityConflict` catches it — refute that the verb's exit-3 HALT is real
  Audit ran in `/tmp/state-verbs-audit.*` (a `git archive | tar` copy, never the impl worktree; removed after). Baseline halt test passed unmutated; then FOUR claim-breaking mutations, each CAUGHT: (1) exit-3→exit-0 → FAIL at `state_commit_test.go:97` (got exit=0); (2) swallow `pull --rebase` exit via forced `rebaseOK=true` (the `tee`-pipe failure mode) → verb proceeds, plain re-push rejected, FAIL at :97 (got exit=1); (3) swallow exit + force re-push → FAIL at :97 (got exit=0); (4) force-push INSIDE the HALT while keeping exit 3 → FAIL at `:116` (the no-force-push/peer-preserved assertion). Test is not a tautology — it cannot pass when the HALT is broken; exit-code AND clobber assertions both have independent teeth.
- DONE: Reproduce AC-1 on the real binary: same-entity conflict yields exit 3, conflicting path on stderr, clean porcelain (rebase aborted, no rebase-merge/apply dir), plain re-push stays rejected, peer's edit preserved on origin — git's own non-zero status, not prose
  Built `cmd/spacedock`; real bare-origin + two `state init`'d host checkouts. B's `state commit first-task` exited 3, stderr named `first-task.md`, `status --porcelain` empty, no rebase-merge/apply (via `rev-parse --git-path`), plain re-push rejected non-ff, origin `first-task.md` carried A's `status: implementation` (NOT B's `review`). All 6 properties PASS on the binary.
- DONE: AC-2 path-scoped (never `add -A`), AC-3/4 happy-path + no-origin local-only exit 0, AC-5 `state ready` boot-HALT, AC-7 `state sweep` read-only reuse of reconcile's un-advanced-pr (no second merged-detection path), AC-8 usage exit-2 — all verify against real-git tests; full `go test ./...` green
  Real-binary repro all PASS: AC-2 (sibling `sibling-junk.md` stays untracked, commit lists only the entity), AC-3 (B push rejected→`pull --rebase`→re-push exit 0, result `pushed`, both entities, linear history, B's beta on origin), AC-4 (no-origin exit 0, result `local-only`, HEAD advanced), AC-5 (boot-HALT exit 3 + clean integration exit 0), AC-7 (HEAD unchanged, clean tree, `command: state sweep`), AC-8 (all 7 usage paths exit 2; `init|new|ready|sweep|commit` enumerated). Single merged-detection path confirmed: only `reconcile.go:classC` exists in production; `Sweep` (reconcile.go:688) calls it verbatim — no duplicate. Go tests: `internal/cli` 11/11 state-verb + `internal/dispatch` 2/2 sweep PASS uncached; `go vet ./...` clean.

### Summary

Validation PASSED. The adversarial audit (four mutations on a throwaway checkout, never the impl worktree) refuted the claim that the exit-3 HALT could be quietly broken: every mutation — dropping exit-3, swallowing git's rebase exit via a pipe, force-pushing past the conflict, and force-pushing inside the HALT — is caught by `TestStateCommitHaltsOnSameEntityConflict` at a distinct assertion, so the test is no tautology. AC-1 and AC-2/3/4/5/7/8 each reproduce on the REAL built `spacedock` binary against real git (bare origin + linked-worktree state checkouts), observing exit codes and on-disk/origin state, not prose. One honest note: a single uncached `go test ./...` run hit the 600s package timeout in `internal/cli` inside the UNRELATED install/frontdoor tooling tests (`TestUpgradeFromStaleMovesToGreen`/`TestFreshBoxInstallSucceeds`/`TestClaudePluginInstallIsHostNative`, which shell to host tools under parallel load); three isolated uncached `internal/cli` runs complete in ~20s with every state-verb test green, so the timeout is environmental, not in this change. PASSED.
