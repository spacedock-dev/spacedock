---
name: review
description: "Get a person's review of an artifact you produced or were handed: a plan, spec, report or other Markdown file. Use when you have a file that someone should read and react to before work continues, or when the user asks to send something for review. Picks local or remote review from this project's recorded choice and never uploads without it."
---

# Review an artifact

Hand one artifact to a person for review. The project's recorded choice decides how:
locally on this machine, or remotely through a share link.
Never upload anything that the project has not agreed to upload.

1. Read the choice:

   ```bash
   spacedock review-mode get --json
   ```

   `mode` is `local`, `remote`, or empty. If the command exits 1, the recorded
   choice is unreadable. Tell the user to run `spacedock review-mode clear`, then
   ask again. Stop.

2. If `mode` is empty, ask the user before anything else. Say what each option
   does. **Local** opens the Subspace viewer on this machine, and nothing leaves it.
   **Remote** uploads the file to review.spacedock.md, and anyone who holds the link
   can open it. Say the answer is remembered for this project only.
   Record their answer, never your own guess:

   ```bash
   spacedock review-mode set local    # or: set remote
   ```

3. Ask the program that will do the work if it supports the hand-off. Run the line
   for the recorded mode.

   Local:

   ```bash
   spacedock-review --supports review-invocation-prepare-v1 >/dev/null 2>&1 || echo "UNSUPPORTED: spacedock-review $(spacedock-review --version 2>/dev/null || echo 'not installed')"
   ```

   Remote:

   ```bash
   subspace-relay-client --supports relay-publish-v1 >/dev/null 2>&1 || echo "UNSUPPORTED: subspace-relay-client $(subspace-relay-client --version 2>/dev/null || echo 'not installed')"
   ```

   If the output starts with `UNSUPPORTED`, stop and publish nothing. Name the
   program and its version from that line. Tell the user to run `spacedock update` in
   a terminal, then restart Claude Code or run `/reload-plugins`, then ask again.

4. Local: invoke the `subspace:r` skill on the file.

5. Remote: invoke the `subspace:remote-review` skill to publish the file. The
   user has already asked for remote review of this project's artifacts, so that
   recorded choice is the user's request to share it.

If the `subspace:r` or `subspace:remote-review` skill is not available, the Subspace
plugin is missing or too old. Give the same remedy as step 3, and stop.
If `spacedock review-mode` is an unknown command, spacedock is too old for this
skill and has no `update` command. Tell the user to update it: `brew upgrade spacedock`
on Homebrew, or run the install line again. Stop.

To change the choice later: `spacedock review-mode set local|remote`, or
`spacedock review-mode clear` to be asked again.
