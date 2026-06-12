# Command reference

The `spacedock` binary has ten subcommands in three groups, plus a top-level
`--version`:

| Command | Group | What it does |
|---------|-------|--------------|
| [`spacedock claude`](#launch-claude-codex-pi) | Launch | Start Claude Code with the first officer loaded |
| [`spacedock codex`](#launch-claude-codex-pi) | Launch | Start Codex (experimental) with the first officer |
| [`spacedock pi`](#launch-claude-codex-pi) | Launch | Start Pi (experimental) with the first officer |
| [`spacedock install`](#setup-install-doctor) | Setup | Install the per-host plugin, then run the compatibility check |
| [`spacedock doctor`](#setup-install-doctor) | Setup | Run the compatibility check alone |
| [`spacedock status`](#status) | Workflow | Read or mutate workflow state: the table, `--next`, `--where`, `--set`, … |
| [`spacedock new`](#new) | Workflow | Create an entity from stdin (alias for `status --new`) |
| [`spacedock state`](#state) | Workflow | Manage a split-root workflow's state checkout |
| [`spacedock completion`](#completion) | Workflow | Print a bash or zsh completion script |
| [`spacedock dispatch`](#dispatch) | Workflow | Build the worker dispatch artifacts the first officer hands an ensign |
| [`spacedock --version`](#-version) | — | Print the binary version and the contract level |

Run `spacedock` with no arguments to print the grouped help, and
`spacedock <command> --help` for a command's own flags. An unknown command or a
stray leading flag exits 2 with a diagnostic on stderr. The verbs are registered
in `internal/cli/cli.go`.

## --version

```bash
spacedock --version
```

Prints the binary version and the contract level, e.g. `spacedock 0.20.0 (contract 1)`.
The `(contract N)` token is load-bearing: the first-officer and ensign skills
read it to check the launcher and the installed plugin agree. The version string
defaults to the `dev` sentinel and is overwritten by the release pipeline's
linker stamp, so an unstamped `go build` reads as a dev build rather than
impersonating a release.

## Launch: claude, codex, pi

`spacedock claude`, `spacedock codex`, and `spacedock pi` each start the named
host with the Spacedock first officer loaded. Claude Code is the primary surface;
Codex and Pi are supported but experimental. The grammar is the same for all
three:

```bash
spacedock claude [task] [spacedock-flags] [-- host-flags]
```

- **The task comes first.** Positionals before `--` join with single spaces into
  the launch prompt handed to the first officer. With no task, a fixed bootstrap
  prompt starts the first officer rather than opening an idle agent.
- **Everything after `--` forwards verbatim to the host** (`claude` / `codex` /
  `pi`): `--model`, `--resume`, `--plugin-dir`, and the like. A task placed
  after `--` is host passthrough, not the launch prompt; the launcher warns on
  stderr when it detects this without altering the assembled argv.
- **The launch is contract-gated.** When no plugin is installed, the launcher
  auto-installs it then launches, so the single command yields a working session.
  `--no-install` opts out and prints the manual install remedy. A contract
  mismatch fails fast (exit 1), since auto-installing would not fix it.

Spacedock-owned launch flags, declared in `internal/cli/frontdoor.go`:

- `--safehouse`: force the safehouse sandbox wrap even without a `.safehouse`
  profile in the working directory. A `.safehouse` profile triggers the wrap
  automatically.
- `--safehouse-enable KEY[,KEY]`, `--safehouse-add-dirs DIR`,
  `--safehouse-add-dirs-ro DIR`: repeatable sandbox knobs whose presence also
  implies sandbox-on.
- `--plugin-dir DIR`: load a local plugin checkout, relaxing the contract gate.
  Repeatable, and accepted both before `--` (parsed by spacedock) and after `--`
  (forwarded verbatim).
- `--skip-contract-check`: bypass the contract gate and launch without resolving
  the installed plugin.
- `--no-install`: refuse to auto-install a missing plugin and print instructions
  instead.

The contract gate is bypassed by `--skip-contract-check` or by any `--plugin-dir`
(the local checkout supersedes the installed plugin). `spacedock pi` loads Pi's
native skills and extension rather than a plugin manifest; its only spacedock
flag is `--plugin-dir`.

## Setup: install, doctor

`spacedock install` installs the per-host plugin, then runs the compatibility
check. `spacedock doctor` runs the check alone.

```bash
spacedock install [--host claude|codex|pi] [--check]
spacedock doctor  [--host claude|codex|pi]
```

- `--host` defaults to `claude`. For `codex`, install prints the `codex plugin`
  commands to run from your shell rather than running them programmatically.
- `install --check` runs the compatibility report without installing.
- `doctor --plugin-manifest PATH` reads a manifest directly instead of resolving
  the installed plugin.

When `doctor` reports the installed plugin is out of date, refresh it with
`spacedock install --host claude`. See [Install Spacedock](../get-started/install.md)
for the full setup path.

## Workflow: status, new, state, completion, dispatch

These commands read and mutate workflow state. `status` and `dispatch` forward
their argv verbatim to their runners; the launcher neither parses nor reorders
their flags.

### status

```bash
spacedock status [args]
```

`spacedock status` resolves the workflow it acts on in this precedence: an
explicit `--workflow-dir DIR` (or the `PIPELINE_DIR` environment variable) is used
verbatim; otherwise it walks up from the working directory to the enclosing
commissioned workflow. With neither, it exits 1 with
`no Spacedock workflow here — pass --workflow-dir or run inside a workflow`. The
exit domain is `{0 success, 1 error}`, so a usage error is exit 1, never 2.

The query forms, parsed in `internal/status/native_runner.go`:

- **No flag**: prints the active-entity status table. `--archived` includes
  archived entities; `--quiet` and `--json` change the rendering.
- `--next` prints the items ready to dispatch (requires a stages block in the
  README).
- `--next-id` computes the next entity id. It accepts `--id-seed` and `--id-actor`
  (valid only with `--next-id` or `--new`).
- `--boot` prints the first-officer boot view (queue, stages, team state). It is
  incompatible with `--next`, `--next-id`, `--archived`, `--where`, and
  `--fields`/`--all-fields`.
- `--validate` validates the workflow. It prints `VALID` and exits 0, or prints the
  errors and exits 1.
- `--resolve REF` resolves a reference to one entity and prints its resolve line.
  With `--root ROOT` it resolves across every workflow under `ROOT`, accepting a
  `workflow::ref` qualifier.
- `--short-id REF` prints an entity's short display id.
- `--set SLUG field=value...` mutates an entity's frontmatter. Bare `completed`
  and `started` auto-fill a timestamp; every other field requires a value.
  Terminal transitions are gated (mod-block, merge-hook, verdict, and the
  require-external-proof guards); `--force` bypasses with a warning.
- `--archive SLUG` archives an entity.
- `--discover [--root ROOT]` prints the commissioned workflows under `ROOT`
  (default: the git toplevel of the working directory). It is incompatible with
  every other flag.
- `--where 'field = value'` filters the table. It supports `=`, `!=`, `field =`
  (empty), and `field !=` (non-empty), and is repeatable.
- `--fields a,b,c` / `--all-fields` chooses the columns. The two are mutually
  exclusive.

Most flags refuse to combine; the runner names the conflict and exits 1 (e.g.
`--set cannot be combined with --next`).

### new

```bash
echo "entity body" | spacedock new [--folder] SLUG
```

A pure alias for `status --new`: it prefixes the argv with `--new` and forwards
it, reading the entity body from stdin and auto-discovering the workflow.
`--folder` creates the entity as a folder rather than a single file.

### state

```bash
spacedock state init|new [--workflow-dir DIR]
```

Manages a split-root workflow's state checkout. `init` resumes a cloned workflow
by fetching the orphan state branch and checking it out as a linked worktree at
the workflow's `state:` path; a present checkout is a no-op that refreshes from
origin. `new` births that branch and worktree around a present split-root README.
An inline workflow has nothing to init. A missing or unknown subcommand exits 2.

### completion

```bash
spacedock completion bash|zsh
```

Prints a static shell-completion script for bash or zsh (exit 0). It completes
the top-level verbs and the common `status` flags. A missing or unknown shell
exits 2.

### dispatch

```bash
spacedock dispatch build | show-stage-def
```

Builds the worker dispatch artifacts the first officer hands an ensign. `build`
assembles the assignment (requires `--workflow-dir` unless `--print-schema` or
validate-only); `show-stage-def` prints a stage's definition (`--workflow-dir`
and `--stage`). A missing or unknown subcommand exits 2. See
[Adding a runtime](../contributing/adding-a-runtime.md) for how `dispatch build`
learns a new host mode.
