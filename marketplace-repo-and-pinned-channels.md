---
id: gpvg343kshg02e34kq7hrk9f
title: Separate marketplace repo + tag-pinned stable / edge channels (Model B core)
status: ideation
source: "Post-flip release-model decision (roadmap 0201, captain 2026-06-09): adopt Model B — stable channel serves a pinned release tag, edge serves next HEAD — grounded by the plugin-distribution research (wcdgsgd88). This is the structural core."
started: 2026-06-13T04:08:48Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0201-post-flip-release-model
group: release-model
sprint-readiness: ready
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
