---
title: Harden the edge-advance verify poll with actions read permission
status: backlog
source: "auto-pre0 validation (#553) — the pre0 verify-or-fail poll's Actions-API read works only because the repo is public (edge-advance job has actions: none). Deferred hardening."
id: bd94c7tey7em06xr8qsbpa27
---

Add `actions: read` to the `edge-advance` job's `permissions:` in `.github/workflows/release.yml` so the verify-or-fail poll's Actions-API run lookup works with a restricted token rather than relying on the repo staying public. Cheap (~1 line). NOTE: touches a workflow file, so the branch must be pushed over SSH (the local keychain OAuth token lacks `workflow` scope).
