# Validation correction gate: clean-runner Git identity

## Capability and reviewed change

PR #565 proves the candidate's offline fixture is not self-contained: commit-bearing recorded-gate tests inherit the developer's global Git identity and fail on a clean runner before exercising the lifecycle.

## Evidence

- GitHub Actions run 30112268591 job 89544401796 failed five commit-bearing recorded-gate cases with `Author identity unknown`.
- The fixture's `git` helper supplies `user.name` and `user.email` only to its own subprocesses.
- `spacedock state commit` launches a separate Git process, so it cannot see those transient `-c` values.
- Existing repository fixtures persist repository-local test identity before invoking product-owned commits.

## Findings

Material evidence defect: local full/race/live green depended on undeclared host configuration. The smallest correction is to persist test-only identity in the temporary state repository. No product, workflow, authority, or acceptance behavior should change.

## Recommendation and decision

Recommendation: **revise**. Return only this fixture setup defect to implementation, rerun the exact clean-config control and required suites, then require fresh validation and a replacement approval.

Decision: revise to implementation for the repository-local fixture identity; approve would merge evidence that is red on a supported clean runner; hold would preserve the open PR without correction.
