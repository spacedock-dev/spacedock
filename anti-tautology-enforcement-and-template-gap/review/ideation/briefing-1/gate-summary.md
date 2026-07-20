# Gate: should we build an automatic checker for fake tests?

**Where you are:** reviewing one design in the current improvement sprint, before any code is written. The full design is the second artifact in this package (the task file itself).

**The problem:** our test suites contained tests that can never fail — they compare a file's text against the same file, or a program's output against a copy of that output. Fake tests give false confidence.

**The question this design settles:** build an automatic detector for such tests, or handle them with review discipline?

**The decision (recommended): no automatic detector.** The designer ran a real experiment: they deliberately broke a working function. A fake test stayed green — proving it checks nothing. A normal test went red. Conclusion: a quick manual spot-check during review reliably exposes fakes, while an automatic detector would wrongly flag many innocent tests (ordinary tests legitimately compare against fixed values) and still miss the subtler fakes.

**What ships instead — seven sentences of rule text, zero code, landing in two files:** the workflow rulebook (`docs/dev/README.md`, the testing-rules section) and the worker instruction file (`skills/ensign/references/ensign-shared-core.md`, the reporting rule):
1. Every task must state what change would make its tests fail. A test nobody can fail is fake by definition.
2. Reports must say what each test actually checks, not just "5 of 5 passed."
3. Your earlier ruling written into the rules: a one-time text search pasted as review evidence is fine; committing that search as a permanent test is banned.

**One separate yes/no for you:** we already run a manual deep review on risky changes. This design proposes ALSO triggering it whenever a test's expected answer comes from the same code being tested. That widens when the review fires, so by our own rules it needs your explicit approval. Droppable without harming the rest.

**Your decision:** approve = the rule text goes to implementation (with or without the widening — say which). Revise = annotate what to change. Hold = we discuss.
