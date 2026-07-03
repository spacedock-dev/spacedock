---
title: "Command reference"
description: "A multi-agent orchestrator where nothing ships without a decision."
doc_version: "0.20.2"
last_updated: "2026-07-03 13:09:31"
---

# Command reference

The `spacedock` binary groups its subcommands into Launch, Setup, and Workflow, plus a top-level `spacedock --version` (the binary version and contract level). For the exact flags of any command, run `spacedock <command> --help`, the always-current source of truth; `spacedock` with no arguments prints the grouped help.

## --version

`spacedock --version` prints the version and contract level, the sandbox posture, then a per-runtime line reporting the installed spacedock plugin version:

```
spacedock 0.20.1 (contract 1)
Sandbox: available, not enabled (no .safehouse profile)
claude: spacedock 0.20.1
codex: spacedock 0.20.0 (disabled)
pi: spacedock ready
```

The `Sandbox:` line is one of `enabled (safehouse)`, `available, not enabled (no .safehouse profile)`, or `unavailable (safehouse not on PATH)`. Each runtime line reads the plugin installed for that host: `spacedock <version>` when a plugin is installed (with `(disabled)` appended only when the host reports it disabled), `spacedock ready` for pi (which launches from skills, not a versioned plugin), `spacedock not installed` when the host is present but carries no plugin, and `not installed` when the host binary itself is absent.

## Launch

`spacedock claude`, `spacedock codex`, and `spacedock pi` start a host with the first officer loaded. Claude Code is the primary surface; Codex and Pi are experimental. The grammar is the same for all three:

```
spacedock claude [task] [spacedock-flags] [-- host-flags]
```

The task comes first and becomes the launch prompt. Anything after `--` forwards verbatim to the host (`--model`, `--resume`, and the like). When no plugin is installed, the launcher auto-installs it and launches, so the single command yields a working session; a contract mismatch fails fast. The sandbox flags (`--safehouse` and its knobs) and the contract-gate flags are listed by `spacedock claude --help`.

An unsandboxed launch carries no safehouse isolation, so per-action permission prompting is friction without a matching safety gain: `spacedock claude` starts in `--permission-mode auto` and `spacedock codex` in `--ask-for-approval on-request`. A sandboxed launch instead skips/bypasses approvals (`--dangerously-skip-permissions` for claude, `--dangerously-bypass-approvals-and-sandbox` for codex) since the sandbox is the gate. Either posture is suppressed when you pass your own mode or a resume.

## Setup

| Command             | What it does                                                  |
| ------------------- | ------------------------------------------------------------- |
| `spacedock install` | Install the per-host plugin, then run the compatibility check |
| `spacedock doctor`  | Run the compatibility check alone                             |

Both take `--host claude|codex|pi` (default `claude`). When `doctor` reports the plugin is out of date, refresh it with `spacedock install`. When the plugin is still contract-compatible but a newer one is available, `doctor` and the front-door launch print an opt-in upgrade hint (`run spacedock install --host <host> to refresh`); the hint never blocks the launch. See [Install Spacedock](../../get-started/install/) for the full setup path.

## Workflow

The first officer runs these against workflow state as it moves entities; you operate through it, not by hand. They are documented here for completeness and for the rare direct use (scripting, debugging, restoring a state checkout on a fresh clone).

| Command                | What it does                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
| ---------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `spacedock status`     | Read or mutate the state: the entity table (omits the SOURCE column by default; `--fields source` or `--all-fields` restores it), `--next`, `--where`, `--set`, `--validate`, `--boot`, `--read <ref-or-path>` (a file's structured frontmatter — including the nested `stages:` taxonomy, and projectable with `--fields` — plus a heading offset/lines map, for section-scoped reads; with `--stage X --checklist` / `--stage X --ac-scan` it extracts a stage report's checklist items with line ranges and per-AC evidence citations for the first officer's gate prep) |
| `spacedock new`        | Create an entity (`new [--folder] SLUG`) from a body on stdin                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| `spacedock dispatch`   | Build the worker dispatch artifacts (`dispatch build`, `dispatch show-stage-def`)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| `spacedock state`      | Manage a [split-root workflow](../../advanced/split-root-state/)'s state checkout (`state init` resumes one on a fresh clone, `state new` births one)                                                                                                                                                                                                                                                                                                                                                                                                                       |
| `spacedock completion` | Print a bash or zsh completion script                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |

[Operate a workflow](../../running-workflows/operating/) covers how the first officer uses `status` on your behalf. Run `spacedock status --help` (and the same for each command) for the full flag list, the mutation guards, and the exit codes.

## Sitemap

- [Frontmatter contract](../frontmatter-contract/index.md)
- [Glossary](../glossary/index.md)
- [Supported sandboxes](../sandbox/index.md)
