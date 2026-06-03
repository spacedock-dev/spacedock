---
title: Pi runtime support — adapt Spacedock to pi-native teams/subagents
status: implementation
score: "0.36"
source: captain (2026-06-03) — dogfood pi support from the prior PR #155 and evaluate pi-agent-teams / pi-subagents as usable ensign constructs
issue: spacedock-dev/spacedock#155
id: s9kcdyb9r5t8addppnnce54j
started: 2026-06-03T21:07:03Z
worktree: .worktrees/spacedock-ensign-pi-runtime-support
---
# Pi runtime support — adapt Spacedock to pi-native teams/subagents

Bring forward the useful parts of the old PR #155 Pi runtime compatibility baseline into Spacedock v1, but design it for the current Go launcher, split-root state model, and current Pi ecosystem.

## Problem

Spacedock v1 currently has stable launcher and skill surfaces for Claude/Codex-oriented dispatch, but Pi does not expose the same runtime tools. A Pi first officer cannot safely assume Claude `Agent`, `SendMessage`, `TeamCreate`, or `TeamDelete` tool signatures exist, and doing so would make Pi support fail at the first dispatch.

The runtime gap is specifically about ensign lifecycle: first-officer code and skill instructions need a Pi-native way to build an assignment, dispatch an ensign, wait for completion, route follow-up when validation rejects work, and shut down or mark the runtime complete. This must fit the current Go launcher, split-root state checkout model, and compatibility-first bootstrap scope without adding PR/mod behavior.

## Context inspected

- Old PR #155 clone at `/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/tmp.SKgFLFccFF/spacedock-pr155/` was present. Its useful Pi baseline is the `scripts/pi_session_registry.py` and `scripts/pi_worker_runtime.py` model: a thin worker-label to Pi session mapping, completion epochs to reject stale evidence, and explicit dispatch/reuse/shutdown commands. Its limitation for v1 is that it is Python test-harness-era code, not a Go launcher/package surface, and it predates the current split-root v1 workflow shape.
- `pi-agent-teams` clone at `/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/tmp.zCMAQuSWp2/pi-agent-teams/` was present. It exposes a Pi `teams` tool, not Claude-compatible tools. Relevant actions include `member_spawn`, `delegate`, `message_dm`, `member_shutdown`, and `team_done`, with task/mailbox/session state under a teams root. This is closest to a long-lived team substrate, but it requires an adapter that speaks its action schema.
- Installed/cached `pi-subagents` artifacts show a parent-side `subagent` tool. It supports a single delegated task, parallel `tasks`, optional worktree isolation, management actions such as `list/status/interrupt`, and acceptance-style completion from a child invocation. It is not a long-lived mailbox team by default, which makes it a better first-slice ensign dispatch primitive than a direct Claude-team emulation.
- This ideation task itself attempted a bounded dogfood run using the Pi subagent construct as the Spacedock ensign. The first dispatch attempt reached the parent-side `subagent(...)` tool but failed before worker start with provider auth (`No API key found for openrouter`). Treat that as a live-harness setup finding, not a product-runtime impossibility: the live Pi runner should isolate Pi's home/session directories and copy the operator's existing Pi OAuth credentials, the same way the Codex live runner uses an isolated `CODEX_HOME` with copied auth.

## Proposed approach

### V1 target shape

Target a small Pi runtime adapter abstraction with `pi-subagents` as the default first implementation and PR #155's session-registry ideas as compatibility scaffolding for follow-up reuse. Treat `pi-agent-teams` as an optional second substrate behind the same adapter, not as the first slice.

The adapter contract should be host-neutral at the Spacedock layer:

1. `BuildDispatch`: reuse the existing `spacedock dispatch build` assignment JSON/prompt shape where possible, adding a Pi host mode only for fields that differ from Claude/Codex.
2. `DispatchEnsign`: launch a Pi ensign using the parent-side `subagent(...)` tool when available. The child receives the same ensign skill contract and dispatch prompt used by other hosts, with Pi-specific instructions documented in `skills/ensign/references/pi-ensign-runtime.md`.
3. `AwaitCompletion`: treat the subagent tool result/acceptance report as completion evidence. For asynchronous or reused-session flows, persist minimal worker state inspired by PR #155: logical worker label, runtime substrate, session/run id if exposed, entity/stage, state, and completion epoch.
4. `RouteFollowup`: for the first slice, prefer a fresh `subagent(...)` dispatch with feedback embedded in the prompt unless a Pi substrate exposes a safe resumable handle. If session reuse is enabled, use the PR #155 epoch guard so old completion text cannot satisfy a new cycle.
5. `Shutdown`: for `pi-subagents`, completion is process-scoped and no long-lived mailbox shutdown is required; mark the worker complete/closed. For `pi-agent-teams`, call `teams({action:"member_shutdown", ...})` or `teams({action:"team_done", all:true})` through an adapter when that substrate is explicitly selected.

### Substrate comparison and decision

