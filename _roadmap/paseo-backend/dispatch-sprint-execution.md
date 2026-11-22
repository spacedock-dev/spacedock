# `paseo-backend` Commander Package — Cold-Boot Execution Dispatch

You are the **Commander**: a fresh, cold-booted First Officer driving sprint `paseo-backend` (`0.29.x`) from approved ideations to merged deliverables on `main` and a pre-cut antipattern audit. Shaping is complete, the independent preflight staff review ([staff-review.md](staff-review.md)) is folded into both tasks, and both ideation gates are closed-approved and consumed (`status: implementation`).

**Captain's operating directive (binding):**
> Assume this sprint's behavior (`paseo-backend`) already works, operate as its design (dispatching your workers through Paseo on this daemon), and identify/record every friction you encounter.

Keep a running friction log at `docs/dev/.spacedock-state/_roadmap/paseo-backend/commander-frictions.md` (path-scoped commits in `.spacedock-state`) capturing every real snag you hit — CLI flags, stdout/stderr quirks, Safehouse denials, skill discovery, SSH push behavior, daemon/worktree lifecycle — and note whether each was fixed in-sprint or deferred for follow-up.

---

## Boot order

1. **Verify your Paseo + SSH environment first (fail-fast):**
   ```bash
   export GIT_SSH_COMMAND='ssh -o BatchMode=yes'
   which -a paseo
   paseo --version && paseo status --json
   ```
   - Confirm `connectedDaemon == "reachable"` and a working `paseo` binary inside Safehouse (this Commander has `--label paseo-safehouse.paseo-access=daemon`, which grants read-only access to `/Applications/Paseo.app` so `/Users/clkao/.local/bin/paseo` works; if `PATH` puts an unbuilt worktree `node_modules/.bin/paseo` ahead of it, pin `PASEO_BIN=/Users/clkao/.local/bin/paseo` and log the friction).
   - Confirm `spacedock --version` (`0.29.0-pre0`) and stay at the repo root (`/Users/clkao/git/spacedock-research/spacedock-v1`).
2. **Load the FO contract** from this repo's `skills/first-officer/SKILL.md` and `references/first-officer-shared-core.md` (loading `fo-dispatch-core.md`, `fo-write-core.md`, `fo-merge-core.md`, `fo-gate-lifecycle`, and `present-gate` at their triggers) and engage `docs/dev` (`spacedock status --boot --identify --json`).
   *(Note: global `~/.gemini/config/plugins/spacedock/` on this host may be older than `0.29.0-pre0` — log any drift you observe.)*
