# 0200-flip pre-flip mechanics — Commander dispatch (cold-boot)

You are the **Commander** for the **pre-flip mechanics** of sprint **0200-flip**, driving **two tasks** to **merged-on-`next`** — `nzb` (release-time e2e-gate) and `k6d` (two-channel devBranch stamp). These are the release-machinery prerequisites the 0.20.0 main-flip rides. The shaping FO carved the sprint, drove ideation, ran the preflight staff review (5 material findings folded + verified), and ran the gates; your job is the **implementation → validation → done** drive.

**You do NOT cut a release.** Unlike a normal sprint, this pair LANDS on `next` and stops there. The outward 0.20.0 flip+cut (`pj`) is driven separately by the **FO + captain** after this pair lands — do not tag, do not flip `main`, do not touch the marketplace `ref`.

This package is self-contained: boot on it cold, then drive.

## Prerequisite (met)

**0.19.9 is cut** (2026-06-08, on `next`). Boot on `next`, which now carries 0.19.9 + the full 0.19.9 cohort. Nothing else gates this pair.

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
spacedock status --workflow-dir docs/dev --where sprint=0200-flip --where 'sprint-readiness != defer'
```
→ **`nzb` gate-release-on-e2e · `k6d` two-channel-release-devbranch-stamp.** Both are ideation-approved (gates passed) with all staff-review findings folded into their bodies; drive each from `implementation`. (`pj` main-flip-0200-marketplace is `sprint-readiness: defer` — it is the capstone driven by the FO + captain, NOT you; do not pick it up.)

## Deliverable & DoD

**Deliverable:** `nzb` + `k6d` `done` / PASSED + **merged to `next`** — the release-machinery prerequisites in place for the flip. **No release cut.**

1. `nzb` + `k6d` `done` / PASSED + merged to `next`.
2. `go test ./...` from the repo root green — **the WHOLE `internal/release` package green** (nzb's AC-4 bar; see the co-edit note below).
3. `nzb`: the `e2e-gate` job + the SHA-match predicate + the waiver all proven by `internal/release` workflow-guard + predicate unit tests; no live run needed at impl time (the gate is exercised for real at the flip).
4. `k6d`: `goreleaser release --snapshot` produces BOTH channel artifacts (stable cask `spacedock` + edge `spacedock@next`) with the correct per-channel `devBranch` ldflag each; a live fresh-HOME front-door smoke confirms a `devBranch=main` binary auto-installs the `main` plugin and `devBranch=next` installs `next` (observed in install argv, both hosts — not a source-grep).

## Drive order — ⚠️ shared-file coordination

The two are conceptually independent, but **both edit `.github/workflows/release.yml`**: `nzb` adds a new `e2e-gate` job + a `needs: e2e-gate` edge on goreleaser; `k6d` changes the "Stamp plugin manifests" step's branch (`next`→`main`). Different regions of the same file → **a parallel two-worktree drive will merge-conflict on `release.yml`.** Land one, rebase the other onto it, OR drive them sequentially. They also both add `internal/release` tests (different functions, same package) — fine once `release.yml` is coherent.

**Both touch CI/release machinery** → **each earns a detached adversarial audit at validation** (README `## validation` → "Detached adversarial audit"), on a throwaway checkout of the merge result. Material findings route back through validation→implementation feedback.

## Per-task build notes

### `nzb` — gate-release-on-e2e (release-gating)
- **Scope:** add an `e2e-gate` job to `release.yml` that goreleaser `needs:`, requiring a `conclusion:success` Runtime-Live-E2E run whose `headSha` == the tagged commit; a pure SHA-match predicate as a `spacedock-release` subcommand (mirroring `journey-costs`); an auditable `SPACEDOCK_E2E_GATE_WAIVER` escape hatch. Ideation banked + spiked (the `gh run list --status success` + `headSha` query is proven against the live repo; a parked run is never `success`).
- **⚠️ CRITICAL co-edit (staff-review M2 — folded into the body):** adding `needs: e2e-gate` changes the goreleaser job header, which **breaks 14 cases across 5 functions in `internal/release/journey_workflow_test.go`** (they anchor on the byte-literal `  goreleaser:\n    runs-on: macos-latest` at ~:187/:231/:238/:274 and assert goreleaser has NO `needs` at ~:424). This is **empirically confirmed** (the edge alone → 0 passed / 14 failed). The co-edit (re-align those anchors to the new header carrying `needs: e2e-gate`, keep the journey-ledger-separation property) **MUST land in the same change** as the `release.yml` edit. AC-4's bar is "`go test ./internal/release/` green over the WHOLE package," NOT "the separation guard survives."
- **Proof:** offline only at impl time — workflow-guard tests over real `release.yml` + predicate unit tests over fixture `gh run list` JSON (pass/block/waiver). The end-to-end gate is exercised at the flip cut.
- **Detached audit:** the `e2e-gate` job + predicate — a weakened SHA-match or a dropped `needs` edge must be caught.

