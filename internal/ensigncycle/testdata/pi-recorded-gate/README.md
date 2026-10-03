# Retained Pi recorded-gate inputs

Preserved Pi session traces for the `recorded-gate-lifecycle` journey on the
`se0 exact-tip` Pi proof. The retained first run invented the Briefing digest
before reading it; the single authorized retry read the canonical snapshot but
stopped at presentation with no child session. Neither is a clean journey, so
these traces are negative evidence: the deterministic grade must reject both for
their specific obligation. No live passing trace is retained here.

## Provenance

Source worktree (read-only, may vanish):

```text
/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-live-lanes-red-on-every-branch/live-artifacts/local-proof/
```

Copied verbatim; no redaction was applied. Any paths embedded in the traces are
the original ephemeral `/private/var/folders/...` and worktree paths, retained as
evidence and never resolved by the tests. The checksums below match the values the
`pi-delegated-gate-continuation-reliability` entity records, and the grade reads
its expected values (canonical id/digest, model, marker) from the traces, not from
constants.

| file | source coordinate | sha256 | lines | bytes |
| --- | --- | --- | --- | --- |
| `first-root.jsonl` | `final-ce4ac943-pi-recorded-gate-pinned/pi-recorded-gate-lifecycle/run/sessions/2026-07-28T15-57-34-288Z_019fa972-16d0-7881-8a94-421ffb1695d3.jsonl` | `4ecc18637c62143b8cbae7fbf584fa3093145aa3f6d8d46daa31f4343dfe85de` | 79 | 249019 |
| `retry-root.jsonl` | `final-ce4ac943-pi-recorded-gate-pinned-retry/pi-recorded-gate-lifecycle/run/sessions/2026-07-28T16-03-39-912Z_019fa977-ab08-773b-ac42-334e8babde6b.jsonl` | `b25a248167a200b7a5ac80ac7034fa2fbdb1377c236d52fd32389ddcea5dfe72` | 42 | 216274 |
| `first-child.jsonl` | `final-ce4ac943-pi-recorded-gate-pinned/pi-recorded-gate-lifecycle/run/sessions/2026-07-28T15-57-34-288Z_019fa972-16d0-7881-8a94-421ffb1695d3/8a89cecb/run-0/session.jsonl` | `03f9a2676fd6b0c64787b2277932bbe86a136b6dd39ab4647ef0e8fadda9b6df` | 21 | 21619 |

Notes:

- `first-root.jsonl`: the captain-facing review invents a digest; the canonical
  snapshot read arrives later in the same stream. The lifecycle itself completes
  (record, consume, one successor, completion, durable-report read).
- `retry-root.jsonl`: reads the canonical snapshot and presents a correct bound
  review, then stops before recording the approval. No child session exists.
- `first-child.jsonl`: the only child session, for the first run.
- Sibling stderr files held launcher banners only, not the session trace, so they
  are not retained.

## Reproduction

```bash
go test ./internal/ensigncycle -run 'TestPiRecordedGate|TestRecordedGateLifecycle' -count=1
```

To verify the retained bytes:

```bash
cd internal/ensigncycle/testdata/pi-recorded-gate && shasum -a 256 first-root.jsonl retry-root.jsonl first-child.jsonl
```
