# Staff report for the Paseo dispatch backend

## What needs to change?

**Reuse `spacedock dispatch build` unchanged.** A Claude-marked subprocess using `--host codex` emitted the same schema-v2 envelope as the Codex FO. The field missing at launch was the Paseo workspace ID. Only Paseo could return that ID after adopting the local checkout. The follow-up resolved the checkout through the existing structured status read. Neither probe needed to parse the assignment or add a binary field. [Cross-host comparison](run-log.md#cross-host-envelope) · [Envelope and stamp](run-log.md#fixture-and-dispatch) · [Worktree lookup](run-log.md#child-return-channel) · [Workspace receipt](run-log.md#workspace-and-spawn)

The following changes are **designs for review**. The run exercised their underlying operations. The selectable binding, real Claude acceptance test, failure fixtures and scoped child permission remain unproved.

| Proposed change | Why existing behavior alone is insufficient | What to reuse |
|---|---|---|
| A selectable FO dispatch binding | The build returned an assignment. The FO still had to turn it into a Paseo workspace and a running worker. | The Codex pointer, canonical name, description, Git checkout and Paseo MCP tools. [Build](run-log.md#fixture-and-dispatch) · [Launch](run-log.md#workspace-and-spawn) |
| A live Claude acceptance test | The host-marker comparison exercised the renderer. The actual FO that received a callback ran Codex. | The split-root fixture, steering word and two-commit assertions. [Host comparison](run-log.md#cross-host-envelope) · [FO identity](run-log.md#runtime-and-access) · [Callback and commits](run-log.md#completion-and-durable-proof) |
| A parent-only daemon permission for ordinary ensigns that need to ask questions | The ordinary worker lacked the send tool. A daemon-labelled child delivered a message, but that label also equipped the probe FO to manage workers. A finish notification cannot carry a question while the worker is still running. | Paseo's existing send operation and authenticated parent relationship. The outstanding dependency is the **Paseo daemon ACL**. [Ordinary worker](run-log.md#steering-and-roster) · [FO grant](run-log.md#runtime-and-access) · [Labelled child](run-log.md#child-return-channel) · [Management operations](run-log.md#shutdown-and-cleanup) |

## From build output to create_agent

**Proposed handoff.** At dispatch, the FO selects the requested backend and calls its normal build with `--host codex`, the current workflow, canonical entity, stage, checklist and applicable stamp flags. On success, the following calls implement `«worker.spawn»`. A failed build follows its existing error path before any Paseo call. [Recorded invocation and output](run-log.md#fixture-and-dispatch)

1. Read `frontmatter.worktree` from `spacedock status --workflow-dir "$workflow" --read "$entity" --json`. Resolve a relative stamp against the workflow's Git root and keep an absolute stamp as supplied. For a worktree-free stage, use the workflow's code checkout. Bind the code directory to the existing branch. The stamped-path case was verified. The worktree-free case needs a fixture. [Structured lookup](run-log.md#child-return-channel) · [Branch and state paths](run-log.md#fixture-and-dispatch)
2. Adopt that directory with `create_workspace({projectId, isolation:"local", path:codeCwd, title:build.description})`, or use an already identified workspace for the same project and physical path. Check the returned cwd and take its `workspaceId`. Paseo classified the adopted checkout as a worktree. Asking Git for a second checkout of the occupied branch failed. [Adoption receipt](run-log.md#workspace-and-spawn) · [Git refusal](run-log.md#fixture-and-dispatch)
3. Send the exact object below for the first fixture with no supplemental notes. This is the **proposed production mapping**, instantiated with the recorded workspace and worker values. The actual probe additionally supplied two tracking labels and appended experiment instructions to the pointer. [Build input fields](run-log.md#fixture-and-dispatch) · [Actual request](run-log.md#workspace-and-spawn)

```json
{
  "provider": "safehouse-codex/gpt-6.2-astra",
  "workspaceId": "wks_8aec77dc0eda8e94",
  "title": "Paseo receipt: implementation",
  "labels": {
    "role": "ensign",
    "spacedock.worker": "paseo-qjyrql-implementation",
    "spacedock.epoch": "qjyrql-1"
  },
  "settings": { "modeId": "auto" },
  "notifyOnFinish": true,
  "initialPrompt": "$spacedock:ensign; then Read /tmp/spacedock-dispatch/paseo-qjyrql-implementation.md and treat its content as your assignment."
}
```

Map `description → title`, `name → labels["spacedock.worker"]` and `prompt → initialPrompt` verbatim. Optional scope notes follow the intact pointer. The FO supplies its completion epoch, requests the session's selected Paseo mode and explicitly names a Codex provider/model available on this daemon. For the acceptance fixture, use the pair shown above. The build returned `model:null`; inheriting a Claude FO model would leave the Codex provider unresolved. Record the daemon's returned agent ID, cwd, model and mode against the entity, stage and epoch. Paseo added the parent labels itself. [Model catalog](run-log.md#runtime-and-access) · [Envelope](run-log.md#fixture-and-dispatch) · [Actual request and returned mode](run-log.md#workspace-and-spawn) · [Roster identity](run-log.md#steering-and-roster)

## Contract ownership

| Component | Proposed implementation |
|---|---|
| Spacedock binary | **No change.** Build owns the stamp, state sync, assignment file and Codex bootstrap. `status --read --json` already exposes the stamp. The worker reads its canonical state path and write scope from the assignment. [Stamp and assignment](run-log.md#fixture-and-dispatch) · [Structured lookup](run-log.md#child-return-channel) |
| FO contract | At the first dispatch load point, bind the worker transport independently of the FO host. For v1, an explicit session instruction selects Paseo, fresh Codex workers and the provider/model. Keep the shared scheduling, checklist, report gate and retry ledger. No new configuration format is required for this lane. [Cross-host rendering](run-log.md#cross-host-envelope) · [Existing gate and retry rules](run-log.md#retry-and-budget-contracts) |
| FO runtime adapters | Keep the host adapter for boot and host-local operations. Add `skills/first-officer/references/paseo-fo-dispatch.md` beside `claude-fo-dispatch.md`. Selecting Paseo loads its worker bindings for spawn, steer, callback, roster, context reading and close. The Claude FO must have agent-scoped Paseo tools on the same daemon. Prove that in the live acceptance test. [Daemon access](run-log.md#runtime-and-access) · [Steer and roster](run-log.md#steering-and-roster) · [Callback and context](run-log.md#completion-and-durable-proof) · [Close](run-log.md#shutdown-and-cleanup) |
| Codex ensign adapter | Keep the existing skill bootstrap, prompt-defined write scope and path-scoped commits. The fresh worker can finish normally with no daemon label. Add a parent-directed clarification binding only when its scoped permission is available. [Assignment](run-log.md#fixture-and-dispatch) · [Completed ordinary worker](run-log.md#completion-and-durable-proof) · [Child messaging](run-log.md#child-return-channel) |

**Proposed recovery.** Bind a finish event to its server and worker IDs. Verify a new exact-stage report against the pre-launch entity revision, retained checklist and clean path-scoped commits. On a lost create receipt, reconcile parent and attempt labels before considering another spawn. A bounded wait returned exit 0 with `status=timeout` during real progress. An observed dispatch failure consults the ledger: nudge a live worker or launch one fresh `-retry` replacement if the old session ended. A second failure holds the entity. Close the owned worker after verification or explicit supersede. Duplicate events and failure/retry cases need fixtures. [Timeout and roster](run-log.md#steering-and-roster) · [Completion proof](run-log.md#completion-and-durable-proof) · [Retry rule](run-log.md#retry-and-budget-contracts) · [Verified closure](run-log.md#shutdown-and-cleanup)

## Child daemon access

**Yes.** The second ensign received `paseo-safehouse.paseo-access=daemon` at creation. Its MCP `send_agent_prompt` receipt and the FO's inbound message shared ID `059cc77b-0f8b-4769-9c93-600311597bc2`. The FO received the message before the finish callback. The unlabelled ensign had reported that the send tool was absent. This proves the transport in both directions under the tested grants. [Original steer](run-log.md#steering-and-roster) · [Child label, inbound message and committed receipt](run-log.md#child-return-channel)

**Required ACL behavior, unproved.** For a child that needs this capability, expose the send tool and authorize delivery only as that child to its daemon-recorded parent on the same server. Derive identity and destination from the authenticated relationship. A caller-supplied ID or editable label must not confer authority. The grant must not permit messaging siblings or other agents, enumerating or reading their activity, creating or controlling agents, managing workspaces, approving permissions, or delegating the grant. Wait for the Paseo daemon ACL before enabling this route by default. Its configuration spelling is unverified. The proposed acceptance fixture checks delivery to the parent and denial for an unrelated fixture recipient and a management operation. [Parent metadata and limits of the test](run-log.md#child-return-channel)

**Paseo finding.** Both `create_agent` calls requested `modeId=auto` and returned `currentModeId=full-access`. The first completed worker still reported `full-access`. Record requested and returned values in the acceptance evidence. The cause is undetermined. [First launch](run-log.md#workspace-and-spawn) · [Finished status](run-log.md#completion-and-durable-proof) · [Second launch](run-log.md#child-return-channel)
