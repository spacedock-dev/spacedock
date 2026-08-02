---
name: debrief
description: "This skill should be used when the user asks to \"debrief\", \"record what happened\", \"session summary\", \"write a debrief\", or wants to capture session activity (commits, task state changes, decisions, issues) into a structured record for the next session."
user-invocable: true
---

# Session Debrief

You are producing a session debrief — a structured record of what happened during a workflow session. The debrief captures commits, task state changes, gate decisions, issues found, and observations. It feeds forward into the next session so the first officer starts with context instead of cold.

Follow these phases in order. Do not skip or combine phases.

---

## Phase 1: Discovery

**Launcher command invariant:** Prefer `${SPACEDOCK_BIN:-spacedock}` for Spacedock helper calls so an explicit front-door launcher remains the command source; when unset, fall back to `spacedock` on `$PATH`. Bare `spacedock` examples are shorthand for that expression.

### Step 1 — Identify the workflow

The user may provide a workflow directory path as an argument. If they didn't, search for it:

1. Run `${SPACEDOCK_BIN:-spacedock} status --discover`. It prints one resolved workflow directory path per line on stdout, scanning from the enclosing git repository root and pruning agent-managed worktrees (`.claude/worktrees/`) and other build/vendor noise.
2. If exactly one path is returned, use it as `{dir}`. If multiple, ask which one. If none, report "No Spacedock workflow found."

Store the confirmed path as `{dir}`.

### Step 1b — Resolve the debrief home

The debrief home is not always `{dir}`. For a split-root workflow (its README declares a `state:` checkout) the established `_debriefs/` history and the auto-syncing convention live in the state checkout, not the definition dir on the code branch. Let the binary tell you which it is rather than re-parsing `state:` yourself:

1. Run `${SPACEDOCK_BIN:-spacedock} status --boot --json --workflow-dir {dir}` once and read three fields the binary already resolves: `state_backend` (`single-root` | `split-root`), `entity_dir` (the absolute resolved debrief home), and `entity_dir_present`.
2. Set `{debrief_root} = entity_dir`. Use `{debrief_root}/_debriefs/` everywhere `_debriefs/` appears below.

- `state_backend == single-root` → `entity_dir` is `{dir}`; debriefs read, write, and commit under `{debrief_root}/_debriefs/` on the current branch (unchanged behavior).
- `state_backend == split-root` → set `{state_checkout} = entity_dir` (`= {debrief_root}`). **Halt-gate:** if `entity_dir_present` is `false` the declared checkout is NOT materialized (an orphan state branch with no linked worktree — a fresh clone or a removed worktree). HALT, report "state not initialized," and stop. Do NOT fall back to writing into `{dir}` on the code branch. Tell the captain to run `spacedock state init` (manual fallback: `git fetch origin {state_branch} && git worktree add {state_checkout} {state_branch}`), then re-run the debrief. Otherwise the reads, the new-file write, and the commit all run in the state checkout on its own branch — never in `{dir}` on the code branch. Resolve the state branch from the checkout itself: `git -C {state_checkout} rev-parse --abbrev-ref HEAD`.

The debrief skill is user-invocable independent of a first-officer boot, so it cannot assume the checkout has converged — hence the explicit `entity_dir_present` gate before any `git -C {state_checkout}` operation.

### Step 2 — Determine session boundaries

Check for an existing debrief to anchor the session start:

1. Look for `{debrief_root}/_debriefs/*.md` files. Sort by filename (lexicographic = chronological).
2. If a prior debrief exists, read the most recent one's YAML frontmatter and extract the `last-commit` field. The current session starts from the next commit after that hash.
3. If no prior debrief exists, fall back: use `git log --oneline --reverse -- {dir} | head -1` to find the first-ever commit in the workflow directory, or use commits from the last 24 hours if the history is large.

Derive `{from_commit}` (the starting boundary). Run:

```bash
git log {from_commit}..HEAD --oneline -- {dir}
```

Present the session boundary to the captain:

> **Session boundary**
>
> Since: `{from_commit}` ({date or description})
> Commits in scope: {count}
>
> Adjust? (confirm or provide a different starting commit)

