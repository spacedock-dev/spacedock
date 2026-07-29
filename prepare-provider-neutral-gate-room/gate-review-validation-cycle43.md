# Validation gate: provider-neutral gate preparation

## Recommendation

Approve s4 at `acae980fc145e624d9e04e7ec9f7fdb599585f6e`.

The candidate gives a First Officer one mechanical command to prepare a review
room. Spacedock chooses the room, writes its frozen request and canonical
Briefing, records a provider or chat decision, verifies the retained evidence,
and consumes an approval exactly once. The provider receives the room path; the
agent does not construct request metadata or output paths.

## Agent journey

1. Run `spacedock gate prepare <ticket> --question ... --artifact ... --summary ...`.
2. Present the returned room through the selected provider or chat.
3. Run `spacedock gate record <ticket> --room <room>` for retained provider
   evidence, or record the captain's chat decision.
4. Run `spacedock gate consume <ticket>` after approval.

The prepared room starts with two authoritative files: `request.json` and the
canonical Briefing. It copies no source payloads. Each source remains pinned by
a Git-root locator, revision, and SHA-256 digest. A provider may add retained
evidence beneath `provider/`.

## Validation evidence

- Fresh validation passed all five acceptance criteria at the exact clean
  candidate.
- `go test ./...`, `go test ./... -race`, formatting, diff checks, and strict
  documentation checks passed.
- Real Codex `gate-guardrail` and `recorded-gate-lifecycle` journeys passed
  locally at the exact candidate in 123.14 seconds and 196.96 seconds.
- PR #573 CI passed offline, docs, macOS install, Ubuntu install, Pi live, Codex
  live, and Opus live.
- The final correction removed 415 lines of gate-review prose parsing. Gate
  tests now judge prepared authority, state, ordering, cardinality, durable
  effect, and Git ancestry instead of model wording.

## Finding disposition

Sonnet live failed only `keep-moving-posture`: it stopped at one gate without
surfacing a separate ticket's in-flight re-shape. The s4 diff contains no
keep-moving change, and Codex, Pi, and Opus passed. This material runtime
finding does not belong to s4; route it to the existing keep-moving reliability
work. It does not weaken the evidence for provider-neutral gate preparation.

## Decision

Approve s4 and merge PR #573. Preserve the Sonnet run as evidence for the
separate keep-moving reliability ticket.
