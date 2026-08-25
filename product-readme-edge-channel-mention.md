---
id: dt8j3pas83725fma2wbez5ss
title: Product README documents no edge channel path
status: ideation
source: "Captain CL work-machine report 2026-08-25: README at v0.27.0 documents only brew tap/install (stable); someone following it for edge parity lands on stable and loses parity silently"
started: 2026-08-25T14:36:24Z
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:dt8j3pas83725fma2wbez5ss:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:dt8j3pas83725fma2wbez5ss-backlog-1
              briefing:
                id: briefing:dt8j3pas83725fma2wbez5ss:backlog:attempt-1:revision-1
                digest: sha256:c35c57a4aa2757547aae765122da7e1cea5396fe267ae9c5f52959b6799c9161
                request-digest: sha256:70b91d9acb281b900da42fd590b89cab6d58cb951ea020386a137234aa5709ab
                room-ref: ./product-readme-edge-channel-mention/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:dt8j3pas83725fma2wbez5ss:backlog:1
                briefing: briefing:dt8j3pas83725fma2wbez5ss:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-25T14:34:27.001703Z"
                decision: approve
                reason: 'Captain chat 2026-08-25: ''file and dispatch'' — approved seeding and immediate dispatch of this task into design; README posture decision reserved to the captain at the ideation gate'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:dt8j3pas83725fma2wbez5ss:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:dt8j3pas83725fma2wbez5ss-ideation-1
              briefing:
                id: briefing:dt8j3pas83725fma2wbez5ss:ideation:attempt-1:revision-1
                digest: sha256:e2a20a599f482863d12a18f79e69871a6c2af2c8a3939372ec381ff5b9c20582
                request-digest: sha256:4933bf8a68692136d9810ae24cf453d18f7c4c777c66bb43854923eabd06d1cb
                room-ref: ./product-readme-edge-channel-mention/review/ideation/briefing-1
---

The product README's install section documents only `brew tap spacedock-dev/homebrew-tap` + `brew install spacedock` — the stable channel. No edge cask, no `SPACEDOCK_CHANNEL=edge`, no pointer to the edge marketplace. A reader seeking edge parity follows the README, lands on stable, and loses parity with no signal. Decide the intended README posture and ship the README diff.

## Problem

README.md lines 60-64 (verified 2026-08-25) show only the stable brew steps. The archived task `marketplace-readme-channel-model-repair` (done 2026-08-24) fixed the MARKETPLACE repo README and docs/site install page; the product README was out of its scope.

Four findings from ideation shape the design:

1. **The install block already defers every other axis.** README.md:60-65 documents exactly one path — macOS Homebrew, stable. The Linux/binary `curl | sh` path, the Codex path, and the Pi path are all absent and deferred to `docs/site/get-started/install.md` via the "Full docs" pointer at README.md:77-80. Channel is not an oversight in an otherwise complete block; it is the one deferred axis carrying no label saying it was deferred.

2. **A bottom-of-section docs link is already there and already failed.** README.md:77-80 links the install guide, which documents edge in both of its tabs. The captain still landed on stable and lost parity silently. The absent thing is not a link to edge — it is a channel signal *adjacent to the commands*. Proximity is the fix, not link presence.

3. **The seed's "edge is a dev convenience" framing conflates two things.** CLAUDE.md's dev-only claim is about the `next` **branch** (a source-build convenience: `go install …@next`, `--plugin-dir`). The edge **channel** is a shipped user-facing release channel: cask `spacedock@next` v0.28.0-pre0 in `spacedock-dev/homebrew-tap`, marketplace `spacedock-edge` tracking main, published `-pre` tags, and `SPACEDOCK_CHANNEL=edge` in `install.sh`. The captain's own work machine runs the edge cask. So "should the README mention edge at all" is not settled against edge by CLAUDE.md; the posture call is genuinely open, which is why it is the captain's at the gate.

4. **The README holds the last unguarded copy of the install commands.** `internal/contractlint/install_hint_drift_test.go` binds the first-officer install-hint prose to install.md's tokens. The product README is not covered, and it has already drifted: README.md:63 says `brew tap spacedock-dev/homebrew-tap` while install.md:9 says `brew tap spacedock-dev/tap`. Both resolve to the same tap (verified below), so nothing is broken today — but nothing would catch it if it were.

## Proposed approach

Both postures are worked out below as concrete diffs against README.md at b04b3ef. **Posture (b) is the recommendation.** The captain rules at the gate; either diff is ready to apply as written.

One finding shrinks both diffs: the README does **not** need to document the plugin install. The binary's build stamp selects its own plugin channel, so installing the edge cask and launching is sufficient — verified live in Risk evidence. Neither posture carries `claude plugin marketplace add` lines.

### Posture (a) — document the edge install path in the README

