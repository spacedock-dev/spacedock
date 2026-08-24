---
id: rb0qr2psxahhndrfarx0avzq
title: Keyless machines fail the plugin install over SSH — surface CLAUDE_CODE_PLUGIN_PREFER_HTTPS
status: backlog
source: "Captain CL fresh-install VM experience report, 2026-08-24: VM with no SSH keys needed CLAUDE_CODE_PLUGIN_PREFER_HTTPS=1 for the marketplace add"
started:
completed:
verdict:
score:
worktree:
issue:
---

`claude plugin marketplace add spacedock-dev/marketplace` uses GitHub shorthand, which clones over SSH by default. A fresh machine with no SSH keys fails the add until `CLAUDE_CODE_PLUGIN_PREFER_HTTPS=1` is set. Verified 2026-08-24: zero mentions of the env var anywhere in this repo or the marketplace README — the HTTPS-not-SSH guidance existed only for the repo clone, so the plugin-install path has an undocumented fresh-machine failure that only appears on keyless hosts.

## Problem

{Ideation fills this in. Seeded: the documented plugin install path fails on exactly the machine class (fresh, keyless) the fresh-install claim targets, and the remedy is discoverable only by already knowing Claude Code internals.}

## Proposed approach

{Ideation fills this in. Two rungs to weigh: (a) document the env var at the point of the marketplace-add command in docs/site/get-started/install.md and the marketplace README troubleshooting; (b) have the launcher fix it instead of documenting it — `spacedock install` / launch-time plugin ensure could set or default the env var (or use an explicit HTTPS URL form) when invoking claude plugin commands, making the doc note unnecessary for the launcher path. Check whether codex/pi hosts have an analogous shorthand-resolution behavior.}

## Out of scope

Marketplace README channel-model repair (separate task). SSH-key provisioning guidance for the repo clone itself.

## Expected surface and tolerance

Estimate net LOC change: ~+15 docs-only for rung (a), or ~+30 across launcher + test for rung (b). Ideation picks.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified. Seeded; ideation refines.

**AC-1 - A keyless fresh machine following the documented install path completes the plugin install without independently discovering the env var.**
Verified by: either the launcher passing HTTPS resolution programmatically (unit test on the env/URL the plugin command is invoked with) or the docs naming the env var inline at the marketplace-add command; fails if the documented path still SSH-fails on a keyless host with no in-path remedy.

## Test plan

{Rung (a): docs review. Rung (b): unit test on the launcher's plugin-command environment/URL construction; one live keyless-profile smoke if cheap.}