- **PR #155 session registry helper:** keep the model, not the Python code. It solves the durable identity problem Spacedock needs for follow-up routing and stale-completion rejection. Reimplement as a small Go package/fixture-driven contract if v1 needs session reuse. Do not make it the primary user-facing runtime because Pi now exposes higher-level tools.
- **`pi-subagents`: choose for first slice.** It already matches Spacedock's ensign mental model: the parent dispatches a bounded child task and receives a completion result. It avoids assuming a long-lived team mailbox and is suitable for live dogfood proof. The tradeoff is that follow-up may be fresh-dispatch first rather than same-session routing until the exposed run/session handles are proven.
- **`pi-agent-teams`: defer as optional team substrate.** It is promising for long-lived workers, DM follow-up, member status, and explicit shutdown, but its API is `teams` action JSON rather than Claude-compatible team tools. First support should document an adapter mapping instead of baking its actions into FO prose.

### Implementation slice

Keep the first implementation compatibility-first and small:

- Add Pi runtime documentation under the skill runtime references and first-officer runtime docs explaining that Pi uses either `subagent(...)` or `teams({...})` adapters, never Claude tool names.
- Add a Go dispatch/front-door hook such as `spacedock pi` or a Pi host mode for `spacedock dispatch build` that emits Pi-specific dispatch instructions and a machine-readable substrate selection (`subagents` default; `teams` optional).
- Add a small internal package for Pi runtime contract parsing/validation: accepted dispatch spec, completion evidence shape, optional worker registry record, and stale completion epoch behavior.
- Add offline tests with fixtures for `subagent` dispatch JSON, `teams` action mapping, registry stale-completion rejection, and split-root entity paths.
- Add one live/dogfood runtime scenario that launches a Pi parent with the subagents extension, dispatches a tiny ensign against a temp split-root workflow, and verifies committed state/output outside the task prose.

### Live harness isolation

Mirror the Codex live runner's isolation pattern instead of requiring a globally-mutated Pi home:

1. Create a temp `PI_CODING_AGENT_DIR` for config, installed packages, settings, and auth.
2. Create a temp `PI_CODING_AGENT_SESSION_DIR` or pass `--session-dir` for session JSONL output.
3. Copy the operator's existing `~/.pi/agent/auth.json` into the temp config dir when present, so subscription/OAuth logins are reused without exposing the real global session/package state to the test.
4. Install or point at only the local test resources needed for the run: Spacedock skills, `pi-subagents` or `pi-agent-teams`, and any local extension path under test.
5. Run the Pi parent/worker from the temp workflow root and assert durable state checkout changes, exit codes, and logs. Do not assert on global `~/.pi/agent` state and do not require the operator to relogin for every live run.

This is the answer to the auth failure observed during dogfood dispatch: isolate the runtime home, but reuse the already-authorized OAuth token file.

## Out of scope

- PR merge, mod behavior, or terminalization behavior specific to Pi.
- Rewriting Claude or Codex runtime behavior.
- Assuming Claude `Agent`, `SendMessage`, `TeamCreate`, or `TeamDelete` signatures in Pi.
- Making `pi-agent-teams` the required first implementation or depending on its long-lived mailbox for the first slice.
- Large runtime orchestration rewrites beyond dispatch shape, package/front-door hooks, skill docs, and tests.

## Riskiest unknown first

The riskiest unknown is whether a Pi parent can dispatch a Spacedock ensign with the available `subagent(...)` surface and get back structured completion evidence strong enough for the FO to advance workflow state. The implementation should start with a live dogfood smoke before building deeper registry semantics: create a temp split-root workflow, dispatch one Pi subagent as `spacedock:ensign` or equivalent configured worker, have it mutate only the state checkout entity body and commit path-scoped, then assert the resulting git log and file state.

No spike is needed for the existence of `pi-agent-teams` action names because the local clone exposes the `teams` tool action schema directly. A spike is still needed before relying on same-session follow-up through either Pi substrate because the currently relevant `pi-subagents` surface is parent-side bounded dispatch, not a long-lived mailbox by default.

## Acceptance criteria

**AC-1 - Pi dispatch uses a Pi-native tool contract instead of Claude-compatible tool names.**
Verified by: focused Go fixture tests for the Pi dispatch builder that fail if generated Pi instructions/specs contain `Agent`, `SendMessage`, `TeamCreate`, or `TeamDelete`, and pass only when the default dispatch target is `subagent` with Pi-specific parameters.

**AC-2 - The default Pi ensign path can complete a real split-root workflow task through `pi-subagents` from an isolated Pi home that reuses existing OAuth credentials.**
Verified by: a live test tag scenario, for example `go test -tags live -run TestLivePiSubagentEnsignSmoke ./internal/ensigncycle -v`, that creates temp `PI_CODING_AGENT_DIR` and session dir, copies `~/.pi/agent/auth.json` into the temp config dir when available, launches Pi with the subagents extension, dispatches a tiny ensign, and asserts the state checkout entity body/git log changed as expected.

