---
id: ea9kke1e8q0wyhx0wjv4cyzr
title: CI PR gating on main (post-flip trunk)
status: done
source: captain (2026-06-13, during the 8p push turn) — pre-cut audit main-gating item pulled forward
started: 2026-06-13T07:01:21Z
completed: 2026-06-13T08:00:24Z
verdict: PASSED
score:
worktree:
issue:
sprint: 0201-post-flip-release-model
group: release-model
sprint-readiness: ready
mod-block:
pr: "#348"
archived: 2026-06-13T08:00:24Z
---

Post-flip, `main` is the release trunk and PRs target it, but `install-e2e.yml` and `runtime-live-e2e.yml` trigger `pull_request` on `[next]` only — so a PR to `main` runs neither the offline `go test ./...` gate nor install-e2e (only `docs.yml`, which already targets `main`, runs). Captain directive (2026-06-13): ensure CI gates fire on `main`-PRs before the sprint relies on the trunk model. This is the pre-cut antipattern-audit's known main-PR-gating item, pulled forward. Ideation pins the exact trigger change and the live-lanes design call.

## Problem

A PR opened against `main` today is checked by only one workflow: `docs.yml`, whose `build` job runs `mkdocs build --strict`. The two workflows that gate code correctness — `install-e2e.yml` (install.sh + checksum gate on Linux and macOS) and `runtime-live-e2e.yml` (whose `offline` job is the secret-free `go build ./...` + `go test ./...` suite) — both trigger `pull_request: branches: [next]`. They never see a `main`-PR. As the sprint moves to the trunk-based release model where PRs target `main`, code can merge to the release trunk without the Go test suite or the install smoke ever running. The `next`-only triggers are a leftover from the pre-flip model and must extend to `main`.

The other three workflows are correctly out of scope: `docs.yml` already triggers on `pull_request: branches: [main]`; `next-publish.yml` is `workflow_dispatch`-only; `release.yml` triggers on `push: tags`. None needs a trigger change.

## Proposed approach

Extend the `pull_request` branch filter in the two next-only workflows to include `main`, changing nothing about what the jobs do — only which PR base branches trigger them.

The non-obvious part is `runtime-live-e2e.yml`. Its `pull_request` trigger is workflow-level: there is no way to trigger only the cheap `offline` job and not the live jobs from the `branches:` filter, because the filter governs the whole workflow. Adding `main` queues all four jobs (`offline`, `claude-live`, `codex-live`, `pi-live`) on a `main`-PR. The safety comes from a different mechanism already in the file: every live job declares an `environment:` (a required-reviewer approval gate, reviewer = clkao), and `needs: offline`. So on a `main`-PR:

- `offline` runs unconditionally (no `environment`, no secret) — this is the `go build ./...` + `go test ./...` gate we want.
- `claude-live` (matrix `CI-E2E` / `CI-E2E-OPUS`), `codex-live` (`CI-E2E-CODEX`), and `pi-live` (`CI-E2E-PI`) each queue but pause in `waiting` on their environment until a maintainer approves that deployment. They do NOT auto-burn API credits on every `main`-PR.

This is exactly the shape `docs.yml` already uses: the `build` job runs on every PR, the credit/Pages-spending `deploy` job is gated (there by `if: github.ref == 'refs/heads/main'`; here by the per-environment approval). So the design decision is: **adopt the FO recommendation — the offline `go test ./...` gate and install-e2e run on every `main`-PR; the live lanes stay environment-gated so a `main`-PR does not auto-spend credits.** The alternative — splitting the offline gate into its own workflow so only it triggers on `main`, leaving the live lanes `next`-only — adds a second workflow file and duplicate setup steps to buy a property the environment gate already gives us for free. Rejected on YAGNI grounds: the live jobs already cannot fire without per-environment approval, so triggering the workflow on `main` is harmless.

`install-e2e.yml` has no live/secret distinction — both matrix legs (`ubuntu-latest`, `macos-latest`) run `goreleaser --snapshot` + `install.sh` locally with no secret and no environment. Adding `main` simply runs the install smoke on `main`-PRs, which is the intent.

### Before / after

`install-e2e.yml` (`on:` block):

```yaml
# before
on:
  workflow_dispatch:
  pull_request:
    branches: [next]

# after
on:
  workflow_dispatch:
  pull_request:
    branches: [next, main]
```

