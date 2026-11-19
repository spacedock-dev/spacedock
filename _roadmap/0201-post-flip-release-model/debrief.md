---
session-date: 2026-06-13
sequence: 1
first-commit: 41bd47ae
last-commit: 8ab302ec
release: v0.20.1
duration: ~1d (long single FO drive, spanned one compaction)
---

# Session Debrief — Sprint 0201-post-flip-release-model (FO drive) — 2026-06-13 #1

FO drive of the post-flip release model to a **cut v0.20.1**. All six gate-approved members landed (gp the centerpiece), plus the doc-owner yw, plus four pre-cut/follow-up hardening PRs the drive surfaced. The captain took the conn (approve + merge + CI authority). The mandated pre-cut antipattern audit ran, found one BLOCKER (a process gate, not a code defect), and the cut shipped clean once it cleared.

## Shipped — sprint members
- **8p** `brew-cask-agentsview-safehouse-deps` — [#347](https://github.com/spacedock-dev/spacedock/pull/347). Brew cask `depends_on` agentsview; safehouse corrected to its real third-party tap (`eugene1g/safehouse/agent-safehouse`) after a cross-tap-`depends_on` spike showed Homebrew 6.0 tap-trust won't auto-resolve it.
- **zrc** `non-sandboxed-auto-mode` — [#349](https://github.com/spacedock-dev/spacedock/pull/349). Unsandboxed `spacedock claude` starts in `--permission-mode auto`; unsandboxed `spacedock codex` in `--ask-for-approval on-request`; sandboxed arms skip/bypass instead.
- **gj** `startup-sandbox-status` — [#350](https://github.com/spacedock-dev/spacedock/pull/350). Launch + `--version` show sandbox posture and per-runtime status.
- **te** `install-refresh-and-upgrade-hint` — [#351](https://github.com/spacedock-dev/spacedock/pull/351). `install --host codex` refreshes a present plugin; an opt-in, non-blocking upgrade hint fires when a contract-compatible-but-behind plugin has a newer version (integer-semver compare, pinned by boundary/pre-release cases after a detached audit caught a lexical-compare regression).
- **gp** `marketplace-repo-and-pinned-channels` (CENTERPIECE) — [#352](https://github.com/spacedock-dev/spacedock/pull/352). Model B: marketplace manifest decoupled to a standalone `spacedock-dev/marketplace` repo; two channels selected by entry name via the `devBranch` stamp (`spacedock` stable / `spacedock-edge` edge); `@branch` shorthand dropped. Detached adversarial audit found 1 Material (a stale plugin-repo ref + dead `@branch` build in `contract.go`'s predates-contract remedy) → fixed in a feedback cycle (dropped the now-dead `branch` param across the contract chain). All four live-e2e hosts green.
- **yw** `readme-and-docs-architecture-cleanup` (DOC OWNER) — [#355](https://github.com/spacedock-dev/spacedock/pull/355). Sole writer of the sprint's user-visible docs: slimmed README to a front door (links `spacedock.md/docs`), removed the banned prose-grep `install_doc_test.go`, and landed the five behaviors' docs on canonical site pages. An adversarial doc-vs-shipped-code fidelity audit (6 skeptics) confirmed every behavior faithful.

## Shipped — pre-cut / follow-up hardening (FO-surfaced, not original sprint entities)
- **ci-pr-gating** — [#348](https://github.com/spacedock-dev/spacedock/pull/348). Gate main-branch PRs on the install-e2e + runtime-live-e2e offline suite.
- **marketplace bridge + LICENSE** — [#353](https://github.com/spacedock-dev/spacedock/pull/353). gp's manifest removal broke the released v0.20.0 install path (default branch `main`; v0.20.0 binary resolves `spacedock-dev/spacedock@main`). Restored a transitional bridge manifest on `main`, re-expressed the absence guard as a PRESENCE guard, and added the Apache-2.0 LICENSE (the README front-door link was 404ing — the blob never merged to main). **The decouple should have been cutover-last, not decouple-first.**
- **da** `version-show-runtime-plugin-version` — [#354](https://github.com/spacedock-dev/spacedock/pull/354). `--version` shows the per-runtime plugin version (read robustly from the resolved manifest, not the fragile `plugin list` probe) and drops the invented "enablement" jargon for plain words. Captain-directed.
- **install-doc repoint** — [#356](https://github.com/spacedock-dev/spacedock/pull/356). Manual `marketplace add` doc pointed at the bridge repo; repointed to canonical `spacedock-dev/marketplace`.
- **stable-ref automation** — [#357](https://github.com/spacedock-dev/spacedock/pull/357). The stable marketplace entry now tracks a moving `stable` branch instead of a per-release tag; `release.yml` auto-advances `stable` to each cut (same-repo push, `GITHUB_TOKEN`). The manifest is now static forever. Replaces the manual post-cut marketplace repoint (design (b) of the deferred `ez` lane). **Validated live on v0.20.1's own release.**
- **live parallelism** — [#358](https://github.com/spacedock-dev/spacedock/pull/358). The four claude shared scenarios fan out with `t.Parallel` (per-scenario `CLAUDE_CONFIG_DIR` isolation) behind the existing `TestLiveEnsignCycle` canary — cuts the claude-opus long pole ~27m→~9m. Compiles `-tags live`; parallel isolation gets its first live exercise next run.

## The cut
`v0.20.1` tagged on `5a8880c9`, annotated body = the release notes, e2e-gate satisfied by a green Runtime Live E2E run bound to that exact SHA. goreleaser produced the GitHub Release + 8 tarballs (stable + edge × darwin/linux × amd64/arm64) + both Homebrew casks bumped to 0.20.1. The stamp step stamped plugin.json→0.20.1 on `main` and auto-advanced `stable` (commit `2b53387a`). End-to-end verified: stable channel resolves 0.20.1, binary stamps 0.20.1.

## Decisions
- **Took the conn** (captain delegated approve + merge + CI). Merged behavior members on full live-e2e green; merged docs/CI/release-machinery-only PRs (#353/#354/#356/#357/#358) on the gating lanes (build/offline/install) **without** burning the env-gated live lanes — those PRs don't change the runtime the live lanes exercise.
- **Marketplace break → surgical bridge, not gp revert** (captain's call): restore the manifest on `main` as a transitional bridge, defer its removal post-cutover, re-express the guard tests for presence.
- **Stable channel → moving `stable` ref (design b)** done as plumbing without an entity (captain: "avoid ceremonies"). Side benefit: it makes the `ez` stamp-then-tag concern moot for the stable channel (stable follows the post-stamp branch, not the tag).
- **Tagged the gated commit directly (no pre-stamp)** — pre-stamping creates a new ungated commit the exact-SHA e2e-gate would block on; the moving `stable` ref handles the plugin version post-tag instead. Same ez quirk as v0.20.0 (tag tree plugin.json briefly behind), invisible to the stable-channel install path.
- **Live parallelism: validate-on-next-run** (captain) — merged on offline-green; first real parallel exercise is the next live dispatch.

## Issues — Workflow
- **The release machinery was decouple-first when it needed to be cutover-last.** gp removing the manifest from `main` broke every in-the-wild v0.20.0 install — caught only by the captain, not by CI or the per-member validation/audit (all of which tested the *new* binary, never the *released* one). Root of the #353 hotfix.
- **The e2e-gate is exact-SHA, so every main-advancing fix re-gates the cut.** Landing #356 then #357 each moved the cut commit and forced a fresh ~19m live-e2e run. Cost real wall-clock; correct behavior, but a friction worth noting (a fix mid-cut is expensive).
- **Crossed-wire with the gp ensign** on marketplace-repo provisioning: an unconditional "Proceed" reached it a beat before a "hold" — net state was exactly what the captain authorized, no harm, but sequencing supersede-shutdown before re-instruction matters.
- **Mis-addressed an ensign message** (used the long emitted agent name, not the 64-char-capped override) — the fold sat in a phantom inbox until the captain noticed the ensign wasn't working. Lesson: address by the *spawned* name.

## Issues — Spacedock
- **CI never tests the RELEASED binary's install path** — the live-e2e lanes build and exercise the PR (next/Model-B) binary, never the stable binary users actually run. This is precisely why gp's manifest removal sailed through green. Highest-leverage CI gap (noted, not filed per captain).
- **Behavior-coverage gaps (noted, not filed):** G1 — a real host *blocking* on a contract mismatch is never live-driven (all four mismatch verdicts offline-only; the live cycle reaches the gate via `--skip-contract-check`). G2 — the release critical path (e2e-gate→stamp→goreleaser→homebrew) has no repeatable live test; it "only proves itself in production on tag push." Neither has a backlog owner. (Estimate: ~40% of behaviors carry a live drive; FO gate-core + ensign cycle solid, release machinery ~15%.)
- **`cmd/spacedock-release` statement coverage is 25%** — the release tooling, thinly unit-tested (the interactive `notes`/claude path is the untested part).
- **Git push of workflow files needs the keyring token directly** — `gh auth git-credential` serves an OAuth-App token lacking `workflow` scope; `git push https://x-access-token:$(gh auth token)@…` works. Bit the #357 push.
- **gofmt go-version skew** — local go1.26 reformats three files (doc-comment list indent, smart-quote) that are clean under CI's pinned go1.22; no workflow gofmt-gates, so it doesn't block. Only `contract_test.go`'s trailing blank line is real 0201 drift.

## Deferred (tracked, not filed this session)
- Bridge-manifest removal from `main` (post-cutover, after v0.20.0 installs migrate).
- Edge auto-advance (`next-publish.yml` retarget at the standalone repo) — `w6`/`ez`/`tw`/`qp` release-machinery lane.
- Parallel-isolation live validation (next live run).
- The known no-drive behavioral drives — `ev` (FO halt/sync/journey), `e3z` (bare-mode), `95` (pre-completion-TeamDelete ban) — already own their gaps; unimplemented.
