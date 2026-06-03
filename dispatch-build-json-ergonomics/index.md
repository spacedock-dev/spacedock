---
id: mtqe6vqr5hmnrjvg5ss4969m
title: dispatch build requires a hand-built JSON blob on stdin, and the contract's inline-echo example anchors the FO on shell heredocs (backticks break them) instead of the Write tool
status: backlog
source: captain (2026-06-03) — observed the FO using `python3 <<'PYEOF'` to escape backticks in checklist strings for `spacedock dispatch build`; asked why it wasn't obvious to just use the Write tool
score: "0.27"
worktree:
started:
completed:
verdict:
issue:
---

`spacedock dispatch build` consumes a JSON spec on stdin (the mandated initial-dispatch path). The FO must hand-assemble that JSON, including a `checklist` array whose strings frequently contain backticks (`` `sharedRuntimeScenarios()` ``) and `$`-vars (`$CLAUDE_CODE_SESSION_ID`). Those characters break bash heredoc/quote escaping, so building the JSON inline in the shell is fragile. The same JSON blob also carries the runtime host; if omitted, the helper silently defaults to Claude and emits `Skill(skill="spacedock:ensign")`, which is wrong in a Codex first-officer session.

## Problem

Three coupled gaps — two binary-ergonomics, one contract-ergonomics:

1. **The binary demands a hand-escaped JSON blob.** There is no flag/file form for the per-dispatch inputs (stage, checklist, scope_notes); the only path is a full JSON document on stdin. Any backtick or `$` in a checklist item must be shell-escaped by the author.

2. **The runtime host is easy to omit and defaults to Claude.** In a Codex FO session, the correct JSON must include `"host":"codex"`. Without it, the current helper defaults to Claude and emits a Claude-only prompt prefix (`Skill(skill="spacedock:ensign")`) plus team-signaling assumptions. The mistake is silent: the command succeeds and produces a prompt that looks plausible until a Codex FO tries to dispatch it.

3. **The contract's documented example actively steers toward the footgun.** The Claude FO runtime dispatch adapter models this as `echo '<json>' | spacedock dispatch build` — inline shell piping. **The example IS the behavior the FO copies.** This session, when backticks broke the heredoc, the FO reached for a `python3 <<'PYEOF'` heredoc (a shell-native escape) rather than the obvious Write tool — *because the contract never surfaces Write as the JSON-authoring path*. The clean path (author the JSON file with the Write tool, then `cat file | dispatch build`) exists but is undocumented, so the FO followed the documented-but-fragile inline-echo pattern and patched its quoting instead of stepping back.

**Why it wasn't obvious to just use Write (the honest answer the captain asked for):** the runtime adapter's dispatch examples — both the production path and the break-glass fallback — are all inline shell piping (`echo '<json>' | ...`). That anchored the FO on a shell mental model. When the shell quoting broke, the nearest in-frame fix was another shell tool (Python heredoc), not a step out to the Write tool. A documented example that demonstrates a fragile pattern propagates that pattern; the FO optimizes against the example it was given.

## Proposed approach

