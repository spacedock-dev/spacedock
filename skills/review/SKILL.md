---
name: review
description: "Get a person's review of an artifact you produced or were handed: a plan, spec, report or other Markdown file. Use when you have a file that someone should read and react to before work continues, or when the user asks to send something for review. Picks the Terminal (TUI) or the Web from this project's recorded choice and never uploads without it. With no argument it shows or changes that choice; tui or web records it; a file path is reviewed; any other text is a prompt to carry out first."
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

## Asking the user

Each question below is asked through the harness's question tool when it has one
(Claude Code: AskUserQuestion; if its schema is not loaded, load it with ToolSearch
first), as one menu with the options named. Where there is no such tool, ask in plain
text with the same options. Wait for the answer; record only what the user chose,
never your guess.

## Pick the branch from the argument

In a natural-language review request, "Spacedock Review" names the product, never
Web. Read the saved choice first: saved Terminal uses Terminal, saved Web uses Web,
and nothing saved runs **Ask for the choice**. Only explicit Web words (`web`,
`browser`, `upload`, `share link`, `review.spacedock.md`) select Web for one review.
Review the named existing file with the requested modifier, such as `decision`;
do not generate a replacement artifact. This does not change the argument branches
below or the four file review keyword forms.

Take the first branch that fits. Compare the reserved words as the whole argument,
ignoring case.

1. No argument: do **Show or change the choice** and nothing else.
2. `tui`, `terminal` or `local`: run `spacedock review-mode set local`, say that
   Terminal (TUI) is recorded for this project, and end with the sentence about
   changing it (**Ask for the choice**, step 4). Open and upload nothing.
3. `web`, `webui` or `remote`: do **Ask for the choice** from step 3. Open and upload
   nothing.
4. An argument that names an existing file (`test -f`), removing enclosing quotes
   if present: do **Review the file**. Test the whole path before parsing keywords.
5. An existing file followed by review keywords: do **File review keywords**.
6. Any other text is a prompt: do **Carry out the prompt, then review**.

## File review keywords

Recognize these four forms before treating text as a prompt. Read trailing keywords
from the right; the remainder is the existing file path. Keywords ignore case;
preserve the path exactly. Quoted paths such as `"plans/launch plan.md"` are one
locator: quote the path when testing it and handing it on, never evaluate it as code.

| Argument | Hand-off |
|---|---|
| `<file> terminal decision` | Terminal once, `subspace:r` with `--mode binding` |
| `<file> web decision` | Web once, `subspace:remote-review`, wait for the owner's decision |
| `<file> decision` | Saved choice or the existing ask-and-save flow; binding on Terminal, owner's decision on Web |
| `<file> terminal questions` | Terminal once, `subspace:r` with `--allow-question` |

Refuse `decision` together with `questions`, in either order: the installed
Subspace supports questions only for feedback, not binding decisions. Refuse Web
with `questions`, with or without `decision`: Web questions are unsupported.
Retain both requested intents in the explanation; do not drop one or substitute
another location. Refuse any other suffix only when the first word after the file
is a location or modifier word (`tui`, `terminal`, `local`, `web`, `webui`, `remote`,
`decision`, `questions`) and the whole request is not one of the four forms.
The four suffixes in the table are exact: a location alone after the file
(`terminal`, `tui`, `local`, `web`, `webui` or `remote`) is refused. Do not
reinterpret it as a bare-file review or infer an omitted modifier.
Give corrected requests only from the table, one per retained intent (so
`decision` with `questions` may offer both halves), each on its own line; use the
actual path, quoting it if needed, for example
`/spacedock:review plan.md terminal decision`. Explain the refused words without repeating the refused command.
Refusals open, upload and generate nothing. Every other text keeps the prompt
branch, including
`plan.md summarize this file into summary.md`.

For a valid form, **Review the file** with its location and modifier. An explicit
location applies only to this review: never set, clear, or temporarily change the
saved choice. An omitted location keeps the existing ask-and-save flow.

## Ask for the choice

Used whenever a choice is needed.

