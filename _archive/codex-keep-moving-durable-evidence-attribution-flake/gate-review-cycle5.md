# Validation gate: 8b keep-moving durable attribution

## Capability

The provider-neutral Git oracle now recognizes both supported engagement histories exposed by PR #585:

- Codex may atomically persist `started`, the implementation Stage Report, and only files newly added beneath that revision's exact-slug gate-room binding.
- Sonnet may persist a known-ticket dispatch batch containing a valid subset of expected tickets plus the exact nonterminal questioned rework; every omitted expected ticket still needs its own independently valid dispatch.

The change does not restore transcript/provider observers, accept arbitrary ticket-owned paths, or weaken report, terminal, archive, ordering, and questioned-history attribution.

## Reviewed snapshot

- Candidate: `7f109271c`
- Incremental correction: 2 existing files, `+111/-15`
- Cumulative branch surface: 10 files, `+971/-1,449` (net `-478`)

## Validation evidence

- 33 per-ticket journey cases and 5 batch cases pass.
- Removing the atomic-room predicate loses only `ready-one`; removing the split-batch predicate loses only `ready-one` and `ready-two`.
- Omitted-ticket partial remains `2/3`; prebound, replaced, modified, outside, and slug-prefix room variants remain red.
- Arbitrary questioned scope, foreign paths, premature terminal fields, archive attribution failures, and questioned terminal/reopen histories remain red.
- Four retained supported live roots each regrade green.
- `gofmt -l` is empty; `git diff --check`, focused tests, `go test ./...`, and `go test ./... -race` pass.

## Findings

No Material or Needs-decision finding remains. Roborev and another live-model run were intentionally not repeated: this correction deterministically encodes the two exact failing CI histories and preserves their adjacent negative controls.

## Recommendation

Approve exact candidate `7f109271c` for PR update. The next CI run must still be based on a main branch containing 26n's separately owned headless gate-stop correction before it can be expected to go fully green.
