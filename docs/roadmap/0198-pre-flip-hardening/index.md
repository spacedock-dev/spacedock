# Sprint 0198 — pre-flip hardening (qa-led)

**Goal:** ship the binary/version/install UX (`qa`) plus pre-flip hardening and survey
polish as **0.19.8**, before the 0.20.0 flip — so a user hitting a missing / incompatible /
stale binary gets a helpful journey (versions, not "contract" jargon), the survey reads
cleaner, and the release path is gated by e2e.

**Deliverable:** spacedock **0.19.8** cut on `next` with the members landed.

## Members

Membership is the query, not this table:

```bash
spacedock status --workflow-dir docs/dev --where sprint=0198-pre-flip-hardening --where 'sprint-readiness != defer'
```

| Entity | group | what it delivers |
|--------|-------|------------------|
| `qa` spacedock-binary-missing-install-journey | binary-ux | **headline** — install/upgrade journey + versions-not-"contract" messages (missing / incompatible / stale) |
| `jh` dispatch-build-flag-form-version-skew | binary-ux | same-contract command-form skew detection (needs ideation + spike) |
| `1p27` survey-scaffold-state-the-fact | survey | SCAFFOLD section states the observed fact, drops the recovered-vs-installed taxonomy |
| `69rk` survey-codex-cwd-workaround | survey | survey-side fallback + hint for Codex sessions agentsview leaves cwd-less |
| `kbr8` migration-check-prune-state-walk | test-hygiene | prune `.spacedock-state` in the migration-check walk (019x audit Material) + drop orphaned survey fixtures |
| `nzb7` gate-release-on-e2e | release-gating | gate the release/flip on the live e2e suite |

> CL is filing more install-path members from manual testing — they tag into this sprint as they land.

## Definition of Done

1. Every `ready` member `done` / PASSED + merged to `next`.
2. `go test ./...` from the repo root green with `.spacedock-state` + `.claude/worktrees` present.
3. spacedock **0.19.8** stamped + cut on `next` (captain-gated release).
4. `qa`'s behavior proven by a **captain-run live drive** at sprint acceptance — drive `spacedock doctor` / `spacedock claude` against a too-old-binary / too-old-plugin manifest and observe the version-bearing, jargon-free messages (mirroring the 019x AC-3 live drives the captain ran). qa's offline ACs prove the *mechanism*; this captain live drive proves the *messages*. (Resolves the preflight's DoD#4-unowned blocker — the live drive is a sprint-level captain step, not an offline member AC.)

## Out of scope

The 0.20.0 flip (`pj`, next sprint). `5h0` (blocked on #315). Fixing agentsview upstream — the Codex-cwd gap is worked around survey-side, not fixed in agentsview.

## Status

Shaping. Ideation dispatched for the design-fork members (`nzb7`, `69rk`, `jh` — the last with a spike); `kbr8`/`1p27` are well-specified (straight-to-implementation candidates, per the `78` precedent); `qa` is at its ideation gate awaiting captain approval. The staff readiness review + the Commander dispatch package follow once ideation settles.
