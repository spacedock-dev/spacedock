---
id: 9243ej9j3a434c3k2ba3g06j
title: Bridge using-claude-team to .178+ merged team semantics (with a legacy-deprecation path)
status: backlog
source: "FO + captain live spike (2026-06-18) on Claude Code 2.1.181: validated the subagent-as-team-member machinery after anthropics/claude-code#68721 removed TeamCreate/TeamDelete (followup comment 4741152246, by bjcoombs, confirms the team->subagent merge is intentional). using-claude-team is built on the now-removed TeamCreate primitive and, on .178+, silently degrades to bare mode (surrendering concurrency) because it gates all dispatch on a TeamCreate that no longer exists."
started:
completed:
verdict:
score:
worktree:
issue: "anthropics/claude-code#68721"
---

Make `skills/using-claude-team` work on Claude Code .178+ (the merged auto-team model, where `TeamCreate`/`TeamDelete` are gone and team membership comes from spawning named background subagents), while remaining compatible with legacy hosts that still expose `TeamCreate` — and mark a clear, triggered path to deprecate the legacy branch once the supported floor is .178+.

## Problem

`using-claude-team` makes `TeamCreate` the mandatory first call and **blocks all `Agent` dispatch until it resolves**. On .178+, `ToolSearch(select:TeamCreate)` returns no match, so the skill falls into **bare mode**: sequential, no concurrent dispatch, no direct ensign chat. Nothing breaks (all workflow functionality is preserved), but spacedock **needlessly surrenders concurrency** by gating on a primitive that no longer exists, while the primitives it actually needs (named background `Agent`, `SendMessage`, Tasks, on-disk team registry) are all present and working. The skill, the `claude-fo-dispatch.md` reconcile step, and the Terminal-Teardown / Degraded-Mode machinery are all written around `TeamCreate`/`TeamDelete` and the `#36806` registry-desync model, which no longer describe the runtime.

## Proposed approach (seed — ideation to flesh out)

Two-mode bridge with a deprecation trigger:

- **Detect mode**, do not assume it. Merged (.178+) vs legacy is observable: `ToolSearch(select:TeamCreate)` empty => merged; present => legacy. Ideation pins the exact detection contract and where it lives (generic skill vs `claude-fo-dispatch.md`).
- **Merged path (.178+):** no `TeamCreate`. Team membership = `Agent(name=…, run_in_background=true)`. Lead<->teammate via `SendMessage(to=name)` / teammate->lead via `SendMessage(to="main")`. Teardown = per-name `SendMessage(shutdown_request)` -> `shutdown_approved` -> `teammate_terminated` (NO `TeamDelete`, no active-member race, no `TERMINAL_TEARDOWN_BOUNDED` apparatus). FO tracks its own ensign roster (it already does — `TaskList` is not used).
- **Legacy path:** keep the existing `TeamCreate`/`TeamDelete` flow intact, clearly fenced and labelled deprecated, with a removal trigger (e.g. "remove when the min supported Claude Code no longer exposes `TeamCreate`, and no live lane still drives the legacy branch").
- **Strip the `claude-team`/NAME_PATTERN leak** from the generic skill (line ~19): the FO no longer names the team, and the introspection helper keys on `leadSessionId`, never the team name, so the compat note is doubly obsolete.
- **Flip `claude-fo-dispatch.md`'s reconcile step** from "pass your `TeamCreate` `{team_name}`" to the already-documented `leadSessionId` **auto-discovery** mode (thread `--session-id`, drop `--team-name`). The Go helper (`internal/claudeteam/reconcile.go`) needs **no change**.

## Out of scope

- The live pty/team-mode harness work — owned by `m4` (`live-team-mode-terminal-harness`, m40mphxan8phr3t3tp03gk89). This task is the **skill + FO-contract prose** layer; m4 is the **live test harness** layer. They must stay reconciled but are separate deliverables.
- The team-mode verdict-omission question (reeppr990pyzzaejmbnyrvt7).
- Any change to the Go `internal/claudeteam` helper — the spike showed it is already robust to the merged model.

## Spike findings — validated live 2026-06-18 on Claude Code 2.1.181

Probed with real `Agent` + `SendMessage` (NOT the Workflow harness — a different code path). Four background subagents spawned (`scout`, `engineer`, `supervisor`, + a blocked `helper`); all reaped cleanly. Backend was `in-process`, model `claude-opus-4-8[1m]`.

