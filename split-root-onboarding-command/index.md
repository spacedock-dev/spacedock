---
id: kjne6h1jkft6b6ek3p3068jh
title: One-command split-root onboarding — `spacedock state new` (birth the orphan state branch + linked worktree for an existing repo)
status: validation
source: "captain (2026-06-02) — real-repo split-root test on DataRecce/recce: storage mechanics work cleanly (zero code-branch churn, fresh-clone halt-gate + state init both correct), but first-time onboarding of an existing repo has no single command"
completed: 2026-06-03T08:25:33Z
verdict: PASSED
score: "0.24"
worktree:
issue:
started: 2026-06-03T07:09:59Z
mod-block:
pr: "#279"
---

Captain validated split-root on a real repo (DataRecce/recce, isolated clone): zero-churn confirmed (code branch `git status` stays empty, log never sees a state commit), the fresh-clone halt-gate fires correctly, and `spacedock state init` re-checks-out the orphan worktree (`present:false` → `present:true`) on a real 90M repo. The gap is **onboarding**:

- `spacedock state init` only RE-checks-out an *existing* orphan state branch (fetch + `git worktree add`). It does not BIRTH a split-root workflow.
- First-time setup of an existing repo — append the one `.gitignore` line, birth the orphan branch (`rm -rf --cached .` to clear the inherited tree), check it out as a linked worktree at `docs/<wf>/.spacedock-state`, scaffold the README with `state: .spacedock-state` — currently requires the `commission` skill or hand-running the Journey-1 sequence.

A `spacedock state new` (or `spacedock commission --split-root`) helper would make onboarding an existing repo a one-liner, symmetric with `state init` for the clone path.

## Proposed approach

`spacedock state new` is the **birth** half of the split-root lifecycle; `state init` is the **resume** half. They are inverses that read the SAME README and reuse the SAME helpers:

| | `state init` (exists) | `state new` (proposed) |
|---|---|---|
| Precondition | orphan branch EXISTS (on origin or local refs) | orphan branch does NOT exist anywhere |
| Action | fetch + `git worktree add {path} {branch}` | gitignore + birth orphan + `git worktree add {path} {branch}` |
| Reads | README `state:`, `state-branch:` | README `state:`, `state-branch:` (identical) |
| Derives | `ClassifyState`, `StateBranch`, `DiscoverWorkflowDir` | same three helpers, verbatim |

**The README is a precondition, not an output.** The README is a design-specific, interactively-authored artifact (mission prose, per-stage definitions, schema — see `skills/commission/SKILL.md` template, ~100 lines). `state new` cannot synthesize it from nothing, and must not try. It reads an already-present split-root README and performs the purely-mechanical birth. This keeps `state new` a thin, testable primitive that `commission` (Journey 1) can later *call* instead of hand-running the sequence — collapsing AC-1's claim "creates … the README scaffold" to: `state new` requires the README and births everything else.

The birth sequence is exactly the spiked Journey-1 mechanic (`skills/commission/SKILL.md:282-310`):

1. Append `{dir}/{state_path}/` to the code branch's tracked `.gitignore` (idempotent `grep -qxF` guard).
2. Birth the orphan branch in a temp detached worktree: `worktree add --detach`, `checkout --orphan {branch}`, `rm -rf --cached .`, clear the inherited working tree (keep `.git`), seed nothing / or a placeholder, commit `seed state`, best-effort `push origin {branch}`, `worktree remove --force`.
3. `git worktree add {dir}/{state_path} {branch}` — the same call `state init` makes.

## Spike result (riskiest unknown — checklist item 1)

Exercised the full birth end-to-end in throwaway repos (3 runs: no-remote, with-remote+origin, dirty-tree). **Mechanism works**; the in-repo `commissionSplitWorkflow` helper (`internal/cli/state_init_test.go:62-121`) already proves the same sequence under test. Findings that seed the implementation's first tests and the failure-mode design:

