---
title: "Edge channel survives the stable-cut window (no binary/skills minor skew between a stable tag and the first prerelease)"
status: ideation
source: "Captain report, 0250 Commander session 2026-07-07: edge tap installs the 0.24.0 binary while the spacedock-edge marketplace serves next-branch skills stamped 0.25.0-pre1. Verified: spacedock@next cask version 0.24.0 (updates only on tag push); origin/next .claude-plugin/plugin.json version 0.25.0-pre1; origin/next shared-core gate line 'require binary minor 0.25'. Result: every edge boot since the 0.24.0 stable cut (Jul 4) aborts at the FO version gate (binary too old, minor 0.24 < required 0.25). Broken by the release flow's own design — 'Advancing the Edge Line' bumps next to the post-release dev pre-version at the stable tag, but the edge BINARY only updates on the next tag push — not by any interim merge. Immediate remediation (separate from this task): push the first prerelease tag to realign."
started: 2026-07-07T08:22:09Z
completed:
verdict:
score: 0.5
worktree:
issue:
id: zr2rbsjsak7xx6tetr3n37hc
---

The release flow guarantees a broken edge channel from every stable cut until the first subsequent prerelease tag: the stable tag bumps next's manifests and contract gate line to the post-release pre-version while the spacedock@next cask keeps serving the stable-tag edge build, so the shipped version gate (same major.minor required) aborts every edge boot in the window. Direction options for ideation: (a) defer next's version bump so it rides the first prerelease tag instead of the stable tag; (b) have the stable tag's pipeline also cut a post-release pre-version edge build + cask bump in the same run, closing the window to minutes; (c) make the version gate tolerate skills exactly one minor ahead when the binary channel is edge — weakest, contract-side complexity. Acceptance sketch: value — after a stable cut, an edge-channel install boots green with zero manual steps (baseline: the 2026-07-04..07 window, every edge boot aborting); mechanism — the chosen pipeline change ships. High-stakes surface (CI and release machinery): detached adversarial audit + the release-flow dry-run treatment per docs/releasing.md.

## Root cause localization (FO forensics, 2026-07-07)

- `.github/workflows/release.yml:322-339`, job `edge-advance`, step "Reconcile the edge line past the stable release" (`if: "!contains(github.ref, '-')"`): computes `spacedock-release dev-preversion X.Y.Z` -> `X.(Y+1).0-pre1` and `stamp-version`s it into `.claude-plugin/plugin.json`, `.codex-plugin/plugin.json`, AND `skills/first-officer/references/first-officer-shared-core.md` (the `FO_PROSE` var) — rewriting the boot version-gate line to a binary minor that does not exist anywhere yet.
- `release.yml:345-353`: the marketplace calendar bump runs on both paths, so on a stable tag it actively triggers every edge installer's re-pull of the incompatible skills.
- The binary half (goreleaser edge build + `spacedock@next` cask bump) only publishes the TAG's version (0.24.0) and does not move again until the next tag push — no 0.25-line binary exists in the same run.
- Latent interaction, not a single bad line: under the retired contract-integer gate, one-minor-ahead skills booted fine (contract 3 == 3). #468 (`511dae11`) replaced it with minor-version coupling BEFORE the 0.24.0 cut; v0.24.0 (Jul 4) was the first stable tag under the new gate, converting the deliberate stamp-ahead ("never masquerade as the shipped stable") into a channel-wide boot abort. docs/releasing.md's own 0.24.0-pre1 note records the drift-cousin of this class.

## Design analysis addendum (captain question + FO analysis, 2026-07-07)

Captain question: should the stable cut always also cut an N+1.0-pre0 edge build, and what happens on a stable .1 patch release?

- **Always-cut-pre0 at the stable tag: sound, with machinery.** Content is exactly the audited stable bits re-stamped (next was just reconciled to the release commit), so nothing unaudited ships and the broken window closes to minutes. Needs an auto-tag (vX.(Y+1).0-pre0) on the stamped next commit + a recursion guard so the auto-tag does not re-trigger edge-advance. Pre-numbering is cosmetic: the boot gate requires same MINOR only (prerelease skew allowed), so a pre0 binary under pre1-stamped skills passes.
- **Maintenance patch releases expose a SECOND latent bug, independent of the original gap.** edge-advance assumes monotonic tag order; a vX.Y.1 cut from an older release line while next is on the (Y+1) dev line breaks three ways: (1) the reconcile's `merge -X theirs` FAVORS THE RELEASE, clobbering next's newer files with old patch-line content on conflict; (2) the stamp writes dev-preversion(X.Y.1) = X.(Y+1).0-pre1, REWINDING next's manifests/gate line if next has advanced past it; (3) under always-cut-pre0, the patch would attempt a colliding vX.(Y+1).0-pre0 auto-tag.
- **Fix for both: a line-ordering guard on edge-advance** — fire only when dev-preversion(tag) >= next's current manifest version; otherwise SKIP with a log line. An old-line patch then updates only the stable cask; its fix reaches edge via normal main->next flow, not the release reconcile. With the guard, always-cut-pre0 is the recommended shape for latest-line stable cuts.

