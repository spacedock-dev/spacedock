# Dispatching Spacedock ensigns through Paseo

A Paseo-launched Codex ensign completed one stamped implementation stage. It committed `receipt.txt` on the task branch at `a5a82f4` and its report in the split-root state checkout at `e6524c9`. The report contained the word delivered mid-run. The first officer (FO) verified both commits, then closed the worker. This was one same-daemon run with no Git remote. [Fixture](run-log.md#fixture-and-dispatch) · [Steering](run-log.md#steering-and-roster) · [Commits](run-log.md#completion-and-durable-proof) · [Closure](run-log.md#shutdown-and-cleanup)

## Dispatch capabilities

| Capability | Support on the tested Paseo installation | Evidence |
|---|---|---|
| `«worker.spawn»` | **Supported.** MCP `create_agent` accepted an explicit provider/model, workspace, title and assignment pointer. It returned an agent ID and the expected code directory. | [Launch request and receipt](run-log.md#workspace-and-spawn) |
| `«addressable-worker»` | **Conditional.** FO steering reached the ordinary worker. It reported no agent-scoped send tool. A second worker with `paseo-safehouse.paseo-access=daemon` sent an addressed message that arrived before completion. A parent-only grant remains unproven. | [FO steer](run-log.md#steering-and-roster) · [Child send, parent receipt and commit](run-log.md#child-return-channel) |
| `«async-dispatch»` | **Supported.** Spawn returned `running`. The FO sent a steer, listed the roster and used a bounded wait before completion. | [Spawn](run-log.md#workspace-and-spawn) · [Concurrent operations](run-log.md#steering-and-roster) |
| `«completion-signal»` | **Supported.** With `notifyOnFinish=true`, the FO received an unsolicited `finished` event with server and agent IDs. MCP status and CLI wait independently reported `idle`. The committed report satisfied the stage gate. | [Event, status and report](run-log.md#completion-and-durable-proof) |
| `«worker.shutdown»` | **Supported.** After file verification, `kill_agent` succeeded. Both status and roster showed `closed` before workspace archival. | [Shutdown and cleanup](run-log.md#shutdown-and-cleanup) |
| `«context-budget»` | **Supported through MCP for this provider.** The finished snapshot reported 30,093 of 258,400 context tokens used (11.6%). CLI `inspect.LastUsage` exposed token counts without the context-window denominator. A bound adapter can apply its reuse threshold and dispatch fresh when the reading is unavailable. | [Completed snapshot and CLI projection](run-log.md#completion-and-durable-proof) · [Reuse rule](run-log.md#retry-and-budget-contracts) |
| `«roster-reconcile»` | **Supported as a roster read.** MCP `list_agents` found the expected worker by cwd, with parent IDs and assignment labels. CLI lookup by epoch found that ID running and later closed. The FO must compare these rows with its dispatch records to classify drift. | [Running roster](run-log.md#steering-and-roster) · [Closed roster](run-log.md#shutdown-and-cleanup) |

## Worktree ownership

Let Spacedock own branch selection and worktree creation. With the fixture on local `main`, `dispatch build --stamp` created the task branch at `f98e81a`. The deliberately older `origin/main` ref pointed to `ff2b560`. A second checkout of the task branch failed with Git exit 128. A separately named branch could avoid that error while giving the ensign a different checkout from the one Spacedock verifies. That last case is a design risk inferred from the stamped paths. [Refs, stamp, assignment and Git refusal](run-log.md#fixture-and-dispatch)

Register the existing physical worktree with `create_workspace(isolation="local", path=...)`. Paseo returned `isolation=worktree` for that request and preserved the exact cwd. Match on physical path and checked-out branch. The worker's `pwd -P`, receipt and report confirmed the intended checkout. [Workspace adoption](run-log.md#workspace-and-spawn) · [Worker and Git evidence](run-log.md#completion-and-durable-proof)

## Cross-runtime dispatch

The helper can already render a Codex assignment independently of the FO's detected host. A subprocess reporting `Runtime: claude (CLAUDECODE)` emitted a byte-identical envelope with `--host codex`. Paseo separately accepted `safehouse-codex/gpt-6.2-astra`. Together these establish the two interfaces needed for Claude-to-Codex dispatch. A live Claude FO and a remote daemon remain untested. [Cross-host comparison](run-log.md#cross-host-envelope) · [Provider selection](run-log.md#workspace-and-spawn)

The FO supplied the target provider/model and workspace ID. The follow-up recovered the physical worker directory from the stamped entity through `status --read --json` and Git. These inputs suffice for the proposed adapter to use the existing Codex envelope unchanged. [Original envelope](run-log.md#fixture-and-dispatch) · [Workspace resolution and launch](run-log.md#child-return-channel)

## Access and write scope

The FO ran under `safehouse-codex` with `role=worker`, `repo=spacedock-v1`, `mandate=paseo-backend` and `paseo-safehouse.paseo-access=daemon`. It adopted, launched, steered, inspected and closed its worker. The follow-up verified child-to-parent delivery under the same daemon access label. [FO access](run-log.md#runtime-and-access) · [Launch](run-log.md#workspace-and-spawn) · [Steer](run-log.md#steering-and-roster) · [Closure](run-log.md#shutdown-and-cleanup) · [Child access](run-log.md#child-return-channel)

The assignment specified the code branch and the canonical state path. The worker committed each deliverable to its assigned checkout. Both launches requested `modeId=auto`, and Paseo returned `currentModeId=full-access`. The first worker’s finished snapshot still reported `full-access`. The runs establish the mode-reporting discrepancy. Its cause is undetermined. [Assignment](run-log.md#fixture-and-dispatch) · [First request](run-log.md#workspace-and-spawn) · [Finished snapshot and commits](run-log.md#completion-and-durable-proof) · [Second request](run-log.md#child-return-channel)

## Completion and recovery

Use a finish event to trigger inspection. Match the server, worker and assignment to a new exact-stage report in the canonical entity. Verify the checklist, path-scoped commit and clean checkout. On reuse, capture the pre-dispatch report revision so an old report cannot satisfy a later callback. The fixture proved the fresh-worker path. The callback contained the agent and message IDs but no Spacedock assignment epoch. [Event and verification](run-log.md#completion-and-durable-proof) · [Epoch label](run-log.md#steering-and-roster)

A one-second CLI wait exited **0** with JSON `status=timeout` while the worker was making progress. Continue waiting on a live assignment. For error or permission notifications, inspect lifecycle state, pending permissions, recent activity and the entity before acting. The spawn receipt advertised both notification types. Only the finish notification was exercised. [Timeout and activity](run-log.md#steering-and-roster) · [Callback contract](run-log.md#workspace-and-spawn) · [Status fields](run-log.md#completion-and-durable-proof)

For a confirmed spawn failure or ended worker missing both signal and report, carry the existing one-retry rule over to Paseo. Reconcile the roster before replacing an ambiguous create. Record the retry on the entity, then nudge an addressable worker or create a fresh `-retry` worker if the old session is gone. On the second failure, hold that entity. A finish event with an incomplete report needs stage review. A wait timeout alone spends no retry. These are proposed bindings for failure cases. [Existing retry rule](run-log.md#retry-and-budget-contracts) · [Live steer and wait](run-log.md#steering-and-roster) · [Closed-session observation](run-log.md#shutdown-and-cleanup)
