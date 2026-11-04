---
title: Preserve orchestrator control in launcher passthrough mode
status: ideation
source: "Captain authorized issue 813 on 2026-11-03: Proceed with the 3. You have the conn to approve gate and open pr"
issue: spacedock-dev/spacedock#813
sprint: paseo-integration
id: xerfqfkayjawte0gjem7v8v1
gates:
    version: 1
    records:
        - id: gate:xerfqfkayjawte0gjem7v8v1:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:xerfqfkayjawte0gjem7v8v1-backlog-1
              briefing:
                id: briefing:xerfqfkayjawte0gjem7v8v1:backlog:attempt-1:revision-1
                digest: sha256:2feac6d7c7831a290eb2866b9ca6bbd0d50eb165039916586156b5fb8a6b4e0d
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:xerfqfkayjawte0gjem7v8v1:backlog:1
                briefing: briefing:xerfqfkayjawte0gjem7v8v1:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-04T07:48:27.359557Z"
                decision: approve
                reason: The seed defines the reported orchestrator failures, retains sandbox behavior, protects ordinary launches, and requires each host boundary to be exercised.
                conn:
                    quote: Proceed with the 3. You have the conn to approve gate and open pr
                    source: Captain message in this Codex conversation on 2026-11-03
              application:
                target-stage: ideation
                state: consumed
started: 2026-11-04T07:48:58Z
---
Paseo can use Spacedock as a provider command while retaining control of its protocol, agent identity, prompts, and permission UI.

## Problem

