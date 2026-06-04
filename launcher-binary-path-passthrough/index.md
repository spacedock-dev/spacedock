---
id: fcfe02dk6nbwk9pw896d3eff
title: Front-door launcher doesn't propagate its own binary path ($0) to the launched claude/codex session — the in-session skill falls back to $PATH spacedock instead of the launching binary
status: validation
source: captain (2026-06-03) — `/path/to/spacedock claude …` does not preserve `$0` for the skill to know which spacedock launched it; need a way to pass env from the spacedock CLI over (safehouse permitting) to claude/codex, and have the skill use the env or default to $PATH
score: "0.28"
worktree: .worktrees/spacedock-ensign-launcher-binary-path-passthrough
started: 2026-06-04T00:00:00Z
completed:
verdict:
issue:
---

When `spacedock claude` or `spacedock codex` launches a first-officer session, the session's operating contract later shells out to `spacedock` for status and dispatch helpers (`spacedock status`, `spacedock dispatch build`, `spacedock dispatch reconcile`, and related checks). Those calls currently resolve through the session's `$PATH`, not through the binary that performed the front-door launch.

That breaks the common development workflow: build a local launcher, run it explicitly, and expect the whole session to exercise that same build. For example:

```sh
/Users/me/src/spacedock/spacedock claude
./spacedock codex --safehouse
```

In both cases, the child Claude/Codex process can load the Spacedock skill but the skill text still says to run bare `spacedock`. If `$PATH` contains an installed release, the in-session FO or dispatched ensigns silently use that installed release instead of the explicit launcher. A stale installed binary can have an older contract, missing dispatch helpers, or different state mutation behavior. The failure mode is especially bad for dispatched ensigns because they do not know how the FO was launched and should not be expected to manually substitute `./spacedock`.

Concrete prior instance: a first-officer session launched while `$PATH` resolved `spacedock` to 0.19.1 even though the repository build was newer; the FO had to invoke `./spacedock` manually throughout. This entity addresses the binary-identity half of the related dev-mode friction where sessions may also load the installed plugin contract rather than the repository's plugin contract.

## Design goal

Make front-door launches compatibility-first and binary-stable:

- If an operator launches through an explicit binary path, the launched session should prefer that same binary for all Spacedock helper calls.
- If no launch-provided binary signal exists, existing behavior stays unchanged: use `spacedock` from `$PATH`.
- The mechanism must work for both Claude and Codex front doors and must be observable in tests without requiring a real interactive agent session.
- The mechanism must degrade safely when a runtime wrapper, especially safehouse, does not propagate the signal.

## Proposed compatibility-first approach

### 1. Add a launch-scoped environment signal

Use a single environment variable, `SPACEDOCK_BIN`, as the launcher-to-session signal. `spacedock claude` and `spacedock codex` should set it in the child process environment before execing the host runtime.

Resolution policy for the value:

1. Start from the launcher process identity, not from `$PATH` inside the child. In production, this should be based on `os.Executable()` rather than a later lookup of `spacedock`.
2. Convert to an absolute path. If the operator invoked `./spacedock`, the child receives `/abs/path/to/spacedock` so subsequent working-directory changes do not break the helper command.
3. Resolve symlinks when feasible (`filepath.EvalSymlinks`) so a symlinked development binary points at the actual artifact. If symlink resolution fails but the absolute path exists and is executable, fall back to the absolute path rather than failing the launch.
4. For unusual launch forms where no stable executable path exists (notably `go run`, deleted/replaced temp executables, or a platform error from `os.Executable()`), omit `SPACEDOCK_BIN` and keep the existing `$PATH` behavior. Do not block launching solely because this enhancement cannot identify the binary.

The launcher should set `SPACEDOCK_BIN` for the child environment even when it was already present in the parent shell. A stale parent value from a previous session must not override the front door that is launching this session. The only accepted override in this design is the explicit binary path used to invoke the front door itself; a separate user-facing override flag is out of scope for the first implementation.

### 2. Extend the launch seam so tests can observe argv and env

The current front-door seam records argv (`hostOps.Launch(argv []string)`) and production `execHost.Launch` passes `os.Environ()` to `syscall.Exec`. The implementation should make the child environment a first-class launch artifact, for example by changing/adding a seam that accepts the final env alongside argv.

Required property: front-door unit tests can assert, without starting Claude or Codex, that:

- both `spacedock claude` and `spacedock codex` launch with `SPACEDOCK_BIN=/resolved/path/to/spacedock`;
- existing argv behavior is unchanged;
- an existing parent `SPACEDOCK_BIN=/old/path` is replaced for the child; and
- when binary resolution fails, the env does not contain a misleading value and the launch still proceeds.

