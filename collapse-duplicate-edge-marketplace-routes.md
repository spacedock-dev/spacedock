---
id: 2dfpswjwbxez1km6439wfrsn
title: Collapse the duplicate marketplace routes so the channel lives in the entry name
status: ideation
source: "Captain directive, CL 2026-08-17: \"we should simplify the marketplace, there should just be one way\", then the chosen shape: \"i think spacedock-edge@marketplace is better, because we now have other things like subspace in the same marketplace.\" Raised from the live cross-host channel investigation; verified against the published marketplace on both refs and against a route-A install cached on the captain's workstation."
started: 2026-08-17T19:38:45Z
completed:
verdict:
score:
worktree:
issue:
---

The edge plugin is reachable by two ids resolving identical content. Collapse to one id per channel. The captain's chosen carrier (channel in the marketplace ENTRY name) was falsified by the ideation spike on codex — see Spike results; the surviving design keeps the channel in the marketplace NAME and deletes the duplicate entry-name route.

## Problem

`spacedock-dev/marketplace` publishes a `marketplace.json` on two refs, and both expose an edge route:

| ref | marketplace `name` | entries |
| --- | --- | --- |
| default | `spacedock` | `spacedock` (→ spacedock.git ref `stable`), `spacedock-edge` (→ ref `next`), `subspace`, `cargento` |
| `edge` | `spacedock-edge` | `spacedock` (→ ref `next`) |

`spacedock-edge@spacedock` and `spacedock@spacedock-edge` deliver identical content from `spacedock.git` ref `next`. Verified live: the captain's workstation holds installs made both ways (`cache/spacedock/spacedock-edge/` and `cache/spacedock-edge/spacedock/`), and the reporting Linux host a third combination.

Two costs. Operators have no single correct answer to "how do I install edge", and which route they picked determines whether a later command recognizes the install. The binary can only express `spacedock@<marketplace>` because `channelEntry` returns a hardcoded `"spacedock"` (internal/cli/host_exec.go:223), so the route-A id is one it cannot construct or recognize.

The decisive argument against encoding the channel in the MARKETPLACE is catalog scope: the marketplace now carries `subspace` and `cargento`, while the edge-branch marketplace carries only `spacedock`. A per-channel marketplace either duplicates the whole catalog per channel or strands every non-spacedock plugin for edge users, and lets the edge catalog drift from stable silently because catalog content is versioned by git branch.

## Proposed approach

**This reverses the captain's chosen carrier — the gate must surface it as a decision.** The captain directed: retire the `edge`-branch marketplace, one marketplace hosting the family, channel in the entry name (`spacedock-edge@spacedock`). The spike (below, run first) falsified that shape on codex: codex requires entry name == manifest name at add time, and BOTH hosts derive the skill namespace from the manifest name — so the edge entry-name shape unavoidably renames every edge skill to `spacedock-edge:*` on codex, breaking every `spacedock:*` cross-reference in skills, agents, and dispatch prompts. The surviving collapse:

1. **Marketplace repo (`spacedock-dev/marketplace`), default ref: remove the `spacedock-edge` entry** (route A, `spacedock-edge@spacedock`, claude-only — codex has always rejected it). The default marketplace keeps `spacedock` (stable), `subspace`, `cargento`; the `edge` branch keeps its single `spacedock` entry tracking `next` (route B, `spacedock@spacedock-edge`). Each channel then resolves through exactly one published id, and both surviving ids are the ones the binary already constructs — zero id changes in the binary.
2. **Binary: add one tolerated cleanup step to the claude install sequence** — `plugin uninstall spacedock-edge@spacedock` — mirroring the codex sequence's both-channel exclusivity sweep, so any `spacedock claude` / `spacedock install --host claude` run migrates a route-A holder automatically (see Migration).
3. **Docs: retired-route migration note in install.md; both marketplace-repo READMEs rewritten** — the default-ref README actively instructs the retired route today and misstates the binary's edge id (see Documentation diffs).

