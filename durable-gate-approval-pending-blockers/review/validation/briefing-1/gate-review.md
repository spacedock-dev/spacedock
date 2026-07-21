# Gate review: durable gate-resolution recorder

Recommendation: **APPROVE validation and proceed to merge/landing.** The recorder is ready: there are no material findings, all seven retained acceptance criteria passed, and the corrected evidence now kills the adversarial failures that caused the first validation rejection.

## What this ships

`spacedock gate record` is the binary-owned writer for durable gate resolutions. It supports `open`, `rebind`, `close`, and `supersede`; `spacedock gate validate` validates and reports the selected record. Recording changes only the entity's `gates` frontmatter, preserves opaque application-owned state, never advances `status` or dispatches, and verifies digest/provider identity before normalizing the provider Briefing id.

## Why the evidence is strong enough

| Outcome | Reproduced proof |
|---|---|
| AC-1, AC-4, AC-6 | Approve/revise/hold survive cold reads with exact identities and digests, obey portable rationale rules, and appear distinctly in text and JSON status. |
| AC-10 | Successful writes are gates-only; digest, pointer, lock, and frozen-state failures leave the entity unchanged. |
| AC-12 | Eight production frontmatters replay exactly, and one fixture carries two logical gates across eight stable attempts with pointer, lineage, application, and extension preservation. |
| AC-13 | Only Review & Gate v1 Resolution fields cross the provider envelope; application and known or future wrapper fields cannot leak. |
| AC-14 | The exact A→B→C→close→new-attempt lifecycle proves rebinding, closure freeze, and supersession lineage. |

Validation cycle 1 correctly rejected a wrapper-field leak and weak lifecycle evidence. Correction `9d279b87` closed both boundaries. A detached audit now turns the full suite red for a success-reporting no-op rebind and disabled supersede, and turns the focused boundary test red for a future-only wrapper-field leak. `go test ./...` and `go test ./... -race` pass; `gofmt -l ./cmd ./internal` is empty.

## Decision

- Material findings: **none**.
- Deferred risk: a future Review & Gate version that adds a required Resolution field will require an allow-list update. That trigger is outside the current v1 promise and becomes material only when Spacedock supports that newer field.
- Approving authorizes the workflow's merge/landing ceremony; that ceremony closes this worktree.

## References

- Recorder product commit: `1095be38`
- Evidence-hardening commit: `9d279b87`
- Validation report state commit: `0c0fb6ca`
- Detailed validation record: `/Users/clkao/git/spacedock-research/spacedock-v1/docs/dev/.spacedock-state/durable-gate-approval-pending-blockers/index.md`
- Product contract on branch `spacedock-ensign/durable-gate-approval-pending-blockers`: `docs/specs/gate-resolution-frontmatter-contract.md`
