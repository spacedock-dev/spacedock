# Install Spacedock

This guide walks a fresh install end to end and names the output you should see
at each step. Every command here is one you can run and check against the stated
result.

Spacedock is two pieces that install separately:

1. **The `spacedock` launcher** — the command you run to start a session.
2. **The host plugin** — the first-officer and ensign agents, loaded by Claude
   Code or Codex.

The recommended setup installs the launcher with Homebrew and adds the plugin
with one command. A from-source build is available for development.

## Install with Homebrew (recommended)

1. **Install the launcher.**

   ```bash
   brew install spacedock-dev/homebrew-tap/spacedock
   ```

2. **Confirm it.**

   ```bash
   spacedock --version
   ```

   Prints the installed version, e.g. `spacedock 0.20.0`.

3. **Add the plugin to Claude Code.**

   ```bash
   spacedock install --host claude
   ```

   Adds the Spacedock plugin to Claude Code and runs a compatibility check that
   reports `OK` when the launcher and plugin match.

4. **Launch.**

   ```bash
   spacedock claude -- "your task"
   ```

   Starts the first officer in Claude Code with your task. When a `.safehouse`
   profile is present in the working directory the launch is wrapped through the
   sandbox.

## Use Codex instead

Codex support is available but experimental; Claude Code is the primary surface.

1. **Install the launcher** (same Homebrew step as above).

2. **Add the plugin.**

   ```bash
   spacedock install --host codex
   ```

   Codex installs plugins from your shell rather than programmatically, so this
   prints the two `codex plugin` commands to run. Run them, then use the
   first-officer skill in your Codex session.

3. **Launch.**

   ```bash
   spacedock codex -- "your task"
   ```

## Build from source (for development)

Use this when you're working on Spacedock itself. It builds the launcher from
the development branch and loads the plugin straight from your checkout, so your
local changes take effect immediately.

1. **Clone and build.**

   ```bash
   git clone --branch next https://github.com/spacedock-dev/spacedock
   cd spacedock
   go build -o spacedock ./cmd/spacedock
   ```

2. **Confirm the binary.**

   ```bash
   ./spacedock --version
   ```

   Prints `spacedock <version>` for your local build.

3. **Launch with the local plugin.**

   ```bash
   ./spacedock claude --plugin-dir "$PWD" -- "your task"
   ```

   `--plugin-dir` loads the first-officer and ensign agents from your checkout
   instead of the installed plugin, so edits to the repo are live.

The `next` branch is the development channel — it has no Homebrew release, so
there is no `brew install` for it. Use the Homebrew path above for a stable
install.

## Keep things in sync

`spacedock doctor` is the compatibility check. If it reports your installed
plugin is out of date, refresh it:

```bash
spacedock install --host claude
```

If the `spacedock` command itself is missing, install the launcher with Homebrew
first, then run `spacedock install --host claude`.

## Command grammar

The front door is `spacedock claude [host-flags…] [--safehouse…] -- "task"`
(and the same shape for `spacedock codex`):

- Flags before `--` pass through to the host (`claude` / `codex`), including
  `--plugin-dir`, `--resume`, `--model`, and the like.
- The bare text after `--` is the launch task, handed to the first officer.
- `--safehouse` forces the launch through the sandbox; a `.safehouse` profile in
  the working directory does the same automatically.
