# Claude parent completion routing

Disposition: FO FIX authorized under captain request “fix 801’s claude issue”; Material outcome defect, prerequisite to AC-3, within #801.

1. Supported workflow: a named background Claude ensign completes and commits its assigned correction/report, then delivers completion to the FO.
2. Harm: mandated team-lead route fails; cycle-limit correction81d738d is committed but unreachable-recipient recovery stalls, watchdog terminates, parent completion and final assertions never occur.
3. Authority: contract[skills/ensign/references/claude-ensign-runtime.md#completion-signal] and captain request. Committed report and actual correlated parent completion remain required.
4. Evidence: retained CI35355801918 cycle-limit SendMessage toolu_012ghCyGy9HyzpjZgUjE85p4 returns success:false. Separate-review-required's same validation worker first fails team-lead then succeeds to main, followed by completed task notification. Exact five-event capture and provenance are retained; no inferred recipient.

The installed host tool description maps main to the main conversation for background subagents. The existing FO adapter already names that route in its opening binding, but later completion instructions contradict it. Fix the existing generator/adapters/examples; no new routing mechanism. The generated-target test replays captured actual call/results and parent completion. It does not simulate a model or claim the candidate completed a live journey.

Local baseline reconciled cleanly from581ed16de to origin735446bf8, full tree equal. No candidate behavior changed during reconciliation.
