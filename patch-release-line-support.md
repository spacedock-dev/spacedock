---
id: d0g21c517b5nvga1ybwckapk
title: Patch-release line support - gate release.yml's main-stamp on line-latestness, fix the stable-branch advance, automate the preversion bump
status: ideation
source: "Captain CL, 2026-08-25, reconciling gr and tw after the v0.27.0 cut: 'reconcile tw and gr and recommend best approach' - supersedes next-independent-release-line (twq68r4y8qg0wetztajtmmzz), whose body described the retired next-branch model; the live incidents of 2026-08-25 are the spec"
started: 2026-08-25T16:33:56Z
completed:
verdict:
score:
worktree:
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

### C. One bounded retry on the main-stamp push

Today's stamp step usually finds no diff and commits nothing. After mechanism B it commits on every
latest-line cut, so a concurrent merge to `main` between fetch and push becomes a realistic
non-fast-forward rejection. One `git pull --rebase origin main`, then re-push.

- **Value AC served:** AC-3, at its margin.
- **Corrected rationale (the earlier draft overstated this).** Without the retry the rejection reds
  `edge-advance` and the recovery is a one-click, idempotent JOB RE-RUN — the pre0 tag is never
  re-minted and the stamp is idempotent — NOT the hand-authored commit AC-3's 6m18s baseline
  measured. The retry therefore buys "no human job re-run", not "no outage". It is kept at that
  smaller value: eight lines against a foreseeable red job on the one step whose whole purpose is to
  need no human.
- *Alternative:* no retry, accept the rare rejection and the re-run. *Why insufficient (weakly):* the
  re-run is manual attention on a path this task exists to make unattended. This is the least
  load-bearing item in the cut and the cheapest to drop at the gate.

### D. Retire the manual ritual

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

**Declared harness limitation.** The replayed steps invoke `go run ./cmd/spacedock-release`, which
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

Estimate net LOC change: **+355, across 11 files** (insertions ~+417, deletions ~-62).
Tolerance: ±30% on the net figure (+249 to +462) and ±2 files.

| File | ins | del | Change |
|---|---|---|---|
| `.github/workflows/release.yml` | 72 | 36 | new gate step in `e2e-gate`; stamp step moved to `edge-advance` with the retry; goreleaser step renamed |
| `internal/release/stable_regression_gate.go` | 32 | 0 | new — the decision over `ComparePreVersion` |
| `cmd/spacedock-release/stable_regression_gate.go` | 34 | 0 | new — thin wrapper over `ManifestVersion` + the decision |
| `cmd/spacedock-release/main.go` | 6 | 0 | dispatch case, usage line, doc comment |
| `internal/release/stable_regression_gate_test.go` | 58 | 0 | new — unit table, including the `>=` boundary |
| `cmd/spacedock-release/stable_regression_gate_test.go` | 34 | 0 | new — exit contract 0/1/2 |
| `internal/release/stable_regression_shell_test.go` | 115 | 0 | new — bare-origin replay of the real step, plus the unguarded baseline |
| `internal/release/edge_advance_wiring_test.go` | 18 | 0 | structural guard that the moved stamp shares the pre0 gate |
| `internal/release/channel_agreement_guard_test.go` | 14 | 8 | two parsers follow the renamed/split steps |
| `internal/release/goreleaser_guard_test.go` | 14 | 0 | pins AC-1's harm premise: `prerelease: auto` and the stable cask's `skip_upload: auto` |
| `docs/releasing.md` | 20 | 18 | step 9 retired; three passages corrected; one short section added |

**Honest note on the estimate, against the ruling that expected "well under half of +350".** It did
not come out that way, and the reason is worth the gate's attention rather than a massaged number.
The dropped mechanism cost roughly 200 insertions, of which about 140 were its test surface. The
added centerpiece costs roughly the same, and for the same reason: the review's own sharpest
requirement — that the fixture prove the UNGUARDED path would publish — can only be met by a
bare-origin replay of the real steps, which is the single largest line item here (115 lines) and is
not compressible without gutting AC-1's baseline. **The lean cut's saving is in authority and risk,
not in lines:** no force-push capability, no new push-path wiring, no ancestry semantics in the
release path, and one fewer irreversible failure mode. If the gate wants the number down, the
available cuts are the goreleaser premise guard (−14), the cmd exit-contract test (−34), and
mechanism C with its test (−26) — about −74, landing near +280. Cutting deeper means cutting AC-1's
proof, which is the thing this reshape exists to add.

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

