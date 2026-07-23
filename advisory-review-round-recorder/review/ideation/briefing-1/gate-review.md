# Ideation gate review: advisory review rounds through 3k

## Capability and change

Extend 3k's existing recorder with one advisory-round operation. It retains an ordinary Review & Gate Briefing/log in a derived round room, publishes a minimal current pointer and one Feedback Cycles projection, and supports exact-prefix append/replay. It reuses 3k's digest, identity, lock, CAS, validation, rollback, and atomic entity writer.

The operation cannot select or close a logical gate, create/consume an application, or change workflow status. It is the deferred persistence half of 02av, not a second recorder or decision model.

## Evidence and reviewed snapshot

- 02av already fixes the semantic mapping: reviewed snapshot → Briefing; reviewer findings → Annotations; reviewer and worker dispositions → separate advisory Resolutions; all-declines differs from no findings.
- The 3j jobs 592/594/597 incident demonstrates the missing value: an adversarial duplicate-member finding caused a 271-addition/137-deletion rewrite, then job 597 passed the unchanged candidate `90aea55` after the finding was declined. Only prose retained that round.
- The proposed minimum replay persists jobs 592/594 plus the worker decline, keeps candidate/product/gate/status bytes unchanged, distinguishes all-declines from no findings, and refuses a bad reviewed-snapshot digest byte-clean.
- The design adds no provider launcher/poller, materiality judge, retry daemon, journal, arbitrary room path, gate/application mutation, or second package/schema.

Reviewed artifacts are the exact entity, the landed gate/round contract, and the 02av materiality fixture identified by URI and SHA in `briefing.json`.

## Findings

- Material: 02av's required advisory disposition can currently be preserved only by hand-authored room records; the 3j incident instead landed prose and lost the structured round.
- Material boundary: round Resolutions must never inherit a selected gate's binding/application semantics.
- Material implementation risk: multi-file room plus pointer publication must reuse 3k's lock/CAS/rollback and remain byte-clean on every returned failure.
- Deferred: ordinary FO gate invocation is owned by 6y; round caller text is limited to the two existing triage triggers.
- Coordination: 6y and this task may both touch the command reference. Product/skill work can proceed concurrently; the later branch must rebase and make the shared doc edit without overwriting the earlier one.

## Recommendation and decision

Recommendation: **approve ideation and proceed to implementation** within the one-recorder boundary. Require the exact 3j decline replay, bad-digest and divergent-replay refusals, existing 3k regressions/race suite, and a surface audit proving no second recorder or workflow effect.

Decision requested: approve, revise with a concrete boundary finding, or hold for a named prerequisite.
