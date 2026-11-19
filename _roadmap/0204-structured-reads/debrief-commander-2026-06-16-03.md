---
session-date: 2026-06-16
sequence: 3
first-commit: b5867dc1
last-commit: 12b67ad5
duration: ~6h
---

# Session Debrief — 2026-06-16 #3

0204 Commander session. Shipped two read-cost members (ezf host-neutral core; hf measure-then-trim, absorbing 4x); m4's live team-mode harness fought a cascade of Claude-Code interactive-team seams and is still open (#390). The deeper yield was FO-discipline: a Proof-policy gate (required CI lanes follow the diff) plus three filed follow-ups exposing where the FO over-trusts itself.

## Shipped
- **ezf** `runtime-neutral-dispatch-core-cleanup` — [#391](https://github.com/spacedock-dev/spacedock/pull/391). Made the host-neutral dispatch core genuinely runtime-neutral: five-leak prose move (Claude/team-only imperatives → the Claude adapter) + a structural contractlint check (literal-token scan, phrase-level adapter-as-subject exemption — not prose-grep).
- **hf** `journeymetrics-read-adoption-metric` — [#392](https://github.com/spacedock-dev/spacedock/pull/392). Surface `status --read` + scoped-`Read` adoption on the journey record, and trim the redundant site-6 dispatch-prompt hint (13 goldens regenerated). Absorbed 4x as the trim half (measure-then-trim).

## Filed (backlog)
- **f5** `journeymetrics-ensign-read-adoption` — hf's metric parses the FO front-door stream, not the dispatched ensign it was built to watch (4 real FO captures read 0/0); fold the ensign sub-agent transcript into the metric.
- **82k** `read-guidance-redundant-with-grep` — the `status --read` section-read guidance is redundant with the grep-anchor rule (entity-body heading maps measured IDENTICAL to grep; grep over-counts only fenced content like the README template; residue is fence-safety/frontmatter).
- **`fo-self-evidence-bar`** (id z25…, no sprint tag) — bind the FO's OWN gate/merge/triage decisions to the evidence bar in the runtime-neutral Working Principles; spans both autonomous AND the present-gate evidence the captain votes on.

## Archived (superseded)
- **4x** `read-hint-adoption-bloat-trim` — `verdict: superseded`; folded into hf as the trim half.

## In flight (NOT shipped)
- **m4** `live-team-mode-terminal-harness` — PR #390 OPEN; status=validation, mod-block=merge:pr-merge, pr=#390. See What's Next.

## Non-PR commits (workflow-only, main)
- `b6e38c25` proof policy: required CI lanes follow the diff, not the FO's read of "relatedness" (+ failures read, not inherited). FO-direct README upkeep.

## Decisions (captain)
- Gave the conn (approve gates / merge / approve CI). Approved m4 ideation→implementation.
- **Fold 4x into hf** (measure-then-trim); dispatched hf in 4x's place (concurrency:3 was full with m4+ezf+hf); 4x archived superseded.
- Push ezf + hf; **merge on the deterministic lanes**. Then corrected (below): claude-live was NOT unrelated to ezf.
- For m4: refresh `~/.claude/benchmark-token` from the Keychain; **green locally first**; then **use CI** (ANTHROPIC_API_KEY — no OAuth quota) for the live proof; **study spacedock-gym** and fix to pass locally.
- Approve hf + file the AC-7 follow-up (→ f5). File the read-guidance redundancy (→ 82k). Strengthen FO proof discipline: #1 README gate (done), #2 the contract Working Principle (→ fo-self-evidence-bar).
- **Corrected** the ensign's "FO writes no transcript" finding — FALSE (this session has a session jsonl); it is a resolver bug, not an architecture wall.

## Issues — Workflow
- m4 harness seams (all real, all caught by green-local-first): macOS `/var`→`/private/var` symlink false-positive in the wrong-root detector (fixed `cf6f2eba`); F30 transcript-pin — the tail flipped to the comm-officer's newer transcript (fixed by pinning the FO session id, `f1d9f120`); CI render-fragility — the pane-text idle gate returned 0/3 on the headless runner (minimal Linux lacks the `tmux-256color` terminfo → blank `capture-pane`), replaced with a gym-grounded on-disk turn-end gate `transcriptReachedIdle` (`bdd86f3a`); the interactive greet-stop (the FO parks for the captain by contract — handled by a bounded captain-nudge, **uncommitted**); and the transcript-resolver not finding the FO transcript (ensign wrongly concluded the FO writes none — **CL corrected: it DOES**).
- `TestLiveZeroDiscoverReportsAndStops` red on main's live lane: the FO broad-searches the filesystem (`find … README.md | head -20 && ls …`) on a zero-discover boot instead of report-and-stop — most likely a CC-2.1.179 behavior shift (CI is unpinned; prior "passes" rested on pinned 2.1.178).
- The pty live step was chained behind the flaky ensign-cycle step (fail-fast skipped it); decoupled with `if: ${{ !cancelled() }}` so each live step's result surfaces independently.

## Issues — Spacedock (framework)
- **The FO contract aims its proof discipline at the ensign's deliverable + the gate review, never the FO's OWN decisions.** Twice this session the FO was too loose with the live lane: (1) merged a Claude-adapter change (ezf, `claude-fo-dispatch.md`) on deterministic lanes while leaving `claude-live` unapproved — that lane is the live drive for that very file; (2) labeled a live-CI red "the known flake" without reading it (it was `TestLiveZeroDiscover`, not `TestLiveEnsignCycle`). → **fo-self-evidence-bar** filed; README gate landed. (Not filed as GH issues — captured as workflow tasks per FO write scope.)
- `status --read` adoption guidance redundant with grep → **82k**. journeymetrics `--read` metric watches the FO not the ensign → **f5**.
- The interactive team-lead FO's session jsonl was hard to locate (resolver found only an empty `session-env/<uuid>/`); the FO DOES write one — a documented "where the interactive team-lead transcript lands" would have prevented both F30 and the false no-transcript conclusion.

## Observations
- **Green-local-first earned its keep.** It caught the symlink + F30 bugs the FO's original "AC-3/4 pending CI" stance would have shipped, and surfaced the CI render-fragility and the FO-broad-search regression. The captain repeatedly corrected the FO's reach-for-a-label-instead-of-reading habit — which is exactly what `fo-self-evidence-bar` encodes.
- The pty harness's correct observation surface is on-disk (the FO transcript + the team `config.json`), not the pane. Every m4 seam was a symptom of observing the wrong surface or the wrong file.
- m4 is team-mode coverage, NOT part of the read-cost DoD — it does not block the v0.20.4 cut.

## What's Next
**v0.20.4 cut — read-cost DoD is MET (e6a/0q/48/j7/6rt prior + ezf/hf this session):**
1. Decide `TestLiveZeroDiscover`: pin CC-2.1.178 in CI to confirm it's the model shift (track the FO-broad-search separately), OR strengthen the "zero → report-and-STOP" guidance to hold under 2.1.179.
2. Run the pre-cut antipattern audit (DoD requirement).
3. Stamp `0.20.3 → 0.20.4` (both plugin manifests) + tag v0.20.4 (mirror `9bd1f46a`).

**m4 (#390 open) — roll to v0.20.5 vs finish now (captain decision):**
- FIRST correct the committed "FO-writes-no-transcript" finding (`12b67ad5`) — it is FALSE. The FO writes a session jsonl; the resolver isn't finding it (likely `session-env/<uuid>/` or a different projects-dir encoding — same class as F30). Fix the resolver to locate the FO transcript; THEN `transcriptReachedIdle` (idle gate) and AC-4 (marker grade) are viable.
- Commit the bounded captain-nudge (shape A) once the resolver is fixed. AC-3 (comm-officer roster via `config.json`) already works live; AC-4 (FO marker) viable once the transcript is located.
- Re-green locally → re-push #390 (the gym-grounded idle gate + the decoupled CI steps should let CI's pty step run).

**Backlog follow-ups (post-cut):** f5, 82k, fo-self-evidence-bar.

**Coordination:** checkout shared with peer sessions (a shaping FO, a 0205 Commander — their fo-tier-delegation / contract-prose-microtest / 0205 commits interleave on the state branch). Commit path-scoped; the team `spacedock-v1-dev-20260616-1227-bcf4a8f7` + the m4 impl ensign + comm-officer were this session's (tear down at exit).
