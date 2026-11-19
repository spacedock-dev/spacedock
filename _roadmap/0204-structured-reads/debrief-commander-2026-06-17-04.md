---
session-date: 2026-06-17
sequence: 4
first-commit: a8b12ab7
last-commit: 44123714
duration: ~10h (2026-06-16 ~21:00 → 2026-06-17 ~07:00 PDT, with idle gaps)
---

# Session Debrief — 2026-06-17 #4 (0204 Commander)

Cut and **published v0.20.4** — the 0204 read-cost sprint's deliverable. The session's real work was disciplined verification: the handoff's central diagnosis was overturned, m4 was deferred to v0.20.5 on a byte-verified upstream regression, and a pre-cut antipattern audit caught + fixed two prose-grep guards (#394) before the cut. The `fo-self-evidence-bar` principle — verify each load-bearing claim from this run's evidence — was applied to the handoff, three successive ensign findings, and the audit verdict alike; it killed false blockers and confirmed real ones.

## Shipped
- **xht** `contractlint-antidrift-guard-hardening` — [#394](https://github.com/spacedock-dev/spacedock/pull/394). Retire two prose-grep doctrine guards (`TestUniversalDoctrineHasSingleSource`, `TestStartupGateGuidanceHasSingleSource`) — literal-phrase `strings.Contains` greps that missed paraphrase drift and violated the contractlint package's own `doc_test.go:11` policy — and close the AC-5 install-script comment-out hole via `executableShellCommands()`. The v0.20.4 pre-cut blocker fix; merged no-CI (admin-squash) per captain.

## Release
- **v0.20.4 cut + published** — 2026-06-17T15:41Z, stamp commit `44123714`, tag `v0.20.4`. Read-cost DoD MET (read helper #386, SOURCE opt-in #387, CI summary #389, host-neutral dispatch #391, commission-template restructure #388, journey-metrics read adoption #392 — across prior + this session). 10 release assets (darwin/linux × amd64/arm64, stable + edge) + checksums + journey-costs; homebrew casks bumped. The e2e-gate was satisfied via a **captain waiver** over the pre-existing zero-discover flake (set → used → **unset** so the next cut gates normally).

## Filed (backlog)
- **tq0** `zero-discover-broad-search-hardening` — the lean-boot zero-discover guard (`detectBroadSearchAtBoot`, #374) flakily reds because the FO still sometimes broad-searches on a zero-discover boot; the prose lever is spent (#374), so this is a residual stochastic discipline gap to harden or quarantine.
- **b5q** `codex-foreground-wait-phrase-check-retire` — `codex_foreground_wait_shape_test.go` (#378) carries the same prose-grep antipattern family (enumerated hyphenation variants standing in for a lifecycle meaning); pre-existing and out of the v0.20.4 surface, so it did not block the cut.
- **xht** `contractlint-antidrift-guard-hardening` — shipped same session (#394, above).

## Non-PR commits (workflow-only)
- `44123714` `release: bump version to spacedock@0.20.4` — the manual release-prep stamp commit (per `docs/releasing.md`); the tag points here.

State-branch transitions (m4 route→defer, the dispatch/advance/validation/archive churn for #394, the filed seeds) are routine state-machine churn, rolled into the sections above. The other main commits in range (`#393` haiku spike, `0205` index) are a peer session's, not this Commander's.

## Decisions (captain)
- **Cut handling: audit-then-cut** — run the pre-cut antipattern audit; on a clean audit, cut on the read-cost DoD.
- **m4: finish-into-v0.20.4 → then DEFER to v0.20.5.** Initially chose to finish m4 into the cut; when the byte-verified upstream block surfaced (2.1.179 dropped interactive team tools), chose to **defer** m4 rather than pin an old claude to green a test of a capability current users can't use. Cut on the read-cost DoD alone; m4's proven work preserved on #390.
- **Audit blockers: fix-first** — the two prose-grep guards fixed before the cut (the DoD requires a clean audit).
- **#394 merge: "push it and merge it, no ci required"** — admin-squash-merged the test-only fix without CI.
- **e2e-gate: waiver** — set `SPACEDOCK_E2E_GATE_WAIVER` over the pre-existing zero-discover flake rather than chase a flaky green; unset after publish.
- **Release note: concise / user-value** (over the `spacedock-release notes` comprehensive draft).
- **Production-impact finding (team mode down on 2.1.178+): note-only** — recorded here + in m4, no separate GH issue.

## Issues — Workflow
- `TestLiveZeroDiscoverReportsAndStops` flakily reds — a **pre-existing, stochastic** FO-boot-discipline gap (the FO sometimes broad-searches on a zero-discover boot). Byte-verified NOT a CC-version cliff: reds on 2.1.177 AND 2.1.179, and flips between sonnet/opus matrix legs of the same run. The handoff's "2.1.179 shift, pin-or-strengthen" framing was **refuted**. Filed `zero-discover-broad-search-hardening`; treated as re-run-grounds, never a merge blocker.
- Two prose-grep contractlint guards on the #388/#391 surface — fixed (#394).
- `codex_foreground_wait_shape_test.go` (#378) — same antipattern family — filed (`codex-foreground-wait-phrase-check-retire`).

## Issues — Spacedock / upstream (note-only per captain; not filed as GH issues)
- **Claude Code 2.1.178+ dropped the native TeamCreate/TeamDelete tools from INTERACTIVE sessions** — broader than #68721's documented *headless* scope, byte-verified on a 2.1.179 tmux FO (the deferred `ToolSearch(select:TeamCreate)` returns no-match; `select:SendMessage` resolves, ruling out a search confound; the "Claude Team" banner was the subscription/org label, not engaged team mode). Team mode (concurrent dispatch, comm-officer, standing teammates) is **down for any spacedock user on 2.1.178+** — silent bare-mode fallback. This session had team mode only because it launched on **2.1.177** before the auto-update. m4's live AC-3/AC-4 are blocked on upstream restore.
- **The harness's tmux-spawned FO suppressed its own on-disk transcript** because the child inherited `CLAUDE_CODE_*` markers (`CLAUDE_CODE_SESSION_ID`, `CLAUDE_CODE_CHILD_SESSION`, …) — the CC v2.1.170 "inherited env → no transcript" bug. Fixed in m4's harness with an env-scrub (`env -u` the nested-session markers, keeping the team flag) — proven live on 2.1.177; the FO transcript persists and the resolver→idle-gate→teardown-grade chain reads it end-to-end. (Earlier "FO writes no transcript" findings were this artifact, twice mis-mechanism'd before the env-scrub nailed it.)

## Observations
- **The self-evidence-bar earned its keep, repeatedly.** Three plausible "blocked" findings either dissolved or were corrected under one more layer of scrutiny before reaching the captain: (1) the handoff's zero-discover diagnosis — refuted; (2) the ensign's "FO writes no transcript" — twice confounded (wrong-place, then inherited-env), each corrected; (3) the ensign's "2.1.179 has no team tools" — the *first* verdict was confounded (deferred-tool eyeball), but the proper ToolSearch-hop re-probe confirmed it **real**. The discipline both killed false blockers and confirmed real ones with byte-level proof — and the same bar applied to my own escalations (I held the "AC-4 blocked" relay until I'd ground-truthed it).
- **The deferred-tool false-negative is a recurring trap.** A fresh FO/ensign eyeballing its init tool list sees no `TeamCreate` (it's deferred) and wrongly concludes "absent." The discriminator is always: run the `ToolSearch(select:…)` hop, then attempt the call.
- **Captain pushback is load-bearing.** "The upstream is `-p`, not interactive — am I missing anything?" forced the ground-truth probe that turned a relayed guess into byte-level proof (the regression *had* spread to interactive on 2.1.179 — beyond what the captain or the docs assumed).
- **The pre-cut aggregate audit caught a real defect the per-entity audits missed** (two prose-grep guards masquerading as proof). The honest resolution — *remove + report the owed test*, per the package's own policy, rather than fake a structural check that still greps prose — was validated by a fresh worker that *tried and failed* to construct a genuine check.
- **FO process slips to watch:** an early `cd` into a worktree drifted my shell cwd (the contract says stay at root, use `git -C`); recovered with no damage. Worth the `fo-self-evidence-bar` family of guards.

## What's Next
- **m4** (#390) → **v0.20.5**, pending upstream restoring interactive team tools (#68721 family) — or a capability-aware skip-not-fatal + a 2.1.177-pinned lane. Harness is proven and preserved on #390; the sprint re-tag (0204 → the v0.20.5 sprint) is a coordination decision for when that sprint forms.
- **0204 backlog (behind the initial gate):** **f5** `journeymetrics-ensign-read-adoption`, **82k** `read-guidance-redundant-with-grep`.
- **Filed this session (untagged backlog):** `zero-discover-broad-search-hardening`, `codex-foreground-wait-phrase-check-retire`.
- **`fo-self-evidence-bar`** (untagged) — the durable principle this session validated end-to-end; strong candidate to build next (its hard part remains the AC: a behavioral FO principle resists a clean code gate).
- **v0.20.4 is shipped; the 0204 read-cost DoD is met.**
