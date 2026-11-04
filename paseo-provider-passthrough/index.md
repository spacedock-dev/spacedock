---
title: Preserve orchestrator control in launcher passthrough mode
status: validation
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
        - id: gate:xerfqfkayjawte0gjem7v8v1:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:xerfqfkayjawte0gjem7v8v1-ideation-1
              briefing:
                id: briefing:xerfqfkayjawte0gjem7v8v1:ideation:attempt-1:revision-1
                digest: sha256:3ff2ebd4a21be5e2945a698c39f2288f5a65ba65eb53ebc5c1a2b4c1ae5efc82
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:xerfqfkayjawte0gjem7v8v1:ideation:1
                briefing: briefing:xerfqfkayjawte0gjem7v8v1:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-04T08:51:20.893949Z"
                decision: approve
                reason: Three revision checklist items are DONE and independent re-review approves the design. Real probes establish protocol, permission, and canonical-setting loss. The revision closes R1 through bounded existing allowances, retains seed values and ordinary behavior, and declares focused candidate proof and the full checks/audit still owed.
                conn:
                    quote: Proceed with the 3. You have the conn to approve gate and open pr
                    source: Captain message in this Codex conversation on 2026-11-03
              application:
                target-stage: implementation
                state: consumed
        - id: gate:xerfqfkayjawte0gjem7v8v1:validation
          stage: validation
          attempts:
            - id: gate-attempt:xerfqfkayjawte0gjem7v8v1-validation-1
              briefing:
                id: briefing:xerfqfkayjawte0gjem7v8v1:validation:attempt-1:revision-1
                digest: sha256:d527206f0419b13c145b0eb9dbd7b791d4c6dafc1a1d34d2e656acaeeb54d3ef
                room-ref: '@review/validation/briefing-1'
started: 2026-11-04T07:48:58Z
worktree: .worktrees/spacedock-ensign-paseo-provider-passthrough
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
- A wrapped passthrough launch preserves the selected host's explicit canonical settings through Safehouse's existing named environment allowance.
  Presence selects `--env-pass` for each applicable name, including an empty value.
  The child receives the exact parent value, including `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=0`.
  Absent variables receive no allowance or synthesized value.

| Selected host | Canonical variables preserved through Safehouse in passthrough |
|---|---|
| Claude | `CLAUDE_CONFIG_DIR`, `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` |
| Codex | `CODEX_HOME` |
| Pi | `PI_CODING_AGENT_DIR` |

The launcher adds allowances only for the selected host's present canonical variables.
Caller-provided native Safehouse environment controls retain their existing behavior.
This correction adds no filesystem grant, general environment forwarding, caller setup requirement, or ordinary-mode behavior change.
Safehouse's existing filesystem rules still determine access to the selected settings path.

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
| A present canonical variable for the selected host | In wrapped passthrough only, adds a named allowance and preserves the exact value, including empty or `0`. |
| An absent canonical variable or another host's variable | Adds no automatic allowance or value. Other explicit Safehouse controls retain their behavior. |

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
| Presence-filtered allowance for selected-host canonical settings | AC-2 | Requiring caller-managed allowances adds setup and loses the promised settings by default. The existing bounded allowance supplies the missing transport. |
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

### R1: wrapped caller settings

The [independent review](ideation-review.md) at state commit `68bc1d8574158c35919f69f41272e70b805d3eed` supplies the following accepted observations.
The First Officer classified R1 as owned Material under `value-ac[AC-2]` and authorized the bounded allowance correction.
These observations were not repeated during this revision.

With identical isolated filesystem access, real Safehouse and Codex `config/read` returned two different effective settings:

| Named allowance | Effective model | Effective approval policy |
|---|---|---|
| Absent | `review-default` | `never` |
| Present for `CODEX_HOME` | `review-selected` | `on-request` |

Both runs supplied `CODEX_HOME` for the selected settings directory.
Only the named allowance differed.
The review also observed all four explicit canonical values disappear under the default wrapper and survive exactly with their named allowances.
The explicit Claude teams value remained `0`, and both sentinel observations reported `APP_SANDBOX_CONTAINER_ID=agent-safehouse`.
The common fixture grant held filesystem access constant. It is a probe input, not a proposed product grant.

