---
id: 7yd3mbsy2am5qggc17sxvz2v
title: Reconcile docs/releasing.md to the live release machinery (one consolidated pass)
status: ideation
sprint: 0230-stable-finalization
source: "0202 Commander drive (2026-06-13). releasing.md 'Cutting a Stable Release' step 3 says to stamp+commit the version before tagging, but v0.20.1 and v0.20.2 both tagged the gated commit directly and let release.yml stamp post-tag (a pre-stamp creates an ungated commit the exact-SHA e2e-gate blocks on)."
group: cleanup
started: 2026-06-21T06:05:14Z
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
