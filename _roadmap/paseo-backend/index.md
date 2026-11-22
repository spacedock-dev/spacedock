# Sprint paseo-backend

**Target release train:** `0.29.x` (movable line — rebindable without changing member frontmatter).
**Members:** `spacedock status --workflow-dir docs/dev --where sprint=paseo-backend` (drivable: `--where 'sprint-readiness != defer'`).

**Goal.** Let any Paseo-hosted first officer (Claude, Codex, Pi, or Antigravity) dispatch fresh ensigns through Paseo and verify the work from Git, instead of through the host's native subagents.

**Deliverable.** One host-neutral Paseo dispatch path (`spacedock dispatch build --host paseo`, `skills/ensign/references/paseo-ensign-runtime.md`, and CLI-first `skills/first-officer/references/paseo-fo-dispatch.md` wired from `PASEO_AGENT_ID`), proven by a real isolated-daemon live smoke (`TestLivePaseoDispatchSmoke`), followed by restart adoption and recovery (`TestLivePaseoAdoptionSmoke`). [Staff report](staff-report.md)

## Why this sprint exists

Two live probes proved both halves of running Antigravity (`agy`) and Codex through Paseo with zero `agy`-specific Go code:
1. A live `antigravity` FO drove a fresh `safehouse-codex` ensign end-to-end via the `paseo` CLI alone (auto-linking `ParentAgentId` via `PASEO_AGENT_ID`, committing both the code receipt and split-root stage report in 49 s, verified via `paseo wait` and archived).
2. A live sandboxed `antigravity` (`gemini-4-neon`) ensign ran inside Safehouse via Paseo's `launch-wrap-hook` (`9f9f39ce1`) + `paseo-safehouse` (`0ff8d55`, `route: wrap`), read a host-neutral prompt, committed the code receipt (`c24eda1`), and committed its path-scoped split-root stage report (`9b1476c`).

Doing this at the Paseo layer supersedes direct per-host integrations such as `27a` (`agy-runtime-support`, archived) and replaces manual prompt-patching and hand-relayed reports. [Git proof](run-log.md#completion-and-durable-proof) · [Cross-host envelope comparison](run-log.md#cross-host-envelope)

## Design choices

- **One `--host paseo` for all Paseo workers.** Under Paseo, every worker (Codex, Antigravity, Claude, Pi) shares the exact same contract: load the `ensign` skill, commit deliverables to the code worktree and `## Stage Report` path-scoped to the state checkout, and finish the turn with `Done: ...`. No worker-runtime branching in `paseo-fo-dispatch.md`.
- **CLI-first FO transport (5 verbs).** Bind `references/paseo-fo-dispatch.md` to `paseo workspace create`, `paseo run --no-wait`, `paseo wait`, `paseo inspect`/`paseo ls`, and `paseo archive` (MCP optional), selected in `skills/first-officer/SKILL.md` when `PASEO_AGENT_ID` is set or Paseo is explicitly chosen.
- **Completion comes from Git state.** `paseo wait` / finish events only wake the FO; polling with `notifyOnFinish=false` is first-class.
- **Real-daemon live tests, zero fake ACP mocks.** To track fast-moving Paseo honestly without maintaining mocks or pin machinery, `internal/ensigncycle/` runs `TestLivePaseoDispatchSmoke` (and later `TestLivePaseoAdoptionSmoke`) against a real isolated `paseo` daemon (`SPACEDOCK_PASEO_BIN`, default `paseo` on `PATH`, canonical `t.TempDir()`, free port, stopped in `t.Cleanup`).
- **Optional CI lane (like `pi-live`, not required).** Runs `TestLivePaseoDispatchSmoke` with `codex/gpt-6-luna` on the existing Codex subscription (`CODEX_AUTH_JSON`). Locally overridable via `SPACEDOCK_PASEO_LIVE_PROVIDER` (e.g. `antigravity/gemini-4-neon`) and `SPACEDOCK_PASEO_BIN`.
- **Ownership & lifecycle.** Spacedock owns the worktree, assignment, and write scope; Paseo owns the workspace and worker handle (`spacedock.*` labels). Archive verified workers; never archive the FO while a child worker is live.

## Sequence

1. **`paseo-dispatch-backend` (local lane + E2E test).**
   - Add `--host paseo` to `spacedock dispatch build` + `skills/ensign/references/paseo-ensign-runtime.md`.
   - Add `skills/first-officer/references/paseo-fo-dispatch.md` (5 CLI verbs) + `PASEO_AGENT_ID` selection in `skills/first-officer/SKILL.md`.
   - Add `TestLivePaseoDispatchSmoke` in `internal/ensigncycle/` (isolated daemon, `SPACEDOCK_PASEO_BIN`, `SPACEDOCK_PASEO_LIVE_PROVIDER` defaulting to `codex/gpt-6-luna`, durable Git + archive assertions) and wire an optional (non-required) CI workflow lane using `codex/gpt-6-luna` on the Codex subscription.
2. **`paseo-worker-adoption-and-recovery` (resilience follow-up).**
   - Persist worker record (server ID, agent ID, stage, worktree, requested/returned provider, mode, model) in the task so a restarted FO adopts live workers instead of double-dispatching.
   - Add timeout, lost-create-receipt, dead-worker, and one-retry-then-hold recovery rules to `paseo-fo-dispatch.md`.
   - Prove with `TestLivePaseoAdoptionSmoke` (kill FO mid-stage, boot fresh FO, assert single worker and clean completion).

## Definition of done

Through the installed binding and `TestLivePaseoDispatchSmoke`, a real Paseo FO:
- builds a host-neutral `--host paseo` assignment and adopts the stamped worktree;
- dispatches a fresh Paseo worker (`codex/gpt-6-luna` in CI, overridable to `antigravity/gemini-4-neon` locally);
- verifies the worker's own code and split-root state commits with notifications off;
- archives the verified worker;
- (after task 2) adopts a running worker across an FO restart without spawning a second one.

## Related work outside this sprint

- **Upstream Paseo launch-wrapper hook & `paseo-safehouse`:** `task/launch-wrap-hook` (`9f9f39ce1` in `paseo`) + `plugin-provider-sandbox` (`0ff8d55` in `paseo-safehouse`, plus follow-up `symlinked-daemon-home`) enable Safehouse-wrapped `antigravity` workers without `safehouse-*` entries. Until merged to production, `antigravity` workers require the test daemon build; `codex/gpt-6-luna` runs on production Paseo today.
- **Paseo-side findings (resolved):** `modeId=auto` returned `full-access` because `paseo-safehouse` (`server/service.ts`) intentionally forces `full-access` / `danger-full-access` for wrapped Codex sessions (Codex cannot nest macOS Seatbelt inside Safehouse).
- **Antigravity skill freshness:** `~/.gemini/config/plugins/spacedock/` can drift behind the `spacedock` binary; keep `agy` skills aligned via `~/.agents/skills` / `spacedock install` / `paseo-spacedock`.
- **`dispatch build` / FO boot guards:** general Spacedock hardening outside this sprint.
- **Remote workers (`serverId`) & child-to-parent questions:** deferred until the local lane and its live E2E test are proven in daily use (and until Paseo ships a parent-only daemon ACL).

## Out of scope

- Remote workers (`serverId`) and expiring remote push credentials (deferred).
- Child-to-parent mid-run questions (blocked on Paseo parent-only daemon ACL).
- Worker reuse across stages.
- Fake ACP stub servers or Paseo version-pinning machinery.
- Several first officers on one workflow, or leases.
- Direct per-host Antigravity CLI support (`27a` / `agy-runtime-support`, superseded and archived).
