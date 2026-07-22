---
title: Ensure the automatic pre0 tag push triggers its release workflow
status: validation
source: "v0.25.0 stable cut, 2026-07-15: edge-advance pushed annotated v0.26.0-pre0, but no release run appeared; manually replaying the identical tag triggered release.yml immediately."
started: 2026-07-21T15:58:31Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-auto-pre0-tag-push-release-trigger
issue:
milestone: 0.27.0
id: 5aqczjeq6rq3mckbc5gyjqe3
mod-block: merge:pr-merge
pr: pr-merge:553
---

A stable release must publish its automatically generated next-minor edge binary without requiring an operator to delete and replay the pre0 tag.

## Problem

The stable `edge-advance` job created and pushed the correct annotated `vX.(Y+1).0-pre0` tag, and the remote tag pointed at the intended green release commit. No `release.yml` run was created for that push. Deleting and manually re-pushing the identical annotated tag immediately created the expected run, which then published the edge release and cask successfully.

The current job can therefore report success while leaving the edge binary behind the newly advanced `next` skills.

## Reproduction

1. Push a stable `vX.Y.Z` tag and let the release workflow reach `edge-advance`.
2. Observe the job create and push annotated tag `vX.(Y+1).0-pre0` using its configured release token.
3. Confirm the remote tag exists on the stable release commit.
4. Observe that no release workflow run exists for the pre0 tag.
5. Delete and manually re-push the same local annotated tag; observe that the release workflow starts immediately.

## Root cause

