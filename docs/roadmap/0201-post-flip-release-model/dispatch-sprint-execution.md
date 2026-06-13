# 0201 (0.20.1) — Commander dispatch (cold-boot)

> **Status: APPROVED 2026-06-13 (captain) — ready for Commander cold-boot drive.** All six members are staff-clean (`staff-review.md` → `zrc` M1 folded; `staff-review-yw.md` → `yw` M1+M2 folded) and the captain approved the six ideation gates. The six members carry `sprint=0201-post-flip-release-model`, `sprint-readiness=ready`, status `ideation` (gate-approved); the Commander advances each ideation→implementation→validation→done as it drives. Drivable set: `spacedock status --workflow-dir docs/dev --where sprint=0201-post-flip-release-model --where 'sprint-readiness != defer'`.

## Prerequisite
The five behavior members + the docs member are ideated and staff-reviewed (`staff-review.md` — NOT-READY → `zrc` M1 folded → clean; `staff-review-yw.md` pending). Entity state is on `spacedock-state/dev`.

## Boot
```bash
# base on the TRUNK (main), not next — the doc site + install-journey-removal live on main
git fetch origin main && git switch -c drive/0201 origin/main && go build -o ./spacedock ./cmd/spacedock
export SPACEDOCK_BIN="$PWD/spacedock" SPACEDOCK_REPO_ROOT="$PWD"
git -C docs/dev/.spacedock-state pull --rebase origin spacedock-state/dev   # gh-HTTPS if SSH down
# rotate the live-auth token before any live run:
security find-generic-password -s "Claude Code-credentials" -w | python3 -c "import sys,json; print(json.load(sys.stdin)['claudeAiOauth']['accessToken'])" > ~/.claude/benchmark-token
./spacedock status --workflow-dir docs/dev --boot
```
**SSH is down** — push via `git -c credential.helper='!gh auth git-credential' push https://github.com/spacedock-dev/spacedock.git <ref>`. **Worktrees + PRs target `main`** (post-flip trunk). The `next`/`main` reconcile is deferred until after `gp` lands (it establishes the clean topology).

## Deliverable & DoD
**0.20.1** = post-flip UX cleanup + the marketplace decouple (Model B, decouple-first; `ez` stamp-then-tag deferred). Done when, merged to `main`:
- **`gp`** decoupling proven: a tag-pinned stable channel stays frozen while edge tracks `next` HEAD, both from one separate marketplace repo; plugin branches carry no manifest. (This is the centerpiece — it also unblocks the deferred `next→main` topology.)
- **`gj`** `--version` (+ banner + `status --boot`) reports 3-way sandbox state + per-runtime install/enablement, first line `spacedock <ver> (contract <N>)` UNCHANGED.
- **`te`** `spacedock install --host codex` refreshes a present plugin; a behind-but-compatible plugin surfaces an opt-in upgrade hint in doctor + the front door.
- **`zrc`** non-sandboxed launch injects auto-mode (`--permission-mode auto` / codex `--ask-for-approval on-request`).
- **`8p`** the Homebrew cask declares `depends_on cask: agentsview` (safehouse out of scope — not a brew package).
- **`yw`** `docs/install-journey.md` gone, README slimmed to a front-door, `install_doc_test.go` removed, the five behaviors documented on the canonical `docs/site/` pages.

