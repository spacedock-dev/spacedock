# P1/P2 focused correction verification

FO FIX authorization for both Material task-owned AC-1 observer findings was supplied in the advance dispatch. Scope: existing native Pi correlation helper, route extractor and two live callers, required focused tests and captured fixtures. P1/P2 findings and original durable bundle/shim evidence remain in `../../../validation/tip35149242496/`.

Source head: d0a6f413af32fa69ac9da4af2c74c0086fcca812. CI merge: 3f23437d91c98bcf4cba93e52d0f1da2761d2cb1, parents 85bd948e626f3522869fb577d8139330f73f176b and d0a6f413af32fa69ac9da4af2c74c0086fcca812. Both source and tested merge have full tree 717d412d4b2b0f2ede3246c6f25f8a49d8d8d252; local source tree was checked again with git rev-parse. No assumption that a merge parent is main. Original record SHA-256 hashes and projection/line offsets accompany the committed fixtures in tip35149242496-provenance.json.

Before helper edits:

```sh
go test ./internal/ensigncycle -run 'TestPiNative(SynchronousCapturedReplay|SameStageCapturedRoutes)$' -count=1 -v
```

Both failed: lifecycle `spawns=1 completed=-1 validation=73`; routes only spawn112 and existing obligation check reports an unfinished round. Exact output: red.log.

Final commands, normal then race sequentially:

```sh
go test ./internal/ensigncycle -run 'TestPiNative|TestPiRejection|TestPiRecorded|TestImplementationLifecycleAndObserverNegativeControls|TestCodexNativeLifecycle|TestAutoContinue|TestSameStage|TestAssertRecordedGateHoldLog' -count=1 -v
go test ./internal/ensigncycle -race -run 'TestPiNative|TestPiRejection|TestPiRecorded|TestImplementationLifecycleAndObserverNegativeControls|TestCodexNativeLifecycle|TestAutoContinue|TestSameStage|TestAssertRecordedGateHoldLog' -count=1 -v
go test -tags live ./internal/ensigncycle -run '^$'
gofmt -w ./cmd ./internal
git diff --check
```

Focused normal/race pass in 5.146s / 7.647s, no race diagnostics. The live-tag command passes with no tests selected (compile only). The exact broad formatting command ran; its two-spacing change in the pre-existing internal/release/runtime_live_evidence_workflow_test.go baseline was inspected and restored from HEAD. Final touched files were formatted again; diff check passes and no unrelated formatting is committed.

All 27 previous native negatives, 21 synchronous negative controls and nine same-stage routing cases run. The two-worker positive uses a separately identified fresh reviewer; stale/wrong calls, reused run, prior worker notice, missing/error child and post-gate completion reject. Tests still reject absent gate and uncommitted report through their existing owners. No duplicate commit parser or policy changes were introduced. Source default has no Git bundle; same-stage bundle verification remains the independent retained triage evidence, not a claimed new native run.

Full normal/race commands are explicitly deferred to the FO's final combined tip with the #801 identity correction. All six same-stage native variants require final tip CI: plain previously failed; review-required, separate-review-required, round-required, round-missing and cycle-limit were unstarted, not skipped or passed. No model/native run, broad suite, push, rebase, CI, or #801 worktree operation occurred during this correction.