Not a race. The `edge-advance` job's "Auto-cut the edge prerelease tag" step
(`.github/workflows/release.yml`, currently the job's last step) constructs the
annotated `vX.(Y+1).0-pre0` tag correctly and pushes it with:

    git push "https://x-access-token:${HOMEBREW_TAP_TOKEN}@github.com/${GITHUB_REPOSITORY}.git" "$PRE0_TAG"

`HOMEBREW_TAP_TOKEN` is a **cross-repo PAT provisioned for
`spacedock-dev/homebrew-tap`** (goreleaser reuses it to bump the tap casks). The
step reuses it for a **same-repo** tag push on the belief — stated in the step's
own comment ("PAT … so the tag push RE-TRIGGERS release.yml") — that any
non-`GITHUB_TOKEN` PAT triggers workflows. That belief is false for this token:
the push lands the remote ref but creates **no** `release.yml` run. It is a
workflow-run-suppressed credential — it behaves like `GITHUB_TOKEN` for
triggering (GitHub does not create a new workflow run for a `push` event
authenticated by a suppressed credential).

Proof it is the credential/event path and not a race: the identical tag object,
deleted and re-pushed under an **operator PAT/ssh** credential, fires `release.yml`
within seconds — twice (both cuts below). Same ref, same peeled commit, same
body; only the pushing credential differs. GitHub's documented rule: pushes made
with `GITHUB_TOKEN` do not create workflow runs, while pushes made with a PAT,
a GitHub App installation token, or an **SSH deploy key** do; `workflow_dispatch`
and `repository_dispatch` always create runs. So the fix is to push the pre0 tag
with a credential whose trigger-capability is unconditional, and to make a
non-triggering credential fail loudly instead of leaving a silent gap.

## Second reproduction (2026-07-21, v0.26.0 stable cut)

Confirms the first occurrence (v0.25.0 cut, 2026-07-15) was not a fluke.

- Release run `29845035279` (tag `v0.26.0`) reached `edge-advance`,
  decision=`advance`. Its auto-cut step ran verbatim `git tag -a "$PRE0_TAG"
  "$RELEASE_COMMIT" -m "$PRE0_BODY"` then the `x-access-token:${HOMEBREW_TAP_TOKEN}`
  push.
- Push SUCCEEDED (remote: `* [new tag] v0.27.0-pre0 -> v0.27.0-pre0`); tag object
  `ac240b4abf8652868dc6b3d9d10a7bb931bf353d`, annotated, peeled to release commit
  `ca136f83`.
- **No** `release.yml` run was created for the pre0 tag. Waited >5 min; nothing
  fired.
- Operator remedy (`git push origin :refs/tags/v0.27.0-pre0` then
  `git push origin v0.27.0-pre0`, under an operator PAT/ssh) fired release run
  `29845875763` within seconds: it published the `v0.27.0-pre0` edge release (4
  edge tarballs) and bumped the `spacedock@next` cask to `v0.27.0-pre0` (tap
  commit `2c601b1c`).
- Run-list confirmation: the ONLY `release.yml` run for each auto-cut pre0 tag is
  the operator replay — `v0.26.0-pre0` → run `29423262944` (2026-07-15), and
  `v0.27.0-pre0` → run `29845875763` (2026-07-21) — each firing minutes after its
  stable run, via a `push` event. The auto-cut push produced zero runs.

Impact: `next`'s skills advance to the new pre-version while the newest edge
binary lags, so `spacedock@next` users hit a version-skew abort until an operator
manually replays the tag. Shipped unfixed in 0.26.0 despite the 0.26.0 milestone
label — recommend re-milestoning to 0.27.0 (frontmatter left unchanged here;
flagged to the FO).

## Proposed approach

Keep the tag construction exactly as it is (annotated, non-empty body, on the
greened `$RELEASE_COMMIT`) — AC-3 is preserved because nothing about the tag
OBJECT changes. Change only two things in the auto-cut step:

1. **Trigger-capable push credential.** Push the pre0 tag over SSH using a
   dedicated **write deploy key** (`EDGE_RELEASE_DEPLOY_KEY` secret) scoped to
   THIS repo — `git@github.com:${GITHUB_REPOSITORY}.git`. The pre0 push is a
   same-repo push, so a repo-scoped deploy key is the minimal-privilege fit (no
   cross-repo reach, unlike the tap PAT it replaces). Deploy-key pushes are the
   one credential GitHub never restricts for triggering, they do not expire, and
   the operator replay already proved an ssh-transport push of this exact tag
   fires the run — so this productionizes the proven positive control.

2. **Verify-or-fail guard.** After the push, poll the Actions API (read-only, via
   the default `GITHUB_TOKEN`) for a `release.yml` run whose `head_branch` equals
   the pre0 tag, for a bounded wait (~2 min, ~10s poll). If none appears, `exit 1`.
   `edge-advance` is a sibling of `goreleaser` (`needs: goreleaser`), so this red
   does NOT unwind the already-published stable release — it surfaces the edge
   handoff as failed instead of silently green. This guard IS the mechanism AC-2
   requires: a suppressed credential now reds the job before the run reports the
   edge handoff complete, rather than leaving the gap the operator discovers later.

**Value AC each mechanism serves.**
- Deploy-key push → **AC-1** (the run fires automatically; edge binary + cask
  publish without an operator replay). Simplest alternative considered: swap the
  https URL to a new trigger-capable **classic PAT** secret (one-line change).
  Insufficient as the primary choice: it repeats the exact class of failure that
  produced this bug (a credential we *believe* triggers, whose trigger-capability
  is conditional on flavor/scope and silently lapses on expiry) — the root cause
  we could not fully explain for the current PAT. A deploy key removes that
  conditionality. The classic-PAT swap is recorded as the smaller-surface fallback
  if deploy-key SSH setup on the runner proves troublesome.
- Verify-or-fail guard → **AC-2** (loud failure on a suppressed/incapable
  credential). Simplest alternative considered: no guard, trust the deploy key.
  Insufficient: AC-2 explicitly requires a suppressed token to FAIL the handoff,
  and without the guard a future credential regression (rotation, revocation,
  someone reverting to the tap PAT) reintroduces the silent gap — this is the
  bug's third-occurrence prevention.

Rejected — `repository_dispatch`/`workflow_dispatch` re-trigger (which GitHub
lets even `GITHUB_TOKEN` fire): it would require every `release.yml` job
(`e2e-gate`, `goreleaser`, `journey-ledger`, `edge-advance`) to resolve the tag
from a dispatch payload instead of `github.ref`, rewriting the core of the
working stable pipeline. That is the LARGEST correction, not the smallest, and
risks the stable path the Test plan says to protect. Rejected — a polling
controller (a separate scheduled workflow that watches for tagless pre0s and
re-fires): the Test plan explicitly disfavors it, and a direct trigger-capable
push satisfies AC-1 without it.

## Acceptance criteria

**AC-1 (VALUE):** A stable `vX.Y.Z` cut automatically publishes the next-minor
`vX.(Y+1).0-pre0` edge release and bumps the `spacedock@next` cask to that
version **without an operator deleting and replaying the tag** — a distinct
`release.yml` run is created on the pre0 commit by the auto-cut push itself.

Measured against baseline (a number that can move the wrong way): the count of
`release.yml` runs created for the auto-pushed pre0 tag goes from **0** (today,
suppressed) to **≥1** (auto-push), and the `spacedock@next` cask version advances
to the pre0 version with **no** operator-replay run in the tag's run history.
Verified by: the in-CI trigger probe (below) proving a runner-embedded deploy-key
tag push fires a run, plus the first live stable cut after the fix — self-proving
because the verify-or-fail guard reds the cut if the run does not appear.

**AC-2:** The pre0 push uses a credential whose trigger-capability is
unconditional (SSH deploy key), and the auto-cut step FAILS (reds `edge-advance`)
before reporting the edge handoff complete when the pre0 push produces no
`release.yml` run — so a workflow-suppressed token shape (the default
`GITHUB_TOKEN`, or the current cross-repo tap PAT) cannot pass silently.

Verified by: a negative control that points the pre0 push at `GITHUB_TOKEN`
(suppressed) and observes the verify-or-fail guard exit non-zero / red the job;
a positive control being the trigger-capable deploy-key push firing the run. The
2026-07-21 operator replay is a real-world positive-control datapoint (a
trigger-capable push of the identical tag fired run `29845875763`).

**AC-3:** The auto-cut pre0 tag remains annotated, non-empty in body, and peeled
to the stable release commit; the fix changes only the push transport/credential
and adds the verify guard — it never alters the tag object's construction, source
tree, or body, and invents no standalone changes.

Verified by: the tag-object construction (`git tag -a "$PRE0_TAG"
"$RELEASE_COMMIT" -m "$PRE0_BODY"`) is unchanged in the diff; a fixture/assertion
that the pushed pre0 tag is annotated, has a non-empty body, and peels to
`$RELEASE_COMMIT` (as `ac240b4a`→`ca136f83` did) after the automated handoff.

## Test plan

1. **Mechanism spike (riskiest path, minutes, do first).** A throwaway
   `workflow_dispatch` scratch job pushes a `zz-trigger-probe-<sha>` throwaway
   tag to this repo over the deploy-key SSH transport, then polls the Actions API
   for a run on that tag; assert a run appears (positive control), then delete the
   probe tag. Re-run the same probe pushing via `${{ secrets.GITHUB_TOKEN }}` and
   assert NO run appears within the window and the verify poll would exit non-zero
   (negative control). This exercises the one thing the manual replay did not: a
   **runner-embedded** trigger-capable push firing a run. Not committed (scratch).
   If deploy-key SSH on the runner is problematic, the same probe validates the
   classic-PAT fallback before committing to it.
   - Cost: low (one probe run per credential); no full cut needed.
2. **Regression guard (Go, release-guard family).** Extend the
   `internal/release/*_workflow_test.go` guards to assert the auto-cut pre0 step
   (a) does NOT authenticate its push with the default `GITHUB_TOKEN`, and (b)
   contains a verify-or-fail poll keyed on the pre0 tag that exits non-zero when
   no run is found. This is a structural regression fence (text over `release.yml`,
   the workflow IS the artifact) in the idiom of `goreleaser_guard_test.go` /
   `e2egate_workflow_test.go` — it prevents a silent revert to a suppressed
   credential. It is NOT the behavioral proof of AC-1; the spike and the live cut
   are.
   - Cost: low; `go test ./internal/release/...`.
3. **Live confirmation (self-proving).** The first real stable cut after merge
   either fires the pre0 run automatically (AC-1 proven live: edge tarballs + cask
   bump appear with no operator-replay run in the tag history) or reds
   `edge-advance` at the verify guard (AC-2 proven live). Retain the tag object,
   run id, release assets, and cask commit as evidence — the same evidence set the
   2026-07-21 reproduction captured.

## Expected surface

Files and LOC this task expects to touch (the ideation-gate baseline; later
rounds calibrate against it):

- `.github/workflows/release.yml` — the "Auto-cut the edge prerelease tag" step
  only: deploy-key SSH setup, push transport swap, verify-or-fail poll, and a
  rewrite of the now-false "PAT re-triggers" comment. **~ +25 to +40 LOC.**
- `docs/releasing.md` — rewrite the "Advancing the Edge Line" claim that the pre0
  tag is "pushed with the PAT (a `GITHUB_TOKEN` push does not fire the pre0
  build)" to describe the deploy-key push + verify-or-fail guard. **~ +8 to +15
  LOC (diff below).**
