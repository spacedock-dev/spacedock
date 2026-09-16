# Stamped owner identity correction

Captain binding `resolution:binding-1789567773656214000` approved the two-file identity proposal at `2026-09-16T14:09:33.656217Z`, pinned to `/private/tmp/naming-ci-identity-scope-review.md`, SHA-256 `6c98b8b64dd35ce7ee2a32927037c7e9a2d147a7918ac591d61a22c8c8d5c35e`. FO separately authorized registration of its deterministic check in the already-counted registry after the normal suite found that omission.

Local code commit: `5e3819df1`. The correction touches three paths, +22/-23 (-1 net): `internal/dispatch/codex_v2_adapter.go`, `internal/ensigncycle/conflict_owner_handoff_live_test.go`, and `docs/runtime-live-ci-registry.md`. Only the first two are new to the task manifest. Final task layer against `c7568815b`: 48 files, +635/-208 (+427 net), within the authorized 48-file/450-net ceiling.

The adapter retains the supplied worker name, isolated spawn arguments, and sanitized-name collision refusal. It no longer claims to reverse the worker name into entity or stage. The conflict-owner fixture reads entity and status from the canonical stamped entity; its tuple, Git registration, and provenance assertions remain. A deterministic regression exercises real stamping and adapter parsing with literal expected entity, stage, worker name and branch. Its registry entry requires exact local selection and makes no native-runtime or CI-lane claim.

## Red-to-green boundary evidence

The test harness first encountered an inherited nonexistent `SPACEDOCK_BIN`; removing that variable exposed an older PATH binary. Those attempts were setup diagnostics, not candidate proof. `go build -o /tmp/semantic-identity-spacedock ./cmd/spacedock` then built the owned checkout. Acceptance commands below explicitly set `SPACEDOCK_BIN=/tmp/semantic-identity-spacedock` where applicable.

- Before the correction, `go test -tags live ./internal/ensigncycle -run '^TestConflictOwnerStampedIdentity$' -count=1` exited 1 in 0.631s. The actual stamped tuple had empty Entity, Stage `implementation`, WorkerName `conflict-owner-implementation`, Branch `conflict-owner`, and the recorded worktree. Log: `/tmp/semantic-identity-approved-red.log`.
- After the correction, the same test plus `TestCodexNativeLifecycleUsesCorrelatedSessionHandle` and `TestCodexNativeLifecycleParentRolloutLookupFailsClosed` exited 0 in 0.509s. Log: `/tmp/semantic-identity-approved-green.log`.
- `go test ./internal/dispatch -run '^(TestCodexMultiAgentV2.*|TestSemantic.*)$' -count=1` exited 0 in 1.325s. This retains isolation, collision and semantic identity controls. Log: `/tmp/semantic-identity-adapter-green.log`.
- `go test -tags live ./internal/ensigncycle -run '^$'` exited 0 in 0.209s, without running a model. Log: `/tmp/semantic-identity-live-compile.log`.
- `gofmt -w ./cmd ./internal` ran. Its unrelated pre-existing alignment-only change in `internal/release/runtime_live_evidence_workflow_test.go` was removed from this candidate. `git diff --check` passed.

## Required suite results

The first `go test ./...` exited 1. It reproduced the known resolver baseline and found a new failure: `TestRuntimeLiveRegistryReconciliation` at line 241 rejected unregistered `TestConflictOwnerStampedIdentity`. All other packages passed, including ensigncycle in 229.611s. Log: `/tmp/semantic-identity-full-normal.log`. This new failure was reported before editing. FO authorized an entry under Targeted implementation proofs in the existing registry file. The focused reconciliation then exited 0 in 0.206s: `/tmp/semantic-identity-registry-green.log`.

The final `go test ./...` exited 1 solely at `TestCodexResolveManifestAgainstInstalledHost`, `internal/cli/codex_resolve_test.go:44`: `spacedock@spacedock` is not installed, but the resolver returns `/Users/clkao/.codex/plugins/cache/spacedock-local/spacedock/0.28.0-pre0/.codex-plugin/plugin.json`. CLI completed in 125.952s; contractlint passed in 0.752s and ensigncycle reused its passing cache. Log: `/tmp/semantic-identity-final-normal.log`.

The sequential final `go test ./... -race` exited 1 solely at the same `TestCodexResolveManifestAgainstInstalledHost`, line 44, with the same installed-manifest error and returned path quoted above. CLI completed in 178.948s, contractlint passed in 6.464s, dispatch in 84.444s and ensigncycle in 256.363s. Every other package passed; no data-race diagnostic appeared. Log: `/tmp/semantic-identity-final-race.log`. The final code worktree is clean at `5e3819df1c545ce76f986eeccc9fc656b07c2725`.

## Remaining held work

The five additional corrections in `naming-ci-consumer-proposal.md` remain HOLD: the Pi semantic dispatch-path extractor and its proof, plus three obsolete-path refusal assertions. They were not edited. This identity correction does not complete those findings or repair the separate Pi notification graders.

The known resolver baseline has the existing FO-authorized DECLINE disposition. Required checks are not wholly green; no new failure is waived by that disposition. No native/model run, push, CI, rebase or PR edit occurred. Approved entity report, gate and frontmatter were preserved.
