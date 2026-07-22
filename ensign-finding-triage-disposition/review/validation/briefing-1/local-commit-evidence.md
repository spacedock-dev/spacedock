# Candidate commit evidence

- Base: `fa240a76cd67fc0ea2552901824722ca8bfa1c73`
- Candidate: `e85eb0cfcc3c243fd94754be2baafa23be302a21`
- Candidate worktree porcelain at package assembly: empty

Commit sequence:

```text
d3efaea01ce03fe497403015901f6cd77c3f71b0  Specify advisory finding triage dispositions
059b12d8f07abda150935649266f6ebbc05f4642  Rebaseline prompt load for triage trigger
9b2093b5d4590afc2f42df1e6893978264a5513d  Make triage materiality fixture conjuncts load-bearing
e85eb0cfcc3c243fd94754be2baafa23be302a21  Reject unknown triage classes
```

Exact base-to-candidate name-status:

```text
M	docs/dev/README.md
A	docs/specs/check-finding-triage-materiality.sh
M	docs/specs/gate-resolution-frontmatter-contract.md
A	docs/specs/testdata/finding-triage-materiality.tsv
M	internal/contractlint/fo_function_reference_invariant_test.go
M	skills/feedback-rejection-flow/SKILL.md
```

Exact base-to-candidate diff stat:

```text
 docs/dev/README.md                                 |  1 +
 docs/specs/check-finding-triage-materiality.sh     | 33 ++++++++++++++++
 docs/specs/gate-resolution-frontmatter-contract.md | 45 ++++++++++++++++++++++
 docs/specs/testdata/finding-triage-materiality.tsv | 11 ++++++
 .../fo_function_reference_invariant_test.go        |  6 +--
 skills/feedback-rejection-flow/SKILL.md            |  2 +
 6 files changed, 95 insertions(+), 3 deletions(-)
```

The exact live replay, adversarial audit, focused check, full/race suites, formatting, and cleanliness outcomes are retained byte-for-byte in `validator-report.md`.

