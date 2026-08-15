---
name: survey
description: Use when arriving at or returning to a project that already has AI-agent history and you want the lay of the land before doing anything else — "survey this project", "what have we been doing here", "catch me up", "orient me", "where did we leave off", "what's the state of this project", or picking up a brownfield repo with several in-flight agent tracks. Reads existing agent session history (read-only), reports the implicit workflow, the open decisions, and how often you had to step in, then offers to commission a spacedock workflow from it.
user-invocable: true
---

# Survey a Project

## Overview

Survey is the first thing you run on unfamiliar ground: it reconstructs what the AI agents in this project have implicitly been doing, from their session history. It reports the inferred workflow, the workstreams, the recent decisions, and — load-bearing — the OPEN decisions (the abandoned or unanswered forks) plus how often the human had to step in. Then it offers to commission a real spacedock workflow with explicit gates from what it found.

It reads **agentsview**'s session DB and is strictly read-only — the recommended queries live in `references/queries.sql` (one labeled query per concern) so nothing is a black box. The decision and interruption signals below are **Claude Code**'s; **Codex** is surfaced too, as its own body section (a workdir-attributed count + workstream clusters + activity), since Codex sessions land with no recorded cwd and need the `exec_command.$.workdir` signal to be scoped to this repo. Gemini and per-file Codex work-by-area remain deferred follow-ups. The report opens with plain value and concrete numbers from the user's own data, then demotes the mode/track vocabulary to a detail section below. The closing move is the discovery → commission bridge: the OPEN decisions become candidate gates, the workstreams become candidate entities, the inferred loop becomes the stage list — and the offer leads with plain value ("this helps you run the repetitive work and stop only at the calls you'd want to make"), keyed to each track's MODE (gate-and-drive for the manual-but-repetitive tracks, book-keeping for exploration and knowledge-work tracks).

Run the four steps in order: **check agentsview → scan → recognize scaffold → report and offer**.

---

## 1. Check agentsview, then sync THIS project only (scoped)

This skill may run in a sandboxed agent that **cannot read `~/.agentsview/` directly** (macOS TCC denies raw FS access to a limited-permission process, even though the `agentsview` binary itself reads it). So do NOT `sqlite3 ~/.agentsview/sessions.db` blindly. Instead, drive the read through the `agentsview` binary into a process-readable data directory under `AGENTSVIEW_DATA_DIR`, then query that copy.

Probe for the binary by **invoking it**, not by walking `PATH`. A `command -v agentsview` (or `which`/`test -x`/`stat`) is an FS-access lookup the sandbox's Seatbelt can deny while still allowing `execve` — so the lookup false-negatives even though `agentsview` runs fine. `agentsview --version` is the execve that survives whatever the survey's real through-the-binary reads survive; `>/dev/null 2>&1` suppresses both the present banner and the absent "not found" so only the `AGENTSVIEW MISSING` sentinel reaches the agent (silent ⇒ present; sentinel ⇒ absent):

```bash
SURVEY_DB_DIR="${SPACEDOCK_SURVEY_DB_DIR:-${TMPDIR:-/tmp}/spacedock-survey}"
DB="$SURVEY_DB_DIR/sessions.db"

if ! agentsview --version >/dev/null 2>&1; then echo "AGENTSVIEW MISSING"; fi
```

If it prints `AGENTSVIEW MISSING`: tell the user agentsview is needed (it ingests the agent logs this skill reads), **ask consent**, and only on a yes run the install (`brew install --cask agentsview`; fallback `curl -fsSL https://agentsview.io/install.sh | bash`). NEVER install without an explicit yes — stop at the consent prompt otherwise.

With the binary present, sync — but **scope the CLAUDE source to this project**. A bare `agentsview sync` enumerates the ENTIRE `~/.claude/projects` history (16k+ sessions, growing): on a real machine that walk exhausts any sane timeout, so the survey data dir ends up empty or partial. The fix is to narrow Claude's source root to just this repo's session directories before syncing, so the walk is bounded and this project's Claude sessions land in seconds. The narrowing applies ONLY to Claude (`CLAUDE_PROJECTS_DIR`); the same `agentsview sync` ALSO walks the other agents' default source dirs unscoped — notably Codex (`~/.codex/sessions`), which is what populates the `codex-presence` count in step 2 (Codex is a small backlog, so leaving it unscoped is cheap). You always run this sync as part of the survey — see below.

