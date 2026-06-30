# 0230 → v0.23.0 stable cut — Commander handoff (ship the race fix, then tag)

> Self-contained handoff. A fresh Commander boots `spacedock:first-officer`, reads this + the sprint state,
> ships ONE filed fix, then drives the v0.23.0 stable cut. Everything else for the cut is DONE.
> Prior session debrief (durable, on the state branch): `docs/dev/.spacedock-state/_debriefs/2026-06-30-01.md`.

## Boot
1. `spacedock claude` → boot `spacedock:first-officer` → `status --boot`. Workflow dir: `docs/dev`. Sprint = `--where sprint=0230-stable-finalization`.
2. **Pull `main` first.** `origin/main` is at `ba18936d` (contract-2; M4 #432 + v9 #433 + #441 + #442 + the contract bump #443 all merged). NOTE: local `main` may be DIVERGED (parallel sessions committed debriefs `e1b27014`/`0031f770` directly to it). Do NOT `reset --hard`; cut from `origin/main` (the releasing.md procedure branches off `origin/main`). Surface the divergence to the captain for their owner to reconcile.

## State at handoff — v0.23.0 is ONE fix + the cut ceremony away
- **12+ members merged** this sprint. The only remaining sprint-0230 work item is the race fix below.
- **pre.4 (`v0.23.0-pre.4`) published + shaken out** — validated M4's release machinery end-to-end (e2e-gate → goreleaser → cask) + the published binary/edge-cask.
- **Contract bump (#443) on main** — the binary reports `contract 2`; both plugin manifests require `>=2,<3`; a v0.22.0 (contract-1) binary now gets a clean "upgrade your binary" abort.
- **Pre-cut audit ✅ clean** (0 ship-blockers) + a **spot-audit of the parallel #441/#442 ✅ clean** (0 ship-blockers; 2 benign majors filed: `3p` below + `w07`).
- **Value gate ⚠️ WAIVED for this tag** (captain decision): #441 grew `first-officer-shared-core` to 29181 (+595 B over the v0.22.0 baseline 28586). Clawback deferred to `kr7`/`z4`/others. Do NOT re-block the tag on the value gate — the captain waived it. (The other 9 tracked FO/ensign files are still ≤ baseline.)
- **Release notes drafted + corrected** — `/tmp/v0.23.0-release-notes.md` (may be reaped; reconstruct from the debrief's PR list — #397–#443 — and the corrected token-clawback framing: do NOT claim "every FO file at/under baseline"; say "clawed back to baseline" + note #441 re-grew shared-core).

## Remaining work (drive order)

1. **Ship `3p` `launcher-signal-forward-race-fix` (sprint 0230) — GATES the cut.**
   Audit-specified fix: `internal/cli/host_exec.go:294` starts the `forwardHostSignals` goroutine BEFORE `cmd.Start()` (:296); the goroutine reads `cmd.Process` (`host_launch_unix.go:43`) before Start writes it — a data race on the front-door teardown path (reproduced under `-race`; benign on shipped arches but real). Fix: keep `signal.Notify` before Start (buffered chan queues early signals), start the goroutine AFTER `cmd.Start()` returns. Drive: dispatch → implementation → validation (detached audit; a `-race` test that fires the signal DURING the fork/exec window, RED-first on the current ordering) → merge. The diff is `internal/cli` (launcher = dispatch/launch path → live lanes nominally required), but the fix is provable deterministically under `-race` and the live behavior is unchanged — judgment call: merge on deterministic + `-race` (like #443) with the final tagged-commit e2e as the authoritative live check, OR run the live lanes if you prefer. Full spec in the entity body.

2. **Cut v0.23.0** from the contract-2 `main` (after `3p` merges), per `docs/releasing.md` "Cutting a Stable Release":
   - Stamp `0.23.0` into `.claude-plugin/plugin.json` + `.codex-plugin/plugin.json` (currently `0.22.0`), commit on a release worktree off `origin/main`, push to `main`.
   - `manifest-tag-gate` check: `go run ./cmd/spacedock-release manifest-tag-gate v0.23.0 .claude-plugin/plugin.json .codex-plugin/plugin.json`.
   - Capture `REL_SHA=$(git rev-parse HEAD)`; **green THAT exact commit**: `gh workflow run "Runtime Live E2E" --ref main`, approve the 4 CI envs (CI-E2E/CI-E2E-CODEX/CI-E2E-PI/CI-E2E-OPUS via `gh api .../pending_deployments`), wait for green. **Expect codex re-runs** — `TestLiveCodexSharedScenarios/rejection-flow` is a reduced-rate liveness flake (took 3 attempts for pre.4); re-run-on-red, never block the tag.
   - Annotated tag `v0.23.0` on `REL_SHA` with the release notes as changelog; push → goreleaser gates on the green run + publishes + advances `stable` to the tagged SHA.
   - Verify: the tagged manifest reads `0.23.0`; a fresh `brew install` of the released cask launches (postflight strips quarantine); `spacedock --version` reports `0.23.0 (contract 2)`.

## Hard-won context (so you don't relearn it)
- **Codex live-lane flakiness:** `rejection-flow` foreground-wait stall is reduced-rate; re-run the failed lane (`gh run rerun <id> --failed`) + re-approve its env. Took up to 3 attempts. Never block the tag on it.
- **The contract-version skew (why #443 exists):** a 0.22 binary + 0.23 skills used to boot clean then break on missing verbs (`state ready`/`merge guard`), because the contract integer stayed dead at 1. #443 bumped it to 2. The durable fix (minor-version coupling) is filed as `kr7`.
- **Workflow-file pushes need the `workflow` OAuth scope** → push via `git -c credential.helper='!gh auth git-credential'`.
- **/tmp gets reaped on long (multi-day) cuts** — don't park load-bearing artifacts there; the debrief (state branch) is durable.
- **Freeze main during the cut** — parallel merges (#441/#442) moved the cut commit + value gate mid-cut. Ask the captain to hold other merges until the tag fires.
- Live e2e token: `~/.claude/benchmark-token` (rotate from keychain on a `max`-account `/login` if 401).

## Filed follow-ups (do NOT pull into v0.23.0)
`kr7` minor-version-compat-coupling · `w07` ac2-reanchor-scenario-falsifiable · still-to-file: the pre-cut audit's 2 majors (stale `docs/dev/_mods/pr-merge.md` event-loop-PR-scan cross-ref; codex-keepalive assertion steerable-by-narration) · `z4` fn-binding-refinements · `3g8` live-e2e-node-runtime-bump.

## Captain-owned
The stable tag is the irreversible step — report immediately before firing it. (The prior captain gave the FO the conn toward the cut + waived the value gate; confirm the new captain's posture.)
