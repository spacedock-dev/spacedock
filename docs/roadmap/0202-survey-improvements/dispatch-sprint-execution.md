# 0202 (0.20.2) — survey + cleanups — Commander dispatch (cold-boot)

> **Status: APPROVED 2026-06-13 (captain) — ready for Commander cold-boot drive.** Five members, all staff-clean (`staff-review.md`: NOT-READY → `5wv` M1 folded → clean; `nd`/`td`/`5ar`/`gf` 0 Material). Gates approved; the Commander advances each ideation→implementation→validation→done. Drivable set: `spacedock status --workflow-dir docs/dev --where sprint=0202-survey-improvements --where 'sprint-readiness != defer'`.

## Boot
```bash
git fetch origin main && git switch -c drive/0202 origin/main && go build -o ./spacedock ./cmd/spacedock
export SPACEDOCK_BIN="$PWD/spacedock" SPACEDOCK_REPO_ROOT="$PWD"
git -C docs/dev/.spacedock-state pull --rebase origin spacedock-state/dev   # gh-HTTPS if SSH down
security find-generic-password -s "Claude Code-credentials" -w | python3 -c "import sys,json; print(json.load(sys.stdin)['claudeAiOauth']['accessToken'])" > ~/.claude/benchmark-token
./spacedock status --workflow-dir docs/dev --boot
```
**Worktrees + PRs target `main`** (post-flip trunk). **SSH is down** — push via `git -c credential.helper='!gh auth git-credential' push https://github.com/spacedock-dev/spacedock.git <ref>`.

## ⚠️ DEADLINE — do `5ar` item #5 FIRST
The node-20 GitHub Actions deprecation hits **~2026-06-16** (days out). `5ar`#5 bumps `actions/checkout@v4→v5`, `actions/setup-go@v5→v6`, `goreleaser/goreleaser-action@v6→v7` (**v6 stays node-20!**) on the cut path (`release.yml:175`, `install-e2e.yml:39`); `deploy-pages`/`upload-pages-artifact` (docs.yml only) have no node-24 release yet — leave them. Land this before the window closes regardless of the rest of the sprint.

## Deliverable & DoD
**0.20.2** = the survey output redesign + dev cleanups. Done when, merged to `main`:
- **`5wv`** the `spacedock:survey` step-4 report is rebuilt in the "value & numbers first" structure (`docs/roadmap/0202-survey-improvements/index.md` mock), every band (R1–R6) proven by a survey run over a constructed fixture, the expected value derived from fixture rows — never a SKILL.md prose-grep. Two NEW queries (`dispatch-fact`, `decision-no-followup`); the fixture is EXTENDED with `tool_calls.message_id` + `messages.ordinal` (R1b) for the no-follow-up join.
- **`nd`** `status --next-id` prints a use-`new` hint on stderr (stdout id unchanged); the FO contract teaches `spacedock new` (proven by a live FO drive).
- **`td`** `status --validate` enforces full mdschema conformance (per-field types/patterns/enums), warn-tier honoring the read-path lockout.
- **`5ar`** five items: checksum-gate tamper test, node-action bump (above), gofmt, darwin doc drift (release.yml header), hasGitEntry comment.
- **`gf`** split-root state-sync degrades to local-only when the state checkout has no `origin` (keep path-scoped commits, skip push/pull, surface "not remotely synced").

## Drive order — ⚠️ coordination
- **`5ar`#5 first** (deadline).
- **`5wv` owns the survey skill alone** (SKILL.md + queries.sql + testdata/survey fixture) — it does NOT collide with the cleanups (different files), so it parallelizes.
- **`internal/status` is the shared package** — `nd` (`--next-id` hint), `td` (`--validate`), `gf` (`boot.go`/`json_commands.go`) all touch it. Sequence or rebase carefully; `go test ./internal/status/` green over the WHOLE package per change. **`td` watch:** the warn-tier stderr must not perturb the exact-match `stdout=="VALID"` golden in `native_validate_test.go` (staff P3).

## Per-member build notes
### `5wv` — survey-output-redesign (survey)
Rewrite the step-4 report template to the locked structure. R1 plain value+numbers + "manual" not "mechanical" + threads-to-pull; R2 recent-window snapshot + one partial-lens caveat; R3 subagent-dispatch fact (new query); R4 collapse empty Codex `(unlabeled)` + strip scratch preamble; R5 knowledge-work archetype (zb#1 Codex-count-first already shipped #335 — do NOT redo); R6 mode-aware framing + conditional branch-aware work-by-area. **AC-1b:** extend the fixture with `message_id`/`ordinal` and prove no-follow-up with a non-vacuous higher-ordinal-Edit mutation. Keep-signal: do NOT regress the inference accuracy or decision-frontier triage.
### `nd` — prefer-new-over-next-id (cleanup) · shipped-contract surface
AC-1 stdout/stderr separation test; AC-2 a live FO drive observing filing via `new` (not a prose-grep). Seams: `native_runner.go`, `new.go`, `cli.go:308`.
### `td` — mdschema-conformance-validator (cleanup) · status-guard surface
Extend `internal/status/validate.go` to full mdschema (mod-block `^[^:]+:[^:]+$`, verdict enum, etc.); warn-tier per `validate.go:138-145` read-path lockout. Known-bad fixtures assert non-zero exit + diagnostic.
### `5ar` — pre-cut-audit-cleanups-0199 (cleanup) · CI/release-machinery surface
#1 checksum tamper test (reds when `install.sh:160-169` gate is stripped; model on `TestGoreleaserBuildGuardRejectsDroppedLinux`). #2 release.yml header darwin→darwin+linux (docs/releasing.md already correct). #3 gofmt — REWORD the `''` comment before `gofmt -w` (go1.26.1 rewrites it to a curly quote). #4 hasGitEntry comment. #5 node-action bump (above).
### `gf` — state-sync-no-origin-local-mode (cleanup)
No-origin detection + local-only mode. Seams real: `boot.go:283-284`, `json_commands.go:195-198`, `build.go:719/498/508`, `state.go:178`. Contract delta: the FO/ensign push/pull guidance gains a no-remote branch.

## Detached adversarial audit (before merge)
High-stakes surfaces here: **`td`** (status mutation/guard), **`5ar`** (CI/release machinery), **`nd`** (shipped contract scaffolding). Run a read-only detached audit on each before merge. `5wv`/`gf` are routine validation.

## Pre-cut antipattern audit (⚠️ before the tag)
All five merged to `main`, `v0.20.2` not yet fired → INDEPENDENT staff-eng reviewer over the assembled sprint. Verify main-PR CI gating exists (same cross-cutting gap flagged in 0.20.1). Ship-blockers fixed pre-cut; non-blockers seed the next sprint.

## Cut
`go test ./...` green from root, then `docs/releasing.md`. Cut `v0.20.2`. Captain authorizes.

## Out of scope (deferred)
The six survey source seeds (subsumed into `5wv` as R1–R6); the broader backlog (`j9`/`ey`/`vc`/`44`/`bw`/`95`/`5wc`/`82`/`xp`/`6a`/`vv`/`e6`).
