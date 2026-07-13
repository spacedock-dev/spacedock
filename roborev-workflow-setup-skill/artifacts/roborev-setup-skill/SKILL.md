---
name: roborev-setup
description: Set up Roborev as independent code-review evidence in an existing Spacedock development workflow. Use when a user asks to adopt or add Roborev to a code workflow with implementation and fresh-validation stages.
user-invocable: true
---

# Spacedock Roborev Setup

Help a Spacedock user adopt Roborev as an evidence producer. Keep Spacedock authoritative over stage state, implementation work, behavioral validation, rejection routing, and human gates.

Never invoke `roborev fix`, `roborev refine`, install the Roborev agent hook, or let Roborev choose the next workflow stage. Roborev must not modify code in a Spacedock-managed workflow.

## Assess fit

Recommend this integration when the project:

- develops code in Git branches or worktrees;
- has implementation followed by fresh validation;
- routes validation defects back to implementation;
- benefits from persistent or multi-reviewer evidence;
- trusts the repository code, or isolates the whole Roborev daemon and its agents.

Do not force it onto prose-only workflows or projects whose trust boundary cannot tolerate external review processes.

## Start with Roborev's own playbook

1. Discover the canonical project root, workflow README, implementation and validation stages, base branch, worktree policy, and split state checkout.
2. From the canonical code repository, run `roborev quickstart`. This must be the first Roborev command. It is read-only and its **Current state** and **Configuration playbook** are the primary setup authority.
3. Use the reported current state to distinguish missing setup from intentionally omitted features. In a Spacedock-managed workflow, leave Roborev fix/refine skills and the agent hook uninstalled even when `quickstart` marks them missing.
4. Explain the integration boundary, the missing setup items that matter, and the proposed edits. Show any proposed `review_guidelines` calibration separately. Get approval before installing software, initializing Roborev, changing repository or workflow files, or changing machine-local configuration.
5. Apply only approved missing steps. Re-run `roborev quickstart` afterward to verify the resulting state.

Do not begin with web search or a tour of Roborev's documentation. When `quickstart` leaves a specific command-syntax gap, use `roborev <command> --help`. Only when that still leaves a named configuration gap should you open the directly relevant official page, such as [Configuration](https://roborev.io/configuration/) or [Subagent Review Panels](https://roborev.io/advanced/subagent-review-panels/). State the gap before opening the page and return to the current-state checklist afterward.

## Configure the workflow boundary

- Configure an inexpensive `quick` panel as `[review] hook_review_panel` for post-commit review and a required multi-reviewer `code_completion` panel for implementation exit. Select available agents with the user; do not invent agent or model names.
- If the workflow has a split state checkout, resolve its actual checked-out branch and add that branch to the state checkout's `excluded_branches`. Do not enqueue a state-only review to probe the exclusion.
- Update the implementation and fresh-validation stage definitions with the exact contract below. Do not leave the panel order only in setup explanation.
- Check daemon reachability and report operator-only work without installing services or authenticating agents on the user's behalf.
- Detect Safehouse or another sandbox. When present, advise the user about the smallest read-only runtime access needed by sandboxed clients. Omit sandbox advice otherwise.

Treat `review_guidelines` as reviewer calibration, not a second instruction file. Derive only context a reviewer cannot infer from the diff and repository: trust boundaries, intentional compatibility posture, known false-positive suppressions, and boundaries on what the review should evaluate. Never duplicate `AGENTS.md`, workflow or stage instructions, component procedures, or required developer and test commands. Show the proposed calibration to the user and get approval before writing it. Keep implementer process guidance in `AGENTS.md` or the generated implementation-stage instructions instead.

## Exclude split workflow state during setup

Roborev may register a split state checkout as a separate repository path, so the code repository's `.roborev.toml` may not govern it. Put the smallest documented configuration in the state checkout and name the actual state branch:

```toml
excluded_branches = ["<state-branch>"]
```

This is a setup action, not a live enqueue test. Never deliberately create a state-only review to prove the exclusion.

## Run an implementation-owned semantic pass before review

Add this pre-review discipline to the generated implementation-stage instructions, never to `review_guidelines`. Before the final candidate commit:

- Trace each changed value or event through every representation and lifecycle phase. Check identity, cardinality, order, exact bytes, attribution, authority, and terminal state.
- Matrix the adjacent variants most likely to invalidate an apparently passing change: empty and terminal states; repeated and out-of-order events; every input path; Unicode, EOF, size, visibility, and layout boundaries.
- Use the repository's canonical validators. When several fields form one record or decision, validate the complete record atomically instead of accepting individually plausible fragments.
- Inspect hot paths and their readers for multiplicative work, blocking I/O, unbounded allocation, and implicit size limits. Add scaling and over-limit proof where the path can grow.
- Ask how the existing tests could pass while observable behavior is still wrong. Assert the exact result plus failure, cleanup, and terminal-state behavior.

