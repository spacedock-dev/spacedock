# First Officer Dispatch Module (Paseo)

The Paseo dispatch parts (deferred from `fo-dispatch-core.md`), read alongside `fo-dispatch-core.md` at the first worker dispatch when `PASEO_AGENT_ID` is set or the captain selects Paseo dispatch. Any Paseo-hosted first officer (Claude, Codex, Pi, or Antigravity) drives fresh workers through five `paseo` CLI verbs and `spacedock dispatch build --host paseo`.

## Common CLI Invariants (Paseo 0.17.1)

- **Always pass `--json`.** Successful commands exit `0` with JSON on stdout; CLI/daemon errors exit `1` with `{"error":{"code":"...","message":"...","details":"..."}}` on stderr.
- **Stdout-only JSON parsing:** `paseo run --workspace <id> --json` writes unconditional non-JSON `Using workspace <id>` (and workspace-creation tips) to stderr even on exit `0`. Callers MUST parse stdout only — NEVER combine streams (`2>&1 | jq` fails JSON decoding on exit `0`). Capture stderr separately (or redirect stderr to a temp file for failure diagnostics) and decode JSON from stdout alone.
- **Preserve `PASEO_AGENT_ID`:** Keep `PASEO_AGENT_ID` in the environment when calling `paseo` CLI so Paseo auto-stamps `paseo.parent-agent-id` and `paseo.parent-server-id` on spawned children. Do not pass `--home` or `--host` during normal FO runs (inherit the session's daemon target).
- **Least-privilege ensign sandbox:** Never pass `paseo-safehouse.paseo-access=daemon` to ordinary ensigns (reserved for FOs until Paseo ships a parent-only daemon ACL).
- **Provider / model / thinking split:** Always pass an explicit `--provider <provider>/<model>` (never rely on `build.model`, which is `null`), and pass thinking effort separately via `--thinking <id>` when requested (for example `--provider antigravity/gemini-4-neon --thinking high` or `--provider codex/<model>`).
- **Frozen shared label contract:**
  - `spacedock.worker` = `{worker_base}` (= `envelope.name` emitted by `spacedock dispatch build --host paseo` or `spacedock dispatch name`). Invariant across cycles and retries (always <= 56 chars); NEVER append `-retry` or `-cycleN` to `spacedock.worker`.
  - `spacedock.epoch` = `{worker_base}-c{cycle}-a{attempt}` (`{worker_base}-c1-a1` for v1 fresh dispatch).

## Pre-Spawn Worktree, Project-ID, and Safehouse Handling

1. **Build envelope:** Run `${SPACEDOCK_BIN:-spacedock} dispatch build --host paseo --stamp --workflow-dir {workflow_dir} --entity-path {entity_path} --stage {stage} --checklist-file -` (omit `--bare-mode`). Retain `name` (`{worker_base}`), `description`, `dispatch_file_path`, and `prompt`. Set `epoch="${name}-c1-a1"`.
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
  --label "spacedock.epoch=${name}-c1-a1" \
  --json "$prompt"
```

- **JSON stdout fields:** `agentId` (UUID), `status` (`"created"|"running"`), `provider`, `cwd`, `title`.
- **Exit codes & branching:**
  - `0`: **parse stdout only** (stderr prints `Using workspace <id>` on exit `0`). Record `agentId`, `workspaceId`, `spacedock.worker=$name`, `spacedock.epoch=${name}-c1-a1`, `cwd`, and requested/returned provider, model, thinking, and mode before waiting. CLI `--no-wait` leaves `notifyOnFinish=false`.
  - `1`: `MISSING_PROVIDER`, `INVALID_OPTIONS`, `INVALID_THINKING_OPTION`, `WORKSPACE_NOT_FOUND`, `AGENT_CREATE_FAILED`.

### Verb 3 — Bounded wait (`wait`)

```bash
paseo wait "$agent_id" --timeout 300 --json
```

- **JSON stdout fields:** `agentId`, `status` (`"idle"|"timeout"|"permission"|"error"`), `message` (includes last 5 activity items on `idle`/`timeout`).
- **Exit codes & branching:**
  - `0` for ALL FOUR statuses — ALWAYS branch on JSON `status`:
    - `timeout` -> keep waiting (normal progress, 0 retries).
    - `idle` -> wake signal ONLY; verify durable Git commits + `## Stage Report: {stage}` before advancing.
    - `permission` -> run `paseo inspect "$agent_id" --json` and escalate to the captain.
    - `error` -> run `paseo inspect "$agent_id" --json` and hold.
  - `1`: `AGENT_NOT_FOUND`, `INVALID_TIMEOUT`, `WAIT_FAILED`.

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

## Runtime implementation

- `«worker.spawn»` -> `${SPACEDOCK_BIN:-spacedock} dispatch build --host paseo --stamp` -> resolve canonical `code_cwd` -> Verb 1 find-or-create (`workspace ls` then `workspace create --isolation local` if absent) -> Verb 2 (`run --no-wait`, stdout-only JSON parse, `spacedock.worker=$name`, `spacedock.epoch=${name}-c1-a1`). Forward the intact helper pointer verbatim (optional ordinary FO scope notes appended after a newline).
- `«addressable-worker»` -> ABSENT in v1 (fresh one-shot workers only; worker reuse and child-to-parent messaging are out of scope).
- `«async-dispatch»` -> ASYNC (`paseo run --no-wait` returns `agentId`; bounded polling via Verb 3 `paseo wait --timeout 300 --json`).
- `«worker-identity»` -> record `agentId`, `workspaceId`, `spacedock.worker` (`{worker_base}`), `spacedock.epoch` (`{worker_base}-c1-a1`), `cwd`, and requested/returned `provider`, `model`, `thinking`, `mode`; canonical model space is `<provider>/<model>` (+ optional `--thinking <id>`).
- `«completion-signal»` -> Verb 3 `status == "idle"` wakes verification; Git state (committed worktree deliverable + path-scoped `## Stage Report: {stage}` commit on canonical entity + clean `git status --porcelain`) is the completion gate.
- `«worker.shutdown»` -> Verb 5 `paseo archive "$agent_id" --json` after Git verification (or `--force` on supersede); guard `agentId != $PASEO_AGENT_ID` and never archive an FO with live children (same-workspace cascade-archive rule).
- `«context-budget»` -> ABSENT (`inspect.LastUsage` omits context-window max and v1 dispatches fresh only).
- `«roster-reconcile»` -> Verb 4 (`paseo ls --global --label "spacedock.worker=..." --json` + `paseo inspect "$agent_id" --json`).
