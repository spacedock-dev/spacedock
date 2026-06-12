# Your first workflow

Your first workflow comes from one of two places: [survey](survey.md) offers
to build one from what it found in your project, or you describe one yourself
to the `/spacedock:commission` skill. A workflow is your work moving through
stages you define, pausing at the gates where you decide.

## Commission a workflow

Describe the work in the launch line:

```bash
spacedock claude "/spacedock:commission ship features through design, implementation, and review"
```

Commission asks a few questions, one at a time: what each work item is, the
stages it moves through, which stages pause for your decision, and the quality
bar for each. You confirm or adjust each proposal, and it presents the design:

> I'll call you Captain — let me know if you prefer something else.
>
> For each run, we process tasks going through the following stages:
>
> a. backlog — new tasks wait here for triage
> b. design — the approach is worked out and written down
> c. implementation — the change is built on an isolated branch
> d. review — the result is checked against the design
> e. done — accepted and closed
>
> If you reject at review, it goes back to implementation for revision.
>
> Our pilot run will be with:
>
> - Add rate limiting to the API
> - Fix the flaky login test
>
> Entity identity will use `sequential`.
>
> All files will be created in `docs/ship-features/` for you to review.
>
> Accept this design, or tell me what to change.

**Nothing is generated until you accept.** Tell it what to change and it
re-presents.

## What gets generated

Everything is plain text in your repo: a README that holds the workflow's
rules (what each stage expects) and one file per work item. The generated
rules are a starting point; tighten them before any work runs, because an
agent working to a vague bar is expensive to correct.

## The pilot run

On accept, commission starts moving your seed items through the stages until
everything is idle or waiting on you. Check where things stand at any time:

```bash
spacedock status --workflow-dir docs/ship-features
```

```
ID     SLUG                           CURRENT              NEXT                 WORKTREE
--     ----                           -------              ----                 --------
001    add-rate-limiting-to-the-api   review               done                 yes
002    fix-the-flaky-login-test       implementation       review               yes
```

When a work item reaches a gate (here, `review`), Spacedock pauses and
presents a stage report: the chosen direction, the evidence behind it, and a
single recommendation. You approve, send it back with feedback, or reject.
[Gates and decisions](../concepts/gates-and-decisions.md) covers the full
machinery.

From there you are running the workflow: approving, sending back, and resuming
in later sessions. [Commission a workflow](../running-workflows/commission.md)
covers every design decision in full;
[Operating a workflow](../running-workflows/operating.md) covers the
day-to-day loop.
