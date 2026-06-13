---
id: sryzghzqazj9s9km6ebqkf5s
title: dispatch reconcile conflates team-management with repo-hygiene (and hardcodes pre-flip trunk `next`)
status: backlog
source: "0202 Commander drive (2026-06-13). Boot reconcile flagged Class-D/E drift against origin/next; Class-E remedy 'reset main->origin/next' would have reverted the entire post-flip trunk. Investigation (captain-prompted) found the deeper cause: reconcile bundles git-hygiene into a team-management helper, so it carries repo/trunk knowledge it shouldn't."
group: cleanup
---

`dispatch reconcile` is two helpers in one coat: roster/team reconciliation AND repo git-hygiene. The repo half hardcodes the pre-flip trunk `next`, which the 2026-06-08 flip silently invalidated.

## Problem

`internal/dispatch/reconcile.go` emits five drift classes:
- **A** (lingering agent), **B** (superseded agent) — genuine team management, sourced from the `~/.claude/teams` roster.
- **C** (un-advanced PR) — entity/PR state.
- **D** (stale branch), **E** (stale local main) — pure git hygiene, hardcoding `origin/next` (`reconcile.go:582`, `:605`, `:616`).

reconcile landed 2026-06-02 (#273), six days BEFORE the 2026-06-08 flip — when `next` was genuinely the integration trunk and `main` tracked it, so D/E were correct then. The flip inverted the model (`main` is now the trunk; `next` is dev-only per `docs/releasing.md`), but the helper was never refit. Post-flip, Class-E's remedy ("reset main->origin/next") would throw away the real trunk and revert to the stale dev branch — the dangerous drift the 0202 Commander refused to act on at boot.

Root cause is the conflation: a helper whose name + roster-loading say "team management" should not carry repo/trunk knowledge. The hardcoded `next` is the symptom; the bundling is the cause.

## Proposed approach (ideation to firm)

Pick one (ideation decides):
- (a) **De-conflate**: split repo-hygiene (D/E) out of `dispatch reconcile` into a separate repo-sync check; reconcile keeps A/B/C (roster + entity state).
- (b) **At minimum parameterize the trunk**: source the trunk branch from workflow/README config (or a flag) instead of the hardcoded `next`, so D/E follow the real trunk post-flip.

Either way, D/E must reference `main` (the post-flip trunk) on this repo, and the team-helper must stop hardcoding a repo branch.

## Acceptance criteria (sketch)

**AC-1 (sketch) — reconcile's repo-hygiene drift references the configured trunk, not a hardcoded `next`.**
Verified by: a reconcile unit/fixture test where a fixture trunk other than `next` (or `main`) drives Class-D/E detection, proving the trunk is config-sourced — expected value from the fixture's configured trunk, not a literal in the code under test.

## Notes
Surfaced by the 0202 drive. Sibling finding: the `pr-merge` mod hardcodes base `next` the same way (separate seed). Both are post-flip stale-trunk debt in the dispatch helpers.
