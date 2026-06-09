# Proof policy

A change is proven by exercising its behavior and observing the outcome — not by asserting that a file contains the right words. This is the discipline the development workflow's ideation and validation stages enforce, and it is the bar any contribution must clear.

## Behavior is proven by exercising it

Prove a claim by running the behavior and checking the result: output bytes, exit code, resulting on-disk state, or a test feeding many inputs and asserting uniform handling. Go unit tests for parser and command behavior, golden fixtures for stable command output, behavior fixtures that drive the binary for command-level claims, and live workflow smoke tests when runtime behavior is the claim.

## No prose-grep acceptance criteria

A string, substring, or regex match over an instruction file the model reads — the first-officer or ensign contract, a workflow README, a skill — never satisfies a behavioral acceptance criterion. The matched string was written by the same implementer the check is supposed to police, so a passing check only asserts "the file contains the text we put in the file." It has no independent source of truth: a valid paraphrase fails it, and an inverted clause passes it.

The one test that settles it: **does the expected value come from somewhere other than the file under test?** If no — the clause is its own expectation — the check is a tautology and proves nothing. If yes — the file is bound to an independent source that can diverge from it — it may be a legitimate invariant. The legitimate case parses real artifacts in code and tests a relationship between independent values; for example, that the plugin manifest's contract range brackets the binary's contract version.

## Acceptance criteria are end-state properties

Acceptance criteria describe properties of the finished entity, not stage actions. Each criterion's "Verified by" must name something outside the entity body that can fail: a test, a command's output or exit code, a file the change produces, or the resulting on-disk state. An acceptance criterion whose only proof is reviewing the entity's own prose can never fail, so it is not an acceptance criterion.

## The detached adversarial audit

For high-stakes surfaces — the front-door launcher, the state-mutation guards, the shipped scaffolding, and the CI/release machinery — a passing validation is necessary but not sufficient. A read-only adversarial audit on a throwaway checkout tries to refute the validation: construct an edit that breaks the claim and confirm the deliverable's own tests catch it. A test that stays green under an edit that breaks the claim is a hole. See [Gates & decisions](../concepts/gates-and-decisions.md) for how findings route back.

The authoritative statement of this policy lives in the development workflow's stage definitions — see [The development workflow](development-workflow.md).
