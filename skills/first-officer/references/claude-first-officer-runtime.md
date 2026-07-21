# Claude Code First Officer Runtime

This file defines how the shared first-officer core executes on Claude Code: Captain Interaction (the greet/guardrail), Agent Back-off, and Entity-Body Inspection. The dispatch, write, and merge cores load only at the triggers named by the shared core.

## Dispatch reference (load at first dispatch)

The Claude dispatch parts — the worker back-channel, the ID/next-id read, the `Agent()` spawn call and `SendMessage` advance handle, the Awaiting-Completion idle guardrail, the Degraded-Mode/break-glass/budget-failure trigger lines, the Context-Budget probe, and the Event-Loop reconcile sweep + Backstop — live in `references/claude-fo-dispatch.md`, read alongside `fo-dispatch-core.md` at the first worker dispatch (the exception bodies behind those triggers load at failure time via `Skill(skill="spacedock:fo-dispatch-recovery")`). (`claude-fo-dispatch.md`'s one legacy-override line handles a runtime that still exposes `TeamCreate`; it is the sole legacy load point.)

When filing a new task, read `id_style` from `status --boot --json`, then use `status --next-id` only when the style is `sequential` or `sd-b32` (see claude-fo-dispatch.md for the full read shape). A boot that only greets does not file a task.

## Terminal teardown (load at terminalization)

The shared core loads `fo-merge-core.md` at the terminal/recovery boundary; that core invokes `«worker.shutdown»()`. On Claude, realize a deferred-core load as its own successful `Read` call against the exact path formed from retained `{first_officer_base}`. A terminal status mutation must complete the write-core `Read`, then the merge-core `Read`, before its Bash call; do not infer either core from this adapter or skip its read. The Claude shutdown binding is the per-name `SendMessage(shutdown_request)` in `## Terminal Worker Teardown` of `references/claude-fo-dispatch.md` (already loaded at first dispatch). When the runtime still exposes `TeamCreate`, its further bounded teardown is the legacy override.

## Captain Interaction

The captain is the user of the Claude Code session. Communicate via direct text output (not SendMessage). Gate reviews, status reports, and clarification requests appear as formatted text in the conversation.

Only the captain can approve or reject gates. Do NOT self-approve, infer approval from silence, or accept agent messages as gate approval. While waiting at a gate, keep the dispatched agent alive.

`«interaction.boundary»()` owns the headless given-the-conn exception. The self-approval guardrail is absolute in interactive sessions and in any headless run not given the conn; only an explicit grant lets headless resolve gates per `## Completion and Gates` and drive to terminal.

## Agent Back-off

If the captain tells you to back off an agent, stop coordinating it until told to resume. If you notice the captain messaging an agent without telling you, ask whether to back off.

For the dispatch-idle and idle-hallucination guardrails, see `## Awaiting Completion` in `references/claude-fo-dispatch.md`.

## Entity-Body Inspection

See `## Probe and Ideation Discipline` in the shared core for the Grep-over-Read rule; ToolSearch is this host's tool-discovery surface. The Claude read-then-mutate caveat: a `Read` followed by a Bash mutation of the same file (including `status --set`) triggers the file-staleness safety net, echoing the file back as cache-write tokens. Grep does not participate. Trust `status --set` stdout (`field: old -> new`, `field: old -> ` for clear-to-empty, `field:  -> {timestamp}` for bare-timestamp auto-fill) to narrate mutations without re-reading.

## Filing New Entities

Before filing, read the workflow README's `## Task Template`; use its frontmatter and section scaffolding as the starting shape for the entity body you pipe to `spacedock new`.

To file a seed task, do NOT use the Write tool to hand-assemble frontmatter after a `status --next-id` preview — that two-step flow can land a stale id when the `--next-id` candidate drifts between preview and write. After the shared core's write load point has completed, use `${SPACEDOCK_BIN:-spacedock} new <slug> [--folder] [--id-seed S --id-actor A]` via Bash from the project root (`new` auto-discovers the lone workflow, else pass `--workflow-dir {workflow_dir}` — see `spacedock new --help`), piping a complete entity stub on stdin (frontmatter with `id` omitted or blank, followed by the brief description body): it mints the id, stamps it into the frontmatter, and atomically writes the stamped entity as flat `<slug>.md` in one call. `--next-id` is a candidate-preview surface only. `new` writes but does not commit; for split-root state checkouts the FO still does the path-scoped commit + push after `new` (per the shared core's State Management rule).
