# Rebased naming validation

Exact tested head: `7535aad706e1d0ae88cad8bd45f5ebf4eadf96b8`, lower layer `342edfa8903e0e822ec545d2ee45401352eb6c60`. Naming scope: 55 files, +817/-354 (+463 net). Source worktree remained clean; no candidate edits, model calls, broad suites, push or CI trigger.

## Current local checks

- `go test ./internal/contractlint -run '^TestRuntimeLiveRegistryReconciliation$' -count=1 -v`: passes against the actual lower-layer checker, including executed offline Claude/Codex scheduler variants (1.521s). This verifies the naming row under the real checker, not the earlier checker revision. The three explicit experiment entries are byte-identical to the lower layer. The independent #801 audit is reused rather than repeated.
- `go test ./internal/ensigncycle -run '^TestConflictOwnerStampedIdentity$' -count=1 -v`: passes without the live tag (1.772s), builds the checkout binary and checks exact semantic entity, stage, worker and branch plus missing/wrong/extra inventory negatives.
- `go test ./internal/dispatch -run '^TestSemantic|Test.*Codex.*' -count=1`: passes (3.938s), preserving semantic naming, long legacy caps, ambiguity refusal, equal-cycle cohorts and Codex adapter identity checks.
- Independent Go overlay `audit_test.go.txt`, selected alone with the live tag, calls the actual `writeConflictOwnerFixture` and `assertConflictOwnerHandoff`; no runner or model is invoked. It builds the current binary, aborts the fixture rebase, writes and commits the marker as Captain, then evaluates the complete existing grader. Semantic and legacy registered branches pass; wrong/missing main and extra branches fail at the branch-inventory boundary (5.654s). Separate exact-predicate controls reject missing owner, empty and duplicate inventories too.
- Restoring the historical order-dependent predicate only through an overlay makes that same complete semantic grader fail on `["conflict-owner" "main"]` (1.765s). This reproduces the CI false failure without weakening authority bytes, marker, author, clean worktree, rebase-abort, worktree count or fresh-envelope assertions.
- `git diff --check 342edfa..HEAD` passes. Implementation formatting evidence is retained; validation made no source edits.

To reproduce the independent grader, map a new `internal/ensigncycle/validation_cycle2_test.go` path through Go's `-overlay` JSON to a copy of `audit_test.go.txt`, then run `go test -overlay <overlay.json> -tags live ./internal/ensigncycle -run '^TestValidationRealHandoffGrader$' -count=1 -v`. The test sets SPACEDOCK_BIN to its own freshly built checkout binary. For the negative control, also overlay conflict_owner_fixture_test.go with exactOwnerBranches changed to require main at index zero.

## Acceptance and limits

AC-1: prior exact public-name/registered-branch/assignment audit is preserved; current semantic and Codex checks pass. AC-2: current collision, long-name, stage-boundary and refusal controls pass; prior independent size/Unicode/1,000-entity matrix remains available. AC-3: current legacy/semantic cohort controls pass; this replay independently establishes legacy registered branch preservation, authoritative stamped identity and full-grader strict membership. AC-4: retained earlier native one-spawn/same-handle reuse and two marker commits remain historical proof. Original native assertions are byte-identical to the prior candidate except the retained artifact-directory argument; the fixture extraction and scheduling declaration do not weaken them.

The earlier targeted-only policy is superseded: native semantic naming is selected by the existing Codex scheduled lane, and deterministic owner identity is part of the default offline suite. Routine selection is proven locally; current native execution is not.

Prior full normal/race evidence belongs to `844fad458` on `bf64ebeca`, not this rebased candidate. Both runs fail solely at TestCodexResolveManifestAgainstInstalledHost, internal/cli/codex_resolve_test.go:44, under the existing FO-authorized baseline DECLINE; they are not wholly green. Final combined-tip suites and both-host live CI remain pending and separately owned. Prior run 35136729906 supplies the branch-order false-failure evidence, not a passing final-tip result. No new material finding; local validation recommends PASSED with these delivery checks still pending.
