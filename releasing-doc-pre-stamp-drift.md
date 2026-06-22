---
id: 7yd3mbsy2am5qggc17sxvz2v
title: Reconcile docs/releasing.md to the live release machinery (one consolidated pass)
status: implementation
sprint: 0230-stable-finalization
source: "0202 Commander drive (2026-06-13). releasing.md 'Cutting a Stable Release' step 3 says to stamp+commit the version before tagging, but v0.20.1 and v0.20.2 both tagged the gated commit directly and let release.yml stamp post-tag (a pre-stamp creates an ungated commit the exact-SHA e2e-gate blocks on)."
group: cleanup
started: 2026-06-21T06:05:14Z
worktree: .worktrees/spacedock-ensign-releasing-doc-pre-stamp-drift
---

`docs/releasing.md`'s "Cutting a Stable Release" procedure (step 3) documents a manual pre-stamp commit before the annotated tag. Actual practice (v0.20.1, v0.20.2) tags the gated commit directly; release.yml stamps the plugin manifests post-tag.

## Problem

The `e2e-gate` job binds the cut to a green Runtime Live E2E run for the EXACT tagged commit SHA. A manual pre-stamp (releasing.md step 3) creates a NEW commit with no green run, which the gate would block. So both recent cuts tagged the already-gated `main` HEAD directly and relied on release.yml's post-tag stamp + the moving `stable` ref. The written procedure and the real procedure diverge — a fresh cutter following step 3 literally would stamp, create an ungated commit, tag it, and be blocked by the e2e-gate.

## Consolidated 0230 member

This is the SINGLE consolidated `docs/releasing.md` reconciliation work item. It absorbs the in-0230 substance of `stamp-then-tag-release-ritual` (the manifest==tag guard → AC-2 below) and `steady-state-stable-release-runbook` (the next->main steady-state runbook → the approach below). All three edit `docs/releasing.md` and its byte-identical mirror `docs/site/contributing/releasing.md`, so they MUST run as one serialized pass, not in parallel — parallel filing collides on the same file.

## Proposed approach

ONE pass over docs/releasing.md (+ its byte-identical mirror docs/site/contributing/releasing.md) reconciling it to the LIVE release machinery, plus ONE in-repo guard that makes the doc's central invariant testable. Premise re-check first changed the scope — see notes; the doc, not the machinery, is the thing that is wrong.

1. Fix the dangerous step 3 (releasing-doc-pre-stamp-drift). Reframe "Cutting a Stable Release" so the cutter tags a commit that ALREADY has a green Runtime Live E2E run for its exact SHA. The e2e-gate binds `run.HeadSha == tagged-commit-SHA` (internal/release/e2egate.go), and Runtime Live E2E is workflow_dispatch-only — so any fresh stamp+commit produces an UNGATED SHA the gate BLOCKS. Rewrite the procedure to: stamp `plugin.json` + `.codex-plugin/plugin.json` to X.Y.Z, commit, push that commit to `main`, get a green Runtime Live E2E run FOR THAT COMMIT (dispatch the workflow on it / confirm green), THEN annotate-tag that same green SHA. State the gate contract explicitly so the cutter knows WHY the tagged commit must be pre-greened.

2. Pin the stamp-then-tag ritual + drop the manual-repoint fiction (stamp-then-tag-release-ritual). The doc currently says the stable marketplace entry is "repointed to the released tag" by a manual commit in the standalone marketplace repo (L21-22, L49-50). That is FALSE: release.yml advances stable as a MOVING BRANCH via `git push origin main:refs/heads/stable` (L221-229), and origin/stable already tracks the v0.22.0 stamp commit. Rewrite "What the Tag Push Does" + step 3's marketplace note to describe the moving-`stable`-branch auto-advance, and describe the post-tag stamp step as idempotent (no-op when the tagged commit already carries the stamp — which it now does). Note the stamp-then-tag inversion is already de-facto reality (v0.20.1..v0.22.0 all tag a matching-manifest stamp commit); the doc just never said so.

