# Adding a runtime

A host is supported when a live or fixture-backed run launches it as a first officer, dispatches an ensign through that host's native agent mechanism, and verifies durable workflow state: process exit, entity body, state-checkout git log, and clean status. A host is not supported because its instructions mention Spacedock, and a substring search over code or prose is not proof of behavior. Spacedock ships `spacedock claude` and `spacedock codex` as proven front doors; adding a host means earning the same level of proof.

This page is the full procedure: what "supported" means, the layers to add support in, the acceptance checklist, and a worked Pi example. Use it when adding a new host such as Pi, or when turning a spike into a supported runtime lane.

## What "supported" means

The supported claim has four observable parts. The runtime must be able to:

- **Launch** the host as the first officer.
- **Delegate** to an ensign through the host's own subagent or team mechanism.
- **Observe** that the ensign completed.
- **Verify** durable state: process exit, entity body, state-checkout git log, and clean status.

If you cannot demonstrate all four against a real or fixture-backed run, the host is a spike, not a runtime.

## Add support in layers

Add support in small layers, each with its own proof at its own abstraction level. The order matters: every later layer assumes the earlier one is proven.

1. **Skill adapters.** Add `skills/first-officer/references/<host>-first-officer-runtime.md` and `skills/ensign/references/<host>-ensign-runtime.md`, and wire both from the matching `SKILL.md` runtime-adapter section. The adapter must name the host's native mechanism. Do not emulate Claude `Agent`, `SendMessage`, `TeamCreate`, or `TeamDelete` unless the host actually provides those tools.
2. **Dispatch host mode.** Teach `spacedock dispatch build` to accept `host: "<host>"` when the assignment shape differs by host. If a generic dispatch mode has no faithful native call shape, reject that host/mode combination before artifact creation; do not silently reinterpret it. Codex fresh dispatch is always named, so `host=codex` rejects `bare_mode`. Its fresh prompt begins `$spacedock:ensign; then Read ...`, which loads the installed Spacedock ensign skill before the child reads the dispatch pointer. Reuse/advance prompts do not repeat this prefix because they target the worker that already loaded the shared contract. Keep entity and worktree paths explicit, especially for split-root workflows (`state: .spacedock-state`). Test the positive shape and the banned-tool negative case.
3. **Runtime contracts and registries.** For long-lived workers, define the minimum worker record (label, substrate, run/session handle, entity, stage, state, completion epoch) and reject stale completion evidence: a previous completion must never satisfy a later assignment. For a host team API, adapt Spacedock lifecycle intents to the host's native action schema.
4. **Launch/install UX.** Add `spacedock <host>` only after a manual or live harness proves the runtime path. Add `spacedock install --host <host>` only when the install path is known and checkable without mutating unrelated global host state. Add `spacedock doctor --host <host>` when there is a manifest, package, or runtime health check to verify.
5. **Live runner.** Prove the host with a live-gated test when the claim is runtime integration. Use a temp workflow fixture, isolated host config and session directories, and copied credentials rather than global host state. Assert process exit, entity content, git log, and clean status. Never pass on transcript phrasing.

## Launcher binary propagation through wrappers

`spacedock claude` and `spacedock codex` attach `SPACEDOCK_BIN` to the process they exec, including the outer `safehouse -- ...` process when safehouse wrapping is active. Spacedock does not modify safehouse internals or assume a private passthrough mechanism; if a wrapper or runtime strips `SPACEDOCK_BIN` before the agent session observes it, the skill contract's `${SPACEDOCK_BIN:-spacedock}` convention degrades to the existing `$PATH` lookup.

## Match the proof to the claim

Use the smallest proof at the same abstraction level as the claim:

- **Text claim** (an adapter mentions the right tool): parse or inspect the real instruction files.
- **Dispatch shape claim:** run `spacedock dispatch build` with a fixture and inspect the emitted JSON or body.
- **Adapter claim:** table-test lifecycle intents to exact host-native payloads.
- **Registry claim:** unit-test persistence and stale-epoch rejection.
- **Runtime claim:** a live-gated host run that mutates a temp workflow and verifies durable state.

The failure mode to guard against: declaring a host supported because its prose looks right. Substring presence is acceptable proof only when the claim itself is about text being present or absent.

## Acceptance checklist