Keep this narrow: do not introduce a general environment-passthrough framework unless it is the smallest way to make the launch seam testable.

### 3. Preserve the signal through safehouse launches

Safehouse wrapping changes the outer executable from the host runtime to `safehouse`, for example:

```text
safehouse --trust-workdir-config ... -- claude ...
safehouse --trust-workdir-config ... -- codex ...
```

`SPACEDOCK_BIN` must be present in the environment of the process that execs safehouse, and the design must verify whether safehouse passes that variable through to the inner host. If safehouse filters environment variables by default, one of these compatibility-first strategies must be used:

- prefer an existing safehouse allowlist/passthrough mechanism that lets the launcher explicitly allow `SPACEDOCK_BIN`; or
- if safehouse cannot pass the variable, document the degraded behavior and keep `$PATH` fallback so safehouse users are not broken.

The first implementation should not modify safehouse itself. It may only use safehouse's documented/current passthrough knobs if available from the launcher. The riskiest unknown is runtime observability: a safehouse-wrapped smoke should prove whether `SPACEDOCK_BIN` survives into a trivial child command or into a live host-visible environment.

### 4. Teach the skill/contract to resolve env-then-PATH

Skill instructions should stop hard-coding only the token `spacedock` as the conceptual command source. They should define a helper convention:

```sh
${SPACEDOCK_BIN:-spacedock} status ...
${SPACEDOCK_BIN:-spacedock} dispatch build ...
```

Prefer a centralized text/invariant approach over ad hoc edits at every command site:

- Shared FO/ensign/debrief contract text should introduce `SPACEDOCK_BIN` once as the launcher-provided command path.
- Canonical command examples should either use the variable form directly or state an invariant that every bare `spacedock` command is shorthand for `${SPACEDOCK_BIN:-spacedock}`.
- Generated dispatch prompts that include fetch commands should emit the env-aware form or be backed by a helper that renders it consistently.

A code-level command renderer is preferable wherever commands are generated by the binary, because it is testable and avoids prose drift. Purely manual skill text changes are acceptable only for the static instructions that are not generated from Go code.

### 5. Keep fallback behavior compatible

When `SPACEDOCK_BIN` is unset, empty, or points to a non-executable path, the contract should fall back to `spacedock` on `$PATH` rather than aborting immediately. The version gate (`spacedock --version` / contract range parse) remains the authoritative failure point if the fallback is absent or incompatible.

If `SPACEDOCK_BIN` is set but stale (path no longer exists after an upgrade, clean, or temp `go run` cleanup), the session should surface the stale path in the failure message and then either:

- fall back to `$PATH` and still run the normal contract version gate; or
- fail with an actionable message if executing the stale path was explicitly requested by a future override.

For this first compatibility task, prefer fallback for stale launcher-provided env because it preserves current sessions under wrappers that may mutate or cache environment. The failure must be observable and not silent: tests should cover the resolver's decision and the contract/version gate should reveal the actual binary used.

## Claude and Codex front-door considerations

- Claude and Codex share the same binary identity problem because both front doors launch an interactive host with Spacedock skill instructions that later shell out.
- Claude has `--agent spacedock:first-officer`; Codex currently relies on the bootstrap prompt naming the FO skill. The env signal should be independent of this difference.
- Safehouse flags differ in inner host argv (`--dangerously-skip-permissions` for Claude, `--dangerously-bypass-approvals-and-sandbox` for Codex), but the `SPACEDOCK_BIN` env should be attached at the outer exec layer in both cases.
- Resume paths (`claude --resume`, `codex resume`) should still receive `SPACEDOCK_BIN`; a resumed session may execute new helper commands even though the bootstrap prompt is suppressed.

## Risks and edge cases

