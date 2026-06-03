---
id: mtqe6vqr5hmnrjvg5ss4969m
title: dispatch build requires a hand-built JSON blob on stdin, and the contract's inline-echo example anchors the FO on shell heredocs (backticks break them) instead of the Write tool
status: ideation
source: captain (2026-06-03) — observed the FO using `python3 <<'PYEOF'` to escape backticks in checklist strings for `spacedock dispatch build`; asked why it wasn't obvious to just use the Write tool
score: "0.27"
worktree:
started: 2026-06-03T21:33:12Z
completed:
verdict:
issue:
---

`spacedock dispatch build` consumes a JSON spec on stdin for initial ensign dispatch. That keeps programmatic callers simple, but it makes the FO hand-assemble a JSON object whose `checklist`, `scope_notes`, or `feedback_context` values often contain backticks (`` `sharedRuntimeScenarios()` ``), shell variables (`$CLAUDE_CODE_SESSION_ID`), and Markdown. Inline shell JSON is fragile, and the current FO runtime docs demonstrate `echo '<json>' | spacedock dispatch build`, which anchors the FO on the fragile path.

The same input also carries the runtime host. If `host` is omitted today, `internal/dispatch/build.go` defaults to Claude and emits `Skill(skill="spacedock:ensign")`, even inside a Codex first-officer session. That mistake is silent: the command exits 0 and returns a plausible dispatch prompt for the wrong host.

## Spike result

The current behavior was reproduced against this entity on 2026-06-03:

- With `"host":"codex"`, `dispatch build` emitted `Read /tmp/spacedock-dispatch/spacedock-ensign-dispatch-build-json-ergonomics-ideation.md and treat its content as your assignment.`
- Without a host, the same request emitted `Skill(skill="spacedock:ensign"); then Read /tmp/spacedock-dispatch/spacedock-ensign-dispatch-build-json-ergonomics-ideation.md and treat its content as your assignment.`

This is the riskiest unknown for the design: Codex correctness depends on host selection, and the current helper proves the default is unsafe.

## Problem

Three coupled gaps — two binary-ergonomics, one contract-ergonomics:

1. **The binary demands a hand-escaped JSON blob for routine FO dispatch.** There is no flag/file form for per-dispatch inputs such as `stage`, `checklist`, `scope_notes`, and `feedback_context`. Any backtick or `$` in a checklist item must be escaped before the helper can read it.

2. **The runtime host is easy to omit and defaults to Claude.** In a Codex FO session, the correct request must select Codex. Without it, the helper emits a Claude-only prompt prefix (`Skill(skill="spacedock:ensign")`) and Claude completion semantics.

3. **The contract's documented example actively steers toward the footgun.** The Claude FO runtime dispatch adapter models this as `echo '<json>' | spacedock dispatch build` — inline shell piping. **The example IS the behavior the FO copies.** This session, when backticks broke the heredoc, the FO reached for a `python3 <<'PYEOF'` heredoc (a shell-native escape) rather than the obvious Write tool — *because the contract never surfaces Write as the JSON-authoring path*. The clean path (author the JSON file with the Write tool, then `cat file | dispatch build`) exists but is undocumented, so the FO followed the documented-but-fragile inline-echo pattern and patched its quoting instead of stepping back.

**Why it wasn't obvious to just use Write (the honest answer the captain asked for):** the runtime adapter's dispatch examples — both the production path and the break-glass fallback — are all inline shell piping (`echo '<json>' | ...`). That anchored the FO on a shell mental model. When the shell quoting broke, the nearest in-frame fix was another shell tool (Python heredoc), not a step out to the Write tool. A documented example that demonstrates a fragile pattern propagates that pattern; the FO optimizes against the example it was given.

## Proposed approach