1. **TeamCreate/TeamDelete are gone** on 2.1.181: `ToolSearch(select:TeamCreate,TeamDelete)` => "No matching deferred tools found". Intentional per #68721 followup.
2. **Named background subagent works with zero TeamCreate.** `Agent(name=…, run_in_background=true)` spawns an addressable teammate.
3. **subagent -> lead channel is OPEN.** A background subagent's `SendMessage(to="main")` was delivered to the root verbatim. (Old model: a raw subagent could only reach the lead via its final return value.)
4. **Standing member works.** A persistent teammate retained state across waves (recalled a secret number 42 -> 84 and an incrementing counter across separate inbound messages), was addressable each round, and replied to the lead each round.
5. **Reap works via lead-initiated shutdown_request.** `SendMessage(shutdown_request)` -> teammate emits `shutdown_approved` -> `teammate_terminated`. This CONTRADICTS the #68721 comment's claim that "the lead's shutdown_request isn't honored" — at least for `in-process` backend on 2.1.181.
6. **Self-exit is NOT available — and this is NOT a regression.** A teammate can only send `shutdown_response` as a REPLY to an inbound `shutdown_request`; it cannot originate termination. Old `TeamCreate` members were identical. The comment's "can't send shutdown_response" was likely observed because no `shutdown_request` was ever sent (no `request_id` to echo).
7. **Roster is FLAT by construction.** A background teammate CANNOT spawn background grandchildren: `Agent(name=…, run_in_background=true)` from a teammate => "Teammates cannot spawn other teammates — the team roster is flat"; nameless `run_in_background=true` => "In-process teammates cannot spawn background agents. Use run_in_background=false for synchronous subagents." So the "orphan-grandchild" hazard cannot occur; recursive *background* depth is bounded to 1. Synchronous (blocking) subagents at depth are still allowed — so an ensign running a fan-out skill/Workflow still works.
8. **On-disk registry is intact + accurate.** The auto-team writes `~/.claude/teams/session-<id>/config.json`, named by **session id** (not the FO scheme), with the SAME schema (`name`/`leadAgentId`/`leadSessionId`/`members[]`/`inboxes/`). `members[]` is a **live pruned roster** (terminated members removed; `inboxes/*.json` linger). Each member records `agentType` = the spawn `subagent_type` (probe recorded `general-purpose`).
9. **The `claude-team` helper survives unchanged.** `internal/claudeteam/reconcile.go` auto-discovers by globbing all team dirs and matching `leadSessionId == current session` AND a member with `agentType == "spacedock:ensign"`; `--team-name` is only an optional override. So the session-id rename does NOT break discovery — the helper never used the team name. Only the CALLER must switch to auto-discovery mode (no `TeamCreate` name to pass).

## Capability matrix — Claude Code 2.1.181 (merged auto-team model)

Tested context: Claude Code **2.1.181**, interactive `cli` entrypoint (**NOT** headless `-p`), `in-process` teammate backend, model `claude-opus-4-8[1m]`, `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`, no tmux. Each ✅/❌ was observed live this session; ❓ rows were NOT exercised and must not be read as either pass or fail.

| Capability | 2.1.181 | Evidence / exact signal |
|---|---|---|
| `TeamCreate` / `TeamDelete` tools | ❌ absent (by design) | `ToolSearch(select:TeamCreate,TeamDelete)` → "No matching deferred tools found" |
| Spawn named background teammate without `TeamCreate` | ✅ | `Agent(name=…, run_in_background=true)` returns a live addressable teammate |
| Teammate → lead message | ✅ | `SendMessage(to="main")` delivered to the root verbatim |
| Lead → teammate message | ✅ | `SendMessage(to=<name>)` reaches the teammate inbox |
| Standing teammate persists + retains state across waves | ✅ | recalled a secret value + an incrementing counter across separate inbound messages |
| Lead-initiated reap via `shutdown_request` | ✅ (in-process) | `shutdown_request` → teammate `shutdown_approved` → `teammate_terminated` |
| Teammate self-originated exit | ❌ (NOT a regression) | `shutdown_response` is valid only as a REPLY to an inbound `shutdown_request`; legacy `TeamCreate` members were identical |
| Background teammate spawns a background **grandchild** | ❌ blocked (flat roster) | "Teammates cannot spawn other teammates — the team roster is flat"; nameless bg → "In-process teammates cannot spawn background agents. Use run_in_background=false for synchronous subagents." |
| Teammate spawns a **synchronous** subagent | ✅ allowed | per the harness's own error guidance (`run_in_background=false`) |
| `TaskList` / `TaskStop` reach background teammates | ❌ | `TaskStop <id>` → "No task found"; `TaskList` → empty |
| On-disk team registry | ✅ present | `~/.claude/teams/session-<id>/config.json`; same schema as the `TeamCreate` era; `members[]` is a live, pruned roster (terminated members removed) |
| Headless **`-p`** team mode (residency past `end_turn`) | ❓ UNTESTED | only interactive `cli` exercised this session |
| **tmux pane-backed** teammate reap | ❓ UNTESTED | all probes were `in-process` |
| **Flag-free** (no `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`) | ❓ UNTESTED | the flag was `=1` throughout |
| `agentType` stamped for a non-`general-purpose` `subagent_type` (e.g. `spacedock:ensign`) | ❓ UNTESTED | probes used `general-purpose`, which the registry recorded faithfully |

