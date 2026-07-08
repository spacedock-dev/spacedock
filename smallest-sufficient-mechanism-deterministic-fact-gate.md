---
id: 1p9kve1s8r7j0ftd23nt9gnt
title: "Sharpen smallest-sufficient-mechanism: deterministic facts don't need adversarial verification, and Ultracode doesn't override the gate"
status: backlog
source: "FO self-diagnosis, 2026-07-08 live session: the FO spawned a 2-agent parallel Workflow to answer a deterministic code-fact question (does function A call function B) that a single direct grep settled in under a minute. Caught and corrected by the captain in-session; this entity formalizes the fix so the same gate hole doesn't reopen."
started:
completed:
verdict:
score:
worktree:
issue:
---

The smallest-sufficient-mechanism principle (`references/fo-smallest-sufficient-mechanism.md`) already bans "Ultracode is on" as a justification for climbing the action-weight ladder, and the FO still walked past it: triggered by the session's Ultracode directive ("use the Workflow tool on every substantive task"), it spawned two parallel general-purpose agents inside a Workflow call to verify whether two Go functions shared code — a fact a single `grep`/`Read` pass settles unambiguously, with only one correct answer either way. The existing "independent adversarial verification" justification bullet is loose enough to rationalize this: it doesn't distinguish a judgment call (where a second independent read can catch something the first missed) from a deterministic fact (where a second read of the same files finds the same thing, adding cost with no added confidence). Four concrete tightenings, sized for the lazy-loaded reference file (not the boot-resident core), consistent with this sprint's leanness discipline:

1. Sharpen "independent adversarial verification" with a discriminating question distinguishing a judgment call from a deterministic fact settled by reading N files.
2. Add an "N agents ≠ N confidence" corollary: before spawning more than one agent for the same question, name what a second agent could find that the first would miss; if the answer is "nothing, same files," spawn one or zero.
3. Resolve the Ultracode/zm interaction explicitly: a session-level thoroughness directive raises the bar on the answer a chosen mechanism must produce, never the weight of the mechanism itself; restating "the session says use Workflow" is not itself a fan-out/isolation/verification justification.
4. Make the one-line justification an artifact: write it as visible reasoning immediately before an Agent/Workflow call for a discretionary task; if it can't be written honestly in one line, the task hasn't cleared the gate.