## Ideation (2026-07-07)

### FLAG FOR THE IMMINENT v0.25.0 CUT (read before cutting)

This fix will almost certainly NOT ship before v0.25.0. Do NOT block the cut on it. Consequences the cut gate must plan around:

1. **v0.25.0 re-opens the window.** The cut stamps next to `0.26.0-pre1` and requires binary minor `0.26`, while the `spacedock@next` cask still serves the `0.25.0` edge build → every edge boot aborts until a `0.26`-line prerelease binary exists. The cutter MUST run the immediate remediation the entity source records: right after the stable cut, push a `v0.26.0-pre1` prerelease tag **on a greened commit** (green Runtime Live E2E for its SHA, else the e2e-gate blocks it) to build the `0.26`-line edge binary + bump the edge cask, realigning the channel. This is manual until always-cut-pre0 lands.
2. **Nothing about v0.25.0 makes this fix harder.** The fix's tests pin version algebra, not the `0.25.0` numbers; the cut and the fix do not interact.
3. **When always-cut-pre0 later ships, expect TWO GitHub releases per stable cut** (the `vX.Y.0` stable + the auto `vX.(Y+1).0-pre0` prerelease on the same commit) — by design, not a double-cut bug.

### Problem statement

The release flow deterministically breaks the edge channel from every latest-line stable cut until the first subsequent prerelease tag (root cause + the retired-integer-gate interaction are in the addenda above). Two independent defects:

- **The window.** `edge-advance` stamps next's manifests + FO gate line PAST the release to `X.(Y+1).0-pre1` at the stable tag, but the edge BINARY (`spacedock@next` cask) only publishes the tag's own `X.Y.0` build — no `X.(Y+1)`-minor binary exists until a later prerelease tag. Skills lead the binary by one minor; the same-minor boot gate aborts.
- **Old-line patch clobber.** `edge-advance` assumes monotonic tag order. A `vX.Y.1` patch cut from an older line while next is on the `(Y+1)` dev line reconciles `-X theirs` (favoring the OLD release) into next and stamps `dev-preversion(X.Y.1)=X.(Y+1).0-pre1`, clobbering next's newer content and rewinding its manifest/gate line.

Keeping the stamp-ahead is deliberate (rejected alternative "defer the bump / option (a)": leaving next's gate at the stable minor turns a LOUD boot-abort into a SILENT incompatibility once next accrues `(Y+1)`-binary-dependent skills — minor-coupling exists to prevent exactly that). So the window must be closed by shipping a real `X.(Y+1)` edge binary, not by lowering the gate.

### Proposed approach — two mechanisms plus three design corrections found during ideation

**Mechanism A — line-ordering guard on BOTH reconcile steps (the load-bearing fix).** A new `spacedock-release edge-advance-decision <tag> <next-plugin.json>` subcommand computes the tag's *target edge version* — `dev-preversion(tag)` for a bare stable tag, the tag's own version for a `-pre` tag (what next would inherit) — reads next's current manifest version, and prints `advance` or `skip`. Each reconcile step runs it right after `git switch -c edge-advance origin/next` (so the manifest read is origin/next's) and, on `skip`, emits `::notice::` and `exit 0` BEFORE any merge/stamp/push. This one guard covers the stable-cut window's pre0 recursion AND all three patch failure modes.

