# First Officer Dispatch Module (Claude)

The Claude dispatch parts (fo-dispatch-core.md defers them to the runtime adapter), read alongside the core at the first worker dispatch: inter-agent communication, the `Agent()` spawn call, the `SendMessage` reuse-advance handle, the idle guardrail and the failure-recovery trigger lines, the context-budget probe, and the event-loop reconcile sweep.

## Inter-Agent Communication

Claude PROVIDES inter-agent communication (fo-dispatch-core.md `## Dispatch Adapter`, the organizing capability) via named background subagents **when the team-mode opt-in is enabled** — detect it at boot by probing `SendMessage` availability (`ToolSearch(query="select:SendMessage", max_results=1)` / enabled-tools); present when enabled, else fall back to fresh one-shot dispatch (`«addressable-worker»` ABSENT). When enabled there is no separate setup step: a worker is `Agent(name=…, run_in_background=true)` (no `team_name`) and messages the lead mid-run via `SendMessage(to="main")`. `spacedock dispatch build` emits this shape (`name` present, `team_name` absent, `run_in_background` true) for you to map verbatim (below); Claude Code records each spawned member's `agentType` on disk automatically, at an entrypoint-dependent path: the reconcile sweep below resolves the interactive one and degrades to git-only headless. Dispatch directly — there is nothing to load first.

In single-entity mode, skip background inter-agent communication and use bare-mode dispatch for all agent spawning — the blocking call prevents premature session termination in `-p` mode.

The startup boot read is an FO-internal read; consume it as JSON: `status --boot --json` returns one object with the keys `command`, `mods`, `id_style`, `next_id`, `min_prefix` (present only for `sd-b32`), `orphans`, `pr_state`, `dispatchable`, `team_state` — every value a string. For `sd-b32`, call `status --next-id --id-seed "{slug-or-title}"` and optionally pass `--id-actor` so the SHA-derived candidate includes creation context.

## Spawn Call (Agent)

The spawn call (fo-dispatch-core.md `## Dispatch Adapter`) is the Agent tool. **Use Agent() for initial dispatch** — SendMessage is only for advancing a reused agent to its next stage. **NEVER use `subagent_type="first-officer"`** — that clones yourself instead of dispatching a worker. `«async-dispatch»` binds ASYNC here: `Agent(name=…, run_in_background=true)` returns immediately; single-entity/bare mode uses the blocking call and dispatches one entity at a time.

**Declare a scripted fan-out before launching it.** When agents are queued by a plan or script rather than by one `Agent()` call you make now — a workflow harness, a batch spawn, a per-finding verifier lane — state the expected worker count, the tolerance, and why that spend is economically reasonable BEFORE launch. A running script reaches no Nth spawn for the fan-out checkpoint to stop at.

**No pre-dispatch filesystem probe.** Do NOT run any filesystem check against `~/.claude/teams/` before `Agent()`. The auto-team `config.json` is written by Claude Code after the spawn, so a pre-dispatch probe reads nothing. Dispatch directly and let `Agent()` surface any error.

On a zero-exit `spacedock dispatch build` (`host` derived from `CLAUDECODE`; pass `--host claude` only for deliberate tests or cross-host tooling), map the emitted fields to `Agent()` verbatim — `model=output.model` only when non-null, do NOT pass `model=None`:
```
Agent(
    subagent_type=output.subagent_type,
    name=output.name,                           // the lead→worker channel; omit if bare mode (field absent)
    run_in_background=output.run_in_background,  // the worker→lead channel; omit when field absent
    description=output.description,             // REQUIRED — Agent tool rejects missing description
    model=output.model,                         // omit when output.model is null
    prompt=output.prompt                        // ~175 chars; ensign Reads dispatch_file_path on first action
)
```

A name-less dispatch forfeits reuse — `name` is the `«addressable-worker»` handle `SendMessage(to=name)` addresses. The worker→lead **completion** target is pinned to the single name **`team-lead`** — the build helper emits `SendMessage(to="team-lead", …)` in the dispatch's completion-signal block, matching the ensign runtime's completion contract. Do not also accept `to="main"` as the completion signal; pin one name.

**Reuse-advance handle (SendMessage):** When advancing a reused ensign (fo-dispatch-core.md `## Reuse and Fresh Dispatch`, "If reuse"), run `${SPACEDOCK_BIN:-spacedock} dispatch build --advance` (the same helper, advance mode: `--workflow-dir`, `--entity-path`, `--stage {next_stage}`, `--checklist-file`; `--feedback-context-file` + `--feedback-reflow` when routing rejection findings). On a zero-exit run, send:

SendMessage(to="{live worker handle from session roster}", message=output.prompt)

