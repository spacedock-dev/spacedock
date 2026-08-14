---
title: Remove boot-time gate withdraw capability probe
status: backlog
source: "Captain directive, 2026-08-15: pedantic compat check for unreleased 0.27.x is overkill; incomplete upgrades already surface a proper error at use time"
score: 0.9
group: tooling
id: dav9qnjhsbbg7k1a8x1260h6
---

Remove the same-minor capability probe added in commit b331baf4f ("fix: reject stale same-minor launchers", 2026-08-09). The binary version gate should rely solely on the minor version match; if a stale build lacks `gate withdraw --reason`, the CLI already returns a clear "unknown subcommand" error at the point of use.

## Scope

Delete the capability check from these four locations:

1. `skills/first-officer/references/first-officer-shared-core.md` — remove the "Compatible minor" probe bullet and the "Missing capability" abort bullet from the Startup section.
2. `skills/integration/testdata/version_gate_flow.sh` — remove the `REQUIRED_CAPABILITY` variable and the entire probe block (the `helpout` / awk check and its failure branch); restore the simple `say "gate passed: spacedock $ver"; exit 0` path.
3. `skills/integration/version_gate_fixture_test.go` — remove the `withdrawCapability` const, the `gate --help` branch from `writeGateLauncher` and `captiveInstall`, the `INVOCATION_LOG` plumbing from `captiveInstall` and `runGateFlow`, the `invocationLog` helper, and both tests: `TestGateFlowRejectsStaleSameMinorBeforeBoot` and `TestGateFlowCompatibleSameMinorProbesThenBoots`.
4. `docs/site/get-started/install.md` — remove the paragraph about "If startup says the installed launcher is missing a required command".

## Acceptance criteria

- **AC-1**: `version_gate_flow.sh` passes the gate on a same-minor binary without probing `gate --help`.
- **AC-2**: All remaining version-gate fixture tests pass (`go test ./skills/integration/ -run TestGateFlow`).
- **AC-3**: No reference to `REQUIRED_CAPABILITY`, `Missing capability`, or the `gate withdraw` probe remains in `skills/` or `docs/site/get-started/install.md`.
- **AC-4**: The FO skill's startup section still enforces the minor-version check (0.27) and the binary-absent / wrong-version abort classes unchanged.
