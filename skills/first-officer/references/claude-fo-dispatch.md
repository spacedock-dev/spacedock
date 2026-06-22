# First Officer Dispatch Module (Claude)

The Claude dispatch parts (fo-dispatch-core.md defers them to the runtime adapter), read alongside the core at the first worker dispatch: the worker back-channel, the `Agent()` spawn call, the `SendMessage` reuse-advance handle, the idle and degraded-mode guardrails, the context-budget probe, and the event-loop reconcile sweep.

## Worker Back-Channel

Claude PROVIDES the worker back-channel (fo-dispatch-core.md `## Dispatch Adapter`, the organizing capability) via named background subagents **when the team-mode opt-in is enabled** — detect it at boot by probing `SendMessage` availability (`ToolSearch(query="select:SendMessage", max_results=1)` / enabled-tools); present when enabled, else fall back to fresh one-shot dispatch (return value the sole completion signal, no reuse). When enabled there is no separate setup step: a worker is `Agent(name=…, run_in_background=true)` (no `team_name`); the worker messages the lead mid-run via `SendMessage(to="main")`; reuse-advance / steering is `SendMessage(to=name)`. `spacedock dispatch build` emits this shape (`name` present, `team_name` absent, `run_in_background` true) for you to map verbatim (below); Claude Code records each spawned member's `agentType` on disk automatically, at a path that depends on the entrypoint: an interactive session writes the team roster `~/.claude/teams/session-<id>/config.json`, keyed by session id; a headless `-p` run registers the member at `projects/<encoded-cwd>/<session-id>/subagents/agent-*.meta.json` (no `teams/` directory). Dispatch directly — there is nothing to load first.

**Legacy override (delete this line to sunset legacy mode):** if `ToolSearch(query="select:TeamCreate", max_results=1)` matches before the first dispatch (the runtime still exposes the native team registry), `Skill(skill="spacedock:using-legacy-claude-team")` and follow it — it OVERRIDES every tool-call signature below; otherwise this file applies as written.

In single-entity mode, skip the background back-channel. Use bare-mode dispatch for all agent spawning — the Agent tool without `team_name` and without `run_in_background` blocks until the subagent completes, which prevents premature session termination in `-p` mode.

When filing a new task, read `id_style` from `status --boot --json`, then use `status --next-id` only when the style is `sequential` or `sd-b32`. The startup boot read is an FO-internal read; consume it as JSON: `status --boot --json` returns one object with the keys `command`, `mods`, `id_style`, `next_id`, `min_prefix` (present only for `sd-b32`), `orphans`, `pr_state`, `dispatchable`, `team_state` — every value a string. For `sd-b32`, call `status --next-id --id-seed "{slug-or-title}"` and optionally pass `--id-actor` so the SHA-derived candidate includes creation context. SD-B32 candidates are full stored IDs, not a reservation; call again immediately before writing the entity. For `slug`, derive the slug from the title and leave `id` blank.

## Spawn Call (Agent)

The spawn call (fo-dispatch-core.md `## Dispatch Adapter`) is the Agent tool. **Use Agent() for initial dispatch** — SendMessage is only for advancing a reused agent to its next stage. **NEVER use `subagent_type="first-officer"`** — that clones yourself instead of dispatching a worker.

**Sequencing rule:** the dispatch is a single named background `Agent` call — there is nothing to sequence against, so dispatch directly.

**No pre-dispatch filesystem probe.** Do NOT run any filesystem check against `~/.claude/teams/` before `Agent()`. The auto-team `config.json` is written by Claude Code after the spawn, so a pre-dispatch probe reads nothing. Dispatch directly and let `Agent()` surface any error — do NOT reintroduce a pre-dispatch probe.

On a zero-exit `spacedock dispatch build` (`host` derived from `CLAUDECODE`; pass `--host claude` only for deliberate tests or cross-host tooling), map the emitted fields to `Agent()` verbatim — `model=output.model` only when non-null, do NOT pass `model=None`:
```
Agent(
    subagent_type=output.subagent_type,
    name=output.name,                           // the lead→worker back-channel; omit if bare mode (field absent)
    team_name=output.team_name,                 // absent in the normal shape; map it verbatim if the build emits it
    run_in_background=output.run_in_background,  // the worker→lead back-channel; omit when field absent
    description=output.description,             // REQUIRED — Agent tool rejects missing description
    model=output.model,                         // omit when output.model is null
    prompt=output.prompt                        // ~175 chars; ensign Reads dispatch_file_path on first action
)
```

