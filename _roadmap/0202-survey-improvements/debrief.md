---
session-date: 2026-06-13
sequence: 1
first-commit: 4d40680c
last-commit: 5605932f
release: v0.20.2
duration: ~1d (Commander cold-boot drive)
---

# Session Debrief — Sprint 0202-survey-improvements (Commander cold-boot drive) — 2026-06-13 #1

Commander cold-boot drive of the survey-output redesign + dev/CI cleanups to a **cut v0.20.2**. All five gate-approved members landed; the captain took the conn (approve + merge + CI authority) partway through. Two member feedback cycles (td test-strength hole, 5wv captain folds) and one CI-caught portability fix all closed cleanly. The mandated pre-cut antipattern audit ran clean (no material findings); the cut shipped on a green live-e2e bound to the exact tagged commit.

## Shipped — sprint members
- **5ar** `pre-cut-audit-cleanups-0199` — [#359](https://github.com/spacedock-dev/spacedock/pull/359). Bumped node-20 GitHub Actions to their node24 majors (deadline 2026-06-16, met 3 days early) with a node24-minimum guard test; a portable checksum-gate tamper test; release.yml header reconciled to the real darwin+linux build; gofmt + hasGitEntry comment.
- **5wv** `survey-output-redesign` — [#360](https://github.com/spacedock-dev/spacedock/pull/360). Redesigned the `spacedock:survey` step-4 report to a value-&-numbers-first spine, folding eight user-feedback seeds (R1–R6 + the book-keeping→per-entity-auto-processing offer reframe + spacedock-incumbent scaffold detection + naming the specific kind of knowledge work).
- **nd** `prefer-new-over-next-id` — [#362](https://github.com/spacedock-dev/spacedock/pull/362). Made `spacedock new` the blessed atomic-create path for filing entities; `status --next-id` emits a use-`new` stderr hint (stdout id unchanged); a `filing` shared runtime scenario (Claude + Codex) grades that the FO files via `new`.
- **td** `mdschema-conformance-validator` — [#363](https://github.com/spacedock-dev/spacedock/pull/363). Extended `status --validate` to schema-driven per-field conformance as a warn tier (prints to stderr, never flips the exit code or gates the read path; `Warning:` prefix pinned).
- **gf** `state-sync-no-origin-local-mode` — [#364](https://github.com/spacedock-dev/spacedock/pull/364). Split-root state sync degrades to local-only when the state checkout has no `origin`: boot exposes `state_remote`, dispatch drops the impossible push/pull while keeping the path-scoped commit.

## The cut
`v0.20.2` tagged on the gated commit `79284ca8` directly (no pre-stamp — matching v0.20.1 practice; releasing.md step 3 is stale, see Issues). The `e2e-gate` confirmed a green Runtime Live E2E run bound to that exact SHA (the drive pre-warmed it: all 5 lanes green, 3 live environments captain-delegated-approved). Tag push fired release.yml clean (e2e-gate · goreleaser · journey-ledger all green): GitHub Release v0.20.2 + 8 tarballs (stable+edge × darwin+linux × amd64+arm64), both Homebrew casks → 0.20.2, plugin manifests stamped 0.20.2 on main (post-tag child `5605932f`), `stable` channel ref advanced to it. The bridge (`ref: main`) and `stable` both serve 0.20.2, so 0.20.0/0.20.1 installs update correctly.

## Decisions
- **Captain delegated the conn** mid-drive: approve gates, merge PRs, and approve CI (including the env-gated live-e2e environments) at the Commander's judgement; only the irreversible `v0.20.2` tag push held for explicit authorize.
- **Fold both extra survey feedbacks into 5wv pre-merge** (rather than file as separate followups): the book-keeping→per-entity-auto-processing offer reframe, the spacedock-incumbent scaffold detection, and naming the specific *type* of knowledge work — three feedback cycles on one isolated skill, since 5wv was not deadline-bound and shares no files with the other members. The pre-filed scaffold-incumbent seed was archived as superseded.
- **Merge gate = the #348 offline suite**, not the flaky env-gated live lanes: main has no branch protection; the live lanes are non-required, slow, and (for 5wv/td/etc.) don't exercise the changed surface. Offline + build + install green was the merge bar; the live-e2e was run deliberately once, on the tag commit, for the e2e-gate.
- **Refused the boot reconcile drift remedy.** The reconcile sweep flagged Class-D/E drift against `origin/next`; its Class-E remedy ("reset main→origin/next") would have reverted the entire post-flip trunk. Acted on none of it — stale pre-flip tooling, not real drift (see Issues — Spacedock).

## Issues — Workflow
None. The five-member pipeline (ideation→implementation→validation→done) ran clean; the two feedback cycles and one CI fix were normal flow.

## Issues — Spacedock
- **`dispatch reconcile` conflates team-management with repo-hygiene, and hardcodes the pre-flip trunk `next`.** Classes A/B (lingering/superseded agents) are genuine team management; classes D/E (stale branch / stale local main) are pure git hygiene that hardcode `origin/next` (`reconcile.go:582/605/616`). reconcile landed `2026-06-02` (#273), 6 days *before* the `2026-06-08` flip — when `next` was genuinely the integration trunk and `main` tracked it. The flip inverted the model (`main` is now trunk, `next` dev-only) but the helper was never refit, so Class-E's "reset main→origin/next" is now exactly backwards. Root cause is the conflation: a team helper shouldn't carry repo/trunk knowledge. Filed (see Filed).
- **The `pr-merge` mod hardcodes base `next` (v0.12.1, pre-flip).** The merge hook opens PRs against `next`; this drive overrode the base to `main` per the dispatch doc. Same stale-trunk class as reconcile — the mod needs a refit so the base branch is `main`/config-driven, not a literal. Filed (see Filed).
- **`docs/releasing.md` step 3 (manual pre-stamp) is stale** vs the actual tag-the-gated-commit practice. Pre-stamping creates a new ungated commit the exact-SHA e2e-gate would block on; v0.20.1 and v0.20.2 both tagged the gated commit directly and let release.yml stamp post-tag. The doc should be reconciled. Filed (see Filed).
- **#358 parallel live-e2e underdelivered on the opus long pole** (a 0201 deliverable, first live exercise this run): all 5 lanes green, but opus was 14.2m vs the projected ~9m (old serial ~27m — a real ~47% cut, just short of estimate). Noted for the live-CI owners; not a 0202 defect.
- **Pre-cut audit non-blockers** (pre-existing, not 0202-introduced): gofmt drift on `internal/release/channel_agreement_guard_test.go` (#352) and `internal/contract/contract_test.go`; duplicated `stateHasOrigin` across `internal/status` and `internal/dispatch`. Filed as a combined cleanup (see Filed).

## Observations
- **CI earned its keep twice.** The #348 main-PR offline gate caught 5ar's load-bearing checksum sub-test relying on macOS `tar` accepting an appended-byte tamper (Linux `tar` rejects it independent of the gate → portable fix: valid tarball + mismatched hash). The detached adversarial audit caught td's `Warning:`-prefix being untested (an `Error:`-misreport would have shipped green). Both are exactly the test-strength holes per-entity validation in isolation cannot see.
- **The shared-root worktree hazard is live and biting.** The cold-boot recipe's `git switch -c drive/0202` switches the *shared* root working tree; concurrent Commanders/Shaping-FOs collide on it. During close-out the shared local `main` was found diverged (7 unpushed 0203 commits from a parallel Shaping FO, 13 behind origin) — this debrief was written from an isolated worktree off `origin/main` to avoid clobbering that unpushed work. The drive itself never depended on the root branch (all work lived in per-entity worktrees + the split-root state checkout). The boot recipe should drop the root-branch switch.
- **The decouple-first/cutover-last lesson from 0201 held this sprint**: the bridge manifest kept 0.20.0-era installs resolving through the cut. The 0.20.0→standalone-marketplace migration (before the bridge is removed) and CI coverage of the *released* binary's install path remain the highest-leverage open release-model gaps (0201-surfaced; not re-filed here).

## Filed (followup seeds — next sprint)
- `dispatch-reconcile-deconflate-repo-hygiene` — split repo-hygiene (D/E) out of the team-reconcile helper, or at minimum source the trunk branch from config instead of the hardcoded pre-flip `next`.
- `pr-merge-mod-base-branch-post-flip` — refit the `pr-merge` mod so its PR base is `main`/config-driven, not the hardcoded pre-flip `next`.
- `releasing-doc-pre-stamp-drift` — reconcile `docs/releasing.md` step 3 with the tag-the-gated-commit practice.
- `code-cleanups-0202-precut` — gofmt drift on two pre-existing files + dedup `stateHasOrigin`.

## What's Next
0202 is closed — all five members archived PASSED, v0.20.2 cut and verified. The followup seeds above are filed to the dev backlog for the next Shaping FO to carve. No 0202 entities remain dispatchable or gate-blocked.
