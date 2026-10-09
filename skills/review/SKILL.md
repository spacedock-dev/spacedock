---
name: review
description: "Get a person's review of an artifact you produced or were handed: a plan, spec, report or other Markdown file. Use when you have a file that someone should read and react to before work continues, or when the user asks to send something for review. Picks local or remote review from this project's recorded choice and never uploads without it."
---

# Review an artifact

Hand one artifact to a person for review. How it is reviewed — locally on this
machine, or remotely through a share link — is this project's recorded choice.
Never upload anything the project has not agreed to upload.

1. Read the choice:

   ```bash
   spacedock review-mode get --json
   ```

   `mode` is `local`, `remote`, or empty.

2. If `mode` is empty, ask the user before anything else. Say what each option
   does: **local** opens the Subspace viewer on this machine and nothing leaves
   it; **remote** uploads the file to review.spacedock.md and gives a link
   anyone holding it can open. Say the answer is remembered for this project
   only. Record their answer, never your own guess:

   ```bash
   spacedock review-mode set local    # or: set remote
   ```

3. Local: invoke the `subspace:r` skill on the file.

4. Remote:

   ```bash
   spacedock remote-review <file> --relay-url https://review.spacedock.md
   ```

   Give the user the share link it prints. Exit 3 means the project has not
   chosen remote, and nothing was uploaded: go back to step 2.

If `spacedock review-mode` is an unknown command, spacedock is too old for this
skill; if `spacedock remote-review` reports `spacedock-remote-review is not
installed`, Subspace is missing. Either way, tell the user to reinstall spacedock,
which installs Subspace with it, and stop.

To change the choice later: `spacedock review-mode set local|remote`, or
`spacedock review-mode clear` to be asked again.
