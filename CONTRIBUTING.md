# Contributing

Thanks for considering a contribution. Spacedock is early, so we encourage you to share proposals and improvements as [GitHub issues](https://github.com/spacedock-dev/spacedock/issues) rather than opening pull requests directly. That lets us discuss the direction before anyone writes code.

## Develop

Spacedock is a Go module (`github.com/spacedock-dev/spacedock`, Go 1.22+). The
launcher binary lives in `cmd/spacedock/` (process entry point); command
routing, usage text, and exit-code behavior live in `internal/cli/`. The
plugin's skills live under `skills/`.

[`AGENTS.md`](AGENTS.md) is the source of truth for the development workflow —
read it before non-trivial work. The baseline gate every change must pass
before you call it done:

```bash
go test ./...          # baseline gate for every change
go test ./... -race    # catch data races
gofmt -w ./cmd ./internal
```

Add focused tests for a change before implementing it. The live runtime E2E
suites (real coding-agent hosts) are separate and gated behind a `live` build
tag — see [`docs/runtime-live-ci.md`](docs/runtime-live-ci.md).

## Build from Source

```bash
go build -o ./spacedock ./cmd/spacedock
```

This drops a `spacedock` binary at the repo root (gitignored, so it never
shows up in `git status`). To put a checkout-built binary on your `PATH`
instead, use `go install ./cmd/spacedock` — note it may shadow a
Homebrew-installed `spacedock`, so prefer the explicit `./spacedock` path
below when you want to be sure which binary you are running.

## Run your Branch

Run your freshly-built binary against the skills in your checkout, so both the
launcher and the skills come from your branch rather than an installed release:

```bash
go build -o ./spacedock ./cmd/spacedock
./spacedock claude --plugin-dir "$PWD" "/spacedock:survey"
```

Replace `claude` with `codex` or `pi` for the respective coding-agent hosts.

`--plugin-dir "$PWD"` loads the local plugin checkout directly and bypasses
installed-plugin resolution — it is the development path, not an install
substitute. It does **not** wrap the launch in the safehouse sandbox (see
[Sandboxing](docs/site/get-started/install.md#sandboxing)). This needs no
install step and no merge to `main`: it exercises the full current-checkout
stack — the launcher binary plus the local skills — straight from your working
tree.

If a launch misbehaves, run `spacedock doctor`.

### Avoid colliding with an installed Spacedock

If you also have Spacedock installed (e.g. via Homebrew), the build-and-run
path above stays isolated from it across all three collision surfaces:

- **Which binary runs** — invoke the local build by explicit path (`./spacedock`),
  not bare `spacedock` (which resolves the installed one on your `PATH`). Avoid
  `go install` for this, since whether its `$GOPATH/bin` copy wins over the
  installed binary depends on `PATH` order.
- **Which skills run** — `--plugin-dir "$PWD"` loads your checkout's skills and
  mutates no host plugin state, so the installed plugin is untouched.
- **Host install state** — a plain `go build` is stamped on the `next` (edge)
  channel, so a released install (stamped `main`) lives under a separate host
  plugin entry (`spacedock@spacedock-edge` vs `spacedock@spacedock`); the two
  coexist. `SPACEDOCK_DEV_BRANCH=main|next` overrides the channel at runtime.

The one combination that *does* overwrite a released install is running the
local binary with both `SPACEDOCK_DEV_BRANCH=main` **and** `spacedock install` —
that targets the same `spacedock@spacedock` entry. The default `next` stamp, and
the `--plugin-dir` path (which never installs), both avoid it.
