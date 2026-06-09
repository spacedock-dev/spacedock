---
id: k6d5xtg9hrxjcajrqyxnfah4
title: Two-channel release (stable→main / edge spacedock@next→next) + per-channel devBranch stamp + next-publish
status: implementation
source: "FO OWED, carried from the 2026-06-08-01 + 2026-06-08-02 debriefs (captain-nodded to file 2026-06-08). z9 (codex-plugin-auto-install) + #311 (Claude auto-install) install the plugin from the shared devBranch; the 0.20.0 flip needs each released channel's binary to install ITS OWN channel's plugin. Flip release-mechanics — a prerequisite of pj (main-flip-0200-marketplace), not a 0198 task."
started: 2026-06-08T22:48:35Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-two-channel-release-devbranch-stamp
issue:
sprint: 0200-flip
group: flip-mechanics
sprint-readiness:
---

Make each released spacedock channel auto-install the plugin from its own channel: the stable binary installs the `main` plugin, the edge binary installs the `next` plugin. Today there is one channel and `devBranch` is hardcoded `next`, so after the flip a stable binary would still install the `next` plugin.

## Problem

The front door auto-installs the plugin from a single shared `devBranch` (`frontdoor.go:49`, today `"next"`). Both `z9` (Codex auto-install) and the existing `#311` (Claude auto-install) consume this same `devBranch`. Two coupled gaps surface at the 0.20.0 flip:

1. **Channel ↔ plugin-source binding.** When the flip publishes a stable release on `main`, the stable binary must install the `main` plugin — but `devBranch` is still `next`, so it would install the edge plugin. The channel-tracking note in the 0198 sprint records this as the qa/z9 ↔ flip dependency: `z9` is correct *as long as it uses `devBranch`*; the retarget is what makes it install the right channel.
2. **Two distribution channels.** There is one brew artifact and one release lane today. The flip wants a stable channel (binary on `main`, installs `main` plugin) AND an edge channel (`spacedock@next`, installs `next` plugin), each stamped with its own `devBranch`, plus a `next`-publish step so edge keeps flowing after stable goes to `main`.

## Riskiest-mechanism spike — channel ↔ plugin-source binding (DONE, 2026-06-08)

The design's whole premise is that `devBranch` is the single knob that picks which channel's plugin a released binary installs. That binding was exercised end-to-end before committing to the rest of the plan, using the EXISTING `SPACEDOCK_DEV_BRANCH` override (`cli.go:478`) so no new code was needed to prove it.

Throwaway Go spike (`internal/cli/spike_channel_binding_test.go`, since deleted) drove the real production path: `applyDevBranchOverride(["SPACEDOCK_DEV_BRANCH=main"])` — the exact call `newClaudeCommand` makes (`cli.go:168`) — then a fresh-HOME-equivalent no-plugin auto-install through `runClaude`, and OBSERVED the `@ref` the install seam issued (not a grep of the constant):

- channel `main` → install seam records branch `main` → `marketplace add spacedock-dev/spacedock@main`
- channel `next` → install seam records branch `next` → `marketplace add spacedock-dev/spacedock@next`

The recorded branch is exactly what `execHost.Install` feeds `installArgvSequence` (`host_exec.go:255`), so the observed value IS the production argv. Binding proven: `devBranch=main` ⇒ stable binary installs the `main` plugin; `devBranch=next` ⇒ edge binary installs the `next` plugin. The auto-install fan-out (`#311` claude, `z9` codex) is correct as-is — it reads `devBranch`; this task only sets the value per channel.

Second unknown also exercised: **can the build pin `devBranch` per channel?** Throwaway goreleaser scratch project confirmed `goreleaser build --snapshot` stamps two builds with distinct `-X …cli.devBranch=main` / `=next` ldflags in ONE invocation, each verified by running the built binary (`devBranch=main` / `devBranch=next`). Env-templated ldflags (`{{ .Env.X }}`) also work, but the two-builds-one-config shape is preferred (one tag-triggered `goreleaser release` produces both channels). goreleaser v2.16.0 present locally.

