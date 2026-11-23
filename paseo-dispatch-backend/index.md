---
id: 936b86qpa8xy4fb6v0n5cafs
title: Opt-in Paseo dispatch backend (--host paseo, CLI-first FO binding, and live smoke)
status: validation
source: captain (2026-11-21)
started: 2026-11-22T17:09:22Z
completed:
verdict:
score: 0.90
sprint: paseo-backend
worktree: .worktrees/spacedock-ensign-paseo-dispatch-backend
issue:
pr:
mod-block:
sprint-readiness: ready
gates:
    version: 1
    records:
        - id: gate:936b86qpa8xy4fb6v0n5cafs:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:936b86qpa8xy4fb6v0n5cafs-backlog-1
              briefing:
                id: briefing:936b86qpa8xy4fb6v0n5cafs:backlog:attempt-1:revision-1
                digest: sha256:e05a1a0037fb0cc9666082e61950d5175302ac815fcfe5c22f1aca2ca1904608
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:936b86qpa8xy4fb6v0n5cafs:backlog:1
                briefing: briefing:936b86qpa8xy4fb6v0n5cafs:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-11-22T17:09:05.573077Z"
                decision: approve
                reason: Captain approved backlog->ideation dispatch on 6769 with agy neon high
              application:
                target-stage: ideation
                state: consumed
        - id: gate:936b86qpa8xy4fb6v0n5cafs:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:936b86qpa8xy4fb6v0n5cafs-ideation-1
              briefing:
                id: briefing:936b86qpa8xy4fb6v0n5cafs:ideation:attempt-1:revision-1
                digest: sha256:f4c799f545c0127d1b2a7df410ba0ff2ed90545dd826873a7baaad2d2d83f4eb
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:936b86qpa8xy4fb6v0n5cafs:ideation:1
                briefing: briefing:936b86qpa8xy4fb6v0n5cafs:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-11-22T18:01:58.330521Z"
                decision: approve
                reason: Captain approved ideation gate after preflight staff-review revisions
              application:
                target-stage: implementation
                state: consumed
---

Let any Paseo-hosted first officer (Claude, Codex, Pi, or Antigravity) dispatch fresh ensigns through Paseo using one host-neutral `--host paseo` contract and 5 `paseo` CLI verbs, backed by a real isolated-daemon live E2E smoke test and an optional (non-required) CI lane using the repo's `codex.exec` model from `internal/release/live_models.txt`.

## Problem

Spacedock currently dispatches workers only through host-specific native subagent tools (Claude `Agent`, Codex `spawn_agent`, Pi `pi-subagents`). When a first officer runs inside Paseo — especially on hosts without a native Spacedock subagent adapter such as Antigravity (`agy`), or across mixed FO/worker runtimes — operators must fall back to `--host codex` with hand-patched prompt overrides, hand-relayed completion signals, or separate per-host Go frontdoors (`27a`, superseded).

Dogfooding this sprint's own dispatches on Paseo 0.17.1 and running the preflight staff review surfaced concrete frictions that a clean Paseo backend must address or bound explicitly:
1. Paseo 0.17.1 splits provider/model from thinking effort (`--provider antigravity/gemini-4-neon --thinking high`), so packing thinking effort into a model ID fails.
2. Without `--host paseo`, Paseo dispatches had to borrow `--host codex` and override the prompt so non-Codex workers would not trip over `$spacedock:ensign` or Codex mailbox directives.
3. Global `~/.gemini/config/plugins/spacedock` can drift behind the repo's skills, so the live harness must stage workflow-free skills from the current checkout via `livePluginDir(t)` rather than relying on global installs or hardcoding `<checkout>/skills/...`.
4. Headless `spacedock state commit` / `git push` can hang on an SSH passphrase prompt when `ssh-agent` has no loaded identities unless `BatchMode=yes` is set.
5. `paseo workspace create` is not idempotent (creates duplicate workspaces on repeated calls), `paseo run --workspace ... --json` writes `Using workspace <id>` to stderr even on exit `0`, and archiving an FO cascade-archives same-workspace children while detaching cross-workspace children.

Because Paseo's CLI flags and JSON shapes evolve quickly, a prose-only FO binding without a real-daemon live E2E smoke test will silently drift.

## Proposed approach

Every mechanism introduced below maps directly to **AC-1 (VALUE)** (a real Paseo FO dispatches a fresh Paseo worker, verifies durable Git commits, and archives the worker) and **AC-2**–**AC-4** that directly serve it:

