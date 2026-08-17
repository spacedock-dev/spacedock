---
id: 2dfpswjwbxez1km6439wfrsn
title: Make channel install non-destructive and give edge installs a working re-pull key
status: ideation
source: "Captain directive, CL 2026-08-17: \"we should simplify the marketplace, there should just be one way\", then the chosen shape: \"i think spacedock-edge@marketplace is better, because we now have other things like subspace in the same marketplace.\" Raised from the live cross-host channel investigation; verified against the published marketplace on both refs and against a route-A install cached on the captain's workstation."
started: 2026-08-17T19:38:45Z
completed:
verdict:
score:
worktree:
issue:
---

One task, end to end: channel installs stop destroying shared marketplace state, edge installs re-pull to a contract that matches their binary, dev codex launches stop paying a measured 40-second copy tax, and the stalled edge release line advances again. Round 1's route collapse stands unchanged underneath: one published id per channel, channel in the marketplace NAME, zero id changes in the binary.

## Problem — the measured failure chain

Five defects, each exercised live 2026-08-17 (hosts: claude 2.1.226, codex-cli 0.147.0; all probes in isolated homes — `CODEX_HOME` / `CLAUDE_CONFIG_DIR`, isolation itself proven with a marker-plugin install that left the real home untouched).

**1. The install sequences destroy shared marketplace state.** `installArgvSequence` step 2 is `plugin marketplace remove <channel-marketplace>`; on stable that names `spacedock`, the marketplace that also hosts `subspace` and `cargento`. Measured on claude 2.1.226: `plugin marketplace remove` with installed dependents succeeds and CASCADE-UNINSTALLS every plugin installed from it (probe: two plugins installed from a test marketplace; after remove, `claude plugin list` reports none). Codex cascades identically, and `codexInstallArgvSequence` removes BOTH channels' marketplaces — the captain's real codex config has `subspace@spacedock` installed, so any codex heal or `spacedock install --host codex` run wipes it.

**2. The documented re-pull key is dead twice over.** `release.yml` and `next-publish.yml` bump the calendar version in the IN-REPO `.claude-plugin/marketplace.json` on `next` — but live edge installs resolve the STANDALONE `spacedock-dev/marketplace` repo's `edge` branch, whose entry version froze at `0.0.2026062001` (2026-06-20; `docs/releasing.md:265` already records the dead field). Deeper, measured: the marketplace ENTRY version is INERT on both current hosts. `claude plugin update` keys on the plugin MANIFEST version at the entry's source ref — bumping the manifest at the source triggers a pull with the entry untouched, bumping the entry with the manifest untouched is a no-op ("already at the latest version"). Codex has no `plugin update` at all; `codex plugin add` re-clones unconditionally. So no calendar bump, in ANY file, can deliver a re-pull on current hosts — and no cross-repo credential is needed to fix this, because the live re-pull key is `next`'s `plugin.json` version, which same-repo automation already stamps.

**3. The edge release line is STALLED (discovered this round).** A manual bare stamp on main (`ff9bb4506`, 2026-08-10, "release: bump version to spacedock@0.27.0" — no automation writes that message) rode the `v0.27.0-pre4` tag into `next`: the manifest-tag gate deliberately skips `-pre` tags (release.yml:163), and the `-X theirs` reconcile copied the bare `0.27.0` manifest onto `next`. `ComparePreVersion` correctly ranks any `0.27.0-preN` BELOW bare `0.27.0`, so every subsequent `-pre` tag's edge-advance decision prints `skip` — proven live: `spacedock-release edge-advance-decision v0.27.0-pre7 <next-manifest>` → `tag v0.27.0-pre7 target edge version 0.27.0-pre7 vs next 0.27.0` → `skip`. Last reconcile on `next` is `dbb0675c4` (pre4); tags pre5–pre7 produced none. Drift measured: `next` is 99 commits behind main, 38 ahead. While stalled, the calendar-bump step (gated on `advance`) never ran either — the re-pull key's writer is parked behind the same skip.

