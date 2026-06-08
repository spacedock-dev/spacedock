---
id: 47rx3x8a809wx35vx6rbqqhv
title: Survey skill body-rendering pass — Codex workstreams (depth ii) + scaffold/work-area/confirm fixes
status: ideation
source: "captain (2026-06-08) — deferred from vh (survey-skill-correctness-pass, 0.19.8). vh's live drive proved the codex-presence COUNT fires (61), but Codex sessions are not surfaced in the survey body, and a sandbox-denied ~/.codex yields a silent confident 0. Captain: file these as a bundled 0.19.9 candidate."
score: "0.25"
started: 2026-06-08T19:37:02Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-survey-codex-and-sandbox-followups
issue:
group: survey
sprint-readiness: ready
sprint: 0199-pre-flip-mechanics
---

**47rx is the survey-skill body-rendering pass for 0.19.9** (`skills/survey/SKILL.md` + `references/queries.sql`). Four deliverables, all body-rendering, all bottoming out on a live drive (survey discipline — no SKILL.md grep as a behavioral AC):

- **D1 — Codex body-surfacing at DEPTH (ii)** (captain 2026-06-08): surface Codex WORKSTREAMS + activity in the body, folded via Codex's own agentsview-DB signals (`first_message` → workstream clusters; `update_plan`/`exec_command` → activity), workdir-attributed and sibling-free (the `$.workdir` discovery). DB-only.
- **C — drop the step-4 confirm gate**: auto-render the full body after the one-line headline instead of stopping to ask "Want me to lay it out?".
- **E — de-narrate SCAFFOLD discovery**: state usage + checked-in-or-not only; drop the "recovered from behavior, not files" / "recovered one-offs" mechanism narration.
- **F — Work By Area leads with PRODUCT paths**: rank product paths (`docs`/`internal`/`cmd`/`skills`/…) first; demote the agent-mechanics paths (`<external>`/`/tmp`-dispatch, `.worktrees`, `.claude*`) to a footnote.

**Deferred (banked, upstream-gated):** D2 (silent-0 caveat) was DROPPED; depth (iii) (per-file Codex work-by-area) and D3 (source-health) are banked upstream-gated follow-ups. **0.19.9 candidate per captain; roadmap/sprint assignment is the captain's.**

## Problem

vh shipped the corrected git-root-basename model and a flagged Codex presence count, but left a gap the live drive surfaced:

- **Codex sessions are absent from the survey BODY.** The report shows a flagged Codex presence COUNT only (matched by `project` name, caveated "may include a same-named sibling"); the decisions / SCAFFOLD / work-by-area / workstreams body is Claude-only. Root reason: agentsview does not record Codex cwd (every Codex session lands `cwd=''`), so a Codex session cannot be cwd-prefix-scoped — and the cwd-prefix is the only thing that excludes a **same-basename sibling repo** (`project` is git-root-basename, which collides: `spacedock`→3 roots, `workspace`→6). Folding Codex into the body naively would corrupt the counts with sibling sessions. This exclusion is correct today (vh AC-2 asserts codex-not-folded-into-scope), but it means a Codex-heavy repo's real work (e.g. this repo's codex-adapter sessions) never reaches the body. **The spike below found a clean attribution signal (`exec_command.$.workdir`) that DOES separate this repo's Codex from a same-basename sibling — so the original "Codex can't be scoped" blocker is resolved. 47rx closes this gap at depth (ii): Codex workstreams + activity in the body (see Scope decision).**
- **(D3) No source-health signal.** The survey cannot detect that a Codex source was denied/disabled vs. genuinely empty. Confirmed upstream-gated below; deferred.

## Scope decision: D1 = depth (ii) — Codex workstreams + activity via DB signals (LOCKED); depth (iii) + D3 deferred upstream-gated

**D2 is dropped (captain, 2026-06-08).** **D1 is firmed to DEPTH (ii) (captain, 2026-06-08).** 47rx ships D1(ii); depth (iii) and D3 are banked as upstream-gated follow-ups.

**D1(ii) — surface Codex WORKSTREAMS + activity in the survey body, folded via Codex's OWN agentsview-DB signals.** The `$.workdir` discovery (Spike Part 1) cleanly scopes Codex to this repo, sibling-free, so attribution is settled. On that attributed set, depth (ii) adds two body signals, both read straight from the synced DB — NO raw-rollout parsing:
  - **Workstreams** — cluster the attributed Codex sessions into ensign-task workstreams from each session's `first_message` (e.g. `journey-cost-ledger`, `orient-workflow-discovery`, `codex-foreground-wait-escape-hint` — see the rendered view + the explicit clustering rule below). These are real tracks the Claude-only body misses entirely.
  - **Activity** — a per-session activity summary from the DB tool-call tallies (`exec_command` command count, `update_plan` plan-step signal, `spawn_agent` multi-agent fan-out). The rendered codex-only view (Spike Part 2) shows this maps coherently — it is NOT garbage.

**Deferred (banked, upstream-gated) — the exploration's value is kept as recorded follow-ups, NOT discarded:**
- **Depth (iii) — per-file Codex WORK-BY-AREA** (which files Codex changed: `internal`/`skills`/`docs`). Out of 0.19.9. agentsview does not ingest structured Codex file edits — it stores `apply_patch` as empty `write_stdin` rows — so the only path is fragile parsing of `apply_patch` markers in the raw `~/.codex/sessions/*.jsonl` rollouts (a plain grep returns garbage; even a JSON parse is brittle). The clean fix is the upstream dependency below. The rendered depth-(iii) view (Spike Part 2, B-3) is banked as the proof-of-value for when the upstream gate clears.
- **D3 — source-health detector.** Distinguishing a denied/disabled Codex source from a genuinely-empty one. Confirmed against agentsview v0.32.1: neither `agentsview sync`, `agentsview projects` (lists project→count but not agent/source), nor any other subcommand exposes a per-source count or source-health signal; an agent-side `[ -r ~/.codex ]` probe faces the same Seatbelt FS-access denial that forced vh's AC-5 `command -v`→`agentsview --version` swap. Needs an upstream source-health signal.

**Upstream agentsview dependency (recorded so the deferred value stays actionable).** Both deferred items resolve cleanly with one upstream `kenn-io/agentsview` change set: (a) **persist Codex `cwd`** (already present in the raw rollout's `session_meta`, just not ingested — would also strengthen D1's attribution from a heuristic to an exact scope); (b) **ingest structured Codex file edits** (extract `apply_patch` `*** Update/Add File:` targets into a queryable shape — unblocks depth (iii) without rollout parsing); (c) **expose a per-source health/count signal** (unblocks D3). These are a dependency for the deferred follow-ups, not a deliverable of 47rx.

**Ship shape.** D1(ii) ships as one task (this entity), all from the agentsview DB — no binary/agentsview changes. Depth (iii) + D3 defer on the upstream dependency above. The AC + test plan below are firmed to depth (ii).

## Spike: D1 feasibility — exercised end-to-end against live agentsview v0.32.1 (the linchpin)

