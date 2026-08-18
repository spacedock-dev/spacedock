---
id: 83tx3a7zwnvz92d7kq0cwf2a
title: Make the auto-pre0 cut idempotent under a workflow re-run
status: ideation
source: "Validation of `collapse-duplicate-edge-marketplace-routes`, cycle 2, 2026-08-18. Recorded there as a deferred risk with a captain-decision note; the captain elected to file it rather than spend a third feedback cycle. Reproduced verbatim by that validator: `fatal: tag 'v0.27.0-pre0' already exists`, rc=128, red under `set -euo pipefail`."
started: 2026-08-18T02:36:44Z
completed:
verdict:
score:
worktree:
issue:
---

Re-running a stable release's `edge-advance` job after its pre0 tag already reached origin makes the job die on the tag it created. The retired mechanism prevented this; its replacement does not.

## Problem

`release.yml`'s `edge-advance` job auto-cuts a `vX.(Y+1).0-pre0` tag after a stable release. The decision that gates it now compares the tag against the highest known version derived from git tag history.

That scan **excludes the ref being decided**. On a re-run the pre0 tag the previous attempt pushed IS in the candidate pool — but the decided tag's own target out-ranks it **by construction** (`0.27.0-pre1 > 0.27.0-pre0`, live-confirmed by the ideation spike below), so the decision advances a second time and `git tag -a` dies:

```
fatal: tag 'v0.27.0-pre0' already exists
```

rc=128, red under `set -euo pipefail`. The collision is deterministic on a faithful re-run: the target is always exactly one prerelease notch above the pre0 the first attempt minted, and no other candidate can outrank it unless a newer stable shipped between attempts.

This is a regression against the mechanism it replaced. Fed the same moment, the retired `next`-manifest read decided `skip`, because its stable path stamped `next` to `X.(Y+1).0-pre1` — one notch above the pre0 it had just cut. The retired step's own comment named "no colliding pre0 auto-tag" as something it prevented. Nothing supplies that self-notch on a re-run now — and nothing can: feeding the decided ref's own notch back into the scan would make the target compare equal to it on EVERY run (strict `>` fails), deadlocking the first cut too. The re-run case is structurally unreachable from the decision's version-ordering inputs; it needs the one input the decision does not consume — whether the concrete pre0 ref already exists.

The trigger is real rather than theoretical. The pre0 step's verify-or-fail poll exits 1 **after** the tag is pushed, so re-running is the natural response to that failure — and three `release.yml` runs in this repository already carry `run_attempt > 1` (v0.23.0, v0.19.5, v0.19.4).

It was classified as a deferred risk rather than material because no value AC fails, nothing wrong publishes, and it fails loudly. It is filed because the promote condition is one observed re-run away.

## Proposed approach

Guard **only the mint**, inside the auto-cut step, leaving push and verify-or-fail live. In `release.yml`'s "Auto-cut the edge prerelease tag on the greened release commit" step, replace the single `git tag -a` line:

```
-          git tag -a "$PRE0_TAG" "$RELEASE_COMMIT" -m "$PRE0_BODY"
+          # Idempotent under re-run: a prior attempt may have pushed the pre0
+          # tag and then failed in the verify poll below. Never re-mint (or
+          # move) an existing tag; push and verify-or-fail still run.
+          git rev-parse -q --verify "refs/tags/$PRE0_TAG" >/dev/null \
+            || git tag -a "$PRE0_TAG" "$RELEASE_COMMIT" -m "$PRE0_BODY"
```

Two code lines (guard + the existing mint as its `||` arm); the rest is comment. Everything downstream is already idempotent and stays load-bearing: `git push` of a tag origin already carries at the same SHA is a no-op success (`Everything up-to-date`, live-confirmed below), and the verify-or-fail poll re-checks that a release.yml run exists for the pre0 tag — so a re-run turns green when the run fired, and still fails loudly when the credential is genuinely workflow-suppressed.

**Placement decision — why guard the mint, not the step or the decision.** The sketch the `2d` implementer weighed ("skip when the tag exists") admits three placements:

