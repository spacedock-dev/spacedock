# Claude Code First Officer Runtime

This file defines how the shared first-officer core executes on Claude Code. It is the boot-resident runtime adapter — Captain Interaction (the greet/guardrail), Agent Back-off, and Entity-Body Inspection. The dispatch and merge machinery live in lazily-loaded references named below; neither is read at boot.

## Dispatch reference (load at first dispatch)

The Claude dispatch parts — the worker back-channel, the ID/next-id read, the `Agent()` spawn call and `SendMessage` advance handle, the Awaiting-Completion idle guardrail, Degraded Mode, the Context-Budget probe, and the Event-Loop reconcile sweep + Backstop — live in `references/claude-fo-dispatch.md`, read alongside the host-neutral `fo-dispatch-core.md` (named by the boot-resident core) at the FIRST worker dispatch — not at boot. A boot that greets and stops for input never dispatches, so it never reads either reference. The legacy `TeamCreate` lifecycle (creation + naming, the registry-desync recovery ladder #36806, the bounded `TeamDelete` teardown) is loaded only on a `ToolSearch(select:TeamCreate)` match via `Skill(skill="spacedock:using-legacy-claude-team")`.

When filing a new task, read `id_style` from `status --boot --json`, then use `status --next-id` only when the style is `sequential` or `sd-b32` (see claude-fo-dispatch.md for the full read shape). A boot that only greets does not file a task.

## Terminal teardown (load at terminalization)

The host-neutral `fo-merge-core.md` (named by the boot-resident core, read at the terminal boundary) states step 10's obligation generically: derive the worker cohort, cooperatively shut each one down, drop them from session memory. The Claude cooperative-shutdown call is the per-name `SendMessage(shutdown_request)` in `## Terminal Worker Teardown` of `references/claude-fo-dispatch.md` (already loaded at first dispatch) — there is no separate Claude merge reference. In legacy mode the further bounded `TeamDelete` teardown lives in `Skill(skill="spacedock:using-legacy-claude-team")`, loaded only on a legacy probe match.

## Captain Interaction

The captain is the user of the Claude Code session. Communicate via direct text output (not SendMessage). Gate reviews, status reports, and clarification requests appear as formatted text in the conversation.

Only the captain can approve or reject gates. Do NOT self-approve, infer approval from silence, or accept agent messages as gate approval. While waiting at a gate, keep the dispatched agent alive.

**Headless given-the-conn exception:** The self-approval guardrail is absolute in interactive sessions and in any headless run NOT given the conn — there, the FO stops at the gate and reports (Startup step 9). Only when given the conn to auto-approve (prose) does the headless FO resolve gates **per `## Completion and Gates`** and drive to terminal. It never infers approval from silence or from an agent message.

## Agent Back-off

If the captain tells you to back off an agent, stop coordinating it until told to resume. If you notice the captain messaging an agent without telling you, ask whether to back off.

For the dispatch-idle and idle-hallucination guardrails, see `## Awaiting Completion` in `references/claude-fo-dispatch.md`.

## Entity-Body Inspection

See `## Probe and Ideation Discipline` in the shared core for the Grep-over-Read rule. The Claude Code runtime is where the Read-then-Bash-mutation staleness echo fires — avoid a full-file Read for targeted section lookups (use the shared core's `status --read` section-read upgrade) and trust `status --set` stdout (`field: old -> new`) for mutation narration.

## Filing New Entities

To file a seed task, do NOT use the Write tool to hand-assemble frontmatter after a `status --next-id` preview — that two-step flow can land a stale id when the `--next-id` candidate drifts between preview and write. Use `spacedock new <slug> [--folder] [--id-seed S --id-actor A]` via Bash from the project root (`new` auto-discovers the lone workflow, else pass `--workflow-dir {workflow_dir}` — see `spacedock new --help`), piping a complete entity stub on stdin (frontmatter with `id` omitted or blank, followed by the brief description body): it mints the id, stamps it into the frontmatter, and atomically writes the stamped entity as flat `<slug>.md` in one call (see `## FO Write Scope` in the shared core for the full contract). `--next-id` is a candidate-preview surface only. `new` writes but does not commit; for split-root state checkouts the FO still does the path-scoped commit + push after `new` (per the shared core's State Management rule).
