# Implementation proof index

All counted before/after live drives used `/opt/homebrew/Caskroom/spacedock@next/0.28.0-pre3/spacedock`, the Codex local isolated OAuth lane, and the assigned worktree plugin source. Skill baseline SHA-256: `9896b43ff641559bcd0965bdb6715030aff91c9285a0a708914954db984eb934`; unchanged candidate skill SHA-256: `cf43979ac026c3a3964956853f122d93b5dbc18aa3d0b478f094024a97cdd75d`. Runtime homes and auth files are excluded.

| Evidence | Outcome | Meaning |
| --- | --- | --- |
| before/ and before.log | Expected failure | Worker completed and committed; mandatory absent round stalled the flow at one attempt. |
| after/ and after.log | Pass | Same launcher, changed skill: one completed worker, corrected committed plan, exactly one fresh open attempt, no round/reviewer. |
| initial-controls/review-required | Pass | Same-stage correction completed; missing required review evidence recorded; no gate. |
| initial-controls/separate-review-required | Pass | Distinct implementation and validation workers completed; missing evidence holds the gate. |
| conventional/ and initial-controls.log | Pass | Existing full correction/re-review flow retained four entries, Cycle projection, reviewer reuse and fresh open gate. |
| initial-controls/cycle-limit | Regraded pass | Original test falsely matched a prose mention. Retained state has no Feedback Cycles heading, records explicit cycle-3 captain escalation, and preserves three closed attempts without a fourth. |
| initial-controls/round-required and round-missing | Setup failure, not acceptance | Fixture lacked explicit identity; round refusal did not establish missing-log enforcement. |
| identity-fixed/round-missing | Pass | Recorder refuses the genuinely missing canonical log; no round, projection, or new gate is invented. |
| identity-fixed/round-required | Partial, overall failure | All four canonical entries publish without projection, and a fresh gate exists; stronger selection grader rejects its omission of the intended plan source. FO authorized fixture alignment; the affected final rerun is recorded separately. |
| round-final/ and round-final.log | Pass | Aligned plan acceptance and explicit review source: completed correction, four-entry round without projection, selected corrected plan, fresh open gate; old attempts unchanged. |
| baseline-normal.log | Existing baseline failure | Only TestCodexResolveManifestAgainstInstalledHost failed before candidate edits. |
| premature-selection-negative-before.log | Expected failure | The old current-file-only rule accepts stale and unrelated selections; the fixed focused test rejects both and accepts current intended sources. |

## Retained selection regrade

Before attempt 1 selects `git-root://state/9e4f264615db4234ce41de611b5a429b2c949588/recorded-gate-task/selected/plan.md`, digest `sha256:6ed9839b7be73d0341c37348f0ddc266c25e09575f06165a100d34286c94777a`, whose bytes remain the incorrect plan. After attempt 2 explicitly selects the plan as a canonical Reference at `git-root://state/31902baa85ecca70a76ae02d4768a6d205fb0587/recorded-gate-task/selected/plan.md`, digest `sha256:9c9f788375c0bd7caf0b764702f68cd519c4dbfe70a721f84b172ec2572ce898`, whose bytes are exactly `KEEP message A; DELETE message B\n`. Both object digests were recomputed from their bundled Git objects; the corrected expectation is an independent literal.

## Registration and verification history

The initial complete normal/race runs each failed only the existing installed-host resolver and the new registration omission; see pre-registration-fix-normal.log and pre-registration-fix-race.log. The temporarily integrated wrapper fixed registration at +217/four files, but FO held its added CI cost before commit. The held-integrated-normal.log and held-integrated-race.log runs additionally hit cumulative ten-minute package timeouts; race also failed four existing 250ms process-budget checks.

All newly failing checks then passed unchanged under isolated serial execution after peer broad suites stopped. See isolated-diagnosis.md and isolated-*.log for exact commands, timings and outcomes. Only the existing resolver failure reproduced. These results support resource contention without erasing failed broad evidence or substituting a focused pass for final full verification.

## Final approved surface

Binding resolution:binding-1789513304916529000 explicitly approves targeted implementation-proof policy and +236 net/five files, bound to /private/tmp/inflight-implementation-decisions.md at sha256:3e4c7f21b5d53b1abfcc5796d19b7bbe7d90ff2f5c5398415f1e11b245fe9c81. The standalone entry point is now TestLiveSameStageRevision. All fixtures, selected-runtime dispatch and assertions are unchanged; no successful live controls were rerun. The normal common-journey selection is unchanged. Passing local controls are mandatory task acceptance, not CI/all-runtime parity evidence.

Final focused checks pass: final-registry.log reconciles actual test registration and suite policies; final-focused.log exercises worker obligations and stale/unrelated selected revisions; final-discovery.log proves the live-tagged targeted entry compiles and is discoverable. Final sequential normal/race verification completed in the FO-reserved isolated broad lane; final-normal.log and final-race.log each fail only the same installed-host resolver baseline. All other packages pass, including ensigncycle at 293.503s normal and 422.191s race; no timeouts or data-race reports. Candidate commit: 5e04442a9c89ef644135721522dbf9bfcc8b7de0, based on 2a7b8719843e40b79545f0bb4def6609cdd9ebbf, +236 net/five files (247 insertions,11 deletions), clean assigned branch. Independent validation belongs to the next dispatched stage.

## Final verification commands

- `go test ./internal/contractlint -run 'TestRuntimeLive(RegistryReconciliation|CommonSuiteTimeouts|CommonFailFastPolicy|GapBindingValidation)$' -count=1`: pass.
- `go test ./internal/ensigncycle -run '^TestSameStage|^TestRejectionTopologyRedsNonConformingShapes$' -count=1`: pass.
- `go test -tags live ./internal/ensigncycle -list '^TestLiveSameStageRevision$'`: targeted entry compiles and is discovered; this is not a live execution claim.
- `go test ./...`: exit1, only known TestCodexResolveManifestAgainstInstalledHost baseline.
- `go test ./... -race`: exit1, same sole baseline; no race report.
- `gofmt -w ./cmd ./internal` was run during implementation; its unrelated pre-existing release-test alignment change was restored to preserve scope. Final changed-file `gofmt -l` and `git diff --check` are clean.

The frozen final diff SHA-256 before commit was ff70c1671779ed7b0d6f931a9e77c2e69e575b3a1189b6c9973a807ff54b9fd4 before and after final suites. No new mechanism, code push, CI trigger or restacking occurred.
