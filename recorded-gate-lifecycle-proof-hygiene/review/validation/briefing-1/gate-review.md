# Validation gate: recorded-gate lifecycle proof hygiene

## Capability delivered

The recorded-gate lifecycle is now judged by executable behavior rather than
shipped-prose self-checks. The structural PR-view allow-list retains its load-bearing
positive and negative controls, and both Codex and Claude guardrail runners diagnose
an archived held-gate entity before attempting to read its active path.

## Exact candidate

Candidate `1fa8ead5acf643bd00233068a1b57130d6b14081` is based on
`4ff98d8cd97ebcf17b6a583070ce69234e24fc87`. It changes exactly four existing test
files, adds 13 lines, removes 63, and changes no product, instruction, command,
provider, compatibility, or lifecycle-AC surface.

## Validation evidence

- Exact-name and diff controls prove the tautological command-text mutant,
  orphaned parser, and shipped-prose map are gone.
- Removing the allowed `gh pr view` tokens fails the structural positive control;
  planting the token outside the allow-list fails the negative control.
- Forced post-run archival makes both Codex and Claude paths fail first with the
  named archive diagnostic.
- Focused lifecycle and contractlint suites, live-tag compilation, gofmt,
  `go test ./...`, and `go test ./... -race` passed.
- Roborev job 2331 passed 2/2 with no findings; independent validation recommends
  PASSED on all four acceptance criteria.

## Deferred evidence defect

Three credentialed journeys clear the new archive assertion but later encounter
pre-existing held-state oracles. Base and candidate are byte-identical at those
sites: canonical v1 encodes open by Resolution absence while an old live oracle
demands explicit `state: open`. This is actionable live-harness debt, but it is not
a regression or obligation of this deletion-first cleanup.

## Recommendation

Approve the exact candidate for merge and terminalization.

## Decision

Approve to merge candidate `1fa8ead5`; revise to return a material finding to
implementation; or hold before merge.
