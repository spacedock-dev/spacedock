# Spacedock v1

Spacedock runs multi-step agent work through plain-text workflows. You define
the stages, Spacedock dispatches the right agent for each stage, and the work
record lives on disk so a long task can survive context limits and resume later.

This repository contains the Go launcher and compatibility bridge for the next
Spacedock command surface.

The first implementation target is conservative:

- provide a `spacedock` binary entry point;
- preserve current `status` behavior through a vendored compatibility path;
- prove per-workflow `.spacedock-state` state checkouts with the README symlink model;
- then replace the symlink dependency with native split-root status handling.

The development workflow for this repo lives in `docs/dev/README.md`. Runtime entities for that workflow live in `docs/dev/.spacedock-state/`, which is intended to be a separate git checkout or nested state repo.

## Install

Two lanes — see [`docs/install-journey.md`](docs/install-journey.md) for the
step-by-step journey with the observable output at each step.

**Stable lane (`main`)** — starting with `v0.20.0`, tagged releases, Homebrew
artifacts, and marketplace plugin installs come from `main`:

```bash
brew tap spacedock-dev/homebrew-tap
brew install spacedock
spacedock install --host claude
```

The no-tap one-liner `brew install spacedock-dev/homebrew-tap/spacedock` is
equivalent. In the stable lane, `spacedock install` resolves the released
marketplace plugin from `spacedock-dev/spacedock` on `main`, not from `next`.

**Dev-only lane (`next` + `--plugin-dir`)** — source build from `next`, the
primary development workflow. `next` has no Homebrew artifact, and `@next` is a
source-build or dev-publish path, not the stable install path:

```bash
git clone --branch next https://github.com/spacedock-dev/spacedock
cd spacedock
go build -o spacedock ./cmd/spacedock
./spacedock claude --plugin-dir "$PWD" -- "your task"
```

`--plugin-dir` loads the repo's own vendored `spacedock:first-officer` /
`spacedock:ensign` skills and relaxes the contract gate.

**Upgrade a stale plugin** — when `spacedock doctor` reports your installed
plugin is out of date (predates this binary's contract), reinstall it:

```bash
spacedock install --host claude
```

On the stable lane, `spacedock install` refreshes the released plugin from
`main`. The old-plugin/no-binary and binary/plugin-skew journeys still need
their release-gate confirmation before the `0.20.0` flip; this docs update does
not claim those upgrade paths are automatic yet.

[safehouse](https://agent-safehouse.dev) is a separate runtime dependency for
sandboxed launches — not installed by either lane. A `.safehouse` profile in the
working directory (or the `--safehouse` / `--safehouse-<key>=…` flags) wraps the
launch through it.

## Usage

```bash
spacedock claude [host-flags…] [--safehouse…] -- "task"   # launch claude --agent spacedock:first-officer
spacedock codex  [host-flags…] [--safehouse…] -- "task"   # launch codex with the spacedock:first-officer skill
spacedock --version                                       # spacedock <version> (contract 1)
spacedock doctor                                          # contract compatibility verdict
```

Flags before `--` pass through to the host; the bare text after `--` is the
launch task. `--skip-contract-check` bypasses the contract gate (bootstrap
only); a `--plugin-dir` launch relaxes it without the flag.

## Commands

```bash
go test ./...
go run ./cmd/spacedock --help
```
