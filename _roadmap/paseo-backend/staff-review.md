# Sprint `paseo-backend` — Independent Preflight Staff Review

**Reviewer:** Independent Staff Engineer (preflight sprint-wide review, refutation + day-1 implementability lens)  
**Date:** 2026-11-22  
**Scope reviewed:**
- `docs/dev/.spacedock-state/_roadmap/paseo-backend/index.md` (sprint goal, DoD, sequence)
- `docs/dev/.spacedock-state/paseo-dispatch-backend/index.md` (`936`)
- `docs/dev/.spacedock-state/paseo-worker-adoption-and-recovery/index.md` (`x9m`)
- Verified directly against `spacedock-v1` (`internal/dispatch/`, `internal/status/`, `internal/contractlint/`, `internal/ensigncycle/`, `internal/release/`, `skills/`, `.github/workflows/runtime-live-e2e.yml`, `docs/runtime-live-ci.md`, `docs/runtime-support.md`) and Paseo 0.17.1 source (`/Users/clkao/git/spacedock-research/paseo/.worktrees/integration/packages/{cli,server}/`).

## Executive Verdict

**`NEEDS REVISION BEFORE LOCKING GATES`**

The sprint's core direction — one host-neutral `--host paseo` dispatch envelope, a CLI-first FO binding wired via `PASEO_AGENT_ID`, real isolated-daemon live smoke tests with zero fake ACP mocks, and serial landing (`936` -> `x9m`) — is sound and directly serves the sprint DoD.

However, **eight Material findings** (six concrete day-1 bugs/proof holes verified against real Paseo 0.17.1 and Spacedock source, one necessity-first over-engineering block in `x9m`, and one cross-member seam/CI-SSOT gap) will break day-1 implementation, produce false greens, or violate repo policy if the ideation gates lock as written. All eight have small, concrete fixes that should be folded into `936` and `x9m` before locking gates.

---

## What Held Under Refutation

Direct source verification confirmed these planned mechanisms are sound:

- **`spacedock dispatch name` exists and matches `x9m`'s usage:** `internal/dispatch/dispatch.go:34-35` routes `name` to `runName` (`internal/dispatch/names.go:102-147`), a read-only query accepting `--workflow-dir`, `--entity-path`, and `--stage`, printing the canonical `worker_base` to stdout.
- **`@getpaseo/cli` npm package name is real:** Paseo `packages/cli/package.json:2` publishes `"name": "@getpaseo/cli"` (`README.md:112`, `public-docs/index.md:24`).
- **CI environment and Codex model match repo conventions:** `CI-E2E-CODEX` and `CODEX_AUTH_JSON` exist in `.github/workflows/runtime-live-e2e.yml:405-409`, and `gpt-6-luna` is the repo's registered `codex.exec` model in `internal/release/live_models.txt:3` (though it must be read from `live_models.txt`, not hardcoded — see M8).
- **Isolated Paseo daemon config keys are exact:** All five `paseo daemon config set` keys (`daemon.listen`, `daemon.relay.enabled`, `features.webUi.enabled`, `features.dictation.enabled`, `features.voiceMode.enabled`) match Paseo 0.17.1's own E2E harness (`packages/cli/src/commands/daemon/lifecycle.e2e.test.ts:90-94`) and avoid removed `--port`/`--no-relay` flags (`local-daemon.ts:83-91`).
- **`paseo wait --json` exit codes and `inspect` PascalCase schema are exact:** `packages/cli/src/commands/agent/wait.ts:102-144` exits `0` for all four statuses (`idle|timeout|permission|error`), and `inspect.ts:14-49, 129-155` emits the exact PascalCase fields `936` documents.
- **H2 `## Paseo Worker` is safe for `status` parsers:** `internal/status/section_read.go:153-164` and `internal/status/gate_extract.go:127, 143-153` terminate `## Stage Report:` at any `## ` heading, so an H2 worker record cannot contaminate stage-report checklist or `--ac-scan` spans.
- **FO byte-cap preservation in `936`:** Wiring Paseo in `skills/first-officer/SKILL.md` (with `{first_officer_base}/references/paseo-fo-dispatch.md`, not `@references/`) leaves `first-officer-shared-core.md` untouched at 23,497 / 23,500 bytes (`internal/contractlint/fo_function_reference_invariant_test.go:46-54`) and satisfies `TestFirstOfficerEntryEagerlyImportsOnlySharedCore` (`first_officer_eager_references_test.go:12-30`).
- **Excluding `PASEO_AGENT_ID` from `internal/runtimehost.Detect`:** Correct; adding it would cause false host-ambiguity errors for Claude/Codex/Pi sessions hosted under Paseo.

