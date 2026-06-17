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

## Stage Report: implementation

- DONE: The two prose-grep guards are honest and doc_test.go:11-compliant; resolution JUSTIFIED.
  Chose **remove + report owed test** for BOTH guards (commit 8b60c568). Justification against doc_test.go:11 + the README "Bad" rule: the guards' targets are MEANINGS (the universal doctrine; the startup-gate guidance), not literals. Proven by probe — appending a paraphrase ("Favor an enforceable code gate rather than a prose-only rule…") to development.md left `TestUniversalDoctrineHasSingleSource` GREEN, i.e. it missed exactly the drift its docstring claimed to catch. A literal-phrase `strings.Contains` used as a proxy for whether a meaning is present IS the banned prose-grep. Genuine-structural does not honestly hold: no token here IS the doctrine (contrast the CLEAN sibling `TestDispatchCoreHasNoClaudeTeamImperative` — `spawn-standing-all`/`--team {team_name}` are command literals that a meaning-inverting paraphrase necessarily drops; the token IS the thing). Honest-narrow to verbatim-absence stays a prose-phrase-absence standing in for the meaning — still the banned grep. So per doc_test.go:11 ("delete the read and report the owed test") I deleted both reads and recorded the OWED PROOF at each removal site: the human gate-review of the commission templates / FO contract (the same review that already owns the qualitative lede judgement in `TestTemplatesLeadWithOutcome`).
- DONE: AC-5 install-script comment-out hole closed via executableShellCommands().
  `TestLiveWorkflowInstallsPinnedGotestsum` now greps `executableShellCommands(script)` (the sibling helper) instead of raw script text. Probe — comment out every VERSION=/sha=/shasum-verify line in install-gotestsum.sh — flips the test RED (was GREEN before the fix). Commit 8b60c568.
- DONE: go build/vet/test green + DETACHED adversarial re-audit confirms CLEAN.
  Worktree: build/vet exit 0; `go test ./...` exit 0 (all packages ok, incl. skills/integration; no FAIL/panic). Detached throwaway checkout (`git worktree add --detach HEAD`): independent build/vet/test green; the two removed guards confirmed ABSENT; paraphrase probe leaves the suite honestly GREEN (no guard falsely claims to catch it); AC-5 comment-out probe REDs; cut-surface scan found the retired `"Prefer a code gate…"` phrase ONLY in a `//` comment, no active prose-grep reintroduced. All probes reverted; detached worktree removed.

### Summary
Both #388/#391 anti-drift guards were literal-phrase `strings.Contains` greps standing in for a MEANING (doctrine / startup-gate guidance) — a paraphrase carried the meaning while dropping the bytes, so each stayed green on the drift its docstring advertised (the banned prose-grep per doc_test.go:11). Since no token IS the doctrine, I retired both per the package's own "delete the read and report the owed test" policy, recording the owed human gate-review at each site; the legitimate structural halves (`TestTemplatesCarryWorkflowSpecificRulesSlot`, `TestTemplatesLeadWithOutcome`) survive. AC-5 now filters the install script through `executableShellCommands()` so a commented-out verification REDs. Out-of-scope note for validation: `codex_foreground_wait_shape_test.go` carries pre-existing phrase checks (commit ab39d5d8, #378) not named by this audit and not touched here. The shipped BEHAVIOR was unchanged — this was hollow test-strength, now honest.