1. **Add a flag/file input mode for FO dispatch.** Keep the existing stdin JSON mode for programmatic callers, but make the FO docs use this command shape:

   ```text
   spacedock dispatch build \
     --workflow-dir {workflow_dir} \
     --host {claude|codex} \
     --entity-path {entity_file_path} \
     --stage {stage} \
     --checklist-file {checklist_file} \
     [--scope-notes-file {scope_notes_file}] \
     [--feedback-context-file {feedback_context_file}] \
     [--team-name {team_name} | --bare-mode] \
     [--feedback-reflow]
   ```

   `--checklist-file` is newline-delimited text: each non-empty line becomes one checklist item and is passed through literally, with no shell interpretation. `--scope-notes-file` and `--feedback-context-file` read the whole file verbatim. This avoids shell-escaping the fields that most often contain Markdown, backticks, and `$` variables. Full JSON stdin remains available for callers that need object-level control or unusual checklist strings.

2. **Make host selection explicit and non-defaulting.** `--host claude|codex` is required in flag/file mode. In stdin JSON mode, either `--host` or a JSON `"host"` field is required; if both are present they must match. Missing host fails with a clear error such as `missing required host (pass --host claude|codex or set JSON field "host")`. Unsupported host keeps the existing loud error. There is no implicit Claude fallback.

3. **Route both input modes through the same validated request object.** The flag/file mode builds the same internal request shape as stdin JSON, then reuses the existing build validation and prompt assembly. That keeps output compatibility stable and prevents two dispatch builders from diverging.

4. **Keep `--print-schema` and `--validate-only`, scoped to JSON mode.** The rationale is concrete: stdin JSON remains a supported contract for programmatic callers, so it needs a discoverable machine schema and a no-dispatch validation path. `--print-schema` emits the canonical schema for the JSON request object, including required `host` and `schema_version: 2`. `--validate-only FILE` validates a JSON request file, exits non-zero on violations, and does not write a dispatch file. These flags are not the primary FO authoring path; they are the compatibility and tooling surface for the JSON contract.

5. **Update FO runtime docs away from inline shell JSON.** The Claude and Codex first-officer runtime docs should demonstrate the flag/file mode first. If they still mention JSON for compatibility or break-glass use, the example must be Write-tool-authored file plus input redirection, not inline `echo '<json>'`. Codex docs must say the FO passes `--host codex` and must never forward a prompt containing `Skill(skill="spacedock:ensign")`.

**Rejected alternative — only document Write-tool JSON files.** That would avoid the immediate shell quoting issue but would leave the FO hand-authoring a large JSON object for every dispatch. The binary should own routine field assembly because it already owns prompt assembly, host-specific first actions, split-root guidance, and mutation guards.

**Rejected alternative — a per-dispatch StructuredOutput tool call.** StructuredOutput is a Workflow-subagent mechanism, not a first-officer main-loop tool. It returns a validated object to context rather than to a reusable file and adds agent-spawn overhead without improving the binary's ingest validation.

## Out of scope

- Removing the JSON-stdin form.
- Redesigning `dispatch build` help routing or usage text beyond adding the new flags. That belongs to `dispatch-build-help-ergonomics`.
- The deferred-tool `SendMessage` hop (related FO-ergonomics friction, separate entity-worthy item).

## Acceptance criteria

**AC-1 — `dispatch build` accepts per-dispatch inputs without a hand-escaped JSON blob.**
Verified by: a Go command-level test driving the flag/file mode with `--checklist-file`, `--scope-notes-file`, and a checklist line containing both a backtick and `$CLAUDE_CODE_SESSION_ID`; the emitted dispatch file body contains those bytes verbatim and the stdout JSON remains valid.

**AC-2 — Stdin JSON remains supported, but missing host no longer defaults to Claude.**
Verified by: Go tests that feed valid JSON-with-host through stdin and assert the current output shape still builds; feed the same JSON without host and assert non-zero exit plus an error naming `host`; and feed conflicting `--host codex` plus JSON `"host":"claude"` and assert non-zero exit.

