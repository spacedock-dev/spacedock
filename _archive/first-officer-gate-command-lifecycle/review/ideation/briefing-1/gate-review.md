# Ideation gate review: recorded First Officer gate lifecycle

## Capability and change

Codify the gate procedure already exercised during the durable-decisions sprint: capability-check the selected `spacedock` binary; bind and validate the canonical Briefing; judge and present concise evidence; record and validate the direct, delegated, or provider decision; require eligibility and one-use consume before approval can advance or dispatch.

The implementation changes the First Officer operating contract, behavioral fixtures/live scenarios, and concise reference documentation. It adds no recorder schema, gate command, application behavior, or production Go path.

## Evidence and reviewed snapshot

- 3k/PR #557 landed the canonical recorder and h1/PR #560 landed eligibility/consume, but neither changed `skills/first-officer`; literal command invocation is absent from the current FO contract.
- The sprint successfully exercised bind, open validation, delegated approval recording, and closure validation. Before h1 landed, later state movement remained manual; the design replaces only that bypass with eligibility/consume.
- A fresh current-worktree build exposes all four operations. The retained repo-root `./spacedock` reports a compatible dev version but lacks the gate command surface, proving version checks alone are insufficient.
- The ideation entity records the demonstrated six-event trace, direct-versus-delegated provenance, fail-closed routing, resume/idempotency behavior, all observed friction, and load-bearing skipped-step mutants.

Reviewed artifacts are the exact entity, landed gate contract, and current FO shared core identified by URI and SHA in `briefing.json`.

## Findings

- Material: the ordinary FO path can still advance through prose/status handling without a recorded and consumed authorization.
- Material: a version-compatible but stale executable can pass readiness while lacking the required commands.
- Material: the procedure must distinguish direct Captain authority from an FO decision rendered under delegated conn.
- Deferred to existing owners: Briefing/Result package generators and provider transport/polling; this task supplies exact input guidance and retained fixtures only.
- Deferred to the concurrent advisory-round task: correction-round recording. The two tasks share no skill surface; their command-reference edits must be serialized.
- Polish: presentation channel/UI remains overridable and is not redesigned here.

## Recommendation and decision

Recommendation: **approve ideation and proceed to implementation** within the declared no-production-Go boundary. Require the real retained 3k package replay, a skipped-step mutant that turns red, the capability-stale launcher control, and Pi-live-capable coverage before validation.

Decision requested: approve, revise with concrete scope findings, or hold for a named prerequisite.
