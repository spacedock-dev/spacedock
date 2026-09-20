# Final naming consumer corrections

Captain binding `resolution:binding-1789570590063113000` approved `/private/tmp/naming-ci-additional-five-review.md` at `2026-09-16T14:56:30.063116Z`, SHA-256 `b77c9ff570c487310082067298d720bc4c991b5461390f3d0698373e179ac098`. FO authorized the five additional naming corrections. The previous identity/registry correction `5e3819df1` is retained.

Local candidate commit: `9653ae579`. Exactly five additional paths changed, +51/-44 (+7 net):

- `internal/ensigncycle/pi_rejection_extractors_test.go`
- `internal/ensigncycle/pi_rejection_extractors_test_test.go`
- `internal/dispatch/build_codex_host_test.go`
- `internal/dispatch/build_input_mode_test.go`
- `internal/dispatch/self_contained_assignment_test.go`

Final task layer against `c7568815b`: 53 files, +686/-252 (+434 net), within the approved 53-file/450-net ceiling. The six-line deterministic registry entry explains the increase from the earlier consolidated +428 forecast. Final categories: 5 production files, 18 test files, 5 documentation/skill files, and 25 golden fixtures. No additional owners changed.

The Pi extractor retains the entire dispatch filename stem for either semantic or historical names. Path, stage, async run correlation and completion controls remain. The three refusal tests now inspect independent literal canonical artifact paths. None of the assertions was weakened.

## Deterministic red-to-green evidence

- Before the regex fix, `go test ./internal/ensigncycle -run '^TestPiRejectionRoutesFreshChain$' -count=1 -v` exited 1 in 0.339s: legacy chain passed; semantic chain produced zero of eight required events. Log: `/tmp/semantic-five-pi-red.log`.
- After the fix, `go test ./internal/ensigncycle -run '^TestPiRejectionRoutes.*$' -count=1 -v` exited 0 in 0.260s. Both full chains pass with exact handle preservation, missing-completion refusal, wrong-branch refusal, and existing handle/run-ID correlation controls. Log: `/tmp/semantic-five-pi-green.log`.
- A temporary overlay fed the actual retained prompt from `dispatch build --stamp --host pi` to the real extractor: `Read /tmp/spacedock-dispatch/thing-implementation.md and treat its content as your assignment.` `TestPiSemanticGeneratedPromptAudit` now passed in 0.311s; before correction it returned no routes. Log: `/tmp/semantic-five-generated-prompt-green.log`. This is deterministic boundary evidence, not a live Pi run.
- The three corrected refusal tests passed together, exit 0 in 0.743s: `TestBuildCodexHostRejectsBareModeBeforeArtifactCreation`, `TestBuildInputFailuresDoNotReadOrStamp`, and `TestBuildWithoutResolvedLauncherFailsBeforeWritingArtifact`. Log: `/tmp/semantic-five-guards-green.log`.
- Temporary overlays injected writes at each canonical path during those refusal checks. All three tests rejected the mutation, collectively exiting 1 in 0.635s. Diagnostics respectively reported an unexpected bare-Codex artifact, mutated dispatch artifact, and an unresolved launcher writing an artifact. Log: `/tmp/semantic-five-guards-mutation.log`. The corresponding pre-correction overlays passed incorrectly; see `/tmp/semantic-stale-guards-audit.log`. Overlay writes restored prior bytes or removed their temporary artifacts through cleanup; candidate code contains no injections.
- `go test -tags live ./internal/ensigncycle -run '^$'` exited 0 in 0.258s. Log: `/tmp/semantic-five-live-compile.log`.
- `gofmt -w ./cmd ./internal` ran; the unrelated baseline alignment hunk in `internal/release/runtime_live_evidence_workflow_test.go` was kept outside the candidate. `git diff --check` passed.

## Final required checks

Commands used `SPACEDOCK_BIN=/tmp/semantic-identity-spacedock`, built from the owned checkout for the preceding identity boundary. These five changes affect test/proof files only.

- `go test ./...` exited 1 solely at `TestCodexResolveManifestAgainstInstalledHost`, `internal/cli/codex_resolve_test.go:44`: `spacedock@spacedock` is not installed, but the resolver returned `/Users/clkao/.codex/plugins/cache/spacedock-local/spacedock/0.28.0-pre0/.codex-plugin/plugin.json`. CLI took 120.924s; contractlint passed in 1.099s, dispatch in 43.229s and ensigncycle in 200.562s. Every other package passed. Log: `/tmp/semantic-five-full-normal.log`.
- Sequential `go test ./... -race` exited 1 solely at the same test, line, error and returned path. CLI took 156.607s; contractlint passed in 3.696s, dispatch in 61.723s and ensigncycle in 237.871s. Every other package passed. No data-race diagnostic appeared. Log: `/tmp/semantic-five-full-race.log`.

Required checks were executed and are not wholly green. The exact known resolver failure was independently confirmed in each final run. Final code worktree is clean at `9653ae57912118f6c577d9abfa03e3752b6beac2`.

## Limits and remaining unrelated findings

All three approved naming findings are addressed: adapter identity inference, Pi filename decoding, and obsolete artifact assertions. The separate Pi notification-grader findings and Claude timeout findings remain outside this correction; no claim is made that these changes repair them. The known installed-manifest resolver baseline retains its FO-authorized DECLINE disposition, which does not waive any new failure.

No native/model run, push, rebase, CI or PR edit occurred. Approved entity report, gate and frontmatter were preserved.
