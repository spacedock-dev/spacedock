# Sprint paseo-backend

**Goal.** Let any Paseo-hosted first officer (Claude, Codex, Pi, or Antigravity) dispatch fresh ensigns through Paseo and verify the work from Git, instead of through the host's native subagents.

**Deliverable.** An opt-in Paseo dispatch binding, `skills/first-officer/references/paseo-fo-dispatch.md` (CLI-first with MCP optional, plus `PASEO_AGENT_ID` selection in `skills/first-officer/SKILL.md`), with a live acceptance test. The existing `spacedock dispatch build --host codex --stamp` supplies Codex assignments unchanged; a host-neutral `--host paseo` mode plus `paseo-ensign-runtime.md` extends the same path to non-Codex Paseo workers (`antigravity`, `opencode`, `copilot`, etc.). [Staff report](staff-report.md)

## Why this sprint exists

A Codex first officer drove a Paseo-launched worker by hand through a full stage, and a Claude-marked build produced the same Codex assignment bytes. [Git proof](run-log.md#completion-and-durable-proof) · [Cross-host envelope comparison](run-log.md#cross-host-envelope) A live `antigravity` (`agy`) first officer running under Paseo also drove a fresh `safehouse-codex` ensign end-to-end via the `paseo` CLI alone (auto-linking `ParentAgentId` via `PASEO_AGENT_ID`, committing both the code receipt and split-root stage report in 49 s, verified via `paseo wait` and archived). The Paseo project has operated this way across a Mac and an exe.dev VM. That work depended on manual steps Spacedock should own: prompts patched to replace the native completion message, the first officer relaying worker reports into tasks, and inlined assignments for remote workers. Doing this at the Paseo layer supersedes per-host direct runtime integrations such as `27a` (`agy-runtime-support`).

## First version

Fresh Codex workers run on the first officer's Paseo daemon via `paseo` CLI (with Paseo MCP optional when available), and non-Codex Paseo workers (`antigravity`, `opencode`, `copilot`, etc.) follow via a host-neutral `--host paseo` build mode. Remote workers follow after the local lane.

- **FO transport is CLI-first.** Bind `references/paseo-fo-dispatch.md` to the `paseo` CLI (`paseo workspace create`, `paseo run --no-wait`, `paseo wait`, `paseo inspect`, `paseo send`, `paseo archive`) with MCP optional, and select it from `skills/first-officer/SKILL.md` when `PASEO_AGENT_ID` is set (or Paseo is explicitly selected). This unlocks any Paseo provider (including `antigravity` / `agy`, which exposes no Paseo MCP tools) as FO with zero host-specific Go code. `paseo run --no-wait` automatically sets `ParentAgentId` from `PASEO_AGENT_ID`.
- **Completion comes from state.** The worker commits its own stage report to the state branch. A finish notification or `paseo wait` only wakes the first officer. Polling is a first-class path, because a first officer started outside Paseo, for example `spacedock claude` in a terminal, has no Paseo agent identity to receive notifications. No step depends on the worker messaging its parent.
- **A worker record in the task.** The first officer records host, server ID, agent ID, stage, worktree, and the requested and returned mode and model. A new first officer adopts a running worker from this record instead of dispatching a second one.
- **Workspace and launch.**
  - Adopt the stamped worktree with `create_workspace`, passing the first officer's own Paseo project ID. Without it, Paseo creates a new project.
  - Launch with an explicit provider, model and mode, and check all three in the returned snapshot. On a Safehouse machine the provider is a `safehouse-*` one: plain `codex` cannot read `~/.codex` from inside Safehouse.
  - Keep Spacedock labels in a `spacedock.*` namespace, apart from Paseo's own `paseo.*` labels.
- **Ownership.** Spacedock owns the worktree, the assignment and its write scope. Paseo owns the workspace and the worker handle.
- **Lifecycle.**
  - Workers are the first officer's Paseo children: it shows as working while they run, and archiving it archives them too. The binding never archives the first officer while a worker is live.
  - Archive each worker after its report is verified.

## Sequence

1. **Live acceptance first, riskiest path.**
   - A real Claude first officer launches a clean Codex worker: a fresh environment where the `spacedock:ensign` skill must be found through the installed plugin.
   - The worker commits its stage report on the state branch itself, which proves its sandbox can reach the state checkout outside the worktree.
   - The first officer advances from that committed report with `notifyOnFinish=false`.
2. **The binding.** Add a skill smoke test for selecting Paseo, then `references/paseo-fo-dispatch.md`. Route the acceptance fixture through it.
3. **Worker record and adoption.** Kill the first officer mid-stage. A new first officer adopts the running worker from the record, with no second dispatch.
4. **Completion, recovery and closure.** Fixtures for:
   - duplicate notifications;
   - a lost create receipt;
   - a worker that ended without a report;
   - a wait timeout during real progress;
   - a second failure that holds only its own task.

   Each fixture also checks that verified workers are archived.
5. **Host-neutral Paseo workers (`--host paseo`).** Add `--host paseo` to `spacedock dispatch build` plus `skills/ensign/references/paseo-ensign-runtime.md` (host-neutral `ensign` skill load, split-root state commit, concise final turn message) so `antigravity` and other Paseo providers can run as ensigns without per-host Go branches. Prove with a live `antigravity` ensign once Paseo's per-session `agent.open.before` `wrap` hook unblocks `paseo-safehouse` on newly stamped worktrees.
6. **Remote workers (`serverId`).** Needed early in practice: the Mac's free disk could not hold an `npm install`. This step covers:
   - the assignment sent inline, because the dispatch file sits in the first officer's `/tmp`;
   - a state remote the worker can reach;
   - push credentials that expire (exe.dev grants 48 hours);
   - paths and free disk that differ per host.
7. **Child-to-parent questions** once Paseo has a parent-only daemon permission. A child with `paseo-safehouse.paseo-access=daemon` already delivered a message, but that label grants too much. [Child return channel](run-log.md#child-return-channel)

## Definition of done

Steps 1 to 4 deliver the local lane. Through the installed binding, a real Claude (or Antigravity) first officer:
- launches a clean Codex worker;
- advances from the worker's own committed report with notifications off;
- adopts a running worker after a restart without dispatching again;
- archives the verified worker.

Step 5 extends the same local lane to non-Codex Paseo workers (`antigravity`, etc.) via `--host paseo` once the upstream `agent.open.before` launch-wrapper hook lands.

The acceptance record includes the requested and returned Paseo modes.

## Related work outside this sprint

- **`dispatch build` guard:** refuse a stamped `worktree:` that is not an existing path. A stamp that held a description, not a path, was carried into the next dispatch.
- **First-officer boot guard:** refuse a status change or dispatch from a session that never booted as first officer.
- **Paseo-side findings (resolved):** `create_agent` requested with `modeId=auto` returned `full-access` because `paseo-safehouse` (`server/service.ts`) intentionally forces `modeId: "full-access"` / `sandbox_mode: "danger-full-access"` for wrapped Codex sessions — Codex cannot nest macOS Seatbelt inside Safehouse. [First launch](run-log.md#workspace-and-spawn)
- **Upstream Paseo launch-wrapper hook & `paseo-safehouse`:** Launching `antigravity` into a freshly stamped worktree currently fails in `paseo-safehouse`'s `agent.create.before` because automatic Safehouse defaults on for new worktree cwds and `antigravity` is a plugin-supplied provider not in `paseo-safehouse`'s `factories` list. Upstream prototype for `agent.open.before` `wrap` (coordinated with `paseo-safehouse` manager `7d469e54` and `paseo` manager `4fc75d8c`) removes per-provider `safehouse-*` entries and enables Safehouse for Antigravity and all SDK ACP providers.
- **Antigravity skill freshness:** `~/.gemini/config/plugins/spacedock/` can drift behind the `spacedock` binary; keep `agy` skills aligned via `~/.agents/skills` / `spacedock install` / `paseo-spacedock`.
- **Paseo daemon permissions:** a parent-only permission for child-to-parent messages, owned by the Paseo daemon access-control work.

## Out of scope

- Worker reuse across stages.
- Several first officers on one workflow, or leases.
- Moving a worker to a new parent across hosts.
- Automatic failover.
- Implementing the Paseo daemon access control.
- Direct per-host Antigravity CLI support (`27a` / `agy-runtime-support`, superseded by this sprint).
