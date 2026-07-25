---
title: Clean recorded-gate lifecycle proof hygiene before v1
status: backlog
source: "Durable-decisions exact-tip close-out audit at deac7f8a, corrected and committed as 4ff98d8c."
score: "0.9"
sprint: durable-decisions
id: fh3n4w4jg7tk015512tn1tsd
---

The recorded-gate lifecycle passes its supported outcome checks, but two test constructs make its evidence harder to trust and diagnose:

- `TestRecordedGateLifecycleCommandTextMutants` reads shipped skill prose outside the contractlint quarantine and deletes superstrings of the same tokens its test-local `procedureEvents` parser searches. Existing real-CLI replay, missing-event controls, refusal/resume matrices, and live journeys already own the behavior. Delete the tautological test and its orphaned helper; do not replace it with another prose check.
- The Codex live runner reads the active entity before comparing it with `resolveRecordedGateEntity`, so archival fails first as a generic read error and the named assertion is tautological. Replace it with the existing explicit pre-read archive-absence check and add the same diagnostic to the Claude runner. This improves attribution; it is not a missing product-outcome fix.

Ideation must verify the smallest current-main surface, inspect the adjacent in-quarantine `layering_restore_test.go` prose mapping without assuming it belongs in the task, and preserve every executable lifecycle control. No product behavior, command surface, prompt obligation, provider rule, compatibility layer, new detector, or new standing cap belongs here. Request Roborev after implementation and triage findings before edits.