````diff
--- a/README.md
+++ b/README.md
@@ -60,8 +60,17 @@
-Install with Homebrew:
+Install with Homebrew. Spacedock ships two channels: **stable** (default) and
+**edge**, which tracks prereleases. Pick one at install time — the two casks
+conflict, and the first launch installs the plugin matching the binary you
+installed.

 ```bash
-brew tap spacedock-dev/homebrew-tap
+brew tap spacedock-dev/tap
 brew install spacedock
+
+# Edge instead of the line above (prereleases; conflicts with the stable cask):
+brew install spacedock-dev/tap/spacedock@next
 ```

+Casks are macOS-only. The [install guide](docs/site/get-started/install.md)
+carries the script-based stable and edge commands for Linux.
+
 Then launch. The first launch sets up the plugin for you, so a single line gets
````

README.md +11/-2, net +9.

### Posture (b) — explicit stable-only statement with the edge pointer (RECOMMENDED)

````diff
--- a/README.md
+++ b/README.md
@@ -60,8 +60,15 @@
-Install with Homebrew:
+Install the stable channel with Homebrew:

 ```bash
-brew tap spacedock-dev/homebrew-tap
+brew tap spacedock-dev/tap
 brew install spacedock
 ```

+These commands install the **stable** channel, and the first launch below
+installs the matching stable plugin. Spacedock also ships an **edge** channel
+that tracks prereleases. The channel is chosen at install time, not with a flag
+later, so if you want edge — or you are matching a teammate who runs it — start
+from the [install guide](docs/site/get-started/install.md) instead of the
+commands above.
+
 Then launch. The first launch sets up the plugin for you, so a single line gets
````

README.md +9/-2, net +7.

Both diffs also normalize the tap token to install.md's spelling (`spacedock-dev/tap`). The two spellings are equivalent — Homebrew strips the `homebrew-` prefix, verified below — so no reader's install breaks; the change removes a divergence between two docs a reader compares side by side.

### Why (b)

1. **(a) breaks the block's editorial shape; (b) matches it.** The install block is deliberately one happy path with platform, host, and channel all deferred (Problem finding 1). Posture (a) adds *channel* branching to a block that still omits *platform* branching, so an edge-seeking Linux reader is stranded by the very block added to stop stranding — unless (a) also imports the script path, which grows the front door further. Posture (b) keeps channel in the same deferred bucket as platform and host, and adds the label that makes the deferral visible.

2. **(a) recreates the failure the previous task in this queue just repaired.** `marketplace-readme-channel-model-repair` (done 2026-08-24) existed because a second, less-maintained copy of the edge install instructions went stale and broke a fresh machine: it named a nonexistent `spacedock-edge` entry and claimed edge tracked `next`. Posture (a) creates a third copy of channel-specific commands in a document maintained on a pitch cadence, not a release cadence. Posture (b) creates none: it adds a channel *label*, which stays true across version bumps.

3. **(b) satisfies AC-1 in full.** The signal lands at the brew block, where the reader is. That is the whole delta against the status quo, whose only edge pointer sits at the bottom of the section and demonstrably did not work (Problem finding 2).

**Cost of (b), accepted:** the edge-seeking reader takes one hop to install.md. In exchange they arrive at platform-aware commands (mac cask *and* the Linux script) instead of a mac-only line the README would have to keep correct.

**What would flip the call to (a):** a decision that edge is a first-class advertised channel for new users rather than an opt-in for people who already know they want it. That is a product-positioning call, and it is the captain's — it is why both diffs are ready.

### New mechanism

One, and it is worth challenging on necessity: a new arm in `internal/contractlint/install_hint_drift_test.go` (AC-2). It serves AC-1's durability — a block labeled "stable" must keep naming the commands install.md documents, or the label becomes a lie. Simplest alternative, no check at all: insufficient by present evidence, since the divergence exists in the file today, no reviewer caught it, and the same commands are already guarded in the FO-prose copy. Second alternative considered and rejected: assert the words "stable"/"edge" appear in the section — a prose-grep that passes without helping any reader. The arm therefore checks only the cross-file token invariant; AC-1's reader-facing property stays a one-off human check. If the gate judges the arm unnecessary, drop AC-2: AC-1 is unaffected and the estimate falls to net +7 across 1 file.

## Risk evidence

**Spike run at ideation, not deferred.** Every command either posture documents was exercised against the shipped install machinery on this host (darwin/arm64, 2026-08-25). Results:

1. **Both tap spellings resolve to the same tap.** `brew tap-info spacedock-dev/homebrew-tap` and `brew tap-info spacedock-dev/tap` return byte-identical output — `spacedock-dev/tap: Installed`, same path `/opt/homebrew/Library/Taps/spacedock-dev/homebrew-tap`, same `HEAD: 4ed86b9fa445e625b2b6d18c2a03b0667a25e19a`. Homebrew strips the `homebrew-` repo prefix, so the README's current spelling is not broken and the normalization to `spacedock-dev/tap` is safe.