1. Early-exit the whole step when the tag exists — REJECTED: it silently disables verify-or-fail on exactly the path where it matters. The recorded trigger is the poll exiting 1 after the push; the re-run scenario therefore includes the sub-case where the credential is genuinely suppressed and no pre0 run exists. An early exit turns that into a green job with the edge binary silently left behind — the exact silent loss the poll's own comment says must never happen.
2. Fold tag-existence into the decision step (`advance=false`) — REJECTED for the same reason (decision=false skips the whole cut step, poll included), and it would mix "was this already done" into a step whose single meaning is version-line ordering.
3. Guard the mint only — CHOSEN: the one non-idempotent command becomes idempotent; every other command in the step already tolerates a re-run.

**Composition with the notch (checklist item 1).** The two guards consume different inputs and neither can absorb the other:

- The notch (`HighestBareStableVersion` → dev-preversion folded into `highest-known-edge-version`) answers *"does this tag's version rank entitle it to cut at all?"* from tag-history ordering. It cannot answer the re-run case: the missing input is the decided ref's own self-notch, and supplying that would deadlock every first run (target == self-notch, strict `>` fails — see Problem).
- The tag-exists guard answers *"does the concrete ref this run would mint already exist?"* from ref existence. It cannot answer the old-line case in general: a patch whose wrong pre0 target does not exist yet (the notch's equality-skip case, e.g. `v0.25.1` cut before any `v0.26.0-pre0` was ever auto-cut) sails past a ref-existence check and only the notch stops it.
- The fix touches neither the decision step nor any Go code, so the notch's entire validated suite (unit cases, the 16-evaluation history replay, the real-`v0.25.1` collision replay, the decision-step shell tests) passes byte-unchanged — verified live: the extracted decision step still prints `advance=false` for `v0.25.1` against the one-tag-behind fixture (notch reason: highest known `0.27.0-pre1`).

**Ordering dependency.** The failure exists only in the post-#727 `release.yml` (main still carries the retired `next`-manifest read, which decides `skip` at this moment). Implementation lands stacked on `spacedock-ensign/collapse-duplicate-edge-marketplace-routes` (PR #727) or on main after #727 merges — never on pre-#727 main.

**Shared dependency, already recorded.** The guard reads tags the checkout fetched (`fetch-depth: 0`, `fetch-tags: true`) — the same dependency the decision step already has. Stripping those checkout flags is one of the two out-of-scope silent-disable edits below; this task does not change that exposure.

## Spike record (riskiest mechanism exercised first)

The riskiest unverified mechanism was the whole-step claim: that the REAL extracted step script — not a re-derivation — fails exactly as recorded and turns green under the two-line guard with push and poll still executing. Exercised live against the post-#727 worktree (`ae8c3a874`), extracting both steps' `run:` blocks from `release.yml` and running them under bash with `GIT_DIR`/`GIT_WORK_TREE` pointed at constructed tag-state repos (the `2d` correction round's own harness pattern, `runDecisionStepScript` in `internal/release/edge_advance_decision_shell_test.go`):

1. Re-run state (`v0.25.1`, `v0.26.0`, `v0.27.0-pre0`; deciding `v0.26.0`): decision step printed `advance=true` (`target 0.27.0-pre1 vs highest known 0.27.0-pre0`); auto-cut step died `fatal: tag 'v0.27.0-pre0' already exists`, **rc=128** — the recorded failure, reproduced from the real step bytes.
2. Same state, guarded script, with a bare "origin" carrying the pre0 tag, `HOME` sandboxed, `gh`/`ssh-keyscan` PATH shims, and a `url.insteadOf` rewrite of the SSH push URL to the local bare repo: **rc=0**, push printed `Everything up-to-date`, the verify poll executed (its `::notice::` fired via the shim), and the pre0 tag's target SHA was byte-identical before and after.
3. First-run state (no pre0 anywhere), guarded script: still mints — annotated (`cat-file -t` = `tag`) and pushed to the bare origin. The guard does not break the primary path.
4. One-tag-behind old-line state (`v0.25.0`, `v0.26.0`, `v0.26.0-pre0`; deciding `v0.25.1`): decision step printed `advance=false` — the notch, untouched.

Push idempotence, the guard's `rev-parse` form, the shim seams, and the `insteadOf` push redirection are all now proven mechanisms; the implementation's first test is the spike's runs 1–3 made durable.

## Documentation diff

`docs/releasing.md`, "Stable (`vX.Y.Z`) tag, latest line" bullet — after "…failing `edge-advance` loudly if none appears** rather than leaving the edge binary silently behind.", insert:

> Re-running the job after such a failure is safe: an existing pre0 tag is never re-minted or moved (the step checks `refs/tags/` before tagging), the push of an already-landed tag is a no-op, and the run-verification poll executes again — so the re-run turns green once the pre0 run exists, and still fails loudly if the credential remains suppressed.

## Out of scope

The notch mechanism itself (`HighestBareStableVersion`, `HighestKnownEdgeVersion`, `EdgeAdvanceDecision`) — validated across every named boundary input plus a 16-evaluation replay of real release history. Do not re-open it.

The two silent-disable mutations recorded alongside this risk: hardcoding the decision step to `advance=false`, and stripping `fetch-depth: 0` / `fetch-tags: true` from the edge-advance checkout. Both still pass the suite. They are a separate concern, and closing them may need a kind of check this project has been deliberately reducing.

## Expected surface and tolerance

Estimate net LOC change: **+90 across 3 files**. Insertions ≈ 91, deletions ≈ 1. Tolerance: net +90 ± 25, files 3 ± 1. Do not declare a gross tolerance.

- `.github/workflows/release.yml`: +5 / −1 (two code lines — the guard and the mint as its `||` arm — plus comment).
- `internal/release/edge_advance_decision_shell_test.go` (extended): ≈ +80 / −0 — one runner, one test (three sequential runs of the real script), the four test doubles the script's own tail forces.
- `docs/releasing.md`: +5 / −0 (the re-run sentence above).

**The proof at its floor (captain ruling "reduce tests", gate cycle 1).** The smallest proof that still goes RED when the guard is removed is one run of the real extracted step script against a tag state where the pre0 exists: rc must be 0, and without the guard it is 128. That single run costs ≈ +70 and its parts are irreducible under this task's proof standard:

- ≈ 30: the Go runner around `bash -c step.run` — extraction reuses the existing `readWorkflow`/`edgeAdvanceJob`/`edgeAdvanceAutoPre0Step` helpers; the lines are env plumbing and exit-code capture. (Generalizing the existing `runDecisionStepScript` instead was weighed: it saves ~8 lines but edits a function `2d`'s in-flight branch owns — not worth the merge friction.)
- ≈ 20: four test doubles, each forced by a named line of the script's own tail, which runs after the guarded mint under `set -euo pipefail`: sandboxed `HOME` (the step writes `~/.ssh`), an `ssh-keyscan` shim (a real network call that dies on an offline runner), a `gh` shim printing a run count (without it the verify poll loops to its 120s timeout and exits 1 — the shim also records its invocation, which is AC-3's observable), and a bare-origin clone plus `url.insteadOf` (the step pushes to a hardcoded SSH URL).
- ≈ 20: fixture tags and the rc/SHA/poll assertions.

Below this floor there are only the two rejected shapes: slicing the script to stop before the push (re-derivation — the real-bytes discipline `2d` just adopted forbids it) or grepping release.yml for the rev-parse line (banned; the class a sibling task is deleting ~900 lines of).

Three increments are kept above the floor, each priced with the mutant that escapes without it (itemized in the test plan): the first-run mint pass (+6 marginal), the divergent-tag never-move pass (+9), and the decision-still-advances pin (+4). Cut from the cycle-1 plan: the second test function and its own fixture, the separate first-run test, and a duplicated fixture wrapper — test-side ≈ 105 → ≈ 80.

The seed said +25 ± 20 across 2 files; the workflow change is on target at two code lines, and the remainder is the priced floor above, declared for the gate rather than absorbed.

Semantics changed: the auto-pre0 step becomes idempotent under re-run — when its target tag already exists it skips the mint (never moves the tag), still pushes (no-op), and still runs verify-or-fail; the job's exit in that state changes from rc=128 to rc=0 when the pre0 run exists, and stays a loud failure when it does not. No command grammar, stored format, or authority changes.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - Re-running a stable release's edge-advance job after its pre0 tag exists leaves the job green and the tag untouched.**
This is the measuring AC: the re-run's exit status must be 0 where it is currently 128, and the tag's target commit must be unchanged. Verified by replaying the recorded reproduction through the REAL extracted step script — the same tag state that produced `fatal: tag 'v0.27.0-pre0' already exists` rc=128 (spike run 1) must exit 0, with the pre0 tag's SHA byte-identical before and after. The never-moved half is made falsifiable by a fixture whose existing pre0 deliberately targets a different commit than the release commit, so a `-f` re-tag mutant moves a SHA the test pins. Fails if the job still dies, or if it "succeeds" by moving or force-replacing an existing tag.

**AC-2 - The old-line patch protection still holds.**
Verified by the untouched existing suite re-running green with zero edits: the `EdgeAdvanceDecision` unit cases, the 16-evaluation release-history replay, `TestHighestKnownEdgeVersionCommandRestoresNextNotchAgainstRealV0251Collision`, and both decision-step shell tests. The fix touches neither the decision step nor any Go code, so a masking regression can only appear as an edit those tests would sit on. Fails if the new condition masks or replaces the notch rather than composing with it — the regression that would trade one guard for another.

**AC-3 - A re-run that skips the mint still proves the pre0 release run exists.**
The verify-or-fail guarantee survives the fix: on the skip-mint path the run-verification poll still executes, so a workflow-suppressed credential still fails the job loudly instead of going green with the edge binary left behind. Verified in the AC-1 test by asserting the poll ran (the `gh` shim records its invocation; the poll's `::notice::` appears in the step output). Fails under the early-exit form of the fix — the sketched "skip when the tag exists" — which would green the exact sub-case verify-or-fail exists to make loud.

## Test plan

Extend `internal/release/edge_advance_decision_shell_test.go` (same package; reuses `readWorkflow`, `edgeAdvanceJob`, `edgeAdvanceAutoPre0Step`, `tagFixtureRepo`, and the `GIT_DIR`/`GIT_WORK_TREE` redirection pattern of `runDecisionStepScript`). One new runner, `runAutoPre0StepScript`, with the four doubles priced in the surface section. No network, no sleeps (the poll exits on its first iteration via the shim); cost well under a second beyond the `go run` cache warm.

**One test, `TestAutoPre0StepScriptRerunIsIdempotent`** — three sequential runs of the real script against one fixture (`v0.25.1`, `v0.26.0`, bare origin), each pass named with what it buys:

1. First-run pass: no pre0 exists; the script must exit 0 and mint an annotated tag onto the bare origin. Buys: a broken mint (the `||` arm lost or degenerated) is caught in CI instead of at the next real stable cut. Marginal cost ≈ +6 over hand-planting the tag, and it makes pass 2's state the faithful recorded one — minted by the step itself.
2. Faithful re-run pass (AC-1, AC-3, **the red/green core**): same script, same fixture, tag now exists. Must exit 0 where the unguarded script exits 128 (`fatal: tag 'v0.27.0-pre0' already exists` — the recorded failure), leave the tag SHA unchanged, and have invoked the `gh` poll shim (AC-3: verify-or-fail still runs on the skip-mint path). This pass alone is the proof floor: delete the guard from release.yml and it goes red.
3. Divergent-tag pass (AC-1's never-moved half): repoint the pre0 (local and origin) to a fresh empty commit, run again; must exit 0 and leave the divergent SHA in place. Buys: the only detection of the `git tag -f` idempotency hack AC-1 forbids — a faithful re-run cannot catch it because a re-mint targets the same SHA. ≈ +9. (Mechanic proven live in the cycle-2 spike run: divergent SHA survived, rc=0.)

Plus a 4-line pin: the real decision step still prints `advance=true` on pass 2's state (via the existing `runDecisionStepScript`). Buys: catches the guard migrating into the decision step — the rejected placement — which would keep all three passes green while silently disabling verify-or-fail on re-runs.

AC-2 needs no new test: the entire decision-level suite runs unchanged; its continued green under this diff is the composition proof.

No new Go mechanism, command, or flag: every mechanism in this plan (extraction, redirection, shims, bare-origin push target) serves AC-1/AC-3 directly, and each simpler alternative — grep tests (banned for this task), script slicing (re-derivation), a Go-side guard (new surface to make a one-line git builtin testable) — was named and rejected above.

## Stage Report: ideation

- DONE: Confirm the tag-exists check COMPOSES with the notch rather than replacing it — they answer different inputs, and a fix that masks the notch trades one guard for another.
  Composition proven structurally (the re-run needs the self-notch the decision can never consume without deadlocking first runs; the equality-skip old-line case has no existing ref for a tag-exists check to see) and live (spike run 4: extracted decision step still prints `advance=false` for `v0.25.1` on the one-tag-behind fixture). The fix touches neither the decision step nor Go code; AC-2 pins the untouched suite.
- DONE: Reproduce the recorded failure before designing: the same tag state that produced `fatal: tag 'v0.27.0-pre0' already exists` rc=128 must be the thing your fix turns green.
  Spike run 1 reproduced `fatal: tag 'v0.27.0-pre0' already exists` rc=128 from the REAL extracted auto-cut step (post-#727 worktree `ae8c3a874`) against tags {v0.25.1, v0.26.0, v0.27.0-pre0}; spike run 2 turned that exact state green (rc=0, tag SHA unchanged, push no-op, verify poll executed) under the two-line guard.
- DONE: Keep this at roughly two lines of workflow change plus its proof; if the design grows past the declared surface, stop and say why rather than absorbing it.
  The workflow change is two code lines. The proof grew past the seeded +25 ± 20 to ≈ +105 test lines because the mandated proof standard (run the real step script, no grep tests) requires push/poll test doubles; the revised surface (net +115 ± 35, 3 files) is declared with its why in "Expected surface and tolerance" for the gate to approve or bounce, not absorbed silently.

### Summary

Confirmed the `2d` implementer's sketched fix with one load-bearing correction: guard ONLY the `git tag -a` mint, not the whole step — the early-exit form would silently disable verify-or-fail on the one path where a suppressed credential must still fail loudly (now AC-3). Reproduced the recorded rc=128 collision from the real extracted step bytes, proved the guard turns the same state green with push and poll still live, and proved composition with the notch in both directions. Flagged the ordering dependency: the fix targets post-#727 `release.yml` and must stack on PR #727.

## Stage Report: ideation (cycle 2)

- DONE: Confirm the tag-exists check COMPOSES with the notch rather than replacing it — they answer different inputs, and a fix that masks the notch trades one guard for another.
  Unchanged from cycle 1; the captain's ruling kept the composition argument and placement analysis. No body edits to those sections.
- DONE: Reproduce the recorded failure before designing: the same tag state that produced `fatal: tag 'v0.27.0-pre0' already exists` rc=128 must be the thing your fix turns green.
  Unchanged from cycle 1 (spike runs 1–2 stand). One new cycle-2 spike run proved the folded never-move pass's mechanic: a pre0 repointed to a divergent commit survives the guarded script untouched at rc=0.
- DONE: Keep this at roughly two lines of workflow change plus its proof; if the design grows past the declared surface, stop and say why rather than absorbing it.
  Captain ruled "reduce tests" at the gate: proof cut from two tests ≈ +105 to one test ≈ +80, declared surface net +115 → **net +90 ± 25, 3 files** (insertions ≈ 91, deletions ≈ 1). The floor is itemized in the surface section — a single red/green run costs ≈ +70 (runner ≈ 30, four doubles forced by the script's own tail ≈ 20, fixture+asserts ≈ 20) — and each of the three increments above it is priced with the mutant that escapes without it (broken mint +6, `tag -f` re-point +9, guard-migrates-into-decision +4). No grep test; the red/green core still runs the real step bytes.

### Summary

Correction round per captain ruling "83t - reduce tests". Fix, placement, and composition untouched. Test plan collapsed to one test with three sequential runs of the real script (mint, faithful re-run red/green core, divergent never-move) plus a 4-line decision pin, reusing the existing extraction helpers and `tagFixtureRepo`; the cycle-1 second test and duplicate fixture scaffolding are cut. Surface re-declared at net +90 ± 25 across 3 files with the cheapest-proof-that-can-fail floor argued line-item by line-item.