This establishes environment transport and an independently changed effective Codex permission setting.
It does not establish effective Claude or Pi settings.
Implementation must exercise the candidate through the same real wrapper and native settings boundaries.
The existing protocol startup checks will also observe each host's selected native configuration.
No new mechanism spike is needed: the independent wrapper and Codex settings observations prove the reused allowance.

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
| `internal/cli/frontdoor.go` | Parse the mode, guard interactive additions, and compose the selected canonical allowances into the existing Claude/Codex wrapper. |
| `internal/cli/launch_options.go` (new) | Resolve sandbox environment fallbacks and presence-filtered canonical allowances for the selected host. |
| `internal/cli/pi.go` | Guard interactive readiness/resources/markers, compose the Pi canonical allowance, and preserve path resolution and other wrapper behavior. |
| `internal/cli/help.go` | Describe provider mode, explicit canonical settings, the sandbox environment settings, and the Pi compatibility alias. |
| `internal/cli/frontdoor_parse_test.go` | Exercise the mode/fence grammar and invalid interactive combinations. |
| `internal/cli/frontdoor_test.go` | Exercise literal Claude/Codex provider argv, plugin inactivity, child environment, quiet launch output, and host exit codes. |
| `internal/cli/codex_multi_agent_test.go` | Exercise caller collaboration overrides in passthrough while retaining ordinary collaboration checks. |
| `internal/cli/pi_frontdoor_test.go` | Exercise native provider argv, absent inherited markers, readiness independence, wrapping, and common flags. |
| `internal/cli/launch_parity_test.go` | Exercise sandbox precedence and canonical allowance presence, exact values, host selection, and ordinary-mode stability. |
| `internal/cli/help_test.go` | Extend existing help coverage for the new mode and settings. |
| `docs/site/reference/command-reference.md` | Document provider mode and qualify interactive collaboration/permission behavior. |
| `docs/site/reference/sandbox.md` | Document sandbox precedence and the exact canonical variables preserved in wrapped provider mode. |
| `docs/site/contributing/adding-a-runtime.md` | Qualify the existing Pi marker claim as interactive behavior and link to provider mode. |

R1 uses these same helper, host, environment-matrix, and documentation files.
The estimate remains +400 net LOC across 13 files, with the same insertion/deletion estimate and tolerance.
No product file or standing mechanism is added by this correction. Implementation must count the actual diff.

Permitted semantics: new pre-fence launch grammar, sandbox environment fallback, passthrough ownership of host arguments and child environment, and skipped interactive setup.
Wrapped passthrough preserves only the selected host's present canonical variables through the existing bounded allowance.
This adds no filesystem access or caller setup requirement.
Stored formats, workflow authority, host permission defaults, ordinary launch behavior without these new inputs, and process supervision remain unchanged.

## Acceptance criteria

**AC-1**: Passthrough Codex app-server and Claude/Pi protocol launches receive valid host argv and a clean stdout channel.
Verified by: literal argv cases in the existing launcher fixtures and candidate protocol startup with the reproduction inputs above. Adding a bootstrap prompt or banner to stdout must fail.

