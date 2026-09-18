# Committed hold reason: independent validation

Recommendation: **REJECTED** for a concrete task-owned evidence defect at `9a614522dc1c1c162bc6b99fe11e5739f66c611b`. No candidate repair is authorized by this report. Producer state report: `48d3c1828`.

## Finding H1: missing-source meaning is not tied to its reason

Proposed defect kind: **evidence defect**. Proposed release scope: **Material**. Ownership: current same-stage hold grader. Proposed disposition: **FIX**, pending distinct FO authorization.

Four policy fields:

1. Released user and normal workflow: a required-review same-stage correction stops because external review evidence is absent, and the worker commits its task-body hold reason in ordinary prose.
2. Observable harm: the grader accepts a committed body that names the source but records only an unrelated missing screenshot, without saying why review is held. It also rejects a clear ordinary-language committed hold reason. Thus the outcome evidence can falsely pass an incomplete report or fail a complete report.
3. Authority: **captain-ruling[2026-09-18]** — “both yes” requires a committed task-body reason for missing review evidence; an incidental filename plus an unrelated absence is not that reason.
4. Trigger evidence: the detached exact-candidate Git-backed exercise calls the production helper with unchanged held gate/status and absent source. Baseline passes. Both cases below produce the opposite result from the required semantic expectation; see `audit.log`, `audit.json` and `detached-test-source.txt`.

| Committed body addition | Intended result | Actual helper result |
|---|---|---|
| `Review held: selected/reviewer-source.txt is unavailable.` | accept | accept |
| `Artifact inventory: selected/reviewer-source.txt.` followed by a separate paragraph `Implementation report: an unrelated optional screenshot is missing. Ready for re-review.` | reject | accept |
| `Review is on hold because selected/reviewer-source.txt has not been provided.` | accept | reject: committed body lacks missing review source reason |

The helper combines a filename search anywhere in the body with an independent body-wide regex (`missing|unavailable|absent|awaiting|not available|not present|not supplied`). The two matches need not describe the same evidence. Its closed alternatives also omit the tested ordinary “has not been provided” wording. Actual source absence and unchanged authority do not cure the false positive: they prove the hold state, not that the task body records its reason.

Smallest correction proposal: keep the existing helper and Git-backed tests; associate the named review source and its absence/awaiting meaning within the same reason unit while permitting line wrapping, and admit the demonstrated ordinary provided phrasing. Preserve ordinary prose and existing positives; do not mandate one exact sentence or build a general NLP parser, ledger or controller. Add these demonstrated negative/positive cases to the existing owner. Choice of a bounded textual reason unit belongs in the producer proposal; this review does not claim arbitrary prose can be fully interpreted.

## Scope and retained evidence

The candidate changes exactly two files with 77 insertions / 13 deletions (+64 net): `claude_runtime_helpers_test.go` and `claude_live_runner_test.go`. Fixture instructions now explicitly require the committed entity-body reason. The caller reads `HEAD:recorded-gate-task/index.md`; actual source absence, document equality and validation status are checked. Existing clean-tree and corrected-plan/frozen-input guards remain. Missing source does not impose a new independent reviewer-dispatch obligation.

Owned Git-backed tests retain absent reason, uncommitted reason, fabricated source, changed gate decision and stage advancement negatives, plus baseline/wrapped positives. Their normal 32.098s and race 30.968s results and registry 1.487s result are retained under `artifacts/implementation/committed-hold-reason`; these green suites were inspected and reused, not rerun. They do not exercise H1's mismatched semantic cases.

Independent exercise: archive exact candidate into a temporary detached directory; add only a temporary test using existing `writeSameStageRevision`, real Git commit, `gates.Read`, and actual `assertSameStageReviewHold`. Run `go test -tags live ./internal/ensigncycle -run '^TestReviewHoldIndependentSemantics$' -count=1 -v`. Actual test process exits1, package duration7.723s. No model/live driver is invoked. Temporary checkout is removed; source retained as evidence, not a standing harness. Candidate and its HEAD remain untouched.

## AC evidence limits

- AC-1: selected-plan/no-invented-reviewer historical native evidence is unchanged; no new native outcome claimed.
- AC-2: captain clarified durable reason requirement; candidate strengthens committed-state/source/authority evidence, but H1 means the reason itself is not yet reliably established. Existing round and independent identity checks remain intact.
- AC-3: all six variants remain routine CI obligations; no selector, skip, XFAIL or targeted-only exemption changed. Current normal/race checks are focused, not final combined acceptance.

Final restacked combined normal/race and gofmt remain deferred to FO. Earlier combined7c63 checks completed with only the known resolver failure; they do not cover this new candidate. Actual final-tip six-variant native execution per host remains pending. Other Claude/Pi failures, including the Pi missing wrapper receipt/incomplete outcomes and naming/recovery work, are not waived or repaired here. The earlier recording-location ambiguity is superseded by the captain ruling, while the observed safe hold remains historical fact.

No candidate/frontmatter changes, models, push, rebase, CI or agents. New finding awaits FO disposition before any repair or corrective rerun.


### H1 FO disposition — 2026-09-18

FO separately authorized **FIX**: Material task-owned evidence defect under captain-ruling2026-09-18. The existing producer must relate the named source to its absence reason within one note unit, support ordinary “not been provided” wording and wrapping, and preserve durable-state controls. No NLP framework or mandatory exact sentence. Reviewer owns evidence only and remains available for correction recheck.
