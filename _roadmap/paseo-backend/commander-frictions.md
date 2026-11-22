# Sprint `paseo-backend` — Commander Friction Log

Running log of real frictions encountered while cold-booting and operating sprint `paseo-backend` as its own design (`paseo` CLI dispatch on the Safehouse-wrapped `127.0.0.1:6769` daemon).

## Summary

| # | Phase | Friction | Severity | Status |
|---|---|---|---|---|
| F1 | Cold boot (Step 1) | `PATH` resolves `paseo` to unbuilt worktree `node_modules/.bin/paseo` (`ERR_MODULE_NOT_FOUND: packages/cli/dist/index.js`) | High | Deferred (dev daemon env) |
| F2 | Cold boot (Step 1) | `/Users/clkao/.local/bin/paseo` → `/Applications/Paseo.app/Contents/Resources/bin/paseo` fails with `Operation not permitted` under Safehouse — **no executable `paseo` CLI exists in this Commander sandbox** | **Blocker** | Open (escalated to Captain) |
| F3 | Cold boot (Step 1) | `PASEO_AGENT_ID` is absent from Commander's shell environment (`PASEO_HOME` and `APP_SANDBOX_CONTAINER_ID` are set, but not `PASEO_AGENT_ID`) | High | Open (needs Paseo/agy env propagation) |
| F4 | Cold boot (Step 2) | `/opt/homebrew/bin/spacedock --version` (`0.29.0-pre0`) omits `contract 3` on Antigravity because `runtimehost.Detect` returns empty (`internal/cli/cli.go:916`), tripping `first-officer-shared-core.md` binary gate | High | Deferred (follow-up candidate) |
| F5 | Cold boot (Step 2) | Global `~/.gemini/config/plugins/spacedock/` (Jun 5) drifts from repo `skills/first-officer/` (`SKILL.md` and `first-officer-shared-core.md` differ) | Medium | Deferred (tracked outside sprint) |

---

## Detailed Entries

### F1 — Unbuilt dev-worktree `paseo` ahead on `PATH`
- **Observed:** `which -a paseo` returned three paths all resolving to `/Users/clkao/git/spacedock-research/paseo/.worktrees/integration/node_modules/.bin/paseo` (symlink to `../@getpaseo/cli/bin/paseo`).
- **Error:** Running `paseo --version` and `paseo status --json` crashed immediately with:
  ```text
  Error [ERR_MODULE_NOT_FOUND]: Cannot find module '/Users/clkao/git/spacedock-research/paseo/.worktrees/integration/packages/cli/dist/index.js' imported from .../packages/cli/bin/paseo
  ```
- **Secondary probes (read-only):**
  - `tsx .../packages/cli/src/index.ts --version` failed with `Error: listen EPERM: operation not permitted /var/folders/.../T/tsx-501/*.pipe` under Safehouse.
  - `bun .../packages/cli/src/index.ts --version` failed with `Cannot find module '@getpaseo/server/process'` (workspace `dist/` unbuilt).
- **Disposition:** Matches `staff-review.md` Commander Cold-Boot Readiness Note #2. Fixable outside sandbox by building `packages/cli` (and workspace deps) in the Paseo integration worktree, or installing a Safehouse-readable `paseo` binary.

### F2 — Safehouse denies `/Applications/Paseo.app` (`Operation not permitted`) — Boot Step 1 Hard Blocker
- **Observed:** `dispatch-sprint-execution.md` expected `--label paseo-safehouse.paseo-access=daemon` to grant read-only access to `/Applications/Paseo.app` so `PASEO_BIN=/Users/clkao/.local/bin/paseo` works.
- **Error:**
  ```text
  ls: /Applications/Paseo.app/Contents/Resources/bin/paseo: Operation not permitted
  /bin/sh: /Users/clkao/.local/bin/paseo: Operation not permitted
  ```
  Checked every directory on `$PATH` plus `~/.local/share/fnm/...`, `~/.bun/bin`, `/opt/homebrew/bin`, `/usr/local/bin` — no other `paseo` binary exists on this host.
