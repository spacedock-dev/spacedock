---
id: 7p4vxcb9edzrph0zsfn75cd9
title: claude install leaves the sibling edge plugin installed and enabled
status: implementation
source: "Captain CL work-machine report 2026-08-25: clean documented stable install on a machine with spacedock@spacedock-edge left both plugins installed and enabled; doctor OK; which plugin serves unpredictable"
started: 2026-08-25T14:36:08Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-claude-install-sibling-channel-cleanup
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:7p4vxcb9edzrph0zsfn75cd9:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:7p4vxcb9edzrph0zsfn75cd9-backlog-1
              briefing:
                id: briefing:7p4vxcb9edzrph0zsfn75cd9:backlog:attempt-1:revision-1
                digest: sha256:4b88d6e7f37061f48b8d0475ee4f5a15476c1728bf602afc4d29856f86b71fb9
                request-digest: sha256:7652910a297d84b88917d9d48ad7e4009f83582811873821de42860027c6a649
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:7p4vxcb9edzrph0zsfn75cd9:backlog:1
                briefing: briefing:7p4vxcb9edzrph0zsfn75cd9:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-25T14:34:14.493119Z"
                decision: approve
                reason: 'Captain chat 2026-08-25: ''file and dispatch'' — approved seeding and immediate dispatch of this task into design after the FO''s verified source report'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:7p4vxcb9edzrph0zsfn75cd9:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:7p4vxcb9edzrph0zsfn75cd9-ideation-1
              briefing:
                id: briefing:7p4vxcb9edzrph0zsfn75cd9:ideation:attempt-1:revision-1
                digest: sha256:4a2309021baa7b1f887ecb610ccf7cda7b76cbda693a455ffe05eca2ce779e8b
                request-digest: sha256:faa415c897b74c56ce8189ab8b79b1daca58007fd611cd511556ab9da0c4d8d2
                room-ref: ./review/ideation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-25T15:09:29.517214Z"
                reason: 'Captain amendment 2026-08-25: fold the product README edge-channel deliverable into this task; design re-gates after the fold'
            - id: gate-attempt:7p4vxcb9edzrph0zsfn75cd9-ideation-2
              briefing:
                id: briefing:7p4vxcb9edzrph0zsfn75cd9:ideation:attempt-2:revision-1
                digest: sha256:2190d896a3e95805a7db88cbb1dd7fef4b9151ecfbf23837282b72fdf2d90f0b
                request-digest: sha256:1f7af3722d949c1b4e5e2a2101f68577e3b01fa80245bae62689cc52d1bc36d8
                room-ref: ./review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:7p4vxcb9edzrph0zsfn75cd9:ideation:2
                briefing: briefing:7p4vxcb9edzrph0zsfn75cd9:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-25T15:31:23.691252Z"
                decision: approve
                reason: 'Captain chat 2026-08-25: ''approve'' on briefing-2 — accepts the folded design as presented, AC-5 kept; approved baseline net +167 across 8 files, tolerance ±55/±1'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:7p4vxcb9edzrph0zsfn75cd9:validation
          stage: validation
          attempts:
            - id: gate-attempt:7p4vxcb9edzrph0zsfn75cd9-validation-1
              briefing:
                id: briefing:7p4vxcb9edzrph0zsfn75cd9:validation:attempt-1:revision-1
                digest: sha256:842cfa5fc2b3930a9dbd10a5bc213fc6ed24c540589e1c283fcd4c8242c6d53b
                request-digest: sha256:43ebe1c775612b378cecfcce6a7b84a16e4c80eb2f7b984ef4d91eb5cbfafa5a
                room-ref: ./review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:7p4vxcb9edzrph0zsfn75cd9:validation:1
                briefing: briefing:7p4vxcb9edzrph0zsfn75cd9:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-25T16:21:56.180302Z"
                decision: approve
                reason: 'Captain chat 2026-08-25: ''push it'' — accepts validation''s PASSED recommendation; deliver to merge'
              application:
                target-stage: done
                state: superseded
review-round:
    id: round:7p4vxcb9edzrph0zsfn75cd9:validation:2
    stage: validation
    cycle: 2
    briefing:
        id: briefing:7p4vxcb9edzrph0zsfn75cd9:validation:round-2:revision-1
        digest: sha256:0ba9df673235a4cf4556be4f9a791e1a60045346af8b8c6e403f499501170d3a
        room-ref: ./review/validation/round-2
---

`spacedock install --host claude` never uninstalls the sibling channel's modern plugin id. A stable install on a machine that already has `spacedock@spacedock-edge` (or the reverse) leaves BOTH channel plugins installed and enabled. Both are entry `spacedock` shipping the same skill set, so which provider serves is unpredictable. Mirror the codex sequence's sibling-channel remove.

**Scope after the captain's 2026-08-25 amendment** ("readme should fold into this. no
doc-only task"): this task ships TWO channel-hygiene deliverables that were separate
tasks, plus one that was cut.

1. **The install fix** — `installArgvSequence` gains the sibling-channel uninstall, so
   a claude install converges on exactly one enabled channel plugin. (AC-1, AC-2, AC-3.)
2. **The product README channel label** — folded in from the superseded
   `product-readme-edge-channel-mention` (`_archive/product-readme-edge-channel-mention.md`),
   whose ideation was complete and gate-verified. Its recommended posture (b) diff,
   value AC, drift arm, and spike evidence are folded below. (AC-4, AC-5.)
3. **Doctor visibility — CUT from this task, routed as a follow-up.** See "Doctor
   visibility" below. This is the one open thread; it was ruled out on necessity and
   on an unowned product decision, not deferred for effort.

The two shipped halves share a subject — which channel a machine ends up on, and
whether the user can tell — and both are proven against the same live install
machinery. They do not share code.

## Problem

`installArgvSequence` (internal/cli/host_exec.go:406) uninstalls only the own-channel id and the retired route-A id (`spacedock-edge@spacedock`). The modern sibling id is never touched. `codexInstallArgvSequence` (host_exec.go:445) already removes the sibling channel first — justified by codex's global skill namespace — and claude has the same practical collision: both channel plugins are entry `spacedock` with identical skill names. Neither safety net catches the dual install: `spacedock doctor` resolves only the running binary's own channel manifest (channel-aware resolver) and exits 0; the binary gate passes because both plugin snapshots can declare the same required binary minor.

Ideation confirms all of this against the source and, for the parts that are runtime
behavior, against a live claude (risk evidence below). One correction to the seed: the
front door does not merely fail to catch the dual install, it cannot heal it — the
auto-heal fires only on NoPluginFound or TooOldPlugin, and a dual install whose
own-channel manifest is Compatible never reaches either, so an affected machine stays
affected until someone runs `spacedock install` explicitly.

### The README half (folded from the superseded task)

`README.md` lines 60-65 (verified 2026-08-25: the "Install with Homebrew:" lead-in at
60 and the fenced command block at 62-65) document only `brew tap
spacedock-dev/homebrew-tap` + `brew install spacedock` — the stable channel. No edge
cask, no `SPACEDOCK_CHANNEL=edge`, no channel label. A reader seeking edge parity
follows the README, lands on stable, and loses parity with no signal. The archived task
`marketplace-readme-channel-model-repair` (done 2026-08-24) fixed the MARKETPLACE repo
README and the docs site install page; the product README was out of its scope.

Four findings from the folded ideation shape that half of the design:

1. **The install block already defers every other axis.** README.md:60-65 documents
   exactly one path — macOS Homebrew, stable. The Linux/binary `curl | sh` path, the
   Codex path, and the Pi path are all absent and deferred to
   `docs/site/get-started/install.md` via the "Full docs" pointer at README.md:77-80.
   Channel is not an oversight in an otherwise complete block; it is the one deferred
   axis carrying no label saying it was deferred.
2. **A bottom-of-section docs link is already there and already failed.**
   README.md:77-80 links the install guide, which documents edge in both of its tabs.
   The captain still landed on stable and lost parity silently. The absent thing is not
   a link to edge — it is a channel signal *adjacent to the commands*. Proximity is the
   fix, not link presence.
3. **"Edge is a dev convenience" conflates two things.** CLAUDE.md's dev-only claim is
   about the `next` **branch** (a source-build convenience: `go install …@next`,
   `--plugin-dir`). The edge **channel** is a shipped user-facing release channel: cask
   `spacedock@next` v0.28.0-pre0 in `spacedock-dev/homebrew-tap`, marketplace
   `spacedock-edge` tracking main, published `-pre` tags, `SPACEDOCK_CHANNEL=edge` in
   `install.sh`. The captain's own work machine runs the edge cask. So "should the
   README mention edge at all" is not settled against edge by CLAUDE.md; the posture
   call is genuinely open, which is why it stays the captain's at the re-gate.
