## Review

I did not write `/Users/clkao/git/spacedock-research/spacedock-v1/reviews/pi-runtime-support-staff-review.md` because the task explicitly says read-only / do not modify files.

- Correct:
  - Pi dispatch is cleanly contained in the existing dispatch builder host-mode seam, and preserves compatibility defaults (`host` defaults to `claude`; only `claude|codex|pi` accepted) in `internal/dispatch/build.go:103-109`.
  - Pi prompt/completion wording avoids Claude team-tool calls and uses dispatch-file + Pi subagent completion semantics in `internal/dispatch/build.go:410-435`, `internal/dispatch/build.go:452-465`, and `internal/dispatch/build.go:484-490`.
  - Offline tests cover banned Claude tool syntax and split-root entity path preservation in `internal/dispatch/build_pi_host_test.go:14-73` and `internal/dispatch/build_pi_host_test.go:75-124`.
  - `internal/piruntime` is small and appropriately separated from CLI/dispatch. The teams adapter maps lifecycle concepts to native `teams` actions, not Claude emulation, in `internal/piruntime/teams.go:24-46`.
  - Stale-completion protection is represented with epoch-scoped evidence and tested in `internal/piruntime/registry.go:64-87` and `internal/piruntime/registry_test.go:10-46`.
  - Live smoke harness is gated behind `//go:build live`, uses temp Pi config/session directories, copies auth into isolated Pi home, and asserts durable state/git evidence rather than transcript prose in `internal/ensigncycle/pi_live_runner_test.go:1-128`, `:200-230`.

- Material:
  - None found.

- Important:
  - `spacedock install --host pi` reports “Pi runtime ready” even when Pi auth is missing. `runInitWithPi` checks `piRuntimeLaunchReady` only (`internal/cli/pi.go:92-96`), and `piRuntimeLaunchReady` excludes `authOK` (`internal/cli/pi.go:258-264`). The test fixture for the “ready” install path also omits auth while expecting success (`internal/cli/pi_frontdoor_test.go:100-123`). This conflicts with the setup text that says authentication is required (`internal/cli/pi.go:102-108`) and with doctor’s health definition including auth. Recommendation: either include `authOK` in install/launch readiness or rename/report this as “resources present, auth missing” so users do not get a false ready state.
  - Docs are stale relative to the frontdoor UX commit: `docs/runtime-support.md:184` says “When `spacedock install --host pi` is added…”, but the branch now adds `install --host pi` and `doctor --host pi`. Update this section so runtime-support docs match the implemented CLI.

- Polish:
  - Registry persistence is direct `os.WriteFile` (`internal/piruntime/registry.go:105-116`). Good enough for the current unintegrated contract package, but if this registry becomes active runtime state, consider temp-file + rename to avoid partial JSON after crash/interruption.
  - `spacedock pi` launch readiness also ignores missing auth (`internal/cli/pi.go:60-66`, `internal/cli/pi.go:258-264`), so missing credentials may surface as a less actionable Pi runtime failure instead of the existing doctor report. This is related to the Important setup finding.

Validation run: targeted cheap tests passed:

```text
go test ./internal/cli ./internal/dispatch ./internal/piruntime -count=1
```

Recommendation: mergeable from an implementation/design standpoint after fixing the auth readiness/reporting inconsistency and stale docs. No material blocker found.