```bash
mkdir -p "$SURVEY_DB_DIR"

# Claude Code stores each project's sessions in ~/.claude/projects/<cwd-with-/-as->,
# so this repo's sessions live under dirs that begin with the dash-encoded cwd. Point
# CLAUDE_PROJECTS_DIR at a symlink farm of just those dirs — the sync then walks only
# this project's Claude sessions (this cwd plus its worktrees), not the whole backlog.
CLAUDE_ROOT="${CLAUDE_PROJECTS_DIR:-$HOME/.claude/projects}"
DASH_CWD=$(pwd | sed 's#/#-#g')                      # ~/.claude/projects dir-name convention
NARROW="$SURVEY_DB_DIR/claude-narrow"
rm -rf "$NARROW"; mkdir -p "$NARROW"
# Enumerate this project's Claude session dirs with find, not a shell glob: an unmatched
# `…-*` glob is a HARD error under zsh's default nomatch (and the skill runs under
# whatever shell the user has). find matches the exact dash-encoded cwd dir plus every
# `<cwd>-*` worktree sibling, and an empty match links nothing without erroring — so the
# "no Claude sessions for this project" path still flows to step 2's "no agent history".
find "$CLAUDE_ROOT" -maxdepth 1 -type d \( -name "$DASH_CWD" -o -name "$DASH_CWD-*" \) -print0 2>/dev/null |
while IFS= read -r -d '' d; do
  ln -s "$d" "$NARROW/$(basename "$d")"
done

AGENTSVIEW_DATA_DIR="$SURVEY_DB_DIR" CLAUDE_PROJECTS_DIR="$NARROW" timeout 300 agentsview sync
```

**Always run this sync as part of the survey — never query a pre-existing `SURVEY_DB_DIR` without re-syncing first.** The data dir persists between runs so the re-sync is incremental (seconds), but it is the sync that refreshes BOTH this project's Claude sessions AND the Codex backlog into the DB. A persisted dir left over from an earlier survey can be missing newer sessions — or, if it predates Codex on disk, missing Codex entirely — so the step-2 `codex-presence` count would read a stale 0. The incremental sync backfills them (no `--full` needed); skipping the sync is what produces a wrong 0. Do not pass `--full` — a full resync re-ingests everything and can fill the disk. If the symlink farm is empty (this project has no Claude sessions under `~/.claude/projects`), the synced DB has no Claude history for it; step 2 reports "no agent history" and stops.

If `agentsview sync` fails (network, disk, permissions), report the exact failure and stop — do not fall back to raw `~/.agentsview/` reads (they fail under TCC).

## 2. Scan the project

**Scope by repo IDENTITY, not by `project` name.** agentsview keys each session's `project` by the **git-root basename** (normalized: non-alphanumerics → `_`), so EVERY checkout of one repo — the root, a subdir, a worktree, the split-root state dir — already shares the SAME `project` key. The risk is the INVERSE: that basename key COLLIDES with a same-basename sibling repo elsewhere on disk, so a `project = <basename>` filter would fold an unrelated repo's sessions into this one. Resolve the repo root once, then scope every query by the absolute cwd-prefix so a same-basename sibling stays out:

```bash
# Repo-root identity: the parent of the common .git dir resolves to the SAME absolute
# path from the repo root, a subdir, the state checkout, or a linked worktree — so the
# cwd-prefix admits every checkout of THIS repo while excluding a same-basename sibling
# elsewhere on disk. --path-format=absolute needs git >= 2.31; if this is not a git repo
# or git is older, fall back to the cwd itself (today's behavior).
if REPO_GIT=$(git rev-parse --path-format=absolute --git-common-dir 2>/dev/null); then
  REPO_ROOT=$(dirname "$REPO_GIT")
else
  REPO_ROOT=$(pwd)
fi
```

Bind `REPO_ROOT` as the `:repo_root` parameter of the queries in `references/queries.sql`. The session-scoping subquery there matches `cwd = :repo_root OR cwd LIKE :repo_root || '/%'` — the cwd AT the root and everything strictly under it. Because `project` collides across same-basename repos, this absolute cwd-prefix is the ONLY column that keeps a same-basename sibling like `…/other/proj` out of scope. Run a labeled query by name (the SQL stays in the reference file; this only extracts and runs it):