[Issue 813](https://github.com/spacedock-dev/spacedock/issues/813) describes three provider protocols: Claude stream-json, Codex app-server, and Pi RPC.
The current launcher treats each provider process as an interactive first officer.
Its additions break Codex startup and replace Claude permission settings inside Safehouse.
Pi also receives launch markers that activate the first officer through its extension.
The orchestrator needs the existing sandbox wrapper while it owns the host session.

## Proposed approach

Add the Boolean flag `--passthrough` before `--` on Claude, Codex, and Pi.
The host argument sequence after `--` remains exact, including order, empty values, and repeated flags.
The inner command contains only the resolved host executable and those arguments.
The existing resident launcher still inherits stdio, forwards signals, and returns the host exit code.

In passthrough mode:

- Claude receives no launch-owned `--agent`, bootstrap prompt, permission flag, adjacent plugin selection, compatibility check, or automatic installation.
  The launcher does not add `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS`. It preserves an existing caller value.
- Codex receives no bootstrap prompt, approval flag, sandbox-bypass flag, collaboration layer, plugin discovery, installation, or compatibility check.
  The launcher also skips its reserved collaboration-key rejection.
  Caller-supplied `-c`, `--config`, `--enable`, and `--disable` values reach Codex unchanged.
- Pi receives no launch-owned extension or skill arguments, task prompt, runtime-readiness report, or first-officer markers.
  Both `PI_SPACEDOCK_LAUNCH` and `PI_SPACEDOCK_LAUNCH_TASKLESS` are absent from the child environment, even if the parent supplied them.
  Installed resources remain subject to Pi's native discovery and the caller's native flags.
- No first-officer banner or stray-prompt advisory appears.
  Launcher errors remain on stderr, and stdout carries host output only.
- `SPACEDOCK_BIN`, Safehouse selection, `--trust-workdir-config`, terminal allowances, and existing directory/profile flags retain their behavior.
  Pi retains its fnm path resolution and the matching sandbox grant for code and dependencies.
  A selected sandbox with no Safehouse executable still fails before launch.

Passthrough rejects task text and Spacedock's `--plugin-dir` before `--`, with a stderr diagnostic before host or plugin activity.
All native prompts and plugin arguments belong after `--` in this mode.
This prevents silent loss of task text or persistent installation from a flag that requests an interactive Spacedock provider.
`--no-install` remains accepted where it already exists and has no additional effect in passthrough.

Every host accepts `--skip-compat-check` before `--`.
For Claude and Codex, the existing interactive bypass remains unchanged.
For Pi, the flag is an accepted compatibility alias with no effect on its ordinary readiness checks.
Passthrough already skips those checks, so it needs no separate bypass flag.
No runtime version floor or protocol support claim changes.

### Sandbox environment and precedence

The two environment settings apply to all three hosts, in both launch modes.
They use the existing sandbox selection and flag translator.
Explicit pre-fence flags take precedence over their corresponding environment setting.

| Input | Effect and precedence |
|---|---|
| `SPACEDOCK_SAFEHOUSE=1` | Requests Safehouse without a profile. `--safehouse=true` or `--safehouse=false` replaces this request. |
| `SPACEDOCK_SAFEHOUSE=0` or an unset/empty value | Adds no sandbox request. A `.safehouse` profile or another sandbox flag still selects Safehouse. |
| Other nonempty Boolean value | Uses `strconv.ParseBool`. Invalid values fail before launch unless an explicit `--safehouse` flag replaces them. |
| `SPACEDOCK_SAFEHOUSE_ENABLE=ssh,docker` | Supplies the capability list only if no `--safehouse-enable` flag occurs. A nonempty list also selects Safehouse. |
| Repeated `--safehouse-enable` flags | Replace the environment list as one group and retain their existing order and comma splitting. |
| `.safehouse`, directory grants, or profile flags | Still select Safehouse independently. Neither `SPACEDOCK_SAFEHOUSE=0` nor `--safehouse=false` disables these triggers. |

An explicit empty capability flag retains its current native validation behavior.
Capability names and profile content still use Safehouse's existing validation.
This task adds no environment settings for directory grants or extra profiles.

### Mechanism and value

| Mechanism | Value AC | Simplest alternative and its limit |
|---|---|---|
| One opt-in flag with local guards around interactive additions | AC-1, AC-2, AC-4 | Detecting protocol tokens couples Spacedock to each host grammar and misses identity-only provider sessions. |
| Skip Codex collaboration injection and its override rejection in this mode | AC-1, AC-2 | Retaining the layer still replaces caller configuration and requires support for Spacedock's selected table. |
| Suppress Pi launch markers and explicit resources | AC-1, AC-2 | Removing only argv prompts leaves the extension free to send a first-officer user message. |
| One small helper for environment fallback and flag presence | AC-3 | Fixed provider argv cannot vary sandbox capabilities per session. Three separate parsers can disagree on precedence. |
| Reject interactive task/provider flags in the new mode | AC-1, AC-2 | Ignoring them loses caller intent. Reinterpreting them recreates prompt assembly or persistent plugin setup. |

Keep the host functions and existing wrapper calls.
The environment helper needs no new host dispatcher, protocol parser, dependency, or process controller.

## Risk evidence

Baseline: `origin/main` at `15e0baea3`, inspected in an isolated snapshot and detached checkout.
The shared root checkout remained unchanged.
Probe versions: Codex 0.166.0, Claude Code 2.1.303, Pi 1.3.0, and Agent Safehouse 0.12.0.
These are observed versions, not new minimums.

1. **Codex protocol failure, reproduced.** Native `codex app-server` returned a JSON-RPC initialize result with `id: 1`.
   The built launcher with `codex --skip-compat-check -- app-server` exited 2 with zero stdout bytes.
   Its native diagnostic was `error: unrecognized subcommand '<bootstrap prompt>'`.
   With the executable beside a checkout, the launcher also installed that checkout in the isolated `CODEX_HOME` before failing.
2. **Claude permission failure, reproduced through real Safehouse.** A wrapped stream-json session with `--permission-mode default` emitted `permissionMode: default`.
   Adding only `--dangerously-skip-permissions` changed the init result to `permissionMode: bypassPermissions`.
   A separate wrapped environment probe returned `APP_SANDBOX_CONTAINER_ID=agent-safehouse`.
   Both permission probes kept stdin open and stopped after the init event.
3. **Pi protocol baseline, exercised.** Native Pi RPC returned `type: response`, `command: get_state`, `success: true`, and the supplied request id.
   The probe used an isolated agent directory, `--no-session`, `--offline`, `--no-extensions`, `--no-skills`, and `--no-mcp`.
   This proves startup without a model call or a Spacedock prerequisite.
4. **Existing seam exercise, four failing subtests.** The existing `fakeHost` and `fakePiRuntimeOps` recorded the actual launcher behavior.
   Codex added six collaboration tokens, an approval pair, and a bootstrap prompt before/after `app-server`.
   Wrapped Claude added skip-permissions, `--agent spacedock:first-officer`, and a bootstrap prompt around the caller arguments.
   Pi added extension/skill arguments and both first-officer markers with value `1`.
   Pi's parser rejected `--skip-compat-check` as an unknown flag.
   Each subtest compared the observed argv or parser result with an independent literal provider expectation.

The smallest hypothesis is that interactive launch additions cause these failures.
The native protocol results and Claude's one-flag comparison support that hypothesis.
No process-supervision or exec-in-place change is needed.

Reproduction inputs for implementation and validation:

- Keep Codex stdin open and send `{"id":1,"method":"initialize","params":{"clientInfo":{"name":"paseo-passthrough-spike","version":"0.0.0"}}}`.
  Read the response for id `1`, then terminate the process.
- Start `safehouse --trust-workdir-config -- claude --print --input-format stream-json --output-format stream-json --verbose --permission-mode default --no-session-persistence`.
  Send `{"type":"user","message":{"role":"user","content":"Reply with OK. Use no tools."}}` and read the system init event.
  Repeat with only `--dangerously-skip-permissions` added after `claude`.
- Start `pi --mode rpc --no-session --offline --no-extensions --no-skills --no-mcp` in an isolated `PI_CODING_AGENT_DIR`.
  Send `{"id":"paseo-spike","type":"get_state"}` and read that response.
- Reuse the existing fake host seams with the same arguments and literal expected host argv.
  For Pi, supply the existing healthy resource fixture so the exercise reaches its current launch path.

The final candidate must repeat protocol startup through `spacedock <host> --passthrough`, rather than claiming native startup alone proves the change.

### Baseline check outcomes

These checks used the unchanged detached `15e0baea3` checkout, separate from the one-off failing spike.
The shared product checkout remained unchanged.

- `go test ./...`: exit 1. The CLI and ensigncycle packages reached the default ten-minute timeout.
  The active tests included `TestClaudeInstallLeavesOnlySelectedChannelEnabled` and `TestRecordedGateCommittedBeforeDispatchResolutionPaths`.
- `go test ./... -race`: exit 1. The same packages reached the ten-minute timeout.
  The active CLI test was `TestDoctorOnHomebrewNamesTheCompassFormulaWhenItsProgramIsMissing`.
- Both runs also failed `TestBuildInputFailuresDoNotReadOrStamp` with `input failure mutated entity, commit, or dispatch artifact`.
  The runs overlapped, and this fixture uses the shared `/tmp/spacedock-dispatch/thing-backlog.md` path.
  Artifact interference remains an unconfirmed explanation.
  A focused repeat during the remaining suite reported the same failure.
- Both Pi extension tests failed with `Library not loaded: /opt/homebrew/opt/llhttp/lib/libllhttp.9.3.dylib` from `/opt/homebrew/Cellar/node/25.9.0_2/bin/node`.
  A focused repeat reproduced both failures.
  A repeat with per-command PATH through `/Users/clkao/.local/share/fnm/node-versions/v24.13.1/installation/bin` passed both extension tests in 0.714 seconds.
  Native Pi RPC also passed through a working Node installation.
- The focused ordinary launcher checks passed in 5.036 seconds.
  Their existing owners cover permission defaults, Codex collaboration, Pi launch markers, binary propagation, fnm resolution, resident parentage, and exit codes.
  Removing a default, marker, env-pass allowance, or resident parent breaks its respective owner.
- `gofmt -w ./cmd ./internal` completed and exposed pre-existing two-line alignment drift in `internal/release/runtime_live_evidence_workflow_test.go`.
  The scratch-only formatting change was restored, and the detached checkout remained clean.

The full suites therefore do not supply a green regression baseline for this stage.
The ideation proof is the reproduced protocol/permission behavior and the explicit implementation test plan.
Full regression proof remains required during implementation and validation.
Local full logs: `/tmp/spike-paseo-provider-passthrough.CZz6MX/baseline-go-test.log` and `baseline-go-test-race.log`.

## Out of scope

A new Paseo provider, a replacement process supervisor, a release, unrelated launcher refactoring, and changes to host-native discovery.
The optional exec-in-place suggestion is not required.
The task adds no protocol-specific command parser, collaboration policy, new readiness prerequisite, or standing enforcement process.

## Expected surface and tolerance

Estimate net LOC change: +400, across 13 files.
Estimated insertions: 470. Estimated deletions: 70.
Tolerance: net +280 through +520 LOC, with at most 15 changed product files.
The entity and its state report are outside this product estimate.
A change outside these semantics requires a design decision even if the line count stays within tolerance.

| Files | Intended change |
|---|---|
| `internal/cli/frontdoor.go` | Parse the mode, record explicit sandbox flag presence, and guard Claude/Codex interactive additions and side effects. |
| `internal/cli/launch_options.go` (new) | Resolve the two environment fallbacks through the existing sandbox representation. |
| `internal/cli/pi.go` | Parse the common flags, guard interactive readiness/resources/markers, and preserve Pi path and wrapper behavior. |
| `internal/cli/help.go` | Describe the mode, the two environment settings, and the Pi compatibility alias. |
| `internal/cli/frontdoor_parse_test.go` | Exercise the mode/fence grammar and invalid interactive combinations. |
| `internal/cli/frontdoor_test.go` | Exercise literal Claude/Codex provider argv, plugin inactivity, child environment, quiet launch output, and host exit codes. |
| `internal/cli/codex_multi_agent_test.go` | Exercise caller collaboration overrides in passthrough while retaining ordinary collaboration checks. |
| `internal/cli/pi_frontdoor_test.go` | Exercise native provider argv, absent inherited markers, readiness independence, wrapping, and common flags. |
| `internal/cli/launch_parity_test.go` | Exercise environment fallback and explicit flag precedence for all hosts and both modes. |
| `internal/cli/help_test.go` | Extend existing help coverage for the new mode and settings. |
| `docs/site/reference/command-reference.md` | Document provider mode and qualify interactive collaboration/permission behavior. |
| `docs/site/reference/sandbox.md` | Document the two environment settings and precedence. |
| `docs/site/contributing/adding-a-runtime.md` | Qualify the existing Pi marker claim as interactive behavior and link to provider mode. |

Permitted semantics: new pre-fence launch grammar, sandbox environment fallback, passthrough ownership of host arguments and child environment, and skipped interactive setup.
Stored formats, workflow authority, host permission defaults, ordinary launch behavior without these new inputs, and process supervision remain unchanged.

## Acceptance criteria

**AC-1**: Passthrough Codex app-server and Claude/Pi protocol launches receive valid host argv and a clean stdout channel.
Verified by: literal argv cases in the existing launcher fixtures and candidate protocol startup with the reproduction inputs above. Adding a bootstrap prompt or banner to stdout must fail.

**AC-2**: The orchestrator retains agent identity, prompting, and permission settings in wrapped and unwrapped passthrough launches.
Verified by: recorded argv/environment fixtures, caller Codex overrides, and candidate Claude stream init with `permissionMode: default` under Safehouse. Injecting a first-officer marker, identity, resource, collaboration value, or permission override must fail.

**AC-3**: Passthrough keeps sandbox wrapping, capability selection, and SPACEDOCK_BIN propagation, including per-session SPACEDOCK_SAFEHOUSE and SPACEDOCK_SAFEHOUSE_ENABLE settings.
Verified by: the cross-host environment/argv matrix, existing scrubbing-wrapper environment proof, and Pi's stable-path fixtures. Losing wrapping, capability propagation, flag precedence, or the selected binary must fail.

**AC-4**: Every host accepts --skip-compat-check, passthrough skips compatibility gating, and ordinary interactive behavior remains stable.
Verified by: missing/incompatible-plugin fixtures that still reach passthrough, Pi common-flag cases, and existing interactive launcher regressions. Plugin activity or an interactive behavior change must fail.

**AC-5**: Help and product documentation explain the supported mode and environment settings.
Verified by: existing command-help behavior tests, the runnable provider examples below, and review of the concrete documentation changes. An unsupported flag, incorrect precedence, or a provider example that cannot start must fail.

## Test plan

Add the focused failing cases before product edits.
Reuse the existing proof owners. No new protocol harness or mandatory live journey is required.
Literal expected provider argv must not come from production constants or assembly helpers.

| Proof owner and type | Observable claim and distinct falsifying edit | Cost |
|---|---|---|
| `frontdoor_parse_test.go`, deterministic | The mode consumes only pre-fence flags and refuses interactive task/provider input before side effects. Accepting and dropping a task breaks this proof. | Seconds. Small table. |
| `frontdoor_test.go`, deterministic | Claude/Codex inner argv exactly matches caller protocol arguments in both wrap branches. Host exit 7 survives. No plugin work or launch text occurs. Restoring any interactive injection or gate call breaks it. | Seconds. Extend existing fake seams. |
| `codex_multi_agent_test.go`, deterministic | Explicit caller collaboration settings reach app-server unchanged. Leaving the reserved-key rejection active breaks this case while plain app-server fixtures can still pass. | Seconds. Extend the current override owner. |
| `pi_frontdoor_test.go`, deterministic | Pi starts without FO prerequisites and emits no launch-owned resources or inherited FO markers. Removing only argv additions while keeping an inherited marker breaks this proof. | Seconds. Reuse the Pi fake and healthy/bare fixtures. |
| `launch_parity_test.go`, deterministic | Both environment settings select the expected wrapper/capabilities, and explicit false/list flags replace environment values across hosts and modes. Appending both capability sources or letting false disable a profile breaks this proof. | Seconds. One shared input table. |
| Existing `TestSafehouseEnvPassForwardsSpacedockBin` and Pi stable-path tests, deterministic | A scrubbing wrapper still resolves the launched binary. Pi code/dependencies remain visible under fnm. Removing the env-pass allowance or stable-path grant breaks these owners. | Seconds. Unchanged tests. |
| Existing permission, banner, collaboration, host launch, and PTY tests, deterministic | Ordinary launches retain their identity, defaults, terminal behavior, signals, and exit codes. Broadly disabling their defaults or changing the supervisor breaks existing owners. | Seconds to minutes. Unchanged tests. |
| Existing help tests plus documentation review | Help exposes the supported mode and Pi alias. Examples agree with the host commands. Removing registration or documenting additive capability precedence breaks the check. | Seconds plus review. |
| One-off candidate protocol startup, real local hosts | Codex returns initialize id 1. Claude reports `default` under Safehouse. Pi returns successful `get_state`. Stdout contains only native protocol records. Adding an app-server prompt or wrapped permission bypass breaks this proof. | Under one minute for startup after local host setup. No model response or full workflow needed. |

The environment matrix covers unset/empty values, force-on, false, invalid Boolean input, an explicit Boolean override, capability fallback, repeated explicit capabilities, profile-only wrapping, and missing Safehouse.
Marker cases begin with both Pi markers already set in the parent.
Plugin cases place a valid checkout beside the selected launcher and record zero discovery/install calls and unchanged isolated host settings.

Run `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` on the implementation worktree.
Run the full suites sequentially because dispatch fixtures share a fixed artifact directory.
Use the documented runtime-support prompt to resolve harness setup friction before declaring a runtime blocker.
Validation must run the existing detached adversarial audit for launcher changes.
Useful adversarial edits restore one bootstrap prompt, one wrapped bypass flag, one Pi marker, or one reversed precedence branch.
The primary owner for that failure must reject each edit.
The First Officer must arrange an independent read-only ideation review before the gate, with necessity reviewed before soundness.

## Product documentation diff

In `docs/site/reference/command-reference.md`, add this section before `## Setup`:

```markdown
### Provider commands

`--passthrough` keeps the sandbox wrapper while your orchestrator controls the host session.
Place the flag before `--` and all native host arguments after it.

| Host | Provider command |
|---|---|
| Claude | `spacedock claude --passthrough --` |
| Codex | `spacedock codex --passthrough --` |
| Pi | `spacedock pi --passthrough --` |

The orchestrator appends its stream-json, `app-server`, or RPC arguments to that command.
For example, `spacedock codex --passthrough -- app-server` serves Codex JSON-RPC on stdio.
Spacedock adds no agent, prompt, permission flag, collaboration setting, or first-officer banner.
It skips plugin setup and compatibility checks.
Pi also receives no launch-owned extension, skill, or first-officer markers.
Native host discovery and caller configuration still apply.
Stdout carries host output only, and the launcher returns the host exit code.

Task text and Spacedock's `--plugin-dir` before `--` are invalid in this mode.
For Claude, native additional plugin arguments can follow `--`.
All three hosts accept `--skip-compat-check`. The flag has no additional effect in passthrough.
On an ordinary Pi launch, it leaves the existing readiness checks active.

Sandbox flags and [per-session environment settings](sandbox.md) apply in both launch modes.
```

Also replace `Every Codex launch, including local-plugin, Safehouse, and resume` with `Every interactive Codex launch, including local-plugin, Safehouse, and resume`.
Replace the next paragraph's opening `Spacedock reserves` with `In interactive mode, Spacedock reserves`.
Insert `These defaults apply to interactive launches.` before the existing unsandboxed permission paragraph.
Keep its existing details unchanged.

In `docs/site/reference/sandbox.md`, extend the Safehouse trigger cell with `SPACEDOCK_SAFEHOUSE=1` and a nonempty `SPACEDOCK_SAFEHOUSE_ENABLE`.
Add this section before `## Additional profiles`:

```markdown
## Per-session settings

These environment settings apply to Claude, Codex, and Pi in both launch modes.

| Setting | Value | Default |
|---|---|---|
| `SPACEDOCK_SAFEHOUSE` | `1` requests Safehouse. `0` adds no request. Other Boolean spellings follow Go's `strconv.ParseBool`. | No request |
| `SPACEDOCK_SAFEHOUSE_ENABLE` | Comma-separated capabilities, such as `ssh,docker`. A nonempty list also requests Safehouse. | No capabilities |

An explicit `--safehouse=true` or `--safehouse=false` replaces the Boolean environment value.
Any `--safehouse-enable` flags replace the environment capability list and retain their existing order.
A `.safehouse` profile or another sandbox flag still requests Safehouse independently.
A false Boolean value does not disable those triggers or an existing parent sandbox.
Invalid Boolean values fail before launch unless an explicit flag replaces them.

`SPACEDOCK_SAFEHOUSE_ENABLE=ssh spacedock pi --passthrough -- --mode rpc` requests the `ssh` capability for that session.
```

In `docs/site/contributing/adding-a-runtime.md`, replace `the marker it sets on every launch` with `the marker it sets on every interactive launch`.
Replace `For Pi, spacedock pi launches the front door through the registered Spacedock package` with `For Pi, an ordinary spacedock pi launch uses the registered Spacedock package` (preserve existing inline code).
Add: `Provider mode suppresses these launch markers and explicit resources. See [Provider commands](../reference/command-reference.md#provider-commands).`

### Feedback Cycles


## Stage Report: ideation

- DONE: Reproduce protocol/argv and permission failures for orchestrator launches and demonstrate the riskiest host boundary through existing fixtures or real protocol startup.
  Four seam subtests failed on the current launcher. Real Codex initialize, wrapped Claude init, and Pi get_state established the native baseline.
- DONE: Produce a concrete passthrough design with unchanged value ACs, environment/flag precedence, counted file/net-LOC estimate and tolerance, semantic changes, proof owners, and product-documentation diff.
  AC-1, AC-2, AC-3, AC-4, and AC-5 retain their statements. The design assigns proof, 13 files, +400 net LOC, tolerance, and documentation changes.
- DONE: Commit a complete ideation Stage Report with the three checklist items and risk evidence; send its exact path and commit to the First Officer.
  This entity is the report artifact. The path-scoped state commit and its exact SHA form the delivery to the First Officer.

### Summary

The probes identify interactive launch additions as the cause of protocol and permission failures.
The proposed mode preserves caller arguments and the existing sandbox wrapper, with no new process controller or collaboration policy.
Baseline suite failures limit full regression proof, and the First Officer owns independent design review before the gate.
