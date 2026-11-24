# First Officer Dispatch Module (Paseo)

The Paseo dispatch parts (deferred from `fo-dispatch-core.md`), read alongside `fo-dispatch-core.md` at the dispatch/resume boundary (before resume reconciliation, `«dispatch.next-action»()`, worker dispatch, or terminal teardown) when Paseo is active (reachable `paseo` daemon with `PASEO_AGENT_ID` set, or explicit captain selection). Any Paseo-hosted first officer (regardless of inner harness) drives fresh workers through five `paseo` CLI verbs and `spacedock dispatch build --host paseo`.

**Zero sidecar/lease files, zero entity-body worker writes:** Reconcile worker state solely from the three pre-existing sources of truth — (1) entity frontmatter + `.spacedock-state` Git history (`status`, `started`, `worktree`, `### Feedback Cycles`, and stage-entry baseline commit SHA), (2) code worktree Git state (expected task branch, clean status), and (3) Paseo daemon labels (`spacedock.worker`, `spacedock.epoch`) + `paseo inspect`. Never write `## Paseo Worker` or `### Dispatch Retries` into entity markdown, and never create sidecar or lease files.

## Common CLI Invariants (Paseo 0.17.1)

- **Always pass `--json`.** Successful commands exit `0` with JSON on stdout; CLI/daemon errors exit `1` with `{"error":{"code":"...","message":"...","details":"..."}}` on stderr.
- **Stdout-only JSON parsing:** `paseo run --workspace <id> --json` writes unconditional non-JSON `Using workspace <id>` (and workspace-creation tips) to stderr even on exit `0`. Callers MUST parse stdout only — NEVER combine streams (`2>&1 | jq` fails JSON decoding on exit `0`). Capture stderr separately (or redirect stderr to a temp file for failure diagnostics) and decode JSON from stdout alone.
- **Preserve `PASEO_AGENT_ID`:** Keep `PASEO_AGENT_ID` in the environment when calling `paseo` CLI so Paseo auto-stamps `paseo.parent-agent-id` and `paseo.parent-server-id` on spawned children. Do not pass `--home` or `--host` during normal FO runs (inherit the session's daemon target).
- **Least-privilege ensign sandbox:** Never pass `paseo-safehouse.paseo-access=daemon` to ordinary ensigns (reserved for FOs until Paseo ships a parent-only daemon ACL).
- **Provider / model / thinking split:** Always pass an explicit `--provider <provider>/<model>` (never rely on `build.model`, which is `null`), and pass thinking effort separately via `--thinking <id>` when requested (for example `--provider antigravity/gemini-4-neon --thinking high` or `--provider codex/<model>`).
- **Shared dispatch directory:** If workers use a private `/tmp`, pass `--dispatch-dir DIR` or set `SPACEDOCK_DISPATCH_DIR`. Use a directory that both processes can read. The flag takes precedence. The default is `/tmp/spacedock-dispatch`.
- **Workspace location:** Paseo creates worktrees under `~/.paseo/worktrees/` with `--isolation worktree`. Spacedock owns the task worktree. Adopt `code_cwd` with `--isolation local`.
- **Thinking setting:** After spawn, run `paseo inspect "$agent_id" --json`. Verify `Thinking` against the requested setting. The `paseo models` metadata alone does not prove the worker setting.
- **Frozen shared label contract:**
  - `spacedock.worker` = `{worker_base}` (= `envelope.name` emitted by `spacedock dispatch build --host paseo` or `spacedock dispatch name`). Invariant across cycles and retries (always <= 56 chars); NEVER append `-retry` or `-cycleN` to `spacedock.worker`.
  - `spacedock.epoch` = `{worker_base}-c{cycle}-a{attempt}`, where `cycle = 1 + count(### Feedback Cycles entries for stage S)` and `attempt` is `1` initially (`a1`) or `2` on the single allowed retry (`a2`).

## Pre-Spawn Worktree, Project-ID, and Safehouse Handling

1. **Build envelope & compute epoch:** Compute `cycle = 1 + count(### Feedback Cycles entries for stage S)` and `attempt` (`1` initially, `2` on single retry). Run `${SPACEDOCK_BIN:-spacedock} dispatch build --host paseo [--stamp] --workflow-dir {workflow_dir} --entity-path {entity_path} --stage {stage} --checklist-file -` (include `--stamp` on initial entry; omit `--stamp` on retry to reuse the existing worktree; always omit `--bare-mode`). Retain `name` (`{worker_base}`), `description`, `dispatch_file_path`, and `prompt`. Set `epoch="${name}-c${cycle}-a${attempt}"`.
2. **Resolve physical `code_cwd` from structured status** (never parse the dispatch body):
   - Run `${SPACEDOCK_BIN:-spacedock} status --workflow-dir {workflow_dir} --read {entity_path} --json` and read `frontmatter.worktree`.
   - If non-empty: join relative paths against `git -C {workflow_dir} rev-parse --show-toplevel`, keep absolute paths as-is, and canonicalize symlinks (`realpath` / `pwd -P`) to produce `code_cwd`.
   - If empty (non-worktree stage such as `ideation`): use canonical `git -C {workflow_dir} rev-parse --show-toplevel`.
3. **Spacedock owns worktrees:** Never pass `--isolation worktree` or `--new-workspace worktree` to Paseo (Git exits 128 on an occupied branch or creates a mismatched branch). Always adopt the existing physical `code_cwd` via find-or-create (`workspace ls` then `workspace create --isolation local`). Paseo classifies an adopted Git worktree as `isolation: "worktree"` while keeping exact `cwd`.
4. **Project ID:** When `workspace ls` finds no matching workspace, resolve `project_id` by matching canonical `git -C {workflow_dir} rev-parse --show-toplevel` against `path` in `paseo project ls --json` (`[{projectId, name, kind, path}]`). When matched, pass `--project "$project_id"` to `workspace create`; when unmatched (for example a fresh isolated repo), omit `--project` so Paseo registers the checkout automatically.
5. **Safehouse sibling state checkout:** In split-root worktree stages under `paseo-safehouse`, the worker writes both `code_cwd` (`.worktrees/spacedock-ensign-{slug}`) and the sibling state checkout (`{workflow_dir}/.spacedock-state`). If Safehouse blocks the state-checkout write on a non-`/tmp` repo, the worker reports the exact diagnostic and stops rather than retrying silently.

## Five Paseo CLI Verbs

### Verb 1 — Find-or-create workspace

`paseo workspace create` is not idempotent (unconditional calls mint duplicate workspaces for the same directory). First list active workspaces:

```bash
paseo workspace ls --json
```

Match an active `WorkspaceRow` whose canonical `cwd` equals `$code_cwd`. Only when none matches, create one:

```bash
paseo workspace create --isolation local --path "$code_cwd" [--project "$project_id"] --title "$description" --json
```

- **JSON stdout fields:** `workspace ls` returns `WorkspaceRow[]`; `workspace create` returns `WorkspaceRow` (`workspaceId` (`wks_...`), `project`, `name`, `isolation` (`"local"|"worktree"`), `cwd` (unshortened), `background`).
- **Exit codes & branching:**
  - `0`: verify matched/returned `cwd == $code_cwd`, retain `workspaceId`.
  - `1`: `DAEMON_NOT_RUNNING` or `WORKSPACE_CREATE_FAILED` — abort spawn before creating an agent.

### Verb 2 — Spawn worker (`run --no-wait`)

```bash
paseo run --no-wait \
  --workspace "$workspace_id" \
  --provider "$provider_model" \
  [--thinking "$thinking_id"] \
  [--mode "$mode_id"] \
  --title "$description" \
  --label "role=ensign" \
  --label "spacedock.worker=$name" \
  --label "spacedock.epoch=$epoch" \
  --json "$prompt"
```

- **JSON stdout fields:** `agentId` (UUID), `status` (`"created"|"running"`), `provider`, `cwd`, `title`.
- **Exit codes & branching:**
  - `0`: **parse stdout only** (stderr prints `Using workspace <id>` on exit `0`). Retain `agentId`, `workspaceId`, `spacedock.worker=$name`, `spacedock.epoch=$epoch`, `cwd`, and requested/returned provider, model, thinking, and mode in session memory before waiting. CLI `--no-wait` leaves `notifyOnFinish=false`.
  - `1`: `MISSING_PROVIDER`, `INVALID_OPTIONS`, `INVALID_THINKING_OPTION`, `WORKSPACE_NOT_FOUND`, `AGENT_CREATE_FAILED` — run `## Roster Reconciliation` for the current in-turn attempt `aK` (`spacedock.epoch={worker_base}-c{cycle}-a{K}`) to distinguish a lost create receipt (`aK` present in the inspect-filtered roster -> `## Recovery Matrix`) from a confirmed `spawn-error` (`aK` absent, even if an earlier attempt such as archived `a1` is present -> Row 8 directly before ordinary roster recovery).

### Verb 3 — Bounded wait (`wait`)

```bash
paseo wait "$agent_id" --timeout 300 --json
```

- **JSON stdout fields:** `agentId`, `status` (`"idle"|"timeout"|"permission"|"error"`), `message` (includes last 5 activity items on `idle`/`timeout`).
- **Exit codes & branching:**
  - `0` for ALL FOUR statuses (`timeout`, `permission`, `idle`, `error`) and `1` (`AGENT_NOT_FOUND`, `INVALID_TIMEOUT`, `WAIT_FAILED`): inspect the worker via Verb 4, evaluate `HasValidNewReport` when ended, and route EVERY outcome through `## Recovery Matrix` (the sole decision owner).

### Verb 4 — Inspect & roster query (`inspect`, `ls`)

```bash
paseo inspect "$agent_id" --json
paseo ls --global [--all] --label "spacedock.worker=$name" --json
```

- **JSON stdout fields:**
  - `inspect` (PascalCase): `Id`, `Name`, `Provider`, `Model`, `Thinking`, `Status`, `Archived`, `ArchivedAt`, `Mode`, `Cwd` (unshortened), `CreatedAt` (ISO-8601), `UpdatedAt`, `LastUsage` (`{InputTokens, OutputTokens, CachedTokens, CostUsd}|null`), `PendingPermissions` (`[{id, tool}]`), `Worktree`, `ParentAgentId`.
  - `ls` (camelCase array): `[{id, shortId, name, provider, thinking, status, cwd, created, background, labels, workspaceId}]`. Note: `ls.cwd` is `~`-shortened, `ls.created` is relative time, and `ls` omits `archived`; use `inspect` for exact path, ISO-8601 timestamp, and `Archived` fields.
- **Exit codes:** `0` for snapshot; `1` for `AGENT_NOT_FOUND`, `INSPECT_FAILED`, `LIST_AGENTS_FAILED`.

### Verb 5 — Archive worker (`archive`)

```bash
paseo archive "$agent_id" --json
```

Pass `--force` ONLY on explicit supersede of a running worker.

- **JSON stdout fields:** `agentId`, `status` (`"archived"`), `archivedAt`.
- **Exit codes & branching:** `0` = worker archived; `1` = `AGENT_RUNNING` (fail-closed without `--force`), `AGENT_ALREADY_ARCHIVED`, `AGENT_NOT_FOUND`, `ARCHIVE_FAILED`.
- **Cascade-archive guard:** Archiving a parent agent cascade-kills and archives every child in the same workspace (`parent.workspaceId === child.workspaceId` — all non-worktree stages share the repo-root workspace) and detaches cross-workspace children by deleting `paseo.parent-agent-id`. NEVER archive `$PASEO_AGENT_ID`, and NEVER archive an FO while any child worker is running (use `paseo stop` to interrupt an FO without killing or detaching its children).

## Roster Reconciliation

Run before spawning or advancing an entity in active stage `S = entity.status` (on boot, compaction resume, `«dispatch.next-action»()`, or after a failed/dropped `paseo run --no-wait`):

1. **Derive `(worker_base, cycle, expected_cwd, baseline_rev)` from entity + `.spacedock-state` Git log:**
   - Read entity frontmatter (`status`, `started`, `worktree`) and `### Feedback Cycles` (read-only).
   - Run read-only `${SPACEDOCK_BIN:-spacedock} dispatch name --workflow-dir {workflow_dir} --entity-path {entity_path} --stage {S}` to get invariant `worker_base`.
   - Compute `cycle = 1 + count(### Feedback Cycles entries for S)` and canonical physical `expected_cwd` (`realpath` / `pwd -P` of `worktree` if non-empty, else canonical repo root).
   - Derive `baseline_rev` from `.spacedock-state` Git log on `{entity_rel}`: the latest `dispatch: {slug} entering {S}` (or backward-route / stage-entry) commit SHA for `(S, cycle)`.
2. **Query and inspect-filter Paseo daemon rows:**
   - Run `paseo ls --global --all --label "spacedock.worker={worker_base}" --json` and keep rows whose `labels["spacedock.epoch"]` matches `{worker_base}-c{cycle}-a*` (retaining each row's exact `spacedock.epoch`, including archived earlier attempts such as `a1`).
   - Call `paseo inspect <id> --json` on each matched row (NEVER filter on `ls[].cwd`, `ls[].created`, or `ls[].status` for archival) and keep ONLY rows where `inspect.Cwd == expected_cwd` AND (`frontmatter.started` is empty OR ISO-8601 `inspect.CreatedAt >= frontmatter.started`).
3. **Check failed-create attempt `aK`, then partition and dispatch into Recovery Matrix:**
   - **After a failed `paseo run --no-wait` for current in-turn attempt `aK` (`epoch = {worker_base}-c{cycle}-a{K}`):** check the inspect-filtered roster for `aK` BEFORE ordinary roster recovery:
     - If NO inspect-filtered row has `labels["spacedock.epoch"] == "{worker_base}-c{cycle}-a{K}"` (even when earlier attempts such as archived `a1` are present) -> route directly to **Row 8 (`spawn-error`)** with current in-turn attempt `aK` (do NOT fall through to ordinary roster recovery / Row 5).
     - If an inspect-filtered row has `labels["spacedock.epoch"] == "{worker_base}-c{cycle}-a{K}"` -> lost create receipt; proceed to partition and adopt/recover through `## Recovery Matrix` below.
   - **Ordinary partition and dispatch:** partition filtered rows into `running` (`!inspect.Archived && inspect.Status == "running"`) and `ended` (`inspect.Archived || inspect.Status != "running"`), and note whether any filtered row has `labels["spacedock.epoch"] == "{worker_base}-c{cycle}-a2"` (retry spent).
   - If 0 filtered rows exist on boot/engage: advance if `HasValidNewReport` already holds, else `SpawnFresh` (`a1`, with `--stamp` if not yet stamped).
   - Otherwise route the partitioned roster directly through `## Recovery Matrix` below.

## Recovery Matrix

Sole decision owner for all `wait`, `inspect`, `archive`, lost-create-receipt, and resume outcomes. Every recovery decision is strictly isolated per entity (one held or failing entity never blocks unrelated entities). Obey the Verb 5 cascade-archive guard on every archival (`$PASEO_AGENT_ID` is never archived, and an FO is never archived while any child runs).

**`HasValidNewReport(entity, stage, baseline_rev)`** — all three Git predicates must hold:
1. `git -C {state_checkout} rev-list {baseline_rev}..HEAD -- {entity_rel}` is non-empty and the committed entity carries a new `## Stage Report: {stage}` passing `${SPACEDOCK_BIN:-spacedock} status --read <entity> --stage {stage} --checklist --json` with non-empty evidence, non-empty `### Summary`, and no `FAILED` items except feedback-review findings.
2. `git -C {state_checkout} status --porcelain -- {entity_rel}` is empty.
3. For worktree stages (`worktree != ""`), the code worktree is on the expected task branch with clean `git -C {worktree} status --porcelain`.

| Row | Observed State | Decision & Action |
|---|---|---|
| 1 | `wait` `status=timeout` + `inspect` `Status=running` & `PendingPermissions=[]` | `ContinueWait`: re-issue `paseo wait <id> --timeout 300 --json` (0 retries). |
| 2 | `wait` `status=permission` OR `inspect.PendingPermissions != []` (precedes Row 3 even when `Status=running`) | `HoldEntity`: keep worker live, stop this entity & surface to captain (0 retries). |
| 3 | 1 running worker (on resume or dropped `run --no-wait` receipt) with `PendingPermissions=[]` | `AdoptRunning`: enter `paseo wait <id> --timeout 300 --json` loop (0 duplicate spawns, 0 retries). |
| 4 | 0 running + `HasValidNewReport == true` | `VerifyAndArchive`: idempotent `paseo archive <id> --json` for any unarchived matched worker (treat `AGENT_ALREADY_ARCHIVED` and `AGENT_NOT_FOUND` as success) -> advance / `gate prepare` / `merge guard`. |
| 5 | 0 running + `HasValidNewReport == false` + `a1` only (no `a2` row; ordinary roster recovery — not a failed in-turn `a2` create) | `ArchiveAndRetry`: idempotent `paseo archive <id> --json` for unarchived `a1` -> reuse workspace (Verb 1) & spawn ONE replacement `a2` (`spacedock.worker={worker_base}`, `spacedock.epoch={worker_base}-c{cycle}-a2`, omit `--stamp`, append WIP check note); if that `a2` create exits `1` without creating `a2`, step 3 routes to Row 8 (`aK=a2`), never back to Row 5. |
| 6 | 0 running + `HasValidNewReport == false` + `a2` row exists | `RetryExhausted`: idempotent `paseo archive <id> --json` for unarchived `a2` -> stop this entity & surface to captain (unrelated entities unaffected). |
| 7 | `> 1` running workers | `AmbiguousRoster`: stop this entity & surface to captain; never spawn or kill unowned agents. |
| 8 | `spawn-error` (`paseo run --no-wait` for current in-turn attempt `aK` exits `1` and Roster Reconciliation confirms no inspect-filtered row with `labels["spacedock.epoch"] == "{worker_base}-c{cycle}-a{K}"`, regardless of earlier attempts such as archived `a1`) | If `aK == a1`, retry ONCE in-turn as `a2` (`spacedock.epoch={worker_base}-c{cycle}-a2`, reusing workspace via Verb 1, omitting `--stamp`); if `aK == a2`, stop this entity & surface to captain (NO third create). |

Late or duplicate completion events for an already-advanced entity (`entity.status != S`) or an already-archived worker with no pending advance are a strict no-op.

## Runtime implementation

- `«worker.spawn»` -> Pre-Spawn (`dispatch build --host paseo [--stamp]`, compute `epoch="${name}-c${cycle}-a${attempt}"`, resolve canonical `code_cwd`) -> Verb 1 find-or-create -> Verb 2 (`run --no-wait`, stdout-only JSON parse, `spacedock.worker=$name`, `spacedock.epoch=$epoch`). Forward the emitted pointer verbatim (optional ordinary FO scope notes appended after a newline), obeying the opening state-invariant above.
- `«addressable-worker»` -> ABSENT in v1 (fresh one-shot workers only; worker reuse and child-to-parent messaging are out of scope).
- `«async-dispatch»` -> ASYNC (`paseo run --no-wait` returns `agentId`; bounded polling via Verb 3 `paseo wait --timeout 300 --json` routed through `## Recovery Matrix`).
- `«worker-identity»` -> Paseo daemon labels (`spacedock.worker={worker_base}`, `spacedock.epoch={worker_base}-c{cycle}-a{attempt}`) + `paseo inspect` (`Id`, `Cwd`, `Provider`/`Model`, `Thinking`, `Mode`, `Status`, `Archived`) + entity/Git state; canonical model space is `<provider>/<model>` (+ optional `--thinking <id>`).
- `«completion-signal»` -> Verb 3 `status == "idle"` wakes verification; `HasValidNewReport` in `## Recovery Matrix` (committed worktree deliverable + path-scoped post-`baseline_rev` `## Stage Report: {stage}` commit + clean `git status --porcelain`) is the completion gate.
- `«worker.shutdown»` -> Verb 5 idempotent `paseo archive "$agent_id" --json` (treat `AGENT_ALREADY_ARCHIVED`/`AGENT_NOT_FOUND` as success; `--force` on supersede only), obeying the cascade-archive guard (`agentId != $PASEO_AGENT_ID`, never archive an FO with running children).
- `«context-budget»` -> ABSENT (`inspect.LastUsage` omits context-window max and v1 dispatches fresh only).
- `«roster-reconcile»` -> `## Roster Reconciliation` + `## Recovery Matrix` (reconcile per the opening three-source-of-truth invariant above).
