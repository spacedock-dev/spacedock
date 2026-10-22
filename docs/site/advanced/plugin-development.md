# Local plugin development

Use a local checkout when you are developing Spacedock plugins. For normal setup, follow [Installation](../get-started/install.md).

Set `SPACEDOCK_MARKETPLACE_SOURCE` to install from a local or alternate
marketplace instead of the default `spacedock-dev/marketplace` — useful for
dogfooding a marketplace change before it reaches the production marketplace:

```bash
SPACEDOCK_MARKETPLACE_SOURCE=/path/to/local/marketplace spacedock install --host codex
```

When the resolved `spacedock` executable sits at the root of a checkout whose
host manifest names `spacedock` and points to an existing skills directory, the
launcher selects that adjacent checkout automatically. A source build at
`./spacedock` therefore needs no `--plugin-dir` for either Claude or Codex. A
release binary in a `bin/` directory, or a checkout with a missing or invalid
host manifest, keeps the normal installed-plugin resolution and compatibility
gate.

An explicit `--plugin-dir <checkout>` before `--` overrides adjacent discovery.
On Claude and Pi this is an ephemeral, install-free override — it bypasses
installed-plugin resolution for that one launch and does not wrap the launch in
the safehouse sandbox; use it for plugin development, not as an install
substitute. On Claude, a `--plugin-dir` after `--` remains a native additional
session plugin: it does not replace Spacedock or suppress the compatibility gate.

Codex has no such flag on its own CLI, so `spacedock codex --plugin-dir
<checkout>` and `spacedock install --host codex --plugin-dir <checkout>` build a
local marketplace from the checkout and install it under the binary's own
channel (`spacedock` stable / `spacedock-edge` edge — matching whatever
`spacedock codex` would otherwise install), then launch. This IS a persistent
install and Spacedock makes it exclusive across Codex channels: the selected
channel replaces any existing stable or edge Spacedock Codex plugin so
`$spacedock:*` skills resolve from the selected install. It is also a
point-in-time snapshot: editing the checkout afterward has no effect until the
command is re-run. The command prints an advisory that names the selected channel
and notes that the reported version reflects the checkout's checked-in manifest,
not necessarily its current HEAD. Codex does not accept a post-`--`
`--plugin-dir`; install additional Codex plugins persistently with
`codex plugin add`.
