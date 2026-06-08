---
id: pj6k0h83g6taszkmt92qvsc0
title: Main-flip milestone — cut 0.20.0 on main + flip the marketplace once 0.19.6 is tested
status: backlog
source: "captain (2026-06-05) — '0.19.6: flip main readiness — once tested i want to make 0.20.0 on main and flip the marketplace.' The capstone of the 0.19.6 line."
score: "0.38"
started:
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

## Acceptance criteria

(To firm up at ideation once the specifics are clarified.)
**AC-1 — Current pre-v1 main is archived before replacement.**
Verified by: tag `v0-archived` exists and points at the pre-flip `origin/main` tip recorded immediately before the force-push.

**AC-2 — `main` is replaced by the prepared 0.20.0-ready line while `next` remains available.**
Verified by: `origin/main` points at the prepared release line, `origin/next` still exists, and the branch update command used `--force-with-lease` or an equivalent guarded replacement.

**AC-3 — 0.20.0 is cut on main and the marketplace serves it.**
Verified by: the 0.20.0 tag/release exists on main (the cut succeeded through the fixed release-gate), and the marketplace install path resolves 0.20.0 from main (an install/resolve exercise confirming the flipped ref, not prose) — exact mechanics per the clarified specifics.

**AC-4 — README/install docs are reconciled for stable main plus dev next.**
Verified by: docs review plus command examples that match observable install behavior: stable install resolves from `main`, dev-only publish/install path still names `next`, and there is no stale wording that presents `next` as the stable release lane.

**AC-5 — Upgrade paths are confirmed for stale plugin / missing binary and outdated binary cases.**
Verified by: isolated host-config exercises or behavior fixtures for (a) old `0.12.x` plugin with no usable binary after plugin refresh, and (b) `0.19.x` plugin/binary skew with an outdated binary. Each must end in an actionable early install/upgrade instruction or a successful launch, not a late dispatch failure.

## Test plan

Per the clarified scope. At minimum:

- Final Runtime Live E2E on the prepared tip, unless the captain explicitly waives it.
- `v0-archived` tag check against the pre-flip `origin/main` SHA.
- Guarded replacement of `origin/main` with the prepared line, while preserving `origin/next`.
- Release-mechanics check that stable release stamping/marketplace refs target `main` and dev-only publish remains possible from `next`.
- README/install-doc reconciliation review against real command output.
- Isolated upgrade-path exercises for old-plugin/no-binary and outdated-binary skew.
- `v0.20.0` annotated tag/release, release action success, Homebrew tap update, and post-flip marketplace install resolving 0.20.0 from `main`.

This is an outward-facing release — captain-gated at each outward step.

## Flip checklist — carve completeness audit (2026-06-08)

An adversarial completeness sweep of the nzb→k6d→pj carve surfaced these flip-touchpoints, each with no prior owner in the carve. Verified against the repo at HEAD. Work them as part of the flip (driven here) unless the pre-flip sprint (nzb + k6d, Commander) lands one first. README reconciliation is already merged; the items below are what remains.

### Release-machinery guards (must hold before the cut)

- [ ] **Archive-tag trigger guard.** `release.yml` triggers on `push: tags: ['v*']`, and `v0-archived` matches `v*` (and even `v[0-9]*`) — pushing the archive tag would fire goreleaser on the legacy main tip and cut a bogus release. Fix: name the archive ref outside the release glob (e.g. `archive/v0`), AND/OR add a `release.yml` guard refusing any tag not matching strict `v[0-9]+.[0-9]+.[0-9]+`, covered by an `internal/release` workflow-guard test (`v[0-9]*` is NOT enough). Ref: `release.yml:7-9`.
- [ ] **Channel-aware version stamp.** The "Stamp plugin manifests to the release version" step is hardcoded `git switch next` → `git push origin next` (`release.yml:156,161,172`). A stable v0.20.0 cut on main must derive its stamp/push branch from the release channel (stable→`main`, edge→`next`); otherwise the displayed plugin version lands on the edge channel while the stable plugin on main stays un-stamped. Cover with the workflow-guard test. (Naturally k6d's two-channel mechanism — confirm k6d owns it, else land it here.)

### Marketplace ref flip (not-green-by-construction without the paired test edit)

- [ ] **Retarget the ref + flip the guard test in ONE commit.** Edit `.claude-plugin/marketplace.json` `source.ref` `next`→`main`, AND in the same commit update `skills/integration/marketplace_manifest_test.go:67` (currently asserts `ref=="next"` is the required value) to expect `main`. Without the paired edit, pj's own "line ALL GREEN" gate blocks the flip by construction. The marketplace ref and the stable binary's devBranch ldflag must agree (both `main`).

### Doc reconciliation (README already merged — these still drift)

- [ ] **`AGENTS.md:28`** says "Cut releases from `next`… Never release from `origin/main`." — directly forbids the post-flip cut. Change to release from `main`; invert the never-clause to name the legacy line. Load-bearing for any agent-run cut.
- [ ] **`docs/releasing.md`** already describes the post-flip world ("cut from main", "stamp on main") while the machinery is still pre-flip. Reconcile it against the ACTUAL post-flip `release.yml`/goreleaser after the stamp fix lands — not the aspirational text there now.

### Existing-user migration (confirm specifics at ideation)

- [ ] **5th upgrade journey: the `next`-pinned 0.19.x user, post-flip.** frontdoor re-pins only on `NoPluginFound`, not `Compatible` (`frontdoor.go:177/325` vs `:185/332`), so a happy user with a working `next`-pinned plugin keeps pulling EDGE updates after the flip. Prove the migration (binary upgrade stamped `main` → install re-pins to `@main`) and DECIDE: accept silent edge-retention, or add a one-time doctor/remedy nudge when a stable binary sees a `next`-pinned ref. (Augments AC-5.)
- [ ] **Calendar-bump on main.** The marketplace `version` calendar key — the only thing making `claude plugin update` re-pull a moving branch — is bumped only by `next-publish.yml` (push origin next). Add a stable calendar-bump-on-`main` path so post-flip `plugin update` re-pulls main against an EXISTING (not just fresh) HOME.

### Runbook decisions (pin before the cut)

- [ ] **e2e-gate headSha path.** nzb's gate binds to the tagged commit's SHA, and `runtime-live-e2e.yml` runs on PR-to-next. If the flip post-stamps a commit on main, its SHA differs from any green next-PR run. Pin ONE: stamp-before-the-e2e-run / tag-the-green-tip / recorded `SPACEDOCK_E2E_GATE_WAIVER` (with reason) — don't leave it a cut-day surprise.
- [ ] **Dev stays on `next`.** Record explicitly that post-flip development/state continues on `next` (so the FO-runtime refs `git … origin next` in `skills/first-officer/references/claude-first-officer-runtime.md` stay correct) and `main` is the stable release lane only.