4. **The README holds the last unguarded copy of the install commands.**
   `internal/contractlint/install_hint_drift_test.go` binds the first-officer
   install-hint prose to install.md's tokens. The product README is not covered, and it
   has already drifted: README.md:63 says `brew tap spacedock-dev/homebrew-tap` while
   install.md:9 says `brew tap spacedock-dev/tap`. Both resolve to the same tap
   (verified in Risk evidence), so nothing is broken today — but nothing would catch it
   if it were. This is the sole motivation for AC-5, and it is why AC-5 is flagged
   separable rather than assumed.

## Proposed approach

Add ONE tolerated step to `installArgvSequence` (internal/cli/host_exec.go:406), at
position 2, directly after the own-channel uninstall:

```go
{argv: []string{"plugin", "uninstall", "spacedock@" + otherChannelMarketplace(devBranch)}, tolerateExit: true},
```

The stable (devBranch=main) shape becomes six commands:

1. `plugin uninstall spacedock@spacedock` — own channel (tolerated)
2. `plugin uninstall spacedock@spacedock-edge` — **sibling channel (tolerated) — NEW**
3. `plugin uninstall spacedock-edge@spacedock` — retired route-A migration (tolerated)
4. `plugin marketplace add <channel source>` — fail-fast
5. `plugin marketplace update spacedock` — tolerated
6. `plugin install spacedock@spacedock` — fail-fast

Edge (any other devBranch) swaps the ids in steps 1 and 2 and targets `spacedock-edge`
in steps 5 and 6. This mirrors `codexInstallArgvSequence`'s existing sibling remove
(host_exec.go:445), which already carries channel exclusivity on codex.

**Position.** Step 2 groups the three uninstalls in channel-relevance order — own,
sibling, retired — and leaves step 1's documented ordering rationale untouched
(uninstall before any marketplace mutation, so no live uninstall is orphaned). The
sibling's marketplace record (`spacedock-edge` on a stable install) is never added,
updated, or removed by this sequence, so the new step has no ordering constraint of
its own and could sit anywhere among the cleanup steps.

**Tolerated, not conditional.** When the sibling is absent, claude exits non-zero
with two DIFFERENT messages depending on whether the sibling's marketplace happens
to be registered (both measured below at claude 2.1.226), so there is no stable
stderr shape to distinguish "nothing to remove" from a real failure. This is the
same tolerance rationale the two existing uninstall steps already carry.

**No new destructive surface.** No step spells `plugin marketplace remove`; the
existing ban is unchanged. `plugin uninstall` is plugin-level — it removes one
`<entry>@<marketplace>` record and leaves the marketplace registration and every
other plugin sourced from it in place (measured below). The sibling step therefore
cannot cascade into a third-party plugin co-hosted on the sibling marketplace, which
is what probe 1 measured `marketplace remove` doing.

### Mechanism justification

**Mechanism:** one tolerated argv step. **Value AC it serves:** AC-1.

- *Simplest alternative — query `claude plugin list --json` first and uninstall only
  when the sibling is present.* Insufficient: it buys no end-state difference (the
  unconditional tolerated step reaches the same state in BOTH the sibling-present and
  fresh-box directions, measured below), costs a host round-trip on every install, and
  moves the decision out of the argv list into `Install`'s control flow — breaking the
  property `installStep`'s own comment names ("the tolerance decision is visible to
  tests via `installArgvSequence` rather than hidden in `Install`'s control flow"),
  which every existing sequence test reads.
- *Alternative — document the manual `claude plugin uninstall <sibling>`.*
  Insufficient: the captain already ran exactly that command by hand. That IS the bug
  report, not the fix.
- *Alternative — `plugin marketplace remove` on the sibling marketplace.* Actively
  harmful and banned by the sequence's design rule: probe 1 measured that command
  cascade-uninstalling every plugin installed from the removed marketplace.

**Second mechanism, folded and FLAGGED FOR NECESSITY:** a new arm in
`internal/contractlint/install_hint_drift_test.go`. **Value AC it serves:** AC-4's
durability — a block labeled "stable" must keep naming the commands install.md
documents, or the label becomes a lie.

- *Simplest alternative — no check at all.* Insufficient by present evidence, narrowly:
  the divergence exists in the file today (Problem finding 4), no reviewer caught it,
  and the same commands are already guarded in the FO-prose copy by this very test file.
- *Alternative — assert the words "stable"/"edge" appear in the section.* Rejected: a
  prose-grep that passes without helping any reader, and banned outright by the
  workflow's Proof policy.

The arm therefore checks ONLY the cross-file token invariant (two independent files that
can diverge, which is the legitimate form); AC-4's reader-facing property stays a
one-off human check. **This mechanism is separable and the captain may drop it at the
re-gate:** dropping AC-5 leaves AC-4 unaffected and lowers the combined estimate to net
+137 across 7 files. The flag is preserved verbatim from the folded ideation rather than
quietly resolved, because the fold is not a ruling on it.

### Doctor visibility: SPLIT OUT, not shipped here

The backlog seed left this open. Decision: **route the doctor half to a separate
task**; this task ships the install fix only.

`doctor` today is `ManifestVerdict(manifestPath, host, binaryVersion, edgeCask)` over
a SINGLE manifest path (internal/contract/doctor.go), fed by
`hostOps.ResolveManifest(host)`, which returns one path for the running binary's own
channel. Making it see a sibling needs, at minimum: a new `hostOps` capability (a seam
every fake host in `internal/cli` implements), a new verdict or a repurposing of the
`Hint` channel, an exit-code decision, and a decision about `gateHost`. That last one
is the blocker — `gateHost`'s default branch fails fast, so a sibling condition that
reaches it as a non-Compatible verdict would make the FRONT DOOR refuse to launch on a
dual-install machine. That is a product/compatibility call this task cannot own; per
the review-finding disposition policy it is `Needs decision`.

Necessity argues for the split as well. After this fix, every `spacedock install
--host claude` and every launcher auto-heal converges to one enabled plugin, so the
population doctor would serve is exactly the machines that already hold a dual install
AND never run install again (a Compatible own-channel manifest never triggers a heal —
so it does not self-clear on launch, which is why the follow-up is worth filing rather
than dropping). One install run clears each such machine, and the fix announces itself:
the sibling step's own claude output (`✔ Successfully uninstalled plugin: spacedock`)
already flows through `runInstallSteps`' combined output into what `spacedock install`
prints, so the user-visible signal costs nothing.

Follow-up seed for the FO to file: *"doctor is blind to a co-installed enabled sibling
spacedock plugin"* — carries the AC-2 seed from this entity's backlog draft, plus a
finding surfaced here: `docs/site/reference/command-reference.md:42` already claims
"For what is installed for each host — plugin versions and enablement — use `spacedock
doctor`", which overstates what doctor does today (one manifest's version; it reads no
enablement at all). That doc line is out of scope here.

## Documentation diffs

Two files, both apply-ready. The first serves the install fix; the second is the folded
README deliverable.

### 1. `docs/site/get-started/install.md` — the exclusivity statement (serves AC-1)

`spacedock install --host claude` gains host-integration behavior a reader would want
stated, so ideation owns the doc change. One paragraph, inserted immediately after the
existing channel paragraph that ends "Codex installs the same way with `codex plugin
add`." (line 82):

```diff
 the `spacedock-edge` marketplace the edge entry resolves from. Codex installs the
 same way with `codex plugin add`.
 
+`spacedock install --host claude` and `spacedock install --host codex` keep exactly
+one channel installed per host: the selected channel's plugin replaces any
+co-installed stable or edge Spacedock plugin. Both channels ship the entry name
+`spacedock` and the same skill names, so leaving both installed makes which one
+serves unpredictable. Installing a channel by hand with the commands above does not
+remove the other; the next `spacedock install` — or the launcher's auto-install — does.
+
 Set `SPACEDOCK_MARKETPLACE_SOURCE` to install from a local or alternate
