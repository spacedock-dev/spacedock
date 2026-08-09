# Gate review: Restore the shared ensign contract at the Codex fresh-dispatch boundary — validation cycle 2

## Chosen direction

Keep candidate `304e09a09864889b1375fe7d41eeabd4d41e5153` unchanged. The implementation and Codex evidence are sound; the validation hold is limited to unavailable external runtime credentials.

## Evidence

- AC-1 through AC-4 have concrete citations. The candidate is clean, based on `8728da3a0`, does not include `113738b20`, and has no kd, `internal/gates`, manifest, assignment, pointer, or oracle paths.
- Focused, full, race, formatting, diff, registry, Codex common, both Codex AutoContinue, Pi fallback, durable-state, and detached adversarial checks passed. All 11 Codex process results are terminal and non-timeout.
- The exact Claude selector stopped at gate-guardrail because its OAuth session expired and could not refresh. The documented default Pi selector stopped at full-ensign-cycle with OpenRouter HTTP 402 credit exhaustion. These are infrastructure evidence failures, not candidate defects, and were not marked SKIPPED.

## Recommendation

REJECTED for an infrastructure/evidence hold only. Rerun the exact Claude lane after OAuth recovery and the default Pi lane after credit recovery; no implementation feedback cycle is requested, and no manifest or candidate change is authorized.