`runtime-live-e2e.yml` (`pull_request` trigger; the surrounding `workflow_dispatch` inputs and the SECURITY MODEL header comment are unchanged):

```yaml
# before
  pull_request:
    branches: [next]

# after
  pull_request:
    branches: [next, main]
```

No comment changes are required: the `runtime-live-e2e.yml` SECURITY MODEL header already says "pull_request on `next`" — implementation should refresh that single phrase to "`next` and `main`" so the comment stays true to the trigger, but the security reasoning (secrets withheld from fork PRs, per-variant environment approval) is branch-agnostic and unchanged.

## No doc-diff

CI trigger config is not user-visible CLI behavior. The docs site describes `spacedock` commands, host integration, and startup behavior — not which branches a GitHub Actions `pull_request` filter matches. This change alters no command surface, no output, no banner. No doc-site file changes.

## Riskiest-mechanism determination

No spike needed: the `pull_request: branches:[main]` trigger mechanism is ALREADY PROVEN on this repo. `docs.yml` triggers on `main`-PRs today, observed live — PR #347 targets `main` and its docs `build` check ran and reported `success` (with `deploy` correctly `skipped`, since `deploy` is `if: github.ref == 'refs/heads/main'` and a PR head ref is not `main`). The design composes only already-proven behavior: the same `pull_request: branches:` filter docs.yml uses, the same per-environment approval gate the live jobs already enforce on `next`-PRs. Nothing in the plan rests on an unexercised mechanism.

## Acceptance criteria

- **AC-1 — A `main`-PR triggers the offline `go test ./...` gate.** After the change, opening (or pushing to) a PR whose base is `main` causes `runtime-live-e2e.yml`'s `offline` job to run and execute `go build ./...` + `go test ./...`.
  - Verified by: the `offline` check appearing and reporting a conclusion on a real `main`-PR's checks list (`gh pr checks <pr>` / the commit's check-runs), where the PR's base ref is `main`. Fails if no `offline` check is present on a `main`-PR. NOT a grep of the workflow file — the expected value (a check-run named `offline` with a conclusion) comes from GitHub's run of the workflow, not from the YAML text under change.

- **AC-2 — A `main`-PR triggers install-e2e.** After the change, a `main`-base PR causes `install-e2e.yml`'s `install` matrix (`ubuntu-latest`, `macos-latest`) to run.
  - Verified by: the `install (ubuntu-latest)` and `install (macos-latest)` check-runs appearing with conclusions on a real `main`-PR's checks list. Fails if the install matrix does not appear on a `main`-PR. Independent source: GitHub's workflow run, not the YAML.

- **AC-3 — The live lanes do NOT auto-run on a `main`-PR.** After the change, a `main`-base PR's `claude-live` / `codex-live` / `pi-live` jobs are blocked in `waiting` on their environment (no API credits spent) until a maintainer approves; they do not start automatically.
  - Verified by: on a real `main`-PR, the live jobs show status `waiting` (pending environment approval) rather than `in_progress`/`completed` without an approval — observable in the workflow run's job list / `gh run view`. Fails if a live job starts without environment approval. Independent source: the run's job states, not the YAML.

- **AC-4 — `docs.yml` already covers `main`-PRs (no change made to it).** `docs.yml` triggers on `pull_request: branches: [main]` and its `build` job gates `main`-PRs, and this task does not modify `docs.yml`.
  - Verified by: PR #347 (base `main`) shows the docs `build` check `success` — already observed. The task's diff touches only `install-e2e.yml` and `runtime-live-e2e.yml`; `git diff` against `origin/main` lists no change to `docs.yml`, `next-publish.yml`, or `release.yml`.

Note: a regex/grep asserting `branches: [next, main]` is present in the changed YAML is NOT an acceptance criterion — the expected string is the same text the implementer wrote, so the check is a tautology (a valid reformat like `branches: ["next", "main"]` would fail it; an inverted edit would still need the string). The behavioral truth — that a `main`-PR actually triggers the jobs — is proven only by observing the checks on a real `main`-PR (AC-1/AC-2/AC-3), which draw their expected values from GitHub's run, a source independent of the YAML under change.

## Test plan

The proof is observable on a real `main`-base PR's checks, not a fixture:

