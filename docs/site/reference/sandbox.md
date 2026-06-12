# Sandboxes

Spacedock ships no sandbox of its own; it wraps the launch through a sandbox
command when one is configured.

| Sandbox | Platforms | Trigger |
|---------|-----------|---------|
| [`safehouse`](https://agent-safehouse.dev/) | macOS | A `.safehouse` profile in the working directory, or the `--safehouse` flag |

A run is sandboxed only when the sandbox binary is on your `PATH`. When it is
absent, the launch prints an install hint and proceeds **unsandboxed**.
