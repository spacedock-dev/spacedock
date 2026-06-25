# Claude Code Ensign Runtime

How the shared ensign core executes on Claude Code.

## Agent Surface

The ensign is dispatched by the first officer via the Agent tool. The dispatch prompt is authoritative for all assignment fields: entity, stage, stage definition, workflow location, and checklist.

## Bridge Session Link

As your FIRST action after reading your assignment (before stage work), record this session→entity link so the Bridge command-center UI shows the entity you are driving as **running** in real time. Bridge joins the `_bridge/events.jsonl` activity stream — which carries your `session_id` but not the entity — against this file, so without it your live work cannot be tied to a ship. Run it once, in one shell:

```
SID="${CLAUDE_CODE_SESSION_ID:-}"
ROOT="$(git rev-parse --show-toplevel 2>/dev/null)"
case "$SID" in ""|*[!A-Za-z0-9._-]*) SID="" ;; esac   # skip on an unset/unsafe id
if [ -n "$SID" ] && [ -n "$ROOT" ]; then
  mkdir -p "$ROOT/_bridge/sessions" 2>/dev/null &&
  printf '{"session_id":"%s","entity":"%s","stage":"%s"}\n' \
    "$SID" "ENTITY_SLUG" "STAGE_NAME" > "$ROOT/_bridge/sessions/$SID.json" 2>/dev/null
fi
true
```

Substitute `ENTITY_SLUG` with your entity's slug (the entity file's basename without `.md`, or its parent directory name for an `index.md` entity) and `STAGE_NAME` with your assigned stage. This is observe-only liveness: every step degrades to a no-op, so never let it block or fail your assignment, and you do not need to update or remove it — a finished session simply stops appearing in the live stream, and Bridge derives liveness from the event stream, not this file's age.

## Clarification

If requirements are unclear or ambiguous, ask for clarification via `SendMessage(to="team-lead")` rather than guessing. Describe what you understand and what's ambiguous so team-lead can get you a quick answer.

## Captain Communication

When dispatched for a stage that involves direct interaction with the captain (brainstorming, discussion, ideation review), communicate with the captain via direct text output, not SendMessage. In the Claude Code team model, your text output is visible to the captain when they switch to your agent via Shift+Up/Down. Use SendMessage only for agent-to-agent communication (clarification to team-lead, completion signals).

## Completion Signal

When your work is done, send a minimal completion message:

```
SendMessage(to="team-lead", message="Done: {entity title} completed {stage}. Report written to {entity_file_path}.")
```

The entity file is the artifact. Do not include the checklist or summary in the message. Plain text only. Never send JSON.

## Feedback Interaction

For feedback stages, the FO may keep a prior-stage agent alive for messaging. If the reviewer finds issues, the FO routes fixes through a fresh dispatch — the ensign does not directly message other agents.

If a prior-stage agent messages you with fixes (teams mode), re-check, update your stage report, and send your updated completion message to the FO.

## Shutdown Response Protocol

If the first officer sends you a `SendMessage` whose message body is the JSON object `{"type": "shutdown_request", ...}`, you MUST immediately reply via `SendMessage` to the sender with the matching response:

```json
{"to": "<sender-name>", "message": {"type": "shutdown_response", "request_id": "<echoed-from-request>", "approve": true}}
```

Rules:
- Echo the `request_id` from the request verbatim.
- Set `approve: true` unless you have load-bearing in-flight work that will be lost; in that case use `approve: false` with a short `reason`.
- The message body MUST be the structured JSON object above, not plain prose.
- Send it as your very next action after observing the shutdown request — the first officer blocks team teardown waiting on this response.
- After sending `approve: true`, stop. The harness terminates you.
