---
title: "An FO's contract-file find-hunt can escape the test fixture and read real, unrelated repos on the host machine"
status: backlog
source: "Surfaced 2026-07-09/10 during entity dp's (fo-deferred-load-point-hunt-vs-skill-addressing) cycle-2 validation, local live AC-1 re-measurement on this development machine. 2 of 4 find-hunts that cycle escalated past a scoped filesystem sweep into a genuine wrong-root wander: an unscoped `find /` (attempt-1 filing scenario) and a similarly unscoped `find / -maxdepth 6 -iname \"*spacedock*\"` (attempt-2 smallest-sufficient-mechanism scenario) both succeeded on this host, discovering and then reading real, unrelated repos -- including this very validation's own live checkout, /Users/clkao/git/spacedock-research/spacedock-v1 (read docs/dev/README.md, grepped internal/status/handlers.go for verdict handling), plus a sibling /Users/clkao/git/spacedock checkout. Both reads were read-only; no writes or commits touched either real repo. Captain directed this be filed as its own entity rather than folded into dp or noted only in a ledger, given the safety implications of a live FO reading arbitrary real repos on a host machine during what should be a sandboxed test fixture run."
started:
completed:
verdict:
score: 0.6
worktree:
issue:
id: pjm3ge4jan4a5gqd84xef3g3
---

A live FO's contract-file find-hunt (the same underlying anti-pattern `sc5`/`live-runner-boot-preamble-hardening` and `dp`/`fo-deferred-load-point-hunt-vs-skill-addressing` document extensively) is not always confined to the plugin/skill directory or the test fixture root. When the FO issues an UNSCOPED filesystem sweep (`find /`, or `find / -maxdepth N -iname "..."` with no path restriction to the fixture or plugin dir), that sweep can succeed on a real development machine — where other, real, unrelated git repositories exist on disk — and the FO then reads real files from those repos to continue its (mistaken) search for a contract reference. On a clean CI runner this would very likely find nothing interesting and terminate quickly; on a developer's local machine, with other real projects checked out, it can return plausible-looking hits and lure the FO into extensively exploring a repo that has nothing to do with the fixture under test.

This is distinct from `wrong-root-detector-symlink-false-positive` (5q): that entity is about a detector FALSE POSITIVE (flagging a wander that did not happen). This is the opposite and more concerning case — a REAL wander that the existing `detectBroadSearchAtBoot`/`detectWrongRootBoot` classifiers may or may not currently catch cleanly when the target is a real, existing, unrelated repo rather than an obviously-wrong path (worth checking as part of ideation: does the detector correctly flag "read real_unrelated_repo/docs/dev/README.md" as a wrong-root wander today, or does it only fire on specific tool-call shapes that this case might not match?).

Read-only in the two observed instances (no writes/commits touched the real repos), but the safety concern is general: an unscoped filesystem sweep gives a live FO access to whatever the host machine's filesystem actually contains, not just the sandboxed fixture. This matters most for local/operator-run live tests (this class of bug is host-specific, not reproducible on a clean CI runner where nothing real exists to find) but the underlying behavior — issuing an unscoped `find /` instead of a scoped search — is itself a production FO behavior, not a test-only artifact.

## Proposed direction (not yet fleshed — for ideation)

Candidates to weigh, not yet chosen:
- Extend the existing anti-hunt discipline (the zero-discovery no-hunt rule, and `dp`'s Move 1 reference-load no-hunt rule) to explicitly forbid an UNSCOPED sweep (no path argument, or a path argument that is a filesystem root / home directory) regardless of what is being searched for — narrower and more general than either prior fix, since this isn't about a specific known file, it's about sweep SCOPE.
- Check whether `detectBroadSearchAtBoot` already classifies an unscoped `find /` as a broad sweep (it likely does, per its own doc comment on `broadSweepTools`) — if so, the live-test gap here may just be that the classification correctly fires but the WANDER through a real repo happens before/alongside classification, which is a sequencing question, not a detection gap.
- Consider whether this warrants a runtime guard (not just a prose rule) given the safety framing — the contract's own principle prefers a code gate over prose-only where one is feasible.

## Acceptance criteria (draft — ideation to firm up)

- A value measure: on a local live re-run of the scenarios that surfaced this (filing, smallest-sufficient-mechanism), an unscoped sweep does not return real-repo results / does not lead the FO into exploring an unrelated real repo.
- Confirm whether existing detectors already correctly flag this shape as a wander (may already be covered — verify before assuming a gap).
- If a contract change is warranted, keep it scoped and evidence-based per the sibling entities' discipline (spike before committing, live-measure the actual claim).
