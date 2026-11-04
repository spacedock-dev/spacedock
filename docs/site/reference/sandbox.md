# Supported sandboxes

| Sandbox | Platforms | Trigger |
|---------|-----------|---------|
| [`safehouse`](https://agent-safehouse.dev/) | macOS | A `.safehouse` profile in the working directory, `--safehouse`, any `--safehouse-*` option, `SPACEDOCK_SAFEHOUSE=1`, or a nonempty `SPACEDOCK_SAFEHOUSE_ENABLE` |

## Per-session settings

These environment settings apply to Claude, Codex, and Pi in both launch modes.

| Setting | Value | Default |
|---|---|---|
| `SPACEDOCK_SAFEHOUSE` | `1` requests Safehouse. `0` adds no request. Other Boolean spellings follow Go's `strconv.ParseBool`. | No request |
| `SPACEDOCK_SAFEHOUSE_ENABLE` | Comma-separated capabilities, such as `ssh,docker`. A nonempty list also requests Safehouse. | No capabilities |

An explicit `--safehouse=true` or `--safehouse=false` replaces the Boolean environment value.
Any `--safehouse-enable` flags replace the environment capability list and retain their existing order.
A `.safehouse` profile or another sandbox flag still requests Safehouse independently.
A false Boolean value does not disable those triggers or an existing parent sandbox.
Invalid Boolean values fail before launch unless an explicit flag replaces them.

`SPACEDOCK_SAFEHOUSE_ENABLE=ssh spacedock pi --passthrough -- --mode rpc` requests the `ssh` capability for that session.

### Caller settings in provider mode

With `--passthrough`, Safehouse preserves explicitly supplied values for the selected host:

| Host | Canonical variables |
|---|---|
| Claude | `CLAUDE_CONFIG_DIR`, `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` |
| Codex | `CODEX_HOME` |
| Pi | `PI_CODING_AGENT_DIR` |

Presence selects the allowance, including an empty value.
Values pass unchanged, so an explicit Claude teams value of `0` stays `0`.
Absent variables receive no synthesized value.
Preserving these variables adds no filesystem grant or general environment forwarding.
Existing Safehouse filesystem rules still control access to the selected settings path.
Ordinary launch behavior remains unchanged.

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
