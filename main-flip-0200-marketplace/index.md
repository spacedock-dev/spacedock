---
id: pj6k0h83g6taszkmt92qvsc0
title: Main-flip milestone — cut 0.20.0 on main + flip the marketplace once 0.19.6 is tested
status: implementation
source: "captain (2026-06-05) — '0.19.6: flip main readiness — once tested i want to make 0.20.0 on main and flip the marketplace.' The capstone of the 0.19.6 line."
score: "0.38"
started: 2026-06-08T22:48:35Z
completed:
verdict:
worktree:
issue:
sprint: 0200-flip
group: flip-capstone
sprint-readiness: defer
---

The 0.19.6 capstone: once the line is tested/green, cut **0.20.0 on main** and **flip the marketplace** so users get the stable release off main instead of the `next` development branch. This is the "flip main readiness" milestone.

## Hard prerequisites (gate this milestone)

- **`release-gate-job-separation-fix` (bqqr) MUST land first** — without it the cut fails the Runtime-Live-E2E gate (like 0.19.5).
- **The 0.19.6 line ALL GREEN** — captain's "tested" gate = all lanes/scenarios green on the branch (gq merged, the pi-line settled, all live + offline green). No partial-green cut.
- **README reconciliation MUST land before the flip.** The root README and install-facing docs must stop describing the pre-flip `next`-only world as the stable path, and must clearly describe the post-flip stable `main` install path plus the retained dev-only `next` path.
- **Upgrade-path confirmation MUST pass before the flip.** Exercise the user journeys that can hit the stable transition:
  - A user with an old `0.12.x` host plugin but no usable `spacedock` binary, where the plugin is refreshed/auto-updated and the first useful instruction is the binary install + `spacedock claude` payoff.
  - A user with a `0.19.x` install whose plugin is newer than the binary or whose binary is outdated for the current plugin surface; the gate/remedy must fail early with an actionable upgrade path rather than failing later in dispatch.

## Direction (captain-clarified 2026-06-05)

