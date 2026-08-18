# Ideation spike evidence: the Claude compaction boundary, captured

Provenance: captured 2026-08-18 (UTC 2026-08-18T23:2x) during ideation for
`force-boot-at-compaction-boundary`. Host: Claude Code v2.1.226, macOS.
Method: a live scratch session driven via tmux, launched with
`claude --settings hook-settings.json --model haiku` where the settings file
registered PreCompact and SessionStart capture hooks (`cat >> events.jsonl`,
no matcher, capturing every source as a control). Conversation was grown to
~28k tokens over five turns, then `/compact` was issued twice (once refused,
once completed). Records below are verbatim; nothing is paraphrased.

## 1. Incident boundary record (durable transcript state)

From the incident session's own transcript,
`~/.claude/projects/-Users-clkao-git-spacedock-research-spacedock-v1/fedfe9c3-f6c8-4ac8-ba26-73620a69d3a1.jsonl`
(line 3130 of 4649) — the compaction that produced the six stale-binding failures:

```json
{"parentUuid":null,"logicalParentUuid":"6b99d305-cee5-406c-ae10-3bd1ec8d4eaf","isSidechain":false,"type":"system","subtype":"compact_boundary","content":"Conversation compacted","level":"info","compactMetadata":{"trigger":"manual","preTokens":801920,"postTokens":19650,"cumulativeDroppedTokens":782270,"durationMs":121692,"preCompactDiscoveredTools":["SendMessage","TaskList","TaskOutput","TaskStop"],"preservedSegment":{"headUuid":"4645e36d-8657-4b0e-9d4d-82667f3e4409","anchorUuid":"5c42641c-3b0b-41b3-97b7-822914a30b8e","tailUuid":"6b99d305-cee5-406c-ae10-3bd1ec8d4eaf"},"preservedMessages":{"anchorUuid":"5c42641c-3b0b-41b3-97b7-822914a30b8e","uuids":["4645e36d-8657-4b0e-9d4d-82667f3e4409","6b99d305-cee5-406c-ae10-3bd1ec8d4eaf"],"allUuids":["4645e36d-8657-4b0e-9d4d-82667f3e4409","8db20dca-608b-4f3e-860a-e292e4d1bf9b","6b99d305-cee5-406c-ae10-3bd1ec8d4eaf"]}},"uuid":"7b52ac75-04d8-4608-b7e1-07da94f5146c","timestamp":"2026-08-18T18:59:30.622Z","userType":"external","entrypoint":"cli","cwd":"/Users/clkao/git/spacedock-research/spacedock-v1","sessionId":"fedfe9c3-f6c8-4ac8-ba26-73620a69d3a1","version":"2.1.226","gitBranch":"main","slug":"curried-brewing-kahn"}
```

Same `sessionId` before and after the boundary; the transcript file does not
change. `trigger:"manual"`, 801,920 tokens dropped to 19,650.

## 2. Spike hook events (all four, in firing order)

```json
{"session_id":"8306a3cd-9558-4be7-ac79-309f96b696d0","transcript_path":"/Users/clkao/.claude/projects/-private-tmp-claude-501--Users-clkao-git-spacedock-research-spacedock-v1-fedfe9c3-f6c8-4ac8-ba26-73620a69d3a1-scratchpad-compaction-spike/8306a3cd-9558-4be7-ac79-309f96b696d0.jsonl","cwd":"/private/tmp/claude-501/-Users-clkao-git-spacedock-research-spacedock-v1/fedfe9c3-f6c8-4ac8-ba26-73620a69d3a1/scratchpad/compaction-spike","hook_event_name":"SessionStart","source":"startup","model":"claude-haiku-4-5-20251001"}
{"session_id":"8306a3cd-9558-4be7-ac79-309f96b696d0","transcript_path":"/Users/clkao/.claude/projects/-private-tmp-claude-501--Users-clkao-git-spacedock-research-spacedock-v1-fedfe9c3-f6c8-4ac8-ba26-73620a69d3a1-scratchpad-compaction-spike/8306a3cd-9558-4be7-ac79-309f96b696d0.jsonl","cwd":"/private/tmp/claude-501/-Users-clkao-git-spacedock-research-spacedock-v1/fedfe9c3-f6c8-4ac8-ba26-73620a69d3a1/scratchpad/compaction-spike","prompt_id":"05f01048-70ee-47f4-88fe-ae289247cca5","hook_event_name":"PreCompact","trigger":"manual","custom_instructions":null}
{"session_id":"8306a3cd-9558-4be7-ac79-309f96b696d0","transcript_path":"/Users/clkao/.claude/projects/-private-tmp-claude-501--Users-clkao-git-spacedock-research-spacedock-v1-fedfe9c3-f6c8-4ac8-ba26-73620a69d3a1-scratchpad-compaction-spike/8306a3cd-9558-4be7-ac79-309f96b696d0.jsonl","cwd":"/private/tmp/claude-501/-Users-clkao-git-spacedock-research-spacedock-v1/fedfe9c3-f6c8-4ac8-ba26-73620a69d3a1/scratchpad/compaction-spike","prompt_id":"1d64fb8f-f63c-45fc-804e-d4781680603a","hook_event_name":"PreCompact","trigger":"manual","custom_instructions":null}
{"session_id":"8306a3cd-9558-4be7-ac79-309f96b696d0","transcript_path":"/Users/clkao/.claude/projects/-private-tmp-claude-501--Users-clkao-git-spacedock-research-spacedock-v1-fedfe9c3-f6c8-4ac8-ba26-73620a69d3a1-scratchpad-compaction-spike/8306a3cd-9558-4be7-ac79-309f96b696d0.jsonl","cwd":"/private/tmp/claude-501/-Users-clkao-git-spacedock-research-spacedock-v1/fedfe9c3-f6c8-4ac8-ba26-73620a69d3a1/scratchpad/compaction-spike","prompt_id":"1d64fb8f-f63c-45fc-804e-d4781680603a","hook_event_name":"SessionStart","source":"compact","model":"claude-haiku-4-5-20251001"}
```

