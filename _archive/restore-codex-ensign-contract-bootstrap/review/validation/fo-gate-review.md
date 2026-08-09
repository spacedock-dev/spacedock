# Gate review: Restore the shared ensign contract at the Codex fresh-dispatch boundary — validation

## Chosen direction

Keep the exact candidate unchanged. The implementation and offline validation are sound; the remaining decision is only whether to accept the live-runtime evidence hold.

## Evidence

- Cycle 4 implementation: 3 DONE, 0 SKIPPED, 0 FAILED; candidate `7d50c7d88` is clean.
- Validation AC-1 through AC-4 each have evidence; focused, full, race, registry, formatting, transport, and detached adversarial checks pass.
- The exact-head live run emits the required `$spacedock:ensign` prompt and final message, but the Codex process hangs after `turn.completed`; the watchdog kills it before durable terminal evidence.
- No `kd`, oracle, assignment-transport, or candidate defect was found.

## Recommendation

Hold the gate for a Captain runtime-lane decision. Do not send the candidate back to implementation for this host-lifecycle failure.