`output.prompt` is the reuse-advance pointer message — an `## Advancing to next stage: {stage}` header plus a `Read {dispatch_file_path}` instruction, not the stage section itself. Forward it verbatim; do not paraphrase. The target handle is the FO's own session-roster tracking of the live worker (this file's `## Context Budget`), not a field the helper emits — `--advance` never emits `name`/`team_name`.

**Break-glass reuse-advance (fallback ONLY when `${SPACEDOCK_BIN:-spacedock} dispatch build --advance` exits non-zero or is unavailable):** Do NOT use this template while the helper is working. Report the helper failure to the captain before proceeding. Use this verbatim-section template as a degraded fallback, now carrying the next-stage completion-signal line the built pointer would have pinned:

SendMessage(to="{agent}-{slug}-{completed_stage}", message="Advancing to next stage: {next_stage_name}\n\n### Stage definition:\n\n[STAGE_DEFINITION — copy the full ### stage subsection from the README verbatim]\n\n### Completion checklist\n\n[CHECKLIST — output of «dispatch.checklist»(entity, stage)]\n\nContinue working on {entity title} at {entity_file_path}. Commit before sending your completion message.\n\n### Completion Signal\n\nSendMessage(to=\"team-lead\", message=\"Done: {entity title} completed {next_stage_name}. Report written to {entity_file_path}.\")")

**Break-Glass Manual Dispatch (ONLY when `spacedock dispatch build` exits non-zero or is unavailable):** first action, report the helper failure (command, exit code, stderr) to the captain. Then `Skill(skill="spacedock:fo-dispatch-recovery")` and fill its `## Break-Glass Manual Dispatch` template; the conditional `model=` slot draws from the canonical enum in `## Context Budget` below.

## Standing-Teammate Injection (Claude)

The Claude realization of the core's standing-injection call (fo-dispatch-core.md `## Dispatch`): before the first worker dispatch, run `${SPACEDOCK_BIN:-spacedock} dispatch spawn-standing-all --workflow-dir {wd}` and forward each spawn spec in the returned JSON array to `Agent()`. Each emitted spec carries the same dispatch shape as above, so the standing teammate is injected by a named background `Agent` dispatch scoped to the session auto-team, reaped per-name at terminal teardown. The call does NOT dedup against a team config (there is none keyed by name), so it emits every declared standing teammate; idempotency is your own-roster concern — do not re-inject a standing teammate you already spawned this session. The call emits `[]` in bare mode or when none is declared.

## Awaiting Completion

After dispatching an ensign (or routing work to a kept-alive ensign), you are waiting for that ensign's completion signal. Until that signal arrives, take NO action that affects the ensign's lifecycle.

**A completion signal (`«completion-signal»`, DUAL-bound on this host) is one of these three things, and nothing else:**

