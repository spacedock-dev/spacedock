# Gate decision: completion-guard diagnostics

## Recommendation

Approve this design for implementation.

The work replaces one generic completion-guard error with four actionable failure classes while preserving existing completion and gate-authority fences.

## What changes

- Missing report, incomplete checklist, untracked entity, and dirty entity failures get distinct remedies.
- A checklist bullet such as `DONE (annotation):` becomes a visible refusal instead of being silently ignored.
- The valid ungated terminal path gains regression coverage and concrete command documentation.

## What remains unchanged

- Command grammar, stored formats, gate authority, exit codes, and local-HEAD durability rules.
- Pending terminal approvals still use `merge guard` as their sole consumer.
- No remote dependency or broader merge-guard redesign is added.

## Cost and proof

- Expected change: +145 net lines across six files, with tolerance of plus or minus 40 lines and one file.
- Real CLI fixtures distinguish the four core failures and require byte-clean refusals.
- Parser and CLI fixtures cover malformed checklist bullets while preserving valid syntax.
- Terminal fixtures cover merge-guard, direct ungated completion, and the pending-approval control.
- Focused, full, and race test suites are required.

## Evidence note

The earlier automated AC scan did not attach a range citation to AC-2 and AC-3. The task defines explicit, falsifiable proof for both. This is not a missing design obligation.

## Captain action

Approve to enter implementation, or hold if the deliberate checklist tightening or six-file surface is unacceptable.