No further spike needed: every remaining piece (the version-stamp branch swap, the second cask, reusing `next-publish.yml`) composes already-proven goreleaser/git behavior.

## Proposed approach

The single mechanism is **set `devBranch` per channel at build time** and add an **edge artifact alongside the stable one**, reusing the existing tag-triggered release flow.

1. **Per-channel `devBranch` stamp (goreleaser two builds).** Replace the single hardcoded `-X …cli.devBranch=next` ldflag (`.goreleaser.yaml:36`) with two `builds:` entries — `spacedock-stable` stamping `devBranch=main`, `spacedock-edge` stamping `devBranch=next` — both also stamping `Version`. Two `archives:` (one per build id) and two `homebrew_casks:` (stable `spacedock`, edge `spacedock@next`) consume them. One `goreleaser release` on the `v*` tag produces BOTH channels. The source default stays `next` (`frontdoor.go:49`) so a `go install …@next` / `--plugin-dir` dev build is unaffected; only the released artifacts carry the per-channel stamp.

2. **Single-target `release.yml` version-stamp branch `next`→`main` (M2 finding — k6d owns this).** The "Stamp plugin manifests to the release version" step hardcodes `git switch next` (`release.yml:161`) / `git push origin next` (`release.yml:172`). `release.yml` triggers ONLY on `push: tags: ['v*']` (`release.yml:7-10`) — a tag push carries NO channel signal, and there is no edge `v*`-tag cut scheme anywhere in the repo. So post-flip exactly ONE thing reaches this step: a stable `v*` tag, which must stamp+push the plugin-manifest `version` to `main` (not the current hardcoded `next`, which would land the displayed version on the edge branch while the stable plugin on `main` stays un-stamped). This is a single-target swap, not a channel branch: change `next`→`main`. There is no "edge cut → next" arm because nothing cuts an edge tag; the edge channel's `version`-display + `plugin update` re-pull stay on `next-publish.yml`'s calendar bump (item 3), with NO involvement from this stamp step. This is the M2 carve-audit finding (pj's flip checklist:128) explicitly tagged "naturally k6d's two-channel mechanism" — **k6d owns it**, confirmed below.

3. **Two-channel brew (reuse, don't rebuild).** Stable cask `spacedock` (unchanged) plus edge cask `spacedock@next`, both emitted by goreleaser from their channel's archive. The edge **publish cadence** reuses the EXISTING `next-publish.yml` (`workflow_dispatch` → `bump-calendar` marketplace.json → push `next`) as-is — it is the on-demand edge re-pull lane and needs no change. k6d does not author a new edge-publish workflow.

4. **Codex channel knob (no marketplace.json safety net).** The codex install issues `plugin marketplace add <source> --ref <devBranch>` (`host_exec.go:275`) — the `--ref` value is the binary's `devBranch` ldflag, and there is NO `.codex-plugin/marketplace.json` calendar key to act as a parallel channel signal (claude has one; codex does not). So for codex the `devBranch` ldflag is the SOLE channel determinant: if it is wrong there is no secondary correction. The per-channel stamp (item 1) already covers codex because both hosts read the same `devBranch`; the codex fresh-HOME smoke (AC-d) is what independently proves the codex `--ref` resolves channel-correct.

## Scope boundary — k6d vs pj (explicit)

**k6d (this task) owns the binary-side channel mechanics:**
- the per-channel `devBranch` ldflag stamp (`.goreleaser.yaml` two builds: stable→`main`, edge→`next`);
- the single-target `release.yml` version-stamp branch swap `next`→`main` (the M2 finding in pj's flip checklist:128, confirmed k6d's here; `release.yml` is tag-only-triggered so there is exactly one stamp target — stable);
- two-channel brew (stable `spacedock` + edge `spacedock@next` casks);
- the edge publish lane = REUSE `next-publish.yml` unchanged.

**pj (the flip capstone) owns the marketplace-side flip:**
- `.claude-plugin/marketplace.json` `source.ref` `next`→`main`, in ONE commit with the paired `skills/integration/marketplace_manifest_test.go:67` edit (`ref=="next"`→`"main"`) that keeps the line green by construction (pj flip checklist:132);
- doc reconciliation (`AGENTS.md:28`, `docs/releasing.md` — DoD 9);
- the archive-tag trigger guard, e2e-gate headSha runbook decision, existing-user migration (DoD 5–10).

**The agreement invariant:** after the flip the marketplace `ref` (pj) and the stable binary's `devBranch` ldflag (k6d) must BOTH be `main`. k6d sets the binary side now (on `next`, lands pre-flip); pj sets the marketplace side at the flip. Neither alone is sufficient: a stable binary stamped `main` against a marketplace still pinned `next` would `marketplace add …@main` (binary wins its own install), but the marketplace's own `update` re-pull would still serve `next` until pj flips the ref. They are sequenced (k6d pre-flip, pj at-flip) precisely so both ends settle on `main`.

**Stable calendar-bump-on-main (pj flip checklist:142) is NOT k6d's.** The marketplace `version` calendar key (the `plugin update` re-pull key) lives in `.claude-plugin/marketplace.json` — the file pj owns. For a stable cut, `bump-calendar` already runs in the local release-prep commit (`docs/releasing.md:38`) that lands on `main` before the tag, so the stable channel's calendar key advances through the normal release-prep flow. A CI calendar-bump-on-`main` path, if wanted, is pj's marketplace-side concern, not k6d's binary-side stamp.

## Out of scope

- The 0.20.0 flip itself (`pj`) — this is its release-mechanics prerequisite, sequenced before the tag.
- The marketplace.json `ref` flip + its paired guard-test edit (`pj`) — see scope boundary above.
- `z9` / `#311` plugin-install behavior — they are correct against `devBranch`; this task only retargets/per-channels the value they read.
- Linux distribution (`v3`).

## Acceptance criteria

Each AC is an end-state property of the finished task with a "Verified by" naming a check OUTSIDE this entity body (Go test, `goreleaser --snapshot`, or live fresh-HOME smoke) that can fail. No AC is satisfied by a source-grep of the `devBranch` constant or by reading this body.

- **AC-a (per-channel devBranch stamp).** `goreleaser release` (or `--snapshot`) on a `v*` cut produces two darwin artifacts: a stable binary whose `devBranch` resolves to `main` and an edge binary whose `devBranch` resolves to `next`, each from the same tag in one invocation.
  - *Verified by:* `goreleaser build --snapshot --clean` over the real `.goreleaser.yaml`, then RUNNING each built binary and asserting its resolved channel — `spacedock --version` plus a fresh-HOME front-door dry-run whose auto-install issues `marketplace add …@main` (stable) / `…@next` (edge), observed in the install argv, NOT a grep. (The spike already proved `devBranch=main` ⇒ `@main` via `SPACEDOCK_DEV_BRANCH`; this AC binds that value to the BUILD stamp.)

- **AC-b (single-target release.yml version-stamp branch `next`→`main`).** The `release.yml` "Stamp plugin manifests to the release version" step stamps+pushes the plugin-manifest `version` to `main` (not the hardcoded `next`). `release.yml` is tag-only-triggered with no edge-tag scheme, so there is exactly ONE stamp target — the stable channel — and this is a single-target swap, not a per-channel branch. The end-state property: a stable `v*` cut lands the displayed plugin version on `main`, never on the edge branch.
  - *Verified by:* a NEW `internal/release` workflow-guard Go test asserting the **agreement invariant across three INDEPENDENT surfaces** — NOT by parsing the branch back out of `release.yml` (that would re-read the value the implementer wrote: the banned tautology shape). The three surfaces that can diverge independently: (1) the `release.yml` version-stamp step's `git switch`/`push` target, (2) `.goreleaser.yaml`'s stable-build `devBranch` ldflag, (3) `.claude-plugin/marketplace.json` `source.ref`. The check parses all three real artifacts and asserts they agree; because they are authored in three different files by three different changes, a drift in any one fails the check — it has an independent source of truth, so it is a real relationship test, not a tautology. **Sequencing (honest):** k6d lands pre-flip on `next`, BEFORE pj flips surface (3) `next`→`main`. So the k6d-landable half is the binary-side pair — (1) `release.yml` stamp target == (2) `.goreleaser.yaml` stable `devBranch` — which k6d makes green now (both `main`). The full tri-surface `== main` agreement goes green only when pj flips the marketplace ref; it is the executable form of the agreement invariant and co-gates the flip (it is RED-by-construction between k6d's merge and pj's flip, exactly like pj's paired `marketplace_manifest_test.go` edit). The guard is net-new — `workflow_exec_guard_test.go` covers only the journey-ledger separation + goreleaser ordering, nothing pins the stamp branch today.

- **AC-c (two-channel brew).** A `goreleaser` run emits a stable `spacedock` cask and an edge `spacedock@next` cask, each pinning its channel's archive url+sha. The edge publish cadence remains the UNCHANGED `next-publish.yml`.
  - *Verified by:* `goreleaser release --snapshot --clean` producing both cask artifacts in `dist/` (the generated `.rb` cask files for `spacedock` and `spacedock@next`), inspected for distinct names and per-channel archive ids — this is the functional proof (both casks really emit). The `next-publish.yml`-untouched check (`git diff` over k6d's change set) is a **NON-REGRESSION assertion** that k6d did not disturb the existing edge lane — it is NOT functional proof of edge flow. Post-flip edge behavior (the calendar re-pull keeping `@next` installers current) is verified under **pj's calendar AC**, not k6d's; k6d exercises nothing that drives the edge publish, so it makes no claim about ongoing edge flow.

- **AC-d (codex fresh-HOME channel smoke — sole-knob host).** A codex-channel binary stamped `devBranch=main` resolves the codex plugin from `--ref main` against a fresh HOME (and `devBranch=next` ⇒ `--ref next`), with NO `.codex-plugin/marketplace.json` calendar-key fallback — the ldflag is the only channel determinant.
  - *Verified by:* a Go test driving `runCodex` (or `runInit --host codex`) through the no-plugin auto-install with `devBranch` set per channel, observing the codex install argv carries `marketplace add … --ref main` / `--ref next` (the seam records branch; `codexInstallArgvSequence` composes `--ref`). Complemented by the AC-a built-binary smoke run a second time with `--host codex` against a fresh HOME. Observed in install argv, not a grep.

## Test plan

The riskiest mechanism (channel ↔ plugin-source binding) is ALREADY exercised (see Spike) — that determination is on the record, so implementation does not re-spike it; it binds the proven value to the build stamp and the four surfaces.

| AC | What proves it | Kind | Est. cost |
|----|----------------|------|-----------|
| a  | `goreleaser build --snapshot` over real `.goreleaser.yaml` + RUN each binary (`--version`, fresh-HOME front-door dry-run observing `@main`/`@next` install argv) | live build + binary smoke | medium (goreleaser ~secs; smoke ~mins) |
| b  | NEW `internal/release` workflow-guard test asserting tri-surface agreement (release.yml stamp target ↔ .goreleaser.yaml stable devBranch ↔ marketplace.json ref) — k6d-landable half is the binary-side pair == `main`; full `==main` tri-surface goes green at pj's flip | Go unit (multi-artifact parse) | low |
| c  | `goreleaser release --snapshot` emits both casks in `dist/` (functional proof); `git diff` shows `next-publish.yml` untouched (NON-REGRESSION only) | live build + diff check | low–medium |
| d  | Go test over `runCodex`/`runInit --host codex` no-plugin arm, observing `--ref main`/`--ref next` install argv; + AC-a binary smoke with `--host codex` | Go unit + live smoke | medium |

Sequencing for implementation: land AC-b first (cheapest, a failing guard test that pins the channel-aware branch — TDD red→green for the `release.yml` edit), then AC-a/c via the `.goreleaser.yaml` two-builds change verified by one `goreleaser --snapshot` (covers a + c artifacts together), then AC-d's codex argv test + the second smoke. The fresh-HOME smokes (a, d) are live front-door dry-runs against a temp `HOME`, observing plugin/install state — never a source-grep of `devBranch`.

Fixtures/tools: goreleaser v2.16.0 (present locally; CI uses `goreleaser-action@v6 ~> v2`). The fresh-HOME smoke needs no real `claude`/`codex` host install — the front-door auto-install argv is observable through the injectable `hostOps` seam (as the spike did), so the channel→argv proof is a hermetic Go test; a fully-live `brew install` + real-host install is the optional belt-and-suspenders end-to-end, gated on host availability (declare if it requires CI/host secrets).

## Stage Report: ideation

- DONE: Exercise the RISKIEST mechanism first and record it — prove the channel→plugin-source binding (binary whose devBranch resolves to `main` auto-installs the `main` plugin, observed in plugin/workflow state, never a grep). Used the existing `SPACEDOCK_DEV_BRANCH` override.
  Throwaway Go spike drove `applyDevBranchOverride(SPACEDOCK_DEV_BRANCH=main)` → fresh-HOME `runClaude` no-plugin auto-install → OBSERVED install seam issuing `marketplace add …@main` (and `@next` for next); the recorded branch is the production `installArgvSequence` argv. Spike since deleted; `internal/cli` 235/235 green after removal. Recorded in body as a fresh-HOME-observed binding, not a constant grep.
- DONE: Scope the build-time per-channel stamp as the remaining engineering (binding already proven by the override).
  Second throwaway spike confirmed `goreleaser build --snapshot` (v2.16.0) stamps two builds with distinct `-X …cli.devBranch=main`/`=next` ldflags in one invocation, each verified by RUNNING the binary; two-builds-one-config chosen over env-templating. Recorded in Proposed approach item 1.
- DONE: Produce build-ready ACs + test plan covering all four channel surfaces, each AC verified OUTSIDE the body (Go test / `goreleaser --snapshot` / live fresh-HOME smoke), never a source-grep — (a) per-channel devBranch stamp; (b) channel-aware release.yml version-stamp branch; (c) two-channel brew reusing next-publish.yml; (d) codex sole-knob fresh-HOME smoke.
  AC-a..d written with per-AC "Verified by" (goreleaser --snapshot + built-binary smoke; NEW internal/release workflow-guard test for release.yml branch; both-cask dist/ check + next-publish.yml diff; codex `--ref` argv Go test + smoke). Test-plan table + TDD sequencing (AC-b guard test red→green first). Confirmed b's guard is net-new — `workflow_exec_guard_test.go` pins only journey-ledger separation/ordering, not the stamp branch.
- DONE: State the explicit k6d/pj scope boundary (k6d: devBranch stamp + release.yml channel-aware version-stamp + two-channel brew + edge publish; pj: marketplace.json ref flip next→main + paired guard-test edit). Marketplace ref and stable binary devBranch must agree (both `main`).
  Wrote a "Scope boundary — k6d vs pj" section: k6d owns binary side (lands pre-flip on `next`), pj owns marketplace side (at the flip). Confirmed via pj flip checklist: M2 version-stamp (line 128) tagged k6d's; ref flip + `marketplace_manifest_test.go:67` edit (line 132) pj's; stable calendar-bump (line 142) is pj's marketplace-side, not k6d's (stable bump already rides `docs/releasing.md:38` release-prep). The agreement invariant (both ends → `main`, sequenced) recorded.

### Summary

Firmed the design for the per-channel `devBranch` stamp. The riskiest mechanism — channel→plugin-source binding — was proven cheaply via the existing `SPACEDOCK_DEV_BRANCH` override (no new spike code needed beyond a throwaway test), and the build-time stamp feasibility was confirmed with a goreleaser scratch project. The single engineering mechanism is two goreleaser builds (stable→`main`, edge→`next` ldflag) plus a channel-aware `release.yml` version-stamp branch and a second `spacedock@next` cask, with `next-publish.yml` reused unchanged as the edge lane. Key decisions: two-builds-one-config over env-templating; AC-b's release.yml branch guard is net-new (nothing pins it today); the k6d/pj boundary is explicit with the both-must-be-`main` agreement invariant, and stable calendar-bump is pj's not k6d's. No implementation done — build-ready design only, per the pre-flip mandate.

### Staff-review fold (2026-06-08)

- M1 (Material) — AC-b/M2 encoded a phantom channel-derived stamp. CONFIRMED against the repo: `release.yml` triggers ONLY on `push: tags: ['v*']` (`release.yml:7-10`); a tag push carries no channel signal and there is no edge `v*`-tag cut scheme. Folded: reframed Proposed-approach item 2 and AC-b to a SINGLE-TARGET swap `next`→`main` (one stamp target — stable); DELETED the "edge cut → next" arm; stated next's calendar re-pull stays on `next-publish.yml` with no stamp-step involvement. Made AC-b's guard NON-TAUTOLOGICAL: replaced "parse the branch back out of release.yml" (the banned re-read shape) with a tri-surface agreement check across three independently-authored artifacts (release.yml stamp target ↔ .goreleaser.yaml stable devBranch ↔ marketplace.json ref). Added the honest sequencing: k6d's pre-flip-landable half is the binary-side pair == `main`; the full tri-surface `==main` goes green at pj's flip (RED-by-construction in between, like pj's paired test edit).
- P2 (polish) — AC-c's "next-publish.yml untouched" mislabeled. Folded: labeled it a NON-REGRESSION assertion (not functional proof of edge flow); noted post-flip edge behavior is verified under pj's calendar AC, not k6d's; dropped any "edge keeps flowing after stable goes to main" claim from the body (k6d exercises nothing that drives edge publish). Retitled AC-c "two-channel brew". Test-plan row c updated to mark the cask emission as functional proof and the diff as non-regression only.
- P5 (polish) — release.yml line citation normalized. Folded: Proposed-approach item 2 now cites `release.yml:161` (`git switch next`) and `release.yml:172` (`git push origin next`) plus `:7-10` (the tag trigger); verified no stale `:156`/`:159` (step name / RELEASE_VERSION assignment) citations remain in the body.

## Stage Report: implementation

- DONE: `.goreleaser.yaml`: replace the single `cli.devBranch=next` ldflag with TWO builds — `spacedock-stable` (ldflag `cli.devBranch=main`) and `spacedock-edge` (ldflag `cli.devBranch=next`), both also stamping `Version` — plus two `archives:` and two `homebrew_casks:` (stable `spacedock`, edge `spacedock@next`); one `goreleaser release --snapshot --clean` emits BOTH channels' artifacts each with the correct per-channel devBranch, verified by RUNNING each built binary and observing its fresh-HOME front-door install argv (`@main` stable / `@next` edge) — AC-a, AC-c.
  `.goreleaser.yaml` now has builds `spacedock-stable`(devBranch=main)+`spacedock-edge`(devBranch=next) sharing a goos/goarch YAML anchor, archives `spacedock-stable`+`spacedock-edge`(`_edge` asset suffix), casks `spacedock`(ids spacedock-stable)+`spacedock@next`(ids spacedock-edge). `goreleaser check` valid; `goreleaser release --snapshot --clean` emitted 8 binaries, both archives, and both casks (`dist/homebrew/Casks/spacedock.rb` + `spacedock@next.rb`, distinct per-channel archive urls). Live built darwin/arm64 binaries RUN: stable front-door auto-install issued `marketplace add spacedock-dev/spacedock@main`, edge issued `…@next` (observed in fake-host-recorded argv, fresh HOME). Commit d740f6ad.
- DONE: `release.yml`: single-target swap of the "Stamp plugin manifests" step `git switch next`/`git push origin next` → `main` (release.yml is tag-only-triggered → exactly ONE stamp target, stable; NOT a per-channel branch), AND land a NEW GREEN `internal/release` guard asserting the binary-side pair agreement: release.yml stamp target == .goreleaser.yaml stable-build devBranch, both `main`. Do NOT assert marketplace.json ref == main (still `next` until pj — that assertion would be RED). `go test ./internal/release/` GREEN — AC-b.
  `release.yml` stamp step now `git fetch origin main`/`git switch main`/`git push origin main` (comment reframed to the single-target rationale). New `internal/release/channel_agreement_guard_test.go`: `TestStableChannelBinaryPairAgreesOnMain` (binary-side pair == main, parsed from two independent artifacts) GREEN, `TestEdgeChannelStampsNext` (channels don't collapse) GREEN, `TestTriSurfaceChannelAgreement` SKIPS while marketplace ref is `next` (pj's surface) — structured so pj removes the skip + extends green-by-construction at the flip. Wrote RED-first (failed: no spacedock-stable build); GREEN after the two edits. `go test ./internal/release/` 75/75 green.
- DONE: AC-d codex sole-knob smoke: a Go test driving `runCodex`/`runInit --host codex` through the no-plugin auto-install with `devBranch` set per channel, observing the install argv carries `marketplace add … --ref main` / `--ref next` (seam records branch; observed in argv, NOT a grep); plus the AC-a fresh-HOME front-door dry-run proving `@main`(stable)/`@next`(edge). `go test ./internal/release/ ./internal/cli/` green.
  New `internal/cli/codex_channel_smoke_test.go`: `TestCodexNoPluginAutoInstallChannelRef` drives real `runCodex` no-plugin auto-install with devBranch=main/next, observes the recorded seam branch AND confirms `codexInstallArgvSequence` threads it into `--ref <branch>`; `TestClaudeNoPluginAutoInstallChannelRef` is the claude `source@branch` analog (AC-a hermetic half). Complemented by the live built-binary codex smoke: stable → `--ref main`, edge → `--ref next`. `go test ./internal/release/ ./internal/cli/` 341/341 green.

### Summary

Bound the proven channel↔plugin-source binding to the BUILD stamp: two goreleaser builds set `cli.devBranch` per channel (stable→main, edge→next), feeding two archives and two casks (`spacedock` + `spacedock@next`), each cask pinning only its channel's archive id. `release.yml`'s stamp step swapped next→main as a single-target (tag-only trigger = one stable target). The new `internal/release` guard asserts the k6d-landable binary-side pair (release.yml stamp target == stable devBranch, both `main`) green from two independent artifacts, with the full tri-surface ==main check skipped until pj flips the marketplace ref — structured for pj to extend green-by-construction. AC-a/c/d proven by a live `goreleaser release --snapshot` emitting both channels plus RUNNING each built binary and observing its fresh-HOME front-door install argv (`@main`/`@next` claude, `--ref main`/`next` codex); AC-b/d also pinned by hermetic Go tests. Scope boundary held: `marketplace.json` ref, `marketplace_manifest_test.go`, `next-publish.yml`, and any new edge-publish workflow all untouched (pj's / out of scope). NOTE for validation: a PRE-EXISTING, unrelated failure in `internal/dispatch` `TestLauncherCommandFallsBackToPathWhenSpacedockBinUnsetEmptyOrUnusable/unset` fails identically on baseline (with my changes stashed) — root cause is the dispatch harness exporting an executable `SPACEDOCK_BIN` into the test's `os.Environ()`, so the "unset" subcase's launcher runs the real `spacedock` instead of the PATH fallback; outside this task's scope, surfaced not fixed.
