---
name: fo-dispatch-recovery
description: "Claude dispatch failure recovery — Degraded Mode (triggers, effects, the verbatim captain report, the cooperative shutdown sweep), Break-Glass Manual Dispatch (the manual `Agent()` template), and Context Budget Failure/Dead Ensign Handling (the budget-unavailable stderr conditions, the recovery clause, dead-ensign bookkeeping). Read ONLY at its three resident triggers inside `claude-fo-dispatch.md` — a second dispatch failure, `/spacedock bare`, or `Agent`/`SendMessage` unavailable (Degraded Mode); a non-zero or unavailable `spacedock dispatch build` (Break-Glass); or a budget-fail/zombie/dead-ensign replacement dispatch (Context Budget) — never at boot, never on the happy path."
user-invocable: false
---

# First Officer Dispatch Recovery (Claude)

The three Claude dispatch exception bodies, each read only at its own failure trigger in `claude-fo-dispatch.md` — never at boot, never on a session where dispatch never fails.

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

## Break-Glass Manual Dispatch

The resident trigger line already covers the first action (report the helper failure — command, exit code, stderr — to the captain before proceeding). Fill this template directly as the degraded fallback:
Populate `{numbered checklist}` with the output of `«dispatch.checklist»(entity, stage)`; do not rebuild its rules here.
```
Agent(
    subagent_type="{dispatch_agent_id}",
    name="{worker_key}-{slug}-{stage}",  // if this exceeds 64 chars, cap it the way `spacedock dispatch build` does: keep the {worker_key} prefix and -{stage} suffix and, on id-style: sd-b32, replace the slug with a fixed-length prefix of the entity id (id-less slug workflows truncate the slug head instead)
    run_in_background=true,              // named background worker, no team_name
    model="{effective_model}",
    prompt="## First action\n\nBefore anything else, invoke your operating contract:\n\n    Skill(skill=\"spacedock:ensign\")\n\nThis loads the shared ensign discipline (stage-report format, background-task polling, worktree ownership, completion signal protocol). Do not paraphrase; call the tool.\n\nYou are working on: {entity title}\n\nStage: {stage}\n\n### Stage definition:\n\n{copy stage subsection from README verbatim}\n\nRead the entity file at {entity_file_path}.\n\n### Completion checklist\n\n{numbered checklist}\n\n### Summary\n{brief description of what was accomplished}\n\n### Stage report\n\nAppend a Stage Report section at the end of the entity file (per the shared-core Stage Report Protocol). Use the title `Stage Report: {stage}`. Account for every checklist item above with a `- DONE:` / `- SKIPPED:` / `- FAILED:` entry. Use the checklist item text verbatim when possible.\n\n### Completion Signal\n\nSendMessage(to=\"team-lead\", message=\"Done: {entity title} completed {stage}. Report written to {entity_file_path}.\")"
)
```
This is the concrete Claude form of fo-dispatch-core.md's Break-Glass template; the contract (what it omits, the conditional `model=` slot, "use only when the helper is unavailable") is stated there. The canonical enum the conditional slot draws from is the resident `## Context Budget` section of `claude-fo-dispatch.md` (already loaded alongside this skill at the first dispatch).

## Context Budget Failure and Dead Ensign Handling

**Budget-unavailable is fail-safe (never silent-reuse).** The probe exits non-zero with no `reuse_ok` field in three conditions; the FO treats every one identically — fresh-dispatch:
- **missing jsonl** — no `agent-*.jsonl` exists for the named member (stderr: `no subagent jsonl found for '{name}'`).
- **unreadable/empty jsonl** — the jsonl exists but carries no assistant entry with non-zero `usage` (stderr: `no assistant entries with usage in {path}`).
- **agent-not-in-team-config** — no team `config.json` lists a member with that name (stderr: `no team config found for member '{name}'`).
A non-zero exit with no `reuse_ok: true` means the FO never silent-reuses on an absent reading.

**Recovery clause** (only when replacing a prior ensign): The prior ensign was shut down due to context budget limits. Its worktree may hold uncommitted changes. Run `git status` and `git diff` first; commit legitimate WIP or reset broken changes.

**Dead ensign handling:**

- `SendMessage(shutdown_request)` is cooperative — do NOT send to dead or unresponsive ensigns.
- Track dead ensigns in session memory; do not route work to dead names.
- Fresh-dispatch under a `-cycleN` suffix when replacing a zombie ensign.
- The post-dispatch config check does NOT detect zombies — zombies pass it. Session memory is the authoritative dead-vs-alive tracker.
