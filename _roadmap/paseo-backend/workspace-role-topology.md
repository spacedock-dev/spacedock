# Organizing Project Roles Across Paseo Workspaces and Daemons

**Recommendation (simplest design that works first):** Collapse **Manager + Shaping FO + Commander** into **one Project FO per repository**, keep **ensigns ephemeral** (same workspace for non-worktree stages, `--background` worktree workspaces archived on teardown), and run **one Paseo daemon per machine** (with per-agent Safehouse policy). Until Safehouse moves to the primary Mac daemon, co-locate a Safehouse-tested repo's FO and ensigns together on port `6769`.

---

## 1. Why the current setup feels messy

Running four overlapping roles across split daemons created four concrete failures during `paseo-backend`:

1. **Duplicate ownership and double-wakes:** A project Manager (`ba80b87`), Shaping FO (`846d9a2`), and Sprint Commander (`3b905937`) all existed for one repo. Tonight both `846d9a2` and `3b905937` woke up on the same friction report and had to pause to avoid colliding on `_roadmap/` and the PR stack.
2. **Split Mac daemons (`6767` vs `6769`):** Putting the Shaping FO on `6767` and the Commander/Gemini ensigns on `6769` broke native `PASEO_AGENT_ID` parent-child links, blocked MCP orchestration tools, forced `env -u PASEO_* paseo --host 127.0.0.1:6769` shell bridges, and required manual watchdog polling.
3. **Workspace sidebar clutter:** Every worktree stage created a visible top-level workspace (`wks_...`) in Paseo unless `--background` was passed and `workspace archive` ran at teardown.
4. **Brittle reverse-RPC across machines (`host_unknown`):** Agents on remote VMs (`studfleet-paseo`) have no registered destination host back to the Mac daemon, so reverse `send` calls fail whenever a workflow expects a direct remote receipt.

---

## 2. Simplest design first (Target State)

### Roles: 2 persistent classes, not 4

| Role | Count & Lifetime | Daemon Access (`paseo-safehouse`) | Paseo Workspace | Owns |
|---|---|---|---|---|
| **MetaFO** | 1 across all projects (persistent) | `daemon` (or hub-only) | Meta workspace | Cross-project briefings, relaying CL directives to the **one** named driver per repo/stack. Never drives repo worktrees. |
| **Project FO** *(replaces Manager + Shaping FO + default Commander)* | 1 per repo (persistent) | `daemon` | Repo root (visible) | Backlog/shaping, dispatching ensigns, presenting gates to CL, PR delivery. Single writer for the repo's active stack. |
| **Ensigns** | Ephemeral per stage | `workspace` (no daemon access) | Same workspace as FO (non-worktree) **or** `--background` local-adopted worktree | Assigned code, reports, audits. Archived immediately on completion; worktree workspaces archived on teardown. |
| **Sprint Commander** *(optional, opt-in only)* | 0 normally; 1 ephemeral during a multi-entity autonomous sprint | `daemon` | Same repo-root workspace (subagent track of Project FO) | Entities matching `sprint=<X>` only. Archived when the sprint ends. Never runs in parallel with a second driver on the same stack. |

### Workspace visibility rule (matches Paseo's `pi-subagents` UX)

- **Non-worktree stages** (`ideation`, read-only audits, staff reviews): spawn in the Project FO's repo-root workspace → appears only in the FO's **subagent track** (zero extra sidebar workspaces).
- **Worktree stages** (`implementation`, `validation`): Spacedock creates `.worktrees/spacedock-ensign-<slug>`; the FO adopts it with `paseo workspace create --isolation local --background --path <worktree> --title "<entity>"`, and runs `paseo workspace archive <wks_id>` when the worktree is removed.

### Daemon topology

- **One Paseo daemon per machine:** Run the Safehouse-capable Paseo daemon (`0.17.4+`) as the single Mac daemon. Enforce privilege per agent (`paseo-access=daemon` on FOs, `paseo-access=workspace` on ensigns), not by splitting daemons. Every Gemini, Codex, and Claude agent shares one parent-child graph and native completion wakes.
- **Remote VMs (`studfleet-paseo`, etc.):** Keep each VM self-contained (1 Project FO + its ensigns on the VM's local daemon). Communicate across machines through durable artifacts (Git branches, PRs, state branch, or files pulled by MetaFO) — never block a remote FO on a reverse `send` to the Mac. Register mutual Paseo hosts (`paseo host add`) only when interactive fleet RPC is explicitly needed.

### Single-driver ownership rule

- Every active task or PR stack has **one named driver** (default: the Project FO; during a sprint: the Sprint Commander).
- MetaFO sends directives **only** to that one driver — never to Manager + Shaping FO + Commander at once.
- Once `79` (`entity-session-claim-lease`) lands, the state branch enforces this lock mechanically across sessions and machines.

---

## 3. Interim setup today (while Gemini + Safehouse stay on port `6769`)

Until Safehouse moves onto the primary `6767` daemon, avoid splitting a single repo across two Mac daemons:

1. **Co-locate Safehouse repos on `6769`:** Run the repo's **single Project FO** (Gemini or Codex, `paseo-access=daemon`) and **all its ensigns** (`paseo-access=workspace`) together on `127.0.0.1:6769`. Retire the separate `6767` Manager / Shaping FO bridge for that repo.
2. **Keep MetaFO one-way (`6767` → `6769` / VM):** MetaFO reads status/briefings and dispatches directives to the one repo FO on `6769` (via `--host 127.0.0.1:6769` or host registration). The `6769` FO never needs to reach back to `6767` mid-run; it writes durable state and replies in its own turn.
3. **Retire `6769` once Safehouse graduates to `6767`:** When the primary Mac daemon runs the Safehouse-enabled build, move the Project FO back to `6767` and delete the second daemon.

---

## 4. Concrete next actions (when CL approves)

1. **Adopt 1 Project FO per repo immediately:** Stop maintaining separate Manager + Shaping FO + Commander agents for `spacedock-v1`; keep one FO on `6769` while Safehouse lives there.
2. **Hide and clean up worktree workspaces:** Add `--background` on worktree `workspace create` and `paseo workspace archive` on worktree teardown in `paseo-fo-dispatch.md` (already tracked as a post-merge follow-up).
3. **Graduate Safehouse to the primary Mac daemon** (after upstream `paseo-safehouse` `tv` and `nc` land) so the `6769` split disappears completely.