- `internal/release/*_workflow_test.go` — new/extended guard asserting the pre0
  step's credential and the verify poll. **~ +40 to +70 LOC.**
- Repo secret `EDGE_RELEASE_DEPLOY_KEY` (an org/repo write deploy key) — a
  provisioning step, not a committed file; documented in the releasing doc.
- Scratch probe workflow for the spike — throwaway, not committed.

**Total committed: 3 files, ~ +75 to +125 LOC.** Tolerance: if the change spills
past ~200 LOC, touches jobs other than `edge-advance`, or requires reworking any
job's `github.ref` resolution (i.e. drifts toward the rejected dispatch rewrite),
re-gate before proceeding.

## Documentation diff (docs/releasing.md)

The "Advancing the Edge Line (`next`)" section, stable-tag bullet, currently ends
(lines ~180-188):

> — then an ANNOTATED `vX.(Y+1).0-pre0` tag is auto-created **on the greened
> release commit** and pushed (via the re-triggering tap PAT). … The auto-tag
> MUST be annotated with a non-empty body (the release-notes extraction step
> rejects a lightweight tag), and MUST be pushed with the PAT (a `GITHUB_TOKEN`
> push does not fire the pre0 build). Expect two GitHub releases per stable cut.

Rewrite the credential clause to:

