# Smallest Sufficient Mechanism

The ladder, the busywork the FO refuses, and worked examples behind the Working-Principles gate. Loaded when the FO faces a discretionary mechanism choice — never at boot.

## The ladder

Two ordered rungs. The FO picks the LOWEST rung that does the job and names in one line why the rung below cannot:

- **Action weight:** do-nothing → in-house Read/Edit → one dispatched worker → a workflow.
- **Publication weight:** direct commit → PR (with CI lanes).

Climbing a rung is justified ONLY by one of:

- **genuine fan-out** — many independent units of work a single in-house pass cannot cover;
- **required isolation** — parallel mutation that would collide in one working tree;
- **independent adversarial verification** — a fresh reviewer whose value is being separate from the author.

Re-running a verification a stage already owns is a climb in the other direction — re-doing work below the smallest rung. It is justified ONLY when the stage's report shows the required check did NOT actually run green: an unapproved, skipped, red, or absent lane is not a pass (the evidence bar above). A green report is consumed once — read it, then a cheap spot-check that the required lane actually ran green against this HEAD. Never a reflexive full re-run.

Never a justification: "it's substantive," "Ultracode is on," "I'm the dispatcher so I don't touch files," or "let me double-check." Ultracode raises the thoroughness of the ANSWER — its coverage and correctness — not the weight of the MECHANISM: a direct in-house multi-file Edit IS the exhaustive-correct answer.

## Scope: a discretionary climb, not a commissioned dispatch

The gate binds the FO's DISCRETIONARY mechanism choice for an ad-hoc task. It does NOT re-gate the standing dispatch a commissioned workflow stage already declares: an `«engage»` loop dispatching ready entities via `«dispatch.next-action»()` is executing an already-declared mechanism — justified when the workflow was commissioned, not re-narrated per entity. The gate fires on a discretionary climb; it never turns "fire and walk away" engage into per-entity justification narration.

The gate binds ALL FO action — release machinery and sprint shaping included, no domain is a framing escape hatch — but "all action" means every discretionary choice, not the commissioned dispatch above.

## Named busywork the FO refuses

Over-orchestration — climbing ABOVE the smallest rung:

- A PR with CI lanes for convention-direct prose. Roadmap and state docs commit directly (the precedent: `index.md` and `dispatch-sprint-execution.md` committed direct, never PR'd).
- Dispatching deterministic edits whose verbatim content the FO already holds.
- A workflow with no fan-out, no isolation, and no verification need.
- Re-formalizing work already done.

Redundant re-verification — re-doing work a stage already owns:

- Re-running a stage's own validation — a worktree + `go build` + a full `go test ./...` — inline at gate time to "double-check" a fresh validation stage whose report already covers it, instead of reading the report and doing a cheap spot-check that the reported check ran green.

The first four climb above the smallest rung; the last re-does a stage's owned work. The same one-line justification is owed before EITHER. Smallest-sufficient is a floor AND a ceiling: don't climb, and don't re-do.

## Worked examples

**Over-orchestration — refused.** Task: apply a handful of known edits to state entities whose verbatim content the FO already holds, and commit a roadmap strategy doc. The wrong climb (an observed baseline): a dynamic workflow plus a dispatched worker for the edits, and a PR with CI lanes for the doc. Smallest sufficient: in-house `Edit` for each file and a direct `git commit` for the doc — no fan-out (one FO applies them all), no isolation need (serial edits, one tree), no independent-verification need (deterministic, content already known). The doc rides the convention-direct precedent.

**Redundant re-verification — refused.** Task: gate a fresh validation stage whose report shows `go test ./...` ran green against this HEAD. The wrong climb (an observed baseline): spin a worktree and re-run `go build` plus a full `go test ./...` inline to "double-check." Smallest sufficient: read the validator's report and do a cheap spot-check that the required lane actually ran green against this HEAD. The full suite is the validator's owned work; re-producing it is redundant. Only a RED, skipped, or absent lane in the report routes back to the validator to run green — the FO never re-runs the suite itself.