---

## Material Findings

### M1. Necessity (`x9m`): Uncalled production Go simulator (`paseo_worker.go`), 22-key YAML schema, and prose-to-simulator `contractlint` binding fail necessity-first and Proof Policy

- **Owner:** `x9m` (`docs/dev/.spacedock-state/paseo-worker-adoption-and-recovery/index.md`)
- **Code & policy proof:**
  - `x9m` (`index.md:48, 252, 264-266`) explicitly introduces **zero CLI subcommands**, yet adds `+110 LOC` of production Go in `internal/dispatch/paseo_worker.go` (`PaseoWorkerRecord`, `PaseoLabels`, `ParseDispatchRetries`, `ReconcilePaseoWorker`), `+120 LOC` of unit tests (`paseo_worker_test.go`), and `+65 LOC` of `contractlint` dual-extraction (`paseo_adoption_binding_test.go`) — ~295 LOC of production/test machinery that **no shipped binary path ever calls**.
  - `docs/dev/README.md:88` ("Trace every mechanism to value") forbids building a second implementation of the system under test; `docs/dev/README.md:163` requires challenging necessity before hardening coherence; `internal/contractlint/boundary_guard_test.go:35-38` explicitly bans **`CODE-BOUND-AS-BEHAVIOR-SUBSTITUTE`** ("asserting an instruction file's prose matches a code value... That is a consistency lint, not a behavior test"). Unlike `internal/contractlint/reconcile_class_binding_test.go:187-196` (where `spacedock dispatch reconcile` is a real executed CLI verb emitting JSON to the FO), the Paseo FO never executes `paseo_worker.go` at runtime — it executes `paseo` CLI directly from `paseo-fo-dispatch.md`. Pinning `paseo-fo-dispatch.md` to an uncalled Go simulator proves only that two files agree, not that a Paseo FO follows the recovery matrix.
  - The 22-required-key fixed-order `## Paseo Worker` schema (`x9m:59-88`) is also over-specified for an LLM-authored entity section: `host: paseo` duplicates the heading; `worktree` duplicates `frontmatter.worktree`; `requested_model` duplicates `requested_provider` (`codex/gpt-6-luna`); `cycle`/`attempt`/`worker_label`/`epoch_label` duplicate deterministic labels; and persisting an intermediate `state: verified` before `paseo archive` (<1s later) adds unnecessary state-repo commit churn.
- **Concrete simplest fix:**
  1. Reduce `## Paseo Worker` to a lean ~10-field YAML record sufficient for the sprint DoD (`index.md:34`): `agent_id`, `workspace_id`, `server_id`, `stage`, `epoch`, `cwd`, `pre_dispatch_state_rev`, `requested` (`provider`, `mode`, `thinking`), `returned` (`provider`, `mode`, `thinking`), `state` (`running | archived | superseded | held`), and optional `hold_reason`. Write it on launch/adoption and terminal outcome (`archived`/`superseded`/`held`), not on intermediate `verified`.
  2. Delete uncalled production `internal/dispatch/paseo_worker.go`, `internal/dispatch/paseo_worker_test.go`, and `internal/contractlint/paseo_adoption_binding_test.go` from `x9m`. Keep lightweight YAML parsing test-local inside `internal/ensigncycle/paseo_adoption_live_test.go`, and keep `TestPaseoWorkerStatusCompatibility` (AC-4) directly over fixture markdown. (If the captain later wants a real `spacedock` CLI helper to own reconciliation, commission it as a real CLI verb rather than shipping dead production Go.)

