# Review reference

<a id="skill-invocation"></a>

## Review requests

Use `/spacedock:review` in Claude Code or `$spacedock:review` in Codex. Both use the same arguments. Replace `plan.md` with your file name.

<!-- Keyword invocation verification: https://github.com/spacedock-dev/spacedock/issues/852 -->

| Arguments | What happens |
| --- | --- |
| `plan.md terminal decision` | Review in the terminal and make a decision |
| `plan.md web decision` | Review in the browser and make a decision, after agreeing to the upload |
| `plan.md decision` | Make a decision using your project's saved review location |
| `plan.md terminal questions` | Review in the terminal with questions enabled |
| `plan.md` | Review the file using your project's saved location |
| `terminal` | Save Terminal as this project's preference |
| `web` | Save Web as this project's preference, after agreeing to uploads |
| No arguments | View or change this project's preference |
| A prompt, such as `draft release notes for v2` | Ask the agent to carry it out, then review the resulting file |

A location supplied with a file applies to that review only. Without a saved preference, your agent asks where to review.

`terminal` also accepts `tui` or `local`. `web` also accepts `webui` or `remote`.

Or ask your agent:

> Send plan.md to me in Spacedock Review for a decision.

## Terminal keys

| Key | Action |
| --- | --- |
| <kbd>c</kbd> | Comment on selected text |
| <kbd>s</kbd> | Suggest replacement text |
| <kbd>Q</kbd> | Ask the agent a question when questions are enabled |
| <kbd>R</kbd> | Refresh the answer to your latest question |
| <kbd>T</kbd> | Set an automatic answer refresh interval; press again to stop |
| <kbd>q</kbd> | Finish, or open the decision screen |
| <kbd>Tab</kbd> | Choose a decision on the finish screen |
| <kbd>Enter</kbd> | Send the decision, or save the note you are editing |
| <kbd>Esc</kbd> | Leave the finish screen, or discard the note you are editing |
| <kbd>?</kbd> | Show all keys |

<a id="decisions"></a>

## Review modes

| Mode | Result |
| --- | --- |
| `feedback` | Comments |
| `binding` | Comments and the owner's decision |

A plain terminal file review defaults to `feedback`. The `decision` keyword requests `binding` mode.

Your agent receives `approve`, `revise` or `hold` for Approve, Request changes and Hold decision. See [Comments and decisions](decisions.md) for what each means.

<a id="commands-your-agent-runs-for-web-reviews"></a>

## Agent commands

Installing Spacedock adds `subspace-relay-client`; a Homebrew install also names it `rr`. Your agent runs these for you.

| Command | Action |
| --- | --- |
| `rr pair` | Pair this computer with your account, once; you sign in and approve in your browser |
| `rr publish <file> --open` | Upload the file, print its share link and open it in your browser |
| `rr decision <share-url> --with-included` | Wait for the owner's decision and included comments |
| `rr feedback <share-url>` | Read the review's comments and decisions from the publishing machine |
| `rr revoke <share-url>` | Revoke access through the share link |
| `rr device` | Check whether this computer is paired with an account |
| `rr unpair` | Unpair this computer. Share links keep working; you decide on its reviews again after pairing |

Use `--timeout` with `rr decision` to set the wait duration, for example `--timeout 110m`.

Your paired computers are listed at `review.spacedock.md/account/computers`, where you can remove one.

## Limits

| Limit | Value |
| --- | --- |
| Share link | Seven days from sharing by default |
| Automatic wait in Claude Code | Up to 110 minutes |
| Automatic wait in Codex | Up to 12 hours |
| Shares per day, unpaired computer | 5 |
| Shares per day, paired computer | 100 per computer and 100 across the account's computers |
| File size | 1 MB per share |
| Feedback sends | 100 per reviewer session, 200 per review |