A new runtime support slice is not done until the entity or PR records evidence for each applicable item:

- Dispatch output uses the host-native contract and excludes incompatible host tool names.
- The first-officer and ensign skills load host runtime adapters.
- Split-root entity paths remain in the state checkout and are not rewritten into a code worktree.
- Follow-up/reuse cannot accept stale completion evidence, if reuse exists.
- Optional team substrates are represented as adapters over their real action schema.
- A live smoke proves the default dispatch path when runtime behavior is the claim.
- Install/launch commands exist only after the underlying mechanism is proven.

## Manifesting from void

When a runtime looks unsupported on first contact, do not read setup friction as proof the product path is impossible. A missing `auth.json`, an extension not auto-discovered in a temp home, or a subagent tool schema that differs from Claude's is harness work, not a blocker. A real blocker is a proven inability to launch, delegate, observe completion, or verify durable state *after* the harness is correct.

Run the first implementation and validation loop under a deliberate "assume it works" prompt so the loop irons out auth, package paths, and tool-shape mismatches before anyone declares a wall:

```text
Assume <runtime> support is supposed to work. Do not treat missing polish, auth setup friction, or tool-shape mismatch as proof the runtime is impossible. In first-officer capacity, iron out the frictions:

- if auth is missing in an isolated harness, copy/reuse the existing host auth file correctly;
- if the dispatch substrate needs a local package/extension path, wire it explicitly;
- if the host tool shape differs from Claude/Codex, adapt to the host-native contract rather than emulating Claude tools;
- if a live test fails due to harness setup, fix the harness and rerun;
- only stop for a real product/design blocker, not for first-contact setup friction.
```

The prompt earns its place by changing the default interpretation of a failure: harness work gets fixed in-loop, and only a proven product or design blocker stops the work.

## The worked Pi runtime

Pi is the worked example of a host taken from spike to supported runtime: the live-smoke mechanism, the exact parent prompt, the install and doctor surface, and the skill load paths.

### Pi live-smoke mechanism

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

For Pi, the concrete "assume it works" prompt was:

```text
Assume Pi support is supposed to work. Do not treat missing polish, auth setup friction, or tool-shape mismatch as proof the runtime is impossible. In FO capacity, iron out the frictions:

- if Pi auth is missing in an isolated harness, copy/reuse the existing Pi OAuth auth file correctly;
- if the dispatch substrate needs a local package/extension path, wire it explicitly;
- if the Pi tool shape differs from Claude/Codex, adapt to the Pi-native contract rather than emulating Claude tools;
- if a live test fails due to harness setup, fix the harness and rerun;
- only stop for a real product/design blocker, not for first-contact setup friction.
```

### Exact Pi parent prompt

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

### Skill install and load paths

For Pi, `spacedock pi` launches the proven front door by loading local resources explicitly:

```text
<spacedock checkout>/skills/first-officer
<spacedock checkout>/skills/ensign
~/.pi/agent/npm/node_modules/pi-subagents/skills/pi-subagents
~/.pi/agent/npm/node_modules/pi-subagents/src/extension/index.ts
```

`spacedock install --host pi` installs Spacedock as a Pi package via `pi install` (idempotent; pass `--check` for a readiness check without installing). It is not a Claude/Codex-style marketplace plugin, and it accepts `--plugin-dir` as the dev-override install source (defaulting to the published package otherwise). Resolve the local skill checkout by running it from the checkout or setting `SPACEDOCK_REPO_ROOT`. `spacedock doctor --host pi` reports the Pi CLI, auth file, `pi-subagents` extension/skill, local Spacedock skill health, and supervisor-talkback setup prerequisites: the `pi-subagents` intercom bridge source, the resolved `PI_INTERCOM_PACKAGE_ROOT` package root, and the `pi-intercom` skill resource. Current `pi-subagents`/`pi-intercom` packages do not expose stable `pi-intercom` or `subagents-doctor` PATH commands, so readiness is based on package/resource paths instead of command shims. These doctor/install checks are necessary setup checks but insufficient to prove live supervisor talkback. Live proof still requires the cq-style `pi-intercom-supervisor-talkback` probe: progress update -> decision request -> supervisor reply -> child resume -> durable marker evidence. Live tests should not mutate global `~/.pi/agent`; they should keep using isolated Pi homes with copied auth.