### Part 1 — attribution heuristic (can Codex be separated from a same-basename sibling?)

**The riskiest unknown** (per the assignment): can a heuristic safely separate THIS repo's Codex sessions from a same-basename-sibling repo's, given agentsview records no Codex cwd and `project`=git-root-basename collides? Spiked end-to-end, NOT re-read.

**Method.** Synced this repo's Claude sessions (narrowed farm) + the unscoped Codex backlog through the agentsview binary into a readable `AGENTSVIEW_DATA_DIR`, then queried the live `sessions.db` directly. This repo keys to `project='spacedock_v1'`; a same-basename SIBLING keys to `project='spacedock'` (a DIFFERENT repo at `/Users/clkao/git/spacedock` — the real collision vh documented, here with 160 Codex sessions vs this repo's 61).

**Result — the `workdir` signal is real and cleanly separates the siblings:**
- All 61 `spacedock_v1` Codex sessions have blank `cwd` AND blank `git_branch` (confirms the session-level cwd is genuinely absent) — but their `exec_command` tool calls carry `$.workdir` (e.g. `/Users/clkao/git/spacedock-research/spacedock-v1/.worktrees/spacedock-ensign-codex-safehouse-launcher`).
- **57/61** `spacedock_v1` sessions have at least one `exec_command` whose `$.workdir` is under the repo-root prefix → confidently THIS repo. **1** session has workdirs but only outside the prefix (a `/tmp` audit run). **3** sessions have no usable workdir (1–4-message near-empty sessions). So the heuristic's recall is **57/61 ≈ 93%**; the 4 misses are tmp-only or near-empty sessions, not sibling confusion.
- **Cross-contamination = 0 (the precision result that matters).** The 160 sibling-`spacedock` Codex sessions ALL carry workdirs under `/Users/clkao/git/spacedock/...`; ZERO of them fall under the `spacedock-v1` prefix. So the prefix-on-workdir scope admits no sibling session — same absolute-prefix discipline as Claude's `cwd`, same exclusion guarantee.
- **Field stability.** Across ALL 60,269 Codex `exec_command` calls in the live DB, 57,144 (95%) carry a non-empty `$.workdir`; for recent sessions (2026-04→) it is ~96% per session. `exec_command` is the ONLY Codex tool carrying a path/workdir (no `apply_patch`/`Edit`/`Write` — Codex edits via shell). So the heuristic does not depend on a rare field.

**Part-1 conclusion.** The attribution heuristic is viable and de-risked: scope Codex by `exec_command.$.workdir` under the repo-root prefix — high precision (0 sibling leak), ~93% recall, on a near-universal field. The 4-session recall gap is the honest cost (a tmp-only or near-empty Codex session is unattributable and excluded — conservative, like the Claude blank-cwd exclusion). The COUNT of this repo's Codex sessions does NOT need the upstream Codex-cwd change.

### Part 2 — RENDERED CONTENT VIEWS (the captain asked to SEE the content, not counts)

All rendered from the real synced agentsview v0.32.1 DB over the 57–58 `$.workdir`-attributed Codex sessions for this repo (the Part-1 set). These are the evidence the D1 depth call rests on. **Correcting my cycle-2 framing: the codex-only content is NOT garbage** — workstreams and tool-call mapping render coherently from the DB; only the file-level work-by-area needs raw-rollout parsing.

#### CODEX-ONLY VIEW (this repo's workdir-attributed Codex sessions, on their own)

**(B-1) WORKSTREAMS — DB-derivable from `first_message`, and coherent.** Codex sessions cluster cleanly into the same ensign-task workstreams the Claude body would show, each spanning ideation/impl/validation:
```
journey-cost-ledger                     8 sessions
codex-foreground-wait-escape-hint       5
orient-workflow-discovery               5
spacedock-claude-auto-install-plugin    4
codex-idle-notification-probe           3
dispatch-build-help-ergonomics          3
dispatch-build-json-ergonomics          3
codex-live-ci                           2
readme-main-flip-reconciliation         2
sweep-guard-reader-axis-invert          2
```
These are REAL workstreams the Claude-only body misses entirely — e.g. `journey-cost-ledger` and the whole `codex-runtime-adapter` line were largely Codex work.

**(B-2) TOOL-CALL MAPPING — does Codex map correctly? YES, agentsview maps Codex's tools coherently:**
```
exec_command  6834   shell cmd + $.workdir (the attribution signal; reads dominate)
write_stdin    195   stdin to a running exec — chars often EMPTY (this is where apply_patch lands)
update_plan    136   STRUCTURED plan steps (status + step text) — a real workstream/decision signal
wait_agent      62   await sub-agent
spawn_agent     54   sub-agent spawn (Codex multi-agent)
close_agent     46   close sub-agent
send_input      24
```
The mapping is correct and meaningful — NOT mis-mapped/garbage. `update_plan` even carries structured plan progression, e.g.:
> `[{"step":"Explore project context","status":"completed"},{"step":"Ask clarifying question","status":"completed"},{"step":"Propose approaches","status":"in_progress"},{"step":"Present UX/agent design","status":"pending"},{"step":"Write design doc","status":"pending"}]`

**(B-3) WORK-BY-AREA — coherent, but RAW-ROLLOUT-derived (not in the agentsview DB).** Codex edits via `apply_patch`, which agentsview stores as empty `write_stdin` rows (no `file_path`). Parsing the `*** Update/Add File:` markers from the raw rollout JSON recovers a clean, coherent work-by-area (456 in-repo file-edits across 58 sessions):
```
internal                   247   e.g. internal/cli/init_test.go, internal/dispatch/build_codex_host_test.go
docs/dev/.spacedock-state   85   e.g. .spacedock-state/codex-runtime-adapter/index.md
skills                      48   e.g. skills/first-officer/SKILL.md, skills/ensign/SKILL.md
docs                        19   e.g. docs/dev/README.md
cmd                          8   e.g. cmd/spacedock-release/main.go
.github                      6   e.g. .github/workflows/release.yml
```
This is the real Codex work signal — and it is coherent. The catch is the SOURCE: it comes from raw rollouts (`~/.codex/sessions/*.jsonl`), NOT the agentsview DB the survey queries. A plain `grep` over those rollouts returns garbage (catches `rtk grep` help text); a proper JSON parse is needed. So depth (iii) is buildable but adds a rollout-parsing dependency the survey skill doesn't have today.

**Codex-only verdict: COHERENT, not garbage.** Workstreams + tool-mapping (depths i/ii) render correctly straight from the DB. File-level work-by-area (depth iii) is coherent too, but only via raw-rollout parsing.

#### COMBINED VIEW (Claude + Codex folded together)

