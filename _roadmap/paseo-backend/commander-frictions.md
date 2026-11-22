# Sprint `paseo-backend` — Commander Friction Log

Running log of real frictions encountered while cold-booting and operating sprint `paseo-backend` as its own design (`paseo` CLI dispatch on the Safehouse-wrapped `127.0.0.1:6769` daemon).

## Summary

| # | Phase | Friction | Severity | Status |
|---|---|---|---|---|
| F1 | Cold boot (Step 1) | `PATH` resolves `paseo` to unbuilt worktree `node_modules/.bin/paseo` (`ERR_MODULE_NOT_FOUND: packages/cli/dist/index.js`) | High | Resolved (`npm run build:server` in Paseo integration worktree) |
| F2 | Cold boot (Step 1) | `/Users/clkao/.local/bin/paseo` → `/Applications/Paseo.app/Contents/Resources/bin/paseo` fails with `Operation not permitted` under Safehouse (`paseo-safehouse`'s `paseoAppBundle()` skipped `/Applications/Paseo.app` by design because `PATH` resolved `paseo` to the dev worktree first) | **Blocker** | Resolved via F1 (`paseo 0.17.1` reachable on `127.0.0.1:6769`) |
| F3 | Cold boot (Step 1) | `PASEO_AGENT_ID` is stripped by `paseo-safehouse` from Commander's shell environment (`PASEO_HOME` and `APP_SANDBOX_CONTAINER_ID` are set, not `PASEO_AGENT_ID`) | High | Filed upstream (`paseo-safehouse` task `tv` / `agent-identity-env`); worked around in-session by explicitly exporting `PASEO_AGENT_ID=3b905937-e237-4060-af8c-87b279dd96d9` |
| F4 | Cold boot (Step 2) | `/opt/homebrew/bin/spacedock --version` (`0.29.0-pre0`) omits `contract 3` on Antigravity because `runtimehost.Detect` returns empty (`internal/cli/cli.go:916`), tripping `first-officer-shared-core.md` binary gate | High | Deferred (post-sprint follow-up; keep `936` faithful to locked ideation gate) |
| F5 | Cold boot (Step 2) | Global `~/.gemini/config/plugins/spacedock/` (Jun 5) drifts from repo `skills/first-officer/` (`SKILL.md` and `first-officer-shared-core.md` differ) | Medium | Deferred (tracked outside sprint; using repo-local skills) |
| F6 | Wave 1 dispatch (`936`) | Before `--host paseo` ships in `936`, Commander must borrow `spacedock dispatch build --host codex --stamp` and wrap it with a host-neutral prompt pointer for the Paseo worker | Medium | Fixed in-sprint by `936` (Wave 2 `x9m` will use `--host paseo --stamp` directly) |

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
- **Disposition:** Resolved. Building `@getpaseo/cli` and workspace deps via `npm run build:server` in `/Users/clkao/git/spacedock-research/paseo/.worktrees/integration` unblocked plain `paseo --version` (`0.17.1`) and `paseo status --json` (`connectedDaemon: "reachable"` on `127.0.0.1:6769`).
- **Root cause note (F2):** `paseo-safehouse`'s `paseoAppBundle()` skipped `/Applications/Paseo.app` by design because `PATH` resolved `paseo` to the dev worktree first (already readable under `add-dirs-ro=~/git/spacedock-research`).

### F3 — `PASEO_AGENT_ID` stripped by `paseo-safehouse`
- **Observed:** `env` inside Commander's shell contains `PASEO_HOME=/Users/clkao/git/spacedock-research/paseo/.worktrees/integration/.dev/paseo-home` and `APP_SANDBOX_CONTAINER_ID=agent-safehouse`, but `PASEO_AGENT_ID` is completely absent.
- **Impact:**
  1. `skills/first-officer/SKILL.md`'s planned `PASEO_AGENT_ID` auto-selection rule will not trigger automatically on Safehouse-wrapped Antigravity without the upstream fix.
  2. `spacedock dispatch build --host paseo`'s `PASEO_AGENT_ID`-based `/tmp/spacedock-dispatch/` filename disambiguation falls back to plain `{derivedName}.md` unless `PASEO_AGENT_ID` is exported.
  3. `paseo run` child spawns will not auto-populate `ParentAgentId` via Paseo CLI's `resolveCallerAgentId` unless `PASEO_AGENT_ID` is present.
- **Disposition:** Confirmed real `paseo-safehouse` bug; filed upstream as task `tv` (`agent-identity-env`). Worked around in this Commander session by explicitly exporting `PASEO_AGENT_ID=3b905937-e237-4060-af8c-87b279dd96d9` on `paseo run` and `spacedock dispatch build --host paseo` commands.

### F6 — Wave 1 bootstrap borrows `--host codex` before `--host paseo` ships
- **Observed:** Dispatching Wave 1 (`936`) before `--host paseo` exists requires `spacedock dispatch build --host codex --stamp` and overriding its Codex-specific `$spacedock:ensign` prompt with a host-neutral prompt pointer when spawning the Paseo worker.
- **Disposition:** Fixed in-sprint by `936` itself; once `936` merges to `main` and `./spacedock` is rebuilt, Wave 2 (`x9m`) switches directly to `--host paseo --stamp` with verbatim prompt forwarding.

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