| Mechanism | Value AC served | Simplest alternative considered | Why alternative is insufficient |
|---|---|---|---|
| `--host paseo` in `internal/dispatch/build.go` + `paseo-ensign-runtime.md` | AC-1, AC-2 | Reuse `--host codex` with FO prompt overrides | Leaks `$spacedock:ensign` and Codex mailbox text into Antigravity/Claude/Pi workers; requires manual prompt patching on every dispatch (observed dogfood friction #2). |
| 5-verb CLI binding `paseo-fo-dispatch.md` wired via `PASEO_AGENT_ID` in `skills/first-officer/SKILL.md` | AC-1, AC-3 | Require Paseo MCP tools on every FO host | Not all FO hosts expose identical Paseo MCP tools; CLI works uniformly across Claude, Codex, Pi, and Antigravity while auto-linking `ParentAgentId` via `PASEO_AGENT_ID`. |
| `TestLivePaseoDispatchSmoke` + shared `paseo_live_helpers_test.go` against an isolated real `paseo` daemon | AC-1 | Offline unit tests + fake ACP/CLI stubs | Fake stubs cannot catch real Paseo flag removals (e.g. removed `--port`/`--no-relay`), stderr/stdout quirks, JSON schema drift, or parent-linkage/archive regressions; extracting package-private helpers also lets `x9m` reuse them cleanly without modifying `936`'s files. |
| Optional `.github/workflows/paseo-live-e2e.yml` reading `codex.exec` from `internal/release/live_models.txt` | AC-1, AC-4 | Local-only manual runs | Local runs drift silently between releases; a dedicated non-blocking (`continue-on-error: true`, `needs: offline`) workflow on existing `CI-E2E-CODEX` + `CODEX_AUTH_JSON`/`OPENAI_API_KEY` gives ongoing signal with zero new secrets, single-source model IDs, and no merge blocking. |

### 1. Exact `--host paseo` envelope contract (`internal/dispatch/build.go` & `dispatch.go`)

- **Host acceptance & refusal (`resolveBuildHost`, `validBuildHost`):**
  - Add `"paseo"` to `validBuildHost(host)`.
  - Update unsupported-host diagnostic to `unsupported host %q (want claude, codex, pi, or paseo)` and ambiguous-markers diagnostic to `...; pass --host claude, codex, pi, or paseo`.
  - **Do NOT add `PASEO_AGENT_ID` to `internal/runtimehost.Detect`.** A Paseo-hosted FO also carries its underlying host marker (`CLAUDECODE`, `CODEX_THREAD_ID`, or `PI_CODING_AGENT_DIR`); adding `PASEO_AGENT_ID` to `markerTable` would trigger host-ambiguity errors for ordinary sessions running inside Paseo. `paseo-fo-dispatch.md` always passes `--host paseo` explicitly, which short-circuits marker detection.
- **Mode & model rules:**
  - `--bare-mode` on `--host paseo` MUST fail before artifact creation with exit `2` and stderr:
    `error: bare_mode is unsupported on host paseo; Paseo worker.spawn requires a named worker label`
  - `--stamp` works identically to other hosts (stamps `started`/`worktree`, commits/syncs state, creates code worktree before assembly).
  - `--advance` remains supported by the host-uniform advance path (though `paseo-fo-dispatch.md` v1 dispatches fresh only).
  - Stage/default `model:` falls into the existing non-`claude` ignore-with-note branch: emits `"model": null` in stdout JSON and logs `[build] declared model '%s' ignored on host paseo: outside paseo's dispatch-settable model space; emitting model=null` to stderr.
- **Dispatch file path disambiguation:**
  - When `host == "paseo"` and `PASEO_AGENT_ID` is non-empty, prepend `pathSafeSessionToken(os.Getenv("PASEO_AGENT_ID")) + "-"` to `dispatchFileName` (mirroring Claude `mergedMode`'s `CLAUDE_CODE_SESSION_ID` keying) so concurrent Paseo FOs never collide on `/tmp/spacedock-dispatch/{derivedName}.md`. When `PASEO_AGENT_ID` is empty, keep plain `{derivedName}.md`.
- **Stdout JSON envelope (`buildOutput`, `schema_version: 2`):**
  - `schema_version`: `2`
  - `subagent_type`: `"spacedock:ensign"` (or stage `agent:` override)
  - `description`: `"{entityTitle}: {stage}"`
  - `fetch_commands`: `["<abs-spacedock> dispatch show-stage-def --workflow-dir <abs-wf> --stage <stage>"]`
  - `dispatch_file_path`: `"/tmp/spacedock-dispatch/[<paseo-token>-]<derivedName>.md"`
  - `prompt`: `Load the Spacedock ensign skill (with references/ensign-shared-core.md and references/paseo-ensign-runtime.md); then Read <dispatch_file_path> and treat its content as your assignment.`
  - `model`: `null`
  - `name`: `derivedName` (always present; the invariant `{worker_base}` used by `spacedock.worker` and `spacedock.epoch`)
  - `run_in_background`, `agent`, `skill`: omitted (`omitempty`).
- **`firstActionBlock("paseo")` (exact dispatch-body opening):**
  ```markdown
  ## First action

  Before anything else, load the Spacedock `ensign` skill (`SKILL.md`, `references/ensign-shared-core.md`, and `references/paseo-ensign-runtime.md`) using your host's skill loader or by reading those files directly, then read this dispatch file and treat its content as your stage-specific assignment.

  This assignment is delivered through Paseo. Work in the assigned code directory (and commit entity/report updates path-scoped in the split-root state checkout when present). Do not invoke host-specific subagent or mailbox tools.
  ```
- **`completionSignalBlock("paseo", entityTitle, stage, entityFileRef)`** (wire `host == "paseo"` into the `if !bareMode && (mergedMode || host == "codex" || host == "paseo")` guard in `runBuildFields`):
  ```markdown
  ### Completion Signal

  When you finish (after all commits and stage report writes are done), end your turn with one concise final message:

      Done: {entityTitle} completed {stage}. Report written to {entityFileRef}.

  The first officer verifies completion from committed Git state and the entity stage report after your turn becomes idle. Do not call host-specific mailbox tools.
  ```
- **Banned host-specific tokens:** Neither `prompt` nor the written dispatch body for `--host paseo` may contain `Skill(skill=`, `SendMessage`, `$spacedock:ensign`, `/skill:ensign`, `pi-subagents`, `TeamCreate`, `TeamDelete`, `Agent(`, or `FO mailbox`.

### 2. Exact `paseo-ensign-runtime.md` contract & `skills/ensign/SKILL.md` wiring

Create `skills/ensign/references/paseo-ensign-runtime.md` following `docs/runtime-support.md`'s ensign binding-block shape (positive bindings only, no re-narration of `ensign-shared-core.md`):
- `Fresh dispatch` -> Load `ensign-shared-core.md` and `paseo-ensign-runtime.md`, Read the `/tmp/spacedock-dispatch/...md` pointer file, run `### Fetch commands` verbatim, and follow the assignment together with any ordinary FO scope notes appended to the prompt.
- `Git sync in headless/sandboxed workers` -> Commit code to the assigned worktree branch and entity/report updates path-scoped in the split-root state checkout. When pushing over SSH without an unlocked agent, use `GIT_SSH_COMMAND='ssh -o BatchMode=yes'` so git fails fast rather than hanging on a passphrase prompt; if push fails due to sandbox or missing credentials, keep the local path-scoped commit clean and note the push failure in `### Summary`.
- `Clarification` -> Ordinary Paseo ensigns run without a daemon send grant; if genuinely blocked, commit safe partial progress (or leave worktree clean), record the blocker in `## Stage Report` (`FAILED:` + `### Summary`), and end your turn so the FO observes `idle` and intervenes.
- `Completion signal` -> After all commits and stage-report writes succeed, end your final turn with one single-line message: `Done: {entity title} completed {stage}. Report written to {entity_file_path}.` Do not call host-specific mailbox tools; Git commits and the stage report are the completion gate.

**Concrete before/after for `skills/ensign/SKILL.md`** (`PASEO_AGENT_ID` MUST appear first so a Paseo-launched Codex/Claude/Pi worker does not match its inner host marker first):
```diff
 Load the runtime adapter for your platform:
+- Paseo (`PASEO_AGENT_ID` env var is set, or the dispatch prompt names `references/paseo-ensign-runtime.md`): read `references/paseo-ensign-runtime.md`
 - Claude Code (`CLAUDECODE` env var is set): read `references/claude-ensign-runtime.md`
 - Codex (`CODEX_THREAD_ID` env var is set): read `references/codex-ensign-runtime.md`
 - Pi (`PI_CODING_AGENT_DIR` is set, or this session is running under Pi without the Claude/Codex markers above): read `references/pi-ensign-runtime.md`
```

### 3. Exact 5-verb `paseo-fo-dispatch.md` CLI contract & `skills/first-officer/SKILL.md` selection

**Selection wiring (concrete before/after for `skills/first-officer/SKILL.md`; do NOT modify `first-officer-shared-core.md`, which sits at 23,497 / 23,500 bytes under `TestFOInstructionComponentCaps`):**
```diff
 ## Runtime adapter
 
 Load the runtime adapter for your platform:
 - Claude Code (`CLAUDECODE` env var is set): read exactly `{first_officer_base}/references/claude-first-officer-runtime.md`; do not list or search the references directory.
 - Codex (`CODEX_THREAD_ID` env var is set): read `references/codex-first-officer-runtime.md`
 - Pi (`PI_CODING_AGENT_DIR` is set, or this session is running under Pi without the Claude/Codex markers above): read `references/pi-first-officer-runtime.md`
+
+**Paseo dispatch backend:** When `PASEO_AGENT_ID` is set in the environment or the captain explicitly selects Paseo dispatch, bind worker dispatch at first dispatch by reading `{first_officer_base}/references/paseo-fo-dispatch.md` alongside `{first_officer_base}/references/fo-dispatch-core.md` instead of the host's native subagent transport (on hosts without a dedicated adapter above, communicate with the captain via direct conversation text).
```
Note: uses `{first_officer_base}/references/paseo-fo-dispatch.md` (not `@references/`), preserving `TestFirstOfficerEntryEagerlyImportsOnlySharedCore` while satisfying `TestUserSkillReferenceClosureResolves`.

**Contractlint closure constraint (Non-blocking #1):** `paseo-fo-dispatch.md` MUST NOT contain `references/paseo-ensign-runtime.md` or `references/ensign-shared-core.md` (or any other `references/...` path that lives under `skills/ensign/` instead of `skills/first-officer/`). `TestBootResidentDeferredLoadPointsResolve` resolves every `references/*.md` match inside `deferredModuleBodies` relative to `skills/first-officer/`, so naming an ensign reference path in `paseo-fo-dispatch.md` breaks closure.

**Exact 5-verb Paseo 0.17.1 CLI contract (`skills/first-officer/references/paseo-fo-dispatch.md`):**

- **Common CLI invariants (verified against Paseo 0.17.1 `packages/cli/src/`):**
  - Always pass `--json`. Successful commands exit `0` with JSON on stdout; CLI/daemon errors exit `1` with `{"error":{"code":"...","message":"...","details":"..."}}` on stderr.
  - **Stdout-only JSON parsing (M4):** `paseo run --workspace <id> --json` writes unconditional non-JSON `Using workspace <id>` (and workspace-creation tips) to stderr even on exit `0` (`packages/cli/src/commands/agent/run.ts:553, 563, 584-588`). Callers MUST parse stdout only — NEVER combine streams (`2>&1 | jq` fails JSON decoding on exit `0`). Capture stderr separately (or redirect `2>/dev/null` only after checking exit status / saving stderr to a temp file for failure diagnostics).
  - Preserve `PASEO_AGENT_ID` in the environment (Paseo CLI's `resolveCallerAgentId` passes it as `callerAgentId`, auto-stamping `paseo.parent-agent-id` and `paseo.parent-server-id` on spawned children). Do not pass `--home` or `--host` during normal FO runs (inherit the session's daemon target).
  - Never pass `paseo-safehouse.paseo-access=daemon` to ordinary ensigns (reserved for FOs until Paseo ships a parent-only daemon ACL).
  - **Provider / model / thinking split:** Always pass an explicit `--provider <provider>/<model>` (never rely on `build.model`, which is `null`), and pass thinking effort separately via `--thinking <id>` when requested (e.g., `--provider antigravity/gemini-4-neon --thinking high`, or `--provider codex/<model>`).
  - **Frozen shared label contract (M7, shared with `x9m`):**
    - `spacedock.worker` = `{worker_base}` (= `envelope.name` emitted by `spacedock dispatch build --host paseo` / `spacedock dispatch name`). Invariant across cycles and retries (always <= 56 chars, valid under `internal/dispatch/names.go`); NEVER append `-retry` or `-cycleN` to `spacedock.worker`.
    - `spacedock.epoch` = `{worker_base}-c{cycle}-a{attempt}` (`{worker_base}-c1-a1` for v1 fresh dispatch in `936`; extended in `x9m` for `cycle >= 1`, `attempt in {1, 2}`).

- **Pre-spawn worktree, project-ID, and Safehouse handling:**
  1. Build envelope: `${SPACEDOCK_BIN:-spacedock} dispatch build --host paseo --stamp --workflow-dir {workflow_dir} --entity-path {entity_path} --stage {stage} --checklist-file -` (omit `--bare-mode`). Retain `name` (`{worker_base}`), `description`, `dispatch_file_path`, `prompt`. Set `epoch="${name}-c1-a1"`.
  2. Resolve physical `code_cwd` from structured status (never parse the dispatch body):
     - Run `${SPACEDOCK_BIN:-spacedock} status --workflow-dir {workflow_dir} --read {entity_path} --json` and read `frontmatter.worktree`.
     - If non-empty: join relative paths against `git -C {workflow_dir} rev-parse --show-toplevel`, keep absolute paths as-is, and canonicalize symlinks (`realpath` / `pwd -P`) to produce `code_cwd`.
     - If empty (non-worktree stage such as `ideation`): use canonical `git -C {workflow_dir} rev-parse --show-toplevel`.
  3. **Spacedock owns worktrees:** Never pass `--isolation worktree` or `--new-workspace worktree` to Paseo (Git exits 128 on an occupied branch or creates a mismatched branch). Always adopt the existing physical `code_cwd` via find-or-create (`workspace ls` then `workspace create --isolation local`). Paseo classifies an adopted Git worktree as `isolation: "worktree"` while keeping exact `cwd`.
  4. **Project ID:** When `workspace ls` finds no matching workspace, resolve `project_id` by matching canonical `git -C {workflow_dir} rev-parse --show-toplevel` against `path` in `paseo project ls --json` (`[{projectId, name, kind, path}]`). When matched, pass `--project "$project_id"` to `workspace create`; when unmatched (e.g. fresh isolated test repo), omit `--project` so Paseo registers the checkout automatically.
  5. **Safehouse sibling state checkout:** In split-root worktree stages under `paseo-safehouse`, the worker writes both `code_cwd` (`.worktrees/spacedock-ensign-{slug}`) and the sibling state checkout (`{workflow_dir}/.spacedock-state`). If Safehouse blocks the state-checkout write on a non-`/tmp` repo, the worker reports the exact diagnostic and stops rather than retrying silently.

- **The 5 CLI verbs (flags, JSON fields, exit codes):**

| # | Verb & exact invocation | JSON stdout fields | Exit codes & branching |
|---|---|---|---|
| 1 | **Find-or-create workspace (M3):** First run `paseo workspace ls --json` and match an active row whose canonical `cwd == $code_cwd`. Only if none matches, run `paseo workspace create --isolation local --path "$code_cwd" [--project "$project_id"] --title "$description" --json` | `workspace ls`: `WorkspaceRow[]`; `workspace create`: `WorkspaceRow` (`workspaceId` (`wks_...`), `project`, `name`, `isolation` (`"local"\|"worktree"`), `cwd` (unshortened), `background`) | `0`: verify matched/returned `cwd == $code_cwd`, retain `workspaceId`. (`workspace create` calls `createWorkspaceForDirectory`, NOT find-or-create, so unconditional calls leak duplicate workspaces on non-worktree stages, sticky-worktree transitions, and retries.) `1`: `DAEMON_NOT_RUNNING` / `WORKSPACE_CREATE_FAILED` -> abort spawn before creating agent. |
| 2 | `paseo run --no-wait --workspace "$workspace_id" --provider "$provider_model" [--thinking "$thinking_id"] [--mode "$mode_id"] --title "$description" --label "role=ensign" --label "spacedock.worker=$name" --label "spacedock.epoch=${name}-c1-a1" --json "$prompt"` | `agentId` (UUID), `status` (`"created"\|"running"`), `provider`, `cwd`, `title` | `0`: **parse stdout only** (M4: stderr prints `Using workspace <id>` on exit `0`); record `agentId`, `workspaceId`, `spacedock.worker=$name`, `spacedock.epoch=${name}-c1-a1`, `cwd`, requested/returned provider/model/thinking/mode before waiting. CLI `--no-wait` leaves `notifyOnFinish=false`. `1`: `MISSING_PROVIDER`, `INVALID_OPTIONS`, `INVALID_THINKING_OPTION`, `WORKSPACE_NOT_FOUND`, `AGENT_CREATE_FAILED`. |
| 3 | `paseo wait "$agent_id" --timeout 300 --json` | `agentId`, `status` (`"idle"\|"timeout"\|"permission"\|"error"`), `message` (includes last 5 activity items on `idle`/`timeout`) | `0` for ALL FOUR statuses — ALWAYS branch on JSON `status`! `timeout` -> keep waiting (normal progress, 0 retries). `idle` -> wake signal ONLY; verify durable Git commits + `## Stage Report: {stage}` before advancing. `permission` -> `paseo inspect` & escalate. `error` -> `paseo inspect` & hold. `1`: `AGENT_NOT_FOUND`, `INVALID_TIMEOUT`, `WAIT_FAILED`. |
| 4 | `paseo inspect "$agent_id" --json` and `paseo ls --global [--all] --label "spacedock.worker=$name" --json` | `inspect` (PascalCase): `Id`, `Name`, `Provider`, `Model`, `Thinking`, `Status`, `Archived`, `ArchivedAt`, `Mode`, `Cwd` (unshortened), `CreatedAt` (ISO-8601), `UpdatedAt`, `LastUsage` (`{InputTokens, OutputTokens, CachedTokens, CostUsd}\|null`), `PendingPermissions` (`[{id, tool}]`), `Worktree`, `ParentAgentId`. `ls` (camelCase array): `[{id, shortId, name, provider, thinking, status, cwd, created, background, labels, workspaceId}]` (note: `ls.cwd` is `~`-shortened and `ls.created` is relative time; use `inspect` for exact path/timestamp/archived fields). | `0`: inspect/roster snapshot. `1`: `AGENT_NOT_FOUND`, `INSPECT_FAILED`, `LIST_AGENTS_FAILED`. |
| 5 | `paseo archive "$agent_id" --json` (`--force` ONLY on explicit supersede of a running worker) | `agentId`, `status` (`"archived"`), `archivedAt` | `0`: worker archived. `1`: `AGENT_RUNNING` (without `--force`, fail-closed guard), `AGENT_ALREADY_ARCHIVED`, `AGENT_NOT_FOUND`, `ARCHIVE_FAILED`. **Cascade-archive guard (M5):** Paseo's `archiveAgentUnlocked` (`agent-manager.ts:2319, 818-827`) cascade-kills and archives every child in the same workspace (`parent.workspaceId === child.workspaceId` — all non-worktree stages share the repo-root workspace!) and detaches cross-workspace children by deleting `paseo.parent-agent-id`. Therefore NEVER archive `$PASEO_AGENT_ID`, and NEVER archive an FO while any child worker is running (use `paseo stop` to interrupt an FO without killing/detaching its children). |

- **`## Runtime implementation` bindings block in `paseo-fo-dispatch.md`** (matching `docs/runtime-support.md`):
  - `«worker.spawn»` -> `dispatch build --host paseo --stamp` -> resolve canonical `code_cwd` -> Verb 1 find-or-create (`workspace ls` then `workspace create --isolation local` if absent) -> Verb 2 (`run --no-wait`, stdout-only JSON parse, `spacedock.worker=$name`, `spacedock.epoch=${name}-c1-a1`). Forward the intact helper pointer verbatim (optional ordinary FO scope notes appended after a newline).
  - `«addressable-worker»` -> ABSENT in v1 (fresh one-shot workers only; worker reuse and child-to-parent messaging are out of scope).
  - `«async-dispatch»` -> ASYNC (`paseo run --no-wait` returns `agentId`; bounded polling via Verb 3 `paseo wait --timeout 300 --json`).
  - `«worker-identity»` -> record `agentId`, `workspaceId`, `spacedock.worker` (`{worker_base}`), `spacedock.epoch` (`{worker_base}-c1-a1`), `cwd`, and requested/returned `provider`, `model`, `thinking`, `mode`; canonical model space is `<provider>/<model>` (+ optional `--thinking <id>`).
  - `«completion-signal»` -> Verb 3 `status == "idle"` wakes verification; Git state (committed worktree deliverable + path-scoped `## Stage Report: {stage}` commit on canonical entity + clean `git status --porcelain`) is the completion gate.
  - `«worker.shutdown»` -> Verb 5 `paseo archive "$agent_id" --json` after Git verification (or `--force` on supersede); guard `agentId != $PASEO_AGENT_ID` and never archive an FO with live children (M5 same-workspace cascade-archive rule).
  - `«context-budget»` -> ABSENT (`inspect.LastUsage` omits context-window max and v1 dispatches fresh only).
  - `«roster-reconcile»` -> Verb 4 (`paseo ls --global --label "spacedock.worker=..." --json` + `paseo inspect "$agent_id" --json`).

### 4. Shared live-test helper seam, `TestLivePaseoDispatchSmoke`, and optional CI workflow

- **Shared package-private helper seam (`internal/ensigncycle/paseo_live_helpers_test.go`, `//go:build live`, M8):**
  Extract the isolated-daemon and fixture helpers into a dedicated file so `x9m` (`paseo_adoption_live_test.go`) consumes them with zero edits to `936`'s files:
  - `startIsolatedPaseoDaemon(t *testing.T, opts ...paseoDaemonOption) *paseoDaemon`
    - Resolves `SPACEDOCK_PASEO_BIN` (default `exec.LookPath("paseo")`); skips locally (`t.Skip`) when missing unless `SPACEDOCK_PASEO_LIVE_REQUIRED=1` (`t.Fatal`).
    - Resolves default provider from SSOT: `"codex/" + release.LiveModels["codex.exec"]` (never hardcodes a model literal in Go!), overridable via `SPACEDOCK_PASEO_LIVE_PROVIDER`, `SPACEDOCK_PASEO_LIVE_THINKING`, `SPACEDOCK_PASEO_LIVE_FO_PROVIDER`, `SPACEDOCK_PASEO_LIVE_TIMEOUT_MINUTES` (default `12`).
    - Resolves canonical temp root via `filepath.EvalSymlinks(t.TempDir())` (eliminating macOS `/var` -> `/private/var` symlink drift) and sets `paseoHome := filepath.Join(canonicalRoot, "paseo-home")`.
    - Picks an ephemeral free TCP port on `127.0.0.1:0` (rejecting `6767` and `6768`).
    - Scrubs ambient env: strips all `PASEO_*` vars, `PORT`, and foreign runtime markers (`CLAUDECODE`, `CODEX_THREAD_ID`, `PI_CODING_AGENT`, `PI_CODING_AGENT_DIR`); sets `PASEO_HOME=paseoHome`, `PATH` prefixed with `filepath.Dir(spacedockBin)`, and `GIT_SSH_COMMAND=ssh -o BatchMode=yes`.
    - **Isolated `CODEX_HOME` & workflow-free skill staging (M6 + Non-blocking #2):**
      - For `codex/*` providers: creates isolated `CODEX_HOME` via `newCodexLiveIsolatedHome(t, repo, artifactDir)` (outside `$TMPDIR` on macOS so Codex sandboxing accepts it), seeds `CODEX_AUTH_JSON` (`0600`) / `OPENAI_API_KEY` / copied `~/.codex/auth.json` (skip locally if absent unless `SPACEDOCK_PASEO_LIVE_REQUIRED=1`), and installs the workflow-free `livePluginDir(t)` copy into that isolated `CODEX_HOME` via `cli.WriteCodexLocalMarketplace` + `codex plugin marketplace add` + `codex plugin add spacedock@spacedock`.
      - For non-Codex local overrides (e.g. `antigravity/*`): copies `livePluginDir(t)`'s workflow-free `skills/` into the fixture repo's `.agents/skills/` (or isolated host skill dir) so named skill discovery resolves from the staged current checkout rather than stale global `~/.gemini/config/plugins/spacedock` or raw `<checkout>/skills/...`.
    - Configures 5 daemon keys via `paseo daemon config set <key> <val> --home "$paseoHome" --json`: `daemon.listen=127.0.0.1:<freePort>`, `daemon.relay.enabled=false`, `features.webUi.enabled=false`, `features.dictation.enabled=false`, `features.voiceMode.enabled=false`.
    - Runs `paseo daemon start --home "$paseoHome" --timeout 30 --json`, registers `t.Cleanup` (`paseo daemon stop --home "$paseoHome" --force --timeout 5`, verify `localDaemon == "stopped"`, copy `paseo.log` into `SPACEDOCK_LIVE_ARTIFACT_DIR` on failure), and polls `paseo status --home "$paseoHome" --json` up to 30s until `connectedDaemon == "reachable"` and `workerPid > 0`.
  - `(*paseoDaemon).runJSON(t *testing.T, dst any, args ...string)` (M4):
    - Executes `paseo` with `--home "$paseoHome"`, attaching **separate `cmd.Stdout` and `cmd.Stderr` `bytes.Buffer`s** (NEVER `cmd.CombinedOutput()`).
    - Fails with both exit code and stderr on non-zero exit; decodes JSON strictly from `stdout` (ignoring benign exit-0 stderr noise such as `Using workspace wks_...`).
  - `setupPaseoSplitRootFixture(t *testing.T, root string) paseoSplitRootFixture`:
    - Initializes `root/repo` (no `origin` remote) with `docs/smoke/README.md` (`state: .spacedock-state`, `id-style: slug`, stages `implementation` with `worktree: true` and `done` terminal), linked state worktree `docs/smoke/.spacedock-state` containing `paseo-smoke/index.md` (`status: implementation`, empty `worktree:`), and returns `{RepoRoot, WorkflowDir, StateRoot, EntityPath, Slug}`.

- **Live smoke (`internal/ensigncycle/paseo_live_smoke_test.go`, `//go:build live`):**
  - Calls `d := startIsolatedPaseoDaemon(t)` and `fx := setupPaseoSplitRootFixture(t, d.CanonicalRoot)`.
  - Adopts `fx.RepoRoot` via `d.runJSON(t, &foWs, "workspace", "create", "--isolation", "local", "--path", fx.RepoRoot, "--title", "paseo-smoke-fo", "--json")`.
  - Launches real FO via `d.runJSON(t, &foRun, "run", "--no-wait", "--workspace", foWs.WorkspaceID, "--provider", d.FOProvider, ... "--env", "SPACEDOCK_BIN="+d.SpacedockBin, "--label", "role=first-officer", "--json", foPrompt)`.
  - **`$foPrompt` (M6 — no hardcoded `<checkout>/skills/...` paths):** invokes the installed `first-officer` skill by name (never naming `paseo-fo-dispatch.md` directly, exercising `skills/first-officer/SKILL.md`'s `PASEO_AGENT_ID` selection), specifies single-entity implementation dispatch for `paseo-smoke` with worker provider `d.WorkerProvider` (+ `--thinking` if set), requires forwarding `dispatch build --host paseo --stamp`'s emitted `prompt` verbatim with zero hardcoded checkout skill paths (exercising `skills/ensign/SKILL.md`'s `PASEO_AGENT_ID`-first precedence and installed skill discovery in the consumer fixture), and requires verifying the commits (`receipt.txt` on task branch + path-scoped `## Stage Report: implementation` in `.spacedock-state`) and archiving the verified worker.
  - Waits on FO via `d.runJSON(t, &waitRes, "wait", foRun.AgentID, "--timeout", "60", "--json")` loop up to `d.Timeout` until `status == "idle"`.
  - **Durable assertions (transcript-independent):**
    1. `frontmatter.worktree` is stamped and resolves to an existing Git worktree checked out on a non-`main` branch.
    2. Task branch latest commit touches `receipt.txt`, `receipt.txt` contains the smoke marker, and `git -C "$codeWorktree" status --porcelain` is empty.
    3. Canonical `paseo-smoke/index.md` contains `## Stage Report: implementation`, `- DONE:`, `### Summary`, and preserves YAML frontmatter `status: implementation`; `git -C "$stateRoot" log -1 --name-only` touches ONLY `paseo-smoke/index.md`; `git -C "$stateRoot" status --porcelain` is empty.
    4. `paseo ls --all --global --label "spacedock.worker=paseo-smoke-implementation" --json` returns exactly 1 worker carrying `spacedock.epoch=paseo-smoke-implementation-c1-a1`; `paseo inspect "$workerId" --json` confirms `ParentAgentId == foRun.AgentID` and `Archived == true` (`ArchivedAt != null`); FO itself remains `Archived == false`.
    5. `t.Cleanup` stops the isolated daemon cleanly (`localDaemon == "stopped"`).

- **Optional non-required CI workflow (`.github/workflows/paseo-live-e2e.yml`, M8):**
  - Top-level `permissions: contents: read` (matching `runtime-live-e2e.yml`).
  - Triggers: `workflow_dispatch` (optional inputs `paseo_version`, `codex_version`, `paseo_provider` empty = `codex/<codex.exec>`) and `pull_request` on `main` (`opened`, `synchronize`, `reopened` — never `pull_request_target`, preserving fork-PR secret isolation) with `paths:` scoped to Paseo/dispatch/skill/workflow surfaces (so unrelated PRs do not queue an unnecessary deployment approval).
  - Job `offline` (`runs-on: ubuntu-latest`, no environment/secrets): checks out repo, reads `codex.exec` from `internal/release/live_models.txt` (emitting `outputs.codex_exec` — zero hardcoded model literals!), runs `go test ./internal/dispatch ./internal/contractlint ./internal/ensigncycle` as the secret-free gate.
  - Job `paseo-live`: `needs: offline`, `runs-on: ubuntu-latest`, `continue-on-error: true` (non-blocking optional lane), `environment: CI-E2E-CODEX`.
  - Steps:
    1. Secret check accepting `CODEX_AUTH_JSON` OR `OPENAI_API_KEY` (failing only when both are empty, matching `runtime-live-e2e.yml:416`).
    2. `actions/checkout@v5`, `actions/setup-go@v6` (`1.26`), `actions/setup-node@v6` (`24`).
    3. Install `@openai/codex` and `@getpaseo/cli` (honoring optional `codex_version`/`paseo_version` inputs).
    4. Build `./spacedock` and export `SPACEDOCK_BIN=$(pwd)/spacedock` and `SPACEDOCK_REPO_ROOT=$GITHUB_WORKSPACE` to `$GITHUB_ENV` (matching `runtime-live-e2e.yml:498-500`).
    5. Install `gotestsum` via `.github/scripts/install-gotestsum.sh`, configure git identity, print `paseo --version`, `codex --version`, and `go version`.
    6. Run:
       ```bash
       SPACEDOCK_PASEO_LIVE_REQUIRED=1 \
       SPACEDOCK_PASEO_LIVE_PROVIDER="${PASEO_PROVIDER:-codex/${{ needs.offline.outputs.codex_exec }}}" \
       gotestsum --jsonfile paseo-live-smoke-detail.jsonl --format pkgname -- \
         -tags live -count=1 -timeout 15m -run '^TestLivePaseoDispatchSmoke$' ./internal/ensigncycle
       ```
    7. Always upload `paseo-live-smoke-detail.jsonl` and `live-artifacts/paseo/**`.

### 5. Concrete CLI help & documentation diffs (ASD-STE100 compliant)

- **`internal/dispatch/dispatch.go`** (`printUsage` and `printBuildUsage`):
  ```diff
  -  spacedock dispatch build --workflow-dir DIR --entity-path FILE --stage STAGE --checklist-file FILE|- [--host claude|codex|pi]
  +  spacedock dispatch build --workflow-dir DIR --entity-path FILE --stage STAGE --checklist-file FILE|- [--host claude|codex|pi|paseo]
  ...
  -    [--host claude|codex|pi] [--bare-mode] [--feedback-reflow] [--stamp | --advance]
  +    [--host claude|codex|pi|paseo] [--bare-mode] [--feedback-reflow] [--stamp | --advance]
  ...
  -  --host HOST                   Override the runtime host (claude|codex|pi). Defaults to the detected runtime.
  +  --host HOST                   Override the runtime host (claude|codex|pi|paseo). Defaults to the detected runtime.
  ...
  -  --bare-mode                   Emit the bare sequential shape (no name, no run_in_background); unsupported on host=codex.
  +  --bare-mode                   Emit the bare sequential shape (no name, no run_in_background); unsupported on host=codex and host=paseo.
  ```
- **`docs/runtime-live-ci.md`** (append a concise `### Paseo live smoke` subsection under `### Local live execution`, following ASD-STE100 and never restating a literal model ID per `docs/runtime-live-ci.md:135`; do not alter existing `claude`/`codex`/`pi` common-suite lines):
  ```markdown
  ### Paseo live smoke

  Run the isolated-daemon Paseo smoke locally with the `codex.exec` model from `internal/release/live_models.txt`:

  ```bash
  go test -tags live -count=1 -timeout 15m -run '^TestLivePaseoDispatchSmoke$' ./internal/ensigncycle -v
  ```

  Override the binary, provider, or thinking level with `SPACEDOCK_PASEO_BIN`, `SPACEDOCK_PASEO_LIVE_PROVIDER`, and `SPACEDOCK_PASEO_LIVE_THINKING`.

  In CI, `.github/workflows/paseo-live-e2e.yml` runs `TestLivePaseoDispatchSmoke` under `CI-E2E-CODEX` with `CODEX_AUTH_JSON` or `OPENAI_API_KEY`. This job is optional (`continue-on-error: true`) and does not block merges.
  ```

## Risk evidence

**No new ideation spike needed for core dispatch transport** — the end-to-end Paseo CLI + Safehouse execution path is already proven by three recorded live probes plus direct source inspection of Paseo 0.17.1 CLI/server (`packages/{cli,server}/src/`) and the 2026-11-22 preflight staff review (`_roadmap/paseo-backend/staff-review.md`):
- **2026-11-18 probe** (`_roadmap/paseo-backend/probe-report.md`, `staff-report.md`, `run-log.md`): verified `dispatch build --stamp`, adopting an existing worktree via `isolation=local` (Paseo returned `isolation=worktree` and preserved `cwd`; Git exited 128 on duplicate branch checkout), `paseo ls --global --label`, `paseo wait --timeout 1 --json` exiting 0 with `status=timeout`, durable code + state commits, and clean archival.
- **2026-11-21 live probe 1 (`agy` FO -> `paseo` CLI -> `safehouse-codex` ensign):** completed in 49 s via CLI alone with automatic `ParentAgentId` linkage via `PASEO_AGENT_ID`, code commit (`a9fb598`), split-root state commit (`cf64887`), `paseo wait` idle, and clean archive.
- **2026-11-21 live probe 2 (sandboxed `antigravity` `gemini-4-neon` ensign):** ran on Paseo `launch-wrap-hook` (`9f9f39ce1`) + `paseo-safehouse` (`0ff8d55`, `route: wrap`) with host-neutral prompt, verified under `safehouse --enable=keychain`, code commit (`c24eda1`), and split-root state commit (`9b1476c`).
- **Paseo 0.17.1 source verification & preflight staff review (`staff-review.md`, `8447132f3`):** verified all 5 verbs' exact flags, JSON fields, exit codes, and `resolveCallerAgentId`, and caught 6 concrete findings folded above: `workspace create` non-idempotency (M3 -> find-or-create via `workspace ls`), `paseo run --workspace --json` exit-0 stderr noise (M4 -> stdout-only JSON parsing), same-workspace cascade-archive (M5 -> explicit FO archive guard), installed-skill discovery via `livePluginDir(t)` (M6), frozen `spacedock.worker`/`spacedock.epoch` label contract (M7), and shared `paseo_live_helpers_test.go` + `live_models.txt` SSOT + CI conventions (M8).
- **Dogfood frictions captured on 2026-11-22:** (1) Paseo 0.17.1 split `--provider <p>/<m>` + `--thinking <id>` -> encoded directly in `paseo-fo-dispatch.md` and `SPACEDOCK_PASEO_LIVE_THINKING`; (2) missing `--host paseo` forcing `--host codex` overrides -> eliminated by AC-2; (3) stale global `~/.gemini/config/plugins/spacedock` -> live smoke stages workflow-free current-checkout skills via `livePluginDir(t)`; (4) SSH passphrase hang when `ssh-agent` has no identities -> mitigated in `paseo-ensign-runtime.md` and live harness via `GIT_SSH_COMMAND='ssh -o BatchMode=yes'`.
- **Two implementation-time checks to run first:**
  1. **Non-`/tmp` Safehouse sibling state-checkout check:** verify a `paseo-safehouse`-wrapped worker in a non-`/tmp` repo worktree (`.worktrees/spacedock-ensign-*`) can write and commit to the sibling `docs/*/.spacedock-state` checkout without a Seatbelt denial; record the exact result in implementation/validation evidence.
  2. **CI `@getpaseo/cli` + `CODEX_AUTH_JSON` check:** verify `npm install -g @getpaseo/cli` + isolated `paseo daemon start` + `codex/<codex.exec>` on Ubuntu in `.github/workflows/paseo-live-e2e.yml` (or locally with `SPACEDOCK_PASEO_LIVE_PROVIDER`).

## Out of scope

- FO restart adoption, persisted worker records, and wait-timeout/dead-worker/lost-receipt recovery matrix (tracked in `paseo-worker-adoption-and-recovery`).
- Remote Paseo workers (`serverId`) and child-to-parent mid-run questions (blocked on Paseo parent-only daemon ACL).
- Fake ACP stub servers or Paseo version-pinning machinery.
- Worker reuse across stages (`--advance` remains host-uniform in `build.go`, unused by `paseo-fo-dispatch.md` v1).
- Global Antigravity skill sync UX (`~/.gemini/config/plugins/spacedock` freshness) and binary-wide `spacedock state commit` SSH `BatchMode` hardening (tracked separately outside this sprint).

## Expected surface and tolerance

Estimate net LOC change: +850 (+863 insertions / -13 deletions), across 17 files.
Tolerance: ±30% net LOC (+595 to +1105 net LOC), 13 to 20 files.

Planned files:
- `internal/dispatch/build.go` (+45 / -5)
- `internal/dispatch/dispatch.go` (+4 / -4)
- `internal/dispatch/build_paseo_host_test.go` (new, +125 / -0)
- `internal/dispatch/build_host_test.go` (+2 / -2)
- `internal/dispatch/build_json_ergonomics_test.go` (+2 / -2)
- `skills/ensign/SKILL.md` (+1 / -0)
- `skills/ensign/references/paseo-ensign-runtime.md` (new, +25 / -0)
- `skills/first-officer/SKILL.md` (+3 / -0)
- `skills/first-officer/references/paseo-fo-dispatch.md` (new, +130 / -0)
- `internal/contractlint/boot_resident_closure_test.go` (+2 / -0)
- `internal/contractlint/structural_checks_test.go` (+10 / -0)
- `internal/contractlint/paseo_live_lane_test.go` (new, +75 / -0)
- `internal/ensigncycle/paseo_live_controls_test.go` (new, +85 / -0)
- `internal/ensigncycle/paseo_live_helpers_test.go` (new, +125 / -0 — shared package-private daemon/JSON/fixture seam consumed by `x9m`)
- `internal/ensigncycle/paseo_live_smoke_test.go` (new, +95 / -0)
- `.github/workflows/paseo-live-e2e.yml` (new, +115 / -0)
- `docs/runtime-live-ci.md` (+16 / -0)

**Declared observable semantic changes:**
- **Command grammar:** `spacedock dispatch build --host` gains `paseo` alongside `claude|codex|pi`; `--host paseo --bare-mode` exits `2` before creating an artifact; unsupported/ambiguous host error messages list `paseo`.
- **Stored formats:** No entity frontmatter or state schema change. Dispatch filename under `/tmp/spacedock-dispatch/` prepends sanitized `PASEO_AGENT_ID` prefix when set on `--host paseo`.
- **Authority:** Unchanged (FO owns frontmatter/state transitions and worker archival; ensigns own assigned code commits and path-scoped stage reports; ordinary ensigns receive no daemon grant).
- **Runtime behavior:** When `PASEO_AGENT_ID` is set (or Paseo is explicitly selected), FO dispatches workers through the 5 `paseo` CLI verbs (with find-or-create workspace reuse, stdout-only JSON parsing, frozen `spacedock.worker`/`spacedock.epoch` labels, and same-workspace cascade-archive guard) instead of native host subagent tools, and Paseo-launched workers load `paseo-ensign-runtime.md` ahead of inner-host markers.

## Acceptance criteria

**AC-1 (VALUE) - A real Paseo-hosted FO discovers installed Spacedock skills by name, dispatches a fresh Paseo worker via verbatim `--host paseo` output, verifies durable Git commits across code worktree and split-root state checkout, and archives the worker cleanly.**
Verified by: `go test -tags live -count=1 -timeout 15m -run '^TestLivePaseoDispatchSmoke$' ./internal/ensigncycle -v` passes against an isolated real `paseo` daemon (defaulting to `"codex/" + release.LiveModels["codex.exec"]`; `SPACEDOCK_PASEO_BIN`, `SPACEDOCK_PASEO_LIVE_PROVIDER`, and `SPACEDOCK_PASEO_LIVE_THINKING` overridable), staging workflow-free current-checkout skills via `livePluginDir(t)` (local Codex marketplace in `newCodexLiveIsolatedHome` for `codex/*`; fixture `.agents/skills/` for non-Codex overrides), invoking `first-officer` by name, forwarding `dispatch build --host paseo`'s emitted `prompt` verbatim, and asserting stamped non-`main` code worktree with clean `receipt.txt` commit, clean path-scoped `## Stage Report: implementation` commit in `.spacedock-state` with intact frontmatter, exactly 1 worker carrying `spacedock.worker=paseo-smoke-implementation` and `spacedock.epoch=paseo-smoke-implementation-c1-a1` with `ParentAgentId == foAgentId` and `Archived == true` (`Archived == false` on the FO), and clean daemon shutdown. Falsifying edit: breaking `SKILL.md`'s `PASEO_AGENT_ID` selection, breaking any of the 5 `paseo` CLI flags in `paseo-fo-dispatch.md`, omitting `paseo archive`, or dropping `PASEO_AGENT_ID` parent linkage fails the test.

**AC-2 - `spacedock dispatch build --host paseo` emits a host-neutral schema-v2 envelope and rejects `--bare-mode` before writing an artifact.**
Verified by: Go unit tests in `internal/dispatch/build_paseo_host_test.go` (`go test ./internal/dispatch -run 'TestBuildPaseo|TestResolveBuildHost'`) assert `--host paseo` succeeds with `--stamp` and split-root worktree/non-worktree entities, emits `model: null` with ignore-note stderr on declared stage models, disambiguates `dispatch_file_path` when `PASEO_AGENT_ID` is set, includes neutral `ensign` skill-load and `Done:` completion instructions, rejects `--bare-mode` with exit `2` and zero artifact created, and contains zero banned host-specific directives (`Skill(skill=`, `SendMessage`, `$spacedock:ensign`, `/skill:ensign`, `pi-subagents`, `TeamCreate`, `TeamDelete`, `Agent(`, `FO mailbox`). Falsifying edit: removing `"paseo"` from `validBuildHost`, allowing `--bare-mode` on Paseo, or reintroducing any host-specific token in `firstActionBlock`/`completionSignalBlock`/`dispatchPointerPrompt` fails the test.

**AC-3 - First-officer and ensign skills wire the Paseo runtime adapters with valid reference closure and without breaching FO byte caps or instruction-read quarantine.**
Verified by: `go test ./...` (including `internal/contractlint` tests `TestFirstOfficerEntryEagerlyImportsOnlySharedCore`, `TestUserSkillReferenceClosureResolves`, `TestBootResidentDeferredLoadPointsResolve`, `TestFOInstructionComponentCaps`, `TestNoInstructionReadsOutsideQuarantine`, and the new Paseo adapter resolution check in `structural_checks_test.go`) passes cleanly. Falsifying edit: deleting or renaming `paseo-fo-dispatch.md` or `paseo-ensign-runtime.md`, mentioning an ensign `references/*.md` path inside `paseo-fo-dispatch.md`, eagerly `@`-importing `paseo-fo-dispatch.md`, exceeding the 23,500-byte cap on `first-officer-shared-core.md`, or reading a skill file from a test outside `internal/contractlint` fails the suite.

**AC-4 - Optional CI workflow runs `TestLivePaseoDispatchSmoke` using `internal/release/live_models.txt`'s `codex.exec` model without blocking merges, and offline controls verify isolated daemon config, stdout-only JSON parsing, and local skip behavior.**
Verified by: `internal/contractlint/paseo_live_lane_test.go` parses `.github/workflows/paseo-live-e2e.yml` to assert `permissions: contents: read`, `pull_request` (not `pull_request_target`), `needs: offline`, `continue-on-error: true`, `environment: CI-E2E-CODEX`, `CODEX_AUTH_JSON` or `OPENAI_API_KEY` check, `SPACEDOCK_BIN`/`SPACEDOCK_REPO_ROOT` exports, `codex.exec` read from `internal/release/live_models.txt` (no hardcoded model literal), and `-run '^TestLivePaseoDispatchSmoke$'`; and `internal/ensigncycle/paseo_live_controls_test.go` asserts isolated daemon config keys (`daemon.listen`, `daemon.relay.enabled`, `features.webUi.enabled`, `features.dictation.enabled`, `features.voiceMode.enabled`), env scrubbing of ambient `PASEO_*` vars, provider/thinking flag splitting, stdout-only JSON parsing in the presence of exit-0 `Using workspace ...` stderr output, and clean skip when `paseo` or auth is absent without `SPACEDOCK_PASEO_LIVE_REQUIRED`. Falsifying edit: removing `continue-on-error: true` or `permissions: contents: read`, hardcoding a model literal instead of reading `live_models.txt`, switching `runJSON` to `CombinedOutput()`, passing removed `--port`/`--no-relay` flags, or failing instead of skipping when `paseo` is absent locally fails these tests.

## Test plan

- **Primary proof owners & cost:**
  1. **Offline Go suite (deterministic, ~15s, no external deps):** `go test ./...` and `go test ./... -race` covering `internal/dispatch/build_paseo_host_test.go` (AC-2), `internal/contractlint/` including `paseo_live_lane_test.go` (AC-3, AC-4), and `internal/ensigncycle/paseo_live_controls_test.go` (AC-4). No pre-existing Paseo live or dispatch proof owner existed; these are the minimal primary owners, each catching a distinct failure mode without duplicating prose-greps.
  2. **Targeted local live smoke (live, ~2–5 min, 1 FO + 1 worker):**
     - Default Codex subscription lane (`"codex/" + release.LiveModels["codex.exec"]`): `go test -tags live -count=1 -timeout 15m -run '^TestLivePaseoDispatchSmoke$' ./internal/ensigncycle -v` (proves AC-1).
     - Optional local Antigravity override: `SPACEDOCK_PASEO_BIN=<path> SPACEDOCK_PASEO_LIVE_PROVIDER=antigravity/gemini-4-neon SPACEDOCK_PASEO_LIVE_THINKING=high go test -tags live -count=1 -timeout 15m -run '^TestLivePaseoDispatchSmoke$' ./internal/ensigncycle -v` (proves split `--thinking`, fixture `.agents/skills/` staging, and cross-provider host neutrality).
  3. **One-off implementation/validation spot-checks:** non-`/tmp` Safehouse sibling `.spacedock-state` write check (recorded in implementation/validation report), detached adversarial audit on high-stakes surfaces (`spike-paseo-dispatch-backend`), and `gofmt -w ./cmd ./internal` + `spacedock status --workflow-dir docs/dev --validate`.

### Feedback Cycles

- Cycle 1: REJECTED — detached adversarial audit (`c8441547`, `spike-paseo-dispatch-backend`) + Commander live Codex smoke; surface 17/+1102 vs estimate 17/+850 (+29.6%); AC unchanged

## Stage Report: ideation

- DONE: Define exact --host paseo envelope contract (build.go prompt/blocks, paseo-ensign-runtime.md) and exact 5-verb paseo-fo-dispatch.md CLI contract (flags, JSON fields, exit codes, PASEO_AGENT_ID selection, project-ID/Safehouse worktree handling) with zero day-1 ambiguities
  Specified exact `build.go`/`dispatch.go` rules (`--host paseo`, `--bare-mode` exit-2 refusal, `PASEO_AGENT_ID` filename disambiguation, verbatim prompt/blocks, `paseo-ensign-runtime.md`, and `SKILL.md` wiring without touching capped `first-officer-shared-core.md` or naming ensign `references/` inside `paseo-fo-dispatch.md` — AC-2, AC-3) plus all 5 Paseo 0.17.1 CLI verbs with M3 find-or-create (`workspace ls` before `workspace create`), M4 stdout-only JSON parsing (`Using workspace ...` stderr noise), M5 same-workspace cascade-archive guard, and M7 frozen `spacedock.worker`/`spacedock.epoch` (`{worker_base}-c1-a1`) labels.
- DONE: Specify TestLivePaseoDispatchSmoke (isolated daemon setup, canonical temp home, SPACEDOCK_PASEO_BIN and SPACEDOCK_PASEO_LIVE_PROVIDER overrides, durable Git + archive assertions) and optional non-required CI workflow on codex/gpt-6-luna using CODEX_AUTH_JSON
  Specified shared package-private `internal/ensigncycle/paseo_live_helpers_test.go` (`startIsolatedPaseoDaemon`, `(*paseoDaemon).runJSON` with separate stdout/stderr buffers, `setupPaseoSplitRootFixture`, `newCodexLiveIsolatedHome`) and `paseo_live_smoke_test.go` using `livePluginDir(t)` workflow-free skill staging + verbatim `dispatch build --host paseo` prompts (M6, AC-1), plus optional `.github/workflows/paseo-live-e2e.yml` (`permissions: contents: read`, `needs: offline`, `continue-on-error: true`, `CI-E2E-CODEX`, `CODEX_AUTH_JSON || OPENAI_API_KEY`, `SPACEDOCK_BIN`/`SPACEDOCK_REPO_ROOT`, and `codex.exec` sourced from `internal/release/live_models.txt` — M8, AC-4).
- DONE: Align task body (Problem, Proposed approach, Risk evidence citing 2026-11-21 live probes + non-/tmp Safehouse sibling state-checkout check, Expected surface net LOC + files, scanner-exact Acceptance criteria with falsifying Verified by clauses, Test plan)
  Folded all `936` findings from `staff-review.md` (`8447132f3`) into `docs/dev/.spacedock-state/paseo-dispatch-backend/index.md`, updating Problem, Proposed approach, Risk evidence (2026-11-18/21 probes + staff review + dogfood frictions + non-`/tmp` Safehouse check), Expected surface (+850 net, +863/-13 across 17 files), scanner-exact `## Acceptance criteria` (AC-1 through AC-4 with falsifying edits), and `## Test plan`.

### Summary

Folded all Preflight Staff Review findings assigned to `936` (M3–M8 plus Non-blocking #1 and #2 from `staff-review.md`) into `docs/dev/.spacedock-state/paseo-dispatch-backend/index.md`. Verb 1 is now explicit find-or-create, stdout-only JSON parsing and same-workspace cascade-archive guards are codified, `TestLivePaseoDispatchSmoke` stages workflow-free skills via `livePluginDir(t)` with zero hardcoded checkout paths, labels are frozen across `936`/`x9m`, and the shared live-helper seam plus `live_models.txt` SSOT and CI hardening are specified. Ready for ideation gate lock.


## Stage Report: implementation

- DONE: Implement --host paseo in internal/dispatch/{build,dispatch}.go (schema-v2 neutral envelope, --bare-mode exit-2 refusal, PASEO_AGENT_ID file disambiguation, ASD-STE100 help/errors) plus paseo-ensign-runtime.md, paseo-fo-dispatch.md (5 CLI verbs with M3 find-or-create, M4 stdout-only JSON, M5 cascade-archive guard, M7 frozen labels), SKILL.md wiring (keeping first-officer-shared-core.md untouched and no ensign references/* in paseo-fo-dispatch.md), and unit/contractlint tests
  Commit `686ea0da5` (`17 files, +1116/-14 = +1102 net LOC`); `TestBuildPaseoHost*`/`TestResolveBuildHost*` assert neutral schema-v2 envelope, exit-2 bare refusal before artifact write, `PASEO_AGENT_ID` filename disambiguation, non-ambiguity with inner host markers, and zero banned tokens (broken by removing `"paseo"` from `validBuildHost`, permitting `--bare-mode`, or leaking host-specific tokens); `TestBootResidentDeferredLoadPointsResolve`/`TestFOInstructionComponentCaps`/`TestPaseoRuntimeAdaptersResolveOnDisk` prove closure and untouched `first-officer-shared-core.md` (`23,497/23,500` B; broken by deleting an adapter, adding ensign `references/*` to `paseo-fo-dispatch.md`, or growing shared core).
- DONE: Implement shared internal/ensigncycle/paseo_live_helpers_test.go (startIsolatedPaseoDaemon, runJSON with separate stdout/stderr buffers, setupPaseoSplitRootFixture, livePluginDir staging), paseo_live_controls_test.go, paseo_live_smoke_test.go (TestLivePaseoDispatchSmoke), .github/workflows/paseo-live-e2e.yml (codex.exec from live_models.txt, permissions/offline gate/continue-on-error), and docs/runtime-live-ci.md
  Commit `686ea0da5`; `TestPaseoLiveConfigEnvAndRunArgs`/`TestPaseoRunJSONAndPreconditions` assert the 5 daemon config keys, `live_models.txt` default provider, ambient `PASEO_*` scrub, split `--thinking`, exit-0 stderr tolerance in `runPaseoJSON`, and local skip vs required fatal (broken by switching to `CombinedOutput()`, passing removed `--port`/`--no-relay`, or failing locally without `paseo`); `TestPaseoLiveWorkflowSecurityAndModelSSOT` + `TestNode24ActionsPinnedAtMinimum` assert `permissions: contents: read`, `pull_request` (no `pull_request_target`), `needs: offline`, `continue-on-error: true`, `@v5` upload-artifact, and zero hardcoded model literals (broken by dropping any guard or hardcoding `codex.exec`); `go test -tags live -c ./internal/ensigncycle` compiles cleanly.
- DONE: Run Day-1 non-/tmp Safehouse sibling .spacedock-state write check, go test ./..., go test ./... -race, gofmt -w ./cmd ./internal, and spacedock status --workflow-dir docs/dev --validate; record live smoke run or exact sandbox unavailability in Stage Report
  Verified Day-1 non-/tmp `.safehouse` workaround (read/write/commit in both worktree and sibling `.spacedock-state`), `gofmt` clean on all changed Go files, full `go test ./...` + `go test ./... -race` green across all 24 packages (plus `-race` re-run on all 4 touched packages after final test LOC trim), `go test -tags live -c` + `go vet -tags live` green, and `spacedock status --workflow-dir docs/dev --validate` returned `VALID`.

### Summary

Implemented the opt-in `--host paseo` dispatch backend, Paseo ensign and 5-verb CLI-first FO adapters, shared isolated-daemon live helpers (`paseo_live_helpers_test.go`), `TestLivePaseoDispatchSmoke`, optional `.github/workflows/paseo-live-e2e.yml`, and `docs/runtime-live-ci.md` in commit `686ea0da5` on branch `paseo-dispatch-backend` (17 files, `+1,102` net LOC within `+850 ± 30%`; `.safehouse` kept untracked/gitignored; zero hardcoded model literals; `first-officer-shared-core.md` untouched). Day-1 Safehouse check: default `paseo-safehouse` without a worktree `.safehouse` denies parent `.git` and sibling `docs/dev/.spacedock-state` on non-`/tmp` repos (`F7`, proven in `c1-a1`), while the gitignored mode-600 `.safehouse` workaround installed in this worktree allowed `git status`, worktree code commits, and path-scoped `.spacedock-state` read/write/commits cleanly. `TestLivePaseoDispatchSmoke` was not executed inside this Safehouse ensign due to named sandbox unavailability (`F8`: ordinary Safehouse ensigns lack daemon access — `paseo status --json` fails with `DAEMON_STATE_READ_FAILED` / `EPERM: operation not permitted` — and unrestricted Paseo/Codex live runtime privileges); live-tag compilation (`go test -tags live -c ./internal/ensigncycle`) is verified green for Commander to execute `TestLivePaseoDispatchSmoke`.


## Stage Report: validation

- DONE: Verify AC-1 through AC-4 against commit 686ea0da5 (surface +1102 net LOC / 17 files, first-officer-shared-core.md untouched, no ensign references/* in paseo-fo-dispatch.md, live_models.txt SSOT, livePluginDir staging, ASD-STE100 compliance in dispatch.go, build.go, and docs/runtime-live-ci.md)
  Commit `686ea0da5` (`17 files, +1116/-14 = +1102 net LOC`, `+29.65%` vs `+850 ± 30%` and `13–20` files); `first-officer-shared-core.md` untouched (`23,497/23,500` B); `internal/runtimehost` untouched; zero `references/*` in `paseo-fo-dispatch.md`; zero hardcoded `gpt-6-luna` across all 17 files (`release.LiveModels["codex.exec"]` + workflow `awk` over `live_models.txt`); `livePluginDir(t)` staging in `setupIsolatedPaseoCodexHome` and `stageNonCodexSkillsIfNeeded` (zero hardcoded `<checkout>/skills/...`); ASD-STE100 compliant across new prose/comments in `docs/runtime-live-ci.md`, `dispatch.go`, and `build.go`.
- DONE: Run offline verification (go test ./..., go test ./... -race, go test -tags live -c ./internal/ensigncycle, gofmt check, spacedock status --workflow-dir docs/dev --validate) and perform the semantic adversarial pass over changed behavior
  `go test ./...` and `go test ./... -race` green across all 24 packages; `go test -tags live -c -o /tmp/ensigncycle-val-live.test ./internal/ensigncycle` and `go vet -tags live` green; `gofmt -l` clean on all 11 changed Go files (`internal/release/runtime_live_evidence_workflow_test.go` is pre-existing untouched from `e68aa5b7a`); `spacedock status --workflow-dir docs/dev --validate` returned `VALID`; exercised live `paseo 0.17.1` flag parsing (`daemon config set` all 5 keys, `daemon start --json` vs bogus-flag control, all 5 CLI verbs `--help`) and compiled `dispatch build --host paseo` edge variants (`--advance`, >36-char and all-punctuation `PASEO_AGENT_ID`, `--bare-mode`, `--stamp --bare-mode`).
- DONE: Record PASSED or REJECTED recommendation with per-AC evidence, falsifying edits, Day-1 Safehouse finding summary, and named sandbox unavailability for live smoke (Commander runs TestLivePaseoDispatchSmoke separately)
  All four ACs accounted below with falsifying edits, `F7` Day-1 Safehouse summary, `F8` named sandbox unavailability, and `## Review-finding disposition` classification (`0 Material`, `1 Deferred risk`, `2 Polish`, `0 Needs decision`); `Recommendation: PASSED` (live AC-1 gate evidence supplied by Commander's concurrent `TestLivePaseoDispatchSmoke` run and detached adversarial audit on `spike-paseo-dispatch-backend`).

### Per-AC Evidence & Falsifying Edits

- **AC-1 (VALUE):** Harness (`paseo_live_helpers_test.go`, `paseo_live_smoke_test.go`) stages workflow-free current-checkout skills via `livePluginDir(t)`, isolates `PASEO_HOME`/`CODEX_HOME`, parses stdout-only JSON, and asserts durable non-`main` `receipt.txt` commit, path-scoped `## Stage Report: implementation` commit, 1 archived worker with `ParentAgentId == foAgentId` and `spacedock.epoch=paseo-smoke-implementation-c1-a1`, and unarchived FO (`go test -tags live -c` + `paseo 0.17.1` CLI flags verified). **Named sandbox unavailability (`F8`):** live execution unavailable in this ordinary Safehouse ensign (`paseo status --json` exits `1` `DAEMON_STATE_READ_FAILED` / `EPERM` on `paseo.pid`; no daemon/live-provider grant); Commander runs `TestLivePaseoDispatchSmoke` separately. Falsifying edit: breaking `SKILL.md` `PASEO_AGENT_ID` selection, any of the 5 `paseo` flags, `paseo archive`, or `ParentAgentId` linkage fails `TestLivePaseoDispatchSmoke`.
- **AC-2:** `TestBuildPaseoHostPromptShapeAndSplitRoot`, `TestBuildPaseoHostWithStampAndBareRejection`, `TestBuildPaseoHostModelDisambiguationAndResolve`, `TestResolveBuildHost*`, and `TestBuildHostResolutionFromFlagAndEnv` pass (normal + `-race`), plus compiled CLI probes on `--advance`, >36-char/all-punctuation `PASEO_AGENT_ID`, and `--bare-mode` exit `2` with zero artifact. Falsifying edit: removing `"paseo"` from `validBuildHost`, allowing `--bare-mode`, or reintroducing any banned host token fails `internal/dispatch`.
- **AC-3:** `TestFirstOfficerEntryEagerlyImportsOnlySharedCore`, `TestUserSkillReferenceClosureResolves`, `TestBootResidentDeferredLoadPointsResolve`, `TestFOInstructionComponentCaps` (`23,497/23,500` B untouched), `TestNoInstructionReadsOutsideQuarantine`, and `TestPaseoRuntimeAdaptersResolveOnDisk` pass (normal + `-race`); `PASEO_AGENT_ID` precedes inner markers in `skills/ensign/SKILL.md`. Falsifying edit: deleting an adapter, naming an ensign `references/*` path in `paseo-fo-dispatch.md`, `@`-importing it, or exceeding the shared-core cap fails `internal/contractlint`.
- **AC-4:** `TestPaseoLiveWorkflowSecurityAndModelSSOT`, `TestNode24ActionsPinnedAtMinimum`, `TestPaseoLiveConfigEnvAndRunArgs`, and `TestPaseoRunJSONAndPreconditions` pass (normal + `-race`). Falsifying edit: dropping `continue-on-error: true` or `permissions: contents: read`, hardcoding `gpt-6-luna`, switching `runPaseoJSON` to `CombinedOutput()`, or failing instead of skipping locally without `paseo` fails these suites.

### Review-Finding Disposition & Safehouse Notes

- **Material:** None.
- **Deferred risk (1 — outcome defect, pre-existing `runBuild` ordering shared with `codex`):** `dispatch build --host paseo --stamp --bare-mode` executes `runStamp` (`internal/dispatch/build.go:123-127`) before `runBuildFields` rejects `--bare-mode` (`build.go:246-248`). Evidence: (1) normal Paseo FO omits `--bare-mode` per `paseo-fo-dispatch.md`; (2) observable harm only if a caller erroneously combines `--stamp` and `--bare-mode`, leaving `started:`/`worktree:` stamped and `.worktrees/` created before exit `2` (no artifact written); (3) `none:` outside promised Paseo/Codex workflow and AC-2 (`--stamp` and `--bare-mode` tested per spec); (4) trigger reproduced in throwaway fixture. Promote-to-material condition: any adapter/caller combining `--stamp --bare-mode` on `paseo`/`codex`, or when `runBuild` hoists host/bare-mode validation ahead of `runStamp`.
- **Polish (2):** (a) Pre-existing semicolons in `--bare-mode` help (`dispatch.go:300`) and error (`build.go:247`, matching sibling `codex` and approved ideation diff verbatim) vs ASD-STE100 Rule 8.1; (b) `paseo_live_helpers_test.go:59` reuses `decideCodexLiveAuthForCI`, whose fatal string names `SPACEDOCK_CODEX_LIVE_REQUIRED` instead of `SPACEDOCK_PASEO_LIVE_REQUIRED` when required CI auth is missing.
- **Needs decision:** None.
- **Day-1 Safehouse (`F7`) & Sandbox (`F8`):** Confirmed worktree remains clean (`686ea0da5`) with untracked gitignored mode-600 `.safehouse` enabling parent `.git` and sibling `.spacedock-state` access on this non-`/tmp` repo (`F7`); live smoke remains sandbox-blocked here (`F8`) and delegates live AC-1 execution to Commander.

### Summary

Independent validation of commit `686ea0da5` on `paseo-dispatch-backend` verified all offline ACs (AC-2, AC-3, AC-4), surface tolerance (`+1,102` net LOC / `17` files), `first-officer-shared-core.md` preservation, `live_models.txt` SSOT, `livePluginDir(t)` staging, ASD-STE100 compliance, full `go test ./...` and `go test ./... -race` across all 24 packages, `go test -tags live -c`, `gofmt`, `spacedock status --validate`, Paseo 0.17.1 CLI flag compatibility, and the semantic adversarial pass with zero Material findings (1 Deferred risk, 2 Polish). Live execution of `TestLivePaseoDispatchSmoke` (AC-1) is unavailable in this Safehouse ensign (`F8`, `EPERM` on `paseo.pid`) and is being executed concurrently by Commander alongside the detached audit on `spike-paseo-dispatch-backend`.

Recommendation: PASSED


## Stage Report: implementation (cycle 2)

- DONE: Fix Material Finding 1 (headless codex --mode full-access in paseoRunArgs/TestLivePaseoDispatchSmoke, daemon.log artifact capture, idle message logging + worker Cwd assertion)
  Commit `e1b971b32`: `paseoRunArgs` appends `--mode full-access` for `codex/*` and `TestLivePaseoDispatchSmoke` includes `--mode full-access` in the FO worker-spawn prompt for `codex/*`; `paseo_live_helpers_test.go` `t.Cleanup` copies `daemon.log` (falling back to `paseo.log`) to `d.ArtifactDir` on failure; `waitForPaseoAgentIdle` logs `res.Message` via `t.Logf` on `idle` and `TestLivePaseoDispatchSmoke` asserts `wi.Cwd == codeWorktree` on `paseo inspect`. `TestPaseoLiveConfigEnvAndRunArgs` asserts exact `paseoRunArgs` flags (`--no-wait`, `--workspace`, `--provider`, `--thinking`, `--mode`, `--title`, `--label`, `--env`, `--json`) and breaks if `--mode full-access` is omitted for `codex/*` or added for non-Codex providers.
- DONE: Fix Material Finding 2 (paseo_live_lane_test.go exact permissions len==1, jobs.paseo-live.env secret bindings, SPACEDOCK_PASEO_LIVE_REQUIRED=1 assertion) and Finding 3 (paseo_live_helpers_test.go pre-start t.Cleanup + fail-closed status check, SPACEDOCK_REPO_ROOT scrub, test temp cleanup)
  Commit `e1b971b32`: `TestPaseoLiveWorkflowSecurityAndModelSSOT` asserts `len(doc.Permissions) == 1 && doc.Permissions["contents"] == "read"`, `jobs.paseo-live.env` bindings for `CODEX_AUTH_JSON` and `OPENAI_API_KEY`, and `SPACEDOCK_PASEO_LIVE_REQUIRED=1` + `needs.offline.outputs.codex_exec` inside the `TestLivePaseoDispatchSmoke` run step (fails if an extra/write permission is added, either secret env binding is dropped, or `SPACEDOCK_PASEO_LIVE_REQUIRED=1` is removed from the smoke step); `startIsolatedPaseoDaemon` registers `t.Cleanup` before `daemon start` and fails closed via `t.Errorf` on `status --json` errors (reporting `stopErr`) or non-`stopped` state; `paseoLiveEnv` scrubs `SPACEDOCK_REPO_ROOT` (asserted in `TestPaseoLiveConfigEnvAndRunArgs`, broken if removed from `cleanEnviron`) and `TestBuildPaseoHostModelDisambiguationAndResolve` registers `t.Cleanup` removing `/tmp/spacedock-dispatch/agent-*` (verified zero leftover `agent-*` files after test run).
- DONE: Keep total net LOC vs main <= +1105 across the same 17 files; run gofmt, go test ./..., go test ./... -race, go test -tags live -c, and spacedock status --validate; commit on paseo-dispatch-backend and append updated implementation Stage Report
  Commit `e1b971b32` on `paseo-dispatch-backend`: `git diff --shortstat main..HEAD` is `17 files changed, 1117 insertions(+), 14 deletions(-)` (`+1,103` net LOC <= `+1,105` cap; `.safehouse` untracked/gitignored; `first-officer-shared-core.md` untouched at `23,497/23,500` B); `gofmt` clean on all changed Go files, `go test ./...` and `go test ./... -race` green across all 24 packages, `go test -tags live -c` + `go vet -tags live ./internal/ensigncycle` green, and `spacedock status --workflow-dir docs/dev --validate` returned `VALID`.

### Summary

Applied only the three FO-authorized Cycle 1 Material fixes in commit `e1b971b32` on `paseo-dispatch-backend` (`17 files, +1,103 net LOC`, within the `+1,105` cap): headless Codex `--mode full-access`, `daemon.log` capture, idle message logging, and worker `Cwd` verification (Finding 1); fail-closed CI workflow permissions, secret env bindings, and `SPACEDOCK_PASEO_LIVE_REQUIRED=1` step checks (Finding 2); and pre-start `t.Cleanup` registration, fail-closed stop/status error reporting, `SPACEDOCK_REPO_ROOT` scrubbing, and `/tmp/spacedock-dispatch/agent-*` cleanup (Finding 3). All offline, race, live-compile, vet, formatting, and workflow validation checks pass; live execution of `TestLivePaseoDispatchSmoke` remains delegated outside this Safehouse ensign (`F8`).
