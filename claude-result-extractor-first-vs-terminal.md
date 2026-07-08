---
id: 6h08n9jrwa9g5kgm3b3fy8vr
title: "Claude final-message extractor returns the first result event, not the terminal one"
status: ideation
source: "Codex-session CI investigation of PR #483's claude-live (opus) red (2026-07-08), independently re-verified file:line by the FO. Confirmed unrelated to PR #483's own commit (43396704, scoped only to codex_liveenv.go/codex_liveenv_test.go) — this is pre-existing shared test-harness infrastructure, not a regression from that change."
started: 2026-07-08T04:10:04Z
completed:
verdict:
score:
worktree:
issue:
---

`extractClaudeFinalMessage` (`internal/ensigncycle/claude_final_message_impl_test.go:55-89`) loops a Claude stream-json transcript and returns on the FIRST `result`-type event with `IsError == false` (line 72's early `return result.Result, nil`) — but the function's own doc comment (lines 44-51) says "a **terminal** result/success event's `result` field... is preferred," implying the LAST one. When a transcript contains more than one success `result` event (observed live: PR #483's failing opus run had two — an early one saying only "All four ensigns are dispatched..." and a later one that actually names the corrected entity and presents the gate), the extractor silently returns the wrong one. This makes any shared-scenario assertion built on the extracted text flake-prone: pass/fail depends on the incidental wording of an intermediate result, not the run's actual final behavior — confirmed by the prior green opus run also having two results, passing only because its first one happened to mention the asserted phrase. Not caused by, and not specific to, PR #483's own change.

Recommended fix (already spiked/diagnosed, not an open design question): change the extractor to retain the LAST non-error success result rather than returning on the first match, add an offline regression test with a fixture stream containing multiple success `result` events asserting the last one wins, and keep the existing loud-failure behavior unchanged for `is_error` / 401 result events (never fall back to a later event when an earlier one is a launch failure).