```bash
SURVEY_SKILL_DIR="${SPACEDOCK_SURVEY_SKILL_DIR:-skills/survey}"
QUERIES="$SURVEY_SKILL_DIR/references/queries.sql"

# :repo_project is the git-root basename agentsview keys `project` by, with the SAME
# normalization agentsview applies (non-alphanumerics → `_`, e.g. spacedock-v1 →
# spacedock_v1). It is read ONLY by codex-presence (Codex sessions land cwd='', so they
# can't be cwd-prefix-scoped — they match by project name alone).
REPO_PROJECT=$(printf '%s' "$(basename "$REPO_ROOT")" | tr -c '[:alnum:]' '_')

run_query() {  # run_query <name> — :repo_root → REPO_ROOT, :repo_project → REPO_PROJECT
  local q
  q=$(awk -v n="$1" '$0=="-- name: "n{f=1;next} /^-- name: /{f=0} f && $0 !~ /^--/{print}' "$QUERIES")
  printf ".param set :repo_root '%s'\n.param set :repo_project '%s'\n%s\n" "$REPO_ROOT" "$REPO_PROJECT" "$q" | sqlite3 "$DB"
}

run_query scoping            # sessions|blank_cwd|span over the cwd-prefix-scoped repo
run_query codex-presence     # flagged Codex count|blank_cwd by project NAME (cwd unrecorded)
run_query codex-scoped       # Codex attributed by exec_command.$.workdir prefix (sibling-free)
run_query codex-workstreams  # cluster codex-scoped sessions into ensign-task workstreams
run_query codex-activity     # exec_command/update_plan/spawn_agent tally over the codex-scoped set
run_query scaffold-usage     # behavioral skill_name family tally (spacedock self EXCLUDED)
run_query work-by-area       # Edit/Write file_path → LOGICAL area (worktree prefix stripped) + kind
run_query decision-open      # AskUserQuestion/ExitPlanMode frontier; OPEN sorts first
run_query decision-no-followup # done decisions with no later Edit/Write (message_id → ordinal join)
run_query mode-classification # classify each git_branch track manual/exploration/knowledge-work/unlabeled
run_query dispatch-fact      # orchestration FACT: distinct in-repo parents that dispatched subagents + total dispatched
```

`scoping` returns `sessions=0` → there is no Claude agent history for this repo; say so and stop. Nothing to discover. (Survey reads Claude history only for now; a repo whose only agent history is Codex/Gemini will report "no agent history" here — surfacing those agents is a deferred follow-up.) Note the `blank_cwd` count in the report if non-zero (sessions agentsview never captured a cwd for, which the repo-root scope cannot place).

`codex-presence` returns the count of `agent='codex'` sessions matching this repo's `project` name, plus the blank-cwd sum among them. Because agentsview does not persist Codex cwd, matches are by git-root-basename `project` only — a same-basename sibling repo will collide, so the count may include unrelated Codex sessions. Treat it as a presence flag only; these sessions are never counted in the Claude `scoping`, `scaffold-usage`, or `work-by-area` sets. The body Codex section is built from the `codex-scoped` set below instead — when `codex-presence > codex-scoped`, the gap is the name-only superset (possible sibling), and the report says so.

**Codex body signals (`codex-scoped`, `codex-workstreams`, `codex-activity`).** agentsview persists no Codex session cwd, but a Codex session's `exec_command` tool calls carry `$.workdir` (the absolute working directory of each shell command). `codex-scoped` attributes a Codex session to THIS repo when it has an `exec_command` whose `$.workdir` is under the repo-root prefix — the workdir analogue of the Claude cwd-prefix scope, so it admits this repo's Codex and EXCLUDES a same-basename sibling (whose workdirs fall under a different prefix). Over that sibling-free set, `codex-workstreams` clusters the sessions into ensign-task workstreams from each `first_message` (the runnable 3-case rule: dispatch-file read → `{TASK}` with the stage suffix stripped; `Spacedock task/entity` backtick → the backtick-quoted `{TASK}`; else `(unlabeled)`), and `codex-activity` tallies the per-tool activity (`exec_command`/`update_plan`/`spawn_agent`). These are the Codex body section (step 4) — the workstreams surface real Codex tracks the Claude-only body misses. All from the agentsview DB; no raw-rollout parsing. (Per-file Codex work-by-area and a source-health signal are deferred — they need an upstream agentsview ingestion change.)

**Track modes (`mode-classification`).** `mode-classification` labels each `git_branch` track `manual` (low veto + gate-pass-dominant + issue→worktree→PR loop markers + code edits — the repetitive-but-substantive drive loop; the label is `manual`, not `mechanical`, because the work is effortful, not trivial — reserve "mechanical" for genuinely trivial edits, which the classifier does not separately detect today), `exploration` (high veto + a rejected/cancelled path + prose/`.md` edits), `knowledge-work` (an intake→process→file→log→close loop + content/ops edits — `.md` and data — dominating code, with no issue→PR loop and no veto-heavy creative signature: a notes/ops shop, not a code repo), or `unlabeled` (none clearly dominant). The report reads the per-track `mode` to label WORKSTREAMS and to pick the right commission offer per track (step 4) — gate-and-drive for manual, thread book-keeping for exploration, batch book-keeping for knowledge-work, generic book-keeping for unlabeled (never a guessed automation pitch).

