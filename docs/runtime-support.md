# Adding runtime support

Runtime support means Spacedock can launch or drive a host as a first officer, dispatch ensigns through that host's native agent mechanism, and prove the resulting workflow state. A host is not supported because its instructions mention Spacedock; it is supported when a live or fixture-backed run exercises the host and verifies durable state.

Use this guide when adding a new host such as Pi, or when turning a spike into a supported runtime lane.

## Runtime layers

Add support in small layers. Each layer should have its own proof.

1. **Skill adapters**
   - Add `skills/first-officer/references/<host>-first-officer-runtime.md`.
   - Add `skills/ensign/references/<host>-ensign-runtime.md`.
   - Wire both from the corresponding `SKILL.md` runtime-adapter section.
   - The adapter must name the host's native mechanism. Do not emulate Claude `Agent`, `SendMessage`, `TeamCreate`, or `TeamDelete` unless the host really provides those tools.

2. **Dispatch host mode**
   - Teach `spacedock dispatch build` to accept `host: "<host>"` when the assignment shape differs by host.
   - Keep entity paths and worktree paths explicit, especially for split-root workflows (`state: .spacedock-state`).
   - Test both positive shape and banned-tool negative cases.

3. **Runtime contracts and registries**
   - If the host has long-lived workers, define the minimum worker record: label, substrate, run/session handle, entity, stage, state, and completion epoch.
   - Reject stale completion evidence after follow-up or reuse. A previous completion must never satisfy a later assignment.
   - If the host has a team API, adapt Spacedock lifecycle intents to the host's native action schema.

4. **Launch/install UX**
   - Add `spacedock <host>` only after the manual/live harness proves the runtime path.
   - Add `spacedock install --host <host>` only when the install path is known and can be checked without mutating unrelated global host state.
   - Add `doctor --host <host>` when there is a manifest, package, or runtime health check to verify.

5. **Live runner**
   - Prove the host with a live-gated test when the claim is runtime integration.
   - Prefer a temp workflow fixture, isolated host config/session dirs, and copied credentials over global host state.
   - Assert process exit, entity content, git log, and clean state. Do not pass by transcript phrasing.

## Acceptance checklist

A new runtime support slice is not done until the entity or PR records evidence for each applicable item:

- Dispatch output uses the host-native contract and excludes incompatible host tool names.
- The first-officer and ensign skills load host runtime adapters.
- Split-root entity paths remain in the state checkout and are not rewritten into a code worktree.
- Follow-up/reuse cannot accept stale completion evidence, if reuse exists.
- Optional team substrates are represented as adapters over their real action schema.
- A live smoke proves the default dispatch path when runtime behavior is the claim.
- Install/launch commands exist only after the underlying mechanism is proven.

## Test strategy

Use the smallest proof at the same abstraction level as the claim:

- **Text claim:** parse or inspect the real instruction files.
- **Dispatch shape claim:** run `spacedock dispatch build` with a fixture and inspect emitted JSON/body.
- **Adapter claim:** table-test lifecycle intents to exact host-native payloads.
- **Registry claim:** unit-test persistence and stale epoch rejection.
- **Runtime claim:** live-gated host run that mutates a temp workflow and verifies durable state.

A substring search over code or prose is not proof of behavior. It is acceptable only when the claim itself is about text being present or absent.

## Manifesting from void

When a runtime seems unsupported on first contact, do not treat setup friction as proof the product path is impossible. Use a deliberate "assume it works" prompt to force the implementation loop to iron out auth, package paths, and tool-shape mismatches before declaring a blocker.

Use this operating prompt for the first implementation/validation loop:

```text
Assume <runtime> support is supposed to work. Do not treat missing polish, auth setup friction, or tool-shape mismatch as proof the runtime is impossible. In first-officer capacity, iron out the frictions:

- if auth is missing in an isolated harness, copy/reuse the existing host auth file correctly;
- if the dispatch substrate needs a local package/extension path, wire it explicitly;
- if the host tool shape differs from Claude/Codex, adapt to the host-native contract rather than emulating Claude tools;
- if a live test fails due to harness setup, fix the harness and rerun;
- only stop for a real product/design blocker, not for first-contact setup friction.
```

For Pi, the concrete version was:

```text
Assume Pi support is supposed to work. Do not treat missing polish, auth setup friction, or tool-shape mismatch as proof the runtime is impossible. In FO capacity, iron out the frictions:

- if Pi auth is missing in an isolated harness, copy/reuse the existing Pi OAuth auth file correctly;
- if the dispatch substrate needs a local package/extension path, wire it explicitly;
- if the Pi tool shape differs from Claude/Codex, adapt to the Pi-native contract rather than emulating Claude tools;
- if a live test fails due to harness setup, fix the harness and rerun;
- only stop for a real product/design blocker, not for first-contact setup friction.
```