**AC-2**: The orchestrator retains agent identity, prompting, and permission settings in wrapped and unwrapped passthrough launches.
Verified by: recorded argv/environment fixtures, caller Codex overrides, candidate Claude init, Pi state, and Codex `config/read` under the real wrapper.
The canonical-variable matrix asserts presence-filtered allowances and exact values, including empty values and Claude's explicit `0`.
Candidate observations must retain the selected native settings and Claude's `permissionMode: default`.
Injecting a first-officer override or omitting a selected canonical allowance must fail, even with exact inner argv.

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
| `launch_parity_test.go`, deterministic | Sandbox environment precedence remains exact. Wrapped passthrough adds only present canonical names for the selected host and preserves value bytes. Omitting an allowance, filtering out empty/`0`, adding an absent value, or changing ordinary allowances breaks this proof. | Seconds. Extend the existing environment matrix. |
| Existing `TestSafehouseEnvPassForwardsSpacedockBin` and Pi stable-path tests, deterministic | A scrubbing wrapper still resolves the launched binary. Pi code/dependencies remain visible under fnm. Removing the env-pass allowance or stable-path grant breaks these owners. | Seconds. Unchanged tests. |
| Existing permission, banner, collaboration, host launch, and PTY tests, deterministic | Ordinary launches retain their identity, defaults, terminal behavior, signals, and exit codes. Broadly disabling their defaults or changing the supervisor breaks existing owners. | Seconds to minutes. Unchanged tests. |
| Existing help tests plus documentation review | Help exposes the supported mode and Pi alias. Examples agree with the host commands. Removing registration or documenting additive capability precedence breaks the check. | Seconds plus review. |
| One-off candidate protocol startup, real local hosts | Protocol startup remains clean. The real wrapper preserves canonical sentinels, and native responses reflect the selected settings. Omitting `CODEX_HOME` changes Codex `config/read` from selected `on-request` to fallback `never`, even with exact argv. | Startup observations only. Reuse the existing plan and independent fixture shape. |

The environment matrix covers unset/empty values, force-on, false, invalid Boolean input, an explicit Boolean override, capability fallback, repeated explicit capabilities, profile-only wrapping, and missing Safehouse.
The canonical-variable matrix covers every selected host, wrapped/unwrapped passthrough, and ordinary mode.
Cases cover present, absent, empty, explicit `0`, and values with spaces or `=`.
Expected variable names and values are literals, independent of production lists.
A different host's variable must not cause an automatic allowance.
Absent Claude teams settings must not become `1`.
These cases extend `launch_parity_test.go`; they require no additional fixture file.
Marker cases begin with both Pi markers already set in the parent.
Plugin cases place a valid checkout beside the selected launcher and record zero discovery/install calls and unchanged isolated host settings.

Candidate observations reuse the independent review's isolated fallback and selected settings directories.
Keep any probe-only fixture grant identical in both observations.
Run the candidate with explicit canonical variables, without caller-supplied environment allowances.
For Codex, send `initialize`, `initialized`, and `config/read` with `{"includeLayers":false}` while stdin remains open.
The selected file contains model `review-selected` and approval policy `on-request`.
The fallback file contains model `review-default` and approval policy `never`.
The candidate must return the selected pair under Safehouse.
For Claude and Pi, use distinct native-supported settings in the selected and fallback directories.
The Claude init and Pi `get_state` responses must reflect the selected configuration.
A real-wrapper sentinel observation must also retain every canonical value exactly, including explicit Claude teams `0`.
All observations retain the existing no-bootstrap and clean-protocol requirements.

Run `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` on the implementation worktree.
Run the full suites sequentially because dispatch fixtures share a fixed artifact directory.
Use the documented runtime-support prompt to resolve harness setup friction before declaring a runtime blocker.
Validation must run the existing detached adversarial audit for launcher changes.
Useful adversarial edits restore a bootstrap prompt, wrapped bypass flag, Pi marker, reversed precedence branch, or omitted canonical allowance.
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
Safehouse preserves the selected host's explicit [canonical settings](sandbox.md#caller-settings-in-provider-mode).
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

### Caller settings in provider mode

With `--passthrough`, Safehouse preserves explicitly supplied values for the selected host:

| Host | Canonical variables |
|---|---|
| Claude | `CLAUDE_CONFIG_DIR`, `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` |
| Codex | `CODEX_HOME` |
| Pi | `PI_CODING_AGENT_DIR` |

