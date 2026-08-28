# Backlog decision: right-size recent proof scaffolding

## Recommendation

APPROVE ideation now. The product fixes are valid, but their proof portfolio duplicates behavior coverage and mislabels a local candidate package as a stable release. Ideation must retain proof for each distinct failure mode while removing redundant committed scaffolding.

## Outcome

The cleanup will reduce CI cost and maintenance without weakening behavior coverage. Each failure mode will have one primary proof owner. One targeted live journey will prove runtime behavior. A manual install from the real release channel will prove published-package provenance.

## Included scope

- Remove or combine redundant tests from PRs #767, #768, #776, #777, #780, and #781.
- Remove source-text and generated-prompt checks used as proxies for runtime behavior.
- Remove the duplicate Pi live journey.
- Stop describing a copied checkout as a published stable package.
- Add a concise workflow rule for choosing between committed tests and manual validation.
- Implement the cleanup as a new layer above commit `37f50588aa37cfa571a88e4aa87b8f5c8f1b39e8`.

## Excluded scope

- Product behavior changes.
- New test frameworks, lanes, fixtures, or recurring gates.
- The older dead-oracle task `xaz`.

## Proof owed by ideation

Ideation must assign one primary proof owner to each acceptance criterion and justify every additional check with a distinct failure mode. It must refine the initial estimate of -450 net lines across 14 files and show that each retained test fails when the behavior it protects is defective.
