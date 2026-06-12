# Commission a workflow

`/spacedock:commission` turns a description of the work you want tracked into a runnable workflow: a directory of markdown entities, a README that is the workflow's living spec, and a first officer ready to dispatch an ensign for each seed entity. You answer a short interactive design pass; the skill generates the files and launches a pilot run as the first officer.

Invoke it from a session started with `spacedock claude`. You can pass the mission inline:

```
/spacedock:commission product idea to simulated customer interview
```

Text after the command name is taken as the workflow mission and presented for confirmation rather than asked from scratch. With no argument, the skill greets you and asks.

## The four things you name

The design pass collects four decisions. The skill derives everything else (the directory path, the entity identity scheme, the rejection routing) from these and shows the full design for your confirmation before writing any file.

1. **The mission and what each entity is.** The first question asks what the workflow is for and what each work item represents. From the entity description the skill derives the entity label used throughout the generated files: "a design idea" becomes label `idea`, plural `ideas`, type `design_idea`. An entity is one work item, a markdown file that moves through the stages.

2. **The stages.** The skill detects the workflow's shape from your mission (shipping code, testing a hypothesis, or iterating on an artifact) and proposes a stage list for you to confirm, modify, add to, or trim. Stage names describe the bucket an entity is sitting in: activity-flavored (`implementation`, `review`, `validation`) or state-flavored (`proposed`, `published`, `accepted`). The skill pushes back on redundant names: `awaiting_validation` reads as "the entity is in awaiting_validation," so it suggests `validation` instead. `done` is the universal terminal and stays as-is.

3. **The gated stages.** A gate is a stage where the workflow pauses for your decision before an entity advances. By default the skill places one gate before the terminal stage. For each gate it also derives a rejection flow, the earlier stage an entity bounces back to when you reject, defaulting to the stage immediately before the gate. You confirm both in the design summary, stated in plain language ("If you reject at `review`, it goes back to `draft` for revision").

4. **The per-stage rules.** Each stage in the generated README carries three bullets that tell a dispatched ensign what "good" means for that stage: **`Outputs`** (what the worker produces), **`Good`** (your quality bar), and **`Bad`** (anti-patterns to avoid). The skill drafts these from the mission, but they are starting prose, not commitments. They are the rules every dispatched agent works from.

You also choose the entity ID style: `sd-b32` (recommended when multiple people or agents create entities across branches or worktrees), `sequential` (single-writer or numeric-continuity workflows), or `slug` (the filename is the identity). See the [frontmatter contract](../reference/frontmatter-contract.md) for what each style stores.

## What gets generated

After you accept the design, the skill writes everything into `docs/{mission-slug}/`:

- `README.md`: the mission, the schema, a per-stage section for every stage (with the `Outputs` / `Good` / `Bad` bullets), and a copy-paste entity template.
- One file per seed entity at `docs/{mission-slug}/{slug}.md`, each with valid YAML frontmatter and the description you gave.
- `_mods/pr-merge.md`, generated only when a stage modifies the repo (a worktree stage) and you accept the `pr-merge` mod, which tracks PR state on the entity's `pr` field instead of modeling merge as its own stage.

If any stage writes code or produces artifacts beyond the entity file, that stage is marked `worktree: true` so each entity gets an isolated branch and the main checkout stays clean.

## Tighten the README before the first dispatch

The README is the workflow's living spec. Before launching, the skill reminds you that the auto-generated per-stage bullets are best-guesses and prompts you to tighten them. You have two ways to do it:

- Open `docs/{mission-slug}/README.md` and edit the bullets under each `### {stage}` heading directly.
- Type `review stages` to have the skill walk you through each stage one at a time, flag the bullets that read as generic, and apply your amendments inline.

Editing here costs minutes; un-editing after agents have been dispatched against vague bullets costs more.

## What the first officer does next

Commission generates the files, then assumes the first-officer role itself to run the pilot; no separate launch step is needed.

1. It loads the [first officer's operating contract](../concepts/operating-model.md) and reads the workflow README you just generated.
2. It runs `spacedock status --boot` to read the workflow's current state.
3. It probes for the team-mode tools and dispatches an ensign for each seed entity that is ready to advance, processing them through the stages.
4. When the workflow goes idle or reaches a gate, it reports the pilot results: which entities moved, which stages they passed through, and any gate waiting on your decision.

To run the workflow in any later session, launch the first officer again:

```
spacedock claude
```

It reads the workflow state, picks up where the last session left off, and dispatches agents for any entity ready for its next stage. [Operating a workflow](operating.md) covers the day-to-day loop: seeing what is ready, dispatching, and handling gate decisions.
