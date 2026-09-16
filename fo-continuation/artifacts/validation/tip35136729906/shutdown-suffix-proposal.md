# Proposed stop-boundary correction to the Claude shutdown suffix

Design preparation only; implementation is not authorized. Inspected at ba8e6eeb4e28e310b17897a0a8b12e151366b8f9. No code change, test, model call, or CI launch was made.

The current suffix has an unconditional completion demand that conflicts semantically with legitimate stops before workflow completion. Narrowing that demand can remove this concrete instruction conflict. It cannot establish that the suffix caused the fabricated approval in run 35136729906, or that a wording change will prevent recurrence.

Current literal, internal/ensigncycle/claude_live_runner_test.go:31:

> Do not shut down your team or prepare your final response until all the work is complete. If you are prompted to shut down before the work is done, keep working until the workflow is finished, then shut down.

Proposed literal:

> Do not shut down your team or prepare your final response merely because a shutdown reminder appears or a worker is dispatched. Continue authorized work and unresolved-worker monitoring until the requested scope is complete or the task or operating contract declares a stop, including a captain stop, an unresolved captain decision, or a blocker with no independent authorized work remaining. This instruction grants no gate approval authority.

The proposal retains the suffix's original purpose: a host shutdown reminder or worker handoff must not prematurely end authorized work. It bounds continuation by the task's existing completion and stop rules. It does not create a new stop detector, authority token, gate mechanism, or runtime adapter. “Requested scope is complete” includes a task that authorizes driving to terminal; an unapproved human gate remains a decision stop. Existing delegated approval remains governed by its actual grant.

## Exact surface and callers

Implementation, if approved, is one existing constant plus its adjacent explanatory comment in internal/ensigncycle/claude_live_runner_test.go. Update the comment to describe timing within authorized scope and declared stops; do not change CLI flags, tools, watcher budgets, launch wiring, fixtures, or grader rules. The same literal has three call sites:

- claude_live_runner_test.go:752 appends it to every Claude shared-scenario launch. This includes both layouts of auto-continue, default headless gate stop, keep-moving, rejection, and full ensign cycle.
- merged_team_mode_live_test.go:109 appends it to TestLiveMergedTeamModeDispatch. That caller explicitly grants conn to resolve gates from reports, so legitimate progression to terminal must remain supported.
- haiku_loop_spike_live_test.go:214 includes it in haikuLoopPrompt, used by TestLiveHaikuLoopSpike. This mechanical experiment has no separate FO contract, directs bare-mode dispatch and terminalization, and already declares HALT on a state commit conflict. Referring to “task or operating contract” preserves this caller's existing specification without imposing a new contract. The spike is not added to CI by this proposal.

The suffix is test-launch prose, not shipped #804 skill prose. The existing Claude live-runner owner should own an approved correction, coordinated with the #804 validation owner. No edit to #804's shared final-response rule is proposed.

## Existing validation and its limits

Reuse the existing Claude cases and assertions; do not add a new framework or a wording-presence test. Existing offline negative controls already passed and need no repeated run for this proposal. They prove the oracle rejects unauthorized resolution, not that the model follows a revised suffix.

If implementation is approved, existing Claude CI cases provide the bounded behavior checks:

- TestLiveCommonAutoContinueAfterImplementation, both layouts: completion proceeds to fresh validation, then leaves the human gate unresolved. A closed gate with consumed=false must remain red.
- TestLiveCommonDefaultHeadlessGateStop: headless operation without conn stops at its human decision boundary.
- TestLiveCommonKeepMovingPosture and TestLiveCommonRejectionFlow: authorized handoffs and correction continue through their required work to the declared stop. Preserve the existing strict rejection topology.
- TestLiveMergedTeamModeDispatch and TestLiveCommonFullEnsignCycle: the explicit authority in their existing prompts still permits the required terminal outcome.

These are existing members of the scheduled Claude lane. The FO owns selection and CI authorization; no new run is requested or triggered here. A later passing run can establish observed compliance at the corrected tip. It cannot prove causation from this one failed sample or measured improvement without a suitable comparison.

## Recorded no-op disposition

FO disposition supersedes the earlier report's topology-choice question: DECLINE narrowing or filtering. Preserve the existing strict topology and red result for the two parent-created no-op agents. Their unnecessary dispatch is observed runtime conduct. No product fix is justified by that sample, and no captain decision is required to leave the existing bar unchanged.

## Decision requested

Approve or decline only the literal replacement and adjacent comment correction above, owned by the Claude live-runner surface. Approval would authorize removing the semantic stop conflict, not declaring a proven causal fix or weakening either failed assertion.
