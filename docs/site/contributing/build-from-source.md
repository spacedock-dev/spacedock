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

    Prints `spacedock <version>` for your local build.

3. **Launch with the local plugin.**

    ```bash
    ./spacedock claude "your task" -- --plugin-dir "$PWD"
    ```

    `--plugin-dir` is a host flag, so it rides after `--`. It loads the
    first-officer and ensign agents from your checkout instead of the installed
    plugin. Edits to the repo are live.

The `next` branch is the development channel. If you want the bleeding edge
without building from source, install the **edge channel** (`spacedock-edge`)
instead — it tracks `next` and ships through the same install path as stable.
For a stable build, use the [Homebrew install](../get-started/install.md).
