# Runtime live owner registry

Human-readable index of the owner bindings on the live journeys in
`internal/ensigncycle/shared_live_runner_test.go`.

The bindings in code (`liveXFail(target, owner)` on each `liveJourney` call) are
the enforced record. This table is a convenience copy for humans and is **not
machine-checked** — no test reads it. Change the code first, then keep this table
in sync by hand.

Each binding names the single owner responsible for repairing a red or flake on
that runtime target. Removing a binding claims the journey is green there.

| Journey | Owner | Observed failure or flake reason |
|---|---|---|
| auto-continue-after-implementation | s0gq9p69nztejw8xp3by4k7f | pi red `validation-worker-not-dispatched` |
| default-headless-gate-stop | penfp034pt9s3cgwp7wg3ykk | pi reds `gate-hold-violation`, `gate-not-held`, `implementation-worker-not-dispatched` |
| ac-value-reanchor | psvqjf0w8xh2txp9604gsvmz | pi flake; passes and fails between runs, previously unowned |
| keep-moving-posture | psvqjf0w8xh2txp9604gsvmz | pi flake; passes and fails between runs |
| keep-moving-posture | 060xp69y61yhrww23g3wvwqy | claude-sonnet flake on `claude-live` |
| owned-conflict-owner-handoff | psvqjf0w8xh2txp9604gsvmz | pi flake; flips between XFAIL and XPASS, no single-repair owner |
| rejection-flow | 6h3teccccn3qh71yqcmjbjx4 | pi rejection-worker topology fault |
| rejection-flow | jfdrr2dbdtt952s4cvvy74k8 | claude-sonnet red `rejection-worker-topology` |
| auto-continue-after-implementation | jfdrr2dbdtt952s4cvvy74k8 | claude-sonnet red `validation-worker-not-dispatched` |
| default-headless-gate-stop | jfdrr2dbdtt952s4cvvy74k8 | claude-sonnet red `implementation-worker-not-dispatched` |
| smallest-sufficient-mechanism | jfdrr2dbdtt952s4cvvy74k8 | claude-sonnet red `smallest-mechanism-violation` |
| keep-moving-posture | unowned (blank; pending analysis) | claude-haiku red `keep-moving-violation` |
| default-headless-gate-stop | unowned (blank; pending analysis) | claude-haiku reds `gate-hold-violation`, `gate-not-held`, `implementation-worker-not-dispatched` |
| withdrawn-gate-recovery | unowned (blank; pending analysis) | claude-haiku red `withdrawn-recovery-boundary` |
| filing | unowned (blank; pending analysis) | claude-haiku red `filing-command-not-observed` |
| break-glass selected-bare (`dispatch_recovery_live_test.go`) | unowned (blank; pending analysis) | claude-haiku red `break-glass-observables` |