Wait for confirmation before proceeding.

### Step 3 — Read workflow metadata

Read `{dir}/README.md` frontmatter to extract:
- Entity label and plural form
- Stage names and ordering
- Which stages are gates, have `feedback-to`, use worktrees

---

## Phase 2: Extract Session Data

All extraction uses git and local files only — no external services.

### 2a. Commits

```bash
git log {from_commit}..HEAD --oneline -- {dir}
git log {from_commit}..HEAD --oneline
```

Run both — the first scoped to `{dir}` for entity-linked commits, the second unscoped to catch workflow-level commits (scaffolding edits, reverts, merge commits) that may live outside `{dir}`.

Split every commit into one of three buckets:

1. **PR squash-merge commits** — any commit whose message ends with `(#NN)` (GitHub's squash-merge pattern). These are the actual PR landings. They are **rolled up into the Shipped section as a PR link** — never enumerated individually.

2. **pr-merge mod landing commits** — commits with the pattern `merge: {slug} done (PASSED) via PR #NN`. These are Spacedock's post-merge state-transition commits emitted by the pr-merge mod after a PR lands. Use them in Phase 2b to resolve shipped entities to their PR numbers. Do NOT list them in the Non-PR section — the information is already captured by the Shipped section's PR link.

3. **Non-PR commits** — workflow-only commits worth surfacing that did not flow through a PR. Include commits with these prefixes:
   - Scaffolding edits: `docs:`, `debrief:`, `refit:`
   - Feedback cycles: `feedback:`
   - Ideation stage reports: `ideation:`
   - Mid-session scope updates: `update:` (entity scope changes, refactors to in-progress specs)
   - Direct merges on main that bypassed the pr-merge mod (e.g. `merge: {slug} done (PASSED) — direct merge from branch`)
   - Reverts: `Revert ...`

   **Suppress as routine state churn** (do NOT list):
   - `dispatch:`, `advance:`, `state:`, `track:` — normal stage-machine transitions; redundant with the Shipped section.
   - `seed:`, `file:` — captured in the Filed section (2f).
   - `merge: {slug} done (PASSED) via PR #NN` — captured in the Shipped section via the PR link.

   Collect the surfaced commits into a list. For each, record the short hash, a compact description (derived from the commit message), and — where useful — a one-clause context note (e.g., "feedback cycle 2 on 116", "reverted in `{hash}`").

   **Consolidate adjacent related commits** onto a single line when they share context. Examples:
   - `Feedback cycle commits on 116: \`88fe41b\`, \`bf93c00\`, \`5306fa4\`, \`31abc13\` (cycle 1, cycle 1 addendum, cycle 2, cycle 2 addendum).`
   - `Ideation stage-report commits on 121 and 125: \`d0a7365\`, \`e1af6cf\`.`

   If a commit was later reverted in the same session, mark it `**[reverted]**` and cross-reference the revert commit on the following line.

Record the first and last commit hashes of the full range for frontmatter.

Do NOT group all commits by entity slug — the old grouped-by-entity format is superseded.

### 2b. Task state changes (Shipped entities)

For each entity file in `{dir}/*.md` and `{dir}/_archive/*.md`, read the YAML frontmatter. Identify entities that reached a terminal/`done` state during this session by cross-referencing with commit messages containing `done:` / `merge: ... done` prefixes.

For each **shipped** entity, resolve the merged PR(s) by scanning the session commit log for the pr-merge squash pattern matching that entity's slug, or by looking at the `merge: {slug} done (PASSED) via PR #NN` commit. Extract the PR number(s).

Extract a one-sentence description from the entity body: read the entity file, locate the problem statement (the paragraph between the closing `---` of the frontmatter and the first `## ` heading). Use the first sentence as the description. Fallback to the entity title if the problem statement is empty.

For each shipped entity, record:
- Entity ID (numeric, e.g. `115`)
- Entity slug (backticked)
- PR number(s) with linked URL (see Phase 3 Step 4 for URL construction)
- One-sentence description

Emit one bullet per shipped entity in this format:

```
- **{id}** `{slug}` — [#{N}]({pr_url}). {one-sentence description}
```

If an entity shipped via multiple PRs (e.g. main + fixups), include all PR links inline:

```
- **{id}** `{slug}` — [#{N1}]({url1}) + [#{N2}]({url2}) fixups. {description}
```

### 2c. Gate decisions

Scan for gate-related activity:
- Entities that advanced past gated stages (approved gates)
- Entities that were rejected at gates (look for `feedback-to` cycles in commit messages or entity files with `### Feedback Cycles` sections)
- Extract reasoning from stage reports where available

### 2d. Issues found

Scan for issues:
- Commits with `fix:` prefix in the workflow directory
- Entity files whose `## Stage Report` sections contain `FAIL` items
- Any notable problems mentioned in stage reports

Categorize each issue into two buckets:

1. **Workflow-specific** — bugs or quirks in the user's pipeline (their stages, entities, agent instructions). These stay as local notes in the debrief.
2. **Spacedock issues** — bugs or limitations in the spacedock framework itself (first-officer template, commission skill, status script, ensign behavior, agent platform issues). These are candidates for GitHub issue filing.

### 2e. What's next

Extract the dispatchable list and overall workflow view via `spacedock status`:

```bash
spacedock status --workflow-dir {dir} --next
spacedock status --workflow-dir {dir}
```

If `spacedock status` is unreachable, raise the error to the captain rather than working around it — the first officer already invoked `spacedock status --boot` to start the session, so an unreachable launcher here indicates a real environmental problem worth surfacing.

Use the status output to populate:
- Entities blocked at gates (waiting for captain)
- Entities with non-empty `worktree` field (in-progress or orphaned)
- Overall workflow state

### 2f. Filed entities (new backlog seeds)

Scan the session commit log for new entity files created during this session. Look for commits with prefix `seed:` or `file:` in `{dir}`:

```bash
git log {from_commit}..HEAD --oneline -- {dir} | grep -E '^[a-f0-9]+ (seed|file):'
```

For each newly-filed entity, read the entity file and extract a one-sentence description from the problem statement (same method as 2b). Emit one bullet per filed entity:

```
- **{id}** `{slug}` — {one-sentence description}
```

If a filed entity also shipped in the same session (appears in both 2b and 2f), still list it in both sections — the Filed section documents what was added to the backlog, the Shipped section documents what landed.

---

## Phase 3: Draft and Review

### Step 1 — Collect the agent testimonial

Ask the driving agent:

> Before answering, self-identify as the agent driving this session:
> - Harness/runtime: the agent harness or runtime you are operating in (for example Claude Code, Codex, or Pi).
> - Model: the model name.
> - Model version/build: the exact version or build when runtime metadata exposes it.
>
> Write `unknown` for any identity field you cannot verify; do not infer or guess. Then, setting the project's subject matter aside: how would you describe the experience using Spacedock versus driving the same work without it? Be honest about friction, costs, or places where the workflow got in your way; this is not a request for praise.

Record the testimonial prose as `{agent_testimonial}` and its agent-supplied identity as `{harness_runtime}`, `{model}`, and `{model_version_build}`. Preserve `unknown` exactly where the agent cannot verify a value. Derive the session-scale counts `{tasks_touched}`, `{workers_dispatched}`, and `{prs_touched_or_merged}` (tasks touched, workers dispatched, PRs touched/merged) separately from the current session data; do not infer agent identity from those data or ask the agent to estimate session scale.

### Step 2 — Present the draft

Before drafting, construct PR URLs by reading the Spacedock plugin manifest (same logic as Step 4 below). Find the `skills/` directory containing this skill file, read its sibling `.codex-plugin/plugin.json` (the `.claude-plugin/plugin.json` mirror carries the same fields for the Claude host), and extract the `repository` field. The `repository` field may be a plain URL string (e.g. `https://github.com/spacedock-dev/spacedock`) or an object (e.g. `{"type": "git", "url": "..."}`) — handle both. Derive `{owner}` and `{repo}` for PR links of the form `https://github.com/{owner}/{repo}/pull/{N}`.

Assemble the extracted data into a draft debrief and present it to the captain:

> **Draft Debrief — {date} #{sequence}**
>
> ## Shipped
> {one bullet per shipped task from 2b, with PR link(s) and one-sentence description}
>
> ## Filed (backlog)
> {one bullet per new entity from 2f, with one-sentence description}
>
> ## Non-PR commits (workflow-only)
> {list of non-PR commits from 2a, each as "- `{hash}` {description} — {optional context}"}
>
> ## Decisions
> {placeholder — captain fills this in}
>
> ## Issues — Workflow
> {workflow-specific issues found, or "None identified."}
>
> ## Issues — Spacedock
> {spacedock framework issues found, or "None identified."}
>
> ## Observations
> {placeholder — captain fills this in}
>
> ## What's Next
> {dispatchable entities, gate-blocked entities, deferred items}
>
> ---
> Add your commentary to **Decisions** and **Observations**, or confirm as-is.

Note: there is no longer a grouped-by-entity "Commits" section. PR-associated commits are rolled up into the Shipped section via their PR link. Only workflow-only commits that did not flow through a PR appear in the "Non-PR commits" section.

### Step 3 — Captain commentary

Wait for the captain to:
- Add decisions (why gates were approved/rejected, scope changes, course corrections)
- Add observations (design insights, process improvements, things to remember)
- Edit any section
- Or confirm the draft as-is

Incorporate the captain's input into the final debrief.

### Step 4 — Handle spacedock issues

For each issue categorized as a spacedock issue, offer to file a GitHub issue:

1. Read the spacedock repo URL from the plugin manifest. Find the Spacedock plugin directory by locating the `skills/` folder that contains this skill file. Read its sibling `.codex-plugin/plugin.json` (or the `.claude-plugin/plugin.json` mirror) and extract the `repository` field.

2. For each spacedock issue, draft an **anonymized** GitHub issue. The issue body must NOT contain:
   - The user's actual mission or workflow purpose
   - Entity titles or content
   - Specific domain details

   Instead, include:
   - Clear description of the bug or limitation
   - Reproduction steps (using generic terms like "entity", "stage", "workflow")
   - Scale context: entity count, stage count, linear vs branching flow
   - Spacedock version (from `commissioned-by` in `{dir}/README.md` frontmatter)

3. Present each draft issue to the captain:

> **File as GitHub issue?**
>
> **Title:** {issue title}
> **Body:**
> {anonymized issue body}
>
> (y/n/edit)

4. If yes: run `gh issue create --repo {repo_url} --title "{title}" --body "{body}"` and record the issue URL.
5. If edit: let the captain modify title/body, then file.
6. If no: note "not filed" in the debrief.

If `gh` is not authenticated or fails, report the error and suggest the captain file manually.

---

## Phase 4: Write Debrief File

### Step 1 — Determine sequence number

Look at existing files in `{debrief_root}/_debriefs/` matching today's date pattern `{YYYY-MM-DD}-*-{$harness}-{$model}.md` (matching against the session's own harness/model suffix); the `*.md` match alone is acceptable when checking for older files that predate the suffix. The sequence number is one more than the highest existing sequence for today, or `1` if none exist. For split-root this counts the state-checkout debriefs only — never `{dir}/_debriefs/` — so numbering continues the established history and a same-named orphan left in the definition dir cannot perturb the sequence.