3. Fold the next->main steady-state runbook (steady-state-stable-release-runbook). Where the doc claims "Stable releases are cut from `main`" and step 1 says "ensure all content is merged to `main` / worktree off `main`", reconcile to the post-flip integration model only IF the live machinery confirms next-integration; the live doc is already main-centric and origin/stable advances from main, so VERIFY before rewriting — if dev still integrates on `next`, document the green-`next`-tip -> advance-`main` step; if integration is genuinely on `main`, keep main-centric and delete the stale next-integration framing instead. Do not invent a next-line if the machinery doesn't have one.

4. Add ONE in-repo guard enforcing the doc's invariant (the stamp-then-tag AC needs a real, divergeable check, not a prose-grep). A test/lint over the cut machinery asserting `plugin.json` version == the release tag's semver for the tagged commit. This is the enforceable proof behind the doc.

Out of scope (flag, do not attempt): the dead `next` marketplace `version` field + tag-vs-branch channel structure live in the standalone spacedock-dev/marketplace repo, unreachable from this repo, so stamp-then-tag's AC-2 cannot be enforced by an in-repo test — note it in the doc and leave the field cleanup to the marketplace-repo task. Edit BOTH releasing.md copies identically (they are byte-identical today).

## Verify integration model first

`docs/releasing.md:3` already says stable releases are "cut from `main`" and origin/stable advances from `main` (release.yml:229), which CONTRADICTS the `steady-state-stable-release-runbook` premise that "dev integrates on `next`". So the steady-state fold is most likely a DELETION of stale next-framing, not an addition of a next->main step. Confirm against the live machinery before rewriting the runbook section — do not document a next-line the machinery doesn't have.

## Acceptance criteria