**4. No automation syncs an existing edge install to its binary.** The version gate compares major.minor only, so an installed `0.27.0-pre1` plugin against a `0.27.0-pre7` binary is Compatible and the launcher never refreshes. Claude does not auto-update plugins (field evidence: the captain's claude edge install sat at pre1 from 2026-07-21 for four weeks of daily use across manifest movements that a version-diff auto-updater would have pulled). Measured consequence, this session: the first officer ran contract `0.27.0-pre1` against a `0.27.0-pre7` binary — 49 changed shared-core lines, 96 changed dispatch-core lines, and a gate-required skill (`fo-gate-lifecycle`) absent. Mid-session (2026-08-17T20:35Z) the real claude install moved to bare `0.27.0` @ next HEAD by an actor outside this worker (isolation-proven) — the exact stranded shape defect 3 publishes.

**5. Codex dev launches pay a ~40s per-launch copy tax (the captain's slowness, root-caused).** `spacedock codex` with `--plugin-dir` or an adjacent checkout runs `installCodexLocalPluginDir` on EVERY launch, before and regardless of the gate. The local marketplace symlinks the REPO ROOT as the plugin; `codex plugin add` copies the symlink target wholesale into the plugin cache and serves skills from that cache (proven: a marker skill added to the live checkout does not appear in `codex debug prompt-input` without a re-add). This checkout is 3.6 GB — 3.0 GB `.worktrees`, 91 MB `.git`, 65 MB `tmp/` — against a plugin surface under 1 MB (skills/ 800 KB). Measured with the exact six-step sequence: 29.0s fresh, 40.3s repeat (the repeat also deletes the prior 3.6 GB copy; sys-time dominated). The captain's real codex cache holds a 4.8 GB copy today. Refuted leads, also measured: the frozen entry version does NOT gate `TooOldPlugin` per launch (codex records the MANIFEST version — listing and cache dir show `0.27.0`, not `0.0.2026062001`; the gate reads the cached manifest → same minor → Compatible → no heal); `codex plugin list` at the gate costs 0.02–0.03s; the installed plugin adds ~0.02s to codex prompt assembly (0.12–0.14s vs 0.11–0.15s empty home). Steady-state codex startup outside the plugin-dir path is not measurably ours: a full real-home `codex exec` roundtrip is 6.2s including the model call.

## Spike results — name match and namespace derivation (round 1, RESOLVED, stands)

Exercised live on this workstation 2026-08-17, before any design was settled. Hosts: codex-cli 0.147.0, claude 2.1.226. All probes in isolated homes; real config restored after.

**(a) Codex enforces entry-name == manifest-name at add time — still true.** `go test ./internal/cli -run TestCodexEntryNameMustMatchPluginName` re-run against the real CLI: PASS. Independently re-proven in a fresh `CODEX_HOME` with a two-entry family marketplace in the captain's exact target shape.

**(b) BOTH hosts derive the skill namespace from the MANIFEST name, not the entry name.** codex: `codex debug prompt-input` flips `spacedock-edge:probe` → `spacedock:probe` when only the cached manifest name changes; no override field exists. claude: accepts entry != manifest, and a live headless session lists the skill under the MANIFEST name.

**Verdict: the entry-name shape (`spacedock-edge@spacedock`) is falsified on codex** — it renames every edge skill to `spacedock-edge:*`. The surviving collapse keeps the channel in the marketplace NAME (route B, `spacedock@spacedock-edge`) and retires the route-A entry. Route A was never installable on codex, so its retirement is claude-only.

## Measurements (round 2) — host semantics and timings

Every mechanism below rests on one of these live probes.

| # | Probe | claude 2.1.226 | codex 0.147.0 |
| --- | --- | --- | --- |
| 1 | `marketplace remove` with installed dependents | succeeds, cascade-uninstalls ALL dependents | same cascade |
| 2 | `marketplace add`, name already registered, same source | no-op, exit 0 ("already on disk"); git snapshot NOT refreshed | silent no-op, exit 0, 0.01s; snapshot NOT refreshed |
| 3 | `marketplace add`, same name, DIFFERENT source | succeeds — re-registers, replacing the source | refuses ("remove it before adding"), exit 0, old source kept |
| 4 | `marketplace update` / `upgrade <name>` | refreshes snapshot from recorded source; works on local-path sources | refreshes git snapshot (11.2s); on a local-path marketplace prints an error but exits 0 |
| 5 | `plugin update`, content moved, manifest version unchanged | no-op ("already at the latest version") | (no update verb) |
| 6 | `plugin update`, manifest version changed at source ref | pulls; SYNCS on any difference including downgrade (pre10→pre9); numeric prerelease ordering correct (pre9→pre10) | (no update verb) |
| 7 | `plugin update`, only the marketplace ENTRY version bumped | inert — no pull | entry version inert everywhere: listing/cache keyed on manifest version |
| 8 | `plugin install`/`add`, already installed | no-op, exit 0; does NOT refresh | RE-CLONES content in place — fresh content without a remove |
| 9 | uninstall + install cycle, same manifest version, content moved | re-clones; cache dir overwritten with fresh content | remove clears config+cache; add re-clones |

Timings (this workstation): codex network heal, full current six-step: 9.0s (`marketplace add` clone 0.7s, `plugin add` clone 5.6s). claude current four-step: 3.5s fresh, 3.3s repeat. Plugin-dir codex six-step against this checkout: 29.0s fresh / 40.3s repeat. `codex plugin list`: 0.02–0.03s. Codex prompt assembly: 0.11–0.15s empty home, 0.12–0.14s with the edge plugin, 0.58–0.69s on the loaded real home.

## Proposed approach

**Round 1's collapse stands as mechanism 1** (route-A entry retirement in the marketplace repo; channel stays in the marketplace NAME; the gate must still surface the directive reversal — see round 1 sections). New mechanisms:

**2. Non-destructive install sequences (both hosts).** No step ever removes a marketplace record. Before/after:

claude — `installArgvSequence`, from:

    plugin uninstall spacedock@<channel-mkt>        (tolerated)
    plugin marketplace remove <channel-mkt>          (tolerated)   ← cascade-uninstalls dependents on stable
    plugin marketplace add <source>                  (fail-fast)
    plugin install spacedock@<channel-mkt>           (fail-fast)

to:

    plugin uninstall spacedock@<channel-mkt>         (tolerated)
    plugin uninstall spacedock-edge@spacedock        (tolerated)   ← route-A migration, round 1
    plugin marketplace add <source>                  (fail-fast)   ← probe 3: natively re-pins a changed source
    plugin marketplace update <channel-mkt>          (tolerated)   ← probe 4: snapshot refresh, the non-destructive re-pin
    plugin install spacedock@<channel-mkt>           (fail-fast)   ← probe 9: unconditional fresh clone

codex — `codexInstallArgvSequence`, from six steps (remove both plugins, remove BOTH marketplaces, add, add) to:

    plugin remove spacedock@<other-channel-mkt>      (tolerated)   ← channel exclusivity, global skill namespace
    plugin remove spacedock@<channel-mkt>            (tolerated)   ← cache hygiene: keeps latestVersionDir single-dir
    plugin marketplace add <source>                  (fail-fast)
    plugin marketplace upgrade <channel-mkt>         (tolerated)   ← probe 4: git-snapshot re-pin; local sources error harmlessly at exit 0
    plugin add spacedock@<channel-mkt>               (fail-fast)   ← probe 8: re-clones content

The re-pin the deleted removes existed to force survives three ways, all probed: claude `marketplace add` replaces a changed source natively (probe 3); `marketplace update`/`upgrade` refresh the snapshot from the recorded source (probe 4); and the plugin content pull is unconditional on both hosts (probes 8/9) — it never depended on marketplace state beyond the entry's source URL/ref. The codex own-channel `plugin remove` stays for a measured reason: without it version dirs accumulate in the cache and `latestVersionDir`'s lexical fallback misorders `0.27.0-pre10` below `0.27.0-pre9`, so the gate could read a stale manifest. Known limitation, accepted and documented: a FUTURE marketplace-source migration on codex (probe 3 refusal) would need its own explicit remove+add step; this task changes no sources.

**3. Edge contract sync at the front door (the scope-4 decision: IN, as a non-blocking refresh).** The blocking gate is UNCHANGED (minor-exact, both directions, stable and edge). New: when `devBranch != "main"` and the gate returns Compatible, the front door compares full versions prerelease-aware (`release.ComparePreVersion` semantics; binary display version stripped of `+build`) and runs one best-effort `Install` before launching when the BINARY IS STRICTLY NEWER, or when the plugin version on the edge channel carries NO prerelease at all (a bare version on the edge line is illegal by construction — the same poison rule as mechanism 5 — and is exactly the stranded shape defects 3–4 published; ComparePreVersion alone would rank it "newer" and never heal it). Plugin-newer-within-the-minor does NOT refresh (a not-yet-upgraded binary must not re-pull ~9s per launch). `--no-install` skips the refresh. A refresh failure warns and launches anyway — no launch that succeeds today is newly blocked. `contract.Result` gains a `PluginVersion` field so the front door can compare without re-reading the manifest. Justification for edge-only: on stable the contract within a minor is interchangeable by policy (D1) and skew windows close at each release; on the edge line the contract IS the product and skew is unbounded — measured at four weeks/six prereleases this session. Known bounded cost, declared: a source-built binary stamped AHEAD of `next`'s manifest re-pulls per launch (~3.5s claude / ~9s codex) until the next `-pre` tag lands; source-build dev loops normally ride `--plugin-dir`, which bypasses the gate entirely.

**4. Codex plugin-dir staging (the 40s fix).** `WriteCodexLocalMarketplace` stops symlinking the repo root and instead stages a filtered copy of the plugin surface into the persistent marketplace dir — include list: `.codex-plugin/`, `.claude-plugin/`, `skills/`, `agents/`, `hooks/`, `hooks.json` (everything the manifests and hook config reference; <1 MB vs 3.6 GB). Each launch re-stages (cheap) and the sequence re-installs from the stage, preserving the dev-loop freshness the per-launch reinstall exists for — freshness MUST flow through the cache because codex serves from it (live-symlink marker probe). Refresh-skipping alternatives are rejected below.

**5. Edge-advance bare-version guard + natural repair.** `EdgeAdvanceDecision` treats a bare (no-prerelease) `nextVersion` as poisoned state and ADVANCES any `-pre` tag over it (with a stderr note), instead of ranking the tag below it. ~6 lines plus tests. Repair needs no manual restamp: the first post-merge `-pre` tag advances, and the `-X theirs` reconcile inherits that tag's correctly-stamped manifests (the manifest-tag gate guarantees the agreement on stamped lines; the guard covers the skipped-`-pre` hole that let `ff9bb4506` poison the line).

**6. Retire the dead calendar machinery.** Drop the `bump-calendar` invocation and its commit from release.yml's edge-advance step (the `git push origin edge-advance:next` stays); delete `next-publish.yml` (its only job was the calendar bump; an out-of-band edge re-pull is now "cut a `-pre` tag", which is automated end to end); delete the `bump-calendar` command and `release.BumpCalendarVersion` with their tests; rewrite `docs/releasing.md`'s re-pull prose to name the real key (the `next` manifest version + the launcher's edge sync). The in-repo `.claude-plugin/marketplace.json` FILE stays: it is a pinned transitional bridge (`channel_agreement_guard_test.go:245` hard-fails on its removal — released v0.20.0 binaries resolve installs from it); it simply stops being bumped.

Rollout order: merge this task to main → cut `v0.27.0-pre8` (the guard makes it advance over bare `0.27.0`; `next` restamps and reconciles; the workflow file at the tagged commit already carries the changes) → the edge cask bumps → every edge launcher heals its own install via mechanism 3 (the bare carve-out covers installs that grabbed `0.27.0` during the stall, including the captain's claude install as of 2026-08-17T20:35Z). The marketplace-repo route-A entry removal follows the binary, per round 1's Migration.

## Migration — retired route `spacedock-edge@spacedock` (round 1, stands)

**Failure mode, exercised.** Retiring the entry from the marketplace manifest strands an installed route-A holder: `✘ failed to load` at session start, `plugin update` exits 1 naming no replacement. The migration must be active, not prose.

**Mechanism (pinned).** The claude sequence's tolerated `plugin uninstall spacedock-edge@spacedock` step (mechanism 2) migrates route-A holders on any `spacedock claude` / `spacedock install --host claude` run; install.md and the marketplace READMEs carry the manual remedy (entry-level uninstall, never `marketplace remove` — now doubly so given probe 1). Rollout: binary first, entry removal second; in the window route A stays installable — harmless.

**AC-4 proof lane (constructible):** `claude plugin install spacedock-edge@spacedock --scope local` in a scratch project from the still-published entry, run the new sequence, assert the route-A record is gone and route B is enabled.

**Population.** Route A appears in no install instruction in this repo; only the marketplace-repo README instructs it. This workstation holds it only as a stale cache dir.

## Documentation diffs

Diffs 1–4 from round 1 stand unchanged: (1) the `spacedock-edge` entry removal from the default-ref `marketplace.json`; (2) the default-ref README rewrite (channels list, install block, retired-route note); (3) the `edge`-ref README wholesale replacement; (4) the install.md retired-route admonition. See the round-1 record in git history (`6ea87a967`) for the full texts; implementation applies them verbatim.

New this round — `docs/releasing.md`: replace the calendar-key sentences (lines ~190–192, ~245, ~251–253) with the measured mechanism. Before: "the marketplace calendar key is bumped (`spacedock-release bump-calendar`) so `claude plugin update` / `codex` re-pull". After: "`next`'s manifests inherit the tag's version through the reconcile — that manifest version is the re-pull key: `claude plugin update` syncs on manifest-version difference, and an edge `spacedock` launch refreshes the installed plugin when its binary is newer (or the installed version carries no prerelease). The marketplace entry's version field is display metadata on current hosts; no calendar bump exists." The "Dev-Only `next` Publishing" section drops the `next-publish` re-pull sentence (workflow deleted; out-of-band re-pull = cut a `-pre` tag). Note 265–268 (the marketplace repo's dead version field) is updated to say the field is inert on current hosts and cosmetic. `docs/site/get-started/install.md` gains only the round-1 admonition.

## Mechanisms considered (necessity)

- **`marketplace update`/`upgrade` as the re-pin (chosen)** vs keeping remove+add: remove cascade-uninstalls dependents (probe 1) — the exercised harm this task exists to stop. vs deleting the remove with no replacement: the snapshot then never refreshes (probe 2), reintroducing the stale-source no-op the remove was added to prevent.
- **Conditional source-compare remove** (read `marketplace list --json`, remove only on mismatch): rejected — claude re-pins a changed source natively on add (probe 3), codex source changes are migration events outside this task, and the extra read adds a hostOps surface for a case with zero current instances.
- **Filtered staging (chosen)** vs skip-the-reinstall-when-unchanged: codex serves skills from the CACHE copy, so skipping the reinstall freezes the dev loop (live-marker probe); vs `git archive` staging: loses uncommitted edits, the dev loop's point. Staging keeps freshness and cuts the copied bytes ~4000×.
- **Cross-repo calendar bump** (deploy key or PAT to push the standalone marketplace repo): DISSOLVED by measurement — the entry version is inert on both current hosts (probe 7), so the push would buy nothing. This answers the credential blocker concretely: no credential is needed; the live re-pull key (`next`'s manifest version) is already writable by same-repo automation.
- **Blocking gate tightening** (reclassify prerelease skew as TooOldPlugin): rejected — the heal-refuse path would newly BLOCK launches (a source binary ahead of `next` re-gates unequal forever); the non-blocking refresh delivers the same freshness with no new refusal state.
- **Tombstone entry** for route A (round 1): stands rejected — more mechanism for less collapse.
- **Deleting the in-repo `marketplace.json` file**: rejected here — a pinned bridge for released v0.20.0 binaries with a guard test; retiring its WRITER is this task's value, the file's end-of-life is a separate decision with its own compatibility clock.

## Out of scope

The prerelease asset-naming defect and artifact-level stamp guard — task `c9`, at its own gate; `.goreleaser.yaml` archive naming and `--version` channel reporting untouched (no coupling found: this task's release.yml edits are confined to the edge-advance job's bump step and the decision guard). Changing how `devBranch` selects a channel. Deleting the in-repo bridge `marketplace.json` (above). Codex startup costs outside the plugin surface (MCP servers, model roundtrip — measured bound 6.2s, not ours to fix). Repairing the captain's real-home installs by hand — the shipped launcher heals them on first launch.

## Expected surface and tolerance

This repo, ~14 files, roughly +330/−230 gross (net ~+100): `internal/cli/host_exec.go` (both sequences + other-channel helper), `internal/cli/codex_marketplace.go` (staging), `internal/cli/frontdoor.go` (edge sync), `internal/contract` (Result.PluginVersion), `internal/release/edge_advance_decision.go` (guard), `cmd/spacedock-release/main.go` + `internal/release/release.go` (bump-calendar deletion), `.github/workflows/release.yml` (−4), `.github/workflows/next-publish.yml` (deleted, −38), `docs/releasing.md`, plus test files for each (sequence fixtures, staging, gate-sync fakes, guard cases, and edits to the wiring tests that currently pin the calendar step). Tolerance: ±35% gross, ±3 files. Outside this repo, unchanged from round 1: the marketplace.json entry removal and both README rewrites.

Semantics changed, declared: host-command traffic of both install sequences (no `marketplace remove` ever; added `marketplace update`/`upgrade` and the route-A uninstall); codex plugin-dir cache content (filtered plugin surface instead of the whole checkout — observable: no `.git`/`.worktrees` in the cache); edge launches may run one extra Install (new subprocess traffic; bounded churn case declared above); release automation (edge-advance advances over a bare `next` version; calendar bump retired; `next-publish.yml` gone); published plugin ids (route A unpublished, round 1). NOT changed: the ids the binary constructs, the marketplace sources it adds, the skill namespace, the blocking gate's verdicts on both channels, command grammar, stored formats, authority boundaries.

## Acceptance criteria

AC-1 through AC-4 stand from round 1 (one published id per channel; binary-constructed ids resolve on both hosts; skill namespace unchanged; route-A holders migrated actively — each with its verification lane recorded above and in round 1).

**AC-5 — A channel install leaves every co-hosted plugin installed and loadable.**
Verified by: the probe-1 lane re-run against the NEW sequences on both hosts — a dependent plugin installed from the shared marketplace, then the full channel install; `plugin list` must still show the dependent installed/enabled. Baseline that moves the wrong way: the CURRENT sequences, which measurably cascade-uninstall it (probe 1).

**AC-6 — An edge install lands a contract matching its binary (the condition that failed this session).**
Verified by: a local-mirror lane — marketplace entry pointing at a local clone of spacedock.git with `next` parked at an older prerelease SHA; install (old manifest version); advance the mirror's `next` to a newer-stamped SHA; one front-door gate resolution with the newer binary; assert the installed plugin version EQUALS the binary version and the new-contract file content is present. Also covers the bare-version carve-out (park the mirror on `ff9bb4506`-shaped content). Baseline: this session's measured four-week pre1-vs-pre7 skew, which the current gate calls Compatible and never heals.

**AC-7 — The codex plugin-dir launch pass costs under 2 seconds against this checkout.**
Verified by: the same `/usr/bin/time` harness that measured the baseline, run on the new staging + sequence. Baseline that moves the wrong way: 29.0s fresh / 40.3s repeat, measured this round on this checkout.

**AC-8 — A `-pre` tag advances the edge line even over a bare `next` manifest version.**
Verified by: unit tests on `EdgeAdvanceDecision` (bare `nextVersion` + `-pre` tag → advance; all existing skip cases unchanged), and live by the first post-merge `-pre` tag producing a reconcile commit on `origin/next` whose manifests carry the tag's version. Baseline: `skip` printed today for v0.27.0-pre7, proven live.

## Test plan

- **Sequences (AC-5, AC-2, migration):** extend the install-sequence fixtures (`install_tolerance_test.go` and channel suites) to pin the exact new argv lists, tolerance flags, and ordering on both hosts, and to assert NO step spells `marketplace remove` (regression lock against reintroduction). Live lane: the isolated-home dependent-survival script from probe 1, both hosts. Fails if a remove returns, tolerance flips, or the route-A uninstall drops.
- **Edge sync (AC-6):** fake-hostOps unit tests — edge + binary strictly newer → exactly one Install then launch; plugin newer → no Install; equal after `+build` strip → no Install; bare plugin version on edge → Install; stable → never; `--no-install` → never; Install failure → launch proceeds with warning. Each fails if the trigger set or the never-blocks property changes. Live lane: the AC-6 local-mirror recipe (no network).
- **Staging (AC-7):** unit test that `WriteCodexLocalMarketplace` produces a real directory (not a symlink) containing exactly the include list and excluding `.git`/`.worktrees`; live timing via the AC-7 harness. Fails if the symlink returns or the include list silently grows the copy.
- **Guard (AC-8):** table-driven `EdgeAdvanceDecision` cases; edit the edge-advance wiring tests that currently pin the calendar-bump step to pin its absence and the preserved push. Fails if a bare `next` version skips again or the push is lost.
- **Regression guards kept:** `TestCodexEntryNameMustMatchPluginName` (why the entry-name shape cannot return); the channel-agreement bridge test (why the in-repo marketplace.json file stays).
- **AC-1/AC-3 lanes** unchanged from round 1 (published-ref fetch assertion; namespace listing on both hosts).

No new test infrastructure; every live lane's commands are the probe commands recorded this round, verbatim.

## Stage Report: ideation

- DONE: Run the codex name-match and namespace-derivation spike BEFORE settling the design, and record the result in the entity body — it can invalidate the captain's chosen entry-name shape outright, so designing first wastes the round.
  Ran first, recorded in "Spike results": name-match holds on codex 0.147.0 (real-CLI test re-run + fresh-CODEX_HOME probe in the exact target shape); BOTH hosts namespace skills by MANIFEST name (codex `debug prompt-input` flip test; claude live listing `sdprobe-manifest:probe`); no codex override field. Captain's shape falsified on codex; design reversed and flagged for the gate.
- DONE: Pin the migration story for installs on the retired route to something exercisable against a real host holding the old id, not asserted in prose.
  Failure mode exercised live (installed plugin goes `failed to load` at session start after entry retirement; `plugin update` exits 1 naming no replacement); mechanism = tolerated `plugin uninstall spacedock-edge@spacedock` in the claude install sequence + remedy text; AC-4 lane pinned to a constructible real route-A install (still published) on this workstation.
- DONE: Produce the concrete doc diff for every install instruction that names an edge plugin id.
  Four diffs in "Documentation diffs": marketplace.json entry removal, both marketplace-repo READMEs (the default README actively instructs the retired route and misstates the binary's edge id), install.md migration note; releasing.md and build-from-source.md verified accurate unchanged via grep + read.

### Summary

The gating spike falsified the captain's entry-name shape before any design work: codex enforces entry==manifest at add time and both hosts namespace skills by manifest name, so `spacedock-edge@spacedock` renames every edge skill to `spacedock-edge:*` on codex. The proposal reverses to deleting the duplicate route-A entry and keeping the channel in the marketplace name — zero id changes in the binary, one tolerated uninstall step as the exercised migration, four concrete doc diffs. The gate must present the directive reversal to the captain, with the catalog-scope worry answered by live evidence that marketplace sources coexist on one host.

## Stage Report: ideation (cycle 2)

- DONE: Test the captain's codex-slowness hypothesis by MEASUREMENT before designing anything for it — confirm or refute it with real timings, because a design built on an assumed cause wastes the round exactly as the entry-name shape would have.
  CONFIRMED with a mechanism none of the dispatch's three leads named: `installCodexLocalPluginDir` runs the full install per launch and codex copies the ROOT-SYMLINKED 3.6 GB checkout into its cache — measured 29.0s fresh / 40.3s repeat with the exact six-step sequence; the 4.8 GB copy sits in the captain's real cache. All three leads refuted by measurement: codex records the MANIFEST version (no per-launch TooOldPlugin from the frozen entry), the gate resolves Compatible on real state, and `codex plugin list` costs 0.02–0.03s.
- DONE: Fix the destructive marketplace removal WITHOUT losing the re-pin it exists to force — name the mechanism that keeps re-pinning working, since deleting the step naively reintroduces the stale-source no-op it was added to prevent.
  Mechanism named and probed: `marketplace update` (claude) / `marketplace upgrade` (codex) refresh the snapshot non-destructively; claude `marketplace add` natively re-pins a changed source; content pulls are unconditional (`plugin install` after uninstall re-clones even at an unchanged version; `codex plugin add` re-clones in place). The cascade harm the removes cause was itself measured first: both hosts cascade-uninstall every dependent (subspace/cargento on claude, subspace@spacedock on codex).
- DONE: Size the cross-repo credential blocker for the bump-calendar repoint concretely — a re-pull key fix that cannot push to the marketplace repo delivers nothing.
  Sized to ZERO: measured on both current hosts, the marketplace ENTRY version is inert (claude `plugin update` keys on the plugin MANIFEST version at the source ref; codex re-clones on add and displays the manifest version), so no marketplace-repo push — and no credential — buys anything. The live re-pull key is `next`'s manifest version, stamped by same-repo automation; the blocker dissolves into unstalling edge-advance (bare-version guard) plus the launcher's edge sync.

### Summary

Scope expanded from the route collapse to the whole install/re-pull value chain, and measurement reshaped every half of it. New discovery: the edge line has been STALLED since 2026-08-10 — a manual bare `0.27.0` stamp rode the pre4 tag into `next` through the manifest-tag gate's `-pre` carve-out, and every later `-pre` tag's edge-advance decision provably skips; this, not just the dead calendar file, is why edge installs froze. The design: non-destructive sequences on both hosts (no `marketplace remove` ever — its cascade-uninstall of co-hosted plugins was measured), a non-blocking edge contract sync at the front door (with a bare-version carve-out that heals the poisoned installs), filtered plugin-dir staging (40.3s → <2s target, the confirmed slowness), a bare-version guard in edge-advance, and retirement of the calendar machinery the measurements proved inert. The gate should still surface the round-1 directive reversal, plus one mid-session observation: the captain's real claude edge install moved to bare `0.27.0` at 2026-08-17T20:35Z by another actor (this worker's isolation proven via marker test).
