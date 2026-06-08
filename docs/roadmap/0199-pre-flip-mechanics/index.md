# Sprint 0199 — pre-flip hardening

**Goal:** ship the final pre-flip hardening as **0.19.9** before the 0.20.0 flip — Linux
binaries (`v3`) and dev-tooling quality (`th`, `jm`) — so the 0.20.0 flip (`pj`) is a clean
cut on a broadly-installable product.

**Deliverable:** spacedock **0.19.9** cut on `next`.

## Members

Membership is the query, not this table:

```bash
spacedock status --workflow-dir docs/dev --where sprint=0199-pre-flip-mechanics --where 'sprint-readiness != defer'
```

| Entity | group | gate | what it delivers |
|--------|-------|------|------------------|
| `v3` ship-linux-binaries | distribution | ideation ✓ | goreleaser linux tarballs + a `curl \| sh` installer + the Linux install doc |
| `th` safehouse-preserves-spacedock-bin | dev-quality | ideation ✓ | front-door re-asserts `SPACEDOCK_BIN` in the wrapped argv (survives safehouse's env-scrub) |
| `jm` entity-label-localization | dev-quality | ideation ✓ (downscoped) | **Layer-1 label localization only** (present-gate + commander-dispatch speak the README `entity-label`) |

**Out of this sprint:**
- `m1` rtk-stale-git-audit-guard → **deferred** (captain): disproportionate for an rtk-only, already-caught issue; the FO does the un-proxied compare ad hoc. Ideation banked.
- `vh` survey-skill-correctness-pass → **moved to 0.19.8** (build-ready; the Commander owns it).
- `k6` two-channel-release-devbranch-stamp → **moves to the flip (`pj`)** (flip-mechanics).
- The cross-workflow `{wf}#{ref}` qualifier (jm's AC-4) → **deferred follow-up** (design banked in jm's body).
- `xp` cross-session-fo-commander-comms → separate design spike. `pj` the 0.20.0 flip itself.

## Definition of Done

1. Every active member (`v3`, `th`, `jm`) `done` / PASSED + merged to `next`.
2. `go test ./...` from the repo root green.
3. `v3`'s Linux path proven: `goreleaser --snapshot` produces `linux/amd64`+`linux/arm64` tarballs AND the installer fetches+checksum-verifies+installs a runnable `spacedock` on Linux and macOS.
4. `th`'s re-assert proven against **real safehouse by the captain** (off-box — this dev environment cannot run real safehouse; the fake-smoke is not sufficient).
5. spacedock **0.19.9** stamped + cut on `next` (captain-gated release).

## Out of scope

The 0.20.0 flip (`pj`) and everything flip-coupled (`k6`). The `vh` survey pass (0.19.8). The `m1` rtk guard (deferred) and the `xp` cross-session channel (spike). The cross-workflow qualifier (jm follow-up).

## Status

**Shaped, ideated, gated — ready for the Commander handoff.** `v3` + `th` ideation captain-approved;
`jm` approved downscoped (Layer-1 label only; cross-workflow qualifier deferred); `m1` deferred;
`vh` moved to 0.19.8. The package
([`dispatch-sprint-execution.md`](dispatch-sprint-execution.md)) is the final cold-boot Commander
dispatch prompt: per-task build notes, the high-stakes detached-audit + th-real-safehouse gates,
and the 0.19.9 release-cut recipe. A separate Commander session drives `implementation → done`.
