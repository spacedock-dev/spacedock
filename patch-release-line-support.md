---
id: d0g21c517b5nvga1ybwckapk
title: Patch-release line support - gate release.yml's main-stamp on line-latestness, fix the stable-branch advance, automate the preversion bump
status: implementation
source: "Captain CL, 2026-08-25, reconciling gr and tw after the v0.27.0 cut: 'reconcile tw and gr and recommend best approach' - supersedes next-independent-release-line (twq68r4y8qg0wetztajtmmzz), whose body described the retired next-branch model; the live incidents of 2026-08-25 are the spec"
started: 2026-08-25T16:33:56Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-patch-release-line-support
issue:
related: "next-post-release-preversion-bump (closed delivered); stamp-then-tag-release-ritual"
gates:
    version: 1
    records:
        - id: gate:d0g21c517b5nvga1ybwckapk:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:d0g21c517b5nvga1ybwckapk-backlog-1
              briefing:
                id: briefing:d0g21c517b5nvga1ybwckapk:backlog:attempt-1:revision-1
                digest: sha256:6217483748c63c0663f4d6a6b96d8b90277d19e216e08534a10f5fca95168c2f
                request-digest: sha256:7c1dce1d1eb8b5e2720ce244ea77f986847eec757358eed954c21f83a61208e1
                room-ref: ./patch-release-line-support/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:d0g21c517b5nvga1ybwckapk:backlog:1
                briefing: briefing:d0g21c517b5nvga1ybwckapk:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-25T16:33:44.49838Z"
                decision: approve
                reason: 'Captain chat 2026-08-25: ''dispatch d0'' — approved seeding into design; the v0.27.0 cut''s live incidents are the spec'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:d0g21c517b5nvga1ybwckapk:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:d0g21c517b5nvga1ybwckapk-ideation-1
              briefing:
                id: briefing:d0g21c517b5nvga1ybwckapk:ideation:attempt-1:revision-1
                digest: sha256:7f5a4070fbf6a244571236e7d1bb75cc9a9a2059047f2728cd0d91e7b63d6dba
                request-digest: sha256:bd53f3a5a3b4a3885bfcac2ba874365e74fb78c27e97c23d651cf42950eeb634
                room-ref: ./patch-release-line-support/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:d0g21c517b5nvga1ybwckapk:ideation:1
                briefing: briefing:d0g21c517b5nvga1ybwckapk:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-25T21:50:55.243501Z"
                decision: revise
                reason: 'Captain chat 2026-08-25: ''leaner'' — accepts the lean-cut direction and takes the body''s full -74 menu: drop the goreleaser premise-guard test (-14; the live ideation verification stands as the premise evidence), drop the cmd exit-contract test (-34; the replay exercises the subcommand through the real step), and drop mechanism C the bounded rebase retry with its test (-26; a rare concurrent-merge race becomes a red job with a one-click idempotent re-run). Landing target near +280. AC-1''s proof is untouchable.'
            - id: gate-attempt:d0g21c517b5nvga1ybwckapk-ideation-2
              briefing:
                id: briefing:d0g21c517b5nvga1ybwckapk:ideation:attempt-2:revision-1
                digest: sha256:89b23eb78b5eb31796828e6034d05d27c688975edf40f3955c32b458a755ac17
                request-digest: sha256:c87923fce8ffe087cadb5ee9fee43018caa633890c3b8c1fc71976c3c9684f94
                room-ref: ./patch-release-line-support/review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:d0g21c517b5nvga1ybwckapk:ideation:2
                briefing: briefing:d0g21c517b5nvga1ybwckapk:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-25T22:30:43.46659Z"
                decision: approve
                reason: 'Captain chat 2026-08-25: ''approve d0 onto 760 stack'' — accepts the trimmed lean cut at the +281 baseline; implementation delivers as a stacked layer on PR #760''s branch'
              application:
                target-stage: implementation
                state: consumed
---

**LEAN CUT (captain ruling, 2026-08-25).** Stop the release-line cut that silently regresses the
binary channels, and stop hand-repairing the edge line after every stable cut. Patch-line DELIVERY —
moving `stable` across release lines — is deliberately NOT in this cut and is deferred whole to a
filed follow-up. See "Deferred: patch-line delivery" below for what was dropped and why.

## Problem

