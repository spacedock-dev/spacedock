---
title: Approved-pending blockers, execution holds, and dispatch eligibility for recorded gates
status: ideation
source: "Split from the gate-recorder task (3k), captain-approved 2026-07-21. Carries 3k's original seed concern (captain design feedback 2026-07-13: an approval evaporated while its dispatch stayed blocked)."
id: h1y616vjh64wc961z5t1031d
sprint: durable-decisions
group: recorder
started: 2026-07-21T01:43:36Z
---

When a recorded approval is current but declared dispatch blockers remain, the entity must report approved-pending (or approved-held under a captain execution hold) and stay non-dispatchable; clearing the final blocker with reviewed content unchanged applies the approval exactly once through the existing transition path; any digest-bound input change first marks it stale. Fail closed: missing, ambiguous, or unqueryable blocker state never reads as satisfied and never consumes an approval. Scope moved from 3k at the split: its ACs 2, 3, 5, 11 and the eligibility subset of AC-6; scheduler rules 4-6 and the one-use guard of rule 10 beyond convention; AC-8's blocked/held/stale/duplicate-pass mutants. Authority for the record shapes stays 3k's gate-resolution-frontmatter-contract.md. Honest standing note from the split: the 0260 production dry run exercised none of this layer — all eight recorded approvals carried zero declared blockers and were consumed by the Commander through the normal path — so this task's own ideation re-examines live need before building (cheapest-check ordering; the recorder must land first regardless).
