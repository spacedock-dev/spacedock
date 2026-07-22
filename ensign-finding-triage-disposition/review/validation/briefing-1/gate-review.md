# Validation gate: finding-triage decline disposition

Recommendation: **APPROVE** exact candidate `e85eb0cfcc3c243fd94754be2baafa23be302a21` for landing.

## Capability

The candidate gives ensigns a bounded, trigger-delivered rule for triaging review findings before fixing them. Material findings remain fixable, correct-but-disproportionate findings can be preserved as explicit advisory declines with promotion conditions, needs-decision findings escalate, and narrowing a value AC remains a captain-owned binding decision. The shape rides the landed gate-resolution vocabulary without adding recorder plumbing, schema, or product code.

## Validation evidence

- Fresh validation passed all 12 checklist items and AC-1/AC-2/AC-3 at exact clean head `e85eb0cf`; the implementation worktree remained clean.
- The required routed live replay discriminated the two value paths: Case A recorded a linked advisory decline with zero product-line change, while material Case B produced a non-zero product fix; neither advisory path advanced entity status.
- The offline checker accepted eight valid fixtures, rejected both red controls, rejected an unknown class, and failed under five independent claim-breaking mutations.
- The detached adversarial audit was clean. Focused checks, `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, and `git diff --check` all passed.
- The exact captain-authorized prompt exception is bounded to +703 bytes on each host. The final branch has 95 insertions across six paths: zero production LOC, 47 fixture/test lines, and 48 docs/skill/spec lines.
- Roborev jobs 548, 554, and 560 are durably triaged. Material evidence defects were fixed; job 560's all-declines outcome is recorded as `0 fixed; 3 declined`, distinct from no findings.

## Deferred risks

- Malformed or externally supplied fixture hardening is deferred until external fixtures are supported or a schema AC requires it.
- Automatic CI wiring is deferred until release proof requires that checker to run outside the focused validation path.
- Advisory `decision` consumer ambiguity is deferred until generic round plumbing branches on that field.

No deferred risk is material to the candidate's declared value ACs today.

## Decision

- `approve`: authorize landing and terminal completion for the exact candidate.
- `revise`: return a concrete material finding to implementation.
- `hold`: retain validation for a named prerequisite.