```

`docs/site/reference/command-reference.md` was considered and declined: its
`spacedock install` row and the paragraph below it describe the compatibility check,
not channel semantics, and the install page is where channel behavior already lives.
Repeating the fact in two files invites drift.

### 2. `README.md` — posture (b), the channel label (serves AC-4)

Folded from the superseded task, which worked out TWO postures as apply-ready diffs and
recommended (b). **Posture (a) — documenting the edge install path in the README — was
NOT selected.** Its full diff and the three-part argument against it are preserved at
`_archive/product-readme-edge-channel-mention.md` (Proposed approach, "Why (b)"); the
short form is that (a) adds *channel* branching to a block that still omits *platform*
branching, stranding an edge-seeking Linux reader in the very block meant to stop
stranding, and it plants a third copy of channel-specific install commands in a
pitch-cadence document — the exact failure `marketplace-readme-channel-model-repair`
existed to repair. **What would flip the call to (a):** a decision that edge is a
first-class advertised channel for new users rather than an opt-in for people who
already know they want it. That is product positioning, and it is the captain's call at
the re-gate — which is why both diffs remain on file.

One folded finding shrinks this diff: the README does **not** need to document the
plugin install. The binary's build stamp selects its own plugin channel, so installing
the edge cask and launching is sufficient (verified live — Risk evidence, README spike
finding 3). Neither posture carries `claude plugin marketplace add` lines.

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

The diff also normalizes the tap token to install.md's spelling (`spacedock-dev/tap`).
The two spellings are equivalent — Homebrew strips the `homebrew-` prefix, verified in
Risk evidence — so no reader's install breaks; the change removes a divergence between
two docs a reader compares side by side, and it is what makes AC-5's arm go green.

**Cost of (b), accepted:** the edge-seeking reader takes one hop to install.md. In
exchange they arrive at platform-aware commands (mac cask *and* the Linux script)
instead of a mac-only line the README would have to keep correct.

## Risk evidence

Backlog: captain live repro on v0.27.0 — after the documented stable install steps on
a machine with an existing edge plugin, `claude` listed spacedock@spacedock v0.27.0
AND spacedock@spacedock-edge v0.27.0-pre8, both enabled; workaround was `claude plugin
uninstall spacedock@spacedock-edge`. Source check 2026-08-25: v0.27.0 and main
installArgvSequence are byte-identical — no sibling uninstall step exists.

### Ideation spike (2026-08-25, claude 2.1.226)

Riskiest unverified mechanism first: claude's REAL `plugin uninstall` behavior for the
sibling id, absent vs installed+enabled. Run against isolated `CLAUDE_CONFIG_DIR` /
`CLAUDE_CODE_PLUGIN_CACHE_DIR` (the `co_hosted_survival_test.go` isolation pattern);
both spike config dirs deleted afterward and the ambient `~/.claude` verified unchanged.

**Standing state — the bug is live on this machine right now.** `claude plugin list
--json` reports `spacedock@spacedock 0.27.0-pre7+dev` enabled AND
`spacedock@spacedock-edge 0.28.0-pre0` enabled, both user scope. Left in place
deliberately as a validation baseline; a captain-owned machine is not this task's to
mutate.

1. **Absent sibling, no marketplaces registered** — `claude plugin uninstall
   spacedock@spacedock-edge` → **exit 1**, `Plugin "spacedock@spacedock-edge" not found
   in installed plugins`.
2. **Absent sibling, sibling marketplace IS registered** — same command → **exit 1**,
   `Plugin "spacedock@spacedock" is not installed in user scope. Use --scope to specify
   the correct scope.` Two distinct non-zero shapes for the same "nothing to remove"
   condition — this is what forces `tolerateExit: true` rather than a fail-fast step or
   a stderr-shape test.
3. **Installed + enabled sibling** — same command → **exit 0**, `✔ Successfully
   uninstalled plugin: spacedock (scope: user)`. Afterward exactly ONE spacedock entry
   remains (`spacedock@spacedock`, still `enabled: true`). Enablement of the survivor
   is NOT disturbed. Both marketplace records (`spacedock`, `spacedock-edge`) survive —
   no cascade. The sibling's plugin cache directory remains on disk with no plugin
   entry pointing at it, which is why the AC counts installed+enabled plugin ENTRIES,
   not cache directories.
4. **End-to-end, current 5-step sequence** (start: edge installed+enabled) → **2
   enabled** (`spacedock@spacedock`, `spacedock@spacedock-edge`). The captain's report
   reproduced hermetically.
5. **End-to-end, proposed 6-step sequence**, same start → **1 enabled**
   (`spacedock@spacedock`). Steps 1 and 3 both exited non-zero as expected and did not
   abort the run.
6. **Fresh box, proposed sequence** (start: 0 spacedock plugins) → **1 enabled**. No
   regression from the extra tolerated step.
7. **Reverse direction** (edge install over an installed+enabled stable sibling) → **1
   enabled** (`spacedock@spacedock-edge`). The mapping works both ways.

### Second spike: the test fixture the value AC depends on

The AC-1 proof needs a local two-channel marketplace, which nothing in the repo builds
yet (`buildLocalMarketplaceWithDependent` makes one marketplace named `spacedock`).
Verified buildable before proposing it: two local directory marketplaces whose
`marketplace.json` names are `spacedock` and `spacedock-edge`, each with the single
`spacedock` entry, both registered successfully under distinct names from distinct
Directory sources — no name/source collision. Re-running items 4 and 5 against that
local fixture reproduced the same discrimination: **current sequence → 2 enabled,
proposed sequence → 1 enabled**. `plugin marketplace update` on a local Directory
source exits 0. The fixture needs no network, so the value test is offline-hermetic.

### README spike (folded, run at the superseded task's ideation, darwin/arm64 2026-08-25)

Every command posture (b) documents was exercised against the shipped install
machinery — spiked, not deferred.

1. **Both tap spellings resolve to the same tap.** `brew tap-info
   spacedock-dev/homebrew-tap` and `brew tap-info spacedock-dev/tap` return
   byte-identical output — `spacedock-dev/tap: Installed`, same path
   `/opt/homebrew/Library/Taps/spacedock-dev/homebrew-tap`, same `HEAD:
   4ed86b9fa445e625b2b6d18c2a03b0667a25e19a`. Homebrew strips the `homebrew-` repo
   prefix, so the README's current spelling is not broken and the normalization to
   `spacedock-dev/tap` is safe for every reader.
2. **Both cask names are real, and the conflict claim is real.** The tap ships exactly
   two casks. `Casks/spacedock.rb`: `cask "spacedock"`, `version "0.27.0"`,
   `conflicts_with cask: ["spacedock@next"]`. `Casks/spacedock@next.rb`: `cask
   "spacedock@next"`, `version "0.28.0-pre0"`, `conflicts_with cask: ["spacedock"]`,
   assets `spacedock_#{version}_darwin_arm64_edge.tar.gz`. Both `brew install` targets
   name shipped casks, and the "two casks conflict" claim is the casks' own declared
   metadata, not an inference.
3. **The binary picks its own plugin channel — live.**
   `/opt/homebrew/Caskroom/spacedock@next/0.28.0-pre0/spacedock --version` prints
   `Channel: edge (spacedock@spacedock-edge)`, and `spacedock doctor` on the same binary
   prints `OK: spacedock binary 0.28.0-pre0 and plugin 0.28.0-pre0 are compatible.` This
   is what removes the plugin commands from the diff: install the edge cask, launch, and
   the matching edge plugin is what you get.
4. **Both channel mappings proven, not just the installed one.** `go test
   ./internal/cli/ -run
   'TestChannelMarketplaceFromDevBranch|TestChannelMarketplaceCarriesTheChannel|TestChannelPluginIDIsEntryAtMarketplace|TestClaudeNoPluginAutoInstallSelectsChannelEntry'`
   — all PASS, each with stable and edge subtests (existing tests, cited not written).
5. **`SPACEDOCK_CHANNEL=edge` is real**, so posture (b)'s pointer has a real destination:
   `install.sh:38` reads `SPACEDOCK_CHANNEL` defaulting to `stable`, `install.sh:49-52`
   accepts only `stable|edge` and dies otherwise, `install.sh:131-132` selects the
   `_edge` asset.

**Honest limit, carried forward unchanged:** no `brew install` transaction was run.
Installing the stable cask would conflict with the edge cask this host runs, and
mutating the host's install is out of proportion for a README diff. Cask names,
versions, asset URLs, and conflict metadata were read from the tap's own shipped `.rb`
files and confirmed via `brew tap-info`; the install transaction itself is unexercised.
For posture (b) — the recommendation — the only commands the README documents are the
stable pair, whose tokens are verified above to name a real tap and a real cask in it.

## Out of scope

