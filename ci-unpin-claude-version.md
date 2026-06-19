---
id: rrrhd7e79w41w1p39r0268e8
title: Unpin CI Claude Code version — run live-e2e on the merged-team floor (retire the #395 pin)
status: backlog
source: 'Captain 2026-06-18 — 9243/#396 (using-claude-team merged-model support) merged + green on 2.1.181. The #395 keystone pinned live-e2e to 2.1.177 (last native-TeamCreate release) ONLY to keep the legacy team contract alive. With the merged contract shipped, the pin should be retired so CI runs the current (merged-team) Claude. Ships in 0.20.5 alongside m4''s merged lane.'
started:
completed:
verdict:
score: 0.6
worktree:
issue:
---

Retire the 2.1.177 CI pin so the live-e2e lane runs current/unpinned Claude Code (the merged-team floor), now that #396 shipped the merged-team contract.

## Problem

`.github/workflows/runtime-live-e2e.yml` pins `SPACEDOCK_PINNED_CLAUDE_VERSION: "2.1.177"` (with `DISABLE_AUTOUPDATER: "1"`) because 2.1.177 was the last release exposing native `TeamCreate`/`TeamDelete`, and the legacy team contract required them (#395 keystone). #396 (`using-claude-team`) re-architected around the merged-team model (named background `Agent` + `SendMessage`, no `TeamCreate`), validated green on 2.1.181. The pin now holds CI on a deprecated team API while real users run the merged floor — CI no longer tests what ships. `internal/release/claude_version_pin_guard_test.go` actively ENFORCES the 2.1.177 pin (it REDs on any other pin / on no-pin install), so an unpin must update or retire that guard in the same change.

## Proposed approach (seed — ideation to flesh out)

- Set the live-e2e Claude install to the current/unpinned (merged-floor) version — drop the `${CLAUDE_VERSION:-$SPACEDOCK_PINNED_CLAUDE_VERSION}` 2.1.177 fallback, or repin to a merged-floor minimum (≥2.1.178). Decide: float-to-latest vs pin-a-merged-floor-minimum (latest tests reality; a floor avoids a surprise upstream regression mid-release — ideation picks one, with rationale).
- Retire / rewrite `claude_version_pin_guard_test.go`: it currently proves "the pin is exactly 2.1.177 with a team-tool rationale comment." Under the new policy it must either be deleted (no pin) or assert the new floor/rationale. The guard must reflect the shipped policy, not the old one.
- Keep `DISABLE_AUTOUPDATER` semantics coherent for a floating/floor version (don't let a job float mid-run if that breaks reproducibility — ideation resolves).
- The legacy interactive pty tests (m4) must SKIP on the merged host (no `TeamCreate`), not RED — coordinate with m4's skip-gating so the unpin doesn't turn the legacy lane red.

## Out of scope

- Deleting the legacy `using-legacy-claude-team` path — that retires later when STABLE Claude catches up to the merged floor (a separate trigger), not here.
- Authoring m4's merged lane (separate m4 work; this task only changes the version + guard).

## Acceptance criteria (provisional — ideation finalizes; proof = behavior, never prose-grep)

**AC-1 — live-e2e installs and runs a merged-floor Claude, not 2.1.177.**
Verified by: a green live-e2e CI run whose "Show tool versions" step reports a ≥2.1.178 (merged-floor) version, and whose team-mode scenario drives the merged path (no `TeamCreate`).

**AC-2 — the pin-guard reflects the new policy.**
Verified by: `go test ./internal/release/...` green under the new version policy (the guard deleted or asserting the new floor/rationale, not the retired 2.1.177 pin) — and demonstrably RED if someone re-pins to a team-tool-broken-for-legacy assumption that contradicts the shipped policy.

**AC-3 — the legacy interactive lane skips (not REDs) on the unpinned host.**
Verified by: the m4 pty team-mode tests SKIP (TeamCreate absent ⇒ merged host) on the unpinned live-e2e run, leaving the lane green via the merged scenario.

## Test plan

Edit the workflow + guard test; trigger a live-e2e run on the unpinned version and confirm AC-1/AC-3 (merged scenario green, legacy skipped). `go test ./internal/release/...` for AC-2. Cost: small diff, one CI run; the risk is the version-policy choice (float vs floor), resolved in ideation.

## Related

- `m4` live-team-mode-terminal-harness — ships its merged lane in 0.20.5 alongside this unpin; the two are the 0.20.5 cut.
- #395 pin keystone (`ea03d094`) — what this retires.
- #396 `using-claude-team` merged-model support — the contract that makes the unpin safe.
- `bare-mode-coverage` — the `-p`-assumes-bare audit that pin-retirement opens (sequenced after).
