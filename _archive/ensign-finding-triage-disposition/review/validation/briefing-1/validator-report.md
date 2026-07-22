## Stage Report: validation

- DONE: Independently validate exact candidate e85eb0cf and the complete four-commit branch against the current entity, approved reframe, captain reconfirmation, and sprint release line; do not rely on the implementation report alone.
  Exact clean head `e85eb0cfcc3c243fd94754be2baafa23be302a21` is four commits over merge-base `fa240a76`; direct diff, entity, durable-decisions DoD, and captain-reconfirmed class/ratchet inputs agree.
- DONE: Reproduce AC-1 as the required live value replay: seed the correct-but-disproportionate finding, exercise the supported triage delivery/record path, and prove an explicit advisory decline with zero product-line fix against the recorded incident-13 nonzero dutiful-fix and 85 prose-only baselines. Do not skip this live replay.
  A fresh routed worker read the shipped trigger: Case A produced a linked decline Annotation plus advisory Resolution and `0/0` product lines; material Case B fixed `product/status.txt` by `+1/0`; both statuses remained `implementation`, unlike the recorded incident-13 fix and 85 prose-only disposition.
- DONE: Attack AC-2 with the offline materiality fixture/check: independently reproduce all valid cases and both red controls, reject unknown classes, and use claim-breaking edits to prove each materiality conjunct and class allowlist is load-bearing.
  `docs/specs/check-finding-triage-materiality.sh` returned 8 ACCEPT/2 REJECT; removing user/workflow, harm, boundary, or trigger independently failed its paired control, and removing the allowlist admitted `defered-risk` and failed.
- DONE: Validate AC-3 end to end: the disposition is advisory and non-advancing; AC narrowing/design reset is captain-owned binding; trigger delivery is present; ensign-shared-core remains unchanged; no product code, recorder round plumbing, new schema field, or second record section entered the diff.
  Live status stayed `implementation`; the trigger and owner-tagged contract carry advisory/no-application semantics and captain binding, while exact-range path/field audits found none of the forbidden surfaces and zero `ensign-shared-core` delta.
- DONE: Verify the exact captain-authorized +703-byte all-host prompt-ratchet rebaseline and prove there was no additional prompt growth or semantic expansion beyond that measured exception.
  The sole changed loaded path is `feedback-rejection-flow/SKILL.md`, `3329→4032` bytes; measured host loads equal Claude `96081`, Codex `75296`, Pi `71426`, each exactly baseline +703, and the ratchet test is green.
- DONE: Audit actual surface against the reconfirmed boundary: zero production LOC, 47 fixture/test lines, 48 docs/skill/spec lines, no more than the approved 95 added-line ceiling after the two recorded correction rounds.
  `git diff --numstat main...HEAD` reports exactly 95 insertions: 0 production LOC, 47 fixture/test insertions (44 fixture/check + 3 test-baseline replacements), and 48 docs/skill/spec insertions across six paths.
- DONE: Verify Roborev jobs 548, 554, and 560 and their durable triage: material evidence defects fixed; every decline names why it is not material and a promotion condition; job 560's all-declines outcome is observably distinct from no findings.
  `roborev show 548/554/560` matches commits `9b2093b5` and `e85eb0cf`; state commits `584cae01`/`55473ae8` precede fixes, and the cycle-3 advisory record gives every decline grounds/promote conditions with job 560 recorded as `0 fixed; 3 declined` rather than absent.
- DONE: Run the detached adversarial audit required for the shipped feedback-rejection skill/contract surface, then run the focused offline check, gofmt -w ./cmd ./internal, go test ./..., go test ./... -race, git diff --check, and cleanliness checks at the exact committed head.
  Detached checkout at `e85eb0cf` was CLEAN; focused check, formatter no-diff, full suite, race suite, and `git diff --check` all exited 0, and the implementation worktree stayed clean at the exact head.
- DONE: AC-1 (VALUE) — live behavior discriminates rather than declining everything.
  Case A's valid four-record JSONL has an explicit linked advisory decline and zero product diff; Case B's material control has a nonzero product fix, with no entity-status advance in either arm.
- DONE: AC-2 — the recorded decline class is falsifiable against independent four-field evidence.
  Both required red controls reject, and five independent claim-breaking mutations each exit 1 on the specific case whose evidence or allowlist they weaken.
- DONE: AC-3 — advisory authority, binding design reset, trigger delivery, and the negative scope boundary hold together.
  Contract/trigger audit plus the live non-advance observation prove the current semantics; exact-range audit confirms no application/plumbing/schema/product expansion.
- DONE: PASSED recommendation with deferred risks separated from material findings.
  No material finding remains. Deferred only: malformed/external fixture hardening (promote on external support/schema AC), CI auto-wiring (promote if release proof requires it), and advisory `decision` consumer ambiguity (promote when round plumbing branches on it).

### Summary

Fresh validation recommends **PASSED** for exact candidate `e85eb0cf`. The required routed live replay moved the value baseline—an explicit advisory decline with zero product-line fix versus a nonzero material control—and the offline evidence, prompt bytes, scope ceiling, Roborev triage, detached audit, full/race suites, and clean-head checks all held independently. No material finding remains; the three detached-audit risks stay deferred under the explicit promotion conditions above.