1. **Binary: an ergonomic input form** so no hand-escaped JSON blob is required — e.g. flags (`--stage`, `--entity`, `--checklist-file FILE`, `--scope-notes-file FILE`) or reading the checklist as a newline-delimited file, so a backtick is never shell-escaped. (This is the "dispatch-build JSON robustness" candidate the roadmap-readiness audit surfaced.) Keep stdin-JSON as a supported form; add the ergonomic one.
2. **Binary: make host explicit.** Add a `--host claude|codex` command flag, or require `host` in the JSON with no silent default. Prefer the flag form when the ergonomic input path lands, so the runtime adapter can put host selection outside the fragile JSON payload. If compatibility keeps a Claude default temporarily, Codex runtime docs and tests must still prove the Codex path cannot omit host silently.
3. **Binary: fail loud on malformed JSON and unsupported/missing host** (clear error naming the offending field), not a partial/silent build or a Claude-shaped dispatch in a Codex session.
4. **Contract: fix the documented pattern.** The FO runtime dispatch adapter should author the JSON via the Write tool to a file, then `cat file | dispatch build --host codex|claude` (or use the new flag form) — NOT inline `echo '<json>'`. The break-glass fallback example should follow suit. A documented clean path beats a documented footgun.
5. **Binary: expose the schema.** `spacedock dispatch build --print-schema` emits the canonical dispatch-spec JSON Schema, and `--validate-only FILE` checks a spec without building — so the schema lives where the contract lives, and any authoring path (Write tool, editor, a future structured-output flow) validates against it. **Rejected alternative — a per-dispatch StructuredOutput tool call:** StructuredOutput is a Workflow-subagent-only mechanism (the `schema:` option on `agent()`), not callable from the FO main loop, and it returns the validated object to the caller's context rather than to a file; it would add agent-spawn overhead over the Write tool with no safety the binary's own ingest-validation doesn't already provide.

## Out of scope

- Removing the JSON-stdin form (keep it for programmatic callers).
- The deferred-tool `SendMessage` hop (related FO-ergonomics friction, separate entity-worthy item).

## Acceptance criteria

**AC-1 — `dispatch build` accepts per-dispatch inputs without a hand-escaped JSON blob.**
Verified by: a Go test driving the new input form (flags or `--checklist-file`) with a checklist item containing a backtick and a `$`-var, asserting the emitted dispatch spec carries the item verbatim.

**AC-2 — Malformed JSON input fails loud, and the dispatch-spec schema is inspectable.**
Verified by: a test feeding truncated/invalid JSON and asserting a non-zero exit with an error naming the problem (not a partial build); plus `spacedock dispatch build --validate-only FILE` returning 0 on a valid spec and non-zero (naming the violation) on an invalid one, and `--print-schema` emitting the canonical `schema_version: 2` JSON Schema so any authoring tool validates against the same contract the build path enforces.

**AC-3 — The FO runtime dispatch-adapter prose documents the Write-to-file (or flag) path, not inline `echo '<json>'`.**
Verified by: a presence/absence oracle over `claude-first-officer-runtime.md` and `codex-first-officer-runtime.md` confirming dispatch examples use file/flag authoring, pass an explicit host, and no longer demonstrate inline-echo JSON as the primary path.

**AC-4 — Codex dispatch cannot silently fall back to a Claude-shaped prompt.**
Verified by: Go tests assert the Codex dispatch path emits the read-dispatch-file prompt and a body without `Skill(skill="spacedock:ensign")` or `SendMessage(to="team-lead")`; a missing/unsupported host in the ergonomic input path either fails loud or is impossible because `--host` is required. A compatibility test may keep legacy JSON-without-host as Claude only if the runtime adapter no longer uses that form.

## Test plan

- Go unit tests for the new input form, explicit host selection, and malformed-input guard (AC-1, AC-2, AC-4). Cost: low.
- Instruction-text invariant for AC-3 (the proven contract-oracle pattern). Cost: trivial.
- High-stakes surface (dispatch machinery the whole workflow rides on) → detached audit before merge.

## Notes

- Lived this session: the FO authored 7 ideation + 2 implementation dispatch JSONs; the first batch used a `python3` heredoc to escape backticks, the captain flagged it, and the FO switched to Write-tool-authored JSON files mid-session — confirming the clean path works and is simply undocumented.
- Lived this session: a Codex FO dispatch-helper probe reproduced the host default. With `"host":"codex"`, `dispatch build` emitted `Read /tmp/... and treat its content as your assignment.` Without host, it emitted `Skill(skill="spacedock:ensign"); then Read ...`, because `internal/dispatch/build.go` defaults missing host to Claude.
- Pairs with `launcher-binary-path-passthrough` (fc) and `extract-team-orchestration-skill` (zd) as a contract/binary ergonomics cluster — all three are "the FO/ensign hot path has a missing or mis-documented abstraction."
