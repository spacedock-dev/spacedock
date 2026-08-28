# Ideation gate: expose sibling plugin conflicts

## Recommendation

Approve and enter implementation. Live Claude and Codex tests reproduced the same split-brain: doctor checked the stable 0.27.1 plugin and reported `OK`, while the runtime selected the enabled edge 0.28.0-pre0 first-officer skill and rejected the 0.27.1 binary before boot.

## Selected design

- Add doctor-only inventory from normalized `<host> plugin list --json` output.
- Print nonfatal `CONFLICT` output when an installed and enabled sibling channel can win runtime selection.
- Print nonfatal `INCOMPLETE` output when doctor cannot inspect enablement.
- Preserve the compatibility exit code. Do not lower the binary floor or add a front-door gate.

## Proof plan

- AC-1: Claude and Codex command tests prove that an enabled sibling produces exact `CONFLICT` output and does not produce a lone `OK` result.
- AC-2: stable 0.27.1 single-channel and dual-channel fixtures still reach launch, with bounded stable-only live boot checks on Claude and Codex.
- AC-3: live-captured Claude and Codex JSON fixtures cover both channels enabled, the selected channel disabled, and the sibling disabled.
- AC-4: command-reference examples match the AC-1 output, and `mkdocs build --strict` passes.

Claude and Codex live validation is required before completion. It must repeat the isolated dual-install reproduction and the stable-only boot check.

## Estimate

The implementation estimate is +255 net lines across 7 files. The accepted tolerance is ±80 net lines and ±2 files.

## Decision effect

Approval enters implementation. The entity remains a Reference.
