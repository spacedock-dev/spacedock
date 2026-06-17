---
title: Harden the #388/#391 contractlint anti-drift guards — prose-grep proofs to genuine checks (+ AC-5 install-script comment-out hole)
status: implementation
source: "Pre-cut antipattern audit for v0.20.4 (2026-06-17, 0204 Commander), confirmed from source. Two contractlint tests are literal-phrase strings.Contains greps that overclaim to be behavioral drift guards and VIOLATE the package's own policy (internal/contractlint/doc_test.go:11: 'Do not add prose-grep checks here ... If behavior matters, test it by running the behavior'). The shipped BEHAVIOR is correct (zero doctrine in templates today); the defect is hollow test-strength on the shipped contract surface. Captain decision: fix-first before the cut."
sprint: 0204-structured-reads
id: xhttp6bs3h6w8afy4xf2zmz8
worktree: .worktrees/spacedock-ensign-contractlint-antidrift-guard-hardening
started: 2026-06-17T07:26:25Z
---

Two pre-cut-audit blockers + one non-blocking hole, all on the v0.20.4 contract/CI surface. Make the checks honest and policy-compliant; the shipped behavior is already correct, so this is test-strength, not a behavior fix.

## Blockers (confirmed from source)
1. **`TestUniversalDoctrineHasSingleSource`** (`internal/contractlint/template_defer_test.go`): `strings.Contains(b, "Prefer a code gate over a prose-only rule")` over a hardcoded literal phrase. A PARAPHRASED doctrine restatement in a template stays GREEN; only a byte-identical copy REDs. Its docstring overclaims ("fails when the doctrine drifts back into a template", "not a presence tautology"). It is the #388 anti-drift proof (templates defer doctrine to the contract), so the deliverable's guarantee is proven by a prose-grep.
2. **`TestStartupGateGuidanceHasSingleSource`** (`internal/contractlint/structural_checks_test.go:295-338`): same shape — `markers := []string{"Contract version gate","per-class remedy","spacedock doctor"}` then `strings.Contains`. Paraphrase slips past GREEN; docstring claims "defect (not a prose property)".

Both violate `doc_test.go:11`. Note: detecting MEANING-drift by machine without interpreting prose is likely impossible — so the honest resolution is a real design call (see below), not a mechanical rewrite.

## Non-blocking hole (fold in — cheap)
3. **`TestLiveWorkflowInstallsPinnedGotestsum`** (`internal/release/cilog_clean_output_workflow_test.go:59-92`): asserts the install script sha256-verifies the tarball via raw `strings.Contains` over the whole script text — defeated by COMMENTING OUT every verify/pin line (inert script stays GREEN). The package already owns the fix: sibling `journey_workflow_test.go` guards use `executableShellCommands()` to defeat exactly this comment-out attack (e.g. `TestReleaseWorkflowGuardRejectsCommentOnlyJourneyCostBuilder`). Apply that helper so a commented-out verification REDs.

## The design call (justify the choice against doc_test.go:11 + the README "Bad" rule)
For blockers 1 & 2, pick and JUSTIFY one resolution per check (the validation gate votes on it):
- **Genuine structural dedup** — reframe to a token/structural property that IS the thing (not a prose proxy for meaning), with a control that REDs on the realistic regression (a paraphrase restatement). Only valid if such a structural property genuinely exists.
- **Honest narrow + owed test** — narrow the check to what it actually verifies (verbatim-phrase absence), strip the overclaiming docstrings, and FILE/identify the owed behavior-or-review for the meaning-drift gap. Must still satisfy doc_test.go:11 (a prose-phrase-absence used as a meaning proxy is itself the banned prose-grep — so this likely means moving the check OUT of the quarantine package or deleting it with the owed test recorded).
- **Remove + report the owed test** — if neither above honestly holds, delete the prose-grep per the package's own policy (doc_test.go:11: "delete the read and report the owed test").

## Acceptance bar (validation = detached re-audit)
- A paraphrased doctrine/guidance restatement in a template no longer passes a guard that CLAIMS to catch drift (either it REDs, or no guard claims to catch it and the docstring is honest + the owed test recorded).
- The install-script comment-out attack (comment every verify/pin line) now REDs (`executableShellCommands()` applied).
- `internal/contractlint/doc_test.go:11` is honored (no prose-grep-as-behavior-substitute remains in the package).
- `go build ./... && go vet ./... && go test ./...` green; the chosen resolution justified in the stage report; a detached adversarial re-audit of the cut surface finds it CLEAN.
