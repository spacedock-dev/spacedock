# m4 readiness — both-semantics live team-mode harness

> FO package for finishing **m4** (`live-team-mode-terminal-harness`, PR #390). Authored 2026-06-18 (shaping FO).
>
> **UPDATE 2026-06-18 (captain re-scope + 9243 merged):** 9243 (`using-claude-team`, #396) is **MERGED** on main and green on Claude Code 2.1.181 locally. The release plan re-scoped: **0.20.5 = m4's merged lane green on unpinned Claude + the CI unpin** (the 0205-layered-fo sprint bumped to 0.20.6). m4's deliverable is the **MERGED lane**; the legacy 2.1.177 interactive lane is **NOT fixed** — legacy support is **best-effort in code** (the conditionally-loaded `using-legacy-claude-team` path), retired when STABLE Claude Code catches up to the merged floor. **The legacy auth diagnosis + probe below are OFF the critical path** (kept for posterity; do not gate the cut on them). m4's local branch is **re-rebased onto main (#396)** — build green.

## The both-semantics seam (captain requirement: "m4 must work in either old and new claude-team semantics")

Per 9243's OQ-6 host-regime split, m4 ends up with **two live lanes** (plus an explicitly-forced bare lane tracked separately as `bare-mode-coverage`):

- **Merged lane (the 0.20.5 deliverable)** — current/**unpinned** claude (≥2.1.178), **headless `claude -p`, no tmux**, auto-spawn team (no `TeamCreate`). Works flag-free headless (9243 OQ-2, reproduced this session). Sidesteps the interactive-pty-auth wall that fails the legacy lane. **Not yet authored** (spec below); base is rebased onto main #396 — authorable now that 9243 is merged.
- **Legacy lane (best-effort, NOT fixed)** — pinned claude **2.1.177**, **interactive tmux/pty**, native `TeamCreate`. The current PR #390 harness, **CI-RED on Linux** (auth/onboarding; see diagnosis). Per the re-scope it is **not fixed and not kept as a separate pinned job** — legacy support stays best-effort in 9243's conditionally-loaded `using-legacy-claude-team` path. When CI unpins in 0.20.5, these interactive pty tests should **skip on a merged host** (`ToolSearch(select:TeamCreate)` empty ⇒ no TeamCreate ⇒ skip, not fail) rather than RED, and the legacy path retires when stable Claude catches up to the merged floor.

## Legacy-lane diagnosis (CI-RED) — OFF the critical path (posterity only)

> The captain re-scope (above) does NOT fix the legacy lane. The diagnosis + probe below are retained for whoever later needs the legacy interactive path to work; they do not gate the 0.20.5 cut. The actionable legacy step now is to **skip these tests on a merged host**, not fix them.

Both pty tests (`TestLivePtyStandingResidencyInjectsCommOfficer`, `TestLivePtyEnsignCycleTeamTeardown`) fail identically at the FIRST gate `waitForSessionFile`: the interactive child writes **no session jsonl** and the FO pane is **0 bytes**, timing out at 4 min. The tmux session stays **alive** the whole time (timeout, not death) → the child **blocks on a pre-transcript screen**. Render is refuted (real 220×50 TTY; idle gate reads on-disk jsonl, not pane text). Headless steps 13/14 pass with the same `ANTHROPIC_API_KEY` → the credential is valid; the divergence is interactive-only. On Linux the stored-login seed is skipped (`seedStoredLoginCredential` is macOS-keychain-only), so the child carries only the env credential interactive claude may not honor as a non-blocking login.

**The fork (unsettled by the existing artifact):** (a) **auth** — interactive 2.1.177 sits at a "Claude API"/login banner; vs (c) **onboarding** — a first-run gate `seedInteractiveClaudeConfig` misses. Why the prior pane dumps were empty: the harness `captureFOPane` title-resolves the FO pane (never titled → `""`), masking the banner.

**Probe (ready to run, not triggered):** branch `spacedock-ensign/m4-pty-auth-probe` (off pushed `31d22d25`; zero 9243 commits) adds `internal/ensigncycle/pty_auth_probe_test.go` (`TestPtyAuthProbe`, gated on `SPACEDOCK_PTY_AUTH_PROBE=1`, 25s budget). On timeout it dumps the **active** pane raw to `live-artifacts/claude/<model>/pty-team-mode/auth-probe/`: `probe-pane.ansi` (`capture-pane -p -e`), `probe-pane-info.txt` (pane pid/command/title), `probe-state.txt`, `probe-pane-fo.txt`. Reading `probe-pane.ansi`: auth/login banner ⇒ (a); theme/trust/login picker ⇒ (c); blank while pane-info shows node/claude live ⇒ render/stdin.

**To run (needs maintainer CI-E2E approval):** `Runtime Live E2E` workflow, `claude-live` job, with `SPACEDOCK_PTY_AUTH_PROBE=1` set, invoking `go test -tags live -run TestPtyAuthProbe ./internal/ensigncycle/`. The existing "Upload live artifacts" step (`if: always()`, globs `live-artifacts/claude/**`) uploads the dumps with no workflow change — only a run step on one model leg. Once the fork is pinned, the targeted fix (Linux stored-login seed from a CI secret for (a), or an extended onboarding seed for (c)) follows.

## Rebase status

- 9243 landed on main as **#396** (squash commit `1d691b45`). m4's local branch is **re-rebased onto main** (`git rebase --onto origin/main a94e4825` — the 11 harness commits replayed onto main, the 3 stale pre-squash merged-mode commits dropped since #396 supersedes them). `go build ./...` OK. m4 touches only `internal/ensigncycle/*` + the CI yaml — disjoint from #396's dispatch/skill files, so the replay was clean.
- **Local + unpushed.** PR #390 still points at the old pushed tip `31d22d25`; the re-rebased branch will be force-pushed when m4's merged lane is authored (a Commander/ensign step, not a shaping-FO push). The legacy auth-probe branch `spacedock-ensign/m4-pty-auth-probe` (off `31d22d25`) is parked, off the critical path.

## Merged-lane test spec (authoring deferred until 9243 settles its contract prose)

**Lane shape:** a headless `claude -p` FO on current/unpinned claude, `--output-format stream-json --verbose`, nested-session env markers `env -u`-scrubbed, **no tmux, no pty, no `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS`** (merged channel is flag-free, resident through spawn→deliver→reap in one `-p` run). Fixture: same `entityFixture()` + `readmeRealisticLifecycle()` (+ `_mods/comm-officer.md` for the AC-7 variant); auth via `isolatedClaudeEnv` (skip-not-fatal without auth).

**Assertions (all rest on shapes 9243 has committed + live-validated):**
1. **No `TeamCreate`** in the FO stream (tool absent on merged host; must not attempt it, must not fall to bare/sequential).
2. **`Agent(name=…, run_in_background=true)`** — at least one merged-shape dispatch (`name` = lead→worker handle; `run_in_background` = worker→lead back-channel).
3. **Auto-team registry** `~/.claude/teams/session-<id>/config.json` (under the run's isolated `CLAUDE_CONFIG_DIR`, matched by `leadSessionId == $CLAUDE_CODE_SESSION_ID`) gains a member with `agentType:"spacedock:ensign"`, `backendType:"in-process"`, **no `team_name`**.
4. **Reconcile** `spacedock dispatch reconcile` (no `--team-name`) auto-discovers the live roster by session id, exit 0.
5. **Teardown** via per-name `SendMessage(shutdown_request)` (no `TeamDelete`, no `TERMINAL_TEARDOWN_BOUNDED`); observe `teammate_terminated` per ensign and `members[]` pruning to empty.
6. **AC-7 variant** `spacedock dispatch spawn-standing-all` (no `--team`) injects `comm-officer` as a named background `Agent`; registry gains a `comm-officer` member (`agentType:"general-purpose"`, `model:"sonnet"`).

**The ONE parameter to fill (not a hard blocker):** the worker→lead completion-signal target. 9243's committed `build.go` emits `SendMessage(to="team-lead")` (AC-6 validated on `team-lead`), but `claude-fo-dispatch.md`'s worker-back-channel prose + AC-1 description say `to="main"` (both route to the lead). Assert against the **committed code target** (`team-lead`) as a parameter; flip only if 9243 reconciles the prose/code split to `main`. Do not assert both.

**Live coverage = a SECOND live-e2e lane on unpinned/current claude** (9243's own recommendation): the merged lane becomes the regression sensor for the merged floor; the pinned-2.1.177 lane stays the legacy sensor. That two-lane split is the live realization of 9243's OQ-6 seam.

## Next steps (the 0.20.5 cut)

0.20.5 = **m4 merged lane green on unpinned Claude + CI unpinned.** The work:

1. **Author the merged lane** — re-rebase done (m4 on main #396). Author the merged-lane test per the spec above (completion target = `team-lead`, per #396's committed `build.go`). Implementation work → a dispatched ensign (or captain-greenlit FO-driven, as m4 has been).
2. **Skip-gate the legacy lane on merged hosts** — the interactive pty tests skip when `ToolSearch(select:TeamCreate)` is empty (merged host), so unpinning makes them SKIP, not RED. Legacy stays best-effort.
3. **Unpin CI** — filed as `ci-unpin-claude-version` (see below): set the live-e2e lane to current/unpinned Claude (≥2.1.178 merged floor), retire the #395 keystone pin + its `claude_version_pin_guard_test.go` enforcement, keep `DISABLE_AUTOUPDATER` semantics sane for a floating version. The pin-guard test currently REDs on any non-2.1.177 pin, so the unpin must update/retire that guard in the same change.
4. **Cut 0.20.5** after the merged lane is green on unpinned CI and the pre-cut antipattern audit is clean.

**Deferred to 0.20.6+:** the legacy auth fix (off critical path, posterity probe parked); `bare-mode-coverage` (forced-bare scenario + `-p`-assumes-bare audit, tied to pin-retirement); the layered-FO verb core / tier / restructure / haiku-drive (the bumped 0205 sprint).

## Folded finding (captain 2026-06-19): contract correction rides in m4 (no separate task)

The headless-registry finding is **folded into m4** rather than filed separately. When the current validation completes, m4 takes one more implementation increment (before its gate/terminal) to: (a) **correct `skills/first-officer/references/claude-fo-dispatch.md` ~line 7** — its claim that "Claude Code writes `~/.claude/teams/session-<id>/config.json` automatically" is false on a headless `-p` merged host (the member record is `projects/<cwd>/<session-id>/subagents/agent-*.meta.json`; the `teams/config.json` registry is the interactive artifact); and (b) **document the consequent reconcile-headless limitation** — the FO's `reconcile` leadSessionId auto-discovery degrades to git-only on a headless `-p` merged host (no `teams/config.json` to glob), noting a `reconcile.go` fix (read the meta.json path for headless roster discovery) as a 0.20.6 candidate for the layered-FO reconcile/next-action work, NOT done here. This expands m4's PR to touch a shipped-contract file (a high-stakes surface) — so the validator's detached adversarial audit applies. Sequencing: fold AFTER the current validator completes (editing the contract file now would race its `go test`, which `contractlint` reads skill files during).
