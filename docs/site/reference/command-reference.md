# Command reference

The `spacedock` binary groups its subcommands into Launch, Setup, and Workflow, plus a top-level `spacedock --version` (the binary version and the host OS/arch, and — inside an agent session — that session's runtime and sandbox state). For the exact flags of any command, run `spacedock <command> --help`, the always-current source of truth; `spacedock` with no arguments prints the grouped help.

## --version

`spacedock --version` reports the binary version and the host OS/arch, and — when it is running inside an agent session — that session's runtime and sandbox state. Outside any session it prints two lines:

    spacedock 0.26.0
    OS: darwin/arm64

Inside a session it also names the host OS/arch, the runtime it detected, the marker that proved it, which session this is, and whether this process is running inside a sandbox:

    spacedock 0.26.0
    OS: darwin/arm64
    Runtime: claude (CLAUDECODE, session afd74765)
    Sandbox: inside (agent-safehouse)
    contract 3

The session identifier is the first eight characters of the host's own session id — the same prefix Claude Code uses to name `~/.claude/teams/session-afd74765` — so you can tell two concurrent sessions apart and match one against its state on disk. Hosts that do not expose a session id, such as pi, omit it:

    Runtime: pi (PI_CODING_AGENT, PI_CODING_AGENT_DIR)

When markers for more than one runtime are set — a nested session can leak them — it reports the ambiguity rather than guessing, and still exits 0:

    Runtime: ambiguous (CODEX_THREAD_ID, CLAUDECODE) — pass --host

Being outside every runtime is a normal state, not a fault — it means a human at a terminal. There is no `Runtime:` line at all in that case, because there is no session to report: the output is the two lines shown above (the version line plus the `OS:` line).

The `Sandbox:` line answers one question — is this process sandboxed? Inside a sandbox it names it; otherwise it reports whether safehouse is available to sandbox future launches:

    Sandbox: inside (agent-safehouse)
    Sandbox: not sandboxed (safehouse available)
    Sandbox: not sandboxed (safehouse not installed)

`spacedock status --boot` reports the same three. The pre-launch banner answers the neighbouring but different question — whether the launch it is about to perform will be wrapped — so its `Sandbox:` line reads in terms of that launch.

The trailing `contract 3` is a frozen compatibility sentinel read only by skill versions predating the current version gate. It prints inside a session only. For what is installed for each host — plugin versions and enablement — use `spacedock doctor`.

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

The launcher and plugin are a version-gated bundle. During a gate lifecycle, the real `gate prepare` invocation is the capability check: a nonzero result halts before presentation or later state effects. Refresh the installed bundle or build the current checkout and select that executable with `SPACEDOCK_BIN`; do not scrape help output or hand-edit gate frontmatter as a fallback. Relative selected-source and room inputs resolve from the launch working directory.

## Workflow

The first officer runs these against workflow state as it moves entities; you operate through it, not by hand. They are documented here for completeness and for the rare direct use (scripting, debugging, restoring a state checkout on a fresh clone).

| Command | What it does |
|---------|--------------|
| `spacedock status` | Read or mutate the state: the entity table (omits the SOURCE column by default; `--fields source` or `--all-fields` restores it), `--next`, `--where`, `--set`, `--validate`, `--boot` (with `--identify`, the first officer's Startup identify — discovers the managed workflow(s), folds in the stage taxonomy and canonical ready-gate scheduling rows, and reports the boot sections; each row carries only `id`, `slug`, `current`, and `readiness`, while entity read and gate commands provide the complete decision record at engage; a gate stage alone is not ready; PR_STATE is a local `pr:` view, live PR state is checked at engage; local reads only, no mutation), `--read <ref-or-path>` (a file's structured frontmatter — including the nested `stages:` taxonomy, and projectable with `--fields` — plus a heading offset/lines map, for section-scoped reads; with `--checklist` / `--ac-scan` it extracts a stage report's checklist items with line ranges and per-AC evidence citations for the first officer's gate prep; `--stage` defaults to the entity's current `status` when omitted (so a bare `--read <entity> --checklist` reads the current stage's report), and `--stage X` reads a non-current stage) |
| `spacedock gate prepare <entity> --question TEXT --artifact REVIEW.md --summary TEXT [--reference FILE ...]` | Derive and bind a recorder-ready room for folder or flat form. Immediately after preparation the room contains exactly `gate-briefing.json` and `request.json`, with no copied sources or association. Selected files are exact local `git-root://<main\|state>/<full-commit>/<path>` objects with raw SHA-256 revisions; there is no fetch, ref requirement, or worktree fallback. Success prints `room`, `briefing`, `digest`, and `state`. |
| `spacedock gate record <entity> --room PATH` | Verify a prepared gate room and record its direct binding Result. The recorder checks the frozen request digest, derives the fixed room evidence paths and complete Artifact/Reference association, and leaves advisory output as retained evidence. The current workflow stage must be an actionable gate, and the bound Briefing must use the canonical v1 stage-qualified identity and name that stage; malformed or mismatched identity fails without mutation. |
| `spacedock gate record <entity> --decision approve\|revise\|hold --actor ID [--reason TEXT]` | Record a chat decision and its derived one-use application. Supported chat actor IDs are `person:captain` and `agent:first-officer`. Delegated First Officer decisions require an evidence reason; the recorder does not accept or authenticate Captain-message text. Recording never advances status or dispatches. The current workflow stage must be an actionable gate, and the bound Briefing must use the canonical v1 stage-qualified identity and name that stage; malformed or mismatched identity fails without mutation. |
| `spacedock gate record <entity> --round STAGE/CYCLE --briefing PATH/briefing.json --log PATH/briefing.review.jsonl [--feedback-cycle FILE]` | For a folder-form entity (`<slug>/index.md`), publish one complete advisory correction round to the immutable derived room `review/<stage>/round-<cycle>`, update the current `review-round` pointer, and optionally project its Feedback Cycles line. `STAGE` must exist in the workflow taxonomy but may differ from current status for historical backfill; cycle-line grammar and verdict binding are validated. Flat entities are refused because review artifacts accumulate beside the entity. Exact replay is a no-op; divergence is refused. |
| `spacedock gate validate <entity>` | Validate the canonical v1 logical-gate selection, ordered attempts, Briefing bindings, frozen closures, typed application, and Resolution shape; report the current recorded decision. Prototype fields and unknown binary-owned fields fail closed. |
| `spacedock gate validate <entity> --round STAGE/CYCLE` | Validate and read an advisory round through its entity pointer and immutable retained room. Every Resolution reports as advisory; no gate selection or workflow state changes. |
| `spacedock gate eligibility <entity>` | Read the current application's fail-closed condition without mutating it or querying dependency entities. |
| `spacedock gate consume <entity>` | Spend an eligible pending approval once and advance status atomically; stale approvals become superseded, while held or blocked approvals are refused. On an approval whose target stage is terminal, consume spends nothing and writes no status: it leaves the application `pending` and returns the route `approved-awaiting-merge` (idempotently, on repeat), and `merge guard` discovers/arms the delivery mechanism when it acts. |
| `spacedock merge guard <slug> --verdict passed\|rejected` | Run the terminal merge ceremony and, with delivery proven, finalize: the sole terminal consumer of a pending terminal-target approval — `application.state: consumed`, the terminal status, `verdict`, and `completed` move in one write, with recorded delivery state (`mod-block`/`pr`) retired in the same write. |
| `spacedock merge guard <slug> --rework` | Delivery requires rework: write the pending terminal-target approval `pending→superseded`, route the entity through the record stage's declared `feedback-to`, and clear `pr`/`mod-block`. Refuses without a pending terminal approval, or with a missing/undefined/terminal `feedback-to`. |
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
