# Independent combined-tip validation — 7c63ad87

Recommendation: **PASSED for bounded deterministic validation of the #806 chronology correction and final combined integration.** This is not full-suite green or final native acceptance. Checklist: 3 DONE, 0 SKIPPED, 0 FAILED; these count completed validation work, not successful command exits.

## Exact candidate and scope

Frozen candidate: 7c63ad87f34a18c80e3e8b24f6d201e227f661e6 on spacedock-ensign/pi-native-completion-evidence. Initial and post-probe worktree status are clean. No candidate commit or persistent source edit was made.

Implementation source 02ef599d016f7c3f48def7f7437f35a3f99163e2 and this restacked candidate differ only in three #801 coverage files, +3/-6. Direct local Git comparison confirms that boundary; FO owns the complete patch-equivalent restack evidence. The #806 correction remains eight files, +438/-7, including 144 blank records preserving native stream indexes. Its 14,698-byte fixture/provenance set retains the actual dispatch/completion/gate evidence. Local SHA-256 checks match all three original parent/child files.

The continuation Git state is reconstructed with the existing fixture owner: there was no original continuation Git bundle. The replay never claims those reconstructed commits are native CI commits. Actual CI chronology and commits remain documented in ../tip35238049892/report.md, state f522daac8. The first report's omitted checklist item is retained as a recovered runtime defect.

## Independent semantic adversarial pass

The direct scan joins successful gate calls/results by tool-call ID and matching briefing identity. Every prepare must follow its preceding completed native worker. The active attempt must be successfully withdrawn before a replacement or repair worker. An unchanged validated report survives withdrawal and can be prepared again without another worker. Current report content is compared against its committed HEAD section by the existing Git durability owner; frontmatter-only changes do not require another validator.

The independently unowned claim tested here was whether the new negative cases discriminate when those guards disappear. Three throwaway Go overlays replaced existing files only for their individual test process; candidate bytes stayed unchanged. The existing test owners supplied all inputs/assertions. No new harness or second green-focused run was added. Overlays were deleted after testing.

| Mutation | Existing control outcome |
| --- | --- |
| Remove per-prepare completion-before-boundary comparison | first-gate-premature detects erroneous acceptance and fails; repair-gate-premature remains rejected by the independent final boundary guard |
| Remove active-attempt guard at new worker and replacement prepare | both omitted-withdrawal and one-validator reprepare-without-withdrawal detect erroneous acceptance and fail |
| Remove current-report versus committed HEAD comparison | uncommitted-repair detects erroneous acceptance and fails |

Each mutant command exits 1 as expected; logs contain “invalid repair chronology/evidence accepted.” These are deliberate fault injections, not failures of the frozen candidate. Exact commands, wall durations and diffs are retained in mutation-results.json and mutation-*.log/.patch.

The owner already proved both captured repair and unchanged-report reprepare positives, plus 16 negatives in 7.918s. Those green focused checks are reused, not rerun for a second opinion. Existing identity/run/owner/task/epoch and error/missing-completion controls remain in the shared verifier. Historical sync/async and Claude/Codex controls, original missing-gate and uncommitted-report negatives remain under their existing test owners and are exercised by the required combined suites. Native completion141 now pairs with replacement144 while first completion68 still precedes first prepare99. No mandatory-new-worker, Briefing-read, whole-entity digest comparator, workflow-wrapper parser or naming change was introduced.

## Integration evidence reused

Independent #801 review 8a3e88ba2 at workflow-owned-same-stage-revision/artifacts/validation/pi-complete-coverage/review.md owns the fatal-child/parallel-parent/gotestsum experiment. Removing Pi -failfast admits all six children while retaining the failing child and nonzero command exit. That audit was not duplicated. It establishes selection/failure preservation, not successful native Pi outcomes. Exact final combined normal/race checks below include the coverage contract-lint change.

## Formatting

Required `gofmt -w ./cmd ./internal` exited 0. Its sole difference was the known unrelated internal/release/runtime_live_evidence_workflow_test.go spacing; exact pre-command bytes were restored. format-diff.txt records it. Candidate status and `git diff --check` are clean. This does not claim the unrelated baseline file is gofmt-clean.

