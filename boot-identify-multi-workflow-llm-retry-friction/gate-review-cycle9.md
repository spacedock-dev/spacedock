# Validation gate: Claude invocation ledger after launcher pinning

## Capability

The multi-workflow and filing live journeys now grade actual `spacedock` argv on Claude after the front door repins `SPACEDOCK_BIN`. Narrated or quoted commands remain unable to satisfy the ledger.

## Reviewed snapshot

- Candidate: `f982e88b656a8966e8739521d7a075e4d9c90a6b`
- Delta for this correction: two test files, `+24/-7`
- Existing PR: `#584`

## Evidence

- Direct regression simulates front-door launcher repinning, starts a real Bash command through `$SPACEDOCK_BIN`, and observes exactly one NUL-delimited argv record.
- Focused live-tagged Claude ledger, filing, multi-workflow, Bash, and zsh tests pass.
- `go test ./...`, `go test ./... -race`, formatting, and diff checks pass.
- Independent validation reproduced the boundary and adversarial variants without changing candidate bytes.
- Roborev job 454 reported no issues.

## Findings and condition

No Material candidate finding remains. Local supported-model execution hit Anthropic API 429 before First Officer work, so this gate does not claim local Claude model-live success. Exact-head Claude Sonnet and Opus CI remain mandatory before merge.

## Recommendation

Approve the candidate for PR update and required CI. Do not merge unless both Claude lanes confirm nonempty execution ledgers and the unrelated 26n-owned headless-gate dependency is resolved.
