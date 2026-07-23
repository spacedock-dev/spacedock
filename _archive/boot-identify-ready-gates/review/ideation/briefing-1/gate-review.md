# Gate review: Expose ready-gate entities in boot identify JSON — ideation

Chosen direction: derive readiness from the selected current-stage gate attempt and application, with a four-field boot scheduling index.

Recommend **approve**.

## Checklist

- DONE: Define the durable readiness reducer.
- DONE: Specify boot, status, and lifecycle ownership.
- DONE: Replace stage-only fixtures and isolate the land plan.

Assessment: 3 done, 0 skipped, 0 failed.

## Evidence

- The five-ticket fixture yields exactly three actionable gates and excludes two still validating.
- AC-1 through AC-6 each have explicit Stage Report evidence.
- `ready_gates` contains only `id`, `slug`, `current`, and `readiness`; full canonical detail remains behind entity read and `gate validate`.
- Human readiness distinguishes `validating`, `awaiting-captain`, `approved-awaiting-merge`, and nonterminal `approved-awaiting-advance`.
- Same-Briefing binding repairs a stale old-stage selection without duplicating an attempt.
- Implementation starts from current main, cherry-picks only counterexample commit `c5a96678`, and must exclude unrelated PR #493 history.

Decision: approve implementation of the durable projection and selection correction, coordinated with `6y`, `xb`, and h1.