- **AC-1 (VALUE)** — A cutter following the reconciled releasing.md produces a commit the e2e-gate ACCEPTS. Verified by: simulate the documented procedure literally and run `spacedock-release e2e-gate <tagged-commit-sha>` (or the EvaluateE2EGate predicate) against a run-list fixture; the reconciled procedure's tagged SHA matches a conclusion:success Runtime Live E2E headSha (PASS), whereas the CURRENT step-3 procedure (stamp+commit-then-tag the new SHA) yields a SHA with no matching green run (BLOCK). Independent baseline that can move the wrong way: the gate is `run.HeadSha == releaseCommit` in internal/release/e2egate.go; greening the wrong SHA blocks the cut.
- **AC-2** — A guard/test asserts the tagged commit's plugin.json version equals the tag's semver. Verified by: an in-repo test over the cut machinery (or a fixture release) comparing two INDEPENDENT values — the git tag semver and the manifest version — that genuinely disagree on the v0.20.0 ordering (tag 0.20.0 -> manifest 0.19.9, FAILS) and agree from v0.20.1 onward (PASSES). Not a tautology: the values come from different sources and the v0.20.0 history proves they can diverge.
- **AC-3** — releasing.md describes the moving-`stable`-branch auto-advance, not a manual repoint-to-tag. Verified by: the reconciled doc states release.yml advances stable via `git push origin main:refs/heads/stable` and the post-tag stamp is idempotent; cross-checked against release.yml L221-229 and `git ls-remote origin refs/heads/stable` (currently f96f4a1a = v0.22.0 stamp commit). The prior 'repointed to the released tag' / 'commit in the standalone marketplace repo' text is gone.
- **AC-4** — docs/releasing.md and docs/site/contributing/releasing.md remain byte-identical after the pass. Verified by: `diff docs/releasing.md docs/site/contributing/releasing.md` is empty (or the project's existing mirror-sync check passes).
- **AC-5** — The dead `next` marketplace version field + channel-structure cleanup is explicitly scoped OUT (standalone marketplace repo, unreachable here) and noted in the doc as deferred to the marketplace-repo task, rather than silently dropped.

## Notes
Pure doc reconciliation (plus the one new AC-2 guard); the machinery is already correct. Surfaced by the 0202 cut.

## Stage Report: implementation

- DONE: Reconcile docs/releasing.md so a cutter tags the e2e-GREEN commit, not a fresh pre-stamp SHA; step 3 stops pre-stamping a new commit and describes the moving-`stable`-branch auto-advance.
  Rewrote step 3 to stamp+commit+push to main, GREEN that exact commit via `gh workflow run "Runtime Live E2E"`, then `git tag -a vX.Y.Z $(git rev-parse origin/main)`; "What the Tag Push Does" now states the gate binds `run.HeadSha == tagged-commit` and Runtime Live E2E is workflow_dispatch-only. Code SHA c97ec639.
- DONE: Add a divergeable guard the cutter can actually fail — assert the tagged commit's plugin.json version == the tag semver.
  New `spacedock-release manifest-tag-gate <tag> <plugin.json>...` + pure predicate `EvaluateManifestTagGate`. Exercised against REAL history: v0.20.0 (manifest 0.19.9) exits 1/BLOCK, v0.22.0 (manifest 0.22.0) exits 0/PASS. 11 new tests, all green.
- DONE: VERIFY FIRST (record the finding) — state exactly what was removed vs added in the steady-state fold; `next` marketplace version-field cleanup OUT of scope.
  CONFIRMED main-centric: live doc already said "cut from `main`" and origin/stable (f96f4a1a) advances from main via `git push origin main:refs/heads/stable` (release.yml L229). The fold was a DELETION of the false "repointed to the released tag" / "commit in the standalone marketplace repo" prose (old L19-23, L49-50) + added the moving-`stable` + idempotent-stamp description. No `next`->main step invented (machinery has none). Dead `next` marketplace version field noted as deferred to the marketplace-repo task (AC-5).

### Summary
Reconciled both releasing.md views to the live e2e-gate machinery in one pass and added the AC-2 manifest/tag guard. Key finding: `docs/site/contributing/releasing.md` is a SYMLINK to `docs/releasing.md` (not a byte-copy), so AC-4 byte-identity holds by construction — editing the target updates both; the mirror never appears in `git status`. AC-1 proven by exercising `EvaluateE2EGate`: tagging the greened SHA PASSES, tagging a fresh ungated post-stamp SHA BLOCKS. AC-2 proven against real release history. Full repo test suite green (0 failures).

## Stage Report: validation

- DONE: Run applicable tests — `go test ./cmd/spacedock-release/... ./internal/release/...` and full `go test ./...`.
  Both green, 0 failures across all 16 packages (e.g. internal/cli 41s, internal/release 7.8s ok); 11 new manifest-tag-gate tests pass.
- DONE: Verify AC-1 (core bug fix) — reconciled procedure GREENs the exact main commit then tags THAT e2e-green SHA; reproduce the EvaluateE2EGate exercise.
  Exercised `EvaluateE2EGate` with one green run for the greened SHA: tag-the-greened-SHA → pass=true ("run 42 matches release commit"); tag-a-fresh-ungated-pre-stamp-SHA → pass=false ("no conclusion:success run found"). Ship path is gated on the tagged SHA (release.yml:124 `git rev-list -1 $GITHUB_REF_NAME`, goreleaser `needs: e2e-gate`).
- DONE: Verify AC-2 (divergeable, non-tautological) — manifest-tag-gate / EvaluateManifestTagGate against REAL history.
  Reproduced via the CLI against the actual manifests at the real tags: v0.20.0 (real plugin.json 0.19.9, both manifests) → exit 1/BLOCK; v0.22.0 (0.22.0) → exit 0/PASS. Tag semver and manifest version are independent artifacts that demonstrably diverged at v0.20.0 — not a tautology.
- DONE: Verify AC-3 — false "repoint to the released tag"/"standalone marketplace repo" prose DELETED; moving-stable auto-advance + idempotent stamp documented.
  Removed phrases grep to NONE in docs/releasing.md; moving-`stable` (`git push origin main:refs/heads/stable`) + idempotent-stamp prose present (L24-33). Cross-checked release.yml L229 + the idempotent `git diff --quiet` skip (L213); live `git ls-remote origin refs/heads/stable` = f96f4a1a = v0.22.0 SHA, as the AC asserts.
- DONE: Verify AC-4 — both releasing.md views byte-identical.
  `docs/site/contributing/releasing.md` is a git-tracked symlink (mode 120000 → `../../releasing.md`); `diff` is empty. Byte-identity holds by construction.
- DONE: Verify AC-5 — dead `next` marketplace version field deferred OUT, noted in the doc.
  Deferral note present at docs/releasing.md L157-160 ("deferred to a marketplace-repo task and is NOT part of this flow").
- DONE: Run the DETACHED adversarial audit over the reconciled doc + the new guard on a throwaway checkout.
  Ran on a throwaway clone (/tmp, since removed). The two pure predicates are individually solid; the audit surfaced wiring/prose gaps recorded below. Cleaned up; no audit artifacts touched the worktree or state checkout.
- FAILED: AC-2 as "enforceable proof behind the doc" (proposed-approach item 4) — the guard is advisory-only.
  `grep -rln manifest-tag-gate .github/workflows/` → ZERO. The gate's only caller is a manual `go run` in docs/releasing.md step 3 (L76). e2e-gate is a real CI job goreleaser `needs:` (release.yml:99,128) with a workflow guard test; manifest-tag-gate has neither. A cutter who skips step 3's line tags a manifest-mismatched commit and nothing blocks the ship — the exact v0.20.0 inversion (tagged a commit whose plugin.json read 0.19.9) recurs undetected.
- FAILED: Doc deliverable carries a broken markdown artifact.
  docs/releasing.md ends with a stray unbalanced ``` code fence on its last line (L163); `grep -c '```'` = 19 (odd). The `## Notes` section is plain bullets with no opening fence, so L163 closes nothing. Affects the mirror via the symlink.

### Feedback Cycles

- **Cycle 1 — validation REJECTED (detached adversarial audit, general-purpose agent on throwaway checkout).** Two ACs land but two material defects remain in the deliverable. Recommendation: REJECTED → route back to implementation.
  1. **BLOCKER (AC-2 intent) — manifest-tag-gate is advisory-only, not CI-enforced.** Add a `manifest-tag-gate` step to the `e2e-gate` (or a sibling gate) job in release.yml — `go run ./cmd/spacedock-release manifest-tag-gate "$GITHUB_REF_NAME" .claude-plugin/plugin.json .codex-plugin/plugin.json` — with goreleaser `needs:` it, PLUS a workflow guard test mirroring `internal/release/e2egate_workflow_test.go`. Carry the same `if: !contains(github.ref, '-')` pre-release skip the stamp step uses (release.yml:204), else a `v0.23.0-pre.1` tag self-blocks against a `0.23.0` manifest (audit confirmed exit 1 on that pairing).
  2. **MINOR — stray ``` code fence at docs/releasing.md:163.** Delete line 163 (the last line). Mirror updates via the symlink.
  3. **MINOR (consider) — step-6 stale-ref trap.** docs/releasing.md L110/L114 tag `$(git rev-parse origin/main)`. The audit's git experiment showed this flips safe↔unsafe on whether an out-of-band `git fetch` ran (step 5's changelog line tempts one). Suggest step 4 capture the greened SHA (`REL_SHA=$(git rev-parse HEAD)`) and step 6 tag `$REL_SHA`, removing the `origin/main` indirection. The e2e-gate still backstops the ship, so this is a clarity fix, not a ship-safety hole.
  - Out of scope (machinery, not this doc task — noted, not blocking): audit MAJOR that release.yml advances `stable` to `main` HEAD at stamp-time rather than the tagged SHA (release.yml:229), and MINOR that `SPACEDOCK_E2E_GATE_WAIVER` is a sticky repo var with no clear-after step. The task scope is the doc + the one in-repo guard; these belong to the marketplace/release-machinery surface and should be filed separately.

- **Cycle 1 — captain ruling (2026-06-21):** REJECT confirmed; route back to implementation. Captain folded the out-of-scope stable-advance MAJOR into this rework (one release.yml pass). Fixes:
  1. Wire `manifest-tag-gate` into release.yml's `e2e-gate` job + goreleaser `needs:` it + a workflow guard test (mirror `e2egate_workflow_test.go`); carry the `if: !contains(github.ref, '-')` pre-release skip.
  2. Delete the stray ``` fence at docs/releasing.md:163.
  3. Step-6 clarity: capture `REL_SHA=$(git rev-parse HEAD)` in step 4 and tag `$REL_SHA` (drop the `origin/main` indirection).
  4. FOLDED IN (the audit MAJOR): release.yml advances `stable` to the TAGGED SHA, not main HEAD at stamp-time (release.yml:229) — fix it so `stable` and the tag point at the same commit.

## Stage Report: implementation (cycle 2)

- DONE: BLOCKER — WIRE manifest-tag-gate into CI so it actually blocks.
  Added a "Gate the cut on the tagged manifest matching the tag semver" step to release.yml's `e2e-gate` job (goreleaser already `needs: e2e-gate`), running `manifest-tag-gate "$GITHUB_REF_NAME" .claude-plugin/plugin.json .codex-plugin/plugin.json` with the `if: "!contains(github.ref, '-')"` pre-release skip. New `internal/release/manifest_tag_gate_workflow_test.go` mirrors e2egate_workflow_test.go: red-then-green wiring guard + 3 adversarial twins (dropped gate step, dropped needs edge, missing pre-release skip). Code SHA e32c27f5.
- DONE: MINOR — delete the stray ``` code fence at docs/releasing.md:163.
  Removed; fence count now even (18), doc renders clean. Mirror updates via the symlink.
- DONE: Step-6 clarity — capture REL_SHA and tag it, dropping the origin/main indirection.
  Step 4 now captures `REL_SHA=$(git rev-parse HEAD)`; step 6 tags `"$REL_SHA"` with a note that this survives a stray `git fetch` moving origin/main.
- DONE: FOLDED IN — advance `stable` to the TAGGED commit SHA, not main HEAD at stamp-time.
  release.yml stamp step now resolves `RELEASE_COMMIT="$(git rev-list -1 "$GITHUB_REF_NAME")"` before switching to main and pushes `"$RELEASE_COMMIT:refs/heads/stable"`. Hardened the existing stable-ref guard (renamed `stableRefPushSource`, `TestStampStepAdvancesStableRefToTaggedCommit` asserts the source is `$RELEASE_COMMIT`) + adversarial twin `TestStableRefGuardRejectsMainSource` rejecting the divergeable `main:refs/heads/stable` form.

### Summary
Closed the cycle-1 reject: the manifest-tag-gate is now CI-enforced in the e2e-gate job goreleaser needs, so a tag/manifest mismatch blocks the cut rather than relying on an advisory manual step — with the pre-release `if` skip so `-pre.N` tags don't self-block. Also captured REL_SHA for the tag target, removed the stray fence, and bound the moving-`stable` advance to the tagged SHA (the folded-in audit MAJOR), each backed by a red-then-green guard. Proven by exercising the binary: exit 1 on v0.20.0/0.19.9, exit 0 on v0.22.0/0.22.0, exit 1 on v0.23.0-pre.1/0.23.0 (the self-block the skip avoids). release.yml is valid YAML; full repo test suite green (0 failures). NOTE: the release.yml edit needs the `workflow` OAuth scope to push — committed locally on the worktree branch; the FO pushes via the scoped credential helper.

## Stage Report: validation (cycle 2)

- DONE: Verify the BLOCKER fix is real CI enforcement — manifest-tag-gate is a STEP in release.yml's e2e-gate, goreleaser `needs:` it, carries the pre-release skip; reproduce the workflow guard red-then-green.
  As-shipped release.yml:134-138 runs the gate in e2e-gate (set -euo pipefail, last command → fatal), goreleaser `needs: e2e-gate` (L141), `if: "!contains(github.ref, '-')"` skip (L135). Reproduced red-then-green by mutating the REAL on-disk release.yml: dropped step → guard RED; dropped `needs: e2e-gate` → RED; dropped pre-release `if` → RED; restored → all GREEN, file byte-identical (empty `git diff`). Guards read the file from disk (readWorkflow/readReleaseWorkflow), so a real edit reds CI.
- DONE: Verify the folded-in stable-advance fix — release.yml advances stable to `$RELEASE_COMMIT:refs/heads/stable`, not main HEAD; stable-ref guard + adversarial twin reject the `main:` form.
  release.yml:249 pushes `"$RELEASE_COMMIT:refs/heads/stable"`; RELEASE_COMMIT resolved via `git rev-list -1 "$GITHUB_REF_NAME"` BEFORE `git switch main` (L225). TestStampStepAdvancesStableRefToTaggedCommit asserts src==`$RELEASE_COMMIT`; regressing the real file to `main:refs/heads/stable` reds it. Both green at HEAD.
- DONE: Verify the rest + re-audit — stray fence deleted; step-4 REL_SHA capture + step-6 tag; binary BLOCKS mismatch / PASSES match / `-pre.N` skip rationale; full `go test ./...`; DETACHED adversarial audit on a throwaway checkout.
  Fence count even (18), no trailing fence; mirror is a symlink → `diff` empty (AC-4). Step 4 captures `REL_SHA=$(git rev-parse HEAD)`, step 6 tags `"$REL_SHA"`. Binary: v0.20.0/0.19.9 exit 1, v0.22.0/0.22.0 exit 0, v0.23.0-pre.1/0.23.0 exit 1 (why the `if` skip exists; a hyphen-free FINAL tag with a mismatched manifest is STILL gated — confirmed v0.23.0/0.22.0 exit 1). `go test ./...` exit 0. Detached audit ran on a throwaway clone of the branch tip in /tmp (since removed); findings below.
- FAILED: AC-3 — the reconciled doc must describe the moving-`stable` auto-advance matching the code; the doc cites the wrong (since-removed) command.
  docs/releasing.md:25 states release.yml runs `git push origin main:refs/heads/stable` — but cycle-2 (the folded-in captain item 4) changed the code to `git push origin "$RELEASE_COMMIT:refs/heads/stable"` (release.yml:249) precisely to stop stable tracking main HEAD, and TestStableRefGuardRejectsMainSource exists to reject the `main:` form. The L25 prose was authored in cycle-1 (c97ec639) when the code DID push `main:`, and cycle-2 changed the code but left the doc. The deliverable's doc now contradicts its own code change.

### Feedback Cycles

- **Cycle 2 — validation REJECTED (detached adversarial audit, general-purpose agent on throwaway /tmp clone).** The four captain-ruled items are all genuinely closed IN THE SHIPPED FILE: the gate is real, fatal CI enforcement (audit confirmed it is NOT bypassable as-shipped); stable advances to the tagged SHA; the fence/REL_SHA/binary all check out; the pre-release skip is correct and necessary. One concrete in-deliverable defect blocks; the rest are guard-hardening findings, not as-shipped holes.
  1. **BLOCKER (AC-3 / captain item 4) — doc contradicts the code it documents.** docs/releasing.md:25 says release.yml runs `git push origin main:refs/heads/stable`; the code does `$RELEASE_COMMIT:refs/heads/stable`. This is the exact `main:` form the folded-in fix removed and `TestStableRefGuardRejectsMainSource` rejects. Fix: update L25 prose to `git push origin "<tagged-commit>:refs/heads/stable"` (or `$RELEASE_COMMIT`), describing the tagged-SHA advance, so the doc matches the code. One-line fix; mirror updates via the symlink.
  2. **Guard-hardening (record, captain may fold in — not an as-shipped hole, same strength as the established sibling).** The workflow guards are syntactic-presence checks, so three subtler FUTURE un-wirings stay GREEN: (a) `if: "!contains(github.ref, '-') && false"` skips the gate for ALL tags yet passes both guards (ifSkipsPreRelease is a substring match); (b) `continue-on-error: true` on the gate step makes a mismatch non-fatal yet passes; (c) reassigning `RELEASE_COMMIT="$(git rev-parse main)"` after `git switch main` re-introduces the stable/main divergence yet passes the stable-ref guard (it matches the variable NAME, not its value). NOTE: the pre-existing, already-merged `e2egate_workflow_test.go` that cycle-1 told the implementer to MIRROR has the SAME blind spots — no workflow guard in the repo checks `if`-truthiness/`continue-on-error`. So this is a class-wide pattern limitation, not a cycle-2 regression; hardening it (assert the gate step has no `continue-on-error` and an `if` that doesn't reduce to false on a stable tag) would lift both guards together.
  3. **MINOR (doc, fail-SAFE) — REL_SHA vs greened-SHA race.** Step 4's `REL_SHA=$(git rev-parse HEAD)` (worktree HEAD) and `gh workflow run "Runtime Live E2E" --ref main` green whatever main HEAD is server-side; a concurrent main push between step 3 and step 4 means the run greens a different SHA than REL_SHA. The e2e-gate (`run.HeadSha == releaseCommit`) BACKSTOPS this (the cut blocks, fail-safe), but the doc's "REL_SHA ... is the SHA this run must go green on" overstates the guarantee. Clarity fix, not a ship-safety hole.

### Summary
Cycle-2 closed the cycle-1 BLOCKER: the manifest-tag-gate is now genuinely CI-enforced (gate step in e2e-gate, goreleaser needs it, fatal under set -e), proven by reproducing the workflow guards red-then-green against the REAL on-disk release.yml, and the folded-in stable-advance binds `stable` to the tagged SHA. The binary blocks the v0.20.0 inversion and passes a match; full suite green. REJECTED on one concrete in-deliverable defect: docs/releasing.md:25 still describes `git push origin main:refs/heads/stable`, the exact form the same commit's code change removed and the new adversarial guard rejects — a stale-prose contradiction (AC-3 / captain item 4). Bundled for the same one-line rework: optional guard-hardening so the workflow guards catch `if: && false` / `continue-on-error` / RELEASE_COMMIT-reassignment un-wirings (the established sibling guard shares these blind spots, so this is class-wide hardening, captain's call).

### Feedback Cycles (cycle 3 — routing)

- **Cycle 3 — captain-authorized route (2026-06-21).** Cycle-2 validation REJECTED on one in-deliverable defect (AC-3): docs/releasing.md narrates `git push origin main:refs/heads/stable` while cycle-2's own code pushes `$RELEASE_COMMIT:refs/heads/stable` (the `main:` form `TestStableRefGuardRejectsMainSource` rejects). The release machinery itself is proven good across two detached audits. The 3-cycle escalation surfaced to the captain, who authorized one final route. SCOPE for this cycle:
  1. **Rebase the branch onto current `origin/main` and RESOLVE the `docs/releasing.md` conflict** by integrating fe's (#431) branch-per-channel marketplace edits with this task's reconciliation — #431 merged edits to the SAME doc after this branch forked, so the cut procedure and the marketplace/channel description must read as one coherent doc.
  2. **Apply the AC-3 fix:** docs/releasing.md describes `git push origin "$RELEASE_COMMIT:refs/heads/stable"` (the tagged-SHA advance), matching release.yml:249 and `TestStableRefGuardRejectsMainSource`. One-line substance; mirror updates via the symlink.
  OUT OF SCOPE (captain ruling — file as a follow-up, do NOT fold in): the class-wide workflow-guard hardening (syntactic-presence checks miss `if: && false` / `continue-on-error` / RELEASE_COMMIT reassignment); the MINOR REL_SHA-vs-greened-SHA doc clarity note.

## Stage Report: implementation (cycle 3)

- DONE: Rebase this branch onto current origin/main and RESOLVE the docs/releasing.md conflict, integrating fe/#431's branch-per-channel marketplace edits WITH this task's cut-procedure reconciliation into one coherent doc.
  Rebased onto origin/main 7338d556 (#431 now the merge-base); both task commits replayed (d3ff3617, f864bc1f). The conflict was the single "What the Tag Push Does" marketplace paragraph. Resolution KEEPS #431's branch-per-channel structure (repo root = marketplace `spacedock`; `edge` branch = `spacedock-edge`; channel lives in the marketplace NAME; binary adds the channel's branch source) AND this task's moving-`stable` mechanism, DROPPING #431's now-false "stable entry repointed to the released tag / commit on the repo root" claim — proven false because release.yml auto-advances `stable`. No conflict markers; release.yml + manifest-tag-gate code re-applied cleanly (only docs/releasing.md conflicted, confirmed via name-overlap check).
- DONE: Fix the AC-3 BLOCKER — docs/releasing.md describes `git push origin "$RELEASE_COMMIT:refs/heads/stable"` (the tagged-SHA advance), matching release.yml:249 and TestStableRefGuardRejectsMainSource; the mirror updates via the symlink.
  Wrote the AC-3 fix directly into the conflict resolution (docs/releasing.md:31-32): "release.yml resolves the tagged commit as `RELEASE_COMMIT` and runs `git push origin \"$RELEASE_COMMIT:refs/heads/stable\"`". `grep main:refs/heads/stable docs/releasing.md` → NONE. Matches release.yml:249 exactly. TestStableRefGuardRejectsMainSource + TestStampStepAdvancesStableRefToTaggedCommit both PASS by name.
- DONE: Self-check before handoff — go test ./... green; both releasing.md views diff-empty (symlink); fence count even; manifest-tag-gate + stable-ref guard tests still green.
  `go test ./...` green across all 16 packages, 0 failures (internal/cli 49s exercises #431's channel-marketplace-source; internal/release 8.9s). `diff docs/releasing.md docs/site/contributing/releasing.md` empty (symlink, mode 120000). Fence count even (18). manifest-tag-gate CI wiring intact: gate step release.yml:138, goreleaser `needs: e2e-gate` L141, workflow guards TestReleaseWorkflowGatesGoreleaserOnManifestTag + ...GuardRejects* all PASS.

### Summary
Closed the cycle-2 AC-3 BLOCKER inside the captain-authorized cycle-3 route. Rebased onto origin/main and resolved the lone docs/releasing.md conflict so #431's branch-per-channel marketplace model and this task's e2e-gate cut-procedure reconciliation read as one coherent doc — keeping both correct contents and dropping only #431's stale "repointed by a hand-edit commit" framing (the live mechanism is release.yml's moving-`stable` advance). The AC-3 fix landed in the same resolution: the doc now narrates `git push origin "$RELEASE_COMMIT:refs/heads/stable"`, matching release.yml:249 and the adversarial TestStableRefGuardRejectsMainSource. The two OUT-OF-SCOPE items (class-wide workflow-guard hardening; REL_SHA-vs-greened-SHA clarity note) were left for follow-up per the captain ruling. Full suite green; mirror byte-identical by symlink construction.
