# Isolated diagnosis — 2026-09-15

Candidate unchanged at +217 net lines/four files. All peer broad suites were stopped per FO. These commands ran serially in the assigned worktree; no code/state edits, commits, live runs, broad reruns or CI.

## quiet-race

`go test ./internal/ensigncycle -race -run '^TestCodexProcess(ActivityResetsQuietBudget|QuietTimeoutPreservesFaultEvidence|RecognizesTerminalTurnBeforeOSExit|RequiresFinalMessageForTerminalTurn)$' -count=1 -parallel=1 -v`

Exit 0; package duration 4.304s. Log: `/tmp/same-stage-isolated-quiet-race.log`.

## overlap-normal

`go test ./internal/ensigncycle -run '^TestDurableKeepMovingRequiresOverlappingJourneys$' -count=1 -parallel=1 -v`

Exit 0; package duration 7.565s. Log: `/tmp/same-stage-isolated-overlap-normal.log`.

## durable-race

`go test ./internal/ensigncycle -race -run '^TestDurableTaskJourneys$' -count=1 -parallel=1 -v`

Exit 0; package duration 42.319s. Log: `/tmp/same-stage-isolated-durable-race.log`.

## freshbox-normal

`go test ./internal/cli -run '^TestFreshBoxInstallSucceeds$' -count=1 -parallel=1 -v`

Exit 0; package duration 2.141s. Log: `/tmp/same-stage-isolated-freshbox-normal.log`.

## statecommit-race

`go test ./internal/cli -race -run '^TestStateCommitRefusesDirtyArchivedEntityBeforePublication$' -count=1 -parallel=1 -v`

Exit 0; package duration 5.897s. Log: `/tmp/same-stage-isolated-statecommit-race.log`.

## resolver-normal

`go test ./internal/cli -run '^TestCodexResolveManifestAgainstInstalledHost$' -count=1 -parallel=1 -v`

Exit 1; package duration 0.762s. Log: `/tmp/same-stage-isolated-resolver-normal.log`.

## Interpretation

All four newly failed 250ms process-budget checks and all four tests active at cumulative package timeouts pass unchanged in isolated serial execution. The known installed-host resolver alone still fails with the same missing spacedock@spacedock versus spacedock-local manifest mismatch.

The prior broad suite timeouts were cumulative ten-minute package deadlines: FreshBox was active only16s and StateCommit10s; overlap86s and durable task journeys109s. The normal timeout stacks show process waits in host installation and git-backed durable-blob resolution. In isolation overlap takes7.32s, fresh-box1.84s, full durable journeys41.11s and state-commit4.62s. These results support resource contention/scheduling delay as the cause, without claiming an exclusive root cause or a green full suite. No data race was reported; the four race-mode failures were timing assertions.

All processes stopped. Registry-policy/CI-scope decision remains pending; no candidate edit is justified by this diagnosis.