1. Ask "Terminal (TUI) or Web?" with two options, Terminal (TUI) and Web, each
   described by its line above. Say the answer is remembered for this project only.

2. Terminal (TUI): `spacedock review-mode set local`.

3. Web: do not record yet. Ask a second question on its own, with two options,
   "Yes, upload it" and "No":

   ```text
   Web uploads the file to review.spacedock.md. Anyone with the link can open it, and it expires automatically after seven days. Agree?
   ```

   Yes: `spacedock review-mode set remote`. No: record nothing, say so, and end
   the turn. Offer the Terminal (TUI) in a sentence of your reply, not as another
   question: the user records it by running `/spacedock:review tui`.

4. In your reply to the final answer, whichever it is (a Terminal (TUI) choice, a Yes,
   or a No while you still offer the Terminal (TUI)), tell the user they can change
   the setting at any time by running `/spacedock:review` with no argument.

## Show or change the choice

Publish nothing and open nothing.

1. Run `spacedock review-mode get --json`. If it exits 1, the recorded choice is
   unreadable: tell the user to run `spacedock review-mode clear`. Stop.

2. Tell the user this project's choice by quoting its line above word for word, or
   that none is recorded.
   Then ask, with four options, Keep, Switch, Clear or Leave: keep it, switch to the
   other one, clear it (the next review asks again), or leave. With none recorded,
   run **Ask for the choice**. This is asked on every entry, however often the
   project has chosen before.

3. Act only on the user's answer. Switching to the Terminal (TUI) runs
   `spacedock review-mode set local`; switching to the Web runs **Ask for the
   choice** from its second question, so it needs the explicit agreement.
   Clearing runs `spacedock review-mode clear`. Keep and leave run nothing. Then
   state what is recorded now and stop.

## Carry out the prompt, then review

Do what the prompt asks, for example writing a Markdown file. Use the path the prompt
names; otherwise create the file in a directory made with `mktemp -d`, so the project
gets no file nobody asked for. Never overwrite an existing file the prompt did not
ask you to change. Then **Review the file** with the resulting file, and tell the user its path. If the prompt produced no
file, say so and stop.

## Review the file

1. Read the choice:

   ```bash
   spacedock review-mode get --json
   ```

   `mode` is `local`, `remote`, or empty. If the command exits 1, the recorded
   choice is unreadable. Tell the user to run `spacedock review-mode clear`, then
   ask again. Stop.

2. Without an explicit location, if `mode` is empty, run **Ask for the choice**
   before anything else. If the user declines, stop and publish nothing.
   With an explicit location, use that location once regardless of the saved mode.
   For explicit Web, unless saved `mode` is already `remote`, ask the seven-day
   upload agreement from **Ask for the choice**, step 3, for this review only.
   Yes authorizes this upload without recording a choice; No stops without opening
   or uploading anything. Never run `review-mode set` or `clear` for a one-off
   location. An unreadable saved choice still stops at step 1.

3. Ask the program that will do the work if it supports the hand-off. Run the line
   for the selected location.

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
   a terminal, then restart Claude Code or run `/reload-plugins`, or in Codex start
   a new session, then ask again.

4. Local: invoke the `subspace:r` skill on the file, adding `--mode binding` for
   `decision` or `--allow-question` for `questions`. Let Subspace own its question
   receiver and capability checks; if it cannot keep a receiver active, stop with
   its remedy rather than opening a review without questions.

5. Remote: invoke the `subspace:remote-review` skill to publish the file. The
   recorded Web agreement or this review's explicit upload agreement authorizes
   sharing. For `decision`, include "wait for the owner's decision" in the handoff;
   the companion owns owner pairing and returning the verified decision and its
   selected comments. Keep its destination disclosure and wait behavior.

For decisions, follow the companion's trusted result: Hold pauses dependent work
and leaves the file unchanged. Comments, question answers, advisory feedback or
successful publication alone never approve dependent work.

If the `subspace:r` or `subspace:remote-review` skill is not available, the Subspace
plugin is missing or too old. Give the same remedy as step 3, and stop.