Each AC's baseline is the SAME replay run against the CURRENT `release.yml`, which must produce the
stated failing result. That baseline is independent of this task's code and can move the wrong way.

**AC-1 (value) — A stable tag older than the release `stable` serves never reaches publication: the
run fails in `e2e-gate`, and on the same fixture the unguarded path DOES reach a published state.**
Verified by: `TestStableRegressionGateBlocksOlderLine` in
`internal/release/stable_regression_shell_test.go` — spike TEST E's fixture (`stable` at the 0.28.0
commit, a CHILD commit stamped 0.27.1 and tagged `v0.27.1`), replaying the real gate step and
asserting a non-zero exit with the tag and both versions named in stderr. Paired with
`TestUnguardedOldLineTagWouldReachStable`, which on the SAME fixture runs the real stable-push step
WITHOUT the gate and asserts the push SUCCEEDS and `stable` then serves 0.27.1 — the independent
baseline, and the proof the block is not redundant with git's own ancestry check. Falsifying change:
flip the gate's comparison to `<= 0` (the pass case reds) or delete the gate step from `e2e-gate`
(the structural check in `stable_regression_shell_test.go` reds). Today's baseline: the gate step
does not exist, the run is green, and `TestUnguardedOldLineTagWouldReachStable` describes exactly
what today's pipeline does.

**AC-2 (value) — A patch tag that PASSES the regression gate still leaves `main` untouched: after an
`edge-advance` replay on a `v0.27.1` tag whose decision output is `advance=false`, `main`'s tip SHA
and the bytes of `.claude-plugin/plugin.json`, `.codex-plugin/plugin.json`, and the FO shared-core
prose are identical to their pre-run values.**
Verified by: `TestPatchTagDoesNotStampMain` in the same file — a fixture repo tagged `v0.26.0`,
`v0.27.0`, `v0.28.0-pre0`, replaying the `decision` step and then the stamp step under the decision
it produces for `v0.26.1`, asserting `main` is byte-identical; plus a structural check in
`edge_advance_wiring_test.go` that the stamp step's `if:` names the same
`steps.decision.outputs.advance` condition as the auto-pre0 step, with its adversarial twin
(diverge them → red). The structural half is load-bearing: a replay drives the step directly and
cannot prove Actions gates it. Falsifying change: remove the gate from the stamp step's `if:`, or
widen it back to `!contains(github.ref, '-')` — the wiring guard reds. Today's baseline: the same
replay stamps `main` DOWN to 0.26.1.

**AC-3 (value) — A latest-line stable cut needs zero human commits to restore the edge line: after a
full `edge-advance` replay on a latest-line stable tag, `main`'s manifest version equals the auto-cut
pre0 tag's version and the FO prose pin equals its major.minor.**
Verified by: `TestLatestLineCutStampsMainToPre0` in the same file — replay on `v0.27.0` against a
fixture whose `main` carries 0.27.0, asserting `main`'s tip has `version: 0.28.0-pre0` and
`These skills require binary minor 0.28`, matching the `v0.28.0-pre0` tag the same job cut; plus
`TestStampRetriesOnceOnConcurrentMainMove`, which advances the fixture's origin `main` between the
step's fetch and its push and asserts the step still lands the commit. Falsifying change: stamp
`$RELEASE_VERSION` instead of the `edge-pre0-version` output — `main` lands on 0.27.0, the
manifest/tag pair disagrees, and the test reds; or delete the rebase retry — the concurrent-move test
reds. Today's baseline: the same replay leaves `main` at 0.27.0 while the pre0 tag is 0.28.0-pre0 —
the exact 6m18s / one-hand-commit mismatch measured at the v0.27.0 cut (pre0 tag 21:55:31, hand
repair b8346ffc9 at 22:01:49).

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

**Command exit contract — `cmd/spacedock-release/stable_regression_gate_test.go` (low cost).** Pass
exits 0, block exits 1 with the tag and both versions in stderr, a missing argument exits 2. The
`if:`-gated step depends on this contract and no other test reaches the usage-error path.

