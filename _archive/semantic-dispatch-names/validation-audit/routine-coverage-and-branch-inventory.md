# Routine coverage and exact branch inventory

Captain policy correction requires routine selection of every native proof unless explicitly exempted. FO authorized native naming promotion, deterministic fixture extraction, and the separately reproduced exact-branch membership correction. The final scope is presented explicitly; the earlier +450 figure was a promotion estimate, not a renewed cap for this amended task.

Candidate: `844fad458`. Owned branch: `spacedock-ensign/semantic-dispatch-names`. Parent: `bf64ebeca949667c2d60e40a9c595ebdf62375b2`. The local rebase encountered only the expected registry conflict. It retained the lower layer's six same-stage variants, strict registry/selector checker, and all three explicit experiments. Naming's obsolete targeted-only entries were removed.

## Actual scope

The task layer has 55 files, +817/-354 (+463 net) against the moved parent. This correction touches five paths, +167/-138 (+29 net):

- `docs/runtime-live-ci-registry.md`: register Codex-specific native handle semantics in the existing runtime-proof category and remove targeted-only entries.
- `internal/ensigncycle/scheduled_live_test.go`: select `TestLiveSemanticNamesCodex` through the existing Codex scheduled lane.
- `internal/ensigncycle/semantic_names_live_test.go`: add proof/fixture bindings and extract the same fixture setup, preserving native name/handle, raw trace, negative controls, two-commit, frontmatter and clean-state assertions.
- `internal/ensigncycle/conflict_owner_handoff_live_test.go`: share offline deterministic helpers and call the exact branch-membership predicate.
- `internal/ensigncycle/conflict_owner_fixture_test.go`: move the stamped-owner regression and required deterministic helpers into the default suite, using the existing checkout-binary builder; test exact branch membership.

The new offline file and scheduler are the two additional owners beyond the previous 53-file task. Actual growth is two lines below the +465 projection because an unused import and its separator moved out of the live file. No workflow lane, checker framework, generic set abstraction, or new native test was added.

## Exact branch finding and proof

CI run `35136729906`, job `104931339842`, failed the existing handoff grader after 117.66s with branches `["conflict-owner" "main"]`. The grader incorrectly assumed main sorted first, which happened to hold for the historical prefixed branch. The recorded owner is `conflict-owner`; all preceding authority-byte, marker, branch, Git-author, clean-worktree, rebase-abort and worktree-inventory assertions had passed. Downloaded artifacts corroborate rebase abort and the reported `c58ab65` follow-up; their process record says terminal=true and timed_out=false.

Before correction, an overlay created the real stamped fixture, aborted its rebase, committed its marker with Captain identity, and called the actual handoff grader. It exited 1 with the identical ordering error in 3.883s: `/tmp/semantic-branch-order-red.log`.

After correction, that same model-free full-grader replay exited 0 in 0.765s, retaining the observed branch order and all prior assertions: `/tmp/semantic-promotion-handoff-green.log`. The shared predicate requires exactly two identities, main and the recorded owner, in either listing order. The offline proof rejects missing main, missing owner, wrong owner and an extra branch. An overlay restoring the old order-dependent predicate failed the offline regression in 0.993s: `/tmp/semantic-promotion-branch-red.log`.

## Coverage and focused evidence

- Default non-live discovery lists `TestConflictOwnerStampedIdentity`: `/tmp/semantic-promotion-default-list.log`. Its selected default run exited 0 in 1.142s: `/tmp/semantic-promotion-offline-green.log`. It builds the current checkout binary and needs no model or installed runtime binary.
- Removing only the Codex scheduling row made the existing lower-layer checker exit 1 in 0.291s: `TestLiveSemanticNamesCodex selected 0 times in codex-live, want 1`. Log: `/tmp/semantic-promotion-selection-red.log`. The row was restored before subsequent checks.
- Registry reconciliation and existing selection/exemption controls passed together in 0.278s: `/tmp/semantic-promotion-registry-green.log`.
- Existing semantic/legacy and Codex adapter isolation/collision controls passed in 1.349s: `/tmp/semantic-promotion-legacy-green.log`. Legacy compatibility remains an offline claim; no new legacy native run is claimed.
- Live-tag no-test compilation passed in 0.293s: `/tmp/semantic-promotion-live-compile.log`. The first extraction compile identified an unused moved-helper import, which was removed within the authorized extraction.
- Required gofmt ran; unrelated baseline alignment formatting stayed outside the candidate. `git diff --check` passed.

## Required broad checks

Both commands used `SPACEDOCK_BIN=/tmp/semantic-identity-spacedock`; the newly offline owner test independently builds the current checkout binary.

- `go test ./...` exited 1 solely at `TestCodexResolveManifestAgainstInstalledHost`, `internal/cli/codex_resolve_test.go:44`: `spacedock@spacedock` is not installed, but the resolver returned `/Users/clkao/.codex/plugins/cache/spacedock-local/spacedock/0.28.0-pre0/.codex-plugin/plugin.json`. CLI completed in 152.288s; contractlint passed in 0.648s, dispatch in 44.681s, and ensigncycle in 248.525s. All other packages passed. Log: `/tmp/semantic-promotion-full-normal.log`.
- Sequential `go test ./... -race` exited 1 solely at that same test, line, error, and returned path. CLI completed in 160.524s; contractlint passed in 9.892s, dispatch in 62.074s, and ensigncycle in 245.215s. All other packages passed; no data-race diagnostic appeared. Log: `/tmp/semantic-promotion-full-race.log`.

The baseline was verified from each run, not assumed. Required checks were executed and are not wholly green. Final local worktree is clean at `844fad458c738d715994cd46fa043eb22c45f994`, with `bf64ebeca` confirmed as ancestor.

While this race run completed, the lower owner committed a stronger executed-scheduler/YAML coverage checker at `342edfa8903e0e822ec545d2ee45401352eb6c60`. This artifact's checker and broad-suite evidence is against `bf64ebeca`, not that later revision. FO will restack and verify the combined tip; no redundant parent-absorption broad rerun was started.

## Limits

No native/model execution occurred. Routine Codex efficacy remains for final-tip CI; deterministic selection evidence does not substitute for that run. Retained earlier native semantic spawn/reuse evidence is unchanged. The unrelated Pi notification and Claude timeout findings remain outside this naming correction. The known local installed-manifest resolver baseline has an existing FO-authorized DECLINE disposition and cannot waive a new failure.

No push, CI trigger, PR change, or other branch mutation occurred. The rebase touched only the registered naming branch. Approved entity report, gate, frontmatter and independent validation history were preserved.
