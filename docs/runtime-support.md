# Adding runtime support

Runtime support means Spacedock can launch or drive a host as a first officer, dispatch ensigns through that host's native agent mechanism, and prove the resulting workflow state. A host is not supported because its instructions mention Spacedock; it is supported when a live or fixture-backed run exercises the host and verifies durable state.

Use this guide when adding a new host such as Pi, or when turning a spike into a supported runtime lane.

## Runtime contract principles

Keep the shared Spacedock contract host-neutral. Shared workflow instructions should name lifecycle capabilities such as `«worker.spawn»`, `«addressable-worker»`, `«completion-signal»`, `«roster-reconcile»`, and `«worker.shutdown»`; runtime adapters own the concrete host tool calls that realize those capabilities.

Write adapters as positive bindings for the host that is running. Prefer "Pi maps `«worker.shutdown»` to ..." over "Pi has no Claude TeamDelete." Avoid negative contrast against another runtime unless the document is deliberately explaining a migration or compatibility hazard.

Do not couple adapter text to mutable procedure step numbers. If a shared procedure says to run `«worker.shutdown»` at the terminal boundary, adapters should bind `«worker.shutdown»`; they should not say "Merge-and-Cleanup step 10" or duplicate the shared teardown sequence.

Keep runtime tool names in runtime binding sections. Host-neutral core text and shared skills should call the capability name, not `SendMessage`, `followup_task`, `subagent`, `TeamDelete`, or any other concrete host tool. If a shared doc must mention host realizations for a transition period, guard it with contractlint so the list stays intentional and cannot drift into imperative procedure text.

Treat probes as the source of runtime truth. Do not file a `v2` or host-variant runtime solely because a version label changed; first probe the live tool surface and update the binding map when the lifecycle capability is the same. A separate runtime file is justified only when the lifecycle semantics differ enough that the shared capability contract would become misleading.

Prefer pseudo-code contracts over narrative instructions. Shared and runtime contracts should look like callable bodies keyed by capability names, with compact fields such as `guard`, `effect`, `done-when`, `block`, and `→` binding/status lines. Use prose only for fuzzy judgment, host quirks proven by probes, or rationale that cannot be encoded as an executable-shaped obligation.