### Step 2 — Calculate duration

Derive approximate session duration from the timestamp of the first and last commits in scope:

```bash
git log {from_commit}..HEAD --format='%ai' -- {dir} | head -1
git log {from_commit}..HEAD --format='%ai' -- {dir} | tail -1
```

Calculate the difference as a human-readable duration (e.g., `~2h30m`).

### Step 3 — Write the file

```bash
mkdir -p {debrief_root}/_debriefs
```

Write the debrief to `{debrief_root}/_debriefs/{date}-{sequence:02d}-{$harness}-{$model}.md`. The filename carries the session's harness and model: `$harness` = the host runtime short name (`pi`, `claude`, `codex`); `$model` = the FO session's live model slug, hyphenated lowercase (e.g. `kimi-k3`, `gpt-5-6-sol`, `claude-opus-4-8`) — on pi, read it via `intercom({action:"list"})`; on claude/codex, from the session's own model identity:

```markdown
---
session-date: {YYYY-MM-DD}
sequence: {N}
first-commit: {hash}
last-commit: {hash}
duration: {approximate duration}
---

# Session Debrief — {YYYY-MM-DD} #{N}

{optional 1-3 sentence narrative framing of the session — headline theme, task counts, phase structure}

## Shipped
{one bullet per shipped task in the format:
"- **{id}** `{slug}` — [#{N}]({pr_url}). {one-sentence description}"
Multiple PRs inline with "+":
"- **{id}** `{slug}` — [#{N1}]({url1}) + [#{N2}]({url2}) fixups. {description}"}

## Filed (backlog)
{one bullet per new entity in the format:
"- **{id}** `{slug}` — {one-sentence description}"
If a filed entity also shipped in the same session, note that inline:
"- **{id}** `{slug}` — shipped same session."}

## Non-PR commits (workflow-only)
{brief lead-in line, e.g. "State transitions and scaffolding that don't belong to a PR:"}

{list of non-PR commits, each as:
"- `{hash}` {description} — {optional context}"
Reverts and reverted commits should be cross-referenced, e.g.:
"- `{hash}` **[reverted]** {description}. {context}"
"- `{hash}` Revert \"{original}\" — {context}"}

{optional trailing note: "All other session commits are rolled up in the shipped PRs above."}

## Decisions
{captain-contributed content, or "_(none recorded)_" if none provided}

## Issues — Workflow
{workflow-specific issues, or "None identified." if empty}

## Issues — Spacedock
{spacedock issues with filing status: "filed as owner/repo#N" or "not filed"}
{or "None identified." if empty}

## Observations
{captain-contributed content, or "_(none recorded)_" if none provided}

## Agent Testimonial
- Date: {YYYY-MM-DD}
- Harness/runtime: {harness_runtime}
- Model: {model}
- Model version/build: {model_version_build}
- Session scale: {tasks_touched} tasks touched; {workers_dispatched} workers dispatched; {prs_touched_or_merged} PRs touched/merged

{agent_testimonial}

## What's Next
{dispatchable entities, gate-blocked entities, deferred items — organized by
tier or category where useful, e.g. "Stalled from prior sessions", "Recommended
next session", "Other backlog", "Ideation"}
```