CLAUDE_CODE_PLUGIN_PREFER_HTTPS documentation (plugin-install-https-keyless-machines /
rb0). The codex install sequence (already removes the sibling). The marketplace repo
README (done: marketplace-readme-channel-model-repair). `install.sh` behavior — read as
evidence that posture (b)'s pointer has a real destination, never edited.

(The former cross-reference to `product-readme-edge-channel-mention` is dropped: that
task is superseded and archived, and its deliverable is now IN scope here.)

Split out by this ideation (see "Doctor visibility" above): doctor/`--version`
visibility into a co-installed enabled sibling, and the overstated
`docs/site/reference/command-reference.md:42` claim that doctor reports per-host
enablement. Also out of scope: healing an existing dual install without an explicit
`spacedock install` run — the launcher's auto-heal fires only on NoPluginFound or
TooOldPlugin, and a dual install whose own-channel manifest is Compatible does not
trigger it. Changing that is the same front-door gate decision the doctor follow-up
owns.

## Expected surface and tolerance

Combined figure for both folded deliverables.

Estimate net LOC change: **+167, across 8 files** (insertions ~+173, deletions ~-6).
Tolerance: **±55 net LOC and ±1 file**.

| File | Net | What | Half |
| --- | --- | --- | --- |
| `internal/cli/host_exec.go` | +6 | the argv step (+1) and the sequence doc comment (5-command → 6-command, sibling rationale, "three cleanup steps" → four) | install |
| `internal/cli/init_devbranch_test.go` | +3 | `TestInstallArgvSequence`: one literal step in each of `wantEdge`/`wantStable`, comment | install |
| `internal/cli/channel_selection_test.go` | +4 | `TestClaudeChannelInstallArgvSequence`: `wantSibling` table field, one literal step, comment | install |
| `internal/cli/install_tolerance_test.go` | +3 | three existing tests each expect one more stub marker; "five steps" → "six steps" in comments | install |
| `internal/cli/sibling_channel_exclusive_test.go` (new) | +108 | the value test plus the two-channel local marketplace fixture builder | install |
| `docs/site/get-started/install.md` | +6 | the exclusivity paragraph | install |
| `README.md` | +7 | posture (b): +9/-2 — stable label, edge pointer, tap token normalized | README |
| `internal/contractlint/install_hint_drift_test.go` | +30 | the AC-5 drift arm plus its doc comment, reusing the existing `installMDSection` helper | README |

Subtotals: install half **+130 across 6 files**; README half **+37 across 2 files**.

**Alternate approved baseline if the captain drops AC-5 at the re-gate: net +137 across
7 files, tolerance ±45.** Record whichever the gate approves, so a later correction
round calibrates against the posture actually approved rather than this superset.

The install half's backlog seed said +40/4 files; ideation revised it upward because the
production change is 6 lines and the rest is the value test's two-channel fixture, which
nothing in the repo builds yet. Tolerance rationale: the fixture builder's final size is
the only real unknown (the spike built it in shell, not Go); the README half's two diffs
are written out line for line and carry almost no uncertainty, which is why folding it
in adds only ±15 to the tolerance. The ±1 file covers putting the fixture helper beside
the existing one in `co_hosted_survival_test.go` instead of a new file.

### Declared semantic changes

- **Runtime behavior (intended, the point of the task):** `spacedock install --host
  claude` now uninstalls the sibling channel's plugin. A captain deliberately holding
  both channels on claude loses that arrangement on the next install or auto-heal.
  Codex has behaved this way since its sequence was written; this makes claude match.
- **Command output:** `spacedock install --host claude` gains claude's own line for the
  extra step (a success line when the sibling was present, a tolerated failure line
  when it was not). No Spacedock-authored output changes.
- **README half: none.** No command grammar, no stored format, no authority boundary,
  no runtime behavior — the output is README prose plus one test arm. The one
  behavior-adjacent edit is the documented tap token's spelling
  (`spacedock-dev/homebrew-tap` → `spacedock-dev/tap`), verified equivalent in Risk
  evidence, so no reader's install changes.
- **Unchanged:** command grammar, flags, stored formats, exit codes, authority, the
  version gate, doctor, the `plugin marketplace remove` ban, and the tolerance
  asymmetry rule (four cleanup steps tolerated, two pinning steps fail-fast).

### Path-to-lane call

**`claude-live` is NOT required.** The mapping keys on what a live lane loads or drives.
`claude-live` runs `go test -tags live ./internal/ensigncycle`, and that package
contains no reference to `Install`, `installArgvSequence`, or any `plugin install` /
`plugin marketplace add` invocation — it drives dispatch and launch scenarios, not the
install path. The diff touches `internal/cli` (install argv + its tests) and one docs
page; it does not touch `skills/**/references/**`, the dispatch core, or any live
lane's own tests. Deterministic lanes (`go test ./...`, docs strict build) plus the
local claude-on-PATH value test are the required set.

**Detached adversarial audit: YES.** `host_exec.go` is front-door launcher machinery —
`installArgvSequence` is what `resolveHealableGate` shells during an auto-heal on
launch — so the four-surface trigger fires. The provenance trigger fires too: the
audit must confirm the sibling id in the argv tests is a literal that diverges from
`otherChannelMarketplace`, not a re-derivation of it. Scope the audit at the install
half; the README half is inside the same diff but is not what the trigger is about.

**Docs strict build covers one doc file, not both.** `mkdocs.yml` does not include
`README.md` (checked — no README entry), so `mkdocs build --strict` validates
`docs/site/get-started/install.md` and not the product README. The README's only
mechanical guard is AC-5's drift arm, which is precisely the gap Problem finding 4
describes — and precisely why dropping AC-5 leaves the README's tokens unguarded again.
Worth the captain's attention when ruling on AC-5's necessity.

**Consequence of the fold the re-gate should see:** the superseded task claimed the
validation small-change fast path (`docs/dev/README.md:192`) for its README diff, and
that claim was correct in isolation. **It no longer applies.** The combined task touches
the front-door launcher, so the whole deliverable carries the heavier posture — full
validation and the detached audit. Folding a low-blast-radius doc change into a
high-stakes-surface change raises the doc change's validation cost; that is a real,
accepted cost of the captain's "no doc-only task" ruling, not an oversight. It does not
change the work, only how it is checked before merge.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - After `spacedock install --host claude` on a host with the sibling channel plugin installed and enabled, exactly one spacedock channel plugin is installed and enabled: the selected channel's.**

The measured quantity is the COUNT of enabled `spacedock@*` entries in `claude plugin
list --json` — claude's own reporting, not a Spacedock value. It moves the wrong way
today: measured 2 under the current sequence and 1 under the proposed one, from an
identical starting state (risk evidence items 4 and 5).

Verified by: a new claude-on-PATH test in `internal/cli` driving the production seam
`execHost{}.Install("claude", <local stable marketplace>, devBranch)` against the
two-channel local marketplace fixture, isolated with `CLAUDE_CONFIG_DIR` and
`CLAUDE_CODE_PLUGIN_CACHE_DIR` (the `co_hosted_survival_test.go` shape), asserting the
enabled-entry count is exactly 1 and its id is the selected channel's. Three subtests:
sibling present (start 1 sibling → end 1, selected channel), fresh box (start 0 → end
1, proving the tolerated non-zero exits do not abort), and the edge direction (start 1
stable → end 1 edge). Falsifying edit: delete the sibling step from
`installArgvSequence` — the sibling-present and edge subtests then count 2 and fail.
Skips when `claude` is not on PATH, so it does not run in CI's `build` job; see the
test plan for how that is settled.

**AC-2 (MEANS, serving AC-1) - The claude install sequence issues the sibling-channel uninstall as a tolerated step targeting the marketplace name the running channel does NOT select, and the tolerance asymmetry is unchanged: all four cleanup steps tolerated, both pinning steps fail-fast.**

Verified by: the literal-argv assertions in `TestInstallArgvSequence`
(`init_devbranch_test.go`) and `TestClaudeChannelInstallArgvSequence`
(`channel_selection_test.go`). The sibling id is written per channel as an independent
literal — `"spacedock@spacedock-edge"` in the stable case, `"spacedock@spacedock"` in
the edge case — never as `"spacedock@" + otherChannelMarketplace(devBranch)`, so the
test and the production function can diverge. `TestInstallArgvSequence`'s existing
cleanup/pinning loop covers the new step's tolerance flag automatically, because the
step is an uninstall. Falsifying edits: misspell the sibling id back to the retired
route-A id `spacedock-edge@spacedock`; derive it from the running channel instead of
the sibling (making stable uninstall itself); or flip `tolerateExit` to false. Each
fails a named assertion.

