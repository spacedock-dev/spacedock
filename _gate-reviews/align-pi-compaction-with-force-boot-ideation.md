# FO Gate Review — `align-pi-compaction-with-force-boot` ideation

## Question
Is the ideation approach sound for aligning the Pi extension's compaction hook with PR #738 (re-read state, not re-inject contract), with well-formed ACs and the spike resolved?

## Verdict: APPROVE → implementation

## Factual claims verified against the code
- `session_compact` sets `injectBootstrap = true` (spacedock.ts:67-68) — confirmed.
- `context` handler injects `FO_BOOTSTRAP_TEXT` (spacedock.ts:81) — confirmed.
- `FO_BOOTSTRAP_MARKER = "SPACEDOCK-FO-BOOTSTRAP-v1"` (spacedock.ts:22) — confirmed.
- `pi.exec(command, args, options): Promise<ExecResult>` IS a method on the ExtensionAPI type (types.d.ts:947) — confirmed. The ensign read the correct declaration.
- The test harness builds `pi = { on(...) }` (spacedock_extension_test.go:36); adding an `exec` mock is feasible.
- `TestSpacedockPiExtensionChildExemption` covers `session_compact` exemption — confirmed.

## Assessment
- **Approach sound**: separate `injectBootRecord` flag for compaction; `context` runs `spacedock status --boot --identify --json` via `pi.exec()` and injects the boot record as a user message; remove `FO_BOOTSTRAP_TEXT` from compaction path (retain for `session_start`). Goes one step beyond Claude/Codex (which only remind) by actually running the read — defensible given `pi.exec()` is available.
- **ACs well-formed**: 5 ACs — value-measuring (AC-1: resumes on boot record, not contract re-injection, vs. baseline that re-injects without reading), mechanism (AC-2: `pi.exec` fired + boot record injected), scope boundary (AC-3: `session_start` unchanged), child exemption (AC-4), dedup (AC-5). Each is testable.
- **Spike resolved**: "no spike needed" grounded in the real ExtensionAPI type declaration (`exec` at line 947) + the `context` handler being async-capable. The actual first-use proof is in implementation (the test harness mock).
- **Scope clean**: observable-semantics declaration correct — only the compaction-boundary context content changes; no FO bootstrap content, gate/dispatch/state mechanics, or command surface touched.

## Implementation notes for the ensign
- The `context` handler becomes async on the boot-record path; the test harness must `await` it. Decide whether `session_start` stays sync or both go async (simpler: both async).
- Dedup check needs a boot-record marker distinct from `FO_BOOTSTRAP_MARKER` (the report uses `[SPACEDOCK-FO-BOOT-v2]`).
- `pi.exec` returns `{stdout, stderr, code, killed}` — parse `stdout` JSON; handle nonzero code (fall back to a minimal directive? or fail loud? — the report should specify).
