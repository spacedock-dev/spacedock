---
session-date: 2026-08-21
sequence: 2
harness: pi
model: glm-5.2-vision-background
model-version-build: glm-5.2-vision-background (lunaroute, thinking max)
first-commit: 84cc857ae
last-commit: 6b8753ad7
duration: ~22h
---

# Session Debrief — 2026-08-21 #2

A conn-held drive to fix 4 pi-live common-journey failures, dispatching their owners onto a single GitHub stack (#748: main ← #725 ← #747 ← #749 ← #750 ← #753). The central discovery: the Pi FO observes worker completion via `subagent_wait` (pi-subagents 0.53.0+), but the worker-lifecycle assert was 10 days stale, crediting only a `subagent status → State: complete` event that never appears. The fix: credit `subagent_wait` completion + bump the CI lane's pi-subagents pin 0.35.1 → 0.53.0. A compaction-alignment gap (the Pi extension's `session_compact` hook re-injects contract text, but PR #738 says re-read state) was filed as a 6th entity.

## Shipped

None merged yet — all 5 stack PRs are open, awaiting review + `CI-E2E-PI` env approval (the `deployments:write` permission is not available to the integration token; env approval is manual, captain-only).

## Filed (backlog)

- **5xww** `pi-subagents-duplicate-extension-load` — issue #746: `spacedock pi` double-loads the pi-subagents extension when the package is also registered, failing with `Tool "subagent" conflicts`.
- **h9nn** `align-pi-compaction-with-force-boot` — the Pi extension's `session_compact` hook re-injects `FO_BOOTSTRAP_TEXT` (a contract pointer), but PR #738 established re-read durable state, not re-inject contract.

## Non-PR commits (workflow-only)

State transitions, gate records, and stage reports across 5 entities in the state checkout — the key ones:
- `84cc857ae` — filed `pi-subagents-duplicate-extension-load` (issue #746 seed)
- `8552e9722` — ideation: `resolve-split-state-gate-artifact-path` (one canonical gate-artifact path)
- `0e3d9ccd7` — ideation: `finish-pi-rejection-flow` (diagnose timeout, propose minimal repair)
- `8abf86de7` — ideation: `repair-pi-default-headless-gate-stop` (async-yield boundary root cause)
- `61734c8ff` — `ntarr`: revert wrong-fix extension; record real root cause (`subagent_wait` vs assert)
- `0cad5157c` — validation: `ntarr` PASSED — assert credits `subagent_wait`; both Pi journeys live PASS
- `da451949e` — dispatch: `align-pi-compaction-with-force-boot` entering ideation

Code commits on the stack branches (not on main): `9dda7ffcb` (747), `903778149` (749), `15eb465e7` (750), `185b53477` (753).

## Decisions

- **`subagent_wait` construct**: `subagent_wait({id, nonBlocking: true})` + return immediately is the preferred completion construct (one-shot wake subscription; Pi wakes the session on completion; captain input resumes the FO's loop; no busy wait, no reinstall). The blocking `subagent_wait({id})` is only the run-to-completion override. `subagent({action:"status", id})` polling is a non-blocking probe, not the completion wait. The Pi adapter doc was updated to prescribe this; the assert credits both the non-blocking `Background task completed` notice (by agent name) and the blocking-wait result (by run id).
- **Compaction alignment**: the Pi extension's `session_compact` hook re-injects contract text (`FO_BOOTSTRAP_TEXT`), but PR #738 (`force-boot-at-compaction-boundary`) established "re-read durable state via `«state.boot»()`, don't re-inject the contract." Filed as `align-pi-compaction-with-force-boot`.
- **PR 725 rebase**: PR 725 was 16 behind `main` (`CONFLICTING`); rebased onto current `main` before stacking. Two conflicts in `internal/release/runtime_live_e2e_test.go` resolved (main simplified `assertNamedLiveEvidence`; 725's parallel-pin check folded in; removed `On`/`WorkflowDispatch` struct field restored minimally).
- **Conn scope**: captain granted conn for pi-related fixes (gates, PR push, CI dispatch/approval) on 2026-08-21.

## Issues — Workflow

- **753 CI red**: `Guard Pi substrate compatibility` hardcoded `expected pi-subagents 0.35.1` while the pin was bumped to 0.53.0 → guard failed → `Install gotestsum` skipped → `gotestsum: command not found` (exit 127). Fixed: updated the guard's hardcoded version to `0.53.0` (commit `185b53477`).
- **Assert 10-day staleness**: `assertWorkerLifecycle`'s Pi completion branch landed 2026-08-10 (`9468d128f`), crediting `subagent status → State: complete`. `subagent_wait` shipped 2026-08-20 (pi-subagents 0.53.0) — 10 days later. The assert never got updated.
- **CI pin lag**: the live runner pinned `PI_SUBAGENTS_VERSION="0.35.1"` (18 versions behind upstream 0.53.0). The 0.35.1 pin couldn't catch the `subagent_wait` staleness. Bumped to 0.53.0 as part of #753.

## Issues — Spacedock

- **`deployments:write` not available**: the `exe-dev-github-integration` GitHub App (used via `GH_HOST=github.int.exe.xyz`) has `actions:write`, `pull_requests:write`, `contents:write` but NOT `deployments:write`. The FO cannot self-approve `CI-E2E-PI` environment deployments. This is a recurring blocker for any pi-lane CI run. Not filed as a GitHub issue (it's an integration-permission config, not a Spacedock bug).
- **`gh stack link` flaky**: the `gh stack` extension's `link` command reports false success and can delete stacks. The REST API (`POST /repos/.../stacks` + `/stacks/{n}/add` with integer arrays) is reliable; `gh stack link` is not. Use the REST API + read back via `stacks/{n}` + GraphQL `stackEntry`.

## Observations

- **Two wrong-fix reworks on `ntarr`**: the first fix (adapter-text binding in `pi-first-officer-runtime.md`) was rejected — live behavior unchanged. The second (Pi-extension tool-intercept blocking gate-presenting bash) was also wrong — the FO already observes completion via `subagent_wait`; blocking its `gate prepare` would be incorrect. The third fix (assert credits `subagent_wait` completion + pin bump) worked. The lesson: diagnose the *actual* completion mechanism in the live transcript before designing the fix; the adapter text and the assert disagreed, and the assert was stale.
- **Worker timeouts at 30min**: glm-5.2-vision-background is too slow for a design-reset implementation that includes live-journey proof runs (~15-40min each). Splitting implementation (no live runs) from validation (FO runs the live journeys directly) worked.
- **SSH agent instability**: the forwarded `clkao` SSH key's socket is unstable (dangling symlink). Probing the live socket fresh per push (`for s in /tmp/ssh-*/agent.*; do SSH_AUTH_SOCK="$s" ssh-add -l 2>/dev/null | grep -q clkao && export SSH_AUTH_SOCK="$s" && break; done`) is the reliable pattern.

## Agent Testimonial

- Date: 2026-08-21
- Harness/runtime: Pi
- Model: glm-5.2-vision-background
- Model version/build: glm-5.2-vision-background (lunaroute, thinking max)
- Session scale: 5 tasks touched; ~12 workers dispatched; 5 PRs touched (725 rebased + 747/749/750/753 created), 0 merged.

The Spacedock workflow structure was load-bearing for this drive: filing the 4 pi-failure owners as distinct entities, driving each through backlog → ideation → implementation → validation with gate reviews, and stacking their PRs onto #748 gave the captain a clear review surface. The conn delegation (FO approves gates, pushes PRs, dispatches CI) let the drive proceed without stopping at every gate. The main friction was the 30-min worker timeout wall — glm-5.2-vision-background is slow, and a design-reset implementation + live-journey proof can't fit in one slice. Splitting implementation from live validation (FO runs the journeys directly) was the workaround, but it meant the FO did validation work that a worker should have owned. The `deployments:write` blocker (can't self-approve `CI-E2E-PI`) is a recurring friction — every pi-lane CI run needs manual captain approval. The `gh stack link` flakiness cost time discovering the REST API alternative. Overall the workflow caught a real 10-day staleness the CI lane couldn't (the 0.35.1 pin masked it), and the live-local-proof discipline (run the journeys before claiming done) caught two wrong fixes the offline tests passed.

## What's Next

- **`CI-E2E-PI` env approvals needed**: 749 (run `324670842`, deployment `6018782736`), 750 (run `32501668107`), 753 (re-run on `185b53477`). Captain-only.
- **`align-pi-compaction-with-force-boot`** (`h9nn`): ideation dispatched (run on 753's tip); needs its stacked PR on 753 once implementation lands.
- **Stack #748 merge order**: 725 → 747 → 749 → 750 → 753 → (compaction-align PR). Each merge re-points the layer above.
- **Live 4-journey sweep on the stack tip**: run all 4 journeys (GateGuardrail, RejectionFlow, DefaultHeadlessGateStop, AutoContinueAfterImplementation) on the 753 tip as a single consolidated green signal before merge. Not yet done (per-layer greens only).
