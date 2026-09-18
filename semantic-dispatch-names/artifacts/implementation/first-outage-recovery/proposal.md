# First-outage recovery proposal

Captain ruling, 2026-09-18: “both yes” requires preserving first-dispatch recovery when no cached worker identity exists. This supersedes naming approach 5 / AC-3's cached-envelope-or-hold rule for the observed build-only outage. The captain approved the behavior; FO subsequently authorized the exact seven-path read-only interface via FIX before candidate mutation.

Proposal baseline: `f5af7282dda9151771ff9888dd246553e38afa82`. Implementation and results are recorded in `results.md`.

## Finding and evidence

- Supported workflow: first selected-team dispatch when every `dispatch build` call fails and no successful envelope exists.
- Harm: the current recovery instruction holds instead of spawning a worker, so no committed worker report is produced.
- Authority: captain-ruling[2026-09-18] preserve first-dispatch recovery, selected mode, safe identity, one worker and committed report without a cached identity.
- Trigger: finding 3 in `workflow-owned-same-stage-revision/artifacts/validation/tip35238049892/report.md` records two real build failures and zero Agent calls. The fixture shim rejects every `dispatch build` but forwards other commands. The runtime followed the old hold rule; this is a naming/recovery contract conflict, not proof of model disobedience.

Classification: Material. Ownership: semantic-dispatch-names and its recovery integration. Disposition: FIX authorized for the concrete mechanism and implemented. No existing failed live outcome is waived.

## Smallest proposed mechanism

Expose `spacedock dispatch name --workflow-dir DIR --entity-path FILE --stage STAGE` as a narrow read-only query. Emit one validated canonical base name on stdout. Reuse the existing `validateWorkerName` and `semanticName` owner, including long-name hashing, stage budget and cross-generation collision refusal. Do not duplicate that logic in skill prose.

The query must require an existing readable entity and workflow README, a declared stage and valid arguments. It must reject worktree-copy entity input and invalid/ambiguous/over-budget identity without emitting a name. It does not read checklist stdin or require a resolved launcher/model. It writes no assignment, stamp, worktree, commit or artifact.

Selected-team recovery retains a validated same-assignment envelope when available. If none exists after build failure, call this query once and feed the successful result to the existing named-background template. Preserve occupancy checks, bounded suffix handling, exact selected transport and one-worker/committed-report obligations. Selected-bare recovery remains unchanged. If this separate identity query fails, report that failure rather than inventing a name or retrying the failed build.

## Existing-surface assessment

`dispatch build` is currently the sole public caller exposing the canonical generated name. A new build flag still hits the fixture's unconditional build failure. `show-stage-def` emits README/context text; `reconcile` describes existing workers and Git state; standing commands generate standing-agent assignments. None currently returns a fresh validated entity/stage identity. Adding identity projection to those unrelated outputs would create a less clear interface. A narrow `dispatch name` route reuses the canonical owner without full build effects, but is still an explicit interface addition for FO decision.

This proposal addresses the evidenced build-only outage. It does not claim that a completely missing launcher executable can provide a safe canonical name without another mechanism.

## Proposed surface and falsifiable tests

Seven existing paths; estimated +160–190 net lines, to be measured after authorization:

- `internal/dispatch/dispatch.go`: route and help for the read-only query.
- `internal/dispatch/names.go`: query validation and canonical-name emission, reusing current derivation.
- `internal/dispatch/names_test.go`: behavior-first query tests, long/sibling/legacy ambiguity and no-mutation negatives.
- `skills/fo-dispatch-recovery/SKILL.md`: first-outage query then existing selected-mode template.
- `skills/integration/dispatch_test.go`: real failing-build shim and uncached identity recovery; retain instruction contract assertions as supporting evidence only.
- `skills/first-officer/references/claude-fo-dispatch.md`: name-source wording consistent with the same canonical query.
- `docs/site/reference/command-reference.md`: document the bounded public query and its effects.

Before implementation, add focused behavior tests in existing owners. Compare query names with independently specified short names and canonical build output for long names. Invalid, missing, ambiguous and insufficient-budget cases must emit no name and preserve entity/Git/artifact bytes. Under the actual first-outage shim, build must stay failed while the separate query succeeds with no prior envelope. Existing mode, duplicate-worker and committed-report grader controls stay intact. Do not add a standing harness, exemption, alternate transport or local model run.

No code mutation, new test run, push, rebase, CI or agent dispatch occurred during this proposal. Final full normal/race/gofmt checks belong once on the final restacked tip, as assigned. The #801 missing-source correction and retained runtime/root failures remain separate.