---

### M2. Day-1 Bug (`x9m`): Step 3 lost-create-receipt reconciliation filters on three `paseo ls --json` fields that do not behave as assumed

- **Owner:** `x9m` (`docs/dev/.spacedock-state/paseo-worker-adoption-and-recovery/index.md`)
- **Code proof:** `x9m` (`index.md:143`) specifies:
  > `paseo ls --all --global --label "spacedock.worker={expected_worker}" --json`, filtered to rows with `realpath(row.cwd) == expected_cwd` and (when `started` is set) `row.created >= started`, preferring non-archived rows (`status != "archived"`).
  
  Against real Paseo 0.17.1 (`packages/cli/src/commands/agent/ls.ts`):
  1. **`row.cwd` is `~`-shortened (`ls.ts:56-62, 114`):** `cwd: shortenPath(agent.cwd)` turns `/Users/<user>/...` into `~/...`. POSIX `realpath("~/...")` fails because `realpath` does not expand `~`.
  2. **`row.created` is human relative time (`ls.ts:44-53, 115`):** `created: relativeTime(agent.createdAt)` emits `"just now"`, `"5 minutes ago"`, `"2 hours ago"`, or `"1 days ago"` — never an ISO-8601 timestamp comparable to `frontmatter.started`.
  3. **`row.status` is NEVER `"archived"` (`ls.ts:29-41, 113, 229-231`):** `status` is the runtime lifecycle (`"running"|"idle"|"closed"|"error"`), and `AgentListItem` omits `archived` and `archivedAt` completely. With `--all`, archived and unarchived closed/idle agents are indistinguishable from `paseo ls --json` alone.
- **Concrete simplest fix:**
  - In Step 3 fallback lookup, run `paseo ls --global --label "spacedock.worker={worker_base}" --json` **without `--all`** first (which excludes archived agents server/CLI-side per `ls.ts:229-231`), falling back to `--all` only if empty.
  - Run `paseo inspect <candidate.id> --json` (`packages/cli/src/commands/agent/inspect.ts:137-141`) on each candidate and filter on `inspect.Cwd` (unshortened physical path), `inspect.CreatedAt` (ISO-8601), and `inspect.Archived` (boolean). Remove all references to `ls[].created >= started`, `realpath(ls[].cwd)`, and `ls[].status != "archived"`.

---

### M3. Day-1 Bug (`936` & `x9m`): `paseo workspace create` is NOT idempotent — repeated calls create duplicate workspaces for the same directory (dogfood friction #6)

