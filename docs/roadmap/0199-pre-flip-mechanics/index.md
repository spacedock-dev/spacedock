# Sprint 0199 — pre-flip hardening

**Goal:** ship the final pre-flip hardening as **0.19.9** before the 0.20.0 flip — Linux
binaries (`v3`) and dev-tooling quality (`th`, `jm`, `m1`) — so the 0.20.0 flip (`pj`) is a
clean cut on a broadly-installable product.

**Deliverable:** spacedock **0.19.9** cut on `next`.

## Members

Membership is the query, not this table:

```bash
spacedock status --workflow-dir docs/dev --where sprint=0199-pre-flip-mechanics --where 'sprint-readiness != defer'
```

| Entity | group | gate | what it delivers |
|--------|-------|------|------------------|
| `v3` ship-linux-binaries | distribution | **ideation** | goreleaser linux tarballs + a `curl \| sh` installer + the Linux install doc |
| `th` safehouse-preserves-spacedock-bin | dev-quality | **ideation** | safehouse-wrapped launch keeps `SPACEDOCK_BIN` (or the documented allowlist) |
| `jm` entity-label-localization | dev-quality | **ideation** | operating voice speaks the README `entity-label` + cross-workflow `{wf}#{ref}` qualification |
| `m1` rtk-stale-git-audit-guard | dev-quality | **ideation** | audit/validation catches rtk-stale-git via an un-proxied SHA verify |

All four are gated at ideation.

**Out of this sprint:**
- `vh` survey-skill-correctness-pass → **moved to 0.19.8** (the Commander's sprint). It is already fully ideated (cycle-1 + cycle-2 verified, build-ready), so it ships with 0.19.8 rather than waiting for 0.19.9. Its frontmatter re-stamp (`0199 → 0198`) is the Commander's, single-writer, to avoid a state-race.
- `k6` two-channel-release-devbranch-stamp → **moves to the flip (`pj`)** — flip-mechanics, not a pre-flip cut.
- `xp` cross-session-fo-commander-comms → a **separate design spike** (bidirectional channel).
- `pj` the 0.20.0 flip itself (the next release).

## Definition of Done

1. Every active member (`v3`, `th`, `jm`, `m1`) `done` / PASSED + merged to `next`.
2. `go test ./...` from the repo root green.
3. `v3`'s Linux path proven: `goreleaser --snapshot` produces `linux/amd64`+`linux/arm64` tarballs AND the installer fetches+checksum-verifies+installs a runnable `spacedock` on Linux and macOS.
4. spacedock **0.19.9** stamped + cut on `next` (captain-gated release).

## Out of scope

The 0.20.0 flip (`pj`) and everything flip-coupled — including `k6` (the 2-channel brew, the `next→main` `devBranch` retarget, next-publish), which moves to the flip. The `vh` survey pass (moved to 0.19.8). The `xp` cross-session channel (separate spike). `rtk` itself (the operator's tooling — `m1` is the spacedock-side defense only).

## Status

**Ideation dispatched.** Sprint carved + membership final (`v3`/`th`/`jm`/`m1`). Each member's
ideation exercises its riskiest mechanism first (`v3` goreleaser-snapshot; `th` safehouse
env-propagation; `m1` reproduce an rtk proxied/un-proxied SHA divergence; `jm`
render-over-fixture). Ideation gates are presented to the captain as each completes; the package
([`dispatch-sprint-execution.md`](dispatch-sprint-execution.md)) is written last and is the
cold-boot Commander dispatch prompt.
