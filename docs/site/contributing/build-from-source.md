# Build from source

Use this when you're working on Spacedock itself. It builds the launcher from the
development branch and loads the plugin from your checkout, so local changes take
effect immediately. For a normal install, see [Install Spacedock](../get-started/install.md).

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

    Keep this build unstamped: do not add a Git tag, revision, or `git describe`
    ldflag. A source build reports the embedded checkout manifest version plus
    `+dev` (for example, `spacedock 0.26.0+dev`), so it stays compatible with the
    adjacent skills even when a future-minor release tag is the nearest Git
    ancestor.

3. **Launch with the adjacent local plugin.**

    ```bash
    ./spacedock claude "your task"
    ```

    The launcher resolves its own executable and validates the adjacent
    `.claude-plugin/plugin.json` before selecting this checkout. Codex does the
    same with `.codex-plugin/plugin.json` through its local marketplace adapter.
    Use `--plugin-dir /another/checkout` before `--` only when you need an
    explicit Spacedock override. On Claude, a plugin directory after `--` is an
    additional host plugin and does not replace Spacedock.

The `next` branch is the development channel. If you want the bleeding edge
without building from source, install the **edge channel** (`spacedock-edge`)
instead — it tracks `next` and ships through the same install path as stable.
For a stable build, use the [Homebrew install](../get-started/install.md).