**AC-3 - Optional `pi-agent-teams` support is represented as an adapter over the `teams` action schema, not as Claude team emulation.**
Verified by: offline contract tests that feed a dispatch/follow-up/shutdown intent and assert the adapter emits `teams` actions such as `member_spawn` or `delegate`, `message_dm`, `member_shutdown`, and `team_done` with required fields.

**AC-4 - Follow-up routing cannot accept stale completion evidence.**
Verified by: unit tests for the Pi worker registry/completion parser that create two epochs for one logical worker and assert an old completion is rejected after a follow-up dispatch/reuse cycle.

**AC-5 - Split-root state paths are preserved in Pi dispatch prompts and tests.**
Verified by: fixture or CLI tests that build a Pi dispatch for `docs/dev` with `state: .spacedock-state` and assert the assigned entity path points at the state checkout while code worktree handling remains separate.

## Test plan

- **Offline dispatch contract tests (low cost):** Add table-driven Go tests for Pi dispatch spec generation. Include negative cases for banned Claude tool names and positive cases for `subagent` default fields and split-root entity paths.
- **Offline teams-adapter tests (low/medium cost):** Add fixtures that translate Spacedock lifecycle intents into `teams` tool action payloads. Assert exact JSON actions and required fields for spawn/delegate, follow-up DM, status/shutdown, and team done.
- **Offline registry tests (low cost):** Reimplement the PR #155 registry semantics in Go only if needed for follow-up/reuse. Test atomic record persistence, active/completed/shutdown transitions, and stale epoch rejection.
- **Live Pi subagent smoke (medium/high cost, required for the runtime claim):** Create an isolated Pi home (`PI_CODING_AGENT_DIR`) and session dir, copy `~/.pi/agent/auth.json` into the temp home when available, load only the local Spacedock skills plus `pi-subagents`, and run a temp split-root workflow where the FO dispatches one ensign through `subagent(...)`. The assertion should inspect output, exit code, state checkout git log, and entity file state; it should not pass by grepping skill prose.
- **Live teams smoke (optional/deferred):** If the optional `pi-agent-teams` adapter is implemented in the same slice, run a small leader session that uses `teams({action:"member_spawn"})`, sends one `message_dm`, and then `member_shutdown`/`team_done`. Otherwise leave this as a documented future substrate test.
- **Baseline Go gate:** Run `go test ./...` for every implementation change; run `go test ./... -race` when registry concurrency or filesystem mutation code is added.

## Stage Report: ideation

- DONE: Decide and document the v1 Pi runtime target shape, including how PR #155's session-registry helper and pi-agent-teams/pi-subagents compare.
  Documented `pi-subagents` as the first-slice default, PR #155 registry semantics as reusable scaffolding, `pi-agent-teams` as an optional `teams` action adapter, and the isolated-Pi-home/OAuth-copy live-harness requirement discovered during dogfood dispatch.
- DONE: Add entity-level acceptance criteria and a test plan with at least one live/dogfood runtime proof and focused offline contract tests.
  Added AC-1 through AC-5 with external verification commands/tests plus a test plan covering live Pi subagent smoke and offline dispatch/teams/registry/split-root tests.
- DONE: Append a `## Stage Report: ideation` section accounting for every checklist item.
  This section accounts for all three assignment checklist items using DONE/SKIPPED/FAILED status labels.

### Summary

The ideation result chooses a compatibility-first Pi runtime shape: default to bounded `pi-subagents` ensign dispatch, preserve PR #155's worker/session registry ideas for safe follow-up semantics, and treat `pi-agent-teams` as an adapter over its native `teams` action schema. The future implementation is scoped to skill/runtime docs, dispatch/front-door hooks, a small contract package, and tests, with PR/mod behavior explicitly out of scope.

## Stage Report: implementation

- DONE: Pi dispatch uses a Pi-native tool contract instead of Claude-compatible tool names.
  Added `host: "pi"` support to `spacedock dispatch build`; `TestBuildPiHostPromptShape` fails on Claude team syntax and passes on the Pi read-dispatch-file / subagent-completion shape.
- DONE: Pi runtime adapters are loadable from shipped skills.
  Added `pi-first-officer-runtime.md` and `pi-ensign-runtime.md`, wired both SKILL.md files, and added `TestPiRuntimeAdaptersAreLoadable`.
- DONE: Live-harness OAuth isolation requirement is captured for the future live smoke.
  The Pi first-officer runtime doc now specifies isolated Pi config/session directories with copied auth, matching the entity's ideation decision.

### Summary

Implemented the first compatibility slice for Pi runtime support in code commit `f935a7b2`: `spacedock dispatch build` now accepts `host: "pi"` and emits Pi-native dispatch guidance without Claude team-tool signatures. The shipped skill surface now advertises and loads Pi first-officer/ensign runtime adapters; focused dispatch and skill integration tests pass.