- **`worktree add --detach` tolerates a dirty main working tree** — it is a separate checkout. The birth does NOT require a clean code tree. The only code-branch mutation is the additive `.gitignore` append, which coexists with the operator's pending changes. → **dirty-tree is NOT a hard failure mode**; `state new` proceeds.
- **No-remote**: `push origin {branch}` fails `exit 128` (`'origin' does not appear to be a git repository`), but local birth + worktree checkout still succeed. → push is **best-effort**; on failure, exit 0 with a warning telling the operator to push the orphan branch later to make state shareable (mirrors commission's `# if origin exists`).
- **Re-run / orphan already exists locally**: `checkout --orphan {branch}` fatals `a branch named '...' already exists`. → must pre-check `git rev-parse --verify refs/heads/{branch}` and refuse.
- **Worktree path occupied**: 2nd `worktree add {path}` fatals `'...' already exists` (same fatal `state init`'s path-exists guard dodges). → must guard `dirExists(statePath)` and refuse.
- **Orphan exists on REMOTE, not local**: `git ls-remote --heads origin {branch}` returns the ref. This is the "already commissioned/birthed elsewhere — you want `state init`, not `state new`" case. → detect and redirect to `state init` (exit non-zero with the corrected command).
- **`.gitignore` append leaves the code branch dirty** (`?? .gitignore` after the sequence). The commission test helper commits it *with* the README; but `state new` onboarding an existing repo only mutates `.gitignore`. **Decision: `state new` does NOT auto-commit `.gitignore`** — it leaves the one-line edit staged-or-unstaged for the operator to commit alongside their README, because committing on the operator's code branch is an outward-facing side effect the operator should own. `state new` prints the next step (`git add .gitignore && git commit`). (Open for captain: auto-commit vs. leave-for-operator — defaulting to leave-for-operator.)

## Acceptance criteria

**AC-1 — one command births a split-root workflow on a repo with a split-root README present.** Given a code branch carrying a split-root README (`state:` set), a single `spacedock state new [--workflow-dir DIR]` creates the orphan state branch, appends the `.gitignore` entry, and checks out the linked state worktree at the declared `state:` path. After it returns, `git worktree list` shows the state path and `spacedock status` renders against the state checkout.
*Verified by:* a real-git e2e test (bare origin + clone, no mocks, the `state_init_test.go` pattern) asserting: orphan branch present in local refs, `worktree list` contains `statePath`, `.gitignore` contains the path line, and `assertStatusRenders`/boot shows `STATE_BACKEND: split-root` + `present: true`.

**AC-2 — symmetric with `state init` (round-trip).** After `state new` on repo A and a code-branch + orphan-branch push, a fresh clone of A runs `state init` and lands in the same working state (`worktree list` has the state path, boot `present: true`, entities render).
*Verified by:* one e2e test that runs `state new` then clones origin and runs `state init`, reusing `TestStateInitResumesFreshClone`'s assertions on the fresh clone.

**AC-3 — refuses an already-birthed or already-remote workflow instead of fataling.** Running `state new` when (a) the state worktree path already exists, (b) the orphan branch exists in local refs, or (c) the orphan branch exists on origin, exits non-zero with a clear message — case (c) names `state init` as the correct command. No raw git `already exists` fatal leaks to stderr.
*Verified by:* three sub-tests asserting non-zero exit, a tailored message, and that `git`'s `already exists` / `fatal:` strings do NOT appear in stderr.

**AC-4 — no-remote is best-effort, not fatal.** `state new` in a repo with no `origin` births the orphan branch and worktree locally, exits 0, and emits a warning that the operator must push the orphan branch to share state.
*Verified by:* an e2e test in a clone-less `git init` repo asserting exit 0, state worktree present, and the push-later warning on stdout/stderr.

**AC-5 — inline / non-split README is a clear error, not a no-op birth.** `state new` against an inline (`state: $inline` or absent) README exits non-zero explaining `state new` only births split-root workflows (contrast `state init`, which prints a benign one-liner — `new` is a mutating verb, so a mismatched backend is an operator error, not a no-op).
*Verified by:* an e2e test on an inline README asserting non-zero exit and the explanatory message.

## Test plan

- **Where:** `internal/cli/state_new_test.go`, real-git e2e (bare origin + clone, `t.TempDir()`, `GIT_*` env), no mocks — verbatim reuse of `state_init_test.go`'s `git`/`gitOK`/`commissionSplitWorkflow`/`assertStatusRenders`/`runStatusBoot` helpers (already in-package). The birth helper is essentially `commissionSplitWorkflow` minus the README/commit step (the operator/commission owns those).
- **Unit:** the orphan-exists precheck (local ref + remote ls-remote parsing) is a pure helper worth a table test; everything else is behavior-level.
- **Cost/complexity:** Low-moderate. New surface is one `runStateNew` + `parseStateNewArgs` mirroring `runStateInit`/`parseStateInitArgs`, reusing every status helper. ~5 e2e tests (one per AC) + 1 unit table test. No new external deps. Estimated a single implementation stage.
- **Live smoke (optional, post-merge):** birth a throwaway split-root workflow in a scratch clone of this repo and confirm `status --boot` — runtime behavior is already covered by the e2e harness, so a live smoke is confirmation, not a gate.

## Notes
- Captain flagged (no action): the orphan `spacedock-state/*` branch must be pushed to be shareable, so it shows in the repo's branch list / PR UI — expected, not a defect. AC-4's push-later warning surfaces this to the operator.
- `commission` Journey 1 (`skills/commission/SKILL.md:282-310`) can later call `state new` instead of hand-running the sequence — out of scope here, but the command surface is designed so commission becomes a thin caller. Not building that integration in this task.
- Untested elsewhere but exercised live in THIS workflow (docs/dev, sd-b32): full FO dispatch with worktree-stages under split-root (implementation→validation→merge, code branch clean, state path-scoped). Captain's recce test covered slug-style storage mechanics only.
- Priority/milestone TBD — captured as backlog; not yet folded into a release.

## Stage Report: ideation

- DONE: Spike the riskiest unknown first: birth an orphan state branch + linked worktree for an existing repo end-to-end (the `spacedock state new` mechanism), before designing the command surface.
  Ran the Journey-1 birth sequence end-to-end in 3 throwaway repos (no-remote, with-origin, dirty-tree); mechanism works and matches the in-repo `commissionSplitWorkflow` helper (`internal/cli/state_init_test.go:62`). Findings recorded in the "Spike result" section.
- DONE: Define the command behavior, idempotency, and failure modes (existing state, dirty tree, no remote).
  Behavior = birth half (inverse of `state init`), reusing `ClassifyState`/`StateBranch`/`DiscoverWorkflowDir`; failure modes for occupied path, local orphan, remote orphan (redirect to `state init`), no-remote (best-effort push), and inline README — each empirically characterized and turned into AC-3/AC-4/AC-5 with verification.

### Summary

Spiked the birth mechanism first (checklist item 1) before designing: the orphan-branch + linked-worktree sequence works, tolerates a dirty main tree, and the in-repo `state_init_test.go` already proves it under test. Key design decision: `state new` is the mechanical inverse of `state init` and does NOT author the README (that's commission's interactive job) — it births everything around an already-present split-root README, reusing the same status helpers. Failure modes (occupied path, existing local/remote orphan, no-remote, inline README) were each empirically reproduced and turned into testable acceptance criteria with a real-git e2e test plan reusing the existing no-mock harness. Two open captain decisions flagged: whether to auto-commit the `.gitignore` edit (defaulted to leave-for-operator) and the `state new` vs. `commission --split-root` surface (designed `state new` as the primitive commission can call).

## Stage Report: implementation

- DONE: Implement `spacedock state new [--workflow-dir DIR]` as the mechanical inverse of `state init`: birth the orphan state branch + linked worktree at the README's `state:` path and append the .gitignore entry, reusing ClassifyState / StateBranch / DiscoverWorkflowDir / the commissionSplitWorkflow helper.
  `runStateNew` + `parseStateNewArgs` in `internal/cli/state.go` (mirror `runStateInit`); wired `state new` into `newStateCommand` (`internal/cli/cli.go`). Birth via `birthOrphanState` (temp detached worktree → `checkout --orphan` → clear inherited tree → seed → best-effort push → `worktree add`), exactly the `commissionSplitWorkflow` mechanic. Reuses `ClassifyState`/`StateBranch`/`DiscoverWorkflowDir`. Commit dfd2daa8.
- DONE: Cover the failure modes (AC-3/AC-4/AC-5): refuse occupied path / existing local-or-remote orphan (redirect to `state init`) / no-remote best-effort / inline README — clear messages, no raw git `fatal:` leaks to stderr.
  Three prechecks (`dirExists`, `rev-parse --verify` local ref, `orphanOnRemote` via `ls-remote`) return non-zero with tailored prose; remote case names `state init`. Inline README errors (AC-5); no-remote births locally + warns push-later (AC-4). `assertNoGitFatal` proves no `fatal:`/`already exists` leak; refusal messages say "already present" to avoid the git substring.
- DONE: Land the real-git e2e tests (no mocks, the state_init_test.go pattern): AC-1 birth, AC-2 round-trip with `state init`, AC-3 refusals.
  `internal/cli/state_new_test.go`: 7 e2e tests (AC-1..AC-5, incl. 3 AC-3 sub-tests) + 1 `lsRemoteHasBranch` table test (6 cases) reusing `git`/`gitOK`/`splitWorkflowReadme`/`runStatusBoot`. 15/15 pass; full repo 848/848; gofmt + vet clean.

### Summary

Implemented `spacedock state new` as the birth half of the split-root lifecycle, reusing the same README interpreters and the spiked `commissionSplitWorkflow` mechanic. Both captain-flagged defaults stand (leave the .gitignore edit for the operator; `state new` as `state init`'s inverse, not `commission --split`) — surface them at the validation gate. A live smoke surfaced a real bug the e2e `strings.Contains` assertion missed: the `.gitignore` entry was built with `filepath.Rel` against git's symlink-resolved `--show-toplevel` (`/var` → `/private/var` on macOS) vs. the unresolved `workflowDir`, yielding a garbage `../../…` path; fixed by deriving the entry from `git rev-parse --show-prefix`, and hardened the test to an exact-line check plus a `git check-ignore` behavioral assert.

## Stage Report: validation

**Recommendation: PASSED.**

- DONE: Reproduce evidence for every AC (AC-1 birth, AC-2 round-trip with state init, AC-3 refusals, AC-4 no-remote, AC-5 inline-README): run state_new_test.go (15 tests) + go test ./...; confirm the symlink-resolved .gitignore-path fix holds via the exact-line + git check-ignore asserts.
  `go test ./internal/cli -run 'TestStateNew|TestLsRemoteHasBranch' -v` → 15/15 named tests pass, zero skips (AC-1 `TestStateNewBirthsSplitRoot`, AC-2 `TestStateNewRoundTripsWithStateInit`, AC-3 `TestStateNewRefusesAlreadyBirthed` ×3 subtests, AC-4 `TestStateNewNoRemoteBestEffort`, AC-5 `TestStateNewInlineErrors`, unit `TestLsRemoteHasBranch` ×6). Full `go test ./...` 848/848, no FAIL/SKIP; `go build`/`go vet`/`gofmt -l` all clean.
- DONE: Confirm the symlink-resolved .gitignore-path fix holds via the exact-line + git check-ignore asserts.
  Live smoke in a real `/tmp`→`/private/tmp` symlinked repo (no remote): `state new` exits 0, births orphan `spacedock-state/dev` + linked worktree, writes the exact line `docs/dev/.spacedock-state/`, `git check-ignore` returns exit 0, `status --boot` shows `STATE_BACKEND: split-root` + `present: true`, code-branch porcelain shows only `?? .gitignore` (zero-churn holds), re-run refuses non-zero with no git fatal leak.
- DONE: Confirm the two captain-flagged defaults are correctly reflected in the deliverable and surface them for the gate decision.
  Verified in code + live: (1) `.gitignore` edit left for the operator — `appendGitignoreEntry` writes via `os.WriteFile`, never stages/commits; `runStateNew` prints "commit the .gitignore edit on your code branch"; the only commit in `state.go` is the orphan's `seed state`. (2) Surface is `spacedock state new` registered alongside `state init` (cli.go:327/329), NOT a `commission --split` flag — no `--split` exists; commission is only referenced in a doc-comment noting the README is its job. Both surfaced below for the gate.

### Detached adversarial audit (high-stakes: on-disk git mutation surface)

Ran on a separate detached checkout of dfd2daa8 (never the implementation worktree). 7 claim-breaking edits, all caught by the deliverable's own tests — no material findings:
1. Revert the symlink fix to `filepath.Rel` against `--show-toplevel` → AC-1 fails at the exact-line check with the garbage `../../../../../../tmp/.../.spacedock-state/` path (confirms the report's claim that the old `strings.Contains` would have green-lit it).
2. Leak `fatal:`/`already exists` into the occupied-path refusal → `assertNoGitFatal` (AC-3) fails.
3. Suppress the no-remote push-later warning → AC-4 fails.
4. Skip the orphan `git push` → AC-2 round-trip fails (fresh-clone `state init` can't fetch the orphan) — proves AC-2 is a true inverse check, not a tautology.
5. Drop the `state init` redirect from the local-orphan refusal → AC-3 local-orphan subtest fails.
6. Drop "split-root" from the inline error → AC-5 fails.
7. Append a `!`-negation after the correct ignore line (exact-line check still passes) → the `git check-ignore` behavioral assert fails — confirms the two AC-1 guards are independent and the functional assert proves behavior, not a spelling match.

### Captain decisions pending at the gate

- **Default A (leave .gitignore edit for the operator):** `state new` writes the `.gitignore` line but does NOT commit it, printing the next step. Alternative was auto-commit on the operator's code branch (rejected as an outward-facing side effect the operator should own).
- **Default B (`state new` vs `commission --split-root`):** shipped as `state new`, the mechanical inverse of `state init` reading the same README; the README stays commission's interactive job. `state new` is the thin primitive commission can later call.

### Polish (non-blocking)

- The shell-completion scripts (`internal/cli/cli.go` bash :518, zsh :538) complete the `state` subcommand with only `init --workflow-dir` — `new` is missing. `verbs_test.go` only asserts `new` appears as a top-level verb, so this gap is untested. Cosmetic (completion is "intentionally static, YAGNI" per the existing comment); does not affect any AC.

### Summary

All five acceptance criteria are verified with reproduced, out-of-body evidence (real-git e2e tests + a live smoke on a real symlinked macOS repo). 15/15 targeted tests and 848/848 full-suite pass; build/vet/gofmt clean. A detached adversarial audit refuted nothing material — all seven claim-breaking edits were caught, including confirmation that the hardened exact-line + `git check-ignore` asserts catch the symlink-path regression the original `strings.Contains` would have missed. Both captain-flagged defaults are correctly reflected and surfaced for the gate. One non-blocking polish item: `state new` is absent from the shell-completion scripts.