**AC-3 (MEANS, serving AC-1's discoverability) - The install documentation states that Spacedock keeps exactly one channel installed per host.**

Verified by: doc diff 1 above applied to `docs/site/get-started/install.md`, reviewed
at the validation gate, with the docs strict build (`mkdocs build --strict`, the docs
lane) as the mechanical check. Stated honestly: no committed test proves prose says
what we wrote — a grep over it would be the banned prose-grep. This AC counts only as
a companion to AC-1, which is where the behavior is measured.

---

*AC-4 and AC-5 are the folded README deliverable. AC-4 is its own value AC — a distinct
end-value from AC-1, measured against its own baseline — not a mechanism serving the
install fix.*

**AC-4 (VALUE) - Every install command block in the product README names the channel it installs, and a reader who wants the other channel is told where it lives — so no reader following the README toward edge silently lands on stable.**

Verified by: a one-off existence check recorded in the validation report, measured
against a stated baseline that fails today. **Baseline (re-verified 2026-08-25 at
b04b3ef):** README.md has 1 install command block (lines 62-65) and 0 channel labels,
and its only edge pointer is the general "Full docs" link at lines 77-80, thirteen lines
below the block — the exact state that produced the captain's report. **After:** 1 of 1
blocks names its channel adjacent to the commands, and the other channel's location is
named in the same paragraph. Presence is the claim, so a one-off check is legitimate
evidence here per the Proof policy's own carve-out ("presence or absence is an existence
fact, and a grep establishes it soundly when that fact is itself the claim") — this is
NOT a standing prose-grep test and must not be committed as one. Falsifying edit: delete
the added paragraph and the section reverts to unlabeled brew steps; the AC fails.

**AC-5 (MEANS, serving AC-4's durability — SEPARABLE, captain may drop at the re-gate) - The product README's Homebrew tap and cask tokens equal the tokens in install.md's `macOS (Homebrew)` tab, and a check fails if either side drifts.**

Verified by: a new arm in `internal/contractlint/install_hint_drift_test.go` that reuses
`installMDSection` to extract `tap` and `formula` from install.md's Homebrew tab, then
asserts README.md's Install section carries `brew tap <tap>` and `brew install
<formula>`. This binds two independent files that can diverge — the legitimate form,
not a prose-grep. **Baseline that moves the wrong way:** the arm FAILS on the pre-change
README (`spacedock-dev/homebrew-tap` at README.md:63 vs install.md:9's
`spacedock-dev/tap` — both re-verified today) and passes after. Falsifying edit: change
the tap token in either file and the arm goes red. Implementation MUST record it red
before green; a drift arm never seen red is a tautology.

**Necessity flag, preserved for the captain's ruling:** this is the only new mechanism
in the folded half, and the folded ideation deliberately flagged it rather than
smuggling it in. Drop AC-5 and AC-4 is unaffected; the estimate becomes net +137 across
7 files. The argument for keeping it is that the divergence exists in the file today and
nothing caught it; the argument against is that one stale token in a README is cheap to
fix by hand and a standing check is a permanent cost for a rare fault.

## Test plan

| Test | Proves | Cost | Lane |
| --- | --- | --- | --- |
| `TestClaudeInstallLeavesOnlySelectedChannelEnabled` (new, `internal/cli`) | AC-1 | ~108 LOC, ~15s, offline | local + any runner with claude on PATH; **skips** in CI `build` |
| `TestInstallArgvSequence` (extend) | AC-2 | +3 LOC, instant | CI `build` (`go test ./...`) |
| `TestClaudeChannelInstallArgvSequence` (extend) | AC-2 | +4 LOC, instant | CI `build` |
| three `TestInstallTolerates…` (extend marker lists) | no regression in `runInstallSteps` tolerance | +3 LOC, ~1s | CI `build` |
| `TestClaudeChannelInstallLeavesCoHostedPluginInstalled` (existing, unchanged) | the sibling step did not reintroduce cascade harm | 0 LOC | local, claude on PATH |
| docs strict build (existing) | AC-3 well-formed | 0 LOC | docs lane |
| one-off existence check vs the stated baseline | AC-4 | 0 LOC, recorded in the validation report | not committed — see AC-4 |
| new arm in `install_hint_drift_test.go` | AC-5 | +30 LOC, instant | CI `build` (`go test ./internal/contractlint/ -run TestInstallHint`) |

**The value test's fixture.** A new helper builds two local directory marketplaces
whose `marketplace.json` names are `spacedock` and `spacedock-edge`, each with the
single `spacedock` entry — proven registrable side by side in the second spike. The
existing `buildLocalMarketplaceWithDependent` is left alone rather than generalized,
so the codex co-hosted test that shares it does not churn.

**The CI skip, settled honestly.** CI's `build` job runs `go test ./...` on a runner
with no `claude`, so the value test skips there — the same posture the existing
`TestClaudeChannelInstallLeavesCoHostedPluginInstalled` has carried since round 2, for
the same class of claim. Installing claude into the `build` job would change CI
machinery (a high-stakes surface) for one test, out of proportion to a 6-line
production fix. The arrangement instead is: AC-2's argv tests guard the sequence
unconditionally in CI, and **validation MUST run `go test ./internal/cli -run
TestClaudeInstallLeavesOnlySelectedChannelEnabled -count=1 -v` locally with claude on
PATH and paste the result** — a SKIP is not a pass and must be reported as a failure to
prove AC-1.

**No test is proposed for these**, deliberately:
- *The sibling step's own tolerance end-to-end through `runInstallSteps`.* A fourth
  `TestInstallTolerates…` against a PATH stub would re-cover the same shared
  `if err != nil { if step.tolerateExit { continue } }` branch with a different
  `failOn` string. The flag is asserted by `TestInstallArgvSequence`'s loop and the
  continue-behavior by the three existing tests; composition covers it because
  `runInstallSteps` special-cases no step. The value test's fresh-box subtest exercises
  the same path against the REAL claude, which is stronger for less code.
- *Cascade safety of the sibling uninstall.* Structural, not empirical: `plugin
  uninstall` names one `<entry>@<marketplace>` record. Confirmed incidentally in spike
  item 3 (both marketplace records survived).

**AC-5's red-then-green obligation.** Implementation MUST run `go test
./internal/contractlint/ -run TestInstallHint` and record the arm FAILING on the
pre-change README first — pasting the failure line showing the two divergent tap tokens
— then edit README.md and record it passing. This ordering is the only sequencing
constraint in the whole task. A drift arm first seen green proves nothing: it would
equally "pass" if it asserted nothing at all.

**Regression:** `go test ./...` once. No `-race` needed for the README half (no
concurrency); the install half introduces none either, but `-race` runs anyway as the
repo's standing command.

**Cross-half interaction: none.** The two halves touch disjoint files and disjoint
tests. The only shared object is `docs/site/get-started/install.md`, which the install
half EDITS (adding the exclusivity paragraph) and which AC-5's arm READS (extracting the
Homebrew tab's tap/formula tokens). Those do not collide — the added paragraph is in the
`## Skills` section, not the `macOS (Homebrew)` tab `installMDSection` parses — but
implementation should run the drift arm after the install.md edit to confirm it, since a
surprise there would be cheap to catch and confusing to debug later.

**Estimated total cost:** ~1.5 hours implementation — ~1 hour for the install half
(dominated by the fixture builder), ~30 minutes for the folded README half.

### Feedback Cycles

- Cycle 1: REWORKED — captain send-back of PR #760 (ASD-STE100 ruling, 2026-08-25; correction commit 2f91b6dd4); surface 8 files/+238 net vs estimate +167±55 (142%); AC unchanged; surface is 16 net LOC past tolerance — cause is the STE conversion itself, baseline predates the prose ruling; captain tolerance decision pending before the validation re-run.
- Cycle 2: REWORKED — captain content-cut directive 2026-08-25 ("cut things we don't need, not just purely convert"; correction commit 4b5b4073a); surface 8 files/+147 net vs estimate +167±55 (88%); AC unchanged; the cycle-1 tolerance breach is closed by deletion, not waived — no captain tolerance decision needed.

## Stage Report: ideation

- DONE: Design lands the concrete sibling-cleanup mechanism: the exact tolerated uninstall argv added to installArgvSequence (mirroring codexInstallArgvSequence's sibling remove) and an explicit decision whether doctor dual-install visibility ships here or splits out — naming the value AC each mechanism serves, the simplest alternative, and why it is insufficient.
  Proposed approach gives the literal step (`{argv: []string{"plugin", "uninstall", "spacedock@" + otherChannelMarketplace(devBranch)}, tolerateExit: true}`) at position 2, the full 6-command shape, and a Mechanism justification naming AC-1 plus three rejected alternatives (conditional pre-query, manual docs, `marketplace remove`). Doctor: SPLIT OUT, with the front-door `gateHost` fail-fast interaction as the reason it is a `Needs decision` this task cannot own, and a follow-up seed for the FO to file.
- DONE: Riskiest unverified mechanism exercised first and recorded in the task body (e.g. claude's real `plugin uninstall` behavior for an absent vs installed+enabled sibling — exit codes, whether enabled state survives), or an auditable `no spike needed: {proven mechanisms}` line.
  Seven-item spike against claude 2.1.226 in an isolated `CLAUDE_CONFIG_DIR`, recorded under Risk evidence: absent sibling exits 1 with TWO different messages (which is what forces tolerance); installed+enabled sibling exits 0 and leaves the survivor still enabled with both marketplace records intact; end-to-end current sequence 1 → 2 enabled vs proposed 1 → 1; fresh box 0 → 1; reverse direction 1 → 1. A second spike proved the two-channel LOCAL marketplace fixture the value test depends on is buildable and offline.
- DONE: Entity-level AC set finalized with a value-measuring AC against an independent baseline, per-AC test plan (literal-argv sequence test not derived from channelPluginID; doctor dual-manifest fixture), net-LOC surface with tolerance, declared semantic changes, and the path-to-lane call on whether claude-live is required.
  AC-1 measures the count of enabled `spacedock@*` entries in claude's own `plugin list --json` (2 today, must be 1 — measured both ways). AC-2 is the literal-argv test with the sibling id as a per-channel literal that can diverge from `otherChannelMarketplace`. AC-3 is the doc paragraph, with its proof limits stated rather than dressed up as a test. Surface: +130 net across 6 files, tolerance ±40 LOC / ±1 file, with a per-file table and the reason the backlog's +40/4 was low. Semantic changes declared (runtime behavior changes intentionally; command grammar, formats, exit codes, authority unchanged). Path-to-lane: claude-live NOT required (`internal/ensigncycle` contains no install-path reference); detached adversarial audit IS required on both triggers.
  The checklist's "doctor dual-manifest fixture" is deliberately absent: the same checklist's first item authorized deciding the doctor half out of scope, and it went to the follow-up with its AC seed intact.

### Summary

The install fix is 6 production lines; the design work was deciding what NOT to build
alongside it. Doctor visibility splits out because making doctor see a sibling requires
a new `hostOps` capability and an exit-code decision that reaches `gateHost`'s
fail-fast branch — which would make the front door refuse to launch on a dual-install
machine. That is a product call, and bundling it would let a doctor disagreement block
the fix that stops new dual installs. Two spikes de-risked the design: claude's real
`plugin uninstall` semantics (the sibling removal is clean, the survivor stays enabled,
nothing cascades) and the two-channel local marketplace fixture the value test needs.
The bug is currently live on the captain's machine — both channels enabled — and was
left in place as a validation baseline rather than quietly fixed.

Two things the gate should weigh. First, the value test skips in CI's `build` job
because no `claude` is installed there; I kept the existing co-hosted test's posture
rather than change CI machinery for one test, and made it a validation obligation to
run it locally and treat a SKIP as a failure. Second, the change removes a co-installed
sibling from a captain's claude on the next install or auto-heal, which is intended but
is a real behavior change for anyone deliberately holding both channels.

### Amendment: README deliverable folded in (2026-08-25)

Captain ruling "readme should fold into this. no doc-only task" withdrew briefing-1 and
gave this task the deliverable from `product-readme-edge-channel-mention`, now
superseded and archived. Its ideation was complete and gate-verified, so this was a fold
of finished design, not a re-design.

- DONE: Fold the posture (b) README diff into the Proposed approach / Documentation diff section.
  "Documentation diffs" is now two numbered diffs — install.md (serves AC-1) and README.md posture (b) (serves AC-4). Posture (a) was NOT selected and is cited, not copied, at `_archive/product-readme-edge-channel-mention.md`, with the short-form argument against it and the "what would flip the call" note kept inline so the captain can still rule without opening the archive. Tap token normalized to `spacedock-dev/tap`.
- DONE: Fold the value AC and the contractlint drift arm AC into Acceptance criteria, renumbered, mechanism→value pairings explicit, drift arm's separability/necessity flag kept.
  AC-1 (VALUE, install) → AC-2, AC-3 as its means; AC-4 (VALUE, README) → AC-5 as its means. Two value ACs because the halves deliver two distinct end-values against two distinct baselines. AC-5 carries the separability flag verbatim plus the arguments on both sides, and states the alternate estimate if dropped. Re-verified both AC baselines against the working tree today rather than trusting the archived numbers: README.md:60/62-65/77-80 and install.md:9 all check out, so the folded citations are accurate as written.
- DONE: Fold the relevant spike evidence into Risk evidence.
  Five findings plus the honest limit, under their own "README spike" heading with provenance and date: tap-spelling equivalence (identical `brew tap-info` output, same HEAD), both cask names and the `conflicts_with` metadata read from the tap's shipped `.rb` files, the live edge binary printing `Channel: edge` (which is why neither diff carries plugin commands), the channel-mapping tests, and `SPACEDOCK_CHANNEL` in install.sh. The "no `brew install` transaction was run" limit is carried forward unchanged rather than quietly dropped.
- DONE: Update Expected surface and tolerance to one combined net estimate, the Test plan, and Out of scope.
  Combined: net +167 across 8 files, tolerance ±55 / ±1 file, with per-half subtotals (+130/6 install, +37/2 README) and the alternate approved baseline (+137/7) if AC-5 is dropped. Test plan gains the AC-4 one-off check, the AC-5 arm with its red-then-green obligation called out as the task's only sequencing constraint, and a cross-half interaction note. Out of scope drops the superseded cross-reference and picks up the archived task's exclusions (marketplace repo README, install.sh); rb0 kept.

### Amendment summary

Two consequences of the fold that the re-gate should see, neither of which is visible
from the diff. First, **the validation posture changed for the README half**: the
superseded task correctly claimed the small-change fast path, and that claim no longer
holds — the combined task touches the front-door launcher, so the whole deliverable
carries full validation and the detached adversarial audit. That is an accepted cost of
"no doc-only task", not an oversight, and it raises the checking cost of the doc change
without changing the work. Second, **`mkdocs.yml` does not include README.md**, so the
docs strict build guards install.md and not the product README; AC-5's drift arm is the
README's only mechanical guard. That sharpens the AC-5 necessity question rather than
settling it — dropping the arm leaves the README's tokens unguarded again, which is the
gap Problem finding 4 documents.

I did not re-litigate the folded design. The one thing I checked rather than assumed was
the cross-half interaction: install.md is edited by one half and parsed by the other's
drift arm, and `installMDSection` stops at the next `## ` heading, so the exclusivity
paragraph in `## Skills` is outside the `macOS (Homebrew)` tab it reads. No collision,
and the test plan says to confirm it in the ordering rather than trust this note.

The doctor cut is unchanged and remains the one open thread from round 1.

## Stage Report: implementation

- DONE: The tolerated sibling-uninstall step at position 2 of installArgvSequence with its updated sequence doc comment.
  `host_exec.go:417` — `{argv: []string{"plugin", "uninstall", "spacedock@" + otherChannelMarketplace(devBranch)}, tolerateExit: true}`; comment now says 6-command, four tolerated cleanup steps, and carries the sibling/no-cascade rationale (commit b56e218a0).
- DONE: The literal-argv test extensions where the sibling id is a per-channel LITERAL that can diverge from otherChannelMarketplace.
  `TestInstallArgvSequence` spells `spacedock@spacedock` in wantEdge and `spacedock@spacedock-edge` in wantStable; `TestClaudeChannelInstallArgvSequence` gains a `wantSibling` table field with both literals. Deriving the sibling from the selected channel (so stable uninstalls itself) fails both. The three `TestInstallTolerates…` marker lists each gained the new step.
- DONE: The new claude-on-PATH value test with the two-channel local marketplace fixture (three subtests: sibling-present, fresh-box, reverse).
  `internal/cli/sibling_channel_exclusive_test.go` — `TestClaudeInstallLeavesOnlySelectedChannelEnabled` drives `execHost{}.Install` against two local directory marketplaces named `spacedock` / `spacedock-edge`, counting enabled `spacedock@*` entries in claude's own `plugin list --json`.
- DONE: The install.md exclusivity paragraph.
  `docs/site/get-started/install.md` +7, inserted after the "Codex installs the same way" paragraph as designed.
- DONE: The README posture (b) diff (stable label, install-guide pointer, tap token normalized).
  README.md +9/-2, exactly the designed diff.
- DONE: The AC-5 contractlint drift arm (captain kept it).
  `TestInstallHintProductReadmeNoDrift` extracts tap/formula from install.md's Homebrew tab and compares them to the README's `## Install` section.
- DONE: Proof run — the AC-5 drift arm recorded RED against the pre-change README first, then green after the README edit.
  RED: `install_hint_drift_test.go:121: README.md's Install section names tap "spacedock-dev/homebrew-tap"; install.md's Homebrew tab names "spacedock-dev/tap"` → `--- FAIL: TestInstallHintProductReadmeNoDrift`. Green after the edit. Re-proved RED a second time against the FINAL arm implementation (by reverting only the README token), because the first RED predated the dedup refactor below.
- DONE: Proof run — the value test run locally with claude on PATH, output pasted (a SKIP is a failure).
  `go test ./internal/cli -run TestClaudeInstallLeavesOnlySelectedChannelEnabled -count=1 -v` against claude 2.1.226: `--- PASS (9.47s)` with `sibling-present (3.64s)`, `fresh-box (2.28s)`, `reverse (3.55s)` all PASS. No SKIP. Falsification RUN, not asserted: deleting the sibling step made it report `enabled spacedock plugins after install = [spacedock@spacedock spacedock@spacedock-edge]` for sibling-present and reverse (fresh-box still passed, as predicted); step restored and re-verified green.
- DONE: Proof run — go test ./... and go test ./... -race green; gofmt -w ./cmd ./internal clean.
  Both run on the FINAL tree after the refactor: `go test ./... exit=0`; `go test ./... -race` exit=0 (all 21 packages ok). gofmt reports nothing for every file in the diff. One pre-existing unformatted file, `internal/release/runtime_live_evidence_workflow_test.go`, is also unformatted at HEAD and outside this surface — deliberately left alone rather than swept into this diff.
- DONE: Surface stays within the approved baseline net +167 across 8 files (tolerance ±55 LOC / ±1 file) with only the declared semantic changes.
  Final: net **+212 across 8 files** (+45 vs baseline, inside the ±55 band; file count exactly on baseline). First measurement came in at +244 (22 past the band) — see Summary for what changed and why it was not absorbed silently.
- DONE: Deliverable committed on the worktree branch with a summary of what was produced and where.
  Commit `b56e218a0` on `spacedock-ensign/claude-install-sibling-channel-cleanup`.

### Summary

The 6-line production fix landed as designed; every spike prediction held against
the live host. The value test earns its keep: with the sibling step deleted it
counts 2 enabled plugins in both sibling directions and 1 on a fresh box, which is
the captain's bug reproduced and the exact discrimination AC-1 asked for — so the
green result is not a tautology. The install.md edit does not disturb AC-5's arm
(confirmed by running the arm after that edit, per the cross-half note), because
`installMDSection` stops at the next `## ` heading.

Two deviations worth the gate's attention. First, the initial surface was +244
across 8 files, 22 LOC past the approved band. Rather than absorb it quietly or
trim comments to fit, I removed real duplication: the brew tap/formula scan is now
one `brewTokens` helper shared by the existing arm 2 and the new arm, the README
section scan reuses the package's existing `markdownSectionFromText`, and a
single-use fixture wrapper was inlined. That brought it to +212 — inside tolerance
on the merits, not by golfing. Second, because that refactor changed the arm after
its first RED run, I re-ran the red-then-green sequence against the shipped
implementation; a red recorded against code that was later rewritten would not
satisfy AC-5's obligation.

Not done here, flagged for validation: `mkdocs build --strict` — mkdocs is not on
this machine, and per the no-hidden-machine-dependencies rule I am not claiming a
lane I did not run. The install.md change is a plain paragraph with inline code
spans and no new links or nav entries. The doctor follow-up remains unfiled and is
still the one open thread.

## Stage Report: validation

- DONE: Every AC (AC-1 through AC-5) reproduced with cited evidence, not trusted from the implementation report: run the local claude-on-PATH value test yourself, reproduce AC-5's red-then-green claim against the shipped arm, run go test ./... and -race, and run mkdocs build --strict.
  AC-1: `go test ./internal/cli -run TestClaudeInstallLeavesOnlySelectedChannelEnabled -count=1 -v` vs claude 2.1.226 — PASS 7.77s (sibling-present 2.92s, fresh-box 2.01s, reverse 2.84s), no SKIP.
  AC-2: argv/tolerance tests green in the fresh suite; the audit's adversarial edits below prove they discriminate, which is what makes the green non-tautological.
  AC-3: install.md exclusivity paragraph present as designed; `uvx --with-requirements docs/requirements.txt mkdocs build --strict` (CI's pinned deps, docs.yml:39-42) exit 0 — closes the lane implementation flagged NOT run.
  AC-4: one-off existence check vs the stated baseline — README's 1 of 1 install blocks is now introduced by "Install the stable channel with Homebrew:" and the same-paragraph pointer names the edge channel and the install guide; the other bash block (line 77) is the launch command, not an install block.
  AC-5: red-then-green reproduced against the SHIPPED arm on the throwaway checkout — reverting only README's tap token yields `install_hint_drift_test.go:121: README.md's Install section names tap "spacedock-dev/homebrew-tap"; install.md's Homebrew tab names "spacedock-dev/tap"` FAIL; shipped tree green; drifting install.md's side instead also reddens the arm, so both files guard.
  Suites: `go test ./... -count=1` and `go test ./... -race -count=1` both exit 0 (20 packages ok each, zero FAIL lines; cli 257s with live tests). gofmt flags only pre-existing `internal/release/runtime_live_evidence_workflow_test.go`, outside this diff.
- DONE: Detached adversarial audit on a THROWAWAY checkout, never the implementation worktree: the four-surface front-door trigger fires (host_exec.go is launcher auto-heal machinery) and the AC-provenance trigger fires; construct adversarial edits the tests should catch and confirm they do; findings enter Review-finding disposition; "refuted nothing material" is a valid recorded outcome.
  Fresh clone of b56e218a0 in the session scratchpad. Front-door confirmed: frontdoor.go:350 — resolveHealableGate → ops.Install → installArgvSequence. Provenance confirmed: sibling ids are per-channel literals (wantStable "spacedock@spacedock-edge", wantEdge "spacedock@spacedock", wantSibling table field), never derived from otherChannelMarketplace.
  Adversarial edits, each caught then restored: (1) delete the sibling step — all 6 CI-lane tests FAIL and the live value test counts 2 (`[spacedock@spacedock spacedock@spacedock-edge]`) in sibling-present and reverse with fresh-box still passing: the captain's bug reproduced; (2) derive the sibling from the SELECTED channel — both argv tests FAIL; (3) flip tolerateExit to false — both argv tests FAIL; (4) misspell to the retired route-A id — both argv tests FAIL; (5) AC-5 red in both drift directions as above.
  Findings: refuted nothing material. One transient observation (codex's sibling remove appearing fail-fast) was an artifact of audit edit 3's own pattern touching host_exec.go:453 in the throwaway; at HEAD codexInstallArgvSequence's sibling remove is tolerateExit: true — no finding.
- DONE: Semantic adversarial pass scaled to the diff plus surface check: measure git diff --numstat against the approved baseline net +167 across 8 files (tolerance ±55/±1); verify only the declared semantic changes; recommend PASSED or REJECTED with material findings, deferred risks, and polish listed separately.
  `git diff --numstat "$(git merge-base main HEAD)"..HEAD` = 243 insertions / 31 deletions = net +212 across 8 files — inside the band, file count on baseline, matching the implementation's figure. All 8 file diffs read; only the declared semantic changes are present: the one tolerated sibling-uninstall step (plus its output line), comment and test extensions, the two doc diffs, and the drift arm with the behavior-preserving brewTokens dedup of arm 2 (same TrimSpace/last-match logic as the inline loop it replaced). The value test asserts cardinality AND exact id from claude's own `plugin list --json` with a seed-state guard, so it cannot pass on a wrong id or a silently failed fixture.

### Recommendation: PASSED

Material findings: none. Polish: none.

Deferred risks:
1. The value test skips in CI (no claude on the build runner) — the settled posture in the test plan; AC-2's argv tests guard the sequence in CI. Promote to material if validation practice drops the mandatory local run while the install sequence changes, or a CI runner gains claude and the test fails there.
2. brewTokens takes the LAST uncommented `brew tap`/`brew install` line per section (pre-existing arm-2 behavior, pure dedup). Today each guarded section has exactly one uncommented line (the edge cask is commented). Promote to material if either section gains a second uncommented brew line whose token can diverge from what the arm compares.

### Summary

All five ACs reproduced first-hand: the live value test passes in all three directions with no SKIP, AC-5's arm went red-then-green against the shipped implementation in both drift directions, both full suites pass fresh (-count=1 and -race), and the strict docs build — the one lane implementation could not run — passes with CI's pinned requirements. The detached audit refuted nothing material: five adversarial edits were each caught by exactly the tests the entity names, including a live reproduction of the captain's two-plugins bug. Surface +212 across 8 files is inside the approved band with only the declared semantic changes. Recommend PASSED.

### Correction round 1 — ASD-STE100 prose (captain send-back of PR #760)

- DONE: Every code comment and user-facing doc line this diff adds or edits conforms to ASD-STE100 (strict mode, simple-english skill catalog, self-check run).
  Commit `2f91b6dd4`. Rewrote: the `installArgvSequence` sequence comment; the `TestInstallArgvSequence` doc comment and its tolerance-loop comment; the `TestClaudeChannelInstallArgvSequence` doc comment; the three `TestInstallTolerates…` doc comments; every comment in `sibling_channel_exclusive_test.go`; and the two new comments in `install_hint_drift_test.go`.
- DONE: Self-check run (the skill's four checks, all four searched over the added lines).
  Zero hits for `should`/`would`/`may`/`might`/`could`, present perfect, contractions, and semicolons. Zero sentences over 25 words, counting backticked, quoted, and parenthetical text as one word each per rules 8.5 and 8.6. Three mid-sentence `when` clauses were moved to condition-first form. One noun `guard` and one rotated verb removed: the check/verify/confirm/validate/ensure set collapses to `asserts` throughout.
- DONE: Behavior and doc meaning are unchanged; untouched pre-existing comments stay as they are.
  Comment text only — no identifier, command, assertion string, or quoted-output change. The README keeps the stable label, the edge pointer, and the `spacedock-dev/tap` token; install.md keeps the exclusivity statement. `go vet` clean.
- DONE: Affected proofs re-run green after the rewrite and recorded.
  `go test ./...` exit 0 (all 21 packages ok). `go test ./... -race` exit 0. `TestInstallHintNoDrift` and `TestInstallHintProductReadmeNoDrift` both PASS — the drift arm still passes after the README prose changed, because the rewrite did not touch the brew tap/formula tokens it reads. The AC-1 value test re-run live with claude 2.1.226: `--- PASS (6.23s)` with `sibling-present`, `fresh-box`, and `reverse` all PASS, no SKIP. gofmt reports nothing for every file in the diff.
- SKIPPED: docs strict build.
  mkdocs is not installed on this machine. Flagged for validation rather than claimed. The install.md change is prose inside an existing section, with no new link and no nav entry.
- DONE: The correction is a new commit on the same branch (b56e218a0 preserved beneath it); state committed.
  `2f91b6dd4` sits on top of `b56e218a0` on `spacedock-ensign/claude-install-sibling-channel-cleanup`, so PR #760 picks it up.

**Scope interpretation, for the gate to rule on.** The assignment says to convert the
comments this diff adds or edits, and not to convert pre-existing comments the diff
does not touch. I applied that at the COMMENT-BLOCK level: where my earlier commit
edited a line inside a doc comment, I converted that whole block. A sentence-level
conversion would leave one paragraph half in STE and half not, which is harder to
read and to validate than either extreme. Comment blocks my diff never touched are
untouched — `codexInstallArgvSequence`, `TestInstallHintNoDrift`, `installMDSection`,
and the codex tests all keep their original prose.

**Surface deviation, captain-visible.** Cumulative against `main` is now net **+238
across 8 files** (325 insertions, 87 deletions). The approved baseline is +167 with a
±55 tolerance, so the band is +112 to +222 and this is **16 net LOC past the top**.
The cause is the STE conversion itself: the candidate stood at +212, inside the band,
and the rewrite added +26. STE costs lines by construction — it splits long sentences,
removes semicolon-joined clauses, and turns the six-command sequence into a numbered
list. The approved baseline predates the 2026-08-25 prose ruling and never carried its
cost. I did not trim load-bearing rationale to fit a number that was set before the
standard existed, and I did not adjust the tolerance myself: only the captain changes
an approved tolerance (README `## Review-finding disposition`, step 5). Recorded here
for that decision.

### Correction round 2 — cut, then convert (captain: "cut things we don't need")

The captain's criticism was correct. Round 1 held content constant and converted it,
which grows text by construction. This round deleted first.

- DONE: Every comment and doc line this diff adds or edits carries a fact a maintainer needs at that site.
  Commit `4b5b4073a`. Deleted as restatements of adjacent code: the 6-item command list in `installArgvSequence` (the argv slice sits directly below it), the command order and id mapping in `TestInstallArgvSequence` (the test literals say it), the stub setup and marker lists in the three tolerance comments (the `writeClaudeStub` call and the want list say it), and the subtest narration in the value test (the case table names them). Deleted as body-resident rationale: probes 3, 4, and 9, the fixture design choice, and the duplicate copy of the sibling-ambiguity reason.
- DONE: Survivors stay strict ASD-STE100.
  What survived is what the code cannot show at that site: the `plugin marketplace remove` cascade, why the uninstalls precede the marketplace add, why cleanup steps are tolerated (two distinct non-zero shapes), the sibling ids being deliberate literals, the falsifying edit for AC-1, and the SKIP-is-a-failure rule. Self-check re-run over the final added lines: zero banned modals, zero present perfect, zero contractions, zero semicolons, zero prose sentences over 25 words, and two remaining mid-sentence `when` clauses moved to condition-first form.
- DONE: Cumulative net at or under the pre-STE +212, inside the approved band.
  **net +147 across 8 files** (234 insertions, 87 deletions) against `main`. That is 65 under the +212 target and 20 under the original +167 baseline, so the round-1 tolerance breach is closed rather than waived. Final numstat: `README.md` +7/-2, `docs/site/get-started/install.md` +7/-0, `internal/cli/channel_selection_test.go` +10/-10, `internal/cli/host_exec.go` +15/-24, `internal/cli/init_devbranch_test.go` +7/-23, `internal/cli/install_tolerance_test.go` +12/-18, `internal/cli/sibling_channel_exclusive_test.go` +128/-0, `internal/contractlint/install_hint_drift_test.go` +48/-10. `host_exec.go` is now net **-9**: the cut reached past my own round-1 text into the original comment it replaced.
- DONE: Behavior, identifiers, commands, and assertion strings unchanged.
  Comment and doc text only. `go vet` clean. The README keeps the stable label, the edge pointer, and the `spacedock-dev/tap` token. install.md keeps the exclusivity statement, with the wording trimmed and every fact intact.
- DONE: Proofs re-run green and recorded.
  `go test ./...` exit 0 (22 packages, no FAIL). `go test ./... -race` exit 0, no FAIL. Both drift arms PASS: `TestInstallHintNoDrift` and `TestInstallHintProductReadmeNoDrift`. The AC-1 live value test re-run with claude 2.1.226: `--- PASS (11.54s)`, `sibling-present` `fresh-box` `reverse` all PASS, no SKIP. gofmt reports nothing for any file in the diff.
- SKIPPED: docs strict build.
  mkdocs is still not installed on this machine. Flagged for validation, not claimed.
- DONE: New commit on the same branch atop 2f91b6dd4; state committed.
  `4b5b4073a` on `spacedock-ensign/claude-install-sibling-channel-cleanup`, above `2f91b6dd4` and `b56e218a0`.

**One judgment call worth the gate's eye.** Cutting to what the code cannot show
means `installArgvSequence` no longer narrates its six steps, and the three tolerance
tests no longer restate their stub setup. A reader now gets the step list from the
argv slice and the reason from the comment. I think that split is right, and it is
the captain's stated direction, but it is a real change in how much a reader gets
from the comment alone.
