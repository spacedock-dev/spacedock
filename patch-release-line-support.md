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
---

One owner for release.yml's stable-tag conditional. Today the machinery cannot ship an old-line patch (a v0.27.1 while main is the 0.28 line): two step-level defects, verified by reading the shipped workflow at the v0.27.0 cut, plus one automation the cut proved missing.

## Problem

`release.yml` has one conditional for stable tags — `!contains(github.ref, '-')` — and it decides
three unrelated questions at once: does `main` get stamped, does `stable` advance, and does the edge
line move. On the only shape the repo has ever cut (a stable release from `main`'s tip) all three
answers coincide, so the collapse has never been visible. A patch on an older line answers them
differently, and the single conditional gets two of the three wrong.

**1. The main stamp has no line-awareness.** "Stamp plugin manifests to the release version"
(`.github/workflows/release.yml:247`) fires on any hyphen-free tag. On a `v0.27.1` tag it switches to
`main` and stamps `main` DOWN to 0.27.1 — rewriting the 0.28.0-pre0 manifests and the FO prose pin
(`These skills require binary minor 0.28`), so every edge user's plugin claims 0.27 while their
binary is 0.28.0-pre0 and the FO binary gate aborts. This is worse than the v0.27.0 incident: it
breaks a currently-working state rather than failing to repair a broken one. The latest-line decision
already exists in the sibling `edge-advance` job, and this step does not consult it.

**2. The stable advance is not the failure the seed described, and the real failure is worse.**
The seed expected the patch push to die non-fast-forward. The spike (Risk evidence below) refutes
that: a patch branched off the current `stable` tip pushes as a clean fast-forward, because `stable`
points at the prior stable release's commit and the patch is that commit's child. The push succeeds.
What breaks is the cut AFTER it. Once `stable` sits on a patch commit that is not on `main`'s
history, the next latest-line stable cut from `main` is non-fast-forward and is REJECTED, so
`stable` freezes on the patch line permanently and every later stable release is invisible to
`spacedock@spacedock` installs. Separately, an old-line patch cut while `stable` already serves a
newer release is also rejected — semantically the right outcome, but it arrives as a red job after
goreleaser has already published, not as a deliberate skip. The step's own comment ("The tagged
commit is on main's history, so this fast-forwards", `release.yml:283`) states a reason that is false
for every patch-line commit; the push happens to work for the wrong reason.

**3. Nothing stamps `main` past the released minor.** The auto-pre0 job tags `vX.(Y+1).0-pre0` and
publishes that edge binary, but leaves `main`'s manifests and FO pin at the released version. Every
edge install then aborts at the FO version gate until a human commits the bump. Observed live at the
v0.27.0 cut: the pre0 tag landed 2026-08-24 21:55:31 PDT, the failing install hit 17f5cd591 at
22:01, and the hand repair b8346ffc9 landed 22:01:49 — a 6m18s outage closed by one human commit.
`docs/releasing.md` step 9 (b04b3effd) is the interim procedure; this item retires it.

Items 1 and 2 are latent — verified by reading and exercising the shipped workflow, never yet
triggered live. Item 3 is the observed incident. All three share one root cause and one fix surface.

## Proposed approach

Split the one overloaded conditional into the two questions it actually conflates, and give each an
owner that already exists or is a thin wrapper over shipped code.

**The two questions are genuinely different.** "Does the edge line move?" and "Does the stable
channel move?" have different answers for the same tag. Cutting `v0.27.1` today: the edge line must
NOT move (`main` is already 0.28.0-pre0, and a 0.28 edge binary is published), but the stable channel
MUST move (stable serves 0.27.0, and patch users need 0.27.1). Any design that answers both with one
predicate gets one of them wrong. This is the necessity argument for two decisions rather than one.

### A. Main stamp — move it into the job that already owns the decision

Move the stamp step out of `goreleaser` and into `edge-advance`, placed AFTER the auto-pre0 step and
carrying the SAME gate the pre0 step carries (`steps.decision.outputs.advance == 'true'`). Change
what it stamps: `edge-pre0-version "$RELEASE_VERSION"` (X.(Y+1).0-pre0), not `$RELEASE_VERSION`.
Rename it to `Stamp main to the next edge prerelease version`, because it no longer stamps the
release version.

- **Value AC served:** AC-1 (an old-line tag never touches `main`) and AC-3 (the edge gate passes
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
- **Ordering inside the job matters.** The pre0 tag is cut and its run is verified FIRST; the stamp
  runs only after that verify poll passes. A failure before the stamp therefore leaves exactly
  today's recoverable state, never a worse one.

### B. Stable advance — its own step, its own predicate, in the goreleaser job

The remainder of the old stamp step becomes a step of its own, `Advance the stable channel ref to
the tagged commit`. It fetches `stable`, reads that branch's OWN `.claude-plugin/plugin.json`
version, and advances only when the tag's version is strictly greater:

```bash
RELEASE_COMMIT="$(git rev-list -1 "$GITHUB_REF_NAME")"
if git fetch origin stable; then
  STABLE_SHA="$(git rev-parse FETCH_HEAD)"
  git show "$STABLE_SHA:.claude-plugin/plugin.json" > "$STABLE_MANIFEST"
  if [ "$(go run ./cmd/spacedock-release stable-advance-decision "$GITHUB_REF_NAME" "$STABLE_MANIFEST")" != "advance" ]; then
    echo "::notice::stable ref NOT advanced for $GITHUB_REF_NAME: the stable channel already serves a newer release"
    exit 0
  fi
  git push --force-with-lease="refs/heads/stable:$STABLE_SHA" origin "$RELEASE_COMMIT:refs/heads/stable"
else
  git push origin "$RELEASE_COMMIT:refs/heads/stable"   # first stable release: the ref does not exist yet
fi
```

- **Value AC served:** AC-2 (the newest stable release reaches stable-channel users, across lines).
- **New mechanism 1 — `stable-advance-decision`.** A new `spacedock-release` subcommand mirroring
  `edge-advance-decision`'s signature and exit contract (prints `advance`/`skip`, exit 0 for both,
  non-zero on unparseable input). It wraps the already-shipped `ComparePreVersion`.
  - *Simplest alternative:* reuse `edge-advance-decision` with `stable`'s manifest as the known
    version. *Why insufficient:* it computes `DevPreVersion(tag)` as the target, so it answers the
    edge question, not the stable one — for `v0.27.1` it compares 0.28.0-pre1 against 0.27.0 and says
    advance for the wrong reason, and its answer diverges from the stable question in exactly the
    cases this task exists to fix.
  - *Second alternative:* compare in shell with `sort -V`. *Why insufficient:* prerelease ordering
    (`pre2` before `pre10`, stable above its own prereleases) is the distinction the shipped
    comparator exists to make; re-deriving it in YAML duplicates the one function that already has
    tests.
- **New mechanism 2 — `--force-with-lease`.** The lease is required, not defensive. Once a patch has
  moved `stable` off `main`'s history, the next latest-line advance is non-fast-forward by
  construction (spike TEST 2), so a plain push can never publish it.
  - *Simplest alternative:* plain push, and treat a non-fast-forward rejection as a skip. *Why
    insufficient:* spike TEST 2 shows that turns the permanent stable-channel freeze from a loud
    failure into a silent one. Git's ancestry check cannot distinguish "older line, correctly
    refused" from "newer line, wrongly refused" — both are non-fast-forward.
  - The version predicate is what makes the force safe: it establishes this tag is newer than what
    `stable` serves BEFORE any force. The lease then guards the read-to-write window and refuses on a
    concurrent move (spike TEST 4). Both are proven, not assumed.
- **The FF-vs-force branch is deliberately absent.** One path — gate, then lease-push — covers both,
  because a fast-forward under a correct lease succeeds too. A branch on ancestry would be a third
  mechanism with no case of its own.

### C. Retire the manual ritual

`docs/releasing.md` step 9 is deleted and step 10 renumbered. Sections "What the Tag Push Does" and
"Advancing the Edge Line" are corrected, and a short "Advancing the Stable Channel" section is added.
The concrete diff is in **Documentation diff** below.

### One bounded retry on the main-stamp push

Today's stamp step usually finds no diff and commits nothing. After this change it commits on every
latest-line cut, so a concurrent merge to `main` between fetch and push becomes a realistic
non-fast-forward rejection — which would reintroduce the manual bump AC-3 exists to remove. One
`git pull --rebase origin main` retry, then re-push. *Alternative:* no retry, accept a rare
rejection. *Why insufficient:* the rejection's cost is precisely the failure mode this task closes.

## Risk evidence

The riskiest unverified mechanism was the stable-ref push: the seed asserted it dies non-fast-forward
for a release-branch commit, and the whole shape of AC-2 depended on that being true. Exercised
first, in a throwaway bare-origin fixture. The script is inlined here rather than left in a
scratchpad so any reviewer can re-run it (bash, not zsh — zsh's `:r` modifier mangles the refspecs):

```bash
set -uo pipefail
S=$(mktemp -d); cd "$S"
git init -q --bare origin.git && git clone -q origin.git work && cd work
git config user.email a@b.c; git config user.name t
c() { echo "$1" > f.txt; git add f.txt; git commit -qm "$1"; }
c A; c B_v0_27_0; git tag v0.27.0; git push -q origin main
B=$(git rev-parse HEAD); git push -q origin "${B}:refs/heads/stable"   # stable at the v0.27.0 commit
c C; c D_0_28_line; git push -q origin main; D=$(git rev-parse HEAD)   # main advances
git switch -q -c release/0.27.1 "$B"; c P_v0_27_1; P=$(git rev-parse HEAD)  # patch off the stable tip
git push origin "${P}:refs/heads/stable"                                     # TEST 1
git push origin "${D}:refs/heads/stable"                                     # TEST 2
git push --force-with-lease="refs/heads/stable:${P}" origin "${D}:refs/heads/stable"  # TEST 3
git push --force-with-lease="refs/heads/stable:${P}" origin "${P}:refs/heads/stable"  # TEST 4 (stale lease)
git push origin "${P}:refs/heads/stable"                                     # TEST 5 (backwards)
git fetch -q origin stable && git rev-parse --short FETCH_HEAD               # TEST 6
git fetch -q origin nosuchbranch                                             # TEST 7
```

Results:

| Test | Setup | Result |
|---|---|---|
| 1 | patch commit P (child of `stable`) → `stable` | **fast-forward SUCCESS** — refutes the seed |
| 2 | later `main` commit D (not a descendant of P) → `stable` | **rejected, non-fast-forward** |
| 3 | D → `stable` with `--force-with-lease=refs/heads/stable:<P>` | **forced update SUCCESS** |
| 4 | same lease, stale expected SHA | **rejected (stale info)**, `stable` unchanged |
| 5 | old-line P → `stable` when `stable` is at D | **rejected, non-fast-forward**, unchanged |
| 6 | `git fetch origin stable` + `rev-parse FETCH_HEAD` | resolves the lease SHA and the manifest read |
| 7 | fetch a `stable` ref that does not exist | **exit 128** — the first-release case needs its own branch |

This inverted item 2's design: the defect is not a failed patch publish, it is a permanently frozen
stable channel on the NEXT cut, and the fix needs a lease rather than a refusal message.

Second spike, against this repo's real 71-tag pool, confirming the seed's claim that the existing
decision helper already ranks the patch case (so no change to `edge-advance-decision` is needed):

```
v0.27.1   known=0.28.0-pre1  -> skip      (correct: main is already the 0.28 line)
v0.28.0   known=0.28.0-pre1  -> advance   (correct: the next latest-line cut)
v0.26.1   known=0.28.0-pre1  -> skip      (correct: two lines back)
```

Third: the harness for AC-1/AC-2/AC-3 is not new. `internal/release/edge_advance_decision_shell_test.go`
already extracts a real `release.yml` step's `run` block and executes it against a fixture git repo
via `GIT_DIR`/`GIT_WORK_TREE` with a stubbed `$GITHUB_OUTPUT`. This task extends that proven harness
with a bare origin remote (spike TEST 1-7 shape) rather than inventing one.

**Declared harness limitation.** The replayed steps invoke `go run ./cmd/spacedock-release`, which
resolves relative to the process working directory. Running the script with the working directory at
the repo root would let `stamp-version` write this checkout's real manifests. The harness therefore
runs with the working directory in the fixture and substitutes that one token for a pre-built binary
path. Every other byte of the script — the gating, the refspecs, the lease, the conditionals — runs
verbatim. This substitution is declared here so validation checks it rather than discovering it.

## Out of scope

The e2e-gate/waiver mechanics (work for branch SHAs today). The marketplace display-version fields
(cosmetic; tidied 2026-08-25, automation optional here). The release notes ritual. The release ritual's
own patch-line steps (step 2 branches off `origin/main`; a patch branches off the prior stable tag
instead) — a documentation change worth making, but it is operator procedure, not machinery, and
folding it in here widens the diff without serving a value AC. Recorded as a follow-up.

## Expected surface and tolerance

Estimate net LOC change: **+350, across 9 files** (insertions ~+420, deletions ~-70).
Tolerance: ±30% on the net figure (+245 to +455) and ±2 files.

| File | Change |
|---|---|
| `.github/workflows/release.yml` | split the stamp step; move the main stamp into `edge-advance`; new stable-advance step |
| `cmd/spacedock-release/stable_advance_decision.go` | new — the subcommand |
| `cmd/spacedock-release/main.go` | dispatch case, usage line, doc comment |
| `internal/release/edge_advance_decision.go` | `StableAdvanceDecision` over the shipped comparator |
| `internal/release/stable_advance_decision_test.go` | new — unit tests for the decision |
| `internal/release/stable_advance_shell_test.go` | new — bare-origin fixture replay of the real steps |
| `internal/release/channel_agreement_guard_test.go` | parsers follow the renamed/split steps |
| `internal/release/edge_advance_wiring_test.go` | structural guard that the stamp shares the pre0 gate |
| `docs/releasing.md` | step 9 retired; three sections corrected |

The test files carry roughly 60% of the insertions. The seeded ~+60 estimate counted the machinery
only and did not include the replay harness.

**Observable semantics this task changes** (cost is lines; these are the boundary):

1. **Which refs a stable tag mutates.** An old-line stable tag will mutate neither `main` nor
   `stable`. Today it mutates both.
2. **What `main` carries after a latest-line cut.** The plugin manifests and the FO prose pin move to
   `X.(Y+1).0-pre0` instead of `X.Y.Z`. This is a change to on-disk content the edge marketplace
   serves.
3. **New authority: CI can rewrite `stable`'s history.** `--force-with-lease` lets the release run
   move `stable` non-fast-forward. Today CI can only extend it. This is the most consequential item
   here and the one that most needs the captain's explicit approval; the version gate and the lease
   bound it, but the authority is genuinely new.
4. **New command grammar:** `spacedock-release stable-advance-decision <tag> <manifest>`.
5. **Failure attribution moves.** A main-stamp failure will red `edge-advance`, not `goreleaser`.
6. **A documented procedure is removed:** `docs/releasing.md` step 9.

Not changed: the e2e-gate, the manifest-tag-gate, goreleaser itself, the journey ledger, the pre0 tag
mechanics, the marketplace repo, and the `edge-advance-decision` logic (reused verbatim).

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

**AC-1 (value) — An old-line patch tag leaves `main` untouched: after a full `edge-advance` replay on
a non-latest-line stable tag, `main`'s tip SHA and the bytes of `.claude-plugin/plugin.json`,
`.codex-plugin/plugin.json`, and the FO shared-core prose are identical to their pre-run values.**
Verified by: `TestOldLineTagDoesNotStampMain` in `internal/release/stable_advance_shell_test.go` —
a fixture repo tagged `v0.26.0`, `v0.27.0`, `v0.28.0-pre0`, replaying the `decision` step and then
the stamp step under the decision it produces for `v0.26.1`, asserting `main` is byte-identical.
Falsifying change: remove the `steps.decision.outputs.advance` gate from the stamp step's `if:`, or
widen it back to `!contains(github.ref, '-')` — the fixture's `main` then moves to 0.27.0-pre0 and
the test reds. Today's baseline: the same replay stamps `main` DOWN to 0.26.1.

**AC-2 (value) — The stable channel serves the newest stable release across a line change: after a
patch cut off the stable line and then a latest-line stable cut from `main`, `refs/heads/stable`
points at the SECOND cut's commit, and no run ever moves it to an older release's commit.**
Verified by: `TestStableAdvanceCrossesLines` in the same file — the spike's TEST 1/2/3 sequence
driven through the real step script against a bare-origin fixture, asserting `stable` equals the
patch commit after cut one and the `main` commit after cut two; plus
`TestStableAdvanceRefusesOlderLine` asserting `stable` is unchanged and the step exits 0 with its
`::notice::` when an old-line tag is cut against a newer `stable`. Falsifying change: drop
`--force-with-lease` (cut two reds non-fast-forward, spike TEST 2) or drop the
`stable-advance-decision` gate (the old-line cut force-moves `stable` backwards and
`TestStableAdvanceRefusesOlderLine` reds). Today's baseline: cut two is rejected and `stable` stays
frozen on the patch commit.

**AC-3 (value) — A latest-line stable cut needs zero human commits to restore the edge line: after a
full `edge-advance` replay on a latest-line stable tag, `main`'s manifest version equals the auto-cut
pre0 tag's version and the FO prose pin equals its major.minor.**
Verified by: `TestLatestLineCutStampsMainToPre0` in the same file — replay on `v0.27.0` against a
fixture whose `main` carries 0.27.0, asserting `main`'s tip has `version: 0.28.0-pre0` and
`These skills require binary minor 0.28`, matching the `v0.28.0-pre0` tag the same job cut.
Falsifying change: stamp `$RELEASE_VERSION` instead of the `edge-pre0-version` output — `main` lands
on 0.27.0, the manifest/tag pair disagrees, and the test reds. Today's baseline: the same replay
leaves `main` at 0.27.0 while the pre0 tag is 0.28.0-pre0 — the exact 6m18s / one-hand-commit
mismatch measured at the v0.27.0 cut (pre0 tag 21:55:31, hand repair b8346ffc9 at 22:01:49).

**AC-4 — The manual ritual is gone, not merely superseded: `docs/releasing.md` contains no post-tag
`main` preversion bump step, and its stable-advance and edge-advance prose describe the gated
behavior.**
Verified by: the Documentation diff below applied, `docs/releasing.md` step count reduced by one with
step 10 renumbered to 9, and the changed prose passing the `simple-english` check the workflow README
requires. Falsifying change: leave step 9 in place — the doc then documents a hand bump the pipeline
performs, and a cutter runs a duplicate stamp. This AC is a mechanism AC and counts only paired with
AC-3, which measures the value it serves.

## Test plan

**Unit — `internal/release/stable_advance_decision_test.go` (low cost).** Table over
`StableAdvanceDecision`: newer patch vs older stable (advance); older line vs newer stable (skip);
equal versions (skip — the boundary is strict `>`, matching `EdgeAdvanceDecision`); a hyphenated tag
(error, since the caller's `if:` guarantees a bare tag); unparseable manifest version (error, so a
miswiring fails loud rather than silently skipping). Plus a `cmd` exit-contract test: `advance` and
`skip` both exit 0, bad input exits non-zero.

**Behavioral replay — `internal/release/stable_advance_shell_test.go` (the bulk of the cost).**
Extends the shipped `edge_advance_decision_shell_test.go` harness with a bare-origin fixture. Each
test extracts the REAL step's `run` block from the on-disk `release.yml` and executes it. Carries the
three AC tests above plus `TestStableAdvanceCreatesMissingStableRef` (spike TEST 7's exit-128 path)
and `TestStableAdvanceLeaseRefusesConcurrentMove` (spike TEST 4). Complexity is moderate and the risk
is front-loaded: the fixture shape is already exercised by the spike script.

**Structural guards.** `edge_advance_wiring_test.go` gains a check that the stamp step's `if:` and
the auto-pre0 step's `if:` name the same `steps.decision.outputs.advance` condition, with the
adversarial twin (diverge them → red). `channel_agreement_guard_test.go`'s `releaseStampTarget` and
`stableRefPushSource` parsers follow the renamed and split steps; `stableRefPushSource` must accept a
flag-bearing push, since `--force-with-lease` breaks its current exact-four-field match. Both keep
their adversarial twins.

**Full suite.** `go test ./...` and `go test ./... -race`, plus `gofmt -w ./cmd ./internal`.

**No live workflow test.** No live lane loads or drives any changed file (see Path-to-lane call). The
one claim a replay cannot make is that GitHub Actions itself evaluates the step `if:` as expected;
that is covered structurally, and the first real latest-line cut after merge is the live
confirmation. Naming this now: **AC-3's live behavior is proven by fixture replay, not by a live cut,
and the next real stable release is where it is observed.** That observation is deliberately NOT an
AC, because it cannot be reproduced by validation on demand.

**Detached adversarial audit targets.** Sharpest questions for the auditor: can the stamp step be
edited to stamp `main` while the pre0 tag was never pushed, with every test still green? Can the
lease be widened to a bare `--force` without a test reding? Does `TestStableAdvanceRefusesOlderLine`
actually fail when the decision gate is removed, or does the fixture's ancestry make git refuse the
push anyway and mask the missing gate? That last one is the likeliest hole in this design's own
tests, and the fixture must be built so the old-line push WOULD succeed if unguarded.

## Sequencing decision for the captain

Two ways to get the #760 fix ("claude install leaves the sibling edge plugin installed and enabled")
to users. Both require this task to ship first — see the one-way constraint below.

**Option 1 — fold into v0.28.0 (recommended).** Land the #760 fix on `main`, ship it in the next
minor. No patch-line cut, no release-branch, no exercise of the new machinery in production.
Cheapest, and it is what the current state supports: **#760 is still OPEN and unmerged**, and its fix
is itself only at ideation in this session, so no v0.27.1 can be cut today regardless.

**Option 2 — cut v0.27.1 from a release branch.** After #760 merges, branch off the `v0.27.0` tag,
cherry-pick the fix, stamp 0.27.1, green that SHA, tag. This is the only way to observe the patch
line end-to-end in production, and it is the scenario the ACs model. It costs a full release cut
(green run, notes, tag) and it is the first time `stable` would leave `main`'s history.

My recommendation is Option 1 for #760, and to treat the first genuinely urgent old-line patch as the
live confirmation. The fixture replays cover the ACs; spending a release cut purely to exercise
machinery the replays already prove is not a good trade.

**The one-way constraint, either way.** Do NOT cut a v0.27.1 on today's machinery. The stamp step
would switch to `main` and stamp it DOWN to 0.27.1, rewriting the 0.28.0-pre0 manifests and the FO
pin to 0.27 — breaking every edge user, whose binary is 0.28.0-pre0. This task ships before any
patch-line cut, not after.

## Documentation diff

`docs/releasing.md` is user-facing documentation, so every "After" block below follows ASD-STE100 per
the workflow README's Prose style section. The text was checked with the `simple-english` rule
catalog, not merely declared compliant: the first draft of this section broke Rule 8.1 (a semicolon),
Rule 6.3 (two sentences at 26 and 32 words), Rule 6.6 (a seven-sentence paragraph), Rule 3.6 (a
passive with a known agent), and GR-1 (a dropped "that"). The blocks below are the corrected text and
are ready to paste. New `release.yml` comments and Go doc comments follow the same rules. The moved
stamp step's existing comment block converts as this task touches it, per the README's
convert-on-touch clause. This task body itself is workflow state, not user-facing documentation, so
it is out of the rule's scope.

**1. `docs/releasing.md`, "What the Tag Push Does" — replace the stamp bullet.**

Before:
```
- stamps the plugin manifests' `version` on `main`, then advances the stable
  channel ref (see below).
```
After:
```
- advances the stable channel ref to the tagged commit. The run advances the ref
  only when the tag is newer than the release that `stable` serves. See
  "Advancing the Stable Channel" below.
- stamps the plugin manifests and the FO prose pin on `main`, but only on a
  latest-line stable tag. The stamp writes the auto-cut `X.(Y+1).0-pre0` version.
  An old-line patch tag does not touch `main`.
```

**2. Replace the paragraph beginning "The stable entry is NOT repointed per release" (lines 40-46).**

After:
```
The stable entry is NOT repointed per release by a hand-edit in the marketplace
repo. `stable` is a MOVING BRANCH in this repo. After the tag fires, release.yml
reads the version that `stable` serves now. The run advances `stable` to the
tagged commit only when the tag is newer than that version. The push uses a
lease, so a patch line and a later minor can both reach the channel.

A fresh `spacedock@spacedock` install resolves whatever `stable` points at. That
push is what publishes the release to the stable channel. No marketplace-repo
commit is necessary.
```

**3. Replace the paragraph beginning "The post-tag manifest stamp is idempotent" (lines 48-51).**
The current text is now false. The stamp writes the pre0 version, which always differs from the
tagged commit's version.

After:
```
The post-tag stamp always makes a commit on a latest-line cut. The stamp writes
the next prerelease version, and the tagged commit carries the release version.
An old-line tag makes no commit.
```

**4. Delete step 9 entirely (lines 183-197), renumber step 10 to step 9.**

**5. "Advancing the Edge Line" — extend both bullets.**

Add to the latest-line bullet:
```
After the job verifies the pre0 run, the same job stamps the manifests and the
FO pin on `main` to the pre0 version. An edge install then passes the FO version
gate with no human step.
```
Add to the old-line bullet:
```
The same decision also gates the `main` stamp. An old-line tag cannot move
`main` backwards.
```

**6. New section "Advancing the Stable Channel"**, after "Advancing the Edge Line".

After:
```
## Advancing the Stable Channel

The `stable` branch is the source that the stable marketplace entry resolves. A
release run advances the branch only when the tag is newer than the release that
`stable` serves now.

- The run fetches `stable` and reads the version in its
  `.claude-plugin/plugin.json`.
- If the tag is not newer, the run writes a `::notice::` and makes no change.
  The stable channel keeps the newer release.
- If the tag is newer, the run pushes the tagged commit to `stable` with
  `--force-with-lease`. The lease holds the SHA that the run read. If `stable`
  moves during the run, the push stops.
- A patch line leaves the history of `main`. The lease is what lets the next
  latest-line release move `stable` back onto the history of `main`.
- On the first stable release the `stable` ref does not exist yet. The run
  creates the ref with a plain push.
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
