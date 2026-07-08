---
id: mzkfah6m06hyh972ba2c9mys
title: "No contract rule requires tracing ownership before acting on shared GitHub state (PRs/CI runs outside this workflow)"
status: backlog
source: "FO self-diagnosis, 2026-07-08 live session. Asked to \"approve relevant ci for all pending ones in the right sequence (of our session's merge train),\" the FO ran `gh pr list --state open` over the WHOLE repository, found two PRs with WAITING live-CI environment approvals (#474, #435, the latter with a \"stacked on #435\" comment giving a plausible-sounding sequence), and approved four GitHub Environment deployments on them. Neither PR traced to a spacedock-tracked entity in this workflow's `.spacedock-state`, and neither branch matched the `{worker_key}/{slug}` naming convention this FO's own dispatches use (one was even under a different person's branch namespace, `iamcxa/status-apply-gate`). The captain: \"why would you think you are approving ci for those? they are not even tracked locally as spacedock entities.\" The approvals cannot be revoked (GitHub has no un-approve API for a granted environment deployment) — a bounded but real cost paid for acting outside scope."
started:
completed:
verdict:
score:
worktree:
issue:
---

**Problem:** the contract has a strong rule for pushing/opening a PR on the FO's OWN tracked entities (`pr-merge.md`'s "Do NOT push or create a PR without explicit captain approval") and a strong rule for verification rigor ON tracked entities (the z25 self-evidence-bar Working Principle). It has nothing in between: no stated requirement to verify that a PR, branch, or CI run actually belongs to this workflow BEFORE evaluating its relevance at all. `gh pr list` / `gh run list` return repo-wide state — every other session's and every other contributor's in-flight work — not this FO's scope, and nothing in the contract says so.

**Cause:** without an explicit ownership-trace checkpoint, a coherent-looking search result (WAITING CI, a plausible "stacked on #435" dependency comment) reads as sufficient evidence of relevance on its own — coherence was substituted for ownership. Compounding this, the existing "ask the human when materially ambiguous" clarification trigger is worded specifically around *dispatch* decisions, so it never fired for a non-dispatch action (CI approval) touching shared infrastructure outside the tracked-entity boundary.

**Recommended fix:** add a Working Principle: before approving, merging, or otherwise mutating state on any PR, branch, or CI run, verify it traces to an entity in this workflow's `.spacedock-state` or a `{worker_key}/{slug}` branch this FO itself dispatched. A `gh` query returning repo-wide results is never itself proof of scope. Separately, broaden the "ask when materially ambiguous" clarification trigger's framing beyond dispatch specifically, so it also covers any outward-facing or shared-infrastructure action (CI approval, merges, deployments) that isn't obviously inside the tracked-entity boundary.
