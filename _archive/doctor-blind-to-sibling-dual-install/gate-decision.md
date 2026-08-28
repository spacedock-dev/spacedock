# Backlog gate: expose cross-channel runtime drift

## Recommendation

Approve ideation as critical 0.27.x work. The launcher can certify one plugin while the host loads another channel's incompatible skill, leaving the operator with a green doctor result and an immediate pre-boot abort.

## Observed failure

A 0.27.1 binary and doctor-selected plugin agree at 0.27.1. The runtime-loaded first-officer skill instead requires 0.28 and aborts before boot. The binary compatibility gate is correct for the manifest it sees; the unresolved defect is that doctor and the runtime can select different installed sources.

## End value

Doctor, install, and runtime selection must expose or prevent a cross-channel skill/manifest mismatch without weakening the binary gate.

## Ideation charge

- Reproduce the mismatch on both Claude and Codex, and identify the exact host and source-selection paths used by doctor, install, and runtime loading.
- Decide whether the smallest safe behavior is nonfatal doctor reporting, install cleanup, front-door enforcement, or a deliberate combination.
- Define stable fixture-backed proof and retain a live host proof where host selection cannot be represented faithfully by fixtures.
- Preserve the binary compatibility floor. Merely lowering the first-officer skill requirement is explicitly rejected because it hides source drift rather than fixing selection or diagnosis.

## Scope boundaries

- Do not weaken or bypass the binary gate.
- Do not assume the existing Claude sibling-cleanup behavior also proves Codex runtime selection.
- Do not prescribe fatal launch behavior until ideation compares its compatibility cost with reporting and cleanup.

## Reference

Seed entity: `doctor-blind-to-sibling-dual-install.md` (`x0petxt7xvr459b6zh4vf4wj`).