**Honest signal accounting.** The `decision-open` rows are the human-decision points; `OPEN` = still needs the human, and you lead the report with those. For the interruption total, count the AskUserQuestion / ExitPlanMode decisions plus the hard-veto markers Claude sessions retain (`[Request interrupted` / `Request interrupted by user` / `doesn't want to proceed` in the message stream), over the same repo-scoped session set; `pct = total*100/user_turns`. Never dress an empty section up as "no decisions" — if a section is empty, say the run found none of that signal.

**No-follow-up decisions (`decision-no-followup`).** This counts the `done` decisions (answered AskUserQuestion / approved ExitPlanMode) that had NO Edit/Write LATER in the same session — a call you settled and then built nothing on. "Later" is the real chronological order (`tool_calls.message_id → messages.ordinal`), not insertion order. It is DISTINCT from BACKLOG (decided-not-shipped, a transcript fork with no repo artifact): this is a same-session gap. It renders as one `BY THE NUMBERS` line.

**Dispatch fact (`dispatch-fact`).** The body EXCLUDES dispatched-subagent sessions, so an orchestrated repo (most work inside subagents) reads as nearly idle. `dispatch-fact` counts the DISTINCT in-repo parent sessions that dispatched subagents and the total dispatched, joining each subagent (`relationship_type='subagent'`) to its parent and keeping only in-repo parents. It surfaces the FACT of orchestration, never subagent CONTENT — the line is dropped when no in-repo session orchestrated.

## 3. Recognize an incumbent scaffold

Recognize the scaffold from TWO signals and reconcile them — a file probe (what is installed on disk) and the behavioral tally (what actually ran), which is the `scaffold-usage` query you already ran in step 2. A file-only probe misses a scaffold that was invoked but isn't checked in; a tally-only read misses one installed but never used. Join them.

**File probe — multi-label, not a single winner.** Probe each scaffold INDEPENDENTLY and name EVERY match; report `none` only when no probe matched (the old single-winner if-ladder hid co-installed scaffolds):

