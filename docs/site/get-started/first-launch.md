# Your first launch

One command orients you in a project you already have:

```bash
spacedock claude "/spacedock:survey"
```

Claude Code opens with Spacedock loaded and surveys the project: what your
agents have been doing, and the decisions still waiting on you.

## What survey reads

Survey reads your recorded agent sessions, read-only, scoped to this repo and
every checkout of it, and nothing else on disk. It reads them through
`agentsview`, a session-history tool; if the tool is missing, survey asks
before installing it. If the repo has no agent history, survey says so and
stops.

## What it reports

A one-line headline (project, sessions, date range, decision and interruption
counts), then:

- **Needs you.** The open decisions: forks raised but never resolved. These
  lead the report because they are the work blocked on you. Threads you are
  deliberately holding are separated from questions awaiting an answer.
- **Inferred workflow.** The loop you have been running without naming it, as
  an arrow chain, with one honest line about it.
- **Workstreams.** Your decisions and prompts clustered into tracks.
- **Work by area.** Where edits actually landed (`src`, `docs`, …);
  config-only paths drop to a footnote.
- **Recent decisions** and **interruptions**: the answered forks, and how
  often you had to step in.
- **Scaffold.** Another agent scaffold in use (superpowers, gsd, another
  `.claude` skill tree) is stated as a fact.
- **Codex** (when present). Codex sessions get their own section.

An empty section says so plainly.

Open decisions are cross-checked against the repo before they reach you: forks
that already shipped are dropped, decided-but-unshipped ones move to a backlog
line, and only the never-decided stay on the frontier. With no repo signal to
check against, open forks are flagged `unverified`.

## The commission offer

Survey ends with an offer, not an action. It asks whether you want it to
[commission](../running-workflows/commission.md) a real Spacedock
[workflow](../concepts/workflows-and-entities.md) built from what it found,
turning the open decisions into
[approval gates](../concepts/gates-and-decisions.md) and the workstreams into
work items. Nothing has changed in your project until you say yes.

The offer matches the work it saw, citing real numbers from the scan:

- **Routine loops** (issue → worktree → PR) get an automation offer: a
  workflow that gates the crucial decisions and lets the agent drive between
  gates.
- **Exploration** (creative or design work where your steering is the point)
  gets a book-keeping offer: structure for the parallel threads and their
  state. There is no automate-the-human-out pitch; the involvement is the work.

On a **yes**, survey hands what it found to commission: the inferred loop
becomes the proposed stages, the workstreams become the seed work items, and
the open forks become the gates. On a **no**, it stops; the survey stands on
its own as orientation.

To define a workflow yourself instead, see
[your first workflow](first-workflow.md).

## The command grammar

Every launch uses the same shape:

```bash
spacedock claude "task" [--safehouse…] [-- host-flags…]
```

- **The task comes first** and becomes the launch prompt. `/spacedock:survey`
  is a skill; a plain sentence describing work works just as well.
- **`--safehouse` forces the [sandbox](../reference/sandbox.md).** A
  `.safehouse` profile in the working directory does it automatically.
- **Anything after `--` goes to Claude Code itself**, including flags like
  `--resume` and `--model`.

`spacedock codex` and `spacedock pi` take the same shape.