- **Owner:** `936` (`docs/dev/.spacedock-state/paseo-dispatch-backend/index.md` — `paseo-fo-dispatch.md` Verb 1 and `TestLivePaseoDispatchSmoke`); consumed by `x9m`.
- **Code proof:**
  - `936` (`index.md:163`) specifies calling `paseo workspace create --isolation local --path "$code_cwd" ... --json` unconditionally on every spawn (dropping `staff-report.md:20`'s *"or use an already identified workspace for the same project and physical path"*).
  - In Paseo 0.17.1, `packages/cli/src/commands/workspace/create.ts:168` -> `packages/server/src/server/session.ts:7558` (`handleWorkspaceCreateLocal`) calls `workspaceProvisioning.createWorkspaceForDirectory(...)` — **NOT** `findOrCreateWorkspaceForDirectory` (`workspace-provisioning-service.ts:395`).
  - `workspace-provisioning-service.ts:285-301` unconditionally mints `generateWorkspaceId()` and upserts a brand-new workspace record on every call, even for an identical directory.
  - Every non-worktree stage (e.g., `ideation` in the repo root, which already has the FO's workspace), every sticky-worktree transition (`implementation` -> `validation`), every feedback rework cycle, and every `-retry` dispatch (`x9m` row 5) will leak duplicate workspaces for the same checkout.
- **Concrete simplest fix:**
  - Make Verb 1 in `paseo-fo-dispatch.md` an explicit **find-or-create**: first call `paseo workspace ls --json` (`packages/cli/src/commands/workspace/ls.ts:18-48`, `shared.ts:33-41` emits unshortened `cwd`, `workspaceId`, `project`, `background`), match an active workspace whose canonical `cwd` equals `$code_cwd`, and call `paseo workspace create --isolation local --path "$code_cwd" ... --json` only when none exists.

---

### M4. Day-1 Bug (`936` & `x9m`): `paseo run --workspace ... --json` writes unconditional non-JSON `Using workspace wks_...` to stderr on exit 0 (dogfood friction #1)

- **Owner:** `936` (`docs/dev/.spacedock-state/paseo-dispatch-backend/index.md` — `paseo-fo-dispatch.md` CLI invariants + shared Paseo CLI runner in `internal/ensigncycle/`).
- **Code proof:**
  - `packages/cli/src/commands/agent/run.ts:553, 563, 584-588` calls `console.error(\`Using workspace ${explicit}\`)` unconditionally whenever `--workspace` is passed — even when `--json` is set! (`with-output.ts:92` writes JSON to stdout, while `run.ts:553` writes non-JSON prose to stderr on exit 0.)
  - `936` (`index.md:144`) states only *"Successful commands exit `0` with JSON on stdout; CLI/daemon errors exit `1` with `{"error":...}` on stderr"*, omitting exit-0 stderr noise.
  - Any shell caller combining streams (`2>&1`), any host Bash tool merging stdout+stderr, or any Go test helper using `cmd.CombinedOutput()` (the pattern already used in `internal/ensigncycle/codex_live_runner_test.go:230` and `dispatch_recovery_assert_test.go:266`) receives `Using workspace wks_...\n{\n  "agentId": ...` and fails JSON decoding on day 1.
- **Concrete simplest fix:**
  - Document explicitly in `paseo-fo-dispatch.md` that `paseo run --json` writes `Using workspace <id>` to stderr even on exit 0, and that callers MUST parse stdout only (`2>/dev/null` when piping to `jq`, or capture stderr separately for failure diagnostics).
  - Mandate separate `cmd.Stdout` and `cmd.Stderr` buffers (never `cmd.CombinedOutput()`) in `paseo_live_smoke_test.go`'s shared `runPaseoJSON` helper used by both `936` and `x9m`.

---

### M5. Day-1 Bug (`x9m` & `936`): `paseo archive "$fo1AgentId" --force` cascade-archives same-workspace children and strips `ParentAgentId` on cross-workspace children

- **Owner:** `x9m` (`TestLivePaseoAdoptionSmoke` Phase 1) and `936` (`paseo-fo-dispatch.md` archival guard).
- **Code proof:**
  - In Paseo 0.17.1 `packages/server/src/server/agent/agent-manager.ts:2319, 2328-2368`: `archiveAgentUnlocked` calls `cascadeArchiveChildren(archivedRecord)`.
  - `shouldDetachFromArchivedParent` (`agent-manager.ts:818-827`) returns `true` **only** when `parent.workspaceId !== child.workspaceId` or `hasOpenAgentTab(child.labels)` is set:
    - **Same workspace (`parent.workspaceId === child.workspaceId`):** `archiveAgentUnlocked(currentChild.id)` **immediately kills and archives the child worker!** That breaks every non-worktree stage (e.g., `ideation` on the repo root once M3 reuses the repo-root workspace).
    - **Different workspace (`parent.workspaceId !== child.workspaceId`):** `detachAgentUnlocked(currentChild.id)` keeps `W1` alive **only by accident of being cross-workspace**, and **deletes `paseo.parent-agent-id`** from `W1`'s labels (`agent-manager.ts:829-837`).
  - `x9m` (`index.md:179`) stops FO-1 mid-stage via `paseo archive "$fo1AgentId" --force` — contradicting the sprint's own invariant ("never archive an FO while any child worker is live", `index.md:25`, `936:167`, `x9m:157`), silently detaching `W1`, and failing outright if ever adapted to a non-worktree stage.
- **Concrete simplest fix:**
  - In `TestLivePaseoAdoptionSmoke` Phase 1, stop FO-1 with `paseo stop "$fo1AgentId" --json` (`packages/cli/src/commands/agent/stop.ts:89` calls `client.cancelAgent(agent.id)`, interrupting FO-1's run without archiving FO-1 or cascading/detaching `W1`), verify `W1` remains `running` with `ParentAgentId == fo1AgentId`, boot fresh FO-2, and archive FO-1 only in `t.Cleanup` after `W1` finishes.
  - Document Paseo's same-workspace cascade-archive rule in `paseo-fo-dispatch.md` as the mechanical reason FOs must never be archived while children run (and note that non-worktree workers share the repo-root workspace).

---

### M6. Proof Gap (`936`): `TestLivePaseoDispatchSmoke` cheats on installed-skill discovery by hardcoding `<checkout>/skills/...` paths (dogfood friction #7)

- **Owner:** `936` (`docs/dev/.spacedock-state/paseo-dispatch-backend/index.md` — `TestLivePaseoDispatchSmoke` & AC-1).
- **Code proof:**
  - `936` (`index.md:208`) specifies:
    > `$foPrompt` instructs the FO to load `<checkout>/skills/first-officer/references/paseo-fo-dispatch.md`, run `spacedock dispatch build --host paseo --stamp`, adopt the stamped worktree, spawn ONE fresh worker via `paseo run --no-wait` (pointing at current-checkout `<checkout>/skills/ensign/SKILL.md` + `paseo-ensign-runtime.md`)...
  - This bypasses three things AC-1 and the sprint DoD claim to prove:
    1. **`skills/first-officer/SKILL.md` `PASEO_AGENT_ID` selection** is never exercised when the prompt names `paseo-fo-dispatch.md` directly.
    2. **`skills/ensign/SKILL.md` `PASEO_AGENT_ID`-first precedence and consumer-repo skill discovery** are never exercised when the FO overrides `dispatch build --host paseo`'s emitted `prompt` (`936:85`, which does NOT embed `<checkout>/skills/ensign/SKILL.md`) with a raw checkout path. A consumer repo with no `skills/` directory will fail to discover `spacedock:ensign` while this smoke stays green.
    3. **Wrong-root discovery hazard (`internal/ensigncycle/live_test.go:63-114`):** Pointing an agent into the real `spacedock-v1` checkout exposes `docs/dev` — the exact wrong-root leak `livePluginDir(t)` was built to eliminate.
- **Concrete simplest fix:**
  - Stage a workflow-free copy of the current checkout's skills using `livePluginDir(t)` (`internal/ensigncycle/live_test.go:74-114`) and install it into the isolated test environment so named skill discovery actually runs:
    - For `codex/*` providers: install via `cli.WriteCodexLocalMarketplace` + `codex plugin marketplace add` + `codex plugin add spacedock@spacedock` into the isolated `CODEX_HOME` (reusing `internal/ensigncycle/codex_marketplace_test.go:29` and mirroring `.github/workflows/runtime-live-e2e.yml:513-548`).
    - For non-Codex local overrides (e.g., `antigravity/*`): stage the workflow-free `skills/` into the consumer fixture's `.agents/skills/` (or isolated host skill dir) so named skill lookup resolves from the staged current checkout rather than stale global `~/.gemini/config/plugins/spacedock` or raw `<checkout>/skills/...`.
  - Have `$foPrompt` invoke the installed `first-officer` skill by name (never hardcoding `paseo-fo-dispatch.md`), and forbid appending raw `<checkout>/skills/ensign/SKILL.md` paths when spawning the worker so `dispatch build --host paseo`'s real emitted `prompt` is tested verbatim.

---

### M7. Cross-Member Seam & `names.go` Bug (`936` & `x9m`): `$epoch` is undefined in `936`, while `x9m`'s `-cycle{cycle}-retry` label breaks `internal/dispatch/names.go` and weakens `TestLivePaseoAdoptionSmoke`

- **Owner:** Both `936` and `x9m`.
- **Code proof:**
  - `936` (`index.md:164, 173`) passes `--label "spacedock.epoch=$epoch"` in Verb 2 and cites `spacedock.epoch` in `«worker-identity»`, but **never defines `$epoch`'s format** — a day-1 implementability gap.
  - `x9m` (`index.md:102`) sets `spacedock.worker` to `{worker_base}-cycle{cycle}-retry` when `cycle >= 2, attempt == 2`. `"-cycle2-retry"` is **13 characters**, violating `internal/dispatch/names.go:15-24, 37-42` (`semanticName` reserves only 8 chars for suffixes so names fit `nameAgentMaxLen = 64`; `validWorkerSuffix` rejects `len(s) > 8` and compound `-cycleN-retry`) and `skills/fo-dispatch-recovery/SKILL.md:27` (max 8-char suffix).
  - Mutating `spacedock.worker` to `{worker_base}-retry` also creates a false-green trap in `TestLivePaseoAdoptionSmoke` (`x9m:182`): querying `paseo ls --all --global --label "spacedock.worker=<worker_base>"` would **miss** a duplicate `-retry` worker spawned with `spacedock.worker=<worker_base>-retry`.
- **Concrete simplest fix:**
  - Freeze ONE shared label contract across `936` and `x9m`:
    - `spacedock.worker`: always `{worker_base}` (canonical `envelope.name` from `dispatch build --host paseo` / `dispatch name`), invariant across retries and cycles — always <= 56 chars and valid under `names.go`.
    - `spacedock.epoch`: `{worker_base}-c{cycle}-a{attempt}` — defined in `936` as `{worker_base}-c1-a1` for v1 fresh dispatch, extended in `x9m` for `cycle >= 1` and `attempt in {1, 2}`.
  - Then `paseo ls --all --global --label "spacedock.worker={worker_base}"` reliably enumerates **all** attempts (including retries) for `(entity, stage)`, while `spacedock.epoch` uniquely identifies each attempt.

---

### M8. Cross-Member Seam & CI Conventions (`936` & `x9m`): Shared live-test helper contract is unspecified, and `936`'s CI workflow violates `live_models.txt` SSOT and omits standard security/env wiring

- **Owner:** `936` (`docs/dev/.spacedock-state/paseo-dispatch-backend/index.md`); `x9m` consumes the shared helper.
- **Code proof:**
  - **Unspecified shared live-test seam:** `x9m` (`index.md:175, 244`) reuses `936`'s isolated-daemon helpers from `internal/ensigncycle/paseo_live_smoke_test.go`, yet `936` defines no helper signatures or file split, leaving `x9m` to modify `936`'s test file (outside `x9m`'s 7-file surface).
  - **Model-ID SSOT violation:** `936` (`index.md:219, 224, 250, 261, 322`) hardcodes `codex/gpt-6-luna` across `.github/workflows/paseo-live-e2e.yml`, `paseo_live_smoke_test.go`, `paseo_live_lane_test.go`, and `docs/runtime-live-ci.md`. Repo policy (`internal/release/livemodels.go:3`: *"No consumer repeats a model ID"*; `docs/runtime-live-ci.md:135`: *"The lane model ids live in one place, `internal/release/live_models.txt`, so this document never restates one"*) requires reading `codex.exec` from `internal/release/live_models.txt`.
  - **CI workflow gaps in `936` (`index.md:217-228`):**
    1. Missing `permissions: contents: read` (standard in `runtime-live-e2e.yml:60-61`).
    2. Secret check mentions only `CODEX_AUTH_JSON` instead of accepting `CODEX_AUTH_JSON` or `OPENAI_API_KEY` (contrast `runtime-live-e2e.yml:416`).
    3. Missing `SPACEDOCK_BIN` and `SPACEDOCK_REPO_ROOT` exports after `go build -o ./spacedock ./cmd/spacedock` (contrast `runtime-live-e2e.yml:498-500`).
    4. No `offline` pre-gate job before requesting `environment: CI-E2E-CODEX` approval on every `pull_request` — adding a second manual deployment approval gate on every PR to `main` even when offline tests fail.
- **Concrete simplest fix:**
  1. In `936`, specify the small package-private helper seam in `internal/ensigncycle/paseo_live_helpers_test.go` (`startIsolatedPaseoDaemon(t)`, `(*paseoDaemon).runJSON(t, args...)`, `setupPaseoSplitRootFixture(t, ...)`) so `x9m` adds `paseo_adoption_live_test.go` with zero edits to `936`'s files.
  2. Default the provider in Go via `"codex/" + release.LiveModels["codex.exec"]`, read `codex.exec` from `internal/release/live_models.txt` in `.github/workflows/paseo-live-e2e.yml`, and omit hardcoded model literals from `docs/runtime-live-ci.md`.
  3. Add `permissions: contents: read`, `CODEX_AUTH_JSON || OPENAI_API_KEY` check, `SPACEDOCK_BIN`/`SPACEDOCK_REPO_ROOT` exports, and either `needs: offline` (matching `runtime-live-e2e.yml`) or `workflow_dispatch` + Paseo-path-filtered `pull_request` trigger to `.github/workflows/paseo-live-e2e.yml`.

---

## Non-Blocking / Deferred Risks / Polish

1. **`TestBootResidentDeferredLoadPointsResolve` trap in `paseo-fo-dispatch.md` (`936`):** `internal/contractlint/boot_resident_closure_test.go:58-73` extracts every `references/*.md` match from `deferredModuleBodies` and resolves it under `skills/first-officer/`. Once `936` adds `paseo-fo-dispatch.md` to `deferredModuleBodies`, mentioning `references/paseo-ensign-runtime.md` or `references/ensign-shared-core.md` inside `paseo-fo-dispatch.md` will fail `TestBootResidentDeferredLoadPointsResolve` (those files live under `skills/ensign/`, not `skills/first-officer/`). Keep ensign reference paths out of `paseo-fo-dispatch.md`.
2. **Isolated `CODEX_HOME` path on macOS (`936`):** `internal/ensigncycle/codex_liveenv.go:125-144` deliberately places isolated `CODEX_HOME` outside `$TMPDIR` (`newCodexLiveIsolatedHome`) because macOS sandboxing can reject system-temp `CODEX_HOME`. Reuse `newCodexLiveIsolatedHome` in `paseo_live_helpers_test.go` if Codex ACP hits temp-dir restrictions on macOS.
3. **Instruction-file read quarantine on `docs/runtime-live-ci.md` (`936`):** `internal/contractlint/instruction_read_detector_test.go:107-110` treats `docs/runtime-live-ci.md` as a quarantined instruction surface. Keep all workflow/doc structural tests inside `internal/contractlint/` (already planned by `936`) and do not alter the `claude`/`codex`/`pi` common-suite lines bound by `TestRuntimeLiveCommonSuiteTimeouts` (`runtime_live_lane_test.go:76-87`).
4. **ASD-STE100 prose compliance (`936`, `x9m`):** Per `docs/dev/README.md:90-100`, all new CLI help lines (`internal/dispatch/dispatch.go`), error diagnostics (`internal/dispatch/build.go`), and `docs/runtime-live-ci.md` additions must follow ASD-STE100 (imperative procedures <= 20 words, descriptions <= 25 words, modals restricted to `can`/`will`/`must`).
5. **Non-`/tmp` Safehouse sibling `.spacedock-state` write check (`936:273`):** Still unproven outside `/tmp`. Run this check first on day 1 of `936` implementation so any `paseo-safehouse` write-scope restriction surfaces before test polish.
6. **Sprint `index.md` sync after M1:** Update `docs/dev/.spacedock-state/_roadmap/paseo-backend/index.md:34` if M1's reduced `## Paseo Worker` schema drops redundant fields (e.g., `worktree` duplicating frontmatter), so sprint DoD and `x9m` stay aligned.
7. **`internal/cli` `--host` lists stay `claude|codex|pi`:** Confirmed `internal/cli/cli.go:517-567`, `internal/cli/help.go:24-25`, and `internal/cli/pi_frontdoor_test.go:121-122` cover `install`/`doctor` host targets only; `936` correctly limits `--host paseo` to `spacedock dispatch build`.

---

## Commander Cold-Boot Readiness Notes

Before packaging `dispatch-sprint-execution.md` for a cold-booted Commander running inside Safehouse with `paseo-safehouse.paseo-access=daemon` on `127.0.0.1:6769`:

1. **Strict serial execution (`936` merged to `main` before `x9m` starts implementation):** Both tasks edit `skills/first-officer/references/paseo-fo-dispatch.md`, `docs/runtime-live-ci.md`, and share the `internal/ensigncycle/` Paseo live helper. Never dispatch `936` and `x9m` implementation concurrently.
2. **Verify a working `paseo` binary before Commander boot (live observation during this review):** Inside this Safehouse environment, `which paseo` resolved to `/Users/clkao/git/spacedock-research/paseo/.worktrees/integration/node_modules/.bin/paseo`, which crashed immediately with `ERR_MODULE_NOT_FOUND` (`packages/cli/dist/index.js` unbuilt), while `/Applications/Paseo.app/Contents/Resources/bin/paseo` returned `Operation not permitted` under Safehouse. The Commander boot recipe MUST build `packages/cli` (`npm run build:server` in the Paseo worktree) or pin a Safehouse-executable `SPACEDOCK_PASEO_BIN` and verify `paseo --version && paseo status --json` against `127.0.0.1:6769` before dispatching `936`.
3. **Export `GIT_SSH_COMMAND='ssh -o BatchMode=yes'` in Commander and worker environments:** `spacedock dispatch build --stamp` invokes `statesync.Publish` -> `git push` (`internal/dispatch/stamp.go:208`), which hangs waiting for an SSH passphrase when `ssh-agent` has no unlocked key inside Safehouse. Export `GIT_SSH_COMMAND='ssh -o BatchMode=yes'` at Commander boot (and pass via `--env` to workers) so push failures fail fast instead of stalling turns.
4. **Live smoke sandbox boundary:** Ordinary dispatched ensigns run without `paseo-safehouse.paseo-access=daemon`, and launching a nested isolated `paseo daemon start` + live `codex`/`antigravity` API calls + reading `~/.codex/auth.json` may be restricted inside Safehouse. The Commander package must explicitly state where `TestLivePaseoDispatchSmoke` and `TestLivePaseoAdoptionSmoke` execute (inside Safehouse worker if permitted, from Commander, or via captain/CI) and follow `docs/dev/README.md:305` (record named unavailability rather than silently skipping).
5. **High-stakes detached adversarial audits & required live CI lanes:**
   - Both `936` (`internal/dispatch/build.go`, `skills/first-officer/**`, `skills/ensign/**`, `.github/workflows/paseo-live-e2e.yml`) and `x9m` (`skills/first-officer/references/fo-write-core.md`, `paseo-fo-dispatch.md`) modify high-stakes surfaces under `docs/dev/README.md:84` — budget a read-only detached adversarial audit on a throwaway checkout (`spike-<slug>`) before merging each PR.
   - Because both tasks touch host-neutral dispatch/FO contract files loaded by every live lane (`docs/dev/README.md:85`), `claude-live`, `codex-live`, and `pi-live` in `runtime-live-e2e.yml` are required green before merge — budget environment approvals and wall time in the Commander handoff.