- **spacedock** — a spacedock WORKFLOW is on disk: any of a `.spacedock-state/` dir (a workflow's state checkout — including the split-root `docs/**/.spacedock-state`), a workflow README carrying spacedock frontmatter (a `README.md` with `commissioned-by:` or `stages:` keys), or a `_mods/` dir (commissioned workflow scaffolding). The DB tally EXCLUDES the `spacedock` family (see Behavioral tally), so a genuine incumbent is recognized HERE, by the on-disk workflow — never by the survey's own `spacedock:survey` self-call. Probe it with (run from the repo root):
  ```bash
  # spacedock incumbent: a workflow on disk (state checkout, workflow README frontmatter, or _mods/).
  # A `spacedock:survey` self-call leaves NO such file, so a survey-only repo prints nothing.
  spacedock_incumbent() {  # echoes "spacedock" iff a workflow is on disk
    if find . -type d \( -name '.spacedock-state' -o -name '_mods' \) -print -quit 2>/dev/null | grep -q .; then
      echo spacedock; return
    fi
    if find . -type f -name 'README.md' -print0 2>/dev/null \
         | xargs -0 grep -lE '^(commissioned-by|stages):' 2>/dev/null | grep -q .; then
      echo spacedock; return
    fi
  }
  ```
- **superpowers** — `.claude/skills/superpowers` exists, `superpowers` appears in `.claude-plugin/`, or a superpowers discipline skill dir is present (`.claude/skills/{brainstorming,writing-plans,executing-plans,subagent-driven-development,…}`);
- **gsd / get-shit-done** — a `.claude/skills/gsd` or `.claude/skills/get-shit-done` dir, a `.claude/commands/gsd` dir, or a `GSD.md` / `gsd.md` / `.gsd` file;
- **similar / unknown** — any other `.claude/skills` or `.claude/commands` tree (name the dirs you found);
- **none** — none of the above is present on disk.

**Behavioral tally.** The `scaffold-usage` rows are a `family → invocations` tally normalized from `tool_calls.skill_name` (`superpowers:brainstorming` and the bare `running-research-spikes` both fold to family `superpowers`); the `spacedock` family is excluded from the TALLY because survey/ensign self-invocation otherwise dominates and would make every surveyed repo read as "uses spacedock" (the false positive). That exclusion is the tally's; it does NOT hide a genuine spacedock incumbent — the spacedock FILE-PROBE above is what distinguishes a real on-disk workflow from the survey's own self-call. So a repo whose ONLY spacedock signal is the `spacedock:survey` self-call (no workflow on disk) is NOT named spacedock; a repo with a workflow on disk IS.

**Join and state the fact.** For each family appearing in either signal, state two observed facts plainly: its invocation count (the behavioral tally) and whether it is checked in on disk (the file probe). Do not narrate HOW the family was discovered (behavior vs files) — state only the usage + on-disk fact. For a family invoked but absent from disk, state the count and that it is not checked in. For example:

> `superpowers: 186 invocations (not checked in). Other one-offs: plan-writing, using-git-worktrees, systematic-debug, simplify, debugging.`

A family on disk and invoked is stated plainly (family + count + present on disk); a family on disk but never invoked is stated as installed-but-not-yet-invoked. The state-the-fact statement drives the comparative benefit in the report (step 4). The probe reads files; the *numbers* come from the scan (step 2).

## 4. Report and offer

Every `{slot}` below is a FILL slot: substitute the real value from the step-2 scan before you show the user. A literal `{slot}` (or a `<…>` angle token) left in what you present is a bug — never show the user an unfilled slot. If a slot has no data (e.g. zero OPEN decisions), drop that line rather than printing an empty slot.

**Cross-check the OPEN frontier against the repo (before you present it).** The `decision-open` query is a TRANSCRIPT-only scan — a fork that read OPEN there may already be shipped (a merged PR / a commit) and over-reports. For each transcript-OPEN fork, cross-reference the repo (`git log`, merged PRs via `gh pr list --state merged` if available, the working tree) and split it:

- **shipped** → **DROP** from the frontier. Evidence: a merged PR or a git-log commit whose subject/body CONFIDENTLY references the fork (its decision header or branch — an exact-ish token match).
- **decided-not-shipped** → move to a **BACKLOG** line (decided, no artifact yet).
- **never-decided** → **true open**, stays on the `THREADS TO PULL` frontier.

**Conservative-match rule.** DROP only on a CONFIDENT repo match; anything less than confident → KEEP on the frontier. A false "still open" is a cheap nudge; a false "shipped" silently hides a real open fork — so the asymmetry favors keeping.

**Mandatory degrade.** When NO repo signal is available (not a git repo, or `git log` / PR lookup fails or is empty), the frontier degrades to transcript-only and EVERY OPEN fork is flagged **`unverified`** in the report — never silently presented as authoritative. The degrade is the default behavior, not an error.

Render the report DIRECTLY in the same turn — do NOT stop and ask first, and do NOT precede it with any scratch-reasoning preamble. The first line the user sees MUST be the `SpaceDock survey —` title; never emit `I have everything I need`, `Let me cross-check …`, `Let me …`, or any "here is my plan" narration before the report. (The cross-check FINDINGS still appear — as report content in `THREADS TO PULL` / `BACKLOG` — but the scratch framing that produced them does not.) The survey is read-only orientation: the body IS the value, and a pre-body confirm/menu is a round-trip with no decision behind it. The ONLY stop in this flow is the end-of-report commission OFFER (the real decision).

The report is **value & numbers first**: a plain "what this gives you" lede + concrete figures from the user's own data lead; the mode/track vocabulary demotes to a labeled detail section below the fold. Every figure is a FILL slot from the step-2 scan (a literal `{slot}` shown is a bug); every figure derives from the surveyed session rows, never templated prose. Emit the report:

```
SpaceDock survey — your last {N} days                          ← {N} from scoping.span
(recent-window snapshot · agent logs only{if blank_cwd>0: · {blank_cwd} sessions had no working dir, not placed})

WHAT THIS GIVES YOU
  {plain language, no jargon: "You steer your agents by hand roughly {interruptions} times in
   this window. About {the manual/repeated share} are the same few moves repeated. A SpaceDock
   workflow can run those repetitive parts for you, stopping only where you'd want a say."}

BY THE NUMBERS
  {interruptions}  hand-steering interruptions                  ← decisions + veto markers (the {V} total)
  {hanging}  hanging threads (started, never closed)            ← count of THREADS TO PULL (post-cross-check OPEN)
  {no_followup}  decisions you made with no follow-up action    ← decision-no-followup (distinct from BACKLOG)
  {sessions}  sessions read{if codex-scoped>0: (Claude {claude} · Codex {codex_scoped} by working dir{if codex-presence>codex-scoped: ; name-match would say {codex_presence} — sibling repos, ignored})}
  {if sessions_that_orchestrated>0: {sessions_that_orchestrated}  sessions dispatched subagents ({subagents_dispatched} dispatched — their work isn't shown here)}

HOW YOU WORK
  {the inferred loop as an arrow chain} — {one honest line naming the dominant mode in PLAIN terms:
   "Mostly manual, repetitive tracks (not trivial — they take real work)." for manual;
   "Mostly exploratory — you steer an iterating agent." for exploration;
   for knowledge-work, NAME THE SPECIFIC TYPE(S) synthesized from the workstream names + areas
   (+ the user's checked-in skills) — never the bare archetype label alone. e.g.
   "People 1-1s & team assessment, positioning, interviews, and strategy mapping — a
   knowledge-work loop run mostly autonomously." Name the actual kinds the workstreams show
   (each track's WORKSTREAMS-cluster name is the source); "knowledge-work" appears only as the
   trailing class, never standing alone}

  ↓ full analysis: modes, work-by-area, what this can't see

═══ everything below is the demoted detail section ═══

THREADS TO PULL   (where you are now + what's still open — only if any fork is OPEN after the repo cross-check)
  {lead with the steady-state — where you are now and the few unresolved threads — and prompt
   the next move ("have you thought about …?"), NOT a narration of past decision history.}
  {if any OPEN exploration/knowledge-work track: held threads (tracked, prioritized — work you're holding, not bottlenecks):}
  ◐ {the open EXPLORATION / knowledge-work forks — deliberately-held threads}
  {if any OPEN manual track: manual backlog (gate-and-drive candidates):}
  ⚠ {the open MANUAL forks — never-decided questions}{if degraded: each flagged unverified (no repo signal)}

BACKLOG   (decided-not-shipped — only if any fork was decided with no artifact yet)
  {decided forks with no shipped artifact yet}                  ← NOT the no-follow-up figure (that's its own query)

RECENT DECISIONS  (answered or shipped)
  {the rest: header — short question}

WORKSTREAMS                                              mode
  {cluster the decisions + prompts into tracks; one line each, status glyph + the mode-classification label (manual / exploration / knowledge-work / unlabeled) per track. For a knowledge-work track, do NOT print the bare label — qualify it with the track's specific type drawn from its workstream name + areas, e.g. "knowledge-work · 1-1s & team assessment" or "knowledge-work · strategy mapping". The label class stays `knowledge-work`; the qualifier names the type.}

WORK BY AREA   (logical areas; worktree edits attributed to their area — F)
  {the product work-by-area buckets (kind=product), by edit count: area — {edits}}
  {if any kind=config: (+ {sum} edits in .claude/.beads/.git config + <external> sibling refs, footnoted)}
  {if the inferred workflow is branch-and-merge (worktree → PR → merge) AND a config bucket out-edits product: caveat — edits counted are the directly-edited branch; product code that lands via merged PRs is under-counted here, so this does NOT mean scaffolding > product}

CODEX   (only if codex-scoped>0; workdir-attributed, distinct from the name-only presence flag)
  {codex_scoped_sessions} Codex sessions attributed to this repo by exec_command working dir
  {if codex-presence>codex-scoped: (codex-presence matches {codex_sessions} by project NAME only — may include a same-named sibling; the workdir-attributed count above is sibling-free)}
  {if the codex-workstreams clusters are ALL (unlabeled): {codex_scoped} Codex sessions, unclassified (ad-hoc shell-driven) — a single honest line, NOT a (unlabeled)-only breakdown}
  {else: workstreams: {the codex-workstreams clusters — workstream → session count; (unlabeled) last}}
  activity: {the codex-activity tally — exec_command {n}, update_plan {n}, spawn_agent {n}}

SCAFFOLD
  {state-the-fact per family: family + invocation count + on-disk presence — e.g. "superpowers: 186 invocations (not checked in). Other one-offs: …"; or "none"}

INTERRUPTIONS
  {if any exploration/knowledge-work track: exploration/knowledge-work tracks: {n} steers across {m} sessions — this IS the work; book-keeping tracks the threads}
  {if any manual track: manual tracks: {n} steps across {m} sessions — gates + autonomy would carry these between your calls}

WHAT THIS CAN'T SEE   (the agent-log corpus is a partial lens)
  · work done outside agent sessions (manual edits, other tools, off-log discussion)
  · the project's history before this {N}-day window
  {if blank_cwd>0: · {blank_cwd} sessions with no recorded working dir}
  {if codex-presence>codex-scoped: · {codex-presence − codex-scoped} Codex sessions matched by name only (possible same-named sibling)}
  {if sessions_that_orchestrated>0: · work inside dispatched subagents (orchestration counted, content not — see "{subagents_dispatched} dispatched")}
```

Each `·` / conditional line drops when its slot is empty (the drop-empty-slot rule). The `WHAT THIS CAN'T SEE` block is the SINGLE consolidated home for the lens caveats — do not also scatter `uncaptured-cwd` or Codex-sibling asides elsewhere in the body.

### The discovery → commission bridge (close every report with this)

After the report, offer spacedock — leading with plain value, then keyed to the MODE of each track (from the `mode-classification` query). The modes call for DIFFERENT things; do NOT make one undifferentiated automate-everything pitch. As in the report above, every `{slot}` is a FILL slot: substitute the real step-2 numbers/forks/track-names before you show the user; a literal `{slot}` in your output is a bug.

**For the MANUAL tracks (mode=manual) — offer GATE-AND-DRIVE.** These are disciplined routine execution (the issue→worktree→PR loop, routine implementation): gate the crucial decisions and let the agent drive the loop between gates. Keep the gate-and-autonomy pitch — it is CORRECT for these. State it tied to the scan (the manual tracks' names + their gate-pass count or the interruption count). The per-scaffold flavor sharpens the offer:
  - **superpowers** maps its disciplines (brainstorming → writing-plans → executing-plans → subagent-driven-development) to stages with the interruption points made EXPLICIT gates.
  - **gsd / get-shit-done** maps its fixed phases to stages + durable entity state so several work items move concurrently, pausing only at gates.
  > For your MANUAL tracks (**{the manual track names}**): a spacedock workflow that runs the repetitive loop for you and stops only at the calls you'd want to make. These passed **{the gate-pass count}** times, so the agent can carry them between your gates.

**For the EXPLORATION tracks (mode=exploration) — lead with ITERATE/STEER, then BOOK-KEEPING; never automation.** These are human-driven creative/exploratory work (writing/content, design exploration, steering an agent that drifts): the involvement IS the point. The offer must speak to the iterate/explore loop FIRST — the agent iterates and you steer; an approval gate is ONE shape that loop can take, not the headline — do NOT lead with "explicit approval gates" for these. Then offer spacedock as structure for the parallel threads: track each draft/path and its state (in-flight / paused-by-choice / abandoned) so several run in parallel without losing which is which. An open thread is tracked-prioritized work, NOT a bottleneck; a cancelled path is a valid tracked outcome, NOT a failure. The exploration offer MUST NOT contain "advances on its own", "without you re-driving each", "minimize involvement", or any automate-the-human-out framing.
  > For your EXPLORATION tracks (**{the exploration track names}**): you steer while the agent iterates. Spacedock acts as book-keeping for the parallel threads, tracking each draft/design path and its state (in-flight / paused-by-choice / abandoned) so you run several at once without losing which is which. The **{the cancelled-path count}** cancelled paths are tracked outcomes, not failures. The involvement is the point; there is no automation here, only structure for the threads.

**For the KNOWLEDGE-WORK tracks (mode=knowledge-work) — offer a PER-THREAD TRACKER THAT AUTO-RUNS THE PROCESSING THEY ALREADY BUILT.** These are a recurring intake→process→commentary loop where the processing is already skill-shaped (an index/update skill, a commentary/summary skill the user wrote). Do NOT pitch this as passive book-keeping and do NOT disclaim automation — that is backwards for this class. Offer long-lived PER-ENTITY threads — one per recurring counterpart, one per advisory topic, plus one-off threads for interviews / strategy passes — and when a new entry lands on a thread, the workflow AUTO-RUNS the user's existing processing/commentary skills (their index/update skill becomes ONE automated step, not the whole job), stopping only at the judgment calls (intake/filing and what the commentary should emphasize). Keep the `knowledge-work` label, but the OFFER is for the recurring loop whose processing is already automatable; a knowledge-work track with no skill-shaped processing degrades to the generic book-keeping offer.
  > For your KNOWLEDGE-WORK tracks (**{the knowledge-work track names}**): a spacedock workflow with a long-lived thread per entity — one per recurring counterpart, one per advisory topic, plus one-offs for interviews and strategy passes. When a new entry lands, it auto-runs the processing you already built (your index/update and commentary skills become automated steps), and stops only at the judgment calls — what to intake and file, and what the commentary should emphasize.

**For UNLABELED tracks (mode=unlabeled) — generic book-keeping**, never a guessed automation pitch (the asymmetry favors not mis-offering: a missed automation offer is cheap; a wrong automation pitch at creative work is the misread to avoid).

If a project carries MULTIPLE modes, make EACH mode's offer (they MUST differ — the manual one keeps the gate-and-drive pitch; the exploration one leads with iterate/steer and carries none of the automate-the-human-out framing; the knowledge-work one names the batch loop). If it carries only one mode, make only that offer. **none** scaffold → the generic spacedock benefit, mode-keyed the same way. Each offer must cite a real scan number (filled track names, gate-pass count, OPEN forks, or cancelled-path count), not a placeholder.

Then make the offer, leading with the plain value:

> Want me to commission a spacedock workflow from this — so the repetitive work runs for you and stops only at the calls you'd want to make{if both manual and exploration modes: (gate-and-drive for the manual tracks; thread book-keeping for the exploration tracks)}?

On a **yes**, invoke commission in batch mode, supplying inputs derived from the scan (commission already accepts batch design inputs in its first message — see its Batch Mode). Assemble:

- **stages** ← the inferred workflow loop (for a detected scaffold, the comparable-workflow stage mapping above is the proposed stage list);
- **seed entities** ← the workstreams;
- **approval gates** ← the true-open forks (each OPEN decision that survived the repo cross-check is a candidate gate);
- **mission / entity** ← inferred from the workstreams and the project.

Hand off by invoking `commission` with those assembled inputs. Survey does NOT generate workflow files itself — file generation stays commission's job; survey only assembles the invocation and hands off.

On a **no**, stop — the survey stands on its own as an orientation.

## Synthesis guidance

- **Project name** = path basename.
- **Value & numbers first.** The report leads with the plain `WHAT THIS GIVES YOU` lede + the `BY THE NUMBERS` block, then demotes the mode/track vocabulary to the below-the-fold detail. Lead with what the user GETS, not with how the classifier labels their work. No scratch-reasoning preamble precedes the title.
- **Workflow + workstreams: infer them**, primarily from the decisions (the `PROMPTS` are sparse/noisy — secondary). Be honest when a track is one-off or stalled.
- **Decisions + stats are data, not invention.** `OPEN` = still needs the human; lead `THREADS TO PULL` with the true-open forks framed as the steady-state ("where you are now + what's still open"), plus a proactive prompt — NOT a narration of past decision history. The transcript scan can't tell shipped from open — that's what the step-4 repo cross-check is for; drop a fork to "shipped" only on a confident match, and flag the whole frontier `unverified` when there's no repo signal.
- **Work-by-area is identity, not a to-do list.** The `work-by-area` buckets say WHAT this project is (where edits land), by LOGICAL area regardless of physical location — a worktree edit is attributed to its area (a worktree `src/` edit is `src`), so worktree-based product work is not hidden. The lead lists product areas (`kind=product`) by edit count; genuine config (`.claude`/`.beads`/`.git`) and an `<external>` sibling-repo path demote to a footnote (`kind=config`) — still counted, just not the project's identity. **Branch-aware caveat:** when the inferred workflow is branch-and-merge (worktree → PR → merge) and a config bucket out-edits product, the edit count over-weights directly-edited scaffolding on the working branch and under-counts product code that lands via merged PRs — so caveat the signal (the counted edits are the directly-edited branch; product is under-counted here) and do NOT conclude "scaffolding > product." Report it separately from the decision frontier (where you stop).
- **Work modes, mode-keyed offers.** `mode-classification` labels each track manual / exploration / knowledge-work / unlabeled. Manual tracks (the issue→worktree→PR loop — repetitive but substantive, not trivial) get the gate-and-drive offer; exploration tracks (creative/content/design steering) get the iterate/steer + book-keeping offer (track the parallel threads + their states) — the involvement IS the point, so NO automation pitch for them and do NOT lead with gates; knowledge-work tracks (the recurring intake→process→commentary loop, where the processing is already skill-shaped) get a per-entity tracker that AUTO-RUNS the user's existing processing/commentary skills and stops only at the judgment calls — NOT passive book-keeping and NOT an automation disclaimer; an unlabeled track gets generic book-keeping, never a guessed automation pitch. The word the report renders is `manual`, not `mechanical` (reserve "mechanical" for genuinely trivial edits).
- **Name the SPECIFIC TYPE of knowledge work.** `knowledge-work` is the underlying archetype CLASS (the `mode-classification` value); it must NEVER stand alone in the user-facing report. Synthesize the specific type(s) — from the workstream cluster names + their work-by-area + the user's checked-in skills — and use that characterization in `HOW YOU WORK` (e.g. "People 1-1s & team assessment, positioning, interviews, and strategy mapping — a knowledge-work loop run mostly autonomously") and in the WORKSTREAMS mode column (e.g. `knowledge-work · 1-1s & team assessment`). The offer already names specifics (1-1s, interviews, strategy passes); extend that same specificity to the report BODY. Do NOT add a sub-type classifier — the type is synthesized from the workstream/area signal the survey already has.
- **Fill every slot, never invent.** Every `{slot}` in the report and the comparison comes from the step-2 numbers; a literal `{slot}` shown to the user is a bug. If a section's signal is empty (no OPEN decisions, no interruptions, no edits), say the run found none — never dress an empty section up as "no decisions."
- **Claude body + a Codex body section.** The Claude body (workflow, workstreams, decisions, work-by-area, scaffold) is built from Claude history. Codex is surfaced too, as its own section (step 4) from the `codex-scoped` set: the workdir-attributed count + the workstream clusters + the activity tally (Gemini and per-file Codex work-by-area remain deferred follow-ups). A repo whose ONLY history is Codex still reports "no agent history" at the `scoping=0` stop (the Claude body has nothing); the Codex section renders alongside a non-empty Claude body, not in place of it.
