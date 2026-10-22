# Review reference

## `/spacedock:review`

In Claude Code, use `/spacedock:review` with the arguments below. In Codex, ask for the same action in plain words, for example "review plan.md with Spacedock Review".

| Argument | What it does | Default |
| --- | --- | --- |
| none | Shows this project's choice and lets you keep, switch or clear it, or leave. Opens nothing. | |
| `tui`, `terminal` or `local` | Records Terminal (TUI) for this project. Opens nothing. | |
| `web`, `webui` or `remote` | Asks you to agree to uploading; records Web only if you agree. Opens nothing. | |
| a file path | Reviews the file in the recorded place. | Asks Terminal or Web when nothing is recorded |
| any other text | Your agent carries it out as a prompt, then reviews the file that results. | |

```text
/spacedock:review plan.md
/spacedock:review web
/spacedock:review
```

## Terminal keys

| Key | What it does |
| --- | --- |
| <kbd>c</kbd> | Comment on the selected text |
| <kbd>s</kbd> | Suggest replacement text for the selection |
| <kbd>q</kbd> | Finish, or open the finish screen when your agent asked for a decision |
| <kbd>Tab</kbd> | On the finish screen, cycle the decision |
| <kbd>Enter</kbd> | On the finish screen, send; while typing a note, save it |
| <kbd>Esc</kbd> | On the finish screen, go back to the document; while typing a note, discard it |
| <kbd>?</kbd> | Show every key |

## Decisions

| You choose | What it means | What your agent does | Stored as |
| --- | --- | --- | --- |
| Approve | Accept this as it is. | Carries on with the task; gains no other permission | `approve` |
| Request changes | Not accepted. Changes are needed; your comments say which. | Edits the file, shares it again under a new link, waits again | `revise` |
| Hold decision | Not accepted, and no changes asked for. Stop here for now. | Stops and changes nothing | `hold` |

A review takes one decision. After it, the review is closed and later comments are not recorded.

## Commands your agent runs for web reviews

`rr` is short for `subspace-relay-client`; where `rr` is taken on your machine, the long name works.

| Command | What it does |
| --- | --- |
| `rr publish <file>` | Uploads the file and prints its share link. `--open-owner-link` opens your owner link in your browser. |
| `rr decision <share-url>` | Waits for the owner's decision. `--with-included` also returns the other reviewers' comments the decision includes; `--timeout` sets how long. |
| `rr feedback <share-url>` | Prints every comment and decision so far. Only the owner's machine can read them. |
| `rr revoke <share-url>` | Ends the link now. |
| `rr device` | Shows whether this machine is linked to an account. |
| `rr unbind` | Unlinks this machine from its account, which ends ownership for every browser that owned its reviews. |

```sh
rr publish plan.md --open-owner-link
rr decision https://review.spacedock.md/r/<share-id> --with-included --timeout 110m
rr feedback https://review.spacedock.md/r/<share-id>
rr device
```

These two end access:

```sh
rr revoke https://review.spacedock.md/r/<share-id>
rr unbind
```

## Limits

| Limit | Value |
| --- | --- |
| Share link | Seven days from sharing, never extended |
| Owner link | One use, 15 minutes |
| Shares per day, machine not linked to an account | 5 |
| Shares per day, machine linked to an account | 100 per machine, and 100 per account across its machines |
| File size | 1 MB per share |
| Feedback sends | 100 per reviewer session, 200 per review |