2. **Both cask names are real, and the conflict claim is real.** The tap ships exactly two casks. `Casks/spacedock.rb`: `cask "spacedock"`, `version "0.27.0"`, `conflicts_with cask: ["spacedock@next"]`. `Casks/spacedock@next.rb`: `cask "spacedock@next"`, `version "0.28.0-pre0"`, `conflicts_with cask: ["spacedock"]`, assets `spacedock_#{version}_darwin_arm64_edge.tar.gz`. So `brew install spacedock` and `brew install spacedock-dev/tap/spacedock@next` both name shipped casks, and posture (a)'s "the two casks conflict" is the casks' own declared metadata.

3. **The binary picks its own plugin channel — live.** `/opt/homebrew/Caskroom/spacedock@next/0.28.0-pre0/spacedock --version` prints `Channel: edge (spacedock@spacedock-edge)`, and `spacedock doctor` on the same binary prints `OK: spacedock binary 0.28.0-pre0 and plugin 0.28.0-pre0 are compatible.` This is what removes the plugin commands from both diffs: install the edge cask, launch, and the matching edge plugin is what you get.

4. **Both channel mappings proven, not just the installed one.** `go test ./internal/cli/ -run 'TestChannelMarketplaceFromDevBranch|TestChannelMarketplaceCarriesTheChannel|TestChannelPluginIDIsEntryAtMarketplace|TestClaudeNoPluginAutoInstallSelectsChannelEntry'` — all PASS, each with stable and edge subtests. They assert stable→marketplace `spacedock`, edge→`spacedock-edge`, entry always `spacedock`, and that auto-install selects the channel entry.

5. **`SPACEDOCK_CHANNEL=edge` is real, for the install.md path the README defers to.** `install.sh:38` reads `SPACEDOCK_CHANNEL` defaulting to `stable`; `install.sh:49-52` accepts only `stable|edge` and dies on anything else; `install.sh:131-132` selects the `_edge` asset. So the destination posture (b) points at documents a real switch.

**Honest limit:** no `brew install` transaction was run. Installing the stable cask would conflict with the edge cask this host runs, and mutating the host's install is out of proportion for a README diff. Cask names, versions, asset URLs, and conflict metadata were read from the tap's own shipped `.rb` files and confirmed via `brew tap-info`; the install transaction itself is unexercised. For posture (b) — the recommendation — the only commands the README documents are the stable pair, whose tokens are verified above to name a real tap and a real cask in that tap.

## Out of scope

The marketplace repo README (done: marketplace-readme-channel-model-repair). install.sh behavior. CLAUDE_CODE_PLUGIN_PREFER_HTTPS documentation (plugin-install-https-keyless-machines / rb0). The claude sibling-plugin cleanup (claude-install-sibling-channel-cleanup).

## Expected surface and tolerance

Estimate net LOC change: +37, across 2 files. Tolerance ±15.

Insertions 39, deletions 2, reported separately per the workflow's net-not-gross rule:

- `README.md` — +9/-2, net +7 (posture (b) diff above). Posture (a) instead: +11/-2, net +9.
- `internal/contractlint/install_hint_drift_test.go` — +30/-0, net +30 (the AC-2 arm plus its doc comment, reusing the existing `installMDSection` helper).

The estimate is above the backlog seed's +15 because the lint arm is the larger half. **If the gate drops AC-2, the figure becomes net +7 across 1 file** — record that as the approved baseline instead, so a later correction round calibrates against the posture actually approved.

**Observable semantics this task may change: none.** No command grammar, no stored format, no authority boundary, no runtime behavior. The output is README prose plus one test arm. The one behavior-adjacent edit is the documented tap token's spelling (`spacedock-dev/homebrew-tap` → `spacedock-dev/tap`), verified equivalent in Risk evidence finding 1, so no reader's install changes.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - Every install command block in the product README names the channel it installs, and a reader who wants the other channel is told where it lives — so no reader following the README toward edge silently lands on stable.**
Verified by: a one-off existence check recorded in the validation report, measured against a stated baseline that fails today. Baseline: README.md has 1 install command block and 0 channel labels, and its only edge pointer is the general "Full docs" link 13 lines below the block — the exact state that produced the captain's 2026-08-25 report. After: 1 of 1 blocks names its channel adjacent to the commands, and the other channel's location is named in the same paragraph. Presence is the claim, so the one-off check is legitimate evidence per the workflow's Proof policy. Falsifying edit: delete the added paragraph (posture b) or the edge command line and its comment (posture a) — the section reverts to unlabeled brew steps and the AC fails.