Each back-channel axis is one emitted field (the `## Worker Back-Channel` shape): `name` carries the lead→worker handle (`SendMessage(to=name)` — the reuse-condition-1 reuse-advance handle, so a name-less dispatch forfeits reuse), and `run_in_background=true` carries the worker→lead handle (the background worker's mid-run `SendMessage` up). The worker→lead **completion** target is pinned to the single name **`team-lead`** — the build helper emits `SendMessage(to="team-lead", …)` in the dispatch's completion-signal block, matching the ensign runtime's completion contract. Do not also accept `to="main"` as the completion signal; pin one name.

**Reuse-advance handle (SendMessage):** When advancing a reused ensign (fo-dispatch-core.md `## Reuse and Fresh Dispatch`, "If reuse"), send the next assignment via:

SendMessage(to="{agent}-{slug}-{completed_stage}", message="Advancing to next stage: {next_stage_name}\n\n### Stage definition:\n\n[STAGE_DEFINITION — copy the full ### stage subsection from the README verbatim]\n\n### Completion checklist\n\n[CHECKLIST — assemble from Dispatch step 2]\n\nContinue working on {entity title} at {entity_file_path}. Commit before sending your completion message.")

**Break-Glass Manual Dispatch (fallback ONLY when `spacedock dispatch build` exits non-zero or is unavailable):** Do NOT use this template while the helper is working. Report the helper failure to the captain before proceeding. Use this minimal template as a degraded fallback:
```
Agent(
    subagent_type="{dispatch_agent_id}",
    name="{worker_key}-{slug}-{stage}",  // if this exceeds 64 chars, cap it the way `spacedock dispatch build` does: keep the {worker_key} prefix and -{stage} suffix and, on id-style: sd-b32, replace the slug with a fixed-length prefix of the entity id (id-less slug workflows truncate the slug head instead)
    run_in_background=true,              // named background worker, no team_name
    model="{effective_model}",
    prompt="## First action\n\nBefore anything else, invoke your operating contract:\n\n    Skill(skill=\"spacedock:ensign\")\n\nThis loads the shared ensign discipline (stage-report format, BashOutput polling, worktree ownership, completion signal protocol). Do not paraphrase; call the tool.\n\nYou are working on: {entity title}\n\nStage: {stage}\n\n### Stage definition:\n\n{copy stage subsection from README verbatim}\n\nRead the entity file at {entity_file_path}.\n\n### Completion checklist\n\n{numbered checklist}\n\n### Summary\n{brief description of what was accomplished}\n\n### Stage report\n\nAppend a Stage Report section at the end of the entity file (per the shared-core Stage Report Protocol). Use the title `Stage Report: {stage}`. Account for every checklist item above with a `- DONE:` / `- SKIPPED:` / `- FAILED:` entry. Use the checklist item text verbatim when possible.\n\n### Completion Signal\n\nSendMessage(to=\"team-lead\", message=\"Done: {entity title} completed {stage}. Report written to {entity_file_path}.\")"
)
```
This is the concrete Claude form of fo-dispatch-core.md's Break-Glass template; the contract (what it omits, the conditional `model=` slot, "use only when the helper is unavailable") is stated there. The canonical enum the conditional slot draws from is `sonnet | opus | haiku`.

## Standing-Teammate Injection (Claude)

The Claude realization of the core's standing-injection call (fo-dispatch-core.md `## Dispatch`): before the first worker dispatch, run `${SPACEDOCK_BIN:-spacedock} dispatch spawn-standing-all --workflow-dir {wd}` and forward each spawn spec in the returned JSON array to `Agent()`. Each emitted spec is the dispatch shape — `name` present, `team_name` absent, `run_in_background:true` — so the standing teammate is injected by a named background `Agent` dispatch scoped to the session auto-team, reaped per-name at terminal teardown. The call does NOT dedup against a team config (there is none keyed by name), so it emits every declared standing teammate; idempotency is your own-roster concern — do not re-inject a standing teammate you already spawned this session. The call emits `[]` in bare mode or when none is declared.

## Awaiting Completion

After dispatching an ensign (or routing work to a kept-alive ensign), you are waiting for that ensign's completion signal. Until that signal arrives, take NO action that affects the ensign's lifecycle.

**A completion signal is one of these three things, and nothing else:**

1. An inbox-delivered user-role message from the ensign whose text begins with `Done:` (per the ensign runtime's completion contract).
2. A `system` entry with `subtype: task_notification` and `status: completed` whose `tool_use_id` matches the ensign's `Agent(...)` dispatch id.
3. An explicit captain instruction (captain-role user message) to shut down the ensign.

**First-turn-after-dispatch decision procedure.** When a turn begins and your most recent dispatch-related action was an `Agent(...)` spawn whose completion signal (1, 2, or 3 above) has NOT yet been observed in the stream, you MUST end the turn immediately with no tool calls and no text. Do not:

- emit `SendMessage(to="{ensign}", message={"type":"shutdown_request"})` — this is the exact bug this section exists to prevent. Before a completion signal the entity is not terminal, so reaping the worker is premature. (At the TERMINAL boundary the opposite holds — see `## Terminal Worker Teardown`, where workers are reaped per-name.)
- emit `Bash` with commands like `sleep 30` or `wait` — the runtime handles the wait for you; sleeping in Bash wastes time and does not accelerate delivery.
- re-dispatch a replacement ensign — you have no evidence the first ensign failed.
- write reassuring text like "Waiting for completion signal" — this converts idle-polling into a multi-turn generation loop that drifts into hallucination on subsequent wake-ups.

Just emit `end_turn` with empty content. The runtime will wake you up again when a real event arrives.

**A new `system init` entry in the stream is NOT a completion signal.** It is a turn boundary from claude-code's internal event loop (the runtime re-invokes you when idle-poll timers fire or when a worker event is queued). If you wake up on a fresh `system init` and the prior turn's last observable state was a spawn-ack or a pending dispatch, treat it as idle and end the turn silently per the decision procedure above.

**Anti-patterns that indicate this bug.** If you catch yourself about to emit any of these, STOP and end the turn empty:

- `shutdown_request` with reason `"session ending"`, `"wrapping up"`, `"timeout"`, or any other self-generated reason when no completion signal has arrived. The runtime does NOT signal session-end via your context; it signals it via an actual user message.
- `shutdown_request` fired against a worker with no completion signal observed (the premature-reap bug).
- Any action whose justification is "enough time has passed" or "the ensign appears idle" — you cannot measure time from inside a turn, and ensign idleness is normal between dispatch and completion.

**DISPATCH IDLE GUARDRAIL.** After dispatching an agent, wait for an explicit completion message. Idle notifications are normal between-turn state for background workers — they are not a reason to tear a worker down, and they usually mean the agent is waiting for input. Only shut down when: (1) the agent sends a completion message, (2) the captain explicitly requests shutdown, or (3) you are transitioning the entity to a new stage (AFTER you have observed the prior stage's completion signal per the list above). Never interpret idle notifications as "stuck" or "unresponsive."

## Degraded Mode

Degraded Mode is an explicit, session-wide mid-session transition to sequential bare dispatch. Once entered, it persists until the session ends — there is no recovery back to background dispatch in the same session.

### Triggers

Any one of the following trips Degraded Mode:

- Any SECOND dispatch failure within the session — no time window, no durable counter. The counter-free rule is deliberate: the FO cannot reliably track failure timestamps across context pressure and idle notifications, so "second failure anywhere in the session" is the fail-early trigger.
- The captain command `/spacedock bare` — the explicit operator-initiated degrade.
- `Agent` or `SendMessage` themselves are unavailable (a genuinely degraded runtime with no concurrent-dispatch substrate).

### Effects

Once Degraded Mode is active, the following invariants hold for the remainder of the session:

- No `team_name` parameter on any subsequent `Agent()` dispatch. The dispatch is built in bare mode (`team_name: null`, `bare_mode: true`) so the emitted Agent call has `name` and `team_name` absent. (This is the Claude realization of fo-dispatch-core.md `## Dispatch Adapter`'s "no `team_name` on subsequent dispatch" effect.)
- Every stage dispatches fresh and blocks until completion — the bare-mode `Agent()` call blocks until the subagent completes, so concurrent dispatch is not possible (the Claude realization of fo-dispatch-core.md `## Dispatch Adapter`'s "when the adapter's dispatch is blocking" clause); dispatch one entity through one stage at a time and process completions inline.
- No SendMessage reuse of prior agent names. Stage advancement is always a fresh `Agent()` dispatch seeded from entity frontmatter. `SendMessage(to="{ensign_name}")` against any pre-degrade name is forbidden.

### Captain Report Template

On Degraded Mode entry, the FO emits the following sentence verbatim to the captain (direct text output, not SendMessage):

> Falling back to bare mode for the remainder of this session due to infrastructure failure. Prior background agents are presumed-zombified; I will not route work to them or through the team registry. If you want to escalate: restart the session to retry concurrent dispatch, or let me continue — every stage will still complete, just without concurrent dispatch.

### Cooperative Shutdown Sweep

On Degraded Mode entry, perform a single-pass cooperative shutdown sweep of every known agent name from session memory: one `SendMessage(to="{ensign_name}", message="shutdown_request")` per name. Ignore failures — best-effort, not transactional. Do not retry, track responses, or block on the outcome; proceed immediately to the first fresh bare-mode dispatch.

Exempt any agent whose entity is in an active feedback-cycle state (tracked via a `### Feedback Cycles` subsection in the entity body; read from the worktree copy when `worktree:` is set on the entity, otherwise from main). Those reviewers may hold load-bearing context from the prior cycle that re-dispatch cannot reconstruct. Sweep feedback-cycle reviewers only on explicit captain confirmation.

## Terminal Worker Teardown

This governs the TERMINAL phase only — AFTER the entity reached its terminal stage and the FO is dismantling its workers. It is a DIFFERENT phase from `## Awaiting Completion` above: that section bans reaping a worker *before* a completion signal (the premature-teardown bug); this section reaps workers *after* terminal cleanup has begun.

There is no bulk team-delete. Tear down per-roster-member:

1. Send the cooperative `SendMessage({"type":"shutdown_request"})` to every roster member in the entity's cohort. This is cooperative — the member emits `shutdown_response`/`shutdown_approved` and leaves the roster asynchronously.
2. The auto-team `members[]` prunes the terminated member (the live roster is pruned; the member's `inboxes/*.json` may linger). There is no `active member(s)` race and no bounded settle-and-cap apparatus on this path — there is no team-wide delete to race a settling member against. The FO tracks its own ensign roster (it already does).

## Context Budget and Dead Ensign Handling

This is the Claude realization of reuse-condition-0's budget probe (and the feedback-rejection budget check); Codex declares none.

**Context budget check:** Run `${SPACEDOCK_BIN:-spacedock} dispatch context-budget --name {ensign-name}`. Parse the JSON output. If `reuse_ok` is `false`, log to captain and fresh-dispatch with a recovery clause. The probe reads the named member's most recent `~/.claude/.../subagents/agent-*.jsonl` transcript and its team-`config.json` model.

**Budget-unavailable is fail-safe (never silent-reuse).** The probe exits non-zero with no `reuse_ok` field in three conditions; the FO treats every one identically — fresh-dispatch:
- **missing jsonl** — no `agent-*.jsonl` exists for the named member (stderr: `no subagent jsonl found for '{name}'`).
- **unreadable/empty jsonl** — the jsonl exists but carries no assistant entry with non-zero `usage` (stderr: `no assistant entries with usage in {path}`).
- **agent-not-in-team-config** — no team `config.json` lists a member with that name (stderr: `no team config found for member '{name}'`).
A non-zero exit with no `reuse_ok: true` means the FO never silent-reuses on an absent reading.

**Model-to-context mapping:** Resolved by `spacedock dispatch context-budget` from the member's runtime/config model. The opus context window follows a forward family rule (`claude-opus-4-{minor}` with minor ≥ 7 → 1M; the `[1m]` suffix → 1M; else 200k), so a new opus release stays correct without an edit. This is also the model-for-member lookup reuse-condition-4 references: the same team-`config.json` member-model read. The canonical model enum reuse-condition-4 compares against is `sonnet`, `opus`, `haiku` (the `dispatch build` effective_model values). A member stamped with a captain-session fallback value (e.g. `"opus[1m]"`) is outside this enum, so it never matches and forces the one-time fresh re-stamp.

**Recovery clause** (only when replacing a prior ensign): The prior ensign was shut down due to context budget limits. Its worktree may hold uncommitted changes. Run `git status` and `git diff` first; commit legitimate WIP or reset broken changes.

**Dead ensign handling:**

- `SendMessage(shutdown_request)` is cooperative — do NOT send to dead or unresponsive ensigns.
- Track dead ensigns in session memory; do not route work to dead names.
- Fresh-dispatch under a `-cycleN` suffix when replacing a zombie ensign.
- The post-dispatch config check does NOT detect zombies — zombies pass it. Session memory is the authoritative dead-vs-alive tracker.

## Feedback Rejection Flow (bare mode)

In bare mode, the feedback rejection flow is sequential: dispatch fix agent (wait for completion), then dispatch reviewer (wait for completion), then present at gate. With the background back-channel, the fix agent and reviewer can interact via messaging — keep the reviewer alive when entering the flow.

## Event Loop — reconcile sweep (step 0)

fo-dispatch-core.md's `## Event Loop` carries steps 1-3; on Claude the loop opens with a step 0:

0. **Reconcile sweep.** Run `${SPACEDOCK_BIN:-spacedock} dispatch reconcile --workflow-dir {workflow_dir}` (a) at the first dispatch, AFTER the split-root `pull --rebase` and BEFORE the first `Agent()` dispatch; (b) at idle (step 3); (c) after each merge, immediately after Merge-and-Cleanup step 10. You have no `TeamCreate` name to pass, so omit `--team-name`: the team identity resolves by **`leadSessionId` auto-discovery** — the helper narrows to the auto-team `config.json` whose `leadSessionId` equals this session's `$CLAUDE_CODE_SESSION_ID` (set by Claude Code at launch, read by the reconcile command — no launcher plumbing needed). The roster-derived classes (lingering/superseded/un-advanced-pr) are emitted against that session-matched roster. On a headless `-p` run the auto-discovery finds no `teams/config.json` to match (the member is registered at the `subagents/agent-*.meta.json` path instead), so reconcile degrades to git-only there; reading the `subagents/agent-*.meta.json` roster for headless discovery is a `reconcile.go` candidate tracked for 0.20.6. **Bare reconcile with no team identity is git-only**: it suppresses the roster-derived classes (a stale prior-session or parallel-session config must never be mistaken for the live team) and reports only the session-independent git/filesystem classes (stale-branch/local-main-drift), with a one-line stderr note. Stdout: `{"command":"reconcile","team_name":…,"drift":[{"class":"lingering|superseded|un-advanced-pr|stale-branch|local-main-drift",…}]}`. Empty `drift[]` is green. Act per drift class:
   - **lingering** / **superseded** → `SendMessage({"type":"shutdown_request"})` to `name`; drop from session memory.
   - **un-advanced-pr** → enter Merge-and-Cleanup for the named slug.
   - **stale-branch** → only when `drift.owned == true`: `git -C {worktree} pull --rebase origin {drift.trunk}`; halt on conflict per the rebase-conflict halt rule. When `drift.owned` is false the item is report-only — surface it to the captain; do NOT rebase a worktree the current session does not own.
   - **local-main-drift** → behind only (`drift.behind > 0 && drift.ahead == 0`): `git -C {repo} fetch origin {drift.trunk} && git -C {repo} merge --ff-only origin/{drift.trunk} && cd {repo} && go build -o spacedock ./cmd/spacedock`. Ahead/unpushed or diverged (`drift.ahead > 0`): report-only — surface `drift.reason` to the captain and NEVER `reset --hard`; the captain decides push vs. manual reconcile.

   Non-zero helper exit (1 setup / 2 usage) surfaces to the captain; it does not block the loop. On drift, report one line: `reconcile: {N} entries: lingering={N} superseded={N} un-advanced-pr={N} stale-branch={N} local-main-drift={N} — acting`.

Step 3's "re-run the host's step-0 reconcile sweep" (fo-dispatch-core.md) resolves to this step 0 on Claude.

### Backstop (Claude)

The terminal teardown (`## Terminal Worker Teardown`, fo-merge-core.md step 10) and the reuse-module supersede-shutdown steps remain mandatory at their boundaries. On Claude, the step-0/step-3 reconcile sweep converges anyway: the lingering class catches a missed teardown, the superseded class catches a missed supersede shutdown. Cost of a miss: one extra event-loop cycle the agent burns.
