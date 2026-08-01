# z5 validation gate review

Entity: `codex-launch-multi-agent-v2`
Stage: `validation`

The fresh Sol/medium validator completed the required full, race, focused,
formatting, credential-gated, and detached-audit checks. The report has three
DONE items and three FAILED items. AC-1 through AC-4 are cited, but the
failures are material.

Findings:

- The lifecycle test calls `codex exec` directly and accepts unordered marker
  presence; a zero-worker fake passes, so the Spacedock front door, order,
  cardinality, same-worker identity, and v2 context are not proven.
- Codex accepts attached short and quoted dotted reserved-key overrides after
  the launcher-owned layer, so `agents.enabled=false` can still downgrade the
  guarantee.
- The disabled-control negative (E-3) is absent.

Science Officer advisory: REVISE. I agree with the classification and route.

Recommendation: REVISE through validation feedback to implementation. The
correction must canonicalize every accepted reserved-key spelling before
side-effects, drive the built `spacedock codex` front door from an isolated
home, parse typed ordered same-worker records, add the disabled negative, and
rerun the required validation suite without weakening the oracle.