3. **Read in full before dispatching:**
   - `docs/dev/.spacedock-state/_roadmap/paseo-backend/index.md`
   - `docs/dev/.spacedock-state/_roadmap/paseo-backend/staff-review.md` (`8447132f3`) — all 8 Material findings (`M1`–`M8`) are already folded into `936` and `x9m`; pay close attention to `Non-Blocking` (#1–#7) and `Commander Cold-Boot Readiness Notes`.
4. **Query the drivable set:**
   ```bash
   spacedock status --workflow-dir docs/dev --where sprint=paseo-backend --where 'sprint-readiness != defer'
   ```
   Expected members (both already `status: implementation` with consumed ideation gates — do **not** re-ask ideation gates):
   - **`936`** (`paseo-dispatch-backend`)
   - **`x9m`** (`paseo-worker-adoption-and-recovery`)

---

## Operating as `paseo-backend` (how to dispatch workers on this daemon)

You are running inside Safehouse on `127.0.0.1:6769` with `paseo-safehouse.paseo-access=daemon` (while workspace daemon access stays Off so child ensigns stay locked down without daemon access).

1. **Stamp & build:**
   - For **`936`** (before `--host paseo` exists): use `spacedock dispatch build --workflow-dir docs/dev --entity-path <abs-entity> --stage <stage> --host codex --stamp --checklist-file -` and dispatch with a host-neutral prompt pointer (log the friction; switch to `--host paseo` as soon as `936` ships!).
   - For **`x9m`** (after `936` merges to `main` and `./spacedock` is rebuilt): use `--host paseo --stamp` directly and forward its emitted `prompt` verbatim.
2. **Workspace find-or-create (M3):**
   - Resolve physical `code_cwd` from `spacedock status --workflow-dir docs/dev --read <slug> --json` (`frontmatter.worktree` resolved against repo root + `realpath`).
   - Run `paseo workspace ls --json` and reuse any active workspace whose canonical `cwd == code_cwd`; call `paseo workspace create --isolation local --path "$code_cwd" [--project "$project_id"] --title "$description" --json` only when absent. Never use `--isolation worktree`.
3. **Spawn fresh worker (M4, M7):**
   - Run `paseo run --no-wait --workspace "$workspace_id" --provider <provider/model> [--thinking <id>] --title "$description" --label "role=ensign" --label "spacedock.worker=$name" --label "spacedock.epoch=${name}-c${cycle}-a${attempt}" --json "$prompt"`.
   - Preferred local workers on this daemon: `antigravity/gemini-4-neon` (`--thinking high`) or `safehouse-codex` (never pass `paseo-safehouse.paseo-access=daemon` to ordinary ensigns).
   - **Parse stdout only** (`paseo run --workspace ... --json` prints `Using workspace <id>` to stderr even on exit `0` — never pipe `2>&1 | jq`).
4. **Wait & verify:**
   - Poll with `paseo wait "$agent_id" --timeout 300 --json` (exits `0` for `idle|timeout|permission|error`; branch on JSON `status`, treat `timeout` as normal progress with 0 retries).
   - Verify durable Git commits on the code worktree branch and path-scoped stage report in `docs/dev/.spacedock-state` before advancing.
5. **Archive verified workers (M5):**
   - Run `paseo archive "$agent_id" --json` after verification (or `--force` only on supersede).
   - **Never** archive `$PASEO_AGENT_ID`, and **never** archive an FO while any child worker runs (Paseo cascade-archives same-workspace children and detaches cross-workspace children).

---

## Strict serial landing order

**Do NOT dispatch `936` and `x9m` implementation concurrently.** Both touch `skills/first-officer/references/paseo-fo-dispatch.md`, `docs/runtime-live-ci.md`, and share `internal/ensigncycle/paseo_live_helpers_test.go`.

### Wave 1 — `936` (`paseo-dispatch-backend`)
- **Surface:** +850 net LOC (+863/-13) across 17 files (tolerance ±30%, 13–20 files).
- **Day-1 priority check:** Verify on day 1 whether a `paseo-safehouse`-wrapped worker in a non-`/tmp` worktree (`.worktrees/spacedock-ensign-paseo-dispatch-backend`) can write and commit to the sibling `docs/dev/.spacedock-state` checkout without a Seatbelt denial; record the exact outcome in the implementation report and `commander-frictions.md`.
- **Watchouts from `staff-review.md`:**
  - Keep `first-officer-shared-core.md` untouched (23,497 / 23,500-byte cap).
  - Keep `references/paseo-ensign-runtime.md` and `references/ensign-shared-core.md` out of `paseo-fo-dispatch.md` (`TestBootResidentDeferredLoadPointsResolve` resolves `references/*.md` under `skills/first-officer/`).
  - Source `codex.exec` from `internal/release/live_models.txt` everywhere (no hardcoded `gpt-6-luna` literals in Go, CI YAML, or `docs/runtime-live-ci.md`).
  - Stage workflow-free skills via `livePluginDir(t)` (`WriteCodexLocalMarketplace` in `newCodexLiveIsolatedHome` for `codex/*`; fixture `.agents/skills/` for non-Codex overrides) — zero hardcoded `<checkout>/skills/...` paths in `TestLivePaseoDispatchSmoke`.
  - Follow ASD-STE100 in CLI help (`dispatch.go`), errors (`build.go`), and `docs/runtime-live-ci.md`.
- **Validation & merge:**
  - `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, `spacedock status --workflow-dir docs/dev --validate`.
  - Run `TestLivePaseoDispatchSmoke` (if ordinary Safehouse ensigns cannot nest `paseo daemon start` or read `~/.codex/auth.json`, execute it from Commander or record the exact named unavailability per `docs/dev/README.md:305` — never silently skip).
  - **Detached adversarial audit** on a throwaway checkout (`spike-paseo-dispatch-backend`) before merge (`internal/dispatch/build.go`, `skills/**`, `.github/workflows/**` are high-stakes surfaces).
  - **Required CI lanes:** `claude-live`, `codex-live`, `pi-live` green on PR before merge (`docs/dev/README.md:85`).
  - Merge PR to `main`, run `spacedock merge guard paseo-dispatch-backend --verdict passed`, clean worktree/branch, and rebuild `./spacedock` so Wave 2 uses `--host paseo`.

### Wave 2 — `x9m` (`paseo-worker-adoption-and-recovery`)
- **Precondition:** `936` is merged to `main` and local `main` is fast-forwarded.
- **Surface:** +278 net LOC (+288/-10) across 5 files (tolerance +195 to +365 net LOC, 4–7 files). **Zero** edits to `936`'s test files (`paseo_live_helpers_test.go` is consumed as-is) and **zero** uncalled production Go simulators.
- **Watchouts from `staff-review.md`:**
  - Lean 10-field H2 `## Paseo Worker` YAML block (4 states: `running | archived | superseded | held`; no intermediate `verified`).
  - Invariant `spacedock.worker={worker_base}` across cycles/retries; `spacedock.epoch={worker_base}-c{cycle}-a{attempt}`.
  - Step 3 lost-create-receipt lookup queries `paseo ls --global` without `--all` first, then filters exclusively via `paseo inspect` (`Cwd`, ISO-8601 `CreatedAt`, boolean `Archived`).
  - `TestLivePaseoAdoptionSmoke` stops FO-1 via `paseo stop "$fo1AgentId" --json` (never `paseo archive --force`), archiving FO-1 only in `t.Cleanup`.
- **Validation & merge:**
  - `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, `spacedock status --workflow-dir docs/dev --validate`.
  - Run `TestLivePaseoAdoptionSmoke` (same sandbox note as Wave 1).
  - **Detached adversarial audit** on throwaway checkout (`spike-paseo-worker-adoption-and-recovery`) before merge (`fo-write-core.md`, `paseo-fo-dispatch.md` are high-stakes surfaces).
  - **Required CI lanes:** `claude-live`, `codex-live`, `pi-live` green before merge.
  - Merge PR to `main`, run `spacedock merge guard paseo-worker-adoption-and-recovery --verdict passed`, clean worktree/branch.

---

## Escalation rules

Present validation gates and PR merges to the captain (never self-approve gates without explicit conn). Also escalate immediately on:
- a 3rd feedback cycle on either task;
- a surface/tolerance breach beyond ideation estimates;
- an irrecoverable Safehouse/Paseo/CI block; or
- a genuine scope fork.

---

## Close-out

1. **Pre-cut antipattern audit:** With both `936` and `x9m` merged to `main` (and no tag cut), dispatch ONE fresh independent staff-eng reviewer (not Commander, not the implementers) over the assembled sprint to catch cross-cutting antipatterns or integration holes before release. Fix ship-blockers; record non-blockers for the next sprint.
2. **Final repo gates:** `go test ./...` and `go test ./... -race` green on `main`, `gofmt -w ./cmd ./internal` clean, `spacedock status --workflow-dir docs/dev --validate` `VALID`.
3. **Friction ledger:** Commit and push `docs/dev/.spacedock-state/_roadmap/paseo-backend/commander-frictions.md` with every friction encountered while dogfooding Paseo dispatch on `6769`.
4. **Release cut:** Do **not** cut a `v0.29.x` tag without explicit captain authorization (follow `docs/releasing.md` once authorized). Seeding the next sprint belongs to the Shaping FO.
