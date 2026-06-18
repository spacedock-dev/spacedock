# m4 readiness — both-semantics live team-mode harness

> FO package for finishing **m4** (`live-team-mode-terminal-harness`, PR #390). Captures the both-semantics plan, the legacy-lane diagnosis + probe, the rebase status, and the merged-lane test spec. Authored 2026-06-18 (shaping FO), against 9243 (`using-claude-team-merged-model-support`) crunching in parallel.

## The both-semantics seam (captain requirement: "m4 must work in either old and new claude-team semantics")

Per 9243's OQ-6 host-regime split, m4 ends up with **two live lanes** (plus an explicitly-forced bare lane tracked separately as `bare-mode-coverage`):

- **Legacy lane** — pinned claude **2.1.177**, **interactive tmux/pty**, native `TeamCreate`. The current PR #390 harness. This lane exists *only* to keep the 2.1.177 pin alive (9243's deprecation trigger retires the legacy path when no live lane drives it). Currently **CI-RED on Linux** (see diagnosis).
- **Merged lane** — current/**unpinned** claude (≥2.1.178), **headless `claude -p`, no tmux**, auto-spawn team (no `TeamCreate`). Works flag-free headless (9243 OQ-2, reproduced this session). Sidesteps the interactive-pty-auth wall that fails the legacy lane. **Not yet authored** (spec below); base is rebased-in, pending 9243 settling.

## Legacy-lane diagnosis (CI-RED) — pinned to a fork, probe ready to settle it

Both pty tests (`TestLivePtyStandingResidencyInjectsCommOfficer`, `TestLivePtyEnsignCycleTeamTeardown`) fail identically at the FIRST gate `waitForSessionFile`: the interactive child writes **no session jsonl** and the FO pane is **0 bytes**, timing out at 4 min. The tmux session stays **alive** the whole time (timeout, not death) → the child **blocks on a pre-transcript screen**. Render is refuted (real 220×50 TTY; idle gate reads on-disk jsonl, not pane text). Headless steps 13/14 pass with the same `ANTHROPIC_API_KEY` → the credential is valid; the divergence is interactive-only. On Linux the stored-login seed is skipped (`seedStoredLoginCredential` is macOS-keychain-only), so the child carries only the env credential interactive claude may not honor as a non-blocking login.

**The fork (unsettled by the existing artifact):** (a) **auth** — interactive 2.1.177 sits at a "Claude API"/login banner; vs (c) **onboarding** — a first-run gate `seedInteractiveClaudeConfig` misses. Why the prior pane dumps were empty: the harness `captureFOPane` title-resolves the FO pane (never titled → `""`), masking the banner.

**Probe (ready to run, not triggered):** branch `spacedock-ensign/m4-pty-auth-probe` (off pushed `31d22d25`; zero 9243 commits) adds `internal/ensigncycle/pty_auth_probe_test.go` (`TestPtyAuthProbe`, gated on `SPACEDOCK_PTY_AUTH_PROBE=1`, 25s budget). On timeout it dumps the **active** pane raw to `live-artifacts/claude/<model>/pty-team-mode/auth-probe/`: `probe-pane.ansi` (`capture-pane -p -e`), `probe-pane-info.txt` (pane pid/command/title), `probe-state.txt`, `probe-pane-fo.txt`. Reading `probe-pane.ansi`: auth/login banner ⇒ (a); theme/trust/login picker ⇒ (c); blank while pane-info shows node/claude live ⇒ render/stdin.

**To run (needs maintainer CI-E2E approval):** `Runtime Live E2E` workflow, `claude-live` job, with `SPACEDOCK_PTY_AUTH_PROBE=1` set, invoking `go test -tags live -run TestPtyAuthProbe ./internal/ensigncycle/`. The existing "Upload live artifacts" step (`if: always()`, globs `live-artifacts/claude/**`) uploads the dumps with no workflow change — only a run step on one model leg. Once the fork is pinned, the targeted fix (Linux stored-login seed from a CI secret for (a), or an extended onboarding seed for (c)) follows.

## Rebase status

- m4's **local** branch is rebased onto 9243's committed tip `a94e4825` (merged-mode `build.go`/`standing.go` spliced under m4's harness commits — `go build ./...` OK, 9243's dispatch unit tests green, m4's harness compiles). 9243 preserved legacy emission byte-identically, so m4's legacy tests stay valid.
- This base is **provisional + unpushed**: 9243 is local-only and being edited in another session. m4 re-rebases when 9243 lands. The legacy probe was branched off the **pushed** tip `31d22d25` to stay pushable without dragging in 9243's unpushed commits.

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

## Next steps

1. **Pin the legacy fork** — captain approves a single-model CI-E2E run of `TestPtyAuthProbe`; read `probe-pane.ansi`; apply the targeted fix.
2. **Author the merged lane** — once 9243 lands (its contract prose settles + it merges to main), re-rebase m4 onto main and author the merged-lane test per the spec above.
3. **Bare lane** — tracked separately (`bare-mode-coverage`), sequenced after these + tied to pin-retirement.