1. **Cost/complexity:** Trivial mechanical edit (two one-token YAML changes plus a one-phrase comment refresh). No Go code, no test code, no fixtures. The verification cost is one real PR to `main` carrying the two workflow edits — its own check run is the test.
2. **What verifies it:** The implementation PR (which targets `main` and includes the trigger edits) is itself the end-to-end exercise. On that PR, observe via `gh pr checks <pr>` / `gh run view`:
   - `offline` check present and runs `go test ./...` (AC-1).
   - `install (ubuntu-latest)` + `install (macos-latest)` present (AC-2).
   - `claude-live` / `codex-live` / `pi-live` in `waiting` on their environments, not auto-started (AC-3).
   - `git diff origin/main` touches only the two files (AC-4).
   This is genuinely end-to-end: a workflow's `pull_request: branches:` change takes effect from the PR head, so the very PR that adds `main` to the filter is gated by the post-change triggers. No separate throwaway PR is needed.
3. **Fixture / CLI / live:** None of the three. This is GitHub Actions trigger config; the only meaningful test is the live PR-check observation above. No `spacedock` binary behavior changes, so no Go unit test, golden fixture, or behavior fixture applies.

## Stage Report: ideation

- DONE: Pin the exact change: add `main` to the `pull_request: branches:` array in install-e2e.yml ([next] -> [next, main]) and runtime-live-e2e.yml ([next] -> [next, main]); confirm docs.yml already covers main; and DECIDE + justify whether the live e2e lanes should fire on main-PRs or stay environment-gated/next-only
  Before/after YAML pinned for both files in "Before / after"; design decision adopted in "Proposed approach" (offline + install-e2e gate main-PRs, live lanes stay environment-gated — split-workflow alternative rejected on YAGNI since the per-environment approval already prevents auto-credit-burn); docs.yml confirmed `pull_request: branches: [main]` on origin/main and out of scope.
- DONE: Write entity-level ACs whose proof is OUTSIDE the entity prose and can fail; record "no doc-diff: CI trigger config is not user-visible CLI behavior."
  AC-1..AC-4 each cite a GitHub-run-observable (check-run presence/conclusion, job `waiting` state, `git diff origin/main` file list) — expected values come from GitHub's workflow run, not the YAML under change; a note explicitly bans the tautological grep-the-YAML check. "No doc-diff" section records CI trigger config is not user-visible CLI behavior.
- DONE: Record the riskiest-mechanism determination
  "Riskiest-mechanism determination" records "no spike needed" — the `pull_request: branches:[main]` mechanism is proven live: PR #347 (base `main`) shows docs `build` = success, `deploy` = skipped (verified via `gh api .../check-runs` this session).

### Summary