The captain's catalog-scope argument for the entry-name shape is answered empirically: marketplace SOURCES are not exclusive on a host. This workstation registers BOTH marketplaces (`spacedock-dev/marketplace` and `@edge` — verified via `claude plugin marketplace list`) with subspace installed from the default one while spacedock runs edge. The edge branch is a one-entry channel pointer, not a catalog copy, so nothing is duplicated per channel and no non-spacedock plugin is stranded for edge users; the silent-drift surface is that single entry.

## Out of scope

The prerelease asset-naming defect and the artifact-level channel guard — separate entity. Changing how `devBranch` selects a channel.

Also surfaced during ideation, recorded for a separate entity, not changed here: the legacy in-repo `.claude-plugin/marketplace.json` (single `spacedock` entry pinned to spacedock.git ref `main` on both `main` and `next` — no edge route, so not part of this collapse) is the file `next-publish.yml` and the edge-advance job calendar-bump, while live edge installs track the STANDALONE repo's `edge` branch, whose entry version is stale at `0.0.2026062001`. The documented "re-pull key" may be bumping a file no current install reads.

## Spike results — name match and namespace derivation (RESOLVED, run before design)

Exercised live on this workstation 2026-08-17, before any design was settled. Hosts: codex-cli 0.147.0, claude 2.1.226 (both current). All probes in isolated homes (`CODEX_HOME`, and local-scope installs in a scratch project for claude); real config restored after.

**(a) Codex enforces entry-name == manifest-name at add time — still true.** `go test ./internal/cli -run TestCodexEntryNameMustMatchPluginName` re-run against the real CLI: PASS (`codex plugin add spacedock-edge@spacedock` with manifest name `spacedock` fails with the "does not match" error; the matching-name shape installs). Independently re-proven in a fresh `CODEX_HOME` with a two-entry family marketplace built in the captain's exact target shape: the `spacedock-edge` entry installs only once its `.codex-plugin/plugin.json` name is also `spacedock-edge`.

**(b) BOTH hosts derive the skill namespace from the MANIFEST name, not the entry name.**
- codex: with entry `spacedock-edge` + manifest `spacedock-edge` installed, `codex debug prompt-input` (renders the model-visible prompt without a model call) lists the probe skill as `spacedock-edge:probe`. Mutating the cached manifest name to `spacedock` (entry unchanged) flips it to `spacedock:probe` — manifest name is the namespace. No override escape hatch: manifest fields `namespace`, `skillNamespace`, `skillsNamespace` probed; none changed the rendered namespace.
- claude: the current CLI still ACCEPTS entry != manifest — probe entry `sdprobe-entry` with manifest `sdprobe-manifest` installed clean (re-confirms, on today's CLI, the falsification of host_exec.go's "claude by construction" claim). A live headless session then listed the skill as `sdprobe-manifest:probe` — manifest name, not entry name.

**Verdict: the entry-name shape is falsified on codex.** The name match forces the edge manifest name to `spacedock-edge`, and the manifest name IS the namespace, so every edge skill renames to `spacedock-edge:*` on codex. On claude alone the shape would survive (mismatch accepted, manifest stays `spacedock`, namespace stays `spacedock:*`), but a one-host shape fails AC-2/AC-3. This is the sink condition this section anticipated; the Proposed approach reverses accordingly.

Route A was never installable on codex (same name-match), so the retirement below is claude-only.

## Migration — retired route `spacedock-edge@spacedock`

**Failure mode, exercised.** On a live claude host holding a structurally identical mismatched install (local scope, scratch project), retiring the entry from the marketplace manifest and refreshing the snapshot makes the INSTALLED plugin stop loading at session start: `claude plugin list` shows `✘ failed to load — Error: Plugin sdprobe-entry not found in marketplace sdprobe`, the skill vanishes from a live session's listing, and `claude plugin update` exits 1 with `Plugin "sdprobe-entry" not found` — no replacement named. A naked entry removal therefore strands route-A holders in exactly AC-4's failure mode; the migration must be active, not prose.

**Mechanism (pinned).**
- The new binary's claude install sequence gains a tolerated `plugin uninstall spacedock-edge@spacedock` step ahead of the existing cleanup pair, so every `spacedock claude` / `spacedock install --host claude` self-heals a route-A holder (removes the retired record, pins route B). Uninstall defaults to user scope — the scope route-A holders installed at.
- install.md and the marketplace READMEs carry the manual remedy naming the replacement: `claude plugin uninstall spacedock-edge@spacedock`, then the standard edge install (see Documentation diffs). Route-A holders' marketplace record (`spacedock`) is shared with subspace/cargento, so the remedy is entry-level uninstall, never `marketplace remove`.
- Rollout order: binary cleanup ships first (edge line via `next`, stable at the next release), the marketplace entry removal second. In the window route A remains installable and resolving — harmless.