> — then an ANNOTATED `vX.(Y+1).0-pre0` tag is auto-created **on the greened
> release commit** and pushed over SSH with a dedicated write **deploy key**
> (`EDGE_RELEASE_DEPLOY_KEY`), scoped to this repo. … The auto-tag MUST be
> annotated with a non-empty body (the release-notes extraction step rejects a
> lightweight tag). The push MUST use a trigger-capable credential: a
> `GITHUB_TOKEN` push — and, as observed on the v0.25.0 and v0.26.0 cuts, the
> cross-repo tap PAT — does NOT create the pre0 `release.yml` run, so the step
> pushes with the deploy key and then **verifies a run was created for the pre0
> tag, failing `edge-advance` loudly if none appears** rather than leaving the
> edge binary silently behind. Expect two GitHub releases per stable cut.

## Spike record

- **Already proven (positive control, no new spike needed for this leg):** the
  2026-07-21 operator replay of the identical `v0.27.0-pre0` tag object under an
  operator PAT/ssh fired release run `29845875763` within seconds (and
  `v0.26.0-pre0` → `29423262944` on 2026-07-15). A trigger-capable push of the
  pre0 tag DOES fire the run and complete the edge publish + cask bump. GitHub's
  documented triggering rules (GITHUB_TOKEN suppressed; PAT / app token / deploy
  key trigger; dispatch events always trigger) corroborate this.
