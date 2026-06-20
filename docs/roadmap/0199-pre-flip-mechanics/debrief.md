---
session-date: "2026-06-08"
sequence: 4
first-commit: 83a95b86
last-commit: 432cda52
duration: ~4h
---

# Session Debrief — 2026-06-08 #4 (0199 Commander drive)

The Commander implementation→cut drive for sprint **0199-pre-flip-mechanics**: six tasks driven implementation → validation (each with a detached adversarial audit) → merge to `next`, then a SHIP-CLEAR sprint-wide staff review. Two tasks pivoted on captain real-use findings. All six shipped; `next` green; 0.19.9 ready to cut.

## Shipped
- **f1** `migration-check-share-walk-helper` — [#331](https://github.com/spacedock-dev/spacedock/pull/331). Shared the migration-check walk-step composition + widened the prune to skip the non-entity `docs/roadmap` debrief tree (also un-RED'd the suite).
- **v3** `ship-linux-binaries` — [#332](https://github.com/spacedock-dev/spacedock/pull/332). Ship Linux binaries + a universal `curl|sh` install path (the release was darwin-only). Closes #321.
- **jm** `entity-label-localization` — [#333](https://github.com/spacedock-dev/spacedock/pull/333). The FO's captain-facing prose speaks the workflow's declared `entity-label`; the shared contract mechanics stay generic.
- **th** `safehouse-preserves-spacedock-bin` — [#334](https://github.com/spacedock-dev/spacedock/pull/334). A safehouse-wrapped launch preserves the launcher-injected `SPACEDOCK_BIN` via safehouse `--env-pass` (inner program stays `claude`/`codex` so the host profile auto-detects).
- **47** `survey-codex-and-sandbox-followups` — [#335](https://github.com/spacedock-dev/spacedock/pull/335). Survey body pass: Codex workstreams (workdir-attributed) + worktree→logical work-area + de-narrated scaffold + no pre-body stop + two-mode commission classification.
- **yq** `frontdoor-launch-ux` — [#336](https://github.com/spacedock-dev/spacedock/pull/336). Front-door UX: honest auto-install message + a pre-launch banner that names the real workflow (host-neutral noise prune); restored personal prompt; deleted a false `requires-contract` doctor note.

## Filed (backlog)
- None filed this session — 47rx / yq / m1 were filed in the prior 0199 shaping sessions (#01–03).

## Non-PR commits (workflow-only)
- `b1bdcf19` 0199 preflight staff review on `next` (shaping FO, before the drive).
- `6ead4fa6` quote the 0198 debrief `session-date` scalar (shaping FO's cut-blocker band-aid; f1 also fixed it at root).
- Feedback cycles: **th cycle 1** (argv-prefix → `--env-pass` after the captain's real-safehouse catch); **yq cycle 1** (test-strength — isolate `hasGitEntry`).
- `z6` `pi-stage-dispatch-uses-build-artifact` — terminalized + archived at boot (PR #300 merged a prior session, left mid-ceremony).

## Decisions
- **th mechanism pivot (captain-caught).** The cycle-1 `/usr/bin/env SPACEDOCK_BIN=<bin>` argv prefix passed all in-env oracles, but the captain's real-safehouse run found it masked the inner program from safehouse → broke claude's profile auto-detection (lost `~/.claude`). Pivoted to safehouse's `--env-pass SPACEDOCK_BIN` (inner program stays `claude`/`codex`); added **AC-6** to guard the host-program invariant; captain re-test confirmed survival + profile preserved.
- **yq live-use reverts.** From real `spacedock claude` runs: restored the personal default bootstrap prompt (reverted the D-neutralization), fixed the banner to find the real `docs/dev` from any launch dir (walk-up missed subdirs; `--discover` was polluted by agent-worktree copies), and deleted the `requires-contract` doctor note (the Claude Code guide confirmed the runtime is silent on unknown plugin.json fields — the note's "load-time warning" claim was false). The discovery noise prune is host-neutral (`.gitignore` dir-patterns + a nested-`.git`-checkout skip), not a `.claude` literal.
- **47rx option-b merge.** Merged on the proven mechanisms (query-smoke artifacts + the both-mode fixture + the classifier audit); **AC-7b/8b's real-mixed-mode render deferred to a captain torahmap spot-check** (the ideation-authorized escalation; torahmap absent from this box, this repo single-mode).
- **The conn.** The captain delegated gate-approval + merge + CI-approval to the Commander mid-drive; the six merges + the sprint-wide staff review ran autonomously (`next` is unprotected → the `offline` gate gates merges; live lanes left unapproved to save quota).
- **m1 deferred** (rtk-only, already-caught — disproportionate to guard); **agy-runtime-support** orphan left flagged (parked mid-implementation, no sprint).

## Issues — Workflow
Sprint-wide staff review = **SHIP-CLEAR** (0 confirmed blockers; suite `go test ./...` 1208 passed / 16 packages, `go vet` + `go build` clean). Non-blocking polish to follow up:
- **Checksum-gate proof isn't go-test-guarded.** install.sh's fail-closed proof lives only in `.github/workflows/install-e2e.yml`, which no `*_test.go` loads — a contributor could delete the checksum check and `go test ./...` stays green (only the install-e2e PR check reds). Add a `workflow_exec_guard`-style test asserting install-e2e.yml carries the tamper step, or a Go end-to-end tamper test.
- **Stale darwin-only prose** in `docs/releasing.md` + the `release.yml` header comment (v3 added linux to the goreleaser build; the release docs still say darwin-only). Doc fix.
- Pre-existing cosmetic `gofmt` hit on `skills/integration/survey_sync_codex_test.go` (a fancy-quote in a comment) — `gofmt -w` to clear (the only gofmt issue on `next`).
- `hasGitEntry`'s only *unmasked* guard is the single `TestDiscoverWorkflowsSkipsNestedCheckout` (the cli banner fixtures double-mask it) — worth a cross-reference comment so a future divergence doesn't silently drop the guard. Minor: the phantom-manifest `NoPluginFound` arm isn't specifically tested; the banner deep-subdir subcase is redundant with the repo-root subcase.

## Issues — Spacedock
- **Team-harness supersede race (minor, not filed).** Fresh-dispatching a validator with the *same name* immediately after superseding the old one races the old one's async shutdown — the `teammate_terminated`/`shutdown_approved` for that name arrived *after* the new spawn. It resolved cleanly (verified the fresh validator alive via the roster + the context-budget probe), but reusing the exact name across a supersede is fragile; a `-cycleN` suffix on the fresh dispatch would avoid it. Candidate spacedock issue — not filed (awaiting captain nod).

## Observations
- The detached adversarial audit earned its keep twice: it caught yq's `hasGitEntry` test-strength hole (a load-bearing function deletable with the suite staying green), and th's real-safehouse verdict-hold caught a ship-blocker that 686 in-env tests sailed past. The "captain runs the un-fakeable leg" pattern (th real-safehouse, 47rx torahmap) is the honest proof bar for environment-gated behavior.
- Two of six tasks pivoted on captain *real-use* findings (th, yq) — the live front-door runs surfaced what the test suites couldn't. Worth front-loading a real launch earlier for front-door-touching tasks.

## What's Next
- **0.19.9 cut** — captain-authorized, imminent (bump the plugin manifests + marketplace date-code, commit, tag `v0.19.9`, push `next` + tag → `release.yml`).
- **Deferred:** m1 `rtk-stale-git-audit-guard` (ideation); `agy-runtime-support` (parked orphan, leave flagged); 47rx **AC-7b/8b** captain torahmap `/spacedock:survey` spot-check (fast-follow).
- **Polish follow-ups** (from the staff review, above): the checksum-gate go-test/workflow-guard proof; the `docs/releasing.md` + `release.yml` darwin-only prose; `gofmt -w` the survey test; the `hasGitEntry` cross-reference comment.
