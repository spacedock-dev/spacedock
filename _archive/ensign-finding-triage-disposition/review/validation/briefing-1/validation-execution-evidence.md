# Validation execution evidence

## Routed live replay

- Case A, correct-but-disproportionate: the worker emitted a linked decline Annotation and an advisory Resolution, changed zero product lines, and left entity status at `implementation`.
- Case B, material control: the worker fixed `product/status.txt` with a `+1/0` product diff and likewise did not advance entity status through the advisory record.
- This distinguishes a present all-declines outcome from no finding and reproduces the value contrast against the archived nonzero dutiful-fix and prose-only baselines.

## Classification attacks

- The checked-in offline fixture/check produced exactly 8 ACCEPT and 2 REJECT outcomes.
- Material-as-declined and unknown-class red controls were rejected.
- Five independent mutations—removing user/workflow effect, harm, promise boundary, trigger, or the class allowlist—each failed the paired control.

## Detached audit and repository checks

- Detached checkout at exact `e85eb0cfcc3c243fd94754be2baafa23be302a21`: CLEAN, no material findings.
- `docs/specs/check-finding-triage-materiality.sh`: passed.
- `gofmt -w ./cmd ./internal`: no diff.
- `go test ./...`: passed.
- `go test ./... -race`: passed.
- `git diff --check`: passed.
- Candidate worktree: clean at exact head.

Deferred audit observations and their promotion conditions are preserved in `gate-review.md` and the exact validator report.
