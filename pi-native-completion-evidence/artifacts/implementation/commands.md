# Local validation commands

All commands executed in the registered `spacedock-ensign/pi-native-completion-evidence` worktree. Broad normal and race runs were sequential, including the final reruns after the authorized conflicting-epoch correction. The initial broad results precede that small correction; final results are `normal.log` and `race.log`.

```sh
SPACEDOCK_PI_AC_ARTIFACT_DIR=/tmp/spacedock-tip-ci/pi/live-artifacts/pi/pi-common/auto-continue-after-implementation--auto-continue/single-root go test ./internal/ensigncycle -run 'TestPiNativeCompletion|TestPiAutoContinueReplayDoubleDispatch|TestImplementationLifecycleAndObserverNegativeControls|TestCodexNativeLifecycle|TestAutoContinue.*Replay|TestAutoContinueBypass' -count=1 -v
go test -tags live ./internal/ensigncycle -run '^$'
gofmt -w ./cmd ./internal
go test ./...
go test ./... -race
```

The live-tag command only compiles: it selects no tests and invokes no model. The new captured replay and 27 negatives are hermetic under ordinary package tests. The optional existing artifact replay additionally exercised the full downloaded original auto parent and child evidence. Provenance beside the committed fixtures records source hashes and record projection; no downloaded Git repository was available, so the auto end-state/Git checks explicitly reconstruct durable state with the captured report.

Formatting was rerun on the final two touched files after the small follow-up correction. The broad formatter’s unrelated existing spacing change in `internal/release/runtime_live_evidence_workflow_test.go` was restored, keeping the task scope unchanged. `git diff --check` passed.

Both broad commands retain the anticipated local `TestCodexResolveManifestAgainstInstalledHost` failure, not an expected-failure pass. Its installed plugin resolver points to the local cached plugin even though the test’s marketplace name is absent. No native/model, network-auth, PR, push, CI, or rebase was performed by this worker. Final host CI remains deferred to the stack tip after independent validation.
