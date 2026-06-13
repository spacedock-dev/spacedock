# Multiple workflows

Related work rarely fits one pipeline. A project can hold several workflows, each its own directory with its own README, stages, and gates; Spacedock finds them all and operates them at the same time. A typical pair: a benchmark project where one workflow tracks the experiment runs and another ships the harness that supports them.

There is no hard dependency declaration between workflows. When an entity in one workflow depends on work in another, annotate it in the entity's frontmatter (a flat custom field naming the other item) and the first officer figures it out at dispatch time.

Workflows in the same repo can also keep their state off your code branch; [split-root state](split-root-state.md) covers that.
