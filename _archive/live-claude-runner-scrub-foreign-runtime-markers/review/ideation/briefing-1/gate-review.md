# Gate review: isolate live-host runtime identity

## Chosen direction

Scrub only the four production host-detector keys at each existing per-host live
environment builder. Preserve credentials, isolated state, and the production rule
that genuinely mixed runtime markers fail closed.

## Artifact summary

This review summarizes the selected harness boundary, the exact six-file/88-LOC
implementation plan, its deterministic and live evidence, and the dependency on landed
6y before implementation.

## Evidence

- A temporary builder probe measured foreign-marker leakage at Claude 3, Codex 2, and
  Pi 2.
- Two retained 6y Claude journeys reproduced Codex+Claude ambiguity before explicit
  `--host claude` recovery.
- Existing production ambiguity-refusal tests remain green and are explicitly
  unchanged.

## Implementation boundary

- Six files under `internal/ensigncycle/`, 88 changed LOC with ±25 tolerance.
- No `internal/dispatch`, CLI, skill, runtime-documentation, or production error change.
- The live proof performs one flag-free `dispatch build` and rejects ambiguity,
  explicit-host recovery, rebuild, or a missing committed successor effect.

## References

- State design/report: state commit
  `9337dbfec1ca5cf21f23e818c1cca0ff455378f0`,
  `live-claude-runner-scrub-foreign-runtime-markers/index.md`.
- Code baseline examined: main commit
  `cc51e518a3420b01fd4b455e9710d38803dc6d3e`.

## Recommendation

Approve the direction, but implement only after 6y lands so the harness repair rebases
over the final recorded-gate oracle without overwriting it.

## Decision

Approve to enter implementation after landed 6y; revise to change the harness boundary;
or hold if explicit-host recovery is acceptable test behavior.
