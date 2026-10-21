---
name: review
description: "Get a person's review of an artifact you produced or were handed: a plan, spec, report or other Markdown file. Use when you have a file that someone should read and react to before work continues, or when the user asks to send something for review. Picks the Terminal (TUI) or the Web from this project's recorded choice and never uploads without it. With no argument it shows or changes that choice; with text that is not a file it writes the artifact first."
---

# Review an artifact

Hand one artifact to a person for review. The project's recorded choice decides the
view and where the file goes. Never upload anything that the project has not agreed
to upload.

The two choices, worded as the view first, then where the file goes. Use these words
whenever you show or ask for the choice:

```text
Terminal (TUI): the file stays on this machine
Web: uploads to review.spacedock.md; anyone with the link can open it
```

The commands record Terminal (TUI) as `local` and Web as `remote`.

If `spacedock review-mode` is an unknown command, spacedock is too old for this
skill and has no `update` command. Tell the user to update it: `brew upgrade spacedock`
on Homebrew, or run the install line again. Stop.

## Pick the branch from the argument

- No argument: do **Show or change the choice** and nothing else.
- An argument that names an existing file (`test -f`): do **Review the file**.
- Any other argument is a request, not a path: do **Write first, then review**.

## Ask for the choice

Used whenever a choice is needed. Record only what the user answers, never your guess.

1. Ask "Terminal (TUI) or Web?" and show the two choices above. Say the answer is
   remembered for this project only.

2. Terminal (TUI): `spacedock review-mode set local`.

3. Web: do not record yet. Ask a second question on its own:

   ```text
   Web uploads the file to review.spacedock.md. Anyone with the link can open it, and it expires automatically after seven days. Agree?
   ```

   Yes: `spacedock review-mode set remote`. No: record nothing, say so, and offer the
   Terminal (TUI) instead; record `local` only if they accept it.

4. In your reply to the final answer, whichever it is (a Terminal (TUI) choice, a Yes,
   or a No while you still offer the Terminal (TUI)), tell the user they can change
   the setting at any time by running `/spacedock:review` with no argument.

## Show or change the choice

Publish nothing and open nothing.

1. Run `spacedock review-mode get --json`. If it exits 1, the recorded choice is
   unreadable: tell the user to run `spacedock review-mode clear`. Stop.

2. Tell the user this project's choice by quoting its line above word for word, or
   that none is recorded.
   Then offer exactly: keep it, switch to the other one, clear it (the next review
   asks again), or leave. With none recorded, run **Ask for the choice**.

3. Act only on the user's answer. Switching to the Terminal (TUI) runs
   `spacedock review-mode set local`; switching to the Web runs **Ask for the
   choice** from its second question, so it needs the explicit agreement.
   Clearing runs `spacedock review-mode clear`. Keep and leave run nothing. Then
   state what is recorded now and stop.

## Write first, then review

Write the artifact the argument asks for as a new Markdown file. Use the path the
argument names; otherwise create the file in a directory made with `mktemp -d`, so
the project gets no file nobody asked for. Never overwrite an existing file. Then
**Review the file** with it, and tell the user its path.

## Review the file

1. Read the choice:

   ```bash
   spacedock review-mode get --json
   ```

   `mode` is `local`, `remote`, or empty. If the command exits 1, the recorded
   choice is unreadable. Tell the user to run `spacedock review-mode clear`, then
   ask again. Stop.

2. If `mode` is empty, run **Ask for the choice** before anything else. If the user
   declines the Web and does not accept the Terminal (TUI), stop and publish nothing.

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