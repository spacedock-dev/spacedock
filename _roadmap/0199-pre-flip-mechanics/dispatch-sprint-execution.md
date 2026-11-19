# 0199 pre-flip-hardening — Commander dispatch (cold-boot)

You are the **Commander** for sprint **0199-pre-flip-mechanics**, driving its three tasks to a **0.19.9** cut on `next`. The shaping FO has carved the sprint, driven ideation, and run the gates; your job is the **implementation → validation → done** drive plus the release cut. Drive autonomously; the captain holds the conn for the gates called out below.

This package is self-contained: boot on it cold, then drive.

## Boot

```bash
git fetch origin next && git reset --hard origin/next && go build -o ./spacedock ./cmd/spacedock
git -C docs/dev/.spacedock-state pull --rebase origin spacedock-state/dev
# rotate the live-auth token before any live run:
security find-generic-password -s "Claude Code-credentials" -w | python3 -c "import sys,json; print(json.load(sys.stdin)['claudeAiOauth']['accessToken'])" > ~/.claude/benchmark-token
./spacedock status --workflow-dir docs/dev --boot
```

Then `TeamCreate` (first team-mode call), standing-teammate discovery, reconcile sweep. Use the freshly-built `./spacedock` (its `--version` cosmetically says 0.19.0 — ignore; contract 1 is what matters).

**Membership is the query, not a table:**
```bash
spacedock status --workflow-dir docs/dev --where sprint=0199-pre-flip-mechanics --where 'sprint-readiness != defer'
```
→ **`v3` ship-linux-binaries · `th` safehouse-preserves-spacedock-bin · `jm` entity-label-localization.** All three are ideation-approved; drive each from `implementation`.

## Deliverable & DoD

**Deliverable:** spacedock **0.19.9** cut on `next`.

1. Every active task (`v3`, `th`, `jm`) `done` / PASSED + merged to `next`.
2. `go test ./...` from the repo root green.
3. `v3`'s Linux path proven: `goreleaser --snapshot` produces `linux/amd64`+`linux/arm64` tarballs AND the installer fetches+checksum-verifies+installs a runnable `spacedock` on Linux and macOS.
4. `th`'s re-assert proven against **real safehouse** by the **captain** (this dev environment cannot run real safehouse — see th's note).
5. spacedock **0.19.9** stamped + cut on `next` (captain-gated release).

## Drive order

The three are **independent** — no cross-task dependency. Run them in parallel (stage concurrency 3). `v3` is the heaviest (release-config + installer + CI). Each ideation body carries a done spike and proof-policy-clean ACs; the implementation's first tests are the fixtures those ACs name.

**All three touch a high-stakes surface** (v3 = release machinery, th = front-door, jm = shipped scaffolding) → **each earns a detached adversarial audit at validation** (README `## validation` → "Detached adversarial audit"), on a throwaway checkout of the merge result. Material findings route back through validation→implementation feedback.

## Per-task build notes