- **Relative launch paths:** `./spacedock claude` must not become unusable after the FO changes directories. Resolve to absolute before exporting.
- **Symlinks:** resolving symlinks improves binary identity but may surprise operators who expected the symlink path. Acceptance should assert the chosen policy, not leave it implicit.
- **`go run`:** `os.Executable()` may point to a temporary build artifact that disappears. Treat as best-effort; do not promise stable `SPACEDOCK_BIN` for `go run` unless a later runtime task explicitly supports it.
- **Deleted/replaced binary after launch:** a long-lived session may keep a stale path. The env-aware resolver should fail loudly or fall back with version-gate proof, not silently run an unrelated binary.
- **Safehouse env filtering:** if safehouse strips `SPACEDOCK_BIN`, safehouse launches degrade to `$PATH`. This must be a known, tested behavior rather than an assumption.
- **Claude vs Codex runtime differences:** the front-door argv shape differs, but env injection should be shared and tested for both hosts.
- **Explicit override:** a future `--spacedock-bin` flag or `SPACEDOCK_BIN` parent override could be useful for advanced debugging, but adding override UX in this task risks stale-env ambiguity. Keep parent env overwritten by the actual launcher for now.
- **Security/trust:** `SPACEDOCK_BIN` is a command path trusted by the launched session. It should come from the launcher itself, not from untrusted workflow state or plugin text.
- **Instruction drift:** some static skill docs and generated dispatch prompts may continue to show bare `spacedock`. Add invariant tests to prevent mixed command forms from reappearing.

## Acceptance criteria

**AC-1 — Front doors inject the launching binary path into child env.**  
External/failable proof: Go tests on the front-door launch seam assert that both `spacedock claude` and `spacedock codex` pass an env containing `SPACEDOCK_BIN` set to the resolved absolute launcher path while preserving the existing argv for plain and safehouse-wrapped launches.

**AC-2 — Parent/stale env does not override the current launcher.**  
External/failable proof: a unit test starts the launcher seam with parent env containing `SPACEDOCK_BIN=/old/spacedock` and asserts the child env contains the current resolved launcher path instead. A separate failure-path test asserts that when launcher path resolution fails, no stale or partial `SPACEDOCK_BIN` is injected.

**AC-3 — In-session command resolution prefers `SPACEDOCK_BIN` and falls back to `$PATH`.**  
External/failable proof: resolver tests or generated-command golden tests show helper invocations use the env-aware command when `SPACEDOCK_BIN` is set, and default to `spacedock` when unset/empty/non-executable. The tests must inspect actual rendered command strings or resolver outputs, not only prose.

**AC-4 — Skill instructions and generated prompts carry a single env-aware invariant.**  
External/failable proof: skill smoke/invariant tests scan the FO, ensign, debrief, and generated dispatch surfaces and fail if canonical Spacedock helper instructions regress to unqualified bare `spacedock` without the documented `${SPACEDOCK_BIN:-spacedock}` convention.

**AC-5 — Safehouse behavior is runtime-observable.**  
External/failable proof: a safehouse fixture or live smoke launches a trivial environment-printing child through the same safehouse wrap shape and records whether `SPACEDOCK_BIN` reaches the inner process. If safehouse strips it, the smoke must prove the session falls back to `$PATH` cleanly and the implementation notes the limitation.

**AC-6 — Existing launches remain compatible when the signal is absent.**  
External/failable proof: tests run the resolver/skill command path with no `SPACEDOCK_BIN` and demonstrate the same `spacedock` command token and version-gate behavior as today.

## Test plan

1. **CLI launch env unit tests**
   - Add/focus front-door tests around `runClaude` and `runCodex` using fake host ops that record argv and env.
   - Cover plain launch, safehouse-wrapped launch, resume launch, parent env overwrite, and binary-resolution failure fallback.
   - Assert no argv regressions while adding env assertions.

2. **Binary path resolver tests**
   - Table-test absolute path, relative path, symlink path, missing executable, `os.Executable` error injection, and parent `SPACEDOCK_BIN` overwrite behavior.
   - Keep platform-specific assumptions minimal; skip symlink tests on platforms/filesystems where symlinks are unavailable.

3. **Skill instruction invariant tests**
   - Extend existing skill text/integration tests so canonical FO/ensign/debrief command instructions define and use the env-aware launcher command convention.
   - Include generated dispatch prompt tests for fetch commands emitted by `spacedock dispatch build`, because ensigns rely on those strings without manual FO correction.

4. **Safehouse smoke/fixture**
   - If safehouse is available in CI or a fixture can stand in, run the same wrap shape with an env-printing child and assert the fate of `SPACEDOCK_BIN`.
   - If live safehouse cannot run in baseline CI, keep a focused manual/live smoke documented with command output and make the unit fallback behavior mandatory.

5. **Compatibility regression tests**
   - Run the normal baseline (`go test ./...`) and front-door-focused tests to prove existing `spacedock claude`/`spacedock codex` argv, contract gate behavior, and fallback `$PATH` behavior are unchanged.

## Out of scope