Fleshed out the ideation body: pinned the two one-token YAML trigger edits ([next] -> [next, main] in install-e2e.yml and runtime-live-e2e.yml), with before/after blocks. Verified all FO findings against origin/main YAML and confirmed the main-PR trigger is already proven live (PR #347's docs build ran and passed, deploy skipped). Adopted the FO's design call: the offline `go test ./...` gate and install-e2e run on main-PRs while the three live lanes stay environment-gated (each `needs: offline` + an approval-gated `environment:`), so a main-PR never auto-burns API credits — same shape as docs.yml's build-runs/deploy-gated split. ACs are proven by observing checks on a real main-PR, never by grepping the workflow file.

## Stage Report: implementation

- DONE: Add `main` to the `pull_request: branches:` array in BOTH install-e2e.yml ([next] -> [next, main]) and runtime-live-e2e.yml ([next] -> [next, main]), and refresh the runtime-live-e2e SECURITY-MODEL header comment phrase from "pull_request on next" to "next and main" (the security reasoning is branch-agnostic, unchanged). `git diff origin/main` must list ONLY those two files — docs.yml, next-publish.yml, release.yml untouched (AC-4).
  Commit 79465a41 on spacedock-ensign/ci-pr-gating-on-main: install-e2e.yml +1/-1, runtime-live-e2e.yml +2/-2 (trigger filter + header phrase "`next`" -> "`next` and `main`"); `git diff --name-only origin/main` lists exactly those two files — docs.yml / next-publish.yml / release.yml untouched (AC-4 satisfied for the diff portion).
- DONE: Run `go test ./...` green — this is a YAML-only edit so Go behavior is unchanged, but confirm no workflow-structure/parse guard (e.g. internal/release workflow_exec_guard or any workflow-validating test) reds on the trigger change. Report the pass count.
  `go test ./...` = 1249 passed in 16 packages; `internal/release` (incl. workflow_exec_guard, which asserts runtime-live-e2e.yml job steps/secrets/scenario commands — none touched by a trigger/comment edit) = 40 passed. No guard reds on the trigger change.
- DONE: Record in the stage report the AC-1/AC-2/AC-3 proof plan to be observed on THIS task's own PR to main (offline `go test` gate + install matrix RUN; the 3 live lanes sit `waiting` on environment approval, NOT auto-started). Add NO tautological YAML-grep test, NO doc-diff (CI trigger config is not user-visible), NO Go/fixture test — per the ideation test plan the PR's own check run is the proof.
  Proof plan recorded below; no YAML-grep test, no doc-diff, no Go/fixture test added (the PR's own check run is the end-to-end proof — a `pull_request: branches:` change takes effect from the PR head).

### AC-1/AC-2/AC-3 proof plan (observed on THIS task's own PR to main)

Open a PR with base `main` carrying commit 79465a41. Because a `pull_request: branches:` change takes effect from the PR head, this PR is itself gated by the post-change triggers — no throwaway PR needed. On that PR observe:
- AC-1: `offline` check (runtime-live-e2e.yml) PRESENT and runs `go build ./...` + `go test ./...`, reporting a conclusion. Source: GitHub's run, via `gh pr checks <pr>` / `gh run view`.
- AC-2: `install (ubuntu-latest)` + `install (macos-latest)` checks (install-e2e.yml) PRESENT with conclusions.
- AC-3: `claude-live` (sonnet + claude-opus-4-8), `codex-live`, `pi-live` jobs sit in `waiting` on their environments (CI-E2E / CI-E2E-OPUS / CI-E2E-CODEX / CI-E2E-PI), NOT auto-started — no API credits spent until a maintainer approves. Observe job states in `gh run view`.
- AC-4 (diff portion): already satisfied — `git diff --name-only origin/main` = only the two workflow files.

### Branch-protection heads-up for the captain (operational, not part of this edit)

On `main`'s branch protection, only `offline` / `install (ubuntu-latest)` / `install (macos-latest)` / docs `build` should be REQUIRED status checks. The three live lanes (`claude-live`, `codex-live`, `pi-live`) must NOT be required: they pause in `waiting` on per-environment approval by design, so requiring them would block every `main`-PR on manual environment approval (and burn API credits on every merge). This is a repo-settings action outside the worktree diff.

### Summary

Mechanical trigger edit landed in commit 79465a41 on spacedock-ensign/ci-pr-gating-on-main: added `main` to the `pull_request: branches:` filter of install-e2e.yml and runtime-live-e2e.yml ([next] -> [next, main]) and refreshed the runtime-live-e2e SECURITY-MODEL header phrase to "`next` and `main`" so the comment stays true (security reasoning unchanged, branch-agnostic). No job behavior changed. `git diff origin/main` touches only those two files (AC-4 diff). `go test ./...` is green at 1249/1249 and the internal/release workflow_exec_guard (40/40) does not red on the trigger change. AC-1/AC-2/AC-3 are proven on this task's own PR to main (offline `go test` + install matrix RUN; the three live lanes sit `waiting` on environment approval, not auto-started) — the PR's own check run is the end-to-end test, so no YAML-grep test, no doc-diff, and no Go/fixture test were added. Flagged for the captain: on main's branch protection only offline/install/docs build should be REQUIRED checks, never the live lanes.

## Stage Report: validation

- DONE: Confirm the trigger diff is exactly correct from the worktree: `pull_request: branches: [next]` -> `[next, main]` in BOTH install-e2e.yml AND runtime-live-e2e.yml, the runtime-live-e2e SECURITY-MODEL header comment refreshed to "next and main" (security reasoning unchanged); `git diff --name-only origin/main` lists ONLY those two files.
  Worktree HEAD = 79465a41, base = origin/main (0ddc9ad6). `git diff origin/main`: install-e2e.yml +1/-1 (`branches: [next]`->`[next, main]`), runtime-live-e2e.yml +2/-2 (trigger `branches: [next]`->`[next, main]` + header comment `# pull_request on \`next\`.`->`# pull_request on \`next\` and \`main\`.`, secrets/permissions reasoning unchanged). `git diff --name-only origin/main` = exactly those two files — docs.yml / next-publish.yml / release.yml untouched.
- DONE: Reproduce `go test ./...` green AND `go test ./internal/release/` (workflow_exec_guard parses runtime-live-e2e.yml job steps/secrets/scenario commands); confirm the edit reds NO workflow-structure guard, and the edited YAML still parses as a valid workflow. Report pass counts.
  `go build ./...` OK; `go test ./...` = 1249 passed in 16 packages, EXIT=0, no FAIL/panic. `go test ./internal/release/` EXIT=0 (matches implementer's 40/40; -v shows 55 PASS subtests, 0 FAIL). The guard reads `../../.github/workflows/runtime-live-e2e.yml` from disk (journey_workflow_test.go:628), i.e. parses the edited file — green, so no workflow-structure guard reds on the trigger/comment edit. `yq -e '.on.pull_request.branches'` parses both files = `[next, main]`; parsed jobs = offline/claude-live/codex-live/pi-live.
- DONE: Confirm the AC-1/AC-2/AC-3 proof PLAN is sound and EXPLICITLY note the definitive proof is observed on THIS task's own PR-to-main checks at the merge boundary; the FO confirms the checks when the PR opens, before merge. Recommend PASSED/REJECTED.
  Plan is sound: a `pull_request: branches:` change takes effect from the PR head, so this task's own main-PR is gated by the post-change triggers (no throwaway PR). Structural backstops verified now in the worktree: `offline` has NO `environment` (runs unconditionally — AC-1/AC-2 gate); `claude-live`/`codex-live`/`pi-live` each declare `needs: offline` AND an `environment:` (CI-E2E/CI-E2E-OPUS via matrix, CI-E2E-CODEX, CI-E2E-PI) → AC-3 live lanes pause `waiting` on approval, not auto-started. Mechanism proven live: PR #347 (base `main`) docs `build` = pass, `deploy` = skipping (`gh pr checks 347`) — the `pull_request: branches:[main]` filter the design reuses already fires on main-PRs. RECOMMEND PASSED.

### Detached adversarial audit

CI machinery → high-stakes surface. On a detached throwaway checkout of the merge result (79465a41, never the implementation worktree; removed after), stripping the `CI-E2E-PI` environment from `pi-live` flips `internal/release` from EXIT=0 to FAIL (`TestWorkflowsPreserveAndPublishJourneyCosts`: "runtime-live-e2e.yml Pi live job is missing its ... CI-E2E-PI environment"). So the environment gates backing AC-3 are protected by a real, non-vacuous guard. No `internal/release` guard asserts on the `branches:`/`pull_request` trigger, so the guards correctly do not constrain (or get weakened by) this trigger edit. No Go/fixture/grep test was added (per dispatch), so there is no new in-repo assertion to refute — the behavioral proof is GitHub's live run on the PR. No material findings.

### Summary

Validated commit 79465a41 from the implementation worktree (base origin/main 0ddc9ad6). The diff is exactly the pinned change: `[next]`->`[next, main]` in both install-e2e.yml and runtime-live-e2e.yml plus the one-phrase SECURITY-MODEL header refresh; `git diff --name-only origin/main` lists only those two files. `go test ./...` = 1249/1249 green and `internal/release` (the workflow_exec_guard that parses the on-disk runtime-live-e2e.yml) = green — no workflow-structure guard reds on the edit, and both files parse via `yq`. AC-4's diff portion is satisfied now and its docs-cover-main half is proven live by PR #347 (docs build pass on a main-PR). AC-1/AC-2/AC-3 are structurally backstopped in the worktree (offline unconditional; live lanes `needs: offline` + environment-gated) and their definitive check-run observation happens on THIS task's own PR-to-main at the merge boundary — the validator cannot open the PR, so the FO must confirm the `offline`/`install` checks ran and the live lanes sit `waiting` (not auto-started) when the PR opens, before merge. A detached adversarial audit refuted the vacuous-guard hypothesis (removing a live-job environment flips the guard red) and found no material issue. RECOMMEND PASSED, with the merge-boundary PR-check confirmation as the FO's remaining gate item.