- **Impact:** Boot Step 1 (`paseo --version && paseo status --json`, confirming `connectedDaemon == "reachable"`) cannot pass. Commander cannot dispatch Wave 1 (`936`) via Paseo on this daemon without either (a) building `packages/cli` (and its workspace deps) in the Paseo integration worktree (which is already on `PATH` and readable inside Safehouse), (b) adjusting `paseo-safehouse`'s `paseo-access=daemon` profile to allow `/Applications/Paseo.app`, or (c) installing a standalone `@getpaseo/cli` binary in a Safehouse-readable path (e.g. `~/.local/bin` or fnm global bin).
- **Disposition:** Hard blocker — escalated immediately to Captain per Boot Step 1 fail-fast and Escalation rules.

### F3 — `PASEO_AGENT_ID` missing from Commander shell environment
- **Observed:** `env` inside Commander's shell contains `PASEO_HOME=/Users/clkao/git/spacedock-research/paseo/.worktrees/integration/.dev/paseo-home` and `APP_SANDBOX_CONTAINER_ID=agent-safehouse`, but `PASEO_AGENT_ID` is completely absent.
- **Impact:**
  1. `skills/first-officer/SKILL.md`'s planned `PASEO_AGENT_ID` auto-selection rule will not trigger automatically on Antigravity (requires explicit captain selection unless fixed).
  2. `spacedock dispatch build --host paseo`'s `PASEO_AGENT_ID`-based `/tmp/spacedock-dispatch/` filename disambiguation will fall back to plain `{derivedName}.md`.
  3. `paseo run` child spawns will not auto-populate `ParentAgentId` via Paseo CLI's `resolveCallerAgentId`, which would also affect `TestLivePaseoDispatchSmoke`'s `ParentAgentId == foRun.AgentID` assertion if tested with an `antigravity` FO.
- **Disposition:** Needs investigation in Paseo's Antigravity provider / `antigravity-cli` tool-shell env forwarding.

### F4 — `spacedock --version` omits `contract 3` on Antigravity
- **Observed:** `/opt/homebrew/bin/spacedock --version` outputs:
  ```text
  spacedock 0.29.0-pre0
  OS: darwin/arm64
  Channel: edge (spacedock@spacedock-edge)
  ```
- **Cause:** `internal/cli/cli.go:916-921` returns early after 3 lines (omitting `Runtime:`, `Sandbox:`, and `contract 3`) whenever `runtimehost.Detect(getenv)` returns empty. Because `runtimehost.Detect` only recognizes `CLAUDECODE`, `CODEX_THREAD_ID`, and `PI_CODING_AGENT_DIR`, an Antigravity (`agy`) session looks like "outside every runtime — a human at a terminal".
- **Impact:** Strictly following `skills/first-officer/references/first-officer-shared-core.md:9` ("the `--version` output must contain `contract 3`. A different level, or a missing token → ABORT") aborts every cold-booted Antigravity FO before boot!
- **Disposition:** Real Spacedock bug for Paseo-hosted non-{Claude,Codex,Pi} FOs; note that `936` explicitly keeps `PASEO_AGENT_ID` out of `runtimehost.Detect` to avoid ambiguity with inner hosts, so `printVersion` should consider emitting `contract 3` when `PASEO_AGENT_ID` or `PASEO_HOME` is present (or always/separately) in a follow-up task.

### F5 — Global `~/.gemini/config/plugins/spacedock` drift
- **Observed:** `~/.gemini/config/plugins/spacedock/` timestamps are `Jun 5 14:04`, and `diff -q` confirms both `skills/first-officer/SKILL.md` and `skills/first-officer/references/first-officer-shared-core.md` differ from `spacedock-v1/skills/first-officer/`.
- **Disposition:** Known outside-sprint friction (`_roadmap/paseo-backend/index.md:68`); Commander loaded the repo-local contract directly as instructed.

### Verified Working at Boot
- `GIT_SSH_COMMAND='ssh -o BatchMode=yes'` works cleanly (already exported in environment; `git -C docs/dev/.spacedock-state fetch --dry-run origin` succeeded with no passphrase hang).
- `spacedock status --boot --identify --json`, `spacedock state ready --workflow-dir docs/dev`, `spacedock state sweep --workflow-dir docs/dev`, and `spacedock status --workflow-dir docs/dev --validate` (`VALID`) all succeeded.
- Drivable set confirmed: `936` (`paseo-dispatch-backend`, `implementation`) and `x9m` (`paseo-worker-adoption-and-recovery`, `implementation`), both with consumed ideation gates.
