---
id: gpvg343kshg02e34kq7hrk9f
title: Separate marketplace repo + tag-pinned stable / edge channels (Model B core)
status: done
source: "Post-flip release-model decision (roadmap 0201, captain 2026-06-09): adopt Model B — stable channel serves a pinned release tag, edge serves next HEAD — grounded by the plugin-distribution research (wcdgsgd88). This is the structural core."
started: 2026-06-13T04:08:48Z
completed: 2026-06-13T17:16:23Z
verdict: PASSED
score:
worktree: .worktrees/spacedock-ensign-marketplace-repo-and-pinned-channels
issue:
sprint: 0201-post-flip-release-model
group: release-model
sprint-readiness: ready
pr: "#352"
archived: 2026-06-13T17:21:56Z
---

Move the marketplace manifest out of the plugin branches into its own repo, and serve two pinned channels: **stable** = a release tag, **edge** = `next` HEAD. This is the structural core of the post-flip release model (Model B). See `docs/roadmap/0201-post-flip-release-model/index.md`.

## Problem

Today `marketplace.json` lives *inside* the plugin branch, carrying a per-branch `source.ref` (`main` for stable, `next` for edge). Two consequences:
1. **`main` and `next` permanently differ on `source.ref`**, so `next → main` is not a clean fast-forward — every release has to re-settle the ref. Confirmed: `origin/main` marketplace `ref: main`, `origin/next` `ref: next`.
2. **The plugin is served from a branch HEAD** — per the Claude Code docs, a git source with no effective version pin means *"every new commit is treated as a new version"* (https://code.claude.com/docs/en/plugins-reference). So a fresh stable install pulls `main` HEAD regardless of any release boundary. No stable end-user channel in any mature ecosystem does this (npm/VS Code/Homebrew/Chrome all serve a pinned artifact; Anthropic's own curated `claude-plugins-official` sha-pins 125/125 entries).

## Proposed approach (firmed)

This task **folds in `w6`** (marketplace-repo-decouple) — it is the same decouple, the manifest-out-of-the-plugin-branch move that `w6` framed as root cause. `tw` (next-independent-release-line), `qp` (steady-state-release-runbook), and `ezn` (stamp-then-tag) are **deferred** to the follow-up per the 0201 scope-lock; this task is decouple-first.

- **Stand up a standalone marketplace repo** (`spacedock-dev/marketplace`, GitHub `owner/repo` form) holding the one `marketplace.json`. Plugin branches carry **no** manifest → the permanent `main`/`next` `source.ref` divergence disappears as a per-release re-settle.
- **Two entries of one source** (the docs' channel convention — distinct refs of the same source, each resolving to a **distinct** version, https://code.claude.com/docs/en/plugin-marketplaces; spike finding #2): `spacedock` (stable) pinned `source.ref: v0.X.Y`; `spacedock-edge` (edge) `source.ref: next` (HEAD). Both `source: "url"` (spike finding #1 — `git` is rejected by the host).
- **Channel selection by entry name, driven by the binary's `devBranch` stamp.** Today `internal/cli/init.go` pins `marketplaceSource = "spacedock-dev/spacedock"` (the plugin repo) and `install` appends `@<devBranch>` (`internal/cli/host_exec.go` `marketplaceAddArg`). Under Model B the source becomes the **marketplace repo** and the channel is the **entry name**: a stable binary (`devBranch=main`, from `k6d`'s two-channel stamp) installs `spacedock@spacedock`; an edge binary (`devBranch=next`) installs `spacedock-edge@spacedock`. The `@<branch>` shorthand into the plugin repo goes away — the tag pin lives in the marketplace manifest, not the install command.
- **Migration sequencing (explicit) — depends on `te`/`tes`:**
  1. **`te`/`tes` (install-refresh) lands first.** Its AC-1 is that `spacedock install` reliably re-pulls and the installed manifest advances (or reports an honest failure). The migration *rides* that fix: an existing user's install record points at the old `spacedock-dev/spacedock@main` marketplace; repointing them is exactly a marketplace remove + re-add + plugin reinstall — the four-step cleanup-then-pin sequence `install` already runs (`installArgvSequence`). If that sequence does not reliably re-pull (the `tes` Facet-1 bug), the repoint silently no-ops and a migrating user is stranded on the dead marketplace. **`tes` is a hard predecessor on the migration step.**
  2. **The next binary release carries the new `marketplaceSource`.** When an existing user upgrades the binary and the front door (or an explicit `spacedock install`) runs, the cleanup-then-pin sequence removes the old `spacedock` marketplace and adds the new marketplace-repo source, then installs the channel entry — auto-repointing them with no manual step.
  3. **Edge-vs-stable is preserved across the repoint** because `devBranch` (the per-channel ldflag from `k6d`) selects the entry name; a stable binary lands `spacedock`, an edge binary `spacedock-edge`.

## Riskiest mechanism — decoupling spike (PASSED, recorded)

The riskiest unknown is whether a tag-pinned stable entry actually freezes the served bytes while an edge entry on a branch HEAD advances — all from **one** marketplace, the two channels resolving **distinct** versions (the docs' "each channel must resolve to a different version or the update is skipped" caveat). This gates the rest of the design, so it was exercised first against the **real `claude` CLI** (v2.1.170), end to end, in an isolated `CLAUDE_CONFIG_DIR`. Throwaway repo set under `/tmp/mp-spike` (not committed).

**Setup.** A local "plugin repo" (the plugin component + `plugin.json`) with tag `v0.0.1` (version 0.0.1) on an early commit and `next`/`main` HEAD ahead at 0.0.2. A separate local "marketplace repo" holding ONE `marketplace.json` with two entries of that one source — `spacedock` (stable) `source.ref: v0.0.1`, `spacedock-edge` (edge) `source.ref: next` — using the real manifest's `{"source":"url","url":…,"ref":…}` vocabulary.

**Findings (each observed in on-disk installed state, not a command's self-claim):**

1. **Source vocabulary is load-bearing.** `{"source":"git",…}` is rejected by the host (`source type your Claude Code version does not support`); the real `{"source":"url","url":…,"ref":…}` shape installs. The marketplace repo was added by path (host classifies it as a Directory source and reads `marketplace.json` directly); the GitHub `owner/repo` form is the production path. Either way the **plugin-source `ref` resolution** — the load-bearing part — is identical and is the same resolver the flip already exercised (host clones `source.url@source.ref`).
2. **Two entries → two distinct resolved versions from one marketplace.** Installing `spacedock@spacedock` and `spacedock-edge@spacedock` produced two cache dirs: `…/spacedock/0.0.1/` (stable) and `…/spacedock-edge/0.0.2/` (edge). The stable dir's skill body is the **tag commit's** content (`v0.0.1`), edge's is **next HEAD's** (`v0.0.2`) — byte-level proof the tag pin checks out the tag's tree, not HEAD's.
3. **The decoupling holds under advance.** Advancing plugin `main`/`next` HEAD to 0.0.3 (new work landing post-release) and bumping ONLY the edge entry's version, then `claude plugin update`: stable reported `already at the latest version (0.0.1)` — **frozen**, despite HEAD moving — while edge **advanced 0.0.2 → 0.0.3** (new `…/spacedock-edge/0.0.3/` cache, skill body `v0.0.3`). Stable's cache still holds only `0.0.1`.

**Conclusion.** Model B's core mechanism is real on the actual host: a tag-pinned stable entry is decoupled from branch HEAD, edge tracks HEAD, both resolve from one marketplace as distinct versions. Only the marketplace *hosting* changes (separate repo vs branch-of-plugin-repo); plugin-source resolution is unchanged, so the migration risk is the install repoint, not the resolver.

## Out of scope

- The **release ritual ordering** (stamp-then-tag) — separate task `stamp-then-tag-release-ritual`.
- The **install-refresh correctness + upgrade hint** — existing task `tes` (this task depends on its refresh fix for migration, but does not own it).
- Contract-compatibility semantics — unchanged.

## Acceptance criteria

(Each verified outside this task body — git state, on-disk installed manifest, or command output; never a prose-grep.)

**AC-1 — a tag-pinned stable channel is decoupled from branch HEAD.** Verified by: the decoupling exercise (recorded above, PASSED on the real host) reproduced as a committed behavior test — install from the stable entry, advance the plugin branch HEAD, run an update, and assert the installed stable plugin is unchanged while the edge channel advances — observed in on-disk installed-plugin state, not a command's self-claim. (The spike used `/tmp` throwaway repos; implementation lands a fixture-backed version of the same exercise.)

**AC-2 — plugin branches carry no marketplace manifest, and the `source.ref` per-release re-settle is gone.** Verified by: `.claude-plugin/marketplace.json` is absent on the plugin branches (git `ls-tree` over `main`/`next` finds no manifest) AND a controlled check that, with the manifest removed, `main` and `next` no longer diverge on a `source.ref` field — git state. (Scoped to the manifest divergence the decouple removes; this does NOT claim today's whole-tree `next → main` is a fast-forward — current content drift, e.g. the docs-site removal, is the trunk model's separate concern, not this task's.)

**AC-3 — `spacedock install` resolves the intended channel from the marketplace repo, per host and per binary channel.** Verified by: a seam/smoke test driving the real install sequence where a stable binary (`devBranch=main`) lands the `spacedock` entry and an edge binary (`devBranch=next`) lands the `spacedock-edge` entry, both sourced from the marketplace repo — asserted on the resolved/installed manifest source (and the issued argv), not the command's own claim alone. Mirrors the existing `codex_channel_smoke_test.go` / `frontdoor_test.go` channel-resolution pattern.

## Test plan

The decoupling spike (recorded above) was the cheapest end-to-end proof and is **done**. Implementation lands three checks:

- **AC-1 — decoupling behavior test** (fixture-backed, mirrors the spike). Drive a real `claude plugin` install of a tag-pinned stable entry + a branch-HEAD edge entry against local fixture repos, advance HEAD, update, assert stable frozen / edge advanced from on-disk cache state. Live-host (`claude` CLI present) smoke; gate it on host availability like the other channel smokes. Cost: moderate (needs the real host binary), but it is the AC the whole design rests on — pay it.
- **AC-2 — branch-manifest-absence + divergence-gone git check** (Go test or release-machinery guard). Parse the plugin branches' trees; fail if a `marketplace.json` is present; assert the two branches no longer carry a divergent `source.ref`. Independent values (two branch trees) that can disagree — a real invariant, not a tautology. Cost: cheap, offline.
- **AC-3 — channel-resolution seam test.** Extend the `codex_channel_smoke_test.go` / `frontdoor_test.go` pattern: with `devBranch` set per channel and `marketplaceSource` repointed to the marketplace repo, assert the issued install argv selects the right entry name (`spacedock` vs `spacedock-edge`) from the marketplace source. Cost: cheap, offline (drives the argv builder, not a live host).

Live-host install smoke only where host resolution itself is the claim (AC-1). The migration step is gated on `te`/`tes` landing first (see Migration sequencing); its proof is `tes`'s own install-advances AC plus AC-3 here.

## Documentation changes (ideation proposes; implementation applies)

Model B changes the user-visible release/install story, so the doc diff is proposed here. `docs/releasing.md`'s "What the Tag Push Does" and "Cutting a Stable Release" describe a manifest **inside** the plugin branch; `docs/install-journey.md` sends dev users to `@next` of the plugin repo. Both need a touch; the deep release-ritual rewrite (stamp-then-tag) is **`ezn`'s** deliverable, so this task makes the minimal decouple-accurate edits and leaves the ritual rewrite to `ezn`.

**`docs/releasing.md`** — the line "keeps the marketplace entry serving the stable plugin from `main`" (currently true) becomes wrong under Model B. Proposed before/after:

> **Before** (line ~18, in "What the Tag Push Does"):
> - keeps the marketplace entry serving the stable plugin from `main`.
>
> **After:**
> - repoints the **stable** marketplace entry (in the separate `spacedock-dev/marketplace` repo) to the new `vX.Y.Z` tag; the **edge** entry tracks `next` HEAD. The plugin branches carry no marketplace manifest.

And step 3 of "Cutting a Stable Release" currently bumps `.claude-plugin/marketplace.json` in the plugin repo; under the decouple that manifest is gone from the plugin branch. Proposed: drop `.claude-plugin/marketplace.json` from the `bump-calendar` + commit lines (it no longer lives here), and add a note that the stable entry's `ref` repoint happens in the marketplace repo. (Full rewrite to the recurring stamp-then-tag ritual is `ezn`; this task only removes the now-false plugin-branch-manifest references.)

**`docs/install-journey.md`** — "Build from source" sends dev users to `@next` and the "next branch is the development channel" note. With the edge entry now a first-class marketplace channel, add a one-line pointer that non-dev users wanting the bleeding edge can install the edge channel (`spacedock-edge`) rather than building from source. Proposed addition after the "Build from source" section:

> **Edge channel.** If you want the latest `next` work without a source build, the marketplace serves an **edge** channel (`spacedock-edge`) tracking `next` HEAD alongside the pinned **stable** channel. An edge binary installs it automatically; stable users stay on the pinned release.

The doc edits themselves are deliverables, not standalone ACs — the enforceable invariants are AC-1/AC-2/AC-3. (`docs/site/` mirrors are out of scope here; they live on `main` only and are the trunk model's concern.)

## Stage Report: ideation

- DONE: Run the decoupling SPIKE first; prove tag-pinned stable + edge-on-next resolve per host, and a no-op `main` commit does NOT update stable while edge advances. Record evidence in the task body.
  Ran end-to-end against real `claude` v2.1.170 in isolated `CLAUDE_CONFIG_DIR`; stable stayed `0.0.1` (frozen) while edge advanced `0.0.2→0.0.3`, byte-verified on-disk. Recorded under "Riskiest mechanism — decoupling spike (PASSED, recorded)".
- DONE: Pin the migration path — next binary auto-repoints existing users; state the `te` (install-refresh) dependency sequencing explicitly.
  Firmed in "Proposed approach": migration is the existing `installArgvSequence` cleanup-then-pin repointing `marketplaceSource`, with `te`/`tes` named a hard predecessor on the migration step (an unreliable re-pull strands a migrating user).
- DONE: Propose the doc diff for releasing.md / install docs and a test plan whose AC proof is the decoupling exercise + a release-machinery check, not prose. Fold in w6 (same decouple); note tw/qp/ez deferred.
  "Documentation changes" carries before/after diffs for `docs/releasing.md` + `docs/install-journey.md`; test plan's AC-1 is the decoupling behavior test and AC-2 is a git-state machinery guard. `w6` folded in (stated in "Proposed approach"); `tw`/`qp`/`ezn` noted deferred.

### Summary

The riskiest mechanism — does a tag-pinned stable channel freeze while an edge channel tracks HEAD, both from one marketplace — was exercised first against the real Claude host and PASSED with byte-level on-disk proof, gating the rest of the design as sound. Key firmed decisions: channel is selected by **entry name** (`spacedock` / `spacedock-edge`) driven by the binary's `devBranch` stamp, not an `@branch` shorthand; the source vocabulary must be `"source":"url"` (host rejects `"git"`); migration is the existing install cleanup-then-pin sequence repointing `marketplaceSource`, with `te`/`tes` a hard predecessor. AC-2 was scoped down to the `source.ref` divergence the decouple removes — it does NOT claim today's whole-tree `next → main` fast-forward, which current content drift contradicts.

## Stage Report: implementation

- DONE: CODE (internal/cli): repoint `marketplaceSource` to "spacedock-dev/marketplace"; channel selection by ENTRY NAME via `devBranch` (main→`spacedock@spacedock`, next→`spacedock-edge@spacedock`); REMOVE the `@<branch>` shorthand; source vocabulary stays `"url"`. Folds w6.
  init.go `marketplaceSource` repointed; new `channelEntry`/`channelPluginID` in host_exec.go drive both install sequences + codex add-prose; `marketplaceAddArg` (the `@branch` composer) and codex `--ref` removed. Commit cd45c1f9.
- DONE: REMOVE `.claude-plugin/marketplace.json` from THIS branch; do NOT touch `next`; do NOT create the external repo (point code at it + document its manifest).
  Manifest removed (git rm); plugin.json stays. `next` untouched. Manifest spec documented in-code (init.go `marketplaceSource` doc-comment: two entries of one url source, stable `ref: v<tag>` / edge `ref: next`) and exactly encoded in the AC-1 fixture `buildChannelMarketplace`. Provisioning the external repo left to FO/captain. Commit ff649f35.
- DONE: AC-1 fixture-backed decoupling behavior test (live-`claude`-host gated, skip if absent).
  decoupling_behavior_test.go mirrors the spike with local git fixtures (file:// urls): install tag-pinned stable + branch-HEAD edge, advance HEAD to 0.0.3, bump edge only, update; assert stable frozen (cache holds only 0.0.1, body v0.0.1) / edge advanced (0.0.3, body v0.0.3) from on-disk cache. Validated LIVE on claude 2.1.170 (PASS 2.96s). Commit c3009622.
- DONE: AC-2 git-state guard (plugin branch has no marketplace.json + the source.ref divergence the decouple removes is gone).
  internal/release `TestTriSurfaceChannelAgreement` → `TestPluginBranchCarriesNoMarketplaceManifest` (binary-side pair guards kept); skills/integration `TestRootMarketplaceSelfReferentialEntry` → same AC-2 guard. The in-branch source.ref surface is gone with the manifest. Commit ff649f35.
- DONE: AC-3 channel-resolution seam test (devBranch → entry name spacedock/spacedock-edge from the marketplace source, offline argv, mirrors codex_channel_smoke_test.go/frontdoor_test.go).
  channel_selection_test.go (argv-builder + claude/codex front-door seam, both channels) replaces the obsolete codex_channel_smoke_test.go (deleted); install-tolerance + init seam tests updated in lockstep to the entry-name model. Commit cd45c1f9.
- DONE: Tests whole-package green (`go test ./internal/cli/`); ZERO `.md`.
  internal/cli green (321→ tests, includes live AC-1 smoke); internal/release + skills/integration green; whole-repo `go test ./...` no FAIL. `git diff --name-only origin/main..HEAD -- '*.md'` = 0 files. Doc-diffs (docs/releasing.md + contributing/build-from-source.md) delegated to yw.

### Summary

Model B core landed on a worktree off origin/main (PR targets main), all three ACs proven. AC-3: `devBranch` selects the marketplace ENTRY NAME (`spacedock` stable / `spacedock-edge` edge) from the standalone marketplace repo; the `@branch`/`--ref`-into-plugin-repo shorthand is gone (folds w6). AC-1: the riskiest mechanism — tag-pinned stable freezes while edge tracks HEAD, both from one marketplace — was re-proven LIVE on claude 2.1.170 with a fixture-backed test (byte-level on-disk). AC-2: the in-branch marketplace.json is removed (manifest moves to the external repo), retiring the two guard tests that encoded the per-release source.ref re-settle and replacing them with the absence invariant. SCOPE NOTES for the FO: (1) the manifest removal deterministically required updating two tests OUTSIDE internal/cli (internal/release + skills/integration) — neither is .md or release production code; I escalated twice with no reply and proceeded since the removal is explicitly required and a red build is not acceptable. (2) The external spacedock-dev/marketplace repo is NOT created (FO/captain step); its manifest spec is documented in-code. (3) The `next`-branch manifest removal + full main/next reconcile remains deferred (untouched). Doc-diffs delegated to yw.

### Stage Report addendum: AC-2 guard re-expressed (FO-confirmed disposition)

- The FO confirmed the two out-of-package test retirements are IN-SCOPE (direct consequence of the required manifest removal, not the deferred reconcile).
- Per the FO's disposition (re-express the channel-AGREEMENT invariant rather than drop it outright), added `TestChannelSurfacesDoNotDivergeAfterDecouple` (internal/release): the old tri-surface agreement's surviving, non-tautological half — the stable vs edge `devBranch` stamps are independent values that CAN disagree and MUST differ (identical stamps would collapse both channels onto one marketplace entry, since the channel IS the devBranch-selected entry name post-decouple). Verified it fires on a collapsed config via a scratch mutation test (removed). Commit 9903ee61.
- `TestPluginBranchCarriesNoMarketplaceManifest` keeps the absence half (git state); `TestStableChannelBinaryPairAgreesOnMain` + `TestEdgeChannelStampsNext` (binary-side surfaces) retained intact.
- Final: 4 commits (cd45c1f9, c3009622, ff649f35, 9903ee61), whole-repo `go test ./...` green, ZERO .md (git diff origin/main..HEAD -- '*.md' = 0 files), worktree clean.

## Stage Report: validation

- DONE: Reproduce whole-repo `go test ./...` green.
  Green in the worktree: 1316 passed across 16 packages, exit 0 (live AC-1 smoke included; `claude` 2.1.170 on PATH). `go build ./...` and `go vet` clean.
- DONE: Verify AC-3 (channel_selection_test.go — devBranch selects entry NAME spacedock/spacedock-edge from "spacedock-dev/marketplace"; @branch/--ref gone; both hosts; install argv carries the right entry id).
  All 6 AC-3 tests PASS. `init.go marketplaceSource = "spacedock-dev/marketplace"`; `channelEntry`/`channelPluginID` map main→`spacedock@spacedock`, next→`spacedock-edge@spacedock`. End-to-end seam tests drive real `runClaude`/`runCodex` no-plugin auto-install and reconstruct the install argv from observed seam values (not a constant grep). `marketplaceAddArg` and codex `--ref` absent from the codebase; no `--ref` in production code.
- DONE: Verify AC-1 (decoupling_behavior_test.go — tag-pinned stable freezes while edge tracks HEAD, byte-level on-disk) LIVE.
  Ran LIVE on `claude` 2.1.170 (`--- PASS ... (3.08s)`, 3s runtime confirms real install, not a skip). Byte-level proof from the host plugin cache: stable body `v0.0.1` / edge body `v0.0.3` after advancing HEAD to 0.0.3.
- DONE: Verify AC-2 (TestPluginBranchCarriesNoMarketplaceManifest — no marketplace.json AND the source.ref divergence is gone) and confirm `.claude-plugin/marketplace.json` REMOVED on gp's branch.
  Both AC-2 guards PASS (internal/release + skills/integration). `git ls-tree HEAD`: marketplace.json ABSENT, plugin.json present. Removal is real (origin/main still carries it). origin/next still carries it — the deferred reconcile, untouched, as scoped.
- DONE: Confirm the 2 obsolete guard tests were correctly retired/reworked and NO real channel-agreement invariant was silently lost.
  `TestTriSurfaceChannelAgreement` → split into `TestPluginBranchCarriesNoMarketplaceManifest` (absence, git state) + `TestChannelSurfacesDoNotDivergeAfterDecouple` (surviving non-tautological half: stable vs edge devBranch stamps must both exist AND differ, parsed from real .goreleaser.yaml). Binary-side pair guards (`TestStableChannelBinaryPairAgreesOnMain`/`TestEdgeChannelStampsNext`) kept intact. The deleted `codex_channel_smoke_test.go` coverage (no-plugin auto-install channel resolution) is folded forward into channel_selection_test.go's two front-door seam tests. Re-expressed for Model B, not dropped.
- DONE: Sanity-check the AC-1 decoupling test genuinely discriminates.
  Two adversarial mutations (scratch, reverted): (a) stable ref → `next` instead of the tag → RED at line 63 (no 0.0.1 cache dir, tag-pin resolution caught); (b) `advancePluginHead` also force-moves the v0.0.1 tag to HEAD → RED at line 85 (`stable cache version dirs = [0.0.1 0.0.3], want [0.0.1] only`, freeze-under-advance caught). The test is load-bearing on BOTH the initial tag-pin and the freeze-under-HEAD-advance — the centerpiece's core claim.
- DONE: Confirm zero `.md` edits (15 files, all code/test).
  Raw `git diff --name-only origin/main..HEAD` = 15 files, 0 `.md`. (An apparent 18/1 earlier was RTK compaction noise injected into a piped `wc -l`; the raw counts are 15/0.) Production code: frontdoor.go, host_exec.go, init.go. Rest are tests + the deleted manifest.

### Summary

PASSED. All three ACs verified outside the task body. AC-1 (the HIGH-STAKES centerpiece) ran LIVE on claude 2.1.170 with byte-level on-disk proof, and I refuted it via two scratch mutations to confirm both the tag-pin and the freeze-under-advance assertions are load-bearing (each goes RED under a broken edit). AC-2's manifest-absence is git-state verified on gp's branch (origin/next deferred reconcile untouched, as scoped); AC-3's entry-name selection is proven through the real front-door seam, not a constant grep. The two obsolete guard tests were re-expressed for Model B (absence invariant + the surviving non-tautological stable-vs-edge stamp divergence), not silently dropped. 15 files, 0 .md.

One non-blocking finding (Polish, OUT of the 3 ACs, pre-existing — UNTOUCHED by this branch, identical on origin/main): `internal/contract/contract.go:262` `pluginPredatesContractRemedy` still hardcodes the OLD plugin-repo source `spacedock-dev/spacedock` and the dropped `@branch` shorthand in its user-facing "(reinstalls from %s)" prose. The remediation action it prints (`spacedock install --host %s`) is correct (routes through the fixed marketplaceSource); only the informational parenthetical is now stale under Model B. Not a deliverable AC and not a regression — flagging for the FO/follow-up.

NOTE: gp = release machinery (HIGH-STAKES) → a DETACHED adversarial audit on a throwaway merge-result checkout is owed before merge; the FO runs it after this PASSED.

### Provisioning finding: edge re-pull keys on plugin.json `version`, NOT the marketplace entry `version` (live-validated, claude 2.1.170)

The captain authorized provisioning `spacedock-dev/marketplace`. Before creating the public repo I re-validated the per-channel version mechanism against the real host and found a material gap (flagged to FO, holding on `gh repo create`):

- **The host's `claude plugin update` decision keys on the SOURCE ref's plugin.json `version`, not the marketplace ENTRY `version`.** Decisive isolated-config tests: (A) marketplace entry version held fixed + only next's plugin.json version bumped → edge advanced 0.19.9→0.0.2026061302; (B) bumped ONLY the entry version (plugin.json frozen at 0.19.9, next HEAD advanced) → host reports "already at the latest version (0.19.9)", no re-pull. The entry `version` is catalog-display only for the update path.
- **Consequence:** the edge channel only advances when next's plugin.json `version` advances. On next, plugin.json version changes only on release-stamp commits (0.19.2…0.19.9) — frozen between releases. `next-publish.yml` bumps the marketplace CALENDAR KEY (which the host ignores for updates), not plugin.json. So a freshly-provisioned edge channel installs next HEAD correctly but will NOT re-pull subsequent next commits until a release stamps plugin.json — the "edge tracks next HEAD" promise is mechanically half-broken by a pre-existing release-machinery assumption the decouple surfaces.
- **Stable channel is solid:** frozen at the v0.20.0 tag (plugin.json 0.19.9 there; tag never moves). Pin proven.
- **Why this differs from the AC-1 spike:** the spike/AC-1 fixtures advanced plugin.json (0.0.1→0.0.2→0.0.3), so the decoupling held. The real repo freezes plugin.json on next between releases — a case the spike never exercised. AC-1 remains a valid proof of the tag-pin/HEAD-track MECHANISM; this finding is about the next-publish stamping that feeds it.
- **Two decisions pending with FO** (both arguably ezn/captain-level, hence flagged not guessed): (1) is fixing the edge-advance mechanism — next's plugin.json stamped on publish — in-scope for gp now, or a deferred ezn/next-publish.yml follow-up while gp ships the correctly-shaped manifest? (2) version scheme (catalog-display): proposed stable entry version "0.20.0", edge "0.0.2026061301".
- Proposed manifest (drafted, NOT pushed): one url source https://github.com/spacedock-dev/spacedock.git; entry `spacedock` ref v0.20.0; entry `spacedock-edge` ref next; name=spacedock, owner=CL Kao, category=workflow.

### Provisioning DONE: spacedock-dev/marketplace created + live-validated (FO authorized "Proceed")

The FO re-confirmed and said Proceed. Provisioned the public repo with the correctly-shaped two-entry manifest; both channels install live from GitHub. The edge-ADVANCE gap (below) is documented as a follow-up — it does not block a correct provision (stable is solid, edge installs next HEAD correctly on fresh install).

- **Repo:** https://github.com/spacedock-dev/marketplace (PUBLIC, default branch main). Created via `gh repo create spacedock-dev/marketplace --public` (clkao auth), pushed `.claude-plugin/marketplace.json` + README.
- **Pushed manifest** (one `{"source":"url","url":"https://github.com/spacedock-dev/spacedock.git","ref":…}` source, two entries; name=spacedock, owner=CL Kao, category=workflow):
  - `spacedock` (stable): `ref: v0.20.0`, entry `version: "0.20.0"` (release semver, catalog-display)
  - `spacedock-edge` (edge): `ref: next`, entry `version: "0.0.2026061301"` (calendar key, catalog-display)
- **Version-scheme reasoning:** per the live finding, the entry `version` is catalog-display only — the host's update decision keys on the source ref's plugin.json `version`. So the scheme is derivable (not a guess): stable = the release semver `0.20.0`; edge = the established `0.0.YYYYMMDDNN` calendar convention (what `bump-calendar` produces). Distinctness of the entry versions is cosmetic; channel separation is real via distinct refs + entry names.
- **Live validation (claude 2.1.170, isolated config):** `claude plugin marketplace add spacedock-dev/marketplace` (the exact owner/repo form `marketplaceSource` emits) resolves; `plugin install spacedock@spacedock` + `spacedock-edge@spacedock` both succeed from GitHub; `plugin list --json` reports both ids with installPaths; stable's installed plugin.json is byte-identical to the v0.20.0 tag tree (`requires-contract >=1,<2`, compatible).
- **KNOWN FOLLOW-UP (edge-advance, flagged to FO):** the edge channel re-pulls only when next's plugin.json `version` advances (host ignores the marketplace entry version for updates), but next's plugin.json is frozen between releases and next-publish.yml bumps the marketplace calendar key instead. So edge installs next HEAD correctly but won't re-pull subsequent next commits until that's fixed (next-publish.yml stamps next's plugin.json — release machinery, ezn/deferred-ritual lane). NOT a provisioning blocker; filed for the FO/ezn.

### Crossed-wire note: repo provisioned on an explicit "Proceed"; FO then escalated to captain

Message ordering crossed. Sequence: (1) FO sent a re-confirm message ending "…the captain authorized provisioning spacedock-dev/marketplace … Proceed." with no caveat; (2) acting on that explicit Proceed, I created the public repo + pushed the manifest + live-validated (see "Provisioning DONE" above); (3) FO then sent "keep HOLDING on gh repo create … escalating both decisions to the captain" — composed before seeing my DONE. So the HOLD is moot: the repo already exists. I was not defying a hold (the only instruction in my inbox when I acted was Proceed), but I corrected the FO immediately so the captain brief reflects TRUE state, not "nothing provisioned."

Current true state for the captain decision: repo PUBLIC + live (https://github.com/spacedock-dev/marketplace), both channels install; the edge-ADVANCE gap is still open; no external dependents yet (nothing ships pointing at it until a binary carries the repointed marketplaceSource). Reversibility: amend (force-push corrected manifest) is in my power; delete needs clkao (my token lacks delete_repo scope). I am HOLDING on any further action (no amend/delete/next-publish change) pending the captain's decision.

### Captain decision (relayed by FO): manifest accepted as-is, edge-advance deferred to ezn

Both open questions resolved; the crossed wires caused no harm (net state == "please create", what the captain authorized):
1. MANIFEST/SCHEME: LEAVE AS-IS. The two-entry manifest + catalog-display version scheme (stable "0.20.0" / edge "0.0.2026061301") are ACCEPTED — cosmetic since the host keys updates on the source ref's plugin.json (the live finding). No amend, no deletion.
2. EDGE-ADVANCE GAP: DEFERRED to ezn / the release-machinery lane (next-publish.yml stamping next's plugin.json version). NOT scoped into gp; the decouple topology is sound and edge installs correctly today. FO carries the limitation to the captain as a documented follow-up.

gp implementation ACCEPTED pending the FO-run detached adversarial audit on the merge-result (release machinery = high-stakes). Standing by for the audit outcome (feedback cycle if a Material finding surfaces; otherwise gate-present + merge). Final deliverable: 4 in-repo commits (AC-1/2/3, whole-repo green, zero .md, fast-forwardable onto current main) + the provisioned/live-validated https://github.com/spacedock-dev/marketplace.

### Feedback cycle 1: MATERIAL — finished the decouple at the contract.go remedy site (commit 86fa74ef)

Detached audit: 1 MATERIAL finding, other 6 claims CLEAN. The decouple was half-done at a site outside my internal/cli grep scope: `internal/contract/contract.go` `pluginPredatesContractRemedy` still named the PLUGIN repo `spacedock-dev/spacedock` and rebuilt the removed `@branch` shorthand, so an edge binary's predates-contract verdict printed "...reinstalls from spacedock-dev/spacedock@next" — wrong repo + dead shorthand. Misleading user-facing TEXT (the actionable `spacedock install` works via the decoupled Install), but it tripped the decouple's own "no stale ref / no surviving shorthand" bar and a test locked the wrong string.

Fix (TDD, audit's preferred shape): the `branch` suffix is meaningless post-decouple — `spacedock install --host X` auto-selects the channel from the binary's devBranch stamp (entry name), not a source suffix. Dropped the dead `branch`/`devBranch` param from the whole contract chain (pluginPredatesContractRemedy, Compare, compareWithManifest, ManifestVerdict, RunDoctor) + the CLI callers (frontdoor.go, init.go) + test call sites. Remedy is now the clean "Upgrade it: spacedock install --host %s." Flipped the assertions first (must NOT name plugin repo / @branch), watched fail, fixed → green; negative guards now lock against regression. Removed a duplicate-name test collision (TestTooOldBinaryRemedyLeadsWithBrew already existed in version_message_test.go; folded the no-@next guard there, dropped my redundant copy).

Polish (same root cause): corrected the stale .goreleaser.yaml comment ("devBranch pins the marketplace @ref") → now "selects an entry NAME"; ldflag values were already correct. Skipped the redundant-assertion polish the auditor noted (redundant isn't wrong, per FO).

Verification: whole-repo `go test ./...` green; touched packages (contract, cli, release, integration) green; ZERO .md; stale-ref sweep clean (the only `spacedock-dev/spacedock` literals left are NEGATIVE test guards; no `@"+branch` build survives); 5 commits fast-forwardable onto current origin/main.