### `v3` — ship-linux-binaries (distribution)
- **Scope:** add `linux` to goreleaser `builds.goos` (amd64+arm64); a universal `curl|sh` `install.sh` (OS/arch detect → fetch latest Release tarball → checksum-verify → install to `~/.local/bin`); document the Linux path + the honest safehouse-on-Linux story in `docs/install-journey.md`.
- **Spike done (build half de-risked):** `goreleaser release --snapshot --skip=publish,homebrew` with the one-line `goos` addition built all four targets, produced the linux tarballs + `checksums.txt`, statically-linked ELF confirmed. The install half's only external (GitHub "latest release" API) is in the AC-3 test.
- **Watch:** AC-3 hits the **live GitHub API** (CI network flakiness — keep the OS-logic test offline via the local-dist override, isolate the network path to AC-3). The `curl|sh` checksum is fetched from the same release (standard limitation — fine, don't gold-plate it). Don't touch the darwin/homebrew_casks blocks (AC-5 unregressed).
- **Detached audit:** release/CI machinery — audit the goreleaser config guard + the install.sh checksum gate (tamper case must stay load-bearing).

### `th` — safehouse-preserves-spacedock-bin (dev-quality)
- **Scope:** on the **wrap path only**, prefix the inner argv with `env SPACEDOCK_BIN=<resolvedLauncherBin>` so the value rides past `--` and survives safehouse's env-scrub. One source for the resolution (shared with `launchEnv`). No prefix on the unwrapped path; no prefix when the bin can't be resolved. The allowlist/doc-only route was considered and **rejected** (spacedock doesn't own safehouse's config).
- **Proof is split — read carefully:** the **argv-composition half is fully testable here** (AC-1…4 are pure Go argv oracles — extend the existing `safehouse_frontdoor_test.go` / `frontdoor_test.go` `want` slices). The **"survives *real* safehouse" half is NOT provable in this environment** — real safehouse does not run here, so AC-5's "end-to-end" smoke uses a *fake* safehouse stub, and that's the best this env can do. **DO NOT present th as validated on the fake smoke alone.**
- **th's validation gate = the captain's real-safehouse run, off-box.** Build + land the argv oracles (AC-1…4) and the corrected env-scrub smoke (AC-5 — note the *current* smoke "passes for the wrong reason": its fake safehouse preserves env; correct it to `unset SPACEDOCK_BIN`). Then **hand to the captain for the real-safehouse confirmation** (does the re-assert survive actual safehouse, and does `/usr/bin/env` exist in the real sandbox). That captain run is th's DoD#4, not a fake-smoke rerun.
- **Detached audit:** front-door — audit that the unwrapped path is unchanged (AC-3) and the no-bin-resolved path emits no blank value (AC-4).

### `jm` — entity-label-localization (dev-quality) — **DOWNSCOPED**
- **Scope (0.19.9):** **Layer-1 label localization ONLY.** The present-gate (`skills/present-gate/SKILL.md`) + commander-dispatch (`docs/roadmap/{sprint}/dispatch-sprint-execution.md`) templates speak the README `entity-label` (AC-1/AC-2); the shared contract (`first-officer-shared-core.md`, `ensign-shared-core.md`, the `dispatch build` ensign package) stays generic "entity" (AC-3, guarded by the existing `dispatch build` golden).
- **DEFERRED — do NOT build:** the cross-workflow `{wf}#{ref}` qualifier (**AC-4** + the two-workflow `dev`+`user-testing` fixture machinery). See jm's `## Feedback Cycles` (cycle 1): the qualifier design is banked in the body but is **out of 0.19.9 scope**. Validation requires **AC-1/AC-2/AC-3 only** — do not flag AC-4 as missing.
- **Watch:** editing the present-gate template is **broad blast-radius** (every workflow renders gates through it). The label resolves from the README via the FO's already-read `entity-label`; prove it by **live render** over an `entity-label: ticket` fixture (the captain-facing gate prose reads "ticket"), never a grep over the template. No spike needed (composes proven render paths).
- **Detached audit:** shipped scaffolding — audit that localization did NOT leak into the shared contract (AC-3 guardrail), and that a non-default-label render is genuinely driven (not a prose assertion).

## The 0.19.9 release cut (captain-gated)

After all three are `done`/PASSED + merged to `next` and `go test ./...` is green from the root:

1. Bump `version` in `.claude-plugin/plugin.json` + `.codex-plugin/plugin.json` (→ `0.19.9`) and the date-code in `.claude-plugin/marketplace.json` (`0.0.YYYYMMDDNN`).
2. Commit `release: bump version to spacedock@0.19.9`.
3. Annotated tag `v0.19.9`; push `next` + the tag → `release.yml` (goreleaser) fires.

`plugin.json` is guarded scaffolding — the **captain authorizes the cut**.

## Out of scope (and why)

- **`m1` rtk-stale-git-audit-guard** — DEFERRED (captain): disproportionate for an rtk-only, already-caught issue; the FO does the un-proxied compare ad hoc. Ideation banked.
- **`vh` survey-skill-correctness-pass** — moved to **0.19.8** (already build-ready; the 0.19.8 Commander owns it).
- **`k6` two-channel/devBranch** — flip-mechanics, moves to the flip (`pj`).
- **The cross-workflow `{wf}#{ref}` qualifier** (jm's AC-4) — deferred follow-up.
- **`xp` cross-session FO↔Commander channel** — separate design spike.
- The 0.20.0 flip itself (`pj`).