- **Still to spike before implementation (Test-plan step 1):** that a
  **runner-embedded** deploy-key push (the CI-embeddable form, not an operator's
  laptop credential) fires the run — the one hop the manual replay did not cover.
  Cheap: a `workflow_dispatch` probe pushing a throwaway tag, minutes-cost, run
  first per the mechanism-validation-before-comprehensive discipline.

## Stage Report: ideation

- DONE: Diagnose why the edge-advance pre0 tag push lands the remote tag but creates no release.yml run; name the root cause (token/event path, not a race).
  Root cause section: the same-repo pre0 push reuses the cross-repo `HOMEBREW_TAP_TOKEN`, a workflow-run-suppressed credential (behaves like GITHUB_TOKEN for triggering); ref lands, no run fires. Proven not-a-race by the identical-object operator replay firing a run twice (`release.yml:403-419`).
- DONE: Design the smallest correction to the credential/event path; prefer a trigger-capable PAT or explicit re-trigger over a polling controller.
  Proposed approach: push the pre0 tag over SSH with a repo-scoped write deploy key + a verify-or-fail guard. Dispatch-rewrite and polling-controller both rejected with reasons; classic-PAT swap kept as smaller-surface fallback.
- DONE: Record the 2026-07-21 second reproduction (v0.26.0 cut) in the task body.
  "Second reproduction" section: `x-access-token:${HOMEBREW_TAP_TOKEN}` pushed `v0.27.0-pre0` (tag `ac240b4a`→`ca136f83`), no run; operator replay fired run `29845875763`, published edge release + bumped `spacedock@next` cask (tap `2c601b1c`). Corroborated by `gh run list`.
- DONE: Finalize AC-1/AC-2/AC-3 and the Test plan against the chosen fix.
  AC-1 measures runs-on-pre0-tag 0→≥1 + cask advance with no replay run; AC-2 names the GITHUB_TOKEN/tap-PAT negative control and the verify-or-fail guard; AC-3 pins tag-object construction unchanged. Test plan layers spike → guard → live self-proving cut.
- DONE: Declare the expected surface (files + LOC) and tolerance.
  Expected surface section: 3 committed files (`release.yml`, `docs/releasing.md`, a `internal/release/*_workflow_test.go` guard) ~ +75 to +125 LOC; deploy-key secret + scratch probe uncommitted; re-gate past ~200 LOC or if it drifts into a github.ref rewrite.
- DONE: Name the value AC each mechanism serves; either spike the riskiest path or record that the manual replay already proves a trigger-capable push fires the run.
  Deploy-key→AC-1, verify-guard→AC-2, each with simplest-alternative-and-why. Spike record: manual replay already proves a trigger-capable push fires the run; the one remaining hop (runner-embedded deploy-key push) is the Test-plan step-1 probe to run first.
- DONE: If the fix changes the documented release process, propose the concrete docs/releasing.md before/after diff.
  "Documentation diff" section: before/after wording for the "Advancing the Edge Line" stable-tag bullet's credential clause (tap PAT → deploy key + verify-or-fail).

### Summary

Diagnosed the pre0-tag-no-run failure as a credential/event-path bug, not a race: the auto-cut step pushes a same-repo tag with a cross-repo tap PAT that is workflow-run-suppressed (empirically behaves like GITHUB_TOKEN), so the ref lands but no `release.yml` run fires; the identical object replayed under an operator PAT/ssh fires the run within seconds — confirmed across both the v0.25.0 and v0.26.0 cuts. The chosen fix is the smallest auth-path correction: push the pre0 tag over SSH with a repo-scoped write deploy key (unconditional trigger-capability, no expiry, minimal privilege) plus a verify-or-fail guard that reds `edge-advance` if no pre0 run appears — the guard is what makes AC-2's suppressed-token negative control fail loudly and prevents a third silent recurrence. One open item for the FO/gate: the milestone reads 0.26.0 but the bug shipped unfixed in 0.26.0 — recommend re-milestoning to 0.27.0 (frontmatter left unchanged per the no-frontmatter rule); and the deploy-key-vs-classic-PAT choice is a real decision the gate may want to weigh (I recommend the deploy key for robustness given the prior credential's unexplained failure).

