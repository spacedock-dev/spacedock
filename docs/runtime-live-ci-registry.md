# Runtime live owner registry

Human-readable index of the owner bindings on the live journeys in
`internal/ensigncycle/shared_live_runner_test.go`.

The bindings in code (`liveXFail(target, owner)` on each `liveJourney` call, or
in a gap function such as `breakGlassSelectedBareGap`) are the enforced record.
`TestRuntimeLiveRegistryMatchesBindings` in `internal/contractlint` reads them and
fails when a binding has no row here or a row has no binding. Change the code and
this table together.

Each binding names the single owner responsible for repairing a red or flake on
that runtime target. Removing a binding claims the journey is green there.

| Journey | Target | Owner | Observed failure or flake reason |
|---|---|---|---|
| auto-continue-after-implementation | pi | s0gq9p69nztejw8xp3by4k7f | pi red `validation-worker-not-dispatched` |
| default-headless-gate-stop | pi | penfp034pt9s3cgwp7wg3ykk | pi reds `gate-hold-violation`, `gate-not-held`, `implementation-worker-not-dispatched` |
| ac-value-reanchor | pi | psvqjf0w8xh2txp9604gsvmz | pi flake; passes and fails between runs, previously unowned |
| keep-moving-posture | pi | psvqjf0w8xh2txp9604gsvmz | pi flake; passes and fails between runs |
| keep-moving-posture | claude-sonnet | 060xp69y61yhrww23g3wvwqy | claude-sonnet flake on `claude-live`: the FO hand batch-commits before `--stamp` (Sonnet 5, run 31894934349), or the host stops the background `questioned` worker before it reports |
| owned-conflict-owner-handoff | pi | psvqjf0w8xh2txp9604gsvmz | pi flake; flips between XFAIL and XPASS, no single-repair owner |
| rejection-flow | pi | 6h3teccccn3qh71yqcmjbjx4 | pi rejection-worker topology fault |
| auto-continue-after-implementation | claude-sonnet | jfdrr2dbdtt952s4cvvy74k8 | claude-sonnet red `validation-worker-not-dispatched` |
| keep-moving-posture | claude-haiku | unowned (blank; pending analysis) | claude-haiku red `keep-moving-violation`; the host can also stop a background worker 600 s after a headless FO ends its turn (any model) |
| default-headless-gate-stop | claude-haiku | unowned (blank; pending analysis) | claude-haiku reds `gate-hold-violation`, `gate-not-held`, `implementation-worker-not-dispatched` |
| withdrawn-gate-recovery | claude-haiku | unowned (blank; pending analysis) | claude-haiku red `withdrawn-recovery-boundary` |
| `breakGlassSelectedBareGap` (break-glass selected-bare, `dispatch_recovery_live_test.go`) | claude-sonnet | 060xp69y61yhrww23g3wvwqy | claude-sonnet flake: the `Agent` call does not stay in bare mode |
| `breakGlassSelectedBareGap` (break-glass selected-bare, `dispatch_recovery_live_test.go`) | claude-haiku | unowned (blank; pending analysis) | claude-haiku red `break-glass-observables` |
