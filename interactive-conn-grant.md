---
title: "Model the interactive conn grant: scope, quotable phrases, and gate-resolution semantics in the shared core"
status: backlog
source: "FO Commander session 2026-07-06/07 (0250 drive): the captain granted 'you have the conn... authorized to approve gates, merges, ci-env' in an INTERACTIVE session. The shipped contract models given-the-conn only for headless runs (claude runtime: 'the self-approval guardrail is absolute in interactive sessions and in any headless run NOT given the conn'), so the FO handled two gates inconsistently under one grant: k7's validation gate conn-resolved, tv's ideation gate presented-and-waited — a self-imposed stall the captain had to clear with a second grant. The gap: interactive conn is contractually void, so its semantics get improvised per gate."
started:
completed:
verdict:
score: 0.4
worktree:
issue:
id: 1jpkr9a86ydk80hmc3as9r88
---

The given-the-conn exception exists only in Startup step 8's headless branch; an interactive captain grant of gate/merge authority has no contract model — no quotable-phrase rule, no scope semantics (whole-workflow vs named members vs sprint), no statement of how conn-resolution interacts with the interactive 'never self-approve' guardrail, and no requirement that the FO still render the gate review as an audit record when resolving under the conn. Proposed direction: extend the shared core's Completion-and-Gates (and the runtime adapters' guardrail paragraphs) to define the interactive conn: the grant must be quotable, its scope is what the captain names (defaulting to the named workflow/sprint goal), it makes gate resolution the FO's action while REQUIRING the rendered gate review in the transcript as the decision record, and it never extends to surfaces the captain reserved (e.g. a release cut). Acceptance sketch: value — a live interactive drive with a quoted conn grant resolves an FO-recommended-approve gate with zero captain round-trips while the transcript carries the full rendered review (baseline: this session's tv stall, one round-trip); mechanism — the contract clauses ship and a live drive observes them.
