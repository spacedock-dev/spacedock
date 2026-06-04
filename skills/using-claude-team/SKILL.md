---
name: using-claude-team
description: "Generic Claude Code multi-agent team-harness discipline — team creation and failure-recovery, awaiting-completion idle guardrails, terminal team teardown, degraded-mode fallback, and the deferred team-tool ToolSearch hop. Invoke when orchestrating any Claude Code agent team (TeamCreate/Agent/SendMessage), independent of any specific workflow."
user-invocable: false
---

# Using Claude Team

This skill carries the generic Claude Code team-harness discipline: how to create a team, recover from registry desync, wait for an agent's completion without prematurely tearing the team down, tear the team down at the terminal boundary, fall back to degraded (bare) mode, and reach the deferred team tools. It contains zero workflow-specific content — a consuming contract invokes this skill to load the lifecycle, then keeps its own decision points inline.

## Deferred Team Tools

The Claude Code team tools (`TeamCreate`, `TeamDelete`, `SendMessage`, and related team-registry tools) are deferred — their schemas are not loaded at session start, so calling one directly fails until its schema is fetched. Before the first call to any team tool, run `ToolSearch(query="select:{ToolName}", max_results=1)` to fetch its schema (e.g. `ToolSearch(query="select:TeamCreate", max_results=1)` before the first `TeamCreate`, `ToolSearch(query="select:SendMessage", max_results=1)` before the first `SendMessage`). Once a tool's schema appears in the ToolSearch result, it is callable exactly like a normal tool. `Agent` is not deferred. Address an agent by its declared `name` via `SendMessage`; your plain text output is NOT visible to other agents.

## Team Creation

At startup (before dispatch):

1. **Probe for TeamCreate and run it first.** `TeamCreate` MUST be the first team-mode tool call in every session, before ANY `Agent` or `SendMessage` invocation. Run `ToolSearch(query="select:TeamCreate", max_results=1)`. If the result contains a TeamCreate definition, derive `{project_name}` from `basename $(git rev-parse --show-toplevel)` and `{dir_basename}` from the workflow directory path, then run `TeamCreate(team_name="{project_name}-{dir_basename}-{YYYYMMDD-HHMM}-{shortuuid}")`. The timestamp token must be lowercase and hyphen-separated — no uppercase, no colons — to stay compatible with Claude Code's NAME_PATTERN and the `claude-team` helper. `{shortuuid}` is eight lowercase alphanumeric characters.
   - **IMPORTANT:** TeamCreate may return a different `team_name` than requested. Always store the returned value and use it for all subsequent calls.
   - **NEVER delete existing team directories** (`rm -rf ~/.claude/teams/...`) — stale directories belong to other sessions.
2. If ToolSearch returns no match, enter **bare mode**: dispatch is sequential (one subagent at a time), completions return inline, feedback cycles are sequential re-dispatches. Report the mode to the captain. All workflow functionality is preserved. When reporting bare mode at startup, append this one-line UX hint verbatim — it goes to a captain who has not opted into team mode and may not know it exists: `Tip: set CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1 and restart this session to enable team mode for concurrent dispatch and direct ensign chat (Shift+Up/Shift+Down).` Do not repeat the hint after startup.

