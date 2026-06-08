# Sprint 0198 — pre-flip hardening

**Goal:** ship the pre-flip essentials as **0.19.8** before the 0.20.0 flip — the
binary/version/install UX (`qa`), the Codex plugin auto-install (`z9`), and the
test-hygiene fix (`kb`). Survey polish and the e2e release-gate are deferred to post-flip.

**Deliverable:** spacedock **0.19.8** cut on `next`.

## Members

Membership is the query, not this table:

```bash
spacedock status --workflow-dir docs/dev --where sprint=0198-pre-flip-hardening --where 'sprint-readiness != defer'
```

| Entity | group | gate | what it delivers |
|--------|-------|------|------------------|
| `qa` spacedock-binary-missing-install-journey | binary-ux | **ideation ✓ approved** | install/upgrade journey + versions-not-"contract" messages (missing / incompatible / stale) |
| `z9` codex-plugin-auto-install | binary-ux | **ideation ✓ approved** | front-door Codex plugin auto-install (mirror #311), channel-tracked via the shared `devBranch` |
| `kb` migration-check-prune-state-walk | test-hygiene | fast-track (no gate) | prune `.spacedock-state` in the migration-check walk + drop orphaned survey fixtures |

**Deferred (out of this sprint):**
- `nzb` gate-release-on-e2e → **post-flip** (useful not critical; per-PR e2e already verifies; when built, the gate must be on-branch before tagging).
- `vh` survey-skill-correctness-pass → **post-flip** (consolidates 69/1p/4t + the agentsview git-root-model fix; design banked).
- `jh` version-skew (no contract bump pre-first-release — self-resolves at the first release); `5h0` (blocked on #315).

## Definition of Done

1. Every active member (`qa`, `z9`, `kb`) `done` / PASSED + merged to `next`.
2. `go test ./...` from the repo root green with `.spacedock-state` + `.claude/worktrees` present (owned by `kb`).
3. `qa`'s behavior proven by a **captain-run live drive** at sprint acceptance (the version-bearing, jargon-free messages observed).
4. `z9`'s front-door auto-install proven by a **detached adversarial audit** (high-stakes surface) + its host-native install test.
5. spacedock **0.19.8** stamped + cut on `next` (captain-gated release).

## Channel-tracking note (the qa/z9 ↔ flip dependency)

`z9` (and the existing `#311` Claude auto-install) install the plugin from the shared
`devBranch` (`frontdoor.go:49`, today `"next"`). **The flip must retarget `devBranch`
`next → main`** (the per-channel `devBranch` stamp / the 2-channel release-mechanics work)
so the released stable binary auto-installs the `main` plugin. That retarget is a **flip
requirement**, NOT a 0198 task — `z9` is correct as long as it uses `devBranch` (confirmed in
its design).

## Out of scope

The 0.20.0 flip (`pj`). The 2-channel brew + `devBranch` retarget (flip-release-mechanics).
Survey polish (`vh`, post-flip).

## Status

**Fully shaped + approved — ready for the Commander handoff.** `qa` + `z9` ideation gates
captain-approved; `kb` fast-track. The package
([`dispatch-sprint-execution.md`](dispatch-sprint-execution.md)) is final; a separate
Commander session drives the three implementations (`qa` → `z9`; `kb` independent).
