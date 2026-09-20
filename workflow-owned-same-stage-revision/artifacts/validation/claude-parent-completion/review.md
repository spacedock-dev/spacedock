# Claude parent completion: independent validation

Recommendation: **REJECTED for C1 evidence-correlation defect** at `a34274756ff34fab2cf91c0c32cd701dbac92a8b`. The actual retained route evidence supports main; this finding concerns the new proof's ability to reject wrong-worker completion, not a claim that the observed successful host handoff was wrong.

## C1 proposal checkpoint

Proposed kind: evidence defect. Proposed scope: Material. Ownership: existing dispatch captured-handoff test. Proposed disposition: FIX, subject to distinct FO authorization; no reviewer candidate change.

Four fields:

- Released user/normal workflow: multiple named Claude background workers report to the same main conversation, and acceptance must match successful delivery with that worker's completion.
- Observable harm: the new regression proof passes when one worker successfully sends but never completes, while another worker completes without a successful send. Thus it cannot establish the claimed correlated completion invariant.
- Authority: contract[skills/ensign/references/claude-ensign-runtime.md#completion-signal] — confirmed delivery followed by the worker finishing its turn must yield the parent's native completion notification; a failed send is not completion.
- Trigger: detached exact-candidate test mutates only captured events: retains original worker A's successful main call/result; adds another worker B's main SendMessage call without result; replaces the completion with B's notification. Actual existing test exits0 although A lacks completion and B lacks confirmed delivery. Baseline exits0; removing completion entirely exits1.

The new test stores call IDs and parents initially, but collapses delivery and completion into maps keyed only by recipient. Two unrelated worker lifecycles can satisfy the final `delivered[main] && completed[main]`. This is a demonstrated owner-identity loss, not hypothetical prose interpretation.

Smallest proposal: keep the existing captured-handoff test and fixture. Require a successful result matched to its exact send ID/parent and a subsequent completion for that same parent. Preserve baseline, failed team-lead, absent completion and wrong-owner negatives; retain this cross-worker case in the existing owner. Do not add an observer/controller or change live graders or workflow policy. No candidate edit or corrective rerun was performed by reviewer.

## Actual detached results

`go test ./internal/dispatch -run '^TestBuildMergedModeCompletionSignal$' -count=1 -v` ran in a temporary archive of the exact candidate. Only the existing capture file was varied. No model/runtime was invoked and the temporary archive was removed.

| Capture | Expected | Actual |
|---|---|---|
| Original five events | pass | exit0,0.768s |
| Remove completion | fail | exit1,0.505s |
| Different worker completion sharing main recipient | fail | exit0,0.497s |

`results.json` retains argv/hash/status; the three JSONL files and logs retain exact mutations and outcomes. The wrong-worker mutation adds one copied SendMessage call with different send/parent identity and moves completion to that parent; it supplies no successful B result. The original successful A result remains unchanged.

## Real evidence and binding inspection

Original source CI35355801918 separate-review-required lines368,370,372,374,386 exactly reproduce the committed five-event fixture; source and fixture SHA-256 both independently match producer provenance. Worker parent `toolu_017EqNbekU7sexR5jyJKxn6b` first sends to team-lead with success:false, then sends to main with success:true, then completes in notification task `aab5ced63617ee2b2`. Success/result IDs and parent identity align in these original bytes. SendMessage input carries to/message fields (and host-expanded type/recipient/content); native responses confirm the actual supported main route. The same source's successful final notification does not prove the corrected candidate has run.

Scope is21 files125insert/63delete(+62): five existing generator/test/adapter owners, fifteen mechanical goldens, one five-event capture. Generated ordinary/advance completion uses the shared completion block targeting main; goldens cover both. Claude adapter ordinary completion, clarification and DISPATCH_FILE_MISSING route align to main; FO's normal completion and break-glass advance align as well. The four cycle-test recipient literals preserve anchored-call and Notify/prose-trap behavior. No independent reviewer, gate, cycle-limit, report or native identity outcome guard was changed. Existing wrong/missing identity and hold tests remain owned by TestSameStage. This inspection is contract coverage, not proof a model follows prose.

Producer final normal/race logs were read and hashed. Both exit1 solely for known TestCodexResolveManifestAgainstInstalledHost; all other package results pass/cache, no race diagnostic. Their earlier first race was interrupted and is not a pass. Focused dispatch/contractlint30.782s/2.081s, TestSameStage18.031s, and cycle owner0.586s are reused, not duplicated. Required formatting was performed per owned record. No wholly green suite is claimed.

## Acceptance limits

AC-1: historical same-stage correction/committed-plan/fresh-gate proof remains unchanged, but current capture alone does not establish fresh candidate runtime success.

AC-2: required review/round, committed reason, worker identity, cycle-limit and held authority outcomes remain strict and unchanged. C1 prevents accepting the new transport proof's asserted completion correlation until fixed.

AC-3: all six variants remain routine supported-host CI obligations; no skip/targeted-only exemption. Actual final-tip Claude cycle-limit must reach correlated parent completion after committed correction/report and retain escalation/held gate. Full801 acceptance waits for all six variants; Pi, roadmap refusal and timeout outcomes are not changed or waived here.

No candidate/frontmatter edits, model calls, push, rebase, CI, gate or round publication occurred. Scope is this report and evidence plus canonical entity report only.


### C1 FO disposition — 2026-09-20

FO separately authorized **FIX** for this Material owned evidence defect in the existing completion regression only. Preserve the exact falsifying mutation; producer must correlate send/result/parent worker identity, then retained reviewer reruns this exact mutant. No new observer or model runs, no further audit expansion, and no reviewer candidate edits.
