# Implementation design-reset gate: one-shot advisory rounds

## Capability and reviewed change

The bounded correction fixed the named structural defects and kept the full and race suites green, but the two-step prefix-append design still measures 670 net production LOC before CLI, above the revised 540-LOC hard stop.

## Evidence

- The second boundary audit found no hidden duplicate recorder: shared entity/frontmatter I/O is 186 LOC, canonical review parsing and triage is 147 LOC, round orchestration is 288 LOC, and model/operation routing is 49 LOC.
- A compact two-step design has a credible 610–630 pre-CLI floor and would require another large ceiling increase.
- Interim reviewer-only persistence and strict-prefix extension are not value acceptance criteria. The operational caller records a round after worker triage.
- A one-shot completed-round operation preserves no-findings versus all-declines, immutable retained evidence, exact replay, divergence refusal, entity CAS, rollback, digest validation, and no gate/application/status effects.
- The current checkpoint remains a green counterexample and must not be compacted by weakening those guards.

## Recommendation and decision

Recommendation: **revise through a narrow design reset**. Replace only mutable prefix-append semantics with one-shot completed-round publication. Preserve every value acceptance criterion and use hard ceilings of 580 net production LOC before CLI and 640 total.

Decision: revise to ideation; do not raise the two-step ceiling and do not park the capability.
