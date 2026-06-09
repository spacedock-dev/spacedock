# Proof policy

In the Spacedock development workflow, behavior is proven by exercising it and observing a durable result — not by matching text in the files the model reads. This page summarizes the proof discipline a contributor must satisfy. The authoritative statement, stage by stage, lives in [the development workflow](development-workflow.md); when this page and that one disagree, that one wins.

## Prove behavior by exercising it

**A behavioral acceptance criterion is proven only by running the behavior and observing the outcome** — output bytes, exit code, resulting on-disk state, or a test that feeds many inputs and asserts uniform handling. In practice that means Go unit tests for parser and command behavior, golden fixtures for `status` output, behavior fixtures that drive the binary for command-level claims, and live workflow smoke tests when runtime behavior is the claim.

A string, substring, or regex match over an instruction file the model reads — the first officer or ensign contract, the workflow README, a skill — never satisfies a behavioral acceptance criterion. The matched string was written by the same implementer the check is meant to police, so a passing check only asserts "the file contains the text we put in the file." It has no independent source of truth: a valid paraphrase fails it, and an inverted clause passes it ("It is a myth that the FO MUST advance…" keeps every matched substring while saying the opposite). Searching code is the same trap one level down — it asserts spelling, and false-passes on a renamed-but-equivalent branch.

So "the contract says to run the command" never proves the agent runs it. "The skill renders the gate" is proven by invoking the skill and observing the rendered gate, never by finding the clause that asks for it.

The one test that settles it: **does the expected value come from somewhere other than the file under test?** If no — the clause is its own expectation — the check is a tautology and is banned as proof. If yes — the file is bound to an independent source that can diverge from it — it may be a legitimate invariant. The legitimate case parses real artifacts in code and tests a relationship between independent values; for example, that the plugin manifest's contract range brackets the binary's contract version. The binary parses that manifest, so it is already outside "files the model reads," and the two versions can disagree — which is exactly what lets the check fail.

## Acceptance criteria are end-state properties

An acceptance criterion names a property of the finished entity, not a stage action. If a draft AC reads as an imperative verb phrase, rewrite it as the end-state property it produces; the imperative belongs in the stage report's checklist instead.

Each AC carries a `Verified by:` clause, and that clause must name something **outside the entity body** that a future reader can reproduce and that can fail:

- a test name,
- a command's output or exit code,
- a file the change produces, or
- the resulting on-disk state.

An AC whose only proof is review of the entity's own prose ("verified by reviewing this task's decision section") can never fail, so it is not an acceptance criterion. Every task must produce a real, checkable change — code, a fixture, on-disk state, or instruction text whose effect a separate check confirms. If a task's only output is a decision with nothing shipped, it does not belong in the dev queue; record the decision in the roadmap instead.

Prefer an AC a code gate can enforce — a guard in the binary, a test that fails on violation — over one the agent is merely instructed to follow. Where a behavior can be guarded by the binary or a failing test, the proof is that gate, not a sentence in a skill file.

When a task has both a text half and a behavioral half — extracting contract prose into a skill, or adding a contract clause — the text half is real authoring work but is not an acceptance criterion on its own, and the text check never stands in for the behavioral one. That invoking the skill renders the gate, or that the FO obeys the clause, is proven only by a live drive that runs it and observes the durable result.

## The instruction-file-read quarantine

**Tests do not read prompt or instruction files except in `internal/contractlint`,** and that package is limited to structural checks a machine can verify without interpreting prose: reference closure, frontmatter validity, structural absence, deduplication, line/count floors, and portability markers. The package's `doc_test.go` states the policy; `boundary_guard_test.go` and `instruction_read_detector_test.go` enforce it — the boundary guard fails on an instruction-file read outside the quarantine.

Prose-grep is banned: a test asserting that a skill, contract, agent file, or the workflow README contains its own wording proves only that the wording is present. Code-bound prose checks are banned too — a prose-to-code consistency lint is not a behavior test and must never substitute for running the behavior. If a deleted read exposed an untested behavior that still matters, record the owed behavior test rather than keeping the read.

## The detached adversarial audit

For high-stakes surfaces, a passing validation is necessary but not sufficient. Before merging such a change, run or dispatch a read-only adversarial audit on a detached, throwaway checkout of the merge result.

- **When it triggers.** The four high-stakes surfaces: the front-door launcher (`spacedock claude` / `codex` / `doctor`), the `status` mutation and guard paths, the shipped contract and scaffolding, and the CI/release machinery. Routine, low-blast-radius changes do not need it; a normal validation suffices.
- **What it produces.** The auditor works on a separate checkout — never the implementation worktree — and never mutates the deliverable. It tries to **refute** the validation: construct an adversarial edit that the deliverable's own tests should catch, and confirm they do. A test that stays green under an edit that breaks the claim is a hole. Findings come in two tiers — `Material:` for a real correctness or test-strength hole, `Polish:` for non-blocking notes. "Refuted nothing material" is itself a valid, recorded outcome.
- **How it is recorded.** Material findings route back through the normal validation → implementation feedback flow as a `### Feedback Cycles` entry naming the audit and its adversarial edit; the gate is not presented as clean until they are closed. A clean audit is noted in the gate's reviewer-findings block.

The audit catches the class of hole where the test passes but would also pass on a broken future edit — which validation, trusting its own green suite, cannot see. Real catches are on the record in the development workflow, including a `strings.Count(...) > 0` check that skipped on zero mentions and a bare `strings.Contains` satisfied by a negated disclaimer.

## See also

- [The development workflow](development-workflow.md) — the authoritative stage-by-stage proof rules and the entity field reference.
- [Agent development](agent-development.md) — the write-scope rules and the "prove runtime claims with durable state" discipline.
