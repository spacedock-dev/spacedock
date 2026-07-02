# Validation: Pi intercom runtime capability probe

Entity: `docs/dev/.spacedock-state/pi-intercom-runtime-capability-probe/index.md`
Implementation worktree: `.worktrees/spacedock-ensign-pi-intercom-runtime-capability-probe`
Product commit inspected: `81a11d066fa369b1e70e25c9e7d46cc2693824b7`

## Recommendation: PASSED

## AC evidence

- AC-1 PASSED: `docs/dev/pi-intercom-runtime-capability-probe.md` documents the reusable probe recipe and evidence contract, and `skills/integration/pi_intercom_runtime_capability_test.go` validates recipe sections, allowed classifications, evidence fields, and interpretation rules.
- AC-2 PASSED: the recipe includes the exact child prompt requiring `contact_supervisor` `progress_update`, `contact_supervisor` `need_decision`, supervisor reply `APPROVED`, and marker content `PI-INTERCOM-SMOKE-APPROVED`; the recipe-shape test asserts these strings.
- AC-3 PASSED: evidence validation distinguishes setup/not-run from behavioral success. `passed` records require bridge-active setup, child tool availability, progress, decision, exact `APPROVED` reply, resume, marker path, and exact marker content; negative tests reject over-claimed passed/setup-only records.
- AC-4 PASSED: docs state bridge-active/setup is necessary but insufficient and must not be interpreted as supervisor talkback proof; invariant tests assert this distinction and reject obvious over-claim wording.
- AC-5 PASSED for static/manual-path scope: the live/manual smoke path and durable marker/evidence requirements are documented, and schema tests accept `passed` only with the required marker and behavioral observations. No live passed evidence was expected for this validation.
- AC-6 PASSED: checked-in `2026-06-04-not-run.json` honestly records `classification: "not_run"` with no progress, decision, resume, reply, or marker success claims, so tests remain reproducible without live Pi runtime access.

## Commands run

- `gofmt -l skills/integration/pi_intercom_runtime_capability_test.go` — passed, no output.
- `go test ./skills/integration -run 'PiIntercom|RuntimeCapability' -count=1` — passed.
- `go test ./skills/integration -count=1` — passed.
- `go test ./... -count=1` — passed.

## Residual risks

- Live child-to-supervisor Pi intercom talkback remains unproven until an operator runs the documented smoke and records durable non-`not_run` evidence.
- The schema permits unknown/null package version/path observations, consistent with the recipe's "use null when unknown" instruction; future live passed evidence should still capture as much resolved setup data as practical.
