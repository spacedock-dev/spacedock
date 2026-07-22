# Exact fixture: cross-logical-gate re-entry

This is the required first implementation behavior fixture for AC-10, AC-12, and AC-14.
Its exact source is state commit `71c61fbc`,
`durable-gate-approval-pending-blockers/index.md`, SHA-256
`c5b66ca3f97860b6e8e5f73b5d2fe3b44ae9e09fd1c3244aa9cf00c959227da8`.
The frontmatter is lines 1-238, SHA-256
`3dfd864747aa6469e2952949a433b05e9a59d466e4187bb5a2c820ea7dbf2e9e`.

## Before

The frozen source must be copied into a discovery-safe testdata path when implementation
begins. Its operative facts are exact, not synthesized:

```yaml
status: ideation
gates:
  version: 1
  current:
    gate: gate:docs-dev:3k:validation
    attempt: gate-attempt:3k-validation-1
  records:
    - id: gate:docs-dev:3k:ideation
      stage: ideation
      current-attempt: gate-attempt:3k-ideation-9
      attempts: # exactly nine, all closed
        # attempts 1-8 are the source bytes from commit 71c61fbc
        - id: gate-attempt:3k-ideation-9
          sequence: 9
          previous-attempt: gate-attempt:3k-ideation-8
          state: closed
          briefing:
            id: briefing:docs-dev:3k:ideation:attempt-9:revision-16
            digest: sha256:c99e7b8597038912b25f2d2f7fccd631649cc3b635fb57aa566d0ad25318aba9
            room-ref: ./review/ideation/briefing-16
          resolution:
            type: Resolution
            id: resolution:fo-delegated-3k-ideation-9
            briefing: briefing:docs-dev:3k:ideation:attempt-9:revision-16
            by: agent:first-officer
            at: "2026-07-21T14:35:06Z"
            decision: approve
          # application, notes, and the full reason remain exactly as in the source
    - id: gate:docs-dev:3k:validation
      stage: validation
      current-attempt: gate-attempt:3k-validation-1
      attempts:
        - id: gate-attempt:3k-validation-1
          sequence: 1
          state: closed
          briefing:
            id: briefing:docs-dev:3k:validation:attempt-1:revision-1
            digest: sha256:0a54f1baec0120c1c93523e6900a6ce28e025c570289e5dfa9835e28099042ac
            digest-domain: canonical-bytes
            room-ref: ./review/validation/briefing-1
          resolution:
            type: Resolution
            id: resolution:first-officer-3k-validation-1-design-reset
            briefing: briefing:docs-dev:3k:validation:attempt-1:revision-1
            by: agent:first-officer
            at: "2026-07-22T02:01:32Z"
            decision: revise
```

The exact ideation-attempt-9 source range (lines 192-213) hashes to
`d5fa69f5fb6a8eaee4698456d4b349d02c9f3742a94d224d13707ecb44058aab`.
The exact validation-attempt-1 source range (lines 218-233) hashes to
`073902d54503af26268665a4471e64f7c805aae988b73a8265cfc726185091e9`.

## Command

```text
spacedock gate record 3k --briefing revision-18.json --workflow-dir docs/dev
```

`revision-18.json` is the complete Briefing beside this fixture. Its canonical-bytes
digest is
`sha256:6b2c4f1388a58f42f7c8610f847ed9e7cce92758c00b201d4eb9f4f89dbedd8b`.

## After

The binary derives the logical gate from `status: ideation`, finds that gate even though
validation is globally selected, appends attempt 10, and selects it. The only logical
changes are:

```diff
 gates:
   current:
-    gate: gate:docs-dev:3k:validation
-    attempt: gate-attempt:3k-validation-1
+    gate: gate:docs-dev:3k:ideation
+    attempt: gate-attempt:3k-ideation-10
   records:
     - id: gate:docs-dev:3k:ideation
-      current-attempt: gate-attempt:3k-ideation-9
+      current-attempt: gate-attempt:3k-ideation-10
       attempts:
         # the nine closed attempts remain byte-identical
+        - id: gate-attempt:3k-ideation-10
+          briefing:
+            id: briefing:docs-dev:3k:ideation:attempt-10:revision-18
+            digest: sha256:6b2c4f1388a58f42f7c8610f847ed9e7cce92758c00b201d4eb9f4f89dbedd8b
+            digest-domain: canonical-bytes
+            room-ref: ./review/ideation/briefing-18
     - id: gate:docs-dev:3k:validation
       # the closed validation record remains byte-identical
```

`status: ideation`, every non-`gates` byte, ideation attempt 9, and validation attempt 1
must compare identical before/after. The appended attempt deliberately has no `sequence`,
`previous-attempt`, or `state`: revision-17's approved minimal-v1 rule stands.

## Required assertions and mutant

The CLI behavior fixture must assert exit 0; global selection and the ideation legacy
`current-attempt` both name attempt 10; there are ten ideation attempts; attempt 10 has the
exact Briefing binding above; both closure byte-range hashes are unchanged; and every byte
outside the three shown edits is unchanged.

The required mutant replaces target-stage gate lookup with the current implementation's
global-current prerequisite (`doc.Current.Gate == targetGate`, equivalently routing through
`currentAttempt` before selecting the target). It must make the full suite fail in this
fixture: the command would report a pointer conflict against validation instead of opening
ideation attempt 10. A one-gate re-entry test cannot kill this mutant, which is why this
dogfooded two-gate state is required.
