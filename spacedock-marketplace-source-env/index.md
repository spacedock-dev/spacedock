---
title: Add SPACEDOCK_MARKETPLACE_SOURCE override for Codex install
status: backlog
score: ""
source: captain request after local Codex marketplace setup
priority: medium
id: z2tjv3570ahjxewv1c309rbc
---

## Problem

`spacedock codex` auto-installs the Spacedock plugin from the hardcoded marketplace source when the canonical `spacedock@spacedock` plugin is missing. That is awkward for local development: a developer can install a local marketplace manually, but the wrapper has no direct override for the marketplace source it uses on the auto-install path.

Add a development override, tentatively `SPACEDOCK_MARKETPLACE_SOURCE`, so the wrapper can install from a local marketplace path or alternate marketplace source without changing `CODEX_HOME` and without editing global Codex config by hand.

## Acceptance criteria

**AC-1** `spacedock codex` and `spacedock install --host codex` honor `SPACEDOCK_MARKETPLACE_SOURCE` on the Codex auto-install/install path, while preserving the current default `spacedock-dev/marketplace` behavior when the env var is unset.

**AC-2** The override is covered by tests that observe the actual install source passed to the Codex install seam; tests must include unset/default and env-override cases.

**AC-3** The override is documented in the Codex/front-door development guidance, including the local-marketplace use case and the caveat that `--plugin-dir` bypasses installed-plugin resolution but does not solve launcher safehouse wrapping.

**AC-4** Validation includes `go test ./internal/cli` and any focused tests added for the install/front-door path.