- Adding PR/mod behavior or changing workflow state semantics.
- Changing safehouse internals or creating a broad safehouse environment-management feature.
- Solving plugin source identity (`--plugin-dir` / installed plugin contract vs repository plugin contract), except to ensure this binary-path signal does not conflict with that work.
- Adding a user-facing binary override flag in the first implementation.
- Promising stable behavior for `go run`-launched sessions beyond best-effort omission/fallback.
- Rewriting all skill command prose by hand if a generated/helper-level invariant can enforce the behavior more narrowly.

## Stage Report: ideation

DONE — Clarified why explicit-path front-door launches can silently switch to a different `$PATH` binary inside Claude/Codex sessions; proposed a compatibility-first `SPACEDOCK_BIN` launch env with env-then-PATH skill/command resolution; identified safehouse, symlink, relative path, `go run`, stale env, and host-runtime edge cases; defined external/failable acceptance criteria and a test plan covering launch env, resolver behavior, skill invariants, generated prompts, and safehouse smoke validation.

## Stage Report: implementation

DONE — Filed implementation result for product commit `ede26b09` (`launch: propagate spacedock binary path`) on branch `spacedock-ensign/launcher-binary-path-passthrough` in `.worktrees/spacedock-ensign-launcher-binary-path-passthrough`.

DONE — AC-1 / AC-2: the implementation changes the front-door launch seam from `Launch(argv)` to `Launch(argv, env)` and injects launch-scoped `SPACEDOCK_BIN` for both Claude and Codex from the current resolved launcher executable. Tests cover resolved absolute/symlink paths, safehouse/resume launch argv preservation, stale parent `SPACEDOCK_BIN` replacement, and resolver-failure omission while launch continues.

DONE — AC-3 / AC-4: generated dispatch fetch commands use `${SPACEDOCK_BIN:-spacedock}` through a shared renderer, golden dispatch outputs were updated, and FO/ensign/debrief skill text now documents the env-aware launcher invariant. Skill integration tests lock the invariant.

DONE — AC-5: front-door tests prove `SPACEDOCK_BIN` is present in the outer safehouse exec environment. Safehouse internals were not changed; `docs/runtime-support.md` documents fallback/degraded behavior if a wrapper strips the variable before the inner runtime observes it.

DONE — AC-6: unset or unresolvable launcher signal preserves launch behavior, and skill/dispatch command rendering falls back to `spacedock` via `${SPACEDOCK_BIN:-spacedock}`.

DONE — Changed files in product commit: `docs/runtime-support.md`; `internal/cli/frontdoor.go`; `internal/cli/frontdoor_test.go`; `internal/cli/host_exec.go`; `internal/cli/pi.go`; `internal/dispatch/build.go`; `internal/dispatch/build_hazards_test.go`; `internal/dispatch/launcher_command.go`; `internal/dispatch/launcher_command_test.go`; `internal/dispatch/native_subcommands_routing_test.go`; dispatch golden files under `internal/dispatch/testdata/golden/`; `internal/ensigncycle/cycle_test.go`; `skills/debrief/SKILL.md`; `skills/ensign/references/ensign-shared-core.md`; `skills/first-officer/references/first-officer-shared-core.md`; `skills/integration/skill_text_test.go`.

DONE — Validation reported by the implementation worker: `gofmt -w ./cmd ./internal`; focused front-door env tests with `go test ./internal/cli -run 'Test(ClaudeFrontDoorInjectsResolvedLauncherBin|ClaudeFrontDoorOmitsStaleLauncherBinWhenResolutionFails|ClaudeFrontDoorLaunchEnvResolvesSymlink|CodexFrontDoorInjectsLauncherBinThroughSafehouseResume|ClaudeFrontDoorLaunchesOnCompatible|CodexFrontDoorLaunchesOnCompatible)' -count=1`; focused dispatch/skill tests with `go test ./internal/dispatch -run 'Test(Build|LauncherCommand)' -count=1` and `go test ./skills/integration -count=1`; full baseline `go test ./... -count=1`; race baseline `go test ./... -race -count=1`.

SKIPPED — No product/source changes were made during this fixback because commit `ede26b09` already contained an implementation that could be honestly reported.

SKIPPED — No live safehouse inner-environment smoke was added in this fixback; implementation evidence remains outer safehouse exec env coverage plus documented fallback if safehouse strips `SPACEDOCK_BIN`.

FAILED — None for the implementation report filing.

Residual risks: no live safehouse inner-env proof is recorded; generated shell fetch commands cover unset/empty `SPACEDOCK_BIN` directly while stale/non-executable launcher-path behavior relies on the documented contract/version-gate fallback guidance.
