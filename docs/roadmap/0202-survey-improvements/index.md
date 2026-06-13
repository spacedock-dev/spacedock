# Sprint 0202 — survey improvements (honest lens + mode-aware framing)

**Goal:** make `spacedock:survey` *honest about what it sees* and *accurate per work-mode*,
without touching the decision-frontier triage that is the survey's emerging strength (and the
0.21.x wedge). Every fix is proven by the survey's **rendered output on a constructed
fixture** — never a prose-grep of `SKILL.md` (the survey-feedback discipline).

> **Status: scope-locked + COALESCED + BROADENED 2026-06-13 (captain); staff-reviewed (NOT-READY → `5wv` M1 folded → clean).** The six survey-cluster seeds all rewrite ONE skill into ONE locked output structure — so they execute as a **single task, `survey-output-redesign` (`5wv`)**, not six worktrees colliding on `SKILL.md`. The six seeds below are the **requirement record** (each a band inside `5wv`), deferred as standalone members (`sprint-readiness: defer`).
>
> **Broadened to a survey + cleanups release** — drivable members:
> - `5wv` survey-output-redesign (group: survey) — the headline.
> - `nd` prefer-`new`-over-`--next-id`, `td` mdschema-conformance-validator, `5ar` pre-cut-audit-cleanups-0199 (node-action deadline ~2026-06-16), `gf` state-sync local-mode degrade (group: cleanup).
>
> Drivable set: `--where sprint=0202-survey-improvements --where 'sprint-readiness != defer'` → the five above. Preflight staff review: `staff-review.md` (cleanups passed clean; `5wv`'s no-follow-up figure was the one Material, now folded). The bands below map to `5wv`'s R1–R6.

## Chosen output structure (captain, 2026-06-13)

**Value & numbers first.** The report opens with plain "what this gives you" + concrete
figures from the user's own data; mode/track vocabulary is demoted to a labeled detail section
below. This is the shared spine all six seeds render into — de-jargoning, plain terms
("manual" not "mechanical"), the lens caveats (`5x`/`za`), the empty-breakdown/scratch fixes
(`h5`), and the mode-aware offers (`zw`/`zb`) all land *inside* this structure. Canonical
reference mock:

```
SpaceDock survey — your last 30 days
(recent-window snapshot · agent logs only · 12 sessions had no working dir, not placed)

WHAT THIS GIVES YOU
  You steer your agents by hand ~100 times a month. About 60 are the same few moves.
  A SpaceDock workflow can run those for you and stop only where you'd want a say.

BY THE NUMBERS
  102  hand-steering interruptions
    6  hanging threads (started, never closed)
    3  decisions you made with no follow-up action
  218  sessions read (Claude 180 · Codex 38, by working dir; name-match would say
       410 — sibling repos, ignored)

HOW YOU WORK
  Plan -> worktree branch -> implement. Mostly manual, repetitive tracks
  (not trivial — they take real work).

  ↓ full analysis: modes, work-by-area, what we can't see
```

Every figure in the lede must derive from the surveyed session rows (an independent source),
never be templated prose — the same external-proof rule the ACs carry. "Threads to pull"
(`9h`#4) is the actionable companion: the hanging-threads / open-decisions list is the
steady-state the lede's numbers point at.

## The through-line

Three real user/author runs (a code repo, a knowledge-work repo, an author self-review) plus
a second round over two more repos converged on one message: **the core inference is good
(keep it), but the survey over-claims its lens and mis-frames work it hasn't mode-classified.**
The fixes group into three bands:

### Band A — the mode classifier (the spine)
- **`zb` knowledge-work-archetype** — name a **third archetype** beside mechanical-drive and
  exploration-steering: the knowledge-work loop (intake → process → file → log → close), so a
  notes/ops repo gets a recognized mode + honest offer instead of falling through to "generic
  book-keeping." (`zb#1`, lead with workdir-attributed Codex count, appears **already
  shipped** — verify against the current skill first; `zb` may narrow to just the archetype.)
- **`zw` iteration-framing + branch-attribution** — make the offer framing **mode-aware**:
  exploration mode leads with iterate/steer (not gates-as-the-whole-model); mechanical mode
  keeps the gate-drive framing (do NOT strip it). Make work-by-area **branch-aware** only as a
  detect-branch-and-merge + caveat (h5 found it didn't bite two repos — lower blast radius
  than first filed).

### Band B — the honest lens
- **`5x` lens-honesty** — own the **recent-window snapshot** (surface evolution where the
  signal exists, else frame "recent window, not whole history") and carry **one** clear
  partial-lens statement ("reflects the agent-log corpus; not captured: …") instead of
  scattered technical asides.
- **`za` subagent-dispatch-fact** (`sprint-readiness: ready`) — report the **fact** of
  orchestration so a spacedock-orchestrated repo isn't read as idle (most work lands in
  excluded subagent sessions). Dispatch fact only, not subagent content. This is the most
  concrete lens-completeness fix — a natural **first** dispatch.