**Behavioral replay — `internal/release/stable_regression_shell_test.go` (the bulk of the cost).**
Extends the shipped `edge_advance_decision_shell_test.go` harness with a bare-origin fixture. Each
test extracts the REAL step's `run` block from the on-disk `release.yml` and executes it. Carries the
four AC tests above plus `TestStableRegressionGatePassesWhenStableRefAbsent` (spike TEST A's exit-2
carve-out) and `TestStableRegressionGateFailsClosedOnUnreadableRemote` (spike TEST B). The
fixture-construction helper is shared across all cases and drives them from one table; its shape is
already exercised by spike 1, so the risk is front-loaded.

**Structural guards.** `edge_advance_wiring_test.go` gains the AC-2 check that the stamp step's `if:`
and the auto-pre0 step's `if:` name the same `steps.decision.outputs.advance` condition, reusing the
shipped `ifHasDecisionGate` helper, with the adversarial twin (diverge them → red).
`channel_agreement_guard_test.go`'s `releaseStampTarget` and `stableRefPushSource` parsers both key
on the step name `"Stamp plugin manifests to the release version"`, which this task splits and
renames — without the edit those tests fail, so this is mandatory, not optional. Both keep their
adversarial twins. `goreleaser_guard_test.go` gains two asserts pinning AC-1's harm premise —
`release.prerelease: auto` and the stable cask's `skip_upload: auto` — so the reason this gate exists
cannot silently evaporate from `.goreleaser.yaml`.

**Full suite.** `go test ./...` and `go test ./... -race`, plus `gofmt -w ./cmd ./internal`.

**No live workflow test.** No live lane loads or drives any changed file (see Path-to-lane call). The
one claim a replay cannot make is that GitHub Actions itself evaluates a step `if:` as expected; that
is covered structurally, and the first real stable cut after merge is the live confirmation. Naming
this now: **AC-3's live behavior is proven by fixture replay, not by a live cut, and the next real
stable release is where it is observed.** That observation is deliberately NOT an AC, because it
cannot be reproduced by validation on demand.

**Detached adversarial audit targets.** Sharpest questions for the auditor: (1) Does
`TestUnguardedOldLineTagWouldReachStable` genuinely succeed at the push, or does the fixture's
ancestry make git refuse it anyway and mask a gate that never fires? This is the same hole the review
found in the earlier design and the reason spike TEST E exists — re-check it against the built
fixture, not against this claim. (2) Can the gate step be edited to compare against the TAG pool
instead of the `stable` ref with every test still green? (3) Can the `ls-remote` fail-closed arm be
flipped to fail-open without a test reding? (4) Can the stamp step be edited to stamp `main` while
the pre0 tag was never pushed, with every test still green?

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
entry" instead of the first draft's "The stable entry is NOT repointed … by a hand-edit". New
`release.yml` comments and Go doc comments follow the same rules. The moved stamp step's existing comment block converts as this task touches
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
  `stable` serves. See "The Stable Regression Gate" below.
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
gate with no human step. If a concurrent merge to `main` rejects the push, the
step rebases one time and pushes again.
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

A stable tag that is older than the release that `stable` serves can damage two
binary channels. goreleaser marks each bare `vX.Y.Z` tag as a full release, so
GitHub moves `/releases/latest` to it. goreleaser also bumps the stable Homebrew
cask, because `skip_upload: auto` skips only a prerelease. An old tag moves both
channels DOWN, and every job stays green. The `e2e-gate` job stops this before
goreleaser starts.

- The gate reads the version in the `.claude-plugin/plugin.json` file that
  `stable` serves.
- If the tag is older than that version, the gate fails the run. goreleaser does
  not start, because it needs the `e2e-gate` job.
- If the tag is the same version or newer, the gate lets the run continue. A
  re-run of a release that already reached `stable` is still possible.
- On the first stable release the `stable` ref does not exist. The gate writes a
  notice and lets the run continue.
- If the gate cannot read the `stable` ref, it fails the run. A read failure is
  not a permission to publish.

The gate does not make a patch line deliverable. A patch that is newer than the
version `stable` serves still publishes, and `stable` then leaves the history of
`main`. The next latest-line release cannot advance `stable` after that. Do not
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
