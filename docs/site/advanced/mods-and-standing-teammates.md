# Mods & standing teammates

Mods extend a workflow without touching the binary. A mod is a markdown file under `{workflow_dir}/_mods/`. There are two kinds, and one file can be both: a **lifecycle hook** that the first officer runs at a named point in the run, and a **standing teammate** declaration that spawns a long-lived specialist agent into the team. Both live in `_mods/*.md`. The difference is which sections the file carries and which binary reads them. `spacedock status` scans the `## Hook:` headings (the `--boot` MODS section, and the merge-hook guard); `spacedock dispatch` parses the standing-teammate sections.

## Lifecycle hooks

A mod hook is a `## Hook: {point}` section that the first officer runs at a fixed point in the run. Three points are supported:

- `startup`: runs once at boot, before the normal dispatch loop.
- `idle`: runs on the idle re-check pass when no entity is ready to dispatch.
- `merge`: runs at the terminal merge boundary for an entity, before any local merge, archival, or status advancement.

Hooks are additive and run alphabetically by mod filename. The body of a hook section is prose the first officer executes; it names the commands to run and the conditions to branch on, in plain markdown. Nothing compiles; the first officer reads the section and acts on it.

A mod can register more than one point. The shipped `pr-merge` mod (`docs/dev/_mods/pr-merge.md`) registers all three: its `## Hook: startup` and `## Hook: idle` sections scan for entities with a pending `pr` and advance any whose PR has merged, and its `## Hook: merge` section opens the code-branch PR, records `pr:` on the entity, and blocks until merge.

### Merge hooks can block, and the mechanism enforces it

A `merge` hook can wait for captain approval before pushing, or for a remote PR to merge. The first officer signals the wait through the entity `mod-block` field, and `spacedock status` enforces the discipline so a blocked entity cannot slip past the gate:

- **Set before invoking.** The first officer sets `mod-block=merge:{mod_name}` before running the merge hook:

  ```bash
  spacedock status --workflow-dir {workflow_dir} --set {slug} mod-block=merge:{mod_name}
  ```

- **Guarded.** `spacedock status --set` refuses any terminal transition while `mod-block` is non-empty. `--archive` refuses too. Pass `--force` to override.
- **Required when a merge hook exists.** Independently of `mod-block`, `status --set` and `status --archive` refuse to terminalize or archive an entity when the workflow registers any merge hook (`_mods/*.md` with a `## Hook: merge` section) *and* the entity's `pr` field is empty *and* `mod-block` is empty. This forces the merge ceremony to leave a truthful signal that a merge actually ran. `merge: local` in the workflow README exempts the `pr` requirement; `verdict=rejected` exempts it too (a rejected entity never runs the ceremony). `--force` bypasses everything.
- **Cleared in its own call.** When the blocking action completes, the first officer clears the block:

  ```bash
  spacedock status --workflow-dir {workflow_dir} --set {slug} mod-block=
  ```

  This clear MUST be standalone. `status --set` exits 1 if `mod-block=` is combined with a terminal field (`status={terminal}`, `completed`, `verdict`, or `worktree=`) in one call. Use two commits.

`mod-block` is read from frontmatter at boot, so a pending merge survives session resume. The first officer picks up which mod is blocking and resumes the wait.

## Standing teammates

A standing teammate is a long-lived specialist agent (a prose polisher, a code reviewer, a translator) declared by a mod with `standing: true` in its frontmatter. It lives in the team for the team's lifetime and is addressed by name. Use one when a recurring specialist judgment is worth a persistent agent rather than a fresh dispatch each time.

### Declaration

One mod file per teammate under `{workflow_dir}/_mods/{name}.md`. The parse contract (see `internal/dispatch/mods.go`):

- **Frontmatter** carries `standing: true` and an optional `description`.
- **`## Hook: startup`** declares the spawn config as `- key: value` bullets. `spacedock dispatch spawn-standing` reads `subagent_type`, `name`, and `model` here; `model` must be one of `sonnet`, `opus`, `haiku`. Backtick-wrapped values are unwrapped.
- **`## Routing Usage`** (optional) is the prose each ensign sees telling it when and how to route to the teammate.
- **`## Agent Prompt`** MUST be the last top-level section. Its body, from the line after the heading to end of file, is the verbatim prompt passed to the spawned agent. Any `## ` heading after it is rejected loudly by `spacedock dispatch spawn-standing`.

### Lifecycle

The first officer drives three `spacedock dispatch` subcommands, all reading `_mods/` directly. Do not grep frontmatter yourself:

```bash
spacedock dispatch list-standing --workflow-dir {wd}     # abs mod paths, one per line, sorted
spacedock dispatch spawn-standing --mod {abs_path} --team {team_name}
spacedock dispatch show-standing --workflow-dir {wd}     # ensign-facing routing block
```

- **Discovery** runs at boot via `list-standing`. It prints the absolute path of each `standing: true` mod, one per line, sorted alphabetically; empty output means none.
- **Spawn is deferred** to the first team-mode dispatch. `spawn-standing` emits an `Agent()` spec for the host to launch, or `{"status": "already-alive", "name": ...}` when the team config already lists that member. Standing teammates are **first-boot-wins**: when several workflows share one team, the first first officer to find the member absent spawns it, and the rest skip. A mod that fails to parse (missing `## Agent Prompt`, an invalid `model`, a trailing heading) is reported and skipped; it does not block the workflow.
- **Routing is best-effort and non-blocking.** Address the teammate by its declared `name`, with a 2-minute timeout. If no reply lands in time, the sender proceeds without the specialist's output. Round-trips of several minutes are normal on long drafts.
- **Teardown is team-scoped.** The teammate dies when the team is torn down (session end, explicit delete, captain shutdown). There is no cross-team or cross-session persistence; mid-session death is detected on the next routing attempt.

Bare (single-entity) mode and degraded mode still run discovery (it is cheap) but skip the spawn pass, because there is no team to spawn into.

### Ensign discovery

Ensigns find standing teammates without the first officer wiring anything per dispatch. When a workflow declares at least one standing teammate, `spacedock dispatch build` appends a `spacedock dispatch show-standing` fetch line to each ensign dispatch. `show-standing` renders a `### Standing teammates available in your team` block, carrying each teammate's `## Routing Usage` body when present and otherwise a one-line fallback, so every dispatched worker learns who to route to.

## The comm-officer prose-polisher

The canonical standing teammate is the **comm-officer**: a standing prose-polisher the first officer routes deliberate drafts through before captain review. By convention it is declared as `_mods/comm-officer.md` with `standing: true` and named `comm-officer`.

The first officer routes through it when composing **deliberate drafts**: PR bodies, gate-review summaries, long narrative entity-body sections, debrief content. It checks team membership first and treats the call as best-effort and non-blocking on the 2-minute timeout; if the comm-officer is absent or silent, the draft ships un-polished. Explicitly **out of scope**: live captain replies, short operational statuses (`pushed`, `tests green`, `PR opened`), tool-call output, commit messages, and transient logs. Polish is a deliberate-draft discipline, not a live-turn reflex.

The comm-officer's prose discipline is light-touch by default: it applies the `elements-of-style:writing-clearly-and-concisely` skill (Strunk) to cut empty words and tighten sentences while preserving the caller's voice, rhythm, and technical vocabulary. It defers to a project voice guide when one exists. For Spacedock's own docs, that guide is the [authoring directive](https://github.com/spacedock-dev/spacedock/blob/next/docs/site/AGENTS.md). The comm-officer and any doc contributor follow it, falling back to plain Strunk only where the guide is silent.
