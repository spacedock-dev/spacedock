# Ideation gate: lifecycle rejection semantics and proof ownership

## Capability and chosen direction

Preserve the validated bind, durable close, one-use consume, and terminal/nonterminal routing core. Restore automatic reviewer-`REJECTED` feedback routing, translate ordinary Captain language into the existing recorder decisions plus a durable reason, and make shipped skills—not live prompts—own the procedure under test.

## Design delta

- Reviewer `REJECTED` at a configured feedback gate routes before Captain presentation and creates no Captain Resolution.
- Captain `approve` maps to `approve`; `redo with feedback` and `reject` with `feedback-to` map to `revise` with distinct reasons.
- Captain `reject` without a correction owner, `hold`, and `not yet` map to `hold` with a nonblank reason and held resume state.
- The deterministic structural mutant owns rejection-branch proof; live rejection is positive integration evidence only.
- The host-neutral prompt retains only goal, fixture, delegated authority, and stop marker.
- Claude review extraction accepts only top-level assistant rows.

## Evidence and boundary

Validation's four counterexamples define the negative controls. Independent staff review approves the corrected mapping and proof split with no material finding. Implementation is limited to seven existing files, intended `+77/-16`, hard stop `+95`, and no eighth file. No recorder schema, product mapper, controller, harness, host adapter, transport contract, or compatibility layer is allowed.

## Recommendation and decision

Recommendation: **approve**. The reset closes the supported outcome defects and makes the live evidence falsifiable without adding machinery.

Decision: approve to resume bounded implementation on the existing branch; revise to change the mapping or proof boundary; or hold for a named prerequisite.
