# Your first launch

One command starts your first session and orients you in a project you already have:

```bash
spacedock claude "/spacedock:survey"
```

This launches the [first officer](../concepts/operating-model.md#three-roles) (the orchestrator agent that runs a Spacedock workflow) inside Claude Code and hands it the survey task. The first launch also sets up the agents Claude Code loads; no separate setup step is needed.

Run it from inside a project that already has some agent history, such as a repo you have been coding in with Claude Code. Survey reads that history; an empty directory has nothing to report.

## What survey reads

Survey reads recorded agent session history through `agentsview`, a session-history tool. It does not parse raw logs by hand. It drives the `agentsview` binary to sync this project's sessions into a process-readable copy, then runs a fixed set of labeled, read-only SQL queries against that copy. The queries live in `skills/survey/references/queries.sql`, one labeled query per concern, so nothing is a black box.

Two behaviors matter when you run it:

- **It scopes to this repo by identity, not by name.** Survey resolves the repo root and scopes every query to that absolute path prefix. Because `agentsview` keys each session's project by the git-root basename, a same-basename sibling repo elsewhere on disk would otherwise fold in; the path-prefix scope keeps it out, and admits every checkout of this repo: the root, a subdir, a worktree.
- **If `agentsview` is missing, it asks before installing.** Survey needs `agentsview` to read the logs. When the binary is absent it tells you so and asks consent; on a yes it installs (`brew install --cask agentsview`, or the install-script fallback). It never installs without an explicit yes. If the sync fails for any reason (network, disk, permissions), it reports the exact failure and stops rather than guessing.

If the repo has no Claude agent history, survey says so plainly and stops. There is nothing to discover.

## What it reports

Survey leads with a one-line headline (the project, the session count, the date range, and the decision and interruption counts), then renders the body in the same turn. The body is the value, so the sections follow without a confirmation pause:

- **Inferred workflow.** The implicit loop reconstructed from the decisions and prompts, as an arrow chain, with one honest line about it.
- **Workstreams.** The decisions and prompts clustered into tracks, each tagged with its work mode (see below).
- **Work by area.** Where edits actually landed, by logical area (`src`, `internal`, `docs`, …) regardless of physical location. A worktree edit counts toward its area, so worktree-based work is not hidden. Genuine config paths (`.claude`, `.beads`, `.git`) and external sibling references demote to a footnote.
- **Needs you.** The open decisions, the forks raised but never resolved. **Survey leads the report with these**, because they are the work blocked on you. Exploration threads you are deliberately holding are separated from mechanical questions awaiting an answer.
- **Recent decisions** and **interruptions**: the answered or shipped forks, and how often you had to step in.
- **Scaffold.** If another agent scaffold is in use (superpowers, gsd / get-shit-done, or another `.claude` skill tree), survey states it as a fact: the family, its invocation count, and whether it is checked in on disk.
- **Codex** (only when present). Codex sessions land with no recorded working directory, so survey attributes them to this repo through each command's working directory and reports them as their own section: a session count, the workstream clusters, and an activity tally.

If a section's signal is empty, survey says the run found none of it. It never dresses an empty section up as "no decisions".

### The open frontier is cross-checked

The open-decision scan reads transcripts only, which cannot tell a shipped fork from a still-open one. Before presenting the frontier, survey cross-references the repo (`git log`, merged PRs, the working tree) and splits each open fork three ways:

- **shipped**: a confident match to a merged PR or commit. Dropped from the frontier.
- **decided, not shipped**: moved to a backlog line.
- **never decided**: stays on the `NEEDS YOU` frontier.

The match is conservative: a fork is dropped only on a confident repo match, because a false "still open" is a cheap nudge while a false "shipped" silently hides real open work. When no repo signal is available, whether because this is not a git repo or because the lookups fail, the frontier degrades to transcript-only and every open fork is flagged `unverified` rather than presented as authoritative.

## The commission offer

Survey ends with an offer, not an action. It asks whether you want it to [commission](../running-workflows/commission.md) a real Spacedock [workflow](../concepts/workflows-and-entities.md) built from what it found, turning the open decisions into [approval gates](../concepts/gates-and-decisions.md) and the workstreams into work items. Nothing has changed in your project until you say yes.

The offer is keyed to each track's work mode, so each track gets a distinct pitch. Survey classifies each track as **mechanical**, **exploration**, or **unlabeled**:

- **Mechanical tracks** (the routine issue → worktree → PR loop) get an **automation** offer: a workflow that gates the crucial decisions and lets the agent drive the loop between gates.
- **Exploration tracks** (creative, content, or design work where your steering is the point) get a **book-keeping** offer: structure for the parallel threads, tracking each draft or path and its state (in-flight / paused-by-choice / abandoned). There is no automate-the-human-out pitch here. The involvement is the work.
- **Unlabeled tracks** get the generic book-keeping offer, never a guessed automation pitch.

A project with both modes gets both offers. Each offer cites a real number from the scan: the track names, the gate-pass count, the open forks, the cancelled-path count.

On a **yes**, survey hands off to commission in batch mode, assembling the inputs from the scan:

- **stages** ← the inferred workflow loop;
- **seed entities** ← the workstreams;
- **approval gates** ← the open forks that survived the repo cross-check;
- **mission and entity** ← inferred from the workstreams and the project.

Survey does not write the workflow files itself; file generation stays commission's job. On a **no**, it stops; the survey stands on its own as orientation.

To define a workflow yourself instead, see [your first workflow](first-workflow.md). If `spacedock` is not yet installed, start with [installing Spacedock](install.md).

## The command grammar

Every launch uses the same shape:

```bash
spacedock claude "task" [--safehouse…] [-- host-flags…]
```

- **The task comes first.** It is handed to the first officer as the launch prompt. Here the task is `/spacedock:survey`, a skill the first officer runs. It could just as well be a plain sentence describing work.
- **`--safehouse` forces the launch through the sandbox.** A `.safehouse` profile in the working directory does the same automatically, so you only pass the flag when you want to force it.
- **Anything after `--` forwards verbatim to the host** (`claude` itself), including flags like `--resume`, `--model`, and `--plugin-dir`.

`spacedock codex "task"` and `spacedock pi "task"` take the same shape for the Codex and Pi harnesses. Claude Code is the primary surface; the examples here use it.