## Drive order — ⚠️ coordination
- **Docs are owned solely by `yw`** — NO behavior member edits a `.md`. That is the anti-collision design; honor it (a member's `git show --stat` must show zero `*.md`).
- **`internal/cli` is the hot code package** — `gj` (--version), `te` (init.go/gate), `zrc` (front-door argv), `gp` (install/channel) all touch it. Sequence or rebase carefully; run `go test ./internal/cli/` green over the WHOLE package per change (the `zrc` M1 lesson).
- **`te` before `gp`'s migration step** (an unreliable refresh strands a migrating user).
- **`gj` ∥ `zrc`** share a concept, not code — parallelize fine (per staff review).
- **`yw` two-wave:** wave-1 structural cleanup (remove install-journey, slim README, kill the test) can lead; wave-2 behavior-docs follow the five members landing.
- **`gp` is the centerpiece** — land it solidly; it gates the topology cleanup.
- **Stray spike worktrees:** `gj`/`zrc`/`8p` ideation spikes left worktrees (`.worktrees/spacedock-ensign-{slug}`, off an earlier `main`) and a stray `worktree:` field on those entities (ideation is non-worktree). Before creating fresh implementation worktrees, `--force`-clear each `worktree:` field (`spacedock status --set {slug} worktree= --force`) and `git worktree remove --force` the stray dirs — then dispatch implementation clean off `main`.

## Per-member build notes
### `gp` — marketplace-repo-and-pinned-channels (release-model) · HIGH-STAKES
Separate marketplace repo; `spacedock`(stable, `ref: vX.Y.Z`) + `spacedock-edge`(next HEAD) entries selected by entry NAME via `devBranch`. `"source":"url"` (host rejects `"git"`). Migration = `installArgvSequence` cleanup-then-pin repoint, gated on `te`. Folds `w6`. ACs: decoupling behavior test (live host) + manifest-absence git check + channel-resolution seam. **Detached adversarial audit before merge** (release machinery).
### `te` — install-refresh-and-upgrade-hint (release-model) · HIGH-STAKES
codex `runInit` arm must call `ops.Install` unconditionally (today it short-circuits to doctor — init.go:47-52). Additive opt-in upgrade hint in doctor + `gateHost`, guarded by the `Version=="dev"` sentinel. **REPLACE `TestInitCodexInstallReadiness/compatible-installed`** — it codifies the bug. **Detached audit before merge** (front-door).
### `gj` — startup-sandbox-status (ux-cleanup)
3-way sandbox state (from `safehouse.Present`/`Available`) on banner / `status --boot` SANDBOX / `--version`. Per-runtime: claude reads the plugin `enabled` field (NEW — not in `pluginListEntry` today), pi via `checkPiRuntime`, probe-error → `enablement unknown`. **AC-3 guards the `--version` first line** the FO/ensign skills parse. Probe seams; no live host in tests.
### `zrc` — non-sandboxed-launch-auto-mode (ux-cleanup)
Inject auto-mode argv on the `!wrap` path. **REQUIRED blast-radius co-edit:** 5 `equalArgv` oracles (frontdoor_test:108; safehouse_frontdoor:159,317; frontdoor_stray_prompt:39,70; plugin_dir_frontdoor:64) + the 2 `--plugin-dir`-is-still-unsandboxed oracles gain the flag; single-insertion-point rule. Success = `go test ./internal/cli/` green whole-package. codex mapping = option A; add `--ask-for-approval` to `valueTakingHostFlags`.
### `8p` — brew-cask-agentsview-safehouse-deps (ux-cleanup)
`depends_on cask: agentsview` in the cask (lands in the separate `homebrew-tap` repo). safehouse omitted (not a Homebrew package — github.com/anthropics/safehouse).
### `yw` — readme-and-docs-architecture-cleanup (docs) · staff-clean (M1+M2 folded)
Sole doc owner. Remove `docs/install-journey.md`; slim README (front-door + docs-site link); remove `install_doc_test.go` (banned prose-grep; `executableShellCommands` helper stays). **Retarget map re-derived against real `origin/main` (HEAD 55158746) — the #343 site is a slim rewrite, not a mirror, so behaviors with no home get a NEW section ADDED:**
- **gj** → ADD `## --version` to `reference/command-reference.md` (today only inline, no example).
- **te** → amend `reference/command-reference.md` `## Setup` (the doctor-refresh line).
- **gp** → `docs/releasing.md` (real, via `contributing/releasing.md` symlink — manifest-in-branch lines) + edge-pointer in `contributing/build-from-source.md`.
- **zrc** → ADD a permission-posture sentence to `reference/command-reference.md` `## Launch` (sole home; phantom `first-launch.md` dropped).
- **8p** → `get-started/install.md` Homebrew.
AC-4 cites only these post-change-checkable anchors. README "full docs" URL = published doc-site URL (TBD — captain, after marketplace).

## Pre-cut antipattern audit (⚠️ before the tag)
With all six merged to `main` and `v0.20.1` NOT yet fired, dispatch an INDEPENDENT staff-eng reviewer over the assembled sprint. **Known cross-cutting item to verify:** CI gates currently fire only on PRs to `next`, not `main` — `main` PRs are ungated. Confirm `main`-PR gating exists (or add it) before relying on the trunk model. Ship-blockers fixed pre-cut; non-blockers seed the next sprint.

## Cut
`go test ./...` green from root, then `docs/releasing.md`. NOTE: the stamp-then-tag ritual (`ez`) is deferred, so the existing cut applies — watch the tag/manifest-version caveat. Cut `v0.20.1`. Captain authorizes.

## Out of scope (deferred)
`ez` (stamp-then-tag), `tw` (next-independent-release-line), `qp` (release-runbook); the `next`/`main` reconcile + roadmap-doc relocation + README URL (all after `gp`).