### Step 4 — Commit the debrief

**Single-root** — commit on the current branch:

```bash
git add {dir}/_debriefs/{date}-{sequence:02d}-{$harness}-{$model}.md
git commit -m "debrief: session {date} #{sequence} — {summary}"
```

**Split-root** — commit path-scoped in the state checkout on its own branch, then push so peers see it (the state checkout is one shared, non-branched index; a bare `git add -A` / `git commit` would sweep up a sibling writer's staged entity, so name the path):

```bash
git -C {state_checkout} add -- _debriefs/{date}-{sequence:02d}-{$harness}-{$model}.md
git -C {state_checkout} commit -m "debrief: session {date} #{sequence} — {summary}" -- _debriefs/{date}-{sequence:02d}-{$harness}-{$model}.md
git -C {state_checkout} push origin {state_branch}
```

On a non-fast-forward push rejection, integrate peers and re-push (the single-file commit replays atop the peer's; disjoint paths → no conflict):

```bash
git -C {state_checkout} pull --rebase origin {state_branch}
git -C {state_checkout} push origin {state_branch}
```

If the `pull --rebase` CONFLICTS, HALT: `git -C {state_checkout} rebase --abort`, report the conflicting path(s), and stop — never force-push or auto-resolve. If the checkout has no `origin` remote, commit locally and skip the push.

Where `{summary}` is a brief count like "8 tasks completed, 3 new filed".

Report the file path to the captain:

> Debrief written to `{debrief_root}/_debriefs/{date}-{sequence:02d}-{$harness}-{$model}.md` and committed.
