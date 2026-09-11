---
title: "Pi live E2E runs as an associated, manually-approved PR check"
status: backlog
source: "Captain directive 2026-09-11 at the s98 validation gate: 'approve the pi ci lane' + 'the pi live is simply skipped and we can't manually trigger that check to be really associated with the PR'."
id: 6v7sdg4pfzwcfnjstt80pvqj
---
Problem: in `.github/workflows/runtime-live-e2e.yml`, `claude-live` and `codex-live` run on every `pull_request` (behind the `CI-E2E` required-reviewer environment), but `pi-live` is gated `if: workflow_dispatch && inputs.live_cadence == 'pi'` — it never runs on a PR, and a manual `workflow_dispatch` run against a PR branch produces no check run on that PR (the workflow's own journey-delta-comment job notes "a workflow_dispatch run has no PR"). So pi live evidence is neither PR-associated nor manually triggerable per-PR; it only exists as an unassociated manual run or the release-time precondition.

Deliverable: pi live E2E becomes an opt-in, PR-associated check. Proposed shape (mirrors the existing claude/codex job pattern): add `labeled` (and `synchronize`) to the `pull_request` trigger types and run `pi-live` when a `live:pi` label is present — `if: (github.event_name == 'pull_request' && contains-label) || inputs.live_cadence == 'pi'` — with its own environment (e.g. `CI-E2E-PI`, required reviewer) holding the pi secret. The label is the cost opt-in; the environment approval is the human gate; the run appears in the PR's Checks tab and can be a required check where warranted. Alternatives considered: workflow_dispatch + out-of-band commit-status on the PR head SHA (associates but bypasses the Checks-tab model and needs a pr_number input); scheduled nightly pi cadence on main (cheap, but never PR-associated).

Acceptance criteria and test plan to be fleshed out at ideation (workflow file change; proof = workflow-syntax validation plus a live dispatched run of the lane on a scratch PR).
