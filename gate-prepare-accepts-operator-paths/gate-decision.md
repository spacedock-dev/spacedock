# Gate decision: make gate paths operator-safe

## Recommendation

Approve ideation. The path-resolution ambiguity has now been reproduced live: a valid project-relative artifact was prefixed with the state root and rejected as a doubled path.

## Why this matters for 0.27.x

- Three of six audited preparations failed this way across two days.
- The gate lifecycle forbids retrying a failed preparation, so one ambiguous path can withhold a gate for the rest of a work window.
- The failure affects both artifacts and references, including workflow documentation.

## Approved scope

- Define deterministic behavior for absolute, cwd-relative, and state-relative Markdown paths.
- Prefer accepting the operator-supplied readable path before applying a root-relative interpretation.
- Align the first-officer wording and binary behavior.
- Cover equivalent path forms with table-driven gate-room tests.

## Exclusions

- Do not change gate-room layout, briefing schema, or the no-retry rule.
- Do not add filesystem-wide searching or permissive symlink handling.

## Proof owed at the next gate

Ideation must specify precedence and ambiguity handling, preserve readable-regular-file and repository-boundary guards, refine the surface estimate, and demonstrate the failing doubled-path baseline.

