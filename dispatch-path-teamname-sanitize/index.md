---
id: e30vayp9ja885x7qfjkm6w2n
title: "Validate/sanitize team_name in the dispatch-file path + fix the bare-mode comment (n1a audit polish)"
status: validation
source: "n1a detached audit (2026-06-04, audit-live-cycle-end-state-determinism @ 9c6b7ae7) — refuted nothing material; two polish items on internal/dispatch/build.go's new team_name-keyed dispatch path (#n1a 1b)."
score: "0.18"
worktree: .worktrees/spacedock-ensign-dispatch-path-teamname-sanitize
started: 2026-06-08T05:11:44Z
completed:
verdict:
issue:
sprint: 019x-pre-flip-cleanups
group: dispatch-hygiene
sprint-readiness: ready
---

n1a's 1b keys the dispatch-file path on `team_name` (`/tmp/spacedock-dispatch/{teamName}-{derivedName}.md`). The audit found two non-blocking polish items in `internal/dispatch/build.go`.

## Problem

1. **Unvalidated input into a filesystem path.** `teamName` is prepended to the dispatch filename (build.go ~570-573) WITHOUT validation, while its sibling `derivedName` IS validated (`namePattern` kebab-only, `nameMaxLen=200`, ~427-438). A `team_name` containing a path separator (`../`) or odd chars would build an unsanitized `/tmp` path. NOT a live exploit today — team names are harness-generated `{project}-{dir}-{YYYYMMDD-HHMM}-{shortuuid}` (kebab-safe) and are FO/harness-controlled, never untrusted external input — so this is defense-in-depth, not a security boundary. Also no length cap on the combined `teamName-derivedName` (~240 chars worst case; within fs limits but a thin margin).
2. **Misleading comment.** build.go ~567 says bare-mode "(no team name)" keeps the plain name, but the code keys solely on `teamName != ""`, independent of `bare_mode` (a `bare_mode:true` + `team_name:"x"` dispatch still gets the `x-` prefix). Harmless (collision-free either way) — the comment conflates "team_name present" with "not bare."

## Proposed approach (ideation firms)

Apply the existing `derivedName` validation (or a path-safe sanitize) to `teamName` before using it in the path, with a combined-length cap; correct the comment to say "team_name present" not "not bare." Smallest reasonable change — mirror the proven validation already next to it.

## Acceptance criteria (seed)

- **AC-1 (seed):** `teamName` is validated/sanitized to the same path-safety as `derivedName` (no separators/odd chars; combined length capped) before path construction — verified by a unit test feeding a path-unsafe team_name and asserting a safe path (or a clean error), and the combined-length cap.
- **AC-2 (seed):** The build.go bare-mode comment matches the actual condition (`team_name != ""`, not `bare_mode`).

## Notes

0.19.6 follow-up; does NOT block 0.19.5 (no live exploit; n1a shipped on the audit's clean-material verdict). Source: n1a (`live-cycle-end-state-determinism`) audit. Low score — pure hardening + comment.

## Stage Report: implementation

- DONE: AC-1: Validate/sanitize teamName to the same path-safety as derivedName (reuse the existing namePattern + nameMaxLen in internal/dispatch/build.go) with a combined-length cap, BEFORE the dispatch-file path is constructed — proven by a unit test that feeds a path-unsafe team_name (e.g. containing `../`) and asserts a safe path or a clean error, plus the combined-length cap.
  internal/dispatch/build.go: teamName checked against `nameMaxLen` + new `teamNamePattern` (same char class as namePattern), combined filename capped at new `dispatchFileNameMaxLen=251`, all before path construction. Tests in build_teamname_path_test.go (9 subtests: 6 path-unsafe inputs incl. `../escape`/`a/b`/`/abs`, combined-length cap, valid-path regression) — all green. Also verified end-to-end via the built CLI: `../escape` and `a/b` error cleanly (exit 1, no traversed path), real-shaped name builds the correct keyed path.
- DONE: AC-2: Fix the build.go bare-mode comment (~line 567) so it matches the actual condition (the code keys on team_name != "", not on bare_mode).
  Comment now reads "Dispatches with no team_name keep the plain derived name." — keyed on team_name presence, not bare_mode.
- DONE: Regression: go test ./... stays green.
  `go test ./...` → 1150 passed in 16 packages; gofmt clean; go vet no issues.

### Summary

Validated team_name for path-safety before it is prepended to the /tmp dispatch filename, mirroring derivedName's existing checks, and added a combined-filename length cap (`dispatchFileNameMaxLen=251`, below the 255-byte fs name limit). Design decision (escalated to FO, proceeded on the path-safety-intent reading after no reply): introduced a dedicated `teamNamePattern` that reuses namePattern's kebab character class but admits a single safe char, because namePattern's two-anchor shape rejects single-char team names (`"t"`) that 19 existing fixtures use as placeholders — this keeps the whole suite green and avoids churning oracle-parity goldens, while rejecting every genuinely path-unsafe input. Corrected the misleading bare-mode comment to key on team_name presence.
