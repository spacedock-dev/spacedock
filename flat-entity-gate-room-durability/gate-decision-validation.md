# Validation decision: flat gate-room durability

- Recommendation: **PASSED**
- Recommended captain decision: **approve**
- Candidate: `e4cd8afec`
- Decision effect: Approval makes the candidate eligible to stack above PR #779 and unblock `wy`.

## Outcome evidence

- **AC-1 passed.** One supported state commit makes a flat round durable on a fresh host. Exact replay changes no bytes and excludes sibling dirt.
- **AC-2 passed.** New gate and round refs use canonical `@review/...`. Flat and folder entities resolve each ref to the same ticket review home.
- **AC-3 passed.** Frozen flat and folder refs keep their meaning. Legacy round refs, Briefing bytes, digests, and `git-root://` identities stay unchanged.
- **AC-4 passed.** Flat and folder rounds publish and replay identically. Publication and replay enforce Artifact containment inside the ticket home.
- **AC-5 passed.** All malformed reserved refs refuse without an authority spend. Each refusal leaves the state tree and gate-lock state unchanged.

## Correction history

Validation cycle 1 rejected `add2d0632`. Prepared-room classification suppressed a malformed canonical-ref error, and gate close spent authority.

Correction `e4cd8afec` propagates that error through retained-authority validation. It also adds `TestMalformedCanonicalRoomRefCannotSpendGateAuthority`.

## Checks and scope

- The focused gates, status, and two-host checks passed.
- `go test ./...` passed.
- `go test ./... -race` passed.
- The checks for task-touched formatting and diff cleanliness passed.
- The final surface has 423 insertions and 173 deletions. This is +250 net lines across 13 files.
- The estimate was +180 net lines across 12 files, with tolerance ±70 lines and ±2 files.
- The deviation is +70 lines and +1 file. This result is within the inclusive tolerance.

The review found no remaining material findings.
