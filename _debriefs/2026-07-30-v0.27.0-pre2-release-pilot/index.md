---
record-type: release-pilot
release: v0.27.0-pre2
date: 2026-07-30
release-commit: 5e7f1ffa08721c062fa9ae82636549a635983e95
base-commit: 4eeb94e9b1f7d2e407961e28c941e422c28749fc
release-run: https://github.com/spacedock-dev/spacedock/actions/runs/30552325976
source-record: source-task-record.md
---

# Spacedock v0.27.0-pre2 release pilot

This is an operational release record, not a development-workflow task. The
release was initially driven through a task-shaped record; that source record
and its generated gate rooms are retained here as reference artifacts so the
mistaken classification does not erase evidence.

## Outcome

`v0.27.0-pre2` was published from release commit
`5e7f1ffa08721c062fa9ae82636549a635983e95`, whose sole parent is
`4eeb94e9b1f7d2e407961e28c941e422c28749fc`. The release commit changes only the
two plugin manifests and the First Officer shared-core version stamp.

The release run completed its waiver gate, goreleaser, edge reconciliation, and
journey-ledger jobs successfully. The one-shot
`SPACEDOCK_E2E_GATE_WAIVER` repository variable was removed and verified
absent after the release gate consumed it.

All published archives matched `checksums.txt`. `origin/next`, the edge plugin
manifests, marketplace calendar key, and remote `spacedock@next` cask advanced
to the pre2 release.

## Provider-free journey

A checksum-verified published arm64 edge archive installed Spacedock and its
Codex plugin at `0.27.0-pre2` into a fresh authenticated `CODEX_HOME`, with
`SPACEDOCK_BIN` unset and no checkout/plugin override.

The retained standalone pilot repository is
`/tmp/spacedock-pre2-pilot.YHcWXZ`. Its history records baseline, Briefing
binding, exact `person:captain` decision, and cold-session consumption. The
second consume failed without changing the entity bytes, proving one-use
authority.

## Findings

- **Sprint-critical:** after consumption, pre2 entered `implementation` but
  projected the later `review` stage and could not build an implementation
  dispatch package without state repair. Owner:
  `dispatch-entered-stage-after-gate-consume` (`gqs`).
- **Deferrable evidence friction:** the Codex sandbox could not open Homebrew's
  host lock path under `/opt/homebrew`; the exact published archive supplied the
  install proof instead. Recheck the cask on an ordinary writable Homebrew host.
- **Deferrable setup friction:** an empty `CODEX_HOME` requires authentication
  bootstrap; the pilot reused only existing authentication.
- **Deferrable fixture friction:** the disposable inline workflow required an
  explicit workflow directory, and `state commit` correctly had no split-root
  state checkout to commit.

## Reference artifacts

- [Original task-shaped record](source-task-record.md)
- [Backlog gate room](review/backlog/briefing-1/)
- [Ideation gate room](review/ideation/briefing-1/)