That prompt matters because it changes the default failure interpretation. A missing `auth.json`, an extension not auto-discovered in a temp home, or a different subagent tool schema is harness work. A real blocker is a proven inability to launch, delegate, observe completion, or verify durable workflow state after the harness is correct.

## Pi live-smoke mechanism

The Pi proof used a live-gated test named:

```bash
go test -tags live -run TestLivePiSubagentEnsignSmoke ./internal/ensigncycle -v -count=1
```

The harness did this:

1. Resolve `pi` from `PATH` and the local Spacedock repo root.
2. Resolve the installed `pi-subagents` package root, defaulting to:

   ```text
   ~/.pi/agent/npm/node_modules/pi-subagents
   ```

3. Create temp runtime state:

   ```text
   PI_CODING_AGENT_DIR=<temp>
   PI_CODING_AGENT_SESSION_DIR=<temp>
   --session-dir <temp>
   HOME=<clean temp>
   ```

4. Copy only the operator's existing OAuth file into the isolated Pi home:

   ```text
   ~/.pi/agent/auth.json -> $PI_CODING_AGENT_DIR/auth.json
   ```

5. Launch `pi --print` with explicit local resources:

   ```text
   --extension ~/.pi/agent/npm/node_modules/pi-subagents/src/extension/index.ts
   --skill ~/.pi/agent/npm/node_modules/pi-subagents/skills/pi-subagents
   --skill <spacedock checkout>/skills/first-officer
   --skill <spacedock checkout>/skills/ensign
   ```

6. Create a temp split-root workflow:
   - `README.md` declares `state: .spacedock-state`.
   - The entity is folder-form in `.spacedock-state/pi-live-smoke/index.md`.
   - Both workflow root and state checkout are git repositories.

7. Ask the Pi parent to call `subagent(...)` exactly once.
8. Require the worker to append a stage report and commit only the state-checkout entity path.
9. Assert durable outcomes:
   - Pi process exits successfully.
   - Entity body contains the exact smoke marker and stage report shape.
   - State checkout git log contains the worker commit.
   - The entity path has no uncommitted changes.

## Exact Pi parent prompt

The live test formats this prompt with repository and temp paths. Keep the structure when debugging Pi runtime support; only substitute the paths and marker.

```text
You are the Spacedock first officer for a live Pi smoke test.

Use the pi-subagents subagent(...) tool exactly once to dispatch one Pi ensign worker. Do not use or mention Claude Agent, SendMessage, TeamCreate, or TeamDelete tools.

Dispatch a worker with agent "delegate" and this task:

Load and follow the local Spacedock ensign skill at <repo>/skills/ensign/SKILL.md and the Pi ensign adapter at <repo>/skills/ensign/references/pi-ensign-runtime.md. This is a split-root Spacedock workflow.

Workflow directory: <workflowRoot>
State checkout: <stateRoot>
Entity file: <entityPath>
Target stage: implementation

Required worker actions:
1. Read the workflow README and entity file.
2. Do not edit YAML frontmatter.
3. Append an implementation stage report to the entity body containing the exact marker PI-LIVE-SUBAGENT-ENSIGN-SMOKE, at least one '- DONE:' item, and a '### Summary' subsection.
4. Commit only the entity path in the state checkout with message 'ensign: pi live smoke'. Use a path-scoped git add/commit for pi-live-smoke/index.md.
5. Return a concise completion result naming the entity file and commit evidence.

After subagent(...) returns, you as first officer must verify the entity file contains PI-LIVE-SUBAGENT-ENSIGN-SMOKE and verify the state checkout git log contains 'ensign: pi live smoke'. Exit successfully only after those durable checks pass.
```

## Skill install and load paths

For Pi, `spacedock pi` launches the proven front door by loading local resources explicitly:

```text
<spacedock checkout>/skills/first-officer
<spacedock checkout>/skills/ensign
~/.pi/agent/npm/node_modules/pi-subagents/skills/pi-subagents
~/.pi/agent/npm/node_modules/pi-subagents/src/extension/index.ts
```

`spacedock install --host pi` is an idempotent readiness check and setup guide for that substrate; it does not install a Claude/Codex-style marketplace plugin and does not accept `--plugin-dir`. Resolve the local skill checkout by running it from the checkout or setting `SPACEDOCK_REPO_ROOT`. `spacedock doctor --host pi` reports the Pi CLI, auth file, `pi-subagents` extension/skill, and local Spacedock skill health; it still accepts `--plugin-dir <spacedock checkout>` for local skill checkout diagnostics. Live tests should not mutate global `~/.pi/agent`; they should keep using isolated Pi homes with copied auth.
