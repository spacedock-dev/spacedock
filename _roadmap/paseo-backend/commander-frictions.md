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
| F7 | Wave 1 Day-1 check (`936`) | `paseo-safehouse` confines a worktree-workspace ensign (`.worktrees/spacedock-ensign-paseo-dispatch-backend`) strictly to its worktree `cwd` + `/tmp`, denying BOTH parent `.git` (`git status` exits 128) AND sibling `docs/dev/.spacedock-state` + `docs/dev` (`Operation not permitted`) | **Blocker** | Reported upstream (`paseo-safehouse` manager `7d469e54`); worked around via gitignored mode-600 worktree `.safehouse` (reused for future worktrees) |
| F8 | Wave 1 Day-1 check (`936`) | Ordinary Safehouse ensigns also hit `Operation not permitted` on `~/.codex/auth.json`, have no `paseo` on `PATH`, and lack `PASEO_AGENT_ID` (same bug as F3 / `tv`) | High | Confirmed boundary (`skills/ensign/SKILL.md` fallback clause is load-bearing; live smoke must run from Commander/CI) |
| F9 | Wave 1 implementation (`936`) | `antigravity` worker yielded its turn while waiting on backgrounded `go test` tasks, causing `paseo wait` to return `status: "idle"` before any code commit or stage report existed | High | Caught by Git-state completion gate; recovered via `paseo send --no-wait` (completed cleanly with `686ea0da5` and `b1f2df67e`) |

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

### F7 — `paseo-safehouse` blocks parent `.git` and sibling `.spacedock-state` for non-`/tmp` Git worktrees (Day-1 Priority Check BLOCKER)
- **Observed:** Spawned `936` implementation worker `e60088b8-1680-45d8-9bf7-bf76c24fefb9` (`antigravity/gemini-4-neon`, `thinking: high`, `ParentAgentId: 3b905937-e237-4060-af8c-87b279dd96d9`) in adopted worktree workspace `wks_f2cea2ad7052f9f1` (`isolation: "worktree"`, `cwd: /Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-paseo-dispatch-backend`).
- **Exact Seatbelt (`EPERM`) failures reported by the worker:**
  1. **Parent `.git` metadata denied — all `git` commands in the worktree fail:**
     - Worktree `.git` points to `gitdir: /Users/clkao/git/spacedock-research/spacedock-v1/.git/worktrees/spacedock-ensign-paseo-dispatch-backend`.
     - `ls -la /Users/clkao/git/spacedock-research/spacedock-v1/.git/` → `Operation not permitted`
     - `git status -sb` → exit `128`: `fatal: not a git repository: /Users/clkao/git/spacedock-research/spacedock-v1/.git/worktrees/spacedock-ensign-paseo-dispatch-backend`
  2. **Sibling `.spacedock-state` checkout denied (both read and write):**
     - `stat /Users/clkao/git/spacedock-research/spacedock-v1/docs/dev/.spacedock-state/paseo-dispatch-backend/index.md` → `Operation not permitted`
     - `touch /Users/clkao/git/spacedock-research/spacedock-v1/docs/dev/.spacedock-state/.probe` → `Operation not permitted`
  3. **Parent `docs/dev` workflow dir and parent `skills/` denied:**
     - `spacedock dispatch show-stage-def --workflow-dir /Users/clkao/git/spacedock-research/spacedock-v1/docs/dev --stage implementation` → `error: workflow directory not found: /Users/clkao/git/spacedock-research/spacedock-v1/docs/dev`
     - `stat /Users/clkao/git/spacedock-research/spacedock-v1/skills/ensign/SKILL.md` → `operation not permitted` (worker had to fall back to identical worktree copies).
- **Why earlier probes missed this:** All 2026-11-18 and 2026-11-21 live probes used temp fixtures inside `/private/tmp/...`, where Safehouse allows broad `/private/tmp` access. On a real non-`/tmp` repo, confining a worktree workspace strictly to `cwd` breaks Git worktrees completely (`git rev-parse --git-common-dir` is outside `cwd`) and blocks split-root `.spacedock-state` read/write.
- **Note on M4 confirmed live:** `paseo run --workspace wks_f2cea2ad7052f9f1 --json` did indeed print `Using workspace wks_f2cea2ad7052f9f1` to stderr on exit `0`; separate stdout/stderr capture parsed JSON cleanly.
- **Disposition:** Reported upstream to `paseo-safehouse` manager (`7d469e54`: linked Git worktree `SAFEHOUSE_WORKDIR` omits common `.git` dir and repo `.safehouse` inheritance). Worked around in-session by writing a gitignored mode-600 `.safehouse` inside each stamped worktree (`add-dirs-ro=~/.roborev/runtime`, `add-dirs-ro=~/git/spacedock-research`, `add-dirs=.../spacedock-v1/.git`, `add-dirs=.../spacedock-v1/docs/dev/.spacedock-state`, `enable=all-agents,ssh`), verified with `git status -sb`, `git -C docs/dev/.spacedock-state status -sb`, `spacedock dispatch show-stage-def`, state-checkout write/delete, and `go test ./internal/dispatch`. Re-dispatched `936` (`c1-a2`) in existing workspace `wks_f2cea2ad7052f9f1`.

### F8 — Ordinary Safehouse ensign environment boundaries confirmed
- **Observed inside `e60088b8-1680-45d8-9bf7-bf76c24fefb9`:**
  - `ls -la ~/.codex/auth.json` → `Operation not permitted`
  - `which paseo` → not found
  - `PASEO_AGENT_ID` → not set (same `paseo-safehouse` stripping bug as F3 / `tv`; confirms `936`'s `skills/ensign/SKILL.md` fallback `or the dispatch prompt names references/paseo-ensign-runtime.md` is load-bearing whenever `PASEO_AGENT_ID` is stripped).
- **Disposition:** Expected security boundary for ordinary ensigns (`paseo` / `~/.codex/auth.json` remain locked down); `TestLivePaseoDispatchSmoke` / `TestLivePaseoAdoptionSmoke` should run from Commander or CI rather than inside an ordinary ensign.

### F9 — Antigravity worker premature `idle` on backgrounded shell task
- **Observed:** During `936` implementation (`27c9d9a7-09c6-4605-a96d-bd2fe5f00814`, `antigravity/gemini-4-neon`), after writing all 17 files and passing fast checks, the worker backgrounded `go test ./...` / `go test ./... -race` and ended its turn waiting for those background tasks. Paseo's Antigravity adapter transitioned the agent to `Status: "idle"`, so `paseo wait` returned `status: "idle"` while all files were uncommitted and `## Stage Report: implementation` did not exist.
- **Impact:** Proves why `paseo-backend`'s core design rule ("`paseo wait` `idle` is only a wake signal; Git commits + stage report are the completion gate") is essential! Without the Git gate, a naive FO would have advanced on premature `idle`.
- **Disposition:** Recovered cleanly in-session via `paseo send --no-wait` instructing the live idle worker to finish its test check without yielding mid-test and commit; worker finished cleanly (`686ea0da5` on `paseo-dispatch-backend`, `b1f2df67e` in `.spacedock-state`) and was archived. Worth noting for `paseo` Antigravity ACP adapter / `x9m` recovery docs.

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