A `→` line states where the capability is realized: `→ shipped` names the binary that already implements it (invoke it directly); `→ runtime-binding` defers to the host adapter (and, for a deferred-module section, names the core file the boot core loads at that load point); `→ prose` marks an obligation with no binary backing — judgment-owned when it names no verb, or a deterministic mechanism still hand-followed when it names a `becomes` verb a binary will later ship (e.g. `«dispatch.next-action»` → `spacedock dispatch next-action`, descoped to roadmap 0222). A deferred-module section (the boot core's dispatch/merge load points) may itself take this compact shape — a `→ runtime-binding` line naming the core file plus `done-when`/`guard` lines — rather than a narrative paragraph.

A deferred module realizes as a core reference file or a non-user-invocable skill, keyed on WHERE its host binding lives, not on whether it loads an adapter section alongside at the trigger. A module whose host binding lives in a deferred reference file stays a core file the boot core names by path via `→ runtime-binding`: dispatch and merge both bind through `references/claude-fo-dispatch.md` (the merge teardown loads there at first dispatch — there is no separate Claude merge reference), so both are references. A module whose host binding is boot-resident, or that has no host binding, realizes as a non-user-invocable skill the boot core invokes via `Skill(skill="spacedock:…")`: the status-viewer surface has no host binding, and the write/id-style surface's `spacedock new` binding lives boot-resident in the runtime adapter's new-entity section, so both are skills even though `spacedock new` is per-host.

A failure-triggered exception body whose trigger lives inside a host adapter's deferred reference realizes as a non-user-invocable skill too: the adapter keeps a resident trigger line that carries the first action and names the skill (`spacedock:fo-dispatch-recovery` — Degraded Mode, break-glass manual dispatch, budget-fail/dead-ensign handling), so the body loads only when its failure actually fires; the trigger line, not the boot core, is the namer.

### Runtime binding-block shape

A first-officer runtime adapter should default to a bindings block, not lifecycle prose. The shared core owns when capabilities are invoked; the runtime file owns how the host realizes them.

Preferred shape:

- `«worker.spawn»` -> host-native spawn call, mapping the shipped `«dispatch.build»` artifact's fields onto it (dispatch-build is its own host-neutral «fn», not part of this binding).
- `«addressable-worker»` -> PRESENT/ABSENT plus worker-to-FO and FO-to-worker message routes.
- `«async-dispatch»` -> async/blocking behavior and wait/poll mechanism.
- `«worker-identity»` -> handle, address, model stamp, and canonical model space.
- `«completion-signal»` -> observable completion signal; note that file verification remains the gate.
- `«worker.shutdown»` -> terminal/supersede shutdown binding.
- `«context-budget»` -> probe binding or ABSENT.
- `«roster-reconcile»` -> reconcile binding or ABSENT.

Keep residual sections short and factual: live-tool probes, harness isolation, compatibility notes, or host-specific guardrails that do not fit a capability name. Do not re-narrate the shared dispatch, await, reuse, gate, or merge lifecycle in the adapter.

If a host-specific rule is load-bearing, first try to attach it to the relevant capability bullet. Add a separate prose section only when the rule spans multiple capabilities or documents a probe/harness concern.

An ensign runtime adapter follows the same default: a `## Runtime implementation` bindings block, not lifecycle prose. The shared ensign core (`ensign-shared-core.md`) owns the discipline (assignment reading, worktree, split-root commit, frontmatter, proof, stage report); the adapter binds only the concerns that differ by host. The ensign block is keyed by the ensign-controlled concern, not the FO's worker-lifecycle capabilities.

Preferred shape:

- `Clarification` -> the host channel for asking the FO a blocking or non-blocking question.
- `Completion signal` -> the observable signal the ensign emits when done; note that the FO file-verifies the stage report as the gate.
- `Captain communication` -> the direct-to-captain channel when the stage involves captain interaction (host-only; omit where the host has no such channel).
- `Shutdown response` -> the cooperative shutdown acknowledgement (host-only; omit where the host has no mailbox shutdown).

Do not re-narrate the shared dispatch, worktree, split-root, frontmatter, path-scoped-commit, or feedback-routing discipline in the ensign adapter — the shared ensign core carries it. Omit a bullet whose concern the host does not have rather than binding it to a negative contrast against another host.

## Runtime layers

Add support in small layers. Each layer should have its own proof.

1. **Skill adapters**
   - Add `skills/first-officer/references/<host>-first-officer-runtime.md`.
   - Add `skills/ensign/references/<host>-ensign-runtime.md`.
   - Wire both from the corresponding `SKILL.md` runtime-adapter section.
   - The adapter must name the host's native mechanism. Do not emulate Claude `Agent`, `SendMessage`, `TeamCreate`, or `TeamDelete` unless the host really provides those tools.

2. **Dispatch host mode**
   - Teach `spacedock dispatch build` to accept `host: "<host>"` when the assignment shape differs by host.
   - If a generic dispatch mode has no faithful native call shape, reject that host/mode combination before artifact creation; do not silently reinterpret it.
   - Codex fresh dispatch is always named, so `host=codex` rejects `bare_mode`. Its fresh prompt begins `$spacedock:ensign; then Read ...`, which loads the installed Spacedock ensign skill before the child reads the dispatch pointer. Reuse/advance prompts do not repeat this prefix because they target the worker that already loaded the shared contract.
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

## Boot guard at the compaction boundary

A compaction-resumed session keeps its session id and transcript; the host
records the boundary durably (a `compact_boundary` record in the session
transcript). `status --boot` writes a one-line per-session receipt; the
authority verbs — `gate record`, `gate consume`, `merge guard` — refuse
(exit 4) when the receipt is missing or older than the latest boundary, until
boot re-runs. Detection resolves per host: Claude Code via
`CLAUDE_CODE_SESSION_ID` plus the recorded transcript path; hosts without a
resolvable identity degrade to a silent no-op (Codex: #595). The guard fails
open on unreadable transcripts and never needs a hook.

## Launcher binary propagation through wrappers

`spacedock claude` and `spacedock codex` attach `SPACEDOCK_BIN` to the host process they spawn, including the outer `safehouse -- ...` process when safehouse wrapping is active. Spacedock does not modify safehouse internals or assume a private passthrough mechanism; if a wrapper or runtime strips `SPACEDOCK_BIN` before the agent session observes it, the skill contract's `${SPACEDOCK_BIN:-spacedock}` convention degrades to the existing `$PATH` lookup.

The launcher stays resident as the host's parent: it spawns the host as a child, inherits the terminal, forwards externally-targeted signals (`SIGTERM`/`SIGHUP`) while letting terminal signals (Ctrl-C, resize) reach the host through the shared foreground process group, and exits with the host's exit code — rather than replacing itself with the host. This keeps the `spacedock <host> …` command legible in process listings and session managers (for example zellij's restart view) and lets the launcher supervise companion processes alongside the session in future. (Unix launch lane; `spacedock <host>` is not a supported launch path on Windows.)

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

`spacedock install --host pi` installs Spacedock as a Pi package via `pi install` (idempotent; pass `--check` for a readiness check without installing). It is not a Claude/Codex-style marketplace plugin, and it accepts `--plugin-dir` as the dev-override install source (defaulting to the published package otherwise). Resolve the local skill checkout by running it from the checkout or setting `SPACEDOCK_REPO_ROOT`. `spacedock doctor --host pi` reports the Pi CLI, auth file, `pi-subagents` extension/skill, local Spacedock skill health, and supervisor-talkback setup prerequisites: the `pi-subagents` intercom bridge source, the resolved `PI_INTERCOM_PACKAGE_ROOT` package root, and the `pi-intercom` skill resource. Current `pi-subagents`/`pi-intercom` packages do not expose stable `pi-intercom` or `subagents-doctor` PATH commands, so readiness is based on package/resource paths instead of command shims. These doctor/install checks are necessary setup checks but insufficient to prove live supervisor talkback. Live proof still requires the cq-style `pi-intercom-supervisor-talkback` probe: progress update -> decision request -> supervisor reply -> child resume -> durable marker evidence. Live tests should not mutate global `~/.pi/agent`; they should keep using isolated Pi homes with copied auth.
