# Implementation live evidence

Date: 2026-07-19
Source SHA: `2148c80c`
Codex CLI: `0.144.6`

## Fresh spawn and reuse probe

The one-off probe ran `codex exec --json --enable multi_agent_v2` with exact
parent-only canary `PARENT_ONLY_CANARY_ONEOFF_2148C80C`. It spawned one worker
with `fork_turns: "none"`, waited for completion, then invoked `followup_task`
against that same task and waited for its second response. Process exit was 0.

- Spawn call: `call_UAyMhZ7Jvqqc5kxSoODeOJm0`
- Spawn task: `fresh_isolation_oneoff`
- Recorded spawn argument: `"fork_turns":"none"`
- Child thread: `019f7a7e-d39a-7d73-9960-eade565bbb61`
- Exact parent-canary scan in child rollout: absent
- Follow-up call: `call_fUaoCElMtqiMJpI2VtLepg7S`
- Follow-up target: `/root/fresh_isolation_oneoff`
- Follow-up activity thread: `019f7a7e-d39a-7d73-9960-eade565bbb61`
- Child triggered turns: 2
- Continuity marker recalled on turn 2: `CHILD_CONTINUITY_ONEOFF_6D42`

The child first reported that no `PARENT_ONLY_CANARY_ONEOFF_` token appeared in
its inherited context. On the follow-up turn it quoted the exact continuity
marker from its prior reply and confirmed the same worker context.

## Evidence hashes

- Public `codex exec --json` stream: `8ef1cfbcc80e7c05de265e24e604b42263f66ea4005d8e4757da6e88596d42c5`
- Parent rollout: `f9481a21b79d8745f0ac1cc00c2acd05f4f0eefe975e96d50f727752d3da9194`
- Child rollout: `f4b77a92e3af10277bd25f2dc1604f338a2e62dfe8d21c127685313369d916cd`

Local raw artifacts were retained under
`/tmp/spacedock-codex-fresh-isolation-oneoff-2148c80c` and the Codex session
rollout directory for forensic inspection. This markdown record is the durable,
sanitized workflow artifact.
