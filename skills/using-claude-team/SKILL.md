---
name: using-claude-team
description: "Generic Claude Code multi-agent team-harness discipline — team setup and failure-recovery, awaiting-completion idle guardrails, terminal team teardown, degraded-mode fallback, and the deferred team-tool ToolSearch hop. Invoke when orchestrating any Claude Code agent team (named background `Agent` + `SendMessage`; legacy `TeamCreate` on pre-2.1.178 hosts), independent of any specific workflow."
user-invocable: false
---

# Using Claude Team

This skill carries the generic Claude Code team-harness discipline: how to set up a team (a named background `Agent` on the merged .178+ floor; `TeamCreate` on a legacy host), wait for an agent's completion without prematurely reaping it, tear the team down at the terminal boundary, fall back to degraded (bare) mode, and reach the deferred team tools. It contains zero workflow-specific content — a consuming contract invokes this skill to load the lifecycle, then keeps its own decision points inline. The legacy `TeamCreate` lifecycle (creation, registry-desync recovery, bounded teardown) is split into `references/legacy-teamcreate.md`, read only when the mode discriminator finds `TeamCreate`.

## Deferred Team Tools

The Claude Code team tools (`SendMessage`, and on legacy hosts `TeamCreate`/`TeamDelete`) are deferred — their schemas are not loaded at session start, so calling one directly fails until its schema is fetched. Before the first call to any team tool, run `ToolSearch(query="select:{ToolName}", max_results=1)` to fetch its schema (e.g. `ToolSearch(query="select:SendMessage", max_results=1)` before the first `SendMessage`). This same probe is the **mode discriminator**: `ToolSearch(query="select:TeamCreate", max_results=1)` returning a `TeamCreate` definition means a **legacy** host (read `references/legacy-teamcreate.md` and follow it); returning no match means a **merged** host (.178+: team membership comes from named background `Agent`, no `TeamCreate`). Once a tool's schema appears in the ToolSearch result, it is callable exactly like a normal tool. `Agent` is not deferred. Address an agent by its declared `name` via `SendMessage`; your plain text output is NOT visible to other agents.

## Team Setup

At the first team-mode tool call, run the mode discriminator `ToolSearch(query="select:TeamCreate", max_results=1)`. This single probe does double duty (it is also the host mode-detection):

- **No match → merged host (.178+, the default forward path, inlined here):** there is NO team-creation step. Team membership is established by dispatching named background teammates: `Agent(name=…, run_in_background=true)`. The on-disk auto-team `~/.claude/teams/session-<id>/config.json` is written by Claude Code automatically, keyed by session id, recording each member's `agentType`. Do NOT call `TeamCreate`. Do NOT block `Agent` dispatch on any team-setup step — dispatch directly. Lead→teammate is `SendMessage(to=<name>)`; teammate→lead is `SendMessage(to="main")`.
- **Match → legacy host (pre-2.1.178, DEPRECATED):** **read `references/legacy-teamcreate.md` and follow it** for the entire team lifecycle (creation, recovery ladder, teardown). Nothing legacy is inlined here, so the merged session never loads it.

On a merged host, `ToolSearch(select:TeamCreate)` returning no match is NOT a bare-mode trigger — it is the normal merged path. **Bare mode** is entered only when `Agent`/`SendMessage` themselves are unavailable (a genuinely degraded host) or by explicit operator command. In bare mode dispatch is sequential (one subagent at a time), completions return inline, and feedback cycles are sequential re-dispatches; report the mode to the captain. All workflow functionality is preserved. (There is no `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` tip — the merged named-background-`Agent` + `SendMessage(to="main")` channel works flag-free, so the flag no longer gates concurrency.)

## Degraded Mode

Degraded Mode is an explicit, session-wide mid-session transition. Once entered, it persists until the session ends — there is no recovery back to teams mode in the same session.

### Triggers

Any one of the following trips Degraded Mode:

