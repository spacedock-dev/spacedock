# First Officer gate review — validation

Entity: simplify-gate-state-v1-schema
Stage: validation
Candidate: `f566f821b`, based on `48a7ea0d97042f0e7aaac258e1b77f16157c5281`.
Scope: 25 files, 167 insertions, 202 deletions, net -35.

Validation report: lines 239–260 record 8 DONE items and 1 FAILED scope item. AC-1 through AC-5 have evidence. Focused, full, race, formatting, and stale-pointer adversarial checks pass.

Roborev job 768: the configured branch-final panel requested changes. It found two HIGH findings: strict v1 upgrade unreadability without one-time normalization, and terminal delivery selecting the target status instead of the source gate record. It also found medium summary and test-strength findings. Full evidence is in `review/validation/roborev-evidence.md`.

Science advisory:

- V-1 is a material boundary/evidence defect: `internal/gates/io.go` and `internal/status/discover.go` are required production consumers but are not covered by the recorded 25-file extension.
- HIGH-1 is material and belongs to the NTH normalization seam. Do not add a compatibility reader or migration.
- HIGH-2 is material only if the Roborev trace reproduces; confirm it before editing.

Recommendation: reject validation and route the exact findings to the JC/NTH implementation seam. Hold candidate bytes until the Captain authorizes the two production paths, NTH normalization is durable, and HIGH-2 is reproduced or cleared.
