# BACKLOG GATE — Pi default extension discovery (`3w1`)

Recommendation: **APPROVE and dispatch ideation.**

## Capability and value

An isolated Pi home does not auto-discover the `pi-subagents` / `pi-intercom`
extensions. The live harness hand-wires this: `piSubagentsPackageRoot` requires
`PI_SUBAGENTS_PACKAGE_ROOT` or falls back to `~/.pi/agent/npm/node_modules/pi-subagents`;
the smoke passes `--extension .../src/extension/index.ts` explicitly. `runtime-support.md:147`
names this as first-contact friction that is supposed to be harness work ("an
extension not auto-discovered in a temp home... is harness work"). This task
makes the isolated home discover the operator's installed extensions the way a
normal home does, so no env var or hard-coded `--extension` path is needed.

## Binding boundaries

- First-contact default reduction, not a journey repair. No journey exercise,
  fixture, or durable assertion changes.
- Touches the same `seedPiLiveAuth` / isolated-home setup seam that pnc just
  landed (models.json mirror). Coordinate to avoid a divergent setup seam.
- Does not change where Pi itself stores or loads extensions.
- No new runtime, fixture, result format, or CI lane.

## Proof direction

Ideation confirms the discovery reads the operator's real installed
extension/package location (Pi home npm node_modules or package manifest), not a
hard-coded absolute path; the explicit `PI_SUBAGENTS_PACKAGE_ROOT` override still
wins; the `--extension` wiring in the smoke stays valid. Implementation proves
`TestLivePiFrontDoorSmoke` passes with no env var exported.

## Decision ask

Approve this first-contact default for ideation, or revise/hold with a concrete
boundary.