This pass is implementation work: the worker fixes what it finds, strengthens proof, and commits the final candidate before entering the quick cost gate. Roborev reviews the resulting change; it does not own this pass or its fixes.

## Put `quick` before `code_completion`

Add this invariant to the implementation stage definition or dispatch checklist:

1. After the final candidate commit, resolve and record the exact candidate `HEAD`.
2. Wait for the already-enqueued post-commit `quick` review for that exact commit with `roborev wait <head>`, then inspect the stored exact-tip job with `roborev show <head> --json`. Verify that it is the expected quick panel before judging findings. Do not start `code_completion` in parallel or before this wait and inspection complete.
3. If the exact-tip quick review is missing, has an execution error, or reports a Medium-or-higher finding, remain in implementation. Repair coverage or fix the finding, commit, and repeat the quick wait for the replacement tip. `wait` may exit nonzero for a Low-only failed verdict; the stored finding severities decide this cost gate, and Low findings remain advisory without forcing another commit.
4. Only after that exact-tip quick result clears the severity floor, launch the required panel against the complete branch. Current CLI syntax is:

```bash
roborev review \
  --repo <canonical-project-root> \
  --branch=<implementation-branch> \
  --base <base-branch> \
  --panel code_completion \
  --min-severity medium \
  --wait
```

5. Preserve the synthesis parent job ID printed by the command, including when `--wait` exits nonzero. Fetch the stored record with `roborev show --job <parent-id> --json`.

Any code-changing commit invalidates both the prior quick result and the prior `code_completion` evidence. The replacement tip starts again at the quick wait. `quick` is a cost gate; it never substitutes for `code_completion` as implementation-exit evidence.

Implementation is complete only when:

- the stored range equals `merge-base(base, head)..head`;
- the reviewed head equals the current implementation tip;
- the panel name is `code_completion`;
- every configured required member appears exactly once and completed without an execution failure;
- the synthesis parent completed with verdict PASS.

The synthesis parent is the review authority. Individual member verdicts are inputs, not separate gates.

Record compact evidence in the implementation Stage Report: Roborev version, parent ID, exact range and head, panel, member execution outcomes, parent verdict, and finding dispositions. When the full synthesis contains findings or would obscure the report, store it as a folder-form entity artifact and link it. Do not duplicate the full synthesis in both places.

## Keep fixes and routing in Spacedock

For every synthesized finding, implementation records one disposition:

- `fix` — implementation changes code and strengthens proof;
- `rebut` — implementation cites repository evidence, then requests a replacement panel;
- `needs decision` — the First Officer stops for a genuine product, contract, or policy choice.

Roborev supplies evidence only. It does not run a fix/refine loop, make commits, advance entity state, or route a rejection. Ordinary review-fix-review iterations remain inside implementation and do not count as validation feedback cycles.

## Verify evidence during fresh validation

The fresh validator fetches the recorded synthesis parent and verifies its frozen range, current head, panel, required-member execution, and passing parent verdict. It does not rerun an unchanged passing panel.

Validation remains behavior-focused: reproduce every acceptance criterion, run required regressions, and exercise claimed runtime behavior. If validation finds a code defect, route it to implementation. Any fixing commit invalidates the prior panel and requires a new exact-tip quick gate and passing `code_completion` replacement before validation resumes.

## Apply failure policy

- Advisory post-commit review fails open for the commit; missing coverage is not a clean review and must be repaired before the final-tip cost gate. Inspect the stored quick result instead of treating `wait` exit status alone as the severity decision.
- The required `code_completion` panel fails closed on findings, absence, daemon failure, database or disk failure, required-member execution failure, or synthesis failure.
- Keep a failed required panel in implementation; do not present it as a human approval gate.
- Do not count ordinary implementation review iterations as validation feedback cycles.

## Respect the sandbox boundary

When a sandbox is present, advise the user to run one machine-local Roborev daemon outside it for trusted repositories. Install and authenticate reviewer agents there. Let sandboxed Spacedock workers act as loopback clients and expose only the runtime path they need, typically:

```text
add-dirs-ro=~/.roborev/runtime
```

Prefer a system or Homebrew binary. Do not expose all of `~/.roborev`, commit machine-local Safehouse settings, or move `ROBOREV_DATA_DIR` into the project. The daemon and reviewer agents do not inherit Safehouse isolation; isolate the whole daemon when reviewing untrusted code.