`release.yml` has one conditional for stable tags — `!contains(github.ref, '-')` — and it decides
several unrelated questions at once. On the only shape the repo has ever cut (a stable release from
`main`'s tip) every answer coincides, so the collapse has never been visible. A tag on an older line
answers them differently.

**1. An old-line stable tag publishes a binary-channel regression with every job green.** Nothing in
the pipeline compares the tag being cut against the release the stable channel already serves. On a
`v0.27.1` tag cut while `stable` serves 0.28.0, goreleaser runs to completion: `release.prerelease:
auto` (`.goreleaser.yaml:105`) marks any hyphen-free tag a normal release, so GitHub flips
`/releases/latest` DOWN to 0.27.1; and `homebrew_casks[spacedock].skip_upload: auto`
(`.goreleaser.yaml:125`) skips the cask bump ONLY for prereleases, so the stable Homebrew cask is
bumped DOWN too. Both are binary channels, both are consumed by `brew upgrade` and by the install
script, and neither is recoverable by a job re-run — the release is published and the cross-repo tap
commit is pushed. Every job is green. This is the sharpest defect in the surface and the centerpiece
of this cut.

**2. The main stamp has no line-awareness.** "Stamp plugin manifests to the release version"
(`.github/workflows/release.yml:247`) fires on any hyphen-free tag. On a `v0.27.1` tag it switches to
`main` and stamps `main` DOWN to 0.27.1 — rewriting the 0.28.0-pre0 manifests and the FO prose pin
(`These skills require binary minor 0.28`), so every edge user's plugin claims 0.27 while their
binary is 0.28.0-pre0 and the FO binary gate aborts. This is worse than the v0.27.0 incident: it
breaks a currently-working state rather than failing to repair a broken one. The latest-line decision
already exists in the sibling `edge-advance` job, and this step does not consult it. Defect 1's gate
does NOT cover this case: a `v0.27.1` cut while `stable` serves 0.27.0 is not a regression of the
stable channel, so it passes the gate — and stamps `main` DOWN anyway.

**3. Nothing stamps `main` past the released minor.** The auto-pre0 job tags `vX.(Y+1).0-pre0` and
publishes that edge binary, but leaves `main`'s manifests and FO pin at the released version. Every
edge install then aborts at the FO version gate until a human commits the bump. Observed live at the
v0.27.0 cut: the pre0 tag landed 2026-08-24 21:55:31 PDT, the failing install hit 17f5cd591 at
22:01, and the hand repair b8346ffc9 landed 22:01:49 — a 6m18s outage closed by one human commit.
`docs/releasing.md` step 9 (b04b3effd) is the interim procedure; this item retires it.

Defects 1 and 2 are latent — verified by reading and exercising the shipped workflow, never yet
triggered live. Defect 3 is the observed incident.

## Proposed approach

Three mechanisms, in descending order of harm removed. Each is either a new step under an existing
job's existing setup, or a move of a shipped step under a shipped gate.

### A. Stop the regressing cut before it publishes — a new gate step in `e2e-gate`

A bare-tag step in the `e2e-gate` job reads the version the stable channel serves NOW and fails the
run when the tag is older than it. goreleaser `needs: e2e-gate`, so a failure here means goreleaser
never starts: no GitHub Release, no `/releases/latest` flip, no cask commit.

```bash
set -euo pipefail
LS=0
git ls-remote --exit-code origin refs/heads/stable >/dev/null || LS=$?
if [ "$LS" -eq 2 ]; then
  echo "::notice::stable-regression gate: no stable ref yet, so there is no release to regress"
  exit 0
elif [ "$LS" -ne 0 ]; then
  echo "::error::stable-regression gate: cannot read refs/heads/stable (git ls-remote exit $LS); failing closed" >&2
  exit 1
fi
git fetch origin stable
STABLE_MANIFEST="$(mktemp)"
git show FETCH_HEAD:.claude-plugin/plugin.json > "$STABLE_MANIFEST"
go run ./cmd/spacedock-release stable-regression-gate "$GITHUB_REF_NAME" "$STABLE_MANIFEST"
```

- **Value AC served:** AC-1 (an old-line tag never publishes).
- **Placement is free.** The `e2e-gate` job already checks out with `fetch-depth: 0`, already sets up
  Go, and already carries two block-the-cut steps (`e2e-gate`, `manifest-tag-gate`) with exactly this
  exit contract. This step is the third of the same family. Placing it in `goreleaser` instead would
  start a job that must then abort; placing it in a new job would add a node to the critical path.
- **The predicate is `>=`, not `>`.** The gate blocks only a STRICTLY older tag. Equality must pass:
  re-running a release whose commit already reached `stable` is a supported, idempotent recovery (the
  auto-pre0 step is written for exactly that re-run). A strict-`>` boundary would turn every re-run
  into a red job. This is the one place this gate's boundary deliberately differs from
  `EdgeAdvanceDecision`'s strict `>`, and the unit table pins it.
- **New mechanism — `stable-regression-gate`.** A new `spacedock-release` subcommand in the shipped
  gate family: read a manifest, compare, exit 0 or 1. It reuses `release.ManifestVersion` and
  `release.ComparePreVersion`, both shipped and both already tested.
  - *Simplest alternative:* no new Go at all — feed the tag version and the stable version to the
    shipped `highest-known-edge-version` and fail when its answer is not the tag. Arithmetically this
    is the same `>=` predicate. *Why insufficient:* `HighestKnownEdgeVersion` SKIPS an unparseable
    candidate rather than erroring (`edge_advance_decision.go:187`), so a corrupt or renamed stable
    manifest would drop out of the comparison, the tag would win by default, and the gate would pass
    the cut it exists to block. A gate must fail loudly on input it cannot read. It also still needs
    a Go call to read the version out of the JSON, so the "no new Go" saving is not real.
  - *Second alternative:* compare the tag against the highest bare stable TAG in the repo, which
    needs no fetch at all. *Why insufficient:* a tag records what was CUT, and the channels this gate
    protects record what was PUBLISHED. A tag whose run died before goreleaser would then block the
    re-cut that repairs it, and a published release whose tag was later deleted would stop being
    counted. `stable` moves only after goreleaser succeeds, so it is the true record of the channel.
- **Failure direction is asymmetric, so the carve-outs are asymmetric.** A missing `stable` ref
  (`ls-remote` exit 2) is the first-stable-release case and must PASS — there is no release to
  regress. Any other non-zero exit is a read failure and must FAIL — an unreadable baseline is not a
  permission to publish, and the cost of a false block is one job re-run against a published
  regression that cannot be re-run away.

### B. Main stamp — move it into the job that already owns the decision

Move the stamp step out of `goreleaser` and into `edge-advance`, placed AFTER the auto-pre0 step and
carrying the SAME gate the pre0 step carries (`steps.decision.outputs.advance == 'true'`). Change
what it stamps: `edge-pre0-version "$RELEASE_VERSION"` (X.(Y+1).0-pre0), not `$RELEASE_VERSION`.
Rename it to `Stamp main to the next edge prerelease version`, because it no longer stamps the
release version. The remainder of the old step — the tagged-commit resolve and
`git push origin "$RELEASE_COMMIT:refs/heads/stable"` — stays in `goreleaser`, renamed to
`Advance the stable channel ref to the tagged commit` and otherwise BYTE-UNCHANGED.

- **Value AC served:** AC-2 (a patch tag never moves `main` backwards) and AC-3 (the edge gate passes
  with no human step).
- **Simplest alternative considered:** leave the step where it is and gate it on monotonicity —
  stamp only when the target version exceeds `main`'s current manifest version. No cross-job move, no
  decision reuse.
- **Why insufficient:** monotonicity gets the ARITHMETIC right and the COUPLING wrong. It would stamp
  `main` to X.(Y+1).0-pre0 whenever that is higher, including on runs where the auto-pre0 cut skipped
  and no X.(Y+1) binary was ever built — pointing every edge user's plugin at a version with no
  binary, which is the exact breakage AC-3 exists to remove. The stamp is only correct when the
  matching binary exists, and the only thing that knows whether it exists is the decision that cuts
  it. Co-locating them in one job under one gate makes disagreement structurally impossible.
- **Alternative also considered and rejected:** promote the decision to a standalone job with
  outputs, consumed by both. Same result, but it adds a job to `goreleaser`'s critical path and about
  25 lines of wiring for no behavior the move does not already give.
- **Why mechanism A does not make this redundant.** A's gate blocks a tag older than what `stable`
  serves. It does NOT block a `v0.27.1` cut while `stable` serves 0.27.0 — correctly, that patch is
  the newest stable release at that moment. That is precisely the tag today's stamp step would use to
  rewrite `main` DOWN to 0.27.1. The two mechanisms cover disjoint cases.
- **Ordering inside the job matters.** The pre0 tag is cut and its run is verified FIRST; the stamp
  runs only after that verify poll passes. A failure before the stamp therefore leaves exactly
  today's recoverable state, never a worse one.

### Not included: a retry on the main-stamp push (cut at briefing-1)

After mechanism B the stamp commits on every latest-line cut, so a concurrent merge to `main` between
fetch and push becomes a realistic non-fast-forward rejection. An earlier draft answered that with one
bounded `git pull --rebase origin main` and a re-push. The captain cut it at briefing-1 and this
records the accepted consequence:

**A concurrent-merge race lands as a red `edge-advance` job.** The recovery is a one-click, idempotent
job re-run — the pre0 tag is never re-minted (the step checks `refs/tags/` before tagging) and the
stamp itself is idempotent. It is NOT the hand-authored commit that AC-3's 6m18s baseline measured, so
the value AC stands without the retry. The race is rare, the failure is loud, and the recovery needs
no authoring. `docs/releasing.md` records the re-run instruction (Documentation diff block 5).

### C. Retire the manual ritual

`docs/releasing.md` step 9 is deleted and step 10 renumbered. The stamp and stable-channel prose is
corrected to the new job layout, and a short "The Stable Regression Gate" section is added. The
concrete diff is in **Documentation diff** below.

- **Value AC served:** AC-4, which counts only paired with AC-3.

## Deferred: patch-line delivery

**Dropped from this cut, whole:** the stable-advance rework — a version-gated `--force-with-lease`
push, its decision wiring on the push path, and the new force authority over `refs/heads/stable`. The
`goreleaser` job keeps today's plain `git push origin "$RELEASE_COMMIT:refs/heads/stable"`, untouched
except for the step rename.

**What this cut therefore does NOT fix, stated precisely.** Mechanism A blocks the REGRESSING class
of old-line cut. It does not block the non-regressing one: a `v0.27.1` cut while `stable` serves
0.27.0 passes the gate, publishes correctly, and advances `stable` to a commit that is NOT on `main`'s
history. The NEXT latest-line cut then fails its stable push non-fast-forward (spike TEST 2), after
goreleaser has published, and the stable channel freezes on the patch line. That failure is LOUD and
post-publication; the failure mechanism A removes is SILENT and post-publication. The lean cut takes
the silent one and leaves the loud one to the follow-up. Until the follow-up ships, the operating
constraint is: do not cut a patch line.

**Follow-up:** `patch-line-stable-delivery-merge-commit`. Direction, from the staff review: advance
`stable` with a MERGE COMMIT — `git commit-tree` on the tagged commit's tree with `stable`'s current
tip as the first parent and the tagged commit as the second. The advance is then always a
fast-forward, so no force and no lease are needed, and a concurrent move of `stable` fails closed on
a plain non-fast-forward rejection. Its cost, which is why it is its own task and not a line here:
`stable`'s tip SHA no longer equals the tag's SHA, so every consumer that assumes the identity —
`stableRefPushSource` in `channel_agreement_guard_test.go`, and any install-path check comparing the
resolved ref to a tag — must be revisited.

**Verified fact recorded for the follow-up's gate.** `refs/heads/stable` has NO branch protection and
this repo has NO rulesets, checked directly on 2026-08-25:
`gh api repos/spacedock-dev/spacedock/branches/stable/protection` returns 404 `Branch not protected`,
and `gh api repos/spacedock-dev/spacedock/rulesets` returns `[]`. A `--force-with-lease` push from CI
would therefore have been NORMALIZED behavior under the credentials the release run already holds,
not a new capability grant. That weakens — it does not remove — the "new authority" objection the
earlier draft raised against itself. The follow-up should weigh the merge-commit approach on its
fail-closed race behavior, not on a protection boundary that does not exist.

## Risk evidence

**Spike 1 (this cycle) — the read path mechanism A depends on, and the fixture shape AC-1 needs.**
The gate's discrimination rests on three unexercised claims: that `git ls-remote --exit-code`
distinguishes "no such ref" from "cannot read the remote", that `git show FETCH_HEAD:<path>` reads
the manifest the stable channel serves, and — the one the review named as the likeliest hole in the
earlier design's own tests — that a version regression can reach `stable` WITHOUT git's ancestry
check refusing it, so a missing gate cannot be masked by a git refusal. Run in a throwaway
bare-origin fixture (bash, not zsh):

| Test | Setup | Result |
|---|---|---|
| A | `ls-remote --exit-code origin refs/heads/stable`, ref absent | **exit 2** — the first-release carve-out is distinguishable |
| B | same against an unreachable remote | **exit 128** — so "fail closed on any non-2" is a real discrimination, not theatre |
| C | same, ref present | exit 0 |
| D | `fetch origin stable` + `git show FETCH_HEAD:.claude-plugin/plugin.json` | prints `{"version":"0.27.0"}` — the baseline read works |
| E | stable at the 0.28.0 commit; push a **child** commit stamped 0.27.1 to `stable` | **fast-forward SUCCESS**, `stable` then serves 0.27.1 — a version regression through a clean fast-forward |
| F | `LS=0; git ls-remote … \|\| LS=$?` under `set -euo pipefail` | captures the code, shell survives — the earlier draft's unbound-variable finding cannot recur in this shape |

TEST E is the load-bearing one. It builds exactly the state the review demanded the fixture prove:
**version inversion WITHOUT ancestry inversion.** Git accepts the push, so the unguarded path really
would publish and really would move the channel down. It is also the realistic mis-cut, not a
contrived one: `docs/releasing.md` step 2 tells the cutter to branch off `origin/main`, so a cutter
hand-stamping 0.27.1 on a branch off today's `main` produces precisely commit shape E.

**Spike 2 (prior cycle, still standing) — the stable-ref push semantics.** A 7-test bare-origin
fixture established that a patch commit that is a CHILD of `stable` fast-forwards cleanly (TEST 1,
refuting the seed's premise that the patch push dies), that a later `main` commit is then rejected
non-fast-forward (TEST 2), and that fetching an absent `stable` exits 128 (TEST 7). TESTs 1, 2 and 7
still carry weight: TEST 2 is the evidence for the deferred follow-up's frozen-channel failure, and
TEST 7 corroborates spike 1's carve-out. TESTs 3 and 4 exercised `--force-with-lease` and are now
moot with that mechanism dropped; they are retained in the follow-up's record, not here.

**Spike 3 (prior cycle, still standing) — the decision helper already ranks the patch case** against
this repo's real 71-tag pool, so mechanism B needs no change to `edge-advance-decision`:

```
v0.27.1   known=0.28.0-pre1  -> skip      (correct: main is already the 0.28 line)
v0.28.0   known=0.28.0-pre1  -> advance   (correct: the next latest-line cut)
v0.26.1   known=0.28.0-pre1  -> skip      (correct: two lines back)
```

**The harness is not new.** `internal/release/edge_advance_decision_shell_test.go` already extracts a
real `release.yml` step's `run` block and executes it against a fixture git repo via
`GIT_DIR`/`GIT_WORK_TREE` with a stubbed `$GITHUB_OUTPUT`. This task extends that proven harness with
a bare origin (spike 1's shape) rather than inventing one.

**Declared harness limitation — MOOT as of the Cycle 1 proof-posture ruling.** No harness ships; the
paragraph below is retained as the record of what was designed, not as a description of the
deliverable. Validation should expect no replay harness and no binary substitution in the diff.

The replayed steps invoke `go run ./cmd/spacedock-release`, which
resolves relative to the process working directory. Running with the working directory at the repo
root would let a replayed stamp write this checkout's real manifests. The harness therefore runs with
the working directory in the fixture and substitutes that one token for a pre-built binary path.
Every other byte of the script — the gating, the refspecs, the conditionals — runs verbatim. This
substitution is declared here so validation checks it rather than discovering it.

## Out of scope

The stable-advance rework and all patch-line delivery (see "Deferred" above). The e2e-gate/waiver
mechanics themselves. The marketplace display-version fields. The release notes ritual. The release
ritual's own patch-line steps (step 2 branches off `origin/main`; a patch would branch off the prior
stable tag instead) — operator procedure, not machinery, and it belongs with the delivery follow-up
that makes a patch line cuttable at all.

## Expected surface and tolerance

Estimate net LOC change: **+281, across 9 files** (insertions ~+343, deletions ~-62).
Tolerance: ±30% on the net figure (+197 to +365) and ±2 files.

| File | ins | del | Change |
|---|---|---|---|
| `.github/workflows/release.yml` | 64 | 36 | new gate step in `e2e-gate`; stamp step moved to `edge-advance`; goreleaser step renamed |
| `internal/release/stable_regression_gate.go` | 32 | 0 | new — the decision over `ComparePreVersion` |
| `cmd/spacedock-release/stable_regression_gate.go` | 34 | 0 | new — thin wrapper over `ManifestVersion` + the decision |
| `cmd/spacedock-release/main.go` | 6 | 0 | dispatch case, usage line, doc comment |
| `internal/release/stable_regression_gate_test.go` | 58 | 0 | new — unit table, including the `>=` boundary |
| `internal/release/stable_regression_shell_test.go` | 97 | 0 | new — bare-origin replay of the real step, plus the unguarded baseline |
| `internal/release/edge_advance_wiring_test.go` | 18 | 0 | structural guard that the moved stamp shares the pre0 gate |
| `internal/release/channel_agreement_guard_test.go` | 14 | 8 | two parsers follow the renamed/split steps |
| `docs/releasing.md` | 20 | 18 | step 9 retired; three passages corrected; one short section added |

**Estimate history, so the gate can audit the number rather than trust it.** Cycle 2 priced this at
+355 across 11 files and declined to massage it to fit the ruling's expectation, offering instead a
priced menu of what could go. The captain took the whole menu at briefing-1, and this is that menu
applied, line for line: the goreleaser premise guard (−14 insertions, one file), the command
exit-contract test (−34, one file), and the stamp-push retry with its test (−26, split 8 in
`release.yml` and 18 in the replay). 417 − 74 = 343 insertions, deletions unchanged at 62, so
**net +281 across 9 files**. Each cut's accepted consequence is recorded where the cut lands, not
only here.

**What is now the floor.** AC-1's proof is untouchable by the captain's own reason, and it dominates
what remains: the bare-origin replay (97 lines) exists to show that the UNGUARDED path really would
publish, which is the review's sharpest requirement and the thing this whole reshape exists to add.
The lean cut's saving was never mostly in lines — it is in authority and risk: no force-push
capability, no new push-path wiring, no ancestry semantics in the release path, and one fewer
irreversible failure mode.

**Observable semantics this task changes:**

1. **A new pre-publication block.** A stable tag older than the release `stable` serves fails the run
   before goreleaser starts. Today it publishes and moves both binary channels down.
2. **What `main` carries after a latest-line cut.** The plugin manifests and the FO prose pin move to
   `X.(Y+1).0-pre0` instead of `X.Y.Z`. This is a change to on-disk content the edge marketplace serves.
3. **New command grammar:** `spacedock-release stable-regression-gate <tag> <manifest>`.
4. **Failure attribution moves.** A main-stamp failure will red `edge-advance`, not `goreleaser`.
5. **A documented procedure is removed:** `docs/releasing.md` step 9.

**Explicitly NOT changed — this is the point of the lean cut.** CI's push authority over
`refs/heads/stable` stays a plain fast-forward push. No force, no lease, no new capability. Also
unchanged: the e2e-gate's own logic, the manifest-tag-gate, goreleaser itself, the journey ledger, the
pre0 tag mechanics, the marketplace repo, and `edge-advance-decision` (reused verbatim).

## Path-to-lane call

The diff touches `.github/workflows/release.yml`, `cmd/spacedock-release/**`, `internal/release/**`,
and `docs/releasing.md`. None of these is a file a live lane loads or drives — not
`skills/**/references/**`, not the dispatch/launch path, not a live lane's own tests — so the
deterministic lanes (build/install/offline) are sufficient and no live lane is required green.

The **detached adversarial audit IS required**: `.github/**` release machinery is one of the Proof
policy's four high-stakes surfaces. Run it on a throwaway checkout before merge. The audit's sharpest
targets are named in the test plan.

## Acceptance criteria

**Proof posture (captain ruling, 2026-08-25; README "Release-machinery proof posture").** A release
failure is loud, so the next real cut is the live test for in-situ shell behavior. Decision logic is
proven by Go unit table, YAML wiring by structural check, and nothing replays workflow shell against
a fixture repository. Each AC below states which of the three carries it.

**AC-1 (value) — A stable tag older than the release `stable` points at never reaches publication:
the run fails in `e2e-gate`, before goreleaser starts.**
Verified by: `TestEvaluateStableRegressionGate` in `internal/release/stable_regression_gate_test.go`
— the decision table over every input class, including the older-tag block rows, the `>=` boundary
that lets a re-run pass, and the error rows that keep the gate loud on input it cannot read. This is
the silent class the posture bullet names: a wrong decision that publishes. Paired with
`TestReleaseStepsSitInTheirOwningJobs` in `stable_regression_wiring_test.go`, which binds the gate
step to `e2e-gate` — the job goreleaser needs, and therefore the placement that gives the gate its
blocking power. Falsifying change: flip the comparison to `<= 0` (the equality row reds — verified
live) or rename/move the gate step out of `e2e-gate` (the presence guard reds — verified live
against a mutated `release.yml`). In-situ: the gate's shell arms (the `ls-remote` exit-2 carve-out
and the fail-closed arm) are observed at the next real cut, per the posture.

**AC-2 (value) — A patch tag that PASSES the regression gate still leaves `main` untouched: the
`main` stamp cannot run on a tag whose latest-line decision is `advance=false`.**
Verified by: `TestMainStampSharesThePre0DecisionGate` in `edge_advance_wiring_test.go` — the stamp
step's `if:` must be EQUAL to the auto-pre0 step's `if:`, not merely carry some gate, so the stamp
and the binary that backs it cannot disagree. Its adversarial twin widens the stamp's `if:` back to
the stable-path-only form the retired goreleaser step carried, and the guard reds. Paired with
`TestReleaseStepsSitInTheirOwningJobs`, which binds the stamp to `edge-advance` — the job that owns
the decision. This AC is wiring, so the structural check is the whole proof: Actions itself
evaluates the `if:`, and per the posture that evaluation is observed at the next real cut.
Falsifying change: remove the decision gate from the stamp's `if:`, or diverge it from the pre0
step's — the guard reds (verified live).

**AC-3 (value) — A latest-line stable cut needs zero human commits to restore the edge line: the
`main` stamp writes the `edge-pre0-version` output, which is the version the same job's auto-cut
tags.**
Verified by: `TestMainStampSharesThePre0DecisionGate` for the co-location and shared gate, plus the
shipped `Pre0EdgeVersion` unit coverage for the version arithmetic — the stamp and the tag call the
SAME command on the SAME input, so agreement is by construction rather than by assertion.
Falsifying change: stamp `$RELEASE_VERSION` instead of the `edge-pre0-version` output; `main` then
lands on the released version and the manifest/tag pair disagrees. Recorded baseline, measured
before this change: on the fixture the v0.27.0 replay left `main` at 0.27.0 while the auto-cut pre0
tag was v0.28.0-pre0 — the exact 6m18s / one-hand-commit mismatch at the v0.27.0 cut (pre0 tag
21:55:31, hand repair b8346ffc9 at 22:01:49). Per the posture, the restored edge line is observed at
the next real cut.
**Scope of the "zero human commits" claim, after the briefing-1 cut of the retry:** it is a claim
about the SUCCESSFUL path, which is the path the 6m18s baseline measured. A concurrent merge to
`main` reds the job instead, and the recovery is a re-run, not a commit — so the AC's measured
quantity (human commits) stays zero either way.

**AC-4 — The manual ritual is gone, not merely superseded: `docs/releasing.md` contains no post-tag
`main` preversion bump step, and its stamp, stable-channel, and edge-advance prose describe the
shipped job layout.**
Verified by: the Documentation diff below applied, `docs/releasing.md` step count reduced by one with
step 10 renumbered to 9, and the changed prose passing the `simple-english` check the workflow README
requires. Falsifying change: leave step 9 in place — the doc then documents a hand bump the pipeline
performs, and a cutter runs a duplicate stamp. This AC is a mechanism AC and counts only paired with
AC-3, which measures the value it serves.

## Test plan

**Unit — `internal/release/stable_regression_gate_test.go` (low cost).** Table over the decision:
older tag vs newer stable (block); newer tag vs older stable (pass); EQUAL versions (pass — the
re-run case, and the boundary that differs from `EdgeAdvanceDecision`'s strict `>`); a prerelease tag
against a stable version (error, since the caller's `if:` guarantees a bare tag); an unparseable
manifest version (error, so a miswiring fails loud rather than silently passing the cut).

**No separate command exit-contract test (cut at briefing-1, −34).** The library table covers the
loud-failure behavior on input the gate cannot parse. Accepted consequence: the usage-error path
(exit 2 on a missing argument) is left unexercised, and it is reachable only by editing the step's
own argument list.

**No behavioral replay (cut by the captain's proof-posture ruling, 2026-08-25).** The bare-origin
replay harness is deleted, not deferred: the six replays, the two-repo fixture builders, and the
binary substitution are all gone. Authority: the README's **Release-machinery proof posture** bullet
— "Do not build replay harnesses that execute workflow shell against fixture repositories", because
a release failure is loud and the next real cut is the live test. Accepted consequence, stated
precisely: the gate's shell arms (the `ls-remote` exit-2 first-release carve-out and the fail-closed
arm on any other non-zero) have no standing test, and a change to either is caught at the next cut
rather than in CI. The silent class — a wrong decision that publishes — stays covered by the unit
table on the decision function, which is the posture bullet's own carve-out.

**Structural guards (the YAML-wiring half of the posture).** `edge_advance_wiring_test.go` carries
the AC-2 check that the stamp step's `if:` EQUALS the auto-pre0 step's, with the adversarial twin
(widen it → red). `stable_regression_wiring_test.go` carries the step-presence guard binding each of
the three added/moved steps to its owning job, with the adversarial twin (rename it away → red).
`channel_agreement_guard_test.go`'s `releaseStampTarget` and `stableRefPushSource` parsers both keyed
on `"Stamp plugin manifests to the release version"`, which this task splits and renames — without
the edit those tests fail, so this is mandatory, not optional. Both keep their adversarial twins.

**No standing guard on AC-1's harm premise (cut at briefing-1, −14).** The premise is recorded
evidence, verified live at ideation: `.goreleaser.yaml:105` `release.prerelease: auto` and
`.goreleaser.yaml:125` `skip_upload: auto` on the stable cask, both read at their shipped line
numbers. Accepted consequence: a future edit to either setting would remove the reason this gate
exists without any test reding. The gate itself stays correct under such an edit — a regressing tag
is still refused — so what is lost is the standing proof of the motive, not the protection.

**Full suite.** `go test ./...` and `go test ./... -race`, plus `gofmt -w ./cmd ./internal`.

**No live workflow test, and the shell behavior is observed at the next cut.** No live lane loads or
drives any changed file (see Path-to-lane call). Under the proof posture, in-situ behavior for all
three mechanisms — the gate's arms, the Actions evaluation of both `if:` conditions, and the stamp's
effect on `main` — is confirmed at the first real stable cut after merge. That observation is
deliberately NOT an AC, because validation cannot reproduce it on demand.

**Detached adversarial audit targets.** The `.github/**` four-surface trigger still applies; the
posture changed what proof exists, not whether the audit runs. Sharpest questions for the auditor,
restated for the surviving proof set: (1) Can the gate step be edited to compare against the TAG pool
instead of the `stable` ref with every test still green? (2) Can the stamp step be edited to stamp
`main` while the pre0 tag was never pushed, with every test still green? (3) Does the unit table
actually cover the silent class, or does a decision input class reach `release.yml` that no row
exercises? (4) Do the two structural guards red under a real mutation of `release.yml`, or only
against their in-test string copies? Targets that depended on the deleted replays — the
unguarded-push baseline and the `ls-remote` fail-open flip — are retired with them, and the
fail-closed arm is now an accepted uncovered path, recorded above.

### Feedback Cycles

- Cycle 1: REWORK — captain proof-posture ruling 2026-08-25 ("we don't need those tests; we would know if the release fails", recorded as the README's release-machinery proof posture): delete the bare-origin replay harness; keep the comparator unit table, the structural wiring guards, and the renamed-step consumer updates; production untouched; ACs re-verify by unit + structural + next-live-cut observation. Supersedes the FAILED +644 surface item; re-measure after the cut.
- Cycle 2: REWORKED — harness deleted whole (b8e186df5, two test files, zero production bytes); surface 10 files/+381 net vs estimate +281±84 (136%, +16 past the ceiling — residual is ordinary per-file under-pricing, harness was 330 of the 363 overage); cycle-1 suite-green claim corrected (grep swallowed the exit code; both suites re-run properly, exit 0); ensigncycle timeout-edge finding preserved, FO-authorized decline, filed separately.

## Sequencing decision for the captain

**#760 rides v0.28.0.** The fix for #760 ("claude install leaves the sibling edge plugin installed and
enabled") lands on `main` and ships in the next minor. No patch-line cut, no release branch, no
exercise of any patch machinery. This is not merely the recommendation — it is the only option the
lean cut leaves open, because patch-line DELIVERY is deferred with the follow-up. A v0.27.1 cut after
this task ships would publish correctly and then freeze the stable channel on the next latest-line
cut (see "Deferred" above).

**v0.27.1 defers with the delivery machinery.** The first genuinely urgent old-line patch is the
trigger to promote `patch-line-stable-delivery-merge-commit` out of the backlog, not to cut against
half the machinery.

**The one-way constraint, unchanged in force and narrowed in reason.** Do NOT cut a v0.27.1 on
today's machinery: the stamp step would switch to `main` and stamp it DOWN to 0.27.1, rewriting the
0.28.0-pre0 manifests and the FO pin to 0.27, breaking every edge user. After this task, that
specific breakage is gated away (AC-2) and the regressing cut is blocked outright (AC-1) — but the
frozen-channel failure remains until the follow-up ships, so the constraint stands either way.

## Documentation diff

`docs/releasing.md` is user-facing documentation, so every "After" block below follows ASD-STE100 per
the workflow README's Prose style section. The text was CHECKED mechanically, not declared compliant:
a script extracts all six "After" blocks from this section and flags Rule 6.3 (sentence over 25
words), Rule 6.6 (paragraph over six sentences), Rule 8.1 (semicolon), and Rule 3.6 (passive with a
stated agent). It reports 6 blocks, 0 violations; the longest sentence is 17 words. The checker was
FALSIFIED before that result was trusted: re-run against a mutant of block 2 carrying a semicolon, a
56-word sentence, a 13-sentence paragraph, and "is repointed by release.yml", it reports all four
violations and exits 1. Rule 3.6 is why block 2 now opens "A hand-edit … does not repoint the stable
entry" instead of the first draft's "The stable entry is NOT repointed … by a hand-edit". The
dictionary rules a script cannot check were applied by hand: "serves" became "points at" — the
phrasing the shipped `docs/releasing.md` already uses, and one word with one meaning — "occurs"
became "happens", and "the stamp is idempotent" became "the stamp gives the same result each time".
New `release.yml` comments and Go doc comments follow the same rules. The moved stamp step's existing comment block converts as this task touches
it, per the README's convert-on-touch clause. This task body itself is workflow state, not
user-facing documentation, so it is out of the rule's scope.

**1. "What the Tag Push Does" — replace the stamp bullet (lines 25-26).**

Before:
```
- stamps the plugin manifests' `version` on `main`, then advances the stable
  channel ref (see below).
```
After:
```
- advances the stable channel ref to the tagged commit. See below.
- stops the cut before publication if the tag is older than the release that
  `stable` points at. See "The Stable Regression Gate" below.
- stamps the plugin manifests and the FO prose pin on `main`, but only on a
  latest-line stable tag. The stamp writes the `X.(Y+1).0-pre0` version that the
  same job auto-cuts. An old-line tag does not change `main`.
```

**2. Replace the paragraph beginning "The stable entry is NOT repointed per release" (lines 40-46).**
The current text cites the "Stamp plugin manifests" step by name, and this task renames it.

After:
```
A hand-edit in the marketplace repo does not repoint the stable entry per
release. `stable` is a MOVING BRANCH in this repo. After the tag fires,
release.yml resolves the tagged commit and pushes it to `refs/heads/stable`. A
fresh `spacedock@spacedock` install resolves whatever `stable` points at. That
push is what publishes the release to the stable channel. No marketplace-repo
commit is necessary.
```

**3. Replace the paragraph beginning "The post-tag manifest stamp is idempotent" (lines 48-51).**
The current text is false in two ways after this task: the stamp is no longer on the stable-tag path,
and it no longer finds "no diff".

After:
```
The `main` stamp is no longer part of the stable-tag path. The `edge-advance`
job stamps `main`, and only on a latest-line tag. The stamp writes the next
prerelease version. The tagged commit keeps the release version. The two values
are different by design, so the stamp always makes one commit on `main`.
```

**4. Delete step 9 entirely (lines 183-197), renumber step 10 to step 9.**

**5. "Advancing the Edge Line" — extend both bullets.**

Add to the latest-line bullet:
```
After the job verifies the pre0 run, the same job stamps the manifests and the
FO pin on `main` to the pre0 version. An edge install then passes the FO version
gate with no human step. If a merge to `main` happens during the run, the push
can fail. Run the job again, because the stamp gives the same result each time.
```
Add to the old-line bullet:
```
The same decision also gates the `main` stamp. An old-line tag cannot move
`main` to a lower version.
```

**6. New section "The Stable Regression Gate"**, after "Advancing the Edge Line".

After:
```
## The Stable Regression Gate

A stable tag that is older than the release that `stable` points at can damage
two binary channels. goreleaser marks each bare `vX.Y.Z` tag as a full release, so
GitHub moves `/releases/latest` to it. goreleaser also bumps the stable Homebrew
cask, because `skip_upload: auto` skips only a prerelease. An old tag moves both
channels DOWN, and every job stays green. The `e2e-gate` job stops this before
goreleaser starts.

- The gate reads the version in the `.claude-plugin/plugin.json` file that
  `stable` points at.
- If the tag is older than that version, the gate fails the run. goreleaser does
  not start, because it needs the `e2e-gate` job.
- If the tag is the same version or newer, the gate lets the run continue. A
  re-run of a release that already reached `stable` is still possible.
- On the first stable release the `stable` ref does not exist. The gate writes a
  notice and lets the run continue.
- If the gate cannot read the `stable` ref, it fails the run. A read failure is
  not a permission to publish.

The gate does not make a patch line deliverable. A patch that is newer than the
version `stable` points at still publishes, and `stable` then leaves the history
of `main`. The next latest-line release cannot advance `stable` after that. Do not
cut a patch line until that work is complete.
```

## Stage Report: ideation

- DONE: Design lands one owner for release.yml's stable-tag conditional with three concrete mechanisms: the line-latestness gate on the manifest-stamp step (state whether it reuses the edge-advance job's highest-known-version logic, or name the simpler alternative and why it is insufficient); a stable-branch advance that delivers a release-branch commit (the current exact-SHA push comment claims fast-forward, which is false off-line); and the automated post-tag main preversion bump replacing manual ritual step 9 — each mechanism names the value AC it serves.
  `## Proposed approach` A/B/C. The stamp gate REUSES the edge-advance decision by moving the step into that job under the same `if:`; the rejected simpler alternative (monotonicity against main's own manifest) is named with its concrete failure — it stamps main to a pre0 whose binary was never built. Each mechanism carries a "Value AC served" line and a "Why insufficient" line for its alternative.
- DONE: Riskiest unverified mechanism exercised first and recorded in the task body (for example: what the stable-ref push actually does with a release-branch commit against a frozen stable ref; what the auto-pre0 tag job stamps today), or an auditable `no spike needed: {proven mechanisms}` line.
  `## Risk evidence` — a 7-test bare-origin fixture spike REFUTED the seed's premise, plus a decision-helper spike over the real 71-tag pool; both re-runnable, script path recorded.
- DONE: Entity-level AC set with a value-measuring AC against a failing-today baseline; a sequencing section presenting v0.27.1-from-release-branch (cherry-pick the #760 fix) versus fold-into-v0.28.0 for the captain to rule at the gate; net-LOC estimate with tolerance; path-to-lane call (.github/** is release machinery — the detached adversarial audit's four-surface trigger applies); all comments and user-facing doc text follow ASD-STE100 per the workflow README's Prose style section.
  AC-1/2/3 each measure resulting on-disk state with the same-replay-on-current-release.yml as the failing baseline and a named falsifying change; AC-3's baseline is the measured 6m18s / one-hand-commit v0.27.0 outage. `## Sequencing decision for the captain`, `## Expected surface and tolerance` (net +350 ±30%, 9 files ±2, six declared semantic changes), `## Path-to-lane call` (audit required, no live lane), `## Documentation diff` (six numbered before/after items).
  ASD-STE100: the six proposed doc-text blocks were CHECKED against the `simple-english` catalog, not declared compliant. The first draft broke five rules (8.1 semicolon, 6.3 twice at 26 and 32 words, 6.6 seven-sentence paragraph, 3.6 passive, GR-1 dropped "that"); all are corrected. The checker was falsified against the pre-fix text and reds on it, so a clean result is not vacuous.

### Summary

The spike inverted item 2 of the seed. A patch-line commit pushed to `stable` fast-forwards cleanly
today — the seed's expected non-fast-forward death does not happen. The real defect is the cut AFTER
a patch: once `stable` leaves `main`'s history, every later latest-line release is rejected and the
stable channel freezes permanently. That changed the fix from "refuse with instructions" to "version
gate plus `--force-with-lease`", and it made the new force authority the single item most needing the
captain's explicit approval at the gate.

The design's load-bearing decision is that "does the edge line move" and "does the stable channel
move" are different questions with different answers for the same tag, so they get two predicates.
The main stamp reuses the existing edge decision by MOVING into the job that owns it, which makes
stamp/binary disagreement structurally impossible rather than merely tested — that is the necessity
argument, and it also means no new decision plumbing for two of the three items.

Two things the gate should weigh. First, the estimate is +350 net, not the seeded +60: roughly 60% is
the replay harness, and the seed counted machinery only. Second, #760 is still OPEN and its own fix
is only at ideation, so no v0.27.1 is cuttable today regardless of the ruling — and cutting one on
today's machinery would stamp `main` DOWN to 0.27.1 and break every edge user, so this task ships
before any patch-line cut either way.

One correction round after the completion signal, before any gate briefing was prepared, so the
body-freeze rule did not bind. The FO relayed today's ASD-STE100 ruling; my first draft had ASSERTED
compliance without checking, and the check found five real violations in the proposed
`docs/releasing.md` text. That is the same failure this workflow's proof policy names — treating a
claim as proven because I wrote it down. Corrected, and the check itself was falsified so it can
fail.

## Stage Report: ideation (cycle 2)

Reshape to the captain's LEAN CUT ruling (2026-08-25, "lean") after an independent staff review
returned needs-rework. The cycle-1 report above is left intact as the record of what was reshaped.

- DONE: The reshaped design carries exactly the lean-cut scope from the scope notes: the pre-publication old-line hard-fail in e2e-gate (with a fixture proving the unguarded path WOULD publish), the main-stamp move into edge-advance stamping pre0, the retired manual step 9, and the corrected-rationale rebase retry — with the stable-advance lease/force rework dropped and its deferral pointer plus the review's verified facts recorded in the body.
  `## Proposed approach` is now A (the `e2e-gate` hard-fail, centerpiece), B (the stamp move), C (the retry, with the corrected "avoids a human job re-run, not an outage" rationale), D (the retired ritual). `## Deferred: patch-line delivery` carries the drop, the `patch-line-stable-delivery-merge-commit` pointer with the merge-commit direction and its SHA-identity cost, and the branch-protection fact. The lease/force mechanism and its two review findings (fail-open fetch, unbound `$STABLE_MANIFEST`) are gone with it, not repaired.
- DONE: ACs, net-LOC estimate with tolerance, test plan, and sequencing section updated to the lean scope; every kept mechanism still names its value AC; all comments and user-facing doc text strict ASD-STE100.
  AC-1 is now the pre-publication block, measured against `TestUnguardedOldLineTagWouldReachStable` on the same fixture; AC-2 is the patch-tag/main case with its structural half declared load-bearing; AC-3 keeps the 6m18s baseline and gains the concurrent-move test; AC-4 unchanged in kind. Estimate is +355 across 11 files. `## Documentation diff` is six blocks, longest sentence 17 words, no semicolons, passives rewritten active — checked by counting, not asserted.
- DONE: The reshape is appended to the Stage Report and state is committed before the completion message.
  This section, plus the path-scoped state commit cited in the completion message.

### Summary

The centerpiece changed. The review's material finding is that an old-line cut regresses the BINARY
channels with all jobs green, and I verified the premise directly in `.goreleaser.yaml`:
`release.prerelease: auto` (line 105) makes any hyphen-free tag a full release, and the stable cask's
`skip_upload: auto` (line 125) skips only prereleases. So `/releases/latest` and the Homebrew cask
both move DOWN, and neither is recoverable by a job re-run. The new gate step in `e2e-gate` blocks
that before goreleaser starts, at the cost of one shipped-family subcommand.

I spiked the gate's own unexercised read path rather than reusing the prior cycle's evidence: six
tests establishing that `ls-remote --exit-code` returns 2 for a missing ref and 128 for an unreadable
remote (so the fail-closed arm is a real discrimination), that `git show FETCH_HEAD:<path>` reads
stable's manifest, and — the load-bearing one — that a commit stamped 0.27.1 pushed as a CHILD of the
0.28.0 commit reaches `stable` as a clean FAST-FORWARD. That is the fixture shape the review demanded:
version inversion without ancestry inversion, so git's own refusal cannot mask a missing gate. I also
re-verified the branch-protection fact rather than relaying it: `stable` is unprotected and the repo
has zero rulesets.

Two things the gate should weigh, and one of them contradicts the ruling. First, the ruling expected
"well under half of +350" and the reshape lands at +355. I did not massage the number. The dropped
mechanism and the added one cost about the same, because both are dominated by the same bare-origin
replay harness, and AC-1's baseline test — the review's own sharpest requirement — is the largest
single line item and is not compressible without gutting the proof. The lean cut's real saving is
authority and risk: no force capability, no new push-path wiring, one fewer irreversible failure
mode. The estimate section names a −74 menu if the gate wants the number down. Second, the scope
note's rationale that "with old-line cuts blocked pre-publication, `stable` never leaves main's
history" is true only while no patch line is cut at all: a v0.27.1 cut while `stable` serves 0.27.0
passes the gate by construction and still ends in a frozen channel on the next latest-line cut. The
body states that precisely and keeps the do-not-cut-a-patch-line constraint until the follow-up ships.

## Stage Report: ideation (cycle 3)

Applied the captain's briefing-1 revise: the full −74 menu the cycle-2 body priced, taken whole.
Cycles 1 and 2 above are left intact as the record.

- DONE: Drop the goreleaser premise-guard test (−14).
  Row removed from the surface table; the test plan's "Structural guards" now carries **No standing guard on AC-1's harm premise** with the accepted consequence — a future edit to `prerelease: auto` or `skip_upload: auto` removes the gate's stated motive with no test reding, while the gate itself stays correct. The premise evidence remains the live ideation verification at `.goreleaser.yaml:105` and `:125`.
- DONE: Drop the cmd exit-contract test (−34).
  Row removed; the test plan records that the replay drives the real step and so exercises the 0-and-1 exits end to end, and that the accepted loss is the exit-2 usage path, reachable only by editing the step's own argument list.
- DONE: Drop mechanism C, the bounded rebase retry, with its test (−26).
  The `### C` section is replaced by **Not included: a retry on the main-stamp push (cut at briefing-1)**, recording the accepted consequence verbatim from the ruling: a rare concurrent-merge race lands as a red `edge-advance` job whose re-run is one-click and idempotent. Old `### D` renumbered to `### C`. Doc-diff block 5 drops the rebase sentence and gains the re-run instruction instead.
- DONE: Estimate, file table, ACs, test plan updated; AC-1's proof untouched.
  **+281 net across 9 files** (insertions ~+343, deletions ~-62), tolerance ±30% (+197 to +365) and ±2 files — arithmetic shown in the body: 417 − 74 = 343, deletions unchanged. AC-3 loses `TestStampRetriesOnceOnConcurrentMainMove` and its falsifier, keeps the 6m18s baseline, and gains a scope note that the measured quantity (human commits) stays zero on the red path too. AC-1, its fixture, and `TestUnguardedOldLineTagWouldReachStable` are byte-unchanged.
- DONE: All user-facing doc text strict ASD-STE100.
  Re-checked after the edits: 6 blocks, 0 violations, longest sentence 17 words, and the falsifier still reds the mutant with all four rule classes. Three hand dictionary fixes the script cannot make are recorded in the Documentation diff preamble.

### Summary

Mechanical application of a menu I had already priced, so the only judgement calls were where each
cut's accepted consequence gets recorded. I put each one at the site of the cut rather than in a
single list, so a later reader meets the consequence where they would look for the mechanism.

Two things worth the gate's eye. First, the AC-3 scope note: with the retry gone, "zero human
commits" is a claim about the successful path. A concurrent merge reds the job instead, and since the
recovery is a re-run rather than an authored commit, the AC's measured quantity still holds — the
claim narrowed, it did not weaken. Second, dropping the premise guard costs the standing proof of
AC-1's motive, not the protection: the gate still refuses a regressing tag if someone edits
`.goreleaser.yaml`, but nothing would then tell a reader why the gate exists.

One process note. My ASD-STE100 re-check first returned a clean exit from a script that had been
overwritten in the shared scratchpad by another agent — its output named
`internal/gates/prepare_channel_test.go`, a file unrelated to this entity. I caught it because the
output shape did not match my checker's, rewrote the checker at a uniquely-named path, and re-ran
both it and its falsifier. Recording it because a green from the wrong program is exactly the failure
mode this workflow's proof policy names, and the shared scratchpad makes it reachable by accident.

## Stage Report: implementation

- DONE: STACK FIRST: before any work, base your branch on the #760 layer
  `git reset --hard origin/spacedock-ensign/claude-install-sibling-channel-cleanup` ran first; the single deliverable commit 284fa8566 sits directly on 4b5b4073a, so the PR bases on that branch, not `main`.
- DONE: The deliverable lands the approved trimmed design exactly
  `e2e-gate` gains "Gate the cut on the tag not regressing the stable channel"; the stamp moved to `edge-advance` as "Stamp main to the next edge prerelease version" under the pre0 step's own `if:`; the goreleaser remainder is renamed "Advance the stable channel ref to the tagged commit" with its two commands unchanged; `docs/releasing.md` step 9 is gone (9 steps, "Clean up" is now 9) and all six doc blocks are applied verbatim. The three menu cuts STAY cut — no `.goreleaser.yaml` premise-guard test, no cmd exit-contract test, no rebase retry anywhere in the diff.
- DONE: Proofs recorded red-then-green per each AC's failing-today baseline
  Baselines captured by replaying TODAY's step on the same fixtures before any workflow edit: AC-1 — today's step exits 0 and `stable` moves 0.28.0 -> 0.27.1, with `e2e-gate` carrying only its two existing gates; AC-2 — today's step commits main 31b57e18 -> 29a6893b, manifest 0.28.0-pre0 -> 0.26.1 and FO pin 0.28 -> 0.26; AC-3 — main stays 0.27.0 while the auto-cut pre0 tag is v0.28.0-pre0. All six replays then failed red against the unedited workflow ("release.yml has no step named ...") and pass green after it.
- DONE: go test ./... and -race green; gofmt clean
  Both full-suite runs exit 0 after the final edits; `gofmt -l ./cmd ./internal` prints nothing.
- DONE: all comments and doc text strict ASD-STE100
  A checker at a uniquely-named scratchpad path (not the shared one that collided in cycle 3) reports 3 blocks, 0 violations, longest sentence 25 words over the added `docs/releasing.md` prose, the added `release.yml` comments, and the new/changed Go comments. It was FALSIFIED first: against a mutant it reports all five classes (6.3 at 41 words, 3.6 "is repointed by", 8.1 semicolon, banned modal "should", 6.6 at 13 sentences) and exits 1. The first clean run was NOT trusted — the real first pass found 27 violations and every one was rewritten, not waived.
- FAILED: surface within the approved net +281 across 9 files (tolerance +197..+365, ±2 files)
  Measured **net +644 across 10 files** (+741/-97) — 229% of the estimate and +279 past the ceiling. File count is IN tolerance (10 of 7..11); the LOC is not. Overage by file: replay harness +233 (330 vs 97), wiring guard +46 (64 vs 18), cmd wrapper +37 (71 vs 34), release.yml +25, gate library +19, docs +15. The 10th file is `runtime_live_evidence_workflow_test.go`, +2/-2 and net zero — pre-existing gofmt drift the mandated `gofmt -w ./cmd ./internal` corrected, not task content.

### Falsifiers, grouped by claim

- AC-1 blocks: `TestStableRegressionGateBlocksOlderLine` asserts the real gate step exits non-zero and names v0.27.1 and 0.28.0. Falsified by flipping the predicate to `<= 0` (the unit table's equality row reds — verified live) or by deleting the step from `e2e-gate` (`assertStepInJob` reds).
- AC-1 is not redundant with git: `TestUnguardedOldLineTagWouldReachStable` asserts the real stable-push step SUCCEEDS on the same fixture and leaves `stable` at 0.27.1. It reds if the fixture ever stops being a fast-forward, which is the exact masking the review named.
- AC-2 leaves main alone: `TestPatchTagDoesNotStampMain` asserts main's tip SHA and all three stamped files are byte-identical after a v0.26.1 replay. `TestMainStampSharesThePre0DecisionGate` asserts the stamp's `if:` equals the pre0 step's; its twin widens the `if:` to the stable-path-only form and the guard reds.
- AC-3 restores the edge line: `TestLatestLineCutStampsMainToPre0` asserts main lands on 0.28.0-pre0 with FO pin 0.28. Stamping `$RELEASE_VERSION` instead reds it.
- Gate boundary and loudness: the unit table pins `>=` (equality passes, the re-run case) and errors on a prerelease tag, an unparseable manifest, and an empty version.
- Renamed-step consumers: `releaseStampTarget` and `stableRefPushSource` now key on the two new names; without the edit their tests fail, so the change is mandatory rather than cosmetic.

### Summary

All four ACs are satisfied and both full-suite runs are green, but the surface
missed its approved ceiling by a wide margin and I did not massage it to fit.
The scope is exactly what was approved — no test, mechanism, or doc block beyond
the trimmed cut, and the three menu cuts are still cut — so the overage is a
mispriced estimate, not added scope. The evidence: this repo's comparable
existing harness, `edge_advance_decision_shell_test.go`, is 245 lines for 3
tests and 3 helpers (82 lines/test); the new one is 330 lines for 6 tests and 9
helpers (55 lines/test) and additionally builds a two-repo bare-origin fixture,
substitutes the binary, and looks steps up across three jobs. It is denser than
the file it extends, so the 97-line estimate could not have been met by any
implementation of the approved test plan.

A priced menu, if the captain wants the number down. Drop the two carve-out
replays (-35), but audit target 3 asks whether the fail-closed arm can be
flipped without a test reding, and after this cut the answer becomes yes. Drop
the job-identity assertions (-15), and AC-1's "delete the gate step" falsifier
stops reding. Halve the comment density in the new test files (-60 to -80),
which removes the falsification rationale the gate reads. Together they reach
about +545, still 194% — so no combination of them lands inside tolerance. This
needs a captain-visible decision on the estimate, not a code cut.

One thing worth the gate's eye beyond the surface. The AC-2 replay can only
drive the stamp step directly, so it cannot prove that Actions honours the `if:`;
that half rests entirely on the wiring guard, which is why the guard compares the
two conditions for equality rather than merely checking that a gate is present.

## Stage Report: implementation (cycle 2)

Applied the captain's proof-posture ruling (README "Release-machinery proof posture", b59c208d8;
entity Cycle 1). The cycle-1 report above is left intact as the record of what was reshaped.

- DONE: DELETE the bare-origin replay harness entirely — the six replay tests, the two-repo fixture builders, the binary substitution. All of it.
  `internal/release/stable_regression_shell_test.go` deleted in b8e186df5. No replay, no fixture repo, no `go build` substitution remains anywhere in the diff.
- DONE: KEEP the comparator unit table, the structural wiring guards, and the renamed-step consumer updates. Production stays untouched.
  Unit table unchanged (`TestEvaluateStableRegressionGate`, 10 rows). The `if:`-equality guard keeps its twin. Step presence was rehoused from the deleted file into `stable_regression_wiring_test.go` with its own twin. The renamed-step parsers are unchanged from cycle 1. `git diff 284fa8566..b8e186df5` touches only those two test files — zero production bytes.
- DONE: Rewrite the AC "Verified by" lines and the test plan to the new posture, noting the README posture bullet as the authority.
  A new **Proof posture** preamble heads the ACs and cites the README bullet. AC-1 now rests on the unit table plus step presence; AC-2 and AC-3 on the `if:`-equality guard plus `Pre0EdgeVersion`. The test plan's replay section is replaced by **No behavioral replay (cut by the captain's proof-posture ruling)** with the accepted consequence stated. Audit targets are restated for the surviving proof set, and the now-moot "Declared harness limitation" carries a MOOT header so validation does not hunt for a harness.
- DONE: Re-measure the surface with numstat and update the FAILED item per the Cycle 1 line.
  **Net +381 across 10 files** (+478/-97), 136% of the +281 estimate and +16 past the +365 ceiling — down from +644/229%. Supersedes the cycle-1 FAILED item. The harness line item went from +330 to 0; the 67-line `stable_regression_wiring_test.go` rehouses the step-presence guard the estimate had counted inside the deleted 97-line file, so that item is now BELOW estimate. The residual ~+100 is the same per-file under-pricing that predates the harness question: cmd wrapper +37, release.yml +25, `if:`-guard +46, gate library +19, docs +15.
- FAILED: go test ./... and -race green
  Correcting my own cycle-1 claim: it was **not verified**. My cycle-1 command piped `go test` into `grep`, so `$?` was the grep's status and the run's real exit code was never read. I reported green on a masked status. Re-run properly and serially: `go test ./... -count=1 -timeout 20m` exits **0** (493s) and `go test ./... -race -timeout 30m` exits **0**. At the DEFAULT timeout `go test ./...` intermittently fails — see the finding below. No test of this task's fails in any run.

### Finding — `internal/ensigncycle` sits on the 10-minute timeout edge (worker proposal, not a fix)

Preserved rather than repaired; this task has no authorization to change it and it is outside the
approved surface.

- **Released user and normal workflow:** a developer or CI running `go test ./...`, the command CLAUDE.md names under Expected Commands.
- **Observable harm:** the suite aborts with `panic: test timed out after 10m0s` in `internal/ensigncycle`, with no product defect behind it.
- **Affected value AC or boundary:** `none:` — no value AC of this task, and no protected boundary. It is a repo-wide dev-loop hazard that predates this branch. Per the policy `none:` cannot establish Material.
- **Trigger evidence:** at the BASE commit 4b5b4073a, cache busted, under `go test ./... -count=1`, `internal/ensigncycle` took **570.9s of its 600s budget (95.2%, 29s of headroom)** and passed. On this branch the same package measured 490.5s (pass) and 600.4s (timeout) on two runs. The package is untouched by this diff and greps clean for any reference to `internal/release`, `spacedock-release`, `release.yml`, or `releasing.md`. This task strictly REDUCES suite time, because the deleted harness ran ~5s and the surviving additions are sub-second. An early cache-hit reading almost misled me here: a base run showed `(cached)` for that package and looked green without executing it, so I re-ran with `-count=1`.

Proposed as three separate facts: materiality **Deferred risk** (real, but the trigger is machine- and load-dependent and pre-existing); ownership **NOT this task**; disposition **decline here and file separately**. Promote-to-material condition: it fails in CI, or any work depends on `go test ./...` being green at the default timeout. Awaiting FO authorization — no candidate bytes changed for it.

### Summary

The ruling removed the single largest line item and the surface fell from +644 to +381, but I want the
gate to see that the two are less connected than they look: the harness was 330 of the 363-line
overage, and deleting it still leaves ~+100 of ordinary per-file under-pricing that was always there.
The number is honest, not massaged, and it is still +16 outside the band.

The more important item is the verification correction. My cycle-1 "both suites green" rested on an
exit code a pipe had swallowed, which is precisely the failure this workflow's proof policy names —
treating a claim as proven because I wrote it down. Both suites do pass once duration is removed as a
variable, and I established that with a base-commit comparison rather than by asserting it, but the
cycle-1 report should not have said green.

One thing worth the gate's eye on the new posture. AC-1's shell arms — the `ls-remote` exit-2
first-release carve-out and the fail-closed arm — now have no standing test, by design. That is
recorded as an accepted consequence at the site of the cut rather than only here, and the first real
stable cut is where those arms are observed.