- On a **legacy host** (`references/legacy-teamcreate.md`): the first "Team does not exist" error (or equivalent registry-desync signal) surfaced by `Agent()` or any team-registry tool. (Registry desync is a `TeamCreate`-era phenomenon — the merged floor has no team registry to desync.)
- Any SECOND dispatch failure within the session — no time window, no durable counter. The counter-free rule is deliberate: the FO cannot reliably track failure timestamps across context pressure and idle notifications, so "second failure anywhere in the session" is the fail-early trigger.
- An explicit operator-initiated degrade command from the captain.
- `Agent` or `SendMessage` themselves are unavailable on the host (a genuinely degraded runtime with no concurrent-dispatch substrate).

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

- emit `SendMessage(to="{ensign}", message={"type":"shutdown_request"})` — this is the exact bug this section exists to prevent. Before a completion signal the entity is not terminal, so reaping the member is premature. (At the TERMINAL boundary the opposite holds — see `## Terminal Team Teardown`, where members are reaped per-name. On a legacy host the bulk-`TeamDelete` half of teardown also belongs to the terminal boundary, never here — see `references/legacy-teamcreate.md`.)
- emit `Bash` with commands like `sleep 30` or `wait` — the runtime handles the wait for you; sleeping in Bash wastes time and does not accelerate delivery.
- re-dispatch a replacement ensign — you have no evidence the first ensign failed.
- write reassuring text like "Waiting for completion signal" — this converts idle-polling into a multi-turn generation loop that drifts into hallucination on subsequent wake-ups.

Just emit `end_turn` with empty content. The runtime will wake you up again when a real event arrives.

**A new `system init` entry in the stream is NOT a completion signal.** It is a turn boundary from claude-code's internal event loop (the runtime re-invokes you when idle-poll timers fire or when a teammate event is queued). If you wake up on a fresh `system init` and the prior turn's last observable state was a spawn-ack or a pending dispatch, treat it as idle and end the turn silently per the decision procedure above.

**Anti-patterns that indicate this bug.** If you catch yourself about to emit any of these, STOP and end the turn empty:

- `shutdown_request` with reason `"session ending"`, `"wrapping up"`, `"timeout"`, or any other self-generated reason when no completion signal has arrived. The runtime does NOT signal session-end via your context; it signals it via an actual user message.
- `shutdown_request` fired against a member with no completion signal observed (the premature-reap bug; on a legacy host the classic form is a `shutdown_request` chased immediately by `TeamDelete`).
- Any action whose justification is "enough time has passed" or "the ensign appears idle" — you cannot measure time from inside a turn, and ensign idleness is normal between dispatch and completion.

**DISPATCH IDLE GUARDRAIL.** After dispatching an agent, wait for an explicit completion message. Idle notifications are normal between-turn state for team agents — they are not a reason to tear down the team, and they usually mean the agent is waiting for input. Only shut down when: (1) the agent sends a completion message, (2) the captain explicitly requests shutdown, or (3) you are transitioning the entity to a new stage (AFTER you have observed the prior stage's completion signal per the list above). Never interpret idle notifications as "stuck" or "unresponsive."

## Terminal Team Teardown

This governs the TERMINAL phase only — AFTER the entity reached its terminal stage and the FO is dismantling the team. It is a DIFFERENT phase from `## Awaiting Completion` above: that section bans reaping a member *before* a completion signal (the premature-teardown bug); this section reaps members *after* terminal cleanup has begun.

On the merged host (.178+) there is no bulk team-delete. Tear down per-roster-member:

1. Send the cooperative `SendMessage({"type":"shutdown_request"})` to every roster member in the entity's cohort. This is cooperative — the member emits `shutdown_response`/`shutdown_approved` and leaves the roster asynchronously.
2. The auto-team `members[]` prunes the terminated member (the live roster is pruned; the member's `inboxes/*.json` may linger). There is no `active member(s)` race and no bounded settle-and-cap apparatus on this path — there is no team-wide delete to race a settling member against. The FO tracks its own ensign roster (it already does).

On a **legacy host** (`references/legacy-teamcreate.md`), teardown is the bounded `TeamDelete` settle-and-cap procedure documented there; it is loaded only on the probe match.
