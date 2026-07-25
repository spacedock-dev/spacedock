# Ideation gate: make gate presentation completion observable

Capability/change: Make `present-gate` the sole owner of the six-field captain review; keep `fo-gate-lifecycle` responsible for bind, close, consume, and waiting for presentation completion.

Test and evidence: Two unchanged goal-only Codex runs loaded both skills, judged the gate, and immediately recorded approval without emitting a qualifying root review. The durable excerpt preserves that exact sequence. The revised tests reject reviews outside bind-to-decision, non-root reviews, zero reviews, and duplicate qualifying root reviews.

Reviewed snapshot: State commit `cf5d5c1e`; failed product tip `37d6980b` remains local and unpushed; PR #565 remains at accepted remote tip `13d70249`.

Findings:

- Material: none after repair. Independent staff review confirms the duplicate-root proof gap is closed.
- Deferred: Claude and Pi have separate extractors; implementation must apply the plural exact-one requirement to their existing focused tests when extraction is not shared.
- Polish: earlier cycle-20 arithmetic is superseded by the appended repair report.

Recommendation: approve. The eight-file repair gives one presenter owner and a falsifiable completion boundary without a command, schema, controller, compatibility path, host copy, or second harness.

Decision ask: approve to repair the existing implementation branch within `+65/-31` incremental intent and a `+175` additions hard stop; revise to change the ownership or proof boundary; hold for a named prerequisite.
