---
id: 7p4vxcb9edzrph0zsfn75cd9
title: claude install leaves the sibling edge plugin installed and enabled
status: ideation
source: "Captain CL work-machine report 2026-08-25: clean documented stable install on a machine with spacedock@spacedock-edge left both plugins installed and enabled; doctor OK; which plugin serves unpredictable"
started: 2026-08-25T14:36:08Z
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
        - id: gate:7p4vxcb9edzrph0zsfn75cd9:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:7p4vxcb9edzrph0zsfn75cd9-backlog-1
              briefing:
                id: briefing:7p4vxcb9edzrph0zsfn75cd9:backlog:attempt-1:revision-1
                digest: sha256:4b88d6e7f37061f48b8d0475ee4f5a15476c1728bf602afc4d29856f86b71fb9
                request-digest: sha256:7652910a297d84b88917d9d48ad7e4009f83582811873821de42860027c6a649
                room-ref: ./claude-install-sibling-channel-cleanup/review/backlog/briefing-1
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
                room-ref: ./claude-install-sibling-channel-cleanup/review/ideation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-25T15:09:29.517214Z"
                reason: 'Captain amendment 2026-08-25: fold the product README edge-channel deliverable into this task; design re-gates after the fold'
---

`spacedock install --host claude` never uninstalls the sibling channel's modern plugin id. A stable install on a machine that already has `spacedock@spacedock-edge` (or the reverse) leaves BOTH channel plugins installed and enabled. Both are entry `spacedock` shipping the same skill set, so which provider serves is unpredictable. Mirror the codex sequence's sibling-channel remove, and give doctor eyes on a co-installed enabled sibling.

## Problem

`installArgvSequence` (internal/cli/host_exec.go:406) uninstalls only the own-channel id and the retired route-A id (`spacedock-edge@spacedock`). The modern sibling id is never touched. `codexInstallArgvSequence` (host_exec.go:445) already removes the sibling channel first — justified by codex's global skill namespace — and claude has the same practical collision: both channel plugins are entry `spacedock` with identical skill names. Neither safety net catches the dual install: `spacedock doctor` resolves only the running binary's own channel manifest (channel-aware resolver) and exits 0; the binary gate passes because both plugin snapshots can declare the same required binary minor.

Ideation confirms all of this against the source and, for the parts that are runtime
behavior, against a live claude (risk evidence below). One correction to the seed: the
front door does not merely fail to catch the dual install, it cannot heal it — the
auto-heal fires only on NoPluginFound or TooOldPlugin, and a dual install whose
own-channel manifest is Compatible never reaches either, so an affected machine stays
affected until someone runs `spacedock install` explicitly.

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

### Documentation diff

`spacedock install --host claude` gains host-integration behavior a reader would want
stated, so ideation owns the doc change. One file, one paragraph, inserted in
`docs/site/get-started/install.md` immediately after the existing channel paragraph
that ends "Codex installs the same way with `codex plugin add`." (line 82):

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

## Out of scope

CLAUDE_CODE_PLUGIN_PREFER_HTTPS documentation (plugin-install-https-keyless-machines / rb0). Product README edge-channel mention (product-readme-edge-channel-mention). The codex install sequence (already removes the sibling).

Split out by this ideation (see "Doctor visibility" above): doctor/`--version`
visibility into a co-installed enabled sibling, and the overstated
`docs/site/reference/command-reference.md:42` claim that doctor reports per-host
enablement. Also out of scope: healing an existing dual install without an explicit
`spacedock install` run — the launcher's auto-heal fires only on NoPluginFound or
TooOldPlugin, and a dual install whose own-channel manifest is Compatible does not
trigger it. Changing that is the same front-door gate decision the doctor follow-up
owns.

## Expected surface and tolerance

Estimate net LOC change: **+130, across 6 files** (insertions ~+134, deletions ~-4).
Tolerance: **±40 net LOC and ±1 file**.

The backlog seed said +40/4 files. Ideation revises it upward: the production change is
6 lines; the rest is the value test's two-channel fixture, which nothing in the repo
builds yet.

| File | Net | What |
| --- | --- | --- |
| `internal/cli/host_exec.go` | +6 | the argv step (+1) and the sequence doc comment (5-command → 6-command, sibling rationale, "three cleanup steps" → four) |
| `internal/cli/init_devbranch_test.go` | +3 | `TestInstallArgvSequence`: one literal step in each of `wantEdge`/`wantStable`, comment |
| `internal/cli/channel_selection_test.go` | +4 | `TestClaudeChannelInstallArgvSequence`: `wantSibling` table field, one literal step, comment |
| `internal/cli/install_tolerance_test.go` | +3 | three existing tests each expect one more stub marker; "five steps" → "six steps" in comments |
| `internal/cli/sibling_channel_exclusive_test.go` (new) | +108 | the value test plus the two-channel local marketplace fixture builder |
| `docs/site/get-started/install.md` | +6 | the exclusivity paragraph |

Tolerance rationale: the fixture builder's final size is the only real unknown (the
spike built it in shell, not Go). The ±1 file covers putting the fixture helper beside
the existing one in `co_hosted_survival_test.go` instead of a new file.

### Declared semantic changes

- **Runtime behavior (intended, the point of the task):** `spacedock install --host
  claude` now uninstalls the sibling channel's plugin. A captain deliberately holding
  both channels on claude loses that arrangement on the next install or auto-heal.
  Codex has behaved this way since its sequence was written; this makes claude match.
- **Command output:** `spacedock install --host claude` gains claude's own line for the
  extra step (a success line when the sibling was present, a tolerated failure line
  when it was not). No Spacedock-authored output changes.
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
`otherChannelMarketplace`, not a re-derivation of it.

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

Verified by: the doc diff above applied to `docs/site/get-started/install.md`, reviewed
at the validation gate, with the docs strict build (`mkdocs build --strict`, the docs
lane) as the mechanical check. Stated honestly: no committed test proves prose says
what we wrote — a grep over it would be the banned prose-grep. This AC counts only as
a companion to AC-1, which is where the behavior is measured.

## Test plan

| Test | Proves | Cost | Lane |
| --- | --- | --- | --- |
| `TestClaudeInstallLeavesOnlySelectedChannelEnabled` (new, `internal/cli`) | AC-1 | ~108 LOC, ~15s, offline | local + any runner with claude on PATH; **skips** in CI `build` |
| `TestInstallArgvSequence` (extend) | AC-2 | +3 LOC, instant | CI `build` (`go test ./...`) |
| `TestClaudeChannelInstallArgvSequence` (extend) | AC-2 | +4 LOC, instant | CI `build` |
| three `TestInstallTolerates…` (extend marker lists) | no regression in `runInstallSteps` tolerance | +3 LOC, ~1s | CI `build` |
| `TestClaudeChannelInstallLeavesCoHostedPluginInstalled` (existing, unchanged) | the sibling step did not reintroduce cascade harm | 0 LOC | local, claude on PATH |
| docs strict build (existing) | AC-3 well-formed | 0 LOC | docs lane |

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

**Estimated total cost:** ~1 hour implementation, dominated by the fixture builder.

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}

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
