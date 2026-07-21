---
title: Fixture workflow READMEs must not be workflow-discovery candidates
status: ideation
source: "Live FO session, 2026-07-21, after the refit-content-propagation fixtures landed on main."
id: ab3ma8m7gsm8tra2ksmcdydq
started: 2026-07-21T16:05:13Z
---

Landing `fixtures/refit-content-propagation/site-workflow/` (a commissioned-shape README used as a template-propagation test fixture) made every auto-discovering command (`state commit`, `new`) ambiguous repo-wide: each now demands `--workflow-dir` for a repo that operationally has ONE workflow. The binary handles the ambiguity correctly (clear error, non-zero exit — verified live), so this is not an error-handling bug; the defect is that a test fixture counts as a real workflow. Fix direction, cheapest first: discovery skips paths under a directory named `fixtures/` (or testdata/), or fixture READMEs carry a marker field discovery respects; additionally, entity-scoped commands could resolve the workflow FROM the entity's state checkout, which is unambiguous by construction. Behavior test: discovery over a tree containing a fixtures-nested commissioned README returns exactly the real workflow. Session note recorded alongside: the FO's own `cmd 2>&1 | tail -1` piping masked a non-zero exit in an && chain — a shell-discipline lesson, not a binary defect.
