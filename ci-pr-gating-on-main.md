---
id: ea9kke1e8q0wyhx0wjv4cyzr
title: CI PR gating on main (post-flip trunk)
status: implementation
source: captain (2026-06-13, during the 8p push turn) — pre-cut audit main-gating item pulled forward
started: 2026-06-13T07:01:21Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-ci-pr-gating-on-main
issue:
sprint: 0201-post-flip-release-model
group: release-model
sprint-readiness: ready
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
