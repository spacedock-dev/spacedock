Spacedock 0.28.1 makes the first install and the first workflow work without surprises.

- The installer and the install guide now use the stable channel, and the installer's closing check no longer reports Review and Compass as missing on a fresh install.
- `spacedock update` checks every host that holds a Spacedock piece. It fails with a clear message when a plugin or a Pi extension package is missing, and it says how to change a marketplace entry that Claude Code refuses.
- Installing on Claude Code adds the channel marketplace before it removes the old plugin, so a failed step no longer leaves you without the plugin.
- Interactive launches check once a day for a newer Spacedock through spacedock.md. Add `--no-update-check` or set `SPACEDOCK_NO_UPDATE_CHECK=1` to skip it.
- New workflows from `/spacedock:commission` give each starting task an Acceptance criteria section, so their first review gate passes.
- When an old local Subspace plugin hides remote review, the review skill now names the plugin and how to replace it, instead of blaming the current relay program.

Full changelog: https://github.com/spacedock-dev/spacedock/compare/v0.28.0...v0.28.1
