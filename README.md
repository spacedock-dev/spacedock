# Spacedock

**Spacedock is the judgement layer for agent work.**

It works inside Claude Code and Codex, so there is no new agent to learn.
It runs on macOS and Linux.

| Feature | What it does | Start it | Code |
| --- | --- | --- | --- |
| **Spacedock Review** | Highlight what needs changing, add your comments, and submit your decision. The requesting agent receives it all. | `/spacedock:review plan.md decision` | [spacedock-dev/subspace](https://github.com/spacedock-dev/subspace) |
| **Spacedock Compass** | See where agent work departed from your goal, with evidence to help you course-correct. | `/spacedock:compass` | [spacedock-dev/compass](https://github.com/spacedock-dev/compass) |
| **Spacedock Workflow** | Define what agents must check, which decisions they can make, and when they need you. Workflow records each decision with its evidence and reason. | `/spacedock:commission` | This repository |

The commands are for Claude Code. In Codex, Review and Compass start with `$`
in place of `/`, for example `$spacedock:review plan.md decision`.

This repository holds Spacedock Workflow, the `spacedock` command and its
installer. One install brings Review and Compass with it.

## Install

```sh
# macOS with Homebrew
brew trust spacedock-dev/tap
brew install --cask spacedock-dev/tap/spacedock

# Linux, or macOS without Homebrew
curl -fsSL https://raw.githubusercontent.com/spacedock-dev/spacedock/main/install.sh | sh
export PATH="${SPACEDOCK_INSTALL_DIR:-$HOME/.local/bin}:$PATH"
```

Then connect your agent with `spacedock install --host claude` or
`spacedock install --host codex`. Run `spacedock doctor` to check the setup
(`spacedock doctor --host codex` for Codex). The
[installation guide](https://spacedock.md/docs/get-started/install/) has the
full steps.

Interactive launches check once a day for a newer Spacedock. The check asks
spacedock.md, which redirects to GitHub's releases API, so both see the
request. Add `--no-update-check` to a launch command, or set
`SPACEDOCK_NO_UPDATE_CHECK=1`, to skip it.

## Your first review

Paste this into Claude Code or Codex:

```text
Install Spacedock using https://spacedock.md/agents.md. Help me choose a plan to review, then open it in Spacedock Review.
```

Your agent installs Spacedock and opens the plan in Review. Comment on the
text, then approve, request changes or hold. The agent receives your decision
and comments.

## For agents

Read [spacedock.md/agents.md](https://spacedock.md/agents.md) for the install
steps and how to bring a first artifact to Review.

## Spacedock Workflow

A workflow gives recurring work a set of stages, such as design,
implementation and review. You set the standards for each stage. You choose
where work pauses for your decision; these pauses are called gates. The stages
and standards live in plain text in your repository.

To start from the work you already did, run `/spacedock:survey` in Claude
Code. It reads your local agent history, without changes, and reports repeated
work, current workstreams and unresolved decisions. Your agent can then
propose a workflow from the findings.

To describe a workflow yourself:

```sh
spacedock claude "/spacedock:commission ship features through design, implementation, and review"
```

Run `/spacedock:debrief` before you end a session. Run `spacedock claude` to
resume.

### What's different

- **The agent does not judge its own work.** Review runs as a separate stage
  with fresh context and no access to the maker's reasoning.
- **Every decision leaves a trail.** Each gate carries a stage report with
  findings, verdicts and artifacts. You decide on evidence, not the transcript.
- **The standards improve as you use them.** Each stage states what good
  means. When a standard is unclear in practice, the agent proposes an edit for
  your approval.
- **Work keeps moving.** Queue many items at once and handle gates as they
  arrive. When an agent runs out of context, a successor continues the work.

## Links

- [Documentation](https://spacedock.md/docs/)
- [Release notes](https://github.com/spacedock-dev/spacedock/releases)
- [Discord](https://discord.gg/3Ux2NKmSw)
- [spacedock-dev/subspace](https://github.com/spacedock-dev/subspace): Spacedock Review
- [spacedock-dev/compass](https://github.com/spacedock-dev/compass): Spacedock Compass
- [Contributing](CONTRIBUTING.md)
- [License: Apache License 2.0](LICENSE)