**AC-4 proof (exercisable now and at implementation).** The retired id is still published, so a real host state is constructible on demand: `claude plugin marketplace add spacedock-dev/marketplace && claude plugin install spacedock-edge@spacedock --scope local` in a scratch project, then run the new binary's install sequence and assert via `claude plugin list` that the route-A record is gone and `spacedock@spacedock-edge` is enabled.

**Population.** Route A appears in no install instruction in this repo (repo-wide grep); only the marketplace-repo README instructs it. This workstation holds route A only as a stale cache dir (`~/.claude/plugins/cache/spacedock/spacedock-edge/0.19.9/`, harmless, `claude plugin prune` territory), not an active install; the live edge install tracks `spacedock-dev/marketplace@edge` (marketplace record verified).

## Documentation diffs

Every install instruction naming an edge plugin id, swept by grep across this repo and fetched from the marketplace repo on both refs.

**1. `spacedock-dev/marketplace` (default ref), `.claude-plugin/marketplace.json` — the collapse itself:**

```diff
@@ after the "spacedock" (stable) entry @@
     },
-    {
-      "name": "spacedock-edge",
-      "source": {
-        "source": "url",
-        "url": "https://github.com/spacedock-dev/spacedock.git",
-        "ref": "next"
-      },
-      "description": "Turn directories of markdown files into structured workflows operated by AI agents (edge channel — tracks next HEAD)",
-      "version": "0.0.2026062001",
-      "category": "workflow"
-    },
     {
       "name": "subspace",
```

No automation writes this file — the marketplace repo has no workflows (verified via repo listing); `next-publish.yml`/edge-advance bump the legacy in-repo copy, which has no `spacedock-edge` entry (see Out of scope).

**2. `spacedock-dev/marketplace` (default ref), `README.md`** — currently instructs `claude plugin install spacedock-edge@spacedock   # edge` and claims "an edge binary installs `spacedock-edge@spacedock`" (false today: the binary installs `spacedock@spacedock-edge`). Replace the channels list, install block, and binary-selection sentence with:

> One channel of the spacedock plugin lives here, plus the rest of the plugin family:
>
> - **`spacedock`** (stable) — pinned to the moving `stable` release branch (`source.ref: stable`); advances when a release tag is cut.
>
> The edge channel is NOT an entry here. It is the `edge` branch of this repo — a marketplace named `spacedock-edge` whose single `spacedock` entry tracks `next` HEAD. The channel lives in the marketplace NAME; the entry name always equals the plugin's manifest name (`spacedock`), so the hosts' entry-name/manifest-name check passes and skills keep the `spacedock:*` namespace on both channels.
>
> The `spacedock` binary selects the channel by its `devBranch` stamp: a stable binary installs `spacedock@spacedock`, an edge binary `spacedock@spacedock-edge`. Install:
>
>     # stable
>     claude plugin marketplace add spacedock-dev/marketplace
>     claude plugin install spacedock@spacedock
>
>     # edge
>     claude plugin marketplace add spacedock-dev/marketplace@edge
>     claude plugin install spacedock@spacedock-edge
>
> Codex: the same shape with `codex plugin marketplace add` / `codex plugin add`.
>
> **Retired route.** Early edge installs used a `spacedock-edge` entry in this marketplace (`spacedock-edge@spacedock`). If `claude plugin list` shows it — including as `failed to load` — run `claude plugin uninstall spacedock-edge@spacedock` and install edge as above, or run `spacedock claude` from an edge binary, which migrates automatically.

(Subspace/cargento sections unchanged.)

**3. `spacedock-dev/marketplace` (`edge` ref), `README.md`** — currently a stale copy of the default README, instructing the retired route and never mentioning that THIS branch is the `spacedock-edge` marketplace. Replace wholesale with:

> # spacedock marketplace — edge channel
>
> The `edge` branch publishes the edge-channel marketplace: NAME `spacedock-edge`, a single `spacedock` entry tracking the spacedock repo's `next` HEAD. The channel lives in the marketplace name; the entry equals the plugin's manifest name (`spacedock`), so skills stay `spacedock:*`.
>
> Install:
>
>     claude plugin marketplace add spacedock-dev/marketplace@edge
>     claude plugin install spacedock@spacedock-edge
>
> Codex: the same shape with `codex plugin marketplace add` / `codex plugin add`. An edge `spacedock` binary installs this id itself.
>
> Stable and the rest of the plugin family (subspace, cargento) live in the default-branch marketplace (`spacedock-dev/marketplace`).

**4. This repo, `docs/site/get-started/install.md`** — the edge instructions (lines 50-61) already name only the surviving route and stay unchanged; insert the migration note after the channel-explanation paragraph ending "...Codex installs the same way with `codex plugin add`." (line 61):

```diff
 the `spacedock-edge` marketplace the edge entry resolves from. Codex installs the
 same way with `codex plugin add`.
+
+!!! note "Retired route: `spacedock-edge@spacedock`"
+    Early edge installs used a `spacedock-edge` entry in the default marketplace.
+    That entry is retired. If `claude plugin list` shows `spacedock-edge@spacedock`
+    — including as `failed to load` after a marketplace refresh — remove it and
+    reinstall through the edge marketplace:
+
+    ```bash
+    claude plugin uninstall spacedock-edge@spacedock
+    claude plugin marketplace add spacedock-dev/marketplace@edge
+    claude plugin install spacedock@spacedock-edge
+    ```
+
+    Running `spacedock claude` from an edge binary performs this cleanup
+    automatically. (The retired route was never installable on Codex.)
```

**Verified no-change:** `docs/releasing.md` (channel-in-marketplace-NAME model, already accurate), `docs/site/contributing/build-from-source.md` (names the channel, not an id, and defers to install.md), and all `internal/cli` comments describing the surviving shape. Repo-wide grep for `spacedock-edge@spacedock` / `install spacedock-edge` / `add spacedock-edge` over `*.md`, `*.go`, `*.json`, `*.sh`: no hits outside historical roadmap/debrief records, which stay as records.

## Mechanisms considered (necessity)

- **Tolerated retired-id uninstall step (chosen)** — serves AC-4 and makes AC-1's collapse real on operator hosts, not just in published data. Simplest alternative: do nothing beyond entry removal + docs — insufficient: the exercised failure mode is a silent `failed to load` at session start plus an update error naming no replacement.
- **Tombstone entry** (keep `spacedock-edge` pointing at a notice-plugin) — satisfies AC-4's "continues to resolve" arm without a binary change, but keeps a second published route alive indefinitely, requires hosting a tombstone artifact, and serves only operators who installed route A yet never run the spacedock binary (~zero: route A was never in this repo's docs). Rejected: more mechanism for less collapse.
- **`spacedock doctor` retired-id check** — redundant once the install sequence self-heals and docs carry the remedy. Rejected (scope).

## Expected surface and tolerance

This repo: ~25 net LOC across 3 files — `internal/cli/host_exec.go` (+1 install step and comment in `installArgvSequence`), `internal/cli` install-sequence tests, `docs/site/get-started/install.md` (+migration note). Tolerance: ±15 LOC, +1 file. Plus marketplace-repo changes outside this repository: the marketplace.json entry removal and both README rewrites above.

Semantics changed: the set of published plugin ids (`spacedock-edge@spacedock` unpublished), and the claude install sequence's command traffic (one added tolerated uninstall). NOT changed: the ids the binary constructs, the marketplace sources it adds, the skill namespace, command grammar, stored formats, authority.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - Each channel resolves through exactly one published plugin id.**
Verified by: fetching `marketplace.json` from every published ref and asserting no two (ref, entry) pairs resolve to the same `spacedock.git` ref. Fails while both `spacedock-edge@spacedock` and `spacedock@spacedock-edge` point at `next` — the current state.