**Mechanism B — always-cut-pre0 at latest-line stable tags (closes the window).** After the stable-path reconcile+stamp+push, a new step auto-creates `vX.(Y+1).0-pre0` **on the greened stable RELEASE_COMMIT** and pushes it via the PAT (`HOMEBREW_TAP_TOKEN`-class token that re-triggers workflows). Its release run reuses the existing prerelease path end-to-end: `e2e-gate` resolves `git rev-list -1 vX.(Y+1).0-pre0` to the stable commit's SHA → matches the SAME green run that cleared the stable cut (no waiver); goreleaser builds+publishes the edge prerelease and bumps `spacedock@next` (`skip_upload: false` already publishes the edge cask on prereleases) to a binary stamped minor `X.(Y+1)` → window closes; the prerelease reconcile's line-ordering guard skips (`pre0 < next`'s `pre1`) → no rewind, no re-tag; recursion terminates at one level (auto-tag is stable-path-only, and its `-pre0` output routes to the prerelease path).

**Correction 1 — boundary is strict `>`, not `>=` (the addendum's `>=` is wrong).** Every legitimate forward stable cut yields `dev-preversion(tag)` strictly greater than next's manifest (a forward cut of `X.(Y+1).0` when next is `X.(Y+1).0-pre1` gives `X.(Y+2).0-pre1 > X.(Y+1).0-pre1`). Equality `dev-preversion(tag) == next` ONLY ever arises for a patch release (e.g. next `0.26.0-pre1`, tag `v0.25.1` → `dev-preversion=0.26.0-pre1`), which must SKIP. `>=` would FIRE on that equality and clobber next via `-X theirs`. Guard MUST use strict `>`.

**Correction 2 — auto-tag goes on the greened stable RELEASE_COMMIT, not the "stamped next commit" (the addendum's placement fails the e2e-gate).** SPIKE PROVED (below): the addendum's "auto-tag on the stamped next commit" places the pre0 tag on an UNGREENED SHA → run 2's `e2e-gate` finds no green Runtime Live E2E run for it → goreleaser (`needs: e2e-gate`) never runs → no edge binary → window stays open. Tagging the stable RELEASE_COMMIT reuses its existing green run. The pre0 build is then the greened stable tree relabeled — exactly the audited bits, nothing unaudited ships (satisfies the addendum's re-stamp intent by construction).

**Correction 3 — goreleaser needs `GORELEASER_CURRENT_TAG: ${{ github.ref_name }}` (load-bearing; default is silently wrong).** SPIKE PROVED: `git describe --tags` on the dual-tagged (`v0.25.0` + `v0.26.0-pre0`) commit returns `v0.25.0` in BOTH tag-creation orders, so goreleaser's default version source would stamp the pre0 edge binary as minor `0.25` and the window would NOT close. Setting `GORELEASER_CURRENT_TAG` to the triggering ref pins the edge binary to `0.26.0-pre0`. Safe on run 1 (only `v0.25.0` present at goreleaser time; the pre0 auto-tag is created later by `needs: goreleaser` edge-advance) and hardens all existing single-tag runs.

Guard comparison uses a new prerelease-aware `ComparePreVersion` in `internal/release` (the existing `contract.semverCompare` is dotted-int only and fails to parse `-preN`, so it cannot order `pre0 < pre1` — the exact distinction the recursion-skip needs).

### Acceptance criteria (each with its test)

- **AC-1 (VALUE, moves the recorded baseline right).** After a modeled latest-line stable cut `vX.Y.0`, the edge binary's stamped minor equals the edge skills' required-binary-minor, so the FO version gate passes with zero manual steps. Baseline that can move the wrong way: the 2026-07-04..07 window, where the post-0.24.0 edge state is binary minor `0.24` vs required `0.25` (mismatch → every edge boot aborts). Target: modeled post-cut state is binary minor `X.(Y+1)` == required `X.(Y+1)`.
  - **Test:** a Go fixture in `internal/release` reconstructs the stable-cut edge-advance + always-cut-pre0 sequence in a temp repo and derives (a) the edge binary version the pre0 run would stamp — from the triggering tag under the `GORELEASER_CURRENT_TAG` rule the fixture READS from the on-disk `.goreleaser.yaml`/`release.yml` env — and (b) the required-minor literal in next's FO shared-core prose after the stable stamp, then asserts equal major.minor. Adversarial twin drops `GORELEASER_CURRENT_TAG` (falls back to git-describe = the stable tag) and asserts the minors DIVERGE (baseline reproduced), so green is not vacuous. The real goreleaser binary build is the CI dry-run spike, not this unit fixture.

- **AC-2 (guard fires forward, skips backward).** `edge-advance` advances next only when the tag's target edge version is strictly greater than next's current manifest version; an older-line tag skips the whole reconcile (no merge, no stamp, no reconciled push) and logs a `::notice::`.
  - **Test:** `edge-advance-decision` table unit test — forward stable (`next 0.26.0-pre1`, `v0.26.0`→advance), old-line patch next-ahead (`0.27.0-pre1`, `v0.25.1`→skip), patch on just-released line (`0.26.0-pre1`, `v0.25.1`→skip), pre0-vs-next (`0.26.0-pre1`, `v0.26.0-pre0`→skip), major bump (`0.27.0-pre1`, `v1.0.0`→advance); plus `ComparePreVersion` ordering cases (`0.26.0-pre1 < 0.26.0`, `pre1 < pre2`, core dominates). A release.yml structure guard asserts both reconcile steps invoke the decision and `exit 0` on `skip` before merging; adversarial twin removes the guard call and reds.

- **AC-3 (three patch failure modes cannot regress next).** For `vX.Y.1` cut from an older line while next is on the `(Y+1)` line: (a) the reconcile is skipped so no `-X theirs` merge clobbers next's newer content; (b) next's manifest/gate-line version is not rewound; (c) no colliding `vX.(Y+1).0-pre0` auto-tag is attempted.
  - **Test:** an `edge_reconcile`-style fixture puts next on `0.27`-dev with exclusive content, cuts `v0.25.1`, runs the guarded edge-advance, and asserts next's tip is byte-identical (no new commit; content, manifest, FO-prose minor all unchanged) and no pre0 tag was created. Adversarial: with the guard removed, the same fixture shows `-X theirs` clobbers next and the stamp rewinds the manifest — proving the guard is load-bearing.

- **AC-4 (always-cut-pre0 closes the window, terminates, ships only audited bits).** A latest-line stable cut auto-tags `vX.(Y+1).0-pre0` on the greened RELEASE_COMMIT and pushes it via the re-triggering PAT; that run publishes an edge prerelease + bumps `spacedock@next` to a binary stamped minor `X.(Y+1)`; its edge-advance takes the prerelease path whose guard skips (`pre0 < pre1`) → no rewind, no re-tag; recursion ends at one level; the pre0 commit's tree == the greened stable tree.
  - **Test:** release.yml structure guard — the always-cut-pre0 step exists ONLY on the `!contains(github.ref,'-')` path, tags `RELEASE_COMMIT` (asserted to be the e2e-gate-greened SHA source, not the next tip), pushes via the PAT env, and the goreleaser step carries `GORELEASER_CURRENT_TAG: ${{ github.ref_name }}`; adversarial twins — (a) auto-tag on the next tip reds an "e2e-gate-reachable SHA" guard, (b) drop `GORELEASER_CURRENT_TAG` reds AC-1's divergence twin, (c) auto-tag on the prerelease path reds a recursion guard. A fixture confirms the pre0 tag routes to the prerelease reconcile → guard-skips → creates no further tag (termination), and that the pre0 commit SHA equals RELEASE_COMMIT (content == greened tree by construction).

### Spike record (riskiest mechanism first — CI/release is high-stakes)

- **DONE (local, this session): tag-collision / version resolution.** `git describe --tags` on a commit carrying both `v0.25.0` and `v0.26.0-pre0` returns `v0.25.0` in BOTH tag-creation orders. Consequence: default goreleaser version detection mis-stamps the pre0 edge binary to minor `0.25` and the window does NOT close → **`GORELEASER_CURRENT_TAG` override is mandatory, not optional** (Correction 3). This invalidated the naive design and is the seed for AC-1's adversarial twin.
- **TO SPIKE in implementation (named dry-run treatment).** (1) goreleaser actually building the edge binary at `vX.(Y+1).0-pre0` with `GORELEASER_CURRENT_TAG` set, on the dual-tagged commit, stamping `cli.Version=0.26.0-pre0` — via `goreleaser release --skip=publish,announce` (or `goreleaser build --single-target`) in a scratch clone, asserting the built binary's `--version` line reports `0.26.0-pre0`; run under the docs/releasing.md release-flow dry-run treatment. (2) `e2e-gate` SHA-reuse: `git rev-list -1 vX.(Y+1).0-pre0` == the greened stable SHA, confirmed against a scratch tag. (3) The auto-tag pushed via the re-triggering PAT actually fires a second release.yml run (vs GITHUB_TOKEN, which would NOT). A **detached adversarial audit** of the release.yml diff precedes the ideation/implementation gate, per the high-stakes-CI-surface policy.

### Test plan summary

Pure-Go unit + fixture tests in `internal/release` (edge-advance-decision table, `ComparePreVersion` ordering, AC-1 minor-equality fixture + divergence twin, AC-3 no-rewind fixture, release.yml structure guards each with an adversarial twin) — cheap, no network, seconds. The goreleaser binary build, cask bump, and e2e-gate SHA-reuse are the CI-level dry-run spike (minutes in a scratch clone), paid before implementation lands. No live release is cut to prove this. Estimated complexity: moderate — one new subcommand + one small compare fn + one release.yml step + guard calls in two existing steps + goreleaser env line, all backed by extensions to the existing `edge_reconcile_test.go` / `workflow_exec_guard_test.go` harness.

### Doc diff — docs/releasing.md "Advancing the Edge Line (`next`)"

Replace the section's **Stable (`vX.Y.Z`) tag** bullet and add two paragraphs:

Before (stable bullet):
> - **Stable (`vX.Y.Z`) tag:** `next` is reconciled the same way, then stamped PAST the release to the post-release dev pre-version (`spacedock-release dev-preversion X.Y.Z` → `X.(Y+1).0-pre1`), so the edge line never masquerades as the stable version it just shipped, then the calendar key is bumped.

After (stable bullet + line-ordering guard + always-cut-pre0):
> - **Stable (`vX.Y.Z`) tag, latest line:** `next` is reconciled the same way, stamped PAST the release to `X.(Y+1).0-pre1`, the calendar key is bumped, and then a `vX.(Y+1).0-pre0` tag is auto-created **on the greened release commit** and pushed (via the re-triggering tap PAT). That prerelease tag's own release run reuses the greened commit's e2e-gate pass, builds+publishes the `X.(Y+1)`-minor edge binary, and bumps the `spacedock@next` cask — so the edge binary's minor catches up to the skills' gate line within minutes instead of waiting for the next hand-cut prerelease. Expect two GitHub releases per stable cut.
> - **Old-line / patch (`vX.Y.1`, or any tag whose target edge version is not strictly greater than `next`'s current manifest version):** `edge-advance` SKIPS the reconcile entirely (`spacedock-release edge-advance-decision` prints `skip`, logged as a `::notice::`). The patch updates only the stable cask; its fix reaches edge through the normal `main`→`next` flow, never through a `-X theirs` reconcile that would clobber `next`'s newer `(Y+1)`-line content or rewind its manifest/gate line. The auto-pre0 step is stable-latest-line-only, so a patch never attempts a colliding pre0 tag.
>
> The `vX.(Y+1).0-pre0` release run does not recurse: its `-pre0` tag routes to the prerelease reconcile, whose line-ordering guard skips (`pre0 < next`'s `pre1`), and the auto-pre0 step runs only on the stable path — so it neither re-tags nor rewinds `next`. goreleaser's version for the edge binary is pinned by `GORELEASER_CURRENT_TAG` to the triggering ref, because `git describe` on the dual-tagged release commit resolves to the stable tag, not the pre0 tag.

## Stage Report: ideation

- DONE: Design per the entity's recorded addenda — line-ordering guard as load-bearing fix + always-cut-pre0 + value AC moving the baseline right
  Body "## Ideation (2026-07-07)" AC-1 (value: post-cut edge binary minor == required minor, vs the 2026-07-04..07 abort baseline) through AC-4; guard covers the window's pre0 recursion AND all three patch modes.
- DONE: Spike-or-record the riskiest mechanism first; name the dry-run treatment + detached audit
  Local spike run: `git describe --tags` on a dual-tagged (`v0.25.0`+`v0.26.0-pre0`) commit returns `v0.25.0` both orders → forced Correction 3 (`GORELEASER_CURRENT_TAG` mandatory). Remaining goreleaser-build / e2e-gate-SHA-reuse spikes named under docs/releasing.md dry-run treatment + a detached release.yml audit.
- DONE: Each of the three patch failure modes gets a named proof + concrete doc diff for the Advancing-the-Edge-Line section
  AC-2/AC-3 name the proofs (old-line skip; no `-X theirs` clobber; no manifest/gate rewind), each with an adversarial twin; doc diff rewrites the stable bullet and adds the old-line/patch + recursion paragraphs.
- DONE: Surface anything the imminent v0.25.0 cut should do differently/avoid
  Prominent "FLAG FOR THE IMMINENT v0.25.0 CUT" callout: don't block the cut; expect the window to re-open; apply the manual `v0.26.0-pre1`-on-a-greened-commit remediation right after.

### Summary

Designed the fix from the recorded root-cause + design-analysis addenda, and the local tag-collision spike forced three corrections to the recorded design: the guard boundary is strict `>` not `>=` (equality only ever hits patch releases, which must skip); the pre0 auto-tag must sit on the greened stable RELEASE_COMMIT (the addendum's "next commit" placement is ungreened → e2e-gate blocks the pre0 build → window never closes); and goreleaser needs `GORELEASER_CURRENT_TAG` because `git describe` on the dual-tagged commit resolves to the stable tag. The value AC (AC-1) measures post-cut binary-minor == required-minor against the recorded abort baseline via a unit fixture with a divergence twin; the goreleaser binary build itself stays a named CI dry-run spike. Flagged prominently that v0.25.0 will re-open the window and needs the manual prerelease realign, since this fix won't land first.