**AC-2 - The product README's Homebrew tap and cask tokens equal the tokens in install.md's `macOS (Homebrew)` tab, and a check fails if either side drifts.**
Verified by: a new arm in `internal/contractlint/install_hint_drift_test.go` that reuses `installMDSection` to extract `tap` and `formula` from install.md's Homebrew tab, then asserts README.md's Install section carries `brew tap <tap>` and `brew install <formula>`. Baseline that can move the wrong way: the arm FAILS on the pre-change README (`spacedock-dev/homebrew-tap` vs install.md's `spacedock-dev/tap`) and passes after. Falsifying edit: change the tap token in either file — the arm goes red. Drop this AC if the gate rules the mechanism unnecessary; see the New mechanism note above.

## Test plan

**Validation small-change fast path applies** (`docs/dev/README.md:192`). A labeled paragraph in a README plus one contractlint arm is low blast radius: no full validation checklist, no detached adversarial audit. Over-validating this diff is its own waste.

- **AC-1** — one-off existence check against the baseline stated in the AC, recorded in the validation report. No fixture, no live workflow run.
- **AC-2** — `go test ./internal/contractlint/ -run TestInstallHint`. Implementation MUST record the arm failing on the pre-change README first (paste the failure line showing the two tap tokens), then passing after the README edit. A drift arm never seen red is a tautology.
- **Regression** — `go test ./...` once, since the change touches a package that ships tests. No `-race` need: no concurrency in the change.
- **Spike** — already run at ideation (see Risk evidence); implementation inherits it and re-runs nothing. The channel-mapping tests cited there are existing tests, not new work.

Estimated cost: under an hour. The only sequencing constraint is AC-2's red-then-green ordering.

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}

## Stage Report: ideation

- DONE: Both README postures worked out as concrete before/after doc diffs in the task body — (a) document the edge install path, (b) explicit stable-only statement with a link to edge docs — with one recommendation; the captain rules on posture at the gate.
  Two unified diffs against README.md at b04b3ef under Proposed approach, each with its own LOC count; posture (b) recommended on three grounds (editorial shape of the install block, the stale-duplicate-copy precedent from `marketplace-readme-channel-model-repair`, and AC-1 satisfaction via proximity). Both are apply-ready, and the note "What would flip the call to (a)" states the positioning question that is the captain's.
- DONE: Every command the chosen diff documents is verified against the shipped install machinery once, or the body records `no spike needed: {what proves them}`.
  Spike run, not deferred — five findings in Risk evidence. Load-bearing ones: `brew tap-info` on both tap spellings returns identical output (same tap, same HEAD 4ed86b9), so the token normalization is safe; both cask names and the `conflicts_with` claim read from the tap's shipped `.rb` files; and the live edge binary prints `Channel: edge (spacedock@spacedock-edge)` with `doctor` reporting binary/plugin 0.28.0-pre0 compatible, which is what let both diffs drop the `claude plugin` lines. The honest limit (no `brew install` transaction run, and why) is recorded in the body.
- DONE: AC and test plan finalized: value-measuring AC (no reader silently lands on stable), validation small-change fast path noted, net-LOC estimate with tolerance.
  AC-1 (VALUE) measures channel-labeled install blocks against a baseline that fails today (1 block, 0 labels, edge pointer 13 lines away) with a named falsifying edit. AC-2 is the drift arm, whose baseline also fails today. Fast path cited at `docs/dev/README.md:192`. Estimate: net +37 across 2 files, tolerance ±15, insertions/deletions reported separately, plus the alternate baseline (net +7 across 1 file) if the gate drops AC-2. Semantic changes declared: none.

### Summary

Recommended posture (b) — label the README's brew block as stable and point at the install guide for edge — over posture (a). The decisive argument is structural rather than stylistic: the install block already defers platform, host, and channel to `install.md`, so posture (a) would add channel branching to a block that still omits platform branching, stranding an edge-seeking Linux reader in the very block meant to stop stranding. Posture (a) would also plant a third copy of channel-specific install commands in a pitch-cadence document, which is the exact failure the previous task in this queue existed to repair.

Two things the gate should look at. First, I pushed back on the seed's framing that CLAUDE.md settles this against edge: CLAUDE.md's dev-only claim is about the `next` *branch*, whereas the edge *channel* is a shipped user-facing channel with a published cask, a marketplace, and `-pre` tags — the captain's own machine runs it. The posture call is genuinely open. Second, the one new mechanism (the contractlint drift arm, AC-2) is flagged for a necessity challenge rather than smuggled in: it is separable, the body states what to drop and what the estimate becomes, and it exists because the README already carries a divergent tap spelling that nothing guards. I also found and fixed a rendering bug in my own body — nested code fences inside the diff blocks were closing the outer fence early, so the diffs are now wrapped in four-backtick fences.