### `k6d` — two-channel-release-devbranch-stamp (flip-mechanics)
- **Scope:** per-channel `devBranch` stamp via **two goreleaser builds** (`spacedock-stable` ldflag `devBranch=main`, `spacedock-edge` ldflag `devBranch=next`) in one `goreleaser release`; two archives + two casks (stable `spacedock`, edge `spacedock@next`); a **single-target** channel-aware `release.yml` version-stamp branch (`next`→`main` — the stamp step at `release.yml:161`/`:172`, today hardcoded `git switch next`/`git push origin next`); REUSE `next-publish.yml` unchanged as the edge lane.
- **Riskiest mechanism PROVEN (spike done, do not re-spike):** `SPACEDOCK_DEV_BRANCH=main` → install issues `marketplace add …@main` (observed in argv); `goreleaser build --snapshot` stamps two builds with distinct ldflags in one invocation. The auto-install fan-out (`#311` claude, `z9` codex) reads `devBranch` and is correct as-is; this task only sets the value per channel.
- **⚠️ The tri-surface guard is partially red until the flip — by design (staff-review M1, folded):** AC-b's guard asserts three independently-authored surfaces agree on `main` — `release.yml` stamp target ↔ `.goreleaser.yaml` stable `devBranch` ↔ `.claude-plugin/marketplace.json` `source.ref`. **k6d lands the binary-side PAIR green now** (release.yml stamp target == goreleaser stable devBranch == `main`). The **full tri-surface `==main` goes green only when `pj` flips the marketplace ref** — it is RED-by-construction between k6d's merge and the flip (exactly like pj's paired `marketplace_manifest_test.go` edit). Do NOT try to make the full tri-surface green here; the marketplace ref is `pj`'s. The codex host's `devBranch` ldflag is its SOLE channel knob (no `.codex-plugin` marketplace calendar key) — its fresh-HOME smoke (AC-d) is the independent proof.
- **Scope boundary (do NOT cross):** k6d owns the binary side. `pj` owns the marketplace `ref` flip + the paired `marketplace_manifest_test.go` edit + the stable calendar-bump-on-`main`. Leave all three to `pj`.
- **Detached audit:** the two-build goreleaser config + the `release.yml` stamp-branch change — `goreleaser --snapshot` must really emit both channels with the right ldflags.

## Out of scope (and why)

- **`pj` main-flip-0200-marketplace** — the flip+cut itself; driven by the **FO + captain** after this pair lands. Do NOT tag, flip `main`, archive, or touch the marketplace `ref`.
- **The marketplace `ref` flip + paired guard-test edit; the stable calendar-bump-on-`main`** — all `pj`'s (see k6d's scope boundary).
- **The strict-semver `release.yml` tag guard** (staff-review P3) — an optional, unowned follow-up; not required by this pair. The archive-tag hazard is handled at the flip by a non-`v*` archive ref (`pj`).
- **The pre-cut cleanups** (`#1` checksum-gate test-strength, `#3` gofmt on `survey_sync_codex_test.go`, `#4` `hasGitEntry` cross-ref comment) — a separate cleanups task, not this pair.
- **`44` bundle / `5w` notarize / `xp` cross-session channel** — deferred / separate.

## On completion

Both `done`/PASSED + merged to `next`, `go test ./...` green (whole `internal/release` package). Report back to the shaping FO — the flip (`pj`) then runs here per `pj`'s 9-step captain-gated runbook. **Do not cut a release.**
