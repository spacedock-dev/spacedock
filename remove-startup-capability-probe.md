---
title: Remove boot-time gate withdraw capability probe
status: backlog
source: "Captain directive, 2026-08-15: pedantic compat check for unreleased 0.27.x is overkill; incomplete upgrades already surface a proper error at use time"
score: 0.9
group: tooling
id: dav9qnjhsbbg7k1a8x1260h6
gates:
    version: 1
    records:
        - id: gate:dav9qnjhsbbg7k1a8x1260h6:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:dav9qnjhsbbg7k1a8x1260h6-backlog-1
              briefing:
                id: briefing:dav9qnjhsbbg7k1a8x1260h6:backlog:attempt-1:revision-1
                digest: sha256:c5ca9ef1e81bd58a76d435a83a27ae7554485d193e7ecdfdbc432e13c7ccd1f4
                request-digest: sha256:fdc1e0b484e075d56c870855a84cc1484e02db4e2acc1144b3db517253a58d75
                room-ref: ./remove-startup-capability-probe/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:dav9qnjhsbbg7k1a8x1260h6:backlog:1
                briefing: briefing:dav9qnjhsbbg7k1a8x1260h6:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T02:53:13.080348Z"
                decision: approve
                reason: 'Captain ruling 2026-08-14 (dispatch them): approved into ideation'
              application:
                target-stage: ideation
                state: pending
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