## Gate: ideation — APPROVED (FO, captain-ruled)

- **Verdict:** ideation design approved for implementation. Design is sound, behavior-first, ACs measure value against a movable baseline, riskiest hop correctly held for spike-first.
- **Credential decision (captain ruling):** SSH write **deploy key** (`EDGE_RELEASE_DEPLOY_KEY`), repo-scoped — chosen over the classic-PAT swap for unconditional trigger-capability and to not repeat the prior credential's unexplained suppression.
- **Milestone:** corrected `0.26.0 → 0.27.0` (shipped unfixed in 0.26.0).
- **BLOCKER before implementation:** the `EDGE_RELEASE_DEPLOY_KEY` repo secret must be provisioned (generate keypair → add repo deploy key with write → store private key as the secret). Test-plan step-1 spike (runner-embedded deploy-key push fires a run) depends on it; implementation holds at ideation until the secret exists.
- **Next on unblock:** dispatch implementation, spike-first (the `workflow_dispatch` probe), then the `release.yml` step + verify-or-fail guard + `docs/releasing.md` diff + `internal/release/*_workflow_test.go` guard.

## Stage Report: implementation

- DONE: Spike FIRST (Test-plan step 1) — runner-embedded deploy-key push fires a run (positive), GITHUB_TOKEN push fires none + verify poll exits non-zero (negative); delete probe tag + scratch job after.
  Push-triggered scratch workflow on branch `zz-trigger-probe-scratch` pushed two tags at HEAD from the runner. Positive: deploy-key SSH push created observer run `29884082483` (driver `29884075373`, success). Negative: GITHUB_TOKEN push landed the ref (GT ref on origin=1) but created 0 runs across the full ~160s window. Both refs landed, so the sole variable is credential-driven run creation. Cleanup verified: both probe tags + scratch branch gone from origin, both runs 404, worktree removed. Go-ahead obtained from team-lead before the outward push.
- DONE: Edit release.yml "Auto-cut the edge prerelease tag" step: deploy-key SSH push (`git@github.com:${GITHUB_REPOSITORY}.git`) + verify-or-fail guard (poll release.yml runs on the pre0 tag, exit 1 if none within ~2 min) + rewrite the false "PAT re-triggers" comment. Apply the docs/releasing.md diff.
  Commit `ff10ed47`, `.github/workflows/release.yml` lines ~389-451: env swapped `HOMEBREW_TAP_TOKEN`→`EDGE_RELEASE_DEPLOY_KEY`+`GH_TOKEN`; SSH deploy-key setup; verify poll filters `head_branch==$PRE0_TAG` (read-only-validated against the real `v0.27.0-pre0` run — count returns 1). Tag object construction unchanged (AC-3). `docs/releasing.md` credential clause rewritten to deploy-key + verify-or-fail per the approved diff.
- DONE: Add the Go regression guard (asserts the pre0 step does NOT authenticate with the default GITHUB_TOKEN and DOES contain the verify-or-fail poll keyed on the pre0 tag). gofmt clean, go test green, go build.
  The pre-existing always-cut-pre0 guard (`internal/release/edge_advance_wiring_test.go`) hard-coded the old `HOMEBREW_TAP_TOKEN` push and broke; extended it to own the new invariants rather than add a duplicating file. `assertAlwaysCutPre0` now rejects a GITHUB_TOKEN pre0 push, requires the deploy-key SSH transport, and requires the verify poll (release.yml runs query keyed on the pre0 tag + `exit 1`). Twins: (d) swap deploy-key→GITHUB_TOKEN push must red; (e) neuter the run-verification query must red. `go test ./internal/release/...` green, gofmt/vet/build clean.

