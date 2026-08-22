# Debrief — 2026-08-21 pi-failure drive + subagent_wait construct

## Scope held
Captain granted conn for **pi-related fixes** (gates, PR push, CI dispatch/approval) on 2026-08-21. Drive: file and dispatch owners for the 4 pi-live common-journey failures, run local target live pi with `lunaroute/glm-5.2-vision-background:max`, open PRs to the stack.

## The 4 pi-failure owners → stack #748 (main ← 725 ← 747 ← 749 ← 750 ← 753)

| Journey | Owner entity | Fix | PR | Local live green |
|---|---|---|---|---|
| `pi-subagents` duplicate load (issue #746) | `pi-subagents-duplicate-extension-load` (`5xww`) | gate explicit `--extension`/`--skill` on package NOT registered | #747 | front-door smoke PASS (CI) |
| `GateGuardrail` | `resolve-split-state-gate-artifact-path` (`mvmpzg`) | one canonical resolver-owned gate-artifact path (split-root) | #749 | PASS (911s, local) |
| `RejectionFlow` (registered XFAIL, owner p17) | `finish-pi-rejection-flow` (`p17`) | async-yield-to-gate-lifecycle binding + Pi-dialect extractors | #750 | PASS (1500s, no timeout) |
| `DefaultHeadlessGateStop` + `AutoContinueAfterImplementation` | `repair-pi-default-headless-gate-stop` (`ntarr`) | assert credits `subagent_wait` completion; pin bump 0.35.1→0.53.0; adapter prescribes `nonBlocking:true` | #753 | PASS (759s / 1966s) |

## The `subagent_wait` construct decision (load-bearing)
- **Preferred: `subagent_wait({id, nonBlocking: true})` + return immediately** — one-shot wake subscription; Pi wakes the session on completion; captain input resumes the FO's loop meanwhile; no busy wait, no reinstall. The Pi analog of Codex's async final-status mailbox notification. The completion surfaces as a `Background task completed: **<agent>**` notice in the parent transcript.
- **Blocking `subagent_wait({id})` is the run-to-completion override only** (FO must report results back in this turn) — returns `Waited … for run "<id>"; done. Outcome: 1 complete.`
- **`subagent({action:"status", id})` polling** is a non-blocking probe, NOT the completion wait (the old adapter text said to poll; that was stale — the FO uses `subagent_wait`, not status polls).
- `assertWorkerLifecycle` now credits BOTH the non-blocking `Background task completed` notice (by spawned agent name) AND the blocking `subagent_wait` result (by run id) AND the old `subagent status → State: complete` path. No assert loosening (`spawns==1, completed>=0, completed<validation, "- DONE:"` unchanged).
- **The assert was 10 days stale**: the `State: complete` branch landed 2026-08-10; `subagent_wait` shipped 2026-08-20 (pi-subagents 0.53.0). The CI lane pinned 0.35.1, so it couldn't catch this — bump to 0.53.0 is part of #753.

## The 753 CI red (root cause + fix)
- **Red**: `Guard Pi substrate compatibility` failed — it hardcoded `expected pi-subagents 0.35.1` while the pin env was bumped to 0.53.0. Guard failed → `Install gotestsum` skipped → downstream `gotestsum: command not found` (exit 127).
- **Fix**: update the guard's hardcoded version to `0.53.0` (commit `185b53477`, pushed). CI re-running.

## Compaction-alignment finding (filed, dispatched)
- PR #738 (`force-boot-at-compaction-boundary`, merged) established: at compaction, re-**read durable state** via one `«state.boot»()`; the **contract does not need re-injecting**.
- The Pi extension's `session_compact` hook re-injects `FO_BOOTSTRAP_TEXT` (a contract pointer) — the thing #738 rejected. **Misaligned.**
- Filed `align-pi-compaction-with-force-boot` (`h9nn`), backlog→ideation consumed under conn, dispatched to stack on 753's tip. NOT yet implemented.

## Credential situation (environment facts)
- `git push` works via the **forwarded clkao SSH key** (probe `for s in /tmp/ssh-*/agent.*; do SSH_AUTH_SOCK="$s" ssh-add -l 2>/dev/null | grep -q clkao && export SSH_AUTH_SOCK="$s" && break; done` per push — the socket is unstable/dangling).
- `gh` as `nomen429` @ github.com is READ-only on spacedock-dev/spacedock.
- `GH_HOST=github.int.exe.xyz` is a proxy authenticating as **clkao (admin)** — repo-scoped GET + most POST writes work (issue comments, PR create, stack create/add, workflow dispatch). It CANNOT do `deployments:write` (env approval 403) — `CI-E2E-PI` approval is **manual, captain-only**.
- Stack API: create via `POST /repos/.../stacks` `{pull_requests:[ints]}`; append via `POST /repos/.../stacks/{n}/add` `{pull_requests:[ints]}` (the gh extension's `gh stack link` is flaky/reports false success — use the REST API + read back via `stacks/{n}` + GraphQL `stackEntry`).

## Open / pending captain
- **`CI-E2E-PI` env approvals** needed for: 749 (run `32467087042`, deployment `6018782736`), 750 (run `32501668107`), 753 (re-run on `185b53477`). I can't self-approve (deployments:write 403).
- `align-pi-compaction-with-force-boot` ideation running; will need its own stacked PR on 753.
