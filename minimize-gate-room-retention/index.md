---
title: Define the minimum recoverable gate-room retention contract
status: backlog
source: "Durable-decisions sprint provider-boundary audit, 2026-07-24."
score: "0.6"
id: 9tzbg7982knh3kwt6mvaskcg
---

The gate presentation contract currently treats Result, review log, inventory, diagnostics, argv, exit status, stderr, and provider-specific death markers as permanently required across success and failure. That preserves evidence but may retain implementation debris after an ordinary successful decision.

Ideation should falsify the minimum recoverable set per outcome:

- Success must preserve the exact binding Result and complete presented inventory needed to verify what was authorized.
- Failure must preserve every produced decision byte plus enough diagnostics to distinguish provider failure, retention failure, invalid Result, and launcher/controller death without relaunch.
- Determine whether argv, stderr, diagnostics, and transport markers add recovery value on success or belong only to failed rooms.
- Preserve post-launch no-fallback, no relaunch, and immutable room history.

Reuse the existing provider conformance matrix. Change behavioral expectations, not its language, exact fixture count, repository topology, or a permanent source-commit pin.