**Diagnostic-only startup probe:** At startup the FO MAY inspect `~/.claude/teams/` with `ls` or `test -f ~/.claude/teams/{project_name}-{dir_basename}*/config.json` to REPORT stale on-disk team directories from prior sessions in the boot summary. This probe is DIAGNOSTIC-ONLY. Its result does NOT gate or skip `TeamCreate` — `TeamCreate` always runs with the fresh-suffixed name regardless of what the probe reports. On-disk state is not evidence of team health (Claude Code bug anthropics/claude-code#36806 leaves config files on disk after the in-memory registry desyncs). Deletion remains forbidden per the NEVER-delete constraint above — the probe surfaces stale directories; it does not authorize removal. No such probe belongs in the pre-dispatch path.

**TeamCreate recovery procedure:** Call TeamDelete in its own message (no other tool calls). Wait for the result. Then call TeamCreate in a subsequent message. Store the returned `team_name`. Do NOT combine TeamDelete, TeamCreate, or Agent dispatch in the same message — Claude Code executes tool calls in a message in parallel, and dependent calls will race. This procedure applies ONLY to the narrow "Already leading team" case at startup (where Claude Code's in-memory slot holds a team the FO wants to replace cleanly). It is NOT a mid-session failure recovery and MUST NOT be invoked in response to "Team does not exist" or any other registry-desync signal — see Degraded Mode below.

**TeamCreate failure recovery (priority-ordered ladder):** If TeamCreate or any subsequent `Agent()` dispatch surfaces "Team does not exist" or any equivalent registry-desync signal mid-session, follow this ladder in order — do NOT retry within the same tier:

1. **Fresh-suffixed TeamCreate.** Attempt one new `TeamCreate` with a fresh name `{project_name}-{dir_basename}-{YYYYMMDD-HHMM}-{shortuuid}` computed at call time (new timestamp, new shortuuid, distinct from any name used earlier this session). Retry to the same team name is banned. Do NOT call `TeamDelete` on the failed team — the registry is already desynced and another `TeamDelete → TeamCreate` cycle will re-contaminate the same slot per anthropics/claude-code#36806. Store the returned `team_name`. All prior agent names are presumed zombified — do not SendMessage them; re-dispatch from entity frontmatter.
2. **Fall back to Degraded Mode per the Degraded Mode section below.** A second dispatch failure (including failure of the tier-1 fresh-suffixed TeamCreate, or a second "Team does not exist" at any point in the session) trips Degraded Mode immediately.
3. **Surface to captain** with an explicit recovery prompt if tiers 1 and 2 both fail (e.g., TeamCreate errors with quota or internal failure on the fresh name, AND Degraded Mode cannot be entered because `Agent` itself is unavailable). Do not silently retry. Do not block indefinitely — report the failure, name the tiers attempted, and wait for captain direction.

**Block all Agent dispatch** until team setup resolves (tier-1 fresh-suffixed TeamCreate succeeds or Degraded Mode is entered). Never dispatch while team state is uncertain.

## Degraded Mode

Degraded Mode is an explicit, session-wide mid-session transition. Once entered, it persists until the session ends — there is no recovery back to teams mode in the same session.

### Triggers

Any one of the following trips Degraded Mode:

- First "Team does not exist" error (or equivalent registry-desync signal) surfaced by `Agent()` or any team-registry tool.
- Any SECOND dispatch failure within the session — no time window, no durable counter. The counter-free rule is deliberate: the FO cannot reliably track failure timestamps across context pressure and idle notifications, so "second failure anywhere in the session" is the fail-early trigger.
- An explicit operator-initiated degrade command from the captain.

### Effects

Once Degraded Mode is active, the following invariants hold for the remainder of the session:

- No `team_name` parameter on any subsequent `Agent()` dispatch. The dispatch is built in bare mode (`team_name: null`, `bare_mode: true`) so the emitted Agent call has `name` and `team_name` absent.
- Every stage dispatches fresh and blocks until completion. No concurrent dispatch; one entity through one stage at a time.
- No SendMessage reuse of prior agent names. Stage advancement is always a fresh `Agent()` dispatch seeded from entity frontmatter. `SendMessage(to="{ensign_name}")` against any pre-degrade name is forbidden.

### Captain Report Template

On Degraded Mode entry, the FO emits the following sentence verbatim to the captain (direct text output, not SendMessage):

> Falling back to bare mode for the remainder of this session due to team-infrastructure failure. Prior team agents are presumed-zombified; I will not route work to them or through the team registry. If you want to escalate: restart the session to retry team mode with a fresh name, or let me continue — every stage will still complete, just without concurrent dispatch.

### Cooperative Shutdown Sweep

On Degraded Mode entry, perform a single-pass cooperative shutdown sweep of every known agent name from session memory: one `SendMessage(to="{ensign_name}", message="shutdown_request")` per name. Ignore failures — best-effort, not transactional. Do not retry, track responses, or block on the outcome; proceed immediately to the first fresh bare-mode dispatch.

Exempt any agent whose entity is in an active feedback-cycle state (tracked via a `### Feedback Cycles` subsection in the entity body; read from the worktree copy when `worktree:` is set on the entity, otherwise from main). Those reviewers may hold load-bearing context from the prior cycle that re-dispatch cannot reconstruct. Sweep feedback-cycle reviewers only on explicit captain confirmation.

## Awaiting Completion

After dispatching an ensign (or routing work to a kept-alive ensign), you are waiting for that ensign's completion signal. Until that signal arrives, take NO action that affects the ensign's lifecycle.

**A completion signal is one of these three things, and nothing else:**

1. An inbox-delivered user-role message from the ensign whose text begins with `Done:` (per the ensign runtime's completion contract).
2. A `system` entry with `subtype: task_notification` and `status: completed` whose `tool_use_id` matches the ensign's `Agent(...)` dispatch id.
3. An explicit captain instruction (captain-role user message) to shut down the ensign.

**First-turn-after-dispatch decision procedure.** When a turn begins and your most recent dispatch-related action was an `Agent(...)` spawn whose completion signal (1, 2, or 3 above) has NOT yet been observed in the stream, you MUST end the turn immediately with no tool calls and no text. Do not:

- emit `SendMessage(to="{ensign}", message={"type":"shutdown_request"})` — this is the exact bug this section exists to prevent.
- emit `TeamDelete` — before a completion signal the entity is not terminal, so tearing down is premature; it fails anyway while members are alive, and retrying it in a loop here is the second-order bug. (At the TERMINAL boundary the opposite holds — see `## Terminal Team Teardown`, where `TeamDelete` is attempted as a bounded best-effort.)
- emit `Bash` with commands like `sleep 30` or `wait` — the runtime handles the wait for you; sleeping in Bash wastes time and does not accelerate delivery.
- re-dispatch a replacement ensign — you have no evidence the first ensign failed.
- write reassuring text like "Waiting for completion signal" — this converts idle-polling into a multi-turn generation loop that drifts into hallucination on subsequent wake-ups.

Just emit `end_turn` with empty content. The runtime will wake you up again when a real event arrives.

**A new `system init` entry in the stream is NOT a completion signal.** It is a turn boundary from claude-code's internal event loop (the runtime re-invokes you when idle-poll timers fire or when a teammate event is queued). If you wake up on a fresh `system init` and the prior turn's last observable state was a spawn-ack or a pending dispatch, treat it as idle and end the turn silently per the decision procedure above.

**Anti-patterns that indicate this bug.** If you catch yourself about to emit any of these, STOP and end the turn empty:

- `shutdown_request` with reason `"session ending"`, `"wrapping up"`, `"timeout"`, or any other self-generated reason when no completion signal has arrived. The runtime does NOT signal session-end via your context; it signals it via an actual user message.
- `shutdown_request` followed immediately by `TeamDelete` (the classic premature-teardown loop).
- Any action whose justification is "enough time has passed" or "the ensign appears idle" — you cannot measure time from inside a turn, and ensign idleness is normal between dispatch and completion.

**DISPATCH IDLE GUARDRAIL.** After dispatching an agent, wait for an explicit completion message. Idle notifications are normal between-turn state for team agents — they are not a reason to tear down the team, and they usually mean the agent is waiting for input. Only shut down when: (1) the agent sends a completion message, (2) the captain explicitly requests shutdown, or (3) you are transitioning the entity to a new stage (AFTER you have observed the prior stage's completion signal per the list above). Never interpret idle notifications as "stuck" or "unresponsive."

## Terminal Team Teardown

This governs the TERMINAL phase only — AFTER the entity reached its terminal stage and the FO is dismantling the team. It is a DIFFERENT phase from `## Awaiting Completion` above: that section bans emitting `TeamDelete` *before* a completion signal (the premature-teardown bug); this section *attempts* `TeamDelete` as a bounded best-effort *after* terminal cleanup has begun. Do not conflate the two — the ban is pre-completion, the bounded attempt is terminal.

When tearing down the team at the terminal boundary:

1. Send the cooperative `SendMessage({"type":"shutdown_request"})` to every roster member in the entity's cohort. This is cooperative — the member acknowledges and leaves the roster asynchronously.
2. Call `TeamDelete`. The cooperative shutdown and `TeamDelete` race: the first `TeamDelete` can fail with `Cannot cleanup team with N active member(s): {name}` because a member you just signalled is still settling out of the team registry.
3. **On that failure, do NOT end the turn, and do NOT immediately re-fire `TeamDelete`.** The member is leaving the roster asynchronously, so a retry fired in the same instant just re-loses the same race — this is exactly how #275 failed (it "retried 3× but all 3 raced before the registry settled, then stopped"). Instead, between attempts you MUST let the registry settle: re-send `SendMessage({"type":"shutdown_request"})` to each still-named active member from the failure message, then **wait for the settle before the next `TeamDelete`** — `Bash("sleep 2")` is the settle (this is the one place a deliberate inter-attempt wait is correct; it is NOT the banned idle-poll sleep of `## Awaiting Completion`, which fires with no pending teardown). Then call `TeamDelete` again. Repeat the settle-then-attempt serially up to a small **attempt cap** (≈5 attempts, each preceded by the ~2s settle). Attempts MUST be serial with a settle between them, never fired back-to-back in one message. In an interactive session the roster clears as the member's session-end propagates, so a `TeamDelete` lands cleanly on an early attempt and the loop exits there — the cap is never reached.
4. **On cap-exhaustion in a non-interactive session, STOP attempting and emit the terminal-status marker — the launcher owns the exit.** An approved-shutdown member can stay listed in the team `members[]` indefinitely (an upstream Claude Code defect: the dead member is never cleared from the roster), so `TeamDelete` keeps failing `active member(s)` against a member that has already terminated. A `claude -p` FO CANNOT exit while `members[]` is non-empty (the harness re-invokes it on `end_turn`), so an FO-driven process exit is impossible, and `retry to success` is unreachable — both re-hang the subprocess exactly like #275. So once the **attempt cap** is reached, the FO STOPS calling `TeamDelete` and emits the verbatim terminal-status marker — `TERMINAL_TEARDOWN_BOUNDED: best-effort teardown exhausted; member(s) stuck in registry; holding for launcher.`. The PROCESS EXIT is the **launcher's** responsibility — the live-e2e cycle's `kill()`, or a real automation's launcher-side timeout — never an FO self-exit. On a subsequent harness re-invocation with the roster still non-empty the FO again runs the bounded best-effort and re-emits the marker; a bounded resume that re-emits the marker is acceptable (the launcher ends the subprocess) — what is forbidden is an UNBOUNDED retry loop that never reaches the marker. In an interactive session the cap is not reached (step 3); if it ever were, also surface the still-active member(s) to the captain alongside the marker.

Emitting the marker at the cap IS the fix: a teardown that gives up silently on the first `active member(s)` failure with NO marker (#275), or one that retries `TeamDelete` unboundedly past the cap and never reaches the marker (#282's live failure), both strand the `claude -p` FO subprocess waiting on an exit the FO cannot cause. The inter-attempt settle is what gives an interactive roster time to clear so `TeamDelete` lands; the **attempt cap** bounds the FAST back-to-back attempts so the non-interactive dead-but-listed member does not spin an unbounded retry loop; the marker is the deliberate terminal status the launcher's `kill()` ends the subprocess on. A failed `TeamDelete` here is the EXPECTED first step of a settling race, not a reason to give up silently; this is the one place the FO actively attempts `TeamDelete` (everywhere else — see `## Awaiting Completion` — emitting it is the bug).