**AC-3 — Codex dispatch cannot silently produce a Claude-shaped prompt.**
Verified by: Go tests for `--host codex` in both flag/file and stdin JSON modes asserting the stdout prompt is the read-dispatch-file form and the dispatch file body contains neither `Skill(skill="spacedock:ensign")` nor `SendMessage(to="team-lead")`; the missing-host test in AC-2 proves the previous silent fallback is gone.

**AC-4 — JSON schema and validation are inspectable without producing a dispatch.**
Verified by: `spacedock dispatch build --print-schema` emitting valid JSON Schema for request schema version 2 with `host` required, and `spacedock dispatch build --validate-only FILE` returning 0 for a valid JSON request and non-zero for malformed JSON, missing host, unsupported host, or checklist violations. A test asserts `--validate-only` does not create or overwrite the deterministic dispatch file.

**AC-5 — FO runtime dispatch docs no longer teach inline shell JSON.**
Verified by: instruction-text invariant tests over `skills/first-officer/references/claude-first-officer-runtime.md` and `skills/first-officer/references/codex-first-officer-runtime.md` confirming primary dispatch examples use the flag/file mode, include an explicit `--host`, and do not contain `echo '<json>' | spacedock dispatch build` or heredoc JSON as the recommended path.

## Test plan

- **Command-level Go tests:** add focused tests under `internal/dispatch` for flag/file input, JSON-with-host compatibility, missing/conflicting host errors, Codex prompt shape, schema printing, and validate-only behavior. Cost: low to medium; fixtures can reuse existing build helpers.
- **Instruction text invariant:** add a small skill smoke/oracle test over both first-officer runtime adapters for the documented dispatch command shape and the banned inline-echo example. Cost: low; this is appropriate because the claim is about instruction text.
- **Regression gate:** run `go test ./...` for baseline and `go test ./... -race` before merge because dispatch prompt generation is a shared runtime surface.
- **Detached audit:** request review before merging because this affects FO-to-ensign dispatch across Claude and Codex.

## Notes

- Lived this session: the FO authored 7 ideation + 2 implementation dispatch JSONs; the first batch used a `python3` heredoc to escape backticks, the captain flagged it, and the FO switched to Write-tool-authored JSON files mid-session — confirming the clean path works and is simply undocumented.
- Lived this session: a Codex FO dispatch-helper probe reproduced the host default. With `"host":"codex"`, `dispatch build` emitted `Read /tmp/... and treat its content as your assignment.` Without host, it emitted `Skill(skill="spacedock:ensign"); then Read ...`, because `internal/dispatch/build.go` defaults missing host to Claude.
- Pairs with `launcher-binary-path-passthrough` (fc) and `extract-team-orchestration-skill` (zd) as a contract/binary ergonomics cluster — all three are "the FO/ensign hot path has a missing or mis-documented abstraction."

## Stage Report: ideation

- DONE: A concrete ergonomic input form is chosen while preserving stdin JSON compatibility.
  Chosen design is flag/file input with `--checklist-file`, `--scope-notes-file`, and `--feedback-context-file`; stdin JSON stays supported for programmatic callers.
- DONE: Runtime host behavior is explicit enough that Codex dispatch cannot silently produce a Claude-shaped prompt.
  Design requires `--host` or JSON `host`, rejects missing/conflicting host, and tests Codex output for absence of `Skill(skill="spacedock:ensign")`.
- DONE: Acceptance criteria and test plan are behavior-level and failable, including schema validation only if it remains part of the chosen design.
  AC-1 through AC-5 name command-level tests, validate-only/schema behavior, and instruction-text invariants tied to observable outputs.

### Summary

Ideation narrowed the task to a concrete flag/file input mode plus a non-defaulting host requirement, while keeping stdin JSON and schema validation as the programmatic compatibility surface. The reproduced Codex host footgun is recorded, and help-routing changes are explicitly left to `dispatch-build-help-ergonomics`.
