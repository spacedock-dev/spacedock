# 3d validation review

Recommendation: **HOLD / REJECTED evidence state**. The 3d candidate satisfies AC-1 through AC-4 and both changed live paths, but the required current-checkout full and race suites remain red on seven stale v1 pilot-manifest paths outside 3d ownership.

## Candidate result

- Commit `17985b60f`; nine files, +401/-170, within the approved tolerance.
- Cycle-1 false-green defect is closed by canonical complete-record validation and workflow-derived routing.
- AC2 positive live branch passed in 194.60 seconds.
- Codex pre-ys `gate-guardrail` passed on `gpt-5.6-luna` with maximum reasoning in 321.98 seconds.
- Detached liveness audit observed 31 ordered JSONL events, one process, and a required 60.44-second stall kill/reap.
- Validation report accounting: 6 DONE, 2 SKIPPED, 3 FAILED.

## Acceptance criteria

- AC-1: satisfied; malformed, duplicate, later-approve, wrong-route, wrong-decision, accepted, and unchanged branches all fail closed.
- AC-2: satisfied; deterministic and live progress evidence stays alive beyond four quiet budgets.
- AC-3: satisfied; silence fails at the 60-second product budget with retained fault evidence.
- AC-4: satisfied; no fixed Codex deadline remains and CI keeps the 40-minute runaway backstop.

## Blocking evidence

Canonical current-checkout `go test ./...` and `go test ./... -race` both fail because seven paths introduced by code commit `9ff2aa50c` now exist only under `_archive` while `v1_pilot_manifest.txt` still names their old active paths. Snapshot `73f41e2a2` predates those moves and is not an acceptable substitute. This is a Material release blocker under `AGENTS.md#Expected Commands`, but it is outside 3d's approved nine-file surface.

Hold 3d at validation until the manifest/current-state owner restores both mandatory suites. Do not route this external repair into 3d implementation.
