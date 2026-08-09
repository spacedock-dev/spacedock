# Gate review: Restore the shared ensign contract at the Codex fresh-dispatch boundary — validation cycle 3

## Chosen direction

Merge unchanged candidate `972f7eabf42a026486c5156e8f831c57754aa484` with current `main` `5efa7d31d` (kd); no rebase is needed.

## Evidence

- Checklist: 3 DONE, 0 SKIPPED, 0 FAILED; see the validation cycle 3 report in the entity file at lines 366-379.
- AC-1 is covered by the exact-head Codex `full-ensign-cycle` journey; AC-2 and AC-3 by the dispatch/bootstrap, process, recorded-gate, and adversarial checks; AC-4 by the cross-host, full, race, and live checks.
- PR #642 run `31296391685` passed offline, Codex, and Claude at the exact head. The auxiliary `journey-delta-comment` job failed only while downloading the Codex artifact after five retries; it is not a required branch-protection check, and no rerun is needed for the runtime evidence.
- `git merge-tree --write-tree origin/main HEAD` produced clean tree `782e44f0b573610e5021e95d42d342536cd13d6e`. The effective merge keeps kd and adds only nv's 16-file patch. Full `go test ./...` and `go test ./... -race` passed on that synthetic tree.

## Recommendation

APPROVE under the delegated Captain conn; invoke the normal terminal merge guard for PR #642, then archive nv.
