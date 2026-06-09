---
id: 5ar2193yw8sv0rcyrt23wxg9
title: Pre-cut audit cleanups (0.19.9) — checksum-gate test, darwin-only doc drift, gofmt, hasGitEntry comment, node-action bump
status: backlog
source: "0.19.9 pre-cut antipattern audit (Commander staff review, 2026-06-08) — four recorded non-blockers, none gated the 0.19.9 tag. Seeded per the roadmap Close step."
started:
completed:
verdict:
score:
worktree:
issue:
sprint:
group:
sprint-readiness:
---

Four non-blocking findings the 0.19.9 pre-cut antipattern audit recorded (none blocked the cut; grouped here as the next-sprint seed). Small, independent.

## Items

1. **Checksum-gate fail-closed has no `*_test.go`.** v3's `curl | sh` installer checksum check is proven only in `install-e2e.yml` — no Go test loads it, so a contributor could delete the checksum verification and `go test` stays green. Add a `workflow_exec_guard`-style test (or a Go e2e tamper test) that goes RED if the checksum gate is removed/weakened. *This is the substantive one — a real test-strength hole in shipped install machinery.*

2. **Darwin-only prose is stale** in `docs/releasing.md` + the `release.yml` header — 0.19.9 added Linux binaries, but both still describe a darwin-only release. Doc fix. *(Backstop: `pj`'s AC-4 doc-reconciliation re-verifies `docs/releasing.md` against the real machinery at the flip; this fixes the platform-language drift sooner so the doc isn't inaccurate in the interim.)*

3. **gofmt drift** on `skills/integration/survey_sync_codex_test.go` (pre-existing, cosmetic). `gofmt -w` to clear.

4. **`hasGitEntry`'s only unmasked guard** is the single `TestDiscoverWorkflowsSkipsNestedCheckout` — add a cross-reference comment so a future edit knows that test is what protects it.

5. **Node-20 GitHub Actions deprecation — time-sensitive (≈2026-06-16).** The `v0.19.9` release run warned that `actions/checkout@v4`, `actions/setup-go@v5`, and `goreleaser/goreleaser-action@v6` are forced to Node-24 starting **2026-06-16** ("may not work as expected"). Bump these actions (and any sibling Node-20 actions across the workflows) to Node-24-compatible versions **before** the 0.20.0 flip cut runs `release.yml` — or before 2026-06-16, whichever comes first. Priority alongside #1, because the flip cut depends on a healthy `release.yml`.

## Out of scope

Anything requiring a behavior redesign. Each item is a localized test/doc/format fix.

## Acceptance criteria

(Ideation fills in. Each verified outside this body: #1 by a guard/tamper test that reds on a removed checksum check; #2 by the reconciled docs matching the real darwin+linux release machinery; #3 by `gofmt -l` clean; #4 by the comment present at the guard site — though #4 is a comment, so it rides #1's behavioral proof, not a standalone AC.)

## Notes

Provenance: 0.19.9 pre-cut audit. Candidate for a small cleanups sprint (cf. `019x-pre-flip-cleanups`) or to ride a near-term point release. `#1` (checksum test-strength) is the priority.