Presence selects the allowance, including an empty value.
Values pass unchanged, so an explicit Claude teams value of `0` stays `0`.
Absent variables receive no synthesized value.
Preserving these variables adds no filesystem grant or general environment forwarding.
Existing Safehouse filesystem rules still control access to the selected settings path.
Ordinary launch behavior remains unchanged.
```

In `internal/cli/help.go`, the `--passthrough` description includes this exact behavior:
`Pass host arguments without interactive additions. Preserve explicit canonical host settings through Safehouse.`

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


## Stage Report: ideation-review

- DONE: Challenge necessity and semantic scope of the proposed three-host launcher mode, including permission ownership, Pi resources, Codex configuration, and sandbox precedence.
  `ideation-review.md` records removal challenges and R1. Real wrapped Codex settings changed from the caller's `on-request` to fallback `never`.
- DONE: Review the independent risk evidence, exact value ACs, regression proof owners, counted surface/tolerance, and proposed documentation for missing promised behavior or unnecessary mechanisms.
  The review retains all five seed value criteria, accepts existing probes, checks the +400/13-file estimate, and identifies the missing settings allowance.
- DONE: Commit and send a concise independent ideation-review artifact and Stage Report with three checklist outcomes, concrete finding triage, and approve/revise recommendation; do not change the design.
  `ideation-review.md` and this appended report form the path-scoped state delivery. The first officer receives the exact commit, results, evidence, and revise recommendation.

### Summary

The proposed mode needs its existing guards and wrapper, with no additional process or enforcement mechanism.
The first officer classified R1 as owned Material and directed bounded forwarding of explicit host settings through the existing Safehouse allowlist.
The recommendation is revise, with the original design and value criteria unchanged.


## Stage Report: ideation (cycle 2)

- DONE: Fold the FO-authorized R1 correction into the design: presence-filtered selected-host config forwarding through the existing Safehouse allowlist, with unchanged ordinary mode and filesystem grants.
  The design preserves the four canonical values by presence, including empty and explicit `0`, through the existing selected-host allowance.
- DONE: Update the environment proof plan and concrete documentation diff while preserving all five seed value AC statements, frontmatter, and the declared mechanism/value boundary.
  AC-1, AC-2, AC-3, AC-4, and AC-5 retain their seed statements. The matrix and docs now include canonical transport and native settings observations.
- DONE: Commit and send an exact-stage ideation report of the revision with all three checklist outcomes, the independent review reference, counted surface, and unchanged scope.
  This revision cites `ideation-review.md` at `68bc1d8574158c35919f69f41272e70b805d3eed`. The path-scoped state commit is the First Officer's delivery artifact.

### Summary

R1 now specifies canonical settings preservation through the existing Safehouse allowance, with unchanged ordinary behavior and filesystem grants.
Accepted independent observations prove the reused allowance and Codex permission effect, while implementation still owes candidate observations through those real boundaries.
The estimate remains +400 net LOC across 13 files, with 470 insertions, 70 deletions, and the same tolerance.


## Stage Report: ideation-review (cycle 2)

- DONE: Re-review the revised design against R1 and verify complete presence-filtered selected-host configuration forwarding, without new caller obligations or filesystem grants.
  `ideation-review-r2.md` finds all four variables, presence rules, empty values, explicit Claude `0`, and the authorized host/mode limits complete.
- DONE: Verify the original value AC statements, ordinary behavior boundary, counted surface/tolerance, environment proof matrix, and concrete documentation remain coherent and unchanged in value.
  The five seed value statements remain exact. The existing proof owners cover R1 within +400 net LOC/13 files and unchanged tolerance.
- DONE: Commit and send a second independent ideation-review artifact and report with three checklist outcomes, remaining finding triage, and approve/revise recommendation.
  `ideation-review-r2.md` and this report form the path-scoped state delivery. The first officer receives the exact commit and approve recommendation.

### Summary

The committed revision closes R1 at the design boundary through the existing selected-host Safehouse allowance.
No remaining design finding requires revision, and the recommendation is approve ideation.
Implementation and validation still owe candidate behavior, full regression evidence, and the inherited audit obligations.

## Stage Report: implementation

- DONE: Deliver the approved behavior and documentation with focused red-first evidence for every value acceptance criterion.
  AC-1, AC-2: literal argv cases and six native startup/settings observations pass. AC-3: sandbox matrix and nine real wrapper sentinels pass. AC-4: ordinary launch regressions pass. AC-5: rendered help and provider examples pass.
- DONE: Commit the scoped candidate after formatting and the complete Go and race suites pass; report denominators and unchanged approved semantics.
  AC-1, AC-2, AC-3, AC-4, AC-5: candidate `2f99a397a36ea77addc69e6481f98777d584f440`. Normal `exit 0, 21/21 test packages passed, 2 no-test packages`. Race `exit 0, 21/21 test packages passed, 2 no-test packages (1533 root passes and 1946 subtest passes)`. Both commands use the assigned environment and a serialized window.
- DONE: Commit and send a complete implementation Stage Report with candidate identity, three checklist outcomes, AC evidence, surface actuals, and any finding proposals.
  AC-1, AC-2, AC-3, AC-4, AC-5: `status --read` returns three DONE items and citations for all five criteria. The state commit and canonical path form the delivery.

### Evidence

- AC-1: `TestProviderHostLaunch`, `TestPiProviderLaunch`, and six native startup observations retain caller argv and JSON stdout. A bootstrap prompt or banner breaks these owners.
- AC-2: `TestCodexProviderAcceptsCallerCollaboration` and `TestProviderCanonicalSettings` retain overrides and exact present values. Restoring a reserved-key rejection or dropping an empty/0 allowance breaks them.
- AC-2 native: wrapped/plain Codex `config/read` returns `review-selected`/`on-request`. Claude init returns `claude-sonnet-5`/`permissionMode: default`. Pi returns `autoCompactionEnabled: false`/`steeringMode: all`.
- AC-3: `TestSandboxEnvironmentPrecedence` covers three hosts and both modes. Real wrapper sentinels retain canonical settings, sandbox identity, binary identity, and exit 7. Reversing precedence or removing an allowance fails.
- AC-4: `TestProviderGrammar` and Pi readiness cases cover the common alias and invalid input. Existing permissions/collaboration and 9 process/Pi path roots plus 11 subtests retain ordinary behavior. Removing defaults or fnm grants fails.
- AC-5: `TestProviderHelp` and existing help owners exercise rendered command help. Native startup exercises provider examples. Documentation records precedence and bounded settings transport.
- Focused red: `/tmp/paseo-provider-passthrough-red.jsonl` has 7 failed roots, 95 failed subtests, and 47 passed subtests before product edits. Grammar-only follow-up records active gate/override/readiness failures.
- Focused green: 7/7 roots and 142/142 subtests pass. Logs: `/tmp/paseo-provider-passthrough-focused.jsonl`, `paseo-provider-passthrough-process.jsonl`, and `paseo-provider-passthrough-help-correction.jsonl`.
- Race skips: `TestCodexIsolatedHomeCollaborationLifecycle` needs the opt-in live flag. `TestUpdateFrontmatterBlockScalarsRewrapped` has documented divergence #L2-b and no live fixture. Neither counts as a passing test. Both remain unverified.
- Full commands: `/tmp/paseo-provider-passthrough-normal-retry.log` records `go test ./...` and `/tmp/paseo-provider-passthrough-race.jsonl` records `go test ./... -race -timeout 30m -json`, sequential with `GOFLAGS=-p=1`, `GOCACHE=/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/paseo-fo-64ye7hx6/gocache`, and PATH prefixed by `/Users/clkao/.local/share/fnm/node-versions/v24.13.1/installation/bin`.
- Native proof: 15/15 observations in `/tmp/paseo-provider-passthrough-proof/native-probes.json`. 9 wrapped sentinel cases plus 6 wrapped/plain native settings startups. Every consumed stdout line parses as JSON.
- Native commands: candidate `host --passthrough [--safehouse-add-dirs=<fixture>] -- <native argv>`. Codex uses `app-server` with initialize/initialized/config-read. Claude uses the specified stream-json print flags. Pi uses RPC/offline/no-session/no-extensions/no-skills.
- Probe setup: native versions are Codex 0.166.0, Claude 2.1.283, Pi 1.0.2, Safehouse 0.12.0. The isolated Claude/Pi wrapper also receives an explicit read-only fnm-installation grant. No product grant changes.
- Probe setup corrections: isolated-HOME Claude initially exits 127 without the fixture-only grant. Pi 1.0.2 rejects `--no-mcp`. The isolated probe has no MCP config. Model-less Pi clamps thinking to off, so documented compaction/steering settings distinguish the native settings files.
- Fixture correction: executable expectations use independent filesystem symlink resolution for `/var` versus `/private/var`. `SPACEDOCK_BIN` still promises the resolved selected executable.
- Formatting: `gofmt -w ./cmd ./internal` completed. Unrelated pre-existing alignment drift in `internal/release/runtime_live_evidence_workflow_test.go` was restored. `git diff --check` passes.
- Surface: `13 files, +612/-240 = +372 net LOC (93% of estimate)` versus +400 net LOC/13 files (+470/-70). Tolerance remains +280..+520 net and at most 15 files. Approved ACs and semantic boundaries remain unchanged. Ordinary permission defaults, stored formats, workflow authority, and process supervision remain unchanged.

### Finding disposition

- Exact observation: `TestPiHelpCarriesSafehouseDetail` reports `pi --help should not declare "--skip-compat-check"`. `/tmp/paseo-provider-passthrough-regression.jsonl` records 66 passed outcomes and this one failure.
- Released user and normal workflow: a developer runs the existing Pi help regression before shipping the approved common flag.
- Observable harm: the stale negative oracle rejects the supported compatibility alias and blocks required suite completion.
- Authority: `value-ac[AC-4]` Every host accepts `--skip-compat-check`, while ordinary Pi readiness stays active.
- Trigger evidence: the HEAD oracle excludes both flags. Candidate help advertises the new alias, and `TestPiProviderLaunch` retains ordinary readiness.
- Proposed materiality: Material. Proposed ownership: this task's approved `pi_frontdoor_test.go`. Proposed disposition: fix only the stale alias expectation and comment.
- FO authorization: `fix` that expectation/comment, retain the `--no-install` exclusion and readiness assertions, and rerun the focused owner. Both focused roots pass after that authorized correction.

- ENOSPC observation: the first normal aggregate exits 1 at `safehouse.test: open .../go-build3647562817/b268/importcfg: no space left on device`. It records 16 passing test packages, 5 build failures, and 2 no-test packages (23 total).
- Released user/workflow: a developer runs the required normal suite before candidate commit. Observable harm: the first run leaves five test packages unbuilt, so its proof is incomplete.
- Authority: `contract[AGENTS.md#expected-commands]` The complete normal and race checks remain required before completion.
- Trigger evidence: six ENOSPC messages in `/tmp/paseo-provider-passthrough-normal.log`, automatic Go work-directory removal, and 421 MiB available afterward.
- Proposed materiality: Material evidence/environment defect. Proposed ownership: FO-session environment, outside product ownership. Proposed disposition: hold the candidate and request bounded recovery without an assertion waiver.
- FO authorization: retry the five packages individually with the same environment and no `-count`. Then repeat the complete normal command and race with 30m headroom. Stop on any recurring ENOSPC. All five retries pass. The complete normal retry and complete race each pass 21/21 test packages without another ENOSPC.

### Summary

The candidate preserves orchestrator control through the existing launcher and Safehouse wrapper, with no new protocol parser or process controller.
Native startup/settings and deterministic proof pass. Validation still owns the detached adversarial audit and required live lanes before merge.

## Stage Report: validation

- DONE: Independently verify each acceptance criterion against the exact committed candidate and report explicit PASS/FAIL with reproducible external evidence.
  AC-1 PASS: exact argv and clean stdout. AC-2 PASS: caller settings and identity. AC-3 PASS: wrapping, precedence, and binary identity. AC-4 PASS: provider bypass and ordinary defaults. AC-5 PASS: rendered help and runnable examples. Evidence follows.
- DONE: Attack the highest-risk semantics and perform any workflow-required detached adversarial audit in a throwaway checkout; report concrete falsifiers and outcomes.
  AC-1, AC-2, AC-3, AC-4: 18/18 independent production mutations failed their existing test owners; 21/21 external process observations passed. Candidate bytes and HEAD stayed unchanged.
- DONE: Commit and send a complete validation Stage Report with the three checklist outcomes, tests and denominators, exact candidate SHA, authorized findings, and PASSED/REJECTED recommendation.
  AC-1, AC-2, AC-3, AC-4, AC-5: this canonical entity contains the report. The path-scoped state commit and exact candidate SHA accompany the delivery to the First Officer. Recommendation: PASSED for PR readiness.

### Evidence

- Candidate: `2f99a397a36ea77addc69e6481f98777d584f440`, branch `paseo-provider-passthrough`. The implementation worktree and detached audit checkout both remain clean at this SHA.
- AC-1 PASS: `TestProviderHostLaunch` and `TestPiProviderLaunch` compare literal complete argv and quiet launcher channels. Six external plain/wrapped launches retained 14 arguments, including empty, Unicode, newline, repeated settings, and a native fence. Adding an identity, prompt, permission flag, or banner failed the audit.
- AC-2 PASS: `TestCodexProviderAcceptsCallerCollaboration` accepts caller overrides. `TestProviderCanonicalSettings` checks selected-host allowances and exact absent/empty/0/spaced values. Native evidence returns Codex `review-selected`/`on-request`, Claude `claude-sonnet-5`/`permissionMode: default`, and Pi `autoCompactionEnabled: false`/`steeringMode: all`, wrapped and plain. Removing a canonical allowance, filtering empty/0, changing host selection, or retaining Pi markers failed.
- AC-3 PASS: `TestSandboxEnvironmentPrecedence`, `TestSafehouseEnvPassForwardsSpacedockBin`, and `TestRunPi_LaunchArgv0_StableUnderWrapWithAddDirs` passed independently. Real-wrapper processes observed `APP_SANDBOX_CONTAINER_ID=agent-safehouse`, the exact selected binary, and exit 7. Ignoring explicit Boolean/capability flags or removing the fnm grant failed.
- AC-4 PASS: provider fixtures skip discovery, installation, compatibility, and Pi readiness. `TestPiProviderLaunch` retains ordinary readiness despite the alias. Existing Claude default, Codex complete collaboration argv, and Pi native-resource owners passed. Removing ordinary Claude permissions or Codex collaboration failed the audit.
- AC-5 PASS: `TestProviderHelp` exercised all three rendered command helps. The concrete documentation describes pre-fence mode selection, native forwarding, environment precedence, and bounded canonical settings. Producer native startup exercised Claude stream-json, Codex app-server, and Pi RPC examples. Changed help, comments, and documentation passed the Simple English review.
- Independent focused command: `go test ./internal/cli -run '^(TestProvider.*|TestCodexProviderAcceptsCallerCollaboration|TestPiProviderLaunch|TestSandboxEnvironmentPrecedence|TestCodexCollaborationLayerCompleteArgv|TestClaudeFrontDoorLaunchesOnCompatible|TestSafehouseEnvPassForwardsSpacedockBin|TestRunPi_LaunchArgv0_StableUnderWrapWithAddDirs|TestPiFrontDoorLaunchesWithNativeResourcePaths)$' -count=1 -json`; 12/12 roots, 148/148 subtests, one package passed, no skips. Log: `/tmp/paseo-provider-passthrough-validation-baseline.jsonl`.
- Required broad checks reused: `go test ./...` passed 21/21 test-bearing packages plus 2 no-test packages; `go test ./... -race -timeout 30m -json` passed the same packages, 1533 roots, and 1946 subtests. Logs: `/tmp/paseo-provider-passthrough-normal-retry.log` and `/tmp/paseo-provider-passthrough-race.jsonl`.
- Unverified skips remain: `TestCodexIsolatedHomeCollaborationLifecycle` requires an opt-in model turn. `TestUpdateFrontmatterBlockScalarsRewrapped` retains documented divergence #L2-b. Neither is counted as passed.
- Environment: every independent Go command used `GOCACHE=/var/folders/h1/vnssm1dj6ks4nzzvx8y29yjm0000gn/T/paseo-fo-64ye7hx6/gocache`, `GOFLAGS=-p=1`, and PATH prefixed with `/Users/clkao/.local/share/fnm/node-versions/v24.13.1/installation/bin`. No full suite ran concurrently; no cache or peer files were cleaned. No independent command reported ENOSPC; final free space was 742 MiB.
- Native evidence reused: `/tmp/paseo-provider-passthrough-proof/native-probes.json` records 15/15 real observations: nine Safehouse sentinels and six native settings startups. Its final Pi observations supersede the documented initial fixture setup failure.
- External independent proof: `/tmp/paseo-provider-passthrough-validation/process-results.json` records 21/21 observations: six exact plain/wrapped argv/environment launches, nine early task/plugin rejections, three invalid-Boolean rejections, and three explicit-Boolean overrides. Rejections added zero child spy events. The exact-SHA binary has SHA-256 `291923bee44eb03d4bab9fbfe71425337b765f9b24f1bd87ced2e269464a8675`.
- Detached audit: `/tmp/paseo-provider-passthrough-validation-audit` was created with `git worktree add --detach ... 2f99a397a36ea77addc69e6481f98777d584f440`. `/tmp/paseo-provider-passthrough-validation/audit.py` applies each production-only edit separately, runs `go test ./internal/cli -run '^OWNER$' -count=1 -json`, and restores the original bytes.
- Audit AC-1/AC-2: bootstrap, Claude identity, both wrapped permission bypasses, and stdout banner failed `TestProviderHostLaunch`; restored collaboration rejection failed `TestCodexProviderAcceptsCallerCollaboration`; inherited Pi markers and unsolicited resources failed `TestPiProviderLaunch`.
- Audit AC-2/AC-4: omitted Codex allowance, filtered empty values, filtered 0 values, wrong-host selection, and expanded ordinary allowances failed `TestProviderCanonicalSettings`. Expected values remain test literals independent of production lists.
- Audit AC-3/AC-4: ignored explicit Boolean/capability flags failed `TestSandboxEnvironmentPrecedence`; removed ordinary Codex collaboration failed `TestCodexCollaborationLayerCompleteArgv`; removed ordinary Claude permissions failed `TestClaudeFrontDoorLaunchesOnCompatible`; removed fnm grant failed `TestRunPi_LaunchArgv0_StableUnderWrapWithAddDirs`.
- Audit total: 18/18 mutations caught by assertion failures, zero surviving mutations, zero compile-failure substitutes. Full per-mutation commands, exit codes, failed tests, and logs are in `/tmp/paseo-provider-passthrough-validation/audit-results.json`.
- Format/surface: producer ran `gofmt -w ./cmd ./internal`; independent `gofmt -l` found no changed Go file. `git diff --check main...HEAD` passed. Actual surface is 13 files, +612/-240, net +372 (93% of +400 estimate), within +280..+520 and maximum 15 files. All five approved ACs and semantic boundaries remain unchanged.

### Findings and remaining proof

- No new outcome defect, evidence defect, deferred risk, or polish finding arose. The detached audit refuted nothing material; no candidate edit or reviewer rerun required disposition.
- Existing authorized corrections are satisfied: the Pi help oracle accepts the compatibility alias while retaining the unsupported `--no-install` exclusion and ordinary readiness. The authorized ENOSPC recovery produced complete green normal/race runs.
- Before merge, Claude, Codex, and Pi live CI still require green evidence for this exact product SHA. The First Officer owns this requirement and holds merge. No provider live-lane evidence exists yet; unrelated reader PR #860 cannot establish it.
- The First Officer reports reader-only run `41817979627` has Pi/Codex OAuth failures before agent turns (`refresh_token_reused`/expired token, HTTP 401). These are environment evidence failures on another candidate. No credential mutation or repeated doomed run was authorized.

### Summary

The candidate satisfies AC-1 through AC-5 with deterministic, native, and independent process evidence. All 18 adversarial edits were caught, and the recommendation is PASSED for PR readiness.
Required host live CI remains pending before merge, and both named skipped checks remain unverified.