Reading order: SessionStart(source:"startup") proves the capture harness works.
The first PreCompact fired at a `/compact` that the host then REFUSED
("Not enough messages to compact") — PreCompact fires before the
enough-to-compact check, so it can false-positive. The second PreCompact and
the SessionStart(source:"compact") bracket the completed compaction.
`session_id` is identical in all four events: session identity does not change
at the compaction boundary.

## 3. Spike transcript boundary record

From the spike session's transcript (same path the hook events name):

```json
{"parentUuid":null,"logicalParentUuid":"5f43d1b8-dfb5-47f3-bffc-d0e357f9ff58","isSidechain":false,"type":"system","subtype":"compact_boundary","content":"Conversation compacted","level":"info","compactMetadata":{"trigger":"manual","preTokens":28139,"postTokens":3673,"cumulativeDroppedTokens":24466,"durationMs":18220,"preservedSegment":{"headUuid":"09cc7c65-d0f3-4b78-b8d6-d873e9de20d7","anchorUuid":"e197a9e2-f29d-47af-95ee-e3e2181b4a9e","tailUuid":"5f43d1b8-dfb5-47f3-bffc-d0e357f9ff58"},"preservedMessages":{"anchorUuid":"e197a9e2-f29d-47af-95ee-e3e2181b4a9e","uuids":["09cc7c65-d0f3-4b78-b8d6-d873e9de20d7","82ec7e97-9ed9-436d-80e7-dea3de608c45","e1599dcf-2ef7-49f3-bcdd-f24df4608194","5f43d1b8-dfb5-47f3-bffc-d0e357f9ff58"],"allUuids":["09cc7c65-d0f3-4b78-b8d6-d873e9de20d7","82ec7e97-9ed9-436d-80e7-dea3de608c45","ac830bc2-06e7-47e5-90c8-7211861ff4a4","e1599dcf-2ef7-49f3-bcdd-f24df4608194","5f43d1b8-dfb5-47f3-bffc-d0e357f9ff58"]}},"uuid":"0cb199a5-6784-45b4-a4fc-1302a4054d6f","timestamp":"2026-08-18T23:26:56.079Z","userType":"external","entrypoint":"cli","cwd":"/private/tmp/claude-501/-Users-clkao-git-spacedock-research-spacedock-v1/fedfe9c3-f6c8-4ac8-ba26-73620a69d3a1/scratchpad/compaction-spike","sessionId":"8306a3cd-9558-4be7-ac79-309f96b696d0","version":"2.1.226","gitBranch":"HEAD","slug":"groovy-watching-dove"}
```

## 4. Session identity is visible to the model-side shell

Inside the spike session, after the completed compaction, the session ran:

    $ echo SID=$CLAUDE_CODE_SESSION_ID CHILD=$CLAUDE_CODE_CHILD_SESSION
    SID=8306a3cd-9558-4be7-ac79-309f96b696d0 CHILD=1

`CLAUDE_CODE_SESSION_ID` matches the transcript filename and survives the
boundary unchanged. (CHILD=1 is inherited noise from the launching
environment; Claude Code overrides SESSION_ID per session but does not clear
CHILD.) The FO's own environment carries the same variable
(`CLAUDE_CODE_SESSION_ID=fedfe9c3-…`), so a binary invoked by the FO can
resolve the running session's transcript with no hook callback.

## 5. Negative findings

- The incident session's transcript contains ZERO SessionStart hook activity
  and none of the registered prose-injection hook's text ("Spacedock: reread
  the authoritative…"), while 216 `stop_hook_summary` records prove hooks in
  general were active in that session. The plugin (0.27.0-pre1 hooks.json,
  matcher `^compact$`, `codex_session_start_compact.sh`) was registered; its
  injection produced no observable effect at the incident boundary, and the
  six failures followed.
- The two July transcripts that contain the hook's text (bef9653f, deed4af4)
  are development chatter — tool_results reading the entity/hook files while
  the hook was being BUILT — not live firings. No local evidence exists that
  the prose-injection hook has ever fired live.
- This repo's `hook-events.jsonl` holds exactly one PreToolUse record and no
  compaction evidence, as the dispatch stated.