## Open questions (for ideation)

- **OQ-1 (load-bearing): does lead `shutdown_request` reap a tmux PANE-BACKED teammate, or only `in-process`?** The spike was all in-process; the #68721 comment's "shutdown not honored" may be pane-specific. m4 spawns teammates into separate tmux panes — directly relevant. This is the one finding that could still be a true per-member-reap regression.
- **OQ-2: does `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` still gate the named-background-subagent + `SendMessage(to="main")` channel on .178+, or is it now flag-free?** During the spike the flag was PRESENT (`=1` in the session env), so flag-free operation is UNTESTED. m4's CI sets the flag. The skill's bare-mode hint tells captains to set it — possibly stale.
- **OQ-3: does a real `subagent_type="spacedock:ensign"` dispatch record `agentType:"spacedock:ensign"` in the auto-team config**, so `reconcile.go`'s `hasEnsign` gate fires? Needs a one-shot real-ensign smoke test (probe used `general-purpose`, which would be correctly excluded).
- **OQ-4: exact mode-detection contract.** Is `ToolSearch(select:TeamCreate)` presence reliable across 2.1.177 <-> .178 <-> .181, and does it belong in the generic skill or the spacedock FO layer?
- **OQ-5: is `--session-id` plumbed** through `spacedock claude`/the launcher so the reconcile sweep can run in `leadSessionId` auto-discovery mode without a `TeamCreate` name?
- **OQ-6: reconcile with m4.** m4's transcript-observation work and this skill rewrite both touch team-mode semantics; ideation must define the seam so they do not collide (skill/contract here, live harness there).

## Provisional acceptance criteria (ideation finalizes; proof = behavior, never prose-grep)

Per this workflow's proof policy a skill change is PASSED only when a LIVE drive observed the behavior — a string match over the skill text proves nothing.

- **AC-1 (provisional):** On a .178+ host, a live FO boot dispatches a teammate concurrently WITHOUT calling `TeamCreate` (no bare-mode fallback). Verified by: a live drive (or the m4 pty harness) observing concurrent dispatch + the auto-team `config.json` on disk.
- **AC-2 (provisional):** On a legacy host that still exposes `TeamCreate`, the legacy path still runs. Verified by: the legacy live lane (or a fixture that stubs `TeamCreate` present) still green.
- **AC-3 (provisional):** Terminal teardown on .178+ reaps every tracked ensign via per-name `shutdown_request` (no `TeamDelete`), leaving `members[]` empty. Verified by: live drive observing `teammate_terminated` for each + empty roster.
- **AC-4 (provisional):** The reconcile sweep resolves the auto-team via `leadSessionId` auto-discovery with no `TeamCreate` name passed. Verified by: `spacedock dispatch reconcile` (no `--team-name`) resolving the live session's roster.
- **AC-5 (provisional):** The legacy branch is fenced + labelled deprecated with an explicit removal trigger. Verified by: the deprecation trigger references an external, checkable condition (supported-version floor / live-lane coverage), not just prose.

## Related

- `m4` — live-team-mode-terminal-harness (m40mphxan8phr3t3tp03gk89, PR #390): live-harness layer of the same regression; OQ-1/OQ-2/OQ-3 are partly answerable there.
- `team-mode-verdict-omission` (reeppr990pyzzaejmbnyrvt7).
- anthropics/claude-code#68721 + comment 4741152246 (intentional team->subagent merge).
- Spike transcript: this FO session (2026-06-18), `~/.claude/teams/session-a15e8296/`.