**(A-1) The two agents did COMPLEMENTARY work — the combined view is MORE informative than either alone:**
```
                 CLAUDE (structured Edit/Write)        CODEX (apply_patch, rollout-derived)
WORK BY AREA      <external> 357  (mostly /tmp/         internal 247
                   spacedock-dispatch handoff files)    docs/.spacedock-state 85
                  docs 277                              skills 48
                  .worktrees 6, .github 3, skills 1     docs 19, cmd 8, .github 6
```
Claude's *visible* work-by-area is dominated by `/tmp` dispatch/handoff scaffolding (`<external>` 357) and `docs` — the orchestration layer. Codex's is `internal` Go code + `skills` + state. (Caveat: Claude ALSO edits `internal/` heavily — 7536 edits — but under `.worktrees/<wt>/internal/...` cwds the work-by-area buckets elsewhere; the headline buckets are what the survey body actually renders.) The honest read: **the combined view tells a fuller story** — it shows Codex carrying a large share of the `internal`/`skills` implementation that the Claude-only body under-represents.

**(A-2) The NAIVE structural fold adds nothing — Codex must be folded via its OWN signals.** Simply broadening the existing work-by-area query to `agent IN ('claude','codex')` yields a view byte-identical to Claude-only (Codex has 0 structured Edit/Write rows). A real combined view must add Codex's signals from THEIR shape: workstreams from `first_message`, activity from `update_plan`/`exec_command`, file edits from rollout `apply_patch`. The combined view is only as deep as the Codex signals you choose to surface (i/ii = DB; iii = rollout).

**Combined verdict: coherent and additive at every depth.** The captain selected **depth (ii)** — count + workstreams + activity, all DB-only (no rollout parsing). Depth (iii) (file-level work-by-area) is banked as the upstream-gated follow-up.

## D3 spike: source-health signal — confirmed absent in agentsview v0.32.1 (upstream-gated)

Checked the binary's surface for a per-source/source-health signal that would let the survey tell "Codex source denied" from "genuinely zero Codex": `agentsview --help` (commands: health, projects, sync, stats, …), `agentsview projects` (emits `PROJECT  SESSIONS` — NO agent/source column, so it cannot say "0 Codex because the source was empty/denied"), `agentsview sync` (no `--full`-independent per-source summary line on stdout/stderr). None exposes per-source health. Combined with the Seatbelt denial on a direct `[ -r ~/.codex ]`, D3 has no in-0.19.9 path — it is upstream-gated on a `kenn-io/agentsview` source-health signal. D3 defers.

## Proposed approach (firmed to depth (ii))

