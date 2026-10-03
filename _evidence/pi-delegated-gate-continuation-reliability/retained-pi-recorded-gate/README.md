# Retained Pi recorded-gate inputs (9w)

Preserved before the source worktree disappears. Copied verbatim; no redaction was applied.

Source worktree (read-only, may vanish):
`/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-live-lanes-red-on-every-branch/live-artifacts/local-proof/`

These are the artifacts `pi-delegated-gate-continuation-reliability` (9w) names in
its Retained evidence and portability section, matched by SHA-256 against the
values that section records.

| file | source | sha256 | lines | bytes |
| --- | --- | --- | --- | --- |
| first-root.jsonl | final-ce4ac943-pi-recorded-gate-pinned/pi-recorded-gate-lifecycle/run/sessions/2026-07-28T15-57-34-288Z_019fa972-16d0-7881-8a94-421ffb1695d3.jsonl | 4ecc18637c62143b8cbae7fbf584fa3093145aa3f6d8d46daa31f4343dfe85de | 79 | 249019 |
| retry-root.jsonl | final-ce4ac943-pi-recorded-gate-pinned-retry/pi-recorded-gate-lifecycle/run/sessions/2026-07-28T16-03-39-912Z_019fa977-ab08-773b-ac42-334e8babde6b.jsonl | b25a248167a200b7a5ac80ac7034fa2fbdb1377c236d52fd32389ddcea5dfe72 | 42 | 216274 |
| first-child.jsonl | final-ce4ac943-pi-recorded-gate-pinned/pi-recorded-gate-lifecycle/run/sessions/2026-07-28T15-57-34-288Z_019fa972-16d0-7881-8a94-421ffb1695d3/8a89cecb/run-0/session.jsonl | 03f9a2676fd6b0c64787b2277932bbe86a136b6dd39ab4647ef0e8fadda9b6df | 21 | 21619 |

Notes carried from the entity body:

- `first-root.jsonl` is the run whose review omits its decision facts; same event acknowledges inventing a digest before reading.
- `retry-root.jsonl` reads canonical sources and ends at the review; its go-test log fails 'Pi child sessions=0, want exactly one'.
- `first-child.jsonl` is the only child session; the retry run has none.
- The sibling stderr files hold launcher banners only. Do not use stderr alone as root-review proof.

Intended repository testdata home, to be placed by a worker under product code:
`internal/ensigncycle/testdata/pi-recorded-gate/`.