### Summary

Spike-first validated the one unproven hop: a runner-embedded deploy-key tag push creates a workflow run while a GITHUB_TOKEN push does not (both refs landed — credential is the sole variable), productionizing the operator-replay positive control. The fix pushes the pre0 tag over the `EDGE_RELEASE_DEPLOY_KEY` deploy key and adds a verify-or-fail poll that reds `edge-advance` if no pre0 `release.yml` run appears, closing AC-1 (auto-publish) and AC-2 (loud failure on a suppressed credential) while leaving the tag object untouched (AC-3). Surface: 3 files, net +68 LOC — one deviation from the declared surface: I extended the existing `edge_advance_wiring_test.go` guard (which broke on the credential change) instead of adding a new `*_workflow_test.go` file, keeping it to 3 files, edge-advance-scoped, no `github.ref` rewrite. AC-1/AC-2 live confirmation (Test-plan step 3) remains the first real stable cut after merge, self-proving via the verify guard.

## Stage Report: validation

- DONE: Do NOT re-run the outward spike; validate the deliverable + accept the documented spike record — positive control observer run 29884082483 (deploy-key push fired a run), negative control 0 runs for the GITHUB_TOKEN push over ~160s; confirm the methodology is sound and the credential was the sole variable (both refs landed).
  Accepted without re-running (`git ls-remote` shows no `zz-trigger-probe` refs left on origin — spike cleaned up). Methodology sound: two tags at the SAME HEAD, deploy-key push → run vs GITHUB_TOKEN push → 0 runs, both refs landed ⇒ credential is the isolated single variable. Clean controlled experiment; productionizes directly into the fix.
- DONE: Verify AC-1 (auto-fire over EDGE_RELEASE_DEPLOY_KEY SSH; spike positive control is behavioral proof, verify-or-fail poll is in-CI enforcement); AC-2 (poll exits 1 when no run appears in window; negative control proves suppression); AC-3 (diff does NOT alter `git tag -a … -m` — annotated, non-empty body, peels to $RELEASE_COMMIT).
  AC-1: guard requires the `git@github.com:${GITHUB_REPOSITORY}.git` SSH transport and rejects GITHUB_TOKEN; positive control run 29884082483 accepted. AC-2: head_branch filter reproduced against LIVE data — `event=push & head_branch==v0.27.0-pre0` returns exactly run 29845875763 (event=push, conclusion=success), count 1 — so the poll keys on the right run by identity. AC-3: `git tag -a "$PRE0_TAG" "$RELEASE_COMMIT" -m "$PRE0_BODY"` is unchanged CONTEXT in the diff; the tagCmd guard (annotated / non-empty body / targets $RELEASE_COMMIT) passes on the real file.
- DONE: Detached adversarial audit on a throwaway checkout (not the impl worktree) — mutation-test the extended guard (twins (d) credential-revert and (e) verify-neuter must each RED); review the SSH-key setup + poll head_branch filter; confirm 3-files/+68 LOC, edge-advance-scoped, no github.ref rewrite.
  Throwaway detached checkout at ff10ed47. `go test ./internal/release/...` GREEN (34s), gofmt+vet clean. Physically reverting the real release.yml on disk and re-running the on-disk guard: (d) GITHUB_TOKEN https push → RED; (e) neuter the poll's `release.yml/runs` query → RED; plus two of my own: (f) poll present but its `exit 1`→`true` → RED (the twin-uncovered `failsOnMiss` check is non-vacuous); (g) revert to the ACTUAL root-cause cross-repo tap PAT → RED (caught by the SSH-transport requirement, not just the GITHUB_TOKEN strawman). Guard is a tautology on none of the four axes. SSH setup is standard (printf key + ssh-keyscan + IdentitiesOnly) and matches the proven spike form. Scope confirmed: 3 files, +68 net (33+5+30), edits confined to the edge-advance auto-cut step, no `github.ref` rewrite of any job; the goreleaser tap-bump `HOMEBREW_TAP_TOKEN` (release.yml:221-226) is untouched and correct, and no `internal/release` test carries a now-stale tap-PAT pre0 assertion.
