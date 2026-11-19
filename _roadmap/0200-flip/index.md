# Sprint 0200 — main-flip (capstone)

**Goal:** cut spacedock **0.20.0 on `main`** and **flip the marketplace** so users install
the stable release off `main`, with `next` retained as the dev/edge channel. NO contract
change — the contract stays at 1; 0.20.0 is the first real (non-dev) release, not a bump.

**Deliverable:** spacedock **0.20.0** cut on `main`; the marketplace resolves `0.20.0` from
`main`; `next` retained as the edge channel.

## Two phases, two drivers (captain-decided)

- **Pre-flip mechanics → Commander.** `nzb` + `k6d` land on `next` (merge, **no intermediate
  release**) so the flip cut rides the enforced e2e-gate and the per-channel install
  mechanics. A cold-boot Commander drives `implementation → done` from the package.
- **The flip → here (Shaping FO + captain).** `pj` is the capstone, driven in the shaping
  session **after** the Commander lands the pre-flip pair — outward-facing and captain-gated
  at every step (archive, force-push, tag, marketplace, brew), so it is not handed to a
  cold-boot Commander. Marked `sprint-readiness: defer` to keep it out of the Commander's
  drivable query; "defer" here means **driven here, not by the Commander** — not dropped.

## Members

Membership is the query, not this table. The Commander's drivable set:

```bash
spacedock status --workflow-dir docs/dev --where sprint=0200-flip --where 'sprint-readiness != defer'
```

| Entity | group | phase | gate | what it delivers |
|--------|-------|-------|------|------------------|
| `nzb` gate-release-on-e2e | release-gating | pre-flip (Commander) | ideation ✓ (banked) | release-time `e2e-gate` job goreleaser `needs:` — a `v*` cut requires a green live-e2e run for the tagged commit (+ auditable `SPACEDOCK_E2E_GATE_WAIVER`) |
| `k6d` two-channel-release-devbranch-stamp | flip-mechanics | pre-flip (Commander) | ideation pending | per-channel `devBranch` stamp (stable→`main` / edge→`next`), channel-aware `release.yml` version-stamp, two-channel brew, reuse the existing `next-publish.yml` as the edge lane |
| `pj` main-flip-0200-marketplace | flip-capstone | the flip (here) | ideation pending | archive `main`→non-`v*` ref, guarded `--force-with-lease` `next`→`main`, cut 0.20.0, marketplace ref flip (+ paired guard-test edit), doc reconciliation, existing-user migration |

**Out of this sprint:**
- `44` bundle-asset-distribution + `5w` notarize-macos-release → **deferred** (captain):
  packaging/notarization is post-flip, not part of the branch/marketplace flip.
- `xp` cross-session-fo-commander-comms → separate design spike.
- README reconciliation → **already merged** (pre-flip); `k6d` adds the edge-channel doc note
  only if it turns out necessary.

## Definition of Done

**Pre-flip (Commander, on `next`):**
1. `nzb` + `k6d` `done` / PASSED + merged to `next`.
2. `go test ./...` from the repo root green.
3. `nzb` proven: workflow-guard tests assert goreleaser `needs: e2e-gate` + the SHA-gated
   query; the predicate's pass/block/waiver unit tests; an observed/dry-run blocking exercise.
4. `k6d` proven (live, fresh HOME): a binary built per channel resolves the channel-correct
   plugin (stable→`main` plugin, edge→`next` plugin) — observed in workflow/plugin state, not
   string-matched; `goreleaser --snapshot` produces both channel artifacts with the correct
   `devBranch` ldflag each.

**The flip (here, captain-gated at each outward step):**
5. Current `origin/main` archived under a **non-`v*`** ref before replacement.
6. `origin/main` replaced by the prepared 0.20.0 line via guarded `--force-with-lease`;
   `origin/next` retained.
7. spacedock **0.20.0** cut on `main` through the enforced e2e-gate / final live run (unless
   captain-waived) — no contract change.
8. Marketplace ref retargeted `next`→`main` **with the paired `marketplace_manifest_test.go`
   edit** so the line stays green; a real install/resolve exercise confirms `0.20.0` resolves
   from `main`.
9. Doc reconciliation: `AGENTS.md` + `docs/releasing.md` describe the post-flip world,
   verified against the real machinery.
10. Upgrade paths confirmed: `0.12.x`-no-binary, `0.19.x`-skew, and the `next`-pinned-user
    post-flip migration (see `pj`'s flip checklist).

The flip checklist in `pj`'s body (the carve completeness audit, 2026-06-08) is the
authoritative task list for DoD 5–10.

## Out of scope

Packaging (`44` bundle / `5w` notarize), the `xp` cross-session channel, and any contract change.

## Status

**Carving (kickoff, 2026-06-08).** Members stamped `sprint=0200-flip`. Next per the lifecycle:
ideate `k6d` + `pj` (riskiest mechanism first — `k6d`'s per-channel-stamp spike; `nzb`'s
ideation is banked, do not re-ideate), preflight staff review, present ideation gates, then
package the Commander dispatch for the pre-flip pair (`nzb` + `k6d`). The flip (`pj`) is held
for the here-driven capstone phase after the pre-flip pair merges.

## Sprint lifecycle checklist

**Shape — Shaping FO**
- [x] **Scope-lock** with the captain — `nzb` + `k6d` (Commander), `pj` (here); `44`/`5w`/`xp` out
- [x] **Carve** — members stamped `sprint`/`group`/`sprint-readiness`; this `index.md` written
- [ ] **Ideate** `k6d` + `pj` — problem / approach / AC + test-plan, riskiest mechanism first (`nzb` banked)
- [ ] **⚠️ Preflight staff review** — independent reviewer over `k6d` + `pj` ideation → `staff-review.md`; fold Material findings before gates lock
- [ ] **Present ideation gates** — checklist + AC cross-check per member; never self-approve *(captain decides)*
- [ ] **Package** — `dispatch-sprint-execution.md` for the pre-flip pair (`nzb` + `k6d`): boot recipe, per-member build notes, in-drive detached-audit gate (high-stakes CI/release surfaces)

**Drive pre-flip — Commander (separate, cold-booted session)**
- [ ] `nzb` + `k6d`: implementation → validation → done; **detached adversarial audit at validation** (both touch CI/release machinery)
- [ ] Merge each to `next` (PR-merge); concurrency-safe state commits
- [ ] No release cut here — the pre-flip mechanics land on `next` for the flip to ride

**The flip — here (Shaping FO + captain)**
- [ ] Work `pj`'s flip checklist (DoD 5–10): archive → guarded force-push → cut 0.20.0 → marketplace ref flip + paired test → doc reconciliation → migration
- [ ] **⚠️ Pre-cut antipattern audit** on the prepared 0.20.0 line before the tag fires
- [ ] **Cut 0.20.0 on `main`** per `docs/releasing.md` (post-flip mechanics) *(captain authorizes each outward step)*

**Close — Shaping FO**
- [ ] Post-flip verification (some release-machinery issues only manifest when the tag fires) + seed any deferred findings (packaging `44`/`5w`, `xp`)