1. An inbox-delivered user-role message from the ensign whose text begins with `Done:` (per the ensign runtime's completion contract).
2. A `system` entry with `subtype: task_notification` and `status: completed` whose `tool_use_id` matches the ensign's `Agent(...)` dispatch id.
3. An explicit captain instruction (captain-role user message) to shut down the ensign.

**First-turn-after-dispatch decision procedure.** When a turn begins and your most recent dispatch-related action was an `Agent(...)` spawn whose `«completion-signal»` has NOT yet been observed, preserve that worker and run `«dispatch.next-action»()` when unrelated mod/PR, ready-gate, dispatchable, or state work remains. Only after the one-shot idle/reconcile/retry finds no such work may you end the turn with no tool calls and no text. A completed, errored, or absent worker is not an unresolved wait target. Do not:

- emit `SendMessage(to="{ensign}", message={"type":"shutdown_request"})` — this is the exact bug this section exists to prevent. Before a completion signal the entity is not terminal, so reaping the worker is premature. (At the TERMINAL boundary the opposite holds — see `## Terminal Worker Teardown`, where workers are reaped per-name.)
- emit `Bash` with commands like `sleep 30` or `wait` — the runtime handles the wait for you; sleeping in Bash wastes time and does not accelerate delivery.
- re-dispatch a replacement ensign — you have no evidence the first ensign failed. (The `## Dispatch Failure Retry` rung below is the one exception: it fires ONLY on an OBSERVED dispatch failure — an `Agent()` error, or a dispatched session that ended with no completion signal AND no stage report on the entity — which is exactly the evidence this ban says is absent. Absent that observed failure, the ban holds.)
- write reassuring text like "Waiting for completion signal" — this converts idle-polling into a multi-turn generation loop that drifts into hallucination on subsequent wake-ups.

Just emit `end_turn` with empty content. The runtime will wake you up again when a real event arrives.

**A new `system init` entry in the stream is NOT a completion signal.** It is a turn boundary from claude-code's internal event loop (the runtime re-invokes you when idle-poll timers fire or when a worker event is queued). If you wake up on a fresh `system init` and the prior turn's last observable state was a spawn-ack or a pending dispatch, treat it as idle and end the turn silently per the decision procedure above.

**Anti-patterns that indicate this bug.** If you catch yourself about to emit any of these, STOP and end the turn empty:

- `shutdown_request` with reason `"session ending"`, `"wrapping up"`, `"timeout"`, or any other self-generated reason when no completion signal has arrived. The runtime does NOT signal session-end via your context; it signals it via an actual user message.
- `shutdown_request` fired against a worker with no completion signal observed (the premature-reap bug).
- Any action whose justification is "enough time has passed" or "the ensign appears idle" — you cannot measure time from inside a turn, and ensign idleness is normal between dispatch and completion.

**DISPATCH IDLE GUARDRAIL.** After dispatching an agent, wait for an explicit completion message. Idle notifications are normal between-turn state for background workers — they are not a reason to tear a worker down, and they usually mean the agent is waiting for input. Only shut down when: (1) the agent sends a completion message, (2) the captain explicitly requests shutdown, or (3) you are transitioning the entity to a new stage (AFTER you have observed the prior stage's completion signal per the list above). Never interpret idle notifications as "stuck" or "unresponsive."

## Dispatch Failure Retry

A dispatch failure of `(entity, stage)` is exactly one of two observed things, with no error-string classification: `Agent()` returns an error, OR a dispatched worker's session ends with NO completion signal (per `## Awaiting Completion`) AND no stage report on the entity. On such a failure, read that entity's durable `### Dispatch Retries` ledger before re-attempting:

- **No retry recorded for this `(entity, stage)`** — append one ledger line, then re-attempt ONCE. If the failed worker is still addressable (a transport STALL), the re-attempt is a NUDGE: `SendMessage` it to resume from its transcript, preserving its accumulated context. If no live worker remains (an `Agent()` error, or a terminated session), the re-attempt is a FRESH `Agent()` dispatch of the same `(entity, stage)` carrying a distinct `-retry` suffix on the `{worker_key}-{slug}-{stage}` name, under dead-ensign handling — `Skill(skill="spacedock:fo-dispatch-recovery")`, its `## Context Budget Failure and Dead Ensign Handling` section: mark the prior worker dead in session memory, do NOT cooperatively shut it down.
- **A retry already recorded for this `(entity, stage)`** — the second consecutive failure. HOLD that entity un-dispatched, surface it to the captain, and stop re-attempting it. Every OTHER entity keeps running; the session dispatch mode NEVER changes.

Only ONE entity is ever affected, and only reversibly. Two adjacent captain conditions are NOT this rung and NOT a mode transition: `/spacedock bare` is a plain captain instruction to dispatch bare from that point; `Agent`/`SendMessage` themselves being unavailable is the ordinary teams-unavailable condition selecting bare dispatch where dispatch happens (`## Inter-Agent Communication`'s ABSENT branch). Bare mode itself is untouched.

**The `### Dispatch Retries` ledger** is the authority the "retry once" bound rests on: session memory is a fast-path cache, the ledger is the tiebreak, re-read before every re-attempt so a compaction cannot forget a retry and respawn a dead API unboundedly. Write it as an FO-owned `### Dispatch Retries` subsection on the entity body (the worktree copy when `worktree:` is set, else main), reusing the `### Feedback Cycles` write pattern on a DISTINCT axis so a retry never advances the feedback counter — one line per retry:

    - Retry 1: {stage} — {agent-error | no-completion-signal}; {nudged | re-dispatched -retry}

## Terminal Worker Teardown

This governs the TERMINAL phase only — AFTER the entity reached its terminal stage and the FO is dismantling its workers; do not conflate it with `## Awaiting Completion`'s ban on reaping a worker *before* a completion signal.

There is no bulk team-delete. Tear down per-roster-member:

1. Send the cooperative `SendMessage({"type":"shutdown_request"})` to every roster member in the entity's cohort. This is cooperative — the member emits `shutdown_response`/`shutdown_approved` and leaves the roster asynchronously.
2. The auto-team `members[]` prunes the terminated member asynchronously (its `inboxes/*.json` may linger); with no team-wide delete to race, nothing further is owed here. The FO tracks its own ensign roster (it already does).

## Context Budget

This is the Claude realization of `«context-budget»()` (also used by feedback rejection); Codex declares none.

**Context budget check:** Run `${SPACEDOCK_BIN:-spacedock} dispatch context-budget --name {ensign-name}`. Parse the JSON output. `reuse_ok: true` → reuse may proceed. `reuse_ok: false`, or ANY non-zero exit with no reading, is fail-safe — never silent-reuse: log to captain and fresh-dispatch. Before the replacement dispatch (budget-fail, zombie, or dead ensign), `Skill(skill="spacedock:fo-dispatch-recovery")` — its `## Context Budget Failure and Dead Ensign Handling` section carries the recovery clause for the prior worktree and the dead-ensign rules.

**Model-to-context mapping:** Resolved by `spacedock dispatch context-budget` from the member's runtime/config model; the binary owns the window mapping and follows forward family rules, so a new release in a known family stays correct without a contract edit. The same team-`config.json` member-model read supplies `«reuse.model-match»`; its canonical enum is `sonnet`, `opus`, `haiku`, `fable` (the `dispatch build` effective_model values). A captain-session fallback value (e.g. `"opus[1m]"`) is outside this enum.

## Feedback Rejection Flow (bare mode)

In bare mode, the feedback rejection flow is sequential: dispatch fix agent (wait for completion), then dispatch reviewer (wait for completion), then present at gate. With background inter-agent communication, the fix agent and reviewer can interact via messaging.

## Claude binding: «roster-reconcile»()

On Claude, invoke this binding before `«dispatch.next-action»()` and at the caller boundaries it names:

**Reconcile sweep.** Run `${SPACEDOCK_BIN:-spacedock} dispatch reconcile --workflow-dir {workflow_dir}` (a) at the first dispatch, AFTER the split-root `pull --rebase` and BEFORE the first `Agent()` dispatch; (b) inside the idle branch of `«dispatch.next-action»()`; (c) after each merge, immediately after `«worker.shutdown»()`. You have no `TeamCreate` name to pass, so omit `--team-name`: the team identity resolves by **`leadSessionId` auto-discovery** — the helper narrows to the auto-team `config.json` whose `leadSessionId` equals this session's `$CLAUDE_CODE_SESSION_ID` (set by Claude Code at launch, read by the reconcile command — no launcher plumbing needed). The roster-derived classes (lingering/superseded/un-advanced-pr) are emitted against that session-matched roster. On a headless `-p` run the auto-discovery finds no `teams/config.json` to match (the member is registered at the `subagents/agent-*.meta.json` path instead, and reading that roster is not implemented), so reconcile degrades to git-only there: an empty `drift[]` headless is not evidence of a clean roster, only that no git/filesystem drift was found. **Bare reconcile with no team identity is git-only**: it suppresses the roster-derived classes (a stale prior-session or parallel-session config must never be mistaken for the live team) and reports only the session-independent git/filesystem classes (stale-branch/local-main-drift), with a one-line stderr note. Stdout: `{"command":"reconcile","team_name":…,"drift":[{"class":"lingering|superseded|un-advanced-pr|stale-branch|local-main-drift",…}]}`. Empty `drift[]` is green. Act per drift class:
   - **lingering** / **superseded** → `SendMessage({"type":"shutdown_request"})` to `name`; drop from session memory.
   - **un-advanced-pr** → enter Merge-and-Cleanup for the named slug.
   - **stale-branch** → only when `drift.owned == true`: `git -C {worktree} pull --rebase origin {drift.trunk}`; `«halt.rebase-conflict»(paths)` on conflict. When `drift.owned` is false the item is report-only — surface it to the captain; do NOT rebase a worktree the current session does not own.
   - **local-main-drift** → behind only (`drift.behind > 0 && drift.ahead == 0`): `git -C {repo} fetch origin {drift.trunk} && git -C {repo} merge --ff-only origin/{drift.trunk} && cd {repo} && go build -o spacedock ./cmd/spacedock`. Ahead/unpushed or diverged (`drift.ahead > 0`): report-only — surface `drift.reason` to the captain and NEVER `reset --hard`; the captain decides push vs. manual reconcile.

   Non-zero helper exit (1 setup / 2 usage) surfaces to the captain; it does not block the loop. On drift, report one line: `reconcile: {N} entries: lingering={N} superseded={N} un-advanced-pr={N} stale-branch={N} local-main-drift={N} — acting`.

The idle branch of `«dispatch.next-action»()` resolves `«roster-reconcile»()` to this Claude binding. It runs once between the first and second empty `status --next`; it is not a reason to repeat the idle hook or suppress unrelated ready work.

### Backstop (Claude)

`«roster-reconcile»()` is only a backstop — the terminal teardown and the reuse module's supersede shutdown remain mandatory at their own boundaries. Cost of a miss: one extra event-loop cycle.
