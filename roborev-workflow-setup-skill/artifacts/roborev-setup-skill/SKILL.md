---
name: spacedock-roborev-setup
description: Set up Roborev as independent code-review evidence in an existing Spacedock development workflow. Use when a user asks to adopt or add Roborev to a code workflow with implementation and fresh-validation stages.
---

# Spacedock Roborev Setup

Help a Spacedock user adopt Roborev as an evidence producer. Keep Spacedock authoritative over stage state, implementation work, behavioral validation, rejection routing, and human gates.

Never invoke `roborev fix`, `roborev refine`, or the Roborev agent hook in a Spacedock-managed workflow. Roborev must not modify code or choose the next workflow stage.

## Assess fit

Recommend this integration when the project:

- develops code in Git branches or worktrees;
- has implementation followed by fresh validation;
- routes validation defects back to implementation;
- benefits from persistent or multi-reviewer evidence;
- trusts the repository code, or isolates the whole Roborev daemon and its agents.

Do not force it onto prose-only workflows or projects whose trust boundary cannot tolerate the external review processes.

## Set up the workflow

1. Discover the canonical project root, workflow README, implementation and validation stages, base branch, worktree policy, and split state checkout.
2. Explain the integration boundary and proposed edits. Get approval before installing software, initializing Roborev, or changing repository configuration.
3. Follow Roborev's current official documentation. Treat it as the authority for commands and configuration schema:
   - [Quick Start](https://roborev.io/quickstart/)
   - [Configuration](https://roborev.io/configuration/)
   - [Subagent Review Panels](https://roborev.io/advanced/subagent-review-panels/)
4. Point the user to those documents to select available agents and configure:
   - an advisory `quick` panel for post-commit review;
   - a required multi-reviewer `code_completion` panel.
5. Run or propose the documented initialization from the canonical code repository, outside Safehouse when the local installation requires it.
6. If the workflow has a split state checkout, identify its actual state branch and add that branch to the checkout's `excluded_branches` during setup. Do not enqueue a state-only review to probe the exclusion.
7. Update both stage definitions:
   - implementation must produce passing, exact-head `code_completion` evidence;
   - fresh validation must verify that evidence without rerunning an unchanged panel;
   - any fixing commit invalidates the old evidence and returns responsibility to implementation.
8. Check daemon reachability and report any operator-only work.
9. Inspect the workflow, launcher configuration, and current runtime for Safehouse or another sandbox. When one is present, explain the boundary and advise the user about the minimum runtime access needed. Omit sandbox-specific changes when workers run directly on the host.

Keep review guidelines focused on the project's current contracts, ownership boundaries, trust model, compatibility policy, and meaningful failure classes. Do not copy Spacedock's own orchestration rules into a project that does not implement them.

## Exclude split workflow state during setup

Roborev may register a split state checkout as a separate repository path, so the code repository's `.roborev.toml` may not govern it. Put the smallest documented configuration in the state checkout and name the actual state branch:

```toml
excluded_branches = ["<state-branch>"]
```

This is a setup action, not a live enqueue test. Never deliberately create a state-only review to prove the exclusion.

## Require `code_completion` at implementation exit

After the code, tests, and implementation diff are ready, run the required panel using the syntax in the current Roborev documentation. A representative invocation is:

```bash
roborev review \
  --repo <canonical-project-root> \
  --branch=<implementation-branch> \
  --base <base-branch> \
  --panel code_completion \
  --min-severity medium \
  --wait
```

Preserve the synthesis parent job ID printed by the command, including when `--wait` exits nonzero for a failed verdict. Fetch the stored record:

```bash
roborev show --job <parent-id> --json
```

Implementation is complete only when:

- the stored range equals `merge-base(base, head)..head`;
- the reviewed head equals the current implementation tip;
- the panel name is `code_completion`;
- every configured required member appears exactly once and completed without an execution failure;
- the synthesis parent completed with verdict PASS.

The synthesis parent is the review authority. Individual member verdicts are inputs, not separate gates.

Record compact evidence in the implementation Stage Report: Roborev version, parent ID, exact range and head, panel, member execution outcomes, parent verdict, and finding dispositions. When the full synthesis contains findings or would obscure the report, store it as a folder-form entity artifact and link it. Do not duplicate the full synthesis in both places.

Quick-review coverage and Low findings are advisory. Note or backlog useful observations without forcing another implementation round. Do not set a repository-wide review severity floor merely to enforce the implementation exit; keep Low findings visible in ambient quick reviews.

## Adjudicate findings inside implementation

For every synthesized finding, record one disposition:

- `fix` — change code and add or strengthen proof;
- `rebut` — cite concrete repository evidence, then ask a replacement panel to decide;
- `needs decision` — stop for a genuine product, contract, or policy choice.

Keep ordinary fix-review iterations inside implementation. After any code change, run `code_completion` again against the new head. Cite both the failed parent and passing replacement in the next Stage Report.

The First Officer checks the dispositions and exact-head evidence. It does not rewrite code, substitute an informal review, or waive a required failure. Escalate only for a real semantic fork or failure to converge.

## Verify evidence during fresh validation

The fresh validator fetches the recorded parent and verifies its frozen range, current head, panel, required-member execution, and passing parent verdict. Do not rerun an unchanged passing panel.

Validation remains behavior-focused: reproduce every acceptance criterion, run the required regressions, and exercise claimed runtime behavior. If validation finds a code defect, route it to implementation. Any fixing commit invalidates the prior panel and requires a passing replacement before validation resumes.

## Apply failure policy

- Advisory post-commit review fails open. Missing coverage is not a clean review.
- The required `code_completion` panel fails closed on findings, absence, daemon failure, database or disk failure, required-member execution failure, or synthesis failure.
- Keep a failed required panel in implementation; do not present it as a human approval gate.
- Do not count ordinary implementation review iterations as validation feedback cycles.

## Respect the sandbox boundary

When a sandbox is present, advise the user to run one machine-local Roborev daemon outside it for trusted repositories. Install and authenticate reviewer agents there. Let sandboxed Spacedock workers act as loopback clients and expose only the runtime path they need, typically:

```text
add-dirs-ro=~/.roborev/runtime
```

Prefer a system or Homebrew binary. Do not expose all of `~/.roborev` or move `ROBOREV_DATA_DIR` into the project. The daemon and reviewer agents do not inherit Safehouse isolation; isolate the whole daemon when reviewing untrusted code.
