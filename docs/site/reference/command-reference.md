# Command reference

The `spacedock` binary groups its subcommands into Launch, Setup, and Workflow, plus a top-level `spacedock --version` (the binary version and contract level). For the exact flags of any command, run `spacedock <command> --help`, the always-current source of truth; `spacedock` with no arguments prints the grouped help.

## --version

`spacedock --version` prints the version and contract level, the sandbox posture, then a per-runtime line reporting the installed spacedock plugin version:

```text
spacedock 0.20.1 (contract 1)
Sandbox: available, not enabled (no .safehouse profile)
claude: spacedock 0.20.1
codex: spacedock 0.20.0 (disabled)
pi: spacedock ready
```

The `Sandbox:` line is one of `enabled (safehouse)`, `available, not enabled (no .safehouse profile)`, or `unavailable (safehouse not on PATH)`. Each runtime line reads the plugin installed for that host: `spacedock <version>` when a plugin is installed (with ` (disabled)` appended only when the host reports it disabled), `spacedock ready` for pi (which launches from skills, not a versioned plugin), `spacedock not installed` when the host is present but carries no plugin, and `not installed` when the host binary itself is absent.

## Launch

`spacedock claude`, `spacedock codex`, and `spacedock pi` start a host with the first officer loaded. Claude Code is the primary surface; Codex and Pi are experimental. The grammar is the same for all three:

```bash
spacedock claude [task] [spacedock-flags] [-- host-flags]
```

The task comes first and becomes the launch prompt. Supported host arguments
after `--` forward verbatim (`--model`, `resume`, and the like). For Codex,
Spacedock keeps its normal launch banner, default approval posture, and bootstrap
prompt unless the forwarded token slice contains an exact token equal to
`resume`. It does not parse Codex option grammar: `spacedock codex -- --model
gpt-5.6-sol` remains a fresh first-officer launch, while `spacedock codex --
--model gpt-5.6-sol resume <id>` stays prompt-free. Use bare `spacedock codex` to
start the first officer. When no plugin is installed, the launcher auto-installs
it and launches, so the single command yields a working session; a contract
mismatch fails fast. The sandbox flags (`--safehouse` and its knobs) and the
contract-gate flags are listed by `spacedock claude --help`.

For Claude and Codex, a valid Spacedock plugin checkout beside the resolved
launcher is selected automatically. `--plugin-dir <checkout>` before `--` is an
explicit Spacedock override and takes precedence. Claude treats a post-`--`
`--plugin-dir` as an additional native session plugin, so the installed or
adjacent Spacedock provider is still gated normally. Codex has no native
session-local flag: a post-`--` `--plugin-dir` is rejected, and additional Codex
plugins must be installed persistently with `codex plugin add`.

An unsandboxed bootstrap launch carries no safehouse isolation, so per-action permission prompting is friction without a matching safety gain: `spacedock claude` starts in `--permission-mode auto` and Codex starts in `--ask-for-approval on-request` unless you supply an approval mode. A sandboxed bootstrap launch instead skips/bypasses approvals (`--dangerously-skip-permissions` for claude, `--dangerously-bypass-approvals-and-sandbox` for codex) since the sandbox is the gate. Claude suppresses its defaults when you pass your own mode or a resume. Codex suppresses its banner and bootstrap prompt only when its forwarded argv contains the exact `resume` token; an explicit approval mode prevents only a duplicate automatic approval flag.

## Setup

| Command | What it does |
|---------|--------------|
| `spacedock install` | Install the per-host plugin, then run the compatibility check |
| `spacedock doctor` | Run the compatibility check alone |

Both take `--host claude|codex|pi` (default `claude`). When `doctor` reports the plugin is out of date, refresh it with `spacedock install`. When the plugin is still contract-compatible but a newer one is available, `doctor` and the front-door launch print an opt-in upgrade hint (`run spacedock install --host <host> to refresh`); the hint never blocks the launch. See [Install Spacedock](../get-started/install.md) for the full setup path.

For source-checkout or retained development launchers, version compatibility is not a command-capability check. Immediately before every gate lifecycle, freshly resolve the launcher and run `spacedock gate --help` exactly once; do not cache the result. It must list `prepare`, `record`, `validate`, `eligibility`, `consume`, the prepare flags `--question`, `--artifact`, `--summary`, `--reference`, and `--workflow-dir`, and the semantic record flags. Missing capability halts before committing selected sources, preparing a room, or mutating state. Refresh the installed launcher or build the current checkout and select that executable with `SPACEDOCK_BIN`; do not hand-edit gate frontmatter as a fallback. Relative selected-source, Briefing, and room inputs resolve from the launch working directory.

## Workflow

The first officer runs these against workflow state as it moves entities; you operate through it, not by hand. They are documented here for completeness and for the rare direct use (scripting, debugging, restoring a state checkout on a fresh clone).

