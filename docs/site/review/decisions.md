# Comments and decisions

Comments identify passages to discuss or change. A decision tells your agent whether to continue, revise or stop.

[Terminal reviews](terminal.md) can collect comments alone or request a decision. In a [browser review](browser.md), reviewers comment and the owner decides.

## Decisions

| You choose | What it means | What your agent does | Stored as |
| --- | --- | --- | --- |
| Approve | Accept this as it is. | Carries on with the task; gains no other permission | `approve` |
| Request changes | Not accepted. Changes are needed; your comments say which. | Edits the file and presents it for review again | `revise` |
| Hold decision | Not accepted, and no changes asked for. Stop here for now. | Stops and changes nothing | `hold` |

A review takes one decision. After it, the review is closed and later comments are not recorded.

For browser reviews, sharing a changed file creates a new link. Send that link to anyone reviewing the revision.
