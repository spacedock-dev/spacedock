---
title: "state ready/sweep/commit auto-discover the lone workflow (drop the --workflow-dir requirement)"
status: ideation
source: "FO boot-ergonomics friction 2026-07-01 (Claude Commander session hit it twice; Pi FO Stage B report shows the '(x2 with --workflow-dir)' annotation). `spacedock state ready`/`sweep`/`commit` fail with 'no workflow here — pass --workflow-dir' when run from the repo root, while `status --boot`/`status --discover`/`spacedock new` already auto-discover the lone commissioned workflow. The inconsistency costs a failed-then-retried call at boot (doubling Stage-B state-read tokens) and is a recurring papercut."
group: tooling
id: 82edd88rq11q2f05z5nhfhj8
started: 2026-07-02T01:23:59Z
---

## Problem
`spacedock state ready`, `state sweep`, and `state commit <slug>` require an explicit `--workflow-dir` when run from the project root — they do NOT auto-discover the single commissioned workflow the way `status --boot`, `status --discover`, and `spacedock new` already do. At FO boot this surfaces as a `no workflow here — pass --workflow-dir` error on the first attempt, then a retry with the flag: wasted round-trips and boot-context tokens (the Pi FO Stage B report's `(x2 with --workflow-dir)` annotation is this friction).

## Desired direction (for ideation to refine)
`state ready`/`sweep`/`commit` (and any other `state` subcommand with the same gap) auto-discover the lone workflow when exactly one is commissioned — the same discovery `status --boot`/`--discover`/`new` use — and require `--workflow-dir` only when discovery is ambiguous (>1) or zero, matching `spacedock new`'s established disambiguation behavior.

## Rough acceptance sketch (ideation tightens into measured ACs + a test plan)
- From the repo root of a single-workflow project, `state ready`/`sweep`/`commit <slug>` succeed with NO `--workflow-dir` (behavioral test: exit 0, the state operation performed) — the boot retry-double-call vanishes.
- With >1 commissioned workflow, the subcommands report the candidates and require `--workflow-dir` (parity with `spacedock new`'s ambiguity behavior) — a two-workflow test.
- No regression to the explicit `--workflow-dir` path.

## Related
- Sibling boot-ergonomics item: `status --read --json --frontmatter` boot-lean projection.
- Discovery precedent to reuse: `status --boot`/`status --discover`/`spacedock new` auto-discovery.
