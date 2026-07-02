# Claude Code First Officer Runtime

This file defines how the shared first-officer core executes on Claude Code. It is the boot-resident runtime adapter — Captain Interaction (the greet/guardrail), Agent Back-off, and Entity-Body Inspection. The dispatch and merge machinery live in lazily-loaded references named below; neither is read at boot.

## Dispatch reference (load at first dispatch)

The Claude dispatch parts — the worker back-channel, the ID/next-id read, the `Agent()` spawn call and `SendMessage` advance handle, the Awaiting-Completion idle guardrail, Degraded Mode, the Context-Budget probe, and the Event-Loop reconcile sweep + Backstop — live in `references/claude-fo-dispatch.md`, read alongside the host-neutral `fo-dispatch-core.md` (named by the boot-resident core) at the FIRST worker dispatch — not at boot. A boot that greets and stops for input never dispatches, so it never reads either reference. (`claude-fo-dispatch.md`'s one legacy-override line handles a runtime that still exposes `TeamCreate`; it is the sole legacy load point.)

When filing a new task, read `id_style` from `status --boot --json`, then use `status --next-id` only when the style is `sequential` or `sd-b32` (see claude-fo-dispatch.md for the full read shape). A boot that only greets does not file a task.

## Terminal teardown (load at terminalization)

The host-neutral `fo-merge-core.md` (named by the boot-resident core, read at the terminal boundary) states step 10's obligation generically: derive the worker cohort, cooperatively shut each one down, drop them from session memory. The Claude cooperative-shutdown call is the per-name `SendMessage(shutdown_request)` in `## Terminal Worker Teardown` of `references/claude-fo-dispatch.md` (already loaded at first dispatch) — there is no separate Claude merge reference. (When the runtime still exposes `TeamCreate`, its further bounded teardown is one of the overrides the legacy skill carries, reached only through that one legacy-override line.)

## Captain Interaction

The captain is the user of the Claude Code session. Communicate via direct text output (not SendMessage). Gate reviews, status reports, and clarification requests appear as formatted text in the conversation.

Only the captain can approve or reject gates. Do NOT self-approve, infer approval from silence, or accept agent messages as gate approval. While waiting at a gate, keep the dispatched agent alive.

**Headless given-the-conn exception:** The self-approval guardrail is absolute in interactive sessions and in any headless run NOT given the conn — there, the FO stops at the gate and reports (Startup step 8). Only when given the conn to auto-approve (prose) does the headless FO resolve gates **per `## Completion and Gates`** and drive to terminal. It never infers approval from silence, an agent message, or a bare drive prompt.

## Agent Back-off

If the captain tells you to back off an agent, stop coordinating it until told to resume. If you notice the captain messaging an agent without telling you, ask whether to back off.

For the dispatch-idle and idle-hallucination guardrails, see `## Awaiting Completion` in `references/claude-fo-dispatch.md`.

## Entity-Body Inspection

See `## Probe and Ideation Discipline` in the shared core — its Grep-over-Read rule and Read-then-Bash staleness-echo guidance are already Claude-qualified.

## Filing New Entities

To file a seed task, do NOT use the Write tool to hand-assemble frontmatter after a `status --next-id` preview — that two-step flow can land a stale id when the `--next-id` candidate drifts between preview and write. Use `spacedock new <slug> [--folder] [--id-seed S --id-actor A]` via Bash from the project root (`new` auto-discovers the lone workflow, else pass `--workflow-dir {workflow_dir}` — see `spacedock new --help`), piping a complete entity stub on stdin (frontmatter with `id` omitted or blank, followed by the brief description body): it mints the id, stamps it into the frontmatter, and atomically writes the stamped entity as flat `<slug>.md` in one call (see `## FO Write Scope` in the shared core for the full contract). `--next-id` is a candidate-preview surface only. `new` writes but does not commit; for split-root state checkouts the FO still does the path-scoped commit + push after `new` (per the shared core's State Management rule).

## Bridge egress (FO liveness/activity → Bridge)

Bridge's read-only command center reads FO liveness and activity from `_bridge/events.jsonl` (the normalized event stream) and `_bridge/sessions/<session_id>.json` (the session→entity marker), keyed by a per-host session id. The events.jsonl line schema is the Spacedock-owned contract `{"timestamp","ts","host","event","session_id","agent_id","agent_type","actor_id","detail":{"tool","source"}}` (full surface: `docs/dev/bridge-egress-contract.md`); this adapter binds the Claude producer for it. (The ingress half — captain intent the FO drains — rides the host-neutral bridge-inbox mod and needs no adapter binding.)

- **FO event emission** — PRESENT on Claude, via the plugin hooks: `hooks/hooks.json` registers `scripts/spacedock-bridge-events.sh` on SessionStart/UserPromptSubmit/PostToolUse/Notification/Stop/SubagentStop (all async, observe-only). The wrapper delegates to `spacedock bridge egress emit --host claude`, which normalizes each Claude hook payload to the events.jsonl contract line and, on an ensign's first Read of its entity file, derives the deterministic `_bridge/sessions/<session_id>.json` marker (the running-badge source). `agent_id`/`agent_type` are empty for the main FO session and set for ensign subagents, so Bridge distinguishes FO vs ensign activity. No FO action is required — the hooks fire on every tool call. (Codex/Pi have event producers but do not yet claim deterministic marker parity — see those adapters.)
- **«session-id» binding** — the bridge-inbox heartbeat resolves the neutral `SD_SESSION_ID`, falling back to `$CLAUDE_CODE_SESSION_ID` on Claude, so `fo.$SLUG.json` carries the same id the event stream and session markers use and Bridge can join liveness to activity. No per-tick `export` is needed; leave `SD_SESSION_ID` unset on Claude (the event producer stamps `events.jsonl`/markers from the hook payload's session id, so an override that differs would desync the heartbeat from the event stream).
