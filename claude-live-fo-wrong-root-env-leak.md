---
title: "Live claude-live FO boots the wrong root — a CI env leak lures the FO off its launch cwd"
status: ideation
source: "0240/0qt PR #446 CI (2026-06-30): claude-live (sonnet) failed TestLiveClaudeSharedScenarios twice at claude_live_runner_test.go:122 — 'FO booted the wrong root: expected the fixture root .../001, but the boot command targets /home/user/spacedock-workflow (run 1) / /tmp (run 2) — a CI env leak likely lured the FO off its launch cwd.' claude-live OPUS PASSED the identical suite both runs, so it is sonnet-lane env/model behavior, not the scenarios or a product bug. Recurred 2/2; forced an admin-override merge of 0qt."
id: 2wbxv8hdq5m754h45ehfd75j
started: 2026-07-02T01:23:59Z
---
The live `claude-live` shared-scenarios harness (`internal/ensigncycle`, `claude_live_runner_test.go`) expects the dispatched FO to boot in the per-scenario fixture root (`/tmp/TestLiveClaudeSharedScenarios…/001`). On the sonnet lane the FO's first boot command instead `cd`s to a leaked path OUTSIDE the fixture — `/home/user/spacedock-workflow` (run 1), `/tmp` (run 2) — and `claude_live_runner_test.go:122` flags "a CI env leak likely lured the FO off its launch cwd."

It recurred on both PR #446 runs while the opus lane passed the identical suite, so it is sonnet-lane environment/model behavior, not a scenario or product defect — and distinct from the handoff's named "sonnet headless-gate" flake.

Impact: it blocks the sonnet live lane for ANY PR (it forced the 0qt admin-override merge). Likely fix: have the harness pin/assert the FO's launch cwd (or scrub the leaking env var) so a leaked cwd cannot lure the FO off the fixture root. Surfaced during the 0240 sprint's 0qt merge; filed at captain request.