- **`h5#1` Codex workstream null-signal** — stop presenting an always-empty
  `workstreams: (unlabeled) — N` breakdown; collapse to one honest unclassified line (or give
  it real signal).

### Band C — output hygiene
- **`h5#2` scratch-reasoning leak** — the model's `I have everything I need. Let me…`
  pre-amble must not precede the rendered report (the cross-check *findings* stay, as report
  content).

### Band D — value-prop legibility (live-run feedback, 2026-06-12)
- **`9h` survey-value-prop-legibility** — a third live run validated accuracy hard but found
  the OUTPUT illegible to a newcomer: lead with plain "this helps you do X" + a concrete number
  from the user's own data ("~N interruptions reducible"), not abstract jargon ("gated
  automation for the mechanical tracks, thread bookkeeping for exploration"); rename
  **"mechanical" → "manual"** for repetitive-but-substantive tracks (reserve "mechanical" for
  trivial); frame the decision-frontier as **"threads to pull, not threads you've lost"**
  (steady-state + unresolved + proactive suggestions, not decision-history narration). Couples
  to `zw` (mode labels) and `5x` (snapshot framing).

## Cross-revisions (bake into ideation, not yet in the seed bodies)

`h5`'s second-round runs explicitly revise the siblings — ideation MUST verify against the
current skill before building:
- `zb#1` (workdir-count-first) — likely **already shipped**; verify, narrow `zb` if so.
- `zw#1` (iteration vs gate) — **narrow to exploration-mode-specific**; mechanical offers keep gates.
- `zw#2` (branch-aware work-by-area) — **conditional** (detect + caveat), did not bite 2 repos.
- `za` — not exercised by the 2nd round (neither repo was orchestrated); still valid for orchestrated repos.

**Pending input:** the author-followup reply (drafted, with the captain, awaiting the author's
answers) folds into the classifier seeds (`zb`/`zw`) — hold those for it, or ideate and revise.

## Definition of done

- Each fix verified by a **survey run over a constructed fixture** rendering the corrected
  output — the expected value derived from the fixture's session rows (independent source),
  never from `SKILL.md` prose.
- The keep-signals do not regress: the recent-workflow inference, the per-area identity
  profile, and the decision-frontier triage all still render.

## The 0.21.x boundary (preserve, do NOT generalize here)

The NEEDS-YOU / BACKLOG / RECENT-DECISIONS triage + open-fork-vs-repo cross-check is a working
**single-repo prototype** of the cross-workflow decision-frontier / ready-room. 0.20.2
**preserves and extends it in-survey only.** Generalizing it (and the "agent log is a partial
lens" principle, which applies equally to `orient` and the ready-room) is **0.21.x
decision-abstraction** — out of scope here.

The 2026-06-12 live run sharpened the 0.21.x thesis: *"historic decisions about things you did
NOT choose pollute the context window; I want the steady-state of what's there now, not how I
got there — short-term vs long-term memory."* The "threads to pull, not threads you've lost"
reframe is the survey-local expression (Band D); the underlying principle — surface present
state + unresolved threads + proactive suggestions, de-emphasize decision-history — is a
**0.21.x decision-abstraction** design input, recorded here for that cycle.

## Cross-cutting items (NOT this sprint — surfaced for the captain)

From the live run, two non-survey asks worth filing separately:
- **Landing-page CTA** — put "run the scan on your own sessions" at the top of the landing
  page (`spacedock-landing`), CTA = run-it and/or email; users are "ready-fire-aim."
- **Bundle the agent-conversation index dependency** so it auto-installs (add to the plugin
  store) rather than being a manual prerequisite.

## Out of scope

- Reproducing surveyed corpus content (private; anonymization discipline — never in this repo).
- Re-opening `69`'s Codex workdir-attribution mechanism (the data is correct; these are framing).
- Generalizing the decision-frontier beyond the survey (0.21.x).

## Provenance

Six real runs/reviews (2026-06-08 → 06-12): a code repo (`zw`), a knowledge-work repo (`zb`),
an author self-review (`5x`), a second round over two repos (`h5`, which revises the cluster),
the `47rx` F-spike that surfaced the subagent-invisibility (`za`), and a live partner-meeting
run on ~200 sessions / ~2200 logs (`9h` — value-prop legibility, terminology, threads-to-pull;
also corroborated `zrc` non-sandbox auto-mode in 0.20.1). Corpus content omitted throughout per
the survey-feedback anonymization discipline; the named live-run source stays in the
uncommitted meeting notes.