**Step 1 — the attribution foundation (`codex-scoped` query).** Add a labeled `-- name: codex-scoped` query to `queries.sql` that scopes `agent='codex'` sessions by `exec_command.$.workdir` under the repo-root prefix (the workdir analogue of the Claude cwd-prefix scope — admits this repo's Codex, excludes a same-basename sibling). This is the de-risked Part-1 mechanism. It STAYS alongside the existing project-NAME-only `codex-presence` flag (the superset that may include siblings) — the two are distinct signals.

**Step 2 — Codex WORKSTREAMS (the depth-(ii) body content).** Over the `codex-scoped` set, cluster sessions into ensign-task workstreams from `first_message`. The clustering rule is explicit (it is NOT a single regex — three cases, in order):
  1. **Dispatch pattern** (~40/58 of this repo's attributed Codex) — `first_message` matches `spacedock-ensign-{TASK}-{stage}` (a dispatch-file read); the workstream label is `{TASK}` (strip the trailing `-ideation`/`-implementation`/`-validation` stage and anchor the end on the `.md`/backtick boundary so trailing text doesn't leak into the label).
  2. **Task/entity pattern** (~9/58) — `first_message` names `Spacedock task \`{TASK}\`` / `Spacedock entity {id} \`{TASK}\``; the label is the backtick-quoted `{TASK}`.
  3. **Unlabeled** (~11/58) — encouragement/meta messages ("You totally got this…", "Captain asked me to tell subagents…") or a null `first_message` carry no task; they fall to an `(unlabeled)` bucket counted but not named. (Never invent a label — an honest `(unlabeled) N` bucket, per the survey's fill-every-slot/never-invent rule.)
  The body renders the resulting `workstream → session-count` clusters (e.g. `journey-cost-ledger 8`, `orient-workflow-discovery 5`).

**Step 3 — Codex ACTIVITY (the depth-(ii) activity summary).** Over the same set, a tool-call activity tally from the DB: `exec_command` command count, `update_plan` plan-step presence, `spawn_agent` multi-agent fan-out — a one-line per-session or aggregate activity signal (exact shape firmed at implementation against the step-4 fence).

**Step 4 — render in the step-4 fence.** A Codex section in the body (under PROJECT / alongside WORKSTREAMS) showing the attributed count + the workstream clusters + the activity summary, clearly marked Codex and workdir-attributed (distinct from the name-only `codex-presence` hint). All from the agentsview DB — NO raw-rollout parsing.

**Deferred (NOT in D1):** per-file Codex work-by-area (depth iii) and D3 — both upstream-gated (see Scope decision).

## C / E / F — survey-body-rendering fixes (firmed)

Three more body-rendering fixes fold into 47rx (captain, 2026-06-08), on the same `SKILL.md` + `queries.sql`. Each is independent of D1 and of each other; each bottoms out on a live drive.

**C — drop the step-4 confirm gate.** Today step 4 (`SKILL.md:145-147`) tells the user the headline then STOPS for a yes ("Want me to lay it out?") before rendering the body. **Recommendation: DROP the confirm gate** — the survey is a read-only orientation; the body IS the value, and a one-line headline followed by a stop-and-ask adds a round-trip with no decision behind it (the captain leans this way). Replace the "tell the user … and wait for a yes" + the question with: render the one-line headline, then auto-render the full body in the same turn. (The discovery→commission OFFER at the end of step 4 — "Want me to commission a spacedock workflow from this?" — STAYS; that gate has a real decision behind it. C drops only the *lay-it-out* ceremony, not the commission offer.) Before/after: `SKILL.md:145` "Tell the user what you found and wait for a yes:" + the `> Found … Want me to lay it out?` block → "Lead with the one-line headline, then render the body directly:" + a `> Found {N} sessions … — here's the lay of the land:` headline that flows straight into the synthesis fence.

**E — de-narrate SCAFFOLD discovery.** vh's D3 dropped the recovered/installed TAXONOMY but 1p kept the "recovered from behavior, not files" + "Other recovered one-offs" phrasing (`SKILL.md:125-127`, fence slot `:157`). Captain's intent: do NOT narrate HOW skills were discovered (behavior vs files) — state only the usage fact. Drop "recovered from behavior, not files" and "recovered one-offs"; render usage + checked-in-or-not plainly. Before/after the `:127` example: `superpowers was recovered from behavior, not files: 186 skill invocations, but no checked-in .claude/skills/superpowers. Other recovered one-offs: plan-writing, …` → `superpowers: 186 invocations (not checked in). Other one-offs: plan-writing, using-git-worktrees, systematic-debug, simplify, debugging.` The `:125` prose ("a family invoked but absent from disk was **recovered from behavior, not files** — say exactly that") and the `:157` fence slot are rewritten to the plain usage+presence fact with no discovery-mechanism narration. The file-probe + behavioral-tally JOIN stays (both signals still read); only the narration drops.

**F — Work By Area leads with PRODUCT paths, demotes agent-mechanics.** Today `work-by-area` (`queries.sql` `#317.2`, fence `SKILL.md:165-166`) buckets by first path segment and ORDERS by edit count, so agent-mechanics paths lead: on this repo's real data the lead bucket is `<external>` (358 edits) — and the spike confirmed that is 100% `/tmp`-dispatch/handoff scaffolding, zero real sibling references — burying the product signal (`docs` 279, etc.). Fix: partition each bucket as **product** vs **mechanics** and order product-first. Mechanics = `<external>` (the `/tmp`-dispatch + sibling bucket), `.worktrees` (worktree copies), `.claude`/`.claude-plugin`/`.codex-plugin` (agent plans/scaffolding). Product = everything else under the repo root (`docs`, `internal`, `cmd`, `skills`, `scripts`, `.github`, `dist`, …). The lead lists product buckets (by edit count); the mechanics buckets demote to a one-line footnote ("+ N edits in agent worktrees/plans/dispatch — execution mechanics, not product"). The partition is in the QUERY (a `kind` column + `ORDER BY (kind='mechanics') ASC, edits DESC`) so a query-smoke can gate the ordering; the fence renders product-first with the mechanics footnote.

## Out of scope

- The upstream agentsview Codex-cwd change (D1 does NOT need it — the `workdir` heuristic replaces that dependency for attribution at depth (ii)).
- **Depth (iii) — per-file Codex WORK-BY-AREA.** Deferred upstream-gated (Scope decision): agentsview stores Codex `apply_patch` edits as empty `write_stdin` rows (no `file_path`), so it needs raw-rollout parsing or a `kenn-io/agentsview` ingestion change. NOT in 47rx.
- D3's source-health detector — upstream-gated on agentsview (above).
- Changing the cwd-scoped Claude body mechanism (correct; the Codex signals are added alongside it without altering it).
- **(E) Removing the file-probe / behavioral-tally JOIN.** E drops only the discovery-mechanism NARRATION, not either signal — both are still read and the usage+presence fact is still stated.
- **(F) FILTERING agent-mechanics edits out of the data.** F only RE-ORDERS (product-first) and footnotes the mechanics buckets; the mechanics edits are still counted and shown — the work-by-area stays an honest full accounting, just led by product.

## Acceptance criteria

- **AC-1 — the `codex-scoped` query attributes Codex by `exec_command.$.workdir` prefix: counts this repo's Codex, excludes a same-basename sibling's.** `queries.sql` gains a labeled `-- name: codex-scoped` query that counts `agent='codex'` sessions having an `exec_command` tool call whose `json_extract(input_json,'$.workdir')` is `:repo_root` or `LIKE :repo_root || '/%'`. Verified by: the query-smoke (extended fixture) — fixture Codex session F gets an `exec_command` row with `$.workdir` under `/repo/proj`, sibling session G gets one under `/sibling/proj`; the `codex-scoped` query returns count=1 (F only, G excluded). Independently non-vacuous: changing G's workdir to fall under `/repo/proj` flips the count to 2 (proves the prefix is load-bearing, not a constant). Expected values come from the fixture rows, not skill prose. Fails if the query counts by `project` (would return 2), drops the prefix, or reads the wrong JSON path.
- **AC-2 — the existing `codex-presence` flag is unchanged and the two Codex signals are distinct.** The project-NAME-only `codex-presence` query (and its blank-cwd sum) is untouched; `codex-scoped` is an ADDITIONAL query, not a replacement. Verified by: query-smoke asserts `codex-presence` still returns the project-name-matched count (2 in the fixture: F + sibling G) while `codex-scoped` returns the workdir-attributed count (1: F only) — the two differ on the same fixture, proving they measure different things. Fails if either query is conflated with the other.
- **AC-3 — the workstream-clustering query/rule clusters real `first_message`s into the correct ensign-task workstreams over the `codex-scoped` set, and buckets unlabeled sessions honestly.** The skill's clustering rule (Step 2: dispatch-pattern → `{TASK}`; task/entity-pattern → backtick `{TASK}`; else `(unlabeled)`) groups the attributed Codex sessions by workstream. Verified by: a query-smoke over fixture rows whose `first_message`s are real-shape samples — at least one dispatch-pattern row (`Read /tmp/spacedock-dispatch/spacedock-ensign-journey-cost-ledger-implementation.md`), one task/entity-pattern row (``You are the implementation worker for Spacedock task `orient-workflow-discovery`.``), and one unlabeled row (`You totally got this. Take your time.`) — asserting the cluster output is `journey-cost-ledger → 1`, `orient-workflow-discovery → 1`, `(unlabeled) → 1`. **Non-tautological:** the expected labels (`journey-cost-ledger`, `orient-workflow-discovery`) are SUBSTRINGS of the fixture `first_message` values, never written in SKILL.md — the test extracts them from the message the same way the shipped rule does, so a broken extractor (wrong boundary → label `journey-cost-ledger-implementation`, or the stage not stripped) reds, and the rule cannot pass by spelling. Independently non-vacuous: a dispatch-pattern row with a DIFFERENT task (`…spacedock-ensign-codex-live-ci-validation.md`) must land in a `codex-live-ci` bucket, not be merged with the others — so the cluster key is proven to be the extracted `{TASK}`, not a constant. Expected values come from the fixture message bytes, not skill prose. (Activity tally — `exec_command`/`update_plan`/`spawn_agent` counts over the set — is pinned by the same query-smoke against fixture tool-call rows.)
- **AC-4 — the live survey report renders the Codex WORKSTREAMS + activity in the body (depth ii), surfacing tracks the Claude-only body misses.** Driving `/spacedock:survey` end-to-end on a repo whose synced DB carries `$.workdir`-attributed Codex sessions, the rendered body shows a Codex section with: the workdir-attributed count (distinct from the name-only `codex-presence` flag), the Codex WORKSTREAM clusters, and the activity summary — and the workstreams include real Codex tracks absent from the Claude-only body (e.g. `journey-cost-ledger`, `orient-workflow-discovery`). Verified by: the live-drive rendered report (the survey's proof bar) — a grep over SKILL.md does NOT satisfy this. **The expected workstream labels come from the real session `first_message`s in the synced DB (an independent source the drive reads), NOT from skill prose** — the drive confirms the rendered labels match the workstreams actually present in this repo's Codex history. Per ideation, the live drive may escalate to a captain-run drive (as vh's AC-4/AC-5 did).

- **AC-5 (C) — the survey renders the full body without stopping for a lay-it-out confirm.** Driving `/spacedock:survey` end-to-end, the run emits the one-line headline and then renders the full synthesis body in the SAME turn, with NO intervening "Want me to lay it out?" stop-and-wait. The end-of-report commission OFFER ("Want me to commission…?") still appears (C drops only the lay-it-out ceremony, not the commission gate). Verified by: the live-drive transcript shows the body rendered without an `AskUserQuestion`/stop between the headline and the synthesis, AND the commission offer still present — a grep over SKILL.md does NOT satisfy this; the proof is the observed single-turn render-through. Fails if the run stops to ask before the body, or if dropping the gate also dropped the commission offer.

- **AC-6 (E) — the rendered SCAFFOLD line states usage + checked-in-or-not with NO discovery-mechanism narration.** Driving `/spacedock:survey` on a repo with an invoked-but-not-checked-in scaffold family, the rendered SCAFFOLD section states each family's invocation count and on-disk presence/absence plainly (e.g. `superpowers: 186 invocations (not checked in). Other one-offs: …`) and contains NONE of the discovery-narration phrases "recovered from behavior, not files" / "recovered one-offs" / "recovered". Verified by: the live-drive rendered SCAFFOLD output shows the usage+presence fact and is asserted to NOT contain the banned narration phrases — the expected content (the family + count + presence) comes from the real scan, and the banned-phrase absence is checked on the RENDERED output, not on SKILL.md. (A SKILL.md grep is the wrong target — the proof is the rendered line.) Fails if the rendered line narrates how the scaffold was discovered.

- **AC-7 (F) — Work By Area leads with PRODUCT paths and does NOT lead with agent-mechanics paths.** The rendered WORK BY AREA section lists product buckets (`docs`/`internal`/`cmd`/`skills`/…) first and demotes the agent-mechanics buckets (`<external>`/`/tmp`-dispatch, `.worktrees`, `.claude*`) to a footnote — even when a mechanics bucket has the most edits (on this repo `<external>`=358 outweighs `docs`=279, yet `docs` must lead). Verified by **two** independent checks: (a) a query-smoke over the fixture asserts the `work-by-area` query tags each bucket `product`/`mechanics` and orders product-first (the partition is the independent expected value — a mechanics bucket with a higher edit count still sorts below every product bucket; non-vacuous: flipping a path from a product to a mechanics segment moves it below the lead); (b) a live drive confirms the rendered section LEADS with a product bucket and relegates `<external>`/`.worktrees`/`.claude*` to the footnote. The product-vs-mechanics partition is the independent oracle (fixture path segments), never a SKILL.md grep. Fails if a mechanics bucket leads, or if the partition mislabels a product path as mechanics.

## Test plan

- **Query-smoke (AC-1, AC-2, AC-3, AC-7a) — fixture, cheap (~seconds), reuses the existing sqlite3-driven harness** (`skills/integration/survey_queries_test.go` + `testdata/survey/fixture-sessions.sql`). Extend the fixture:
  - *Attribution (AC-1/AC-2):* add `exec_command` tool_call rows to the existing Codex sessions F (`$.workdir`=`/repo/proj/.worktrees/wt`, under the repo prefix) and G (`$.workdir`=`/sibling/proj`, the same-basename sibling, OUTSIDE the prefix). Add the `codex-scoped` query. Assert `codex-scoped`=1 (F only); `codex-presence` unchanged at 2 (F+G by name); the two differ. Non-vacuousness: repoint G's workdir under `/repo/proj` → `codex-scoped` flips 1→2.
  - *Clustering (AC-3):* give the in-repo Codex sessions real-shape `first_message`s — a dispatch-pattern (`…spacedock-ensign-journey-cost-ledger-implementation.md`), a task/entity-pattern (``Spacedock task `orient-workflow-discovery``), an unlabeled (`You totally got this…`), and a SECOND distinct dispatch task (`…spacedock-ensign-codex-live-ci-validation.md`). Assert the cluster output is `journey-cost-ledger 1`, `orient-workflow-discovery 1`, `codex-live-ci 1`, `(unlabeled) 1` — and (non-vacuousness) that the stage suffix is stripped (NOT `journey-cost-ledger-implementation`) and the two dispatch tasks do NOT merge. Expected labels are substrings of the fixture messages, not SKILL.md text.
  - *Work-by-area partition (AC-7a, F):* the existing fixture already has product (`/repo/proj/internal/…`, `/repo/proj/skills/…`) and mechanics (`<external>` via `/sibling/otherlib/…`) Edit/Write rows; add a `.worktrees`-prefixed and a `.claude`-prefixed row, and ensure a mechanics bucket out-counts a product bucket. Assert the `work-by-area` query returns a `kind` (`product`/`mechanics`) column and orders ALL product buckets before ALL mechanics buckets regardless of edit count. Non-vacuousness: re-segment a `.worktrees` row to `internal` → it moves into the product lead.
  Cost: LOW. TDD: write the failing assertions first (queries/rule missing → FAIL), then add GREEN.
- **Live drive (AC-4, AC-5, AC-6, AC-7b) — minutes; the ONLY proof for the rendered-body behaviors.** ONE end-to-end pass over the real synced DB exercises all four rendered changes at once: seed via the SKILL.md §1 recipe, drive `/spacedock:survey`, and observe in the single rendered run — (AC-4) the Codex section (count + workstream clusters like `journey-cost-ledger`/`orient-workflow-discovery` + activity); (AC-5) the body renders through from the headline with NO lay-it-out stop, commission offer still present; (AC-6) the SCAFFOLD line states usage+presence with none of the "recovered…"/"recovered one-offs" narration; (AC-7b) WORK BY AREA leads with product (`docs`/…) and footnotes `<external>`/`.worktrees`/`.claude*` (on this repo `<external>`=358 must NOT lead). Reuses the spike's sync recipe; agentsview v0.32.1 present. Expected values come from the real survey output / scan numbers, not skill prose. May escalate to a captain-run drive (as vh's AC-4/AC-5 did). Cost: LOW–MEDIUM.
- **No depth-(iii) / D3 test** — both deferred/upstream-gated; nothing ships for them, so they have no AC (a decision-with-nothing-shipped does not belong in the queue — the deferrals + the upstream dependency are recorded in the Scope decision for the roadmap).
- **Estimated total cost/complexity: LOW–MEDIUM.** D1: `codex-scoped` query + clustering rule + activity tally + fixture rows. C: drop the confirm-gate prose. E: rewrite the SCAFFOLD narration to plain usage+presence. F: a `kind` partition + product-first ordering in the `work-by-area` query + the footnote in the fence. All four are `SKILL.md`/`queries.sql` body-rendering changes proven by the shared query-smoke + ONE combined live drive — NO agentsview/binary/rollout-parse changes (those are the deferred depth-(iii) path).

## Notes

Provenance: vh feedback cycle 1 + the captain's 2026-06-08 live-drive findings, plus THIS ideation's D1 spike + rendered content views + the captain's C/E/F body-rendering observations. Captain decisions (2026-06-08): D2 DROPPED; D1 firmed to DEPTH (ii); C/E/F folded in as additional body-rendering deliverables; depth (iii) + D3 banked upstream-gated. 47rx is the survey-skill body-rendering pass (D1 + C + E + F), all on `SKILL.md`/`queries.sql`. The D1 spike OVERTURNS the original upstream-gated premise (the `$.workdir` field makes attribution + the depth-(ii) content DB-derivable). C/E/F are grounded in real data: F's `<external>`-leads problem reproduced (358 edits, 100% `/tmp`-dispatch noise, burying `docs`=279) and the product-first partition prototyped in-query. Per the survey discipline, a grep over SKILL.md never satisfies a behavioral AC — every deliverable bottoms out on the combined live drive (AC-4/5/6/7b) plus the query-smoke for the gate-able query halves (AC-1/2/3/7a), whose expected values are fixture/scan-derived, not skill prose.

## Stage Report: ideation

- DONE: Scope decision (the central output): determine whether 47rx ships as ONE task or SPLITS, and what is buildable in 0.19.9 vs upstream-gated — D1, D2, D3.
  Recorded in "Scope decision". Verdict: D1 SHIPPABLE in 0.19.9 (the `workdir` heuristic — spike overturned the upstream-gated premise), D2 SHIPPABLE in 0.19.9 (small), D3 UPSTREAM-GATED (no agentsview source-health signal — confirmed). D1+D2 ship as ONE task (same step-4 fence + queries.sql + fixture, same as vh's four); D3 deferred. Evidence is the two spikes below.
- DONE: Exercise the riskiest mechanism FIRST — D1's attribution path — end-to-end, OR record D1 as upstream-gated with evidence.
  Spiked end-to-end against live agentsview v0.32.1 ("Spike: D1 attribution heuristic"). Found `exec_command.$.workdir` carries the real absolute cwd that the session-level `cwd` lacks. Measured on real data: this repo (`spacedock_v1`, 61 Codex) vs same-basename sibling (`spacedock`, 160 Codex) — workdir-prefix scope gives 57/61 recall and ZERO sibling cross-contamination (0/160 sibling sessions under this prefix). Field present on 95% of all 60k Codex exec_command calls. D1 is viable WITHOUT the upstream Codex-cwd change.
- DONE: Produce AC + test plan proven by a LIVE DRIVE of the survey behavior — never a grep over SKILL.md.
  Three ACs (AC-1 `codex-scoped` workdir-attribution query proven by query-smoke over the extended fixture + non-vacuousness flip; AC-2 the two Codex signals stay distinct; AC-3 the live-drive rendered D1 Codex line + D2 codex-presence=0 caveat) + a 4-bullet test plan. AC-1/AC-2 oracle from independent fixture rows; AC-3 from real survey output. D3 has no AC (deferred, nothing ships). The D2 codex-presence=0 case was confirmed against the real DB (`0|0` for a non-matching project; current fence renders nothing → the silent confident 0).

### Summary

47rx SPLITS by verdict, ships D1+D2 together in 0.19.9, defers D3 as upstream-gated. The linchpin was the D1 attribution spike: I overturned the task body's premise that surfacing Codex needs the upstream agentsview Codex-cwd change by finding that `exec_command` tool calls persist a `$.workdir` field — the absolute working directory the session-level `cwd` drops. On real data it separates this repo's 61 Codex sessions from the same-basename sibling's 160 with zero cross-contamination (the exact collision vh documented), at 93% recall on a 95%-present field. D2 (the honest caveat on a real codex-presence=0) is small and confirmed; D3 (source-health detector) has no in-0.19.9 path — agentsview exposes no per-source signal and the FS-access probe is Seatbelt-denied. ACs are proven by the query-smoke (independent fixture oracle) and a live drive (rendered report), per the survey discipline — no SKILL.md grep stands as a behavioral AC.

## Stage Report: ideation (cycle 2)

Captain follow-up: two concrete D1-feasibility exercises (combined view + codex-only mapping) folded into the spike. Run against the same live agentsview v0.32.1 DB. They refined the D1 scope call.

- DONE: Exercise A — combined Claude+Codex view over this repo; capture the rendered output.
  Recorded in Spike Part 2 / Exercise A. The NAIVE combined fold (`agent IN ('claude','codex')` on the work-by-area query) renders a view IDENTICAL to Claude-only — Codex contributes ZERO rows (0 structured Edit/Write + blank cwd). Rendered combined work-by-area: `<external>` 357, `docs` 276, … all `agent=claude`. The only coherent combined signal is the COUNT line: Claude 51 (cwd-scoped) · Codex 58 (workdir-attributed) · 61 name-only. Verdict: the combined view is coherent ONLY at the count level; folding Codex work CONTENT into the body adds nothing without upstream ingestion.
- DONE: Exercise B — codex-only mapping check; is Codex session data meaningful on its own?
  Recorded in Spike Part 2 / Exercise B. COUNT + attribution: clean (58 workdir-attributed, sibling-free). Tool-call mapping: Codex's tools are `exec_command`(6745)/`write_stdin`(193)/`update_plan`(126)/`spawn_agent`… — **0 Edit/Write/Read/apply_patch** (the columns survey body queries read); dominant verb is `rtk`(5883, CL's token-proxy) wrapping git/sed/go/rg READS. Edits via `apply_patch` are stored by agentsview as EMPTY `write_stdin` rows (no `file_path`); edited paths exist only as `*** Update File:` text in raw rollouts (parsed 484 targets → internal/docs/skills, but a plain grep returns garbage). Verdict: codex-only COUNT is meaningful; codex-only WORK-BY-AREA/SCAFFOLD/DECISION is NOT cleanly derivable from agentsview today.
- DONE: Refine the D1 scope call from the two exercises.
  D1 SPLITS: **D1-count** (attributed Codex COUNT line) ships in 0.19.9 — the workdir heuristic is clean; **D1-body** (Codex work/scaffold/decision content folded into the body) DEFERS as upstream-gated — agentsview ingests no structured Codex edit/skill/decision signal. Updated the Scope decision, firmed approach, Out of scope, and AC-3 to ship D1-count + D2, defer D1-body + D3. The ACs were already count-scoped, so no AC was weakened.

### Summary (cycle 2)

The captain's two exercises converted my earlier "D1 buildable" into a sharper split. Attribution (Part 1) is clean — workdir separates this repo's Codex from a same-basename sibling with zero leak — so a Codex COUNT in the body (D1-count) ships in 0.19.9. But surfacing Codex WORK CONTENT (D1-body) is upstream-gated: Exercise A showed the naive combined fold is byte-identical to Claude-only (Codex contributes nothing), and Exercise B showed why — agentsview captures 0 structured Codex file edits (apply_patch → empty write_stdin), Codex's recoverable signal is overwhelmingly `rtk`-proxied reads, and edited paths survive only as raw-rollout text the survey doesn't (and shouldn't) parse. Net 0.19.9 scope: D1-count + D2 ship as one task; D1-body + D3 defer as upstream-gated on a `kenn-io/agentsview` change. The rendered combined view and codex-only mapping are both in the spike section as the evidence the call rests on.

## Stage Report: ideation (cycle 3)

Captain update (2026-06-08): DROP D2; RENDER the two content views; HOLD D1 scope pending captain review.

- DONE: Drop D2 (silent-0 honesty caveat) entirely from scope, body, AC, and test plan.
  Removed D2 from the intro, Problem, Scope decision, firmed approach, Out of scope, AC-3, test plan, and Notes. 47rx is now D1 (Codex body-surfacing) + D3 (deferred). No D2 query/caveat/AC remains. The D3-spike's stale "D2 is the shippable honesty answer" line was corrected.
- DONE: Render the COMBINED view (Claude + Codex folded) and the CODEX-ONLY view as actual CONTENT, not counts.
  Spike Part 2 rewritten as "RENDERED CONTENT VIEWS" against the real synced DB (57–58 workdir-attributed sessions). CODEX-ONLY: (B-1) workstreams from first_message cluster coherently into ensign tasks (journey-cost-ledger 8, codex-foreground-wait-escape-hint 5, orient-workflow-discovery 5, …); (B-2) tool-call mapping is CORRECT (exec_command/update_plan/spawn_agent map coherently — update_plan even carries structured plan steps), NOT garbage; (B-3) work-by-area (internal 247, skills 48, docs/.spacedock-state 85) is coherent but RAW-ROLLOUT-derived (agentsview stores Codex apply_patch as empty write_stdin). COMBINED: (A-1) the two agents did complementary work — Claude's visible buckets are /tmp-dispatch + docs orchestration, Codex's are internal/skills implementation — so the combined view tells a fuller story; (A-2) the naive structural fold adds nothing (Codex needs its own signals folded in).
- DONE: Mark D1's scope (count-only vs content-surfacing) PENDING captain review; do NOT lock the re-scope.
  Scope decision rewritten: D1 attribution is PROVEN; the CONTENT-DEPTH is the captain's pending call across three rendered depths — (i) count, (ii) +workstreams/activity (DB-only, my suggested sweet spot), (iii) +file work-by-area (adds rollout parsing). The firmed approach, AC-3, and test plan are written for the settled foundation + depth (i)/(ii) and explicitly flag that the FO's post-views feedback round firms the chosen depth.

### Summary (cycle 3)

Dropped D2 per the captain and rendered the two CONTENT views the captain asked to see. The headline correction to my cycle-2 framing: the codex-only view is COHERENT, not garbage — workstreams (ensign-task clusters from first_message) and tool-call mapping (exec_command/update_plan/spawn_agent, including structured update_plan steps) render correctly straight from the agentsview DB; only the file-level work-by-area needs raw-rollout parsing. The combined view is additive and complementary: Claude's visible work is /tmp-dispatch + docs orchestration while Codex carries the internal/skills implementation, so folding Codex in tells a fuller story than the Claude-only body. D1's depth (i count / ii +workstreams+activity / iii +work-by-area) is left UNLOCKED for the captain; the FO will run a feedback round to firm the chosen depth. D1 attribution itself is proven; D3 stays deferred/upstream-gated.

## Stage Report: ideation (cycle 4)

FO feedback round: captain picked DEPTH (ii). Firmed the ideation to it (independent staff review comes after).

- DONE: Lock D1 to depth (ii) — Codex workstreams + activity via Codex's OWN DB signals (first_message → clusters; update_plan/exec_command → activity), workdir-attributed, DB-only (NO rollout parsing).
  Scope decision, intro, Problem, combined-verdict, and the firmed approach all rewritten from "PENDING/captain-picks" to depth-(ii)-LOCKED. The firmed approach now spells the build in 4 steps (codex-scoped query → workstreams → activity → step-4 fence) with the clustering rule made EXPLICIT (3 ordered cases: dispatch-pattern → {TASK}; task/entity-pattern → backtick {TASK}; else (unlabeled)).
- DONE: Defer depth (iii) + D3 as banked upstream-gated follow-ups; record the upstream-agentsview dependency so the exploration's value is kept.
  Scope decision banks depth (iii) (per-file work-by-area, needs fragile rollout apply_patch parsing) and D3 (source-health) with a concrete recorded upstream `kenn-io/agentsview` dependency: (a) persist Codex cwd, (b) ingest structured Codex file edits, (c) expose per-source health — each tied to the follow-up it unblocks. The rendered depth-(iii) view (B-3) is banked as proof-of-value.
- DONE: Scrub D2 entirely (last net-call message still said "D1-count + D2").
  No D2 query/caveat/AC remains in scope, body, AC, or test plan (grep-verified: only the intended "D2 dropped" provenance notes survive). The stale "D1-count"/"D1-body" labels are gone (cleaned cycle 3). Net call is now D1(ii) + deferrals — no D2.
- DONE: Firm AC to depth (ii), proof-policy-clean, with a NON-TAUTOLOGICAL workstream-clustering AC.
  AC-1 (codex-scoped $.workdir attribution, sibling-excluded, non-vacuous) carries over. NEW AC-3 (clustering): the rule clusters real-shape fixture first_messages into the correct task labels — expected labels are SUBSTRINGS of the fixture message bytes (independent source), never SKILL.md text, so a broken extractor (stage not stripped, wrong boundary) reds and the two distinct dispatch tasks must not merge. AC-4 (live drive) proves the rendered Codex workstreams in the body, labels from the real synced DB's first_messages. Validated the clustering SQL by exercising it against fixture-shape messages — dispatch/unlabeled/second-task extract correctly; caught + corrected a backtick off-by-one in the task/entity rule (the implementation's first TDD test will pin it).

### Summary (cycle 4)

Firmed 47rx to the captain's depth (ii): the survey body surfaces Codex WORKSTREAMS (clustered from first_message) + activity (update_plan/exec_command tallies) over the $.workdir-attributed, sibling-free Codex set — all from the agentsview DB, no rollout parsing. Depth (iii) (per-file work-by-area) and D3 (source-health) are banked as upstream-gated follow-ups with a concrete recorded `kenn-io/agentsview` dependency (persist Codex cwd / ingest structured edits / per-source health), so the exploration's value is preserved as actionable next steps. D2 is fully scrubbed. The ACs are firmed and proof-policy-clean: attribution + clustering by query-smoke over an extended fixture whose expected labels are message-byte substrings (non-tautological — a SKILL.md grep cannot satisfy them), and the rendered body workstreams by a live drive whose labels come from the real DB. I exercised the clustering SQL against real fixture-shape messages to confirm the mechanism is sound (and corrected a backtick off-by-one). Ready for independent staff review.

## Stage Report: ideation (cycle 5)

Captain folded three more survey-body-rendering fixes (C, E, F) into 47rx on top of the firmed D1 depth-(ii). 47rx is now the survey-skill body pass: D1 + C + E + F.

- DONE: Fold C/E/F as additional deliverables ON TOP OF D1 depth-(ii) (kept intact), same SKILL.md + queries.sql, each with its own AC, each bottoming out on a live drive.
  Added the "C / E / F — survey-body-rendering fixes (firmed)" section with concrete before/after wording for each; added AC-5 (C), AC-6 (E), AC-7 (F, two-part: query-smoke partition + live-drive ordering); extended the test plan (F's partition joins the query-smoke; C/E/F live-drives join the ONE combined live drive); updated intro, Out of scope, Notes. D1's AC-1..AC-4 and the depth-(ii) approach are unchanged.
- DONE: C — evaluate dropping the step-4 confirm gate; recommend.
  RECOMMEND DROP. The survey is read-only orientation; the one-line headline → stop-and-ask "Want me to lay it out?" is a round-trip with no decision behind it. Drop it; auto-render the body after the headline. The end-of-report commission OFFER stays (real decision). AC-5 proves the single-turn render-through + commission-offer-still-present by live drive.
- DONE: E — drop the discovery-mechanism narration; state usage only.
  Firmed the before/after: `recovered from behavior, not files: 186 … Other recovered one-offs: …` → `superpowers: 186 invocations (not checked in). Other one-offs: …`. The file-probe+tally JOIN stays (both signals read); only "recovered from behavior, not files"/"recovered one-offs" narration drops. AC-6 asserts the rendered SCAFFOLD line has usage+presence and NONE of the banned phrases (checked on rendered output, not SKILL.md).
- DONE: F — Work By Area leads with PRODUCT, demotes agent-mechanics; partition is the independent expected value.
  Grounded in real data: reproduced the problem (`<external>`=358 leads, and proved that bucket is 100% `/tmp`-dispatch noise; `docs`=279 buried) and PROTOTYPED the fix in-query — a `kind` (product/mechanics) column + `ORDER BY (kind='mechanics') ASC, edits DESC` puts product first even when a mechanics bucket out-counts it. Mechanics = `<external>`/`.worktrees`/`.claude*`; product = the rest under repo root. AC-7a (query-smoke: partition + product-first ordering, non-vacuous) + AC-7b (live drive: rendered lead is product, mechanics footnoted).

### Summary (cycle 5)

47rx grew from Codex-only into the full survey-skill body-rendering pass: D1 (Codex depth ii) + C (drop the lay-it-out confirm gate) + E (de-narrate SCAFFOLD discovery) + F (Work By Area leads with product, demotes agent-mechanics). The D1 depth-(ii) firming is intact; C/E/F are independent body-rendering changes on the same SKILL.md/queries.sql, each with its own AC and all bottoming out on ONE combined live drive (plus the shared query-smoke for the gate-able query halves — D1's codex-scoped/clustering and F's product/mechanics partition). C is a recommend-DROP with the commission offer preserved; E drops only the discovery narration (the JOIN stays); F only re-orders (mechanics edits still counted, just footnoted) — each scoped to avoid over-reach. F's partition was prototyped in-query against real data (product leads despite `<external>`=358 outweighing `docs`=279), so the design rests on an exercised mechanism. AC-7's product-vs-mechanics partition and AC-3's workstream labels are both non-tautological (fixture/scan-derived oracles, never a SKILL.md grep). Ready for the independent staff review.

## Stage Report: implementation

- DONE: D1 (depth ii) — `codex-scoped` workdir-attribution query (AC-1) + workstream CLUSTERING as a labeled query block (AC-3) + activity tally; render the Codex section in step-4 fence.
  Added `-- name: codex-scoped` (#321, `exec_command.$.workdir` prefix), `-- name: codex-workstreams` (#322, the 3-case clustering rule AS a runnable query block, not prose), `-- name: codex-activity` (#323) to `references/queries.sql`. SKILL.md step-2 runs all three; step-4 fence renders a CODEX section. Worktree commit f0f1b9f6.
- DONE: C — drop the step-4 lay-it-out confirm gate; auto-render body; KEEP commission offer.
  SKILL.md §4 heading `Confirm, then report and offer` → `Report and offer`; the "wait for a yes" + "Want me to lay it out?" question replaced with a render-through headline (`— here's the lay of the land:`) flowing straight into the fence. Commission offer at §4 close preserved.
- DONE: E — de-narrate SCAFFOLD discovery; KEEP the file-probe/behavioral-tally JOIN.
  SKILL.md §3 join-prose + fence slot rewritten to plain usage+presence (`superpowers: 186 invocations (not checked in)`); dropped "recovered from behavior, not files"/"recovered one-offs"/"recovered" from the rendered surface (also cleaned the stray term-of-art in queries.sql + fixture comments). Both probe signals still read.
- DONE: F — work-by-area product-first, mechanics footnoted (edits STILL counted, just re-ordered).
  `work-by-area` query gains a `kind` (product/mechanics) column + `ORDER BY (kind='mechanics') ASC, edits DESC` (AC-7a). Mechanics = `<external>`/`.worktrees`/`.claude*`. Fence renders product buckets first + a mechanics-sum footnote.
- DONE: Proof — query-smoke (AC-1/2/3/7a) green, TDD failing-first then green; ONE combined live drive (AC-4/5/6/7b) over the real synced DB.
  Query-smoke: `skills/integration/survey_queries_test.go` 11/11 (full integration 23/23), fixture extended with Codex F/F2/F3/F4 + G exec_command/update_plan/spawn_agent rows and `.worktrees`/`.claude` work-by-area rows; non-vacuous (G-workdir-repoint flips codex-scoped 4→5; re-segment a `.worktrees` row → product). Live drive against real agentsview v0.32.1 (`spacedock_v1`, 53 Claude + 61 codex-presence): codex-scoped=58 (sibling-free, distinct from 61); workstreams `journey-cost-ledger 8`/`orient-workflow-discovery 5`/… `(unlabeled) 15` last (labels from real `first_message`s); activity `exec_command 6834`/`update_plan 136`/`spawn_agent 54`; SCAFFOLD `superpowers: 6 invocations (not checked in)` with NO banned narration (asserted on rendered output); WORK BY AREA leads `docs 279` with `<external> 375` (highest count) demoted to footnote; headline renders through to the body with no confirm stop, commission offer present. Did NOT need to escalate to a captain-run drive.
- SKIPPED: depth (iii) (per-file Codex work-by-area) + D3 (source-health).
  Out of scope (upstream-gated on a `kenn-io/agentsview` change) per the captain-locked ideation; nothing shipped for them.

### Summary

Shipped the survey-skill body-rendering pass: D1 surfaces Codex at depth (ii) as its own body section (workdir-attributed count + `first_message` workstream clusters + `exec_command`/`update_plan`/`spawn_agent` activity), all from the agentsview DB with no rollout parsing; C drops the lay-it-out confirm gate (commission offer kept); E de-narrates SCAFFOLD to a plain usage+presence fact (JOIN kept); F re-orders work-by-area product-first with a `kind` partition (mechanics still counted, footnoted). The LOAD-BEARING clustering (AC-3) and product/mechanics partition (AC-7a) live as runnable labeled query blocks in `queries.sql`, so the query-smoke executes the artifact — not a SKILL.md grep. Query-smoke is green TDD (failing-first, fixture-derived non-vacuous oracles), and the ONE combined live drive over the real synced DB proves the four rendered-body behaviors (AC-4/5/6/7b) with expected values from real survey output, not skill prose. Deliverable committed on the worktree branch (f0f1b9f6); four files: SKILL.md, queries.sql, the test, the fixture. Routes to independent validation.