**AC-2 - The id the binary constructs is the id the marketplace publishes, on every channel and every supported host.**
Verified by: a test asserting `channelPluginID` output for each devBranch value exists as a real (marketplace, entry) pair in the published marketplace data, plus an install exercised on both claude and codex. Fails if the surviving route is one a supported host cannot resolve.

**AC-3 - The skill and agent namespace is unchanged by the migration.**
Verified by: installing the surviving edge id on a real host and asserting the first-officer skill still resolves as `spacedock:first-officer`. Fails if the entry rename leaks into namespacing — the outcome that would invalidate this approach rather than merely complicate it.

**AC-4 - An operator already installed via the retired route is not silently left broken.**
Verified by: the documented migration path exercised against a host holding the retired id, showing it either continues to resolve or produces an actionable message naming the replacement. Fails if a retired route degrades to plugin-not-found with no remedy text.

## Test plan

- **AC-1:** at the implementation gate, fetch `marketplace.json` from every published ref of `spacedock-dev/marketplace` and assert no two (ref, entry) pairs resolve the same spacedock.git ref. Live check (network), not a CI unit; fails today (both routes resolve `next`), passes after the entry removal. Cost: one script invocation.
- **AC-2 / migration step:** extend the existing claude install-sequence tests (`install_tolerance_test.go` and the channel-selection suites) to assert `installArgvSequence` contains the tolerated `plugin uninstall spacedock-edge@spacedock` step ahead of the marketplace add on both channels. Fails if the step is dropped, reordered after pinning, or made fail-fast. Codex live-lane coverage already exists (`runtime-live-e2e.yml` name-match job exercises a real `codex plugin add` in the surviving shape).
- **AC-3:** after implementation install on both hosts, assert the skill listing still shows `spacedock:first-officer` — claude via a headless `-p` listing (mechanism proven in this spike), codex via `codex debug prompt-input` (no model call). Cost: minutes, no new fixtures.
- **AC-4:** the live migration lane pinned in Migration — construct a real route-A install at local scope from the still-published entry, run the new binary's install path, assert the retired record is gone and route B is enabled. Fails if the sequence leaves the retired record or the remedy text names a wrong command.
- **Regression guard:** `TestCodexEntryNameMustMatchPluginName` stays — it is the on-record proof of why the entry-name shape cannot return and why no codex migration exists.

No new test infrastructure; the spike's probe commands seed the live checks verbatim.

## Stage Report: ideation

- DONE: Run the codex name-match and namespace-derivation spike BEFORE settling the design, and record the result in the entity body — it can invalidate the captain's chosen entry-name shape outright, so designing first wastes the round.
  Ran first, recorded in "Spike results": name-match holds on codex 0.147.0 (real-CLI test re-run + fresh-CODEX_HOME probe in the exact target shape); BOTH hosts namespace skills by MANIFEST name (codex `debug prompt-input` flip test; claude live listing `sdprobe-manifest:probe`); no codex override field. Captain's shape falsified on codex; design reversed and flagged for the gate.
- DONE: Pin the migration story for installs on the retired route to something exercisable against a real host holding the old id, not asserted in prose.
  Failure mode exercised live (installed plugin goes `failed to load` at session start after entry retirement; `plugin update` exits 1 naming no replacement); mechanism = tolerated `plugin uninstall spacedock-edge@spacedock` in the claude install sequence + remedy text; AC-4 lane pinned to a constructible real route-A install (still published) on this workstation.
- DONE: Produce the concrete doc diff for every install instruction that names an edge plugin id.
  Four diffs in "Documentation diffs": marketplace.json entry removal, both marketplace-repo READMEs (the default README actively instructs the retired route and misstates the binary's edge id), install.md migration note; releasing.md and build-from-source.md verified accurate unchanged via grep + read.

### Summary

The gating spike falsified the captain's entry-name shape before any design work: codex enforces entry==manifest at add time and both hosts namespace skills by manifest name, so `spacedock-edge@spacedock` renames every edge skill to `spacedock-edge:*` on codex. The proposal reverses to deleting the duplicate route-A entry and keeping the channel in the marketplace name — zero id changes in the binary, one tolerated uninstall step as the exercised migration, four concrete doc diffs. The gate must present the directive reversal to the captain, with the catalog-scope worry answered by live evidence that marketplace sources coexist on one host.