| Command | What it does |
|---------|--------------|
| `spacedock status` | Read or mutate the state: the entity table (omits the SOURCE column by default; `--fields source` or `--all-fields` restores it), `--next`, `--where`, `--set`, `--validate`, `--boot` (with `--identify`, the first officer's Startup identify — discovers the managed workflow(s), folds in the stage taxonomy and canonical ready-gate scheduling rows, and reports the boot sections; each row carries only `id`, `slug`, `current`, and `readiness`, while entity read and gate commands provide the complete decision record at engage; a gate stage alone is not ready; PR_STATE is a local `pr:` view, live PR state is checked at engage; local reads only, no mutation), `--read <ref-or-path>` (a file's structured frontmatter — including the nested `stages:` taxonomy, and projectable with `--fields` — plus a heading offset/lines map, for section-scoped reads; with `--checklist` / `--ac-scan` it extracts a stage report's checklist items with line ranges and per-AC evidence citations for the first officer's gate prep; `--stage` defaults to the entity's current `status` when omitted (so a bare `--read <entity> --checklist` reads the current stage's report), and `--stage X` reads a non-current stage) |
| `spacedock gate prepare <entity> --question TEXT --artifact REVIEW.md --summary TEXT [--reference FILE ...]` | Derive and bind a recorder-ready room for folder or flat form. Immediately after preparation the room contains exactly `gate-briefing.json` and `request.json`, with no copied sources or association. Selected files are exact local `git-root://<main\|state>/<full-commit>/<path>` objects with raw SHA-256 revisions; there is no fetch, ref requirement, or worktree fallback. Success prints `room`, `briefing`, `digest`, and `state`. |
| `spacedock gate record <entity> --briefing PATH` | Bind the exact canonical Briefing file; no basename is canonical. Request-less and advisory Briefings do not require or synthesize a summary. The recorder derives the current-stage logical gate and whether this opens, rebinds, or supersedes an attempt. |
| `spacedock gate record <entity> --room PATH` | Verify a prepared gate room and record its direct binding Result. The recorder checks the frozen request digest, derives fixed provider paths and the complete Artifact/Reference association, and leaves advisory output as retained evidence. |
| `spacedock gate record <entity> --decision approve\|revise\|hold --actor ID [--reason TEXT]` | Record a chat decision and its derived one-use application. Supported chat actor IDs are `person:captain` and `agent:first-officer`. Delegated First Officer decisions require an evidence reason; the recorder does not accept or authenticate Captain-message text. Recording never advances status or dispatches. |
| `spacedock gate record <entity> --round STAGE/CYCLE --briefing PATH/briefing.json --log PATH/briefing.review.jsonl [--feedback-cycle FILE]` | For a folder-form entity (`<slug>/index.md`), publish one complete advisory correction round to the immutable derived room `review/<stage>/round-<cycle>`, update the current `review-round` pointer, and optionally project its Feedback Cycles line. `STAGE` must exist in the workflow taxonomy but may differ from current status for historical backfill; cycle-line grammar and verdict binding are validated. Flat entities are refused because review artifacts accumulate beside the entity. Exact replay is a no-op; divergence is refused. |
| `spacedock gate validate <entity>` | Validate the canonical v1 logical-gate selection, ordered attempts, Briefing bindings, frozen closures, typed application, and Resolution shape; report the current recorded decision. Prototype fields and unknown binary-owned fields fail closed. |
| `spacedock gate validate <entity> --round STAGE/CYCLE` | Validate and read an advisory round through its entity pointer and immutable retained room. Every Resolution reports as advisory; no gate selection or workflow state changes. |
| `spacedock gate eligibility <entity>` | Read the current application's fail-closed condition without mutating it or querying dependency entities. |
| `spacedock gate consume <entity>` | Atomically advance status and spend an eligible pending approval once; stale approvals become superseded, while held or blocked approvals are refused. |
| `spacedock new` | Create an entity (`new [--folder] SLUG`) from a body on stdin |
| `spacedock dispatch` | Build the worker dispatch artifacts (`dispatch build`, `dispatch show-stage-def`) |
| `spacedock state` | Manage a [split-root workflow](../advanced/split-root-state.md)'s state checkout (`state init` resumes one on a fresh clone, `state new` births one) |
| `spacedock completion` | Print a bash or zsh completion script |

### `state commit`

`spacedock state commit <slug>` commits and synchronizes one active or clean
archived entity from a split-root state checkout. The operand is one canonical top-level entity slug,
not a nested path. A flat entity commits exactly `<slug>.md` plus its `<slug>/`
companion directory when present or tracked; a folder-form entity
commits every changed, new, or deleted non-ignored path below `<slug>/`, including
reports and artifacts. Dirty sibling entities and unrelated top-level state paths
remain untouched. Archived scope is publish-only: dirt or identity/shape collisions refuse before movement; clean unpublished history is pushed, and fully published state is a no-op.
Because archived recovery never creates or amends a commit, `-m MSG` applies only
to active entities and has no effect for an archived slug.

[Operate a workflow](../running-workflows/operating.md) covers how the first officer uses `status` on your behalf. Run `spacedock status --help` (and the same for each command) for the full flag list, the mutation guards, and the exit codes.
