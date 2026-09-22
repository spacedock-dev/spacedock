# Supported sandboxes

| Sandbox | Platforms | Trigger |
|---------|-----------|---------|
| [`safehouse`](https://agent-safehouse.dev/) | macOS | A `.safehouse` profile in the working directory, `--safehouse`, or any `--safehouse-*` option |

## Additional profiles

| Option | Type | Description | Default |
|--------|------|-------------|---------|
| `--safehouse-append-profile PATH` | Optional, repeatable path | Append a safehouse policy file for this launch. | None |

```bash
spacedock claude --safehouse-append-profile=file.sb
spacedock codex --safehouse-append-profile="profiles/local rules.sb"
spacedock pi --safehouse-append-profile=first.sb --safehouse-append-profile=second.sb
```

Both `--safehouse-append-profile=PATH` and `--safehouse-append-profile PATH` select safehouse.
Place the option before `--`; tokens after `--` go to the coding agent.
Relative paths start in the directory where you launch Spacedock.
Each occurrence supplies one path; repeats retain their order.
Safehouse loads project-config profiles before these profiles, then applies its final write protections.
Safehouse reports invalid profile paths or content, and Spacedock returns the failure.
An existing sandbox stays active; this option still requests a safehouse launch and cannot relax the parent sandbox.
