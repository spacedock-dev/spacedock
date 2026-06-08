# Sprint 0199 — pre-flip mechanics

**Goal:** ship the final pre-flip hardening as **0.19.9** before the 0.20.0 flip — the
flip-mechanics (`k6`), Linux binaries (`v3`), and dev-tooling quality (`th`, `jm`, `m1`) —
so the 0.20.0 flip (`pj`) is a clean, low-risk cut.

**Deliverable:** spacedock **0.19.9** cut on `next`.

## Members

Membership is the query, not this table:

```bash
spacedock status --workflow-dir docs/dev --where sprint=0199-pre-flip-mechanics --where 'sprint-readiness != defer'
```

| Entity | group | gate | what it delivers |
|--------|-------|------|------------------|
| `k6` two-channel-release-devbranch-stamp | flip-mechanics | **ideation + staff review** | per-channel `devBranch` stamp so a binary resolves its channel's plugin (the flip-mechanics linchpin) |
| `v3` ship-linux-binaries | distribution | **ideation** | goreleaser linux tarballs + a `curl \| sh` installer + the Linux install doc |
| `th` safehouse-preserves-spacedock-bin | dev-quality | **ideation** | safehouse-wrapped launch keeps `SPACEDOCK_BIN` (or the documented allowlist) |
| `jm` entity-label-localization | dev-quality | **ideation** | operating voice speaks the README `entity-label` + cross-workflow `{wf}#{ref}` qualification |
| `m1` rtk-stale-git-audit-guard | dev-quality | **ideation** | audit/validation catches rtk-stale-git via an un-proxied SHA verify |

All five are gated at ideation — none is as clear-cut as a fast-track (no-gate) item. `k6` additionally earns an independent staff review (high-stakes front-door / release machinery).

**Out of this sprint:**
- `xp` cross-session-fo-commander-comms → **a separate design spike**, not a 0.19.9 deliverable. It is a bidirectional-channel design (the receive-half mutual-injection round-trip is its first spike); it feeds a later build rather than bloating this cut.
- `pj` the 0.20.0 flip itself (the next release).

## The k6 boundary (load-bearing — pins 0.19.9 vs. flip)

`k6` is flip-mechanics, but **only its stamp mechanism lands in 0.19.9**; the channel split executes at the flip:

- **0.19.9 (this sprint):** the binary resolves `devBranch` from a per-channel **stamped** value instead of the hardcoded `frontdoor.go:49` constant — proven by building with the stamp and confirming the channel-correct plugin installs against a fresh HOME. Stamping `devBranch=next` for the current single channel is a no-op in behavior; the *mechanism* is what ships and is tested.
- **0.20.0 flip (`pj`):** the 2-channel brew (stable→`main` / edge `spacedock@next`→`next`), the `next→main` retarget, and the next-publish step USE that mechanism.

Ideation pins the exact 0.19.9-vs-flip boundary and may split `k6` if the brew/retarget half is better carried by `pj`. The 0.19.9 deliverable is the riskiest, front-door-touching half: the stamp.

## Definition of Done

1. Every active member (`k6`, `v3`, `th`, `jm`, `m1`) `done` / PASSED + merged to `next`.
2. `go test ./...` from the repo root green.
3. `k6`'s stamped binary installs the channel-correct plugin — proven by a **live front-door run** against a fresh HOME (not a grep of the constant) — plus a **detached adversarial audit** (high-stakes front-door / release surface).
4. `v3`'s Linux path proven: `goreleaser --snapshot` produces `linux/amd64`+`linux/arm64` tarballs AND the installer fetches+checksum-verifies+installs a runnable `spacedock` on Linux and macOS.
5. spacedock **0.19.9** stamped + cut on `next` (captain-gated release).

## Out of scope

The 0.20.0 flip (`pj`) and everything flip-coupled: the 2-channel brew, the `next→main` `devBranch` retarget, and next-publish (those USE `k6`'s mechanism at the flip). The `xp` cross-session channel (separate spike). `rtk` itself (the operator's tooling — `m1` is the spacedock-side defense only).

## Status

**Shaping in progress.** Sprint carved (members stamped). Next: drive ideation per member —
exercise each riskiest mechanism first (`k6`'s per-channel-stamp resolution is the linchpin
spike; `v3` goreleaser-snapshot; `th` safehouse env-propagation), independent staff review on
`k6`, then present ideation gates. The package
([`dispatch-sprint-execution.md`](dispatch-sprint-execution.md)) is written last and is the
cold-boot Commander dispatch prompt.
