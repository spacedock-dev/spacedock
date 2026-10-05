# Independent validation

Candidate: ba8e6eeb4e28e310b17897a0a8b12e151366b8f9 on spacedock-ensign/pi-native-completion-evidence. Parent: 49d0c7a962ae500168a3a11b6d91136a59ebdb19. Candidate untouched; final worktree clean. Recommendation: PASSED for local validation; final tip CI pending.

## Authority and scope

Captain binding resolution:binding-1789576723761874000, approved brief SHA-256 fb4c087fbde4cad12db9080553c5a85a84a7ada2f504da760151d8c03fbcc3d5. Independent validation dispatch forbids candidate/frontmatter changes, native/model runs, push, CI, PR and rebase. Implementation proposal records distinct FO FIX authorizations for G1 and conflicting epoch; implementation report records FO DECLINED for the anticipated installed-manifest defect. This review found no new finding requiring disposition.

Exact diff: 12 files, +424/-11. One 150-line Pi evidence helper, existing parser fields and lifecycle call sites, existing replay test owner, four JSONL projections and provenance. No runtime controller, synthetic completion event, additional Git parser, instruction or production change. Both live callers pass retained artifact roots. Existing commit and gate checks retain ownership.

## AC evidence

AC-1: PASS locally. Source SHA-256 and every projected field/ordered content block were checked against all four retained original files and exact source lines (provenance.log). Parent toolCallId binds result runId and ownerSessionId. Child session_info binds agent/run/epoch; exact user task and cwd bind assignment; child terminal stop and times precede the native notification. Child directory UUID is not used as run identity. Native completion index is exactly 3 in the projected parent, not the end of the child stream.

focused.log contains the two captured positive replays, 27 attribution/order negatives and the full original auto parent replay. independent.log adds seven independent probes through the existing grader: exact index; completion prose removed; stale notice followed by redispatch; notice before spawn result; malformed terminal child record; full original default parent/stdout/stderr; 1,000 inert parent records retaining exact index 1003. All seven pass. The last probe checks representation/position stability with a larger parent without creating a standing scanner. Helper work is parent scanning plus retained child reads per notification; no nested child-event lifecycle expansion.

AC-2: PASS locally. Actual auto captured report with reconstructed committed Git end state fails `entity has no gates record`; uncommitted report fails the existing durable-commit owner. Reconstruction is explicit: retained CI artifacts do not include the original Git repository. Historical Pi status and wait, Claude replay, Codex correlated lifecycle and cross-host gate bypass controls pass in focused.log. Completion observation therefore changes attribution without declaring the recorded auto journey green.

Three deliberate detached Go overlays prove the existing tests are falsifiable; candidate bytes stayed unchanged. Removing identity checks makes the wrong-run negative fail with `invalid completion passed`; bypassing absent-gate validation fails with `captured missing gate must still fail: <nil>`; bypassing the Git check fails with `uncommitted report passed: <nil>`. Each mutation run exits 1 as intended (identity-mutation.log, gate-mutation.log, commit-mutation.log). These are test-strength experiments, not candidate failures.

AC-3: PARTIAL / deferred final tip host CI. Reused actual final implementation normal.log and race.log, each exit 1 solely for TestCodexResolveManifestAgainstInstalledHost at codex_resolve_test.go:44, matching the explicitly declined pre-existing installed-plugin resolver condition. ensigncycle passes in 354.805s normal and 335.058s race, with no race diagnostic. No new broad rerun was needed after an unchanged diff. Implementation ran required gofmt -w ./cmd ./internal; independent gofmt -l on all changed Go files and git diff --check are clean. Implementation compile-only live-tag run selected no tests. No native/model run or CI was invoked here; skipped/cancelled/not-run checks are not called green.

## Commands

Focused reproduction (exit 0, 12 top-level tests, no skips):
`SPACEDOCK_PI_AC_ARTIFACT_DIR=/tmp/spacedock-tip-ci/pi/live-artifacts/pi/pi-common/auto-continue-after-implementation--auto-continue/single-root go test ./internal/ensigncycle -run 'TestPiNativeCompletion|TestPiAutoContinueReplayDoubleDispatch|TestImplementationLifecycleAndObserverNegativeControls|TestCodexNativeLifecycle|TestAutoContinue.*Replay|TestAutoContinueBypass' -count=1 -v`

Independent probes use Go -overlay to append independent-probes.txt to the existing pi_auto_continue_double_dispatch_replay_test.go owner, then `go test -overlay <mapping.json> ./internal/ensigncycle -run TestIndependentNativeBoundary -count=1 -v` (exit 0). Temporary overlay source/mapping files were inside the assigned worktree and removed after the runs; the candidate was never edited. Mutation overlays replace only the indicated helper predicates/check return in a detached copy, run the existing focused test owner and produce the retained expected red logs.

## Findings and delivery

No new material, deferred-risk or polish finding. Known installed-manifest defect remains declined, with red broad results retained. Final native/host tip CI is an explicitly deferred AC obligation, not a risk or a passed check. Local delivery review may proceed under the dispatch's allowance; release/CI proof remains FO-owned. Checklist: 3 DONE, 1 SKIPPED (final tip CI), 0 FAILED.
