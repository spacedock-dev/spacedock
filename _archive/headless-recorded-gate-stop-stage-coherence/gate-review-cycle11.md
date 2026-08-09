# Validation gate: stage-coherent headless recorded-gate stop

## Capability

The shared Sonnet headless journey now starts at a real queued initial stage, dispatches implementation once, records its implementation report, enters gated validation, prepares and commits one canonical validation Briefing, and stops with the attempt open and unresolved. The test observes provider-neutral state and command evidence; it does not parse model transcripts or provider output.

## Candidate

- Reviewed commit: `aeb3009b005ff16318add590cb0b0a58a81cfd27`
- Surface: three test files, 88 insertions and 12 deletions
- Product command, storage, authority, and runtime semantics: unchanged

## Evidence

- Focused queued-fixture, direct-gate transition/authority mutants, and success/failure evidence-retention tests passed.
- `go test ./...`, `go test ./... -race`, formatting, and diff checks passed.
- The retained cycle-10 supported Sonnet workflow regrades green against the final oracle: implementation dispatch/report → validation transition → one prepare/commit → open unresolved attempt.
- The final delta after that live run changes only the provider-neutral oracle/mutants (`+8/-11`); fixture, prompt, runner, and retention code are byte-unchanged, so no redundant model run was spent.

## Findings and correction history

No Material, deferred-risk, or polish findings remain. Validation caught and removed an invented validation-worker/report requirement that was not present in AC-2. The final candidate is re-anchored to the literal accepted value and is smaller than the over-specified candidate.

## Recommendation

Approve validation and consume the authorization to update PR #583 and run required CI.
