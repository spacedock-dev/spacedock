# First-outage recovery implementation

FO authorized the exact narrow `dispatch name` interface and seven existing paths, estimated +160–190 net lines. Code commit `bfc97adbd` on baseline `f5af7282dda9151771ff9888dd246553e38afa82` changes seven files, +186/-4 (+182 net). No new owner, controller, framework, discovery surface or state mechanism was added.

The read-only route reuses `validateWorkerName`; it requires readable inputs and a declared stage, rejects worktree copies and unsafe/ambiguous identities, and emits only the canonical base. Selected-team recovery queries it once when no validated same-assignment envelope exists. Selected mode, suffix/occupancy refusal, exactly-one-worker and committed-report obligations remain intact. Complete executable absence remains an unavailable dependency, not repaired or tested here.

## Exact verification

1. Before implementation, `go test ./internal/dispatch ./skills/integration -run '^(TestDispatchName.*|TestFirstOutageCanonicalName)$' -count=1` exited 1: the new positive query and actual outage shim failed because `name` was an unknown dispatch subcommand (`red.log`).
2. Initial green: `go test ./internal/dispatch ./skills/integration -run '^(TestDispatchName.*|TestFirstOutageCanonicalName|TestSemantic.*|TestCodexMultiAgentV2.*)$' -count=1` exited 0 (`green.log`).
3. After strengthening no-mutation assertions, `go test ./internal/dispatch ./skills/integration ./internal/ensigncycle -run '^(TestDispatchName.*|TestFirstOutageCanonicalName|TestSemantic.*|TestCodexMultiAgentV2.*|TestAssertBreakGlass.*|TestAssertCompleteRecoveryReport.*|TestDispatchRecoveryPromptsSelectIntendedDispatchModes)$' -count=1` exited 0: dispatch 10.877s, integration 2.347s, ensigncycle 3.178s (`focused.log`).
4. Detached Go overlay replacing the emitted name with `wrong-` plus the canonical name: `go test -overlay=/tmp/semantic-first-outage-overlays/wrong-name.json ./internal/dispatch -run '^TestDispatchNameReadOnly/ci-duration-hints$' -count=1` exited 1 with `wrong-ci-duration-hints-backlog` (`wrong-name.log`).
5. Detached Go overlay replacing only the query's `validateWorkerName` call with `semanticName(status.EntitySlug(entityPath), *stage)`: `go test -overlay=/tmp/semantic-first-outage-overlays/collision.json ./internal/dispatch -run '^TestDispatchNameRefusesUnsafeIdentity/ambiguous$' -count=1` exited 1, rejecting the mutant's successful ambiguous identity (`collision.log`). Candidate files stayed unchanged during both overlays.
6. `gofmt -w internal/dispatch/dispatch.go internal/dispatch/names.go internal/dispatch/names_test.go skills/integration/dispatch_test.go` and `git diff --check` completed successfully. Full-tree formatting remains assigned to final combined verification.

The first-outage fixture builds the real checkout binary, fails both build attempts through the actual forwarding shim, and successfully queries the uncached literal identity between them. It checks entity bytes, checkout status, commit and dispatch artifact preservation. Missing inputs, invalid flags, unknown stage, worktree copy, invalid component, insufficient budget and cross-generation ambiguity refuse without output or mutation. Existing long-name, legacy, suffix and collision controls remain unchanged.

## Limits and recommendation

Offline identity recovery and existing mode/duplicate-worker/committed-report grader controls pass. This does not claim a new native worker run or live committed report. The retained live first-outage failure is not waived; independent validation and final-tip host CI must prove the corrected runtime behavior. Final combined normal/race/full-tree formatting are pending under FO ownership. No broad run, model run, CI, push, rebase or new agent occurred. Prior resolver baseline and separate runtime/root findings remain visible; no new failure is exempted.

Recommend independent validation of this bounded correction, then final combined checks and native CI. Both dispatched implementation items are complete; native efficacy remains a downstream verification obligation.
