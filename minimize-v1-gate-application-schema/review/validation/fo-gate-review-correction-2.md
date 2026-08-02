# Validation gate review — correction cycle 2

## Recommendation

HOLD for Captain scope authorization. The correction closes AC-1, AC-2, and AC-3 with independent evidence, but the cumulative candidate surface is 15 files and +344 insertions (excluding the unrelated README contamination), above the approved AC-4 ceiling of 14 files and +220 insertions.

## Evidence

- AC-1: decoded YAML-node assertions prove approval applications contain exactly `{target-stage,state}` and hold/revise contain no application; the detached `policy` leaf fails.
- AC-2: the checked-in manifest iterates 16 active and 15 archived paths with strict read/validation; detached `action: advance` state fails the iterator; state is clean at `812b7a47e` plus report commits.
- AC-3: the real recorded CLI asserts `application=advance/consumed`; the detached wrong-action mutation fails.
- AC-4 checks are green (`gofmt`, focused suites, `go test ./... -count=1`, and `go test ./... -race`), but the report discloses original implementation 14 files/+138/-177 plus correction evidence 3 files/+207/-2, yielding the cumulative 15/+344 surface.

## Science Officer advisory

HOLD for Captain scope authorization. Choose one explicit direction: accept the evidence-only cumulative surface of 15 files/+344, or send NTH back to compress below 14 files/+220. Do not treat each correction round as a fresh allowance.

## Decision effect

Approve only with explicit Captain authorization for the revised cumulative tolerance; otherwise revise to implementation for compression. WJ remains held until NTH lands.