- **The flip = `next` becomes `origin/main`.** Make `main` the release branch carrying what `next` holds; the marketplace serves from `main` (the stable channel) instead of `next`. This is the "flip main readiness."
- **Version = 0.20.0, the FIRST actual release.** NO contract change — the contract stays at 1; 0.20.0 is the first real (non-dev) release, not a contract bump.
- **Bundling is OUT — figure it out later.** The `44 bundle-asset-distribution` (plugin-into-binary, --plugin-dir) path is NOT part of this milestone; the captain explicitly deferred it. This milestone is the branch/marketplace flip + the cut, not the packaging mechanism.
- Note the related stale-ref bug `s0cq install-marketplace-ref-refresh` — a clean marketplace flip likely needs that fix so the new `main` ref actually replaces the old `next` pin (don't let the flip no-op on a stale ref).

## Branch mechanics (captain-clarified 2026-06-06)

- **Archive current `main` first.** Before replacing it, tag the current `origin/main` tip as `v0-archived` so the pre-v1 history remains reachable by a named ref.
- **Replace `main` with `next`.** The actual flip is a deliberate non-fast-forward update: force-push the prepared `next` tip (or the 0.20.0 release-prep commit based on it) to `origin/main`.
- **Keep `next`.** Do not delete `next`; keep it as the dev-only release/publish channel after the stable `main` lane exists.
- **Stable release mechanics move to `main`; dev mechanics stay on `next`.** The `main` marketplace/ref, release post-stamp target, and released binary install pin must serve from `main`. The `next` branch can retain a dev-only publish path for pre-stable testing.

## Out of scope

The release-gate machinery fix itself (bqqr). The packaging/notarization tasks (44, 5w) unless they're prerequisites the captain names.

## Riskiest-mechanism determination — no spike needed

The flip's risk is **sequencing/runbook discipline**, not an unproven mechanism. Every step it relies on is already proven; the determination is recorded here rather than left silent:

- **Marketplace ref re-pin actually replaces the stale pin** — `s0cq` (install-marketplace-ref-refresh) PASSED 2026-06-03. A clean ref flip `next`→`main` re-pins; it does not no-op on the old pin.
- **The e2e-gate distinguishes a fully-green live run from a parked/offline-only run, bound to the release commit** — `nzb`'s spike (2026-06-08) exercised `gh run list --workflow "Runtime Live E2E" --status success` + `--json headSha` against the live `spacedock-dev/spacedock` repo: a parked run is never `success`; the query is the same one the release runner already uses; `headSha` binds a green run to a commit. Banked, do not re-spike.
- **Guarded non-fast-forward branch replacement** — standard git `--force-with-lease` (CAS on the expected old tip). The pre-flip `origin/main` (`8c069d95` at ideation) and `origin/next` (`d13f355e`) are divergent, so the flip is a real non-ff update, exactly the case `--force-with-lease` guards.

**E2E-gate headSha path for the cut — DECIDED: tag-the-green-tip.** `runtime-live-e2e.yml` carries `workflow_dispatch` (confirmed at `runtime-live-e2e.yml` `on:`), so the live matrix can be dispatched against an arbitrary ref. The cut runs the live e2e via `workflow_dispatch` on the *exact prepared `main` tip*, waits for it green, then tags **that same commit** `v0.20.0`. `release.yml`'s `e2e-gate` (nzb) resolves the tagged commit SHA and finds the matching green run → goreleaser proceeds. The channel-aware stamp (`k6d`) commits its stamp AFTER the tag fires (it is a step inside the tag-triggered `release.yml`), so no new commit lands between green-run and tag — the tagged SHA and the green-run SHA are identical by construction. The `SPACEDOCK_E2E_GATE_WAIVER` (recorded reason) remains the captain's auditable escape hatch if the prepared tip cannot get a green live run in time, but it is the fallback, not the plan.

**Boundary with the pre-flip pair (gates this milestone).** Two release-machinery fixes the flip *rides* are owned by the Commander-driven pre-flip pair, NOT by this entity — confirmed against HEAD so they are not double-owned:
- The `release.yml` version-stamp retarget (today hardcoded `git switch next` / `git push origin next` at `release.yml:161,172` — retarget to `main`, SINGLE-TARGET because `release.yml` is tag-push-only; not a per-channel arm) and the `devBranch` ldflag (today hardcoded `-X …cli.devBranch=next` at `.goreleaser.yaml:36` — stable stamps `main`) are **`k6d`**'s deliverable.
- The `e2e-gate` job + SHA-match predicate + waiver are **`nzb`**'s deliverable.

This milestone is gated on `nzb` + `k6d` landing on `next` first. If either slips, the flip's AC-3/AC-7 cannot be satisfied as written — surface to the captain rather than reimplementing the pre-flip mechanics here.

## Marketplace resolver verification (2026-06-08, FO + captain)

The flip's load-bearing unknown — does the host resolve the plugin from the marketplace `@ref` checkout, or from the manifest `source.url@source.ref`? — was exercised end-to-end before any marketplace change. A local self-referential git marketplace (`claude` 2.1.165) was added at branch `chan-a`, whose `marketplace.json` declared `source: {url, ref: chan-b}` (deliberately divergent); `claude plugin install` resolved the plugin content from **`chan-b`** (the manifest `source.ref`), NOT the `chan-a` checkout. **`source.ref` is authoritative for plugin content.** Consequence: each branch's `marketplace.json` must point at itself — `next`→`next` (edge), `main`→`main` (stable). The original AC-6 (fold `source.ref: main` onto `next`) would have cloned edge users' plugin from `main` — an edge break — and is retired in favor of a `main`-only post-flip settle.

**Deferred (captain 2026-06-08):** making `next` an independent release line (its own `0.21.0` version/cask, decoupled from the shared `v*` tag) is a SEPARATE follow-up task — NOT this milestone. This milestone cuts `main` `0.20.0` and flips; the edge cask continues to ride the shared tag + `next-publish.yml` calendar bumps. A standalone marketplace repo (one manifest, explicit per-channel refs) is the longer-term simplification that would moot the branch-local `source.ref` entirely — also deferred.

## Acceptance criteria

Each AC binds to a numbered item of the "Flip checklist — carve completeness audit (2026-06-08)" below (the repo-verified task list for DoD 5–10) and is verified by something OUTSIDE this entity body — a tag/branch state, command output/exit code, a produced commit, or a Go test that fails on violation.

**AC-1 — Archive ref under a NON-`v*` name; the archive cannot fire goreleaser.** (Checklist: *Archive-tag trigger guard* + DoD 5.)
The pre-flip `origin/main` tip is preserved under a ref OUTSIDE the `release.yml` `v*` trigger glob (decided: `archive/v0` — `v0-archived` is rejected because it matches `v*`/`v[0-9]*`). The non-`v*` ref alone is sufficient to prevent the archive firing goreleaser — proven by exercise: `archive/v0` does not match the `v*` glob while `v0-archived` does (recorded in the stage report).
Verified by: `git rev-parse archive/v0` equals the recorded pre-flip `origin/main` SHA, AND `archive/v0` is not in the `v*` namespace (`git tag -l 'v*' | grep -x archive/v0` is empty). This is the mandatory proof.
Optional follow-up (NOT a mandatory AC, NOT owned by this flip — `release.yml` step authoring is out of pj's and k6d's scope): a `release.yml` defense-in-depth guard refusing any pushed tag not matching strict `v[0-9]+.[0-9]+.[0-9]+`, covered by an `internal/release` workflow-guard test (the audit flagged `v[0-9]*` as too weak). File it as a separate hardening task; the flip does not depend on it because the non-`v*` ref already prevents firing.

**AC-2 — `origin/main` carries the prepared 0.20.0 line; `origin/next` retained; replacement was guarded.** (Checklist: DoD 6.)
Verified by: post-flip `git rev-parse origin/main` equals the prepared-line SHA recorded immediately pre-push; `git rev-parse origin/next` still resolves (branch retained); and the update was performed with `git push --force-with-lease=main:<recorded-old-main-sha> origin <prepared>:main` (CAS guard against a racing write to `main`).

**AC-3 — 0.20.0 is cut on `main` through the enforced e2e-gate; the marketplace resolves it from `main`.** (Checklist: *Marketplace ref flip* + DoD 7+8.)
Verified by: (a) the `v0.20.0` tag exists and the GitHub Release published (goreleaser succeeded → `release.yml`'s `e2e-gate` found a green Runtime-Live-E2E run whose `headSha` == the tagged commit, OR a recorded `SPACEDOCK_E2E_GATE_WAIVER`); (b) a real install/resolve exercise — a fresh-HOME `spacedock claude`/`init` auto-install — resolves `0.20.0` from the `main` plugin (observed in plugin state, not by grepping a constant); (c) `main`'s `marketplace.json` `source.ref` is `main` post-flip (see AC-6), so the stable binary resolving `@main` clones the `main` plugin.

**AC-4 — `AGENTS.md` + `docs/releasing.md` describe the post-flip world, verified against the real machinery.** (Checklist: *Doc reconciliation* + DoD 9. README clause DROPPED — README reconciliation already merged.)
Verified by: `AGENTS.md:28` no longer says "Cut releases from `next`… Never release from `origin/main`" — it directs cuts from `main` and names the legacy line in any retained never-clause; and `docs/releasing.md` (which already describes the post-flip world aspirationally) matches the ACTUAL post-flip `release.yml`/`.goreleaser.yaml` after `k6d`'s stamp/devBranch fixes land — checked by reading the reconciled docs against the real workflow files (a divergence between doc and machinery is the failure). No live test; the failable check is doc-vs-machinery agreement on the post-flip refs.

**AC-5 — Three upgrade journeys end in an actionable early upgrade path or a successful launch, never a late dispatch failure.** (Checklist: *Existing-user migration* 5th journey + DoD 10.)
Verified by: isolated host-config exercises or behavior fixtures for (a) old `0.12.x` plugin with no usable binary after plugin refresh — first useful instruction is the binary install + `spacedock claude` payoff; (b) `0.19.x` plugin/binary skew with an outdated binary — the gate/remedy fails EARLY with an actionable upgrade path, not later in dispatch; (c) **the 5th journey** — a `next`-pinned `0.19.x` user post-flip: a stable binary stamped `devBranch=main` is installed, and on the next launch the plugin re-pins to `@main`. Because `frontdoor.go:175` (`Compatible`) is a no-op "proceed to launch" and re-pin (`ops.Install`) fires only on `NoPluginFound` at `frontdoor.go:177/185`, a happy `next`-pinned user does NOT auto-re-pin — so this AC also carries the DECISION (AC-8 below) and, if a nudge is chosen, a fixture proving the `next`-pinned-ref-under-stable-binary nudge fires.

**AC-6 — Marketplace `source.ref` is branch-local: `main` points at `main`, `next` stays `next` — set on `main` AFTER the flip, never folded onto `next`.** (Checklist: *Settle `main`'s marketplace ref post-flip*.)
The host clones the plugin from the manifest's `source.url@source.ref` (VERIFIED — see "Marketplace resolver verification" above); the marketplace `@ref` only selects WHICH branch's `marketplace.json` is read. So each branch's manifest must point at itself: `next`→`next` (edge), `main`→`main` (stable). `next` already does (no fold change). `main`'s `source.ref` is set to `main` by the post-flip `main`-only settle commit (runbook step 8), riding the post-tag `main`-tip window with k6d's stamp.
Verified by: post-flip `git show origin/main:.claude-plugin/marketplace.json` has `source.ref: main` AND `git show origin/next:.claude-plugin/marketplace.json` still has `source.ref: next` (the failable per-branch check); `marketplace_manifest_test.go` accepts `source.ref ∈ {next, main}` so the shared test passes on both channels (the relaxed assertion no longer catches `next` wrongly set to `main` — the per-branch git-state check above is the real proof). The `main` `source.ref` (`main`) and the stable binary's `devBranch` ldflag (`main`, via `k6d`) agree.

**AC-7 — A stable calendar-bump lands on `main` AFTER the flip, making post-flip `plugin update` re-pull `main` against an existing HOME.** (Checklist: *Calendar-bump on main*.)
The calendar-key bump is applied on `main` ONLY after the flip force-push lands (runbook step 8) — it is deliberately HELD OUT of the pre-flip fold. Reason: the calendar `version` key is the ONLY thing that fires a `@next` user's `claude plugin update` re-pull. If the bump rode the pre-flip fold on `next` alongside the ref flip, there would be a transient window where `next` serves `marketplace.json{ref:main, bumped-key}` while `origin/main` is still the OLD pre-flip tip — a `@next` user re-pulling in that window would read `ref:main` and resolve the payload from the STALE old main. Holding the bump out keeps the ref-flip-on-`next` dormant (no re-pull trigger) until `main` actually carries the flipped line.
Verified by: a stable calendar-bump path (mirroring `next-publish.yml`'s `bump-calendar` → push, but targeting `main`) produces a changed `marketplace.json` calendar key on `main` post-flip, and an against-EXISTING-HOME `claude plugin update` then re-pulls `main` (not just a fresh-HOME install). The single mechanism is: stable bump targets `main`, post-flip. (The edge channel keeps its own bump via `next-publish.yml` on `next` — a separate `workflow_dispatch`, unchanged.)

**AC-8 — Dev-stays-on-`next` recorded; the `next`-pinned-user post-flip behavior is a deliberate decision, not an accident.** (Checklist: *Dev stays on `next`* + the 5th-journey decision.)
Verified by: the roadmap/this entity records that post-flip development and Spacedock-state continue on `next` (so the FO-runtime `git … origin next` refs in `claude-first-officer-runtime.md:178-179` stay correct and `main` is stable-release-only) — and the 5th-journey behavior is DECIDED: either (i) accept silent edge-retention for already-happy `next`-pinned users (recorded with rationale), or (ii) add a one-time doctor/remedy nudge when a stable binary sees a `next`-pinned ref, proven by a fixture. The failable artifact is the decision plus, for (ii), the fixture; for (i), the recorded acceptance is a roadmap entry a reviewer can check against the shipped behavior.

## Flip runbook — ordered, captain-gated (implementation spine)

Every outward-facing step is a captain gate (📡). The flip is driven HERE (Shaping FO + captain), not by a cold-boot Commander. **Precondition: `nzb` + `k6d` merged to `next` and green** (this milestone's hard gate).

**Tag-the-green-tip invariant + the `next` freeze window (why the order is what it is):** the green e2e run, the force-pushed `main` tip, and the tagged `v0.20.0` commit MUST all be the SAME commit. The wrinkle: `runtime-live-e2e.yml`'s `workflow_dispatch` has NO per-SHA input (its inputs are `claude_version`/`codex_version`/`effort`) — a dispatch runs against the selected ref's TIP at dispatch time, not the recorded `$PREPARED`. Because `next` keeps moving (dev-stays-on-`next`, AC-8), a concurrent dev push between recording `$PREPARED` and the e2e run would make the green run's `headSha` ≠ `$PREPARED`. Therefore **`next` is FROZEN from the moment `$PREPARED` is recorded (step 3) through the tag (step 7)** — branch-protection lock or a captain-announced dev quiesce for the flip window. And after the e2e goes green (step 4), the green run's `headSha` is RE-VERIFIED to equal `$PREPARED` before proceeding; a mismatch means the freeze leaked — re-record `$PREPARED` and re-dispatch (burning another live run). Nothing may advance the tip between the green run and the tag. The only commit that lands after the tag is `release.yml`'s own post-tag stamp step (`k6d`), which is downstream of the gate and does not affect the tagged SHA. The marketplace `source.ref: main` and the calendar-bump are deliberately NOT in the pre-flip fold (the ref would break edge on `next`; the bump would open a stale-payload window — see AC-6/AC-7); both land on `main` post-flip (step 8).

1. **Pre-cut antipattern audit** on the prepared `next` line (the 0.20.0 content) before anything outward fires. (Internal.)
2. **Fold the 0.20.0 content that must ride the tagged commit onto the prepared `next` line.** A normal commit/PR onto `next`, reviewed before it goes green: doc reconciliation — `AGENTS.md:28` + `docs/releasing.md` against the real post-flip `release.yml`/`.goreleaser.yaml` (AC-4); and relax `marketplace_manifest_test.go` to accept `source.ref ∈ {next, main}` (shared file, two channels) + fix the now-stale `TestTriSurfaceChannelAgreement` comment (keep its skip). (Internal — on `next`, not yet outward.) *(NO marketplace `source.ref` change on `next` — it stays `next` for edge; the host clones the plugin from `source.ref` (verified), so `source.ref: main` on `next` would break edge. `main`'s `source.ref: main` is set post-flip in step 8. The calendar-bump is also NOT folded here — see AC-7.)*
3. **Freeze `next` and record the SHAs.** Engage the freeze: branch-protection lock on `next` or a captain-announced dev quiesce for the flip window — no dev push to `next` until after the tag (step 7). Then `OLD_MAIN=$(git rev-parse origin/main)`; `PREPARED=$(git rev-parse origin/next)`. `$PREPARED` is the single commit the green run, the flip, and the tag all reference. (Internal — pins AC-1's archive target and AC-2's CAS, and the freeze holds it stable.)
4. 📡 **Run the live e2e against `next` and wait green, then RE-VERIFY the headSha.** `workflow_dispatch` of `runtime-live-e2e.yml` against `next` (it runs `next`'s tip = `$PREPARED` because the freeze holds it); approve each spending environment; confirm all 5 legs `success`. THEN re-verify `gh run view <id> --json headSha` == `$PREPARED`. On mismatch (the freeze leaked a push), re-record `$PREPARED` (step 3) and re-dispatch — this burns another live run, stated here so it is not a cut-day surprise. The matching green run's `headSha` == `$PREPARED` is what `release.yml`'s `e2e-gate` matches at the tag. (Feeds AC-3 via tag-the-green-tip.) *Captain gate: authorize the live API spend.*
5. 📡 **Archive current `origin/main` under a non-`v*` ref.** `git tag archive/v0 $OLD_MAIN && git push origin archive/v0` (refuse `v0-archived` — matches `v*`). (AC-1.) *Captain gate: confirm the archive target SHA before pushing.*
6. 📡 **Guarded flip `next`→`main`.** `git push --force-with-lease=main:$OLD_MAIN origin $PREPARED:main`; verify `origin/next` still resolves and `origin/main` now == `$PREPARED`. (AC-2.) *Captain gate: the non-ff replacement of `main`.*
7. 📡 **Cut 0.20.0 — tag the green tip.** Tag `$PREPARED` (now `origin/main`'s tip, the green-run commit) `v0.20.0`, annotated; push the tag. `release.yml`'s `e2e-gate` resolves the tagged SHA == `$PREPARED`, finds the matching green run from step 4 → goreleaser publishes the Release + bumps the brew tap; the version-stamp step (`k6d`, retargeted to `main`) commits the stable stamp on `main` AFTER the tag (downstream of the gate, does not alter the tagged commit). (AC-3.) *Captain gate: the actual release.* **After the tag is pushed, the freeze on `next` may be released** (dev resumes on `next`).
8. 📡 **Settle `main` (post-flip, main-only): `source.ref: main` + calendar-bump.** A pj-owned `main`-only commit (a) sets `.claude-plugin/marketplace.json` `source.ref` `next`→`main` (the stable channel's plugin source — `main` serves from its branch tip, so this post-tag commit resolves stable correctly), and (b) applies the stable calendar-key bump — the trigger that makes existing-HOME `@…` users' `claude plugin update` re-pull `main` (held until here to avoid the stale-payload window, AC-7). `next` keeps `source.ref: next` (edge unchanged). *Captain gate (pushes to `main`).*
9. **Post-flip verification** — fresh-HOME install resolves `0.20.0` from `main` (AC-3b); existing-HOME `plugin update` re-pulls `main` against the post-flip bumped calendar key (AC-7); the three upgrade journeys incl. the 5th (AC-5); confirm dev resumed on `next` and record dev-stays-on-`next` (AC-8). *(Some release-machinery issues only manifest once the tag fires — verify after.)*

## Test plan

What verifies each AC, its cost, and the fixture/CLI/live split:

- **AC-1** — `git rev-parse archive/v0` == recorded `OLD_MAIN` + `git tag -l 'v*' | grep -x archive/v0` empty (free, offline). This is the full mandatory proof; the non-`v*` ref alone prevents firing (proven by exercise). The strict-semver `release.yml` guard is an OPTIONAL follow-up task (out of pj/k6d scope) — if filed, it reuses `workflow_exec_guard_test.go` with an adversarial reject variant.
- **AC-2** — `git rev-parse origin/main`/`origin/next` post-flip (free); the `--force-with-lease=…:$OLD_MAIN` form is the guarded replacement. Done live at the flip (captain-gated).
- **AC-3** — goreleaser success (the gate held) + a fresh-HOME `spacedock claude`/`init` install resolving `0.20.0` from `main`, observed in plugin state. Live, at the flip. The e2e-gate mechanism itself is `nzb`-proven; this exercises the end-to-end cut.
- **AC-4** — read reconciled `AGENTS.md`/`docs/releasing.md` against the real post-flip workflow files; a doc-vs-machinery divergence fails. Optionally anchor the manual read to a machine fact with a small `internal/release` assertion that the post-flip `release.yml` stamp target and `.goreleaser.yaml` `devBranch` are both `main` (parses the real workflow YAML, fails if k6d's fixes regress) — so the doc claim "stamps on `main`" is backed by a failing-on-violation test, not a human read alone.
- **AC-5** — isolated host-config exercises / behavior fixtures for the three journeys (0.12.x-no-binary, 0.19.x-skew, 5th `next`-pinned post-flip). Reuses the existing front-door/doctor fixture patterns. Fixture-level; the 5th may need a live front-door smoke for the re-pin observation.
- **AC-6** — post-flip `git show origin/main:…source.ref` == `main` and `git show origin/next:…source.ref` == `next` (free, the failable per-branch check); `go test ./skills/integration/` green on `next` with the relaxed `∈ {next,main}` assertion. The resolver mechanism is verified (see "Marketplace resolver verification"), not re-spiked. Offline.
- **AC-7** — a post-flip `bump-calendar` on `main` producing a changed calendar key + an existing-HOME `plugin update` re-pull resolving from `main`. Live-ish; reuses `next-publish.yml`'s `bump-calendar` tool. The single pinned mechanism: stable bump targets `main`, applied AFTER the flip (runbook step 8) to avoid the stale-payload window.
- **AC-8** — roadmap/entity record (decision) + (for the nudge option) a doctor/remedy fixture firing on a `next`-pinned ref under a stable binary. Fixture-level if nudge chosen; recorded-decision otherwise.

This is an outward-facing release — captain-gated at each outward step (📡 above). Final Runtime Live E2E on the prepared tip per runbook step 4, unless the captain records a `SPACEDOCK_E2E_GATE_WAIVER`.

## Flip checklist — carve completeness audit (2026-06-08)

An adversarial completeness sweep of the nzb→k6d→pj carve surfaced these flip-touchpoints, each with no prior owner in the carve. Verified against the repo at HEAD. Work them as part of the flip (driven here) unless the pre-flip sprint (nzb + k6d, Commander) lands one first. README reconciliation is already merged; the items below are what remains.

### Release-machinery guards (must hold before the cut)

- [ ] **Archive-tag trigger guard.** `release.yml` triggers on `push: tags: ['v*']`, and `v0-archived` matches `v*` (and even `v[0-9]*`) — pushing the archive tag would fire goreleaser on the legacy main tip and cut a bogus release. Fix: name the archive ref outside the release glob (e.g. `archive/v0`), AND/OR add a `release.yml` guard refusing any tag not matching strict `v[0-9]+.[0-9]+.[0-9]+`, covered by an `internal/release` workflow-guard test (`v[0-9]*` is NOT enough). Ref: `release.yml:7-9`.
- [ ] **Channel-aware version stamp.** The "Stamp plugin manifests to the release version" step is hardcoded `git switch next` → `git push origin next` (`release.yml:161,172`). A stable v0.20.0 cut on main must stamp/push to `main`; otherwise the stable plugin on main stays un-stamped. Cover with the workflow-guard test. (k6d's mechanism — confirm k6d owns it, else land it here.) **Superseded framing (staff review, 2026-06-08):** the original "stable→`main`, edge→`next`" two-arm framing is corrected — `release.yml` is tag-push-only (`on: push: tags: ['v*']`, a single trigger), so its stamp is SINGLE-TARGET `main` (the cut lane after the flip), not a per-channel branch within `release.yml`. The edge channel keeps flowing via `next-publish.yml` (a separate `workflow_dispatch` calendar-bump on `next`), not through a release.yml edge arm. pj and k6d share this one source of truth: release.yml stamp → `main`.

### Marketplace ref (branch-local — settled on `main` post-flip, NOT folded onto `next`)

- [ ] **Settle `main`'s `source.ref: main` post-flip (main-only commit, step 8).** After the force-push, a pj-owned `main`-only commit sets `.claude-plugin/marketplace.json` `source.ref` `next`→`main` (+ calendar bump). `next` keeps `source.ref: next` — it is the edge channel's manifest, and the host clones the plugin from `source.ref` (VERIFIED 2026-06-08), so `source.ref: main` on `next` would break edge. On `next`, relax `skills/integration/marketplace_manifest_test.go` to accept `source.ref ∈ {next, main}`; keep `TestTriSurfaceChannelAgreement`'s skip (it asserts the `main` agreement only on a `main`-ref tree) and fix its now-stale comment.

### Doc reconciliation (README already merged — these still drift)

- [ ] **`AGENTS.md:28`** says "Cut releases from `next`… Never release from `origin/main`." — directly forbids the post-flip cut. Change to release from `main`; invert the never-clause to name the legacy line. Load-bearing for any agent-run cut.
- [ ] **`docs/releasing.md`** already describes the post-flip world ("cut from main", "stamp on main") while the machinery is still pre-flip. Reconcile it against the ACTUAL post-flip `release.yml`/goreleaser after the stamp fix lands — not the aspirational text there now.

### Existing-user migration (confirm specifics at ideation)

- [ ] **5th upgrade journey: the `next`-pinned 0.19.x user, post-flip.** frontdoor re-pins only on `NoPluginFound`, not `Compatible` (`frontdoor.go:177/325` vs `:185/332`), so a happy user with a working `next`-pinned plugin keeps pulling EDGE updates after the flip. Prove the migration (binary upgrade stamped `main` → install re-pins to `@main`) and DECIDE: accept silent edge-retention, or add a one-time doctor/remedy nudge when a stable binary sees a `next`-pinned ref. (Augments AC-5.)
- [ ] **Calendar-bump on main.** The marketplace `version` calendar key — the only thing making `claude plugin update` re-pull a moving branch — is bumped only by `next-publish.yml` (push origin next). Add a stable calendar-bump-on-`main` path so post-flip `plugin update` re-pulls main against an EXISTING (not just fresh) HOME.

### Runbook decisions (pin before the cut)

- [ ] **e2e-gate headSha path.** nzb's gate binds to the tagged commit's SHA, and `runtime-live-e2e.yml` runs on PR-to-next. If the flip post-stamps a commit on main, its SHA differs from any green next-PR run. Pin ONE: stamp-before-the-e2e-run / tag-the-green-tip / recorded `SPACEDOCK_E2E_GATE_WAIVER` (with reason) — don't leave it a cut-day surprise.
- [ ] **Dev stays on `next`.** Record explicitly that post-flip development/state continues on `next` (so the FO-runtime refs `git … origin next` in `skills/first-officer/references/claude-first-officer-runtime.md` stay correct) and `main` is the stable release lane only.

## Stage Report: ideation

- DONE: Firm the flip ACs against the 9-item "Flip checklist — carve completeness audit (2026-06-08)" already in this entity's body; each surviving AC binds to a checklist item and is verified OUTSIDE this entity body. Drop the README half of AC-4; keep AGENTS.md + docs/releasing.md reconciliation; keep upgrade-path AC-5 and ADD the 5th journey.
  8 ACs written, each bound to a named checklist item (AC-1→archive-tag guard/DoD5; AC-2→DoD6; AC-3→marketplace-flip/DoD7-8; AC-4→doc-reconciliation/DoD9, README clause DROPPED; AC-5→migration incl. 5th journey/DoD10; AC-6→paired ref+test commit; AC-7→calendar-bump-on-main; AC-8→dev-stays-on-next + 5th-journey decision). Each "Verified by" names a git ref/command, a Go test, a fixture, or doc-vs-machinery agreement — none rests on reading this body.
- DONE: Produce the ordered, captain-gated flip RUNBOOK as the implementation spine.
  8-step runbook with 📡 captain gates on every outward step, ordered around the tag-the-green-tip invariant (green run, force-pushed `main` tip, and tagged commit must all be ONE SHA): antipattern audit → fold ALL 0.20.0 content onto the prepared `next` line (marketplace ref+test in one commit, calendar-bump, doc reconciliation) so the tip is final before the e2e → record `$PREPARED` → live e2e on `$PREPARED` (workflow_dispatch) green → archive `archive/v0` → guarded `--force-with-lease=main:$OLD_MAIN` flip → tag-the-green-tip cut 0.20.0 → post-flip verification.
- DONE: Record the riskiest-mechanism determination explicitly (no spike needed, naming proven mechanisms) AND decide the e2e-gate headSha path.
  "No spike needed" recorded citing s0cq (marketplace re-pin PASSED 2026-06-03), nzb (e2e-gate query proven against live repo, banked), git `--force-with-lease`. HeadSha path DECIDED: tag-the-green-tip — `runtime-live-e2e.yml` carries `workflow_dispatch` (verified at `on:`), so dispatch the live matrix on the prepared `main` tip, get green, tag THAT commit; k6d's stamp commits post-tag so no SHA drift between green-run and tag. `SPACEDOCK_E2E_GATE_WAIVER` is the recorded fallback.

### Summary

Firmed the flip from a placeholder spec into 8 checklist-bound ACs + a 10-step captain-gated runbook + an explicit no-spike determination. Three load-bearing claims were proven by exercise, not assertion: (1) flipping the marketplace ref `next`→`main` alone FAILS `marketplace_manifest_test.go:68` ("source.ref = main, want next"), proving AC-6's paired-commit "green-by-construction" requirement is real; (2) `archive/v0` does NOT match the `release.yml` `v*` glob while `v0-archived` DOES, confirming the audit's archive-name fix; (3) `internal/release` + `skills/integration` baselines green so implementers start clean. Decided the e2e-gate headSha path (tag-the-green-tip via workflow_dispatch) so the cut is not a cut-day surprise — and on self-review caught and fixed a sequencing bug in my own first runbook draft: the marketplace-flip/calendar/doc commits must be folded into the prepared tip BEFORE the live e2e, not pushed to `main` after the green run, or the green-run SHA / force-pushed tip / tagged commit would diverge and the released 0.20.0 would not carry the `main` marketplace ref. Confirmed the `release.yml:159-172` channel-aware stamp and `.goreleaser.yaml:36` per-channel devBranch are `k6d`'s deliverable (not double-owned here), and the e2e-gate is `nzb`'s — this milestone is gated on both landing on `next` first. The 5th-journey finding is verified against `frontdoor.go:175/177/185` (re-pin fires only on NoPluginFound, not Compatible), so a happy `next`-pinned user does not auto-migrate — carried as AC-8's explicit decision (accept silent edge-retention vs. add a nudge).

### Staff-review fold

Folded the preflight staff-review findings (2026-06-08):

- **M3 (Material) — `next` freeze window.** `runtime-live-e2e.yml`'s `workflow_dispatch` has NO per-SHA input (inputs are `claude_version`/`codex_version`/`effort`, confirmed by reading the `on:` block) — it runs the selected ref's TIP at dispatch time, and `next` moves concurrently. Extended the tag-the-green-tip invariant: `next` is FROZEN from step 3 (record `$PREPARED`) through step 7 (tag), via branch-protection lock or captain-announced dev quiesce; added a post-step-4 `headSha == $PREPARED` re-verify with explicit re-record+re-dispatch (burns another live run) on mismatch. Updated the invariant note + runbook steps 3, 4, 7.
- **M4+M5 (Material, one fix) — calendar-bump moved post-flip.** Holding the marketplace ref flip AND the calendar-key bump on `next` through the flip window opened a stale-payload window: `next` would serve `marketplace.json{ref:main, bumped-key}` while `origin/main` is still the OLD tip, so a `@next` `plugin update` (fired by the bumped key) resolves the payload from stale old main. Fix: dropped the calendar-bump from the pre-flip fold (runbook step 2); added a post-flip calendar-bump-on-`main` step (runbook step 8). The ref flip stays in the fold but is dormant (no calendar-key change → no `@next` re-pull). Pinned AC-7 to the single mechanism (stable bump targets `main`, post-flip), resolving the prior "implementation call" punt.
- **P1 (polish) — stamp is single-target `main`.** `release.yml` is tag-push-only (`on: push: tags: ['v*']`), so k6d's version-stamp has no edge arm — it is SINGLE-TARGET `main`, with edge flowing via `next-publish.yml` (separate `workflow_dispatch`). Corrected the checklist "Channel-aware version stamp" item (superseded-framing note) and the boundary section so pj and k6d share one source of truth.
- **P3 (polish) — strict-semver guard demoted to optional.** AC-1's mandatory proof is now ONLY the non-`v*` archive ref (`archive/v0`), which alone prevents the archive firing (proven by exercise). The strict-semver `release.yml` guard step (unowned — out of pj/k6d scope) is an explicit OPTIONAL follow-up task, not a blocking AC.
- **P4 (polish, optional) — AC-4 anchor.** Added an optional small `internal/release` assertion (post-flip stamp target + `.goreleaser.yaml` devBranch are `main`) so the doc-vs-machinery check has a failing-on-violation artifact, not a human read alone.
- **P5 (polish) — citation normalized.** The hardcoded `git switch/push next` lines are now cited as `release.yml:161,172` throughout (was `:159-172`/`:156,161,172`).