- DONE: Note the live confirmation (AC-1 end-to-end) is DEFERRED to the next real stable cut (self-proving via the verify-or-fail guard); semantic adversarial pass; classify findings; recommend PASSED/REJECTED.
  Live confirmation legitimately deferred — the guard fires the pre0 run OR reds edge-advance, so it self-proves at the next cut. Semantic pass done: the poll cannot self-match its stable-tag driver run (head_branch filter excludes it), and it fails LOUD (exit 1), not silent, on a read/parse glitch. No material findings; 3 deferred risks below.

### Findings

Material: NONE. AC-1/AC-2/AC-3 all carry valid, non-self-referential evidence; the regression guard is proven non-tautological.

Deferred risks (do not block the gate):
1. **[evidence/config] Verify-poll Actions-API read rides on the repo staying PUBLIC.** Workflow `permissions: contents: write` gives the edge-advance token `actions: none`; the poll's `gh api .../actions/workflows/release.yml/runs` works today only because an unauthenticated read of this public repo returns HTTP 200 (a `actions:none` token ≥ anon for public read). Trigger→material: repo made private (or GitHub tightening public-repo GITHUB_TOKEN Actions reads) ⇒ poll 403s ⇒ runs=0 ⇒ false `exit 1` on EVERY cut, breaking AC-1. The restricted-token read was never exercised in-CI (the spike used a separate scratch workflow). RECOMMEND a cheap ~1-line hardening — add `actions: read` to the edge-advance job's `permissions:` — before/at the next cut; it removes the visibility dependency entirely. Failure is loud, not silent, so it does not block the gate.
2. **[outcome, edge] Stale-run false-pass on re-cut.** The poll matches ANY run with `head_branch==PRE0_TAG`, not specifically the one THIS push created. Normal flow cuts each pre0 tag once (fresh minor) so no prior run exists; trigger is re-cutting an already-published minor (a prior run lingers in the last 50). Outside the normal promise. Revisit if the release process ever re-cuts an advanced minor.
3. **[operational] EDGE_RELEASE_DEPLOY_KEY provisioning** (ideation BLOCKER) is not verifiable from the repo. If unprovisioned at the next cut, the push/poll fails loud (reds edge-advance) — self-guarding, not a silent gap — but the secret must be confirmed present before relying on AC-1 live.

Polish (optional): the in-file twins cover (a)-(e) but not the `failsOnMiss` check; my (f) mutation proved it reds, so it is not a hole — an in-file (f) twin would simply document it.

### Recommendation

PASSED. All promised value ACs carry valid evidence; the regression guard reds under four independent claim-breaking edits including the actual root-cause credential; scope is as declared (3 files / +68 / edge-advance-scoped / no github.ref rewrite). Live confirmation of AC-1 end-to-end is legitimately deferred to the next stable cut (self-proving via the verify-or-fail guard). Strongly recommend the deferred-risk-1 `actions: read` one-liner as cheap insurance on this high-stakes surface before the next cut; deferred risks 2 and 3 are listed with their revisit conditions.

### Summary

Validated the deploy-key + verify-or-fail fix against AC-1/AC-2/AC-3 and ran a detached adversarial audit on a throwaway checkout. The extended `edge_advance_wiring_test.go` guard is non-tautological: physically reverting the real release.yml reds it on all four claim-breaking axes — GITHUB_TOKEN push (d), neutered poll (e), poll-without-`exit 1` (f), and the actual root-cause cross-repo tap PAT (g). Accepted the documented spike record (methodology sound, credential the sole variable) and reproduced the poll's head_branch filter against live run 29845875763. No material findings; recommend PASSED with three deferred risks — chiefly a cheap `actions: read` hardening since the poll's Actions read currently rides on the repo being public and was never exercised under the restricted in-CI token.