## Acceptance evidence and limits

- AC-1: captured per-attempt chronology, unchanged-report reprepare positive and strict native identity controls are supported locally. Independent mutations establish that removal of order/withdrawal guards is detected. Existing sync result and async notification handling remains intact. WorkflowScript/runs.run and revival-epoch support remain HOLD, not silently accepted.
- AC-2: current report durability is enforced; bypassing it fails the uncommitted-repair control. Historical absent-gate evidence remains red through its existing owner. No host-specific product/runtime behavior changed. The report omission in the original first attempt remains visible rather than being erased by the repaired final state.
- AC-3: local deterministic evidence is recorded below with actual exits. Final native host CI is pending. CI35238049892 remains red; Pi plain remains incomplete and the other five variants were unstarted. Coverage correction supplies no final-tip native pass. All six same-stage variants still require actual execution on each supported host.

## Held outcomes and authority

No new material, deferred-risk or polish finding arose from this bounded review. Existing FO dispositions remain: Pi plain incomplete report/provider error HOLD, without an authentication-root-cause claim; workflow-wrapper/revival observer gap and Pi fixture Claude-host wording HOLD; Claude missing-source recording, roadmap refusal and naming/recovery findings HOLD. None was repaired, rerun or waived. Earlier default-timeout failures remain retained in final-tip-0997; this review uses the previously authorized serial 30m observation budget and does not alter code, timeout policy or pass thresholds.

No local model/live run, CI, push, rebase, gate action, frontmatter mutation or extra agent was used. FO owns publication and any later native run.

## Required combined commands and actual results

Both ran once, sequentially, on exact frozen 7c63ad87f34a18c80e3e8b24f6d201e227f661e6, using natural Go cache behavior and no forced count. The authorized observation budget is command-only.

| Command | Exit | Wall time | Ensign-cycle result |
| --- | --- | --- | --- |
| `go test ./... -p 1 -timeout 30m` | 1 | 824.409s | PASS, 330.118s |
| `go test ./... -race -p 1 -timeout 30m` | 1 | 645.671s | PASS, 246.778s |

Every other package passed in both logs. Neither run timed out; the race log has no race diagnostic. Each sole failing test is TestCodexResolveManifestAgainstInstalledHost, with the exact previously FO-declined failure:

```text
codex_resolve_test.go:44: spacedock@spacedock not installed in codex, but resolver returned "/Users/clkao/.codex/plugins/cache/spacedock-local/spacedock/0.28.0-pre0/.codex-plugin/plugin.json"
```

The resolver remains red; its prior DECLINED disposition does not turn either suite green. No new failure required classification or correction. Final candidate HEAD matches the initial HEAD, tracked/untracked worktree status is clean, and diff-check exits 0. No overlay remains.

## Stage Report: validation (cycle 3)

- DONE: Independently verify per-attempt native completion and current committed report without forcing a new worker for unchanged-report gate re-preparation.
  Three independent fault-injection overlays prove early-gate, withdrawal and current-report durability controls discriminate; existing two positives and identity/error negatives are reused and covered by combined suites.
- DONE: Verify exact combined tip integration and required normal/race/format checks, retaining genuine failures and reusing already-owned coverage evidence.
  Exact 7c63ad87: sequential normal/race exit1 solely for the retained resolver failure; ensigncycle passes 330.118s/246.778s. gofmt ran; unrelated baseline spacing restored exactly; independent #801 selection audit reused.
- DONE: Record bounded recommendation, all AC evidence/limits and unresolved live outcomes without approval, push or CI.
  AC-1/AC-2 deterministic proof recorded; AC-3 native host acceptance pending. Original Pi/Claude held outcomes and five unstarted Pi variants remain explicit.

### Summary

Recommend PASSED for this bounded deterministic correction and integration review; no new finding. Both full suites remain exit1 for the exact FO-declined resolver defect. Final native host runs and all six same-stage variants per host are still owed. Candidate/frontmatter remain unchanged; report and artifacts are committed locally for FO synchronization and fresh-gate handling.
