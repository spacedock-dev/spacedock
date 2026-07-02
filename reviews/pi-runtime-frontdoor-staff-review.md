## Review

- Correct:
  - `spacedock pi` is registered and routes through Pi-native extension/skill paths, not Claude/Codex plugin/team tools. Evidence: `internal/cli/cli.go:196-219`, `internal/cli/pi.go:68-76`, `internal/cli/pi_frontdoor_test.go:53-97`.
  - Fixback commit `3bcaf487` is narrow and appropriate: it only updates the live smoke prompt/test for the exact required heading in `internal/ensigncycle/pi_live_runner_test.go`.
  - Cheap package tests passed locally: `go test ./internal/cli ./internal/ensigncycle -count=1`.

- Material:
  - None found.

- Important:
  - `spacedock pi` and `install --host pi` treat the runtime as ready even when Pi auth is missing. `checkPiRuntime` records auth at `internal/cli/pi.go:247-250`, but `piRuntimeLaunchReady` excludes `authOK` at `internal/cli/pi.go:258-259`. `runPi` gates launch on that incomplete readiness at `internal/cli/pi.go:61-66`, and `runInitWithPi` prints `Pi runtime ready` on the same incomplete predicate at `internal/cli/pi.go:93-96`. The test fixture also encodes this gap by asserting ready without adding auth to `statOK` at `internal/cli/pi_frontdoor_test.go:108-111`. This conflicts with the path/auth handling goal and the install instructions that say authentication is required.
  - `--plugin-dir` is now accepted and silently stripped for non-Pi install/doctor paths. `parsePiSetupArgs` accepts `--plugin-dir` before host dispatch, then non-Pi paths call `stripPluginDirArg` at `internal/cli/pi.go:89-90`. This changes existing Claude/Codex behavior from “unknown argument” to “ignored”, which weakens compatibility and can hide operator mistakes such as `spacedock install --plugin-dir ./checkout`.

- Polish:
  - Top-level help still describes `doctor` as checking “the installed plugin” for all hosts, including Pi, even though Pi intentionally has no plugin marketplace path. See `internal/cli/help.go:16-22`. Consider host-neutral wording.

Recommendation: not ready to merge as-is. The fixback itself is sound, but I would require the auth readiness and non-Pi `--plugin-dir` compatibility issues to be fixed or explicitly accepted before merging. I did not write the requested review file because the task also specified read-only/no modifications.